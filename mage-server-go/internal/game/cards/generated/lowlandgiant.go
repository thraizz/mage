package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lowland Giant", NewLowlandGiant)
}

// NewLowlandGiant creates a Lowland Giant
// {2}{R}{R} - CREATURE
func NewLowlandGiant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lowland Giant")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
