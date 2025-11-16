package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/effects"
	"go.uber.org/zap/zaptest"
)

// TestCombatDynamicFlying verifies granted flying works in combat
func TestCombatDynamicFlying(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-dynamic-flying"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Alice has a normal creature, Bob has a flying blocker
	gameState.mu.Lock()
	attackerID := "grizzly-bears"
	blockerID := "flying-blocker"
	
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
		Abilities:    []EngineAbilityView{}, // No flying initially
	}
	
	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Flying Blocker",
		Type:         "Creature - Bird",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "1",
		Toughness:    "1",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityFlying, Text: "Flying"},
		},
	}
	
	// Grant flying to attacker dynamically
	effect := effects.NewEffectBuilder("giant-growth").
		Targeting(attackerID).
		UntilEndOfCombat().
		GrantAbility(abilityFlying)
	
	effectManager := effects.NewEffectManager(gameState.layerSystem)
	effectManager.AddEffect(effect)
	
	gameState.mu.Unlock()
	
	// Verify attacker now has flying (including granted)
	gameState.mu.RLock()
	attacker := gameState.cards[attackerID]
	hasFlying := engine.hasAbilityWithEffects(gameState, attacker, abilityFlying)
	gameState.mu.RUnlock()
	
	if !hasFlying {
		t.Fatal("Attacker should have flying after effect")
	}
	
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
	
	// Flying blocker should be able to block flying attacker
	canBlock, err := engine.CanBlock(gameID, blockerID, attackerID)
	if err != nil {
		t.Fatalf("failed to check can block: %v", err)
	}
	
	if !canBlock {
		t.Error("Flying blocker should be able to block flying attacker")
	}
}

// TestCombatDynamicVigilance verifies granted vigilance works in combat
func TestCombatDynamicVigilance(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-dynamic-vigilance"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Alice has a creature without vigilance
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
		Abilities:    []EngineAbilityView{}, // No vigilance initially
	}
	
	// Grant vigilance dynamically
	effect := effects.NewEffectBuilder("inspired-charge").
		Targeting(creatureID).
		UntilEndOfCombat().
		GrantAbility(abilityVigilance)
	
	effectManager := effects.NewEffectManager(gameState.layerSystem)
	effectManager.AddEffect(effect)
	
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
	if err := engine.DeclareAttacker(gameID, creatureID, "Bob", "Alice"); err != nil {
		t.Fatalf("failed to declare attacker: %v", err)
	}
	
	// Creature should NOT be tapped (vigilance prevents tapping)
	gameState.mu.RLock()
	creature := gameState.cards[creatureID]
	isTapped := creature.Tapped
	gameState.mu.RUnlock()
	
	if isTapped {
		t.Error("Creature with vigilance should not be tapped after attacking")
	}
}

// TestCombatDynamicFirstStrike verifies granted first strike works
func TestCombatDynamicFirstStrike(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-dynamic-first-strike"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Alice has a creature without first strike
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
		Abilities:    []EngineAbilityView{}, // No first strike initially
	}
	
	// Grant first strike dynamically
	effect := effects.NewEffectBuilder("battle-mastery").
		Targeting(creatureID).
		UntilEndOfTurn().
		GrantAbility(abilityFirstStrike)
	
	effectManager := effects.NewEffectManager(gameState.layerSystem)
	effectManager.AddEffect(effect)
	
	creature := gameState.cards[creatureID]
	hasFirstStrike := engine.hasFirstStrikeWithEffects(gameState, creature)
	gameState.mu.Unlock()
	
	if !hasFirstStrike {
		t.Error("Creature should have first strike after effect")
	}
	
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
	if err := engine.DeclareAttacker(gameID, creatureID, "Bob", "Alice"); err != nil {
		t.Fatalf("failed to declare attacker: %v", err)
	}
	
	// Check if there's first strike in combat
	hasFirstStrikeInCombat, err := engine.HasFirstOrDoubleStrike(gameID)
	if err != nil {
		t.Fatalf("failed to check first strike: %v", err)
	}
	
	if !hasFirstStrikeInCombat {
		t.Error("Combat should have first strike creature")
	}
}

// TestCombatDynamicAbilityCleanup verifies effects clean up correctly
func TestCombatDynamicAbilityCleanup(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-dynamic-cleanup"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: Alice has a creature
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
		Abilities:    []EngineAbilityView{},
	}
	
	// Grant flying until end of combat
	effect := effects.NewEffectBuilder("temporary-flight").
		Targeting(creatureID).
		UntilEndOfCombat().
		GrantAbility(abilityFlying)
	
	effectManager := effects.NewEffectManager(gameState.layerSystem)
	effectManager.AddEffect(effect)
	
	// Verify creature has flying
	creature := gameState.cards[creatureID]
	hasFlying := engine.hasAbilityWithEffects(gameState, creature, abilityFlying)
	if !hasFlying {
		t.Fatal("Creature should have flying before combat ends")
	}
	
	gameState.mu.Unlock()
	
	// Setup and end combat
	if err := engine.ResetCombat(gameID); err != nil {
		t.Fatalf("failed to reset combat: %v", err)
	}
	
	if err := engine.SetAttacker(gameID, "Alice"); err != nil {
		t.Fatalf("failed to set attacker: %v", err)
	}
	
	if err := engine.SetDefenders(gameID); err != nil {
		t.Fatalf("failed to set defenders: %v", err)
	}
	
	// End combat (should clean up effects)
	if err := engine.EndCombat(gameID); err != nil {
		t.Fatalf("failed to end combat: %v", err)
	}
	
	// Verify flying is gone
	gameState.mu.RLock()
	creature = gameState.cards[creatureID]
	hasFlying = engine.hasAbilityWithEffects(gameState, creature, abilityFlying)
	gameState.mu.RUnlock()
	
	if hasFlying {
		t.Error("Creature should not have flying after combat ends")
	}
}
