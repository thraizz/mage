package rules

import (
	"fmt"
	"sync"
)

// ManaAbilityActivationContext tracks mana ability activation state
// Per Rule 605.3c: Once a player begins to activate a mana ability,
// that ability can't be activated again until it has resolved.
type ManaAbilityActivationContext struct {
	mu                    sync.RWMutex
	activatingAbilities   map[string]bool // abilityID -> is currently activating
	resolvedThisWindow    map[string]int  // abilityID -> count of resolutions in current window
}

// NewManaAbilityActivationContext creates a new mana ability activation context
func NewManaAbilityActivationContext() *ManaAbilityActivationContext {
	return &ManaAbilityActivationContext{
		activatingAbilities: make(map[string]bool),
		resolvedThisWindow:  make(map[string]int),
	}
}

// BeginActivation marks the start of activating a mana ability
func (maac *ManaAbilityActivationContext) BeginActivation(abilityID string) error {
	maac.mu.Lock()
	defer maac.mu.Unlock()

	if maac.activatingAbilities[abilityID] {
		return fmt.Errorf("mana ability %s is already being activated", abilityID)
	}

	maac.activatingAbilities[abilityID] = true
	return nil
}

// EndActivation marks the end of activating a mana ability (it has resolved)
func (maac *ManaAbilityActivationContext) EndActivation(abilityID string) {
	maac.mu.Lock()
	defer maac.mu.Unlock()

	delete(maac.activatingAbilities, abilityID)
	maac.resolvedThisWindow[abilityID]++
}

// CanActivate returns true if the mana ability can be activated
func (maac *ManaAbilityActivationContext) CanActivate(abilityID string) bool {
	maac.mu.RLock()
	defer maac.mu.RUnlock()

	return !maac.activatingAbilities[abilityID]
}

// IsActivating returns true if the mana ability is currently activating
func (maac *ManaAbilityActivationContext) IsActivating(abilityID string) bool {
	maac.mu.RLock()
	defer maac.mu.RUnlock()

	return maac.activatingAbilities[abilityID]
}

// GetResolutionCount returns how many times an ability has resolved in this window
func (maac *ManaAbilityActivationContext) GetResolutionCount(abilityID string) int {
	maac.mu.RLock()
	defer maac.mu.RUnlock()

	return maac.resolvedThisWindow[abilityID]
}

// ResetWindow clears the resolution tracking for a new window
func (maac *ManaAbilityActivationContext) ResetWindow() {
	maac.mu.Lock()
	defer maac.mu.Unlock()

	maac.resolvedThisWindow = make(map[string]int)
}

// Reset clears all state
func (maac *ManaAbilityActivationContext) Reset() {
	maac.mu.Lock()
	defer maac.mu.Unlock()

	maac.activatingAbilities = make(map[string]bool)
	maac.resolvedThisWindow = make(map[string]int)
}

// ManaAbility represents a mana ability that can be activated
type ManaAbility struct {
	ID           string
	SourceID     string // Card/permanent that has this ability
	ControllerID string // Player who controls the source
	Text         string // Ability text
	Activate     func() error // Function to execute the ability
}

// TriggeredManaAbility represents a triggered mana ability
// Per Rule 605.4a: Triggered mana abilities don't go on stack,
// they resolve immediately after the mana ability that triggered them
type TriggeredManaAbility struct {
	ID           string
	SourceID     string
	ControllerID string
	TriggerID    string // ID of the mana ability that triggered this
	Text         string
	Resolve      func() error
}

// ManaAbilityManager manages mana ability activation and triggered mana abilities
type ManaAbilityManager struct {
	mu                     sync.RWMutex
	activationContext      *ManaAbilityActivationContext
	triggeredQueue         []TriggeredManaAbility // Queue of triggered mana abilities waiting to resolve
	canActivateDuringCast  bool                    // Whether mana abilities can be activated during casting
	canActivateDuringResolve bool                  // Whether mana abilities can be activated during resolution
}

// NewManaAbilityManager creates a new mana ability manager
func NewManaAbilityManager() *ManaAbilityManager {
	return &ManaAbilityManager{
		activationContext:      NewManaAbilityActivationContext(),
		triggeredQueue:         make([]TriggeredManaAbility, 0, 8),
		canActivateDuringCast:  true,  // Rule 117.1d - can activate when casting
		canActivateDuringResolve: false, // Generally false unless in payment window
	}
}

// ActivateManaAbility activates a mana ability
// Per Rule 605.3b: Mana abilities don't go on the stack, they resolve immediately
func (mam *ManaAbilityManager) ActivateManaAbility(ability ManaAbility) error {
	mam.mu.Lock()
	defer mam.mu.Unlock()

	// Check if ability can be activated (Rule 605.3c)
	if !mam.activationContext.CanActivate(ability.ID) {
		return fmt.Errorf("mana ability %s cannot be activated (already activating)", ability.ID)
	}

	// Mark as activating
	if err := mam.activationContext.BeginActivation(ability.ID); err != nil {
		return err
	}

	// Execute the ability immediately (doesn't go on stack)
	if err := ability.Activate(); err != nil {
		// Failed to activate, mark as no longer activating
		mam.activationContext.EndActivation(ability.ID)
		return fmt.Errorf("failed to activate mana ability %s: %w", ability.ID, err)
	}

	// Mark as resolved
	mam.activationContext.EndActivation(ability.ID)

	return nil
}

// QueueTriggeredManaAbility adds a triggered mana ability to the queue
// Per Rule 605.4a: Triggered mana abilities resolve immediately
func (mam *ManaAbilityManager) QueueTriggeredManaAbility(ability TriggeredManaAbility) {
	mam.mu.Lock()
	defer mam.mu.Unlock()

	mam.triggeredQueue = append(mam.triggeredQueue, ability)
}

// ResolveTriggeredManaAbilities resolves all queued triggered mana abilities
// Per Rule 605.4a: They resolve immediately after the triggering mana ability
func (mam *ManaAbilityManager) ResolveTriggeredManaAbilities() error {
	mam.mu.Lock()
	defer mam.mu.Unlock()

	for len(mam.triggeredQueue) > 0 {
		// Take the first triggered ability
		ability := mam.triggeredQueue[0]
		mam.triggeredQueue = mam.triggeredQueue[1:]

		// Resolve it immediately
		if err := ability.Resolve(); err != nil {
			return fmt.Errorf("failed to resolve triggered mana ability %s: %w", ability.ID, err)
		}

		// Note: Resolving a triggered mana ability might trigger more mana abilities
		// The loop will continue until all are resolved
	}

	return nil
}

// HasPendingTriggeredAbilities returns true if there are triggered mana abilities waiting
func (mam *ManaAbilityManager) HasPendingTriggeredAbilities() bool {
	mam.mu.RLock()
	defer mam.mu.RUnlock()

	return len(mam.triggeredQueue) > 0
}

// SetCanActivateDuringCast sets whether mana abilities can be activated during casting
func (mam *ManaAbilityManager) SetCanActivateDuringCast(can bool) {
	mam.mu.Lock()
	defer mam.mu.Unlock()
	mam.canActivateDuringCast = can
}

// SetCanActivateDuringResolve sets whether mana abilities can be activated during resolution
func (mam *ManaAbilityManager) SetCanActivateDuringResolve(can bool) {
	mam.mu.Lock()
	defer mam.mu.Unlock()
	mam.canActivateDuringResolve = can
}

// CanActivate returns true if mana abilities can currently be activated
func (mam *ManaAbilityManager) CanActivate(duringCast bool, duringResolve bool) bool {
	mam.mu.RLock()
	defer mam.mu.RUnlock()

	if duringCast {
		return mam.canActivateDuringCast
	}
	if duringResolve {
		return mam.canActivateDuringResolve
	}
	return true // Can always activate when player has priority (Rule 117.1d)
}

// ResetWindow resets tracking for a new payment window
func (mam *ManaAbilityManager) ResetWindow() {
	mam.mu.Lock()
	defer mam.mu.Unlock()
	mam.activationContext.ResetWindow()
}

// Reset clears all state
func (mam *ManaAbilityManager) Reset() {
	mam.mu.Lock()
	defer mam.mu.Unlock()
	mam.activationContext.Reset()
	mam.triggeredQueue = mam.triggeredQueue[:0]
	mam.canActivateDuringCast = true
	mam.canActivateDuringResolve = false
}
