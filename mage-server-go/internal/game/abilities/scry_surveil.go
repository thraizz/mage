package abilities

import (
	"context"

	"github.com/google/uuid"
)

// ScryEffect allows the controller to scry N (look at top N cards, put any on bottom, rest on top in any order)
type ScryEffect struct {
	amount         int
	showEffectHint bool
}

// NewScryEffect creates a scry effect
func NewScryEffect(amount int) *ScryEffect {
	return &ScryEffect{
		amount:         amount,
		showEffectHint: true,
	}
}

// NewScryEffectNoHint creates a scry effect without the reminder text hint
func NewScryEffectNoHint(amount int) *ScryEffect {
	return &ScryEffect{
		amount:         amount,
		showEffectHint: false,
	}
}

// Apply executes the scry effect
// TODO: Implement player.Scry() method with UI integration for card ordering
func (e *ScryEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO Phase 1: Get controller from source
	// TODO Phase 2: Look at top N cards of library
	// TODO Phase 3: UI prompt for player to choose cards to put on bottom
	// TODO Phase 4: UI prompt for player to order remaining cards on top
	// TODO Phase 5: Apply the reordering
	return nil
}

// GetDescription returns a text description of the effect
func (e *ScryEffect) GetDescription() string {
	text := "scry " + numberToText(e.amount)

	if e.showEffectHint {
		if e.amount == 1 {
			text += ". (Look at the top card of your library. You may put that card on the bottom.)"
		} else {
			text += ". (Look at the top " + numberToText(e.amount) +
				" cards of your library, then put any number of them on the bottom and the rest on top in any order.)"
		}
	}

	return text
}

// SurveilEffect allows the controller to surveil N (look at top N cards, put any in graveyard, rest on top in any order)
type SurveilEffect struct {
	amount         int
	showEffectHint bool
}

// NewSurveilEffect creates a surveil effect
func NewSurveilEffect(amount int) *SurveilEffect {
	return &SurveilEffect{
		amount:         amount,
		showEffectHint: true,
	}
}

// NewSurveilEffectNoHint creates a surveil effect without the reminder text hint
func NewSurveilEffectNoHint(amount int) *SurveilEffect {
	return &SurveilEffect{
		amount:         amount,
		showEffectHint: false,
	}
}

// Apply executes the surveil effect
// TODO: Implement player.Surveil() method with UI integration for card selection
func (e *SurveilEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO Phase 1: Get controller from source
	// TODO Phase 2: Look at top N cards of library
	// TODO Phase 3: UI prompt for player to choose cards to put in graveyard
	// TODO Phase 4: Move selected cards to graveyard
	// TODO Phase 5: UI prompt for player to order remaining cards on top
	// TODO Phase 6: Apply the reordering
	return nil
}

// GetDescription returns a text description of the effect
func (e *SurveilEffect) GetDescription() string {
	text := "surveil " + numberToText(e.amount)

	if e.showEffectHint {
		if e.amount == 1 {
			text += ". (Look at the top card of your library. You may put that card into your graveyard.)"
		} else {
			text += ". (Look at the top " + numberToText(e.amount) +
				" cards of your library, then put any number of them into your graveyard and the rest on top in any order.)"
		}
	}

	return text
}
