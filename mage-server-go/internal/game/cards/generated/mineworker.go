package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mine Worker", NewMineWorker)
}

// NewMineWorker creates a Mine Worker
// {2} - ARTIFACT CREATURE
func NewMineWorker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mine Worker")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ASSEMBLY_WORKER"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
