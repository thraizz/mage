package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestCombatUnblockableCannotBeBlocked tests that unblockable creatures cannot be blocked
func TestCombatUnblockableCannotBeBlocked(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-unblockable"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Unblockable attacker and normal blocker
	gameState.mu.Lock()
	attackerID := "unblockable-attacker"
	blockerID := "normal-blocker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Invisible Stalker",
		Type:         "Creature - Human Rogue",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "1",
		Toughness:    "1",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityUnblockable, Text: "This creature can't be blocked"},
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
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Normal blocker cannot block unblockable attacker
	canBlock, err := engine.CanBlock(gameID, blockerID, attackerID)
	if err != nil {
		t.Fatalf("failed to check can block: %v", err)
	}
	if canBlock {
		t.Error("normal creature should not be able to block unblockable creature")
	}

	// Attempting to declare block should fail
	err = engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	if err == nil {
		t.Error("declaring blocker should fail for unblockable attacker")
	}
}

// TestCombatUnblockableDealsFullDamage tests that unblockable creatures deal full damage to player
func TestCombatUnblockableDealsFullDamage(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-unblockable-damage"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Unblockable attacker
	gameState.mu.Lock()
	attackerID := "unblockable-attacker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Phantom Warrior",
		Type:         "Creature - Illusion Warrior",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityUnblockable, Text: "This creature can't be blocked"},
		},
	}

	initialBobLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()

	// Full combat flow
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// No blockers can be declared
	engine.AcceptBlockers(gameID)

	// Damage
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

// TestCombatUnblockableVsFlyingBlocker tests that even flying/reach blockers can't block unblockable
func TestCombatUnblockableVsFlyingBlocker(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-unblockable-vs-flying"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Unblockable attacker and flying blocker
	gameState.mu.Lock()
	attackerID := "unblockable-attacker"
	blockerID := "flying-blocker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Invisible Stalker",
		Type:         "Creature - Human Rogue",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "1",
		Toughness:    "1",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityUnblockable, Text: "This creature can't be blocked"},
		},
	}

	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Wind Drake",
		Type:         "Creature - Drake",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
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
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Flying blocker cannot block unblockable attacker
	canBlock, err := engine.CanBlock(gameID, blockerID, attackerID)
	if err != nil {
		t.Fatalf("failed to check can block: %v", err)
	}
	if canBlock {
		t.Error("flying creature should not be able to block unblockable creature")
	}
}

// TestCombatUnblockableWithOtherAbilities tests unblockable combined with other combat abilities
func TestCombatUnblockableWithOtherAbilities(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-unblockable-combo"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Unblockable + lifelink + first strike attacker
	gameState.mu.Lock()
	attackerID := "unblockable-attacker"
	blockerID := "blocker"

	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Elite Phantom",
		Type:         "Creature - Spirit Soldier",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "3",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityUnblockable, Text: "This creature can't be blocked"},
			{ID: abilityLifelink, Text: "Lifelink"},
			{ID: abilityFirstStrike, Text: "First strike"},
		},
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

	initialBobLife := gameState.players["Bob"].Life
	initialAliceLife := gameState.players["Alice"].Life
	gameState.mu.Unlock()

	// Combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Cannot block
	canBlock, _ := engine.CanBlock(gameID, blockerID, attackerID)
	if canBlock {
		t.Error("wall should not be able to block unblockable creature")
	}

	// No blockers declared
	engine.AcceptBlockers(gameID)

	// First strike damage
	engine.AssignCombatDamage(gameID, true)
	engine.ApplyCombatDamage(gameID)

	// Verify damage and lifelink
	gameState.mu.RLock()
	bobLife := gameState.players["Bob"].Life
	aliceLife := gameState.players["Alice"].Life

	if bobLife != initialBobLife-3 {
		t.Errorf("expected Bob to lose 3 life, lost %d", initialBobLife-bobLife)
	}

	if aliceLife != initialAliceLife+3 {
		t.Errorf("expected Alice to gain 3 life from lifelink, gained %d", aliceLife-initialAliceLife)
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}

// TestCombatMultipleUnblockableAttackers tests multiple unblockable creatures attacking
func TestCombatMultipleUnblockableAttackers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-multi-unblockable"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Two unblockable attackers
	gameState.mu.Lock()
	attacker1ID := "unblockable-1"
	attacker2ID := "unblockable-2"
	blockerID := "blocker"

	gameState.cards[attacker1ID] = &internalCard{
		ID:           attacker1ID,
		Name:         "Invisible Stalker",
		Type:         "Creature - Human Rogue",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "1",
		Toughness:    "1",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityUnblockable, Text: "This creature can't be blocked"},
		},
	}

	gameState.cards[attacker2ID] = &internalCard{
		ID:           attacker2ID,
		Name:         "Phantom Warrior",
		Type:         "Creature - Illusion Warrior",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityUnblockable, Text: "This creature can't be blocked"},
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
	engine.DeclareAttacker(gameID, attacker1ID, "Bob", "Alice")
	engine.DeclareAttacker(gameID, attacker2ID, "Bob", "Alice")

	// Cannot block either attacker
	canBlock1, _ := engine.CanBlock(gameID, blockerID, attacker1ID)
	canBlock2, _ := engine.CanBlock(gameID, blockerID, attacker2ID)

	if canBlock1 || canBlock2 {
		t.Error("blocker should not be able to block either unblockable creature")
	}

	// No blockers
	engine.AcceptBlockers(gameID)

	// Damage
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify total damage (1 + 2 = 3)
	gameState.mu.RLock()
	bobLife := gameState.players["Bob"].Life
	if bobLife != initialBobLife-3 {
		t.Errorf("expected Bob to lose 3 life, lost %d", initialBobLife-bobLife)
	}
	gameState.mu.RUnlock()

	engine.EndCombat(gameID)
}

// TestCombatNormalCreatureCanBeBlocked tests that normal creatures can still be blocked
func TestCombatNormalCreatureCanBeBlocked(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-normal-blockable"
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
	attackerID := "normal-attacker"
	blockerID := "normal-blocker"

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

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Normal blocker can block normal attacker
	canBlock, err := engine.CanBlock(gameID, blockerID, attackerID)
	if err != nil {
		t.Fatalf("failed to check can block: %v", err)
	}
	if !canBlock {
		t.Error("normal creature should be able to block normal creature")
	}

	// Declare block should succeed
	if err := engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob"); err != nil {
		t.Fatalf("failed to declare blocker: %v", err)
	}

	// Verify block
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	if !blocker.Blocking {
		t.Error("normal creature should be blocking")
	}
	gameState.mu.RUnlock()
}
