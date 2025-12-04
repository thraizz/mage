package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("March From The Black Gate", NewMarchFromTheBlackGate)
}

// NewMarchFromTheBlackGate creates a March From The Black Gate
// {1}{B} - ENCHANTMENT
func NewMarchFromTheBlackGate(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "March From The Black Gate")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
