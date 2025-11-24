package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Equilibrium", NewEquilibrium)
}

// NewEquilibrium creates a Equilibrium
// {1}{U}{U} - ENCHANTMENT
func NewEquilibrium(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Equilibrium")
	card.ManaCost = "{1}{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new ReturnToHandTargetEffect(), new GenericManaCos...)
	// card.AddAbility(ability0)
	return card, nil
}
