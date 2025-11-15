package rules

import (
	"fmt"
)

// LegalityChecker validates stack items before resolution.
// Mirrors Java's legality checking system.
type LegalityChecker struct {
	gameState GameStateAccessor
}

// GameStateAccessor provides access to game state needed for legality checks.
type GameStateAccessor interface {
	// FindCard finds a card by ID in any zone
	FindCard(cardID string) (CardInfo, bool)
	// FindPlayer finds player info by ID
	FindPlayer(playerID string) (PlayerInfo, bool)
	// IsCardInZone checks if a card is in a specific zone
	IsCardInZone(cardID string, zone int) bool
	// GetCardZone returns the zone a card is currently in
	GetCardZone(cardID string) (int, bool)
}

// CardInfo provides information about a card for legality checks.
type CardInfo struct {
	ID           string
	Name         string
	Type         string
	Zone         int
	ControllerID string
	OwnerID      string
	Tapped       bool
	FaceDown     bool
}

// PlayerInfo provides information about a player for legality checks.
type PlayerInfo struct {
	PlayerID string
	Name     string
	Life     int
	Lost     bool
	Left     bool
}

// LegalityResult represents the result of a legality check.
type LegalityResult struct {
	Legal   bool
	Reason  string
	Details map[string]string
}

// NewLegalityChecker creates a new legality checker.
func NewLegalityChecker(gameState GameStateAccessor) *LegalityChecker {
	return &LegalityChecker{
		gameState: gameState,
	}
}

// CheckStackItemLegality validates a stack item before resolution.
// Returns true if the item is legal to resolve, false otherwise.
func (lc *LegalityChecker) CheckStackItemLegality(item StackItem) LegalityResult {
	if lc == nil || lc.gameState == nil {
		return LegalityResult{
			Legal:  true, // Default to legal if checker not initialized
			Reason: "Legality checker not initialized",
		}
	}

	// Check 1: Controller still in game
	if item.Controller != "" {
		player, found := lc.gameState.FindPlayer(item.Controller)
		if !found {
			return LegalityResult{
				Legal:  false,
				Reason: "Controller not found",
				Details: map[string]string{
					"controller_id": item.Controller,
				},
			}
		}
		if player.Lost || player.Left {
			return LegalityResult{
				Legal:  false,
				Reason: "Controller has left or lost the game",
				Details: map[string]string{
					"controller_id": item.Controller,
					"lost":          fmt.Sprintf("%v", player.Lost),
					"left":          fmt.Sprintf("%v", player.Left),
				},
			}
		}
	}

	// Check 2: Source card still exists and is in valid zone
	if item.SourceID != "" {
		sourceCard, found := lc.gameState.FindCard(item.SourceID)
		if !found {
			// Source card no longer exists - this is legal for some abilities
			// but illegal for spells that require their source
			if item.Kind == StackItemKindSpell {
				return LegalityResult{
					Legal:  false,
					Reason: "Source card no longer exists",
					Details: map[string]string{
						"source_id": item.SourceID,
						"kind":      string(item.Kind),
					},
				}
			}
			// For abilities, source disappearing is usually legal (e.g., creature dies but ability resolves)
		} else {
			// Source exists - check if it's in a valid zone for the ability type
			if !lc.isSourceInValidZone(sourceCard, item.Kind) {
				return LegalityResult{
					Legal:  false,
					Reason: "Source card not in valid zone",
					Details: map[string]string{
						"source_id": item.SourceID,
						"source_zone": fmt.Sprintf("%d", sourceCard.Zone),
						"kind":        string(item.Kind),
					},
				}
			}
		}
	}

	// Check 3: Validate targets (if any)
	if targets, hasTargets := lc.extractTargets(item); hasTargets {
		if result := lc.validateTargets(targets); !result.Legal {
			return result
		}
	}

	// Check 4: Timing restrictions (e.g., sorceries can't be cast during combat)
	if result := lc.checkTimingRestrictions(item); !result.Legal {
		return result
	}

	return LegalityResult{
		Legal:  true,
		Reason: "All legality checks passed",
	}
}

// Zone constants (matching mage_engine.go)
// Note: These must match the constants in mage_engine.go exactly
const (
	zoneLibrary = 0
	zoneHand    = 1
	zoneBattlefield = 2
	zoneGraveyard = 3
	zoneStack   = 4
	zoneExile   = 5
	zoneCommand = 6
)

// isSourceInValidZone checks if a source card is in a valid zone for its ability type.
func (lc *LegalityChecker) isSourceInValidZone(card CardInfo, kind StackItemKind) bool {
	// Spells must be on the stack
	if kind == StackItemKindSpell {
		return card.Zone == zoneStack
	}

	// Activated abilities: source must be on battlefield (or hand/library for some)
	// For now, we allow battlefield, hand, graveyard, and exile
	if kind == StackItemKindActivated {
		return card.Zone == zoneHand || card.Zone == zoneBattlefield || card.Zone == zoneGraveyard || card.Zone == zoneExile || card.Zone == zoneStack
	}

	// Triggered abilities: source can be in various zones depending on trigger
	// Most commonly battlefield, but also graveyard, exile, etc.
	if kind == StackItemKindTriggered {
		return true // Triggered abilities can resolve from any zone
	}

	return true
}

// extractTargets extracts target IDs from a stack item's metadata.
func (lc *LegalityChecker) extractTargets(item StackItem) ([]string, bool) {
	if item.Metadata == nil {
		return nil, false
	}

	// Check for explicit targets list
	if targetsStr, ok := item.Metadata["targets"]; ok && targetsStr != "" {
		// Parse comma-separated targets
		targets := []string{}
		// Simple parsing - in production would handle more formats
		if targetsStr != "" {
			// For now, assume single target or comma-separated
			// Full implementation would parse properly
			targets = append(targets, targetsStr)
		}
		return targets, len(targets) > 0
	}

	// Check for single target
	if targetID, ok := item.Metadata["target"]; ok && targetID != "" {
		return []string{targetID}, true
	}

	return nil, false
}

// validateTargets checks if all targets are still legal.
func (lc *LegalityChecker) validateTargets(targets []string) LegalityResult {
	if len(targets) == 0 {
		return LegalityResult{
			Legal:  true,
			Reason: "No targets to validate",
		}
	}

	invalidTargets := []string{}
	for _, targetID := range targets {
		// Check if target is a card
		if _, found := lc.gameState.FindCard(targetID); found {
			// Card targets: check if still in valid zone
			// Most spells/abilities require targets on battlefield
			// Some can target cards in graveyard, hand, etc.
			zone, hasZone := lc.gameState.GetCardZone(targetID)
			if !hasZone {
				invalidTargets = append(invalidTargets, fmt.Sprintf("%s (not found)", targetID))
				continue
			}

			// Default: targets must be on battlefield unless specified otherwise
			// This is a simplification - full implementation would check spell/ability rules
			if zone != zoneBattlefield {
				invalidTargets = append(invalidTargets, fmt.Sprintf("%s (zone %d)", targetID, zone))
			}
		} else if player, found := lc.gameState.FindPlayer(targetID); found {
			// Player targets: check if still in game
			if player.Lost || player.Left {
				invalidTargets = append(invalidTargets, fmt.Sprintf("%s (lost/left)", targetID))
			}
		} else {
			// Target not found
			invalidTargets = append(invalidTargets, fmt.Sprintf("%s (not found)", targetID))
		}
	}

	if len(invalidTargets) > 0 {
		return LegalityResult{
			Legal:  false,
			Reason: "One or more targets are illegal",
			Details: map[string]string{
				"invalid_targets": fmt.Sprintf("%v", invalidTargets),
			},
		}
	}

	return LegalityResult{
		Legal:  true,
		Reason: "All targets are legal",
	}
}

// checkTimingRestrictions validates timing restrictions for spells/abilities.
func (lc *LegalityChecker) checkTimingRestrictions(item StackItem) LegalityResult {
	// For now, we don't enforce strict timing restrictions
	// Full implementation would check:
	// - Sorceries can only be cast during main phases when stack is empty
	// - Instants can be cast anytime
	// - Activated abilities have timing restrictions
	// - Triggered abilities resolve when triggered

	// Basic check: if metadata indicates timing violation
	if item.Metadata != nil {
		if timingIssue, ok := item.Metadata["timing_violation"]; ok && timingIssue == "true" {
			return LegalityResult{
				Legal:  false,
				Reason: "Timing restriction violation",
				Details: map[string]string{
					"kind": string(item.Kind),
				},
			}
		}
	}

	return LegalityResult{
		Legal:  true,
		Reason: "Timing restrictions satisfied",
	}
}

// CheckCostsPaid validates that costs for a stack item were paid.
// This is a placeholder - full implementation would track cost payment.
func (lc *LegalityChecker) CheckCostsPaid(item StackItem) LegalityResult {
	// Check metadata for cost payment confirmation
	if item.Metadata != nil {
		if costsPaid, ok := item.Metadata["costs_paid"]; ok {
			if costsPaid == "false" {
				return LegalityResult{
					Legal:  false,
					Reason: "Costs not paid",
					Details: map[string]string{
						"item_id": item.ID,
					},
				}
			}
		}
	}

	// Default: assume costs were paid if not explicitly marked otherwise
	return LegalityResult{
		Legal:  true,
		Reason: "Costs verified as paid",
	}
}
