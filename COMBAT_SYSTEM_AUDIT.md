# Combat System Audit - Java vs Go Implementation

## Executive Summary

**Status: ❌ NOT IMPLEMENTED**

The Go implementation has **NO combat system**. Only basic view structures exist for displaying combat state, but there is no actual combat logic, attacker/blocker declaration, damage assignment, or combat phases.

**Coverage: ~0%**

## Java Implementation Analysis

### Combat.java (1,956 lines)

The Java `Combat` class is a comprehensive combat management system with:

#### Core Data Structures
- `List<CombatGroup> groups` - Active combat groups
- `List<CombatGroup> formerGroups` - Historical combat groups
- `Map<UUID, CombatGroup> blockingGroups` - Blocker to group mapping
- `Set<UUID> defenders` - All possible defenders (players, planeswalkers, battles)
- `Map<UUID, Set<UUID>> numberCreaturesDefenderAttackedBy` - Attack tracking
- `UUID attackingPlayerId` - Current attacker
- `Map<UUID, Set<UUID>> creatureMustBlockAttackers` - Block requirements
- `Map<UUID, Set<UUID>> creaturesForcedToAttack` - Attack requirements
- `int maxAttackers` - Maximum attackers limit
- `HashSet<UUID> attackersTappedByAttack` - Tap tracking
- `boolean useToughnessForDamage` - Damage calculation mode
- `List<FilterCreaturePermanent> useToughnessForDamageFilters` - Damage filters

#### Key Methods (50+ methods)

**Setup & State Management:**
- `reset(Game game)` - Reset combat state
- `clear()` - Clear all combat data
- `setAttacker(UUID playerId)` - Set attacking player
- `setDefenders(Game game)` - Identify all possible defenders
- `checkForRemoveFromCombat(Game game)` - Remove non-creatures

**Attacker Declaration:**
- `selectAttackers(Game game)` - Main attacker selection flow
- `resumeSelectAttackers(Game game)` - Resume after interruption
- `declareAttacker(UUID creatureId, UUID defenderId, UUID playerId, Game game)` - Declare single attacker
- `addAttackingCreature(UUID creatureId, Game game)` - Add attacker (e.g., tokens)
- `addAttackerToCombat(UUID attackerId, UUID defenderId, Game game)` - Internal add
- `canDefenderBeAttacked(UUID attackerId, UUID defenderId, Game game)` - Validation
- `removeAttacker(UUID attackerId, Game game)` - Remove attacker

**Blocker Declaration:**
- `selectBlockers(Game game)` - Main blocker selection flow
- `selectBlockers(Player blockController, Ability source, Game game)` - With controller override
- `resumeSelectBlockers(Game game)` - Resume after interruption
- `addBlockingGroup(UUID blockerId, UUID attackerId, UUID playerId, Game game)` - Declare blocker
- `acceptBlockers(Game game)` - Finalize blockers
- `checkBlockRestrictions(Player defender, Game game)` - Validate blocks (before)
- `checkBlockRequirementsAfter(Player player, Player controller, Game game)` - Validate requirements (after)
- `checkBlockRestrictionsAfter(Player player, Player controller, Game game)` - Validate restrictions (after)
- `removeBlocker(UUID blockerId, Game game)` - Remove blocker
- `removeBlockerGromGroup(UUID blockerId, CombatGroup groupToUnblock, Game game)` - Remove from specific group

**Combat Resolution:**
- `endCombat(Game game)` - End combat phase
- `hasFirstOrDoubleStrike(Game game)` - Check for first/double strike
- `removeFromCombat(UUID creatureId, Game game, boolean withEvent)` - Remove creature
- `removeDefendingPermanentFromCombat(UUID permanentId, Game game)` - Remove defender

**Queries:**
- `getAttackers()` - All attackers
- `getBlockers()` - All blockers
- `getDefenders()` - All defenders
- `getGroups()` - All combat groups
- `getBlockingGroups()` - All blocking groups
- `findGroup(UUID attackerId)` - Find group by attacker
- `findGroupOfBlocker(UUID blockerId)` - Find group by blocker
- `getDefenderId(UUID attackerId)` - Get defender for attacker
- `getDefendingPlayerId(UUID attackingCreatureId, Game game)` - Get defending player
- `getPlayerDefenders(Game game)` - Get player defenders only
- `attacksAlone()` - Check single attacker
- `noAttackers()` - Check no attackers
- `isPlaneswalkerAttacked(UUID defenderId, Game game)` - Check planeswalker attacked

**Special Mechanics:**
- `useToughnessForDamage(Permanent permanent, Game game)` - Check toughness damage
- `setUseToughnessForDamage(boolean useToughnessForDamage)` - Set toughness mode
- `addUseToughnessForDamageFilter(FilterCreaturePermanent filter)` - Add filter
- Banding support (multiple methods)
- Forced attack/block support
- Attack/block restrictions and requirements

### CombatGroup.java (897 lines)

Represents a single combat group (1+ attackers vs 1 defender + blockers):

#### Core Data
- `UUID defenderId` - Who/what is being attacked
- `boolean defenderIsPermanent` - Is defender a permanent (planeswalker/battle)
- `UUID defendingPlayerId` - Controlling player
- `List<UUID> attackers` - Attacking creatures
- `List<UUID> formerAttackers` - Historical attackers
- `List<UUID> blockers` - Blocking creatures
- `List<UUID> unblockedAttackers` - Unblocked attackers
- `boolean blocked` - Is group blocked
- `Map<UUID, Integer> attackerOrder` - Damage assignment order for attackers
- `Map<UUID, Integer> blockerOrder` - Damage assignment order for blockers
- `boolean canDamageDefenderDirectly` - Trample/etc.
- `Map<UUID, UUID> attackerBlockedBy` - Attacker to blocker mapping
- `Map<UUID, Set<UUID>> blockedAttackers` - Blocker to attackers mapping

#### Key Methods (28+ methods)

**Damage Assignment:**
- `assignDamageToBlockers(boolean first, Game game)` - Assign attacker damage
- `assignDamageToAttackers(boolean first, Game game)` - Assign blocker damage
- `applyDamage(Game game)` - Apply all damage
- `attackerAssignsCombatDamage(Game game)` - Check attacker assigns
- `defenderAssignsCombatDamage(Game game)` - Check defender assigns
- `assignsDefendingPlayerAndOrDefendingCreaturesDividedDamage(...)` - Divided damage

**State Management:**
- `addBlocker(UUID blockerId, UUID playerId, Game game)` - Add blocker
- `addBlockerToGroup(UUID blockerId, UUID playerId, Game game)` - Add to group
- `remove(UUID creatureId)` - Remove creature
- `removeAttackedPermanent(UUID permanentId)` - Remove defender
- `acceptBlockers(Game game)` - Finalize blockers
- `setBlocked(boolean blocked, Game game)` - Set blocked status
- `changeDefenderPostDeclaration(UUID newDefenderId, Game game)` - Change defender

**Queries:**
- `hasFirstOrDoubleStrike(Game game)` - Check first/double strike
- `getDefenderId()` - Get defender
- `getDefendingPlayerId()` - Get defending player
- `getAttackers()` - Get attackers
- `getFormerAttackers()` - Get former attackers
- `getBlockers()` - Get blockers
- `getBlocked()` - Is blocked
- `canBlock(Permanent blocker, Game game)` - Can blocker block
- `checkSoleBlockerAfter(Permanent blocker, Game game)` - Check sole blocker
- `checkBlockRestrictions(Game game, Player defender, int blockersCount)` - Validate blocks

**Utilities:**
- `dealsDamageThisStep(Permanent perm, boolean first, Game game)` - Static damage check
- `copy()` - Deep copy
- `toString()` - Debug string

## Go Implementation Analysis

### Current State

**EngineCombatView struct (4 lines):**
```go
type EngineCombatView struct {
    AttackingPlayerID string
    Groups            []EngineCombatGroupView
}
```

**EngineCombatGroupView struct (5 lines):**
```go
type EngineCombatGroupView struct {
    Attackers         []string
    Blockers          []string
    DefendingPlayerID string
}
```

**Total Implementation: 9 lines (view structures only)**

### What's Missing

**Everything:**
- ❌ No combat state management
- ❌ No attacker declaration
- ❌ No blocker declaration
- ❌ No damage assignment
- ❌ No damage application
- ❌ No combat groups
- ❌ No attack/block restrictions
- ❌ No attack/block requirements
- ❌ No first strike / double strike
- ❌ No trample
- ❌ No vigilance
- ❌ No banding
- ❌ No forced attack/block
- ❌ No combat damage prevention
- ❌ No combat triggers
- ❌ No combat events
- ❌ No defender selection
- ❌ No blocker ordering
- ❌ No damage ordering
- ❌ No combat cleanup

### internalCard Fields

The `internalCard` struct has a `Tapped` field but no:
- `Attacking` field
- `Blocking` field
- `AttackingWhat` field
- `BlockingWhat` field
- `BlockedBy` field
- `Blocking` list field

## Combat Flow in Java

### 1. Beginning of Combat Step
```
Combat.reset(game)
Combat.setAttacker(activePlayerId)
Combat.setDefenders(game)
Fire COMBAT_PHASE_PRE event
```

### 2. Declare Attackers Step
```
Combat.selectAttackers(game)
  → Player chooses attackers
  → For each attacker:
      Combat.declareAttacker(creatureId, defenderId, playerId, game)
        → Validate can attack
        → Validate can attack defender
        → Tap creature (unless vigilance)
        → Create/add to CombatGroup
        → Fire ATTACKER_DECLARED event
  → Fire DECLARED_ATTACKERS event
```

### 3. Declare Blockers Step
```
Combat.selectBlockers(game)
  → For each defending player:
      Player chooses blockers
      For each blocker:
        Combat.addBlockingGroup(blockerId, attackerId, playerId, game)
          → Validate can block
          → Add to CombatGroup
          → Fire BLOCKER_DECLARED event
  → Combat.acceptBlockers(game)
      → Check block requirements
      → Check block restrictions
      → Set blocked status
      → Order blockers
      → Fire DECLARED_BLOCKERS event
```

### 4. Combat Damage Step (First Strike)
```
If Combat.hasFirstOrDoubleStrike(game):
  For each CombatGroup:
    group.assignDamageToBlockers(true, game)
    group.assignDamageToAttackers(true, game)
    group.applyDamage(game)
  Fire COMBAT_DAMAGE_APPLIED event
  Check state-based actions
```

### 5. Combat Damage Step (Normal)
```
For each CombatGroup:
  group.assignDamageToBlockers(false, game)
  group.assignDamageToAttackers(false, game)
  group.applyDamage(game)
Fire COMBAT_DAMAGE_APPLIED event
Check state-based actions
```

### 6. End of Combat Step
```
Combat.endCombat(game)
  → Move groups to formerGroups
  → Clear current groups
  → Clear blockers
  → Keep attackers for "attacked this turn" tracking
Fire END_COMBAT_STEP_PRE event
```

## Required Implementation

### Phase 1: Core Combat Structure

**1. Combat State (internal struct):**
```go
type combatState struct {
    attackingPlayerID string
    groups            []*combatGroup
    formerGroups      []*combatGroup
    blockingGroups    map[string]*combatGroup // blockerID -> group
    defenders         map[string]bool
    attackers         map[string]bool
    blockers          map[string]bool
    attackersTapped   map[string]bool
}
```

**2. Combat Group (internal struct):**
```go
type combatGroup struct {
    defenderID         string
    defenderIsPermanent bool
    defendingPlayerID  string
    attackers          []string
    formerAttackers    []string
    blockers           []string
    blocked            bool
    attackerOrder      map[string]int
    blockerOrder       map[string]int
}
```

**3. Add to internalCard:**
```go
type internalCard struct {
    // ... existing fields ...
    Attacking    bool
    Blocking     bool
    AttackingWhat string // defenderID
    BlockingWhat  []string // attackerIDs
}
```

**4. Add to engineGameState:**
```go
type engineGameState struct {
    // ... existing fields ...
    combat *combatState
}
```

### Phase 2: Core Combat Methods

**Engine Methods:**
```go
func (e *MageEngine) ResetCombat(gameID string) error
func (e *MageEngine) SetAttacker(gameID, playerID string) error
func (e *MageEngine) SetDefenders(gameID string) error
func (e *MageEngine) DeclareAttacker(gameID, creatureID, defenderID, playerID string) error
func (e *MageEngine) DeclareBlocker(gameID, blockerID, attackerID, playerID string) error
func (e *MageEngine) AcceptBlockers(gameID string) error
func (e *MageEngine) AssignCombatDamage(gameID string, firstStrike bool) error
func (e *MageEngine) ApplyCombatDamage(gameID string) error
func (e *MageEngine) EndCombat(gameID string) error
```

**Combat Group Methods:**
```go
func (cg *combatGroup) addAttacker(attackerID string)
func (cg *combatGroup) addBlocker(blockerID, playerID string)
func (cg *combatGroup) setBlocked(blocked bool)
func (cg *combatGroup) assignDamageToBlockers(firstStrike bool, gameState *engineGameState)
func (cg *combatGroup) assignDamageToAttackers(firstStrike bool, gameState *engineGameState)
func (cg *combatGroup) applyDamage(gameState *engineGameState)
```

### Phase 3: Combat Events

**New Events:**
```go
EventBeginCombat
EventDeclareAttackersStepPre
EventAttackerDeclared
EventDeclaredAttackers
EventDeclareBlockersStepPre
EventBlockerDeclared
EventDeclaredBlockers
EventCombatDamageStepPre
EventCombatDamageAssigned
EventCombatDamageApplied
EventEndCombatStepPre
EventEndCombat
```

### Phase 4: Combat Validation

**Validation Methods:**
```go
func (e *MageEngine) CanAttack(gameID, creatureID string) (bool, error)
func (e *MageEngine) CanAttackDefender(gameID, creatureID, defenderID string) (bool, error)
func (e *MageEngine) CanBlock(gameID, blockerID, attackerID string) (bool, error)
func (e *MageEngine) CheckBlockRequirements(gameID, playerID string) error
func (e *MageEngine) CheckBlockRestrictions(gameID, playerID string) error
```

### Phase 5: Special Combat Mechanics

**Advanced Features:**
- First strike / Double strike damage
- Trample damage
- Vigilance (no tap on attack)
- Flying / Reach restrictions
- Menace / Deathtouch
- Lifelink
- Banding (complex)
- Forced attack/block
- Can't attack/block restrictions

### Phase 6: Damage Assignment

**Player Choices:**
- Order blockers for each attacker
- Order attackers for each blocker
- Assign damage amounts
- Trample overflow to defender

### Phase 7: Combat Integration

**Turn Structure Integration:**
```go
// In TurnManager or game loop:
case StepBeginCombat:
    engine.ResetCombat(gameID)
    engine.SetAttacker(gameID, activePlayerID)
    engine.SetDefenders(gameID)

case StepDeclareAttackers:
    // Wait for player to declare attackers
    // engine.DeclareAttacker() called per attacker

case StepDeclareBlockers:
    // Wait for defenders to declare blockers
    // engine.DeclareBlocker() called per blocker
    // engine.AcceptBlockers() when done

case StepFirstStrikeDamage:
    if engine.HasFirstOrDoubleStrike(gameID) {
        engine.AssignCombatDamage(gameID, true)
        engine.ApplyCombatDamage(gameID)
    }

case StepCombatDamage:
    engine.AssignCombatDamage(gameID, false)
    engine.ApplyCombatDamage(gameID)

case StepEndCombat:
    engine.EndCombat(gameID)
```

## Complexity Analysis

### Java Combat System
- **Lines of Code:** ~2,850 (Combat.java + CombatGroup.java)
- **Methods:** ~78 methods
- **Complexity:** Very High
- **Features:** Complete MTG combat rules
- **Special Mechanics:** Banding, forced attack/block, restrictions, requirements

### Estimated Go Implementation
- **Estimated Lines:** ~2,500-3,000 lines
- **Estimated Methods:** ~60-70 methods
- **Estimated Time:** 2-3 weeks for full implementation
- **Testing Required:** Extensive (100+ test cases)

## Priority Breakdown

### P0 - Critical (Must Have)
1. Combat state management
2. Attacker declaration
3. Blocker declaration
4. Damage assignment (basic)
5. Damage application
6. Combat cleanup

### P1 - High (Should Have)
1. First strike / Double strike
2. Trample
3. Vigilance
4. Flying / Reach
5. Combat events
6. Combat triggers

### P2 - Medium (Nice to Have)
1. Menace
2. Deathtouch
3. Lifelink
4. Block ordering
5. Damage ordering
6. Planeswalker/Battle attacks

### P3 - Low (Future)
1. Banding
2. Forced attack/block
3. Complex restrictions
4. Complex requirements

## Recommendations

### Option 1: Full Implementation (Recommended)
Implement complete combat system matching Java:
- **Pros:** Feature complete, production ready
- **Cons:** Large effort (2-3 weeks)
- **When:** Before beta/production release

### Option 2: Basic Implementation
Implement P0 + P1 features only:
- **Pros:** Covers 90% of use cases
- **Cons:** Missing edge cases
- **When:** For alpha testing

### Option 3: Defer
Continue without combat:
- **Pros:** Focus on other systems
- **Cons:** Can't play real games
- **When:** Early development only

## Test Coverage Requirements

### Minimum Tests (P0)
1. Declare single attacker
2. Declare multiple attackers
3. Declare single blocker
4. Declare multiple blockers
5. Unblocked damage to player
6. Blocked damage to creatures
7. Creature death from combat
8. Player death from combat
9. Tapping on attack
10. Combat cleanup

### Comprehensive Tests (P0-P2)
1. All minimum tests
2. First strike damage
3. Double strike damage
4. Trample damage
5. Vigilance (no tap)
6. Flying vs non-flying
7. Reach blocking flying
8. Multiple blockers ordering
9. Multiple attackers ordering
10. Damage assignment validation
11. Combat events firing
12. Combat triggers
13. Removal during combat
14. Blink during combat
15. Phase out during combat
16. ... (50+ more scenarios)

## Summary

**Current Status:** ❌ **0% Complete**

The Go implementation has **no combat system**. Only view structures exist for displaying combat state to clients. The entire combat logic, from attacker/blocker declaration through damage assignment and application, is missing.

**Estimated Effort:** 2-3 weeks for full implementation

**Blocking:** Yes - combat is essential for playing Magic: The Gathering

**Recommendation:** Implement at least P0 + P1 features before any production use.
