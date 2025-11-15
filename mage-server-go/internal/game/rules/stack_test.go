package rules

import "testing"

func TestStackManagerPushPop(t *testing.T) {
	sm := NewStackManager()

	firstResolved := false
	secondResolved := false

	sm.Push(StackItem{
		ID:          "first",
		Controller:  "Alice",
		Description: "First Spell",
		Kind:        StackItemKindSpell,
		Metadata:    map[string]string{"card_name": "First"},
		Resolve: func() error {
			firstResolved = true
			return nil
		},
	})

	sm.Push(StackItem{
		ID:          "second",
		Controller:  "Bob",
		Description: "Second Spell",
		Kind:        StackItemKindTriggered,
		Resolve: func() error {
			secondResolved = true
			return nil
		},
	})

	item, err := sm.Pop()
	if err != nil {
		t.Fatalf("unexpected error popping top: %v", err)
	}
	if item.ID != "second" {
		t.Fatalf("expected LIFO order (second), got %s", item.ID)
	}
	if err := item.Resolve(); err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if !secondResolved {
		t.Fatalf("expected second resolve to run")
	}

	item, err = sm.Pop()
	if err != nil {
		t.Fatalf("unexpected error popping second item: %v", err)
	}
	if item.ID != "first" {
		t.Fatalf("expected remaining item to be first, got %s", item.ID)
	}
	if err := item.Resolve(); err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if !firstResolved {
		t.Fatalf("expected first resolve to run")
	}

	if !sm.IsEmpty() {
		t.Fatalf("expected stack to be empty")
	}
}

func TestStackManagerRemove(t *testing.T) {
	sm := NewStackManager()

	sm.Push(StackItem{ID: "first"})
	sm.Push(StackItem{ID: "second"})
	sm.Push(StackItem{ID: "third"})

	item, ok := sm.Remove("second")
	if !ok {
		t.Fatalf("expected to remove existing item")
	}
	if item.ID != "second" {
		t.Fatalf("expected removed ID second, got %s", item.ID)
	}

	top, _ := sm.Pop()
	if top.ID != "third" {
		t.Fatalf("expected third to remain on top, got %s", top.ID)
	}
}
