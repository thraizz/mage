package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Escape To The Wilds", NewEscapeToTheWilds)
}

// NewEscapeToTheWilds creates a Escape To The Wilds
// {3}{R}{G} - SORCERY
func NewEscapeToTheWilds(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Escape To The Wilds")
	card.ManaCost = "{3}{R}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
