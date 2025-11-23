package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Conclave Tribunal", NewConclaveTribunal)
}

// NewConclaveTribunal creates a Conclave Tribunal
// {3}{W} - ENCHANTMENT
func NewConclaveTribunal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Conclave Tribunal")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
