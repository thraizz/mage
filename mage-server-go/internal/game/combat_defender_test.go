package game

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestCombatDefenderCannotAttack verifies that creatures with defender cannot attack
func TestCombatDefenderCannotAttack(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-defender-cannot-attack"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Creature with defender
	gameState.mu.Lock()
	defenderCreatureID := "wall-of-stone"

	gameState.cards[defenderCreatureID] = &internalCard{
		ID:           defenderCreatureID,
		Name:         "Wall of Stone",
		Type:         "Creature - Wall",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "0",
		Toughness:    "7",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityDefender, Text: "Defender"},
		},
	}
	gameState.mu.Unlock()

	// Try to attack with defender creature
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	err := engine.DeclareAttacker(gameID, defenderCreatureID, "Bob", "Alice")

	// Should fail
	if err == nil {
		t.Error("Expected error when declaring attacker with defender, got nil")
	}

	// Verify creature is not attacking
	gameState.mu.RLock()
	creature := gameState.cards[defenderCreatureID]
	isAttacking := gameState.combat.attackers[defenderCreatureID]
	gameState.mu.RUnlock()

	if creature.Attacking {
		t.Error("Creature with defender should not be attacking")
	}

	if isAttacking {
		t.Error("Creature with defender should not be in attackers set")
	}

	engine.EndCombat(gameID)
}

// TestCombatDefenderCanBlock verifies that creatures with defender can still block
func TestCombatDefenderCanBlock(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-defender-can-block"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Alice has attacker, Bob has defender
	gameState.mu.Lock()
	attackerID := "grizzly-bears"
	defenderCreatureID := "wall-of-stone"

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

	gameState.cards[defenderCreatureID] = &internalCard{
		ID:           defenderCreatureID,
		Name:         "Wall of Stone",
		Type:         "Creature - Wall",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "0",
		Toughness:    "7",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityDefender, Text: "Defender"},
		},
	}
	gameState.mu.Unlock()

	// Combat: Alice attacks, Bob blocks with defender
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	if err := engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice"); err != nil {
		t.Fatalf("Failed to declare attacker: %v", err)
	}

	// Defender should be able to block
	if err := engine.DeclareBlocker(gameID, defenderCreatureID, attackerID, "Bob"); err != nil {
		t.Fatalf("Creature with defender should be able to block, got error: %v", err)
	}

	// Verify blocking
	gameState.mu.RLock()
	blocker := gameState.cards[defenderCreatureID]
	isBlocking := gameState.combat.blockers[defenderCreatureID]
	gameState.mu.RUnlock()

	if !blocker.Blocking {
		t.Error("Defender creature should be blocking")
	}

	if !isBlocking {
		t.Error("Defender creature should be in blockers set")
	}

	engine.AcceptBlockers(gameID)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify combat damage was applied
	gameState.mu.RLock()
	blockerCard := gameState.cards[defenderCreatureID]
	attackerCard := gameState.cards[attackerID]
	gameState.mu.RUnlock()

	if blockerCard.Damage != 2 {
		t.Errorf("Expected defender to have 2 damage, got %d", blockerCard.Damage)
	}

	// Attacker should have 0 damage (wall has 0 power)
	if attackerCard.Damage != 0 {
		t.Errorf("Expected attacker to have 0 damage, got %d", attackerCard.Damage)
	}

	engine.EndCombat(gameID)
}

// TestCombatNormalCreatureCanAttack verifies normal creatures without defender can attack
func TestCombatNormalCreatureCanAttack(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-normal-can-attack"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Normal creature without defender
	gameState.mu.Lock()
	creatureID := "grizzly-bears"

	gameState.cards[creatureID] = &internalCard{
		ID:           creatureID,
		Name:         "Grizzly Bears",
		Type:         "Creature - Bear",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	gameState.mu.Unlock()

	// Should be able to attack
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	if err := engine.DeclareAttacker(gameID, creatureID, "Bob", "Alice"); err != nil {
		t.Fatalf("Normal creature should be able to attack, got error: %v", err)
	}

	// Verify attacking
	gameState.mu.RLock()
	creature := gameState.cards[creatureID]
	isAttacking := gameState.combat.attackers[creatureID]
	gameState.mu.RUnlock()

	if !creature.Attacking {
		t.Error("Normal creature should be attacking")
	}

	if !isAttacking {
		t.Error("Normal creature should be in attackers set")
	}

	engine.EndCombat(gameID)
}

// TestCombatDefenderFullFlow verifies defender in a complete combat scenario
func TestCombatDefenderFullFlow(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)

	gameID := "test-defender-full-flow"
	players := []string{"Alice", "Bob"}

	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()

	// Setup: Alice has normal attacker, Bob has defender and normal blocker
	gameState.mu.Lock()
	attackerID := "serra-angel"
	defenderWallID := "wall-of-stone"
	normalBlockerID := "grizzly-bears"

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

	gameState.cards[defenderWallID] = &internalCard{
		ID:           defenderWallID,
		Name:         "Wall of Stone",
		Type:         "Creature - Wall",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "0",
		Toughness:    "7",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityDefender, Text: "Defender"},
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

	initialBobLife := gameState.players["Bob"].Life
	gameState.mu.Unlock()

	// Combat: Alice attacks
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)

	// Alice attacks with Serra Angel
	if err := engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice"); err != nil {
		t.Fatalf("Failed to declare attacker: %v", err)
	}

	// Bob tries to attack with wall (should fail)
	engine.SetAttacker(gameID, "Bob")
	if err := engine.DeclareAttacker(gameID, defenderWallID, "Alice", "Bob"); err == nil {
		t.Error("Wall with defender should not be able to attack")
	}

	// Back to Alice's attack, Bob blocks with wall
	engine.SetAttacker(gameID, "Alice")
	if err := engine.DeclareBlocker(gameID, defenderWallID, attackerID, "Bob"); err != nil {
		t.Fatalf("Wall should be able to block: %v", err)
	}

	engine.AcceptBlockers(gameID)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)

	// Verify: Wall takes 4 damage, attacker takes 0, Bob takes no damage
	gameState.mu.RLock()
	wall := gameState.cards[defenderWallID]
	attacker := gameState.cards[attackerID]
	bobLife := gameState.players["Bob"].Life
	gameState.mu.RUnlock()

	if wall.Damage != 4 {
		t.Errorf("Expected wall to have 4 damage, got %d", wall.Damage)
	}

	if attacker.Damage != 0 {
		t.Errorf("Expected attacker to have 0 damage (wall has 0 power), got %d", attacker.Damage)
	}

	if bobLife != initialBobLife {
		t.Errorf("Expected Bob to have %d life (blocked), got %d", initialBobLife, bobLife)
	}

	engine.EndCombat(gameID)
}
