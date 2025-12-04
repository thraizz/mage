package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jedit Ojanen Mercenary", NewJeditOjanenMercenary)
}

// NewJeditOjanenMercenary creates a Jedit Ojanen Mercenary
// {1}{W}{U} - CREATURE
func NewJeditOjanenMercenary(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jedit Ojanen Mercenary")
	card.ManaCost = "{1}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new CreateTokenEffect(new CatWarr...)
	// card.AddAbility(ability0)
	return card, nil
}
