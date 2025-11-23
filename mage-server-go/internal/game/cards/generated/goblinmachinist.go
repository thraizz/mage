package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Goblin Machinist", NewGoblinMachinist)
}

// NewGoblinMachinist creates a Goblin Machinist
// {4}{R} - CREATURE
func NewGoblinMachinist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Goblin Machinist")
	card.ManaCost = "{4}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN"}
	card.Power = "0"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
