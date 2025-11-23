package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rod Of Spanking", NewRodOfSpanking)
}

// NewRodOfSpanking creates a Rod Of Spanking
//
//	-
func NewRodOfSpanking(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rod Of Spanking")
	card.ManaCost = ""
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
