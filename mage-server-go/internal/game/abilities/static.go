package abilities

import (
	"context"

	"github.com/google/uuid"
)

// ========================================
// Static Ability System
// ========================================

// StaticAbility represents a continuous effect that is always active
// Java: mage.abilities.common.SimpleStaticAbility
// MTG Rules: 604 (Static Abilities), 611 (Continuous Effects)
type StaticAbility struct {
	baseAbility
	zone    Zone               // Where must the source be for this to be active?
	effects []ContinuousEffect // Effects that are continuously applied
}

// NewStaticAbility creates a new static ability
func NewStaticAbility(sourceID uuid.UUID, zone Zone, effects []ContinuousEffect) *StaticAbility {
	text := "static ability"
	if len(effects) > 0 {
		text = effects[0].GetDescription()
	}

	return &StaticAbility{
		baseAbility: newBaseAbility(sourceID, text),
		zone:        zone,
		effects:     effects,
	}
}

// GetType returns the ability type
func (a *StaticAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

// CanActivate checks if this static ability is active
func (a *StaticAbility) CanActivate(ctx context.Context, game GameContext) bool {
	// Static abilities are always "active" - they don't activate
	// The game engine checks the zone to determine if effects apply
	return true
}

// Resolve does nothing (static abilities don't resolve)
func (a *StaticAbility) Resolve(ctx context.Context, game GameContext) error {
	// Static abilities create continuous effects
	// They don't resolve like activated/triggered abilities
	return nil
}

// GetZone returns the zone this ability functions in
func (a *StaticAbility) GetZone() Zone {
	return a.zone
}

// GetEffects returns the continuous effects
func (a *StaticAbility) GetEffects() []ContinuousEffect {
	return a.effects
}

// IsActive checks if this ability is active based on source zone
func (a *StaticAbility) IsActive(sourceZone Zone) bool {
	return sourceZone == a.zone
}

// ========================================
// Zone System
// ========================================

// Zone represents where a card/permanent can be
// MTG Rules: 400 (Zones)
type Zone int

const (
	ZoneLibrary Zone = iota
	ZoneHand
	ZoneBattlefield
	ZoneGraveyard
	ZoneStack
	ZoneExile
	ZoneCommand
	ZoneOutside // Outside the game (sideboard, etc.)
)

func (z Zone) String() string {
	switch z {
	case ZoneLibrary:
		return "Library"
	case ZoneHand:
		return "Hand"
	case ZoneBattlefield:
		return "Battlefield"
	case ZoneGraveyard:
		return "Graveyard"
	case ZoneStack:
		return "Stack"
	case ZoneExile:
		return "Exile"
	case ZoneCommand:
		return "Command"
	case ZoneOutside:
		return "Outside"
	default:
		return "Unknown"
	}
}

// Note: Duration is defined in effects.go to avoid circular dependencies

// ========================================
// Continuous Effect Interface
// ========================================

// ContinuousEffect represents an effect that is applied continuously
// Java: mage.abilities.effects.ContinuousEffect
// MTG Rules: 611 (Continuous Effects), 613 (Interaction of Continuous Effects / Layers)
type ContinuousEffect interface {
	Effect

	// GetLayer returns which layer this effect applies in
	// MTG Rules 613: Layer system determines order of continuous effects
	GetLayer() Layer

	// GetDuration returns how long this effect lasts
	GetDuration() Duration

	// IsActive checks if this effect is currently active
	IsActive(ctx context.Context, game GameContext, source uuid.UUID) bool
}

// ========================================
// Layer System
// ========================================

// Layer represents the order in which continuous effects are applied
// Java: mage.constants.Layer
// MTG Rules: 613 (Interaction of Continuous Effects)
//
// The layer system ensures continuous effects are applied in the correct order:
// 1. Copy effects
// 2. Control-changing effects
// 3. Text-changing effects
// 4. Type-changing effects
// 5. Color-changing effects
// 6. Ability-adding/removing effects
// 7. Power/toughness effects
type Layer int

const (
	// LayerCopyEffects - Layer 1: Copy effects (Clone, etc.)
	LayerCopyEffects Layer = iota + 1

	// LayerControlChanging - Layer 2: Control-changing effects
	LayerControlChanging

	// LayerTextChanging - Layer 3: Text-changing effects
	LayerTextChanging

	// LayerTypeChanging - Layer 4: Type-changing effects
	LayerTypeChanging

	// LayerColorChanging - Layer 5: Color-changing effects
	LayerColorChanging

	// LayerAbilityAddingRemoving - Layer 6: Ability-adding/removing effects
	LayerAbilityAddingRemoving

	// LayerPowerToughnessEffects - Layer 7: Power/toughness effects
	LayerPowerToughnessEffects

	// LayerOther - Used for effects that don't fit standard layers
	LayerOther
)

func (l Layer) String() string {
	switch l {
	case LayerCopyEffects:
		return "Copy Effects"
	case LayerControlChanging:
		return "Control Changing"
	case LayerTextChanging:
		return "Text Changing"
	case LayerTypeChanging:
		return "Type Changing"
	case LayerColorChanging:
		return "Color Changing"
	case LayerAbilityAddingRemoving:
		return "Ability Adding/Removing"
	case LayerPowerToughnessEffects:
		return "Power/Toughness Effects"
	default:
		return "Unknown"
	}
}

// ========================================
// Base Continuous Effect
// ========================================

// baseContinuousEffect provides common fields for continuous effects
type baseContinuousEffect struct {
	layer    Layer
	duration Duration
}

// GetLayer returns the layer
func (e *baseContinuousEffect) GetLayer() Layer {
	return e.layer
}

// GetDuration returns the duration
func (e *baseContinuousEffect) GetDuration() Duration {
	return e.duration
}

// IsActive checks if the effect is active (default: always active)
func (e *baseContinuousEffect) IsActive(ctx context.Context, game GameContext, source uuid.UUID) bool {
	// Default: check duration
	// Subclasses can override for more complex logic
	switch e.duration {
	case DurationPermanent:
		return true
	case DurationWhileOnBattlefield:
		// TODO: Check if source is on battlefield
		return true
	default:
		return true
	}
}

// ========================================
// Static Ability Builder
// ========================================

// StaticAbilityBuilder provides a fluent API for building static abilities
type StaticAbilityBuilder struct {
	sourceID uuid.UUID
	zone     Zone
	effects  []ContinuousEffect
}

// NewStaticAbilityBuilder creates a new static ability builder
func NewStaticAbilityBuilder(sourceID uuid.UUID, zone Zone) *StaticAbilityBuilder {
	return &StaticAbilityBuilder{
		sourceID: sourceID,
		zone:     zone,
		effects:  make([]ContinuousEffect, 0),
	}
}

// AddEffect adds a continuous effect to this static ability
func (b *StaticAbilityBuilder) AddEffect(effect ContinuousEffect) *StaticAbilityBuilder {
	b.effects = append(b.effects, effect)
	return b
}

// Build constructs the static ability
func (b *StaticAbilityBuilder) Build() *StaticAbility {
	return NewStaticAbility(b.sourceID, b.zone, b.effects)
}

// ========================================
// Simple Static Ability (convenience function)
// ========================================

// NewSimpleStaticAbility creates a simple static ability that works on the battlefield
// This is the most common case - equivalent to Java's SimpleStaticAbility
func NewSimpleStaticAbility(sourceID uuid.UUID, zone Zone) *StaticAbilityBuilder {
	return NewStaticAbilityBuilder(sourceID, zone)
}
