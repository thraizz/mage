package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("March Of The World Ooze", NewMarchOfTheWorldOoze)
}

// NewMarchOfTheWorldOoze creates a March Of The World Ooze
// {3}{G}{G}{G} - ENCHANTMENT
func NewMarchOfTheWorldOoze(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "March Of The World Ooze")
	card.ManaCost = "{3}{G}{G}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}