package abilities

import (
	"context"
	"math"

	"github.com/google/uuid"
)

// ========================================
// Source Permanent Power Value
// ========================================

// SourcePermanentPowerValueType represents different modes for power value calculation
type SourcePermanentPowerValueType int

const (
	// PowerValueRaw returns the exact power value (can be negative)
	PowerValueRaw SourcePermanentPowerValueType = iota
	// PowerValueNotNegative returns max(0, power)
	PowerValueNotNegative
	// PowerValuePositive returns max(1, power)
	PowerValuePositive
)

// sourcePermanentPowerValue calculates a value based on the source permanent's power
// Java: mage.abilities.dynamicvalue.common.SourcePermanentPowerValue
type sourcePermanentPowerValue struct {
	valueType SourcePermanentPowerValueType
}

// SourcePermanentPowerValue provides pre-defined instances for common uses
// Java uses: SourcePermanentPowerValue.NOT_NEGATIVE, etc.
var SourcePermanentPowerValue = struct {
	RAW          DynamicValue
	NOT_NEGATIVE DynamicValue
	POSITIVE     DynamicValue
}{
	RAW:          &sourcePermanentPowerValue{valueType: PowerValueRaw},
	NOT_NEGATIVE: &sourcePermanentPowerValue{valueType: PowerValueNotNegative},
	POSITIVE:     &sourcePermanentPowerValue{valueType: PowerValuePositive},
}

func (v *sourcePermanentPowerValue) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	// Get the source permanent's power
	power := game.GetPermanentPower(ctx, source)

	switch v.valueType {
	case PowerValueNotNegative:
		if power < 0 {
			return 0
		}
		return power
	case PowerValuePositive:
		if power < 1 {
			return 1
		}
		return power
	default:
		return power
	}
}

func (v *sourcePermanentPowerValue) GetMessage() string {
	return "source permanent's power"
}

func (v *sourcePermanentPowerValue) Copy() DynamicValue {
	return &sourcePermanentPowerValue{valueType: v.valueType}
}

// ========================================
// Source Permanent Toughness Value
// ========================================

// sourcePermanentToughnessValue calculates a value based on the source permanent's toughness
// Java: mage.abilities.dynamicvalue.common.SourcePermanentToughnessValue
type sourcePermanentToughnessValue struct {
	notNegative bool
}

// SourcePermanentToughnessValue provides pre-defined instances
var SourcePermanentToughnessValue = struct {
	RAW          DynamicValue
	NOT_NEGATIVE DynamicValue
}{
	RAW:          &sourcePermanentToughnessValue{notNegative: false},
	NOT_NEGATIVE: &sourcePermanentToughnessValue{notNegative: true},
}

func (v *sourcePermanentToughnessValue) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	toughness := game.GetPermanentToughness(ctx, source)

	if v.notNegative && toughness < 0 {
		return 0
	}
	return toughness
}

func (v *sourcePermanentToughnessValue) GetMessage() string {
	return "source permanent's toughness"
}

func (v *sourcePermanentToughnessValue) Copy() DynamicValue {
	return &sourcePermanentToughnessValue{notNegative: v.notNegative}
}

// ========================================
// Greatest Among Permanents Value
// ========================================

// GreatestAmongPermanentsValueType represents what to find the greatest of
type GreatestAmongPermanentsValueType int

const (
	// GreatestPower finds the greatest power among permanents
	GreatestPower GreatestAmongPermanentsValueType = iota
	// GreatestToughness finds the greatest toughness among permanents
	GreatestToughness
	// GreatestManaValue finds the greatest mana value among permanents
	GreatestManaValue
	// GreatestManaValueControlled finds greatest mana value among permanents you control
	GreatestManaValueControlled
)

// greatestAmongPermanentsValue finds the greatest value among permanents
// Java: mage.abilities.dynamicvalue.common.GreatestAmongPermanentsValue
type greatestAmongPermanentsValue struct {
	valueType   GreatestAmongPermanentsValueType
	controlOnly bool // true = only your permanents, false = all permanents
}

// GreatestAmongPermanentsValue provides pre-defined instances
var GreatestAmongPermanentsValue = struct {
	POWER_CONTROLLED_CREATURES      DynamicValue
	TOUGHNESS_CONTROLLED_CREATURES  DynamicValue
	MANAVALUE_CONTROLLED_PERMANENTS DynamicValue
	MANAVALUE_ALL_PERMANENTS        DynamicValue
}{
	POWER_CONTROLLED_CREATURES:      &greatestAmongPermanentsValue{valueType: GreatestPower, controlOnly: true},
	TOUGHNESS_CONTROLLED_CREATURES:  &greatestAmongPermanentsValue{valueType: GreatestToughness, controlOnly: true},
	MANAVALUE_CONTROLLED_PERMANENTS: &greatestAmongPermanentsValue{valueType: GreatestManaValueControlled, controlOnly: true},
	MANAVALUE_ALL_PERMANENTS:        &greatestAmongPermanentsValue{valueType: GreatestManaValue, controlOnly: false},
}

func (v *greatestAmongPermanentsValue) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	dvGame, ok := game.(DynamicValueGameContext)
	if !ok {
		return 0
	}

	var permanents []PermanentInfo

	if v.controlOnly {
		controllerID, getControllerErr := dvGame.GetControllerID(source)
		if getControllerErr != nil {
			return 0
		}
		var getPermanentsErr error
		permanents, getPermanentsErr = dvGame.GetPermanentsControlledBy(ctx, controllerID)
		if getPermanentsErr != nil {
			return 0
		}
	} else {
		// TODO: Get all permanents on battlefield
		permanents = nil
	}

	greatest := 0
	for _, perm := range permanents {
		var value int
		switch v.valueType {
		case GreatestPower:
			value = game.GetPermanentPower(ctx, perm.GetID())
		case GreatestToughness:
			value = game.GetPermanentToughness(ctx, perm.GetID())
		case GreatestManaValue, GreatestManaValueControlled:
			value = parseManaValue(perm.GetManaCost())
		}
		if value > greatest {
			greatest = value
		}
	}

	return greatest
}

func (v *greatestAmongPermanentsValue) GetMessage() string {
	switch v.valueType {
	case GreatestPower:
		if v.controlOnly {
			return "the greatest power among creatures you control"
		}
		return "the greatest power among creatures"
	case GreatestToughness:
		if v.controlOnly {
			return "the greatest toughness among creatures you control"
		}
		return "the greatest toughness among creatures"
	case GreatestManaValue, GreatestManaValueControlled:
		if v.controlOnly {
			return "the greatest mana value among permanents you control"
		}
		return "the greatest mana value among all permanents"
	}
	return ""
}

func (v *greatestAmongPermanentsValue) Copy() DynamicValue {
	return &greatestAmongPermanentsValue{valueType: v.valueType, controlOnly: v.controlOnly}
}

// parseManaValue parses a mana cost string and returns its mana value
func parseManaValue(manaCost string) int {
	// Simple parsing - count symbols
	// TODO: Implement proper mana value calculation
	cost, err := ParseManaCost(manaCost)
	if err != nil {
		return 0
	}
	return cost.ConvertedManaCost()
}

// ========================================
// Target Permanent Power Value
// ========================================

// TargetPermanentPowerValue calculates a value based on a target permanent's power
// Java: mage.abilities.dynamicvalue.common.TargetPermanentPowerCount
type TargetPermanentPowerValue struct {
	notNegative bool
}

var TargetPermanentPowerValueInstance = &TargetPermanentPowerValue{notNegative: false}
var TargetPermanentPowerNotNegativeValue = &TargetPermanentPowerValue{notNegative: true}

func (v *TargetPermanentPowerValue) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	// TODO: Need to get the target from the ability
	// This requires knowing what the target is, which would come from the ability's target list
	return 0
}

func (v *TargetPermanentPowerValue) GetMessage() string {
	return "target permanent's power"
}

func (v *TargetPermanentPowerValue) Copy() DynamicValue {
	return &TargetPermanentPowerValue{notNegative: v.notNegative}
}

// ========================================
// Saved Damage Value
// ========================================

// SavedDamageValue represents damage that was dealt and saved for use by effects
// Java: mage.abilities.dynamicvalue.common.SavedDamageValue
// Used for effects like "whenever ~ deals damage, draw that many cards"
type SavedDamageValue struct{}

var SavedDamageValueInstance = &SavedDamageValue{}

func (v *SavedDamageValue) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	// TODO: Retrieve the saved damage amount from the game context
	// This is typically set when damage is dealt and accessed by effects that reference it
	return 0
}

func (v *SavedDamageValue) GetMessage() string {
	return "damage dealt"
}

func (v *SavedDamageValue) Copy() DynamicValue {
	return SavedDamageValueInstance
}

// ========================================
// X Value (Variable)
// ========================================

// xValue represents the X in variable costs like {X}{X}{R}
// Java: mage.abilities.dynamicvalue.common.GetXValue
type xValue struct {
	multiplier int
}

// GetXValue returns the X value from a spell or ability
// Java: GetXValue.instance
var GetXValue DynamicValue = &xValue{multiplier: 1}

// NewXValue creates a new X value with a multiplier
// e.g., for Fireball which does X damage but costs {X}{R}, use multiplier 1
// e.g., for a spell that does 2X damage, use multiplier 2
func NewXValue(multiplier int) DynamicValue {
	return &xValue{multiplier: multiplier}
}

func (v *xValue) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	// TODO: Get the X value from the spell/ability being resolved
	// This requires the game context to track the X value that was chosen
	return 0
}

func (v *xValue) GetMessage() string {
	if v.multiplier == 1 {
		return "X"
	}
	return "X times multiplier"
}

func (v *xValue) Copy() DynamicValue {
	return &xValue{multiplier: v.multiplier}
}

// xValue alias for generated code compatibility
// Some generated code uses lowercase xValue
var xValueInstance = GetXValue

// ========================================
// Multiplied Value
// ========================================

// MultipliedValue multiplies another dynamic value
type MultipliedValue struct {
	inner      DynamicValue
	multiplier int
}

func NewMultipliedValue(inner DynamicValue, multiplier int) *MultipliedValue {
	return &MultipliedValue{inner: inner, multiplier: multiplier}
}

func (v *MultipliedValue) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	return v.inner.Calculate(ctx, game, source) * v.multiplier
}

func (v *MultipliedValue) GetMessage() string {
	return v.inner.GetMessage()
}

func (v *MultipliedValue) Copy() DynamicValue {
	return &MultipliedValue{inner: v.inner.Copy(), multiplier: v.multiplier}
}

// ========================================
// Added Value
// ========================================

// AddedValue adds a constant to another dynamic value
type AddedValue struct {
	inner DynamicValue
	bonus int
}

func NewAddedValue(inner DynamicValue, bonus int) *AddedValue {
	return &AddedValue{inner: inner, bonus: bonus}
}

func (v *AddedValue) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	return v.inner.Calculate(ctx, game, source) + v.bonus
}

func (v *AddedValue) GetMessage() string {
	return v.inner.GetMessage()
}

func (v *AddedValue) Copy() DynamicValue {
	return &AddedValue{inner: v.inner.Copy(), bonus: v.bonus}
}

// ========================================
// Clamped Value (Min/Max)
// ========================================

// ClampedValue clamps another dynamic value to a range
type ClampedValue struct {
	inner DynamicValue
	min   *int
	max   *int
}

func NewClampedValue(inner DynamicValue, min, max *int) *ClampedValue {
	return &ClampedValue{inner: inner, min: min, max: max}
}

func NewMinValue(inner DynamicValue, min int) *ClampedValue {
	return &ClampedValue{inner: inner, min: &min, max: nil}
}

func NewMaxValue(inner DynamicValue, max int) *ClampedValue {
	return &ClampedValue{inner: inner, min: nil, max: &max}
}

func (v *ClampedValue) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	value := v.inner.Calculate(ctx, game, source)

	if v.min != nil && value < *v.min {
		return *v.min
	}
	if v.max != nil && value > *v.max {
		return *v.max
	}
	return value
}

func (v *ClampedValue) GetMessage() string {
	return v.inner.GetMessage()
}

func (v *ClampedValue) Copy() DynamicValue {
	var minCopy, maxCopy *int
	if v.min != nil {
		minVal := *v.min
		minCopy = &minVal
	}
	if v.max != nil {
		maxVal := *v.max
		maxCopy = &maxVal
	}
	return &ClampedValue{inner: v.inner.Copy(), min: minCopy, max: maxCopy}
}

// ========================================
// Signed Value (for P/T modifications)
// ========================================

// SignedCorrectionValue calculates a value ensuring it shows proper sign
// Java: mage.abilities.dynamicvalue.common.SignInversionDynamicValue
type SignedCorrectionValue struct {
	inner DynamicValue
}

func NewSignedCorrectionValue(inner DynamicValue) *SignedCorrectionValue {
	return &SignedCorrectionValue{inner: inner}
}

func (v *SignedCorrectionValue) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	return -v.inner.Calculate(ctx, game, source)
}

func (v *SignedCorrectionValue) GetMessage() string {
	return v.inner.GetMessage()
}

func (v *SignedCorrectionValue) Copy() DynamicValue {
	return &SignedCorrectionValue{inner: v.inner.Copy()}
}

// ========================================
// Permanent Count Value
// ========================================

// PermanentCount counts permanents matching a filter
// Java: mage.abilities.dynamicvalue.common.PermanentsOnBattlefieldCount
type PermanentCount struct {
	filter      TargetFilter
	controlOnly bool
}

func NewPermanentCount(filter TargetFilter, controlOnly bool) *PermanentCount {
	return &PermanentCount{filter: filter, controlOnly: controlOnly}
}

func (v *PermanentCount) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	dvGame, ok := game.(DynamicValueGameContext)
	if !ok {
		return 0
	}

	var permanents []PermanentInfo

	if v.controlOnly {
		controllerID, getControllerErr := dvGame.GetControllerID(source)
		if getControllerErr != nil {
			return 0
		}
		var getPermanentsErr error
		permanents, getPermanentsErr = dvGame.GetPermanentsControlledBy(ctx, controllerID)
		if getPermanentsErr != nil {
			return 0
		}
	} else {
		// TODO: Get all permanents
		permanents = nil
	}

	count := 0
	for _, perm := range permanents {
		if v.filter == nil || v.filter.Matches(perm.GetID(), game) {
			count++
		}
	}

	return count
}

func (v *PermanentCount) GetMessage() string {
	if v.filter != nil {
		if v.controlOnly {
			return "the number of " + v.filter.GetDescription() + "s you control"
		}
		return "the number of " + v.filter.GetDescription() + "s"
	}
	if v.controlOnly {
		return "the number of permanents you control"
	}
	return "the number of permanents"
}

func (v *PermanentCount) Copy() DynamicValue {
	return &PermanentCount{filter: v.filter, controlOnly: v.controlOnly}
}

// ========================================
// Dynamic Boost Effect
// ========================================

// DynamicBoostEffect modifies power/toughness by dynamic values
// Java: mage.abilities.effects.common.continuous.BoostSourceEffect with dynamic values
type DynamicBoostEffect struct {
	powerValue     DynamicValue
	toughnessValue DynamicValue
	duration       Duration
}

func NewDynamicBoostEffect(power, toughness DynamicValue, duration Duration) *DynamicBoostEffect {
	return &DynamicBoostEffect{
		powerValue:     power,
		toughnessValue: toughness,
		duration:       duration,
	}
}

func (e *DynamicBoostEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Apply dynamic boost as continuous effect
	return nil
}

func (e *DynamicBoostEffect) GetDescription() string {
	// Note: using math.MaxInt to avoid importing "math" in generated code that uses these
	_ = math.MaxInt
	return "gets boost based on dynamic values"
}
