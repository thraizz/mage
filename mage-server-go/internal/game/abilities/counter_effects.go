package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

// ========================================
// Add Counters Source Effect
// ========================================

// AddCountersSourceEffect adds counters to the source permanent.
// Mirrors Java AddCountersSourceEffect.
type AddCountersSourceEffect struct {
	counter       *counters.Counter
	informPlayers bool
}

// NewAddCountersSourceEffect creates an effect that adds counters to the source permanent.
func NewAddCountersSourceEffect(counter *counters.Counter) *AddCountersSourceEffect {
	return &AddCountersSourceEffect{
		counter:       counter.Copy(),
		informPlayers: false,
	}
}

// NewAddCountersSourceEffectInform creates an effect that adds counters to the source permanent
// and informs players.
func NewAddCountersSourceEffectInform(counter *counters.Counter, informPlayers bool) *AddCountersSourceEffect {
	return &AddCountersSourceEffect{
		counter:       counter.Copy(),
		informPlayers: informPlayers,
	}
}

func (e *AddCountersSourceEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	if e.counter == nil || e.counter.Count <= 0 {
		return nil
	}

	// Cast to CounterGameContext for counter operations
	counterGame, ok := game.(CounterGameContext)
	if !ok {
		return fmt.Errorf("game context does not support counter operations")
	}

	// Get the source permanent
	permanent, err := counterGame.GetPermanent(source)
	if err != nil {
		return fmt.Errorf("failed to get source permanent: %w", err)
	}

	// Add counters to the permanent
	newCounter := e.counter.Copy()
	if err := counterGame.AddCountersToPermanent(permanent, newCounter); err != nil {
		return fmt.Errorf("failed to add counters to source: %w", err)
	}

	if e.informPlayers {
		// Inform players about the counter addition
		counterGame.InformPlayers(fmt.Sprintf("Added %d %s counter(s) to source permanent",
			newCounter.Count, newCounter.Name))
	}

	return nil
}

func (e *AddCountersSourceEffect) GetDescription() string {
	if e.counter == nil {
		return ""
	}
	if e.counter.Count == 1 {
		return fmt.Sprintf("put a %s counter on {this}", e.counter.Name)
	}
	return fmt.Sprintf("put %d %s counters on {this}", e.counter.Count, e.counter.Name)
}

// ========================================
// Add Counters Target Effect
// ========================================

// AddCountersTargetEffect adds counters to target permanents or players.
// Mirrors Java AddCountersTargetEffect.
type AddCountersTargetEffect struct {
	counter *counters.Counter
}

// NewAddCountersTargetEffect creates an effect that adds counters to target permanent(s) or player(s).
func NewAddCountersTargetEffect(counter *counters.Counter) *AddCountersTargetEffect {
	return &AddCountersTargetEffect{
		counter: counter.Copy(),
	}
}

func (e *AddCountersTargetEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	if e.counter == nil || e.counter.Count <= 0 {
		return nil
	}

	if len(targets) == 0 {
		return nil
	}

	// Cast to CounterGameContext for counter operations
	counterGame, ok := game.(CounterGameContext)
	if !ok {
		return fmt.Errorf("game context does not support counter operations")
	}

	affectedTargets := 0
	for _, targetID := range targets {
		newCounter := e.counter.Copy()

		// Try to add to permanent first
		permanent, permErr := counterGame.GetPermanent(targetID)
		if permErr == nil {
			if err := counterGame.AddCountersToPermanent(permanent, newCounter); err != nil {
				return fmt.Errorf("failed to add counters to permanent: %w", err)
			}
			counterGame.InformPlayers(fmt.Sprintf("Put %d %s counter(s) on permanent",
				newCounter.Count, newCounter.Name))
			affectedTargets++
			continue
		}

		// Try to add to player
		player, playerErr := counterGame.GetPlayer(targetID)
		if playerErr == nil {
			if err := counterGame.AddCountersToPlayer(player, newCounter); err != nil {
				return fmt.Errorf("failed to add counters to player: %w", err)
			}
			counterGame.InformPlayers(fmt.Sprintf("Put %d %s counter(s) on player",
				newCounter.Count, newCounter.Name))
			affectedTargets++
			continue
		}

		// Try to add to card (in graveyard, hand, etc.)
		card, cardErr := counterGame.GetCard(targetID)
		if cardErr == nil {
			if err := counterGame.AddCountersToCard(card, newCounter); err != nil {
				return fmt.Errorf("failed to add counters to card: %w", err)
			}
			counterGame.InformPlayers(fmt.Sprintf("Put %d %s counter(s) on card",
				newCounter.Count, newCounter.Name))
			affectedTargets++
			continue
		}

		// If we couldn't find the target as permanent, player, or card, that's an error
		return fmt.Errorf("target %s not found as permanent, player, or card", targetID)
	}

	return nil
}

func (e *AddCountersTargetEffect) GetDescription() string {
	if e.counter == nil {
		return ""
	}
	if e.counter.Count == 1 {
		return fmt.Sprintf("put a %s counter on target", e.counter.Name)
	}
	return fmt.Sprintf("put %d %s counters on target", e.counter.Count, e.counter.Name)
}

// ========================================
// Add Counters All Effect
// ========================================

// AddCountersAllEffect adds counters to all permanents matching a filter.
// Mirrors Java AddCountersAllEffect.
type AddCountersAllEffect struct {
	counter     *counters.Counter
	filterFunc  func(interface{}) bool // Function to filter which permanents get counters
	description string
}

// NewAddCountersAllEffect creates an effect that adds counters to all permanents matching the filter.
func NewAddCountersAllEffect(counter *counters.Counter, filterFunc func(interface{}) bool, description string) *AddCountersAllEffect {
	return &AddCountersAllEffect{
		counter:     counter.Copy(),
		filterFunc:  filterFunc,
		description: description,
	}
}

// NewAddCountersAllCreatures creates an effect that adds counters to all creatures.
func NewAddCountersAllCreatures(counter *counters.Counter) *AddCountersAllEffect {
	return &AddCountersAllEffect{
		counter:     counter.Copy(),
		filterFunc:  nil, // Will be set to creature filter by game context
		description: "each creature",
	}
}

func (e *AddCountersAllEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	if e.counter == nil || e.counter.Count <= 0 {
		return nil
	}

	// Cast to CounterGameContext for counter operations
	counterGame, ok := game.(CounterGameContext)
	if !ok {
		return fmt.Errorf("game context does not support counter operations")
	}

	// Get all permanents from the battlefield
	permanents, err := counterGame.GetAllPermanents()
	if err != nil {
		return fmt.Errorf("failed to get permanents: %w", err)
	}

	affectedCount := 0
	for _, permanent := range permanents {
		// Apply filter if provided
		if e.filterFunc != nil && !e.filterFunc(permanent) {
			continue
		}

		newCounter := e.counter.Copy()
		if err := counterGame.AddCountersToPermanent(permanent, newCounter); err != nil {
			return fmt.Errorf("failed to add counters to permanent: %w", err)
		}
		affectedCount++
	}

	if affectedCount > 0 {
		counterGame.InformPlayers(fmt.Sprintf("Put %d %s counter(s) on %d permanent(s)",
			e.counter.Count, e.counter.Name, affectedCount))
	}

	return nil
}

func (e *AddCountersAllEffect) GetDescription() string {
	if e.counter == nil {
		return ""
	}

	targetDesc := "each permanent"
	if e.description != "" {
		targetDesc = e.description
	}

	if e.counter.Count == 1 {
		return fmt.Sprintf("put a %s counter on %s", e.counter.Name, targetDesc)
	}
	return fmt.Sprintf("put %d %s counters on %s", e.counter.Count, e.counter.Name, targetDesc)
}

// ========================================
// Extended GameContext Interface for Counters
// ========================================

// CounterGameContext extends GameContext with counter-specific methods.
// This should be implemented by the game engine.
type CounterGameContext interface {
	GameContext

	// GetPermanent retrieves a permanent by ID
	GetPermanent(id uuid.UUID) (interface{}, error)

	// AddCountersToPermanent adds counters to a permanent
	AddCountersToPermanent(permanent interface{}, counter *counters.Counter) error

	// AddCountersToPlayer adds counters to a player
	AddCountersToPlayer(player interface{}, counter *counters.Counter) error

	// AddCountersToCard adds counters to a card (in any zone)
	AddCountersToCard(card interface{}, counter *counters.Counter) error

	// GetAllPermanents returns all permanents on the battlefield
	GetAllPermanents() ([]interface{}, error)

	// InformPlayers sends a message to all players
	InformPlayers(message string)
}
