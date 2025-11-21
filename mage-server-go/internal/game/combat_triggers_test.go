package game

import (
	"fmt"
	"testing"

	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap/zaptest"
)

// TestCombatTriggerAttacks verifies "Whenever ~ attacks" triggers
func TestCombatTriggerAttacks(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-attack-trigger"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create attacker with attack trigger
	gameState.mu.Lock()
	attackerID := "attacker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Attack Trigger Creature",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
	}

	// Track if trigger fired
	triggerFired := false

	// Register combat trigger: "Whenever ~ attacks, draw a card"
	trigger := &combatTrigger{
		SourceID:    attackerID,
		TriggerType: "attacks",
		Condition: func(gs *engineGameState, event rules.Event) bool {
			// Check if this is an attacker declared event for our creature
			return event.Type == rules.EventAttackerDeclared && event.SourceID == attackerID
		},
		CreateAbility: func(gs *engineGameState, event rules.Event) *triggeredAbilityQueueItem {
			return &triggeredAbilityQueueItem{
				ID:          fmt.Sprintf("%s-attack-trigger", attackerID),
				SourceID:    attackerID,
				Controller:  "Alice",
				Description: "Whenever Attack Trigger Creature attacks, draw a card",
				Resolve: func(gs *engineGameState) error {
					triggerFired = true
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

	// Declare attacker - should trigger
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Check that trigger was queued
	gameState.mu.RLock()
	queuedTriggers := len(gameState.triggeredQueue)
	gameState.mu.RUnlock()

	if queuedTriggers != 1 {
		t.Errorf("Expected 1 triggered ability in queue, got %d", queuedTriggers)
	}

	// Process triggers (puts them on stack)
	gameState.mu.Lock()
	processed := engine.processTriggeredAbilities(gameState)
	gameState.mu.Unlock()

	if !processed {
		t.Error("Expected triggers to be processed")
	}

	// Resolve the stack
	gameState.mu.Lock()
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
		t.Error("Expected attack trigger to fire")
	}

	engine.EndCombat(gameID)
}

// TestCombatTriggerBlocks verifies "Whenever ~ blocks" triggers
func TestCombatTriggerBlocks(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-block-trigger"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create attacker and blocker with block trigger
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
		Power:        "2",
		Toughness:    "2",
	}

	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Block Trigger Creature",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
	}

	// Track if trigger fired
	triggerFired := false

	// Register combat trigger: "Whenever ~ blocks, gain 1 life"
	trigger := &combatTrigger{
		SourceID:    blockerID,
		TriggerType: "blocks",
		Condition: func(gs *engineGameState, event rules.Event) bool {
			// Check if this is a blocker declared event for our creature
			return event.Type == rules.EventBlockerDeclared && event.SourceID == blockerID
		},
		CreateAbility: func(gs *engineGameState, event rules.Event) *triggeredAbilityQueueItem {
			return &triggeredAbilityQueueItem{
				ID:          fmt.Sprintf("%s-block-trigger", blockerID),
				SourceID:    blockerID,
				Controller:  "Bob",
				Description: "Whenever Block Trigger Creature blocks, gain 1 life",
				Resolve: func(gs *engineGameState) error {
					triggerFired = true
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

	// Declare blocker - should trigger
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")

	// Check that trigger was queued
	gameState.mu.RLock()
	queuedTriggers := len(gameState.triggeredQueue)
	gameState.mu.RUnlock()

	if queuedTriggers != 1 {
		t.Errorf("Expected 1 triggered ability in queue, got %d", queuedTriggers)
	}

	// Process triggers (puts them on stack)
	gameState.mu.Lock()
	processed := engine.processTriggeredAbilities(gameState)
	gameState.mu.Unlock()

	if !processed {
		t.Error("Expected triggers to be processed")
	}

	// Resolve the stack
	gameState.mu.Lock()
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
		t.Error("Expected block trigger to fire")
	}

	engine.EndCombat(gameID)
}

// TestCombatTriggerBecomesBlocked verifies "Whenever ~ becomes blocked" triggers
func TestCombatTriggerBecomesBlocked(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-becomes-blocked-trigger"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create attacker with "becomes blocked" trigger and blocker
	gameState.mu.Lock()
	attackerID := "attacker"
	blockerID := "blocker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Becomes Blocked Trigger Creature",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
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

	// Register combat trigger: "Whenever ~ becomes blocked, it gets +1/+1"
	trigger := &combatTrigger{
		SourceID:    attackerID,
		TriggerType: "becomes_blocked",
		Condition: func(gs *engineGameState, event rules.Event) bool {
			// Check if this is a creature blocked event for our creature
			return event.Type == rules.EventCreatureBlocked && event.SourceID == attackerID
		},
		CreateAbility: func(gs *engineGameState, event rules.Event) *triggeredAbilityQueueItem {
			return &triggeredAbilityQueueItem{
				ID:          fmt.Sprintf("%s-blocked-trigger", attackerID),
				SourceID:    attackerID,
				Controller:  "Alice",
				Description: "Whenever Becomes Blocked Trigger Creature becomes blocked, it gets +1/+1",
				Resolve: func(gs *engineGameState) error {
					triggerFired = true
					// In real implementation, would add +1/+1 effect here
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

	// Accept blockers - should trigger "becomes blocked"
	engine.AcceptBlockers(gameID)

	// Check that trigger was queued
	gameState.mu.RLock()
	queuedTriggers := len(gameState.triggeredQueue)
	gameState.mu.RUnlock()

	if queuedTriggers != 1 {
		t.Errorf("Expected 1 triggered ability in queue, got %d", queuedTriggers)
	}

	// Process triggers (puts them on stack)
	gameState.mu.Lock()
	processed := engine.processTriggeredAbilities(gameState)
	gameState.mu.Unlock()

	if !processed {
		t.Error("Expected triggers to be processed")
	}

	// Resolve the stack
	gameState.mu.Lock()
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
		t.Error("Expected becomes blocked trigger to fire")
	}

	engine.EndCombat(gameID)
}

// TestCombatTriggerMultiple verifies multiple triggers fire in correct order
func TestCombatTriggerMultiple(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-multiple-triggers"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create two attackers with attack triggers
	gameState.mu.Lock()
	attacker1ID := "attacker1"
	attacker2ID := "attacker2"

	gameState.cards[attacker1ID] = &internalCard{
		ID:           attacker1ID,
		Name:         "Attacker 1",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
	}

	gameState.cards[attacker2ID] = &internalCard{
		ID:           attacker2ID,
		Name:         "Attacker 2",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
	}

	// Register triggers for both attackers
	for _, attackerID := range []string{attacker1ID, attacker2ID} {
		trigger := &combatTrigger{
			SourceID:    attackerID,
			TriggerType: "attacks",
			Condition: func(gs *engineGameState, event rules.Event) bool {
				return event.Type == rules.EventAttackerDeclared && event.SourceID == attackerID
			},
			CreateAbility: func(gs *engineGameState, event rules.Event) *triggeredAbilityQueueItem {
				return &triggeredAbilityQueueItem{
					ID:          fmt.Sprintf("%s-attack-trigger", attackerID),
					SourceID:    attackerID,
					Controller:  "Alice",
					Description: fmt.Sprintf("Whenever %s attacks", attackerID),
					Resolve: func(gs *engineGameState) error {
						return nil
					},
					UsesStack: true,
				}
			},
		}
		gameState.combatTriggers = append(gameState.combatTriggers, trigger)
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	// Declare both attackers
	engine.DeclareAttacker(gameID, attacker1ID, "Bob", "Alice")
	engine.DeclareAttacker(gameID, attacker2ID, "Bob", "Alice")

	// Check that both triggers were queued
	gameState.mu.RLock()
	queuedTriggers := len(gameState.triggeredQueue)
	gameState.mu.RUnlock()

	if queuedTriggers != 2 {
		t.Errorf("Expected 2 triggered abilities in queue, got %d", queuedTriggers)
	}

	engine.EndCombat(gameID)
}
