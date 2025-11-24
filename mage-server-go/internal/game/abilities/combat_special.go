package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Special Combat Mechanics (Rules 702.21-702.22)
// This file implements rare and complex combat mechanics:
// - Banding (Rule 702.21)
// - Bands with other (Rule 702.22)
// - Flanking (Rule 702.24)
// - Rampage (Rule 702.23)
// - Bushido (Rule 702.44)
// - Exalted (Rule 702.83)

// ===== Banding (Rule 702.21) =====

// BandingAbility represents the Banding keyword
// Rule 702.21: Banding
type BandingAbility struct {
	baseAbility
}

// NewBandingAbility creates a Banding ability
func NewBandingAbility(source uuid.UUID) *BandingAbility {
	return &BandingAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *BandingAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *BandingAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return false
}

func (a *BandingAbility) Resolve(ctx context.Context, game GameContext) error {
	// Banding is a static ability
	// Rule 702.21j: Creatures with banding can attack in a band
	// Rule 702.21k: Defending player assigns combat damage for banded attackers
	return nil
}

func (a *BandingAbility) String() string {
	return "Banding"
}

// BandsWithOtherAbility represents "Bands with other [quality]"
// Rule 702.22: Bands with other
type BandsWithOtherAbility struct {
	baseAbility
	quality string // e.g., "Dinosaurs", "Legends"
}

// NewBandsWithOtherAbility creates a "Bands with other" ability
func NewBandsWithOtherAbility(source uuid.UUID, quality string) *BandsWithOtherAbility {
	return &BandsWithOtherAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		quality: quality,
	}
}

func (a *BandsWithOtherAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *BandsWithOtherAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return false
}

func (a *BandsWithOtherAbility) Resolve(ctx context.Context, game GameContext) error {
	// Bands with other allows banding with creatures of specific quality
	return nil
}

func (a *BandsWithOtherAbility) String() string {
	return fmt.Sprintf("Bands with other %s", a.quality)
}

// BandedGroup represents a group of creatures attacking as a band
// Rule 702.21j: Creatures with banding and up to one without can band
type BandedGroup struct {
	groupID    uuid.UUID
	creatures  []uuid.UUID
	hasBanding map[uuid.UUID]bool // which creatures have banding
	controller uuid.UUID          // who controls damage assignment
}

// NewBandedGroup creates a new banded attacking group
func NewBandedGroup(controller uuid.UUID) *BandedGroup {
	return &BandedGroup{
		groupID:    uuid.New(),
		creatures:  make([]uuid.UUID, 0),
		hasBanding: make(map[uuid.UUID]bool),
		controller: controller,
	}
}

// AddCreature adds a creature to the band
func (bg *BandedGroup) AddCreature(creature uuid.UUID, hasBanding bool) error {
	// Rule 702.21j: Can have at most one creature without banding
	if !hasBanding {
		// Check if we already have a non-banding creature
		for c, hasB := range bg.hasBanding {
			if !hasB && c != creature {
				return fmt.Errorf("band can have at most one creature without banding")
			}
		}
	}

	bg.creatures = append(bg.creatures, creature)
	bg.hasBanding[creature] = hasBanding
	return nil
}

// GetCreatures returns all creatures in the band
func (bg *BandedGroup) GetCreatures() []uuid.UUID {
	return bg.creatures
}

// IsValidBand checks if the band composition is valid
func (bg *BandedGroup) IsValidBand() bool {
	if len(bg.creatures) == 0 {
		return false
	}

	nonBandingCount := 0
	for _, hasBanding := range bg.hasBanding {
		if !hasBanding {
			nonBandingCount++
		}
	}

	// Rule 702.21j: At most one without banding
	return nonBandingCount <= 1
}

// ===== Flanking (Rule 702.24) =====

// FlankingAbility represents the Flanking keyword
// Rule 702.24: Flanking
type FlankingAbility struct {
	baseAbility
}

// NewFlankingAbility creates a Flanking ability
func NewFlankingAbility(source uuid.UUID) *FlankingAbility {
	return &FlankingAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *FlankingAbility) GetType() AbilityType {
	return AbilityTypeTriggered
}

func (a *FlankingAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return false // Triggered, not activated
}

func (a *FlankingAbility) Resolve(ctx context.Context, game GameContext) error {
	// When blocked by creature without flanking, that creature gets -1/-1
	return nil
}

func (a *FlankingAbility) String() string {
	return "Flanking"
}

// FlankingTrigger represents a flanking trigger
type FlankingTrigger struct {
	*TriggeredAbility
	blockerID uuid.UUID
}

// NewFlankingTrigger creates a flanking trigger
func NewFlankingTrigger(source, blocker uuid.UUID) *FlankingTrigger {
	trigger := &FlankingTrigger{
		blockerID: blocker,
	}

	// Create the triggered ability
	trigger.TriggeredAbility = NewTriggeredAbilityBuilder(source).
		SetTrigger(NewBlocksTrigger(source)).
		AddEffect(NewBoostEffect(-1, -1)). // -1/-1 until end of turn
		SetOptional(false).
		Build()

	return trigger
}

// ===== Rampage (Rule 702.23) =====

// RampageAbility represents the Rampage keyword
// Rule 702.23: Rampage N
type RampageAbility struct {
	baseAbility
	rampageAmount int // +N/+N for each blocker beyond the first
}

// NewRampageAbility creates a Rampage ability
func NewRampageAbility(source uuid.UUID, amount int) *RampageAbility {
	return &RampageAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		rampageAmount: amount,
	}
}

func (a *RampageAbility) GetType() AbilityType {
	return AbilityTypeTriggered
}

func (a *RampageAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return false // Triggered, not activated
}

func (a *RampageAbility) Resolve(ctx context.Context, game GameContext) error {
	// When blocked, gets +N/+N for each creature blocking beyond the first
	return nil
}

func (a *RampageAbility) String() string {
	return fmt.Sprintf("Rampage %d", a.rampageAmount)
}

// GetRampageBonus calculates the rampage bonus
func (a *RampageAbility) GetRampageBonus(blockerCount int) (int, int) {
	if blockerCount <= 1 {
		return 0, 0
	}

	bonus := (blockerCount - 1) * a.rampageAmount
	return bonus, bonus
}

// ===== Bushido (Rule 702.44) =====

// BushidoAbility represents the Bushido keyword
// Rule 702.44: Bushido N
type BushidoAbility struct {
	baseAbility
	bushidoAmount int // +N/+N when blocks or becomes blocked
}

// NewBushidoAbility creates a Bushido ability
func NewBushidoAbility(source uuid.UUID, amount int) *BushidoAbility {
	return &BushidoAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		bushidoAmount: amount,
	}
}

func (a *BushidoAbility) GetType() AbilityType {
	return AbilityTypeTriggered
}

func (a *BushidoAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return false // Triggered, not activated
}

func (a *BushidoAbility) Resolve(ctx context.Context, game GameContext) error {
	// When blocks or becomes blocked, gets +N/+N until end of turn
	return nil
}

func (a *BushidoAbility) String() string {
	return fmt.Sprintf("Bushido %d", a.bushidoAmount)
}

// ===== Exalted (Rule 702.83) =====

// ExaltedAbility represents the Exalted keyword
// Rule 702.83: Exalted
type ExaltedAbility struct {
	baseAbility
	exaltedCount int // Some cards have "Exalted, Exalted" (multiple instances)
}

// NewExaltedAbility creates an Exalted ability
func NewExaltedAbility(source uuid.UUID) *ExaltedAbility {
	return &ExaltedAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		exaltedCount: 1,
	}
}

func (a *ExaltedAbility) GetType() AbilityType {
	return AbilityTypeTriggered
}

func (a *ExaltedAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return false // Triggered, not activated
}

func (a *ExaltedAbility) Resolve(ctx context.Context, game GameContext) error {
	// Whenever a creature you control attacks alone, it gets +1/+1
	return nil
}

func (a *ExaltedAbility) String() string {
	if a.exaltedCount > 1 {
		return fmt.Sprintf("Exalted (x%d)", a.exaltedCount)
	}
	return "Exalted"
}

// ===== Shadow (Rule 702.28) =====

// ShadowAbility represents the Shadow keyword
// Rule 702.28: Shadow
type ShadowAbility struct {
	baseAbility
}

// NewShadowAbility creates a Shadow ability
func NewShadowAbility(source uuid.UUID) *ShadowAbility {
	return &ShadowAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *ShadowAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *ShadowAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return false
}

func (a *ShadowAbility) Resolve(ctx context.Context, game GameContext) error {
	// Can only block or be blocked by creatures with shadow
	return nil
}

func (a *ShadowAbility) String() string {
	return "Shadow"
}

// ===== Horsemanship (Rule 702.31) =====

// HorsemanshipAbility represents the Horsemanship keyword
// Rule 702.31: Horsemanship (like flying for Portal Three Kingdoms)
type HorsemanshipAbility struct {
	baseAbility
}

// NewHorsemanshipAbility creates a Horsemanship ability
func NewHorsemanshipAbility(source uuid.UUID) *HorsemanshipAbility {
	return &HorsemanshipAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *HorsemanshipAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *HorsemanshipAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return false
}

func (a *HorsemanshipAbility) Resolve(ctx context.Context, game GameContext) error {
	// Can only be blocked by creatures with horsemanship
	return nil
}

func (a *HorsemanshipAbility) String() string {
	return "Horsemanship"
}

// ===== Fear (Rule 702.36) =====

// FearAbility represents the Fear keyword
// Rule 702.36: Fear (can't be blocked except by black/artifact)
type FearAbility struct {
	baseAbility
}

// NewFearAbility creates a Fear ability
func NewFearAbility(source uuid.UUID) *FearAbility {
	return &FearAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *FearAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *FearAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return false
}

func (a *FearAbility) Resolve(ctx context.Context, game GameContext) error {
	// Can't be blocked except by artifact creatures and/or black creatures
	return nil
}

func (a *FearAbility) String() string {
	return "Fear"
}

// ===== Intimidate (Rule 702.13) =====

// IntimidateAbility represents the Intimidate keyword
// Rule 702.13: Intimidate (can't be blocked except by artifact/same color)
type IntimidateAbility struct {
	baseAbility
}

// NewIntimidateAbility creates an Intimidate ability
func NewIntimidateAbility(source uuid.UUID) *IntimidateAbility {
	return &IntimidateAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *IntimidateAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *IntimidateAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return false
}

func (a *IntimidateAbility) Resolve(ctx context.Context, game GameContext) error {
	// Can't be blocked except by artifact creatures and/or creatures that share a color
	return nil
}

func (a *IntimidateAbility) String() string {
	return "Intimidate"
}

// ===== Phasing (Rule 702.26) =====

// PhasingAbility represents the Phasing keyword
// Rule 702.26: Phasing
type PhasingAbility struct {
	baseAbility
	phaseState PhasingState
}

// PhasingState tracks whether a permanent is phased in or out
type PhasingState int

const (
	PhasingPhasedIn  PhasingState = iota // Normal (phased in)
	PhasingPhasedOut                     // Phased out (treated as though it doesn't exist)
)

// NewPhasingAbility creates a Phasing ability
func NewPhasingAbility(source uuid.UUID) *PhasingAbility {
	return &PhasingAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		phaseState: PhasingPhasedIn,
	}
}

func (a *PhasingAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *PhasingAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return false
}

func (a *PhasingAbility) Resolve(ctx context.Context, game GameContext) error {
	// Phases in/out during untap step
	return nil
}

func (a *PhasingAbility) String() string {
	return "Phasing"
}

// PhaseOut phases out a permanent
func (a *PhasingAbility) PhaseOut() {
	a.phaseState = PhasingPhasedOut
}

// PhaseIn phases in a permanent
func (a *PhasingAbility) PhaseIn() {
	a.phaseState = PhasingPhasedIn
}

// IsPhasedOut checks if phased out
func (a *PhasingAbility) IsPhasedOut() bool {
	return a.phaseState == PhasingPhasedOut
}

// ===== Combat Triggers for Special Abilities =====

// AttackAloneTrigger triggers when a creature attacks alone
type AttackAloneTrigger struct {
	*TriggeredAbility
}

// NewAttackAloneTrigger creates an "attacks alone" trigger
func NewAttackAloneTrigger(source uuid.UUID) *AttackAloneTrigger {
	trigger := &AttackAloneTrigger{}

	trigger.TriggeredAbility = NewTriggeredAbilityBuilder(source).
		SetTrigger(NewAttacksTrigger(source)).
		SetOptional(false).
		Build()

	return trigger
}

// ShouldTrigger checks if this should trigger (only when attacking alone)
func (t *AttackAloneTrigger) ShouldTrigger(attackerCount int) bool {
	return attackerCount == 1
}

// BlockedByMultipleTrigger triggers when blocked by multiple creatures
type BlockedByMultipleTrigger struct {
	*TriggeredAbility
	blockerCount int
}

// NewBlockedByMultipleTrigger creates a "blocked by multiple" trigger
func NewBlockedByMultipleTrigger(source uuid.UUID) *BlockedByMultipleTrigger {
	trigger := &BlockedByMultipleTrigger{
		blockerCount: 0,
	}

	trigger.TriggeredAbility = NewTriggeredAbilityBuilder(source).
		SetTrigger(NewBecomesBlockedTrigger(source)).
		SetOptional(false).
		Build()

	return trigger
}

// ShouldTrigger checks if this should trigger (only when blocked by 2+)
func (t *BlockedByMultipleTrigger) ShouldTrigger(blockerCount int) bool {
	t.blockerCount = blockerCount
	return blockerCount >= 2
}

// GetBlockerCount returns the number of blockers
func (t *BlockedByMultipleTrigger) GetBlockerCount() int {
	return t.blockerCount
}

// ===== Helper Functions =====

// HasBanding checks if a creature has banding
func HasBanding(permanentID uuid.UUID, game GameContext) bool {
	// TODO: Check permanent for BandingAbility
	return false
}

// HasShadow checks if a creature has shadow
func HasShadow(permanentID uuid.UUID, game GameContext) bool {
	// TODO: Check permanent for ShadowAbility
	return false
}

// HasHorsemanship checks if a creature has horsemanship
func HasHorsemanship(permanentID uuid.UUID, game GameContext) bool {
	// TODO: Check permanent for HorsemanshipAbility
	return false
}

// HasFear checks if a creature has fear
func HasFear(permanentID uuid.UUID, game GameContext) bool {
	// TODO: Check permanent for FearAbility
	return false
}

// HasIntimidate checks if a creature has intimidate
func HasIntimidate(permanentID uuid.UUID, game GameContext) bool {
	// TODO: Check permanent for IntimidateAbility
	return false
}

// CanBlockWithShadow checks if a shadow creature can block another creature
func CanBlockWithShadow(blocker, attacker uuid.UUID, game GameContext) bool {
	blockerHasShadow := HasShadow(blocker, game)
	attackerHasShadow := HasShadow(attacker, game)

	// Shadow creatures can only block other shadow creatures
	if blockerHasShadow {
		return attackerHasShadow
	}

	// Non-shadow creatures can't block shadow creatures
	return !attackerHasShadow
}

// ===== Integration with Existing Systems =====

// These special combat mechanics integrate with:
// - Combat system: Modify attack/block legality and damage assignment
// - Triggered abilities: Most trigger on combat events
// - Static abilities: Banding, Shadow, Fear, Intimidate
// - State-based actions: Flanking creates state-based P/T changes
// - Turn structure: Phasing occurs during untap step

// Important rules:
// - Rule 702.21: Banding (complex damage assignment rules)
// - Rule 702.22: Bands with other
// - Rule 702.23: Rampage
// - Rule 702.24: Flanking
// - Rule 702.26: Phasing
// - Rule 702.28: Shadow
// - Rule 702.31: Horsemanship
// - Rule 702.36: Fear
// - Rule 702.13: Intimidate
// - Rule 702.44: Bushido
// - Rule 702.83: Exalted

// Example usage for Banding:
// 1. Three creatures with banding attack together
// 2. Create BandedGroup for the three creatures
// 3. ValidateValidBand() confirms composition
// 4. During damage assignment, defending player (not attacker) assigns damage
// 5. All damage to band is assigned by defending player

// Example usage for Flanking:
// 1. Creature with Flanking attacks
// 2. Blocked by creature without flanking
// 3. FlankingTrigger fires
// 4. Blocker gets -1/-1 until end of turn
// 5. Combat damage calculated with reduced toughness

// Example usage for Exalted:
// 1. You control 3 creatures with Exalted
// 2. Declare exactly 1 attacker
// 3. Three Exalted triggers fire
// 4. Attacker gets +3/+3 until end of turn
// 5. Combat proceeds with boosted creature

// Example usage for Shadow:
// 1. Creature with Shadow attacks
// 2. Opponent declares blockers
// 3. Check CanBlockWithShadow() for each potential blocker
// 4. Only creatures with Shadow can block
// 5. If no shadow blockers, damage goes through unblocked
