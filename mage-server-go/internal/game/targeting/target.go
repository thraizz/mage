package targeting

import (
	"fmt"
	"strings"
)

// TargetType represents the type of target a spell or ability can have.
type TargetType string

const (
	// TargetTypeCreature targets creatures
	TargetTypeCreature TargetType = "CREATURE"
	// TargetTypePlayer targets players
	TargetTypePlayer TargetType = "PLAYER"
	// TargetTypeSpell targets spells on the stack
	TargetTypeSpell TargetType = "SPELL"
	// TargetTypePermanent targets permanents (creatures, artifacts, enchantments, etc.)
	TargetTypePermanent TargetType = "PERMANENT"
	// TargetTypeArtifact targets artifacts
	TargetTypeArtifact TargetType = "ARTIFACT"
	// TargetTypeEnchantment targets enchantments
	TargetTypeEnchantment TargetType = "ENCHANTMENT"
	// TargetTypeLand targets lands
	TargetTypeLand TargetType = "LAND"
	// TargetTypePlaneswalker targets planeswalkers
	TargetTypePlaneswalker TargetType = "PLANESWALKER"
)

// TargetRequirement defines what targets a spell or ability requires.
type TargetRequirement struct {
	// Type specifies what kind of target is required
	Type TargetType
	// MinTargets is the minimum number of targets required (usually 1)
	MinTargets int
	// MaxTargets is the maximum number of targets allowed (usually 1, but can be "up to X")
	MaxTargets int
	// Optional indicates if targets are optional (e.g., "up to X targets")
	Optional bool
	// Description is a human-readable description of the target requirement
	Description string
}

// TargetSelection represents a player's target selection for a spell or ability.
type TargetSelection struct {
	// Targets is a list of target IDs (can be card IDs or player IDs)
	Targets []string
	// Requirement is the requirement this selection satisfies
	Requirement TargetRequirement
}

// IsComplete checks if the target selection meets the requirement.
func (ts *TargetSelection) IsComplete() bool {
	if ts == nil {
		return false
	}
	count := len(ts.Targets)
	if ts.Requirement.Optional {
		return count >= ts.Requirement.MinTargets && count <= ts.Requirement.MaxTargets
	}
	return count >= ts.Requirement.MinTargets && count <= ts.Requirement.MaxTargets
}

// Validate checks if the target selection is valid.
func (ts *TargetSelection) Validate() error {
	if ts == nil {
		return fmt.Errorf("target selection is nil")
	}
	count := len(ts.Targets)
	if count < ts.Requirement.MinTargets {
		return fmt.Errorf("not enough targets: need at least %d, got %d", ts.Requirement.MinTargets, count)
	}
	if count > ts.Requirement.MaxTargets {
		return fmt.Errorf("too many targets: need at most %d, got %d", ts.Requirement.MaxTargets, count)
	}
	return nil
}

// ParseTargetRequirements parses target requirements from a card's rules text or metadata.
// This is a simplified parser - full implementation would parse actual card text.
func ParseTargetRequirements(cardType string, rulesText string) []TargetRequirement {
	requirements := []TargetRequirement{}
	
	// Simple heuristic-based parsing
	// In a full implementation, this would parse actual card rules text
	text := strings.ToLower(rulesText)
	
	// Check for common targeting patterns
	if strings.Contains(text, "target creature") {
		requirements = append(requirements, TargetRequirement{
			Type:        TargetTypeCreature,
			MinTargets:  1,
			MaxTargets:  1,
			Optional:    false,
			Description: "target creature",
		})
	}
	if strings.Contains(text, "target player") {
		requirements = append(requirements, TargetRequirement{
			Type:        TargetTypePlayer,
			MinTargets:  1,
			MaxTargets:  1,
			Optional:    false,
			Description: "target player",
		})
	}
	if strings.Contains(text, "target spell") {
		requirements = append(requirements, TargetRequirement{
			Type:        TargetTypeSpell,
			MinTargets:  1,
			MaxTargets:  1,
			Optional:    false,
			Description: "target spell",
		})
	}
	if strings.Contains(text, "target permanent") {
		requirements = append(requirements, TargetRequirement{
			Type:        TargetTypePermanent,
			MinTargets:  1,
			MaxTargets:  1,
			Optional:    false,
			Description: "target permanent",
		})
	}
	
	// Check for "up to X targets" patterns
	if strings.Contains(text, "up to") {
		for _, req := range requirements {
			req.Optional = true
			req.MinTargets = 0
		}
	}
	
	return requirements
}

// FormatTargets formats target IDs into a human-readable string for metadata storage.
func FormatTargets(targets []string) string {
	if len(targets) == 0 {
		return ""
	}
	return strings.Join(targets, ",")
}

// ParseTargets parses target IDs from a formatted string.
func ParseTargets(formatted string) []string {
	if formatted == "" {
		return []string{}
	}
	return strings.Split(formatted, ",")
}
