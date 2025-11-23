package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Arcanis The Omnipotent", NewArcanisTheOmnipotent)
}

// NewArcanisTheOmnipotent creates a Arcanis The Omnipotent
// {3}{U}{U}{U} - CREATURE
func NewArcanisTheOmnipotent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arcanis The Omnipotent")
	card.ManaCost = "{3}{U}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
