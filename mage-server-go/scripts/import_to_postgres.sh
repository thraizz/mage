#!/bin/bash
# Import H2 SQL export to PostgreSQL database
# Imports into the existing 'mage' database

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
INPUT_SQL="${PROJECT_ROOT}/data/cards.sql"

# PostgreSQL connection settings
DB_NAME="${DB_NAME:-mage}"
DB_USER="${DB_USER:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

echo "=== Import to PostgreSQL ==="
echo "Input: $INPUT_SQL"
echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"
echo ""

# Check if SQL file exists
if [ ! -f "$INPUT_SQL" ]; then
    echo "ERROR: SQL file not found: $INPUT_SQL"
    echo "Run first: ./scripts/h2_to_sql.sh"
    exit 1
fi

# Check if psql is installed
if ! command -v psql &> /dev/null; then
    echo "ERROR: psql not found"
    echo "Install PostgreSQL client"
    exit 1
fi

# Test connection
echo "Testing database connection..."
if ! PAGER=cat psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" > /dev/null 2>&1; then
    echo "ERROR: Cannot connect to PostgreSQL"
    echo "Check connection settings and ensure database exists"
    exit 1
fi
echo "✓ Connected to database"

# Convert H2 SQL to PostgreSQL-compatible SQL
echo "Converting SQL to PostgreSQL format..."
python3 "$SCRIPT_DIR/convert_h2_to_postgres.py" "$INPUT_SQL" "${INPUT_SQL}.postgres.tmp"

# Import to PostgreSQL
echo "Importing to PostgreSQL..."
PAGER=cat psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "${INPUT_SQL}.postgres.tmp" 2>&1 | grep -v "^$" || true

# Cleanup temp file
rm -f "${INPUT_SQL}.postgres.tmp"

# Verify import
CARD_COUNT=$(PAGER=cat psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM card;" 2>/dev/null | xargs || echo "0")
SET_COUNT=$(PAGER=cat psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM expansion;" 2>/dev/null | xargs || echo "0")

echo "✓ Import complete"
echo ""
echo "=== Database Statistics ==="
echo "Cards: $CARD_COUNT"
echo "Sets: $SET_COUNT"
echo ""

# Test queries
echo "Sample cards:"
PAGER=cat psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT name, setcode, manacosts FROM card WHERE name LIKE 'Lightning%' LIMIT 5;" 2>/dev/null || true
echo ""

# Create indexes for performance
echo "Creating indexes..."
PAGER=cat psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<'EOF'
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
echo ""
echo "Database ready!"
echo ""
echo "Test with:"
echo "  PAGER=cat psql -d $DB_NAME -c 'SELECT COUNT(*) FROM card;'"
echo "  PAGER=cat psql -d $DB_NAME -c \"SELECT name, manacosts FROM card WHERE name = 'Lightning Bolt';\""
