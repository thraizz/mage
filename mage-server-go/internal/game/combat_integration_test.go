package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap/zaptest"
)

// TestCombatFullFlowUnblocked tests a complete combat flow with an unblocked attacker
// This is an integration test covering: Reset -> Declare Attackers -> Damage -> Cleanup
func TestCombatFullFlowUnblocked(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-full-combat-unblocked"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create a 3/3 attacker for Alice
	gameState.mu.Lock()
	attackerID := "attacker-1"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Trained Armodon",
		Type:         "Creature - Elephant",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "3",
		Toughness:    "3",
		Tapped:       false,
	}

	initialBobLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()

	// Phase 1: Reset Combat (Beginning of Combat)
	if err := engine.ResetCombat(gameID); err != nil {
		t.Fatalf("failed to reset combat: %v", err)
	}

	// Verify initial state
	gameState.mu.RLock()
	if gameState.cards[attackerID].Attacking {
		t.Error("creature should not be attacking after reset")
	}
	gameState.mu.RUnlock()

	// Phase 2: Declare Attackers Step
	if err := engine.SetAttacker(gameID, "Alice"); err != nil {
		t.Fatalf("failed to set attacker: %v", err)
	}

	if err := engine.SetDefenders(gameID); err != nil {
		t.Fatalf("failed to set defenders: %v", err)
	}

	if err := engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice"); err != nil {
		t.Fatalf("failed to declare attacker: %v", err)
	}

	// Verify attacker is attacking
	gameState.mu.RLock()
	attacker := gameState.cards[attackerID]
	if !attacker.Attacking {
		t.Error("creature should be attacking after declaration")
	}
	if attacker.AttackingWhat != "Bob" {
		t.Errorf("creature should be attacking Bob, got %s", attacker.AttackingWhat)
	}
	if !attacker.Tapped {
		t.Error("attacking creature should be tapped")
	}
	gameState.mu.RUnlock()

	// Phase 3: Declare Blockers Step (no blockers)
	// Skip - no blockers declared

	// Phase 4: Combat Damage Step
	if err := engine.AssignCombatDamage(gameID, false); err != nil {
		t.Fatalf("failed to assign combat damage: %v", err)
	}

	if err := engine.ApplyCombatDamage(gameID); err != nil {
		t.Fatalf("failed to apply combat damage: %v", err)
	}

	// Verify damage was dealt to Bob
	gameState.mu.RLock()
	bobLife := gameState.players["Bob"].Life
	expectedLife := initialBobLife - 3
	if bobLife != expectedLife {
		t.Errorf("expected Bob's life to be %d, got %d", expectedLife, bobLife)
	}

	// Verify attacker took no damage
	attacker = gameState.cards[attackerID]
	if attacker.Damage != 0 {
		t.Errorf("unblocked attacker should have no damage, got %d", attacker.Damage)
	}
	gameState.mu.RUnlock()

	// Phase 5: End of Combat Step
	if err := engine.EndCombat(gameID); err != nil {
		t.Fatalf("failed to end combat: %v", err)
	}

	// Verify combat cleanup
	gameState.mu.RLock()
	attacker = gameState.cards[attackerID]
	if attacker.Attacking {
		t.Error("creature should not be attacking after combat ends")
	}
	if attacker.AttackingWhat != "" {
		t.Error("AttackingWhat should be cleared after combat ends")
	}
	if attacker.Damage != 0 {
		t.Error("damage should be cleared after combat ends")
	}

	// Verify creature still exists (didn't die)
	if attacker.Zone != zoneBattlefield {
		t.Error("attacker should still be on battlefield")
	}

	// Verify historical tracking
	if len(gameState.combat.formerGroups) == 0 {
		t.Error("former groups should be preserved after combat")
	}
	gameState.mu.RUnlock()

	// Verify "attacked this turn" query
	attacked, err := engine.GetAttackedThisTurn(gameID, attackerID)
	if err != nil {
		t.Fatalf("failed to check attacked this turn: %v", err)
	}
	if !attacked {
		t.Error("creature should show as attacked this turn")
	}
}

// TestCombatFullFlowBlocked tests a complete combat flow with blockers and lethal damage
// This is an integration test covering: Reset -> Attackers -> Blockers -> Damage -> Death -> Cleanup
func TestCombatFullFlowBlocked(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-full-combat-blocked"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Create a 2/2 attacker and a 2/2 blocker
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
		Name:         "Grizzly Bears",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}

	initialBobLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()

	// Phase 1: Reset Combat
	if err := engine.ResetCombat(gameID); err != nil {
		t.Fatalf("failed to reset combat: %v", err)
	}

	// Phase 2: Declare Attackers
	if err := engine.SetAttacker(gameID, "Alice"); err != nil {
		t.Fatalf("failed to set attacker: %v", err)
	}

	if err := engine.SetDefenders(gameID); err != nil {
		t.Fatalf("failed to set defenders: %v", err)
	}

	if err := engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice"); err != nil {
		t.Fatalf("failed to declare attacker: %v", err)
	}

	// Verify attacker state
	gameState.mu.RLock()
	if !gameState.cards[attackerID].Attacking {
		t.Error("creature should be attacking")
	}
	gameState.mu.RUnlock()

	// Phase 3: Declare Blockers
	if err := engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob"); err != nil {
		t.Fatalf("failed to declare blocker: %v", err)
	}

	if err := engine.AcceptBlockers(gameID); err != nil {
		t.Fatalf("failed to accept blockers: %v", err)
	}

	// Verify blocker state
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	if !blocker.Blocking {
		t.Error("creature should be blocking")
	}
	if len(blocker.BlockingWhat) != 1 || blocker.BlockingWhat[0] != attackerID {
		t.Error("blocker should be blocking the attacker")
	}
	gameState.mu.RUnlock()

	// Phase 4: Combat Damage Step
	if err := engine.AssignCombatDamage(gameID, false); err != nil {
		t.Fatalf("failed to assign combat damage: %v", err)
	}

	// Verify damage is marked (before applying)
	gameState.mu.RLock()
	attacker := gameState.cards[attackerID]
	blocker = gameState.cards[blockerID]
	if attacker.Damage != 2 {
		t.Errorf("attacker should have 2 damage marked, got %d", attacker.Damage)
	}
	if blocker.Damage != 2 {
		t.Errorf("blocker should have 2 damage marked, got %d", blocker.Damage)
	}
	gameState.mu.RUnlock()

	if err := engine.ApplyCombatDamage(gameID); err != nil {
		t.Fatalf("failed to apply combat damage: %v", err)
	}

	// Verify Bob took no damage (attacker was blocked)
	gameState.mu.RLock()
	bobLife := gameState.players["Bob"].Life
	if bobLife != initialBobLife {
		t.Errorf("Bob should not have taken damage (attacker was blocked), life: %d -> %d", initialBobLife, bobLife)
	}

	// Verify both creatures died (lethal damage)
	attacker = gameState.cards[attackerID]
	blocker = gameState.cards[blockerID]

	if attacker.Zone != zoneGraveyard {
		t.Errorf("attacker should be in graveyard, in zone %d", attacker.Zone)
	}
	if blocker.Zone != zoneGraveyard {
		t.Errorf("blocker should be in graveyard, in zone %d", blocker.Zone)
	}
	gameState.mu.RUnlock()

	// Phase 5: End of Combat
	if err := engine.EndCombat(gameID); err != nil {
		t.Fatalf("failed to end combat: %v", err)
	}

	// Verify combat cleanup (even for dead creatures)
	gameState.mu.RLock()
	attacker = gameState.cards[attackerID]
	blocker = gameState.cards[blockerID]

	if attacker.Attacking {
		t.Error("dead attacker should not be attacking")
	}
	if blocker.Blocking {
		t.Error("dead blocker should not be blocking")
	}

	// Verify historical tracking
	if len(gameState.combat.formerGroups) == 0 {
		t.Error("former groups should be preserved")
	}
	gameState.mu.RUnlock()

	// Verify "attacked this turn" still works for dead creature
	attacked, err := engine.GetAttackedThisTurn(gameID, attackerID)
	if err != nil {
		t.Fatalf("failed to check attacked this turn: %v", err)
	}
	if !attacked {
		t.Error("dead creature should still show as attacked this turn")
	}
}

// TestCombatFullFlowMultipleBlockers tests combat with multiple blockers on one attacker
func TestCombatFullFlowMultipleBlockers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-full-combat-multi-blockers"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: 5/5 attacker vs two 2/2 blockers
	gameState.mu.Lock()
	attackerID := "attacker-1"
	blocker1ID := "blocker-1"
	blocker2ID := "blocker-2"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Craw Wurm",
		Type:         "Creature - Wurm",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "6",
		Toughness:    "4",
		Tapped:       false,
	}

	gameState.cards[blocker1ID] = &internalCard{
		ID:           blocker1ID,
		Name:         "Grizzly Bears",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}

	gameState.cards[blocker2ID] = &internalCard{
		ID:           blocker2ID,
		Name:         "Grizzly Bears",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}

	initialBobLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()

	// Full combat flow
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blocker1ID, attackerID, "Bob")
	engine.DeclareBlocker(gameID, blocker2ID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)

	// Verify both blockers are blocking
	gameState.mu.RLock()
	if !gameState.cards[blocker1ID].Blocking {
		t.Error("blocker 1 should be blocking")
	}
	if !gameState.cards[blocker2ID].Blocking {
		t.Error("blocker 2 should be blocking")
	}
	gameState.mu.RUnlock()

	// Damage phase
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify results
	gameState.mu.RLock()
	bobLife := gameState.players["Bob"].Life

	// Bob should take no damage (attacker was blocked)
	if bobLife != initialBobLife {
		t.Errorf("Bob should not have taken damage, life: %d -> %d", initialBobLife, bobLife)
	}

	// Both blockers should be dead (took 3 damage each, have 2 toughness)
	blocker1 := gameState.cards[blocker1ID]
	blocker2 := gameState.cards[blocker2ID]
	if blocker1.Zone != zoneGraveyard {
		t.Error("blocker 1 should be in graveyard")
	}
	if blocker2.Zone != zoneGraveyard {
		t.Error("blocker 2 should be in graveyard")
	}

	// Attacker should have taken 4 damage (2+2) and be dead (6/4 with 4 damage = lethal)
	attacker := gameState.cards[attackerID]
	if attacker.Zone != zoneGraveyard {
		t.Error("attacker should be in graveyard (took lethal damage)")
	}
	gameState.mu.RUnlock()

	// Cleanup
	engine.EndCombat(gameID)

	// Verify cleanup (even dead creatures get cleaned up)
	gameState.mu.RLock()
	attacker = gameState.cards[attackerID]
	if attacker.Attacking {
		t.Error("dead attacker should not be attacking after cleanup")
	}
	if attacker.Damage != 0 {
		t.Error("dead attacker damage should be cleared after cleanup")
	}
	// Verify it's still in graveyard
	if attacker.Zone != zoneGraveyard {
		t.Error("attacker should remain in graveyard after cleanup")
	}
	gameState.mu.RUnlock()
}

// TestCombatFullFlowWithEvents tests that all expected events are fired during combat
func TestCombatFullFlowWithEvents(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-full-combat-events"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Track events
	var events []rules.EventType
	eventsMu := make(chan struct{}, 1)
	eventsMu <- struct{}{}

	// Subscribe to all combat-related events
	gameState.eventBus.SubscribeTyped(rules.EventAttackerDeclared, func(e rules.Event) {
		<-eventsMu
		events = append(events, e.Type)
		eventsMu <- struct{}{}
	})
	gameState.eventBus.SubscribeTyped(rules.EventBlockerDeclared, func(e rules.Event) {
		<-eventsMu
		events = append(events, e.Type)
		eventsMu <- struct{}{}
	})
	gameState.eventBus.SubscribeTyped(rules.EventDeclaredBlockers, func(e rules.Event) {
		<-eventsMu
		events = append(events, e.Type)
		eventsMu <- struct{}{}
	})
	gameState.eventBus.SubscribeTyped(rules.EventDamagePlayer, func(e rules.Event) {
		<-eventsMu
		events = append(events, e.Type)
		eventsMu <- struct{}{}
	})
	gameState.eventBus.SubscribeTyped(rules.EventEndCombatStep, func(e rules.Event) {
		<-eventsMu
		events = append(events, e.Type)
		eventsMu <- struct{}{}
	})

	// Setup creatures
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
		Name:         "Wall of Stone",
		Type:         "Creature - Wall",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "0",
		Toughness:    "8",
		Tapped:       false,
	}
	gameState.mu.Unlock()

	// Run full combat flow
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	engine.EndCombat(gameID)

	// Verify events were fired
	<-eventsMu
	if len(events) == 0 {
		t.Error("no events were fired during combat")
	}

	// Check for expected events
	hasAttackerDeclared := false
	hasBlockerDeclared := false
	hasDeclaredBlockers := false
	hasEndCombat := false

	for _, evt := range events {
		switch evt {
		case rules.EventAttackerDeclared:
			hasAttackerDeclared = true
		case rules.EventBlockerDeclared:
			hasBlockerDeclared = true
		case rules.EventDeclaredBlockers:
			hasDeclaredBlockers = true
		case rules.EventEndCombatStep:
			hasEndCombat = true
		}
	}

	if !hasAttackerDeclared {
		t.Error("EventAttackerDeclared was not fired")
	}
	if !hasBlockerDeclared {
		t.Error("EventBlockerDeclared was not fired")
	}
	if !hasDeclaredBlockers {
		t.Error("EventDeclaredBlockers was not fired")
	}
	if !hasEndCombat {
		t.Error("EventEndCombatStep was not fired")
	}
}
