package game

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap"
)

// ContinuousEffectsManager manages the recalculation of continuous effects layers
// This implements MTG Rule 613 - Continuous Effects
type ContinuousEffectsManager struct {
	layerSystem     *effects.LayerSystem
	abilityRegistry *AbilityRegistry
	logger          *zap.Logger
	mu              sync.RWMutex
	effectMappings  map[string][]string // abilityID -> []effectID
}

// NewContinuousEffectsManager creates a new continuous effects manager
func NewContinuousEffectsManager(
	layerSystem *effects.LayerSystem,
	abilityRegistry *AbilityRegistry,
	logger *zap.Logger,
) *ContinuousEffectsManager {
	return &ContinuousEffectsManager{
		layerSystem:     layerSystem,
		abilityRegistry: abilityRegistry,
		logger:          logger,
		effectMappings:  make(map[string][]string),
	}
}

// RecalculateAll recalculates all continuous effects for all permanents
// This is called when state-based actions are checked (Rule 704.3)
func (cem *ContinuousEffectsManager) RecalculateAll(
	ctx context.Context,
	state rules.GameStateReader,
) error {
	cem.mu.Lock()
	defer cem.mu.Unlock()

	cem.logger.Debug("recalculating continuous effects")

	// Clear existing effect mappings
	// We'll rebuild from scratch each time
	cem.effectMappings = make(map[string][]string)

	// If no state provided, nothing to process
	if state == nil {
		return nil
	}

	// Get all permanents on the battlefield
	permanents := state.GetAllPermanents()

	// Process each permanent and register its static abilities
	for i := range permanents {
		if err := cem.processPermanent(&permanents[i]); err != nil {
			cem.logger.Error("failed to process permanent for continuous effects",
				zap.String("permanent", permanents[i].ID.String()),
				zap.Error(err))
			continue
		}
	}

	cem.logger.Debug("continuous effects recalculation complete",
		zap.Int("permanents_processed", len(permanents)))

	return nil
}

// processPermanent processes a single permanent's static abilities
func (cem *ContinuousEffectsManager) processPermanent(permanent *rules.Permanent) error {
	// Get all abilities for this permanent from the registry
	// permanent.ID is already a uuid.UUID
	allAbilities := cem.abilityRegistry.GetAbilitiesBySource(permanent.ID)

	// Filter for static abilities
	for _, ability := range allAbilities {
		if ability.GetType() != abilities.AbilityTypeStatic {
			continue
		}

		// Type assert to static ability
		staticAbility, ok := ability.(*abilities.StaticAbility)
		if !ok {
			continue
		}

		// Only process battlefield abilities
		if staticAbility.GetZone() != abilities.ZoneBattlefield {
			continue
		}

		// Convert static ability effects to layer effects
		if err := cem.registerStaticAbility(staticAbility, permanent); err != nil {
			cem.logger.Warn("failed to register static ability",
				zap.String("ability", ability.GetID().String()),
				zap.Error(err))
			continue
		}
	}

	return nil
}

// registerStaticAbility converts a static ability's effects into layer effects
func (cem *ContinuousEffectsManager) registerStaticAbility(
	staticAbility *abilities.StaticAbility,
	permanent *rules.Permanent,
) error {
	abilityID := staticAbility.GetID().String()

	// Get the effects from the static ability
	abilityEffects := staticAbility.GetEffects()

	var registeredEffectIDs []string

	for _, effect := range abilityEffects {
		// Convert abilities.Effect to effects.ContinuousEffect
		continuousEffect := cem.convertToContinuousEffect(effect, permanent)
		if continuousEffect == nil {
			// Not a continuous effect (might be one-shot effect)
			continue
		}

		// Register the effect with the layer system
		effectID := cem.layerSystem.AddEffect(continuousEffect)
		registeredEffectIDs = append(registeredEffectIDs, effectID)
	}

	// Track which layer effects came from this ability
	if len(registeredEffectIDs) > 0 {
		cem.effectMappings[abilityID] = registeredEffectIDs
	}

	return nil
}

// convertToContinuousEffect converts an abilities.Effect to an effects.ContinuousEffect
// Returns nil if the effect is not a continuous effect
func (cem *ContinuousEffectsManager) convertToContinuousEffect(
	effect abilities.Effect,
	source *rules.Permanent,
) effects.ContinuousEffect {
	// Check if this is a continuous effect (has layer and duration)
	continuousEffect, ok := effect.(abilities.ContinuousEffect)
	if !ok {
		// Not a continuous effect (one-shot effect)
		return nil
	}

	// Only process permanent duration effects from static abilities
	// Temporary effects (until end of turn, etc.) are handled differently
	if continuousEffect.GetDuration() != abilities.DurationWhileOnBattlefield &&
		continuousEffect.GetDuration() != abilities.DurationPermanent {
		return nil
	}

	// Type switch on the effect to determine its layer effect type
	switch e := effect.(type) {
	case *abilities.BoostEffect:
		// Create a PT boost layer effect
		return effects.NewSimplePTBoostEffect(
			source.ID.String(),
			source.ControllerID.String(),
			e.Power,
			e.Toughness,
			true, // includeSelf - adjust based on effect's filter
		)

	// TODO: Add more effect type conversions:
	// - AddAbilityEffect -> Layer 6 effect
	// - SetPowerToughnessEffect -> Layer 7b effect
	// - ChangeTypeEffect -> Layer 4 effect
	// - ChangeColorEffect -> Layer 5 effect
	// - GainControlEffect -> Layer 2 effect

	default:
		// Not yet implemented for layer system
		cem.logger.Debug("continuous effect type not yet implemented for layer conversion",
			zap.String("effect_type", fmt.Sprintf("%T", effect)))
		return nil
	}
}

// ApplyToCard applies all continuous effects to a single card
// This calculates the card's current characteristics based on layers
func (cem *ContinuousEffectsManager) ApplyToCard(
	ctx context.Context,
	card *Card,
) error {
	cem.mu.RLock()
	defer cem.mu.RUnlock()

	// Create a snapshot from the card
	snapshot := cem.createSnapshot(card)

	// Apply all layer effects to the snapshot
	cem.layerSystem.Apply(snapshot)

	// Update the card with the snapshot results
	cem.applySnapshotToCard(snapshot, card)

	return nil
}

// ApplyToPermanent applies all continuous effects to a single permanent
// Note: This is a stub that needs to be implemented properly
// Currently permanents are stored separately from cards, so we can't
// directly apply card effects. This will be refactored when we unify
// the permanent/card representation.
func (cem *ContinuousEffectsManager) ApplyToPermanent(
	ctx context.Context,
	permanent *rules.Permanent,
	state rules.GameStateReader,
) error {
	cem.mu.RLock()
	defer cem.mu.RUnlock()

	// TODO: Once we unify permanent and card representations,
	// this should apply layer effects to modify the permanent's
	// characteristics (power, toughness, types, abilities, etc.)

	cem.logger.Debug("ApplyToPermanent called (stub)",
		zap.String("permanent", permanent.ID.String()))

	return nil
}

// createSnapshot creates a layer snapshot from a card
func (cem *ContinuousEffectsManager) createSnapshot(card *Card) *effects.Snapshot {
	// Parse power and toughness
	power := 0
	toughness := 0
	hasPower := false
	hasToughness := false

	if card.IsCreature() {
		// Parse power/toughness strings (e.g., "2", "3", "*", "1+*")
		// For now, simple integer parsing
		if card.Power != "" && card.Power != "*" {
			// Try to parse as int
			// TODO: Handle more complex P/T expressions
			power = 0 // placeholder
			hasPower = true
		}
		if card.Toughness != "" && card.Toughness != "*" {
			toughness = 0 // placeholder
			hasToughness = true
		}
	}

	return effects.NewSnapshot(
		card.ID.String(),
		card.ControllerID.String(),
		card.Types,
		power,
		toughness,
		hasPower,
		hasToughness,
	)
}

// applySnapshotToCard applies the snapshot results back to the card
func (cem *ContinuousEffectsManager) applySnapshotToCard(snapshot *effects.Snapshot, card *Card) {
	// Update power/toughness if the card is a creature
	if card.IsCreature() && snapshot.HasBasePower && snapshot.HasBaseTough {
		// Update the card's effective P/T
		// The snapshot has already been processed through all layers
		card.Power = string(rune(snapshot.Power + '0'))         // Convert int to string
		card.Toughness = string(rune(snapshot.Toughness + '0')) // Convert int to string
	}

	// Update types if they changed
	if len(snapshot.Types) > 0 {
		card.Types = snapshot.Types
	}
}

// RemoveAbilityEffects removes all layer effects associated with an ability
// This is called when a permanent with static abilities leaves the battlefield
func (cem *ContinuousEffectsManager) RemoveAbilityEffects(abilityID uuid.UUID) {
	cem.mu.Lock()
	defer cem.mu.Unlock()

	abilityIDStr := abilityID.String()
	effectIDs, ok := cem.effectMappings[abilityIDStr]
	if !ok {
		return
	}

	// Remove each effect from the layer system
	for _, effectID := range effectIDs {
		cem.layerSystem.RemoveEffect(effectID)
	}

	// Remove the mapping
	delete(cem.effectMappings, abilityIDStr)

	cem.logger.Debug("removed ability effects from layers",
		zap.String("ability", abilityIDStr),
		zap.Int("effects_removed", len(effectIDs)))
}

// RemoveSourceEffects removes all layer effects from a source (permanent)
// This is called when a permanent leaves the battlefield
func (cem *ContinuousEffectsManager) RemoveSourceEffects(sourceID uuid.UUID) {
	cem.mu.Lock()
	defer cem.mu.Unlock()

	// Get all abilities from this source
	abilities := cem.abilityRegistry.GetAbilitiesBySource(sourceID)

	for _, ability := range abilities {
		abilityIDStr := ability.GetID().String()
		effectIDs, ok := cem.effectMappings[abilityIDStr]
		if !ok {
			continue
		}

		// Remove each effect from the layer system
		for _, effectID := range effectIDs {
			cem.layerSystem.RemoveEffect(effectID)
		}

		// Remove the mapping
		delete(cem.effectMappings, abilityIDStr)
	}

	cem.logger.Debug("removed source effects from layers",
		zap.String("source", sourceID.String()),
		zap.Int("abilities_processed", len(abilities)))
}

// GetCharacteristics returns the current characteristics of a card after layers
func (cem *ContinuousEffectsManager) GetCharacteristics(
	ctx context.Context,
	card *Card,
) (*CardCharacteristics, error) {
	cem.mu.RLock()
	defer cem.mu.RUnlock()

	// Create and apply snapshot
	snapshot := cem.createSnapshot(card)
	cem.layerSystem.Apply(snapshot)

	// Return the characteristics
	return &CardCharacteristics{
		Power:     snapshot.Power,
		Toughness: snapshot.Toughness,
		Types:     snapshot.Types,
	}, nil
}

// CardCharacteristics represents the current characteristics of a card
type CardCharacteristics struct {
	Power     int
	Toughness int
	Types     []string
	// TODO: Add more characteristics as needed:
	// - Colors
	// - Abilities
	// - Subtypes
	// - Supertypes
}
