# H2 to SQL Conversion - Quick Start

## TL;DR

Convert Java MAGE's H2 card database to SQL format for use in your Go server.

## One-Line Commands

### For SQLite (Recommended)

```bash
./scripts/h2_to_sql.sh && ./scripts/import_to_sqlite.sh
```

**Result:** `data/cards.db` with 31,000+ cards

### For PostgreSQL

```bash
./scripts/h2_to_sql.sh && ./scripts/import_to_postgres.sh
```

**Result:** Cards in your `mage` PostgreSQL database

## What You Get

### Card Data
- ✅ 31,000+ Magic cards
- ✅ All card properties (name, cost, types, P/T, etc.)
- ✅ Rules text for each card
- ✅ Set information
- ✅ Rarity, colors, etc.

### What's NOT Included
- ❌ Ability implementations (you'll build these in Go)
- ❌ Card images (separate download)
- ❌ Game logic (you'll implement this)

## Prerequisites

### Required
- Java MAGE server (to generate H2 database)
- Python 3 (for conversion scripts)
- SQLite3 or PostgreSQL

### Check Prerequisites

```bash
# Check if Java H2 database exists
ls -lh ../Mage.Server/db/cards.h2.mv.db

# If not found, generate it:
cd ../Mage.Server
mvn exec:java
# Wait for startup, then Ctrl+C

# Check Python
python3 --version

# Check SQLite (if using SQLite)
sqlite3 --version

# Check PostgreSQL (if using PostgreSQL)
psql --version
```

## Step-by-Step

### Step 1: Export H2 to SQL (2 minutes)

```bash
cd mage-server-go
./scripts/h2_to_sql.sh
```

**Output:** `data/cards.sql` (~45MB)

### Step 2: Import to Database (1 minute)

**Option A: SQLite**
```bash
./scripts/import_to_sqlite.sh
```

**Option B: PostgreSQL**
```bash
./scripts/import_to_postgres.sh
```

### Step 3: Verify (30 seconds)

**SQLite:**
```bash
sqlite3 data/cards.db "SELECT COUNT(*) FROM card;"
# Should show: 31234 (or similar)

sqlite3 data/cards.db "SELECT name, manacosts FROM card WHERE name = 'Lightning Bolt' LIMIT 3;"
# Should show Lightning Bolt entries
```

**PostgreSQL:**
```bash
PAGER=cat psql -d mage -c "SELECT COUNT(*) FROM card;"
# Should show: 31234 (or similar)

PAGER=cat psql -d mage -c "SELECT name, manacosts FROM card WHERE name = 'Lightning Bolt' LIMIT 3;"
# Should show Lightning Bolt entries
```

## What's in the Database?

### Tables

**`card`** - All Magic cards
- `name` - Card name (e.g., "Lightning Bolt")
- `setcode` - Set code (e.g., "LEA", "ISD")
- `cardnumber` - Card number in set
- `classname` - Java class name (for ability lookup)
- `manacosts` - Mana cost (e.g., "{2}{U}")
- `manavalue` - Converted mana cost
- `power`, `toughness` - Creature stats
- `types`, `subtypes`, `supertypes` - Card types
- `rules` - Rules text
- `rarity` - Rarity (COMMON, UNCOMMON, RARE, MYTHIC)
- ... and 20+ more fields

**`expansion`** - Set information
- `code` - Set code
- `name` - Set name
- `releasedate` - Release date
- `settype` - Set type (EXPANSION, CORE, etc.)

### Example Queries

```sql
-- Find all Lightning Bolt printings
SELECT name, setcode, cardnumber, manacosts 
FROM card 
WHERE name = 'Lightning Bolt';

-- Find all cards in Innistrad
SELECT name, types, rarity 
FROM card 
WHERE setcode = 'ISD'
ORDER BY cardnumber;

-- Find all creatures with power 5+
SELECT name, power, toughness, manacosts
FROM card
WHERE types LIKE '%Creature%'
  AND CAST(power AS INTEGER) >= 5;

-- Find all cards with "flying"
SELECT name, types, rules
FROM card
WHERE rules LIKE '%flying%'
LIMIT 10;

-- Count cards by rarity
SELECT rarity, COUNT(*) 
FROM card 
GROUP BY rarity;
```

## Using in Go

### SQLite Example

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // Open database (read-only)
    db, err := sql.Open("sqlite3", "data/cards.db?mode=ro")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Query a card
    var name, manaCosts string
    err = db.QueryRow(
        "SELECT name, manacosts FROM card WHERE name = ? LIMIT 1",
        "Lightning Bolt",
    ).Scan(&name, &manaCosts)
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("%s costs %s\n", name, manaCosts)
    // Output: Lightning Bolt costs {R}
}
```

### PostgreSQL Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    // Connect to database
    pool, err := pgxpool.New(context.Background(),
        "postgres://postgres:postgres@localhost:5432/mage")
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()
    
    // Query a card
    var name, manaCosts string
    err = pool.QueryRow(context.Background(),
        "SELECT name, manacosts FROM card WHERE name = $1 LIMIT 1",
        "Lightning Bolt",
    ).Scan(&name, &manaCosts)
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("%s costs %s\n", name, manaCosts)
    // Output: Lightning Bolt costs {R}
}
```

## Troubleshooting

### "H2 database not found"

```bash
# Generate it by running Java MAGE server
cd ../Mage.Server
mvn exec:java
# Wait for "Server started", then Ctrl+C
```

### "sqlite3 not found"

```bash
# macOS
brew install sqlite3

# Ubuntu
sudo apt-get install sqlite3
```

### "Cannot connect to PostgreSQL"

```bash
# Check if PostgreSQL is running
pg_isready

# Create database if needed
createdb mage

# Test connection
PAGER=cat psql -d mage -c "SELECT 1;"
```

### "Python not found"

```bash
# macOS (should be pre-installed)
python3 --version

# Ubuntu
sudo apt-get install python3
```

## Next Steps

After importing cards:

1. ✅ **You have card metadata** - All 31,000+ cards in database
2. ⬜ **Implement ability system** - See `ABILITY_IMPLEMENTATION_STRATEGY.md`
3. ⬜ **Build card loader** - Load cards from DB into game engine
4. ⬜ **Port abilities** - Implement card abilities in Go
5. ⬜ **Test in engine** - Play games with real cards

## Files Created

```
mage-server-go/
├── data/
│   ├── cards.sql          # Intermediate SQL export
│   └── cards.db           # SQLite database (if using SQLite)
└── scripts/
    ├── h2_to_sql.sh       # Export script
    ├── import_to_sqlite.sh # SQLite import
    └── import_to_postgres.sh # PostgreSQL import
```

## Performance

- **Export time:** ~30 seconds
- **Import time:** ~30-60 seconds
- **Database size:** ~48MB (SQLite) or ~60MB (PostgreSQL)
- **Query speed:** <1ms for indexed lookups
- **Total cards:** 31,000+

## Summary

**What you did:**
1. Exported Java's H2 database to SQL
2. Converted to SQLite/PostgreSQL format
3. Imported into your database

**What you have now:**
- ✅ All Magic card metadata
- ✅ Fast, indexed database
- ✅ Ready for Go integration

**What's next:**
- Implement ability system
- Load cards into game engine
- Start playing games!

## Questions?

- **Full documentation:** `scripts/README_H2_CONVERSION.md`
- **Ability strategy:** `ABILITY_IMPLEMENTATION_STRATEGY.md`
- **Card architecture:** `CARD_DATA_ARCHITECTURE.md`
- **Practical guide:** `PRACTICAL_CARD_STRATEGY.md`
