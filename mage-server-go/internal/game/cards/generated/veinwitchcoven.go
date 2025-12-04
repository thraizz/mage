package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Veinwitch Coven", NewVeinwitchCoven)
}

// NewVeinwitchCoven creates a Veinwitch Coven
// {2}{B} - CREATURE
func NewVeinwitchCoven(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Veinwitch Coven")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE", "WARLOCK"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: GainLifeControllerTriggeredAbility
	//   - Effect: DoIfCostPaid(                 new ReturnFromGraveyardToHandTarg...)
	// card.AddAbility(ability0)
	return card, nil
}
