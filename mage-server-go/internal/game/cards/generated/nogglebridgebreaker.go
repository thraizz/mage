package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Noggle Bridgebreaker", NewNoggleBridgebreaker)
}

// NewNoggleBridgebreaker creates a Noggle Bridgebreaker
// {2}{U/R}{U/R} - CREATURE
func NewNoggleBridgebreaker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Noggle Bridgebreaker")
	card.ManaCost = "{2}{U/R}{U/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"NOGGLE", "ROGUE"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
