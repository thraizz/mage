package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Najal The Storm Runner", NewNajalTheStormRunner)
}

// NewNajalTheStormRunner creates a Najal The Storm Runner
// {2}{U}{U}{R} - CREATURE
func NewNajalTheStormRunner(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Najal The Storm Runner")
	card.ManaCost = "{2}{U}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"EFREET", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new CreateDelayedTriggeredAbility...)
	// card.AddAbility(ability0)
	return card, nil
}
