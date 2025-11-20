package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAttackTracking_DirectPlayerAttack tests tracking when a player is attacked directly
func TestAttackTracking_DirectPlayerAttack(t *testing.T) {
	h := NewCombatTestHarness(t, "game-1", []string{"Alice", "Bob"})

	// Create attacker
	attackerID := h.CreateAttacker("attacker", "Grizzly Bears", "Alice", "2", "2")

	// Setup combat and declare attacker targeting Bob
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, "Bob", "Alice")
	require.NoError(t, err)

	// Verify attack is tracked
	attacked, err := h.engine.HasPlayerAttackedPlayerOrPlaneswalker(h.gameID, "Alice", "Bob")
	require.NoError(t, err)
	assert.True(t, attacked, "Alice should have attacked Bob")

	// Verify other player not attacked
	attacked, err = h.engine.HasPlayerAttackedPlayerOrPlaneswalker(h.gameID, "Bob", "Alice")
	require.NoError(t, err)
	assert.False(t, attacked, "Bob should not have attacked Alice")
}

// TestAttackTracking_PlaneswalkerAttack tests tracking when a planeswalker is attacked
func TestAttackTracking_PlaneswalkerAttack(t *testing.T) {
	h := NewCombatTestHarness(t, "game-2", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create a planeswalker controlled by Bob
	planeswalkerID := "jace"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Jace, the Mind Sculptor",
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

	// Create attacker
	attackerID := h.CreateAttacker("attacker", "Grizzly Bears", "Alice", "2", "2")

	// Setup combat and declare attacker targeting the planeswalker
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Verify planeswalker attack is tracked as attacking Bob (the controller)
	attacked, err := h.engine.HasPlayerAttackedPlayerOrPlaneswalker(h.gameID, "Alice", "Bob")
	require.NoError(t, err)
	assert.True(t, attacked, "Alice should have attacked Bob's planeswalker")
}

// TestAttackTracking_BothPlayerAndPlaneswalker tests attacking both a player and their planeswalker
func TestAttackTracking_BothPlayerAndPlaneswalker(t *testing.T) {
	h := NewCombatTestHarness(t, "game-3", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create a planeswalker controlled by Bob
	planeswalkerID := "liliana"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Liliana of the Veil",
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

	// Create two attackers
	attacker1 := h.CreateAttacker("attacker1", "Bear", "Alice", "2", "2")
	attacker2 := h.CreateAttacker("attacker2", "Wolf", "Alice", "3", "3")

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Attack planeswalker with first creature
	err = h.engine.DeclareAttacker(h.gameID, attacker1, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Attack player with second creature
	err = h.engine.DeclareAttacker(h.gameID, attacker2, "Bob", "Alice")
	require.NoError(t, err)

	// Verify both attacks are tracked
	attacked, err := h.engine.HasPlayerAttackedPlayerOrPlaneswalker(h.gameID, "Alice", "Bob")
	require.NoError(t, err)
	assert.True(t, attacked, "Alice should have attacked Bob (both direct and planeswalker)")
}

// TestAttackTracking_NotAttackingOwnPlaneswalker tests that attacking your own planeswalker doesn't count
func TestAttackTracking_NotAttackingOwnPlaneswalker(t *testing.T) {
	h := NewCombatTestHarness(t, "game-4", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create a planeswalker controlled by Alice (can't be attacked by Alice)
	planeswalkerID := "nissa"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Nissa, Who Shakes the World",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Alice",
		ControllerID: "Alice",
		Loyalty:      "5",
		Counters:     counters.NewCounters(),
	}
	planeswalker.Counters.AddCounter(counters.NewCounter("loyalty", 5))
	gameState.cards[planeswalkerID] = planeswalker
	gameState.mu.Unlock()

	// Create attacker
	h.CreateAttacker("attacker", "Elf", "Alice", "2", "2")

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Verify Alice's own planeswalker is NOT in defenders
	assert.False(t, gameState.combat.defenders[planeswalkerID], "Alice's own planeswalker should not be in defenders")

	// Verify Alice has not attacked Alice (can't attack own planeswalker)
	attacked, err := h.engine.HasPlayerAttackedPlayerOrPlaneswalker(h.gameID, "Alice", "Alice")
	require.NoError(t, err)
	assert.False(t, attacked, "Alice should not have attacked herself")
}

// TestAttackTracking_MultiplePlaneswalkers tests attacking multiple planeswalkers of the same controller
func TestAttackTracking_MultiplePlaneswalkers(t *testing.T) {
	h := NewCombatTestHarness(t, "game-5", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create two planeswalkers controlled by Bob
	pw1ID := "jace"
	pw2ID := "teferi"
	gameState.mu.Lock()
	pw1 := &internalCard{
		ID:           pw1ID,
		Name:         "Jace",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Loyalty:      "4",
		Counters:     counters.NewCounters(),
	}
	pw1.Counters.AddCounter(counters.NewCounter("loyalty", 4))
	gameState.cards[pw1ID] = pw1

	pw2 := &internalCard{
		ID:           pw2ID,
		Name:         "Teferi",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Loyalty:      "5",
		Counters:     counters.NewCounters(),
	}
	pw2.Counters.AddCounter(counters.NewCounter("loyalty", 5))
	gameState.cards[pw2ID] = pw2
	gameState.mu.Unlock()

	// Create two attackers
	attacker1 := h.CreateAttacker("attacker1", "Bear", "Alice", "2", "2")
	attacker2 := h.CreateAttacker("attacker2", "Wolf", "Alice", "3", "3")

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Attack both planeswalkers
	err = h.engine.DeclareAttacker(h.gameID, attacker1, pw1ID, "Alice")
	require.NoError(t, err)
	err = h.engine.DeclareAttacker(h.gameID, attacker2, pw2ID, "Alice")
	require.NoError(t, err)

	// Verify attacks are tracked (both count as attacking Bob)
	attacked, err := h.engine.HasPlayerAttackedPlayerOrPlaneswalker(h.gameID, "Alice", "Bob")
	require.NoError(t, err)
	assert.True(t, attacked, "Alice should have attacked Bob's planeswalkers")
}

// TestAttackTracking_NoAttacksThisTurn tests the default case when no attacks have occurred
func TestAttackTracking_NoAttacksThisTurn(t *testing.T) {
	h := NewCombatTestHarness(t, "game-6", []string{"Alice", "Bob"})

	// Create attacker but don't declare it
	h.CreateAttacker("attacker", "Grizzly Bears", "Alice", "2", "2")

	// Setup combat but don't declare any attackers
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Verify no attacks tracked
	attacked, err := h.engine.HasPlayerAttackedPlayerOrPlaneswalker(h.gameID, "Alice", "Bob")
	require.NoError(t, err)
	assert.False(t, attacked, "Alice should not have attacked Bob (no attacks declared)")
}
