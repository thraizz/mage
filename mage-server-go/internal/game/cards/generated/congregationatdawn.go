package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Congregation At Dawn", NewCongregationAtDawn)
}

// NewCongregationAtDawn creates a Congregation At Dawn
// {G}{G}{W} - INSTANT
func NewCongregationAtDawn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Congregation At Dawn")
	card.ManaCost = "{G}{G}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
