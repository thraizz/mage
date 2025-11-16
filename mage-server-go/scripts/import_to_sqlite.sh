#!/bin/bash
# Import H2 SQL export to SQLite database
# Creates a portable cards.db file

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
INPUT_SQL="${PROJECT_ROOT}/data/cards.sql"
OUTPUT_DB="${PROJECT_ROOT}/data/cards.db"

echo "=== Import to SQLite ==="
echo "Input: $INPUT_SQL"
echo "Output: $OUTPUT_DB"
echo ""

# Check if SQL file exists
if [ ! -f "$INPUT_SQL" ]; then
    echo "ERROR: SQL file not found: $INPUT_SQL"
    echo "Run first: ./scripts/h2_to_sql.sh"
    exit 1
fi

# Check if sqlite3 is installed
if ! command -v sqlite3 &> /dev/null; then
    echo "ERROR: sqlite3 not found"
    echo "Install with:"
    echo "  macOS: brew install sqlite3"
    echo "  Ubuntu: sudo apt-get install sqlite3"
    exit 1
fi

# Remove existing database
if [ -f "$OUTPUT_DB" ]; then
    echo "Removing existing database..."
    rm "$OUTPUT_DB"
fi

# Convert H2 SQL to SQLite-compatible SQL
echo "Converting SQL to SQLite format..."
python3 "$SCRIPT_DIR/convert_h2_to_sqlite.py" "$INPUT_SQL" "${OUTPUT_DB}.tmp.sql"

# Import to SQLite
echo "Importing to SQLite..."
sqlite3 "$OUTPUT_DB" < "${OUTPUT_DB}.tmp.sql"

# Cleanup temp file
rm -f "${OUTPUT_DB}.tmp.sql"

# Verify import
CARD_COUNT=$(sqlite3 "$OUTPUT_DB" "SELECT COUNT(*) FROM card;" 2>/dev/null || echo "0")
SET_COUNT=$(sqlite3 "$OUTPUT_DB" "SELECT COUNT(*) FROM expansion;" 2>/dev/null || echo "0")

echo "✓ Import complete"
echo ""
echo "=== Database Statistics ==="
echo "Cards: $CARD_COUNT"
echo "Sets: $SET_COUNT"
echo "Size: $(du -h "$OUTPUT_DB" | cut -f1)"
echo ""

# Test queries
echo "Sample cards:"
sqlite3 "$OUTPUT_DB" "SELECT name, setCode, manaCosts FROM card WHERE name LIKE 'Lightning%' LIMIT 5;" 2>/dev/null || true
echo ""

# Create indexes for performance
echo "Creating indexes..."
sqlite3 "$OUTPUT_DB" <<'EOF'
CREATE INDEX IF NOT EXISTS idx_card_name ON card(name);
CREATE INDEX IF NOT EXISTS idx_card_setcode ON card(setCode);
CREATE INDEX IF NOT EXISTS idx_card_classname ON card(className);
CREATE INDEX IF NOT EXISTS idx_card_name_set ON card(name, setCode);
CREATE INDEX IF NOT EXISTS idx_card_types ON card(types);
CREATE INDEX IF NOT EXISTS idx_card_manavalue ON card(manaValue);
ANALYZE;
EOF

echo "✓ Indexes created"
echo ""
echo "Database ready: $OUTPUT_DB"
echo ""
echo "Test with:"
echo "  sqlite3 $OUTPUT_DB 'SELECT COUNT(*) FROM card;'"
echo "  sqlite3 $OUTPUT_DB \"SELECT name, manaCosts FROM card WHERE name = 'Lightning Bolt';\""
