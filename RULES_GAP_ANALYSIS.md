# MTG Comprehensive Rules - Go Engine Gap Analysis

**Date**: 2025-01-24
**Rules Version**: September 19, 2025
**Engine Status**: Based on GO_PORT_TASKS.md audit

This document identifies gaps between the MTG Comprehensive Rules and the Go engine implementation.

---

## Executive Summary

**Total Keyword Abilities in Rules**: 188 (Rule 702.1 through 702.188)
**Implemented Keywords**: 14 basic keywords
**Implementation Rate**: ~7.4%

**Critical Gaps**:
- 174 keyword abilities missing
- Replacement/prevention effects system absent
- Copy effects system absent
- Several special card types not implemented
- Multiplayer rules not implemented
- Advanced mechanics (phasing, protection, etc.) missing

---

## 1. Core Game Rules (Rules 100-123)

### ✅ Fully Implemented

- **Rule 103**: Starting the game (mulligan system complete)
- **Rule 104**: Ending the game (win/loss conditions)
- **Rule 117**: Timing and priority (priority windows, passing)
- **Rule 118**: Costs (mana costs, additional costs, payment)
- **Rule 119**: Life (gaining, losing tracked)
- **Rule 120**: Damage (damage dealing, combat damage)
- **Rule 121**: Drawing a card
- **Rule 122**: Counters (70+ counter types, placement, removal)
- **Rule 123**: Zones (all 8 zones tracked: library, hand, battlefield, graveyard, stack, exile, ante, command)

### ⚠️ Partially Implemented

- **Rule 106**: Mana
  - ✅ Mana pool exists (`internal/game/mana/pool.go`)
  - ✅ Mana types (WUBRG, colorless, generic)
  - ❌ Hybrid mana parsing incomplete
  - ❌ Phyrexian mana not implemented
  - ❌ Snow mana not tracked
  - ❌ X costs partially implemented (TODOs in mage_engine.go)

- **Rule 107**: Numbers and symbols
  - ✅ Basic number handling
  - ❌ Infinity symbol not implemented (Rule 702.186)
  - ❌ X, Y, Z variable costs incomplete

- **Rule 111**: Tokens
  - ✅ Token creation effect exists (`CreateTokenEffect`)
  - ✅ 711 token types registered
  - ❌ Token characteristics modification incomplete
  - ❌ Token copying not implemented

- **Rule 115**: Targets
  - ✅ Target selection flow exists (`engine_targeting.go`)
  - ✅ Target validation basic implementation
  - ❌ Hexproof checking stubbed (TODO line 121 in targeting/validator.go)
  - ❌ Shroud not checked
  - ❌ Protection not implemented
  - ❌ Ward not implemented

- **Rule 116**: Special Actions
  - ✅ Playing lands
  - ✅ Morph/face-down creature actions (Rule 116.2b)
  - ✅ Companion (Rule 116.2g)
  - ✅ Plot (Rule 116.2k)
  - ❌ Suspend (Rule 116.2e) - not mentioned
  - ❌ Foretell (Rule 116.2f) - not mentioned
  - ❌ Meld (Rule 116.2h) - not mentioned
  - ❌ Escape (Rule 116.2i) - not mentioned
  - ❌ Adventure (Rule 116.2j) - not mentioned
  - ❌ Outlaws ability (Rule 116.2n) - not mentioned
  - ❌ Craft (Rule 116.2p) - not mentioned
  - ❌ Disguise (Rule 116.2q) - not mentioned

### ❌ Not Implemented

- **Rule 101**: The Magic Golden Rules
  - Card-specific rules override general rules (needs explicit checking)

- **Rule 102**: Players
  - ✅ Starting life tracked
  - ❌ Maximum hand size not enforced (default 7)
  - ❌ Shared team resources (Two-Headed Giant, etc.)

- **Rule 105**: Colors
  - ❌ Color identity for Commander not tracked
  - ❌ Multicolor vs hybrid distinction unclear

- **Rule 108-110**: Cards, Objects, Permanents
  - ❌ Card characteristics modification incomplete
  - ❌ Permanent timestamp tracking (for layer ordering)
  - ❌ Continuous effects dependency ordering

- **Rule 112**: Spells
  - ✅ Basic casting implemented
  - ❌ Alternative costs incomplete
  - ❌ Additional costs incomplete
  - ❌ Cost reductions incomplete

- **Rule 113**: Abilities
  - ✅ 6 ability types defined (spell, activated, triggered, static, mana, keyword)
  - ❌ Ability removal/granting edge cases incomplete

- **Rule 114**: Emblems
  - ❌ Emblem objects not implemented
  - ❌ Command zone emblem tracking missing

---

## 2. Turn Structure (Rules 500-514)

### ✅ Fully Implemented

- **Rule 500-505**: Beginning phase (untap, upkeep, draw) - complete turn structure in `turn.go`
- **Rule 506-511**: Combat phase (all 6 combat steps implemented)
- **Rule 512-514**: Ending phase (end step, cleanup step)
- **Rule 502**: Untap step with phasing stub (phasing not implemented)
- **Rule 510**: Combat damage step with first strike/double strike

### ⚠️ Partially Implemented

- **Rule 502.1**: Phasing
  - ❌ No phasing implementation (0 matches in codebase)
  - Creatures don't phase in/out during untap step

- **Rule 514**: Cleanup Step
  - ✅ Discard to max hand size
  - ✅ Remove damage from creatures
  - ❌ Simultaneous trigger handling in cleanup incomplete
  - ❌ Until end of turn duration expiration

### ❌ Not Implemented

- **Rule 500.11**: Turn-based action priority (some actions don't use stack)
- **Rule 505.4a**: Playing land limit per turn (tracking unclear)

---

## 3. Spell Casting & Abilities (Rules 601-609)

### ✅ Fully Implemented

- **Rule 601**: Casting spells (sequence complete in `engine_abilities.go`)
- **Rule 602**: Activating activated abilities (complete workflow)
- **Rule 603**: Handling triggered abilities (APNAP order, queue system)
- **Rule 604**: Handling static abilities (layer system integrated)
- **Rule 605**: Mana abilities (special timing, no stack, immediate resolution)
- **Rule 608**: Resolving spells and abilities (full resolution sequence)

### ⚠️ Partially Implemented

- **Rule 606**: Loyalty abilities
  - ✅ Loyalty counters exist
  - ✅ Planeswalker combat damage reduces loyalty
  - ❌ Loyalty ability activation restrictions (once per turn per planeswalker) not tracked
  - ❌ Loyalty cost payment (add/remove loyalty) incomplete

- **Rule 607**: Linked abilities
  - ✅ Basic linked ability support for casting during resolution
  - ❌ Paired abilities (e.g., exile/return pairs) not fully implemented
  - ❌ Characteristic-defining abilities incomplete

- **Rule 609**: Effects
  - ✅ 40+ effect types implemented
  - ❌ Effect modification (doubling, prevention) incomplete
  - ❌ Effect ordering edge cases

### ❌ Not Implemented

- **Rule 601.2f**: Alternative costs (Flashback, Escape, etc.)
- **Rule 601.2h**: Cost increases/reductions
  - ❌ Cost modification effects not tracked
  - ❌ Convoke, Delve, Improvise reduce costs - not implemented

---

## 4. Continuous Effects (Rules 610-616) - CRITICAL GAP

### ✅ Implemented

- **Rule 611**: Continuous effects
  - ✅ Duration system (UntilEndOfTurn, WhileOnBattlefield, etc.)
  - ✅ Layer system framework exists

- **Rule 613**: Interaction of continuous effects (LAYERS)
  - ✅ 7-layer system implemented (`engine_layers.go`)
  - ✅ Layer 1-7 constants defined
  - ✅ Automatic recalculation before SBAs
  - ⚠️ Timestamp ordering within layers incomplete
  - ⚠️ Dependency system not implemented (Rule 613.8)

### ❌ Not Implemented (CRITICAL)

- **Rule 612**: Text-changing effects
  - ❌ No text-changing effect system (0 matches)
  - Example: "Change target of spell" effects missing

- **Rule 614**: Replacement effects
  - ❌ No replacement effect system (0 matches)
  - Examples: "If ~ would die, instead...", "enters with counters", doubling effects
  - **CRITICAL**: This affects hundreds of cards

- **Rule 615**: Prevention effects
  - ❌ No prevention effect system (0 matches)
  - Examples: "Prevent the next 3 damage...", protection
  - **CRITICAL**: This affects combat and damage resolution

- **Rule 616**: Interaction of replacement/prevention effects
  - ❌ Ordering and application not implemented

**Impact**: Replacement and prevention effects are fundamental to MTG. Without them:
- Protection doesn't work
- Damage prevention doesn't work
- Doubling effects don't work
- "Would die" triggers don't work
- ETB replacement effects don't work

---

## 5. State-Based Actions (Rule 704)

### ✅ Fully Implemented (Rule 704.5)

Implemented in `rules/state_based_actions.go`:
- 704.5a: Player with 0 or less life loses
- 704.5b: Player with 10+ poison counters loses
- 704.5c: Player attempted to draw from empty library loses
- 704.5d: Token in non-battlefield zone ceases to exist
- 704.5f: Creature with toughness 0 or less is put into graveyard
- 704.5g: Creature with lethal damage is destroyed
- 704.5h: Creature with deathtouch damage is destroyed
- 704.5i: Planeswalker with loyalty 0 is put into graveyard
- 704.5j: Legend rule (one legendary permanent per name per player)
- 704.5k: Planeswalker uniqueness rule (deprecated but checked)
- 704.5m: Aura attached to illegal object is put into graveyard
- 704.5n: Equipment attached to illegal object becomes unattached
- 704.5p: Creature/planeswalker attached to object unattaches
- 704.5q: +1/+1 and -1/-1 counter annihilation

### ❌ Not Implemented

- 704.5e: Copy of spell not on stack ceases to exist (copy effects not implemented)
- 704.5r: Battle with no defense counters is put into graveyard (battles not implemented)
- 704.5s: Creature with -X/-0 from suspected status (not implemented)
- 704.5t: World rule (one World enchantment) - obsolete but in rules
- 704.6: Simultaneous SBA handling edge cases

---

## 6. Card Types (Rules 300-315)

### ✅ Implemented

- **Rule 302**: Creatures (power/toughness, combat)
- **Rule 303**: Enchantments (Auras with enchant ability)
- **Rule 304**: Instants (instant-speed casting)
- **Rule 305**: Lands (basic land types, mana abilities)
- **Rule 307**: Sorceries (sorcery-speed casting)

### ⚠️ Partially Implemented

- **Rule 301**: Artifacts
  - ✅ Basic artifacts work
  - ❌ Artifact creatures type line unclear
  - ❌ Artifact lands special rules unclear

- **Rule 306**: Planeswalkers
  - ✅ Planeswalker type exists
  - ✅ Combat damage to planeswalkers (85 mentions in code)
  - ✅ Loyalty tracking
  - ❌ Loyalty ability activation (once per turn limit) not enforced
  - ❌ Planeswalker uniqueness rule implementation unclear
  - ❌ Loyalty costs incomplete

### ❌ Not Implemented (CRITICAL for Modern Cards)

- **Rule 308**: Kindreds (formerly Tribal)
  - ❌ Kindred card type not implemented
  - ❌ Creature type sharing not tracked

- **Rule 309**: Dungeons
  - ❌ No dungeon implementation (0 matches)
  - ❌ Venture into dungeon mechanic missing
  - ❌ Undercity not implemented

- **Rule 310**: Battles
  - ⚠️ 347 mentions in code but marked TODO
  - ❌ Defense counters not fully implemented
  - ❌ Siege battles not attackable (TODO line 4544)
  - ❌ Battle protectors not tracked

- **Rule 311-315**: Special formats
  - ❌ Planes (Planechase) - 0 matches
  - ❌ Phenomena (Planechase) - 0 matches
  - ❌ Vanguards - 0 matches
  - ❌ Schemes (Archenemy) - 0 matches
  - ❌ Conspiracies (Conspiracy Draft) - 0 matches

---

## 7. Keyword Abilities (Rule 702) - MAJOR GAP

**Total Keywords in Rules**: 188 (702.1 through 702.188)
**Implemented**: 14 keywords

### ✅ Implemented (14 keywords)

From `abilities/keyword.go`:
1. Flying (702.9)
2. First Strike (702.7)
3. Double Strike (702.4)
4. Deathtouch (702.2)
5. Haste (702.10)
6. Hexproof (702.11) - **Checking not implemented**
7. Indestructible (702.12) - **Not enforced in SBAs**
8. Lifelink (702.15)
9. Menace (702.111)
10. Reach (702.17)
11. Trample (702.19)
12. Vigilance (702.20)
13. Defender (702.3)
14. Flash (702.8)

### ❌ Missing Critical Keywords (174 keywords)

**High Priority Missing** (commonly used):
- Protection (702.16) - affects targeting, damage, blocking
- Ward (702.21) - triggers when targeted
- Shroud (702.18) - can't be targeted
- Intimidate (702.13) - blocking restriction
- Landwalk (702.14) - unblockable conditional
- Shadow (702.28) - can only block/be blocked by shadow
- Horsemanship (702.31) - can only be blocked by horsemanship
- Fear (702.36) - blocking restriction
- Prowess (702.107) - triggers on noncreature spell
- Toxic (702.164) - poison counters on combat damage
- Backup (702.165) - enters with effect on another creature

**Ability-Granting Keywords** (very common):
- Kicker (702.33) - additional cost for additional effect
- Flashback (702.34) - cast from graveyard
- Cycling (702.29) - discard to draw
- Convoke (702.51) - tap creatures to help pay
- Delve (702.66) - exile cards from graveyard to help pay
- Escape (702.138) - cast from graveyard with exile cost
- Madness (702.35) - cast when discarded
- Cascade (702.85) - cast free spells
- Storm (702.40) - copies for each spell cast
- Suspend (702.62) - exile with time counters

**Modal/Transform Keywords**:
- Morph (702.37) - ✅ IMPLEMENTED (face_down.go exists)
- Disguise (702.168) - variant of morph
- Megamorph (702.37) - morph with +1/+1 counter
- Manifest (not in 702 - in Rule 701.34) - ⚠️ Implemented in face_down.go
- Transform/DFC (712) - double-faced cards
- Meld (709) - combine two cards
- Mutate (702.140) - merge creatures
- Daybound/Nightbound (702.145) - transform cycle

**Cost Reduction Keywords**:
- Affinity (702.41) - cost reduction
- Improvise (702.126) - tap artifacts to pay
- Emerge (702.119) - sacrifice creature to reduce cost
- Offering (702.48) - sacrifice for cost reduction

**Token/Counter Keywords**:
- Modular (702.43) - +1/+1 counters, transfer on death
- Graft (702.58) - move +1/+1 counters
- Fabricate (702.123) - choose +1/+1 counters or tokens
- Training (702.149) - gets +1/+1 counter when attacks
- Renown (702.112) - gets +1/+1 counters on combat damage

**Graveyard Mechanics**:
- Unearth (702.84) - return from graveyard
- Embalm (702.128) - create token copy from graveyard
- Eternalize (702.129) - create 4/4 token from graveyard
- Encore (702.141) - create tokens for each opponent
- Disturb (702.146) - cast transformed from graveyard

**Combat Keywords**:
- Banding (702.22) - ⚠️ Partial implementation (29 mentions, damage control only)
- Flanking (702.25) - gives -1/-1 when blocked
- Bushido (702.45) - pump when blocks/blocked
- Rampage (702.23) - pump for each blocker
- Provoke (702.39) - force creature to block
- Ninjutsu (702.49) - swap attacking creature

**Mana Keywords**:
- Convoke (702.51) - tap creatures for mana
- Delve (702.66) - exile for generic mana
- Improvise (702.126) - tap artifacts for mana

**Triggered Keywords**:
- Exalted (702.83) - pump sole attacker
- Extort (702.101) - pay mana when casting spell
- Heroic (702.108) - trigger on targeting
- Constellation (not keyword - ability word)
- Revolt (not keyword - ability word)

**Recent Keywords** (2020+):
- Companion (702.139) - ✅ IMPLEMENTED
- Foretell (702.143) - exile for later casting
- Boast (702.142) - activate after attack
- Disturb (702.146) - transform from graveyard
- Cleave (702.148) - choose to ignore text
- Prototype (702.160) - alternative casting
- Craft (702.167) - exile to transform
- Plot (702.170) - ✅ IMPLEMENTED
- Saddle (702.171) - tap creatures to activate
- Spree (702.172) - choose modes with costs
- Offspring (702.175) - create token copy
- Impending (702.176) - suspend variant

---

## 8. Advanced Rules (Rules 700-732)

### ❌ Not Implemented (CRITICAL)

- **Rule 700.4, 707**: Copy effects
  - ❌ No copy effect system (0 matches for "CopyEffect")
  - ❌ Clone effects don't work
  - ❌ Copy spell effects don't work
  - ❌ Populate (create token copy) incomplete

- **Rule 701**: Keyword actions
  - ✅ Basic actions: destroy, exile, tap, untap
  - ❌ Attach/detach incomplete
  - ❌ Transform not implemented
  - ❌ Meld not implemented
  - ❌ Manifest partially implemented
  - ❌ Explore, surveil, adapt, amass, etc. missing

- **Rule 705**: Flipping a coin
  - ❌ No coin flip system

- **Rule 706**: Rolling a die
  - ❌ No die rolling system

- **Rule 708**: Face-down spells/permanents
  - ✅ IMPLEMENTED (`face_down.go` with morph, manifest, cloak, disguise)

- **Rule 709-713**: Special card layouts
  - ❌ Split cards (709) - not implemented
  - ❌ Flip cards (710) - not implemented
  - ❌ Leveler cards (711) - not implemented
  - ❌ Double-faced cards (712) - 7 mentions, likely stubs
  - ❌ Substitute cards (713) - not needed for digital

- **Rule 714**: Saga cards
  - ❌ No saga implementation (0 matches)
  - ❌ Lore counters not tracked
  - ❌ Chapter abilities not implemented

- **Rule 715**: Adventurer cards
  - ❌ No adventure implementation (0 matches)
  - ❌ Adventure frame not supported

- **Rule 716**: Class cards
  - ❌ No class implementation
  - ❌ Class levels not tracked

- **Rule 717**: Attraction cards
  - ❌ Not implemented (Un-sets)

- **Rule 718**: Prototype cards
  - ❌ Not implemented

- **Rule 719-721**: Recent mechanics
  - ❌ Case cards (719) - not implemented
  - ❌ Omen cards (720) - not implemented
  - ❌ Station cards (721) - not implemented

- **Rule 722**: Controlling another player
  - ❌ Mindslaver effects not implemented

- **Rule 723**: Ending turns/phases
  - ❌ Time Stop effects not implemented
  - ❌ End the turn effects missing

- **Rule 724**: The Monarch
  - ❌ Monarch designation not tracked
  - ❌ Monarch draw not implemented

- **Rule 725**: The Initiative
  - ❌ Initiative not tracked
  - ❌ Undercity not implemented

- **Rule 726**: Restarting the game
  - ❌ Karn Liberated ultimate not possible

- **Rule 727**: Rad counters
  - ❌ Rad counters not in counter list
  - ❌ Milling from rad not implemented

- **Rule 728**: Subgames
  - ❌ Shahrazad not supported

- **Rule 729**: Merging with permanents
  - ❌ Mutations not implemented

- **Rule 730**: Day and Night
  - ❌ Day/Night tracking not implemented
  - ❌ Daybound/Nightbound not implemented

- **Rule 731**: Taking shortcuts
  - ✅ Partially implemented (passing priority until event)

- **Rule 732**: Handling illegal actions
  - ⚠️ Basic error handling exists
  - ❌ Rewind mechanism incomplete

---

## 9. Multiplayer Rules (Rules 800-811) - COMPLETELY MISSING

### ❌ Not Implemented

- **Rule 800**: General multiplayer rules
  - ❌ Turn order for 3+ players
  - ❌ Attack restrictions in multiplayer
  - ❌ Effects that target "each opponent"

- **Rule 801**: Limited range of influence
  - ❌ Not implemented (0 matches)

- **Rule 802**: Attack multiple players
  - ❌ Not implemented

- **Rule 803**: Attack left/right options
  - ❌ Not implemented

- **Rule 804**: Deploy creatures option
  - ❌ Not implemented

- **Rule 805**: Shared team turns
  - ❌ Not implemented

- **Rule 806**: Free-for-All variant
  - ⚠️ Game type registered but rules not enforced

- **Rule 807**: Grand Melee variant
  - ❌ Not implemented

- **Rule 808**: Team vs. Team
  - ❌ Not implemented

- **Rule 809**: Emperor variant
  - ❌ Not implemented

- **Rule 810**: Two-Headed Giant
  - ⚠️ Game type registered but rules not enforced
  - ❌ Shared life total not implemented
  - ❌ Shared turn structure not implemented

- **Rule 811**: Alternating Teams
  - ❌ Not implemented

**Impact**: Multiplayer games beyond 1v1 duel won't work correctly.

---

## 10. Casual Variants (Rules 900-905)

### ❌ Not Implemented

- **Rule 901**: Planechase
  - ❌ Plane cards not supported
  - ❌ Planar die not implemented

- **Rule 902**: Vanguard
  - ❌ Vanguard cards not supported

- **Rule 903**: Commander (EDH)
  - ⚠️ Game type registered: `CommanderFreeForAll`, `CommanderDuel`, `Brawl`
  - ❌ Commander rules not enforced:
    - ❌ Color identity restrictions
    - ❌ Commander damage (21 damage rule)
    - ❌ Command zone casting cost increase
    - ❌ 100-card singleton deck
    - ❌ Commander death replacement (return to command zone)

- **Rule 904**: Archenemy
  - ❌ Scheme cards not supported

- **Rule 905**: Conspiracy Draft
  - ❌ Conspiracy cards not supported

---

## 11. TODOs Found in Engine Code

From `grep TODO` on engine files (47 TODOs found):

### High Priority TODOs

1. **Targeting** (4 TODOs)
   - Hexproof, shroud, protection checking not implemented
   - Target legality incomplete
   - Filter matching stubbed

2. **Combat** (11 TODOs)
   - AsThoughEffectType.ATTACK_AS_HASTE for haste effects
   - AsThoughEffectType.ATTACK for defender attacking
   - Summoning sickness tracking incomplete
   - "Can't attack" restrictions missing
   - Battle type not attackable
   - Suspected status not implemented
   - Shadow, intimidate restrictions missing
   - Protection checking missing
   - Dragon blocking exception missing
   - "Can block while tapped" abilities missing

3. **P/T Calculation** (2 TODOs)
   - Dynamic power calculation stubbed (returns 0)
   - Dynamic toughness calculation stubbed (returns 0)
   - Characteristic-defining abilities (CDA) missing

4. **Layer System** (3 TODOs)
   - More effect type conversions needed
   - Unify permanent and card representations
   - Complex P/T expressions not parsed

5. **Stack Resolution** (4 TODOs)
   - Spell resolution zone changes incomplete
   - Card movement to graveyard/battlefield stubbed
   - Specialized resolution context needed
   - Target information in context missing

6. **Zone Changes** (2 TODOs)
   - Dies events for creatures/planeswalkers not triggered
   - Discard events not triggered

7. **Cost Reductions** (1 TODO)
   - Cost modifications not implemented

---

## Critical Gaps Summary

### Tier 1: Blocks Core Gameplay (Must Fix)

1. **Replacement Effects** (Rule 614)
   - Affects hundreds of cards
   - Prevents many ETB effects, death replacement, doubling effects

2. **Prevention Effects** (Rule 615)
   - Protection doesn't work
   - Damage prevention doesn't work

3. **Copy Effects** (Rule 707)
   - Clone effects broken
   - Token copy effects incomplete
   - Spell copy effects missing

4. **Hexproof/Shroud/Protection** (Rules 702.11, 702.18, 702.16)
   - Targeting restrictions not enforced
   - Combat restrictions not checked

5. **Dynamic P/T** (Rule 613.4)
   - Characteristic-defining abilities not working
   - * / * creatures broken

6. **Indestructible** (Rule 702.12)
   - Not enforced in state-based actions

### Tier 2: Breaks Many Cards (Should Fix Soon)

7. **174 Missing Keywords** (Rule 702)
   - Only 14/188 keywords implemented
   - Kicker, Flashback, Cycling, Prowess, etc. all missing

8. **Alternative Costs** (Rule 601.2f)
   - Flashback, Escape, Madness, etc. can't work

9. **Cost Modifications** (Rule 601.2h)
   - Convoke, Delve, Improvise don't reduce costs

10. **Loyalty Abilities** (Rule 606)
    - Can't activate more than once per turn (not enforced)

11. **Special Card Types**
    - Battles (partial, 347 mentions but TODOs)
    - Sagas (completely missing)
    - Adventures (completely missing)
    - Double-faced cards (7 mentions, likely stubs)

### Tier 3: Breaks Modern Sets (Future Work)

12. **Transform/Meld/Modal Cards** (Rules 709-713)
    - Double-faced cards incomplete
    - Meld not implemented
    - Split cards not implemented

13. **Monarch/Initiative** (Rules 724-725)
    - Not tracked or implemented

14. **Day/Night** (Rule 730)
    - Innistrad werewolves don't work

15. **Token Abilities**
    - Token modification incomplete
    - Populate mechanic missing

### Tier 4: Breaks Multiplayer (Lower Priority)

16. **Multiplayer Rules** (Rules 800-811)
    - 3+ player games won't work correctly
    - Team formats broken

17. **Commander Rules** (Rule 903)
    - Color identity not enforced
    - Commander damage not tracked
    - Command zone tax not implemented

---

## Recommendations

### Phase 1: Fix Tier 1 Gaps (Critical)

1. Implement replacement effect system (Rule 614)
2. Implement prevention effect system (Rule 615)
3. Add hexproof/shroud/protection checking
4. Implement copy effect system (Rule 707)
5. Fix dynamic P/T calculation (CDA abilities)
6. Enforce indestructible in SBAs

### Phase 2: Expand Keyword Coverage (High Priority)

7. Implement top 30 most-used keywords:
   - Protection, Ward, Shroud, Prowess, Toxic
   - Kicker, Flashback, Cycling
   - Convoke, Delve, Improvise (cost reduction)
   - Cascade, Storm (stack mechanics)
   - Exalted, Heroic, Extort (triggered)

### Phase 3: Modern Card Support (Medium Priority)

8. Implement alternative cost system
9. Implement cost modification system
10. Complete loyalty ability system
11. Implement Battles fully
12. Implement Sagas
13. Implement Adventures
14. Complete double-faced card support

### Phase 4: Format Support (Lower Priority)

15. Implement Commander-specific rules
16. Implement multiplayer turn structure
17. Implement Two-Headed Giant
18. Add Monarch/Initiative tracking
19. Add Day/Night tracking

---

## Testing Requirements

For each gap closure:

1. **Unit tests** for the rule system
2. **Integration tests** with real cards
3. **Regression tests** against Java XMage behavior
4. **Comprehensive rules tests** from Comprehensive Rules examples

---

## Conclusion

The Go engine has **excellent coverage of basic MTG gameplay** (combat, stack, turns, SBAs) but is missing:

- **174 of 188 keyword abilities** (92% missing)
- **Replacement/prevention effect systems** (critical)
- **Copy effect system** (critical)
- **Most alternative cost mechanics** (affects hundreds of cards)
- **Modern card types** (Battles partial, Sagas/Adventures missing)
- **Multiplayer support** (completely absent)

**Estimated implementation time**:
- Tier 1 fixes: 6-8 weeks
- Tier 2 fixes: 12-16 weeks
- Tier 3 fixes: 8-12 weeks
- Tier 4 fixes: 6-8 weeks

**Total**: ~32-44 weeks of work to reach comprehensive rules coverage.

The good news: The architecture is solid and extensible. Most gaps are "more of the same" rather than fundamental redesigns.
