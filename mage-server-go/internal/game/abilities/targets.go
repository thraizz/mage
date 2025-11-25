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

	// Tag is used to distinguish between multiple target requirements
	// Java: Target.setTargetTag(int) / Target.getTargetTag()
	// Example: Resourceful Defense has tag=1 for source permanent, tag=2 for destination
	// Tag 0 means no tag (default)
	Tag int

	// ChosenTargets stores the targets actually chosen for this requirement
	ChosenTargets []uuid.UUID
}

// NewTargetRequirement creates a new target requirement
func NewTargetRequirement(min, max int, filter TargetFilter) *TargetRequirement {
	return &TargetRequirement{
		MinTargets:    min,
		MaxTargets:    max,
		Filter:        filter,
		Description:   filter.GetDescription(),
		Tag:           0,
		ChosenTargets: make([]uuid.UUID, 0),
	}
}

// NewTaggedTargetRequirement creates a target requirement with a specific tag
// Java: target.setTargetTag(tag)
// Tags are used when an ability has multiple distinct targets
func NewTaggedTargetRequirement(min, max int, filter TargetFilter, tag int) *TargetRequirement {
	return &TargetRequirement{
		MinTargets:    min,
		MaxTargets:    max,
		Filter:        filter,
		Description:   filter.GetDescription(),
		Tag:           tag,
		ChosenTargets: make([]uuid.UUID, 0),
	}
}

// SetTag sets the tag for this target requirement
func (tr *TargetRequirement) SetTag(tag int) *TargetRequirement {
	tr.Tag = tag
	return tr
}

// GetTag returns the tag for this target requirement
func (tr *TargetRequirement) GetTag() int {
	return tr.Tag
}

// AddChosenTarget adds a chosen target to this requirement
func (tr *TargetRequirement) AddChosenTarget(targetID uuid.UUID) {
	tr.ChosenTargets = append(tr.ChosenTargets, targetID)
}

// GetChosenTargets returns the targets chosen for this requirement
func (tr *TargetRequirement) GetChosenTargets() []uuid.UUID {
	return tr.ChosenTargets
}

// ClearChosenTargets removes all chosen targets
func (tr *TargetRequirement) ClearChosenTargets() {
	tr.ChosenTargets = make([]uuid.UUID, 0)
}

// GetFirstTarget returns the first chosen target, or uuid.Nil if none
func (tr *TargetRequirement) GetFirstTarget() uuid.UUID {
	if len(tr.ChosenTargets) > 0 {
		return tr.ChosenTargets[0]
	}
	return uuid.Nil
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

// NewPermanentTarget is a convenience function for creating a single permanent target
func NewPermanentTarget() TargetFilter {
	return NewPermanentTargetFilter()
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

// ========================================
// Card Filter (for cards in hand/graveyard)
// ========================================

// CardFilter determines what cards match (used for discard, search, etc.)
// This is similar to TargetFilter but for cards not on battlefield
type CardFilter interface {
	// Matches checks if a card matches this filter
	Matches(cardID uuid.UUID, game GameContext) bool

	// GetDescription returns a description of what this filter matches
	GetDescription() string
}

// ArtifactCardFilter matches artifact cards
type ArtifactCardFilter struct{}

// NewArtifactCardFilter creates a new artifact card filter
func NewArtifactCardFilter() *ArtifactCardFilter {
	return &ArtifactCardFilter{}
}

// Matches checks if the card is an artifact
func (f *ArtifactCardFilter) Matches(cardID uuid.UUID, game GameContext) bool {
	// TODO: Check if card is an artifact
	return true
}

// GetDescription returns the description
func (f *ArtifactCardFilter) GetDescription() string {
	return "artifact card"
}

// CreatureCardFilter matches creature cards
type CreatureCardFilter struct{}

// NewCreatureCardFilter creates a new creature card filter
func NewCreatureCardFilter() *CreatureCardFilter {
	return &CreatureCardFilter{}
}

// Matches checks if the card is a creature
func (f *CreatureCardFilter) Matches(cardID uuid.UUID, game GameContext) bool {
	// TODO: Check if card is a creature
	return true
}

// GetDescription returns the description
func (f *CreatureCardFilter) GetDescription() string {
	return "creature card"
}

// LandCardFilter matches land cards
type LandCardFilter struct{}

// NewLandCardFilter creates a new land card filter
func NewLandCardFilter() *LandCardFilter {
	return &LandCardFilter{}
}

// Matches checks if the card is a land
func (f *LandCardFilter) Matches(cardID uuid.UUID, game GameContext) bool {
	// TODO: Check if card is a land
	return true
}

// GetDescription returns the description
func (f *LandCardFilter) GetDescription() string {
	return "land card"
}

// AnyCardFilter matches any card
type AnyCardFilter struct{}

// NewAnyCardFilter creates a new any card filter
func NewAnyCardFilter() *AnyCardFilter {
	return &AnyCardFilter{}
}

// Matches always returns true
func (f *AnyCardFilter) Matches(cardID uuid.UUID, game GameContext) bool {
	return true
}

// GetDescription returns the description
func (f *AnyCardFilter) GetDescription() string {
	return "card"
}

// ========================================
// Has Counters Filter
// ========================================

// HasCountersFilter matches permanents that have counters on them
// Java: Filtering by counter presence (used for effects that care about counters)
type HasCountersFilter struct {
	counterType string       // Empty = any counter, otherwise specific type
	inner       TargetFilter // Optional inner filter to combine with
}

// NewHasCountersFilter creates a filter for permanents with any counters
func NewHasCountersFilter() *HasCountersFilter {
	return &HasCountersFilter{
		counterType: "",
		inner:       nil,
	}
}

// NewHasCounterTypeFilter creates a filter for permanents with specific counter type
func NewHasCounterTypeFilter(counterType string) *HasCountersFilter {
	return &HasCountersFilter{
		counterType: counterType,
		inner:       nil,
	}
}

// NewHasCountersFilterWithInner creates a filter combining counter check with inner filter
func NewHasCountersFilterWithInner(counterType string, inner TargetFilter) *HasCountersFilter {
	return &HasCountersFilter{
		counterType: counterType,
		inner:       inner,
	}
}

// Matches returns true if the target has the specified counters and passes inner filter
func (f *HasCountersFilter) Matches(targetID uuid.UUID, game GameContext) bool {
	// Apply inner filter first if present
	if f.inner != nil && !f.inner.Matches(targetID, game) {
		return false
	}

	// Get all counters on the permanent
	counters := game.GetAllCountersOnPermanent(nil, targetID)
	if counters == nil {
		return false
	}

	// Check for counters
	if f.counterType == "" {
		// Any counter will do
		for _, count := range counters {
			if count > 0 {
				return true
			}
		}
		return false
	}

	// Check for specific counter type
	count, exists := counters[f.counterType]
	return exists && count > 0
}

// GetDescription returns the filter description
func (f *HasCountersFilter) GetDescription() string {
	if f.counterType != "" {
		desc := fmt.Sprintf("permanent with %s counters", f.counterType)
		if f.inner != nil {
			desc = f.inner.GetDescription() + " with " + f.counterType + " counters"
		}
		return desc
	}
	if f.inner != nil {
		return f.inner.GetDescription() + " with counters on it"
	}
	return "permanent with counters on it"
}

// ========================================
// Another Permanent Filter
// ========================================

// AnotherPermanentFilter excludes the source permanent from valid targets
// Java: AnotherPredicate.instance used with FilterControlledPermanent
// Used for "another target permanent" targeting
type AnotherPermanentFilter struct {
	sourceID uuid.UUID    // The source permanent to exclude
	inner    TargetFilter // Optional inner filter to combine with
}

// NewAnotherPermanentFilter creates a filter that excludes the source permanent
func NewAnotherPermanentFilter(sourceID uuid.UUID) *AnotherPermanentFilter {
	return &AnotherPermanentFilter{
		sourceID: sourceID,
		inner:    nil,
	}
}

// NewAnotherPermanentFilterWithInner creates a filter that excludes source and applies inner filter
func NewAnotherPermanentFilterWithInner(sourceID uuid.UUID, inner TargetFilter) *AnotherPermanentFilter {
	return &AnotherPermanentFilter{
		sourceID: sourceID,
		inner:    inner,
	}
}

// Matches returns true if the target is not the source and passes inner filter
func (f *AnotherPermanentFilter) Matches(targetID uuid.UUID, game GameContext) bool {
	// Exclude the source permanent
	if targetID == f.sourceID {
		return false
	}

	// If there's an inner filter, apply it
	if f.inner != nil {
		return f.inner.Matches(targetID, game)
	}

	// Otherwise, any other permanent matches
	return true
}

// GetDescription returns the filter description
func (f *AnotherPermanentFilter) GetDescription() string {
	if f.inner != nil {
		return "another " + f.inner.GetDescription()
	}
	return "another permanent"
}

// ========================================
// Multi-Target Support
// ========================================

// MultiTargetRequirements manages multiple target requirements for abilities
// that need to select different types of targets (e.g., "source" and "destination")
// Java: Ability can have multiple Target objects with different tags
type MultiTargetRequirements struct {
	requirements []*TargetRequirement
}

// NewMultiTargetRequirements creates a new multi-target requirements container
func NewMultiTargetRequirements() *MultiTargetRequirements {
	return &MultiTargetRequirements{
		requirements: make([]*TargetRequirement, 0),
	}
}

// AddRequirement adds a target requirement
func (mtr *MultiTargetRequirements) AddRequirement(req *TargetRequirement) *MultiTargetRequirements {
	mtr.requirements = append(mtr.requirements, req)
	return mtr
}

// AddTaggedRequirement adds a target requirement with a specific tag
func (mtr *MultiTargetRequirements) AddTaggedRequirement(min, max int, filter TargetFilter, tag int) *MultiTargetRequirements {
	mtr.requirements = append(mtr.requirements, NewTaggedTargetRequirement(min, max, filter, tag))
	return mtr
}

// GetRequirementByTag returns the requirement with the specified tag, or nil
func (mtr *MultiTargetRequirements) GetRequirementByTag(tag int) *TargetRequirement {
	for _, req := range mtr.requirements {
		if req.Tag == tag {
			return req
		}
	}
	return nil
}

// GetTargetsByTag returns the chosen targets for the requirement with the specified tag
func (mtr *MultiTargetRequirements) GetTargetsByTag(tag int) []uuid.UUID {
	req := mtr.GetRequirementByTag(tag)
	if req != nil {
		return req.GetChosenTargets()
	}
	return nil
}

// GetFirstTargetByTag returns the first chosen target for the requirement with the specified tag
func (mtr *MultiTargetRequirements) GetFirstTargetByTag(tag int) uuid.UUID {
	req := mtr.GetRequirementByTag(tag)
	if req != nil {
		return req.GetFirstTarget()
	}
	return uuid.Nil
}

// GetAllRequirements returns all requirements
func (mtr *MultiTargetRequirements) GetAllRequirements() []*TargetRequirement {
	return mtr.requirements
}

// ClearAllChosenTargets clears chosen targets from all requirements
func (mtr *MultiTargetRequirements) ClearAllChosenTargets() {
	for _, req := range mtr.requirements {
		req.ClearChosenTargets()
	}
}
