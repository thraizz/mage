package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Glistening Goremonger", NewGlisteningGoremonger)
}

// NewGlisteningGoremonger creates a Glistening Goremonger
//   - CREATURE
func NewGlisteningGoremonger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Glistening Goremonger")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "DEVIL"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
