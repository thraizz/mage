package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vigean Intuition", NewVigeanIntuition)
}

// NewVigeanIntuition creates a Vigean Intuition
// {3}{G}{U} - INSTANT
func NewVigeanIntuition(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vigean Intuition")
	card.ManaCost = "{3}{G}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
