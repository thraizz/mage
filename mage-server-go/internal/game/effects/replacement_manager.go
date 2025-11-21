package effects

import (
	"fmt"
	"sync"

	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap"
)

// ReplacementManager manages all active replacement and prevention effects in a game.
// Implements the replacement effect resolution algorithm from Rule 616.
//
// Key responsibilities:
// - Track all active replacement effects
// - Apply replacement effects to events in the correct order
// - Handle multiple replacement effects on the same event
// - Prevent effects from applying twice to the same event
// - Support self-replacement effect priority (Rule 614.15)
type ReplacementManager struct {
	mu      sync.RWMutex
	effects map[string]ReplacementEffect // All active replacement effects by ID
	logger  *zap.Logger
}

// NewReplacementManager creates a new replacement effect manager
func NewReplacementManager(logger *zap.Logger) *ReplacementManager {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &ReplacementManager{
		effects: make(map[string]ReplacementEffect),
		logger:  logger,
	}
}

// AddEffect adds a replacement effect to the manager
func (rm *ReplacementManager) AddEffect(effect ReplacementEffect) {
	if effect == nil {
		rm.logger.Warn("attempted to add nil replacement effect")
		return
	}

	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.effects[effect.ID()] = effect

	rm.logger.Debug("added replacement effect",
		zap.String("effect_id", effect.ID()),
		zap.String("source_id", effect.SourceID()),
		zap.Bool("self_replacement", effect.IsSelfReplacement()),
		zap.Bool("self_scope", effect.HasSelfScope()))
}

// RemoveEffect removes a replacement effect from the manager
func (rm *ReplacementManager) RemoveEffect(effectID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	delete(rm.effects, effectID)

	rm.logger.Debug("removed replacement effect",
		zap.String("effect_id", effectID))
}

// GetEffect retrieves a replacement effect by ID
func (rm *ReplacementManager) GetEffect(effectID string) (ReplacementEffect, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	effect, ok := rm.effects[effectID]
	return effect, ok
}

// GetEffects returns all active replacement effects
func (rm *ReplacementManager) GetEffects() []ReplacementEffect {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	effects := make([]ReplacementEffect, 0, len(rm.effects))
	for _, effect := range rm.effects {
		effects = append(effects, effect)
	}

	return effects
}

// ClearEffects removes all replacement effects
func (rm *ReplacementManager) ClearEffects() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.effects = make(map[string]ReplacementEffect)

	rm.logger.Debug("cleared all replacement effects")
}

// CleanupExpiredEffects removes effects that have expired
// This should be called at appropriate times (e.g., end of turn, end of combat)
func (rm *ReplacementManager) CleanupExpiredEffects(currentDuration Duration) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	removed := 0
	for id, effect := range rm.effects {
		// Check if effect has expired based on duration
		// For now, we'll implement simple duration-based cleanup
		// More complex duration tracking (until end of turn, etc.) would require game state
		if shouldRemoveEffect(effect, currentDuration) {
			delete(rm.effects, id)
			removed++
		}
	}

	if removed > 0 {
		rm.logger.Debug("cleaned up expired replacement effects",
			zap.Int("removed", removed))
	}
}

// shouldRemoveEffect determines if an effect should be removed
// This is a simplified implementation - full implementation would require game state
func shouldRemoveEffect(effect ReplacementEffect, currentDuration Duration) bool {
	// OneUse effects are removed after first application (tracked externally)
	if effect.Duration() == DurationOneUse {
		return false // Don't auto-remove, let caller handle it
	}

	// Permanent effects never expire
	if effect.Duration() == DurationPermanent {
		return false
	}

	// For other durations, we'd need game state to determine expiration
	// This is a placeholder for future implementation
	return false
}

// ReplaceEvent applies all applicable replacement effects to an event.
// Implements the replacement effect resolution algorithm from Rule 616.
//
// The algorithm:
// 1. Find all replacement effects that could apply to this event
// 2. Separate self-replacement effects (applied first per Rule 614.15)
// 3. Apply one replacement effect at a time
// 4. Track which effects have already been applied to prevent double-application
// 5. Repeat until no more effects can apply
//
// Returns the modified event (which may be completely replaced/prevented)
func (rm *ReplacementManager) ReplaceEvent(event rules.Event, gameID string, choosingPlayerID string) rules.Event {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	// Track which effects have already been applied to prevent double-application (Rule 614.5)
	appliedEffects := make(map[string]bool)
	if event.AppliedEffects != nil {
		for _, effectID := range event.AppliedEffects {
			appliedEffects[effectID] = true
		}
	}

	// Keep applying effects until none are applicable
	maxIterations := 100 // Safety limit to prevent infinite loops
	iteration := 0

	for iteration < maxIterations {
		iteration++

		// Find all applicable effects for this event
		applicable := rm.findApplicableEffects(event, gameID, appliedEffects)

		if len(applicable) == 0 {
			// No more effects to apply
			break
		}

		// Separate self-replacement effects (they must be applied first per Rule 614.15)
		selfReplacementEffects := make([]ReplacementEffect, 0)
		otherEffects := make([]ReplacementEffect, 0)

		for _, effect := range applicable {
			if effect.IsSelfReplacement() {
				selfReplacementEffects = append(selfReplacementEffects, effect)
			} else {
				otherEffects = append(otherEffects, effect)
			}
		}

		// Choose which effect to apply
		var chosenEffect ReplacementEffect

		if len(selfReplacementEffects) > 0 {
			// Must choose a self-replacement effect first (Rule 616.1a)
			if len(selfReplacementEffects) == 1 {
				chosenEffect = selfReplacementEffects[0]
			} else {
				// Multiple self-replacement effects - choosing player selects
				// For now, just take the first one
				// TODO: Implement player choice mechanism
				chosenEffect = selfReplacementEffects[0]
			}
		} else {
			// No self-replacement effects, choose from other effects
			if len(otherEffects) == 1 {
				chosenEffect = otherEffects[0]
			} else {
				// Multiple effects - choosing player selects (Rule 616.1)
				// The affected object's controller or affected player chooses
				// For now, just take the first one
				// TODO: Implement player choice mechanism
				chosenEffect = otherEffects[0]
			}
		}

		// Apply the chosen effect
		replacedEvent, completelyReplaced := chosenEffect.ReplaceEvent(event, gameID)
		event = replacedEvent

		// Mark this effect as applied
		appliedEffects[chosenEffect.ID()] = true
		if event.AppliedEffects == nil {
			event.AppliedEffects = make([]string, 0)
		}
		event.AppliedEffects = append(event.AppliedEffects, chosenEffect.ID())

		rm.logger.Debug("applied replacement effect",
			zap.String("effect_id", chosenEffect.ID()),
			zap.String("event_type", string(event.Type)),
			zap.Bool("completely_replaced", completelyReplaced),
			zap.Int("iteration", iteration))

		// If the effect completely replaced the event, stop processing
		if completelyReplaced {
			rm.logger.Debug("event completely replaced, stopping replacement chain",
				zap.String("event_type", string(event.Type)))
			break
		}

		// If effect has OneUse duration, remove it
		if chosenEffect.Duration() == DurationOneUse {
			rm.logger.Debug("removing one-use replacement effect",
				zap.String("effect_id", chosenEffect.ID()))
			// Note: Can't modify during read lock, would need separate cleanup pass
			// For now, just mark in the event
			if event.Metadata == nil {
				event.Metadata = make(map[string]string)
			}
			event.Metadata["consumed_effect_"+chosenEffect.ID()] = "true"
		}
	}

	if iteration >= maxIterations {
		rm.logger.Error("replacement effect loop exceeded maximum iterations",
			zap.String("event_type", string(event.Type)),
			zap.Int("max_iterations", maxIterations))
	}

	return event
}

// findApplicableEffects returns all effects that could apply to the given event
// and haven't already been applied
func (rm *ReplacementManager) findApplicableEffects(
	event rules.Event,
	gameID string,
	appliedEffects map[string]bool,
) []ReplacementEffect {
	applicable := make([]ReplacementEffect, 0)

	for _, effect := range rm.effects {
		// Skip if already applied to this event
		if appliedEffects[effect.ID()] {
			continue
		}

		// Check if effect checks this event type (fast filter)
		if !effect.ChecksEventType(event.Type) {
			continue
		}

		// Check if effect applies to this specific event (detailed check)
		if !effect.Applies(event, gameID) {
			continue
		}

		// Handle self-scope check (Rule 614.12)
		// Self-scope effects can apply to their own source
		// Non-self-scope effects cannot
		if !effect.HasSelfScope() && event.SourceID == effect.SourceID() {
			continue
		}

		applicable = append(applicable, effect)
	}

	return applicable
}

// GetApplicableEffects returns all replacement effects that could apply to the given event type
// This is useful for checking if any effects might modify an event before it's published
func (rm *ReplacementManager) GetApplicableEffects(eventType rules.EventType, gameID string) []ReplacementEffect {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	applicable := make([]ReplacementEffect, 0)

	for _, effect := range rm.effects {
		if effect.ChecksEventType(eventType) {
			applicable = append(applicable, effect)
		}
	}

	return applicable
}

// HasApplicableEffects checks if there are any replacement effects that could apply to the given event type
func (rm *ReplacementManager) HasApplicableEffects(eventType rules.EventType) bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	for _, effect := range rm.effects {
		if effect.ChecksEventType(eventType) {
			return true
		}
	}

	return false
}

// Stats returns statistics about the replacement manager
func (rm *ReplacementManager) Stats() ReplacementManagerStats {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	stats := ReplacementManagerStats{
		TotalEffects:          len(rm.effects),
		SelfReplacementCount:  0,
		PreventionEffectCount: 0,
	}

	for _, effect := range rm.effects {
		if effect.IsSelfReplacement() {
			stats.SelfReplacementCount++
		}
		if _, ok := effect.(PreventionEffect); ok {
			stats.PreventionEffectCount++
		}
	}

	return stats
}

// ReplacementManagerStats contains statistics about the replacement manager
type ReplacementManagerStats struct {
	TotalEffects          int
	SelfReplacementCount  int
	PreventionEffectCount int
}

// String returns a string representation of the stats
func (s ReplacementManagerStats) String() string {
	return fmt.Sprintf("ReplacementManager[total=%d, self=%d, prevention=%d]",
		s.TotalEffects, s.SelfReplacementCount, s.PreventionEffectCount)
}
