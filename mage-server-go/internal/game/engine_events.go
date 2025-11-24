package game

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/rules"
)

// EventAdapter converts rules events to abilities events for trigger checking.
// This bridges the existing rules.Event system with the new abilities.GameEvent system.
type EventAdapter struct{}

// NewEventAdapter creates a new event adapter.
func NewEventAdapter() *EventAdapter {
	return &EventAdapter{}
}

// ConvertEvent converts a rules.Event to an abilities.GameEvent.
func (ea *EventAdapter) ConvertEvent(re *rules.Event) *abilities.GameEvent {
	if re == nil {
		return nil
	}

	ae := &abilities.GameEvent{
		Type:     ea.mapEventType(re.Type),
		SourceID: parseUUID(re.SourceID),
		TargetID: parseUUID(re.TargetID),
		PlayerID: parseUUID(re.PlayerID),
		Amount:   re.Amount,
	}

	// Note: The rules.Event struct uses a single Zone field
	// For zone change events, the event type determines the direction
	// (e.g., EventEntersBattlefield, EventLeavesBattlefield)
	// The Zone field in rules.Event is an int, not easily mappable to abilities.Zone
	// For now, we'll set default zones based on event type in mapEventType

	return ae
}

// mapEventType maps rules.EventType to abilities.EventType.
func (ea *EventAdapter) mapEventType(rt rules.EventType) abilities.EventType {
	switch rt {
	// Zone changes - need to inspect FromZone/ToZone for specific mapping
	case rules.EventZoneChange:
		// Will need to check FromZone and ToZone at call site
		return abilities.EventEntersBattlefield // Default

	// Damage events
	case rules.EventDamagePlayer, rules.EventDamagedPlayer:
		return abilities.EventDamageDealt
	case rules.EventDamagedBatchForOnePlayer:
		return abilities.EventDamageDealt

	// Life events
	case rules.EventGainLife, rules.EventGainedLife:
		return abilities.EventLifeGained
	case rules.EventLoseLife, rules.EventLostLife:
		return abilities.EventLifeLost

	// Card draw
	case rules.EventDrawCard, rules.EventDrewCard:
		return abilities.EventCardDrawn

	// Discard
	case rules.EventDiscardCard, rules.EventDiscardedCard:
		return abilities.EventCardDiscarded

	// Spell/ability
	case rules.EventCastSpell, rules.EventSpellCast:
		return abilities.EventSpellCast
	case rules.EventActivateAbility, rules.EventActivatedAbility:
		return abilities.EventAbilityActivated
	case rules.EventTriggeredAbility:
		return abilities.EventEntersBattlefield // Placeholder

	// Phase/step
	case rules.EventChangePhase, rules.EventPhaseChanged:
		return abilities.EventPhaseBegin
	case rules.EventChangeStep, rules.EventStepChanged:
		return abilities.EventStepBegin

	// Untap step
	case rules.EventUntapStep:
		return abilities.EventUntapped

	// Counter events
	case rules.EventCounterAdded:
		return abilities.EventCounterAdded

	// Attack/block
	case rules.EventAttackerDeclared:
		return abilities.EventAttackerDeclared
	case rules.EventBlockerDeclared:
		return abilities.EventBlockerDeclared

	default:
		// For unmapped events, return a generic event type
		return abilities.EventType(0) // Unknown
	}
}

// mapZone maps string zone names to abilities.Zone constants.
func (ea *EventAdapter) mapZone(zoneName string) abilities.Zone {
	switch zoneName {
	case "LIBRARY":
		return abilities.ZoneLibrary
	case "HAND":
		return abilities.ZoneHand
	case "BATTLEFIELD":
		return abilities.ZoneBattlefield
	case "GRAVEYARD":
		return abilities.ZoneGraveyard
	case "STACK":
		return abilities.ZoneStack
	case "EXILE":
		return abilities.ZoneExile
	case "COMMAND":
		return abilities.ZoneCommand
	case "OUTSIDE":
		return abilities.ZoneOutside
	default:
		return abilities.ZoneBattlefield // Default
	}
}

// DetermineZoneChangeEventType determines the specific event type for a zone change.
// This handles the "enters the battlefield", "dies", "leaves the battlefield" cases.
func (ea *EventAdapter) DetermineZoneChangeEventType(fromZone, toZone string, cardTypes []string) abilities.EventType {
	// Dies: battlefield → graveyard (creature/planeswalker)
	if fromZone == "BATTLEFIELD" && toZone == "GRAVEYARD" {
		for _, ct := range cardTypes {
			if ct == "CREATURE" || ct == "PLANESWALKER" {
				return abilities.EventDies
			}
		}
		return abilities.EventLeavesBattlefield
	}

	// Enters the battlefield: any zone → battlefield
	if toZone == "BATTLEFIELD" {
		return abilities.EventEntersBattlefield
	}

	// Leaves the battlefield: battlefield → any other zone
	if fromZone == "BATTLEFIELD" {
		return abilities.EventLeavesBattlefield
	}

	// Generic zone change
	return abilities.EventType(0) // Unknown
}

// parseUUID safely parses a string to UUID, returning zero UUID on error.
func parseUUID(s string) uuid.UUID {
	if s == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}
