package abilities

import (
	"context"

	"github.com/google/uuid"
)

// ExileTargetEffect exiles target permanent(s) or card(s)
type ExileTargetEffect struct {
	exileZone         string
	exileID           *uuid.UUID
	onlyFromZone      string // Optional zone restriction
	toSourceExileZone bool   // Exile to source-specific zone
}

// NewExileTargetEffect creates a basic exile target effect
func NewExileTargetEffect() *ExileTargetEffect {
	return &ExileTargetEffect{
		toSourceExileZone: false,
	}
}

// NewExileTargetEffectWithZone creates an exile effect with a named exile zone
func NewExileTargetEffectWithZone(exileID uuid.UUID, exileZone string) *ExileTargetEffect {
	return &ExileTargetEffect{
		exileZone:         exileZone,
		exileID:           &exileID,
		toSourceExileZone: false,
	}
}

// NewExileTargetEffectWithText creates an exile effect with custom text
func NewExileTargetEffectWithText(effectText string) *ExileTargetEffect {
	// TODO: Store effectText for GetDescription()
	return NewExileTargetEffect()
}

// Apply executes the exile effect
// TODO: Implement exile zone management and card movement
func (e *ExileTargetEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO Phase 1: Get controller from source
	// TODO Phase 2: For each target, get the card/permanent
	// TODO Phase 3: Check zone restrictions if onlyFromZone is set
	// TODO Phase 4: Determine exile zone (general or source-specific)
	// TODO Phase 5: Move cards to exile zone
	// TODO Phase 6: Handle stack spells and copies specially
	return nil
}

// GetDescription returns a text description of the effect
func (e *ExileTargetEffect) GetDescription() string {
	if e.exileZone != "" {
		return "exile target card (" + e.exileZone + ")"
	}
	return "exile target permanent"
}

// ExileSourceEffect exiles the source permanent
type ExileSourceEffect struct {
	toUniqueExileZone bool
}

// NewExileSourceEffect creates a basic exile source effect
func NewExileSourceEffect() *ExileSourceEffect {
	return &ExileSourceEffect{
		toUniqueExileZone: false,
	}
}

// NewExileSourceEffectUnique creates an exile effect that uses a unique zone per source
func NewExileSourceEffectUnique(toUniqueExileZone bool) *ExileSourceEffect {
	return &ExileSourceEffect{
		toUniqueExileZone: toUniqueExileZone,
	}
}

// Apply executes the exile source effect
// TODO: Implement source card tracking and exile zone management
func (e *ExileSourceEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO Phase 1: Get controller from source
	// TODO Phase 2: Get source card/permanent
	// TODO Phase 3: Verify card still exists and is phased in
	// TODO Phase 4: Determine exile zone (unique or general)
	// TODO Phase 5: Move source to exile
	return nil
}

// GetDescription returns a text description of the effect
func (e *ExileSourceEffect) GetDescription() string {
	return "exile {this}"
}

// ExileAllEffect exiles all permanents matching a filter
type ExileAllEffect struct {
	filter    TargetFilter
	forSource bool // Whether to track exile zone per source
}

// NewExileAllEffect creates an exile all effect with a filter
func NewExileAllEffect(filter TargetFilter) *ExileAllEffect {
	return &ExileAllEffect{
		filter:    filter,
		forSource: false,
	}
}

// NewExileAllEffectForSource creates an exile all effect with source-specific zone
func NewExileAllEffectForSource(filter TargetFilter, forSource bool) *ExileAllEffect {
	return &ExileAllEffect{
		filter:    filter,
		forSource: forSource,
	}
}

// Apply executes the exile all effect
// TODO: Implement battlefield filtering and mass exile
func (e *ExileAllEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO Phase 1: Get controller from source
	// TODO Phase 2: Get all permanents on battlefield matching filter
	// TODO Phase 3: Determine exile zone (source-specific or general)
	// TODO Phase 4: Move all matching permanents to exile
	return nil
}

// GetDescription returns a text description of the effect
func (e *ExileAllEffect) GetDescription() string {
	return "exile all permanents"
}
