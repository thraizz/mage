package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPlaneswalkerCombat_BasicAttack tests basic attack on a planeswalker
func TestPlaneswalkerCombat_BasicAttack(t *testing.T) {
	h := NewCombatTestHarness(t, "game-1", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create a planeswalker controlled by Bob with 4 loyalty
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

	// Create a 2/2 attacker for Alice
	attackerID := h.CreateAttacker("attacker", "Grizzly Bears", "Alice", "2", "2")

	// Setup combat with Alice attacking
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Verify planeswalker is in defenders
	assert.True(t, gameState.combat.defenders[planeswalkerID], "planeswalker should be in defenders")

	// Declare attacker targeting the planeswalker
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Verify attack declaration
	attacker := gameState.cards[attackerID]
	assert.True(t, attacker.Attacking, "creature should be attacking")
	assert.Equal(t, planeswalkerID, attacker.AttackingWhat, "creature should be attacking planeswalker")

	// Assign and apply combat damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify planeswalker loyalty reduced by 2
	assert.Equal(t, 2, planeswalker.Counters.GetCount("loyalty"), "planeswalker should have 2 loyalty (4-2)")
}

// TestPlaneswalkerCombat_Death tests planeswalker with 0 loyalty
func TestPlaneswalkerCombat_Death(t *testing.T) {
	h := NewCombatTestHarness(t, "game-2", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create a planeswalker with 3 loyalty
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

	// Create a 5/5 attacker (deals lethal to planeswalker)
	attackerID := h.CreateAttacker("attacker", "Serra Angel", "Alice", "5", "5")

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Attack planeswalker
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Apply damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify planeswalker has 0 loyalty (3 - 5 = 0, minimum)
	// State-based actions would move it to graveyard (Rule 306.9)
	loyalty := planeswalker.Counters.GetCount("loyalty")
	assert.Equal(t, 0, loyalty, "planeswalker should have 0 loyalty")
}

// TestPlaneswalkerCombat_Lifelink tests lifelink with planeswalker damage
func TestPlaneswalkerCombat_Lifelink(t *testing.T) {
	h := NewCombatTestHarness(t, "game-3", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create planeswalker with 5 loyalty
	planeswalkerID := "teferi"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Teferi, Time Raveler",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Loyalty:      "5",
		Counters:     counters.NewCounters(),
	}
	planeswalker.Counters.AddCounter(counters.NewCounter("loyalty", 5))
	gameState.cards[planeswalkerID] = planeswalker

	// Set Alice's life to 15 for easier testing
	gameState.players["Alice"].Life = 15
	gameState.mu.Unlock()

	// Create attacker with lifelink
	attackerID := h.CreateCreature(CreatureSpec{
		ID:         "attacker",
		Name:       "Vampire Nighthawk",
		Power:      "3",
		Toughness:  "3",
		Controller: "Alice",
		Abilities:  []string{abilityLifelink},
	})

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Attack planeswalker
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Apply damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify planeswalker loyalty reduced
	assert.Equal(t, 2, planeswalker.Counters.GetCount("loyalty"), "planeswalker should have 2 loyalty (5-3)")

	// Verify Alice gained 3 life from lifelink
	assert.Equal(t, 18, gameState.players["Alice"].Life, "Alice should have 18 life (15+3 from lifelink)")
}

// TestPlaneswalkerCombat_Deathtouch tests deathtouch with planeswalker
func TestPlaneswalkerCombat_Deathtouch(t *testing.T) {
	h := NewCombatTestHarness(t, "game-4", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create planeswalker with 6 loyalty
	planeswalkerID := "chandra"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Chandra, Torch of Defiance",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Loyalty:      "6",
		Counters:     counters.NewCounters(),
	}
	planeswalker.Counters.AddCounter(counters.NewCounter("loyalty", 6))
	gameState.cards[planeswalkerID] = planeswalker
	gameState.mu.Unlock()

	// Create 1/1 attacker with deathtouch
	attackerID := h.CreateCreature(CreatureSpec{
		ID:         "attacker",
		Name:       "Typhoid Rats",
		Power:      "1",
		Toughness:  "1",
		Controller: "Alice",
		Abilities:  []string{abilityDeathtouch},
	})

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Attack planeswalker
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Verify lethal damage calculation (deathtouch makes 1 damage lethal)
	lethal := h.engine.getLethalDamageWithAttacker(gameState, planeswalker, attackerID)
	assert.Equal(t, 1, lethal, "with deathtouch, 1 damage should be lethal to planeswalker")

	// Apply damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify planeswalker lost 1 loyalty (deathtouch doesn't automatically kill planeswalkers)
	assert.Equal(t, 5, planeswalker.Counters.GetCount("loyalty"), "planeswalker should have 5 loyalty (6-1)")
}

// TestPlaneswalkerCombat_Multiple tests attacking multiple targets including planeswalker
func TestPlaneswalkerCombat_Multiple(t *testing.T) {
	h := NewCombatTestHarness(t, "game-5", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create planeswalker
	planeswalkerID := "garruk"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Garruk Wildspeaker",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Loyalty:      "3",
		Counters:     counters.NewCounters(),
	}
	planeswalker.Counters.AddCounter(counters.NewCounter("loyalty", 3))
	gameState.cards[planeswalkerID] = planeswalker

	// Set Bob's life to 20
	gameState.players["Bob"].Life = 20
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

	// Apply damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify planeswalker took 2 damage
	assert.Equal(t, 1, planeswalker.Counters.GetCount("loyalty"), "planeswalker should have 1 loyalty (3-2)")

	// Verify Bob took 3 damage
	assert.Equal(t, 17, gameState.players["Bob"].Life, "Bob should have 17 life (20-3)")
}

// TestPlaneswalkerCombat_CannotAttackOwnPlaneswalker tests that you can't attack your own planeswalker
func TestPlaneswalkerCombat_CannotAttackOwnPlaneswalker(t *testing.T) {
	h := NewCombatTestHarness(t, "game-6", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create planeswalker controlled by Alice
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

	// Create attacker for Alice
	attackerID := h.CreateAttacker("attacker", "Elf", "Alice", "2", "2")

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Verify Alice's own planeswalker is NOT in defenders
	assert.False(t, gameState.combat.defenders[planeswalkerID], "Alice's own planeswalker should not be in defenders")

	// Attempt to attack own planeswalker should fail
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	assert.Error(t, err, "should not be able to attack own planeswalker")
}

// TestPlaneswalkerCombat_WithBlocker tests blocking for planeswalker
func TestPlaneswalkerCombat_WithBlocker(t *testing.T) {
	h := NewCombatTestHarness(t, "game-7", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create planeswalker
	planeswalkerID := "ajani"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Ajani, Strength of the Pride",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Loyalty:      "5",
		Counters:     counters.NewCounters(),
	}
	planeswalker.Counters.AddCounter(counters.NewCounter("loyalty", 5))
	gameState.cards[planeswalkerID] = planeswalker
	gameState.mu.Unlock()

	// Create attacker and blocker
	attackerID := h.CreateAttacker("attacker", "Dragon", "Alice", "4", "4")
	blockerID := h.CreateBlocker("blocker", "Angel", "Bob", "3", "3")

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Attack planeswalker
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Block the attacker
	err = h.engine.DeclareBlocker(h.gameID, blockerID, attackerID, "Bob")
	require.NoError(t, err)

	// Apply damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify planeswalker took no damage (was blocked)
	assert.Equal(t, 5, planeswalker.Counters.GetCount("loyalty"), "planeswalker should still have 5 loyalty (blocked)")

	// Verify blocker and attacker damaged each other
	attacker := gameState.cards[attackerID]
	blocker := gameState.cards[blockerID]
	assert.Equal(t, 3, attacker.Damage, "attacker should have 3 damage")
	assert.Equal(t, 4, blocker.Damage, "blocker should have 4 damage (lethal)")
}

// TestPlaneswalkerCombat_TrampleOver tests trample over planeswalkers ability
func TestPlaneswalkerCombat_TrampleOver(t *testing.T) {
	h := NewCombatTestHarness(t, "game-8", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create planeswalker with 4 loyalty
	planeswalkerID := "garruk"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Garruk Wildspeaker",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Loyalty:      "4",
		Counters:     counters.NewCounters(),
	}
	planeswalker.Counters.AddCounter(counters.NewCounter("loyalty", 4))
	gameState.cards[planeswalkerID] = planeswalker

	// Set Bob's life to 20
	gameState.players["Bob"].Life = 20
	gameState.mu.Unlock()

	// Create 7/7 attacker with trample over planeswalkers (like Thrasta)
	attackerID := h.CreateCreature(CreatureSpec{
		ID:         "attacker",
		Name:       "Thrasta, Tempest's Roar",
		Power:      "7",
		Toughness:  "7",
		Controller: "Alice",
		Abilities:  []string{abilityTrample, abilityTrampleOverPlaneswalkers},
	})

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Attack planeswalker
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Apply damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify planeswalker took 4 damage (all loyalty counters)
	assert.Equal(t, 0, planeswalker.Counters.GetCount("loyalty"), "planeswalker should have 0 loyalty")

	// Verify Bob took 3 excess damage (7 power - 4 loyalty = 3)
	assert.Equal(t, 17, gameState.players["Bob"].Life, "Bob should have 17 life (20-3 from trample over)")
}

// TestPlaneswalkerCombat_TrampleOverLessThanLethal tests trample over planeswalkers with less than lethal damage
func TestPlaneswalkerCombat_TrampleOverLessThanLethal(t *testing.T) {
	h := NewCombatTestHarness(t, "game-9", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create planeswalker with 5 loyalty
	planeswalkerID := "jace"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Jace, the Mind Sculptor",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Loyalty:      "5",
		Counters:     counters.NewCounters(),
	}
	planeswalker.Counters.AddCounter(counters.NewCounter("loyalty", 5))
	gameState.cards[planeswalkerID] = planeswalker

	// Set Bob's life to 20
	gameState.players["Bob"].Life = 20
	gameState.mu.Unlock()

	// Create 3/3 attacker with trample over planeswalkers (less than lethal)
	attackerID := h.CreateCreature(CreatureSpec{
		ID:         "attacker",
		Name:       "Trampler",
		Power:      "3",
		Toughness:  "3",
		Controller: "Alice",
		Abilities:  []string{abilityTrampleOverPlaneswalkers},
	})

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Attack planeswalker
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Apply damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify planeswalker took 3 damage (not lethal)
	assert.Equal(t, 2, planeswalker.Counters.GetCount("loyalty"), "planeswalker should have 2 loyalty (5-3)")

	// Verify Bob took NO damage (not enough to kill planeswalker)
	assert.Equal(t, 20, gameState.players["Bob"].Life, "Bob should still have 20 life (no excess)")
}

// TestPlaneswalkerCombat_TrampleOverWithLifelink tests trample over planeswalkers with lifelink
func TestPlaneswalkerCombat_TrampleOverWithLifelink(t *testing.T) {
	h := NewCombatTestHarness(t, "game-10", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create planeswalker with 3 loyalty
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

	// Set players' life
	gameState.players["Alice"].Life = 15
	gameState.players["Bob"].Life = 20
	gameState.mu.Unlock()

	// Create 6/6 attacker with trample over planeswalkers and lifelink
	attackerID := h.CreateCreature(CreatureSpec{
		ID:         "attacker",
		Name:       "Lifelink Trampler",
		Power:      "6",
		Toughness:  "6",
		Controller: "Alice",
		Abilities:  []string{abilityTrampleOverPlaneswalkers, abilityLifelink},
	})

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Attack planeswalker
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Apply damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify planeswalker took 3 damage
	assert.Equal(t, 0, planeswalker.Counters.GetCount("loyalty"), "planeswalker should have 0 loyalty")

	// Verify Bob took 3 excess damage
	assert.Equal(t, 17, gameState.players["Bob"].Life, "Bob should have 17 life (20-3)")

	// Verify Alice gained 6 life from lifelink (full damage dealt, both to planeswalker and player)
	assert.Equal(t, 21, gameState.players["Alice"].Life, "Alice should have 21 life (15+6 from lifelink)")
}

// TestPlaneswalkerCombat_RegularTrampleDoesNotCarryOver tests that regular trample doesn't work on planeswalkers
func TestPlaneswalkerCombat_RegularTrampleDoesNotCarryOver(t *testing.T) {
	h := NewCombatTestHarness(t, "game-11", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create planeswalker with 3 loyalty
	planeswalkerID := "chandra"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Chandra, Torch of Defiance",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Loyalty:      "3",
		Counters:     counters.NewCounters(),
	}
	planeswalker.Counters.AddCounter(counters.NewCounter("loyalty", 3))
	gameState.cards[planeswalkerID] = planeswalker

	// Set Bob's life to 20
	gameState.players["Bob"].Life = 20
	gameState.mu.Unlock()

	// Create 7/7 attacker with ONLY regular trample (not trample over planeswalkers)
	attackerID := h.CreateCreature(CreatureSpec{
		ID:         "attacker",
		Name:       "Colossal Dreadmaw",
		Power:      "7",
		Toughness:  "7",
		Controller: "Alice",
		Abilities:  []string{abilityTrample}, // Only regular trample
	})

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Attack planeswalker
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Apply damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify planeswalker took only 3 damage (reduced to 0 loyalty)
	assert.Equal(t, 0, planeswalker.Counters.GetCount("loyalty"), "planeswalker should have 0 loyalty")

	// Verify Bob took NO damage (regular trample doesn't carry over to planeswalker controller)
	assert.Equal(t, 20, gameState.players["Bob"].Life, "Bob should still have 20 life (regular trample doesn't work)")
}

// TestPlaneswalkerCombat_TrampleOverWithDeathtouch tests trample over planeswalkers with deathtouch
func TestPlaneswalkerCombat_TrampleOverWithDeathtouch(t *testing.T) {
	h := NewCombatTestHarness(t, "game-12", []string{"Alice", "Bob"})
	gameState := h.GetGameState()

	// Create planeswalker with 5 loyalty
	planeswalkerID := "teferi"
	gameState.mu.Lock()
	planeswalker := &internalCard{
		ID:           planeswalkerID,
		Name:         "Teferi, Time Raveler",
		Type:         "Planeswalker",
		Zone:         zoneBattlefield,
		OwnerID:      "Bob",
		ControllerID: "Bob",
		Loyalty:      "5",
		Counters:     counters.NewCounters(),
	}
	planeswalker.Counters.AddCounter(counters.NewCounter("loyalty", 5))
	gameState.cards[planeswalkerID] = planeswalker

	// Set Bob's life to 20
	gameState.players["Bob"].Life = 20
	gameState.mu.Unlock()

	// Create 3/3 attacker with trample over planeswalkers and deathtouch
	attackerID := h.CreateCreature(CreatureSpec{
		ID:         "attacker",
		Name:       "Deathtouch Trampler",
		Power:      "3",
		Toughness:  "3",
		Controller: "Alice",
		Abilities:  []string{abilityTrampleOverPlaneswalkers, abilityDeathtouch},
	})

	// Setup combat
	err := h.engine.SetAttacker(h.gameID, "Alice")
	require.NoError(t, err)
	err = h.engine.SetDefenders(h.gameID)
	require.NoError(t, err)

	// Attack planeswalker
	err = h.engine.DeclareAttacker(h.gameID, attackerID, planeswalkerID, "Alice")
	require.NoError(t, err)

	// Verify lethal damage calculation with deathtouch
	lethal := h.engine.getLethalDamageWithAttacker(gameState, planeswalker, attackerID)
	assert.Equal(t, 1, lethal, "with deathtouch, 1 damage should be lethal to planeswalker")

	// Apply damage
	err = h.engine.AssignCombatDamage(h.gameID, false)
	require.NoError(t, err)
	err = h.engine.ApplyCombatDamage(h.gameID)
	require.NoError(t, err)

	// Verify planeswalker took 1 damage (lethal with deathtouch)
	assert.Equal(t, 4, planeswalker.Counters.GetCount("loyalty"), "planeswalker should have 4 loyalty (5-1)")

	// Verify Bob took 2 excess damage (3 power - 1 lethal = 2 trample over)
	assert.Equal(t, 18, gameState.players["Bob"].Life, "Bob should have 18 life (20-2 from trample over)")
}
