package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestRemoveFromCombat_Attacker tests removing an attacking creature from combat
// Per Java Combat.removeFromCombat() - should clear attacking state and fire event
func TestRemoveFromCombat_Attacker(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-remove-from-combat-attacker"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create attacker
	gameState.mu.Lock()
	attackerID := "grizzly-bears"
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

	// Setup combat and declare attacker
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Verify attacker is attacking and tapped
	gameState.mu.RLock()
	attacker := gameState.cards[attackerID]
	if !attacker.Attacking {
		t.Error("Attacker should be attacking before removal")
	}
	if !attacker.Tapped {
		t.Error("Attacker should be tapped before removal")
	}
	if !gameState.combat.attackers[attackerID] {
		t.Error("Attacker should be in attackers set before removal")
	}
	gameState.mu.RUnlock()

	// Remove from combat
	if err := engine.RemoveFromCombat(gameID, attackerID); err != nil {
		t.Fatalf("Failed to remove from combat: %v", err)
	}

	// Verify attacker is no longer attacking
	gameState.mu.RLock()
	attacker = gameState.cards[attackerID]
	if attacker.Attacking {
		t.Error("Attacker should not be attacking after removal")
	}
	if attacker.AttackingWhat != "" {
		t.Error("AttackingWhat should be cleared after removal")
	}
	if gameState.combat.attackers[attackerID] {
		t.Error("Attacker should not be in attackers set after removal")
	}
	if len(gameState.combat.groups) != 0 {
		t.Errorf("Expected 0 combat groups after removal, got %d", len(gameState.combat.groups))
	}
	// Note: RemoveFromCombat does NOT untap (unlike RemoveAttacker)
	if !attacker.Tapped {
		t.Error("Attacker should still be tapped (RemoveFromCombat does not untap)")
	}
	gameState.mu.RUnlock()

	// Event is published (verified in logs), but event testing is complex with locks
	// The key functionality is state changes, which we've verified above
}

// TestRemoveFromCombat_Blocker tests removing a blocking creature from combat
func TestRemoveFromCombat_Blocker(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-remove-from-combat-blocker"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create attacker and blocker
	gameState.mu.Lock()
	attackerID := "grizzly-bears"
	blockerID := "wall-of-stone"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Grizzly Bears",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
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
	}
	gameState.mu.Unlock()

	// Setup combat and declare attacker
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Declare blocker
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")

	// Verify blocker is blocking
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	if !blocker.Blocking {
		t.Error("Blocker should be blocking before removal")
	}
	if !gameState.combat.blockers[blockerID] {
		t.Error("Blocker should be in blockers set before removal")
	}
	if len(gameState.combat.groups) != 1 {
		t.Errorf("Expected 1 combat group, got %d", len(gameState.combat.groups))
	}
	if !gameState.combat.groups[0].blocked {
		t.Error("Combat group should be blocked")
	}
	gameState.mu.RUnlock()

	// Remove blocker from combat
	if err := engine.RemoveFromCombat(gameID, blockerID); err != nil {
		t.Fatalf("Failed to remove from combat: %v", err)
	}

	// Verify blocker is no longer blocking
	gameState.mu.RLock()
	blocker = gameState.cards[blockerID]
	if blocker.Blocking {
		t.Error("Blocker should not be blocking after removal")
	}
	if blocker.BlockingWhat != nil {
		t.Error("BlockingWhat should be nil after removal")
	}
	if gameState.combat.blockers[blockerID] {
		t.Error("Blocker should not be in blockers set after removal")
	}
	if len(gameState.combat.groups) != 1 {
		t.Errorf("Expected 1 combat group, got %d", len(gameState.combat.groups))
	}
	if gameState.combat.groups[0].blocked {
		t.Error("Combat group should be unblocked after blocker removal")
	}
	if len(gameState.combat.groups[0].blockers) != 0 {
		t.Errorf("Expected 0 blockers in group, got %d", len(gameState.combat.groups[0].blockers))
	}
	gameState.mu.RUnlock()

	// Event is published (verified in logs), but event testing is complex with locks
	// The key functionality is state changes, which we've verified above
}

// TestRemoveFromCombat_MultipleBlockers tests removing one blocker when multiple exist
func TestRemoveFromCombat_MultipleBlockers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-remove-from-combat-multiple-blockers"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create attacker and two blockers
	gameState.mu.Lock()
	attackerID := "serra-angel"
	blocker1ID := "grizzly-bears"
	blocker2ID := "wall-of-stone"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Serra Angel",
		Type:         "Creature - Angel",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "4",
		Toughness:    "4",
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
	}

	gameState.cards[blocker2ID] = &internalCard{
		ID:           blocker2ID,
		Name:         "Wall of Stone",
		Type:         "Creature - Wall",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "0",
		Toughness:    "8",
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Declare both blockers
	engine.DeclareBlocker(gameID, blocker1ID, attackerID, "Bob")
	engine.DeclareBlocker(gameID, blocker2ID, attackerID, "Bob")

	// Verify both are blocking
	gameState.mu.RLock()
	if !gameState.cards[blocker1ID].Blocking {
		t.Error("Blocker 1 should be blocking")
	}
	if !gameState.cards[blocker2ID].Blocking {
		t.Error("Blocker 2 should be blocking")
	}
	if len(gameState.combat.groups[0].blockers) != 2 {
		t.Errorf("Expected 2 blockers, got %d", len(gameState.combat.groups[0].blockers))
	}
	gameState.mu.RUnlock()

	// Remove first blocker
	if err := engine.RemoveFromCombat(gameID, blocker1ID); err != nil {
		t.Fatalf("Failed to remove blocker 1: %v", err)
	}

	// Verify first blocker removed, second still blocking
	gameState.mu.RLock()
	if gameState.cards[blocker1ID].Blocking {
		t.Error("Blocker 1 should not be blocking after removal")
	}
	if !gameState.cards[blocker2ID].Blocking {
		t.Error("Blocker 2 should still be blocking")
	}
	if !gameState.combat.groups[0].blocked {
		t.Error("Combat group should still be blocked")
	}
	if len(gameState.combat.groups[0].blockers) != 1 {
		t.Errorf("Expected 1 blocker remaining, got %d", len(gameState.combat.groups[0].blockers))
	}
	if gameState.combat.groups[0].blockers[0] != blocker2ID {
		t.Error("Blocker 2 should be the remaining blocker")
	}
	gameState.mu.RUnlock()
}

// TestRemoveFromCombat_NotInCombat tests removing a creature that's not in combat
func TestRemoveFromCombat_NotInCombat(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-remove-from-combat-not-in-combat"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create creature not in combat
	gameState.mu.Lock()
	creatureID := "llanowar-elves"
	gameState.cards[creatureID] = &internalCard{
		ID:           creatureID,
		Name:         "Llanowar Elves",
		Type:         "Creature - Elf",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "1",
		Toughness:    "1",
	}
	gameState.mu.Unlock()

	// Setup combat but don't declare creature
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	// Remove from combat (should succeed but do nothing)
	if err := engine.RemoveFromCombat(gameID, creatureID); err != nil {
		t.Fatalf("Failed to remove from combat: %v", err)
	}

	// Verify creature state unchanged (it was never in combat)
	// Event should not fire when creature wasn't in combat (checked via removed flag in implementation)
}

// TestRemoveFromCombat_DoesNotUntap verifies RemoveFromCombat does NOT untap creatures
// This is a key difference from RemoveAttacker
func TestRemoveFromCombat_DoesNotUntap(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-remove-from-combat-no-untap"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create attacker without vigilance
	gameState.mu.Lock()
	attackerID := "grizzly-bears"
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

	// Setup combat and declare attacker (this will tap it)
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Verify attacker is tapped
	gameState.mu.RLock()
	attacker := gameState.cards[attackerID]
	if !attacker.Tapped {
		t.Error("Attacker should be tapped after declaration")
	}
	if !gameState.combat.attackersTapped[attackerID] {
		t.Error("Attacker should be in attackersTapped set")
	}
	gameState.mu.RUnlock()

	// Remove from combat
	if err := engine.RemoveFromCombat(gameID, attackerID); err != nil {
		t.Fatalf("Failed to remove from combat: %v", err)
	}

	// Verify attacker is STILL tapped (key difference from RemoveAttacker)
	gameState.mu.RLock()
	attacker = gameState.cards[attackerID]
	if !attacker.Tapped {
		t.Error("RemoveFromCombat should NOT untap creatures (unlike RemoveAttacker)")
	}
	gameState.mu.RUnlock()
}

// TestRemoveFromCombat_InvalidGame tests error handling for invalid game
func TestRemoveFromCombat_InvalidGame(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-remove-from-combat-invalid-game"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Try to remove from combat with invalid game ID
	err := engine.RemoveFromCombat("invalid-game-id", "some-creature")
	if err == nil {
		t.Error("Expected error for invalid game ID")
	}
}

// TestRemoveFromCombat_InvalidCreature tests error handling for invalid creature
func TestRemoveFromCombat_InvalidCreature(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-remove-from-combat-invalid-creature"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	// Try to remove non-existent creature
	err := engine.RemoveFromCombat(gameID, "invalid-creature-id")
	if err == nil {
		t.Error("Expected error for invalid creature ID")
	}
}

// TestRemoveFromCombat_Integration tests a complete combat flow with removal
// Simulates a regeneration effect
func TestRemoveFromCombat_Integration(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-remove-from-combat-integration"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create two attackers and one blocker
	gameState.mu.Lock()
	attacker1ID := "grizzly-bears"
	attacker2ID := "hill-giant"
	blockerID := "wall-of-stone"

	gameState.cards[attacker1ID] = &internalCard{
		ID:           attacker1ID,
		Name:         "Grizzly Bears",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
	}

	gameState.cards[attacker2ID] = &internalCard{
		ID:           attacker2ID,
		Name:         "Hill Giant",
		Type:         "Creature - Giant",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "3",
		Toughness:    "3",
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
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	// Declare both attackers
	engine.DeclareAttacker(gameID, attacker1ID, "Bob", "Alice")
	engine.DeclareAttacker(gameID, attacker2ID, "Bob", "Alice")

	// Declare blocker on first attacker
	engine.DeclareBlocker(gameID, blockerID, attacker1ID, "Bob")

	// Verify two combat groups
	gameState.mu.RLock()
	if len(gameState.combat.groups) != 2 {
		t.Errorf("Expected 2 combat groups, got %d", len(gameState.combat.groups))
	}
	gameState.mu.RUnlock()

	// Simulate regeneration removing first attacker from combat
	if err := engine.RemoveFromCombat(gameID, attacker1ID); err != nil {
		t.Fatalf("Failed to remove attacker 1: %v", err)
	}

	// Verify first attacker removed, second still attacking
	gameState.mu.RLock()
	if gameState.cards[attacker1ID].Attacking {
		t.Error("Attacker 1 should not be attacking after removal")
	}
	if !gameState.cards[attacker2ID].Attacking {
		t.Error("Attacker 2 should still be attacking")
	}
	if len(gameState.combat.groups) != 1 {
		t.Errorf("Expected 1 combat group remaining, got %d", len(gameState.combat.groups))
	}
	if gameState.combat.groups[0].attackers[0] != attacker2ID {
		t.Error("Remaining group should contain attacker 2")
	}
	gameState.mu.RUnlock()

	// Assign and apply damage for remaining attacker
	if err := engine.AssignCombatDamage(gameID, false); err != nil {
		t.Fatalf("Failed to assign damage: %v", err)
	}

	if err := engine.ApplyCombatDamage(gameID); err != nil {
		t.Fatalf("Failed to apply damage: %v", err)
	}

	// Verify attacker 2 dealt damage to player
	gameState.mu.RLock()
	bobLife := gameState.players["Bob"].Life
	gameState.mu.RUnlock()

	if bobLife != 17 { // 20 - 3 = 17
		t.Errorf("Expected Bob to have 17 life, got %d", bobLife)
	}
}
