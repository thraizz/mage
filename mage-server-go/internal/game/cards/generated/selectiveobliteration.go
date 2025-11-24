package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Selective Obliteration", NewSelectiveObliteration)
}

// NewSelectiveObliteration creates a Selective Obliteration
// {3}{C}{C} - SORCERY
func NewSelectiveObliteration(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Selective Obliteration")
	card.ManaCost = "{3}{C}{C}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
