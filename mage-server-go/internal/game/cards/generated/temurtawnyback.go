package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Temur Tawnyback", NewTemurTawnyback)
}

// NewTemurTawnyback creates a Temur Tawnyback
// {2/G}{2/U}{2/R} - CREATURE
func NewTemurTawnyback(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Temur Tawnyback")
	card.ManaCost = "{2/G}{2/U}{2/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAST"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
