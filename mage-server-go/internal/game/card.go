package game

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

// Card represents a Magic card with its game state
// This is the public-facing Card type used by the card factory system
type Card struct {
	// Identity
	ID      uuid.UUID
	Name    string
	OwnerID uuid.UUID

	// Card characteristics
	ManaCost      string
	Types         []string // CREATURE, INSTANT, SORCERY, etc.
	Subtypes      []string // HUMAN, WIZARD, etc.
	Supertypes    []string // LEGENDARY, BASIC, etc.
	Color         string   // Color identity
	Power         string   // Power (for creatures)
	Toughness     string   // Toughness (for creatures)
	Loyalty       string   // Starting loyalty (for planeswalkers)
	RulesText     string   // Rules text
	FlavorText    string   // Flavor text
	Rarity        string   // common, uncommon, rare, mythic
	CardNumber    int      // Card number in set
	SetCode       string   // Set code
	CardClassName string   // Java class name

	// Game state
	Zone         Zone               // Current zone
	ControllerID uuid.UUID          // Current controller
	Tapped       bool               // Is tapped
	FaceDown     bool               // Is face down
	Flipped      bool               // Is flipped
	Transformed  bool               // Is transformed
	Counters     *counters.Counters // Counters on the card

	// Combat state
	Attacking         bool        // Is attacking
	Blocking          bool        // Is blocking
	AttackingWhat     uuid.UUID   // What this creature is attacking
	BlockingWhat      []uuid.UUID // What this creature is blocking
	SummoningSickness bool        // Has summoning sickness

	// Damage tracking
	Damage        int               // Damage marked on this card
	DamageSources map[uuid.UUID]int // Damage by source

	// Abilities (will be populated by card implementations)
	Abilities []interface{} // Can store abilities.Ability or other ability types
}

// Zone represents where a card is located
type Zone int

const (
	ZoneLibrary Zone = iota
	ZoneHand
	ZoneBattlefield
	ZoneGraveyard
	ZoneExile
	ZoneStack
	ZoneCommand
)

// NewCard creates a new card with basic initialization
func NewCard(ownerID uuid.UUID, name string) *Card {
	return &Card{
		ID:            uuid.New(),
		Name:          name,
		OwnerID:       ownerID,
		ControllerID:  ownerID,
		Zone:          ZoneLibrary,
		Counters:      counters.NewCounters(),
		Abilities:     make([]interface{}, 0),
		DamageSources: make(map[uuid.UUID]int),
		BlockingWhat:  make([]uuid.UUID, 0),
	}
}

// IsCreature returns true if this is a creature
func (c *Card) IsCreature() bool {
	return contains(c.Types, "CREATURE")
}

// IsInstant returns true if this is an instant
func (c *Card) IsInstant() bool {
	return contains(c.Types, "INSTANT")
}

// IsSorcery returns true if this is a sorcery
func (c *Card) IsSorcery() bool {
	return contains(c.Types, "SORCERY")
}

// IsLand returns true if this is a land
func (c *Card) IsLand() bool {
	return contains(c.Types, "LAND")
}

// IsArtifact returns true if this is an artifact
func (c *Card) IsArtifact() bool {
	return contains(c.Types, "ARTIFACT")
}

// IsEnchantment returns true if this is an enchantment
func (c *Card) IsEnchantment() bool {
	return contains(c.Types, "ENCHANTMENT")
}

// IsPlaneswalker returns true if this is a planeswalker
func (c *Card) IsPlaneswalker() bool {
	return contains(c.Types, "PLANESWALKER")
}

// IsLegendary returns true if this is legendary
func (c *Card) IsLegendary() bool {
	return contains(c.Supertypes, "LEGENDARY")
}

// IsBasic returns true if this is a basic land
func (c *Card) IsBasic() bool {
	return contains(c.Supertypes, "BASIC")
}

// AddAbility adds an ability to this card
func (c *Card) AddAbility(ability interface{}) {
	c.Abilities = append(c.Abilities, ability)
}

// GetAbilities returns all abilities of this card
func (c *Card) GetAbilities() []interface{} {
	return c.Abilities
}

// GetAbilityCount returns the number of abilities on this card
func (c *Card) GetAbilityCount() int {
	return len(c.Abilities)
}

// ToInternal converts this Card to the internal engine format
func (c *Card) ToInternal() *internalCard {
	return &internalCard{
		ID:                c.ID.String(),
		Name:              c.Name,
		DisplayName:       c.Name,
		ManaCost:          c.ManaCost,
		Type:              joinTypes(c.Types, c.Subtypes, c.Supertypes),
		SubTypes:          c.Subtypes,
		SuperTypes:        c.Supertypes,
		Color:             c.Color,
		Power:             c.Power,
		Toughness:         c.Toughness,
		Loyalty:           c.Loyalty,
		CardNumber:        c.CardNumber,
		ExpansionSet:      c.SetCode,
		Rarity:            c.Rarity,
		RulesText:         c.RulesText,
		Tapped:            c.Tapped,
		Flipped:           c.Flipped,
		Transformed:       c.Transformed,
		FaceDown:          c.FaceDown,
		Zone:              int(c.Zone),
		ControllerID:      c.ControllerID.String(),
		OwnerID:           c.OwnerID.String(),
		Counters:          c.Counters,
		Attacking:         c.Attacking,
		Blocking:          c.Blocking,
		AttackingWhat:     c.AttackingWhat.String(),
		BlockingWhat:      uuidSliceToStringSlice(c.BlockingWhat),
		Damage:            c.Damage,
		DamageSources:     uuidMapToStringMap(c.DamageSources),
		SummoningSickness: c.SummoningSickness,
		// TODO: Convert abilities
		Abilities: []EngineAbilityView{},
	}
}

// Helper functions

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func joinTypes(types, subtypes, supertypes []string) string {
	result := ""

	if len(supertypes) > 0 {
		for i, st := range supertypes {
			if i > 0 {
				result += " "
			}
			result += st
		}
		result += " "
	}

	for i, t := range types {
		if i > 0 {
			result += " "
		}
		result += t
	}

	if len(subtypes) > 0 {
		result += " — "
		for i, st := range subtypes {
			if i > 0 {
				result += " "
			}
			result += st
		}
	}

	return result
}

func uuidSliceToStringSlice(uuids []uuid.UUID) []string {
	result := make([]string, len(uuids))
	for i, u := range uuids {
		result[i] = u.String()
	}
	return result
}

func uuidMapToStringMap(m map[uuid.UUID]int) map[string]int {
	result := make(map[string]int, len(m))
	for k, v := range m {
		result[k.String()] = v
	}
	return result
}
