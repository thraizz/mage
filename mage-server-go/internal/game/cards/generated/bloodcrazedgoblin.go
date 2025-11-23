package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bloodcrazed Goblin", NewBloodcrazedGoblin)
}

// NewBloodcrazedGoblin creates a Bloodcrazed Goblin
// {R} - CREATURE
func NewBloodcrazedGoblin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bloodcrazed Goblin")
	card.ManaCost = "{R}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
