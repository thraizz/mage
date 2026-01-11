#!/bin/bash
# Import H2 SQL export to PostgreSQL database
# Imports into the existing 'mage' database

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Input SQL (H2 export)
INPUT_SQL_DEFAULT="${PROJECT_ROOT}/data/cards.sql"
INPUT_SQL="${INPUT_SQL:-$INPUT_SQL_DEFAULT}"

MODE="direct"
if [[ "${1:-}" == "--docker" ]]; then
  MODE="docker"
  shift
elif [[ "${1:-}" == "--direct" ]]; then
  MODE="direct"
  shift
fi

# PostgreSQL connection settings
DB_NAME="${DB_NAME:-mage}"
DB_USER="${DB_USER:-mage}"
DB_PASSWORD="${DB_PASSWORD:-mage}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DOCKER_CONTAINER="${DOCKER_CONTAINER:-mage-postgres}"

echo "=== Import to PostgreSQL ==="
echo "Input: $INPUT_SQL"
echo "Database: $DB_NAME"
if [[ "$MODE" == "docker" ]]; then
    echo "Mode: docker"
    echo "Container: $DOCKER_CONTAINER"
else
    echo "Mode: direct"
    echo "Host: $DB_HOST:$DB_PORT"
fi
echo "User: $DB_USER"
echo ""

# Check if SQL file exists
if [ ! -f "$INPUT_SQL" ]; then
    echo "ERROR: SQL file not found: $INPUT_SQL"
    echo "Run first: ./scripts/h2_to_sql.sh"
    exit 1
fi

psql_base_query() {
    local query="$1"
    if [[ "$MODE" == "docker" ]]; then
        docker exec -i "$DOCKER_CONTAINER" env PAGER=cat \
            psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME" -tAc "$query"
    else
        PAGER=cat PGPASSWORD="$DB_PASSWORD" \
            psql -v ON_ERROR_STOP=1 -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "$query"
    fi
}

psql_base_file() {
    local file_path="$1"
    if [[ "$MODE" == "docker" ]]; then
        docker exec -i "$DOCKER_CONTAINER" env PAGER=cat \
            psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME" -f "$file_path"
    else
        PAGER=cat PGPASSWORD="$DB_PASSWORD" \
            psql -v ON_ERROR_STOP=1 -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file_path"
    fi
}

psql_base_stdin() {
    if [[ "$MODE" == "docker" ]]; then
        docker exec -i "$DOCKER_CONTAINER" env PAGER=cat \
            psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME"
    else
        PAGER=cat PGPASSWORD="$DB_PASSWORD" \
            psql -v ON_ERROR_STOP=1 -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME"
    fi
}

# Test connection
echo "Testing database connection..."
if ! psql_base_query "SELECT 1;" > /dev/null 2>&1; then
    echo "ERROR: Cannot connect to PostgreSQL"
    echo "Check connection settings and ensure database exists"
    exit 1
fi
echo "✓ Connected to database"

# Convert H2 SQL to PostgreSQL-compatible SQL
echo "Converting SQL to PostgreSQL format..."
TMP_PG="${INPUT_SQL}.postgres.tmp"
python3 "$SCRIPT_DIR/convert_h2_to_postgres.py" "$INPUT_SQL" "$TMP_PG"

# Import to PostgreSQL
echo "Importing to PostgreSQL..."
if [[ "$MODE" == "docker" ]]; then
    CONTAINER_TMP="/tmp/$(basename "$TMP_PG")"
    docker cp "$TMP_PG" "$DOCKER_CONTAINER:$CONTAINER_TMP" >/dev/null
    psql_base_file "$CONTAINER_TMP" 2>&1 | grep -v "^$" || true
    docker exec -i "$DOCKER_CONTAINER" rm -f "$CONTAINER_TMP" >/dev/null 2>&1 || true
else
    psql_base_file "$TMP_PG" 2>&1 | grep -v "^$" || true
fi

# Cleanup temp file
rm -f "$TMP_PG"

# Verify import
CARD_COUNT=$(psql_base_query "SELECT COUNT(*) FROM card;" 2>/dev/null | xargs || echo "0")
SET_COUNT=$(psql_base_query "SELECT COUNT(*) FROM expansion;" 2>/dev/null | xargs || echo "0")

echo "✓ Import complete"
echo ""
echo "=== Database Statistics ==="
echo "Cards: $CARD_COUNT"
echo "Sets: $SET_COUNT"
echo ""

# Test queries
echo "Sample cards:"
psql_base_stdin >/dev/null 2>&1 <<'SQL' || true
SELECT name, setcode, manacosts FROM card WHERE name LIKE 'Lightning%' LIMIT 5;
SQL
echo ""

# Create indexes for H2 card table
echo "Creating indexes on H2 card table..."
psql_base_stdin <<'EOF'
CREATE INDEX IF NOT EXISTS idx_card_name ON card(name);
CREATE INDEX IF NOT EXISTS idx_card_setcode ON card(setcode);
CREATE INDEX IF NOT EXISTS idx_card_classname ON card(classname);
CREATE INDEX IF NOT EXISTS idx_card_name_set ON card(name, setcode);
CREATE INDEX IF NOT EXISTS idx_card_types ON card(types);
CREATE INDEX IF NOT EXISTS idx_card_manavalue ON card(manavalue);
ANALYZE card;
ANALYZE expansion;
EOF
echo "✓ Indexes created"

# Migrate data to Go server's cards table
echo ""
echo "Migrating data to Go server's cards table..."
psql_base_stdin <<'EOF'
-- Clear existing cards (they're auto-generated placeholders)
TRUNCATE cards CASCADE;

-- Migrate from H2 card table to Go cards table with column mapping
INSERT INTO cards (card_number, set_code, name, card_type, mana_cost, power, toughness, rules_text, rarity, card_class_name)
SELECT 
    cardnumber,
    setcode,
    name,
    types,
    -- Remove @@@ delimiter from mana costs
    REPLACE(manacosts, '@@@', ''),
    power,
    toughness,
    rules,
    rarity,
    classname
FROM card;

-- Update indexes
ANALYZE cards;
EOF

MIGRATED_COUNT=$(psql_base_query "SELECT COUNT(*) FROM cards;" 2>/dev/null | xargs || echo "0")
echo "✓ Migrated $MIGRATED_COUNT cards to Go server format"

echo ""
echo "Database ready!"
echo ""
echo "Test with:"
echo "  PAGER=cat psql -d $DB_NAME -c 'SELECT COUNT(*) FROM cards;'"
echo "  PAGER=cat psql -d $DB_NAME -c \"SELECT name, mana_cost FROM cards WHERE name = 'Lightning Bolt';\""
