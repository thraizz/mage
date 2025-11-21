package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestCombatDamageUnblocked tests damage to player from unblocked attacker
func TestCombatDamageUnblocked(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-unblocked-damage"
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

	// Record initial life
	initialLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Assign and apply damage
	if err := engine.AssignCombatDamage(gameID, false); err != nil {
		t.Fatalf("failed to assign combat damage: %v", err)
	}

	if err := engine.ApplyCombatDamage(gameID); err != nil {
		t.Fatalf("failed to apply combat damage: %v", err)
	}

	// Verify Bob took damage
	gameState.mu.RLock()
	bobLife := gameState.players["Bob"].Life
	gameState.mu.RUnlock()

	expectedLife := initialLife - 2
	if bobLife != expectedLife {
		t.Errorf("expected Bob's life to be %d, got %d", expectedLife, bobLife)
	}
}

// TestCombatDamageBlocked tests damage assignment when attacker is blocked
func TestCombatDamageBlocked(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-blocked-damage"
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

	// Record initial life
	initialLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()

	// Setup combat with blocker
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")

	// Assign and apply damage
	if err := engine.AssignCombatDamage(gameID, false); err != nil {
		t.Fatalf("failed to assign combat damage: %v", err)
	}

	if err := engine.ApplyCombatDamage(gameID); err != nil {
		t.Fatalf("failed to apply combat damage: %v", err)
	}

	// Verify Bob took no damage (attacker was blocked)
	gameState.mu.RLock()
	bobLife := gameState.players["Bob"].Life
	if bobLife != initialLife {
		t.Errorf("expected Bob's life to remain %d, got %d", initialLife, bobLife)
	}

	// Verify blocker took damage
	blocker := gameState.cards[blockerID]
	if blocker.Damage != 2 {
		t.Errorf("expected blocker to have 2 damage, got %d", blocker.Damage)
	}

	// Verify blocker is still alive (0/8 with 2 damage)
	if blocker.Zone != zoneBattlefield {
		t.Error("blocker should still be on battlefield")
	}
	gameState.mu.RUnlock()
}

// TestCombatDamageLethal tests creature death from lethal damage
func TestCombatDamageLethal(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-lethal-damage"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create attacker and blocker (both 2/2)
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
		Name:         "Grizzly Bears",
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

	// Assign and apply damage
	if err := engine.AssignCombatDamage(gameID, false); err != nil {
		t.Fatalf("failed to assign combat damage: %v", err)
	}

	if err := engine.ApplyCombatDamage(gameID); err != nil {
		t.Fatalf("failed to apply combat damage: %v", err)
	}

	// Verify both creatures died (2 damage >= 2 toughness)
	gameState.mu.RLock()
	attacker, attackerExists := gameState.cards[attackerID]
	blocker, blockerExists := gameState.cards[blockerID]

	if !attackerExists || attacker.Zone != zoneGraveyard {
		t.Error("attacker should be in graveyard")
	}

	if !blockerExists || blocker.Zone != zoneGraveyard {
		t.Error("blocker should be in graveyard")
	}
	gameState.mu.RUnlock()
}

// TestCombatDamageMultipleBlockers tests damage split among multiple blockers
func TestCombatDamageMultipleBlockers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-multi-blocker-damage"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create one 4/4 attacker and two 2/2 blockers
	gameState.mu.Lock()
	attackerID := "attacker-1"
	blocker1ID := "blocker-1"
	blocker2ID := "blocker-2"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Serra Angel",
		Type:         "Creature - Angel",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "4",
		Toughness:    "4",
		Tapped:       false,
	}

	gameState.cards[blocker1ID] = &internalCard{
		ID:           blocker1ID,
		Name:         "Grizzly Bears",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}

	gameState.cards[blocker2ID] = &internalCard{
		ID:           blocker2ID,
		Name:         "Llanowar Elves",
		Type:         "Creature - Elf",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "1",
		Toughness:    "1",
		Tapped:       false,
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blocker1ID, attackerID, "Bob")
	engine.DeclareBlocker(gameID, blocker2ID, attackerID, "Bob")

	// Assign and apply damage
	if err := engine.AssignCombatDamage(gameID, false); err != nil {
		t.Fatalf("failed to assign combat damage: %v", err)
	}

	if err := engine.ApplyCombatDamage(gameID); err != nil {
		t.Fatalf("failed to apply combat damage: %v", err)
	}

	// Verify damage was split (4 power / 2 blockers = 2 each)
	gameState.mu.RLock()
	blocker1 := gameState.cards[blocker1ID]
	blocker2 := gameState.cards[blocker2ID]

	totalDamage := blocker1.Damage + blocker2.Damage
	if totalDamage != 4 {
		t.Errorf("expected total damage to be 4, got %d", totalDamage)
	}

	// Both blockers should have taken damage
	if blocker1.Damage == 0 || blocker2.Damage == 0 {
		t.Error("both blockers should have taken damage")
	}

	// Verify attacker took damage from both blockers (2+1=3)
	attacker := gameState.cards[attackerID]
	if attacker.Damage != 3 {
		t.Errorf("expected attacker to have 3 damage, got %d", attacker.Damage)
	}
	gameState.mu.RUnlock()
}

// TestCombatDamageFirstStrike tests that first strike doesn't deal damage yet
func TestCombatDamageFirstStrike(t *testing.T) {
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

	initialLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Assign damage with firstStrike=true (should not deal damage yet)
	if err := engine.AssignCombatDamage(gameID, true); err != nil {
		t.Fatalf("failed to assign combat damage: %v", err)
	}

	if err := engine.ApplyCombatDamage(gameID); err != nil {
		t.Fatalf("failed to apply combat damage: %v", err)
	}

	// Verify no damage was dealt (first strike not implemented yet)
	gameState.mu.RLock()
	bobLife := gameState.players["Bob"].Life
	if bobLife != initialLife {
		t.Errorf("expected no damage during first strike step, Bob's life changed from %d to %d", initialLife, bobLife)
	}
	gameState.mu.RUnlock()
}
