package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap/zaptest"
)

// TestCombatEventsBeginCombat tests that EventBeginCombatStep is fired
func TestCombatEventsBeginCombat(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-begin-combat"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Subscribe to begin combat event
	eventFired := false
	gameState.mu.Lock()
	gameState.eventBus.SubscribeTyped(rules.EventBeginCombatStep, func(event rules.Event) {
		eventFired = true
	})
	gameState.mu.Unlock()
	
	// Reset combat (should fire event)
	engine.ResetCombat(gameID)
	
	if !eventFired {
		t.Error("EventBeginCombatStep should have been fired")
	}
}

// TestCombatEventsDeclareAttackers tests attacker declaration events
func TestCombatEventsDeclareAttackers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-declare-attackers"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Attacker
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
	
	// Subscribe to events
	preEventFired := false
	attackerDeclaredFired := false
	defenderAttackedFired := false
	declaredAttackersFired := false
	
	gameState.mu.Lock()
	gameState.eventBus.SubscribeTyped(rules.EventDeclareAttackersStepPre, func(event rules.Event) {
		preEventFired = true
	})
	gameState.eventBus.SubscribeTyped(rules.EventAttackerDeclared, func(event rules.Event) {
		attackerDeclaredFired = true
		if event.SourceID != attackerID {
			t.Errorf("expected attacker ID %s, got %s", attackerID, event.SourceID)
		}
	})
	gameState.eventBus.SubscribeTyped(rules.EventDefenderAttacked, func(event rules.Event) {
		defenderAttackedFired = true
		if event.TargetID != "Bob" {
			t.Errorf("expected defender ID Bob, got %s", event.TargetID)
		}
	})
	gameState.eventBus.SubscribeTyped(rules.EventDeclaredAttackers, func(event rules.Event) {
		declaredAttackersFired = true
	})
	gameState.mu.Unlock()
	
	// Combat flow
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.FinishDeclaringAttackers(gameID)
	
	// Verify events
	if !preEventFired {
		t.Error("EventDeclareAttackersStepPre should have been fired")
	}
	if !attackerDeclaredFired {
		t.Error("EventAttackerDeclared should have been fired")
	}
	if !defenderAttackedFired {
		t.Error("EventDefenderAttacked should have been fired")
	}
	if !declaredAttackersFired {
		t.Error("EventDeclaredAttackers should have been fired")
	}
}

// TestCombatEventsDeclareBlockers tests blocker declaration events
func TestCombatEventsDeclareBlockers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-declare-blockers"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Attacker and blocker
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
		Name:         "Runeclaw Bear",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Subscribe to events
	preEventFired := false
	blockerDeclaredFired := false
	declaredBlockersFired := false
	
	gameState.mu.Lock()
	gameState.eventBus.SubscribeTyped(rules.EventDeclareBlockersStepPre, func(event rules.Event) {
		preEventFired = true
	})
	gameState.eventBus.SubscribeTyped(rules.EventBlockerDeclared, func(event rules.Event) {
		blockerDeclaredFired = true
		if event.SourceID != blockerID {
			t.Errorf("expected blocker ID %s, got %s", blockerID, event.SourceID)
		}
	})
	gameState.eventBus.SubscribeTyped(rules.EventDeclaredBlockers, func(event rules.Event) {
		declaredBlockersFired = true
	})
	gameState.mu.Unlock()
	
	// Combat flow
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)
	
	// Verify events
	if !preEventFired {
		t.Error("EventDeclareBlockersStepPre should have been fired")
	}
	if !blockerDeclaredFired {
		t.Error("EventBlockerDeclared should have been fired")
	}
	if !declaredBlockersFired {
		t.Error("EventDeclaredBlockers should have been fired")
	}
}

// TestCombatEventsCombatDamage tests combat damage events
func TestCombatEventsCombatDamage(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-combat-damage"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Attacker
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
	
	// Subscribe to events
	preEventFired := false
	appliedEventFired := false
	
	gameState.mu.Lock()
	gameState.eventBus.SubscribeTyped(rules.EventCombatDamageStepPre, func(event rules.Event) {
		preEventFired = true
	})
	gameState.eventBus.SubscribeTyped(rules.EventCombatDamageApplied, func(event rules.Event) {
		appliedEventFired = true
	})
	gameState.mu.Unlock()
	
	// Combat flow
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.AcceptBlockers(gameID)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	
	// Verify events
	if !preEventFired {
		t.Error("EventCombatDamageStepPre should have been fired")
	}
	if !appliedEventFired {
		t.Error("EventCombatDamageApplied should have been fired")
	}
}

// TestCombatEventsEndCombat tests end combat events
func TestCombatEventsEndCombat(t *testing.T) {
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
	
	// Setup: Attacker
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
	
	// Subscribe to events
	preEventFired := false
	endEventFired := false
	
	gameState.mu.Lock()
	gameState.eventBus.SubscribeTyped(rules.EventEndCombatStepPre, func(event rules.Event) {
		preEventFired = true
	})
	gameState.eventBus.SubscribeTyped(rules.EventEndCombatStep, func(event rules.Event) {
		endEventFired = true
	})
	gameState.mu.Unlock()
	
	// Combat flow
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.AcceptBlockers(gameID)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	engine.EndCombat(gameID)
	
	// Verify events
	if !preEventFired {
		t.Error("EventEndCombatStepPre should have been fired")
	}
	if !endEventFired {
		t.Error("EventEndCombatStep should have been fired")
	}
}

// TestCombatEventsFullFlow tests all combat events in a complete flow
func TestCombatEventsFullFlow(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-full-flow"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Attacker and blocker
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
		Name:         "Runeclaw Bear",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Track all events
	eventsFired := make(map[rules.EventType]bool)
	
	gameState.mu.Lock()
	eventTypes := []rules.EventType{
		rules.EventBeginCombatStep,
		rules.EventDeclareAttackersStepPre,
		rules.EventAttackerDeclared,
		rules.EventDefenderAttacked,
		rules.EventDeclaredAttackers,
		rules.EventDeclareBlockersStepPre,
		rules.EventBlockerDeclared,
		rules.EventDeclaredBlockers,
		rules.EventCombatDamageStepPre,
		rules.EventCombatDamageApplied,
		rules.EventEndCombatStepPre,
		rules.EventEndCombatStep,
	}
	
	for _, eventType := range eventTypes {
		et := eventType // Capture for closure
		gameState.eventBus.SubscribeTyped(et, func(event rules.Event) {
			eventsFired[et] = true
		})
	}
	gameState.mu.Unlock()
	
	// Full combat flow
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.FinishDeclaringAttackers(gameID)
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	engine.EndCombat(gameID)
	
	// Verify all events fired
	for _, eventType := range eventTypes {
		if !eventsFired[eventType] {
			t.Errorf("Event %s should have been fired", eventType)
		}
	}
}
