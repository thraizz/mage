package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestCombatTrampleUnblocked tests that trample creatures deal damage when unblocked
func TestCombatTrampleUnblocked(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-trample-unblocked"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Trample attacker
	gameState.mu.Lock()
	attackerID := "trample-attacker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Colossal Dreadmaw",
		Type:         "Creature - Dinosaur",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "6",
		Toughness:    "6",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityTrample, Text: "Trample"},
		},
	}

	initialBobLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()

	// Full combat flow
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// No blockers
	engine.AcceptBlockers(gameID)

	// Damage
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify damage dealt to Bob
	gameState.mu.RLock()
	bobLife := gameState.players["Bob"].Life
	if bobLife != initialBobLife-6 {
		t.Errorf("expected Bob to lose 6 life, lost %d", initialBobLife-bobLife)
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}

// TestCombatTrampleOverBlocker tests that trample damage carries over to defender
func TestCombatTrampleOverBlocker(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-trample-over"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: 6/6 trample attacker vs 2/2 blocker
	gameState.mu.Lock()
	attackerID := "trample-attacker"
	blockerID := "small-blocker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Colossal Dreadmaw",
		Type:         "Creature - Dinosaur",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "6",
		Toughness:    "6",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityTrample, Text: "Trample"},
		},
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

	// Combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)

	// Damage
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify:
	// - Blocker takes 2 damage (lethal) and dies
	// - Remaining 4 damage tramples through to Bob
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	bobLife := gameState.players["Bob"].Life

	if blocker.Zone != zoneGraveyard {
		t.Error("blocker should be in graveyard")
	}

	expectedDamage := 4 // 6 power - 2 lethal to blocker = 4 trample
	if bobLife != initialBobLife-expectedDamage {
		t.Errorf("expected Bob to lose %d life from trample, lost %d", expectedDamage, initialBobLife-bobLife)
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}

// TestCombatTrampleExactLethal tests trample with exact lethal damage to blocker
func TestCombatTrampleExactLethal(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-trample-exact"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: 2/2 trample attacker vs 2/2 blocker (exact lethal, no trample through)
	gameState.mu.Lock()
	attackerID := "trample-attacker"
	blockerID := "equal-blocker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Trample Bear",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityTrample, Text: "Trample"},
		},
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

	// Combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)

	// Damage
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify:
	// - Blocker takes 2 damage (lethal) and dies
	// - No trample damage (all damage was lethal)
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	attacker := gameState.cards[attackerID]
	bobLife := gameState.players["Bob"].Life

	if blocker.Zone != zoneGraveyard {
		t.Error("blocker should be in graveyard")
	}

	if attacker.Zone != zoneGraveyard {
		t.Error("attacker should be in graveyard (took 2 damage)")
	}

	if bobLife != initialBobLife {
		t.Errorf("Bob should not lose life (no trample through), lost %d", initialBobLife-bobLife)
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}

// TestCombatTrampleMultipleBlockers tests trample with multiple blockers
func TestCombatTrampleMultipleBlockers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-trample-multi"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: 10/10 trample attacker vs two 2/2 blockers
	gameState.mu.Lock()
	attackerID := "trample-attacker"
	blocker1ID := "blocker-1"
	blocker2ID := "blocker-2"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Giant Trampler",
		Type:         "Creature - Giant",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "10",
		Toughness:    "10",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityTrample, Text: "Trample"},
		},
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
		Name:         "Runeclaw Bear",
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

	// Combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blocker1ID, attackerID, "Bob")
	engine.DeclareBlocker(gameID, blocker2ID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)

	// Damage
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify:
	// - Both blockers take 2 damage (lethal) and die
	// - Remaining 6 damage tramples through to Bob (10 - 2 - 2 = 6)
	gameState.mu.RLock()
	blocker1 := gameState.cards[blocker1ID]
	blocker2 := gameState.cards[blocker2ID]
	bobLife := gameState.players["Bob"].Life

	if blocker1.Zone != zoneGraveyard {
		t.Error("blocker 1 should be in graveyard")
	}

	if blocker2.Zone != zoneGraveyard {
		t.Error("blocker 2 should be in graveyard")
	}

	expectedDamage := 6 // 10 power - 2 - 2 = 6 trample
	if bobLife != initialBobLife-expectedDamage {
		t.Errorf("expected Bob to lose %d life from trample, lost %d", expectedDamage, initialBobLife-bobLife)
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}

// TestCombatNoTrampleBlocked tests that creatures without trample don't deal damage to defender when blocked
func TestCombatNoTrampleBlocked(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-no-trample"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: 6/6 normal attacker vs 2/2 blocker (no trample)
	gameState.mu.Lock()
	attackerID := "normal-attacker"
	blockerID := "small-blocker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Big Creature",
		Type:         "Creature - Beast",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "6",
		Toughness:    "6",
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

	// Combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)

	// Damage
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify:
	// - Blocker takes damage and dies
	// - Bob takes NO damage (no trample)
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	bobLife := gameState.players["Bob"].Life

	if blocker.Zone != zoneGraveyard {
		t.Error("blocker should be in graveyard")
	}

	if bobLife != initialBobLife {
		t.Errorf("Bob should not lose life (no trample), lost %d", initialBobLife-bobLife)
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}

// TestCombatTrampleWithFirstStrike tests trample combined with first strike
func TestCombatTrampleWithFirstStrike(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-trample-first-strike"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: 4/4 trample + first strike attacker vs 2/2 blocker
	gameState.mu.Lock()
	attackerID := "trample-first-strike"
	blockerID := "normal-blocker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Elite Trampler",
		Type:         "Creature - Soldier",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "4",
		Toughness:    "4",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityTrample, Text: "Trample"},
			{ID: abilityFirstStrike, Text: "First strike"},
		},
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

	// Combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)

	// First strike damage
	engine.AssignCombatDamage(gameID, true)
	engine.ApplyCombatDamage(gameID)

	// Normal damage (blocker is dead, so no damage back)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify:
	// - Blocker dies in first strike step
	// - 2 trample damage goes through to Bob in first strike step
	// - Attacker survives (blocker dealt no damage)
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	attacker := gameState.cards[attackerID]
	bobLife := gameState.players["Bob"].Life

	if blocker.Zone != zoneGraveyard {
		t.Error("blocker should be in graveyard")
	}

	if attacker.Zone != zoneBattlefield {
		t.Error("attacker should still be on battlefield")
	}

	if attacker.Damage != 0 {
		t.Errorf("attacker should have no damage, has %d", attacker.Damage)
	}

	expectedDamage := 2 // 4 power - 2 lethal = 2 trample
	if bobLife != initialBobLife-expectedDamage {
		t.Errorf("expected Bob to lose %d life from trample, lost %d", expectedDamage, initialBobLife-bobLife)
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}

// TestCombatTrampleInsufficientPower tests trample when attacker has less power than blocker toughness
func TestCombatTrampleInsufficientPower(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-trample-insufficient"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: 2/2 trample attacker vs 5/5 blocker
	gameState.mu.Lock()
	attackerID := "weak-trampler"
	blockerID := "big-blocker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Small Trampler",
		Type:         "Creature - Beast",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityTrample, Text: "Trample"},
		},
	}

	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Big Wall",
		Type:         "Creature - Wall",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "0",
		Toughness:    "5",
		Tapped:       false,
	}

	initialBobLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()

	// Combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)

	// Damage
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify:
	// - Blocker takes 2 damage (not lethal, survives)
	// - No trample damage (all damage assigned to blocker)
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	bobLife := gameState.players["Bob"].Life

	if blocker.Zone != zoneBattlefield {
		t.Error("blocker should still be on battlefield")
	}

	if blocker.Damage != 2 {
		t.Errorf("blocker should have 2 damage, has %d", blocker.Damage)
	}

	if bobLife != initialBobLife {
		t.Errorf("Bob should not lose life (no trample through), lost %d", initialBobLife-bobLife)
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}
