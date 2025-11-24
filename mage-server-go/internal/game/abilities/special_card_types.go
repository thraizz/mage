package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

// This file implements special card types and layouts
// These cards have unique rules and characteristics that differ from normal cards

// ===== Split Cards (Rule 709) =====

// SplitCardState tracks the state of a split card
// Rule 709: Split cards have two card faces on a single card
type SplitCardState struct {
	cardID      uuid.UUID
	leftHalfID  uuid.UUID // ID of left spell half
	rightHalfID uuid.UUID // ID of right spell half
	leftName    string
	rightName   string
	canFuse     bool // If true, this is a Fuse split card
}

// NewSplitCardState creates a new split card state
func NewSplitCardState(cardID uuid.UUID, leftName, rightName string, canFuse bool) *SplitCardState {
	return &SplitCardState{
		cardID:    cardID,
		leftName:  leftName,
		rightName: rightName,
		canFuse:   canFuse,
	}
}

// GetLeftHalfID returns the left half spell ID
func (scs *SplitCardState) GetLeftHalfID() uuid.UUID {
	return scs.leftHalfID
}

// GetRightHalfID returns the right half spell ID
func (scs *SplitCardState) GetRightHalfID() uuid.UUID {
	return scs.rightHalfID
}

// CanFuse returns whether this split card has Fuse
func (scs *SplitCardState) CanFuse() bool {
	return scs.canFuse
}

// SplitCardCastChoice represents the choice of which half(ves) to cast
type SplitCardCastChoice struct {
	cardID      uuid.UUID
	castLeft    bool
	castRight   bool
	castingBoth bool // For Fuse
}

// NewSplitCardCastChoice creates a cast choice for a split card
func NewSplitCardCastChoice(cardID uuid.UUID, left, right bool) *SplitCardCastChoice {
	return &SplitCardCastChoice{
		cardID:      cardID,
		castLeft:    left,
		castRight:   right,
		castingBoth: left && right,
	}
}

// FuseAbility represents the Fuse keyword
// Rule 702.101a: Fuse is a static ability found on some split cards that applies
// while the card is in a player's hand. If a player casts a split card with fuse
// from their hand, the player may choose to cast both halves of that split card.
type FuseAbility struct {
	baseAbility
}

// NewFuseAbility creates a Fuse ability
func NewFuseAbility(source uuid.UUID) *FuseAbility {
	return &FuseAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *FuseAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *FuseAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *FuseAbility) Resolve(ctx context.Context, game GameContext) error {
	// Fuse doesn't resolve - it modifies casting options
	// Allows casting both halves as a single spell
	return nil
}

func (a *FuseAbility) String() string {
	return "Fuse"
}

// ===== Adventure Cards (Rule 715) =====

// AdventureCardState tracks the state of an Adventure card
// Rule 715: Adventure cards have a creature face and an Adventure face
type AdventureCardState struct {
	cardID          uuid.UUID
	creatureFaceID  uuid.UUID
	adventureFaceID uuid.UUID
	creatureName    string
	adventureName   string
	isInExile       bool      // Adventure sent the creature to exile
	exileZoneID     uuid.UUID // Special exile zone for this Adventure
}

// NewAdventureCardState creates a new Adventure card state
func NewAdventureCardState(cardID uuid.UUID, creatureName, adventureName string) *AdventureCardState {
	return &AdventureCardState{
		cardID:        cardID,
		creatureName:  creatureName,
		adventureName: adventureName,
		isInExile:     false,
	}
}

// GetCreatureFaceID returns the creature face ID
func (acs *AdventureCardState) GetCreatureFaceID() uuid.UUID {
	return acs.creatureFaceID
}

// GetAdventureFaceID returns the Adventure face ID
func (acs *AdventureCardState) GetAdventureFaceID() uuid.UUID {
	return acs.adventureFaceID
}

// IsInAdventureExile returns whether card is in exile from Adventure
func (acs *AdventureCardState) IsInAdventureExile() bool {
	return acs.isInExile
}

// AdventureAbility represents the Adventure mechanic
// Rule 715.3a: When casting a spell with Adventure, a player may choose to cast
// the Adventure (the effect) instead of the creature
type AdventureAbility struct {
	baseAbility
	adventureName string
	adventureCost *ManaCost
}

// NewAdventureAbility creates an Adventure ability
func NewAdventureAbility(source uuid.UUID, adventureName string, cost *ManaCost) *AdventureAbility {
	return &AdventureAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		adventureName: adventureName,
		adventureCost: cost,
	}
}

func (a *AdventureAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *AdventureAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *AdventureAbility) Resolve(ctx context.Context, game GameContext) error {
	// Adventure doesn't resolve - it modifies casting options
	// Rule 715.4: If Adventure would go anywhere but exile, it goes to exile instead
	return nil
}

func (a *AdventureAbility) GetAdventureName() string {
	return a.adventureName
}

func (a *AdventureAbility) GetAdventureCost() *ManaCost {
	return a.adventureCost
}

func (a *AdventureAbility) String() string {
	return fmt.Sprintf("Adventure — %s %s", a.adventureName, a.adventureCost.String())
}

// ===== Saga Cards (Rule 714) =====

// SagaState tracks the state of a Saga enchantment
// Rule 714: Sagas are enchantments with chapter abilities
type SagaState struct {
	permanentID    uuid.UUID
	currentChapter int
	finalChapter   int
	hasTriggered   map[int]bool // Track which chapters have triggered
}

// NewSagaState creates a new Saga state
func NewSagaState(permanentID uuid.UUID, finalChapter int) *SagaState {
	return &SagaState{
		permanentID:    permanentID,
		currentChapter: 0,
		finalChapter:   finalChapter,
		hasTriggered:   make(map[int]bool),
	}
}

// GetCurrentChapter returns the current chapter (lore counter count)
func (ss *SagaState) GetCurrentChapter() int {
	return ss.currentChapter
}

// AddLoreCounter adds a lore counter and returns new chapter
func (ss *SagaState) AddLoreCounter() int {
	ss.currentChapter++
	return ss.currentChapter
}

// HasReachedFinalChapter checks if saga has reached its final chapter
func (ss *SagaState) HasReachedFinalChapter() bool {
	return ss.currentChapter >= ss.finalChapter
}

// HasTriggeredChapter checks if a chapter has already triggered
func (ss *SagaState) HasTriggeredChapter(chapter int) bool {
	return ss.hasTriggered[chapter]
}

// MarkChapterTriggered marks a chapter as having triggered
func (ss *SagaState) MarkChapterTriggered(chapter int) {
	ss.hasTriggered[chapter] = true
}

// SagaAbility represents a Saga enchantment's overall ability
// Rule 714.2a: Saga is an enchantment type, not an ability
// This represents the framework that adds lore counters and triggers chapters
type SagaAbility struct {
	baseAbility
	chapterAbilities map[int]*TriggeredAbility // Chapter number -> ability
	finalChapter     int
}

// NewSagaAbility creates a Saga ability framework
func NewSagaAbility(source uuid.UUID, finalChapter int) *SagaAbility {
	return &SagaAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		chapterAbilities: make(map[int]*TriggeredAbility),
		finalChapter:     finalChapter,
	}
}

// AddChapterAbility adds a chapter ability
func (a *SagaAbility) AddChapterAbility(chapter int, ability *TriggeredAbility) {
	a.chapterAbilities[chapter] = ability
}

// GetChapterAbility returns the ability for a chapter
func (a *SagaAbility) GetChapterAbility(chapter int) (*TriggeredAbility, bool) {
	ability, ok := a.chapterAbilities[chapter]
	return ability, ok
}

func (a *SagaAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *SagaAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *SagaAbility) Resolve(ctx context.Context, game GameContext) error {
	// Saga doesn't resolve - it creates triggered abilities
	// Rule 714.3a: At precombat main phase, add lore counter and trigger chapter
	return nil
}

func (a *SagaAbility) String() string {
	return fmt.Sprintf("Saga (Chapters I-%s)", romanNumeral(a.finalChapter))
}

// SagaChapterTrigger represents a chapter ability trigger
// Rule 714.3b: Triggered ability triggers when lore counter added for that chapter
type SagaChapterTrigger struct {
	*TriggeredAbility
	chapter    int
	sagaSource uuid.UUID
}

// NewSagaChapterTrigger creates a chapter trigger
func NewSagaChapterTrigger(chapter int, sagaSource uuid.UUID, effect Effect) *SagaChapterTrigger {
	trigger := &SagaChapterTrigger{
		chapter:    chapter,
		sagaSource: sagaSource,
	}

	// Build triggered ability
	// TODO: Implement custom trigger for lore counter addition
	trigger.TriggeredAbility = NewTriggeredAbilityBuilder(sagaSource).
		AddEffect(effect).
		Build()

	return trigger
}

// GetChapter returns the chapter number
func (sct *SagaChapterTrigger) GetChapter() int {
	return sct.chapter
}

// ===== Class Cards (Rule 716) =====

// ClassState tracks the state of a Class enchantment
// Rule 716: Classes are enchantments with level-up abilities
type ClassState struct {
	permanentID  uuid.UUID
	currentLevel int // 1, 2, or 3
}

// NewClassState creates a new Class state
func NewClassState(permanentID uuid.UUID) *ClassState {
	return &ClassState{
		permanentID:  permanentID,
		currentLevel: 1, // Always starts at level 1
	}
}

// GetCurrentLevel returns the current class level
func (cs *ClassState) GetCurrentLevel() int {
	return cs.currentLevel
}

// LevelUp increases the class level
func (cs *ClassState) LevelUp() error {
	if cs.currentLevel >= 3 {
		return fmt.Errorf("class already at maximum level")
	}
	cs.currentLevel++
	return nil
}

// ClassAbility represents a Class enchantment's overall ability
// Rule 716.2: A Class card has a class level bar with three class levels
type ClassAbility struct {
	baseAbility
	level2Cost      *ManaCost
	level3Cost      *ManaCost
	level1Abilities []Ability
	level2Abilities []Ability
	level3Abilities []Ability
}

// NewClassAbility creates a Class ability
func NewClassAbility(source uuid.UUID, level2Cost, level3Cost *ManaCost) *ClassAbility {
	return &ClassAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		level2Cost:      level2Cost,
		level3Cost:      level3Cost,
		level1Abilities: make([]Ability, 0),
		level2Abilities: make([]Ability, 0),
		level3Abilities: make([]Ability, 0),
	}
}

// AddLevel1Ability adds an ability available at level 1+
func (a *ClassAbility) AddLevel1Ability(ability Ability) {
	a.level1Abilities = append(a.level1Abilities, ability)
}

// AddLevel2Ability adds an ability available at level 2+
func (a *ClassAbility) AddLevel2Ability(ability Ability) {
	a.level2Abilities = append(a.level2Abilities, ability)
}

// AddLevel3Ability adds an ability available at level 3
func (a *ClassAbility) AddLevel3Ability(ability Ability) {
	a.level3Abilities = append(a.level3Abilities, ability)
}

// GetAbilitiesForLevel returns all active abilities for a given level
func (a *ClassAbility) GetAbilitiesForLevel(level int) []Ability {
	abilities := make([]Ability, 0)

	// Level 1 abilities always active
	abilities = append(abilities, a.level1Abilities...)

	// Level 2 abilities if level 2+
	if level >= 2 {
		abilities = append(abilities, a.level2Abilities...)
	}

	// Level 3 abilities if level 3
	if level >= 3 {
		abilities = append(abilities, a.level3Abilities...)
	}

	return abilities
}

func (a *ClassAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *ClassAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *ClassAbility) Resolve(ctx context.Context, game GameContext) error {
	// Class doesn't resolve - it grants abilities based on level
	return nil
}

func (a *ClassAbility) String() string {
	return fmt.Sprintf("Class (Level up: %s, then %s)", a.level2Cost.String(), a.level3Cost.String())
}

// ClassLevelUpAbility represents the activated ability to level up a Class
// Rule 716.4: Level up is an activated ability
type ClassLevelUpAbility struct {
	*ActivatedAbility
	targetLevel int
	cost        *ManaCost
}

// NewClassLevelUpAbility creates a level up activated ability
func NewClassLevelUpAbility(source uuid.UUID, targetLevel int, cost *ManaCost) *ClassLevelUpAbility {
	ability := &ClassLevelUpAbility{
		targetLevel: targetLevel,
		cost:        cost,
	}

	// Build activated ability
	ability.ActivatedAbility = NewActivatedAbilityBuilder(source).
		AddManaCost(cost.String()).
		AddEffect(NewClassLevelUpEffect(source, targetLevel)).
		Build()

	return ability
}

// GetTargetLevel returns the level this ability levels up to
func (clua *ClassLevelUpAbility) GetTargetLevel() int {
	return clua.targetLevel
}

// ClassLevelUpEffect represents the effect of leveling up a Class
type ClassLevelUpEffect struct {
	description string
	source      uuid.UUID
	targetLevel int
}

// NewClassLevelUpEffect creates a level up effect
func NewClassLevelUpEffect(source uuid.UUID, targetLevel int) *ClassLevelUpEffect {
	return &ClassLevelUpEffect{
		description: "Level up",
		source:      source,
		targetLevel: targetLevel,
	}
}

// Apply implements the Effect interface
func (e *ClassLevelUpEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. Get Class permanent
	// 2. Check current level
	// 3. If current level is exactly one less than target, level up
	// 4. Otherwise, can't activate
	//
	// Rule 716.4b: You can't activate a level up ability if the Class is already
	// the level immediately below the indicated level or higher

	return fmt.Errorf("class level up not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *ClassLevelUpEffect) GetDescription() string {
	return fmt.Sprintf("%s (Gain level %d)", e.description, e.targetLevel)
}

// ===== Flip Cards (Rule 710) =====

// FlipCardState tracks the state of a Flip card
// Rule 710: Flip cards have a card image that can be turned 180 degrees
type FlipCardState struct {
	permanentID uuid.UUID
	isFlipped   bool
	normalName  string
	flippedName string
}

// NewFlipCardState creates a new Flip card state
func NewFlipCardState(permanentID uuid.UUID, normalName, flippedName string) *FlipCardState {
	return &FlipCardState{
		permanentID: permanentID,
		isFlipped:   false,
		normalName:  normalName,
		flippedName: flippedName,
	}
}

// IsFlipped returns whether the card is flipped
func (fcs *FlipCardState) IsFlipped() bool {
	return fcs.isFlipped
}

// Flip flips the card
func (fcs *FlipCardState) Flip() {
	fcs.isFlipped = true
}

// ===== Leveler Cards (Rule 711) =====

// LevelerState tracks the state of a Leveler creature
// Rule 711: Leveler cards have level up abilities and level symbols
type LevelerState struct {
	permanentID   uuid.UUID
	levelCounters int
	levelRanges   []LevelRange
}

// LevelRange represents a level range and its P/T/abilities
type LevelRange struct {
	minLevel  int
	maxLevel  int
	power     int
	toughness int
	abilities []Ability
}

// NewLevelerState creates a new Leveler state
func NewLevelerState(permanentID uuid.UUID) *LevelerState {
	return &LevelerState{
		permanentID:   permanentID,
		levelCounters: 0,
		levelRanges:   make([]LevelRange, 0),
	}
}

// GetLevelCounters returns the number of level counters
func (ls *LevelerState) GetLevelCounters() int {
	return ls.levelCounters
}

// AddLevelCounters adds level counters
func (ls *LevelerState) AddLevelCounters(count int) {
	ls.levelCounters += count
}

// GetCurrentLevelRange returns the current level range
func (ls *LevelerState) GetCurrentLevelRange() *LevelRange {
	for _, lr := range ls.levelRanges {
		if ls.levelCounters >= lr.minLevel && ls.levelCounters <= lr.maxLevel {
			return &lr
		}
	}
	return nil
}

// LevelUpAbility represents the Level Up keyword
// Rule 702.87: Level up is an activated ability
type LevelUpAbility struct {
	*ActivatedAbility
	cost *ManaCost
}

// NewLevelUpAbility creates a Level Up ability
func NewLevelUpAbility(source uuid.UUID, cost *ManaCost) *LevelUpAbility {
	ability := &LevelUpAbility{
		cost: cost,
	}

	// Build activated ability
	ability.ActivatedAbility = NewActivatedAbilityBuilder(source).
		AddManaCost(cost.String()).
		AddEffect(NewAddLevelCounterEffect(source)).
		Build()

	return ability
}

// AddLevelCounterEffect represents adding a level counter
type AddLevelCounterEffect struct {
	description string
	source      uuid.UUID
}

// NewAddLevelCounterEffect creates an effect to add a level counter
func NewAddLevelCounterEffect(source uuid.UUID) *AddLevelCounterEffect {
	return &AddLevelCounterEffect{
		description: "Level up",
		source:      source,
	}
}

// Apply implements the Effect interface
func (e *AddLevelCounterEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Get counter game context
	counterGame, ok := game.(CounterGameContext)
	if !ok {
		return fmt.Errorf("game context does not support counter operations")
	}

	// Get permanent
	permanent, err := counterGame.GetPermanent(e.source)
	if err != nil {
		return fmt.Errorf("failed to get permanent: %w", err)
	}

	// Add level counter
	counter := counters.NewCounter("level", 1)
	if err := counterGame.AddCountersToPermanent(permanent, counter); err != nil {
		return fmt.Errorf("failed to add level counter: %w", err)
	}

	return nil
}

// GetDescription returns the effect description
func (e *AddLevelCounterEffect) GetDescription() string {
	return e.description
}

// ===== Helper Functions =====

// romanNumeral converts a number to Roman numerals (for Sagas)
func romanNumeral(n int) string {
	if n <= 0 || n > 10 {
		return fmt.Sprintf("%d", n)
	}

	numerals := []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X"}
	return numerals[n]
}

// ===== Common Card Examples =====

// Split Cards:
// - Fire // Ice: "{1}{R} Instant // {1}{U} Instant"
// - Wear // Tear: "{1}{R} Instant // {W} Instant"
// - Turn // Burn (with Fuse): "{2}{U} Instant // {1}{R} Instant, Fuse"
// - Breaking // Entering (with Fuse): "{U}{B} Sorcery // {4}{B}{R} Sorcery, Fuse"

// Adventure Cards:
// - Brazen Borrower // Petty Theft: "{1}{U}{U} Creature - Faerie Rogue 3/1 // {1}{U} Instant"
// - Bonecrusher Giant // Stomp: "{2}{R} Creature - Giant 4/3 // {1}{R} Instant"
// - Murderous Rider // Swift End: "{1}{B}{B} Creature - Zombie Knight 2/3 // {1}{B}{B} Instant"

// Saga Cards:
// - The Eldest Reborn: "{4}{B} Enchantment - Saga"
//   I: Each opponent sacrifices a creature or planeswalker
//   II: Each opponent discards a card
//   III: Put target creature or planeswalker from a graveyard onto battlefield under your control
// - History of Benalia: "{1}{W}{W} Enchantment - Saga"
//   I, II: Create a 2/2 white Knight creature token with vigilance
//   III: Knights you control get +2/+1 until end of turn

// Class Cards:
// - Barbarian Class: "{R} Enchantment - Class"
//   Level 1: If a source you control would deal damage to an opponent, it deals that much plus 1 instead
//   Level 2 ({1}{R}): Creatures you control have menace
//   Level 3 ({2}{R}): When this Class becomes level 3, deal 2 damage to any target

// Flip Cards:
// - Nezumi Graverobber // Nighteyes the Desecrator
// - Kitsune Mystic // Autumn-Tail, Kitsune Sage

// Leveler Cards:
// - Figure of Destiny: "{R/W} Creature - Kithkin 1/1, Level up {R/W}"
//   LEVEL 1-3: 2/2
//   LEVEL 4-7: 4/4, First strike
//   LEVEL 8+: 8/8, First strike, flying
// - Student of Warfare: "{W} Creature - Human Knight 1/1, Level up {W}"
