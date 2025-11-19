package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMustAttackIfAble tests the "attacks if able" combat mechanic
func TestMustAttackIfAble(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create a creature that must attack
	attacker := harness.createCreature(harness.player1, "Attacker", "2", "2")

	// Add "must attack if able" effect
	effect := effects.NewMustAttackEffect(attacker, []string{attacker}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat
	harness.setupCombat(harness.player1)

	// Process forced attackers
	err := harness.engine.processForcedAttackers(harness.gameState)
	require.NoError(t, err)

	// Verify creature was forced to attack
	attackerCard := harness.gameState.cards[attacker]
	assert.True(t, attackerCard.Attacking, "creature with 'attacks if able' should be attacking")
	assert.Contains(t, harness.gameState.combat.attackers, attacker, "attacker should be in combat.attackers")
	assert.NotEmpty(t, attackerCard.AttackingWhat, "attacker should be attacking something")
}

// TestMustAttackIfAble_Tapped tests that tapped creatures can't be forced to attack
func TestMustAttackIfAble_Tapped(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create a tapped creature that must attack
	attacker := harness.createCreature(harness.player1, "Attacker", "2", "2")
	harness.gameState.cards[attacker].Tapped = true

	// Add "must attack if able" effect
	effect := effects.NewMustAttackEffect(attacker, []string{attacker}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat
	harness.setupCombat(harness.player1)

	// Process forced attackers
	err := harness.engine.processForcedAttackers(harness.gameState)
	require.NoError(t, err)

	// Verify creature was NOT forced to attack (can't attack because tapped)
	attackerCard := harness.gameState.cards[attacker]
	assert.False(t, attackerCard.Attacking, "tapped creature should not be forced to attack")
	assert.NotContains(t, harness.gameState.combat.attackers, attacker, "tapped creature should not be in attackers")
}

// TestMustAttackIfAble_SummoningSickness tests that creatures with summoning sickness can't be forced to attack
func TestMustAttackIfAble_SummoningSickness(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create a creature with summoning sickness that must attack
	attacker := harness.createCreature(harness.player1, "Attacker", "2", "2")
	harness.gameState.cards[attacker].SummoningSickness = true

	// Add "must attack if able" effect
	effect := effects.NewMustAttackEffect(attacker, []string{attacker}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat
	harness.setupCombat(harness.player1)

	// Process forced attackers
	err := harness.engine.processForcedAttackers(harness.gameState)
	require.NoError(t, err)

	// Verify creature was NOT forced to attack (can't attack due to summoning sickness)
	attackerCard := harness.gameState.cards[attacker]
	assert.False(t, attackerCard.Attacking, "creature with summoning sickness should not be forced to attack")
	assert.NotContains(t, harness.gameState.combat.attackers, attacker, "creature with summoning sickness should not be in attackers")
}

// TestMustAttackIfAble_MultipleCreatures tests multiple creatures with "attacks if able"
func TestMustAttackIfAble_MultipleCreatures(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create multiple creatures that must attack
	attacker1 := harness.createCreature(harness.player1, "Attacker 1", "2", "2")
	attacker2 := harness.createCreature(harness.player1, "Attacker 2", "3", "3")
	attacker3 := harness.createCreature(harness.player1, "Attacker 3", "1", "1")

	// Add "must attack if able" effects
	effect1 := effects.NewMustAttackEffect(attacker1, []string{attacker1}, effects.DurationWhileOnBattlefield)
	effect2 := effects.NewMustAttackEffect(attacker2, []string{attacker2}, effects.DurationWhileOnBattlefield)
	effect3 := effects.NewMustAttackEffect(attacker3, []string{attacker3}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect1)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect2)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect3)

	// Setup combat
	harness.setupCombat(harness.player1)

	// Process forced attackers
	err := harness.engine.processForcedAttackers(harness.gameState)
	require.NoError(t, err)

	// Verify all creatures were forced to attack
	assert.True(t, harness.gameState.cards[attacker1].Attacking, "attacker 1 should be attacking")
	assert.True(t, harness.gameState.cards[attacker2].Attacking, "attacker 2 should be attacking")
	assert.True(t, harness.gameState.cards[attacker3].Attacking, "attacker 3 should be attacking")
	assert.Len(t, harness.gameState.combat.attackers, 3, "should have 3 attackers")
}

// TestMustAttackIfAble_AlreadyAttacking tests that already-attacking creatures are skipped
func TestMustAttackIfAble_AlreadyAttacking(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create a creature and manually declare it as attacking
	attacker := harness.createCreature(harness.player1, "Attacker", "2", "2")

	// Setup combat
	harness.setupCombat(harness.player1)

	// Manually declare attacker
	err := harness.engine.DeclareAttacker(harness.gameID, attacker, harness.player2, harness.player1)
	require.NoError(t, err)

	// Add "must attack if able" effect AFTER already attacking
	effect := effects.NewMustAttackEffect(attacker, []string{attacker}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Get initial attack state
	initialAttackingWhat := harness.gameState.cards[attacker].AttackingWhat

	// Process forced attackers (should not re-declare)
	err = harness.engine.processForcedAttackers(harness.gameState)
	require.NoError(t, err)

	// Verify creature is still attacking the same target (not re-declared)
	attackerCard := harness.gameState.cards[attacker]
	assert.True(t, attackerCard.Attacking, "creature should still be attacking")
	assert.Equal(t, initialAttackingWhat, attackerCard.AttackingWhat, "should attack same target")
}

// TestMustAttackIfAble_WithDefender tests creature with defender can't be forced to attack
func TestMustAttackIfAble_WithDefender(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create a creature with defender that must attack
	attacker := harness.createCreature(harness.player1, "Defender", "0", "4")
	harness.addAbility(attacker, abilityDefender)

	// Add "must attack if able" effect
	effect := effects.NewMustAttackEffect(attacker, []string{attacker}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat
	harness.setupCombat(harness.player1)

	// Process forced attackers
	err := harness.engine.processForcedAttackers(harness.gameState)
	require.NoError(t, err)

	// Verify creature was NOT forced to attack (has defender)
	attackerCard := harness.gameState.cards[attacker]
	assert.False(t, attackerCard.Attacking, "creature with defender should not be forced to attack")
	assert.NotContains(t, harness.gameState.combat.attackers, attacker, "creature with defender should not be in attackers")
}

// TestMustAttackIfAble_Integration tests the full integration with turn structure
func TestMustAttackIfAble_Integration(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create creatures
	attacker := harness.createCreature(harness.player1, "Forced Attacker", "3", "3")
	blocker := harness.createCreature(harness.player2, "Blocker", "2", "2")

	// Add "must attack if able" effect
	effect := effects.NewMustAttackEffect(attacker, []string{attacker}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat
	harness.setupCombat(harness.player1)

	// Process forced attackers (called automatically in declare attackers step)
	err := harness.engine.processForcedAttackers(harness.gameState)
	require.NoError(t, err)

	// Verify attacker is attacking
	assert.True(t, harness.gameState.cards[attacker].Attacking, "creature should be forced to attack")

	// Declare blocker
	err = harness.engine.DeclareBlocker(harness.gameID, blocker, attacker, harness.player2)
	require.NoError(t, err)

	// Verify blocking
	assert.True(t, harness.gameState.cards[blocker].Blocking, "blocker should be blocking")

	// Assign and apply damage
	err = harness.engine.AssignCombatDamage(harness.gameID, false)
	require.NoError(t, err)
	err = harness.engine.ApplyCombatDamage(harness.gameID)
	require.NoError(t, err)

	// Verify damage
	assert.Equal(t, 3, harness.gameState.cards[blocker].Damage, "blocker should take 3 damage")
	assert.Equal(t, 2, harness.gameState.cards[attacker].Damage, "attacker should take 2 damage")
}
