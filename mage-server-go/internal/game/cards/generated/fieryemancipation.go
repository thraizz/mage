package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fiery Emancipation", NewFieryEmancipation)
}

// NewFieryEmancipation creates a Fiery Emancipation
// {3}{R}{R}{R} - ENCHANTMENT
func NewFieryEmancipation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fiery Emancipation")
	card.ManaCost = "{3}{R}{R}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
