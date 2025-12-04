package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Charge Of The Forever Beast", NewChargeOfTheForeverBeast)
}

// NewChargeOfTheForeverBeast creates a Charge Of The Forever Beast
// {2}{G} - SORCERY
func NewChargeOfTheForeverBeast(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Charge Of The Forever Beast")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
