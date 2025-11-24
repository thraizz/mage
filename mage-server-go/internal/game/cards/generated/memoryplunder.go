package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Memory Plunder", NewMemoryPlunder)
}

// NewMemoryPlunder creates a Memory Plunder
// {U/B}{U/B}{U/B}{U/B} - INSTANT
func NewMemoryPlunder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Memory Plunder")
	card.ManaCost = "{U/B}{U/B}{U/B}{U/B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}