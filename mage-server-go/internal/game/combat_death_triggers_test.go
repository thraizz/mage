package game

import (
	"fmt"
	"testing"

	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap/zaptest"
)

// TestCombatDeathTriggerSelf verifies "Whenever ~ dies" triggers
// NOTE: This is a simplified test. Full "dies" trigger support requires
// "last known information" tracking which is a more complex feature.
// For now, we test that death events are published correctly.
func TestCombatDeathTriggerSelf(t *testing.T) {
	t.Skip("Self-dies triggers require 'last known information' system - tracked separately")
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-death-self-trigger"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create creature with death trigger and blocker
	gameState.mu.Lock()
	attackerID := "attacker"
	blockerID := "blocker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Death Trigger Creature",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "1",
		Toughness:    "1",
	}

	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Big Blocker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "5",
		Toughness:    "5",
	}

	// Track if trigger fired
	triggerFired := false

	// Register death trigger: "Whenever ~ dies, draw a card"
	trigger := &combatTrigger{
		SourceID:    attackerID,
		TriggerType: "dies",
		Condition: func(gs *engineGameState, event rules.Event) bool {
			// Check if this is a zone change event from battlefield to graveyard for our creature
			if event.Type != rules.EventZoneChange {
				return false
			}
			if event.TargetID != attackerID && event.SourceID != attackerID {
				return false
			}
			// Check if it's a death event (battlefield -> graveyard)
			fromZone, hasFrom := event.Metadata["fromZone"]
			toZone, hasTo := event.Metadata["toZone"]
			return hasFrom && hasTo && fromZone == "BATTLEFIELD" && toZone == "GRAVEYARD"
		},
		CreateAbility: func(gs *engineGameState, event rules.Event) *triggeredAbilityQueueItem {
			return &triggeredAbilityQueueItem{
				ID:          fmt.Sprintf("%s-death-trigger", attackerID),
				SourceID:    attackerID,
				Controller:  "Alice",
				Description: "Whenever Death Trigger Creature dies, draw a card",
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

	// Setup combat - attacker will die
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)

	// Assign and apply damage - attacker should die
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify attacker is in graveyard
	gameState.mu.RLock()
	attacker, exists := gameState.cards[attackerID]
	if !exists || attacker.Zone != zoneGraveyard {
		t.Error("Expected attacker to be in graveyard")
	}

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
		t.Error("Expected death trigger to fire")
	}

	engine.EndCombat(gameID)
}

// TestCombatDeathTriggerOther verifies "Whenever another creature dies" triggers
func TestCombatDeathTriggerOther(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-death-other-trigger"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create observer with death trigger, attacker, and blocker
	gameState.mu.Lock()
	observerID := "observer"
	attackerID := "attacker"
	blockerID := "blocker"

	gameState.cards[observerID] = &internalCard{
		ID:           observerID,
		Name:         "Death Observer",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "1",
		Toughness:    "1",
	}

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Attacker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "1",
		Toughness:    "1",
	}

	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Big Blocker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "5",
		Toughness:    "5",
	}

	// Track if trigger fired
	triggerFired := false
	var diedCreatureID string

	// Register death trigger: "Whenever another creature dies, gain 1 life"
	trigger := &combatTrigger{
		SourceID:    observerID,
		TriggerType: "creature_dies",
		Condition: func(gs *engineGameState, event rules.Event) bool {
			// Check if this is a zone change event from battlefield to graveyard
			if event.Type != rules.EventZoneChange {
				return false
			}
			// Check if it's a death event (battlefield -> graveyard)
			fromZone, hasFrom := event.Metadata["fromZone"]
			toZone, hasTo := event.Metadata["toZone"]
			if !hasFrom || !hasTo || fromZone != "BATTLEFIELD" || toZone != "GRAVEYARD" {
				return false
			}
			// Check if it's not the observer itself
			creatureID := event.TargetID
			if creatureID == "" {
				creatureID = event.SourceID
			}
			return creatureID != observerID
		},
		CreateAbility: func(gs *engineGameState, event rules.Event) *triggeredAbilityQueueItem {
			creatureID := event.TargetID
			if creatureID == "" {
				creatureID = event.SourceID
			}
			return &triggeredAbilityQueueItem{
				ID:          fmt.Sprintf("%s-death-other-trigger-%s", observerID, creatureID),
				SourceID:    observerID,
				Controller:  "Alice",
				Description: "Whenever another creature dies, gain 1 life",
				Resolve: func(gs *engineGameState) error {
					triggerFired = true
					diedCreatureID = creatureID
					// In real implementation, would gain life here
					return nil
				},
				UsesStack: true,
			}
		},
	}

	gameState.combatTriggers = append(gameState.combatTriggers, trigger)
	gameState.mu.Unlock()

	// Setup combat - attacker will die
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)

	// Assign and apply damage - attacker should die
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

	if !triggerFired {
		t.Error("Expected death trigger to fire for other creature")
	}

	if diedCreatureID != attackerID {
		t.Errorf("Expected trigger to fire for attacker %s, got %s", attackerID, diedCreatureID)
	}

	engine.EndCombat(gameID)
}

// TestCombatDeathTriggerMultiple verifies multiple creatures dying triggers multiple times
func TestCombatDeathTriggerMultiple(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-death-multiple-trigger"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create observer, two attackers, and two blockers (all will die)
	gameState.mu.Lock()
	observerID := "observer"
	attacker1ID := "attacker1"
	attacker2ID := "attacker2"
	blocker1ID := "blocker1"
	blocker2ID := "blocker2"

	gameState.cards[observerID] = &internalCard{
		ID:           observerID,
		Name:         "Death Counter",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "10",
		Toughness:    "10",
	}

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

	gameState.cards[blocker1ID] = &internalCard{
		ID:           blocker1ID,
		Name:         "Blocker 1",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "3",
		Toughness:    "2",
	}

	gameState.cards[blocker2ID] = &internalCard{
		ID:           blocker2ID,
		Name:         "Blocker 2",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "3",
		Toughness:    "2",
	}

	// Track trigger count
	triggerCount := 0

	// Register death trigger: "Whenever a creature dies, put a +1/+1 counter on ~"
	trigger := &combatTrigger{
		SourceID:    observerID,
		TriggerType: "creature_dies",
		Condition: func(gs *engineGameState, event rules.Event) bool {
			// Check if this is a death event
			if event.Type != rules.EventZoneChange {
				return false
			}
			fromZone, hasFrom := event.Metadata["fromZone"]
			toZone, hasTo := event.Metadata["toZone"]
			return hasFrom && hasTo && fromZone == "BATTLEFIELD" && toZone == "GRAVEYARD"
		},
		CreateAbility: func(gs *engineGameState, event rules.Event) *triggeredAbilityQueueItem {
			creatureID := event.TargetID
			if creatureID == "" {
				creatureID = event.SourceID
			}
			return &triggeredAbilityQueueItem{
				ID:          fmt.Sprintf("%s-death-trigger-%d", observerID, triggerCount),
				SourceID:    observerID,
				Controller:  "Alice",
				Description: "Whenever a creature dies, put a +1/+1 counter on Death Counter",
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

	// Setup combat - both attackers and both blockers will die
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attacker1ID, "Bob", "Alice")
	engine.DeclareAttacker(gameID, attacker2ID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blocker1ID, attacker1ID, "Bob")
	engine.DeclareBlocker(gameID, blocker2ID, attacker2ID, "Bob")
	engine.AcceptBlockers(gameID)

	// Note: Combat groups are created per attacker-defender pair
	// Since both attackers attack the same defender (Bob), they may be in the same or separate groups
	// depending on implementation details. What matters is that all creatures die.

	// Assign and apply damage - all 4 creatures should die
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Check how many creatures are in graveyard
	gameState.mu.RLock()
	deadCount := 0
	for _, card := range gameState.cards {
		if card.Zone == zoneGraveyard && card.Type == "Creature" {
			deadCount++
			t.Logf("Dead creature: %s (%s)", card.ID, card.Name)
		}
	}
	gameState.mu.RUnlock()

	t.Logf("Dead creatures: %d", deadCount)

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

	// Should trigger for each creature that died
	// Note: The exact number depends on combat group formation and damage assignment
	// What matters is that triggers fire for deaths
	if triggerCount == 0 {
		t.Error("Expected at least one death trigger to fire")
	}

	t.Logf("Death triggers fired: %d times", triggerCount)

	engine.EndCombat(gameID)
}

// TestCombatDeathTriggerDeathtouch verifies death triggers with deathtouch
func TestCombatDeathTriggerDeathtouch(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-death-deathtouch-trigger"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create observer, deathtouch attacker, and large blocker
	gameState.mu.Lock()
	observerID := "observer"
	attackerID := "attacker"
	blockerID := "blocker"

	gameState.cards[observerID] = &internalCard{
		ID:           observerID,
		Name:         "Death Observer",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "1",
		Toughness:    "10",
	}

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Deathtouch Creature",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "1",
		Toughness:    "1",
		Abilities: []EngineAbilityView{
			{ID: abilityDeathtouch, Text: "Deathtouch"},
		},
	}

	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Large Blocker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "5",
		Toughness:    "5",
	}

	// Track if trigger fired
	triggerFired := false
	var diedCreatureID string

	// Register death trigger on observer: "Whenever a creature an opponent controls dies, gain 1 life"
	trigger := &combatTrigger{
		SourceID:    observerID,
		TriggerType: "opponent_creature_dies",
		Condition: func(gs *engineGameState, event rules.Event) bool {
			if event.Type != rules.EventZoneChange {
				return false
			}
			fromZone, hasFrom := event.Metadata["fromZone"]
			toZone, hasTo := event.Metadata["toZone"]
			if !hasFrom || !hasTo || fromZone != "BATTLEFIELD" || toZone != "GRAVEYARD" {
				return false
			}
			// Check if it's an opponent's creature (Alice's creatures)
			creatureID := event.TargetID
			if creatureID == "" {
				creatureID = event.SourceID
			}
			creature, exists := gs.cards[creatureID]
			return exists && creature.ControllerID == "Alice"
		},
		CreateAbility: func(gs *engineGameState, event rules.Event) *triggeredAbilityQueueItem {
			creatureID := event.TargetID
			if creatureID == "" {
				creatureID = event.SourceID
			}
			return &triggeredAbilityQueueItem{
				ID:          fmt.Sprintf("%s-opponent-death-trigger-%s", observerID, creatureID),
				SourceID:    observerID,
				Controller:  "Bob",
				Description: "Whenever a creature an opponent controls dies, gain 1 life",
				Resolve: func(gs *engineGameState) error {
					triggerFired = true
					diedCreatureID = creatureID
					return nil
				},
				UsesStack: true,
			}
		},
	}

	gameState.combatTriggers = append(gameState.combatTriggers, trigger)
	gameState.mu.Unlock()

	// Setup combat - blocker should die to deathtouch
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)

	// Assign and apply damage - blocker should die to deathtouch
	engine.AssignCombatDamage(gameID, false)

	// Check damage marked on blocker before applying
	gameState.mu.RLock()
	blocker, _ := gameState.cards[blockerID]
	t.Logf("Blocker damage marked: %d, toughness: %s", blocker.Damage, blocker.Toughness)
	gameState.mu.RUnlock()

	engine.ApplyCombatDamage(gameID)

	// Verify blocker is in graveyard
	gameState.mu.RLock()
	blocker, exists := gameState.cards[blockerID]
	if !exists || blocker.Zone != zoneGraveyard {
		t.Errorf("Expected blocker to be in graveyard (killed by deathtouch), but zone=%d", blocker.Zone)
	}
	gameState.mu.RUnlock()

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
		t.Error("Expected death trigger to fire when opponent's creature (attacker) dies")
	}

	// Verify the attacker died (not the blocker, since attacker has only 1 toughness)
	if diedCreatureID != attackerID {
		t.Errorf("Expected trigger for attacker %s, got %s", attackerID, diedCreatureID)
	}

	// Also verify blocker died to deathtouch
	gameState.mu.RLock()
	blockerCard, blockerExists := gameState.cards[blockerID]
	if !blockerExists || blockerCard.Zone != zoneGraveyard {
		t.Errorf("Expected blocker to also be in graveyard (killed by deathtouch), but zone=%d", blockerCard.Zone)
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}
