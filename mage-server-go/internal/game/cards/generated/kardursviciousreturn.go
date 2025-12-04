package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kardurs Vicious Return", NewKardursViciousReturn)
}

// NewKardursViciousReturn creates a Kardurs Vicious Return
// {2}{B}{R} - ENCHANTMENT
func NewKardursViciousReturn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kardurs Vicious Return")
	card.ManaCost = "{2}{B}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
