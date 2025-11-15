package rules

import (
	"testing"
	"time"
)

func TestTriggerManagerHandle(t *testing.T) {
	manager := NewTriggerManager()

	callCount := 0
	manager.Register(AbilityTrigger{
		EventType: EventSpellCast,
		Condition: func(e Event) bool {
			return e.Metadata["card_name"] == "Lightning Bolt"
		},
		Build: func(e Event) StackItem {
			callCount++
			return StackItem{
				Controller:  e.Controller,
				Description: "Deal 3 damage",
				Kind:        StackItemKindTriggered,
			}
		},
	})

	items := manager.Handle(Event{
		Type:       EventSpellCast,
		Controller: "Alice",
		Timestamp:  time.Now(),
		Metadata: map[string]string{
			"card_name": "Lightning Bolt",
		},
	})

	if len(items) != 1 {
		t.Fatalf("expected 1 stack item, got %d", len(items))
	}
	if items[0].Controller != "Alice" {
		t.Fatalf("expected controller Alice, got %s", items[0].Controller)
	}
	if callCount != 1 {
		t.Fatalf("expected build to be called once, got %d", callCount)
	}
}
