package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Frost Giant", NewFrostGiant)
}

// NewFrostGiant creates a Frost Giant
// {3}{R}{R}{R} - CREATURE
func NewFrostGiant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Frost Giant")
	card.ManaCost = "{3}{R}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
