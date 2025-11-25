package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lightning Rift", NewLightningRift)
}

// NewLightningRift creates a Lightning Rift
// {1}{R} - ENCHANTMENT
func NewLightningRift(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lightning Rift")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: CycleAllTriggeredAbility
	//   - Effect: DoIfCostPaid(new DamageTargetEffect(2), new GenericManaCost(1))
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewAnyTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new DamageTargetEffect(2), new GenericManaCost(1))
	// card.AddAbility(ability1)
	return card, nil
}
