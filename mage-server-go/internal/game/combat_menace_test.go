package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestMenaceBlockedByOne verifies menace creature can't be blocked by a single creature
func TestMenaceBlockedByOne(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-menace-blocked-by-one"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create menace attacker and single blocker
	gameState.mu.Lock()
	attackerID := "attacker"
	blockerID := "blocker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Menace Creature",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Abilities: []EngineAbilityView{
			{ID: abilityMenace, Text: "Menace"},
		},
	}

	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Single Blocker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "3",
		Toughness:    "3",
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Try to block with single creature
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")

	// Accept blockers - should remove the single blocker due to menace
	engine.AcceptBlockers(gameID)

	// Verify blocker was removed
	gameState.mu.RLock()
	blocker, exists := gameState.cards[blockerID]
	if !exists {
		t.Fatal("Blocker card not found")
	}

	if blocker.Blocking {
		t.Error("Expected blocker to be removed due to menace")
	}

	// Verify combat group shows no blockers
	if len(gameState.combat.groups) != 1 {
		t.Fatalf("Expected 1 combat group, got %d", len(gameState.combat.groups))
	}

	group := gameState.combat.groups[0]
	if len(group.blockers) != 0 {
		t.Errorf("Expected 0 blockers in group, got %d", len(group.blockers))
	}

	if group.blocked {
		t.Error("Expected group to be unblocked")
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}

// TestMenaceBlockedByTwo verifies menace creature CAN be blocked by two creatures
func TestMenaceBlockedByTwo(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-menace-blocked-by-two"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create menace attacker and two blockers
	gameState.mu.Lock()
	attackerID := "attacker"
	blocker1ID := "blocker1"
	blocker2ID := "blocker2"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Menace Creature",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "3",
		Toughness:    "3",
		Abilities: []EngineAbilityView{
			{ID: abilityMenace, Text: "Menace"},
		},
	}

	gameState.cards[blocker1ID] = &internalCard{
		ID:           blocker1ID,
		Name:         "Blocker 1",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
	}

	gameState.cards[blocker2ID] = &internalCard{
		ID:           blocker2ID,
		Name:         "Blocker 2",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Block with two creatures
	engine.DeclareBlocker(gameID, blocker1ID, attackerID, "Bob")
	engine.DeclareBlocker(gameID, blocker2ID, attackerID, "Bob")

	// Accept blockers - should keep both blockers
	engine.AcceptBlockers(gameID)

	// Verify both blockers are still blocking
	gameState.mu.RLock()
	blocker1, exists := gameState.cards[blocker1ID]
	if !exists {
		t.Fatal("Blocker 1 card not found")
	}

	if !blocker1.Blocking {
		t.Error("Expected blocker 1 to still be blocking")
	}

	blocker2, exists := gameState.cards[blocker2ID]
	if !exists {
		t.Fatal("Blocker 2 card not found")
	}

	if !blocker2.Blocking {
		t.Error("Expected blocker 2 to still be blocking")
	}

	// Verify combat group shows two blockers
	if len(gameState.combat.groups) != 1 {
		t.Fatalf("Expected 1 combat group, got %d", len(gameState.combat.groups))
	}

	group := gameState.combat.groups[0]
	if len(group.blockers) != 2 {
		t.Errorf("Expected 2 blockers in group, got %d", len(group.blockers))
	}

	if !group.blocked {
		t.Error("Expected group to be blocked")
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}

// TestMenaceBlockedByThree verifies menace creature can be blocked by three or more creatures
func TestMenaceBlockedByThree(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-menace-blocked-by-three"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create menace attacker and three blockers
	gameState.mu.Lock()
	attackerID := "attacker"
	blocker1ID := "blocker1"
	blocker2ID := "blocker2"
	blocker3ID := "blocker3"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Menace Creature",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "4",
		Toughness:    "4",
		Abilities: []EngineAbilityView{
			{ID: abilityMenace, Text: "Menace"},
		},
	}

	gameState.cards[blocker1ID] = &internalCard{
		ID:           blocker1ID,
		Name:         "Blocker 1",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
	}

	gameState.cards[blocker2ID] = &internalCard{
		ID:           blocker2ID,
		Name:         "Blocker 2",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
	}

	gameState.cards[blocker3ID] = &internalCard{
		ID:           blocker3ID,
		Name:         "Blocker 3",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Block with three creatures
	engine.DeclareBlocker(gameID, blocker1ID, attackerID, "Bob")
	engine.DeclareBlocker(gameID, blocker2ID, attackerID, "Bob")
	engine.DeclareBlocker(gameID, blocker3ID, attackerID, "Bob")

	// Accept blockers - should keep all three blockers
	engine.AcceptBlockers(gameID)

	// Verify all three blockers are still blocking
	gameState.mu.RLock()
	for i, blockerID := range []string{blocker1ID, blocker2ID, blocker3ID} {
		blocker, exists := gameState.cards[blockerID]
		if !exists {
			t.Fatalf("Blocker %d card not found", i+1)
		}

		if !blocker.Blocking {
			t.Errorf("Expected blocker %d to still be blocking", i+1)
		}
	}

	// Verify combat group shows three blockers
	if len(gameState.combat.groups) != 1 {
		t.Fatalf("Expected 1 combat group, got %d", len(gameState.combat.groups))
	}

	group := gameState.combat.groups[0]
	if len(group.blockers) != 3 {
		t.Errorf("Expected 3 blockers in group, got %d", len(group.blockers))
	}

	if !group.blocked {
		t.Error("Expected group to be blocked")
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}

// TestMenaceUnblocked verifies menace creature can attack unblocked
func TestMenaceUnblocked(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-menace-unblocked"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create menace attacker with no blockers
	gameState.mu.Lock()
	attackerID := "attacker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Menace Creature",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "3",
		Toughness:    "3",
		Abilities: []EngineAbilityView{
			{ID: abilityMenace, Text: "Menace"},
		},
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// No blockers declared
	engine.AcceptBlockers(gameID)

	// Verify combat group shows no blockers
	gameState.mu.RLock()
	if len(gameState.combat.groups) != 1 {
		t.Fatalf("Expected 1 combat group, got %d", len(gameState.combat.groups))
	}

	group := gameState.combat.groups[0]
	if len(group.blockers) != 0 {
		t.Errorf("Expected 0 blockers in group, got %d", len(group.blockers))
	}

	if group.blocked {
		t.Error("Expected group to be unblocked")
	}
	gameState.mu.RUnlock()

	// Assign and apply damage - should deal damage to player
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify player took damage
	gameState.mu.RLock()
	bob, exists := gameState.players["Bob"]
	if !exists {
		t.Fatal("Bob not found")
	}

	expectedLife := 20 - 3 // Starting life minus attacker power
	if bob.Life != expectedLife {
		t.Errorf("Expected Bob to have %d life, got %d", expectedLife, bob.Life)
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}
