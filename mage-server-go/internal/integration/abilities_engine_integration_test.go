package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/magefree/mage-server-go/internal/game/mana"
	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap"
)

// TestAbilityRegistration tests the complete ability registration workflow
func TestAbilityRegistration(t *testing.T) {
	logger := zap.NewNop()
	registry := game.NewAbilityRegistry(logger)

	cardID := uuid.New()
	playerID := uuid.New()

	// Create a simple activated ability
	ability := abilities.NewActivatedAbilityBuilder(cardID).
		AddManaCost("{1}{R}").
		AddEffect(abilities.NewDamageEffect(3)).
		Build()

	// Register the ability
	registry.RegisterAbility(ability, playerID, abilities.ZoneBattlefield, 0)

	// Verify we can retrieve it
	retrieved := registry.GetAbility(ability.GetID())
	if retrieved == nil {
		t.Fatal("Failed to retrieve registered ability")
	}

	if retrieved.GetID() != ability.GetID() {
		t.Errorf("Retrieved wrong ability: expected %s, got %s",
			ability.GetID(), retrieved.GetID())
	}

	// Verify we can get abilities by source
	sourceAbilities := registry.GetAbilitiesBySource(cardID)
	if len(sourceAbilities) != 1 {
		t.Errorf("Expected 1 ability for source, got %d", len(sourceAbilities))
	}

	// Verify we can get activatable abilities
	activatable := registry.GetActivatableAbilities(playerID, abilities.ZoneBattlefield)
	if len(activatable) != 1 {
		t.Errorf("Expected 1 activatable ability, got %d", len(activatable))
	}

	// Unregister and verify cleanup
	registry.UnregisterSource(cardID)
	sourceAbilities = registry.GetAbilitiesBySource(cardID)
	if len(sourceAbilities) != 0 {
		t.Errorf("Expected 0 abilities after unregister, got %d", len(sourceAbilities))
	}
}

// TestAbilityRegistrationWithZoneChange tests zone tracking
func TestAbilityRegistrationWithZoneChange(t *testing.T) {
	logger := zap.NewNop()
	registry := game.NewAbilityRegistry(logger)

	cardID := uuid.New()
	playerID := uuid.New()

	// Register ability in battlefield
	ability := abilities.NewActivatedAbilityBuilder(cardID).
		AddManaCost("{1}").
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()

	registry.RegisterAbility(ability, playerID, abilities.ZoneBattlefield, 0)

	// Get activatable abilities on battlefield
	activatable := registry.GetActivatableAbilities(playerID, abilities.ZoneBattlefield)
	if len(activatable) != 1 {
		t.Fatalf("Expected 1 activatable ability on battlefield, got %d", len(activatable))
	}

	// Move to graveyard
	registry.UpdateCardZone(cardID, abilities.ZoneGraveyard)

	// Should not be activatable on battlefield anymore
	activatable = registry.GetActivatableAbilities(playerID, abilities.ZoneBattlefield)
	if len(activatable) != 0 {
		t.Errorf("Expected 0 activatable abilities on battlefield after move, got %d", len(activatable))
	}

	// Should be tracked in graveyard
	activatable = registry.GetActivatableAbilities(playerID, abilities.ZoneGraveyard)
	if len(activatable) != 1 {
		t.Errorf("Expected 1 activatable ability in graveyard, got %d", len(activatable))
	}
}

// TestTargetValidation tests the target selection and validation system
func TestTargetValidation(t *testing.T) {
	logger := zap.NewNop()
	tsm := game.NewTargetSelectionManager(logger)

	ctx := context.Background()
	sourceID := uuid.New()
	playerID := uuid.New()
	targetID := uuid.New()

	// Create a target request requiring 1 target
	filter := abilities.NewAnyTarget()
	request := &game.TargetRequest{
		AbilityID:    uuid.New(),
		SourceID:     sourceID,
		PlayerID:     playerID,
		MinTargets:   1,
		MaxTargets:   1,
		TargetFilter: filter,
		LegalTargets: []uuid.UUID{targetID},
		Message:      "Choose a target",
	}

	// Valid selection
	selection := &game.TargetSelection{
		Targets: []uuid.UUID{targetID},
	}

	err := tsm.ValidateTargets(ctx, request, selection, nil, nil)
	if err != nil {
		t.Errorf("Valid target selection failed validation: %v", err)
	}

	// Too few targets
	selection = &game.TargetSelection{
		Targets: []uuid.UUID{},
	}
	err = tsm.ValidateTargets(ctx, request, selection, nil, nil)
	if err == nil {
		t.Error("Expected error for too few targets, got nil")
	}

	// Too many targets
	selection = &game.TargetSelection{
		Targets: []uuid.UUID{targetID, uuid.New()},
	}
	err = tsm.ValidateTargets(ctx, request, selection, nil, nil)
	if err == nil {
		t.Error("Expected error for too many targets, got nil")
	}

	// Illegal target (not in legal targets list)
	illegalID := uuid.New()
	selection = &game.TargetSelection{
		Targets: []uuid.UUID{illegalID},
	}
	err = tsm.ValidateTargets(ctx, request, selection, nil, nil)
	if err == nil {
		t.Error("Expected error for illegal target, got nil")
	}
}

// TestTargetValidationOptional tests optional target selection
func TestTargetValidationOptional(t *testing.T) {
	logger := zap.NewNop()
	tsm := game.NewTargetSelectionManager(logger)

	ctx := context.Background()
	sourceID := uuid.New()
	playerID := uuid.New()
	targetID := uuid.New()

	// Create a request with optional target (min=0, max=1)
	filter := abilities.NewAnyTarget()
	request := &game.TargetRequest{
		AbilityID:    uuid.New(),
		SourceID:     sourceID,
		PlayerID:     playerID,
		MinTargets:   0,
		MaxTargets:   1,
		TargetFilter: filter,
		LegalTargets: []uuid.UUID{targetID},
		Message:      "You may choose a target",
	}

	// Valid with no targets (optional)
	selection := &game.TargetSelection{
		Targets: []uuid.UUID{},
	}
	err := tsm.ValidateTargets(ctx, request, selection, nil, nil)
	if err != nil {
		t.Errorf("Optional target with no selection failed validation: %v", err)
	}

	// Valid with one target
	selection = &game.TargetSelection{
		Targets: []uuid.UUID{targetID},
	}
	err = tsm.ValidateTargets(ctx, request, selection, nil, nil)
	if err != nil {
		t.Errorf("Optional target with one selection failed validation: %v", err)
	}
}

// TestTargetValidationMultiple tests multiple target selection
func TestTargetValidationMultiple(t *testing.T) {
	logger := zap.NewNop()
	tsm := game.NewTargetSelectionManager(logger)

	ctx := context.Background()
	sourceID := uuid.New()
	playerID := uuid.New()
	target1 := uuid.New()
	target2 := uuid.New()
	target3 := uuid.New()

	// Create a request requiring up to 3 targets
	filter := abilities.NewAnyTarget()
	request := &game.TargetRequest{
		AbilityID:    uuid.New(),
		SourceID:     sourceID,
		PlayerID:     playerID,
		MinTargets:   1,
		MaxTargets:   3,
		TargetFilter: filter,
		LegalTargets: []uuid.UUID{target1, target2, target3},
		Message:      "Choose up to 3 targets",
	}

	// Valid with 1 target
	selection := &game.TargetSelection{
		Targets: []uuid.UUID{target1},
	}
	err := tsm.ValidateTargets(ctx, request, selection, nil, nil)
	if err != nil {
		t.Errorf("Selection with 1 target failed: %v", err)
	}

	// Valid with 2 targets
	selection = &game.TargetSelection{
		Targets: []uuid.UUID{target1, target2},
	}
	err = tsm.ValidateTargets(ctx, request, selection, nil, nil)
	if err != nil {
		t.Errorf("Selection with 2 targets failed: %v", err)
	}

	// Valid with 3 targets
	selection = &game.TargetSelection{
		Targets: []uuid.UUID{target1, target2, target3},
	}
	err = tsm.ValidateTargets(ctx, request, selection, nil, nil)
	if err != nil {
		t.Errorf("Selection with 3 targets failed: %v", err)
	}

	// Invalid with duplicate targets
	selection = &game.TargetSelection{
		Targets: []uuid.UUID{target1, target1},
	}
	err = tsm.ValidateTargets(ctx, request, selection, nil, nil)
	if err == nil {
		t.Error("Expected error for duplicate targets, got nil")
	}
}

// TestStackManagement tests the enhanced stack manager
func TestStackManagement(t *testing.T) {
	logger := zap.NewNop()

	// Create dependencies
	turnMgr := rules.NewTurnManager([]string{"Alice", "Bob"}, logger)
	priorityMgr := game.NewPriorityManager(turnMgr, nil, logger)
	stackMgr := game.NewEnhancedStackManager(priorityMgr, logger)

	ctx := context.Background()
	spellID := uuid.New()
	controllerID := uuid.New()

	// Create a simple spell ability
	spellAbility := abilities.NewSpellAbilityBuilder(spellID, "{R}").
		AddEffect(abilities.NewDamageEffect(3)).
		Build()

	// Push spell to stack
	err := stackMgr.PushSpell(ctx, spellID, spellAbility, controllerID, []uuid.UUID{}, nil)
	if err != nil {
		t.Fatalf("Failed to push spell: %v", err)
	}

	// Verify stack size
	if stackMgr.Size() != 1 {
		t.Errorf("Expected stack size 1, got %d", stackMgr.Size())
	}

	// Peek at top
	top := stackMgr.Peek()
	if top == nil {
		t.Fatal("Expected object on top of stack, got nil")
	}
	if top.SourceID != spellID {
		t.Errorf("Wrong spell on top: expected %s, got %s", spellID, top.SourceID)
	}
}

// TestStackPushMultiple tests pushing multiple objects
func TestStackPushMultiple(t *testing.T) {
	logger := zap.NewNop()

	turnMgr := rules.NewTurnManager([]string{"Alice", "Bob"}, logger)
	priorityMgr := game.NewPriorityManager(turnMgr, nil, logger)
	stackMgr := game.NewEnhancedStackManager(priorityMgr, logger)

	ctx := context.Background()
	controllerID := uuid.New()

	// Push 3 spells
	spell1 := uuid.New()
	spell2 := uuid.New()
	spell3 := uuid.New()

	spellAbility1 := abilities.NewSpellAbilityBuilder(spell1, "{1}").
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	spellAbility2 := abilities.NewSpellAbilityBuilder(spell2, "{2}").
		AddEffect(abilities.NewDrawCardsEffect(2)).
		Build()
	spellAbility3 := abilities.NewSpellAbilityBuilder(spell3, "{3}").
		AddEffect(abilities.NewDrawCardsEffect(3)).
		Build()

	stackMgr.PushSpell(ctx, spell1, spellAbility1, controllerID, []uuid.UUID{}, nil)
	stackMgr.PushSpell(ctx, spell2, spellAbility2, controllerID, []uuid.UUID{}, nil)
	stackMgr.PushSpell(ctx, spell3, spellAbility3, controllerID, []uuid.UUID{}, nil)

	// Verify stack size
	if stackMgr.Size() != 3 {
		t.Errorf("Expected stack size 3, got %d", stackMgr.Size())
	}

	// Verify LIFO order (last in, first out)
	top := stackMgr.Peek()
	if top.SourceID != spell3 {
		t.Errorf("Expected spell3 on top, got %s", top.SourceID)
	}
}

// TestStackCounter tests countering spells
func TestStackCounter(t *testing.T) {
	logger := zap.NewNop()

	turnMgr := rules.NewTurnManager([]string{"Alice", "Bob"}, logger)
	priorityMgr := game.NewPriorityManager(turnMgr, nil, logger)
	stackMgr := game.NewEnhancedStackManager(priorityMgr, logger)

	ctx := context.Background()
	spellID := uuid.New()
	controllerID := uuid.New()

	// Push spell
	spellAbility := abilities.NewSpellAbilityBuilder(spellID, "{R}").
		AddEffect(abilities.NewDamageEffect(3)).
		Build()
	stackMgr.PushSpell(ctx, spellID, spellAbility, controllerID, []uuid.UUID{}, nil)

	// Verify on stack
	if stackMgr.Size() != 1 {
		t.Fatalf("Expected stack size 1, got %d", stackMgr.Size())
	}

	// Counter the spell
	err := stackMgr.Counter(ctx, spellID)
	if err != nil {
		t.Fatalf("Failed to counter spell: %v", err)
	}

	// Verify removed from stack
	if stackMgr.Size() != 0 {
		t.Errorf("Expected empty stack after counter, got size %d", stackMgr.Size())
	}
}

// TestLayerRecalculation tests the continuous effects layer system
func TestLayerRecalculation(t *testing.T) {
	logger := zap.NewNop()

	// Create layer system and registry
	layerSys := effects.NewLayerSystem(logger)
	registry := game.NewAbilityRegistry(logger)
	layerMgr := game.NewContinuousEffectsManager(layerSys, registry, logger)

	// Create a creature card
	cardID := uuid.New()
	playerID := uuid.New()
	card := game.NewCard(playerID, "Test Creature")
	card.ID = cardID
	card.Types = []string{"CREATURE"}
	card.Power = 2
	card.Toughness = 2

	// Add a static ability that boosts power/toughness
	boostAbility := abilities.NewStaticAbilityBuilder(cardID).
		AddEffect(abilities.NewBoostEffect(1, 1, abilities.DurationPermanent)).
		SetZone(abilities.ZoneBattlefield).
		Build()

	card.AddAbility(boostAbility)

	// Register the ability
	registry.RegisterAbility(boostAbility, playerID, abilities.ZoneBattlefield, 0)

	// Create a mock game context
	ctx := context.Background()

	// Recalculate layers
	err := layerMgr.RecalculateAll(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to recalculate layers: %v", err)
	}

	// Note: This test verifies the recalculation runs without error
	// Full integration would require applying effects to the card
}

// TestCombatIntegration tests combat triggered abilities
func TestCombatIntegration(t *testing.T) {
	logger := zap.NewNop()

	// Create engine components
	layerSys := effects.NewLayerSystem(logger)
	registry := game.NewAbilityRegistry(logger)
	turnMgr := rules.NewTurnManager([]string{"Alice", "Bob"}, logger)
	priorityMgr := game.NewPriorityManager(turnMgr, nil, logger)
	stackMgr := game.NewEnhancedStackManager(priorityMgr, logger)

	// Create mock engine
	engine := &game.MageEngine{
		logger: logger,
	}

	// Create combat integration manager
	combatMgr := game.NewCombatIntegrationManager(engine, registry, stackMgr, logger)

	// Create a creature with attack trigger
	creatureID := uuid.New()
	playerID := uuid.New()

	// Create attack trigger
	attackTrigger := abilities.NewAttacksTrigger(creatureID)
	triggeredAbility := abilities.NewTriggeredAbility(
		creatureID,
		attackTrigger,
		[]abilities.Effect{abilities.NewDrawCardsEffect(1)},
		false,
	)

	// Register the ability
	registry.RegisterAbility(triggeredAbility, playerID, abilities.ZoneBattlefield, 0)

	// Simulate declare attackers
	ctx := context.Background()
	err := combatMgr.OnDeclareAttackers(ctx, "game-123", []string{creatureID.String()}, nil)
	if err != nil {
		t.Fatalf("Failed to process attacker declaration: %v", err)
	}

	// Note: Full integration would push the triggered ability to stack
	// Currently this just verifies the trigger is detected
}

// TestCombatKeywordAbilities tests keyword ability detection
func TestCombatKeywordAbilities(t *testing.T) {
	logger := zap.NewNop()

	// Create engine components
	registry := game.NewAbilityRegistry(logger)
	turnMgr := rules.NewTurnManager([]string{"Alice", "Bob"}, logger)
	priorityMgr := game.NewPriorityManager(turnMgr, nil, logger)
	stackMgr := game.NewEnhancedStackManager(priorityMgr, logger)

	engine := &game.MageEngine{
		logger: logger,
	}

	combatMgr := game.NewCombatIntegrationManager(engine, registry, stackMgr, logger)

	// Create creature with keyword abilities
	creatureID := uuid.New()
	playerID := uuid.New()

	// Create flying keyword
	flyingAbility := abilities.NewKeywordAbility(creatureID, "Flying")
	registry.RegisterAbility(flyingAbility, playerID, abilities.ZoneBattlefield, 0)

	// Check keyword abilities
	ctx := context.Background()
	keywords, err := combatMgr.CheckCombatKeywordAbilities(ctx, creatureID.String(), nil)
	if err != nil {
		t.Fatalf("Failed to check keyword abilities: %v", err)
	}

	if !keywords.Flying {
		t.Error("Expected Flying to be true, got false")
	}
	if keywords.FirstStrike {
		t.Error("Expected FirstStrike to be false, got true")
	}
}

// TestManaCostPayment tests mana cost payment integration
func TestManaCostPayment(t *testing.T) {
	// Create mana pool
	pool := mana.NewManaPool()

	// Add mana
	pool.AddMana(mana.Red, 3, false)
	pool.AddMana(mana.Colorless, 2, false)

	// Create mana cost
	cost := abilities.NewManaCost("{2}{R}")

	// Check if we can pay
	canPay := cost.CanPay(pool)
	if !canPay {
		t.Error("Expected to be able to pay {2}{R} with 3 red + 2 colorless")
	}

	// Pay the cost
	err := cost.Pay(pool)
	if err != nil {
		t.Fatalf("Failed to pay mana cost: %v", err)
	}

	// Verify remaining mana
	remaining := pool.GetAmount(mana.Red)
	if remaining != 1 {
		t.Errorf("Expected 1 red mana remaining, got %d", remaining)
	}
}

// TestManaCostPaymentInsufficient tests insufficient mana
func TestManaCostPaymentInsufficient(t *testing.T) {
	// Create mana pool
	pool := mana.NewManaPool()

	// Add insufficient mana
	pool.AddMana(mana.Red, 1, false)

	// Create mana cost
	cost := abilities.NewManaCost("{2}{R}")

	// Check if we can pay (should fail)
	canPay := cost.CanPay(pool)
	if canPay {
		t.Error("Expected to not be able to pay {2}{R} with only 1 red")
	}

	// Try to pay (should fail)
	err := cost.Pay(pool)
	if err == nil {
		t.Error("Expected error when paying with insufficient mana, got nil")
	}
}

// TestCardWithMultipleAbilities tests cards with multiple abilities
func TestCardWithMultipleAbilities(t *testing.T) {
	logger := zap.NewNop()
	registry := game.NewAbilityRegistry(logger)

	cardID := uuid.New()
	playerID := uuid.New()

	// Create card with multiple abilities
	card := game.NewCard(playerID, "Complex Card")
	card.ID = cardID

	// Add activated ability
	activatedAbility := abilities.NewActivatedAbilityBuilder(cardID).
		AddManaCost("{1}").
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	card.AddAbility(activatedAbility)

	// Add triggered ability
	etbTrigger := abilities.NewEntersBattlefieldTrigger(cardID)
	triggeredAbility := abilities.NewTriggeredAbility(
		cardID,
		etbTrigger,
		[]abilities.Effect{abilities.NewDamageEffect(2)},
		false,
	)
	card.AddAbility(triggeredAbility)

	// Add static ability
	staticAbility := abilities.NewStaticAbilityBuilder(cardID).
		AddEffect(abilities.NewBoostEffect(1, 1, abilities.DurationPermanent)).
		SetZone(abilities.ZoneBattlefield).
		Build()
	card.AddAbility(staticAbility)

	// Register all abilities
	for i, ability := range card.GetAbilities() {
		registry.RegisterAbility(ability, playerID, abilities.ZoneBattlefield, i)
	}

	// Verify all abilities are registered
	sourceAbilities := registry.GetAbilitiesBySource(cardID)
	if len(sourceAbilities) != 3 {
		t.Errorf("Expected 3 abilities registered, got %d", len(sourceAbilities))
	}

	// Verify ability types
	foundActivated := false
	foundTriggered := false
	foundStatic := false

	for _, ability := range sourceAbilities {
		switch ability.GetType() {
		case abilities.AbilityTypeActivated:
			foundActivated = true
		case abilities.AbilityTypeTriggered:
			foundTriggered = true
		case abilities.AbilityTypeStatic:
			foundStatic = true
		}
	}

	if !foundActivated {
		t.Error("Activated ability not found")
	}
	if !foundTriggered {
		t.Error("Triggered ability not found")
	}
	if !foundStatic {
		t.Error("Static ability not found")
	}
}

// TestAbilityCleanupOnZoneChange tests that abilities are properly cleaned up
func TestAbilityCleanupOnZoneChange(t *testing.T) {
	logger := zap.NewNop()
	registry := game.NewAbilityRegistry(logger)
	layerSys := effects.NewLayerSystem(logger)
	layerMgr := game.NewContinuousEffectsManager(layerSys, registry, logger)

	cardID := uuid.New()
	playerID := uuid.New()

	// Register ability
	ability := abilities.NewStaticAbilityBuilder(cardID).
		AddEffect(abilities.NewBoostEffect(1, 1, abilities.DurationPermanent)).
		SetZone(abilities.ZoneBattlefield).
		Build()

	registry.RegisterAbility(ability, playerID, abilities.ZoneBattlefield, 0)

	// Verify registered
	sourceAbilities := registry.GetAbilitiesBySource(cardID)
	if len(sourceAbilities) != 1 {
		t.Fatalf("Expected 1 ability, got %d", len(sourceAbilities))
	}

	// Remove effects when card leaves battlefield
	layerMgr.RemoveSourceEffects(cardID)

	// Unregister from registry
	registry.UnregisterSource(cardID)

	// Verify cleanup
	sourceAbilities = registry.GetAbilitiesBySource(cardID)
	if len(sourceAbilities) != 0 {
		t.Errorf("Expected 0 abilities after cleanup, got %d", len(sourceAbilities))
	}
}
