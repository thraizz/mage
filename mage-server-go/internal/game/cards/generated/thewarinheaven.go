package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The War In Heaven", NewTheWarInHeaven)
}

// NewTheWarInHeaven creates a The War In Heaven
// {3}{B}{B}{B} - ENCHANTMENT
func NewTheWarInHeaven(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The War In Heaven")
	card.ManaCost = "{3}{B}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}