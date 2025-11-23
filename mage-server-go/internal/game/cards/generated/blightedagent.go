package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Blighted Agent", NewBlightedAgent)
}

// NewBlightedAgent creates a Blighted Agent
// {1}{U} - CREATURE
func NewBlightedAgent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Blighted Agent")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"CREATURE"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
