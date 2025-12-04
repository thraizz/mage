package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("King Narfis Betrayal", NewKingNarfisBetrayal)
}

// NewKingNarfisBetrayal creates a King Narfis Betrayal
// {1}{U}{B} - ENCHANTMENT
func NewKingNarfisBetrayal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "King Narfis Betrayal")
	card.ManaCost = "{1}{U}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
