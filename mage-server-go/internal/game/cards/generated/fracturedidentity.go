package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fractured Identity", NewFracturedIdentity)
}

// NewFracturedIdentity creates a Fractured Identity
// {3}{W}{U} - SORCERY
func NewFracturedIdentity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fractured Identity")
	card.ManaCost = "{3}{W}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
