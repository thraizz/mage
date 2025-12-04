package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wanderbrine Rootcutters", NewWanderbrineRootcutters)
}

// NewWanderbrineRootcutters creates a Wanderbrine Rootcutters
// {2}{U/B}{U/B} - CREATURE
func NewWanderbrineRootcutters(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wanderbrine Rootcutters")
	card.ManaCost = "{2}{U/B}{U/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MERFOLK", "ROGUE"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
