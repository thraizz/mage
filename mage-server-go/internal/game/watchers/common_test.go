package watchers

import (
	"testing"
	"time"

	"github.com/magefree/mage-server-go/internal/game/rules"
)

func TestSpellsCastWatcher(t *testing.T) {
	watcher := NewSpellsCastWatcher()

	// Test initial state
	if watcher.ConditionMet() {
		t.Fatal("watcher should not have condition met initially")
	}
	if watcher.GetCount("player1") != 0 {
		t.Fatalf("expected 0 spells cast, got %d", watcher.GetCount("player1"))
	}

	// Watch a spell cast event
	event := rules.NewEvent(rules.EventSpellCast, "spell1", "spell1", "player1")
	watcher.Watch(event)

	if !watcher.ConditionMet() {
		t.Fatal("watcher should have condition met after spell cast")
	}
	if watcher.GetCount("player1") != 1 {
		t.Fatalf("expected 1 spell cast, got %d", watcher.GetCount("player1"))
	}

	// Watch another spell cast
	event2 := rules.NewEvent(rules.EventSpellCast, "spell2", "spell2", "player1")
	watcher.Watch(event2)

	if watcher.GetCount("player1") != 2 {
		t.Fatalf("expected 2 spells cast, got %d", watcher.GetCount("player1"))
	}

	// Test reset
	watcher.Reset()
	if watcher.ConditionMet() {
		t.Fatal("watcher should not have condition met after reset")
	}
	if watcher.GetCount("player1") != 0 {
		t.Fatalf("expected 0 spells cast after reset, got %d", watcher.GetCount("player1"))
	}
}

func TestCreaturesDiedWatcher(t *testing.T) {
	watcher := NewCreaturesDiedWatcher()

	// Test initial state
	if watcher.ConditionMet() {
		t.Fatal("watcher should not have condition met initially")
	}

	// Watch a creature dies event
	event := rules.Event{
		Type:        rules.EventPermanentDies,
		TargetID:    "creature1",
		SourceID:    "creature1",
		Controller:  "player1",
		PlayerID:    "player1",
		Timestamp:   time.Now(),
		Metadata: map[string]string{
			"owner_id": "player1",
		},
	}
	watcher.Watch(event)

	if !watcher.ConditionMet() {
		t.Fatal("watcher should have condition met after creature dies")
	}
	if watcher.GetAmountByController("player1") != 1 {
		t.Fatalf("expected 1 creature died for controller, got %d", watcher.GetAmountByController("player1"))
	}
	if watcher.GetAmountByOwner("player1") != 1 {
		t.Fatalf("expected 1 creature died for owner, got %d", watcher.GetAmountByOwner("player1"))
	}

	// Test reset
	watcher.Reset()
	if watcher.ConditionMet() {
		t.Fatal("watcher should not have condition met after reset")
	}
	if watcher.GetAmountByController("player1") != 0 {
		t.Fatalf("expected 0 creatures died after reset, got %d", watcher.GetAmountByController("player1"))
	}
}

func TestCardsDrawnWatcher(t *testing.T) {
	watcher := NewCardsDrawnWatcher()

	// Watch a card drawn event
	event := rules.NewEvent(rules.EventDrewCard, "card1", "card1", "player1")
	watcher.Watch(event)

	if watcher.GetCount("player1") != 1 {
		t.Fatalf("expected 1 card drawn, got %d", watcher.GetCount("player1"))
	}

	// Watch another card drawn
	event2 := rules.NewEvent(rules.EventDrewCard, "card2", "card2", "player1")
	watcher.Watch(event2)

	if watcher.GetCount("player1") != 2 {
		t.Fatalf("expected 2 cards drawn, got %d", watcher.GetCount("player1"))
	}

	// Test reset
	watcher.Reset()
	if watcher.GetCount("player1") != 0 {
		t.Fatalf("expected 0 cards drawn after reset, got %d", watcher.GetCount("player1"))
	}
}

func TestPermanentsEnteredWatcher(t *testing.T) {
	watcher := NewPermanentsEnteredWatcher()

	// Watch a permanent enters event
	event := rules.NewEvent(rules.EventEntersTheBattlefield, "permanent1", "permanent1", "player1")
	watcher.Watch(event)

	entered := watcher.GetPermanentsEntered("player1")
	if len(entered) != 1 {
		t.Fatalf("expected 1 permanent entered, got %d", len(entered))
	}
	if entered[0] != "permanent1" {
		t.Fatalf("expected permanent1, got %s", entered[0])
	}

	// Test reset
	watcher.Reset()
	entered = watcher.GetPermanentsEntered("player1")
	if len(entered) != 0 {
		t.Fatalf("expected 0 permanents entered after reset, got %d", len(entered))
	}
}

func TestWatcherCopy(t *testing.T) {
	watcher := NewSpellsCastWatcher()
	event := rules.NewEvent(rules.EventSpellCast, "spell1", "spell1", "player1")
	watcher.Watch(event)

	copy := watcher.Copy()
	if copy == nil {
		t.Fatal("copy should not be nil")
	}

	copyWatcher, ok := copy.(*SpellsCastWatcher)
	if !ok {
		t.Fatal("copy should be *SpellsCastWatcher")
	}

	// Copy should have same condition
	if copyWatcher.ConditionMet() != watcher.ConditionMet() {
		t.Fatal("copy should have same condition")
	}

	// Copy should have same data
	if copyWatcher.GetCount("player1") != watcher.GetCount("player1") {
		t.Fatal("copy should have same spell count")
	}

	// Modifying copy shouldn't affect original
	copyWatcher.Watch(rules.NewEvent(rules.EventSpellCast, "spell2", "spell2", "player1"))
	if watcher.GetCount("player1") != 1 {
		t.Fatal("modifying copy shouldn't affect original")
	}
	if copyWatcher.GetCount("player1") != 2 {
		t.Fatal("copy should have updated count")
	}
}
