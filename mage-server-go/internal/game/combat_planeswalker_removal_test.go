package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPlaneswalkerRemoval_DuringCombat tests handling when a planeswalker is removed while being attacked
func TestPlaneswalkerRemoval_DuringCombat(t *testing.T) {
	h := NewCombatTestHarness(t, "game-1", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create a planeswalker with 3 loyalty
	planeswalkerID := "jace"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Jace",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Loyalty:      "3",
		Counters:     counters.NewCounters(),
	}
	planeswalker.Counters.AddCounter(counters.NewCounter("loyalty", 3))
	gameState.cards[planeswalkerID] = planeswalker
	gameState.mu.Unlock()

	// Create attacker
	attackerID := h.CreateAttacker("attacker", "Grizzly Bears", "Alice", "2", "2")

	// Setup combat and attack the planeswalker
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Verify attacker is attacking the planeswalker
	attacker := gameState.cards[attackerID]
	assert.True(t, attacker.Attacking, "creature should be attacking")
	assert.Equal(t, planeswalkerID, attacker.AttackingWhat, "creature should be attacking planeswalker")

	// Simulate planeswalker being destroyed (e.g., by instant-speed removal)
	gameState.mu.Lock()
	planeswalker.Zone = zoneGraveyard
	delete(gameState.cards, planeswalkerID)
	gameState.mu.Unlock()

	// The attacker remains in combat but attacking a now-invalid target
	// This is correct behavior - the creature is still attacking, just the defender is gone
	// Damage will not be dealt since the planeswalker no longer exists
	assert.True(t, attacker.Attacking, "creature should still be in attacking state")
	assert.Equal(t, planeswalkerID, attacker.AttackingWhat, "creature still has attacking target recorded")

	// When damage is assigned, it should handle the missing planeswalker gracefully
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)

	// Damage assignment should not crash even though planeswalker is gone
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)
}

// TestPlaneswalkerRemoval_StateBasedActionDuringCombat tests planeswalker death at 0 loyalty during combat
func TestPlaneswalkerRemoval_StateBasedActionDuringCombat(t *testing.T) {
	h := NewCombatTestHarness(t, "game-2", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create a planeswalker with 2 loyalty (will die from 2 damage)
	planeswalkerID := "liliana"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Liliana",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Loyalty:      "2",
		Counters:     counters.NewCounters(),
	}
	planeswalker.Counters.AddCounter(counters.NewCounter("loyalty", 2))
	gameState.cards[planeswalkerID] = planeswalker
	gameState.battlefield = append(gameState.battlefield, planeswalker)
	gameState.mu.Unlock()

	// Create attacker with 2 power (will deal lethal)
	attackerID := h.CreateAttacker("attacker", "Grizzly Bears", "Alice", "2", "2")

	// Setup combat and attack the planeswalker
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Apply damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify planeswalker has 0 loyalty
	assert.Equal(t, 0, planeswalker.Counters.GetCount("loyalty"), "planeswalker should have 0 loyalty")

	// State-based actions would move it to graveyard (Rule 704.5i)
	// The actual state-based action processing happens outside combat damage
	// but we verify the loyalty is correctly at 0
}

// TestPlaneswalkerRemoval_MultipleAttackersOneBlinks tests when a planeswalker blinks while multiple creatures attack it
func TestPlaneswalkerRemoval_MultipleAttackersOneBlinks(t *testing.T) {
	h := NewCombatTestHarness(t, "game-3", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create a planeswalker
	planeswalkerID := "garruk"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Garruk",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Loyalty:      "4",
		Counters:     counters.NewCounters(),
	}
	planeswalker.Counters.AddCounter(counters.NewCounter("loyalty", 4))
	gameState.cards[planeswalkerID] = planeswalker
	gameState.mu.Unlock()

	// Create two attackers
	attacker1 := h.CreateAttacker("attacker1", "Bear", "Alice", "2", "2")
	attacker2 := h.CreateAttacker("attacker2", "Wolf", "Alice", "3", "3")

	// Setup combat and attack with both
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attacker1, planeswalkerID, "Alice")
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attacker2, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Both attackers should be attacking the planeswalker
	assert.Equal(t, planeswalkerID, gameState.cards[attacker1].AttackingWhat)
	assert.Equal(t, planeswalkerID, gameState.cards[attacker2].AttackingWhat)

	// Simulate planeswalker being bounced/blinked
	gameState.mu.Lock()
	planeswalker.Zone = zoneHand
	delete(gameState.cards, planeswalkerID)
	gameState.mu.Unlock()

	// Both attackers remain attacking (the target just no longer exists)
	assert.True(t, gameState.cards[attacker1].Attacking, "attacker 1 should still be attacking")
	assert.True(t, gameState.cards[attacker2].Attacking, "attacker 2 should still be attacking")

	// Damage assignment should not crash
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)
}
