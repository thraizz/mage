package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Glacierwood Siege", NewGlacierwoodSiege)
}

// NewGlacierwoodSiege creates a Glacierwood Siege
// {1}{G}{U} - ENCHANTMENT
func NewGlacierwoodSiege(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Glacierwood Siege")
	card.ManaCost = "{1}{G}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}