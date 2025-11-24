package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Strength Bobblehead", NewStrengthBobblehead)
}

// NewStrengthBobblehead creates a Strength Bobblehead
// {3} - ARTIFACT
func NewStrengthBobblehead(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Strength Bobblehead")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"BOBBLEHEAD"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}