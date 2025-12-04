package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Windcrag Siege", NewWindcragSiege)
}

// NewWindcragSiege creates a Windcrag Siege
// {1}{R}{W} - ENCHANTMENT
func NewWindcragSiege(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Windcrag Siege")
	card.ManaCost = "{1}{R}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
