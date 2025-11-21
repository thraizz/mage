package abilities

import (
	"fmt"

	"github.com/google/uuid"
)

// Target represents something that can be targeted by a spell or ability
type Target interface {
	// GetID returns the ID of the target
	GetID() uuid.UUID

	// IsValid checks if this is a valid target
	IsValid() bool

	// String returns a text description
	String() string
}

// TargetRequirement specifies what can be targeted
type TargetRequirement struct {
	// MinTargets is the minimum number of targets required
	MinTargets int

	// MaxTargets is the maximum number of targets allowed
	MaxTargets int

	// Filter specifies what types of things can be targeted
	Filter TargetFilter

	// Description is a human-readable description
	Description string
}

// NewTargetRequirement creates a new target requirement
func NewTargetRequirement(min, max int, filter TargetFilter) *TargetRequirement {
	return &TargetRequirement{
		MinTargets:  min,
		MaxTargets:  max,
		Filter:      filter,
		Description: filter.GetDescription(),
	}
}

// TargetFilter determines what can be targeted
type TargetFilter interface {
	// Matches checks if a target matches this filter
	Matches(targetID uuid.UUID, game GameContext) bool

	// GetDescription returns a description of what this filter matches
	GetDescription() string
}

// ========================================
// Any Target
// ========================================

// AnyTargetFilter matches any target (creature, player, planeswalker, or battle)
type AnyTargetFilter struct{}

func NewAnyTargetFilter() *AnyTargetFilter {
	return &AnyTargetFilter{}
}

func (f *AnyTargetFilter) Matches(targetID uuid.UUID, game GameContext) bool {
	// TODO: Check if target is a creature, player, planeswalker, or battle
	return true
}

func (f *AnyTargetFilter) GetDescription() string {
	return "any target"
}

// ========================================
// Creature Target
// ========================================

// CreatureTargetFilter matches only creatures
type CreatureTargetFilter struct {
	Subtype string // Optional subtype filter (e.g., "Human", "Wizard")
}

func NewCreatureTargetFilter() *CreatureTargetFilter {
	return &CreatureTargetFilter{}
}

func NewCreatureTargetFilterWithSubtype(subtype string) *CreatureTargetFilter {
	return &CreatureTargetFilter{Subtype: subtype}
}

func (f *CreatureTargetFilter) Matches(targetID uuid.UUID, game GameContext) bool {
	// TODO: Check if target is a creature
	// TODO: Check subtype if specified
	return true
}

func (f *CreatureTargetFilter) GetDescription() string {
	if f.Subtype != "" {
		return fmt.Sprintf("target %s creature", f.Subtype)
	}
	return "target creature"
}

// ========================================
// Player Target
// ========================================

// PlayerTargetFilter matches only players
type PlayerTargetFilter struct {
	Opponent bool // If true, only opponents
}

func NewPlayerTargetFilter() *PlayerTargetFilter {
	return &PlayerTargetFilter{Opponent: false}
}

func NewOpponentTargetFilter() *PlayerTargetFilter {
	return &PlayerTargetFilter{Opponent: true}
}

func (f *PlayerTargetFilter) Matches(targetID uuid.UUID, game GameContext) bool {
	// TODO: Check if target is a player
	// TODO: Check if opponent if specified
	return true
}

func (f *PlayerTargetFilter) GetDescription() string {
	if f.Opponent {
		return "target opponent"
	}
	return "target player"
}

// ========================================
// Permanent Target
// ========================================

// PermanentTargetFilter matches any permanent
type PermanentTargetFilter struct {
	Type string // Optional type filter (e.g., "creature", "artifact")
}

func NewPermanentTargetFilter() *PermanentTargetFilter {
	return &PermanentTargetFilter{}
}

func NewPermanentTargetFilterWithType(cardType string) *PermanentTargetFilter {
	return &PermanentTargetFilter{Type: cardType}
}

func (f *PermanentTargetFilter) Matches(targetID uuid.UUID, game GameContext) bool {
	// TODO: Check if target is a permanent
	// TODO: Check type if specified
	return true
}

func (f *PermanentTargetFilter) GetDescription() string {
	if f.Type != "" {
		return fmt.Sprintf("target %s", f.Type)
	}
	return "target permanent"
}

// ========================================
// Spell Target
// ========================================

// SpellTargetFilter matches spells on the stack
type SpellTargetFilter struct{}

func NewSpellTargetFilter() *SpellTargetFilter {
	return &SpellTargetFilter{}
}

func (f *SpellTargetFilter) Matches(targetID uuid.UUID, game GameContext) bool {
	// TODO: Check if target is a spell on the stack
	return true
}

func (f *SpellTargetFilter) GetDescription() string {
	return "target spell"
}

// ========================================
// Planeswalker Target
// ========================================

// PlaneswalkerTargetFilter matches planeswalkers
type PlaneswalkerTargetFilter struct{}

func NewPlaneswalkerTargetFilter() *PlaneswalkerTargetFilter {
	return &PlaneswalkerTargetFilter{}
}

func (f *PlaneswalkerTargetFilter) Matches(targetID uuid.UUID, game GameContext) bool {
	// TODO: Check if target is a planeswalker
	return true
}

func (f *PlaneswalkerTargetFilter) GetDescription() string {
	return "target planeswalker"
}

// ========================================
// Artifact Target
// ========================================

// ArtifactTargetFilter matches artifacts
type ArtifactTargetFilter struct{}

func NewArtifactTargetFilter() *ArtifactTargetFilter {
	return &ArtifactTargetFilter{}
}

func (f *ArtifactTargetFilter) Matches(targetID uuid.UUID, game GameContext) bool {
	// TODO: Check if target is an artifact
	return true
}

func (f *ArtifactTargetFilter) GetDescription() string {
	return "target artifact"
}

// ========================================
// Enchantment Target
// ========================================

// EnchantmentTargetFilter matches enchantments
type EnchantmentTargetFilter struct{}

func NewEnchantmentTargetFilter() *EnchantmentTargetFilter {
	return &EnchantmentTargetFilter{}
}

func (f *EnchantmentTargetFilter) Matches(targetID uuid.UUID, game GameContext) bool {
	// TODO: Check if target is an enchantment
	return true
}

func (f *EnchantmentTargetFilter) GetDescription() string {
	return "target enchantment"
}

// ========================================
// Land Target
// ========================================

// LandTargetFilter matches lands
type LandTargetFilter struct{}

func NewLandTargetFilter() *LandTargetFilter {
	return &LandTargetFilter{}
}

func (f *LandTargetFilter) Matches(targetID uuid.UUID, game GameContext) bool {
	// TODO: Check if target is a land
	return true
}

func (f *LandTargetFilter) GetDescription() string {
	return "target land"
}

// ========================================
// Composite Filter (AND/OR)
// ========================================

// AndFilter matches if all sub-filters match
type AndFilter struct {
	Filters []TargetFilter
}

func NewAndFilter(filters ...TargetFilter) *AndFilter {
	return &AndFilter{Filters: filters}
}

func (f *AndFilter) Matches(targetID uuid.UUID, game GameContext) bool {
	for _, filter := range f.Filters {
		if !filter.Matches(targetID, game) {
			return false
		}
	}
	return true
}

func (f *AndFilter) GetDescription() string {
	if len(f.Filters) == 0 {
		return ""
	}
	if len(f.Filters) == 1 {
		return f.Filters[0].GetDescription()
	}

	desc := f.Filters[0].GetDescription()
	for i := 1; i < len(f.Filters); i++ {
		desc += " " + f.Filters[i].GetDescription()
	}
	return desc
}

// OrFilter matches if any sub-filter matches
type OrFilter struct {
	Filters []TargetFilter
}

func NewOrFilter(filters ...TargetFilter) *OrFilter {
	return &OrFilter{Filters: filters}
}

func (f *OrFilter) Matches(targetID uuid.UUID, game GameContext) bool {
	for _, filter := range f.Filters {
		if filter.Matches(targetID, game) {
			return true
		}
	}
	return false
}

func (f *OrFilter) GetDescription() string {
	if len(f.Filters) == 0 {
		return ""
	}
	if len(f.Filters) == 1 {
		return f.Filters[0].GetDescription()
	}

	desc := f.Filters[0].GetDescription()
	for i := 1; i < len(f.Filters)-1; i++ {
		desc += ", " + f.Filters[i].GetDescription()
	}
	desc += " or " + f.Filters[len(f.Filters)-1].GetDescription()
	return desc
}
