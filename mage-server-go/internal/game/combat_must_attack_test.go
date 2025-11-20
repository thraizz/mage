package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMustAttackIfAble tests the "attacks if able" combat mechanic
func TestMustAttackIfAble(t *testing.T) {
	h := NewCombatTestHarness(t, "game-1", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create a creature that must attack
	attackerID := h.CreateAttacker("attacker", "Forced Attacker", "Alice", "2", "2")

	// Add "must attack if able" effect
	effect := effects.NewMustAttackEffect(attackerID, []string{attackerID}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Process forced attackers
	err = h.engine.processForcedAttackers(gameState)
	require.NoError(t, err)

	// Verify creature was forced to attack
	attackerCard := gameState.cards[attackerID]
	assert.True(t, attackerCard.Attacking, "creature with 'attacks if able' should be attacking")
	assert.Contains(t, gameState.combat.attackers, attackerID, "attacker should be in combat.attackers")
	assert.NotEmpty(t, attackerCard.AttackingWhat, "attacker should be attacking something")
}

// TestMustAttackIfAble_Tapped tests that tapped creatures can't be forced to attack
func TestMustAttackIfAble_Tapped(t *testing.T) {
	h := NewCombatTestHarness(t, "game-2", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create a tapped creature that must attack
	attackerID := h.CreateCreature(CreatureSpec{
		ID:         "attacker",
		Name:       "Tapped Attacker",
		Power:      "2",
		Toughness:  "2",
		Controller: "Alice",
		Tapped:     true,
	})

	// Add "must attack if able" effect
	effect := effects.NewMustAttackEffect(attackerID, []string{attackerID}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Process forced attackers
	err = h.engine.processForcedAttackers(gameState)
	require.NoError(t, err)

	// Verify creature was NOT forced to attack (can't attack because tapped)
	attackerCard := gameState.cards[attackerID]
	assert.False(t, attackerCard.Attacking, "tapped creature should not be forced to attack")
	assert.NotContains(t, gameState.combat.attackers, attackerID, "tapped creature should not be in attackers")
}

// TestMustAttackIfAble_SummoningSickness tests that creatures with summoning sickness can't be forced to attack
func TestMustAttackIfAble_SummoningSickness(t *testing.T) {
	h := NewCombatTestHarness(t, "game-3", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create a creature with summoning sickness
	attackerID := h.CreateAttacker("attacker", "New Creature", "Alice", "2", "2")
	gameState.mu.Lock()
	gameState.cards[attackerID].SummoningSickness = true
	gameState.mu.Unlock()

	// Add "must attack if able" effect
	effect := effects.NewMustAttackEffect(attackerID, []string{attackerID}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Process forced attackers
	err = h.engine.processForcedAttackers(gameState)
	require.NoError(t, err)

	// Verify creature was NOT forced to attack (can't attack due to summoning sickness)
	attackerCard := gameState.cards[attackerID]
	assert.False(t, attackerCard.Attacking, "creature with summoning sickness should not be forced to attack")
	assert.NotContains(t, gameState.combat.attackers, attackerID, "creature with summoning sickness should not be in attackers")
}

// TestMustAttackIfAble_MultipleCreatures tests multiple creatures with "attacks if able"
func TestMustAttackIfAble_MultipleCreatures(t *testing.T) {
	h := NewCombatTestHarness(t, "game-4", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create multiple creatures that must attack
	attacker1 := h.CreateAttacker("attacker1", "Forced 1", "Alice", "2", "2")
	attacker2 := h.CreateAttacker("attacker2", "Forced 2", "Alice", "3", "3")
	attacker3 := h.CreateAttacker("attacker3", "Forced 3", "Alice", "1", "1")

	// Add "must attack if able" effects
	effect1 := effects.NewMustAttackEffect(attacker1, []string{attacker1}, effects.DurationWhileOnBattlefield)
	effect2 := effects.NewMustAttackEffect(attacker2, []string{attacker2}, effects.DurationWhileOnBattlefield)
	effect3 := effects.NewMustAttackEffect(attacker3, []string{attacker3}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect1)
	gameState.layerSystem.AddEffect(effect2)
	gameState.layerSystem.AddEffect(effect3)

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Process forced attackers
	err = h.engine.processForcedAttackers(gameState)
	require.NoError(t, err)

	// Verify all creatures were forced to attack
	assert.True(t, gameState.cards[attacker1].Attacking, "attacker 1 should be attacking")
	assert.True(t, gameState.cards[attacker2].Attacking, "attacker 2 should be attacking")
	assert.True(t, gameState.cards[attacker3].Attacking, "attacker 3 should be attacking")
	assert.Len(t, gameState.combat.attackers, 3, "should have 3 attackers")
}

// TestMustAttackIfAble_AlreadyAttacking tests that already-attacking creatures are skipped
func TestMustAttackIfAble_AlreadyAttacking(t *testing.T) {
	h := NewCombatTestHarness(t, "game-5", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create a creature and manually declare it as attacking
	attackerID := h.CreateAttacker("attacker", "Already Attacking", "Alice", "2", "2")

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Manually declare attacker
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)

	// Add "must attack if able" effect AFTER already attacking
	effect := effects.NewMustAttackEffect(attackerID, []string{attackerID}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Get initial attack state
	initialAttackingWhat := gameState.cards[attackerID].AttackingWhat

	// Process forced attackers (should not re-declare)
	err = h.engine.processForcedAttackers(gameState)
	require.NoError(t, err)

	// Verify creature is still attacking the same target (not re-declared)
	attackerCard := gameState.cards[attackerID]
	assert.True(t, attackerCard.Attacking, "creature should still be attacking")
	assert.Equal(t, initialAttackingWhat, attackerCard.AttackingWhat, "should attack same target")
}

// TestMustAttackIfAble_WithDefender tests creature with defender can't be forced to attack
func TestMustAttackIfAble_WithDefender(t *testing.T) {
	h := NewCombatTestHarness(t, "game-6", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create a creature with defender that must attack
	attackerID := h.CreateCreature(CreatureSpec{
		ID:         "defender",
		Name:       "Wall",
		Power:      "0",
		Toughness:  "4",
		Controller: "Alice",
		Abilities:  []string{abilityDefender},
	})

	// Add "must attack if able" effect
	effect := effects.NewMustAttackEffect(attackerID, []string{attackerID}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Process forced attackers
	err = h.engine.processForcedAttackers(gameState)
	require.NoError(t, err)

	// Verify creature was NOT forced to attack (has defender)
	attackerCard := gameState.cards[attackerID]
	assert.False(t, attackerCard.Attacking, "creature with defender should not be forced to attack")
	assert.NotContains(t, gameState.combat.attackers, attackerID, "creature with defender should not be in attackers")
}

// TestMustAttackIfAble_Integration tests the full integration with turn structure
func TestMustAttackIfAble_Integration(t *testing.T) {
	h := NewCombatTestHarness(t, "game-7", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create creatures
	attackerID := h.CreateAttacker("attacker", "Forced Attacker", "Alice", "3", "3")
	blockerID := h.CreateBlocker("blocker", "Blocker", "Bob", "2", "2")

	// Add "must attack if able" effect
	effect := effects.NewMustAttackEffect(attackerID, []string{attackerID}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Process forced attackers (called automatically in declare attackers step)
	err = h.engine.processForcedAttackers(gameState)
	require.NoError(t, err)

	// Verify attacker is attacking
	assert.True(t, gameState.cards[attackerID].Attacking, "creature should be forced to attack")

	// Declare blocker
	err = h.engine.DeclareBlocker(h.gameID, blockerID, attackerID, "Bob")
	require.NoError(t, err)

	// Verify blocking
	assert.True(t, gameState.cards[blockerID].Blocking, "blocker should be blocking")

	// Assign and apply damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify damage
	assert.Equal(t, 3, gameState.cards[blockerID].Damage, "blocker should take 3 damage")
	assert.Equal(t, 2, gameState.cards[attackerID].Damage, "attacker should take 2 damage")
}
