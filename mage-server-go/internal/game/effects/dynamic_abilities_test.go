package effects

import (
	"testing"
)

// TestDynamicAbilityGranting_Flying verifies flying can be granted dynamically
func TestDynamicAbilityGranting_Flying(t *testing.T) {
	system := NewLayerSystem()
	manager := NewEffectManager(system)
	
	creatureID := "grizzly-bears"
	spellID := "giant-growth"
	
	// Initially creature doesn't have flying
	if manager.HasGrantedAbility(creatureID, "FlyingAbility") {
		t.Error("Creature should not have flying initially")
	}
	
	// Cast spell that grants flying until end of turn
	effect := NewEffectBuilder(spellID).
		Targeting(creatureID).
		UntilEndOfTurn().
		GrantAbility("FlyingAbility")
	
	effectID := manager.AddEffect(effect)
	
	// Now creature should have flying
	if !manager.HasGrantedAbility(creatureID, "FlyingAbility") {
		t.Error("Creature should have flying after spell")
	}
	
	// Simulate end of turn cleanup
	CleanupEndOfTurnEffects(system)
	
	// Flying should be gone
	if manager.HasGrantedAbility(creatureID, "FlyingAbility") {
		t.Error("Creature should not have flying after end of turn")
	}
	
	_ = effectID // Suppress unused warning
}

// TestDynamicAbilityGranting_Vigilance verifies vigilance can be granted
func TestDynamicAbilityGranting_Vigilance(t *testing.T) {
	system := NewLayerSystem()
	manager := NewEffectManager(system)
	
	creatureID := "grizzly-bears"
	spellID := "inspired-charge"
	
	// Grant vigilance until end of combat
	effect := NewEffectBuilder(spellID).
		Targeting(creatureID).
		UntilEndOfCombat().
		GrantAbility("VigilanceAbility")
	
	manager.AddEffect(effect)
	
	// Creature should have vigilance
	if !manager.HasGrantedAbility(creatureID, "VigilanceAbility") {
		t.Error("Creature should have vigilance")
	}
	
	// Simulate end of combat cleanup
	CleanupEndOfCombatEffects(system)
	
	// Vigilance should be gone
	if manager.HasGrantedAbility(creatureID, "VigilanceAbility") {
		t.Error("Creature should not have vigilance after combat")
	}
}

// TestDynamicAbilityGranting_MultipleAbilities verifies multiple abilities
func TestDynamicAbilityGranting_MultipleAbilities(t *testing.T) {
	system := NewLayerSystem()
	manager := NewEffectManager(system)
	
	creatureID := "grizzly-bears"
	
	// Grant flying
	effect1 := NewEffectBuilder("spell-1").
		Targeting(creatureID).
		UntilEndOfTurn().
		GrantAbility("FlyingAbility")
	manager.AddEffect(effect1)
	
	// Grant vigilance
	effect2 := NewEffectBuilder("spell-2").
		Targeting(creatureID).
		UntilEndOfTurn().
		GrantAbility("VigilanceAbility")
	manager.AddEffect(effect2)
	
	// Grant first strike
	effect3 := NewEffectBuilder("spell-3").
		Targeting(creatureID).
		UntilEndOfTurn().
		GrantAbility("FirstStrikeAbility")
	manager.AddEffect(effect3)
	
	// Creature should have all three abilities
	if !manager.HasGrantedAbility(creatureID, "FlyingAbility") {
		t.Error("Should have flying")
	}
	if !manager.HasGrantedAbility(creatureID, "VigilanceAbility") {
		t.Error("Should have vigilance")
	}
	if !manager.HasGrantedAbility(creatureID, "FirstStrikeAbility") {
		t.Error("Should have first strike")
	}
	
	// Get all effects
	effects := manager.GetEffectsForCard(creatureID)
	if len(effects) != 3 {
		t.Errorf("Expected 3 effects, got %d", len(effects))
	}
}

// TestDynamicAbilityGranting_MultipleTargets verifies one spell affecting multiple creatures
func TestDynamicAbilityGranting_MultipleTargets(t *testing.T) {
	system := NewLayerSystem()
	manager := NewEffectManager(system)
	
	creature1 := "bear-1"
	creature2 := "bear-2"
	creature3 := "bear-3"
	
	// One spell grants flying to all three creatures
	effect := NewEffectBuilder("mass-flight").
		Targeting(creature1, creature2, creature3).
		UntilEndOfTurn().
		GrantAbility("FlyingAbility")
	
	manager.AddEffect(effect)
	
	// All three should have flying
	if !manager.HasGrantedAbility(creature1, "FlyingAbility") {
		t.Error("Creature 1 should have flying")
	}
	if !manager.HasGrantedAbility(creature2, "FlyingAbility") {
		t.Error("Creature 2 should have flying")
	}
	if !manager.HasGrantedAbility(creature3, "FlyingAbility") {
		t.Error("Creature 3 should have flying")
	}
}

// TestDynamicAbilityGranting_Permanent verifies permanent ability grants
func TestDynamicAbilityGranting_Permanent(t *testing.T) {
	system := NewLayerSystem()
	manager := NewEffectManager(system)
	
	creatureID := "grizzly-bears"
	enchantmentID := "enchantment-1"
	
	// Enchantment grants flying permanently
	effect := NewEffectBuilder(enchantmentID).
		Targeting(creatureID).
		Permanent().
		GrantAbility("FlyingAbility")
	
	manager.AddEffect(effect)
	
	// Creature should have flying
	if !manager.HasGrantedAbility(creatureID, "FlyingAbility") {
		t.Error("Creature should have flying")
	}
	
	// End of turn cleanup shouldn't remove it
	CleanupEndOfTurnEffects(system)
	if !manager.HasGrantedAbility(creatureID, "FlyingAbility") {
		t.Error("Permanent effect should survive end of turn")
	}
	
	// End of combat cleanup shouldn't remove it
	CleanupEndOfCombatEffects(system)
	if !manager.HasGrantedAbility(creatureID, "FlyingAbility") {
		t.Error("Permanent effect should survive end of combat")
	}
	
	// Only removing the enchantment should remove the ability
	manager.RemoveEffectsFromSource(enchantmentID)
	if manager.HasGrantedAbility(creatureID, "FlyingAbility") {
		t.Error("Flying should be gone after enchantment removed")
	}
}

// TestDynamicAbilityGranting_WhileOnBattlefield verifies source-dependent grants
func TestDynamicAbilityGranting_WhileOnBattlefield(t *testing.T) {
	system := NewLayerSystem()
	manager := NewEffectManager(system)
	
	creatureID := "grizzly-bears"
	lordID := "lord-of-the-skies"
	
	// Lord grants flying while it's on battlefield
	effect := NewEffectBuilder(lordID).
		Targeting(creatureID).
		WhileOnBattlefield().
		GrantAbility("FlyingAbility")
	
	manager.AddEffect(effect)
	
	// Creature should have flying
	if !manager.HasGrantedAbility(creatureID, "FlyingAbility") {
		t.Error("Creature should have flying while lord is on battlefield")
	}
	
	// Lord leaves battlefield
	CleanupSourceLeftBattlefieldEffects(system, lordID)
	
	// Flying should be gone
	if manager.HasGrantedAbility(creatureID, "FlyingAbility") {
		t.Error("Flying should be gone after lord leaves battlefield")
	}
}

// TestDynamicAbilityGranting_LayerOrdering verifies abilities are granted in Layer 6
func TestDynamicAbilityGranting_LayerOrdering(t *testing.T) {
	system := NewLayerSystem()
	
	creatureID := "grizzly-bears"
	
	// Create a grant ability effect
	effect := NewEffectBuilder("spell-1").
		Targeting(creatureID).
		UntilEndOfTurn().
		GrantAbility("FlyingAbility")
	
	// Verify it's in Layer 6 (Ability)
	if effect.Layer() != LayerAbility {
		t.Errorf("GrantAbilityEffect should be in LayerAbility, got %v", effect.Layer())
	}
	
	// Add to system
	system.AddEffect(effect)
	
	// Verify it's stored in the correct layer
	system.mu.RLock()
	layerEffects := system.effects[LayerAbility]
	system.mu.RUnlock()
	
	if len(layerEffects) != 1 {
		t.Errorf("Expected 1 effect in LayerAbility, got %d", len(layerEffects))
	}
}

// TestDynamicAbilityGranting_EffectStacking verifies multiple effects stack
func TestDynamicAbilityGranting_EffectStacking(t *testing.T) {
	system := NewLayerSystem()
	manager := NewEffectManager(system)
	
	creatureID := "grizzly-bears"
	
	// Two different spells both grant flying
	effect1 := NewEffectBuilder("spell-1").
		Targeting(creatureID).
		UntilEndOfTurn().
		GrantAbility("FlyingAbility")
	
	effect2 := NewEffectBuilder("spell-2").
		Targeting(creatureID).
		UntilEndOfCombat().
		GrantAbility("FlyingAbility")
	
	manager.AddEffect(effect1)
	manager.AddEffect(effect2)
	
	// Creature should have flying
	if !manager.HasGrantedAbility(creatureID, "FlyingAbility") {
		t.Error("Creature should have flying from stacked effects")
	}
	
	// End of combat - one effect expires
	CleanupEndOfCombatEffects(system)
	
	// Should still have flying from the other effect
	if !manager.HasGrantedAbility(creatureID, "FlyingAbility") {
		t.Error("Creature should still have flying from remaining effect")
	}
	
	// End of turn - other effect expires
	CleanupEndOfTurnEffects(system)
	
	// Now flying should be gone
	if manager.HasGrantedAbility(creatureID, "FlyingAbility") {
		t.Error("Flying should be gone after all effects expire")
	}
}
