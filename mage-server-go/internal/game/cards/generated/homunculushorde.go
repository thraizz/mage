package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Homunculus Horde", NewHomunculusHorde)
}

// NewHomunculusHorde creates a Homunculus Horde
// {3}{U} - CREATURE
func NewHomunculusHorde(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Homunculus Horde")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HOMUNCULUS"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
