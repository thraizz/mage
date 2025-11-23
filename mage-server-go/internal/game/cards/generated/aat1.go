package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("A A T1", NewAAT1)
}

// NewAAT1 creates a A A T1
// {1}{W}{U}{B} - ARTIFACT CREATURE
func NewAAT1(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "A A T1")
	card.ManaCost = "{1}{W}{U}{B}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"DROID", "CONSTRUCT"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new LoseLifeTargetEffect(1), new ManaCostsImpl<>("...)
	// card.AddAbility(ability0)
	return card, nil
}
