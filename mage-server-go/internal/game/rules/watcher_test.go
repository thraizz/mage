package rules

import (
	"testing"
)

func TestWatcherRegistry(t *testing.T) {
	registry := NewWatcherRegistry()

	// Create a simple test watcher
	testWatcher := &testWatcherImpl{
		BaseWatcher: NewBaseWatcher(WatcherScopeGame),
		condition:   false,
	}
	testWatcher.SetKey("TestWatcher")

	registry.AddWatcher(testWatcher)

	// Test retrieval
	retrieved := registry.GetWatcher("TestWatcher")
	if retrieved == nil {
		t.Fatal("should retrieve TestWatcher")
	}

	// Test GetWatchersByScope
	gameWatchers := registry.GetWatchersByScope(WatcherScopeGame)
	if len(gameWatchers) != 1 {
		t.Fatalf("expected 1 game watcher, got %d", len(gameWatchers))
	}

	// Test notification
	event := NewEvent(EventSpellCast, "spell1", "spell1", "player1")
	registry.NotifyWatchers(event)

	if !testWatcher.ConditionMet() {
		t.Fatal("testWatcher should have condition met")
	}

	// Test reset
	registry.ResetWatchers()
	if testWatcher.ConditionMet() {
		t.Fatal("watcher should not have condition met after reset")
	}

	// Test removal
	registry.RemoveWatcher("TestWatcher")
	retrieved = registry.GetWatcher("TestWatcher")
	if retrieved != nil {
		t.Fatal("watcher should be removed")
	}
}

// testWatcherImpl is a simple test watcher implementation
type testWatcherImpl struct {
	*BaseWatcher
	condition bool
}

func (t *testWatcherImpl) Watch(event Event) {
	if event.Type == EventSpellCast {
		t.condition = true
		t.SetCondition(true)
	}
}

func (t *testWatcherImpl) Reset() {
	t.BaseWatcher.Reset()
	t.condition = false
}

func (t *testWatcherImpl) ConditionMet() bool {
	return t.condition
}

func (t *testWatcherImpl) Copy() Watcher {
	return &testWatcherImpl{
		BaseWatcher: NewBaseWatcher(t.GetScope()),
		condition:   t.condition,
	}
}

func TestWatcherScope(t *testing.T) {
	if WatcherScopeGame.String() != "GAME" {
		t.Fatalf("expected GAME, got %s", WatcherScopeGame.String())
	}
	if WatcherScopePlayer.String() != "PLAYER" {
		t.Fatalf("expected PLAYER, got %s", WatcherScopePlayer.String())
	}
	if WatcherScopeCard.String() != "CARD" {
		t.Fatalf("expected CARD, got %s", WatcherScopeCard.String())
	}
}

func TestBaseWatcher(t *testing.T) {
	bw := NewBaseWatcher(WatcherScopeGame)
	bw.SetKey("test_key")

	if bw.GetKey() != "test_key" {
		t.Fatalf("expected test_key, got %s", bw.GetKey())
	}
	if bw.GetScope() != WatcherScopeGame {
		t.Fatalf("expected GAME scope, got %v", bw.GetScope())
	}
	if bw.ConditionMet() {
		t.Fatal("should not have condition met initially")
	}

	bw.SetCondition(true)
	if !bw.ConditionMet() {
		t.Fatal("should have condition met after SetCondition")
	}

	bw.Reset()
	if bw.ConditionMet() {
		t.Fatal("should not have condition met after reset")
	}

	bw.SetControllerID("player1")
	if bw.GetControllerID() != "player1" {
		t.Fatalf("expected player1, got %s", bw.GetControllerID())
	}

	bw.SetSourceID("card1")
	if bw.GetSourceID() != "card1" {
		t.Fatalf("expected card1, got %s", bw.GetSourceID())
	}
}

func TestWatcherRegistryIntegration(t *testing.T) {
	registry := NewWatcherRegistry()
	eventBus := NewEventBus()

	// Wire registry to event bus
	eventBus.Subscribe(func(event Event) {
		registry.NotifyWatchers(event)
	})

	// Add test watcher
	testWatcher := &testWatcherImpl{
		BaseWatcher: NewBaseWatcher(WatcherScopeGame),
		condition:   false,
	}
	testWatcher.SetKey("TestWatcher")
	registry.AddWatcher(testWatcher)

	// Publish spell cast event
	event := NewEvent(EventSpellCast, "spell1", "spell1", "player1")
	eventBus.Publish(event)

	if !testWatcher.ConditionMet() {
		t.Fatal("testWatcher should have condition met")
	}

	// Test reset
	registry.ResetWatchers()
	if testWatcher.ConditionMet() {
		t.Fatal("testWatcher should not have condition met after reset")
	}
}
