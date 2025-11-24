package game

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/rules"
	"go.uber.org/zap"
)

// StackObject represents an object on the stack (spell or ability)
type StackObject struct {
	ID         uuid.UUID
	Kind       StackObjectKind
	Controller uuid.UUID
	SourceID   uuid.UUID
	Ability    abilities.Ability
	Targets    []uuid.UUID
	Choices    map[string]interface{} // Modal choices, X values, etc.
	CardID     uuid.UUID              // For spells, the card being cast
}

// StackObjectKind indicates what type of object this is
type StackObjectKind int

const (
	StackObjectSpell StackObjectKind = iota
	StackObjectActivatedAbility
	StackObjectTriggeredAbility
)

func (k StackObjectKind) String() string {
	switch k {
	case StackObjectSpell:
		return "Spell"
	case StackObjectActivatedAbility:
		return "Activated Ability"
	case StackObjectTriggeredAbility:
		return "Triggered Ability"
	default:
		return "Unknown"
	}
}

// EnhancedStackManager wraps the rules.StackManager and adds abilities integration
type EnhancedStackManager struct {
	stack       *rules.StackManager
	objects     map[string]*StackObject // Map stack item ID to stack object
	priorityMgr *PriorityManager
	logger      *zap.Logger
}

// NewEnhancedStackManager creates a new enhanced stack manager
func NewEnhancedStackManager(priorityMgr *PriorityManager, logger *zap.Logger) *EnhancedStackManager {
	return &EnhancedStackManager{
		stack:       rules.NewStackManager(),
		objects:     make(map[string]*StackObject),
		priorityMgr: priorityMgr,
		logger:      logger,
	}
}

// PushSpell adds a spell to the stack
func (esm *EnhancedStackManager) PushSpell(
	ctx context.Context,
	cardID uuid.UUID,
	controller uuid.UUID,
	ability abilities.Ability,
	targets []uuid.UUID,
	choices map[string]interface{},
) (uuid.UUID, error) {
	stackObjID := uuid.New()

	stackObject := &StackObject{
		ID:         stackObjID,
		Kind:       StackObjectSpell,
		Controller: controller,
		SourceID:   ability.GetSourceID(),
		Ability:    ability,
		Targets:    targets,
		Choices:    choices,
		CardID:     cardID,
	}

	// Store the stack object
	esm.objects[stackObjID.String()] = stackObject

	// Create the stack item for the rules engine
	stackItem := rules.StackItem{
		ID:          stackObjID.String(),
		Controller:  controller.String(),
		Description: ability.String(),
		Kind:        rules.StackItemKindSpell,
		SourceID:    ability.GetSourceID().String(),
		Metadata:    make(map[string]string),
		Resolve:     nil, // Will be set below
		// onRemove will be set to clean up our objects map
	}

	// Set the resolve function
	stackItem.Resolve = func() error {
		return esm.resolveStackObject(ctx, stackObjID)
	}

	// Push to the stack
	esm.stack.Push(stackItem)

	esm.logger.Info("spell added to stack",
		zap.String("id", stackObjID.String()),
		zap.String("card", cardID.String()),
		zap.String("controller", controller.String()))

	return stackObjID, nil
}

// PushActivatedAbility adds an activated ability to the stack
func (esm *EnhancedStackManager) PushActivatedAbility(
	ctx context.Context,
	controller uuid.UUID,
	ability abilities.Ability,
	targets []uuid.UUID,
	choices map[string]interface{},
) (uuid.UUID, error) {
	stackObjID := uuid.New()

	stackObject := &StackObject{
		ID:         stackObjID,
		Kind:       StackObjectActivatedAbility,
		Controller: controller,
		SourceID:   ability.GetSourceID(),
		Ability:    ability,
		Targets:    targets,
		Choices:    choices,
	}

	esm.objects[stackObjID.String()] = stackObject

	stackItem := rules.StackItem{
		ID:          stackObjID.String(),
		Controller:  controller.String(),
		Description: ability.String(),
		Kind:        rules.StackItemKindActivated,
		SourceID:    ability.GetSourceID().String(),
		Metadata:    make(map[string]string),
		Resolve: func() error {
			return esm.resolveStackObject(ctx, stackObjID)
		},
	}

	esm.stack.Push(stackItem)

	esm.logger.Info("activated ability added to stack",
		zap.String("id", stackObjID.String()),
		zap.String("source", ability.GetSourceID().String()),
		zap.String("controller", controller.String()))

	return stackObjID, nil
}

// PushTriggeredAbility adds a triggered ability to the stack
func (esm *EnhancedStackManager) PushTriggeredAbility(
	ctx context.Context,
	controller uuid.UUID,
	ability abilities.Ability,
	targets []uuid.UUID,
) (uuid.UUID, error) {
	stackObjID := uuid.New()

	stackObject := &StackObject{
		ID:         stackObjID,
		Kind:       StackObjectTriggeredAbility,
		Controller: controller,
		SourceID:   ability.GetSourceID(),
		Ability:    ability,
		Targets:    targets,
		Choices:    make(map[string]interface{}),
	}

	esm.objects[stackObjID.String()] = stackObject

	stackItem := rules.StackItem{
		ID:          stackObjID.String(),
		Controller:  controller.String(),
		Description: ability.String(),
		Kind:        rules.StackItemKindTriggered,
		SourceID:    ability.GetSourceID().String(),
		Metadata:    make(map[string]string),
		Resolve: func() error {
			return esm.resolveStackObject(ctx, stackObjID)
		},
	}

	esm.stack.Push(stackItem)

	esm.logger.Info("triggered ability added to stack",
		zap.String("id", stackObjID.String()),
		zap.String("source", ability.GetSourceID().String()),
		zap.String("controller", controller.String()))

	return stackObjID, nil
}

// ResolveTop resolves the top item on the stack
func (esm *EnhancedStackManager) ResolveTop(ctx context.Context, gameCtx abilities.GameContext, state rules.GameStateReader) error {
	// Get the top item
	stackItem, ok := esm.stack.Peek()
	if !ok {
		return fmt.Errorf("stack is empty")
	}

	esm.logger.Info("resolving stack object",
		zap.String("id", stackItem.ID),
		zap.String("kind", string(stackItem.Kind)),
		zap.String("description", stackItem.Description))

	// Resolve the item
	if stackItem.Resolve != nil {
		if err := stackItem.Resolve(); err != nil {
			esm.logger.Error("failed to resolve stack object",
				zap.String("id", stackItem.ID),
				zap.Error(err))
			return fmt.Errorf("failed to resolve: %w", err)
		}
	}

	// Pop the item from the stack
	if _, err := esm.stack.Pop(); err != nil {
		return fmt.Errorf("failed to pop stack: %w", err)
	}

	// Clean up the stack object
	delete(esm.objects, stackItem.ID)

	// After resolving, check state-based actions (Rule 608.2k)
	if err := esm.priorityMgr.AfterSpellResolves(ctx, state); err != nil {
		return fmt.Errorf("post-resolution SBA check failed: %w", err)
	}

	esm.logger.Info("stack object resolved successfully",
		zap.String("id", stackItem.ID))

	return nil
}

// resolveStackObject resolves a stack object by its ID
func (esm *EnhancedStackManager) resolveStackObject(ctx context.Context, id uuid.UUID) error {
	stackObject, ok := esm.objects[id.String()]
	if !ok {
		return fmt.Errorf("stack object not found: %s", id.String())
	}

	esm.logger.Debug("resolving ability",
		zap.String("id", id.String()),
		zap.String("kind", stackObject.Kind.String()))

	// For spells, we need to handle the card moving from stack to appropriate zone
	// For now, we just resolve the ability
	// TODO: Implement full spell resolution (move card to graveyard, battlefield, etc.)

	// Create a game context that can access the stack object's targets and choices
	// For now, we use the passed gameCtx
	// TODO: Create a specialized context that includes target and choice information

	// Resolve the ability
	// Note: The ability's Resolve method will apply all effects
	// We would pass targets here if Resolve took them
	// For now, abilities need to be refactored to accept targets in Resolve
	if err := stackObject.Ability.Resolve(ctx, nil); err != nil {
		return fmt.Errorf("failed to resolve ability: %w", err)
	}

	return nil
}

// Counter removes a spell or ability from the stack (counterspell effect)
func (esm *EnhancedStackManager) Counter(stackObjectID uuid.UUID) error {
	stackItem, ok := esm.stack.Remove(stackObjectID.String())
	if !ok {
		return fmt.Errorf("stack object not found: %s", stackObjectID.String())
	}

	// Clean up our tracking
	delete(esm.objects, stackObjectID.String())

	esm.logger.Info("stack object countered",
		zap.String("id", stackObjectID.String()),
		zap.String("kind", string(stackItem.Kind)))

	// For spells, the card goes to graveyard (Rule 701.5a)
	// TODO: Move card to graveyard

	return nil
}

// IsEmpty returns true if the stack is empty
func (esm *EnhancedStackManager) IsEmpty() bool {
	return esm.stack.IsEmpty()
}

// GetTop returns the top stack object without removing it
func (esm *EnhancedStackManager) GetTop() (*StackObject, bool) {
	stackItem, ok := esm.stack.Peek()
	if !ok {
		return nil, false
	}

	stackObject, ok := esm.objects[stackItem.ID]
	return stackObject, ok
}

// GetAll returns all stack objects (bottom to top)
func (esm *EnhancedStackManager) GetAll() []*StackObject {
	items := esm.stack.List()
	objects := make([]*StackObject, 0, len(items))

	for _, item := range items {
		if obj, ok := esm.objects[item.ID]; ok {
			objects = append(objects, obj)
		}
	}

	return objects
}

// RemoveIllegalTargets removes stack objects that have illegal targets
// This is called during state-based actions check (Rule 608.2b)
func (esm *EnhancedStackManager) RemoveIllegalTargets(ctx context.Context, state rules.GameStateReader) []uuid.UUID {
	removed := []uuid.UUID{}

	// Check each stack object for illegal targets
	for idStr, obj := range esm.objects {
		if esm.hasIllegalTargets(obj, state) {
			// Remove from stack
			if _, ok := esm.stack.Remove(idStr); ok {
				removed = append(removed, obj.ID)
				delete(esm.objects, idStr)

				esm.logger.Info("removed stack object with illegal targets",
					zap.String("id", idStr))
			}
		}
	}

	return removed
}

// hasIllegalTargets checks if a stack object has illegal targets
func (esm *EnhancedStackManager) hasIllegalTargets(obj *StackObject, state rules.GameStateReader) bool {
	// A target is illegal if:
	// 1. The target no longer exists
	// 2. The target is no longer a legal target (e.g., gained shroud)
	// 3. For Aura spells, the target is not a legal permanent to enchant

	// TODO: Implement full target legality checking
	// For now, just check if targets exist

	for _, targetID := range obj.Targets {
		// Check if target still exists
		if _, ok := state.GetPermanent(targetID); !ok {
			// Check if it's a player
			if _, ok := state.GetPlayer(targetID); !ok {
				// Target doesn't exist
				return true
			}
		}
	}

	return false
}

// Clear removes all items from the stack
// Used for game cleanup or special effects
func (esm *EnhancedStackManager) Clear() {
	items := esm.stack.List()
	for _, item := range items {
		delete(esm.objects, item.ID)
	}

	// Clear the underlying stack
	for !esm.stack.IsEmpty() {
		_, _ = esm.stack.Pop()
	}

	esm.logger.Info("stack cleared")
}
