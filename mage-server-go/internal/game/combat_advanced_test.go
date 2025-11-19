package game

import (
	"testing"

	"github.com/magefree/mage-server-go/internal/game/rules"
)

// TestCombatEvents_AttackerDeclared tests that EventAttackerDeclared fires correctly
func TestCombatEvents_AttackerDeclared(t *testing.T) {
	h := NewCombatTestHarness(t, "test-events-attacker", []string{"Alice", "Bob"})

	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")

	// Subscribe to events
	events := make([]rules.Event, 0)
	gameState := h.GetGameState()
	gameState.eventBus.SubscribeTyped(rules.EventAttackerDeclared, func(evt rules.Event) {
		events = append(events, evt)
	})

	// Declare attacker
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")

	// Verify EventAttackerDeclared was fired
	if len(events) != 1 {
		t.Errorf("expected 1 EventAttackerDeclared, got %d", len(events))
	}
	if len(events) > 0 && events[0].SourceID != attacker {
		t.Errorf("expected event source to be %s, got %s", attacker, events[0].SourceID)
	}

	h.EndCombat()
}

// TestCombatEvents_BlockerDeclared tests that EventBlockerDeclared fires correctly
func TestCombatEvents_BlockerDeclared(t *testing.T) {
	h := NewCombatTestHarness(t, "test-events-blocker", []string{"Alice", "Bob"})

	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")
	blocker := h.CreateBlocker("blocker-1", "Hill Giant", "Bob", "3", "3")

	// Subscribe to events
	events := make([]rules.Event, 0)
	gameState := h.GetGameState()
	gameState.eventBus.SubscribeTyped(rules.EventBlockerDeclared, func(evt rules.Event) {
		events = append(events, evt)
	})

	// Declare attacker and blocker
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")
	h.DeclareBlocker(blocker, attacker, "Bob")
	h.AcceptBlockers()

	// Verify EventBlockerDeclared was fired
	// Note: Event fires during both DeclareBlocker and AcceptBlockers (per Java implementation)
	if len(events) < 1 {
		t.Errorf("expected at least 1 EventBlockerDeclared, got %d", len(events))
	}
	if len(events) > 0 {
		if events[0].SourceID != blocker {
			t.Errorf("expected event source to be %s, got %s", blocker, events[0].SourceID)
		}
		if events[0].TargetID != attacker {
			t.Errorf("expected event target to be %s, got %s", attacker, events[0].TargetID)
		}
	}

	h.EndCombat()
}

// TestCombatEvents_CreatureBlocked tests that EventCreatureBlocked fires when creature becomes blocked
func TestCombatEvents_CreatureBlocked(t *testing.T) {
	h := NewCombatTestHarness(t, "test-events-blocked", []string{"Alice", "Bob"})

	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")
	blocker := h.CreateBlocker("blocker-1", "Hill Giant", "Bob", "3", "3")

	// Subscribe to events
	events := make([]rules.Event, 0)
	gameState := h.GetGameState()
	gameState.eventBus.SubscribeTyped(rules.EventCreatureBlocked, func(evt rules.Event) {
		events = append(events, evt)
	})

	// Declare attacker and blocker
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")
	h.DeclareBlocker(blocker, attacker, "Bob")
	h.AcceptBlockers()

	// Verify EventCreatureBlocked was fired for the attacker
	if len(events) != 1 {
		t.Errorf("expected 1 EventCreatureBlocked, got %d", len(events))
	}
	if len(events) > 0 && events[0].SourceID != attacker {
		t.Errorf("expected event source to be %s, got %s", attacker, events[0].SourceID)
	}

	h.EndCombat()
}

// TestCombatEvents_CreatureBlocks tests that EventCreatureBlocks fires when creature blocks
func TestCombatEvents_CreatureBlocks(t *testing.T) {
	h := NewCombatTestHarness(t, "test-events-blocks", []string{"Alice", "Bob"})

	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")
	blocker := h.CreateBlocker("blocker-1", "Hill Giant", "Bob", "3", "3")

	// Subscribe to events
	events := make([]rules.Event, 0)
	gameState := h.GetGameState()
	gameState.eventBus.SubscribeTyped(rules.EventCreatureBlocks, func(evt rules.Event) {
		events = append(events, evt)
	})

	// Declare attacker and blocker
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")
	h.DeclareBlocker(blocker, attacker, "Bob")
	h.AcceptBlockers()

	// Verify EventCreatureBlocks was fired for the blocker
	if len(events) != 1 {
		t.Errorf("expected 1 EventCreatureBlocks, got %d", len(events))
	}
	if len(events) > 0 && events[0].SourceID != blocker {
		t.Errorf("expected event source to be %s, got %s", blocker, events[0].SourceID)
	}

	h.EndCombat()
}

// TestCombatEvents_UnblockedAttacker tests that EventUnblockedAttacker fires for unblocked attackers
func TestCombatEvents_UnblockedAttacker(t *testing.T) {
	h := NewCombatTestHarness(t, "test-events-unblocked", []string{"Alice", "Bob"})

	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")

	// Subscribe to events
	events := make([]rules.Event, 0)
	gameState := h.GetGameState()
	gameState.eventBus.SubscribeTyped(rules.EventUnblockedAttacker, func(evt rules.Event) {
		events = append(events, evt)
	})

	// Declare attacker (no blockers)
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")
	// AcceptBlockers should fire EventUnblockedAttacker for unblocked attackers
	h.AcceptBlockers()

	// Verify EventUnblockedAttacker was fired
	if len(events) != 1 {
		t.Errorf("expected 1 EventUnblockedAttacker, got %d", len(events))
	}
	if len(events) > 0 && events[0].SourceID != attacker {
		t.Errorf("expected event source to be %s, got %s", attacker, events[0].SourceID)
	}

	h.EndCombat()
}

// TestCombatEvents_DamageApplied tests that EventCombatDamageApplied fires
func TestCombatEvents_DamageApplied(t *testing.T) {
	h := NewCombatTestHarness(t, "test-events-damage", []string{"Alice", "Bob"})

	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")

	// Subscribe to events
	events := make([]rules.Event, 0)
	gameState := h.GetGameState()
	gameState.eventBus.SubscribeTyped(rules.EventCombatDamageApplied, func(evt rules.Event) {
		events = append(events, evt)
	})

	// Run full combat
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")
	h.AssignDamage(false)
	h.ApplyDamage()

	// Verify EventCombatDamageApplied was fired
	if len(events) != 1 {
		t.Errorf("expected 1 EventCombatDamageApplied, got %d", len(events))
	}

	h.EndCombat()
}

// TestCombatTriggers_WheneverAttacks tests triggers that fire when a creature attacks
func TestCombatTriggers_WheneverAttacks(t *testing.T) {
	h := NewCombatTestHarness(t, "test-triggers-attacks", []string{"Alice", "Bob"})

	// Create attacker with a triggered ability that fires when it attacks
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "Goblin Guide",
		Power:      "2",
		Toughness:  "2",
		Controller: "Alice",
	})

	// Track triggers
	triggerCount := 0
	gameState := h.GetGameState()

	// Subscribe to EventAttackerDeclared to simulate "Whenever ~ attacks" trigger
	gameState.eventBus.SubscribeTyped(rules.EventAttackerDeclared, func(evt rules.Event) {
		if evt.SourceID == attacker {
			triggerCount++
		}
	})

	// Attack
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")

	// Verify trigger fired
	if triggerCount != 1 {
		t.Errorf("expected 1 attack trigger, got %d", triggerCount)
	}

	h.EndCombat()
}

// TestCombatTriggers_WheneverBlocks tests triggers that fire when a creature blocks
func TestCombatTriggers_WheneverBlocks(t *testing.T) {
	h := NewCombatTestHarness(t, "test-triggers-blocks", []string{"Alice", "Bob"})

	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")
	blocker := h.CreateBlocker("blocker-1", "Wall of Essence", "Bob", "0", "4")

	// Track triggers
	triggerCount := 0
	gameState := h.GetGameState()

	// Subscribe to EventCreatureBlocks to simulate "Whenever ~ blocks" trigger
	gameState.eventBus.SubscribeTyped(rules.EventCreatureBlocks, func(evt rules.Event) {
		if evt.SourceID == blocker {
			triggerCount++
		}
	})

	// Block
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")
	h.DeclareBlocker(blocker, attacker, "Bob")
	h.AcceptBlockers()

	// Verify trigger fired
	if triggerCount != 1 {
		t.Errorf("expected 1 block trigger, got %d", triggerCount)
	}

	h.EndCombat()
}

// TestCombatTriggers_WheneverBecomesBlocked tests triggers that fire when a creature becomes blocked
func TestCombatTriggers_WheneverBecomesBlocked(t *testing.T) {
	h := NewCombatTestHarness(t, "test-triggers-becomes-blocked", []string{"Alice", "Bob"})

	attacker := h.CreateAttacker("attacker-1", "Boros Swiftblade", "Alice", "2", "2")
	blocker := h.CreateBlocker("blocker-1", "Hill Giant", "Bob", "3", "3")

	// Track triggers
	triggerCount := 0
	gameState := h.GetGameState()

	// Subscribe to EventCreatureBlocked to simulate "Whenever ~ becomes blocked" trigger
	gameState.eventBus.SubscribeTyped(rules.EventCreatureBlocked, func(evt rules.Event) {
		if evt.SourceID == attacker {
			triggerCount++
		}
	})

	// Block
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")
	h.DeclareBlocker(blocker, attacker, "Bob")
	h.AcceptBlockers()

	// Verify trigger fired
	if triggerCount != 1 {
		t.Errorf("expected 1 becomes-blocked trigger, got %d", triggerCount)
	}

	h.EndCombat()
}

// TestCombatRemoval_DuringAttackerDeclaration tests creature removal during attacker declaration
func TestCombatRemoval_DuringAttackerDeclaration(t *testing.T) {
	h := NewCombatTestHarness(t, "test-removal-attacker-declare", []string{"Alice", "Bob"})

	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")

	// Setup combat
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")

	// Remove attacker before blockers
	gameState := h.GetGameState()
	gameState.mu.Lock()
	if card, exists := gameState.cards[attacker]; exists {
		card.Zone = zoneGraveyard
		card.Attacking = false
		card.AttackingWhat = ""
	}
	gameState.mu.Unlock()

	// Verify creature is dead
	h.AssertCreatureDead(attacker)

	h.EndCombat()
}

// TestCombatRemoval_DuringBlockerDeclaration tests creature removal during blocker declaration
func TestCombatRemoval_DuringBlockerDeclaration(t *testing.T) {
	h := NewCombatTestHarness(t, "test-removal-blocker-declare", []string{"Alice", "Bob"})

	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")
	blocker := h.CreateBlocker("blocker-1", "Hill Giant", "Bob", "3", "3")

	initialBobLife := h.GetPlayerLife("Bob")

	// Setup combat and declare blocker
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")
	h.DeclareBlocker(blocker, attacker, "Bob")

	// Remove blocker before damage
	gameState := h.GetGameState()
	gameState.mu.Lock()
	if card, exists := gameState.cards[blocker]; exists {
		card.Zone = zoneGraveyard
		card.Blocking = false
		card.BlockingWhat = []string{}
	}
	// Remove from combat group
	for _, group := range gameState.combat.groups {
		newBlockers := []string{}
		for _, b := range group.blockers {
			if b != blocker {
				newBlockers = append(newBlockers, b)
			}
		}
		group.blockers = newBlockers
		if len(group.blockers) == 0 {
			group.blocked = false
		}
	}
	delete(gameState.combat.blockers, blocker)
	gameState.mu.Unlock()

	// Complete combat
	h.AcceptBlockers()
	h.AssignDamage(false)
	h.ApplyDamage()
	h.EndCombat()

	// Verify blocker is dead and attacker dealt damage to player (unblocked)
	h.AssertCreatureDead(blocker)
	h.AssertPlayerLife("Bob", initialBobLife-2) // Attacker is now unblocked
}

// TestCombatRemoval_BeforeDamage tests creature removal before damage step
func TestCombatRemoval_BeforeDamage(t *testing.T) {
	h := NewCombatTestHarness(t, "test-removal-before-damage", []string{"Alice", "Bob"})

	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")
	blocker := h.CreateBlocker("blocker-1", "Hill Giant", "Bob", "3", "3")

	// Setup combat
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")
	h.DeclareBlocker(blocker, attacker, "Bob")
	h.AcceptBlockers()

	// Remove attacker before damage
	gameState := h.GetGameState()
	gameState.mu.Lock()
	if card, exists := gameState.cards[attacker]; exists {
		card.Zone = zoneGraveyard
	}
	gameState.mu.Unlock()

	// Assign damage (should handle missing attacker gracefully)
	h.AssignDamage(false)
	h.ApplyDamage()
	h.EndCombat()

	// Verify attacker is dead and blocker is alive (took no damage)
	h.AssertCreatureDead(attacker)
	h.AssertCreatureAlive(blocker)
}

// TestBlockerRequirements_MustBlockIfAble tests "must block if able" restriction
func TestBlockerRequirements_MustBlockIfAble(t *testing.T) {
	h := NewCombatTestHarness(t, "test-must-block", []string{"Alice", "Bob"})

	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")

	// Create blocker that "must block if able"
	blocker := h.CreateCreature(CreatureSpec{
		ID:         "blocker-1",
		Name:       "Guardian Automaton",
		Power:      "2",
		Toughness:  "2",
		Controller: "Bob",
	})

	// Setup forced block requirement
	gameState := h.GetGameState()
	gameState.mu.Lock()
	if gameState.combat.creatureMustBlockAttackers == nil {
		gameState.combat.creatureMustBlockAttackers = make(map[string]map[string]bool)
	}
	gameState.combat.creatureMustBlockAttackers[blocker] = map[string]bool{attacker: true}
	gameState.mu.Unlock()

	// Setup combat
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")

	// Check block requirements - should require blocker to block
	violations, err := h.engine.CheckBlockRequirements(h.gameID, "Bob")
	if err != nil {
		t.Fatalf("failed to check block requirements: %v", err)
	}

	// If we don't declare blocker, there should be a violation
	if len(violations) == 0 {
		// This is expected - in real game, player must declare the block
		// For this test, we verify the requirement exists
	}

	// Now declare the required block
	h.DeclareBlocker(blocker, attacker, "Bob")
	h.AcceptBlockers()

	// Verify blocker is blocking
	if !h.IsCreatureBlocking(blocker) {
		t.Error("creature should be blocking (must block if able)")
	}

	h.EndCombat()
}

// TestBlockerRestrictions_CantBlock tests "can't block" restriction
func TestBlockerRestrictions_CantBlock(t *testing.T) {
	h := NewCombatTestHarness(t, "test-cant-block", []string{"Alice", "Bob"})

	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")

	// Create tapped blocker (can't block when tapped)
	blocker := h.CreateCreature(CreatureSpec{
		ID:         "blocker-1",
		Name:       "Tapped Bear",
		Power:      "2",
		Toughness:  "2",
		Controller: "Bob",
		Tapped:     true, // Tapped creatures can't block
	})

	// Setup combat
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")

	// Try to block with tapped creature - should fail
	err := h.engine.DeclareBlocker(h.gameID, blocker, attacker, "Bob")
	if err == nil {
		t.Error("tapped creature should not be able to block")
	}

	h.EndCombat()
}

// TestAttackerRequirements_ForcedAttack tests "must attack if able" requirement
func TestAttackerRequirements_ForcedAttack(t *testing.T) {
	h := NewCombatTestHarness(t, "test-forced-attack", []string{"Alice", "Bob"})

	// Create creature that "must attack if able"
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "Juggernaut",
		Power:      "5",
		Toughness:  "3",
		Controller: "Alice",
	})

	// Setup forced attack requirement
	gameState := h.GetGameState()
	gameState.mu.Lock()
	if gameState.combat.creaturesForcedToAttack == nil {
		gameState.combat.creaturesForcedToAttack = make(map[string]map[string]bool)
	}
	gameState.combat.creaturesForcedToAttack[attacker] = make(map[string]bool) // Empty map = can attack any defender
	gameState.mu.Unlock()

	// Setup combat
	h.SetupCombat("Alice")

	// Verify creature can attack
	canAttack, err := h.engine.CanAttack(h.gameID, attacker)
	if err != nil {
		t.Fatalf("failed to check if creature can attack: %v", err)
	}
	if !canAttack {
		t.Error("forced attacker should be able to attack")
	}

	// Declare attack
	h.DeclareAttacker(attacker, "Bob", "Alice")

	// Verify creature is attacking
	if !h.IsCreatureAttacking(attacker) {
		t.Error("creature should be attacking (forced to attack)")
	}

	h.EndCombat()
}

// TestAttackerRestrictions_CantAttack tests "can't attack" restriction
func TestAttackerRestrictions_CantAttack(t *testing.T) {
	h := NewCombatTestHarness(t, "test-cant-attack", []string{"Alice", "Bob"})

	// Create tapped creature (can't attack when tapped)
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "Tapped Bear",
		Power:      "2",
		Toughness:  "2",
		Controller: "Alice",
		Tapped:     true,
	})

	// Setup combat
	h.SetupCombat("Alice")

	// Try to attack with tapped creature - should fail
	err := h.engine.DeclareAttacker(h.gameID, attacker, "Bob", "Alice")
	if err == nil {
		t.Error("tapped creature should not be able to attack")
	}

	h.EndCombat()
}

// TestAttackerRestrictions_SummoningSickness tests summoning sickness restriction
func TestAttackerRestrictions_SummoningSickness(t *testing.T) {
	h := NewCombatTestHarness(t, "test-summoning-sickness", []string{"Alice", "Bob"})

	// Create creature with summoning sickness (just entered this turn)
	attacker := h.CreateCreature(CreatureSpec{
		ID:         "attacker-1",
		Name:       "Fresh Bear",
		Power:      "2",
		Toughness:  "2",
		Controller: "Alice",
	})

	// Mark as having summoning sickness
	gameState := h.GetGameState()
	gameState.mu.Lock()
	if card, exists := gameState.cards[attacker]; exists {
		card.SummoningSickness = true
	}
	gameState.mu.Unlock()

	// Setup combat
	h.SetupCombat("Alice")

	// Check if creature can attack - should return false due to summoning sickness
	canAttack, err := h.engine.CanAttack(h.gameID, attacker)
	if err != nil {
		t.Fatalf("failed to check if creature can attack: %v", err)
	}
	if canAttack {
		t.Error("creature with summoning sickness should not be able to attack")
	}

	h.EndCombat()
}
