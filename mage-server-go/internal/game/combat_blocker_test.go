package game

import (
	"sync"
	"testing"
	"time"

	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap/zaptest"
)

// TestCombatBlockerDeclaration tests basic blocker declaration
func TestCombatBlockerDeclaration(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-game-blocker"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create two creatures - one for Alice (attacker), one for Bob (blocker)
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

	// Setup combat
	if err := engine.ResetCombat(gameID); err != nil {
		t.Fatalf("failed to reset combat: %v", err)
	}

	if err := engine.SetAttacker(gameID, "Alice"); err != nil {
		t.Fatalf("failed to set attacker: %v", err)
	}

	if err := engine.SetDefenders(gameID); err != nil {
		t.Fatalf("failed to set defenders: %v", err)
	}

	// Declare attacker
	if err := engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice"); err != nil {
		t.Fatalf("failed to declare attacker: %v", err)
	}

	// Test CanBlock
	canBlock, err := engine.CanBlock(gameID, blockerID, attackerID)
	if err != nil {
		t.Fatalf("CanBlock failed: %v", err)
	}
	if !canBlock {
		t.Fatal("blocker should be able to block attacker")
	}

	// Declare blocker
	if err := engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob"); err != nil {
		t.Fatalf("failed to declare blocker: %v", err)
	}

	// Verify blocker state
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	if !blocker.Blocking {
		t.Error("blocker should be marked as blocking")
	}
	if len(blocker.BlockingWhat) != 1 || blocker.BlockingWhat[0] != attackerID {
		t.Errorf("blocker should be blocking attacker, got: %v", blocker.BlockingWhat)
	}

	// Verify combat group
	if len(gameState.combat.groups) != 1 {
		t.Fatalf("expected 1 combat group, got %d", len(gameState.combat.groups))
	}
	group := gameState.combat.groups[0]
	if !group.blocked {
		t.Error("combat group should be marked as blocked")
	}
	if len(group.blockers) != 1 || group.blockers[0] != blockerID {
		t.Errorf("combat group should have blocker, got: %v", group.blockers)
	}
	gameState.mu.RUnlock()
}

// TestCombatMultipleBlockers tests multiple blockers on one attacker
func TestCombatMultipleBlockers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-game-multi-blockers"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create one attacker and two blockers
	gameState.mu.Lock()
	attackerID := "attacker-1"
	blocker1ID := "blocker-1"
	blocker2ID := "blocker-2"

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
		Name:         "Llanowar Elves",
		Type:         "Creature - Elf",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "1",
		Toughness:    "1",
		Tapped:       false,
	}
	gameState.mu.Unlock()

	// Setup combat and declare attacker
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Declare both blockers
	if err := engine.DeclareBlocker(gameID, blocker1ID, attackerID, "Bob"); err != nil {
		t.Fatalf("failed to declare first blocker: %v", err)
	}

	if err := engine.DeclareBlocker(gameID, blocker2ID, attackerID, "Bob"); err != nil {
		t.Fatalf("failed to declare second blocker: %v", err)
	}

	// Verify both blockers are in the combat group
	gameState.mu.RLock()
	group := gameState.combat.groups[0]
	if len(group.blockers) != 2 {
		t.Errorf("expected 2 blockers, got %d", len(group.blockers))
	}
	if !group.blocked {
		t.Error("combat group should be marked as blocked")
	}
	gameState.mu.RUnlock()
}

// TestCombatBlockerValidation tests blocker validation rules
func TestCombatBlockerValidation(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-game-blocker-validation"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create attacker and various blockers
	gameState.mu.Lock()
	attackerID := "attacker-1"
	tappedBlockerID := "tapped-blocker"
	wrongControllerBlockerID := "wrong-controller-blocker"
	nonCreatureID := "non-creature"

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

	gameState.cards[tappedBlockerID] = &internalCard{
		ID:           tappedBlockerID,
		Name:         "Tapped Wall",
		Type:         "Creature - Wall",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "0",
		Toughness:    "5",
		Tapped:       true, // Tapped!
	}

	gameState.cards[wrongControllerBlockerID] = &internalCard{
		ID:           wrongControllerBlockerID,
		Name:         "Alice's Creature",
		Type:         "Creature - Human",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice", // Controlled by attacker!
		Power:        "1",
		Toughness:    "1",
		Tapped:       false,
	}

	gameState.cards[nonCreatureID] = &internalCard{
		ID:           nonCreatureID,
		Name:         "Lightning Bolt",
		Type:         "Instant",
		Zone:         zoneBattlefield, // Shouldn't be on battlefield, but testing
		OwnerID:      "Bob",
		ControllerID: "Bob",
	}
	gameState.mu.Unlock()

	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Test 1: Tapped creature can't block
	canBlock, err := engine.CanBlock(gameID, tappedBlockerID, attackerID)
	if err != nil {
		t.Fatalf("CanBlock failed for tapped blocker: %v", err)
	}
	if canBlock {
		t.Error("tapped creature should not be able to block")
	}

	err = engine.DeclareBlocker(gameID, tappedBlockerID, attackerID, "Bob")
	if err == nil {
		t.Error("should not be able to declare tapped creature as blocker")
	}

	// Test 2: Wrong controller can't block
	canBlock, err = engine.CanBlock(gameID, wrongControllerBlockerID, attackerID)
	if err != nil {
		t.Fatalf("CanBlock failed for wrong controller: %v", err)
	}
	if canBlock {
		t.Error("creature controlled by attacker should not be able to block")
	}

	// Test 3: Non-creature can't block
	canBlock, err = engine.CanBlock(gameID, nonCreatureID, attackerID)
	if err != nil {
		t.Fatalf("CanBlock failed for non-creature: %v", err)
	}
	if canBlock {
		t.Error("non-creature should not be able to block")
	}
}

// TestCombatRemoveBlocker tests removing a blocker
func TestCombatRemoveBlocker(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-game-remove-blocker"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create attacker and blocker
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

	// Setup combat and declare blocker
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")

	// Verify blocker is blocking
	gameState.mu.RLock()
	blocker := gameState.cards[blockerID]
	if !blocker.Blocking {
		t.Fatal("blocker should be blocking before removal")
	}
	group := gameState.combat.groups[0]
	if !group.blocked {
		t.Fatal("group should be blocked before removal")
	}
	gameState.mu.RUnlock()

	// Remove blocker
	if err := engine.RemoveBlocker(gameID, blockerID); err != nil {
		t.Fatalf("failed to remove blocker: %v", err)
	}

	// Verify blocker is no longer blocking
	gameState.mu.RLock()
	blocker = gameState.cards[blockerID]
	if blocker.Blocking {
		t.Error("blocker should not be blocking after removal")
	}
	if blocker.BlockingWhat != nil && len(blocker.BlockingWhat) > 0 {
		t.Error("blocker should not be blocking anything after removal")
	}

	// Verify group is no longer blocked
	group = gameState.combat.groups[0]
	if group.blocked {
		t.Error("group should not be blocked after blocker removal")
	}
	if len(group.blockers) != 0 {
		t.Errorf("group should have no blockers, got %d", len(group.blockers))
	}
	gameState.mu.RUnlock()
}

// TestCombatAcceptBlockers tests the AcceptBlockers method and events
func TestCombatAcceptBlockers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-game-accept-blockers"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Create attacker and blocker
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

	// Setup event tracking
	var events []rules.Event
	eventsMu := sync.Mutex{}

	gameState.eventBus.SubscribeTyped(rules.EventBlockerDeclared, func(event rules.Event) {
		eventsMu.Lock()
		events = append(events, event)
		eventsMu.Unlock()
	})

	gameState.eventBus.SubscribeTyped(rules.EventCreatureBlocked, func(event rules.Event) {
		eventsMu.Lock()
		events = append(events, event)
		eventsMu.Unlock()
	})

	gameState.eventBus.SubscribeTyped(rules.EventCreatureBlocks, func(event rules.Event) {
		eventsMu.Lock()
		events = append(events, event)
		eventsMu.Unlock()
	})

	gameState.eventBus.SubscribeTyped(rules.EventDeclaredBlockers, func(event rules.Event) {
		eventsMu.Lock()
		events = append(events, event)
		eventsMu.Unlock()
	})

	// Setup combat and declare blocker
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")

	// Clear events from attacker declaration
	eventsMu.Lock()
	events = nil
	eventsMu.Unlock()

	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")

	// Accept blockers
	if err := engine.AcceptBlockers(gameID); err != nil {
		t.Fatalf("failed to accept blockers: %v", err)
	}

	// Give events time to propagate
	time.Sleep(10 * time.Millisecond)

	// Verify events were fired
	eventsMu.Lock()
	defer eventsMu.Unlock()

	// Should have:
	// - BLOCKER_DECLARED (from DeclareBlocker)
	// - BLOCKER_DECLARED (from AcceptBlockers)
	// - CREATURE_BLOCKED
	// - CREATURE_BLOCKS
	// - DECLARED_BLOCKERS

	if len(events) < 4 {
		t.Errorf("expected at least 4 events, got %d", len(events))
	}

	// Check for specific event types
	hasBlockerDeclared := false
	hasCreatureBlocked := false
	hasCreatureBlocks := false
	hasDeclaredBlockers := false

	for _, event := range events {
		switch event.Type {
		case rules.EventBlockerDeclared:
			hasBlockerDeclared = true
		case rules.EventCreatureBlocked:
			hasCreatureBlocked = true
		case rules.EventCreatureBlocks:
			hasCreatureBlocks = true
		case rules.EventDeclaredBlockers:
			hasDeclaredBlockers = true
		}
	}

	if !hasBlockerDeclared {
		t.Error("missing BLOCKER_DECLARED event")
	}
	if !hasCreatureBlocked {
		t.Error("missing CREATURE_BLOCKED event")
	}
	if !hasCreatureBlocks {
		t.Error("missing CREATURE_BLOCKS event")
	}
	if !hasDeclaredBlockers {
		t.Error("missing DECLARED_BLOCKERS event")
	}
}
