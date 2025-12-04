package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lightning Cloud", NewLightningCloud)
}

// NewLightningCloud creates a Lightning Cloud
// {3}{R} - ENCHANTMENT
func NewLightningCloud(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lightning Cloud")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SpellCastAllTriggeredAbility
	//   - Effect: DoIfCostPaid(new DamageTargetEffect(1), new ManaCostsImpl<>("{R...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewAnyTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new DamageTargetEffect(1), new ManaCostsImpl<>("{R...)
	// card.AddAbility(ability1)
	return card, nil
}
