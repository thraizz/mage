package game

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// NullEngine is a stub game engine implementation that logs player actions.
type NullEngine struct {
	logger *zap.Logger

	mu    sync.RWMutex
	games map[string]*nullGameState
}

type nullGameState struct {
	GameType string
	Players  []string
	Actions  []PlayerAction
}

// NullGameView represents a snapshot of the null engine state.
type NullGameView struct {
	GameID   string
	GameType string
	Players  []string
	Actions  []PlayerAction
}

// NewNullEngine creates a new null engine.
func NewNullEngine(logger *zap.Logger) *NullEngine {
	return &NullEngine{
		logger: logger,
		games:  make(map[string]*nullGameState),
	}
}

// StartGame initializes a new game state.
func (n *NullEngine) StartGame(gameID string, players []string, gameType string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.games[gameID] = &nullGameState{
		GameType: gameType,
		Players:  append([]string(nil), players...),
		Actions:  make([]PlayerAction, 0, 32),
	}

	if n.logger != nil {
		n.logger.Info("null engine started game",
			zap.String("game_id", gameID),
			zap.Strings("players", players),
			zap.String("game_type", gameType),
		)
	}

	return nil
}

// ProcessAction records the action for later inspection.
func (n *NullEngine) ProcessAction(gameID string, action PlayerAction) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	state, ok := n.games[gameID]
	if !ok {
		return fmt.Errorf("game %s not found", gameID)
	}

	state.Actions = append(state.Actions, action)
	if len(state.Actions) > 200 {
		state.Actions = state.Actions[len(state.Actions)-200:]
	}

	if n.logger != nil {
		n.logger.Debug("null engine processed action",
			zap.String("game_id", gameID),
			zap.String("player_id", action.PlayerID),
			zap.String("action_type", action.ActionType),
		)
	}
	return nil
}

// GetGameView returns a snapshot of the recorded actions.
func (n *NullEngine) GetGameView(gameID, _ string) (interface{}, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	state, ok := n.games[gameID]
	if !ok {
		return nil, fmt.Errorf("game %s not found", gameID)
	}

	actions := make([]PlayerAction, len(state.Actions))
	copy(actions, state.Actions)

	return NullGameView{
		GameID:   gameID,
		GameType: state.GameType,
		Players:  append([]string(nil), state.Players...),
		Actions:  actions,
	}, nil
}

// EndGame removes the game state.
func (n *NullEngine) EndGame(gameID string, winner string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	delete(n.games, gameID)

	if n.logger != nil {
		n.logger.Info("null engine ended game",
			zap.String("game_id", gameID),
			zap.String("winner", winner),
		)
	}

	return nil
}

// PauseGame logs a pause event.
func (n *NullEngine) PauseGame(gameID string) error {
	if n.logger != nil {
		n.logger.Info("null engine pause game", zap.String("game_id", gameID))
	}
	return nil
}

// ResumeGame logs a resume event.
func (n *NullEngine) ResumeGame(gameID string) error {
	if n.logger != nil {
		n.logger.Info("null engine resume game", zap.String("game_id", gameID))
	}
	return nil
}
