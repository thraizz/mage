package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Corruption Of Towashi", NewCorruptionOfTowashi)
}

// NewCorruptionOfTowashi creates a Corruption Of Towashi
// {4}{U} - ENCHANTMENT
func NewCorruptionOfTowashi(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Corruption Of Towashi")
	card.ManaCost = "{4}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
