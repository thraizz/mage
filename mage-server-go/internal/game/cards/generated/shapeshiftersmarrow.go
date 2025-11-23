package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shapeshifters Marrow", NewShapeshiftersMarrow)
}

// NewShapeshiftersMarrow creates a Shapeshifters Marrow
// {2}{U}{U} - ENCHANTMENT
func NewShapeshiftersMarrow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shapeshifters Marrow")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
