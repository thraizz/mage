package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Noggle Bandit", NewNoggleBandit)
}

// NewNoggleBandit creates a Noggle Bandit
// {1}{U/R}{U/R} - CREATURE
func NewNoggleBandit(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Noggle Bandit")
	card.ManaCost = "{1}{U/R}{U/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"NOGGLE", "ROGUE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}