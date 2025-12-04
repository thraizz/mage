package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Loamdragger Giant", NewLoamdraggerGiant)
}

// NewLoamdraggerGiant creates a Loamdragger Giant
// {4}{R/G}{R/G}{R/G} - CREATURE
func NewLoamdraggerGiant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Loamdragger Giant")
	card.ManaCost = "{4}{R/G}{R/G}{R/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT", "WARRIOR"}
	card.Power = "7"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
