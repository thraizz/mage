package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Enchanted Evening", NewEnchantedEvening)
}

// NewEnchantedEvening creates a Enchanted Evening
// {3}{W/U}{W/U} - ENCHANTMENT
func NewEnchantedEvening(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Enchanted Evening")
	card.ManaCost = "{3}{W/U}{W/U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
