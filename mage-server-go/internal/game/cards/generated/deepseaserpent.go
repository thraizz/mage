package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Deep Sea Serpent", NewDeepSeaSerpent)
}

// NewDeepSeaSerpent creates a Deep Sea Serpent
// {4}{U}{U} - CREATURE
func NewDeepSeaSerpent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Deep Sea Serpent")
	card.ManaCost = "{4}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SERPENT"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
