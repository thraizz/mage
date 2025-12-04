package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Phasing Of Zhalfir", NewThePhasingOfZhalfir)
}

// NewThePhasingOfZhalfir creates a The Phasing Of Zhalfir
// {2}{U}{U} - ENCHANTMENT
func NewThePhasingOfZhalfir(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Phasing Of Zhalfir")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
