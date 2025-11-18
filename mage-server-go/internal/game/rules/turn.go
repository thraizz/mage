package rules

import (
	"fmt"
	"strings"
)

// Phase represents the broad phases of a Magic: The Gathering turn.
type Phase int

const (
	PhaseBeginning Phase = iota
	PhasePrecombatMain
	PhaseCombat
	PhasePostcombatMain
	PhaseEnding
)

var phaseNames = map[Phase]string{
	PhaseBeginning:      "BEGINNING",
	PhasePrecombatMain:  "PRECOMBAT_MAIN",
	PhaseCombat:         "COMBAT",
	PhasePostcombatMain: "POSTCOMBAT_MAIN",
	PhaseEnding:         "ENDING",
}

func (p Phase) String() string {
	if name, ok := phaseNames[p]; ok {
		return name
	}
	return fmt.Sprintf("PHASE_%d", int(p))
}

// Step represents the individual steps that comprise a turn.
type Step int

const (
	StepUntap Step = iota
	StepUpkeep
	StepDraw
	StepMain1
	StepBeginCombat
	StepDeclareAttackers
	StepDeclareBlockers
	StepFirstStrikeDamage
	StepCombatDamage
	StepEndCombat
	StepMain2
	StepEnd
	StepCleanup
)

var stepNames = map[Step]string{
	StepUntap:             "UNTAP",
	StepUpkeep:            "UPKEEP",
	StepDraw:              "DRAW",
	StepMain1:             "MAIN1",
	StepBeginCombat:       "BEGIN_COMBAT",
	StepDeclareAttackers:  "DECLARE_ATTACKERS",
	StepDeclareBlockers:   "DECLARE_BLOCKERS",
	StepFirstStrikeDamage: "FIRST_STRIKE_DAMAGE",
	StepCombatDamage:      "COMBAT_DAMAGE",
	StepEndCombat:         "END_COMBAT",
	StepMain2:             "MAIN2",
	StepEnd:               "END",
	StepCleanup:           "CLEANUP",
}

func (s Step) String() string {
	if name, ok := stepNames[s]; ok {
		return name
	}
	return fmt.Sprintf("STEP_%d", int(s))
}

type turnEntry struct {
	phase Phase
	step  Step
}

// baseTurnSequence is the default turn structure without first strike damage step
var baseTurnSequence = []turnEntry{
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

// buildTurnSequence creates the turn sequence, optionally including StepFirstStrikeDamage
// if hasFirstStrike is true
func buildTurnSequence(hasFirstStrike bool) []turnEntry {
	sequence := make([]turnEntry, len(baseTurnSequence))
	copy(sequence, baseTurnSequence)

	if !hasFirstStrike {
		return sequence
	}

	// Insert StepFirstStrikeDamage before StepCombatDamage
	// Find the index of StepCombatDamage
	damageIdx := -1
	for i, entry := range sequence {
		if entry.step == StepCombatDamage {
			damageIdx = i
			break
		}
	}

	if damageIdx == -1 {
		// Shouldn't happen, return as-is
		return sequence
	}

	// Create a new sequence with space for the first strike step
	newSequence := make([]turnEntry, len(sequence)+1)
	copy(newSequence, sequence[:damageIdx])
	newSequence[damageIdx] = turnEntry{PhaseCombat, StepFirstStrikeDamage}
	copy(newSequence[damageIdx+1:], sequence[damageIdx:])

	return newSequence
}

// TurnManager tracks active/priority player and turn progression.
type TurnManager struct {
	orderIndex      int
	turnNumber      int
	activePlayer    string
	priorityPlayer  string
	sequence        []turnEntry // Dynamic turn sequence
	hasFirstStrike  bool         // Whether current turn sequence includes first strike step
}

// NewTurnManager creates a new turn manager initialized at turn 1, untap step.
func NewTurnManager(activePlayer string) *TurnManager {
	active := strings.TrimSpace(activePlayer)
	return &TurnManager{
		orderIndex:     0,
		turnNumber:     1,
		activePlayer:   active,
		priorityPlayer: active,
		sequence:       buildTurnSequence(false), // Start without first strike step
		hasFirstStrike: false,
	}
}

// CurrentPhase returns the phase currently in progress.
func (tm *TurnManager) CurrentPhase() Phase {
	return tm.sequence[tm.orderIndex].phase
}

// CurrentStep returns the step currently in progress.
func (tm *TurnManager) CurrentStep() Step {
	return tm.sequence[tm.orderIndex].step
}

// TurnNumber returns the current turn number (1-based).
func (tm *TurnManager) TurnNumber() int {
	return tm.turnNumber
}

// ActivePlayer returns the player who currently has the turn.
func (tm *TurnManager) ActivePlayer() string {
	return tm.activePlayer
}

// PriorityPlayer returns the player who currently has priority.
func (tm *TurnManager) PriorityPlayer() string {
	return tm.priorityPlayer
}

// SetPriority sets the player who currently has priority.
func (tm *TurnManager) SetPriority(player string) {
	tm.priorityPlayer = strings.TrimSpace(player)
}

// AdvanceStep advances to the next step in the turn structure.
// When the end of the structure is reached, the turn number is incremented
// and the active player is rotated to nextActivePlayer if provided.
func (tm *TurnManager) AdvanceStep(nextActivePlayer string) (Phase, Step) {
	tm.orderIndex++
	if tm.orderIndex >= len(tm.sequence) {
		tm.orderIndex = 0
		tm.turnNumber++
		if next := strings.TrimSpace(nextActivePlayer); next != "" {
			tm.activePlayer = next
		}
		// Reset sequence for new turn (no first strike by default)
		tm.sequence = buildTurnSequence(false)
		tm.hasFirstStrike = false
	}

	// Priority always reverts to active player at the start of a step.
	tm.priorityPlayer = tm.activePlayer

	return tm.CurrentPhase(), tm.CurrentStep()
}

// SetHasFirstStrike updates the turn sequence to include/exclude first strike damage step.
// This should be called after StepBeginCombat when creatures with first strike are declared.
func (tm *TurnManager) SetHasFirstStrike(hasFirstStrike bool) {
	if tm.hasFirstStrike == hasFirstStrike {
		return // No change
	}

	// We're currently at StepDeclareBlockers, need to rebuild sequence
	// First, save current state
	oldOrderIndex := tm.orderIndex

	// Rebuild the sequence
	newSequence := buildTurnSequence(hasFirstStrike)

	if !tm.hasFirstStrike && hasFirstStrike {
		// Going from no-first-strike to first-strike: orderIndex stays same
		// (the first strike step will be inserted AFTER current declare blockers)
		tm.orderIndex = oldOrderIndex
	} else if tm.hasFirstStrike && !hasFirstStrike {
		// Going from first-strike to no-first-strike: need to adjust orderIndex
		// if we're at or past the first strike step
		if oldOrderIndex >= len(newSequence) {
			tm.orderIndex = len(newSequence) - 1
		}
	}

	tm.sequence = newSequence
	tm.hasFirstStrike = hasFirstStrike
}

// GetSequence returns the current turn sequence for testing/inspection
func (tm *TurnManager) GetSequence() []turnEntry {
	return tm.sequence
}

// turnEntry needs Step() method for testing
type turnEntryWrapper interface {
	Step() Step
}

// Implement Step() on turnEntry
func (te turnEntry) Step() Step {
	return te.step
}
