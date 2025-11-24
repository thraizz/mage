package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Blood Of The Martyr", NewBloodOfTheMartyr)
}

// NewBloodOfTheMartyr creates a Blood Of The Martyr
// {W}{W}{W} - INSTANT
func NewBloodOfTheMartyr(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Blood Of The Martyr")
	card.ManaCost = "{W}{W}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}