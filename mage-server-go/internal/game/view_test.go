package game

import (
	"fmt"
	"sync"
	"testing"

	"go.uber.org/zap"
)

// view_test.go guards the fix for the GetGameView data race documented on
// buildGameView.
//
// The view used to alias live GameState memory by pointer while GetGameView's
// deferred RUnlock released the lock on return, so every caller read the
// running game with no lock held. The tests below assert the two properties
// that make that impossible: the returned view is unaffected by later
// mutations of the state it was built from, and concurrent readers of returned
// views do not race concurrent writers (checked under -race).

// viewTestGame starts a two-player game and returns the engine plus its live
// state, reaching past the public API on purpose -- the point of these tests is
// the relationship between the view and the authoritative state.
func viewTestGame(t *testing.T) (*GameEngine, *GameState, string) {
	t.Helper()

	e := NewGameEngine(zap.NewNop())
	const gameID = "view-aliasing-test"
	if err := e.StartGame(gameID, []string{"p1", "p2"}, "Commander Free For All"); err != nil {
		t.Fatalf("StartGame: %v", err)
	}
	state := e.games[gameID]

	// Put something in every zone the view exposes, so each copy path is
	// actually exercised rather than copying an empty slice.
	card := func(id, name string) *Card {
		return &Card{
			ID:         id,
			Name:       name,
			Counters:   []Counter{{Name: "+1/+1", Count: 1}},
			AttachedTo: []string{"someone"},
		}
	}
	state.Battlefield = append(state.Battlefield, card("bf-1", "Grizzly Bears"))
	state.Exile = append(state.Exile, card("ex-1", "Exiled Thing"))
	state.Stack = append(state.Stack, card("st-1", "Stacked Thing"))
	state.Command = append(state.Command, card("cm-1", "Marath"))
	state.Log = append(state.Log, LogEntry{Kind: "test", Message: "hello"})

	p1 := state.Players["p1"]
	p1.Hand = append(p1.Hand, card("hand-1", "Island"))
	p1.Library = append(p1.Library, card("lib-1", "Forest"))
	p1.Graveyard = append(p1.Graveyard, card("gy-1", "Shock"))
	p1.ManaPool.Red = 2

	p2 := state.Players["p2"]
	p2.Graveyard = append(p2.Graveyard, card("gy-2", "Bolt"))
	p2.Library = append(p2.Library, card("lib-2", "Mountain"))
	p2.RevealedTopCard = true

	return e, state, gameID
}

// TestGetGameViewDoesNotAliasState is the direct regression test: take a view,
// mutate everything reachable from the state it was built from, and assert the
// view did not move.
func TestGetGameViewDoesNotAliasState(t *testing.T) {
	e, state, gameID := viewTestGame(t)

	raw, err := e.GetGameView(gameID, "p1")
	if err != nil {
		t.Fatalf("GetGameView: %v", err)
	}
	v, ok := raw.(*PlaytestGameView)
	if !ok {
		t.Fatalf("GetGameView returned %T, want *PlaytestGameView", raw)
	}

	// Snapshot the things we are about to assert on.
	bfLen := len(v.Battlefield)
	exLen := len(v.Exile)
	stLen := len(v.Stack)
	cmLen := len(v.Command)
	logLen := len(v.Log)
	handLen := len(v.Me.Hand)
	libLen := len(v.Me.Library)
	gyLen := len(v.Me.Graveyard)
	life := v.Me.Life
	red := v.Me.ManaPool.Red
	tapped := v.Battlefield[0].Tapped
	counterCount := v.Battlefield[0].Counters[0].Count
	attached := v.Battlefield[0].AttachedTo[0]
	oppGYLen := len(v.Opponents[0].Graveyard)
	topName := v.Opponents[0].TopCard.Name

	// Now mutate the live state in every way the view could have aliased.
	state.Battlefield = append(state.Battlefield, &Card{ID: "bf-2", Name: "Added Later"})
	state.Exile = append(state.Exile, &Card{ID: "ex-2"})
	state.Stack = append(state.Stack, &Card{ID: "st-2"})
	state.Command = append(state.Command, &Card{ID: "cm-2"})
	state.Log = append(state.Log, LogEntry{Kind: "test", Message: "later"})

	state.Battlefield[0].Tapped = true
	state.Battlefield[0].Name = "Mutated"
	state.Battlefield[0].Counters[0].Count = 99
	state.Battlefield[0].AttachedTo[0] = "someone-else"

	p1 := state.Players["p1"]
	p1.Life = 3
	p1.ManaPool.Red = 99
	p1.Hand = append(p1.Hand, &Card{ID: "hand-2"})
	p1.Library = append(p1.Library, &Card{ID: "lib-3"})
	p1.Graveyard = append(p1.Graveyard, &Card{ID: "gy-3"})

	p2 := state.Players["p2"]
	p2.Graveyard = append(p2.Graveyard, &Card{ID: "gy-4"})
	p2.Library[0].Name = "Mutated Top"

	// Nothing the view holds may have moved.
	eq := func(what string, got, want interface{}) {
		t.Helper()
		if got != want {
			t.Errorf("%s: view saw %v after mutation, want %v (view aliases live state)", what, got, want)
		}
	}
	eq("len(Battlefield)", len(v.Battlefield), bfLen)
	eq("len(Exile)", len(v.Exile), exLen)
	eq("len(Stack)", len(v.Stack), stLen)
	eq("len(Command)", len(v.Command), cmLen)
	eq("len(Log)", len(v.Log), logLen)
	eq("len(Me.Hand)", len(v.Me.Hand), handLen)
	eq("len(Me.Library)", len(v.Me.Library), libLen)
	eq("len(Me.Graveyard)", len(v.Me.Graveyard), gyLen)
	eq("Me.Life", v.Me.Life, life)
	eq("Me.ManaPool.Red", v.Me.ManaPool.Red, red)
	eq("Battlefield[0].Tapped", v.Battlefield[0].Tapped, tapped)
	eq("Battlefield[0].Name", v.Battlefield[0].Name, "Grizzly Bears")
	eq("Battlefield[0].Counters[0].Count", v.Battlefield[0].Counters[0].Count, counterCount)
	eq("Battlefield[0].AttachedTo[0]", v.Battlefield[0].AttachedTo[0], attached)
	eq("len(Opponents[0].Graveyard)", len(v.Opponents[0].Graveyard), oppGYLen)
	eq("Opponents[0].TopCard.Name", v.Opponents[0].TopCard.Name, topName)

	// And the reverse direction: writing through the view must not reach the
	// engine. A consumer holding an aliased view could corrupt a running game.
	v.Battlefield[0].Name = "Written Through View"
	v.Me.ManaPool.Red = 42
	if state.Battlefield[0].Name != "Mutated" {
		t.Errorf("writing through view changed live battlefield card name to %q", state.Battlefield[0].Name)
	}
	if p1.ManaPool.Red != 99 {
		t.Errorf("writing through view changed live mana pool to %d", p1.ManaPool.Red)
	}
}

// TestGetGameViewConcurrentWithMutation is the -race half. Without the deep
// copy this reports a write/read pair between Mulligan (or any actions.go
// mutation) and a reader walking a previously-returned view.
func TestGetGameViewConcurrentWithMutation(t *testing.T) {
	e, _, gameID := viewTestGame(t)

	const iterations = 200
	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Writer: real engine mutations under the write lock.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = e.Mulligan(gameID, "p1")
			_ = e.ModifyLife(gameID, "p2", i%3-1)
			_ = e.AddCounter(gameID, "p1", "bf-1", "+1/+1", 1)
			select {
			case <-stop:
				return
			default:
			}
		}
	}()

	// Readers: fetch views and keep reading them after GetGameView returned,
	// which is exactly the window the old code left unsynchronized.
	for r := 0; r < 3; r++ {
		wg.Add(1)
		go func(seat string) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				raw, err := e.GetGameView(gameID, seat)
				if err != nil {
					continue
				}
				v := raw.(*PlaytestGameView)
				sink := 0
				for _, c := range v.Battlefield {
					sink += len(c.Name) + len(c.Counters) + len(c.AttachedTo)
				}
				if v.Me != nil {
					sink += v.Me.Life + len(v.Me.Hand) + len(v.Me.Library) + v.Me.ManaPool.Red
				}
				for _, o := range v.Opponents {
					sink += o.Life + len(o.Graveyard)
				}
				sink += len(v.Log) + len(v.Exile) + len(v.Stack) + len(v.Command)
				_ = sink
			}
		}([]string{"p1", "p2"}[r%2])
	}

	wg.Wait()
	close(stop)
}

// benchGameState builds a 4-player state of roughly mid-game size: a full
// commander library behind each seat, a hand, a board, a graveyard, and a log.
// This is the shape broadcast fans out over on every single mutation.
func benchGameState(b *testing.B) (*GameEngine, string) {
	b.Helper()

	e := NewGameEngine(zap.NewNop())
	const gameID = "bench"
	seats := []string{"p1", "p2", "p3", "p4"}
	if err := e.StartGame(gameID, seats, "Commander Free For All"); err != nil {
		b.Fatalf("StartGame: %v", err)
	}
	state := e.games[gameID]

	mk := func(id string) *Card {
		return &Card{
			ID: id, Name: "Card " + id, DisplayName: "Card " + id,
			ManaCost: "{2}{G}", Type: "Creature — Beast", RulesText: "Trample.",
			Counters: []Counter{{Name: "+1/+1", Count: 2}},
		}
	}
	for _, s := range seats {
		p := state.Players[s]
		for i := 0; i < 90; i++ {
			p.Library = append(p.Library, mk(fmt.Sprintf("%s-lib-%d", s, i)))
		}
		p.LibraryCount = len(p.Library)
		for i := 0; i < 7; i++ {
			p.Hand = append(p.Hand, mk(fmt.Sprintf("%s-hand-%d", s, i)))
		}
		p.HandCount = len(p.Hand)
		for i := 0; i < 10; i++ {
			p.Graveyard = append(p.Graveyard, mk(fmt.Sprintf("%s-gy-%d", s, i)))
		}
		for i := 0; i < 8; i++ {
			state.Battlefield = append(state.Battlefield, mk(fmt.Sprintf("%s-bf-%d", s, i)))
		}
	}
	for i := 0; i < 200; i++ {
		state.Log = append(state.Log, LogEntry{Kind: "play", Message: "something happened"})
	}
	return e, gameID
}

// BenchmarkBuildGameView measures one player's view.
func BenchmarkBuildGameView(b *testing.B) {
	e, gameID := benchGameState(b)
	state := e.games[gameID]
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.buildGameView(state, "p1")
	}
}

// BenchmarkBroadcastFanout measures what broadcast actually costs: one view per
// player, per mutation, on a 4-player game.
func BenchmarkBroadcastFanout(b *testing.B) {
	e, gameID := benchGameState(b)
	state := e.games[gameID]
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for playerID := range state.Players {
			_ = e.buildGameView(state, playerID)
		}
	}
}
