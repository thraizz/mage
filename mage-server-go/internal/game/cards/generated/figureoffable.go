package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Figure Of Fable", NewFigureOfFable)
}

// NewFigureOfFable creates a Figure Of Fable
// {G/W} - CREATURE
func NewFigureOfFable(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Figure Of Fable")
	card.ManaCost = "{G/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KITHKIN"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
