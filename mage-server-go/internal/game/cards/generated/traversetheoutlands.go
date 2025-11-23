package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Traverse The Outlands", NewTraverseTheOutlands)
}

// NewTraverseTheOutlands creates a Traverse The Outlands
// {4}{G} - SORCERY
func NewTraverseTheOutlands(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Traverse The Outlands")
	card.ManaCost = "{4}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
