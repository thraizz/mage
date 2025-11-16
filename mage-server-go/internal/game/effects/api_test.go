package effects

import (
	"testing"
)

// TestEffectBuilder_GrantAbility verifies fluent API for ability granting
func TestEffectBuilder_GrantAbility(t *testing.T) {
	effect := NewEffectBuilder("source-1").
		Targeting("creature-1", "creature-2").
		UntilEndOfTurn().
		GrantAbility("FlyingAbility")
	
	if effect == nil {
		t.Fatal("Failed to create effect")
	}
	
	if effect.GetAbilityID() != "FlyingAbility" {
		t.Errorf("Expected FlyingAbility, got %s", effect.GetAbilityID())
	}
	
	if effect.GetDuration() != DurationEndOfTurn {
		t.Errorf("Expected EndOfTurn, got %v", effect.GetDuration())
	}
	
	targetIDs := effect.GetTargetIDs()
	if len(targetIDs) != 2 {
		t.Errorf("Expected 2 targets, got %d", len(targetIDs))
	}
}

// TestEffectBuilder_CantAttack verifies fluent API for attack restriction
func TestEffectBuilder_CantAttack(t *testing.T) {
	effect := NewEffectBuilder("source-1").
		Targeting("creature-1").
		UntilEndOfCombat().
		CantAttack()
	
	if effect == nil {
		t.Fatal("Failed to create effect")
	}
	
	if effect.GetDuration() != DurationEndOfCombat {
		t.Errorf("Expected EndOfCombat, got %v", effect.GetDuration())
	}
}

// TestEffectBuilder_CantBlock verifies fluent API for block restriction
func TestEffectBuilder_CantBlock(t *testing.T) {
	effect := NewEffectBuilder("source-1").
		Targeting("creature-1").
		Permanent().
		CantBlock()
	
	if effect == nil {
		t.Fatal("Failed to create effect")
	}
	
	if effect.GetDuration() != DurationPermanent {
		t.Errorf("Expected Permanent, got %v", effect.GetDuration())
	}
}

// TestEffectBuilder_MustAttack verifies fluent API for attack requirement
func TestEffectBuilder_MustAttack(t *testing.T) {
	effect := NewEffectBuilder("source-1").
		Targeting("creature-1").
		WhileOnBattlefield().
		MustAttack()
	
	if effect == nil {
		t.Fatal("Failed to create effect")
	}
	
	if effect.GetDuration() != DurationWhileOnBattlefield {
		t.Errorf("Expected WhileOnBattlefield, got %v", effect.GetDuration())
	}
}

// TestEffectManager_AddRemove verifies effect management
func TestEffectManager_AddRemove(t *testing.T) {
	system := NewLayerSystem()
	manager := NewEffectManager(system)
	
	effect := NewEffectBuilder("source-1").
		Targeting("creature-1").
		UntilEndOfTurn().
		GrantAbility("FlyingAbility")
	
	effectID := manager.AddEffect(effect)
	if effectID == "" {
		t.Fatal("Failed to add effect")
	}
	
	// Verify effect was added
	effects := manager.GetEffectsForCard("creature-1")
	if len(effects) != 1 {
		t.Errorf("Expected 1 effect, got %d", len(effects))
	}
	
	// Remove effect
	manager.RemoveEffect(effectID)
	
	// Verify effect was removed
	effects = manager.GetEffectsForCard("creature-1")
	if len(effects) != 0 {
		t.Errorf("Expected 0 effects after removal, got %d", len(effects))
	}
}

// TestEffectManager_RemoveBySource verifies source-based removal
func TestEffectManager_RemoveBySource(t *testing.T) {
	system := NewLayerSystem()
	manager := NewEffectManager(system)
	
	sourceID := "source-1"
	
	// Add multiple effects from same source
	effect1 := NewEffectBuilder(sourceID).
		Targeting("creature-1").
		UntilEndOfTurn().
		GrantAbility("FlyingAbility")
	
	effect2 := NewEffectBuilder(sourceID).
		Targeting("creature-2").
		UntilEndOfTurn().
		CantAttack()
	
	// Add effect from different source
	effect3 := NewEffectBuilder("source-2").
		Targeting("creature-3").
		UntilEndOfTurn().
		CantBlock()
	
	manager.AddEffect(effect1)
	manager.AddEffect(effect2)
	manager.AddEffect(effect3)
	
	// Verify all effects were added
	if len(manager.GetEffectsForCard("creature-1")) != 1 {
		t.Error("Effect 1 not added")
	}
	if len(manager.GetEffectsForCard("creature-2")) != 1 {
		t.Error("Effect 2 not added")
	}
	if len(manager.GetEffectsForCard("creature-3")) != 1 {
		t.Error("Effect 3 not added")
	}
	
	// Remove effects from source-1
	manager.RemoveEffectsFromSource(sourceID)
	
	// Verify only source-1 effects were removed
	if len(manager.GetEffectsForCard("creature-1")) != 0 {
		t.Error("Effect 1 should be removed")
	}
	if len(manager.GetEffectsForCard("creature-2")) != 0 {
		t.Error("Effect 2 should be removed")
	}
	if len(manager.GetEffectsForCard("creature-3")) != 1 {
		t.Error("Effect 3 should remain")
	}
}

// TestEffectManager_HasCantAttack verifies attack restriction checking
func TestEffectManager_HasCantAttack(t *testing.T) {
	system := NewLayerSystem()
	manager := NewEffectManager(system)
	
	creatureID := "creature-1"
	
	// Initially no restriction
	if manager.HasCantAttackEffect(creatureID) {
		t.Error("Should not have can't attack effect initially")
	}
	
	// Add can't attack effect
	effect := NewEffectBuilder("source-1").
		Targeting(creatureID).
		UntilEndOfTurn().
		CantAttack()
	
	manager.AddEffect(effect)
	
	// Should now have restriction
	if !manager.HasCantAttackEffect(creatureID) {
		t.Error("Should have can't attack effect after adding")
	}
}

// TestEffectManager_HasCantBlock verifies block restriction checking
func TestEffectManager_HasCantBlock(t *testing.T) {
	system := NewLayerSystem()
	manager := NewEffectManager(system)
	
	creatureID := "creature-1"
	
	// Initially no restriction
	if manager.HasCantBlockEffect(creatureID) {
		t.Error("Should not have can't block effect initially")
	}
	
	// Add can't block effect
	effect := NewEffectBuilder("source-1").
		Targeting(creatureID).
		UntilEndOfTurn().
		CantBlock()
	
	manager.AddEffect(effect)
	
	// Should now have restriction
	if !manager.HasCantBlockEffect(creatureID) {
		t.Error("Should have can't block effect after adding")
	}
}

// TestEffectManager_HasMustAttack verifies attack requirement checking
func TestEffectManager_HasMustAttack(t *testing.T) {
	system := NewLayerSystem()
	manager := NewEffectManager(system)
	
	creatureID := "creature-1"
	
	// Initially no requirement
	if manager.HasMustAttackEffect(creatureID) {
		t.Error("Should not have must attack effect initially")
	}
	
	// Add must attack effect
	effect := NewEffectBuilder("source-1").
		Targeting(creatureID).
		Permanent().
		MustAttack()
	
	manager.AddEffect(effect)
	
	// Should now have requirement
	if !manager.HasMustAttackEffect(creatureID) {
		t.Error("Should have must attack effect after adding")
	}
}

// TestEffectManager_HasGrantedAbility verifies granted ability checking
func TestEffectManager_HasGrantedAbility(t *testing.T) {
	system := NewLayerSystem()
	manager := NewEffectManager(system)
	
	creatureID := "creature-1"
	abilityID := "FlyingAbility"
	
	// Initially no granted ability
	if manager.HasGrantedAbility(creatureID, abilityID) {
		t.Error("Should not have granted ability initially")
	}
	
	// Add ability grant effect
	effect := NewEffectBuilder("source-1").
		Targeting(creatureID).
		UntilEndOfTurn().
		GrantAbility(abilityID)
	
	manager.AddEffect(effect)
	
	// Should now have granted ability
	if !manager.HasGrantedAbility(creatureID, abilityID) {
		t.Error("Should have granted ability after adding")
	}
	
	// Should not have different ability
	if manager.HasGrantedAbility(creatureID, "VigilanceAbility") {
		t.Error("Should not have vigilance ability")
	}
}
