package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CloneManaCost creates a deep copy of a ManaCost
func CloneManaCost(cost *ManaCost) *ManaCost {
	if cost == nil || cost.Mana == nil {
		return &ManaCost{Mana: NewMana()}
	}
	return &ManaCost{
		Mana: &Mana{
			White:     cost.Mana.White,
			Blue:      cost.Mana.Blue,
			Black:     cost.Mana.Black,
			Red:       cost.Mana.Red,
			Green:     cost.Mana.Green,
			Colorless: cost.Mana.Colorless,
			Generic:   cost.Mana.Generic,
		},
	}
}

// CostModifier represents a modifier that affects the cost of spells or abilities
// MTG Rule 601.2f: The total cost is determined by starting with the mana cost,
// adding all additional costs, adding cost increases, applying cost reductions
type CostModifier interface {
	// ModifyCost modifies a cost and returns the modified cost
	ModifyCost(ctx context.Context, game GameContext, cost *ManaCost) (*ManaCost, error)

	// GetType returns the type of modification
	GetType() CostModifierType

	// GetPriority returns the order in which this modifier should be applied
	// Lower numbers are applied first
	GetPriority() int

	// Applies checks if this modifier applies to the given spell/ability
	Applies(ctx context.Context, game GameContext, source uuid.UUID, cardTypes []string) bool

	// String returns a human-readable description
	String() string
}

// CostModifierType indicates what kind of cost modification this is
type CostModifierType int

const (
	CostModifierIncrease    CostModifierType = iota // Cost increases (Thalia, Sphere of Resistance)
	CostModifierReduction                           // Cost reductions (Goblin Electromancer, Affinity)
	CostModifierAlternative                         // Alternative costs (Flashback, Evoke, Force of Will)
	CostModifierAdditional                          // Additional costs (Kicker, Buyback, Replicate)
)

// CostModificationManager manages all cost modifiers in a game
type CostModificationManager struct {
	modifiers []CostModifier
}

// NewCostModificationManager creates a new cost modification manager
func NewCostModificationManager() *CostModificationManager {
	return &CostModificationManager{
		modifiers: make([]CostModifier, 0),
	}
}

// RegisterModifier registers a cost modifier
func (cmm *CostModificationManager) RegisterModifier(modifier CostModifier) {
	cmm.modifiers = append(cmm.modifiers, modifier)
}

// UnregisterModifier removes a cost modifier
func (cmm *CostModificationManager) UnregisterModifier(modifier CostModifier) {
	for i, m := range cmm.modifiers {
		if m == modifier {
			cmm.modifiers = append(cmm.modifiers[:i], cmm.modifiers[i+1:]...)
			return
		}
	}
}

// CalculateTotalCost calculates the total cost of a spell or ability
// Rule 601.2f: Total cost = mana cost + additional costs + cost increases - cost reductions
func (cmm *CostModificationManager) CalculateTotalCost(
	ctx context.Context,
	game GameContext,
	baseCost *ManaCost,
	source uuid.UUID,
	cardTypes []string,
) (*ManaCost, error) {
	if baseCost == nil {
		return nil, fmt.Errorf("base cost cannot be nil")
	}

	// Start with base cost
	totalCost := CloneManaCost(baseCost)

	// Apply modifiers in priority order
	// Rule 601.2f: Additional costs → Increases → Reductions
	applicableModifiers := make([]CostModifier, 0)
	for _, modifier := range cmm.modifiers {
		if modifier.Applies(ctx, game, source, cardTypes) {
			applicableModifiers = append(applicableModifiers, modifier)
		}
	}

	// Sort by priority (lower priority first)
	sortModifiersByPriority(applicableModifiers)

	// Apply each modifier
	for _, modifier := range applicableModifiers {
		modified, err := modifier.ModifyCost(ctx, game, totalCost)
		if err != nil {
			return nil, fmt.Errorf("failed to apply cost modifier: %w", err)
		}
		totalCost = modified
	}

	return totalCost, nil
}

// sortModifiersByPriority sorts modifiers by their priority
func sortModifiersByPriority(modifiers []CostModifier) {
	// Bubble sort (simple, but fine for small lists)
	n := len(modifiers)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if modifiers[j].GetPriority() > modifiers[j+1].GetPriority() {
				modifiers[j], modifiers[j+1] = modifiers[j+1], modifiers[j]
			}
		}
	}
}

// ===== Cost Increase Modifiers =====

// GenericCostIncreaseModifier increases the generic mana cost by a fixed amount
// Examples: Thalia, Guardian of Thraben (+{1}), Sphere of Resistance (+{1})
type GenericCostIncreaseModifier struct {
	source         uuid.UUID
	increaseAmount int
	appliesToCard  func(cardTypes []string) bool
	description    string
}

// NewGenericCostIncreaseModifier creates a cost increase modifier
func NewGenericCostIncreaseModifier(
	source uuid.UUID,
	amount int,
	appliesToCard func(cardTypes []string) bool,
	description string,
) *GenericCostIncreaseModifier {
	return &GenericCostIncreaseModifier{
		source:         source,
		increaseAmount: amount,
		appliesToCard:  appliesToCard,
		description:    description,
	}
}

func (m *GenericCostIncreaseModifier) ModifyCost(ctx context.Context, game GameContext, cost *ManaCost) (*ManaCost, error) {
	modified := CloneManaCost(cost)
	modified.Mana.Generic += m.increaseAmount
	return modified, nil
}

func (m *GenericCostIncreaseModifier) GetType() CostModifierType {
	return CostModifierIncrease
}

func (m *GenericCostIncreaseModifier) GetPriority() int {
	return 100 // Cost increases before reductions
}

func (m *GenericCostIncreaseModifier) Applies(ctx context.Context, game GameContext, source uuid.UUID, cardTypes []string) bool {
	// Don't apply to itself
	if source == m.source {
		return false
	}
	return m.appliesToCard(cardTypes)
}

func (m *GenericCostIncreaseModifier) String() string {
	return m.description
}

// ===== Cost Reduction Modifiers =====

// GenericCostReductionModifier reduces the generic mana cost by a fixed amount
// Examples: Goblin Electromancer (-{1} for instants/sorceries), Ruby Medallion (-{1} for red spells)
type GenericCostReductionModifier struct {
	source          uuid.UUID
	reductionAmount int
	appliesToCard   func(cardTypes []string) bool
	description     string
}

// NewGenericCostReductionModifier creates a cost reduction modifier
func NewGenericCostReductionModifier(
	source uuid.UUID,
	amount int,
	appliesToCard func(cardTypes []string) bool,
	description string,
) *GenericCostReductionModifier {
	return &GenericCostReductionModifier{
		source:          source,
		reductionAmount: amount,
		appliesToCard:   appliesToCard,
		description:     description,
	}
}

func (m *GenericCostReductionModifier) ModifyCost(ctx context.Context, game GameContext, cost *ManaCost) (*ManaCost, error) {
	modified := CloneManaCost(cost)
	modified.Mana.Generic -= m.reductionAmount
	if modified.Mana.Generic < 0 {
		modified.Mana.Generic = 0
	}
	return modified, nil
}

func (m *GenericCostReductionModifier) GetType() CostModifierType {
	return CostModifierReduction
}

func (m *GenericCostReductionModifier) GetPriority() int {
	return 200 // Cost reductions after increases
}

func (m *GenericCostReductionModifier) Applies(ctx context.Context, game GameContext, source uuid.UUID, cardTypes []string) bool {
	// Don't apply to itself
	if source == m.source {
		return false
	}
	return m.appliesToCard(cardTypes)
}

func (m *GenericCostReductionModifier) String() string {
	return m.description
}

// ===== Alternative Cost Modifiers =====

// AlternativeCostModifier represents an alternative way to cast a spell
// Examples: Flashback, Evoke, Force of Will, Delve
type AlternativeCostModifier struct {
	name        string
	cost        *ManaCost
	condition   func(ctx context.Context, game GameContext) bool
	description string
}

// NewAlternativeCostModifier creates an alternative cost modifier
func NewAlternativeCostModifier(
	name string,
	cost *ManaCost,
	condition func(ctx context.Context, game GameContext) bool,
	description string,
) *AlternativeCostModifier {
	return &AlternativeCostModifier{
		name:        name,
		cost:        cost,
		condition:   condition,
		description: description,
	}
}

func (m *AlternativeCostModifier) ModifyCost(ctx context.Context, game GameContext, cost *ManaCost) (*ManaCost, error) {
	// Alternative costs replace the original cost entirely
	return CloneManaCost(m.cost), nil
}

func (m *AlternativeCostModifier) GetType() CostModifierType {
	return CostModifierAlternative
}

func (m *AlternativeCostModifier) GetPriority() int {
	return 0 // Alternative costs applied first (replaces base cost)
}

func (m *AlternativeCostModifier) Applies(ctx context.Context, game GameContext, source uuid.UUID, cardTypes []string) bool {
	if m.condition == nil {
		return true
	}
	return m.condition(ctx, game)
}

func (m *AlternativeCostModifier) String() string {
	return m.description
}

func (m *AlternativeCostModifier) GetName() string {
	return m.name
}

// ===== Additional Cost Modifiers =====

// AdditionalCostModifier represents an additional cost to cast a spell
// Examples: Kicker, Buyback, Replicate, Multikicker
type AdditionalCostModifier struct {
	name        string
	cost        *ManaCost
	optional    bool
	repeatable  bool // Can be paid multiple times (Multikicker, Replicate)
	description string
}

// NewAdditionalCostModifier creates an additional cost modifier
func NewAdditionalCostModifier(
	name string,
	cost *ManaCost,
	optional bool,
	repeatable bool,
	description string,
) *AdditionalCostModifier {
	return &AdditionalCostModifier{
		name:        name,
		cost:        cost,
		optional:    optional,
		repeatable:  repeatable,
		description: description,
	}
}

func (m *AdditionalCostModifier) ModifyCost(ctx context.Context, game GameContext, cost *ManaCost) (*ManaCost, error) {
	// Additional costs are added to the total cost
	modified := CloneManaCost(cost)
	modified.Mana.Generic += m.cost.Mana.Generic
	modified.Mana.White += m.cost.Mana.White
	modified.Mana.Blue += m.cost.Mana.Blue
	modified.Mana.Black += m.cost.Mana.Black
	modified.Mana.Red += m.cost.Mana.Red
	modified.Mana.Green += m.cost.Mana.Green
	modified.Mana.Colorless += m.cost.Mana.Colorless
	return modified, nil
}

func (m *AdditionalCostModifier) GetType() CostModifierType {
	return CostModifierAdditional
}

func (m *AdditionalCostModifier) GetPriority() int {
	return 50 // Additional costs before increases
}

func (m *AdditionalCostModifier) Applies(ctx context.Context, game GameContext, source uuid.UUID, cardTypes []string) bool {
	// Additional costs need to be explicitly chosen by the player
	// This is tracked separately in the casting process
	return false // By default, don't auto-apply
}

func (m *AdditionalCostModifier) String() string {
	return m.description
}

func (m *AdditionalCostModifier) GetName() string {
	return m.name
}

func (m *AdditionalCostModifier) IsOptional() bool {
	return m.optional
}

func (m *AdditionalCostModifier) IsRepeatable() bool {
	return m.repeatable
}

// ===== Affinity Modifier =====

// AffinityCostModifier represents Affinity cost reduction
// Examples: Affinity for artifacts, Affinity for Islands
type AffinityCostModifier struct {
	source       uuid.UUID
	perPermanent int
	filter       func(card interface{}) bool
	description  string
}

// NewAffinityCostModifier creates an Affinity cost modifier
func NewAffinityCostModifier(
	source uuid.UUID,
	reductionPerPermanent int,
	filter func(card interface{}) bool,
	description string,
) *AffinityCostModifier {
	return &AffinityCostModifier{
		source:       source,
		perPermanent: reductionPerPermanent,
		filter:       filter,
		description:  description,
	}
}

func (m *AffinityCostModifier) ModifyCost(ctx context.Context, game GameContext, cost *ManaCost) (*ManaCost, error) {
	// Count matching permanents
	// TODO: Get permanents from game context and count matches
	count := 0 // Placeholder

	reduction := count * m.perPermanent

	modified := CloneManaCost(cost)
	modified.Mana.Generic -= reduction
	if modified.Mana.Generic < 0 {
		modified.Mana.Generic = 0
	}
	return modified, nil
}

func (m *AffinityCostModifier) GetType() CostModifierType {
	return CostModifierReduction
}

func (m *AffinityCostModifier) GetPriority() int {
	return 200 // Cost reductions
}

func (m *AffinityCostModifier) Applies(ctx context.Context, game GameContext, source uuid.UUID, cardTypes []string) bool {
	// Don't apply to itself
	if source == m.source {
		return false
	}
	return true
}

func (m *AffinityCostModifier) String() string {
	return m.description
}

// ===== Convoke Modifier =====

// ConvokeCostModifier represents Convoke cost reduction
// Rule 702.51: Each creature you tap while casting this spell pays for {1} or one mana of that creature's color
type ConvokeCostModifier struct {
	source          uuid.UUID
	tappedCreatures []uuid.UUID // Creatures tapped for convoke
}

// NewConvokeCostModifier creates a Convoke modifier
func NewConvokeCostModifier(source uuid.UUID) *ConvokeCostModifier {
	return &ConvokeCostModifier{
		source:          source,
		tappedCreatures: make([]uuid.UUID, 0),
	}
}

func (m *ConvokeCostModifier) SetTappedCreatures(creatures []uuid.UUID) {
	m.tappedCreatures = creatures
}

func (m *ConvokeCostModifier) ModifyCost(ctx context.Context, game GameContext, cost *ManaCost) (*ManaCost, error) {
	// Each tapped creature reduces generic cost by 1
	// (Colored mana reduction requires knowing creature colors)
	reduction := len(m.tappedCreatures)

	modified := CloneManaCost(cost)
	modified.Mana.Generic -= reduction
	if modified.Mana.Generic < 0 {
		modified.Mana.Generic = 0
	}
	return modified, nil
}

func (m *ConvokeCostModifier) GetType() CostModifierType {
	return CostModifierReduction
}

func (m *ConvokeCostModifier) GetPriority() int {
	return 200
}

func (m *ConvokeCostModifier) Applies(ctx context.Context, game GameContext, source uuid.UUID, cardTypes []string) bool {
	return source == m.source && len(m.tappedCreatures) > 0
}

func (m *ConvokeCostModifier) String() string {
	return "Convoke"
}

// ===== X Cost Support =====

// XCost represents a variable cost in a spell (e.g., "Fireball costs {X}{R}")
type XCost struct {
	Value int // The chosen value for X
}

// NewXCost creates an X cost with the chosen value
func NewXCost(value int) *XCost {
	return &XCost{Value: value}
}

// AddToManaCost adds the X cost to a mana cost
func (x *XCost) AddToManaCost(cost *ManaCost) *ManaCost {
	modified := CloneManaCost(cost)
	modified.Mana.Generic += x.Value
	return modified
}

// String returns the X value
func (x *XCost) String() string {
	return fmt.Sprintf("X=%d", x.Value)
}
