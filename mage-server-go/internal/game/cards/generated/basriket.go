package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Basri Ket", NewBasriKet)
}

// NewBasriKet creates a Basri Ket
// {1}{W}{W} - PLANESWALKER
func NewBasriKet(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Basri Ket")
	card.ManaCost = "{1}{W}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"BASRI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
