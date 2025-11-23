#!/bin/bash
# Parallel transpilation test - fixed version

set -e

JAVA_CARDS_DIR="/Users/aron/dev/opensource/mage/Mage.Sets/src/mage/cards"
OUTPUT_DIR="internal/game/cards/generated"
LIMIT=${1:-50}

echo "Starting parallel transpilation ($LIMIT cards)..."

# Use simpler xargs approach - pass just the basename
find "$JAVA_CARDS_DIR" -name "*.java" -type f | head -n $LIMIT | while read -r java_file; do
    basename "$java_file" .java
done | xargs -P 8 -I {} python3 scripts/transpile_cards.py --card="{}" --output="$OUTPUT_DIR" >/dev/null 2>&1

# Count results
success=$(find "$OUTPUT_DIR" -name "*.go" | wc -l)
echo "Generated $success Go files"
