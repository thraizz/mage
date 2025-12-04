package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Order Of The Golden Cricket", NewOrderOfTheGoldenCricket)
}

// NewOrderOfTheGoldenCricket creates a Order Of The Golden Cricket
// {1}{W} - CREATURE
func NewOrderOfTheGoldenCricket(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Order Of The Golden Cricket")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KITHKIN", "KNIGHT"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: AttacksTriggeredAbility
	//   - Effect: DoIfCostPaid(                 new GainAbilitySourceEffect(Flyin...)
	// card.AddAbility(ability0)
	return card, nil
}
