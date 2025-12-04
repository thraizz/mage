package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Struggle For Project Purity", NewStruggleForProjectPurity)
}

// NewStruggleForProjectPurity creates a Struggle For Project Purity
// {3}{U} - ENCHANTMENT
func NewStruggleForProjectPurity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Struggle For Project Purity")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
