# Transpiler Statistics Tracking - Summary

## What Changed

The transpiler now tracks detailed statistics about successes and failures, saving them to `transpile_stats.json`.

## How to Use

### Single Card
```bash
python3 scripts/transpile_cards.py --card=LightningBolt --stats=transpile_stats.json
```

### Batch Mode
```bash
# Test with 50 cards
./scripts/batch_transpile.sh 50

# Full batch
./scripts/batch_transpile.sh

# Custom limit (using Python directly)
python3 scripts/transpile_cards.py --batch --limit=1000 --stats=transpile_stats.json
```

## Output Format

### JSON Structure
```json
{
  "timestamp": "2025-11-23T11:00:13.600246",
  "summary": {
    "total": 1000,
    "successful": 994,
    "failed": 6,
    "success_rate": "99.40%"
  },
  "error_categories": {
    "ParseError": {
      "count": 6,
      "cards": ["RagnarokDivineDeliverance", "RainOfRiches", ...]
    }
  },
  "failed_cards": [
    {
      "card_name": "RagnarokDivineDeliverance",
      "java_file": "/path/to/RagnarokDivineDeliverance.java",
      "error_type": "ParseError",
      "error_message": "Could not find class name in /path/to/..."
    }
  ]
}
```

### Console Output
```
=== Transpilation Summary ===
Total: 1000
Successful: 994
Failed: 6
Success rate: 99.40%

=== Error Breakdown ===
ParseError: 6 cards
```

## Error Types

### ParseError
Cards that failed to parse the Java source file.

**Common Causes:**
- Unsupported card type (not extending CardImpl, SplitCard, ModalDoubleFacedCard, or AdventureCard)
- Malformed Java file
- Missing class declaration

**Example:**
```json
{
  "error_type": "ParseError",
  "error_message": "Could not find class name in /path/to/Card.java"
}
```

### ValueError
Invalid card data or missing required fields.

**Common Causes:**
- Missing mana cost
- Invalid card types
- Malformed ability data

### Other Exceptions
Any other Python exception during transpilation.

## Analyzing Results

### View Summary
```bash
cat transpile_stats.json | python3 -m json.tool | grep -A 10 '"summary"'
```

### List Failed Cards
```bash
cat transpile_stats.json | python3 -m json.tool | grep -A 1 '"card_name"' | grep -v "^--$"
```

### Count Errors by Type
```bash
cat transpile_stats.json | python3 -m json.tool | grep -B 1 '"count"'
```

### Export Failed Cards to CSV
```bash
python3 << 'PYEOF'
import json
with open('transpile_stats.json') as f:
    stats = json.load(f)
print("Card Name,Error Type,Error Message")
for card in stats['failed_cards']:
    print(f"{card['card_name']},{card['error_type']},{card['error_message']}")
PYEOF
```

## Recent Test Results (1000 cards)

- **Total**: 1000
- **Successful**: 994 (99.4%)
- **Failed**: 6 (0.6%)

**Failures:**
- 6 ParseError (cards extending unsupported base classes)

## Next Steps

When you see failures:

1. **Check error_categories** to see common patterns
2. **Look at specific error messages** in failed_cards
3. **Update transpile_cards.py** to handle new patterns
4. **Re-run affected cards** to verify fixes
5. **Run full batch** once fixes are verified

## Implementation Details

### Code Location
- Main transpiler: `scripts/transpile_cards.py`
- Stats function: `save_stats()` (line ~1758)
- Result class: `TranspileResult` (line ~1699)

### Key Changes
1. Added `TranspileResult` dataclass for structured results
2. Modified `transpile_card()` to return result objects
3. Added `save_stats()` to aggregate and save statistics
4. Updated batch mode to collect results and generate stats
5. Enhanced error handling with categorization
