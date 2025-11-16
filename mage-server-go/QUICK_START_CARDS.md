# Quick Start: Populating Cards in Go Engine

## TL;DR - Do I Need to Populate Cards?

**Yes, eventually.** But the approach depends on what you're building:

### For Basic Server Functionality (Now)
- ✅ **Database schema exists** - cards table is ready
- ✅ **Repository layer exists** - `internal/repository/cards.go`
- ⚠️ **Database is empty** - no card data yet
- ⚠️ **No ability system** - cards would be metadata only

### For Game Engine (Later)
- ❌ **Need card data** - 30,000+ cards from Java
- ❌ **Need ability system** - how cards actually work
- ❌ **Need card loader** - instantiate cards with abilities
- ❌ **Need deck validation** - check legality, limits

## Current State

```
┌─────────────────────────────────────────────────────────┐
│ Java MAGE Server                                        │
│                                                         │
│  ┌──────────────┐      ┌─────────────────┐            │
│  │ 31,000+ Java │ ───▶ │ H2 Database     │            │
│  │ Card Classes │      │ (cards.h2.mv.db)│            │
│  └──────────────┘      └─────────────────┘            │
│         │                      │                        │
│         │ Reflection           │ Metadata               │
│         ▼                      ▼                        │
│  ┌──────────────────────────────────┐                  │
│  │ CardRepository (in-memory cache) │                  │
│  └──────────────────────────────────┘                  │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│ Go MAGE Server (Current)                                │
│                                                         │
│  ┌──────────────┐      ┌─────────────────┐            │
│  │ Empty        │      │ PostgreSQL      │            │
│  │              │      │ (cards table)   │            │
│  │              │      │ [EMPTY]         │            │
│  └──────────────┘      └─────────────────┘            │
│                               │                         │
│                               ▼                         │
│  ┌──────────────────────────────────┐                  │
│  │ CardRepository (with cache)      │                  │
│  │ - GetByID()     ✅               │                  │
│  │ - GetByName()   ✅               │                  │
│  │ - SearchByName() ✅              │                  │
│  └──────────────────────────────────┘                  │
└─────────────────────────────────────────────────────────┘
```

## Step-by-Step: Populate Card Database

### Step 1: Export from Java (5 minutes)

```bash
cd mage-server-go

# Export card metadata from Java's H2 database
./scripts/export_java_cards.sh
```

This will:
- Find the Java H2 database at `../Mage.Server/db/cards.h2.mv.db`
- Export all cards to `data/cards_export.csv`
- Download H2 JAR if needed

**Expected output:**
```
=== MAGE Card Data Export ===
Java DB path: /path/to/Mage.Server/db/cards.h2.mv.db
Output CSV: /path/to/mage-server-go/data/cards_export.csv
Exporting cards from H2 database...
✓ Export complete: data/cards_export.csv
Total cards exported: 31,234
```

### Step 2: Import to PostgreSQL (2 minutes)

```bash
# Make sure PostgreSQL is running
# Make sure database 'mage' exists

# Import cards
go run scripts/import_cards.go data/cards_export.csv
```

**Expected output:**
```
=== MAGE Card Data Import ===
CSV file: /path/to/data/cards_export.csv
Connecting to database...
✓ Database connection established
Found 31234 cards in CSV
Parsed 31234 valid cards
Importing cards...
Progress: 5000/31234 cards imported
Progress: 10000/31234 cards imported
...
Progress: 31234/31234 cards imported

=== Import Complete ===
✓ Successfully imported: 31234 cards
Time taken: 12.3s
Rate: 2540 cards/second

Total cards in database: 31234
```

### Step 3: Verify Import (30 seconds)

```bash
# Check total count
PAGER=cat psql -d mage -c "SELECT COUNT(*) FROM cards;"

# View sample cards
PAGER=cat psql -d mage -c "
  SELECT name, set_code, mana_cost, rarity 
  FROM cards 
  WHERE name LIKE 'Lightning%' 
  LIMIT 5;
"

# Search by name
PAGER=cat psql -d mage -c "
  SELECT name, set_code, card_type 
  FROM cards 
  WHERE name ILIKE '%delver%';
"
```

**Expected output:**
```
 count 
-------
 31234

       name       | set_code | mana_cost | rarity 
------------------+----------+-----------+--------
 Lightning Bolt   | LEA      | {R}       | COMMON
 Lightning Bolt   | LEB      | {R}       | COMMON
 Lightning Bolt   | 2ED      | {R}       | COMMON
```

## What You Get After Import

### ✅ Card Metadata Available

```go
// In your Go code, you can now:
card, err := cardRepo.GetByName(ctx, "Lightning Bolt")
// Returns: Card with name, mana cost, types, P/T, rules text, etc.

cards, err := cardRepo.SearchByName(ctx, "delver", 10)
// Returns: All cards matching "delver"

cards, err := cardRepo.GetBySetCode(ctx, "ISD")
// Returns: All cards from Innistrad set
```

### ⚠️ What's Still Missing

1. **Abilities are not functional** - cards have rules text but no executable code
2. **No card instantiation** - can't create card objects for games
3. **No deck loading** - can validate deck lists but can't load into games
4. **No game integration** - engine can't use cards yet

## Next Steps: Making Cards Playable

### Option A: Start with Test Cards (Recommended)

Create a small set of hardcoded cards for testing:

```go
// internal/game/cards/test_cards.go
package cards

func NewLightningBolt(id uuid.UUID) *Card {
    card := &Card{
        ID:       id,
        Name:     "Lightning Bolt",
        ManaCost: "{R}",
        Types:    []CardType{TypeInstant},
    }
    
    // Add damage ability
    ability := &ActivatedAbility{
        Cost: ManaCost{Red: 1},
        Effect: &DamageEffect{
            Amount: 3,
            Target: TargetAnyTarget,
        },
    }
    card.AddAbility(ability)
    
    return card
}
```

**Pros:**
- Quick to implement
- Full control
- Easy to test

**Cons:**
- Only works for a few cards
- Not scalable

### Option B: Ability Template System (Medium-term)

Store ability configurations in database:

```sql
-- Add ability storage
ALTER TABLE cards ADD COLUMN abilities_json JSONB;

-- Example: Lightning Bolt
UPDATE cards 
SET abilities_json = '[
  {
    "type": "spell",
    "effects": [
      {
        "type": "damage",
        "amount": 3,
        "target": "any"
      }
    ]
  }
]'
WHERE name = 'Lightning Bolt';
```

Then load abilities dynamically:

```go
func (l *CardLoader) LoadCard(name string) (*Card, error) {
    // Get from database
    cardInfo, _ := l.repo.GetByName(ctx, name)
    
    // Create card
    card := &Card{Name: cardInfo.Name, ...}
    
    // Parse and attach abilities
    abilities := l.parseAbilities(cardInfo.AbilitiesJSON)
    card.Abilities = abilities
    
    return card, nil
}
```

**Pros:**
- Scalable to many cards
- No code changes for new cards
- Can update abilities without recompiling

**Cons:**
- Limited to template-able abilities
- Complex abilities need custom code

### Option C: Port Card Classes (Long-term)

Translate Java card classes to Go:

```go
// internal/cards/innistrad/delver_of_secrets.go
func init() {
    RegisterCard("mage.cards.d.DelverOfSecrets", NewDelverOfSecrets)
}

func NewDelverOfSecrets(id uuid.UUID, setInfo CardSetInfo) *Card {
    card := &Card{
        ID:        id,
        Name:      "Delver of Secrets",
        ManaCost:  "{U}",
        Types:     []CardType{TypeCreature},
        Subtypes:  []string{"Human", "Wizard"},
        Power:     "1",
        Toughness: "1",
    }
    
    // At the beginning of your upkeep, look at the top card of your library.
    // You may reveal that card. If an instant or sorcery card is revealed 
    // this way, transform Delver of Secrets.
    trigger := &TriggeredAbility{
        Trigger: TriggerBeginningOfUpkeep,
        Effect:  &DelverTransformEffect{},
    }
    card.AddAbility(trigger)
    
    return card
}
```

**Pros:**
- Full parity with Java
- All cards eventually supported
- Type-safe, compiled

**Cons:**
- Massive effort (31,000+ cards)
- Requires ability system port
- Ongoing maintenance

## Recommended Approach

### Phase 1: Database Population (This Week)
✅ Run export script  
✅ Run import script  
✅ Verify card data  

### Phase 2: Test Card Set (Next Week)
- Manually create 10-20 simple cards
- Implement basic abilities (damage, draw, counter)
- Test in game engine

### Phase 3: Ability Templates (Month 1)
- Design ability template system
- Add ability JSON to database
- Implement dynamic card loader

### Phase 4: Priority Cards (Month 2-3)
- Port 200 most-used cards
- Focus on Standard format
- Build automated tests

### Phase 5: Full Port (Month 4+)
- Build Java-to-Go transpiler
- Generate all cards
- Maintain parity

## FAQ

**Q: Can I use the Go server without populating cards?**  
A: Yes, for user management, chat, lobbies. But you can't play games.

**Q: Do I need all 31,000 cards?**  
A: No. Start with 20-50 cards for testing. Add more as needed.

**Q: Can I mix hardcoded and database cards?**  
A: Yes! Use hardcoded for complex cards, templates for simple ones.

**Q: How does the Java server handle new sets?**  
A: Developers write Java classes for new cards. Server scans and imports on startup.

**Q: Can I auto-generate Go cards from Java?**  
A: Yes, but it's complex. Better to start manual, then build tooling.

**Q: What about tokens and emblems?**  
A: They're in the database too. Import them the same way.

**Q: How do I handle card errata?**  
A: Update the database. For abilities, update the card class or template.

## Summary

**Right now:**
1. ✅ Database schema ready
2. ✅ Repository layer ready
3. ⚠️ Database empty (run scripts to populate)
4. ❌ Ability system needed
5. ❌ Card loader needed

**To populate:**
```bash
# 1. Export from Java
./scripts/export_java_cards.sh

# 2. Import to PostgreSQL  
go run scripts/import_cards.go

# 3. Verify
PAGER=cat psql -d mage -c "SELECT COUNT(*) FROM cards;"
```

**To make playable:**
- Design ability system (see CARD_DATA_ARCHITECTURE.md)
- Create test cards manually
- Build card loader with ability attachment
- Integrate with game engine

**Timeline:**
- Database population: 10 minutes
- Test cards: 1 week
- Ability system: 2-4 weeks
- Priority cards: 2-3 months
- Full port: 6-12 months
