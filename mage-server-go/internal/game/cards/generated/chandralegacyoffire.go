package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chandra Legacy Of Fire", NewChandraLegacyOfFire)
}

// NewChandraLegacyOfFire creates a Chandra Legacy Of Fire
// {4}{R} - PLANESWALKER
func NewChandraLegacyOfFire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chandra Legacy Of Fire")
	card.ManaCost = "{4}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"CHANDRA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
