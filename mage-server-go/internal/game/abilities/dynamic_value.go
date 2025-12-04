package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DynamicValue calculates values at runtime based on game state.
// This mirrors Java's mage.abilities.dynamicvalue.DynamicValue interface.
type DynamicValue interface {
	// Calculate computes the value based on current game state
	Calculate(ctx context.Context, game GameContext, source uuid.UUID) int

	// GetMessage returns a text description of what this value represents
	GetMessage() string

	// Copy returns a copy of this dynamic value
	Copy() DynamicValue
}

// ========================================
// Static Value
// ========================================

// StaticValue is a DynamicValue that always returns a fixed integer.
// Mirrors Java's mage.abilities.dynamicvalue.common.StaticValue.
type StaticValue struct {
	value int
}

// NewStaticValue creates a new StaticValue with the given integer.
func NewStaticValue(value int) *StaticValue {
	return &StaticValue{value: value}
}

func (v *StaticValue) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	return v.value
}

func (v *StaticValue) GetMessage() string {
	return fmt.Sprintf("%d", v.value)
}

func (v *StaticValue) Copy() DynamicValue {
	return &StaticValue{value: v.value}
}

func (v *StaticValue) GetValue() int {
	return v.value
}

// toDynamicValue converts various types to DynamicValue.
// - If v is already a DynamicValue, returns it as-is.
// - If v is an int, wraps it in a StaticValue.
// - If v is nil, returns a StaticValue of 0.
func toDynamicValue(v interface{}) DynamicValue {
	if v == nil {
		return NewStaticValue(0)
	}
	switch val := v.(type) {
	case int:
		return NewStaticValue(val)
	case DynamicValue:
		// This also handles *StaticValue since it implements DynamicValue
		return val
	default:
		return NewStaticValue(0)
	}
}

// ========================================
// Dynamic Value Game Context
// ========================================

// DynamicValueGameContext extends GameContext with methods needed for dynamic value calculations.
// This should be implemented by the game engine.
type DynamicValueGameContext interface {
	GameContext

	// GetControllerID returns the controller ID of a permanent/card
	GetControllerID(objectID uuid.UUID) (uuid.UUID, error)

	// GetPermanentsControlledBy returns all permanents controlled by a player
	GetPermanentsControlledBy(ctx context.Context, playerID uuid.UUID) ([]PermanentInfo, error)

	// GetCardsInGraveyard returns all cards in a player's graveyard
	GetCardsInGraveyard(ctx context.Context, playerID uuid.UUID) ([]CardTypeInfo, error)

	// GetPlayerCounters returns the number of counters of a specific type on a player
	GetPlayerCounters(ctx context.Context, playerID uuid.UUID, counterType string) int

	// GetManaSpentToCast returns the total mana spent to cast a spell
	GetManaSpentToCast(ctx context.Context, spellID uuid.UUID) int
}

// PermanentInfo provides information about a permanent for dynamic value calculations.
type PermanentInfo interface {
	GetID() uuid.UUID
	GetManaCost() string
	GetTypes() []string
	GetSubtypes() []string
}

// CardTypeInfo provides card type information for filtering.
type CardTypeInfo interface {
	GetID() uuid.UUID
	GetTypes() []string
	GetSubtypes() []string
}
