package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Null Chamber", NewNullChamber)
}

// NewNullChamber creates a Null Chamber
// {3}{W} - ENCHANTMENT
func NewNullChamber(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Null Chamber")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"WORLD"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
