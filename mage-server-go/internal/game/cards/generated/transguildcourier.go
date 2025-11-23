package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Transguild Courier", NewTransguildCourier)
}

// NewTransguildCourier creates a Transguild Courier
// {4} - ARTIFACT CREATURE
func NewTransguildCourier(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Transguild Courier")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"GOLEM"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
