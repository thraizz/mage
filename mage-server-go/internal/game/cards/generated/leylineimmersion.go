package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Leyline Immersion", NewLeylineImmersion)
}

// NewLeylineImmersion creates a Leyline Immersion
// {3}{G} - ENCHANTMENT
func NewLeylineImmersion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Leyline Immersion")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
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
		AddEffect(abilities.NewGainAbilityAttachedEffect(new ConditionalAnyColorManaAbility(5, new LeylineImmersionManaBuilder()), AttachmentType.AURA)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}