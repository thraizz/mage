package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Heed The Mists", NewHeedTheMists)
}

// NewHeedTheMists creates a Heed The Mists
// {3}{U}{U} - SORCERY
func NewHeedTheMists(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Heed The Mists")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.Subtypes = []string{"ARCANE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
