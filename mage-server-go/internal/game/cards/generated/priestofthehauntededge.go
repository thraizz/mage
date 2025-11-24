package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Priest Of The Haunted Edge", NewPriestOfTheHauntedEdge)
}

// NewPriestOfTheHauntedEdge creates a Priest Of The Haunted Edge
// {1}{B} - CREATURE
func NewPriestOfTheHauntedEdge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Priest Of The Haunted Edge")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"SNOW"}
	card.Power = "0"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}