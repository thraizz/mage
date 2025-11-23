package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rent Is Due", NewRentIsDue)
}

// NewRentIsDue creates a Rent Is Due
// {W} - ENCHANTMENT
func NewRentIsDue(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rent Is Due")
	card.ManaCost = "{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeSourceEffect()
	//   - DoIfCostPaid(                 new DrawCardSourceControllerEffec...)
	// card.AddAbility(ability0)
	return card, nil
}
