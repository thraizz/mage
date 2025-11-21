package effects

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestReplacementManager_AddRemoveEffect(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	effect := NewDamagePreventionEffect("source1", "target1", "", 5, DurationUntilEndOfTurn)

	// Add effect
	rm.AddEffect(effect)

	// Verify it was added
	retrieved, ok := rm.GetEffect(effect.ID())
	require.True(t, ok)
	assert.Equal(t, effect.ID(), retrieved.ID())

	// Remove effect
	rm.RemoveEffect(effect.ID())

	// Verify it was removed
	_, ok = rm.GetEffect(effect.ID())
	assert.False(t, ok)
}

func TestReplacementManager_GetEffects(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	effect1 := NewDamagePreventionEffect("source1", "target1", "", 5, DurationUntilEndOfTurn)
	effect2 := NewDoubleAmountReplacementEffect("source2", []rules.EventType{rules.EventGainLife}, "player1", "", DurationPermanent)

	rm.AddEffect(effect1)
	rm.AddEffect(effect2)

	effects := rm.GetEffects()
	assert.Len(t, effects, 2)
}

func TestReplacementManager_ClearEffects(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	effect1 := NewDamagePreventionEffect("source1", "target1", "", 5, DurationUntilEndOfTurn)
	effect2 := NewDoubleAmountReplacementEffect("source2", []rules.EventType{rules.EventGainLife}, "player1", "", DurationPermanent)

	rm.AddEffect(effect1)
	rm.AddEffect(effect2)

	assert.Len(t, rm.GetEffects(), 2)

	rm.ClearEffects()

	assert.Len(t, rm.GetEffects(), 0)
}

func TestReplacementManager_ReplaceEvent_SingleEffect(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	// Add a damage prevention effect
	effect := NewDamagePreventionEffect("shield", "target1", "", 0, DurationPermanent)
	rm.AddEffect(effect)

	// Create a damage event
	event := rules.NewEventWithAmount(rules.EventDamagePermanent, "target1", "attacker1", "player1", 5)

	// Replace the event
	modifiedEvent := rm.ReplaceEvent(event, "game1", "player1")

	// Damage should be prevented
	assert.Equal(t, 0, modifiedEvent.Amount)

	// Effect should be marked as applied
	require.Len(t, modifiedEvent.AppliedEffects, 1)
	assert.Equal(t, effect.ID(), modifiedEvent.AppliedEffects[0])
}

func TestReplacementManager_ReplaceEvent_MultipleEffects(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	// Add two doubling effects (simulating two "double life gain" effects)
	// Use different source IDs so self-scope check doesn't interfere
	effect1 := NewDoubleAmountReplacementEffect("doubler1", []rules.EventType{rules.EventGainLife}, "player1", "", DurationPermanent)
	effect2 := NewDoubleAmountReplacementEffect("doubler2", []rules.EventType{rules.EventGainLife}, "player1", "", DurationPermanent)

	rm.AddEffect(effect1)
	rm.AddEffect(effect2)

	// Create a life gain event from a different source
	event := rules.NewEventWithAmount(rules.EventGainLife, "player1", "spell1", "player1", 2)

	// Replace the event
	modifiedEvent := rm.ReplaceEvent(event, "game1", "player1")

	// Both effects should apply: 2 * 2 * 2 = 8
	assert.Equal(t, 8, modifiedEvent.Amount)

	// Both effects should be marked as applied
	assert.Len(t, modifiedEvent.AppliedEffects, 2)
}

func TestReplacementManager_ReplaceEvent_PreventDoubleApplication(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	// Add a doubling effect
	effect := NewDoubleAmountReplacementEffect("source1", []rules.EventType{rules.EventGainLife}, "player1", "", DurationPermanent)
	rm.AddEffect(effect)

	// Create an event that already has this effect applied
	event := rules.NewEventWithAmount(rules.EventGainLife, "player1", "source1", "player1", 4)
	event.AppliedEffects = []string{effect.ID()}

	// Replace the event
	modifiedEvent := rm.ReplaceEvent(event, "game1", "player1")

	// Amount should not change (effect already applied)
	assert.Equal(t, 4, modifiedEvent.Amount)

	// Should still have only one applied effect
	assert.Len(t, modifiedEvent.AppliedEffects, 1)
}

func TestReplacementManager_ReplaceEvent_CompleteReplacement(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	// Add two effects: damage prevention and doubling
	prevention := NewDamagePreventionEffect("shield", "target1", "", 0, DurationPermanent)
	doubling := NewDoubleAmountReplacementEffect("doubler", []rules.EventType{rules.EventDamagePermanent}, "target1", "", DurationPermanent)

	rm.AddEffect(prevention)
	rm.AddEffect(doubling)

	// Create a damage event
	event := rules.NewEventWithAmount(rules.EventDamagePermanent, "target1", "attacker1", "player1", 5)

	// Replace the event
	modifiedEvent := rm.ReplaceEvent(event, "game1", "player1")

	// Damage should be prevented (completely replaced)
	// If prevention is applied first, damage becomes 0 and doubling won't apply
	// OR if doubling is applied first, damage becomes 10, then prevention makes it 0
	assert.Equal(t, 0, modifiedEvent.Amount)

	// At least prevention should be applied
	require.NotEmpty(t, modifiedEvent.AppliedEffects)
}

func TestReplacementManager_ReplaceEvent_SelfReplacementFirst(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	// Add a regular effect and a self-replacement effect
	regularEffect := NewDoubleAmountReplacementEffect("regular", []rules.EventType{rules.EventGainLife}, "player1", "", DurationPermanent)

	// Create a self-replacement effect (we need to use the base effect directly for testing)
	selfEffect := &testSelfReplacementEffect{
		BaseReplacementEffect: NewBaseReplacementEffect("self", DurationOneUse, true, false),
		multiplier:            3,
	}

	rm.AddEffect(regularEffect)
	rm.AddEffect(selfEffect)

	// Create a life gain event
	event := rules.NewEventWithAmount(rules.EventGainLife, "player1", "source1", "player1", 2)

	// Replace the event
	modifiedEvent := rm.ReplaceEvent(event, "game1", "player1")

	// Self-replacement should apply first (x3), then regular (x2)
	// 2 * 3 * 2 = 12
	assert.Equal(t, 12, modifiedEvent.Amount)

	// Both effects should be marked as applied
	assert.Len(t, modifiedEvent.AppliedEffects, 2)

	// First applied should be the self-replacement effect
	assert.Equal(t, selfEffect.ID(), modifiedEvent.AppliedEffects[0])
}

func TestReplacementManager_ReplaceEvent_SelfScopeCheck(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	// Add an effect WITHOUT self-scope
	effect := NewZoneChangeReplacementEffect("source1", -1, -1, 5, "", "", "", DurationPermanent, false)
	rm.AddEffect(effect)

	// Create an event from the same source
	event := rules.NewEvent(rules.EventZoneChange, "target1", "source1", "player1")

	// Replace the event
	modifiedEvent := rm.ReplaceEvent(event, "game1", "player1")

	// Effect should NOT apply (no self-scope, event is from the same source)
	assert.Empty(t, modifiedEvent.AppliedEffects)
}

func TestReplacementManager_ReplaceEvent_WithSelfScope(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	// Add an effect WITH self-scope
	effect := NewZoneChangeReplacementEffect("source1", -1, 2, 5, "source1", "", "", DurationPermanent, true)
	rm.AddEffect(effect)

	// Create an event from the same source
	event := rules.NewEvent(rules.EventZoneChange, "source1", "source1", "player1")
	event.Zone = 2

	// Replace the event
	modifiedEvent := rm.ReplaceEvent(event, "game1", "player1")

	// Effect SHOULD apply (has self-scope)
	assert.NotEmpty(t, modifiedEvent.AppliedEffects)
	assert.Equal(t, 5, modifiedEvent.Zone) // Changed to exile
}

func TestReplacementManager_HasApplicableEffects(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	// Initially no effects
	assert.False(t, rm.HasApplicableEffects(rules.EventDamagePermanent))

	// Add a damage prevention effect
	effect := NewDamagePreventionEffect("shield", "target1", "", 5, DurationUntilEndOfTurn)
	rm.AddEffect(effect)

	// Now should have applicable effects for damage
	assert.True(t, rm.HasApplicableEffects(rules.EventDamagePermanent))
	assert.True(t, rm.HasApplicableEffects(rules.EventDamagePlayer))

	// Should not have effects for other event types
	assert.False(t, rm.HasApplicableEffects(rules.EventDrawCard))
}

func TestReplacementManager_GetApplicableEffects(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	damageEffect := NewDamagePreventionEffect("shield", "target1", "", 5, DurationUntilEndOfTurn)
	lifeEffect := NewDoubleAmountReplacementEffect("doubler", []rules.EventType{rules.EventGainLife}, "player1", "", DurationPermanent)

	rm.AddEffect(damageEffect)
	rm.AddEffect(lifeEffect)

	// Get effects for damage
	damageEffects := rm.GetApplicableEffects(rules.EventDamagePermanent, "game1")
	assert.Len(t, damageEffects, 1)

	// Get effects for life gain
	lifeEffects := rm.GetApplicableEffects(rules.EventGainLife, "game1")
	assert.Len(t, lifeEffects, 1)

	// Get effects for draw (none)
	drawEffects := rm.GetApplicableEffects(rules.EventDrawCard, "game1")
	assert.Len(t, drawEffects, 0)
}

func TestReplacementManager_Stats(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	prevention1 := NewDamagePreventionEffect("shield1", "target1", "", 5, DurationUntilEndOfTurn)
	prevention2 := NewDamagePreventionEffect("shield2", "target2", "", 0, DurationPermanent)
	doubling := NewDoubleAmountReplacementEffect("doubler", []rules.EventType{rules.EventGainLife}, "player1", "", DurationPermanent)

	rm.AddEffect(prevention1)
	rm.AddEffect(prevention2)
	rm.AddEffect(doubling)

	stats := rm.Stats()

	assert.Equal(t, 3, stats.TotalEffects)
	assert.Equal(t, 2, stats.PreventionEffectCount)
	assert.Equal(t, 0, stats.SelfReplacementCount) // None of these are self-replacement

	// Test string representation
	str := stats.String()
	assert.Contains(t, str, "total=3")
	assert.Contains(t, str, "prevention=2")
}

func TestReplacementManager_NilEffect(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	// Should not panic when adding nil effect
	rm.AddEffect(nil)

	// Should still be empty
	assert.Len(t, rm.GetEffects(), 0)
}

func TestReplacementManager_RemoveNonExistentEffect(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	// Should not panic when removing non-existent effect
	rm.RemoveEffect("nonexistent")

	// Should still be empty
	assert.Len(t, rm.GetEffects(), 0)
}

func TestReplacementManager_MaxIterations(t *testing.T) {
	logger := zap.NewNop()
	rm := NewReplacementManager(logger)

	// Create an event
	event := rules.NewEventWithAmount(rules.EventGainLife, "player1", "source1", "player1", 2)

	// Add many doubling effects (but they won't infinite loop due to applied tracking)
	for i := 0; i < 10; i++ {
		effect := NewDoubleAmountReplacementEffect(
			"source"+string(rune(i)),
			[]rules.EventType{rules.EventGainLife},
			"player1",
			"",
			DurationPermanent,
		)
		rm.AddEffect(effect)
	}

	// Replace the event
	modifiedEvent := rm.ReplaceEvent(event, "game1", "player1")

	// All effects should apply: 2 * 2^10 = 2048
	assert.Equal(t, 2048, modifiedEvent.Amount)

	// All 10 effects should be marked as applied
	assert.Len(t, modifiedEvent.AppliedEffects, 10)
}

// testSelfReplacementEffect is a test helper that implements a self-replacement effect
type testSelfReplacementEffect struct {
	*BaseReplacementEffect
	multiplier int
}

func (e *testSelfReplacementEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventGainLife
}

func (e *testSelfReplacementEffect) Applies(event rules.Event, gameID string) bool {
	return event.Type == rules.EventGainLife && event.TargetID == "player1"
}

func (e *testSelfReplacementEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	event.Amount = event.Amount * e.multiplier
	return event, false
}
