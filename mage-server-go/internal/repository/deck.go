package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// DeckRepository handles deck database operations
type DeckRepository struct {
	db *DB
}

// NewDeckRepository creates a new deck repository
func NewDeckRepository(db *DB) *DeckRepository {
	return &DeckRepository{db: db}
}

// Create creates a new deck
func (r *DeckRepository) Create(ctx context.Context, d *Deck) error {
	mainDeckJSON, err := d.MainDeckJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal main deck: %w", err)
	}

	sideboardJSON, err := d.SideboardJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal sideboard: %w", err)
	}

	commandersJSON, err := d.CommandersJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize commanders: %w", err)
	}

	query := `
		INSERT INTO decks (user_id, name, format, description, main_deck, sideboard, commanders)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err = r.db.Pool.QueryRow(ctx, query,
		d.UserID, d.Name, d.Format, d.Description, mainDeckJSON, sideboardJSON, commandersJSON).
		Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create deck: %w", err)
	}

	return nil
}

// GetByID retrieves a deck by ID
func (r *DeckRepository) GetByID(ctx context.Context, id int64) (*Deck, error) {
	query := `
		SELECT id, user_id, name, format, description, main_deck, sideboard, commanders,
		       created_at, updated_at
		FROM decks
		WHERE id = $1
	`

	d := &Deck{}
	var mainDeckJSON, sideboardJSON, commandersJSON string

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.UserID, &d.Name, &d.Format, &d.Description,
		&mainDeckJSON, &sideboardJSON, &commandersJSON, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("deck not found")
		}
		return nil, fmt.Errorf("failed to get deck: %w", err)
	}

	if err := d.SetMainDeckFromJSON(mainDeckJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal main deck: %w", err)
	}
	if err := d.SetSideboardFromJSON(sideboardJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sideboard: %w", err)
	}
	if err := d.SetCommandersFromJSON(commandersJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal commanders: %w", err)
	}

	return d, nil
}

// GetByUser retrieves all decks for a user
func (r *DeckRepository) GetByUser(ctx context.Context, userID int64) ([]*Deck, error) {
	query := `
		SELECT id, user_id, name, format, description, main_deck, sideboard, commanders,
		       created_at, updated_at
		FROM decks
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query decks: %w", err)
	}
	defer rows.Close()

	var decks []*Deck
	for rows.Next() {
		d := &Deck{}
		var mainDeckJSON, sideboardJSON, commandersJSON string

		err := rows.Scan(
			&d.ID, &d.UserID, &d.Name, &d.Format, &d.Description,
			&mainDeckJSON, &sideboardJSON, &commandersJSON, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deck: %w", err)
		}

		if err := d.SetMainDeckFromJSON(mainDeckJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal main deck: %w", err)
		}
		if err := d.SetSideboardFromJSON(sideboardJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal sideboard: %w", err)
		}
		if err := d.SetCommandersFromJSON(commandersJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal commanders: %w", err)
		}

		decks = append(decks, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating decks: %w", err)
	}

	return decks, nil
}

// Update updates an existing deck
func (r *DeckRepository) Update(ctx context.Context, d *Deck) error {
	mainDeckJSON, err := d.MainDeckJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal main deck: %w", err)
	}

	sideboardJSON, err := d.SideboardJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal sideboard: %w", err)
	}

	query := `
		UPDATE decks
		SET name = $1, format = $2, description = $3, main_deck = $4, sideboard = $5
		WHERE id = $6
		RETURNING updated_at
	`

	err = r.db.Pool.QueryRow(ctx, query,
		d.Name, d.Format, d.Description, mainDeckJSON, sideboardJSON, d.ID).
		Scan(&d.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("deck not found")
		}
		return fmt.Errorf("failed to update deck: %w", err)
	}

	return nil
}

// Delete deletes a deck by ID
func (r *DeckRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM decks WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete deck: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("deck not found")
	}

	return nil
}

// DeleteByUserAndID deletes a deck by ID and user ID (for ownership check)
func (r *DeckRepository) DeleteByUserAndID(ctx context.Context, userID int64, deckID int64) error {
	query := `DELETE FROM decks WHERE id = $1 AND user_id = $2`

	result, err := r.db.Pool.Exec(ctx, query, deckID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete deck: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("deck not found or not owned by user")
	}

	return nil
}

// GetByUserAndFormat retrieves decks for a user filtered by format
func (r *DeckRepository) GetByUserAndFormat(ctx context.Context, userID int64, format string) ([]*Deck, error) {
	query := `
		SELECT id, user_id, name, format, description, main_deck, sideboard, commanders,
		       created_at, updated_at
		FROM decks
		WHERE user_id = $1 AND format = $2
		ORDER BY updated_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, format)
	if err != nil {
		return nil, fmt.Errorf("failed to query decks: %w", err)
	}
	defer rows.Close()

	var decks []*Deck
	for rows.Next() {
		d := &Deck{}
		var mainDeckJSON, sideboardJSON, commandersJSON string

		err := rows.Scan(
			&d.ID, &d.UserID, &d.Name, &d.Format, &d.Description,
			&mainDeckJSON, &sideboardJSON, &commandersJSON, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deck: %w", err)
		}

		if err := d.SetMainDeckFromJSON(mainDeckJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal main deck: %w", err)
		}
		if err := d.SetSideboardFromJSON(sideboardJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal sideboard: %w", err)
		}
		if err := d.SetCommandersFromJSON(commandersJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal commanders: %w", err)
		}

		decks = append(decks, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating decks: %w", err)
	}

	return decks, nil
}
