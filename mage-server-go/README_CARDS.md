# Card System Overview

This document provides a high-level overview of the card system architecture and how to get started.

## Quick Links

- **[CARD_DATA_ARCHITECTURE.md](CARD_DATA_ARCHITECTURE.md)** - Detailed explanation of how cards work and population strategies
- **[QUICK_START_CARDS.md](QUICK_START_CARDS.md)** - Step-by-step guide to populate your database
- **[internal/game/ability_example.go](internal/game/ability_example.go)** - Reference implementation of ability system

## Current Status

### ✅ What's Ready

1. **Database Schema** (`migrations/002_create_cards_table.up.sql`)
   - Cards table with all necessary fields
   - Indexes for fast queries
   - Full-text search support

2. **Repository Layer** (`internal/repository/cards.go`)
   - `GetByID(id)` - Fetch card by ID
   - `GetByName(name)` - Fetch all printings of a card
   - `SearchByName(term)` - Full-text search
   - `GetBySetCode(code)` - Get all cards from a set
   - `Create(card)` - Add new cards
   - In-memory caching (10,000 cards)

3. **Import Scripts**
   - `scripts/export_java_cards.sh` - Export from Java H2 database
   - `scripts/import_cards.go` - Import to PostgreSQL

### ⚠️ What's Missing

1. **Card Data** - Database is empty, needs population
2. **Ability System** - Cards have metadata but no game logic
3. **Card Loader** - Can't instantiate cards for games yet
4. **Deck System** - Can't load decks into games

## How Java MAGE Works

```
┌─────────────────────────────────────────────────┐
│ Java Card Class (31,000+ files)                │
│                                                 │
│  public class LightningBolt extends CardImpl { │
│    public LightningBolt(...) {                 │
│      super(..., "{R}");                        │
│      this.getSpellAbility().addEffect(         │
│        new DamageTargetEffect(3)               │
│      );                                         │
│    }                                            │
│  }                                              │
└─────────────────────────────────────────────────┘
                    │
                    │ Reflection at startup
                    ▼
┌─────────────────────────────────────────────────┐
│ H2 Database (cards.h2.mv.db)                   │
│                                                 │
│  - Metadata: name, cost, P/T, types            │
│  - Rules text (extracted from abilities)       │
│  - Class name for instantiation                │
└─────────────────────────────────────────────────┘
                    │
                    │ During game
                    ▼
┌─────────────────────────────────────────────────┐
│ Card Instance (in-game object)                 │
│                                                 │
│  - Metadata from database                      │
│  - Abilities from Java class constructor       │
│  - Game state (tapped, damage, counters)       │
└─────────────────────────────────────────────────┘
```

**Key Insight:** Java stores metadata in database, but abilities are in code.

## Go Engine Strategy

### Phase 1: Database Population (Now)

**Goal:** Get card metadata into PostgreSQL

```bash
# 1. Export from Java
./scripts/export_java_cards.sh

# 2. Import to PostgreSQL
go run scripts/import_cards.go

# 3. Verify
PAGER=cat psql -d mage -c "SELECT COUNT(*) FROM cards;"
```

**Result:** 31,000+ cards in database with metadata

### Phase 2: Test Cards (Week 1-2)

**Goal:** Create 10-20 hardcoded cards for testing

```go
// internal/game/cards/test_cards.go
func NewLightningBolt(id uuid.UUID) *CardInstance {
    card := &CardInstance{
        Name:     "Lightning Bolt",
        ManaCost: "{R}",
        Types:    []CardType{TypeInstant},
    }
    
    // Add damage ability
    card.AddAbility(&DamageAbility{
        Amount: 3,
        Target: TargetAny,
    })
    
    return card
}
```

**Result:** Can test basic game mechanics

### Phase 3: Ability System (Week 3-4)

**Goal:** Design flexible ability framework

See `internal/game/ability_example.go` for reference implementation.

Key components:
- `Ability` interface (activated, triggered, static)
- `Effect` interface (damage, draw, counter, etc.)
- `Cost` interface (mana, tap, sacrifice, etc.)
- `Target` interface (card, player, spell)

**Result:** Can implement most common abilities

### Phase 4: Card Loader (Week 5-6)

**Goal:** Load cards from database with abilities

```go
// internal/game/card_loader.go
func (l *CardLoader) LoadCard(name string) (*CardInstance, error) {
    // Get metadata from database
    cardInfo, _ := l.repo.GetByName(ctx, name)
    
    // Create instance
    card := &CardInstance{
        Name:      cardInfo.Name,
        ManaCost:  cardInfo.ManaCost,
        Types:     parseTypes(cardInfo.CardType),
        // ...
    }
    
    // Attach abilities (from registry or templates)
    abilities := l.loadAbilities(cardInfo.CardClassName)
    card.Abilities = abilities
    
    return card, nil
}
```

**Result:** Can load any card from database

### Phase 5: Priority Cards (Month 2-3)

**Goal:** Port 200 most-used cards manually

Focus on:
- Basic lands (5)
- Simple creatures (50)
- Simple instants/sorceries (50)
- Common keywords (flying, trample, etc.)
- Standard format staples (100)

**Result:** Can play meaningful games

### Phase 6: Full Port (Month 4+)

**Goal:** All 31,000+ cards playable

Options:
1. Manual port (slow but accurate)
2. Java-to-Go transpiler (fast but complex)
3. Hybrid (templates + custom code)

**Result:** Full parity with Java

## Architecture Comparison

### Java Approach
- ✅ Every card is a class
- ✅ Type-safe, compiled
- ✅ Full control over behavior
- ❌ Massive codebase (31,000+ files)
- ❌ Requires recompile for changes
- ❌ Hard to auto-generate

### Go Options

#### Option A: Port Classes (Like Java)
```go
// internal/cards/lightning_bolt.go
func NewLightningBolt(id uuid.UUID) *CardInstance {
    // Full control, type-safe
}
```
- ✅ Full control
- ✅ Type-safe
- ❌ Massive effort

#### Option B: Ability Templates (Hybrid)
```json
{
  "name": "Lightning Bolt",
  "abilities": [
    {
      "type": "spell",
      "effects": [{"type": "damage", "amount": 3}]
    }
  ]
}
```
- ✅ Scalable
- ✅ No recompile
- ❌ Limited to templates

#### Option C: Mixed Approach (Recommended)
- Hardcode complex cards
- Templates for simple cards
- Generate common patterns

## Database Schema

### Cards Table (Existing)

```sql
CREATE TABLE cards (
    id SERIAL PRIMARY KEY,
    card_number VARCHAR(50),
    set_code VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    card_type VARCHAR(255),
    mana_cost VARCHAR(255),
    power VARCHAR(10),
    toughness VARCHAR(10),
    rules_text TEXT,
    rarity VARCHAR(50),
    card_class_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Potential Additions

```sql
-- Store ability configurations
ALTER TABLE cards ADD COLUMN abilities_json JSONB;

-- Deck storage
CREATE TABLE decks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    name VARCHAR(255),
    format VARCHAR(50)
);

CREATE TABLE deck_cards (
    deck_id INTEGER REFERENCES decks(id),
    card_id INTEGER REFERENCES cards(id),
    quantity INTEGER,
    is_sideboard BOOLEAN
);
```

## API Usage Examples

### Query Cards

```go
// Get card by name
card, err := cardRepo.GetByName(ctx, "Lightning Bolt")

// Search cards
cards, err := cardRepo.SearchByName(ctx, "delver", 10)

// Get all cards from a set
cards, err := cardRepo.GetBySetCode(ctx, "ISD")
```

### Load Card for Game

```go
// Create card instance
card, err := cardLoader.LoadCard("Lightning Bolt", "LEA")

// Card has metadata + abilities
fmt.Println(card.Name)      // "Lightning Bolt"
fmt.Println(card.ManaCost)  // "{R}"
fmt.Println(len(card.Abilities)) // 1

// Execute ability
ability := card.Abilities[0]
if ability.CanExecute(game, card) {
    ability.Execute(game, card, targets)
}
```

### Validate Deck

```go
// Load deck from list
deck, err := deckLoader.LoadDeck(deckList)

// Validate for format
err = deckValidator.ValidateDeck(deck, "Standard")
// Checks: size, legality, card limits
```

## Next Steps

### Immediate (This Week)
1. ✅ Read this document
2. ✅ Read QUICK_START_CARDS.md
3. ⬜ Run export script
4. ⬜ Run import script
5. ⬜ Verify card data

### Short-term (Next 2 Weeks)
1. ⬜ Review ability_example.go
2. ⬜ Design ability system for your needs
3. ⬜ Implement 5-10 test cards
4. ⬜ Test in game engine

### Medium-term (Next 2 Months)
1. ⬜ Implement ability templates
2. ⬜ Port 50 priority cards
3. ⬜ Build card loader
4. ⬜ Integrate with game engine

### Long-term (Next 6 Months)
1. ⬜ Port 200+ cards
2. ⬜ Build transpiler (optional)
3. ⬜ Achieve parity with Java
4. ⬜ Maintain for new sets

## FAQ

**Q: Do I need to populate cards to run the server?**  
A: No, but you can't play games without cards.

**Q: Can I use the server with just metadata?**  
A: Yes, for deck validation, card search, collection management.

**Q: How long to populate the database?**  
A: ~10 minutes (export + import)

**Q: How long to make cards playable?**  
A: 2-4 weeks for basic ability system + test cards

**Q: Do I need all 31,000 cards?**  
A: No, start with 20-50 for testing

**Q: Can I auto-generate from Java?**  
A: Possible but complex. Start manual, build tooling later.

**Q: What about new sets?**  
A: Re-export from Java, import new cards

**Q: How does Java handle errata?**  
A: Update card class, regenerate database

## Resources

- **Java Card Repository:** `Mage/src/main/java/mage/cards/repository/CardRepository.java`
- **Java Card Scanner:** `Mage/src/main/java/mage/cards/repository/CardScanner.java`
- **Java Card Classes:** `Mage.Sets/src/mage/sets/*/` (31,000+ files)
- **Go Card Repository:** `internal/repository/cards.go`
- **Go Ability Example:** `internal/game/ability_example.go`

## Summary

**Current State:**
- ✅ Database schema ready
- ✅ Repository layer ready
- ⚠️ Database empty (run scripts)
- ❌ Ability system needed
- ❌ Card loader needed

**To Populate:**
```bash
./scripts/export_java_cards.sh
go run scripts/import_cards.go
```

**To Make Playable:**
1. Design ability system
2. Create test cards
3. Build card loader
4. Integrate with engine

**Timeline:**
- Database: 10 minutes
- Test cards: 1-2 weeks
- Ability system: 2-4 weeks
- Priority cards: 2-3 months
- Full port: 6-12 months
