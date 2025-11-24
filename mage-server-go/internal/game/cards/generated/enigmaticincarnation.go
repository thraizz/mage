package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Enigmatic Incarnation", NewEnigmaticIncarnation)
}

// NewEnigmaticIncarnation creates a Enigmatic Incarnation
// {2}{G}{U} - ENCHANTMENT
func NewEnigmaticIncarnation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Enigmatic Incarnation")
	card.ManaCost = "{2}{G}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}