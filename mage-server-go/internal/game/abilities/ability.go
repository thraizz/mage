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

	// TODO: Add more methods as needed
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
