package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Automated Assembly Line", NewAutomatedAssemblyLine)
}

// NewAutomatedAssemblyLine creates a Automated Assembly Line
// {1}{W} - ARTIFACT
func NewAutomatedAssemblyLine(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Automated Assembly Line")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}