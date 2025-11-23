package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Goblin Rimerunner", NewGoblinRimerunner)
}

// NewGoblinRimerunner creates a Goblin Rimerunner
// {2}{R} - CREATURE
func NewGoblinRimerunner(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Goblin Rimerunner")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "WARRIOR"}
	card.Supertypes = []string{"SNOW"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
