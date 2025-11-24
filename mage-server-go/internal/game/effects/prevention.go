package effects

import (
	"fmt"
	"strings"

	"github.com/magefree/mage-server-go/internal/game/rules"
)

// Rule 615: Prevention Effects
// Prevention effects are like replacement effects that specifically prevent damage.
// They watch for damage events and prevent all or part of the damage.

// TargetedPreventionEffect prevents damage to/from a specific target
// Example: "Prevent all damage that would be dealt to target creature this turn"
type TargetedPreventionEffect struct {
	*BasePreventionEffect
	targetID    string // Target that damage is prevented to/from
	preventTo   bool   // Prevent damage TO this target
	preventFrom bool   // Prevent damage FROM this target
	amount      int    // Amount to prevent (0 = all)
	description string
}

// NewTargetedPreventionEffect creates a targeted prevention effect
func NewTargetedPreventionEffect(
	sourceID, targetID string,
	preventTo, preventFrom bool,
	amount int,
	duration Duration,
	description string,
) *TargetedPreventionEffect {
	return &TargetedPreventionEffect{
		BasePreventionEffect: NewBasePreventionEffect(sourceID, duration, amount),
		targetID:             strings.TrimSpace(targetID),
		preventTo:            preventTo,
		preventFrom:          preventFrom,
		amount:               amount,
		description:          description,
	}
}

func (e *TargetedPreventionEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventDamagePlayer ||
		eventType == rules.EventDamagePermanent ||
		eventType == rules.EventDamagedPlayer ||
		eventType == rules.EventDamagedPermanent
}

func (e *TargetedPreventionEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Check if preventing damage TO target
	if e.preventTo && event.TargetID == e.targetID {
		return true
	}

	// Check if preventing damage FROM target
	if e.preventFrom && event.SourceID == e.targetID {
		return true
	}

	return false
}

func (e *TargetedPreventionEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	originalAmount := event.Amount

	if e.amount == 0 {
		// Prevent all damage
		event.Amount = 0
	} else {
		// Prevent up to shield amount
		prevented := e.ReduceShield(event.Amount)
		event.Amount -= prevented
	}

	// Add metadata
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["prevented_by"] = e.ID()
	event.Metadata["prevented_amount"] = fmt.Sprintf("%d", originalAmount-event.Amount)
	event.Metadata["description"] = e.description

	return event, event.Amount == 0
}

// ProtectionPreventionEffect prevents damage from sources with specific qualities
// Example: Protection from red, Protection from creatures
type ProtectionPreventionEffect struct {
	*BasePreventionEffect
	protectedID string                                      // Permanent with protection
	checkSource func(event rules.Event, gameID string) bool // Check if source matches protection
	description string
}

// NewProtectionPreventionEffect creates a protection-based prevention effect
func NewProtectionPreventionEffect(
	sourceID, protectedID string,
	checkSource func(event rules.Event, gameID string) bool,
	description string,
) *ProtectionPreventionEffect {
	return &ProtectionPreventionEffect{
		BasePreventionEffect: NewBasePreventionEffect(sourceID, DurationPermanent, 0),
		protectedID:          strings.TrimSpace(protectedID),
		checkSource:          checkSource,
		description:          description,
	}
}

func (e *ProtectionPreventionEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventDamagePlayer ||
		eventType == rules.EventDamagePermanent ||
		eventType == rules.EventDamagedPlayer ||
		eventType == rules.EventDamagedPermanent
}

func (e *ProtectionPreventionEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Must be damage to the protected permanent
	if event.TargetID != e.protectedID {
		return false
	}

	// Check if source matches protection criteria
	if e.checkSource != nil {
		return e.checkSource(event, gameID)
	}

	return true
}

func (e *ProtectionPreventionEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Protection prevents ALL damage from matching sources
	event.Amount = 0

	// Add metadata
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["prevented_by_protection"] = e.ID()
	event.Metadata["description"] = e.description

	return event, true // Completely prevented
}

// IndestructiblePreventionEffect prevents lethal damage/destruction
// Example: Indestructible prevents destruction and lethal damage
type IndestructiblePreventionEffect struct {
	*BaseReplacementEffect
	permanentID string
	description string
}

// NewIndestructiblePreventionEffect creates an indestructible effect
func NewIndestructiblePreventionEffect(sourceID, permanentID string) *IndestructiblePreventionEffect {
	return &IndestructiblePreventionEffect{
		BaseReplacementEffect: NewBaseReplacementEffect(sourceID, DurationPermanent, false, false),
		permanentID:           strings.TrimSpace(permanentID),
		description:           "Indestructible",
	}
}

func (e *IndestructiblePreventionEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventDestroy ||
		eventType == rules.EventZoneChange // Prevent going to graveyard from lethal damage
}

func (e *IndestructiblePreventionEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Must be our permanent
	if event.TargetID != e.permanentID {
		return false
	}

	// For zone changes, only prevent if going to graveyard from lethal damage/destruction
	if event.Type == rules.EventZoneChange {
		if event.Zone != ZoneGraveyard {
			return false
		}
		// Check if it's from destruction/lethal damage (would need metadata)
		if event.Metadata != nil {
			reason := event.Metadata["reason"]
			return reason == "destroyed" || reason == "lethal_damage"
		}
	}

	return true
}

func (e *IndestructiblePreventionEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Prevent destruction
	if event.Type == rules.EventDestroy {
		// Add metadata
		if event.Metadata == nil {
			event.Metadata = make(map[string]string)
		}
		event.Metadata["indestructible"] = "true"
		event.Metadata["prevented_destruction"] = e.ID()

		return event, true // Completely prevent destruction
	}

	// Prevent zone change to graveyard from lethal damage
	if event.Type == rules.EventZoneChange && event.Zone == ZoneGraveyard {
		// Keep on battlefield
		event.Zone = ZoneBattlefield

		if event.Metadata == nil {
			event.Metadata = make(map[string]string)
		}
		event.Metadata["indestructible"] = "true"
		event.Metadata["stayed_on_battlefield"] = e.ID()

		return event, false // Modified but not completely replaced
	}

	return event, false
}

// ShieldPreventionEffect provides a damage prevention shield
// Example: "Prevent the next 3 damage that would be dealt to any target"
type ShieldPreventionEffect struct {
	*BasePreventionEffect
	targetIDs   []string // Multiple targets that can be protected
	anyTarget   bool     // If true, protects any target
	description string
}

// NewShieldPreventionEffect creates a shield effect
func NewShieldPreventionEffect(
	sourceID string,
	shieldAmount int,
	targetIDs []string,
	anyTarget bool,
	duration Duration,
	description string,
) *ShieldPreventionEffect {
	return &ShieldPreventionEffect{
		BasePreventionEffect: NewBasePreventionEffect(sourceID, duration, shieldAmount),
		targetIDs:            targetIDs,
		anyTarget:            anyTarget,
		description:          description,
	}
}

func (e *ShieldPreventionEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventDamagePlayer ||
		eventType == rules.EventDamagePermanent ||
		eventType == rules.EventDamagedPlayer ||
		eventType == rules.EventDamagedPermanent
}

func (e *ShieldPreventionEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Check if shield is exhausted
	if e.GetShield() <= 0 {
		return false
	}

	// Check if any target is allowed
	if e.anyTarget {
		return true
	}

	// Check if damage is to one of our protected targets
	for _, targetID := range e.targetIDs {
		if event.TargetID == targetID {
			return true
		}
	}

	return false
}

func (e *ShieldPreventionEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Prevent damage up to shield amount
	prevented := e.ReduceShield(event.Amount)
	event.Amount -= prevented

	// Add metadata
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["prevented_by_shield"] = e.ID()
	event.Metadata["prevented_amount"] = fmt.Sprintf("%d", prevented)
	event.Metadata["remaining_shield"] = fmt.Sprintf("%d", e.GetShield())
	event.Metadata["description"] = e.description

	return event, event.Amount == 0
}

// AbsorbPreventionEffect absorbs damage and converts it to something else
// Example: "If damage would be dealt to ~, prevent that damage and put a +1/+1 counter on it instead"
type AbsorbPreventionEffect struct {
	*BasePreventionEffect
	targetID    string                                 // Target that absorbs damage
	onAbsorb    func(event rules.Event, gameID string) // What happens when damage is absorbed
	description string
}

// NewAbsorbPreventionEffect creates an absorb effect
func NewAbsorbPreventionEffect(
	sourceID, targetID string,
	onAbsorb func(event rules.Event, gameID string),
	duration Duration,
	description string,
) *AbsorbPreventionEffect {
	return &AbsorbPreventionEffect{
		BasePreventionEffect: NewBasePreventionEffect(sourceID, duration, 0),
		targetID:             strings.TrimSpace(targetID),
		onAbsorb:             onAbsorb,
		description:          description,
	}
}

func (e *AbsorbPreventionEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventDamagePlayer ||
		eventType == rules.EventDamagePermanent ||
		eventType == rules.EventDamagedPlayer ||
		eventType == rules.EventDamagedPermanent
}

func (e *AbsorbPreventionEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Must be damage to our target
	return event.TargetID == e.targetID
}

func (e *AbsorbPreventionEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	originalAmount := event.Amount

	// Prevent all damage
	event.Amount = 0

	// Trigger absorb effect
	if e.onAbsorb != nil {
		e.onAbsorb(event, gameID)
	}

	// Add metadata
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["absorbed_by"] = e.ID()
	event.Metadata["absorbed_amount"] = fmt.Sprintf("%d", originalAmount)
	event.Metadata["description"] = e.description

	return event, true // Completely prevented
}

// RedirectPreventionEffect prevents damage and redirects it elsewhere
// Example: "Prevent all damage that would be dealt to you. Redirect that damage to target creature"
type RedirectPreventionEffect struct {
	*BasePreventionEffect
	originalTargetID string                                      // Original damage target
	newTargetID      string                                      // Where to redirect damage
	condition        func(event rules.Event, gameID string) bool // Optional condition
	description      string
}

// NewRedirectPreventionEffect creates a redirect effect
func NewRedirectPreventionEffect(
	sourceID, originalTargetID, newTargetID string,
	condition func(event rules.Event, gameID string) bool,
	duration Duration,
	description string,
) *RedirectPreventionEffect {
	return &RedirectPreventionEffect{
		BasePreventionEffect: NewBasePreventionEffect(sourceID, duration, 0),
		originalTargetID:     strings.TrimSpace(originalTargetID),
		newTargetID:          strings.TrimSpace(newTargetID),
		condition:            condition,
		description:          description,
	}
}

func (e *RedirectPreventionEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventDamagePlayer ||
		eventType == rules.EventDamagePermanent ||
		eventType == rules.EventDamagedPlayer ||
		eventType == rules.EventDamagedPermanent
}

func (e *RedirectPreventionEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Must be damage to original target
	if event.TargetID != e.originalTargetID {
		return false
	}

	// Check condition if provided
	if e.condition != nil && !e.condition(event, gameID) {
		return false
	}

	return true
}

func (e *RedirectPreventionEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Redirect the damage
	oldTarget := event.TargetID
	event.TargetID = e.newTargetID

	// Add metadata
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["redirected_by"] = e.ID()
	event.Metadata["original_target"] = oldTarget
	event.Metadata["new_target"] = e.newTargetID
	event.Metadata["description"] = e.description

	// Not completely prevented - damage still happens, just to different target
	return event, false
}

// UnpreventableDamageModifier marks damage as unpreventable
// Example: "Damage can't be prevented this turn", Skullcrack, Leyline of Punishment
type UnpreventableDamageModifier struct {
	*BaseReplacementEffect
	sourceID    string                                      // Specific damage source (empty = all)
	targetID    string                                      // Specific damage target (empty = all)
	condition   func(event rules.Event, gameID string) bool // Optional condition
	description string
}

// NewUnpreventableDamageModifier creates an unpreventable damage modifier
func NewUnpreventableDamageModifier(
	effectSourceID, damageSourceID, targetID string,
	condition func(event rules.Event, gameID string) bool,
	duration Duration,
	description string,
) *UnpreventableDamageModifier {
	return &UnpreventableDamageModifier{
		BaseReplacementEffect: NewBaseReplacementEffect(effectSourceID, duration, false, false),
		sourceID:              strings.TrimSpace(damageSourceID),
		targetID:              strings.TrimSpace(targetID),
		condition:             condition,
		description:           description,
	}
}

func (e *UnpreventableDamageModifier) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventDamagePlayer ||
		eventType == rules.EventDamagePermanent ||
		eventType == rules.EventDamagedPlayer ||
		eventType == rules.EventDamagedPermanent
}

func (e *UnpreventableDamageModifier) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Check source if specified
	if e.sourceID != "" && event.SourceID != e.sourceID {
		return false
	}

	// Check target if specified
	if e.targetID != "" && event.TargetID != e.targetID {
		return false
	}

	// Check condition if provided
	if e.condition != nil && !e.condition(event, gameID) {
		return false
	}

	return true
}

func (e *UnpreventableDamageModifier) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Mark damage as unpreventable
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["unpreventable"] = "true"
	event.Metadata["unpreventable_source"] = e.ID()
	event.Metadata["description"] = e.description

	// Not completely replaced - damage still happens with unpreventable flag
	return event, false
}

// DamageDoublePreventionAdjuster modifies damage but preserves prevention logic
// Example: "If a source would deal damage, it deals double that damage instead"
// This applies BEFORE prevention effects (Rule 615.10)
type DamageDoublePreventionAdjuster struct {
	*BaseReplacementEffect
	multiplier  float64                                     // Damage multiplier (2.0 = double)
	sourceCheck func(event rules.Event, gameID string) bool // Check if applies to source
	description string
}

// NewDamageDoublePreventionAdjuster creates a damage adjustment effect
func NewDamageDoublePreventionAdjuster(
	effectSourceID string,
	multiplier float64,
	sourceCheck func(event rules.Event, gameID string) bool,
	duration Duration,
	description string,
) *DamageDoublePreventionAdjuster {
	return &DamageDoublePreventionAdjuster{
		BaseReplacementEffect: NewBaseReplacementEffect(effectSourceID, duration, false, false),
		multiplier:            multiplier,
		sourceCheck:           sourceCheck,
		description:           description,
	}
}

func (e *DamageDoublePreventionAdjuster) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventDamagePlayer ||
		eventType == rules.EventDamagePermanent ||
		eventType == rules.EventDamagedPlayer ||
		eventType == rules.EventDamagedPermanent
}

func (e *DamageDoublePreventionAdjuster) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Check if applies to this source
	if e.sourceCheck != nil {
		return e.sourceCheck(event, gameID)
	}

	return true
}

func (e *DamageDoublePreventionAdjuster) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Apply multiplier
	originalAmount := event.Amount
	event.Amount = int(float64(event.Amount) * e.multiplier)

	// Add metadata
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["damage_multiplied_by"] = e.ID()
	event.Metadata["original_amount"] = fmt.Sprintf("%d", originalAmount)
	event.Metadata["multiplier"] = fmt.Sprintf("%.1f", e.multiplier)
	event.Metadata["description"] = e.description

	// Not completely replaced - damage still happens with modified amount
	return event, false
}
