package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

// Rule 701.27: Proliferate
// To proliferate means to choose any number of permanents and/or players that have
// a counter, then give each one additional counter of each kind that permanent or
// player already has.

// ProliferateEffect represents the Proliferate keyword action
type ProliferateEffect struct {
	description string
}

// NewProliferateEffect creates a new Proliferate effect
func NewProliferateEffect() *ProliferateEffect {
	return &ProliferateEffect{
		description: "Proliferate",
	}
}

// Apply implements the Effect interface
func (e *ProliferateEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. Identify all permanents and players with counters
	// 2. Let controller choose any number of them
	// 3. For each chosen permanent/player:
	//    - Identify all counter types on it
	//    - Add one counter of each type
	//
	// Examples:
	// - Permanent with two +1/+1 counters → gets one more +1/+1 counter
	// - Player with 3 poison counters → gets one more poison counter
	// - Permanent with +1/+1 counter AND -1/-1 counter → gets one of each
	// - Planeswalker with 4 loyalty counters → gets one more loyalty counter

	// Get counter game context
	counterGame, ok := game.(CounterGameContext)
	if !ok {
		return fmt.Errorf("game context does not support counter operations")
	}

	// Get controller
	// TODO: Get controller from source permanent
	_ = counterGame // Use counter game context when implemented

	// Phase 1: Collect all permanents and players with counters
	// TODO: Query game state for permanents with counters
	// TODO: Query game state for players with counters (poison, energy, experience, etc.)

	// Phase 2: Let player choose which ones to proliferate
	// TODO: UI interaction for player choice

	// Phase 3: For each chosen target, add one counter of each type
	// TODO: Iterate through chosen permanents/players
	// TODO: Get counter types on each
	// TODO: Add one counter of each type

	return fmt.Errorf("proliferate effect not yet fully implemented in game context")
}

// GetDescription returns a text description of the effect
func (e *ProliferateEffect) GetDescription() string {
	return e.description + ". (Choose any number of permanents and/or players, then give each another counter of each kind already there.)"
}

// ProliferateHelper provides helper methods for proliferate implementation
type ProliferateHelper struct {
	game CounterGameContext
}

// NewProliferateHelper creates a helper for proliferate operations
func NewProliferateHelper(game CounterGameContext) *ProliferateHelper {
	return &ProliferateHelper{
		game: game,
	}
}

// GetPermanentsWithCounters returns all permanents that have counters
func (ph *ProliferateHelper) GetPermanentsWithCounters() []uuid.UUID {
	// Query game state for all permanents with any counter
	// Return list of permanent IDs
	return []uuid.UUID{} // Placeholder
}

// GetPlayersWithCounters returns all players that have counters
func (ph *ProliferateHelper) GetPlayersWithCounters() []uuid.UUID {
	// Query game state for all players with counters
	// (Poison counters, Energy counters, Experience counters, etc.)
	return []uuid.UUID{} // Placeholder
}

// GetCounterTypesOnPermanent returns all counter types on a permanent
func (ph *ProliferateHelper) GetCounterTypesOnPermanent(permanentID uuid.UUID) ([]*counters.Counter, error) {
	// Query permanent's counters and return list of types
	// Each counter should have Count=1 for proliferate (add one of each type)
	return []*counters.Counter{}, nil // Placeholder
}

// GetCounterTypesOnPlayer returns all counter types on a player
func (ph *ProliferateHelper) GetCounterTypesOnPlayer(playerID uuid.UUID) ([]*counters.Counter, error) {
	// Query player's counters and return list of types
	// Each counter should have Count=1 for proliferate
	return []*counters.Counter{}, nil // Placeholder
}

// AddCounterToPermanent adds one counter of the specified type to a permanent
func (ph *ProliferateHelper) AddCounterToPermanent(permanentID uuid.UUID, counter *counters.Counter) error {
	// Add one counter to the permanent
	permanent, err := ph.game.GetPermanent(permanentID)
	if err != nil {
		return fmt.Errorf("failed to get permanent: %w", err)
	}

	// Add counter with Count=1
	counterToAdd := counters.NewCounter(counter.Name, 1)
	if err := ph.game.AddCountersToPermanent(permanent, counterToAdd); err != nil {
		return fmt.Errorf("failed to add counter: %w", err)
	}

	return nil
}

// AddCounterToPlayer adds one counter of the specified type to a player
func (ph *ProliferateHelper) AddCounterToPlayer(playerID uuid.UUID, counter *counters.Counter) error {
	// Add one counter to the player
	// TODO: Implement player counter system
	return fmt.Errorf("player counters not yet implemented")
}

// ProliferateTargets performs proliferate on the chosen targets
func (ph *ProliferateHelper) ProliferateTargets(permanents []uuid.UUID, players []uuid.UUID) error {
	// Proliferate chosen permanents
	for _, permanentID := range permanents {
		// Get all counter types on this permanent
		counters, err := ph.GetCounterTypesOnPermanent(permanentID)
		if err != nil {
			return fmt.Errorf("failed to get counters on permanent %s: %w", permanentID, err)
		}

		// Add one counter of each type
		for _, counter := range counters {
			if err := ph.AddCounterToPermanent(permanentID, counter); err != nil {
				return fmt.Errorf("failed to add counter to permanent %s: %w", permanentID, err)
			}
		}
	}

	// Proliferate chosen players
	for _, playerID := range players {
		// Get all counter types on this player
		counters, err := ph.GetCounterTypesOnPlayer(playerID)
		if err != nil {
			return fmt.Errorf("failed to get counters on player %s: %w", playerID, err)
		}

		// Add one counter of each type
		for _, counter := range counters {
			if err := ph.AddCounterToPlayer(playerID, counter); err != nil {
				return fmt.Errorf("failed to add counter to player %s: %w", playerID, err)
			}
		}
	}

	return nil
}

// ===== Common Proliferate Card Examples =====

// Example cards using Proliferate:
// - Viral Drake: "{3}{U}, {T}: Proliferate"
// - Contagion Engine: "{6}, Artifact, "{4}, {T}: Proliferate, then proliferate again"
// - Tezzeret's Gambit: "{1}{U/P}{U/P}, Sorcery, "Draw two cards, then proliferate"
// - Karn's Bastion: Land, "{T}: Add {C}. {4}, {T}: Proliferate"
// - Steady Progress: "{2}{U}, Instant, "Scry 1, then proliferate"
// - Grateful Apparition: "{1}{W}, Flying, Whenever ~ deals combat damage to a player or planeswalker, proliferate
// - Flux Channeler: "{2}{U}, Whenever you cast a noncreature spell, proliferate
