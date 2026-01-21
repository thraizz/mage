# ⚠️ ARCHIVED - Migration Complete (2026-01-20)

**This document describes a historical migration that is now complete.**

For current card data architecture, see [CARD_DATA_ARCHITECTURE.md](CARD_DATA_ARCHITECTURE.md).

---

# Scryfall Data Migration - Making scryfall_cards Primary

**Date**: 2026-01-19
**Status**: ✅ Complete (Superseded by XMage dependency removal on 2026-01-20)

---

## Problem

Backend was querying the old `cards` table (87,962 XMage cards) instead of `scryfall_cards` (101,852 Scryfall cards), causing deck upload failures for new cards like:
- Hexing Squelcher (Lorwyn Eclipsed)
- Path of the Pyromancer (Modern Horizons Commanders)

---

## Solution

Updated `internal/repository/cards.go` to query `scryfall_cards` as the primary data source.

### Changes Made

All query functions now use `scryfall_cards WHERE lang='en'`:

1. **GetByID** - Queries scryfall_cards
2. **GetByName** - Queries scryfall_cards
3. **GetByNameCaseInsensitive** ⭐ - Queries scryfall_cards (used for deck validation)
4. **SearchByName** - Queries scryfall_cards
5. **GetBySetCode** - Queries scryfall_cards

### Field Mapping

Scryfall → Card struct:
- `collector_number` → `CardNumber`
- `type_line` → `CardType`
- `oracle_text` → `RulesText`
- `id` (UUID) → `ID` (converted to int64 via bit conversion)
- `flavor_text` → NULL (not in Scryfall bulk data)
- `card_class_name` → "" (XMage-specific)

### Exception: GetByClassName

This function still queries the old `cards` table (or `cards_xmage_backup`) because Java class names are XMage-specific and not in Scryfall data. It tries `cards_xmage_backup` first, falls back to `cards` if that doesn't exist.

---

## Deployment

### Local Development

```bash
# Build the updated server
cd mage-server-go
go build -o mage-server ./cmd/server

# Rebuild Docker image
docker compose -f ../docker-compose.prod.yml build mage-server
```

### Production

```bash
# SSH into production server
ssh hkdebiandocker@192.168.178.24

# Navigate to project
cd ~/gomage

# Pull latest code
git pull origin master

# Rebuild and restart
docker compose build mage-server
docker compose up -d mage-server

# Verify
docker compose logs -f mage-server
```

---

## Verification

### Test Deck Upload

Try uploading a deck with "Hexing Squelcher":

```
Commander:
1 Atraxa, Praetors' Voice

1 Hexing Squelcher
1 Path of the Pyromancer
98 other cards...
```

**Before**: ❌ `invalid card names: Hexing Squelcher, Path of the Pyromancer`
**After**: ✅ Deck uploaded successfully

### Database Queries

```bash
# Check scryfall_cards has the cards
docker compose exec postgres psql -U mage -d mage -c \
  "SELECT name, set_code FROM scryfall_cards
   WHERE name IN ('Hexing Squelcher', 'Path of the Pyromancer')
   AND lang='en';"
```

Should return:
```
       name       | set_code
------------------+----------
 Hexing Squelcher | ecl
 Hexing Squelcher | ecl
 Path of the Pyromancer | moc
 Path of the Pyromancer | moc
```

---

## Benefits

✅ **All 101k+ Scryfall cards** now available for deck uploads
✅ **No more "card not found" errors** for new sets
✅ **Automatic updates** - just re-import Scryfall data
✅ **Backward compatible** - old code still works
✅ **No view needed** - queries scryfall_cards directly

---

## Database State

### Before Migration

- `cards` table: 87,962 rows (old XMage data)
- `scryfall_cards` table: 101,852 rows (Scryfall data)
- Backend: Queried `cards` → missing 14k+ cards

### After Migration

- `cards` table: Still exists (can be renamed to `cards_xmage_backup`)
- `scryfall_cards` table: 101,852 rows (PRIMARY source)
- Backend: Queries `scryfall_cards` → all cards available

### Optional: Rename Old Cards Table

```bash
# On production server
docker compose exec postgres psql -U mage -d mage -c \
  "ALTER TABLE cards RENAME TO cards_xmage_backup;"
```

This preserves the XMage data for `GetByClassName()` while making it clear it's legacy data.

---

## Code Changes

### File: `internal/repository/cards.go`

**Before**:
```go
func (r *CardRepository) GetByNameCaseInsensitive(ctx context.Context, name string) ([]*Card, error) {
    query := `
        SELECT ...
        FROM cards  -- Old table (87k cards)
        WHERE LOWER(TRIM(name)) = LOWER(TRIM($1))
    `
    // ...
}
```

**After**:
```go
func (r *CardRepository) GetByNameCaseInsensitive(ctx context.Context, name string) ([]*Card, error) {
    query := `
        SELECT ('x' || substring(id::text, 1, 16))::bit(64)::bigint as id,
               COALESCE(collector_number, ''),
               set_code,
               name,
               COALESCE(type_line, ''),
               -- ... map Scryfall fields to Card struct
        FROM scryfall_cards  -- New table (101k+ cards)
        WHERE lang = 'en'
          AND (LOWER(TRIM(name)) = LOWER(TRIM($1))
               OR LOWER(name) LIKE LOWER($1) || ' //%')
        ORDER BY released_at DESC, set_code, collector_number
    `
    // ...
}
```

---

## Testing

### Unit Test

```go
func TestCardRepository_GetByNameCaseInsensitive_Scryfall(t *testing.T) {
    // Setup test database with scryfall_cards
    repo := NewCardRepository(testDB, logger)

    // Test recent card (ECL set)
    cards, err := repo.GetByNameCaseInsensitive(context.Background(), "Hexing Squelcher")
    assert.NoError(t, err)
    assert.NotEmpty(t, cards)
    assert.Equal(t, "Hexing Squelcher", cards[0].Name)
    assert.Equal(t, "ecl", cards[0].SetCode)

    // Test MOC set card
    cards, err = repo.GetByNameCaseInsensitive(context.Background(), "Path of the Pyromancer")
    assert.NoError(t, err)
    assert.NotEmpty(t, cards)
}
```

### Integration Test

```bash
# Test deck upload via API
curl -X POST http://localhost:17171/deck/save \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-session",
    "deck_name": "Test Deck",
    "format": "Commander",
    "deck": {
      "commanders": [{"name": "Atraxa, Praetors'\'' Voice", "quantity": 1}],
      "main_deck": [
        {"name": "Hexing Squelcher", "quantity": 1},
        {"name": "Path of the Pyromancer", "quantity": 1}
      ]
    }
  }'

# Expected: {"success": true, "deck_id": 123}
```

---

## Rollback Plan

If issues occur:

### Option 1: Revert Code

```bash
git revert <commit-hash>
docker compose build mage-server
docker compose up -d mage-server
```

### Option 2: Create Compatibility View

```sql
-- Make 'cards' a view of scryfall_cards
ALTER TABLE cards RENAME TO cards_old;

CREATE VIEW cards AS
SELECT
    ('x' || substring(id::text, 1, 16))::bit(64)::bigint as id,
    collector_number as card_number,
    set_code,
    name,
    type_line as card_type,
    mana_cost,
    power,
    toughness,
    oracle_text as rules_text,
    NULL::text as flavor_text,
    NULL::text as original_text,
    NULL::text as original_type,
    CASE WHEN collector_number ~ '^[0-9]+$' THEN collector_number::bigint ELSE NULL END as cn,
    name as card_name,
    rarity,
    '' as card_class_name,
    created_at
FROM scryfall_cards
WHERE lang = 'en';
```

---

## Future Improvements

1. **Remove legacy cards table** - Once confirmed GetByClassName is not needed
2. **Add card index** - Create indexes on scryfall_cards(name) for faster lookups
3. **Cache layer** - Keep in-memory cache of frequently used cards
4. **JSON index** - Consider loading Scryfall JSON into memory (see DB_JSON_Findings.md)

---

## Related Documentation

- **DB_JSON_Findings.md** - Analysis of database vs JSON approach
- **PRODUCTION_CARD_UPDATES.md** - How to update Scryfall data
- **SCRYFALL_MIGRATION_GUIDE.md** - Original migration from XMage to Scryfall

---

## Summary

✅ Backend now queries `scryfall_cards` as primary source
✅ All 101k+ cards available for deck uploads
✅ No more missing card errors
✅ Backward compatible with legacy code
✅ Build successful, ready for deployment
