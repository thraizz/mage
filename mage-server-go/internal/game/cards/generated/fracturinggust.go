package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fracturing Gust", NewFracturingGust)
}

// NewFracturingGust creates a Fracturing Gust
// {2}{G/W}{G/W}{G/W} - INSTANT
func NewFracturingGust(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fracturing Gust")
	card.ManaCost = "{2}{G/W}{G/W}{G/W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
