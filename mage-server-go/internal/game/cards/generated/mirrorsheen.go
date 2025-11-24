package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mirror Sheen", NewMirrorSheen)
}

// NewMirrorSheen creates a Mirror Sheen
// {1}{U/R}{U/R} - ENCHANTMENT
func NewMirrorSheen(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mirror Sheen")
	card.ManaCost = "{1}{U/R}{U/R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
