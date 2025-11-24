package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Kenriths Royal Funeral", NewTheKenrithsRoyalFuneral)
}

// NewTheKenrithsRoyalFuneral creates a The Kenriths Royal Funeral
// {2}{W}{B} - ENCHANTMENT
func NewTheKenrithsRoyalFuneral(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Kenriths Royal Funeral")
	card.ManaCost = "{2}{W}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}