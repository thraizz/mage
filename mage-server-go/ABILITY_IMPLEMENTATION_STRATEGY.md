# Implementing ALL Abilities for MVP

## The Challenge

You want to implement all abilities so players can use any deck. This is ambitious but achievable with the right strategy.

**The Numbers:**
- 31,000+ unique cards
- ~15,000 unique abilities (many cards share abilities)
- ~500 keyword abilities (flying, trample, etc.)
- ~2,000 unique triggered abilities
- ~3,000 unique activated abilities
- ~1,000 unique static abilities

## The Java Approach (What We're Porting)

Java MAGE has **every card as a Java class**. Here's the breakdown:

```
Mage.Sets/
├── src/mage/sets/
│   ├── innistrad/
│   │   ├── DelverOfSecrets.java      (200 lines)
│   │   ├── SnapcasterMage.java       (150 lines)
│   │   └── ... (hundreds more)
│   ├── ravnica/
│   └── ... (100+ sets)
└── Total: 31,163 Java files

Mage/src/main/java/mage/abilities/
├── keyword/
│   ├── FlyingAbility.java
│   ├── TrampleAbility.java
│   └── ... (500+ keywords)
├── effects/
│   ├── common/
│   │   ├── DamageTargetEffect.java
│   │   ├── DrawCardSourceControllerEffect.java
│   │   └── ... (1000+ effects)
└── ... (massive ability framework)
```

**Key Insight:** Java has been developed over **15+ years** by hundreds of contributors.

## Realistic MVP Strategy

### Phase 1: Ability Framework (2-3 weeks)

Build the **infrastructure** that all abilities use:

```go
// Core ability system
- Ability interface (activated, triggered, static)
- Effect interface (damage, draw, counter, etc.)
- Cost interface (mana, tap, sacrifice)
- Target interface (card, player, spell)
- Condition interface (if X, then Y)
- Duration interface (until end of turn, permanent)
```

**Deliverable:** Framework that can express most abilities

### Phase 2: Common Effects Library (2-3 weeks)

Implement the **building blocks** used by 80% of cards:

```go
// ~100 common effects that cover most cards
- DamageEffect (deal X damage)
- DrawCardsEffect (draw X cards)
- GainLifeEffect (gain X life)
- CounterSpellEffect (counter target spell)
- DestroyEffect (destroy target permanent)
- ExileEffect (exile target card)
- ReturnToHandEffect (return to owner's hand)
- PutCounterEffect (add +1/+1 counters)
- TapEffect (tap target permanent)
- UntapEffect (untap target permanent)
- DiscardEffect (discard X cards)
- SacrificeEffect (sacrifice a permanent)
- SearchLibraryEffect (search for a card)
- ShuffleEffect (shuffle library)
- MillEffect (mill X cards)
- CreateTokenEffect (create token)
- GainControlEffect (gain control of target)
- BoostEffect (get +X/+X)
- ... (80 more)
```

**Deliverable:** Library of reusable effects

### Phase 3: Keyword Abilities (2-3 weeks)

Implement the **500 keyword abilities**:

```go
// Keywords are the easiest - they're standardized
- Flying, First Strike, Double Strike
- Trample, Vigilance, Haste
- Deathtouch, Lifelink, Hexproof
- Flash, Defender, Menace
- Reach, Protection, Indestructible
- ... (490 more)
```

**Deliverable:** All keyword abilities working

### Phase 4: Card Generation System (3-4 weeks)

Build a **transpiler** to convert Java cards to Go:

```
Java Card Class → Parser → Go Card Function
```

**Example:**

```java
// Java: LightningBolt.java
public final class LightningBolt extends CardImpl {
    public LightningBolt(UUID ownerId, CardSetInfo setInfo) {
        super(ownerId, setInfo, new CardType[]{CardType.INSTANT}, "{R}");
        this.getSpellAbility().addEffect(new DamageTargetEffect(3));
        this.getSpellAbility().addTarget(new TargetAnyTarget());
    }
}
```

↓ Transpiler ↓

```go
// Go: lightning_bolt.go
func NewLightningBolt(id uuid.UUID, setInfo CardSetInfo) *Card {
    card := NewCard(id, setInfo, []CardType{TypeInstant}, "{R}")
    card.SpellAbility.AddEffect(NewDamageTargetEffect(3))
    card.SpellAbility.AddTarget(NewTargetAnyTarget())
    return card
}
```

**Deliverable:** Automated card generation

### Phase 5: Generate All Cards (1-2 weeks)

Run the transpiler on all 31,000 cards:

```bash
# Generate all cards
./tools/transpiler \
  --input ../Mage.Sets/src/mage/sets/ \
  --output internal/cards/generated/ \
  --format go

# Result:
# internal/cards/generated/
#   ├── innistrad/
#   │   ├── delver_of_secrets.go
#   │   ├── snapcaster_mage.go
#   │   └── ...
#   └── ... (31,000 files)
```

**Deliverable:** All cards in Go

### Phase 6: Manual Fixes (2-4 weeks)

Some cards will need manual fixes:

- Complex interactions
- Custom mechanics
- Edge cases
- Transpiler limitations

**Estimate:** ~5-10% of cards need manual attention (~1,500-3,000 cards)

## Transpiler Design

### Architecture

```
┌─────────────────────────────────────────────────────┐
│ Transpiler Pipeline                                 │
│                                                     │
│  Java Source → Parser → AST → Analyzer → Generator │
│                                                     │
│  1. Parse Java file                                 │
│  2. Extract card properties                         │
│  3. Map abilities to Go equivalents                 │
│  4. Generate Go code                                │
│  5. Format and write file                           │
└─────────────────────────────────────────────────────┘
```

### Implementation

```go
// tools/transpiler/main.go
type Transpiler struct {
    parser      *JavaParser
    mapper      *AbilityMapper
    generator   *GoGenerator
    validator   *Validator
}

func (t *Transpiler) TranspileCard(javaFile string) (*GoCard, error) {
    // 1. Parse Java
    ast, err := t.parser.Parse(javaFile)
    if err != nil {
        return nil, err
    }
    
    // 2. Extract card info
    cardInfo := t.extractCardInfo(ast)
    
    // 3. Map abilities
    abilities := t.mapper.MapAbilities(ast.Abilities)
    
    // 4. Generate Go code
    goCode := t.generator.Generate(cardInfo, abilities)
    
    // 5. Validate
    if err := t.validator.Validate(goCode); err != nil {
        return nil, err
    }
    
    return goCode, nil
}
```

### Ability Mapping

```go
// Map Java abilities to Go
var abilityMap = map[string]string{
    "DamageTargetEffect":           "NewDamageTargetEffect",
    "DrawCardSourceControllerEffect": "NewDrawCardsEffect",
    "GainLifeEffect":               "NewGainLifeEffect",
    "CounterTargetEffect":          "NewCounterSpellEffect",
    // ... 1000+ mappings
}
```

### Example Transpilation

**Input (Java):**
```java
public final class ShockingGrasp extends CardImpl {
    public ShockingGrasp(UUID ownerId, CardSetInfo setInfo) {
        super(ownerId, setInfo, new CardType[]{CardType.INSTANT}, "{1}{R}");
        
        // Shocking Grasp deals 2 damage to target creature.
        this.getSpellAbility().addEffect(new DamageTargetEffect(2));
        this.getSpellAbility().addTarget(new TargetCreaturePermanent());
    }
}
```

**Output (Go):**
```go
package generated

func NewShockingGrasp(id uuid.UUID, setInfo CardSetInfo) *Card {
    card := NewCard(id, setInfo, []CardType{TypeInstant}, "{1}{R}")
    
    // Shocking Grasp deals 2 damage to target creature.
    spellAbility := card.GetSpellAbility()
    spellAbility.AddEffect(NewDamageTargetEffect(2))
    spellAbility.AddTarget(NewTargetCreaturePermanent())
    
    return card
}
```

## Timeline for Full Implementation

### Aggressive Timeline (3-4 months full-time)

**Month 1:**
- Week 1-2: Ability framework
- Week 3-4: Common effects library (100 effects)

**Month 2:**
- Week 1-2: Keyword abilities (500 keywords)
- Week 3-4: Transpiler design + prototype

**Month 3:**
- Week 1-2: Transpiler implementation
- Week 3-4: Generate all cards + testing

**Month 4:**
- Week 1-2: Manual fixes for complex cards
- Week 3-4: Integration testing + bug fixes

**Result:** All 31,000 cards playable

### Realistic Timeline (6-9 months part-time)

**Months 1-2:**
- Ability framework
- Common effects library
- Keyword abilities

**Months 3-4:**
- Transpiler design
- Transpiler implementation
- Test on 100 sample cards

**Months 5-6:**
- Generate all cards
- Automated testing
- Fix transpiler issues

**Months 7-9:**
- Manual fixes
- Edge cases
- Integration testing
- Performance optimization

## Prioritization Strategy

### Tier 1: Essential for Testing (Week 1-4)
- 20 cards manually implemented
- Basic keywords (flying, trample, haste)
- Common effects (damage, draw, destroy)
- Can play simple games

### Tier 2: Standard Format (Month 2-3)
- 500 cards from current Standard
- All Standard-legal mechanics
- Can play competitive decks

### Tier 3: Modern Format (Month 4-5)
- 2,000 most-played Modern cards
- All Modern mechanics
- Can play most popular decks

### Tier 4: Full Card Pool (Month 6+)
- All 31,000 cards
- All mechanics
- Full Legacy/Vintage/Commander support

## Recommended Approach: Hybrid

**Don't port everything manually. Don't fully automate either.**

### 1. Build Framework (Manual)
- Core ability system
- Common effects
- Keywords

### 2. Generate Simple Cards (Automated)
- Vanilla creatures
- Simple spells
- Basic abilities

### 3. Port Complex Cards (Manual)
- Planeswalkers
- Transforming cards
- Complex interactions

### 4. Iterate
- Generate → Test → Fix → Regenerate

## Practical Steps to Start

### Week 1: Foundation

```bash
# 1. Set up SQLite card database
./scripts/build_card_database.sh

# 2. Implement ability framework
# internal/game/ability.go
# internal/game/effect.go
# internal/game/cost.go
# internal/game/target.go

# 3. Implement 10 common effects
# - Damage
# - Draw
# - Destroy
# - Counter
# - Gain life
# - Discard
# - Sacrifice
# - Tap/Untap
# - Exile
# - Return to hand
```

### Week 2: First Cards

```bash
# 4. Manually implement 20 test cards
# - 5 creatures (vanilla + keywords)
# - 5 instants (simple effects)
# - 5 sorceries (simple effects)
# - 5 basic lands

# 5. Test in game engine
# - Cast spells
# - Resolve abilities
# - Verify effects work
```

### Week 3-4: Transpiler Prototype

```bash
# 6. Build transpiler for simple cards
# - Parse Java constructors
# - Extract card properties
# - Map to Go equivalents
# - Generate Go code

# 7. Test on 100 sample cards
# - Verify generated code compiles
# - Test in game engine
# - Fix issues
```

### Month 2+: Scale Up

```bash
# 8. Generate all cards
./tools/transpiler --generate-all

# 9. Automated testing
go test ./internal/cards/generated/...

# 10. Manual fixes
# - Review failed tests
# - Fix complex cards
# - Update transpiler
```

## Is This Realistic?

**Yes, with caveats:**

✅ **Achievable:**
- Ability framework: 2-3 weeks
- Common effects: 2-3 weeks
- Keywords: 2-3 weeks
- Transpiler: 3-4 weeks
- Generation: 1-2 weeks
- **Total: 3-4 months full-time**

⚠️ **Challenges:**
- Complex interactions
- Edge cases
- Testing thoroughness
- Performance optimization
- Bug fixes

❌ **Not realistic:**
- Doing it manually (15+ years of work)
- Perfect parity immediately (needs iteration)
- Zero bugs (needs extensive testing)

## Recommended MVP Scope

**For a playable MVP, you don't need ALL cards immediately:**

### MVP Option 1: Standard Format Only
- ~500 cards
- 2-3 months
- Can play current competitive decks
- **Recommended for MVP**

### MVP Option 2: Modern Format
- ~2,000 cards
- 4-5 months
- Can play most popular decks
- Good balance of scope and effort

### MVP Option 3: Full Card Pool
- ~31,000 cards
- 6-9 months
- Can play any deck
- **Your stated goal**

## My Recommendation

**Start with MVP Option 1, build toward Option 3:**

1. **Month 1-2:** Framework + 500 Standard cards
   - Playable MVP
   - Test with real users
   - Validate approach

2. **Month 3-4:** Transpiler + Modern cards
   - Expand to 2,000 cards
   - Prove automation works
   - Refine transpiler

3. **Month 5-6:** Generate all cards
   - Full 31,000 cards
   - Manual fixes
   - Comprehensive testing

**This gives you:**
- ✅ Playable MVP in 2 months
- ✅ Validation of approach
- ✅ Iterative development
- ✅ Full card pool in 6 months

## Next Steps

1. **This week:** Set up SQLite card database
2. **Week 1-2:** Implement ability framework
3. **Week 3-4:** Implement 20 test cards manually
4. **Week 5-6:** Build transpiler prototype
5. **Month 2:** Generate Standard cards (500)
6. **Month 3+:** Scale to full card pool

**Want to proceed?** I can help you:
1. Set up SQLite card database
2. Design the ability framework
3. Implement the first 20 cards
4. Build the transpiler prototype
