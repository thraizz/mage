package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestCombatFlyingCannotBeBlocked tests that flying creatures cannot be blocked by non-flying/non-reach creatures
func TestCombatFlyingCannotBeBlocked(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-flying-unblockable"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Flying attacker and normal blocker
	gameState.mu.Lock()
	flyingAttackerID := "flying-attacker"
	normalBlockerID := "normal-blocker"
	
	gameState.cards[flyingAttackerID] = &internalCard{
		ID:           flyingAttackerID,
		Name:         "Wind Drake",
		Type:         "Creature - Drake",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityFlying, Text: "Flying"},
		},
	}
	
	gameState.cards[normalBlockerID] = &internalCard{
		ID:           normalBlockerID,
		Name:         "Grizzly Bears",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, flyingAttackerID, "Bob", "Alice")
	
	// Normal blocker cannot block flying attacker
	canBlock, err := engine.CanBlock(gameID, normalBlockerID, flyingAttackerID)
	if err != nil {
		t.Fatalf("failed to check can block: %v", err)
	}
	if canBlock {
		t.Error("normal creature should not be able to block flying creature")
	}
	
	// Attempting to declare block should fail
	err = engine.DeclareBlocker(gameID, normalBlockerID, flyingAttackerID, "Bob")
	if err == nil {
		t.Error("declaring blocker should fail for flying attacker")
	}
}

// TestCombatFlyingBlockedByFlying tests that flying creatures can block flying creatures
func TestCombatFlyingBlockedByFlying(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-flying-vs-flying"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Flying attacker and flying blocker
	gameState.mu.Lock()
	flyingAttackerID := "flying-attacker"
	flyingBlockerID := "flying-blocker"
	
	gameState.cards[flyingAttackerID] = &internalCard{
		ID:           flyingAttackerID,
		Name:         "Wind Drake",
		Type:         "Creature - Drake",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityFlying, Text: "Flying"},
		},
	}
	
	gameState.cards[flyingBlockerID] = &internalCard{
		ID:           flyingBlockerID,
		Name:         "Storm Crow",
		Type:         "Creature - Bird",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "1",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityFlying, Text: "Flying"},
		},
	}
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, flyingAttackerID, "Bob", "Alice")
	
	// Flying blocker can block flying attacker
	canBlock, err := engine.CanBlock(gameID, flyingBlockerID, flyingAttackerID)
	if err != nil {
		t.Fatalf("failed to check can block: %v", err)
	}
	if !canBlock {
		t.Error("flying creature should be able to block flying creature")
	}
	
	// Declare block should succeed
	if err := engine.DeclareBlocker(gameID, flyingBlockerID, flyingAttackerID, "Bob"); err != nil {
		t.Fatalf("failed to declare blocker: %v", err)
	}
	
	// Verify block
	gameState.mu.RLock()
	blocker := gameState.cards[flyingBlockerID]
	if !blocker.Blocking {
		t.Error("flying creature should be blocking")
	}
	gameState.mu.RUnlock()
}

// TestCombatFlyingBlockedByReach tests that reach creatures can block flying creatures
func TestCombatFlyingBlockedByReach(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-flying-vs-reach"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Flying attacker and reach blocker
	gameState.mu.Lock()
	flyingAttackerID := "flying-attacker"
	reachBlockerID := "reach-blocker"
	
	gameState.cards[flyingAttackerID] = &internalCard{
		ID:           flyingAttackerID,
		Name:         "Wind Drake",
		Type:         "Creature - Drake",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityFlying, Text: "Flying"},
		},
	}
	
	gameState.cards[reachBlockerID] = &internalCard{
		ID:           reachBlockerID,
		Name:         "Giant Spider",
		Type:         "Creature - Spider",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "4",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityReach, Text: "Reach"},
		},
	}
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, flyingAttackerID, "Bob", "Alice")
	
	// Reach blocker can block flying attacker
	canBlock, err := engine.CanBlock(gameID, reachBlockerID, flyingAttackerID)
	if err != nil {
		t.Fatalf("failed to check can block: %v", err)
	}
	if !canBlock {
		t.Error("reach creature should be able to block flying creature")
	}
	
	// Declare block should succeed
	if err := engine.DeclareBlocker(gameID, reachBlockerID, flyingAttackerID, "Bob"); err != nil {
		t.Fatalf("failed to declare blocker: %v", err)
	}
	
	// Verify block
	gameState.mu.RLock()
	blocker := gameState.cards[reachBlockerID]
	if !blocker.Blocking {
		t.Error("reach creature should be blocking")
	}
	gameState.mu.RUnlock()
}

// TestCombatReachCanBlockNormal tests that reach doesn't prevent blocking normal creatures
func TestCombatReachCanBlockNormal(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-reach-vs-normal"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Normal attacker and reach blocker
	gameState.mu.Lock()
	normalAttackerID := "normal-attacker"
	reachBlockerID := "reach-blocker"
	
	gameState.cards[normalAttackerID] = &internalCard{
		ID:           normalAttackerID,
		Name:         "Grizzly Bears",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	
	gameState.cards[reachBlockerID] = &internalCard{
		ID:           reachBlockerID,
		Name:         "Giant Spider",
		Type:         "Creature - Spider",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "4",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityReach, Text: "Reach"},
		},
	}
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, normalAttackerID, "Bob", "Alice")
	
	// Reach blocker can block normal attacker
	canBlock, err := engine.CanBlock(gameID, reachBlockerID, normalAttackerID)
	if err != nil {
		t.Fatalf("failed to check can block: %v", err)
	}
	if !canBlock {
		t.Error("reach creature should be able to block normal creature")
	}
	
	// Declare block should succeed
	if err := engine.DeclareBlocker(gameID, reachBlockerID, normalAttackerID, "Bob"); err != nil {
		t.Fatalf("failed to declare blocker: %v", err)
	}
	
	// Verify block
	gameState.mu.RLock()
	blocker := gameState.cards[reachBlockerID]
	if !blocker.Blocking {
		t.Error("reach creature should be blocking")
	}
	gameState.mu.RUnlock()
}

// TestCombatFlyingFullFlow tests a complete combat with flying creatures
func TestCombatFlyingFullFlow(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-flying-full"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Flying attacker, normal blocker (can't block), reach blocker (can block)
	gameState.mu.Lock()
	flyingAttackerID := "flying-attacker"
	normalBlockerID := "normal-blocker"
	reachBlockerID := "reach-blocker"
	
	gameState.cards[flyingAttackerID] = &internalCard{
		ID:           flyingAttackerID,
		Name:         "Serra Angel",
		Type:         "Creature - Angel",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "4",
		Toughness:    "4",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityFlying, Text: "Flying"},
			{ID: abilityVigilance, Text: "Vigilance"},
		},
	}
	
	gameState.cards[normalBlockerID] = &internalCard{
		ID:           normalBlockerID,
		Name:         "Grizzly Bears",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	
	gameState.cards[reachBlockerID] = &internalCard{
		ID:           reachBlockerID,
		Name:         "Giant Spider",
		Type:         "Creature - Spider",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "4",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityReach, Text: "Reach"},
		},
	}
	
	initialBobLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, flyingAttackerID, "Bob", "Alice")
	
	// Normal blocker cannot block
	canBlockNormal, _ := engine.CanBlock(gameID, normalBlockerID, flyingAttackerID)
	if canBlockNormal {
		t.Error("normal creature should not be able to block flying creature")
	}
	
	// Reach blocker can block
	canBlockReach, _ := engine.CanBlock(gameID, reachBlockerID, flyingAttackerID)
	if !canBlockReach {
		t.Error("reach creature should be able to block flying creature")
	}
	
	// Block with reach creature
	engine.DeclareBlocker(gameID, reachBlockerID, flyingAttackerID, "Bob")
	engine.AcceptBlockers(gameID)
	
	// Combat damage
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	
	// Verify damage: 4/4 flying attacker vs 2/4 reach blocker
	// Attacker takes 2 damage (survives), blocker takes 4 damage (dies)
	gameState.mu.RLock()
	attacker := gameState.cards[flyingAttackerID]
	blocker := gameState.cards[reachBlockerID]
	
	if attacker.Zone != zoneBattlefield {
		t.Error("attacker should still be on battlefield")
	}
	if attacker.Damage != 2 {
		t.Errorf("attacker should have 2 damage, has %d", attacker.Damage)
	}
	
	if blocker.Zone != zoneGraveyard {
		t.Error("blocker should be in graveyard")
	}
	
	// Bob's life should be unchanged (attacker was blocked)
	bobLife := gameState.players["Bob"].Life
	if bobLife != initialBobLife {
		t.Errorf("Bob's life should be unchanged, was %d, now %d", initialBobLife, bobLife)
	}
	gameState.mu.RUnlock()
	
	// End combat
	engine.EndCombat(gameID)
	
	// Verify attacker still has vigilance (untapped)
	gameState.mu.RLock()
	attacker = gameState.cards[flyingAttackerID]
	if attacker.Tapped {
		t.Error("vigilance creature should still be untapped")
	}
	gameState.mu.RUnlock()
}

// TestCombatFlyingUnblocked tests flying creature dealing damage when unblocked
func TestCombatFlyingUnblocked(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-flying-unblocked"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Flying attacker only
	gameState.mu.Lock()
	flyingAttackerID := "flying-attacker"
	
	gameState.cards[flyingAttackerID] = &internalCard{
		ID:           flyingAttackerID,
		Name:         "Wind Drake",
		Type:         "Creature - Drake",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityFlying, Text: "Flying"},
		},
	}
	
	initialBobLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()
	
	// Full combat flow
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, flyingAttackerID, "Bob", "Alice")
	
	// No blockers declared
	engine.AcceptBlockers(gameID)
	
	// Damage (unblocked)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	
	// Verify damage dealt to Bob
	gameState.mu.RLock()
	bobLife := gameState.players["Bob"].Life
	if bobLife != initialBobLife-2 {
		t.Errorf("expected Bob to lose 2 life, lost %d", initialBobLife-bobLife)
	}
	gameState.mu.RUnlock()
	
	// End combat
	engine.EndCombat(gameID)
}

// TestCombatNormalCanBlockNormal tests that normal creatures can block each other
func TestCombatNormalCanBlockNormal(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-normal-vs-normal"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Normal attacker and normal blocker
	gameState.mu.Lock()
	normalAttackerID := "normal-attacker"
	normalBlockerID := "normal-blocker"
	
	gameState.cards[normalAttackerID] = &internalCard{
		ID:           normalAttackerID,
		Name:         "Grizzly Bears",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	
	gameState.cards[normalBlockerID] = &internalCard{
		ID:           normalBlockerID,
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
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, normalAttackerID, "Bob", "Alice")
	
	// Normal blocker can block normal attacker
	canBlock, err := engine.CanBlock(gameID, normalBlockerID, normalAttackerID)
	if err != nil {
		t.Fatalf("failed to check can block: %v", err)
	}
	if !canBlock {
		t.Error("normal creature should be able to block normal creature")
	}
	
	// Declare block should succeed
	if err := engine.DeclareBlocker(gameID, normalBlockerID, normalAttackerID, "Bob"); err != nil {
		t.Fatalf("failed to declare blocker: %v", err)
	}
	
	// Verify block
	gameState.mu.RLock()
	blocker := gameState.cards[normalBlockerID]
	if !blocker.Blocking {
		t.Error("normal creature should be blocking")
	}
	gameState.mu.RUnlock()
}
