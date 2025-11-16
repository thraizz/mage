package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestCombatFirstStrike tests that first strike creatures deal damage before normal creatures
func TestCombatFirstStrike(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-first-strike"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: 2/2 first strike attacker vs 2/2 normal blocker
	gameState.mu.Lock()
	attackerID := "attacker-1"
	blockerID := "blocker-1"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "First Strike Bear",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityFirstStrike, Text: "First strike"},
		},
	}
	
	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Normal Bear",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	
	// Check if combat has first strike
	hasFirstStrike, err := engine.HasFirstOrDoubleStrike(gameID)
	if err != nil {
		t.Fatalf("failed to check first strike: %v", err)
	}
	if !hasFirstStrike {
		t.Error("combat should have first strike")
	}
	
	// First strike damage step
	if err := engine.AssignCombatDamage(gameID, true); err != nil {
		t.Fatalf("failed to assign first strike damage: %v", err)
	}
	
	// Verify first striker dealt damage, blocker did not
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	attacker := gameState.cards[attackerID]
	
	if blocker.Damage != 2 {
		t.Errorf("blocker should have 2 damage from first strike, got %d", blocker.Damage)
	}
	if attacker.Damage != 0 {
		t.Errorf("attacker should have no damage yet (blocker hasn't dealt damage), got %d", attacker.Damage)
	}
	gameState.mu.RUnlock()
	
	// Apply first strike damage
	if err := engine.ApplyCombatDamage(gameID); err != nil {
		t.Fatalf("failed to apply first strike damage: %v", err)
	}
	
	// Verify blocker died
	gameState.mu.RLock()
	blocker = gameState.cards[blockerID]
	if blocker.Zone != zoneGraveyard {
		t.Error("blocker should be dead from first strike damage")
	}
	gameState.mu.RUnlock()
	
	// Normal damage step
	if err := engine.AssignCombatDamage(gameID, false); err != nil {
		t.Fatalf("failed to assign normal damage: %v", err)
	}
	
	// Verify attacker takes no additional damage (blocker is dead)
	gameState.mu.RLock()
	attacker = gameState.cards[attackerID]
	if attacker.Damage != 0 {
		t.Errorf("attacker should still have no damage (blocker died before dealing damage), got %d", attacker.Damage)
	}
	
	// Verify attacker survived
	if attacker.Zone != zoneBattlefield {
		t.Error("attacker should still be alive")
	}
	gameState.mu.RUnlock()
}

// TestCombatDoubleStrike tests that double strike creatures deal damage in both steps
func TestCombatDoubleStrike(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-double-strike"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: 2/2 double strike attacker vs 5/5 normal blocker
	gameState.mu.Lock()
	attackerID := "attacker-1"
	blockerID := "blocker-1"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Double Strike Bear",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityDoubleStrike, Text: "Double strike"},
		},
	}
	
	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Big Bear",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "5",
		Toughness:    "5",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	
	// First strike damage step
	if err := engine.AssignCombatDamage(gameID, true); err != nil {
		t.Fatalf("failed to assign first strike damage: %v", err)
	}
	
	// Verify double striker dealt damage in first strike step
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	attacker := gameState.cards[attackerID]
	
	if blocker.Damage != 2 {
		t.Errorf("blocker should have 2 damage from first strike, got %d", blocker.Damage)
	}
	if attacker.Damage != 0 {
		t.Errorf("attacker should have no damage yet, got %d", attacker.Damage)
	}
	gameState.mu.RUnlock()
	
	// Apply first strike damage (blocker survives with 5 toughness)
	if err := engine.ApplyCombatDamage(gameID); err != nil {
		t.Fatalf("failed to apply first strike damage: %v", err)
	}
	
	// Verify blocker survived
	gameState.mu.RLock()
	blocker = gameState.cards[blockerID]
	if blocker.Zone != zoneBattlefield {
		t.Error("blocker should still be alive after first strike")
	}
	gameState.mu.RUnlock()
	
	// Normal damage step
	if err := engine.AssignCombatDamage(gameID, false); err != nil {
		t.Fatalf("failed to assign normal damage: %v", err)
	}
	
	// Verify double striker dealt damage again AND blocker dealt damage
	gameState.mu.RLock()
	blocker = gameState.cards[blockerID]
	attacker = gameState.cards[attackerID]
	
	// Blocker should have 4 total damage (2 from first strike + 2 from normal)
	if blocker.Damage != 4 {
		t.Errorf("blocker should have 4 total damage, got %d", blocker.Damage)
	}
	
	// Attacker should have 5 damage from blocker
	if attacker.Damage != 5 {
		t.Errorf("attacker should have 5 damage from blocker, got %d", attacker.Damage)
	}
	gameState.mu.RUnlock()
	
	// Apply normal damage
	if err := engine.ApplyCombatDamage(gameID); err != nil {
		t.Fatalf("failed to apply normal damage: %v", err)
	}
	
	// Verify attacker died, blocker survived
	gameState.mu.RLock()
	blocker = gameState.cards[blockerID]
	attacker = gameState.cards[attackerID]
	
	if blocker.Zone != zoneBattlefield {
		t.Error("blocker should survive (4 damage < 5 toughness)")
	}
	if attacker.Zone != zoneGraveyard {
		t.Error("attacker should be dead (5 damage > 2 toughness)")
	}
	gameState.mu.RUnlock()
}

// TestCombatFirstStrikeVsFirstStrike tests first strike vs first strike (both deal damage simultaneously)
func TestCombatFirstStrikeVsFirstStrike(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-first-vs-first"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: 2/2 first strike vs 2/2 first strike
	gameState.mu.Lock()
	attackerID := "attacker-1"
	blockerID := "blocker-1"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "First Strike Bear",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityFirstStrike, Text: "First strike"},
		},
	}
	
	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "First Strike Bear",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityFirstStrike, Text: "First strike"},
		},
	}
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	
	// First strike damage step
	if err := engine.AssignCombatDamage(gameID, true); err != nil {
		t.Fatalf("failed to assign first strike damage: %v", err)
	}
	
	// Verify both dealt damage simultaneously
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	attacker := gameState.cards[attackerID]
	
	if blocker.Damage != 2 {
		t.Errorf("blocker should have 2 damage, got %d", blocker.Damage)
	}
	if attacker.Damage != 2 {
		t.Errorf("attacker should have 2 damage, got %d", attacker.Damage)
	}
	gameState.mu.RUnlock()
	
	// Apply damage
	if err := engine.ApplyCombatDamage(gameID); err != nil {
		t.Fatalf("failed to apply first strike damage: %v", err)
	}
	
	// Verify both died
	gameState.mu.RLock()
	blocker = gameState.cards[blockerID]
	attacker = gameState.cards[attackerID]
	
	if blocker.Zone != zoneGraveyard {
		t.Error("blocker should be dead")
	}
	if attacker.Zone != zoneGraveyard {
		t.Error("attacker should be dead")
	}
	gameState.mu.RUnlock()
	
	// Normal damage step should have no effect
	if err := engine.AssignCombatDamage(gameID, false); err != nil {
		t.Fatalf("failed to assign normal damage: %v", err)
	}
}

// TestCombatNoFirstStrike tests that HasFirstOrDoubleStrike returns false when no first strikers
func TestCombatNoFirstStrike(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-no-first-strike"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Normal creatures
	gameState.mu.Lock()
	attackerID := "attacker-1"
	blockerID := "blocker-1"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Normal Bear",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	
	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Normal Bear",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	
	// Check if combat has first strike
	hasFirstStrike, err := engine.HasFirstOrDoubleStrike(gameID)
	if err != nil {
		t.Fatalf("failed to check first strike: %v", err)
	}
	if hasFirstStrike {
		t.Error("combat should not have first strike")
	}
	
	// First strike step should deal no damage
	if err := engine.AssignCombatDamage(gameID, true); err != nil {
		t.Fatalf("failed to assign first strike damage: %v", err)
	}
	
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	attacker := gameState.cards[attackerID]
	
	if blocker.Damage != 0 {
		t.Errorf("blocker should have no damage in first strike step, got %d", blocker.Damage)
	}
	if attacker.Damage != 0 {
		t.Errorf("attacker should have no damage in first strike step, got %d", attacker.Damage)
	}
	gameState.mu.RUnlock()
}
