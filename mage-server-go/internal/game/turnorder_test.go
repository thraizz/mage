package game

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
)

// turnorder_test.go guards the determinism fix described on GameState.TurnOrder.
//
// Before that fix NextTurn derived the seating order by ranging over
// state.Players, a map, so Go's randomized map iteration gave a different turn
// order on every call. Nothing in the engine noticed, because nothing in the
// engine depended on sequencing -- but a bot simulation does, and an
// irreproducible bot game cannot be debugged.

// fullTurnCycle plays len(players) NextTurn calls in a fresh game and returns
// the sequence of active players it observed.
func fullTurnCycle(t *testing.T, players []string) []string {
	t.Helper()

	e := NewGameEngine(zap.NewNop())
	gameID := "turn-order-test"
	if err := e.StartGame(gameID, players, "Commander Free For All"); err != nil {
		t.Fatalf("StartGame: %v", err)
	}

	seen := make([]string, 0, len(players)*2)
	seen = append(seen, e.games[gameID].ActivePlayerID)
	// Two full cycles, so a bug that only shows up on wrap-around is caught too.
	for i := 0; i < len(players)*2; i++ {
		state := e.games[gameID]
		if err := e.NextTurn(gameID, state.ActivePlayerID); err != nil {
			t.Fatalf("NextTurn: %v", err)
		}
		seen = append(seen, e.games[gameID].ActivePlayerID)
	}
	return seen
}

func TestNextTurnOrderIsStableAcrossRuns(t *testing.T) {
	players := []string{"delta", "alpha", "charlie", "bravo"}

	want := fullTurnCycle(t, players)

	// 100 independent games. Go re-randomizes map iteration per range statement,
	// per process, so the pre-fix code failed this within a handful of runs.
	for run := 0; run < 100; run++ {
		got := fullTurnCycle(t, players)
		if len(got) != len(want) {
			t.Fatalf("run %d: got %d turns, want %d", run, len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("run %d: turn %d active player = %q, want %q\n got: %v\nwant: %v",
					run, i, got[i], want[i], got, want)
			}
		}
	}
}

func TestNextTurnFollowsSeatingOrderNotAlphabetical(t *testing.T) {
	// Deliberately not in alphabetical order: turn order is seating order.
	players := []string{"delta", "alpha", "charlie", "bravo"}

	got := fullTurnCycle(t, players)
	want := []string{
		"delta", "alpha", "charlie", "bravo",
		"delta", "alpha", "charlie", "bravo",
		"delta",
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("turn %d active player = %q, want %q (got %v)", i, got[i], want[i], got)
		}
	}
}

func TestNextTurnStableWithoutExplicitTurnOrder(t *testing.T) {
	// A GameState that predates the TurnOrder field -- e.g. restored from an
	// older snapshot -- must still cycle in a stable order.
	e := NewGameEngine(zap.NewNop())
	gameID := "legacy-state"

	build := func() *GameState {
		st := NewGameState(gameID, []string{"p1", "p2", "p3"}, nil)
		st.TurnOrder = nil // simulate the pre-fix serialized shape
		return st
	}

	var want []string
	for run := 0; run < 100; run++ {
		e.games[gameID] = build()
		got := make([]string, 0, 6)
		for i := 0; i < 6; i++ {
			if err := e.NextTurn(gameID, e.games[gameID].ActivePlayerID); err != nil {
				t.Fatalf("NextTurn: %v", err)
			}
			got = append(got, e.games[gameID].ActivePlayerID)
		}
		if want == nil {
			want = got
			continue
		}
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Fatalf("run %d: got %v, want %v", run, got, want)
		}
	}
}

func TestStartGameWithDecksRecordsTurnOrder(t *testing.T) {
	e := NewGameEngine(zap.NewNop())
	players := []string{"zeta", "alpha", "mike"}
	decks := map[string]DeckList{
		"zeta":  {MainDeck: []string{"Forest"}},
		"alpha": {MainDeck: []string{"Island"}},
		"mike":  {MainDeck: []string{"Swamp"}},
	}
	if err := e.StartGameWithDecks("g", players, "Commander Free For All", decks); err != nil {
		t.Fatalf("StartGameWithDecks: %v", err)
	}
	got := e.games["g"].TurnOrder
	for i, want := range players {
		if got[i] != want {
			t.Fatalf("TurnOrder[%d] = %q, want %q (got %v)", i, got[i], want, got)
		}
	}
}
