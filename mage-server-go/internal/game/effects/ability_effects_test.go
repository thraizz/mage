package effects

import (
	"testing"
)

// TestGrantAbilityEffect_Basic verifies basic ability granting effect
func TestGrantAbilityEffect_Basic(t *testing.T) {
	effect := NewGrantAbilityEffect("source-1", "FlyingAbility", []string{"creature-1"}, "EndOfTurn")
	
	if effect == nil {
		t.Fatal("Failed to create grant ability effect")
	}
	
	if effect.ID() == "" {
		t.Error("Effect ID should not be empty")
	}
	
	if effect.Layer() != LayerAbility {
		t.Errorf("Expected layer %d, got %d", LayerAbility, effect.Layer())
	}
	
	if effect.GetAbilityID() != "FlyingAbility" {
		t.Errorf("Expected FlyingAbility, got %s", effect.GetAbilityID())
	}
}

// TestGrantAbilityEffect_AppliesTo verifies target matching
func TestGrantAbilityEffect_AppliesTo(t *testing.T) {
	effect := NewGrantAbilityEffect("source-1", "FlyingAbility", []string{"creature-1", "creature-2"}, "EndOfTurn")
	
	// Test matching target
	snapshot1 := &Snapshot{CardID: "creature-1"}
	if !effect.AppliesTo(snapshot1) {
		t.Error("Effect should apply to creature-1")
	}
	
	snapshot2 := &Snapshot{CardID: "creature-2"}
	if !effect.AppliesTo(snapshot2) {
		t.Error("Effect should apply to creature-2")
	}
	
	// Test non-matching target
	snapshot3 := &Snapshot{CardID: "creature-3"}
	if effect.AppliesTo(snapshot3) {
		t.Error("Effect should not apply to creature-3")
	}
}

// TestCantAttackEffect_Basic verifies basic can't attack effect
func TestCantAttackEffect_Basic(t *testing.T) {
	effect := NewCantAttackEffect("source-1", []string{"creature-1"}, "EndOfTurn")
	
	if effect == nil {
		t.Fatal("Failed to create can't attack effect")
	}
	
	if effect.ID() == "" {
		t.Error("Effect ID should not be empty")
	}
	
	// Verify target matching
	snapshot := &Snapshot{CardID: "creature-1"}
	if !effect.AppliesTo(snapshot) {
		t.Error("Effect should apply to creature-1")
	}
	
	snapshotOther := &Snapshot{CardID: "creature-2"}
	if effect.AppliesTo(snapshotOther) {
		t.Error("Effect should not apply to creature-2")
	}
}

// TestCantBlockEffect_Basic verifies basic can't block effect
func TestCantBlockEffect_Basic(t *testing.T) {
	effect := NewCantBlockEffect("source-1", []string{"creature-1"}, "EndOfTurn")
	
	if effect == nil {
		t.Fatal("Failed to create can't block effect")
	}
	
	if effect.ID() == "" {
		t.Error("Effect ID should not be empty")
	}
	
	// Verify target matching
	snapshot := &Snapshot{CardID: "creature-1"}
	if !effect.AppliesTo(snapshot) {
		t.Error("Effect should apply to creature-1")
	}
}

// TestMustAttackEffect_Basic verifies basic must attack effect
func TestMustAttackEffect_Basic(t *testing.T) {
	effect := NewMustAttackEffect("source-1", []string{"creature-1"}, "EndOfTurn")
	
	if effect == nil {
		t.Fatal("Failed to create must attack effect")
	}
	
	if effect.ID() == "" {
		t.Error("Effect ID should not be empty")
	}
	
	// Verify target matching
	snapshot := &Snapshot{CardID: "creature-1"}
	if !effect.AppliesTo(snapshot) {
		t.Error("Effect should apply to creature-1")
	}
}

// TestLayerSystem_Integration verifies integration with existing layer system
func TestLayerSystem_Integration(t *testing.T) {
	system := NewLayerSystem()
	
	// Add a PT boost effect (existing)
	ptEffect := NewSimplePTBoostEffect("source-1", "Alice", 1, 1, false)
	system.AddEffect(ptEffect)
	
	// Add an ability grant effect (new)
	abilityEffect := NewGrantAbilityEffect("source-2", "FlyingAbility", []string{"creature-1"}, "EndOfTurn")
	system.AddEffect(abilityEffect)
	
	// Create a snapshot
	snapshot := NewSnapshot("creature-1", "Alice", []string{"Creature"}, 2, 2, true, true)
	
	// Apply effects
	system.Apply(snapshot)
	
	// Verify PT boost was applied (Layer 7)
	if snapshot.Power != 3 || snapshot.Toughness != 3 {
		t.Errorf("Expected 3/3, got %d/%d", snapshot.Power, snapshot.Toughness)
	}
	
	// Note: Ability granting would need more infrastructure to test fully
}
