package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mayaels Aria", NewMayaelsAria)
}

// NewMayaelsAria creates a Mayaels Aria
// {R}{G}{W} - ENCHANTMENT
func NewMayaelsAria(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mayaels Aria")
	card.ManaCost = "{R}{G}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
