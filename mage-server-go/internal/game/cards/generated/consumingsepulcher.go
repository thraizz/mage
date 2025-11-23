package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Consuming Sepulcher", NewConsumingSepulcher)
}

// NewConsumingSepulcher creates a Consuming Sepulcher
//   - ARTIFACT
func NewConsumingSepulcher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Consuming Sepulcher")
	card.ManaCost = ""
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
