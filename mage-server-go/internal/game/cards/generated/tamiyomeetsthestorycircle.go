package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tamiyo Meets The Story Circle", NewTamiyoMeetsTheStoryCircle)
}

// NewTamiyoMeetsTheStoryCircle creates a Tamiyo Meets The Story Circle
// {1}{U} - ENCHANTMENT
func NewTamiyoMeetsTheStoryCircle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tamiyo Meets The Story Circle")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
