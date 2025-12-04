package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Iron Fist Of The Empire", NewIronFistOfTheEmpire)
}

// NewIronFistOfTheEmpire creates a Iron Fist Of The Empire
// {U}{B}{R} - ENCHANTMENT
func NewIronFistOfTheEmpire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Iron Fist Of The Empire")
	card.ManaCost = "{U}{B}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
