package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestCheckForRemoveFromCombat_AttackerLosesCreatureType tests automatic removal when attacker loses creature type
// Per Java RemoveFromCombatTest.test_LeavesCombatIfNoLongerACreature() - validates creature type loss during combat
func TestCheckForRemoveFromCombat_AttackerLosesCreatureType(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-type-loss-attacker"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create an attacker
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
	}
	gameState.mu.Unlock()

	// Setup combat and declare attacker
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Verify attacker is attacking
	gameState.mu.RLock()
	attacker := gameState.cards[attackerID]
	if !attacker.Attacking {
		t.Error("Attacker should be attacking before type loss")
	}
	if !gameState.combat.attackers[attackerID] {
		t.Error("Attacker should be in attackers set before type loss")
	}
	gameState.mu.RUnlock()

	// Creature loses creature type (e.g., Ambush Commander dies)
	gameState.mu.Lock()
	attacker.Type = "Land - Forest" // No longer a creature
	gameState.mu.Unlock()

	// Check for removal (this would happen automatically during combat steps)
	if err := engine.CheckForRemoveFromCombat(gameID); err != nil {
		t.Fatalf("CheckForRemoveFromCombat failed: %v", err)
	}

	// Verify attacker was removed from combat
	gameState.mu.RLock()
	attacker = gameState.cards[attackerID]
	if attacker.Attacking {
		t.Error("Attacker should not be attacking after losing creature type")
	}
	if gameState.combat.attackers[attackerID] {
		t.Error("Attacker should not be in attackers set after losing creature type")
	}
	if len(gameState.combat.groups) != 0 {
		t.Errorf("Expected 0 combat groups after removal, got %d", len(gameState.combat.groups))
	}
	gameState.mu.RUnlock()
}

// TestCheckForRemoveFromCombat_BlockerLosesCreatureType tests automatic removal when blocker loses creature type
func TestCheckForRemoveFromCombat_BlockerLosesCreatureType(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-type-loss-blocker"
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

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")

	// Verify blocker is blocking
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	if !blocker.Blocking {
		t.Error("Blocker should be blocking before type loss")
	}
	if !gameState.combat.blockers[blockerID] {
		t.Error("Blocker should be in blockers set before type loss")
	}
	if !gameState.combat.groups[0].blocked {
		t.Error("Combat group should be blocked before type loss")
	}
	gameState.mu.RUnlock()

	// Blocker loses creature type
	gameState.mu.Lock()
	blocker.Type = "Artifact" // No longer a creature
	gameState.mu.Unlock()

	// Check for removal
	if err := engine.CheckForRemoveFromCombat(gameID); err != nil {
		t.Fatalf("CheckForRemoveFromCombat failed: %v", err)
	}

	// Verify blocker was removed from combat
	gameState.mu.RLock()
	blocker = gameState.cards[blockerID]
	if blocker.Blocking {
		t.Error("Blocker should not be blocking after losing creature type")
	}
	if gameState.combat.blockers[blockerID] {
		t.Error("Blocker should not be in blockers set after losing creature type")
	}
	if gameState.combat.groups[0].blocked {
		t.Error("Combat group should be unblocked after blocker loses creature type")
	}
	if len(gameState.combat.groups[0].blockers) != 0 {
		t.Errorf("Expected 0 blockers, got %d", len(gameState.combat.groups[0].blockers))
	}
	gameState.mu.RUnlock()
}

// TestCheckForRemoveFromCombat_MultipleCreaturesLoseType tests removal of multiple creatures simultaneously
func TestCheckForRemoveFromCombat_MultipleCreaturesLoseType(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-type-loss-multiple"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create multiple attackers
	gameState.mu.Lock()
	attacker1ID := "animated-forest-1"
	attacker2ID := "animated-forest-2"
	attacker3ID := "grizzly-bears"

	// All three start as creatures
	for _, id := range []string{attacker1ID, attacker2ID, attacker3ID} {
		gameState.cards[id] = &internalCard{
			ID:           id,
			Name:         "Animated Forest",
			Type:         "Land Creature - Forest",
			Zone:         zoneBattlefield,
			OwnerID:      "Alice",
			ControllerID: "Alice",
			Power:        "1",
			Toughness:    "1",
		}
	}
	gameState.mu.Unlock()

	// Setup combat and declare all attackers
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attacker1ID, "Bob", "Alice")
	engine.DeclareAttacker(gameID, attacker2ID, "Bob", "Alice")
	engine.DeclareAttacker(gameID, attacker3ID, "Bob", "Alice")

	// Verify all are attacking
	gameState.mu.RLock()
	if len(gameState.combat.attackers) != 3 {
		t.Errorf("Expected 3 attackers, got %d", len(gameState.combat.attackers))
	}
	if len(gameState.combat.groups) != 3 {
		t.Errorf("Expected 3 combat groups, got %d", len(gameState.combat.groups))
	}
	gameState.mu.RUnlock()

	// Two creatures lose creature type (e.g., Ambush Commander dies)
	gameState.mu.Lock()
	gameState.cards[attacker1ID].Type = "Land - Forest"
	gameState.cards[attacker2ID].Type = "Land - Forest"
	// attacker3 remains a creature
	gameState.mu.Unlock()

	// Check for removal
	if err := engine.CheckForRemoveFromCombat(gameID); err != nil {
		t.Fatalf("CheckForRemoveFromCombat failed: %v", err)
	}

	// Verify only the two non-creatures were removed
	gameState.mu.RLock()
	if gameState.combat.attackers[attacker1ID] {
		t.Error("Attacker 1 should not be in combat after losing creature type")
	}
	if gameState.combat.attackers[attacker2ID] {
		t.Error("Attacker 2 should not be in combat after losing creature type")
	}
	if !gameState.combat.attackers[attacker3ID] {
		t.Error("Attacker 3 should still be in combat")
	}
	if len(gameState.combat.attackers) != 1 {
		t.Errorf("Expected 1 attacker remaining, got %d", len(gameState.combat.attackers))
	}
	if len(gameState.combat.groups) != 1 {
		t.Errorf("Expected 1 combat group remaining, got %d", len(gameState.combat.groups))
	}
	gameState.mu.RUnlock()
}

// TestCheckForRemoveFromCombat_NoRemovalIfStillCreature tests that creatures remain in combat if they keep creature type
func TestCheckForRemoveFromCombat_NoRemovalIfStillCreature(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-type-no-loss"
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
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Change type but keep creature
	gameState.mu.Lock()
	gameState.cards[attackerID].Type = "Artifact Creature - Bear Golem"
	gameState.mu.Unlock()

	// Check for removal
	if err := engine.CheckForRemoveFromCombat(gameID); err != nil {
		t.Fatalf("CheckForRemoveFromCombat failed: %v", err)
	}

	// Verify attacker is still in combat
	gameState.mu.RLock()
	if !gameState.combat.attackers[attackerID] {
		t.Error("Attacker should still be in combat with modified creature type")
	}
	if !gameState.cards[attackerID].Attacking {
		t.Error("Attacker should still be attacking with modified creature type")
	}
	gameState.mu.RUnlock()
}

// TestCheckForRemoveFromCombat_EmptyGame tests check on empty game state
func TestCheckForRemoveFromCombat_EmptyGame(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-type-loss-empty"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Check with no combat active
	if err := engine.CheckForRemoveFromCombat(gameID); err != nil {
		t.Fatalf("CheckForRemoveFromCombat should not fail on empty combat: %v", err)
	}

	// Verify state unchanged
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	gameState.mu.RLock()
	if len(gameState.combat.attackers) != 0 {
		t.Errorf("Expected 0 attackers, got %d", len(gameState.combat.attackers))
	}
	if len(gameState.combat.blockers) != 0 {
		t.Errorf("Expected 0 blockers, got %d", len(gameState.combat.blockers))
	}
	gameState.mu.RUnlock()
}

// TestCheckForRemoveFromCombat_InvalidGame tests error handling
func TestCheckForRemoveFromCombat_InvalidGame(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	// Try to check non-existent game
	err := engine.CheckForRemoveFromCombat("non-existent-game")
	if err == nil {
		t.Error("Expected error for non-existent game")
	}
}
