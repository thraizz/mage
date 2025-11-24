package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Urabrasks Forge", NewUrabrasksForge)
}

// NewUrabrasksForge creates a Urabrasks Forge
// {2}{R} - ARTIFACT
func NewUrabrasksForge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Urabrasks Forge")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
