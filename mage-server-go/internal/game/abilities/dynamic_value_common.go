package abilities

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

// ========================================
// Devotion Count
// ========================================

// ColoredManaSymbol represents a colored mana symbol
type ColoredManaSymbol string

const (
	ManaWhite ColoredManaSymbol = "W"
	ManaBlue  ColoredManaSymbol = "U"
	ManaBlack ColoredManaSymbol = "B"
	ManaRed   ColoredManaSymbol = "R"
	ManaGreen ColoredManaSymbol = "G"
)

// DevotionCount counts colored mana symbols in mana costs of permanents you control.
// Mirrors Java's mage.abilities.dynamicvalue.common.DevotionCount.
type DevotionCount struct {
	colors []ColoredManaSymbol
}

// Pre-defined devotion count instances for common uses
var (
	DevotionW  = &DevotionCount{colors: []ColoredManaSymbol{ManaWhite}}
	DevotionU  = &DevotionCount{colors: []ColoredManaSymbol{ManaBlue}}
	DevotionB  = &DevotionCount{colors: []ColoredManaSymbol{ManaBlack}}
	DevotionR  = &DevotionCount{colors: []ColoredManaSymbol{ManaRed}}
	DevotionG  = &DevotionCount{colors: []ColoredManaSymbol{ManaGreen}}
	DevotionWU = &DevotionCount{colors: []ColoredManaSymbol{ManaWhite, ManaBlue}}
	DevotionWB = &DevotionCount{colors: []ColoredManaSymbol{ManaWhite, ManaBlack}}
	DevotionUB = &DevotionCount{colors: []ColoredManaSymbol{ManaBlue, ManaBlack}}
	DevotionUR = &DevotionCount{colors: []ColoredManaSymbol{ManaBlue, ManaRed}}
	DevotionBR = &DevotionCount{colors: []ColoredManaSymbol{ManaBlack, ManaRed}}
	DevotionBG = &DevotionCount{colors: []ColoredManaSymbol{ManaBlack, ManaGreen}}
	DevotionRG = &DevotionCount{colors: []ColoredManaSymbol{ManaRed, ManaGreen}}
	DevotionRW = &DevotionCount{colors: []ColoredManaSymbol{ManaRed, ManaWhite}}
	DevotionGW = &DevotionCount{colors: []ColoredManaSymbol{ManaGreen, ManaWhite}}
	DevotionGU = &DevotionCount{colors: []ColoredManaSymbol{ManaGreen, ManaBlue}}
)

// NewDevotionCount creates a new devotion count for the specified colors.
func NewDevotionCount(colors ...ColoredManaSymbol) *DevotionCount {
	return &DevotionCount{colors: colors}
}

func (d *DevotionCount) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	dvGame, ok := game.(DynamicValueGameContext)
	if !ok {
		return 0
	}

	controllerID, err := dvGame.GetControllerID(source)
	if err != nil {
		return 0
	}

	permanents, err := dvGame.GetPermanentsControlledBy(ctx, controllerID)
	if err != nil {
		return 0
	}

	devotion := 0
	for _, perm := range permanents {
		manaCost := perm.GetManaCost()
		devotion += d.countSymbolsInManaCost(manaCost)
	}

	return devotion
}

// countSymbolsInManaCost counts how many of the devotion colors appear in the mana cost.
func (d *DevotionCount) countSymbolsInManaCost(manaCost string) int {
	count := 0
	for _, color := range d.colors {
		// Count occurrences of the color symbol in the mana cost
		// Mana costs are typically formatted like "{2}{B}{B}" or "{W}{U}"
		count += strings.Count(manaCost, string(color))
	}
	return count
}

func (d *DevotionCount) GetMessage() string {
	if len(d.colors) == 1 {
		return "your devotion to " + colorName(d.colors[0])
	}

	colorNames := make([]string, len(d.colors))
	for i, c := range d.colors {
		colorNames[i] = colorName(c)
	}
	return "your devotion to " + strings.Join(colorNames, " and ")
}

func (d *DevotionCount) Copy() DynamicValue {
	colorsCopy := make([]ColoredManaSymbol, len(d.colors))
	copy(colorsCopy, d.colors)
	return &DevotionCount{colors: colorsCopy}
}

func colorName(c ColoredManaSymbol) string {
	switch c {
	case ManaWhite:
		return "white"
	case ManaBlue:
		return "blue"
	case ManaBlack:
		return "black"
	case ManaRed:
		return "red"
	case ManaGreen:
		return "green"
	default:
		return string(c)
	}
}

// ========================================
// Mana Spent To Cast Count
// ========================================

// ManaSpentToCastCount returns the total mana spent to cast a spell.
// Mirrors Java's mage.abilities.dynamicvalue.common.ManaSpentToCastCount.
type ManaSpentToCastCount struct{}

// ManaSpentToCastCountInstance is the singleton instance.
var ManaSpentToCastCountInstance = &ManaSpentToCastCount{}

func (m *ManaSpentToCastCount) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	dvGame, ok := game.(DynamicValueGameContext)
	if !ok {
		return 0
	}

	return dvGame.GetManaSpentToCast(ctx, source)
}

func (m *ManaSpentToCastCount) GetMessage() string {
	return "the amount of mana spent to cast it"
}

func (m *ManaSpentToCastCount) Copy() DynamicValue {
	return ManaSpentToCastCountInstance
}

// ========================================
// Counters Controller Count
// ========================================

// CountersControllerCount counts counters of a specific type on the controller.
// Mirrors Java's mage.abilities.dynamicvalue.common.CountersControllerCount.
type CountersControllerCount struct {
	counterType string
}

// NewCountersControllerCount creates a new counter count for the specified counter type.
func NewCountersControllerCount(counterType string) *CountersControllerCount {
	return &CountersControllerCount{counterType: counterType}
}

func (c *CountersControllerCount) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	dvGame, ok := game.(DynamicValueGameContext)
	if !ok {
		return 0
	}

	controllerID, err := dvGame.GetControllerID(source)
	if err != nil {
		return 0
	}

	return dvGame.GetPlayerCounters(ctx, controllerID, c.counterType)
}

func (c *CountersControllerCount) GetMessage() string {
	return "the number of " + c.counterType + " counters you have"
}

func (c *CountersControllerCount) Copy() DynamicValue {
	return &CountersControllerCount{counterType: c.counterType}
}

// ========================================
// Cards In Controller Graveyard Count
// ========================================

// GraveyardCardFilter is a function that filters cards based on some criteria.
// Used for dynamic value calculations on graveyard cards.
type GraveyardCardFilter func(card CardTypeInfo) bool

// CardsInControllerGraveyardCount counts cards in your graveyard matching a filter.
// Mirrors Java's mage.abilities.dynamicvalue.common.CardsInControllerGraveyardCount.
type CardsInControllerGraveyardCount struct {
	filter     GraveyardCardFilter
	multiplier int
	filterDesc string
}

// NewCardsInControllerGraveyardCount creates a counter for all cards in graveyard.
func NewCardsInControllerGraveyardCount() *CardsInControllerGraveyardCount {
	return &CardsInControllerGraveyardCount{
		filter:     nil, // nil means all cards
		multiplier: 1,
		filterDesc: "card",
	}
}

// NewCardsInControllerGraveyardCountFiltered creates a counter with a filter.
func NewCardsInControllerGraveyardCountFiltered(filter GraveyardCardFilter, filterDesc string) *CardsInControllerGraveyardCount {
	return &CardsInControllerGraveyardCount{
		filter:     filter,
		multiplier: 1,
		filterDesc: filterDesc,
	}
}

// NewCardsInControllerGraveyardCountArtifactOrCreature creates a counter for artifacts and creatures.
func NewCardsInControllerGraveyardCountArtifactOrCreature() *CardsInControllerGraveyardCount {
	return &CardsInControllerGraveyardCount{
		filter: func(card CardTypeInfo) bool {
			types := card.GetTypes()
			for _, t := range types {
				if t == "ARTIFACT" || t == "CREATURE" {
					return true
				}
			}
			return false
		},
		multiplier: 1,
		filterDesc: "artifact and/or creature card",
	}
}

func (c *CardsInControllerGraveyardCount) Calculate(ctx context.Context, game GameContext, source uuid.UUID) int {
	dvGame, ok := game.(DynamicValueGameContext)
	if !ok {
		return 0
	}

	controllerID, err := dvGame.GetControllerID(source)
	if err != nil {
		return 0
	}

	cards, err := dvGame.GetCardsInGraveyard(ctx, controllerID)
	if err != nil {
		return 0
	}

	count := 0
	for _, card := range cards {
		if c.filter == nil || c.filter(card) {
			count++
		}
	}

	return count * c.multiplier
}

func (c *CardsInControllerGraveyardCount) GetMessage() string {
	return c.filterDesc + " in your graveyard"
}

func (c *CardsInControllerGraveyardCount) Copy() DynamicValue {
	return &CardsInControllerGraveyardCount{
		filter:     c.filter,
		multiplier: c.multiplier,
		filterDesc: c.filterDesc,
	}
}
