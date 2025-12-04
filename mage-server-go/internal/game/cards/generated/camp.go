package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("C A M P", NewCAMP)
}

// NewCAMP creates a C A M P
// {3} - ARTIFACT
func NewCAMP(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "C A M P")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"FORTIFICATION"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
