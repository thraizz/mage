package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Isu The Abominable", NewIsuTheAbominable)
}

// NewIsuTheAbominable creates a Isu The Abominable
// {3}{U}{U} - CREATURE
func NewIsuTheAbominable(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Isu The Abominable")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"YETI"}
	card.Supertypes = []string{"LEGENDARY", "SNOW"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                         new AddCountersSourceEffe...)
	// card.AddAbility(ability0)
	return card, nil
}
