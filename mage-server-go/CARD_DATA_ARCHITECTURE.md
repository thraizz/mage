# Card Data Architecture & Population Strategy

## Overview

This document explains how MAGE loads card data and outlines the strategy for populating the Go engine with cards, abilities, and decks.

## How the Java MAGE Server Works

### 1. **Card Data Source: Java Classes**

The Java MAGE server does **NOT** use external JSON/XML files for card data. Instead:

- **Every card is a Java class** in the `Mage.Sets` module (31,163+ Java files)
- Each card class extends `CardImpl` and defines its properties programmatically
- Cards are organized by expansion set (e.g., `mage.sets.innistrad.DelverOfSecrets`)
- Abilities are attached to cards via Java code in their constructors

Example card structure:
```java
public final class LightningBolt extends CardImpl {
    public LightningBolt(UUID ownerId, CardSetInfo setInfo) {
        super(ownerId, setInfo, new CardType[]{CardType.INSTANT}, "{R}");
        
        // Lightning Bolt deals 3 damage to any target.
        this.getSpellAbility().addEffect(new DamageTargetEffect(3));
        this.getSpellAbility().addTarget(new TargetAnyTarget());
    }
}
```

### 2. **Card Scanning & Database Population**

At server startup, the Java server:

1. **Scans all card classes** using reflection (`CardScanner.scan()`)
2. **Instantiates each card** to extract metadata (name, mana cost, types, etc.)
3. **Stores metadata in H2 database** (`CardRepository` → `cards.h2.mv.db`)
4. **Caches card info** for fast lookups during gameplay

The database schema (`CardInfo` table) stores:
- Basic properties: name, set code, card number, rarity
- Game properties: mana cost, power/toughness, types, subtypes
- Rules text (extracted from abilities)
- Color identity, frame style, art variants
- Special flags: split cards, double-faced, flip cards, etc.

### 3. **Card Instantiation During Games**

When a game needs a card:
1. **Lookup in CardRepository** by name/set/number
2. **Get the Java class name** from `CardInfo.className`
3. **Instantiate via reflection**: `CardImpl.createCard(className, cardSetInfo)`
4. The card's constructor **attaches all abilities** programmatically

### 4. **Ability System**

Abilities are **NOT stored in the database**. They are:
- Defined in card Java classes
- Instantiated when cards are created
- Attached to card objects at runtime
- Executed by the game engine

## Go Engine Population Strategy

### Option 1: **Port Card Classes to Go** (Recommended Long-Term)

**Approach:**
- Translate Java card classes to Go structs/functions
- Each card becomes a Go function that returns a configured `Card` struct
- Abilities are Go functions attached to cards

**Pros:**
- Full control over card behavior
- Type-safe, compiled code
- No external dependencies
- Can optimize for performance

**Cons:**
- Massive initial effort (31,000+ cards)
- Requires ongoing maintenance for new sets
- Need to port ability system

**Implementation:**
```go
// internal/cards/innistrad/delver_of_secrets.go
func NewDelverOfSecrets(id uuid.UUID, setInfo CardSetInfo) *Card {
    card := &Card{
        ID:       id,
        Name:     "Delver of Secrets",
        ManaCost: "{U}",
        Types:    []CardType{TypeCreature},
        Subtypes: []string{"Human", "Wizard"},
        Power:    "1",
        Toughness: "1",
    }
    
    // Add transform ability
    card.AddTriggeredAbility(NewDelverTransformAbility())
    
    return card
}
```

### Option 2: **Generate Go Code from Java** (Recommended Short-Term)

**Approach:**
- Write a Java-to-Go transpiler/generator
- Parse Java card classes and generate equivalent Go code
- Automate for new set releases

**Pros:**
- Automates the massive porting effort
- Can regenerate when Java cards change
- Maintains parity with Java implementation

**Cons:**
- Complex transpiler logic
- May not handle all edge cases
- Still requires ability system port

**Implementation:**
```bash
# Transpiler tool
./tools/card-transpiler \
  --input Mage.Sets/src/mage/sets/ \
  --output mage-server-go/internal/cards/ \
  --format go
```

### Option 3: **Use Card Metadata Database** (Recommended Immediate)

**Approach:**
- Populate PostgreSQL with card metadata from Java's H2 database
- Store basic card properties (name, cost, P/T, types)
- For abilities: use a scripting system or ability templates

**Pros:**
- Quick to implement
- Leverages existing database schema
- Easy to query and filter cards
- Good for deck validation and card lookups

**Cons:**
- Doesn't include ability logic
- Still need to implement game mechanics separately
- May need external ability definitions

**Implementation:**

#### Step 1: Export from Java H2 Database
```bash
# Export cards from Java server's H2 database
java -cp h2.jar org.h2.tools.Script \
  -url jdbc:h2:./db/cards \
  -script cards_export.sql
```

#### Step 2: Import to PostgreSQL
```sql
-- Transform and import to PostgreSQL
COPY cards (
  name, set_code, card_number, card_type, mana_cost,
  power, toughness, rules_text, rarity, card_class_name
)
FROM '/path/to/cards_export.csv'
DELIMITER ',' CSV HEADER;
```

#### Step 3: Use in Go Server
```go
// Fetch card metadata
card, err := cardRepo.GetByName(ctx, "Lightning Bolt")

// For gameplay, need to attach abilities separately
card.Abilities = LoadAbilitiesForCard(card.CardClassName)
```

### Option 4: **Hybrid Approach** (Recommended Overall)

**Phase 1: Metadata Database (Immediate)**
- Import card metadata to PostgreSQL
- Use for deck validation, card search, collection management
- Implement basic card properties (cost, types, P/T)

**Phase 2: Ability Templates (Short-Term)**
- Define common ability patterns as Go templates
- Store ability configurations in database
- Example: `DamageEffect{Amount: 3, Target: "any"}`

**Phase 3: Core Card Port (Medium-Term)**
- Port high-priority cards manually (Standard format, popular cards)
- Use ability template system for simple cards
- Focus on cards needed for testing

**Phase 4: Full Generation (Long-Term)**
- Build transpiler to generate Go from Java
- Automate for all cards
- Maintain parity with Java implementation

## Recommended Implementation Plan

### Phase 1: Database Population (Week 1-2)

1. **Export Java card data:**
   ```bash
   # Run from Java server
   java -jar mage-server.jar --export-cards cards_export.csv
   ```

2. **Create import script:**
   ```go
   // mage-server-go/scripts/import_cards.go
   func ImportCardsFromCSV(filepath string) error {
       // Read CSV
       // Transform to PostgreSQL format
       // Bulk insert into cards table
   }
   ```

3. **Populate database:**
   ```bash
   cd mage-server-go
   go run scripts/import_cards.go ../cards_export.csv
   ```

4. **Verify import:**
   ```bash
   PAGER=cat psql -d mage -c "SELECT COUNT(*) FROM cards;"
   # Should show ~31,000+ cards
   ```

### Phase 2: Ability System Design (Week 3-4)

1. **Define ability interfaces:**
   ```go
   // internal/game/ability.go
   type Ability interface {
       Execute(game *Game, source *Card, targets []Target) error
       CanExecute(game *Game, source *Card) bool
       GetCost() Cost
   }
   
   type TriggeredAbility interface {
       Ability
       GetTrigger() Trigger
   }
   
   type StaticAbility interface {
       Ability
       Apply(game *Game, source *Card)
   }
   ```

2. **Implement common abilities:**
   ```go
   // internal/game/abilities/damage.go
   type DamageEffect struct {
       Amount int
       Target TargetType
   }
   
   func (e *DamageEffect) Execute(game *Game, source *Card, targets []Target) error {
       for _, target := range targets {
           target.DealDamage(e.Amount, source)
       }
       return nil
   }
   ```

3. **Create ability registry:**
   ```go
   // internal/game/ability_registry.go
   var abilityRegistry = map[string]AbilityFactory{
       "damage_target": NewDamageTargetAbility,
       "draw_cards": NewDrawCardsAbility,
       "counter_spell": NewCounterSpellAbility,
   }
   ```

### Phase 3: Card Loading System (Week 5-6)

1. **Card loader with abilities:**
   ```go
   // internal/game/card_loader.go
   type CardLoader struct {
       cardRepo     *repository.CardRepository
       abilityReg   *AbilityRegistry
   }
   
   func (l *CardLoader) LoadCard(name string, setCode string) (*Card, error) {
       // Get metadata from database
       cardInfo, err := l.cardRepo.GetByName(ctx, name)
       if err != nil {
           return nil, err
       }
       
       // Create card instance
       card := &Card{
           Name:      cardInfo.Name,
           ManaCost:  cardInfo.ManaCost,
           Types:     parseTypes(cardInfo.CardType),
           Power:     cardInfo.Power,
           Toughness: cardInfo.Toughness,
       }
       
       // Attach abilities (from registry or database)
       abilities := l.loadAbilitiesForCard(cardInfo.CardClassName)
       card.Abilities = abilities
       
       return card, nil
   }
   ```

2. **Deck loading:**
   ```go
   // internal/game/deck_loader.go
   func (l *DeckLoader) LoadDeck(deckList DeckList) (*Deck, error) {
       deck := &Deck{Cards: make([]*Card, 0)}
       
       for _, entry := range deckList.Entries {
           card, err := l.cardLoader.LoadCard(entry.Name, entry.SetCode)
           if err != nil {
               return nil, err
           }
           
           for i := 0; i < entry.Count; i++ {
               deck.Cards = append(deck.Cards, card.Copy())
           }
       }
       
       return deck, nil
   }
   ```

### Phase 4: Priority Cards Port (Week 7-12)

Manually port high-priority cards for testing:

1. **Basic lands** (5 cards)
2. **Simple creatures** (50 cards)
3. **Simple instants/sorceries** (50 cards)
4. **Common keywords** (flying, trample, haste, etc.)
5. **Standard format staples** (100 cards)

This gives ~200 cards for comprehensive testing.

## Database Schema Enhancements

Add to existing `cards` table:

```sql
-- Add ability configuration column
ALTER TABLE cards ADD COLUMN abilities_json JSONB;

-- Store ability templates
CREATE TABLE card_abilities (
    id SERIAL PRIMARY KEY,
    card_id INTEGER REFERENCES cards(id),
    ability_type VARCHAR(100) NOT NULL,
    ability_config JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast ability lookup
CREATE INDEX idx_card_abilities_card_id ON card_abilities(card_id);
```

Example ability storage:
```json
{
  "type": "activated",
  "cost": "{T}",
  "effects": [
    {
      "type": "damage",
      "amount": 1,
      "target": "any"
    }
  ]
}
```

## Deck Management

### Deck Storage

```sql
-- Deck tables
CREATE TABLE decks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    format VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE deck_cards (
    id SERIAL PRIMARY KEY,
    deck_id INTEGER REFERENCES decks(id),
    card_id INTEGER REFERENCES cards(id),
    quantity INTEGER NOT NULL,
    is_sideboard BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_deck_cards_deck_id ON deck_cards(deck_id);
```

### Deck Validation

```go
// internal/game/deck_validator.go
type DeckValidator struct {
    cardRepo *repository.CardRepository
}

func (v *DeckValidator) ValidateDeck(deck *Deck, format string) error {
    // Check minimum deck size
    if len(deck.MainDeck) < 60 {
        return errors.New("deck must have at least 60 cards")
    }
    
    // Check card legality in format
    for _, card := range deck.AllCards() {
        if !v.isLegalInFormat(card, format) {
            return fmt.Errorf("card %s is not legal in %s", card.Name, format)
        }
    }
    
    // Check card limits (4-of rule)
    counts := make(map[string]int)
    for _, card := range deck.AllCards() {
        if !card.IsBasicLand() {
            counts[card.Name]++
            if counts[card.Name] > 4 {
                return fmt.Errorf("too many copies of %s", card.Name)
            }
        }
    }
    
    return nil
}
```

## Summary

**Immediate Actions:**
1. ✅ Database schema exists (`002_create_cards_table.up.sql`)
2. ✅ Card repository exists (`internal/repository/cards.go`)
3. 🔄 **Need to populate database** with card data from Java
4. 🔄 **Need to design ability system** for Go
5. 🔄 **Need to implement card loader** that attaches abilities

**Next Steps:**
1. Export card data from Java H2 database
2. Write import script for PostgreSQL
3. Design ability template system
4. Port 200 priority cards manually
5. Build card loader with ability attachment
6. Implement deck validation and loading

**Long-term:**
- Build Java-to-Go transpiler for automatic card generation
- Maintain parity with new set releases
- Optimize card loading and caching
