package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kurkesh Onakke Ancient", NewKurkeshOnakkeAncient)
}

// NewKurkeshOnakkeAncient creates a Kurkesh Onakke Ancient
// {2}{R}{R} - CREATURE
func NewKurkeshOnakkeAncient(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kurkesh Onakke Ancient")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OGRE", "SPIRIT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new CopyStackObjectEffect(), new ManaCostsImpl<>("...)
	// card.AddAbility(ability0)
	return card, nil
}