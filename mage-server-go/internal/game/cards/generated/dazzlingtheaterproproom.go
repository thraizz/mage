package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dazzling Theater Prop Room", NewDazzlingTheaterPropRoom)
}

// NewDazzlingTheaterPropRoom creates a Dazzling Theater Prop Room
// {3}{W} - ENCHANTMENT
func NewDazzlingTheaterPropRoom(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dazzling Theater Prop Room")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"ROOM"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
