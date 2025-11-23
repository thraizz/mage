package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Strength Of Isolation", NewStrengthOfIsolation)
}

// NewStrengthOfIsolation creates a Strength Of Isolation
// {1}{W} - ENCHANTMENT
func NewStrengthOfIsolation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Strength Of Isolation")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect(ProtectionAbility.from(ObjectColor.BLACK), AttachmentType.AURA)).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomeBoostCreature)).
		AddEffect(abilities.NewBoostEnchantedEffect(1, 2)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(ProtectionAbility.from(ObjectColor.BLACK), AttachmentType.AURA)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEnchantedEffect(1, 2)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(ProtectionAbility.from(ObjectColor.BLACK), AttachmentType.AURA)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
