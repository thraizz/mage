package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Copy Enchantment", NewCopyEnchantment)
}

// NewCopyEnchantment creates a Copy Enchantment
// {2}{U} - ENCHANTMENT
func NewCopyEnchantment(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Copy Enchantment")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(StaticFilters.FILTER_PERMANENT_ENCHANTMENT)
	// card.AddAbility(ability0)
	return card, nil
}
