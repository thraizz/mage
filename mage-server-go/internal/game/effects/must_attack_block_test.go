package effects

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMustAttackEffect tests the MustAttackEffect creation and behavior
func TestMustAttackEffect(t *testing.T) {
	sourceID := "creature-1"
	targetIDs := []string{"creature-1"}
	duration := DurationEndOfTurn

	effect := NewMustAttackEffect(sourceID, targetIDs, duration)

	require.NotNil(t, effect)
	assert.Equal(t, sourceID, effect.GetSourceID())
	assert.Equal(t, duration, effect.GetDuration())
	assert.Equal(t, targetIDs, effect.GetTargetIDs())
	assert.NotEmpty(t, effect.ID())
	assert.Equal(t, LayerPowerToughness, effect.Layer())
}

// TestMustAttackEffectAppliesTo tests the AppliesTo logic
func TestMustAttackEffectAppliesTo(t *testing.T) {
	sourceID := "source-1"
	targetID := "creature-1"
	effect := NewMustAttackEffect(sourceID, []string{targetID}, DurationEndOfTurn)

	// Test that it applies to the target creature
	snapshot := &Snapshot{CardID: targetID}
	assert.True(t, effect.AppliesTo(snapshot), "should apply to target creature")

	// Test that it doesn't apply to a different creature
	differentSnapshot := &Snapshot{CardID: "creature-2"}
	assert.False(t, effect.AppliesTo(differentSnapshot), "should not apply to different creature")

	// Test nil snapshot
	assert.False(t, effect.AppliesTo(nil), "should not apply to nil snapshot")
}

// TestMustBeBlockedEffect tests the MustBeBlockedEffect creation and behavior
func TestMustBeBlockedEffect(t *testing.T) {
	sourceID := "attacker-1"
	attackerID := "" // Empty means source itself
	targetIDs := []string{} // Empty means all able blockers
	duration := DurationWhileOnBattlefield

	effect := NewMustBeBlockedEffect(sourceID, attackerID, targetIDs, duration)

	require.NotNil(t, effect)
	assert.Equal(t, sourceID, effect.GetSourceID())
	assert.Equal(t, duration, effect.GetDuration())
	assert.Equal(t, targetIDs, effect.GetTargetIDs())
	assert.Equal(t, sourceID, effect.GetAttackerID(), "should return source when attackerID is empty")
	assert.NotEmpty(t, effect.ID())
	assert.Equal(t, LayerPowerToughness, effect.Layer())
}

// TestMustBeBlockedEffectWithSpecificAttacker tests effect with specific attacker
func TestMustBeBlockedEffectWithSpecificAttacker(t *testing.T) {
	sourceID := "enchantment-1"
	attackerID := "attacker-1"
	targetIDs := []string{"blocker-1", "blocker-2"}
	duration := DurationEndOfTurn

	effect := NewMustBeBlockedEffect(sourceID, attackerID, targetIDs, duration)

	require.NotNil(t, effect)
	assert.Equal(t, attackerID, effect.GetAttackerID(), "should return specified attacker ID")
	assert.Equal(t, targetIDs, effect.GetTargetIDs())
}

// TestMustBeBlockedEffectAppliesTo tests the AppliesTo logic
func TestMustBeBlockedEffectAppliesTo(t *testing.T) {
	sourceID := "attacker-1"
	targetIDs := []string{"blocker-1", "blocker-2"}
	effect := NewMustBeBlockedEffect(sourceID, "", targetIDs, DurationWhileOnBattlefield)

	// Test that it applies to target blockers
	snapshot1 := &Snapshot{CardID: "blocker-1"}
	assert.True(t, effect.AppliesTo(snapshot1), "should apply to blocker-1")

	snapshot2 := &Snapshot{CardID: "blocker-2"}
	assert.True(t, effect.AppliesTo(snapshot2), "should apply to blocker-2")

	// Test that it doesn't apply to a non-target creature
	differentSnapshot := &Snapshot{CardID: "blocker-3"}
	assert.False(t, effect.AppliesTo(differentSnapshot), "should not apply to non-target blocker")

	// Test nil snapshot
	assert.False(t, effect.AppliesTo(nil), "should not apply to nil snapshot")
}

// TestMustBeBlockedEffectAppliesToAll tests that empty targetIDs applies to all
func TestMustBeBlockedEffectAppliesToAll(t *testing.T) {
	sourceID := "attacker-1"
	targetIDs := []string{} // Empty = all able blockers
	effect := NewMustBeBlockedEffect(sourceID, "", targetIDs, DurationWhileOnBattlefield)

	// Test that it applies to any creature when targetIDs is empty
	snapshot1 := &Snapshot{CardID: "blocker-1"}
	assert.True(t, effect.AppliesTo(snapshot1), "should apply to any blocker when targetIDs empty")

	snapshot2 := &Snapshot{CardID: "blocker-2"}
	assert.True(t, effect.AppliesTo(snapshot2), "should apply to any blocker when targetIDs empty")

	// Test nil snapshot still returns false
	assert.False(t, effect.AppliesTo(nil), "should not apply to nil snapshot")
}

// TestMustAttackEffectDifferentTargets tests effect with multiple targets
func TestMustAttackEffectDifferentTargets(t *testing.T) {
	sourceID := "enchantment-1"
	targetIDs := []string{"creature-1", "creature-2", "creature-3"}
	effect := NewMustAttackEffect(sourceID, targetIDs, DurationEndOfTurn)

	// All targets should be affected
	for _, targetID := range targetIDs {
		snapshot := &Snapshot{CardID: targetID}
		assert.True(t, effect.AppliesTo(snapshot), "should apply to target %s", targetID)
	}

	// Non-target should not be affected
	nonTarget := &Snapshot{CardID: "creature-4"}
	assert.False(t, effect.AppliesTo(nonTarget), "should not apply to non-target")
}

// TestMustBeBlockedEffectDuration tests different durations
func TestMustBeBlockedEffectDuration(t *testing.T) {
	sourceID := "attacker-1"

	tests := []struct {
		name     string
		duration Duration
	}{
		{"EndOfTurn", DurationEndOfTurn},
		{"WhileOnBattlefield", DurationWhileOnBattlefield},
		{"EndOfCombat", DurationEndOfCombat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			effect := NewMustBeBlockedEffect(sourceID, "", []string{}, tt.duration)
			assert.Equal(t, tt.duration, effect.GetDuration())
		})
	}
}

// TestMustAttackEffectApplyNoOp tests that Apply is a no-op
func TestMustAttackEffectApplyNoOp(t *testing.T) {
	effect := NewMustAttackEffect("source-1", []string{"creature-1"}, DurationEndOfTurn)
	snapshot := &Snapshot{CardID: "creature-1", Power: 3, Toughness: 3}

	originalPower := snapshot.Power
	originalToughness := snapshot.Toughness

	// Apply should not modify the snapshot
	effect.Apply(snapshot)

	assert.Equal(t, originalPower, snapshot.Power, "power should not change")
	assert.Equal(t, originalToughness, snapshot.Toughness, "toughness should not change")
}

// TestMustBeBlockedEffectApplyNoOp tests that Apply is a no-op
func TestMustBeBlockedEffectApplyNoOp(t *testing.T) {
	effect := NewMustBeBlockedEffect("attacker-1", "", []string{"blocker-1"}, DurationWhileOnBattlefield)
	snapshot := &Snapshot{CardID: "blocker-1", Power: 2, Toughness: 2}

	originalPower := snapshot.Power
	originalToughness := snapshot.Toughness

	// Apply should not modify the snapshot
	effect.Apply(snapshot)

	assert.Equal(t, originalPower, snapshot.Power, "power should not change")
	assert.Equal(t, originalToughness, snapshot.Toughness, "toughness should not change")
}
