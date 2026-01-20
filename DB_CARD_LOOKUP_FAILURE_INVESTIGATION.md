# Card Lookup Failure Investigation

**Date**: 2026-01-20
**Issue**: Backend rejects all cards during deck upload with "invalid card names" error
**Status**: Root cause identified

---

## Problem Statement

When uploading a 67-card deck list (including basic lands and common staples like Sol Ring, Lightning Bolt, etc.), the backend returns:

```
invalid card names: Agate Instigator, Ancient Tomb, Anger, Arena of Glory, ..., Mountain, ...
[ALL 67 cards rejected]
```

---

## Investigation Timeline

### 1. Database Verification ✅

**Production database status:**
- PostgreSQL container: Running and healthy
- `scryfall_cards` table: **101,852 English cards present**
- Test queries work: `Mountain`, `Lightning Bolt`, `Hexing Squelcher` all found

```sql
SELECT name FROM scryfall_cards
WHERE lang='en' AND LOWER(TRIM(name)) = 'hexing squelcher';
-- Returns: Hexing Squelcher ✓
```

### 2. Backend Code Analysis ✅

**Current production deployment:**
- Docker image built: 2026-01-20 11:00 AM CET
- Code version: Commit `ec5b1c50e81` (Jan 12, 2026)
- Repository layer: Queries from `cards` table/view (OLD code)

**Uncommitted local changes (not deployed):**
- Modified `cards.go` to query `scryfall_cards` directly
- Removed `card_class_name` field
- Changed column mappings (e.g., `collector_number` ← `card_number`)

**Key finding:** Production is running OLD code that queries the `cards` view, not the updated code.

### 3. Database Schema Discovery 🔍

**Tables and views in production:**
```
Tables:
- card (87,962 rows) ← OLD XMage data table
- scryfall_cards (101,852 rows) ← Scryfall data
- cards_xmage_backup ← Backup table

Views:
- cards ← Compatibility view (BROKEN - see below)
- canonical_cards
```

**Migration status:**
```sql
SELECT * FROM schema_migrations;
-- Latest: 008 (create_active_games_table)
-- MISSING: 009 (create_scryfall_tables) ← Not recorded but table exists!
-- MISSING: 010 (drop_xmage_backup) ← Not applied
```

**Critical finding:** Migrations 009 and 010 haven't been officially applied, but `scryfall_cards` table exists (manually imported).

---

## Root Cause 🔴

### The `cards` View is BROKEN

The production `cards` view was created with this ID conversion:

```sql
CREATE VIEW cards AS
SELECT
    ('x' || substring(id::text, 1, 16))::bit(64)::bigint as id,
    ...
FROM scryfall_cards
WHERE lang = 'en';
```

**Problem:** Scryfall UUIDs contain hyphens:
```
id = '6d3ecc2e-7e49-4822-8330-32e57f05ee5e'
substring(id::text, 1, 16) = '6d3ecc2e-7e49-48'
'x' || '6d3ecc2e-7e49-48' = 'x6d3ecc2e-7e49-48'
Converting to bit(64) → ERROR: "-" is not a valid hexadecimal digit
```

**Test query that fails:**
```sql
SELECT id FROM cards WHERE name = 'Mountain' LIMIT 1;
-- ERROR: "-" is not a valid hexadecimal digit
```

**Why it affects deck uploads:**
1. Backend calls `cardRepo.GetByNameCaseInsensitive(ctx, "Mountain")`
2. Repository queries `SELECT id, ... FROM cards WHERE LOWER(TRIM(name)) = ...`
3. View tries to convert UUID → bigint with hyphens
4. PostgreSQL throws error
5. Backend catches error, marks card as "invalid"
6. ALL cards fail the same way

---

## Evidence

### UUID Format in scryfall_cards
```sql
SELECT id FROM scryfall_cards LIMIT 1;
-- 6d3ecc2e-7e49-4822-8330-32e57f05ee5e
--         ^    ^    (hyphens present)
```

### Error When Querying View
```bash
$ psql -c "SELECT id FROM cards WHERE name = 'Mountain' LIMIT 1;"
ERROR:  "-" is not a valid hexadecimal digit
```

### Direct Table Query Works
```bash
$ psql -c "SELECT id FROM scryfall_cards WHERE name = 'Mountain' LIMIT 1;"
-- Returns UUID successfully ✓
```

---

## Fix Required

### Option 1: Fix the View (Quick Fix)
Replace the view with proper UUID→bigint conversion that removes hyphens:

```sql
DROP VIEW IF EXISTS cards CASCADE;
CREATE VIEW cards AS
SELECT
    -- Remove hyphens from UUID before conversion
    ('x' || REPLACE(substring(id::text, 1, 18), '-', ''))::bit(64)::bigint as id,
    collector_number AS card_number,
    set_code,
    name,
    type_line AS card_type,
    COALESCE(mana_cost, '') as mana_cost,
    COALESCE(power, '') as power,
    COALESCE(toughness, '') as toughness,
    COALESCE(oracle_text, '') as rules_text,
    NULL::text as flavor_text,
    NULL::text as original_text,
    NULL::text as original_type,
    CASE WHEN collector_number ~ '^[0-9]+$' THEN collector_number::bigint ELSE NULL END as cn,
    name as card_name,
    rarity,
    ''::text as card_class_name,
    created_at
FROM scryfall_cards
WHERE lang = 'en';
```

### Option 2: Deploy Updated Code (Recommended)
1. Commit and push the updated `cards.go` code
2. Deploy new backend build that queries `scryfall_cards` directly
3. Apply migrations 009 and 010 properly
4. Remove the broken `cards` view dependency

---

## Timeline of Changes

| Date | Event |
|------|-------|
| Jan 11 | Production database initialized with migrations 001-008 |
| Jan 12 | Code updated to use Scryfall (commit ec5b1c50e81) |
| Jan 12 | `scryfall_cards` table manually imported (migration 009 not recorded) |
| Jan 12 | `cards` view created with BROKEN UUID conversion |
| Jan 20 09:25 | Local code updated to remove XMage dependencies (UNCOMMITTED) |
| Jan 20 11:00 | Production backend rebuilt (still using OLD code from Jan 12) |
| Jan 20 12:22 | Card lookup failures discovered |

---

## Recommendations

1. **Immediate fix:** Recreate the `cards` view with proper UUID conversion (see Option 1)
2. **Long-term fix:** Deploy updated code that queries `scryfall_cards` directly
3. **Migration hygiene:** Properly apply and record migrations 009 and 010
4. **Testing:** Add integration tests for card lookup before deployment

---

## Related Files

- `/Users/aron/dev/opensource/mage/mage-server-go/internal/repository/cards.go` (uncommitted changes)
- `/Users/aron/dev/opensource/mage/mage-server-go/migrations/009_create_scryfall_tables.up.sql`
- `/Users/aron/dev/opensource/mage/mage-server-go/migrations/010_drop_xmage_backup.up.sql`
- `/Users/aron/dev/opensource/mage/CARD_DATA_ARCHITECTURE.md`
- `/Users/aron/dev/opensource/mage/update-cards-prod.sh`

---

## Conclusion

The card lookup failures are caused by a **broken UUID→bigint conversion in the `cards` view**. The view was created during the Scryfall migration but failed to account for hyphens in PostgreSQL UUID format. Every card query fails with the same hexadecimal conversion error, causing all 67 cards in the deck to be rejected.

**Fix:** Replace the view with proper hyphen removal before conversion, or deploy the updated code that bypasses the view entirely.
