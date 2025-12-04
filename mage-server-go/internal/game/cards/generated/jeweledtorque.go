package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jeweled Torque", NewJeweledTorque)
}

// NewJeweledTorque creates a Jeweled Torque
// {2} - ARTIFACT
func NewJeweledTorque(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jeweled Torque")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                         new GainLifeEffect(2),   ...)
	// card.AddAbility(ability0)
	return card, nil
}
