# Complete Card Implementation Plan

## Goal: Implement ALL 31,000+ Cards for MVP

This document outlines the complete strategy to implement every Magic card in the Go engine, using SQLite for card data storage.

---

## Phase 1: SQLite Card Database (Week 1)

### Step 1: Export & Convert H2 Database

```bash
cd mage-server-go

# 1. Export Java H2 to SQL
./scripts/h2_to_sql.sh

# 2. Import to SQLite
./scripts/import_to_sqlite.sh

# 3. Verify
sqlite3 data/cards.db "SELECT COUNT(*) FROM card;"
# Expected: 31,234 cards
```

**Deliverable:** `data/cards.db` with all card metadata

### Step 2: Create SQLite Card Repository

```go
// internal/repository/card_db.go
package repository

import (
    "database/sql"
    "fmt"
    "sync"
    "time"
    _ "github.com/mattn/go-sqlite3"
    "go.uber.org/zap"
)

type CardDB struct {
    db     *sql.DB
    cache  *cardCache
    logger *zap.Logger
}

func NewCardDB(dbPath string, logger *zap.Logger) (*CardDB, error) {
    // Open read-only with shared cache
    db, err := sql.Open("sqlite3", dbPath+"?mode=ro&cache=shared&_journal_mode=WAL")
    if err != nil {
        return nil, fmt.Errorf("failed to open card database: %w", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(time.Hour)
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping card database: %w", err)
    }
    
    return &CardDB{
        db:     db,
        cache:  newCardCache(10000),
        logger: logger,
    }, nil
}

func (db *CardDB) GetByName(name string) ([]*CardInfo, error) {
    // Check cache first
    if cards, ok := db.cache.get(name); ok {
        return cards, nil
    }
    
    query := `
        SELECT name, setcode, cardnumber, classname, manacosts, manavalue,
               power, toughness, startingloyalty, types, subtypes, supertypes,
               rules, rarity, black, blue, green, red, white
        FROM card
        WHERE name = ?
        ORDER BY setcode, cardnumber
    `
    
    rows, err := db.db.Query(query, name)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    cards := make([]*CardInfo, 0)
    for rows.Next() {
        card := &CardInfo{}
        err := rows.Scan(
            &card.Name, &card.SetCode, &card.CardNumber, &card.ClassName,
            &card.ManaCosts, &card.ManaValue, &card.Power, &card.Toughness,
            &card.StartingLoyalty, &card.Types, &card.Subtypes, &card.Supertypes,
            &card.Rules, &card.Rarity, &card.Black, &card.Blue, &card.Green,
            &card.Red, &card.White,
        )
        if err != nil {
            return nil, err
        }
        cards = append(cards, card)
    }
    
    // Cache results
    db.cache.set(name, cards)
    
    return cards, nil
}

func (db *CardDB) GetByClassName(className string) (*CardInfo, error) {
    query := `
        SELECT name, setcode, cardnumber, classname, manacosts, manavalue,
               power, toughness, startingloyalty, types, subtypes, supertypes,
               rules, rarity, black, blue, green, red, white
        FROM card
        WHERE classname = ?
        LIMIT 1
    `
    
    card := &CardInfo{}
    err := db.db.QueryRow(query, className).Scan(
        &card.Name, &card.SetCode, &card.CardNumber, &card.ClassName,
        &card.ManaCosts, &card.ManaValue, &card.Power, &card.Toughness,
        &card.StartingLoyalty, &card.Types, &card.Subtypes, &card.Supertypes,
        &card.Rules, &card.Rarity, &card.Black, &card.Blue, &card.Green,
        &card.Red, &card.White,
    )
    
    if err != nil {
        return nil, err
    }
    
    return card, nil
}

func (db *CardDB) SearchByName(searchTerm string, limit int) ([]*CardInfo, error) {
    query := `
        SELECT name, setcode, cardnumber, classname, manacosts, manavalue,
               power, toughness, startingloyalty, types, subtypes, supertypes,
               rules, rarity, black, blue, green, red, white
        FROM card
        WHERE name LIKE ?
        ORDER BY name
        LIMIT ?
    `
    
    rows, err := db.db.Query(query, "%"+searchTerm+"%", limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    cards := make([]*CardInfo, 0)
    for rows.Next() {
        card := &CardInfo{}
        err := rows.Scan(
            &card.Name, &card.SetCode, &card.CardNumber, &card.ClassName,
            &card.ManaCosts, &card.ManaValue, &card.Power, &card.Toughness,
            &card.StartingLoyalty, &card.Types, &card.Subtypes, &card.Supertypes,
            &card.Rules, &card.Rarity, &card.Black, &card.Blue, &card.Green,
            &card.Red, &card.White,
        )
        if err != nil {
            return nil, err
        }
        cards = append(cards, card)
    }
    
    return cards, nil
}

func (db *CardDB) PreloadAllCards() error {
    query := "SELECT classname FROM card"
    rows, err := db.db.Query(query)
    if err != nil {
        return err
    }
    defer rows.Close()
    
    count := 0
    for rows.Next() {
        count++
    }
    
    db.logger.Info("preloaded card database", zap.Int("cards", count))
    return nil
}

type CardInfo struct {
    Name            string
    SetCode         string
    CardNumber      string
    ClassName       string
    ManaCosts       string
    ManaValue       int
    Power           string
    Toughness       string
    StartingLoyalty string
    Types           string
    Subtypes        string
    Supertypes      string
    Rules           string
    Rarity          string
    Black           bool
    Blue            bool
    Green           bool
    Red             bool
    White           bool
}

type cardCache struct {
    items map[string][]*CardInfo
    mu    sync.RWMutex
}

func newCardCache(maxSize int) *cardCache {
    return &cardCache{
        items: make(map[string][]*CardInfo),
    }
}

func (c *cardCache) get(key string) ([]*CardInfo, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    cards, ok := c.items[key]
    return cards, ok
}

func (c *cardCache) set(key string, cards []*CardInfo) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = cards
}
```

**Deliverable:** Card database access layer

---

## Phase 2: Ability Framework (Week 2-3)

### Core Ability System

```go
// internal/game/ability/ability.go
package ability

type Ability interface {
    GetID() uuid.UUID
    GetType() AbilityType
    CanExecute(ctx context.Context, game *Game, source *Card) bool
    Execute(ctx context.Context, game *Game, source *Card, targets []Target) error
}

type AbilityType string

const (
    AbilityTypeActivated AbilityType = "activated"
    AbilityTypeTriggered AbilityType = "triggered"
    AbilityTypeStatic    AbilityType = "static"
    AbilityTypeSpell     AbilityType = "spell"
    AbilityTypeKeyword   AbilityType = "keyword"
)

// internal/game/ability/activated.go
type ActivatedAbility struct {
    ID          uuid.UUID
    Cost        Cost
    Effects     []Effect
    Targets     []TargetRequirement
    Timing      ActivationTiming
    UsesStack   bool
}

// internal/game/ability/triggered.go
type TriggeredAbility struct {
    ID          uuid.UUID
    Trigger     Trigger
    Condition   Condition
    Effects     []Effect
    Targets     []TargetRequirement
    Optional    bool
}

// internal/game/ability/static.go
type StaticAbility struct {
    ID          uuid.UUID
    Effect      ContinuousEffect
    Condition   Condition
    Layer       EffectLayer
}

// internal/game/ability/keyword.go
type KeywordAbility struct {
    ID      uuid.UUID
    Keyword Keyword
}

type Keyword string

const (
    KeywordFlying        Keyword = "flying"
    KeywordFirstStrike   Keyword = "first_strike"
    KeywordDoubleStrike  Keyword = "double_strike"
    KeywordTrample       Keyword = "trample"
    KeywordVigilance     Keyword = "vigilance"
    KeywordHaste         Keyword = "haste"
    KeywordDeathtouch    Keyword = "deathtouch"
    KeywordLifelink      Keyword = "lifelink"
    KeywordHexproof      Keyword = "hexproof"
    KeywordIndestructible Keyword = "indestructible"
    // ... 490+ more keywords
)
```

### Effect Library

```go
// internal/game/effect/effect.go
package effect

type Effect interface {
    Apply(ctx context.Context, game *Game, source *Card, targets []Target) error
    GetDescription() string
}

// internal/game/effect/damage.go
type DamageEffect struct {
    Amount       int
    SourceDamage bool
    Targets      []Target
}

// internal/game/effect/draw.go
type DrawCardsEffect struct {
    Amount int
}

// internal/game/effect/destroy.go
type DestroyEffect struct {
    Regenerate bool
}

// internal/game/effect/counter.go
type CounterSpellEffect struct{}

// internal/game/effect/life.go
type GainLifeEffect struct {
    Amount int
}

type LoseLifeEffect struct {
    Amount int
}

// internal/game/effect/discard.go
type DiscardEffect struct {
    Amount int
    Random bool
}

// internal/game/effect/sacrifice.go
type SacrificeEffect struct {
    Filter CardFilter
}

// internal/game/effect/exile.go
type ExileEffect struct {
    FromZone Zone
}

// internal/game/effect/return.go
type ReturnToHandEffect struct{}

type ReturnToBattlefieldEffect struct{}

// internal/game/effect/search.go
type SearchLibraryEffect struct {
    Filter   CardFilter
    Reveal   bool
    ToHand   bool
    ToField  bool
}

// internal/game/effect/shuffle.go
type ShuffleEffect struct{}

// internal/game/effect/mill.go
type MillEffect struct {
    Amount int
}

// internal/game/effect/token.go
type CreateTokenEffect struct {
    TokenName string
    Amount    int
}

// internal/game/effect/control.go
type GainControlEffect struct {
    Duration Duration
}

// internal/game/effect/boost.go
type BoostEffect struct {
    PowerBoost     int
    ToughnessBoost int
    Duration       Duration
}

// internal/game/effect/tap.go
type TapEffect struct{}

type UntapEffect struct{}

// internal/game/effect/counter_permanent.go
type AddCounterEffect struct {
    CounterType  CounterType
    Amount       int
}

type RemoveCounterEffect struct {
    CounterType  CounterType
    Amount       int
}

// ... 80+ more effect types
```

### Cost System

```go
// internal/game/cost/cost.go
package cost

type Cost interface {
    CanPay(game *Game, playerID uuid.UUID) bool
    Pay(game *Game, playerID uuid.UUID) error
    String() string
}

// internal/game/cost/mana.go
type ManaCost struct {
    Generic   int
    White     int
    Blue      int
    Black     int
    Red       int
    Green     int
    Colorless int
    Hybrid    []HybridMana
    Phyrexian []PhyrexianMana
}

// internal/game/cost/tap.go
type TapCost struct{}

// internal/game/cost/sacrifice.go
type SacrificeTargetCost struct {
    Filter CardFilter
    Amount int
}

// internal/game/cost/discard.go
type DiscardCost struct {
    Amount int
    Random bool
}

// internal/game/cost/life.go
type PayLifeCost struct {
    Amount int
}

// internal/game/cost/exile.go
type ExileFromGraveyardCost struct {
    Amount int
}

// internal/game/cost/composite.go
type CompositeCost struct {
    Costs []Cost
}
```

### Target System

```go
// internal/game/target/target.go
package target

type Target interface {
    IsValid(game *Game) bool
    GetID() uuid.UUID
}

type TargetRequirement struct {
    Type     TargetType
    MinCount int
    MaxCount int
    Filter   CardFilter
}

type TargetType string

const (
    TargetTypeAny        TargetType = "any"
    TargetTypeCreature   TargetType = "creature"
    TargetTypePlayer     TargetType = "player"
    TargetTypePermanent  TargetType = "permanent"
    TargetTypeSpell      TargetType = "spell"
    TargetTypePlaneswalker TargetType = "planeswalker"
    TargetTypeArtifact   TargetType = "artifact"
    TargetTypeEnchantment TargetType = "enchantment"
    TargetTypeLand       TargetType = "land"
)

type CardTarget struct {
    CardID uuid.UUID
}

type PlayerTarget struct {
    PlayerID uuid.UUID
}

type SpellTarget struct {
    SpellID uuid.UUID
}
```

**Deliverable:** Complete ability framework

---

## Phase 3: Card Generator/Transpiler (Week 4-6)

### Transpiler Architecture

```go
// tools/transpiler/main.go
package main

type Transpiler struct {
    parser      *JavaParser
    mapper      *AbilityMapper
    generator   *GoGenerator
    validator   *Validator
    logger      *zap.Logger
}

func (t *Transpiler) TranspileCard(javaFile string) (*GoCard, error) {
    // 1. Parse Java file
    ast, err := t.parser.Parse(javaFile)
    if err != nil {
        return nil, err
    }
    
    // 2. Extract card info
    cardInfo := extractCardInfo(ast)
    
    // 3. Map abilities to Go
    abilities := t.mapper.MapAbilities(ast.Abilities)
    
    // 4. Generate Go code
    goCode := t.generator.Generate(cardInfo, abilities)
    
    // 5. Validate
    if err := t.validator.Validate(goCode); err != nil {
        return nil, err
    }
    
    return goCode, nil
}

// tools/transpiler/parser.go
type JavaParser struct {
    // Parse Java source files
}

func (p *JavaParser) Parse(filepath string) (*JavaAST, error) {
    // Read Java file
    content, err := os.ReadFile(filepath)
    if err != nil {
        return nil, err
    }
    
    // Extract card properties
    ast := &JavaAST{
        ClassName:  extractClassName(content),
        SuperClass: extractSuperClass(content),
        ManaCost:   extractManaCost(content),
        Types:      extractTypes(content),
        Subtypes:   extractSubtypes(content),
        Power:      extractPower(content),
        Toughness:  extractToughness(content),
        Abilities:  extractAbilities(content),
    }
    
    return ast, nil
}

// tools/transpiler/mapper.go
type AbilityMapper struct {
    // Map Java abilities to Go equivalents
}

var abilityMap = map[string]string{
    // Java class -> Go function
    "DamageTargetEffect":              "NewDamageTargetEffect",
    "DrawCardSourceControllerEffect":  "NewDrawCardsEffect",
    "GainLifeEffect":                  "NewGainLifeEffect",
    "CounterTargetEffect":             "NewCounterSpellEffect",
    "DestroyTargetEffect":             "NewDestroyEffect",
    "ExileTargetEffect":               "NewExileEffect",
    "ReturnToHandTargetEffect":        "NewReturnToHandEffect",
    "PutCounterTargetEffect":          "NewAddCounterEffect",
    "TapTargetEffect":                 "NewTapEffect",
    "UntapTargetEffect":               "NewUntapEffect",
    "DiscardTargetEffect":             "NewDiscardEffect",
    "SacrificeTargetEffect":           "NewSacrificeEffect",
    "SearchLibraryPutInHandEffect":    "NewSearchLibraryEffect",
    "ShuffleLibrarySourceEffect":      "NewShuffleEffect",
    "PutLibraryIntoGraveTargetEffect": "NewMillEffect",
    "CreateTokenEffect":               "NewCreateTokenEffect",
    "GainControlTargetEffect":         "NewGainControlEffect",
    "BoostTargetEffect":               "NewBoostEffect",
    // ... 1000+ mappings
}

// tools/transpiler/generator.go
type GoGenerator struct {
    // Generate Go code from AST
}

func (g *GoGenerator) Generate(cardInfo *CardInfo, abilities []*Ability) string {
    tmpl := `package generated

import (
    "github.com/google/uuid"
    "github.com/magefree/mage-server-go/internal/game"
    "github.com/magefree/mage-server-go/internal/game/ability"
    "github.com/magefree/mage-server-go/internal/game/effect"
)

func New{{.Name}}(id uuid.UUID, setInfo game.CardSetInfo) *game.Card {
    card := game.NewCard(id, setInfo, []game.CardType{ {{.Types}} }, "{{.ManaCost}}")
    {{if .Power}}card.Power = {{.Power}}{{end}}
    {{if .Toughness}}card.Toughness = {{.Toughness}}{{end}}
    
    {{range .Abilities}}
    {{.}}
    {{end}}
    
    return card
}
`
    
    // Execute template
    return executeTemplate(tmpl, cardInfo, abilities)
}
```

### Batch Generation

```bash
# tools/transpiler/generate_all.sh
#!/bin/bash

JAVA_SETS_DIR="../Mage.Sets/src/mage/sets"
OUTPUT_DIR="internal/cards/generated"

echo "Generating all cards..."

# Find all Java card files
find "$JAVA_SETS_DIR" -name "*.java" | while read -r java_file; do
    # Extract set and card name
    set=$(basename $(dirname "$java_file"))
    card=$(basename "$java_file" .java)
    
    # Create output directory
    mkdir -p "$OUTPUT_DIR/$set"
    
    # Generate Go file
    ./tools/transpiler/transpiler \
        --input "$java_file" \
        --output "$OUTPUT_DIR/$set/${card}.go" \
        --package "generated"
    
    echo "Generated: $set/$card"
done

echo "✓ Generated all cards"
```

**Deliverable:** Transpiler that converts Java cards to Go

---

## Phase 4: Generate All Cards (Week 7-8)

### Generate All 31,000 Cards

```bash
cd mage-server-go

# Generate all cards
./tools/transpiler/generate_all.sh

# Result:
# internal/cards/generated/
#   ├── innistrad/
#   │   ├── delver_of_secrets.go
#   │   ├── snapcaster_mage.go
#   │   └── ... (hundreds more)
#   ├── ravnica/
#   │   └── ...
#   └── ... (100+ sets, 31,000+ files)
```

### Card Registry

```go
// internal/cards/registry.go
package cards

import (
    "fmt"
    "github.com/google/uuid"
    "github.com/magefree/mage-server-go/internal/game"
)

type CardFactory func(id uuid.UUID, setInfo game.CardSetInfo) *game.Card

var Registry = make(map[string]CardFactory)

func RegisterCard(className string, factory CardFactory) {
    Registry[className] = factory
}

func CreateCard(className string, id uuid.UUID, setInfo game.CardSetInfo) (*game.Card, error) {
    factory, ok := Registry[className]
    if !ok {
        return nil, fmt.Errorf("card not found: %s", className)
    }
    return factory(id, setInfo), nil
}

// Auto-register all generated cards
func init() {
    // Generated cards register themselves
    // via init() functions in each file
}
```

### Auto-Registration

```go
// internal/cards/generated/innistrad/delver_of_secrets.go
package innistrad

import (
    "github.com/google/uuid"
    "github.com/magefree/mage-server-go/internal/cards"
    "github.com/magefree/mage-server-go/internal/game"
)

func init() {
    cards.RegisterCard("mage.cards.d.DelverOfSecrets", NewDelverOfSecrets)
}

func NewDelverOfSecrets(id uuid.UUID, setInfo game.CardSetInfo) *game.Card {
    // Generated card implementation
    card := game.NewCard(id, setInfo, []game.CardType{game.TypeCreature}, "{U}")
    card.Subtypes = []string{"Human", "Wizard"}
    card.Power = 1
    card.Toughness = 1
    
    // Add triggered ability
    // ...
    
    return card
}
```

**Deliverable:** All 31,000+ cards in Go

---

## Phase 5: Card Loader (Week 9)

### Load Cards from Database + Registry

```go
// internal/game/card_loader.go
package game

import (
    "context"
    "fmt"
    "github.com/google/uuid"
    "github.com/magefree/mage-server-go/internal/cards"
    "github.com/magefree/mage-server-go/internal/repository"
)

type CardLoader struct {
    cardDB   *repository.CardDB
    registry map[string]cards.CardFactory
    logger   *zap.Logger
}

func NewCardLoader(cardDB *repository.CardDB, logger *zap.Logger) *CardLoader {
    return &CardLoader{
        cardDB:   cardDB,
        registry: cards.Registry,
        logger:   logger,
    }
}

func (l *CardLoader) LoadCard(name string, setCode string) (*Card, error) {
    // 1. Get card metadata from SQLite
    cardInfos, err := l.cardDB.GetByName(name)
    if err != nil {
        return nil, fmt.Errorf("failed to get card info: %w", err)
    }
    
    if len(cardInfos) == 0 {
        return nil, fmt.Errorf("card not found: %s", name)
    }
    
    // 2. Find specific printing
    var cardInfo *repository.CardInfo
    if setCode != "" {
        for _, info := range cardInfos {
            if info.SetCode == setCode {
                cardInfo = info
                break
            }
        }
    } else {
        // Use first printing
        cardInfo = cardInfos[0]
    }
    
    if cardInfo == nil {
        return nil, fmt.Errorf("card not found: %s [%s]", name, setCode)
    }
    
    // 3. Create card from registry
    factory, ok := l.registry[cardInfo.ClassName]
    if !ok {
        return nil, fmt.Errorf("card class not registered: %s", cardInfo.ClassName)
    }
    
    setInfo := CardSetInfo{
        Name:       cardInfo.Name,
        SetCode:    cardInfo.SetCode,
        CardNumber: cardInfo.CardNumber,
        Rarity:     cardInfo.Rarity,
    }
    
    card := factory(uuid.New(), setInfo)
    
    return card, nil
}

func (l *CardLoader) LoadDeck(deckList *DeckList) (*Deck, error) {
    deck := &Deck{
        MainDeck: make([]*Card, 0),
        Sideboard: make([]*Card, 0),
    }
    
    // Load main deck
    for _, entry := range deckList.MainDeck {
        for i := 0; i < entry.Quantity; i++ {
            card, err := l.LoadCard(entry.CardName, entry.SetCode)
            if err != nil {
                return nil, fmt.Errorf("failed to load card %s: %w", entry.CardName, err)
            }
            deck.MainDeck = append(deck.MainDeck, card)
        }
    }
    
    // Load sideboard
    for _, entry := range deckList.Sideboard {
        for i := 0; i < entry.Quantity; i++ {
            card, err := l.LoadCard(entry.CardName, entry.SetCode)
            if err != nil {
                return nil, fmt.Errorf("failed to load card %s: %w", entry.CardName, err)
            }
            deck.Sideboard = append(deck.Sideboard, card)
        }
    }
    
    return deck, nil
}
```

**Deliverable:** System to load any card from database + registry

---

## Phase 6: Testing & Fixes (Week 10-12)

### Automated Testing

```go
// internal/cards/generated/test_all_cards.go
package generated_test

import (
    "testing"
    "github.com/google/uuid"
    "github.com/magefree/mage-server-go/internal/cards"
    "github.com/magefree/mage-server-go/internal/game"
)

func TestAllCardsInstantiate(t *testing.T) {
    // Test that every card in registry can be instantiated
    for className, factory := range cards.Registry {
        t.Run(className, func(t *testing.T) {
            setInfo := game.CardSetInfo{
                Name:    "Test Card",
                SetCode: "TST",
            }
            
            card := factory(uuid.New(), setInfo)
            
            if card == nil {
                t.Errorf("card is nil: %s", className)
            }
            
            if card.Name == "" {
                t.Errorf("card has no name: %s", className)
            }
        })
    }
}

func TestCardAbilities(t *testing.T) {
    // Test that cards have expected abilities
    tests := []struct {
        className     string
        expectedAbilities int
    }{
        {"mage.cards.l.LightningBolt", 1},
        {"mage.cards.g.GrizzlyBears", 0},
        {"mage.cards.l.LlanowarElves", 1},
    }
    
    for _, tt := range tests {
        t.Run(tt.className, func(t *testing.T) {
            factory := cards.Registry[tt.className]
            card := factory(uuid.New(), game.CardSetInfo{})
            
            if len(card.Abilities) != tt.expectedAbilities {
                t.Errorf("expected %d abilities, got %d",
                    tt.expectedAbilities, len(card.Abilities))
            }
        })
    }
}
```

### Manual Fixes

```bash
# Identify cards that need manual fixes
go test ./internal/cards/generated/... -v | grep FAIL > failed_cards.txt

# Categories of failures:
# 1. Complex abilities (5-10% of cards)
# 2. Custom mechanics (2-3% of cards)
# 3. Transpiler limitations (1-2% of cards)

# Estimated: 1,500-3,000 cards need manual attention
```

**Deliverable:** All cards working, tests passing

---

## Timeline Summary

| Phase | Duration | Deliverable |
|-------|----------|-------------|
| **Phase 1: SQLite DB** | Week 1 | Card database with 31,000+ cards |
| **Phase 2: Ability Framework** | Week 2-3 | Complete ability system |
| **Phase 3: Transpiler** | Week 4-6 | Java-to-Go card generator |
| **Phase 4: Generate Cards** | Week 7-8 | All 31,000+ cards in Go |
| **Phase 5: Card Loader** | Week 9 | Load cards from DB + registry |
| **Phase 6: Testing & Fixes** | Week 10-12 | All cards working |
| **Total** | **12 weeks** | **All cards playable** |

---

## File Structure

```
mage-server-go/
├── data/
│   ├── cards.db                      # SQLite database (50MB)
│   └── cards.db.sha256               # Checksum
├── internal/
│   ├── repository/
│   │   └── card_db.go                # SQLite card repository
│   ├── game/
│   │   ├── card.go                   # Card struct
│   │   ├── card_loader.go            # Load cards from DB
│   │   ├── ability/
│   │   │   ├── ability.go            # Base interfaces
│   │   │   ├── activated.go          # Activated abilities
│   │   │   ├── triggered.go          # Triggered abilities
│   │   │   ├── static.go             # Static abilities
│   │   │   └── keyword.go            # Keyword abilities
│   │   ├── effect/
│   │   │   ├── effect.go             # Base interface
│   │   │   ├── damage.go             # Damage effects
│   │   │   ├── draw.go               # Draw effects
│   │   │   └── ... (80+ effect types)
│   │   ├── cost/
│   │   │   ├── cost.go               # Base interface
│   │   │   ├── mana.go               # Mana costs
│   │   │   └── ... (10+ cost types)
│   │   └── target/
│   │       ├── target.go             # Base interface
│   │       └── ... (5+ target types)
│   └── cards/
│       ├── registry.go               # Card registry
│       └── generated/
│           ├── innistrad/
│           │   ├── delver_of_secrets.go
│           │   └── ... (hundreds more)
│           └── ... (100+ sets)
├── tools/
│   └── transpiler/
│       ├── main.go                   # Transpiler entry point
│       ├── parser.go                 # Java parser
│       ├── mapper.go                 # Ability mapper
│       ├── generator.go              # Go code generator
│       └── generate_all.sh           # Batch generation script
└── scripts/
    ├── h2_to_sql.sh                  # Export H2 to SQL
    ├── import_to_sqlite.sh           # Import to SQLite
    └── convert_h2_to_sqlite.py       # SQL converter
```

---

## Quick Start Commands

```bash
# Week 1: Set up database
./scripts/h2_to_sql.sh
./scripts/import_to_sqlite.sh
sqlite3 data/cards.db "SELECT COUNT(*) FROM card;"

# Week 2-3: Implement ability framework
# (Manual coding - see Phase 2)

# Week 4-6: Build transpiler
cd tools/transpiler
go build -o transpiler main.go

# Week 7-8: Generate all cards
./tools/transpiler/generate_all.sh

# Week 9: Test card loading
go run cmd/test_card_loader/main.go

# Week 10-12: Run tests and fix issues
go test ./internal/cards/generated/... -v
```

---

## Success Metrics

- ✅ **31,000+ cards** in SQLite database
- ✅ **31,000+ Go files** generated
- ✅ **All cards** instantiate without errors
- ✅ **95%+ cards** work correctly
- ✅ **5% cards** fixed manually
- ✅ **All tests** passing
- ✅ **Any deck** can be loaded and played

---

## Maintenance

### Adding New Sets

```bash
# 1. Update Java MAGE
cd ../Mage
git pull
mvn clean install

# 2. Regenerate H2 database
cd ../Mage.Server
mvn exec:java
# Ctrl+C after startup

# 3. Re-export to SQLite
cd ../mage-server-go
./scripts/h2_to_sql.sh
./scripts/import_to_sqlite.sh

# 4. Regenerate new cards
./tools/transpiler/generate_all.sh --new-only

# 5. Test new cards
go test ./internal/cards/generated/... -v
```

### Updating Abilities

```bash
# If ability system changes:
# 1. Update ability framework
# 2. Update transpiler mappings
# 3. Regenerate affected cards
./tools/transpiler/generate_all.sh --force
```

---

## End Goal

**You will have:**
- ✅ SQLite database with all card metadata
- ✅ Complete ability framework in Go
- ✅ Transpiler to convert Java cards to Go
- ✅ All 31,000+ cards as Go code
- ✅ Card loader to instantiate any card
- ✅ Ability to play any Magic deck
- ✅ Full parity with Java MAGE

**Timeline:** 12 weeks (3 months)

**Effort:** Full-time work or 6 months part-time

**Result:** Complete card implementation for MVP
