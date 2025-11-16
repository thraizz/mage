package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap/zaptest"
)

// TestCombatDeathtouchTrample verifies deathtouch + trample interaction
// With deathtouch, only 1 damage is lethal, so excess tramples through
func TestCombatDeathtouchTrample(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-deathtouch-trample"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	// Get game state
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: 4/4 attacker with deathtouch and trample vs 5/5 blocker
	gameState.mu.Lock()
	attackerID := "deathtouch-attacker"
	blockerID := "big-blocker"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Deathtouch Trampler",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "4",
		Toughness:    "4",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityDeathtouch, Text: "Deathtouch"},
			{ID: abilityTrample, Text: "Trample"},
		},
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
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	
	// Verify: Blocker should have 1 damage (lethal with deathtouch)
	// Bob should have lost 3 life (4 power - 1 lethal = 3 trample)
	gameState.mu.RLock()
	blockerCard := gameState.cards[blockerID]
	bobLife := gameState.players["Bob"].Life
	gameState.mu.RUnlock()
	
	if blockerCard.Damage != 1 {
		t.Errorf("Expected blocker to have 1 damage (deathtouch lethal), got %d", blockerCard.Damage)
	}
	
	expectedLife := initialBobLife - 3 // 4 power - 1 lethal = 3 trample
	if bobLife != expectedLife {
		t.Errorf("Expected Bob to have %d life (3 trample damage), got %d", expectedLife, bobLife)
	}
	
	engine.EndCombat(gameID)
}

// TestCombatDeathtouchTrampleMultipleBlockers verifies deathtouch with multiple blockers
// Each blocker only needs 1 damage to be lethal
func TestCombatDeathtouchTrampleMultipleBlockers(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-deathtouch-multi"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: 5/5 attacker with deathtouch and trample vs three 3/3 blockers
	gameState.mu.Lock()
	attackerID := "deathtouch-attacker"
	blocker1ID := "blocker1"
	blocker2ID := "blocker2"
	blocker3ID := "blocker3"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Deathtouch Trampler",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "5",
		Toughness:    "5",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityDeathtouch, Text: "Deathtouch"},
			{ID: abilityTrample, Text: "Trample"},
		},
	}
	
	for _, blockerID := range []string{blocker1ID, blocker2ID, blocker3ID} {
		gameState.cards[blockerID] = &internalCard{
			ID:           blockerID,
			Name:         "Blocker",
			Type:         "Creature",
			Zone:         zoneBattlefield,
			OwnerID:      "Bob",
			ControllerID: "Bob",
			Power:        "3",
			Toughness:    "3",
			Tapped:       false,
		}
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
	engine.DeclareBlocker(gameID, blocker3ID, attackerID, "Bob")
	engine.AcceptBlockers(gameID)
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	
	// Verify: Each blocker gets 1 damage (lethal with deathtouch)
	// Bob loses 2 life (5 power - 3 blockers * 1 damage = 2 trample)
	gameState.mu.RLock()
	b1 := gameState.cards[blocker1ID]
	b2 := gameState.cards[blocker2ID]
	b3 := gameState.cards[blocker3ID]
	bobLife := gameState.players["Bob"].Life
	gameState.mu.RUnlock()
	
	if b1.Damage != 1 {
		t.Errorf("Expected blocker1 to have 1 damage, got %d", b1.Damage)
	}
	if b2.Damage != 1 {
		t.Errorf("Expected blocker2 to have 1 damage, got %d", b2.Damage)
	}
	if b3.Damage != 1 {
		t.Errorf("Expected blocker3 to have 1 damage, got %d", b3.Damage)
	}
	
	expectedLife := initialBobLife - 2 // 5 power - 3 lethal = 2 trample
	if bobLife != expectedLife {
		t.Errorf("Expected Bob to have %d life, got %d", expectedLife, bobLife)
	}
	
	engine.EndCombat(gameID)
}

// TestCombatDeathtouchNoTrample verifies deathtouch without trample
// Excess damage doesn't trample through
func TestCombatDeathtouchNoTrample(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-deathtouch-no-trample"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup: 4/4 attacker with deathtouch (no trample) vs 1/1 blocker
	gameState.mu.Lock()
	attackerID := "deathtouch-attacker"
	blockerID := "small-blocker"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Deathtouch Creature",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "4",
		Toughness:    "4",
		Tapped:       false,
		Abilities: []EngineAbilityView{
			{ID: abilityDeathtouch, Text: "Deathtouch"},
		},
	}
	
	gameState.cards[blockerID] = &internalCard{
		ID:           blockerID,
		Name:         "Small Blocker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Power:        "1",
		Toughness:    "1",
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
	engine.AssignCombatDamage(gameID, false)
	engine.ApplyCombatDamage(gameID)
	
	// Verify: Blocker gets all 4 damage (even though 1 is lethal)
	// Bob takes no damage (no trample)
	gameState.mu.RLock()
	blockerCard := gameState.cards[blockerID]
	bobLife := gameState.players["Bob"].Life
	gameState.mu.RUnlock()
	
	// Without trample, all damage goes to blocker
	if blockerCard.Damage == 0 {
		t.Error("Expected blocker to have damage")
	}
	
	if bobLife != initialBobLife {
		t.Errorf("Expected Bob to have %d life (no trample), got %d", initialBobLife, bobLife)
	}
	
	engine.EndCombat(gameID)
}

// TestCombatEventUnblockedAttacker verifies EventUnblockedAttacker fires
func TestCombatEventUnblockedAttacker(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-event-unblocked"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup attacker
	gameState.mu.Lock()
	attackerID := "attacker"
	
	gameState.cards[attackerID] = &internalCard{
		ID:           attackerID,
		Name:         "Attacker",
		Type:         "Creature",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Power:        "2",
		Toughness:    "2",
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Subscribe to event
	eventFired := false
	gameState.eventBus.SubscribeTyped(rules.EventUnblockedAttacker, func(e rules.Event) {
		if e.SourceID == attackerID {
			eventFired = true
		}
	})
	
	// Attack without blockers
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.AcceptBlockers(gameID) // No blockers declared
	
	if !eventFired {
		t.Error("Expected EventUnblockedAttacker to fire for unblocked attacker")
	}
	
	engine.EndCombat(gameID)
}

// TestCombatEventRemovedFromCombat verifies EventRemovedFromCombat fires
func TestCombatEventRemovedFromCombat(t *testing.T) {
	logger := zaptest.NewLogger(t)
	engine := NewMageEngine(logger)
	
	gameID := "test-event-removed"
	players := []string{"Alice", "Bob"}
	
	if err := engine.StartGame(gameID, players, "Duel"); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	
	engine.mu.RLock()
	gameState := engine.games[gameID]
	engine.mu.RUnlock()
	
	// Setup creatures
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
		Tapped:       false,
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
		Tapped:       false,
	}
	gameState.mu.Unlock()
	
	// Subscribe to event
	eventFired := false
	gameState.eventBus.SubscribeTyped(rules.EventRemovedFromCombat, func(e rules.Event) {
		if e.TargetID == blockerID {
			eventFired = true
		}
	})
	
	// Set up combat
	engine.ResetCombat(gameID)
	engine.SetAttacker(gameID, "Alice")
	engine.SetDefenders(gameID)
	engine.DeclareAttacker(gameID, attackerID, "Bob", "Alice")
	engine.DeclareBlocker(gameID, blockerID, attackerID, "Bob")
	
	// Remove blocker from combat
	engine.RemoveBlocker(gameID, blockerID)
	
	if !eventFired {
		t.Error("Expected EventRemovedFromCombat to fire when blocker removed")
	}
	
	engine.EndCombat(gameID)
}
