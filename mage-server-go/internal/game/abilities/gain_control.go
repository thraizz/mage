package abilities

import (
	"context"

	"github.com/google/uuid"
)

// GainControlTargetEffect allows a player to gain control of target permanent(s)
type GainControlTargetEffect struct {
	duration            Duration
	controllingPlayerID *uuid.UUID
	fixedControl        bool
}

// NewGainControlTargetEffect creates a gain control effect with the specified duration
func NewGainControlTargetEffect(duration Duration) *GainControlTargetEffect {
	return &GainControlTargetEffect{
		duration:     duration,
		fixedControl: false,
	}
}

// NewGainControlTargetEffectFixed creates a gain control effect where the controlling player is fixed
func NewGainControlTargetEffectFixed(duration Duration, fixedControl bool) *GainControlTargetEffect {
	return &GainControlTargetEffect{
		duration:     duration,
		fixedControl: fixedControl,
	}
}

// NewGainControlTargetEffectPlayer creates a gain control effect for a specific player
func NewGainControlTargetEffectPlayer(duration Duration, controllingPlayerID uuid.UUID) *GainControlTargetEffect {
	return &GainControlTargetEffect{
		duration:            duration,
		controllingPlayerID: &controllingPlayerID,
		fixedControl:        true,
	}
}

// Apply executes the gain control effect
// TODO: Implement continuous effect system and permanent.ChangeController() method
func (e *GainControlTargetEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO Phase 1: Get target permanent(s) from targets
	// TODO Phase 2: Determine controlling player (source controller or specified player)
	// TODO Phase 3: Create continuous effect that changes control
	// TODO Phase 4: Add continuous effect to game with specified duration
	// TODO Phase 5: Handle control change events and state-based actions
	return nil
}

// GetDescription returns a text description of the effect
func (e *GainControlTargetEffect) GetDescription() string {
	text := "gain control of target permanent"

	switch e.duration {
	case DurationUntilEndOfTurn:
		text += " until end of turn"
	case DurationUntilEndOfCombat:
		text += " until end of combat"
	case DurationWhileOnBattlefield:
		text += " for as long as you control this permanent"
	}

	return text
}

// GainControlAllEffect allows a player to gain control of all permanents matching a filter
type GainControlAllEffect struct {
	duration            Duration
	filter              TargetFilter
	controllingPlayerID *uuid.UUID
}

// NewGainControlAllEffect creates a gain control all effect
func NewGainControlAllEffect(duration Duration, filter TargetFilter) *GainControlAllEffect {
	return &GainControlAllEffect{
		duration: duration,
		filter:   filter,
	}
}

// NewGainControlAllEffectPlayer creates a gain control all effect for a specific player
func NewGainControlAllEffectPlayer(duration Duration, filter TargetFilter, controllingPlayerID uuid.UUID) *GainControlAllEffect {
	return &GainControlAllEffect{
		duration:            duration,
		filter:              filter,
		controllingPlayerID: &controllingPlayerID,
	}
}

// Apply executes the gain control all effect
// TODO: Implement battlefield filtering and continuous effect creation
func (e *GainControlAllEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO Phase 1: Get all permanents on battlefield matching filter
	// TODO Phase 2: Determine controlling player
	// TODO Phase 3: For each permanent, create a GainControlTargetEffect
	// TODO Phase 4: Add all continuous effects to game
	return nil
}

// GetDescription returns a text description of the effect
func (e *GainControlAllEffect) GetDescription() string {
	text := "gain control of all permanents"

	switch e.duration {
	case DurationUntilEndOfTurn:
		text += " until end of turn"
	case DurationUntilEndOfCombat:
		text += " until end of combat"
	case DurationWhileOnBattlefield:
		text += " for as long as you control this permanent"
	}

	return text
}
