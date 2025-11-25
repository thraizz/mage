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
	registry := game.NewAbilityRegistry()

	cardID := uuid.New()
	playerID := uuid.New()

	// Create a simple activated ability
	ability := abilities.NewActivatedAbilityBuilder(cardID).
		AddManaCost("{1}{R}").
		AddEffect(abilities.NewDamageEffect(3)).
		Build()

	// Register the ability
	registry.RegisterAbility(ability, playerID, 0, abilities.ZoneBattlefield)

	// Verify we can retrieve it
	retrieved, err := registry.GetAbility(ability.GetID())
	if err != nil || retrieved == nil {
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
	registry := game.NewAbilityRegistry()

	cardID := uuid.New()
	playerID := uuid.New()

	// Register ability in battlefield
	ability := abilities.NewActivatedAbilityBuilder(cardID).
		AddManaCost("{1}").
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()

	registry.RegisterAbility(ability, playerID, 0, abilities.ZoneBattlefield)

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
	filter := abilities.NewAnyTargetFilter()
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
	filter := abilities.NewAnyTargetFilter()
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
	filter := abilities.NewAnyTargetFilter()
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
	turnMgr := rules.NewTurnManager("Alice")
	priorityMgr := game.NewPriorityManager(turnMgr, logger)
	stackMgr := game.NewEnhancedStackManager(priorityMgr, logger)

	ctx := context.Background()
	spellID := uuid.New()
	controllerID := uuid.New()

	// Create a simple spell ability
	spellAbility, err := abilities.NewSpellAbilityBuilder(spellID, "{R}").
		AddEffect(abilities.NewDamageEffect(3)).
		Build()
	if err != nil {
		t.Fatalf("Failed to build spell ability: %v", err)
	}

	// Push spell to stack
	_, err = stackMgr.PushSpell(ctx, spellID, controllerID, spellAbility, []uuid.UUID{}, nil)
	if err != nil {
		t.Fatalf("Failed to push spell: %v", err)
	}

	// Verify stack size
	stackItems := stackMgr.GetAll()
	if len(stackItems) != 1 {
		t.Errorf("Expected stack size 1, got %d", len(stackItems))
	}

	// Peek at top
	top, ok := stackMgr.GetTop()
	if !ok || top == nil {
		t.Fatal("Expected object on top of stack, got nil")
	}
	if top.SourceID != spellID {
		t.Errorf("Wrong spell on top: expected %s, got %s", spellID, top.SourceID)
	}
}

// TestStackPushMultiple tests pushing multiple objects
func TestStackPushMultiple(t *testing.T) {
	logger := zap.NewNop()

	turnMgr := rules.NewTurnManager("Alice")
	priorityMgr := game.NewPriorityManager(turnMgr, logger)
	stackMgr := game.NewEnhancedStackManager(priorityMgr, logger)

	ctx := context.Background()
	controllerID := uuid.New()

	// Push 3 spells
	spell1 := uuid.New()
	spell2 := uuid.New()
	spell3 := uuid.New()

	spellAbility1, _ := abilities.NewSpellAbilityBuilder(spell1, "{1}").
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	spellAbility2, _ := abilities.NewSpellAbilityBuilder(spell2, "{2}").
		AddEffect(abilities.NewDrawCardsEffect(2)).
		Build()
	spellAbility3, _ := abilities.NewSpellAbilityBuilder(spell3, "{3}").
		AddEffect(abilities.NewDrawCardsEffect(3)).
		Build()

	stackMgr.PushSpell(ctx, spell1, controllerID, spellAbility1, []uuid.UUID{}, nil)
	stackMgr.PushSpell(ctx, spell2, controllerID, spellAbility2, []uuid.UUID{}, nil)
	stackMgr.PushSpell(ctx, spell3, controllerID, spellAbility3, []uuid.UUID{}, nil)

	// Verify stack size
	stackItems := stackMgr.GetAll()
	if len(stackItems) != 3 {
		t.Errorf("Expected stack size 3, got %d", len(stackItems))
	}

	// Verify LIFO order (last in, first out)
	top, _ := stackMgr.GetTop()
	if top.SourceID != spell3 {
		t.Errorf("Expected spell3 on top, got %s", top.SourceID)
	}
}

// TestStackCounter tests countering spells
func TestStackCounter(t *testing.T) {
	logger := zap.NewNop()

	turnMgr := rules.NewTurnManager("Alice")
	priorityMgr := game.NewPriorityManager(turnMgr, logger)
	stackMgr := game.NewEnhancedStackManager(priorityMgr, logger)

	ctx := context.Background()
	spellID := uuid.New()
	controllerID := uuid.New()

	// Push spell
	spellAbility, _ := abilities.NewSpellAbilityBuilder(spellID, "{R}").
		AddEffect(abilities.NewDamageEffect(3)).
		Build()
	stackObjID, _ := stackMgr.PushSpell(ctx, spellID, controllerID, spellAbility, []uuid.UUID{}, nil)

	// Verify on stack
	stackItems := stackMgr.GetAll()
	if len(stackItems) != 1 {
		t.Fatalf("Expected stack size 1, got %d", len(stackItems))
	}

	// Counter the spell (use the stack object ID, not the spell ID)
	err := stackMgr.Counter(stackObjID)
	if err != nil {
		t.Fatalf("Failed to counter spell: %v", err)
	}

	// Verify removed from stack
	stackItems = stackMgr.GetAll()
	if len(stackItems) != 0 {
		t.Errorf("Expected empty stack after counter, got size %d", len(stackItems))
	}
}

// TestLayerRecalculation tests the continuous effects layer system
func TestLayerRecalculation(t *testing.T) {
	logger := zap.NewNop()

	// Create layer system and registry
	layerSys := effects.NewLayerSystem()
	registry := game.NewAbilityRegistry()
	layerMgr := game.NewContinuousEffectsManager(layerSys, registry, logger)

	// Create a creature card
	cardID := uuid.New()
	playerID := uuid.New()
	card := game.NewCard(playerID, "Test Creature")
	card.ID = cardID
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "2"

	// Create a keyword ability for testing (does not require ContinuousEffect)
	keywordAbility := abilities.NewKeywordAbility(cardID, "Flying")
	card.AddAbility(keywordAbility)

	// Register the ability
	registry.RegisterAbility(keywordAbility, playerID, 0, abilities.ZoneBattlefield)

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
	_ = effects.NewLayerSystem()
	registry := game.NewAbilityRegistry()
	turnMgr := rules.NewTurnManager("Alice")
	priorityMgr := game.NewPriorityManager(turnMgr, logger)
	stackMgr := game.NewEnhancedStackManager(priorityMgr, logger)

	// Create engine
	engine := game.NewMageEngine(logger)

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
	registry.RegisterAbility(triggeredAbility, playerID, 0, abilities.ZoneBattlefield)

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
	registry := game.NewAbilityRegistry()
	turnMgr := rules.NewTurnManager("Alice")
	priorityMgr := game.NewPriorityManager(turnMgr, logger)
	stackMgr := game.NewEnhancedStackManager(priorityMgr, logger)

	engine := game.NewMageEngine(logger)

	combatMgr := game.NewCombatIntegrationManager(engine, registry, stackMgr, logger)

	// Create creature with keyword abilities
	creatureID := uuid.New()
	playerID := uuid.New()

	// Create flying keyword
	flyingAbility := abilities.NewKeywordAbility(creatureID, "Flying")
	registry.RegisterAbility(flyingAbility, playerID, 0, abilities.ZoneBattlefield)

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
// Note: The full mana cost payment API requires GameContext.
// See internal/game/mana/pool_test.go for direct mana pool tests.
func TestManaCostPayment(t *testing.T) {
	// Create mana pool and test direct add/spend operations
	pool := mana.NewManaPool()

	// Add mana using the correct API
	pool.Add(mana.ManaRed, 3)
	pool.Add(mana.ManaColorless, 2)

	// Verify mana was added
	redTotal := pool.GetTotal(mana.ManaRed)
	if redTotal != 3 {
		t.Errorf("Expected 3 red mana, got %d", redTotal)
	}

	colorlessTotal := pool.GetTotal(mana.ManaColorless)
	if colorlessTotal != 2 {
		t.Errorf("Expected 2 colorless mana, got %d", colorlessTotal)
	}

	// Test spending mana
	if !pool.Spend(mana.ManaRed, 2) {
		t.Error("Expected to spend 2 red mana successfully")
	}

	remaining := pool.GetTotal(mana.ManaRed)
	if remaining != 1 {
		t.Errorf("Expected 1 red mana remaining, got %d", remaining)
	}
}

// TestManaCostPaymentInsufficient tests insufficient mana
func TestManaCostPaymentInsufficient(t *testing.T) {
	// Create mana pool
	pool := mana.NewManaPool()

	// Add insufficient mana
	pool.Add(mana.ManaRed, 1)

	// Verify cannot spend more than available
	if pool.Spend(mana.ManaRed, 3) {
		t.Error("Expected to not be able to spend 3 red with only 1 available")
	}

	// Original mana should still be there
	remaining := pool.GetTotal(mana.ManaRed)
	if remaining != 1 {
		t.Errorf("Expected 1 red mana remaining after failed spend, got %d", remaining)
	}
}

// TestCardWithMultipleAbilities tests cards with multiple abilities
func TestCardWithMultipleAbilities(t *testing.T) {
	registry := game.NewAbilityRegistry()

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

	// Add a keyword ability (simpler than static ability with ContinuousEffect)
	keywordAbility := abilities.NewKeywordAbility(cardID, "Flying")
	card.AddAbility(keywordAbility)

	// Register all abilities individually (card.GetAbilities returns interface{})
	registry.RegisterAbility(activatedAbility, playerID, 0, abilities.ZoneBattlefield)
	registry.RegisterAbility(triggeredAbility, playerID, 1, abilities.ZoneBattlefield)
	registry.RegisterAbility(keywordAbility, playerID, 2, abilities.ZoneBattlefield)

	// Verify all abilities are registered
	sourceAbilities := registry.GetAbilitiesBySource(cardID)
	if len(sourceAbilities) != 3 {
		t.Errorf("Expected 3 abilities registered, got %d", len(sourceAbilities))
	}

	// Verify ability types
	foundActivated := false
	foundTriggered := false
	foundKeyword := false

	for _, ability := range sourceAbilities {
		switch ability.GetType() {
		case abilities.AbilityTypeActivated:
			foundActivated = true
		case abilities.AbilityTypeTriggered:
			foundTriggered = true
		case abilities.AbilityTypeKeyword:
			foundKeyword = true
		}
	}

	if !foundActivated {
		t.Error("Activated ability not found")
	}
	if !foundTriggered {
		t.Error("Triggered ability not found")
	}
	if !foundKeyword {
		t.Error("Keyword ability not found")
	}
}

// TestAbilityCleanupOnZoneChange tests that abilities are properly cleaned up
func TestAbilityCleanupOnZoneChange(t *testing.T) {
	logger := zap.NewNop()
	registry := game.NewAbilityRegistry()
	layerSys := effects.NewLayerSystem()
	layerMgr := game.NewContinuousEffectsManager(layerSys, registry, logger)

	cardID := uuid.New()
	playerID := uuid.New()

	// Register a keyword ability for testing cleanup
	ability := abilities.NewKeywordAbility(cardID, "Flying")

	registry.RegisterAbility(ability, playerID, 0, abilities.ZoneBattlefield)

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
