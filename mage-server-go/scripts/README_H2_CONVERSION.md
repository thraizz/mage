# H2 Database Conversion Guide

This directory contains scripts to convert the Java MAGE H2 database to SQL format and import it into SQLite or PostgreSQL.

## Overview

```
Java H2 Database → SQL Export → SQLite/PostgreSQL
(cards.h2.mv.db)   (cards.sql)   (cards.db / mage database)
```

## Quick Start

### Option 1: SQLite (Recommended for Development)

```bash
# 1. Export H2 to SQL
./scripts/h2_to_sql.sh

# 2. Import to SQLite
./scripts/import_to_sqlite.sh

# 3. Verify
sqlite3 data/cards.db "SELECT COUNT(*) FROM card;"
```

**Result:** `data/cards.db` (portable, ~50MB)

### Option 2: PostgreSQL (Recommended for Production)

```bash
# 1. Export H2 to SQL
./scripts/h2_to_sql.sh

# 2. Import to PostgreSQL
./scripts/import_to_postgres.sh

# 3. Verify
PAGER=cat psql -d mage -c "SELECT COUNT(*) FROM card;"
```

**Result:** Cards in your existing PostgreSQL database

## Detailed Steps

### Step 1: Export H2 to SQL

```bash
cd mage-server-go
./scripts/h2_to_sql.sh
```

**What it does:**
1. Finds Java H2 database at `../Mage.Server/db/cards.h2.mv.db`
2. Downloads H2 JAR if needed
3. Exports to `data/cards.sql` using H2's Script tool

**Output:**
```
=== H2 to SQL Export ===
Source: /path/to/Mage.Server/db/cards.h2.mv.db
Output: /path/to/mage-server-go/data/cards.sql

Downloading H2 JAR...
✓ Downloaded H2 JAR
Exporting H2 database to SQL...
✓ Exported to SQL

=== Export Complete ===
✓ Output: data/cards.sql
✓ Size: 45M
✓ Lines: 234567
```

**Troubleshooting:**

If H2 database not found:
```bash
# Generate it by running Java MAGE server once
cd ../Mage.Server
mvn clean install
mvn exec:java
# Wait for server to start, then stop it (Ctrl+C)
# Database will be created at db/cards.h2.mv.db
```

### Step 2a: Import to SQLite

```bash
./scripts/import_to_sqlite.sh
```

**What it does:**
1. Reads `data/cards.sql`
2. Converts H2 SQL to SQLite format
3. Creates `data/cards.db`
4. Creates indexes for performance

**Output:**
```
=== Import to SQLite ===
Input: data/cards.sql
Output: data/cards.db

Converting SQL to SQLite format...
✓ Conversion complete
  Total lines: 234567
  Converted: 234000
  Skipped: 567
  CREATE TABLE: 2
  INSERT: 31234

Importing to SQLite...
✓ Import complete

=== Database Statistics ===
Cards: 31234
Sets: 567
Size: 48M

Sample cards:
Lightning Bolt|LEA|{R}
Lightning Bolt|LEB|{R}
Lightning Bolt|2ED|{R}

Creating indexes...
✓ Indexes created

Database ready: data/cards.db
```

**Verify:**
```bash
# Count cards
sqlite3 data/cards.db "SELECT COUNT(*) FROM card;"

# Search for a card
sqlite3 data/cards.db "SELECT name, setCode, manaCosts FROM card WHERE name = 'Lightning Bolt';"

# List all tables
sqlite3 data/cards.db ".tables"

# Show schema
sqlite3 data/cards.db ".schema card"
```

### Step 2b: Import to PostgreSQL

```bash
# Set connection details (optional, defaults shown)
export DB_NAME=mage
export DB_USER=postgres
export DB_HOST=localhost
export DB_PORT=5432

# Import
./scripts/import_to_postgres.sh
```

**What it does:**
1. Reads `data/cards.sql`
2. Converts H2 SQL to PostgreSQL format
3. Imports into existing `mage` database
4. Creates indexes for performance

**Output:**
```
=== Import to PostgreSQL ===
Input: data/cards.sql
Database: mage
Host: localhost:5432
User: postgres

Testing database connection...
✓ Connected to database
Converting SQL to PostgreSQL format...
✓ Conversion complete
  CREATE TABLE: 2
  INSERT: 31234
  CREATE INDEX: 5
  CREATE SEQUENCE: 2
  ALTER SEQUENCE: 2
  Other: 0
  Output: data/cards.sql.postgres.tmp

Importing to PostgreSQL...
✓ Import complete

=== Database Statistics ===
Cards: 31234
Sets: 567

Sample cards:
     name      | setcode | manacosts 
---------------+---------+-----------
 Lightning Bolt| lea     | {R}
 Lightning Bolt| leb     | {R}
 Lightning Bolt| 2ed     | {R}

Creating indexes...
✓ Indexes created

Database ready!
```

**Verify:**
```bash
# Count cards
PAGER=cat psql -d mage -c "SELECT COUNT(*) FROM card;"

# Search for a card
PAGER=cat psql -d mage -c "SELECT name, setcode, manacosts FROM card WHERE name = 'Lightning Bolt';"

# List tables
PAGER=cat psql -d mage -c "\dt"

# Show schema
PAGER=cat psql -d mage -c "\d card"
```

## Database Schema

### Main Tables

**`card` table:**
```sql
- name (TEXT) - Card name
- setcode (TEXT) - Set code (e.g., "ISD")
- cardnumber (TEXT) - Card number in set
- classname (TEXT) - Java class name
- manacosts (TEXT) - Mana cost (e.g., "{2}{U}")
- manavalue (INTEGER) - Converted mana cost
- power (TEXT) - Power (for creatures)
- toughness (TEXT) - Toughness (for creatures)
- startingloyalty (TEXT) - Loyalty (for planeswalkers)
- types (TEXT) - Card types
- subtypes (TEXT) - Subtypes
- supertypes (TEXT) - Supertypes
- rules (TEXT) - Rules text
- rarity (TEXT) - Rarity
- ... (30+ more fields)
```

**`expansion` table:**
```sql
- code (TEXT) - Set code
- name (TEXT) - Set name
- releasedate (TEXT) - Release date
- settype (TEXT) - Set type
- blockname (TEXT) - Block name
- ... (more fields)
```

### Indexes Created

```sql
CREATE INDEX idx_card_name ON card(name);
CREATE INDEX idx_card_setcode ON card(setcode);
CREATE INDEX idx_card_classname ON card(classname);
CREATE INDEX idx_card_name_set ON card(name, setcode);
CREATE INDEX idx_card_types ON card(types);
CREATE INDEX idx_card_manavalue ON card(manavalue);
```

## Using in Go

### SQLite

```go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

// Open database
db, err := sql.Open("sqlite3", "data/cards.db?mode=ro&cache=shared")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Query cards
rows, err := db.Query("SELECT name, manacosts FROM card WHERE name = ?", "Lightning Bolt")
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var name, manaCosts string
    rows.Scan(&name, &manaCosts)
    fmt.Printf("%s: %s\n", name, manaCosts)
}
```

### PostgreSQL

```go
import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
)

// Connect to database
pool, err := pgxpool.New(context.Background(), 
    "postgres://postgres:postgres@localhost:5432/mage")
if err != nil {
    log.Fatal(err)
}
defer pool.Close()

// Query cards
rows, err := pool.Query(context.Background(),
    "SELECT name, manacosts FROM card WHERE name = $1", "Lightning Bolt")
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var name, manaCosts string
    rows.Scan(&name, &manaCosts)
    fmt.Printf("%s: %s\n", name, manaCosts)
}
```

## Conversion Details

### H2 → SQLite

**Data type conversions:**
- `VARCHAR_IGNORECASE` → `TEXT`
- `VARCHAR(n)` → `TEXT`
- `INTEGER` → `INTEGER`
- `BOOLEAN` → `INTEGER` (0/1)
- `TIMESTAMP` → `INTEGER`

**Boolean conversions:**
- `TRUE` → `1`
- `FALSE` → `0`

**Removed:**
- `CACHED` clause
- `NOT PERSISTENT` clause
- Sequence statements
- User/schema creation

### H2 → PostgreSQL

**Data type conversions:**
- `VARCHAR_IGNORECASE` → `TEXT`
- `BOOLEAN` → `BOOLEAN` (kept as-is)
- `DOUBLE` → `DOUBLE PRECISION`
- `BLOB` → `BYTEA`

**Name conversions:**
- Table names → lowercase
- Column names → lowercase
- (PostgreSQL convention)

**Sequence conversions:**
- `START WITH` → `START`

## Troubleshooting

### H2 database not found

**Problem:** `ERROR: Java H2 database not found`

**Solution:**
```bash
cd ../Mage.Server
mvn clean install
mvn exec:java
# Wait for startup, then Ctrl+C
cd ../mage-server-go
./scripts/h2_to_sql.sh
```

### H2 JAR download fails

**Problem:** Cannot download H2 JAR

**Solution:**
```bash
# Download manually
curl -L -o scripts/h2.jar \
  https://repo1.maven.org/maven2/com/h2database/h2/2.2.224/h2-2.2.224.jar
```

### SQLite not installed

**Problem:** `ERROR: sqlite3 not found`

**Solution:**
```bash
# macOS
brew install sqlite3

# Ubuntu/Debian
sudo apt-get install sqlite3

# Verify
sqlite3 --version
```

### PostgreSQL connection fails

**Problem:** `ERROR: Cannot connect to PostgreSQL`

**Solution:**
```bash
# Check PostgreSQL is running
pg_isready

# Create database if needed
createdb mage

# Test connection
PAGER=cat psql -d mage -c "SELECT 1;"

# Set password if needed
export PGPASSWORD=your_password
```

### Import fails with encoding errors

**Problem:** `UnicodeDecodeError` or encoding issues

**Solution:**
```bash
# The scripts use UTF-8 encoding
# Ensure your terminal supports UTF-8

# macOS/Linux
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8

# Re-run import
./scripts/import_to_sqlite.sh
```

## File Structure

```
mage-server-go/
├── scripts/
│   ├── h2_to_sql.sh                  # Export H2 to SQL
│   ├── import_to_sqlite.sh           # Import to SQLite
│   ├── import_to_postgres.sh         # Import to PostgreSQL
│   ├── convert_h2_to_sqlite.py       # H2→SQLite converter
│   ├── convert_h2_to_postgres.py     # H2→PostgreSQL converter
│   ├── h2.jar                        # H2 database JAR (auto-downloaded)
│   └── README_H2_CONVERSION.md       # This file
├── data/
│   ├── cards.sql                     # Exported SQL (intermediate)
│   ├── cards.db                      # SQLite database (if using SQLite)
│   └── cards.db.sha256               # Checksum (if using SQLite)
```

## Performance

### SQLite
- **Size:** ~48MB for 31,000 cards
- **Import time:** ~30 seconds
- **Query time:** <1ms for indexed lookups
- **Startup:** ~10ms to open + load

### PostgreSQL
- **Size:** ~60MB for 31,000 cards (with indexes)
- **Import time:** ~45 seconds
- **Query time:** <5ms for indexed lookups
- **Startup:** Connection pool

## Updating Cards

When Java MAGE adds new cards:

```bash
# 1. Update Java MAGE
cd ../Mage
git pull
mvn clean install

# 2. Run Java server to regenerate H2 database
cd ../Mage.Server
mvn exec:java
# Ctrl+C after startup

# 3. Re-export and re-import
cd ../mage-server-go
./scripts/h2_to_sql.sh
./scripts/import_to_sqlite.sh  # or import_to_postgres.sh
```

## Next Steps

After importing cards:

1. **Implement card repository:**
   ```go
   // internal/repository/card_db.go
   type CardDB struct {
       db *sql.DB
   }
   ```

2. **Build card loader:**
   ```go
   // internal/game/card_loader.go
   func LoadCard(name string) (*Card, error)
   ```

3. **Implement abilities:**
   See `ABILITY_IMPLEMENTATION_STRATEGY.md`

4. **Test in game engine:**
   Load cards and play games!
