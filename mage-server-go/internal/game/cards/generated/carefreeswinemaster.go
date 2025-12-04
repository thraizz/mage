package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Carefree Swinemaster", NewCarefreeSwinemaster)
}

// NewCarefreeSwinemaster creates a Carefree Swinemaster
// {2}{G} - CREATURE
func NewCarefreeSwinemaster(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Carefree Swinemaster")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GNOME", "RANGER"}
	card.Power = "1"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new CreateTokenEffect(                 new Boar2To...)
	// card.AddAbility(ability0)
	return card, nil
}
