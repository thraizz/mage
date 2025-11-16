package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestCheckBlockRequirementsMustBlock verifies must-block requirement detection
func TestCheckBlockRequirementsMustBlock(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-must-block"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create attacker and blocker
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
		Name:         "Blocker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
	}
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	
	// Set up must-block requirement (after combat setup)
	gameState.mu.Lock()
	gameState.combat.creatureMustBlockAttackers[blockerID] = map[string]bool{attackerID: true}
	gameState.mu.Unlock()
	
	// Check requirements without declaring blocker
	violations, err := engine.CheckBlockRequirements(gameID, "Bob")
	if err != nil {
		t.Fatalf("CheckBlockRequirements returned error: %v", err)
	}
	
	if len(violations) == 0 {
		t.Error("Expected violation for must-block requirement, got none")
	}
	
	// Now declare the blocker
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	
	// Check requirements again
	violations, err = engine.CheckBlockRequirements(gameID, "Bob")
	if err != nil {
		t.Fatalf("CheckBlockRequirements returned error: %v", err)
	}
	
	if len(violations) != 0 {
		t.Errorf("Expected no violations after blocking, got: %v", violations)
	}
	
	engine.EndCombat(gameID)
}

// TestCheckBlockRequirementsMinBlockers verifies minimum blocker requirement (menace)
func TestCheckBlockRequirementsMinBlockers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-min-blockers"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create attacker with menace and two blockers
	gameState.mu.Lock()
	attackerID := "attacker"
	blocker1ID := "blocker1"
	blocker2ID := "blocker2"
	
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
	
	gameState.cards[blocker1ID] = &internalCard{
		ID:           blocker1ID,
		Name:         "Blocker 1",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
	}
	
	gameState.cards[blocker2ID] = &internalCard{
		ID:           blocker2ID,
		Name:         "Blocker 2",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
	}
	
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	
	// Set up menace requirement (must be blocked by 2+ creatures) after combat setup
	gameState.mu.Lock()
	gameState.combat.minBlockersPerAttacker[attackerID] = 2
	gameState.mu.Unlock()
	
	// Declare only one blocker (violates menace)
	engine.DeclareBlocker(gameID, blocker1ID, attackerID, "Bob")
	
	// Check requirements
	violations, err := engine.CheckBlockRequirements(gameID, "Bob")
	if err != nil {
		t.Fatalf("CheckBlockRequirements returned error: %v", err)
	}
	
	if len(violations) == 0 {
		t.Error("Expected violation for menace requirement, got none")
	}
	
	// Add second blocker
	engine.DeclareBlocker(gameID, blocker2ID, attackerID, "Bob")
	
	// Check requirements again
	violations, err = engine.CheckBlockRequirements(gameID, "Bob")
	if err != nil {
		t.Fatalf("CheckBlockRequirements returned error: %v", err)
	}
	
	if len(violations) != 0 {
		t.Errorf("Expected no violations with 2 blockers, got: %v", violations)
	}
	
	engine.EndCombat(gameID)
}

// TestCheckBlockRestrictionsMaxBlockers verifies maximum blocker restriction
func TestCheckBlockRestrictionsMaxBlockers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-max-blockers"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create attacker and three blockers
	gameState.mu.Lock()
	attackerID := "attacker"
	blocker1ID := "blocker1"
	blocker2ID := "blocker2"
	blocker3ID := "blocker3"
	
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
	
	gameState.cards[blocker1ID] = &internalCard{
		ID:           blocker1ID,
		Name:         "Blocker 1",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
	}
	
	gameState.cards[blocker2ID] = &internalCard{
		ID:           blocker2ID,
		Name:         "Blocker 2",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
	}
	
	gameState.cards[blocker3ID] = &internalCard{
		ID:           blocker3ID,
		Name:         "Blocker 3",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
	}
	
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	
	// Set up max blockers restriction (can only be blocked by 1 creature) after combat setup
	gameState.mu.Lock()
	gameState.combat.maxBlockersPerAttacker[attackerID] = 1
	gameState.mu.Unlock()
	
	// Declare two blockers (violates restriction)
	engine.DeclareBlocker(gameID, blocker1ID, attackerID, "Bob")
	engine.DeclareBlocker(gameID, blocker2ID, attackerID, "Bob")
	
	// Check restrictions
	violations, err := engine.CheckBlockRestrictions(gameID, "Bob")
	if err != nil {
		t.Fatalf("CheckBlockRestrictions returned error: %v", err)
	}
	
	if len(violations) == 0 {
		t.Error("Expected violation for max blockers restriction, got none")
	}
	
	engine.EndCombat(gameID)
}

// TestValidateAttackerCountMax verifies maximum attacker count validation
func TestValidateAttackerCountMax(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-max-attackers"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create three attackers
	gameState.mu.Lock()
	attacker1ID := "attacker1"
	attacker2ID := "attacker2"
	attacker3ID := "attacker3"
	
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
	
	gameState.cards[attacker3ID] = &internalCard{
		ID:           attacker3ID,
		Name:         "Attacker 3",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
	}
	
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	
	// Set up max attackers restriction after combat setup
	gameState.mu.Lock()
	gameState.combat.maxAttackers = 2
	gameState.mu.Unlock()
	
	// Declare three attackers (violates restriction)
	engine.DeclareAttacker(gameID, attacker1ID, "Bob", "Alice")
	engine.DeclareAttacker(gameID, attacker2ID, "Bob", "Alice")
	engine.DeclareAttacker(gameID, attacker3ID, "Bob", "Alice")
	
	// Validate attacker count
	violations, err := engine.ValidateAttackerCount(gameID)
	if err != nil {
		t.Fatalf("ValidateAttackerCount returned error: %v", err)
	}
	
	if len(violations) == 0 {
		t.Error("Expected violation for max attackers, got none")
	}
	
	engine.EndCombat(gameID)
}

// TestValidateAttackerCountForcedAttack verifies forced attack validation
func TestValidateAttackerCountForcedAttack(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-forced-attack"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Create creature that must attack
	gameState.mu.Lock()
	creatureID := "forced-attacker"
	
	gameState.cards[creatureID] = &internalCard{
		ID:           creatureID,
		Name:         "Forced Attacker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
	}
	
	gameState.mu.Unlock()
	
	// Setup combat without declaring the forced attacker
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	
	// Set up forced attack (must attack any defender) after combat setup
	gameState.mu.Lock()
	gameState.combat.creaturesForcedToAttack[creatureID] = make(map[string]bool)
	gameState.mu.Unlock()
	
	// Validate attacker count
	violations, err := engine.ValidateAttackerCount(gameID)
	if err != nil {
		t.Fatalf("ValidateAttackerCount returned error: %v", err)
	}
	
	if len(violations) == 0 {
		t.Error("Expected violation for forced attack, got none")
	}
	
	// Now declare the attacker
	engine.DeclareAttacker(gameID, creatureID, "Bob", "Alice")
	
	// Validate again
	violations, err = engine.ValidateAttackerCount(gameID)
	if err != nil {
		t.Fatalf("ValidateAttackerCount returned error: %v", err)
	}
	
	if len(violations) != 0 {
		t.Errorf("Expected no violations after declaring forced attacker, got: %v", violations)
	}
	
	engine.EndCombat(gameID)
}
