package abilities_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManaCost(t *testing.T) {
	t.Run("Parse simple mana cost", func(t *testing.T) {
		cost, err := abilities.ParseManaCost("{R}")
		require.NoError(t, err)
		assert.Equal(t, 1, cost.Mana.Red)
		assert.Equal(t, "{R}", cost.String())
	})

	t.Run("Parse generic mana cost", func(t *testing.T) {
		cost, err := abilities.ParseManaCost("{2}")
		require.NoError(t, err)
		assert.Equal(t, 2, cost.Mana.Generic)
		assert.Equal(t, "{2}", cost.String())
	})

	t.Run("Parse complex mana cost", func(t *testing.T) {
		cost, err := abilities.ParseManaCost("{2}{U}{U}")
		require.NoError(t, err)
		assert.Equal(t, 2, cost.Mana.Generic)
		assert.Equal(t, 2, cost.Mana.Blue)
		assert.Equal(t, "{2}{U}{U}", cost.String())
	})

	t.Run("Parse mixed colors", func(t *testing.T) {
		cost, err := abilities.ParseManaCost("{W}{U}{B}{R}{G}")
		require.NoError(t, err)
		assert.Equal(t, 1, cost.Mana.White)
		assert.Equal(t, 1, cost.Mana.Blue)
		assert.Equal(t, 1, cost.Mana.Black)
		assert.Equal(t, 1, cost.Mana.Red)
		assert.Equal(t, 1, cost.Mana.Green)
	})

	t.Run("Parse colorless mana", func(t *testing.T) {
		cost, err := abilities.ParseManaCost("{C}")
		require.NoError(t, err)
		assert.Equal(t, 1, cost.Mana.Colorless)
		assert.Equal(t, "{C}", cost.String())
	})
}

func TestEffects(t *testing.T) {
	t.Run("DamageEffect description", func(t *testing.T) {
		effect := abilities.NewDamageEffect(3)
		assert.Equal(t, "deals 3 damage", effect.GetDescription())
	})

	t.Run("DrawCardsEffect description", func(t *testing.T) {
		effect := abilities.NewDrawCardsEffect(2)
		assert.Equal(t, "draw 2 cards", effect.GetDescription())

		effect = abilities.NewDrawCardsEffect(1)
		assert.Equal(t, "draw a card", effect.GetDescription())
	})

	t.Run("DestroyEffect description", func(t *testing.T) {
		effect := abilities.NewDestroyEffect()
		assert.Equal(t, "destroy target", effect.GetDescription())

		effect = abilities.NewDestroyEffectNoRegen()
		assert.Equal(t, "destroy target. It can't be regenerated", effect.GetDescription())
	})

	t.Run("BoostEffect description", func(t *testing.T) {
		effect := abilities.NewBoostEffect(3, 3)
		assert.Equal(t, "gets +3/+3", effect.GetDescription())

		effect = abilities.NewBoostEffect(-1, -1)
		assert.Equal(t, "gets -1/-1", effect.GetDescription())
	})
}

func TestTargetFilters(t *testing.T) {
	t.Run("AnyTargetFilter description", func(t *testing.T) {
		filter := abilities.NewAnyTargetFilter()
		assert.Equal(t, "any target", filter.GetDescription())
	})

	t.Run("CreatureTargetFilter description", func(t *testing.T) {
		filter := abilities.NewCreatureTargetFilter()
		assert.Equal(t, "target creature", filter.GetDescription())

		filter = abilities.NewCreatureTargetFilterWithSubtype("Human")
		assert.Equal(t, "target Human creature", filter.GetDescription())
	})

	t.Run("PlayerTargetFilter description", func(t *testing.T) {
		filter := abilities.NewPlayerTargetFilter()
		assert.Equal(t, "target player", filter.GetDescription())

		filter = abilities.NewOpponentTargetFilter()
		assert.Equal(t, "target opponent", filter.GetDescription())
	})
}

func TestSpellAbilityBuilder(t *testing.T) {
	sourceID := uuid.New()

	t.Run("Build simple damage spell", func(t *testing.T) {
		ability, err := abilities.NewSpellAbilityBuilder(sourceID, "{R}").
			AddEffect(abilities.NewDamageEffect(3)).
			AddTarget(abilities.NewAnyTargetFilter()).
			Build()

		require.NoError(t, err)
		assert.NotNil(t, ability)
		assert.Equal(t, abilities.AbilityTypeSpell, ability.GetType())
		assert.Equal(t, "{R}", ability.GetManaCost().String())
	})

	t.Run("Build draw spell", func(t *testing.T) {
		ability, err := abilities.NewSpellAbilityBuilder(sourceID, "{2}{U}").
			AddEffect(abilities.NewDrawCardsEffect(2)).
			Build()

		require.NoError(t, err)
		assert.NotNil(t, ability)
	})
}

func TestActivatedAbilityBuilder(t *testing.T) {
	sourceID := uuid.New()

	t.Run("Build tap-for-mana ability", func(t *testing.T) {
		mana := abilities.NewMana()
		mana.Green = 1

		ability := abilities.NewActivatedAbilityBuilder(sourceID).
			AddTapCost().
			AddEffect(abilities.NewAddManaEffect(mana)).
			SetManaAbility().
			Build()

		assert.NotNil(t, ability)
		assert.Equal(t, abilities.AbilityTypeMana, ability.GetType())
		assert.False(t, ability.UsesStack)
	})

	t.Run("Build damage ability", func(t *testing.T) {
		ability := abilities.NewActivatedAbilityBuilder(sourceID).
			AddTapCost().
			AddEffect(abilities.NewDamageEffect(1)).
			AddTarget(abilities.NewAnyTargetFilter()).
			Build()

		assert.NotNil(t, ability)
		assert.Equal(t, abilities.AbilityTypeActivated, ability.GetType())
		assert.True(t, ability.UsesStack)
	})

	t.Run("Build ability with mana cost", func(t *testing.T) {
		ability := abilities.NewActivatedAbilityBuilder(sourceID).
			AddManaCost("{2}{R}").
			AddEffect(abilities.NewDamageEffect(5)).
			Build()

		assert.NotNil(t, ability)
		assert.Len(t, ability.Costs, 1)
	})
}

func TestConvenienceFunctions(t *testing.T) {
	sourceID := uuid.New()

	t.Run("BuildSimpleManaAbility for green", func(t *testing.T) {
		ability := abilities.BuildSimpleManaAbility(sourceID, "G")
		assert.NotNil(t, ability)
		assert.Equal(t, abilities.AbilityTypeMana, ability.GetType())
	})

	t.Run("BuildSimpleDamageAbility", func(t *testing.T) {
		ability := abilities.BuildSimpleDamageAbility(sourceID, 1)
		assert.NotNil(t, ability)
		assert.Equal(t, abilities.AbilityTypeActivated, ability.GetType())
	})
}
