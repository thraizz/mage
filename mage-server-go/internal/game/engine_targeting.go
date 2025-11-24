package game

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap"
)

// TargetSelectionManager handles target selection for spells and abilities
type TargetSelectionManager struct {
	logger *zap.Logger
}

// NewTargetSelectionManager creates a new target selection manager
func NewTargetSelectionManager(logger *zap.Logger) *TargetSelectionManager {
	return &TargetSelectionManager{
		logger: logger,
	}
}

// TargetRequest represents a request for the player to select targets
type TargetRequest struct {
	AbilityID    uuid.UUID
	SourceID     uuid.UUID
	PlayerID     uuid.UUID
	MinTargets   int
	MaxTargets   int
	TargetFilter abilities.TargetFilter
	LegalTargets []uuid.UUID // Precalculated legal targets
	Message      string
}

// TargetSelection represents the player's target choices
type TargetSelection struct {
	Targets []uuid.UUID
}

// ValidateTargets validates that the provided targets are legal
func (tsm *TargetSelectionManager) ValidateTargets(
	ctx context.Context,
	request *TargetRequest,
	selection *TargetSelection,
	state rules.GameStateReader,
	gameCtx abilities.GameContext,
) error {
	// Check number of targets
	numTargets := len(selection.Targets)
	if numTargets < request.MinTargets {
		return fmt.Errorf("not enough targets: need at least %d, got %d", request.MinTargets, numTargets)
	}
	if numTargets > request.MaxTargets {
		return fmt.Errorf("too many targets: maximum %d, got %d", request.MaxTargets, numTargets)
	}

	// Check each target is legal
	for _, targetID := range selection.Targets {
		if !tsm.isLegalTarget(ctx, targetID, request.TargetFilter, state, gameCtx) {
			return fmt.Errorf("illegal target: %s", targetID.String())
		}
	}

	// Check for duplicate targets (most spells don't allow this)
	if tsm.hasDuplicates(selection.Targets) {
		return fmt.Errorf("duplicate targets not allowed")
	}

	return nil
}

// GetLegalTargets returns all legal targets for a given filter
func (tsm *TargetSelectionManager) GetLegalTargets(
	ctx context.Context,
	filter abilities.TargetFilter,
	state rules.GameStateReader,
	gameCtx abilities.GameContext,
) []uuid.UUID {
	legalTargets := make([]uuid.UUID, 0)

	// For now, we'll return a simplified implementation
	// TODO: Implement full target legality checking based on filter type

	tsm.logger.Debug("calculating legal targets",
		zap.String("filter", filter.GetDescription()))

	// This would need to:
	// 1. Get all permanents/players in game
	// 2. Filter by the target filter criteria
	// 3. Check for targeting restrictions (hexproof, shroud, protection, etc.)
	// 4. Return list of legal target IDs

	return legalTargets
}

// isLegalTarget checks if a single target is legal
func (tsm *TargetSelectionManager) isLegalTarget(
	ctx context.Context,
	targetID uuid.UUID,
	filter abilities.TargetFilter,
	state rules.GameStateReader,
	gameCtx abilities.GameContext,
) bool {
	// Check if target exists
	permanent, hasPermanent := state.GetPermanent(targetID)
	player, hasPlayer := state.GetPlayer(targetID)

	if !hasPermanent && !hasPlayer {
		return false
	}

	// Check if filter matches
	if hasPermanent {
		// Convert permanent to a format the filter can check
		// TODO: Implement proper permanent to filter checking
		_ = permanent
		// return filter.Matches(permanent)
	}

	if hasPlayer {
		// Check if filter accepts players
		// TODO: Implement player filter checking
		_ = player
	}

	// Check targeting restrictions
	// TODO: Implement hexproof, shroud, protection checks

	return true
}

// hasDuplicates checks if a slice has duplicate UUIDs
func (tsm *TargetSelectionManager) hasDuplicates(targets []uuid.UUID) bool {
	seen := make(map[uuid.UUID]bool)
	for _, target := range targets {
		if seen[target] {
			return true
		}
		seen[target] = true
	}
	return false
}

// CreateTargetRequest creates a target request from an ability
func (tsm *TargetSelectionManager) CreateTargetRequest(
	ctx context.Context,
	ability abilities.Ability,
	playerID uuid.UUID,
	state rules.GameStateReader,
	gameCtx abilities.GameContext,
) (*TargetRequest, error) {
	// Get target requirements from the ability
	// For activated abilities, get the target requirement
	var targetReq *abilities.TargetRequirement
	var targetFilter abilities.TargetFilter

	switch a := ability.(type) {
	case *abilities.ActivatedAbility:
		targetReq = a.Targets
	case *abilities.SpellAbility:
		targetReq = a.GetTargets()
	case *abilities.TriggeredAbility:
		targetReq = a.GetTargets()
	default:
		return nil, fmt.Errorf("ability type does not support targeting")
	}

	if targetReq == nil {
		return nil, fmt.Errorf("ability has no target requirement")
	}

	// Extract filter from requirement
	// TODO: Get filter from TargetRequirement
	// For now, assume we have a way to get it
	// targetFilter = targetReq.GetFilter()

	// Calculate legal targets
	legalTargets := tsm.GetLegalTargets(ctx, targetFilter, state, gameCtx)

	request := &TargetRequest{
		AbilityID:    ability.GetID(),
		SourceID:     ability.GetSourceID(),
		PlayerID:     playerID,
		MinTargets:   1, // TODO: Get from targetReq
		MaxTargets:   1, // TODO: Get from targetReq
		TargetFilter: targetFilter,
		LegalTargets: legalTargets,
		Message:      fmt.Sprintf("Choose target for %s", ability.String()),
	}

	tsm.logger.Info("created target request",
		zap.String("ability", ability.GetID().String()),
		zap.Int("legal_targets", len(legalTargets)))

	return request, nil
}

// AutoSelectTargets automatically selects targets when there's only one legal choice
// Returns nil if automatic selection is not possible
func (tsm *TargetSelectionManager) AutoSelectTargets(
	ctx context.Context,
	request *TargetRequest,
) *TargetSelection {
	// If exactly one target is required and exactly one is legal, auto-select
	if request.MinTargets == 1 && request.MaxTargets == 1 && len(request.LegalTargets) == 1 {
		tsm.logger.Info("auto-selecting single legal target",
			zap.String("target", request.LegalTargets[0].String()))

		return &TargetSelection{
			Targets: []uuid.UUID{request.LegalTargets[0]},
		}
	}

	// If no targets are required and none are provided, auto-select empty
	if request.MinTargets == 0 && len(request.LegalTargets) == 0 {
		tsm.logger.Debug("auto-selecting no targets")
		return &TargetSelection{
			Targets: []uuid.UUID{},
		}
	}

	return nil
}

// CheckTargetingRestrictions checks if a permanent can be targeted
// This includes checking for hexproof, shroud, protection, etc.
func (tsm *TargetSelectionManager) CheckTargetingRestrictions(
	ctx context.Context,
	permanentID uuid.UUID,
	sourceID uuid.UUID,
	controller uuid.UUID,
	state rules.GameStateReader,
) error {
	permanent, ok := state.GetPermanent(permanentID)
	if !ok {
		return fmt.Errorf("permanent not found")
	}

	// TODO: Implement full targeting restriction checks:
	// 1. Hexproof - can't be targeted by opponents' spells/abilities
	// 2. Shroud - can't be targeted by any spells/abilities
	// 3. Protection from [quality] - can't be targeted by sources with that quality
	// 4. Ward - requires additional cost to target

	_ = permanent
	_ = sourceID
	_ = controller

	return nil
}

// DivideAmongTargets handles division effects (e.g., "distribute 3 damage among any number of targets")
type DivisionRequest struct {
	TotalAmount  int
	Targets      []uuid.UUID
	MinPerTarget int
	MaxPerTarget int
}

type DivisionSelection struct {
	Amounts map[uuid.UUID]int // targetID -> amount
}

// ValidateDivision validates a division selection
func (tsm *TargetSelectionManager) ValidateDivision(
	request *DivisionRequest,
	selection *DivisionSelection,
) error {
	totalAssigned := 0

	// Check each target
	for _, targetID := range request.Targets {
		amount, ok := selection.Amounts[targetID]
		if !ok {
			amount = 0
		}

		// Check min/max per target
		if amount < request.MinPerTarget {
			return fmt.Errorf("target %s: amount %d less than minimum %d",
				targetID.String(), amount, request.MinPerTarget)
		}
		if request.MaxPerTarget > 0 && amount > request.MaxPerTarget {
			return fmt.Errorf("target %s: amount %d exceeds maximum %d",
				targetID.String(), amount, request.MaxPerTarget)
		}

		totalAssigned += amount
	}

	// Check total
	if totalAssigned != request.TotalAmount {
		return fmt.Errorf("total assigned %d does not match required %d",
			totalAssigned, request.TotalAmount)
	}

	return nil
}
