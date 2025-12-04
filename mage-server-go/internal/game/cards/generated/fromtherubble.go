package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("From The Rubble", NewFromTheRubble)
}

// NewFromTheRubble creates a From The Rubble
// {4}{W}{W} - ENCHANTMENT
func NewFromTheRubble(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "From The Rubble")
	card.ManaCost = "{4}{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfEndStepTriggeredAbility
	//   - Effect: ReturnFromGraveyardToBattlefieldWithCounterTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
