package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Forging The Tyrite Sword", NewForgingTheTyriteSword)
}

// NewForgingTheTyriteSword creates a Forging The Tyrite Sword
// {1}{R}{W} - ENCHANTMENT
func NewForgingTheTyriteSword(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Forging The Tyrite Sword")
	card.ManaCost = "{1}{R}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
