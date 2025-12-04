package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Swarm Intelligence", NewSwarmIntelligence)
}

// NewSwarmIntelligence creates a Swarm Intelligence
// {6}{U} - ENCHANTMENT
func NewSwarmIntelligence(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Swarm Intelligence")
	card.ManaCost = "{6}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
