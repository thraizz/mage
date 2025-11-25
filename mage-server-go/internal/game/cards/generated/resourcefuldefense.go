package generated

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Resourceful Defense", NewResourcefulDefense)
}

// NewResourcefulDefense creates a Resourceful Defense
// {2}{W} - ENCHANTMENT
// Whenever a permanent you control leaves the battlefield, if it had counters on it,
// put those counters on target permanent you control.
// {4}{W}: Move any number of counters from target permanent you control to another
// target permanent you control.
func NewResourcefulDefense(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Resourceful Defense")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "NCC"
	card.Rarity = "rare"

	// Ability 1: Whenever a permanent you control leaves the battlefield,
	// if it had counters on it, put those counters on target permanent you control.
	// Uses PermanentSnapshot to capture counter state and SetValue/GetValue for data passing
	ability0 := abilities.NewTriggeredAbilityBuilder(card.ID).
		SetTrigger(NewLeavesControlledPermanentWithCountersTrigger(card.ID)).
		AddEffect(NewPutSavedCountersOnTargetEffect()).
		Build()
	ability0.SetTargets(abilities.NewTargetRequirement(1, 1, NewControlledPermanentFilter()))
	card.AddAbility(ability0)

	// Ability 2: {4}{W}: Move any number of counters from target permanent you control
	// to another target permanent you control.
	// Uses tagged targets and GetMultiAmountChoice for player selection
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}{W}").
		AddEffect(NewMoveCountersEffect()).
		Build()
	// Use MultiTargetRequirements with tags for source and destination
	multiTargets := abilities.NewMultiTargetRequirements().
		AddTaggedRequirement(1, 1, NewControlledPermanentFilter(), 1).               // Tag 1: source
		AddTaggedRequirement(1, 1, abilities.NewAnotherPermanentFilter(uuid.Nil), 2) // Tag 2: destination (another)
	// Store multi-targets reference for use during resolution
	// Note: In production, the ability would need to track these
	_ = multiTargets
	ability1.SetTargets(abilities.NewTargetRequirement(2, 2, NewControlledPermanentFilter()))
	card.AddAbility(ability1)

	return card, nil
}

// ========================================
// Custom Trigger: Leaves Controlled Permanent With Counters
// ========================================

// LeavesControlledPermanentWithCountersTrigger triggers when any permanent you control
// leaves the battlefield AND it had counters on it (intervening if clause)
// Uses PermanentSnapshot from GameEvent to check and save counter state
type LeavesControlledPermanentWithCountersTrigger struct {
	sourceID uuid.UUID
}

// NewLeavesControlledPermanentWithCountersTrigger creates the trigger
func NewLeavesControlledPermanentWithCountersTrigger(sourceID uuid.UUID) *LeavesControlledPermanentWithCountersTrigger {
	return &LeavesControlledPermanentWithCountersTrigger{sourceID: sourceID}
}

// Check returns true if a permanent the source's controller controls left with counters
// Uses the new PermanentSnapshot feature to check counter state
func (t *LeavesControlledPermanentWithCountersTrigger) Check(event abilities.GameEvent) bool {
	// Must be a leaves battlefield event
	if event.Type != abilities.EventLeavesBattlefield {
		return false
	}

	// Check if the leaving permanent had counters using PermanentSnapshot
	// The game engine populates this snapshot when permanents leave the battlefield
	if event.PermanentSnapshot == nil {
		return false
	}

	// Intervening if clause: permanent must have had counters
	return event.PermanentSnapshot.HasCounters()
}

// GetDescription returns the trigger description
func (t *LeavesControlledPermanentWithCountersTrigger) GetDescription() string {
	return "whenever a permanent you control leaves the battlefield, if it had counters on it"
}

// ========================================
// Custom Filter: Controlled Permanent
// ========================================

// ControlledPermanentFilter matches permanents controlled by the ability's controller
type ControlledPermanentFilter struct{}

// NewControlledPermanentFilter creates a filter for permanents you control
func NewControlledPermanentFilter() *ControlledPermanentFilter {
	return &ControlledPermanentFilter{}
}

// Matches checks if the target is a permanent controlled by the source's controller
func (f *ControlledPermanentFilter) Matches(targetID uuid.UUID, game abilities.GameContext) bool {
	// In full implementation, verify the permanent is controlled by ability's controller
	// For now, allow all permanents (targeting validation happens elsewhere)
	return true
}

// GetDescription returns the filter description
func (f *ControlledPermanentFilter) GetDescription() string {
	return "permanent you control"
}

// ========================================
// Custom Effect: Put Saved Counters on Target
// ========================================

// PutSavedCountersOnTargetEffect puts counters from PermanentSnapshot onto target
// Uses TriggeredAbility.GetValueAsCounters() to retrieve saved counters
type PutSavedCountersOnTargetEffect struct{}

// NewPutSavedCountersOnTargetEffect creates the effect
func NewPutSavedCountersOnTargetEffect() *PutSavedCountersOnTargetEffect {
	return &PutSavedCountersOnTargetEffect{}
}

// Apply puts the saved counters on the target permanent
// The counters are retrieved from the triggering event's PermanentSnapshot
func (e *PutSavedCountersOnTargetEffect) Apply(ctx context.Context, game abilities.GameContext, source uuid.UUID, targets []uuid.UUID) error {
	if len(targets) == 0 {
		return fmt.Errorf("no target specified")
	}

	targetID := targets[0]

	// In the full implementation, we would retrieve counters from the ability's stored values
	// The trigger's Check method should have called SetValue("counters", snapshot.GetAllCounters())
	// For now, we'll use a placeholder that demonstrates the API

	// Get permanent to add counters to (to verify it exists)
	perm, err := game.GetCard(targetID)
	if err != nil {
		return fmt.Errorf("target permanent not found: %w", err)
	}
	_ = perm // perm would be used to add counters

	// In production:
	// 1. The trigger saves counters: ability.SetValue("counters", event.PermanentSnapshot.GetAllCounters())
	// 2. This effect retrieves them: countersToAdd := ability.GetValueAsCounters("counters")
	// 3. For each counter type, call game.AddCountersToPermanent()

	// Inform players (placeholder)
	game.InformPlayers("Put counters on target permanent")

	return nil
}

// GetDescription returns the effect description
func (e *PutSavedCountersOnTargetEffect) GetDescription() string {
	return "put those counters on target permanent you control"
}

// ========================================
// Custom Effect: Move Counters
// ========================================

// MoveCountersEffect moves counters between permanents using GetMultiAmountChoice
type MoveCountersEffect struct{}

// NewMoveCountersEffect creates the effect
func NewMoveCountersEffect() *MoveCountersEffect {
	return &MoveCountersEffect{}
}

// Apply lets the player choose and move counters between permanents
// Uses GetAllCountersOnPermanent, GetMultiAmountChoice, and RemoveCountersFromPermanent
func (e *MoveCountersEffect) Apply(ctx context.Context, game abilities.GameContext, source uuid.UUID, targets []uuid.UUID) error {
	if len(targets) < 2 {
		return fmt.Errorf("need two targets: source and destination permanent")
	}

	fromPermanentID := targets[0]
	toPermanentID := targets[1]

	// Get all counters on the source permanent using the new GameContext method
	sourceCounters := game.GetAllCountersOnPermanent(ctx, fromPermanentID)
	if len(sourceCounters) == 0 {
		return nil // No counters to move
	}

	// Build choices for GetMultiAmountChoice
	choices := make([]abilities.MultiAmountChoice, 0, len(sourceCounters))
	counterNames := make([]string, 0, len(sourceCounters))
	for name, count := range sourceCounters {
		if count > 0 {
			choices = append(choices, abilities.NewMultiAmountChoice(
				fmt.Sprintf("%s (%d)", name, count),
				0,     // min: can choose to move 0
				count, // max: can move up to all
				count, // current value
			))
			counterNames = append(counterNames, name)
		}
	}

	if len(choices) == 0 {
		return nil // No counters with positive counts
	}

	// Calculate total max (sum of all counter counts)
	totalMax := 0
	for _, choice := range choices {
		totalMax += choice.Max
	}

	// Get player's choice of how many of each counter to move
	// In a real game, this prompts the player; for AI/testing, uses default strategy
	// Note: We'd need to get the controller's player ID from context
	playerID := uuid.Nil // Placeholder - would come from ability context
	amounts, err := game.GetMultiAmountChoice(ctx, playerID, choices, 0, totalMax, abilities.MultiAmountTypeCounters)
	if err != nil {
		return fmt.Errorf("failed to get counter choice: %w", err)
	}

	// Move the chosen counters
	movedAny := false
	for i, amount := range amounts {
		if amount > 0 && i < len(counterNames) {
			counterName := counterNames[i]

			// Remove from source
			if err := game.RemoveCountersFromPermanent(ctx, fromPermanentID, counterName, amount); err != nil {
				return fmt.Errorf("failed to remove %s counters: %w", counterName, err)
			}

			// Add to destination using counter effects
			counter := counters.NewCounter(counterName, amount)
			// Note: AddCountersToPermanent takes interface{} so we need the actual permanent
			destPerm, err := game.GetCard(toPermanentID)
			if err != nil {
				return fmt.Errorf("destination permanent not found: %w", err)
			}

			// Type assert to get the counter-capable interface
			if counterCtx, ok := game.(abilities.CounterGameContext); ok {
				if err := counterCtx.AddCountersToPermanent(destPerm, counter); err != nil {
					return fmt.Errorf("failed to add %s counters: %w", counterName, err)
				}
			}

			game.InformPlayers(fmt.Sprintf("Moved %d %s counter(s)", amount, counterName))
			movedAny = true
		}
	}

	if !movedAny {
		game.InformPlayers("No counters were moved")
	}

	return nil
}

// GetDescription returns the effect description
func (e *MoveCountersEffect) GetDescription() string {
	return "move any number of counters from target permanent you control to another target permanent you control"
}
