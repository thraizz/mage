package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Searing Meditation", NewSearingMeditation)
}

// NewSearingMeditation creates a Searing Meditation
// {1}{R}{W} - ENCHANTMENT
func NewSearingMeditation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Searing Meditation")
	card.ManaCost = "{1}{R}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new DamageTargetEffect(2), new GenericManaCost(2))
	// card.AddAbility(ability0)
	return card, nil
}