package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Winters Night", NewWintersNight)
}

// NewWintersNight creates a Winters Night
// {R}{G}{W} - ENCHANTMENT
func NewWintersNight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Winters Night")
	card.ManaCost = "{R}{G}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"WORLD"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
