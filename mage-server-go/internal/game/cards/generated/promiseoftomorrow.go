package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Promise Of Tomorrow", NewPromiseOfTomorrow)
}

// NewPromiseOfTomorrow creates a Promise Of Tomorrow
// {2}{W} - ENCHANTMENT
func NewPromiseOfTomorrow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Promise Of Tomorrow")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfEndStepTriggeredAbility
	//   - Effect: SacrificeSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
