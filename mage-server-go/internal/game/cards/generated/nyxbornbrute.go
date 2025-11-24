package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nyxborn Brute", NewNyxbornBrute)
}

// NewNyxbornBrute creates a Nyxborn Brute
// {3}{R}{R} - ENCHANTMENT CREATURE
func NewNyxbornBrute(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nyxborn Brute")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"CYCLOPS"}
	card.Power = "7"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
