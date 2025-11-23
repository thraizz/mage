package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Phyrexian Arena", NewPhyrexianArena)
}

// NewPhyrexianArena creates a Phyrexian Arena
// {1}{B}{B} - ENCHANTMENT
func NewPhyrexianArena(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Phyrexian Arena")
	card.ManaCost = "{1}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
