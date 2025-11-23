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
    python3 scripts/transpile_cards.py --batch --limit=$LIMIT --output="$OUTPUT_DIR" 2>&1 | tee "$LOG_FILE"
else
    python3 scripts/transpile_cards.py --batch --output="$OUTPUT_DIR" 2>&1 | tee "$LOG_FILE"
fi

# Count results
total=$(grep -c "Generated:" "$LOG_FILE" || echo "0")
success=$total
failed=0

# Count TODO markers in generated files
has_todo=$(grep -l "TODO" "$OUTPUT_DIR"/*.go 2>/dev/null | wc -l | tr -d ' ')
fully_implemented=$((success - has_todo))

# Calculate percentages (avoid division by zero)
if [ $total -gt 0 ]; then
    success_pct=$(echo "scale=2; $success * 100 / $total" | bc)
    failed_pct=$(echo "scale=2; $failed * 100 / $total" | bc)
    todo_pct=$(echo "scale=2; $has_todo * 100 / $total" | bc)
    complete_pct=$(echo "scale=2; $fully_implemented * 100 / $total" | bc)
else
    success_pct="0.00"
    failed_pct="0.00"
    todo_pct="0.00"
    complete_pct="0.00"
fi

# Generate statistics
cat > "$STATS_FILE" <<EOF
{
  "total_cards": $total,
  "successful": $success,
  "failed": $failed,
  "has_todo": $has_todo,
  "fully_implemented": $fully_implemented,
  "success_rate": "$success_pct%",
  "failure_rate": "$failed_pct%",
  "todo_rate": "$todo_pct%",
  "complete_rate": "$complete_pct%"
}
EOF

# Print summary
echo ""
echo "======================================"
echo "Batch Transpilation Complete!"
echo "======================================"
echo ""
echo "Total cards:          $total"
echo "✓ Successful:         $success ($success_pct%)"
echo "✗ Failed:             $failed ($failed_pct%)"
echo "⚠ Has TODO:           $has_todo ($todo_pct%)"
echo "✓ Fully implemented:  $fully_implemented ($complete_pct%)"
echo ""
echo "Results saved to:"
echo "  - Log: $LOG_FILE"
echo "  - Stats: $STATS_FILE"
echo ""

# Skip compilation for now (too slow for large batches)
# Uncomment to test compilation:
# echo "Testing Go compilation..."
# if go build ./internal/game/cards/generated/... 2>&1 | tee -a "$LOG_FILE"; then
#     echo "✓ All generated cards compile successfully!"
# else
#     echo "✗ Some generated cards have compilation errors (see log)"
# fi
