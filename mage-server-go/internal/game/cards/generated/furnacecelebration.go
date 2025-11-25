package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Furnace Celebration", NewFurnaceCelebration)
}

// NewFurnaceCelebration creates a Furnace Celebration
// {1}{R}{R} - ENCHANTMENT
func NewFurnaceCelebration(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Furnace Celebration")
	card.ManaCost = "{1}{R}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SacrificePermanentTriggeredAbility
	//   - Effect: DoIfCostPaid(                 new DamageTargetEffect(2), new Ge...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewAnyTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
