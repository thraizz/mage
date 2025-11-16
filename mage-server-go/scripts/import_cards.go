package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CardImport represents a card record from the CSV export
type CardImport struct {
	Name              string
	SetCode           string
	CardNumber        string
	ClassName         string
	Power             string
	Toughness         string
	StartingLoyalty   string
	StartingDefense   string
	ManaValue         int
	Rarity            string
	Types             string
	Subtypes          string
	Supertypes        string
	ManaCosts         string
	Rules             string
	Black             bool
	Blue              bool
	Green             bool
	Red               bool
	White             bool
	FrameColor        string
	FrameStyle        string
	VariousArt        bool
}

func main() {
	ctx := context.Background()

	// Get CSV file path from args or use default
	csvPath := "data/cards_export.csv"
	if len(os.Args) > 1 {
		csvPath = os.Args[1]
	}

	// Get absolute path
	absPath, err := filepath.Abs(csvPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	fmt.Println("=== MAGE Card Data Import ===")
	fmt.Printf("CSV file: %s\n", absPath)

	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("CSV file not found: %s\nRun: ./scripts/export_java_cards.sh", absPath)
	}

	// Connect to PostgreSQL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/mage?sslmode=disable"
	}

	fmt.Printf("Connecting to database...\n")
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("✓ Database connection established")

	// Read CSV file
	file, err := os.Open(absPath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV: %v", err)
	}

	if len(records) < 2 {
		log.Fatal("CSV file is empty or has no data rows")
	}

	fmt.Printf("Found %d cards in CSV\n", len(records)-1) // -1 for header

	// Parse and import cards
	cards := make([]*CardImport, 0, len(records)-1)
	for i, record := range records[1:] { // Skip header
		if len(record) < 23 {
			log.Printf("Warning: Skipping row %d - insufficient columns", i+2)
			continue
		}

		card := &CardImport{
			Name:            record[0],
			SetCode:         record[1],
			CardNumber:      record[2],
			ClassName:       record[3],
			Power:           record[4],
			Toughness:       record[5],
			StartingLoyalty: record[6],
			StartingDefense: record[7],
			Rarity:          record[9],
			Types:           record[10],
			Subtypes:        record[11],
			Supertypes:      record[12],
			ManaCosts:       record[13],
			Rules:           record[14],
			FrameColor:      record[19],
			FrameStyle:      record[20],
		}

		// Parse integer fields
		if manaValue, err := strconv.Atoi(record[8]); err == nil {
			card.ManaValue = manaValue
		}

		// Parse boolean fields
		card.Black = parseBool(record[15])
		card.Blue = parseBool(record[16])
		card.Green = parseBool(record[17])
		card.Red = parseBool(record[18])
		card.White = parseBool(record[19])
		card.VariousArt = parseBool(record[22])

		cards = append(cards, card)
	}

	fmt.Printf("Parsed %d valid cards\n", len(cards))

	// Check if cards already exist
	var existingCount int64
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM cards").Scan(&existingCount)
	if err != nil {
		log.Fatalf("Failed to check existing cards: %v", err)
	}

	if existingCount > 0 {
		fmt.Printf("Warning: Database already contains %d cards\n", existingCount)
		fmt.Print("Do you want to clear and reimport? (yes/no): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) == "yes" {
			fmt.Println("Clearing existing cards...")
			_, err = pool.Exec(ctx, "TRUNCATE cards RESTART IDENTITY CASCADE")
			if err != nil {
				log.Fatalf("Failed to clear cards: %v", err)
			}
			fmt.Println("✓ Existing cards cleared")
		} else {
			fmt.Println("Import cancelled")
			return
		}
	}

	// Import cards in batches
	fmt.Println("Importing cards...")
	batchSize := 1000
	imported := 0
	failed := 0

	startTime := time.Now()

	for i := 0; i < len(cards); i += batchSize {
		end := i + batchSize
		if end > len(cards) {
			end = len(cards)
		}

		batch := cards[i:end]
		
		// Begin transaction
		tx, err := pool.Begin(ctx)
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			failed += len(batch)
			continue
		}

		for _, card := range batch {
			// Combine types, subtypes, supertypes into card_type field
			cardType := buildCardType(card.Types, card.Subtypes, card.Supertypes)

			_, err := tx.Exec(ctx, `
				INSERT INTO cards (
					card_number, set_code, name, card_type, mana_cost,
					power, toughness, rules_text, flavor_text, original_text,
					original_type, cn, card_name, rarity, card_class_name
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			`,
				card.CardNumber,
				card.SetCode,
				card.Name,
				cardType,
				card.ManaCosts,
				card.Power,
				card.Toughness,
				card.Rules,
				"", // flavor_text (not in export)
				"", // original_text (not in export)
				"", // original_type (not in export)
				0,  // cn (not in export)
				card.Name,
				card.Rarity,
				card.ClassName,
			)

			if err != nil {
				log.Printf("Failed to insert card %s: %v", card.Name, err)
				failed++
			} else {
				imported++
			}
		}

		// Commit transaction
		if err := tx.Commit(ctx); err != nil {
			log.Printf("Failed to commit batch: %v", err)
			tx.Rollback(ctx)
			failed += len(batch)
		}

		// Progress update
		if (i+batchSize)%5000 == 0 || end == len(cards) {
			fmt.Printf("Progress: %d/%d cards imported\n", imported, len(cards))
		}
	}

	duration := time.Since(startTime)

	fmt.Println("\n=== Import Complete ===")
	fmt.Printf("✓ Successfully imported: %d cards\n", imported)
	if failed > 0 {
		fmt.Printf("✗ Failed to import: %d cards\n", failed)
	}
	fmt.Printf("Time taken: %s\n", duration)
	fmt.Printf("Rate: %.0f cards/second\n", float64(imported)/duration.Seconds())

	// Verify import
	var finalCount int64
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM cards").Scan(&finalCount)
	if err == nil {
		fmt.Printf("\nTotal cards in database: %d\n", finalCount)
	}

	fmt.Println("\nNext steps:")
	fmt.Println("  1. Verify: PAGER=cat psql -d mage -c 'SELECT COUNT(*) FROM cards;'")
	fmt.Println("  2. Test query: PAGER=cat psql -d mage -c \"SELECT name, set_code, mana_cost FROM cards LIMIT 10;\"")
	fmt.Println("  3. Implement ability system (see CARD_DATA_ARCHITECTURE.md)")
}

func parseBool(s string) bool {
	return strings.ToLower(s) == "true" || s == "1"
}

func buildCardType(types, subtypes, supertypes string) string {
	parts := []string{}

	if supertypes != "" {
		parts = append(parts, supertypes)
	}

	if types != "" {
		parts = append(parts, types)
	}

	result := strings.Join(parts, " ")

	if subtypes != "" {
		result += " — " + subtypes
	}

	return result
}
