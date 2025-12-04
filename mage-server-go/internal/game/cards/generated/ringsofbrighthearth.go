package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rings Of Brighthearth", NewRingsOfBrighthearth)
}

// NewRingsOfBrighthearth creates a Rings Of Brighthearth
// {3} - ARTIFACT
func NewRingsOfBrighthearth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rings Of Brighthearth")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new CopyStackObjectEffect(), new ManaCostsImpl<>("...)
	// card.AddAbility(ability0)
	return card, nil
}
