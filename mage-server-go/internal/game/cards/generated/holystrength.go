package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Holy Strength", NewHolyStrength)
}

// NewHolyStrength creates a Holy Strength
// {W} - ENCHANTMENT
func NewHolyStrength(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Holy Strength")
	card.ManaCost = "{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomeBoostCreature)).
		AddEffect(abilities.NewBoostEnchantedEffect(1, 2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEnchantedEffect(1, 2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
