# MTG Game Engine - Gap Analysis

## Executive Summary

This document identifies missing components in the Go port of XMage's game engine based on the MTG Comprehensive Rules and comparison with the complete rule set.

## Analysis Methodology

- ✅ **Implemented**: Component exists with tests
- ⚠️ **Partial**: Component exists but incomplete
- ❌ **Missing**: Component not implemented
- 🔧 **In Progress**: Currently being worked on

---

## 1. Core Game Concepts (Rules 100-123)

### ✅ Fully Implemented
- **100. General** - Basic game rules framework exists
- **102. Players** - Player management with life, poison, energy
- **103. Starting the Game** - Game initialization in `StartGame()`
- **104. Ending the Game** - Win/loss conditions with `checkIfGameIsOver()`
- **105. Colors** - Color system exists
- **106. Mana** - `internal/game/mana/pool.go` with full mana pool support
- **107. Numbers and Symbols** - Parsing in place
- **108. Cards** - Complete card structure in `card.go`
- **109. Objects** - Object model implemented
- **110. Permanents** - Permanent tracking on battlefield
- **111. Tokens** - `internal/game/token/` package
- **122. Counters** - `internal/game/counters/` package

### ⚠️ Partially Implemented
- **119. Life** - Basic life tracking, missing life loss triggers
- **120. Damage** - Damage dealing exists, missing prevention/replacement
- **121. Drawing a Card** - Basic draw, missing empty library loss tracking

### ❌ Missing
- **101. The Magic Golden Rules** - No explicit rule priority system
- **113. Abilities** - ✅ Now integrated (abilities package)
- **116. Special Actions** - Partially in `rules/special_action.go` but incomplete:
  - Missing: Play a land special action
  - Missing: Turn face-down creatures face up
  - Missing: Suspend cards from hand (Time Spiral mechanic)
- **117. Timing and Priority** - ⚠️ Exists but needs integration with abilities
- **118. Costs** - ⚠️ Basic cost payment, missing complex cost types
- **123. Stickers** - ❌ Not implemented (Un-set mechanic)

---

## 2. Turn Structure (Rules 500-514)

### ✅ Fully Implemented
- **500. General** - Turn structure in `rules/turn.go`
- **501. Beginning Phase** - Beginning phase tracking
- **502. Untap Step** - Untap logic
- **503. Upkeep Step** - Upkeep triggers
- **504. Draw Step** - Draw step logic
- **505. Main Phase** - Main phase tracking
- **506. Combat Phase** - Full combat phase structure
- **507. Beginning of Combat Step** - Begin combat triggers
- **508. Declare Attackers Step** - ✅ Attacker declaration with triggers
- **509. Declare Blockers Step** - ✅ Blocker declaration with triggers
- **510. Combat Damage Step** - Combat damage resolution
- **511. End of Combat Step** - End combat cleanup
- **512. Ending Phase** - End phase tracking
- **513. End Step** - End step triggers
- **514. Cleanup Step** - Cleanup with discard to hand size

### ⚠️ Needs Enhancement
- **510. Combat Damage Step** - Missing first strike/double strike dynamic step insertion

---

## 3. Spells, Abilities, and Effects (Rules 600-618)

### ✅ Fully Implemented
- **601. Casting Spells** - ✅ Complete in `engine_abilities.go`
- **602. Activating Activated Abilities** - ✅ Complete in `engine_abilities.go`
- **605. Mana Abilities** - ✅ Implemented with no-stack handling
- **608. Resolving Spells and Abilities** - ✅ Resolution logic in stack manager
- **613. Interaction of Continuous Effects** - ✅ Layer system in `engine_layers.go`

### ⚠️ Partially Implemented
- **603. Handling Triggered Abilities** - Basic triggers exist, needs full integration:
  - ✅ ETB/LTB triggers
  - ✅ Combat triggers (attacks, blocks, damage)
  - ⚠️ Missing: "At beginning of X step" triggers
  - ⚠️ Missing: "Whenever you cast a spell" triggers
  - ⚠️ Missing: State triggers (threshold, metalcraft, etc.)
- **604. Handling Static Abilities** - Static abilities exist but need more coverage:
  - ✅ P/T modification
  - ⚠️ Missing: Rules-changing effects ("As long as...")
  - ⚠️ Missing: Continuous type-changing effects
- **606. Loyalty Abilities** - Planeswalker abilities partially implemented
- **609. Effects** - Effect framework exists but incomplete
- **610. One-Shot Effects** - Basic one-shot effects work
- **611. Continuous Effects** - Layer system exists but missing some effect types
- **612. Text-Changing Effects** - ❌ Not implemented

### ❌ Missing
- **607. Linked Abilities** - No linked ability tracking (Champion, Exert, etc.)
- **614. Replacement Effects** - ⚠️ Exists in `effects/replacement.go` but incomplete:
  - Missing: "Instead" replacement effects
  - Missing: "As X enters the battlefield" effects
  - Missing: "If X would" effects
- **615. Prevention Effects** - ❌ No prevention effect system
- **616. Interaction of Replacement/Prevention** - ❌ No interaction system

---

## 4. Additional Rules (Rules 700-729)

### ✅ Fully Implemented
- **704. State-Based Actions** - Comprehensive SBA checker:
  - Life ≤ 0
  - Poison ≥ 10
  - Toughness ≤ 0
  - Lethal damage
  - Deathtouch damage
  - Planeswalker loyalty = 0
  - Legend rule
  - +1/+1 and -1/-1 annihilation

### ⚠️ Partially Implemented
- **701. Keyword Actions** - Many missing:
  - ✅ Destroy
  - ✅ Exile
  - ✅ Tap/Untap
  - ✅ Counter (spell)
  - ⚠️ Missing: Mill, Surveil, Scry, Fateseal
  - ⚠️ Missing: Proliferate, Amass, Adapt
  - ⚠️ Missing: Transform, Meld, Mutate
  - ⚠️ Missing: Venture into dungeon, Initiative
  - ⚠️ Missing: 50+ other keyword actions

- **702. Keyword Abilities** - Basic keywords exist, many missing:
  - ✅ Flying, First Strike, Double Strike
  - ✅ Deathtouch, Lifelink, Vigilance
  - ✅ Trample, Defender, Reach
  - ✅ Menace, Hexproof, Shroud
  - ⚠️ Missing: Flash, Haste, Indestructible
  - ⚠️ Missing: Protection from X
  - ⚠️ Missing: Convoke, Delve, Affinity
  - ⚠️ Missing: 100+ other keyword abilities

- **703. Turn-Based Actions** - Partially implemented:
  - ✅ Untap step actions
  - ✅ Draw step actions
  - ⚠️ Missing: Phasing
  - ⚠️ Missing: Saga counter placement

### ❌ Missing Completely
- **705. Flipping a Coin** - ❌ No coin flip system
- **706. Rolling a Die** - ❌ No die rolling system
- **707. Copying Objects** - ❌ No copy effects (Clone, etc.)
- **708. Face-Down Spells and Permanents** - ❌ No morph/manifest support
- **709-721. Special Card Types** - ❌ Missing:
  - Split cards, Flip cards, Transform cards
  - Leveler cards, Double-faced cards
  - Saga cards, Adventure cards, Class cards
  - Attraction cards, Prototype cards, Case cards
- **722. Controlling Another Player** - ❌ No Mindslaver effects
- **723. Ending Turns and Phases** - ❌ No Time Stop effects
- **724. The Monarch** - ❌ Not implemented
- **725. The Initiative** - ❌ Not implemented
- **726. Restarting the Game** - ❌ No Karn Liberated restart
- **727. Rad Counters** - ❌ Not implemented (Fallout)
- **728. Subgames** - ❌ Not implemented (Shahrazad)
- **729. Merging with Permanents** - ❌ Not implemented (Mutate)

---

## 5. Critical Missing Systems

### 🔴 HIGH PRIORITY - Game-Breaking if Missing

#### 1. **Replacement Effects System** (Rule 614)
**Status**: ⚠️ Partial - exists in `effects/replacement.go` but needs completion

**Missing**:
- "As X enters the battlefield" replacement (Thalia's Lieutenant, etc.)
- "If X would draw" replacements (Abundance, etc.)
- "If X would die" replacements (Undying, Persist, etc.)
- "Instead" effect application
- Multiple replacement ordering (Rule 616)

**Impact**: Cards like Doubling Season, Anointed Procession, Elesh Norn don't work

**Effort**: Medium - framework exists, needs effect implementations

---

#### 2. **Prevention Effects System** (Rule 615)
**Status**: ❌ Missing entirely

**Missing**:
- "Prevent the next X damage" effects
- Prevention shields
- Prevention sources tracking
- Interaction with replacement effects

**Impact**: All protection effects, damage prevention spells broken

**Effort**: Medium - similar to replacement effects

---

#### 3. **Text-Changing Effects** (Rule 612)
**Status**: ❌ Missing entirely

**Missing**:
- Color/type word replacement
- Name replacement
- Ability text modification
- Layer 3 processing

**Impact**: Sleight of Mind, Glamerdye, Artificial Evolution don't work

**Effort**: Low - niche mechanic

---

#### 4. **Linked Abilities** (Rule 607)
**Status**: ❌ Missing entirely

**Missing**:
- Ability pairing tracking
- Imprint (Isochron Scepter)
- Champion (Champion of the Parish)
- Exert paired triggers
- Modal DFC ability linking

**Impact**: Dozens of cards broken

**Effort**: Medium - needs ability tracking infrastructure

---

#### 5. **Copy Effects** (Rule 707)
**Status**: ❌ Missing entirely

**Missing**:
- Copy object creation
- Copiable values tracking
- Copy modification (with exceptions)
- Layer 1 copy effects
- Token copies

**Impact**: Clone, Rite of Replication, Helm of the Host broken

**Effort**: High - complex rule interactions

---

#### 6. **Modal Spell Support**
**Status**: ⚠️ Mentioned in TODO but not implemented

**Missing**:
- Mode selection UI/API
- Mode validation
- Multiple mode selection
- Modal ability framework

**Impact**: All Charms, Commands, modal spells broken

**Effort**: Low - straightforward implementation

---

#### 7. **Cost Modifications** (Rule 601.2f-2h)
**Status**: ⚠️ Basic costs work, modifications missing

**Missing**:
- Cost increases (Thalia, Guardian of Thraben)
- Cost reductions (Goblin Electromancer, Affinity)
- Alternative costs (Flashback, Force of Will)
- Additional costs (Kicker, Buyback, Replicate)
- X cost support

**Impact**: Hundreds of cards broken

**Effort**: Medium - cost calculation framework needs extension

---

#### 8. **Face-Down Spells/Permanents** (Rule 708)
**Status**: ❌ Missing entirely

**Missing**:
- Morph mechanic
- Manifest mechanic
- Megamorph
- Face-down card state tracking
- Face-up action handling

**Impact**: All morph/manifest cards broken

**Effort**: High - complex state tracking

---

#### 9. **Transform/DFC Support** (Rule 712)
**Status**: ❌ Missing entirely

**Missing**:
- Double-faced card representation
- Transform action
- Day/night cycle
- Daybound/Nightbound
- Modal DFC support

**Impact**: All transform cards broken (Werewolves, Planeswalkers, etc.)

**Effort**: High - requires card data model changes

---

### 🟡 MEDIUM PRIORITY - Affects Many Cards

#### 10. **Keyword Action Support** (Rule 701)
**Status**: ⚠️ ~10% implemented

**Critical Missing Actions**:
- Scry, Surveil, Fateseal
- Mill
- Proliferate
- Adapt, Amass, Bolster
- Convoke, Delve, Improvise
- Connive, Casualty
- Fight (creatures deal damage to each other)

**Impact**: Hundreds of cards

**Effort**: Low-Medium per action

---

#### 11. **Keyword Ability Support** (Rule 702)
**Status**: ⚠️ ~15% implemented

**Critical Missing Abilities**:
- Flash, Haste
- Indestructible
- Protection from X
- Hexproof from X
- Ward X
- Cascade, Storm, Affinity, Convoke, Delve
- Phasing
- Flanking, Bushido, Rampage
- Cumulative upkeep
- Echo, Fading, Vanishing

**Impact**: Thousands of cards

**Effort**: Low-Medium per ability

---

#### 12. **Triggered Ability Completeness**
**Status**: ⚠️ ~40% implemented

**Missing Trigger Types**:
- "At the beginning of [phase/step]"
- "Whenever you cast a spell"
- "Whenever X becomes the target"
- "Whenever X is dealt damage"
- "Whenever a creature you control dies"
- State-based triggers (Threshold, Metalcraft, etc.)
- "At end of turn" triggers

**Impact**: Thousands of cards

**Effort**: Low per trigger type

---

### 🟢 LOW PRIORITY - Niche Mechanics

#### 13. **Special Card Types**
- Split cards (Fire // Ice)
- Saga cards
- Adventure cards
- Class cards
- Attractions, Stickers
- Prototype cards
- Case cards

**Impact**: Specific card sets

**Effort**: Medium-High per type

---

#### 14. **Special Game States**
- The Monarch
- The Initiative
- Dungeons
- City's Blessing
- Day/Night
- Rad Counters

**Impact**: Specific mechanics

**Effort**: Low-Medium per state

---

#### 15. **Random Events**
- Coin flips
- Die rolls

**Impact**: Specific cards

**Effort**: Low

---

## 6. Integration Gaps

### ⚠️ Engine Component Integration Issues

#### 1. **Triggered Ability → Stack**
**Status**: ⚠️ Triggers detected but not pushed to stack

**Missing**:
- Controller information for triggers
- Automatic trigger resolution
- "You may" optional trigger handling
- Simultaneous trigger ordering (APNAP)

**Location**: `engine_combat.go` lines 172-174, 219-220, 274-275, 323-324, 369-370

---

#### 2. **Target Selection → Spell Casting**
**Status**: ⚠️ System exists but not integrated

**Missing**:
- UI/API for target selection
- Integration with `CastSpell()`
- Multi-target validation
- Target redirection

---

#### 3. **Layer Effects → Card State**
**Status**: ⚠️ Layer system exists but missing effect types

**Missing Layer Effect Conversions**:
- Layer 2: Control-changing (Threaten, Act of Treason)
- Layer 4: Type-changing (Opalescence, March of the Machines)
- Layer 5: Color-changing (Painter's Servant, Moonlace)
- Layer 6: Ability-adding (Archetype of Courage, etc.)
- Layer 7b: Set P/T (Chameleon Colossus, Tarmogoyf)

**Effort**: Low-Medium per effect type

---

#### 4. **Mana Pool → Mana Abilities**
**Status**: ✅ Integrated but needs floating mana expiration

**Missing**:
- Mana pool emptying at step boundaries
- Mana burn (obsolete but could be optional)
- Conditional mana (can only be spent on X)

---

## 7. Performance & Scalability Issues

### Identified Bottlenecks

1. **Layer Recalculation**
   - **Current**: Recalculated before every SBA check
   - **Issue**: O(n) permanents × O(m) effects per check
   - **Solution**: Dirty tracking, incremental updates

2. **Ability Registry**
   - **Current**: Single `sync.RWMutex` for 30,400+ cards
   - **Issue**: Lock contention with many permanents
   - **Solution**: Sharding, lock-free data structures

3. **Event Processing**
   - **Current**: Linear event processing
   - **Issue**: Slow with many watchers/triggers
   - **Solution**: Event subscription filtering

---

## 8. Testing Gaps

### Test Coverage Analysis

**Excellent Coverage** (30+ tests each):
- Combat system (30+ test files)
- Turn structure
- Priority and stack

**Good Coverage** (5-10 tests):
- Abilities (effects, counters, tokens)
- Rules engine (triggers, watchers)

**Minimal Coverage** (<5 tests):
- Replacement effects
- Layer system integration
- Ability activation workflow

**No Coverage**:
- Prevention effects
- Text-changing effects
- Copy effects
- Face-down permanents
- Transform/DFC
- Most keyword actions
- Most keyword abilities

---

## 9. Priority Roadmap

### Phase 1: Critical Game Functionality (2-3 weeks)
1. ✅ Complete ability activation integration (DONE)
2. ✅ Fix circular import dependencies (DONE)
3. Push triggered abilities to stack
4. Modal spell support
5. Cost modification system
6. Replacement effects completion

### Phase 2: Common Card Support (3-4 weeks)
1. Prevention effects system
2. Copy effects (Clone, etc.)
3. Critical keyword actions (Scry, Mill, Proliferate)
4. Critical keyword abilities (Flash, Haste, Indestructible, Ward)
5. Linked abilities tracking
6. "At beginning of" triggers

### Phase 3: Advanced Mechanics (4-6 weeks)
1. Face-down permanents (Morph/Manifest)
2. Transform/DFC support
3. Text-changing effects
4. Complete keyword action library
5. Complete keyword ability library
6. Special card types (Split, Adventure, Saga)

### Phase 4: Niche Mechanics (2-3 weeks)
1. Special game states (Monarch, Initiative, Dungeons)
2. Random events (coin flips, die rolls)
3. Advanced mechanics (Subgames, Turn ending, Player control)
4. Un-set support (Stickers, Attractions)

### Phase 5: Optimization & Polish (2-3 weeks)
1. Layer recalculation optimization
2. Ability registry sharding
3. Event system optimization
4. Comprehensive integration tests
5. Performance benchmarking

---

## 10. Estimated Completion

**Current State**: ~40% complete for competitive play
- Core engine: 70%
- Spell casting: 80%
- Combat: 85%
- Abilities: 60%
- Effects: 35%
- Keywords: 15%
- Special mechanics: 5%

**To Minimum Viable Game (MVP)**: 6-8 weeks
- Replacement effects
- Prevention effects
- Cost modifications
- Modal spells
- Basic keyword library (top 30 most common)
- Triggered ability integration

**To Comprehensive Support**: 14-18 weeks
- All critical missing systems
- 90%+ keyword action/ability coverage
- Special card type support
- Transform/DFC support
- Optimization pass

**To 100% Rules Coverage**: 24-30 weeks
- All MTG comprehensive rules
- All special mechanics
- Un-set support
- Performance optimization
- Full test coverage

---

## 11. Recommendations

### Immediate Actions
1. **Push Triggered Abilities to Stack** - Required for basic gameplay
2. **Modal Spell Support** - Blocks many common cards
3. **Cost Modifications** - Affects hundreds of cards
4. **Replacement Effects Completion** - Critical for EDH/Commander

### Architecture Improvements
1. **Effect Type System** - Unified interface for all effect types
2. **Ability Builder Enhancement** - Support all ability patterns
3. **Event Subscription** - Efficient trigger/watcher system
4. **Layer Effect Registry** - Automated layer effect registration

### Development Process
1. **Test-Driven Development** - Write tests before implementing
2. **Incremental Integration** - Small, testable changes
3. **Performance Monitoring** - Benchmark critical paths
4. **Documentation Updates** - Keep CLAUDE.md synchronized

---

## Conclusion

The Go port engine has a **solid foundation** with:
- ✅ Complete turn structure
- ✅ Comprehensive combat system
- ✅ Stack and priority management
- ✅ Ability activation workflows
- ✅ Layer system framework

**Critical gaps** blocking MVP:
- Replacement/prevention effects
- Cost modifications
- Triggered ability stack integration
- Modal spell support
- Keyword library expansion

**Estimated effort to MVP**: 6-8 weeks with focused development.

The codebase is well-structured and follows good patterns, making feature additions straightforward. The main challenge is the **breadth of MTG rules** rather than architectural complexity.
