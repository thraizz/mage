package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// MatchHistoryRepository handles match history database operations
type MatchHistoryRepository struct {
	db *DB
}

// NewMatchHistoryRepository creates a new match history repository
func NewMatchHistoryRepository(db *DB) *MatchHistoryRepository {
	return &MatchHistoryRepository{db: db}
}

// SaveMatch saves a completed match to the database
func (r *MatchHistoryRepository) SaveMatch(ctx context.Context, match *MatchHistory) error {
	playersJSON, err := match.PlayersJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal players: %w", err)
	}

	query := `
		INSERT INTO match_history (
			game_id, table_id, tournament_id, players, game_type,
			start_time, end_time, duration_seconds,
			winner_id, winner_name, match_options, replay_data
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at
	`

	err = r.db.Pool.QueryRow(ctx, query,
		match.GameID, match.TableID, match.TournamentID, playersJSON, match.GameType,
		match.StartTime, match.EndTime, match.DurationSeconds,
		match.WinnerID, match.WinnerName, match.MatchOptions, match.ReplayData).
		Scan(&match.ID, &match.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save match: %w", err)
	}

	return nil
}

// GetMatchByID retrieves a match by ID
func (r *MatchHistoryRepository) GetMatchByID(ctx context.Context, id int64) (*MatchHistory, error) {
	query := `
		SELECT id, game_id, table_id, tournament_id, players, game_type,
		       start_time, end_time, duration_seconds,
		       winner_id, winner_name, match_options, replay_data, created_at
		FROM match_history
		WHERE id = $1
	`

	match := &MatchHistory{}
	var playersJSON string

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&match.ID, &match.GameID, &match.TableID, &match.TournamentID,
		&playersJSON, &match.GameType,
		&match.StartTime, &match.EndTime, &match.DurationSeconds,
		&match.WinnerID, &match.WinnerName, &match.MatchOptions, &match.ReplayData,
		&match.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("match not found")
		}
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	if err := match.SetPlayersFromJSON(playersJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal players: %w", err)
	}

	return match, nil
}

// GetMatchesByUser retrieves all matches for a user (paginated)
func (r *MatchHistoryRepository) GetMatchesByUser(ctx context.Context, userID int64, limit, offset int) ([]*MatchHistory, error) {
	query := `
		SELECT id, game_id, table_id, tournament_id, players, game_type,
		       start_time, end_time, duration_seconds,
		       winner_id, winner_name, match_options, created_at
		FROM match_history
		WHERE players @> $1::jsonb
		ORDER BY end_time DESC
		LIMIT $2 OFFSET $3
	`

	// Create a JSONB query for user_id
	userQuery := fmt.Sprintf(`[{"user_id": %d}]`, userID)

	rows, err := r.db.Pool.Query(ctx, query, userQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query matches: %w", err)
	}
	defer rows.Close()

	var matches []*MatchHistory
	for rows.Next() {
		match := &MatchHistory{}
		var playersJSON string

		err := rows.Scan(
			&match.ID, &match.GameID, &match.TableID, &match.TournamentID,
			&playersJSON, &match.GameType,
			&match.StartTime, &match.EndTime, &match.DurationSeconds,
			&match.WinnerID, &match.WinnerName, &match.MatchOptions,
			&match.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan match: %w", err)
		}

		if err := match.SetPlayersFromJSON(playersJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal players: %w", err)
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating matches: %w", err)
	}

	return matches, nil
}

// GetRecentMatches retrieves the most recent matches (for lobby display)
func (r *MatchHistoryRepository) GetRecentMatches(ctx context.Context, limit int) ([]*MatchHistory, error) {
	query := `
		SELECT id, game_id, table_id, tournament_id, players, game_type,
		       start_time, end_time, duration_seconds,
		       winner_id, winner_name, match_options, created_at
		FROM match_history
		ORDER BY end_time DESC
		LIMIT $1
	`

	rows, err := r.db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent matches: %w", err)
	}
	defer rows.Close()

	var matches []*MatchHistory
	for rows.Next() {
		match := &MatchHistory{}
		var playersJSON string

		err := rows.Scan(
			&match.ID, &match.GameID, &match.TableID, &match.TournamentID,
			&playersJSON, &match.GameType,
			&match.StartTime, &match.EndTime, &match.DurationSeconds,
			&match.WinnerID, &match.WinnerName, &match.MatchOptions,
			&match.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan match: %w", err)
		}

		if err := match.SetPlayersFromJSON(playersJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal players: %w", err)
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating matches: %w", err)
	}

	return matches, nil
}

// CountMatchesByUser returns the total number of matches for a user
func (r *MatchHistoryRepository) CountMatchesByUser(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM match_history
		WHERE players @> $1::jsonb
	`

	userQuery := fmt.Sprintf(`[{"user_id": %d}]`, userID)

	var count int
	err := r.db.Pool.QueryRow(ctx, query, userQuery).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count matches: %w", err)
	}

	return count, nil
}

// GetMatchesByGameType retrieves matches by game type (paginated)
func (r *MatchHistoryRepository) GetMatchesByGameType(ctx context.Context, gameType string, limit, offset int) ([]*MatchHistory, error) {
	query := `
		SELECT id, game_id, table_id, tournament_id, players, game_type,
		       start_time, end_time, duration_seconds,
		       winner_id, winner_name, match_options, created_at
		FROM match_history
		WHERE game_type = $1
		ORDER BY end_time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, gameType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query matches by game type: %w", err)
	}
	defer rows.Close()

	var matches []*MatchHistory
	for rows.Next() {
		match := &MatchHistory{}
		var playersJSON string

		err := rows.Scan(
			&match.ID, &match.GameID, &match.TableID, &match.TournamentID,
			&playersJSON, &match.GameType,
			&match.StartTime, &match.EndTime, &match.DurationSeconds,
			&match.WinnerID, &match.WinnerName, &match.MatchOptions,
			&match.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan match: %w", err)
		}

		if err := match.SetPlayersFromJSON(playersJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal players: %w", err)
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating matches: %w", err)
	}

	return matches, nil
}

// GetMatchesByTournament retrieves all matches from a specific tournament
func (r *MatchHistoryRepository) GetMatchesByTournament(ctx context.Context, tournamentID string) ([]*MatchHistory, error) {
	query := `
		SELECT id, game_id, table_id, tournament_id, players, game_type,
		       start_time, end_time, duration_seconds,
		       winner_id, winner_name, match_options, created_at
		FROM match_history
		WHERE tournament_id = $1
		ORDER BY start_time ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tournament matches: %w", err)
	}
	defer rows.Close()

	var matches []*MatchHistory
	for rows.Next() {
		match := &MatchHistory{}
		var playersJSON string

		err := rows.Scan(
			&match.ID, &match.GameID, &match.TableID, &match.TournamentID,
			&playersJSON, &match.GameType,
			&match.StartTime, &match.EndTime, &match.DurationSeconds,
			&match.WinnerID, &match.WinnerName, &match.MatchOptions,
			&match.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan match: %w", err)
		}

		if err := match.SetPlayersFromJSON(playersJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal players: %w", err)
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating matches: %w", err)
	}

	return matches, nil
}
