package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rowan Kenrith", NewRowanKenrith)
}

// NewRowanKenrith creates a Rowan Kenrith
// {4}{R}{R} - PLANESWALKER
func NewRowanKenrith(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rowan Kenrith")
	card.ManaCost = "{4}{R}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"ROWAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}