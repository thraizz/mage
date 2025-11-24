package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Boggart Shenanigans", NewBoggartShenanigans)
}

// NewBoggartShenanigans creates a Boggart Shenanigans
// {2}{R} - KINDRED ENCHANTMENT
func NewBoggartShenanigans(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Boggart Shenanigans")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"KINDRED", "ENCHANTMENT"}
	card.Subtypes = []string{"GOBLIN"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}