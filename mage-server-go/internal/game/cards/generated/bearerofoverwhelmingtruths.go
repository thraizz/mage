package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bearer Of Overwhelming Truths", NewBearerOfOverwhelmingTruths)
}

// NewBearerOfOverwhelmingTruths creates a Bearer Of Overwhelming Truths
//   - CREATURE
func NewBearerOfOverwhelmingTruths(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bearer Of Overwhelming Truths")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
