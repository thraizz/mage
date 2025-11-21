package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestCanAttackBasic verifies basic CanAttack functionality
func TestCanAttackBasic(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-can-attack-basic"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create untapped creature without summoning sickness
	gameState.mu.Lock()
	creatureID := "grizzly-bears"

	gameState.cards[creatureID] = &internalCard{
		ID:                creatureID,
		Name:              "Grizzly Bears",
		Type:              "Creature - Bear",
		Zone:              zoneBattlefield,
		OwnerID:           "Alice",
		ControllerID:      "Alice",
		Power:             "2",
		Toughness:         "2",
		Tapped:            false,
		SummoningSickness: false,
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	// Test: Can attack
	canAttack, err := engine.CanAttack(gameID, creatureID)
	if err != nil {
		t.Fatalf("CanAttack returned error: %v", err)
	}

	if !canAttack {
		t.Error("Untapped creature without summoning sickness should be able to attack")
	}

	engine.EndCombat(gameID)
}

// TestCanAttackTapped verifies tapped creature cannot attack
func TestCanAttackTapped(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-can-attack-tapped"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create tapped creature
	gameState.mu.Lock()
	creatureID := "grizzly-bears"

	gameState.cards[creatureID] = &internalCard{
		ID:                creatureID,
		Name:              "Grizzly Bears",
		Type:              "Creature - Bear",
		Zone:              zoneBattlefield,
		OwnerID:           "Alice",
		ControllerID:      "Alice",
		Power:             "2",
		Toughness:         "2",
		Tapped:            true,
		SummoningSickness: false,
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	// Test: Cannot attack
	canAttack, err := engine.CanAttack(gameID, creatureID)
	if err != nil {
		t.Fatalf("CanAttack returned error: %v", err)
	}

	if canAttack {
		t.Error("Tapped creature should not be able to attack")
	}

	engine.EndCombat(gameID)
}

// TestCanAttackSummoningSickness verifies creature with summoning sickness cannot attack
func TestCanAttackSummoningSickness(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-can-attack-summoning-sickness"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create creature with summoning sickness
	gameState.mu.Lock()
	creatureID := "grizzly-bears"

	gameState.cards[creatureID] = &internalCard{
		ID:                creatureID,
		Name:              "Grizzly Bears",
		Type:              "Creature - Bear",
		Zone:              zoneBattlefield,
		OwnerID:           "Alice",
		ControllerID:      "Alice",
		Power:             "2",
		Toughness:         "2",
		Tapped:            false,
		SummoningSickness: true,
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	// Test: Cannot attack
	canAttack, err := engine.CanAttack(gameID, creatureID)
	if err != nil {
		t.Fatalf("CanAttack returned error: %v", err)
	}

	if canAttack {
		t.Error("Creature with summoning sickness should not be able to attack")
	}

	engine.EndCombat(gameID)
}

// TestCanAttackDefender verifies creature with defender cannot attack
func TestCanAttackDefenderAbility(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-can-attack-defender"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create creature with defender
	gameState.mu.Lock()
	creatureID := "wall-of-wood"

	gameState.cards[creatureID] = &internalCard{
		ID:                creatureID,
		Name:              "Wall of Wood",
		Type:              "Creature - Wall",
		Zone:              zoneBattlefield,
		OwnerID:           "Alice",
		ControllerID:      "Alice",
		Power:             "0",
		Toughness:         "3",
		Tapped:            false,
		SummoningSickness: false,
		Abilities: []EngineAbilityView{
			{ID: abilityDefender, Text: "Defender"},
		},
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	// Test: Cannot attack
	canAttack, err := engine.CanAttack(gameID, creatureID)
	if err != nil {
		t.Fatalf("CanAttack returned error: %v", err)
	}

	if canAttack {
		t.Error("Creature with defender should not be able to attack")
	}

	engine.EndCombat(gameID)
}

// TestCanAttackDefenderSpecific verifies CanAttackDefender with specific defender
func TestCanAttackDefenderSpecific(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-can-attack-defender-specific"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create creature
	gameState.mu.Lock()
	creatureID := "grizzly-bears"

	gameState.cards[creatureID] = &internalCard{
		ID:                creatureID,
		Name:              "Grizzly Bears",
		Type:              "Creature - Bear",
		Zone:              zoneBattlefield,
		OwnerID:           "Alice",
		ControllerID:      "Alice",
		Power:             "2",
		Toughness:         "2",
		Tapped:            false,
		SummoningSickness: false,
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	// Test: Can attack Bob
	canAttack, err := engine.CanAttackDefender(gameID, creatureID, "Bob")
	if err != nil {
		t.Fatalf("CanAttackDefender returned error: %v", err)
	}

	if !canAttack {
		t.Error("Creature should be able to attack Bob")
	}

	engine.EndCombat(gameID)
}

// TestCanAttackDefenderInvalid verifies error for invalid defender
func TestCanAttackDefenderInvalid(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-can-attack-defender-invalid"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create creature
	gameState.mu.Lock()
	creatureID := "grizzly-bears"

	gameState.cards[creatureID] = &internalCard{
		ID:                creatureID,
		Name:              "Grizzly Bears",
		Type:              "Creature - Bear",
		Zone:              zoneBattlefield,
		OwnerID:           "Alice",
		ControllerID:      "Alice",
		Power:             "2",
		Toughness:         "2",
		Tapped:            false,
		SummoningSickness: false,
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	// Test: Cannot attack invalid defender
	canAttack, err := engine.CanAttackDefender(gameID, creatureID, "InvalidPlayer")

	// Should return error
	if err == nil {
		t.Error("Expected error for invalid defender, got nil")
	}

	if canAttack {
		t.Error("Should not be able to attack invalid defender")
	}

	engine.EndCombat(gameID)
}

// TestCanAttackDefenderTapped verifies tapped creature cannot attack specific defender
func TestCanAttackDefenderTapped(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-can-attack-defender-tapped"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create tapped creature
	gameState.mu.Lock()
	creatureID := "grizzly-bears"

	gameState.cards[creatureID] = &internalCard{
		ID:                creatureID,
		Name:              "Grizzly Bears",
		Type:              "Creature - Bear",
		Zone:              zoneBattlefield,
		OwnerID:           "Alice",
		ControllerID:      "Alice",
		Power:             "2",
		Toughness:         "2",
		Tapped:            true,
		SummoningSickness: false,
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	// Test: Cannot attack Bob
	canAttack, err := engine.CanAttackDefender(gameID, creatureID, "Bob")
	if err != nil {
		t.Fatalf("CanAttackDefender returned error: %v", err)
	}

	if canAttack {
		t.Error("Tapped creature should not be able to attack Bob")
	}

	engine.EndCombat(gameID)
}
