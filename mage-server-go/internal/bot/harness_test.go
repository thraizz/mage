package bot

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/magefree/mage-server-go/internal/game"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// harness_test.go stands up a complete 4-player Commander game with nothing but
// the engine: no database, no session manager, no table manager, no gRPC, no
// websocket.
//
// That is not just convenience. MatchStart on the real path refuses to start
// without at least two players AND a live websocket session per seat
// (internal/server/grpc_game.go:93-110), so a bot simulation that went through
// it would need a fake transport for every seat. game.Manager.CreateGame needs
// only a logger, and EngineAdapter.StartGameWithDecks needs only the deck map,
// so the whole thing is four objects and a goroutine.
//
// KNOWN INHERITED DEFECTS. These are asserted, not fixed. They belong to a
// Commander-rules task; fixing them here would hide which layer broke when a
// bot game goes wrong:
//
//	20 starting life, not 40   internal/game/game_state.go:149
//	commanders shuffled into the library, not put in the command zone
//	                           internal/game/game_engine.go:94 (MainDeck + Commanders)
//	libraries are never shuffled at game start
//	cards are name-only -- ManaCost/Type/Power/Toughness all blank
//	                           internal/game/game_engine.go:94-114
//
// TestInheritedCommanderDefects below pins every one of them, so the day
// somebody fixes them a test tells them which bot assumptions to revisit.

// ---------------------------------------------------------------------------
// Fixture card pool
// ---------------------------------------------------------------------------

// simOracle is the in-package decklist oracle. It exists so that no test in
// this package needs a live Postgres: ScryfallOracle satisfies the same
// OracleLookup interface in production, and swapping it in is a one-line
// change in RunnerConfig.
func simOracle() MapOracle {
	return MapOracle{
		"Forest":   {Name: "Forest", TypeLine: "Basic Land — Forest", OracleText: "({T}: Add {G}.)"},
		"Island":   {Name: "Island", TypeLine: "Basic Land — Island", OracleText: "({T}: Add {U}.)"},
		"Mountain": {Name: "Mountain", TypeLine: "Basic Land — Mountain", OracleText: "({T}: Add {R}.)"},
		"Plains":   {Name: "Plains", TypeLine: "Basic Land — Plains", OracleText: "({T}: Add {W}.)"},
		"Swamp":    {Name: "Swamp", TypeLine: "Basic Land — Swamp", OracleText: "({T}: Add {B}.)"},
		"Command Tower": {
			Name: "Command Tower", TypeLine: "Land",
			OracleText: "{T}: Add one mana of any color in your commander's color identity.",
		},
		"Llanowar Elves": {
			Name: "Llanowar Elves", ManaCost: "{G}", TypeLine: "Creature — Elf Druid",
			OracleText: "{T}: Add {G}.", Power: "1", Toughness: "1",
		},
		"Grizzly Bears": {
			Name: "Grizzly Bears", ManaCost: "{1}{G}", TypeLine: "Creature — Bear",
			Power: "2", Toughness: "2",
		},
		"Hill Giant": {
			Name: "Hill Giant", ManaCost: "{3}{R}", TypeLine: "Creature — Giant",
			Power: "3", Toughness: "3",
		},
		"Serra Angel": {
			Name: "Serra Angel", ManaCost: "{3}{W}{W}", TypeLine: "Creature — Angel",
			OracleText: "Flying, vigilance", Power: "4", Toughness: "4",
		},
		"Air Elemental": {
			Name: "Air Elemental", ManaCost: "{3}{U}{U}", TypeLine: "Creature — Elemental",
			OracleText: "Flying", Power: "4", Toughness: "4",
		},
		"Bog Wraith": {
			Name: "Bog Wraith", ManaCost: "{3}{B}", TypeLine: "Creature — Wraith",
			OracleText: "Swampwalk", Power: "3", Toughness: "3",
		},
		"Lightning Bolt": {
			Name: "Lightning Bolt", ManaCost: "{R}", TypeLine: "Instant",
			OracleText: "Lightning Bolt deals 3 damage to any target.",
		},
		"Giant Growth": {
			Name: "Giant Growth", ManaCost: "{G}", TypeLine: "Instant",
			OracleText: "Target creature gets +3/+3 until end of turn.",
		},
		"Divination": {
			Name: "Divination", ManaCost: "{2}{U}", TypeLine: "Sorcery",
			OracleText: "Draw two cards.",
		},
		"Marath, Will of the Wild": {
			Name: "Marath, Will of the Wild", ManaCost: "{R}{G}{W}",
			TypeLine: "Legendary Creature — Elemental Beast",
			Power:    "0", Toughness: "0",
		},
	}
}

// simDeck builds a 99-card main deck plus one commander, deterministically.
//
// The composition is 38 lands / 61 spells -- a normal Commander curve -- drawn
// from the fixture pool above. It is intentionally NOT randomized: the engine
// never shuffles (see the defect list), so deck order is the only thing
// deciding what a bot draws, and a fixed order keeps a seeded simulation
// reproducible end to end.
func simDeck(seatIndex int) game.DeckList {
	lands := []string{"Forest", "Island", "Mountain", "Plains", "Swamp", "Command Tower"}
	spells := []string{
		"Llanowar Elves", "Grizzly Bears", "Hill Giant", "Serra Angel",
		"Air Elemental", "Bog Wraith", "Lightning Bolt", "Giant Growth", "Divination",
	}

	main := make([]string, 0, 99)
	// Rotate the pool per seat so the four decks are not identical, without
	// introducing randomness.
	for i := 0; i < 99; i++ {
		if i%5 < 2 {
			main = append(main, lands[(i+seatIndex)%len(lands)])
		} else {
			main = append(main, spells[(i+seatIndex)%len(spells)])
		}
	}
	return game.DeckList{
		MainDeck:   main,
		Commanders: []string{"Marath, Will of the Wild"},
	}
}

// ---------------------------------------------------------------------------
// Headless game setup
// ---------------------------------------------------------------------------

type simGame struct {
	mgr     *game.Manager
	adapter *game.EngineAdapter
	engine  *game.GameEngine
	g       *game.Game
	seats   []string
}

func (s *simGame) close() { s.mgr.RemoveGame(s.g.ID) } // closes ActionQueue, ending ProcessGameActions

func newSimGame(t testing.TB, logger *zap.Logger, tableID string, seats []string) *simGame {
	t.Helper()

	engine := game.NewGameEngine(logger)
	adapter := game.NewEngineAdapter(engine, logger)
	mgr := game.NewManager(logger)

	g := mgr.CreateGame(tableID, "Commander Free For All", seats)

	decks := make(map[string]game.DeckList, len(seats))
	for i, s := range seats {
		decks[s] = simDeck(i)
	}
	if err := adapter.StartGameWithDecks(g, decks); err != nil {
		t.Fatalf("StartGameWithDecks: %v", err)
	}
	go adapter.ProcessGameActions(g)

	return &simGame{mgr: mgr, adapter: adapter, engine: engine, g: g, seats: seats}
}

func fourSeats() []string { return []string{"bot-a", "bot-b", "bot-c", "bot-d"} }

// runOneGame plays a single headless 4-bot game and reports whether it reached
// a terminal state.
func runOneGame(t testing.TB, logger *zap.Logger, seed int64, maxTurns int) (bool, Stats, error) {
	t.Helper()

	seats := fourSeats()
	sg := newSimGame(t, logger, fmt.Sprintf("table-%d", seed), seats)
	defer sg.close()

	runner := NewBotRunner(RunnerConfig{
		GameID:  sg.g.ID,
		Actions: sg.mgr,
		Views:   sg.adapter,
		Oracle:  simOracle(),
		// Pacing is entirely zero: a hundred games with human-like delays would
		// take hours. Phase 4 fills these in.
		Pacing:   Pacing{},
		MaxTurns: maxTurns,
		Logger:   logger,
	})
	for i, s := range seats {
		runner.AddSeat(s, NewRandomPolicy(seed*1000+int64(i)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	done, err := runner.Run(ctx)
	return done, runner.Stats(), err
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestHeadlessCommanderGameCompletes(t *testing.T) {
	logger := zaptest.NewLogger(t)
	done, stats, err := runOneGame(t, logger, 1, 200)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	t.Logf("completed=%v turns=%d macros=%d failed=%d steps=%d",
		done, stats.Turns, stats.MacrosExecuted, stats.MacrosFailed, stats.StepsSent)
	if !done {
		t.Fatalf("game did not reach a terminal state within 200 turns")
	}
	if stats.MacrosFailed != 0 {
		t.Fatalf("%d macro(s) failed to land: %v", stats.MacrosFailed, stats.FailedMacros)
	}
	if stats.MacrosExecuted == 0 {
		t.Fatal("no macros executed at all")
	}
}

// TestBotSimCompletionRate is the Phase 3 milestone: random bots must be able to
// finish a game. Set BOT_SIM_GAMES to change the sample size; `make bot-sim`
// sets it to 100.
func TestBotSimCompletionRate(t *testing.T) {
	n := 20
	if v := os.Getenv("BOT_SIM_GAMES"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			t.Fatalf("BOT_SIM_GAMES=%q: %v", v, err)
		}
		n = parsed
	}
	if testing.Short() {
		n = 3
	}

	// zap.NewNop, not zaptest: a hundred games produce a lot of log lines and
	// zaptest buffers all of them into the test output.
	logger := zap.NewNop()

	completed, failedMacros := 0, 0
	var reasons []string
	start := time.Now()
	for i := 0; i < n; i++ {
		done, stats, err := runOneGame(t, logger, int64(i+1), 200)
		if err != nil {
			reasons = append(reasons, fmt.Sprintf("game %d: %v", i, err))
			continue
		}
		failedMacros += stats.MacrosFailed
		if done {
			completed++
		} else {
			reasons = append(reasons, fmt.Sprintf("game %d: turn cap at turn %d", i, stats.Turns))
		}
		if len(stats.FailedMacros) > 0 {
			reasons = append(reasons, fmt.Sprintf("game %d: failed macros %v", i, stats.FailedMacros))
		}
	}

	rate := float64(completed) / float64(n) * 100
	t.Logf("completion rate: %d/%d = %.1f%% in %s (failed macros: %d)",
		completed, n, rate, time.Since(start).Round(time.Millisecond), failedMacros)
	for _, r := range reasons {
		t.Logf("  %s", r)
	}
	if rate < 95 {
		t.Fatalf("completion rate %.1f%% is below the 95%% Phase 3 gate", rate)
	}
	if failedMacros != 0 {
		t.Fatalf("%d macro(s) failed to land across %d games", failedMacros, n)
	}
}

// TestEveryLegalMoveExecutes is verification item 8: every macro LegalMoves
// offers must actually run against the engine.
//
// It runs one seat alone so the game log is an exact per-step counter -- the
// only read-back signal available, because ProcessGameActions swallows errors
// (anti-pattern 8). A step that names an unimplemented verb, or a card the
// engine cannot find, appends nothing and is caught here.
func TestEveryLegalMoveExecutes(t *testing.T) {
	logger := zaptest.NewLogger(t)
	seats := fourSeats()
	sg := newSimGame(t, logger, "table-macro-exec", seats)
	defer sg.close()

	me := seats[0]
	oracle := simOracle()

	viewOf := func() *SafeView {
		t.Helper()
		raw, err := sg.adapter.GetGameView(sg.g.ID, me)
		if err != nil {
			t.Fatalf("view: %v", err)
		}
		v, err := RedactErr(raw.(*game.PlaytestGameView), me)
		if err != nil {
			t.Fatalf("view: %v", err)
		}
		Enrich(context.Background(), v, oracle)
		return v
	}

	// Send one step and assert the log grew by exactly one -- i.e. the engine
	// really performed it.
	sendStep := func(step string) {
		t.Helper()
		before := len(viewOf().Log)
		if err := sg.mgr.SendPlayerAction(sg.g.ID, me, "SEND_STRING", step); err != nil {
			t.Fatalf("enqueue %q: %v", step, err)
		}
		deadline := time.Now().Add(3 * time.Second)
		for {
			if len(viewOf().Log) == before+1 {
				return
			}
			if time.Now().After(deadline) {
				t.Fatalf("step %q produced no engine state change (it was silently swallowed)", step)
			}
			time.Sleep(time.Millisecond)
		}
	}

	// Phase 1: the mulligan decision set.
	pre := viewOf()
	if pre.Me.KeptHand {
		t.Fatal("expected the seat to start on a mulligan decision")
	}
	mullMoves := LegalMoves(pre)
	if len(mullMoves) != 2 {
		t.Fatalf("mulligan decision offered %d macros, want 2 (keep, mulligan)", len(mullMoves))
	}
	for _, m := range mullMoves {
		if m.KindOf() == KindKeepHand {
			continue
		}
		for _, step := range m.Steps {
			sendStep(step)
		}
	}
	// Now keep, so the rest of the move set unlocks.
	sendStep("KEEP_HAND:" + me)

	// Phase 2: every macro of the full move set, one at a time, each against a
	// freshly read view (a macro can invalidate the ones after it -- moving a
	// card out of hand makes the "Discard <that card>" macro stale).
	seenKinds := map[Kind]int{}
	for round := 0; round < 40; round++ {
		v := viewOf()
		if v.ActivePlayerID != me {
			t.Fatalf("seat lost priority unexpectedly")
		}
		moves := LegalMoves(v)
		if len(moves) == 0 {
			t.Fatal("no legal moves offered")
		}
		// Prefer a macro class not yet exercised, so the loop covers the whole
		// vocabulary rather than whatever the round-robin index lands on.
		m := moves[round%len(moves)]
		for _, cand := range moves {
			if cand.KindOf() != KindPassTurn && seenKinds[cand.KindOf()] == 0 {
				m = cand
				break
			}
		}
		if m.KindOf() == KindPassTurn {
			// Passing hands the turn to another seat, ending this test's
			// exclusive access. Verify it last.
			continue
		}
		seenKinds[m.KindOf()]++
		for _, step := range m.Steps {
			sendStep(step)
		}
	}

	// Every macro class the fixture can produce must have been exercised.
	for _, want := range []Kind{
		KindPlayLand, KindCast, KindTap, KindUntap,
		KindMoveZone, KindModifyLife, KindDraw,
	} {
		if seenKinds[want] == 0 {
			t.Errorf("macro kind %q was never exercised", want)
		}
	}
	t.Logf("exercised macro kinds: %v", seenKinds)

	// Finally, pass the turn.
	sendStep("NEXT_TURN:" + me)
	if got := viewOf().ActivePlayerID; got == me {
		t.Fatalf("NEXT_TURN did not advance the active player (still %q)", got)
	}
}

// TestInheritedCommanderDefects pins the known-broken Commander behaviour so
// that fixing it is a deliberate, visible change rather than a silent one that
// quietly alters every bot simulation. Phase 3 asserts around these; it does
// not fix them.
func TestInheritedCommanderDefects(t *testing.T) {
	logger := zaptest.NewLogger(t)
	seats := fourSeats()
	sg := newSimGame(t, logger, "table-defects", seats)
	defer sg.close()

	// No bots are running in this test, so a plain read is safe.
	raw, err := sg.adapter.GetGameView(sg.g.ID, seats[0])
	if err != nil {
		t.Fatalf("GetGameView: %v", err)
	}
	pv := raw.(*game.PlaytestGameView)

	// DEFECT 1: starting life is hardcoded to 20 (internal/game/game_state.go:149).
	// Commander is 40. internal/plugin, which would know that, is referenced
	// only by a blank import (anti-pattern 6).
	if pv.Me.Life != 20 {
		t.Errorf("starting life = %d; the 20-life defect appears to be fixed -- "+
			"update the bot harness and the Commander-rules task", pv.Me.Life)
	}
	for _, o := range pv.Opponents {
		if o.Life != 20 {
			t.Errorf("opponent %s starting life = %d, want the inherited 20", o.PlayerID, o.Life)
		}
	}

	// DEFECT 2: commanders are appended to the main deck and shuffled into the
	// library instead of starting in the command zone
	// (internal/game/game_engine.go:94).
	if len(pv.Command) != 0 {
		t.Errorf("command zone has %d card(s); the commander-zone defect appears "+
			"to be fixed -- revisit the harness", len(pv.Command))
	}
	commanderInLibrary := false
	for _, c := range pv.Me.Library {
		if c.Name == "Marath, Will of the Wild" {
			commanderInLibrary = true
			break
		}
	}
	if !commanderInLibrary && !strings.Contains(handNames(pv), "Marath") {
		t.Errorf("commander is in neither library nor hand; expected the inherited " +
			"behaviour of shuffling it into the deck")
	}

	// DEFECT 3: the library is never shuffled at game start -- it is exactly
	// the decklist order. That is why simDeck is deterministic: deck order is
	// the only thing deciding draws.
	want := simDeck(0).MainDeck
	got := make([]string, 0, len(pv.Me.Hand)+len(pv.Me.Library))
	for _, c := range pv.Me.Hand {
		got = append(got, c.Name)
	}
	for _, c := range pv.Me.Library {
		got = append(got, c.Name)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("card %d is %q, want %q -- libraries appear to be shuffled now; "+
				"the harness assumed decklist order", i, got[i], want[i])
			break
		}
	}

	// DEFECT 4: cards are name-only. ManaCost/Type/Power/Toughness are blank,
	// which is exactly why bot.Enrich exists.
	for _, c := range pv.Me.Hand {
		if c.ManaCost != "" || c.Type != "" || c.Power != "" {
			t.Errorf("card %q arrived with printed characteristics (cost=%q type=%q power=%q); "+
				"the name-only defect appears fixed -- Enrich can be simplified",
				c.Name, c.ManaCost, c.Type, c.Power)
			break
		}
	}
}

func handNames(pv *game.PlaytestGameView) string {
	var b strings.Builder
	for _, c := range pv.Me.Hand {
		b.WriteString(c.Name)
		b.WriteByte(' ')
	}
	return b.String()
}
