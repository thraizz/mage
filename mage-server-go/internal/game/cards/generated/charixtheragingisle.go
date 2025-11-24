package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Charix The Raging Isle", NewCharixTheRagingIsle)
}

// NewCharixTheRagingIsle creates a Charix The Raging Isle
// {2}{U}{U} - CREATURE
func NewCharixTheRagingIsle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Charix The Raging Isle")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"LEVIATHAN", "CRAB"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "17"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
