package game

import (
	"context"

	"github.com/magefree/mage-server-go/internal/repository"
)

// PersistenceAdapter adapts the ActiveGameRepository to the PersistenceRepository interface
// This allows the game engine to persist state without directly depending on the repository package
type PersistenceAdapter struct {
	repo *repository.ActiveGameRepository
}

// NewPersistenceAdapter creates a new persistence adapter
func NewPersistenceAdapter(repo *repository.ActiveGameRepository) *PersistenceAdapter {
	return &PersistenceAdapter{repo: repo}
}

// SaveGameState persists the current game state to the database
func (pa *PersistenceAdapter) SaveGameState(ctx context.Context, gameID, tableID, gameType string, players []string, gameState []byte, turnNumber int, state string) error {
	activeGame := &repository.ActiveGame{
		GameID:     gameID,
		TableID:    tableID,
		GameType:   gameType,
		Players:    players,
		GameState:  gameState,
		TurnNumber: turnNumber,
		State:      state,
	}

	return pa.repo.SaveGameState(ctx, activeGame)
}

// DeleteActiveGame removes a game from persistence (when finished)
func (pa *PersistenceAdapter) DeleteActiveGame(ctx context.Context, gameID string) error {
	return pa.repo.DeleteActiveGame(ctx, gameID)
}

// GetRepository returns the underlying repository for direct access (e.g., for loading all games)
func (pa *PersistenceAdapter) GetRepository() *repository.ActiveGameRepository {
	return pa.repo
}
