package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wings Of Aesthir", NewWingsOfAesthir)
}

// NewWingsOfAesthir creates a Wings Of Aesthir
// {W}{U} - ENCHANTMENT
func NewWingsOfAesthir(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wings Of Aesthir")
	card.ManaCost = "{W}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEnchantedEffect(1, 0)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FlyingAbility", abilities.AttachmentTypeAura)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FirstStrikeAbility", abilities.AttachmentTypeAura)).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomeBenefit)).
		AddEffect(abilities.NewBoostEnchantedEffect(1, 0)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FlyingAbility", abilities.AttachmentTypeAura)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FirstStrikeAbility", abilities.AttachmentTypeAura)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
