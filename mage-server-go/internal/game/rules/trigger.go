package rules

import (
	"sync"

	"github.com/google/uuid"
)

// AbilityTrigger encapsulates the logic for reacting to a specific event and
// producing stack items when the conditions are satisfied.
type AbilityTrigger struct {
	ID         string
	SourceID   string
	Controller string
	EventType  EventType
	Condition  func(Event) bool
	Build      func(Event) StackItem
	Once       bool
}

// TriggerManager stores and evaluates ability triggers against events.
type TriggerManager struct {
	mu        sync.Mutex
	triggers  map[string]AbilityTrigger
	listeners map[int]struct{}
}

// NewTriggerManager creates an empty trigger manager.
func NewTriggerManager() *TriggerManager {
	return &TriggerManager{
		triggers:  make(map[string]AbilityTrigger),
		listeners: make(map[int]struct{}),
	}
}

// Register adds a new trigger to the manager.
func (tm *TriggerManager) Register(trigger AbilityTrigger) string {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if trigger.ID == "" {
		trigger.ID = uuid.NewString()
	}
	tm.triggers[trigger.ID] = trigger
	return trigger.ID
}

// Unregister removes a trigger by ID.
func (tm *TriggerManager) Unregister(id string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	delete(tm.triggers, id)
}

// Handle evaluates the provided event against all registered triggers and
// returns the stack items they produce.
func (tm *TriggerManager) Handle(event Event) []StackItem {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if len(tm.triggers) == 0 {
		return nil
	}

	var (
		stackItems []StackItem
		toRemove   []string
	)

	for id, trigger := range tm.triggers {
		if trigger.EventType != event.Type {
			continue
		}
		if trigger.Condition != nil && !trigger.Condition(event) {
			continue
		}
		if trigger.Build == nil {
			continue
		}

		item := trigger.Build(event)
		if item.ID == "" {
			item.ID = uuid.NewString()
		}
		stackItems = append(stackItems, item)

		if trigger.Once {
			toRemove = append(toRemove, id)
		}
	}

	for _, id := range toRemove {
		delete(tm.triggers, id)
	}

	return stackItems
}
