package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Synchronized Charge", NewSynchronizedCharge)
}

// NewSynchronizedCharge creates a Synchronized Charge
// {1}{G} - SORCERY
func NewSynchronizedCharge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Synchronized Charge")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
