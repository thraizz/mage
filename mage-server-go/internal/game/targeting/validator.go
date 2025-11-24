package targeting

import (
	"fmt"
	"strings"
)

// TargetValidator validates that selected targets are legal.
type TargetValidator struct {
	gameState TargetGameStateAccessor
}

// TargetGameStateAccessor provides access to game state needed for target validation.
type TargetGameStateAccessor interface {
	// FindCardForTarget finds a card by ID in any zone
	FindCardForTarget(cardID string) (TargetCardInfo, bool)
	// FindPlayerForTarget finds player info by ID
	FindPlayerForTarget(playerID string) (TargetPlayerInfo, bool)
	// IsCardInZone checks if a card is in a specific zone
	IsCardInZone(cardID string, zone int) bool
	// GetCardZone returns the zone a card is currently in
	GetCardZone(cardID string) (int, bool)
	// GetStackItemsForTarget returns all items currently on the stack
	GetStackItemsForTarget() []TargetStackItem
	// HasKeywordAbility checks if a card has a specific keyword ability
	HasKeywordAbility(cardID, keyword string) bool
	// GetProtectionQualities returns protection qualities (colors/types) for a card
	GetProtectionQualities(cardID string) []string
	// GetCardColor returns the color(s) of a card
	GetCardColor(cardID string) []string
}

// TargetCardInfo provides information about a card for target validation.
type TargetCardInfo struct {
	ID           string
	Name         string
	Type         string
	Zone         int
	ControllerID string
	OwnerID      string
	Tapped       bool
	FaceDown     bool
}

// TargetPlayerInfo provides information about a player for target validation.
type TargetPlayerInfo struct {
	PlayerID string
	Name     string
	Life     int
	Lost     bool
	Left     bool
}

// TargetStackItem provides information about a stack item for target validation.
type TargetStackItem struct {
	ID         string
	Controller string
	Kind       string
}

// NewTargetValidator creates a new target validator.
func NewTargetValidator(gameState TargetGameStateAccessor) *TargetValidator {
	return &TargetValidator{
		gameState: gameState,
	}
}

// ValidateTarget checks if a single target ID is valid for the given requirement.
func (tv *TargetValidator) ValidateTarget(targetID string, requirement TargetRequirement) error {
	if tv == nil || tv.gameState == nil {
		return fmt.Errorf("target validator not initialized")
	}

	// Check if target is a player
	player, isPlayer := tv.gameState.FindPlayerForTarget(targetID)
	if isPlayer {
		if requirement.Type != TargetTypePlayer {
			return fmt.Errorf("target %s is a player but requirement is %s", targetID, requirement.Type)
		}
		if player.Lost || player.Left {
			return fmt.Errorf("target player %s has left or lost the game", targetID)
		}
		return nil
	}

	// Check if target is a card
	card, isCard := tv.gameState.FindCardForTarget(targetID)
	if !isCard {
		return fmt.Errorf("target %s not found", targetID)
	}

	// Validate card type matches requirement
	switch requirement.Type {
	case TargetTypeCreature:
		if !strings.Contains(strings.ToLower(card.Type), "creature") {
			return fmt.Errorf("target %s is not a creature", card.Name)
		}
	case TargetTypeSpell:
		// Spells must be on the stack
		if card.Zone != 4 { // zoneStack
			return fmt.Errorf("target %s is not on the stack", card.Name)
		}
	case TargetTypePermanent:
		if card.Zone != 1 { // zoneBattlefield
			return fmt.Errorf("target %s is not a permanent", card.Name)
		}
	case TargetTypeArtifact:
		if !strings.Contains(strings.ToLower(card.Type), "artifact") {
			return fmt.Errorf("target %s is not an artifact", card.Name)
		}
	case TargetTypeEnchantment:
		if !strings.Contains(strings.ToLower(card.Type), "enchantment") {
			return fmt.Errorf("target %s is not an enchantment", card.Name)
		}
	case TargetTypeLand:
		if !strings.Contains(strings.ToLower(card.Type), "land") {
			return fmt.Errorf("target %s is not a land", card.Name)
		}
	case TargetTypePlaneswalker:
		if !strings.Contains(strings.ToLower(card.Type), "planeswalker") {
			return fmt.Errorf("target %s is not a planeswalker", card.Name)
		}
	case TargetTypePlayer:
		return fmt.Errorf("target %s is a card but requirement is player", card.Name)
	}

	// Check for targeting restrictions (hexproof, shroud, protection)
	if err := tv.checkTargetingRestrictions(card, requirement); err != nil {
		return err
	}

	return nil
}

// checkTargetingRestrictions checks if a card can be legally targeted
// Implements Rule 702.11 (Hexproof), 702.18 (Shroud), 702.16 (Protection)
func (tv *TargetValidator) checkTargetingRestrictions(card TargetCardInfo, requirement TargetRequirement) error {
	// Rule 702.18: Shroud - can't be the target of spells or abilities
	if tv.gameState.HasKeywordAbility(card.ID, "SHROUD") {
		return fmt.Errorf("target %s has shroud and can't be targeted", card.Name)
	}

	// Rule 702.11: Hexproof - can't be targeted by spells or abilities opponents control
	if tv.gameState.HasKeywordAbility(card.ID, "HEXPROOF") {
		if requirement.ControllerID != "" && requirement.ControllerID != card.ControllerID {
			return fmt.Errorf("target %s has hexproof and can't be targeted by opponents", card.Name)
		}
	}

	// Rule 702.16: Protection - can't be targeted by sources with the protected quality
	protectionQualities := tv.gameState.GetProtectionQualities(card.ID)
	if len(protectionQualities) > 0 && requirement.SourceID != "" {
		// Check if source matches any protection quality
		sourceColors := tv.gameState.GetCardColor(requirement.SourceID)

		for _, quality := range protectionQualities {
			// Check color protection
			for _, color := range sourceColors {
				if strings.EqualFold(quality, color) {
					return fmt.Errorf("target %s has protection from %s", card.Name, quality)
				}
			}

			// Check type protection (e.g., "protection from creatures")
			sourceInfo, found := tv.gameState.FindCardForTarget(requirement.SourceID)
			if found {
				if strings.Contains(strings.ToLower(sourceInfo.Type), strings.ToLower(quality)) {
					return fmt.Errorf("target %s has protection from %s", card.Name, quality)
				}
			}
		}
	}

	return nil
}

// ValidateTargetSelection validates an entire target selection against its requirements.
func (tv *TargetValidator) ValidateTargetSelection(selection *TargetSelection) error {
	if tv == nil {
		return fmt.Errorf("target validator not initialized")
	}

	if selection == nil {
		return fmt.Errorf("target selection is nil")
	}

	// Check if selection meets requirement counts
	if err := selection.Validate(); err != nil {
		return err
	}

	// Validate each individual target
	for _, targetID := range selection.Targets {
		if err := tv.ValidateTarget(targetID, selection.Requirement); err != nil {
			return fmt.Errorf("invalid target %s: %v", targetID, err)
		}
	}

	// Check for duplicate targets
	seen := make(map[string]bool)
	for _, targetID := range selection.Targets {
		if seen[targetID] {
			return fmt.Errorf("duplicate target: %s", targetID)
		}
		seen[targetID] = true
	}

	return nil
}
