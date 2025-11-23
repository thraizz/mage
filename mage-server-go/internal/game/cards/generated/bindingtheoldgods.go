package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Binding The Old Gods", NewBindingTheOldGods)
}

// NewBindingTheOldGods creates a Binding The Old Gods
// {2}{B}{G} - ENCHANTMENT
func NewBindingTheOldGods(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Binding The Old Gods")
	card.ManaCost = "{2}{B}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
