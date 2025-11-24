package game

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap"
)

// CombatIntegrationManager integrates the combat system with the abilities system
// This handles triggering combat-related abilities at the appropriate times
type CombatIntegrationManager struct {
	engine          *MageEngine
	abilityRegistry *AbilityRegistry
	stackMgr        *EnhancedStackManager
	logger          *zap.Logger
}

// NewCombatIntegrationManager creates a new combat integration manager
func NewCombatIntegrationManager(
	engine *MageEngine,
	abilityRegistry *AbilityRegistry,
	stackMgr *EnhancedStackManager,
	logger *zap.Logger,
) *CombatIntegrationManager {
	return &CombatIntegrationManager{
		engine:          engine,
		abilityRegistry: abilityRegistry,
		stackMgr:        stackMgr,
		logger:          logger,
	}
}

// OnDeclareAttackers is called when attackers are declared
// This triggers "When X attacks" abilities
func (cim *CombatIntegrationManager) OnDeclareAttackers(
	ctx context.Context,
	gameID string,
	attackers []string,
	state rules.GameStateReader,
) error {
	cim.logger.Debug("processing attacker declarations",
		zap.Int("attacker_count", len(attackers)))

	// For each attacking creature, check for attack triggers
	for _, attackerID := range attackers {
		if err := cim.triggerAttackAbilities(ctx, gameID, attackerID, state); err != nil {
			cim.logger.Error("failed to trigger attack abilities",
				zap.String("attacker", attackerID),
				zap.Error(err))
			return fmt.Errorf("failed to trigger attack abilities: %w", err)
		}
	}

	return nil
}

// OnDeclareBlockers is called when blockers are declared
// This triggers "When X blocks" and "When X becomes blocked" abilities
func (cim *CombatIntegrationManager) OnDeclareBlockers(
	ctx context.Context,
	gameID string,
	blocks map[string][]string, // attackerID -> []blockerID
	state rules.GameStateReader,
) error {
	cim.logger.Debug("processing blocker declarations",
		zap.Int("block_count", len(blocks)))

	// For each blocker, trigger block abilities
	for attackerID, blockers := range blocks {
		// Trigger "becomes blocked" on the attacker
		if len(blockers) > 0 {
			if err := cim.triggerBecomesBlockedAbilities(ctx, gameID, attackerID, state); err != nil {
				return fmt.Errorf("failed to trigger becomes blocked abilities: %w", err)
			}
		}

		// Trigger "blocks" on each blocker
		for _, blockerID := range blockers {
			if err := cim.triggerBlockAbilities(ctx, gameID, blockerID, attackerID, state); err != nil {
				return fmt.Errorf("failed to trigger block abilities: %w", err)
			}
		}
	}

	return nil
}

// OnCombatDamage is called when combat damage is dealt
// This triggers "When X deals combat damage" abilities
func (cim *CombatIntegrationManager) OnCombatDamage(
	ctx context.Context,
	gameID string,
	damageEvents []CombatDamageEvent,
	state rules.GameStateReader,
) error {
	cim.logger.Debug("processing combat damage",
		zap.Int("damage_event_count", len(damageEvents)))

	// For each damage event, trigger damage abilities
	for _, event := range damageEvents {
		if err := cim.triggerDamageAbilities(ctx, gameID, event, state); err != nil {
			return fmt.Errorf("failed to trigger damage abilities: %w", err)
		}
	}

	return nil
}

// OnCreatureDies is called when a creature dies in combat
// This triggers "When X dies" abilities
func (cim *CombatIntegrationManager) OnCreatureDies(
	ctx context.Context,
	gameID string,
	creatureID string,
	state rules.GameStateReader,
) error {
	cim.logger.Debug("processing creature death",
		zap.String("creature", creatureID))

	// Trigger dies abilities
	if err := cim.triggerDiesAbilities(ctx, gameID, creatureID, state); err != nil {
		return fmt.Errorf("failed to trigger dies abilities: %w", err)
	}

	return nil
}

// triggerAttackAbilities triggers "When X attacks" abilities
func (cim *CombatIntegrationManager) triggerAttackAbilities(
	ctx context.Context,
	gameID string,
	attackerID string,
	state rules.GameStateReader,
) error {
	// Get the attacker UUID
	attackerUUID, err := uuid.Parse(attackerID)
	if err != nil {
		return fmt.Errorf("invalid attacker ID: %w", err)
	}

	// Get all abilities for this permanent
	allAbilities := cim.abilityRegistry.GetAbilitiesBySource(attackerUUID)

	// Create the attack event
	gameEvent := abilities.GameEvent{
		Type:     abilities.EventAttackerDeclared,
		SourceID: attackerUUID,
	}

	// Filter for triggered abilities that trigger on attack
	for _, ability := range allAbilities {
		if ability.GetType() != abilities.AbilityTypeTriggered {
			continue
		}

		triggeredAbility, ok := ability.(*abilities.TriggeredAbility)
		if !ok {
			continue
		}

		// Check if this ability triggers on attack
		if triggeredAbility.CheckTrigger(gameEvent) {
			cim.logger.Info("triggering attack ability",
				zap.String("ability", ability.GetID().String()),
				zap.String("attacker", attackerID))

			// Get controller from metadata
			metadata, err := cim.abilityRegistry.GetMetadata(ability.GetID())
			if err != nil {
				cim.logger.Error("failed to get ability metadata",
					zap.String("ability", ability.GetID().String()),
					zap.Error(err))
				continue
			}

			// Collect targets if the ability requires them
			targets := []uuid.UUID{}
			// TODO: If ability has target requirements, prompt for targets
			// For now, triggered abilities with auto-targeting would work

			// Push triggered ability to stack
			if _, err := cim.stackMgr.PushTriggeredAbility(ctx, metadata.Controller, ability, targets); err != nil {
				cim.logger.Error("failed to push triggered ability to stack",
					zap.String("ability", ability.GetID().String()),
					zap.Error(err))
				return err
			}
		}
	}

	return nil
}

// triggerBecomesBlockedAbilities triggers "When X becomes blocked" abilities
func (cim *CombatIntegrationManager) triggerBecomesBlockedAbilities(
	ctx context.Context,
	gameID string,
	attackerID string,
	state rules.GameStateReader,
) error {
	// Get the attacker UUID
	attackerUUID, err := uuid.Parse(attackerID)
	if err != nil {
		return fmt.Errorf("invalid attacker ID: %w", err)
	}

	// Get all abilities for this permanent
	allAbilities := cim.abilityRegistry.GetAbilitiesBySource(attackerUUID)

	// Create the becomes blocked event
	gameEvent := abilities.GameEvent{
		Type:     abilities.EventBlockerDeclared,
		TargetID: attackerUUID, // The attacker is the target of the block
	}

	// Filter for "becomes blocked" triggers
	for _, ability := range allAbilities {
		if ability.GetType() != abilities.AbilityTypeTriggered {
			continue
		}

		triggeredAbility, ok := ability.(*abilities.TriggeredAbility)
		if !ok {
			continue
		}

		// Check if this ability triggers on becoming blocked
		if triggeredAbility.CheckTrigger(gameEvent) {
			cim.logger.Info("triggering becomes blocked ability",
				zap.String("ability", ability.GetID().String()),
				zap.String("attacker", attackerID))

			// Get controller from metadata
			metadata, err := cim.abilityRegistry.GetMetadata(ability.GetID())
			if err != nil {
				cim.logger.Error("failed to get ability metadata",
					zap.String("ability", ability.GetID().String()),
					zap.Error(err))
				continue
			}

			// Collect targets if needed
			targets := []uuid.UUID{}

			// Push triggered ability to stack
			if _, err := cim.stackMgr.PushTriggeredAbility(ctx, metadata.Controller, ability, targets); err != nil {
				cim.logger.Error("failed to push triggered ability to stack",
					zap.String("ability", ability.GetID().String()),
					zap.Error(err))
				return err
			}
		}
	}

	return nil
}

// triggerBlockAbilities triggers "When X blocks" abilities
func (cim *CombatIntegrationManager) triggerBlockAbilities(
	ctx context.Context,
	gameID string,
	blockerID string,
	attackerID string,
	state rules.GameStateReader,
) error {
	// Get the blocker UUID
	blockerUUID, err := uuid.Parse(blockerID)
	if err != nil {
		return fmt.Errorf("invalid blocker ID: %w", err)
	}

	// Get the attacker UUID
	attackerUUID, err := uuid.Parse(attackerID)
	if err != nil {
		return fmt.Errorf("invalid attacker ID: %w", err)
	}

	// Get all abilities for this permanent
	allAbilities := cim.abilityRegistry.GetAbilitiesBySource(blockerUUID)

	// Create the block event
	gameEvent := abilities.GameEvent{
		Type:     abilities.EventBlockerDeclared,
		SourceID: blockerUUID,
		TargetID: attackerUUID, // The attacker being blocked
	}

	// Filter for "blocks" triggers
	for _, ability := range allAbilities {
		if ability.GetType() != abilities.AbilityTypeTriggered {
			continue
		}

		triggeredAbility, ok := ability.(*abilities.TriggeredAbility)
		if !ok {
			continue
		}

		// Check if this ability triggers on blocking
		if triggeredAbility.CheckTrigger(gameEvent) {
			cim.logger.Info("triggering block ability",
				zap.String("ability", ability.GetID().String()),
				zap.String("blocker", blockerID),
				zap.String("attacker", attackerID))

			// Get controller from metadata
			metadata, err := cim.abilityRegistry.GetMetadata(ability.GetID())
			if err != nil {
				cim.logger.Error("failed to get ability metadata",
					zap.String("ability", ability.GetID().String()),
					zap.Error(err))
				continue
			}

			// Collect targets if needed
			targets := []uuid.UUID{}

			// Push triggered ability to stack
			if _, err := cim.stackMgr.PushTriggeredAbility(ctx, metadata.Controller, ability, targets); err != nil {
				cim.logger.Error("failed to push triggered ability to stack",
					zap.String("ability", ability.GetID().String()),
					zap.Error(err))
				return err
			}
		}
	}

	return nil
}

// triggerDamageAbilities triggers "When X deals damage" abilities
func (cim *CombatIntegrationManager) triggerDamageAbilities(
	ctx context.Context,
	gameID string,
	event CombatDamageEvent,
	state rules.GameStateReader,
) error {
	// Get the source UUID
	sourceUUID, err := uuid.Parse(event.SourceID)
	if err != nil {
		return fmt.Errorf("invalid source ID: %w", err)
	}

	// Get all abilities for this permanent
	allAbilities := cim.abilityRegistry.GetAbilitiesBySource(sourceUUID)

	// Filter for damage triggers
	for _, ability := range allAbilities {
		if ability.GetType() != abilities.AbilityTypeTriggered {
			continue
		}

		triggeredAbility, ok := ability.(*abilities.TriggeredAbility)
		if !ok {
			continue
		}

		// Check if this ability triggers on damage
		gameEvent := abilities.GameEvent{
			Type:     abilities.EventDamageDealt,
			SourceID: sourceUUID,
			TargetID: uuid.MustParse(event.TargetID),
			Amount:   event.Amount,
		}

		if triggeredAbility.CheckTrigger(gameEvent) {
			cim.logger.Info("triggering damage ability",
				zap.String("ability", ability.GetID().String()),
				zap.String("source", event.SourceID),
				zap.Int("damage", event.Amount))

			// Get controller from metadata
			metadata, err := cim.abilityRegistry.GetMetadata(ability.GetID())
			if err != nil {
				cim.logger.Error("failed to get ability metadata",
					zap.String("ability", ability.GetID().String()),
					zap.Error(err))
				continue
			}

			// Collect targets if needed
			targets := []uuid.UUID{}

			// Push triggered ability to stack
			if _, err := cim.stackMgr.PushTriggeredAbility(ctx, metadata.Controller, ability, targets); err != nil {
				cim.logger.Error("failed to push triggered ability to stack",
					zap.String("ability", ability.GetID().String()),
					zap.Error(err))
				return err
			}
		}
	}

	return nil
}

// triggerDiesAbilities triggers "When X dies" abilities
func (cim *CombatIntegrationManager) triggerDiesAbilities(
	ctx context.Context,
	gameID string,
	creatureID string,
	state rules.GameStateReader,
) error {
	// Get the creature UUID
	creatureUUID, err := uuid.Parse(creatureID)
	if err != nil {
		return fmt.Errorf("invalid creature ID: %w", err)
	}

	// Get all abilities for this permanent
	allAbilities := cim.abilityRegistry.GetAbilitiesBySource(creatureUUID)

	// Filter for dies triggers
	for _, ability := range allAbilities {
		if ability.GetType() != abilities.AbilityTypeTriggered {
			continue
		}

		triggeredAbility, ok := ability.(*abilities.TriggeredAbility)
		if !ok {
			continue
		}

		// Check if this ability triggers on death
		gameEvent := abilities.GameEvent{
			Type:     abilities.EventDies,
			SourceID: creatureUUID,
		}

		if triggeredAbility.CheckTrigger(gameEvent) {
			cim.logger.Info("triggering dies ability",
				zap.String("ability", ability.GetID().String()),
				zap.String("creature", creatureID))

			// Get controller from metadata
			metadata, err := cim.abilityRegistry.GetMetadata(ability.GetID())
			if err != nil {
				cim.logger.Error("failed to get ability metadata",
					zap.String("ability", ability.GetID().String()),
					zap.Error(err))
				continue
			}

			// Collect targets if needed
			targets := []uuid.UUID{}

			// Push triggered ability to stack
			if _, err := cim.stackMgr.PushTriggeredAbility(ctx, metadata.Controller, ability, targets); err != nil {
				cim.logger.Error("failed to push triggered ability to stack",
					zap.String("ability", ability.GetID().String()),
					zap.Error(err))
				return err
			}
		}
	}

	return nil
}

// CombatDamageEvent represents a combat damage event
type CombatDamageEvent struct {
	SourceID string
	TargetID string
	Amount   int
	IsCombat bool
}

// CheckCombatKeywordAbilities checks for keyword abilities that affect combat
// This includes flying, first strike, trample, etc.
func (cim *CombatIntegrationManager) CheckCombatKeywordAbilities(
	ctx context.Context,
	creatureID string,
	state rules.GameStateReader,
) (*CombatKeywordAbilities, error) {
	// Get the creature UUID
	creatureUUID, err := uuid.Parse(creatureID)
	if err != nil {
		return nil, fmt.Errorf("invalid creature ID: %w", err)
	}

	// Get all abilities for this permanent
	allAbilities := cim.abilityRegistry.GetAbilitiesBySource(creatureUUID)

	keywords := &CombatKeywordAbilities{}

	// Check for keyword abilities
	for _, ability := range allAbilities {
		if ability.GetType() != abilities.AbilityTypeKeyword {
			continue
		}

		// TODO: Extract keyword ability type and set flags
		// For now, we'll need to check the ability's string representation
		// This would be better with a proper keyword ability interface

		abilityText := ability.String()
		switch abilityText {
		case "Flying":
			keywords.Flying = true
		case "First strike":
			keywords.FirstStrike = true
		case "Double strike":
			keywords.DoubleStrike = true
		case "Deathtouch":
			keywords.Deathtouch = true
		case "Lifelink":
			keywords.Lifelink = true
		case "Trample":
			keywords.Trample = true
		case "Vigilance":
			keywords.Vigilance = true
		case "Menace":
			keywords.Menace = true
		case "Defender":
			keywords.Defender = true
		}
	}

	return keywords, nil
}

// CombatKeywordAbilities represents keyword abilities that affect combat
type CombatKeywordAbilities struct {
	Flying       bool
	FirstStrike  bool
	DoubleStrike bool
	Deathtouch   bool
	Lifelink     bool
	Trample      bool
	Vigilance    bool
	Menace       bool
	Defender     bool
}

// CanAttack checks if a creature can attack based on its abilities
func (cim *CombatIntegrationManager) CanAttack(
	ctx context.Context,
	creatureID string,
	state rules.GameStateReader,
) (bool, string, error) {
	keywords, err := cim.CheckCombatKeywordAbilities(ctx, creatureID, state)
	if err != nil {
		return false, "", err
	}

	// Defender can't attack
	if keywords.Defender {
		return false, "creature has defender", nil
	}

	// TODO: Check for other restrictions (summoning sickness, tapped, etc.)
	// These are already handled by the MageEngine, but we could add
	// ability-based restrictions here

	return true, "", nil
}

// CanBlock checks if a creature can block based on its abilities
func (cim *CombatIntegrationManager) CanBlock(
	ctx context.Context,
	blockerID string,
	attackerID string,
	state rules.GameStateReader,
) (bool, string, error) {
	// Get keyword abilities for both creatures
	blockerKeywords, err := cim.CheckCombatKeywordAbilities(ctx, blockerID, state)
	if err != nil {
		return false, "", err
	}

	attackerKeywords, err := cim.CheckCombatKeywordAbilities(ctx, attackerID, state)
	if err != nil {
		return false, "", err
	}

	// Flying can only be blocked by creatures with flying or reach
	if attackerKeywords.Flying {
		if !blockerKeywords.Flying {
			// TODO: Check for reach
			return false, "cannot block flying without flying or reach", nil
		}
	}

	return true, "", nil
}
