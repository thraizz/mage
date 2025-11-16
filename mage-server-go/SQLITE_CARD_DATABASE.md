# SQLite Card Database Strategy

## Why SQLite for Cards?

### Perfect Fit for Card Data

**Cards are read-only reference data:**
- ✅ 31,000+ cards loaded at startup
- ✅ Never modified during gameplay
- ✅ Same data across all servers
- ✅ Can be distributed as a single file
- ✅ No connection pooling needed
- ✅ Extremely fast for reads

**Benefits:**
1. **Portable** - Single `cards.db` file
2. **Fast** - In-process, no network overhead
3. **Simple** - No server setup required
4. **Distributable** - Ship with the server binary
5. **Cacheable** - Entire DB fits in memory (~50MB)
6. **Version-controlled** - Can track in git with LFS

### Architecture: Dual Database

```
┌─────────────────────────────────────────────────────┐
│ Go MAGE Server                                      │
│                                                     │
│  ┌──────────────────┐      ┌──────────────────┐   │
│  │ SQLite           │      │ PostgreSQL       │   │
│  │ (cards.db)       │      │ (game data)      │   │
│  │                  │      │                  │   │
│  │ - Cards          │      │ - Users          │   │
│  │ - Sets           │      │ - Stats          │   │
│  │ - Abilities      │      │ - Decks          │   │
│  │ (read-only)      │      │ - Games          │   │
│  │                  │      │ - Sessions       │   │
│  └──────────────────┘      └──────────────────┘   │
│         │                           │              │
│         │                           │              │
│         ▼                           ▼              │
│  ┌──────────────────────────────────────────────┐ │
│  │ Game Engine                                  │ │
│  │ - Loads cards from SQLite                    │ │
│  │ - Saves game state to PostgreSQL             │ │
│  └──────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

## Implementation Plan

### Step 1: Copy Java's H2 Database to SQLite

Java MAGE already has all cards in H2. We'll convert it to SQLite:

```bash
# Export from Java H2 to SQL
java -cp h2.jar org.h2.tools.Script \
  -url jdbc:h2:./Mage.Server/db/cards \
  -script cards.sql

# Convert H2 SQL to SQLite format
# (H2 and SQLite have slightly different SQL dialects)
./scripts/h2_to_sqlite.sh cards.sql cards.db
```

### Step 2: Create SQLite Card Repository

We'll create a separate repository for SQLite cards:

```go
// internal/repository/card_db.go
type CardDB struct {
    db     *sql.DB
    cache  *cardCache
    logger *zap.Logger
}

func NewCardDB(dbPath string, logger *zap.Logger) (*CardDB, error) {
    db, err := sql.Open("sqlite3", dbPath+"?mode=ro&cache=shared")
    if err != nil {
        return nil, err
    }
    
    // Enable optimizations
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(10)
    
    return &CardDB{
        db:     db,
        cache:  newCardCache(10000),
        logger: logger,
    }, nil
}
```

### Step 3: Preload All Cards at Startup

Since cards are read-only, load them all into memory:

```go
func (db *CardDB) PreloadAllCards(ctx context.Context) error {
    rows, err := db.db.QueryContext(ctx, "SELECT * FROM cards")
    if err != nil {
        return err
    }
    defer rows.Close()
    
    for rows.Next() {
        card := &Card{}
        // Scan into card
        db.cache.set(card.Name, card)
    }
    
    db.logger.Info("preloaded all cards", 
        zap.Int("count", len(db.cache.items)))
    return nil
}
```

## Database Schema

### Cards Table (SQLite)

```sql
CREATE TABLE cards (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    set_code TEXT NOT NULL,
    card_number TEXT,
    mana_cost TEXT,
    mana_value INTEGER,
    power TEXT,
    toughness TEXT,
    loyalty TEXT,
    defense TEXT,
    types TEXT,
    subtypes TEXT,
    supertypes TEXT,
    rules_text TEXT,
    flavor_text TEXT,
    rarity TEXT,
    card_class_name TEXT NOT NULL,
    frame_style TEXT,
    color_identity TEXT,
    is_split_card BOOLEAN,
    is_double_faced BOOLEAN,
    is_flip_card BOOLEAN,
    created_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- Indexes for fast lookups
CREATE INDEX idx_cards_name ON cards(name);
CREATE INDEX idx_cards_set_code ON cards(set_code);
CREATE INDEX idx_cards_class_name ON cards(card_class_name);
CREATE INDEX idx_cards_name_set ON cards(name, set_code);

-- Full-text search
CREATE VIRTUAL TABLE cards_fts USING fts5(
    name, 
    rules_text, 
    types, 
    subtypes,
    content=cards
);
```

### Ability Templates Table (SQLite)

For storing ability configurations:

```sql
CREATE TABLE card_abilities (
    id INTEGER PRIMARY KEY,
    card_id INTEGER NOT NULL,
    ability_type TEXT NOT NULL, -- 'activated', 'triggered', 'static', 'spell'
    ability_config TEXT NOT NULL, -- JSON configuration
    FOREIGN KEY (card_id) REFERENCES cards(id)
);

CREATE INDEX idx_abilities_card_id ON card_abilities(card_id);
```

Example ability config:
```json
{
  "type": "spell",
  "cost": "{R}",
  "effects": [
    {
      "type": "damage",
      "amount": 3,
      "target": {
        "type": "any",
        "count": 1
      }
    }
  ]
}
```

### Sets Table (SQLite)

```sql
CREATE TABLE sets (
    id INTEGER PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    release_date TEXT,
    set_type TEXT,
    block_name TEXT,
    has_boosters BOOLEAN,
    has_basic_lands BOOLEAN
);

CREATE INDEX idx_sets_code ON sets(code);
```

## Distribution Strategy

### Option 1: Ship Pre-built Database (Recommended)

**Pros:**
- Users don't need Java
- Instant startup
- Guaranteed consistency

**Process:**
1. Build `cards.db` from Java during release
2. Include in server distribution
3. Update with each new set release

```bash
# Release process
./scripts/build_card_database.sh  # Exports from Java, builds SQLite
cp cards.db mage-server-go/data/
./scripts/build_release.sh        # Includes cards.db in binary
```

### Option 2: Build on First Startup

**Pros:**
- Always up-to-date with Java
- Smaller distribution

**Cons:**
- Requires Java installation
- Slower first startup

```go
func (s *Server) ensureCardDatabase() error {
    if _, err := os.Stat("data/cards.db"); os.IsNotExist(err) {
        s.logger.Info("card database not found, building from Java...")
        if err := buildCardDatabaseFromJava(); err != nil {
            return err
        }
    }
    return nil
}
```

### Option 3: Download from CDN

**Pros:**
- Smallest distribution
- Easy updates

**Process:**
```go
func (s *Server) downloadCardDatabase() error {
    url := "https://cdn.mage.com/cards/v54/cards.db"
    resp, err := http.Get(url)
    // Download and verify checksum
}
```

## File Structure

```
mage-server-go/
├── data/
│   ├── cards.db              # SQLite card database (50MB)
│   ├── cards.db.sha256       # Checksum for verification
│   └── version.txt           # Database version
├── internal/
│   └── repository/
│       ├── card_db.go        # SQLite card repository
│       ├── card_loader.go    # Loads cards with abilities
│       └── db.go             # PostgreSQL for game data
```

## Performance Characteristics

### SQLite (Cards)
- **Size:** ~50MB for 31,000 cards
- **Startup:** ~100ms to open + preload
- **Query:** <1ms for indexed lookups
- **Memory:** ~100MB with full cache
- **Concurrent reads:** Excellent (read-only)

### PostgreSQL (Game Data)
- **Size:** Grows with user data
- **Startup:** Connection pool
- **Query:** <10ms for indexed lookups
- **Memory:** Configurable pool
- **Concurrent writes:** Excellent

## Migration from Current Setup

### Current (PostgreSQL Only)
```go
// All data in PostgreSQL
db := NewDB(ctx, cfg.Database, logger)
cardRepo := NewCardRepository(db, logger)
userRepo := NewUserRepository(db, logger)
```

### New (Dual Database)
```go
// Cards in SQLite (read-only)
cardDB := NewCardDB("data/cards.db", logger)
cardDB.PreloadAllCards(ctx)

// Game data in PostgreSQL (read-write)
gameDB := NewDB(ctx, cfg.Database, logger)
userRepo := NewUserRepository(gameDB, logger)
statsRepo := NewStatsRepository(gameDB, logger)

// Card loader combines SQLite data + ability system
cardLoader := NewCardLoader(cardDB, abilityRegistry, logger)
```

## Advantages for MVP

1. **No PostgreSQL setup for cards** - Just ship the file
2. **Faster development** - No migrations for card schema changes
3. **Portable** - Copy `cards.db` between environments
4. **Testable** - Easy to create test databases
5. **Versionable** - Track card DB versions explicitly
6. **Distributable** - Can host on CDN for updates

## Summary

**Recommendation: Use SQLite for cards**

✅ **Do:**
- Store all card metadata in SQLite
- Store ability configurations in SQLite
- Preload cards at startup
- Ship `cards.db` with server
- Use PostgreSQL for game data (users, stats, sessions)

❌ **Don't:**
- Store user data in SQLite
- Store game state in SQLite
- Modify cards.db at runtime
- Use SQLite for concurrent writes

**Next Steps:**
1. Create conversion script (H2 → SQLite)
2. Implement `CardDB` with SQLite
3. Build `cards.db` from Java
4. Update card loader to use SQLite
5. Keep PostgreSQL for game data
