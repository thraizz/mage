package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chandra Hopes Beacon", NewChandraHopesBeacon)
}

// NewChandraHopesBeacon creates a Chandra Hopes Beacon
// {4}{R}{R} - PLANESWALKER
func NewChandraHopesBeacon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chandra Hopes Beacon")
	card.ManaCost = "{4}{R}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
