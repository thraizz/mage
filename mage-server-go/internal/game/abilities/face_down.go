package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Rule 708: Face-Down Spells and Permanents
// This file implements the mechanics for face-down permanents including:
// - Morph (Rule 702.37)
// - Manifest (Rule 701.34)
// - Megamorph (Rule 702.37)
// - Cloak (Rule 702.162)

// ===== Face-Down Card State =====

// FaceDownType indicates how a card became face down
type FaceDownType int

const (
	FaceDownMorph     FaceDownType = iota // Cast face down using morph
	FaceDownManifest                      // Put onto battlefield face down via manifest
	FaceDownMegamorph                     // Cast face down using megamorph
	FaceDownCloak                         // Put onto battlefield face down via cloak
	FaceDownDisguise                      // Cast face down using disguise
)

// FaceDownState tracks the state of a face-down permanent
// Rule 708.2: Face-down spells and permanents have no characteristics except those listed
type FaceDownState struct {
	// Identity
	permanentID  uuid.UUID
	actualCard   uuid.UUID // The real card ID (hidden from opponents)
	ownerID      uuid.UUID
	controllerID uuid.UUID

	// Face-down characteristics (Rule 708.2)
	// A face-down permanent is a 2/2 creature with no text, name, subtypes,
	// mana cost, or color
	isCreature    bool // Always true for morph/manifest
	power         int  // Always 2
	toughness     int  // Always 2
	hasNoText     bool // Always true
	hasNoName     bool // Always true
	hasNoSubtypes bool // Always true
	hasNoManaCost bool // Always true
	hasNoColor    bool // Always true (colorless)

	// How it became face down
	faceDownType FaceDownType

	// Costs and options
	morphCost       *ManaCost // For morph/megamorph
	alternativeCost *ManaCost // For disguise
	canTurnFaceUp   bool      // Can use special action to turn face up

	// Megamorph tracking
	isMegamorph bool // If true, put +1/+1 counter when turned face up
}

// NewFaceDownState creates a new face-down permanent state
func NewFaceDownState(permanentID, actualCard, ownerID, controllerID uuid.UUID, faceDownType FaceDownType) *FaceDownState {
	return &FaceDownState{
		permanentID:   permanentID,
		actualCard:    actualCard,
		ownerID:       ownerID,
		controllerID:  controllerID,
		faceDownType:  faceDownType,
		isCreature:    true,
		power:         2,
		toughness:     2,
		hasNoText:     true,
		hasNoName:     true,
		hasNoSubtypes: true,
		hasNoManaCost: true,
		hasNoColor:    true,
		canTurnFaceUp: false,
	}
}

// GetPower returns the face-down power (always 2)
func (fds *FaceDownState) GetPower() int {
	return fds.power
}

// GetToughness returns the face-down toughness (always 2)
func (fds *FaceDownState) GetToughness() int {
	return fds.toughness
}

// IsFaceDown checks if a permanent is face down
func (fds *FaceDownState) IsFaceDown() bool {
	return true // This struct only exists for face-down permanents
}

// CanTurnFaceUp checks if the permanent can be turned face up
func (fds *FaceDownState) CanTurnFaceUp() bool {
	return fds.canTurnFaceUp
}

// GetMorphCost returns the cost to turn this face up
func (fds *FaceDownState) GetMorphCost() *ManaCost {
	return fds.morphCost
}

// GetFaceDownType returns how this became face down
func (fds *FaceDownState) GetFaceDownType() FaceDownType {
	return fds.faceDownType
}

// IsMegamorph returns whether this is megamorph
func (fds *FaceDownState) IsMegamorph() bool {
	return fds.isMegamorph
}

// ===== Morph Ability (Rule 702.37) =====

// MorphAbility represents the Morph keyword
// Rule 702.37a: Morph is a static ability that functions in any zone from which
// you could play the card it's on, and the morph effect works any time the card
// is face down.
type MorphAbility struct {
	baseAbility
	morphCost   *ManaCost
	isMegamorph bool // Megamorph variant
}

// NewMorphAbility creates a Morph ability
func NewMorphAbility(source uuid.UUID, cost *ManaCost) *MorphAbility {
	return &MorphAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		morphCost:   cost,
		isMegamorph: false,
	}
}

// NewMegamorphAbility creates a Megamorph ability
// Rule 702.37c: Megamorph is a variant of morph. When turned face up,
// it puts a +1/+1 counter on the creature
func NewMegamorphAbility(source uuid.UUID, cost *ManaCost) *MorphAbility {
	return &MorphAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		morphCost:   cost,
		isMegamorph: true,
	}
}

func (a *MorphAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *MorphAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true // Always active
}

func (a *MorphAbility) Resolve(ctx context.Context, game GameContext) error {
	// Morph doesn't resolve - it modifies casting and provides turn face up option
	// Rule 702.37b: You may cast a card with morph face down as a 2/2 creature
	// Rule 702.37d: You may turn a face-down permanent face up any time you have priority
	return nil
}

func (a *MorphAbility) GetMorphCost() *ManaCost {
	return a.morphCost
}

func (a *MorphAbility) IsMegamorph() bool {
	return a.isMegamorph
}

func (a *MorphAbility) String() string {
	if a.isMegamorph {
		return fmt.Sprintf("Megamorph %s", a.morphCost.String())
	}
	return fmt.Sprintf("Morph %s", a.morphCost.String())
}

// ===== Manifest Effect (Rule 701.34) =====

// ManifestEffect puts cards onto the battlefield face down
// Rule 701.34a: To manifest a card, turn it face down. It becomes a 2/2 face-down
// creature card with no text, no name, no subtypes, and no mana cost.
type ManifestEffect struct {
	description string
	source      uuid.UUID
	cardSource  ManifestSource // Where to get the card(s) from
	count       int            // Number of cards to manifest
}

// ManifestSource indicates where to manifest cards from
type ManifestSource int

const (
	ManifestFromLibrary   ManifestSource = iota // Top of library (most common)
	ManifestFromHand                            // From hand
	ManifestFromGraveyard                       // From graveyard (rare)
)

// NewManifestEffect creates a Manifest effect
func NewManifestEffect(source uuid.UUID, cardSource ManifestSource, count int) *ManifestEffect {
	return &ManifestEffect{
		description: "Manifest",
		source:      source,
		cardSource:  cardSource,
		count:       count,
	}
}

// Apply implements the Effect interface
func (e *ManifestEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. Get the card(s) from the specified zone
	// 2. For each card:
	//    a. Create FaceDownState
	//    b. Put onto battlefield face down as 2/2 creature
	//    c. If it's a creature card, it can be turned face up by paying mana cost
	//    d. If not a creature, it remains face down

	// Rule 701.34b: The face-down permanent is a 2/2 creature with no text, name,
	// subtypes, or mana cost. These values are the copiable values of that object's
	// characteristics.

	// Rule 701.34c: If a manifested creature card has morph, its controller may turn
	// that card face up using the morph rules

	// Rule 701.34d: Any time you have priority, you may turn a manifested creature face up
	// if it's a creature card by revealing it and paying its mana cost

	return fmt.Errorf("manifest effect not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *ManifestEffect) GetDescription() string {
	return fmt.Sprintf("%s %d card(s)", e.description, e.count)
}

// ===== Cloak Effect (Rule 702.162) =====

// CloakEffect is similar to Manifest but specific to Outlaws of Thunder Junction
// Rule 702.162a: Cloak [cost] means "You may cast this card face down as a 2/2
// creature spell with no text, name, creature type, or mana cost by paying {3}
// rather than paying its mana cost."
type CloakEffect struct {
	description string
	source      uuid.UUID
	cloakCost   *ManaCost
}

// NewCloakEffect creates a Cloak effect
func NewCloakEffect(source uuid.UUID, cost *ManaCost) *CloakEffect {
	return &CloakEffect{
		description: "Cloak",
		source:      source,
		cloakCost:   cost,
	}
}

// Apply implements the Effect interface
func (e *CloakEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Similar to Manifest but specifically for casting face down
	return fmt.Errorf("cloak effect not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *CloakEffect) GetDescription() string {
	return fmt.Sprintf("%s %s", e.description, e.cloakCost.String())
}

// ===== Turn Face Up Special Action (Rule 708.8) =====

// TurnFaceUpAction represents the special action to turn a face-down permanent face up
// Rule 708.8: Any time you have priority, you may turn a face-down permanent you control
// with a morph ability face up. This is a special action; it doesn't use the stack.
type TurnFaceUpAction struct {
	permanentID uuid.UUID
	player      uuid.UUID
	cost        *ManaCost
	isMegamorph bool
}

// NewTurnFaceUpAction creates a turn face up special action
func NewTurnFaceUpAction(permanentID, player uuid.UUID, cost *ManaCost, isMegamorph bool) *TurnFaceUpAction {
	return &TurnFaceUpAction{
		permanentID: permanentID,
		player:      player,
		cost:        cost,
		isMegamorph: isMegamorph,
	}
}

// CanPerform checks if the player can turn the permanent face up
func (tfa *TurnFaceUpAction) CanPerform(ctx context.Context, game GameContext) bool {
	// Requirements:
	// 1. Player has priority
	// 2. Player controls the permanent
	// 3. Permanent is face down
	// 4. Player can pay the cost
	return false // Placeholder
}

// Perform executes the turn face up action
func (tfa *TurnFaceUpAction) Perform(ctx context.Context, game GameContext) error {
	// Steps:
	// 1. Pay the cost
	// 2. Reveal the card to all players
	// 3. Turn the permanent face up (restore all characteristics)
	// 4. If megamorph, put a +1/+1 counter on it
	// 5. Trigger any "when turned face up" abilities

	// Rule 708.8a: This is a special action; it doesn't use the stack
	// Rule 708.8b: If you can't pay the cost, you can't turn the permanent face up

	return fmt.Errorf("turn face up action not yet fully implemented in game context")
}

// GetPermanentID returns the permanent being turned face up
func (tfa *TurnFaceUpAction) GetPermanentID() uuid.UUID {
	return tfa.permanentID
}

// GetCost returns the cost to turn face up
func (tfa *TurnFaceUpAction) GetCost() *ManaCost {
	return tfa.cost
}

// IsMegamorph returns whether this is megamorph
func (tfa *TurnFaceUpAction) IsMegamorph() bool {
	return tfa.isMegamorph
}

// ===== Face-Down Casting Support =====

// CastFaceDownOption represents the option to cast a spell face down
type CastFaceDownOption struct {
	cardID       uuid.UUID
	faceDownType FaceDownType
	cost         *ManaCost // Cost to cast face down (usually {3})
}

// NewCastFaceDownOption creates a face-down casting option
func NewCastFaceDownOption(cardID uuid.UUID, faceDownType FaceDownType, cost *ManaCost) *CastFaceDownOption {
	return &CastFaceDownOption{
		cardID:       cardID,
		faceDownType: faceDownType,
		cost:         cost,
	}
}

// GetCardID returns the card being cast face down
func (cfdo *CastFaceDownOption) GetCardID() uuid.UUID {
	return cfdo.cardID
}

// GetFaceDownType returns how it's being cast face down
func (cfdo *CastFaceDownOption) GetFaceDownType() FaceDownType {
	return cfdo.faceDownType
}

// GetCost returns the cost to cast face down
func (cfdo *CastFaceDownOption) GetCost() *ManaCost {
	return cfdo.cost
}

// ===== Helper Functions =====

// IsPermanentFaceDown checks if a permanent is face down
func IsPermanentFaceDown(permanentID uuid.UUID, game GameContext) bool {
	// This would check the permanent's state
	return false // Placeholder
}

// GetFaceDownState retrieves the face-down state of a permanent
func GetFaceDownState(permanentID uuid.UUID, game GameContext) (*FaceDownState, error) {
	// This would retrieve the FaceDownState from game context
	return nil, fmt.Errorf("not yet implemented")
}

// CanTurnFaceUp checks if a player can turn a permanent face up
func CanTurnFaceUp(permanentID, player uuid.UUID, game GameContext) bool {
	// Requirements:
	// 1. Permanent is face down
	// 2. Player controls it
	// 3. Player has priority
	// 4. Permanent has morph or is a manifested creature
	// 5. Player can pay the cost
	return false // Placeholder
}

// GetTurnFaceUpCost returns the cost to turn a permanent face up
func GetTurnFaceUpCost(permanentID uuid.UUID, game GameContext) (*ManaCost, error) {
	faceDownState, err := GetFaceDownState(permanentID, game)
	if err != nil {
		return nil, err
	}

	// If it has morph, return morph cost
	if faceDownState.morphCost != nil {
		return faceDownState.morphCost, nil
	}

	// If it's manifested and is a creature card, return mana cost
	// Otherwise, can't be turned face up
	return nil, fmt.Errorf("permanent cannot be turned face up")
}

// ===== Integration with Existing Systems =====

// These face-down mechanics integrate with:
// - Casting system: Option to cast face down instead of normally
// - Permanent tracking: FaceDownState tracked alongside permanents
// - Priority system: Turn face up is a special action performed at priority
// - Counter system: Megamorph adds +1/+1 counter when turned face up
// - Triggered abilities: "When turned face up" triggers

// Example usage for Morph:
// When casting a card with morph, player can choose to:
// 1. Cast normally for mana cost
// 2. Cast face down for {3} (creates FaceDownState with morphCost)
//
// Later, at any time with priority:
// 3. Pay morph cost to turn face up (special action, doesn't use stack)
// 4. If megamorph, put +1/+1 counter on it

// Example usage for Manifest:
// When a spell says "manifest the top card of your library":
// 1. Take top card
// 2. Create FaceDownState for it
// 3. Put onto battlefield face down as 2/2 creature
// 4. If it's a creature card with morph, can turn up using morph
// 5. If it's a creature card without morph, can turn up by paying mana cost
// 6. If it's not a creature card, remains face down

// ===== Common Face-Down Card Examples =====

// Example cards using Morph:
// - Willbender: "{1}{U}, Morph {1}{U}"
// - Akroma, Angel of Fury: "{5}{R}{R}{R}, Morph {3}{R}{R}{R}"
// - Birchlore Rangers: "{G}, Morph {G}"
// - Exalted Angel: "{4}{W}{W}, Morph {2}{W}{W}"

// Example cards using Megamorph:
// - Den Protector: "{1}{G}, Megamorph {1}{G}"
// - Ire Shaman: "{1}{R}, Megamorph {R}"

// Example cards using Manifest:
// - Qarsi High Priest: "Manifest the top card of your library"
// - Whisperwood Elemental: "At the beginning of your end step, manifest the top card"
// - Soul Summons: "{1}{W}, Sorcery, Manifest the top two cards of your library"

// Example cards using Cloak:
// - Concealed Courtyard: "Cloak {3}"
// - Thunder Junction cards with Cloak mechanic
