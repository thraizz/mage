package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Barktooth Warbeard", NewBarktoothWarbeard)
}

// NewBarktoothWarbeard creates a Barktooth Warbeard
// {4}{B}{R}{R} - CREATURE
func NewBarktoothWarbeard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Barktooth Warbeard")
	card.ManaCost = "{4}{B}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
