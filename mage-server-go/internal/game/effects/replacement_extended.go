package effects

import (
	"fmt"
	"strings"

	"github.com/magefree/mage-server-go/internal/game/rules"
)

// EntersBattlefieldReplacementEffect modifies how a permanent enters the battlefield
// Example: "As ~ enters the battlefield, choose a color" or "~ enters tapped"
// Rule 614.12: These effects have self-scope
type EntersBattlefieldReplacementEffect struct {
	*BaseReplacementEffect
	cardID      string                                      // Card entering battlefield
	modifyCard  func(event rules.Event) rules.Event         // Function to modify the card/event
	condition   func(event rules.Event, gameID string) bool // Condition for applying
	description string
}

// NewEntersBattlefieldReplacementEffect creates an ETB replacement effect
func NewEntersBattlefieldReplacementEffect(
	sourceID string,
	cardID string,
	modifyCard func(event rules.Event) rules.Event,
	condition func(event rules.Event, gameID string) bool,
	description string,
) *EntersBattlefieldReplacementEffect {
	return &EntersBattlefieldReplacementEffect{
		BaseReplacementEffect: NewBaseReplacementEffect(sourceID, DurationPermanent, false, true), // Has self-scope
		cardID:                strings.TrimSpace(cardID),
		modifyCard:            modifyCard,
		condition:             condition,
		description:           description,
	}
}

func (e *EntersBattlefieldReplacementEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventZoneChange ||
		eventType == rules.EventEntersBattlefield
}

func (e *EntersBattlefieldReplacementEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Check if entering battlefield
	if event.Zone != ZoneBattlefield {
		return false
	}

	// Check specific card if specified
	if e.cardID != "" && event.TargetID != e.cardID {
		return false
	}

	// Check condition if provided
	if e.condition != nil && !e.condition(event, gameID) {
		return false
	}

	return true
}

func (e *EntersBattlefieldReplacementEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Apply the modification function
	if e.modifyCard != nil {
		event = e.modifyCard(event)
	}

	// Add metadata
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["etb_replacement"] = e.ID()
	event.Metadata["description"] = e.description

	// Not completely replaced - card still enters, just modified
	return event, false
}

// Zone constants (these should ideally come from a shared package)
const (
	ZoneLibrary     = 0
	ZoneHand        = 1
	ZoneBattlefield = 2
	ZoneGraveyard   = 3
	ZoneStack       = 4
	ZoneExile       = 5
	ZoneCommand     = 6
	ZoneOutside     = 7
)

// DrawReplacementEffect replaces or modifies card draws
// Examples: "If you would draw a card, instead...", Abundance, Underworld Dreams triggers
type DrawReplacementEffect struct {
	*BaseReplacementEffect
	playerID    string                                      // Player drawing (empty = any)
	replacement func(event rules.Event) rules.Event         // What happens instead
	condition   func(event rules.Event, gameID string) bool // Optional condition
	description string
}

// NewDrawReplacementEffect creates a draw replacement effect
func NewDrawReplacementEffect(
	sourceID string,
	playerID string,
	replacement func(event rules.Event) rules.Event,
	condition func(event rules.Event, gameID string) bool,
	description string,
) *DrawReplacementEffect {
	return &DrawReplacementEffect{
		BaseReplacementEffect: NewBaseReplacementEffect(sourceID, DurationPermanent, false, false),
		playerID:              strings.TrimSpace(playerID),
		replacement:           replacement,
		condition:             condition,
		description:           description,
	}
}

func (e *DrawReplacementEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventDrawCard ||
		eventType == rules.EventDrawCards
}

func (e *DrawReplacementEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Check player if specified
	if e.playerID != "" && event.TargetID != e.playerID {
		return false
	}

	// Check condition if provided
	if e.condition != nil && !e.condition(event, gameID) {
		return false
	}

	return true
}

func (e *DrawReplacementEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Apply replacement
	if e.replacement != nil {
		event = e.replacement(event)
	}

	// Add metadata
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["draw_replacement"] = e.ID()
	event.Metadata["description"] = e.description

	// Check if event was completely replaced (no cards drawn)
	completelyReplaced := event.Amount == 0

	return event, completelyReplaced
}

// DiesReplacementEffect replaces death with something else
// Examples: Undying, Persist, "If ~ would die, regenerate it instead"
type DiesReplacementEffect struct {
	*BaseReplacementEffect
	cardID        string                                      // Specific card (empty = any)
	controllerID  string                                      // Controller (empty = any)
	cardTypeCheck string                                      // Card type requirement
	replacement   func(event rules.Event) rules.Event         // What happens instead
	condition     func(event rules.Event, gameID string) bool // Optional condition
	oneTimeUse    bool                                        // Remove after first use (Undying/Persist)
	description   string
}

// NewDiesReplacementEffect creates a dies replacement effect
func NewDiesReplacementEffect(
	sourceID string,
	cardID, controllerID, cardTypeCheck string,
	replacement func(event rules.Event) rules.Event,
	condition func(event rules.Event, gameID string) bool,
	oneTimeUse bool,
	description string,
) *DiesReplacementEffect {
	duration := DurationPermanent
	if oneTimeUse {
		duration = DurationOneUse
	}

	return &DiesReplacementEffect{
		BaseReplacementEffect: NewBaseReplacementEffect(sourceID, duration, false, false),
		cardID:                strings.TrimSpace(cardID),
		controllerID:          strings.TrimSpace(controllerID),
		cardTypeCheck:         strings.TrimSpace(cardTypeCheck),
		replacement:           replacement,
		condition:             condition,
		oneTimeUse:            oneTimeUse,
		description:           description,
	}
}

func (e *DiesReplacementEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventZoneChange && eventType == rules.EventDies
}

func (e *DiesReplacementEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Check if going to graveyard (dies)
	if event.Zone != ZoneGraveyard {
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

	// Check condition if provided
	if e.condition != nil && !e.condition(event, gameID) {
		return false
	}

	return true
}

func (e *DiesReplacementEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Apply replacement
	if e.replacement != nil {
		event = e.replacement(event)
	}

	// Add metadata
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["dies_replacement"] = e.ID()
	event.Metadata["description"] = e.description

	if e.oneTimeUse {
		event.Metadata["one_time_use"] = "true"
	}

	// Check if completely replaced (card didn't go to graveyard)
	completelyReplaced := event.Zone != ZoneGraveyard

	return event, completelyReplaced
}

// TokenCreationReplacementEffect modifies token creation
// Example: "If an effect would create one or more tokens, it creates twice that many instead" (Doubling Season)
type TokenCreationReplacementEffect struct {
	*BaseReplacementEffect
	multiplier  int                                         // Multiplier for token count (2 = double)
	filter      func(event rules.Event, gameID string) bool // Filter for which tokens
	description string
}

// NewTokenCreationReplacementEffect creates a token creation replacement effect
func NewTokenCreationReplacementEffect(
	sourceID string,
	multiplier int,
	filter func(event rules.Event, gameID string) bool,
	description string,
) *TokenCreationReplacementEffect {
	return &TokenCreationReplacementEffect{
		BaseReplacementEffect: NewBaseReplacementEffect(sourceID, DurationPermanent, false, false),
		multiplier:            multiplier,
		filter:                filter,
		description:           description,
	}
}

func (e *TokenCreationReplacementEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventCreateToken ||
		eventType == rules.EventCreateTokens
}

func (e *TokenCreationReplacementEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Check filter if provided
	if e.filter != nil && !e.filter(event, gameID) {
		return false
	}

	return true
}

func (e *TokenCreationReplacementEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Multiply the token count
	event.Amount = event.Amount * e.multiplier

	// Add metadata
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["token_multiplier"] = fmt.Sprintf("%d", e.multiplier)
	event.Metadata["description"] = e.description

	// Not completely replaced - tokens still created, just more of them
	return event, false
}

// CounterPlacementReplacementEffect modifies counter placement
// Example: "If one or more counters would be put on a permanent, twice that many are put instead" (Doubling Season, Vorinclex)
type CounterPlacementReplacementEffect struct {
	*BaseReplacementEffect
	multiplier   int                                         // Multiplier for counter count
	counterTypes []string                                    // Counter types this applies to (empty = all)
	filter       func(event rules.Event, gameID string) bool // Filter for which permanents
	description  string
}

// NewCounterPlacementReplacementEffect creates a counter placement replacement effect
func NewCounterPlacementReplacementEffect(
	sourceID string,
	multiplier int,
	counterTypes []string,
	filter func(event rules.Event, gameID string) bool,
	description string,
) *CounterPlacementReplacementEffect {
	return &CounterPlacementReplacementEffect{
		BaseReplacementEffect: NewBaseReplacementEffect(sourceID, DurationPermanent, false, false),
		multiplier:            multiplier,
		counterTypes:          counterTypes,
		filter:                filter,
		description:           description,
	}
}

func (e *CounterPlacementReplacementEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventAddCounter ||
		eventType == rules.EventAddCounters
}

func (e *CounterPlacementReplacementEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Check counter type if specified
	if len(e.counterTypes) > 0 {
		counterType := event.Metadata["counter_type"]
		found := false
		for _, ct := range e.counterTypes {
			if ct == counterType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check filter if provided
	if e.filter != nil && !e.filter(event, gameID) {
		return false
	}

	return true
}

func (e *CounterPlacementReplacementEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Multiply the counter count
	event.Amount = event.Amount * e.multiplier

	// Add metadata
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["counter_multiplier"] = fmt.Sprintf("%d", e.multiplier)
	event.Metadata["description"] = e.description

	// Not completely replaced - counters still placed, just more of them
	return event, false
}

// DiscardReplacementEffect modifies or replaces discard events
// Example: "If you would discard a card, exile it instead" (Library of Leng effect)
type DiscardReplacementEffect struct {
	*BaseReplacementEffect
	playerID    string                                      // Player discarding (empty = any)
	newZone     int                                         // New destination zone (-1 = no change)
	condition   func(event rules.Event, gameID string) bool // Optional condition
	description string
}

// NewDiscardReplacementEffect creates a discard replacement effect
func NewDiscardReplacementEffect(
	sourceID string,
	playerID string,
	newZone int,
	condition func(event rules.Event, gameID string) bool,
	description string,
) *DiscardReplacementEffect {
	return &DiscardReplacementEffect{
		BaseReplacementEffect: NewBaseReplacementEffect(sourceID, DurationPermanent, false, false),
		playerID:              strings.TrimSpace(playerID),
		newZone:               newZone,
		condition:             condition,
		description:           description,
	}
}

func (e *DiscardReplacementEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventDiscard ||
		eventType == rules.EventDiscardCard
}

func (e *DiscardReplacementEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Check player if specified
	if e.playerID != "" && event.Controller != e.playerID {
		return false
	}

	// Check condition if provided
	if e.condition != nil && !e.condition(event, gameID) {
		return false
	}

	return true
}

func (e *DiscardReplacementEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Change destination zone if specified
	if e.newZone >= 0 {
		event.Zone = e.newZone
	}

	// Add metadata
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["discard_replacement"] = e.ID()
	event.Metadata["description"] = e.description

	if e.newZone >= 0 {
		event.Metadata["new_zone"] = fmt.Sprintf("%d", e.newZone)
	}

	// Not completely replaced - card still discarded (or moved elsewhere)
	return event, false
}

// RegenerationEffect is a special prevention/replacement effect for regeneration
// Rule 701.15: Regeneration creates a replacement effect that lasts until end of turn
type RegenerationEffect struct {
	*BasePreventionEffect
	permanentID  string // Permanent that would be destroyed
	usedThisTurn bool   // Track if already used this turn
}

// NewRegenerationEffect creates a regeneration effect
func NewRegenerationEffect(sourceID, permanentID string) *RegenerationEffect {
	return &RegenerationEffect{
		BasePreventionEffect: NewBasePreventionEffect(sourceID, DurationEndOfTurn, 0),
		permanentID:          strings.TrimSpace(permanentID),
		usedThisTurn:         false,
	}
}

func (e *RegenerationEffect) ChecksEventType(eventType rules.EventType) bool {
	return eventType == rules.EventDestroy ||
		eventType == rules.EventZoneChange // When going to graveyard from lethal damage
}

func (e *RegenerationEffect) Applies(event rules.Event, gameID string) bool {
	if !e.ChecksEventType(event.Type) {
		return false
	}

	// Already used this turn
	if e.usedThisTurn {
		return false
	}

	// Check if it's our permanent
	if e.permanentID != "" && event.TargetID != e.permanentID {
		return false
	}

	// For zone changes, only apply if going to graveyard from battlefield
	if event.Type == rules.EventZoneChange {
		if event.Zone != ZoneGraveyard {
			return false
		}
		// TODO: Check that it's coming from battlefield
	}

	return true
}

func (e *RegenerationEffect) ReplaceEvent(event rules.Event, gameID string) (rules.Event, bool) {
	// Mark as used
	e.usedThisTurn = true

	// Instead of being destroyed/going to graveyard:
	// 1. Tap the permanent (if not already tapped)
	// 2. Remove all damage from it
	// 3. Remove it from combat (if it's in combat)

	// Add metadata
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["regenerated"] = "true"
	event.Metadata["regeneration_effect"] = e.ID()
	event.Metadata["tap_permanent"] = "true"
	event.Metadata["remove_damage"] = "true"
	event.Metadata["remove_from_combat"] = "true"

	// For destroy events, prevent destruction
	if event.Type == rules.EventDestroy {
		// Completely prevent the destruction
		return event, true
	}

	// For zone change to graveyard, prevent the zone change
	if event.Type == rules.EventZoneChange && event.Zone == ZoneGraveyard {
		// Keep permanent on battlefield
		event.Zone = ZoneBattlefield
		return event, false // Modified but not completely replaced
	}

	return event, false
}

func (e *RegenerationEffect) ResetForNewTurn() {
	e.usedThisTurn = false
}
