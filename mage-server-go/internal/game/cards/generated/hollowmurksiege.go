package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hollowmurk Siege", NewHollowmurkSiege)
}

// NewHollowmurkSiege creates a Hollowmurk Siege
// {B}{G} - ENCHANTMENT
func NewHollowmurkSiege(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hollowmurk Siege")
	card.ManaCost = "{B}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
