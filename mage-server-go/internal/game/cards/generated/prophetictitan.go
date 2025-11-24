package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Prophetic Titan", NewPropheticTitan)
}

// NewPropheticTitan creates a Prophetic Titan
// {4}{U}{R} - CREATURE
func NewPropheticTitan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Prophetic Titan")
	card.ManaCost = "{4}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT", "WIZARD"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
