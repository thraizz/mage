package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Combat Damage System (Rules 510, 702.19, 702.2)
// This file implements advanced combat damage mechanics including:
// - Damage prevention effects
// - Damage replacement effects
// - Trample calculations with deathtouch
// - Combat damage triggers
// - Damage redirection

// ===== Combat Damage Effects =====

// PreventCombatDamageEffect prevents combat damage to a target
// Rule 615: Prevention effects
type PreventCombatDamageEffect struct {
	description       string
	source            uuid.UUID
	targetID          uuid.UUID
	amount            int  // -1 = all damage
	preventNextDamage bool // Only prevent next damage instance
	duration          Duration
}

// NewPreventCombatDamageEffect creates a prevention effect
func NewPreventCombatDamageEffect(source, target uuid.UUID, amount int, duration Duration) *PreventCombatDamageEffect {
	return &PreventCombatDamageEffect{
		description:       "Prevent combat damage",
		source:            source,
		targetID:          target,
		amount:            amount,
		preventNextDamage: false,
		duration:          duration,
	}
}

// Apply implements the Effect interface
func (e *PreventCombatDamageEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Prevention effects are applied during damage resolution, not here
	// This registers the prevention effect in the game's replacement effect system
	return nil
}

// GetDescription returns effect description
func (e *PreventCombatDamageEffect) GetDescription() string {
	if e.amount == -1 {
		return fmt.Sprintf("Prevent all combat damage that would be dealt to %s", e.targetID)
	}
	return fmt.Sprintf("Prevent the next %d combat damage that would be dealt to %s", e.amount, e.targetID)
}

// ===== Damage Replacement Effects =====

// ReplaceCombatDamageEffect replaces combat damage with different damage
// Rule 614: Replacement effects
type ReplaceCombatDamageEffect struct {
	description  string
	source       uuid.UUID
	targetID     uuid.UUID
	multiplier   float64 // 2.0 for double, 0.5 for half
	addAmount    int     // Additional damage to add
	duration     Duration
	replacesOnce bool // Only replace next damage
}

// NewReplaceCombatDamageEffect creates a replacement effect
func NewReplaceCombatDamageEffect(source, target uuid.UUID, multiplier float64, addAmount int, duration Duration) *ReplaceCombatDamageEffect {
	return &ReplaceCombatDamageEffect{
		description:  "Replace combat damage",
		source:       source,
		targetID:     target,
		multiplier:   multiplier,
		addAmount:    addAmount,
		duration:     duration,
		replacesOnce: false,
	}
}

// Apply implements the Effect interface
func (e *ReplaceCombatDamageEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Replacement effects are applied during damage resolution
	return nil
}

// GetDescription returns effect description
func (e *ReplaceCombatDamageEffect) GetDescription() string {
	if e.multiplier != 1.0 {
		return fmt.Sprintf("Combat damage dealt to %s is multiplied by %.1f", e.targetID, e.multiplier)
	}
	return fmt.Sprintf("Combat damage dealt to %s is increased by %d", e.targetID, e.addAmount)
}

// ===== Damage Redirection Effects =====

// RedirectCombatDamageEffect redirects damage to a different target
// Rule 614.9: Damage redirection
type RedirectCombatDamageEffect struct {
	description      string
	source           uuid.UUID
	fromTarget       uuid.UUID
	toTarget         uuid.UUID
	maxAmount        int // Maximum damage to redirect (-1 = all)
	duration         Duration
	redirectNextOnly bool
}

// NewRedirectCombatDamageEffect creates a redirection effect
func NewRedirectCombatDamageEffect(source, from, to uuid.UUID, maxAmount int, duration Duration) *RedirectCombatDamageEffect {
	return &RedirectCombatDamageEffect{
		description:      "Redirect combat damage",
		source:           source,
		fromTarget:       from,
		toTarget:         to,
		maxAmount:        maxAmount,
		duration:         duration,
		redirectNextOnly: false,
	}
}

// Apply implements the Effect interface
func (e *RedirectCombatDamageEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Redirection is applied during damage resolution
	return nil
}

// GetDescription returns effect description
func (e *RedirectCombatDamageEffect) GetDescription() string {
	if e.maxAmount == -1 {
		return fmt.Sprintf("Redirect all damage from %s to %s", e.fromTarget, e.toTarget)
	}
	return fmt.Sprintf("Redirect up to %d damage from %s to %s", e.maxAmount, e.fromTarget, e.toTarget)
}

// ===== Trample Damage Calculation =====

// TrampleDamageCalculator calculates trample damage with special interactions
// Rule 702.19: Trample
type TrampleDamageCalculator struct {
	attackerID       uuid.UUID
	attackerPower    int
	blockers         []BlockerInfo
	hasDeathtouch    bool
	hasDoubleStrike  bool
	hasFirstStrike   bool
	trampleDamageAll bool // Trample over planeswalkers
}

// BlockerInfo represents information about a blocker for damage calculation
type BlockerInfo struct {
	BlockerID    uuid.UUID
	Toughness    int
	DamageMarked int
	HasIndestruc bool
	LethalDamage int // Calculated lethal damage considering all factors
}

// NewTrampleDamageCalculator creates a calculator for trample damage
func NewTrampleDamageCalculator(attackerID uuid.UUID, power int, hasDeathtouch bool) *TrampleDamageCalculator {
	return &TrampleDamageCalculator{
		attackerID:       attackerID,
		attackerPower:    power,
		blockers:         make([]BlockerInfo, 0),
		hasDeathtouch:    hasDeathtouch,
		hasDoubleStrike:  false,
		hasFirstStrike:   false,
		trampleDamageAll: false,
	}
}

// AddBlocker adds a blocker to the calculation
func (tc *TrampleDamageCalculator) AddBlocker(blockerID uuid.UUID, toughness, damageMarked int, hasIndestruc bool) {
	// Rule 702.19b: Assign lethal damage to blockers before trampling
	// Rule 702.2c: Any amount of damage from deathtouch source is lethal
	lethalDamage := toughness - damageMarked

	if tc.hasDeathtouch && !hasIndestruc {
		// With deathtouch, only 1 damage is lethal (unless they have indestructible)
		if lethalDamage > 1 {
			lethalDamage = 1
		}
	}

	// Indestructible creatures still require toughness worth of damage
	if lethalDamage < 0 {
		lethalDamage = 0
	}

	tc.blockers = append(tc.blockers, BlockerInfo{
		BlockerID:    blockerID,
		Toughness:    toughness,
		DamageMarked: damageMarked,
		HasIndestruc: hasIndestruc,
		LethalDamage: lethalDamage,
	})
}

// CalculateTrampleDamage calculates optimal damage assignment for trample
// Returns: map of blocker ID to damage, and excess damage to trample through
func (tc *TrampleDamageCalculator) CalculateTrampleDamage() (map[uuid.UUID]int, int) {
	assignment := make(map[uuid.UUID]int)
	remainingPower := tc.attackerPower

	// Rule 702.19b: Must assign lethal to each blocker before trampling through
	for _, blocker := range tc.blockers {
		if remainingPower <= 0 {
			break
		}

		damageToAssign := blocker.LethalDamage
		if damageToAssign > remainingPower {
			damageToAssign = remainingPower
		}

		if damageToAssign > 0 {
			assignment[blocker.BlockerID] = damageToAssign
			remainingPower -= damageToAssign
		}
	}

	// Remaining damage tramples through to defender
	trampleDamage := remainingPower
	if trampleDamage < 0 {
		trampleDamage = 0
	}

	return assignment, trampleDamage
}

// CalculateMaximumTrampleDamage calculates maximum damage that can trample through
// This allows players to assign more than lethal to blockers
func (tc *TrampleDamageCalculator) CalculateMaximumTrampleDamage(assignment map[uuid.UUID]int) int {
	totalAssignedToBlockers := 0
	for _, damage := range assignment {
		totalAssignedToBlockers += damage
	}

	trampleDamage := tc.attackerPower - totalAssignedToBlockers
	if trampleDamage < 0 {
		trampleDamage = 0
	}

	return trampleDamage
}

// ValidateTrampleAssignment validates a damage assignment for trample rules
func (tc *TrampleDamageCalculator) ValidateTrampleAssignment(assignment map[uuid.UUID]int) error {
	totalAssigned := 0

	// Check each blocker has at least lethal assigned
	for _, blocker := range tc.blockers {
		assigned, exists := assignment[blocker.BlockerID]
		if !exists {
			assigned = 0
		}

		if assigned < blocker.LethalDamage {
			return fmt.Errorf("must assign at least %d lethal damage to blocker %s, assigned %d",
				blocker.LethalDamage, blocker.BlockerID, assigned)
		}

		totalAssigned += assigned
	}

	// Check total doesn't exceed power
	if totalAssigned > tc.attackerPower {
		return fmt.Errorf("cannot assign more damage (%d) than attacker's power (%d)",
			totalAssigned, tc.attackerPower)
	}

	return nil
}

// ===== Combat Damage Assignment Helpers =====

// DamageAssignmentOrder represents the order for assigning damage
type DamageAssignmentOrder struct {
	sourceID      uuid.UUID
	targets       []uuid.UUID // Ordered list of targets
	playerChooses bool        // True if player controls source and chooses order
}

// NewDamageAssignmentOrder creates a damage assignment order
func NewDamageAssignmentOrder(sourceID uuid.UUID, targets []uuid.UUID, playerChooses bool) *DamageAssignmentOrder {
	return &DamageAssignmentOrder{
		sourceID:      sourceID,
		targets:       targets,
		playerChooses: playerChooses,
	}
}

// GetTargets returns the ordered target list
func (dao *DamageAssignmentOrder) GetTargets() []uuid.UUID {
	return dao.targets
}

// SetOrder sets a new damage assignment order
func (dao *DamageAssignmentOrder) SetOrder(targets []uuid.UUID) error {
	if len(targets) != len(dao.targets) {
		return fmt.Errorf("must specify order for all %d targets, got %d", len(dao.targets), len(targets))
	}

	// Validate all targets are in the original list
	targetMap := make(map[uuid.UUID]bool)
	for _, t := range dao.targets {
		targetMap[t] = true
	}

	for _, t := range targets {
		if !targetMap[t] {
			return fmt.Errorf("target %s not in original target list", t)
		}
	}

	dao.targets = targets
	return nil
}

// ===== Combat Damage Context =====

// CombatDamageContext tracks all damage being dealt in a combat damage step
type CombatDamageContext struct {
	isFirstStrike bool
	damageEvents  []CombatDamageEvent
	replacements  []ReplacementEffect
	preventions   []PreventionEffect
}

// CombatDamageEvent represents a single instance of combat damage
type CombatDamageEvent struct {
	sourceID     uuid.UUID
	targetID     uuid.UUID
	amount       int
	isCombat     bool
	isLifelink   bool
	isDeathtouch bool
	isTrample    bool
	prevented    int // Amount prevented so far
	replaced     bool
	finalAmount  int // Final amount after replacements/prevention
}

// ReplacementEffect represents an active replacement effect
type ReplacementEffect struct {
	effectID  uuid.UUID
	sourceID  uuid.UUID
	applies   func(*CombatDamageEvent) bool
	replace   func(*CombatDamageEvent) *CombatDamageEvent
	usedOnce  bool
	timesUsed int
	maxUses   int // -1 = unlimited
}

// PreventionEffect represents an active prevention effect
type PreventionEffect struct {
	effectID       uuid.UUID
	sourceID       uuid.UUID
	applies        func(*CombatDamageEvent) bool
	preventAmount  int // -1 = all
	preventedTotal int
	usedOnce       bool
}

// NewCombatDamageContext creates a new combat damage context
func NewCombatDamageContext(isFirstStrike bool) *CombatDamageContext {
	return &CombatDamageContext{
		isFirstStrike: isFirstStrike,
		damageEvents:  make([]CombatDamageEvent, 0),
		replacements:  make([]ReplacementEffect, 0),
		preventions:   make([]PreventionEffect, 0),
	}
}

// AddDamageEvent adds a damage event to the context
func (cdc *CombatDamageContext) AddDamageEvent(event CombatDamageEvent) {
	cdc.damageEvents = append(cdc.damageEvents, event)
}

// AddReplacementEffect adds a replacement effect
func (cdc *CombatDamageContext) AddReplacementEffect(effect ReplacementEffect) {
	cdc.replacements = append(cdc.replacements, effect)
}

// AddPreventionEffect adds a prevention effect
func (cdc *CombatDamageContext) AddPreventionEffect(effect PreventionEffect) {
	cdc.preventions = append(cdc.preventions, effect)
}

// ProcessDamageEvents processes all damage events with replacements and prevention
// Rule 614: Replacement effects are applied before prevention
func (cdc *CombatDamageContext) ProcessDamageEvents() []CombatDamageEvent {
	processed := make([]CombatDamageEvent, 0)

	for _, event := range cdc.damageEvents {
		// Apply replacement effects first (Rule 614.1)
		modifiedEvent := event
		for i := range cdc.replacements {
			replacement := &cdc.replacements[i]

			// Skip if already used and one-time only
			if replacement.usedOnce && replacement.timesUsed > 0 {
				continue
			}

			// Skip if max uses reached
			if replacement.maxUses != -1 && replacement.timesUsed >= replacement.maxUses {
				continue
			}

			// Apply if applicable
			if replacement.applies(&modifiedEvent) {
				modifiedEvent = *replacement.replace(&modifiedEvent)
				replacement.timesUsed++
			}
		}

		// Apply prevention effects (Rule 615.1)
		preventedAmount := 0
		for i := range cdc.preventions {
			prevention := &cdc.preventions[i]

			// Skip if already used and one-time only
			if prevention.usedOnce && prevention.preventedTotal > 0 {
				continue
			}

			// Apply if applicable
			if prevention.applies(&modifiedEvent) {
				amountToPrevent := prevention.preventAmount
				if amountToPrevent == -1 {
					// Prevent all
					amountToPrevent = modifiedEvent.amount - preventedAmount
				} else {
					// Prevent up to specified amount
					remaining := prevention.preventAmount - prevention.preventedTotal
					if remaining > modifiedEvent.amount-preventedAmount {
						amountToPrevent = modifiedEvent.amount - preventedAmount
					} else {
						amountToPrevent = remaining
					}
				}

				preventedAmount += amountToPrevent
				prevention.preventedTotal += amountToPrevent

				if preventedAmount >= modifiedEvent.amount {
					break
				}
			}
		}

		modifiedEvent.prevented = preventedAmount
		modifiedEvent.finalAmount = modifiedEvent.amount - preventedAmount

		if modifiedEvent.finalAmount > 0 {
			processed = append(processed, modifiedEvent)
		}
	}

	return processed
}

// ===== Helper Functions =====

// CalculateLethalDamage calculates lethal damage considering current damage and effects
func CalculateLethalDamage(toughness, markedDamage int, hasDeathtouch, hasIndestructible bool) int {
	if hasIndestructible {
		// Indestructible creatures still require full toughness worth
		return toughness - markedDamage
	}

	if hasDeathtouch {
		// Any damage from deathtouch is lethal
		if markedDamage >= 1 {
			return 0 // Already lethal
		}
		return 1
	}

	lethal := toughness - markedDamage
	if lethal < 0 {
		lethal = 0
	}
	return lethal
}

// ===== Integration with Existing Systems =====

// These combat damage mechanics integrate with:
// - Combat system: Damage assignment during combat damage step
// - Replacement effects: Layer-independent replacement (Rule 614)
// - Prevention effects: Applied after replacements (Rule 615)
// - State-based actions: Lethal damage checking (Rule 704.5g)
// - Triggered abilities: "Whenever X deals combat damage" triggers
// - Layer system: Does not use layers (instant replacement at damage time)

// Important rules:
// - Rule 510.1: Combat damage step mechanics
// - Rule 510.1c: Attacker divides damage among blockers
// - Rule 510.1d: Blocker divides damage among attackers
// - Rule 614: Replacement effects (applied before prevention)
// - Rule 615: Prevention effects (applied after replacement)
// - Rule 702.2: Deathtouch (any amount is lethal)
// - Rule 702.19: Trample (excess tramples through)
// - Rule 702.4: Double strike (deals damage twice)
// - Rule 702.7: First strike (separate damage step)

// Example usage for preventing combat damage (Fog):
// 1. Cast Fog (instant): "Prevent all combat damage that would be dealt this turn"
// 2. Create PreventCombatDamageEffect with amount=-1, duration=UntilEndOfTurn
// 3. Register as prevention effect in game's replacement effect system
// 4. During combat damage step, effect prevents all damage
// 5. Effect expires at end of turn

// Example usage for trample + deathtouch:
// 1. Attacker has 5 power, trample, deathtouch
// 2. Blocked by two 3/3 creatures
// 3. Create TrampleDamageCalculator with hasDeathtouch=true
// 4. AddBlocker for each 3/3 (lethal=1 due to deathtouch)
// 5. CalculateTrampleDamage returns: {blocker1: 1, blocker2: 1}, trample: 3
// 6. Player can assign 1 to each blocker, 3 tramples through

// Example usage for damage doubling (Furnace of Rath):
// 1. Furnace of Rath: "If a source would deal damage, it deals double that damage instead"
// 2. Create ReplaceCombatDamageEffect with multiplier=2.0
// 3. Register as replacement effect
// 4. All combat damage events are doubled before being dealt
// 5. Works with trample, lifelink, deathtouch, etc.
