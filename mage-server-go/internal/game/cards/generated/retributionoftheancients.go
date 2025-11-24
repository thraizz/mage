package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Retribution Of The Ancients", NewRetributionOfTheAncients)
}

// NewRetributionOfTheAncients creates a Retribution Of The Ancients
// {B} - ENCHANTMENT
func NewRetributionOfTheAncients(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Retribution Of The Ancients")
	card.ManaCost = "{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(xValue, xValue)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
