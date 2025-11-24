package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Disturbing Mirth", NewDisturbingMirth)
}

// NewDisturbingMirth creates a Disturbing Mirth
// {B}{R} - ENCHANTMENT
func NewDisturbingMirth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Disturbing Mirth")
	card.ManaCost = "{B}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new DrawCardSourceControllerEffect(2), new Sacrifi...)
	// card.AddAbility(ability0)
	return card, nil
}