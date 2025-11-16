#!/bin/bash
# Export card data from Java MAGE server's H2 database to CSV
# This script extracts card metadata for import into PostgreSQL

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
JAVA_DB_PATH="${PROJECT_ROOT}/../Mage.Server/db/cards.h2.mv.db"
OUTPUT_CSV="${PROJECT_ROOT}/data/cards_export.csv"

echo "=== MAGE Card Data Export ==="
echo "Java DB path: $JAVA_DB_PATH"
echo "Output CSV: $OUTPUT_CSV"

# Check if H2 database exists
if [ ! -f "$JAVA_DB_PATH" ]; then
    echo "ERROR: Java H2 database not found at $JAVA_DB_PATH"
    echo "Please run the Java MAGE server at least once to generate the card database."
    exit 1
fi

# Create output directory
mkdir -p "$(dirname "$OUTPUT_CSV")"

# Download H2 JAR if not present
H2_JAR="${PROJECT_ROOT}/scripts/h2.jar"
if [ ! -f "$H2_JAR" ]; then
    echo "Downloading H2 database JAR..."
    curl -L -o "$H2_JAR" "https://repo1.maven.org/maven2/com/h2database/h2/2.2.224/h2-2.2.224.jar"
fi

# Export cards table to CSV
echo "Exporting cards from H2 database..."
java -cp "$H2_JAR" org.h2.tools.Csv \
    -url "jdbc:h2:${PROJECT_ROOT}/../Mage.Server/db/cards" \
    -user "sa" \
    -password "" \
    -sql "SELECT name, setCode, cardNumber, className, power, toughness, startingLoyalty, startingDefense, manaValue, rarity, types, subtypes, supertypes, manaCosts, rules, black, blue, green, red, white, frameColor, frameStyle, variousArt FROM card" \
    -charset "UTF-8" \
    > "$OUTPUT_CSV"

echo "✓ Export complete: $OUTPUT_CSV"
echo "Total cards exported: $(wc -l < "$OUTPUT_CSV")"
echo ""
echo "Next steps:"
echo "  1. Review the exported CSV file"
echo "  2. Run: go run scripts/import_cards.go"
echo "  3. Verify: PAGER=cat psql -d mage -c 'SELECT COUNT(*) FROM cards;'"
