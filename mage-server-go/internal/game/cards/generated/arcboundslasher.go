package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Arcbound Slasher", NewArcboundSlasher)
}

// NewArcboundSlasher creates a Arcbound Slasher
// {4}{R} - ARTIFACT CREATURE
func NewArcboundSlasher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arcbound Slasher")
	card.ManaCost = "{4}{R}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"CAT"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
