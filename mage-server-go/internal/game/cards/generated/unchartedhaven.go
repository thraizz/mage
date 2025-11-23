package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Uncharted Haven", NewUnchartedHaven)
}

// NewUnchartedHaven creates a Uncharted Haven
//   - LAND
func NewUnchartedHaven(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Uncharted Haven")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
