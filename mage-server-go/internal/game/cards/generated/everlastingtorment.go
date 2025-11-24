package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Everlasting Torment", NewEverlastingTorment)
}

// NewEverlastingTorment creates a Everlasting Torment
// {2}{B/R} - ENCHANTMENT
func NewEverlastingTorment(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Everlasting Torment")
	card.ManaCost = "{2}{B/R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}