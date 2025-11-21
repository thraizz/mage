package rules

import (
	"fmt"
	"sync"
)

// SpecialActionType represents types of special actions
// Per Rule 116: Special actions don't use the stack
type SpecialActionType string

const (
	// SpecialActionPlayLand plays a land (Rule 116.2a)
	// Can only be taken during main phase with empty stack
	SpecialActionPlayLand SpecialActionType = "PLAY_LAND"

	// SpecialActionTurnFaceUp turns a face-down creature face up (Rule 116.2b)
	// Can be taken any time player has priority
	SpecialActionTurnFaceUp SpecialActionType = "TURN_FACE_UP"

	// SpecialActionEndEffect ends a continuous effect (Rule 116.2c)
	// Can be taken any time player has priority
	SpecialActionEndEffect SpecialActionType = "END_EFFECT"

	// SpecialActionIgnoreStaticAbility ignores a static ability (Rule 116.2d)
	// Can be taken any time player has priority
	SpecialActionIgnoreStaticAbility SpecialActionType = "IGNORE_STATIC_ABILITY"

	// SpecialActionDiscardCirclingVultures discards Circling Vultures (Rule 116.2e)
	// Can be taken any time player could cast an instant
	SpecialActionDiscardCirclingVultures SpecialActionType = "DISCARD_CIRCLING_VULTURES"

	// SpecialActionSuspend exiles a card with suspend (Rule 116.2f)
	// Can be taken any time player has priority and could cast the card
	SpecialActionSuspend SpecialActionType = "SUSPEND"

	// SpecialActionCompanion pays {3} to put companion into hand (Rule 116.2g)
	// Can only be taken during main phase with empty stack, once per game
	SpecialActionCompanion SpecialActionType = "COMPANION"

	// SpecialActionForetell pays {2} to exile with foretell (Rule 116.2h)
	// Can be taken any time player has priority during their turn
	SpecialActionForetell SpecialActionType = "FORETELL"

	// SpecialActionPlot exiles a card with plot (Rule 116.2k)
	// Can be taken any time player has priority during their turn with empty stack
	SpecialActionPlot SpecialActionType = "PLOT"

	// SpecialActionUnlock pays unlock cost for a locked half (Rule 116.2m)
	// Can only be taken during main phase with empty stack
	SpecialActionUnlock SpecialActionType = "UNLOCK"
)

// SpecialAction represents a special action that can be taken
type SpecialAction struct {
	Type         SpecialActionType
	PlayerID     string
	SourceID     string // Card/permanent involved
	Description  string
	Execute      func() error
	CanTake      func() bool // Additional legality check
}

// SpecialActionRestriction defines when a special action can be taken
type SpecialActionRestriction struct {
	RequiresMainPhase   bool // Must be main phase
	RequiresEmptyStack  bool // Stack must be empty
	RequiresOwnTurn     bool // Must be own turn
	RequiresPriority    bool // Must have priority (all special actions require this)
	OncePerGame         bool // Can only be taken once per game
	AdditionalCheck     func() bool // Additional custom checks
}

// GetRestrictions returns the restrictions for a special action type
func GetRestrictions(actionType SpecialActionType) SpecialActionRestriction {
	switch actionType {
	case SpecialActionPlayLand:
		// Rule 116.2a: main phase, empty stack, once per turn
		return SpecialActionRestriction{
			RequiresMainPhase:  true,
			RequiresEmptyStack: true,
			RequiresPriority:   true,
		}

	case SpecialActionTurnFaceUp:
		// Rule 116.2b: any time with priority
		return SpecialActionRestriction{
			RequiresPriority: true,
		}

	case SpecialActionEndEffect:
		// Rule 116.2c: any time with priority
		return SpecialActionRestriction{
			RequiresPriority: true,
		}

	case SpecialActionIgnoreStaticAbility:
		// Rule 116.2d: any time with priority
		return SpecialActionRestriction{
			RequiresPriority: true,
		}

	case SpecialActionDiscardCirclingVultures:
		// Rule 116.2e: any time could cast instant (has priority)
		return SpecialActionRestriction{
			RequiresPriority: true,
		}

	case SpecialActionSuspend:
		// Rule 116.2f: any time with priority, if could begin casting
		return SpecialActionRestriction{
			RequiresPriority: true,
		}

	case SpecialActionCompanion:
		// Rule 116.2g: main phase, empty stack, once per game
		return SpecialActionRestriction{
			RequiresMainPhase:  true,
			RequiresEmptyStack: true,
			RequiresPriority:   true,
			OncePerGame:        true,
		}

	case SpecialActionForetell:
		// Rule 116.2h: any time with priority during own turn
		return SpecialActionRestriction{
			RequiresPriority: true,
			RequiresOwnTurn:  true,
		}

	case SpecialActionPlot:
		// Rule 116.2k: own turn, empty stack, with priority
		return SpecialActionRestriction{
			RequiresEmptyStack: true,
			RequiresPriority:   true,
			RequiresOwnTurn:    true,
		}

	case SpecialActionUnlock:
		// Rule 116.2m: main phase, empty stack, with priority
		return SpecialActionRestriction{
			RequiresMainPhase:  true,
			RequiresEmptyStack: true,
			RequiresPriority:   true,
		}

	default:
		// Default: requires priority only
		return SpecialActionRestriction{
			RequiresPriority: true,
		}
	}
}

// SpecialActionManager manages special actions
type SpecialActionManager struct {
	mu                  sync.RWMutex
	takenThisGame       map[string]map[SpecialActionType]bool // playerID -> actionType -> taken
	takenThisTurn       map[string]map[SpecialActionType]int  // playerID -> actionType -> count
	availableActions    []SpecialAction
	canTakeDuringResolve bool
}

// NewSpecialActionManager creates a new special action manager
func NewSpecialActionManager() *SpecialActionManager {
	return &SpecialActionManager{
		takenThisGame:       make(map[string]map[SpecialActionType]bool),
		takenThisTurn:       make(map[string]map[SpecialActionType]int),
		availableActions:    make([]SpecialAction, 0, 8),
		canTakeDuringResolve: false, // Generally can't take during resolution
	}
}

// CanTakeAction checks if a special action can be taken
func (sam *SpecialActionManager) CanTakeAction(
	action SpecialAction,
	hasPriority bool,
	isMainPhase bool,
	isEmptyStack bool,
	isOwnTurn bool,
) bool {
	restrictions := GetRestrictions(action.Type)

	// Check basic restrictions
	if restrictions.RequiresPriority && !hasPriority {
		return false
	}
	if restrictions.RequiresMainPhase && !isMainPhase {
		return false
	}
	if restrictions.RequiresEmptyStack && !isEmptyStack {
		return false
	}
	if restrictions.RequiresOwnTurn && !isOwnTurn {
		return false
	}

	// Check once-per-game restriction
	if restrictions.OncePerGame {
		sam.mu.RLock()
		if playerActions, exists := sam.takenThisGame[action.PlayerID]; exists {
			if playerActions[action.Type] {
				sam.mu.RUnlock()
				return false
			}
		}
		sam.mu.RUnlock()
	}

	// Check additional custom restrictions
	if action.CanTake != nil && !action.CanTake() {
		return false
	}

	return true
}

// TakeAction executes a special action
// Per Rule 116.3: Player receives priority afterward
func (sam *SpecialActionManager) TakeAction(action SpecialAction) error {
	sam.mu.Lock()
	defer sam.mu.Unlock()

	// Execute the action
	if err := action.Execute(); err != nil {
		return fmt.Errorf("failed to execute special action %s: %w", action.Type, err)
	}

	// Track that action was taken
	restrictions := GetRestrictions(action.Type)
	if restrictions.OncePerGame {
		if sam.takenThisGame[action.PlayerID] == nil {
			sam.takenThisGame[action.PlayerID] = make(map[SpecialActionType]bool)
		}
		sam.takenThisGame[action.PlayerID][action.Type] = true
	}

	// Track per-turn count
	if sam.takenThisTurn[action.PlayerID] == nil {
		sam.takenThisTurn[action.PlayerID] = make(map[SpecialActionType]int)
	}
	sam.takenThisTurn[action.PlayerID][action.Type]++

	return nil
}

// GetActionsTakenThisTurn returns the count of a specific action taken this turn
func (sam *SpecialActionManager) GetActionsTakenThisTurn(playerID string, actionType SpecialActionType) int {
	sam.mu.RLock()
	defer sam.mu.RUnlock()

	if playerActions, exists := sam.takenThisTurn[playerID]; exists {
		return playerActions[actionType]
	}
	return 0
}

// ResetTurn clears per-turn tracking
func (sam *SpecialActionManager) ResetTurn() {
	sam.mu.Lock()
	defer sam.mu.Unlock()
	sam.takenThisTurn = make(map[string]map[SpecialActionType]int)
}

// SetCanTakeDuringResolve sets whether special actions can be taken during resolution
func (sam *SpecialActionManager) SetCanTakeDuringResolve(can bool) {
	sam.mu.Lock()
	defer sam.mu.Unlock()
	sam.canTakeDuringResolve = can
}

// CanTakeDuringResolve returns whether special actions can be taken during resolution
func (sam *SpecialActionManager) CanTakeDuringResolve() bool {
	sam.mu.RLock()
	defer sam.mu.RUnlock()
	return sam.canTakeDuringResolve
}

// Reset clears all state
func (sam *SpecialActionManager) Reset() {
	sam.mu.Lock()
	defer sam.mu.Unlock()
	sam.takenThisGame = make(map[string]map[SpecialActionType]bool)
	sam.takenThisTurn = make(map[string]map[SpecialActionType]int)
	sam.availableActions = sam.availableActions[:0]
	sam.canTakeDuringResolve = false
}
