package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestCombatEndCombat tests that EndCombat clears combat state properly
func TestCombatEndCombat(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-end-combat"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Create attacker and blocker
	gameState.mu.Lock()
	attackerID := "attacker-1"
	blockerID := "blocker-1"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Grizzly Bears",
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
		Name:         "Wall of Stone",
		Type:         "Creature - Wall",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "0",
		Toughness:    "8",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	
	// Verify combat is active
	gameState.mu.RLock()
	attacker := gameState.cards[attackerID]
	blocker := gameState.cards[blockerID]
	
	if !attacker.Attacking {
		t.Error("attacker should be attacking before EndCombat")
	}
	if !blocker.Blocking {
		t.Error("blocker should be blocking before EndCombat")
	}
	if len(gameState.combat.groups) == 0 {
		t.Error("combat groups should exist before EndCombat")
	}
	gameState.mu.RUnlock()
	
	// End combat
	if err := engine.EndCombat(gameID); err != nil {
		t.Fatalf("failed to end combat: %v", err)
	}
	
	// Verify combat flags are cleared
	gameState.mu.RLock()
	attacker = gameState.cards[attackerID]
	blocker = gameState.cards[blockerID]
	
	if attacker.Attacking {
		t.Error("attacker should not be attacking after EndCombat")
	}
	if attacker.Blocking {
		t.Error("attacker should not be blocking after EndCombat")
	}
	if attacker.AttackingWhat != "" {
		t.Error("attacker.AttackingWhat should be cleared after EndCombat")
	}
	if attacker.BlockingWhat != nil {
		t.Error("attacker.BlockingWhat should be cleared after EndCombat")
	}
	
	if blocker.Attacking {
		t.Error("blocker should not be attacking after EndCombat")
	}
	if blocker.Blocking {
		t.Error("blocker should not be blocking after EndCombat")
	}
	if blocker.AttackingWhat != "" {
		t.Error("blocker.AttackingWhat should be cleared after EndCombat")
	}
	if blocker.BlockingWhat != nil {
		t.Error("blocker.BlockingWhat should be cleared after EndCombat")
	}
	
	// Verify combat groups are cleared
	if len(gameState.combat.groups) != 0 {
		t.Errorf("combat groups should be cleared, got %d groups", len(gameState.combat.groups))
	}
	
	// Verify former groups are preserved
	if len(gameState.combat.formerGroups) == 0 {
		t.Error("former groups should be preserved after EndCombat")
	}
	
	gameState.mu.RUnlock()
}

// TestCombatEndCombatDamageCleared tests that damage is cleared after combat
func TestCombatEndCombatDamageCleared(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-end-combat-damage"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Create attacker and blocker
	gameState.mu.Lock()
	attackerID := "attacker-1"
	blockerID := "blocker-1"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Grizzly Bears",
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
		Name:         "Wall of Stone",
		Type:         "Creature - Wall",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "0",
		Toughness:    "8",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Setup combat and assign damage
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AssignCombatDamage(gameID, false)
	
	// Verify damage is marked
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	if blocker.Damage == 0 {
		t.Error("blocker should have damage before EndCombat")
	}
	gameState.mu.RUnlock()
	
	// End combat
	if err := engine.EndCombat(gameID); err != nil {
		t.Fatalf("failed to end combat: %v", err)
	}
	
	// Verify damage is cleared
	gameState.mu.RLock()
	blocker = gameState.cards[blockerID]
	if blocker.Damage != 0 {
		t.Errorf("blocker damage should be cleared after EndCombat, got %d", blocker.Damage)
	}
	if blocker.DamageSources != nil {
		t.Error("blocker damage sources should be cleared after EndCombat")
	}
	gameState.mu.RUnlock()
}

// TestCombatAttackedThisTurn tests GetAttackedThisTurn tracking
func TestCombatAttackedThisTurn(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-attacked-this-turn"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Create attacker
	gameState.mu.Lock()
	attackerID := "attacker-1"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Grizzly Bears",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Before combat - should not have attacked
	attacked, err := engine.GetAttackedThisTurn(gameID, attackerID)
	if err != nil {
		t.Fatalf("failed to check attacked this turn: %v", err)
	}
	if attacked {
		t.Error("creature should not have attacked before combat")
	}
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	
	// During combat - should have attacked
	attacked, err = engine.GetAttackedThisTurn(gameID, attackerID)
	if err != nil {
		t.Fatalf("failed to check attacked this turn: %v", err)
	}
	if !attacked {
		t.Error("creature should have attacked during combat")
	}
	
	// End combat
	if err := engine.EndCombat(gameID); err != nil {
		t.Fatalf("failed to end combat: %v", err)
	}
	
	// After combat - should still show as attacked this turn
	attacked, err = engine.GetAttackedThisTurn(gameID, attackerID)
	if err != nil {
		t.Fatalf("failed to check attacked this turn: %v", err)
	}
	if !attacked {
		t.Error("creature should still show as attacked after combat ends")
	}
}

// TestCombatEndCombatMultipleGroups tests EndCombat with multiple combat groups
func TestCombatEndCombatMultipleGroups(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-end-combat-multi"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Create multiple attackers
	gameState.mu.Lock()
	attacker1ID := "attacker-1"
	attacker2ID := "attacker-2"
	
	gameState.cards[attacker1ID] = &internalCard{
		ID:           attacker1ID,
		Name:         "Grizzly Bears",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	
	gameState.cards[attacker2ID] = &internalCard{
		ID:           attacker2ID,
		Name:         "Llanowar Elves",
		Type:         "Creature - Elf",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "1",
		Toughness:    "1",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Setup combat with multiple attackers
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attacker1ID, "Bob", "Alice")
	engine.DeclareAttacker(gameID, attacker2ID, "Bob", "Alice")
	
	// Verify both are attacking
	gameState.mu.RLock()
	if !gameState.cards[attacker1ID].Attacking {
		t.Error("attacker 1 should be attacking")
	}
	if !gameState.cards[attacker2ID].Attacking {
		t.Error("attacker 2 should be attacking")
	}
	groupCount := len(gameState.combat.groups)
	gameState.mu.RUnlock()
	
	// End combat
	if err := engine.EndCombat(gameID); err != nil {
		t.Fatalf("failed to end combat: %v", err)
	}
	
	// Verify both attackers are cleared
	gameState.mu.RLock()
	if gameState.cards[attacker1ID].Attacking {
		t.Error("attacker 1 should not be attacking after EndCombat")
	}
	if gameState.cards[attacker2ID].Attacking {
		t.Error("attacker 2 should not be attacking after EndCombat")
	}
	
	// Verify former groups preserved all groups
	if len(gameState.combat.formerGroups) != groupCount {
		t.Errorf("expected %d former groups, got %d", groupCount, len(gameState.combat.formerGroups))
	}
	
	// Both should show as attacked this turn
	gameState.mu.RUnlock()
	
	attacked1, _ := engine.GetAttackedThisTurn(gameID, attacker1ID)
	attacked2, _ := engine.GetAttackedThisTurn(gameID, attacker2ID)
	
	if !attacked1 {
		t.Error("attacker 1 should show as attacked this turn")
	}
	if !attacked2 {
		t.Error("attacker 2 should show as attacked this turn")
	}
}
