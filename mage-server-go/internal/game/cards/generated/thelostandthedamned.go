package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Lost And The Damned", NewTheLostAndTheDamned)
}

// NewTheLostAndTheDamned creates a The Lost And The Damned
// {1}{U}{R} - ENCHANTMENT
func NewTheLostAndTheDamned(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Lost And The Damned")
	card.ManaCost = "{1}{U}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}