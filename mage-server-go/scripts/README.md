# MAGE Scripts Guide

This directory contains various scripts for working with the MAGE project.

## Card Transpilation Scripts

### Main Scripts (Use These)

**`transpile_cards.py`** - The main transpiler script
Converts Java card implementations to Go code.

```bash
# Transpile a single card
python3 scripts/transpile_cards.py --card=LightningBolt

# Transpile all cards
python3 scripts/transpile_cards.py --batch

# Transpile first 1000 cards
python3 scripts/transpile_cards.py --batch --limit=1000

# Specify output location and stats file
python3 scripts/transpile_cards.py --batch \
  --output=internal/game/cards/generated \
  --stats=transpile_stats.json
```

**`batch_transpile.sh`** - Shell wrapper for batch transpilation
Uses the Python transpiler in batch mode with statistics tracking.

```bash
# Transpile all cards
./scripts/batch_transpile.sh

# Transpile first 50 cards (for testing)
./scripts/batch_transpile.sh 50

# Transpile first 1000 cards
./scripts/batch_transpile.sh 1000
```

**Output:**
- Generated Go files: `internal/game/cards/generated/*.go`
- Statistics: `transpile_stats.json` (includes error breakdown and TODO analysis)
- Log: `transpile_results.log`

**`analyze_transpile_stats.py`** - Analyze transpilation statistics

View and query statistics from `transpile_stats.json`:

```bash
# Show full analysis
python3 scripts/analyze_transpile_stats.py

# Show summary only
python3 scripts/analyze_transpile_stats.py --summary

# Show TODO analysis
python3 scripts/analyze_transpile_stats.py --todos

# Show error details
python3 scripts/analyze_transpile_stats.py --errors

# Export failed cards to CSV
python3 scripts/analyze_transpile_stats.py --export-csv failed_cards.csv
```

## Database Scripts

**`build_sqlite_card_db.sh`** - Build SQLite database from card data

**`import_to_postgres.sh`** - Import card data to PostgreSQL

**`import_to_sqlite.sh`** - Import card data to SQLite

**`h2_to_sql.sh`** - Convert H2 database to SQL

**`export_java_cards.sh`** - Export Java card metadata

## Protobuf Scripts

**`generate_proto.sh`** - Generate Go code from .proto files

```bash
make proto  # Preferred - use the Makefile target
# OR
./scripts/generate_proto.sh
```

## Understanding Transpile Statistics

The `transpile_stats.json` file provides detailed information about the transpilation process:

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
      "cards": ["Card1", "Card2", ...]
    }
  },
  "failed_cards": [
    {
      "card_name": "CardName",
      "java_file": "/path/to/Card.java",
      "error_type": "ParseError",
      "error_message": "Could not find class name..."
    }
  ]
}
```

### Common Error Types

- **ParseError**: Failed to parse Java file (usually missing class name or unsupported card type)
- **ValueError**: Invalid card data or missing required fields
- **IOError**: File system issues

## Workflow

### Full Transpilation Workflow

```bash
# 1. Clean previous transpilation
rm -rf internal/game/cards/generated
rm transpile_stats.json

# 2. Run batch transpilation
./scripts/batch_transpile.sh

# 3. Check results
cat transpile_stats.json | python3 -m json.tool

# 4. Test compilation
go build ./internal/game/cards/generated/...

# 5. Fix errors if needed
# Check failed_cards in transpile_stats.json
# Update transpile_cards.py to handle new patterns
```

### Development Workflow

```bash
# 1. Make changes to transpile_cards.py

# 2. Test with sample cards
./scripts/batch_transpile.sh 50  # Quick test with 50 cards

# 3. Check for errors
cat transpile_stats.json

# 4. If good, test with more cards
./scripts/batch_transpile.sh 1000  # Test with 1000

# 5. Then full transpilation
./scripts/batch_transpile.sh
```

## Performance Notes

- **Python batch mode**: ~200x faster than shell-based transpilation
- **1000 cards**: ~15 seconds
- **50,000+ cards**: ~10 minutes
- Stats tracking adds negligible overhead (~1%)
