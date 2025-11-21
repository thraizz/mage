package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap/zaptest"
)

// TestCombatDamageAssignedEvent verifies EventCombatDamageAssigned fires
func TestCombatDamageAssignedEvent(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-damage-assigned-event"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create attacker and blocker
	gameState.mu.Lock()
	attackerID := "grizzly-bears"
	blockerID := "wall-of-wood"

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
		Name:         "Wall of Wood",
		Type:         "Creature - Wall",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "0",
		Toughness:    "3",
		Tapped:       false,
	}

	// Subscribe to event
	assignedEventFired := false
	gameState.eventBus.SubscribeTyped(rules.EventCombatDamageAssigned, func(e rules.Event) {
		assignedEventFired = true
	})
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)

	// Assign damage
	if err := engine.AssignCombatDamage(gameID, false); err != nil {
		t.Fatalf("Failed to assign combat damage: %v", err)
	}

	// Verify event fired
	if !assignedEventFired {
		t.Error("Expected EventCombatDamageAssigned to fire")
	}

	engine.EndCombat(gameID)
}

// TestCombatDamageAppliedEvent verifies EventCombatDamageApplied fires
func TestCombatDamageAppliedEvent(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-damage-applied-event"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create attacker and blocker
	gameState.mu.Lock()
	attackerID := "grizzly-bears"
	blockerID := "wall-of-wood"

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
		Name:         "Wall of Wood",
		Type:         "Creature - Wall",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "0",
		Toughness:    "3",
		Tapped:       false,
	}

	// Subscribe to event
	appliedEventFired := false
	gameState.eventBus.SubscribeTyped(rules.EventCombatDamageApplied, func(e rules.Event) {
		appliedEventFired = true
	})
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
	if err := engine.ApplyCombatDamage(gameID); err != nil {
		t.Fatalf("Failed to apply combat damage: %v", err)
	}

	// Verify event fired
	if !appliedEventFired {
		t.Error("Expected EventCombatDamageApplied to fire")
	}

	engine.EndCombat(gameID)
}

// TestCombatDamageBothEvents verifies both events fire in sequence
func TestCombatDamageBothEvents(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-both-damage-events"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create attacker and blocker
	gameState.mu.Lock()
	attackerID := "grizzly-bears"
	blockerID := "wall-of-wood"

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
		Name:         "Wall of Wood",
		Type:         "Creature - Wall",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "0",
		Toughness:    "3",
		Tapped:       false,
	}

	// Track event order
	eventOrder := []string{}

	gameState.eventBus.SubscribeTyped(rules.EventCombatDamageAssigned, func(e rules.Event) {
		eventOrder = append(eventOrder, "assigned")
	})

	gameState.eventBus.SubscribeTyped(rules.EventCombatDamageApplied, func(e rules.Event) {
		eventOrder = append(eventOrder, "applied")
	})
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

	// Verify both events fired in correct order
	if len(eventOrder) != 2 {
		t.Errorf("Expected 2 events, got %d", len(eventOrder))
	}

	if len(eventOrder) >= 1 && eventOrder[0] != "assigned" {
		t.Errorf("Expected first event to be 'assigned', got '%s'", eventOrder[0])
	}

	if len(eventOrder) >= 2 && eventOrder[1] != "applied" {
		t.Errorf("Expected second event to be 'applied', got '%s'", eventOrder[1])
	}

	engine.EndCombat(gameID)
}

// TestCombatDamageFirstStrikeEvents verifies events fire for first strike damage
func TestCombatDamageFirstStrikeEvents(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-first-strike-damage-events"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create first strike attacker and blocker
	gameState.mu.Lock()
	attackerID := "first-strike-bear"
	blockerID := "wall-of-wood"

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
			{ID: abilityFirstStrike, Text: "First Strike"},
		},
	}

	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Wall of Wood",
		Type:         "Creature - Wall",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "0",
		Toughness:    "3",
		Tapped:       false,
	}

	// Track events
	assignedCount := 0
	appliedCount := 0

	gameState.eventBus.SubscribeTyped(rules.EventCombatDamageAssigned, func(e rules.Event) {
		assignedCount++
	})

	gameState.eventBus.SubscribeTyped(rules.EventCombatDamageApplied, func(e rules.Event) {
		appliedCount++
	})
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)

	// First strike damage
	engine.AssignCombatDamage(gameID, true)
	engine.ApplyCombatDamage(gameID)

	// Normal damage (should still fire even if no creatures deal damage)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify events fired twice (once for first strike, once for normal)
	if assignedCount != 2 {
		t.Errorf("Expected EventCombatDamageAssigned to fire 2 times, got %d", assignedCount)
	}

	if appliedCount != 2 {
		t.Errorf("Expected EventCombatDamageApplied to fire 2 times, got %d", appliedCount)
	}

	engine.EndCombat(gameID)
}
