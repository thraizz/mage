package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flame Blitz", NewFlameBlitz)
}

// NewFlameBlitz creates a Flame Blitz
// {R} - ENCHANTMENT
func NewFlameBlitz(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flame Blitz")
	card.ManaCost = "{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(                 5, StaticFilters.FILTER_PERMANENT...)
	// card.AddAbility(ability0)
	return card, nil
}
