# Priority Windows Implementation

This document describes the implementation of priority windows for casting during stack resolution, per MTG Comprehensive Rules 116, 117, 605, and 608.

## Overview

The priority windows system allows players to take specific actions during spell/ability resolution, including:
- Activating mana abilities during cost payment
- Taking special actions when allowed
- Casting spells during resolution (nested casting)
- Making choices during resolution

## Architecture

### Core Components

#### 1. Resolution Context (`priority.go`)
Tracks which spell/ability is currently resolving and manages nested resolution.

**Key Features:**
- Nested resolution tracking with configurable depth limit (default: 10)
- Permission flags for mana abilities and special actions
- Thread-safe with mutex protection

**API:**
```go
rc := NewResolutionContext()
rc.BeginResolution("spell-id")
rc.SetAllowManaAbilities(true)
// ... resolution logic ...
rc.EndResolution("spell-id")
```

#### 2. Mana Ability Manager (`mana_ability.go`)
Manages mana ability activation and triggered mana abilities.

**Key Features:**
- **Rule 605.3b**: Mana abilities resolve immediately (don't use stack)
- **Rule 605.3c**: Once activated, can't be activated again until resolved
- **Rule 605.4a**: Triggered mana abilities resolve immediately after triggering ability
- **Rule 117.1d**: Can activate during casting or when rule requests mana payment

**API:**
```go
mam := NewManaAbilityManager()
ability := ManaAbility{
    ID: "forest-tap",
    Activate: func() error {
        // Add mana to pool
        return nil
    },
}
mam.ActivateManaAbility(ability)
mam.ResolveTriggeredManaAbilities()
```

#### 3. Special Action Manager (`special_action.go`)
Manages special actions per Rule 116.

**Key Features:**
- 10 special action types (play land, morph, companion, etc.)
- Per-action restrictions (main phase, empty stack, own turn, etc.)
- Once-per-game tracking (e.g., companion)
- **Rule 116.3**: Player receives priority after special action

**API:**
```go
sam := NewSpecialActionManager()
action := SpecialAction{
    Type: SpecialActionPlayLand,
    Execute: func() error {
        // Play the land
        return nil
    },
}
if sam.CanTakeAction(action, hasPriority, isMainPhase, isEmptyStack, isOwnTurn) {
    sam.TakeAction(action)
}
```

#### 4. Payment Window Manager (`payment_window.go`)
Manages cost payment during spell/ability casting.

**Key Features:**
- Three payment steps: BEFORE (special mana like Convoke), NORMAL, AFTER
- Per Java implementation: after special mana payment, normal mana abilities blocked
- Tracks which costs have been paid
- Integrates with mana ability activation

**API:**
```go
pwm := NewPaymentWindowManager()
costs := []Cost{{Type: CostTypeMana, Amount: 3}}
state := NewPaymentState("spell-id", "player-1", costs)
pwm.BeginPayment(state)
// ... activate mana abilities, pay costs ...
pwm.EndPayment("spell-id")
```

#### 5. Choice Manager (`payment_window.go`)
Manages player choices during resolution per Rule 608.2.

**Key Features:**
- 9 choice types (mode, target, X value, color, etc.)
- Min/max choice validation
- **Rule 608.2e**: APNAP order for multiplayer choices

**API:**
```go
cm := NewChoiceManager()
choice := Choice{
    Type: ChoiceTypeMode,
    Prompt: "Choose a mode",
    Options: []string{"draw 2", "deal 2"},
    MinChoices: 1,
    MaxChoices: 1,
}
cm.AddChoice(choice)
next := cm.GetNextChoice()
cm.MakeChoice([]string{"draw 2"})
```

#### 6. Priority Window Manager (`priority.go`)
Coordinates different priority windows during resolution.

**Key Features:**
- Manages active priority window
- Validates allowed actions per window type
- Window history for debugging/replay

**API:**
```go
pwm := NewPriorityWindowManager()
window := PriorityWindow{
    Type: PriorityWindowManaPayment,
    AllowedActions: []ActionType{ActionActivateMana},
}
pwm.OpenWindow(window)
// ... player takes actions ...
pwm.CloseWindow()
```

## MTG Rules Coverage

### Rule 116: Special Actions
- ✅ 116.2a: Play land (main phase, empty stack)
- ✅ 116.2b: Turn face up (anytime with priority)
- ✅ 116.2c: End effect (anytime with priority)
- ✅ 116.2d: Ignore static ability (anytime with priority)
- ✅ 116.2e: Discard Circling Vultures (anytime could cast instant)
- ✅ 116.2f: Suspend (with priority, could begin casting)
- ✅ 116.2g: Companion (main phase, empty stack, once per game)
- ✅ 116.2h: Foretell (own turn with priority)
- ✅ 116.2k: Plot (own turn, empty stack, with priority)
- ✅ 116.2m: Unlock (main phase, empty stack)
- ✅ 116.3: Player receives priority after special action

### Rule 117: Timing and Priority
- ✅ 117.1d: Activate mana abilities when casting/paying costs
- ✅ 117.1c: Take special actions when have priority
- ✅ 117.2e: No player has priority while spell/ability resolving
- ✅ 117.5: State-based actions and triggers before priority

### Rule 605: Mana Abilities
- ✅ 605.3a: Can activate when casting, paying, or rule asks for mana
- ✅ 605.3b: Resolve immediately, don't use stack
- ✅ 605.3c: Can't reactivate until resolved
- ✅ 605.4a: Triggered mana abilities resolve immediately

### Rule 608: Resolving Spells and Abilities
- ✅ 608.2: Resolution steps and choices
- ✅ 608.2e: APNAP order for multiplayer choices
- ✅ 707.12: Cast copies during resolution

## Test Coverage

### Priority Context Tests (11 tests)
- Basic resolution tracking
- Nested resolution (up to 10 levels)
- Maximum depth enforcement
- Mana ability permission toggling
- Special action permission toggling
- Reset functionality

### Mana Ability Tests (9 tests)
- Activation tracking per Rule 605.3c
- Immediate resolution per Rule 605.3b
- Triggered mana abilities per Rule 605.4a
- Cascading triggered abilities
- Permission during cast vs. resolve
- Integration with payment windows
- Error handling

### Special Action Tests (10 tests)
- Restriction validation for all 10 action types
- Main phase + empty stack requirements
- Own turn requirements
- Once-per-game enforcement (companion)
- Custom restriction checks
- Per-turn tracking
- Multi-player tracking

### Payment Window Tests (10 tests)
- Payment state tracking
- Payment step progression (BEFORE → NORMAL → AFTER)
- Partial payment tracking
- Mana ability integration
- Convoke blocking normal mana abilities
- Payment window lifecycle
- Error handling

### Choice Manager Tests (5 tests)
- Choice validation (min/max)
- Sequential choices
- APNAP order documentation
- Error conditions

**Total: 45 comprehensive tests**

## Integration Example

```go
// During spell casting with mana payment
func CastSpell(spellID, playerID string, cost ManaCost) error {
    // 1. Create payment state
    costs := []Cost{{Type: CostTypeMana, Amount: cost.TotalCMC()}}
    paymentState := NewPaymentState(spellID, playerID, costs)

    // 2. Begin payment window
    pwm := NewPaymentWindowManager()
    if err := pwm.BeginPayment(paymentState); err != nil {
        return err
    }

    // 3. Open priority window for mana payment
    priorityWM := NewPriorityWindowManager()
    window := PriorityWindow{
        Type: PriorityWindowManaPayment,
        PlayerID: playerID,
        AllowedActions: []ActionType{ActionActivateMana, ActionSpecialAction},
    }
    priorityWM.OpenWindow(window)

    // 4. Player activates mana abilities
    mam := NewManaAbilityManager()
    mam.SetCanActivateDuringCast(true)

    // Activate Forest for {G}
    forestTap := ManaAbility{
        ID: "forest-tap",
        Activate: func() error {
            // Add {G} to player's mana pool
            return nil
        },
    }
    mam.ActivateManaAbility(forestTap) // Resolves immediately (Rule 605.3b)

    // 5. Pay mana from pool
    paymentState.AddManaPaid(1)
    paymentState.MarkCostPaid(CostTypeMana)

    // 6. Close windows
    priorityWM.CloseWindow()
    pwm.EndPayment(spellID)

    // 7. Put spell on stack
    // ...

    return nil
}
```

## Future Enhancements

While the core infrastructure is complete, future work could include:

1. **Integration with Game Engine**: Wire these managers into the main MageEngine
2. **UI/Prompt System**: Connect choice manager to player interface
3. **Replay Support**: Leverage window history for game replay
4. **Performance Metrics**: Track window opening frequency
5. **Additional Special Actions**: Add new special actions as sets release

## File Structure

```
internal/game/rules/
├── priority.go              # Resolution context and priority windows
├── priority_test.go         # 11 tests
├── mana_ability.go          # Mana ability manager
├── mana_ability_test.go     # 9 tests
├── special_action.go        # Special action manager
├── special_action_test.go   # 10 tests
├── payment_window.go        # Payment and choice managers
└── payment_window_test.go   # 15 tests
```

## References

- MTG Comprehensive Rules (September 2025)
- Java MAGE implementation (`GameImpl.playPriority()`, `HumanPlayer.playManaAbilities()`)
- Rules 116 (Special Actions), 117 (Priority), 605 (Mana Abilities), 608 (Resolution)
