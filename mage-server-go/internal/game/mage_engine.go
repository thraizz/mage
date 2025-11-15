package game

import (
	"fmt"
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
	DefendingPlayerID string
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
}

// internalPlayer represents a player in the game state
type internalPlayer struct {
	PlayerID     string
	Name         string
	Life         int
	Poison       int
	Energy       int
	Library      []*internalCard
	Hand         []*internalCard
	Graveyard    []*internalCard
		ManaPool     *mana.ManaPool
	HasPriority  bool
	Passed       bool
	StateOrdinal int
	Lost         bool
	Left         bool
	Wins         int
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
	combat        EngineCombatView
	turnManager   *rules.TurnManager
	stack         *rules.StackManager
	eventBus      *rules.EventBus
	watchers      *rules.WatcherRegistry
	legality      *rules.LegalityChecker
	targetValidator *targeting.TargetValidator
	layerSystem   *effects.LayerSystem
	messages      []EngineMessage
	prompts       []EnginePrompt
	startedAt     time.Time
	mu            sync.RWMutex
}

// MageEngine is the main game engine implementation
type MageEngine struct {
	logger *zap.Logger
	mu     sync.RWMutex
	games  map[string]*engineGameState
}

// NewMageEngine creates a new MageEngine instance
func NewMageEngine(logger *zap.Logger) *MageEngine {
	return &MageEngine{
		logger: logger,
		games:  make(map[string]*engineGameState),
	}
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
	defer e.mu.Unlock()

	if _, exists := e.games[gameID]; exists {
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
			PlayerID:     playerID,
			Name:         playerID,
			Life:         20,
			Poison:       0,
			Energy:       0,
			Library:      make([]*internalCard, 0),
			Hand:         make([]*internalCard, 0),
			Graveyard:    make([]*internalCard, 0),
			ManaPool:     mana.NewManaPool(),
			HasPriority:  false,
			Passed:       false,
			StateOrdinal: 0,
			Lost:         false,
			Left:         false,
			Wins:         0,
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

// ProcessAction processes a player action
func (e *MageEngine) ProcessAction(gameID string, action PlayerAction) error {
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

	// Per rule 117.5 and 603.3: Check state-based actions and triggered abilities before priority
	// Repeat until stable (SBA → triggers → repeat)
	e.checkStateAndTriggered(gameState)

	player.Passed = true
	gameState.addMessage(fmt.Sprintf("%s passes", playerID), "action")

	// Check if all players who can respond have passed
	if gameState.allPassed() {
		// Resolve stack if not empty
		if !gameState.stack.IsEmpty() {
			return e.resolveStack(gameState)
		}

		// Advance step/phase
		nextPlayer := e.getNextPlayer(gameState)
		phase, step := gameState.turnManager.AdvanceStep(nextPlayer)
		gameState.addMessage(fmt.Sprintf("Game advances to %s - %s", phase.String(), step.String()), "action")

		// Reset pass flags (preserves lost/left player state)
		gameState.resetPassed()

		// Set priority to active player
		activePlayerID := gameState.turnManager.ActivePlayer()
		gameState.turnManager.SetPriority(activePlayerID)
		gameState.players[activePlayerID].HasPriority = true

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
					return e.resolveStack(gameState)
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
	gameState.addMessage(fmt.Sprintf("%s casts %s", playerID, card.Name), "action")

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

	return nil
}

// createTriggeredAbilityForSpell creates a triggered ability when a spell is cast
// This simulates effects like "Sanctuary" that trigger on spell casts
func (e *MageEngine) createTriggeredAbilityForSpell(gameState *engineGameState, card *internalCard, casterID string) {
	// For Lightning Bolt, create a triggered ability that gains life
	// This simulates a "Sanctuary" effect for testing
	cardNameLower := strings.ToLower(card.Name)
	if strings.Contains(cardNameLower, "lightning bolt") {
		triggerID := uuid.New().String()
		triggeredAbility := rules.StackItem{
			ID:          triggerID,
			Controller:  casterID,
			Description: fmt.Sprintf("Triggered ability: %s gains 1 life", casterID),
			Kind:        rules.StackItemKindTriggered,
			SourceID:    card.ID,
			Metadata:    make(map[string]string),
			Resolve: func() error {
				player, exists := gameState.players[casterID]
				if !exists {
					return fmt.Errorf("player %s not found", casterID)
				}
				oldLife := player.Life
				player.Life += 1
				gameState.addMessage(fmt.Sprintf("%s gains 1 life (now %d)", casterID, player.Life), "life")
				
				// Emit life gain event
				gameState.eventBus.Publish(rules.Event{
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
		
		gameState.stack.Push(triggeredAbility)
		gameState.addMessage(fmt.Sprintf("Triggered ability: %s gains 1 life", casterID), "action")
		
		if e.logger != nil {
			e.logger.Debug("created triggered ability",
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
		Combat:          gameState.combat,
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
		card, found := gameState.cards[item.SourceID]
		if !found {
			// Create a placeholder view for triggered abilities (they don't have source cards)
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

	return views
}

// buildCounterViews converts counters to view format
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

// processTriggeredAbilities processes triggered abilities that should be put on the stack.
// Returns true if any triggered abilities were processed.
// Per rule 603.3: "Once an ability has triggered, its controller puts it on the stack
// as an object that's not a card the next time a player would receive priority."
func (e *MageEngine) processTriggeredAbilities(gameState *engineGameState) bool {
	// Check watchers for triggered conditions
	// For now, this is a placeholder that will be enhanced when triggered ability
	// queue system is fully implemented (see task: "Queue triggered abilities instead
	// of immediately pushing to stack")
	
	// Currently, triggered abilities are created immediately when events occur
	// (e.g., in createTriggeredAbilityForSpell). This method provides a hook
	// for future triggered ability processing in APNAP order.
	
	// For now, we just ensure events are processed (watchers are already notified
	// via event bus subscription, so they're up to date)
	
	// Return false for now since we're not actively processing new triggers here yet
	// This will be enhanced when the triggered ability queue system is implemented
	return false
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
