package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Knight Of New Benalia", NewKnightOfNewBenalia)
}

// NewKnightOfNewBenalia creates a Knight Of New Benalia
// {1}{W} - CREATURE
func NewKnightOfNewBenalia(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Knight Of New Benalia")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "KNIGHT"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
