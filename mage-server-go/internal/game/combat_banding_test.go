package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBanding_DefenderAssignsDamage tests Rule 702.22j
// When blocked by creature with banding, defending player assigns attacker's damage
func TestBanding_DefenderAssignsDamage(t *testing.T) {
	h := NewCombatTestHarness(t, "game-1", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create attacker with 6 power
	attackerID := h.CreateAttacker("attacker", "Giant Bear", "Alice", "6", "6")

	// Create two blockers - one with banding
	blocker1 := h.CreateBlocker("blocker1", "Banding Bear", "Bob", "2", "2")
	blocker2 := h.CreateBlocker("blocker2", "Normal Bear", "Bob", "2", "2")

	// Give blocker1 banding ability
	gameState.mu.Lock()
	gameState.cards[blocker1].Abilities = append(gameState.cards[blocker1].Abilities, EngineAbilityView{
		ID:   "BandingAbility",
		Text: "Banding",
		Rule: "Banding",
	})
	gameState.mu.Unlock()

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blocker1, attackerID, "Bob")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blocker2, attackerID, "Bob")
	require.NoError(t, err)

	// Rule 702.22j: Defending player (Bob) assigns damage, not attacking player (Alice)
	damageMap := map[string]int{
		blocker1: 5, // Bob chooses to assign most damage to banding creature
		blocker2: 1,
	}

	// Alice (attacker) should NOT be allowed to assign (would normally control)
	err = h.engine.AssignAttackerDamage(h.gameID, attackerID, "Alice", damageMap)
	assert.Error(t, err, "attacking player should not control damage assignment")
	assert.Contains(t, err.Error(), "defending player must assign")

	// Bob (defender) should be allowed to assign
	err = h.engine.AssignAttackerDamage(h.gameID, attackerID, "Bob", damageMap)
	require.NoError(t, err, "defending player should control damage assignment")

	// Apply combat damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify damage was assigned as Bob specified
	gameState.mu.RLock()
	defer gameState.mu.RUnlock()
	assert.Equal(t, 5, gameState.cards[blocker1].Damage, "blocker1 should have 5 damage (Bob's choice)")
	assert.Equal(t, 1, gameState.cards[blocker2].Damage, "blocker2 should have 1 damage (Bob's choice)")
}

// TestBanding_AttackerAssignsDamage tests Rule 702.22k
// When blocking creature with banding, attacking player assigns blocker's damage
func TestBanding_AttackerAssignsDamage(t *testing.T) {
	h := NewCombatTestHarness(t, "game-2", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create TWO attackers with banding (so blocker blocks banding creatures)
	attacker1 := h.CreateAttacker("attacker1", "Banding Bear 1", "Alice", "3", "3")
	attacker2 := h.CreateAttacker("attacker2", "Normal Bear", "Alice", "3", "3")

	// Give attacker1 banding (Rule 702.22k: blocker is blocking creature with banding)
	gameState.mu.Lock()
	gameState.cards[attacker1].Abilities = append(gameState.cards[attacker1].Abilities, EngineAbilityView{
		ID:   "BandingAbility",
		Text: "Banding",
		Rule: "Banding",
	})
	gameState.mu.Unlock()

	// Create blocker with 6 power
	blockerID := h.CreateBlocker("blocker", "Giant", "Bob", "6", "6")

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attacker1, "Bob", "Alice")
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attacker2, "Bob", "Alice")
	require.NoError(t, err)

	// Blocker blocks both attackers
	err = h.engine.DeclareBlocker(h.gameID, blockerID, attacker1, "Bob")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blockerID, attacker2, "Bob")
	require.NoError(t, err)

	// Rule 702.22k: Attacking player (Alice) assigns blocker's damage, not defending player (Bob)
	damageMap := map[string]int{
		attacker1: 5, // Alice chooses to kill attacker1
		attacker2: 1,
	}

	// Bob (defender) should NOT be allowed to assign (would normally control his blocker)
	err = h.engine.AssignBlockerDamage(h.gameID, blockerID, "Bob", damageMap)
	assert.Error(t, err, "defending player should not control damage assignment")
	assert.Contains(t, err.Error(), "attacking player must assign")

	// Alice (attacker) should be allowed to assign
	err = h.engine.AssignBlockerDamage(h.gameID, blockerID, "Alice", damageMap)
	require.NoError(t, err, "attacking player should control damage assignment")

	// Apply combat damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify damage was assigned as Alice specified
	gameState.mu.RLock()
	defer gameState.mu.RUnlock()
	assert.Equal(t, 5, gameState.cards[attacker1].Damage, "attacker1 should have 5 damage (Alice's choice)")
	assert.Equal(t, 1, gameState.cards[attacker2].Damage, "attacker2 should have 1 damage (Alice's choice)")
}

// TestBanding_NormalCombatWithoutBanding tests that normal combat still works
func TestBanding_NormalCombatWithoutBanding(t *testing.T) {
	h := NewCombatTestHarness(t, "game-3", []string{"Alice", "Bob"})

	// Create attacker
	attackerID := h.CreateAttacker("attacker", "Bear", "Alice", "4", "4")

	// Create two blockers WITHOUT banding
	blocker1 := h.CreateBlocker("blocker1", "Bear 1", "Bob", "2", "2")
	blocker2 := h.CreateBlocker("blocker2", "Bear 2", "Bob", "2", "2")

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blocker1, attackerID, "Bob")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blocker2, attackerID, "Bob")
	require.NoError(t, err)

	// Without banding: Attacking player (Alice) controls damage assignment
	damageMap := map[string]int{
		blocker1: 3,
		blocker2: 1,
	}

	// Alice should be allowed to assign (normal case)
	err = h.engine.AssignAttackerDamage(h.gameID, attackerID, "Alice", damageMap)
	require.NoError(t, err, "attacking player should control damage assignment (no banding)")

	// Bob should NOT be allowed (no banding)
	damageMap2 := map[string]int{
		blocker1: 2,
		blocker2: 2,
	}
	err = h.engine.AssignAttackerDamage(h.gameID, attackerID, "Bob", damageMap2)
	assert.Error(t, err, "defending player should not control without banding")
	assert.Contains(t, err.Error(), "attacking player must assign")
}

// TestBanding_OnlyOneBandingNeeded tests that only one blocker needs banding
func TestBanding_OnlyOneBandingNeeded(t *testing.T) {
	h := NewCombatTestHarness(t, "game-4", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create attacker
	attackerID := h.CreateAttacker("attacker", "Bear", "Alice", "6", "6")

	// Create three blockers - only ONE has banding
	blocker1 := h.CreateBlocker("blocker1", "Normal 1", "Bob", "2", "2")
	blocker2 := h.CreateBlocker("blocker2", "Banding Bear", "Bob", "2", "2")
	blocker3 := h.CreateBlocker("blocker3", "Normal 2", "Bob", "2", "2")

	// Give only blocker2 banding
	gameState.mu.Lock()
	gameState.cards[blocker2].Abilities = append(gameState.cards[blocker2].Abilities, EngineAbilityView{
		ID:   "BandingAbility",
		Text: "Banding",
		Rule: "Banding",
	})
	gameState.mu.Unlock()

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blocker1, attackerID, "Bob")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blocker2, attackerID, "Bob")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blocker3, attackerID, "Bob")
	require.NoError(t, err)

	// Even with only one banding blocker, defending player controls damage
	damageMap := map[string]int{
		blocker1: 1,
		blocker2: 1,
		blocker3: 4,
	}

	// Bob (defender) should control assignment
	err = h.engine.AssignAttackerDamage(h.gameID, attackerID, "Bob", damageMap)
	require.NoError(t, err, "defending player should control (one blocker has banding)")

	// Apply combat damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify damage
	gameState.mu.RLock()
	defer gameState.mu.RUnlock()
	assert.Equal(t, 1, gameState.cards[blocker1].Damage)
	assert.Equal(t, 1, gameState.cards[blocker2].Damage)
	assert.Equal(t, 4, gameState.cards[blocker3].Damage)
}
