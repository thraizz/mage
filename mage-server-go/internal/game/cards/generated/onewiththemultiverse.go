package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("One With The Multiverse", NewOneWithTheMultiverse)
}

// NewOneWithTheMultiverse creates a One With The Multiverse
// {6}{U}{U} - ENCHANTMENT
func NewOneWithTheMultiverse(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "One With The Multiverse")
	card.ManaCost = "{6}{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
