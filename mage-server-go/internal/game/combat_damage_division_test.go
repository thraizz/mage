package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDamageDivision_AttackerMultipleBlockersDefault tests default damage division for attacker with multiple blockers
func TestDamageDivision_AttackerMultipleBlockersDefault(t *testing.T) {
	h := NewCombatTestHarness(t, "game-1", []string{"Alice", "Bob"})

	// Create attacker with 6 power
	attackerID := h.CreateAttacker("attacker", "Giant Growth Bear", "Alice", "6", "6")

	// Create three blockers with different toughness
	blocker1 := h.CreateBlocker("blocker1", "Small Bear", "Bob", "1", "1")
	blocker2 := h.CreateBlocker("blocker2", "Medium Bear", "Bob", "2", "2")
	blocker3 := h.CreateBlocker("blocker3", "Large Bear", "Bob", "3", "3")

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)

	// Declare all three creatures as blockers
	err = h.engine.DeclareBlocker(h.gameID, blocker1, attackerID, "Bob")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blocker2, attackerID, "Bob")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blocker3, attackerID, "Bob")
	require.NoError(t, err)

	// Assign and apply combat damage (using default division)
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Default division: 6 power / 3 blockers = 2 damage each
	gameState := h.GetGameState()
	assert.Equal(t, 2, gameState.cards[blocker1].Damage, "blocker1 should have 2 damage")
	assert.Equal(t, 2, gameState.cards[blocker2].Damage, "blocker2 should have 2 damage")
	assert.Equal(t, 2, gameState.cards[blocker3].Damage, "blocker3 should have 2 damage")
}

// TestDamageDivision_AttackerCustomAssignment tests custom damage assignment for attacker
func TestDamageDivision_AttackerCustomAssignment(t *testing.T) {
	h := NewCombatTestHarness(t, "game-2", []string{"Alice", "Bob"})

	// Create attacker with 6 power
	attackerID := h.CreateAttacker("attacker", "Bear", "Alice", "6", "6")

	// Create three blockers
	blocker1 := h.CreateBlocker("blocker1", "Bear 1", "Bob", "1", "1")
	blocker2 := h.CreateBlocker("blocker2", "Bear 2", "Bob", "2", "2")
	blocker3 := h.CreateBlocker("blocker3", "Bear 3", "Bob", "3", "3")

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

	// Custom assignment: 4 damage to blocker1, 1 to blocker2, 1 to blocker3
	damageMap := map[string]int{
		blocker1: 4,
		blocker2: 1,
		blocker3: 1,
	}
	err = h.engine.AssignAttackerDamage(h.gameID, attackerID, damageMap)
	require.NoError(t, err)

	// Apply combat damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify custom assignment was used
	gameState := h.GetGameState()
	assert.Equal(t, 4, gameState.cards[blocker1].Damage, "blocker1 should have 4 damage")
	assert.Equal(t, 1, gameState.cards[blocker2].Damage, "blocker2 should have 1 damage")
	assert.Equal(t, 1, gameState.cards[blocker3].Damage, "blocker3 should have 1 damage")
}

// TestDamageDivision_TrampleCustomAssignment tests custom damage assignment with trample
func TestDamageDivision_TrampleCustomAssignment(t *testing.T) {
	h := NewCombatTestHarness(t, "game-3", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create attacker with 6 power and trample
	attackerID := h.CreateAttacker("attacker", "Trample Bear", "Alice", "6", "6")
	gameState.mu.Lock()
	gameState.cards[attackerID].Abilities = append(gameState.cards[attackerID].Abilities, EngineAbilityView{
		ID:   "TrampleAbility",
		Text: "Trample",
		Rule: "Trample",
	})
	gameState.mu.Unlock()

	// Create two blockers with 2 toughness each
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

	// Custom assignment with trample: 2 to blocker1, 1 to blocker2, 3 tramples through
	damageMap := map[string]int{
		blocker1: 2,
		blocker2: 1,
	}
	err = h.engine.AssignAttackerDamage(h.gameID, attackerID, damageMap)
	require.NoError(t, err)

	// Apply combat damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify assignments
	gameState.mu.RLock()
	defer gameState.mu.RUnlock()
	assert.Equal(t, 2, gameState.cards[blocker1].Damage, "blocker1 should have 2 damage")
	assert.Equal(t, 1, gameState.cards[blocker2].Damage, "blocker2 should have 1 damage")
	assert.Equal(t, 17, gameState.players["Bob"].Life, "Bob should have taken 3 trample damage (20 - 3 = 17)")
}

// TestDamageDivision_ValidationTotalMismatch tests validation of total damage
func TestDamageDivision_ValidationTotalMismatch(t *testing.T) {
	h := NewCombatTestHarness(t, "game-4", []string{"Alice", "Bob"})

	// Create attacker with 6 power
	attackerID := h.CreateAttacker("attacker", "Bear", "Alice", "6", "6")

	// Create two blockers
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

	// Try to assign only 5 damage (should fail - must assign all 6)
	damageMap := map[string]int{
		blocker1: 3,
		blocker2: 2,
	}
	err = h.engine.AssignAttackerDamage(h.gameID, attackerID, damageMap)
	assert.Error(t, err, "should reject assignment that doesn't total to power")
	assert.Contains(t, err.Error(), "must assign all damage")
}

// TestDamageDivision_ValidationInvalidTarget tests validation of blocker targets
func TestDamageDivision_ValidationInvalidTarget(t *testing.T) {
	h := NewCombatTestHarness(t, "game-5", []string{"Alice", "Bob"})

	// Create attacker
	attackerID := h.CreateAttacker("attacker", "Bear", "Alice", "6", "6")

	// Create two blockers
	blocker1 := h.CreateBlocker("blocker1", "Bear 1", "Bob", "2", "2")
	blocker2 := h.CreateBlocker("blocker2", "Bear 2", "Bob", "2", "2")

	// Create a third creature that is NOT blocking
	nonBlockerID := h.CreateBlocker("non-blocker", "Idle Bear", "Bob", "3", "3")

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

	// Try to assign damage to non-blocker (should fail)
	damageMap := map[string]int{
		blocker1:     2,
		blocker2:     2,
		nonBlockerID: 2,
	}
	err = h.engine.AssignAttackerDamage(h.gameID, attackerID, damageMap)
	assert.Error(t, err, "should reject assignment to non-blocker")
	assert.Contains(t, err.Error(), "is not blocking this attacker")
}

// TestDamageDivision_BlockerMultipleAttackers tests blocker dividing damage among multiple attackers
func TestDamageDivision_BlockerMultipleAttackers(t *testing.T) {
	h := NewCombatTestHarness(t, "game-6", []string{"Alice", "Bob"})

	// Create three attackers
	attacker1 := h.CreateAttacker("attacker1", "Bear 1", "Alice", "2", "2")
	attacker2 := h.CreateAttacker("attacker2", "Bear 2", "Alice", "2", "2")
	attacker3 := h.CreateAttacker("attacker3", "Bear 3", "Alice", "2", "2")

	// Create blocker with 6 power
	blockerID := h.CreateBlocker("blocker", "Giant Bear", "Bob", "6", "6")

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Declare all three attackers
	err = h.engine.DeclareAttacker(h.gameID, attacker1, "Bob", "Alice")
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attacker2, "Bob", "Alice")
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attacker3, "Bob", "Alice")
	require.NoError(t, err)

	// Block all three with single blocker
	err = h.engine.DeclareBlocker(h.gameID, blockerID, attacker1, "Bob")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blockerID, attacker2, "Bob")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blockerID, attacker3, "Bob")
	require.NoError(t, err)

	// Custom assignment: 4 damage to attacker1, 1 each to others
	damageMap := map[string]int{
		attacker1: 4,
		attacker2: 1,
		attacker3: 1,
	}
	err = h.engine.AssignBlockerDamage(h.gameID, blockerID, damageMap)
	require.NoError(t, err)

	// Apply combat damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify damage assignments
	gameState := h.GetGameState()
	assert.Equal(t, 4, gameState.cards[attacker1].Damage, "attacker1 should have 4 damage")
	assert.Equal(t, 1, gameState.cards[attacker2].Damage, "attacker2 should have 1 damage")
	assert.Equal(t, 1, gameState.cards[attacker3].Damage, "attacker3 should have 1 damage")

	// Blocker should also take damage from all three attackers (6 damage total)
	assert.Equal(t, 6, gameState.cards[blockerID].Damage, "blocker should have 6 damage")
}

// TestDamageDivision_BlockerValidationTotal tests blocker damage assignment validation
func TestDamageDivision_BlockerValidationTotal(t *testing.T) {
	h := NewCombatTestHarness(t, "game-7", []string{"Alice", "Bob"})

	// Create two attackers
	attacker1 := h.CreateAttacker("attacker1", "Bear 1", "Alice", "2", "2")
	attacker2 := h.CreateAttacker("attacker2", "Bear 2", "Alice", "2", "2")

	// Create blocker with 6 power
	blockerID := h.CreateBlocker("blocker", "Giant Bear", "Bob", "6", "6")

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attacker1, "Bob", "Alice")
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attacker2, "Bob", "Alice")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blockerID, attacker1, "Bob")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blockerID, attacker2, "Bob")
	require.NoError(t, err)

	// Try to assign only 4 damage (should fail - must assign all 6)
	damageMap := map[string]int{
		attacker1: 2,
		attacker2: 2,
	}
	err = h.engine.AssignBlockerDamage(h.gameID, blockerID, damageMap)
	assert.Error(t, err, "should reject assignment that doesn't total to power")
	assert.Contains(t, err.Error(), "must assign all damage")
}

// TestDamageDivision_DefaultEvenSplit tests default even split with remainder
func TestDamageDivision_DefaultEvenSplit(t *testing.T) {
	h := NewCombatTestHarness(t, "game-8", []string{"Alice", "Bob"})

	// Create attacker with 7 power (odd number)
	attackerID := h.CreateAttacker("attacker", "Bear", "Alice", "7", "7")

	// Create three blockers
	blocker1 := h.CreateBlocker("blocker1", "Bear 1", "Bob", "3", "3")
	blocker2 := h.CreateBlocker("blocker2", "Bear 2", "Bob", "3", "3")
	blocker3 := h.CreateBlocker("blocker3", "Bear 3", "Bob", "3", "3")

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

	// Apply combat damage (using default division)
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Default division: 7 / 3 = 2 each, with 1 remainder to first blocker
	// So: 3, 2, 2
	gameState := h.GetGameState()
	assert.Equal(t, 3, gameState.cards[blocker1].Damage, "blocker1 should have 3 damage (gets remainder)")
	assert.Equal(t, 2, gameState.cards[blocker2].Damage, "blocker2 should have 2 damage")
	assert.Equal(t, 2, gameState.cards[blocker3].Damage, "blocker3 should have 2 damage")
}

// TestDamageDivision_TrampleDefaultLethalAssignment tests trample default assigns lethal to each
func TestDamageDivision_TrampleDefaultLethalAssignment(t *testing.T) {
	h := NewCombatTestHarness(t, "game-9", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create attacker with 10 power and trample
	attackerID := h.CreateAttacker("attacker", "Trample Bear", "Alice", "10", "10")
	gameState.mu.Lock()
	gameState.cards[attackerID].Abilities = append(gameState.cards[attackerID].Abilities, EngineAbilityView{
		ID:   "TrampleAbility",
		Text: "Trample",
		Rule: "Trample",
	})
	gameState.mu.Unlock()

	// Create three blockers with 1, 2, and 3 toughness
	blocker1 := h.CreateBlocker("blocker1", "Small Bear", "Bob", "1", "1")
	blocker2 := h.CreateBlocker("blocker2", "Medium Bear", "Bob", "2", "2")
	blocker3 := h.CreateBlocker("blocker3", "Large Bear", "Bob", "3", "3")

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

	// Apply combat damage (using default - should assign lethal to each, trample rest)
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify: 1 + 2 + 3 = 6 damage to blockers, 4 tramples through
	gameState.mu.RLock()
	defer gameState.mu.RUnlock()
	assert.Equal(t, 1, gameState.cards[blocker1].Damage, "blocker1 should have 1 damage (lethal)")
	assert.Equal(t, 2, gameState.cards[blocker2].Damage, "blocker2 should have 2 damage (lethal)")
	assert.Equal(t, 3, gameState.cards[blocker3].Damage, "blocker3 should have 3 damage (lethal)")
	assert.Equal(t, 16, gameState.players["Bob"].Life, "Bob should have taken 4 trample damage (20 - 4 = 16)")
}
