package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ========================================
// Condition Interface
// ========================================

// Condition represents a condition that can be checked at runtime
// Java: mage.abilities.condition.Condition
type Condition interface {
	// Check returns true if the condition is met
	Check(ctx context.Context, game GameContext, source uuid.UUID) bool

	// GetDescription returns a text description of this condition
	GetDescription() string
}

// ========================================
// Conditional Effect
// ========================================

// ConditionalEffect applies an effect only if a condition is met
// Java: mage.abilities.effects.common.conditional.ConditionalOneShotEffect
type ConditionalEffect struct {
	effect          Effect    // The effect to apply if condition is met
	condition       Condition // The condition to check
	conditionText   string    // Text representation of condition
	otherwiseText   string    // Text for what happens if condition not met
	otherwiseEffect Effect    // Optional effect if condition not met
}

// NewConditionalEffect creates a new conditional effect
// effect: The effect to apply if condition is true
// conditionText: A string describing the condition (for cases where we don't have a Condition object)
// Java: new ConditionalOneShotEffect(effect, condition, ruleText)
func NewConditionalEffect(effect Effect, conditionText string) *ConditionalEffect {
	return &ConditionalEffect{
		effect:        effect,
		condition:     nil,
		conditionText: conditionText,
	}
}

// NewConditionalEffectWithCondition creates a conditional effect with a proper condition
func NewConditionalEffectWithCondition(effect Effect, condition Condition) *ConditionalEffect {
	return &ConditionalEffect{
		effect:        effect,
		condition:     condition,
		conditionText: condition.GetDescription(),
	}
}

// NewConditionalEffectWithOtherwise creates a conditional effect with an else clause
func NewConditionalEffectWithOtherwise(effect Effect, condition Condition, otherwiseEffect Effect) *ConditionalEffect {
	return &ConditionalEffect{
		effect:          effect,
		condition:       condition,
		conditionText:   condition.GetDescription(),
		otherwiseEffect: otherwiseEffect,
	}
}

// Apply applies the effect if the condition is met
func (e *ConditionalEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	conditionMet := false

	if e.condition != nil {
		conditionMet = e.condition.Check(ctx, game, source)
	} else {
		// If no condition object, assume condition is met (condition text is informational only)
		// This handles cases where the condition isn't implemented yet
		conditionMet = true
	}

	if conditionMet && e.effect != nil {
		return e.effect.Apply(ctx, game, source, targets)
	} else if !conditionMet && e.otherwiseEffect != nil {
		return e.otherwiseEffect.Apply(ctx, game, source, targets)
	}

	return nil
}

// GetDescription returns a description of the conditional effect
func (e *ConditionalEffect) GetDescription() string {
	effectDesc := ""
	if e.effect != nil {
		effectDesc = e.effect.GetDescription()
	}

	if e.conditionText != "" {
		return fmt.Sprintf("if %s, %s", e.conditionText, effectDesc)
	}
	return effectDesc
}

// ========================================
// Common Conditions
// ========================================

// StaticCondition always returns the same value
type StaticCondition struct {
	value       bool
	description string
}

// NewStaticCondition creates a condition that always returns the given value
func NewStaticCondition(value bool, description string) *StaticCondition {
	return &StaticCondition{value: value, description: description}
}

func (c *StaticCondition) Check(ctx context.Context, game GameContext, source uuid.UUID) bool {
	return c.value
}

func (c *StaticCondition) GetDescription() string {
	return c.description
}

// SourceTappedCondition checks if the source is tapped
// Java: mage.abilities.condition.common.SourceTappedCondition
type SourceTappedCondition struct {
	inverted bool // if true, checks for NOT tapped
}

var SourceTappedConditionInstance = &SourceTappedCondition{inverted: false}
var SourceUntappedConditionInstance = &SourceTappedCondition{inverted: true}

func (c *SourceTappedCondition) Check(ctx context.Context, game GameContext, source uuid.UUID) bool {
	// TODO: Check if source is tapped
	// tapped := game.IsTapped(source)
	// return tapped != c.inverted
	return false
}

func (c *SourceTappedCondition) GetDescription() string {
	if c.inverted {
		return "source is untapped"
	}
	return "source is tapped"
}

// ControlsCreatureCondition checks if the player controls a creature
// Java: various creature control conditions
type ControlsCreatureCondition struct {
	withPower     *int   // optional: minimum power requirement
	withToughness *int   // optional: minimum toughness requirement
	subtype       string // optional: specific creature subtype
}

func NewControlsCreatureCondition() *ControlsCreatureCondition {
	return &ControlsCreatureCondition{}
}

func NewControlsCreatureConditionWithPower(minPower int) *ControlsCreatureCondition {
	return &ControlsCreatureCondition{withPower: &minPower}
}

func NewControlsCreatureConditionWithSubtype(subtype string) *ControlsCreatureCondition {
	return &ControlsCreatureCondition{subtype: subtype}
}

func (c *ControlsCreatureCondition) Check(ctx context.Context, game GameContext, source uuid.UUID) bool {
	// TODO: Implement creature control check
	return false
}

func (c *ControlsCreatureCondition) GetDescription() string {
	if c.withPower != nil {
		return fmt.Sprintf("you control a creature with power %d or greater", *c.withPower)
	}
	if c.subtype != "" {
		return fmt.Sprintf("you control a %s", c.subtype)
	}
	return "you control a creature"
}

// GreatestPowerCondition checks if a creature has the greatest power among creatures
// Java: mage.abilities.condition.common.GreatestPowerCondition
type GreatestPowerCondition struct{}

var GreatestPowerConditionInstance = &GreatestPowerCondition{}

func (c *GreatestPowerCondition) Check(ctx context.Context, game GameContext, source uuid.UUID) bool {
	// TODO: Check if source has greatest power
	return false
}

func (c *GreatestPowerCondition) GetDescription() string {
	return "this creature has the greatest power or is tied for greatest power"
}

// InvertCondition inverts another condition
type InvertCondition struct {
	inner Condition
}

func NewInvertCondition(inner Condition) *InvertCondition {
	return &InvertCondition{inner: inner}
}

func (c *InvertCondition) Check(ctx context.Context, game GameContext, source uuid.UUID) bool {
	return !c.inner.Check(ctx, game, source)
}

func (c *InvertCondition) GetDescription() string {
	return fmt.Sprintf("not %s", c.inner.GetDescription())
}

// AndCondition combines multiple conditions with AND logic
type AndCondition struct {
	conditions []Condition
}

func NewAndCondition(conditions ...Condition) *AndCondition {
	return &AndCondition{conditions: conditions}
}

func (c *AndCondition) Check(ctx context.Context, game GameContext, source uuid.UUID) bool {
	for _, cond := range c.conditions {
		if !cond.Check(ctx, game, source) {
			return false
		}
	}
	return true
}

func (c *AndCondition) GetDescription() string {
	if len(c.conditions) == 0 {
		return "true"
	}
	if len(c.conditions) == 1 {
		return c.conditions[0].GetDescription()
	}
	desc := c.conditions[0].GetDescription()
	for i := 1; i < len(c.conditions); i++ {
		desc += " and " + c.conditions[i].GetDescription()
	}
	return desc
}

// OrCondition combines multiple conditions with OR logic
type OrCondition struct {
	conditions []Condition
}

func NewOrCondition(conditions ...Condition) *OrCondition {
	return &OrCondition{conditions: conditions}
}

func (c *OrCondition) Check(ctx context.Context, game GameContext, source uuid.UUID) bool {
	for _, cond := range c.conditions {
		if cond.Check(ctx, game, source) {
			return true
		}
	}
	return false
}

func (c *OrCondition) GetDescription() string {
	if len(c.conditions) == 0 {
		return "false"
	}
	if len(c.conditions) == 1 {
		return c.conditions[0].GetDescription()
	}
	desc := c.conditions[0].GetDescription()
	for i := 1; i < len(c.conditions); i++ {
		desc += " or " + c.conditions[i].GetDescription()
	}
	return desc
}

// ========================================
// Conditional Continuous Effect
// ========================================

// ConditionalContinuousEffect is a continuous effect that only applies if a condition is met
// Java: mage.abilities.effects.common.continuous.ConditionalContinuousEffect
type ConditionalContinuousEffect struct {
	baseContinuousEffect
	effect    Effect    // The continuous effect to apply
	condition Condition // The condition to check
}

// NewConditionalContinuousEffect creates a new conditional continuous effect
func NewConditionalContinuousEffect(effect Effect, condition Condition, duration Duration) *ConditionalContinuousEffect {
	return &ConditionalContinuousEffect{
		baseContinuousEffect: baseContinuousEffect{
			layer:    LayerOther,
			duration: duration,
		},
		effect:    effect,
		condition: condition,
	}
}

// Apply applies the continuous effect if condition is met
func (e *ConditionalContinuousEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	if e.condition != nil && !e.condition.Check(ctx, game, source) {
		return nil // Condition not met, don't apply effect
	}

	if e.effect != nil {
		return e.effect.Apply(ctx, game, source, targets)
	}
	return nil
}

// GetDescription returns a description of the conditional continuous effect
func (e *ConditionalContinuousEffect) GetDescription() string {
	effectDesc := ""
	if e.effect != nil {
		effectDesc = e.effect.GetDescription()
	}

	condDesc := ""
	if e.condition != nil {
		condDesc = e.condition.GetDescription()
	}

	if condDesc != "" {
		return fmt.Sprintf("as long as %s, %s", condDesc, effectDesc)
	}
	return effectDesc
}
