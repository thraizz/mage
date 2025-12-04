package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Case Of The Gateway Express", NewCaseOfTheGatewayExpress)
}

// NewCaseOfTheGatewayExpress creates a Case Of The Gateway Express
// {1}{W} - ENCHANTMENT
func NewCaseOfTheGatewayExpress(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Case Of The Gateway Express")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"CASE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
