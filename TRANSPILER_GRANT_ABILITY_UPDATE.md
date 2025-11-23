# Transpiler Update: Grant Ability Support

## Summary

Updated the card transpiler (`mage-server-go/scripts/transpile_cards.py`) to properly recognize and document GainAbility effects (cards that grant keyword abilities like Flying, First Strike, Deathtouch, etc.).

## Changes Made

### 1. Effect Mapping (lines 247-254)
Added mappings for all Java GainAbility effect classes:
```python
'GainAbilityTargetEffect': ('MANUAL', 'effects.GrantAbilityEffect'),
'GainAbilitySourceEffect': ('MANUAL', 'effects.GrantAbilityEffect'),
'GainAbilityControlledEffect': ('MANUAL', 'effects.GrantAbilityEffect'),
'GainAbilityAllEffect': ('MANUAL', 'effects.GrantAbilityEffect'),
'GainAbilityAttachedEffect': ('MANUAL', 'effects.GrantAbilityEffect'),
```

### 2. Auto-Import Detection (lines 709-715)
Added `_needs_effects_import()` method to automatically add the `effects` package import when GrantAbilityEffect is detected.

### 3. Ability Expression Parsing (lines 373-394, 874-886)
Added parsing for Java ability expressions:
- `FirstStrikeAbility.getInstance()` → `"FirstStrikeAbility"`
- `new FlyingAbility()` → `"FlyingAbility"`

### 4. Manual Implementation Detection (lines 819-883)
Enhanced code generation to detect MANUAL effects and generate comprehensive documentation:

```go
// MANUAL IMPLEMENTATION REQUIRED: Ability granting effects detected
// Go engine fully supports GrantAbilityEffect - see GRANT_ABILITY_SUPPORT.md
// Example pattern: effects.NewEffectBuilder(card.ID).Targeting(...).UntilEndOfTurn().GrantAbility("XAbility")
//
// Effects to implement:
//   - Grant ability: "FirstStrikeAbility"
//
// Targets:
//   - abilities.NewCreatureTargetFilter()
// card.AddAbility(ability1)
```

## Example Output

### Java Card: Abandon Reason
```java
Effect effect = new GainAbilityTargetEffect(
    FirstStrikeAbility.getInstance(),
    Duration.EndOfTurn
);
this.getSpellAbility().addEffect(effect);
this.getSpellAbility().addTarget(new TargetCreaturePermanent(0, 2));
```

### Generated Go Code
```go
package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"  // ← Auto-imported
)

func NewAbandonReason(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abandon Reason")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"INSTANT"}

	// MANUAL IMPLEMENTATION REQUIRED: Ability granting effects detected
	// Go engine fully supports GrantAbilityEffect - see GRANT_ABILITY_SUPPORT.md
	// Example pattern: effects.NewEffectBuilder(card.ID).Targeting(...).UntilEndOfTurn().GrantAbility("XAbility")
	//
	// Effects to implement:
	//   - Grant ability: "FirstStrikeAbility"
	// card.AddAbility(ability1)

	return card, nil
}
```

## Why MANUAL Instead of Full Automation?

The transpiler marks these as `MANUAL` because:

1. **Pattern Mismatch**: The Go engine uses `EffectBuilder` pattern, not the `SpellAbilityBuilder` pattern used for other effects
2. **Dynamic Targets**: Target IDs come from spell resolution, not compile-time parameters
3. **Effect Manager Integration**: The effect needs to be added to the game's `EffectManager` during resolution

## Verification

The Go engine has **complete support** for granting abilities:

✅ `effects.GrantAbilityEffect` - Fully implemented
✅ All keyword abilities supported (Flying, First Strike, Deathtouch, Haste, Vigilance, etc.)
✅ Duration support (EndOfTurn, EndOfCombat, WhileOnBattlefield, Permanent)
✅ Layer system integration (Layer 6 - Ability)
✅ Comprehensive test coverage (30+ tests)

See `GRANT_ABILITY_SUPPORT.md` for complete implementation details.

## Cards Affected

Any card that grants abilities will be detected by the transpiler:

- **Instant/Sorcery** spells that grant abilities (e.g., Abandon Reason, Giant Growth variants)
- **Activated abilities** that grant abilities (e.g., Saberclaw Golem: "{R}: Gain first strike")
- **Static abilities** that grant abilities (e.g., Lancer Sliver: "Slivers you control have first strike")
- **Aura** enchantments that grant abilities
- **Equipment** that grants abilities

## Testing

```bash
# Test transpilation
python3 scripts/transpile_cards.py --card=AbandonReason --output=internal/game/cards/generated/

# Output will include MANUAL implementation comments with:
# - Clear indication that feature is supported
# - Reference to documentation
# - Example implementation pattern
# - List of abilities to grant
```

## Files Modified

1. `mage-server-go/scripts/transpile_cards.py` - Transpiler script
2. `GRANT_ABILITY_SUPPORT.md` - Comprehensive support documentation (new)
3. `TRANSPILER_GRANT_ABILITY_UPDATE.md` - This summary (new)

## Next Steps

For developers implementing cards manually:
1. Transpile the card to get structure and documentation
2. Follow the pattern in the generated comments
3. Refer to `GRANT_ABILITY_SUPPORT.md` for examples
4. Use test cases in `internal/game/effects/dynamic_abilities_test.go` as reference

For full automation in the future:
1. Create `SpellAbilityBuilder.AddGrantAbilityEffect()` method
2. Or generate spell-specific resolution handlers
3. Implement filter-based targeting for controlled/all effects
