package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("An Zerrin Ruins", NewAnZerrinRuins)
}

// NewAnZerrinRuins creates a An Zerrin Ruins
// {2}{R}{R} - ENCHANTMENT
func NewAnZerrinRuins(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "An Zerrin Ruins")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
