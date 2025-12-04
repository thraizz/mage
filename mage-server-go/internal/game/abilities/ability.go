package abilities

import (
	"context"

	"github.com/google/uuid"
)

// Ability represents any ability on a card
// Mirrors Java's Ability interface
type Ability interface {
	// GetID returns the unique ID of this ability
	GetID() uuid.UUID

	// GetType returns the type of ability
	GetType() AbilityType

	// GetSourceID returns the ID of the card that has this ability
	GetSourceID() uuid.UUID

	// CanActivate checks if this ability can be activated/triggered
	CanActivate(ctx context.Context, game GameContext) bool

	// Resolve resolves this ability
	Resolve(ctx context.Context, game GameContext) error

	// String returns a text representation of the ability
	String() string
}

// AbilityType indicates what kind of ability this is
type AbilityType int

const (
	// AbilityTypeActivated is an activated ability (cost: effect)
	AbilityTypeActivated AbilityType = iota

	// AbilityTypeTriggered is a triggered ability (when X, do Y)
	AbilityTypeTriggered

	// AbilityTypeStatic is a static ability (continuous effect)
	AbilityTypeStatic

	// AbilityTypeSpell is a spell ability (on instants/sorceries)
	AbilityTypeSpell

	// AbilityTypeMana is a mana ability (produces mana)
	AbilityTypeMana

	// AbilityTypeKeyword is a keyword ability (flying, trample, etc.)
	AbilityTypeKeyword
)

func (t AbilityType) String() string {
	switch t {
	case AbilityTypeActivated:
		return "Activated"
	case AbilityTypeTriggered:
		return "Triggered"
	case AbilityTypeStatic:
		return "Static"
	case AbilityTypeSpell:
		return "Spell"
	case AbilityTypeMana:
		return "Mana"
	case AbilityTypeKeyword:
		return "Keyword"
	default:
		return "Unknown"
	}
}

// GameContext provides access to the game state during ability resolution
// This is a minimal interface to avoid circular dependencies with the game package
type GameContext interface {
	// GetCard retrieves a card by ID
	GetCard(id uuid.UUID) (interface{}, error)

	// GetPlayer retrieves a player by ID
	GetPlayer(id uuid.UUID) (interface{}, error)

	// DealDamage deals damage from source to target
	DealDamage(sourceID, targetID uuid.UUID, amount int) error

	// DrawCards has a player draw cards
	DrawCards(playerID uuid.UUID, amount int) error

	// DestroyPermanent destroys a permanent
	DestroyPermanent(permanentID uuid.UUID) error

	// AddMana adds mana to a player's mana pool
	AddMana(playerID uuid.UUID, mana *Mana) error

	// GetManaPool returns a player's mana pool for cost payment
	GetManaPool(playerID uuid.UUID) ManaPoolInterface

	// TapPermanent taps a permanent
	TapPermanent(permanentID uuid.UUID) error

	// UntapPermanent untaps a permanent
	UntapPermanent(permanentID uuid.UUID) error

	// IsPermanentTapped checks if a permanent is tapped
	IsPermanentTapped(permanentID uuid.UUID) bool

	// SacrificePermanent sacrifices a permanent (moves to graveyard, triggers dies events)
	SacrificePermanent(permanentID uuid.UUID) error

	// DiscardCard discards a card from a player's hand
	DiscardCard(playerID uuid.UUID, cardID uuid.UUID) error

	// GainLife has a player gain life
	GainLife(playerID uuid.UUID, amount int) error

	// GetPlayerHand returns the cards in a player's hand
	GetPlayerHand(playerID uuid.UUID) ([]interface{}, error)

	// GetPermanentsControlledByPlayer returns all permanents controlled by a player
	GetPermanentsControlledByPlayer(playerID uuid.UUID) ([]interface{}, error)

	// GetControllerID returns the controller ID of a permanent/card
	GetControllerID(objectID uuid.UUID) (uuid.UUID, error)

	// GetCardColors returns the colors of a card (e.g., ["W", "U"])
	GetCardColors(cardID uuid.UUID) ([]string, error)

	// CDA support methods
	// GetAllCardsInZone returns all cards in a specific zone (for Tarmogoyf, Lord of Extinction, etc.)
	GetAllCardsInZone(ctx context.Context, zone int) []CardInfo

	// GetCreaturesControlledBy returns all creatures controlled by a player (for "creatures you control" CDAs)
	GetCreaturesControlledBy(ctx context.Context, playerID uuid.UUID) []CardInfo

	// GetPlayerHandForCDA returns cards in a player's hand for CDA calculations (for Maro, etc.)
	// Note: This is separate from GetPlayerHand above to avoid signature conflicts
	GetPlayerHandForCDA(ctx context.Context, playerID uuid.UUID) []CardInfo

	// GetCountersOnPermanent returns the number of a specific counter type on a permanent
	GetCountersOnPermanent(ctx context.Context, permanentID uuid.UUID, counterType string) int

	// GetAllCountersOnPermanent returns all counters on a permanent as a map (name -> count)
	// Java: permanent.getCounters(game).values()
	// Used by effects like Resourceful Defense that need to move all counter types
	GetAllCountersOnPermanent(ctx context.Context, permanentID uuid.UUID) map[string]int

	// GetPermanentPower returns the power of a permanent (creature)
	// Returns 0 for non-creatures
	GetPermanentPower(ctx context.Context, permanentID uuid.UUID) int

	// GetPermanentToughness returns the toughness of a permanent (creature)
	// Returns 0 for non-creatures
	GetPermanentToughness(ctx context.Context, permanentID uuid.UUID) int

	// RemoveCountersFromPermanent removes counters from a permanent
	// Java: permanent.removeCounters(counterName, amount, source, game)
	// Returns error if the permanent doesn't exist
	RemoveCountersFromPermanent(ctx context.Context, permanentID uuid.UUID, counterName string, amount int) error

	// GetMultiAmountChoice asks the player to distribute amounts among multiple options
	// Java: player.getMultiAmountWithIndividualConstraints()
	// Used for effects like "move any number of counters" where player chooses distribution
	// choices: list of options with constraints (name, min, max, current value)
	// totalMin: minimum total that must be distributed
	// totalMax: maximum total that can be distributed
	// Returns: list of chosen amounts (one per choice), or nil if player cancelled
	GetMultiAmountChoice(
		ctx context.Context,
		playerID uuid.UUID,
		choices []MultiAmountChoice,
		totalMin, totalMax int,
		choiceType MultiAmountType,
	) ([]int, error)
}

// MultiAmountChoice represents a single option in a multi-amount choice
// Java: MultiAmountMessage
type MultiAmountChoice struct {
	Name    string // Display name (e.g., "+1/+1 (3)")
	Min     int    // Minimum value for this option
	Max     int    // Maximum value for this option
	Current int    // Current value on the permanent (for display)
}

// NewMultiAmountChoice creates a new choice option
func NewMultiAmountChoice(name string, min, max, current int) MultiAmountChoice {
	return MultiAmountChoice{
		Name:    name,
		Min:     min,
		Max:     max,
		Current: current,
	}
}

// MultiAmountType indicates what type of multi-amount choice this is
// Java: MultiAmountType enum
type MultiAmountType int

const (
	// MultiAmountTypeCounters is for choosing counter amounts
	MultiAmountTypeCounters MultiAmountType = iota
	// MultiAmountTypeDamage is for distributing damage
	MultiAmountTypeDamage
	// MultiAmountTypeMana is for choosing mana distribution
	MultiAmountTypeMana
	// MultiAmountTypeGeneric is for other distributions
	MultiAmountTypeGeneric
)

// CardInfo provides minimal card information for CDA calculations
type CardInfo interface {
	GetID() uuid.UUID
	GetName() string
	GetTypes() []string
	GetSubtypes() []string
	GetPower() int
	GetToughness() int
}

// ManaPoolInterface provides access to a player's mana pool
// This interface allows the abilities system to check and spend mana
// without depending on the concrete mana.ManaPool type
type ManaPoolInterface interface {
	// GetAmount returns the amount of a specific mana type (including floating)
	GetAmount(manaType string) int

	// SpendMana attempts to spend mana from the pool
	// Returns error if insufficient mana
	SpendMana(manaType string, amount int) error
}

// Mana represents mana in a player's mana pool
type Mana struct {
	White     int
	Blue      int
	Black     int
	Red       int
	Green     int
	Colorless int
	Generic   int // Can be paid with any color or colorless
}

// Total returns the total amount of mana
func (m *Mana) Total() int {
	return m.White + m.Blue + m.Black + m.Red + m.Green + m.Colorless + m.Generic
}

// Add adds another mana to this mana
func (m *Mana) Add(other *Mana) {
	m.White += other.White
	m.Blue += other.Blue
	m.Black += other.Black
	m.Red += other.Red
	m.Green += other.Green
	m.Colorless += other.Colorless
	m.Generic += other.Generic
}

// NewMana creates a new empty mana
func NewMana() *Mana {
	return &Mana{}
}

// baseAbility provides common fields for all abilities
type baseAbility struct {
	id       uuid.UUID
	sourceID uuid.UUID
	text     string
}

func (b *baseAbility) GetID() uuid.UUID {
	return b.id
}

func (b *baseAbility) GetSourceID() uuid.UUID {
	return b.sourceID
}

func (b *baseAbility) String() string {
	return b.text
}

// newBaseAbility creates a new base ability
func newBaseAbility(sourceID uuid.UUID, text string) baseAbility {
	return baseAbility{
		id:       uuid.New(),
		sourceID: sourceID,
		text:     text,
	}
}
