package rules

import "testing"

func TestTurnManagerSequence(t *testing.T) {
	tm := NewTurnManager("Alice")

	expected := []struct {
		phase Phase
		step  Step
	}{
		{PhaseBeginning, StepUntap},
		{PhaseBeginning, StepUpkeep},
		{PhaseBeginning, StepDraw},
		{PhasePrecombatMain, StepMain1},
		{PhaseCombat, StepBeginCombat},
		{PhaseCombat, StepDeclareAttackers},
		{PhaseCombat, StepDeclareBlockers},
		{PhaseCombat, StepCombatDamage},
		{PhaseCombat, StepEndCombat},
		{PhasePostcombatMain, StepMain2},
		{PhaseEnding, StepEnd},
		{PhaseEnding, StepCleanup},
	}

	for i, exp := range expected {
		if tm.CurrentPhase() != exp.phase {
			t.Fatalf("step %d: expected phase %s, got %s", i, exp.phase, tm.CurrentPhase())
		}
		if tm.CurrentStep() != exp.step {
			t.Fatalf("step %d: expected step %s, got %s", i, exp.step, tm.CurrentStep())
		}
		if i < len(expected)-1 {
			tm.AdvanceStep("")
		}
	}
}

func TestTurnManagerAdvanceWrapsTurn(t *testing.T) {
	tm := NewTurnManager("Alice")

	// Advance through all but the last step to remain on turn 1.
	for i := 0; i < 11; i++ {
		tm.AdvanceStep("")
		if tm.TurnNumber() != 1 {
			t.Fatalf("expected to remain on turn 1, got turn %d at step %d", tm.TurnNumber(), i)
		}
		if tm.ActivePlayer() != "Alice" {
			t.Fatalf("expected active player to remain Alice during turn, got %s", tm.ActivePlayer())
		}
	}

	phase, step := tm.AdvanceStep("Bob")
	if tm.TurnNumber() != 2 {
		t.Fatalf("expected turn number 2 after wrap, got %d", tm.TurnNumber())
	}
	if tm.ActivePlayer() != "Bob" {
		t.Fatalf("expected active player Bob after wrap, got %s", tm.ActivePlayer())
	}
	if tm.PriorityPlayer() != "Bob" {
		t.Fatalf("expected priority player Bob after wrap, got %s", tm.PriorityPlayer())
	}
	if phase != PhaseBeginning || step != StepUntap {
		t.Fatalf("expected new turn to start at BEGINNING/UNTAP, got %s/%s", phase, step)
	}
}
