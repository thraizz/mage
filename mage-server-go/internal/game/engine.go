package game

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// Engine implements the rules-light game engine
// Copied from playtest-game.ts patterns, adapted for multiplayer server sync
type Engine struct {
	mu       sync.RWMutex
	games    map[string]*EngineGameState
	notifyFn EngineNotificationHandler
	logger   *zap.Logger
}

// EngineNotificationHandler is called when the engine needs to broadcast state changes
type EngineNotificationHandler interface {
	NotifyGameStateChange(playerID string, gameView interface{})
	NotifyGameEvent(gameID string, eventType string, data interface{})
}

// NewEngine creates a new rules-light engine
func NewEngine(logger *zap.Logger) *Engine {
	return &Engine{
		games:  make(map[string]*EngineGameState),
		logger: logger,
	}
}

// SetNotificationHandler sets the callback for game notifications
func (e *Engine) SetNotificationHandler(handler EngineNotificationHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.notifyFn = handler
}

// StartGame initializes and starts a game (without decks - for backwards compatibility)
// Implements GameEngine interface
func (e *Engine) StartGame(gameID string, players []string, gameType string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Create player names map (use player IDs as names for now)
	playerNames := make(map[string]string)
	for _, playerID := range players {
		playerNames[playerID] = playerID
	}

	state := NewEngineGameState(gameID, players, playerNames)
	state.IsInitialized = true
	e.games[gameID] = state

	e.logger.Info("game started",
		zap.String("game_id", gameID),
		zap.String("game_type", gameType),
		zap.Int("player_count", len(players)))

	// Broadcast initial state
	e.broadcast(gameID)
	return nil
}

// StartGameWithDecks initializes and starts a game with player-submitted decks
// Implements GameEngine interface
func (e *Engine) StartGameWithDecks(gameID string, players []string, gameType string, decks map[string]DeckList) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Create player names map
	playerNames := make(map[string]string)
	for _, playerID := range players {
		playerNames[playerID] = playerID
	}

	state := NewEngineGameState(gameID, players, playerNames)

	// Load decks into each player's library
	for playerID, deck := range decks {
		player, ok := state.Players[playerID]
		if !ok {
			continue
		}

		// Create cards from deck list (main deck + commanders)
		allCards := append(deck.MainDeck, deck.Commanders...)
		library := make([]*EngineCard, 0, len(allCards))

		for _, cardName := range allCards {
			card := &EngineCard{
				ID:           fmt.Sprintf("%s-%s-%d", gameID, playerID, len(library)),
				Name:         cardName,
				DisplayName:  cardName,
				OwnerID:      playerID,
				ControllerID: playerID,
				Zone:         ZoneLibraryStr,
				FaceDown:     true,
				Counters:     make([]EngineCounter, 0),
				AttachedTo:   make([]string, 0),
			}
			library = append(library, card)
		}

		player.Library = library
		player.LibraryCount = len(library)

		// Draw opening hand (7 cards)
		handSize := 7
		if len(player.Library) < handSize {
			handSize = len(player.Library)
		}

		openingHand := make([]*EngineCard, handSize)
		copy(openingHand, player.Library[:handSize])
		player.Library = player.Library[handSize:]

		for _, card := range openingHand {
			card.Zone = ZoneHandStr
			card.FaceDown = false
		}

		player.Hand = openingHand
		player.HandCount = len(openingHand)
		player.LibraryCount = len(player.Library)
	}

	state.IsInitialized = true
	e.games[gameID] = state

	e.logger.Info("game started with decks",
		zap.String("game_id", gameID),
		zap.String("game_type", gameType),
		zap.Int("player_count", len(players)))

	// Broadcast initial state
	e.broadcast(gameID)
	return nil
}

// ProcessAction processes a player action
// Implements GameEngine interface
func (e *Engine) ProcessAction(gameID string, action PlayerAction) error {
	// Route action to appropriate handler based on ActionType
	// This is where direct-actions commands would be parsed and dispatched

	actionType := action.ActionType
	playerID := action.PlayerID

	// Parse action data
	data, ok := action.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid action data format")
	}

	switch actionType {
	case "DRAW":
		count := getIntFromData(data, "count", 1)
		return e.DrawCards(gameID, playerID, count)

	case "PLAY":
		cardID := getStringFromData(data, "cardId", "")
		tapped := getBoolFromData(data, "tapped", false)
		return e.PlayCard(gameID, playerID, cardID, tapped)

	case "MOVE":
		cardID := getStringFromData(data, "cardId", "")
		zone := getStringFromData(data, "zone", "")
		return e.MoveCard(gameID, playerID, cardID, zone)

	case "TAP":
		cardID := getStringFromData(data, "cardId", "")
		tapped := getBoolFromData(data, "tapped", true)
		return e.TapCard(gameID, playerID, cardID, tapped)

	case "UNTAP_ALL":
		return e.UntapAll(gameID, playerID)

	case "FLIP":
		cardID := getStringFromData(data, "cardId", "")
		faceDown := getBoolFromData(data, "faceDown", true)
		return e.FlipCard(gameID, playerID, cardID, faceDown)

	case "MODIFY_LIFE":
		delta := getIntFromData(data, "delta", 0)
		targetPlayerID := getStringFromData(data, "targetPlayerId", playerID)
		return e.ModifyLife(gameID, targetPlayerID, delta)

	case "SET_COUNTER":
		targetPlayerID := getStringFromData(data, "targetPlayerId", playerID)
		counterType := getStringFromData(data, "counterType", "")
		value := getIntFromData(data, "value", 0)
		return e.SetPlayerCounter(gameID, targetPlayerID, counterType, value)

	case "SHUFFLE":
		return e.ShuffleLibrary(gameID, playerID)

	case "CREATE_TOKEN":
		name := getStringFromData(data, "name", "Token")
		types := getStringFromData(data, "types", "Creature")
		power := getStringFromData(data, "power", "1")
		toughness := getStringFromData(data, "toughness", "1")
		color := getStringFromData(data, "color", "")
		return e.CreateToken(gameID, playerID, name, types, power, toughness, color)

	case "ADD_COUNTER":
		cardID := getStringFromData(data, "cardId", "")
		counterName := getStringFromData(data, "counterName", "")
		amount := getIntFromData(data, "amount", 1)
		return e.AddCounter(gameID, playerID, cardID, counterName, amount)

	case "REMOVE_COUNTER":
		cardID := getStringFromData(data, "cardId", "")
		counterName := getStringFromData(data, "counterName", "")
		amount := getIntFromData(data, "amount", 1)
		return e.RemoveCounter(gameID, playerID, cardID, counterName, amount)

	case "SET_CARD_COUNTER":
		cardID := getStringFromData(data, "cardId", "")
		counterName := getStringFromData(data, "counterName", "")
		amount := getIntFromData(data, "amount", 0)
		return e.SetCounter(gameID, playerID, cardID, counterName, amount)

	case "MILL":
		count := getIntFromData(data, "count", 1)
		return e.MillCards(gameID, playerID, count)

	case "SCRY":
		scryCount := getIntFromData(data, "scryCount", 1)
		keepOnTop := getStringSliceFromData(data, "keepOnTop")
		putToBottom := getStringSliceFromData(data, "putToBottom")
		return e.ScryCards(gameID, playerID, scryCount, keepOnTop, putToBottom)

	case "SET_REVEALED_TOP":
		revealed := getBoolFromData(data, "revealed", true)
		return e.SetRevealedTop(gameID, playerID, revealed)

	case "NEXT_TURN":
		return e.NextTurn(gameID, playerID)

	case "MULLIGAN":
		return e.Mulligan(gameID, playerID)

	case "KEEP_HAND":
		return e.KeepHand(gameID, playerID)

	default:
		return fmt.Errorf("unknown action type: %s", actionType)
	}
}

// GetGameView returns the current game view for a player
// Implements GameEngine interface
func (e *Engine) GetGameView(gameID, playerID string) (interface{}, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	state := e.games[gameID]
	if state == nil {
		return nil, fmt.Errorf("game not found: %s", gameID)
	}

	// Delegate to view.go for hidden information filtering
	return e.buildGameView(state, playerID), nil
}

// EndGame ends a game
// Implements GameEngine interface
func (e *Engine) EndGame(gameID string, winner string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	msg := fmt.Sprintf("Game ended. Winner: %s", winner)
	e.appendLog(state, "gameEnd", msg)
	e.broadcast(gameID)

	e.logger.Info("game ended",
		zap.String("game_id", gameID),
		zap.String("winner", winner))

	return nil
}

// PauseGame pauses a game
// Implements GameEngine interface
func (e *Engine) PauseGame(gameID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	e.appendLog(state, "pause", "Game paused")
	e.broadcast(gameID)

	e.logger.Info("game paused", zap.String("game_id", gameID))
	return nil
}

// ResumeGame resumes a paused game
// Implements GameEngine interface
func (e *Engine) ResumeGame(gameID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	e.appendLog(state, "resume", "Game resumed")
	e.broadcast(gameID)

	e.logger.Info("game resumed", zap.String("game_id", gameID))
	return nil
}

// broadcast sends game state updates to all players
func (e *Engine) broadcast(gameID string) {
	if e.notifyFn == nil {
		return
	}

	state := e.games[gameID]
	if state == nil {
		return
	}

	// Send personalized view to each player
	for playerID := range state.Players {
		view := e.buildGameView(state, playerID)
		e.notifyFn.NotifyGameStateChange(playerID, view)
	}
}

// Helper functions for extracting data from action maps

func getStringFromData(data map[string]interface{}, key, defaultValue string) string {
	if val, ok := data[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

func getIntFromData(data map[string]interface{}, key string, defaultValue int) int {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case int64:
			return int(v)
		}
	}
	return defaultValue
}

func getBoolFromData(data map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := data[key]; ok {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return defaultValue
}

func getStringSliceFromData(data map[string]interface{}, key string) []string {
	if val, ok := data[key]; ok {
		if sliceVal, ok := val.([]interface{}); ok {
			result := make([]string, 0, len(sliceVal))
			for _, item := range sliceVal {
				if strItem, ok := item.(string); ok {
					result = append(result, strItem)
				}
			}
			return result
		}
		if sliceVal, ok := val.([]string); ok {
			return sliceVal
		}
	}
	return []string{}
}
