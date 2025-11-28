package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// ActiveGameRepository handles active game database operations
type ActiveGameRepository struct {
	db *DB
}

// NewActiveGameRepository creates a new active game repository
func NewActiveGameRepository(db *DB) *ActiveGameRepository {
	return &ActiveGameRepository{db: db}
}

// SaveGameState saves or updates an active game's state
func (r *ActiveGameRepository) SaveGameState(ctx context.Context, game *ActiveGame) error {
	playersJSON, err := game.PlayersJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal players: %w", err)
	}

	query := `
		INSERT INTO active_games (
			game_id, table_id, game_type, players, game_state, turn_number, state, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (game_id) DO UPDATE SET
			game_state = EXCLUDED.game_state,
			turn_number = EXCLUDED.turn_number,
			state = EXCLUDED.state,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`

	err = r.db.Pool.QueryRow(ctx, query,
		game.GameID, game.TableID, game.GameType, playersJSON,
		game.GameState, game.TurnNumber, game.State).
		Scan(&game.ID, &game.CreatedAt, &game.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save active game: %w", err)
	}

	return nil
}

// LoadGameState retrieves an active game by game ID
func (r *ActiveGameRepository) LoadGameState(ctx context.Context, gameID string) (*ActiveGame, error) {
	query := `
		SELECT id, game_id, table_id, game_type, players, game_state, 
		       turn_number, state, created_at, updated_at
		FROM active_games
		WHERE game_id = $1
	`

	game := &ActiveGame{}
	var playersJSON string

	err := r.db.Pool.QueryRow(ctx, query, gameID).Scan(
		&game.ID, &game.GameID, &game.TableID, &game.GameType,
		&playersJSON, &game.GameState, &game.TurnNumber, &game.State,
		&game.CreatedAt, &game.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Game not found - not an error
		}
		return nil, fmt.Errorf("failed to load active game: %w", err)
	}

	if err := game.SetPlayersFromJSON(playersJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal players: %w", err)
	}

	return game, nil
}

// GetActiveGamesForPlayer retrieves all active games for a player
func (r *ActiveGameRepository) GetActiveGamesForPlayer(ctx context.Context, username string) ([]*ActiveGame, error) {
	// Query for games where player is in the players array and state is not FINISHED
	query := `
		SELECT id, game_id, table_id, game_type, players, 
		       turn_number, state, created_at, updated_at
		FROM active_games
		WHERE players ? $1 AND state != 'FINISHED'
		ORDER BY updated_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to query active games for player: %w", err)
	}
	defer rows.Close()

	var games []*ActiveGame
	for rows.Next() {
		game := &ActiveGame{}
		var playersJSON string

		err := rows.Scan(
			&game.ID, &game.GameID, &game.TableID, &game.GameType,
			&playersJSON, &game.TurnNumber, &game.State,
			&game.CreatedAt, &game.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan active game: %w", err)
		}

		if err := game.SetPlayersFromJSON(playersJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal players: %w", err)
		}

		// Don't load game_state for list queries (it can be large)
		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating active games: %w", err)
	}

	return games, nil
}

// DeleteActiveGame removes an active game from the database
func (r *ActiveGameRepository) DeleteActiveGame(ctx context.Context, gameID string) error {
	query := `DELETE FROM active_games WHERE game_id = $1`

	result, err := r.db.Pool.Exec(ctx, query, gameID)
	if err != nil {
		return fmt.Errorf("failed to delete active game: %w", err)
	}

	if result.RowsAffected() == 0 {
		return nil // Game didn't exist - not an error
	}

	return nil
}

// LoadAllActiveGames retrieves all active games (for server startup restoration)
func (r *ActiveGameRepository) LoadAllActiveGames(ctx context.Context) ([]*ActiveGame, error) {
	query := `
		SELECT id, game_id, table_id, game_type, players, game_state, 
		       turn_number, state, created_at, updated_at
		FROM active_games
		WHERE state != 'FINISHED'
		ORDER BY created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all active games: %w", err)
	}
	defer rows.Close()

	var games []*ActiveGame
	for rows.Next() {
		game := &ActiveGame{}
		var playersJSON string

		err := rows.Scan(
			&game.ID, &game.GameID, &game.TableID, &game.GameType,
			&playersJSON, &game.GameState, &game.TurnNumber, &game.State,
			&game.CreatedAt, &game.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan active game: %w", err)
		}

		if err := game.SetPlayersFromJSON(playersJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal players: %w", err)
		}

		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating active games: %w", err)
	}

	return games, nil
}

// UpdateGameState updates just the game state and turn number
func (r *ActiveGameRepository) UpdateGameState(ctx context.Context, gameID string, gameState []byte, turnNumber int, state string) error {
	query := `
		UPDATE active_games
		SET game_state = $2, turn_number = $3, state = $4, updated_at = NOW()
		WHERE game_id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, gameID, gameState, turnNumber, state)
	if err != nil {
		return fmt.Errorf("failed to update game state: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("game not found: %s", gameID)
	}

	return nil
}

// CleanupStaleGames removes games that haven't been updated in the specified duration
func (r *ActiveGameRepository) CleanupStaleGames(ctx context.Context, maxAge time.Duration) (int64, error) {
	query := `
		DELETE FROM active_games
		WHERE updated_at < $1
	`

	cutoff := time.Now().Add(-maxAge)
	result, err := r.db.Pool.Exec(ctx, query, cutoff)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup stale games: %w", err)
	}

	return result.RowsAffected(), nil
}

// CountActiveGames returns the total number of active (non-finished) games
func (r *ActiveGameRepository) CountActiveGames(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM active_games WHERE state != 'FINISHED'`

	var count int
	err := r.db.Pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active games: %w", err)
	}

	return count, nil
}
