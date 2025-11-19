package game

import (
	"testing"
)

// TestCombatView_Empty tests empty combat view before combat starts
func TestCombatView_Empty(t *testing.T) {
	h := NewCombatTestHarness(t, "test-combat-view-empty", []string{"Alice", "Bob"})

	// Get combat view before combat starts
	combatView, err := h.engine.GetCombatView(h.gameID)
	if err != nil {
		t.Fatalf("failed to get combat view: %v", err)
	}

	// Should have no attacking player and no groups
	if combatView.AttackingPlayerID != "" {
		t.Errorf("expected no attacking player, got %s", combatView.AttackingPlayerID)
	}
	if len(combatView.Groups) != 0 {
		t.Errorf("expected 0 combat groups, got %d", len(combatView.Groups))
	}
}

// TestCombatView_AfterSetup tests combat view after combat setup
func TestCombatView_AfterSetup(t *testing.T) {
	h := NewCombatTestHarness(t, "test-combat-view-setup", []string{"Alice", "Bob"})

	// Setup combat
	h.SetupCombat("Alice")

	// Get combat view
	combatView, err := h.engine.GetCombatView(h.gameID)
	if err != nil {
		t.Fatalf("failed to get combat view: %v", err)
	}

	// Should have Alice as attacking player but no groups yet
	if combatView.AttackingPlayerID != "Alice" {
		t.Errorf("expected Alice as attacking player, got %s", combatView.AttackingPlayerID)
	}
	if len(combatView.Groups) != 0 {
		t.Errorf("expected 0 combat groups before attackers declared, got %d", len(combatView.Groups))
	}
}

// TestCombatView_SingleAttacker tests combat view with one attacker
func TestCombatView_SingleAttacker(t *testing.T) {
	h := NewCombatTestHarness(t, "test-combat-view-attacker", []string{"Alice", "Bob"})

	// Create attacker
	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")

	// Setup combat and declare attacker
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")

	// Get combat view
	combatView, err := h.engine.GetCombatView(h.gameID)
	if err != nil {
		t.Fatalf("failed to get combat view: %v", err)
	}

	// Verify combat view
	if combatView.AttackingPlayerID != "Alice" {
		t.Errorf("expected Alice as attacking player, got %s", combatView.AttackingPlayerID)
	}
	if len(combatView.Groups) != 1 {
		t.Fatalf("expected 1 combat group, got %d", len(combatView.Groups))
	}

	group := combatView.Groups[0]
	if len(group.Attackers) != 1 {
		t.Errorf("expected 1 attacker in group, got %d", len(group.Attackers))
	}
	if len(group.Attackers) > 0 && group.Attackers[0] != attacker {
		t.Errorf("expected attacker %s, got %s", attacker, group.Attackers[0])
	}
	if group.DefenderID != "Bob" {
		t.Errorf("expected defender Bob, got %s", group.DefenderID)
	}
	if group.Blocked {
		t.Error("expected group to be unblocked")
	}
	if len(group.Blockers) != 0 {
		t.Errorf("expected 0 blockers, got %d", len(group.Blockers))
	}

	h.EndCombat()
}

// TestCombatView_MultipleAttackers tests combat view with multiple attackers
func TestCombatView_MultipleAttackers(t *testing.T) {
	h := NewCombatTestHarness(t, "test-combat-view-multi-attackers", []string{"Alice", "Bob"})

	// Create attackers
	attacker1 := h.CreateAttacker("attacker-1", "Bears", "Alice", "2", "2")
	attacker2 := h.CreateAttacker("attacker-2", "Giant", "Alice", "3", "3")

	// Setup combat and declare attackers
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker1, "Bob", "Alice")
	h.DeclareAttacker(attacker2, "Bob", "Alice")

	// Get combat view
	combatView, err := h.engine.GetCombatView(h.gameID)
	if err != nil {
		t.Fatalf("failed to get combat view: %v", err)
	}

	// Should have 2 combat groups (one per attacker)
	if len(combatView.Groups) != 2 {
		t.Fatalf("expected 2 combat groups, got %d", len(combatView.Groups))
	}

	// Verify each group has one attacker
	for i, group := range combatView.Groups {
		if len(group.Attackers) != 1 {
			t.Errorf("group %d: expected 1 attacker, got %d", i, len(group.Attackers))
		}
		if group.DefenderID != "Bob" {
			t.Errorf("group %d: expected defender Bob, got %s", i, group.DefenderID)
		}
		if group.Blocked {
			t.Errorf("group %d: expected unblocked", i)
		}
	}

	h.EndCombat()
}

// TestCombatView_WithBlockers tests combat view with blockers
func TestCombatView_WithBlockers(t *testing.T) {
	h := NewCombatTestHarness(t, "test-combat-view-blockers", []string{"Alice", "Bob"})

	// Create attacker and blocker
	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")
	blocker := h.CreateBlocker("blocker-1", "Hill Giant", "Bob", "3", "3")

	// Setup combat and declare blocker
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")
	h.DeclareBlocker(blocker, attacker, "Bob")
	h.AcceptBlockers()

	// Get combat view
	combatView, err := h.engine.GetCombatView(h.gameID)
	if err != nil {
		t.Fatalf("failed to get combat view: %v", err)
	}

	// Verify combat view shows blocker
	if len(combatView.Groups) != 1 {
		t.Fatalf("expected 1 combat group, got %d", len(combatView.Groups))
	}

	group := combatView.Groups[0]
	if !group.Blocked {
		t.Error("expected group to be blocked")
	}
	if len(group.Blockers) != 1 {
		t.Fatalf("expected 1 blocker, got %d", len(group.Blockers))
	}
	if group.Blockers[0] != blocker {
		t.Errorf("expected blocker %s, got %s", blocker, group.Blockers[0])
	}

	h.EndCombat()
}

// TestCombatView_MultipleBlockers tests combat view with multiple blockers on one attacker
func TestCombatView_MultipleBlockers(t *testing.T) {
	h := NewCombatTestHarness(t, "test-combat-view-multi-blockers", []string{"Alice", "Bob"})

	// Create attacker and blockers
	attacker := h.CreateAttacker("attacker-1", "Big Creature", "Alice", "5", "5")
	blocker1 := h.CreateBlocker("blocker-1", "Bear 1", "Bob", "2", "2")
	blocker2 := h.CreateBlocker("blocker-2", "Bear 2", "Bob", "2", "2")

	// Setup combat and declare blockers
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")
	h.DeclareBlocker(blocker1, attacker, "Bob")
	h.DeclareBlocker(blocker2, attacker, "Bob")
	h.AcceptBlockers()

	// Get combat view
	combatView, err := h.engine.GetCombatView(h.gameID)
	if err != nil {
		t.Fatalf("failed to get combat view: %v", err)
	}

	// Verify combat view shows both blockers
	if len(combatView.Groups) != 1 {
		t.Fatalf("expected 1 combat group, got %d", len(combatView.Groups))
	}

	group := combatView.Groups[0]
	if !group.Blocked {
		t.Error("expected group to be blocked")
	}
	if len(group.Blockers) != 2 {
		t.Fatalf("expected 2 blockers, got %d", len(group.Blockers))
	}

	// Verify both blockers are present (order may vary)
	blockerMap := make(map[string]bool)
	for _, b := range group.Blockers {
		blockerMap[b] = true
	}
	if !blockerMap[blocker1] {
		t.Errorf("blocker1 %s not found in group", blocker1)
	}
	if !blockerMap[blocker2] {
		t.Errorf("blocker2 %s not found in group", blocker2)
	}

	h.EndCombat()
}

// TestCombatView_InGameView tests that combat view is included in game view
func TestCombatView_InGameView(t *testing.T) {
	h := NewCombatTestHarness(t, "test-combat-view-game-view", []string{"Alice", "Bob"})

	// Create attacker
	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")

	// Setup combat and declare attacker
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker, "Bob", "Alice")

	// Get game view
	viewRaw, err := h.engine.GetGameView(h.gameID, "Alice")
	if err != nil {
		t.Fatalf("failed to get game view: %v", err)
	}

	view, ok := viewRaw.(*EngineGameView)
	if !ok {
		t.Fatal("game view is not of type *EngineGameView")
	}

	// Verify combat view is populated in game view
	if view.Combat.AttackingPlayerID != "Alice" {
		t.Errorf("expected Alice as attacking player in game view, got %s", view.Combat.AttackingPlayerID)
	}
	if len(view.Combat.Groups) != 1 {
		t.Fatalf("expected 1 combat group in game view, got %d", len(view.Combat.Groups))
	}

	group := view.Combat.Groups[0]
	if len(group.Attackers) != 1 {
		t.Errorf("expected 1 attacker in game view group, got %d", len(group.Attackers))
	}
	if len(group.Attackers) > 0 && group.Attackers[0] != attacker {
		t.Errorf("expected attacker %s in game view, got %s", attacker, group.Attackers[0])
	}

	h.EndCombat()
}

// TestCombatView_AfterCombat tests combat view after combat ends
func TestCombatView_AfterCombat(t *testing.T) {
	h := NewCombatTestHarness(t, "test-combat-view-after", []string{"Alice", "Bob"})

	// Create attacker
	attacker := h.CreateAttacker("attacker-1", "Grizzly Bears", "Alice", "2", "2")

	// Run full combat
	h.RunFullCombat("Alice", map[string]string{
		attacker: "Bob",
	}, nil)

	// Get combat view after combat
	combatView, err := h.engine.GetCombatView(h.gameID)
	if err != nil {
		t.Fatalf("failed to get combat view: %v", err)
	}

	// Combat should be cleared (no current groups)
	if len(combatView.Groups) != 0 {
		t.Errorf("expected 0 combat groups after combat ends, got %d", len(combatView.Groups))
	}
	// Attacking player may still be set, but groups should be empty
}

// TestCombatView_ComplexScenario tests combat view in complex multi-group scenario
func TestCombatView_ComplexScenario(t *testing.T) {
	h := NewCombatTestHarness(t, "test-combat-view-complex", []string{"Alice", "Bob"})

	// Create multiple attackers and blockers
	attacker1 := h.CreateAttacker("attacker-1", "Bears", "Alice", "2", "2")
	attacker2 := h.CreateAttacker("attacker-2", "Giant", "Alice", "3", "3")
	attacker3 := h.CreateAttacker("attacker-3", "Dragon", "Alice", "4", "4")
	blocker1 := h.CreateBlocker("blocker-1", "Wall", "Bob", "0", "4")
	blocker2 := h.CreateBlocker("blocker-2", "Knight", "Bob", "2", "2")

	// Setup complex combat scenario
	h.SetupCombat("Alice")
	h.DeclareAttacker(attacker1, "Bob", "Alice")
	h.DeclareAttacker(attacker2, "Bob", "Alice")
	h.DeclareAttacker(attacker3, "Bob", "Alice")
	h.DeclareBlocker(blocker1, attacker1, "Bob") // Wall blocks Bears
	h.DeclareBlocker(blocker2, attacker2, "Bob") // Knight blocks Giant
	// Dragon unblocked
	h.AcceptBlockers()

	// Get combat view
	combatView, err := h.engine.GetCombatView(h.gameID)
	if err != nil {
		t.Fatalf("failed to get combat view: %v", err)
	}

	// Should have 3 combat groups
	if len(combatView.Groups) != 3 {
		t.Fatalf("expected 3 combat groups, got %d", len(combatView.Groups))
	}

	// Count blocked and unblocked groups
	blockedCount := 0
	unblockedCount := 0
	for _, group := range combatView.Groups {
		if group.Blocked {
			blockedCount++
		} else {
			unblockedCount++
		}
	}

	if blockedCount != 2 {
		t.Errorf("expected 2 blocked groups, got %d", blockedCount)
	}
	if unblockedCount != 1 {
		t.Errorf("expected 1 unblocked group, got %d", unblockedCount)
	}

	h.EndCombat()
}
