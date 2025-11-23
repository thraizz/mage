package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Azami Lady Of Scrolls", NewAzamiLadyOfScrolls)
}

// NewAzamiLadyOfScrolls creates a Azami Lady Of Scrolls
// {2}{U}{U}{U} - CREATURE
func NewAzamiLadyOfScrolls(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Azami Lady Of Scrolls")
	card.ManaCost = "{2}{U}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
