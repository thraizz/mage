package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CharacteristicDefiningAbility represents abilities that define characteristics of an object
// Rule 604.3: Characteristic-defining abilities function in all zones
// They're a subset of static abilities that define characteristics like power, toughness, color, etc.
//
// Examples:
// - Tarmogoyf: "Tarmogoyf's power is equal to the number of card types among cards in all graveyards..."
// - Lord of Extinction: "Lord of Extinction's power and toughness are each equal to the number of cards in all graveyards."
// - Kavu Chameleon: "Kavu Chameleon's power and toughness are each equal to the number of creatures you control."
type CharacteristicDefiningAbility interface {
	Ability

	// DefinesPower returns true if this CDA defines power
	DefinesPower() bool

	// DefinesToughness returns true if this CDA defines toughness
	DefinesToughness() bool

	// DefinesColor returns true if this CDA defines color
	DefinesColor() bool

	// DefinesTypes returns true if this CDA defines types
	DefinesTypes() bool

	// CalculatePower calculates the power value (if DefinesPower is true)
	CalculatePower(ctx context.Context, game GameContext) (int, error)

	// CalculateToughness calculates the toughness value (if DefinesToughness is true)
	CalculateToughness(ctx context.Context, game GameContext) (int, error)

	// CalculateColor calculates the color (if DefinesColor is true)
	CalculateColor(ctx context.Context, game GameContext) ([]string, error)

	// CalculateTypes calculates the types (if DefinesTypes is true)
	CalculateTypes(ctx context.Context, game GameContext) ([]string, error)
}

// BaseCDA provides common functionality for CDAs
type BaseCDA struct {
	id           uuid.UUID
	source       uuid.UUID
	definesPower bool
	definesTough bool
	definesColor bool
	definesTypes bool
	description  string
}

// NewBaseCDA creates a new base CDA
func NewBaseCDA(source uuid.UUID, definesPower, definesTough, definesColor, definesTypes bool, description string) *BaseCDA {
	return &BaseCDA{
		id:           uuid.New(),
		source:       source,
		definesPower: definesPower,
		definesTough: definesTough,
		definesColor: definesColor,
		definesTypes: definesTypes,
		description:  description,
	}
}

// Implement Ability interface
func (c *BaseCDA) GetID() uuid.UUID       { return c.id }
func (c *BaseCDA) GetType() AbilityType   { return AbilityTypeStatic }
func (c *BaseCDA) GetSource() uuid.UUID   { return c.source }
func (c *BaseCDA) GetSourceID() uuid.UUID { return c.source }
func (c *BaseCDA) String() string         { return c.description }

func (c *BaseCDA) CanActivate(ctx context.Context, game GameContext) bool {
	return false // CDAs don't activate
}

func (c *BaseCDA) Resolve(ctx context.Context, game GameContext) error {
	return nil // CDAs don't resolve
}

// Implement CharacteristicDefiningAbility interface
func (c *BaseCDA) DefinesPower() bool     { return c.definesPower }
func (c *BaseCDA) DefinesToughness() bool { return c.definesTough }
func (c *BaseCDA) DefinesColor() bool     { return c.definesColor }
func (c *BaseCDA) DefinesTypes() bool     { return c.definesTypes }

// Default implementations (should be overridden)
func (c *BaseCDA) CalculatePower(ctx context.Context, game GameContext) (int, error) {
	return 0, fmt.Errorf("CalculatePower not implemented for this CDA")
}

func (c *BaseCDA) CalculateToughness(ctx context.Context, game GameContext) (int, error) {
	return 0, fmt.Errorf("CalculateToughness not implemented for this CDA")
}

func (c *BaseCDA) CalculateColor(ctx context.Context, game GameContext) ([]string, error) {
	return nil, fmt.Errorf("CalculateColor not implemented for this CDA")
}

func (c *BaseCDA) CalculateTypes(ctx context.Context, game GameContext) ([]string, error) {
	return nil, fmt.Errorf("CalculateTypes not implemented for this CDA")
}

// ===== Tarmogoyf CDA =====

// TarmogoyfCDA implements Tarmogoyf's characteristic-defining ability
// "Tarmogoyf's power is equal to the number of card types among cards in all graveyards
// and its toughness is equal to that number plus 1."
type TarmogoyfCDA struct {
	*BaseCDA
}

// NewTarmogoyfCDA creates a Tarmogoyf CDA
func NewTarmogoyfCDA(source uuid.UUID) *TarmogoyfCDA {
	return &TarmogoyfCDA{
		BaseCDA: NewBaseCDA(
			source,
			true,  // defines power
			true,  // defines toughness
			false, // doesn't define color
			false, // doesn't define types
			"Tarmogoyf's power is equal to the number of card types among cards in all graveyards and its toughness is equal to that number plus 1.",
		),
	}
}

func (t *TarmogoyfCDA) CalculatePower(ctx context.Context, game GameContext) (int, error) {
	count := t.countCardTypesInGraveyards(ctx, game)
	return count, nil
}

func (t *TarmogoyfCDA) CalculateToughness(ctx context.Context, game GameContext) (int, error) {
	count := t.countCardTypesInGraveyards(ctx, game)
	return count + 1, nil
}

func (t *TarmogoyfCDA) countCardTypesInGraveyards(ctx context.Context, game GameContext) int {
	// Card types in Magic: Artifact, Creature, Enchantment, Instant, Land, Planeswalker, Sorcery, Tribal/Kindred
	// Also Battle (newer)
	typesFound := make(map[string]bool)

	// Get all cards in all graveyards
	graveyardCards := game.GetAllCardsInZone(ctx, int(ZoneGraveyard))

	for _, card := range graveyardCards {
		// Check each type
		cardTypes := card.GetTypes()
		for _, cardType := range cardTypes {
			typesFound[cardType] = true
		}
	}

	return len(typesFound)
}

// ===== Lord of Extinction CDA =====

// LordOfExtinctionCDA implements Lord of Extinction's characteristic-defining ability
// "Lord of Extinction's power and toughness are each equal to the number of cards in all graveyards."
type LordOfExtinctionCDA struct {
	*BaseCDA
}

// NewLordOfExtinctionCDA creates a Lord of Extinction CDA
func NewLordOfExtinctionCDA(source uuid.UUID) *LordOfExtinctionCDA {
	return &LordOfExtinctionCDA{
		BaseCDA: NewBaseCDA(
			source,
			true,  // defines power
			true,  // defines toughness
			false, // doesn't define color
			false, // doesn't define types
			"Lord of Extinction's power and toughness are each equal to the number of cards in all graveyards.",
		),
	}
}

func (l *LordOfExtinctionCDA) CalculatePower(ctx context.Context, game GameContext) (int, error) {
	count := l.countCardsInGraveyards(ctx, game)
	return count, nil
}

func (l *LordOfExtinctionCDA) CalculateToughness(ctx context.Context, game GameContext) (int, error) {
	count := l.countCardsInGraveyards(ctx, game)
	return count, nil
}

func (l *LordOfExtinctionCDA) countCardsInGraveyards(ctx context.Context, game GameContext) int {
	graveyardCards := game.GetAllCardsInZone(ctx, int(ZoneGraveyard))
	return len(graveyardCards)
}

// ===== Creatures You Control CDA =====

// CreaturesYouControlCDA implements abilities like "* / * where * is creatures you control"
// Examples: Kavu Chameleon, Multani, Yavimaya's Avatar
type CreaturesYouControlCDA struct {
	*BaseCDA
	controllerID uuid.UUID
}

// NewCreaturesYouControlCDA creates a CDA based on creatures you control
func NewCreaturesYouControlCDA(source, controllerID uuid.UUID) *CreaturesYouControlCDA {
	return &CreaturesYouControlCDA{
		BaseCDA: NewBaseCDA(
			source,
			true,  // defines power
			true,  // defines toughness
			false, // doesn't define color
			false, // doesn't define types
			"This creature's power and toughness are each equal to the number of creatures you control.",
		),
		controllerID: controllerID,
	}
}

func (c *CreaturesYouControlCDA) CalculatePower(ctx context.Context, game GameContext) (int, error) {
	count := c.countCreaturesControlled(ctx, game)
	return count, nil
}

func (c *CreaturesYouControlCDA) CalculateToughness(ctx context.Context, game GameContext) (int, error) {
	count := c.countCreaturesControlled(ctx, game)
	return count, nil
}

func (c *CreaturesYouControlCDA) countCreaturesControlled(ctx context.Context, game GameContext) int {
	creatures := game.GetCreaturesControlledBy(ctx, c.controllerID)
	return len(creatures)
}

// ===== Hands You Control CDA =====

// HandSizeCDA implements abilities like "* / * where * is cards in your hand"
// Examples: Maro, Masumaro, First to Live
type HandSizeCDA struct {
	*BaseCDA
	controllerID uuid.UUID
}

// NewHandSizeCDA creates a CDA based on hand size
func NewHandSizeCDA(source, controllerID uuid.UUID) *HandSizeCDA {
	return &HandSizeCDA{
		BaseCDA: NewBaseCDA(
			source,
			true,  // defines power
			true,  // defines toughness
			false, // doesn't define color
			false, // doesn't define types
			"This creature's power and toughness are each equal to the number of cards in your hand.",
		),
		controllerID: controllerID,
	}
}

func (h *HandSizeCDA) CalculatePower(ctx context.Context, game GameContext) (int, error) {
	count := h.countCardsInHand(ctx, game)
	return count, nil
}

func (h *HandSizeCDA) CalculateToughness(ctx context.Context, game GameContext) (int, error) {
	count := h.countCardsInHand(ctx, game)
	return count, nil
}

func (h *HandSizeCDA) countCardsInHand(ctx context.Context, game GameContext) int {
	hand := game.GetPlayerHandForCDA(ctx, h.controllerID)
	return len(hand)
}

// ===== Generic Counter CDA =====

// CountersCDA implements abilities where P/T equals counters of a certain type
// Examples: Primordial Hydra (X/X where X is +1/+1 counters)
type CountersCDA struct {
	*BaseCDA
	counterType string
	powerOnly   bool
}

// NewCountersCDA creates a CDA based on counter count
func NewCountersCDA(source uuid.UUID, counterType string, powerOnly bool) *CountersCDA {
	description := fmt.Sprintf("This permanent's power and toughness are each equal to the number of %s counters on it.", counterType)
	if powerOnly {
		description = fmt.Sprintf("This permanent's power is equal to the number of %s counters on it.", counterType)
	}

	return &CountersCDA{
		BaseCDA: NewBaseCDA(
			source,
			true,       // defines power
			!powerOnly, // defines toughness unless power-only
			false,      // doesn't define color
			false,      // doesn't define types
			description,
		),
		counterType: counterType,
		powerOnly:   powerOnly,
	}
}

func (c *CountersCDA) CalculatePower(ctx context.Context, game GameContext) (int, error) {
	count := game.GetCountersOnPermanent(ctx, c.source, c.counterType)
	return count, nil
}

func (c *CountersCDA) CalculateToughness(ctx context.Context, game GameContext) (int, error) {
	if c.powerOnly {
		return 0, fmt.Errorf("this CDA doesn't define toughness")
	}
	count := game.GetCountersOnPermanent(ctx, c.source, c.counterType)
	return count, nil
}

// ===== Generic Calculation CDA =====

// CalculationCDA is a generic CDA that uses provided calculation functions
type CalculationCDA struct {
	*BaseCDA
	powerCalc     func(ctx context.Context, game GameContext) (int, error)
	toughnessCalc func(ctx context.Context, game GameContext) (int, error)
	colorCalc     func(ctx context.Context, game GameContext) ([]string, error)
	typesCalc     func(ctx context.Context, game GameContext) ([]string, error)
}

// NewCalculationCDA creates a generic CDA with custom calculation functions
func NewCalculationCDA(
	source uuid.UUID,
	definesPower, definesTough, definesColor, definesTypes bool,
	description string,
	powerCalc func(ctx context.Context, game GameContext) (int, error),
	toughnessCalc func(ctx context.Context, game GameContext) (int, error),
	colorCalc func(ctx context.Context, game GameContext) ([]string, error),
	typesCalc func(ctx context.Context, game GameContext) ([]string, error),
) *CalculationCDA {
	return &CalculationCDA{
		BaseCDA:       NewBaseCDA(source, definesPower, definesTough, definesColor, definesTypes, description),
		powerCalc:     powerCalc,
		toughnessCalc: toughnessCalc,
		colorCalc:     colorCalc,
		typesCalc:     typesCalc,
	}
}

func (c *CalculationCDA) CalculatePower(ctx context.Context, game GameContext) (int, error) {
	if c.powerCalc != nil {
		return c.powerCalc(ctx, game)
	}
	return c.BaseCDA.CalculatePower(ctx, game)
}

func (c *CalculationCDA) CalculateToughness(ctx context.Context, game GameContext) (int, error) {
	if c.toughnessCalc != nil {
		return c.toughnessCalc(ctx, game)
	}
	return c.BaseCDA.CalculateToughness(ctx, game)
}

func (c *CalculationCDA) CalculateColor(ctx context.Context, game GameContext) ([]string, error) {
	if c.colorCalc != nil {
		return c.colorCalc(ctx, game)
	}
	return c.BaseCDA.CalculateColor(ctx, game)
}

func (c *CalculationCDA) CalculateTypes(ctx context.Context, game GameContext) ([]string, error) {
	if c.typesCalc != nil {
		return c.typesCalc(ctx, game)
	}
	return c.BaseCDA.CalculateTypes(ctx, game)
}
