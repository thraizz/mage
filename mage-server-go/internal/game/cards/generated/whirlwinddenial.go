package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Whirlwind Denial", NewWhirlwindDenial)
}

// NewWhirlwindDenial creates a Whirlwind Denial
// {2}{U} - INSTANT
func NewWhirlwindDenial(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Whirlwind Denial")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
