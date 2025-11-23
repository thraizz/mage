package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Invoke Calamity", NewInvokeCalamity)
}

// NewInvokeCalamity creates a Invoke Calamity
// {1}{R}{R}{R}{R} - INSTANT
func NewInvokeCalamity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invoke Calamity")
	card.ManaCost = "{1}{R}{R}{R}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
