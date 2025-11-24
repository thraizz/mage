package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Elder Spawn", NewElderSpawn)
}

// NewElderSpawn creates a Elder Spawn
// {4}{U}{U}{U} - CREATURE
func NewElderSpawn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elder Spawn")
	card.ManaCost = "{4}{U}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPAWN"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}