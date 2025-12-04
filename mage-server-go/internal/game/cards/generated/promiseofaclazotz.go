package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Promise Of Aclazotz", NewPromiseOfAclazotz)
}

// NewPromiseOfAclazotz creates a Promise Of Aclazotz
// {1}{B} - ENCHANTMENT
func NewPromiseOfAclazotz(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Promise Of Aclazotz")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfEndStepTriggeredAbility
	//   - Effect: DoIfCostPaid(                 new PopulateEffect(),            ...)
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new CreateTokenEffect(new Vampire...)
	// card.AddAbility(ability1)
	return card, nil
}
