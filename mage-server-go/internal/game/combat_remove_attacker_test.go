package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap/zaptest"
)

// TestCombatRemoveAttacker verifies removing an attacker from combat
func TestCombatRemoveAttacker(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-remove-attacker"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create attacker
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

	// Declare attacker
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	if err := engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice"); err != nil {
		t.Fatalf("Failed to declare attacker: %v", err)
	}

	// Verify attacker is attacking and tapped
	gameState.mu.RLock()
	attacker := gameState.cards[attackerID]
	isAttacking := gameState.combat.attackers[attackerID]
	wasTapped := gameState.combat.attackersTapped[attackerID]
	gameState.mu.RUnlock()

	if !attacker.Attacking {
		t.Error("Attacker should be attacking before removal")
	}
	if !isAttacking {
		t.Error("Attacker should be in attackers set before removal")
	}
	if !attacker.Tapped {
		t.Error("Attacker should be tapped before removal")
	}
	if !wasTapped {
		t.Error("Attacker should be in attackersTapped set before removal")
	}

	// Remove attacker
	if err := engine.RemoveAttacker(gameID, attackerID); err != nil {
		t.Fatalf("Failed to remove attacker: %v", err)
	}

	// Verify attacker is no longer attacking and is untapped
	gameState.mu.RLock()
	attacker = gameState.cards[attackerID]
	isAttacking = gameState.combat.attackers[attackerID]
	wasTapped = gameState.combat.attackersTapped[attackerID]
	groupCount := len(gameState.combat.groups)
	gameState.mu.RUnlock()

	if attacker.Attacking {
		t.Error("Attacker should not be attacking after removal")
	}
	if isAttacking {
		t.Error("Attacker should not be in attackers set after removal")
	}
	if attacker.Tapped {
		t.Error("Attacker should be untapped after removal (was tapped by attack)")
	}
	if wasTapped {
		t.Error("Attacker should not be in attackersTapped set after removal")
	}
	if groupCount != 0 {
		t.Errorf("Expected 0 combat groups after removing only attacker, got %d", groupCount)
	}

	engine.EndCombat(gameID)
}

// TestCombatRemoveAttackerWithVigilance verifies removing attacker with vigilance
func TestCombatRemoveAttackerWithVigilance(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-remove-attacker-vigilance"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create attacker with vigilance
	gameState.mu.Lock()
	attackerID := "serra-angel"

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
		Abilities: []EngineAbilityView{
			{ID: abilityVigilance, Text: "Vigilance"},
		},
	}
	gameState.mu.Unlock()

	// Declare attacker
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	if err := engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice"); err != nil {
		t.Fatalf("Failed to declare attacker: %v", err)
	}

	// Verify attacker is attacking and NOT tapped (vigilance)
	gameState.mu.RLock()
	attacker := gameState.cards[attackerID]
	wasTapped := gameState.combat.attackersTapped[attackerID]
	gameState.mu.RUnlock()

	if attacker.Tapped {
		t.Error("Attacker with vigilance should not be tapped")
	}
	if wasTapped {
		t.Error("Attacker with vigilance should not be in attackersTapped set")
	}

	// Remove attacker
	if err := engine.RemoveAttacker(gameID, attackerID); err != nil {
		t.Fatalf("Failed to remove attacker: %v", err)
	}

	// Verify attacker is still untapped (was never tapped)
	gameState.mu.RLock()
	attacker = gameState.cards[attackerID]
	gameState.mu.RUnlock()

	if attacker.Tapped {
		t.Error("Attacker should still be untapped after removal")
	}

	engine.EndCombat(gameID)
}

// TestCombatRemoveAttackerMultiple verifies removing one of multiple attackers
func TestCombatRemoveAttackerMultiple(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-remove-attacker-multiple"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create two attackers
	gameState.mu.Lock()
	attacker1ID := "grizzly-bears"
	attacker2ID := "serra-angel"

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
		Name:         "Serra Angel",
		Type:         "Creature - Angel",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "4",
		Toughness:    "4",
		Tapped:       false,
	}
	gameState.mu.Unlock()

	// Declare both attackers
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	if err := engine.DeclareAttacker(gameID, attacker1ID, "Bob", "Alice"); err != nil {
		t.Fatalf("Failed to declare attacker 1: %v", err)
	}
	if err := engine.DeclareAttacker(gameID, attacker2ID, "Bob", "Alice"); err != nil {
		t.Fatalf("Failed to declare attacker 2: %v", err)
	}

	// Verify both attacking
	gameState.mu.RLock()
	attackerCount := len(gameState.combat.attackers)
	gameState.mu.RUnlock()

	if attackerCount != 2 {
		t.Errorf("Expected 2 attackers, got %d", attackerCount)
	}

	// Remove first attacker
	if err := engine.RemoveAttacker(gameID, attacker1ID); err != nil {
		t.Fatalf("Failed to remove attacker: %v", err)
	}

	// Verify only second attacker remains
	gameState.mu.RLock()
	attacker1 := gameState.cards[attacker1ID]
	attacker2 := gameState.cards[attacker2ID]
	isAttacking1 := gameState.combat.attackers[attacker1ID]
	isAttacking2 := gameState.combat.attackers[attacker2ID]
	groupCount := len(gameState.combat.groups)
	gameState.mu.RUnlock()

	if attacker1.Attacking {
		t.Error("First attacker should not be attacking after removal")
	}
	if !attacker2.Attacking {
		t.Error("Second attacker should still be attacking")
	}
	if isAttacking1 {
		t.Error("First attacker should not be in attackers set")
	}
	if !isAttacking2 {
		t.Error("Second attacker should still be in attackers set")
	}
	if groupCount != 1 {
		t.Errorf("Expected 1 combat group remaining, got %d", groupCount)
	}

	engine.EndCombat(gameID)
}

// TestCombatRemoveAttackerEvent verifies EventRemovedFromCombat fires
func TestCombatRemoveAttackerEvent(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-remove-attacker-event"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create attacker
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

	// Subscribe to event
	eventFired := false
	gameState.eventBus.SubscribeTyped(rules.EventRemovedFromCombat, func(e rules.Event) {
		if e.TargetID == attackerID {
			eventFired = true
		}
	})

	// Declare and remove attacker
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.RemoveAttacker(gameID, attackerID)

	if !eventFired {
		t.Error("Expected EventRemovedFromCombat to fire when attacker removed")
	}

	engine.EndCombat(gameID)
}

// TestCombatRemoveAttackerNotAttacking verifies error when removing non-attacker
func TestCombatRemoveAttackerNotAttacking(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-remove-attacker-not-attacking"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create creature but don't attack
	gameState.mu.Lock()
	creatureID := "grizzly-bears"

	gameState.cards[creatureID] = &internalCard{
		ID:           creatureID,
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

	// Try to remove non-attacker
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	err := engine.RemoveAttacker(gameID, creatureID)

	// Should fail
	if err == nil {
		t.Error("Expected error when removing non-attacking creature, got nil")
	}

	engine.EndCombat(gameID)
}
