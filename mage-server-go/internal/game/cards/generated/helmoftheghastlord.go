package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Helm Of The Ghastlord", NewHelmOfTheGhastlord)
}

// NewHelmOfTheGhastlord creates a Helm Of The Ghastlord
// {3}{U/B} - ENCHANTMENT
func NewHelmOfTheGhastlord(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Helm Of The Ghastlord")
	card.ManaCost = "{3}{U/B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1)
	// card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEnchantedEffect(1, 1)).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(new DealsDamageToAPlayerTriggeredAbility(new DrawCardSourceControllerEffect(1),false), AttachmentType.AURA)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1)
	// card.AddAbility(ability3)
	return card, nil
}