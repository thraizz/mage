package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kyoshi Warrior Guard", NewKyoshiWarriorGuard)
}

// NewKyoshiWarriorGuard creates a Kyoshi Warrior Guard
// {1}{W} - CREATURE
func NewKyoshiWarriorGuard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kyoshi Warrior Guard")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WARRIOR", "ALLY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
