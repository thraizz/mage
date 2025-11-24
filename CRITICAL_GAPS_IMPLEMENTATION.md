# Critical Gaps - Implementation Status & Fixes

**Date**: 2025-01-24
**Status**: Systems exist but need integration and bug fixes

## Discovery Summary

After analyzing the codebase against RULES_GAP_ANALYSIS.md, I found that **all Tier 1 critical systems already exist** but have integration gaps:

1. ✅ **Replacement Effects** (Rule 614) - EXISTS in `/internal/game/effects/replacement.go`
2. ✅ **Prevention Effects** (Rule 615) - EXISTS in `/internal/game/effects/prevention.go`
3. ✅ **Copy Effects** (Rule 707) - EXISTS in `/internal/game/effects/copy.go`

## Critical Integration Gaps

### 1. Replacement/Prevention Effects Not Wired to Engine

**Files exist**:
- `effects/replacement.go` - Full implementation with `ReplacementEffectManager`
- `effects/replacement_manager.go` - Manager with tests
- `effects/prevention.go` - Prevention effects with protection support

**Problem**: `grep "ReplacementEffectManager" mage_engine.go` returns **0 matches**

**Impact**:
- ETB replacement effects don't work (enters with counters)
- Death replacement effects don't work (shields, totem armor)
- Doubling effects don't work
- Damage prevention doesn't work
- Protection doesn't work

**Fix Required**: Integrate `ReplacementEffectManager` into `MageEngine` struct

### 2. Targeting Protection Checks Stubbed

**File**: `/internal/game/targeting/validator.go`
**Line 121**: `// TODO: Check for hexproof, protection, shroud, etc.`

**Problem**: ValidateTarget() doesn't check keyword abilities

**Impact**:
- Hexproof creatures can be targeted
- Shroud permanents can be targeted
- Protection doesn't prevent targeting

**Fix Required**: Add keyword ability checking to ValidateTarget()

### 3. Dynamic P/T Returns Zero

**File**: `/internal/game/mage_engine.go`
**Lines 6874, 6893**: `return 0, nil // TODO: Calculate dynamic power/toughness`

**Problem**: Creatures with `*/*` or `X/X` P/T return 0/0

**Impact**:
- Tarmogoyf always 0/0
- Lord of Extinction always 0/0
- Kavu Chameleon always 0/0
- All Characteristic-Defining Abilities broken

**Fix Required**: Implement CDA (Characteristic-Defining Ability) system

### 4. Indestructible Not Enforced

**File**: `/internal/game/rules/state_based_actions.go`
**Problem**: No check for indestructible keyword before destroying creatures

**Impact**: Indestructible creatures die to lethal damage

**Fix Required**: Check for indestructible before applying 704.5g and 704.5h

## Implementation Plan

### Phase 1: Integration (High Priority)

#### Task 1: Integrate Replacement Effects Manager
- [ ] Add `replacementEffects *effects.ReplacementEffectManager` to `MageEngine` struct
- [ ] Add `preventionEffects *effects.PreventionEffectManager` to `MageEngine` struct
- [ ] Initialize managers in `NewMageEngine()`
- [ ] Wire event handling to call `ApplyReplacements()` before executing events
- [ ] Add cleanup hooks for expired effects

#### Task 2: Wire Copy Effects to Layer System
- [ ] Integrate `CopyEffect` with `ContinuousEffectsManager`
- [ ] Ensure copy effects apply in Layer 1
- [ ] Test Clone, Phantasmal Image, Copy Enchantment

### Phase 2: Bug Fixes (High Priority)

#### Task 3: Implement Targeting Protection Checks
- [ ] Add `HasKeywordAbility(cardID, keyword)` method to game state accessor
- [ ] Implement hexproof check in `ValidateTarget()`
- [ ] Implement shroud check in `ValidateTarget()`
- [ ] Implement protection check in `ValidateTarget()`
- [ ] Implement ward check (triggers but doesn't prevent targeting)

#### Task 4: Implement Characteristic-Defining Abilities (CDA)
- [ ] Create CDA system for `*/*` P/T
- [ ] Implement Tarmogoyf CDA (creatures + instants/sorceries in all graveyards)
- [ ] Implement Lord of Extinction CDA (creatures in all graveyards)
- [ ] Implement Kavu Chameleon CDA (other creatures you control)
- [ ] Wire CDA calculations into `GetCreaturePower()` and `GetCreatureToughness()`

#### Task 5: Enforce Indestructible in SBAs
- [ ] Add indestructible check before Rule 704.5g (lethal damage)
- [ ] Add indestructible check before Rule 704.5h (deathtouch damage)
- [ ] Ensure destroy effects check indestructible
- [ ] Test with Darksteel Colossus, Blightsteel Colossus

### Phase 3: Testing (High Priority)

#### Task 6: Test Replacement Effects
- [ ] Test ETB with counters (Renown, Modular)
- [ ] Test doubling effects (Doubling Season)
- [ ] Test death replacement (Totem Armor, Shield counters)
- [ ] Test Panharmonicon (double ETB triggers)

#### Task 7: Test Prevention Effects
- [ ] Test damage prevention spells
- [ ] Test Protection from [color]
- [ ] Test Protection from [type]
- [ ] Test Ward triggers

#### Task 8: Test Copy Effects
- [ ] Test Clone entering as copy
- [ ] Test Phantasmal Image
- [ ] Test Copy Enchantment
- [ ] Test Progenitor Mimic

#### Task 9: Test CDA
- [ ] Test Tarmogoyf P/T calculation
- [ ] Test Lord of Extinction
- [ ] Test * / * creatures with layers
- [ ] Test X/X creatures (Primordial Hydra, etc.)

## Detailed Implementation

### 1. Integrate Replacement Effects into MageEngine

```go
// In mage_engine.go, add to MageEngine struct:
type MageEngine struct {
    // ... existing fields ...

    // Critical systems integration
    replacementEffects *effects.ReplacementEffectManager
    preventionEffects  *effects.PreventionEffectManager
}

// In NewMageEngine():
engine := &MageEngine{
    // ... existing initialization ...
    replacementEffects: effects.NewReplacementEffectManager(),
    preventionEffects:  effects.NewPreventionEffectManager(),
}
```

### 2. Wire Events to Replacement System

```go
// Before any zone change event:
func (e *MageEngine) MoveCard(gameID, cardID string, fromZone, toZone int) error {
    // Create replacement event
    event := &rules.Event{
        Type:     rules.EventZoneChange,
        SourceID: cardID,
        Data: map[string]interface{}{
            "fromZone": fromZone,
            "toZone":   toZone,
        },
    }

    // Apply replacement effects
    result, err := e.replacementEffects.ApplyReplacements(ctx, event)
    if err != nil {
        return err
    }

    if result.Prevented {
        return nil // Event completely replaced
    }

    // Execute modified event
    // ... rest of move logic with result.NewData ...
}
```

### 3. Add Targeting Protection Checks

```go
// In targeting/validator.go, replace TODO at line 121:

// Check for hexproof, protection, shroud, etc.
if err := tv.checkTargetingRestrictions(card, requirement); err != nil {
    return err
}

func (tv *TargetValidator) checkTargetingRestrictions(card TargetCardInfo, req TargetRequirement) error {
    // Check shroud (can't be targeted)
    if tv.cardHasKeyword(card.ID, "SHROUD") {
        return fmt.Errorf("target %s has shroud and can't be targeted", card.Name)
    }

    // Check hexproof (can't be targeted by opponent's spells/abilities)
    if tv.cardHasKeyword(card.ID, "HEXPROOF") {
        if req.ControllerID != card.ControllerID {
            return fmt.Errorf("target %s has hexproof and can't be targeted by opponent", card.Name)
        }
    }

    // Check protection
    if protection := tv.getProtectionQuality(card.ID); protection != "" {
        if tv.sourceMatchesProtection(req.SourceID, protection) {
            return fmt.Errorf("target %s has protection from %s", card.Name, protection)
        }
    }

    return nil
}
```

### 4. Implement CDA System

```go
// In mage_engine.go, replace dynamic P/T TODOs:

func (e *MageEngine) GetCreaturePower(gameID, creatureID string) (int, error) {
    creature, exists := e.getCreatureInternal(gameID, creatureID)
    if !exists {
        return 0, fmt.Errorf("creature %s not found", creatureID)
    }

    // Check for CDA (Characteristic-Defining Ability)
    if creature.Power == "*" || creature.Power == "X" {
        return e.calculateCDAPower(gameID, creatureID, creature)
    }

    // ... rest of existing logic ...
}

func (e *MageEngine) calculateCDAPower(gameID, creatureID string, creature *internalCard) (int, error) {
    // Check for CDA abilities on this creature
    for _, abilityID := range creature.Abilities {
        if cda, isCDA := e.getCharacteristicDefiningAbility(gameID, abilityID); isCDA {
            power, err := cda.CalculatePower(e, gameID, creatureID)
            if err != nil {
                return 0, err
            }
            return power, nil
        }
    }

    // No CDA found, return 0 as default
    return 0, nil
}
```

### 5. Enforce Indestructible in SBAs

```go
// In rules/state_based_actions.go, modify lethal damage check:

func (c *Checker) checkCreatureLethalDamage(game GameState) []Action {
    actions := make([]Action, 0)

    for _, creature := range game.GetBattlefield() {
        // Skip if indestructible (Rule 702.12)
        if creature.HasKeywordAbility("INDESTRUCTIBLE") {
            continue
        }

        // Rule 704.5g: Creature with lethal damage
        if creature.Damage >= creature.Toughness && creature.Toughness > 0 {
            actions = append(actions, DestroyCreatureAction{CreatureID: creature.ID})
        }

        // Rule 704.5h: Creature with deathtouch damage
        if creature.Damage > 0 && creature.DamagedByDeathtouch {
            actions = append(actions, DestroyCreatureAction{CreatureID: creature.ID})
        }
    }

    return actions
}
```

## Testing Strategy

### Unit Tests

Create tests for each fix in isolation:

1. `replacement_integration_test.go` - Test replacement effects in engine
2. `targeting_protection_test.go` - Test hexproof/shroud/protection targeting
3. `cda_calculation_test.go` - Test * / * P/T calculations
4. `indestructible_sba_test.go` - Test indestructible preventing destruction

### Integration Tests

Test complete card interactions:

1. `tarmogoyf_test.go` - Tarmogoyf P/T changes with graveyard
2. `clone_test.go` - Clone copying creatures
3. `doubling_season_test.go` - Doubling Season with +1/+1 counters and tokens
4. `darksteel_colossus_test.go` - Indestructible + lethal damage + deathtouch

### Regression Tests

Ensure existing functionality still works:
- Combat damage
- Stack resolution
- State-based actions
- Turn structure

## Success Criteria

✅ **Replacement Effects**: Doubling Season doubles counters and tokens
✅ **Prevention Effects**: Protection prevents targeting and damage
✅ **Copy Effects**: Clone enters as copy of target creature
✅ **Targeting**: Hexproof/Shroud prevent targeting
✅ **CDA**: Tarmogoyf has correct P/T based on graveyards
✅ **Indestructible**: Darksteel Colossus survives lethal damage

## Estimated Timeline

- **Phase 1** (Integration): 2-3 days
- **Phase 2** (Bug Fixes): 3-4 days
- **Phase 3** (Testing): 2-3 days

**Total**: 7-10 days for all Tier 1 critical gaps

## Files to Modify

1. `/internal/game/mage_engine.go` - Add managers, wire events
2. `/internal/game/targeting/validator.go` - Add protection checks
3. `/internal/game/rules/state_based_actions.go` - Add indestructible check
4. `/internal/game/abilities/characteristic_defining.go` - NEW FILE for CDA
5. Various test files for validation

## Notes

- The heavy lifting is already done (all systems implemented)
- Main work is integration and wiring
- Most changes are additive (low risk of breaking existing functionality)
- Test coverage already exists for individual systems
- Need integration tests to verify end-to-end functionality
