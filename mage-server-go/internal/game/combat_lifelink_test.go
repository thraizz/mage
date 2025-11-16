package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestLifelinkBasic verifies basic lifelink functionality
func TestLifelinkBasic(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-lifelink"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create attacker with lifelink
	gameState.mu.Lock()
	attackerID := "attacker"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Lifelink Attacker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "3",
		Toughness:    "3",
		Abilities: []EngineAbilityView{
			{ID: abilityLifelink, Text: "Lifelink"},
		},
	}
	
	// Record initial life
	initialLife := gameState.players["Alice"].Life
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.AcceptBlockers(gameID)
	
	// Assign and apply damage (unblocked attacker deals 3 to Bob)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	
	// Check that Alice gained 3 life from lifelink
	gameState.mu.RLock()
	aliceLife := gameState.players["Alice"].Life
	bobLife := gameState.players["Bob"].Life
	gameState.mu.RUnlock()
	
	expectedLife := initialLife + 3
	if aliceLife != expectedLife {
		t.Errorf("Expected Alice to have %d life (gained 3 from lifelink), got %d", expectedLife, aliceLife)
	}
	
	expectedBobLife := 20 - 3 // Starting life minus damage
	if bobLife != expectedBobLife {
		t.Errorf("Expected Bob to have %d life, got %d", expectedBobLife, bobLife)
	}
	
	engine.EndCombat(gameID)
}

// TestLifelinkBlocked verifies lifelink when attacker is blocked
func TestLifelinkBlocked(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-lifelink-blocked"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create attacker with lifelink and blocker
	gameState.mu.Lock()
	attackerID := "attacker"
	blockerID := "blocker"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Lifelink Attacker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "4",
		Toughness:    "4",
		Abilities: []EngineAbilityView{
			{ID: abilityLifelink, Text: "Lifelink"},
		},
	}
	
	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Blocker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
	}
	
	// Record initial life
	initialLife := gameState.players["Alice"].Life
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)
	
	// Assign and apply damage
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	
	// Check that Alice gained 4 life from lifelink (damage to blocker)
	gameState.mu.RLock()
	aliceLife := gameState.players["Alice"].Life
	bobLife := gameState.players["Bob"].Life
	gameState.mu.RUnlock()
	
	expectedLife := initialLife + 4
	if aliceLife != expectedLife {
		t.Errorf("Expected Alice to have %d life (gained 4 from lifelink), got %d", expectedLife, aliceLife)
	}
	
	// Bob should not have lost life (damage went to blocker)
	if bobLife != 20 {
		t.Errorf("Expected Bob to have 20 life, got %d", bobLife)
	}
	
	engine.EndCombat(gameID)
}

// TestLifelinkBlocker verifies lifelink on a blocking creature
func TestLifelinkBlocker(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-lifelink-blocker"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create attacker and blocker with lifelink
	gameState.mu.Lock()
	attackerID := "attacker"
	blockerID := "blocker"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Attacker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "3",
		Toughness:    "3",
	}
	
	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Lifelink Blocker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "5",
		Abilities: []EngineAbilityView{
			{ID: abilityLifelink, Text: "Lifelink"},
		},
	}
	
	// Record initial life
	initialBobLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)
	
	// Assign and apply damage
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	
	// Check that Bob gained 2 life from lifelink (blocker dealt 2 damage to attacker)
	gameState.mu.RLock()
	bobLife := gameState.players["Bob"].Life
	gameState.mu.RUnlock()
	
	expectedLife := initialBobLife + 2
	if bobLife != expectedLife {
		t.Errorf("Expected Bob to have %d life (gained 2 from lifelink blocker), got %d", expectedLife, bobLife)
	}
	
	engine.EndCombat(gameID)
}

// TestLifelinkTrample verifies lifelink with trample
func TestLifelinkTrample(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-lifelink-trample"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create attacker with lifelink and trample, and a small blocker
	gameState.mu.Lock()
	attackerID := "attacker"
	blockerID := "blocker"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Lifelink Trampler",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "5",
		Toughness:    "5",
		Abilities: []EngineAbilityView{
			{ID: abilityLifelink, Text: "Lifelink"},
			{ID: abilityTrample, Text: "Trample"},
		},
	}
	
	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Small Blocker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "1",
		Toughness:    "2",
	}
	
	// Record initial life
	initialLife := gameState.players["Alice"].Life
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)
	
	// Assign and apply damage
	// Attacker should deal 2 to blocker (lethal) and 3 trample to Bob
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	
	// Check that Alice gained 5 life from lifelink (2 to blocker + 3 to Bob)
	gameState.mu.RLock()
	aliceLife := gameState.players["Alice"].Life
	bobLife := gameState.players["Bob"].Life
	gameState.mu.RUnlock()
	
	expectedLife := initialLife + 5
	if aliceLife != expectedLife {
		t.Errorf("Expected Alice to have %d life (gained 5 from lifelink with trample), got %d", expectedLife, aliceLife)
	}
	
	expectedBobLife := 20 - 3 // Starting life minus trample damage
	if bobLife != expectedBobLife {
		t.Errorf("Expected Bob to have %d life, got %d", expectedBobLife, bobLife)
	}
	
	engine.EndCombat(gameID)
}

// TestLifelinkMultipleBlockers verifies lifelink with multiple blockers
func TestLifelinkMultipleBlockers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-lifelink-multiple"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create attacker with lifelink and two blockers
	gameState.mu.Lock()
	attackerID := "attacker"
	blocker1ID := "blocker1"
	blocker2ID := "blocker2"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Lifelink Attacker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "6",
		Toughness:    "6",
		Abilities: []EngineAbilityView{
			{ID: abilityLifelink, Text: "Lifelink"},
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
	
	// Record initial life
	initialLife := gameState.players["Alice"].Life
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blocker1ID, attackerID, "Bob")
	engine.DeclareBlocker(gameID, blocker2ID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)
	
	// Assign and apply damage
	// Attacker deals 6 damage split between two blockers (3 each)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	
	// Check that Alice gained 6 life from lifelink (total damage dealt)
	gameState.mu.RLock()
	aliceLife := gameState.players["Alice"].Life
	gameState.mu.RUnlock()
	
	expectedLife := initialLife + 6
	if aliceLife != expectedLife {
		t.Errorf("Expected Alice to have %d life (gained 6 from lifelink), got %d", expectedLife, aliceLife)
	}
	
	engine.EndCombat(gameID)
}
