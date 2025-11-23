package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Serpent Of Yawning Depths", NewSerpentOfYawningDepths)
}

// NewSerpentOfYawningDepths creates a Serpent Of Yawning Depths
// {4}{U}{U} - ENCHANTMENT CREATURE
func NewSerpentOfYawningDepths(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Serpent Of Yawning Depths")
	card.ManaCost = "{4}{U}{U}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"SERPENT"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
