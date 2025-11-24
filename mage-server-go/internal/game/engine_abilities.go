package game

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap"
)

// AbilityActivationManager handles player activation of abilities
type AbilityActivationManager struct {
	priorityMgr     *PriorityManager
	stackMgr        *EnhancedStackManager
	abilityRegistry *AbilityRegistry
	logger          *zap.Logger
}

// NewAbilityActivationManager creates a new ability activation manager
func NewAbilityActivationManager(
	priorityMgr *PriorityManager,
	stackMgr *EnhancedStackManager,
	abilityRegistry *AbilityRegistry,
	logger *zap.Logger,
) *AbilityActivationManager {
	return &AbilityActivationManager{
		priorityMgr:     priorityMgr,
		stackMgr:        stackMgr,
		abilityRegistry: abilityRegistry,
		logger:          logger,
	}
}

// ActivateAbility activates an activated ability for a player
// This implements the full activation sequence from MTG Rules 602
func (aam *AbilityActivationManager) ActivateAbility(
	ctx context.Context,
	gameID uuid.UUID,
	playerID uuid.UUID,
	abilityID uuid.UUID,
	targets []uuid.UUID,
	state rules.GameStateReader,
	gameCtx abilities.GameContext,
) error {
	aam.logger.Info("activating ability",
		zap.String("game", gameID.String()),
		zap.String("player", playerID.String()),
		zap.String("ability", abilityID.String()))

	// 602.1a: The player announces that they are activating the ability
	// (implicit - this function call represents the announcement)

	// 602.1b: The ability is moved to the stack
	// TODO: Get the ability from the permanent/card
	// For now, we assume the ability is passed in or retrieved from game state

	// 602.2: The player checks timing restrictions
	// 602.2a: Check if the player has priority
	if aam.priorityMgr.GetPriorityPlayer() != playerID.String() {
		return fmt.Errorf("player does not have priority")
	}

	// 602.2b: Check timing restrictions based on ability type
	// - Activated abilities can only be played when the player could cast a sorcery
	//   UNLESS they are mana abilities
	// - Mana abilities can be activated any time the player has priority
	// TODO: Implement timing restriction checks

	// 602.3: The player chooses modes (if the ability is modal)
	// TODO: Implement mode selection

	// 602.3a: The player chooses targets for the ability
	// This is passed in via the targets parameter
	// TODO: Validate targets are legal

	// 602.3b: The player chooses how to divide effects among targets
	// TODO: Implement division logic (e.g., "distribute 3 damage among any number of targets")

	// 602.3c: The player makes other choices required by the ability
	// TODO: Implement choice making

	// 602.3d: The player determines the total cost of the ability
	// TODO: Get ability costs and calculate total cost with cost increases/reductions

	// 602.3e-602.3i: Advanced cost modifications
	// TODO: Implement cost modifications, alternative costs, additional costs

	// 602.2: Determine the ability can be activated
	// This includes checking if costs can be paid
	ability, err := aam.getAbility(ctx, abilityID, state)
	if err != nil {
		return fmt.Errorf("failed to get ability: %w", err)
	}

	// Check if ability can be activated (timing, costs, etc.)
	if !ability.CanActivate(ctx, gameCtx) {
		return fmt.Errorf("ability cannot be activated at this time")
	}

	// 602.2f: The player pays all costs
	// Costs are paid in this order:
	// 1. Mana costs
	// 2. Tap/untap costs
	// 3. All other costs in any order
	if err := aam.payCosts(ctx, ability, playerID, gameCtx); err != nil {
		return fmt.Errorf("failed to pay costs: %w", err)
	}

	aam.logger.Debug("costs paid successfully")

	// 602.2g: The ability becomes activated and goes on the stack
	choices := make(map[string]interface{})
	// TODO: Populate choices from modal selection, X values, etc.

	stackID, err := aam.stackMgr.PushActivatedAbility(ctx, playerID, ability, targets, choices)
	if err != nil {
		return fmt.Errorf("failed to push ability to stack: %w", err)
	}

	aam.logger.Info("ability added to stack",
		zap.String("ability", abilityID.String()),
		zap.String("stack_id", stackID.String()))

	// 602.2h: The ability can be responded to
	// Priority is given to the active player after the ability is on the stack
	// This is handled by the priority system

	// Reset priority pass counter since something was added to the stack
	aam.priorityMgr.ResetPassCount()

	// Give priority to the active player
	activePlayer := aam.priorityMgr.GetActivePlayer()
	if err := aam.priorityMgr.GivePriority(ctx, activePlayer, state); err != nil {
		return fmt.Errorf("failed to give priority: %w", err)
	}

	aam.logger.Info("ability activated successfully",
		zap.String("ability", abilityID.String()))

	return nil
}

// CastSpell casts a spell (instant or sorcery)
// This implements the full spell casting sequence from MTG Rules 601
// chosenModes: for modal spells, the mode IDs chosen (nil for non-modal)
func (aam *AbilityActivationManager) CastSpell(
	ctx context.Context,
	gameID uuid.UUID,
	playerID uuid.UUID,
	cardID uuid.UUID,
	targets []uuid.UUID,
	chosenModes []int,
	state rules.GameStateReader,
	gameCtx abilities.GameContext,
) error {
	aam.logger.Info("casting spell",
		zap.String("game", gameID.String()),
		zap.String("player", playerID.String()),
		zap.String("card", cardID.String()))

	// 601.2: To cast a spell, a player follows these steps

	// 601.2a: Get the card from game state
	cardInterface, err := gameCtx.GetCard(cardID)
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	card, ok := cardInterface.(*Card)
	if !ok {
		return fmt.Errorf("invalid card type")
	}

	// Verify the card is a spell (instant or sorcery)
	if !card.IsInstant() && !card.IsSorcery() {
		return fmt.Errorf("card is not a spell")
	}

	// 601.2b: Check if the player has priority
	if aam.priorityMgr.GetPriorityPlayer() != playerID.String() {
		return fmt.Errorf("player does not have priority")
	}

	// 601.2c: Check timing restrictions
	// - Instants can be cast any time the player has priority
	// - Sorceries can only be cast during main phase when stack is empty
	if card.IsSorcery() {
		if err := aam.checkSorceryTiming(); err != nil {
			return fmt.Errorf("cannot cast sorcery: %w", err)
		}
	}

	// Get the spell ability from the card
	// Spells have exactly one spell ability (the main ability that resolves)
	spellAbility, err := aam.getSpellAbility(card)
	if err != nil {
		return fmt.Errorf("failed to get spell ability: %w", err)
	}

	// 601.2d: The player chooses modes (if the spell is modal)
	choices := make(map[string]interface{})

	// Check if this is a modal spell
	if modalSpell, ok := spellAbility.(*abilities.ModalSpellAbility); ok {
		// Validate and set chosen modes
		if chosenModes == nil || len(chosenModes) == 0 {
			return fmt.Errorf("modal spell requires mode selection")
		}

		// Validate chosen modes
		if err := modalSpell.SetChosenModes(chosenModes); err != nil {
			return fmt.Errorf("invalid mode selection: %w", err)
		}

		// Check that all chosen modes can be selected
		availableModes := modalSpell.GetAvailableModes(ctx, gameCtx)
		availableMap := make(map[int]bool)
		for _, mode := range availableModes {
			availableMap[mode.ID] = true
		}

		for _, modeID := range chosenModes {
			if !availableMap[modeID] {
				return fmt.Errorf("mode %d cannot be chosen at this time", modeID)
			}
		}

		aam.logger.Debug("modes chosen for spell",
			zap.String("card", cardID.String()),
			zap.Ints("modes", chosenModes))

		// Store chosen modes in choices map for stack
		choices["chosen_modes"] = chosenModes
	} else if chosenModes != nil && len(chosenModes) > 0 {
		// Non-modal spell but modes were provided
		return fmt.Errorf("spell is not modal but modes were provided")
	}

	// 601.2e: The player chooses targets
	// Validate targets are legal
	if spellAbility.GetType() == abilities.AbilityTypeSpell {
		if sa, ok := spellAbility.(*abilities.SpellAbility); ok {
			if sa.GetTargets() != nil {
				// TODO: Validate targets using TargetSelectionManager
				// For now, we accept the provided targets
			}
		}
	}

	// 601.2f: The player determines the total cost
	// TODO: Calculate total cost with modifications (cost increases/reductions)

	// 601.2g-601.2i: Cost modifications, alternative costs, additional costs
	// TODO: Implement cost modifications

	// 601.2h: The player pays all costs
	// Mana costs first, then other costs
	if err := aam.paySpellCosts(ctx, card, spellAbility, playerID, gameCtx); err != nil {
		return fmt.Errorf("failed to pay costs: %w", err)
	}

	// Move card to stack zone
	// This is done by the stack manager when pushing the spell
	card.Zone = ZoneStack

	// 601.2i: The spell becomes cast and goes on the stack
	stackID, err := aam.stackMgr.PushSpell(ctx, cardID, playerID, spellAbility, targets, choices)
	if err != nil {
		return fmt.Errorf("failed to push spell to stack: %w", err)
	}

	aam.logger.Info("spell added to stack",
		zap.String("card", cardID.String()),
		zap.String("stack_id", stackID.String()))

	// Reset priority pass counter (something was added to stack)
	aam.priorityMgr.ResetPassCount()

	// Give priority to active player
	activePlayer := aam.priorityMgr.GetActivePlayer()
	if err := aam.priorityMgr.GivePriority(ctx, activePlayer, state); err != nil {
		return fmt.Errorf("failed to give priority: %w", err)
	}

	aam.logger.Info("spell cast successfully",
		zap.String("card", cardID.String()))

	return nil
}

// payCosts pays all costs for an ability in the correct order
func (aam *AbilityActivationManager) payCosts(
	ctx context.Context,
	ability abilities.Ability,
	playerID uuid.UUID,
	gameCtx abilities.GameContext,
) error {
	// Get the ability costs
	// For activated abilities, costs are stored in the ability
	activatedAbility, ok := ability.(*abilities.ActivatedAbility)
	if !ok {
		return fmt.Errorf("ability is not an activated ability")
	}

	costs := activatedAbility.GetCosts()

	// Pay costs in the correct order (Rule 118.8):
	// 1. Mana costs first
	// 2. Then all other costs in any order

	// Separate mana costs from other costs
	var manaCosts []abilities.Cost
	var otherCosts []abilities.Cost

	for _, cost := range costs {
		if _, ok := cost.(*abilities.ManaCost); ok {
			manaCosts = append(manaCosts, cost)
		} else {
			otherCosts = append(otherCosts, cost)
		}
	}

	// Pay mana costs first
	for _, cost := range manaCosts {
		aam.logger.Debug("paying mana cost",
			zap.String("cost", cost.String()))

		if err := cost.Pay(ctx, gameCtx, playerID); err != nil {
			return fmt.Errorf("failed to pay mana cost: %w", err)
		}
	}

	// Pay other costs
	for _, cost := range otherCosts {
		aam.logger.Debug("paying cost",
			zap.String("cost", cost.String()))

		if err := cost.Pay(ctx, gameCtx, playerID); err != nil {
			return fmt.Errorf("failed to pay cost %s: %w", cost.String(), err)
		}
	}

	return nil
}

// getAbility retrieves an ability from the ability registry
func (aam *AbilityActivationManager) getAbility(
	ctx context.Context,
	abilityID uuid.UUID,
	state rules.GameStateReader,
) (abilities.Ability, error) {
	// Retrieve the ability from the registry
	ability, err := aam.abilityRegistry.GetAbility(abilityID)
	if err != nil {
		return nil, fmt.Errorf("ability not found in registry: %w", err)
	}

	return ability, nil
}

// ActivateManaAbility activates a mana ability (special case - doesn't use stack)
// MTG Rules 605: Mana abilities don't use the stack and resolve immediately
func (aam *AbilityActivationManager) ActivateManaAbility(
	ctx context.Context,
	gameID uuid.UUID,
	playerID uuid.UUID,
	abilityID uuid.UUID,
	state rules.GameStateReader,
	gameCtx abilities.GameContext,
) error {
	aam.logger.Info("activating mana ability",
		zap.String("game", gameID.String()),
		zap.String("player", playerID.String()),
		zap.String("ability", abilityID.String()))

	// Get the mana ability
	ability, err := aam.getAbility(ctx, abilityID, state)
	if err != nil {
		return fmt.Errorf("failed to get ability: %w", err)
	}

	// Verify this is a mana ability
	if ability.GetType() != abilities.AbilityTypeMana {
		return fmt.Errorf("ability is not a mana ability")
	}

	// Check if costs can be paid
	if !ability.CanActivate(ctx, gameCtx) {
		return fmt.Errorf("mana ability cannot be activated")
	}

	// Pay costs
	if err := aam.payCosts(ctx, ability, playerID, gameCtx); err != nil {
		return fmt.Errorf("failed to pay costs: %w", err)
	}

	// Resolve the mana ability immediately (doesn't use stack)
	if err := ability.Resolve(ctx, gameCtx); err != nil {
		return fmt.Errorf("failed to resolve mana ability: %w", err)
	}

	aam.logger.Info("mana ability activated and resolved",
		zap.String("ability", abilityID.String()))

	return nil
}

// CheckActivationRestrictions checks if an ability can be activated
// based on timing restrictions, costs, and game state
func (aam *AbilityActivationManager) CheckActivationRestrictions(
	ctx context.Context,
	ability abilities.Ability,
	playerID uuid.UUID,
	gameCtx abilities.GameContext,
) error {
	// Check priority
	if aam.priorityMgr.GetPriorityPlayer() != playerID.String() {
		return fmt.Errorf("player does not have priority")
	}

	// Check ability type and timing restrictions
	switch ability.GetType() {
	case abilities.AbilityTypeMana:
		// Mana abilities can be activated any time player has priority
		return nil

	case abilities.AbilityTypeActivated:
		// Activated abilities require sorcery timing unless they specify otherwise
		// Check if it's a main phase and stack is empty
		currentStep := aam.priorityMgr.GetCurrentStep()
		activePlayer := aam.priorityMgr.GetActivePlayer()

		// Sorcery-speed abilities can only be activated during main phases
		// when stack is empty and it's the active player's turn
		if currentStep != rules.StepMain1 && currentStep != rules.StepMain2 {
			return fmt.Errorf("can only activate during main phase")
		}

		if activePlayer != playerID.String() {
			return fmt.Errorf("can only activate during your turn")
		}

		// TODO: Check if stack is empty

		return nil

	default:
		return fmt.Errorf("ability type cannot be activated by player")
	}
}

// getSpellAbility retrieves the spell ability from a card
// Instants and sorceries have exactly one spell ability
func (aam *AbilityActivationManager) getSpellAbility(card *Card) (abilities.Ability, error) {
	// Look through card's abilities for the spell ability
	for _, abilityInterface := range card.GetAbilities() {
		ability, ok := abilityInterface.(abilities.Ability)
		if !ok {
			continue
		}

		if ability.GetType() == abilities.AbilityTypeSpell {
			return ability, nil
		}
	}

	return nil, fmt.Errorf("card has no spell ability")
}

// checkSorceryTiming checks if it's valid timing to cast a sorcery
// Rule 307.1: Sorceries can only be cast during main phase when stack is empty
func (aam *AbilityActivationManager) checkSorceryTiming() error {
	currentStep := aam.priorityMgr.GetCurrentStep()

	// Must be during a main phase
	if currentStep != rules.StepMain1 && currentStep != rules.StepMain2 {
		return fmt.Errorf("can only cast sorcery during main phase")
	}

	// TODO: Check if stack is empty
	// stackSize := aam.stackMgr.Size()
	// if stackSize > 0 {
	// 	return fmt.Errorf("cannot cast sorcery while stack is not empty")
	// }

	// Must be active player's turn
	activePlayer := aam.priorityMgr.GetActivePlayer()
	priorityPlayer := aam.priorityMgr.GetPriorityPlayer()
	if activePlayer != priorityPlayer {
		return fmt.Errorf("can only cast sorcery during your turn")
	}

	return nil
}

// paySpellCosts pays the costs to cast a spell
// This is similar to payCosts but handles spell-specific cost payment
func (aam *AbilityActivationManager) paySpellCosts(
	ctx context.Context,
	card *Card,
	spellAbility abilities.Ability,
	playerID uuid.UUID,
	gameCtx abilities.GameContext,
) error {
	// Get the spell ability's mana cost
	var manaCost *abilities.ManaCost

	if sa, ok := spellAbility.(*abilities.SpellAbility); ok {
		manaCost = sa.GetManaCost()
	}

	// Pay mana cost if present
	if manaCost != nil {
		aam.logger.Debug("paying spell mana cost",
			zap.String("cost", manaCost.String()))

		// Check if player can pay the cost
		if !manaCost.CanPay(ctx, gameCtx, playerID) {
			return fmt.Errorf("cannot pay mana cost: insufficient mana")
		}

		// Pay the cost
		if err := manaCost.Pay(ctx, gameCtx, playerID); err != nil {
			return fmt.Errorf("failed to pay mana cost: %w", err)
		}
	}

	// TODO: Pay additional costs (sacrifice, discard, etc.)
	// These would come from the spell or from game effects

	return nil
}
