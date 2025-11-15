package counters

import (
	"fmt"
	"time"

	"github.com/magefree/mage-server-go/internal/game/rules"
)

// CounterOperations provides counter manipulation functions for game objects.
type CounterOperations struct {
	eventBus *rules.EventBus
}

// NewCounterOperations creates a new CounterOperations instance.
func NewCounterOperations(eventBus *rules.EventBus) *CounterOperations {
	return &CounterOperations{
		eventBus: eventBus,
	}
}

// AddCounterToCard adds counters to a card and emits appropriate events.
func (co *CounterOperations) AddCounterToCard(cardID string, counter *Counter, controllerID string, timestamp time.Time) {
	if counter == nil || counter.Count <= 0 {
		return
	}

	// Emit counter added event
	co.eventBus.Publish(rules.Event{
		Type:       rules.EventCounterAdded,
		ID:         fmt.Sprintf("event-counter-added-%s-%d", cardID, timestamp.UnixNano()),
		TargetID:   cardID,
		SourceID:   cardID,
		Controller: controllerID,
		PlayerID:   controllerID,
		Amount:     counter.Count,
		Data:       counter.Name,
		Timestamp:  timestamp,
		Metadata: map[string]string{
			"counter_name": counter.Name,
			"counter_count": fmt.Sprintf("%d", counter.Count),
		},
		Description: fmt.Sprintf("Added %d %s counter(s) to %s", counter.Count, counter.Name, cardID),
	})
}

// RemoveCounterFromCard removes counters from a card and emits appropriate events.
func (co *CounterOperations) RemoveCounterFromCard(cardID string, counterName string, amount int, controllerID string, timestamp time.Time) bool {
	if amount <= 0 {
		return false
	}

	// Emit counter removed event
	co.eventBus.Publish(rules.Event{
		Type:       rules.EventCounterRemoved,
		ID:         fmt.Sprintf("event-counter-removed-%s-%d", cardID, timestamp.UnixNano()),
		TargetID:   cardID,
		SourceID:   cardID,
		Controller: controllerID,
		PlayerID:   controllerID,
		Amount:     amount,
		Data:       counterName,
		Timestamp:  timestamp,
		Metadata: map[string]string{
			"counter_name": counterName,
			"counter_count": fmt.Sprintf("%d", amount),
		},
		Description: fmt.Sprintf("Removed %d %s counter(s) from %s", amount, counterName, cardID),
	})

	return true
}

// AddCountersToCard adds multiple counters to a card (batch operation).
func (co *CounterOperations) AddCountersToCard(cardID string, counters []*Counter, controllerID string, timestamp time.Time) {
	if len(counters) == 0 {
		return
	}

	// Emit batch counters added event
	co.eventBus.Publish(rules.Event{
		Type:       rules.EventCountersAdded,
		ID:         fmt.Sprintf("event-counters-added-%s-%d", cardID, timestamp.UnixNano()),
		TargetID:   cardID,
		SourceID:   cardID,
		Controller: controllerID,
		PlayerID:   controllerID,
		Timestamp:  timestamp,
		Metadata: map[string]string{
			"counter_count": fmt.Sprintf("%d", len(counters)),
		},
		Description: fmt.Sprintf("Added %d counter(s) to %s", len(counters), cardID),
	})

	// Also emit individual counter added events
	for _, counter := range counters {
		co.AddCounterToCard(cardID, counter, controllerID, timestamp)
	}
}

// RemoveCountersFromCard removes multiple counters from a card (batch operation).
func (co *CounterOperations) RemoveCountersFromCard(cardID string, counters map[string]int, controllerID string, timestamp time.Time) {
	if len(counters) == 0 {
		return
	}

	// Emit batch counters removed event
	co.eventBus.Publish(rules.Event{
		Type:       rules.EventCountersRemoved,
		ID:         fmt.Sprintf("event-counters-removed-%s-%d", cardID, timestamp.UnixNano()),
		TargetID:   cardID,
		SourceID:   cardID,
		Controller: controllerID,
		PlayerID:   controllerID,
		Timestamp:  timestamp,
		Metadata: map[string]string{
			"counter_count": fmt.Sprintf("%d", len(counters)),
		},
		Description: fmt.Sprintf("Removed counters from %s", cardID),
	})

	// Also emit individual counter removed events
	for counterName, amount := range counters {
		co.RemoveCounterFromCard(cardID, counterName, amount, controllerID, timestamp)
	}
}
