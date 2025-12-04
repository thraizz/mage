package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flooded Woodlands", NewFloodedWoodlands)
}

// NewFloodedWoodlands creates a Flooded Woodlands
// {2}{U}{B} - ENCHANTMENT
func NewFloodedWoodlands(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flooded Woodlands")
	card.ManaCost = "{2}{U}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
