package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Berserkers Of Blood Ridge", NewBerserkersOfBloodRidge)
}

// NewBerserkersOfBloodRidge creates a Berserkers Of Blood Ridge
// {4}{R} - CREATURE
func NewBerserkersOfBloodRidge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Berserkers Of Blood Ridge")
	card.ManaCost = "{4}{R}"
	card.Types = []string{"CREATURE"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
