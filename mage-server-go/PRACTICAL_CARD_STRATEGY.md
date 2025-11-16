# Practical Card Strategy: The Fastest Path to MVP

## TL;DR - The Pragmatic Approach

**Don't convert databases. Don't build transpilers yet. Start simple.**

### What You Actually Need for MVP

1. **Card metadata** - Already in Java's H2 database ✅
2. **Ability system** - Build in Go (2-3 weeks)
3. **~500 Standard cards** - Manually port (4-6 weeks)
4. **Transpiler** - Build after MVP works (optional)

## The Simplest Path: Use Java's H2 Database Directly

### Option 1: Read H2 from Go (Recommended for Quick Start)

**Why:** Java already has all cards. Just read them.

```go
// Use H2 JDBC driver via CGo or pure Go H2 reader
import "github.com/jmrobles/h2go"

db, err := h2go.Open("../Mage.Server/db/cards")
rows, err := db.Query("SELECT * FROM card WHERE name = ?", "Lightning Bolt")
```

**Pros:**
- ✅ No conversion needed
- ✅ Always up-to-date with Java
- ✅ Start coding today

**Cons:**
- ⚠️ Dependency on H2 format
- ⚠️ Slightly slower than SQLite

### Option 2: One-Time Export to SQLite

**Why:** Portable, fast, no dependencies.

```bash
# Run once
./scripts/export_h2_to_sqlite.sh

# Result: data/cards.db (50MB)
# Ship with your server
```

**Pros:**
- ✅ Portable single file
- ✅ Fast reads
- ✅ No Java dependency

**Cons:**
- ⚠️ Need to re-export for updates
- ⚠️ One-time conversion effort

### Option 3: PostgreSQL (Current Plan)

**Why:** You already have PostgreSQL for game data.

**Pros:**
- ✅ One database for everything
- ✅ Familiar tooling

**Cons:**
- ⚠️ Overkill for read-only data
- ⚠️ Need PostgreSQL running

## My Recommendation: Hybrid Approach

```
┌─────────────────────────────────────────────────────┐
│ Development (Use Java's H2)                         │
│                                                     │
│  Go Server → Read H2 directly → Fast iteration     │
│              (no conversion)                        │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ Production (Use SQLite)                             │
│                                                     │
│  Go Server → Read SQLite → Portable, fast          │
│              (ship cards.db with binary)            │
└─────────────────────────────────────────────────────┘
```

## Ability Implementation: The Real Work

**This is where you spend your time, not database conversion.**

### Week 1-2: Core Framework

```go
// internal/game/ability/
├── ability.go          # Base interfaces
├── activated.go        # Activated abilities
├── triggered.go        # Triggered abilities
├── static.go           # Static abilities
├── keyword.go          # Keyword abilities
└── spell.go            # Spell abilities

// internal/game/effect/
├── effect.go           # Base interface
├── damage.go           # Deal damage
├── draw.go             # Draw cards
├── destroy.go          # Destroy permanents
├── counter.go          # Counter spells
├── life.go             # Gain/lose life
└── ... (20 more common effects)

// internal/game/cost/
├── cost.go             # Base interface
├── mana.go             # Mana costs
├── tap.go              # Tap costs
├── sacrifice.go        # Sacrifice costs
└── ... (10 more cost types)

// internal/game/target/
├── target.go           # Base interface
├── card.go             # Card targets
├── player.go           # Player targets
├── spell.go            # Spell targets
└── ... (5 more target types)
```

**Deliverable:** Framework that can express most abilities

### Week 3-4: First 20 Cards (Manual)

**Start with the simplest cards to validate your framework:**

```go
// internal/cards/manual/basic_lands.go
func NewPlains() *Card { ... }
func NewIsland() *Card { ... }
func NewSwamp() *Card { ... }
func NewMountain() *Card { ... }
func NewForest() *Card { ... }

// internal/cards/manual/simple_creatures.go
func NewGrizzlyBears() *Card {
    // 2/2 for {1}{G}, no abilities
}

func NewEliteVanguard() *Card {
    // 2/1 for {W}, no abilities
}

// internal/cards/manual/simple_spells.go
func NewLightningBolt() *Card {
    card := NewInstant("{R}")
    card.AddEffect(NewDamageEffect(3, TargetAny))
    return card
}

func NewDivination() *Card {
    card := NewSorcery("{2}{U}")
    card.AddEffect(NewDrawCardsEffect(2))
    return card
}

func NewMurder() *Card {
    card := NewInstant("{1}{B}{B}")
    card.AddEffect(NewDestroyEffect(TargetCreature))
    return card
}
```

**Test these 20 cards thoroughly:**
- Cast spells
- Resolve abilities
- Target selection
- Mana payment
- Stack interaction

**Deliverable:** 20 working cards, validated framework

### Week 5-8: Expand to 100 Cards

**Add more complex cards:**

```go
// Creatures with keywords
func NewShivanDragon() *Card {
    card := NewCreature("{4}{R}{R}", 5, 5)
    card.AddKeyword(KeywordFlying)
    card.AddActivatedAbility(
        ManaCost("{R}"),
        NewBoostEffect(1, 0, UntilEndOfTurn),
    )
    return card
}

// Triggered abilities
func NewLlanowarElves() *Card {
    card := NewCreature("{G}", 1, 1)
    card.AddActivatedAbility(
        TapCost(),
        NewAddManaEffect(Mana{Green: 1}),
    )
    return card
}

// Modal spells
func NewCharmSpell() *Card {
    card := NewInstant("{2}{G}")
    card.AddModalChoice(
        NewDestroyEffect(TargetArtifact),
        NewBoostEffect(3, 3, UntilEndOfTurn),
        NewDrawCardsEffect(2),
    )
    return card
}
```

**Deliverable:** 100 cards covering common patterns

### Week 9-12: Standard Format (~500 cards)

**Now you can play real decks:**

```
Standard Format Cards:
- ~500 unique cards
- All current mechanics
- Competitive decks playable

Implementation:
- 50% simple (use templates)
- 30% medium (adapt from similar cards)
- 20% complex (custom implementation)
```

**Deliverable:** Playable Standard format

## Transpiler: Build After MVP Works

**Don't build the transpiler first. Build it after you have:**
1. ✅ Working ability framework
2. ✅ 100+ manually implemented cards
3. ✅ Clear patterns identified
4. ✅ Validated approach

**Then the transpiler writes itself:**

```go
// You'll know exactly what patterns to generate
// because you've written 100 cards manually

type CardPattern struct {
    Type     string  // "simple_spell", "creature_with_keywords", etc.
    Template string  // Go code template
}

// Simple spell pattern (covers 30% of cards)
var simpleSpellPattern = `
func New{{.Name}}(id uuid.UUID) *Card {
    card := New{{.CardType}}("{{.ManaCost}}")
    {{range .Effects}}
    card.AddEffect({{.}})
    {{end}}
    return card
}
`

// Creature with keywords pattern (covers 20% of cards)
var creatureKeywordPattern = `
func New{{.Name}}(id uuid.UUID) *Card {
    card := NewCreature("{{.ManaCost}}", {{.Power}}, {{.Toughness}})
    {{range .Keywords}}
    card.AddKeyword({{.}})
    {{end}}
    return card
}
`
```

## Realistic Timeline

### Aggressive (Full-time, 3 months)

**Month 1:**
- Week 1-2: Ability framework
- Week 3-4: 20 test cards

**Month 2:**
- Week 1-2: 100 cards
- Week 3-4: 300 cards

**Month 3:**
- Week 1-2: 500 Standard cards
- Week 3-4: Testing + polish

**Result:** Playable Standard format

### Realistic (Part-time, 6 months)

**Months 1-2:**
- Ability framework
- 50 cards

**Months 3-4:**
- 200 more cards
- Patterns identified

**Months 5-6:**
- 500 Standard cards
- Transpiler prototype

**Result:** Playable Standard format + automation path

### Extended (Full card pool, 12 months)

**Months 1-6:**
- Standard format (500 cards)

**Months 7-9:**
- Transpiler
- Modern format (2,000 cards)

**Months 10-12:**
- Full card pool (31,000 cards)
- Testing + fixes

**Result:** Full parity with Java

## What to Do This Week

### Day 1: Database Setup (2 hours)

```bash
# Option A: Use Java's H2 directly
go get github.com/jmrobles/h2go
# Test reading cards from H2

# Option B: Export to SQLite
./scripts/export_h2_to_sqlite.sh
# Test reading cards from SQLite

# Option C: Import to PostgreSQL
./scripts/import_cards.go
# Test reading cards from PostgreSQL
```

**Pick ONE. Don't overthink it.**

### Day 2-3: Ability Framework (8 hours)

```go
// Start with the absolute minimum:

// 1. Ability interface
type Ability interface {
    Execute(game *Game, source *Card) error
}

// 2. One effect
type DamageEffect struct {
    Amount int
    Target Target
}

// 3. One card
func NewLightningBolt() *Card {
    card := NewInstant("{R}")
    card.AddEffect(&DamageEffect{3, TargetAny})
    return card
}

// 4. Test it
game.CastSpell(lightningBolt, target)
// Does it work? Yes? Great!
// No? Fix it.
```

### Day 4-5: Five More Cards (8 hours)

```go
// Add 5 more cards to validate framework:
- Shock (damage variant)
- Giant Growth (boost effect)
- Divination (draw effect)
- Grizzly Bears (vanilla creature)
- Llanowar Elves (mana ability)

// Each card teaches you something:
// - Different effects
// - Different targets
// - Different costs
// - Permanent vs spell
```

### End of Week: Review

**You should have:**
- ✅ Card database accessible
- ✅ Basic ability framework
- ✅ 5-10 working cards
- ✅ Tests passing
- ✅ Clear path forward

**You should NOT have:**
- ❌ All 31,000 cards
- ❌ Perfect transpiler
- ❌ Complete ability system
- ❌ All keywords

**That comes later.**

## The Truth About All Abilities

**You asked: "I want to implement all abilities for MVP"**

**Reality check:**

1. **You don't need ALL abilities for MVP**
   - Standard format: ~200 unique abilities
   - Modern format: ~500 unique abilities
   - Full card pool: ~15,000 unique abilities

2. **Many abilities are variations**
   - "Deal 3 damage" vs "Deal 5 damage" = same ability, different parameter
   - "Draw 2 cards" vs "Draw 3 cards" = same ability, different parameter
   - ~80% of abilities are variations of ~100 base patterns

3. **Transpiler makes sense after patterns are clear**
   - Write 100 cards manually
   - Identify patterns
   - Generate the rest

## My Advice

**Start small. Iterate. Scale.**

1. **This week:** Database + framework + 5 cards
2. **This month:** 50 cards
3. **Month 2:** 200 cards
4. **Month 3:** 500 cards (playable MVP)
5. **Month 4+:** Transpiler + full card pool

**Don't try to build everything at once.**

**Build what you need, when you need it.**

## Next Steps

**Want to start today?**

1. Pick a database approach (I recommend: use Java's H2 directly for now)
2. Implement the ability framework (start with 3 interfaces)
3. Implement Lightning Bolt (one card, fully working)
4. Test it in your game engine
5. Add 4 more cards

**By end of week:**
- 5 working cards
- Validated approach
- Clear path to 500 cards

**Want help with any of these steps?** I can:
- Set up H2 database reading in Go
- Design the ability framework
- Implement the first 5 cards
- Create tests for validation
