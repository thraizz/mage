package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gate To Phyrexia", NewGateToPhyrexia)
}

// NewGateToPhyrexia creates a Gate To Phyrexia
// {B}{B} - ENCHANTMENT
func NewGateToPhyrexia(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gate To Phyrexia")
	card.ManaCost = "{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
