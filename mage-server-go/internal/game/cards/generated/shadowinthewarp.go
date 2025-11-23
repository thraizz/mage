package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shadow In The Warp", NewShadowInTheWarp)
}

// NewShadowInTheWarp creates a Shadow In The Warp
// {1}{R}{G} - ENCHANTMENT
func NewShadowInTheWarp(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shadow In The Warp")
	card.ManaCost = "{1}{R}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
