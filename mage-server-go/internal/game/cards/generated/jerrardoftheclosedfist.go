package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jerrard Of The Closed Fist", NewJerrardOfTheClosedFist)
}

// NewJerrardOfTheClosedFist creates a Jerrard Of The Closed Fist
// {3}{R}{G}{G} - CREATURE
func NewJerrardOfTheClosedFist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jerrard Of The Closed Fist")
	card.ManaCost = "{3}{R}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "KNIGHT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}