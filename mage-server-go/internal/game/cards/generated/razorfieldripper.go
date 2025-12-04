package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Razorfield Ripper", NewRazorfieldRipper)
}

// NewRazorfieldRipper creates a Razorfield Ripper
// {2}{W} - ARTIFACT CREATURE
func NewRazorfieldRipper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Razorfield Ripper")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"EQUIPMENT", "RHINO"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: RazorfieldRipperTriggeredAbility
	//   - Effect: GetEnergyCountersControllerEffect()
	// card.AddAbility(ability0)
	return card, nil
}
