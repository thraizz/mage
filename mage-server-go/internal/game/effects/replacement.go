package effects

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/rules"
)

// ReplacementEffect represents an effect that can replace or modify an event before it happens.
// Implements Rule 614 - Replacement Effects from the Comprehensive Rules.
//
// Key concepts:
// - Replacement effects apply continuously as events happen (not locked in ahead of time)
// - They watch for a particular event and completely or partially replace it
// - Self-replacement effects (from the resolving spell/ability) are applied first
// - Multiple replacement effects are chosen by the affected player/controller
// - A replacement effect gets only one opportunity per event (doesn't invoke itself repeatedly)
type ReplacementEffect interface {
	// ID returns the unique identifier for this effect
	ID() string

	// SourceID returns the ID of the source object/ability that created this effect
	SourceID() string

	// Duration returns how long this effect lasts
	Duration() Duration

	// ChecksEventType returns true if this effect cares about the given event type
	// This is the first filter - only effects that care about the event type are considered
	ChecksEventType(eventType rules.EventType) bool

	// Applies returns true if this effect applies to the given event
	// This is the second filter - checks specific conditions beyond just event type
	Applies(event rules.Event, gameID string) bool

	// ReplaceEvent modifies or replaces the event
	// Returns the modified event and a boolean indicating if the event was completely replaced
	// If true is returned, no further replacement effects apply to this event
	// If false is returned, the modified event can still be replaced by other effects
	ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool)

	// IsSelfReplacement returns true if this is a self-replacement effect (Rule 614.15)
	// Self-replacement effects are effects of a resolving spell/ability that replace
	// part of that spell or ability's own effect
	// Self-replacement effects are applied before other replacement effects
	IsSelfReplacement() bool

	// HasSelfScope returns true if this effect applies to events from its own source
	// Used primarily for "enters the battlefield" effects (Rule 614.12)
	HasSelfScope() bool
}

// BaseReplacementEffect provides common functionality for replacement effects
type BaseReplacementEffect struct {
	id               string
	sourceID         string
	duration         Duration
	selfReplacement  bool
	selfScope        bool
}

// NewBaseReplacementEffect creates a new base replacement effect
func NewBaseReplacementEffect(sourceID string, duration Duration, selfReplacement, selfScope bool) *BaseReplacementEffect {
	source := strings.TrimSpace(sourceID)
	seed := fmt.Sprintf("%s|replacement|%s|%t|%t|%d", source, duration, selfReplacement, selfScope, uuid.New().ID())
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(seed)).String()

	return &BaseReplacementEffect{
		id:              id,
		sourceID:        source,
		duration:        duration,
		selfReplacement: selfReplacement,
		selfScope:       selfScope,
	}
}

// ID returns the unique identifier
func (e *BaseReplacementEffect) ID() string {
	return e.id
}

// SourceID returns the source ID
func (e *BaseReplacementEffect) SourceID() string {
	return e.sourceID
}

// Duration returns the duration
func (e *BaseReplacementEffect) Duration() Duration {
	return e.duration
}

// IsSelfReplacement returns whether this is a self-replacement effect
func (e *BaseReplacementEffect) IsSelfReplacement() bool {
	return e.selfReplacement
}

// HasSelfScope returns whether this effect has self-scope
func (e *BaseReplacementEffect) HasSelfScope() bool {
	return e.selfScope
}

// PreventionEffect represents an effect that prevents damage or other events.
// Implements Rule 615 - Prevention Effects from the Comprehensive Rules.
//
// Key concepts:
// - Prevention effects are like replacement effects but specifically prevent damage
// - They apply continuously as events happen
// - They can prevent all or part of damage
// - Some have shields (e.g., "prevent the next 3 damage")
// - Unpreventable damage still triggers prevention effects (for additional effects) but doesn't reduce shields
type PreventionEffect interface {
	ReplacementEffect // Prevention effects are a specialized type of replacement effect

	// GetShield returns the remaining shield amount (0 if no shield/unlimited)
	GetShield() int

	// ReduceShield reduces the shield by the given amount
	// Returns the actual amount reduced
	ReduceShield(amount int) int
}

// BasePreventionEffect provides common functionality for prevention effects
type BasePreventionEffect struct {
	*BaseReplacementEffect
	shield int // Remaining shield amount (0 = unlimited or no shield)
}

// NewBasePreventionEffect creates a new base prevention effect
func NewBasePreventionEffect(sourceID string, duration Duration, shield int) *BasePreventionEffect {
	return &BasePreventionEffect{
		BaseReplacementEffect: NewBaseReplacementEffect(sourceID, duration, false, false),
		shield:                shield,
	}
}

// GetShield returns the remaining shield amount
func (e *BasePreventionEffect) GetShield() int {
	return e.shield
}

// ReduceShield reduces the shield by the given amount
// Returns the actual amount reduced (may be less than requested if shield is exhausted)
func (e *BasePreventionEffect) ReduceShield(amount int) int {
	// Shield of 0 means no shield system (unlimited prevention)
	// If we want to track actual prevention, we'd need different logic
	if e.shield <= 0 {
		// No shield system in place
		return 0
	}

	reduced := amount
	if reduced > e.shield {
		reduced = e.shield
	}

	e.shield -= reduced
	return reduced
}

// DamagePreventionEffect prevents damage from being dealt
// Example: "Prevent the next 3 damage that would be dealt to target creature"
type DamagePreventionEffect struct {
	*BasePreventionEffect
	targetID    string // Target that damage is prevented to (empty = any)
	sourceCheck string // Source that damage must come from (empty = any)
	amount      int    // Amount to prevent (0 = all)
}

// NewDamagePreventionEffect creates a damage prevention effect
func NewDamagePreventionEffect(sourceID, targetID, sourceCheck string, amount int, duration Duration) *DamagePreventionEffect {
	return &DamagePreventionEffect{
		BasePreventionEffect: NewBasePreventionEffect(sourceID, duration, amount),
		targetID:             strings.TrimSpace(targetID),
		sourceCheck:          strings.TrimSpace(sourceCheck),
		amount:               amount,
	}
}

// ChecksEventType checks if this effect cares about damage events
func (e *DamagePreventionEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventDamagePlayer ||
		eventType == rules.EventDamagePermanent ||
		eventType == rules.EventDamagedPlayer ||
		eventType == rules.EventDamagedPermanent
}

// Applies checks if this effect applies to the given damage event
func (e *DamagePreventionEffect) Applies(event rules.Event, gameID string) bool {
	// Check event type
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Check target if specified
	if e.targetID != "" && event.TargetID != e.targetID {
		return false
	}

	// Check source if specified
	if e.sourceCheck != "" && event.SourceID != e.sourceCheck {
		return false
	}

	// Check if we've already consumed our shield
	if e.GetShield() > 0 && e.amount > 0 {
		// Still has shield
		return true
	}

	// Unlimited prevention or no shield system
	return e.amount == 0
}

// ReplaceEvent prevents or reduces damage
func (e *DamagePreventionEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	if e.amount == 0 {
		// Prevent all damage
		event.Amount = 0
		return event, true // Completely replaced (no damage)
	}

	// Prevent up to shield amount
	prevented := e.ReduceShield(event.Amount)
	event.Amount -= prevented

	// If all damage prevented, event is completely replaced
	return event, event.Amount == 0
}

// ZoneChangeReplacementEffect replaces where a card goes during a zone change
// Example: "If a creature would die, exile it instead"
type ZoneChangeReplacementEffect struct {
	*BaseReplacementEffect
	fromZone      int    // Zone card is coming from (-1 = any)
	toZone        int    // Zone card is going to (-1 = any)
	newZone       int    // Zone to send card to instead
	cardID        string // Specific card (empty = any card matching conditions)
	controllerID  string // Controller of the card (empty = any)
	cardTypeCheck string // Card type requirement (empty = any)
}

// NewZoneChangeReplacementEffect creates a zone change replacement effect
func NewZoneChangeReplacementEffect(
	sourceID string,
	fromZone, toZone, newZone int,
	cardID, controllerID, cardTypeCheck string,
	duration Duration,
	selfScope bool,
) *ZoneChangeReplacementEffect {
	return &ZoneChangeReplacementEffect{
		BaseReplacementEffect: NewBaseReplacementEffect(sourceID, duration, false, selfScope),
		fromZone:              fromZone,
		toZone:                toZone,
		newZone:               newZone,
		cardID:                strings.TrimSpace(cardID),
		controllerID:          strings.TrimSpace(controllerID),
		cardTypeCheck:         strings.TrimSpace(cardTypeCheck),
	}
}

// ChecksEventType checks if this effect cares about zone change events
func (e *ZoneChangeReplacementEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventZoneChange ||
		eventType == rules.EventZoneChangeGroup ||
		eventType == rules.EventZoneChangeBatch
}

// Applies checks if this effect applies to the given zone change event
func (e *ZoneChangeReplacementEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Check specific card if specified
	if e.cardID != "" && event.TargetID != e.cardID {
		return false
	}

	// Check controller if specified
	if e.controllerID != "" && event.Controller != e.controllerID {
		return false
	}

	// For now, we don't check card type (would require game state access)
	// TODO: Add game state parameter to Applies() for type checking

	// Check zone (stored in event.Zone for destination, would need metadata for source)
	// For now, accept if any zone matches
	// TODO: Add proper zone tracking to events

	return true
}

// ReplaceEvent changes the destination zone
func (e *ZoneChangeReplacementEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Change the zone to the new zone
	event.Zone = e.newZone

	// Add to metadata for tracking
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["replacement_effect"] = e.ID()
	event.Metadata["original_zone"] = fmt.Sprintf("%d", e.toZone)
	event.Metadata["new_zone"] = fmt.Sprintf("%d", e.newZone)

	// Not completely replaced - the zone change still happens, just to a different zone
	return event, false
}

// DoubleAmountReplacementEffect doubles an amount in an event
// Example: "If you would gain life, you gain twice that much life instead"
type DoubleAmountReplacementEffect struct {
	*BaseReplacementEffect
	eventTypes   []rules.EventType // Event types this applies to
	targetID     string            // Specific target (empty = any)
	controllerID string            // Controller requirement (empty = any)
}

// NewDoubleAmountReplacementEffect creates a doubling replacement effect
func NewDoubleAmountReplacementEffect(
	sourceID string,
	eventTypes []rules.EventType,
	targetID, controllerID string,
	duration Duration,
) *DoubleAmountReplacementEffect {
	return &DoubleAmountReplacementEffect{
		BaseReplacementEffect: NewBaseReplacementEffect(sourceID, duration, false, false),
		eventTypes:            eventTypes,
		targetID:              strings.TrimSpace(targetID),
		controllerID:          strings.TrimSpace(controllerID),
	}
}

// ChecksEventType checks if this effect cares about the given event type
func (e *DoubleAmountReplacementEffect) ChecksEventType(eventType rules.EventType) bool {
	for _, et := range e.eventTypes {
		if et == eventType {
			return true
		}
	}
	return false
}

// Applies checks if this effect applies to the given event
func (e *DoubleAmountReplacementEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Check target if specified
	if e.targetID != "" && event.TargetID != e.targetID {
		return false
	}

	// Check controller if specified
	if e.controllerID != "" && event.Controller != e.controllerID {
		return false
	}

	return true
}

// ReplaceEvent doubles the amount
func (e *DoubleAmountReplacementEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	event.Amount = event.Amount * 2

	// Add to metadata for tracking
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["doubled_by"] = e.ID()

	// Not completely replaced - event still happens with modified amount
	return event, false
}
