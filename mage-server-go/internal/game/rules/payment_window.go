package rules

import (
	"fmt"
	"sync"
)

// PaymentStep represents a step in paying costs
type PaymentStep string

const (
	// PaymentStepBefore allows special mana abilities before normal mana (e.g., Convoke)
	PaymentStepBefore PaymentStep = "BEFORE"
	// PaymentStepNormal allows normal mana ability activation
	PaymentStepNormal PaymentStep = "NORMAL"
	// PaymentStepAfter is after special mana abilities, normal mana blocked
	PaymentStepAfter PaymentStep = "AFTER"
)

// CostType represents different types of costs that can be paid
type CostType string

const (
	// CostTypeMana is mana cost
	CostTypeMana CostType = "MANA"
	// CostTypeTap is tapping cost
	CostTypeTap CostType = "TAP"
	// CostTypeSacrifice is sacrificing permanents
	CostTypeSacrifice CostType = "SACRIFICE"
	// CostTypeDiscard is discarding cards
	CostTypeDiscard CostType = "DISCARD"
	// CostTypeLife is paying life
	CostTypeLife CostType = "LIFE"
	// CostTypeExile is exiling cards
	CostTypeExile CostType = "EXILE"
	// CostTypeOther is other costs
	CostTypeOther CostType = "OTHER"
)

// Cost represents a cost that must be paid
type Cost struct {
	Type        CostType
	Amount      int    // Numeric amount (for mana, life, etc.)
	Description string // Human-readable description
	Paid        bool   // Whether this cost has been paid
}

// PaymentState tracks the state of paying costs for a spell/ability
type PaymentState struct {
	mu                sync.RWMutex
	spellOrAbilityID  string
	controllerID      string
	costs             []Cost
	currentStep       PaymentStep
	allowManaAbilities bool
	allowSpecialActions bool
	paidCosts         map[CostType]bool
	totalManaPaid     int
	manaRemaining     int
}

// NewPaymentState creates a new payment state
func NewPaymentState(spellOrAbilityID string, controllerID string, costs []Cost) *PaymentState {
	return &PaymentState{
		spellOrAbilityID:   spellOrAbilityID,
		controllerID:       controllerID,
		costs:              costs,
		currentStep:        PaymentStepBefore,
		allowManaAbilities: true,
		allowSpecialActions: true,
		paidCosts:          make(map[CostType]bool),
		totalManaPaid:      0,
		manaRemaining:      0,
	}
}

// GetCurrentStep returns the current payment step
func (ps *PaymentState) GetCurrentStep() PaymentStep {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.currentStep
}

// SetCurrentStep sets the current payment step
func (ps *PaymentState) SetCurrentStep(step PaymentStep) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.currentStep = step

	// Update what's allowed based on step
	// Per Java implementation: after special mana abilities, normal mana is blocked
	if step == PaymentStepAfter {
		ps.allowManaAbilities = false
	}
}

// CanActivateManaAbilities returns whether mana abilities can be activated
func (ps *PaymentState) CanActivateManaAbilities() bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.allowManaAbilities
}

// CanTakeSpecialActions returns whether special actions can be taken
func (ps *PaymentState) CanTakeSpecialActions() bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.allowSpecialActions
}

// MarkCostPaid marks a cost as paid
func (ps *PaymentState) MarkCostPaid(costType CostType) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.paidCosts[costType] = true
}

// IsCostPaid returns whether a cost has been paid
func (ps *PaymentState) IsCostPaid(costType CostType) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.paidCosts[costType]
}

// AddManaPaid adds to the total mana paid
func (ps *PaymentState) AddManaPaid(amount int) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.totalManaPaid += amount
	if ps.manaRemaining >= amount {
		ps.manaRemaining -= amount
	}
}

// GetTotalManaPaid returns the total mana paid
func (ps *PaymentState) GetTotalManaPaid() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.totalManaPaid
}

// SetManaRemaining sets the remaining mana to pay
func (ps *PaymentState) SetManaRemaining(amount int) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.manaRemaining = amount
}

// GetManaRemaining returns the remaining mana to pay
func (ps *PaymentState) GetManaRemaining() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.manaRemaining
}

// IsFullyPaid returns whether all costs have been paid
func (ps *PaymentState) IsFullyPaid() bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, cost := range ps.costs {
		if !cost.Paid {
			return false
		}
	}
	return ps.manaRemaining <= 0
}

// GetCosts returns a copy of all costs
func (ps *PaymentState) GetCosts() []Cost {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	costs := make([]Cost, len(ps.costs))
	copy(costs, ps.costs)
	return costs
}

// PaymentWindowManager manages payment windows during spell/ability casting
type PaymentWindowManager struct {
	mu               sync.RWMutex
	activePayment    *PaymentState
	paymentHistory   []string // IDs of spells/abilities paid for
}

// NewPaymentWindowManager creates a new payment window manager
func NewPaymentWindowManager() *PaymentWindowManager {
	return &PaymentWindowManager{
		paymentHistory: make([]string, 0, 16),
	}
}

// BeginPayment starts a new payment window
func (pwm *PaymentWindowManager) BeginPayment(state *PaymentState) error {
	pwm.mu.Lock()
	defer pwm.mu.Unlock()

	if pwm.activePayment != nil {
		return fmt.Errorf("payment already in progress for %s", pwm.activePayment.spellOrAbilityID)
	}

	pwm.activePayment = state
	return nil
}

// EndPayment closes the current payment window
func (pwm *PaymentWindowManager) EndPayment(spellOrAbilityID string) error {
	pwm.mu.Lock()
	defer pwm.mu.Unlock()

	if pwm.activePayment == nil {
		return fmt.Errorf("no payment in progress")
	}

	if pwm.activePayment.spellOrAbilityID != spellOrAbilityID {
		return fmt.Errorf("payment mismatch: expected %s, got %s",
			pwm.activePayment.spellOrAbilityID, spellOrAbilityID)
	}

	pwm.paymentHistory = append(pwm.paymentHistory, spellOrAbilityID)
	pwm.activePayment = nil
	return nil
}

// GetActivePayment returns the current payment state
func (pwm *PaymentWindowManager) GetActivePayment() *PaymentState {
	pwm.mu.RLock()
	defer pwm.mu.RUnlock()
	return pwm.activePayment
}

// IsPaymentInProgress returns true if a payment is in progress
func (pwm *PaymentWindowManager) IsPaymentInProgress() bool {
	pwm.mu.RLock()
	defer pwm.mu.RUnlock()
	return pwm.activePayment != nil
}

// Reset clears all payment state
func (pwm *PaymentWindowManager) Reset() {
	pwm.mu.Lock()
	defer pwm.mu.Unlock()
	pwm.activePayment = nil
	pwm.paymentHistory = pwm.paymentHistory[:0]
}

// ChoiceType represents types of choices during resolution
type ChoiceType string

const (
	// ChoiceTypeMode chooses mode for modal spells/abilities
	ChoiceTypeMode ChoiceType = "MODE"
	// ChoiceTypeTarget chooses targets
	ChoiceTypeTarget ChoiceType = "TARGET"
	// ChoiceTypeX chooses value for X
	ChoiceTypeX ChoiceType = "X_VALUE"
	// ChoiceTypeColor chooses a color
	ChoiceTypeColor ChoiceType = "COLOR"
	// ChoiceTypeNumber chooses a number
	ChoiceTypeNumber ChoiceType = "NUMBER"
	// ChoiceTypeYesNo chooses yes/no
	ChoiceTypeYesNo ChoiceType = "YES_NO"
	// ChoiceTypeCard chooses a card
	ChoiceTypeCard ChoiceType = "CARD"
	// ChoiceTypePlayer chooses a player
	ChoiceTypePlayer ChoiceType = "PLAYER"
	// ChoiceTypeOther other choice type
	ChoiceTypeOther ChoiceType = "OTHER"
)

// Choice represents a choice that must be made during resolution
type Choice struct {
	Type        ChoiceType
	PlayerID    string
	Prompt      string
	Options     []string
	MinChoices  int
	MaxChoices  int
	Result      []string // Chosen options
	Made        bool
}

// ChoiceManager manages choices during resolution
// Per Rule 608.2: Choices are made during resolution
type ChoiceManager struct {
	mu              sync.RWMutex
	pendingChoices  []Choice
	madeChoices     []Choice
	currentChoice   *Choice
}

// NewChoiceManager creates a new choice manager
func NewChoiceManager() *ChoiceManager {
	return &ChoiceManager{
		pendingChoices: make([]Choice, 0, 8),
		madeChoices:    make([]Choice, 0, 16),
	}
}

// AddChoice adds a pending choice
func (cm *ChoiceManager) AddChoice(choice Choice) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.pendingChoices = append(cm.pendingChoices, choice)
}

// GetNextChoice returns the next pending choice
func (cm *ChoiceManager) GetNextChoice() *Choice {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if len(cm.pendingChoices) == 0 {
		return nil
	}

	cm.currentChoice = &cm.pendingChoices[0]
	cm.pendingChoices = cm.pendingChoices[1:]
	return cm.currentChoice
}

// MakeChoice records a choice
func (cm *ChoiceManager) MakeChoice(result []string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.currentChoice == nil {
		return fmt.Errorf("no choice in progress")
	}

	// Validate choice count
	if len(result) < cm.currentChoice.MinChoices {
		return fmt.Errorf("too few choices (need at least %d)", cm.currentChoice.MinChoices)
	}
	if cm.currentChoice.MaxChoices > 0 && len(result) > cm.currentChoice.MaxChoices {
		return fmt.Errorf("too many choices (max %d)", cm.currentChoice.MaxChoices)
	}

	cm.currentChoice.Result = result
	cm.currentChoice.Made = true
	cm.madeChoices = append(cm.madeChoices, *cm.currentChoice)
	cm.currentChoice = nil

	return nil
}

// HasPendingChoices returns true if there are pending choices
func (cm *ChoiceManager) HasPendingChoices() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.pendingChoices) > 0
}

// Reset clears all choice state
func (cm *ChoiceManager) Reset() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.pendingChoices = cm.pendingChoices[:0]
	cm.madeChoices = cm.madeChoices[:0]
	cm.currentChoice = nil
}
