package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Harald Unites The Elves", NewHaraldUnitesTheElves)
}

// NewHaraldUnitesTheElves creates a Harald Unites The Elves
// {2}{B}{G} - ENCHANTMENT
func NewHaraldUnitesTheElves(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Harald Unites The Elves")
	card.ManaCost = "{2}{B}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
