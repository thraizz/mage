package game

import (
	"fmt"
	"testing"

	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap/zaptest"
)

// TestCombatDamageTriggerPlayer verifies "Whenever ~ deals combat damage to a player" triggers
func TestCombatDamageTriggerPlayer(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-damage-player-trigger"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create attacker with damage trigger
	gameState.mu.Lock()
	attackerID := "attacker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Damage Trigger Creature",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "3",
		Toughness:    "3",
	}

	// Track if trigger fired and damage amount
	triggerFired := false
	var damageDealt int

	// Register combat damage trigger: "Whenever ~ deals combat damage to a player, draw a card"
	trigger := &combatTrigger{
		SourceID:    attackerID,
		TriggerType: "deals_damage_player",
		Condition: func(gs *engineGameState, event rules.Event) bool {
			// Check if this is a damaged player event from our creature
			return event.Type == rules.EventDamagedPlayer &&
				event.SourceID == attackerID &&
				event.Flag == true // Combat damage
		},
		CreateAbility: func(gs *engineGameState, event rules.Event) *triggeredAbilityQueueItem {
			return &triggeredAbilityQueueItem{
				ID:          fmt.Sprintf("%s-damage-player-trigger", attackerID),
				SourceID:    attackerID,
				Controller:  "Alice",
				Description: "Whenever Damage Trigger Creature deals combat damage to a player, draw a card",
				Resolve: func(gs *engineGameState) error {
					triggerFired = true
					damageDealt = event.Amount
					// In real implementation, would draw a card here
					return nil
				},
				UsesStack: true,
			}
		},
	}

	gameState.combatTriggers = append(gameState.combatTriggers, trigger)
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.AcceptBlockers(gameID)

	// Assign and apply damage - should trigger
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Check that trigger was queued
	gameState.mu.RLock()
	queuedTriggers := len(gameState.triggeredQueue)
	gameState.mu.RUnlock()

	if queuedTriggers != 1 {
		t.Errorf("Expected 1 triggered ability in queue, got %d", queuedTriggers)
	}

	// Process and resolve triggers
	gameState.mu.Lock()
	engine.processTriggeredAbilities(gameState)
	for !gameState.stack.IsEmpty() {
		item, err := gameState.stack.Pop()
		if err != nil {
			t.Fatalf("Failed to pop from stack: %v", err)
		}
		if item.Resolve != nil {
			if err := item.Resolve(); err != nil {
				t.Fatalf("Failed to resolve: %v", err)
			}
		}
	}
	gameState.mu.Unlock()

	if !triggerFired {
		t.Error("Expected damage to player trigger to fire")
	}

	if damageDealt != 3 {
		t.Errorf("Expected 3 damage dealt, got %d", damageDealt)
	}

	engine.EndCombat(gameID)
}

// TestCombatDamageTriggerCreature verifies "Whenever ~ deals combat damage to a creature" triggers
func TestCombatDamageTriggerCreature(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-damage-creature-trigger"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create attacker with damage trigger and blocker
	gameState.mu.Lock()
	attackerID := "attacker"
	blockerID := "blocker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Damage Trigger Creature",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "4",
		Toughness:    "4",
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

	// Track if trigger fired
	triggerFired := false
	var damageDealt int

	// Register combat damage trigger: "Whenever ~ deals combat damage to a creature, gain 1 life"
	trigger := &combatTrigger{
		SourceID:    attackerID,
		TriggerType: "deals_damage_creature",
		Condition: func(gs *engineGameState, event rules.Event) bool {
			// Check if this is a damaged permanent event from our creature
			return event.Type == rules.EventDamagedPermanent &&
				event.SourceID == attackerID &&
				event.Flag == true // Combat damage
		},
		CreateAbility: func(gs *engineGameState, event rules.Event) *triggeredAbilityQueueItem {
			return &triggeredAbilityQueueItem{
				ID:          fmt.Sprintf("%s-damage-creature-trigger", attackerID),
				SourceID:    attackerID,
				Controller:  "Alice",
				Description: "Whenever Damage Trigger Creature deals combat damage to a creature, gain 1 life",
				Resolve: func(gs *engineGameState) error {
					triggerFired = true
					damageDealt = event.Amount
					// In real implementation, would gain life here
					return nil
				},
				UsesStack: true,
			}
		},
	}

	gameState.combatTriggers = append(gameState.combatTriggers, trigger)
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)

	// Assign and apply damage - should trigger
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Check that trigger was queued
	gameState.mu.RLock()
	queuedTriggers := len(gameState.triggeredQueue)
	gameState.mu.RUnlock()

	if queuedTriggers != 1 {
		t.Errorf("Expected 1 triggered ability in queue, got %d", queuedTriggers)
	}

	// Process and resolve triggers
	gameState.mu.Lock()
	engine.processTriggeredAbilities(gameState)
	for !gameState.stack.IsEmpty() {
		item, err := gameState.stack.Pop()
		if err != nil {
			t.Fatalf("Failed to pop from stack: %v", err)
		}
		if item.Resolve != nil {
			if err := item.Resolve(); err != nil {
				t.Fatalf("Failed to resolve: %v", err)
			}
		}
	}
	gameState.mu.Unlock()

	if !triggerFired {
		t.Error("Expected damage to creature trigger to fire")
	}

	if damageDealt != 4 {
		t.Errorf("Expected 4 damage dealt, got %d", damageDealt)
	}

	engine.EndCombat(gameID)
}

// TestCombatDamageTriggerAny verifies "Whenever ~ deals combat damage" triggers (any target)
func TestCombatDamageTriggerAny(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-damage-any-trigger"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create attacker with generic damage trigger
	gameState.mu.Lock()
	attackerID := "attacker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Generic Damage Trigger",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
	}

	// Track trigger count
	triggerCount := 0

	// Register combat damage trigger: "Whenever ~ deals combat damage, put a +1/+1 counter on it"
	trigger := &combatTrigger{
		SourceID:    attackerID,
		TriggerType: "deals_damage_any",
		Condition: func(gs *engineGameState, event rules.Event) bool {
			// Check if this is any damaged event from our creature (player or permanent)
			return (event.Type == rules.EventDamagedPlayer || event.Type == rules.EventDamagedPermanent) &&
				event.SourceID == attackerID &&
				event.Flag == true // Combat damage
		},
		CreateAbility: func(gs *engineGameState, event rules.Event) *triggeredAbilityQueueItem {
			return &triggeredAbilityQueueItem{
				ID:          fmt.Sprintf("%s-damage-any-trigger-%d", attackerID, triggerCount),
				SourceID:    attackerID,
				Controller:  "Alice",
				Description: "Whenever Generic Damage Trigger deals combat damage, put a +1/+1 counter on it",
				Resolve: func(gs *engineGameState) error {
					triggerCount++
					// In real implementation, would add +1/+1 counter here
					return nil
				},
				UsesStack: true,
			}
		},
	}

	gameState.combatTriggers = append(gameState.combatTriggers, trigger)
	gameState.mu.Unlock()

	// Setup combat - unblocked attacker
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.AcceptBlockers(gameID)

	// Assign and apply damage - should trigger once for player damage
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Process and resolve triggers
	gameState.mu.Lock()
	engine.processTriggeredAbilities(gameState)
	for !gameState.stack.IsEmpty() {
		item, err := gameState.stack.Pop()
		if err != nil {
			t.Fatalf("Failed to pop from stack: %v", err)
		}
		if item.Resolve != nil {
			if err := item.Resolve(); err != nil {
				t.Fatalf("Failed to resolve: %v", err)
			}
		}
	}
	gameState.mu.Unlock()

	if triggerCount != 1 {
		t.Errorf("Expected trigger to fire once, fired %d times", triggerCount)
	}

	engine.EndCombat(gameID)
}

// TestCombatDamageTriggerFirstStrike verifies damage triggers with first strike
func TestCombatDamageTriggerFirstStrike(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-damage-first-strike-trigger"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create first strike attacker with damage trigger
	gameState.mu.Lock()
	attackerID := "attacker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "First Strike Damage Trigger",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "3",
		Toughness:    "3",
		Abilities: []EngineAbilityView{
			{ID: abilityFirstStrike, Text: "First Strike"},
		},
	}

	// Track trigger count
	triggerCount := 0

	// Register combat damage trigger
	trigger := &combatTrigger{
		SourceID:    attackerID,
		TriggerType: "deals_damage_player",
		Condition: func(gs *engineGameState, event rules.Event) bool {
			return event.Type == rules.EventDamagedPlayer &&
				event.SourceID == attackerID &&
				event.Flag == true
		},
		CreateAbility: func(gs *engineGameState, event rules.Event) *triggeredAbilityQueueItem {
			return &triggeredAbilityQueueItem{
				ID:          fmt.Sprintf("%s-damage-trigger-%d", attackerID, triggerCount),
				SourceID:    attackerID,
				Controller:  "Alice",
				Description: "Damage trigger",
				Resolve: func(gs *engineGameState) error {
					triggerCount++
					return nil
				},
				UsesStack: true,
			}
		},
	}

	gameState.combatTriggers = append(gameState.combatTriggers, trigger)
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.AcceptBlockers(gameID)

	// First strike damage step - should trigger
	engine.AssignCombatDamage(gameID, true)
	engine.ApplyCombatDamage(gameID)

	// Process triggers
	gameState.mu.Lock()
	engine.processTriggeredAbilities(gameState)
	for !gameState.stack.IsEmpty() {
		item, err := gameState.stack.Pop()
		if err != nil {
			t.Fatalf("Failed to pop from stack: %v", err)
		}
		if item.Resolve != nil {
			if err := item.Resolve(); err != nil {
				t.Fatalf("Failed to resolve: %v", err)
			}
		}
	}
	gameState.mu.Unlock()

	if triggerCount != 1 {
		t.Errorf("Expected trigger to fire once in first strike, fired %d times", triggerCount)
	}

	// Normal damage step - should NOT trigger again (already dealt damage)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	if triggerCount != 1 {
		t.Errorf("Expected trigger to fire only once total, fired %d times", triggerCount)
	}

	engine.EndCombat(gameID)
}
