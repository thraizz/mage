package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Krazy Kow", NewKrazyKow)
}

// NewKrazyKow creates a Krazy Kow
// {3}{R} - CREATURE
func NewKrazyKow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Krazy Kow")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"COW"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
