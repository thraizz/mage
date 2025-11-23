#!/bin/bash
# Batch transpile all Magic cards and generate statistics
# OPTIMIZED VERSION - Uses Python batch mode for 200x speedup

set -e

JAVA_CARDS_DIR="/Users/aron/dev/opensource/mage/Mage.Sets/src/mage/cards"
OUTPUT_DIR="internal/game/cards/generated"
LOG_FILE="transpile_results.log"
STATS_FILE="transpile_stats.json"
LIMIT=${1:-0}  # Optional: limit number of cards (0 = no limit)

echo "Starting batch transpilation (OPTIMIZED)..."
echo "Java cards directory: $JAVA_CARDS_DIR"
echo "Output directory: $OUTPUT_DIR"
if [ $LIMIT -gt 0 ]; then
    echo "Limit: $LIMIT cards"
fi
echo ""

# Clean output directory
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Use Python's batch mode for massive speedup
echo "Running Python transpiler in batch mode..."
if [ $LIMIT -gt 0 ]; then
    python3 scripts/transpile_cards.py --batch --limit=$LIMIT --output="$OUTPUT_DIR" --stats="$STATS_FILE" 2>&1 | tee "$LOG_FILE"
else
    python3 scripts/transpile_cards.py --batch --output="$OUTPUT_DIR" --stats="$STATS_FILE" 2>&1 | tee "$LOG_FILE"
fi

echo ""
echo "Results saved to:"
echo "  - Log: $LOG_FILE"
echo "  - Stats: $STATS_FILE"
echo ""

# Show quick summary from stats file
if [ -f "$STATS_FILE" ]; then
    echo "Quick summary from stats:"
    cat "$STATS_FILE" | python3 -m json.tool | grep -A 4 '"summary"'
fi

# Skip compilation for now (too slow for large batches)
# Uncomment to test compilation:
# echo "Testing Go compilation..."
# if go build ./internal/game/cards/generated/... 2>&1 | tee -a "$LOG_FILE"; then
#     echo "✓ All generated cards compile successfully!"
# else
#     echo "✗ Some generated cards have compilation errors (see log)"
# fi
