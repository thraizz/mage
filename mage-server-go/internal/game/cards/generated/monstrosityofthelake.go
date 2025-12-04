package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Monstrosity Of The Lake", NewMonstrosityOfTheLake)
}

// NewMonstrosityOfTheLake creates a Monstrosity Of The Lake
// {4}{U} - CREATURE
func NewMonstrosityOfTheLake(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Monstrosity Of The Lake")
	card.ManaCost = "{4}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KRAKEN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new TapAllEffect(StaticFilters.FI...)
	// card.AddAbility(ability0)
	return card, nil
}
