package abilities

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// ActivatedAbility represents an activated ability (cost: effect)
type ActivatedAbility struct {
	baseAbility
	Costs         []Cost
	Effects       []Effect
	Targets       *TargetRequirement
	TimingRule    TimingRule
	UsesStack     bool
	IsManaAbility bool
}

// TimingRule specifies when an ability can be activated
type TimingRule int

const (
	// TimingAny means the ability can be activated any time the player has priority
	TimingAny TimingRule = iota

	// TimingSorcery means the ability can only be activated during main phase with empty stack
	TimingSorcery

	// TimingInstant means the ability can be activated any time the player could cast an instant
	TimingInstant
)

// NewActivatedAbility creates a new activated ability
func NewActivatedAbility(sourceID uuid.UUID, costs []Cost, effects []Effect) *ActivatedAbility {
	text := buildAbilityText(costs, effects)
	return &ActivatedAbility{
		baseAbility:   newBaseAbility(sourceID, text),
		Costs:         costs,
		Effects:       effects,
		TimingRule:    TimingAny,
		UsesStack:     true,
		IsManaAbility: false,
	}
}

// GetType returns the ability type
func (a *ActivatedAbility) GetType() AbilityType {
	if a.IsManaAbility {
		return AbilityTypeMana
	}
	return AbilityTypeActivated
}

// CanActivate checks if this ability can be activated
func (a *ActivatedAbility) CanActivate(ctx context.Context, game GameContext) bool {
	// TODO: Check timing restrictions
	// TODO: Check if costs can be paid
	return true
}

// Resolve resolves this ability
func (a *ActivatedAbility) Resolve(ctx context.Context, game GameContext) error {
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

// buildAbilityText constructs the text representation of an ability
func buildAbilityText(costs []Cost, effects []Effect) string {
	costStr := buildCostString(costs)
	effectStr := buildEffectString(effects)

	if costStr == "" {
		return effectStr
	}

	return fmt.Sprintf("%s: %s", costStr, effectStr)
}

func buildCostString(costs []Cost) string {
	if len(costs) == 0 {
		return ""
	}

	parts := make([]string, len(costs))
	for i, cost := range costs {
		parts[i] = cost.String()
	}

	return strings.Join(parts, ", ")
}

func buildEffectString(effects []Effect) string {
	if len(effects) == 0 {
		return ""
	}

	parts := make([]string, len(effects))
	for i, effect := range effects {
		parts[i] = effect.GetDescription()
	}

	return strings.Join(parts, ". ")
}

// SetTargets sets the target requirement for this ability
func (a *ActivatedAbility) SetTargets(targets *TargetRequirement) {
	a.Targets = targets
	// Rebuild text with target info
	a.text = buildAbilityText(a.Costs, a.Effects)
	if targets != nil {
		a.text += " " + targets.Description
	}
}

// SetTimingRule sets when this ability can be activated
func (a *ActivatedAbility) SetTimingRule(rule TimingRule) {
	a.TimingRule = rule
}

// SetManaAbility marks this as a mana ability
func (a *ActivatedAbility) SetManaAbility(isMana bool) {
	a.IsManaAbility = isMana
	if isMana {
		// Mana abilities don't use the stack
		a.UsesStack = false
	}
}
