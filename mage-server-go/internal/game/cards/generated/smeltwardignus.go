package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Smelt Ward Ignus", NewSmeltWardIgnus)
}

// NewSmeltWardIgnus creates a Smelt Ward Ignus
// {1}{R} - CREATURE
func NewSmeltWardIgnus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Smelt Ward Ignus")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
