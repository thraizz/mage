package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Battle For Bretagard", NewBattleForBretagard)
}

// NewBattleForBretagard creates a Battle For Bretagard
// {1}{G}{W} - ENCHANTMENT
func NewBattleForBretagard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Battle For Bretagard")
	card.ManaCost = "{1}{G}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
