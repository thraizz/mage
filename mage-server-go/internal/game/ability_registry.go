package game

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/abilities"
)

// AbilityRegistry manages the mapping between ability IDs and ability objects
// This is used to retrieve abilities from cards and permanents for activation
type AbilityRegistry struct {
	mu        sync.RWMutex
	abilities map[string]abilities.Ability // abilityID -> ability
	sources   map[string][]string          // sourceID -> []abilityID
	metadata  map[string]*AbilityMetadata  // abilityID -> metadata
}

// AbilityMetadata stores additional information about an ability
type AbilityMetadata struct {
	SourceID   uuid.UUID
	Controller uuid.UUID
	Index      int            // Index in the source's ability list
	Zone       abilities.Zone // Zone where ability functions
}

// NewAbilityRegistry creates a new ability registry
func NewAbilityRegistry() *AbilityRegistry {
	return &AbilityRegistry{
		abilities: make(map[string]abilities.Ability),
		sources:   make(map[string][]string),
		metadata:  make(map[string]*AbilityMetadata),
	}
}

// RegisterAbility registers an ability with the registry
func (ar *AbilityRegistry) RegisterAbility(
	ability abilities.Ability,
	controller uuid.UUID,
	index int,
	zone abilities.Zone,
) {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	abilityID := ability.GetID().String()
	sourceID := ability.GetSourceID().String()

	// Store the ability
	ar.abilities[abilityID] = ability

	// Track abilities by source
	if _, ok := ar.sources[sourceID]; !ok {
		ar.sources[sourceID] = make([]string, 0, 4)
	}
	ar.sources[sourceID] = append(ar.sources[sourceID], abilityID)

	// Store metadata
	ar.metadata[abilityID] = &AbilityMetadata{
		SourceID:   ability.GetSourceID(),
		Controller: controller,
		Index:      index,
		Zone:       zone,
	}
}

// GetAbility retrieves an ability by its ID
func (ar *AbilityRegistry) GetAbility(abilityID uuid.UUID) (abilities.Ability, error) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	ability, ok := ar.abilities[abilityID.String()]
	if !ok {
		return nil, fmt.Errorf("ability %s not found", abilityID.String())
	}

	return ability, nil
}

// GetAbilitiesBySource retrieves all abilities for a given source (card/permanent)
func (ar *AbilityRegistry) GetAbilitiesBySource(sourceID uuid.UUID) []abilities.Ability {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	abilityIDs, ok := ar.sources[sourceID.String()]
	if !ok {
		return []abilities.Ability{}
	}

	result := make([]abilities.Ability, 0, len(abilityIDs))
	for _, abilityID := range abilityIDs {
		if ability, ok := ar.abilities[abilityID]; ok {
			result = append(result, ability)
		}
	}

	return result
}

// GetMetadata retrieves metadata for an ability
func (ar *AbilityRegistry) GetMetadata(abilityID uuid.UUID) (*AbilityMetadata, error) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	metadata, ok := ar.metadata[abilityID.String()]
	if !ok {
		return nil, fmt.Errorf("ability metadata not found for %s", abilityID.String())
	}

	return metadata, nil
}

// UnregisterSource removes all abilities for a given source
// This should be called when a card/permanent leaves the battlefield or changes zones
func (ar *AbilityRegistry) UnregisterSource(sourceID uuid.UUID) {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	abilityIDs, ok := ar.sources[sourceID.String()]
	if !ok {
		return
	}

	// Remove all abilities for this source
	for _, abilityID := range abilityIDs {
		delete(ar.abilities, abilityID)
		delete(ar.metadata, abilityID)
	}

	// Remove the source entry
	delete(ar.sources, sourceID.String())
}

// GetActivatableAbilities returns all abilities that can be activated by a player
// This filters for activated abilities and checks zone restrictions
func (ar *AbilityRegistry) GetActivatableAbilities(
	playerID uuid.UUID,
	zone abilities.Zone,
) []abilities.Ability {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	result := make([]abilities.Ability, 0)

	for abilityID, ability := range ar.abilities {
		metadata := ar.metadata[abilityID]

		// Check if player controls the source
		if metadata.Controller != playerID {
			continue
		}

		// Check ability type (only activated and mana abilities can be activated)
		abilityType := ability.GetType()
		if abilityType != abilities.AbilityTypeActivated && abilityType != abilities.AbilityTypeMana {
			continue
		}

		// Check if ability is in the correct zone (using metadata)
		if metadata.Zone != zone {
			continue
		}

		// For static abilities, also check the ability's own zone property
		if staticAbility, ok := ability.(*abilities.StaticAbility); ok {
			if staticAbility.GetZone() != zone {
				continue
			}
		}

		result = append(result, ability)
	}

	return result
}

// RegisterCardAbilities registers all abilities from a card
func (ar *AbilityRegistry) RegisterCardAbilities(card *LegacyCard) {
	if card.Abilities == nil || len(card.Abilities) == 0 {
		return
	}

	currentZone := convertZone(card.Zone)

	for i, abilityInterface := range card.Abilities {
		// Type assert to abilities.Ability
		if ability, ok := abilityInterface.(abilities.Ability); ok {
			ar.RegisterAbility(ability, card.ControllerID, i, currentZone)
		}
	}
}

// UpdateCardZone updates the zone information for all abilities of a card
// This should be called when a card changes zones
func (ar *AbilityRegistry) UpdateCardZone(cardID uuid.UUID, newZone abilities.Zone) {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	abilityIDs, ok := ar.sources[cardID.String()]
	if !ok {
		return
	}

	for _, abilityID := range abilityIDs {
		if metadata, ok := ar.metadata[abilityID]; ok {
			metadata.Zone = newZone
		}
	}
}

// convertZone converts from game.Zone to abilities.Zone
func convertZone(z Zone) abilities.Zone {
	switch z {
	case ZoneLibrary:
		return abilities.ZoneLibrary
	case ZoneHand:
		return abilities.ZoneHand
	case ZoneBattlefield:
		return abilities.ZoneBattlefield
	case ZoneGraveyard:
		return abilities.ZoneGraveyard
	case ZoneStack:
		return abilities.ZoneStack
	case ZoneExile:
		return abilities.ZoneExile
	case ZoneCommand:
		return abilities.ZoneCommand
	default:
		return abilities.ZoneBattlefield
	}
}
