package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Surgespanner", NewSurgespanner)
}

// NewSurgespanner creates a Surgespanner
// {2}{U}{U} - CREATURE
func NewSurgespanner(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Surgespanner")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MERFOLK", "WIZARD"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BecomesTappedSourceTriggeredAbility
	//   - Effect: DoIfCostPaid(new ReturnToHandTargetEffect(), new ManaCostsImpl<...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new ReturnToHandTargetEffect(), new ManaCostsImpl<...)
	// card.AddAbility(ability1)
	return card, nil
}
