#!/usr/bin/env bash
#
# Download Scryfall bulk data (all cards)
#
# Usage:
#   ./scripts/download_scryfall_bulk.sh
#
# Output: data/scryfall-all-cards-{date}.json
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
DATA_DIR="$GO_ROOT/data"
DATE=$(date +%Y%m%d%H%M%S)

mkdir -p "$DATA_DIR"

echo "=== Scryfall Bulk Data Download ==="

# Get bulk data info from Scryfall API
echo "Fetching bulk data info from Scryfall API..."
BULK_INFO=$(curl -s https://api.scryfall.com/bulk-data/all-cards)

# Extract download URL and updated timestamp
DOWNLOAD_URL=$(echo "$BULK_INFO" | grep -o '"download_uri":"[^"]*"' | cut -d'"' -f4)
UPDATED_AT=$(echo "$BULK_INFO" | grep -o '"updated_at":"[^"]*"' | head -1 | cut -d'"' -f4)
FILE_SIZE=$(echo "$BULK_INFO" | grep -o '"compressed_size":[0-9]*' | cut -d':' -f2)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "ERROR: Failed to get download URL from Scryfall API"
    exit 1
fi

echo "Download URL: $DOWNLOAD_URL"
echo "Last updated: $UPDATED_AT"
echo "Compressed size: $(numfmt --to=iec-i --suffix=B $FILE_SIZE 2>/dev/null || echo "${FILE_SIZE} bytes")"

OUTPUT_FILE="$DATA_DIR/scryfall-all-cards-${DATE}.json"

# Download with progress
echo "Downloading to: $OUTPUT_FILE"
curl -L --progress-bar "$DOWNLOAD_URL" -o "$OUTPUT_FILE"

# Verify file was downloaded
if [ ! -f "$OUTPUT_FILE" ]; then
    echo "ERROR: Download failed"
    exit 1
fi

ACTUAL_SIZE=$(stat -f%z "$OUTPUT_FILE" 2>/dev/null || stat -c%s "$OUTPUT_FILE" 2>/dev/null || echo "unknown")
echo ""
echo "✓ Download complete"
echo "  File: $OUTPUT_FILE"
echo "  Size: $(numfmt --to=iec-i --suffix=B $ACTUAL_SIZE 2>/dev/null || echo "$ACTUAL_SIZE bytes")"
echo ""

# Create symlink to latest
LATEST_LINK="$DATA_DIR/scryfall-all-cards-latest.json"
rm -f "$LATEST_LINK"
ln -s "$(basename "$OUTPUT_FILE")" "$LATEST_LINK"
echo "✓ Created symlink: $LATEST_LINK -> $(basename "$OUTPUT_FILE")"
echo ""

# Show some stats
echo "Quick stats:"
TOTAL_CARDS=$(grep -o '"id"' "$OUTPUT_FILE" | wc -l | tr -d ' ')
echo "  Total cards: $TOTAL_CARDS"

echo ""
echo "Next step:"
echo "  go run ./cmd/scryfall-import/main.go --input=\"$OUTPUT_FILE\""
