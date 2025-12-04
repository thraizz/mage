package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("All Is Dust", NewAllIsDust)
}

// NewAllIsDust creates a All Is Dust
// {7} - KINDRED SORCERY
func NewAllIsDust(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "All Is Dust")
	card.ManaCost = "{7}"
	card.Types = []string{"KINDRED", "SORCERY"}
	card.Subtypes = []string{"ELDRAZI"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
