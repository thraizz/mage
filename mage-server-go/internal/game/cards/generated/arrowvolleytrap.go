package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Arrow Volley Trap", NewArrowVolleyTrap)
}

// NewArrowVolleyTrap creates a Arrow Volley Trap
// {3}{W}{W} - INSTANT
func NewArrowVolleyTrap(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arrow Volley Trap")
	card.ManaCost = "{3}{W}{W}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"TRAP"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
