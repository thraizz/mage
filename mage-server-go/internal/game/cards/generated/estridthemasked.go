package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Estrid The Masked", NewEstridTheMasked)
}

// NewEstridTheMasked creates a Estrid The Masked
// {1}{G}{W}{U} - PLANESWALKER
func NewEstridTheMasked(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Estrid The Masked")
	card.ManaCost = "{1}{G}{W}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"ESTRID"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
