package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestOrderBlockersBasic verifies basic blocker ordering
func TestOrderBlockersBasic(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-order-blockers-basic"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create attacker and three blockers
	gameState.mu.Lock()
	attackerID := "attacker"
	blocker1ID := "blocker1"
	blocker2ID := "blocker2"
	blocker3ID := "blocker3"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Attacker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "5",
		Toughness:    "5",
		Tapped:       false,
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
		Tapped:       false,
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
		Tapped:       false,
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
	engine.DeclareBlocker(gameID, blocker3ID, attackerID, "Bob")
	
	// Get initial blocker order
	gameState.mu.RLock()
	var initialOrder []string
	for _, group := range gameState.combat.groups {
		if len(group.attackers) > 0 && group.attackers[0] == attackerID {
			initialOrder = append([]string{}, group.blockers...)
			break
		}
	}
	gameState.mu.RUnlock()
	
	// Verify initial order (should be declaration order)
	if len(initialOrder) != 3 {
		t.Fatalf("Expected 3 blockers, got %d", len(initialOrder))
	}
	if initialOrder[0] != blocker1ID || initialOrder[1] != blocker2ID || initialOrder[2] != blocker3ID {
		t.Errorf("Initial order incorrect: got %v, expected [%s, %s, %s]", 
			initialOrder, blocker1ID, blocker2ID, blocker3ID)
	}
	
	// Reorder blockers (reverse order)
	newOrder := []string{blocker3ID, blocker2ID, blocker1ID}
	if err := engine.OrderBlockers(gameID, attackerID, newOrder); err != nil {
		t.Fatalf("Failed to order blockers: %v", err)
	}
	
	// Verify new order
	gameState.mu.RLock()
	var finalOrder []string
	for _, group := range gameState.combat.groups {
		if len(group.attackers) > 0 && group.attackers[0] == attackerID {
			finalOrder = append([]string{}, group.blockers...)
			break
		}
	}
	gameState.mu.RUnlock()
	
	if len(finalOrder) != 3 {
		t.Fatalf("Expected 3 blockers after reorder, got %d", len(finalOrder))
	}
	if finalOrder[0] != blocker3ID || finalOrder[1] != blocker2ID || finalOrder[2] != blocker1ID {
		t.Errorf("Final order incorrect: got %v, expected %v", finalOrder, newOrder)
	}
	
	engine.EndCombat(gameID)
}

// TestOrderBlockersInvalidCount verifies error when blocker count doesn't match
func TestOrderBlockersInvalidCount(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-order-blockers-invalid-count"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create attacker and two blockers
	gameState.mu.Lock()
	attackerID := "attacker"
	blocker1ID := "blocker1"
	blocker2ID := "blocker2"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Attacker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "5",
		Toughness:    "5",
		Tapped:       false,
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
		Tapped:       false,
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
	
	// Try to order with wrong count (only 1 blocker instead of 2)
	wrongOrder := []string{blocker1ID}
	err := engine.OrderBlockers(gameID, attackerID, wrongOrder)
	
	if err == nil {
		t.Error("Expected error when ordering with wrong blocker count, got nil")
	}
	
	engine.EndCombat(gameID)
}

// TestOrderBlockersInvalidBlocker verifies error when blocker not in group
func TestOrderBlockersInvalidBlocker(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-order-blockers-invalid-blocker"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create attacker and two blockers
	gameState.mu.Lock()
	attackerID := "attacker"
	blocker1ID := "blocker1"
	blocker2ID := "blocker2"
	invalidBlockerID := "invalid-blocker"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Attacker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "5",
		Toughness:    "5",
		Tapped:       false,
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
		Tapped:       false,
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
	
	// Try to order with invalid blocker
	wrongOrder := []string{blocker1ID, invalidBlockerID}
	err := engine.OrderBlockers(gameID, attackerID, wrongOrder)
	
	if err == nil {
		t.Error("Expected error when ordering with invalid blocker, got nil")
	}
	
	engine.EndCombat(gameID)
}

// TestOrderBlockersDamageAssignment verifies damage is assigned in specified order
func TestOrderBlockersDamageAssignment(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-order-blockers-damage"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create 4/4 attacker with trample and three 1/1 blockers
	gameState.mu.Lock()
	attackerID := "attacker"
	blocker1ID := "blocker1"
	blocker2ID := "blocker2"
	blocker3ID := "blocker3"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Attacker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "4",
		Toughness:    "4",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityTrample, Text: "Trample"},
		},
	}
	
	gameState.cards[blocker1ID] = &internalCard{
		ID:           blocker1ID,
		Name:         "Blocker 1",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "1",
		Toughness:    "1",
		Tapped:       false,
	}
	
	gameState.cards[blocker2ID] = &internalCard{
		ID:           blocker2ID,
		Name:         "Blocker 2",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "1",
		Toughness:    "1",
		Tapped:       false,
	}
	
	gameState.cards[blocker3ID] = &internalCard{
		ID:           blocker3ID,
		Name:         "Blocker 3",
		Type:         "Creature",
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
	engine.DeclareBlocker(gameID, blocker3ID, attackerID, "Bob")
	
	// Reorder blockers: blocker3, blocker1, blocker2
	newOrder := []string{blocker3ID, blocker1ID, blocker2ID}
	if err := engine.OrderBlockers(gameID, attackerID, newOrder); err != nil {
		t.Fatalf("Failed to order blockers: %v", err)
	}
	
	engine.AcceptBlockers(gameID)
	
	// Assign and apply damage
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	
	// Verify damage was assigned in order: blocker3, blocker1, blocker2
	// With 4 power and trample, should assign 1 to each blocker (lethal), 1 tramples through
	gameState.mu.RLock()
	blocker1 := gameState.cards[blocker1ID]
	blocker2 := gameState.cards[blocker2ID]
	blocker3 := gameState.cards[blocker3ID]
	bobPlayer := gameState.players["Bob"]
	gameState.mu.RUnlock()
	
	// Each blocker should have 1 damage
	if blocker1.Damage != 1 {
		t.Errorf("Blocker 1 should have 1 damage, got %d", blocker1.Damage)
	}
	if blocker2.Damage != 1 {
		t.Errorf("Blocker 2 should have 1 damage, got %d", blocker2.Damage)
	}
	if blocker3.Damage != 1 {
		t.Errorf("Blocker 3 should have 1 damage, got %d", blocker3.Damage)
	}
	
	// Bob should have taken 1 trample damage
	if bobPlayer.Life != 19 {
		t.Errorf("Bob should have 19 life (1 trample damage), got %d", bobPlayer.Life)
	}
	
	engine.EndCombat(gameID)
}
