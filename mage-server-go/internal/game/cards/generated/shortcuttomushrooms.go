package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shortcut To Mushrooms", NewShortcutToMushrooms)
}

// NewShortcutToMushrooms creates a Shortcut To Mushrooms
// {1}{G} - ENCHANTMENT
func NewShortcutToMushrooms(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shortcut To Mushrooms")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
