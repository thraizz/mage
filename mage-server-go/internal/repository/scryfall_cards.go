package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// ScryfallCard represents a card from Scryfall bulk data
type ScryfallCard struct {
	ID              uuid.UUID
	OracleID        uuid.UUID
	Name            string
	Lang            string
	SetCode         string
	SetName         string
	CollectorNumber string
	Layout          string
	TypeLine        string
	OracleText      sql.NullString
	ManaCost        sql.NullString
	CMC             float64
	Power           sql.NullString
	Toughness       sql.NullString
	Loyalty         sql.NullString
	Defense         sql.NullString
	Colors          []string
	ColorIdentity   []string
	Rarity          string
	Keywords        []string
	CardFaces       json.RawMessage // JSONB field
	Legalities      json.RawMessage // JSONB field
	ImageURISmall   sql.NullString
	ImageURINormal  sql.NullString
	ImageURILarge   sql.NullString
	ImageURIPNG     sql.NullString
	ReleasedAt      sql.NullTime
	Reprint         bool
	Digital         bool
	Promo           bool
	EDHRECRank      sql.NullInt64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ScryfallSet represents a Magic: The Gathering set from Scryfall
type ScryfallSet struct {
	ID         uuid.UUID
	Code       string
	Name       string
	SetType    string
	ReleasedAt sql.NullTime
	CardCount  int
	Digital    bool
	IconURI    sql.NullString
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ScryfallCardRepository handles Scryfall card database operations
type ScryfallCardRepository struct {
	db     *DB
	logger *zap.Logger
}

// NewScryfallCardRepository creates a new Scryfall card repository
func NewScryfallCardRepository(db *DB, logger *zap.Logger) *ScryfallCardRepository {
	return &ScryfallCardRepository{
		db:     db,
		logger: logger,
	}
}

// GetByOracleID retrieves the latest English printing of a card by oracle ID
func (r *ScryfallCardRepository) GetByOracleID(ctx context.Context, oracleID uuid.UUID) (*ScryfallCard, error) {
	query := `
		SELECT id, oracle_id, name, lang, set_code, set_name, collector_number,
		       layout, type_line, oracle_text, mana_cost, cmc,
		       power, toughness, loyalty, defense,
		       colors, color_identity, rarity, keywords,
		       card_faces, legalities,
		       image_uri_small, image_uri_normal, image_uri_large, image_uri_png,
		       released_at, reprint, digital, promo, edhrec_rank,
		       created_at, updated_at
		FROM scryfall_cards
		WHERE oracle_id = $1 AND lang = 'en'
		ORDER BY released_at DESC
		LIMIT 1
	`

	card := &ScryfallCard{}
	err := r.db.Pool.QueryRow(ctx, query, oracleID).Scan(
		&card.ID, &card.OracleID, &card.Name, &card.Lang, &card.SetCode, &card.SetName, &card.CollectorNumber,
		&card.Layout, &card.TypeLine, &card.OracleText, &card.ManaCost, &card.CMC,
		&card.Power, &card.Toughness, &card.Loyalty, &card.Defense,
		pq.Array(&card.Colors), pq.Array(&card.ColorIdentity), &card.Rarity, pq.Array(&card.Keywords),
		&card.CardFaces, &card.Legalities,
		&card.ImageURISmall, &card.ImageURINormal, &card.ImageURILarge, &card.ImageURIPNG,
		&card.ReleasedAt, &card.Reprint, &card.Digital, &card.Promo, &card.EDHRECRank,
		&card.CreatedAt, &card.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get card by oracle ID: %w", err)
	}

	return card, nil
}

// GetByName retrieves all printings of a card by exact name
// Handles both single-faced cards and multi-faced cards (DFC, split, etc.)
func (r *ScryfallCardRepository) GetByName(ctx context.Context, name string) ([]*ScryfallCard, error) {
	// First try exact match
	cards, err := r.getByNameQuery(ctx, "name = $1", name)
	if err == nil && len(cards) > 0 {
		return cards, nil
	}

	// For multi-faced cards, try matching the front face name
	// Modal DFCs, transforms, etc. are stored as "Front Name // Back Name"
	// So we search for cards where the name starts with the query + " //"
	cards, err = r.getByNameQuery(ctx, "name LIKE $1 || ' //%'", name)
	if err == nil && len(cards) > 0 {
		return cards, nil
	}

	// Also check if the query itself contains "//" (user provided full name)
	if len(cards) == 0 && err == nil {
		return r.getByNameQuery(ctx, "name = $1", name)
	}

	return cards, err
}

// GetByNameCaseInsensitive retrieves cards by name (case-insensitive)
// Handles both single-faced cards and multi-faced cards (DFC, split, etc.)
func (r *ScryfallCardRepository) GetByNameCaseInsensitive(ctx context.Context, name string) ([]*ScryfallCard, error) {
	// First try exact match (case-insensitive)
	cards, err := r.getByNameQuery(ctx, "LOWER(name) = LOWER($1)", name)
	if err == nil && len(cards) > 0 {
		return cards, nil
	}

	// For multi-faced cards, try matching the front face name
	cards, err = r.getByNameQuery(ctx, "LOWER(name) LIKE LOWER($1) || ' //%'", name)
	if err == nil && len(cards) > 0 {
		return cards, nil
	}

	return cards, err
}

// SearchByName performs a full-text search on card names
func (r *ScryfallCardRepository) SearchByName(ctx context.Context, searchTerm string, limit int) ([]*ScryfallCard, error) {
	query := `
		SELECT id, oracle_id, name, lang, set_code, set_name, collector_number,
		       layout, type_line, oracle_text, mana_cost, cmc,
		       power, toughness, loyalty, defense,
		       colors, color_identity, rarity, keywords,
		       card_faces, legalities,
		       image_uri_small, image_uri_normal, image_uri_large, image_uri_png,
		       released_at, reprint, digital, promo, edhrec_rank,
		       created_at, updated_at
		FROM scryfall_cards
		WHERE lang = 'en' AND name ILIKE $1
		ORDER BY name, released_at DESC
		LIMIT $2
	`

	return r.queryCards(ctx, query, "%"+searchTerm+"%", limit)
}

// GetBySetAndNumber retrieves a specific card printing by set and collector number
func (r *ScryfallCardRepository) GetBySetAndNumber(ctx context.Context, setCode, collectorNumber string) (*ScryfallCard, error) {
	query := `
		SELECT id, oracle_id, name, lang, set_code, set_name, collector_number,
		       layout, type_line, oracle_text, mana_cost, cmc,
		       power, toughness, loyalty, defense,
		       colors, color_identity, rarity, keywords,
		       card_faces, legalities,
		       image_uri_small, image_uri_normal, image_uri_large, image_uri_png,
		       released_at, reprint, digital, promo, edhrec_rank,
		       created_at, updated_at
		FROM scryfall_cards
		WHERE set_code = $1 AND collector_number = $2
		LIMIT 1
	`

	card := &ScryfallCard{}
	err := r.db.Pool.QueryRow(ctx, query, setCode, collectorNumber).Scan(
		&card.ID, &card.OracleID, &card.Name, &card.Lang, &card.SetCode, &card.SetName, &card.CollectorNumber,
		&card.Layout, &card.TypeLine, &card.OracleText, &card.ManaCost, &card.CMC,
		&card.Power, &card.Toughness, &card.Loyalty, &card.Defense,
		pq.Array(&card.Colors), pq.Array(&card.ColorIdentity), &card.Rarity, pq.Array(&card.Keywords),
		&card.CardFaces, &card.Legalities,
		&card.ImageURISmall, &card.ImageURINormal, &card.ImageURILarge, &card.ImageURIPNG,
		&card.ReleasedAt, &card.Reprint, &card.Digital, &card.Promo, &card.EDHRECRank,
		&card.CreatedAt, &card.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get card by set and number: %w", err)
	}

	return card, nil
}

// GetCanonicalCards retrieves unique cards (one per oracle_id) from canonical_cards view
func (r *ScryfallCardRepository) GetCanonicalCards(ctx context.Context, limit int) ([]*ScryfallCard, error) {
	query := `
		SELECT id, oracle_id, name, 'en' as lang, set_code, set_name, collector_number,
		       layout, type_line, oracle_text, mana_cost, cmc,
		       power, toughness, loyalty, defense,
		       colors, color_identity, rarity, keywords,
		       card_faces, legalities,
		       NULL as image_uri_small, image_uri_normal, NULL as image_uri_large, NULL as image_uri_png,
		       released_at, false as reprint, false as digital, false as promo, NULL as edhrec_rank,
		       NOW() as created_at, NOW() as updated_at
		FROM canonical_cards
		ORDER BY name
		LIMIT $1
	`

	return r.queryCards(ctx, query, limit)
}

// CountCards returns the total number of cards (optionally filtered by language)
func (r *ScryfallCardRepository) CountCards(ctx context.Context, lang string) (int, error) {
	var query string
	var args []interface{}

	if lang != "" {
		query = "SELECT COUNT(*) FROM scryfall_cards WHERE lang = $1"
		args = []interface{}{lang}
	} else {
		query = "SELECT COUNT(*) FROM scryfall_cards"
	}

	var count int
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count cards: %w", err)
	}

	return count, nil
}

// CountUniqueCards returns the number of unique cards (by oracle_id)
func (r *ScryfallCardRepository) CountUniqueCards(ctx context.Context, lang string) (int, error) {
	var query string
	var args []interface{}

	if lang != "" {
		query = "SELECT COUNT(DISTINCT oracle_id) FROM scryfall_cards WHERE lang = $1"
		args = []interface{}{lang}
	} else {
		query = "SELECT COUNT(DISTINCT oracle_id) FROM scryfall_cards"
	}

	var count int
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count unique cards: %w", err)
	}

	return count, nil
}

// Helper methods

func (r *ScryfallCardRepository) getByNameQuery(ctx context.Context, whereClause string, args ...interface{}) ([]*ScryfallCard, error) {
	query := fmt.Sprintf(`
		SELECT id, oracle_id, name, lang, set_code, set_name, collector_number,
		       layout, type_line, oracle_text, mana_cost, cmc,
		       power, toughness, loyalty, defense,
		       colors, color_identity, rarity, keywords,
		       card_faces, legalities,
		       image_uri_small, image_uri_normal, image_uri_large, image_uri_png,
		       released_at, reprint, digital, promo, edhrec_rank,
		       created_at, updated_at
		FROM scryfall_cards
		WHERE lang = 'en' AND %s
		ORDER BY released_at DESC
	`, whereClause)

	return r.queryCards(ctx, query, args...)
}

func (r *ScryfallCardRepository) queryCards(ctx context.Context, query string, args ...interface{}) ([]*ScryfallCard, error) {
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query cards: %w", err)
	}
	defer rows.Close()

	cards := make([]*ScryfallCard, 0)
	for rows.Next() {
		card := &ScryfallCard{}
		err := rows.Scan(
			&card.ID, &card.OracleID, &card.Name, &card.Lang, &card.SetCode, &card.SetName, &card.CollectorNumber,
			&card.Layout, &card.TypeLine, &card.OracleText, &card.ManaCost, &card.CMC,
			&card.Power, &card.Toughness, &card.Loyalty, &card.Defense,
			pq.Array(&card.Colors), pq.Array(&card.ColorIdentity), &card.Rarity, pq.Array(&card.Keywords),
			&card.CardFaces, &card.Legalities,
			&card.ImageURISmall, &card.ImageURINormal, &card.ImageURILarge, &card.ImageURIPNG,
			&card.ReleasedAt, &card.Reprint, &card.Digital, &card.Promo, &card.EDHRECRank,
			&card.CreatedAt, &card.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// ScryfallSetRepository handles Scryfall set database operations
type ScryfallSetRepository struct {
	db     *DB
	logger *zap.Logger
}

// NewScryfallSetRepository creates a new Scryfall set repository
func NewScryfallSetRepository(db *DB, logger *zap.Logger) *ScryfallSetRepository {
	return &ScryfallSetRepository{
		db:     db,
		logger: logger,
	}
}

// GetByCode retrieves a set by its code
func (r *ScryfallSetRepository) GetByCode(ctx context.Context, code string) (*ScryfallSet, error) {
	query := `
		SELECT id, code, name, set_type, released_at, card_count, digital, icon_uri, created_at, updated_at
		FROM scryfall_sets
		WHERE code = $1
	`

	set := &ScryfallSet{}
	err := r.db.Pool.QueryRow(ctx, query, code).Scan(
		&set.ID, &set.Code, &set.Name, &set.SetType, &set.ReleasedAt,
		&set.CardCount, &set.Digital, &set.IconURI, &set.CreatedAt, &set.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get set: %w", err)
	}

	return set, nil
}

// List retrieves all sets
func (r *ScryfallSetRepository) List(ctx context.Context) ([]*ScryfallSet, error) {
	query := `
		SELECT id, code, name, set_type, released_at, card_count, digital, icon_uri, created_at, updated_at
		FROM scryfall_sets
		ORDER BY released_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list sets: %w", err)
	}
	defer rows.Close()

	sets := make([]*ScryfallSet, 0)
	for rows.Next() {
		set := &ScryfallSet{}
		err := rows.Scan(
			&set.ID, &set.Code, &set.Name, &set.SetType, &set.ReleasedAt,
			&set.CardCount, &set.Digital, &set.IconURI, &set.CreatedAt, &set.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan set: %w", err)
		}
		sets = append(sets, set)
	}

	return sets, nil
}
