package game

import (
	"context"
	"fmt"

	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap"
)

// PriorityManager handles priority passing and state-based action checks
// This integrates the SBA system with the turn structure
type PriorityManager struct {
	turnMgr    *rules.TurnManager
	sbaChecker *rules.StateBasedActions
	layerMgr   *ContinuousEffectsManager
	logger     *zap.Logger

	// Priority passing state
	passedCount  int // How many consecutive players have passed
	totalPlayers int // Total number of players in game
}

// NewPriorityManager creates a new priority manager
func NewPriorityManager(turnMgr *rules.TurnManager, logger *zap.Logger) *PriorityManager {
	return &PriorityManager{
		turnMgr:      turnMgr,
		sbaChecker:   rules.NewStateBasedActions(),
		logger:       logger,
		passedCount:  0,
		totalPlayers: 2, // TODO: Get from game state
	}
}

// SetLayerManager sets the continuous effects layer manager
// This must be called after initialization to enable layer recalculation
func (pm *PriorityManager) SetLayerManager(layerMgr *ContinuousEffectsManager) {
	pm.layerMgr = layerMgr
}

// CheckStateBasedActions performs SBA checks and executes all resulting actions
// This should be called:
// 1. When a player would receive priority
// 2. After a spell/ability finishes resolving
func (pm *PriorityManager) CheckStateBasedActions(ctx context.Context, state rules.GameStateReader) error {
	pm.logger.Debug("checking state-based actions")

	// 1. Recalculate continuous effects layers (Rule 613)
	// This must happen BEFORE checking SBAs, as SBAs depend on current characteristics
	if pm.layerMgr != nil {
		if err := pm.layerMgr.RecalculateAll(ctx, state); err != nil {
			pm.logger.Error("failed to recalculate layers", zap.Error(err))
			return fmt.Errorf("layer recalculation failed: %w", err)
		}
	}

	// 2. Check for state-based actions
	actions := pm.sbaChecker.Check(state)

	if len(actions) == 0 {
		pm.logger.Debug("no state-based actions to perform")
		return nil
	}

	pm.logger.Info("performing state-based actions",
		zap.Int("count", len(actions)))

	// 3. Execute all actions
	// Note: SBAs are executed simultaneously, not one at a time
	for _, action := range actions {
		pm.logger.Debug("executing SBA",
			zap.String("action", action.GetDescription()))

		if err := action.Execute(state); err != nil {
			pm.logger.Error("failed to execute SBA",
				zap.String("action", action.GetDescription()),
				zap.Error(err))
			return fmt.Errorf("failed to execute SBA: %w", err)
		}
	}

	// 4. After performing SBAs, check again (Rule 704.3)
	// If any SBAs were performed, repeat the check
	if len(actions) > 0 {
		pm.logger.Debug("rechecking SBAs after execution")
		return pm.CheckStateBasedActions(ctx, state)
	}

	return nil
}

// GivePriority gives priority to the specified player
// This performs SBA checks before giving priority (Rule 117.5)
func (pm *PriorityManager) GivePriority(ctx context.Context, playerID string, state rules.GameStateReader) error {
	pm.logger.Debug("giving priority",
		zap.String("player", playerID),
		zap.String("phase", pm.turnMgr.CurrentPhase().String()),
		zap.String("step", pm.turnMgr.CurrentStep().String()))

	// 1. Check state-based actions (Rule 704.3)
	if err := pm.CheckStateBasedActions(ctx, state); err != nil {
		return fmt.Errorf("SBA check failed: %w", err)
	}

	// 2. Check for triggered abilities (TODO: implement)
	// if err := pm.checkTriggeredAbilities(ctx, state); err != nil {
	// 	return fmt.Errorf("trigger check failed: %w", err)
	// }

	// 3. Set priority player
	pm.turnMgr.SetPriority(playerID)
	pm.passedCount = 0 // Reset passed count when giving priority

	return nil
}

// PassPriority handles a player passing priority
// Returns true if all players have passed and the stack should resolve/step should advance
func (pm *PriorityManager) PassPriority(ctx context.Context, state rules.GameStateReader) (bool, error) {
	pm.logger.Debug("player passed priority",
		zap.String("player", pm.turnMgr.PriorityPlayer()))

	pm.passedCount++

	// If all players have passed, we're done with this priority round
	if pm.passedCount >= pm.totalPlayers {
		pm.logger.Debug("all players passed priority")
		return true, nil
	}

	// Pass priority to next player
	// TODO: Get next player from game state
	nextPlayer := pm.getNextPlayer()

	// Give priority to next player (this will check SBAs)
	if err := pm.GivePriority(ctx, nextPlayer, state); err != nil {
		return false, fmt.Errorf("failed to give priority: %w", err)
	}

	return false, nil
}

// AfterSpellResolves should be called after a spell/ability resolves from the stack
// This performs SBA checks and triggers check (Rule 608.2k)
func (pm *PriorityManager) AfterSpellResolves(ctx context.Context, state rules.GameStateReader) error {
	pm.logger.Debug("spell/ability resolved, checking SBAs")

	// Check state-based actions after resolution
	if err := pm.CheckStateBasedActions(ctx, state); err != nil {
		return fmt.Errorf("post-resolution SBA check failed: %w", err)
	}

	// Check for triggered abilities
	// TODO: Implement trigger checking

	return nil
}

// AdvanceToNextStep advances to the next step and gives priority to the active player
func (pm *PriorityManager) AdvanceToNextStep(ctx context.Context, state rules.GameStateReader, nextActivePlayer string) error {
	phase, step := pm.turnMgr.AdvanceStep(nextActivePlayer)

	pm.logger.Info("advancing to next step",
		zap.String("phase", phase.String()),
		zap.String("step", step.String()),
		zap.Int("turn", pm.turnMgr.TurnNumber()))

	// Perform step-specific actions
	if err := pm.performStepActions(ctx, step, state); err != nil {
		return fmt.Errorf("step actions failed: %w", err)
	}

	// Give priority to active player (unless it's untap step)
	if step != rules.StepUntap {
		if err := pm.GivePriority(ctx, pm.turnMgr.ActivePlayer(), state); err != nil {
			return fmt.Errorf("failed to give priority: %w", err)
		}
	}

	return nil
}

// performStepActions performs the automatic actions for each step
func (pm *PriorityManager) performStepActions(ctx context.Context, step rules.Step, state rules.GameStateReader) error {
	switch step {
	case rules.StepUntap:
		// Untap all permanents controlled by active player
		// TODO: Implement untap logic with static abilities that prevent untap
		pm.logger.Debug("performing untap step")

	case rules.StepUpkeep:
		// No automatic actions (triggers will be checked before priority)
		pm.logger.Debug("performing upkeep step")

	case rules.StepDraw:
		// Active player draws a card
		// TODO: Implement draw logic
		pm.logger.Debug("performing draw step")

	case rules.StepCleanup:
		// 1. Active player discards down to hand size
		// 2. Remove damage from all permanents
		// 3. Remove "until end of turn" effects
		// TODO: Implement cleanup logic
		pm.logger.Debug("performing cleanup step")
	}

	return nil
}

// getNextPlayer returns the next player in turn order
// TODO: This should be integrated with the game state to handle multiplayer properly
func (pm *PriorityManager) getNextPlayer() string {
	// Placeholder: For 2-player games, alternate players
	current := pm.turnMgr.PriorityPlayer()
	if current == pm.turnMgr.ActivePlayer() {
		return "opponent" // TODO: Get from game state
	}
	return pm.turnMgr.ActivePlayer()
}

// ResetPassCount resets the priority pass counter
// This should be called when a spell/ability is added to the stack
func (pm *PriorityManager) ResetPassCount() {
	pm.passedCount = 0
}

// GetCurrentStep returns the current step
func (pm *PriorityManager) GetCurrentStep() rules.Step {
	return pm.turnMgr.CurrentStep()
}

// GetCurrentPhase returns the current phase
func (pm *PriorityManager) GetCurrentPhase() rules.Phase {
	return pm.turnMgr.CurrentPhase()
}

// GetActivePlayer returns the active player
func (pm *PriorityManager) GetActivePlayer() string {
	return pm.turnMgr.ActivePlayer()
}

// GetPriorityPlayer returns the player who currently has priority
func (pm *PriorityManager) GetPriorityPlayer() string {
	return pm.turnMgr.PriorityPlayer()
}
