package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Macabre Mockery", NewMacabreMockery)
}

// NewMacabreMockery creates a Macabre Mockery
// {2}{B}{R} - INSTANT
func NewMacabreMockery(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Macabre Mockery")
	card.ManaCost = "{2}{B}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect("sacrifice " + permanent.getLogName(), controller....)
	// card.AddAbility(ability0)
	return card, nil
}
