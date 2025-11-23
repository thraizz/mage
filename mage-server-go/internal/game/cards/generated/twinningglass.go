package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Twinning Glass", NewTwinningGlass)
}

// NewTwinningGlass creates a Twinning Glass
// {4} - ARTIFACT
func NewTwinningGlass(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Twinning Glass")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
