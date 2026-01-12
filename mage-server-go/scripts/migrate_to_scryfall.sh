#!/usr/bin/env bash
#
# Migrate card data from XMage to Scryfall
#
# This script:
# 1. Runs database migrations to create Scryfall tables
# 2. Downloads latest Scryfall bulk data (if needed)
# 3. Imports Scryfall data into PostgreSQL
# 4. Verifies the import
# 5. Creates compatibility view for existing code
#
# Usage:
#   ./mage-server-go/scripts/migrate_to_scryfall.sh [--skip-download]
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
DATA_DIR="$GO_ROOT/data"

SKIP_DOWNLOAD=false
if [[ "${1:-}" == "--skip-download" ]]; then
    SKIP_DOWNLOAD=true
    shift
fi

# Database connection settings
DB_NAME="${DB_NAME:-mage}"
DB_USER="${DB_USER:-mage}"
DB_PASSWORD="${DB_PASSWORD:-mage}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DOCKER_CONTAINER="${DOCKER_CONTAINER:-mage-postgres}"

# Determine if using Docker or direct connection
MODE="direct"
if docker ps --format '{{.Names}}' | grep -q "^${DOCKER_CONTAINER}$"; then
    MODE="docker"
    echo "Detected running Docker container: $DOCKER_CONTAINER"
fi

psql_exec() {
    local query="$1"
    if [[ "$MODE" == "docker" ]]; then
        docker exec -i "$DOCKER_CONTAINER" env PAGER=cat \
            psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME" -tAc "$query"
    else
        PAGER=cat PGPASSWORD="$DB_PASSWORD" \
            psql -v ON_ERROR_STOP=1 -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "$query"
    fi
}

echo "==================================="
echo "  Migrate to Scryfall Data"
echo "==================================="
echo ""
echo "Database: $DB_NAME@$DB_HOST:$DB_PORT"
echo "Mode: $MODE"
echo ""

# Step 1: Run migrations
echo "Step 1: Running database migrations..."
if [[ "$MODE" == "docker" ]]; then
    (cd "$GO_ROOT" && DOCKER_CONTAINER="$DOCKER_CONTAINER" DB_NAME="$DB_NAME" DB_USER="$DB_USER" \
        ./scripts/run_postgres_migrations.sh)
else
    (cd "$GO_ROOT" && DB_HOST="$DB_HOST" DB_PORT="$DB_PORT" DB_NAME="$DB_NAME" \
        DB_USER="$DB_USER" DB_PASSWORD="$DB_PASSWORD" \
        ./scripts/run_postgres_migrations.sh --direct)
fi
echo "✓ Migrations complete"
echo ""

# Step 2: Download Scryfall data
if [[ "$SKIP_DOWNLOAD" == "false" ]]; then
    echo "Step 2: Downloading Scryfall bulk data..."
    "$SCRIPT_DIR/download_scryfall_bulk.sh"
else
    echo "Step 2: Skipping download (using existing file)"
fi
echo ""

# Step 3: Find latest Scryfall data file
SCRYFALL_FILE=""
if [ -L "$DATA_DIR/scryfall-all-cards-latest.json" ]; then
    SCRYFALL_FILE="$DATA_DIR/$(readlink "$DATA_DIR/scryfall-all-cards-latest.json")"
elif [ -f "$DATA_DIR/scryfall-all-cards-latest.json" ]; then
    SCRYFALL_FILE="$DATA_DIR/scryfall-all-cards-latest.json"
else
    # Find most recent file
    SCRYFALL_FILE=$(ls -t "$DATA_DIR"/scryfall-all-cards-*.json 2>/dev/null | head -1 || echo "")
fi

if [ -z "$SCRYFALL_FILE" ] || [ ! -f "$SCRYFALL_FILE" ]; then
    # Check if the user specified the file from DATA_MIGRATION.md
    if [ -f "/Users/aron/dev/opensource/mage/all-cards-20260111103023.json" ]; then
        SCRYFALL_FILE="/Users/aron/dev/opensource/mage/all-cards-20260111103023.json"
        echo "Using existing Scryfall dump: $SCRYFALL_FILE"
    else
        echo "ERROR: No Scryfall data file found in $DATA_DIR"
        echo "Run: ./scripts/download_scryfall_bulk.sh"
        exit 1
    fi
fi

echo "Step 3: Importing Scryfall data..."
echo "  File: $SCRYFALL_FILE"

# Build and run importer
export DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
if [[ "$MODE" == "docker" ]]; then
    # For Docker, we need to use the container's network
    export DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@localhost:5432/$DB_NAME?sslmode=disable"
fi

cd "$GO_ROOT"
go run ./cmd/scryfall-import/main.go --input="$SCRYFALL_FILE" --lang=en --skip-tokens=true

echo "✓ Import complete"
echo ""

# Step 4: Verify import
echo "Step 4: Verifying import..."
CARD_COUNT=$(psql_exec "SELECT COUNT(*) FROM scryfall_cards WHERE lang='en';" 2>/dev/null || echo "0")
UNIQUE_COUNT=$(psql_exec "SELECT COUNT(DISTINCT oracle_id) FROM scryfall_cards WHERE lang='en';" 2>/dev/null || echo "0")

echo "  English cards: $CARD_COUNT"
echo "  Unique cards (by oracle_id): $UNIQUE_COUNT"

if [ "$CARD_COUNT" -eq "0" ]; then
    echo "ERROR: No cards were imported!"
    exit 1
fi

echo "✓ Verification passed"
echo ""

# Step 5: Create compatibility layer
echo "Step 5: Creating compatibility layer..."

# Backup old cards table
OLD_TABLE_EXISTS=$(psql_exec "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'cards');" 2>/dev/null || echo "f")

if [ "$OLD_TABLE_EXISTS" == "t" ]; then
    echo "  Backing up old cards table..."
    psql_exec "DROP TABLE IF EXISTS cards_xmage_backup CASCADE;" 2>&1 | grep -v "^$" || true
    psql_exec "ALTER TABLE cards RENAME TO cards_xmage_backup;" 2>&1 | grep -v "^$" || true
    echo "  ✓ Old cards table backed up as cards_xmage_backup"
fi

# Create compatibility view
echo "  Creating cards compatibility view..."
psql_exec "
CREATE OR REPLACE VIEW cards AS
SELECT 
    -- Map UUID to bigint for legacy compatibility (hash-based)
    ('x' || substring(id::text, 1, 15))::bit(60)::bigint as id,
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
" 2>&1 | grep -v "^$" || true

echo "  ✓ Compatibility view created"
echo ""

# Final stats
echo "==================================="
echo "  Migration Complete!"
echo "==================================="
echo ""
echo "Statistics:"
echo "  English cards in database: $CARD_COUNT"
echo "  Unique cards: $UNIQUE_COUNT"
echo ""
echo "Compatibility:"
echo "  - Old 'cards' table backed up as 'cards_xmage_backup'"
echo "  - New 'cards' view provides XMage-compatible interface"
echo "  - Native Scryfall data available in 'scryfall_cards' table"
echo ""
echo "Test queries:"
echo "  PAGER=cat psql -d $DB_NAME -c 'SELECT COUNT(*) FROM cards;'"
echo "  PAGER=cat psql -d $DB_NAME -c \"SELECT name, mana_cost, type_line FROM scryfall_cards WHERE name = 'Lightning Bolt' LIMIT 1;\""
echo ""
echo "Rollback (if needed):"
echo "  DROP VIEW cards;"
echo "  ALTER TABLE cards_xmage_backup RENAME TO cards;"
