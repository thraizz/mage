package abilities

import (
	"context"

	"github.com/google/uuid"
)

// Effect represents what an ability does when it resolves
type Effect interface {
	// Apply applies this effect
	Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error

	// GetDescription returns a text description of this effect
	GetDescription() string
}

// Duration represents how long an effect lasts
type Duration int

const (
	// DurationUntilEndOfTurn lasts until end of turn (most common)
	DurationUntilEndOfTurn Duration = iota

	// DurationPermanent lasts forever (until removed)
	DurationPermanent

	// DurationWhileOnBattlefield lasts while the source is on battlefield
	DurationWhileOnBattlefield

	// DurationUntilEndOfCombat lasts until end of combat
	DurationUntilEndOfCombat

	// DurationEndOfTurn alias for DurationUntilEndOfTurn
	DurationEndOfTurn = DurationUntilEndOfTurn

	// DurationEndOfCombat alias for DurationUntilEndOfCombat
	DurationEndOfCombat = DurationUntilEndOfCombat
)
