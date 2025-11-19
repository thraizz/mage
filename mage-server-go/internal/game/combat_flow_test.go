package game

import (
	"testing"
)

// TestCombatFlow_SingleAttackerNoBlockers tests basic combat with one attacker and no blockers
// Per GO_PORT_TASKS.md: "Test single attacker, no blockers (damage to player)"
func TestCombatFlow_SingleAttackerNoBlockers(t *testing.T) {
	h := NewCombatTestHarness(t, "test-single-attacker", []string{"Alice", "Bob"})

	// Create a 3/3 attacker for Alice
	attacker := h.CreateAttacker("attacker-1", "Trained Armodon", "Alice", "3", "3")

	// Get initial life totals
	initialBobLife := h.GetPlayerLife("Bob")

	// Run combat: attacker attacks Bob, no blockers
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, nil)

	// Verify results
	h.AssertPlayerLife("Bob", initialBobLife-3) // Bob should lose 3 life
	h.AssertCreatureDamage(attacker, 0)          // Attacker should have no damage
	h.AssertCreatureAlive(attacker)              // Attacker should still be alive
	h.AssertCreatureTapped(attacker, true)       // Attacker should be tapped
}

// TestCombatFlow_SingleAttackerSingleBlocker tests combat with one attacker blocked by one blocker
// Per GO_PORT_TASKS.md: "Test single attacker, single blocker (damage to creatures)"
func TestCombatFlow_SingleAttackerSingleBlocker(t *testing.T) {
	h := NewCombatTestHarness(t, "test-single-blocker", []string{"Alice", "Bob"})

	// Create a 3/3 attacker for Alice
	attacker := h.CreateAttacker("attacker-1", "Trained Armodon", "Alice", "3", "3")

	// Create a 2/2 blocker for Bob
	blocker := h.CreateBlocker("blocker-1", "Grizzly Bears", "Bob", "2", "2")

	// Get initial life totals
	initialBobLife := h.GetPlayerLife("Bob")

	// Run combat: attacker attacks Bob, blocker blocks
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, map[string]string{
		blocker: attacker,
	})

	// Verify results
	h.AssertPlayerLife("Bob", initialBobLife)   // Bob should take no damage (blocked)
	// Note: damage is cleared after EndCombat, so we can't check attacker damage here
	h.AssertCreatureDead(blocker)               // Blocker should be dead (2 toughness, 3 damage)
	h.AssertCreatureAlive(attacker)             // Attacker should still be alive (3 toughness, 2 damage)
}

// TestCombatFlow_MultipleAttackersNoBlockers tests multiple attackers with no blockers
// Per GO_PORT_TASKS.md: "Test multiple attackers, no blockers"
func TestCombatFlow_MultipleAttackersNoBlockers(t *testing.T) {
	h := NewCombatTestHarness(t, "test-multiple-attackers", []string{"Alice", "Bob"})

	// Create three attackers for Alice
	attacker1 := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")
	attacker2 := h.CreateAttacker("attacker-2", "Trained Armodon", "Alice", "3", "3")
	attacker3 := h.CreateAttacker("attacker-3", "Hill Giant", "Alice", "3", "3")

	// Get initial life totals
	initialBobLife := h.GetPlayerLife("Bob")

	// Run combat: all three attack Bob, no blockers
	h.RunFullCombat("Alice", map[string]string{
		attacker1: "Bob",
		attacker2: "Bob",
		attacker3: "Bob",
	}, nil)

	// Verify results
	expectedDamage := 2 + 3 + 3 // Total 8 damage
	h.AssertPlayerLife("Bob", initialBobLife-expectedDamage)
	h.AssertCreatureAlive(attacker1)
	h.AssertCreatureAlive(attacker2)
	h.AssertCreatureAlive(attacker3)
}

// TestCombatFlow_MultipleAttackersMultipleBlockers tests multiple attackers and blockers
// Per GO_PORT_TASKS.md: "Test multiple attackers, multiple blockers"
func TestCombatFlow_MultipleAttackersMultipleBlockers(t *testing.T) {
	h := NewCombatTestHarness(t, "test-multiple-blockers", []string{"Alice", "Bob"})

	// Create three 3/3 attackers for Alice
	attacker1 := h.CreateAttacker("attacker-1", "Trained Armodon", "Alice", "3", "3")
	attacker2 := h.CreateAttacker("attacker-2", "Hill Giant", "Alice", "3", "3")
	attacker3 := h.CreateAttacker("attacker-3", "Centaur Courser", "Alice", "3", "3")

	// Create two 2/2 blockers for Bob
	blocker1 := h.CreateBlocker("blocker-1", "Grizzly Bears", "Bob", "2", "2")
	blocker2 := h.CreateBlocker("blocker-2", "Grizzly Bears", "Bob", "2", "2")

	// Get initial life totals
	initialBobLife := h.GetPlayerLife("Bob")

	// Run combat: three attackers, two blockers
	h.RunFullCombat("Alice", map[string]string{
		attacker1: "Bob",
		attacker2: "Bob",
		attacker3: "Bob",
	}, map[string]string{
		blocker1: attacker1, // Blocks attacker1
		blocker2: attacker2, // Blocks attacker2
		// attacker3 is unblocked
	})

	// Verify results
	h.AssertPlayerLife("Bob", initialBobLife-3) // Only attacker3 (unblocked) deals damage
	h.AssertCreatureDead(blocker1)              // Blocker1 dies (2 toughness, 3 damage)
	h.AssertCreatureDead(blocker2)              // Blocker2 dies (2 toughness, 3 damage)
	// Note: damage is cleared after EndCombat
	h.AssertCreatureAlive(attacker1)
	h.AssertCreatureAlive(attacker2)
	h.AssertCreatureAlive(attacker3)
}

// TestCombatFlow_LethalDamage tests creatures dying from lethal combat damage
// Per GO_PORT_TASKS.md: "Test creature death from combat damage"
func TestCombatFlow_LethalDamage(t *testing.T) {
	h := NewCombatTestHarness(t, "test-lethal-damage", []string{"Alice", "Bob"})

	// Create a 5/5 attacker for Alice
	attacker := h.CreateAttacker("attacker-1", "Craw Wurm", "Alice", "6", "4")

	// Create a 4/4 blocker for Bob
	blocker := h.CreateBlocker("blocker-1", "Hill Giant", "Bob", "3", "3")

	// Run combat
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, map[string]string{
		blocker: attacker,
	})

	// Verify results
	h.AssertCreatureDead(blocker)  // Blocker dies (3 toughness, 6 damage)
	h.AssertCreatureAlive(attacker) // Attacker survives (4 toughness, 3 damage)
}

// TestCombatFlow_MutualDestruction tests both attacker and blocker dying
func TestCombatFlow_MutualDestruction(t *testing.T) {
	h := NewCombatTestHarness(t, "test-mutual-destruction", []string{"Alice", "Bob"})

	// Create a 3/3 attacker for Alice
	attacker := h.CreateAttacker("attacker-1", "Hill Giant", "Alice", "3", "3")

	// Create a 3/3 blocker for Bob
	blocker := h.CreateBlocker("blocker-1", "Hill Giant", "Bob", "3", "3")

	// Run combat
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, map[string]string{
		blocker: attacker,
	})

	// Verify both creatures die
	h.AssertCreatureDead(attacker) // Both have 3 toughness and take 3 damage
	h.AssertCreatureDead(blocker)
}

// TestCombatFlow_PlayerDeath tests player losing from combat damage
// Per GO_PORT_TASKS.md: "Test player death from combat damage"
func TestCombatFlow_PlayerDeath(t *testing.T) {
	h := NewCombatTestHarness(t, "test-player-death", []string{"Alice", "Bob"})

	// Create a massive attacker for Alice
	attacker := h.CreateAttacker("attacker-1", "Serra Angel", "Alice", "25", "25")

	// Get Bob's initial life (should be 20)
	initialBobLife := h.GetPlayerLife("Bob")

	// Run combat
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, nil)

	// Verify Bob lost life (test doesn't verify game-over state, just damage)
	expectedLife := initialBobLife - 25
	if expectedLife < 0 {
		expectedLife = 0 // Life can't go below 0
	}

	actualLife := h.GetPlayerLife("Bob")
	if actualLife > 0 || expectedLife > 0 {
		h.AssertPlayerLife("Bob", expectedLife)
	}
}

// TestCombatFlow_NoValidAttackers tests combat when player has no creatures that can attack
func TestCombatFlow_NoValidAttackers(t *testing.T) {
	h := NewCombatTestHarness(t, "test-no-attackers", []string{"Alice", "Bob"})

	// Create tapped creatures that can't attack
	h.CreateCreature(CreatureSpec{
		ID:         "creature-1",
		Name:       "Tapped Bear",
		Power:      "2",
		Toughness:  "2",
		Controller: "Alice",
		Tapped:     true,
	})

	// Setup combat - should succeed even with no valid attackers
	h.SetupCombat("Alice")

	// No attackers declared, just end combat
	h.EndCombat()

	// Bob should take no damage
	initialLife := 20 // Default starting life
	h.AssertPlayerLife("Bob", initialLife)
}

// TestCombatFlow_CombatCleanup tests that combat state is properly cleaned up
func TestCombatFlow_CombatCleanup(t *testing.T) {
	h := NewCombatTestHarness(t, "test-combat-cleanup", []string{"Alice", "Bob"})

	// Create an attacker
	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")

	// Run combat
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, nil)

	// Verify combat cleanup
	gameState := h.GetGameState()
	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	attackerCard := gameState.cards[attacker]

	// Creature should no longer be attacking
	if attackerCard.Attacking {
		t.Error("creature should not be attacking after combat ends")
	}

	// AttackingWhat should be cleared
	if attackerCard.AttackingWhat != "" {
		t.Error("AttackingWhat should be cleared after combat ends")
	}

	// Damage should be cleared
	if attackerCard.Damage != 0 {
		t.Error("damage should be cleared after combat ends")
	}

	// Former groups should be preserved for "attacked this turn" queries
	if len(gameState.combat.formerGroups) == 0 {
		t.Error("former groups should be preserved after combat")
	}
}
