package abilities

import (
	"github.com/google/uuid"
)

// ========================================
// Spell Ability Builder
// ========================================

// SpellAbilityBuilder provides a fluent API for building spell abilities
type SpellAbilityBuilder struct {
	sourceID uuid.UUID
	manaCost string
	effects  []Effect
	targets  *TargetRequirement
}

// NewSpellAbilityBuilder creates a new spell ability builder
func NewSpellAbilityBuilder(sourceID uuid.UUID, manaCost string) *SpellAbilityBuilder {
	return &SpellAbilityBuilder{
		sourceID: sourceID,
		manaCost: manaCost,
		effects:  make([]Effect, 0),
	}
}

// AddEffect adds an effect to this spell
func (b *SpellAbilityBuilder) AddEffect(effect Effect) *SpellAbilityBuilder {
	b.effects = append(b.effects, effect)
	return b
}

// AddTarget sets the target requirement for this spell
func (b *SpellAbilityBuilder) AddTarget(filter TargetFilter) *SpellAbilityBuilder {
	b.targets = NewTargetRequirement(1, 1, filter)
	return b
}

// AddTargets sets the target requirement with custom min/max
func (b *SpellAbilityBuilder) AddTargets(min, max int, filter TargetFilter) *SpellAbilityBuilder {
	b.targets = NewTargetRequirement(min, max, filter)
	return b
}

// Build constructs the spell ability
func (b *SpellAbilityBuilder) Build() (*SpellAbility, error) {
	ability, err := NewSpellAbility(b.sourceID, b.manaCost, b.effects)
	if err != nil {
		return nil, err
	}

	if b.targets != nil {
		ability.SetTargets(b.targets)
	}

	return ability, nil
}

// ========================================
// Activated Ability Builder
// ========================================

// ActivatedAbilityBuilder provides a fluent API for building activated abilities
type ActivatedAbilityBuilder struct {
	sourceID      uuid.UUID
	costs         []Cost
	effects       []Effect
	targets       *TargetRequirement
	timingRule    TimingRule
	isManaAbility bool
}

// NewActivatedAbilityBuilder creates a new activated ability builder
func NewActivatedAbilityBuilder(sourceID uuid.UUID) *ActivatedAbilityBuilder {
	return &ActivatedAbilityBuilder{
		sourceID:      sourceID,
		costs:         make([]Cost, 0),
		effects:       make([]Effect, 0),
		timingRule:    TimingAny,
		isManaAbility: false,
	}
}

// AddCost adds a cost to this ability
func (b *ActivatedAbilityBuilder) AddCost(cost Cost) *ActivatedAbilityBuilder {
	b.costs = append(b.costs, cost)
	return b
}

// AddManaCost adds a mana cost (parsed from string like "{2}{U}")
func (b *ActivatedAbilityBuilder) AddManaCost(costStr string) *ActivatedAbilityBuilder {
	cost, err := ParseManaCost(costStr)
	if err == nil {
		b.costs = append(b.costs, cost)
	}
	return b
}

// AddTapCost adds a tap cost (uses the source permanent)
func (b *ActivatedAbilityBuilder) AddTapCost() *ActivatedAbilityBuilder {
	b.costs = append(b.costs, NewTapCostWithSource(b.sourceID))
	return b
}

// AddSacrificeSourceCost adds a cost to sacrifice the source permanent
func (b *ActivatedAbilityBuilder) AddSacrificeSourceCost() *ActivatedAbilityBuilder {
	b.costs = append(b.costs, NewSacrificeSourceCost())
	return b
}

// AddEffect adds an effect to this ability
func (b *ActivatedAbilityBuilder) AddEffect(effect Effect) *ActivatedAbilityBuilder {
	b.effects = append(b.effects, effect)
	return b
}

// AddTarget sets the target requirement
func (b *ActivatedAbilityBuilder) AddTarget(filter TargetFilter) *ActivatedAbilityBuilder {
	b.targets = NewTargetRequirement(1, 1, filter)
	return b
}

// AddTargets sets the target requirement with custom min/max
func (b *ActivatedAbilityBuilder) AddTargets(min, max int, filter TargetFilter) *ActivatedAbilityBuilder {
	b.targets = NewTargetRequirement(min, max, filter)
	return b
}

// SetTimingRule sets when this ability can be activated
func (b *ActivatedAbilityBuilder) SetTimingRule(rule TimingRule) *ActivatedAbilityBuilder {
	b.timingRule = rule
	return b
}

// SetManaAbility marks this as a mana ability
func (b *ActivatedAbilityBuilder) SetManaAbility() *ActivatedAbilityBuilder {
	b.isManaAbility = true
	return b
}

// Build constructs the activated ability
func (b *ActivatedAbilityBuilder) Build() *ActivatedAbility {
	ability := NewActivatedAbility(b.sourceID, b.costs, b.effects)

	if b.targets != nil {
		ability.SetTargets(b.targets)
	}

	ability.SetTimingRule(b.timingRule)
	ability.SetManaAbility(b.isManaAbility)

	return ability
}

// ========================================
// Convenience Functions
// ========================================

// BuildSimpleManaAbility is a convenience function for building common mana abilities
// Example: BuildSimpleManaAbility(sourceID, "T: Add {G}")
func BuildSimpleManaAbility(sourceID uuid.UUID, color string) *ActivatedAbility {
	mana := NewMana()
	switch color {
	case "W":
		mana.White = 1
	case "U":
		mana.Blue = 1
	case "B":
		mana.Black = 1
	case "R":
		mana.Red = 1
	case "G":
		mana.Green = 1
	case "C":
		mana.Colorless = 1
	default:
		mana.Generic = 1
	}

	return NewActivatedAbilityBuilder(sourceID).
		AddTapCost().
		AddEffect(NewAddManaEffect(mana)).
		SetManaAbility().
		Build()
}

// BuildSimpleDamageAbility is a convenience function for abilities that deal damage
// Example: BuildSimpleDamageAbility(sourceID, "{T}: Deal 1 damage to any target", 1)
func BuildSimpleDamageAbility(sourceID uuid.UUID, amount int) *ActivatedAbility {
	return NewActivatedAbilityBuilder(sourceID).
		AddTapCost().
		AddEffect(NewDamageEffect(amount)).
		AddTarget(NewAnyTargetFilter()).
		Build()
}
