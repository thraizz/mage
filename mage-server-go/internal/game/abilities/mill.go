package abilities

import (
	"context"

	"github.com/google/uuid"
)

// MillCardsTargetEffect mills N cards from target player's library to their graveyard
type MillCardsTargetEffect struct {
	amount int
}

// NewMillCardsTargetEffect creates a mill effect targeting a player
func NewMillCardsTargetEffect(amount int) *MillCardsTargetEffect {
	return &MillCardsTargetEffect{
		amount: amount,
	}
}

// Apply executes the mill effect
// TODO: Implement zone management and player.MillCards() method
func (e *MillCardsTargetEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO Phase 1: Get target player from targets
	// TODO Phase 2: Call player.MillCards(amount)
	// TODO Phase 3: Move cards from library to graveyard
	return nil
}

// GetDescription returns a text description of the effect
func (e *MillCardsTargetEffect) GetDescription() string {
	if e.amount == 1 {
		return "target player mills a card"
	}
	return "target player mills " + numberToText(e.amount) + " cards"
}

// MillCardsControllerEffect mills N cards from the controller's library to their graveyard
type MillCardsControllerEffect struct {
	amount int
}

// NewMillCardsControllerEffect creates a mill effect for the controller
func NewMillCardsControllerEffect(amount int) *MillCardsControllerEffect {
	return &MillCardsControllerEffect{
		amount: amount,
	}
}

// Apply executes the mill effect
// TODO: Implement zone management and player.MillCards() method
func (e *MillCardsControllerEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO Phase 1: Get controller from source
	// TODO Phase 2: Call controller.MillCards(amount)
	// TODO Phase 3: Move cards from library to graveyard
	return nil
}

// GetDescription returns a text description of the effect
func (e *MillCardsControllerEffect) GetDescription() string {
	if e.amount == 1 {
		return "mill a card"
	}
	return "mill " + numberToText(e.amount) + " cards"
}

// Helper function to convert numbers to text
func numberToText(n int) string {
	switch n {
	case 1:
		return "one"
	case 2:
		return "two"
	case 3:
		return "three"
	case 4:
		return "four"
	case 5:
		return "five"
	case 6:
		return "six"
	case 7:
		return "seven"
	case 8:
		return "eight"
	case 9:
		return "nine"
	case 10:
		return "ten"
	default:
		return string(rune(n + '0'))
	}
}
