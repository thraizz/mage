package game

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/magefree/mage-server-go/internal/game/mana"
	"github.com/magefree/mage-server-go/internal/game/rules"
	"github.com/magefree/mage-server-go/internal/game/targeting"
	"github.com/magefree/mage-server-go/internal/plugin"
	"github.com/magefree/mage-server-go/internal/repository"
	"go.uber.org/zap"
)

// Zone constants matching Java implementation
const (
	zoneLibrary     = 0
	zoneHand        = 1
	zoneBattlefield = 2
	zoneGraveyard   = 3
	zoneStack       = 4
	zoneExile       = 5
	zoneCommand     = 6
)

// zoneToString converts a zone constant to a string
func zoneToString(zone int) string {
	switch zone {
	case zoneLibrary:
		return "LIBRARY"
	case zoneHand:
		return "HAND"
	case zoneBattlefield:
		return "BATTLEFIELD"
	case zoneGraveyard:
		return "GRAVEYARD"
	case zoneStack:
		return "STACK"
	case zoneExile:
		return "EXILE"
	case zoneCommand:
		return "COMMAND"
	default:
		return "UNKNOWN"
	}
}

// Context helpers for CDA calculations
// These allow CDAs to access game state through context

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const gameIDContextKey contextKey = "gameID"

// withGameID adds a game ID to the context for CDA calculations
func withGameID(ctx context.Context, gameID string) context.Context {
	return context.WithValue(ctx, gameIDContextKey, gameID)
}

// extractGameIDFromContext retrieves the game ID from context
func extractGameIDFromContext(ctx context.Context) string {
	if gameID, ok := ctx.Value(gameIDContextKey).(string); ok {
		return gameID
	}
	return ""
}

// Ability ID constants matching Java keyword abilities
const (
	abilityFirstStrike              = "FirstStrikeAbility"
	abilityDoubleStrike             = "DoubleStrikeAbility"
	abilityVigilance                = "VigilanceAbility"
	abilityFlying                   = "FlyingAbility"
	abilityReach                    = "ReachAbility"
	abilityTrample                  = "TrampleAbility"
	abilityTrampleOverPlaneswalkers = "TrampleOverPlaneswalkersAbility"
	abilityDeathtouch               = "DeathtouchAbility"
	abilityDefender                 = "DefenderAbility"
	abilityLifelink                 = "LifelinkAbility"
	abilityMenace                   = "MenaceAbility"
	abilityUnblockable              = "CantBeBlockedSourceAbility"
	abilityBanding                  = "BandingAbility"
	abilityHaste                    = "HasteAbility"
)

// PassUntilType represents different auto-pass modes
type PassUntilType int

const (
	PassUntilNone          PassUntilType = iota // No auto-pass
	PassUntilEndOfTurn                          // Pass until end of current turn
	PassUntilNextTurn                           // Pass until start of next turn
	PassUntilStackResolved                      // Pass until current stack resolves
	PassUntilMyNextTurn                         // Pass until player's next upkeep (F6)
)

// EngineGameView represents the complete game state view for a player
type EngineGameView struct {
	GameID         string
	State          GameState
	Phase          string
	Step           string
	Turn           int
	ActivePlayerID string
	PriorityPlayer string
	Players        []EnginePlayerView
	Battlefield    []EngineCardView
	Stack          []EngineCardView
	Exile          []EngineCardView
	Command        []EngineCardView
	Revealed       []EngineRevealedView
	LookedAt       []EngineLookedAtView
	Combat         EngineCombatView
	StartedAt      time.Time
	Messages       []EngineMessage
	Prompts        []EnginePrompt

	// Pre-computed display values (server source of truth)
	ActivePlayerName     string
	PriorityPlayerName   string
	GameFormat           string
	IsMulliganPhase      bool
	LandsPlayedThisTurn  int
	LandsAllowedThisTurn int

	// Pending library search (if any) - only visible to the searching player
	PendingLibrarySearch *EngineLibrarySearchView
}

// EngineLibrarySearchView represents a pending library search for the UI
type EngineLibrarySearchView struct {
	PlayerID    string           // Player who is searching
	Message     string           // Description of what to search for
	Destination string           // Where selected card goes: "hand", "battlefield", "top", "graveyard"
	Cards       []EngineCardView // Cards in the library to choose from
	CanCancel   bool             // Whether the player can cancel the search
}

// EnginePlayerView represents a player's view in the game
type EnginePlayerView struct {
	PlayerID            string
	Name                string
	Life                int
	Poison              int
	Energy              int
	LibraryCount        int
	HandCount           int
	Hand                []EngineCardView
	Graveyard           []EngineCardView
	ManaPool            EngineManaPoolView
	HasPriority         bool
	Passed              bool
	StateOrdinal        int
	Lost                bool
	Left                bool
	Wins                int
	KeptHand            bool // Whether player has kept their hand during mulligan phase
	HasAvailableActions bool // Server-computed: does this player have any legal actions right now?
}

// EngineCardView represents a card in any zone
type EngineCardView struct {
	ID                string
	Name              string
	DisplayName       string
	ManaCost          string
	Type              string
	SubTypes          []string
	SuperTypes        []string
	Color             string
	Power             string
	Toughness         string
	Loyalty           string
	CardNumber        int
	ExpansionSet      string
	Rarity            string
	RulesText         string
	Tapped            bool
	Flipped           bool
	Transformed       bool
	FaceDown          bool
	Zone              int
	ControllerID      string
	OwnerID           string
	AttachedToCard    []string
	Abilities         []EngineAbilityView
	Counters          []EngineCounterView
	AvailableActions  []EngineCardAction // Server-computed available actions
	SummoningSickness bool               // Creature has summoning sickness
}

// EngineCardAction represents an available action for a card
type EngineCardAction struct {
	ActionType     string // "CAST_SPELL", "PLAY_LAND", "ACTIVATE_ABILITY", "ACTIVATE_MANA_ABILITY"
	ActionID       string // For abilities with multiple options
	DisplayText    string // "Cast", "Play Land", "Tap: Add {G}"
	IsEnabled      bool   // Can perform right now?
	DisabledReason string // "Not enough mana", "Wrong phase"
}

// EngineAbilityView represents an ability on a card
type EngineAbilityView struct {
	ID   string
	Text string
	Rule string
}

// EngineCounterView represents counters on a card
type EngineCounterView struct {
	Name  string
	Count int
}

// EngineManaPoolView represents a player's mana pool
type EngineManaPoolView struct {
	White     int
	Blue      int
	Black     int
	Red       int
	Green     int
	Colorless int
}

// EngineRevealedView represents revealed cards
type EngineRevealedView struct {
	Name  string
	Cards []EngineCardView
}

// EngineLookedAtView represents looked-at cards
type EngineLookedAtView struct {
	Name  string
	Cards []EngineCardView
}

// EngineCombatView represents combat state
type EngineCombatView struct {
	AttackingPlayerID string
	Groups            []EngineCombatGroupView
}

// EngineCombatGroupView represents a combat group
type EngineCombatGroupView struct {
	Attackers         []string
	Blockers          []string
	DefenderID        string
	DefendingPlayerID string
	Blocked           bool
}

// combatState tracks all combat-related state for a game
// Per Java Combat class
type combatState struct {
	attackingPlayerID string
	groups            []*combatGroup
	formerGroups      []*combatGroup
	blockingGroups    map[string]*combatGroup // blockerID -> group
	defenders         map[string]bool         // all possible defenders (players, planeswalkers, battles)
	attackers         map[string]bool         // all attacking creatures
	blockers          map[string]bool         // all blocking creatures
	attackersTapped   map[string]bool         // creatures tapped by attack
	firstStrikers     map[string]bool         // creatures that dealt damage in first strike step
	// Combat requirements/restrictions tracking (Java: Combat lines 70-74)
	creaturesForcedToAttack    map[string]map[string]bool // creatureID -> set of defenderIDs it must attack (empty = any)
	creatureMustBlockAttackers map[string]map[string]bool // blockerID -> set of attackerIDs it must block
	maxAttackers               int                        // maximum number of attackers allowed (-1 = no limit)
	minBlockersPerAttacker     map[string]int             // attackerID -> minimum blockers required
	maxBlockersPerAttacker     map[string]int             // attackerID -> maximum blockers allowed
	// Attack tracking for triggers (Java: PlayersAttackedThisTurnWatcher)
	playersAttackedThisTurn                 map[string]map[string]bool // attackingPlayerID -> set of playerIDs attacked
	planeswalkerControllersAttackedThisTurn map[string]map[string]bool // attackingPlayerID -> set of playerIDs whose planeswalkers were attacked
}

// combatGroup represents a single combat group (attackers vs defender + blockers)
// Per Java CombatGroup class
type combatGroup struct {
	defenderID          string         // player, planeswalker, or battle being attacked
	defenderIsPermanent bool           // is defender a permanent (vs player)
	defendingPlayerID   string         // controller of defending permanents
	attackers           []string       // attacking creature IDs
	formerAttackers     []string       // historical attackers (for "attacked this turn")
	blockers            []string       // blocking creature IDs
	blocked             bool           // is this group blocked
	attackerOrder       map[string]int // damage assignment order for attackers (deprecated - kept for compatibility)
	blockerOrder        map[string]int // damage assignment order for blockers (deprecated - kept for compatibility)
	// Modern damage division (Rule 510.1c-d: players divide damage as they choose, no ordering required)
	attackerDamageAssignments map[string]map[string]int // attackerID -> (blockerID -> damage)
	blockerDamageAssignments  map[string]map[string]int // blockerID -> (attackerID -> damage)
}

// newCombatState creates a new combat state
func newCombatState() *combatState {
	return &combatState{
		groups:                                  make([]*combatGroup, 0),
		formerGroups:                            make([]*combatGroup, 0),
		blockingGroups:                          make(map[string]*combatGroup),
		defenders:                               make(map[string]bool),
		attackers:                               make(map[string]bool),
		blockers:                                make(map[string]bool),
		attackersTapped:                         make(map[string]bool),
		firstStrikers:                           make(map[string]bool),
		creaturesForcedToAttack:                 make(map[string]map[string]bool),
		creatureMustBlockAttackers:              make(map[string]map[string]bool),
		maxAttackers:                            -1, // no limit by default
		minBlockersPerAttacker:                  make(map[string]int),
		maxBlockersPerAttacker:                  make(map[string]int),
		playersAttackedThisTurn:                 make(map[string]map[string]bool),
		planeswalkerControllersAttackedThisTurn: make(map[string]map[string]bool),
	}
}

// newCombatGroup creates a new combat group
func newCombatGroup(defenderID string, defenderIsPermanent bool, defendingPlayerID string) *combatGroup {
	return &combatGroup{
		defenderID:                defenderID,
		defenderIsPermanent:       defenderIsPermanent,
		defendingPlayerID:         defendingPlayerID,
		attackers:                 make([]string, 0),
		formerAttackers:           make([]string, 0),
		blockers:                  make([]string, 0),
		blocked:                   false,
		attackerOrder:             make(map[string]int),
		blockerOrder:              make(map[string]int),
		attackerDamageAssignments: make(map[string]map[string]int),
		blockerDamageAssignments:  make(map[string]map[string]int),
	}
}

// EngineMessage represents a game log message
type EngineMessage struct {
	Text              string
	Color             string
	Timestamp         time.Time
	BookmarkID        int  // Snapshot ID taken before this message was added (0 = no snapshot)
	RollbackAvailable bool // Whether this state can be rolled back to
}

// EnginePrompt represents a prompt for player input
type EnginePrompt struct {
	PlayerID  string
	Text      string
	Options   []string
	Timestamp time.Time
}

// PendingXValueRequest represents an active X value selection request
// The game engine waits for the player to respond via SEND_INTEGER
type PendingXValueRequest struct {
	// PlayerID is the player who needs to select the X value
	PlayerID string
	// SourceID is the ID of the card/ability requiring X
	SourceID string
	// SourceName is the human-readable name of the source
	SourceName string
	// Message is a human-readable description shown to the player
	Message string
	// MinValue is the minimum allowed X value
	MinValue int
	// MaxValue is the maximum allowed X value
	MaxValue int
	// Timestamp when the request was created
	Timestamp time.Time
	// OnComplete is called when X value selection is complete
	OnComplete func(xValue int) error
	// OnCancel is called when the player cancels (if allowed)
	OnCancel func() error
}

// PendingTargetRequest represents an active target selection request
// The game engine waits for the player to respond via SEND_UUID
type PendingTargetRequest struct {
	// PlayerID is the player who needs to select targets
	PlayerID string
	// SourceID is the ID of the spell/ability requiring targets
	SourceID string
	// Requirement defines what kind of targets are needed
	Requirement targeting.TargetRequirement
	// ValidTargetIDs is a list of IDs that are valid targets
	ValidTargetIDs []string
	// SelectedTargetIDs contains targets selected so far (for multi-target)
	SelectedTargetIDs []string
	// Message is a human-readable description shown to the player
	Message string
	// Required indicates if the player MUST select targets (can't cancel)
	Required bool
	// Timestamp when the request was created
	Timestamp time.Time
	// OnComplete is called when target selection is complete
	OnComplete func(selectedTargets []string) error
	// OnCancel is called when the player cancels target selection (if allowed)
	OnCancel func() error
}

// PendingLibrarySearchRequest represents an active library search request
// The game engine shows the player their library and waits for card selection
type PendingLibrarySearchRequest struct {
	// PlayerID is the player searching their library
	PlayerID string
	// SearchingPlayerID is whose library is being searched (usually same as PlayerID)
	SearchingPlayerID string
	// Message describes what the player is searching for
	Message string
	// Destination is where the selected card goes: "hand", "battlefield", "top", "graveyard"
	Destination string
	// CardFilter optional filter description (e.g., "creature", "land", "any")
	CardFilter string
	// Shuffle indicates if library should be shuffled after search
	Shuffle bool
	// Required indicates if the player MUST select a card (can't cancel without selecting)
	Required bool
	// Timestamp when the request was created
	Timestamp time.Time
}

// internalCard represents a card in the game state
type internalCard struct {
	ID             string
	Name           string
	DisplayName    string
	ManaCost       string
	Type           string
	SubTypes       []string
	SuperTypes     []string
	Color          string
	Power          string
	Toughness      string
	Loyalty        string
	CardNumber     int
	ExpansionSet   string
	Rarity         string
	RulesText      string
	Tapped         bool
	Flipped        bool
	Transformed    bool
	FaceDown       bool
	Zone           int
	ControllerID   string
	OwnerID        string
	AttachedToCard []string
	Abilities      []EngineAbilityView
	Counters       *counters.Counters
	// Combat fields
	Attacking     bool     // Is this creature attacking
	Blocking      bool     // Is this creature blocking
	AttackingWhat string   // ID of what this creature is attacking (player/planeswalker/battle)
	BlockingWhat  []string // IDs of creatures this creature is blocking
	// Banding fields (Rule 702.22)
	BandedCards []string // IDs of creatures in the same attacking band (bidirectional)
	// Damage tracking
	Damage        int            // Damage marked on this creature
	DamageSources map[string]int // Damage by source ID
	// Status fields
	SummoningSickness bool // Does this creature have summoning sickness
	IsToken           bool // Is this a token (doesn't go to graveyard when destroyed)
	IsCommander       bool // Is this a commander card
	// Spell metadata
	Metadata map[string]string // Generic metadata for storing targets, choices, etc.
}

// LastKnownInfo stores the state of a permanent at the moment it left a zone
// Java: GameImpl.getLastKnownInformation(), lki/lkiExtended maps
// MTG Rules: 113.7a, 400.7 (Objects that change zones are tracked)
type LastKnownInfo struct {
	ID           string              // Permanent ID
	Name         string              // Card name
	ControllerID string              // Controller at time of leaving
	OwnerID      string              // Owner
	Types        string              // Type line
	SubTypes     []string            // Subtypes
	Power        string              // Power (for creatures)
	Toughness    string              // Toughness (for creatures)
	Counters     map[string]int      // All counters at time of leaving (name -> count)
	Tapped       bool                // Was it tapped
	Zone         int                 // Zone it was in
	ZoneCounter  int                 // Zone change counter (increments each zone change)
	Timestamp    time.Time           // When it left
	Abilities    []EngineAbilityView // Abilities it had
}

// copyCountersToMap converts a Counters object to a simple map for LKI storage
func copyCountersToMap(c *counters.Counters) map[string]int {
	if c == nil {
		return make(map[string]int)
	}
	result := make(map[string]int)
	for name, counter := range c.GetAll() {
		result[name] = counter.Count
	}
	return result
}

// createLKIFromCard creates a LastKnownInfo snapshot from an internalCard
func createLKIFromCard(card *internalCard, zoneCounter int) *LastKnownInfo {
	return &LastKnownInfo{
		ID:           card.ID,
		Name:         card.Name,
		ControllerID: card.ControllerID,
		OwnerID:      card.OwnerID,
		Types:        card.Type,
		SubTypes:     append([]string{}, card.SubTypes...),
		Power:        card.Power,
		Toughness:    card.Toughness,
		Counters:     copyCountersToMap(card.Counters),
		Tapped:       card.Tapped,
		Zone:         card.Zone,
		ZoneCounter:  zoneCounter,
		Timestamp:    time.Now(),
		Abilities:    append([]EngineAbilityView{}, card.Abilities...),
	}
}

// storeLKI stores Last Known Information for a permanent leaving the battlefield
// This must be called BEFORE the permanent is removed from the battlefield
func (gs *engineGameState) storeLKI(card *internalCard) {
	// Increment zone counter for this permanent
	gs.lkiZoneCounter[card.ID]++
	zoneCounter := gs.lkiZoneCounter[card.ID]

	// Create and store the LKI snapshot
	gs.lki[card.ID] = createLKIFromCard(card, zoneCounter)
}

// getLKI retrieves the Last Known Information for a permanent
// Returns nil if no LKI exists for the given ID
func (gs *engineGameState) getLKI(permanentID string) *LastKnownInfo {
	return gs.lki[permanentID]
}

// internalPlayer represents a player in the game state
type internalPlayer struct {
	PlayerID            string
	Name                string
	Life                int
	Poison              int
	Energy              int
	Library             []*internalCard
	Hand                []*internalCard
	Graveyard           []*internalCard
	ManaPool            *mana.ManaPool
	HasPriority         bool
	Passed              bool
	StateOrdinal        int
	Lost                bool
	Left                bool
	Wins                int
	Quit                bool           // Player quit the match
	TimerTimeout        bool           // Player lost due to timer timeout
	IdleTimeout         bool           // Player lost due to idle timeout
	Conceded            bool           // Player conceded
	StoredBookmark      int            // Bookmark ID for player undo (-1 = no undo available)
	MulliganCount       int            // Number of times player has mulliganed
	KeptHand            bool           // Whether player has kept their hand
	CommanderDamage     map[string]int // Tracks combat damage from each commander (commander ID -> damage)
	LandsPlayedThisTurn int            // Number of lands played this turn
	LandsPerTurn        int            // Maximum lands allowed per turn (default 1)
	PassUntil           PassUntilType  // Auto-pass mode for this player
	PassUntilTurn       int            // Turn number to pass until (for PassUntilMyNextTurn)
}

// triggeredAbilityQueueItem represents a triggered ability waiting to be put on the stack
type triggeredAbilityQueueItem struct {
	ID          string
	SourceID    string
	Controller  string
	Description string
	Resolve     func(*engineGameState) error
	UsesStack   bool // If false, executes immediately without going on stack
}

// combatTrigger represents a combat-related trigger condition
// Per Java TriggeredAbility pattern (checkEventType + checkTrigger)
type combatTrigger struct {
	SourceID      string                                                         // Card with the trigger
	TriggerType   string                                                         // Type of trigger (attacks, blocks, etc.)
	Condition     func(*engineGameState, rules.Event) bool                       // Check if trigger should fire
	CreateAbility func(*engineGameState, rules.Event) *triggeredAbilityQueueItem // Create the triggered ability
}

// gameAnalytics tracks metrics for a game
type gameAnalytics struct {
	maxStackDepth      int               // Maximum stack depth reached
	totalStackItems    int               // Total items put on stack
	actionsPerTurn     map[int]int       // Actions taken per turn number
	turnStartTimes     map[int]time.Time // Turn start times
	priorityPassCount  int               // Total priority passes
	spellsCast         int               // Total spells cast
	abilitiesActivated int               // Total abilities activated
	triggersProcessed  int               // Total triggered abilities processed
	gameStartTime      time.Time         // When game started
}

// engineGameState represents the internal state of a game
type engineGameState struct {
	gameID             string
	gameType           string
	gameTypeConfig     plugin.GameType       // Configured rules for this game type
	gameRules          plugin.GameRules      // Game rules (starting life, deck size, etc.)
	behaviors          []plugin.GameBehavior // Format-specific behaviors (commander, etc.)
	state              GameState
	players            map[string]*internalPlayer
	playerOrder        []string
	cards              map[string]*internalCard
	battlefield        []*internalCard
	exile              []*internalCard
	command            []*internalCard
	revealed           []EngineRevealedView
	lookedAt           []EngineLookedAtView
	combat             *combatState // Internal combat state
	turnManager        *rules.TurnManager
	stack              *rules.StackManager
	eventBus           *rules.EventBus
	watchers           *rules.WatcherRegistry
	legality           *rules.LegalityChecker
	targetValidator    *targeting.TargetValidator
	layerSystem        *effects.LayerSystem
	abilityRegistry    *AbilityRegistry             // Maps ability IDs to ability objects for retrieval
	triggeredQueue     []*triggeredAbilityQueueItem // Queue of triggered abilities waiting to be put on stack
	combatTriggers     []*combatTrigger             // Registered combat triggers (for cards with combat-related abilities)
	simultaneousEvents []rules.Event                // Queue of events that happened simultaneously
	concedingPlayers   []string                     // Queue of players requesting concession
	analytics          *gameAnalytics               // Game metrics and analytics
	messages           []EngineMessage
	prompts            []EnginePrompt
	startedAt          time.Time
	startingPlayerID   string // Player who won the coin flip and goes first (skips first draw)
	firstTurnDrawDone  bool   // Whether the first turn draw skip has been applied

	// Target selection system
	// When a spell/ability needs targets, we store the pending request here
	// and wait for the player to respond via SEND_UUID
	pendingTargetRequest *PendingTargetRequest

	// X value selection system
	// When a spell/ability needs X value input, we store the pending request here
	// and wait for the player to respond via SEND_INTEGER
	pendingXValueRequest *PendingXValueRequest

	// Library search system
	// When a player searches their library, we store the pending request here
	// and wait for the player to respond via SEND_UUID with the selected card
	pendingLibrarySearch *PendingLibrarySearchRequest

	// Last Known Information (LKI) system
	// Java: GameImpl.lki (Map<Zone, Map<UUID, MageObject>>)
	// Stores permanent state when it leaves the battlefield for triggered abilities
	lki            map[string]*LastKnownInfo // permanentID -> LKI snapshot
	lkiZoneCounter map[string]int            // permanentID -> zone change counter

	// Message-level rollback system
	// Tracks pending rollback requests awaiting opponent consent
	pendingRollbackRequest *PendingRollbackRequest
	// Maps message index to bookmark ID for rollback
	messageBookmarks map[int]int
	// Counter for unique message IDs
	nextMessageID int
	// Current action bookmark ID - set when processing an action, messages use this for rollback
	currentActionBookmark int

	mu sync.RWMutex
}

// PendingRollbackRequest represents an active rollback request awaiting opponent consent
type PendingRollbackRequest struct {
	RequestID         string    // UUID for tracking the request
	RequestingPlayer  string    // Player ID who requested the rollback
	TargetMessageID   int       // Message ID to rollback to
	TargetBookmarkID  int       // Bookmark ID associated with that message
	TargetMessageText string    // Text of the target message for display
	Timestamp         time.Time // When the request was made
}

// GameNotification represents a notification that can be sent to UI/websocket clients
type GameNotification struct {
	Type      string                 // Type of notification (e.g., "PRIORITY_CHANGE", "STACK_UPDATE", "COMBAT_UPDATE")
	GameID    string                 // Game ID
	PlayerID  string                 // Target player ID (empty for broadcast)
	Timestamp time.Time              // When the notification was created
	Data      map[string]interface{} // Notification-specific data
}

// NotificationHandler is a function that handles game notifications
type NotificationHandler func(notification GameNotification)

// gameStateSnapshot represents a complete snapshot of game state for rollback
type gameStateSnapshot struct {
	// Core game state
	GameID         string
	GameType       string
	State          GameState
	TurnNumber     int
	ActivePlayer   string
	PriorityPlayer string

	// Players - deep copy of all player data
	Players     map[string]*internalPlayer
	PlayerOrder []string

	// Cards - deep copy of all cards
	Cards       map[string]*internalCard
	Battlefield []*internalCard
	Exile       []*internalCard
	Command     []*internalCard

	// Stack state
	StackItems []rules.StackItem

	// Other state
	Messages  []EngineMessage
	Prompts   []EnginePrompt
	Timestamp time.Time
}

// PersistenceRepository interface for persisting game state to database
// This allows the engine to save/restore games without direct DB dependency
type PersistenceRepository interface {
	// SaveGameState persists the current game state
	SaveGameState(ctx context.Context, gameID, tableID, gameType string, players []string, gameState []byte, turnNumber int, state string) error
	// DeleteActiveGame removes a game from persistence (when finished)
	DeleteActiveGame(ctx context.Context, gameID string) error
}

// CardBuilderFunc is a function type for building cards from the registry
// This allows the cards package to inject card builders without import cycles
type CardBuilderFunc func(cardName string, ownerID uuid.UUID) (*Card, error)

// MageEngine is the main game engine implementation
type MageEngine struct {
	logger              *zap.Logger
	mu                  sync.RWMutex
	games               map[string]*engineGameState
	notificationHandler atomic.Value            // Stores NotificationHandler; uses atomic to avoid deadlock with gameState.mu
	cardRepo            CardRepositoryInterface // Optional card repository for looking up card metadata
	persistenceRepo     PersistenceRepository   // Optional persistence repository for crash recovery
	cardBuilder         CardBuilderFunc         // Optional card builder for Go-implemented cards

	// State bookmarking for rollback/undo
	// Maps gameID -> list of bookmarked states
	bookmarks map[string][]*gameStateSnapshot

	// Turn rollback system (separate from action bookmarks)
	// Maps gameID -> map[turnNumber -> snapshot]
	// Keeps last 4 turns for player-requested rollback
	turnSnapshots    map[string]map[int]*gameStateSnapshot
	rollbackTurnsMax int  // Maximum turns to keep for rollback (default 4)
	rollbackAllowed  bool // Whether turn rollback is enabled (default true)
	bookmarksMax     int  // Maximum bookmarks per game to prevent memory growth (default 100)

	// Replay recording system
	// Records step-by-step game state for replay and spectator synchronization
	replayRecorder *ReplayRecorder

	// Critical game systems (Tier 1 gaps integration)
	// Replacement effects (Rule 614) - ETB effects, doubling, death replacement
	replacementEffects map[string]*effects.ReplacementManager // gameID -> manager

	// TODO: Prevention effects (Rule 615) - Damage prevention, protection
	// Prevention effects exist but need a manager similar to ReplacementManager
	// preventionEffects map[string]*effects.PreventionManager // gameID -> manager
}

// CardRepositoryInterface interface for looking up card metadata
type CardRepositoryInterface interface {
	GetByName(ctx context.Context, name string) ([]*repository.Card, error)
}

// NewMageEngine creates a new MageEngine instance
func NewMageEngine(logger *zap.Logger) *MageEngine {
	return &MageEngine{
		logger:             logger,
		games:              make(map[string]*engineGameState),
		bookmarks:          make(map[string][]*gameStateSnapshot),
		turnSnapshots:      make(map[string]map[int]*gameStateSnapshot),
		rollbackTurnsMax:   4,                                            // Keep last 4 turns
		rollbackAllowed:    true,                                         // Enable turn rollback by default
		bookmarksMax:       100,                                          // Keep last 100 bookmarks per game
		replayRecorder:     NewReplayRecorder(logger, "replays"),         // Default replay directory
		replacementEffects: make(map[string]*effects.ReplacementManager), // Per-game replacement effects
	}
}

// SetNotificationHandler sets the handler for game notifications
// This allows external systems (UI, websockets) to receive real-time game updates
// Uses atomic.Value to avoid lock contention with emitNotification
func (e *MageEngine) SetNotificationHandler(handler NotificationHandler) {
	e.notificationHandler.Store(handler)
}

// SetCardRepository sets the card repository for looking up card metadata
func (e *MageEngine) SetCardRepository(repo CardRepositoryInterface) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cardRepo = repo
}

// SetPersistenceRepository sets the persistence repository for crash recovery
// When set, game state will be persisted to database at turn boundaries
func (e *MageEngine) SetPersistenceRepository(repo PersistenceRepository) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.persistenceRepo = repo
}

// SetCardBuilder sets the card builder function for creating Go-implemented cards
// This allows the cards package to inject its registry without import cycles
func (e *MageEngine) SetCardBuilder(builder CardBuilderFunc) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cardBuilder = builder
}

// emitNotification sends a notification to the registered handler
// This method is safe to call while holding gameState locks because:
//  1. It uses atomic.Value to read the handler without acquiring any locks
//  2. The handler is called in a separate goroutine, so it doesn't block
//  3. The goroutine can safely call back into the engine (e.g., GetGameView)
//     because it runs asynchronously after emitNotification returns
//
// IMPORTANT: This method must NOT acquire e.mu to avoid deadlock:
// - ProcessAction holds gameState.mu.Lock(), releases it, calls BookmarkState
// - BookmarkState acquires e.mu.Lock(), then tries gameState.mu.RLock()
// - Concurrent ProcessAction holds gameState.mu.Lock(), calls notify* -> emitNotification
// - If emitNotification tried e.mu.RLock(), it would deadlock (AB-BA pattern)
func (e *MageEngine) emitNotification(notification GameNotification) {
	// Use atomic.Load to read handler without any mutex
	// This prevents deadlock when called while holding gameState.mu
	handlerVal := e.notificationHandler.Load()
	if handlerVal == nil {
		return
	}

	handler, ok := handlerVal.(NotificationHandler)
	if !ok || handler == nil {
		return
	}

	// Call handler in a goroutine to avoid blocking game logic
	// The goroutine runs asynchronously, so it can safely acquire locks
	// (e.g., call GetGameView) after emitNotification returns
	go handler(notification)
}

// notifyPriorityChange notifies that priority has changed
func (e *MageEngine) notifyPriorityChange(gameID, playerID string, data map[string]interface{}) {
	e.emitNotification(GameNotification{
		Type:      "PRIORITY_CHANGE",
		GameID:    gameID,
		PlayerID:  playerID,
		Timestamp: time.Now(),
		Data:      data,
	})
}

// notifyStackUpdate notifies that the stack has changed
func (e *MageEngine) notifyStackUpdate(gameID string, data map[string]interface{}) {
	e.emitNotification(GameNotification{
		Type:      "STACK_UPDATE",
		GameID:    gameID,
		PlayerID:  "", // Broadcast to all players
		Timestamp: time.Now(),
		Data:      data,
	})
}

// notifyGameStateChange notifies that the game state has changed
func (e *MageEngine) notifyGameStateChange(gameID string, data map[string]interface{}) {
	e.emitNotification(GameNotification{
		Type:      "GAME_STATE_CHANGE",
		GameID:    gameID,
		PlayerID:  "", // Broadcast to all players
		Timestamp: time.Now(),
		Data:      data,
	})
}

// notifyPhaseChange notifies that the phase/step has changed
func (e *MageEngine) notifyPhaseChange(gameID string, data map[string]interface{}) {
	e.emitNotification(GameNotification{
		Type:      "PHASE_CHANGE",
		GameID:    gameID,
		PlayerID:  "", // Broadcast to all players
		Timestamp: time.Now(),
		Data:      data,
	})
}

// notifyChoicePrompt notifies that a player needs to make a choice
// This is used for combat declarations, mode choices, and other player selections
func (e *MageEngine) notifyChoicePrompt(gameID, playerID, message string, choices []string) {
	e.emitNotification(GameNotification{
		Type:      "GAME_CHOOSE_CHOICE",
		GameID:    gameID,
		PlayerID:  playerID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": message,
			"choices": choices,
		},
	})
}

// notifyGameError notifies a specific player about an error
func (e *MageEngine) notifyGameError(gameID, playerID string, errorMsg string) {
	e.emitNotification(GameNotification{
		Type:      "GAME_ERROR",
		GameID:    gameID,
		PlayerID:  playerID, // Send only to the player who caused the error
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"error": errorMsg,
		},
	})
}

// notifyPlayerAction notifies about a player action
func (e *MageEngine) notifyPlayerAction(gameID, playerID string, data map[string]interface{}) {
	e.emitNotification(GameNotification{
		Type:      "PLAYER_ACTION",
		GameID:    gameID,
		PlayerID:  "", // Broadcast to all players
		Timestamp: time.Now(),
		Data:      data,
	})
}

// notifyTrigger notifies about a triggered ability
func (e *MageEngine) notifyTrigger(gameID string, data map[string]interface{}) {
	e.emitNotification(GameNotification{
		Type:      "TRIGGER",
		GameID:    gameID,
		PlayerID:  "", // Broadcast to all players
		Timestamp: time.Now(),
		Data:      data,
	})
}

// notifyTargetRequest notifies a player that they need to select targets
// This sends a GAME_TARGET notification to the UI
func (e *MageEngine) notifyTargetRequest(gameID, playerID string, data map[string]interface{}) {
	e.emitNotification(GameNotification{
		Type:      "GAME_TARGET",
		GameID:    gameID,
		PlayerID:  playerID, // Send only to the player who needs to select
		Timestamp: time.Now(),
		Data:      data,
	})
}

// notifyXValueRequest notifies a player that they need to select an X value
// This sends a GAME_XMANA notification to the UI
func (e *MageEngine) notifyXValueRequest(gameID, playerID string, data map[string]interface{}) {
	e.emitNotification(GameNotification{
		Type:      "GAME_XMANA",
		GameID:    gameID,
		PlayerID:  playerID, // Send only to the player who needs to select
		Timestamp: time.Now(),
		Data:      data,
	})
}

// RequestXValueSelection initiates X value selection for a spell or ability.
// This sets up the pending X value request and notifies the player via GAME_XMANA.
// The game engine will wait for SEND_INTEGER response from the player.
func (e *MageEngine) RequestXValueSelection(
	gameState *engineGameState,
	gameID, playerID, sourceID, sourceName string,
	message string,
	minValue, maxValue int,
	onComplete func(xValue int) error,
	onCancel func() error,
) error {
	// Validate player exists
	_, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// Check if there's already a pending X value request
	if gameState.pendingXValueRequest != nil {
		return fmt.Errorf("X value selection already in progress for player %s", gameState.pendingXValueRequest.PlayerID)
	}

	// Create the pending request
	gameState.pendingXValueRequest = &PendingXValueRequest{
		PlayerID:   playerID,
		SourceID:   sourceID,
		SourceName: sourceName,
		Message:    message,
		MinValue:   minValue,
		MaxValue:   maxValue,
		Timestamp:  time.Now(),
		OnComplete: onComplete,
		OnCancel:   onCancel,
	}

	if e.logger != nil {
		e.logger.Info("requesting X value selection",
			zap.String("player", playerID),
			zap.String("source", sourceName),
			zap.Int("min", minValue),
			zap.Int("max", maxValue),
			zap.String("message", message))
	}

	// Notify the player via WebSocket
	e.notifyXValueRequest(gameID, playerID, map[string]interface{}{
		"message":   message,
		"available": maxValue,
		"min":       minValue,
		"max":       maxValue,
		"source_id": sourceID,
	})

	return nil
}

// RequestTargetSelection initiates target selection for a spell or ability.
// This sets up the pending target request and notifies the player via GAME_TARGET.
// The game engine will wait for SEND_UUID responses from the player.
//
// Parameters:
//   - gameID: The game ID
//   - playerID: The player who needs to select targets
//   - sourceID: The ID of the spell/ability requiring targets
//   - requirement: The target requirement specification
//   - message: Human-readable message to show the player
//   - required: Whether the player must select targets (can't cancel)
//   - onComplete: Callback when target selection is complete
//   - onCancel: Callback when player cancels (if allowed)
//
// Returns error if game not found or player doesn't have priority
func (e *MageEngine) RequestTargetSelection(
	gameID, playerID, sourceID string,
	requirement targeting.TargetRequirement,
	message string,
	required bool,
	onComplete func(selectedTargets []string) error,
	onCancel func() error,
) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	gameState, exists := e.games[gameID]
	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	// Verify the player exists
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// Check if there's already a pending target request
	if gameState.pendingTargetRequest != nil {
		return fmt.Errorf("there is already a pending target request")
	}

	// Find all valid targets based on the requirement
	validTargetIDs := e.findValidTargets(gameState, requirement)

	// If no valid targets and targeting is required, return error
	if len(validTargetIDs) == 0 && requirement.MinTargets > 0 {
		return fmt.Errorf("no valid targets available for %s", requirement.Description)
	}

	// Create the pending target request
	gameState.pendingTargetRequest = &PendingTargetRequest{
		PlayerID:          playerID,
		SourceID:          sourceID,
		Requirement:       requirement,
		ValidTargetIDs:    validTargetIDs,
		SelectedTargetIDs: make([]string, 0),
		Message:           message,
		Required:          required,
		Timestamp:         time.Now(),
		OnComplete:        onComplete,
		OnCancel:          onCancel,
	}

	// Build target card views for the UI
	targetViews := make([]map[string]interface{}, 0, len(validTargetIDs))
	for _, targetID := range validTargetIDs {
		if card, found := gameState.cards[targetID]; found {
			targetViews = append(targetViews, map[string]interface{}{
				"id":         card.ID,
				"name":       card.Name,
				"type":       card.Type,
				"zone":       zoneToString(card.Zone),
				"controller": card.ControllerID,
			})
		} else if _, isPlayer := gameState.players[targetID]; isPlayer {
			targetViews = append(targetViews, map[string]interface{}{
				"id":   targetID,
				"name": targetID,
				"type": "player",
			})
		}
	}

	// Log the target request
	if e.logger != nil {
		e.logger.Info("requesting target selection",
			zap.String("game_id", gameID),
			zap.String("player_id", playerID),
			zap.String("source_id", sourceID),
			zap.String("message", message),
			zap.Int("valid_targets", len(validTargetIDs)),
			zap.Int("min_targets", requirement.MinTargets),
			zap.Int("max_targets", requirement.MaxTargets),
			zap.Bool("required", required),
		)
	}

	// Notify the player that they need to select targets
	e.notifyTargetRequest(gameID, playerID, map[string]interface{}{
		"message":     message,
		"targets":     targetViews,
		"required":    required,
		"min_targets": requirement.MinTargets,
		"max_targets": requirement.MaxTargets,
		"source_id":   sourceID,
	})

	// Add a prompt for the player
	gameState.addPrompt(playerID, message, validTargetIDs)
	_ = player // Mark as used

	return nil
}

// findValidTargets finds all valid targets for a given requirement
func (e *MageEngine) findValidTargets(gameState *engineGameState, requirement targeting.TargetRequirement) []string {
	validTargets := make([]string, 0)

	switch requirement.Type {
	case targeting.TargetTypeCreature:
		// Find all creatures on the battlefield
		for _, card := range gameState.battlefield {
			if card == nil {
				continue
			}
			if strings.Contains(strings.ToLower(card.Type), "creature") {
				if err := gameState.targetValidator.ValidateTarget(card.ID, requirement); err == nil {
					validTargets = append(validTargets, card.ID)
				}
			}
		}

	case targeting.TargetTypePlayer:
		// All players who haven't lost/left are valid targets
		for playerID, player := range gameState.players {
			if !player.Lost && !player.Left {
				validTargets = append(validTargets, playerID)
			}
		}

	case targeting.TargetTypePermanent:
		// All permanents on the battlefield
		for _, card := range gameState.battlefield {
			if card == nil {
				continue
			}
			if err := gameState.targetValidator.ValidateTarget(card.ID, requirement); err == nil {
				validTargets = append(validTargets, card.ID)
			}
		}

	case targeting.TargetTypeSpell:
		// All spells on the stack
		stackItems := gameState.stack.List()
		for _, item := range stackItems {
			if item.Kind == "spell" {
				validTargets = append(validTargets, item.ID)
			}
		}

	case targeting.TargetTypeArtifact:
		for _, card := range gameState.battlefield {
			if card == nil {
				continue
			}
			if strings.Contains(strings.ToLower(card.Type), "artifact") {
				if err := gameState.targetValidator.ValidateTarget(card.ID, requirement); err == nil {
					validTargets = append(validTargets, card.ID)
				}
			}
		}

	case targeting.TargetTypeEnchantment:
		for _, card := range gameState.battlefield {
			if card == nil {
				continue
			}
			if strings.Contains(strings.ToLower(card.Type), "enchantment") {
				if err := gameState.targetValidator.ValidateTarget(card.ID, requirement); err == nil {
					validTargets = append(validTargets, card.ID)
				}
			}
		}

	case targeting.TargetTypeLand:
		for _, card := range gameState.battlefield {
			if card == nil {
				continue
			}
			if strings.Contains(strings.ToLower(card.Type), "land") {
				if err := gameState.targetValidator.ValidateTarget(card.ID, requirement); err == nil {
					validTargets = append(validTargets, card.ID)
				}
			}
		}

	case targeting.TargetTypePlaneswalker:
		for _, card := range gameState.battlefield {
			if card == nil {
				continue
			}
			if strings.Contains(strings.ToLower(card.Type), "planeswalker") {
				if err := gameState.targetValidator.ValidateTarget(card.ID, requirement); err == nil {
					validTargets = append(validTargets, card.ID)
				}
			}
		}
	}

	return validTargets
}

// CancelTargetSelection cancels the current pending target request
// Returns error if no pending request or if request is required (can't cancel)
func (e *MageEngine) CancelTargetSelection(gameID, playerID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	gameState, exists := e.games[gameID]
	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	if gameState.pendingTargetRequest == nil {
		return fmt.Errorf("no pending target request")
	}

	if gameState.pendingTargetRequest.PlayerID != playerID {
		return fmt.Errorf("target request is for player %s, not %s", gameState.pendingTargetRequest.PlayerID, playerID)
	}

	if gameState.pendingTargetRequest.Required {
		return fmt.Errorf("target selection is required and cannot be cancelled")
	}

	// Call the cancel callback
	if gameState.pendingTargetRequest.OnCancel != nil {
		if err := gameState.pendingTargetRequest.OnCancel(); err != nil {
			return fmt.Errorf("cancel callback failed: %v", err)
		}
	}

	// Clear the pending request
	gameState.pendingTargetRequest = nil

	if e.logger != nil {
		e.logger.Info("target selection cancelled",
			zap.String("game_id", gameID),
			zap.String("player_id", playerID),
		)
	}

	return nil
}

// StartGame initializes a new game state
// Deprecated: Use StartGameWithDecks instead for proper deck loading
func (e *MageEngine) StartGame(gameID string, players []string, gameType string) error {
	return e.StartGameWithDecks(gameID, players, gameType, nil)
}

// StartGameWithDecks initializes a new game state with player-submitted decks
func (e *MageEngine) StartGameWithDecks(gameID string, players []string, gameType string, decks map[string]DeckList) error {
	if e.logger != nil {
		e.logger.Info("[ENGINE] StartGameWithDecks called",
			zap.String("game_id", gameID),
			zap.Strings("players", players),
			zap.String("game_type", gameType),
			zap.Bool("decks_provided", decks != nil),
			zap.Int("decks_count", len(decks)),
		)
		// Log each deck received
		for playerID, deck := range decks {
			e.logger.Info("[ENGINE] Deck received for player",
				zap.String("player", playerID),
				zap.Int("main_deck_size", len(deck.MainDeck)),
				zap.Int("sideboard_size", len(deck.Sideboard)),
				zap.Int("commander_count", len(deck.Commanders)),
			)
			if len(deck.MainDeck) > 0 {
				first5 := deck.MainDeck
				if len(first5) > 5 {
					first5 = first5[:5]
				}
				e.logger.Info("[ENGINE] First 5 cards in deck",
					zap.String("player", playerID),
					zap.Strings("cards", first5),
				)
			}
		}
	}

	if gameID == "" {
		return fmt.Errorf("gameID is required")
	}
	if len(players) < 2 {
		return fmt.Errorf("at least 2 players required")
	}

	e.mu.Lock()
	// Note: We manually unlock before calling notifications to avoid deadlock
	// Do not use defer here

	if _, exists := e.games[gameID]; exists {
		e.mu.Unlock()
		return fmt.Errorf("game %s already exists", gameID)
	}

	// Look up game type configuration from registry
	var gameTypeConfig plugin.GameType
	gameRules := plugin.DefaultGameRules()
	var behaviors []plugin.GameBehavior

	if gt, err := plugin.GetGameType(gameType); err == nil {
		gameTypeConfig = gt
		gameRules = plugin.GetRulesForGameType(gt)
		behaviors = plugin.GetBehaviorsForGameType(gt)
		if e.logger != nil {
			e.logger.Info("using game type configuration",
				zap.String("game_type", gameType),
				zap.Int("starting_life", gameRules.StartingLife),
				zap.Int("minimum_deck_size", gameRules.MinimumDeckSize),
				zap.Int("starting_hand_size", gameRules.StartingHandSize),
				zap.Int("behavior_count", len(behaviors)),
			)
		}
	} else {
		if e.logger != nil {
			e.logger.Warn("game type not found in registry, using defaults",
				zap.String("game_type", gameType),
				zap.Error(err),
			)
		}
	}

	// Coin flip to determine starting player (MTG Rule 103.2)
	// The winner of the flip chooses to play first or draw first
	coinFlipResult, err := rand.Int(rand.Reader, big.NewInt(int64(len(players))))
	if err != nil {
		e.mu.Unlock()
		return fmt.Errorf("failed to perform coin flip: %w", err)
	}
	startingPlayerIndex := int(coinFlipResult.Int64())
	startingPlayerID := players[startingPlayerIndex]

	// Create game state - start in MULLIGAN state
	gameState := &engineGameState{
		gameID:           gameID,
		gameType:         gameType,
		gameTypeConfig:   gameTypeConfig,
		gameRules:        gameRules,
		behaviors:        behaviors,
		state:            GameStateMulligan, // Start in mulligan phase
		players:          make(map[string]*internalPlayer),
		playerOrder:      make([]string, len(players)),
		cards:            make(map[string]*internalCard),
		battlefield:      make([]*internalCard, 0),
		exile:            make([]*internalCard, 0),
		command:          make([]*internalCard, 0),
		revealed:         make([]EngineRevealedView, 0),
		lookedAt:         make([]EngineLookedAtView, 0),
		combat:           newCombatState(),
		startingPlayerID: startingPlayerID,
		lki:              make(map[string]*LastKnownInfo),
		lkiZoneCounter:   make(map[string]int),
		analytics: &gameAnalytics{
			actionsPerTurn: make(map[int]int),
			turnStartTimes: make(map[int]time.Time),
			gameStartTime:  time.Now(),
		},
		messages:         make([]EngineMessage, 0),
		prompts:          make([]EnginePrompt, 0),
		startedAt:        time.Now(),
		messageBookmarks: make(map[int]int),
		nextMessageID:    1,
	}

	// Initialize supporting systems
	gameState.stack = rules.NewStackManager()
	gameState.eventBus = rules.NewEventBus()
	gameState.watchers = rules.NewWatcherRegistry()
	gameState.layerSystem = effects.NewLayerSystem()
	gameState.abilityRegistry = NewAbilityRegistry()

	// Initialize per-game effect managers (Rule 614)
	e.replacementEffects[gameID] = effects.NewReplacementManager(e.logger)

	// Create players
	for i, playerID := range players {
		gameState.playerOrder[i] = playerID
		gameState.players[playerID] = &internalPlayer{
			PlayerID:            playerID,
			Name:                playerID,
			Life:                gameRules.StartingLife, // Use game type starting life
			Poison:              0,
			Energy:              0,
			Library:             make([]*internalCard, 0),
			Hand:                make([]*internalCard, 0),
			Graveyard:           make([]*internalCard, 0),
			ManaPool:            mana.NewManaPool(),
			HasPriority:         false,
			Passed:              false,
			StateOrdinal:        0,
			Lost:                false,
			Left:                false,
			Wins:                0,
			StoredBookmark:      -1,                   // No undo available initially
			MulliganCount:       0,                    // No mulligans yet
			KeptHand:            false,                // Haven't kept hand yet
			CommanderDamage:     make(map[string]int), // Track commander damage from each commander
			LandsPlayedThisTurn: 0,                    // No lands played yet
			LandsPerTurn:        1,                    // Default: 1 land per turn
		}

		// Load player's deck if provided, otherwise use test deck
		var deckCardNames []string
		var commanderNames []string

		if e.logger != nil {
			e.logger.Info("[ENGINE] Loading deck for player",
				zap.String("player", playerID),
				zap.Bool("decks_map_exists", decks != nil),
			)
		}

		if decks != nil {
			playerDeck, ok := decks[playerID]
			if e.logger != nil {
				e.logger.Info("[ENGINE] Checking deck map for player",
					zap.String("player", playerID),
					zap.Bool("found_in_map", ok),
					zap.Int("main_deck_size_if_found", len(playerDeck.MainDeck)),
				)
			}
			if ok && len(playerDeck.MainDeck) > 0 {
				deckCardNames = playerDeck.MainDeck
				commanderNames = playerDeck.Commanders
				if e.logger != nil {
					e.logger.Info("[ENGINE] USING PLAYER'S SUBMITTED DECK",
						zap.String("player", playerID),
						zap.Int("main_deck_size", len(deckCardNames)),
						zap.Int("commander_count", len(commanderNames)),
					)
				}
			}
		}

		// Fall back to test deck if no deck provided
		if len(deckCardNames) == 0 {
			if e.logger != nil {
				e.logger.Warn("[ENGINE] NO DECK FOUND - FALLING BACK TO TEST DECK!",
					zap.String("player", playerID),
				)
			}
			testDeckNames := []string{"Lightning Bolt", "Lightning Bolt", "Counterspell", "Counterspell", "Shock", "Shock",
				"Lightning Bolt", "Counterspell", "Shock", "Lightning Bolt"}
			for j := 0; j < 60; j++ {
				deckCardNames = append(deckCardNames, testDeckNames[j%len(testDeckNames)])
			}
		}

		// Create cards in library from deck
		for j, cardName := range deckCardNames {
			card := e.createStarterCard(fmt.Sprintf("%s-card-%d", playerID, j), playerID, cardName)
			gameState.cards[card.ID] = card
			gameState.players[playerID].Library = append(gameState.players[playerID].Library, card)
			card.Zone = zoneLibrary
		}

		// Create commander cards and move to command zone if applicable
		// Check if this game type has commander behavior
		hasCommanderBehavior := gameState.hasCommanderBehavior()
		if hasCommanderBehavior && len(commanderNames) > 0 {
			for j, commanderName := range commanderNames {
				card := e.createStarterCard(fmt.Sprintf("%s-commander-%d", playerID, j), playerID, commanderName)
				card.IsCommander = true
				gameState.cards[card.ID] = card
				gameState.command = append(gameState.command, card)
				card.Zone = zoneCommand
				if e.logger != nil {
					e.logger.Info("commander moved to command zone",
						zap.String("player", playerID),
						zap.String("commander", commanderName),
						zap.String("card_id", card.ID),
					)
				}
			}
		}

		// Shuffle library using Fisher-Yates with crypto/rand
		e.shuffleLibrary(gameState.players[playerID])
	}

	// Initialize turn manager with starting player (determined by coin flip)
	gameState.turnManager = rules.NewTurnManager(startingPlayerID)
	gameState.players[startingPlayerID].HasPriority = true

	// Initialize legality checker and target validator
	gameState.legality = rules.NewLegalityChecker(gameState)
	gameState.targetValidator = targeting.NewTargetValidator(gameState)

	// Wire up event bus to watchers
	gameState.eventBus.Subscribe(func(event rules.Event) {
		gameState.watchers.NotifyWatchers(event)
	})

	// Add coin flip message
	gameState.addMessage(fmt.Sprintf("%s wins the coin flip and will play first", startingPlayerID), "action")

	// Draw initial hands (7 cards each)
	for _, playerID := range players {
		player := gameState.players[playerID]
		for j := 0; j < 7 && len(player.Library) > 0; j++ {
			// Draw from top of library
			card := player.Library[len(player.Library)-1]
			player.Library = player.Library[:len(player.Library)-1]
			card.Zone = zoneHand
			player.Hand = append(player.Hand, card)
		}
		gameState.addMessage(fmt.Sprintf("%s draws opening hand of %d cards", playerID, len(player.Hand)), "action")
	}

	// Prompt each player for mulligan decision
	for _, playerID := range players {
		gameState.addPrompt(playerID, "Keep hand or mulligan?", []string{"KEEP", "MULLIGAN"})
	}

	e.games[gameID] = gameState

	// Release lock before sending notifications to avoid deadlock
	// Notifications may trigger callbacks that need to acquire locks
	e.mu.Unlock()

	// Record initial replay state
	// Per Java GameImpl.init() line 1246: saveState(false) after initialization
	gameState.mu.RLock()
	e.recordReplayState(gameState)
	gameState.mu.RUnlock()

	// Save initial turn snapshot (turn 1)
	// Per Java: save state at start of each turn
	if err := e.SaveTurnSnapshot(gameID, 1); err != nil {
		if e.logger != nil {
			e.logger.Warn("failed to save initial turn snapshot",
				zap.String("game_id", gameID),
				zap.Error(err),
			)
		}
	}

	if e.logger != nil {
		e.logger.Info("mage engine started game",
			zap.String("game_id", gameID),
			zap.Strings("players", players),
			zap.String("game_type", gameType),
			zap.String("starting_player", startingPlayerID),
		)
	}

	// Notify game start - players need to make mulligan decisions
	// Game stays in GameStateMulligan until all players have kept their hands
	e.notifyGameStateChange(gameID, map[string]interface{}{
		"state":           "mulligan",
		"game_type":       gameType,
		"players":         players,
		"starting_player": startingPlayerID,
	})

	if e.logger != nil {
		e.logger.Info("game in mulligan phase, waiting for player decisions",
			zap.String("game_id", gameID),
			zap.String("starting_player", startingPlayerID),
		)
	}

	return nil
}

// shuffleLibrary performs a Fisher-Yates shuffle using crypto/rand
func (e *MageEngine) shuffleLibrary(player *internalPlayer) {
	n := len(player.Library)
	for i := n - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			// Fallback to deterministic swap on error (shouldn't happen)
			continue
		}
		j := int(jBig.Int64())
		player.Library[i], player.Library[j] = player.Library[j], player.Library[i]
	}
}

// createStarterCard creates a simple starter card for testing
// Priority: 1) Go card registry (properly implemented cards)
//  2. JSON repository lookup (fallback for unimplemented cards)
//  3. Default values (last resort)
func (e *MageEngine) createStarterCard(id, ownerID, cardName string) *internalCard {
	if cardName == "" {
		cardName = "Lightning Bolt"
	}

	// 1. Check if card has a Go implementation via the injected card builder
	if e.cardBuilder != nil {
		ownerUUID, err := uuid.Parse(ownerID)
		if err == nil {
			gameCard, err := e.cardBuilder(cardName, ownerUUID)
			if err == nil && gameCard != nil {
				// Convert to internal card format
				internal := gameCard.ToInternal()
				internal.ID = id // Use the provided ID
				internal.Metadata = make(map[string]string)

				if e.logger != nil {
					e.logger.Info("created card from Go registry",
						zap.String("card_name", cardName),
						zap.String("id", id),
						zap.Int("ability_count", len(gameCard.Abilities)))
				}

				return internal
			}
		}
	}

	// 2. Default values (fallback if repository lookup fails)
	manaCost := "{R}"
	cardType := "Instant"
	subTypes := []string{}
	superTypes := []string{}
	color := "Red"
	power := ""
	toughness := ""
	loyalty := ""
	cardNumber := 1
	expansionSet := "M21"
	rarity := "Common"
	rulesText := fmt.Sprintf("%s deals damage.", cardName)

	// 3. Try to look up card metadata from repository if available
	if e.cardRepo != nil {
		ctx := context.Background()
		cards, err := e.cardRepo.GetByName(ctx, cardName)
		if err == nil && len(cards) > 0 {
			// Use first printing
			cardData := cards[0]
			manaCost = cardData.ManaCost
			cardType = cardData.CardType
			power = cardData.Power
			toughness = cardData.Toughness
			cardNumber, _ = strconv.Atoi(cardData.CardNumber)
			expansionSet = cardData.SetCode
			rarity = cardData.Rarity
			if cardData.RulesText != "" {
				rulesText = cardData.RulesText
			}

			// Parse card type string to extract types, subtypes, and supertypes
			// Format is typically: "Supertype Type — Subtype" or "Type — Subtype"
			typeParts := strings.Split(cardType, " — ")
			if len(typeParts) > 1 {
				// Has subtypes
				mainType := strings.TrimSpace(typeParts[0])
				subtypeStr := strings.TrimSpace(typeParts[1])
				subTypes = strings.Fields(subtypeStr)

				// Check for supertypes (Legendary, Basic, etc.)
				mainTypeParts := strings.Fields(mainType)
				if len(mainTypeParts) > 1 {
					// First parts are supertypes, last is the main type
					superTypes = mainTypeParts[:len(mainTypeParts)-1]
					cardType = mainTypeParts[len(mainTypeParts)-1]
				} else {
					cardType = mainType
				}
			} else {
				// No subtypes, but might have supertypes
				typeParts := strings.Fields(cardType)
				if len(typeParts) > 1 {
					superTypes = typeParts[:len(typeParts)-1]
					cardType = typeParts[len(typeParts)-1]
				}
			}

			// Parse color from mana cost (simple heuristic)
			if strings.Contains(manaCost, "{W}") {
				color = "White"
			} else if strings.Contains(manaCost, "{U}") {
				color = "Blue"
			} else if strings.Contains(manaCost, "{B}") {
				color = "Black"
			} else if strings.Contains(manaCost, "{R}") {
				color = "Red"
			} else if strings.Contains(manaCost, "{G}") {
				color = "Green"
			} else if manaCost == "" {
				color = "Colorless"
			} else {
				color = "Multicolor"
			}
		}
	}

	return &internalCard{
		ID:           id,
		Name:         cardName,
		DisplayName:  cardName,
		ManaCost:     manaCost,
		Type:         cardType,
		SubTypes:     subTypes,
		SuperTypes:   superTypes,
		Color:        color,
		Power:        power,
		Toughness:    toughness,
		Loyalty:      loyalty,
		CardNumber:   cardNumber,
		ExpansionSet: expansionSet,
		Rarity:       rarity,
		RulesText:    rulesText,
		Tapped:       false,
		Flipped:      false,
		Transformed:  false,
		FaceDown:     false,
		Zone:         zoneLibrary,
		Metadata:     make(map[string]string),
		ControllerID: ownerID,
		OwnerID:      ownerID,
		Counters:     counters.NewCounters(),
	}
}

// ProcessAction processes a player action with automatic error recovery
// Per Java GameImpl.playPriority(): creates bookmark before action, restores on error
func (e *MageEngine) ProcessAction(gameID string, action PlayerAction) (err error) {
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] ProcessAction entering",
			zap.String("game_id", gameID),
			zap.String("player_id", action.PlayerID),
			zap.String("action_type", action.ActionType))
	}

	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] ProcessAction acquiring gameState.mu.Lock",
			zap.String("game_id", gameID))
	}
	gameState.mu.Lock()
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] ProcessAction acquired gameState.mu.Lock",
			zap.String("game_id", gameID))
	}
	lockHeld := true // Track lock state for panic safety
	defer func() {
		if lockHeld {
			if e.logger != nil {
				e.logger.Info("[LOCK-DEBUG] ProcessAction releasing gameState.mu.Lock (defer)",
					zap.String("game_id", gameID))
			}
			gameState.mu.Unlock()
		}
	}()

	if gameState.state == GameStateFinished {
		return fmt.Errorf("game %s has ended", gameID)
	}

	// Create bookmark before processing action for error recovery
	// Per Java GameImpl.playPriority() line 1728: rollbackBookmarkOnPriorityStart = bookmarkState()
	var bookmarkID int
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] ProcessAction releasing gameState.mu for BookmarkState",
			zap.String("game_id", gameID))
	}
	gameState.mu.Unlock() // Temporarily unlock to call BookmarkState
	lockHeld = false
	bookmarkID, bookmarkErr := e.BookmarkState(gameID)
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] ProcessAction re-acquiring gameState.mu after BookmarkState",
			zap.String("game_id", gameID))
	}
	gameState.mu.Lock() // Re-acquire lock
	lockHeld = true
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] ProcessAction re-acquired gameState.mu after BookmarkState",
			zap.String("game_id", gameID))
	}

	if bookmarkErr != nil {
		if e.logger != nil {
			e.logger.Warn("failed to create bookmark before action",
				zap.String("game_id", gameID),
				zap.Error(bookmarkErr),
			)
		}
		// Continue without bookmark - error recovery won't be available
		bookmarkID = 0
	} else {
		// Set player's stored bookmark for undo
		// Per Java PlayerImpl.setStoredBookmark(): enables undo button
		if player, exists := gameState.players[action.PlayerID]; exists {
			player.StoredBookmark = bookmarkID
		}
	}

	// Defer error recovery: if action fails and we have a bookmark, restore state
	defer func() {
		if err != nil && bookmarkID > 0 {
			// Restore to bookmarked state on error
			// Per Java GameImpl.playPriority() line 1800: restoreState(rollbackBookmarkOnPriorityStart, "Game error: " + e)
			gameState.mu.Unlock() // Temporarily unlock to call RestoreState
			restoreErr := e.RestoreState(gameID, bookmarkID, fmt.Sprintf("Error recovery: %v", err))
			gameState.mu.Lock() // Re-acquire lock

			if restoreErr != nil {
				if e.logger != nil {
					e.logger.Error("failed to restore state after error",
						zap.String("game_id", gameID),
						zap.Int("bookmark_id", bookmarkID),
						zap.Error(err),
						zap.Error(restoreErr),
					)
				}
			} else {
				if e.logger != nil {
					e.logger.Info("auto-restored game state after error",
						zap.String("game_id", gameID),
						zap.Int("bookmark_id", bookmarkID),
						zap.Error(err),
					)
				}
				// Update error message to indicate restoration
				err = fmt.Errorf("action failed and state restored: %w", err)
			}

			// Notify the player about the error
			// Use the original error message (before wrapping) for cleaner user display
			gameState.mu.Unlock() // Temporarily unlock to emit notification
			e.notifyGameError(gameID, action.PlayerID, err.Error())
			gameState.mu.Lock() // Re-acquire lock
		} else if err != nil {
			// Error occurred but no bookmark to restore - still notify the player
			gameState.mu.Unlock() // Temporarily unlock to emit notification
			e.notifyGameError(gameID, action.PlayerID, err.Error())
			gameState.mu.Lock() // Re-acquire lock
		} else if bookmarkID > 0 {
			// Action succeeded, check if any player is using this bookmark
			// If so, don't remove it (player undo takes precedence)
			// Per Java: bookmark is kept if player stored it for undo
			bookmarkInUse := false
			for _, player := range gameState.players {
				if player.StoredBookmark == bookmarkID {
					bookmarkInUse = true
					break
				}
			}

			if !bookmarkInUse {
				// Remove the bookmark since no player is using it
				gameState.mu.Unlock() // Temporarily unlock to call RemoveBookmark
				e.RemoveBookmark(gameID, bookmarkID)
				gameState.mu.Lock() // Re-acquire lock
			}
		}

		// Clear current action bookmark after action completes
		gameState.currentActionBookmark = 0
	}()

	// Route action by type
	switch action.ActionType {
	case "PLAYER_ACTION":
		return e.handlePlayerAction(gameState, action)
	case "SEND_STRING":
		return e.handleStringAction(gameState, action)
	case "SEND_INTEGER":
		return e.handleIntegerAction(gameState, action)
	case "SEND_UUID":
		return e.handleUUIDAction(gameState, action)
	case "SPECIAL_ACTION":
		return e.handleSpecialAction(gameState, action)
	case "ACTIVATE_ABILITY":
		return e.handleActivateAbilityAction(gameState, action)
	default:
		return fmt.Errorf("unknown action type: %s", action.ActionType)
	}
}

// handlePlayerAction handles PLAYER_ACTION type actions
func (e *MageEngine) handlePlayerAction(gameState *engineGameState, action PlayerAction) error {
	dataStr, ok := action.Data.(string)
	if !ok {
		return fmt.Errorf("PLAYER_ACTION data must be string")
	}

	dataStr = strings.ToUpper(strings.TrimSpace(dataStr))

	switch dataStr {
	case "PASS":
		return e.handlePass(gameState, action.PlayerID)
	case "KEEP":
		return e.handleKeepHand(gameState, action.PlayerID)
	case "MULLIGAN":
		return e.handleMulligan(gameState, action.PlayerID)
	case "PASS_UNTIL_END_OF_TURN":
		return e.handlePassUntil(gameState, action.PlayerID, PassUntilEndOfTurn)
	case "PASS_UNTIL_NEXT_TURN":
		return e.handlePassUntil(gameState, action.PlayerID, PassUntilNextTurn)
	case "PASS_UNTIL_STACK_RESOLVED":
		return e.handlePassUntil(gameState, action.PlayerID, PassUntilStackResolved)
	case "PASS_UNTIL_MY_NEXT_TURN":
		return e.handlePassUntil(gameState, action.PlayerID, PassUntilMyNextTurn)
	default:
		return fmt.Errorf("unknown player action: %s", dataStr)
	}
}

// handleKeepHand handles a player choosing to keep their hand during mulligan
func (e *MageEngine) handleKeepHand(gameState *engineGameState, playerID string) error {
	if e.logger != nil {
		e.logger.Info("[MULLIGAN] handleKeepHand called",
			zap.String("game_id", gameState.gameID),
			zap.String("player_id", playerID),
			zap.String("current_state", gameState.state.String()),
		)
	}

	player, exists := gameState.players[playerID]
	if !exists {
		if e.logger != nil {
			e.logger.Warn("[MULLIGAN] handleKeepHand failed - player not found",
				zap.String("game_id", gameState.gameID),
				zap.String("player_id", playerID),
			)
		}
		return fmt.Errorf("player %s not found", playerID)
	}

	if e.logger != nil {
		e.logger.Info("[MULLIGAN] Player state before keep",
			zap.String("game_id", gameState.gameID),
			zap.String("player_id", playerID),
			zap.Bool("already_kept_hand", player.KeptHand),
			zap.Int("hand_size", len(player.Hand)),
			zap.Int("mulligan_count", player.MulliganCount),
		)
	}

	if gameState.state != GameStateMulligan {
		if e.logger != nil {
			e.logger.Warn("[MULLIGAN] handleKeepHand failed - not in mulligan phase",
				zap.String("game_id", gameState.gameID),
				zap.String("player_id", playerID),
				zap.String("actual_state", gameState.state.String()),
			)
		}
		return fmt.Errorf("not in mulligan phase")
	}

	if player.KeptHand {
		if e.logger != nil {
			e.logger.Warn("[MULLIGAN] handleKeepHand failed - player already kept hand",
				zap.String("game_id", gameState.gameID),
				zap.String("player_id", playerID),
			)
		}
		return fmt.Errorf("player %s has already kept hand", playerID)
	}

	player.KeptHand = true
	gameState.addMessage(fmt.Sprintf("%s keeps their hand (%d cards)", playerID, len(player.Hand)), "action")

	if e.logger != nil {
		e.logger.Info("[MULLIGAN] Player kept hand successfully",
			zap.String("game_id", gameState.gameID),
			zap.String("player_id", playerID),
			zap.Int("hand_size", len(player.Hand)),
		)
	}

	// Check if all players have kept their hands
	allKept := true
	playersNotKept := []string{}
	for pid, p := range gameState.players {
		if !p.KeptHand {
			allKept = false
			playersNotKept = append(playersNotKept, pid)
		}
	}

	if e.logger != nil {
		e.logger.Info("[MULLIGAN] Checking if all players kept",
			zap.String("game_id", gameState.gameID),
			zap.Bool("all_kept", allKept),
			zap.Strings("players_still_deciding", playersNotKept),
		)
	}

	if allKept {
		// All players have kept - transition to main game
		e.completeMulliganPhase(gameState)
	}

	// Notify state change
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{
		"player":    playerID,
		"action":    "keep",
		"hand_size": len(player.Hand),
	})

	return nil
}

// handleMulligan handles a player choosing to mulligan
func (e *MageEngine) handleMulligan(gameState *engineGameState, playerID string) error {
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	if gameState.state != GameStateMulligan {
		return fmt.Errorf("not in mulligan phase")
	}

	if player.KeptHand {
		return fmt.Errorf("player %s has already kept hand", playerID)
	}

	// Increment mulligan count
	player.MulliganCount++

	// Calculate new hand size (London Mulligan: draw 7, then put X on bottom)
	// For simplicity, we'll draw 7 and then they'll discard to (7 - mulliganCount)
	newHandSize := 7 - player.MulliganCount
	if newHandSize < 0 {
		newHandSize = 0
	}

	// Return current hand to library
	for _, card := range player.Hand {
		if card == nil {
			continue
		}
		card.Zone = zoneLibrary
		player.Library = append(player.Library, card)
	}
	player.Hand = make([]*internalCard, 0)

	// Shuffle library
	e.shuffleLibrary(player)

	// Draw new hand of 7 cards
	for i := 0; i < 7 && len(player.Library) > 0; i++ {
		card := player.Library[len(player.Library)-1]
		player.Library = player.Library[:len(player.Library)-1]
		card.Zone = zoneHand
		player.Hand = append(player.Hand, card)
	}

	gameState.addMessage(fmt.Sprintf("%s takes mulligan #%d (drew 7 cards, must put %d on bottom)",
		playerID, player.MulliganCount, player.MulliganCount), "action")

	// If mulligan count reached hand size, auto-keep (they'd have 0 cards)
	if newHandSize == 0 {
		player.KeptHand = true
		gameState.addMessage(fmt.Sprintf("%s forced to keep (no cards left to draw)", playerID), "action")
	} else {
		// Prompt for cards to put on bottom (simplified: just keep immediately for now)
		// TODO: Implement London Mulligan bottom selection
		// For now, auto-put the last N cards on bottom
		cardsToBottom := player.MulliganCount
		if cardsToBottom > 0 && len(player.Hand) >= cardsToBottom {
			for i := 0; i < cardsToBottom; i++ {
				card := player.Hand[len(player.Hand)-1]
				player.Hand = player.Hand[:len(player.Hand)-1]
				card.Zone = zoneLibrary
				// Put on bottom of library (beginning of slice)
				player.Library = append([]*internalCard{card}, player.Library...)
			}
			gameState.addMessage(fmt.Sprintf("%s puts %d cards on bottom of library (hand now %d cards)",
				playerID, cardsToBottom, len(player.Hand)), "action")
		}

		// Add prompt for next mulligan decision
		gameState.addPrompt(playerID, fmt.Sprintf("Keep hand (%d cards) or mulligan again?", len(player.Hand)), []string{"KEEP", "MULLIGAN"})
	}

	// Check if all players have kept their hands
	allKept := true
	for _, p := range gameState.players {
		if !p.KeptHand {
			allKept = false
			break
		}
	}

	if allKept {
		// All players have kept - transition to main game
		e.completeMulliganPhase(gameState)
	}

	// Notify state change
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{
		"player":         playerID,
		"action":         "mulligan",
		"mulligan_count": player.MulliganCount,
		"hand_size":      len(player.Hand),
		"kept":           player.KeptHand,
	})

	return nil
}

// completeMulliganPhase transitions from mulligan to main game
func (e *MageEngine) completeMulliganPhase(gameState *engineGameState) {
	if e.logger != nil {
		e.logger.Info("[MULLIGAN] completeMulliganPhase called - transitioning to IN_PROGRESS",
			zap.String("game_id", gameState.gameID),
			zap.String("old_state", gameState.state.String()),
		)
	}

	gameState.state = GameStateInProgress

	// Give priority to the starting player
	startingPlayerID := gameState.playerOrder[0]
	for _, playerID := range gameState.playerOrder {
		if player, exists := gameState.players[playerID]; exists && player.HasPriority {
			startingPlayerID = playerID
			break
		}
	}

	// Set active player based on starting player determination
	gameState.turnManager.SetPriority(startingPlayerID)

	gameState.addMessage("Mulligan phase complete, game begins", "system")
	gameState.addMessage(fmt.Sprintf("Turn 1 - %s's turn (Beginning Phase)", startingPlayerID), "phase")

	// Notify all players that game is now in progress
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{
		"state":           "in_progress",
		"starting_player": startingPlayerID,
	})

	if e.logger != nil {
		e.logger.Info("[MULLIGAN] Mulligan phase completed, game now IN_PROGRESS",
			zap.String("game_id", gameState.gameID),
			zap.String("new_state", gameState.state.String()),
			zap.String("priority_player", startingPlayerID),
		)
	}

	// Persist the game state immediately after mulligan completion.
	// This is critical for crash recovery - if the server restarts before the next
	// turn snapshot, the database will have the correct IN_PROGRESS state.
	// Use a goroutine to avoid lock issues (caller holds gameState.mu).
	gameID := gameState.gameID
	go func() {
		if e.logger != nil {
			e.logger.Info("[MULLIGAN] Persisting game state after mulligan completion (async)",
				zap.String("game_id", gameID),
			)
		}
		if err := e.PersistGameState(gameID); err != nil {
			if e.logger != nil {
				e.logger.Error("[MULLIGAN] FAILED to persist game state after mulligan completion",
					zap.String("game_id", gameID),
					zap.Error(err),
				)
			}
		} else if e.logger != nil {
			e.logger.Info("[MULLIGAN] Successfully persisted game state after mulligan completion",
				zap.String("game_id", gameID),
			)
		}
	}()
}

// handlePass handles a pass action
func (e *MageEngine) handlePass(gameState *engineGameState, playerID string) error {
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// RULES-LIGHT: Priority check removed - any player can pass at any time
	// The UI guides turn flow, but players control when to advance

	// Check for concessions before priority
	e.checkConcede(gameState)
	if e.checkIfGameIsOver(gameState) {
		return nil
	}

	// Per rule 117.5 and 603.3: Check state-based actions and triggered abilities before priority
	// Repeat until stable (SBA → triggers → repeat)
	e.checkStateAndTriggered(gameState)

	// Record replay state before priority
	// Per Java GameImpl.playPriority() line 1740: saveState(false) before each priority
	e.recordReplayState(gameState)

	player.Passed = true
	gameState.trackPriorityPass()
	gameState.trackAction()
	gameState.addMessage(fmt.Sprintf("%s passes", playerID), "action")

	// Check if all players who can respond have passed
	if gameState.allPassed() {
		// Resolve stack if not empty
		if !gameState.stack.IsEmpty() {
			err := e.resolveStack(gameState)
			// Check for concessions after stack resolution
			e.checkConcede(gameState)
			if e.checkIfGameIsOver(gameState) {
				return nil
			}
			return err
		}

		// Advance step/phase
		nextPlayer := e.getNextPlayer(gameState)
		oldTurn := gameState.turnManager.TurnNumber()
		phase, step := gameState.turnManager.AdvanceStep(nextPlayer)
		newTurn := gameState.turnManager.TurnNumber()
		gameState.addMessage(fmt.Sprintf("Game advances to %s - %s", phase.String(), step.String()), "action")

		if e.logger != nil {
			e.logger.Info("[TURN_ADVANCE] Step advanced",
				zap.String("game_id", gameState.gameID),
				zap.Int("old_turn", oldTurn),
				zap.Int("new_turn", newTurn),
				zap.String("phase", phase.String()),
				zap.String("step", step.String()),
				zap.String("next_active_player", nextPlayer),
				zap.String("actual_active_player", gameState.turnManager.ActivePlayer()),
			)
		}

		// Save turn snapshot if we advanced to a new turn
		// Per Java GameImpl.saveRollBackGameState(): save at start of each turn
		if newTurn > oldTurn {
			gameState.mu.Unlock() // Temporarily unlock to call SaveTurnSnapshot
			e.SaveTurnSnapshot(gameState.gameID, newTurn)
			gameState.mu.Lock() // Re-acquire lock
		}

		// Cleanup continuous effects at end of turn
		// Per Java: ContinuousEffects.removeEndOfTurnEffects() in cleanup step
		if step == rules.StepCleanup && gameState.layerSystem != nil {
			effects.CleanupEndOfTurnEffects(gameState.layerSystem)
		}

		// Get active player
		activePlayerID := gameState.turnManager.ActivePlayer()

		// Handle step-specific actions (untap, draw, cleanup, etc.)
		// Per MTG Rules 502-514 for turn-based actions
		e.handleStepBegin(gameState, step, activePlayerID)

		// Handle combat step initialization
		// Per Java BeginCombatStep.beginStep() and DeclareAttackersStep.beginStep()
		e.handleCombatStepBegin(gameState, step, activePlayerID)

		// Notify phase change
		e.notifyPhaseChange(gameState.gameID, map[string]interface{}{
			"phase":         phase.String(),
			"step":          step.String(),
			"active_player": activePlayerID,
			"turn":          gameState.turnManager.TurnNumber(),
		})

		// Reset pass flags (preserves lost/left player state)
		gameState.resetPassed()

		// Set priority to active player
		gameState.turnManager.SetPriority(activePlayerID)
		gameState.players[activePlayerID].HasPriority = true

		// Notify priority change
		e.notifyPriorityChange(gameState.gameID, activePlayerID, map[string]interface{}{
			"active_player": activePlayerID,
			"phase":         gameState.turnManager.CurrentPhase().String(),
			"step":          gameState.turnManager.CurrentStep().String(),
		})

		// Per rule 117.5: Check state-based actions before priority
		// Repeat until no more state-based actions occur
		for e.checkStateBasedActions(gameState) {
			// Continue checking until stable
		}

		// Emit phase/step change events
		gameState.eventBus.Publish(rules.NewEvent(rules.EventChangePhase, "", "", activePlayerID))
		gameState.eventBus.Publish(rules.NewEvent(rules.EventChangeStep, "", "", activePlayerID))

		// Check if active player should auto-pass (F6 etc.)
		if e.shouldAutoPass(gameState, activePlayerID) {
			return e.handlePass(gameState, activePlayerID)
		}
	} else {
		// Pass priority to next player
		nextPlayerID := e.getNextPlayerWithPriority(gameState, playerID)
		if nextPlayerID == "" {
			// No valid next player, all players who can respond have passed
			if gameState.allPassed() {
				if !gameState.stack.IsEmpty() {
					err := e.resolveStack(gameState)
					// Check for concessions after stack resolution
					e.checkConcede(gameState)
					if e.checkIfGameIsOver(gameState) {
						return nil
					}
					return err
				}
			}
			// Advance step/phase
			nextPlayer := e.getNextPlayer(gameState)
			phase, step := gameState.turnManager.AdvanceStep(nextPlayer)
			gameState.addMessage(fmt.Sprintf("Game advances to %s - %s", phase.String(), step.String()), "action")
			// Reset pass flags (preserves lost/left player state)
			gameState.resetPassed()
			// Set priority to active player
			activePlayerID := gameState.turnManager.ActivePlayer()

			// Handle step-specific actions (untap, draw, cleanup, etc.)
			// Per MTG Rules 502-514 for turn-based actions
			e.handleStepBegin(gameState, step, activePlayerID)

			// Handle combat step initialization
			// Per Java BeginCombatStep.beginStep() and DeclareAttackersStep.beginStep()
			e.handleCombatStepBegin(gameState, step, activePlayerID)

			// Per rule 117.5: Check state-based actions before priority
			// Repeat until no more state-based actions occur
			for e.checkStateBasedActions(gameState) {
				// Continue checking until stable
			}

			gameState.turnManager.SetPriority(activePlayerID)
			gameState.players[activePlayerID].HasPriority = true

			// Notify phase change and priority
			e.notifyPhaseChange(gameState.gameID, map[string]interface{}{
				"phase":         phase.String(),
				"step":          step.String(),
				"active_player": activePlayerID,
				"turn":          gameState.turnManager.TurnNumber(),
			})
			e.notifyPriorityChange(gameState.gameID, activePlayerID, map[string]interface{}{
				"active_player": activePlayerID,
				"phase":         phase.String(),
				"step":          step.String(),
			})
			return nil
		}
		// Per rule 117.5: Check state-based actions before priority
		// Repeat until no more state-based actions occur
		for e.checkStateBasedActions(gameState) {
			// Continue checking until stable
		}

		player.HasPriority = false
		gameState.turnManager.SetPriority(nextPlayerID)
		gameState.players[nextPlayerID].HasPriority = true
		gameState.players[nextPlayerID].Passed = false
		gameState.addPrompt(nextPlayerID, "You have priority. Pass?", []string{"PASS", "CAST"})

		// Notify priority change to update clients
		e.notifyPriorityChange(gameState.gameID, nextPlayerID, map[string]interface{}{
			"active_player": gameState.turnManager.ActivePlayer(),
			"phase":         gameState.turnManager.CurrentPhase().String(),
			"step":          gameState.turnManager.CurrentStep().String(),
		})

		// Check if next player should auto-pass (F6 etc.)
		if e.shouldAutoPass(gameState, nextPlayerID) {
			// Auto-pass for this player
			return e.handlePass(gameState, nextPlayerID)
		}
	}

	return nil
}

// handlePassUntil sets up auto-pass mode for a player
// The player will automatically pass priority until the specified condition is met
func (e *MageEngine) handlePassUntil(gameState *engineGameState, playerID string, passType PassUntilType) error {
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// RULES-LIGHT: Priority check removed - players can set auto-pass at any time

	// Set up pass-until mode
	player.PassUntil = passType

	// Record the current turn number for turn-based auto-pass modes
	// This is used to detect when the turn has advanced
	switch passType {
	case PassUntilMyNextTurn:
		// Pass until the next time it's this player's turn and we're in upkeep
		player.PassUntilTurn = gameState.turnManager.TurnNumber()
	case PassUntilNextTurn:
		// Pass until any new turn begins (next player's turn)
		player.PassUntilTurn = gameState.turnManager.TurnNumber()
		if e.logger != nil {
			e.logger.Info("[PASS_UNTIL_NEXT_TURN] Setting up auto-pass",
				zap.String("player_id", playerID),
				zap.Int("pass_until_turn", player.PassUntilTurn),
				zap.Int("current_turn", gameState.turnManager.TurnNumber()),
				zap.String("current_step", gameState.turnManager.CurrentStep().String()),
				zap.String("active_player", gameState.turnManager.ActivePlayer()),
			)
		}
	case PassUntilEndOfTurn:
		// Pass until cleanup step of the current turn
		player.PassUntilTurn = gameState.turnManager.TurnNumber()
	}

	var passDesc string
	switch passType {
	case PassUntilEndOfTurn:
		passDesc = "until end of turn"
	case PassUntilNextTurn:
		passDesc = "until next turn"
	case PassUntilStackResolved:
		passDesc = "until stack resolves"
	case PassUntilMyNextTurn:
		passDesc = "until their next turn"
	}

	gameState.addMessage(fmt.Sprintf("%s will pass %s", playerID, passDesc), "action")

	// Immediately pass priority
	return e.handlePass(gameState, playerID)
}

// shouldAutoPass checks if a player should automatically pass based on their PassUntil mode
func (e *MageEngine) shouldAutoPass(gameState *engineGameState, playerID string) bool {
	player, exists := gameState.players[playerID]
	if !exists || player.PassUntil == PassUntilNone {
		return false
	}

	switch player.PassUntil {
	case PassUntilEndOfTurn:
		// Auto-pass until cleanup step of current turn
		// Stop at cleanup step (end of turn) on the same turn
		currentTurn := gameState.turnManager.TurnNumber()
		isCleanup := gameState.turnManager.CurrentStep() == rules.StepCleanup

		// If turn has advanced beyond when we set this, stop auto-passing
		if currentTurn > player.PassUntilTurn {
			player.PassUntil = PassUntilNone
			return false
		}

		// Stop auto-passing when we reach cleanup on the same turn
		if isCleanup && currentTurn == player.PassUntilTurn {
			player.PassUntil = PassUntilNone
			return false
		}
		return true

	case PassUntilNextTurn:
		// Auto-pass until a new turn begins (turn number increases)
		// This passes through all remaining phases of the current turn
		// and stops at the beginning of the NEXT turn (any player's turn)
		currentTurn := gameState.turnManager.TurnNumber()
		currentStep := gameState.turnManager.CurrentStep()
		activePlayer := gameState.turnManager.ActivePlayer()

		if e.logger != nil {
			e.logger.Info("[PASS_UNTIL_NEXT_TURN] shouldAutoPass check",
				zap.String("player_id", playerID),
				zap.Int("current_turn", currentTurn),
				zap.Int("pass_until_turn", player.PassUntilTurn),
				zap.String("current_step", currentStep.String()),
				zap.String("active_player", activePlayer),
				zap.Bool("should_stop", currentTurn > player.PassUntilTurn),
			)
		}

		// Stop auto-passing when the turn number has advanced
		if currentTurn > player.PassUntilTurn {
			if e.logger != nil {
				e.logger.Info("[PASS_UNTIL_NEXT_TURN] STOPPING auto-pass - turn advanced",
					zap.String("player_id", playerID),
					zap.Int("current_turn", currentTurn),
					zap.Int("pass_until_turn", player.PassUntilTurn),
				)
			}
			player.PassUntil = PassUntilNone
			return false
		}
		return true

	case PassUntilStackResolved:
		// Auto-pass until stack is empty
		if gameState.stack.IsEmpty() {
			player.PassUntil = PassUntilNone
			return false
		}
		return true

	case PassUntilMyNextTurn:
		// Auto-pass until it's this player's upkeep again
		isMyTurn := gameState.turnManager.ActivePlayer() == playerID
		isUpkeep := gameState.turnManager.CurrentStep() == rules.StepUpkeep
		turnAdvanced := gameState.turnManager.TurnNumber() > player.PassUntilTurn

		// Stop auto-passing when it's our upkeep (and turn has advanced)
		if isMyTurn && isUpkeep && turnAdvanced {
			player.PassUntil = PassUntilNone // Clear the auto-pass
			return false
		}
		return true
	}

	return false
}

// clearPassUntilOnAction clears a player's auto-pass mode when they take an action
func (e *MageEngine) clearPassUntilOnAction(gameState *engineGameState, playerID string) {
	if player, exists := gameState.players[playerID]; exists {
		player.PassUntil = PassUntilNone
	}
}

// handleStringAction handles SEND_STRING type actions (spell casting or passing)
func (e *MageEngine) handleStringAction(gameState *engineGameState, action PlayerAction) error {
	spellName, ok := action.Data.(string)
	if !ok {
		return fmt.Errorf("SEND_STRING data must be string")
	}

	// Check for special action strings first
	spellNameUpper := strings.ToUpper(strings.TrimSpace(spellName))

	// Handle mulligan phase actions
	if spellNameUpper == "KEEP" {
		return e.handleKeepHand(gameState, action.PlayerID)
	}
	if spellNameUpper == "MULLIGAN" {
		return e.handleMulligan(gameState, action.PlayerID)
	}

	// Check if this is a pass action (some tests use "Pass" as SEND_STRING)
	if spellNameUpper == "PASS" {
		return e.handlePass(gameState, action.PlayerID)
	}

	// RULES-LIGHT: Handle direct manipulation commands
	// These allow players to directly manipulate game state
	if strings.HasPrefix(spellNameUpper, "TAP:") {
		cardID := strings.TrimPrefix(spellName, "TAP:")
		cardID = strings.TrimPrefix(cardID, "tap:")
		return e.handleDirectTap(gameState, cardID, true)
	}
	if strings.HasPrefix(spellNameUpper, "UNTAP:") {
		cardID := strings.TrimPrefix(spellName, "UNTAP:")
		cardID = strings.TrimPrefix(cardID, "untap:")
		return e.handleDirectTap(gameState, cardID, false)
	}
	if spellNameUpper == "UNTAP_ALL" {
		return e.handleDirectUntapAll(gameState, action.PlayerID)
	}
	if strings.HasPrefix(spellNameUpper, "FLIP:") {
		parts := strings.SplitN(spellName[5:], ":", 2)
		if len(parts) < 1 {
			return fmt.Errorf("FLIP requires cardId")
		}
		faceDown := len(parts) < 2 || parts[1] == "true"
		return e.handleDirectFlip(gameState, parts[0], faceDown)
	}
	if strings.HasPrefix(spellNameUpper, "TRANSFORM:") {
		cardID := strings.TrimPrefix(spellName, "TRANSFORM:")
		cardID = strings.TrimPrefix(cardID, "transform:")
		return e.handleDirectTransform(gameState, cardID)
	}
	if strings.HasPrefix(spellNameUpper, "MOVE:") {
		parts := strings.SplitN(spellName[5:], ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("MOVE requires cardId:zone")
		}
		return e.handleDirectMove(gameState, action.PlayerID, parts[0], parts[1])
	}
	if strings.HasPrefix(spellNameUpper, "SET_COUNTER:") {
		parts := strings.SplitN(spellName[12:], ":", 3)
		if len(parts) != 3 {
			return fmt.Errorf("SET_COUNTER requires cardId:type:amount")
		}
		amount := 0
		fmt.Sscanf(parts[2], "%d", &amount)
		return e.handleDirectSetCounter(gameState, parts[0], parts[1], amount)
	}
	if strings.HasPrefix(spellNameUpper, "MODIFY_COUNTER:") {
		parts := strings.SplitN(spellName[15:], ":", 3)
		if len(parts) != 3 {
			return fmt.Errorf("MODIFY_COUNTER requires cardId:type:delta")
		}
		delta := 0
		fmt.Sscanf(parts[2], "%d", &delta)
		return e.handleDirectModifyCounter(gameState, parts[0], parts[1], delta)
	}
	if strings.HasPrefix(spellNameUpper, "CREATE_TOKEN:") {
		// Format: CREATE_TOKEN:name:types:power:toughness:color:abilities(comma-sep)
		parts := strings.SplitN(spellName[13:], ":", 6)
		if len(parts) < 5 {
			return fmt.Errorf("CREATE_TOKEN requires name:types:power:toughness:color[:abilities]")
		}
		abilities := []string{}
		if len(parts) >= 6 && parts[5] != "" {
			abilities = strings.Split(parts[5], ",")
		}
		return e.handleDirectCreateToken(gameState, action.PlayerID, parts[0], parts[1], parts[2], parts[3], parts[4], abilities)
	}
	if strings.HasPrefix(spellNameUpper, "DESTROY_TOKEN:") {
		cardID := strings.TrimPrefix(spellName, "DESTROY_TOKEN:")
		cardID = strings.TrimPrefix(cardID, "destroy_token:")
		return e.handleDirectDestroyToken(gameState, cardID)
	}
	if strings.HasPrefix(spellNameUpper, "SET_LIFE:") {
		parts := strings.SplitN(spellName[9:], ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("SET_LIFE requires playerId:amount")
		}
		amount := 0
		fmt.Sscanf(parts[1], "%d", &amount)
		return e.handleDirectSetLife(gameState, parts[0], amount)
	}
	if strings.HasPrefix(spellNameUpper, "MODIFY_LIFE:") {
		parts := strings.SplitN(spellName[12:], ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("MODIFY_LIFE requires playerId:delta")
		}
		delta := 0
		fmt.Sscanf(parts[1], "%d", &delta)
		return e.handleDirectModifyLife(gameState, parts[0], delta)
	}
	if strings.HasPrefix(spellNameUpper, "DRAW:") {
		parts := strings.SplitN(spellName[5:], ":", 2)
		count := 1
		if len(parts) == 2 {
			fmt.Sscanf(parts[1], "%d", &count)
		}
		playerID := action.PlayerID
		if len(parts) >= 1 && parts[0] != "" {
			playerID = parts[0]
		}
		return e.handleDirectDraw(gameState, playerID, count)
	}

	// RULES-LIGHT: Handle stack tracking (card stays in place, just added to visual stack)
	if strings.HasPrefix(spellNameUpper, "STACK_ADD:") {
		cardID := strings.TrimPrefix(spellName, "STACK_ADD:")
		cardID = strings.TrimPrefix(cardID, "stack_add:")
		return e.handleDirectStackAdd(gameState, action.PlayerID, cardID)
	}
	if strings.HasPrefix(spellNameUpper, "STACK_REMOVE:") {
		itemID := strings.TrimPrefix(spellName, "STACK_REMOVE:")
		itemID = strings.TrimPrefix(itemID, "stack_remove:")
		return e.handleDirectStackRemove(gameState, action.PlayerID, itemID)
	}

	// RULES-LIGHT: Handle turn and phase control commands
	if spellNameUpper == "NEXT_TURN" {
		return e.handleDirectNextTurn(gameState, action.PlayerID)
	}
	if spellNameUpper == "CLEAR_COMBAT" {
		return e.handleDirectClearCombat(gameState)
	}
	if spellNameUpper == "SHUFFLE" {
		return e.handleDirectShuffle(gameState, action.PlayerID)
	}
	if strings.HasPrefix(spellNameUpper, "SHUFFLE:") {
		playerID := strings.TrimPrefix(spellName, "SHUFFLE:")
		playerID = strings.TrimPrefix(playerID, "shuffle:")
		return e.handleDirectShuffle(gameState, playerID)
	}
	if strings.HasPrefix(spellNameUpper, "SET_PLAYER_COUNTER:") {
		parts := strings.SplitN(spellName[19:], ":", 3)
		if len(parts) != 3 {
			return fmt.Errorf("SET_PLAYER_COUNTER requires playerId:type:amount")
		}
		amount := 0
		fmt.Sscanf(parts[2], "%d", &amount)
		return e.handleDirectSetPlayerCounter(gameState, parts[0], parts[1], amount)
	}
	if strings.HasPrefix(spellNameUpper, "SEARCH_LIBRARY:") {
		parts := strings.SplitN(spellName[15:], ":", 3)
		if len(parts) < 2 {
			return fmt.Errorf("SEARCH_LIBRARY requires destination:shuffle[:message]")
		}
		shuffle := true
		if len(parts) >= 2 && strings.ToLower(parts[1]) == "false" {
			shuffle = false
		}
		message := ""
		if len(parts) >= 3 {
			message = parts[2]
		}
		return e.handleDirectSearchLibrary(gameState, action.PlayerID, parts[0], shuffle, message)
	}

	playerID := action.PlayerID
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// RULES-LIGHT: Priority check removed - any player can cast spells at any time
	// Players are trusted to follow timing rules

	// Check state-based actions (advisory - may provide UI hints)
	e.checkStateAndTriggered(gameState)

	// Find card in hand
	var card *internalCard
	for _, c := range player.Hand {
		if c == nil {
			continue
		}
		if strings.EqualFold(c.Name, spellName) {
			card = c
			break
		}
	}

	if card == nil {
		return fmt.Errorf("card %s not found in hand", spellName)
	}

	// Rules-light approach: We don't enforce targeting.
	// Per MAGE_ENGINE_ARCHITECTURE.md: "Assist, don't enforce"
	// Players communicate targets via chat and/or ping cards to indicate targets.
	// Just inform them what the spell targets (for reference) and proceed with cast.
	targetRequirements := targeting.ParseTargetRequirements(card.Type, card.RulesText)
	if len(targetRequirements) > 0 {
		req := targetRequirements[0]
		gameState.addMessage(fmt.Sprintf("💡 %s targets: %s (communicate via chat)", card.Name, req.Description), "info")
	}

	// Proceed directly with casting - players handle targeting themselves
	return e.proceedWithSpellCast(gameState, playerID, card)
}

// proceedWithSpellCast handles the actual spell casting after targets are selected (if any)
// This moves the card to stack and sets up the resolution callback
func (e *MageEngine) proceedWithSpellCast(gameState *engineGameState, playerID string, card *internalCard) error {
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// Parse and pay mana cost (best-effort, rules-light approach)
	// Per MAGE_ENGINE_ARCHITECTURE.md: We don't enforce mana costs, just assist
	if card.ManaCost != "" {
		cost, err := mana.ParseCost(card.ManaCost)
		if err != nil {
			gameState.addMessage(fmt.Sprintf("⚠️ %s: Could not parse mana cost %s", card.Name, card.ManaCost), "warning")
		} else {
			// Notify players if mana is insufficient (but allow the cast)
			if !cost.CanPay(player.ManaPool, 0) {
				gameState.addMessage(fmt.Sprintf("⚠️ %s cast without paying full mana cost %s", card.Name, card.ManaCost), "warning")
			}

			// Pay whatever mana is available (best-effort)
			_ = e.payManaCost(player.ManaPool, cost)
		}
	}

	// Move card from hand to stack
	player.Hand = e.removeCardFromSlice(player.Hand, card.ID)
	card.Zone = zoneStack

	// Create stack item with resolve function that looks up card by ID
	// This ensures we get the current card reference, not a stale closure
	cardID := card.ID
	stackItem := rules.StackItem{
		ID:          card.ID,
		Controller:  playerID,
		Description: fmt.Sprintf("%s casts %s", playerID, card.Name),
		Kind:        rules.StackItemKindSpell,
		SourceID:    card.ID,
		Metadata:    make(map[string]string),
		Resolve: func() error {
			// Look up card by ID to ensure we have the current reference
			resolveCard, found := gameState.cards[cardID]
			if !found {
				return fmt.Errorf("card %s not found in game state", cardID)
			}
			return e.resolveSpell(gameState, resolveCard)
		},
	}

	// Copy any target information to stack item metadata
	if targets, ok := card.Metadata["targets"]; ok && targets != "" {
		stackItem.Metadata["targets"] = targets
	}

	gameState.stack.Push(stackItem)
	gameState.trackStackItem()
	gameState.trackStackDepth()
	gameState.trackSpellCast()
	gameState.trackAction()
	gameState.addMessage(fmt.Sprintf("%s casts %s", playerID, card.Name), "action")

	// Notify stack update
	e.notifyStackUpdate(gameState.gameID, map[string]interface{}{
		"action":      "spell_cast",
		"player_id":   playerID,
		"card_name":   card.Name,
		"card_id":     cardID,
		"stack_depth": len(gameState.stack.List()),
	})

	// Emit spell cast event
	spellCastEvent := rules.Event{
		Type:        rules.EventSpellCast,
		ID:          uuid.New().String(),
		TargetID:    card.ID,
		SourceID:    card.ID,
		Controller:  playerID,
		PlayerID:    playerID,
		Timestamp:   time.Now(),
		Metadata:    make(map[string]string),
		Description: fmt.Sprintf("%s casts %s", playerID, card.Name),
	}
	gameState.eventBus.Publish(spellCastEvent)

	// Check for triggered abilities (e.g., "whenever you cast a spell")
	// Create a triggered ability for Lightning Bolt (for testing - simulates a "Sanctuary" effect)
	// Triggered abilities go on top of the stack (LIFO - last in, first out)
	e.createTriggeredAbilityForSpell(gameState, card, playerID)

	// Per MTG rules 117.3c: After a player casts a spell, activates an ability, or takes a special action,
	// that player retains priority and may take another action. Priority only passes when the player
	// explicitly passes or when a spell/ability resolves.
	// Reset all players' passed flags (preserves lost/left player state)
	gameState.resetPassed()

	// Per rule 117.5 and 603.3: Check state-based actions and triggered abilities before priority
	// Repeat until stable (SBA → triggers → repeat)
	e.checkStateAndTriggered(gameState)

	// Caster retains priority after casting
	player.HasPriority = true
	player.Passed = false
	gameState.turnManager.SetPriority(playerID)
	gameState.addPrompt(playerID, "You have priority. Cast another spell or pass?", []string{"PASS", "CAST"})

	// Notify clients that caster retains priority after casting
	// Per MTG rule 117.3c: After casting, the player retains priority
	e.notifyPriorityChange(gameState.gameID, playerID, map[string]interface{}{
		"active_player": gameState.turnManager.ActivePlayer(),
		"phase":         gameState.turnManager.CurrentPhase().String(),
		"step":          gameState.turnManager.CurrentStep().String(),
		"reason":        "spell_cast",
	})

	return nil
}

// payManaCost pays a mana cost from a player's mana pool (best-effort)
// Per MAGE_ENGINE_ARCHITECTURE.md: Rules-light approach - pay what's available, don't enforce
// It pays colored mana first, then generic mana from any available source
func (e *MageEngine) payManaCost(pool *mana.ManaPool, cost *mana.ManaCost) error {
	// Pay colored mana first - spend what's available (best-effort)
	if cost.White > 0 {
		pool.Spend(mana.ManaWhite, cost.White) // Best-effort, ignore if insufficient
	}
	if cost.Blue > 0 {
		pool.Spend(mana.ManaBlue, cost.Blue)
	}
	if cost.Black > 0 {
		pool.Spend(mana.ManaBlack, cost.Black)
	}
	if cost.Red > 0 {
		pool.Spend(mana.ManaRed, cost.Red)
	}
	if cost.Green > 0 {
		pool.Spend(mana.ManaGreen, cost.Green)
	}
	if cost.Colorless > 0 {
		pool.Spend(mana.ManaColorless, cost.Colorless)
	}

	// Pay generic mana from any available source
	// Prefer colorless first, then colors in WUBRG order
	genericRemaining := cost.Generic
	if genericRemaining > 0 {
		// Try colorless first
		colorlessAvail := pool.GetTotal(mana.ManaColorless)
		if colorlessAvail > 0 {
			spend := colorlessAvail
			if spend > genericRemaining {
				spend = genericRemaining
			}
			pool.Spend(mana.ManaColorless, spend)
			genericRemaining -= spend
		}

		// Then try each color
		manaTypes := []mana.ManaType{mana.ManaWhite, mana.ManaBlue, mana.ManaBlack, mana.ManaRed, mana.ManaGreen}
		for _, mt := range manaTypes {
			if genericRemaining <= 0 {
				break
			}
			avail := pool.GetTotal(mt)
			if avail > 0 {
				spend := avail
				if spend > genericRemaining {
					spend = genericRemaining
				}
				pool.Spend(mt, spend)
				genericRemaining -= spend
			}
		}
		// Rules-light: Don't error if generic cost couldn't be fully paid
	}

	return nil
}

// handleIntegerAction handles SEND_INTEGER type actions
func (e *MageEngine) handleIntegerAction(gameState *engineGameState, action PlayerAction) error {
	var value int
	switch v := action.Data.(type) {
	case int:
		value = v
	case int32:
		value = int(v)
	case int64:
		value = int(v)
	case float64:
		value = int(v)
	case float32:
		value = int(v)
	default:
		return fmt.Errorf("SEND_INTEGER data must be numeric, got %T", action.Data)
	}

	playerID := action.PlayerID

	// Check if there's a pending X value request
	if gameState.pendingXValueRequest != nil {
		req := gameState.pendingXValueRequest
		if req.PlayerID != playerID {
			return fmt.Errorf("X value request is for player %s, not %s", req.PlayerID, playerID)
		}

		// Validate X value is within bounds
		if value < req.MinValue {
			return fmt.Errorf("X value %d is below minimum %d", value, req.MinValue)
		}
		if value > req.MaxValue {
			return fmt.Errorf("X value %d is above maximum %d", value, req.MaxValue)
		}

		if e.logger != nil {
			e.logger.Info("received X value selection",
				zap.String("player", playerID),
				zap.String("source", req.SourceName),
				zap.Int("value", value))
		}

		// Clear the pending request before calling OnComplete
		onComplete := req.OnComplete
		gameState.pendingXValueRequest = nil

		// Execute the callback with the chosen value
		if onComplete != nil {
			if err := onComplete(value); err != nil {
				return fmt.Errorf("failed to complete X value selection: %w", err)
			}
		}

		return nil
	}

	// No pending X value request - legacy behavior (for testing/compatibility)
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// For now, treat integer as life change (for testing)
	oldLife := player.Life
	player.Life += value
	gameState.addMessage(fmt.Sprintf("%s's life changes by %d (now %d)", playerID, value, player.Life), "life")

	// Emit life change event
	if value < 0 {
		gameState.eventBus.Publish(rules.Event{
			Type:        rules.EventLostLife,
			ID:          uuid.New().String(),
			TargetID:    playerID,
			PlayerID:    playerID,
			Amount:      -value,
			Timestamp:   time.Now(),
			Metadata:    make(map[string]string),
			Description: fmt.Sprintf("%s's life changes from %d to %d", playerID, oldLife, player.Life),
		})
	} else {
		gameState.eventBus.Publish(rules.Event{
			Type:        rules.EventGainedLife,
			ID:          uuid.New().String(),
			TargetID:    playerID,
			PlayerID:    playerID,
			Amount:      value,
			Timestamp:   time.Now(),
			Metadata:    make(map[string]string),
			Description: fmt.Sprintf("%s's life changes from %d to %d", playerID, oldLife, player.Life),
		})
	}

	return nil
}

// handleUUIDAction handles SEND_UUID type actions (e.g., selecting targets, countering spells)
func (e *MageEngine) handleUUIDAction(gameState *engineGameState, action PlayerAction) error {
	uuidStr, ok := action.Data.(string)
	if !ok {
		return fmt.Errorf("SEND_UUID data must be string")
	}

	playerID := action.PlayerID
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// Check if there's a pending target request for this player
	if gameState.pendingTargetRequest != nil && gameState.pendingTargetRequest.PlayerID == playerID {
		return e.handleTargetSelection(gameState, playerID, uuidStr)
	}

	// Check if there's a pending library search for this player
	if gameState.pendingLibrarySearch != nil && gameState.pendingLibrarySearch.PlayerID == playerID {
		return e.handleLibrarySearchSelection(gameState, player, uuidStr)
	}

	// Per MAGE_ENGINE_ARCHITECTURE.md: Rules-light - don't enforce priority
	if gameState.turnManager.PriorityPlayer() != playerID {
		gameState.addMessage(fmt.Sprintf("⚠️ %s acts without priority", playerID), "warning")
	}

	// Per rule 117.5 and 603.3: Check state-based actions and triggered abilities before priority
	// Repeat until stable (SBA → triggers → repeat)
	e.checkStateAndTriggered(gameState)

	// Check if UUID refers to a spell on the stack that can be countered
	stackItems := gameState.stack.List()
	for _, item := range stackItems {
		if item.ID == uuidStr || item.SourceID == uuidStr {
			// Counter the spell by removing it from stack
			removedItem, found := gameState.stack.Remove(item.ID)
			if found {
				gameState.addMessage(fmt.Sprintf("%s counters %s", playerID, removedItem.Description), "action")

				// Move countered spell to graveyard
				if card, found := gameState.cards[removedItem.SourceID]; found {
					card.Zone = zoneGraveyard
					if controller, exists := gameState.players[removedItem.Controller]; exists {
						controller.Graveyard = append(controller.Graveyard, card)
					}
				}

				// Emit counter event
				gameState.eventBus.Publish(rules.Event{
					Type:        rules.EventStackItemRemoved,
					ID:          uuid.New().String(),
					TargetID:    removedItem.ID,
					SourceID:    removedItem.SourceID,
					Controller:  playerID,
					Timestamp:   time.Now(),
					Description: fmt.Sprintf("%s counters %s", playerID, removedItem.Description),
				})

				// Pass priority to next player
				nextPlayerID := e.getNextPlayerWithPriority(gameState, playerID)
				if nextPlayerID != "" && nextPlayerID != playerID {
					// Per rule 117.5: Check state-based actions before priority
					// Repeat until no more state-based actions occur
					for e.checkStateBasedActions(gameState) {
						// Continue checking until stable
					}

					player.HasPriority = false
					gameState.turnManager.SetPriority(nextPlayerID)
					gameState.players[nextPlayerID].HasPriority = true
					gameState.players[nextPlayerID].Passed = false

					// Notify clients of priority change after countering
					e.notifyPriorityChange(gameState.gameID, nextPlayerID, map[string]interface{}{
						"active_player": gameState.turnManager.ActivePlayer(),
						"phase":         gameState.turnManager.CurrentPhase().String(),
						"step":          gameState.turnManager.CurrentStep().String(),
						"reason":        "spell_countered",
					})
				}

				return nil
			}
		}
	}

	return fmt.Errorf("UUID %s not found on stack", uuidStr)
}

// handleTargetSelection processes a target selection from a player
// This is called when there's a pending target request and the player sends a SEND_UUID
func (e *MageEngine) handleTargetSelection(gameState *engineGameState, playerID, targetID string) error {
	req := gameState.pendingTargetRequest
	if req == nil {
		return fmt.Errorf("no pending target request")
	}

	if req.PlayerID != playerID {
		return fmt.Errorf("target request is for player %s, not %s", req.PlayerID, playerID)
	}

	// Check if this is a cancel request (empty UUID or special cancel marker)
	if targetID == "" || targetID == "CANCEL" {
		if req.Required {
			return fmt.Errorf("target selection is required and cannot be cancelled")
		}
		// Call cancel callback
		if req.OnCancel != nil {
			if err := req.OnCancel(); err != nil {
				return fmt.Errorf("cancel callback failed: %v", err)
			}
		}
		gameState.pendingTargetRequest = nil
		gameState.addMessage(fmt.Sprintf("%s cancelled target selection", playerID), "action")

		if e.logger != nil {
			e.logger.Info("target selection cancelled",
				zap.String("game_id", gameState.gameID),
				zap.String("player_id", playerID),
			)
		}
		return nil
	}

	// Validate that the target is in the valid targets list
	isValid := false
	for _, validID := range req.ValidTargetIDs {
		if validID == targetID {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("target %s is not a valid target", targetID)
	}

	// Check if this target is already selected (toggle off)
	alreadySelected := false
	for i, selected := range req.SelectedTargetIDs {
		if selected == targetID {
			alreadySelected = true
			// Remove from selected targets
			req.SelectedTargetIDs = append(req.SelectedTargetIDs[:i], req.SelectedTargetIDs[i+1:]...)
			gameState.addMessage(fmt.Sprintf("%s deselected target", playerID), "action")
			break
		}
	}

	if !alreadySelected {
		// Check if we've already selected the maximum number of targets
		if len(req.SelectedTargetIDs) >= req.Requirement.MaxTargets {
			return fmt.Errorf("maximum targets already selected (%d)", req.Requirement.MaxTargets)
		}

		// Add to selected targets
		req.SelectedTargetIDs = append(req.SelectedTargetIDs, targetID)

		// Get target name for message
		targetName := targetID
		if card, found := gameState.cards[targetID]; found {
			targetName = card.Name
		}
		gameState.addMessage(fmt.Sprintf("%s selected target: %s", playerID, targetName), "action")
	}

	if e.logger != nil {
		e.logger.Debug("target selection updated",
			zap.String("game_id", gameState.gameID),
			zap.String("player_id", playerID),
			zap.String("target_id", targetID),
			zap.Int("selected_count", len(req.SelectedTargetIDs)),
			zap.Int("min_targets", req.Requirement.MinTargets),
			zap.Int("max_targets", req.Requirement.MaxTargets),
		)
	}

	// Check if we have enough targets to complete
	if len(req.SelectedTargetIDs) >= req.Requirement.MinTargets {
		// Check if we've hit max or if this is a single-target requirement
		if len(req.SelectedTargetIDs) >= req.Requirement.MaxTargets || req.Requirement.MaxTargets == 1 {
			// Auto-complete single target selections
			return e.completeTargetSelection(gameState)
		}
	}

	// Notify client of updated selection
	e.notifyTargetRequest(gameState.gameID, playerID, map[string]interface{}{
		"message":          req.Message,
		"targets":          req.ValidTargetIDs,
		"selected_targets": req.SelectedTargetIDs,
		"required":         req.Required,
		"min_targets":      req.Requirement.MinTargets,
		"max_targets":      req.Requirement.MaxTargets,
		"source_id":        req.SourceID,
		"can_confirm":      len(req.SelectedTargetIDs) >= req.Requirement.MinTargets,
	})

	return nil
}

// completeTargetSelection finalizes the target selection and calls the completion callback
func (e *MageEngine) completeTargetSelection(gameState *engineGameState) error {
	req := gameState.pendingTargetRequest
	if req == nil {
		return fmt.Errorf("no pending target request to complete")
	}

	// Validate final selection
	if len(req.SelectedTargetIDs) < req.Requirement.MinTargets {
		return fmt.Errorf("not enough targets selected: need %d, have %d", req.Requirement.MinTargets, len(req.SelectedTargetIDs))
	}
	if len(req.SelectedTargetIDs) > req.Requirement.MaxTargets {
		return fmt.Errorf("too many targets selected: max %d, have %d", req.Requirement.MaxTargets, len(req.SelectedTargetIDs))
	}

	// Copy the selected targets
	selectedTargets := make([]string, len(req.SelectedTargetIDs))
	copy(selectedTargets, req.SelectedTargetIDs)

	if e.logger != nil {
		e.logger.Info("completing target selection",
			zap.String("game_id", gameState.gameID),
			zap.String("player_id", req.PlayerID),
			zap.Strings("selected_targets", selectedTargets),
		)
	}

	// Call the completion callback
	if req.OnComplete != nil {
		if err := req.OnComplete(selectedTargets); err != nil {
			return fmt.Errorf("completion callback failed: %v", err)
		}
	}

	// Clear the pending request
	gameState.pendingTargetRequest = nil

	// Build target names for message
	targetNames := make([]string, 0, len(selectedTargets))
	for _, targetID := range selectedTargets {
		if card, found := gameState.cards[targetID]; found {
			targetNames = append(targetNames, card.Name)
		} else {
			targetNames = append(targetNames, targetID)
		}
	}
	gameState.addMessage(fmt.Sprintf("%s confirmed targets: %s", req.PlayerID, strings.Join(targetNames, ", ")), "action")

	return nil
}

// ConfirmTargetSelection manually completes target selection (e.g., when player clicks "Confirm")
// Used for multi-target scenarios where auto-complete doesn't trigger
func (e *MageEngine) ConfirmTargetSelection(gameID, playerID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	gameState, exists := e.games[gameID]
	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	if gameState.pendingTargetRequest == nil {
		return fmt.Errorf("no pending target request")
	}

	if gameState.pendingTargetRequest.PlayerID != playerID {
		return fmt.Errorf("target request is for player %s, not %s", gameState.pendingTargetRequest.PlayerID, playerID)
	}

	return e.completeTargetSelection(gameState)
}

// resolveStack resolves all items on the stack
func (e *MageEngine) resolveStack(gameState *engineGameState) error {
	// RULES-LIGHT: Resolve only the top item, not all items
	// This gives players control over resolution order and timing

	if gameState.stack.IsEmpty() {
		return nil
	}

	item, err := gameState.stack.Pop()
	if err != nil {
		return fmt.Errorf("failed to pop from stack: %w", err)
	}

	if e.logger != nil {
		e.logger.Debug("resolving top stack item (rules-light mode)",
			zap.String("item_id", item.ID),
			zap.String("source_id", item.SourceID),
			zap.String("description", item.Description),
			zap.Int("remaining_items", len(gameState.stack.List())),
		)
	}

	// RULES-LIGHT: Skip legality check - players handle fizzling manually
	// Resolve the item
	gameState.addMessage(fmt.Sprintf("%s resolves", item.Description), "action")

	if item.Resolve != nil {
		if err := item.Resolve(); err != nil {
			gameState.addMessage(fmt.Sprintf("Error resolving %s: %v", item.Description, err), "action")
			if e.logger != nil {
				e.logger.Error("failed to resolve stack item",
					zap.String("item_id", item.ID),
					zap.Error(err),
				)
			}
		} else {
			gameState.addMessage(fmt.Sprintf("%s resolved successfully", item.Description), "action")
		}
	}

	// Emit stack resolution event
	gameState.eventBus.Publish(rules.Event{
		Type:        rules.EventStackItemResolved,
		ID:          uuid.New().String(),
		TargetID:    item.ID,
		SourceID:    item.SourceID,
		Controller:  item.Controller,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("%s resolved", item.Description),
	})

	// Reset pass flags after stack resolution (preserves lost/left player state)
	gameState.resetPassed()

	// Priority returns to active player
	activePlayerID := gameState.turnManager.ActivePlayer()

	// Per Java GameImpl.resolve() lines 1857-1860: Process simultaneous events after stack resolution
	// This handles events that occurred during resolution (e.g., multiple creatures dying)
	for gameState.hasSimultaneousEvents() {
		e.handleSimultaneousEvents(gameState)
	}

	// Per rule 117.5 and 603.3: Check state-based actions and triggered abilities before priority
	// Repeat until stable (SBA → triggers → repeat)
	e.checkStateAndTriggered(gameState)

	gameState.turnManager.SetPriority(activePlayerID)
	gameState.players[activePlayerID].HasPriority = true
	gameState.addPrompt(activePlayerID, "You have priority. Pass?", []string{"PASS", "CAST"})

	// Notify clients that priority has changed after stack resolution
	// Per MTG rule 117.3b: After a spell/ability resolves, active player receives priority
	e.notifyPriorityChange(gameState.gameID, activePlayerID, map[string]interface{}{
		"active_player": activePlayerID,
		"phase":         gameState.turnManager.CurrentPhase().String(),
		"step":          gameState.turnManager.CurrentStep().String(),
		"reason":        "stack_resolved",
	})

	return nil
}

// handleSpecialAction handles SPECIAL_ACTION type actions (play land, foretell, etc.)
// Per MTG Rule 116: Special actions don't use the stack
func (e *MageEngine) handleSpecialAction(gameState *engineGameState, action PlayerAction) error {
	payload, ok := action.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("SPECIAL_ACTION data must be a map")
	}

	actionType, ok := payload["action_type"].(string)
	if !ok {
		return fmt.Errorf("action_type is required")
	}

	playerID := action.PlayerID
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// sourceID is optional for some actions (like SEARCH_LIBRARY)
	sourceID, _ := payload["source_id"].(string)

	switch actionType {
	case "PLAY_LAND":
		if sourceID == "" {
			return fmt.Errorf("source_id is required for PLAY_LAND")
		}
		return e.handlePlayLand(gameState, player, sourceID)
	case "ADVANCE_PHASE":
		return e.handleAdvancePhase(gameState, player)
	case "ACTIVATE_MANA_ABILITY":
		if sourceID == "" {
			return fmt.Errorf("source_id is required for ACTIVATE_MANA_ABILITY")
		}
		return e.handleActivateManaAbility(gameState, player, sourceID)
	case "SEARCH_LIBRARY":
		destination, _ := payload["destination"].(string)
		if destination == "" {
			destination = "hand" // Default to searching to hand
		}
		shuffle := true // Default to shuffling after search
		if shuffleVal, ok := payload["shuffle"].(bool); ok {
			shuffle = shuffleVal
		}
		message, _ := payload["message"].(string)
		if message == "" {
			message = "Search your library for a card"
		}
		return e.handleSearchLibrary(gameState, player, destination, shuffle, message)
	default:
		return fmt.Errorf("unknown special action: %s", actionType)
	}
}

// handleActivateAbilityAction handles ACTIVATE_ABILITY type actions
// Per MTG Rule 602: Activating an activated ability follows a specific sequence
func (e *MageEngine) handleActivateAbilityAction(gameState *engineGameState, action PlayerAction) error {
	payload, ok := action.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("ACTIVATE_ABILITY data must be a map")
	}

	cardID, ok := payload["card_id"].(string)
	if !ok {
		return fmt.Errorf("card_id is required")
	}

	abilityID, ok := payload["ability_id"].(string)
	if !ok {
		return fmt.Errorf("ability_id is required")
	}

	// Parse targets (may be empty)
	var targets []string
	if targetsRaw, ok := payload["targets"]; ok {
		if targetsSlice, ok := targetsRaw.([]interface{}); ok {
			for _, t := range targetsSlice {
				if str, ok := t.(string); ok {
					targets = append(targets, str)
				}
			}
		} else if targetsStrSlice, ok := targetsRaw.([]string); ok {
			targets = targetsStrSlice
		}
	}

	playerID := action.PlayerID
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	return e.handleActivateAbility(gameState, player, cardID, abilityID, targets)
}

// handleActivateAbility activates a non-mana activated ability
// Per MTG Rule 602: The full activation sequence
// Supports both registered abilities (UUID IDs) and rules-text-parsed abilities (cardID-ability-N format)
func (e *MageEngine) handleActivateAbility(gameState *engineGameState, player *internalPlayer, cardID, abilityIDStr string, targets []string) error {
	// RULES-LIGHT: Priority and controller checks removed - any player can activate abilities
	// Players are trusted to follow the rules

	// Find the card/permanent on battlefield
	var card *internalCard
	for _, c := range gameState.battlefield {
		if c.ID == cardID {
			card = c
			break
		}
	}
	if card == nil {
		return fmt.Errorf("permanent %s not found on battlefield", cardID)
	}

	// Check if this is a synthetic ability ID from rules text parsing (format: cardID-ability-N)
	var abilityText string
	var hasTapCost bool
	var isSorcerySpeed bool
	var resolveFunc func() error

	if strings.Contains(abilityIDStr, "-ability-") {
		// This is a synthetic ability parsed from rules text
		// Extract the ability index
		parts := strings.Split(abilityIDStr, "-ability-")
		if len(parts) != 2 {
			return fmt.Errorf("invalid synthetic ability ID format: %s", abilityIDStr)
		}
		abilityIndex := 0
		if _, err := fmt.Sscanf(parts[1], "%d", &abilityIndex); err != nil {
			return fmt.Errorf("invalid ability index in ID: %s", abilityIDStr)
		}

		// Parse activated abilities from rules text
		parsedAbilities := parseActivatedAbilitiesFromText(card.RulesText)
		if abilityIndex >= len(parsedAbilities) {
			return fmt.Errorf("ability index %d out of range for card %s", abilityIndex, card.Name)
		}

		parsedAbility := parsedAbilities[abilityIndex]
		abilityText = parsedAbility.FullText
		hasTapCost = parsedAbility.HasTapCost
		isSorcerySpeed = parsedAbility.IsSorcerySpeed

		// Check if ability needs X value selection
		if parsedAbility.HasXValue {
			// Determine max X based on ability type
			maxX := 0
			if parsedAbility.XValueType == "reveal_white" {
				// Count white cards in hand for Martyr of Sands
				for _, c := range player.Hand {
					if isCardWhite(c) {
						maxX++
					}
				}
			} else {
				// Default max based on available mana
				if player.ManaPool != nil {
					maxX = player.ManaPool.GetTotalMana()
				}
			}

			if maxX > 0 {
				// Check if we're waiting for X value
				if gameState.pendingXValueRequest != nil && gameState.pendingXValueRequest.PlayerID == player.PlayerID {
					// Already waiting for this player - they shouldn't be activating again
					return fmt.Errorf("already waiting for X value selection")
				}

				// Create ability activation context for callback
				capturedCard := card
				capturedPlayer := player
				capturedAbility := parsedAbility
				capturedGameState := gameState

				// Store X value when player chooses it
				var chosenX int

				// Create the pending X value request
				message := fmt.Sprintf("Choose X value for %s", card.Name)
				if parsedAbility.XValueType == "reveal_white" {
					message = fmt.Sprintf("How many white cards to reveal? (0-%d)", maxX)
				}

				gameState.pendingXValueRequest = &PendingXValueRequest{
					PlayerID:   player.PlayerID,
					SourceID:   cardID,
					SourceName: card.Name,
					Message:    message,
					MinValue:   0,
					MaxValue:   maxX,
					Timestamp:  time.Now(),
					OnComplete: func(xValue int) error {
						chosenX = xValue
						if e.logger != nil {
							e.logger.Info("X value chosen, completing ability activation",
								zap.String("card", capturedCard.Name),
								zap.Int("x_value", chosenX))
						}

						// Store X value in the ability for later use during resolution
						capturedAbility.chosenXValue = chosenX

						// Continue with ability activation - pay costs and push to stack
						return e.completeAbilityActivation(capturedGameState, capturedPlayer, capturedCard, capturedAbility, abilityIDStr)
					},
					OnCancel: func() error {
						if e.logger != nil {
							e.logger.Info("X value selection cancelled",
								zap.String("card", capturedCard.Name))
						}
						return nil
					},
				}

				if e.logger != nil {
					e.logger.Info("requesting X value selection",
						zap.String("card", card.Name),
						zap.String("type", parsedAbility.XValueType),
						zap.Int("max_x", maxX))
				}

				// Send prompt to client
				e.notifyXValueRequest(gameState.gameID, player.PlayerID, map[string]interface{}{
					"message":   message,
					"available": maxX,
					"min":       0,
					"max":       maxX,
					"source_id": cardID,
				})

				// Return nil - activation is pending, waiting for X value
				return nil
			}
		}

		// Pay mana cost if present
		// NOTE: We work directly with player.ManaPool here instead of using GameContext
		// because we already hold gameState.mu.Lock() and GameContext methods try to acquire
		// locks, causing a deadlock.
		// Per MAGE_ENGINE_ARCHITECTURE.md: Rules-light - pay what's available, don't enforce
		if parsedAbility.ManaCostString != "" {
			manaCost, err := abilities.ParseManaCost(parsedAbility.ManaCostString)
			if err == nil && manaCost != nil && manaCost.Mana != nil {
				// Best-effort payment (rules-light: don't block if insufficient)
				_ = payManaCostDirect(player.ManaPool, manaCost.Mana)
				if e.logger != nil {
					e.logger.Info("attempted mana cost payment for ability",
						zap.String("card", card.Name),
						zap.String("cost", parsedAbility.ManaCostString))
				}
			}
		}

		// Handle sacrifice cost - this moves the card to graveyard as part of the cost
		hasSacrificeCost := parsedAbility.HasSacrificeCost
		if hasSacrificeCost {
			// Move card from battlefield to graveyard
			if err := e.sacrificePermanent(gameState, card); err != nil {
				return fmt.Errorf("failed to sacrifice permanent: %w", err)
			}
			if e.logger != nil {
				e.logger.Info("sacrificed permanent for ability cost",
					zap.String("card", card.Name))
			}
		}

		// Create resolve function for rules-text abilities
		// Captures card and ability info for deferred execution when ability resolves
		capturedCard := card
		capturedPlayer := player
		capturedAbility := parsedAbility
		resolveFunc = func() error {
			if e.logger != nil {
				e.logger.Info("resolving activated ability from rules text",
					zap.String("card", capturedCard.Name),
					zap.String("ability", capturedAbility.FullText))
			}

			// Execute the effect based on parsed ability info
			return e.executeRulesTextAbilityEffect(gameState, capturedCard, capturedPlayer, capturedAbility)
		}
	} else {
		// This is a registered ability with UUID
		cardUUID, err := uuid.Parse(cardID)
		if err != nil {
			return fmt.Errorf("invalid card_id: %w", err)
		}

		abilityUUID, err := uuid.Parse(abilityIDStr)
		if err != nil {
			return fmt.Errorf("invalid ability_id: %w", err)
		}

		// Get the ability from the registry
		if gameState.abilityRegistry == nil {
			return fmt.Errorf("ability registry not initialized")
		}

		ability, err := gameState.abilityRegistry.GetAbility(abilityUUID)
		if err != nil {
			return fmt.Errorf("ability not found: %w", err)
		}

		// Verify ability belongs to this card
		if ability.GetSourceID() != cardUUID {
			return fmt.Errorf("ability does not belong to this permanent")
		}

		// Get activated ability type
		activatedAbility, ok := ability.(*abilities.ActivatedAbility)
		if !ok {
			return fmt.Errorf("ability is not an activated ability")
		}

		// Skip mana abilities - they should use SPECIAL_ACTION/ACTIVATE_MANA_ABILITY
		if activatedAbility.IsManaAbility {
			return fmt.Errorf("use ACTIVATE_MANA_ABILITY for mana abilities")
		}

		abilityText = ability.String()
		isSorcerySpeed = activatedAbility.GetTimingRule() == abilities.TimingSorcery

		// Check for tap cost in registered ability
		for _, cost := range activatedAbility.GetCosts() {
			if _, isTapCost := cost.(*abilities.TapCost); isTapCost {
				hasTapCost = true
				break
			}
		}

		// Pay non-tap costs for registered abilities
		// NOTE: For registered abilities, we need to handle costs directly to avoid deadlock
		// The lock is already held, so GameContext methods would deadlock
		// Per MAGE_ENGINE_ARCHITECTURE.md: Rules-light - pay what's available, don't enforce
		for _, cost := range activatedAbility.GetCosts() {
			if _, isTapCost := cost.(*abilities.TapCost); isTapCost {
				continue // Handle tap cost separately below
			}
			// Handle mana cost directly (best-effort)
			if manaCost, isManaCost := cost.(*abilities.ManaCost); isManaCost && manaCost.Mana != nil {
				_ = payManaCostDirect(player.ManaPool, manaCost.Mana)
			}
			// TODO: Handle other cost types (sacrifice, discard) directly
		}

		// Create resolve function for registered abilities
		// The resolve function will create a GameContext when called during stack resolution
		// At that point, the lock should not be held
		capturedAbility := ability
		resolveFunc = func() error {
			gameUUID, _ := uuid.Parse(gameState.gameID)
			gameCtx := NewGameContext(gameUUID, e, e.logger)
			return capturedAbility.Resolve(context.Background(), gameCtx)
		}
	}

	// Per MAGE_ENGINE_ARCHITECTURE.md: Rules-light - log timing info but don't enforce
	if isSorcerySpeed {
		currentStep := gameState.turnManager.CurrentStep()
		activePlayer := gameState.turnManager.ActivePlayer()

		if currentStep != rules.StepMain1 && currentStep != rules.StepMain2 {
			gameState.addMessage(fmt.Sprintf("⚠️ %s ability activated outside main phase", card.Name), "warning")
		}
		if activePlayer != player.PlayerID {
			gameState.addMessage(fmt.Sprintf("⚠️ %s ability activated during opponent's turn", card.Name), "warning")
		}
	}

	// Check if permanent is tapped and ability requires tap
	if hasTapCost {
		if card.Tapped {
			return fmt.Errorf("permanent is already tapped")
		}
		// Pay the tap cost
		card.Tapped = true
	}

	// Create stack item for the ability
	stackItemID := uuid.New().String()

	stackItem := rules.StackItem{
		ID:          stackItemID,
		Controller:  player.PlayerID,
		Description: fmt.Sprintf("%s: %s", card.Name, abilityText),
		Kind:        rules.StackItemKindActivated,
		SourceID:    cardID,
		Metadata: map[string]string{
			"ability_id": abilityIDStr,
			"card_name":  card.Name,
		},
		Resolve: resolveFunc,
	}

	// Push ability to stack
	gameState.stack.Push(stackItem)
	gameState.trackStackItem()
	gameState.trackStackDepth()
	gameState.trackAction()

	if e.logger != nil {
		e.logger.Info("activated ability added to stack",
			zap.String("player", player.PlayerID),
			zap.String("card", card.Name),
			zap.String("ability", abilityIDStr))
	}

	// Add to game log
	gameState.addMessage(fmt.Sprintf("%s activates ability of %s", player.PlayerID, card.Name), "action")

	// Notify stack update
	e.notifyStackUpdate(gameState.gameID, map[string]interface{}{
		"action":     "ability_activated",
		"player_id":  player.PlayerID,
		"ability_id": abilityIDStr,
		"card_id":    cardID,
		"card_name":  card.Name,
	})

	return nil
}

// completeAbilityActivation finishes ability activation after X value is chosen
// This is called from the X value selection callback
func (e *MageEngine) completeAbilityActivation(
	gameState *engineGameState,
	player *internalPlayer,
	card *internalCard,
	parsedAbility ParsedActivatedAbility,
	abilityIDStr string,
) error {
	if e.logger != nil {
		e.logger.Info("completing ability activation after X selection",
			zap.String("card", card.Name),
			zap.Int("x_value", parsedAbility.chosenXValue))
	}

	// Pay mana cost if present (best-effort, rules-light)
	// Per MAGE_ENGINE_ARCHITECTURE.md: Don't enforce mana costs
	if parsedAbility.ManaCostString != "" {
		manaCost, err := abilities.ParseManaCost(parsedAbility.ManaCostString)
		if err == nil && manaCost != nil && manaCost.Mana != nil {
			_ = payManaCostDirect(player.ManaPool, manaCost.Mana)
			if e.logger != nil {
				e.logger.Info("attempted mana cost payment for ability",
					zap.String("card", card.Name),
					zap.String("cost", parsedAbility.ManaCostString))
			}
		}
	}

	// Handle sacrifice cost
	if parsedAbility.HasSacrificeCost {
		if err := e.sacrificePermanent(gameState, card); err != nil {
			return fmt.Errorf("failed to sacrifice permanent: %w", err)
		}
		if e.logger != nil {
			e.logger.Info("sacrificed permanent for ability cost",
				zap.String("card", card.Name))
		}
	}

	// Handle tap cost
	if parsedAbility.HasTapCost {
		if card.Tapped {
			return fmt.Errorf("permanent is already tapped")
		}
		card.Tapped = true
	}

	// Create resolve function
	capturedCard := card
	capturedPlayer := player
	capturedAbility := parsedAbility
	resolveFunc := func() error {
		if e.logger != nil {
			e.logger.Info("resolving activated ability from rules text",
				zap.String("card", capturedCard.Name),
				zap.String("ability", capturedAbility.FullText),
				zap.Int("x_value", capturedAbility.chosenXValue))
		}
		return e.executeRulesTextAbilityEffect(gameState, capturedCard, capturedPlayer, capturedAbility)
	}

	// Create stack item for the ability
	stackItemID := uuid.New().String()
	stackItem := rules.StackItem{
		ID:          stackItemID,
		Controller:  player.PlayerID,
		Description: fmt.Sprintf("%s: %s (X=%d)", card.Name, parsedAbility.FullText, parsedAbility.chosenXValue),
		Kind:        rules.StackItemKindActivated,
		SourceID:    card.ID,
		Metadata: map[string]string{
			"ability_id": abilityIDStr,
			"card_name":  card.Name,
			"x_value":    fmt.Sprintf("%d", parsedAbility.chosenXValue),
		},
		Resolve: resolveFunc,
	}

	// Push ability to stack
	gameState.stack.Push(stackItem)
	gameState.trackStackItem()
	gameState.trackStackDepth()
	gameState.trackAction()

	if e.logger != nil {
		e.logger.Info("activated ability added to stack after X selection",
			zap.String("player", player.PlayerID),
			zap.String("card", card.Name),
			zap.Int("x_value", parsedAbility.chosenXValue))
	}

	// Add to game log
	gameState.addMessage(fmt.Sprintf("%s activates ability of %s (X=%d)", player.PlayerID, card.Name, parsedAbility.chosenXValue), "action")

	// Notify stack update
	e.notifyStackUpdate(gameState.gameID, map[string]interface{}{
		"action":     "ability_activated",
		"player_id":  player.PlayerID,
		"ability_id": abilityIDStr,
		"card_id":    card.ID,
		"card_name":  card.Name,
		"x_value":    parsedAbility.chosenXValue,
	})

	return nil
}

// executeRulesTextAbilityEffect executes the effect of a rules-text-parsed ability
// This handles common patterns like life gain, damage, etc.
func (e *MageEngine) executeRulesTextAbilityEffect(
	gameState *engineGameState,
	card *internalCard,
	player *internalPlayer,
	ability ParsedActivatedAbility,
) error {
	effectText := strings.ToLower(ability.EffectText)

	if e.logger != nil {
		e.logger.Info("executing rules-text ability effect",
			zap.String("card", card.Name),
			zap.String("effect", ability.EffectText))
	}

	// Handle Martyr of Sands and similar "gain X life" abilities
	if strings.Contains(effectText, "gain") && strings.Contains(effectText, "life") {
		lifeGain := 0

		// Check for "three times X" pattern (Martyr of Sands)
		if strings.Contains(effectText, "three times x") || strings.Contains(effectText, "3x") {
			// Use the chosen X value if ability has X, otherwise fallback to auto-count
			if ability.HasXValue && ability.chosenXValue >= 0 {
				// Use the player's chosen X value
				lifeGain = ability.chosenXValue * 3

				if e.logger != nil {
					e.logger.Info("Martyr of Sands effect: using chosen X value",
						zap.Int("x_value", ability.chosenXValue),
						zap.Int("life_gain", lifeGain))
				}
			} else {
				// Fallback: Count white cards in the player's hand (legacy behavior)
				whiteCardCount := 0
				for _, c := range player.Hand {
					if isCardWhite(c) {
						whiteCardCount++
					}
				}
				lifeGain = whiteCardCount * 3

				if e.logger != nil {
					e.logger.Info("Martyr of Sands effect: auto-counting white cards",
						zap.Int("white_cards", whiteCardCount),
						zap.Int("life_gain", lifeGain))
				}
			}
		} else {
			// Try to parse fixed life gain amount from effect text
			// Pattern: "gain N life" where N is a number
			lifeGain = parseLifeGainAmount(effectText)
		}

		if lifeGain > 0 {
			player.Life += lifeGain
			gameState.addMessage(fmt.Sprintf("%s gains %d life (now %d)", player.PlayerID, lifeGain, player.Life), "life")

			if e.logger != nil {
				e.logger.Info("player gained life",
					zap.String("player", player.PlayerID),
					zap.Int("amount", lifeGain),
					zap.Int("new_life", player.Life))
			}
		}
	}

	// Handle damage effects
	if strings.Contains(effectText, "deal") && strings.Contains(effectText, "damage") {
		// TODO: Parse damage amount and target
		// For now, log that this needs implementation
		if e.logger != nil {
			e.logger.Info("damage effect not yet fully implemented",
				zap.String("effect", ability.EffectText))
		}
	}

	// Handle draw effects
	if strings.Contains(effectText, "draw") && strings.Contains(effectText, "card") {
		// Parse number of cards to draw
		numCards := parseDrawAmount(effectText)
		if numCards > 0 && len(player.Library) >= numCards {
			for i := 0; i < numCards; i++ {
				if len(player.Library) > 0 {
					drawnCard := player.Library[len(player.Library)-1]
					player.Library = player.Library[:len(player.Library)-1]
					player.Hand = append(player.Hand, drawnCard)
					drawnCard.Zone = zoneHand
				}
			}
			gameState.addMessage(fmt.Sprintf("%s draws %d card(s)", player.PlayerID, numCards), "draw")
		}
	}

	// Log the resolution
	gameState.addMessage(fmt.Sprintf("%s's ability resolves: %s", card.Name, ability.EffectText), "ability")

	return nil
}

// isCardWhite checks if a card is white based on its color field
func isCardWhite(card *internalCard) bool {
	if card == nil {
		return false
	}
	colorLower := strings.ToLower(card.Color)
	return strings.Contains(colorLower, "white") || colorLower == "w"
}

// parseLifeGainAmount extracts the life gain amount from effect text
func parseLifeGainAmount(effectText string) int {
	// Common patterns: "gain 3 life", "gain N life"
	effectLower := strings.ToLower(effectText)

	// Number words to digits
	numberWords := map[string]int{
		"one": 1, "two": 2, "three": 3, "four": 4, "five": 5,
		"six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10,
	}

	// Try to find pattern "gain N life"
	for word, num := range numberWords {
		if strings.Contains(effectLower, "gain "+word+" life") {
			return num
		}
	}

	// Try numeric pattern
	for i := 1; i <= 20; i++ {
		if strings.Contains(effectLower, fmt.Sprintf("gain %d life", i)) {
			return i
		}
	}

	return 0
}

// parseDrawAmount extracts the number of cards to draw from effect text
func parseDrawAmount(effectText string) int {
	effectLower := strings.ToLower(effectText)

	// Common patterns: "draw a card", "draw N cards"
	if strings.Contains(effectLower, "draw a card") {
		return 1
	}

	numberWords := map[string]int{
		"two": 2, "three": 3, "four": 4, "five": 5,
	}

	for word, num := range numberWords {
		if strings.Contains(effectLower, "draw "+word+" card") {
			return num
		}
	}

	// Try numeric pattern
	for i := 2; i <= 10; i++ {
		if strings.Contains(effectLower, fmt.Sprintf("draw %d card", i)) {
			return i
		}
	}

	return 0
}

// canPayManaCostDirect checks if a mana pool can pay the given mana cost
// This function does NOT acquire any locks - caller must hold appropriate locks
// Used to avoid deadlock when we already hold gameState.mu.Lock()
func canPayManaCostDirect(pool *mana.ManaPool, manaCost *abilities.Mana) bool {
	if pool == nil || manaCost == nil {
		return manaCost == nil // No cost is always payable
	}

	// Check colored mana requirements
	if pool.White < manaCost.White ||
		pool.Blue < manaCost.Blue ||
		pool.Black < manaCost.Black ||
		pool.Red < manaCost.Red ||
		pool.Green < manaCost.Green {
		return false
	}

	// Calculate available generic mana (after paying colored)
	availableGeneric := (pool.White - manaCost.White) +
		(pool.Blue - manaCost.Blue) +
		(pool.Black - manaCost.Black) +
		(pool.Red - manaCost.Red) +
		(pool.Green - manaCost.Green) +
		pool.Colorless

	return availableGeneric >= manaCost.Colorless
}

// payManaCostDirect pays a mana cost from a mana pool
// This function does NOT acquire any locks - caller must hold appropriate locks
// Used to avoid deadlock when we already hold gameState.mu.Lock()
// payManaCostDirect pays a mana cost from a mana pool (best-effort, rules-light)
// Per MAGE_ENGINE_ARCHITECTURE.md: Don't enforce mana costs, pay what's available
// This function does NOT acquire any locks - caller must hold appropriate locks
func payManaCostDirect(pool *mana.ManaPool, manaCost *abilities.Mana) error {
	if pool == nil || manaCost == nil {
		return nil // No pool or no cost - nothing to do
	}

	// Pay colored mana (best-effort - pay what's available, don't go negative)
	if manaCost.White > 0 && pool.White > 0 {
		pay := manaCost.White
		if pool.White < pay {
			pay = pool.White
		}
		pool.White -= pay
	}
	if manaCost.Blue > 0 && pool.Blue > 0 {
		pay := manaCost.Blue
		if pool.Blue < pay {
			pay = pool.Blue
		}
		pool.Blue -= pay
	}
	if manaCost.Black > 0 && pool.Black > 0 {
		pay := manaCost.Black
		if pool.Black < pay {
			pay = pool.Black
		}
		pool.Black -= pay
	}
	if manaCost.Red > 0 && pool.Red > 0 {
		pay := manaCost.Red
		if pool.Red < pay {
			pay = pool.Red
		}
		pool.Red -= pay
	}
	if manaCost.Green > 0 && pool.Green > 0 {
		pay := manaCost.Green
		if pool.Green < pay {
			pay = pool.Green
		}
		pool.Green -= pay
	}

	// Pay generic/colorless mana from colorless first, then from excess colored
	remaining := manaCost.Colorless
	if remaining > 0 {
		if pool.Colorless >= remaining {
			pool.Colorless -= remaining
			remaining = 0
		} else {
			remaining -= pool.Colorless
			pool.Colorless = 0
		}
	}

	// Pay remaining from any color (prioritize excess)
	colors := []*int{&pool.White, &pool.Blue, &pool.Black, &pool.Red, &pool.Green}
	for _, color := range colors {
		if remaining <= 0 {
			break
		}
		if *color > 0 {
			take := remaining
			if *color < take {
				take = *color
			}
			*color -= take
			remaining -= take
		}
	}
	// Rules-light: Don't error if cost couldn't be fully paid

	return nil
}

// handlePlayLand handles playing a land from hand
// Per MTG Rule 116.2a: Playing a land is a special action (doesn't use stack)
// Per MTG Rule 305.1: Can only play during main phase, with empty stack, once per turn
func (e *MageEngine) handlePlayLand(gameState *engineGameState, player *internalPlayer, cardID string) error {
	// RULES-LIGHT: Priority, timing, and land-per-turn checks removed
	// Players control when lands are played - the UI provides guidance
	_ = gameState.turnManager.CurrentPhase() // Keep for potential UI hints

	// Find card in hand by ID
	var card *internalCard
	var idx int = -1
	for i, c := range player.Hand {
		if c == nil {
			continue
		}
		if c.ID == cardID {
			card = c
			idx = i
			break
		}
	}

	if card == nil {
		return fmt.Errorf("card %s not found in hand", cardID)
	}

	// Verify it's a land
	if !strings.Contains(strings.ToLower(card.Type), "land") {
		return fmt.Errorf("card %s is not a land", card.Name)
	}

	// Remove from hand
	player.Hand = append(player.Hand[:idx], player.Hand[idx+1:]...)

	// Move to battlefield
	card.Zone = zoneBattlefield
	card.ControllerID = player.PlayerID
	card.Tapped = false
	gameState.battlefield = append(gameState.battlefield, card)

	// Increment lands played this turn
	player.LandsPlayedThisTurn++

	// Add message
	gameState.addMessage(fmt.Sprintf("%s plays %s", player.Name, card.Name), "action")

	// Publish land played event
	gameState.eventBus.Publish(rules.NewEvent(rules.EventLandPlayed, card.ID, "", player.PlayerID))

	// Notify state change
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{
		"player": player.PlayerID,
		"action": "play_land",
		"card":   card.Name,
	})

	if e.logger != nil {
		e.logger.Debug("land played",
			zap.String("player", player.PlayerID),
			zap.String("card", card.Name),
			zap.Int("lands_played_this_turn", player.LandsPlayedThisTurn),
		)
	}

	return nil
}

// handleActivateManaAbility activates a mana ability on a permanent
// Per MTG Rule 605: Mana abilities don't use the stack and resolve immediately
func (e *MageEngine) handleActivateManaAbility(gameState *engineGameState, player *internalPlayer, permanentID string) error {
	// Per MAGE_ENGINE_ARCHITECTURE.md: Rules-light - don't enforce priority
	// Mana abilities can technically be activated during mana payment anyway
	// (No warning needed - mana abilities are commonly activated without priority)

	// Find permanent on battlefield
	var permanent *internalCard
	for _, card := range gameState.battlefield {
		if card == nil {
			continue
		}
		if card.ID == permanentID {
			permanent = card
			break
		}
	}

	if permanent == nil {
		return fmt.Errorf("permanent %s not found on battlefield", permanentID)
	}

	// Verify controller
	if permanent.ControllerID != player.PlayerID {
		return fmt.Errorf("you don't control this permanent")
	}

	// Check if already tapped
	if permanent.Tapped {
		return fmt.Errorf("permanent is already tapped")
	}

	// Parse mana production from rules text first
	production := parseManaAbilityFromText(permanent.RulesText)

	// Fall back to basic land subtypes if no rules text mana ability
	if production == nil {
		production = &ManaProduction{}
		for _, st := range permanent.SubTypes {
			switch strings.ToUpper(st) {
			case "PLAINS":
				production.White++
			case "ISLAND":
				production.Blue++
			case "SWAMP":
				production.Black++
			case "MOUNTAIN":
				production.Red++
			case "FOREST":
				production.Green++
			}
		}
	}

	if production.Total() == 0 {
		return fmt.Errorf("permanent %s has no mana ability", permanent.Name)
	}

	// Tap the permanent (pay the cost)
	permanent.Tapped = true

	// Add mana to player's pool
	if production.White > 0 {
		player.ManaPool.Add(mana.ManaWhite, production.White)
	}
	if production.Blue > 0 {
		player.ManaPool.Add(mana.ManaBlue, production.Blue)
	}
	if production.Black > 0 {
		player.ManaPool.Add(mana.ManaBlack, production.Black)
	}
	if production.Red > 0 {
		player.ManaPool.Add(mana.ManaRed, production.Red)
	}
	if production.Green > 0 {
		player.ManaPool.Add(mana.ManaGreen, production.Green)
	}
	if production.Colorless > 0 {
		player.ManaPool.Add(mana.ManaColorless, production.Colorless)
	}

	manaSymbol := production.String()
	gameState.addMessage(fmt.Sprintf("%s taps %s for %s", player.Name, permanent.Name, manaSymbol), "mana")

	// Publish event for triggers
	gameState.eventBus.Publish(rules.Event{
		Type:        rules.EventTapped,
		ID:          uuid.New().String(),
		TargetID:    permanent.ID,
		SourceID:    permanent.ID,
		Controller:  player.PlayerID,
		PlayerID:    player.PlayerID,
		Description: fmt.Sprintf("%s tapped for mana", permanent.Name),
	})

	// Notify state change
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{
		"player":    player.PlayerID,
		"action":    "activate_mana_ability",
		"permanent": permanent.Name,
		"mana":      manaSymbol,
	})

	if e.logger != nil {
		e.logger.Debug("mana ability activated",
			zap.String("player", player.PlayerID),
			zap.String("permanent", permanent.Name),
			zap.String("mana", manaSymbol),
			zap.Int("total", production.Total()),
		)
	}

	return nil
}

// handleAdvancePhase manually advances to the next phase/step
// This is a debug/development feature to allow manual turn progression
func (e *MageEngine) handleAdvancePhase(gameState *engineGameState, player *internalPlayer) error {
	// Per MAGE_ENGINE_ARCHITECTURE.md: Rules-light - don't enforce priority
	if gameState.turnManager.PriorityPlayer() != player.PlayerID {
		gameState.addMessage(fmt.Sprintf("⚠️ %s advances phase without priority", player.Name), "warning")
	}

	// Get next player in turn order for turn transitions
	nextPlayer := gameState.playerOrder[0]
	for i, pid := range gameState.playerOrder {
		if pid == gameState.turnManager.ActivePlayer() {
			nextPlayer = gameState.playerOrder[(i+1)%len(gameState.playerOrder)]
			break
		}
	}

	// Advance to next step
	oldPhase := gameState.turnManager.CurrentPhase()
	oldStep := gameState.turnManager.CurrentStep()
	newPhase, newStep := gameState.turnManager.AdvanceStep(nextPlayer)

	// Reset lands played on new turn
	if newStep == rules.StepUntap {
		e.performUntapStep(gameState, gameState.turnManager.ActivePlayer())
	}

	gameState.addMessage(fmt.Sprintf("Advanced from %s/%s to %s/%s",
		oldPhase.String(), oldStep.String(),
		newPhase.String(), newStep.String()), "phase")

	// Notify phase change
	e.notifyPhaseChange(gameState.gameID, map[string]interface{}{
		"old_phase":       oldPhase.String(),
		"old_step":        oldStep.String(),
		"new_phase":       newPhase.String(),
		"new_step":        newStep.String(),
		"active_player":   gameState.turnManager.ActivePlayer(),
		"priority_player": gameState.turnManager.PriorityPlayer(),
		"turn":            gameState.turnManager.TurnNumber(),
	})

	if e.logger != nil {
		e.logger.Info("phase advanced",
			zap.String("game_id", gameState.gameID),
			zap.String("player", player.PlayerID),
			zap.String("new_phase", newPhase.String()),
			zap.String("new_step", newStep.String()),
		)
	}

	return nil
}

// handleSearchLibrary initiates a library search for a player
// This shows the library contents to the player and waits for them to select a card
func (e *MageEngine) handleSearchLibrary(gameState *engineGameState, player *internalPlayer, destination string, shuffle bool, message string) error {
	// Check if there's already a pending library search
	if gameState.pendingLibrarySearch != nil {
		return fmt.Errorf("there is already a pending library search")
	}

	// Log the search initiation
	gameState.addMessage(fmt.Sprintf("%s searches their library", player.Name), "action")

	// Create the pending library search request
	gameState.pendingLibrarySearch = &PendingLibrarySearchRequest{
		PlayerID:          player.PlayerID,
		SearchingPlayerID: player.PlayerID,
		Message:           message,
		Destination:       destination,
		Shuffle:           shuffle,
		Required:          false, // Player can cancel by selecting nothing
		Timestamp:         time.Now(),
	}

	// Build library card views for the UI
	libraryViews := make([]map[string]interface{}, 0, len(player.Library))
	for _, card := range player.Library {
		if card == nil {
			continue
		}
		libraryViews = append(libraryViews, map[string]interface{}{
			"id":        card.ID,
			"name":      card.Name,
			"type":      card.Type,
			"mana_cost": card.ManaCost,
			"rules":     card.RulesText,
		})
	}

	// Notify the player to select a card from their library
	e.emitNotification(GameNotification{
		Type:     "LIBRARY_SEARCH",
		GameID:   gameState.gameID,
		PlayerID: player.PlayerID,
		Data: map[string]interface{}{
			"message":     message,
			"destination": destination,
			"cards":       libraryViews,
			"can_cancel":  true,
		},
	})

	if e.logger != nil {
		e.logger.Info("library search initiated",
			zap.String("game_id", gameState.gameID),
			zap.String("player", player.PlayerID),
			zap.String("destination", destination),
			zap.Int("library_size", len(player.Library)),
		)
	}

	return nil
}

// handleLibrarySearchSelection processes a card selection from a library search
func (e *MageEngine) handleLibrarySearchSelection(gameState *engineGameState, player *internalPlayer, cardID string) error {
	req := gameState.pendingLibrarySearch
	if req == nil {
		return fmt.Errorf("no pending library search")
	}

	if req.PlayerID != player.PlayerID {
		return fmt.Errorf("library search is for player %s, not %s", req.PlayerID, player.PlayerID)
	}

	// Handle cancellation (empty cardID or "CANCEL")
	if cardID == "" || strings.ToUpper(cardID) == "CANCEL" {
		gameState.pendingLibrarySearch = nil
		gameState.addMessage(fmt.Sprintf("%s finishes searching (found nothing)", player.Name), "action")

		// Still shuffle if required
		if req.Shuffle {
			e.shuffleLibrary(player)
			gameState.addMessage(fmt.Sprintf("%s shuffles their library", player.Name), "action")
		}
		return nil
	}

	// Find the card in the library
	var selectedCard *internalCard
	var cardIndex int = -1
	for i, card := range player.Library {
		if card != nil && card.ID == cardID {
			selectedCard = card
			cardIndex = i
			break
		}
	}

	if selectedCard == nil {
		return fmt.Errorf("card %s not found in library", cardID)
	}

	// Remove card from library
	player.Library = append(player.Library[:cardIndex], player.Library[cardIndex+1:]...)

	// Move card to destination
	switch req.Destination {
	case "hand":
		selectedCard.Zone = zoneHand
		player.Hand = append(player.Hand, selectedCard)
		gameState.addMessage(fmt.Sprintf("%s finds %s and puts it into their hand", player.Name, selectedCard.Name), "action")

	case "battlefield":
		selectedCard.Zone = zoneBattlefield
		selectedCard.ControllerID = player.PlayerID
		gameState.battlefield = append(gameState.battlefield, selectedCard)
		gameState.addMessage(fmt.Sprintf("%s finds %s and puts it onto the battlefield", player.Name, selectedCard.Name), "action")

	case "top":
		selectedCard.Zone = zoneLibrary
		// Put on top of library (end of slice since we draw from end)
		player.Library = append(player.Library, selectedCard)
		gameState.addMessage(fmt.Sprintf("%s finds %s and puts it on top of their library", player.Name, selectedCard.Name), "action")

	case "graveyard":
		selectedCard.Zone = zoneGraveyard
		player.Graveyard = append(player.Graveyard, selectedCard)
		gameState.addMessage(fmt.Sprintf("%s finds %s and puts it into their graveyard", player.Name, selectedCard.Name), "action")

	default:
		// Default to hand
		selectedCard.Zone = zoneHand
		player.Hand = append(player.Hand, selectedCard)
		gameState.addMessage(fmt.Sprintf("%s finds %s and puts it into their hand", player.Name, selectedCard.Name), "action")
	}

	// Shuffle library if required
	if req.Shuffle {
		e.shuffleLibrary(player)
		gameState.addMessage(fmt.Sprintf("%s shuffles their library", player.Name), "action")
	}

	// Clear the pending request
	gameState.pendingLibrarySearch = nil

	// Notify state change
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{
		"player":      player.PlayerID,
		"action":      "library_search_complete",
		"card":        selectedCard.Name,
		"destination": req.Destination,
	})

	if e.logger != nil {
		e.logger.Info("library search completed",
			zap.String("game_id", gameState.gameID),
			zap.String("player", player.PlayerID),
			zap.String("card", selectedCard.Name),
			zap.String("destination", req.Destination),
		)
	}

	return nil
}

// CancelLibrarySearch cancels an ongoing library search
func (e *MageEngine) CancelLibrarySearch(gameID, playerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	if gameState.pendingLibrarySearch == nil {
		return fmt.Errorf("no pending library search")
	}

	if gameState.pendingLibrarySearch.PlayerID != playerID {
		return fmt.Errorf("library search is for player %s, not %s", gameState.pendingLibrarySearch.PlayerID, playerID)
	}

	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// Shuffle if required
	if gameState.pendingLibrarySearch.Shuffle {
		e.shuffleLibrary(player)
		gameState.addMessage(fmt.Sprintf("%s shuffles their library", player.Name), "action")
	}

	gameState.addMessage(fmt.Sprintf("%s finishes searching (found nothing)", player.Name), "action")
	gameState.pendingLibrarySearch = nil

	return nil
}

// resolveSpell resolves a spell on the stack
// Per Java Spell.resolve(): instant/sorcery goes to graveyard, permanents go to battlefield
func (e *MageEngine) resolveSpell(gameState *engineGameState, card *internalCard) error {
	if card == nil {
		return fmt.Errorf("card is nil")
	}

	if e.logger != nil {
		e.logger.Debug("resolving spell",
			zap.String("card_id", card.ID),
			zap.String("card_name", card.Name),
			zap.Int("current_zone", card.Zone),
			zap.String("card_type", card.Type),
		)
	}

	// Determine where the card should go based on its type
	// Per Java: instant/sorcery -> graveyard, permanents (creature, artifact, enchantment, planeswalker, land) -> battlefield
	cardType := strings.ToLower(card.Type)

	// Check if it's a permanent type
	isPermanent := strings.Contains(cardType, "creature") ||
		strings.Contains(cardType, "artifact") ||
		strings.Contains(cardType, "enchantment") ||
		strings.Contains(cardType, "planeswalker") ||
		strings.Contains(cardType, "land")

	if isPermanent {
		// Move to battlefield
		// Per Java: controller.moveCards(card, Zone.BATTLEFIELD, ability, game)
		if err := e.moveCard(gameState, card, zoneBattlefield, card.ControllerID); err != nil {
			return fmt.Errorf("failed to move permanent to Battlefield: %w", err)
		}

		// Apply layer system for power/toughness if it's a creature
		if strings.Contains(cardType, "creature") {
			power, _ := e.parsePowerToughness(card.Power)
			toughness, _ := e.parsePowerToughness(card.Toughness)
			snapshot := effects.NewSnapshot(card.ID, card.ControllerID, []string{"Creature"}, power, toughness, true, true)
			gameState.layerSystem.Apply(snapshot)
			card.Power = fmt.Sprintf("%d", snapshot.Power)
			card.Toughness = fmt.Sprintf("%d", snapshot.Toughness)
		}
	} else {
		// Move instant/sorcery to graveyard
		// Per Java: controller.moveCards(card, Zone.GRAVEYARD, ability, game)
		if err := e.moveCard(gameState, card, zoneGraveyard, ""); err != nil {
			return fmt.Errorf("failed to move spell to graveyard: %w", err)
		}
	}

	// Reset stored bookmark for the controller after spell resolves
	// Per Java PlayerImpl line 1550: resetStoredBookmark(game) after spell resolution
	// This makes the spell resolution irreversible
	if player, exists := gameState.players[card.ControllerID]; exists {
		if player.StoredBookmark != -1 {
			gameState.mu.Unlock() // Temporarily unlock to call ResetPlayerStoredBookmark
			e.ResetPlayerStoredBookmark(gameState.gameID, card.ControllerID)
			gameState.mu.Lock() // Re-acquire lock
		}
	}

	return nil
}

// createTriggeredAbilityForSpell creates a triggered ability when a spell is cast
// This simulates effects like "Sanctuary" that trigger on spell casts
// Per new implementation: adds to triggered queue instead of immediately to stack
func (e *MageEngine) createTriggeredAbilityForSpell(gameState *engineGameState, card *internalCard, casterID string) {
	// For Lightning Bolt, create a triggered ability that gains life
	// This simulates a "Sanctuary" effect for testing
	cardNameLower := strings.ToLower(card.Name)
	if strings.Contains(cardNameLower, "lightning bolt") {
		triggerID := uuid.New().String()

		// Create triggered ability queue item
		triggeredAbility := &triggeredAbilityQueueItem{
			ID:          triggerID,
			SourceID:    card.ID,
			Controller:  casterID,
			Description: fmt.Sprintf("Triggered ability: %s gains 1 life", casterID),
			UsesStack:   true, // This ability uses the stack
			Resolve: func(gs *engineGameState) error {
				player, exists := gs.players[casterID]
				if !exists {
					return fmt.Errorf("player %s not found", casterID)
				}
				oldLife := player.Life
				player.Life += 1
				gs.addMessage(fmt.Sprintf("%s gains 1 life (now %d)", casterID, player.Life), "life")

				// Emit life gain event
				gs.eventBus.Publish(rules.Event{
					Type:        rules.EventGainedLife,
					ID:          uuid.New().String(),
					TargetID:    casterID,
					PlayerID:    casterID,
					Amount:      1,
					Timestamp:   time.Now(),
					Description: fmt.Sprintf("%s gains 1 life (from %d to %d)", casterID, oldLife, player.Life),
				})
				return nil
			},
		}

		// Add to triggered queue instead of directly to stack
		// Per rule 603.3: triggered abilities are put on stack before priority
		gameState.triggeredQueue = append(gameState.triggeredQueue, triggeredAbility)
		gameState.addMessage(fmt.Sprintf("Triggered: %s gains 1 life (queued)", casterID), "action")

		if e.logger != nil {
			e.logger.Debug("queued triggered ability",
				zap.String("trigger_id", triggerID),
				zap.String("spell_id", card.ID),
				zap.String("controller", casterID),
			)
		}
	}
}

// GetGameView returns the current game view for a player
func (e *MageEngine) GetGameView(gameID, playerID string) (interface{}, error) {
	if e.logger != nil {
		e.logger.Info("GetGameView entering", zap.String("game_id", gameID), zap.String("player_id", playerID))
	}

	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("game %s not found", gameID)
	}

	if e.logger != nil {
		e.logger.Info("GetGameView acquiring gameState.mu.RLock")
	}
	gameState.mu.RLock()
	if e.logger != nil {
		e.logger.Info("GetGameView acquired gameState.mu.RLock")
	}
	defer gameState.mu.RUnlock()

	// Log current game state for debugging restoration issues
	if e.logger != nil {
		e.logger.Info("[GAMEVIEW] Current game state",
			zap.String("game_id", gameID),
			zap.String("player_id", playerID),
			zap.String("state", gameState.state.String()),
			zap.Bool("is_mulligan", gameState.state == GameStateMulligan),
			zap.Int("turn", gameState.turnManager.TurnNumber()),
			zap.String("phase", gameState.turnManager.CurrentPhase().String()),
			zap.String("active_player", gameState.turnManager.ActivePlayer()),
			zap.Int("prompt_count", len(gameState.prompts)),
		)

		// Log each player's mulligan status
		for pid, player := range gameState.players {
			e.logger.Info("[GAMEVIEW] Player mulligan status",
				zap.String("game_id", gameID),
				zap.String("player_id", pid),
				zap.Bool("kept_hand", player.KeptHand),
				zap.Int("mulligan_count", player.MulliganCount),
				zap.Int("hand_size", len(player.Hand)),
			)
		}
	}

	// Get player names for display
	activePlayerName := ""
	priorityPlayerName := ""
	if activePlayer, exists := gameState.players[gameState.turnManager.ActivePlayer()]; exists {
		activePlayerName = activePlayer.Name
	}
	if priorityPlayer, exists := gameState.players[gameState.turnManager.PriorityPlayer()]; exists {
		priorityPlayerName = priorityPlayer.Name
	}

	// Get requesting player's land tracking
	landsPlayedThisTurn := 0
	landsAllowedThisTurn := 1
	if requestingPlayer, exists := gameState.players[playerID]; exists {
		landsPlayedThisTurn = requestingPlayer.LandsPlayedThisTurn
		landsAllowedThisTurn = requestingPlayer.LandsPerTurn
	}

	view := &EngineGameView{
		GameID:         gameID,
		State:          gameState.state,
		Phase:          gameState.turnManager.CurrentPhase().String(),
		Step:           gameState.turnManager.CurrentStep().String(),
		Turn:           gameState.turnManager.TurnNumber(),
		ActivePlayerID: gameState.turnManager.ActivePlayer(),
		PriorityPlayer: gameState.turnManager.PriorityPlayer(),
		Players:        e.buildPlayerViewsWithActions(gameState, playerID),
		Battlefield:    e.buildBattlefieldViewsWithActions(gameState, playerID),
		Stack:          e.buildStackViews(gameState),
		Exile:          e.buildCardViews(gameState.exile),
		Command:        e.buildCardViews(gameState.command),
		Revealed:       gameState.revealed,
		LookedAt:       gameState.lookedAt,
		Combat:         e.buildCombatView(gameState),
		StartedAt:      gameState.startedAt,
		Messages:       make([]EngineMessage, len(gameState.messages)),
		Prompts:        make([]EnginePrompt, len(gameState.prompts)),

		// Pre-computed display values
		ActivePlayerName:     activePlayerName,
		PriorityPlayerName:   priorityPlayerName,
		GameFormat:           gameState.gameType,
		IsMulliganPhase:      gameState.state == GameStateMulligan,
		LandsPlayedThisTurn:  landsPlayedThisTurn,
		LandsAllowedThisTurn: landsAllowedThisTurn,
	}

	copy(view.Messages, gameState.messages)
	copy(view.Prompts, gameState.prompts)

	// Add pending library search if it's for this player
	if gameState.pendingLibrarySearch != nil && gameState.pendingLibrarySearch.PlayerID == playerID {
		req := gameState.pendingLibrarySearch
		searchingPlayer := gameState.players[req.SearchingPlayerID]
		if searchingPlayer != nil {
			view.PendingLibrarySearch = &EngineLibrarySearchView{
				PlayerID:    req.PlayerID,
				Message:     req.Message,
				Destination: req.Destination,
				Cards:       e.buildCardViews(searchingPlayer.Library),
				CanCancel:   !req.Required,
			}
		}
	}

	return view, nil
}

// buildPlayerViews builds player views (without available actions)
func (e *MageEngine) buildPlayerViews(gameState *engineGameState, requestingPlayerID string) []EnginePlayerView {
	views := make([]EnginePlayerView, 0, len(gameState.playerOrder))

	for _, playerID := range gameState.playerOrder {
		player := gameState.players[playerID]
		if player == nil {
			if e.logger != nil {
				e.logger.Error("[GAMEVIEW] Player not found in players map (buildPlayerViews)",
					zap.String("game_id", gameState.gameID),
					zap.String("player_id", playerID),
					zap.Strings("player_order", gameState.playerOrder),
				)
			}
			continue
		}
		view := EnginePlayerView{
			PlayerID:     player.PlayerID,
			Name:         player.Name,
			Life:         player.Life,
			Poison:       player.Poison,
			Energy:       player.Energy,
			LibraryCount: len(player.Library),
			HandCount:    len(player.Hand),
			Graveyard:    e.buildCardViews(player.Graveyard),
			ManaPool: EngineManaPoolView{
				White:     player.ManaPool.GetTotal(mana.ManaWhite),
				Blue:      player.ManaPool.GetTotal(mana.ManaBlue),
				Black:     player.ManaPool.GetTotal(mana.ManaBlack),
				Red:       player.ManaPool.GetTotal(mana.ManaRed),
				Green:     player.ManaPool.GetTotal(mana.ManaGreen),
				Colorless: player.ManaPool.GetTotal(mana.ManaColorless),
			},
			// Derive HasPriority from turn manager to ensure consistency
			HasPriority:  gameState.turnManager.PriorityPlayer() == playerID,
			Passed:       player.Passed,
			StateOrdinal: player.StateOrdinal,
			Lost:         player.Lost,
			Left:         player.Left,
			Wins:         player.Wins,
			KeptHand:     player.KeptHand,
		}

		// Only show hand to the owning player
		if playerID == requestingPlayerID {
			view.Hand = e.buildCardViews(player.Hand)
		} else {
			view.Hand = make([]EngineCardView, 0, len(player.Hand))
			for _, card := range player.Hand {
				if card == nil {
					continue
				}
				view.Hand = append(view.Hand, EngineCardView{
					ID:       card.ID,
					FaceDown: true,
					Zone:     zoneHand,
				})
			}
		}

		views = append(views, view)
	}

	return views
}

// buildPlayerViewsWithActions builds player views with available actions for the requesting player's hand
func (e *MageEngine) buildPlayerViewsWithActions(gameState *engineGameState, requestingPlayerID string) []EnginePlayerView {
	views := make([]EnginePlayerView, 0, len(gameState.playerOrder))

	for _, playerID := range gameState.playerOrder {
		player := gameState.players[playerID]
		if player == nil {
			if e.logger != nil {
				e.logger.Error("[GAMEVIEW] Player not found in players map",
					zap.String("game_id", gameState.gameID),
					zap.String("player_id", playerID),
					zap.Strings("player_order", gameState.playerOrder),
				)
			}
			continue
		}
		view := EnginePlayerView{
			PlayerID:     player.PlayerID,
			Name:         player.Name,
			Life:         player.Life,
			Poison:       player.Poison,
			Energy:       player.Energy,
			LibraryCount: len(player.Library),
			HandCount:    len(player.Hand),
			Graveyard:    e.buildCardViews(player.Graveyard),
			ManaPool: EngineManaPoolView{
				White:     player.ManaPool.GetTotal(mana.ManaWhite),
				Blue:      player.ManaPool.GetTotal(mana.ManaBlue),
				Black:     player.ManaPool.GetTotal(mana.ManaBlack),
				Red:       player.ManaPool.GetTotal(mana.ManaRed),
				Green:     player.ManaPool.GetTotal(mana.ManaGreen),
				Colorless: player.ManaPool.GetTotal(mana.ManaColorless),
			},
			// Derive HasPriority from turn manager to ensure consistency
			HasPriority:  gameState.turnManager.PriorityPlayer() == playerID,
			Passed:       player.Passed,
			StateOrdinal: player.StateOrdinal,
			Lost:         player.Lost,
			Left:         player.Left,
			Wins:         player.Wins,
			KeptHand:     player.KeptHand,
		}

		// Only show hand to the owning player
		if playerID == requestingPlayerID {
			// Build card views with available actions
			view.Hand = e.buildCardViewsWithActions(gameState, player.Hand, playerID)

			// Compute HasAvailableActions for the requesting player
			view.HasAvailableActions = e.playerHasAvailableActions(gameState, playerID, view.Hand)
		} else {
			view.Hand = make([]EngineCardView, 0, len(player.Hand))
			for _, card := range player.Hand {
				if card == nil {
					continue
				}
				view.Hand = append(view.Hand, EngineCardView{
					ID:       card.ID,
					FaceDown: true,
					Zone:     zoneHand,
				})
			}
		}

		views = append(views, view)
	}

	return views
}

// playerHasAvailableActions checks if a player has any legal actions they can take right now
func (e *MageEngine) playerHasAvailableActions(gameState *engineGameState, playerID string, handCards []EngineCardView) bool {
	player, exists := gameState.players[playerID]
	if !exists || player.Lost || player.Left {
		return false
	}

	// Player must have priority to take actions
	hasPriority := gameState.turnManager.PriorityPlayer() == playerID
	if !hasPriority {
		return false
	}

	// Check hand cards for enabled actions
	for _, card := range handCards {
		for _, action := range card.AvailableActions {
			if action.IsEnabled {
				return true
			}
		}
	}

	// Check battlefield permanents for activated abilities
	for _, card := range gameState.battlefield {
		if card == nil {
			continue
		}
		if card.ControllerID != playerID {
			continue
		}
		actions := e.getAvailableActionsForPermanent(gameState, card, playerID)
		for _, action := range actions {
			if action.IsEnabled {
				return true
			}
		}
	}

	// TODO: Check for special actions (e.g., suspend, morph face-up, etc.)

	return false
}

// getAvailableActionsForPermanent computes available actions for a permanent on the battlefield
func (e *MageEngine) getAvailableActionsForPermanent(gameState *engineGameState, card *internalCard, playerID string) []EngineCardAction {
	var actions []EngineCardAction

	// Only controller can activate abilities on their permanents
	if card.ControllerID != playerID {
		if e.logger != nil {
			e.logger.Debug("getAvailableActionsForPermanent: not controller",
				zap.String("card", card.Name),
				zap.String("cardController", card.ControllerID),
				zap.String("requestingPlayer", playerID))
		}
		return actions
	}

	hasPriority := gameState.turnManager.PriorityPlayer() == playerID

	if e.logger != nil {
		e.logger.Debug("getAvailableActionsForPermanent: checking mana ability",
			zap.String("card", card.Name),
			zap.String("cardID", card.ID),
			zap.String("rulesText", card.RulesText),
			zap.Bool("hasPriority", hasPriority),
			zap.Strings("subTypes", card.SubTypes))
	}

	// Check for mana abilities - these can be activated any time player has priority
	// or during mana payment (Rule 605.3a)
	if manaAbility := e.getManaAbilityAction(card, hasPriority); manaAbility != nil {
		if e.logger != nil {
			e.logger.Debug("getAvailableActionsForPermanent: mana ability found",
				zap.String("card", card.Name),
				zap.String("displayText", manaAbility.DisplayText),
				zap.Bool("isEnabled", manaAbility.IsEnabled))
		}
		actions = append(actions, *manaAbility)
	} else if e.logger != nil {
		e.logger.Debug("getAvailableActionsForPermanent: no mana ability found",
			zap.String("card", card.Name))
	}

	// Check for activated abilities (non-mana) from the registry
	if gameState.abilityRegistry != nil {
		cardUUID, err := uuid.Parse(card.ID)
		if err == nil {
			registeredAbilities := gameState.abilityRegistry.GetAbilitiesBySource(cardUUID)
			for _, ability := range registeredAbilities {
				// Only process activated abilities (not mana abilities)
				if ability.GetType() == abilities.AbilityTypeActivated {
					action := e.getActivatedAbilityAction(gameState, card, ability, hasPriority)
					if action != nil {
						actions = append(actions, *action)
					}
				}
			}
		}
	}

	// Parse activated abilities from rules text (for cards not yet registered)
	activatedAbilityActions := e.parseActivatedAbilitiesFromText(gameState, card, hasPriority)
	actions = append(actions, activatedAbilityActions...)

	return actions
}

// parseActivatedAbilitiesFromText parses activated abilities from a card's rules text
// and returns EngineCardActions for each detected ability.
// This handles abilities like "{2}, {T}: Draw a card" or "Sacrifice {this}: You gain 3 life"
func (e *MageEngine) parseActivatedAbilitiesFromText(gameState *engineGameState, card *internalCard, hasPriority bool) []EngineCardAction {
	var actions []EngineCardAction

	if card.RulesText == "" {
		return actions
	}

	// Split rules text by ability separator
	abilityLines := strings.Split(card.RulesText, "@@@")

	for i, line := range abilityLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this looks like an activated ability (has a colon with cost before it)
		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 || colonIdx == 0 {
			continue
		}

		costPart := strings.TrimSpace(line[:colonIdx])
		effectPart := strings.TrimSpace(line[colonIdx+1:])

		// Skip if it doesn't look like a cost (should contain mana symbols, {T}, or sacrifice/discard keywords)
		if !looksLikeCost(costPart) {
			continue
		}

		// Skip if this is a mana ability (effect adds mana)
		if isManaAbilityEffect(effectPart) {
			continue
		}

		// Skip triggered abilities that have embedded activated abilities (e.g., Channel)
		if strings.HasPrefix(strings.ToLower(line), "when") ||
			strings.HasPrefix(strings.ToLower(line), "whenever") ||
			strings.HasPrefix(strings.ToLower(line), "at the") {
			continue
		}

		// Skip keyword abilities with reminder text
		if isKeywordReminder(line) {
			continue
		}

		// Skip equip abilities - they need a target creature selection flow
		if strings.HasPrefix(strings.ToLower(costPart), "equip") {
			continue
		}

		// This looks like a non-mana activated ability
		displayText := line
		abilityID := fmt.Sprintf("%s-ability-%d", card.ID, i)

		// Determine if ability can be activated
		isEnabled := hasPriority
		disabledReason := ""

		if !hasPriority {
			isEnabled = false
			disabledReason = "You don't have priority"
		} else {
			// Check tap cost
			if strings.Contains(costPart, "{T}") && card.Tapped {
				isEnabled = false
				disabledReason = "Already tapped"
			}

			// Check for summoning sickness on creatures with tap abilities
			if strings.Contains(costPart, "{T}") && card.SummoningSickness && strings.Contains(strings.ToUpper(card.Type), "CREATURE") {
				isEnabled = false
				disabledReason = "Summoning sickness"
			}

			// Check for "only as a sorcery" restriction
			if strings.Contains(strings.ToLower(effectPart), "only as a sorcery") {
				currentStep := gameState.turnManager.CurrentStep()
				activePlayer := gameState.turnManager.ActivePlayer()
				if currentStep != rules.StepMain1 && currentStep != rules.StepMain2 {
					isEnabled = false
					disabledReason = "Can only activate as a sorcery"
				} else if activePlayer != card.ControllerID {
					isEnabled = false
					disabledReason = "Can only activate during your turn"
				} else if !gameState.stack.IsEmpty() {
					isEnabled = false
					disabledReason = "Stack must be empty"
				}
			}
		}

		if e.logger != nil {
			e.logger.Debug("parseActivatedAbilitiesFromText: found ability",
				zap.String("card", card.Name),
				zap.String("displayText", displayText),
				zap.Bool("isEnabled", isEnabled),
				zap.String("disabledReason", disabledReason))
		}

		actions = append(actions, EngineCardAction{
			ActionType:     "ACTIVATE_ABILITY",
			ActionID:       abilityID,
			DisplayText:    displayText,
			IsEnabled:      isEnabled,
			DisabledReason: disabledReason,
		})
	}

	return actions
}

// looksLikeCost checks if a string looks like an activated ability cost
func looksLikeCost(s string) bool {
	s = strings.ToLower(s)
	// Contains mana symbols
	if strings.Contains(s, "{") && strings.Contains(s, "}") {
		return true
	}
	// Contains tap symbol (case insensitive variations)
	if strings.Contains(s, "{t}") {
		return true
	}
	// Contains sacrifice keyword
	if strings.Contains(s, "sacrifice") {
		return true
	}
	// Contains discard keyword
	if strings.Contains(s, "discard") {
		return true
	}
	// Contains pay life
	if strings.Contains(s, "pay") {
		return true
	}
	// Contains reveal
	if strings.Contains(s, "reveal") {
		return true
	}
	// Contains remove counters
	if strings.Contains(s, "remove") && strings.Contains(s, "counter") {
		return true
	}
	// Contains exile
	if strings.Contains(s, "exile") {
		return true
	}
	return false
}

// isManaAbilityEffect checks if an effect text produces mana (making it a mana ability)
func isManaAbilityEffect(effectPart string) bool {
	effectLower := strings.ToLower(effectPart)
	// Check for "Add {X}" pattern where X is a mana symbol
	if strings.Contains(effectLower, "add {") {
		return true
	}
	// Check for "add one mana" pattern
	if strings.Contains(effectLower, "add one mana") || strings.Contains(effectLower, "add mana") {
		return true
	}
	return false
}

// isKeywordReminder checks if a line is just a keyword with reminder text
func isKeywordReminder(line string) bool {
	keywords := []string{
		"vigilance", "lifelink", "flying", "trample", "haste", "deathtouch",
		"first strike", "double strike", "hexproof", "indestructible", "reach",
		"menace", "flash", "defender", "evolve",
	}
	lineLower := strings.ToLower(line)
	for _, kw := range keywords {
		if strings.HasPrefix(lineLower, kw) {
			return true
		}
	}
	return false
}

// getActivatedAbilityAction creates an EngineCardAction for an activated ability
func (e *MageEngine) getActivatedAbilityAction(gameState *engineGameState, card *internalCard, ability abilities.Ability, hasPriority bool) *EngineCardAction {
	// Check if this is an activated ability
	activatedAbility, ok := ability.(*abilities.ActivatedAbility)
	if !ok {
		return nil
	}

	// Skip mana abilities - they're handled separately
	if activatedAbility.IsManaAbility {
		return nil
	}

	// Build display text from ability
	displayText := ability.String()
	if displayText == "" {
		displayText = "Activate ability"
	}

	// Determine if the ability can be activated
	isEnabled := hasPriority
	disabledReason := ""

	if !hasPriority {
		isEnabled = false
		disabledReason = "You don't have priority"
	} else if card.Tapped {
		// Check if the ability requires tapping (has tap cost)
		for _, cost := range activatedAbility.GetCosts() {
			if _, isTapCost := cost.(*abilities.TapCost); isTapCost {
				isEnabled = false
				disabledReason = "Already tapped"
				break
			}
		}
	}

	// Check timing restrictions
	if isEnabled && activatedAbility.GetTimingRule() == abilities.TimingSorcery {
		// Sorcery-speed abilities can only be activated during main phases
		currentStep := gameState.turnManager.CurrentStep()
		activePlayer := gameState.turnManager.ActivePlayer()
		if currentStep != rules.StepMain1 && currentStep != rules.StepMain2 {
			isEnabled = false
			disabledReason = "Can only activate during main phase"
		} else if activePlayer != card.ControllerID {
			isEnabled = false
			disabledReason = "Can only activate during your turn"
		}
	}

	if e.logger != nil {
		e.logger.Debug("getActivatedAbilityAction: found activated ability",
			zap.String("card", card.Name),
			zap.String("abilityID", ability.GetID().String()),
			zap.String("displayText", displayText),
			zap.Bool("isEnabled", isEnabled),
			zap.String("disabledReason", disabledReason))
	}

	return &EngineCardAction{
		ActionType:     "ACTIVATE_ABILITY",
		ActionID:       ability.GetID().String(),
		DisplayText:    displayText,
		IsEnabled:      isEnabled,
		DisabledReason: disabledReason,
	}
}

// ManaProduction represents the mana produced by a mana ability
type ManaProduction struct {
	White     int
	Blue      int
	Black     int
	Red       int
	Green     int
	Colorless int
}

// Total returns the total amount of mana produced
func (mp *ManaProduction) Total() int {
	return mp.White + mp.Blue + mp.Black + mp.Red + mp.Green + mp.Colorless
}

// String returns a display string like "{C}{C}" or "{G}"
func (mp *ManaProduction) String() string {
	var parts []string
	for i := 0; i < mp.White; i++ {
		parts = append(parts, "{W}")
	}
	for i := 0; i < mp.Blue; i++ {
		parts = append(parts, "{U}")
	}
	for i := 0; i < mp.Black; i++ {
		parts = append(parts, "{B}")
	}
	for i := 0; i < mp.Red; i++ {
		parts = append(parts, "{R}")
	}
	for i := 0; i < mp.Green; i++ {
		parts = append(parts, "{G}")
	}
	for i := 0; i < mp.Colorless; i++ {
		parts = append(parts, "{C}")
	}
	return strings.Join(parts, "")
}

// parseManaAbilityFromText extracts mana production from rules text
// Supports patterns like:
// - "{T}: Add {W}" - single colored mana
// - "{T}: Add {C}" - single colorless
// - "{T}: Add {C}{C}" - multiple mana (Sol Ring)
// - "{T}: Add {G} or {W}" - choice abilities (returns first option for now)
func parseManaAbilityFromText(rulesText string) *ManaProduction {
	if rulesText == "" {
		return nil
	}

	// Normalize the text - handle HTML entities and separators
	text := strings.ReplaceAll(rulesText, "@@@", "\n")
	text = strings.ReplaceAll(text, "&mdash;", "—")

	// Look for tap-for-mana pattern: {T}: Add {X}...
	// The pattern matches lines that start with {T}: Add and captures the mana symbols
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check if this line is a tap-for-mana ability
		if !strings.HasPrefix(line, "{T}: Add ") && !strings.HasPrefix(line, "{T}: Add{") {
			// Also check for "{T}: Add " after other text (like "I -")
			if !strings.Contains(line, "{T}: Add ") {
				continue
			}
			// Extract just the mana ability part
			idx := strings.Index(line, "{T}: Add ")
			if idx >= 0 {
				line = line[idx:]
			}
		}

		// Extract the mana symbols after "Add "
		addIdx := strings.Index(line, "Add ")
		if addIdx < 0 {
			continue
		}

		manaStr := line[addIdx+4:]

		// Parse mana symbols - stop at period, "or", comma, or end of line
		production := &ManaProduction{}
		i := 0
		for i < len(manaStr) {
			if manaStr[i] == '{' {
				endBrace := strings.Index(manaStr[i:], "}")
				if endBrace < 0 {
					break
				}
				symbol := strings.ToUpper(manaStr[i+1 : i+endBrace])
				switch symbol {
				case "W":
					production.White++
				case "U":
					production.Blue++
				case "B":
					production.Black++
				case "R":
					production.Red++
				case "G":
					production.Green++
				case "C":
					production.Colorless++
				}
				i += endBrace + 1
			} else if manaStr[i] == '.' || manaStr[i] == ',' {
				break
			} else if strings.HasPrefix(strings.ToLower(manaStr[i:]), " or ") {
				// Stop at "or" - we only handle the first option for now
				break
			} else {
				i++
			}
		}

		if production.Total() > 0 {
			return production
		}
	}

	return nil
}

// ParsedActivatedAbility represents an activated ability parsed from rules text
type ParsedActivatedAbility struct {
	FullText         string // The complete ability text (cost: effect)
	CostText         string // Just the cost part
	EffectText       string // Just the effect part
	HasTapCost       bool   // Whether the ability requires tapping
	HasSacrificeCost bool   // Whether the ability requires sacrificing this permanent
	ManaCostString   string // The mana cost portion (e.g., "{1}", "{2}{B}")
	IsSorcerySpeed   bool   // Whether the ability can only be activated at sorcery speed
	AbilityIndex     int    // The index of this ability in the card's rules text
	HasXValue        bool   // Whether the ability uses X in cost or effect
	XValueType       string // What X represents: "reveal_white", "mana", "cards", etc.
	chosenXValue     int    // The X value chosen by the player (only valid after selection)
}

// parseActivatedAbilitiesFromText parses activated abilities from a card's rules text
// This is used for ability execution, returning structured data about each ability
func parseActivatedAbilitiesFromText(rulesText string) []ParsedActivatedAbility {
	var abilities []ParsedActivatedAbility

	if rulesText == "" {
		return abilities
	}

	// Split rules text by ability separator
	abilityLines := strings.Split(rulesText, "@@@")

	for i, line := range abilityLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this looks like an activated ability (has a colon with cost before it)
		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 || colonIdx == 0 {
			continue
		}

		costPart := strings.TrimSpace(line[:colonIdx])
		effectPart := strings.TrimSpace(line[colonIdx+1:])

		// Skip if it doesn't look like a cost (should contain mana symbols, {T}, or sacrifice/discard keywords)
		if !looksLikeCost(costPart) {
			continue
		}

		// Skip if this is a mana ability (effect adds mana)
		if isManaAbilityEffect(effectPart) {
			continue
		}

		// Skip triggered abilities that have embedded activated abilities
		lowerLine := strings.ToLower(line)
		if strings.HasPrefix(lowerLine, "when") ||
			strings.HasPrefix(lowerLine, "whenever") ||
			strings.HasPrefix(lowerLine, "at the") {
			continue
		}

		// Skip keyword abilities with reminder text
		if isKeywordReminder(line) {
			continue
		}

		// Skip equip abilities - they need a target creature selection flow
		if strings.HasPrefix(strings.ToLower(costPart), "equip") {
			continue
		}

		// Check for tap cost
		hasTapCost := strings.Contains(costPart, "{T}")

		// Check for sacrifice cost
		costLower := strings.ToLower(costPart)
		hasSacrificeCost := strings.Contains(costLower, "sacrifice {this}") ||
			strings.Contains(costLower, "sacrifice this") ||
			strings.Contains(costLower, ", sacrifice {this}")

		// Extract mana cost (e.g., "{1}", "{2}{B}", etc.)
		manaCost := extractManaCostFromCostText(costPart)

		// Check for sorcery speed restriction
		isSorcerySpeed := strings.Contains(strings.ToLower(effectPart), "only as a sorcery") ||
			strings.Contains(strings.ToLower(effectPart), "activate only as a sorcery")

		// Check for X value in cost or effect
		hasXValue := false
		xValueType := ""
		lowerCost := strings.ToLower(costPart)
		lowerEffect := strings.ToLower(effectPart)

		// Check for X in cost (e.g., "Reveal X white cards")
		if strings.Contains(lowerCost, " x ") || strings.Contains(lowerCost, "reveal x") {
			hasXValue = true
			if strings.Contains(lowerCost, "white card") {
				xValueType = "reveal_white"
			} else if strings.Contains(lowerCost, "card") {
				xValueType = "reveal_cards"
			}
		}
		// Check for X in effect referencing cost (e.g., "three times X life")
		if strings.Contains(lowerEffect, "times x") || strings.Contains(lowerEffect, " x ") {
			hasXValue = true
		}
		// Check for mana X (e.g., "{X}")
		if strings.Contains(costPart, "{X}") {
			hasXValue = true
			xValueType = "mana"
		}

		abilities = append(abilities, ParsedActivatedAbility{
			FullText:         line,
			CostText:         costPart,
			EffectText:       effectPart,
			HasTapCost:       hasTapCost,
			HasSacrificeCost: hasSacrificeCost,
			ManaCostString:   manaCost,
			IsSorcerySpeed:   isSorcerySpeed,
			AbilityIndex:     i,
			HasXValue:        hasXValue,
			XValueType:       xValueType,
		})
	}

	return abilities
}

// extractManaCostFromCostText extracts mana symbols from a cost string
// e.g., "{1}, {T}, Sacrifice {this}" -> "{1}"
// e.g., "{2}{B}, Discard a card" -> "{2}{B}"
func extractManaCostFromCostText(costText string) string {
	var manaParts []string

	// Split by comma to get individual cost components
	parts := strings.Split(costText, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		// Check if this part is a mana cost (contains { and } and has mana symbols)
		if strings.Contains(part, "{") && strings.Contains(part, "}") {
			// Skip tap symbols and other non-mana costs
			if part == "{T}" || part == "{Q}" {
				continue
			}
			// Skip "Sacrifice {this}" type costs
			if strings.Contains(strings.ToLower(part), "sacrifice") {
				continue
			}
			// This looks like a mana cost
			manaParts = append(manaParts, part)
		}
	}

	return strings.Join(manaParts, "")
}

// getManaAbilityAction checks if a permanent has a mana ability and returns the action
func (e *MageEngine) getManaAbilityAction(card *internalCard, hasPriority bool) *EngineCardAction {
	// First try to parse mana ability from rules text
	production := parseManaAbilityFromText(card.RulesText)

	// If no mana ability found in rules text, check for basic land subtypes
	if production == nil {
		production = &ManaProduction{}
		for _, st := range card.SubTypes {
			switch strings.ToUpper(st) {
			case "PLAINS":
				production.White++
			case "ISLAND":
				production.Blue++
			case "SWAMP":
				production.Black++
			case "MOUNTAIN":
				production.Red++
			case "FOREST":
				production.Green++
			}
		}

		// If no mana production from subtypes either, no mana ability
		if production.Total() == 0 {
			return nil
		}
	}

	// Check if the permanent is tapped
	canActivate := true
	reason := ""

	if card.Tapped {
		canActivate = false
		reason = "Already tapped"
	} else if !hasPriority {
		// Mana abilities can technically be activated without priority during payment
		// but for now, require priority for simplicity
		canActivate = false
		reason = "You don't have priority"
	}

	displayText := fmt.Sprintf("{T}: Add %s", production.String())

	return &EngineCardAction{
		ActionType:     "ACTIVATE_MANA_ABILITY",
		ActionID:       card.ID, // Use card ID as the ability identifier
		DisplayText:    displayText,
		IsEnabled:      canActivate,
		DisabledReason: reason,
	}
}

// buildCardViewsWithActions converts internal cards to view cards with available actions
func (e *MageEngine) buildCardViewsWithActions(gameState *engineGameState, cards []*internalCard, playerID string) []EngineCardView {
	views := make([]EngineCardView, 0, len(cards))
	for _, card := range cards {
		if card == nil {
			continue
		}
		views = append(views, EngineCardView{
			ID:                card.ID,
			Name:              card.Name,
			DisplayName:       card.DisplayName,
			ManaCost:          card.ManaCost,
			Type:              card.Type,
			SubTypes:          append([]string(nil), card.SubTypes...),
			SuperTypes:        append([]string(nil), card.SuperTypes...),
			Color:             card.Color,
			Power:             card.Power,
			Toughness:         card.Toughness,
			Loyalty:           card.Loyalty,
			CardNumber:        card.CardNumber,
			ExpansionSet:      card.ExpansionSet,
			Rarity:            card.Rarity,
			RulesText:         card.RulesText,
			Tapped:            card.Tapped,
			Flipped:           card.Flipped,
			Transformed:       card.Transformed,
			FaceDown:          card.FaceDown,
			Zone:              card.Zone,
			ControllerID:      card.ControllerID,
			OwnerID:           card.OwnerID,
			AttachedToCard:    append([]string(nil), card.AttachedToCard...),
			Abilities:         append([]EngineAbilityView(nil), card.Abilities...),
			Counters:          e.buildCounterViews(card.Counters),
			AvailableActions:  e.getAvailableActionsForCard(gameState, card, playerID),
			SummoningSickness: card.SummoningSickness,
		})
	}
	return views
}

// getAvailableActionsForCard computes available actions for a card in hand
func (e *MageEngine) getAvailableActionsForCard(gameState *engineGameState, card *internalCard, playerID string) []EngineCardAction {
	var actions []EngineCardAction

	player, exists := gameState.players[playerID]
	if !exists {
		return actions
	}

	hasPriority := gameState.turnManager.PriorityPlayer() == playerID
	currentPhase := gameState.turnManager.CurrentPhase()
	isMainPhase := currentPhase == rules.PhasePrecombatMain || currentPhase == rules.PhasePostcombatMain
	stackEmpty := gameState.stack.IsEmpty()

	// Check if card is a land
	isLand := strings.Contains(strings.ToLower(card.Type), "land")

	if isLand {
		// Land play action
		canPlay := true
		reason := ""

		if !hasPriority {
			canPlay = false
			reason = "You don't have priority"
		} else if !isMainPhase {
			canPlay = false
			reason = "Only during main phase"
		} else if !stackEmpty {
			canPlay = false
			reason = "Stack must be empty"
		} else if player.LandsPlayedThisTurn >= player.LandsPerTurn {
			canPlay = false
			reason = "Already played a land this turn"
		}

		actions = append(actions, EngineCardAction{
			ActionType:     "PLAY_LAND",
			DisplayText:    "Play Land",
			IsEnabled:      canPlay,
			DisabledReason: reason,
		})
	} else {
		// Spell casting action
		canCast := true
		reason := ""

		if !hasPriority {
			canCast = false
			reason = "You don't have priority"
		}

		// Check if spell requires targets and if valid targets exist
		if canCast {
			targetRequirements := targeting.ParseTargetRequirements(card.Type, card.RulesText)
			for _, req := range targetRequirements {
				// Only check non-optional targets (e.g., "target creature" not "up to X targets")
				if !req.Optional && req.MinTargets > 0 {
					validTargets := e.findValidTargets(gameState, req)
					if len(validTargets) < req.MinTargets {
						canCast = false
						reason = "No valid targets"
						break
					}
				}
			}
		}

		// TODO: Add mana cost checking, timing restrictions (sorcery vs instant), etc.

		actions = append(actions, EngineCardAction{
			ActionType:     "CAST_SPELL",
			DisplayText:    "Cast " + card.Name,
			IsEnabled:      canCast,
			DisabledReason: reason,
		})
	}

	// TODO: Add activated abilities from the card

	return actions
}

// buildBattlefieldViewsWithActions builds battlefield card views with available actions for a player
func (e *MageEngine) buildBattlefieldViewsWithActions(gameState *engineGameState, playerID string) []EngineCardView {
	views := make([]EngineCardView, 0, len(gameState.battlefield))
	for _, card := range gameState.battlefield {
		if card == nil {
			continue
		}
		views = append(views, EngineCardView{
			ID:                card.ID,
			Name:              card.Name,
			DisplayName:       card.DisplayName,
			ManaCost:          card.ManaCost,
			Type:              card.Type,
			SubTypes:          append([]string(nil), card.SubTypes...),
			SuperTypes:        append([]string(nil), card.SuperTypes...),
			Color:             card.Color,
			Power:             card.Power,
			Toughness:         card.Toughness,
			Loyalty:           card.Loyalty,
			CardNumber:        card.CardNumber,
			ExpansionSet:      card.ExpansionSet,
			Rarity:            card.Rarity,
			RulesText:         card.RulesText,
			Tapped:            card.Tapped,
			Flipped:           card.Flipped,
			Transformed:       card.Transformed,
			FaceDown:          card.FaceDown,
			Zone:              card.Zone,
			ControllerID:      card.ControllerID,
			OwnerID:           card.OwnerID,
			AttachedToCard:    append([]string(nil), card.AttachedToCard...),
			Abilities:         append([]EngineAbilityView(nil), card.Abilities...),
			Counters:          e.buildCounterViews(card.Counters),
			AvailableActions:  e.getAvailableActionsForPermanent(gameState, card, playerID),
			SummoningSickness: card.SummoningSickness,
		})
	}
	return views
}

// buildCardViews converts internal cards to view cards
func (e *MageEngine) buildCardViews(cards []*internalCard) []EngineCardView {
	views := make([]EngineCardView, 0, len(cards))
	for _, card := range cards {
		if card == nil {
			continue
		}
		views = append(views, EngineCardView{
			ID:                card.ID,
			Name:              card.Name,
			DisplayName:       card.DisplayName,
			ManaCost:          card.ManaCost,
			Type:              card.Type,
			SubTypes:          append([]string(nil), card.SubTypes...),
			SuperTypes:        append([]string(nil), card.SuperTypes...),
			Color:             card.Color,
			Power:             card.Power,
			Toughness:         card.Toughness,
			Loyalty:           card.Loyalty,
			CardNumber:        card.CardNumber,
			ExpansionSet:      card.ExpansionSet,
			Rarity:            card.Rarity,
			RulesText:         card.RulesText,
			Tapped:            card.Tapped,
			Flipped:           card.Flipped,
			Transformed:       card.Transformed,
			FaceDown:          card.FaceDown,
			Zone:              card.Zone,
			ControllerID:      card.ControllerID,
			OwnerID:           card.OwnerID,
			AttachedToCard:    append([]string(nil), card.AttachedToCard...),
			Abilities:         append([]EngineAbilityView(nil), card.Abilities...),
			Counters:          e.buildCounterViews(card.Counters),
			SummoningSickness: card.SummoningSickness,
		})
	}
	return views
}

// buildStackViews builds stack item views
// Stack.List() returns items bottom-to-top (topmost last), so last item is top of stack
// IMPORTANT: All stack views use item.ID (stack item ID) as their ID, not the source card ID.
// This allows clients to reference stack items for removal via STACK_REMOVE.
func (e *MageEngine) buildStackViews(gameState *engineGameState) []EngineCardView {
	items := gameState.stack.List()
	views := make([]EngineCardView, 0, len(items))

	// Items are already in correct order (bottom to top, topmost last)
	for _, item := range items {
		// Check if this is a triggered ability (not a spell)
		if item.Kind == "TRIGGERED" || item.Kind == rules.StackItemKindTriggered {
			// Create a view for triggered ability using its description
			views = append(views, EngineCardView{
				ID:           item.ID,
				Name:         item.Description,
				DisplayName:  item.Description,
				Zone:         zoneStack,
				ControllerID: item.Controller,
			})
		} else {
			// This is a spell - use the card view but with stack item ID
			card, found := gameState.cards[item.SourceID]
			if !found {
				// Create a placeholder view if card not found
				views = append(views, EngineCardView{
					ID:           item.ID,
					Name:         item.Description,
					DisplayName:  item.Description,
					Zone:         zoneStack,
					ControllerID: item.Controller,
				})
			} else {
				cardView := e.buildCardViews([]*internalCard{card})[0]
				// Use stack item ID, not card ID, so clients can reference for STACK_REMOVE
				cardView.ID = item.ID
				cardView.Zone = zoneStack
				views = append(views, cardView)
			}
		}
	}

	return views
}

// buildCounterViews converts counters to view format
func (e *MageEngine) buildCombatView(gameState *engineGameState) EngineCombatView {
	view := EngineCombatView{
		AttackingPlayerID: gameState.combat.attackingPlayerID,
		Groups:            make([]EngineCombatGroupView, 0, len(gameState.combat.groups)),
	}

	for _, group := range gameState.combat.groups {
		groupView := EngineCombatGroupView{
			Attackers:         make([]string, len(group.attackers)),
			Blockers:          make([]string, len(group.blockers)),
			DefenderID:        group.defenderID,
			DefendingPlayerID: group.defendingPlayerID,
			Blocked:           group.blocked,
		}
		copy(groupView.Attackers, group.attackers)
		copy(groupView.Blockers, group.blockers)
		view.Groups = append(view.Groups, groupView)
	}

	return view
}

func (e *MageEngine) buildCounterViews(counters *counters.Counters) []EngineCounterView {
	if counters == nil {
		return []EngineCounterView{}
	}

	allCounters := counters.GetAll()
	views := make([]EngineCounterView, 0, len(allCounters))
	for name, counter := range allCounters {
		views = append(views, EngineCounterView{
			Name:  name,
			Count: counter.Count,
		})
	}
	return views
}

// GetGameAnalytics returns analytics for a game
func (e *MageEngine) GetGameAnalytics(gameID string) (map[string]interface{}, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	gameState, exists := e.games[gameID]
	if !exists {
		return nil, fmt.Errorf("game %s not found", gameID)
	}

	return gameState.getAnalyticsSummary(), nil
}

// PlayerConcede handles a player conceding the game
// Per Java GameImpl.setConcedingPlayer() and PlayerImpl.concede()
func (e *MageEngine) PlayerConcede(gameID, playerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// Add to conceding players queue if not already there
	alreadyQueued := false
	for _, pid := range gameState.concedingPlayers {
		if pid == playerID {
			alreadyQueued = true
			break
		}
	}
	if !alreadyQueued {
		gameState.concedingPlayers = append(gameState.concedingPlayers, playerID)
	}

	// Mark player as conceded
	player.Conceded = true

	if e.logger != nil {
		e.logger.Info("player conceded",
			zap.String("game_id", gameID),
			zap.String("player_id", playerID),
			zap.String("player_name", player.Name),
		)
	}

	// Process concession immediately (in Java this is done on next priority check)
	e.checkConcede(gameState)
	e.checkIfGameIsOver(gameState)

	return nil
}

// PlayerQuit handles a player quitting the match
// Per Java PlayerImpl.quit()
func (e *MageEngine) PlayerQuit(gameID, playerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	player, exists := gameState.players[playerID]
	if !exists {
		gameState.mu.Unlock()
		return fmt.Errorf("player %s not found", playerID)
	}

	player.Quit = true
	gameState.addMessage(fmt.Sprintf("%s quits the match", player.Name), "system")
	gameState.mu.Unlock()

	if e.logger != nil {
		e.logger.Info("player quit",
			zap.String("game_id", gameID),
			zap.String("player_id", playerID),
			zap.String("player_name", player.Name),
		)
	}

	// Quitting also triggers concession
	return e.PlayerConcede(gameID, playerID)
}

// PlayerTimerTimeout handles a player timing out
// Per Java PlayerImpl.timerTimeout()
func (e *MageEngine) PlayerTimerTimeout(gameID, playerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	player, exists := gameState.players[playerID]
	if !exists {
		gameState.mu.Unlock()
		return fmt.Errorf("player %s not found", playerID)
	}

	player.Quit = true
	player.TimerTimeout = true
	gameState.addMessage(fmt.Sprintf("%s loses due to timer timeout", player.Name), "system")
	gameState.mu.Unlock()

	if e.logger != nil {
		e.logger.Info("player timer timeout",
			zap.String("game_id", gameID),
			zap.String("player_id", playerID),
			zap.String("player_name", player.Name),
		)
	}

	// Timer timeout also triggers concession
	return e.PlayerConcede(gameID, playerID)
}

// PlayerIdleTimeout handles a player idling out
func (e *MageEngine) PlayerIdleTimeout(gameID, playerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	player, exists := gameState.players[playerID]
	if !exists {
		gameState.mu.Unlock()
		return fmt.Errorf("player %s not found", playerID)
	}

	player.Quit = true
	player.IdleTimeout = true
	gameState.addMessage(fmt.Sprintf("%s loses due to idle timeout", player.Name), "system")
	gameState.mu.Unlock()

	if e.logger != nil {
		e.logger.Info("player idle timeout",
			zap.String("game_id", gameID),
			zap.String("player_id", playerID),
			zap.String("player_name", player.Name),
		)
	}

	// Idle timeout also triggers concession
	return e.PlayerConcede(gameID, playerID)
}

// checkConcede processes all players in the conceding queue
// Per Java GameImpl.checkConcede()
func (e *MageEngine) checkConcede(gameState *engineGameState) {
	for len(gameState.concedingPlayers) > 0 {
		// Pop first player from queue
		playerID := gameState.concedingPlayers[0]
		gameState.concedingPlayers = gameState.concedingPlayers[1:]

		// Process their leave
		e.playerLeave(gameState, playerID)
	}
}

// playerLeave handles a player leaving the game
// Per Java PlayerImpl.leave() and GameImpl.leave()
func (e *MageEngine) playerLeave(gameState *engineGameState, playerID string) {
	player, exists := gameState.players[playerID]
	if !exists {
		return
	}

	// Mark player as left and lost
	player.Left = true
	player.Lost = true
	player.Passed = true

	// Emit player lost event
	lostEvent := rules.Event{
		Type:      rules.EventLost,
		ID:        uuid.New().String(),
		PlayerID:  playerID,
		Timestamp: time.Now(),
	}
	gameState.eventBus.Publish(lostEvent)

	gameState.addMessage(fmt.Sprintf("%s has lost the game", player.Name), "system")

	if e.logger != nil {
		e.logger.Info("player left game",
			zap.String("game_id", gameState.gameID),
			zap.String("player_id", playerID),
			zap.String("player_name", player.Name),
		)
	}

	// Per rule 800.4a: When a player leaves the game, all objects owned by that player leave the game
	e.removePlayerObjects(gameState, playerID)
}

// removePlayerObjects removes all objects owned by a player from the game
// Per Java GameImpl.leave() lines 3356-3420
func (e *MageEngine) removePlayerObjects(gameState *engineGameState, playerID string) {
	// Remove permanents from battlefield
	remainingBattlefield := make([]*internalCard, 0)
	for _, card := range gameState.battlefield {
		if card == nil {
			continue
		}
		if card.OwnerID != playerID {
			remainingBattlefield = append(remainingBattlefield, card)
		}
	}
	gameState.battlefield = remainingBattlefield

	// Clear player's zones per rule 800.4a
	if player, exists := gameState.players[playerID]; exists {
		player.Hand = make([]*internalCard, 0)
		player.Library = make([]*internalCard, 0)
		player.Graveyard = make([]*internalCard, 0)
	}

	// Remove from exile
	remainingExile := make([]*internalCard, 0)
	for _, card := range gameState.exile {
		if card == nil {
			continue
		}
		if card.OwnerID != playerID {
			remainingExile = append(remainingExile, card)
		}
	}
	gameState.exile = remainingExile

	// Remove from command zone
	remainingCommand := make([]*internalCard, 0)
	for _, card := range gameState.command {
		if card == nil {
			continue
		}
		if card.OwnerID != playerID {
			remainingCommand = append(remainingCommand, card)
		}
	}
	gameState.command = remainingCommand

	// Remove from stack
	stackItems := gameState.stack.List()
	for _, item := range stackItems {
		if item.Controller == playerID {
			gameState.stack.Remove(item.ID)
		}
	}
}

// checkIfGameIsOver checks if the game should end
// Per Java GameImpl.checkIfGameIsOver()
func (e *MageEngine) checkIfGameIsOver(gameState *engineGameState) bool {
	if gameState.state == GameStateFinished {
		return true
	}

	// Count remaining and losing players
	remainingPlayers := 0
	numLosers := 0
	var lastRemainingPlayer *internalPlayer

	for _, pid := range gameState.playerOrder {
		player := gameState.players[pid]
		if !player.Left {
			remainingPlayers++
			lastRemainingPlayer = player
		}
		if player.Lost {
			numLosers++
		}
	}

	// Game ends if only one player remains or all players have lost
	if remainingPlayers <= 1 || numLosers == len(gameState.playerOrder) {
		if remainingPlayers == 1 && lastRemainingPlayer != nil {
			// Single winner
			lastRemainingPlayer.Wins++
			gameState.state = GameStateFinished
			gameState.addMessage(fmt.Sprintf("%s wins the game!", lastRemainingPlayer.Name), "system")

			// Notify game end
			e.notifyGameStateChange(gameState.gameID, map[string]interface{}{
				"state":     "finished",
				"winner_id": lastRemainingPlayer.PlayerID,
				"winner":    lastRemainingPlayer.Name,
			})

			if e.logger != nil {
				e.logger.Info("game ended",
					zap.String("game_id", gameState.gameID),
					zap.String("winner", lastRemainingPlayer.Name),
				)
			}
		} else {
			// Draw or all players lost
			gameState.state = GameStateFinished
			gameState.addMessage("Game ended in a draw", "system")

			// Notify game end
			e.notifyGameStateChange(gameState.gameID, map[string]interface{}{
				"state":  "finished",
				"result": "draw",
			})

			if e.logger != nil {
				e.logger.Info("game ended in draw",
					zap.String("game_id", gameState.gameID),
				)
			}
		}
		return true
	}

	return false
}

// EndGame ends a game
func (e *MageEngine) EndGame(gameID string, winner string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	gameState, exists := e.games[gameID]
	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	gameState.state = GameStateFinished
	gameState.addMessage(fmt.Sprintf("Game ended. Winner: %s", winner), "action")

	if e.logger != nil {
		e.logger.Info("mage engine ended game",
			zap.String("game_id", gameID),
			zap.String("winner", winner),
		)
	}

	return nil
}

// PauseGame pauses a game
func (e *MageEngine) PauseGame(gameID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Validate state
	if gameState.state == GameStatePaused {
		return fmt.Errorf("game %s is already paused", gameID)
	}
	if gameState.state == GameStateFinished {
		return fmt.Errorf("game %s has ended, cannot pause", gameID)
	}

	gameState.state = GameStatePaused
	gameState.addMessage("Game paused", "action")

	if e.logger != nil {
		e.logger.Info("mage engine paused game", zap.String("game_id", gameID))
	}

	return nil
}

// ResumeGame resumes a paused game
func (e *MageEngine) ResumeGame(gameID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	if gameState.state != GameStatePaused {
		return fmt.Errorf("game %s is not paused", gameID)
	}

	gameState.state = GameStateInProgress
	gameState.addMessage("Game resumed", "action")

	if e.logger != nil {
		e.logger.Info("mage engine resumed game", zap.String("game_id", gameID))
	}

	return nil
}

// checkStateAndTriggered checks state-based actions and processes triggered abilities
// until the game state is stable. This is called before each priority per rule 117.5 and 603.3.
// Per Java implementation: runs SBA → triggers → repeat until stable.
// Returns true if anything happened (SBA or triggers processed).
func (e *MageEngine) checkStateAndTriggered(gameState *engineGameState) bool {
	somethingHappened := false
	maxIterations := 100 // Safety limit to prevent infinite loops

	for i := 0; i < maxIterations; i++ {
		// First check state-based actions
		sbaHappened := e.checkStateBasedActions(gameState)

		// Then process triggered abilities
		triggeredHappened := e.processTriggeredAbilities(gameState)

		// If nothing happened, we're stable
		if !sbaHappened && !triggeredHappened {
			break
		}

		somethingHappened = true

		// If we hit the limit, log a warning
		if i == maxIterations-1 {
			if e.logger != nil {
				e.logger.Warn("checkStateAndTriggered hit iteration limit",
					zap.Int("iterations", maxIterations),
				)
			}
		}
	}

	return somethingHappened
}

// checkStateAndTriggeredAfterResolution checks state-based actions and processes triggered
// abilities after a stack item resolves. This is called after each resolution to ensure
// SBAs and triggers are handled before the next item resolves.
// This is a convenience wrapper around checkStateAndTriggered().
func (e *MageEngine) checkStateAndTriggeredAfterResolution(gameState *engineGameState) {
	e.checkStateAndTriggered(gameState)
}

// processTriggeredAbilities processes triggered abilities from the queue in APNAP order.
// Returns true if any triggered abilities were processed.
// Per rule 603.3: "Once an ability has triggered, its controller puts it on the stack
// as an object that's not a card the next time a player would receive priority."
// Per Java implementation: processes abilities in APNAP order (Active Player, Non-Active Player).
func (e *MageEngine) processTriggeredAbilities(gameState *engineGameState) bool {
	if len(gameState.triggeredQueue) == 0 {
		return false
	}

	played := false
	activePlayerID := gameState.turnManager.ActivePlayer()

	// Process in APNAP order: Active Player first, then Non-Active Players in turn order
	// Per Java GameImpl.checkTriggered() line 2332: for (UUID playerId : state.getPlayerList(state.getActivePlayerId()))
	playerOrder := e.getPlayerListStartingWithActive(gameState, activePlayerID)

	for _, playerID := range playerOrder {
		player := gameState.players[playerID]
		if player == nil {
			continue
		}

		// Process all triggered abilities for this player
		// Per Java: while (player.canRespond()) - player can die or win caused by triggered abilities
		for player.canRespond() {
			// Get triggered abilities for this player
			abilities := e.getTriggeredAbilitiesForPlayer(gameState, playerID)
			if len(abilities) == 0 {
				break
			}

			// Per Java lines 2339-2347: Process non-stack abilities first
			// (e.g., Banisher Priest return exiled creature)
			for i := len(abilities) - 1; i >= 0; i-- {
				ability := abilities[i]
				if !ability.UsesStack {
					// Remove from queue
					e.removeTriggeredAbility(gameState, ability.ID)

					// Execute immediately
					if ability.Resolve != nil {
						if err := ability.Resolve(gameState); err != nil {
							if e.logger != nil {
								e.logger.Error("failed to execute non-stack triggered ability",
									zap.String("ability_id", ability.ID),
									zap.Error(err),
								)
							}
						} else {
							played = true
						}
					}

					// Remove from local list
					abilities = append(abilities[:i], abilities[i+1:]...)
				}
			}

			if len(abilities) == 0 {
				break
			}

			// Per Java lines 2351-2360: If only one ability, put it on stack
			// If multiple, player chooses order (for now, we process in queue order)
			if len(abilities) == 1 {
				ability := abilities[0]
				e.removeTriggeredAbility(gameState, ability.ID)

				// Put on stack
				if err := e.putTriggeredAbilityOnStack(gameState, ability); err != nil {
					if e.logger != nil {
						e.logger.Error("failed to put triggered ability on stack",
							zap.String("ability_id", ability.ID),
							zap.Error(err),
						)
					}
				} else {
					played = true
				}
			} else {
				// Multiple abilities - for now, process in queue order
				// In full implementation, player would choose order
				for _, ability := range abilities {
					e.removeTriggeredAbility(gameState, ability.ID)

					if err := e.putTriggeredAbilityOnStack(gameState, ability); err != nil {
						if e.logger != nil {
							e.logger.Error("failed to put triggered ability on stack",
								zap.String("ability_id", ability.ID),
								zap.Error(err),
							)
						}
					} else {
						played = true
					}
				}
				break
			}
		}
	}

	return played
}

// addSimultaneousEvent adds an event to the simultaneous events queue
// These events will be processed together after stack resolution
func (gameState *engineGameState) addSimultaneousEvent(event rules.Event) {
	gameState.simultaneousEvents = append(gameState.simultaneousEvents, event)
}

// hasSimultaneousEvents returns true if there are events waiting to be processed
func (gameState *engineGameState) hasSimultaneousEvents() bool {
	return len(gameState.simultaneousEvents) > 0
}

// handleSimultaneousEvents processes all simultaneous events
// Per Java GameState.handleSimultaneousEvent(): processes events that happened at the same time
// This allows triggers to see all events that occurred together (e.g., multiple creatures dying)
func (e *MageEngine) handleSimultaneousEvents(gameState *engineGameState) {
	if !gameState.hasSimultaneousEvents() {
		return
	}

	// Copy events to process (new events might be added during processing)
	eventsToHandle := make([]rules.Event, len(gameState.simultaneousEvents))
	copy(eventsToHandle, gameState.simultaneousEvents)
	gameState.simultaneousEvents = nil

	// Process each event through the event bus
	// This allows watchers and triggers to respond to the events
	for _, event := range eventsToHandle {
		gameState.eventBus.Publish(event)
	}

	if e.logger != nil && len(eventsToHandle) > 0 {
		e.logger.Debug("processed simultaneous events",
			zap.Int("count", len(eventsToHandle)),
		)
	}
}

// trackStackDepth updates stack depth metrics
func (gameState *engineGameState) trackStackDepth() {
	if gameState.analytics == nil {
		return
	}

	currentDepth := len(gameState.stack.List())
	if currentDepth > gameState.analytics.maxStackDepth {
		gameState.analytics.maxStackDepth = currentDepth
	}
}

// trackStackItem increments the total stack items counter
func (gameState *engineGameState) trackStackItem() {
	if gameState.analytics != nil {
		gameState.analytics.totalStackItems++
	}
}

// trackPriorityPass increments the priority pass counter
func (gameState *engineGameState) trackPriorityPass() {
	if gameState.analytics != nil {
		gameState.analytics.priorityPassCount++
	}
}

// trackSpellCast increments the spells cast counter
func (gameState *engineGameState) trackSpellCast() {
	if gameState.analytics != nil {
		gameState.analytics.spellsCast++
	}
}

// trackTriggerProcessed increments the triggers processed counter
func (gameState *engineGameState) trackTriggerProcessed() {
	if gameState.analytics != nil {
		gameState.analytics.triggersProcessed++
	}
}

// trackAction increments the action count for the current turn
func (gameState *engineGameState) trackAction() {
	if gameState.analytics == nil {
		return
	}

	currentTurn := gameState.turnManager.TurnNumber()
	gameState.analytics.actionsPerTurn[currentTurn]++
}

// trackTurnStart records the start time of a turn
func (gameState *engineGameState) trackTurnStart() {
	if gameState.analytics == nil {
		return
	}

	currentTurn := gameState.turnManager.TurnNumber()
	gameState.analytics.turnStartTimes[currentTurn] = time.Now()
}

// getAnalyticsSummary returns a summary of game analytics
func (gameState *engineGameState) getAnalyticsSummary() map[string]interface{} {
	if gameState.analytics == nil {
		return nil
	}

	// Calculate average response time per turn
	var totalTurnTime time.Duration
	turnCount := 0
	currentTurn := gameState.turnManager.TurnNumber()

	for turn := 1; turn < currentTurn; turn++ {
		if startTime, exists := gameState.analytics.turnStartTimes[turn]; exists {
			if endTime, exists := gameState.analytics.turnStartTimes[turn+1]; exists {
				totalTurnTime += endTime.Sub(startTime)
				turnCount++
			}
		}
	}

	var avgTurnTime float64
	if turnCount > 0 {
		avgTurnTime = totalTurnTime.Seconds() / float64(turnCount)
	}

	// Calculate total game time
	gameTime := time.Since(gameState.analytics.gameStartTime).Seconds()

	return map[string]interface{}{
		"max_stack_depth":         gameState.analytics.maxStackDepth,
		"total_stack_items":       gameState.analytics.totalStackItems,
		"priority_pass_count":     gameState.analytics.priorityPassCount,
		"spells_cast":             gameState.analytics.spellsCast,
		"abilities_activated":     gameState.analytics.abilitiesActivated,
		"triggers_processed":      gameState.analytics.triggersProcessed,
		"actions_per_turn":        gameState.analytics.actionsPerTurn,
		"avg_turn_time_seconds":   avgTurnTime,
		"total_game_time_seconds": gameTime,
		"current_turn":            currentTurn,
	}
}

// getPlayerListStartingWithActive returns the player list starting with the active player
// and continuing in turn order. This is used for APNAP (Active Player, Non-Active Player) ordering.
func (e *MageEngine) getPlayerListStartingWithActive(gameState *engineGameState, activePlayerID string) []string {
	result := make([]string, 0, len(gameState.playerOrder))

	// Find active player index
	activeIndex := -1
	for i, pid := range gameState.playerOrder {
		if pid == activePlayerID {
			activeIndex = i
			break
		}
	}

	if activeIndex == -1 {
		// Active player not found, return normal order
		return gameState.playerOrder
	}

	// Start with active player, then continue in turn order
	for i := 0; i < len(gameState.playerOrder); i++ {
		idx := (activeIndex + i) % len(gameState.playerOrder)
		result = append(result, gameState.playerOrder[idx])
	}

	return result
}

// getTriggeredAbilitiesForPlayer returns all triggered abilities controlled by the specified player
func (e *MageEngine) getTriggeredAbilitiesForPlayer(gameState *engineGameState, playerID string) []*triggeredAbilityQueueItem {
	result := make([]*triggeredAbilityQueueItem, 0)
	for _, ability := range gameState.triggeredQueue {
		if ability.Controller == playerID {
			result = append(result, ability)
		}
	}
	return result
}

// removeTriggeredAbility removes a triggered ability from the queue
func (e *MageEngine) removeTriggeredAbility(gameState *engineGameState, abilityID string) {
	for i, ability := range gameState.triggeredQueue {
		if ability.ID == abilityID {
			gameState.triggeredQueue = append(gameState.triggeredQueue[:i], gameState.triggeredQueue[i+1:]...)
			return
		}
	}
}

// RegisterCombatTrigger registers a combat trigger for a card
// Per Java: Cards add TriggeredAbilities to their abilities list
func (e *MageEngine) RegisterCombatTrigger(gameID string, trigger *combatTrigger) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	gameState.combatTriggers = append(gameState.combatTriggers, trigger)

	if e.logger != nil {
		e.logger.Debug("registered combat trigger",
			zap.String("game_id", gameID),
			zap.String("source_id", trigger.SourceID),
			zap.String("trigger_type", trigger.TriggerType),
		)
	}

	return nil
}

// checkCombatTriggers checks all registered combat triggers for a given event
// Per Java: TriggeredAbilities.checkTriggers() called when events fire
func (e *MageEngine) checkCombatTriggers(gameState *engineGameState, event rules.Event) {
	for _, trigger := range gameState.combatTriggers {
		// Check if the source card still exists and is on battlefield
		source, exists := gameState.cards[trigger.SourceID]
		if !exists || source.Zone != zoneBattlefield {
			continue
		}

		// Check if the trigger condition is met
		if trigger.Condition != nil && trigger.Condition(gameState, event) {
			// Create and queue the triggered ability
			if trigger.CreateAbility != nil {
				ability := trigger.CreateAbility(gameState, event)
				if ability != nil {
					gameState.triggeredQueue = append(gameState.triggeredQueue, ability)

					if e.logger != nil {
						e.logger.Debug("combat trigger fired",
							zap.String("source_id", trigger.SourceID),
							zap.String("trigger_type", trigger.TriggerType),
							zap.String("ability_id", ability.ID),
						)
					}
				}
			}
		}
	}
}

// putTriggeredAbilityOnStack puts a triggered ability on the stack
func (e *MageEngine) putTriggeredAbilityOnStack(gameState *engineGameState, ability *triggeredAbilityQueueItem) error {
	// Wrap the resolve function to match StackItem signature
	resolveFunc := func() error {
		if ability.Resolve != nil {
			return ability.Resolve(gameState)
		}
		return nil
	}

	// Create stack item for triggered ability
	item := rules.StackItem{
		ID:          ability.ID,
		SourceID:    ability.SourceID,
		Controller:  ability.Controller,
		Description: ability.Description,
		Kind:        "TRIGGERED",
		Resolve:     resolveFunc,
	}

	// Push to stack
	gameState.stack.Push(item)
	gameState.trackStackItem()
	gameState.trackStackDepth()
	gameState.trackTriggerProcessed()

	// Notify trigger
	e.notifyTrigger(gameState.gameID, map[string]interface{}{
		"ability_id":  ability.ID,
		"source_id":   ability.SourceID,
		"controller":  ability.Controller,
		"description": ability.Description,
		"uses_stack":  ability.UsesStack,
	})

	if e.logger != nil {
		e.logger.Debug("put triggered ability on stack",
			zap.String("ability_id", ability.ID),
			zap.String("source_id", ability.SourceID),
			zap.String("controller", ability.Controller),
			zap.String("description", ability.Description),
		)
	}

	return nil
}

// checkStateBasedActions checks and applies state-based actions per rule 117.5
// Returns true if any state-based actions were performed
// Per rule 117.5: "Each time a player would get priority, the game first performs all
// applicable state-based actions as a single event (see rule 704, "State-Based Actions"),
// then repeats this process until no state-based actions are performed."
func (e *MageEngine) checkStateBasedActions(gameState *engineGameState) bool {
	somethingHappened := false

	// Check player loss conditions (704.5a/704.5b/704.5c)
	for _, player := range gameState.players {
		if player.Lost || player.Left {
			continue
		}

		// 704.5a: If a player has 0 or less life, they lose the game
		if player.Life <= 0 {
			player.Lost = true
			gameState.addMessage(fmt.Sprintf("%s loses the game (life <= 0)", player.PlayerID), "action")
			somethingHappened = true
			if e.logger != nil {
				e.logger.Info("player lost due to life",
					zap.String("player_id", player.PlayerID),
					zap.Int("life", player.Life),
				)
			}
			continue
		}

		// 704.5b: If a player has 10 or more poison counters, they lose the game
		if player.Poison >= 10 {
			player.Lost = true
			gameState.addMessage(fmt.Sprintf("%s loses the game (poison >= 10)", player.PlayerID), "action")
			somethingHappened = true
			if e.logger != nil {
				e.logger.Info("player lost due to poison",
					zap.String("player_id", player.PlayerID),
					zap.Int("poison", player.Poison),
				)
			}
			continue
		}

		// Commander rule 903.10a: Check via CommanderBehavior if attached
		// A player that's been dealt 21 or more combat damage by the same commander loses
		if cb := gameState.getCommanderBehavior(); cb != nil && cb.IsCommanderDamageEnabled() && player.CommanderDamage != nil {
			threshold := cb.GetDamageThreshold()
			for commanderID, damage := range player.CommanderDamage {
				if damage >= threshold {
					player.Lost = true
					commanderName := commanderID
					if commander, exists := gameState.cards[commanderID]; exists {
						commanderName = commander.Name
					}
					gameState.addMessage(fmt.Sprintf("%s loses the game (commander damage from %s >= %d)",
						player.PlayerID, commanderName, threshold), "action")
					somethingHappened = true
					if e.logger != nil {
						e.logger.Info("player lost due to commander damage",
							zap.String("player_id", player.PlayerID),
							zap.String("commander_id", commanderID),
							zap.Int("damage", damage),
							zap.Int("threshold", threshold),
						)
					}
					break
				}
			}
			if player.Lost {
				continue
			}
		}

		// 704.5c: If a player would draw a card from an empty library, they lose the game
		// Note: This is typically handled when the draw would occur, but we check here too
		if len(player.Library) == 0 {
			// Only lose if they would draw (this is usually handled during draw step)
			// For now, we'll skip this check as it's typically handled elsewhere
		}
	}

	// Check permanents on battlefield
	creaturesToRemove := make([]*internalCard, 0)
	planeswalkersToRemove := make([]*internalCard, 0)

	for _, card := range gameState.battlefield {
		if card == nil {
			continue
		}
		if card.Zone != zoneBattlefield {
			continue
		}

		// 704.5f: If a creature has toughness 0 or less, it's put into its owner's graveyard
		if strings.Contains(strings.ToLower(card.Type), "creature") {
			toughness, err := e.parsePowerToughness(card.Toughness)
			if err == nil && toughness <= 0 {
				creaturesToRemove = append(creaturesToRemove, card)
				gameState.addMessage(fmt.Sprintf("%s dies (toughness <= 0)", card.Name), "action")
				somethingHappened = true
				if e.logger != nil {
					e.logger.Info("creature dies due to zero toughness",
						zap.String("card_id", card.ID),
						zap.String("card_name", card.Name),
						zap.Int("toughness", toughness),
					)
				}
				continue
			}

			// 704.5g: If a creature has been dealt damage greater than or equal to its toughness,
			// it's destroyed (dies). Note: We need to track damage on creatures for this.
			// For now, we'll skip this as it requires damage tracking infrastructure.
		}

		// 704.5i: If a planeswalker has loyalty 0, it's put into its owner's graveyard
		// Per Rule 306.5c: The loyalty of a planeswalker on the battlefield is equal to the number of loyalty counters on it
		if e.isPlaneswalker(card) {
			loyalty := 0
			if card.Counters != nil {
				loyalty = card.Counters.GetCount("loyalty")
			}
			if loyalty <= 0 {
				planeswalkersToRemove = append(planeswalkersToRemove, card)
				gameState.addMessage(fmt.Sprintf("%s dies (loyalty <= 0)", card.Name), "action")
				somethingHappened = true
				if e.logger != nil {
					e.logger.Info("planeswalker dies due to zero loyalty",
						zap.String("card_id", card.ID),
						zap.String("card_name", card.Name),
						zap.Int("loyalty", loyalty),
					)
				}
			}
		}
	}

	// Remove creatures that died
	for _, card := range creaturesToRemove {
		e.moveCardToGraveyard(gameState, card)
	}

	// Remove planeswalkers that died
	for _, card := range planeswalkersToRemove {
		e.moveCardToGraveyard(gameState, card)
	}

	// Emit events for state-based actions
	if somethingHappened {
		gameState.eventBus.Publish(rules.Event{
			Type:        rules.EventStateBasedActions,
			ID:          uuid.New().String(),
			Timestamp:   time.Now(),
			Description: "State-based actions performed",
		})
	}

	return somethingHappened
}

// parsePowerToughness parses a power/toughness string to an integer
func (e *MageEngine) parsePowerToughness(value string) (int, error) {
	if value == "" {
		return 0, fmt.Errorf("empty value")
	}
	// Remove any non-numeric characters except minus sign
	cleaned := strings.TrimSpace(value)
	var result int
	_, err := fmt.Sscanf(cleaned, "%d", &result)
	return result, err
}

// parseLoyalty parses a loyalty string to an integer
func (e *MageEngine) parseLoyalty(value string) (int, error) {
	return e.parsePowerToughness(value)
}

// moveCardToGraveyard moves a card from battlefield to graveyard
// moveCard moves a card from its current zone to a target zone with proper event emission.
// This is the central function for all zone changes, matching Java's moveCards() behavior.
// Per Java implementation: cards are removed from source zone, added to target zone, and zone change events are emitted.
func (e *MageEngine) moveCard(gameState *engineGameState, card *internalCard, targetZone int, controllerID string) error {
	if card == nil {
		return fmt.Errorf("card is nil")
	}

	sourceZone := card.Zone

	// Remove from source zone
	switch sourceZone {
	case zoneStack:
		// Stack removal is handled by StackManager.Pop(), so we don't need to remove here
		// Just update the zone tracking
	case zoneBattlefield:
		// Remove from battlefield
		for i, bfCard := range gameState.battlefield {
			if bfCard.ID == card.ID {
				gameState.battlefield = append(gameState.battlefield[:i], gameState.battlefield[i+1:]...)
				break
			}
		}
	case zoneHand:
		// Remove from hand
		if player, exists := gameState.players[card.OwnerID]; exists {
			player.Hand = e.removeCardFromSlice(player.Hand, card.ID)
		}
	case zoneGraveyard:
		// Remove from graveyard
		if player, exists := gameState.players[card.OwnerID]; exists {
			player.Graveyard = e.removeCardFromSlice(player.Graveyard, card.ID)
		}
	case zoneExile:
		// Remove from exile
		for i, exCard := range gameState.exile {
			if exCard.ID == card.ID {
				gameState.exile = append(gameState.exile[:i], gameState.exile[i+1:]...)
				break
			}
		}
	case zoneLibrary:
		// Remove from library
		if player, exists := gameState.players[card.OwnerID]; exists {
			player.Library = e.removeCardFromSlice(player.Library, card.ID)
		}
	case zoneCommand:
		// Remove from command zone
		for i, cmdCard := range gameState.command {
			if cmdCard.ID == card.ID {
				gameState.command = append(gameState.command[:i], gameState.command[i+1:]...)
				break
			}
		}
	}

	// Update card zone and controller
	card.Zone = targetZone
	if controllerID != "" {
		card.ControllerID = controllerID
	}

	// Add to target zone
	switch targetZone {
	case zoneBattlefield:
		gameState.battlefield = append(gameState.battlefield, card)

		// Set summoning sickness for creatures entering the battlefield (unless they have haste)
		// Per MTG Rule 302.6: A creature can't attack unless it has been under its controller's
		// control continuously since the start of their most recent turn
		if e.isCreature(card) {
			if !e.hasAbility(card, abilityHaste) {
				card.SummoningSickness = true
			}
		}

		// Emit enters battlefield event
		gameState.eventBus.Publish(rules.Event{
			Type:        rules.EventEntersTheBattlefield,
			ID:          uuid.New().String(),
			TargetID:    card.ID,
			SourceID:    card.ID,
			Controller:  card.ControllerID,
			PlayerID:    card.ControllerID,
			Zone:        zoneBattlefield,
			Timestamp:   time.Now(),
			Description: fmt.Sprintf("%s enters the battlefield", card.Name),
		})
	case zoneGraveyard:
		// Add to owner's graveyard (cards always go to owner's graveyard, not controller's)
		if player, exists := gameState.players[card.OwnerID]; exists {
			player.Graveyard = append(player.Graveyard, card)
		}

		// If moving from battlefield, emit dies event
		if sourceZone == zoneBattlefield {
			gameState.eventBus.Publish(rules.Event{
				Type:        rules.EventPermanentDies,
				ID:          uuid.New().String(),
				TargetID:    card.ID,
				SourceID:    card.ID,
				Controller:  card.ControllerID,
				PlayerID:    card.OwnerID,
				Zone:        zoneGraveyard,
				Timestamp:   time.Now(),
				Description: fmt.Sprintf("%s dies", card.Name),
			})
		}
	case zoneHand:
		if player, exists := gameState.players[card.OwnerID]; exists {
			player.Hand = append(player.Hand, card)
		}
	case zoneExile:
		gameState.exile = append(gameState.exile, card)
	case zoneLibrary:
		if player, exists := gameState.players[card.OwnerID]; exists {
			player.Library = append(player.Library, card)
		}
	case zoneCommand:
		gameState.command = append(gameState.command, card)
	case zoneStack:
		// Stack additions are handled by StackManager.Push(), not here
		return fmt.Errorf("cannot move card to stack via moveCard, use StackManager.Push()")
	}

	// Emit zone change event
	gameState.eventBus.Publish(rules.Event{
		Type:        rules.EventZoneChange,
		ID:          uuid.New().String(),
		TargetID:    card.ID,
		SourceID:    card.ID,
		Controller:  card.ControllerID,
		PlayerID:    card.OwnerID,
		Zone:        targetZone,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("%s moved from zone %d to zone %d", card.Name, sourceZone, targetZone),
		Metadata: map[string]string{
			"source_zone": fmt.Sprintf("%d", sourceZone),
			"target_zone": fmt.Sprintf("%d", targetZone),
		},
	})

	if e.logger != nil {
		e.logger.Debug("moved card",
			zap.String("card_id", card.ID),
			zap.String("card_name", card.Name),
			zap.Int("source_zone", sourceZone),
			zap.Int("target_zone", targetZone),
		)
	}

	return nil
}

func (e *MageEngine) moveCardToGraveyard(gameState *engineGameState, card *internalCard) {
	// Use the unified moveCard function
	if err := e.moveCard(gameState, card, zoneGraveyard, ""); err != nil {
		if e.logger != nil {
			e.logger.Error("failed to move card to graveyard",
				zap.String("card_id", card.ID),
				zap.Error(err),
			)
		}
	}
}

// sacrificePermanent handles sacrificing a permanent (as a cost or effect)
// Per MTG Rules: Sacrificing is moving a permanent to the graveyard
func (e *MageEngine) sacrificePermanent(gameState *engineGameState, card *internalCard) error {
	if card == nil {
		return fmt.Errorf("cannot sacrifice nil card")
	}

	// Make sure the card is on the battlefield
	if card.Zone != zoneBattlefield {
		return fmt.Errorf("can only sacrifice permanents on the battlefield")
	}

	// Move to graveyard (tokens cease to exist instead)
	if err := e.moveCard(gameState, card, zoneGraveyard, ""); err != nil {
		return fmt.Errorf("failed to move sacrificed permanent to graveyard: %w", err)
	}

	// Add game message
	gameState.addMessage(fmt.Sprintf("%s was sacrificed", card.Name), "sacrifice")

	// TODO: Emit "sacrificed" event for death triggers

	return nil
}

// Helper methods for engineGameState

func (s *engineGameState) addMessage(text, color string) {
	// If there's a current action bookmark, use it for rollback
	if s.currentActionBookmark > 0 {
		s.addMessageWithBookmark(text, color, s.currentActionBookmark, true)
	} else {
		s.addMessageWithBookmark(text, color, 0, false)
	}
}

// addMessageWithBookmark adds a message with an associated bookmark for rollback
func (s *engineGameState) addMessageWithBookmark(text, color string, bookmarkID int, rollbackAvailable bool) {
	messageID := s.nextMessageID
	s.nextMessageID++

	msg := EngineMessage{
		Text:              text,
		Color:             color,
		Timestamp:         time.Now(),
		BookmarkID:        bookmarkID,
		RollbackAvailable: rollbackAvailable,
	}
	s.messages = append(s.messages, msg)

	// Track message-to-bookmark mapping if bookmark was provided
	if bookmarkID > 0 {
		if s.messageBookmarks == nil {
			s.messageBookmarks = make(map[int]int)
		}
		s.messageBookmarks[messageID] = bookmarkID
	}

	// Keep only last 1000 messages
	if len(s.messages) > 1000 {
		// Calculate how many messages we're removing
		removeCount := len(s.messages) - 1000
		// Clean up bookmark mappings for removed messages
		for i := 0; i < removeCount; i++ {
			oldMsgID := s.nextMessageID - len(s.messages) + i
			delete(s.messageBookmarks, oldMsgID)
		}
		s.messages = s.messages[removeCount:]
	}
}

func (s *engineGameState) addPrompt(playerID, text string, options []string) {
	s.prompts = append(s.prompts, EnginePrompt{
		PlayerID:  playerID,
		Text:      text,
		Options:   options,
		Timestamp: time.Now(),
	})
}

// AddMessageWithRollback adds a game message with an associated rollback bookmark.
// This creates a snapshot of the game state BEFORE adding the message, allowing
// rollback to the state just before this action occurred.
// NOTE: This method temporarily releases gameState.mu to call BookmarkState,
// so callers must be prepared for state changes.
func (e *MageEngine) AddMessageWithRollback(gameState *engineGameState, text, color string) {
	// Create bookmark of current state (before the message/action)
	// We need to release the lock temporarily since BookmarkState acquires e.mu
	gameState.mu.Unlock()
	bookmarkID, err := e.BookmarkState(gameState.gameID)
	gameState.mu.Lock()

	if err != nil {
		if e.logger != nil {
			e.logger.Warn("failed to create message bookmark",
				zap.String("game_id", gameState.gameID),
				zap.String("message", text),
				zap.Error(err),
			)
		}
		// Fall back to message without bookmark
		gameState.addMessage(text, color)
		return
	}

	// Add message with the bookmark
	gameState.addMessageWithBookmark(text, color, bookmarkID, true)

	if e.logger != nil {
		e.logger.Debug("added message with rollback bookmark",
			zap.String("game_id", gameState.gameID),
			zap.Int("bookmark_id", bookmarkID),
			zap.String("message", text),
		)
	}
}

// hasCommanderBehavior returns true if this game has a CommanderBehavior attached.
// This determines whether commander cards should be moved to the command zone.
func (s *engineGameState) hasCommanderBehavior() bool {
	for _, behavior := range s.behaviors {
		if _, ok := behavior.(*plugin.CommanderBehavior); ok {
			return true
		}
	}
	return false
}

// getCommanderBehavior returns the CommanderBehavior if one is attached, nil otherwise.
func (s *engineGameState) getCommanderBehavior() *plugin.CommanderBehavior {
	for _, behavior := range s.behaviors {
		if cb, ok := behavior.(*plugin.CommanderBehavior); ok {
			return cb
		}
	}
	return nil
}

// buildAttackerPromptOptions builds prompt options for declaring attackers
// Returns options in format: ["ATTACK:creatureID:defenderID", ..., "DONE_ATTACKING"]
func (e *MageEngine) buildAttackerPromptOptions(gameState *engineGameState) []string {
	options := make([]string, 0)
	attackingPlayerID := gameState.combat.attackingPlayerID

	// Find all creatures controlled by attacking player
	for _, card := range gameState.cards {
		if card.Zone != zoneBattlefield {
			continue
		}
		if card.ControllerID != attackingPlayerID {
			continue
		}
		if !e.isCreature(card) {
			continue
		}

		// Skip creatures already declared as attackers
		if card.Attacking {
			continue
		}

		// Check if creature can attack (use internal version since lock is already held)
		if !e.canAttackInternal(gameState, card) {
			continue
		}

		// For each valid defender, add an option (use internal version since lock is already held)
		for defenderID := range gameState.combat.defenders {
			canAttackDefender, _ := e.canAttackDefenderInternal(gameState, card, defenderID)
			if canAttackDefender {
				option := fmt.Sprintf("ATTACK:%s:%s", card.ID, defenderID)
				options = append(options, option)
			}
		}
	}

	// Always add option to finish declaring attackers
	options = append(options, "DONE_ATTACKING")

	return options
}

// buildBlockerPromptOptions builds prompt options for declaring blockers
// Returns options in format: ["BLOCK:blockerID:attackerID", ..., "DONE_BLOCKING"]
func (e *MageEngine) buildBlockerPromptOptions(gameState *engineGameState, defendingPlayerID string) []string {
	options := make([]string, 0)

	// Find all creatures controlled by defending player
	for _, card := range gameState.cards {
		if card.Zone != zoneBattlefield {
			continue
		}
		if card.ControllerID != defendingPlayerID {
			continue
		}
		if !e.isCreature(card) {
			continue
		}

		// Skip creatures already declared as blockers
		if card.Blocking {
			continue
		}

		// For each attacker, check if this creature can block it
		for attackerID := range gameState.combat.attackers {
			canBlock, _ := e.CanBlock(gameState.gameID, card.ID, attackerID)
			if canBlock {
				option := fmt.Sprintf("BLOCK:%s:%s", card.ID, attackerID)
				options = append(options, option)
			}
		}
	}

	// Always add option to finish declaring blockers
	options = append(options, "DONE_BLOCKING")

	return options
}

func (e *MageEngine) getNextPlayer(gameState *engineGameState) string {
	if len(gameState.playerOrder) == 0 {
		return ""
	}
	activeIndex := -1
	for i, pid := range gameState.playerOrder {
		if pid == gameState.turnManager.ActivePlayer() {
			activeIndex = i
			break
		}
	}
	if activeIndex == -1 {
		return gameState.playerOrder[0]
	}
	nextIndex := (activeIndex + 1) % len(gameState.playerOrder)
	return gameState.playerOrder[nextIndex]
}

func (e *MageEngine) getNextPlayerWithPriority(gameState *engineGameState, currentPlayerID string) string {
	currentIndex := -1
	for i, pid := range gameState.playerOrder {
		if pid == currentPlayerID {
			currentIndex = i
			break
		}
	}
	if currentIndex == -1 {
		if len(gameState.playerOrder) > 0 {
			return gameState.playerOrder[0]
		}
		return ""
	}

	// Find next player who hasn't lost or left
	for i := 1; i <= len(gameState.playerOrder); i++ {
		nextIndex := (currentIndex + i) % len(gameState.playerOrder)
		nextPlayerID := gameState.playerOrder[nextIndex]
		player := gameState.players[nextPlayerID]
		if !player.Lost && !player.Left {
			return nextPlayerID
		}
	}

	// All players have lost or left
	return ""
}

// resetPassed resets all players' passed flags, preserving the state for lost/left players.
// Per Java implementation: passed = loses || hasLeft()
// This ensures lost/left players remain passed and don't receive priority.
func (gameState *engineGameState) resetPassed() {
	for _, pid := range gameState.playerOrder {
		p := gameState.players[pid]
		// Set passed = true if player has lost or left, false otherwise
		p.Passed = p.Lost || p.Left
	}
}

// canRespond checks if a player can respond to game actions.
// Per Java implementation: returns true if player is in game (not lost, not left).
// In Java: isInGame() && !abort && !Thread.currentThread().isInterrupted()
// For now, we check Lost and Left; can be extended with Won, Drew, Quit, Abort fields if needed.
func (p *internalPlayer) canRespond() bool {
	return !p.Lost && !p.Left
}

// allPassed checks if all players who can respond have passed.
// Per Java implementation: only considers players who canRespond().
// Returns true if all responding players have passed, false otherwise.
func (gameState *engineGameState) allPassed() bool {
	for _, pid := range gameState.playerOrder {
		p := gameState.players[pid]
		// Only consider players who can respond
		if !p.Passed && p.canRespond() {
			return false
		}
	}
	return true
}

func (e *MageEngine) removeCardFromSlice(cards []*internalCard, cardID string) []*internalCard {
	for i, card := range cards {
		if card == nil {
			continue
		}
		if card.ID == cardID {
			return append(cards[:i], cards[i+1:]...)
		}
	}
	return cards
}

// ChangeControl changes the controller of a permanent on the battlefield
// Returns true if control was successfully changed, false otherwise
// Per Java PermanentImpl.changeControllerId(): emits GAIN_CONTROL and LOSE_CONTROL events
func (e *MageEngine) ChangeControl(gameID, cardID, newControllerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Find the card
	card, found := gameState.cards[cardID]
	if !found {
		return fmt.Errorf("card %s not found", cardID)
	}

	// Verify card is on battlefield
	if card.Zone != zoneBattlefield {
		return fmt.Errorf("card %s is not on battlefield (zone %d)", cardID, card.Zone)
	}

	// Verify new controller exists and is in game
	newController, exists := gameState.players[newControllerID]
	if !exists {
		return fmt.Errorf("player %s not found", newControllerID)
	}
	if newController.Lost || newController.Left {
		return fmt.Errorf("player %s is not in game", newControllerID)
	}

	oldControllerID := card.ControllerID

	// Only emit events if control is actually changing
	if oldControllerID != newControllerID {
		// Emit LOSE_CONTROL event for old controller
		loseControlEvent := rules.Event{
			Type:        rules.EventLoseControl,
			ID:          uuid.New().String(),
			TargetID:    cardID,
			SourceID:    cardID,
			Controller:  oldControllerID,
			PlayerID:    oldControllerID,
			Timestamp:   time.Now(),
			Description: fmt.Sprintf("%s loses control of %s", oldControllerID, card.Name),
			Metadata: map[string]string{
				"old_controller": oldControllerID,
				"new_controller": newControllerID,
			},
		}
		gameState.eventBus.Publish(loseControlEvent)

		// Change the controller
		card.ControllerID = newControllerID

		// Emit GAIN_CONTROL event for new controller
		gainControlEvent := rules.Event{
			Type:        rules.EventGainControl,
			ID:          uuid.New().String(),
			TargetID:    cardID,
			SourceID:    cardID,
			Controller:  newControllerID,
			PlayerID:    newControllerID,
			Timestamp:   time.Now(),
			Description: fmt.Sprintf("%s gains control of %s", newControllerID, card.Name),
			Metadata: map[string]string{
				"old_controller": oldControllerID,
				"new_controller": newControllerID,
			},
		}
		gameState.eventBus.Publish(gainControlEvent)

		gameState.addMessage(fmt.Sprintf("%s gains control of %s", newControllerID, card.Name), "action")

		if e.logger != nil {
			e.logger.Info("control changed",
				zap.String("game_id", gameID),
				zap.String("card_id", cardID),
				zap.String("card_name", card.Name),
				zap.String("old_controller", oldControllerID),
				zap.String("new_controller", newControllerID),
			)
		}
	}

	return nil
}

// createSnapshot creates a deep copy snapshot of the current game state
// This is used for bookmarking and rollback functionality
func (e *MageEngine) createSnapshot(gameState *engineGameState) *gameStateSnapshot {
	snapshot := &gameStateSnapshot{
		GameID:         gameState.gameID,
		GameType:       gameState.gameType,
		State:          gameState.state,
		TurnNumber:     gameState.turnManager.TurnNumber(),
		ActivePlayer:   gameState.turnManager.ActivePlayer(),
		PriorityPlayer: gameState.turnManager.PriorityPlayer(),
		PlayerOrder:    make([]string, len(gameState.playerOrder)),
		Players:        make(map[string]*internalPlayer),
		Cards:          make(map[string]*internalCard),
		Battlefield:    make([]*internalCard, 0, len(gameState.battlefield)),
		Exile:          make([]*internalCard, 0, len(gameState.exile)),
		Command:        make([]*internalCard, 0, len(gameState.command)),
		StackItems:     make([]rules.StackItem, 0),
		Messages:       make([]EngineMessage, len(gameState.messages)),
		Prompts:        make([]EnginePrompt, len(gameState.prompts)),
		Timestamp:      time.Now(),
	}

	// Copy player order
	copy(snapshot.PlayerOrder, gameState.playerOrder)

	// Deep copy players
	for id, player := range gameState.players {
		playerCopy := &internalPlayer{
			PlayerID:            player.PlayerID,
			Name:                player.Name,
			Life:                player.Life,
			Poison:              player.Poison,
			Energy:              player.Energy,
			Library:             make([]*internalCard, len(player.Library)),
			Hand:                make([]*internalCard, len(player.Hand)),
			Graveyard:           make([]*internalCard, len(player.Graveyard)),
			ManaPool:            player.ManaPool.Copy(),
			HasPriority:         player.HasPriority,
			Passed:              player.Passed,
			StateOrdinal:        player.StateOrdinal,
			Lost:                player.Lost,
			Left:                player.Left,
			Wins:                player.Wins,
			Quit:                player.Quit,
			TimerTimeout:        player.TimerTimeout,
			IdleTimeout:         player.IdleTimeout,
			Conceded:            player.Conceded,
			StoredBookmark:      player.StoredBookmark,
			MulliganCount:       player.MulliganCount,
			KeptHand:            player.KeptHand,
			LandsPlayedThisTurn: player.LandsPlayedThisTurn,
			LandsPerTurn:        player.LandsPerTurn,
		}
		snapshot.Players[id] = playerCopy
	}

	// Deep copy all cards
	for id, card := range gameState.cards {
		cardCopy := e.copyCard(card)
		snapshot.Cards[id] = cardCopy

		// Update player zone references
		if player, exists := snapshot.Players[card.OwnerID]; exists {
			ownerPlayer := gameState.players[card.OwnerID]
			if ownerPlayer == nil {
				continue
			}
			switch card.Zone {
			case zoneLibrary:
				for i, c := range ownerPlayer.Library {
					if c != nil && c.ID == card.ID {
						player.Library[i] = cardCopy
						break
					}
				}
			case zoneHand:
				for i, c := range ownerPlayer.Hand {
					if c != nil && c.ID == card.ID {
						player.Hand[i] = cardCopy
						break
					}
				}
			case zoneGraveyard:
				for i, c := range ownerPlayer.Graveyard {
					if c != nil && c.ID == card.ID {
						player.Graveyard[i] = cardCopy
						break
					}
				}
			case zoneBattlefield:
				snapshot.Battlefield = append(snapshot.Battlefield, cardCopy)
			case zoneExile:
				snapshot.Exile = append(snapshot.Exile, cardCopy)
			case zoneCommand:
				snapshot.Command = append(snapshot.Command, cardCopy)
			}
		}
	}

	// Copy stack items
	if gameState.stack != nil {
		snapshot.StackItems = append(snapshot.StackItems, gameState.stack.List()...)
	}

	// Copy messages and prompts
	copy(snapshot.Messages, gameState.messages)
	copy(snapshot.Prompts, gameState.prompts)

	return snapshot
}

// copyCard creates a deep copy of a card
func (e *MageEngine) copyCard(card *internalCard) *internalCard {
	if card == nil {
		return nil
	}

	return &internalCard{
		ID:                card.ID,
		Name:              card.Name,
		DisplayName:       card.DisplayName,
		ManaCost:          card.ManaCost,
		Type:              card.Type,
		SubTypes:          append([]string(nil), card.SubTypes...),
		SuperTypes:        append([]string(nil), card.SuperTypes...),
		Color:             card.Color,
		Power:             card.Power,
		Toughness:         card.Toughness,
		Loyalty:           card.Loyalty,
		CardNumber:        card.CardNumber,
		ExpansionSet:      card.ExpansionSet,
		Rarity:            card.Rarity,
		RulesText:         card.RulesText,
		Tapped:            card.Tapped,
		Flipped:           card.Flipped,
		Transformed:       card.Transformed,
		FaceDown:          card.FaceDown,
		Zone:              card.Zone,
		ControllerID:      card.ControllerID,
		OwnerID:           card.OwnerID,
		AttachedToCard:    append([]string(nil), card.AttachedToCard...),
		Abilities:         append([]EngineAbilityView(nil), card.Abilities...),
		Counters:          card.Counters.Copy(),
		Metadata:          copyMetadata(card.Metadata),
		SummoningSickness: card.SummoningSickness,
		IsToken:           card.IsToken,
		IsCommander:       card.IsCommander,
	}
}

// copyMetadata creates a deep copy of a metadata map
func copyMetadata(src map[string]string) map[string]string {
	if src == nil {
		return make(map[string]string)
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// BookmarkState creates a bookmark of the current game state and returns the bookmark ID
// The bookmark can be used later to restore the game to this state
// Per Java GameImpl.bookmarkState(): saves state and returns index for later restoration
func (e *MageEngine) BookmarkState(gameID string) (int, error) {
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] BookmarkState acquiring e.mu.Lock",
			zap.String("game_id", gameID))
	}
	e.mu.Lock()
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] BookmarkState acquired e.mu.Lock",
			zap.String("game_id", gameID))
	}
	defer func() {
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] BookmarkState releasing e.mu.Lock",
				zap.String("game_id", gameID))
		}
		e.mu.Unlock()
	}()

	gameState, exists := e.games[gameID]
	if !exists {
		return 0, fmt.Errorf("game %s not found", gameID)
	}

	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] BookmarkState acquiring gameState.mu.RLock",
			zap.String("game_id", gameID))
	}
	gameState.mu.RLock()
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] BookmarkState acquired gameState.mu.RLock",
			zap.String("game_id", gameID))
	}
	snapshot := e.createSnapshot(gameState)
	gameState.mu.RUnlock()
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] BookmarkState released gameState.mu.RUnlock",
			zap.String("game_id", gameID))
	}

	// Add snapshot to bookmarks
	if e.bookmarks[gameID] == nil {
		e.bookmarks[gameID] = make([]*gameStateSnapshot, 0)
	}
	e.bookmarks[gameID] = append(e.bookmarks[gameID], snapshot)
	bookmarkID := len(e.bookmarks[gameID])

	// Cleanup old bookmarks if we exceed the max
	// When we trim, we also need to update messageBookmarks to reflect new IDs
	if len(e.bookmarks[gameID]) > e.bookmarksMax {
		trimCount := len(e.bookmarks[gameID]) - e.bookmarksMax
		e.bookmarks[gameID] = e.bookmarks[gameID][trimCount:]

		// Update messageBookmarks to reflect trimmed bookmark IDs
		// Bookmarks were removed from the front, so subtract trimCount from all bookmark IDs
		gameState.mu.Lock()
		for msgID, oldBookmarkID := range gameState.messageBookmarks {
			newBookmarkID := oldBookmarkID - trimCount
			if newBookmarkID > 0 {
				gameState.messageBookmarks[msgID] = newBookmarkID
			} else {
				// This bookmark was trimmed, remove the mapping
				delete(gameState.messageBookmarks, msgID)
			}
		}
		gameState.mu.Unlock()

		if e.logger != nil {
			e.logger.Debug("trimmed old bookmarks",
				zap.String("game_id", gameID),
				zap.Int("trimmed", trimCount),
				zap.Int("remaining", len(e.bookmarks[gameID])),
			)
		}
	}

	if e.logger != nil {
		e.logger.Debug("bookmarked game state",
			zap.String("game_id", gameID),
			zap.Int("bookmark_id", bookmarkID),
			zap.Int("turn", snapshot.TurnNumber),
			zap.Int("total_bookmarks", len(e.bookmarks[gameID])),
		)
	}

	return bookmarkID, nil
}

// RestoreState restores the game to a previously bookmarked state
// Returns error if bookmark doesn't exist or restoration fails
// Per Java GameImpl.restoreState(): rolls back to saved state and removes newer bookmarks
func (e *MageEngine) RestoreState(gameID string, bookmarkID int, context string) error {
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] RestoreState acquiring e.mu.Lock",
			zap.String("game_id", gameID),
			zap.Int("bookmark_id", bookmarkID))
	}
	e.mu.Lock()
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] RestoreState acquired e.mu.Lock",
			zap.String("game_id", gameID))
	}
	defer func() {
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] RestoreState releasing e.mu.Lock",
				zap.String("game_id", gameID))
		}
		e.mu.Unlock()
	}()

	gameState, exists := e.games[gameID]
	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	bookmarks := e.bookmarks[gameID]
	if bookmarks == nil || bookmarkID < 1 || bookmarkID > len(bookmarks) {
		return fmt.Errorf("bookmark %d not found for game %s", bookmarkID, gameID)
	}

	snapshot := bookmarks[bookmarkID-1]

	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] RestoreState acquiring gameState.mu.Lock",
			zap.String("game_id", gameID))
	}
	gameState.mu.Lock()
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] RestoreState acquired gameState.mu.Lock",
			zap.String("game_id", gameID))
	}
	defer func() {
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] RestoreState releasing gameState.mu.Lock",
				zap.String("game_id", gameID))
		}
		gameState.mu.Unlock()
	}()

	// Restore game state from snapshot
	gameState.state = snapshot.State
	gameState.gameType = snapshot.GameType

	// Restore players
	gameState.players = make(map[string]*internalPlayer)
	for id, player := range snapshot.Players {
		gameState.players[id] = player
	}
	gameState.playerOrder = append([]string(nil), snapshot.PlayerOrder...)

	// Restore cards
	gameState.cards = make(map[string]*internalCard)
	for id, card := range snapshot.Cards {
		gameState.cards[id] = card
	}

	// Restore zones
	gameState.battlefield = append([]*internalCard(nil), snapshot.Battlefield...)
	gameState.exile = append([]*internalCard(nil), snapshot.Exile...)
	gameState.command = append([]*internalCard(nil), snapshot.Command...)

	// Restore stack
	gameState.stack = rules.NewStackManager()
	for _, item := range snapshot.StackItems {
		gameState.stack.Push(item)
	}

	// Restore messages and prompts
	gameState.messages = append([]EngineMessage(nil), snapshot.Messages...)
	gameState.prompts = append([]EnginePrompt(nil), snapshot.Prompts...)

	// Remove this bookmark and all newer bookmarks
	e.bookmarks[gameID] = bookmarks[:bookmarkID-1]

	gameState.addMessage(fmt.Sprintf("Game restored to turn %d (%s)", snapshot.TurnNumber, context), "system")

	if e.logger != nil {
		e.logger.Info("restored game state",
			zap.String("game_id", gameID),
			zap.Int("bookmark_id", bookmarkID),
			zap.Int("turn", snapshot.TurnNumber),
			zap.String("context", context),
		)
	}

	return nil
}

// RemoveBookmark removes a bookmark and all newer bookmarks
// Per Java GameImpl.removeBookmark(): cleanup after restoration
func (e *MageEngine) RemoveBookmark(gameID string, bookmarkID int) error {
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] RemoveBookmark acquiring e.mu.Lock",
			zap.String("game_id", gameID),
			zap.Int("bookmark_id", bookmarkID))
	}
	e.mu.Lock()
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] RemoveBookmark acquired e.mu.Lock",
			zap.String("game_id", gameID))
	}
	defer func() {
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] RemoveBookmark releasing e.mu.Lock",
				zap.String("game_id", gameID))
		}
		e.mu.Unlock()
	}()

	bookmarks := e.bookmarks[gameID]
	if bookmarks == nil || bookmarkID < 1 || bookmarkID > len(bookmarks) {
		return fmt.Errorf("bookmark %d not found for game %s", bookmarkID, gameID)
	}

	// Remove this bookmark and all newer ones
	e.bookmarks[gameID] = bookmarks[:bookmarkID-1]

	if e.logger != nil {
		e.logger.Debug("removed bookmark",
			zap.String("game_id", gameID),
			zap.Int("bookmark_id", bookmarkID),
		)
	}

	return nil
}

// ClearBookmarks removes all bookmarks for a game
// Used when game ends or for cleanup
func (e *MageEngine) ClearBookmarks(gameID string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.bookmarks, gameID)

	if e.logger != nil {
		e.logger.Debug("cleared all bookmarks",
			zap.String("game_id", gameID),
		)
	}
}

// SetPlayerStoredBookmark sets a player's stored bookmark for undo
// Per Java PlayerImpl.setStoredBookmark(): enables undo button for player
func (e *MageEngine) SetPlayerStoredBookmark(gameID, playerID string, bookmarkID int) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	player.StoredBookmark = bookmarkID

	if e.logger != nil {
		e.logger.Debug("set player stored bookmark",
			zap.String("game_id", gameID),
			zap.String("player_id", playerID),
			zap.Int("bookmark_id", bookmarkID),
		)
	}

	return nil
}

// ResetPlayerStoredBookmark clears a player's stored bookmark and removes it from the bookmark list
// Per Java PlayerImpl.resetStoredBookmark(): disables undo button for player
func (e *MageEngine) ResetPlayerStoredBookmark(gameID, playerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	player, exists := gameState.players[playerID]
	if !exists {
		gameState.mu.Unlock()
		return fmt.Errorf("player %s not found", playerID)
	}

	bookmarkID := player.StoredBookmark
	player.StoredBookmark = -1
	gameState.mu.Unlock()

	// Remove the bookmark if it exists
	if bookmarkID != -1 {
		e.RemoveBookmark(gameID, bookmarkID)
	}

	if e.logger != nil {
		e.logger.Debug("reset player stored bookmark",
			zap.String("game_id", gameID),
			zap.String("player_id", playerID),
			zap.Int("old_bookmark_id", bookmarkID),
		)
	}

	return nil
}

// Undo performs a player-initiated undo operation
// Per Java GameImpl.undo(): restores to player's stored bookmark if available
func (e *MageEngine) Undo(gameID, playerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	player, exists := gameState.players[playerID]
	if !exists {
		gameState.mu.Unlock()
		return fmt.Errorf("player %s not found", playerID)
	}

	bookmarkID := player.StoredBookmark
	gameState.mu.Unlock()

	if bookmarkID == -1 {
		return fmt.Errorf("no undo available for player %s", playerID)
	}

	// Restore to the stored bookmark
	if err := e.RestoreState(gameID, bookmarkID, fmt.Sprintf("player %s undo", playerID)); err != nil {
		return fmt.Errorf("failed to undo: %w", err)
	}

	// Clear the stored bookmark
	if err := e.SetPlayerStoredBookmark(gameID, playerID, -1); err != nil {
		return fmt.Errorf("failed to clear stored bookmark: %w", err)
	}

	if e.logger != nil {
		e.logger.Info("player undo",
			zap.String("game_id", gameID),
			zap.String("player_id", playerID),
			zap.Int("bookmark_id", bookmarkID),
		)
	}

	// Notify players of the undo
	e.notifyGameStateChange(gameID, map[string]interface{}{
		"type":      "undo",
		"player_id": playerID,
	})

	return nil
}

// RequestRollback initiates a rollback request to a specific message.
// This requires opponent consent in multiplayer games.
// Returns the request ID that can be used to track the response.
func (e *MageEngine) RequestRollback(gameID, playerID string, messageID int) (string, error) {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Check if there's already a pending rollback request
	if gameState.pendingRollbackRequest != nil {
		return "", fmt.Errorf("a rollback request is already pending")
	}

	// Find the message and its bookmark
	if messageID < 1 || messageID > len(gameState.messages) {
		return "", fmt.Errorf("message %d not found", messageID)
	}

	message := gameState.messages[messageID-1]
	if !message.RollbackAvailable {
		return "", fmt.Errorf("rollback not available for message %d", messageID)
	}

	bookmarkID := message.BookmarkID
	if bookmarkID <= 0 {
		return "", fmt.Errorf("no bookmark available for message %d", messageID)
	}

	// Generate a unique request ID
	requestID := uuid.New().String()

	// Get requesting player's name
	requestingPlayerName := playerID
	if player, exists := gameState.players[playerID]; exists {
		requestingPlayerName = player.Name
	}

	// Create the pending request
	gameState.pendingRollbackRequest = &PendingRollbackRequest{
		RequestID:         requestID,
		RequestingPlayer:  playerID,
		TargetMessageID:   messageID,
		TargetBookmarkID:  bookmarkID,
		TargetMessageText: message.Text,
		Timestamp:         time.Now(),
	}

	if e.logger != nil {
		e.logger.Info("rollback request created",
			zap.String("game_id", gameID),
			zap.String("request_id", requestID),
			zap.String("player_id", playerID),
			zap.Int("message_id", messageID),
			zap.Int("bookmark_id", bookmarkID),
		)
	}

	// Notify all other players about the rollback request
	e.emitNotification(GameNotification{
		Type:      "ROLLBACK_REQUEST",
		GameID:    gameID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"request_id":             requestID,
			"requesting_player_id":   playerID,
			"requesting_player_name": requestingPlayerName,
			"target_message_id":      messageID,
			"target_message_text":    message.Text,
		},
	})

	return requestID, nil
}

// RespondToRollback handles a player's response to a rollback request.
// If approved by all opponents, the rollback is performed.
func (e *MageEngine) RespondToRollback(gameID, playerID string, requestID string, approved bool) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()

	// Verify the pending request exists and matches
	if gameState.pendingRollbackRequest == nil {
		gameState.mu.Unlock()
		return fmt.Errorf("no pending rollback request")
	}

	if gameState.pendingRollbackRequest.RequestID != requestID {
		gameState.mu.Unlock()
		return fmt.Errorf("rollback request ID mismatch")
	}

	// Don't allow the requesting player to respond to their own request
	if gameState.pendingRollbackRequest.RequestingPlayer == playerID {
		gameState.mu.Unlock()
		return fmt.Errorf("cannot respond to your own rollback request")
	}

	request := gameState.pendingRollbackRequest
	respondingPlayerName := playerID
	if player, exists := gameState.players[playerID]; exists {
		respondingPlayerName = player.Name
	}

	// Clear the pending request
	gameState.pendingRollbackRequest = nil
	gameState.mu.Unlock()

	if e.logger != nil {
		e.logger.Info("rollback response received",
			zap.String("game_id", gameID),
			zap.String("request_id", requestID),
			zap.String("responding_player", playerID),
			zap.Bool("approved", approved),
		)
	}

	if !approved {
		// Notify all players that the rollback was denied
		e.emitNotification(GameNotification{
			Type:      "ROLLBACK_DENIED",
			GameID:    gameID,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"request_id":             requestID,
				"responding_player_id":   playerID,
				"responding_player_name": respondingPlayerName,
			},
		})
		return nil
	}

	// Rollback approved - perform the rollback
	return e.RollbackToMessage(gameID, request.TargetMessageID, request.RequestingPlayer)
}

// RollbackToMessage performs a rollback to the state before a specific message was added.
// This uses the bookmark associated with the message.
func (e *MageEngine) RollbackToMessage(gameID string, messageID int, initiatedBy string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()

	// Find the message and its bookmark
	if messageID < 1 || messageID > len(gameState.messages) {
		gameState.mu.Unlock()
		return fmt.Errorf("message %d not found", messageID)
	}

	message := gameState.messages[messageID-1]
	bookmarkID := message.BookmarkID

	if bookmarkID <= 0 {
		gameState.mu.Unlock()
		return fmt.Errorf("no bookmark available for message %d", messageID)
	}

	initiatedByName := initiatedBy
	if player, exists := gameState.players[initiatedBy]; exists {
		initiatedByName = player.Name
	}

	gameState.mu.Unlock()

	// Perform the rollback using the existing RestoreState mechanism
	context := fmt.Sprintf("rollback to message %d by %s", messageID, initiatedByName)
	if err := e.RestoreState(gameID, bookmarkID, context); err != nil {
		return fmt.Errorf("failed to rollback to message %d: %w", messageID, err)
	}

	if e.logger != nil {
		e.logger.Info("rollback to message completed",
			zap.String("game_id", gameID),
			zap.Int("message_id", messageID),
			zap.Int("bookmark_id", bookmarkID),
			zap.String("initiated_by", initiatedBy),
		)
	}

	// Notify all players about the completed rollback
	e.emitNotification(GameNotification{
		Type:      "ROLLBACK_COMPLETE",
		GameID:    gameID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"target_message_id": messageID,
			"initiated_by":      initiatedBy,
			"initiated_by_name": initiatedByName,
		},
	})

	return nil
}

// CancelRollbackRequest cancels any pending rollback request for a game
func (e *MageEngine) CancelRollbackRequest(gameID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	if gameState.pendingRollbackRequest == nil {
		return nil // No pending request to cancel
	}

	requestID := gameState.pendingRollbackRequest.RequestID
	gameState.pendingRollbackRequest = nil

	if e.logger != nil {
		e.logger.Info("rollback request cancelled",
			zap.String("game_id", gameID),
			zap.String("request_id", requestID),
		)
	}

	return nil
}

// GetPendingRollbackRequest returns the pending rollback request for a game, if any
func (e *MageEngine) GetPendingRollbackRequest(gameID string) (*PendingRollbackRequest, error) {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	return gameState.pendingRollbackRequest, nil
}

// SaveTurnSnapshot saves a snapshot at the start of a turn for turn rollback
// Per Java GameImpl.saveRollBackGameState(): keeps last N turns for rollback
func (e *MageEngine) SaveTurnSnapshot(gameID string, turnNumber int) error {
	if !e.rollbackAllowed {
		return nil // Turn rollback disabled
	}

	e.mu.Lock()
	gameState, exists := e.games[gameID]
	if !exists {
		e.mu.Unlock()
		return fmt.Errorf("game %s not found", gameID)
	}
	persistenceRepo := e.persistenceRepo
	e.mu.Unlock()

	gameState.mu.RLock()
	snapshot := e.createSnapshot(gameState)
	tableID := gameState.gameID // Table ID is stored as gameID for now (can be enhanced)
	gameType := gameState.gameType
	players := make([]string, len(gameState.playerOrder))
	copy(players, gameState.playerOrder)
	state := gameState.state.String()
	gameState.mu.RUnlock()

	e.mu.Lock()

	// Initialize turn snapshots map for this game if needed
	if e.turnSnapshots[gameID] == nil {
		e.turnSnapshots[gameID] = make(map[int]*gameStateSnapshot)
	}

	// Save snapshot for this turn
	e.turnSnapshots[gameID][turnNumber] = snapshot

	// Remove old snapshots beyond the max
	toDelete := turnNumber - e.rollbackTurnsMax
	if toDelete > 0 {
		delete(e.turnSnapshots[gameID], toDelete)
	}

	if e.logger != nil {
		e.logger.Debug("saved turn snapshot",
			zap.String("game_id", gameID),
			zap.Int("turn", turnNumber),
			zap.Int("snapshots_kept", len(e.turnSnapshots[gameID])),
		)
	}

	e.mu.Unlock()

	// Persist to database if repository is configured
	if persistenceRepo != nil {
		serializedState, err := snapshot.SerializeToBytes()
		if err != nil {
			if e.logger != nil {
				e.logger.Error("failed to serialize game state for persistence",
					zap.String("game_id", gameID),
					zap.Error(err),
				)
			}
			return nil // Don't fail the game for persistence errors
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := persistenceRepo.SaveGameState(ctx, gameID, tableID, gameType, players, serializedState, turnNumber, state); err != nil {
			if e.logger != nil {
				e.logger.Error("failed to persist game state to database",
					zap.String("game_id", gameID),
					zap.Error(err),
				)
			}
			// Don't fail the game for persistence errors
		} else if e.logger != nil {
			e.logger.Debug("persisted game state to database",
				zap.String("game_id", gameID),
				zap.Int("turn", turnNumber),
				zap.Int("state_size", len(serializedState)),
			)
		}
	}

	return nil
}

// PersistGameState manually persists the current game state to the database
// This can be called at important moments like game start, mulligan complete, etc.
func (e *MageEngine) PersistGameState(gameID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	persistenceRepo := e.persistenceRepo
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	if persistenceRepo == nil {
		return nil // No persistence configured
	}

	gameState.mu.RLock()
	snapshot := e.createSnapshot(gameState)
	tableID := gameState.gameID
	gameType := gameState.gameType
	players := make([]string, len(gameState.playerOrder))
	copy(players, gameState.playerOrder)
	turnNumber := gameState.turnManager.TurnNumber()
	state := gameState.state.String()
	gameState.mu.RUnlock()

	serializedState, err := snapshot.SerializeToBytes()
	if err != nil {
		return fmt.Errorf("failed to serialize game state: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := persistenceRepo.SaveGameState(ctx, gameID, tableID, gameType, players, serializedState, turnNumber, state); err != nil {
		return fmt.Errorf("failed to persist game state: %w", err)
	}

	if e.logger != nil {
		e.logger.Debug("persisted game state",
			zap.String("game_id", gameID),
			zap.Int("turn", turnNumber),
			zap.String("state", state),
		)
	}

	return nil
}

// DeletePersistedGame removes a game from the persistence database
// This should be called when a game finishes
func (e *MageEngine) DeletePersistedGame(gameID string) error {
	e.mu.RLock()
	persistenceRepo := e.persistenceRepo
	e.mu.RUnlock()

	if persistenceRepo == nil {
		return nil // No persistence configured
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return persistenceRepo.DeleteActiveGame(ctx, gameID)
}

// LoadGameFromSnapshot restores a game from a serialized snapshot
// Used for server restart recovery - recreates the full game state
func (e *MageEngine) LoadGameFromSnapshot(gameID, tableID, gameType string, players []string, serializedState []byte) error {
	if e.logger != nil {
		e.logger.Info("[RESTORE] LoadGameFromSnapshot starting",
			zap.String("game_id", gameID),
			zap.String("table_id", tableID),
			zap.String("game_type", gameType),
			zap.Strings("players", players),
			zap.Int("serialized_state_size", len(serializedState)),
		)
	}

	// Deserialize the snapshot
	snapshot, err := DeserializeFromBytes(serializedState)
	if err != nil {
		if e.logger != nil {
			e.logger.Error("[RESTORE] Failed to deserialize snapshot",
				zap.String("game_id", gameID),
				zap.Error(err),
			)
		}
		return fmt.Errorf("failed to deserialize game state: %w", err)
	}

	if e.logger != nil {
		e.logger.Info("[RESTORE] Snapshot deserialized successfully",
			zap.String("game_id", gameID),
			zap.String("snapshot_game_id", snapshot.GameID),
			zap.String("snapshot_state", snapshot.State.String()),
			zap.Int("snapshot_turn", snapshot.TurnNumber),
			zap.String("snapshot_active_player", snapshot.ActivePlayer),
			zap.String("snapshot_priority_player", snapshot.PriorityPlayer),
			zap.Int("snapshot_player_count", len(snapshot.Players)),
			zap.Int("snapshot_card_count", len(snapshot.Cards)),
			zap.Int("snapshot_battlefield_count", len(snapshot.Battlefield)),
			zap.Int("snapshot_prompt_count", len(snapshot.Prompts)),
			zap.Int("snapshot_message_count", len(snapshot.Messages)),
		)

		// Log detailed player state from snapshot
		for playerID, player := range snapshot.Players {
			e.logger.Info("[RESTORE] Snapshot player state",
				zap.String("game_id", gameID),
				zap.String("player_id", playerID),
				zap.String("player_name", player.Name),
				zap.Int("life", player.Life),
				zap.Int("hand_size", len(player.Hand)),
				zap.Int("library_size", len(player.Library)),
				zap.Int("graveyard_size", len(player.Graveyard)),
				zap.Bool("kept_hand", player.KeptHand),
				zap.Int("mulligan_count", player.MulliganCount),
				zap.Bool("passed", player.Passed),
				zap.Bool("lost", player.Lost),
				zap.Bool("left", player.Left),
			)
		}

		// Log prompts from snapshot
		for i, prompt := range snapshot.Prompts {
			e.logger.Info("[RESTORE] Snapshot prompt",
				zap.String("game_id", gameID),
				zap.Int("prompt_index", i),
				zap.String("player_id", prompt.PlayerID),
				zap.String("text", prompt.Text),
				zap.Strings("options", prompt.Options),
			)
		}
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.games[gameID]; exists {
		if e.logger != nil {
			e.logger.Warn("[RESTORE] Game already exists, cannot restore",
				zap.String("game_id", gameID),
			)
		}
		return fmt.Errorf("game %s already exists", gameID)
	}

	// Look up game type configuration from registry
	var gameTypeConfig plugin.GameType
	gameRules := plugin.DefaultGameRules()
	var behaviors []plugin.GameBehavior

	if gt, err := plugin.GetGameType(gameType); err == nil {
		gameTypeConfig = gt
		gameRules = plugin.GetRulesForGameType(gt)
		behaviors = plugin.GetBehaviorsForGameType(gt)
	}

	// Create a new game state structure
	gameState := &engineGameState{
		gameID:           gameID,
		gameType:         gameType,
		gameTypeConfig:   gameTypeConfig,
		gameRules:        gameRules,
		behaviors:        behaviors,
		state:            snapshot.State,
		players:          make(map[string]*internalPlayer),
		playerOrder:      make([]string, len(snapshot.PlayerOrder)),
		cards:            make(map[string]*internalCard),
		battlefield:      make([]*internalCard, 0),
		exile:            make([]*internalCard, 0),
		command:          make([]*internalCard, 0),
		revealed:         make([]EngineRevealedView, 0),
		lookedAt:         make([]EngineLookedAtView, 0),
		combat:           newCombatState(),
		triggeredQueue:   make([]*triggeredAbilityQueueItem, 0),
		combatTriggers:   make([]*combatTrigger, 0),
		concedingPlayers: make([]string, 0),
		analytics: &gameAnalytics{
			actionsPerTurn: make(map[int]int),
			turnStartTimes: make(map[int]time.Time),
			gameStartTime:  time.Now(),
		},
		messages:         make([]EngineMessage, 0),
		prompts:          make([]EnginePrompt, 0),
		startedAt:        time.Now(), // Use current time for restored game
		lki:              make(map[string]*LastKnownInfo),
		lkiZoneCounter:   make(map[string]int),
		messageBookmarks: make(map[int]int),
		nextMessageID:    1,
	}

	// Restore player order
	copy(gameState.playerOrder, snapshot.PlayerOrder)

	// Restore players
	for id, player := range snapshot.Players {
		gameState.players[id] = player
	}

	// Restore cards
	for id, card := range snapshot.Cards {
		gameState.cards[id] = card
	}

	// Restore zones
	gameState.battlefield = append(gameState.battlefield, snapshot.Battlefield...)
	gameState.exile = append(gameState.exile, snapshot.Exile...)
	gameState.command = append(gameState.command, snapshot.Command...)

	// Initialize turn manager and restore turn state
	activePlayer := snapshot.ActivePlayer
	if activePlayer == "" && len(players) > 0 {
		activePlayer = players[0]
	}
	gameState.turnManager = rules.NewTurnManager(activePlayer)
	gameState.turnManager.RestoreTurnState(snapshot.TurnNumber, snapshot.ActivePlayer, snapshot.PriorityPlayer)

	// Initialize stack manager and restore stack
	gameState.stack = rules.NewStackManager()
	for _, item := range snapshot.StackItems {
		gameState.stack.Push(item)
	}

	// Initialize event bus
	gameState.eventBus = rules.NewEventBus()

	// Initialize watcher registry
	gameState.watchers = rules.NewWatcherRegistry()

	// Initialize legality checker and target validator
	gameState.legality = rules.NewLegalityChecker(gameState)
	gameState.targetValidator = targeting.NewTargetValidator(gameState)

	// Initialize layer system
	gameState.layerSystem = effects.NewLayerSystem()

	// Initialize ability registry
	gameState.abilityRegistry = NewAbilityRegistry()

	// Restore messages
	gameState.messages = append(gameState.messages, snapshot.Messages...)

	// Restore prompts
	gameState.prompts = append(gameState.prompts, snapshot.Prompts...)

	// Add a restoration message
	gameState.addMessage(fmt.Sprintf("Game restored from server persistence (turn %d)", snapshot.TurnNumber), "system")

	// Register the game
	e.games[gameID] = gameState

	// Initialize turn snapshots and bookmarks for this game
	e.turnSnapshots[gameID] = make(map[int]*gameStateSnapshot)
	e.bookmarks[gameID] = make([]*gameStateSnapshot, 0)

	// Initialize replacement effects manager for this game
	e.replacementEffects[gameID] = effects.NewReplacementManager(e.logger)

	// Fix inconsistent state: if state is MULLIGAN but all players have kept their hands,
	// the game was restored from a snapshot taken before mulligan completion was persisted.
	// In this case, automatically transition to IN_PROGRESS to fix the inconsistent state.
	if e.logger != nil {
		e.logger.Info("[RESTORE] Checking for inconsistent mulligan state",
			zap.String("game_id", gameID),
			zap.String("current_state", gameState.state.String()),
			zap.Bool("is_mulligan_state", gameState.state == GameStateMulligan),
		)
	}

	if gameState.state == GameStateMulligan {
		allKept := true
		playersNotKept := []string{}
		playersKept := []string{}
		for playerID, player := range gameState.players {
			if !player.KeptHand {
				allKept = false
				playersNotKept = append(playersNotKept, playerID)
			} else {
				playersKept = append(playersKept, playerID)
			}
		}

		if e.logger != nil {
			e.logger.Info("[RESTORE] Mulligan state analysis",
				zap.String("game_id", gameID),
				zap.Bool("all_kept", allKept),
				zap.Int("total_players", len(gameState.players)),
				zap.Strings("players_kept", playersKept),
				zap.Strings("players_not_kept", playersNotKept),
			)
		}

		if allKept && len(gameState.players) > 0 {
			if e.logger != nil {
				e.logger.Warn("[RESTORE] FIXING INCONSISTENT STATE: MULLIGAN state but all players kept hands",
					zap.String("game_id", gameID),
					zap.Int("turn", snapshot.TurnNumber),
					zap.String("old_state", "MULLIGAN"),
					zap.String("new_state", "IN_PROGRESS"),
					zap.Int("prompt_count_before_clear", len(gameState.prompts)),
				)
			}
			gameState.state = GameStateInProgress
			gameState.addMessage("Game state corrected: mulligan phase was already complete", "system")

			// Clear any stale mulligan prompts
			gameState.prompts = nil

			if e.logger != nil {
				e.logger.Info("[RESTORE] State fix complete",
					zap.String("game_id", gameID),
					zap.String("final_state", gameState.state.String()),
				)
			}
		} else if !allKept {
			// Legitimate mulligan state - some players haven't kept hands
			if e.logger != nil {
				e.logger.Info("[RESTORE] Legitimate MULLIGAN state - waiting for players to keep hands",
					zap.String("game_id", gameID),
					zap.Strings("players_pending", playersNotKept),
				)
			}
		}
	}

	if e.logger != nil {
		e.logger.Info("[RESTORE] Game restoration complete",
			zap.String("game_id", gameID),
			zap.String("game_type", gameType),
			zap.Int("turn", snapshot.TurnNumber),
			zap.String("final_state", gameState.state.String()),
			zap.Int("player_count", len(players)),
			zap.Int("battlefield_cards", len(gameState.battlefield)),
			zap.Int("stack_items", len(gameState.stack.List())),
			zap.Int("active_prompts", len(gameState.prompts)),
		)
	}

	return nil
}

// GetRestoredGameIDs returns a list of all game IDs currently loaded in the engine
// Useful for verifying restoration completed correctly
func (e *MageEngine) GetRestoredGameIDs() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	ids := make([]string, 0, len(e.games))
	for id := range e.games {
		ids = append(ids, id)
	}
	return ids
}

// CanRollbackTurns checks if it's possible to rollback N turns
// Per Java GameImpl.canRollbackTurns(): validates rollback is possible
func (e *MageEngine) CanRollbackTurns(gameID string, turnsToRollback int) (bool, error) {
	if !e.rollbackAllowed {
		return false, fmt.Errorf("turn rollback is disabled")
	}

	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("game %s not found", gameID)
	}

	currentTurn := gameState.turnManager.TurnNumber()
	targetTurn := currentTurn - turnsToRollback

	if targetTurn < 1 {
		return false, nil // Can't rollback before turn 1
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	turnSnaps := e.turnSnapshots[gameID]
	if turnSnaps == nil {
		return false, nil
	}

	_, exists = turnSnaps[targetTurn]
	return exists, nil
}

// RollbackTurns rolls back the game to N turns ago
// Per Java GameImpl.rollbackTurns(): requires all players to agree (not implemented yet)
func (e *MageEngine) RollbackTurns(gameID string, turnsToRollback int) error {
	if !e.rollbackAllowed {
		return fmt.Errorf("turn rollback is disabled")
	}

	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	currentTurn := gameState.turnManager.TurnNumber()
	targetTurn := currentTurn - turnsToRollback

	if targetTurn < 1 {
		return fmt.Errorf("cannot rollback to turn %d (before game start)", targetTurn)
	}

	e.mu.RLock()
	turnSnaps := e.turnSnapshots[gameID]
	snapshot, exists := turnSnaps[targetTurn]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no snapshot available for turn %d", targetTurn)
	}

	// Restore game state from turn snapshot
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Restore game state from snapshot
	gameState.state = snapshot.State
	gameState.gameType = snapshot.GameType

	// Restore players
	gameState.players = make(map[string]*internalPlayer)
	for id, player := range snapshot.Players {
		// Clear all player stored bookmarks on turn rollback
		// Per Java: resetStoredBookmark for all players
		player.StoredBookmark = -1
		gameState.players[id] = player
	}
	gameState.playerOrder = append([]string(nil), snapshot.PlayerOrder...)

	// Restore cards
	gameState.cards = make(map[string]*internalCard)
	for id, card := range snapshot.Cards {
		gameState.cards[id] = card
	}

	// Restore zones
	gameState.battlefield = append([]*internalCard(nil), snapshot.Battlefield...)
	gameState.exile = append([]*internalCard(nil), snapshot.Exile...)
	gameState.command = append([]*internalCard(nil), snapshot.Command...)

	// Restore stack
	gameState.stack = rules.NewStackManager()
	for _, item := range snapshot.StackItems {
		gameState.stack.Push(item)
	}

	// Restore messages and prompts
	gameState.messages = append([]EngineMessage(nil), snapshot.Messages...)
	gameState.prompts = append([]EnginePrompt(nil), snapshot.Prompts...)

	// Clear all action bookmarks (they're invalid after turn rollback)
	// Per Java: savedStates.clear() and gameStates.clear()
	e.mu.Lock()
	delete(e.bookmarks, gameID)
	e.bookmarks[gameID] = make([]*gameStateSnapshot, 0)
	e.mu.Unlock()

	gameState.addMessage(fmt.Sprintf("Game rolled back to start of turn %d", targetTurn), "system")

	if e.logger != nil {
		e.logger.Info("rolled back turns",
			zap.String("game_id", gameID),
			zap.Int("from_turn", currentTurn),
			zap.Int("to_turn", targetTurn),
			zap.Int("turns_rolled_back", turnsToRollback),
		)
	}

	// Notify players of the rollback
	e.notifyGameStateChange(gameID, map[string]interface{}{
		"type":              "turn_rollback",
		"from_turn":         currentTurn,
		"to_turn":           targetTurn,
		"turns_rolled_back": turnsToRollback,
	})

	return nil
}

// CleanupGame removes a game and frees all associated resources
// Per Java GameImpl.cleanUp(): dispose of game resources, clear watchers, remove listeners
func (e *MageEngine) CleanupGame(gameID string) error {
	e.mu.Lock()
	gameState, exists := e.games[gameID]
	if !exists {
		e.mu.Unlock()
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()

	// Clear all bookmarks
	delete(e.bookmarks, gameID)

	// Clear turn snapshots
	delete(e.turnSnapshots, gameID)

	// Clear watchers
	if gameState.watchers != nil {
		gameState.watchers.Clear()
	}

	// Clear effect managers (Rule 614)
	delete(e.replacementEffects, gameID)

	// Remove game from engine
	delete(e.games, gameID)

	gameState.mu.Unlock()
	e.mu.Unlock()

	// Clear replay from memory (saves should be done before cleanup)
	e.replayRecorder.ClearReplay(gameID)

	if e.logger != nil {
		e.logger.Info("cleaned up game",
			zap.String("game_id", gameID),
		)
	}

	// Notify cleanup complete (safe to call after releasing locks)
	e.notifyGameStateChange(gameID, map[string]interface{}{
		"type": "game_cleanup",
	})

	return nil
}

// StartMulligan transitions game to mulligan phase
// Per Java GameImpl.start(): mulligan phase before main game
func (e *MageEngine) StartMulligan(gameID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	gameState.state = GameStateMulligan

	if e.logger != nil {
		e.logger.Info("started mulligan phase",
			zap.String("game_id", gameID),
		)
	}

	e.notifyGameStateChange(gameID, map[string]interface{}{
		"type":  "mulligan_started",
		"state": "MULLIGAN",
	})

	return nil
}

// PlayerMulligan performs a mulligan for a player (London mulligan)
// Per Java LondonMulligan.mulligan(): shuffle hand into library, draw N-1 cards
func (e *MageEngine) PlayerMulligan(gameID, playerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	if gameState.state != GameStateMulligan {
		return fmt.Errorf("game is not in mulligan phase")
	}

	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	if player.KeptHand {
		return fmt.Errorf("player has already kept their hand")
	}

	// Shuffle hand back into library
	player.Library = append(player.Library, player.Hand...)
	player.Hand = make([]*internalCard, 0)

	// Shuffle library using proper crypto/rand shuffle
	e.shuffleLibrary(player)

	// Increment mulligan count
	player.MulliganCount++

	// Draw N - mulliganCount cards (London mulligan)
	handSize := 7 - player.MulliganCount
	if handSize < 0 {
		handSize = 0
	}

	for i := 0; i < handSize && len(player.Library) > 0; i++ {
		card := player.Library[0]
		player.Library = player.Library[1:]
		card.Zone = zoneHand
		player.Hand = append(player.Hand, card)
	}

	gameState.addMessage(fmt.Sprintf("%s mulligans to %d cards", player.Name, handSize), "mulligan")

	if e.logger != nil {
		e.logger.Info("player mulliganed",
			zap.String("game_id", gameID),
			zap.String("player_id", playerID),
			zap.Int("mulligan_count", player.MulliganCount),
			zap.Int("hand_size", handSize),
		)
	}

	e.notifyGameStateChange(gameID, map[string]interface{}{
		"type":           "player_mulligan",
		"player_id":      playerID,
		"mulligan_count": player.MulliganCount,
		"hand_size":      handSize,
	})

	return nil
}

// PlayerKeepHand indicates player is keeping their current hand
// Per Java LondonMulligan.endMulligan(): finalize mulligan, bottom cards if needed
func (e *MageEngine) PlayerKeepHand(gameID, playerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	if gameState.state != GameStateMulligan {
		return fmt.Errorf("game is not in mulligan phase")
	}

	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	player.KeptHand = true

	gameState.addMessage(fmt.Sprintf("%s keeps their hand", player.Name), "mulligan")

	if e.logger != nil {
		e.logger.Info("player kept hand",
			zap.String("game_id", gameID),
			zap.String("player_id", playerID),
			zap.Int("mulligan_count", player.MulliganCount),
		)
	}

	e.notifyGameStateChange(gameID, map[string]interface{}{
		"type":      "player_keep_hand",
		"player_id": playerID,
	})

	return nil
}

// EndMulligan ends the mulligan phase and starts the main game
// Per Java GameImpl.endMulligan(): transition to main game after all players keep
func (e *MageEngine) EndMulligan(gameID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	if gameState.state != GameStateMulligan {
		return fmt.Errorf("game is not in mulligan phase")
	}

	// Check all players have kept their hands
	for _, player := range gameState.players {
		if !player.KeptHand {
			return fmt.Errorf("not all players have kept their hands")
		}
	}

	// Transition to main game
	gameState.state = GameStateInProgress

	gameState.addMessage("Mulligan phase complete, game starting", "system")

	if e.logger != nil {
		e.logger.Info("mulligan phase ended",
			zap.String("game_id", gameID),
		)
	}

	e.notifyGameStateChange(gameID, map[string]interface{}{
		"type":  "mulligan_ended",
		"state": "IN_PROGRESS",
	})

	return nil
}

// Combat System Implementation
// Per Java Combat class

// processForcedAttackers processes "attacks if able" effects and automatically declares forced attackers
// Per Java Combat.checkAttackRequirements() (lines 474-591)
// NOTE: This function is called while holding gameState.mu.Lock(), so it uses internal
// helper methods that don't acquire locks to avoid self-deadlock.
func (e *MageEngine) processForcedAttackers(gameState *engineGameState) error {
	activePlayerID := gameState.combat.attackingPlayerID

	// Find all creatures that must attack
	for _, card := range gameState.cards {
		if card.Zone != zoneBattlefield {
			continue
		}
		if card.ControllerID != activePlayerID {
			continue
		}
		if !e.isCreature(card) {
			continue
		}
		if card.Attacking {
			continue // Already attacking
		}

		// Check if creature has "attacks if able" effect
		if !e.hasMustAttackEffect(gameState, card.ID) {
			continue
		}

		// Check if creature can attack (using internal version that doesn't acquire locks)
		if !e.canAttackInternal(gameState, card) {
			continue // Can't attack due to restrictions (tapped, summoning sickness, etc.)
		}

		// Find valid defenders this creature can attack
		validDefenders := make([]string, 0)
		for defenderID := range gameState.combat.defenders {
			if canAttack, _ := e.canAttackDefenderInternal(gameState, card, defenderID); canAttack {
				validDefenders = append(validDefenders, defenderID)
			}
		}

		if len(validDefenders) == 0 {
			continue // No valid targets
		}

		// Pick the first valid defender (in a real implementation, player would choose)
		// For now, just attack the first defender in the map
		defenderID := validDefenders[0]

		// Declare the attacker (using internal version that doesn't acquire locks)
		if err := e.declareAttackerInternal(gameState, card.ID, defenderID, activePlayerID); err != nil {
			if e.logger != nil {
				e.logger.Warn("failed to declare forced attacker",
					zap.String("game_id", gameState.gameID),
					zap.String("creature_id", card.ID),
					zap.Error(err),
				)
			}
			continue
		}

		// Track this as a forced attack
		if gameState.combat.creaturesForcedToAttack[card.ID] == nil {
			gameState.combat.creaturesForcedToAttack[card.ID] = make(map[string]bool)
		}
		gameState.combat.creaturesForcedToAttack[card.ID][defenderID] = true

		if e.logger != nil {
			e.logger.Debug("declared forced attacker",
				zap.String("game_id", gameState.gameID),
				zap.String("creature_id", card.ID),
				zap.String("defender_id", defenderID),
			)
		}
	}

	return nil
}

// processMustBeBlockedRequirements processes "must be blocked if able" effects
// Per Java Combat.retrieveMustBlockAttackerRequirements() (lines 848-891)
func (e *MageEngine) processMustBeBlockedRequirements(gameState *engineGameState) error {
	// For each attacker, check if it has a "must be blocked" effect
	for attackerID := range gameState.combat.attackers {
		attacker, exists := gameState.cards[attackerID]
		if !exists || !attacker.Attacking {
			continue
		}

		// Get all "must be blocked" effects on this attacker
		mbEffects := e.getMustBeBlockedEffects(gameState, attackerID)
		if len(mbEffects) == 0 {
			continue
		}

		// Find the defender being attacked
		defenderID := attacker.AttackingWhat
		if defenderID == "" {
			continue
		}

		// Determine defending player
		defendingPlayerID := ""
		if defender, exists := gameState.players[defenderID]; exists {
			defendingPlayerID = defender.PlayerID
		} else {
			// Defender is a permanent (planeswalker/battle), find its controller
			if defenderCard, exists := gameState.cards[defenderID]; exists {
				defendingPlayerID = defenderCard.ControllerID
			}
		}

		if defendingPlayerID == "" {
			continue
		}

		// Find all potential blockers controlled by defending player
		for _, blocker := range gameState.cards {
			if blocker.Zone != zoneBattlefield {
				continue
			}
			if blocker.ControllerID != defendingPlayerID {
				continue
			}
			if !e.isCreature(blocker) {
				continue
			}

			// Check if this blocker can block the attacker
			canBlock, _ := e.CanBlock(gameState.gameID, blocker.ID, attackerID)
			if !canBlock {
				continue
			}

			// Add this blocker to the must-block list for this attacker
			if gameState.combat.creatureMustBlockAttackers[blocker.ID] == nil {
				gameState.combat.creatureMustBlockAttackers[blocker.ID] = make(map[string]bool)
			}
			gameState.combat.creatureMustBlockAttackers[blocker.ID][attackerID] = true

			if e.logger != nil {
				e.logger.Debug("creature must block attacker",
					zap.String("game_id", gameState.gameID),
					zap.String("blocker_id", blocker.ID),
					zap.String("attacker_id", attackerID),
				)
			}
		}
	}

	return nil
}

// handleStepBegin handles step-specific actions when entering a new step
// Per MTG Rules 502-514 for turn-based actions
func (e *MageEngine) handleStepBegin(gameState *engineGameState, step rules.Step, activePlayerID string) {
	switch step {
	case rules.StepUntap:
		// Per MTG Rule 502.2: Untap all tapped permanents controlled by active player
		// Skip this on turn 1 (first player's first turn)
		e.performUntapStep(gameState, activePlayerID)

	case rules.StepUpkeep:
		// Per MTG Rule 503: No turn-based actions, but triggers happen here
		gameState.eventBus.Publish(rules.NewEvent(rules.EventUpkeepStep, "", "", activePlayerID))

	case rules.StepDraw:
		// Per MTG Rule 504.1: Active player draws a card
		// Per MTG Rule 103.7a: The starting player skips their draw step on turn 1
		e.performDrawStep(gameState, activePlayerID)

	case rules.StepCleanup:
		// Per MTG Rule 514.1-514.3: Discard to hand size, remove damage, end "until end of turn" effects
		e.performCleanupStep(gameState, activePlayerID)
	}
}

// performUntapStep untaps all permanents controlled by the active player
// Per MTG Rule 502.2
func (e *MageEngine) performUntapStep(gameState *engineGameState, activePlayerID string) {
	// Reset lands played this turn for active player (beginning of turn)
	if player, exists := gameState.players[activePlayerID]; exists {
		player.LandsPlayedThisTurn = 0
	}

	untappedCount := 0
	summoningSicknessCleared := 0
	for _, card := range gameState.battlefield {
		if card == nil {
			continue
		}
		if card.ControllerID == activePlayerID {
			// Untap tapped permanents
			if card.Tapped {
				// TODO: Check for "doesn't untap" effects
				card.Tapped = false
				untappedCount++
			}
			// Clear summoning sickness for creatures controlled by the active player
			// Per MTG Rule 302.6: A creature can attack/tap if it has been under its controller's
			// control continuously since the start of their most recent turn
			if card.SummoningSickness {
				card.SummoningSickness = false
				summoningSicknessCleared++
			}
		}
	}

	if untappedCount > 0 {
		gameState.addMessage(fmt.Sprintf("%s untaps %d permanents", activePlayerID, untappedCount), "action")
	}

	gameState.eventBus.Publish(rules.NewEvent(rules.EventUntapStep, "", "", activePlayerID))

	if e.logger != nil {
		e.logger.Debug("untap step performed",
			zap.String("player", activePlayerID),
			zap.Int("untapped", untappedCount),
			zap.Int("summoning_sickness_cleared", summoningSicknessCleared),
		)
	}
}

// performDrawStep has the active player draw a card
// Per MTG Rule 504.1 and 103.7a (starting player skips first draw)
func (e *MageEngine) performDrawStep(gameState *engineGameState, activePlayerID string) {
	// Per MTG Rule 103.7a: The starting player skips their draw step on the first turn
	turnNumber := gameState.turnManager.TurnNumber()
	if turnNumber == 1 && activePlayerID == gameState.startingPlayerID && !gameState.firstTurnDrawDone {
		gameState.firstTurnDrawDone = true
		gameState.addMessage(fmt.Sprintf("%s skips their first turn draw", activePlayerID), "action")
		gameState.eventBus.Publish(rules.NewEvent(rules.EventDrawStep, "", "", activePlayerID))

		if e.logger != nil {
			e.logger.Debug("first turn draw skipped",
				zap.String("player", activePlayerID),
			)
		}
		return
	}

	player, exists := gameState.players[activePlayerID]
	if !exists {
		if e.logger != nil {
			e.logger.Error("player not found for draw step", zap.String("player", activePlayerID))
		}
		return
	}

	if len(player.Library) == 0 {
		// Player loses if they can't draw from empty library
		// Per MTG Rule 704.5b
		player.Lost = true
		gameState.addMessage(fmt.Sprintf("%s cannot draw from empty library and loses the game", activePlayerID), "action")
		if e.logger != nil {
			e.logger.Info("player loses by drawing from empty library",
				zap.String("player", activePlayerID),
			)
		}
		return
	}

	// Draw from top of library
	card := player.Library[len(player.Library)-1]
	player.Library = player.Library[:len(player.Library)-1]
	card.Zone = zoneHand
	player.Hand = append(player.Hand, card)

	gameState.addMessage(fmt.Sprintf("%s draws a card", activePlayerID), "action")
	gameState.eventBus.Publish(rules.NewEvent(rules.EventDrawStep, "", "", activePlayerID))

	if e.logger != nil {
		e.logger.Debug("draw step performed",
			zap.String("player", activePlayerID),
			zap.String("card", card.Name),
			zap.Int("hand_size", len(player.Hand)),
		)
	}
}

// performCleanupStep handles end of turn cleanup
// Per MTG Rule 514.1-514.3
func (e *MageEngine) performCleanupStep(gameState *engineGameState, activePlayerID string) {
	player, exists := gameState.players[activePlayerID]
	if !exists {
		return
	}

	// 514.1: Discard to maximum hand size (7 by default)
	maxHandSize := 7
	if len(player.Hand) > maxHandSize {
		discardCount := len(player.Hand) - maxHandSize
		actualDiscarded := 0
		// TODO: Let player choose which cards to discard
		// For now, discard from the end of hand
		for i := 0; i < discardCount && len(player.Hand) > maxHandSize; i++ {
			card := player.Hand[len(player.Hand)-1]
			player.Hand = player.Hand[:len(player.Hand)-1]
			if card == nil {
				continue
			}
			card.Zone = zoneGraveyard
			player.Graveyard = append(player.Graveyard, card)
			actualDiscarded++
		}
		if actualDiscarded > 0 {
			gameState.addMessage(fmt.Sprintf("%s discards %d cards to hand size", activePlayerID, actualDiscarded), "action")
		}
	}

	// 514.2: Remove all damage from creatures
	for _, card := range gameState.battlefield {
		if card == nil {
			continue
		}
		if strings.Contains(card.Type, "Creature") && card.Damage > 0 {
			card.Damage = 0
		}
	}

	// 514.3: "Until end of turn" and "this turn" effects end
	// This is handled by effects.CleanupEndOfTurnEffects called in handlePass

	gameState.eventBus.Publish(rules.NewEvent(rules.EventCleanupStep, "", "", activePlayerID))

	if e.logger != nil {
		e.logger.Debug("cleanup step performed",
			zap.String("player", activePlayerID),
		)
	}
}

// handleCombatStepBegin handles combat initialization when entering combat steps
// Per Java BeginCombatStep.beginStep() and DeclareAttackersStep.beginStep()
func (e *MageEngine) handleCombatStepBegin(gameState *engineGameState, step rules.Step, activePlayerID string) {
	switch step {
	case rules.StepBeginCombat:
		// Per Java BeginCombatStep.beginStep() (line 28-33)
		// 507.1: At the start of the combat phase, if an opponent controls a permanent with
		// "At the beginning of combat" triggered ability, or if a turn-based action of a player
		// other than the active player occurs at the beginning of combat, the active player gets priority

		// Create new combat state (equivalent to game.getCombat().clear())
		gameState.combat = newCombatState()

		// Clear combat flags on all cards
		for _, card := range gameState.cards {
			card.Attacking = false
			card.Blocking = false
			card.AttackingWhat = ""
			card.BlockingWhat = nil
		}

		// Set the attacking player (equivalent to game.getCombat().setAttacker(activePlayerId))
		gameState.combat.attackingPlayerID = activePlayerID

		// Set defenders (equivalent to game.getCombat().setDefenders(game))
		// Clear previous defenders
		gameState.combat.defenders = make(map[string]bool)

		// Add all opponents as defenders
		for playerID := range gameState.players {
			if playerID != activePlayerID {
				gameState.combat.defenders[playerID] = true
			}
		}

		// Fire begin combat event
		gameState.eventBus.Publish(rules.NewEvent(rules.EventBeginCombatStep, "", "", ""))

		if e.logger != nil {
			e.logger.Debug("begin combat step initialized",
				zap.String("game_id", gameState.gameID),
				zap.String("attacker", activePlayerID),
				zap.Int("defenders", len(gameState.combat.defenders)),
			)
		}

	case rules.StepDeclareAttackers:
		// Per Java DeclareAttackersStep.beginStep() (line 33-36)
		// This is where selectAttackers() would be called in Java
		// In our implementation, attacker selection is handled via player actions
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] StepDeclareAttackers publishing event",
				zap.String("game_id", gameState.gameID))
		}
		gameState.eventBus.Publish(rules.NewEvent(rules.EventDeclareAttackersStepPre, "", "", activePlayerID))

		// Check for creatures that lost creature type (Per Java Combat.checkForRemoveFromCombat())
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] StepDeclareAttackers releasing gameState.mu for CheckForRemoveFromCombat",
				zap.String("game_id", gameState.gameID))
		}
		gameState.mu.Unlock()
		if err := e.CheckForRemoveFromCombat(gameState.gameID); err != nil && e.logger != nil {
			e.logger.Error("failed to check for removal from combat",
				zap.String("game_id", gameState.gameID),
				zap.Error(err),
			)
		}
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] StepDeclareAttackers re-acquiring gameState.mu after CheckForRemoveFromCombat",
				zap.String("game_id", gameState.gameID))
		}
		gameState.mu.Lock()
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] StepDeclareAttackers re-acquired gameState.mu",
				zap.String("game_id", gameState.gameID))
		}

		// Process forced attackers ("attacks if able" effects)
		// Per Java Combat.checkAttackRequirements()
		if err := e.processForcedAttackers(gameState); err != nil && e.logger != nil {
			e.logger.Error("failed to process forced attackers",
				zap.String("game_id", gameState.gameID),
				zap.Error(err),
			)
		}

		// Generate prompt for attacking player to declare attackers
		options := e.buildAttackerPromptOptions(gameState)
		var promptMessage string
		if len(options) > 1 { // More than just "DONE_ATTACKING"
			promptMessage = "Declare attackers (select creatures to attack)"
			gameState.addPrompt(activePlayerID, promptMessage, options)
		} else {
			promptMessage = "No creatures can attack"
			options = []string{"DONE_ATTACKING"}
			gameState.addPrompt(activePlayerID, promptMessage, options)
		}

		// Send choice prompt notification to client
		e.notifyChoicePrompt(gameState.gameID, activePlayerID, promptMessage, options)

		if e.logger != nil {
			e.logger.Debug("declare attackers step initialized",
				zap.String("game_id", gameState.gameID),
				zap.String("active_player", activePlayerID),
				zap.Int("available_options", len(options)),
			)
		}

	case rules.StepDeclareBlockers:
		// Fire the pre-step event for declare blockers
		gameState.eventBus.Publish(rules.NewEvent(rules.EventDeclareBlockersStepPre, "", "", activePlayerID))

		// Check for creatures that lost creature type (Per Java Combat.checkForRemoveFromCombat())
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] StepDeclareBlockers releasing gameState.mu for CheckForRemoveFromCombat",
				zap.String("game_id", gameState.gameID))
		}
		gameState.mu.Unlock()
		if err := e.CheckForRemoveFromCombat(gameState.gameID); err != nil && e.logger != nil {
			e.logger.Error("failed to check for removal from combat",
				zap.String("game_id", gameState.gameID),
				zap.Error(err),
			)
		}
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] StepDeclareBlockers re-acquiring gameState.mu after CheckForRemoveFromCombat",
				zap.String("game_id", gameState.gameID))
		}
		gameState.mu.Lock()
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] StepDeclareBlockers re-acquired gameState.mu",
				zap.String("game_id", gameState.gameID))
		}

		// Process "must be blocked if able" requirements
		// Per Java Combat.retrieveMustBlockAttackerRequirements()
		if err := e.processMustBeBlockedRequirements(gameState); err != nil && e.logger != nil {
			e.logger.Error("failed to process must-be-blocked requirements",
				zap.String("game_id", gameState.gameID),
				zap.Error(err),
			)
		}

		// Generate prompts for each defending player to declare blockers
		// Defending players are all opponents of the active player
		for playerID := range gameState.players {
			if playerID != activePlayerID {
				options := e.buildBlockerPromptOptions(gameState, playerID)
				var promptMessage string
				if len(options) > 1 { // More than just "DONE_BLOCKING"
					promptMessage = "Declare blockers (select creatures to block)"
					gameState.addPrompt(playerID, promptMessage, options)
				} else {
					promptMessage = "No creatures can block or no attackers"
					options = []string{"DONE_BLOCKING"}
					gameState.addPrompt(playerID, promptMessage, options)
				}
				// Send choice prompt notification to client
				e.notifyChoicePrompt(gameState.gameID, playerID, promptMessage, options)
			}
		}

		// After blockers are declared, check if there are creatures with first/double strike
		// If so, update the turn sequence to include the first strike damage step
		// Note: Use internal version since we already hold gameState.mu.Lock()
		if e.hasFirstOrDoubleStrikeInternal(gameState) {
			gameState.turnManager.SetHasFirstStrike(true)
			if e.logger != nil {
				e.logger.Debug("first strike damage step added to turn sequence",
					zap.String("game_id", gameState.gameID),
				)
			}
		}

		if e.logger != nil {
			e.logger.Debug("declare blockers step initialized",
				zap.String("game_id", gameState.gameID),
			)
		}

	case rules.StepFirstStrikeDamage:
		// First strike damage step
		// Fire the pre-step event for first strike combat damage
		gameState.eventBus.Publish(rules.NewEvent(rules.EventCombatDamageStepPre, "", "", activePlayerID))

		// Check for creatures that lost creature type (Per Java Combat.checkForRemoveFromCombat())
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] StepFirstStrikeDamage releasing gameState.mu for CheckForRemoveFromCombat",
				zap.String("game_id", gameState.gameID))
		}
		gameState.mu.Unlock()
		if err := e.CheckForRemoveFromCombat(gameState.gameID); err != nil && e.logger != nil {
			e.logger.Error("failed to check for removal from combat",
				zap.String("game_id", gameState.gameID),
				zap.Error(err),
			)
		}
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] StepFirstStrikeDamage re-acquiring gameState.mu after CheckForRemoveFromCombat",
				zap.String("game_id", gameState.gameID))
		}
		gameState.mu.Lock()
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] StepFirstStrikeDamage re-acquired gameState.mu",
				zap.String("game_id", gameState.gameID))
		}

		// Automatically assign and apply first strike damage
		// Note: Use internal versions since we already hold gameState.mu.Lock()
		if err := e.assignCombatDamageInternal(gameState, true); err == nil {
			if err := e.applyCombatDamageInternal(gameState); err != nil && e.logger != nil {
				e.logger.Error("failed to apply first strike damage",
					zap.String("game_id", gameState.gameID),
					zap.Error(err),
				)
			}
		} else if e.logger != nil {
			e.logger.Error("failed to assign first strike damage",
				zap.String("game_id", gameState.gameID),
				zap.Error(err),
			)
		}

		if e.logger != nil {
			e.logger.Debug("first strike damage step initialized and executed",
				zap.String("game_id", gameState.gameID),
			)
		}

	case rules.StepCombatDamage:
		// Fire the pre-step event for combat damage
		gameState.eventBus.Publish(rules.NewEvent(rules.EventCombatDamageStepPre, "", "", activePlayerID))

		// Check for creatures that lost creature type (Per Java Combat.checkForRemoveFromCombat())
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] StepCombatDamage releasing gameState.mu for CheckForRemoveFromCombat",
				zap.String("game_id", gameState.gameID))
		}
		gameState.mu.Unlock()
		if err := e.CheckForRemoveFromCombat(gameState.gameID); err != nil && e.logger != nil {
			e.logger.Error("failed to check for removal from combat",
				zap.String("game_id", gameState.gameID),
				zap.Error(err),
			)
		}
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] StepCombatDamage re-acquiring gameState.mu after CheckForRemoveFromCombat",
				zap.String("game_id", gameState.gameID))
		}
		gameState.mu.Lock()
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] StepCombatDamage re-acquired gameState.mu",
				zap.String("game_id", gameState.gameID))
		}

		// Automatically assign and apply normal damage
		// Note: Use internal versions since we already hold gameState.mu.Lock()
		if err := e.assignCombatDamageInternal(gameState, false); err == nil {
			if err := e.applyCombatDamageInternal(gameState); err != nil && e.logger != nil {
				e.logger.Error("failed to apply normal combat damage",
					zap.String("game_id", gameState.gameID),
					zap.Error(err),
				)
			}
		} else if e.logger != nil {
			e.logger.Error("failed to assign normal combat damage",
				zap.String("game_id", gameState.gameID),
				zap.Error(err),
			)
		}

		if e.logger != nil {
			e.logger.Debug("combat damage step initialized and executed",
				zap.String("game_id", gameState.gameID),
			)
		}

	case rules.StepEndCombat:
		// Fire the pre-step event for end of combat
		gameState.eventBus.Publish(rules.NewEvent(rules.EventEndCombatStepPre, "", "", activePlayerID))

		// End combat and clean up combat state
		// Note: Use internal version since we already hold gameState.mu.Lock()
		if err := e.endCombatInternal(gameState); err != nil && e.logger != nil {
			e.logger.Error("failed to end combat",
				zap.String("game_id", gameState.gameID),
				zap.Error(err),
			)
		}

		if e.logger != nil {
			e.logger.Debug("end combat step initialized",
				zap.String("game_id", gameState.gameID),
			)
		}
	}
}

// ResetCombat clears all combat state at the beginning of combat
// Per Java Combat.reset()
func (e *MageEngine) ResetCombat(gameID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Create new combat state
	gameState.combat = newCombatState()

	// Clear combat flags on all cards
	for _, card := range gameState.cards {
		card.Attacking = false
		card.Blocking = false
		card.AttackingWhat = ""
		card.BlockingWhat = nil
	}

	if e.logger != nil {
		e.logger.Debug("reset combat", zap.String("game_id", gameID))
	}

	// Fire begin combat event
	gameState.eventBus.Publish(rules.NewEvent(rules.EventBeginCombatStep, "", "", ""))

	return nil
}

// SetAttacker sets the attacking player for this combat
// Per Java Combat.setAttacker()
func (e *MageEngine) SetAttacker(gameID, playerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	if _, exists := gameState.players[playerID]; !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	gameState.combat.attackingPlayerID = playerID

	if e.logger != nil {
		e.logger.Debug("set attacking player",
			zap.String("game_id", gameID),
			zap.String("player_id", playerID),
		)
	}

	return nil
}

// SetDefenders identifies all possible defenders (players, planeswalkers, battles)
// Per Java Combat.setDefenders()
func (e *MageEngine) SetDefenders(gameID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	attackingPlayerID := gameState.combat.attackingPlayerID
	if attackingPlayerID == "" {
		return fmt.Errorf("no attacking player set")
	}

	// Clear previous defenders
	gameState.combat.defenders = make(map[string]bool)

	// Add all opponents as defenders
	for playerID := range gameState.players {
		if playerID != attackingPlayerID {
			gameState.combat.defenders[playerID] = true
		}
	}

	// Add planeswalkers controlled by opponents (Rule 306.6, 508.1b)
	// Per Java Combat.setDefenders() - adds planeswalkers to defenders map
	for _, card := range gameState.cards {
		if card.Zone != zoneBattlefield {
			continue
		}

		// Check if this is a planeswalker
		if !e.isPlaneswalker(card) {
			continue
		}

		// Check if controlled by an opponent
		if card.ControllerID == attackingPlayerID {
			continue // Can't attack your own planeswalkers
		}

		// Add planeswalker as a defender
		gameState.combat.defenders[card.ID] = true
	}

	// TODO: Add battles that can be attacked when battle system is implemented

	if e.logger != nil {
		e.logger.Debug("set defenders",
			zap.String("game_id", gameID),
			zap.Int("defender_count", len(gameState.combat.defenders)),
		)
	}

	return nil
}

// CanAttack checks if a creature can attack (any defender)
// Per Java Permanent.canAttack(null, game) and canAttackInPrinciple(null, game)
func (e *MageEngine) CanAttack(gameID, creatureID string) (bool, error) {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	creature, exists := gameState.cards[creatureID]
	if !exists {
		return false, fmt.Errorf("creature %s not found", creatureID)
	}

	// Basic checks (Java: Permanent.canAttack line 1485)
	if creature.Tapped {
		return false, nil
	}

	// Check if can attack in principle (Java: canAttackInPrinciple line 1504)
	// Check summoning sickness
	// TODO: Implement AsThoughEffectType.ATTACK_AS_HASTE for haste effects
	if creature.SummoningSickness {
		return false, nil
	}

	// Check defender ability (Java: line 1527)
	// TODO: Implement AsThoughEffectType.ATTACK for effects that allow defender to attack
	if e.hasAbility(creature, abilityDefender) {
		return false, nil
	}

	// Check for continuous effects that prevent attacking
	// Per Java: RestrictionEffect.applies() and canAttack() checks
	if e.hasCantAttackEffect(gameState, creatureID) {
		return false, nil
	}

	// Check if can attack at least one defender (Java: line 1516-1522)
	// If no specific defender, check if can attack ANY defender
	for defenderID := range gameState.combat.defenders {
		canAttack, _ := e.canAttackDefenderInternal(gameState, creature, defenderID)
		if canAttack {
			return true, nil
		}
	}

	return false, nil
}

// CanAttackDefender checks if a creature can attack a specific defender
// Per Java Permanent.canAttack(defenderId, game) and canAttackInPrinciple(defenderId, game)
func (e *MageEngine) CanAttackDefender(gameID, creatureID, defenderID string) (bool, error) {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	creature, exists := gameState.cards[creatureID]
	if !exists {
		return false, fmt.Errorf("creature %s not found", creatureID)
	}

	// Basic checks (Java: Permanent.canAttack line 1485)
	if creature.Tapped {
		return false, nil
	}

	return e.canAttackDefenderInternal(gameState, creature, defenderID)
}

// canAttackInternal checks if a creature can attack (any defender) without acquiring locks
// This is an internal helper for use when gameState.mu is already held
// Per Java Permanent.canAttack(null, game) and canAttackInPrinciple(null, game)
func (e *MageEngine) canAttackInternal(gameState *engineGameState, creature *internalCard) bool {
	// RULES-LIGHT: All creatures can attack - players handle restrictions manually
	// Original checks (tapped, summoning sickness, defender, cant-attack effects) are removed
	// The UI can still show hints about these conditions
	_ = gameState // Keep parameter for interface compatibility
	return creature != nil && creature.Zone == zoneBattlefield
}

// canAttackDefenderInternal checks if a creature can attack a specific defender (internal helper)
// Per Java Permanent.canAttackInPrinciple(defenderId, game)
func (e *MageEngine) canAttackDefenderInternal(gameState *engineGameState, creature *internalCard, defenderID string) (bool, error) {
	// RULES-LIGHT: Any creature can attack any defender - players handle restrictions
	// Original checks (summoning sickness, defender ability, restriction effects) are removed
	_ = creature   // Keep parameter for interface compatibility
	_ = defenderID // Keep parameter for interface compatibility
	_ = gameState  // Keep parameter for interface compatibility
	return true, nil
}

// declareAttackerInternal declares a creature as an attacker without acquiring locks
// This is an internal helper for use when gameState.mu is already held
// RULES-LIGHT: Validation removed - players control what attacks what
func (e *MageEngine) declareAttackerInternal(gameState *engineGameState, creatureID, defenderID, playerID string) error {
	// Basic existence check only
	creature, exists := gameState.cards[creatureID]
	if !exists {
		return fmt.Errorf("creature %s not found", creatureID)
	}

	// RULES-LIGHT: Removed validation for:
	// - Player being the attacking player
	// - Controller matching player
	// - Zone being battlefield
	// - Tapped status
	// - Defender ability

	// Fire declare attackers step pre event (before first attacker)
	if len(gameState.combat.attackers) == 0 {
		gameState.eventBus.Publish(rules.NewEvent(rules.EventDeclareAttackersStepPre, "", "", playerID))
	}

	// RULES-LIGHT: Accept any defender - no validation
	if !gameState.combat.defenders[defenderID] {
		// Add it as a valid defender on the fly
		gameState.combat.defenders[defenderID] = true
	}

	// Determine if defender is a permanent (planeswalker/battle) or player
	defenderIsPermanent := false
	defendingPlayerID := defenderID

	if defenderCard, exists := gameState.cards[defenderID]; exists {
		defenderIsPermanent = true
		defendingPlayerID = defenderCard.ControllerID
	}

	group := newCombatGroup(defenderID, defenderIsPermanent, defendingPlayerID)
	group.attackers = append(group.attackers, creatureID)
	gameState.combat.groups = append(gameState.combat.groups, group)
	gameState.combat.attackers[creatureID] = true

	// Tap creature (unless it has vigilance)
	hasVigilance := e.hasAbilityWithEffects(gameState, creature, abilityVigilance)
	if !hasVigilance && !creature.Tapped {
		creature.Tapped = true
		gameState.combat.attackersTapped[creatureID] = true
	}

	// Set creature combat state
	creature.Attacking = true
	creature.AttackingWhat = defenderID

	// Fire attacker declared event
	event := rules.NewEvent(rules.EventAttackerDeclared, creatureID, creatureID, playerID)
	event.Metadata["defender_id"] = defenderID
	gameState.eventBus.Publish(event)

	// Check for combat triggers
	e.checkCombatTriggers(gameState, event)

	// Fire defender attacked event
	defenderEvent := rules.NewEvent(rules.EventDefenderAttacked, defenderID, creatureID, playerID)
	defenderEvent.Metadata["attacker_id"] = creatureID
	gameState.eventBus.Publish(defenderEvent)

	// Track attacks
	attackingPlayerID := gameState.combat.attackingPlayerID
	if defenderIsPermanent {
		if defenderCard, exists := gameState.cards[defenderID]; exists && e.isPlaneswalker(defenderCard) {
			controllerID := defenderCard.ControllerID
			if gameState.combat.planeswalkerControllersAttackedThisTurn[attackingPlayerID] == nil {
				gameState.combat.planeswalkerControllersAttackedThisTurn[attackingPlayerID] = make(map[string]bool)
			}
			gameState.combat.planeswalkerControllersAttackedThisTurn[attackingPlayerID][controllerID] = true
		}
	} else {
		if gameState.combat.playersAttackedThisTurn[attackingPlayerID] == nil {
			gameState.combat.playersAttackedThisTurn[attackingPlayerID] = make(map[string]bool)
		}
		gameState.combat.playersAttackedThisTurn[attackingPlayerID][defenderID] = true
	}

	gameState.addMessage(fmt.Sprintf("%s attacks", creature.Name), "combat")

	return nil
}

// DeclareAttacker declares a creature as an attacker
// Per Java Combat.declareAttacker()
func (e *MageEngine) DeclareAttacker(gameID, creatureID, defenderID, playerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Validate player
	if playerID != gameState.combat.attackingPlayerID {
		return fmt.Errorf("player %s is not the attacking player", playerID)
	}

	// Validate creature exists and is controlled by player
	creature, exists := gameState.cards[creatureID]
	if !exists {
		return fmt.Errorf("creature %s not found", creatureID)
	}

	if creature.ControllerID != playerID {
		return fmt.Errorf("creature %s is not controlled by player %s", creatureID, playerID)
	}

	// Validate creature is on battlefield
	if creature.Zone != zoneBattlefield {
		return fmt.Errorf("creature %s is not on battlefield", creatureID)
	}

	// Validate creature can attack (not tapped, not summoning sick)
	if creature.Tapped {
		return fmt.Errorf("creature %s is tapped", creatureID)
	}

	// Check for defender ability (Java: PermanentImpl.canAttackInPrinciple line 1527)
	// Creatures with defender can't attack unless they have an effect allowing them to
	if e.hasAbility(creature, abilityDefender) {
		// TODO: Check for AsThoughEffectType.ATTACK effects that allow defender to attack
		return fmt.Errorf("creature %s has defender and cannot attack", creatureID)
	}

	// TODO: Check summoning sickness when we track turn entered
	// TODO: Check for "can't attack" restrictions
	// TODO: Check for "must attack" requirements

	// Fire declare attackers step pre event (before first attacker)
	if len(gameState.combat.attackers) == 0 {
		gameState.eventBus.Publish(rules.NewEvent(rules.EventDeclareAttackersStepPre, "", "", playerID))
	}

	// Validate defender exists
	if !gameState.combat.defenders[defenderID] {
		return fmt.Errorf("invalid defender %s", defenderID)
	}

	// TODO: Validate can attack this specific defender (protection, etc.)

	// Create a new combat group for this attacker
	// Per MTG rules and Java implementation: each attacking creature gets its own combat group
	// Blockers may later be assigned to this group during declare blockers step

	// Determine if defender is a permanent (planeswalker/battle) or player
	// Rule 508.1b, 306.6: Creatures can attack players, planeswalkers, or battles
	defenderIsPermanent := false
	defendingPlayerID := defenderID

	// Check if defender is a planeswalker or battle
	if defenderCard, exists := gameState.cards[defenderID]; exists {
		defenderIsPermanent = true
		defendingPlayerID = defenderCard.ControllerID
	}

	group := newCombatGroup(defenderID, defenderIsPermanent, defendingPlayerID)
	group.attackers = append(group.attackers, creatureID)
	gameState.combat.groups = append(gameState.combat.groups, group)
	gameState.combat.attackers[creatureID] = true

	// Tap creature (unless it has vigilance)
	// Check both base and granted vigilance
	hasVigilance := e.hasAbilityWithEffects(gameState, creature, abilityVigilance)
	if !hasVigilance && !creature.Tapped {
		creature.Tapped = true
		gameState.combat.attackersTapped[creatureID] = true
	}

	// Set creature combat state
	creature.Attacking = true
	creature.AttackingWhat = defenderID

	// Fire attacker declared event
	event := rules.NewEvent(rules.EventAttackerDeclared, creatureID, creatureID, playerID)
	event.Metadata["defender_id"] = defenderID
	gameState.eventBus.Publish(event)

	// Check for combat triggers (e.g., "Whenever ~ attacks")
	e.checkCombatTriggers(gameState, event)

	// Fire defender attacked event
	defenderEvent := rules.NewEvent(rules.EventDefenderAttacked, defenderID, creatureID, playerID)
	defenderEvent.Metadata["attacker_id"] = creatureID
	gameState.eventBus.Publish(defenderEvent)

	// Track attacks for "attacked this turn" queries
	// Per Java PlayersAttackedThisTurnWatcher (lines 62-73)
	attackingPlayerID := gameState.combat.attackingPlayerID
	if defenderIsPermanent {
		// Attacking a planeswalker - track the controller
		if defenderCard, exists := gameState.cards[defenderID]; exists && e.isPlaneswalker(defenderCard) {
			controllerID := defenderCard.ControllerID
			if gameState.combat.planeswalkerControllersAttackedThisTurn[attackingPlayerID] == nil {
				gameState.combat.planeswalkerControllersAttackedThisTurn[attackingPlayerID] = make(map[string]bool)
			}
			gameState.combat.planeswalkerControllersAttackedThisTurn[attackingPlayerID][controllerID] = true
		}
	} else {
		// Attacking a player directly
		if gameState.combat.playersAttackedThisTurn[attackingPlayerID] == nil {
			gameState.combat.playersAttackedThisTurn[attackingPlayerID] = make(map[string]bool)
		}
		gameState.combat.playersAttackedThisTurn[attackingPlayerID][defenderID] = true
	}

	gameState.addMessage(fmt.Sprintf("%s attacks", creature.Name), "combat")

	if e.logger != nil {
		e.logger.Debug("declared attacker",
			zap.String("game_id", gameID),
			zap.String("creature_id", creatureID),
			zap.String("defender_id", defenderID),
		)
	}

	return nil
}

// FinishDeclaringAttackers signals that all attackers have been declared
// Fires the DECLARED_ATTACKERS event
func (e *MageEngine) FinishDeclaringAttackers(gameID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Fire DECLARED_ATTACKERS event
	declaredEvent := rules.NewEvent(rules.EventDeclaredAttackers, "", "", gameState.combat.attackingPlayerID)
	gameState.eventBus.Publish(declaredEvent)

	// Check for combat triggers (e.g., "Whenever one or more creatures attack")
	e.checkCombatTriggers(gameState, declaredEvent)

	return nil
}

// GetCombatView builds the combat view for display
func (e *MageEngine) GetCombatView(gameID string) (EngineCombatView, error) {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return EngineCombatView{}, fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	view := EngineCombatView{
		AttackingPlayerID: gameState.combat.attackingPlayerID,
		Groups:            make([]EngineCombatGroupView, 0, len(gameState.combat.groups)),
	}

	for _, group := range gameState.combat.groups {
		groupView := EngineCombatGroupView{
			Attackers:         make([]string, len(group.attackers)),
			Blockers:          make([]string, len(group.blockers)),
			DefenderID:        group.defenderID,
			DefendingPlayerID: group.defendingPlayerID,
			Blocked:           group.blocked,
		}
		copy(groupView.Attackers, group.attackers)
		copy(groupView.Blockers, group.blockers)
		view.Groups = append(view.Groups, groupView)
	}

	return view, nil
}

// CanBlock checks if a creature can block a specific attacker
// Per Java PermanentImpl.canBlock()
func (e *MageEngine) CanBlock(gameID, blockerID, attackerID string) (bool, error) {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	// Get blocker
	blocker, exists := gameState.cards[blockerID]
	if !exists {
		return false, fmt.Errorf("blocker %s not found", blockerID)
	}

	// Get attacker
	_, exists = gameState.cards[attackerID]
	if !exists {
		return false, fmt.Errorf("attacker %s not found", attackerID)
	}

	// Basic checks
	// 1. Blocker must be untapped (or have ability to block while tapped)
	if blocker.Tapped {
		// TODO: Check for "can block while tapped" abilities
		return false, nil
	}

	// 2. Blocker must be a creature
	if !strings.Contains(blocker.Type, "Creature") {
		return false, nil
	}

	// 3. Blocker must not be a battle
	// TODO: Check for battle type when implemented

	// 4. Blocker must not be suspected
	// TODO: Check for suspected status when implemented

	// 4a. Check for continuous effects that prevent blocking
	// Per Java: RestrictionEffect.applies() and canBlock() checks
	if e.hasCantBlockEffect(gameState, blockerID) {
		return false, nil
	}

	// 5. Blocker must be on battlefield
	if blocker.Zone != zoneBattlefield {
		return false, nil
	}

	// 6. Attacker must be attacking
	if !gameState.combat.attackers[attackerID] {
		return false, nil
	}

	// 7. Controller of blocker must be opponent of attacker's controller
	// Find the group this attacker is in to get the defending player
	var defendingPlayerID string
	for _, group := range gameState.combat.groups {
		for _, aid := range group.attackers {
			if aid == attackerID {
				defendingPlayerID = group.defendingPlayerID
				break
			}
		}
		if defendingPlayerID != "" {
			break
		}
	}

	if defendingPlayerID == "" {
		return false, fmt.Errorf("attacker %s not found in any combat group", attackerID)
	}

	// Blocker must be controlled by the defending player
	if blocker.ControllerID != defendingPlayerID {
		return false, nil
	}

	// Get attacker for evasion checks
	attacker := gameState.cards[attackerID]

	// Unblockable check: If attacker has "can't be blocked" ability, it cannot be blocked by any creature
	// Per Rule 509.1b and Java CantBeBlockedSourceEffect.canBeBlocked() which returns false
	if e.hasAbilityWithEffects(gameState, attacker, abilityUnblockable) {
		return false, nil
	}

	// Flying restriction: creatures with flying can only be blocked by creatures with flying or reach
	// Exception: Dragons can be blocked by non-flying creatures with special abilities (AsThoughEffectType.BLOCK_DRAGON)
	// Check both base and granted abilities
	if e.hasAbilityWithEffects(gameState, attacker, abilityFlying) {
		if !e.hasAbilityWithEffects(gameState, blocker, abilityFlying) && !e.hasAbilityWithEffects(gameState, blocker, abilityReach) {
			// TODO: Check for AsThoughEffectType.BLOCK_DRAGON and attacker.hasSubtype(SubType.DRAGON)
			// This requires implementing:
			// 1. Subtype checking system
			// 2. AsThough effects / continuous effects system
			return false, nil
		}
	}

	// TODO: Check other restriction effects (shadow, intimidate, etc.)
	// TODO: Check protection

	return true, nil
}

// DeclareBlocker declares a creature as a blocker for an attacker
// Per Java PlayerImpl.declareBlocker()
func (e *MageEngine) DeclareBlocker(gameID, blockerID, attackerID, playerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	// Fire declare blockers step pre event (before first blocker)
	gameState.mu.RLock()
	hasBlockers := len(gameState.combat.blockers) > 0
	gameState.mu.RUnlock()

	if !hasBlockers {
		gameState.mu.Lock()
		gameState.eventBus.Publish(rules.NewEvent(rules.EventDeclareBlockersStepPre, "", "", playerID))
		gameState.mu.Unlock()
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Validate blocker can block this attacker
	canBlock, err := e.canBlockInternal(gameState, blockerID, attackerID)
	if err != nil {
		return err
	}
	if !canBlock {
		return fmt.Errorf("creature %s cannot block attacker %s", blockerID, attackerID)
	}

	// Find the combat group for this attacker
	var group *combatGroup
	for _, g := range gameState.combat.groups {
		for _, aid := range g.attackers {
			if aid == attackerID {
				group = g
				break
			}
		}
		if group != nil {
			break
		}
	}

	if group == nil {
		return fmt.Errorf("attacker %s not found in any combat group", attackerID)
	}

	// Validate player controls the blocker
	blocker, exists := gameState.cards[blockerID]
	if !exists {
		return fmt.Errorf("blocker %s not found", blockerID)
	}

	if blocker.ControllerID != playerID {
		return fmt.Errorf("player %s does not control blocker %s", playerID, blockerID)
	}

	// Check if blocker is already blocking
	if blocker.Blocking {
		// In MTG, a creature can block multiple attackers in some cases
		// For now, we'll allow it but track it properly
		// TODO: Check for restrictions on multiple blocks
	}

	// Add blocker to the group
	group.blockers = append(group.blockers, blockerID)
	group.blocked = true
	gameState.combat.blockers[blockerID] = true

	// Update blocker's blocking status
	blocker.Blocking = true
	if blocker.BlockingWhat == nil {
		blocker.BlockingWhat = []string{}
	}
	blocker.BlockingWhat = append(blocker.BlockingWhat, attackerID)

	// Add to blocking groups map (blocker -> group)
	gameState.combat.blockingGroups[blockerID] = group

	// Fire BLOCKER_DECLARED event
	blockerEvent := rules.Event{
		Type:       rules.EventBlockerDeclared,
		SourceID:   blockerID,
		TargetID:   attackerID,
		PlayerID:   playerID,
		Controller: playerID,
	}
	gameState.eventBus.Publish(blockerEvent)

	// Check for combat triggers (e.g., "Whenever ~ blocks")
	e.checkCombatTriggers(gameState, blockerEvent)

	if e.logger != nil {
		e.logger.Debug("blocker declared",
			zap.String("game_id", gameID),
			zap.String("blocker_id", blockerID),
			zap.String("attacker_id", attackerID),
			zap.String("player_id", playerID),
		)
	}

	return nil
}

// canBlockInternal is an internal version of CanBlock that works with locked state
// RULES-LIGHT: All blocking is allowed - players handle restrictions manually
func (e *MageEngine) canBlockInternal(gameState *engineGameState, blockerID, attackerID string) (bool, error) {
	// Basic existence checks only
	_, blockerExists := gameState.cards[blockerID]
	if !blockerExists {
		return false, fmt.Errorf("blocker %s not found", blockerID)
	}

	_, attackerExists := gameState.cards[attackerID]
	if !attackerExists {
		return false, fmt.Errorf("attacker %s not found", attackerID)
	}

	// RULES-LIGHT: All other restrictions (tapped, flying/reach, unblockable) are removed
	// Players are trusted to follow blocking rules manually
	return true, nil
}

// RemoveBlocker removes a blocker from combat
// Per Java CombatGroup.remove()
func (e *MageEngine) RemoveBlocker(gameID, blockerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Find the combat group this blocker is in
	group, exists := gameState.combat.blockingGroups[blockerID]
	if !exists {
		return fmt.Errorf("blocker %s is not blocking", blockerID)
	}

	// Remove blocker from group
	for i, bid := range group.blockers {
		if bid == blockerID {
			group.blockers = append(group.blockers[:i], group.blockers[i+1:]...)
			break
		}
	}

	// Update blocked status
	if len(group.blockers) == 0 {
		group.blocked = false
	}

	// Remove from blocking groups map
	delete(gameState.combat.blockingGroups, blockerID)

	// Remove from global blockers set
	delete(gameState.combat.blockers, blockerID)

	// Update blocker card state
	blocker, exists := gameState.cards[blockerID]
	if exists {
		blocker.Blocking = false
		blocker.BlockingWhat = nil
	}

	// Fire REMOVED_FROM_COMBAT event (Java: Combat.removeFromCombat)
	gameState.eventBus.Publish(rules.NewEvent(rules.EventRemovedFromCombat, blockerID, "", ""))

	if e.logger != nil {
		e.logger.Debug("blocker removed",
			zap.String("game_id", gameID),
			zap.String("blocker_id", blockerID),
		)
	}

	return nil
}

// RemoveAttacker removes an attacker from combat
// Per Java Combat.removeAttacker()
func (e *MageEngine) RemoveAttacker(gameID, attackerID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Check if creature is actually attacking
	if !gameState.combat.attackers[attackerID] {
		return fmt.Errorf("creature %s is not attacking", attackerID)
	}

	// Find and remove from combat group
	var groupToRemove *combatGroup
	for _, group := range gameState.combat.groups {
		for i, aid := range group.attackers {
			if aid == attackerID {
				// Remove attacker from group
				group.attackers = append(group.attackers[:i], group.attackers[i+1:]...)

				// If group is now empty, mark for removal
				if len(group.attackers) == 0 {
					groupToRemove = group
				}
				break
			}
		}
	}

	// Remove empty group
	if groupToRemove != nil {
		// Move to former groups
		gameState.combat.formerGroups = append(gameState.combat.formerGroups, groupToRemove)

		// Remove from active groups
		for i, g := range gameState.combat.groups {
			if g == groupToRemove {
				gameState.combat.groups = append(gameState.combat.groups[:i], gameState.combat.groups[i+1:]...)
				break
			}
		}
	}

	// Remove from global attackers set
	delete(gameState.combat.attackers, attackerID)

	// Update attacker card state
	attacker, exists := gameState.cards[attackerID]
	if exists {
		attacker.Attacking = false
		attacker.AttackingWhat = ""

		// Untap if it was tapped by attack (Java: attackersTappedByAttack check)
		if gameState.combat.attackersTapped[attackerID] {
			attacker.Tapped = false
			delete(gameState.combat.attackersTapped, attackerID)
		}
	}

	// Fire REMOVED_FROM_COMBAT event (Java: Combat.removeFromCombat)
	gameState.eventBus.Publish(rules.NewEvent(rules.EventRemovedFromCombat, attackerID, "", ""))

	if e.logger != nil {
		e.logger.Debug("attacker removed",
			zap.String("game_id", gameID),
			zap.String("attacker_id", attackerID),
		)
	}

	return nil
}

// RemoveFromCombat removes a creature from combat completely, clearing both attacking and blocking state
// This is a general-purpose removal that handles both attackers and blockers
// Per Java Combat.removeFromCombat() - used for effects like regeneration, control changes, and phasing
func (e *MageEngine) RemoveFromCombat(gameID, creatureID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	creature, exists := gameState.cards[creatureID]
	if !exists {
		return fmt.Errorf("creature %s not found", creatureID)
	}

	removed := false

	// Remove as attacker if attacking
	if gameState.combat.attackers[creatureID] {
		// Find and remove from combat group
		var groupToRemove *combatGroup
		for _, group := range gameState.combat.groups {
			for i, aid := range group.attackers {
				if aid == creatureID {
					// Remove attacker from group
					group.attackers = append(group.attackers[:i], group.attackers[i+1:]...)

					// If group is now empty, mark for removal
					if len(group.attackers) == 0 {
						groupToRemove = group
					}
					removed = true
					break
				}
			}
		}

		// Remove empty group
		if groupToRemove != nil {
			// Move to former groups
			gameState.combat.formerGroups = append(gameState.combat.formerGroups, groupToRemove)

			// Remove from active groups
			for i, g := range gameState.combat.groups {
				if g == groupToRemove {
					gameState.combat.groups = append(gameState.combat.groups[:i], gameState.combat.groups[i+1:]...)
					break
				}
			}
		}

		// Remove from global attackers set
		delete(gameState.combat.attackers, creatureID)

		// Clear attacking state (do NOT untap - that's specific to RemoveAttacker)
		creature.Attacking = false
		creature.AttackingWhat = ""
	}

	// Remove as blocker if blocking
	if gameState.combat.blockers[creatureID] {
		// Find the combat group this blocker is in
		group, exists := gameState.combat.blockingGroups[creatureID]
		if exists {
			// Remove blocker from group
			for i, bid := range group.blockers {
				if bid == creatureID {
					group.blockers = append(group.blockers[:i], group.blockers[i+1:]...)
					break
				}
			}

			// Update blocked status
			if len(group.blockers) == 0 {
				group.blocked = false
			}

			// Remove from blocking groups map
			delete(gameState.combat.blockingGroups, creatureID)
			removed = true
		}

		// Remove from global blockers set
		delete(gameState.combat.blockers, creatureID)

		// Clear blocking state
		creature.Blocking = false
		creature.BlockingWhat = nil
	}

	// Fire REMOVED_FROM_COMBAT event only if creature was actually in combat
	if removed {
		gameState.eventBus.Publish(rules.NewEvent(rules.EventRemovedFromCombat, creatureID, "", ""))

		if e.logger != nil {
			e.logger.Debug("creature removed from combat",
				zap.String("game_id", gameID),
				zap.String("creature_id", creatureID),
			)
		}
	}

	return nil
}

// CheckForRemoveFromCombat checks all attacking and blocking creatures and removes those that are no longer creatures
// Per Java Combat.checkForRemoveFromCombat() - called during combat steps to enforce rule that non-creatures can't attack/block
func (e *MageEngine) CheckForRemoveFromCombat(gameID string) error {
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] CheckForRemoveFromCombat acquiring e.mu.RLock",
			zap.String("game_id", gameID))
	}
	e.mu.RLock()
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] CheckForRemoveFromCombat acquired e.mu.RLock",
			zap.String("game_id", gameID))
	}
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] CheckForRemoveFromCombat released e.mu.RUnlock",
			zap.String("game_id", gameID))
	}

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] CheckForRemoveFromCombat acquiring gameState.mu.Lock",
			zap.String("game_id", gameID))
	}
	gameState.mu.Lock()
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] CheckForRemoveFromCombat acquired gameState.mu.Lock",
			zap.String("game_id", gameID))
	}
	defer func() {
		if e.logger != nil {
			e.logger.Info("[LOCK-DEBUG] CheckForRemoveFromCombat releasing gameState.mu (defer)",
				zap.String("game_id", gameID))
		}
		gameState.mu.Unlock()
	}()

	// Collect all attackers and blockers that need to be removed
	// We collect first, then remove, to avoid modifying maps while iterating
	toRemove := make([]string, 0)

	// Check all attackers
	for creatureID := range gameState.combat.attackers {
		creature, exists := gameState.cards[creatureID]
		if exists && !e.isCreature(creature) {
			toRemove = append(toRemove, creatureID)
		}
	}

	// Check all blockers
	for creatureID := range gameState.combat.blockers {
		creature, exists := gameState.cards[creatureID]
		if exists && !e.isCreature(creature) {
			toRemove = append(toRemove, creatureID)
		}
	}

	// Remove all non-creatures from combat
	// We need to temporarily unlock to call RemoveFromCombat which takes its own lock
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] CheckForRemoveFromCombat releasing gameState.mu for RemoveFromCombat loop",
			zap.String("game_id", gameID),
			zap.Int("creatures_to_remove", len(toRemove)))
	}
	gameState.mu.Unlock()
	for _, creatureID := range toRemove {
		if err := e.RemoveFromCombat(gameID, creatureID); err != nil && e.logger != nil {
			e.logger.Warn("failed to remove non-creature from combat",
				zap.String("game_id", gameID),
				zap.String("creature_id", creatureID),
				zap.Error(err),
			)
		}
	}
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] CheckForRemoveFromCombat re-acquiring gameState.mu after RemoveFromCombat loop",
			zap.String("game_id", gameID))
	}
	gameState.mu.Lock()
	if e.logger != nil {
		e.logger.Info("[LOCK-DEBUG] CheckForRemoveFromCombat re-acquired gameState.mu",
			zap.String("game_id", gameID))
	}

	return nil
}

// OrderBlockers sets the damage assignment order for blockers on a specific attacker
// Per Java: The attacking player chooses the order in which damage is assigned to blockers
// This is typically done during the declare blockers step, before damage assignment
func (e *MageEngine) OrderBlockers(gameID, attackerID string, blockerOrder []string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Find the combat group for this attacker
	var targetGroup *combatGroup
	for _, group := range gameState.combat.groups {
		for _, aid := range group.attackers {
			if aid == attackerID {
				targetGroup = group
				break
			}
		}
		if targetGroup != nil {
			break
		}
	}

	if targetGroup == nil {
		return fmt.Errorf("attacker %s not found in combat", attackerID)
	}

	// Validate that all blockers in the order are actually blocking this attacker
	if len(blockerOrder) != len(targetGroup.blockers) {
		return fmt.Errorf("blocker order length (%d) does not match actual blocker count (%d)",
			len(blockerOrder), len(targetGroup.blockers))
	}

	// Create a set of current blockers for validation
	currentBlockers := make(map[string]bool)
	for _, bid := range targetGroup.blockers {
		currentBlockers[bid] = true
	}

	// Validate that all provided blockers are actually blocking
	for _, bid := range blockerOrder {
		if !currentBlockers[bid] {
			return fmt.Errorf("blocker %s is not blocking attacker %s", bid, attackerID)
		}
	}

	// Update the blocker order
	targetGroup.blockers = blockerOrder

	if e.logger != nil {
		e.logger.Debug("blocker order set",
			zap.String("game_id", gameID),
			zap.String("attacker_id", attackerID),
			zap.Strings("blocker_order", blockerOrder),
		)
	}

	return nil
}

// AcceptBlockers finalizes the blocker declarations and fires events
// Per Java Combat.acceptBlockers()
func (e *MageEngine) AcceptBlockers(gameID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Validate menace and other blocking restrictions
	// Per Java CombatGroup.acceptBlockers() lines 710-718
	for _, group := range gameState.combat.groups {
		if len(group.attackers) == 0 {
			continue
		}

		for _, attackerID := range group.attackers {
			attacker, exists := gameState.cards[attackerID]
			if !exists {
				continue
			}

			minBlockedBy := e.getMinBlockedBy(attacker)
			if minBlockedBy > 1 && len(group.blockers) > 0 && len(group.blockers) < minBlockedBy {
				// Menace violation - remove all blockers from this attacker
				if e.logger != nil {
					e.logger.Debug("menace violation - removing blockers",
						zap.String("attacker_id", attackerID),
						zap.Int("blockers", len(group.blockers)),
						zap.Int("min_required", minBlockedBy),
					)
				}

				// Remove all blockers from this group
				for _, blockerID := range group.blockers {
					delete(gameState.combat.blockers, blockerID)
					if blocker, exists := gameState.cards[blockerID]; exists {
						blocker.Blocking = false
						blocker.BlockingWhat = []string{}
					}
				}
				group.blockers = []string{}
				group.blocked = false
			}
		}
	}

	// Fire BLOCKER_DECLARED events for each blocker-attacker pair
	// Per Java CombatGroup.acceptBlockers()
	for _, group := range gameState.combat.groups {
		if len(group.attackers) == 0 {
			continue
		}

		for _, blockerID := range group.blockers {
			blocker, exists := gameState.cards[blockerID]
			if !exists {
				continue
			}

			for _, attackerID := range group.attackers {
				gameState.eventBus.Publish(rules.Event{
					Type:       rules.EventBlockerDeclared,
					SourceID:   blockerID,
					TargetID:   attackerID,
					PlayerID:   blocker.ControllerID,
					Controller: blocker.ControllerID,
				})
			}
		}

		// Fire CREATURE_BLOCKED event for each attacker that is blocked
		if len(group.blockers) > 0 {
			for _, attackerID := range group.attackers {
				blockedEvent := rules.Event{
					Type:     rules.EventCreatureBlocked,
					SourceID: attackerID,
				}
				gameState.eventBus.Publish(blockedEvent)
				// Check for combat triggers (e.g., "Whenever ~ becomes blocked")
				e.checkCombatTriggers(gameState, blockedEvent)
			}
		}
	}

	// Fire CREATURE_BLOCKS event for each blocker
	// Per Java Combat.acceptBlockers()
	for blockerID := range gameState.combat.blockers {
		blocksEvent := rules.Event{
			Type:     rules.EventCreatureBlocks,
			SourceID: blockerID,
		}
		gameState.eventBus.Publish(blocksEvent)
		// Check for combat triggers (e.g., "Whenever ~ blocks")
		e.checkCombatTriggers(gameState, blocksEvent)
	}

	// Fire DECLARED_BLOCKERS event for each defending player
	defendingPlayers := make(map[string]bool)
	for _, group := range gameState.combat.groups {
		defendingPlayers[group.defendingPlayerID] = true
	}

	for playerID := range defendingPlayers {
		gameState.eventBus.Publish(rules.Event{
			Type:       rules.EventDeclaredBlockers,
			PlayerID:   playerID,
			Controller: playerID,
		})
	}

	// Fire UNBLOCKED_ATTACKER event for each unblocked attacker
	// Per Java Combat.acceptBlockers() - fires after blockers are declared
	for _, group := range gameState.combat.groups {
		if len(group.attackers) > 0 && !group.blocked {
			for _, attackerID := range group.attackers {
				gameState.eventBus.Publish(rules.Event{
					Type:       rules.EventUnblockedAttacker,
					SourceID:   attackerID,
					PlayerID:   gameState.combat.attackingPlayerID,
					Controller: gameState.combat.attackingPlayerID,
				})
			}
		}
	}

	if e.logger != nil {
		e.logger.Debug("blockers accepted",
			zap.String("game_id", gameID),
			zap.Int("blocker_count", len(gameState.combat.blockers)),
		)
	}

	return nil
}

// CheckBlockRequirements validates that all blocking requirements are met
// Per Java Combat.checkBlockRequirements()
// Returns list of violations (empty if all requirements met)
func (e *MageEngine) CheckBlockRequirements(gameID, playerID string) ([]string, error) {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	violations := make([]string, 0)

	// Check creatures that must block specific attackers
	for blockerID, requiredAttackers := range gameState.combat.creatureMustBlockAttackers {
		blocker, exists := gameState.cards[blockerID]
		if !exists || blocker.ControllerID != playerID {
			continue
		}

		// Check if blocker is actually blocking
		if !gameState.combat.blockers[blockerID] {
			// TODO: Check if blocker CAN block (not tapped, etc.)
			// For now, just report the violation
			if len(requiredAttackers) > 0 {
				for attackerID := range requiredAttackers {
					violations = append(violations, fmt.Sprintf("creature %s must block attacker %s", blockerID, attackerID))
				}
			} else {
				violations = append(violations, fmt.Sprintf("creature %s must block if able", blockerID))
			}
		}
	}

	// Check minimum blockers per attacker (e.g., menace)
	for attackerID, minBlockers := range gameState.combat.minBlockersPerAttacker {
		// Count blockers for this attacker
		blockerCount := 0
		for _, group := range gameState.combat.groups {
			for _, aid := range group.attackers {
				if aid == attackerID {
					blockerCount = len(group.blockers)
					break
				}
			}
		}

		if blockerCount > 0 && blockerCount < minBlockers {
			violations = append(violations,
				fmt.Sprintf("attacker %s requires at least %d blockers, but has %d",
					attackerID, minBlockers, blockerCount))
		}
	}

	return violations, nil
}

// CheckBlockRestrictions validates that no blocking restrictions are violated
// Per Java Combat.checkBlockRestrictions()
// Returns list of violations (empty if no restrictions violated)
func (e *MageEngine) CheckBlockRestrictions(gameID, playerID string) ([]string, error) {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	violations := make([]string, 0)

	// Check maximum blockers per attacker
	for attackerID, maxBlockers := range gameState.combat.maxBlockersPerAttacker {
		// Count blockers for this attacker
		blockerCount := 0
		for _, group := range gameState.combat.groups {
			for _, aid := range group.attackers {
				if aid == attackerID {
					blockerCount = len(group.blockers)
					break
				}
			}
		}

		if blockerCount > maxBlockers {
			violations = append(violations,
				fmt.Sprintf("attacker %s can be blocked by at most %d creatures, but has %d",
					attackerID, maxBlockers, blockerCount))
		}
	}

	// TODO: Check other restrictions from continuous effects
	// - "can't block" effects
	// - "can't block creature X" effects
	// - Protection from color/type restrictions

	return violations, nil
}

// ValidateAttackerCount checks if the number of attackers meets requirements
// Per Java Combat validation
func (e *MageEngine) ValidateAttackerCount(gameID string) ([]string, error) {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	violations := make([]string, 0)

	// Check maximum attackers
	if gameState.combat.maxAttackers >= 0 {
		attackerCount := len(gameState.combat.attackers)
		if attackerCount > gameState.combat.maxAttackers {
			violations = append(violations,
				fmt.Sprintf("maximum %d attackers allowed, but %d are attacking",
					gameState.combat.maxAttackers, attackerCount))
		}
	}

	// Check forced attackers
	for creatureID, requiredDefenders := range gameState.combat.creaturesForcedToAttack {
		creature, exists := gameState.cards[creatureID]
		if !exists {
			continue
		}

		// Check if creature is attacking
		if !gameState.combat.attackers[creatureID] {
			// TODO: Check if creature CAN attack (not tapped, summoning sickness, etc.)
			// For now, just report the violation
			if len(requiredDefenders) > 0 {
				violations = append(violations,
					fmt.Sprintf("creature %s must attack one of: %v", creatureID, requiredDefenders))
			} else {
				violations = append(violations, fmt.Sprintf("creature %s must attack if able", creatureID))
			}
		} else if len(requiredDefenders) > 0 {
			// Check if attacking the correct defender
			attackingCorrectDefender := false
			for _, group := range gameState.combat.groups {
				for _, aid := range group.attackers {
					if aid == creatureID && requiredDefenders[group.defenderID] {
						attackingCorrectDefender = true
						break
					}
				}
			}

			if !attackingCorrectDefender {
				violations = append(violations,
					fmt.Sprintf("creature %s must attack one of: %v", creatureID, requiredDefenders))
			}
		}

		_ = creature // avoid unused warning
	}

	return violations, nil
}

// AssignCombatDamage assigns combat damage for all combat groups
// Per Java CombatDamageStep.beginStep()
func (e *MageEngine) AssignCombatDamage(gameID string, firstStrike bool) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	return e.assignCombatDamageInternal(gameState, firstStrike)
}

// assignCombatDamageInternal assigns combat damage without acquiring locks
// Caller must hold gameState.mu.Lock()
func (e *MageEngine) assignCombatDamageInternal(gameState *engineGameState, firstStrike bool) error {
	// Fire combat damage step pre event
	gameState.eventBus.Publish(rules.NewEvent(rules.EventCombatDamageStepPre, "", "", ""))

	// Assign damage to blockers (attackers dealing damage)
	for _, group := range gameState.combat.groups {
		if err := e.assignDamageToBlockers(gameState, group, firstStrike); err != nil {
			return err
		}
	}

	// Assign damage to attackers (blockers dealing damage)
	// Track which blockers have already dealt damage (since a blocker can be in multiple groups)
	processedBlockers := make(map[string]bool)
	for _, group := range gameState.combat.groups {
		if len(group.blockers) > 0 {
			if err := e.assignDamageToAttackers(gameState, group, firstStrike, processedBlockers); err != nil {
				return err
			}
		}
	}

	// Fire combat damage assigned event
	gameState.eventBus.Publish(rules.NewEvent(rules.EventCombatDamageAssigned, "", "", ""))

	if e.logger != nil {
		e.logger.Debug("combat damage assigned",
			zap.String("game_id", gameState.gameID),
			zap.Bool("first_strike", firstStrike),
		)
	}

	return nil
}

// assignDamageToBlockers handles attacker damage to blockers or defender
// Per Java CombatGroup.assignDamageToBlockers()
func (e *MageEngine) assignDamageToBlockers(gameState *engineGameState, group *combatGroup, firstStrike bool) error {
	if len(group.attackers) == 0 {
		return nil
	}

	// Get the attacker (should only be one per group)
	attackerID := group.attackers[0]
	attacker, exists := gameState.cards[attackerID]
	if !exists {
		return nil
	}

	// Check if attacker deals damage this step (first strike check)
	if !e.dealsDamageThisStep(gameState, attacker, firstStrike) {
		return nil
	}

	// Record first striker if dealing damage in first strike step
	if firstStrike && e.hasFirstOrDoubleStrikeWithEffects(gameState, attacker) {
		e.recordFirstStrikingCreature(gameState, attackerID)
	}

	// Get attacker's power
	power, err := e.getCreaturePower(gameState, attacker)
	if err != nil {
		power = 0
	}

	hasTrample := e.hasAbility(attacker, abilityTrample)

	// Check if there are any live blockers
	liveBlockers := 0
	for _, blockerID := range group.blockers {
		if blocker, exists := gameState.cards[blockerID]; exists && blocker.Zone == zoneBattlefield {
			liveBlockers++
		}
	}

	if len(group.blockers) == 0 {
		// Never blocked - deal damage to defender
		return e.dealDamageToDefender(gameState, attacker, group.defenderID, power)
	}

	if liveBlockers == 0 {
		// Was blocked but all blockers are dead (e.g., from first strike)
		// With trample, remaining damage goes through
		// Without trample, no damage goes through
		if hasTrample {
			return e.dealDamageToDefender(gameState, attacker, group.defenderID, power)
		}
		return nil
	}

	// Blocked - assign damage to blockers
	// Rule 510.1c: Player divides damage as they choose among blockers
	// Use stored damage assignment or compute default

	var damageAssignment map[string]int
	if assignment, exists := group.attackerDamageAssignments[attackerID]; exists {
		damageAssignment = assignment
	} else {
		// No explicit assignment - use default
		damageAssignment = e.computeDefaultAttackerDamageAssignment(gameState, attackerID, group.blockers)
	}

	// Apply the damage assignment
	totalAssigned := 0
	for blockerID, damage := range damageAssignment {
		blocker, exists := gameState.cards[blockerID]
		if !exists || blocker.Zone != zoneBattlefield {
			continue
		}

		if damage > 0 {
			e.markDamageWithLifelink(gameState, blocker, damage, attackerID)
			totalAssigned += damage
		}
	}

	// With trample, excess damage goes to defender
	if hasTrample {
		trampleDamage := power - totalAssigned
		if trampleDamage > 0 {
			return e.dealDamageToDefender(gameState, attacker, group.defenderID, trampleDamage)
		}
	}

	return nil
}

// assignDamageToAttackers handles blocker damage to attackers
// Per Java CombatGroup.assignDamageToAttackers()
func (e *MageEngine) assignDamageToAttackers(gameState *engineGameState, group *combatGroup, firstStrike bool, processedBlockers map[string]bool) error {
	if len(group.blockers) == 0 {
		return nil
	}

	// For each blocker, deal damage to the attacker(s) it's blocking
	for _, blockerID := range group.blockers {
		// Skip if we've already processed this blocker (can be in multiple groups)
		if processedBlockers[blockerID] {
			continue
		}

		blocker, exists := gameState.cards[blockerID]
		if !exists {
			continue
		}

		// Dead creatures don't deal damage
		if blocker.Zone != zoneBattlefield {
			continue
		}

		// Check if blocker deals damage this step
		if !e.dealsDamageThisStep(gameState, blocker, firstStrike) {
			continue
		}

		// Record first striker if dealing damage in first strike step
		if firstStrike && e.hasFirstOrDoubleStrikeWithEffects(gameState, blocker) {
			e.recordFirstStrikingCreature(gameState, blockerID)
		}

		// Rule 510.1d: Blocker divides damage as controller chooses among attackers
		// If blocker blocks multiple attackers, get assignment from any group (should be same in all)
		var damageAssignment map[string]int
		if assignment, exists := group.blockerDamageAssignments[blockerID]; exists {
			damageAssignment = assignment
		} else {
			// No explicit assignment - collect all attackers this blocker is blocking
			blockingAttackers := make([]string, 0)
			for _, g := range gameState.combat.groups {
				for _, bid := range g.blockers {
					if bid == blockerID {
						blockingAttackers = append(blockingAttackers, g.attackers...)
					}
				}
			}
			damageAssignment = e.computeDefaultBlockerDamageAssignment(gameState, blockerID, blockingAttackers)
		}

		// Apply the damage assignment
		for attackerID, damage := range damageAssignment {
			attacker, exists := gameState.cards[attackerID]
			if !exists || attacker.Zone != zoneBattlefield {
				continue
			}

			if damage > 0 {
				e.markDamageWithLifelink(gameState, attacker, damage, blockerID)
			}
		}

		// Mark this blocker as processed
		processedBlockers[blockerID] = true
	}

	return nil
}

// ApplyCombatDamage applies all marked damage
// Per Java CombatGroup.applyDamage()
func (e *MageEngine) ApplyCombatDamage(gameID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	return e.applyCombatDamageInternal(gameState)
}

// applyCombatDamageInternal applies combat damage without acquiring locks
// Caller must hold gameState.mu.Lock()
func (e *MageEngine) applyCombatDamageInternal(gameState *engineGameState) error {
	// Apply damage to all creatures in combat
	for _, group := range gameState.combat.groups {
		// Apply damage to attackers
		for _, attackerID := range group.attackers {
			if err := e.applyDamageToCreature(gameState, attackerID); err != nil {
				return err
			}
		}

		// Apply damage to blockers
		for _, blockerID := range group.blockers {
			if err := e.applyDamageToCreature(gameState, blockerID); err != nil {
				return err
			}
		}
	}

	if e.logger != nil {
		e.logger.Debug("combat damage applied",
			zap.String("game_id", gameState.gameID),
		)
	}

	// Fire combat damage applied event
	gameState.eventBus.Publish(rules.NewEvent(rules.EventCombatDamageApplied, "", "", ""))

	return nil
}

// AssignAttackerDamage assigns how an attacker divides its damage among blockers
// Rule 510.1c: A blocked creature assigns its combat damage divided as its controller chooses among blockers
// Rule 702.22j: When blocked by banding creature, DEFENDING player assigns (not attacking player)
// Per Java CombatGroup.blockerDamage() multi-amount dialog
func (e *MageEngine) AssignAttackerDamage(gameID, attackerID, playerID string, damageMap map[string]int) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Find the combat group for this attacker
	var group *combatGroup
	for _, g := range gameState.combat.groups {
		for _, aID := range g.attackers {
			if aID == attackerID {
				group = g
				break
			}
		}
		if group != nil {
			break
		}
	}

	if group == nil {
		return fmt.Errorf("attacker %s not found in combat", attackerID)
	}

	// Rule 702.22j: When blocked by banding creature, defending player controls damage assignment
	if e.defenderControlsDamageAssignment(gameState, group) {
		// Defending player must assign damage, not attacking player
		if playerID != group.defendingPlayerID {
			return fmt.Errorf("defending player must assign damage (blocked by banding creature)")
		}
	} else {
		// Normal case: attacking player assigns damage
		if playerID != gameState.combat.attackingPlayerID {
			return fmt.Errorf("attacking player must assign damage")
		}
	}

	// Validate the damage assignment
	attacker, exists := gameState.cards[attackerID]
	if !exists {
		return fmt.Errorf("attacker %s not found", attackerID)
	}

	power, err := e.getCreaturePower(gameState, attacker)
	if err != nil {
		power = 0
	}

	// Validate total damage equals power
	totalAssigned := 0
	for _, damage := range damageMap {
		totalAssigned += damage
	}

	// With trample, player can assign less than full power (excess tramples through)
	hasTrample := e.hasAbility(attacker, abilityTrample)
	if hasTrample {
		if totalAssigned > power {
			return fmt.Errorf("cannot assign more damage (%d) than creature's power (%d)", totalAssigned, power)
		}
	} else {
		if totalAssigned != power {
			return fmt.Errorf("must assign all damage (%d) to blockers, assigned %d", power, totalAssigned)
		}
	}

	// Validate all targets are valid blockers
	for blockerID := range damageMap {
		found := false
		for _, bid := range group.blockers {
			if bid == blockerID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("creature %s is not blocking this attacker", blockerID)
		}
	}

	// Store the assignment
	group.attackerDamageAssignments[attackerID] = damageMap

	if e.logger != nil {
		e.logger.Debug("attacker damage assigned",
			zap.String("attacker_id", attackerID),
			zap.Any("damage_map", damageMap),
		)
	}

	return nil
}

// AssignBlockerDamage assigns how a blocker divides its damage among attackers it's blocking
// Rule 510.1d: A blocking creature assigns its combat damage divided as its controller chooses among attackers
// Rule 702.22k: When blocking banding attacker, ATTACKING player assigns (not defending player)
// Per Java CombatGroup.attackerDamage() multi-amount dialog
func (e *MageEngine) AssignBlockerDamage(gameID, blockerID, playerID string, damageMap map[string]int) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Validate the blocker exists
	blocker, exists := gameState.cards[blockerID]
	if !exists {
		return fmt.Errorf("blocker %s not found", blockerID)
	}

	// Find a combat group this blocker is in (to check for banding attackers)
	var checkGroup *combatGroup
	for _, group := range gameState.combat.groups {
		for _, bid := range group.blockers {
			if bid == blockerID {
				checkGroup = group
				break
			}
		}
		if checkGroup != nil {
			break
		}
	}

	// Rule 702.22k: When blocking banding attacker, attacking player controls damage assignment
	if checkGroup != nil && e.attackerControlsDamageAssignment(gameState, checkGroup) {
		// Attacking player must assign damage, not defending player
		if playerID != gameState.combat.attackingPlayerID {
			return fmt.Errorf("attacking player must assign damage (blocking banding attacker)")
		}
	} else {
		// Normal case: defending player assigns damage
		if playerID != blocker.ControllerID {
			return fmt.Errorf("blocker's controller must assign damage")
		}
	}

	power, err := e.getCreaturePower(gameState, blocker)
	if err != nil {
		power = 0
	}

	// Validate total damage equals power
	totalAssigned := 0
	for _, damage := range damageMap {
		totalAssigned += damage
	}

	if totalAssigned != power {
		return fmt.Errorf("must assign all damage (%d), assigned %d", power, totalAssigned)
	}

	// Collect all attackers this blocker is blocking (may be in multiple groups)
	blockingAttackers := make(map[string]bool)
	for _, group := range gameState.combat.groups {
		for _, bid := range group.blockers {
			if bid == blockerID {
				for _, aid := range group.attackers {
					blockingAttackers[aid] = true
				}
			}
		}
	}

	// Validate all targets are valid attackers being blocked
	for attackerID := range damageMap {
		if !blockingAttackers[attackerID] {
			return fmt.Errorf("creature %s is not being blocked by this blocker", attackerID)
		}
	}

	// Store the assignment in all groups this blocker is in
	for _, group := range gameState.combat.groups {
		for _, bid := range group.blockers {
			if bid == blockerID {
				group.blockerDamageAssignments[blockerID] = damageMap
				break
			}
		}
	}

	if e.logger != nil {
		e.logger.Debug("blocker damage assigned",
			zap.String("blocker_id", blockerID),
			zap.Any("damage_map", damageMap),
		)
	}

	return nil
}

// computeDefaultAttackerDamageAssignment computes the default damage assignment for an attacker
// Rule 510.1c: With trample, assign lethal to each blocker (excess tramples); without, divide among blockers
// Per Java CombatGroup.blockerDamage() defaultDamage calculation
func (e *MageEngine) computeDefaultAttackerDamageAssignment(gameState *engineGameState, attackerID string, blockers []string) map[string]int {
	attacker, exists := gameState.cards[attackerID]
	if !exists {
		return make(map[string]int)
	}

	power, err := e.getCreaturePower(gameState, attacker)
	if err != nil {
		power = 0
	}

	hasTrample := e.hasAbility(attacker, abilityTrample)
	assignment := make(map[string]int)

	if hasTrample {
		// With trample: assign lethal damage to each blocker in order
		remainingDamage := power
		for _, blockerID := range blockers {
			blocker, exists := gameState.cards[blockerID]
			if !exists || blocker.Zone != zoneBattlefield {
				continue
			}

			lethalDamage := e.getLethalDamageWithAttacker(gameState, blocker, attackerID)
			damageToAssign := lethalDamage
			if damageToAssign > remainingDamage {
				damageToAssign = remainingDamage
			}

			if damageToAssign > 0 {
				assignment[blockerID] = damageToAssign
				remainingDamage -= damageToAssign
			}

			if remainingDamage <= 0 {
				break
			}
		}
		// Remaining damage tramples through (not part of blocker assignment)
	} else {
		// Without trample: divide damage evenly among blockers
		if len(blockers) == 0 {
			return assignment
		}

		damagePerBlocker := power / len(blockers)
		remainingDamage := power % len(blockers)

		for i, blockerID := range blockers {
			_, exists := gameState.cards[blockerID]
			if !exists {
				continue
			}

			damage := damagePerBlocker
			if i == 0 {
				damage += remainingDamage // Give remainder to first blocker
			}

			if damage > 0 {
				assignment[blockerID] = damage
			}
		}
	}

	return assignment
}

// computeDefaultBlockerDamageAssignment computes the default damage assignment for a blocker
// Rule 510.1d: Blocker assigns damage divided among attackers it's blocking
// Per Java CombatGroup.attackerDamage() defaultDamage calculation
func (e *MageEngine) computeDefaultBlockerDamageAssignment(gameState *engineGameState, blockerID string, attackers []string) map[string]int {
	blocker, exists := gameState.cards[blockerID]
	if !exists {
		return make(map[string]int)
	}

	power, err := e.getCreaturePower(gameState, blocker)
	if err != nil {
		power = 0
	}

	assignment := make(map[string]int)

	if len(attackers) == 0 {
		return assignment
	}

	// Divide damage evenly among attackers
	damagePerAttacker := power / len(attackers)
	remainingDamage := power % len(attackers)

	for i, attackerID := range attackers {
		_, exists := gameState.cards[attackerID]
		if !exists {
			continue
		}

		damage := damagePerAttacker
		if i == 0 {
			damage += remainingDamage // Give remainder to first attacker
		}

		if damage > 0 {
			assignment[attackerID] = damage
		}
	}

	return assignment
}

// EndCombat ends combat phase, clearing combat flags and moving to former groups
// Per Java Combat.endCombat()
func (e *MageEngine) EndCombat(gameID string) error {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	return e.endCombatInternal(gameState)
}

// endCombatInternal ends combat phase without acquiring locks
// Caller must hold gameState.mu.Lock()
func (e *MageEngine) endCombatInternal(gameState *engineGameState) error {
	// Fire end combat step pre event
	gameState.eventBus.Publish(rules.NewEvent(rules.EventEndCombatStepPre, "", "", ""))

	// Clear combat flags on all creatures in combat
	for _, group := range gameState.combat.groups {
		// Clear attacker flags
		for _, attackerID := range group.attackers {
			if creature, exists := gameState.cards[attackerID]; exists {
				creature.Attacking = false
				creature.Blocking = false
				creature.AttackingWhat = ""
				creature.BlockingWhat = nil
				// Clear damage tracking
				creature.Damage = 0
				creature.DamageSources = nil
			}
		}

		// Clear blocker flags
		for _, blockerID := range group.blockers {
			if creature, exists := gameState.cards[blockerID]; exists {
				creature.Attacking = false
				creature.Blocking = false
				creature.AttackingWhat = ""
				creature.BlockingWhat = nil
				// Clear damage tracking
				creature.Damage = 0
				creature.DamageSources = nil
			}
		}

		// Move attackers to formerAttackers for "attacked this turn" queries
		group.formerAttackers = append([]string{}, group.attackers...)
	}

	// Move current groups to former groups (for historical queries)
	gameState.combat.formerGroups = append([]*combatGroup{}, gameState.combat.groups...)

	// Clear current combat state
	gameState.combat.groups = nil
	gameState.combat.blockingGroups = make(map[string]*combatGroup)
	gameState.combat.attackers = make(map[string]bool)
	gameState.combat.blockers = make(map[string]bool)
	gameState.combat.attackersTapped = make(map[string]bool)
	// Keep defenders for queries
	// Keep attackingPlayerID for queries

	// Cleanup continuous effects that expire at end of combat
	// Per Java: ContinuousEffects.removeEndOfCombatEffects()
	if gameState.layerSystem != nil {
		effects.CleanupEndOfCombatEffects(gameState.layerSystem)
	}

	// Fire end combat event
	gameState.eventBus.Publish(rules.Event{
		Type: rules.EventEndCombatStep,
	})

	if e.logger != nil {
		e.logger.Debug("ended combat",
			zap.String("game_id", gameState.gameID),
			zap.Int("former_groups", len(gameState.combat.formerGroups)),
		)
	}

	return nil
}

// GetAttackedThisTurn returns whether a creature attacked this turn
// Checks formerGroups for historical attack data
func (e *MageEngine) GetAttackedThisTurn(gameID, creatureID string) (bool, error) {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	// Check current combat
	if gameState.combat.attackers[creatureID] {
		return true, nil
	}

	// Check former groups
	for _, group := range gameState.combat.formerGroups {
		for _, attackerID := range group.formerAttackers {
			if attackerID == creatureID {
				return true, nil
			}
		}
	}

	return false, nil
}

// HasFirstOrDoubleStrike returns whether any creature in combat has first strike or double strike
// Per Java Combat.hasFirstOrDoubleStrike()
func (e *MageEngine) HasFirstOrDoubleStrike(gameID string) (bool, error) {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	return e.hasFirstOrDoubleStrikeInternal(gameState), nil
}

// hasFirstOrDoubleStrikeInternal checks if any creature in combat has first or double strike
// This internal version does NOT acquire locks - caller must hold gameState.mu
func (e *MageEngine) hasFirstOrDoubleStrikeInternal(gameState *engineGameState) bool {
	// Check all creatures in combat groups
	for _, group := range gameState.combat.groups {
		// Check attackers
		for _, attackerID := range group.attackers {
			if attacker, exists := gameState.cards[attackerID]; exists {
				if e.hasFirstOrDoubleStrikeWithEffects(gameState, attacker) {
					return true
				}
			}
		}

		// Check blockers
		for _, blockerID := range group.blockers {
			if blocker, exists := gameState.cards[blockerID]; exists {
				if e.hasFirstOrDoubleStrikeWithEffects(gameState, blocker) {
					return true
				}
			}
		}
	}

	return false
}

// Helper methods

// dealsDamageThisStep checks if a creature deals damage this combat damage step
// Per Java CombatGroup.dealsDamageThisStep()
func (e *MageEngine) dealsDamageThisStep(gameState *engineGameState, creature *internalCard, firstStrike bool) bool {
	if creature == nil {
		return false
	}

	if firstStrike {
		// In first strike step, only creatures with first strike or double strike deal damage
		if e.hasFirstOrDoubleStrike(creature) {
			// Record that this creature dealt damage in first strike step
			// (This is done in assignDamageToBlockers/assignDamageToAttackers)
			return true
		}
		return false
	} else {
		// In normal damage step:
		// - Creatures with double strike deal damage again
		// - Creatures without first/double strike deal damage for the first time
		// - Creatures that already dealt damage in first strike step don't deal damage again (unless double strike)
		return e.hasDoubleStrike(creature) || !e.wasFirstStrikingCreatureInCombat(gameState, creature.ID)
	}
}

// hasAbility checks if a creature has a specific ability by ID
// hasAbility checks if a card has a specific ability
// TODO: This should also check for abilities granted by continuous effects
// Java equivalent: permanent.getAbilities(game).containsKey(abilityId)
// Requires implementing:
// 1. ContinuousEffects system to track temporary ability grants
// 2. Layer system for effect ordering (Layer 6 for abilities)
// 3. Effect duration tracking (until end of turn, until end of combat, etc.)
func (e *MageEngine) isCreature(card *internalCard) bool {
	if card == nil {
		return false
	}
	return strings.Contains(strings.ToLower(card.Type), "creature")
}

func (e *MageEngine) isPlaneswalker(card *internalCard) bool {
	if card == nil {
		return false
	}
	return strings.Contains(strings.ToLower(card.Type), "planeswalker")
}

// hasBanding checks if a creature has the banding ability
// Per Java CombatGroup.hasBanding()
func (e *MageEngine) hasBanding(card *internalCard) bool {
	if card == nil {
		return false
	}
	return e.hasAbility(card, abilityBanding)
}

// HasPlayerAttackedPlayerOrPlaneswalker checks if a player attacked another player or their planeswalkers this turn
// Per Java PlayersAttackedThisTurnWatcher.hasPlayerAttackedPlayerOrControlledPlaneswalker()
// Returns true if attacker attacked defender directly OR attacked a planeswalker controlled by defender
func (e *MageEngine) HasPlayerAttackedPlayerOrPlaneswalker(gameID, attackingPlayerID, defendingPlayerID string) (bool, error) {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	// Check if attacker attacked the player directly
	if playerSet, ok := gameState.combat.playersAttackedThisTurn[attackingPlayerID]; ok {
		if playerSet[defendingPlayerID] {
			return true, nil
		}
	}

	// Check if attacker attacked a planeswalker controlled by the defender
	if controllerSet, ok := gameState.combat.planeswalkerControllersAttackedThisTurn[attackingPlayerID]; ok {
		if controllerSet[defendingPlayerID] {
			return true, nil
		}
	}

	return false, nil
}

// defenderControlsDamageAssignment checks if the defending player controls damage assignment
// Rule 702.22j: When blocked by creature with banding, defending player assigns attacker's damage
// Per Java CombatGroup.defenderAssignsCombatDamage()
func (e *MageEngine) defenderControlsDamageAssignment(gameState *engineGameState, group *combatGroup) bool {
	// Check if any blocker has banding
	for _, blockerID := range group.blockers {
		if blocker, exists := gameState.cards[blockerID]; exists {
			if e.hasBanding(blocker) {
				return true
			}
		}
	}
	return false
}

// attackerControlsDamageAssignment checks if the attacking player controls damage assignment
// Rule 702.22k: When blocking creature with banding, attacking player assigns blocker's damage
// Per Java CombatGroup.attackerAssignsCombatDamage()
func (e *MageEngine) attackerControlsDamageAssignment(gameState *engineGameState, group *combatGroup) bool {
	// Check if any attacker has banding
	for _, attackerID := range group.attackers {
		if attacker, exists := gameState.cards[attackerID]; exists {
			if e.hasBanding(attacker) {
				return true
			}
		}
	}
	return false
}

func (e *MageEngine) hasAbility(creature *internalCard, abilityID string) bool {
	if creature == nil {
		return false
	}

	// Check base abilities
	for _, ability := range creature.Abilities {
		if ability.ID == abilityID {
			return true
		}
	}

	return false
}

// hasAbilityWithEffects checks if a creature has a specific ability, including granted abilities
// This version checks both base abilities and abilities granted by continuous effects
// Per Java: permanent.getAbilities(game).containsKey(abilityId)
func (e *MageEngine) hasAbilityWithEffects(gameState *engineGameState, creature *internalCard, abilityID string) bool {
	if creature == nil {
		return false
	}

	// Check base abilities
	if e.hasAbility(creature, abilityID) {
		return true
	}

	// Check continuous effects for granted abilities
	// Per Java: GainAbilityTargetEffect in Layer 6
	if gameState != nil && gameState.layerSystem != nil {
		hasGranted := gameState.layerSystem.HasEffectType(creature.ID, func(effect effects.ContinuousEffect) bool {
			if grantEffect, ok := effect.(*effects.GrantAbilityEffect); ok {
				return grantEffect.GetAbilityID() == abilityID
			}
			return false
		})
		if hasGranted {
			return true
		}
	}

	return false
}

// hasCantAttackEffect checks if a creature is affected by a "can't attack" continuous effect
// Per Java: RestrictionEffect.applies() for attack restrictions
func (e *MageEngine) hasCantAttackEffect(gameState *engineGameState, creatureID string) bool {
	if gameState == nil || gameState.layerSystem == nil {
		return false
	}

	// Check if any CantAttackEffect applies to this creature
	return gameState.layerSystem.HasEffectType(creatureID, func(effect effects.ContinuousEffect) bool {
		_, isCantAttack := effect.(*effects.CantAttackEffect)
		return isCantAttack
	})
}

// hasCantBlockEffect checks if a creature is affected by a "can't block" continuous effect
// Per Java: RestrictionEffect.applies() for block restrictions
func (e *MageEngine) hasCantBlockEffect(gameState *engineGameState, creatureID string) bool {
	if gameState == nil || gameState.layerSystem == nil {
		return false
	}

	// Check if any CantBlockEffect applies to this creature
	return gameState.layerSystem.HasEffectType(creatureID, func(effect effects.ContinuousEffect) bool {
		_, isCantBlock := effect.(*effects.CantBlockEffect)
		return isCantBlock
	})
}

// hasMustAttackEffect checks if a creature is affected by a "must attack if able" continuous effect
// Per Java: RequirementEffect.applies() for attack requirements
func (e *MageEngine) hasMustAttackEffect(gameState *engineGameState, creatureID string) bool {
	if gameState == nil || gameState.layerSystem == nil {
		return false
	}

	// Check if any MustAttackEffect applies to this creature
	return gameState.layerSystem.HasEffectType(creatureID, func(effect effects.ContinuousEffect) bool {
		_, isMustAttack := effect.(*effects.MustAttackEffect)
		return isMustAttack
	})
}

// hasMustBeBlockedEffect checks if an attacker is affected by a "must be blocked if able" continuous effect
// Per Java: RequirementEffect.applies() with mustBlock() returning true
func (e *MageEngine) hasMustBeBlockedEffect(gameState *engineGameState, attackerID string) bool {
	if gameState == nil || gameState.layerSystem == nil {
		return false
	}

	// Check if any MustBeBlockedEffect applies to this attacker
	return gameState.layerSystem.HasEffectType(attackerID, func(effect effects.ContinuousEffect) bool {
		_, isMustBeBlocked := effect.(*effects.MustBeBlockedEffect)
		return isMustBeBlocked
	})
}

// getMustBeBlockedEffects returns all MustBeBlockedEffects for an attacker
// Per Java: ContinuousEffects.getApplicableRequirementEffects() filtering for mustBlock
func (e *MageEngine) getMustBeBlockedEffects(gameState *engineGameState, attackerID string) []*effects.MustBeBlockedEffect {
	if gameState == nil || gameState.layerSystem == nil {
		return nil
	}

	result := make([]*effects.MustBeBlockedEffect, 0)
	allEffects := gameState.layerSystem.GetEffectsForCard(attackerID)
	for _, effect := range allEffects {
		if mbEffect, ok := effect.(*effects.MustBeBlockedEffect); ok {
			result = append(result, mbEffect)
		}
	}

	return result
}

// getMinBlockedBy returns the minimum number of blockers required to block this creature
// Per Java: Permanent.getMinBlockedBy() - default 1, menace sets to 2
func (e *MageEngine) getMinBlockedBy(creature *internalCard) int {
	if e.hasAbility(creature, abilityMenace) {
		return 2
	}
	// TODO: Support other effects that modify minBlockedBy (e.g., "can't be blocked except by 3 or more creatures")
	return 1
}

// hasFirstStrike checks if a creature has first strike (base abilities only)
func (e *MageEngine) hasFirstStrike(creature *internalCard) bool {
	return e.hasAbility(creature, abilityFirstStrike)
}

// hasDoubleStrike checks if a creature has double strike (base abilities only)
func (e *MageEngine) hasDoubleStrike(creature *internalCard) bool {
	return e.hasAbility(creature, abilityDoubleStrike)
}

// hasFirstOrDoubleStrike checks if a creature has first strike or double strike (base abilities only)
func (e *MageEngine) hasFirstOrDoubleStrike(creature *internalCard) bool {
	return e.hasFirstStrike(creature) || e.hasDoubleStrike(creature)
}

// hasFirstStrikeWithEffects checks if a creature has first strike (including granted)
func (e *MageEngine) hasFirstStrikeWithEffects(gameState *engineGameState, creature *internalCard) bool {
	return e.hasAbilityWithEffects(gameState, creature, abilityFirstStrike)
}

// hasDoubleStrikeWithEffects checks if a creature has double strike (including granted)
func (e *MageEngine) hasDoubleStrikeWithEffects(gameState *engineGameState, creature *internalCard) bool {
	return e.hasAbilityWithEffects(gameState, creature, abilityDoubleStrike)
}

// hasFirstOrDoubleStrikeWithEffects checks if a creature has first strike or double strike (including granted)
func (e *MageEngine) hasFirstOrDoubleStrikeWithEffects(gameState *engineGameState, creature *internalCard) bool {
	return e.hasFirstStrikeWithEffects(gameState, creature) || e.hasDoubleStrikeWithEffects(gameState, creature)
}

// recordFirstStrikingCreature records that a creature dealt damage in first strike step
func (e *MageEngine) recordFirstStrikingCreature(gameState *engineGameState, creatureID string) {
	if gameState.combat != nil {
		gameState.combat.firstStrikers[creatureID] = true
	}
}

// wasFirstStrikingCreature checks if a creature dealt damage in first strike step
// This needs to be called with the game state context
func (e *MageEngine) wasFirstStrikingCreature(creature *internalCard) bool {
	// Note: This method should ideally take gameState as parameter
	// For now, it returns false to allow normal damage for creatures without first strike
	// The actual tracking happens in recordFirstStrikingCreature during first strike step
	return false
}

// wasFirstStrikingCreatureInCombat checks if a creature dealt damage in first strike step (with game state)
func (e *MageEngine) wasFirstStrikingCreatureInCombat(gameState *engineGameState, creatureID string) bool {
	if gameState.combat == nil {
		return false
	}
	return gameState.combat.firstStrikers[creatureID]
}

// getCreaturePower gets the power of a creature
// Supports dynamic power calculation via CDAs (Characteristic-Defining Abilities) for cards with * power
func (e *MageEngine) getCreaturePower(gameState *engineGameState, creature *internalCard) (int, error) {
	if creature.Power == "" {
		return 0, nil
	}

	// Check for dynamic power (*, X, or contains *)
	if creature.Power == "*" || creature.Power == "X" || containsStar(creature.Power) {
		// Try to calculate via CDA (Rule 604.3)
		power, err := e.calculateCDAPower(gameState, creature)
		if err == nil {
			return power, nil
		}
		// Fall back to 0 if no CDA found or calculation failed
		return 0, nil
	}

	power, err := strconv.Atoi(creature.Power)
	if err != nil {
		return 0, err
	}

	return power, nil
}

// getCreatureToughness gets the toughness of a creature
// Supports dynamic toughness calculation via CDAs (Characteristic-Defining Abilities) for cards with * toughness
func (e *MageEngine) getCreatureToughness(gameState *engineGameState, creature *internalCard) (int, error) {
	if creature.Toughness == "" {
		return 0, nil
	}

	// Check for dynamic toughness (*, X, or contains *)
	if creature.Toughness == "*" || creature.Toughness == "X" || containsStar(creature.Toughness) {
		// Try to calculate via CDA (Rule 604.3)
		toughness, err := e.calculateCDAToughness(gameState, creature)
		if err == nil {
			return toughness, nil
		}
		// Fall back to 0 if no CDA found or calculation failed
		return 0, nil
	}

	toughness, err := strconv.Atoi(creature.Toughness)
	if err != nil {
		return 0, err
	}

	return toughness, nil
}

// containsStar checks if a string contains an asterisk
func containsStar(s string) bool {
	for _, ch := range s {
		if ch == '*' {
			return true
		}
	}
	return false
}

// calculateCDAPower calculates dynamic power via Characteristic-Defining Abilities
// Rule 604.3: CDAs function in all zones and define characteristics
func (e *MageEngine) calculateCDAPower(gameState *engineGameState, creature *internalCard) (int, error) {
	if gameState.abilityRegistry == nil {
		return 0, fmt.Errorf("ability registry not initialized")
	}

	// Create a GameContext for the CDA to use
	gameIDUUID, err := uuid.Parse(gameState.gameID)
	if err != nil {
		return 0, fmt.Errorf("invalid game ID: %w", err)
	}
	gameContext := NewGameContext(gameIDUUID, e, e.logger)
	ctx := withGameID(context.Background(), gameState.gameID)

	// Check each ability on the creature
	for _, abilityView := range creature.Abilities {
		// Parse ability ID
		abilityID, err := uuid.Parse(abilityView.ID)
		if err != nil {
			continue // Skip invalid IDs
		}

		// Retrieve the actual ability object from the registry
		ability, err := gameState.abilityRegistry.GetAbility(abilityID)
		if err != nil {
			continue // Ability not in registry
		}

		// Check if this is a CDA that defines power
		if cda, ok := ability.(abilities.CharacteristicDefiningAbility); ok {
			if cda.DefinesPower() {
				// Calculate power using the CDA
				power, err := cda.CalculatePower(ctx, gameContext)
				if err == nil {
					return power, nil
				}
			}
		}
	}

	// No CDA found that defines power
	return 0, fmt.Errorf("no CDA found for dynamic power calculation")
}

// calculateCDAToughness calculates dynamic toughness via Characteristic-Defining Abilities
// Rule 604.3: CDAs function in all zones and define characteristics
func (e *MageEngine) calculateCDAToughness(gameState *engineGameState, creature *internalCard) (int, error) {
	if gameState.abilityRegistry == nil {
		return 0, fmt.Errorf("ability registry not initialized")
	}

	// Create a GameContext for the CDA to use
	gameIDUUID, err := uuid.Parse(gameState.gameID)
	if err != nil {
		return 0, fmt.Errorf("invalid game ID: %w", err)
	}
	gameContext := NewGameContext(gameIDUUID, e, e.logger)
	ctx := withGameID(context.Background(), gameState.gameID)

	// Check each ability on the creature
	for _, abilityView := range creature.Abilities {
		// Parse ability ID
		abilityID, err := uuid.Parse(abilityView.ID)
		if err != nil {
			continue // Skip invalid IDs
		}

		// Retrieve the actual ability object from the registry
		ability, err := gameState.abilityRegistry.GetAbility(abilityID)
		if err != nil {
			continue // Ability not in registry
		}

		// Check if this is a CDA that defines toughness
		if cda, ok := ability.(abilities.CharacteristicDefiningAbility); ok {
			if cda.DefinesToughness() {
				// Calculate toughness using the CDA
				toughness, err := cda.CalculateToughness(ctx, gameContext)
				if err == nil {
					return toughness, nil
				}
			}
		}
	}

	// No CDA found that defines toughness
	return 0, fmt.Errorf("no CDA found for dynamic toughness calculation")
}

// getLethalDamage calculates the amount of damage needed to destroy a creature
// This is toughness minus damage already marked on the creature
// Deprecated: Use getLethalDamageWithAttacker instead
// Note: This deprecated function cannot calculate dynamic toughness properly without gameState
func (e *MageEngine) getLethalDamage(creature *internalCard, attackerID string) int {
	// Cannot pass gameState to deprecated function, dynamic toughness will be 0
	toughness, err := e.getCreatureToughness(nil, creature)
	if err != nil {
		return 0
	}

	lethal := toughness - creature.Damage
	if lethal < 0 {
		lethal = 0
	}

	return lethal
}

// getLethalDamageWithAttacker calculates the amount of damage needed to destroy a creature/planeswalker
// considering deathtouch on the attacker. Per Java PermanentImpl.getLethalDamage()
func (e *MageEngine) getLethalDamageWithAttacker(gameState *engineGameState, creature *internalCard, attackerID string) int {
	// For planeswalkers, lethal damage = current loyalty (Rule 306.9)
	// Planeswalkers have no toughness, damage removes loyalty counters
	if e.isPlaneswalker(creature) {
		if creature.Counters == nil {
			return 0
		}
		loyalty := creature.Counters.GetCount("loyalty")

		// With deathtouch, 1 damage is enough to remove all loyalty
		if attackerID != "" {
			if attacker, exists := gameState.cards[attackerID]; exists {
				if e.hasAbility(attacker, abilityDeathtouch) && loyalty > 1 {
					return 1
				}
			}
		}

		return loyalty
	}

	// For creatures, calculate based on toughness
	toughness, err := e.getCreatureToughness(gameState, creature)
	if err != nil {
		return 0
	}

	lethal := toughness - creature.Damage
	if lethal < 0 {
		lethal = 0
	}

	// TODO: For battles, lethal = min(lethal, defense counters)

	// Check for deathtouch on attacker (Java: attacker.getAbilities(game).containsKey(DeathtouchAbility.getInstance().getId()))
	if attackerID != "" {
		if attacker, exists := gameState.cards[attackerID]; exists {
			if e.hasAbility(attacker, abilityDeathtouch) {
				// With deathtouch, any amount of damage is lethal
				if lethal > 1 {
					lethal = 1
				}
			}
		}
	}

	return lethal
}

// markDamage marks damage on a creature from a source
func (e *MageEngine) markDamage(creature *internalCard, amount int, sourceID string) {
	if amount <= 0 {
		return
	}

	// Initialize damage sources map if needed
	if creature.DamageSources == nil {
		creature.DamageSources = make(map[string]int)
	}

	// Add damage
	creature.Damage += amount
	creature.DamageSources[sourceID] += amount
}

// markDamageWithLifelink marks damage and handles lifelink
// Per Java PermanentImpl.markDamage() lines 1119-1126
func (e *MageEngine) markDamageWithLifelink(gameState *engineGameState, creature *internalCard, amount int, sourceID string) {
	if amount <= 0 {
		return
	}

	// Mark the damage
	e.markDamage(creature, amount, sourceID)

	// Check if source has lifelink
	source, exists := gameState.cards[sourceID]
	if exists && e.hasAbility(source, abilityLifelink) {
		// Gain life equal to damage dealt
		controller, exists := gameState.players[source.ControllerID]
		if exists {
			controller.Life += amount

			if e.logger != nil {
				e.logger.Debug("lifelink triggered",
					zap.String("source_id", sourceID),
					zap.String("controller", source.ControllerID),
					zap.Int("life_gained", amount),
				)
			}
		}
	}

	// Fire damaged permanent event for triggers
	// Per Java: DAMAGED_PERMANENT event with flag=true for combat damage
	damagedEvent := rules.Event{
		Type:       rules.EventDamagedPermanent,
		TargetID:   creature.ID,
		SourceID:   sourceID,
		Amount:     amount,
		Controller: creature.ControllerID,
		Flag:       true, // Combat damage
	}
	gameState.eventBus.Publish(damagedEvent)

	// Check for combat damage triggers (e.g., "Whenever ~ deals combat damage to a creature")
	e.checkCombatTriggers(gameState, damagedEvent)
}

// dealDamageToDefender deals damage to a defending player or permanent
// Per Java CombatGroup.defenderDamage()
func (e *MageEngine) dealDamageToDefender(gameState *engineGameState, attacker *internalCard, defenderID string, amount int) error {
	if amount <= 0 {
		return nil
	}

	// Check if defender is a permanent (planeswalker/battle) or player
	if defender, exists := gameState.cards[defenderID]; exists {
		// Defender is a permanent

		// Rule 306.8, 120.3c: Damage dealt to planeswalker removes loyalty counters
		if e.isPlaneswalker(defender) {
			// Apply excess damage from "trample over planeswalkers" ability (Rule 702.19d)
			// Example: Thrasta, Tempest's Roar
			if e.hasAbility(attacker, abilityTrampleOverPlaneswalkers) {
				lethalDamage := e.getLethalDamageWithAttacker(gameState, defender, attacker.ID)

				if lethalDamage >= amount {
					// Normal damage - all damage to planeswalker
					if defender.Counters != nil {
						defender.Counters.RemoveCounter("loyalty", amount)
					}

					// Handle lifelink for the full amount
					if e.hasAbility(attacker, abilityLifelink) {
						controller, exists := gameState.players[attacker.ControllerID]
						if exists {
							controller.Life += amount
						}
					}

					// Fire damaged permanent event
					gameState.eventBus.Publish(rules.Event{
						Type:       rules.EventDamagedPermanent,
						TargetID:   defender.ID,
						SourceID:   attacker.ID,
						Amount:     amount,
						Controller: defender.ControllerID,
						Flag:       true, // Combat damage
					})

					if e.logger != nil {
						e.logger.Debug("damage dealt to planeswalker",
							zap.String("planeswalker_id", defender.ID),
							zap.String("attacker_id", attacker.ID),
							zap.Int("damage", amount),
							zap.Int("loyalty_remaining", defender.Counters.GetCount("loyalty")),
						)
					}
				} else {
					// Damage with excess - lethal to planeswalker, excess to controller
					if defender.Counters != nil {
						defender.Counters.RemoveCounter("loyalty", lethalDamage)
					}

					// Handle lifelink for lethal damage to planeswalker
					if e.hasAbility(attacker, abilityLifelink) {
						controller, exists := gameState.players[attacker.ControllerID]
						if exists {
							controller.Life += lethalDamage
						}
					}

					// Fire damaged permanent event for planeswalker
					gameState.eventBus.Publish(rules.Event{
						Type:       rules.EventDamagedPermanent,
						TargetID:   defender.ID,
						SourceID:   attacker.ID,
						Amount:     lethalDamage,
						Controller: defender.ControllerID,
						Flag:       true, // Combat damage
					})

					if e.logger != nil {
						e.logger.Debug("trample over planeswalkers: lethal damage to planeswalker",
							zap.String("planeswalker_id", defender.ID),
							zap.String("attacker_id", attacker.ID),
							zap.Int("lethal_damage", lethalDamage),
							zap.Int("loyalty_remaining", defender.Counters.GetCount("loyalty")),
						)
					}

					// Recursively deal excess damage to planeswalker's controller
					excessDamage := amount - lethalDamage
					if excessDamage > 0 {
						controllerID := defender.ControllerID
						if e.logger != nil {
							e.logger.Debug("trample over planeswalkers: excess damage to controller",
								zap.String("controller_id", controllerID),
								zap.Int("excess_damage", excessDamage),
							)
						}
						return e.dealDamageToDefender(gameState, attacker, controllerID, excessDamage)
					}
				}
			} else {
				// Normal damage to planeswalker (no trample over)
				if defender.Counters != nil {
					defender.Counters.RemoveCounter("loyalty", amount)
				}

				// Handle lifelink
				if e.hasAbility(attacker, abilityLifelink) {
					controller, exists := gameState.players[attacker.ControllerID]
					if exists {
						controller.Life += amount
					}
				}

				// Fire damaged permanent event for triggers
				damagedEvent := rules.Event{
					Type:       rules.EventDamagedPermanent,
					TargetID:   defender.ID,
					SourceID:   attacker.ID,
					Amount:     amount,
					Controller: defender.ControllerID,
					Flag:       true, // Combat damage
				}
				gameState.eventBus.Publish(damagedEvent)

				if e.logger != nil {
					e.logger.Debug("damage dealt to planeswalker",
						zap.String("planeswalker_id", defender.ID),
						zap.String("attacker_id", attacker.ID),
						zap.Int("damage", amount),
						zap.Int("loyalty_remaining", defender.Counters.GetCount("loyalty")),
					)
				}
			}
			return nil
		} else {
			// For other permanents (battles, creatures), mark damage normally
			e.markDamageWithLifelink(gameState, defender, amount, attacker.ID)
		}
		return nil
	}

	// Defender is a player (or was a permanent that has left the battlefield)
	player, exists := gameState.players[defenderID]
	if !exists {
		// Defender not found - likely a permanent that left the battlefield during combat
		// This is legal; damage simply isn't dealt
		if e.logger != nil {
			e.logger.Debug("defender not found (likely removed from battlefield)",
				zap.String("defender_id", defenderID),
				zap.String("attacker_id", attacker.ID),
			)
		}
		return nil
	}

	// Deal damage to player
	player.Life -= amount

	// Track commander damage if this is combat damage from a commander (Rule 903.10a)
	// "A player that's been dealt 21 or more combat damage by the same commander over the
	// course of the game loses the game."
	if cb := gameState.getCommanderBehavior(); cb != nil && cb.IsCommanderDamageEnabled() && attacker.IsCommander {
		if player.CommanderDamage == nil {
			player.CommanderDamage = make(map[string]int)
		}
		player.CommanderDamage[attacker.ID] += amount

		if e.logger != nil {
			e.logger.Info("commander damage dealt",
				zap.String("commander_id", attacker.ID),
				zap.String("commander_name", attacker.Name),
				zap.String("player", defenderID),
				zap.Int("damage", amount),
				zap.Int("total_from_commander", player.CommanderDamage[attacker.ID]),
				zap.Int("threshold", cb.GetDamageThreshold()),
			)
		}
	}

	// Handle lifelink
	if e.hasAbility(attacker, abilityLifelink) {
		controller, exists := gameState.players[attacker.ControllerID]
		if exists {
			controller.Life += amount

			if e.logger != nil {
				e.logger.Debug("lifelink triggered on player damage",
					zap.String("attacker_id", attacker.ID),
					zap.String("controller", attacker.ControllerID),
					zap.Int("life_gained", amount),
				)
			}
		}
	}

	// Fire damage event (before damage is dealt)
	gameState.eventBus.Publish(rules.Event{
		Type:       rules.EventDamagePlayer,
		TargetID:   defenderID,
		SourceID:   attacker.ID,
		Amount:     amount,
		Controller: attacker.ControllerID,
	})

	// Fire damaged event (after damage is dealt) for triggers
	// Per Java: DAMAGED_PLAYER event with flag=true for combat damage
	damagedEvent := rules.Event{
		Type:       rules.EventDamagedPlayer,
		TargetID:   defenderID,
		SourceID:   attacker.ID,
		Amount:     amount,
		Controller: attacker.ControllerID,
		Flag:       true, // Combat damage
	}
	gameState.eventBus.Publish(damagedEvent)

	// Check for combat damage triggers (e.g., "Whenever ~ deals combat damage")
	e.checkCombatTriggers(gameState, damagedEvent)

	return nil
}

// applyDamageToCreature applies marked damage to a creature and checks for death
func (e *MageEngine) applyDamageToCreature(gameState *engineGameState, creatureID string) error {
	creature, exists := gameState.cards[creatureID]
	if !exists {
		return nil
	}

	if creature.Damage == 0 {
		return nil
	}

	// Get creature's toughness
	toughness, err := e.getCreatureToughness(gameState, creature)
	if err != nil {
		toughness = 0
	}

	// Check if any damage source has deathtouch
	hasDeathtouch := false
	for sourceID := range creature.DamageSources {
		if source, exists := gameState.cards[sourceID]; exists {
			if e.hasAbility(source, abilityDeathtouch) {
				hasDeathtouch = true
				break
			}
		}
	}

	// Check if creature dies (damage >= toughness OR any deathtouch damage)
	shouldDie := (creature.Damage >= toughness && toughness > 0) || (hasDeathtouch && creature.Damage > 0)

	if shouldDie {
		// Creature dies - move to graveyard
		previousZone := creature.Zone
		if err := e.moveCard(gameState, creature, zoneGraveyard, ""); err != nil {
			return err
		}

		// Fire death event (zone change from battlefield to graveyard)
		// Per Java: ZONE_CHANGE event where isDiesEvent() checks fromZone==BATTLEFIELD && toZone==GRAVEYARD
		deathEvent := rules.Event{
			Type:       rules.EventZoneChange,
			TargetID:   creatureID,
			SourceID:   creatureID,
			Controller: creature.ControllerID,
			Zone:       zoneGraveyard,
			Metadata: map[string]string{
				"fromZone": zoneToString(previousZone),
				"toZone":   zoneToString(zoneGraveyard),
			},
		}
		gameState.eventBus.Publish(deathEvent)

		// Check for death triggers (e.g., "Whenever ~ dies" or "Whenever a creature dies")
		e.checkCombatTriggers(gameState, deathEvent)
	}

	return nil
}

// GameStateAccessor implementation for engineGameState

func (s *engineGameState) FindCard(cardID string) (rules.CardInfo, bool) {
	card, found := s.cards[cardID]
	if !found {
		return rules.CardInfo{}, false
	}
	return rules.CardInfo{
		ID:           card.ID,
		Name:         card.Name,
		Type:         card.Type,
		Zone:         card.Zone,
		ControllerID: card.ControllerID,
		OwnerID:      card.OwnerID,
		Tapped:       card.Tapped,
		FaceDown:     card.FaceDown,
	}, true
}

func (s *engineGameState) FindPlayer(playerID string) (rules.PlayerInfo, bool) {
	player, found := s.players[playerID]
	if !found {
		return rules.PlayerInfo{}, false
	}
	return rules.PlayerInfo{
		PlayerID: player.PlayerID,
		Name:     player.Name,
		Life:     player.Life,
		Lost:     player.Lost,
		Left:     player.Left,
	}, true
}

func (s *engineGameState) IsCardInZone(cardID string, zone int) bool {
	card, found := s.cards[cardID]
	if !found {
		return false
	}
	return card.Zone == zone
}

func (s *engineGameState) GetCardZone(cardID string) (int, bool) {
	card, found := s.cards[cardID]
	if !found {
		return 0, false
	}
	return card.Zone, true
}

// TargetGameStateAccessor implementation (already in file, but ensuring completeness)

func (s *engineGameState) FindCardForTarget(cardID string) (targeting.TargetCardInfo, bool) {
	card, found := s.FindCard(cardID)
	if !found {
		return targeting.TargetCardInfo{}, false
	}
	return targeting.TargetCardInfo{
		ID:           card.ID,
		Name:         card.Name,
		Type:         card.Type,
		Zone:         card.Zone,
		ControllerID: card.ControllerID,
		OwnerID:      card.OwnerID,
		Tapped:       card.Tapped,
		FaceDown:     card.FaceDown,
	}, true
}

func (s *engineGameState) FindPlayerForTarget(playerID string) (targeting.TargetPlayerInfo, bool) {
	player, found := s.FindPlayer(playerID)
	if !found {
		return targeting.TargetPlayerInfo{}, false
	}
	return targeting.TargetPlayerInfo{
		PlayerID: player.PlayerID,
		Name:     player.Name,
		Life:     player.Life,
		Lost:     player.Lost,
		Left:     player.Left,
	}, true
}

func (s *engineGameState) GetStackItemsForTarget() []targeting.TargetStackItem {
	if s.stack == nil {
		return []targeting.TargetStackItem{}
	}
	items := s.stack.List()
	result := make([]targeting.TargetStackItem, len(items))
	for i, item := range items {
		result[i] = targeting.TargetStackItem{
			ID:         item.ID,
			Controller: item.Controller,
			Kind:       string(item.Kind),
		}
	}
	return result
}

// HasKeywordAbility checks if a card has a specific keyword ability
// TODO: Implement full keyword ability checking once ability system is fully integrated
func (s *engineGameState) HasKeywordAbility(cardID, keyword string) bool {
	// Stub implementation - will be replaced with full ability system integration
	return false
}

// GetProtectionQualities returns protection qualities (colors/types) for a card
// TODO: Implement full protection quality checking once ability system is fully integrated
func (s *engineGameState) GetProtectionQualities(cardID string) []string {
	// Stub implementation - will be replaced with full ability system integration
	return []string{}
}

// GetCardColor returns the color(s) of a card
// TODO: Enhance to include color-changing effects from layer system
func (s *engineGameState) GetCardColor(cardID string) []string {
	// Find the card
	for _, card := range s.battlefield {
		if card == nil {
			continue
		}
		if card.ID == cardID {
			return parseColors(card.Color)
		}
	}

	// Check other zones if needed
	for _, player := range s.players {
		if player == nil {
			continue
		}
		for _, card := range player.Hand {
			if card == nil {
				continue
			}
			if card.ID == cardID {
				return parseColors(card.Color)
			}
		}
		for _, card := range player.Graveyard {
			if card == nil {
				continue
			}
			if card.ID == cardID {
				return parseColors(card.Color)
			}
		}
	}

	return []string{}
}

// parseColors converts a color string like "WU" to []string{"White", "Blue"}
func parseColors(colorStr string) []string {
	colors := []string{}
	for _, ch := range colorStr {
		switch ch {
		case 'W':
			colors = append(colors, "White")
		case 'U':
			colors = append(colors, "Blue")
		case 'B':
			colors = append(colors, "Black")
		case 'R':
			colors = append(colors, "Red")
		case 'G':
			colors = append(colors, "Green")
		}
	}
	return colors
}

// ====================================
// Direct Manipulation Handlers (Rules-Light)
// ====================================

// handleDirectTap directly taps or untaps a card
func (e *MageEngine) handleDirectTap(gameState *engineGameState, cardID string, tapped bool) error {
	card, exists := gameState.cards[cardID]
	if !exists {
		return fmt.Errorf("card %s not found", cardID)
	}

	card.Tapped = tapped
	action := "tapped"
	if !tapped {
		action = "untapped"
	}
	gameState.addMessage(fmt.Sprintf("%s was %s", card.Name, action), "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": action, "card_id": cardID})
	return nil
}

// handleDirectUntapAll untaps all permanents controlled by a player
func (e *MageEngine) handleDirectUntapAll(gameState *engineGameState, playerID string) error {
	untappedCount := 0
	for _, card := range gameState.battlefield {
		if card == nil {
			continue
		}
		if card.ControllerID == playerID && card.Tapped {
			card.Tapped = false
			untappedCount++
		}
	}

	if untappedCount > 0 {
		gameState.addMessage(fmt.Sprintf("%s untapped all permanents (%d)", playerID, untappedCount), "action")
		e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": "untap_all", "player_id": playerID, "count": untappedCount})
	}
	return nil
}

// handleDirectFlip directly flips a card face-up or face-down
func (e *MageEngine) handleDirectFlip(gameState *engineGameState, cardID string, faceDown bool) error {
	card, exists := gameState.cards[cardID]
	if !exists {
		return fmt.Errorf("card %s not found", cardID)
	}

	card.FaceDown = faceDown
	action := "turned face-down"
	if !faceDown {
		action = "turned face-up"
	}
	gameState.addMessage(fmt.Sprintf("%s was %s", card.Name, action), "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": "flip", "card_id": cardID})
	return nil
}

// handleDirectTransform directly transforms a double-faced card
func (e *MageEngine) handleDirectTransform(gameState *engineGameState, cardID string) error {
	card, exists := gameState.cards[cardID]
	if !exists {
		return fmt.Errorf("card %s not found", cardID)
	}

	card.Transformed = !card.Transformed
	gameState.addMessage(fmt.Sprintf("%s transformed", card.Name), "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": "transform", "card_id": cardID})
	return nil
}

// handleDirectMove directly moves a card to a new zone
func (e *MageEngine) handleDirectMove(gameState *engineGameState, playerID, cardID, zoneName string) error {
	card, exists := gameState.cards[cardID]
	if !exists {
		return fmt.Errorf("card %s not found", cardID)
	}

	// Parse target zone
	targetZone := directParseZone(strings.ToUpper(zoneName))
	if targetZone == -1 {
		return fmt.Errorf("invalid zone: %s", zoneName)
	}

	oldZone := card.Zone

	// Remove from old zone
	e.directRemoveCardFromZone(gameState, card)

	// Move to new zone
	card.Zone = targetZone

	// Add to new zone
	player, exists := gameState.players[card.OwnerID]
	if !exists {
		player = gameState.players[playerID]
	}
	switch targetZone {
	case zoneHand:
		if player != nil {
			player.Hand = append(player.Hand, card)
		}
	case zoneGraveyard:
		if player != nil {
			player.Graveyard = append(player.Graveyard, card)
		}
	case zoneLibrary:
		if player != nil {
			// Top of library (prepend)
			player.Library = append([]*internalCard{card}, player.Library...)
		}
		card.Zone = zoneLibrary // Set actual zone
	case zoneLibraryBottom:
		if player != nil {
			// Bottom of library (append)
			player.Library = append(player.Library, card)
		}
		card.Zone = zoneLibrary // Set actual zone (library, just positioned at bottom)
	case zoneBattlefield:
		gameState.battlefield = append(gameState.battlefield, card)
	case zoneExile:
		gameState.exile = append(gameState.exile, card)
	case zoneCommand:
		gameState.command = append(gameState.command, card)
	}

	displayZoneName := zoneName
	if targetZone == zoneLibraryBottom {
		displayZoneName = "bottom of library"
	}
	gameState.addMessage(fmt.Sprintf("%s moved from %s to %s", card.Name, zoneToString(oldZone), displayZoneName), "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": "move", "card_id": cardID, "zone": zoneName})
	return nil
}

// handleDirectStackAdd adds a card to the visual stack WITHOUT moving it from its current zone.
// This is for rules-light manual tracking - the card stays on battlefield/hand/etc but appears in the stack.
// Per rules-light design: players manually track what's "on the stack" for communication purposes.
func (e *MageEngine) handleDirectStackAdd(gameState *engineGameState, playerID, cardID string) error {
	card, exists := gameState.cards[cardID]
	if !exists {
		return fmt.Errorf("card %s not found", cardID)
	}

	// Create a stack item for visual tracking (card stays in its current zone)
	stackItemID := uuid.New().String()
	stackItem := rules.StackItem{
		ID:          stackItemID,
		Controller:  playerID,
		Description: fmt.Sprintf("%s puts %s on stack", playerID, card.Name),
		Kind:        rules.StackItemKindSpell, // Use spell kind for visual display
		SourceID:    cardID,
		Metadata:    make(map[string]string),
	}
	stackItem.Metadata["manual_add"] = "true"
	stackItem.Metadata["source_zone"] = zoneToString(card.Zone)

	gameState.stack.Push(stackItem)

	gameState.addMessage(fmt.Sprintf("%s adds %s to the stack", playerID, card.Name), "action")
	e.notifyStackUpdate(gameState.gameID, map[string]interface{}{
		"action":      "stack_add",
		"player_id":   playerID,
		"card_id":     cardID,
		"card_name":   card.Name,
		"stack_depth": len(gameState.stack.List()),
	})
	return nil
}

// handleDirectStackRemove removes an item from the stack (for manual resolution tracking)
func (e *MageEngine) handleDirectStackRemove(gameState *engineGameState, playerID, itemID string) error {
	removedItem, found := gameState.stack.Remove(itemID)
	if !found {
		return fmt.Errorf("stack item %s not found", itemID)
	}

	gameState.addMessage(fmt.Sprintf("%s removes %s from the stack", playerID, removedItem.Description), "action")
	e.notifyStackUpdate(gameState.gameID, map[string]interface{}{
		"action":      "stack_remove",
		"player_id":   playerID,
		"item_id":     itemID,
		"stack_depth": len(gameState.stack.List()),
	})
	return nil
}

// handleDirectSetCounter sets a counter on a card to a specific value
func (e *MageEngine) handleDirectSetCounter(gameState *engineGameState, cardID, counterType string, amount int) error {
	card, exists := gameState.cards[cardID]
	if !exists {
		return fmt.Errorf("card %s not found", cardID)
	}

	if card.Counters == nil {
		card.Counters = counters.NewCounters()
	}

	// Remove existing counters of this type first
	card.Counters.RemoveCounter(counterType, card.Counters.GetCount(counterType))

	// Add new amount if positive
	if amount > 0 {
		card.Counters.AddCounter(counters.NewCounter(counterType, amount))
	}

	gameState.addMessage(fmt.Sprintf("Set %d %s counters on %s", amount, counterType, card.Name), "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": "set_counter", "card_id": cardID})
	return nil
}

// handleDirectModifyCounter adds or removes counters from a card
func (e *MageEngine) handleDirectModifyCounter(gameState *engineGameState, cardID, counterType string, delta int) error {
	card, exists := gameState.cards[cardID]
	if !exists {
		return fmt.Errorf("card %s not found", cardID)
	}

	if card.Counters == nil {
		card.Counters = counters.NewCounters()
	}

	oldAmount := card.Counters.GetCount(counterType)
	newAmount := oldAmount + delta
	if newAmount < 0 {
		newAmount = 0
	}

	if delta > 0 {
		card.Counters.AddCounter(counters.NewCounter(counterType, delta))
	} else if delta < 0 {
		card.Counters.RemoveCounter(counterType, -delta)
	}

	action := "added"
	absDelta := delta
	if delta < 0 {
		action = "removed"
		absDelta = -delta
	}
	gameState.addMessage(fmt.Sprintf("%s %d %s counters on %s (now %d)", action, absDelta, counterType, card.Name, newAmount), "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": "modify_counter", "card_id": cardID})
	return nil
}

// handleDirectCreateToken creates a new token on the battlefield
func (e *MageEngine) handleDirectCreateToken(gameState *engineGameState, playerID, name, types, power, toughness, color string, abilities []string) error {
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	tokenID := uuid.New().String()
	token := &internalCard{
		ID:           tokenID,
		Name:         name,
		DisplayName:  name + " Token",
		Type:         types,
		Power:        power,
		Toughness:    toughness,
		Color:        color,
		Zone:         zoneBattlefield,
		ControllerID: playerID,
		OwnerID:      playerID,
		IsToken:      true,
		Counters:     counters.NewCounters(),
		Metadata:     make(map[string]string),
	}

	// Parse abilities into the card
	for _, ability := range abilities {
		ability = strings.TrimSpace(ability)
		if ability != "" {
			token.RulesText += ability + "\n"
		}
	}

	gameState.cards[tokenID] = token
	gameState.battlefield = append(gameState.battlefield, token)

	gameState.addMessage(fmt.Sprintf("%s created a %s token", player.Name, name), "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": "create_token", "token_id": tokenID})
	return nil
}

// handleDirectDestroyToken removes a token from the game
func (e *MageEngine) handleDirectDestroyToken(gameState *engineGameState, cardID string) error {
	card, exists := gameState.cards[cardID]
	if !exists {
		return fmt.Errorf("card %s not found", cardID)
	}

	if !card.IsToken {
		return fmt.Errorf("%s is not a token", card.Name)
	}

	name := card.Name
	e.directRemoveCardFromZone(gameState, card)
	delete(gameState.cards, cardID)

	gameState.addMessage(fmt.Sprintf("%s token was destroyed", name), "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": "destroy_token", "card_id": cardID})
	return nil
}

// handleDirectSetLife sets a player's life total directly
func (e *MageEngine) handleDirectSetLife(gameState *engineGameState, playerID string, amount int) error {
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	oldLife := player.Life
	player.Life = amount

	gameState.addMessage(fmt.Sprintf("%s's life changed from %d to %d", player.Name, oldLife, amount), "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": "set_life", "player_id": playerID})
	return nil
}

// handleDirectModifyLife adds or removes life from a player
func (e *MageEngine) handleDirectModifyLife(gameState *engineGameState, playerID string, delta int) error {
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	player.Life += delta
	action := "gained"
	absDelta := delta
	if delta < 0 {
		action = "lost"
		absDelta = -delta
	}

	gameState.addMessage(fmt.Sprintf("%s %s %d life (now %d)", player.Name, action, absDelta, player.Life), "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": "modify_life", "player_id": playerID})
	return nil
}

// handleDirectDraw draws cards for a player
func (e *MageEngine) handleDirectDraw(gameState *engineGameState, playerID string, count int) error {
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	drawn := 0
	for i := 0; i < count && len(player.Library) > 0; i++ {
		card := player.Library[0]
		player.Library = player.Library[1:]
		card.Zone = zoneHand
		player.Hand = append(player.Hand, card)
		drawn++
	}

	gameState.addMessage(fmt.Sprintf("%s drew %d card(s)", player.Name, drawn), "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": "draw", "player_id": playerID, "count": drawn})
	return nil
}

// handleDirectNextTurn advances to the next turn
// RULES-LIGHT: This allows players to manually advance turns without going through each phase
func (e *MageEngine) handleDirectNextTurn(gameState *engineGameState, requestingPlayerID string) error {
	// Clear combat state first
	e.handleDirectClearCombat(gameState)

	// Get current turn number for logging
	oldTurn := gameState.turnManager.TurnNumber()
	oldActivePlayer := gameState.turnManager.ActivePlayer()

	// Determine next active player (rotate in player order)
	playerIDs := make([]string, 0, len(gameState.players))
	for pid := range gameState.players {
		playerIDs = append(playerIDs, pid)
	}
	// Sort to ensure consistent order
	sort.Strings(playerIDs)

	nextPlayerIndex := 0
	for i, pid := range playerIDs {
		if pid == oldActivePlayer {
			nextPlayerIndex = (i + 1) % len(playerIDs)
			break
		}
	}
	nextActivePlayer := playerIDs[nextPlayerIndex]

	// Advance the turn manager to the beginning of a new turn
	// We'll advance step by step until we reach the cleanup step and wrap to a new turn
	for {
		phase, step := gameState.turnManager.AdvanceStep(nextActivePlayer)
		if phase == rules.PhaseBeginning && step == rules.StepUntap {
			// We've wrapped to a new turn
			break
		}
	}

	newTurn := gameState.turnManager.TurnNumber()
	newActivePlayer := gameState.turnManager.ActivePlayer()

	// Reset lands played for the new active player
	if player, exists := gameState.players[newActivePlayer]; exists {
		player.LandsPlayedThisTurn = 0
	}

	// Perform untap step for the new active player
	e.performUntapStep(gameState, newActivePlayer)

	// Set priority to the new active player
	gameState.turnManager.SetPriority(newActivePlayer)
	if player, exists := gameState.players[newActivePlayer]; exists {
		player.HasPriority = true
	}

	// Reset priority for old active player
	if oldActivePlayer != newActivePlayer {
		if oldPlayer, exists := gameState.players[oldActivePlayer]; exists {
			oldPlayer.HasPriority = false
		}
	}

	// Clear any auto-pass settings that are turn-based
	for _, player := range gameState.players {
		if player.PassUntil == PassUntilNextTurn {
			player.PassUntil = PassUntilNone
		}
		// Check PassUntilMyNextTurn
		if player.PassUntil == PassUntilMyNextTurn && player.PlayerID == newActivePlayer {
			player.PassUntil = PassUntilNone
		}
	}

	gameState.addMessage(fmt.Sprintf("Turn %d → Turn %d (%s's turn)", oldTurn, newTurn, newActivePlayer), "turn")

	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{
		"action":        "next_turn",
		"old_turn":      oldTurn,
		"new_turn":      newTurn,
		"active_player": newActivePlayer,
		"requested_by":  requestingPlayerID,
	})

	return nil
}

// handleDirectClearCombat clears all combat state
// RULES-LIGHT: This allows players to manually clear combat
func (e *MageEngine) handleDirectClearCombat(gameState *engineGameState) error {
	// Clear attacking/blocking status from all creatures
	for _, card := range gameState.battlefield {
		if card.Attacking {
			card.Attacking = false
			card.AttackingWhat = ""
		}
		if card.Blocking {
			card.Blocking = false
			card.BlockingWhat = nil
		}
	}

	// Reset combat state
	gameState.combat = newCombatState()

	gameState.addMessage("Combat cleared", "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": "clear_combat"})
	return nil
}

// handleDirectShuffle shuffles a player's library
func (e *MageEngine) handleDirectShuffle(gameState *engineGameState, playerID string) error {
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// Use existing shuffleLibrary which uses crypto/rand
	e.shuffleLibrary(player)

	gameState.addMessage(fmt.Sprintf("%s shuffles their library", player.Name), "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{"action": "shuffle", "player_id": playerID})
	return nil
}

// handleDirectSetPlayerCounter sets a counter on a player (poison, energy, etc.)
func (e *MageEngine) handleDirectSetPlayerCounter(gameState *engineGameState, playerID, counterType string, amount int) error {
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// Handle common counter types
	switch strings.ToLower(counterType) {
	case "poison":
		player.Poison = amount
	case "energy":
		player.Energy = amount
	default:
		// For other counter types, just log but don't store (struct doesn't support generic counters)
		gameState.addMessage(fmt.Sprintf("%s set %s counters to %d (advisory only)", player.Name, counterType, amount), "action")
		e.notifyGameStateChange(gameState.gameID, map[string]interface{}{
			"action":       "set_player_counter",
			"player_id":    playerID,
			"counter_type": counterType,
			"amount":       amount,
		})
		return nil
	}

	gameState.addMessage(fmt.Sprintf("%s now has %d %s counters", player.Name, amount, counterType), "action")
	e.notifyGameStateChange(gameState.gameID, map[string]interface{}{
		"action":       "set_player_counter",
		"player_id":    playerID,
		"counter_type": counterType,
		"amount":       amount,
	})
	return nil
}

// handleDirectSearchLibrary initiates a library search for a player
// RULES-LIGHT: This reveals the library contents to the player so they can use MOVE commands
func (e *MageEngine) handleDirectSearchLibrary(gameState *engineGameState, playerID, destination string, shuffle bool, message string) error {
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	if len(player.Library) == 0 {
		gameState.addMessage(fmt.Sprintf("%s's library is empty", player.Name), "action")
		return nil
	}

	// Check if there's already a pending library search
	if gameState.pendingLibrarySearch != nil {
		return fmt.Errorf("there is already a pending library search")
	}

	promptMsg := "Search your library"
	if message != "" {
		promptMsg = fmt.Sprintf("Search your library for %s", message)
	}

	gameState.addMessage(fmt.Sprintf("%s searches their library", player.Name), "action")

	// Create the pending library search request - this populates GameView.PendingLibrarySearch
	gameState.pendingLibrarySearch = &PendingLibrarySearchRequest{
		PlayerID:          playerID,
		SearchingPlayerID: playerID,
		Message:           promptMsg,
		Destination:       destination,
		Shuffle:           shuffle,
		Required:          false, // Player can cancel
		Timestamp:         time.Now(),
	}

	return nil
}

// directParseZone converts a zone name string to zone constant
// Special zone constant for bottom of library (not a real zone, but used for move logic)
const zoneLibraryBottom = 100

func directParseZone(zoneName string) int {
	switch strings.ToUpper(zoneName) {
	case "HAND":
		return zoneHand
	case "LIBRARY", "DECK":
		return zoneLibrary
	case "LIBRARY_BOTTOM", "DECK_BOTTOM", "BOTTOM":
		return zoneLibraryBottom
	case "GRAVEYARD", "GRAVE":
		return zoneGraveyard
	case "BATTLEFIELD", "PLAY":
		return zoneBattlefield
	case "STACK":
		return zoneStack
	case "EXILE", "EXILED":
		return zoneExile
	case "COMMAND", "COMMAND_ZONE":
		return zoneCommand
	default:
		return -1
	}
}

// directRemoveCardFromZone removes a card from its current zone
func (e *MageEngine) directRemoveCardFromZone(gameState *engineGameState, card *internalCard) {
	// Remove from player zones
	for _, player := range gameState.players {
		switch card.Zone {
		case zoneHand:
			player.Hand = directRemoveCard(player.Hand, card.ID)
		case zoneLibrary:
			player.Library = directRemoveCard(player.Library, card.ID)
		case zoneGraveyard:
			player.Graveyard = directRemoveCard(player.Graveyard, card.ID)
		}
	}
	// Remove from game-level zones
	switch card.Zone {
	case zoneBattlefield:
		gameState.battlefield = directRemoveCard(gameState.battlefield, card.ID)
	case zoneExile:
		gameState.exile = directRemoveCard(gameState.exile, card.ID)
	case zoneCommand:
		gameState.command = directRemoveCard(gameState.command, card.ID)
	}
}

// directRemoveCard removes a card by ID from a slice
func directRemoveCard(cards []*internalCard, cardID string) []*internalCard {
	result := make([]*internalCard, 0, len(cards))
	for _, c := range cards {
		if c.ID != cardID {
			result = append(result, c)
		}
	}
	return result
}

// ====================================
// Replay Recording System
// ====================================

// StartReplayRecording enables replay recording for a game
// Per Java GameImpl.saveState(): recording can be enabled/disabled
func (e *MageEngine) StartReplayRecording(gameID string) error {
	e.mu.RLock()
	_, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	e.replayRecorder.StartRecording(gameID)
	return nil
}

// StopReplayRecording disables replay recording for a game
func (e *MageEngine) StopReplayRecording(gameID string) {
	e.replayRecorder.StopRecording(gameID)
}

// SaveReplayToFile saves the recorded replay to disk
// Per Java ReplayManager: saves replay for later playback
func (e *MageEngine) SaveReplayToFile(gameID string) error {
	return e.replayRecorder.SaveReplay(gameID)
}

// LoadReplayFromFile loads a previously saved replay
func (e *MageEngine) LoadReplayFromFile(gameID string) (*Replay, error) {
	return e.replayRecorder.LoadReplay(gameID)
}

// GetReplay returns the current replay for a game (if recording)
func (e *MageEngine) GetReplay(gameID string) (*Replay, bool) {
	return e.replayRecorder.GetReplay(gameID)
}

// IsRecordingReplay checks if a game is being recorded
func (e *MageEngine) IsRecordingReplay(gameID string) bool {
	return e.replayRecorder.IsRecording(gameID)
}

// recordReplayState records the current game state to replay if recording is enabled
// Per Java GameImpl.saveState(): saves state at key points
func (e *MageEngine) recordReplayState(gameState *engineGameState) {
	if !e.replayRecorder.IsRecording(gameState.gameID) {
		return
	}

	snapshot := e.createSnapshot(gameState)
	e.replayRecorder.RecordState(gameState.gameID, snapshot)
}
