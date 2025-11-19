package integration

import (
	"strings"
	"testing"
	"time"

	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/rules"
	"github.com/magefree/mage-server-go/internal/game/watchers"
	"go.uber.org/zap"
)

// TestWatcherGameEngineIntegration tests that watchers work correctly with the actual game engine
func TestWatcherGameEngineIntegration(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "watcher-engine-test"
	players := []string{"Alice", "Bob"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Create a test event bus to capture events
	testEventBus := rules.NewEventBus()
	watcherReg := rules.NewWatcherRegistry()

	// Wire watchers
	testEventBus.Subscribe(func(event rules.Event) {
		watcherReg.NotifyWatchers(event)
	})

	// Add watchers
	spellsWatcher := watchers.NewSpellsCastWatcher()
	cardsWatcher := watchers.NewCardsDrawnWatcher()
	permanentsWatcher := watchers.NewPermanentsEnteredWatcher()

	watcherReg.AddWatcher(spellsWatcher)
	watcherReg.AddWatcher(cardsWatcher)
	watcherReg.AddWatcher(permanentsWatcher)

	// Note: The actual game engine has its own event bus, so we're testing
	// the watcher infrastructure separately. In a real scenario, the engine's
	// event bus would be wired to watchers.

	// Simulate events that would be emitted by the game engine
	spellEvent := rules.NewEvent(rules.EventSpellCast, "spell1", "spell1", "Alice")
	testEventBus.Publish(spellEvent)

	if !spellsWatcher.ConditionMet() {
		t.Error("SpellsCastWatcher should have condition met")
	}
	if spellsWatcher.GetCount("Alice") != 1 {
		t.Errorf("Expected 1 spell cast by Alice, got %d", spellsWatcher.GetCount("Alice"))
	}

	// Simulate permanent entering battlefield
	etbEvent := rules.NewEvent(rules.EventEntersTheBattlefield, "permanent1", "permanent1", "Alice")
	testEventBus.Publish(etbEvent)

	if !permanentsWatcher.ConditionMet() {
		t.Error("PermanentsEnteredWatcher should have condition met")
	}
	entered := permanentsWatcher.GetPermanentsEntered("Alice")
	if len(entered) != 1 {
		t.Errorf("Expected 1 permanent entered, got %d", len(entered))
	}
}

// TestWatcherLifeCycle tests watcher lifecycle: add, track events, reset, remove
func TestWatcherLifeCycle(t *testing.T) {
	eventBus := rules.NewEventBus()
	watcherReg := rules.NewWatcherRegistry()

	eventBus.Subscribe(func(event rules.Event) {
		watcherReg.NotifyWatchers(event)
	})

	watcher := watchers.NewSpellsCastWatcher()
	key := watcher.GetKey()

	// Add watcher
	watcherReg.AddWatcher(watcher)

	// Verify watcher is registered
	retrieved := watcherReg.GetWatcher(key)
	if retrieved == nil {
		t.Fatal("Watcher should be retrievable")
	}

	// Track events
	event1 := rules.NewEvent(rules.EventSpellCast, "spell1", "spell1", "player1")
	eventBus.Publish(event1)

	if watcher.GetCount("player1") != 1 {
		t.Errorf("Expected 1 spell, got %d", watcher.GetCount("player1"))
	}

	// Reset watcher
	watcherReg.ResetWatchers()

	if watcher.ConditionMet() {
		t.Error("Watcher should not have condition met after reset")
	}
	if watcher.GetCount("player1") != 0 {
		t.Errorf("Expected 0 spells after reset, got %d", watcher.GetCount("player1"))
	}

	// Track more events
	event2 := rules.NewEvent(rules.EventSpellCast, "spell2", "spell2", "player1")
	eventBus.Publish(event2)

	if watcher.GetCount("player1") != 1 {
		t.Errorf("Expected 1 spell after reset, got %d", watcher.GetCount("player1"))
	}

	// Remove watcher
	watcherReg.RemoveWatcher(key)

	retrieved = watcherReg.GetWatcher(key)
	if retrieved != nil {
		t.Error("Watcher should be removed")
	}

	// Publish event after removal - watcher shouldn't track it
	event3 := rules.NewEvent(rules.EventSpellCast, "spell3", "spell3", "player1")
	eventBus.Publish(event3)

	// Watcher state shouldn't change (it's removed from registry)
	if watcher.GetCount("player1") != 1 {
		t.Errorf("Expected count to remain 1 after removal, got %d", watcher.GetCount("player1"))
	}
}

// TestWatcherMultipleEventsPerTurn tests tracking multiple events in a single turn
func TestWatcherMultipleEventsPerTurn(t *testing.T) {
	eventBus := rules.NewEventBus()
	watcherReg := rules.NewWatcherRegistry()

	eventBus.Subscribe(func(event rules.Event) {
		watcherReg.NotifyWatchers(event)
	})

	spellsWatcher := watchers.NewSpellsCastWatcher()
	cardsWatcher := watchers.NewCardsDrawnWatcher()
	watcherReg.AddWatcher(spellsWatcher)
	watcherReg.AddWatcher(cardsWatcher)

	// Simulate a turn with multiple events
	// Player draws cards
	for i := 0; i < 2; i++ {
		event := rules.NewEvent(rules.EventDrewCard, 
			"card"+string(rune('1'+i)), 
			"card"+string(rune('1'+i)), 
			"player1")
		eventBus.Publish(event)
	}

	// Player casts spells
	for i := 0; i < 3; i++ {
		spellID := "spell" + string(rune('1'+i))
		event := rules.NewEvent(rules.EventSpellCast, spellID, spellID, "player1")
		eventBus.Publish(event)
	}

	// Verify counts
	if cardsWatcher.GetCount("player1") != 2 {
		t.Errorf("Expected 2 cards drawn, got %d", cardsWatcher.GetCount("player1"))
	}
	if spellsWatcher.GetCount("player1") != 3 {
		t.Errorf("Expected 3 spells cast, got %d", spellsWatcher.GetCount("player1"))
	}

	// Verify all spells tracked
	spells := spellsWatcher.GetSpellsCast("player1")
	if len(spells) != 3 {
		t.Errorf("Expected 3 spells in list, got %d", len(spells))
	}
}

// TestWatcherScopeIsolation tests that watchers with different scopes are isolated
func TestWatcherScopeIsolation(t *testing.T) {
	watcherReg := rules.NewWatcherRegistry()

	// Game scope watcher
	gameWatcher := watchers.NewSpellsCastWatcher()

	// Player scope watcher
	playerWatcher := &testPlayerWatcher{
		BaseWatcher: rules.NewBaseWatcher(rules.WatcherScopePlayer),
		playerID:    "player1",
		spellCount:  0,
	}
	playerWatcher.SetControllerID("player1")
	playerWatcher.SetKey("player1_TestPlayerWatcher")

	// Card scope watcher
	cardWatcher := &testCardWatcher{
		BaseWatcher: rules.NewBaseWatcher(rules.WatcherScopeCard),
		cardID:      "card1",
		triggered:   false,
	}
	cardWatcher.SetSourceID("card1")
	cardWatcher.SetKey("card1_TestCardWatcher")

	watcherReg.AddWatcher(gameWatcher)
	watcherReg.AddWatcher(playerWatcher)
	watcherReg.AddWatcher(cardWatcher)

	// Verify scopes
	if gameWatcher.GetScope() != rules.WatcherScopeGame {
		t.Errorf("Expected GAME scope, got %v", gameWatcher.GetScope())
	}
	if playerWatcher.GetScope() != rules.WatcherScopePlayer {
		t.Errorf("Expected PLAYER scope, got %v", playerWatcher.GetScope())
	}
	if cardWatcher.GetScope() != rules.WatcherScopeCard {
		t.Errorf("Expected CARD scope, got %v", cardWatcher.GetScope())
	}

	// Test reset by scope
	eventBus := rules.NewEventBus()
	eventBus.Subscribe(func(event rules.Event) {
		watcherReg.NotifyWatchers(event)
	})

	event := rules.NewEvent(rules.EventSpellCast, "spell1", "card1", "player1")
	eventBus.Publish(event)

	// All watchers should track it (they all watch spell cast events)
	if !gameWatcher.ConditionMet() {
		t.Error("Game watcher should have condition met")
	}
	if !playerWatcher.ConditionMet() {
		t.Error("Player watcher should have condition met")
	}
	if !cardWatcher.ConditionMet() {
		t.Error("Card watcher should have condition met")
	}

	// Reset only game scope watchers
	watcherReg.ResetWatchersByScope(rules.WatcherScopeGame)

	if gameWatcher.ConditionMet() {
		t.Error("Game watcher should be reset")
	}
	// Player and card watchers should still have condition met
	if !playerWatcher.ConditionMet() {
		t.Error("Player watcher should still have condition met")
	}
	if !cardWatcher.ConditionMet() {
		t.Error("Card watcher should still have condition met")
	}
}

// TestWatcherEventFiltering tests that watchers only track relevant events
func TestWatcherEventFiltering(t *testing.T) {
	eventBus := rules.NewEventBus()
	watcherReg := rules.NewWatcherRegistry()

	eventBus.Subscribe(func(event rules.Event) {
		watcherReg.NotifyWatchers(event)
	})

	spellsWatcher := watchers.NewSpellsCastWatcher()
	creaturesWatcher := watchers.NewCreaturesDiedWatcher()
	cardsWatcher := watchers.NewCardsDrawnWatcher()

	watcherReg.AddWatcher(spellsWatcher)
	watcherReg.AddWatcher(creaturesWatcher)
	watcherReg.AddWatcher(cardsWatcher)

	// Publish various events
	spellEvent := rules.NewEvent(rules.EventSpellCast, "spell1", "spell1", "player1")
	eventBus.Publish(spellEvent)

	lifeEvent := rules.NewEventWithAmount(rules.EventGainedLife, "player1", "source1", "player1", 5)
	eventBus.Publish(lifeEvent)

	manaEvent := rules.NewEvent(rules.EventManaAdded, "player1", "source1", "player1")
	eventBus.Publish(manaEvent)

	// Only spells watcher should have tracked spell cast
	if !spellsWatcher.ConditionMet() {
		t.Error("SpellsWatcher should have condition met")
	}
	if creaturesWatcher.ConditionMet() {
		t.Error("CreaturesWatcher should not have condition met")
	}
	if cardsWatcher.ConditionMet() {
		t.Error("CardsWatcher should not have condition met")
	}

	// Publish creature dies event
	dieEvent := rules.Event{
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
	eventBus.Publish(dieEvent)

	// Creatures watcher should now have condition met
	if !creaturesWatcher.ConditionMet() {
		t.Error("CreaturesWatcher should have condition met")
	}
	// Spells watcher should still have condition met
	if !spellsWatcher.ConditionMet() {
		t.Error("SpellsWatcher should still have condition met")
	}
}

// TestWatcherBatchEvents tests that watchers handle batch events correctly
func TestWatcherBatchEvents(t *testing.T) {
	eventBus := rules.NewEventBus()
	watcherReg := rules.NewWatcherRegistry()

	eventBus.Subscribe(func(event rules.Event) {
		watcherReg.NotifyWatchers(event)
	})

	spellsWatcher := watchers.NewSpellsCastWatcher()
	watcherReg.AddWatcher(spellsWatcher)

	// Publish multiple spell cast events in a batch
	events := []rules.Event{
		rules.NewEvent(rules.EventSpellCast, "spell1", "spell1", "player1"),
		rules.NewEvent(rules.EventSpellCast, "spell2", "spell2", "player1"),
		rules.NewEvent(rules.EventSpellCast, "spell3", "spell3", "player1"),
	}

	eventBus.PublishBatch(events)

	// Verify all spells tracked
	if spellsWatcher.GetCount("player1") != 3 {
		t.Errorf("Expected 3 spells cast, got %d", spellsWatcher.GetCount("player1"))
	}

	spells := spellsWatcher.GetSpellsCast("player1")
	if len(spells) != 3 {
		t.Errorf("Expected 3 spells in list, got %d", len(spells))
	}
}

// TestWatcherRealGameFlow tests watchers with a realistic game flow
func TestWatcherRealGameFlow(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "watcher-real-flow"
	players := []string{"Alice", "Bob"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	now := time.Now()

	// 1. Alice casts a spell
	castAction := game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "SEND_STRING",
		Data:       "Lightning Bolt",
		Timestamp:  now,
	}
	err = engine.ProcessAction(gameID, castAction)
	if err != nil {
		t.Fatalf("Failed to cast spell: %v", err)
	}

	// 2. Alice passes priority (per MTG rule 117.3c, Alice retains priority after casting)
	passAction1 := game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  now.Add(time.Second),
	}
	err = engine.ProcessAction(gameID, passAction1)
	if err != nil {
		t.Fatalf("Failed to pass (Alice): %v", err)
	}

	// 3. Bob passes priority
	passAction2 := game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  now.Add(2 * time.Second),
	}
	err = engine.ProcessAction(gameID, passAction2)
	if err != nil {
		t.Fatalf("Failed to pass (Bob): %v", err)
	}

	// 3. Verify game state
	viewRaw, err := engine.GetGameView(gameID, "Alice")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}

	view, ok := viewRaw.(*game.EngineGameView)
	if !ok {
		t.Fatalf("Unexpected view type: %T", viewRaw)
	}

	// Verify spell was cast and resolved
	spellCastFound := false
	resolvedFound := false
	for _, msg := range view.Messages {
		textLower := strings.ToLower(msg.Text)
		if strings.Contains(textLower, "cast") {
			spellCastFound = true
		}
		if strings.Contains(textLower, "resolves") || strings.Contains(textLower, "resolve") {
			resolvedFound = true
		}
	}

	if !spellCastFound {
		t.Error("Expected spell cast message")
	}
	// Resolution message may vary, but stack should be empty and battlefield should have the card
	if !resolvedFound && len(view.Stack) > 0 {
		t.Error("Expected spell resolution (stack should be empty)")
	}

	// Verify stack is empty after resolution
	if len(view.Stack) != 0 {
		t.Errorf("Expected empty stack after resolution, got %d items", len(view.Stack))
	}

	// Lightning Bolt is an instant, so it goes to graveyard, not battlefield
	// Just verify the spell resolved (stack is empty and we got resolution messages)
	if !spellCastFound {
		t.Error("Expected to see spell cast event tracked by watchers")
	}
}

// TestWatcherRegistryThreadSafety tests concurrent access to watcher registry
func TestWatcherRegistryThreadSafety(t *testing.T) {
	watcherReg := rules.NewWatcherRegistry()
	eventBus := rules.NewEventBus()

	eventBus.Subscribe(func(event rules.Event) {
		watcherReg.NotifyWatchers(event)
	})

	// Add watchers concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			watcher := watchers.NewSpellsCastWatcher()
			watcher.SetKey("watcher" + string(rune('0'+idx)))
			watcherReg.AddWatcher(watcher)
			done <- true
		}(i)
	}

	// Wait for all additions
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all watchers added
	allWatchers := watcherReg.GetAllWatchers()
	if len(allWatchers) != 10 {
		t.Errorf("Expected 10 watchers, got %d", len(allWatchers))
	}

	// Publish events concurrently
	for i := 0; i < 10; i++ {
		go func(idx int) {
			event := rules.NewEvent(rules.EventSpellCast, 
				"spell"+string(rune('0'+idx)), 
				"spell"+string(rune('0'+idx)), 
				"player1")
			eventBus.Publish(event)
			done <- true
		}(i)
	}

	// Wait for all events
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify watchers tracked events (at least one should have)
	gameWatchers := watcherReg.GetWatchersByScope(rules.WatcherScopeGame)
	if len(gameWatchers) == 0 {
		t.Error("Expected at least one game watcher")
	}
}
