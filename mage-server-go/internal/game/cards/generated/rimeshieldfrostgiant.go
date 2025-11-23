package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rimeshield Frost Giant", NewRimeshieldFrostGiant)
}

// NewRimeshieldFrostGiant creates a Rimeshield Frost Giant
// {3}{U}{U} - CREATURE
func NewRimeshieldFrostGiant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rimeshield Frost Giant")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT", "WARRIOR"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
