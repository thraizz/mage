package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SpellAbility represents the main ability of an instant or sorcery
type SpellAbility struct {
	baseAbility
	ManaCost *ManaCost
	Effects  []Effect
	Targets  *TargetRequirement
}

// NewSpellAbility creates a new spell ability
func NewSpellAbility(sourceID uuid.UUID, manaCost string, effects []Effect) (*SpellAbility, error) {
	cost, err := ParseManaCost(manaCost)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mana cost: %w", err)
	}

	text := buildEffectString(effects)
	return &SpellAbility{
		baseAbility: newBaseAbility(sourceID, text),
		ManaCost:    cost,
		Effects:     effects,
	}, nil
}

// GetType returns the ability type
func (a *SpellAbility) GetType() AbilityType {
	return AbilityTypeSpell
}

// CanActivate checks if this spell can be cast
func (a *SpellAbility) CanActivate(ctx context.Context, game GameContext) bool {
	// TODO: Check timing restrictions (instant/sorcery speed)
	// TODO: Check if mana cost can be paid
	// TODO: Check if targets are available
	return true
}

// Resolve resolves this spell
func (a *SpellAbility) Resolve(ctx context.Context, game GameContext) error {
	// Get targets (for now, just use empty slice)
	targets := []uuid.UUID{}

	// Apply each effect
	for _, effect := range a.Effects {
		if err := effect.Apply(ctx, game, a.sourceID, targets); err != nil {
			return fmt.Errorf("failed to apply effect: %w", err)
		}
	}

	return nil
}

// SetTargets sets the target requirement for this spell
func (a *SpellAbility) SetTargets(targets *TargetRequirement) {
	a.Targets = targets
	// Rebuild text with target info
	a.text = buildEffectString(a.Effects)
	if targets != nil {
		a.text += " " + targets.Description
	}
}

// GetManaCost returns the mana cost of this spell
func (a *SpellAbility) GetManaCost() *ManaCost {
	return a.ManaCost
}

// GetTargets returns the target requirement for this spell
func (a *SpellAbility) GetTargets() *TargetRequirement {
	return a.Targets
}
