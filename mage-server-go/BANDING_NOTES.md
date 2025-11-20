# Banding Implementation Notes

## Status: Partially Implemented (P3 - Low Priority)

**Implemented:** Damage assignment control (Rules 702.22j-k) - the most gameplay-relevant part of banding.

**Not Implemented:** Band formation, block propagation, "bands with other" variants.

Banding is one of Magic's most complex and rarely-used mechanics. This document outlines what has been implemented and what would be needed for full implementation.

## What's Implemented

### Damage Assignment Control (Rules 702.22j-k)

**✅ Rule 702.22j**: When an attacker is blocked by a creature with banding, the **defending player** assigns the attacker's damage (instead of attacking player).

**✅ Rule 702.22k**: When a blocker is blocking a creature with banding, the **attacking player** assigns the blocker's damage (instead of defending player).

**Implementation:**
- `abilityBanding` constant for detection
- `hasBanding()` helper method
- `defenderControlsDamageAssignment()` checks if defender controls attacker damage assignment
- `attackerControlsDamageAssignment()` checks if attacker controls blocker damage assignment
- `AssignAttackerDamage()` validates correct player based on banding
- `AssignBlockerDamage()` validates correct player based on banding
- `BandedCards []string` field on internalCard (for future use)

**Tests:**
- `TestBanding_DefenderAssignsDamage` - Verifies Rule 702.22j
- `TestBanding_AttackerAssignsDamage` - Verifies Rule 702.22k
- `TestBanding_NormalCombatWithoutBanding` - Ensures normal combat still works
- `TestBanding_OnlyOneBandingNeeded` - Only one creature needs banding

All tests passing!

## Rules Summary

### Core Banding Rules (Rule 702.22)

**702.22a** - Banding is a static ability that modifies combat rules

**702.22c** - Band Formation:
- 1+ creatures with banding + up to 1 without banding can form a band
- All creatures in band must attack same target
- Player declares bands during attack declaration

**702.22d** - All creatures in attacking band must attack same player/planeswalker/battle

**702.22e** - Band persists through combat even if banding is removed

**702.22f** - Creature removed from combat is removed from band

**702.22h** - If one band member blocked, ALL members blocked by same blocker(s)

**702.22i** - If one member would become blocked by effect, entire band blocked

**702.22j** - **DEFENDING PLAYER ASSIGNS DAMAGE** when blocked by banding creature
- Exception to normal Rule 510.1c
- Defending player divides attacker's damage among banding blockers

**702.22k** - **ATTACKING PLAYER ASSIGNS DAMAGE** when blocker blocks banding attacker
- Exception to normal Rule 510.1d
- Attacking player divides blocker's damage among banding attackers

### "Bands with Other" Variant (Rule 702.22b)

Special form that allows banding with creatures sharing a quality:
- By subtype: "Bands with other Dinosaurs"
- By supertype: "Bands with other legendary creatures"
- By name: "Bands with other creatures named Wolves of the Hunt"

Requires at least 2 creatures with the matching quality.

## Java MAGE Implementation Analysis

### Data Structures

**Permanent.java** stores bands:
```java
protected List<UUID> bandedCards = new ArrayList<>();
```

Bidirectional relationship: If A bands with B, then:
- A.bandedCards contains B's UUID
- B.bandedCards contains A's UUID

### Key Methods

1. **Combat.handleBanding()** (lines 338-469)
   - Called during attacker declaration
   - Prompts player to select creatures to band with
   - Validates band formation rules
   - Creates bidirectional links

2. **CombatGroup.setBlocked()** (lines 764-779)
   - Propagates blocked status across all band members
   - Implements Rule 702.22h

3. **CombatGroup.defenderAssignsCombatDamage()** (lines 850-868)
   - Returns true if defender controls damage assignment
   - Checks for banding blockers (Rule 702.22j)

4. **CombatGroup.attackerAssignsCombatDamage()** (lines 831-842)
   - Returns true if attacker controls damage assignment
   - Checks for banding attackers (Rule 702.22k)

5. **CombatGroup.appliesBandsWithOther()** (lines 116-164)
   - Validates "bands with other" combinations
   - Checks subtype/supertype/name matching

## Implementation Requirements for Go Port

### Phase 1: Basic Infrastructure (Required)

1. **Add ability constant:**
```go
abilityBanding = "BandingAbility"
```

2. **Add band tracking to internalCard:**
```go
type internalCard struct {
    // ... existing fields
    BandedCards []string // IDs of creatures in same band
}
```

3. **Add helper method:**
```go
func (e *MageEngine) hasBanding(card *internalCard) bool {
    return e.hasAbility(card, abilityBanding)
}
```

### Phase 2: Band Formation (Complex)

1. **API for band declaration:**
```go
// FormAttackingBand creates a band of attacking creatures
// attackerIDs must all attack same target, follow banding rules
func (e *MageEngine) FormAttackingBand(gameID string, attackerIDs []string, playerID string) error
```

2. **Validation logic:**
   - All creatures attack same target
   - All have banding OR at most 1 without
   - Not already in another band
   - All controlled by same player

3. **Bidirectional linking:**
   - Each creature's BandedCards contains all others
   - Must maintain consistency when adding/removing

### Phase 3: Block Propagation (Medium)

Update `DeclareBlocker()` to propagate blocking across band:

```go
// When blocking any band member, all members become blocked
if len(attacker.BandedCards) > 0 {
    for _, bandedID := range attacker.BandedCards {
        // Mark all banded creatures as blocked
    }
}
```

### Phase 4: Damage Assignment Control (High Impact)

This is the most gameplay-relevant part.

1. **Add control checking methods:**
```go
func (e *MageEngine) defenderControlsDamageAssignment(group *combatGroup) bool {
    // Rule 702.22j: Check if any blocker has banding
    for _, blockerID := range group.blockers {
        if blocker, exists := gameState.cards[blockerID]; exists {
            if e.hasBanding(blocker) {
                return true
            }
        }
    }
    return false
}

func (e *MageEngine) attackerControlsDamageAssignment(group *combatGroup) bool {
    // Rule 702.22k: Check if any attacker has banding
    for _, attackerID := range group.attackers {
        if attacker, exists := gameState.cards[attackerID]; exists {
            if e.hasBanding(attacker) {
                return true
            }
        }
    }
    return false
}
```

2. **Modify damage assignment:**

In `AssignAttackerDamage()`:
```go
// Check if defending player controls assignment
if e.defenderControlsDamageAssignment(group) {
    // Validate playerID is defending player, not attacking player
    if playerID != group.defendingPlayerID {
        return fmt.Errorf("defending player must assign damage (banding)")
    }
}
```

In `AssignBlockerDamage()`:
```go
// Check if attacking player controls assignment
if e.attackerControlsDamageAssignment(group) {
    // Validate playerID is attacking player
    if playerID != gameState.combat.attackingPlayerID {
        return fmt.Errorf("attacking player must assign damage (banding)")
    }
}
```

### Phase 5: "Bands with Other" (Very Complex)

Requires:
1. New ability type with quality parameter (subtype/supertype/name)
2. Special validation in band formation
3. Minimum 2 creatures with matching quality
4. Complex filtering logic

**Recommendation:** Skip until basic banding works

### Phase 6: Edge Cases

1. **Removal during combat:**
   - Remove from band when removed from combat
   - Update all bandedCards lists

2. **Banding granted/lost mid-combat:**
   - Rule 702.22e: Band persists even if banding removed
   - Track "formed as band" separately from "has banding"

3. **Multiple bands:**
   - Each creature can only be in one band
   - Need to track which creatures are already banded

## Testing Requirements

### Unit Tests Needed

1. **Damage assignment with banding blocker:**
   - Attacker blocked by creature with banding
   - Defending player assigns damage (not attacking player)
   - Verify damage division works correctly

2. **Damage assignment with banding attacker:**
   - Banding attacker blocked by multiple creatures
   - Attacking player assigns blocker's damage
   - Verify control switches correctly

3. **Band formation validation:**
   - Accept: all banding + up to 1 non-banding
   - Reject: 2+ non-banding creatures
   - Reject: creatures attacking different targets
   - Reject: creature already in another band

4. **Block propagation:**
   - When one band member blocked, all blocked
   - Blockers assigned to all band members

### Integration Tests Needed

1. Full combat with banding attackers
2. Full combat with banding blockers
3. Band formation -> blocking -> damage assignment flow
4. Removal of band member during combat

## Estimated Complexity

- **Basic infrastructure:** 2-3 hours
- **Band formation:** 8-10 hours (includes UI/API design)
- **Block propagation:** 3-4 hours
- **Damage assignment control:** 4-5 hours
- **"Bands with other":** 10-15 hours
- **Comprehensive testing:** 8-10 hours
- **Edge cases:** 5-8 hours

**Total: 40-55 hours of development time**

## Recommendation

Given that banding is:
- **P3 - Low Priority** in the task list
- Rarely used in modern Magic (appears in ~50 cards total, mostly from early sets)
- Extremely complex to implement fully
- Requires extensive UI work for band formation

**Recommended approach:**
1. ✅ Document requirements (this file)
2. ✅ Note as "not implemented" in GO_PORT_TASKS.md
3. ⏳ Implement only if specifically needed for a card
4. ⏳ Start with Phase 4 (damage assignment control) as it's highest impact
5. ⏳ Skip full band formation until there's a clear use case

## Alternative: Simplified Implementation

If basic banding support is needed:

### Minimum Viable Implementation

1. **Add banding ability detection** ✓ (trivial)
2. **Implement Rule 702.22j** (defending player assigns damage)
3. **Implement Rule 702.22k** (attacking player assigns damage)
4. **Skip band formation entirely** (require manual setup in tests)
5. **Skip "bands with other"** (too complex, rarely used)

This gives the gameplay impact (different damage assignment) without the complexity of band formation.

**Estimated time: 6-8 hours** (vs 40-55 for full implementation)

## Cards Affected

~50 cards in total, mostly from Alpha/Beta/Unlimited/Arabian Nights:

**Common banding cards:**
- Benalish Hero
- Icatian Phalanx
- Mesa Pegasus
- Kjeldoran Warrior

**"Bands with other" cards:**
- Wolves of the Hunt (token)
- Old Fogey (Un-set, "bands with other Dinosaurs")

**Grants banding:**
- Banding Sliver
- Tolaria
- Kjeldoran Outpost

## See Also

- `internal/game/mage_engine.go` - damage assignment logic
- `internal/game/combat_damage_division_test.go` - existing damage division tests
- Java MAGE: `mage/abilities/keyword/BandingAbility.java`
- Java MAGE: `mage/game/combat/Combat.java` (handleBanding method)
