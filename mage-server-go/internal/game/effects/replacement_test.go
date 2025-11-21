package effects

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/rules"
	"github.com/stretchr/testify/assert"
)

func TestDamagePreventionEffect_PreventAll(t *testing.T) {
	effect := NewDamagePreventionEffect("source1", "target1", "", 0, DurationUntilEndOfTurn)

	// Create a damage event
	event := rules.NewEventWithAmount(rules.EventDamagePermanent, "target1", "attacker1", "player1", 5)

	// Check if effect applies
	assert.True(t, effect.ChecksEventType(event.Type))
	assert.True(t, effect.Applies(event, "game1"))

	// Apply the effect
	modifiedEvent, completelyReplaced := effect.ReplaceEvent(event, "game1")

	// Should prevent all damage
	assert.True(t, completelyReplaced)
	assert.Equal(t, 0, modifiedEvent.Amount)
}

func TestDamagePreventionEffect_PreventPartial(t *testing.T) {
	effect := NewDamagePreventionEffect("source1", "target1", "", 3, DurationUntilEndOfTurn)

	// Create a damage event with 5 damage
	event := rules.NewEventWithAmount(rules.EventDamagePermanent, "target1", "attacker1", "player1", 5)

	// Apply the effect
	modifiedEvent, completelyReplaced := effect.ReplaceEvent(event, "game1")

	// Should prevent 3 damage, leaving 2
	assert.False(t, completelyReplaced)
	assert.Equal(t, 2, modifiedEvent.Amount)

	// Shield should be consumed
	assert.Equal(t, 0, effect.GetShield())
}

func TestDamagePreventionEffect_Shield(t *testing.T) {
	effect := NewDamagePreventionEffect("source1", "target1", "", 5, DurationUntilEndOfTurn)

	// First damage event - 3 damage
	event1 := rules.NewEventWithAmount(rules.EventDamagePermanent, "target1", "attacker1", "player1", 3)
	modifiedEvent1, _ := effect.ReplaceEvent(event1, "game1")

	// Should prevent all 3 damage
	assert.Equal(t, 0, modifiedEvent1.Amount)
	assert.Equal(t, 2, effect.GetShield()) // 2 shield remaining

	// Second damage event - 4 damage
	event2 := rules.NewEventWithAmount(rules.EventDamagePermanent, "target1", "attacker1", "player1", 4)
	modifiedEvent2, completelyReplaced := effect.ReplaceEvent(event2, "game1")

	// Should prevent 2 damage, leaving 2
	assert.False(t, completelyReplaced)
	assert.Equal(t, 2, modifiedEvent2.Amount)
	assert.Equal(t, 0, effect.GetShield()) // Shield exhausted

	// Third damage event - should not apply since shield is exhausted
	event3 := rules.NewEventWithAmount(rules.EventDamagePermanent, "target1", "attacker1", "player1", 2)
	assert.False(t, effect.Applies(event3, "game1")) // Should not apply anymore
}

func TestDamagePreventionEffect_TargetFilter(t *testing.T) {
	effect := NewDamagePreventionEffect("source1", "target1", "", 0, DurationUntilEndOfTurn)

	// Event targeting target1 - should apply
	event1 := rules.NewEventWithAmount(rules.EventDamagePermanent, "target1", "attacker1", "player1", 5)
	assert.True(t, effect.Applies(event1, "game1"))

	// Event targeting target2 - should not apply
	event2 := rules.NewEventWithAmount(rules.EventDamagePermanent, "target2", "attacker1", "player1", 5)
	assert.False(t, effect.Applies(event2, "game1"))
}

func TestDamagePreventionEffect_SourceFilter(t *testing.T) {
	effect := NewDamagePreventionEffect("source1", "", "specific_source", 0, DurationUntilEndOfTurn)

	// Event from specific_source - should apply
	event1 := rules.NewEventWithAmount(rules.EventDamagePermanent, "target1", "specific_source", "player1", 5)
	assert.True(t, effect.Applies(event1, "game1"))

	// Event from other_source - should not apply
	event2 := rules.NewEventWithAmount(rules.EventDamagePermanent, "target1", "other_source", "player1", 5)
	assert.False(t, effect.Applies(event2, "game1"))
}

func TestZoneChangeReplacementEffect_ExileInsteadOfGraveyard(t *testing.T) {
	// "If a creature would die, exile it instead"
	effect := NewZoneChangeReplacementEffect(
		"source1",
		2,  // from battlefield
		3,  // to graveyard
		5,  // exile instead
		"", // any card
		"", // any controller
		"creature",
		DurationPermanent,
		false, // not self-scope
	)

	event := rules.NewEvent(rules.EventZoneChange, "creature1", "source1", "player1")
	event.Zone = 3 // going to graveyard

	assert.True(t, effect.ChecksEventType(event.Type))
	assert.True(t, effect.Applies(event, "game1"))

	modifiedEvent, completelyReplaced := effect.ReplaceEvent(event, "game1")

	// Should change destination to exile
	assert.False(t, completelyReplaced) // Zone change still happens
	assert.Equal(t, 5, modifiedEvent.Zone)
	assert.Equal(t, "5", modifiedEvent.Metadata["new_zone"])
}

func TestZoneChangeReplacementEffect_SpecificCard(t *testing.T) {
	effect := NewZoneChangeReplacementEffect(
		"source1",
		-1, // any from zone
		-1, // any to zone
		5,  // exile
		"specific_card", // only this card
		"",
		"",
		DurationPermanent,
		false,
	)

	// Event for specific_card - should apply
	event1 := rules.NewEvent(rules.EventZoneChange, "specific_card", "source1", "player1")
	assert.True(t, effect.Applies(event1, "game1"))

	// Event for other_card - should not apply
	event2 := rules.NewEvent(rules.EventZoneChange, "other_card", "source1", "player1")
	assert.False(t, effect.Applies(event2, "game1"))
}

func TestDoubleAmountReplacementEffect_DoubleLife(t *testing.T) {
	effect := NewDoubleAmountReplacementEffect(
		"source1",
		[]rules.EventType{rules.EventGainLife, rules.EventGainedLife},
		"player1",
		"",
		DurationPermanent,
	)

	event := rules.NewEventWithAmount(rules.EventGainLife, "player1", "source1", "player1", 5)

	assert.True(t, effect.ChecksEventType(event.Type))
	assert.True(t, effect.Applies(event, "game1"))

	modifiedEvent, completelyReplaced := effect.ReplaceEvent(event, "game1")

	// Should double the amount
	assert.False(t, completelyReplaced)
	assert.Equal(t, 10, modifiedEvent.Amount)
	assert.Equal(t, effect.ID(), modifiedEvent.Metadata["doubled_by"])
}

func TestDoubleAmountReplacementEffect_TargetFilter(t *testing.T) {
	effect := NewDoubleAmountReplacementEffect(
		"source1",
		[]rules.EventType{rules.EventGainLife},
		"player1", // only player1
		"",
		DurationPermanent,
	)

	// Event for player1 - should apply
	event1 := rules.NewEventWithAmount(rules.EventGainLife, "player1", "source1", "player1", 5)
	assert.True(t, effect.Applies(event1, "game1"))

	// Event for player2 - should not apply
	event2 := rules.NewEventWithAmount(rules.EventGainLife, "player2", "source1", "player2", 5)
	assert.False(t, effect.Applies(event2, "game1"))
}

func TestDoubleAmountReplacementEffect_MultipleDoubling(t *testing.T) {
	// Test Rule 614.5 - effect doesn't invoke itself repeatedly
	// Two separate doubling effects should stack (2 * 2 = 4), not infinite

	effect1 := NewDoubleAmountReplacementEffect(
		"source1",
		[]rules.EventType{rules.EventGainLife},
		"player1",
		"",
		DurationPermanent,
	)

	effect2 := NewDoubleAmountReplacementEffect(
		"source2",
		[]rules.EventType{rules.EventGainLife},
		"player1",
		"",
		DurationPermanent,
	)

	event := rules.NewEventWithAmount(rules.EventGainLife, "player1", "source1", "player1", 2)

	// Apply first effect
	event1, _ := effect1.ReplaceEvent(event, "game1")
	assert.Equal(t, 4, event1.Amount) // 2 * 2 = 4

	// Mark first effect as applied
	event1.AppliedEffects = []string{effect1.ID()}

	// Apply second effect
	event2, _ := effect2.ReplaceEvent(event1, "game1")
	assert.Equal(t, 8, event2.Amount) // 4 * 2 = 8

	// The doubling effects should stack, giving 8 (not infinite)
	// This matches the Comprehensive Rules example:
	// "A creature that normally deals 2 damage will deal 8 damage—not just 4, and not an infinite amount"
}

func TestBaseReplacementEffect_SelfReplacement(t *testing.T) {
	effect := NewBaseReplacementEffect("source1", DurationOneUse, true, false)

	assert.True(t, effect.IsSelfReplacement())
	assert.False(t, effect.HasSelfScope())
	assert.Equal(t, DurationOneUse, effect.Duration())
}

func TestBaseReplacementEffect_SelfScope(t *testing.T) {
	effect := NewBaseReplacementEffect("source1", DurationPermanent, false, true)

	assert.False(t, effect.IsSelfReplacement())
	assert.True(t, effect.HasSelfScope())
}

func TestBasePreventionEffect_UnlimitedShield(t *testing.T) {
	effect := NewBasePreventionEffect("source1", DurationPermanent, 0)

	// Shield of 0 means unlimited (no shield system)
	assert.Equal(t, 0, effect.GetShield())

	// Reducing shield returns 0 since there's no shield system
	reduced := effect.ReduceShield(10)
	assert.Equal(t, 0, reduced)
	assert.Equal(t, 0, effect.GetShield()) // Still 0 (no shield system)
}

func TestBasePreventionEffect_LimitedShield(t *testing.T) {
	effect := NewBasePreventionEffect("source1", DurationOneUse, 5)

	assert.Equal(t, 5, effect.GetShield())

	// Reduce by 3
	reduced := effect.ReduceShield(3)
	assert.Equal(t, 3, reduced)
	assert.Equal(t, 2, effect.GetShield())

	// Try to reduce by 5, but only 2 remaining
	reduced = effect.ReduceShield(5)
	assert.Equal(t, 2, reduced)
	assert.Equal(t, 0, effect.GetShield())

	// Shield exhausted - can't reduce anymore
	reduced = effect.ReduceShield(10)
	assert.Equal(t, 0, reduced) // Can't reduce when shield is at 0
	assert.Equal(t, 0, effect.GetShield())
}

func TestReplacementEffect_IDs(t *testing.T) {
	effect1 := NewBaseReplacementEffect("source1", DurationPermanent, false, false)
	effect2 := NewBaseReplacementEffect("source1", DurationPermanent, false, false)

	// IDs should be unique even with same parameters (due to UUID seed)
	assert.NotEqual(t, effect1.ID(), effect2.ID())

	// Source ID should match
	assert.Equal(t, "source1", effect1.SourceID())
	assert.Equal(t, "source1", effect2.SourceID())
}

func TestDamagePreventionEffect_PlayerDamage(t *testing.T) {
	effect := NewDamagePreventionEffect("shield", "player1", "", 0, DurationPermanent)

	// Test with EventDamagePlayer
	event := rules.NewEventWithAmount(rules.EventDamagePlayer, "player1", "creature1", "player1", 10)

	assert.True(t, effect.ChecksEventType(event.Type))
	assert.True(t, effect.Applies(event, "game1"))

	modifiedEvent, completelyReplaced := effect.ReplaceEvent(event, "game1")

	assert.True(t, completelyReplaced)
	assert.Equal(t, 0, modifiedEvent.Amount)
}

func TestZoneChangeReplacementEffect_SelfScope(t *testing.T) {
	// Test self-scope for "enters the battlefield" effects
	effect := NewZoneChangeReplacementEffect(
		"card1",    // source
		-1,         // any from zone
		2,          // to battlefield
		2,          // stay on battlefield (modify properties)
		"card1",    // the card itself
		"",
		"",
		DurationPermanent,
		true, // has self-scope
	)

	assert.True(t, effect.HasSelfScope())

	// Event for the source card itself - should apply due to self-scope
	event := rules.NewEvent(rules.EventZoneChange, "card1", "card1", "player1")
	event.Zone = 2

	assert.True(t, effect.Applies(event, "game1"))
}
