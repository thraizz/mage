package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Card represents a Magic card
type Card struct {
	ID           int64
	CardNumber   string
	SetCode      string
	Name         string
	CardType     string
	ManaCost     string
	Power        string
	Toughness    string
	RulesText    string
	FlavorText   sql.NullString // Nullable
	OriginalText sql.NullString // Nullable
	OriginalType sql.NullString // Nullable
	CN           sql.NullInt64  // Nullable - basic lands have NULL collector numbers
	CardName     sql.NullString // Nullable - basic lands have NULL card_name
	Rarity       string
	CreatedAt    time.Time
}

// CardRepository handles card database operations
type CardRepository struct {
	db     *DB
	cache  *cardCache
	logger *zap.Logger
}

// NewCardRepository creates a new card repository
func NewCardRepository(db *DB, logger *zap.Logger) *CardRepository {
	return &CardRepository{
		db:     db,
		cache:  newCardCache(10000), // Cache up to 10k cards
		logger: logger,
	}
}

// GetByID retrieves a card by ID
func (r *CardRepository) GetByID(ctx context.Context, id int64) (*Card, error) {
	// Check cache first
	if card, ok := r.cache.get(fmt.Sprintf("id:%d", id)); ok {
		return card, nil
	}

	// Query scryfall_cards table (primary data source)
	// Map Scryfall fields to Card struct fields
	query := `
		SELECT ('x' || REPLACE(substring(id::text, 1, 18), '-', ''))::bit(64)::bigint as id,
		       COALESCE(collector_number, ''),
		       set_code,
		       name,
		       COALESCE(type_line, ''),
		       COALESCE(mana_cost, ''),
		       COALESCE(power, ''),
		       COALESCE(toughness, ''),
		       COALESCE(oracle_text, ''),
		       NULL::text as flavor_text,
		       NULL::text as original_text,
		       NULL::text as original_type,
		       CASE WHEN collector_number ~ '^[0-9]+$' THEN collector_number::bigint ELSE NULL END as cn,
		       name as card_name,
		       COALESCE(rarity, ''),
		       created_at
		FROM scryfall_cards
		WHERE ('x' || REPLACE(substring(id::text, 1, 18), '-', ''))::bit(64)::bigint = $1
		  AND lang = 'en'
		LIMIT 1
	`

	card := &Card{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&card.ID, &card.CardNumber, &card.SetCode, &card.Name, &card.CardType,
		&card.ManaCost, &card.Power, &card.Toughness, &card.RulesText,
		&card.FlavorText, &card.OriginalText, &card.OriginalType, &card.CN,
		&card.CardName, &card.Rarity, &card.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	// Cache the result
	r.cache.set(fmt.Sprintf("id:%d", id), card)

	return card, nil
}

// GetByName retrieves cards by name
// Handles both single-faced cards and multi-faced cards (DFC, split, etc.)
func (r *CardRepository) GetByName(ctx context.Context, name string) ([]*Card, error) {
	// Query scryfall_cards table (primary data source)
	query := `
		SELECT ('x' || REPLACE(substring(id::text, 1, 18), '-', ''))::bit(64)::bigint as id,
		       COALESCE(collector_number, ''),
		       set_code,
		       name,
		       COALESCE(type_line, ''),
		       COALESCE(mana_cost, ''),
		       COALESCE(power, ''),
		       COALESCE(toughness, ''),
		       COALESCE(oracle_text, ''),
		       NULL::text as flavor_text,
		       NULL::text as original_text,
		       NULL::text as original_type,
		       CASE WHEN collector_number ~ '^[0-9]+$' THEN collector_number::bigint ELSE NULL END as cn,
		       name as card_name,
		       COALESCE(rarity, ''),
		       created_at
		FROM scryfall_cards
		WHERE lang = 'en'
		  AND (name = $1 OR name LIKE $1 || ' //%')
		ORDER BY released_at DESC, set_code, collector_number
	`

	rows, err := r.db.Pool.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to query cards: %w", err)
	}
	defer rows.Close()

	cards := make([]*Card, 0)
	for rows.Next() {
		card := &Card{}
		err := rows.Scan(
			&card.ID, &card.CardNumber, &card.SetCode, &card.Name, &card.CardType,
			&card.ManaCost, &card.Power, &card.Toughness, &card.RulesText,
			&card.FlavorText, &card.OriginalText, &card.OriginalType, &card.CN,
			&card.CardName, &card.Rarity, &card.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// GetByNameCaseInsensitive retrieves cards by name (case-insensitive)
// Handles both single-faced cards and multi-faced cards (DFC, split, etc.)
func (r *CardRepository) GetByNameCaseInsensitive(ctx context.Context, name string) ([]*Card, error) {
	// Trim and normalize the input name
	name = strings.TrimSpace(name)
	// Normalize Unicode apostrophes to ASCII apostrophe
	// U+2019 (') -> U+0027 (')
	// U+2018 (') -> U+0027 (')
	name = strings.ReplaceAll(name, "'", "'")
	name = strings.ReplaceAll(name, "'", "'")

	// Query scryfall_cards table (primary data source)
	query := `
		SELECT ('x' || REPLACE(substring(id::text, 1, 18), '-', ''))::bit(64)::bigint as id,
		       COALESCE(collector_number, ''),
		       set_code,
		       name,
		       COALESCE(type_line, ''),
		       COALESCE(mana_cost, ''),
		       COALESCE(power, ''),
		       COALESCE(toughness, ''),
		       COALESCE(oracle_text, ''),
		       NULL::text as flavor_text,
		       NULL::text as original_text,
		       NULL::text as original_type,
		       CASE WHEN collector_number ~ '^[0-9]+$' THEN collector_number::bigint ELSE NULL END as cn,
		       name as card_name,
		       COALESCE(rarity, ''),
		       created_at
		FROM scryfall_cards
		WHERE lang = 'en'
		  AND (LOWER(TRIM(name)) = LOWER(TRIM($1))
		       OR LOWER(name) LIKE LOWER($1) || ' //%')
		ORDER BY released_at DESC, set_code, collector_number
	`

	rows, err := r.db.Pool.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to query cards: %w", err)
	}
	defer rows.Close()

	cards := make([]*Card, 0)
	for rows.Next() {
		card := &Card{}
		err := rows.Scan(
			&card.ID, &card.CardNumber, &card.SetCode, &card.Name, &card.CardType,
			&card.ManaCost, &card.Power, &card.Toughness, &card.RulesText,
			&card.FlavorText, &card.OriginalText, &card.OriginalType, &card.CN,
			&card.CardName, &card.Rarity, &card.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// SearchByName performs a full-text search on card names
func (r *CardRepository) SearchByName(ctx context.Context, searchTerm string, limit int) ([]*Card, error) {
	// Query scryfall_cards table (primary data source)
	query := `
		SELECT ('x' || REPLACE(substring(id::text, 1, 18), '-', ''))::bit(64)::bigint as id,
		       COALESCE(collector_number, ''),
		       set_code,
		       name,
		       COALESCE(type_line, ''),
		       COALESCE(mana_cost, ''),
		       COALESCE(power, ''),
		       COALESCE(toughness, ''),
		       COALESCE(oracle_text, ''),
		       NULL::text as flavor_text,
		       NULL::text as original_text,
		       NULL::text as original_type,
		       CASE WHEN collector_number ~ '^[0-9]+$' THEN collector_number::bigint ELSE NULL END as cn,
		       name as card_name,
		       COALESCE(rarity, ''),
		       created_at
		FROM scryfall_cards
		WHERE lang = 'en'
		  AND name ILIKE $1
		ORDER BY name, released_at DESC
		LIMIT $2
	`

	rows, err := r.db.Pool.Query(ctx, query, "%"+searchTerm+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search cards: %w", err)
	}
	defer rows.Close()

	cards := make([]*Card, 0)
	for rows.Next() {
		card := &Card{}
		err := rows.Scan(
			&card.ID, &card.CardNumber, &card.SetCode, &card.Name, &card.CardType,
			&card.ManaCost, &card.Power, &card.Toughness, &card.RulesText,
			&card.FlavorText, &card.OriginalText, &card.OriginalType, &card.CN,
			&card.CardName, &card.Rarity, &card.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// GetBySetCode retrieves all cards from a set
func (r *CardRepository) GetBySetCode(ctx context.Context, setCode string) ([]*Card, error) {
	// Query scryfall_cards table (primary data source)
	query := `
		SELECT ('x' || REPLACE(substring(id::text, 1, 18), '-', ''))::bit(64)::bigint as id,
		       COALESCE(collector_number, ''),
		       set_code,
		       name,
		       COALESCE(type_line, ''),
		       COALESCE(mana_cost, ''),
		       COALESCE(power, ''),
		       COALESCE(toughness, ''),
		       COALESCE(oracle_text, ''),
		       NULL::text as flavor_text,
		       NULL::text as original_text,
		       NULL::text as original_type,
		       CASE WHEN collector_number ~ '^[0-9]+$' THEN collector_number::bigint ELSE NULL END as cn,
		       name as card_name,
		       COALESCE(rarity, ''),
		       created_at
		FROM scryfall_cards
		WHERE lang = 'en'
		  AND set_code = $1
		ORDER BY collector_number
	`

	rows, err := r.db.Pool.Query(ctx, query, setCode)
	if err != nil {
		return nil, fmt.Errorf("failed to query cards: %w", err)
	}
	defer rows.Close()

	cards := make([]*Card, 0)
	for rows.Next() {
		card := &Card{}
		err := rows.Scan(
			&card.ID, &card.CardNumber, &card.SetCode, &card.Name, &card.CardType,
			&card.ManaCost, &card.Power, &card.Toughness, &card.RulesText,
			&card.FlavorText, &card.OriginalText, &card.OriginalType, &card.CN,
			&card.CardName, &card.Rarity, &card.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// Create creates a new card
func (r *CardRepository) Create(ctx context.Context, card *Card) error {
	query := `
		INSERT INTO cards (card_number, set_code, name, card_type, mana_cost,
		                   power, toughness, rules_text, flavor_text, original_text,
		                   original_type, cn, card_name, rarity)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		card.CardNumber, card.SetCode, card.Name, card.CardType, card.ManaCost,
		card.Power, card.Toughness, card.RulesText,
		getNullStringValue(card.FlavorText),
		getNullStringValue(card.OriginalText),
		getNullStringValue(card.OriginalType),
		getNullInt64Value(card.CN), getNullStringValue(card.CardName), card.Rarity,
	).Scan(&card.ID, &card.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create card: %w", err)
	}

	return nil
}

// getNullStringValue returns the string value from sql.NullString, or nil if not valid
func getNullStringValue(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}

// getNullInt64Value returns the int64 value from sql.NullInt64, or nil if not valid
func getNullInt64Value(ni sql.NullInt64) interface{} {
	if ni.Valid {
		return ni.Int64
	}
	return nil
}

// cardCache is a simple in-memory cache for cards
type cardCache struct {
	items   map[string]*cacheItem
	maxSize int
	mu      sync.RWMutex
}

type cacheItem struct {
	card      *Card
	expiresAt time.Time
}

func newCardCache(maxSize int) *cardCache {
	return &cardCache{
		items:   make(map[string]*cacheItem),
		maxSize: maxSize,
	}
}

func (c *cardCache) get(key string) (*Card, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.card, true
}

func (c *cardCache) set(key string, card *Card) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Simple eviction: remove oldest items if cache is full
	if len(c.items) >= c.maxSize {
		// Remove first item (simple FIFO)
		for k := range c.items {
			delete(c.items, k)
			break
		}
	}

	c.items[key] = &cacheItem{
		card:      card,
		expiresAt: time.Now().Add(24 * time.Hour),
	}
}

func (c *cardCache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*cacheItem)
}
