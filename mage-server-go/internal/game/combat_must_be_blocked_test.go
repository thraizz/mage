package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMustBeBlockedIfAble tests the "must be blocked if able" combat mechanic (lure)
func TestMustBeBlockedIfAble(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create creatures
	attacker := harness.createCreature(harness.player1, "Attacker with Lure", "3", "3")
	blocker := harness.createCreature(harness.player2, "Blocker", "2", "2")

	// Add "must be blocked if able" effect (Lure)
	effect := effects.NewMustBeBlockedEffect(attacker, "", []string{}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat and declare attacker
	harness.setupCombat(harness.player1)
	err := harness.engine.DeclareAttacker(harness.gameID, attacker, harness.player2, harness.player1)
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = harness.engine.processMustBeBlockedRequirements(harness.gameState)
	require.NoError(t, err)

	// Verify blocker is required to block the attacker
	requirements, ok := harness.gameState.combat.creatureMustBlockAttackers[blocker]
	require.True(t, ok, "blocker should have blocking requirements")
	assert.True(t, requirements[attacker], "blocker must block the attacker")
}

// TestMustBeBlockedIfAble_MultipleBlockers tests lure with multiple potential blockers
func TestMustBeBlockedIfAble_MultipleBlockers(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create creatures
	attacker := harness.createCreature(harness.player1, "Attacker with Lure", "4", "4")
	blocker1 := harness.createCreature(harness.player2, "Blocker 1", "2", "2")
	blocker2 := harness.createCreature(harness.player2, "Blocker 2", "3", "3")
	blocker3 := harness.createCreature(harness.player2, "Blocker 3", "1", "1")

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attacker, "", []string{}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat and declare attacker
	harness.setupCombat(harness.player1)
	err := harness.engine.DeclareAttacker(harness.gameID, attacker, harness.player2, harness.player1)
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = harness.engine.processMustBeBlockedRequirements(harness.gameState)
	require.NoError(t, err)

	// Verify all able blockers are required to block
	assert.True(t, harness.gameState.combat.creatureMustBlockAttackers[blocker1][attacker], "blocker 1 must block")
	assert.True(t, harness.gameState.combat.creatureMustBlockAttackers[blocker2][attacker], "blocker 2 must block")
	assert.True(t, harness.gameState.combat.creatureMustBlockAttackers[blocker3][attacker], "blocker 3 must block")
}

// TestMustBeBlockedIfAble_TappedBlocker tests that tapped creatures can't block even with lure
func TestMustBeBlockedIfAble_TappedBlocker(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create creatures
	attacker := harness.createCreature(harness.player1, "Attacker with Lure", "3", "3")
	blocker := harness.createCreature(harness.player2, "Tapped Blocker", "2", "2")
	harness.gameState.cards[blocker].Tapped = true

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attacker, "", []string{}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat and declare attacker
	harness.setupCombat(harness.player1)
	err := harness.engine.DeclareAttacker(harness.gameID, attacker, harness.player2, harness.player1)
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = harness.engine.processMustBeBlockedRequirements(harness.gameState)
	require.NoError(t, err)

	// Verify tapped blocker is NOT required to block (can't block because tapped)
	_, hasRequirement := harness.gameState.combat.creatureMustBlockAttackers[blocker]
	assert.False(t, hasRequirement, "tapped blocker should not be required to block")
}

// TestMustBeBlockedIfAble_FlyingAttacker tests lure with flying (only reach/flying can block)
func TestMustBeBlockedIfAble_FlyingAttacker(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create creatures
	attacker := harness.createCreature(harness.player1, "Flying Attacker with Lure", "3", "3")
	harness.addAbility(attacker, abilityFlying)

	groundBlocker := harness.createCreature(harness.player2, "Ground Blocker", "2", "2")
	flyingBlocker := harness.createCreature(harness.player2, "Flying Blocker", "2", "2")
	harness.addAbility(flyingBlocker, abilityFlying)

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attacker, "", []string{}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat and declare attacker
	harness.setupCombat(harness.player1)
	err := harness.engine.DeclareAttacker(harness.gameID, attacker, harness.player2, harness.player1)
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = harness.engine.processMustBeBlockedRequirements(harness.gameState)
	require.NoError(t, err)

	// Verify only flying blocker is required to block
	_, groundHasRequirement := harness.gameState.combat.creatureMustBlockAttackers[groundBlocker]
	flyingHasRequirement := harness.gameState.combat.creatureMustBlockAttackers[flyingBlocker][attacker]

	assert.False(t, groundHasRequirement, "ground blocker cannot block flying attacker")
	assert.True(t, flyingHasRequirement, "flying blocker must block flying attacker")
}

// TestMustBeBlockedIfAble_WithMenace tests lure interaction with menace
func TestMustBeBlockedIfAble_WithMenace(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create creatures
	attacker := harness.createCreature(harness.player1, "Attacker with Lure and Menace", "4", "4")
	harness.addAbility(attacker, abilityMenace)

	blocker1 := harness.createCreature(harness.player2, "Blocker 1", "2", "2")
	blocker2 := harness.createCreature(harness.player2, "Blocker 2", "2", "2")

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attacker, "", []string{}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat and declare attacker
	harness.setupCombat(harness.player1)
	err := harness.engine.DeclareAttacker(harness.gameID, attacker, harness.player2, harness.player1)
	require.NoError(t, err)

	// Set menace requirement (minimum 2 blockers)
	harness.gameState.combat.minBlockersPerAttacker[attacker] = 2

	// Process must-be-blocked requirements
	err = harness.engine.processMustBeBlockedRequirements(harness.gameState)
	require.NoError(t, err)

	// Verify both blockers are required to block (to satisfy menace)
	assert.True(t, harness.gameState.combat.creatureMustBlockAttackers[blocker1][attacker], "blocker 1 must block")
	assert.True(t, harness.gameState.combat.creatureMustBlockAttackers[blocker2][attacker], "blocker 2 must block")
}

// TestMustBeBlockedIfAble_SpecificBlockers tests lure targeting specific blockers
func TestMustBeBlockedIfAble_SpecificBlockers(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create creatures
	attacker := harness.createCreature(harness.player1, "Attacker", "3", "3")
	targetBlocker := harness.createCreature(harness.player2, "Target Blocker", "2", "2")
	otherBlocker := harness.createCreature(harness.player2, "Other Blocker", "2", "2")

	// Add "must be blocked if able" effect targeting only one blocker
	effect := effects.NewMustBeBlockedEffect(attacker, "", []string{targetBlocker}, effects.DurationEndOfTurn)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat and declare attacker
	harness.setupCombat(harness.player1)
	err := harness.engine.DeclareAttacker(harness.gameID, attacker, harness.player2, harness.player1)
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = harness.engine.processMustBeBlockedRequirements(harness.gameState)
	require.NoError(t, err)

	// Verify only target blocker is required to block
	targetHasRequirement := harness.gameState.combat.creatureMustBlockAttackers[targetBlocker][attacker]
	_, otherHasRequirement := harness.gameState.combat.creatureMustBlockAttackers[otherBlocker]

	assert.True(t, targetHasRequirement, "target blocker must block")
	assert.False(t, otherHasRequirement, "other blocker should not be required to block")
}

// TestMustBeBlockedIfAble_CheckBlockRequirements tests validation of blocking requirements
func TestMustBeBlockedIfAble_CheckBlockRequirements(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create creatures
	attacker := harness.createCreature(harness.player1, "Attacker with Lure", "3", "3")
	blocker := harness.createCreature(harness.player2, "Blocker", "2", "2")

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attacker, "", []string{}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat and declare attacker
	harness.setupCombat(harness.player1)
	err := harness.engine.DeclareAttacker(harness.gameID, attacker, harness.player2, harness.player1)
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = harness.engine.processMustBeBlockedRequirements(harness.gameState)
	require.NoError(t, err)

	// Check block requirements without declaring blocker
	violations, err := harness.engine.CheckBlockRequirements(harness.gameID, harness.player2)
	require.NoError(t, err)

	// Should have violation because blocker didn't block
	assert.NotEmpty(t, violations, "should have blocking requirement violation")
}

// TestMustBeBlockedIfAble_CheckBlockRequirements_Satisfied tests validation when requirement is met
func TestMustBeBlockedIfAble_CheckBlockRequirements_Satisfied(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create creatures
	attacker := harness.createCreature(harness.player1, "Attacker with Lure", "3", "3")
	blocker := harness.createCreature(harness.player2, "Blocker", "2", "2")

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attacker, "", []string{}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat and declare attacker
	harness.setupCombat(harness.player1)
	err := harness.engine.DeclareAttacker(harness.gameID, attacker, harness.player2, harness.player1)
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = harness.engine.processMustBeBlockedRequirements(harness.gameState)
	require.NoError(t, err)

	// Declare blocker to satisfy requirement
	err = harness.engine.DeclareBlocker(harness.gameID, blocker, attacker, harness.player2)
	require.NoError(t, err)

	// Check block requirements after declaring blocker
	violations, err := harness.engine.CheckBlockRequirements(harness.gameID, harness.player2)
	require.NoError(t, err)

	// Should have no violations because blocker blocked
	assert.Empty(t, violations, "should have no blocking requirement violations")
}

// TestMustBeBlockedIfAble_Integration tests full combat flow with lure
func TestMustBeBlockedIfAble_Integration(t *testing.T) {
	harness := newCombatTestHarness(t)
	defer harness.cleanup()

	// Create creatures
	attacker := harness.createCreature(harness.player1, "Attacker with Lure", "5", "5")
	blocker1 := harness.createCreature(harness.player2, "Blocker 1", "2", "2")
	blocker2 := harness.createCreature(harness.player2, "Blocker 2", "3", "3")

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attacker, "", []string{}, effects.DurationWhileOnBattlefield)
	harness.engine.layerSystem.AddEffect(harness.gameID, effect)

	// Setup combat
	harness.setupCombat(harness.player1)

	// Declare attacker
	err := harness.engine.DeclareAttacker(harness.gameID, attacker, harness.player2, harness.player1)
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = harness.engine.processMustBeBlockedRequirements(harness.gameState)
	require.NoError(t, err)

	// Declare blockers (forced by lure)
	err = harness.engine.DeclareBlocker(harness.gameID, blocker1, attacker, harness.player2)
	require.NoError(t, err)
	err = harness.engine.DeclareBlocker(harness.gameID, blocker2, attacker, harness.player2)
	require.NoError(t, err)

	// Verify blocking setup
	assert.True(t, harness.gameState.cards[blocker1].Blocking, "blocker 1 should be blocking")
	assert.True(t, harness.gameState.cards[blocker2].Blocking, "blocker 2 should be blocking")

	// Assign and apply damage
	err = harness.engine.AssignCombatDamage(harness.gameID, false)
	require.NoError(t, err)
	err = harness.engine.ApplyCombatDamage(harness.gameID)
	require.NoError(t, err)

	// Verify damage (5 damage split between blockers, each taking their toughness)
	assert.Equal(t, 2, harness.gameState.cards[blocker1].Damage, "blocker 1 should take lethal damage")
	assert.Equal(t, 3, harness.gameState.cards[blocker2].Damage, "blocker 2 should take lethal damage")
	assert.Equal(t, 5, harness.gameState.cards[attacker].Damage, "attacker should take 5 damage from both blockers")
}
