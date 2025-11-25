package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Predatory Urge", NewPredatoryUrge)
}

// NewPredatoryUrge creates a Predatory Urge
// {3}{G} - ENCHANTMENT
func NewPredatoryUrge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Predatory Urge")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DamageEachOtherEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())
	// card.AddAbility(ability1)
	// TODO: Implement activated ability with unmapped effects
	//   - DamageEachOtherEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability2)
	return card, nil
}
