package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Rule 712: Double-Faced Cards (DFCs)
// This file implements the mechanics for transformable cards including:
// - Transform (Rule 701.28)
// - Modal Double-Faced Cards (MDFCs)
// - Transforming Double-Faced Cards (TDFCs)
// - Day/Night cycle tracking

// ===== Card Face Types =====

// CardFaceType indicates what type of double-faced card this is
type CardFaceType int

const (
	CardFaceSingle    CardFaceType = iota // Normal single-faced card
	CardFaceTransform                     // Transforming DFC (Innistrad-style)
	CardFaceModal                         // Modal DFC (Zendikar Rising-style)
	CardFaceMeld                          // Meld card (pairs of cards that become one)
	CardFaceAdventure                     // Adventure card (creature + instant/sorcery)
	CardFaceSplit                         // Split card (two spells on one card)
	CardFaceFuse                          // Fuse split card (can cast both halves)
	CardFaceFlip                          // Flip card (old-style, rotates 180°)
	CardFaceLeveler                       // Leveler (creature with level up ability)
)

// FacePosition indicates which face is currently showing
type FacePosition int

const (
	FaceFront FacePosition = iota // Front face showing
	FaceBack                      // Back face showing
)

// ===== Double-Faced Card State =====

// DFCState tracks the state of a double-faced card
type DFCState struct {
	// Card identity
	cardID       uuid.UUID
	permanentID  uuid.UUID // If on battlefield
	ownerID      uuid.UUID
	controllerID uuid.UUID

	// Face information
	faceType        CardFaceType
	currentFace     FacePosition
	frontFaceID     uuid.UUID // Card ID of front face
	backFaceID      uuid.UUID // Card ID of back face
	frontFaceName   string
	backFaceName    string
	canTransform    bool
	transformCost   *ManaCost // For cards that cost mana to transform (rare)
	transformsToBFB bool      // Back face becomes front face (Bed to Breakfast)

	// Meld information (for meld cards only)
	meldPartner  uuid.UUID // The other half of the meld pair
	meldResult   uuid.UUID // The combined card they become
	isMeldResult bool      // If true, this is the result of melding

	// Modal DFC information
	isModal       bool         // Modal DFCs are chosen during casting
	chosenFace    FacePosition // Which face was chosen to cast
	canCastEither bool         // Can cast either face from current zone

	// Transform tracking
	hasTransformedThisTurn bool
	transformCount         int // Number of times transformed (for tracking)
}

// NewDFCState creates a new double-faced card state
func NewDFCState(cardID, ownerID, controllerID uuid.UUID, faceType CardFaceType) *DFCState {
	return &DFCState{
		cardID:         cardID,
		ownerID:        ownerID,
		controllerID:   controllerID,
		faceType:       faceType,
		currentFace:    FaceFront, // Always starts front face up
		canTransform:   false,
		transformCount: 0,
	}
}

// GetCurrentFace returns which face is showing
func (dfc *DFCState) GetCurrentFace() FacePosition {
	return dfc.currentFace
}

// GetFrontFaceID returns the front face card ID
func (dfc *DFCState) GetFrontFaceID() uuid.UUID {
	return dfc.frontFaceID
}

// GetBackFaceID returns the back face card ID
func (dfc *DFCState) GetBackFaceID() uuid.UUID {
	return dfc.backFaceID
}

// GetCurrentFaceID returns the currently showing face's card ID
func (dfc *DFCState) GetCurrentFaceID() uuid.UUID {
	if dfc.currentFace == FaceFront {
		return dfc.frontFaceID
	}
	return dfc.backFaceID
}

// GetFaceType returns the type of double-faced card
func (dfc *DFCState) GetFaceType() CardFaceType {
	return dfc.faceType
}

// CanTransform checks if the card can transform
func (dfc *DFCState) CanTransform() bool {
	return dfc.canTransform
}

// HasTransformedThisTurn checks if transformed this turn
func (dfc *DFCState) HasTransformedThisTurn() bool {
	return dfc.hasTransformedThisTurn
}

// GetTransformCount returns how many times card has transformed
func (dfc *DFCState) GetTransformCount() int {
	return dfc.transformCount
}

// IsModal checks if this is a modal DFC
func (dfc *DFCState) IsModal() bool {
	return dfc.isModal
}

// GetChosenFace returns which face was chosen for modal DFC
func (dfc *DFCState) GetChosenFace() FacePosition {
	return dfc.chosenFace
}

// IsFrontFace checks if front face is showing
func (dfc *DFCState) IsFrontFace() bool {
	return dfc.currentFace == FaceFront
}

// IsBackFace checks if back face is showing
func (dfc *DFCState) IsBackFace() bool {
	return dfc.currentFace == FaceBack
}

// ===== Transform Ability =====

// TransformAbility represents the ability to transform
// Rule 701.28: Transform means turn a double-faced permanent upside down
type TransformAbility struct {
	baseAbility
	transformType TransformType
	condition     func(ctx context.Context, game GameContext) bool
}

// TransformType indicates what triggers the transform
type TransformType int

const (
	TransformManual      TransformType = iota // Activated ability
	TransformTriggered                        // Triggered ability
	TransformConditional                      // Conditional (e.g., day/night)
)

// NewTransformAbility creates a Transform ability
func NewTransformAbility(source uuid.UUID, transformType TransformType) *TransformAbility {
	return &TransformAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		transformType: transformType,
	}
}

// SetCondition sets the condition for conditional transforms
func (a *TransformAbility) SetCondition(condition func(ctx context.Context, game GameContext) bool) {
	a.condition = condition
}

func (a *TransformAbility) GetType() AbilityType {
	if a.transformType == TransformManual {
		return AbilityTypeActivated
	}
	return AbilityTypeTriggered
}

func (a *TransformAbility) CanActivate(ctx context.Context, game GameContext) bool {
	// Check if source is a transformable card
	// Check if it's on the battlefield
	// Check if condition is met (if applicable)
	if a.condition != nil {
		return a.condition(ctx, game)
	}
	return true
}

func (a *TransformAbility) Resolve(ctx context.Context, game GameContext) error {
	// Transform doesn't use the stack for most cards
	// It's a special action or triggered ability
	return nil
}

func (a *TransformAbility) String() string {
	return "Transform"
}

// ===== Transform Effect (Rule 701.28) =====

// TransformEffect performs the transform action
// Rule 701.28a: Only permanents represented by transforming double-faced cards can transform
// Rule 701.28b: Although transforming a permanent uses the word "becomes," it doesn't use the stack
type TransformEffect struct {
	description string
	source      uuid.UUID
	targetCount int // Number of targets to transform
}

// NewTransformEffect creates a Transform effect
func NewTransformEffect(source uuid.UUID) *TransformEffect {
	return &TransformEffect{
		description: "Transform",
		source:      source,
		targetCount: 0, // 0 means transform source itself
	}
}

// NewTransformTargetEffect creates a Transform effect for targets
func NewTransformTargetEffect(source uuid.UUID, targetCount int) *TransformEffect {
	return &TransformEffect{
		description: "Transform target permanent(s)",
		source:      source,
		targetCount: targetCount,
	}
}

// Apply implements the Effect interface
func (e *TransformEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. Identify the permanent(s) to transform
	// 2. For each permanent:
	//    a. Verify it's a transformable DFC
	//    b. Flip to other face (front → back or back → front)
	//    c. Update all characteristics to match new face
	//    d. Trigger any "when transformed" abilities
	//
	// Rule 701.28c: If a permanent transform while its front face is showing, it transforms
	// to its back face. If it transforms while its back face is showing, it transforms to
	// its front face.
	//
	// Rule 701.28d: If a spell or ability instructs a player to transform a permanent that
	// isn't represented by a transforming DFC, nothing happens.

	return fmt.Errorf("transform effect not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *TransformEffect) GetDescription() string {
	return e.description
}

// ===== Day/Night Cycle =====

// DayNightState represents the global day/night cycle
// Rule 702.145: Daybound and Nightbound
type DayNightState struct {
	isDay          bool
	isNight        bool
	hasBeenSet     bool // Whether day/night has been set this game
	lastChangeTurn int  // Turn number of last change
}

// NewDayNightState creates a new day/night state
func NewDayNightState() *DayNightState {
	return &DayNightState{
		isDay:      false,
		isNight:    false,
		hasBeenSet: false,
	}
}

// IsDay returns whether it's day
func (dns *DayNightState) IsDay() bool {
	return dns.isDay
}

// IsNight returns whether it's night
func (dns *DayNightState) IsNight() bool {
	return dns.isNight
}

// HasBeenSet returns whether day/night has been established
func (dns *DayNightState) HasBeenSet() bool {
	return dns.hasBeenSet
}

// SetDay sets it to day
func (dns *DayNightState) SetDay() {
	dns.isDay = true
	dns.isNight = false
	dns.hasBeenSet = true
}

// SetNight sets it to night
func (dns *DayNightState) SetNight() {
	dns.isDay = false
	dns.isNight = true
	dns.hasBeenSet = true
}

// BecomesDay transitions to day (triggers abilities)
func (dns *DayNightState) BecomesDay(turn int) {
	if !dns.isDay {
		dns.SetDay()
		dns.lastChangeTurn = turn
		// TODO: Trigger "when it becomes day" abilities
	}
}

// BecomesNight transitions to night (triggers abilities)
func (dns *DayNightState) BecomesNight(turn int) {
	if !dns.isNight {
		dns.SetNight()
		dns.lastChangeTurn = turn
		// TODO: Trigger "when it becomes night" abilities
	}
}

// DayboundAbility represents the Daybound keyword
// Rule 702.145a: Daybound appears on the front faces of some transforming DFCs
type DayboundAbility struct {
	baseAbility
}

// NewDayboundAbility creates a Daybound ability
func NewDayboundAbility(source uuid.UUID) *DayboundAbility {
	return &DayboundAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *DayboundAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *DayboundAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *DayboundAbility) Resolve(ctx context.Context, game GameContext) error {
	// Daybound doesn't resolve - it establishes day/night and triggers transform
	// Rule 702.145b: If it becomes night, transform all Daybound permanents
	return nil
}

func (a *DayboundAbility) String() string {
	return "Daybound"
}

// NightboundAbility represents the Nightbound keyword
// Rule 702.145c: Nightbound appears on the back faces of some transforming DFCs
type NightboundAbility struct {
	baseAbility
}

// NewNightboundAbility creates a Nightbound ability
func NewNightboundAbility(source uuid.UUID) *NightboundAbility {
	return &NightboundAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *NightboundAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *NightboundAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *NightboundAbility) Resolve(ctx context.Context, game GameContext) error {
	// Nightbound doesn't resolve - it prevents transform during day
	// Rule 702.145d: If it becomes day, transform all Nightbound permanents
	return nil
}

func (a *NightboundAbility) String() string {
	return "Nightbound"
}

// ===== Meld Support (Rule 713) =====

// MeldAbility represents the Meld mechanic
// Rule 713: A meld card is one member of a meld pair
type MeldAbility struct {
	baseAbility
	partnerName string    // Name of the other half
	partnerID   uuid.UUID // Card ID of the other half
	resultName  string    // Name of the combined card
	resultID    uuid.UUID // Card ID of the combined card
}

// NewMeldAbility creates a Meld ability
func NewMeldAbility(source uuid.UUID, partnerName, resultName string) *MeldAbility {
	return &MeldAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		partnerName: partnerName,
		resultName:  resultName,
	}
}

func (a *MeldAbility) GetType() AbilityType {
	return AbilityTypeTriggered // Usually triggered
}

func (a *MeldAbility) CanActivate(ctx context.Context, game GameContext) bool {
	// Check if both halves are on battlefield
	// Check if trigger condition is met
	return false // Placeholder
}

func (a *MeldAbility) Resolve(ctx context.Context, game GameContext) error {
	// Meld process:
	// 1. Exile both meld cards
	// 2. Put the meld result onto battlefield (combined single card)
	// 3. It's a single permanent with characteristics from the result card
	return fmt.Errorf("meld not yet fully implemented")
}

func (a *MeldAbility) GetPartnerName() string {
	return a.partnerName
}

func (a *MeldAbility) GetResultName() string {
	return a.resultName
}

func (a *MeldAbility) String() string {
	return fmt.Sprintf("Meld (with %s into %s)", a.partnerName, a.resultName)
}

// MeldEffect performs the meld action
type MeldEffect struct {
	description string
	source      uuid.UUID
	partner     uuid.UUID
	result      uuid.UUID
}

// NewMeldEffect creates a Meld effect
func NewMeldEffect(source, partner, result uuid.UUID) *MeldEffect {
	return &MeldEffect{
		description: "Meld",
		source:      source,
		partner:     partner,
		result:      result,
	}
}

// Apply implements the Effect interface
func (e *MeldEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Rule 713.4: To meld two cards:
	// 1. The two cards must be on the battlefield
	// 2. Exile them
	// 3. The back faces of those cards become a single permanent on the battlefield
	// 4. The new permanent has the characteristics of the back faces melded together

	return fmt.Errorf("meld effect not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *MeldEffect) GetDescription() string {
	return e.description
}

// ===== Modal DFC Support =====

// ModalDFCChoice represents the choice of which face to cast
type ModalDFCChoice struct {
	cardID          uuid.UUID
	chosenFace      FacePosition
	alternativeCost *ManaCost // If casting with alternative cost
}

// NewModalDFCChoice creates a modal DFC choice
func NewModalDFCChoice(cardID uuid.UUID, face FacePosition) *ModalDFCChoice {
	return &ModalDFCChoice{
		cardID:     cardID,
		chosenFace: face,
	}
}

// GetCardID returns the card ID
func (mdc *ModalDFCChoice) GetCardID() uuid.UUID {
	return mdc.cardID
}

// GetChosenFace returns the chosen face
func (mdc *ModalDFCChoice) GetChosenFace() FacePosition {
	return mdc.chosenFace
}

// ===== Helper Functions =====

// IsTransformable checks if a card can transform
func IsTransformable(cardID uuid.UUID, game GameContext) bool {
	// Check if card is a transforming DFC
	return false // Placeholder
}

// GetDFCState retrieves the DFC state of a card
func GetDFCState(cardID uuid.UUID, game GameContext) (*DFCState, error) {
	// Retrieve DFC state from game context
	return nil, fmt.Errorf("not yet implemented")
}

// Transform performs the transform action on a permanent
func Transform(permanentID uuid.UUID, game GameContext) error {
	// Steps:
	// 1. Get DFC state
	// 2. Verify it's transformable
	// 3. Flip to other face
	// 4. Update characteristics
	// 5. Trigger transform abilities
	return fmt.Errorf("not yet implemented")
}

// CanTransform checks if a permanent can transform
func CanTransform(permanentID uuid.UUID, game GameContext) bool {
	// Requirements:
	// 1. Permanent must be a transforming DFC
	// 2. It must be on the battlefield
	return false // Placeholder
}

// ===== Integration with Existing Systems =====

// These transform mechanics integrate with:
// - Casting system: Modal DFCs choose face during casting
// - Permanent tracking: DFCState tracked alongside permanents
// - Zone changes: DFCs always enter front face up (Rule 712.7)
// - Copy effects: Copies of DFCs copy the currently showing face only
// - Triggered abilities: "When transformed" triggers
// - Day/Night cycle: Global state tracked per game

// Example usage for Transform:
// Card: Delver of Secrets (front) / Insectile Aberration (back)
// Front: "{U} Creature - Human Wizard 1/1"
// "At the beginning of your upkeep, look at the top card of your library.
//  You may reveal that card. If an instant or sorcery card is revealed this way,
//  transform Delver of Secrets."
// Back: "Creature - Human Insect 3/2, Flying"

// Example usage for Modal DFC:
// Card: Tangled Florahedron (front) / Tangled Vale (back)
// Front: "{1}{G} Creature - Elemental 2/2"
// Back: "Land, {T}: Add {G}"
// Player chooses which face to cast when casting from hand

// Example usage for Meld:
// Card: Bruna, the Fading Light + Gisela, the Broken Blade
// → Brisela, Voice of Nightmares (combined card)

// Example usage for Day/Night:
// When Daybound permanent enters during day, it stays front face
// When night begins, transform all Daybound permanents
// When Nightbound permanent enters during night, enter on back face
// When day begins, transform all Nightbound permanents

// ===== Common Transform Card Examples =====

// Transforming DFCs:
// - Delver of Secrets / Insectile Aberration
// - Huntmaster of the Fells / Ravager of the Fells
// - Garruk Relentless / Garruk, the Veil-Cursed (planeswalker)
// - Arlinn Kord / Arlinn, Embraced by the Moon (planeswalker)

// Modal DFCs:
// - Tangled Florahedron / Tangled Vale
// - Valki, God of Lies / Tibalt, Cosmic Impostor
// - Agadeem's Awakening / Agadeem, the Undercrypt
// - Sea Gate Restoration / Sea Gate, Reborn

// Meld pairs:
// - Bruna, the Fading Light + Gisela, the Broken Blade → Brisela, Voice of Nightmares
// - Graf Rats + Midnight Scavengers → Chittering Host
// - Hanweir Battlements + Hanweir Garrison → Hanweir, the Writhing Township

// Day/Night cards:
// - Celestus Sanctifier (Daybound)
// - Gavony Dawnguard (Daybound)
// - Firmament Sage (Nightbound - back face)
