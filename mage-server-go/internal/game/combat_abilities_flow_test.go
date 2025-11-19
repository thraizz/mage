package game

import (
	"testing"
)

// TestCombatFlow_FirstStrikeDamage tests first strike combat damage mechanics
// Per GO_PORT_TASKS.md: "Test first strike damage (kill before normal damage)"
func TestCombatFlow_FirstStrikeDamage(t *testing.T) {
	h := NewCombatTestHarness(t, "test-first-strike", []string{"Alice", "Bob"})

	// Create a 2/2 first strike attacker for Alice
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "First Strike Knight",
		Power:      "2",
		Toughness:  "2",
		Controller: "Alice",
		Abilities:  []string{abilityFirstStrike},
	})

	// Create a 2/2 blocker for Bob (no first strike)
	blocker := h.CreateBlocker("blocker-1", "Grizzly Bears", "Bob", "2", "2")

	// Get initial life totals
	initialBobLife := h.GetPlayerLife("Bob")

	// Run combat: attacker with first strike attacks, blocker blocks
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, map[string]string{
		blocker: attacker,
	})

	// Verify results
	h.AssertPlayerLife("Bob", initialBobLife) // Bob takes no damage (blocked)
	h.AssertCreatureDead(blocker)              // Blocker dies from first strike damage
	h.AssertCreatureAlive(attacker)            // Attacker survives (takes no damage - blocker died before dealing damage)
}

// TestCombatFlow_DoubleStrike tests double strike damage in both steps
// Per GO_PORT_TASKS.md: "Test double strike damage (damage in both steps)"
func TestCombatFlow_DoubleStrike(t *testing.T) {
	h := NewCombatTestHarness(t, "test-double-strike", []string{"Alice", "Bob"})

	// Create a 2/2 double strike attacker for Alice
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "Double Strike Knight",
		Power:      "2",
		Toughness:  "2",
		Controller: "Alice",
		Abilities:  []string{abilityDoubleStrike},
	})

	// Get initial life totals (unblocked, should deal damage twice)
	initialBobLife := h.GetPlayerLife("Bob")

	// Run combat: attacker with double strike attacks unblocked
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, nil)

	// Verify results - double strike deals damage twice (2 + 2 = 4 total)
	h.AssertPlayerLife("Bob", initialBobLife-4)
	h.AssertCreatureAlive(attacker)
}

// TestCombatFlow_Vigilance tests that vigilance prevents tapping
// Per GO_PORT_TASKS.md: "Test vigilance (no tap on attack)"
func TestCombatFlow_Vigilance(t *testing.T) {
	h := NewCombatTestHarness(t, "test-vigilance", []string{"Alice", "Bob"})

	// Create a creature with vigilance for Alice
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "Vigilant Knight",
		Power:      "3",
		Toughness:  "3",
		Controller: "Alice",
		Abilities:  []string{abilityVigilance},
	})

	// Run combat
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, nil)

	// Verify attacker is NOT tapped (vigilance)
	h.AssertCreatureTapped(attacker, false)
}

// TestCombatFlow_FlyingReach tests flying and reach interaction
// Per GO_PORT_TASKS.md: "Test flying/reach restrictions"
func TestCombatFlow_FlyingReach(t *testing.T) {
	h := NewCombatTestHarness(t, "test-flying-reach", []string{"Alice", "Bob"})

	// Create a flying attacker for Alice
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "Wind Drake",
		Power:      "2",
		Toughness:  "2",
		Controller: "Alice",
		Abilities:  []string{abilityFlying},
	})

	// Create a non-flying blocker for Bob (can't block flying)
	blocker1 := h.CreateBlocker("blocker-1", "Grizzly Bears", "Bob", "2", "2")

	// Create a reach blocker for Bob (can block flying)
	blocker2 := h.CreateCreature(CreatureSpec{
		ID:         "blocker-2",
		Name:       "Giant Spider",
		Power:      "2",
		Toughness:  "4",
		Controller: "Bob",
		Abilities:  []string{abilityReach},
	})

	// Setup combat
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")

	// Try to block with non-flying creature - should fail
	err := h.engine.DeclareBlocker(h.gameID, blocker1, attacker, "Bob")
	if err == nil {
		t.Error("non-flying creature should not be able to block flying creature")
	}

	// Block with reach creature - should succeed
	h.DeclareBlocker(blocker2, attacker, "Bob")
	h.AcceptBlockers()

	// Complete combat
	h.AssignDamage(false)
	h.ApplyDamage()
	h.EndCombat()

	// Verify the reach blocker blocked successfully
	// Attacker (2/2) takes 2 damage from blocker (2/4) and dies
	// Blocker (2/4) takes 2 damage from attacker and survives
	h.AssertCreatureDead(attacker)
	h.AssertCreatureAlive(blocker2)
}

// TestCombatFlow_Deathtouch tests deathtouch makes any damage lethal
func TestCombatFlow_Deathtouch(t *testing.T) {
	h := NewCombatTestHarness(t, "test-deathtouch", []string{"Alice", "Bob"})

	// Create a 1/1 deathtouch attacker for Alice
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "Typhoid Rats",
		Power:      "1",
		Toughness:  "1",
		Controller: "Alice",
		Abilities:  []string{abilityDeathtouch},
	})

	// Create a 5/5 blocker for Bob
	blocker := h.CreateBlocker("blocker-1", "Craw Wurm", "Bob", "6", "4")

	// Run combat
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, map[string]string{
		blocker: attacker,
	})

	// Verify deathtouch killed the large creature with just 1 damage
	h.AssertCreatureDead(blocker)   // Dies from deathtouch
	h.AssertCreatureDead(attacker)  // Dies from blocker's damage
}

// TestCombatFlow_Lifelink tests lifelink grants life equal to damage dealt
func TestCombatFlow_Lifelink(t *testing.T) {
	h := NewCombatTestHarness(t, "test-lifelink", []string{"Alice", "Bob"})

	// Create a 3/3 lifelink attacker for Alice
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "Vampire Nighthawk",
		Power:      "3",
		Toughness:  "3",
		Controller: "Alice",
		Abilities:  []string{abilityLifelink},
	})

	// Get initial life totals
	initialAliceLife := h.GetPlayerLife("Alice")

	// Run combat (unblocked)
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, nil)

	// Verify Alice gained 3 life from lifelink
	h.AssertPlayerLife("Alice", initialAliceLife+3)
}

// TestCombatFlow_Defender tests defender ability prevents attacking
func TestCombatFlow_Defender(t *testing.T) {
	h := NewCombatTestHarness(t, "test-defender", []string{"Alice", "Bob"})

	// Create a creature with defender for Alice
	defender := h.CreateCreature(CreatureSpec{
		ID:         "defender-1",
		Name:       "Wall of Stone",
		Power:      "0",
		Toughness:  "8",
		Controller: "Alice",
		Abilities:  []string{abilityDefender},
	})

	// Setup combat
	h.SetupCombat("Alice")

	// Try to attack with defender - should fail
	err := h.engine.DeclareAttacker(h.gameID, defender, "Bob", "Alice")
	if err == nil {
		t.Error("creature with defender should not be able to attack")
	}

	// End combat
	h.EndCombat()
}

// TestCombatFlow_Menace tests menace requires two blockers
func TestCombatFlow_Menace(t *testing.T) {
	h := NewCombatTestHarness(t, "test-menace", []string{"Alice", "Bob"})

	// Create a menace attacker for Alice
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "Menace Creature",
		Power:      "3",
		Toughness:  "3",
		Controller: "Alice",
		Abilities:  []string{abilityMenace},
	})

	// Create one blocker for Bob
	blocker1 := h.CreateBlocker("blocker-1", "Grizzly Bears", "Bob", "2", "2")

	// Get Bob's initial life
	initialBobLife := h.GetPlayerLife("Bob")

	// Setup combat
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")

	// Declare one blocker - this is illegal for menace
	h.DeclareBlocker(blocker1, attacker, "Bob")

	// AcceptBlockers should succeed but automatically remove the illegal blocker
	// Per Java implementation: menace violations result in blockers being removed, not an error
	h.AcceptBlockers()

	// Verify blocker was removed (no longer blocking)
	if h.IsCreatureBlocking(blocker1) {
		t.Error("blocker should have been removed due to menace violation")
	}

	// Complete combat
	h.AssignDamage(false)
	h.ApplyDamage()
	h.EndCombat()

	// Verify attacker dealt damage to player (unblocked due to menace)
	h.AssertPlayerLife("Bob", initialBobLife-3)
	h.AssertCreatureAlive(attacker)
	h.AssertCreatureAlive(blocker1) // Blocker survives (wasn't in combat)
}

// TestCombatFlow_MultipleAbilities tests creature with multiple combat abilities
func TestCombatFlow_MultipleAbilities(t *testing.T) {
	h := NewCombatTestHarness(t, "test-multiple-abilities", []string{"Alice", "Bob"})

	// Create a creature with flying, vigilance, and lifelink
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "Serra Angel",
		Power:      "4",
		Toughness:  "4",
		Controller: "Alice",
		Abilities:  []string{abilityFlying, abilityVigilance, abilityLifelink},
	})

	// Get initial life
	initialAliceLife := h.GetPlayerLife("Alice")

	// Run combat (unblocked due to flying)
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, nil)

	// Verify all abilities worked
	h.AssertCreatureTapped(attacker, false)         // Vigilance - not tapped
	h.AssertPlayerLife("Alice", initialAliceLife+4) // Lifelink - gained 4 life
	// Flying means it could attack (we don't test blocking restrictions here)
}

// TestCombatFlow_TrampleUnblocked tests trample with no blockers (all damage to player)
func TestCombatFlow_TrampleUnblocked(t *testing.T) {
	h := NewCombatTestHarness(t, "test-trample-unblocked", []string{"Alice", "Bob"})

	// Create a trample attacker
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "Craw Wurm",
		Power:      "6",
		Toughness:  "4",
		Controller: "Alice",
		Abilities:  []string{abilityTrample},
	})

	initialBobLife := h.GetPlayerLife("Bob")

	// Run combat (unblocked)
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, nil)

	// All damage goes to player
	h.AssertPlayerLife("Bob", initialBobLife-6)
}

// TestCombatFlow_TrampleBlocked tests trample with blocker (excess damage to player)
func TestCombatFlow_TrampleBlocked(t *testing.T) {
	h := NewCombatTestHarness(t, "test-trample-blocked", []string{"Alice", "Bob"})

	// Create a 6/4 trample attacker
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "Craw Wurm",
		Power:      "6",
		Toughness:  "4",
		Controller: "Alice",
		Abilities:  []string{abilityTrample},
	})

	// Create a 2/2 blocker
	blocker := h.CreateBlocker("blocker-1", "Grizzly Bears", "Bob", "2", "2")

	initialBobLife := h.GetPlayerLife("Bob")

	// Run combat
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, map[string]string{
		blocker: attacker,
	})

	// 2 damage to kill blocker, 4 damage tramples over to Bob
	h.AssertCreatureDead(blocker)
	h.AssertPlayerLife("Bob", initialBobLife-4) // Trample overflow
	h.AssertCreatureAlive(attacker)              // Attacker survives
}
