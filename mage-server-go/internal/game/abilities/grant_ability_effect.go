package abilities

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

// TODO: GameContext needs to be enhanced to provide access to layer system
// For now, we'll use the existing GameContext interface from ability.go

// GrantAbilityEffect is an abilities.Effect that grants a keyword ability to targets
// This bridges abilities.Effect (one-shot) with effects.GrantAbilityEffect (continuous)
type GrantAbilityEffect struct {
	abilityID string           // The ability being granted (e.g., "FlyingAbility")
	duration  effects.Duration // How long the ability lasts
}

// NewGrantAbilityEffect creates a new effect that grants an ability
func NewGrantAbilityEffect(abilityID string, duration effects.Duration) *GrantAbilityEffect {
	return &GrantAbilityEffect{
		abilityID: abilityID,
		duration:  duration,
	}
}

// Apply applies the effect by adding a continuous effect to the layer system
func (e *GrantAbilityEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	if len(targets) == 0 {
		return fmt.Errorf("no targets for grant ability effect")
	}

	// Get the layer system from game context
	// For now, we'll assume GameContext has a method to get the effects manager
	// This will need to be implemented in the actual GameContext interface

	// Convert UUID targets to string IDs for effects system
	targetIDs := make([]string, len(targets))
	for i, id := range targets {
		targetIDs[i] = id.String()
	}

	// Create the continuous effect
	continuousEffect := effects.NewGrantAbilityEffect(
		source.String(),
		e.abilityID,
		targetIDs,
		e.duration,
	)

	// Add to layer system (this will need proper integration with GameContext)
	// For now, this is a placeholder that shows the pattern
	_ = continuousEffect

	// TODO: Add continuousEffect to game's layer system
	// This requires GameContext to expose the layer system or effects manager
	// game.GetEffectManager().AddEffect(continuousEffect)

	return nil
}

// GetDescription returns a description of the effect
func (e *GrantAbilityEffect) GetDescription() string {
	// Convert ability ID to readable text (e.g., "FlyingAbility" → "flying")
	abilityName := e.abilityID
	if len(abilityName) > 7 && abilityName[len(abilityName)-7:] == "Ability" {
		abilityName = abilityName[:len(abilityName)-7]
	}
	// Make lowercase for readability
	abilityName = strings.ToLower(abilityName)

	return fmt.Sprintf("gains %s until %s", abilityName, e.duration)
}

// String returns a description of the effect (alias for GetDescription)
func (e *GrantAbilityEffect) String() string {
	return e.GetDescription()
}
