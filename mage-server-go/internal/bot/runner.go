package bot

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/magefree/mage-server-go/internal/game"
	"go.uber.org/zap"
)

// runner.go drives bot seats.
//
// INTEGRATION IS PULL-BASED, ON PURPOSE. The engine has a push path --
// GameEngine.SetNotificationHandler / EngineAdapter.SetNotificationCallback --
// and this package deliberately does not use it (anti-pattern 2 and 5):
//
//   - GameEngine.broadcast (internal/game/game_engine.go:342) runs while e.mu
//     is held for WRITING; every mutation in actions.go is
//     Lock(); defer Unlock(); ...; broadcast(). A handler that synchronously
//     calls GetGameView or SendPlayerAction from inside broadcast deadlocks on
//     the engine's own mutex. The codebase already documents the hazard at
//     internal/server/grpc.go:188-190.
//   - notifyFn is a single field and SetNotificationCallback overwrites it. The
//     websocket layer already owns that one registration; a bot that grabbed it
//     would silently unsubscribe every human client.
//
// Polling EngineAdapter.GetGameView (internal/game/manager.go:482) takes the
// read lock, needs zero engine changes, and cannot deadlock. Phase 7 can
// revisit if latency ever matters; for a headless simulation it never will.

// Policy chooses one macro from the set offered.
//
// Phase 3 ships RandomPolicy. Phase 5's LLMPolicy implements the same
// interface, which is the whole point of the split: the loop is proven with a
// policy that cannot possibly be the reason a game fails to finish.
type Policy interface {
	Pick(ctx context.Context, v *SafeView, moves []Macro) (Macro, error)
}

// ErrNoMoves is returned by a Policy handed an empty move set.
var ErrNoMoves = errors.New("bot: no moves available")

// RandomPolicy picks uniformly at random from the offered macros.
//
// The RNG is injectable and seeded rather than package-global. A bot simulation
// that cannot be replayed is not a debugging tool, and math/rand's global
// source is both unseeded-per-run and shared with the engine's own shuffles.
type RandomPolicy struct {
	mu  sync.Mutex
	rng *rand.Rand
}

// NewRandomPolicy returns a RandomPolicy over a fresh source seeded with seed.
func NewRandomPolicy(seed int64) *RandomPolicy {
	return &RandomPolicy{rng: rand.New(rand.NewSource(seed))}
}

// NewRandomPolicyWithRand returns a RandomPolicy over a caller-supplied source.
func NewRandomPolicyWithRand(r *rand.Rand) *RandomPolicy {
	return &RandomPolicy{rng: r}
}

// Pick implements Policy.
func (p *RandomPolicy) Pick(_ context.Context, _ *SafeView, moves []Macro) (Macro, error) {
	if len(moves) == 0 {
		return Macro{}, ErrNoMoves
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return moves[p.rng.Intn(len(moves))], nil
}

// ActionSender is the subset of *game.Manager the runner needs.
type ActionSender interface {
	SendPlayerAction(gameID, playerID, actionType string, data interface{}) error
}

// ViewSource is the subset of *game.EngineAdapter the runner needs.
type ViewSource interface {
	GetGameView(gameID, playerID string) (interface{}, error)
}

// Pacing controls how fast a bot acts.
//
// Every field defaults to zero here, and zero means "as fast as the scheduler
// allows". Phase 3's completion-rate test runs a hundred full games; human-like
// delays would turn that into hours. Phase 4 adds the log-normal pacer that
// makes bots look human in the real client, and it does it by filling these in
// -- the runner already puts the delays in the right places, between
// SendPlayerAction calls in the bot's own goroutine and never inside an engine
// callback.
type Pacing struct {
	// Poll is how long to wait between GetGameView calls when there is nothing
	// to do. Zero yields to the scheduler instead of sleeping.
	Poll time.Duration
	// PreAction is the "thinking" delay before a macro's first step.
	PreAction time.Duration
	// BetweenSteps staggers a macro's steps so a watching client animates each
	// one separately.
	BetweenSteps time.Duration
}

func (p Pacing) wait(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		runtime.Gosched()
		return ctx.Err() == nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

// RunnerConfig configures a BotRunner.
type RunnerConfig struct {
	GameID  string
	Actions ActionSender
	Views   ViewSource

	// Oracle supplies the printed characteristics the engine never sets. May be
	// nil, in which case only basic lands are recognised and nothing is
	// castable.
	Oracle OracleLookup

	Pacing Pacing

	// MaxActionsPerTurn bounds a seat's action loop before the turn is passed
	// regardless of what the policy wants. It is what guarantees the simulation
	// makes forward progress: without it a policy that never picks "Pass the
	// turn" hangs the game forever. Default 12.
	MaxActionsPerTurn int

	// MaxTurns stops the simulation if the game has not reached a terminal
	// state. Default 200. Hitting it is a non-completion, not a success.
	MaxTurns int

	// StepTimeout bounds the read-back wait after a macro's steps are sent.
	// Default 2s.
	StepTimeout time.Duration

	Logger *zap.Logger
}

func (c *RunnerConfig) applyDefaults() {
	if c.MaxActionsPerTurn <= 0 {
		c.MaxActionsPerTurn = 12
	}
	if c.MaxTurns <= 0 {
		c.MaxTurns = 200
	}
	if c.StepTimeout <= 0 {
		c.StepTimeout = 2 * time.Second
	}
	if c.Logger == nil {
		c.Logger = zap.NewNop()
	}
}

// Stats is what a simulation reports about one game.
type Stats struct {
	MacrosExecuted int
	MacrosFailed   int
	StepsSent      int
	Turns          int
	// FailedMacros records the labels of macros whose steps did not all land.
	// A macro appearing here is an engine-level failure, not a policy choice:
	// ProcessGameActions logs errors and never returns them (anti-pattern 8),
	// so the only way to see one is by reading state back.
	FailedMacros []string
}

func (s *Stats) merge(o Stats) {
	s.MacrosExecuted += o.MacrosExecuted
	s.MacrosFailed += o.MacrosFailed
	s.StepsSent += o.StepsSent
	if o.Turns > s.Turns {
		s.Turns = o.Turns
	}
	s.FailedMacros = append(s.FailedMacros, o.FailedMacros...)
}

// BotRunner runs one goroutine per bot seat.
type BotRunner struct {
	cfg   RunnerConfig
	seats []*seat

	mu    sync.Mutex
	stats Stats
}

// NewBotRunner creates a runner. Seats are added with AddSeat.
func NewBotRunner(cfg RunnerConfig) *BotRunner {
	cfg.applyDefaults()
	return &BotRunner{cfg: cfg}
}

// AddSeat registers a bot seat and the policy that drives it.
func (r *BotRunner) AddSeat(botID string, p Policy) {
	r.seats = append(r.seats, &seat{r: r, botID: botID, policy: p})
}

// Stats returns the aggregate statistics of the run.
func (r *BotRunner) Stats() Stats {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := r.stats
	out.FailedMacros = append([]string(nil), r.stats.FailedMacros...)
	return out
}

// Run starts every seat and blocks until the game reaches a terminal state, the
// turn cap is hit, or ctx is cancelled.
//
// It returns whether the game completed, i.e. reached a state with at most one
// living seat.
func (r *BotRunner) Run(ctx context.Context) (completed bool, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for _, s := range r.seats {
		wg.Add(1)
		go func(s *seat) {
			defer wg.Done()
			defer cancel() // any seat seeing the end ends the game for everyone
			s.run(ctx)
		}(s)
	}
	wg.Wait()

	v, verr := r.view(r.seatID(0))
	if verr != nil {
		return false, verr
	}
	if v == nil {
		return false, errors.New("bot: no view available at end of run")
	}
	r.mu.Lock()
	r.stats.Turns = v.Turn
	r.mu.Unlock()
	return IsTerminal(v), nil
}

func (r *BotRunner) seatID(i int) string {
	if i >= len(r.seats) {
		return ""
	}
	return r.seats[i].botID
}

// view polls the engine and redacts the result.
//
// No external locking is needed. GameEngine.buildGameView deep-copies under the
// engine's read lock, so the *game.PlaytestGameView handed back here already
// owns its memory and Redact's own copy cannot race a concurrent mutation. It
// did not always -- see the Phase 7 note in docs/BOT_PLAYERS_PLAN.md for the
// bot-side workaround this replaced.
func (r *BotRunner) view(botID string) (*SafeView, error) {
	raw, err := r.cfg.Views.GetGameView(r.cfg.GameID, botID)
	if err != nil {
		return nil, err
	}
	pv, ok := raw.(*game.PlaytestGameView)
	if !ok || pv == nil {
		return nil, fmt.Errorf("bot: unexpected view type %T", raw)
	}
	sv, err := RedactErr(pv, botID)
	if err != nil {
		return nil, err
	}
	Enrich(context.Background(), sv, r.cfg.Oracle)
	return sv, nil
}

// IsTerminal reports whether the game is over: at most one seat still at a
// positive life total.
//
// The engine has no win condition of its own -- EndGame just writes a log line
// -- so "terminal" has to be defined by the caller. Life is the only resource
// the rules-light engine actually tracks per seat, so it is the one this uses.
func IsTerminal(v *SafeView) bool {
	return LivingSeats(v) <= 1
}

// LivingSeats counts seats at a positive life total, the viewer included.
func LivingSeats(v *SafeView) int {
	if v == nil || v.Me == nil {
		return 0
	}
	n := 0
	if v.Me.Life > 0 {
		n++
	}
	for _, o := range v.Opponents {
		if o != nil && o.Life > 0 {
			n++
		}
	}
	return n
}

// seat is one bot's goroutine state.
type seat struct {
	r      *BotRunner
	botID  string
	policy Policy

	stats Stats

	// lastUpkeepTurn is the turn number this seat last performed its untap /
	// draw step on, so the upkeep runs exactly once per turn.
	lastUpkeepTurn int
	// landDropTurn is the turn number this seat last played a land on.
	landDropTurn int
}

func (s *seat) run(ctx context.Context) {
	defer func() {
		s.r.mu.Lock()
		s.r.stats.merge(s.stats)
		s.r.mu.Unlock()
	}()

	for ctx.Err() == nil {
		v, err := s.r.view(s.botID)
		if err != nil {
			s.r.cfg.Logger.Error("bot: view failed",
				zap.String("bot", s.botID), zap.Error(err))
			return
		}
		if IsTerminal(v) || v.Turn > s.r.cfg.MaxTurns {
			return
		}

		switch {
		case !v.Me.KeptHand:
			// Mulligan decisions are made off-turn, by every seat at once.
			s.act(ctx, v)
		case v.ActivePlayerID == s.botID:
			s.takeTurn(ctx, v)
		default:
			if !s.r.cfg.Pacing.wait(ctx, s.r.cfg.Pacing.Poll) {
				return
			}
		}
	}
}

// takeTurn plays out this seat's turn and passes it.
func (s *seat) takeTurn(ctx context.Context, v *SafeView) {
	turn := v.Turn

	// A dead seat still has to pass, or the game stalls on its turn forever.
	if v.Me.Life <= 0 {
		s.exec(ctx, macro(KindPassTurn, "Pass the turn (eliminated)", "NEXT_TURN:"+s.botID))
		return
	}

	if s.lastUpkeepTurn != turn {
		s.lastUpkeepTurn = turn
		s.upkeep(ctx, v)
	}

	for i := 0; i < s.r.cfg.MaxActionsPerTurn && ctx.Err() == nil; i++ {
		cur, err := s.r.view(s.botID)
		if err != nil || cur == nil {
			return
		}
		if IsTerminal(cur) || cur.ActivePlayerID != s.botID || cur.Turn != turn {
			return
		}
		if !s.act(ctx, cur) {
			return
		}
	}

	// Action budget spent. Passing is not the policy's call any more: the cap
	// is what guarantees the simulation terminates.
	s.exec(ctx, macro(KindPassTurn, "Pass the turn (action budget spent)", "NEXT_TURN:"+s.botID))
}

// upkeep is turn structure, not a decision: untap everything, then draw.
//
// The engine exposes UNTAP_ALL only as a ProcessAction action type, not as a
// string command (§0.6), and the runner speaks exclusively in SEND_STRING
// steps, so this untaps card by card.
func (s *seat) upkeep(ctx context.Context, v *SafeView) {
	var steps []string
	for _, c := range controlledBy(v.Battlefield, s.botID) {
		if c.Tapped {
			steps = append(steps, "UNTAP:"+c.ID)
		}
	}
	if len(steps) > 0 {
		s.exec(ctx, macro(KindUntap, fmt.Sprintf("Untap %d permanent(s)", len(steps)), steps...))
	}
	if v.Me.LibraryCount > 0 {
		s.exec(ctx, macro(KindDraw, "Draw for turn", "DRAW:"+s.botID+":1"))
	}
}

// act offers the seat's moves to the policy and executes the pick. It reports
// whether the seat should keep acting this turn.
func (s *seat) act(ctx context.Context, v *SafeView) bool {
	moves := s.available(v)
	if len(moves) == 0 {
		return false
	}
	m, err := s.policy.Pick(ctx, v, moves)
	if err != nil {
		if !errors.Is(err, ErrNoMoves) {
			s.r.cfg.Logger.Error("bot: policy failed",
				zap.String("bot", s.botID), zap.Error(err))
		}
		return false
	}
	if !s.r.cfg.Pacing.wait(ctx, s.r.cfg.Pacing.PreAction) {
		return false
	}
	s.exec(ctx, m)
	if m.KindOf() == KindPlayLand {
		s.landDropTurn = v.Turn
	}
	return m.KindOf() != KindPassTurn
}

// available is LegalMoves plus the turn-structure filter the engine does not
// enforce: one land drop per turn.
func (s *seat) available(v *SafeView) []Macro {
	all := LegalMoves(v)
	if s.landDropTurn != v.Turn {
		return all
	}
	out := all[:0:0]
	for _, m := range all {
		if m.KindOf() == KindPlayLand {
			continue
		}
		out = append(out, m)
	}
	return out
}

// exec sends a macro's steps and then VERIFIES THEM BY READING STATE BACK.
//
// This is anti-pattern 8 in practice. SendPlayerAction only enqueues; the queue
// is drained by EngineAdapter.ProcessGameActions (internal/game/manager.go:439),
// which logs any error from ProcessAction and never returns it. A step naming a
// card that has moved zones, or a verb the engine does not implement, therefore
// "succeeds" from the sender's point of view and changes nothing.
//
// Every mutation in actions.go ends with appendLog followed by broadcast, and
// no failure path reaches appendLog, so the game log's length is an exact
// counter of steps that actually took effect. exec records it before sending
// and waits for it to advance by len(Steps).
func (s *seat) exec(ctx context.Context, m Macro) bool {
	if len(m.Steps) == 0 {
		return true
	}
	before, err := s.logLen()
	if err != nil {
		return false
	}

	for i, step := range m.Steps {
		if i > 0 && !s.r.cfg.Pacing.wait(ctx, s.r.cfg.Pacing.BetweenSteps) {
			return false
		}
		if err := s.r.cfg.Actions.SendPlayerAction(
			s.r.cfg.GameID, s.botID, "SEND_STRING", step); err != nil {
			s.stats.MacrosFailed++
			s.stats.FailedMacros = append(s.stats.FailedMacros,
				fmt.Sprintf("%s: enqueue %q: %v", m.Label, step, err))
			return false
		}
		s.stats.StepsSent++
	}

	landed, aborted := s.awaitLogGrowth(ctx, before+len(m.Steps))
	if aborted {
		// The run was cancelled out from under us -- another seat saw the game
		// end. That is a shutdown race, not an engine failure, and counting it
		// as one would make every completed game report a spurious defect.
		return false
	}
	if !landed {
		s.stats.MacrosFailed++
		s.stats.FailedMacros = append(s.stats.FailedMacros, m.Label)
		s.r.cfg.Logger.Warn("bot: macro did not fully land",
			zap.String("bot", s.botID),
			zap.String("macro", m.Label),
			zap.Strings("steps", m.Steps))
		return false
	}
	s.stats.MacrosExecuted++
	return true
}

func (s *seat) logLen() (int, error) {
	v, err := s.r.view(s.botID)
	if err != nil {
		return 0, err
	}
	return len(v.Log), nil
}

// awaitLogGrowth waits for the game log to reach want entries.
//
// Other seats can also append while this waits -- only during the mulligan
// phase, since afterwards a seat acts only on its own turn -- which makes the
// check more permissive, never less. It never reports success for a step that
// did not run when the seat is acting alone, which is the case that matters.
// It returns (landed, aborted). aborted means ctx was cancelled while waiting,
// which is a shutdown race rather than a failed step.
func (s *seat) awaitLogGrowth(ctx context.Context, want int) (landed, aborted bool) {
	deadline := time.Now().Add(s.r.cfg.StepTimeout)
	for {
		n, err := s.logLen()
		if err != nil {
			return false, false
		}
		if n >= want {
			return true, false
		}
		if ctx.Err() != nil {
			return false, true
		}
		if time.Now().After(deadline) {
			return false, false
		}
		runtime.Gosched()
	}
}
