package rules

import (
	"errors"
	"sync"
)

// StackItemKind describes the type of object on the stack.
type StackItemKind string

const (
	// StackItemKindSpell represents a spell cast by a player.
	StackItemKindSpell StackItemKind = "SPELL"
	// StackItemKindActivated represents an activated ability.
	StackItemKindActivated StackItemKind = "ACTIVATED"
	// StackItemKindTriggered represents a triggered ability.
	StackItemKindTriggered StackItemKind = "TRIGGERED"
)

// StackItem represents a single object on the stack.
type StackItem struct {
	ID          string
	Controller  string
	Description string
	Kind        StackItemKind
	SourceID    string
	Metadata    map[string]string
	Resolve     func() error
	onRemove    func()
}

// StackManager manages the game stack.
type StackManager struct {
	mu    sync.Mutex
	items []StackItem
}

// NewStackManager creates a new stack manager.
func NewStackManager() *StackManager {
	return &StackManager{
		items: make([]StackItem, 0, 16),
	}
}

// Push adds an item to the top of the stack.
func (sm *StackManager) Push(item StackItem) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.items = append(sm.items, item)
}

// Pop removes the top item from the stack.
func (sm *StackManager) Pop() (StackItem, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if len(sm.items) == 0 {
		return StackItem{}, errors.New("stack empty")
	}

	idx := len(sm.items) - 1
	item := sm.items[idx]
	sm.items = sm.items[:idx]
	return item, nil
}

// Remove deletes an item from anywhere in the stack by ID.
func (sm *StackManager) Remove(id string) (StackItem, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	for idx := len(sm.items) - 1; idx >= 0; idx-- {
		if sm.items[idx].ID == id {
			item := sm.items[idx]
			sm.items = append(sm.items[:idx], sm.items[idx+1:]...)
			return item, true
		}
	}
	return StackItem{}, false
}

// Peek returns the top item without removing it.
func (sm *StackManager) Peek() (StackItem, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if len(sm.items) == 0 {
		return StackItem{}, false
	}
	return sm.items[len(sm.items)-1], true
}

// List returns a copy of all stack items (topmost last).
func (sm *StackManager) List() []StackItem {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	cpy := make([]StackItem, len(sm.items))
	copy(cpy, sm.items)
	return cpy
}

// IsEmpty returns whether the stack is empty.
func (sm *StackManager) IsEmpty() bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return len(sm.items) == 0
}

// RemoveIllegalItems removes all illegal items from the stack using the provided legality checker.
// Returns the IDs of removed items.
func (sm *StackManager) RemoveIllegalItems(checker *LegalityChecker) []string {
	if checker == nil {
		return nil
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var removedIDs []string
	validItems := make([]StackItem, 0, len(sm.items))

	for _, item := range sm.items {
		result := checker.CheckStackItemLegality(item)
		if !result.Legal {
			removedIDs = append(removedIDs, item.ID)
			if item.onRemove != nil {
				item.onRemove()
			}
		} else {
			validItems = append(validItems, item)
		}
	}

	sm.items = validItems
	return removedIDs
}
