package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Restoration Of Eiganjo", NewTheRestorationOfEiganjo)
}

// NewTheRestorationOfEiganjo creates a The Restoration Of Eiganjo
// {2}{W} - ENCHANTMENT
func NewTheRestorationOfEiganjo(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Restoration Of Eiganjo")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
