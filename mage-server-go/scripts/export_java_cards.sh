#!/bin/bash
#
# Export card metadata from Java source files to CSV
#
# This script parses Java card files and extracts metadata without
# needing to run the Java server or access an H2 database.
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
JAVA_CARDS_DIR="$PROJECT_ROOT/../Mage.Sets/src/mage/cards"
OUTPUT_FILE="$PROJECT_ROOT/data/cards_export.csv"

echo "=== MAGE Card Metadata Exporter ==="
echo "Java cards directory: $JAVA_CARDS_DIR"
echo "Output file: $OUTPUT_FILE"
echo ""

# Check if Java cards directory exists
if [ ! -d "$JAVA_CARDS_DIR" ]; then
    echo "Error: Java cards directory not found: $JAVA_CARDS_DIR"
    echo "Please ensure Mage.Sets is checked out in the parent directory"
    exit 1
fi

# Count total cards
TOTAL_CARDS=$(find "$JAVA_CARDS_DIR" -name "*.java" | wc -l | tr -d ' ')
echo "Found $TOTAL_CARDS Java card files"
echo ""

# Create CSV header
echo "Creating CSV with metadata..."
echo "name,set_code,card_number,class_name,power,toughness,starting_loyalty,starting_defense,mana_value,rarity,types,subtypes,supertypes,mana_costs,rules,black,blue,green,red,white,frame_color,frame_style,various_art" > "$OUTPUT_FILE"

# Use Python to parse Java files (more reliable than bash regex)
python3 "$SCRIPT_DIR/parse_java_cards.py" "$JAVA_CARDS_DIR" >> "$OUTPUT_FILE"

# Count exported cards
EXPORTED=$(wc -l < "$OUTPUT_FILE" | tr -d ' ')
EXPORTED=$((EXPORTED - 1))  # Subtract header

echo ""
echo "=== Export Complete ==="
echo "✓ Exported: $EXPORTED cards"
echo "✓ Output: $OUTPUT_FILE"
echo ""
echo "Next steps:"
echo "  1. Review: head -20 $OUTPUT_FILE"
echo "  2. Import: cd mage-server-go && go run scripts/import_cards.go"
