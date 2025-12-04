package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ========================================
// Kicker Ability
// ========================================

// KickerAbility represents the "Kicker {cost}" ability
// Java: mage.abilities.keyword.KickerAbility
// MTG Rules: 702.32 (Kicker)
//
// Kicker is a static ability that represents an optional additional cost.
// A card is "kicked" if its kicker cost was paid during casting.
type KickerAbility struct {
	baseAbility
	kickerCost *ManaCost // The additional cost to kick
	kicked     bool      // Whether the spell was kicked when cast
}

// NewKickerAbility creates a new kicker ability
// cost is the kicker mana cost (e.g., "{2}" for "Kicker {2}")
// Java: new KickerAbility("{2}")
func NewKickerAbility(sourceID uuid.UUID, cost string) *KickerAbility {
	manaCost, err := ParseManaCost(cost)
	if err != nil {
		// If cost parsing fails, create with zero cost
		manaCost = &ManaCost{}
	}

	text := fmt.Sprintf("Kicker %s", cost)
	return &KickerAbility{
		baseAbility: newBaseAbility(sourceID, text),
		kickerCost:  manaCost,
		kicked:      false,
	}
}

// GetType returns the ability type
func (a *KickerAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

// CanActivate always returns true (kicker decision is made during casting)
func (a *KickerAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

// Resolve does nothing (kicker is handled during spell casting)
func (a *KickerAbility) Resolve(ctx context.Context, game GameContext) error {
	return nil
}

// GetKickerCost returns the kicker cost
func (a *KickerAbility) GetKickerCost() *ManaCost {
	return a.kickerCost
}

// IsKicked returns whether the spell was kicked
func (a *KickerAbility) IsKicked() bool {
	return a.kicked
}

// SetKicked marks the spell as kicked
func (a *KickerAbility) SetKicked(kicked bool) {
	a.kicked = kicked
}

// ========================================
// Multikicker Ability
// ========================================

// MultikickerAbility represents the "Multikicker {cost}" ability
// Java: mage.abilities.keyword.MultikickerAbility
// MTG Rules: 702.32 (Kicker - Multikicker)
//
// Multikicker is a variant of kicker that can be paid any number of times.
type MultikickerAbility struct {
	KickerAbility
	kickCount int // Number of times the spell was kicked
}

// NewMultikickerAbility creates a new multikicker ability
func NewMultikickerAbility(sourceID uuid.UUID, cost string) *MultikickerAbility {
	manaCost, err := ParseManaCost(cost)
	if err != nil {
		manaCost = &ManaCost{}
	}

	text := fmt.Sprintf("Multikicker %s", cost)
	return &MultikickerAbility{
		KickerAbility: KickerAbility{
			baseAbility: newBaseAbility(sourceID, text),
			kickerCost:  manaCost,
			kicked:      false,
		},
		kickCount: 0,
	}
}

// GetKickCount returns the number of times the spell was kicked
func (a *MultikickerAbility) GetKickCount() int {
	return a.kickCount
}

// SetKickCount sets the number of times the spell was kicked
func (a *MultikickerAbility) SetKickCount(count int) {
	a.kickCount = count
	a.kicked = count > 0
}

// ========================================
// Kicked Condition
// ========================================

// KickedCondition represents the condition "if this spell was kicked"
// Java: mage.abilities.condition.common.KickedCondition
type KickedCondition struct{}

// KickedConditionInstance is the singleton instance for checking if a spell was kicked
var KickedConditionInstance = &KickedCondition{}

// Check returns true if the spell was kicked
func (c *KickedCondition) Check(ctx context.Context, game GameContext, source uuid.UUID) bool {
	// TODO: Look up the spell and check if it was kicked
	// This requires accessing the spell on the stack or the card's kicker ability
	return false
}

// GetDescription returns the condition description
func (c *KickedCondition) GetDescription() string {
	return "this spell was kicked"
}
