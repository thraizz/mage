package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Expressive Iteration", NewExpressiveIteration)
}

// NewExpressiveIteration creates a Expressive Iteration
// {U}{R} - SORCERY
func NewExpressiveIteration(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Expressive Iteration")
	card.ManaCost = "{U}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
