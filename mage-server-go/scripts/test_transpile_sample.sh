#!/bin/bash
# Test transpiler on a random sample of cards

JAVA_CARDS_DIR="/Users/aron/dev/opensource/mage/Mage.Sets/src/mage/cards"
OUTPUT_DIR="internal/game/cards/generated"
SAMPLE_SIZE=100

echo "Testing transpiler on $SAMPLE_SIZE random cards..."
echo ""

# Clean output directory
rm -rf "$OUTPUT_DIR"/*.go
mkdir -p "$OUTPUT_DIR"

# Initialize counters
total=0
success=0
failed=0
has_todo=0

# Get random sample of cards
find "$JAVA_CARDS_DIR" -name "*.java" -type f | sort -R | head -n $SAMPLE_SIZE | while read java_file; do
    total=$((total + 1))
    card_name=$(basename "$java_file" .java)

    printf "[$total/$SAMPLE_SIZE] Testing $card_name... "

    # Try to transpile
    if python3 scripts/transpile_cards.py --card="$card_name" --output="$OUTPUT_DIR" > /dev/null 2>&1; then
        # Check if generated file has TODO comments
        output_file="$OUTPUT_DIR/$(echo $card_name | tr '[:upper:]' '[:lower:]').go"
        if [ -f "$output_file" ] && grep -q "TODO" "$output_file"; then
            echo "⚠ TODO"
            has_todo=$((has_todo + 1))
        else
            echo "✓ OK"
        fi
        success=$((success + 1))
    else
        echo "✗ FAIL"
        failed=$((failed + 1))
    fi
done

# Note: The counters don't persist outside the while loop due to subshell
# So we'll count the actual files
success=$(find "$OUTPUT_DIR" -name "*.go" -not -name "*_test.go" | wc -l | tr -d ' ')
has_todo=$(find "$OUTPUT_DIR" -name "*.go" -not -name "*_test.go" -exec grep -l "TODO" {} \; | wc -l | tr -d ' ')
failed=$((SAMPLE_SIZE - success))
complete=$((success - has_todo))

# Calculate percentages
success_pct=$(echo "scale=1; $success * 100 / $SAMPLE_SIZE" | bc)
failed_pct=$(echo "scale=1; $failed * 100 / $SAMPLE_SIZE" | bc)
todo_pct=$(echo "scale=1; $has_todo * 100 / $SAMPLE_SIZE" | bc)
complete_pct=$(echo "scale=1; $complete * 100 / $SAMPLE_SIZE" | bc)

# Print summary
echo ""
echo "======================================"
echo "Sample Test Results"
echo "======================================"
echo "Sample size:          $SAMPLE_SIZE"
echo "✓ Successful:         $success ($success_pct%)"
echo "✗ Failed:             $failed ($failed_pct%)"
echo "⚠ Has TODO:           $has_todo ($todo_pct%)"
echo "✓ Fully implemented:  $complete ($complete_pct%)"
echo ""

# Test compilation
echo "Testing Go compilation..."
if go build ./internal/game/cards/generated/... 2>&1 | head -20; then
    echo "✓ Generated cards compile!"
else
    echo "✗ Compilation errors found"
fi

echo ""
echo "Extrapolated to full set (30,600 cards):"
echo "  - Fully implemented: ~$((complete * 306)) cards"
echo "  - Need TODO work:    ~$((has_todo * 306)) cards"
echo "  - Failed/Complex:    ~$((failed * 306)) cards"
