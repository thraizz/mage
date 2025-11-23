#!/bin/bash
# Test version - transpile limited cards

set -e

JAVA_CARDS_DIR="/Users/aron/dev/opensource/mage/Mage.Sets/src/mage/cards"
OUTPUT_DIR="internal/game/cards/generated"
LOG_FILE="transpile_results.log"

echo "Starting batch transpilation (TEST - LIMITED CARDS)..."

# Initialize counters
total=0
success=0
failed=0

# Create log file
echo "Transpilation Log - $(date)" > "$LOG_FILE"

# Find first N Java card files and transpile (limit for testing)
LIMIT=${1:-50}  # Default to 50 cards, or use first argument

for java_file in $(find "$JAVA_CARDS_DIR" -name "*.java" -type f | head -n $LIMIT | sort); do
    total=$((total + 1))
    card_name=$(basename "$java_file" .java)

    # Try to transpile
    if python3 scripts/transpile_cards.py --card="$card_name" --output="$OUTPUT_DIR" >> "$LOG_FILE" 2>&1; then
        success=$((success + 1))
    else
        failed=$((failed + 1))
    fi

    # Progress indicator every 10 cards
    if [ $((total % 10)) -eq 0 ]; then
        echo "Processed $total cards... (Success: $success, Failed: $failed)"
    fi
done

# Print summary
echo "Total: $total | Success: $success | Failed: $failed"
