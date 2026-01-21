#!/bin/sh
# Docker entrypoint for MAGE server
# This script runs on container startup and:
# 1. Checks if Scryfall card data needs to be imported
# 2. Imports data if needed
# 3. Creates compatibility view
# 4. Starts the MAGE server

set -e

echo "======================================"
echo "MAGE Server - Docker Entrypoint"
echo "======================================"
echo ""

# Database connection details (from environment variables)
DB_HOST="${DATABASE_HOST:-postgres}"
DB_PORT="${DATABASE_PORT:-5432}"
DB_NAME="${DATABASE_NAME:-mage}"
DB_USER="${DATABASE_USER:-mage}"
DB_PASSWORD="${DATABASE_PASSWORD:-mage}"
DB_SSLMODE="${DATABASE_SSLMODE:-disable}"

# Construct database URL
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

# Wait for database to be ready
echo "[1/4] Waiting for database to be ready..."
max_attempts=30
attempt=0
while [ $attempt -lt $max_attempts ]; do
    if pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" >/dev/null 2>&1; then
        echo "✓ Database is ready"
        break
    fi
    attempt=$((attempt + 1))
    echo "  Waiting... ($attempt/$max_attempts)"
    sleep 2
done

if [ $attempt -eq $max_attempts ]; then
    echo "ERROR: Database did not become ready in time"
    exit 1
fi
echo ""

# Check if scryfall_cards table exists and has data
echo "[2/4] Checking Scryfall card data..."
CARD_COUNT=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc \
    "SELECT COUNT(*) FROM scryfall_cards WHERE lang='en';" 2>/dev/null || echo "0")

if [ "$CARD_COUNT" = "0" ] || [ -z "$CARD_COUNT" ]; then
    echo "⚠ No Scryfall card data found. Starting import..."
    echo ""

    # Find Scryfall data file
    SCRYFALL_FILE=""
    if [ -L "/app/data/scryfall-all-cards-latest.json" ]; then
        # Follow symlink
        SCRYFALL_FILE="/app/data/scryfall-all-cards-latest.json"
    elif [ -f "/app/data/scryfall-all-cards-latest.json" ]; then
        SCRYFALL_FILE="/app/data/scryfall-all-cards-latest.json"
    else
        # Find most recent file
        SCRYFALL_FILE=$(ls -t /app/data/scryfall-all-cards-*.json 2>/dev/null | head -1 || echo "")
    fi

    if [ -z "$SCRYFALL_FILE" ] || [ ! -f "$SCRYFALL_FILE" ]; then
        echo "ERROR: No Scryfall data file found in /app/data/"
        echo "Please ensure a Scryfall JSON file exists in the data directory"
        echo ""
        echo "You can download Scryfall data with:"
        echo "  ./scripts/download_scryfall_bulk.sh"
        echo ""
        echo "Starting server without card data (you can import later)..."
    else
        FILE_SIZE=$(du -h "$SCRYFALL_FILE" | cut -f1)
        echo "  Found: $(basename "$SCRYFALL_FILE") ($FILE_SIZE)"
        echo "  This will take 2-3 minutes. Please wait..."
        echo ""

        # Run import
        /app/scryfall-import \
            --input="$SCRYFALL_FILE" \
            --lang=en \
            --skip-tokens=true \
            --batch=1000 \
            --db="$DATABASE_URL"

        echo ""
        echo "✓ Scryfall import complete"
    fi
else
    echo "✓ Found $CARD_COUNT English Scryfall cards"
fi
echo ""

# Create or update compatibility view
echo "[3/4] Creating compatibility view..."
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<'SQL'
-- Backup old cards table if it exists and isn't already backed up
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'cards' AND table_type = 'BASE TABLE') THEN
        DROP TABLE IF EXISTS cards_xmage_backup CASCADE;
        ALTER TABLE cards RENAME TO cards_xmage_backup;
    END IF;
END
$$;

-- Create compatibility view
DROP VIEW IF EXISTS cards CASCADE;
CREATE VIEW cards AS
SELECT
    (('x' || REPLACE(substring(id::text, 1, 18), '-', ''))::bit(64)::bigint) as id,
    collector_number as card_number,
    set_code,
    name,
    type_line as card_type,
    COALESCE(mana_cost, '') as mana_cost,
    COALESCE(power, '') as power,
    COALESCE(toughness, '') as toughness,
    COALESCE(oracle_text, '') as rules_text,
    '' as flavor_text,
    '' as original_text,
    '' as original_type,
    0::bigint as cn,
    name as card_name,
    rarity,
    COALESCE(name, '') as card_class_name,
    created_at
FROM scryfall_cards
WHERE lang = 'en'
  AND layout NOT IN ('token', 'art_series', 'emblem', 'double_faced_token')
  AND digital = false;
SQL

echo "✓ Compatibility view created"
echo ""

# Verify setup
echo "[4/4] Verifying setup..."
FINAL_COUNT=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc \
    "SELECT COUNT(*) FROM cards;" 2>/dev/null || echo "0")
echo "✓ Cards available through compatibility view: $FINAL_COUNT"
echo ""

echo "======================================"
echo "Starting MAGE Server"
echo "======================================"
echo ""

# Start the server with all passed arguments
exec /app/mage-server "$@"
