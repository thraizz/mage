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
	StepCombatDamage
	StepEndCombat
	StepMain2
	StepEnd
	StepCleanup
)

var stepNames = map[Step]string{
	StepUntap:            "UNTAP",
	StepUpkeep:           "UPKEEP",
	StepDraw:             "DRAW",
	StepMain1:            "MAIN1",
	StepBeginCombat:      "BEGIN_COMBAT",
	StepDeclareAttackers: "DECLARE_ATTACKERS",
	StepDeclareBlockers:  "DECLARE_BLOCKERS",
	StepCombatDamage:     "COMBAT_DAMAGE",
	StepEndCombat:        "END_COMBAT",
	StepMain2:            "MAIN2",
	StepEnd:              "END",
	StepCleanup:          "CLEANUP",
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

var turnSequence = []turnEntry{
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

// TurnManager tracks active/priority player and turn progression.
type TurnManager struct {
	orderIndex     int
	turnNumber     int
	activePlayer   string
	priorityPlayer string
}

// NewTurnManager creates a new turn manager initialized at turn 1, untap step.
func NewTurnManager(activePlayer string) *TurnManager {
	active := strings.TrimSpace(activePlayer)
	return &TurnManager{
		orderIndex:     0,
		turnNumber:     1,
		activePlayer:   active,
		priorityPlayer: active,
	}
}

// CurrentPhase returns the phase currently in progress.
func (tm *TurnManager) CurrentPhase() Phase {
	return turnSequence[tm.orderIndex].phase
}

// CurrentStep returns the step currently in progress.
func (tm *TurnManager) CurrentStep() Step {
	return turnSequence[tm.orderIndex].step
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
	if tm.orderIndex >= len(turnSequence) {
		tm.orderIndex = 0
		tm.turnNumber++
		if next := strings.TrimSpace(nextActivePlayer); next != "" {
			tm.activePlayer = next
		}
	}

	// Priority always reverts to active player at the start of a step.
	tm.priorityPlayer = tm.activePlayer

	return tm.CurrentPhase(), tm.CurrentStep()
}
