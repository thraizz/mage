package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Assaultron Dominator", NewAssaultronDominator)
}

// NewAssaultronDominator creates a Assaultron Dominator
// {1}{R} - ARTIFACT CREATURE
func NewAssaultronDominator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Assaultron Dominator")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new AssaultronDominatorEffect(), new PayEnergyCost...)
	// card.AddAbility(ability0)
	return card, nil
}
