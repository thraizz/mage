package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Phyrexian Tyranny", NewPhyrexianTyranny)
}

// NewPhyrexianTyranny creates a Phyrexian Tyranny
// {U}{B}{R} - ENCHANTMENT
func NewPhyrexianTyranny(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Phyrexian Tyranny")
	card.ManaCost = "{U}{B}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
