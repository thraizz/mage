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

func TestWatcherIntegration_SpellCasting(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "watcher-test-spells"
	players := []string{"Alice", "Bob"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Cast a spell
	action := game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "SEND_STRING",
		Data:       "Lightning Bolt",
		Timestamp:  time.Now(),
	}

	err = engine.ProcessAction(gameID, action)
	if err != nil {
		t.Fatalf("Failed to cast spell: %v", err)
	}

	// Get game view to verify spell was cast
	viewRaw, err := engine.GetGameView(gameID, "Alice")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}

	view, ok := viewRaw.(*game.EngineGameView)
	if !ok {
		t.Fatalf("Unexpected view type: %T", viewRaw)
	}

	// Verify spell cast message exists
	spellCastFound := false
	for _, msg := range view.Messages {
		textLower := strings.ToLower(msg.Text)
		if strings.Contains(textLower, "cast") || strings.Contains(textLower, "lightning bolt") {
			spellCastFound = true
			break
		}
	}
	if !spellCastFound {
		t.Error("Expected spell cast message in game view")
	}

	// Verify stack contains the spell
	if len(view.Stack) == 0 {
		t.Error("Expected spell on stack")
	}
}

func TestWatcherIntegration_MultipleWatchers(t *testing.T) {
	// Test multiple watchers tracking different events through event bus
	eventBus := rules.NewEventBus()
	watcherReg := rules.NewWatcherRegistry()

	// Wire watchers to event bus
	eventBus.Subscribe(func(event rules.Event) {
		watcherReg.NotifyWatchers(event)
	})

	// Add all common watchers
	spellsWatcher := watchers.NewSpellsCastWatcher()
	creaturesWatcher := watchers.NewCreaturesDiedWatcher()
	cardsWatcher := watchers.NewCardsDrawnWatcher()
	permanentsWatcher := watchers.NewPermanentsEnteredWatcher()

	watcherReg.AddWatcher(spellsWatcher)
	watcherReg.AddWatcher(creaturesWatcher)
	watcherReg.AddWatcher(cardsWatcher)
	watcherReg.AddWatcher(permanentsWatcher)

	// Publish spell cast event
	spellEvent := rules.NewEvent(rules.EventSpellCast, "spell1", "spell1", "player1")
	eventBus.Publish(spellEvent)

	// Verify spells watcher tracked it
	if !spellsWatcher.ConditionMet() {
		t.Error("SpellsCastWatcher should have condition met")
	}
	if spellsWatcher.GetCount("player1") != 1 {
		t.Errorf("Expected 1 spell cast, got %d", spellsWatcher.GetCount("player1"))
	}

	// Other watchers should not have tracked spell cast
	if creaturesWatcher.ConditionMet() {
		t.Error("CreaturesDiedWatcher should not have condition met")
	}
	if cardsWatcher.ConditionMet() {
		t.Error("CardsDrawnWatcher should not have condition met")
	}
	if permanentsWatcher.ConditionMet() {
		t.Error("PermanentsEnteredWatcher should not have condition met")
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

	// Verify creatures watcher tracked it
	if !creaturesWatcher.ConditionMet() {
		t.Error("CreaturesDiedWatcher should have condition met")
	}
	if creaturesWatcher.GetAmountByController("player1") != 1 {
		t.Errorf("Expected 1 creature died, got %d", creaturesWatcher.GetAmountByController("player1"))
	}

	// Publish card drawn event
	drawEvent := rules.NewEvent(rules.EventDrewCard, "card1", "card1", "player1")
	eventBus.Publish(drawEvent)

	// Verify cards watcher tracked it
	if !cardsWatcher.ConditionMet() {
		t.Error("CardsDrawnWatcher should have condition met")
	}
	if cardsWatcher.GetCount("player1") != 1 {
		t.Errorf("Expected 1 card drawn, got %d", cardsWatcher.GetCount("player1"))
	}

	// Publish permanent enters event
	etbEvent := rules.NewEvent(rules.EventEntersTheBattlefield, "permanent1", "permanent1", "player1")
	eventBus.Publish(etbEvent)

	// Verify permanents watcher tracked it
	if !permanentsWatcher.ConditionMet() {
		t.Error("PermanentsEnteredWatcher should have condition met")
	}
	entered := permanentsWatcher.GetPermanentsEntered("player1")
	if len(entered) != 1 {
		t.Errorf("Expected 1 permanent entered, got %d", len(entered))
	}
}

func TestWatcherIntegration_EndToEnd(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "watcher-test-e2e"
	players := []string{"Alice", "Bob"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Simulate a game flow with multiple events
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

	// 4. Spell resolves (Lightning Bolt is an instant, goes to graveyard not battlefield)
	// The watchers should have tracked:
	// - Spell cast event
	// - Zone change events (hand -> stack -> graveyard)

	// Verify game state
	viewRaw, err := engine.GetGameView(gameID, "Alice")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}

	view, ok := viewRaw.(*game.EngineGameView)
	if !ok {
		t.Fatalf("Unexpected view type: %T", viewRaw)
	}

	// Check that events were logged
	if len(view.Messages) == 0 {
		t.Fatal("Expected game messages")
	}

	// Verify spell was cast (check messages)
	spellCastFound := false
	for _, msg := range view.Messages {
		textLower := strings.ToLower(msg.Text)
		if strings.Contains(textLower, "cast") || strings.Contains(textLower, "lightning bolt") {
			spellCastFound = true
		}
	}
	if !spellCastFound {
		t.Error("Expected spell cast message")
	}

	// Verify stack is empty (spell resolved)
	if len(view.Stack) != 0 {
		t.Errorf("Expected empty stack after resolution, got %d items", len(view.Stack))
	}
}

func TestWatcherIntegration_WatcherReset(t *testing.T) {
	// Test watcher reset functionality
	eventBus := rules.NewEventBus()
	watcherReg := rules.NewWatcherRegistry()

	eventBus.Subscribe(func(event rules.Event) {
		watcherReg.NotifyWatchers(event)
	})

	spellsWatcher := watchers.NewSpellsCastWatcher()
	watcherReg.AddWatcher(spellsWatcher)

	// Cast multiple spells
	for i := 0; i < 3; i++ {
		spellID := "spell" + string(rune('1'+i))
		event := rules.NewEvent(rules.EventSpellCast, spellID, spellID, "player1")
		eventBus.Publish(event)
	}

	// Verify spells tracked
	if spellsWatcher.GetCount("player1") != 3 {
		t.Errorf("Expected 3 spells cast, got %d", spellsWatcher.GetCount("player1"))
	}
	if !spellsWatcher.ConditionMet() {
		t.Error("Watcher should have condition met")
	}

	// Reset watchers
	watcherReg.ResetWatchers()

	// Verify reset
	if spellsWatcher.ConditionMet() {
		t.Error("Watcher should not have condition met after reset")
	}
	if spellsWatcher.GetCount("player1") != 0 {
		t.Errorf("Expected 0 spells after reset, got %d", spellsWatcher.GetCount("player1"))
	}

	// Cast another spell after reset
	event := rules.NewEvent(rules.EventSpellCast, "spell4", "spell4", "player1")
	eventBus.Publish(event)

	// Should track new spell
	if spellsWatcher.GetCount("player1") != 1 {
		t.Errorf("Expected 1 spell after reset, got %d", spellsWatcher.GetCount("player1"))
	}
}

func TestWatcherIntegration_CreatureDies(t *testing.T) {
	// Test creature dies watcher
	eventBus := rules.NewEventBus()
	watcherReg := rules.NewWatcherRegistry()

	eventBus.Subscribe(func(event rules.Event) {
		watcherReg.NotifyWatchers(event)
	})

	creaturesWatcher := watchers.NewCreaturesDiedWatcher()
	watcherReg.AddWatcher(creaturesWatcher)

	// Simulate creature entering battlefield first
	etbEvent := rules.NewEvent(rules.EventEntersTheBattlefield, "creature1", "creature1", "player1")
	eventBus.Publish(etbEvent)

	// Creature dies (goes to graveyard from battlefield)
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

	// Verify watcher tracked it
	if !creaturesWatcher.ConditionMet() {
		t.Error("CreaturesDiedWatcher should have condition met")
	}
	if creaturesWatcher.GetAmountByController("player1") != 1 {
		t.Errorf("Expected 1 creature died for controller, got %d", creaturesWatcher.GetAmountByController("player1"))
	}
	if creaturesWatcher.GetAmountByOwner("player1") != 1 {
		t.Errorf("Expected 1 creature died for owner, got %d", creaturesWatcher.GetAmountByOwner("player1"))
	}
	if creaturesWatcher.GetTotalAmount() != 1 {
		t.Errorf("Expected 1 total creature died, got %d", creaturesWatcher.GetTotalAmount())
	}

	// Test multiple creatures dying
	dieEvent2 := rules.Event{
		Type:        rules.EventPermanentDies,
		TargetID:    "creature2",
		SourceID:    "creature2",
		Controller:  "player1",
		PlayerID:    "player1",
		Timestamp:   time.Now(),
		Metadata: map[string]string{
			"owner_id": "player1",
		},
	}
	eventBus.Publish(dieEvent2)

	if creaturesWatcher.GetAmountByController("player1") != 2 {
		t.Errorf("Expected 2 creatures died, got %d", creaturesWatcher.GetAmountByController("player1"))
	}
	if creaturesWatcher.GetTotalAmount() != 2 {
		t.Errorf("Expected 2 total creatures died, got %d", creaturesWatcher.GetTotalAmount())
	}
}

func TestWatcherIntegration_EventBusWatcherIntegration(t *testing.T) {
	// Test that watchers properly integrate with the event bus
	eventBus := rules.NewEventBus()
	watcherReg := rules.NewWatcherRegistry()

	// Wire watchers to event bus
	eventBus.Subscribe(func(event rules.Event) {
		watcherReg.NotifyWatchers(event)
	})

	// Add watchers
	spellsWatcher := watchers.NewSpellsCastWatcher()
	creaturesWatcher := watchers.NewCreaturesDiedWatcher()
	cardsWatcher := watchers.NewCardsDrawnWatcher()
	permanentsWatcher := watchers.NewPermanentsEnteredWatcher()

	watcherReg.AddWatcher(spellsWatcher)
	watcherReg.AddWatcher(creaturesWatcher)
	watcherReg.AddWatcher(cardsWatcher)
	watcherReg.AddWatcher(permanentsWatcher)

	// Publish spell cast event
	spellEvent := rules.NewEvent(rules.EventSpellCast, "spell1", "spell1", "player1")
	eventBus.Publish(spellEvent)

	if !spellsWatcher.ConditionMet() {
		t.Error("SpellsCastWatcher should have condition met")
	}
	if spellsWatcher.GetCount("player1") != 1 {
		t.Errorf("Expected 1 spell cast, got %d", spellsWatcher.GetCount("player1"))
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

	if !creaturesWatcher.ConditionMet() {
		t.Error("CreaturesDiedWatcher should have condition met")
	}
	if creaturesWatcher.GetAmountByController("player1") != 1 {
		t.Errorf("Expected 1 creature died, got %d", creaturesWatcher.GetAmountByController("player1"))
	}

	// Publish card drawn event
	drawEvent := rules.NewEvent(rules.EventDrewCard, "card1", "card1", "player1")
	eventBus.Publish(drawEvent)

	if !cardsWatcher.ConditionMet() {
		t.Error("CardsDrawnWatcher should have condition met")
	}
	if cardsWatcher.GetCount("player1") != 1 {
		t.Errorf("Expected 1 card drawn, got %d", cardsWatcher.GetCount("player1"))
	}

	// Publish permanent enters event
	etbEvent := rules.NewEvent(rules.EventEntersTheBattlefield, "permanent1", "permanent1", "player1")
	eventBus.Publish(etbEvent)

	if !permanentsWatcher.ConditionMet() {
		t.Error("PermanentsEnteredWatcher should have condition met")
	}
	entered := permanentsWatcher.GetPermanentsEntered("player1")
	if len(entered) != 1 {
		t.Errorf("Expected 1 permanent entered, got %d", len(entered))
	}

	// Test reset
	watcherReg.ResetWatchers()

	if spellsWatcher.ConditionMet() {
		t.Error("SpellsCastWatcher should not have condition met after reset")
	}
	if creaturesWatcher.ConditionMet() {
		t.Error("CreaturesDiedWatcher should not have condition met after reset")
	}
	if cardsWatcher.ConditionMet() {
		t.Error("CardsDrawnWatcher should not have condition met after reset")
	}
	if permanentsWatcher.ConditionMet() {
		t.Error("PermanentsEnteredWatcher should not have condition met after reset")
	}

	// Verify counts are reset
	if spellsWatcher.GetCount("player1") != 0 {
		t.Errorf("Expected 0 spells cast after reset, got %d", spellsWatcher.GetCount("player1"))
	}
	if creaturesWatcher.GetAmountByController("player1") != 0 {
		t.Errorf("Expected 0 creatures died after reset, got %d", creaturesWatcher.GetAmountByController("player1"))
	}
	if cardsWatcher.GetCount("player1") != 0 {
		t.Errorf("Expected 0 cards drawn after reset, got %d", cardsWatcher.GetCount("player1"))
	}
	if len(permanentsWatcher.GetPermanentsEntered("player1")) != 0 {
		t.Errorf("Expected 0 permanents entered after reset, got %d", len(permanentsWatcher.GetPermanentsEntered("player1")))
	}
}

func TestWatcherIntegration_MultiplePlayers(t *testing.T) {
	eventBus := rules.NewEventBus()
	watcherReg := rules.NewWatcherRegistry()

	eventBus.Subscribe(func(event rules.Event) {
		watcherReg.NotifyWatchers(event)
	})

	spellsWatcher := watchers.NewSpellsCastWatcher()
	watcherReg.AddWatcher(spellsWatcher)

	// Player 1 casts spells
	event1 := rules.NewEvent(rules.EventSpellCast, "spell1", "spell1", "player1")
	eventBus.Publish(event1)

	event2 := rules.NewEvent(rules.EventSpellCast, "spell2", "spell2", "player1")
	eventBus.Publish(event2)

	// Player 2 casts a spell
	event3 := rules.NewEvent(rules.EventSpellCast, "spell3", "spell3", "player2")
	eventBus.Publish(event3)

	// Verify counts per player
	if spellsWatcher.GetCount("player1") != 2 {
		t.Errorf("Expected 2 spells cast by player1, got %d", spellsWatcher.GetCount("player1"))
	}
	if spellsWatcher.GetCount("player2") != 1 {
		t.Errorf("Expected 1 spell cast by player2, got %d", spellsWatcher.GetCount("player2"))
	}

	// Verify all spells tracked
	allSpells := spellsWatcher.GetSpellsCast("player1")
	if len(allSpells) != 2 {
		t.Errorf("Expected 2 spells in list, got %d", len(allSpells))
	}
	if allSpells[0] != "spell1" || allSpells[1] != "spell2" {
		t.Error("Spell IDs not tracked correctly")
	}
}

func TestWatcherIntegration_WatcherScopes(t *testing.T) {
	watcherReg := rules.NewWatcherRegistry()

	// Game scope watcher
	gameWatcher := watchers.NewSpellsCastWatcher()
	watcherReg.AddWatcher(gameWatcher)

	// Player scope watcher (would need a player-specific watcher implementation)
	playerWatcher := &testPlayerWatcher{
		BaseWatcher: rules.NewBaseWatcher(rules.WatcherScopePlayer),
		playerID:    "player1",
		spellCount:  0,
	}
	playerWatcher.SetControllerID("player1")
	playerWatcher.SetKey("player1_TestPlayerWatcher")
	watcherReg.AddWatcher(playerWatcher)

	// Card scope watcher (would need a card-specific watcher implementation)
	cardWatcher := &testCardWatcher{
		BaseWatcher: rules.NewBaseWatcher(rules.WatcherScopeCard),
		cardID:      "card1",
		triggered:   false,
	}
	cardWatcher.SetSourceID("card1")
	cardWatcher.SetKey("card1_TestCardWatcher")
	watcherReg.AddWatcher(cardWatcher)

	// Test retrieval by scope
	gameWatchers := watcherReg.GetWatchersByScope(rules.WatcherScopeGame)
	if len(gameWatchers) != 1 {
		t.Errorf("Expected 1 game watcher, got %d", len(gameWatchers))
	}

	playerWatchers := watcherReg.GetWatchersByScope(rules.WatcherScopePlayer)
	if len(playerWatchers) != 1 {
		t.Errorf("Expected 1 player watcher, got %d", len(playerWatchers))
	}

	cardWatchers := watcherReg.GetWatchersByScope(rules.WatcherScopeCard)
	if len(cardWatchers) != 1 {
		t.Errorf("Expected 1 card watcher, got %d", len(cardWatchers))
	}

	// Test reset by scope
	watcherReg.ResetWatchersByScope(rules.WatcherScopeGame)
	if gameWatcher.ConditionMet() {
		t.Error("Game watcher should be reset")
	}
	// Player and card watchers should still have their state
}

func TestWatcherIntegration_ConcurrentWatchers(t *testing.T) {
	eventBus := rules.NewEventBus()
	watcherReg := rules.NewWatcherRegistry()

	eventBus.Subscribe(func(event rules.Event) {
		watcherReg.NotifyWatchers(event)
	})

	// Add multiple watchers of different types
	spellsWatcher := watchers.NewSpellsCastWatcher()
	creaturesWatcher := watchers.NewCreaturesDiedWatcher()
	cardsWatcher := watchers.NewCardsDrawnWatcher()

	watcherReg.AddWatcher(spellsWatcher)
	watcherReg.AddWatcher(creaturesWatcher)
	watcherReg.AddWatcher(cardsWatcher)

	// Publish spell cast event
	spellEvent := rules.NewEvent(rules.EventSpellCast, "spell1", "spell1", "player1")
	eventBus.Publish(spellEvent)

	// Spells watcher should have tracked it
	if spellsWatcher.GetCount("player1") != 1 {
		t.Errorf("SpellsWatcher: Expected 1 spell, got %d", spellsWatcher.GetCount("player1"))
	}
	// Other watchers should not have tracked spell cast
	if creaturesWatcher.ConditionMet() {
		t.Error("CreaturesWatcher should not have condition met")
	}
	if cardsWatcher.ConditionMet() {
		t.Error("CardsWatcher should not have condition met")
	}

	// Publish multiple different events
	drawEvent := rules.NewEvent(rules.EventDrewCard, "card1", "card1", "player1")
	eventBus.Publish(drawEvent)

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

	// Verify each watcher tracked its relevant event
	if spellsWatcher.GetCount("player1") != 1 {
		t.Errorf("SpellsWatcher: Expected 1 spell, got %d", spellsWatcher.GetCount("player1"))
	}
	if !cardsWatcher.ConditionMet() {
		t.Error("CardsWatcher should have condition met")
	}
	if cardsWatcher.GetCount("player1") != 1 {
		t.Errorf("CardsWatcher: Expected 1 card drawn, got %d", cardsWatcher.GetCount("player1"))
	}
	if !creaturesWatcher.ConditionMet() {
		t.Error("CreaturesWatcher should have condition met")
	}
	if creaturesWatcher.GetAmountByController("player1") != 1 {
		t.Errorf("CreaturesWatcher: Expected 1 creature died, got %d", creaturesWatcher.GetAmountByController("player1"))
	}
}

func TestWatcherIntegration_WatcherCopy(t *testing.T) {
	watcher := watchers.NewSpellsCastWatcher()

	event := rules.NewEvent(rules.EventSpellCast, "spell1", "spell1", "player1")
	watcher.Watch(event)

	// Copy the watcher
	copy := watcher.Copy()
	if copy == nil {
		t.Fatal("Copy should not be nil")
	}

	copyWatcher, ok := copy.(*watchers.SpellsCastWatcher)
	if !ok {
		t.Fatal("Copy should be *SpellsCastWatcher")
	}

	// Verify copy has same state
	if copyWatcher.ConditionMet() != watcher.ConditionMet() {
		t.Error("Copy should have same condition")
	}
	if copyWatcher.GetCount("player1") != watcher.GetCount("player1") {
		t.Error("Copy should have same count")
	}

	// Modify copy - shouldn't affect original
	copyWatcher.Watch(rules.NewEvent(rules.EventSpellCast, "spell2", "spell2", "player1"))
	if watcher.GetCount("player1") != 1 {
		t.Error("Modifying copy shouldn't affect original")
	}
	if copyWatcher.GetCount("player1") != 2 {
		t.Error("Copy should have updated count")
	}
}

// Helper functions and test watchers

// testPlayerWatcher is a test watcher for player scope
type testPlayerWatcher struct {
	*rules.BaseWatcher
	playerID   string
	spellCount int
}

func (t *testPlayerWatcher) Watch(event rules.Event) {
	if event.Type == rules.EventSpellCast && event.PlayerID == t.playerID {
		t.spellCount++
		t.SetCondition(true)
	}
}

func (t *testPlayerWatcher) Reset() {
	t.BaseWatcher.Reset()
	t.spellCount = 0
}

func (t *testPlayerWatcher) ConditionMet() bool {
	return t.spellCount > 0
}

func (t *testPlayerWatcher) Copy() rules.Watcher {
	return &testPlayerWatcher{
		BaseWatcher: rules.NewBaseWatcher(t.GetScope()),
		playerID:    t.playerID,
		spellCount:  t.spellCount,
	}
}

// testCardWatcher is a test watcher for card scope
type testCardWatcher struct {
	*rules.BaseWatcher
	cardID    string
	triggered bool
}

func (t *testCardWatcher) Watch(event rules.Event) {
	if event.Type == rules.EventSpellCast && event.SourceID == t.cardID {
		t.triggered = true
		t.SetCondition(true)
	}
}

func (t *testCardWatcher) Reset() {
	t.BaseWatcher.Reset()
	t.triggered = false
}

func (t *testCardWatcher) ConditionMet() bool {
	return t.triggered
}

func (t *testCardWatcher) Copy() rules.Watcher {
	return &testCardWatcher{
		BaseWatcher: rules.NewBaseWatcher(t.GetScope()),
		cardID:      t.cardID,
		triggered:   t.triggered,
	}
}
