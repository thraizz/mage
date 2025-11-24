package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Endurance Bobblehead", NewEnduranceBobblehead)
}

// NewEnduranceBobblehead creates a Endurance Bobblehead
// {3} - ARTIFACT
func NewEnduranceBobblehead(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Endurance Bobblehead")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"BOBBLEHEAD"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}