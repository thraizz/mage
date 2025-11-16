package effects

import (
	"testing"
)

// mockDurationEffect is a test helper that implements EffectWithDuration
type mockDurationEffect struct {
	id       string
	layer    Layer
	duration Duration
	sourceID string
}

func (m *mockDurationEffect) ID() string {
	return m.id
}

func (m *mockDurationEffect) Layer() Layer {
	return m.layer
}

func (m *mockDurationEffect) AppliesTo(*Snapshot) bool {
	return true
}

func (m *mockDurationEffect) Apply(*Snapshot) {
	// No-op for test
}

func (m *mockDurationEffect) GetDuration() Duration {
	return m.duration
}

func (m *mockDurationEffect) GetSourceID() string {
	return m.sourceID
}

// TestCleanupEndOfCombatEffects verifies end of combat cleanup
func TestCleanupEndOfCombatEffects(t *testing.T) {
	system := NewLayerSystem()
	
	// Add an end-of-combat effect
	combatEffect := &mockDurationEffect{
		id:       "combat-effect",
		layer:    LayerAbility,
		duration: DurationEndOfCombat,
		sourceID: "source-1",
	}
	system.AddEffect(combatEffect)
	
	// Add an end-of-turn effect (should not be removed)
	turnEffect := &mockDurationEffect{
		id:       "turn-effect",
		layer:    LayerAbility,
		duration: DurationEndOfTurn,
		sourceID: "source-2",
	}
	system.AddEffect(turnEffect)
	
	// Add a permanent effect (should not be removed)
	permEffect := &mockDurationEffect{
		id:       "perm-effect",
		layer:    LayerAbility,
		duration: DurationPermanent,
		sourceID: "source-3",
	}
	system.AddEffect(permEffect)
	
	// Verify all effects were added
	system.mu.RLock()
	initialCount := len(system.effects[LayerAbility])
	system.mu.RUnlock()
	
	if initialCount != 3 {
		t.Fatalf("Expected 3 effects, got %d", initialCount)
	}
	
	// Cleanup end of combat effects
	CleanupEndOfCombatEffects(system)
	
	// Verify only combat effect was removed
	system.mu.RLock()
	finalCount := len(system.effects[LayerAbility])
	_, combatExists := system.effects[LayerAbility]["combat-effect"]
	_, turnExists := system.effects[LayerAbility]["turn-effect"]
	_, permExists := system.effects[LayerAbility]["perm-effect"]
	system.mu.RUnlock()
	
	if finalCount != 2 {
		t.Errorf("Expected 2 effects after cleanup, got %d", finalCount)
	}
	
	if combatExists {
		t.Error("Combat effect should have been removed")
	}
	
	if !turnExists {
		t.Error("Turn effect should not have been removed")
	}
	
	if !permExists {
		t.Error("Permanent effect should not have been removed")
	}
}

// TestCleanupEndOfTurnEffects verifies end of turn cleanup
func TestCleanupEndOfTurnEffects(t *testing.T) {
	system := NewLayerSystem()
	
	// Add an end-of-turn effect
	turnEffect := &mockDurationEffect{
		id:       "turn-effect",
		layer:    LayerAbility,
		duration: DurationEndOfTurn,
		sourceID: "source-1",
	}
	system.AddEffect(turnEffect)
	
	// Add an end-of-combat effect (should not be removed)
	combatEffect := &mockDurationEffect{
		id:       "combat-effect",
		layer:    LayerAbility,
		duration: DurationEndOfCombat,
		sourceID: "source-2",
	}
	system.AddEffect(combatEffect)
	
	// Verify both effects were added
	system.mu.RLock()
	initialCount := len(system.effects[LayerAbility])
	system.mu.RUnlock()
	
	if initialCount != 2 {
		t.Fatalf("Expected 2 effects, got %d", initialCount)
	}
	
	// Cleanup end of turn effects
	CleanupEndOfTurnEffects(system)
	
	// Verify only turn effect was removed
	system.mu.RLock()
	finalCount := len(system.effects[LayerAbility])
	_, turnExists := system.effects[LayerAbility]["turn-effect"]
	_, combatExists := system.effects[LayerAbility]["combat-effect"]
	system.mu.RUnlock()
	
	if finalCount != 1 {
		t.Errorf("Expected 1 effect after cleanup, got %d", finalCount)
	}
	
	if turnExists {
		t.Error("Turn effect should have been removed")
	}
	
	if !combatExists {
		t.Error("Combat effect should not have been removed")
	}
}

// TestCleanupSourceLeftBattlefield verifies source-dependent cleanup
func TestCleanupSourceLeftBattlefield(t *testing.T) {
	system := NewLayerSystem()
	
	sourceID := "source-1"
	
	// Add a WhileOnBattlefield effect
	whileEffect := &mockDurationEffect{
		id:       "while-effect",
		layer:    LayerAbility,
		duration: DurationWhileOnBattlefield,
		sourceID: sourceID,
	}
	system.AddEffect(whileEffect)
	
	// Add a permanent effect from same source (should not be removed)
	permEffect := &mockDurationEffect{
		id:       "perm-effect",
		layer:    LayerAbility,
		duration: DurationPermanent,
		sourceID: sourceID,
	}
	system.AddEffect(permEffect)
	
	// Add an effect from different source (should not be removed)
	otherEffect := &mockDurationEffect{
		id:       "other-effect",
		layer:    LayerAbility,
		duration: DurationWhileOnBattlefield,
		sourceID: "source-2",
	}
	system.AddEffect(otherEffect)
	
	// Verify all effects were added
	system.mu.RLock()
	initialCount := len(system.effects[LayerAbility])
	system.mu.RUnlock()
	
	if initialCount != 3 {
		t.Fatalf("Expected 3 effects, got %d", initialCount)
	}
	
	// Cleanup effects from source-1 leaving battlefield
	CleanupSourceLeftBattlefieldEffects(system, sourceID)
	
	// Verify only while-effect was removed
	system.mu.RLock()
	finalCount := len(system.effects[LayerAbility])
	_, whileExists := system.effects[LayerAbility]["while-effect"]
	_, permExists := system.effects[LayerAbility]["perm-effect"]
	_, otherExists := system.effects[LayerAbility]["other-effect"]
	system.mu.RUnlock()
	
	if finalCount != 2 {
		t.Errorf("Expected 2 effects after cleanup, got %d", finalCount)
	}
	
	if whileExists {
		t.Error("While effect should have been removed")
	}
	
	if !permExists {
		t.Error("Permanent effect should not have been removed")
	}
	
	if !otherExists {
		t.Error("Other source effect should not have been removed")
	}
}
