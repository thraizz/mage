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

// Pacing controls how fast a bot acts, and is the per-seat human-likeness
// knob: two seats at one table can carry different Pacing values, which is how
// personas get different speeds and different voices.
//
// THE ZERO VALUE IS "NO DELAY AT ALL". Human is nil, the three fixed durations
// are zero, and every wait below short-circuits to runtime.Gosched(). That is
// not an accident to be tidied away later: Phase 3's completion-rate test runs
// a hundred full games and `make bot-sim` runs it from the Makefile, and
// human-like delays would turn a few seconds into hours. Paced runs opt in.
//
// Human, when non-nil, supersedes the fixed durations with log-normal draws
// weighted by decision type (see pace.go). The fixed fields remain for a caller
// who wants a flat delay without a distribution -- a slowed-down debugging run,
// say -- and for the zero case.
//
// Every wait these fields drive is performed by the bot's own goroutine between
// SendPlayerAction calls. None of them is reachable from a notification handler
// or an engine callback; see the pace.go header and anti-pattern 2.
type Pacing struct {
	// Poll is how long to wait between GetGameView calls when there is nothing
	// to do. Zero yields to the scheduler instead of sleeping.
	Poll time.Duration
	// PreAction is the "thinking" delay before a macro's first step.
	PreAction time.Duration
	// BetweenSteps staggers a macro's steps so a watching client animates each
	// one separately.
	BetweenSteps time.Duration

	// Human, when non-nil, replaces the three durations above with log-normal
	// draws. It also owns the hesitation and chat-frequency rolls.
	Human *HumanPacer

	// Chat is where this seat's lines come from. Nil means the seat plays in
	// silence. Phase 5 swaps CannedChat for the model's own `why` text without
	// touching the runner.
	Chat ChatSource
}

// preActionDelay is the "thinking" pause before a macro of this kind, weighted
// by how hard the decision would be for a person.
func (p Pacing) preActionDelay(k Kind) time.Duration {
	if p.Human != nil {
		return p.Human.PreActionDelay(k)
	}
	return p.PreAction
}

// stepDelay is the stagger between two commands of one macro.
func (p Pacing) stepDelay() time.Duration {
	if p.Human != nil {
		return p.Human.StepDelay()
	}
	return p.BetweenSteps
}

// pollDelay is the idle sleep while another seat has the turn.
func (p Pacing) pollDelay() time.Duration {
	if p.Human != nil {
		return p.Human.PollDelay()
	}
	return p.Poll
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

	// Pacing is the default for seats added with AddSeat. AddSeatWithPacing
	// overrides it per seat.
	Pacing Pacing

	// Chat is the table-wide half of chat: where lines go. The other half --
	// what a given seat says -- lives on that seat's Pacing.Chat, because voice
	// is a per-persona property and the sink is not.
	Chat ChatConfig

	// Presence is the activity/heartbeat signal a human client emits, if this
	// deployment has one for a bot to emit too. See ChatConfig's note and the
	// Presence doc: for an in-process headless bot this is nil, and there is
	// nothing to emit.
	Presence Presence

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

// ChatConfig points a runner's seats at a chat room.
//
// Sink is satisfied directly by *chat.Manager (internal/chat/manager.go:152:
// SendMessage(roomID, username, text string) error). RoomID is the table's chat
// room -- the same ID grpc_table.go passes to chatMgr, which for a real table
// is the table's RoomID.
//
// Both zero means no chat, which is the headless-simulation case.
type ChatConfig struct {
	Sink   ChatSink
	RoomID string
}

func (c ChatConfig) enabled() bool { return c.Sink != nil && c.RoomID != "" }

// Presence is the activity/heartbeat signal a human client emits so the server
// does not expire its session.
//
// FINDING (Phase 4, task 4): a headless bot has nothing to emit. The only
// activity signal in this codebase is session.Manager.UpdateActivity /
// Session.UpdateActivity (internal/session/session.go:73), which pushes out a
// session's lease so CleanupExpiredSessions does not reap it. It is driven from
// exactly three places, all of them gRPC entry points
// (internal/server/interceptors.go:105, internal/server/grpc.go:961, :1013,
// :1031) and all of them keyed by a session ID that a gRPC client established.
// A bot that lives in the server process has no session, no lease, and nothing
// to keep alive; there is also no client-visible presence indicator to feed --
// nothing in mage-client-web renders LastActivity or an online/typing state.
//
// So this is an interface with no production implementation on purpose. It is
// the seam for the day bots run out-of-process as real gRPC clients (Open
// Question 2 in the plan), at which point the transport's own keepalive is what
// implements it. Building a fake session for an in-process bot would create a
// lease to expire and a reaper to fight, for no observable benefit.
type Presence interface {
	// Touch records that the bot is still alive and playing.
	Touch()
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

// AddSeat registers a bot seat and the policy that drives it, using the
// runner's default Pacing.
func (r *BotRunner) AddSeat(botID string, p Policy) {
	r.AddSeatWithPacing(botID, p, r.cfg.Pacing)
}

// AddSeatWithPacing registers a seat with its own pacing.
//
// This is how personas differ: give each seat a HumanPacer over its own profile
// and its own seed, and four bots at one table think, fidget and talk at four
// different rhythms instead of moving in lockstep. Lockstep is a tell.
func (r *BotRunner) AddSeatWithPacing(botID string, p Policy, pacing Pacing) {
	r.seats = append(r.seats, &seat{r: r, botID: botID, policy: p, pacing: pacing})
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
	pacing Pacing

	stats Stats

	// chatTurn / chatCount enforce MAX_CHAT_MESSAGES_PER_TURN, and
	// lastChatTurn enforces the "at least one message per 2 turn cycles"
	// floor. All three are touched only by this seat's goroutine.
	chatTurn     int
	chatCount    int
	lastChatTurn int

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
		// Presence, such as it is. Touch is a no-op seam for an in-process bot
		// -- see the Presence doc -- but it belongs here, in the seat's own
		// loop, and not on any path the engine could be holding a lock on.
		if s.r.cfg.Presence != nil {
			s.r.cfg.Presence.Touch()
		}

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
			if !s.pacing.wait(ctx, s.pacing.pollDelay()) {
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
		if !s.pause(ctx, KindPassTurn) {
			return
		}
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
	if !s.pause(ctx, KindPassTurn) {
		return
	}
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
		if !s.pause(ctx, KindUntap) {
			return
		}
		s.exec(ctx, macro(KindUntap, fmt.Sprintf("Untap %d permanent(s)", len(steps)), steps...))
	}
	if v.Me.LibraryCount > 0 {
		if !s.pause(ctx, KindDraw) {
			return
		}
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
	// Fidget first, then think, then act -- the order a person does it in.
	// fidget is state-neutral (see its doc) so it cannot change what m means.
	if !s.fidget(ctx, v) {
		return false
	}
	if !s.pause(ctx, m.KindOf()) {
		return false
	}
	s.exec(ctx, m)
	if m.KindOf() == KindPlayLand {
		s.landDropTurn = v.Turn
	}
	s.maybeChat(ctx, v, m)
	return m.KindOf() != KindPassTurn
}

// pause is the pre-action "thinking" delay, weighted by how hard a person would
// find this decision. It runs in the seat's own goroutine, before the first
// SendPlayerAction of the macro, and reports whether the seat should continue.
func (s *seat) pause(ctx context.Context, k Kind) bool {
	d := s.pacing.preActionDelay(k)
	// A thinking policy has already spent some of this pause for real. Deduct
	// it -- see ThinkTimeReporter.
	if r, ok := s.policy.(ThinkTimeReporter); ok {
		if d -= r.LastThinkTime(); d < 0 {
			d = 0
		}
	}
	return s.pacing.wait(ctx, d)
}

// fidget occasionally taps one of the seat's untapped lands and immediately
// untaps it again, which to anyone watching the board reads as someone working
// out whether they can afford something.
//
// IT MUST NOT CHANGE THE GAME. Both verbs are implemented string commands
// (TAP:<id>, UNTAP:<id>, §0.6), the target is chosen from this seat's own
// untapped lands, and the untap restores exactly the card the tap touched --
// so the state after a fidget is identical to the state before it, and the
// mana solver in moves.go sees the same untapped sources either way. Nothing
// else is touched: no zone moves, no life, no cards that were already tapped.
//
// It returns false only if the run was cancelled mid-fidget.
func (s *seat) fidget(ctx context.Context, v *SafeView) bool {
	if s.pacing.Human == nil || !s.pacing.Human.Hesitate() {
		return ctx.Err() == nil
	}
	var lands []*SafeCard
	for _, c := range controlledBy(v.Battlefield, s.botID) {
		if !c.Tapped && IsLand(c) {
			lands = append(lands, c)
		}
	}
	if len(lands) == 0 {
		return ctx.Err() == nil
	}
	c := lands[s.pacing.Human.Pick(len(lands))]
	s.exec(ctx, macro(KindTap, "Fidget with "+displayName(c), "TAP:"+c.ID, "UNTAP:"+c.ID))
	return ctx.Err() == nil
}

// maybeChat applies mage-bench's chat cadence (§0.4 and reference/
// system-prompt.md "## Chat"): at most MaxChatMessagesPerTurn lines in one
// turn, and at least one line every ChatEveryNCycles turn cycles, where a cycle
// is one turn per seat.
//
// The line itself comes from the seat's ChatSource, which in Phase 4 is
// CannedChat and in Phase 5 is the model's own `why` text. The cadence lives
// here rather than in the source so that swapping the source cannot
// accidentally make a bot spam the table.
func (s *seat) maybeChat(ctx context.Context, v *SafeView, m Macro) {
	if s.pacing.Chat == nil || !s.r.cfg.Chat.enabled() || v == nil {
		return
	}
	if v.Turn != s.chatTurn {
		s.chatTurn = v.Turn
		s.chatCount = 0
	}
	if s.chatCount >= MaxChatMessagesPerTurn {
		return
	}

	cycle := 1 + len(v.Opponents)
	due := v.Turn-s.lastChatTurn >= ChatEveryNCycles*cycle

	line, ok := s.pacing.Chat.Line(ctx, v, m, due)
	if !ok || line == "" {
		return
	}
	// chat.Manager takes its own lock and never touches the engine, so this
	// cannot stall a broadcast. It still runs in the seat's goroutine.
	if err := s.r.cfg.Chat.Sink.SendMessage(s.r.cfg.Chat.RoomID, s.botID, line); err != nil {
		s.r.cfg.Logger.Warn("bot: chat failed",
			zap.String("bot", s.botID), zap.Error(err))
		return
	}
	s.chatCount++
	s.lastChatTurn = v.Turn
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
		// The stagger is what makes a tap-out visible: every step broadcasts
		// (game_engine.go:342), so a client animates each one separately only
		// if they do not arrive in the same frame. This sleep is in the bot's
		// goroutine BETWEEN SendPlayerAction calls -- never inside a handler,
		// never while the engine holds a lock (anti-pattern 2).
		if i > 0 && !s.pacing.wait(ctx, s.pacing.stepDelay()) {
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

// ThinkTimeReporter is implemented by a Policy whose Pick spends real
// wall-clock time deciding -- which is to say, by LLMPolicy.
//
// PACING MUST OVERLAP THINKING, NOT FOLLOW IT (plan Phase 5, task 4). A
// persona's pre-action pause exists to make the seat look like a person
// weighing a decision. An LLM seat is already weighing it, for seconds at a
// time; adding the full pause on top would make it slower than a human rather
// than more like one, and the tell would be the wrong way round.
//
// So a policy that thinks reports how long it thought, and seat.pause spends
// only the remainder. RandomPolicy does not implement this, returns nothing,
// and is paced exactly as it was in Phase 4.
type ThinkTimeReporter interface {
	// LastThinkTime is how long the most recent Pick took.
	LastThinkTime() time.Duration
}
