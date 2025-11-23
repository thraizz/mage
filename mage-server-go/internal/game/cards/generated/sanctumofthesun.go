package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sanctum Of The Sun", NewSanctumOfTheSun)
}

// NewSanctumOfTheSun creates a Sanctum Of The Sun
//   - LAND
func NewSanctumOfTheSun(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sanctum Of The Sun")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
