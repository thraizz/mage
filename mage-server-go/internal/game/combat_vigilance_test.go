package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestCombatVigilance tests that vigilance creatures don't tap when attacking
func TestCombatVigilance(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-vigilance"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Creature with vigilance
	gameState.mu.Lock()
	attackerID := "attacker-1"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Vigilant Knight",
		Type:         "Creature - Knight",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityVigilance, Text: "Vigilance"},
		},
	}
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	
	// Declare attacker
	if err := engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice"); err != nil {
		t.Fatalf("failed to declare attacker: %v", err)
	}
	
	// Verify creature is attacking but NOT tapped
	gameState.mu.RLock()
	attacker := gameState.cards[attackerID]
	
	if !attacker.Attacking {
		t.Error("creature should be attacking")
	}
	if attacker.Tapped {
		t.Error("vigilance creature should not be tapped when attacking")
	}
	
	// Verify not tracked as tapped by attack
	if gameState.combat.attackersTapped[attackerID] {
		t.Error("vigilance creature should not be in attackersTapped map")
	}
	gameState.mu.RUnlock()
}

// TestCombatNoVigilance tests that normal creatures tap when attacking
func TestCombatNoVigilance(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-no-vigilance"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Normal creature without vigilance
	gameState.mu.Lock()
	attackerID := "attacker-1"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Normal Knight",
		Type:         "Creature - Knight",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Setup combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	
	// Declare attacker
	if err := engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice"); err != nil {
		t.Fatalf("failed to declare attacker: %v", err)
	}
	
	// Verify creature is attacking AND tapped
	gameState.mu.RLock()
	attacker := gameState.cards[attackerID]
	
	if !attacker.Attacking {
		t.Error("creature should be attacking")
	}
	if !attacker.Tapped {
		t.Error("normal creature should be tapped when attacking")
	}
	
	// Verify tracked as tapped by attack
	if !gameState.combat.attackersTapped[attackerID] {
		t.Error("normal creature should be in attackersTapped map")
	}
	gameState.mu.RUnlock()
}

// TestCombatVigilanceCanBlock tests that vigilance creatures can block after attacking
func TestCombatVigilanceCanBlock(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-vigilance-block"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Vigilance creature and opponent's attacker
	gameState.mu.Lock()
	vigilantID := "vigilant-1"
	opponentAttackerID := "opponent-attacker"
	
	gameState.cards[vigilantID] = &internalCard{
		ID:           vigilantID,
		Name:         "Vigilant Knight",
		Type:         "Creature - Knight",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityVigilance, Text: "Vigilance"},
		},
	}
	
	gameState.cards[opponentAttackerID] = &internalCard{
		ID:           opponentAttackerID,
		Name:         "Opponent Bear",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Alice's turn: Attack with vigilance creature
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, vigilantID, "Bob", "Alice")
	
	// Verify vigilant creature is untapped
	gameState.mu.RLock()
	vigilant := gameState.cards[vigilantID]
	if vigilant.Tapped {
		t.Error("vigilance creature should be untapped after attacking")
	}
	gameState.mu.RUnlock()
	
	// End Alice's combat
	engine.EndCombat(gameID)
	
	// Bob's turn: Attack with opponent creature
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Bob")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, opponentAttackerID, "Alice", "Bob")
	
	// Alice can block with vigilant creature (it's untapped)
	canBlock, err := engine.CanBlock(gameID, vigilantID, opponentAttackerID)
	if err != nil {
		t.Fatalf("failed to check can block: %v", err)
	}
	if !canBlock {
		t.Error("vigilance creature should be able to block (it's untapped)")
	}
	
	// Actually declare the block
	if err := engine.DeclareBlocker(gameID, vigilantID, opponentAttackerID, "Alice"); err != nil {
		t.Fatalf("failed to declare blocker: %v", err)
	}
	
	// Verify block was successful
	gameState.mu.RLock()
	vigilant = gameState.cards[vigilantID]
	if !vigilant.Blocking {
		t.Error("vigilance creature should be blocking")
	}
	gameState.mu.RUnlock()
}

// TestCombatNoVigilanceCannotBlock tests that tapped creatures cannot block
func TestCombatNoVigilanceCannotBlock(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-no-vigilance-block"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Normal creature and opponent's attacker
	gameState.mu.Lock()
	normalID := "normal-1"
	opponentAttackerID := "opponent-attacker"
	
	gameState.cards[normalID] = &internalCard{
		ID:           normalID,
		Name:         "Normal Knight",
		Type:         "Creature - Knight",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	
	gameState.cards[opponentAttackerID] = &internalCard{
		ID:           opponentAttackerID,
		Name:         "Opponent Bear",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Alice's turn: Attack with normal creature
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, normalID, "Bob", "Alice")
	
	// Verify normal creature is tapped
	gameState.mu.RLock()
	normal := gameState.cards[normalID]
	if !normal.Tapped {
		t.Error("normal creature should be tapped after attacking")
	}
	gameState.mu.RUnlock()
	
	// End Alice's combat
	engine.EndCombat(gameID)
	
	// Bob's turn: Attack with opponent creature
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Bob")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, opponentAttackerID, "Alice", "Bob")
	
	// Alice cannot block with tapped creature
	canBlock, err := engine.CanBlock(gameID, normalID, opponentAttackerID)
	if err != nil {
		t.Fatalf("failed to check can block: %v", err)
	}
	if canBlock {
		t.Error("tapped creature should not be able to block")
	}
}

// TestCombatVigilanceFullFlow tests a complete combat with vigilance
func TestCombatVigilanceFullFlow(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-vigilance-full"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Vigilance attacker
	gameState.mu.Lock()
	attackerID := "attacker-1"
	
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
		Abilities: []EngineAbilityView{
			{ID: abilityVigilance, Text: "Vigilance"},
			{ID: abilityFlying, Text: "Flying"},
		},
	}
	
	initialBobLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()
	
	// Full combat flow
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	
	// Verify untapped
	gameState.mu.RLock()
	if gameState.cards[attackerID].Tapped {
		t.Error("vigilance creature should be untapped")
	}
	gameState.mu.RUnlock()
	
	// Damage (unblocked)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	
	// Verify damage dealt
	gameState.mu.RLock()
	bobLife := gameState.players["Bob"].Life
	if bobLife != initialBobLife-4 {
		t.Errorf("expected Bob to lose 4 life, lost %d", initialBobLife-bobLife)
	}
	gameState.mu.RUnlock()
	
	// End combat
	engine.EndCombat(gameID)
	
	// Verify still untapped after combat
	gameState.mu.RLock()
	attacker := gameState.cards[attackerID]
	if attacker.Tapped {
		t.Error("vigilance creature should still be untapped after combat")
	}
	if attacker.Attacking {
		t.Error("creature should not be attacking after combat ends")
	}
	gameState.mu.RUnlock()
}
