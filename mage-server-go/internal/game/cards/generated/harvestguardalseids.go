package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Harvestguard Alseids", NewHarvestguardAlseids)
}

// NewHarvestguardAlseids creates a Harvestguard Alseids
// {2}{W} - ENCHANTMENT CREATURE
func NewHarvestguardAlseids(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Harvestguard Alseids")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"NYMPH"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
