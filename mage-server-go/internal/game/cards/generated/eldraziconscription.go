package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eldrazi Conscription", NewEldraziConscription)
}

// NewEldraziConscription creates a Eldrazi Conscription
// {8} - KINDRED ENCHANTMENT
func NewEldraziConscription(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eldrazi Conscription")
	card.ManaCost = "{8}"
	card.Types = []string{"KINDRED", "ENCHANTMENT"}
	card.Subtypes = []string{"ELDRAZI", "AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomeBoostCreature)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEnchantedEffect(10, 10)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("TrampleAbility", abilities.AttachmentTypeAura)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("TrampleAbility", abilities.AttachmentTypeAura)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
