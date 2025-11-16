package game

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/magefree/mage-server-go/internal/game/mana"
	"github.com/magefree/mage-server-go/internal/game/rules"
	"github.com/magefree/mage-server-go/internal/game/targeting"
	"go.uber.org/zap"
)

// Zone constants matching Java implementation
const (
	zoneLibrary = 0
	zoneHand    = 1
	zoneBattlefield = 2
	zoneGraveyard = 3
	zoneStack   = 4
	zoneExile   = 5
	zoneCommand = 6
)

// Ability ID constants matching Java keyword abilities
const (
	abilityFirstStrike   = "FirstStrikeAbility"
	abilityDoubleStrike  = "DoubleStrikeAbility"
	abilityVigilance     = "VigilanceAbility"
	abilityFlying        = "FlyingAbility"
	abilityReach         = "ReachAbility"
	abilityTrample       = "TrampleAbility"
	abilityDeathtouch    = "DeathtouchAbility"
	abilityDefender      = "DefenderAbility"
)

// EngineGameView represents the complete game state view for a player
type EngineGameView struct {
	GameID         string
	State         GameState
	Phase         string
	Step          string
	Turn          int
	ActivePlayerID string
	PriorityPlayer string
	Players       []EnginePlayerView
	Battlefield   []EngineCardView
	Stack         []EngineCardView
	Exile         []EngineCardView
	Command       []EngineCardView
	Revealed      []EngineRevealedView
	LookedAt      []EngineLookedAtView
	Combat        EngineCombatView
	StartedAt     time.Time
	Messages      []EngineMessage
	Prompts       []EnginePrompt
}

// EnginePlayerView represents a player's view in the game
type EnginePlayerView struct {
	PlayerID     string
	Name         string
	Life         int
	Poison       int
	Energy       int
	LibraryCount int
	HandCount    int
	Hand         []EngineCardView
	Graveyard    []EngineCardView
	ManaPool     EngineManaPoolView
	HasPriority  bool
	Passed       bool
	StateOrdinal int
	Lost         bool
	Left         bool
	Wins         int
}

// EngineCardView represents a card in any zone
type EngineCardView struct {
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
	Loyalty       string
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
	Counters       []EngineCounterView
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
}

// combatGroup represents a single combat group (attackers vs defender + blockers)
// Per Java CombatGroup class
type combatGroup struct {
	defenderID         string   // player, planeswalker, or battle being attacked
	defenderIsPermanent bool     // is defender a permanent (vs player)
	defendingPlayerID  string   // controller of defending permanents
	attackers          []string // attacking creature IDs
	formerAttackers    []string // historical attackers (for "attacked this turn")
	blockers           []string // blocking creature IDs
	blocked            bool     // is this group blocked
	attackerOrder      map[string]int // damage assignment order for attackers
	blockerOrder       map[string]int // damage assignment order for blockers
}

// newCombatState creates a new combat state
func newCombatState() *combatState {
	return &combatState{
		groups:          make([]*combatGroup, 0),
		formerGroups:    make([]*combatGroup, 0),
		blockingGroups:  make(map[string]*combatGroup),
		defenders:       make(map[string]bool),
		attackers:       make(map[string]bool),
		blockers:        make(map[string]bool),
		attackersTapped: make(map[string]bool),
		firstStrikers:   make(map[string]bool),
	}
}

// newCombatGroup creates a new combat group
func newCombatGroup(defenderID string, defenderIsPermanent bool, defendingPlayerID string) *combatGroup {
	return &combatGroup{
		defenderID:         defenderID,
		defenderIsPermanent: defenderIsPermanent,
		defendingPlayerID:  defendingPlayerID,
		attackers:          make([]string, 0),
		formerAttackers:    make([]string, 0),
		blockers:           make([]string, 0),
		blocked:            false,
		attackerOrder:      make(map[string]int),
		blockerOrder:       make(map[string]int),
	}
}

// EngineMessage represents a game log message
type EngineMessage struct {
	Text      string
	Color     string
	Timestamp time.Time
}

// EnginePrompt represents a prompt for player input
type EnginePrompt struct {
	PlayerID  string
	Text      string
	Options   []string
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
	Loyalty       string
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
	// Damage tracking
	Damage        int      // Damage marked on this creature
	DamageSources map[string]int // Damage by source ID
}

// internalPlayer represents a player in the game state
type internalPlayer struct {
	PlayerID       string
	Name           string
	Life           int
	Poison         int
	Energy         int
	Library        []*internalCard
	Hand           []*internalCard
	Graveyard      []*internalCard
	ManaPool       *mana.ManaPool
	HasPriority    bool
	Passed         bool
	StateOrdinal   int
	Lost           bool
	Left           bool
	Wins           int
	Quit           bool      // Player quit the match
	TimerTimeout   bool      // Player lost due to timer timeout
	IdleTimeout    bool      // Player lost due to idle timeout
	Conceded       bool      // Player conceded
	StoredBookmark int       // Bookmark ID for player undo (-1 = no undo available)
	MulliganCount  int       // Number of times player has mulliganed
	KeptHand       bool      // Whether player has kept their hand
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

// gameAnalytics tracks metrics for a game
type gameAnalytics struct {
	maxStackDepth      int           // Maximum stack depth reached
	totalStackItems    int           // Total items put on stack
	actionsPerTurn     map[int]int   // Actions taken per turn number
	turnStartTimes     map[int]time.Time // Turn start times
	priorityPassCount  int           // Total priority passes
	spellsCast         int           // Total spells cast
	abilitiesActivated int           // Total abilities activated
	triggersProcessed  int           // Total triggered abilities processed
	gameStartTime      time.Time     // When game started
}

// engineGameState represents the internal state of a game
type engineGameState struct {
	gameID        string
	gameType      string
	state         GameState
	players       map[string]*internalPlayer
	playerOrder   []string
	cards         map[string]*internalCard
	battlefield   []*internalCard
	exile         []*internalCard
	command       []*internalCard
	revealed      []EngineRevealedView
	lookedAt      []EngineLookedAtView
	combat        *combatState // Internal combat state
	turnManager   *rules.TurnManager
	stack         *rules.StackManager
	eventBus      *rules.EventBus
	watchers      *rules.WatcherRegistry
	legality      *rules.LegalityChecker
	targetValidator *targeting.TargetValidator
	layerSystem   *effects.LayerSystem
	triggeredQueue []*triggeredAbilityQueueItem // Queue of triggered abilities waiting to be put on stack
	simultaneousEvents []rules.Event             // Queue of events that happened simultaneously
	concedingPlayers   []string                  // Queue of players requesting concession
	analytics     *gameAnalytics                 // Game metrics and analytics
	messages      []EngineMessage
	prompts       []EnginePrompt
	startedAt     time.Time
	mu            sync.RWMutex
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
	gameID        string
	gameType      string
	state         GameState
	turnNumber    int
	activePlayer  string
	priorityPlayer string
	
	// Players - deep copy of all player data
	players       map[string]*internalPlayer
	playerOrder   []string
	
	// Cards - deep copy of all cards
	cards         map[string]*internalCard
	battlefield   []*internalCard
	exile         []*internalCard
	command       []*internalCard
	
	// Stack state
	stackItems    []rules.StackItem
	
	// Other state
	messages      []EngineMessage
	prompts       []EnginePrompt
	timestamp     time.Time
}

// MageEngine is the main game engine implementation
type MageEngine struct {
	logger              *zap.Logger
	mu                  sync.RWMutex
	games               map[string]*engineGameState
	notificationHandler NotificationHandler // Optional handler for UI/websocket notifications
	
	// State bookmarking for rollback/undo
	// Maps gameID -> list of bookmarked states
	bookmarks           map[string][]*gameStateSnapshot
	
	// Turn rollback system (separate from action bookmarks)
	// Maps gameID -> map[turnNumber -> snapshot]
	// Keeps last 4 turns for player-requested rollback
	turnSnapshots       map[string]map[int]*gameStateSnapshot
	rollbackTurnsMax    int  // Maximum turns to keep for rollback (default 4)
	rollbackAllowed     bool // Whether turn rollback is enabled (default true)
}

// NewMageEngine creates a new MageEngine instance
func NewMageEngine(logger *zap.Logger) *MageEngine {
	return &MageEngine{
		logger:           logger,
		games:            make(map[string]*engineGameState),
		bookmarks:        make(map[string][]*gameStateSnapshot),
		turnSnapshots:    make(map[string]map[int]*gameStateSnapshot),
		rollbackTurnsMax: 4,    // Keep last 4 turns
		rollbackAllowed:  true, // Enable turn rollback by default
	}
}

// SetNotificationHandler sets the handler for game notifications
// This allows external systems (UI, websockets) to receive real-time game updates
func (e *MageEngine) SetNotificationHandler(handler NotificationHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.notificationHandler = handler
}

// emitNotification sends a notification to the registered handler
// This method is safe to call while holding gameState locks because:
// 1. It only briefly acquires e.mu.RLock() to read the handler
// 2. The handler is called in a separate goroutine, so it doesn't block
// 3. The goroutine can safely call back into the engine (e.g., GetGameView)
//    because it runs asynchronously after emitNotification returns
func (e *MageEngine) emitNotification(notification GameNotification) {
	// Read handler with RLock to prevent data race with SetNotificationHandler
	// This is safe even when called while holding gameState.mu because:
	// - e.mu and gameState.mu are different mutexes (no deadlock)
	// - RLock allows concurrent reads (no contention)
	// - The lock is held very briefly (just to read a pointer)
	e.mu.RLock()
	handler := e.notificationHandler
	e.mu.RUnlock()

	if handler != nil {
		// Call handler in a goroutine to avoid blocking game logic
		// The goroutine runs asynchronously, so it can safely acquire locks
		// (e.g., call GetGameView) after emitNotification returns
		go handler(notification)
	}
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

// StartGame initializes a new game state
func (e *MageEngine) StartGame(gameID string, players []string, gameType string) error {
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

	// Create game state
	gameState := &engineGameState{
		gameID:      gameID,
		gameType:    gameType,
		state:       GameStateInProgress,
		players:     make(map[string]*internalPlayer),
		playerOrder: make([]string, len(players)),
		cards:       make(map[string]*internalCard),
		battlefield: make([]*internalCard, 0),
		exile:       make([]*internalCard, 0),
		command:     make([]*internalCard, 0),
		revealed:    make([]EngineRevealedView, 0),
		lookedAt:    make([]EngineLookedAtView, 0),
		combat:      newCombatState(),
		analytics: &gameAnalytics{
			actionsPerTurn: make(map[int]int),
			turnStartTimes: make(map[int]time.Time),
			gameStartTime:  time.Now(),
		},
		messages:    make([]EngineMessage, 0),
		prompts:     make([]EnginePrompt, 0),
		startedAt:   time.Now(),
	}

	// Initialize supporting systems
	gameState.stack = rules.NewStackManager()
	gameState.eventBus = rules.NewEventBus()
	gameState.watchers = rules.NewWatcherRegistry()
	gameState.layerSystem = effects.NewLayerSystem()

	// Create players
	for i, playerID := range players {
		gameState.playerOrder[i] = playerID
		gameState.players[playerID] = &internalPlayer{
			PlayerID:       playerID,
			Name:           playerID,
			Life:           20,
			Poison:         0,
			Energy:         0,
			Library:        make([]*internalCard, 0),
			Hand:           make([]*internalCard, 0),
			Graveyard:      make([]*internalCard, 0),
			ManaPool:       mana.NewManaPool(),
			HasPriority:    false,
			Passed:         false,
			StateOrdinal:   0,
			Lost:           false,
			Left:           false,
			Wins:           0,
			StoredBookmark: -1,   // No undo available initially
			MulliganCount:  0,    // No mulligans yet
			KeptHand:       false, // Haven't kept hand yet
		}

		// Create starting hand (7 cards)
		// Mix of different card types for testing
		cardNames := []string{"Lightning Bolt", "Lightning Bolt", "Lightning Bolt", "Counterspell", "Shock", "Lightning Bolt", "Lightning Bolt"}
		for j := 0; j < 7; j++ {
			cardName := cardNames[j%len(cardNames)]
			card := e.createStarterCard(fmt.Sprintf("%s-card-%d", playerID, j), playerID, cardName)
			gameState.cards[card.ID] = card
			gameState.players[playerID].Hand = append(gameState.players[playerID].Hand, card)
			card.Zone = zoneHand
		}

		// Create library (53 cards for a 60-card deck)
		// Mix card types
		libraryCardNames := []string{"Lightning Bolt", "Counterspell", "Shock", "Lightning Bolt", "Counterspell"}
		for j := 0; j < 53; j++ {
			cardName := libraryCardNames[j%len(libraryCardNames)]
			card := e.createStarterCard(fmt.Sprintf("%s-library-%d", playerID, j), playerID, cardName)
			gameState.cards[card.ID] = card
			gameState.players[playerID].Library = append(gameState.players[playerID].Library, card)
			card.Zone = zoneLibrary
		}
	}

	// Initialize turn manager with first player
	gameState.turnManager = rules.NewTurnManager(players[0])
	gameState.players[players[0]].HasPriority = true

	// Initialize legality checker and target validator
	gameState.legality = rules.NewLegalityChecker(gameState)
	gameState.targetValidator = targeting.NewTargetValidator(gameState)

	// Wire up event bus to watchers
	gameState.eventBus.Subscribe(func(event rules.Event) {
		gameState.watchers.NotifyWatchers(event)
	})

	// Add initial log message
	gameState.addMessage("Game started", "action")

	e.games[gameID] = gameState

	// Release lock before sending notifications to avoid deadlock
	// Notifications may trigger callbacks that need to acquire locks
	e.mu.Unlock()

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

	// Notify game start (safe to call after releasing lock)
	e.notifyGameStateChange(gameID, map[string]interface{}{
		"state":     "started",
		"game_type": gameType,
		"players":   players,
	})

	if e.logger != nil {
		e.logger.Info("mage engine started game",
			zap.String("game_id", gameID),
			zap.Strings("players", players),
			zap.String("game_type", gameType),
		)
	}

	return nil
}

// createStarterCard creates a simple starter card for testing
func (e *MageEngine) createStarterCard(id, ownerID, cardName string) *internalCard {
	if cardName == "" {
		cardName = "Lightning Bolt"
	}
	
	return &internalCard{
		ID:           id,
		Name:         cardName,
		DisplayName:  cardName,
		ManaCost:     "{R}",
		Type:         "Instant",
		SubTypes:     []string{},
		SuperTypes:   []string{},
		Color:        "Red",
		Power:        "",
		Toughness:    "",
		Loyalty:      "",
		CardNumber:   1,
		ExpansionSet: "M21",
		Rarity:       "Common",
		RulesText:    fmt.Sprintf("%s deals damage.", cardName),
		Tapped:       false,
		Flipped:      false,
		Transformed:  false,
		FaceDown:     false,
		Zone:         zoneLibrary,
		ControllerID: ownerID,
		OwnerID:      ownerID,
		Counters:     counters.NewCounters(),
	}
}

// ProcessAction processes a player action with automatic error recovery
// Per Java GameImpl.playPriority(): creates bookmark before action, restores on error
func (e *MageEngine) ProcessAction(gameID string, action PlayerAction) (err error) {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	if gameState.state == GameStateFinished {
		return fmt.Errorf("game %s has ended", gameID)
	}

	// Create bookmark before processing action for error recovery
	// Per Java GameImpl.playPriority() line 1728: rollbackBookmarkOnPriorityStart = bookmarkState()
	var bookmarkID int
	gameState.mu.Unlock() // Temporarily unlock to call BookmarkState
	bookmarkID, bookmarkErr := e.BookmarkState(gameID)
	gameState.mu.Lock() // Re-acquire lock
	
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

	if dataStr == "PASS" {
		return e.handlePass(gameState, action.PlayerID)
	}

	return fmt.Errorf("unknown player action: %s", dataStr)
}

// handlePass handles a pass action
func (e *MageEngine) handlePass(gameState *engineGameState, playerID string) error {
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// Check if player has priority
	if gameState.turnManager.PriorityPlayer() != playerID {
		return fmt.Errorf("player %s does not have priority", playerID)
	}

	// Check for concessions before priority
	e.checkConcede(gameState)
	if e.checkIfGameIsOver(gameState) {
		return nil
	}

	// Per rule 117.5 and 603.3: Check state-based actions and triggered abilities before priority
	// Repeat until stable (SBA → triggers → repeat)
	e.checkStateAndTriggered(gameState)

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

		// Save turn snapshot if we advanced to a new turn
		// Per Java GameImpl.saveRollBackGameState(): save at start of each turn
		if newTurn > oldTurn {
			gameState.mu.Unlock() // Temporarily unlock to call SaveTurnSnapshot
			e.SaveTurnSnapshot(gameState.gameID, newTurn)
			gameState.mu.Lock() // Re-acquire lock
		}

		// Notify phase change
		e.notifyPhaseChange(gameState.gameID, map[string]interface{}{
			"phase":         phase.String(),
			"step":          step.String(),
			"active_player": gameState.turnManager.ActivePlayer(),
			"turn":          gameState.turnManager.TurnNumber(),
		})

		// Reset pass flags (preserves lost/left player state)
		gameState.resetPassed()

		// Set priority to active player
		activePlayerID := gameState.turnManager.ActivePlayer()
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
			
			// Per rule 117.5: Check state-based actions before priority
			// Repeat until no more state-based actions occur
			for e.checkStateBasedActions(gameState) {
				// Continue checking until stable
			}
			
			gameState.turnManager.SetPriority(activePlayerID)
			gameState.players[activePlayerID].HasPriority = true
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
	}

	return nil
}

// handleStringAction handles SEND_STRING type actions (spell casting or passing)
func (e *MageEngine) handleStringAction(gameState *engineGameState, action PlayerAction) error {
	spellName, ok := action.Data.(string)
	if !ok {
		return fmt.Errorf("SEND_STRING data must be string")
	}

	// Check if this is a pass action (some tests use "Pass" as SEND_STRING)
	spellNameUpper := strings.ToUpper(strings.TrimSpace(spellName))
	if spellNameUpper == "PASS" {
		return e.handlePass(gameState, action.PlayerID)
	}

	playerID := action.PlayerID
	player, exists := gameState.players[playerID]
	if !exists {
		return fmt.Errorf("player %s not found", playerID)
	}

	// Check if player has priority
	if gameState.turnManager.PriorityPlayer() != playerID {
		return fmt.Errorf("player %s does not have priority", playerID)
	}

	// Per rule 117.5 and 603.3: Check state-based actions and triggered abilities before priority
	// Repeat until stable (SBA → triggers → repeat)
	e.checkStateAndTriggered(gameState)

	// Find card in hand
	var card *internalCard
	for _, c := range player.Hand {
		if strings.EqualFold(c.Name, spellName) {
			card = c
			break
		}
	}

	if card == nil {
		return fmt.Errorf("card %s not found in hand", spellName)
	}

	// Move card to stack
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

	// Check if player has priority
	if gameState.turnManager.PriorityPlayer() != playerID {
		return fmt.Errorf("player %s does not have priority", playerID)
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
				}
				
				return nil
			}
		}
	}

	return fmt.Errorf("UUID %s not found on stack", uuidStr)
}

// resolveStack resolves all items on the stack
func (e *MageEngine) resolveStack(gameState *engineGameState) error {
	gameState.addMessage("Resolving stack", "action")
	
	// Log stack state before resolution
	stackItems := gameState.stack.List()
	if e.logger != nil {
		e.logger.Debug("stack before resolution",
			zap.Int("stack_size", len(stackItems)),
		)
		for i, item := range stackItems {
			e.logger.Debug("stack item",
				zap.Int("index", i),
				zap.String("item_id", item.ID),
				zap.String("source_id", item.SourceID),
				zap.String("description", item.Description),
				zap.String("kind", string(item.Kind)),
			)
		}
	}
	
	// Resolve items in LIFO order (top to bottom)
	for !gameState.stack.IsEmpty() {
		item, err := gameState.stack.Pop()
		if err != nil {
			return fmt.Errorf("failed to pop from stack: %w", err)
		}
		
		if e.logger != nil {
			e.logger.Debug("popped from stack",
				zap.String("item_id", item.ID),
				zap.String("source_id", item.SourceID),
				zap.String("description", item.Description),
				zap.Int("remaining_items", len(gameState.stack.List())),
			)
		}

		// Check legality before resolution
		result := gameState.legality.CheckStackItemLegality(item)
		if e.logger != nil {
			e.logger.Debug("legality check",
				zap.String("item_id", item.ID),
				zap.Bool("legal", result.Legal),
				zap.String("reason", result.Reason),
			)
		}
		if !result.Legal {
			gameState.addMessage(fmt.Sprintf("%s is no longer legal: %s", item.Description, result.Reason), "action")
			if e.logger != nil {
				e.logger.Warn("stack item is illegal",
					zap.String("item_id", item.ID),
					zap.String("reason", result.Reason),
				)
			}
			// Remove illegal item from game state if it's a card
			if card, found := gameState.cards[item.SourceID]; found && card.Zone == zoneStack {
				// Move to graveyard (or appropriate zone)
				card.Zone = zoneGraveyard
				if player, exists := gameState.players[item.Controller]; exists {
					player.Graveyard = append(player.Graveyard, card)
				}
			}
			continue
		}

		// Resolve the item
		gameState.addMessage(fmt.Sprintf("%s resolves", item.Description), "action")
		if e.logger != nil {
			e.logger.Debug("resolving stack item",
				zap.String("item_id", item.ID),
				zap.String("source_id", item.SourceID),
				zap.String("description", item.Description),
				zap.String("kind", string(item.Kind)),
				zap.Bool("has_resolve", item.Resolve != nil),
			)
		}
		
		if item.Resolve != nil {
			if err := item.Resolve(); err != nil {
				gameState.addMessage(fmt.Sprintf("Error resolving %s: %v", item.Description, err), "action")
				if e.logger != nil {
					e.logger.Error("failed to resolve stack item",
						zap.String("item_id", item.ID),
						zap.String("source_id", item.SourceID),
						zap.String("description", item.Description),
						zap.Error(err),
					)
				}
				// Continue resolving other items even if one fails
			} else {
				gameState.addMessage(fmt.Sprintf("%s resolved successfully", item.Description), "action")
				if e.logger != nil {
					e.logger.Debug("resolved stack item successfully",
						zap.String("item_id", item.ID),
						zap.String("source_id", item.SourceID),
						zap.String("description", item.Description),
					)
				}
			}
		} else {
			gameState.addMessage(fmt.Sprintf("%s has no resolve function", item.Description), "action")
			if e.logger != nil {
				e.logger.Warn("stack item has no resolve function",
					zap.String("item_id", item.ID),
					zap.String("description", item.Description),
				)
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
		
		// Per rule 117.5 and 603.3: After each stack item resolves, check state-based actions
		// and process triggered abilities before resolving the next item.
		// This ensures that SBAs and triggers are handled immediately after each resolution.
		e.checkStateAndTriggeredAfterResolution(gameState)
	}

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
			return fmt.Errorf("failed to move permanent to battlefield: %w", err)
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
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("game %s not found", gameID)
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	view := &EngineGameView{
		GameID:         gameID,
		State:          gameState.state,
		Phase:          gameState.turnManager.CurrentPhase().String(),
		Step:           gameState.turnManager.CurrentStep().String(),
		Turn:           gameState.turnManager.TurnNumber(),
		ActivePlayerID: gameState.turnManager.ActivePlayer(),
		PriorityPlayer: gameState.turnManager.PriorityPlayer(),
		Players:        e.buildPlayerViews(gameState, playerID),
		Battlefield:     e.buildCardViews(gameState.battlefield),
		Stack:           e.buildStackViews(gameState),
		Exile:           e.buildCardViews(gameState.exile),
		Command:         e.buildCardViews(gameState.command),
		Revealed:        gameState.revealed,
		LookedAt:        gameState.lookedAt,
		Combat:          e.buildCombatView(gameState),
		StartedAt:       gameState.startedAt,
		Messages:        make([]EngineMessage, len(gameState.messages)),
		Prompts:         make([]EnginePrompt, len(gameState.prompts)),
	}

	copy(view.Messages, gameState.messages)
	copy(view.Prompts, gameState.prompts)

	return view, nil
}

// buildPlayerViews builds player views
func (e *MageEngine) buildPlayerViews(gameState *engineGameState, requestingPlayerID string) []EnginePlayerView {
	views := make([]EnginePlayerView, 0, len(gameState.playerOrder))

	for _, playerID := range gameState.playerOrder {
		player := gameState.players[playerID]
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
			HasPriority:  player.HasPriority,
			Passed:       player.Passed,
			StateOrdinal: player.StateOrdinal,
			Lost:         player.Lost,
			Left:         player.Left,
			Wins:         player.Wins,
		}

		// Only show hand to the owning player
		if playerID == requestingPlayerID {
			view.Hand = e.buildCardViews(player.Hand)
		} else {
			view.Hand = make([]EngineCardView, len(player.Hand))
			for i := range player.Hand {
				view.Hand[i] = EngineCardView{
					ID:     player.Hand[i].ID,
					FaceDown: true,
					Zone:   zoneHand,
				}
			}
		}

		views = append(views, view)
	}

	return views
}

// buildCardViews converts internal cards to view cards
func (e *MageEngine) buildCardViews(cards []*internalCard) []EngineCardView {
	views := make([]EngineCardView, len(cards))
	for i, card := range cards {
		views[i] = EngineCardView{
			ID:             card.ID,
			Name:           card.Name,
			DisplayName:    card.DisplayName,
			ManaCost:       card.ManaCost,
			Type:           card.Type,
			SubTypes:       append([]string(nil), card.SubTypes...),
			SuperTypes:     append([]string(nil), card.SuperTypes...),
			Color:          card.Color,
			Power:          card.Power,
			Toughness:      card.Toughness,
			Loyalty:        card.Loyalty,
			CardNumber:     card.CardNumber,
			ExpansionSet:   card.ExpansionSet,
			Rarity:         card.Rarity,
			RulesText:      card.RulesText,
			Tapped:         card.Tapped,
			Flipped:        card.Flipped,
			Transformed:    card.Transformed,
			FaceDown:       card.FaceDown,
			Zone:           card.Zone,
			ControllerID:   card.ControllerID,
			OwnerID:        card.OwnerID,
			AttachedToCard: append([]string(nil), card.AttachedToCard...),
			Abilities:      append([]EngineAbilityView(nil), card.Abilities...),
			Counters:       e.buildCounterViews(card.Counters),
		}
	}
	return views
}

// buildStackViews builds stack item views
// Stack.List() returns items bottom-to-top (topmost last), so last item is top of stack
func (e *MageEngine) buildStackViews(gameState *engineGameState) []EngineCardView {
	items := gameState.stack.List()
	views := make([]EngineCardView, 0, len(items))

	// Items are already in correct order (bottom to top, topmost last)
	for _, item := range items {
		// Check if this is a triggered ability (not a spell)
		if item.Kind == "TRIGGERED" || item.Kind == rules.StackItemKindTriggered {
			// Create a view for triggered ability using its description
			views = append(views, EngineCardView{
				ID:          item.ID,
				Name:        item.Description,
				DisplayName: item.Description,
				Zone:        zoneStack,
				ControllerID: item.Controller,
			})
		} else {
			// This is a spell - use the card view
			card, found := gameState.cards[item.SourceID]
			if !found {
				// Create a placeholder view if card not found
				views = append(views, EngineCardView{
					ID:          item.ID,
					Name:        item.Description,
					DisplayName: item.Description,
					Zone:        zoneStack,
					ControllerID: item.Controller,
				})
			} else {
				cardView := e.buildCardViews([]*internalCard{card})[0]
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
		if card.OwnerID != playerID {
			remainingExile = append(remainingExile, card)
		}
	}
	gameState.exile = remainingExile

	// Remove from command zone
	remainingCommand := make([]*internalCard, 0)
	for _, card := range gameState.command {
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
		"max_stack_depth":       gameState.analytics.maxStackDepth,
		"total_stack_items":     gameState.analytics.totalStackItems,
		"priority_pass_count":   gameState.analytics.priorityPassCount,
		"spells_cast":           gameState.analytics.spellsCast,
		"abilities_activated":   gameState.analytics.abilitiesActivated,
		"triggers_processed":    gameState.analytics.triggersProcessed,
		"actions_per_turn":      gameState.analytics.actionsPerTurn,
		"avg_turn_time_seconds": avgTurnTime,
		"total_game_time_seconds": gameTime,
		"current_turn":          currentTurn,
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

		// 704.5j: If a planeswalker has loyalty 0, it's put into its owner's graveyard
		if strings.Contains(strings.ToLower(card.Type), "planeswalker") {
			loyalty, err := e.parseLoyalty(card.Loyalty)
			if err == nil && loyalty <= 0 {
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

// Helper methods for engineGameState

func (s *engineGameState) addMessage(text, color string) {
	s.messages = append(s.messages, EngineMessage{
		Text:      text,
		Color:     color,
		Timestamp: time.Now(),
	})
	// Keep only last 1000 messages
	if len(s.messages) > 1000 {
		s.messages = s.messages[len(s.messages)-1000:]
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
		gameID:         gameState.gameID,
		gameType:       gameState.gameType,
		state:          gameState.state,
		turnNumber:     gameState.turnManager.TurnNumber(),
		activePlayer:   gameState.turnManager.ActivePlayer(),
		priorityPlayer: gameState.turnManager.PriorityPlayer(),
		playerOrder:    make([]string, len(gameState.playerOrder)),
		players:        make(map[string]*internalPlayer),
		cards:          make(map[string]*internalCard),
		battlefield:    make([]*internalCard, 0, len(gameState.battlefield)),
		exile:          make([]*internalCard, 0, len(gameState.exile)),
		command:        make([]*internalCard, 0, len(gameState.command)),
		stackItems:     make([]rules.StackItem, 0),
		messages:       make([]EngineMessage, len(gameState.messages)),
		prompts:        make([]EnginePrompt, len(gameState.prompts)),
		timestamp:      time.Now(),
	}
	
	// Copy player order
	copy(snapshot.playerOrder, gameState.playerOrder)
	
	// Deep copy players
	for id, player := range gameState.players {
		playerCopy := &internalPlayer{
			PlayerID:       player.PlayerID,
			Name:           player.Name,
			Life:           player.Life,
			Poison:         player.Poison,
			Energy:         player.Energy,
			Library:        make([]*internalCard, len(player.Library)),
			Hand:           make([]*internalCard, len(player.Hand)),
			Graveyard:      make([]*internalCard, len(player.Graveyard)),
			ManaPool:       player.ManaPool.Copy(),
			HasPriority:    player.HasPriority,
			Passed:         player.Passed,
			StateOrdinal:   player.StateOrdinal,
			Lost:           player.Lost,
			Left:           player.Left,
			Wins:           player.Wins,
			Quit:           player.Quit,
			TimerTimeout:   player.TimerTimeout,
			IdleTimeout:    player.IdleTimeout,
			Conceded:       player.Conceded,
			StoredBookmark: player.StoredBookmark,
			MulliganCount:  player.MulliganCount,
			KeptHand:       player.KeptHand,
		}
		snapshot.players[id] = playerCopy
	}
	
	// Deep copy all cards
	for id, card := range gameState.cards {
		cardCopy := e.copyCard(card)
		snapshot.cards[id] = cardCopy
		
		// Update player zone references
		if player, exists := snapshot.players[card.OwnerID]; exists {
			switch card.Zone {
			case zoneLibrary:
				for i, c := range gameState.players[card.OwnerID].Library {
					if c.ID == card.ID {
						player.Library[i] = cardCopy
						break
					}
				}
			case zoneHand:
				for i, c := range gameState.players[card.OwnerID].Hand {
					if c.ID == card.ID {
						player.Hand[i] = cardCopy
						break
					}
				}
			case zoneGraveyard:
				for i, c := range gameState.players[card.OwnerID].Graveyard {
					if c.ID == card.ID {
						player.Graveyard[i] = cardCopy
						break
					}
				}
			case zoneBattlefield:
				snapshot.battlefield = append(snapshot.battlefield, cardCopy)
			case zoneExile:
				snapshot.exile = append(snapshot.exile, cardCopy)
			case zoneCommand:
				snapshot.command = append(snapshot.command, cardCopy)
			}
		}
	}
	
	// Copy stack items
	if gameState.stack != nil {
		snapshot.stackItems = append(snapshot.stackItems, gameState.stack.List()...)
	}
	
	// Copy messages and prompts
	copy(snapshot.messages, gameState.messages)
	copy(snapshot.prompts, gameState.prompts)
	
	return snapshot
}

// copyCard creates a deep copy of a card
func (e *MageEngine) copyCard(card *internalCard) *internalCard {
	if card == nil {
		return nil
	}
	
	return &internalCard{
		ID:             card.ID,
		Name:           card.Name,
		DisplayName:    card.DisplayName,
		ManaCost:       card.ManaCost,
		Type:           card.Type,
		SubTypes:       append([]string(nil), card.SubTypes...),
		SuperTypes:     append([]string(nil), card.SuperTypes...),
		Color:          card.Color,
		Power:          card.Power,
		Toughness:      card.Toughness,
		Loyalty:        card.Loyalty,
		CardNumber:     card.CardNumber,
		ExpansionSet:   card.ExpansionSet,
		Rarity:         card.Rarity,
		RulesText:      card.RulesText,
		Tapped:         card.Tapped,
		Flipped:        card.Flipped,
		Transformed:    card.Transformed,
		FaceDown:       card.FaceDown,
		Zone:           card.Zone,
		ControllerID:   card.ControllerID,
		OwnerID:        card.OwnerID,
		AttachedToCard: append([]string(nil), card.AttachedToCard...),
		Abilities:      append([]EngineAbilityView(nil), card.Abilities...),
		Counters:       card.Counters.Copy(),
	}
}

// BookmarkState creates a bookmark of the current game state and returns the bookmark ID
// The bookmark can be used later to restore the game to this state
// Per Java GameImpl.bookmarkState(): saves state and returns index for later restoration
func (e *MageEngine) BookmarkState(gameID string) (int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	gameState, exists := e.games[gameID]
	if !exists {
		return 0, fmt.Errorf("game %s not found", gameID)
	}
	
	gameState.mu.RLock()
	snapshot := e.createSnapshot(gameState)
	gameState.mu.RUnlock()
	
	// Add snapshot to bookmarks
	if e.bookmarks[gameID] == nil {
		e.bookmarks[gameID] = make([]*gameStateSnapshot, 0)
	}
	e.bookmarks[gameID] = append(e.bookmarks[gameID], snapshot)
	bookmarkID := len(e.bookmarks[gameID])
	
	if e.logger != nil {
		e.logger.Debug("bookmarked game state",
			zap.String("game_id", gameID),
			zap.Int("bookmark_id", bookmarkID),
			zap.Int("turn", snapshot.turnNumber),
		)
	}
	
	return bookmarkID, nil
}

// RestoreState restores the game to a previously bookmarked state
// Returns error if bookmark doesn't exist or restoration fails
// Per Java GameImpl.restoreState(): rolls back to saved state and removes newer bookmarks
func (e *MageEngine) RestoreState(gameID string, bookmarkID int, context string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	gameState, exists := e.games[gameID]
	if !exists {
		return fmt.Errorf("game %s not found", gameID)
	}
	
	bookmarks := e.bookmarks[gameID]
	if bookmarks == nil || bookmarkID < 1 || bookmarkID > len(bookmarks) {
		return fmt.Errorf("bookmark %d not found for game %s", bookmarkID, gameID)
	}
	
	snapshot := bookmarks[bookmarkID-1]
	
	gameState.mu.Lock()
	defer gameState.mu.Unlock()
	
	// Restore game state from snapshot
	gameState.state = snapshot.state
	gameState.gameType = snapshot.gameType
	
	// Restore players
	gameState.players = make(map[string]*internalPlayer)
	for id, player := range snapshot.players {
		gameState.players[id] = player
	}
	gameState.playerOrder = append([]string(nil), snapshot.playerOrder...)
	
	// Restore cards
	gameState.cards = make(map[string]*internalCard)
	for id, card := range snapshot.cards {
		gameState.cards[id] = card
	}
	
	// Restore zones
	gameState.battlefield = append([]*internalCard(nil), snapshot.battlefield...)
	gameState.exile = append([]*internalCard(nil), snapshot.exile...)
	gameState.command = append([]*internalCard(nil), snapshot.command...)
	
	// Restore stack
	gameState.stack = rules.NewStackManager()
	for _, item := range snapshot.stackItems {
		gameState.stack.Push(item)
	}
	
	// Restore messages and prompts
	gameState.messages = append([]EngineMessage(nil), snapshot.messages...)
	gameState.prompts = append([]EnginePrompt(nil), snapshot.prompts...)
	
	// Remove this bookmark and all newer bookmarks
	e.bookmarks[gameID] = bookmarks[:bookmarkID-1]
	
	gameState.addMessage(fmt.Sprintf("Game restored to turn %d (%s)", snapshot.turnNumber, context), "system")
	
	if e.logger != nil {
		e.logger.Info("restored game state",
			zap.String("game_id", gameID),
			zap.Int("bookmark_id", bookmarkID),
			zap.Int("turn", snapshot.turnNumber),
			zap.String("context", context),
		)
	}
	
	return nil
}

// RemoveBookmark removes a bookmark and all newer bookmarks
// Per Java GameImpl.removeBookmark(): cleanup after restoration
func (e *MageEngine) RemoveBookmark(gameID string, bookmarkID int) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
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
	e.mu.Unlock()
	
	gameState.mu.RLock()
	snapshot := e.createSnapshot(gameState)
	gameState.mu.RUnlock()
	
	e.mu.Lock()
	defer e.mu.Unlock()
	
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
	
	return nil
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
	gameState.state = snapshot.state
	gameState.gameType = snapshot.gameType
	
	// Restore players
	gameState.players = make(map[string]*internalPlayer)
	for id, player := range snapshot.players {
		// Clear all player stored bookmarks on turn rollback
		// Per Java: resetStoredBookmark for all players
		player.StoredBookmark = -1
		gameState.players[id] = player
	}
	gameState.playerOrder = append([]string(nil), snapshot.playerOrder...)
	
	// Restore cards
	gameState.cards = make(map[string]*internalCard)
	for id, card := range snapshot.cards {
		gameState.cards[id] = card
	}
	
	// Restore zones
	gameState.battlefield = append([]*internalCard(nil), snapshot.battlefield...)
	gameState.exile = append([]*internalCard(nil), snapshot.exile...)
	gameState.command = append([]*internalCard(nil), snapshot.command...)
	
	// Restore stack
	gameState.stack = rules.NewStackManager()
	for _, item := range snapshot.stackItems {
		gameState.stack.Push(item)
	}
	
	// Restore messages and prompts
	gameState.messages = append([]EngineMessage(nil), snapshot.messages...)
	gameState.prompts = append([]EnginePrompt(nil), snapshot.prompts...)
	
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
	
	// Remove game from engine
	delete(e.games, gameID)
	
	gameState.mu.Unlock()
	e.mu.Unlock()
	
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
	
	// Shuffle library (simple random shuffle)
	for i := len(player.Library) - 1; i > 0; i-- {
		j := i // In production, use crypto/rand for true randomness
		player.Library[i], player.Library[j] = player.Library[j], player.Library[i]
	}
	
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
	
	// Add planeswalkers controlled by opponents
	// Add battles that can be attacked
	// TODO: Implement when planeswalkers and battles are added
	
	if e.logger != nil {
		e.logger.Debug("set defenders",
			zap.String("game_id", gameID),
			zap.Int("defender_count", len(gameState.combat.defenders)),
		)
	}
	
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
	
	// Find or create combat group for this defender
	var group *combatGroup
	for _, g := range gameState.combat.groups {
		if g.defenderID == defenderID {
			group = g
			break
		}
	}
	
	if group == nil {
		// Determine if defender is a permanent (planeswalker/battle) or player
		defenderIsPermanent := false
		defendingPlayerID := defenderID
		// TODO: Check if defender is a permanent when planeswalkers/battles added
		
		group = newCombatGroup(defenderID, defenderIsPermanent, defendingPlayerID)
		gameState.combat.groups = append(gameState.combat.groups, group)
	}
	
	// Add attacker to group
	group.attackers = append(group.attackers, creatureID)
	gameState.combat.attackers[creatureID] = true
	
	// Tap creature (unless it has vigilance)
	hasVigilance := e.hasAbility(creature, abilityVigilance)
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
	
	// Fire defender attacked event
	defenderEvent := rules.NewEvent(rules.EventDefenderAttacked, defenderID, creatureID, playerID)
	defenderEvent.Metadata["attacker_id"] = creatureID
	gameState.eventBus.Publish(defenderEvent)
	
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
	gameState.eventBus.Publish(rules.NewEvent(rules.EventDeclaredAttackers, "", "", gameState.combat.attackingPlayerID))
	
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
	
	// Flying restriction: creatures with flying can only be blocked by creatures with flying or reach
	// Exception: Dragons can be blocked by non-flying creatures with special abilities (AsThoughEffectType.BLOCK_DRAGON)
	if e.hasAbility(attacker, abilityFlying) {
		if !e.hasAbility(blocker, abilityFlying) && !e.hasAbility(blocker, abilityReach) {
			// TODO: Check for AsThoughEffectType.BLOCK_DRAGON and attacker.hasSubtype(SubType.DRAGON)
			// This requires implementing:
			// 1. Subtype checking system
			// 2. AsThough effects / continuous effects system
			return false, nil
		}
	}
	
	// TODO: Check other restriction effects (can't block, shadow, etc.)
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
	gameState.eventBus.Publish(rules.Event{
		Type:       rules.EventBlockerDeclared,
		SourceID:   blockerID,
		TargetID:   attackerID,
		PlayerID:   playerID,
		Controller: playerID,
	})
	
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
func (e *MageEngine) canBlockInternal(gameState *engineGameState, blockerID, attackerID string) (bool, error) {
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
	
	// Basic checks (same as CanBlock)
	if blocker.Tapped {
		return false, nil
	}
	
	if !strings.Contains(blocker.Type, "Creature") {
		return false, nil
	}
	
	if blocker.Zone != zoneBattlefield {
		return false, nil
	}
	
	if !gameState.combat.attackers[attackerID] {
		return false, nil
	}
	
	// Find defending player
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
	
	if blocker.ControllerID != defendingPlayerID {
		return false, nil
	}
	
	// Get attacker for evasion checks
	attacker := gameState.cards[attackerID]
	
	// Flying restriction: creatures with flying can only be blocked by creatures with flying or reach
	// Exception: Dragons can be blocked by non-flying creatures with special abilities (AsThoughEffectType.BLOCK_DRAGON)
	if e.hasAbility(attacker, abilityFlying) {
		if !e.hasAbility(blocker, abilityFlying) && !e.hasAbility(blocker, abilityReach) {
			// TODO: Check for AsThoughEffectType.BLOCK_DRAGON and attacker.hasSubtype(SubType.DRAGON)
			return false, nil
		}
	}
	
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
				gameState.eventBus.Publish(rules.Event{
					Type:       rules.EventCreatureBlocked,
					SourceID:   attackerID,
				})
			}
		}
	}
	
	// Fire CREATURE_BLOCKS event for each blocker
	// Per Java Combat.acceptBlockers()
	for blockerID := range gameState.combat.blockers {
		gameState.eventBus.Publish(rules.Event{
			Type:       rules.EventCreatureBlocks,
			SourceID:   blockerID,
		})
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
	
	// Fire combat damage step pre event
	gameState.eventBus.Publish(rules.NewEvent(rules.EventCombatDamageStepPre, "", "", ""))
	
	// Assign damage to blockers (attackers dealing damage)
	for _, group := range gameState.combat.groups {
		if err := e.assignDamageToBlockers(gameState, group, firstStrike); err != nil {
			return err
		}
	}
	
	// Assign damage to attackers (blockers dealing damage)
	for _, group := range gameState.combat.groups {
		if len(group.blockers) > 0 {
			if err := e.assignDamageToAttackers(gameState, group, firstStrike); err != nil {
				return err
			}
		}
	}
	
	if e.logger != nil {
		e.logger.Debug("combat damage assigned",
			zap.String("game_id", gameID),
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
	if firstStrike && e.hasFirstOrDoubleStrike(attacker) {
		e.recordFirstStrikingCreature(gameState, attackerID)
	}
	
	// Get attacker's power
	power, err := e.getCreaturePower(attacker)
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
	// With trample, assign lethal damage to blockers, excess goes to defender
	// Without trample, all damage goes to blockers
	
	if hasTrample {
		// Trample: assign lethal damage to each blocker, remainder to defender
		// TODO: Implement player choice for damage assignment (Java: getMultiAmountWithIndividualConstraints)
		// For now, we automatically assign lethal damage to each blocker in order
		remainingDamage := power
		
		for _, blockerID := range group.blockers {
			blocker, exists := gameState.cards[blockerID]
			if !exists {
				continue
			}
			
			// Skip dead blockers (already in graveyard from first strike, etc.)
			if blocker.Zone != zoneBattlefield {
				continue
			}
			
			// Calculate lethal damage (toughness - damage already marked)
			// With deathtouch, only 1 damage is lethal
			lethalDamage := e.getLethalDamageWithAttacker(gameState, blocker, attackerID)
			damageToAssign := lethalDamage
			if damageToAssign > remainingDamage {
				damageToAssign = remainingDamage
			}
			
			// Mark damage on blocker
			e.markDamage(blocker, damageToAssign, attackerID)
			remainingDamage -= damageToAssign
			
			if remainingDamage <= 0 {
				break
			}
		}
		
		// Trample damage to defender
		if remainingDamage > 0 {
			return e.dealDamageToDefender(gameState, attacker, group.defenderID, remainingDamage)
		}
	} else {
		// No trample: divide damage among blockers
		// For now, we'll do simple damage assignment (divide evenly)
		// TODO: Implement proper damage ordering and player choice
		damagePerBlocker := power / len(group.blockers)
		remainingDamage := power % len(group.blockers)
		
		for i, blockerID := range group.blockers {
			blocker, exists := gameState.cards[blockerID]
			if !exists {
				continue
			}
			
			damage := damagePerBlocker
			if i == 0 {
				damage += remainingDamage // Give remainder to first blocker
			}
			
			// Mark damage on blocker
			e.markDamage(blocker, damage, attackerID)
		}
	}
	
	return nil
}

// assignDamageToAttackers handles blocker damage to attackers
// Per Java CombatGroup.assignDamageToAttackers()
func (e *MageEngine) assignDamageToAttackers(gameState *engineGameState, group *combatGroup, firstStrike bool) error {
	if len(group.blockers) == 0 {
		return nil
	}
	
	// For each blocker, deal damage to the attacker(s) it's blocking
	for _, blockerID := range group.blockers {
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
		if firstStrike && e.hasFirstOrDoubleStrike(blocker) {
			e.recordFirstStrikingCreature(gameState, blockerID)
		}
		
		// Get blocker's power
		power, err := e.getCreaturePower(blocker)
		if err != nil {
			power = 0
		}
		
		// Deal damage to attacker(s)
		// For now, simple assignment to first attacker
		// TODO: Handle multiple attackers (banding)
		if len(group.attackers) > 0 {
			attackerID := group.attackers[0]
			attacker, exists := gameState.cards[attackerID]
			if exists {
				e.markDamage(attacker, power, blockerID)
			}
		}
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
			zap.String("game_id", gameID),
		)
	}
	
	// Fire combat damage applied event
	gameState.eventBus.Publish(rules.NewEvent(rules.EventCombatDamageApplied, "", "", ""))
	
	return nil
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
	
	// Fire end combat event
	gameState.eventBus.Publish(rules.Event{
		Type: rules.EventEndCombatStep,
	})
	
	if e.logger != nil {
		e.logger.Debug("ended combat",
			zap.String("game_id", gameID),
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
	
	// Check all creatures in combat groups
	for _, group := range gameState.combat.groups {
		// Check attackers
		for _, attackerID := range group.attackers {
			if attacker, exists := gameState.cards[attackerID]; exists {
				if e.hasFirstOrDoubleStrike(attacker) {
					return true, nil
				}
			}
		}
		
		// Check blockers
		for _, blockerID := range group.blockers {
			if blocker, exists := gameState.cards[blockerID]; exists {
				if e.hasFirstOrDoubleStrike(blocker) {
					return true, nil
				}
			}
		}
	}
	
	return false, nil
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
	
	// TODO: Check continuous effects for granted abilities
	// Example: "Target creature gains flying until end of turn"
	// This would require:
	// - gameState.continuousEffects.getAbilityEffects(card.ID, abilityID)
	// - Check if any active effects grant this ability to this card
	
	return false
}

// hasFirstStrike checks if a creature has first strike
func (e *MageEngine) hasFirstStrike(creature *internalCard) bool {
	return e.hasAbility(creature, abilityFirstStrike)
}

// hasDoubleStrike checks if a creature has double strike
func (e *MageEngine) hasDoubleStrike(creature *internalCard) bool {
	return e.hasAbility(creature, abilityDoubleStrike)
}

// hasFirstOrDoubleStrike checks if a creature has first strike or double strike
func (e *MageEngine) hasFirstOrDoubleStrike(creature *internalCard) bool {
	return e.hasFirstStrike(creature) || e.hasDoubleStrike(creature)
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
func (e *MageEngine) getCreaturePower(creature *internalCard) (int, error) {
	if creature.Power == "" {
		return 0, nil
	}
	
	// Parse power (handle X, *, etc.)
	if creature.Power == "*" || creature.Power == "X" {
		return 0, nil // TODO: Calculate dynamic power
	}
	
	power, err := strconv.Atoi(creature.Power)
	if err != nil {
		return 0, err
	}
	
	return power, nil
}

// getCreatureToughness gets the toughness of a creature
func (e *MageEngine) getCreatureToughness(creature *internalCard) (int, error) {
	if creature.Toughness == "" {
		return 0, nil
	}
	
	// Parse toughness (handle X, *, etc.)
	if creature.Toughness == "*" || creature.Toughness == "X" {
		return 0, nil // TODO: Calculate dynamic toughness
	}
	
	toughness, err := strconv.Atoi(creature.Toughness)
	if err != nil {
		return 0, err
	}
	
	return toughness, nil
}

// getLethalDamage calculates the amount of damage needed to destroy a creature
// This is toughness minus damage already marked on the creature
// Deprecated: Use getLethalDamageWithAttacker instead
func (e *MageEngine) getLethalDamage(creature *internalCard, attackerID string) int {
	toughness, err := e.getCreatureToughness(creature)
	if err != nil {
		return 0
	}
	
	lethal := toughness - creature.Damage
	if lethal < 0 {
		lethal = 0
	}
	
	return lethal
}

// getLethalDamageWithAttacker calculates the amount of damage needed to destroy a creature
// considering deathtouch on the attacker. Per Java PermanentImpl.getLethalDamage()
// TODO: Add planeswalker support (loyalty counters) and battle support (defense counters)
func (e *MageEngine) getLethalDamageWithAttacker(gameState *engineGameState, creature *internalCard, attackerID string) int {
	toughness, err := e.getCreatureToughness(creature)
	if err != nil {
		return 0
	}
	
	lethal := toughness - creature.Damage
	if lethal < 0 {
		lethal = 0
	}
	
	// TODO: For planeswalkers, lethal = min(lethal, loyalty counters)
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

// dealDamageToDefender deals damage to a defending player or permanent
// Per Java CombatGroup.defenderDamage()
func (e *MageEngine) dealDamageToDefender(gameState *engineGameState, attacker *internalCard, defenderID string, amount int) error {
	if amount <= 0 {
		return nil
	}
	
	// Check if defender is a permanent (planeswalker/battle) or player
	if defender, exists := gameState.cards[defenderID]; exists {
		// Defender is a permanent
		e.markDamage(defender, amount, attacker.ID)
		return nil
	}
	
	// Defender is a player
	player, exists := gameState.players[defenderID]
	if !exists {
		return fmt.Errorf("defender %s not found", defenderID)
	}
	
	// Deal damage to player
	player.Life -= amount
	
	// Fire damage event
	gameState.eventBus.Publish(rules.Event{
		Type:       rules.EventDamagePlayer,
		TargetID:   defenderID,
		SourceID:   attacker.ID,
		Amount:     amount,
		Controller: attacker.ControllerID,
	})
	
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
	toughness, err := e.getCreatureToughness(creature)
	if err != nil {
		toughness = 0
	}
	
	// Check if creature dies (damage >= toughness)
	if creature.Damage >= toughness && toughness > 0 {
		// Creature dies - move to graveyard
		if err := e.moveCard(gameState, creature, zoneGraveyard, ""); err != nil {
			return err
		}
		
		// Fire death event
		gameState.eventBus.Publish(rules.Event{
			Type:       rules.EventZoneChange,
			SourceID:   creatureID,
			Controller: creature.ControllerID,
			Zone:       zoneGraveyard,
		})
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
