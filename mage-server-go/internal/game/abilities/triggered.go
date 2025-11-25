package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ========================================
// Triggered Ability System
// ========================================

// TriggeredAbility represents a "when X happens, do Y" ability
// Java: mage.abilities.TriggeredAbilityImpl
// MTG Rules: 603 (Triggered Abilities)
type TriggeredAbility struct {
	baseAbility
	trigger  TriggerCondition
	effects  []Effect
	targets  *TargetRequirement
	optional bool // "you may" triggers

	// values stores runtime data passed from trigger condition to effects
	// Java: Uses setValue/getValue on Effects to pass data like counters
	// Example: ResourcefulDefense trigger stores counters, effect retrieves them
	values map[string]interface{}
}

// NewTriggeredAbility creates a new triggered ability
func NewTriggeredAbility(sourceID uuid.UUID, trigger TriggerCondition, effects []Effect, optional bool) *TriggeredAbility {
	text := trigger.GetDescription()
	if optional {
		text = "you may " + text
	}

	return &TriggeredAbility{
		baseAbility: newBaseAbility(sourceID, text),
		trigger:     trigger,
		effects:     effects,
		optional:    optional,
		values:      make(map[string]interface{}),
	}
}

// GetType returns the ability type
func (a *TriggeredAbility) GetType() AbilityType {
	return AbilityTypeTriggered
}

// CanActivate checks if this trigger matches the current event
func (a *TriggeredAbility) CanActivate(ctx context.Context, game GameContext) bool {
	// Triggered abilities are automatically put on stack when they trigger
	// This method would be called by the game engine's event system
	return true
}

// Resolve resolves this triggered ability
func (a *TriggeredAbility) Resolve(ctx context.Context, game GameContext) error {
	// If optional and player declines, do nothing
	if a.optional {
		// TODO: Ask player if they want to trigger
		// For now, assume yes
	}

	// Resolve each effect in order
	for _, effect := range a.effects {
		targets := []uuid.UUID{} // TODO: Get actual targets
		if err := effect.Apply(ctx, game, a.sourceID, targets); err != nil {
			return fmt.Errorf("failed to resolve effect: %w", err)
		}
	}

	return nil
}

// SetTargets sets the target requirement
func (a *TriggeredAbility) SetTargets(targets *TargetRequirement) {
	a.targets = targets
}

// GetTargets returns the target requirement for this ability
func (a *TriggeredAbility) GetTargets() *TargetRequirement {
	return a.targets
}

// CheckTrigger checks if this ability should trigger for the given event
func (a *TriggeredAbility) CheckTrigger(event GameEvent) bool {
	return a.trigger.Check(event)
}

// IsOptional returns whether this is an optional trigger
func (a *TriggeredAbility) IsOptional() bool {
	return a.optional
}

// SetValue stores a value that can be retrieved later by effects
// Java: this.getEffects().setValue(key, value)
// This is used to pass data from the trigger condition to the effect
// Example: A trigger saves the counters from a leaving permanent, then the effect retrieves them
func (a *TriggeredAbility) SetValue(key string, value interface{}) {
	if a.values == nil {
		a.values = make(map[string]interface{})
	}
	a.values[key] = value
}

// GetValue retrieves a value stored by SetValue
// Java: this.getValue(key)
// Returns nil if the key doesn't exist
func (a *TriggeredAbility) GetValue(key string) interface{} {
	if a.values == nil {
		return nil
	}
	return a.values[key]
}

// GetValueAsCounters retrieves a value as a counter map (convenience method)
// Returns nil if the key doesn't exist or isn't a map[string]int
func (a *TriggeredAbility) GetValueAsCounters(key string) map[string]int {
	value := a.GetValue(key)
	if value == nil {
		return nil
	}
	if counters, ok := value.(map[string]int); ok {
		return counters
	}
	return nil
}

// ClearValues removes all stored values (called after resolution)
func (a *TriggeredAbility) ClearValues() {
	a.values = make(map[string]interface{})
}

// ========================================
// Trigger Conditions
// ========================================

// TriggerCondition represents when a triggered ability triggers
type TriggerCondition interface {
	// Check returns true if this condition matches the event
	Check(event GameEvent) bool

	// GetDescription returns a text description of the trigger
	GetDescription() string
}

// GameEvent represents an event in the game that can trigger abilities
type GameEvent struct {
	Type     EventType
	SourceID uuid.UUID // The object that caused the event
	TargetID uuid.UUID // The object affected by the event
	PlayerID uuid.UUID // The player involved (if applicable)
	Amount   int       // Numeric value (damage, counters, etc.)
	Zone     Zone      // Zone information
	FromZone Zone      // Previous zone (for zone change events)
	ToZone   Zone      // New zone (for zone change events)

	// PermanentSnapshot contains the state of a permanent at the moment it left a zone
	// Java: ZoneChangeEvent.getTarget() returns the Permanent object
	// This is populated for EventLeavesBattlefield and EventDies events
	PermanentSnapshot *PermanentSnapshot
}

// PermanentSnapshot captures the state of a permanent at a specific moment
// Used for triggered abilities that need to know the state of a permanent
// after it has left the battlefield (e.g., Resourceful Defense needs counters)
// Java: Similar to how ZoneChangeEvent stores the Permanent reference
type PermanentSnapshot struct {
	ID           uuid.UUID      // Permanent's ID
	Name         string         // Card name
	ControllerID uuid.UUID      // Controller at time of snapshot
	OwnerID      uuid.UUID      // Owner
	Types        []string       // Card types
	SubTypes     []string       // Subtypes
	Power        int            // Power (for creatures)
	Toughness    int            // Toughness (for creatures)
	Counters     map[string]int // All counters (name -> count)
	Tapped       bool           // Was it tapped
	Abilities    []string       // Ability descriptions
}

// NewPermanentSnapshot creates a PermanentSnapshot from basic data
func NewPermanentSnapshot(id, controllerID, ownerID uuid.UUID, name string) *PermanentSnapshot {
	return &PermanentSnapshot{
		ID:           id,
		Name:         name,
		ControllerID: controllerID,
		OwnerID:      ownerID,
		Counters:     make(map[string]int),
		Types:        make([]string, 0),
		SubTypes:     make([]string, 0),
		Abilities:    make([]string, 0),
	}
}

// HasCounters returns true if the snapshot has any counters
func (ps *PermanentSnapshot) HasCounters() bool {
	if ps == nil || ps.Counters == nil {
		return false
	}
	for _, count := range ps.Counters {
		if count > 0 {
			return true
		}
	}
	return false
}

// GetCounterCount returns the count of a specific counter type
func (ps *PermanentSnapshot) GetCounterCount(counterName string) int {
	if ps == nil || ps.Counters == nil {
		return 0
	}
	return ps.Counters[counterName]
}

// GetAllCounters returns a copy of all counters
func (ps *PermanentSnapshot) GetAllCounters() map[string]int {
	if ps == nil || ps.Counters == nil {
		return make(map[string]int)
	}
	result := make(map[string]int)
	for name, count := range ps.Counters {
		result[name] = count
	}
	return result
}

// EventType represents the type of game event
type EventType int

const (
	EventEntersBattlefield EventType = iota
	EventLeavesBattlefield
	EventDies
	EventTapped
	EventUntapped
	EventDamageDealt
	EventDamageTaken
	EventCounterAdded
	EventCounterRemoved
	EventAttackerDeclared
	EventBlockerDeclared
	EventSpellCast
	EventAbilityActivated
	EventCardDrawn
	EventCardDiscarded
	EventLifeGained
	EventLifeLost
	EventTokenCreated
	EventPhaseBegin
	EventPhaseEnd
	EventStepBegin
	EventStepEnd
	// Add more as needed
)

// ========================================
// Specific Trigger Types
// ========================================

// EntersBattlefieldTrigger triggers when a permanent enters the battlefield
// Java: mage.abilities.common.EntersBattlefieldTriggeredAbility
// MTG Rules: 603.6a
type EntersBattlefieldTrigger struct {
	sourceID uuid.UUID
}

// NewEntersBattlefieldTrigger creates a new ETB trigger
func NewEntersBattlefieldTrigger(sourceID uuid.UUID) *EntersBattlefieldTrigger {
	return &EntersBattlefieldTrigger{sourceID: sourceID}
}

// Check returns true if this is an ETB event for this permanent
func (t *EntersBattlefieldTrigger) Check(event GameEvent) bool {
	return event.Type == EventEntersBattlefield && event.SourceID == t.sourceID
}

// GetDescription returns the trigger description
func (t *EntersBattlefieldTrigger) GetDescription() string {
	return "when this permanent enters the battlefield"
}

// LeavesBattlefieldTrigger triggers when a permanent leaves the battlefield
// Java: mage.abilities.common.LeavesBattlefieldTriggeredAbility
type LeavesBattlefieldTrigger struct {
	sourceID uuid.UUID
}

// NewLeavesBattlefieldTrigger creates a new LTB trigger
func NewLeavesBattlefieldTrigger(sourceID uuid.UUID) *LeavesBattlefieldTrigger {
	return &LeavesBattlefieldTrigger{sourceID: sourceID}
}

// Check returns true if this is an LTB event for this permanent
func (t *LeavesBattlefieldTrigger) Check(event GameEvent) bool {
	return event.Type == EventLeavesBattlefield && event.SourceID == t.sourceID
}

// GetDescription returns the trigger description
func (t *LeavesBattlefieldTrigger) GetDescription() string {
	return "when this permanent leaves the battlefield"
}

// DiesTrigger triggers when a creature dies (goes to graveyard from battlefield)
// Java: mage.abilities.common.DiesTriggeredAbility
// MTG Rules: 700.4 (Dies = creature goes from battlefield to graveyard)
type DiesTrigger struct {
	sourceID uuid.UUID
}

// NewDiesTrigger creates a new dies trigger
func NewDiesTrigger(sourceID uuid.UUID) *DiesTrigger {
	return &DiesTrigger{sourceID: sourceID}
}

// Check returns true if this creature died
func (t *DiesTrigger) Check(event GameEvent) bool {
	return event.Type == EventDies && event.SourceID == t.sourceID
}

// GetDescription returns the trigger description
func (t *DiesTrigger) GetDescription() string {
	return "when this creature dies"
}

// BecomesTappedTrigger triggers when a permanent becomes tapped
type BecomesTappedTrigger struct {
	sourceID uuid.UUID
}

// NewBecomesTappedTrigger creates a new becomes tapped trigger
func NewBecomesTappedTrigger(sourceID uuid.UUID) *BecomesTappedTrigger {
	return &BecomesTappedTrigger{sourceID: sourceID}
}

// Check returns true if this permanent became tapped
func (t *BecomesTappedTrigger) Check(event GameEvent) bool {
	return event.Type == EventTapped && event.SourceID == t.sourceID
}

// GetDescription returns the trigger description
func (t *BecomesTappedTrigger) GetDescription() string {
	return "when this permanent becomes tapped"
}

// BecomesUntappedTrigger triggers when a permanent becomes untapped
type BecomesUntappedTrigger struct {
	sourceID uuid.UUID
}

// NewBecomesUntappedTrigger creates a new becomes untapped trigger
func NewBecomesUntappedTrigger(sourceID uuid.UUID) *BecomesUntappedTrigger {
	return &BecomesUntappedTrigger{sourceID: sourceID}
}

// Check returns true if this permanent became untapped
func (t *BecomesUntappedTrigger) Check(event GameEvent) bool {
	return event.Type == EventUntapped && event.SourceID == t.sourceID
}

// GetDescription returns the trigger description
func (t *BecomesUntappedTrigger) GetDescription() string {
	return "when this permanent becomes untapped"
}

// DealsDamageTrigger triggers when this source deals damage
type DealsDamageTrigger struct {
	sourceID uuid.UUID
}

// NewDealsDamageTrigger creates a new deals damage trigger
func NewDealsDamageTrigger(sourceID uuid.UUID) *DealsDamageTrigger {
	return &DealsDamageTrigger{sourceID: sourceID}
}

// Check returns true if this source dealt damage
func (t *DealsDamageTrigger) Check(event GameEvent) bool {
	return event.Type == EventDamageDealt && event.SourceID == t.sourceID
}

// GetDescription returns the trigger description
func (t *DealsDamageTrigger) GetDescription() string {
	return "when this creature deals damage"
}

// AttacksTrigger triggers when this creature attacks
// Java: mage.abilities.common.AttacksTriggeredAbility
type AttacksTrigger struct {
	sourceID uuid.UUID
}

// NewAttacksTrigger creates a new attacks trigger
func NewAttacksTrigger(sourceID uuid.UUID) *AttacksTrigger {
	return &AttacksTrigger{sourceID: sourceID}
}

// Check returns true if this creature attacked
func (t *AttacksTrigger) Check(event GameEvent) bool {
	return event.Type == EventAttackerDeclared && event.SourceID == t.sourceID
}

// GetDescription returns the trigger description
func (t *AttacksTrigger) GetDescription() string {
	return "when this creature attacks"
}

// BlocksTrigger triggers when this creature blocks
// Java: mage.abilities.common.BlocksTriggeredAbility
type BlocksTrigger struct {
	sourceID uuid.UUID
}

// NewBlocksTrigger creates a new blocks trigger
func NewBlocksTrigger(sourceID uuid.UUID) *BlocksTrigger {
	return &BlocksTrigger{sourceID: sourceID}
}

// Check returns true if this creature blocked
func (t *BlocksTrigger) Check(event GameEvent) bool {
	return event.Type == EventBlockerDeclared && event.SourceID == t.sourceID
}

// GetDescription returns the trigger description
func (t *BlocksTrigger) GetDescription() string {
	return "when this creature blocks"
}

// BecomesBlockedTrigger triggers when this creature becomes blocked
// Java: mage.abilities.common.BecomesBlockedTriggeredAbility
type BecomesBlockedTrigger struct {
	sourceID uuid.UUID
}

// NewBecomesBlockedTrigger creates a new becomes blocked trigger
func NewBecomesBlockedTrigger(sourceID uuid.UUID) *BecomesBlockedTrigger {
	return &BecomesBlockedTrigger{sourceID: sourceID}
}

// Check returns true if this creature became blocked
func (t *BecomesBlockedTrigger) Check(event GameEvent) bool {
	return event.Type == EventBlockerDeclared && event.TargetID == t.sourceID
}

// GetDescription returns the trigger description
func (t *BecomesBlockedTrigger) GetDescription() string {
	return "when this creature becomes blocked"
}

// DealsCombatDamageTrigger triggers when this creature deals combat damage
// Java: mage.abilities.common.DealsCombatDamageTriggeredAbility
type DealsCombatDamageTrigger struct {
	sourceID uuid.UUID
}

// NewDealsCombatDamageTrigger creates a new deals combat damage trigger
func NewDealsCombatDamageTrigger(sourceID uuid.UUID) *DealsCombatDamageTrigger {
	return &DealsCombatDamageTrigger{sourceID: sourceID}
}

// Check returns true if this creature dealt combat damage
func (t *DealsCombatDamageTrigger) Check(event GameEvent) bool {
	// Check if this is combat damage from this source
	return event.Type == EventDamageDealt && event.SourceID == t.sourceID
	// TODO: Add combat damage flag to GameEvent to distinguish from non-combat damage
}

// GetDescription returns the trigger description
func (t *DealsCombatDamageTrigger) GetDescription() string {
	return "when this creature deals combat damage"
}

// ========================================
// Triggered Ability Builder
// ========================================

// TriggeredAbilityBuilder provides a fluent API for building triggered abilities
type TriggeredAbilityBuilder struct {
	sourceID uuid.UUID
	trigger  TriggerCondition
	effects  []Effect
	targets  *TargetRequirement
	optional bool
}

// NewTriggeredAbilityBuilder creates a new triggered ability builder
func NewTriggeredAbilityBuilder(sourceID uuid.UUID) *TriggeredAbilityBuilder {
	return &TriggeredAbilityBuilder{
		sourceID: sourceID,
		effects:  make([]Effect, 0),
		optional: false,
	}
}

// SetTrigger sets the trigger condition
func (b *TriggeredAbilityBuilder) SetTrigger(trigger TriggerCondition) *TriggeredAbilityBuilder {
	b.trigger = trigger
	return b
}

// AddEffect adds an effect to this triggered ability
func (b *TriggeredAbilityBuilder) AddEffect(effect Effect) *TriggeredAbilityBuilder {
	b.effects = append(b.effects, effect)
	return b
}

// AddTarget sets the target requirement
func (b *TriggeredAbilityBuilder) AddTarget(filter TargetFilter) *TriggeredAbilityBuilder {
	b.targets = NewTargetRequirement(1, 1, filter)
	return b
}

// AddTargets sets the target requirement with custom min/max
func (b *TriggeredAbilityBuilder) AddTargets(min, max int, filter TargetFilter) *TriggeredAbilityBuilder {
	b.targets = NewTargetRequirement(min, max, filter)
	return b
}

// SetOptional makes this trigger optional ("you may")
func (b *TriggeredAbilityBuilder) SetOptional(optional bool) *TriggeredAbilityBuilder {
	b.optional = optional
	return b
}

// Build constructs the triggered ability
func (b *TriggeredAbilityBuilder) Build() *TriggeredAbility {
	ability := NewTriggeredAbility(b.sourceID, b.trigger, b.effects, b.optional)

	if b.targets != nil {
		ability.SetTargets(b.targets)
	}

	return ability
}
