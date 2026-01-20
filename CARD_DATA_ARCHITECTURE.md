# Card Data Architecture

**Last Updated**: 2026-01-20
**Status**: ✅ Production Ready

This document describes the card data architecture for the Mage project.

---

## Overview

Mage uses **Scryfall** as the exclusive source for Magic: The Gathering card data. All card metadata (names, types, mana costs, oracle text, etc.) is stored in the `scryfall_cards` PostgreSQL table.

### Key Principles

1. **Single Source of Truth**: All card data comes from Scryfall bulk data API
2. **No Legacy Dependencies**: XMage-specific code and data have been completely removed
3. **Regular Updates**: Card database is updated via automated downloads from Scryfall
4. **Complete Coverage**: 100k+ English cards from all Magic sets

---

## Database Schema

### Primary Table: `scryfall_cards`

Created in migration `009_create_scryfall_tables.up.sql`.

```sql
CREATE TABLE scryfall_cards (
    id UUID PRIMARY KEY,                    -- Scryfall card ID
    oracle_id UUID,                         -- Oracle ID (same across printings)
    name TEXT NOT NULL,                     -- Card name
    lang VARCHAR(10) DEFAULT 'en',          -- Language code
    released_at DATE,                       -- Release date
    uri TEXT,                               -- Scryfall API URI
    scryfall_uri TEXT,                      -- Scryfall web page
    layout VARCHAR(50),                     -- Card layout (normal, split, etc.)
    mana_cost TEXT,                         -- Mana cost string
    cmc NUMERIC,                            -- Converted mana cost
    type_line TEXT,                         -- Full type line
    oracle_text TEXT,                       -- Oracle rules text
    power VARCHAR(10),                      -- Creature power
    toughness VARCHAR(10),                  -- Creature toughness
    loyalty VARCHAR(10),                    -- Planeswalker loyalty
    defense VARCHAR(10),                    -- Battle defense
    colors JSONB,                           -- Card colors array
    color_identity JSONB,                   -- Commander color identity
    keywords JSONB,                         -- Keyword abilities
    legalities JSONB,                       -- Format legalities
    set_code VARCHAR(10),                   -- Set code (e.g., "mkm")
    set_name TEXT,                          -- Set name
    collector_number VARCHAR(50),           -- Collector number
    rarity VARCHAR(20),                     -- Rarity
    image_uris JSONB,                       -- Image URLs
    card_faces JSONB,                       -- Multi-faced card data
    created_at TIMESTAMP DEFAULT NOW()
);
```

**Indexes:**
- Primary key on `id`
- Index on `name` for fast lookups
- Index on `oracle_id` for finding all printings
- Index on `set_code` for set queries
- Full-text search indexes using pg_trgm

---

## Data Flow

### 1. Download Scryfall Bulk Data

```bash
cd mage-server-go
./scripts/download_scryfall_bulk.sh
```

This script:
- Downloads the latest Scryfall bulk data JSON (~200MB compressed)
- Extracts to `data/scryfall-all-cards.json`
- Verifies data integrity

### 2. Import to Database

```bash
cd mage-server-go
go run ./cmd/import-scryfall/main.go
```

The import process:
- Reads `data/scryfall-all-cards.json`
- Truncates `scryfall_cards` table
- Bulk inserts all card data
- Rebuilds indexes
- Takes ~2-3 minutes for 100k+ cards

### 3. Query from Application

The application queries via `CardRepository` in `internal/repository/cards.go`:

```go
// Get card by name
cards, err := cardRepo.GetByName(ctx, "Lightning Bolt")

// Search cards
results, err := cardRepo.SearchByName(ctx, "bolt", 10)

// Get cards by set
setCards, err := cardRepo.GetBySetCode(ctx, "lea")
```

---

## Field Mappings

### Scryfall → Repository Card Struct

| Scryfall Field | Card Struct Field | Notes |
|---------------|------------------|-------|
| `id` | `ID` | UUID converted to int64 via bit conversion |
| `name` | `Name` | Card name |
| `collector_number` | `CardNumber` | Collector number as string |
| `set_code` | `SetCode` | Set code (uppercase) |
| `type_line` | `CardType` | Full type line |
| `mana_cost` | `ManaCost` | Mana cost string (e.g., "{2}{U}") |
| `power` | `Power` | Creature power |
| `toughness` | `Toughness` | Creature toughness |
| `oracle_text` | `RulesText` | Oracle rules text |
| `rarity` | `Rarity` | common/uncommon/rare/mythic |

**Nullable fields** mapped to `sql.NullString` or `sql.NullInt64`:
- `FlavorText` - Always NULL (not in bulk data)
- `OriginalText` - Always NULL
- `OriginalType` - Always NULL
- `CN` - Collector number as integer (when numeric)

---

## Query Patterns

### Basic Queries

```sql
-- Get card by exact name
SELECT * FROM scryfall_cards
WHERE lang = 'en' AND name = 'Lightning Bolt'
ORDER BY released_at DESC;

-- Search by partial name
SELECT * FROM scryfall_cards
WHERE lang = 'en' AND name ILIKE '%bolt%'
ORDER BY name;

-- Get all cards from a set
SELECT * FROM scryfall_cards
WHERE lang = 'en' AND set_code = 'lea'
ORDER BY collector_number;
```

### Multi-Faced Cards

Split/DFC cards are handled by checking the name pattern:

```sql
-- Handles "Fire // Ice" and "Fire" queries
SELECT * FROM scryfall_cards
WHERE lang = 'en'
  AND (name = 'Fire' OR name LIKE 'Fire //%');
```

### Case-Insensitive Search

```go
// Normalizes Unicode apostrophes and performs case-insensitive match
cards, err := cardRepo.GetByNameCaseInsensitive(ctx, "Jace's Ingenuity")
```

---

## Production Updates

### Updating Card Data in Production

1. **Download latest Scryfall data:**
   ```bash
   ./update-cards-prod.sh
   ```

2. **What it does:**
   - SSHs into production server
   - Downloads latest Scryfall bulk data
   - Runs import process
   - Verifies import succeeded
   - Restarts services if needed

3. **Frequency:**
   - Run after new set releases
   - Run monthly to catch oracle text updates
   - Run after ban/restriction announcements

### Rollback

Database migration 010 can be rolled back:

```bash
cd mage-server-go
migrate -path migrations -database "postgres://..." down 1
```

**Note:** Rollback recreates the `cards_xmage_backup` table structure but **does not restore data**. The rollback is primarily for development/testing.

---

## Benefits Over XMage Data

### Completeness
- **Scryfall**: 101,852 English cards (Jan 2026)
- **XMage**: 87,962 cards (missing recent sets)

### Data Quality
- **Official oracle text**: Directly from Scryfall API
- **Consistent formatting**: Standardized type lines, mana costs
- **Regular updates**: Weekly Scryfall bulk data updates

### Features
- **Multi-faced cards**: Proper support for DFC, MDFC, split, flip, etc.
- **Card faces**: Individual face data for transform cards
- **Legalities**: Format-legal status for all formats
- **Images**: High-resolution card images via image_uris

### Simplicity
- **Single table**: No complex JOINs or compatibility views
- **No class names**: No XMage Java class mapping
- **Standard schema**: Matches Scryfall API structure

---

## Removed Components (2026-01-20)

The following XMage-specific code was removed:

### Repository Layer
- ❌ `CardRepository.GetByClassName()` - Queried XMage backup table
- ❌ `Card.CardClassName` field - XMage Java class name

### Game Engine
- ❌ `Factory.CreateCardByClassName()` - Created cards by Java class
- ❌ `cardRegistry.GetByClassName()` - Registry lookup by class
- ❌ `CardInfo.CardClassName` - Card metadata field
- ❌ `game.Card.CardClassName` - Runtime card field

### Database
- ❌ `cards_xmage_backup` table - Dropped in migration 010
- ❌ `cards` table compatibility view - No longer needed

### Scripts
- ❌ `migrate_to_scryfall.sh` - Archived (migration complete)
- ❌ `rollback_from_scryfall.sh` - Archived (no rollback path)

---

## Migration History

1. **Migration 002** (Legacy): Created original `cards` table from XMage H2 database
2. **Migration 009** (Jan 2026): Created `scryfall_cards` table and imported Scryfall data
3. **Migration 010** (Jan 2026): Dropped `cards_xmage_backup` table, completed XMage removal

See [SCRYFALL_PRIMARY_MIGRATION.md](SCRYFALL_PRIMARY_MIGRATION.md) (archived) for historical context.

---

## References

- **Scryfall Bulk Data**: https://scryfall.com/docs/api/bulk-data
- **Download Script**: `mage-server-go/scripts/download_scryfall_bulk.sh`
- **Import Tool**: `mage-server-go/cmd/import-scryfall/`
- **Repository Code**: `mage-server-go/internal/repository/cards.go`
- **Migration Scripts**: `mage-server-go/migrations/`
