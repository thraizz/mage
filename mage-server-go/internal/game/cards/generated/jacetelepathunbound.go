package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jace Telepath Unbound", NewJaceTelepathUnbound)
}

// NewJaceTelepathUnbound creates a Jace Telepath Unbound
//   - PLANESWALKER
func NewJaceTelepathUnbound(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jace Telepath Unbound")
	card.ManaCost = ""
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"JACE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
