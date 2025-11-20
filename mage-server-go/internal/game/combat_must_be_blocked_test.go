package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMustBeBlockedIfAble tests the "must be blocked if able" combat mechanic (lure)
func TestMustBeBlockedIfAble(t *testing.T) {
	h := NewCombatTestHarness(t, "game-1", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create creatures
	attackerID := h.CreateAttacker("attacker", "Attacker with Lure", "Alice", "3", "3")
	blockerID := h.CreateBlocker("blocker", "Blocker", "Bob", "2", "2")

	// Add "must be blocked if able" effect (Lure)
	effect := effects.NewMustBeBlockedEffect(attackerID, "", []string{}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat and declare attacker
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = h.engine.processMustBeBlockedRequirements(gameState)
	require.NoError(t, err)

	// Verify blocker is required to block the attacker
	requirements, ok := gameState.combat.creatureMustBlockAttackers[blockerID]
	require.True(t, ok, "blocker should have blocking requirements")
	assert.True(t, requirements[attackerID], "blocker must block the attacker")
}

// TestMustBeBlockedIfAble_MultipleBlockers tests lure with multiple potential blockers
func TestMustBeBlockedIfAble_MultipleBlockers(t *testing.T) {
	h := NewCombatTestHarness(t, "game-2", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create creatures
	attackerID := h.CreateAttacker("attacker", "Attacker with Lure", "Alice", "4", "4")
	blocker1 := h.CreateBlocker("blocker1", "Blocker 1", "Bob", "2", "2")
	blocker2 := h.CreateBlocker("blocker2", "Blocker 2", "Bob", "3", "3")
	blocker3 := h.CreateBlocker("blocker3", "Blocker 3", "Bob", "1", "1")

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attackerID, "", []string{}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat and declare attacker
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = h.engine.processMustBeBlockedRequirements(gameState)
	require.NoError(t, err)

	// Verify all able blockers are required to block
	assert.True(t, gameState.combat.creatureMustBlockAttackers[blocker1][attackerID], "blocker 1 must block")
	assert.True(t, gameState.combat.creatureMustBlockAttackers[blocker2][attackerID], "blocker 2 must block")
	assert.True(t, gameState.combat.creatureMustBlockAttackers[blocker3][attackerID], "blocker 3 must block")
}

// TestMustBeBlockedIfAble_TappedBlocker tests that tapped creatures can't block even with lure
func TestMustBeBlockedIfAble_TappedBlocker(t *testing.T) {
	h := NewCombatTestHarness(t, "game-3", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create creatures
	attackerID := h.CreateAttacker("attacker", "Attacker with Lure", "Alice", "3", "3")
	blockerID := h.CreateCreature(CreatureSpec{
		ID:         "blocker",
		Name:       "Tapped Blocker",
		Power:      "2",
		Toughness:  "2",
		Controller: "Bob",
		Tapped:     true,
	})

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attackerID, "", []string{}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat and declare attacker
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = h.engine.processMustBeBlockedRequirements(gameState)
	require.NoError(t, err)

	// Verify tapped blocker is NOT required to block (can't block because tapped)
	_, hasRequirement := gameState.combat.creatureMustBlockAttackers[blockerID]
	assert.False(t, hasRequirement, "tapped blocker should not be required to block")
}

// TestMustBeBlockedIfAble_FlyingAttacker tests lure with flying (only reach/flying can block)
func TestMustBeBlockedIfAble_FlyingAttacker(t *testing.T) {
	h := NewCombatTestHarness(t, "game-4", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create creatures
	attackerID := h.CreateCreature(CreatureSpec{
		ID:         "attacker",
		Name:       "Flying Attacker with Lure",
		Power:      "3",
		Toughness:  "3",
		Controller: "Alice",
		Abilities:  []string{abilityFlying},
	})

	groundBlockerID := h.CreateBlocker("ground_blocker", "Ground Blocker", "Bob", "2", "2")
	flyingBlockerID := h.CreateCreature(CreatureSpec{
		ID:         "flying_blocker",
		Name:       "Flying Blocker",
		Power:      "2",
		Toughness:  "2",
		Controller: "Bob",
		Abilities:  []string{abilityFlying},
	})

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attackerID, "", []string{}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat and declare attacker
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = h.engine.processMustBeBlockedRequirements(gameState)
	require.NoError(t, err)

	// Verify only flying blocker is required to block
	_, groundHasRequirement := gameState.combat.creatureMustBlockAttackers[groundBlockerID]
	flyingHasRequirement := gameState.combat.creatureMustBlockAttackers[flyingBlockerID][attackerID]

	assert.False(t, groundHasRequirement, "ground blocker cannot block flying attacker")
	assert.True(t, flyingHasRequirement, "flying blocker must block flying attacker")
}

// TestMustBeBlockedIfAble_WithMenace tests lure interaction with menace
func TestMustBeBlockedIfAble_WithMenace(t *testing.T) {
	h := NewCombatTestHarness(t, "game-5", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create creatures
	attackerID := h.CreateCreature(CreatureSpec{
		ID:         "attacker",
		Name:       "Attacker with Lure and Menace",
		Power:      "4",
		Toughness:  "4",
		Controller: "Alice",
		Abilities:  []string{abilityMenace},
	})

	blocker1 := h.CreateBlocker("blocker1", "Blocker 1", "Bob", "2", "2")
	blocker2 := h.CreateBlocker("blocker2", "Blocker 2", "Bob", "2", "2")

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attackerID, "", []string{}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat and declare attacker
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)

	// Set menace requirement (minimum 2 blockers)
	gameState.combat.minBlockersPerAttacker[attackerID] = 2

	// Process must-be-blocked requirements
	err = h.engine.processMustBeBlockedRequirements(gameState)
	require.NoError(t, err)

	// Verify both blockers are required to block (to satisfy menace)
	assert.True(t, gameState.combat.creatureMustBlockAttackers[blocker1][attackerID], "blocker 1 must block")
	assert.True(t, gameState.combat.creatureMustBlockAttackers[blocker2][attackerID], "blocker 2 must block")
}

// TestMustBeBlockedIfAble_SpecificBlockers tests lure targeting specific blockers
func TestMustBeBlockedIfAble_SpecificBlockers(t *testing.T) {
	h := NewCombatTestHarness(t, "game-6", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create creatures
	attackerID := h.CreateAttacker("attacker", "Attacker", "Alice", "3", "3")
	targetBlockerID := h.CreateBlocker("target_blocker", "Target Blocker", "Bob", "2", "2")
	otherBlockerID := h.CreateBlocker("other_blocker", "Other Blocker", "Bob", "2", "2")

	// Add "must be blocked if able" effect targeting only one blocker
	effect := effects.NewMustBeBlockedEffect(attackerID, "", []string{targetBlockerID}, effects.DurationEndOfTurn)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat and declare attacker
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = h.engine.processMustBeBlockedRequirements(gameState)
	require.NoError(t, err)

	// Verify only target blocker is required to block
	targetHasRequirement := gameState.combat.creatureMustBlockAttackers[targetBlockerID][attackerID]
	_, otherHasRequirement := gameState.combat.creatureMustBlockAttackers[otherBlockerID]

	assert.True(t, targetHasRequirement, "target blocker must block")
	assert.False(t, otherHasRequirement, "other blocker should not be required to block")
}

// TestMustBeBlockedIfAble_CheckBlockRequirements tests validation of blocking requirements
func TestMustBeBlockedIfAble_CheckBlockRequirements(t *testing.T) {
	h := NewCombatTestHarness(t, "game-7", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create creatures
	attackerID := h.CreateAttacker("attacker", "Attacker with Lure", "Alice", "3", "3")
	_ = h.CreateBlocker("blocker", "Blocker", "Bob", "2", "2")

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attackerID, "", []string{}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat and declare attacker
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = h.engine.processMustBeBlockedRequirements(gameState)
	require.NoError(t, err)

	// Check block requirements without declaring blocker
	violations, err := h.engine.CheckBlockRequirements(h.gameID, "Bob")
	require.NoError(t, err)

	// Should have violation because blocker didn't block
	assert.NotEmpty(t, violations, "should have blocking requirement violation")
}

// TestMustBeBlockedIfAble_CheckBlockRequirements_Satisfied tests validation when requirement is met
func TestMustBeBlockedIfAble_CheckBlockRequirements_Satisfied(t *testing.T) {
	h := NewCombatTestHarness(t, "game-8", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create creatures
	attackerID := h.CreateAttacker("attacker", "Attacker with Lure", "Alice", "3", "3")
	blockerID := h.CreateBlocker("blocker", "Blocker", "Bob", "2", "2")

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attackerID, "", []string{}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat and declare attacker
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = h.engine.processMustBeBlockedRequirements(gameState)
	require.NoError(t, err)

	// Declare blocker to satisfy requirement
	err = h.engine.DeclareBlocker(h.gameID, blockerID, attackerID, "Bob")
	require.NoError(t, err)

	// Check block requirements after declaring blocker
	violations, err := h.engine.CheckBlockRequirements(h.gameID, "Bob")
	require.NoError(t, err)

	// Should have no violations because blocker blocked
	assert.Empty(t, violations, "should have no blocking requirement violations")
}

// TestMustBeBlockedIfAble_Integration tests full combat flow with lure
func TestMustBeBlockedIfAble_Integration(t *testing.T) {
	h := NewCombatTestHarness(t, "game-9", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create creatures
	attackerID := h.CreateAttacker("attacker", "Attacker with Lure", "Alice", "5", "5")
	blocker1 := h.CreateBlocker("blocker1", "Blocker 1", "Bob", "2", "2")
	blocker2 := h.CreateBlocker("blocker2", "Blocker 2", "Bob", "3", "3")

	// Add "must be blocked if able" effect
	effect := effects.NewMustBeBlockedEffect(attackerID, "", []string{}, effects.DurationWhileOnBattlefield)
	gameState.layerSystem.AddEffect(effect)

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Declare attacker
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)

	// Process must-be-blocked requirements
	err = h.engine.processMustBeBlockedRequirements(gameState)
	require.NoError(t, err)

	// Declare blockers (forced by lure)
	err = h.engine.DeclareBlocker(h.gameID, blocker1, attackerID, "Bob")
	require.NoError(t, err)
	err = h.engine.DeclareBlocker(h.gameID, blocker2, attackerID, "Bob")
	require.NoError(t, err)

	// Verify blocking setup
	assert.True(t, gameState.cards[blocker1].Blocking, "blocker 1 should be blocking")
	assert.True(t, gameState.cards[blocker2].Blocking, "blocker 2 should be blocking")

	// Assign and apply damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify damage (5 damage split between blockers, each taking their toughness)
	assert.Equal(t, 2, gameState.cards[blocker1].Damage, "blocker 1 should take lethal damage")
	assert.Equal(t, 3, gameState.cards[blocker2].Damage, "blocker 2 should take lethal damage")
	assert.Equal(t, 5, gameState.cards[attackerID].Damage, "attacker should take 5 damage from both blockers")
}
