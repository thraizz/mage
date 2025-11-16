#!/bin/bash
# Build SQLite card database from Java MAGE's H2 database
# This creates a portable cards.db file that can be shipped with the server

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
JAVA_DB_PATH="${PROJECT_ROOT}/../Mage.Server/db/cards"
OUTPUT_DB="${PROJECT_ROOT}/data/cards.db"
TEMP_SQL="${PROJECT_ROOT}/data/cards_export.sql"

echo "=== Building SQLite Card Database ==="
echo "Java H2 database: ${JAVA_DB_PATH}.h2.mv.db"
echo "Output SQLite: $OUTPUT_DB"
echo ""

# Check if H2 database exists
if [ ! -f "${JAVA_DB_PATH}.h2.mv.db" ]; then
    echo "ERROR: Java H2 database not found at ${JAVA_DB_PATH}.h2.mv.db"
    echo "Please run the Java MAGE server at least once to generate the card database."
    echo ""
    echo "To generate the Java database:"
    echo "  cd ../Mage.Server"
    echo "  mvn exec:java"
    exit 1
fi

# Create output directory
mkdir -p "$(dirname "$OUTPUT_DB")"

# Download H2 JAR if not present
H2_JAR="${SCRIPT_DIR}/h2.jar"
if [ ! -f "$H2_JAR" ]; then
    echo "Downloading H2 database JAR..."
    curl -L -o "$H2_JAR" "https://repo1.maven.org/maven2/com/h2database/h2/2.2.224/h2-2.2.224.jar"
    echo "✓ H2 JAR downloaded"
fi

# Export H2 database to SQL
echo "Exporting H2 database to SQL..."
java -cp "$H2_JAR" org.h2.tools.Script \
    -url "jdbc:h2:${JAVA_DB_PATH}" \
    -user "sa" \
    -password "" \
    -script "$TEMP_SQL"

echo "✓ Exported to SQL"

# Check if sqlite3 is installed
if ! command -v sqlite3 &> /dev/null; then
    echo "ERROR: sqlite3 command not found"
    echo "Please install SQLite3:"
    echo "  macOS: brew install sqlite3"
    echo "  Ubuntu: sudo apt-get install sqlite3"
    exit 1
fi

# Remove existing database
if [ -f "$OUTPUT_DB" ]; then
    echo "Removing existing database..."
    rm "$OUTPUT_DB"
fi

# Create SQLite database
echo "Creating SQLite database..."
sqlite3 "$OUTPUT_DB" <<'EOF'
-- Create cards table
CREATE TABLE cards (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    set_code TEXT NOT NULL,
    card_number TEXT,
    card_number_as_int INTEGER,
    class_name TEXT NOT NULL,
    mana_cost TEXT,
    mana_value INTEGER DEFAULT 0,
    power TEXT,
    toughness TEXT,
    starting_loyalty TEXT,
    starting_defense TEXT,
    types TEXT,
    subtypes TEXT,
    supertypes TEXT,
    rules_text TEXT,
    rarity TEXT,
    frame_style TEXT,
    frame_color TEXT,
    color_identity TEXT,
    is_black BOOLEAN DEFAULT 0,
    is_blue BOOLEAN DEFAULT 0,
    is_green BOOLEAN DEFAULT 0,
    is_red BOOLEAN DEFAULT 0,
    is_white BOOLEAN DEFAULT 0,
    is_colorless BOOLEAN DEFAULT 0,
    is_split_card BOOLEAN DEFAULT 0,
    is_double_faced BOOLEAN DEFAULT 0,
    is_flip_card BOOLEAN DEFAULT 0,
    is_modal_double_faced BOOLEAN DEFAULT 0,
    various_art BOOLEAN DEFAULT 0,
    created_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- Create indexes
CREATE INDEX idx_cards_name ON cards(name);
CREATE INDEX idx_cards_set_code ON cards(set_code);
CREATE INDEX idx_cards_class_name ON cards(class_name);
CREATE INDEX idx_cards_name_set ON cards(name, set_code);
CREATE INDEX idx_cards_mana_value ON cards(mana_value);
CREATE INDEX idx_cards_rarity ON cards(rarity);

-- Create full-text search
CREATE VIRTUAL TABLE cards_fts USING fts5(
    name, 
    rules_text, 
    types, 
    subtypes,
    content=cards,
    content_rowid=id
);

-- Create sets table
CREATE TABLE sets (
    id INTEGER PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    release_date TEXT,
    set_type TEXT,
    block_name TEXT,
    has_boosters BOOLEAN DEFAULT 0,
    has_basic_lands BOOLEAN DEFAULT 0
);

CREATE INDEX idx_sets_code ON sets(code);
CREATE INDEX idx_sets_type ON sets(set_type);

-- Metadata table
CREATE TABLE metadata (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

INSERT INTO metadata (key, value) VALUES 
    ('version', '1.0'),
    ('created_at', datetime('now')),
    ('source', 'mage-h2-export');
EOF

echo "✓ SQLite schema created"

# Convert H2 SQL to SQLite format and import
echo "Converting and importing data..."
python3 - <<'PYTHON_EOF'
import re
import sqlite3
import sys

# Read H2 SQL export
with open('data/cards_export.sql', 'r', encoding='utf-8') as f:
    sql = f.read()

# Connect to SQLite
conn = sqlite3.connect('data/cards.db')
cursor = conn.cursor()

# Parse INSERT statements
# H2 format: INSERT INTO CARD(...) VALUES(...);
# We need to convert to SQLite format

card_count = 0
set_count = 0

# Extract card inserts
card_pattern = r"INSERT INTO CARD\s*\((.*?)\)\s*VALUES\s*\((.*?)\);"
for match in re.finditer(card_pattern, sql, re.IGNORECASE | re.DOTALL):
    columns = match.group(1)
    values = match.group(2)
    
    # Parse columns and values
    cols = [c.strip().lower() for c in columns.split(',')]
    vals = []
    
    # Simple value parsing (handles strings, numbers, booleans)
    in_string = False
    current_val = ""
    for char in values:
        if char == "'" and (not current_val or current_val[-1] != '\\'):
            in_string = not in_string
        if char == ',' and not in_string:
            vals.append(current_val.strip())
            current_val = ""
        else:
            current_val += char
    if current_val:
        vals.append(current_val.strip())
    
    # Map H2 columns to SQLite columns
    col_map = {
        'name': 'name',
        'setcode': 'set_code',
        'cardnumber': 'card_number',
        'classname': 'class_name',
        'manacosts': 'mana_cost',
        'manavalue': 'mana_value',
        'power': 'power',
        'toughness': 'toughness',
        'startingloyalty': 'starting_loyalty',
        'startingdefense': 'starting_defense',
        'types': 'types',
        'subtypes': 'subtypes',
        'supertypes': 'supertypes',
        'rules': 'rules_text',
        'rarity': 'rarity',
        'framestyle': 'frame_style',
        'framecolor': 'frame_color',
        'black': 'is_black',
        'blue': 'is_blue',
        'green': 'is_green',
        'red': 'is_red',
        'white': 'is_white',
    }
    
    # Build insert statement
    sqlite_cols = []
    sqlite_vals = []
    for i, col in enumerate(cols):
        if col in col_map and i < len(vals):
            sqlite_cols.append(col_map[col])
            val = vals[i]
            # Convert boolean values
            if val.upper() in ('TRUE', 'FALSE'):
                val = '1' if val.upper() == 'TRUE' else '0'
            sqlite_vals.append(val)
    
    if sqlite_cols:
        try:
            query = f"INSERT INTO cards ({','.join(sqlite_cols)}) VALUES ({','.join(['?' for _ in sqlite_vals])})"
            # Clean values (remove quotes)
            clean_vals = [v.strip("'") if v.startswith("'") else v for v in sqlite_vals]
            cursor.execute(query, clean_vals)
            card_count += 1
        except Exception as e:
            print(f"Warning: Failed to insert card: {e}", file=sys.stderr)

# Similar for sets (if present in export)
# ... (simplified for now)

conn.commit()
conn.close()

print(f"Imported {card_count} cards")
print(f"Imported {set_count} sets")
PYTHON_EOF

# Verify import
CARD_COUNT=$(sqlite3 "$OUTPUT_DB" "SELECT COUNT(*) FROM cards;")
echo "✓ Imported $CARD_COUNT cards"

# Update FTS index
echo "Building full-text search index..."
sqlite3 "$OUTPUT_DB" "INSERT INTO cards_fts(cards_fts) VALUES('rebuild');"
echo "✓ FTS index built"

# Optimize database
echo "Optimizing database..."
sqlite3 "$OUTPUT_DB" "VACUUM; ANALYZE;"
echo "✓ Database optimized"

# Get database size
DB_SIZE=$(du -h "$OUTPUT_DB" | cut -f1)
echo ""
echo "=== Build Complete ==="
echo "✓ Database: $OUTPUT_DB"
echo "✓ Size: $DB_SIZE"
echo "✓ Cards: $CARD_COUNT"
echo ""

# Create checksum
SHA256=$(shasum -a 256 "$OUTPUT_DB" | cut -d' ' -f1)
echo "$SHA256" > "${OUTPUT_DB}.sha256"
echo "✓ Checksum: ${OUTPUT_DB}.sha256"
echo ""

# Test query
echo "Sample cards:"
sqlite3 "$OUTPUT_DB" "SELECT name, set_code, mana_cost FROM cards WHERE name LIKE 'Lightning%' LIMIT 5;"
echo ""

echo "Next steps:"
echo "  1. Verify: sqlite3 $OUTPUT_DB 'SELECT COUNT(*) FROM cards;'"
echo "  2. Test FTS: sqlite3 $OUTPUT_DB \"SELECT name FROM cards_fts WHERE cards_fts MATCH 'flying' LIMIT 5;\""
echo "  3. Use in Go: cardDB := NewCardDB(\"$OUTPUT_DB\", logger)"

# Cleanup
rm -f "$TEMP_SQL"
