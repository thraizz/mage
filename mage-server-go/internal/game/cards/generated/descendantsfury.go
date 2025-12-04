package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Descendants Fury", NewDescendantsFury)
}

// NewDescendantsFury creates a Descendants Fury
// {3}{R} - ENCHANTMENT
func NewDescendantsFury(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Descendants Fury")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: OneOrMoreCombatDamagePlayerTriggeredAbility
	//   - Effect: DoIfCostPaid(                         new DescendantsFuryEffect...)
	// card.AddAbility(ability0)
	return card, nil
}
