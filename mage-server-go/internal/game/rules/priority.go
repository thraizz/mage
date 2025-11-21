package rules

import (
	"fmt"
	"sync"
)

// ResolutionContext tracks what spell/ability is currently resolving
// and allows nested resolution (e.g., casting copies during resolution)
type ResolutionContext struct {
	mu                sync.RWMutex
	resolvingStack    []string // Stack of resolving item IDs (innermost at end)
	depth             int      // Current resolution depth
	maxDepth          int      // Maximum allowed depth (prevent infinite recursion)
	allowManaAbilities bool     // Whether mana abilities can be activated
	allowSpecialActions bool    // Whether special actions can be taken
}

// NewResolutionContext creates a new resolution context
func NewResolutionContext() *ResolutionContext {
	return &ResolutionContext{
		resolvingStack: make([]string, 0, 8),
		depth:          0,
		maxDepth:       10, // Prevent infinite recursion
		allowManaAbilities: false,
		allowSpecialActions: false,
	}
}

// BeginResolution marks the start of resolving a stack item
func (rc *ResolutionContext) BeginResolution(itemID string) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.depth >= rc.maxDepth {
		return fmt.Errorf("maximum resolution depth (%d) exceeded", rc.maxDepth)
	}

	rc.resolvingStack = append(rc.resolvingStack, itemID)
	rc.depth++
	return nil
}

// EndResolution marks the end of resolving a stack item
func (rc *ResolutionContext) EndResolution(itemID string) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.depth == 0 {
		return fmt.Errorf("no item currently resolving")
	}

	// Verify we're ending the correct resolution
	if len(rc.resolvingStack) > 0 {
		current := rc.resolvingStack[len(rc.resolvingStack)-1]
		if current != itemID {
			return fmt.Errorf("resolution mismatch: expected %s, got %s", current, itemID)
		}
		rc.resolvingStack = rc.resolvingStack[:len(rc.resolvingStack)-1]
	}

	rc.depth--
	return nil
}

// IsResolving returns true if something is currently resolving
func (rc *ResolutionContext) IsResolving() bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.depth > 0
}

// GetCurrentResolvingID returns the ID of the currently resolving item (innermost)
func (rc *ResolutionContext) GetCurrentResolvingID() string {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if len(rc.resolvingStack) == 0 {
		return ""
	}
	return rc.resolvingStack[len(rc.resolvingStack)-1]
}

// GetDepth returns the current resolution depth
func (rc *ResolutionContext) GetDepth() int {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.depth
}

// SetAllowManaAbilities enables/disables mana ability activation during resolution
func (rc *ResolutionContext) SetAllowManaAbilities(allow bool) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.allowManaAbilities = allow
}

// CanActivateManaAbilities returns true if mana abilities can be activated
func (rc *ResolutionContext) CanActivateManaAbilities() bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.allowManaAbilities
}

// SetAllowSpecialActions enables/disables special actions during resolution
func (rc *ResolutionContext) SetAllowSpecialActions(allow bool) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.allowSpecialActions = allow
}

// CanTakeSpecialActions returns true if special actions can be taken
func (rc *ResolutionContext) CanTakeSpecialActions() bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.allowSpecialActions
}

// Reset clears all resolution state
func (rc *ResolutionContext) Reset() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.resolvingStack = rc.resolvingStack[:0]
	rc.depth = 0
	rc.allowManaAbilities = false
	rc.allowSpecialActions = false
}

// PriorityWindow represents a window during resolution where players can take actions
type PriorityWindow struct {
	Type        PriorityWindowType
	PlayerID    string // Player who has the window
	Context     string // Description of what's happening
	AllowedActions []ActionType // What actions are allowed in this window
}

// PriorityWindowType represents the type of priority window
type PriorityWindowType string

const (
	// PriorityWindowManaPayment allows mana abilities during cost payment
	PriorityWindowManaPayment PriorityWindowType = "MANA_PAYMENT"
	// PriorityWindowChoice allows choices during resolution
	PriorityWindowChoice PriorityWindowType = "CHOICE"
	// PriorityWindowTarget allows target selection during resolution
	PriorityWindowTarget PriorityWindowType = "TARGET"
	// PriorityWindowSpecialAction allows special actions
	PriorityWindowSpecialAction PriorityWindowType = "SPECIAL_ACTION"
	// PriorityWindowNestedCast allows casting spells during resolution
	PriorityWindowNestedCast PriorityWindowType = "NESTED_CAST"
)

// ActionType represents types of actions players can take
type ActionType string

const (
	// ActionActivateMana activates a mana ability
	ActionActivateMana ActionType = "ACTIVATE_MANA"
	// ActionSpecialAction takes a special action
	ActionSpecialAction ActionType = "SPECIAL_ACTION"
	// ActionCastSpell casts a spell
	ActionCastSpell ActionType = "CAST_SPELL"
	// ActionActivateAbility activates an ability
	ActionActivateAbility ActionType = "ACTIVATE_ABILITY"
	// ActionMakeChoice makes a choice
	ActionMakeChoice ActionType = "MAKE_CHOICE"
	// ActionSelectTarget selects a target
	ActionSelectTarget ActionType = "SELECT_TARGET"
)

// PriorityWindowManager manages priority windows during resolution
type PriorityWindowManager struct {
	mu            sync.RWMutex
	activeWindow  *PriorityWindow
	windowHistory []PriorityWindow // For debugging/replay
}

// NewPriorityWindowManager creates a new priority window manager
func NewPriorityWindowManager() *PriorityWindowManager {
	return &PriorityWindowManager{
		windowHistory: make([]PriorityWindow, 0, 16),
	}
}

// OpenWindow opens a new priority window
func (pwm *PriorityWindowManager) OpenWindow(window PriorityWindow) error {
	pwm.mu.Lock()
	defer pwm.mu.Unlock()

	if pwm.activeWindow != nil {
		return fmt.Errorf("priority window already open: %v", pwm.activeWindow.Type)
	}

	pwm.activeWindow = &window
	pwm.windowHistory = append(pwm.windowHistory, window)
	return nil
}

// CloseWindow closes the active priority window
func (pwm *PriorityWindowManager) CloseWindow() {
	pwm.mu.Lock()
	defer pwm.mu.Unlock()
	pwm.activeWindow = nil
}

// GetActiveWindow returns the currently active window
func (pwm *PriorityWindowManager) GetActiveWindow() *PriorityWindow {
	pwm.mu.RLock()
	defer pwm.mu.RUnlock()
	return pwm.activeWindow
}

// IsActionAllowed returns true if the given action is allowed in the current window
func (pwm *PriorityWindowManager) IsActionAllowed(action ActionType) bool {
	pwm.mu.RLock()
	defer pwm.mu.RUnlock()

	if pwm.activeWindow == nil {
		return false
	}

	for _, allowed := range pwm.activeWindow.AllowedActions {
		if allowed == action {
			return true
		}
	}
	return false
}

// Reset clears all window state
func (pwm *PriorityWindowManager) Reset() {
	pwm.mu.Lock()
	defer pwm.mu.Unlock()
	pwm.activeWindow = nil
	pwm.windowHistory = pwm.windowHistory[:0]
}
