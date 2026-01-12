package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// ScryfallCard represents a card from Scryfall bulk data
type ScryfallCard struct {
	ID              string            `json:"id"`
	OracleID        string            `json:"oracle_id"`
	Name            string            `json:"name"`
	Lang            string            `json:"lang"`
	ReleasedAt      string            `json:"released_at"`
	Layout          string            `json:"layout"`
	ImageURIs       *ImageURIs        `json:"image_uris,omitempty"`
	ManaCost        string            `json:"mana_cost,omitempty"`
	CMC             float64           `json:"cmc"`
	TypeLine        string            `json:"type_line"`
	OracleText      string            `json:"oracle_text,omitempty"`
	Power           *string           `json:"power,omitempty"`
	Toughness       *string           `json:"toughness,omitempty"`
	Loyalty         *string           `json:"loyalty,omitempty"`
	Defense         *string           `json:"defense,omitempty"`
	Colors          []string          `json:"colors"`
	ColorIdentity   []string          `json:"color_identity"`
	Keywords        []string          `json:"keywords"`
	Legalities      map[string]string `json:"legalities"`
	Set             string            `json:"set"`
	SetName         string            `json:"set_name"`
	CollectorNumber string            `json:"collector_number"`
	Rarity          string            `json:"rarity"`
	Digital         bool              `json:"digital"`
	Reprint         bool              `json:"reprint"`
	Promo           bool              `json:"promo"`
	CardFaces       []CardFace        `json:"card_faces,omitempty"`
	EDHRECRank      *int              `json:"edhrec_rank,omitempty"`
}

type ImageURIs struct {
	Small      string `json:"small"`
	Normal     string `json:"normal"`
	Large      string `json:"large"`
	PNG        string `json:"png"`
	ArtCrop    string `json:"art_crop"`
	BorderCrop string `json:"border_crop"`
}

type CardFace struct {
	Name       string     `json:"name"`
	ManaCost   string     `json:"mana_cost,omitempty"`
	TypeLine   string     `json:"type_line"`
	OracleText string     `json:"oracle_text,omitempty"`
	Power      *string    `json:"power,omitempty"`
	Toughness  *string    `json:"toughness,omitempty"`
	Loyalty    *string    `json:"loyalty,omitempty"`
	Colors     []string   `json:"colors,omitempty"`
	ImageURIs  *ImageURIs `json:"image_uris,omitempty"`
}

// ScryfallSet represents a set from Scryfall
type ScryfallSet struct {
	ID         string `json:"id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	SetType    string `json:"set_type"`
	ReleasedAt string `json:"released_at,omitempty"`
	CardCount  int    `json:"card_count"`
	Digital    bool   `json:"digital"`
	IconSVGURI string `json:"icon_svg_uri,omitempty"`
}

func main() {
	var (
		inputFile  = flag.String("input", "", "Path to Scryfall bulk data JSON file")
		dbURL      = flag.String("db", "", "PostgreSQL connection URL (or use DATABASE_URL env)")
		batchSize  = flag.Int("batch", 1000, "Number of cards to insert per batch")
		langFilter = flag.String("lang", "en", "Language filter (empty for all)")
		skipTokens = flag.Bool("skip-tokens", true, "Skip tokens and art cards")
	)
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: --input flag is required")
		fmt.Println("Usage: scryfall-import --input=/path/to/all-cards.json")
		os.Exit(1)
	}

	// Setup logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Get database URL
	databaseURL := *dbURL
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		databaseURL = "postgres://mage:mage@localhost:5432/mage?sslmode=disable"
	}

	ctx := context.Background()

	// Connect to database
	logger.Info("Connecting to database", zap.String("url", maskPassword(databaseURL)))
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}
	logger.Info("Database connection established")

	// Run import
	importer := &Importer{
		pool:       pool,
		logger:     logger,
		batchSize:  *batchSize,
		langFilter: *langFilter,
		skipTokens: *skipTokens,
	}

	startTime := time.Now()
	stats, err := importer.Import(ctx, *inputFile)
	if err != nil {
		logger.Fatal("Import failed", zap.Error(err))
	}

	duration := time.Since(startTime)
	logger.Info("Import completed",
		zap.Int("total_cards", stats.TotalCards),
		zap.Int("imported_cards", stats.ImportedCards),
		zap.Int("skipped_cards", stats.SkippedCards),
		zap.Int("failed_cards", stats.FailedCards),
		zap.Duration("duration", duration),
		zap.Float64("cards_per_second", float64(stats.ImportedCards)/duration.Seconds()),
	)
}

type Importer struct {
	pool       *pgxpool.Pool
	logger     *zap.Logger
	batchSize  int
	langFilter string
	skipTokens bool
}

type ImportStats struct {
	TotalCards    int
	ImportedCards int
	SkippedCards  int
	FailedCards   int
}

func (imp *Importer) Import(ctx context.Context, filename string) (*ImportStats, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file size for progress tracking
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	imp.logger.Info("Starting import",
		zap.String("file", filename),
		zap.Float64("size_gb", float64(fileSize)/(1024*1024*1024)),
	)

	stats := &ImportStats{}

	// Use streaming JSON decoder
	dec := json.NewDecoder(file)

	// Read opening bracket
	if _, err := dec.Token(); err != nil {
		return nil, fmt.Errorf("failed to read opening bracket: %w", err)
	}

	batch := make([]*ScryfallCard, 0, imp.batchSize)
	lastProgress := time.Now()

	// Read array elements
	for dec.More() {
		var card ScryfallCard
		if err := dec.Decode(&card); err != nil {
			imp.logger.Warn("Failed to decode card", zap.Error(err))
			stats.FailedCards++
			continue
		}

		stats.TotalCards++

		// Apply filters
		if imp.langFilter != "" && card.Lang != imp.langFilter {
			stats.SkippedCards++
			continue
		}

		if imp.skipTokens && shouldSkipLayout(card.Layout) {
			stats.SkippedCards++
			continue
		}

		batch = append(batch, &card)

		// Insert batch when full
		if len(batch) >= imp.batchSize {
			if err := imp.insertBatch(ctx, batch); err != nil {
				imp.logger.Error("Failed to insert batch", zap.Error(err))
				stats.FailedCards += len(batch)
			} else {
				stats.ImportedCards += len(batch)
			}
			batch = batch[:0]

			// Progress reporting every 5 seconds
			if time.Since(lastProgress) > 5*time.Second {
				imp.logger.Info("Progress",
					zap.Int("total", stats.TotalCards),
					zap.Int("imported", stats.ImportedCards),
					zap.Int("skipped", stats.SkippedCards),
				)
				lastProgress = time.Now()
			}
		}
	}

	// Insert remaining cards
	if len(batch) > 0 {
		if err := imp.insertBatch(ctx, batch); err != nil {
			imp.logger.Error("Failed to insert final batch", zap.Error(err))
			stats.FailedCards += len(batch)
		} else {
			stats.ImportedCards += len(batch)
		}
	}

	return stats, nil
}

func (imp *Importer) insertBatch(ctx context.Context, cards []*ScryfallCard) error {
	tx, err := imp.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, card := range cards {
		// Prepare image URIs
		var imgSmall, imgNormal, imgLarge, imgPNG *string
		if card.ImageURIs != nil {
			if card.ImageURIs.Small != "" {
				imgSmall = &card.ImageURIs.Small
			}
			if card.ImageURIs.Normal != "" {
				imgNormal = &card.ImageURIs.Normal
			}
			if card.ImageURIs.Large != "" {
				imgLarge = &card.ImageURIs.Large
			}
			if card.ImageURIs.PNG != "" {
				imgPNG = &card.ImageURIs.PNG
			}
		}

		// Convert card_faces to JSON
		var cardFacesJSON []byte
		if len(card.CardFaces) > 0 {
			cardFacesJSON, _ = json.Marshal(card.CardFaces)
		}

		// Convert legalities to JSON
		var legalitiesJSON []byte
		if len(card.Legalities) > 0 {
			legalitiesJSON, _ = json.Marshal(card.Legalities)
		}

		// Upsert card
		_, err := tx.Exec(ctx, `
			INSERT INTO scryfall_cards (
				id, oracle_id, name, lang,
				set_code, set_name, collector_number,
				layout, type_line, oracle_text, mana_cost, cmc,
				power, toughness, loyalty, defense,
				colors, color_identity, rarity, keywords,
				card_faces, legalities,
				image_uri_small, image_uri_normal, image_uri_large, image_uri_png,
				released_at, reprint, digital, promo,
				edhrec_rank
			) VALUES (
				$1, $2, $3, $4,
				$5, $6, $7,
				$8, $9, $10, $11, $12,
				$13, $14, $15, $16,
				$17, $18, $19, $20,
				$21, $22,
				$23, $24, $25, $26,
				$27, $28, $29, $30,
				$31
			)
			ON CONFLICT (id) DO UPDATE SET
				oracle_id = EXCLUDED.oracle_id,
				name = EXCLUDED.name,
				type_line = EXCLUDED.type_line,
				oracle_text = EXCLUDED.oracle_text,
				mana_cost = EXCLUDED.mana_cost,
				cmc = EXCLUDED.cmc,
				power = EXCLUDED.power,
				toughness = EXCLUDED.toughness,
				loyalty = EXCLUDED.loyalty,
				defense = EXCLUDED.defense,
				colors = EXCLUDED.colors,
				color_identity = EXCLUDED.color_identity,
				keywords = EXCLUDED.keywords,
				card_faces = EXCLUDED.card_faces,
				legalities = EXCLUDED.legalities,
				image_uri_small = EXCLUDED.image_uri_small,
				image_uri_normal = EXCLUDED.image_uri_normal,
				image_uri_large = EXCLUDED.image_uri_large,
				image_uri_png = EXCLUDED.image_uri_png,
				updated_at = CURRENT_TIMESTAMP
		`,
			parseUUID(card.ID), parseUUID(card.OracleID), card.Name, card.Lang,
			card.Set, card.SetName, card.CollectorNumber,
			card.Layout, card.TypeLine, nullIfEmpty(card.OracleText), nullIfEmpty(card.ManaCost), card.CMC,
			card.Power, card.Toughness, card.Loyalty, card.Defense,
			card.Colors, card.ColorIdentity, card.Rarity, card.Keywords,
			cardFacesJSON, legalitiesJSON,
			imgSmall, imgNormal, imgLarge, imgPNG,
			parseDate(card.ReleasedAt), card.Reprint, card.Digital, card.Promo,
			card.EDHRECRank,
		)

		if err != nil {
			return fmt.Errorf("failed to insert card %s: %w", card.Name, err)
		}
	}

	return tx.Commit(ctx)
}

func shouldSkipLayout(layout string) bool {
	skipLayouts := []string{
		"token", "art_series", "emblem", "double_faced_token",
	}
	for _, skip := range skipLayouts {
		if layout == skip {
			return true
		}
	}
	return false
}

func parseUUID(s string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return u
}

func parseDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func maskPassword(url string) string {
	// Simple password masking for logs
	if idx := strings.Index(url, "://"); idx >= 0 {
		if idx2 := strings.Index(url[idx+3:], "@"); idx2 >= 0 {
			return url[:idx+3] + "***:***" + url[idx+3+idx2:]
		}
	}
	return url
}
