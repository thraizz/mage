package rules

import (
	"testing"
	"time"
)

func TestEventBusSubscribeTyped(t *testing.T) {
	bus := NewEventBus()

	spellCastCount := 0
	lifeGainCount := 0

	handle1 := bus.SubscribeTyped(EventSpellCast, func(e Event) {
		spellCastCount++
	})

	handle2 := bus.SubscribeTyped(EventGainedLife, func(e Event) {
		lifeGainCount++
	})

	// Publish spell cast event
	bus.Publish(NewEvent(EventSpellCast, "card1", "card1", "player1"))
	if spellCastCount != 1 {
		t.Fatalf("expected spell cast count 1, got %d", spellCastCount)
	}
	if lifeGainCount != 0 {
		t.Fatalf("expected life gain count 0, got %d", lifeGainCount)
	}

	// Publish life gain event
	bus.Publish(NewEventWithAmount(EventGainedLife, "player1", "source1", "player1", 5))
	if spellCastCount != 1 {
		t.Fatalf("expected spell cast count still 1, got %d", spellCastCount)
	}
	if lifeGainCount != 1 {
		t.Fatalf("expected life gain count 1, got %d", lifeGainCount)
	}

	// Unsubscribe spell cast listener
	bus.UnsubscribeTyped(handle1)

	// Publish spell cast again - should not increment
	bus.Publish(NewEvent(EventSpellCast, "card2", "card2", "player1"))
	if spellCastCount != 1 {
		t.Fatalf("expected spell cast count still 1 after unsubscribe, got %d", spellCastCount)
	}

	// Life gain should still work
	bus.Publish(NewEventWithAmount(EventGainedLife, "player1", "source2", "player1", 3))
	if lifeGainCount != 2 {
		t.Fatalf("expected life gain count 2, got %d", lifeGainCount)
	}

	// Unsubscribe life gain listener
	bus.UnsubscribeTyped(handle2)

	// Publish life gain again - should not increment
	bus.Publish(NewEventWithAmount(EventGainedLife, "player1", "source3", "player1", 2))
	if lifeGainCount != 2 {
		t.Fatalf("expected life gain count still 2 after unsubscribe, got %d", lifeGainCount)
	}
}

func TestEventBusSubscribeAll(t *testing.T) {
	bus := NewEventBus()

	allEventCount := 0
	handle := bus.Subscribe(func(e Event) {
		allEventCount++
	})

	bus.Publish(NewEvent(EventSpellCast, "card1", "card1", "player1"))
	bus.Publish(NewEvent(EventGainedLife, "player1", "source1", "player1"))
	bus.Publish(NewEvent(EventZoneChange, "card2", "card2", "player1"))

	if allEventCount != 3 {
		t.Fatalf("expected all event count 3, got %d", allEventCount)
	}

	bus.Unsubscribe(handle)

	bus.Publish(NewEvent(EventSpellCast, "card3", "card3", "player1"))
	if allEventCount != 3 {
		t.Fatalf("expected all event count still 3 after unsubscribe, got %d", allEventCount)
	}
}

func TestEventIsBatch(t *testing.T) {
	if !EventZoneChangeBatch.IsBatch() {
		t.Fatal("EventZoneChangeBatch should be a batch event")
	}
	if !EventDamagedBatchForAll.IsBatch() {
		t.Fatal("EventDamagedBatchForAll should be a batch event")
	}
	if EventSpellCast.IsBatch() {
		t.Fatal("EventSpellCast should not be a batch event")
	}
	if EventGainedLife.IsBatch() {
		t.Fatal("EventGainedLife should not be a batch event")
	}
}

func TestEventFields(t *testing.T) {
	evt := NewEventWithAmount(EventDamagedPlayer, "player1", "source1", "player1", 5)
	evt.Flag = true // Combat damage
	evt.Data = "combat"
	evt.Zone = 2 // Battlefield
	evt.Targets = []string{"player1", "player2"}
	evt.Metadata["damage_type"] = "combat"
	evt.Description = "Player takes 5 combat damage"

	if evt.Type != EventDamagedPlayer {
		t.Fatalf("expected type EventDamagedPlayer, got %s", evt.Type)
	}
	if evt.Amount != 5 {
		t.Fatalf("expected amount 5, got %d", evt.Amount)
	}
	if !evt.Flag {
		t.Fatal("expected flag true")
	}
	if evt.Data != "combat" {
		t.Fatalf("expected data 'combat', got %s", evt.Data)
	}
	if len(evt.Targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(evt.Targets))
	}
}

func TestEventBusPublishBatch(t *testing.T) {
	bus := NewEventBus()

	count := 0
	bus.Subscribe(func(e Event) {
		count++
	})

	events := []Event{
		NewEvent(EventSpellCast, "card1", "card1", "player1"),
		NewEvent(EventGainedLife, "player1", "source1", "player1"),
		NewEvent(EventZoneChange, "card2", "card2", "player1"),
	}

	bus.PublishBatch(events)

	if count != 3 {
		t.Fatalf("expected count 3 after batch publish, got %d", count)
	}
}

func TestEventTimestamp(t *testing.T) {
	before := time.Now()
	evt := NewEvent(EventSpellCast, "card1", "card1", "player1")
	after := time.Now()

	if evt.Timestamp.Before(before) || evt.Timestamp.After(after) {
		t.Fatal("event timestamp should be between before and after")
	}
}
