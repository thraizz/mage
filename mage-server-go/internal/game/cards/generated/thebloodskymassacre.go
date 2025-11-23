package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Bloodsky Massacre", NewTheBloodskyMassacre)
}

// NewTheBloodskyMassacre creates a The Bloodsky Massacre
// {1}{B}{R} - ENCHANTMENT
func NewTheBloodskyMassacre(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Bloodsky Massacre")
	card.ManaCost = "{1}{B}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
