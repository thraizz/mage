package cards_test

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/stretchr/testify/assert"
)

func TestCardInfo(t *testing.T) {
	t.Run("IsCreature", func(t *testing.T) {
		info := &cards.CardInfo{Types: []string{"CREATURE"}}
		assert.True(t, info.IsCreature())

		info = &cards.CardInfo{Types: []string{"INSTANT"}}
		assert.False(t, info.IsCreature())
	})

	t.Run("IsInstant", func(t *testing.T) {
		info := &cards.CardInfo{Types: []string{"INSTANT"}}
		assert.True(t, info.IsInstant())

		info = &cards.CardInfo{Types: []string{"SORCERY"}}
		assert.False(t, info.IsInstant())
	})

	t.Run("IsSorcery", func(t *testing.T) {
		info := &cards.CardInfo{Types: []string{"SORCERY"}}
		assert.True(t, info.IsSorcery())

		info = &cards.CardInfo{Types: []string{"INSTANT"}}
		assert.False(t, info.IsSorcery())
	})

	t.Run("IsEnchantment", func(t *testing.T) {
		info := &cards.CardInfo{Types: []string{"ENCHANTMENT"}}
		assert.True(t, info.IsEnchantment())

		info = &cards.CardInfo{Types: []string{"CREATURE"}}
		assert.False(t, info.IsEnchantment())
	})

	t.Run("IsArtifact", func(t *testing.T) {
		info := &cards.CardInfo{Types: []string{"ARTIFACT"}}
		assert.True(t, info.IsArtifact())

		info = &cards.CardInfo{Types: []string{"CREATURE"}}
		assert.False(t, info.IsArtifact())
	})

	t.Run("IsPlaneswalker", func(t *testing.T) {
		info := &cards.CardInfo{Types: []string{"PLANESWALKER"}}
		assert.True(t, info.IsPlaneswalker())

		info = &cards.CardInfo{Types: []string{"CREATURE"}}
		assert.False(t, info.IsPlaneswalker())
	})

	t.Run("IsLand", func(t *testing.T) {
		info := &cards.CardInfo{Types: []string{"LAND"}}
		assert.True(t, info.IsLand())

		info = &cards.CardInfo{Types: []string{"CREATURE"}}
		assert.False(t, info.IsLand())
	})

	t.Run("IsLegendary", func(t *testing.T) {
		info := &cards.CardInfo{Supertypes: []string{"LEGENDARY"}}
		assert.True(t, info.IsLegendary())

		info = &cards.CardInfo{Supertypes: []string{}}
		assert.False(t, info.IsLegendary())
	})

	t.Run("IsBasic", func(t *testing.T) {
		info := &cards.CardInfo{Supertypes: []string{"BASIC"}}
		assert.True(t, info.IsBasic())

		info = &cards.CardInfo{Supertypes: []string{}}
		assert.False(t, info.IsBasic())
	})

	t.Run("Multiple types", func(t *testing.T) {
		// Artifact Creature
		info := &cards.CardInfo{Types: []string{"ARTIFACT", "CREATURE"}}
		assert.True(t, info.IsCreature())
		assert.True(t, info.IsArtifact())
		assert.False(t, info.IsInstant())
	})

	t.Run("Legendary creature", func(t *testing.T) {
		info := &cards.CardInfo{
			Types:      []string{"CREATURE"},
			Subtypes:   []string{"HUMAN", "WIZARD"},
			Supertypes: []string{"LEGENDARY"},
		}
		assert.True(t, info.IsCreature())
		assert.True(t, info.IsLegendary())
		assert.False(t, info.IsBasic())
	})

	t.Run("Basic land", func(t *testing.T) {
		info := &cards.CardInfo{
			Types:      []string{"LAND"},
			Subtypes:   []string{"PLAINS"},
			Supertypes: []string{"BASIC"},
		}
		assert.True(t, info.IsLand())
		assert.True(t, info.IsBasic())
		assert.False(t, info.IsLegendary())
	})
}
