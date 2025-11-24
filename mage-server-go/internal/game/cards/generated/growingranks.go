package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Growing Ranks", NewGrowingRanks)
}

// NewGrowingRanks creates a Growing Ranks
// {2}{G/W}{G/W} - ENCHANTMENT
func NewGrowingRanks(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Growing Ranks")
	card.ManaCost = "{2}{G/W}{G/W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}