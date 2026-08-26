package bot

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"time"
)

// pace.go makes a bot look like a person to anyone watching the real client.
//
// This is an EXECUTOR problem, not a policy problem, which is why it lands
// before the LLM exists. A model that picks brilliant lines but fires four
// commands inside the same millisecond still reads as a robot: the client
// animates each broadcast (internal/game/game_engine.go:342), and four
// broadcasts in one frame collapse into one visual jump. Conversely a random
// policy that taps its lands one at a time, pauses before a block, and
// occasionally stares at the board for half a minute reads as a distracted
// human. Phase 5 only has to add decision quality on top of this.
//
// WHERE THE DELAYS HAPPEN. Every wait in this file is performed by the bot's
// own goroutine, between SendPlayerAction calls, in seat.act / seat.exec /
// seat.run. Nothing here is ever called from a notification handler or an
// engine callback -- see anti-pattern 2 in docs/BOT_PLAYERS_PLAN.md §0.8:
// GameEngine.broadcast runs with e.mu held for WRITING, so a sleep on that path
// would stall every mutation on the table. The runner does not subscribe to the
// push path at all (see the header comment in runner.go), so there is no
// callback for a sleep to leak into.
//
// WHY LOG-NORMAL. Uniform jitter is the classic mistake: it is bounded on both
// ends, so the bot never once agonises, and a viewer reads the tight band as
// "machine with a delay()" rather than "person". Human decision latency is
// right-skewed -- most decisions are fast, a few take an order of magnitude
// longer -- which is what a log-normal describes. On top of the per-kind
// log-normal this pacer mixes in a low-probability "tank" component (§ Phase 4:
// ~5% of decisions blow past 25s) because the natural sigma that would produce
// a 25s tail from a 2.4s median would also smear the everyday case into
// uselessness. A two-component log-normal mixture keeps the body tight and the
// tail long, which is what the real thing looks like.

// Chat cadence, from mage-bench's measured constants (§0.4) and the "## Chat"
// section of reference/system-prompt.md: talk at least once every two turn
// cycles, and never more than twice in one turn.
const (
	// MaxChatMessagesPerTurn is upstream's MAX_CHAT_MESSAGES_PER_TURN.
	MaxChatMessagesPerTurn = 2
	// ChatEveryNCycles is the upper bound on silence, in turn cycles, where a
	// cycle is one turn for every seat at the table.
	ChatEveryNCycles = 2
)

// DelayBand is one log-normal component: a median and a spread on the log
// scale.
//
// Median is the p50 exactly (for a log-normal the median is exp(mu), so mu is
// just log(Median) and no mean/median correction is needed). Sigma is the
// standard deviation of the underlying normal, so the p95 is
// Median*exp(1.645*Sigma) and the p05 is Median/exp(1.645*Sigma) -- the band is
// multiplicative and symmetric in log space, which is the whole point.
type DelayBand struct {
	Median time.Duration
	Sigma  float64
}

// PacingProfile is one bot's sense of time. It is a plain value so a caller can
// copy a preset and tweak a field; two seats at the same table can and should
// carry different profiles.
type PacingProfile struct {
	// Name is for logs and for persona-flavoured chat.
	Name string

	// Bands is the pre-action "thinking" delay per decision class. A Kind with
	// no entry falls back to Default.
	Bands map[Kind]DelayBand
	// Default covers any Kind missing from Bands.
	Default DelayBand

	// Step staggers the commands within one macro, so a three-land tap-out
	// visibly taps three lands instead of blinking.
	Step DelayBand
	// StepMin and StepMax clamp a Step draw. The stagger has to stay inside the
	// window where it reads as "a hand moving" rather than "a stall".
	StepMin, StepMax time.Duration

	// Poll is the idle sleep between GetGameView calls while it is somebody
	// else's turn. Invisible to a viewer; it exists so an idle seat does not
	// spin a core.
	Poll time.Duration

	// TankChance is the probability that a decision additionally draws from the
	// tank component -- the ~5% of turns where a human gets up, re-reads a
	// card, or has to think properly.
	TankChance float64
	// TankMedian / TankSigma describe that component.
	TankMedian time.Duration
	TankSigma  float64

	// HesitateChance is the probability of a visible fidget (tap a land, untap
	// it again) before acting. See seat.fidget.
	HesitateChance float64

	// ChatChance is the probability of an unprompted line after an action. The
	// cadence floor in ChatEveryNCycles is enforced separately and overrides
	// this; the cap in MaxChatMessagesPerTurn overrides both.
	ChatChance float64

	// MaxDelay caps every draw. A log-normal is unbounded above and one draw in
	// a few thousand would otherwise park the table for ten minutes.
	MaxDelay time.Duration

	// Speed scales every duration in the profile. 1 is the profile as written,
	// 0.5 is twice as fast. It is the one dial a caller needs to make a paced
	// sim finish in test time without losing the shape of the distribution.
	Speed float64
}

// AveragePace is the reference profile: the medians below are the midpoints of
// the bands Phase 4 specifies (land drop / untap 0.5-1.5s, routine cast 1.5-4s,
// combat or targeted removal 4-10s, mulligan 8-20s), and each Sigma is chosen
// so that the p05..p95 interval lands roughly on the quoted band.
func AveragePace() PacingProfile {
	routine := DelayBand{Median: 2400 * time.Millisecond, Sigma: 0.42}
	quick := DelayBand{Median: 800 * time.Millisecond, Sigma: 0.35}
	// "Combat or targeted removal": in this rules-light engine removal is
	// expressed as a manual zone move, so KindMoveZone belongs here rather than
	// with the routine plays.
	weighty := DelayBand{Median: 6300 * time.Millisecond, Sigma: 0.30}
	// Mulligan is the one decision with no board to read, so it is pure
	// deliberation and the slowest thing a human does all game.
	agonising := DelayBand{Median: 12 * time.Second, Sigma: 0.35}

	return PacingProfile{
		Name: "average",
		Bands: map[Kind]DelayBand{
			KindPlayLand:   quick,
			KindUntap:      quick,
			KindTap:        quick,
			KindDraw:       quick,
			KindCast:       routine,
			KindModifyLife: routine,
			KindPassTurn:   DelayBand{Median: 1200 * time.Millisecond, Sigma: 0.40},
			KindAttack:     weighty,
			KindMoveZone:   weighty,
			KindMulligan:   agonising,
			KindKeepHand:   agonising,
		},
		Default:        routine,
		Step:           DelayBand{Median: 320 * time.Millisecond, Sigma: 0.30},
		StepMin:        120 * time.Millisecond,
		StepMax:        900 * time.Millisecond,
		Poll:           250 * time.Millisecond,
		TankChance:     0.05,
		TankMedian:     30 * time.Second,
		TankSigma:      0.35,
		HesitateChance: 0.07,
		ChatChance:     0.15,
		MaxDelay:       45 * time.Second,
		Speed:          1,
	}
}

// BriskPace is a fast, decisive persona. Same shape, less of it.
func BriskPace() PacingProfile {
	p := AveragePace()
	p.Name = "brisk"
	p.Speed = 0.55
	p.TankChance = 0.02
	p.HesitateChance = 0.03
	p.ChatChance = 0.10
	return p
}

// DeliberatePace is a slow, chatty, fidgety persona.
func DeliberatePace() PacingProfile {
	p := AveragePace()
	p.Name = "deliberate"
	p.Speed = 1.6
	p.TankChance = 0.08
	p.HesitateChance = 0.12
	p.ChatChance = 0.25
	return p
}

// HumanPacer draws the delays for one bot.
//
// The RNG is injectable and seeded, exactly as RandomPolicy's is: a simulation
// whose timing cannot be replayed is not a debugging tool. Every draw goes
// through this one lock, so a pacer shared between seats stays deterministic
// only in aggregate -- give each seat its own.
type HumanPacer struct {
	mu  sync.Mutex
	rng *rand.Rand
	p   PacingProfile
}

// NewHumanPacer returns a pacer over profile, seeded with seed.
func NewHumanPacer(profile PacingProfile, seed int64) *HumanPacer {
	return NewHumanPacerWithRand(profile, rand.New(rand.NewSource(seed)))
}

// NewHumanPacerWithRand returns a pacer over a caller-supplied source.
func NewHumanPacerWithRand(profile PacingProfile, r *rand.Rand) *HumanPacer {
	if profile.Speed <= 0 {
		profile.Speed = 1
	}
	if profile.MaxDelay <= 0 {
		profile.MaxDelay = 45 * time.Second
	}
	return &HumanPacer{rng: r, p: profile}
}

// Profile returns the pacer's profile.
func (h *HumanPacer) Profile() PacingProfile {
	if h == nil {
		return PacingProfile{}
	}
	return h.p
}

// PreActionDelay is how long the bot appears to think before committing to a
// macro of this kind.
func (h *HumanPacer) PreActionDelay(k Kind) time.Duration {
	if h == nil {
		return 0
	}
	h.mu.Lock()
	defer h.mu.Unlock()

	band, ok := h.p.Bands[k]
	if !ok {
		band = h.p.Default
	}
	d := h.drawLocked(band)

	// The tank component. Taking the max rather than adding keeps the body of
	// the distribution untouched: a tanked decision is a long one, not an
	// ordinary one plus a constant.
	if h.p.TankChance > 0 && h.rng.Float64() < h.p.TankChance {
		if t := h.drawLocked(DelayBand{Median: h.p.TankMedian, Sigma: h.p.TankSigma}); t > d {
			d = t
		}
	}
	return h.clamp(d)
}

// StepDelay is the stagger between two commands of the same macro.
func (h *HumanPacer) StepDelay() time.Duration {
	if h == nil {
		return 0
	}
	h.mu.Lock()
	defer h.mu.Unlock()

	d := h.drawLocked(h.p.Step)
	lo, hi := h.scale(h.p.StepMin), h.scale(h.p.StepMax)
	if lo > 0 && d < lo {
		d = lo
	}
	if hi > 0 && d > hi {
		d = hi
	}
	return h.clamp(d)
}

// PollDelay is the idle sleep while waiting for another seat.
func (h *HumanPacer) PollDelay() time.Duration {
	if h == nil || h.p.Poll <= 0 {
		return 0
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	// Light jitter only. Nobody watches this, but identical poll periods across
	// four seats produce a thundering herd on the engine's read lock.
	return h.clamp(h.drawLocked(DelayBand{Median: h.p.Poll, Sigma: 0.25}))
}

// Hesitate reports whether the bot should visibly fidget before its next
// action.
func (h *HumanPacer) Hesitate() bool { return h.roll(h.hesitateChance()) }

// ChatRoll reports whether the bot feels like saying something unprompted.
func (h *HumanPacer) ChatRoll() bool { return h.roll(h.chatChance()) }

// Pick returns a uniformly random index in [0,n), for callers that need a
// choice on the pacer's own deterministic stream (canned chat lines, which land
// of several to fidget with). Returns 0 for n <= 0.
func (h *HumanPacer) Pick(n int) int {
	if h == nil || n <= 0 {
		return 0
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.rng.Intn(n)
}

func (h *HumanPacer) hesitateChance() float64 {
	if h == nil {
		return 0
	}
	return h.p.HesitateChance
}

func (h *HumanPacer) chatChance() float64 {
	if h == nil {
		return 0
	}
	return h.p.ChatChance
}

func (h *HumanPacer) roll(chance float64) bool {
	if h == nil || chance <= 0 {
		return false
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.rng.Float64() < chance
}

// drawLocked samples one log-normal delay. Caller holds h.mu.
func (h *HumanPacer) drawLocked(b DelayBand) time.Duration {
	median := h.scale(b.Median)
	if median <= 0 {
		return 0
	}
	if b.Sigma <= 0 {
		return median
	}
	// X = median * exp(sigma*Z), Z ~ N(0,1). log(X) is normal with mean
	// log(median) -- so median is the p50, by construction.
	return time.Duration(float64(median) * math.Exp(b.Sigma*h.rng.NormFloat64()))
}

func (h *HumanPacer) scale(d time.Duration) time.Duration {
	if d <= 0 {
		return 0
	}
	return time.Duration(float64(d) * h.p.Speed)
}

func (h *HumanPacer) clamp(d time.Duration) time.Duration {
	if d < 0 {
		return 0
	}
	if max := h.scale(h.p.MaxDelay); max > 0 && d > max {
		return max
	}
	return d
}

// ---------------------------------------------------------------------------
// Chat
// ---------------------------------------------------------------------------

// ChatSink is the subset of *chat.Manager (internal/chat/manager.go:152) the
// bot needs: SendMessage(roomID, username, text) error.
//
// It is an interface rather than the concrete manager for the usual reason --
// internal/bot must stay testable without a session manager -- and because a
// headless simulation has no chat manager at all, in which case this is nil and
// the bot plays in silence.
//
// NOTE ON THE ROOM. chat.Manager.SendMessage is a no-op (a warning, and a nil
// error) when the room does not exist, and it fans out only to sessions of
// users who have JoinRoom'd. So a bot's line reaches a human iff the table's
// chat room already exists and the human is in it -- which is exactly the state
// the gRPC table path produces (internal/server/grpc_table.go:185). The bot
// itself never joins: it has no session to deliver anything to.
type ChatSink interface {
	SendMessage(roomID, username, text string) error
}

// ChatSource produces the line a bot says.
//
// This is the seam Phase 5 replaces. Today it is CannedChat over a small pool
// of persona-flavoured lines; in Phase 5 the same interface is fed by the
// model's own `why` text, and nothing at the call site changes. The `due` flag
// tells the source that the cadence floor has been reached and silence is no
// longer acceptable -- a source may return false when due is false, but should
// try to produce something when it is true.
type ChatSource interface {
	Line(ctx context.Context, v *SafeView, m Macro, due bool) (string, bool)
}

// CannedChat is the Phase 4 ChatSource: a small varied pool, keyed by the kind
// of decision just made so the line at least relates to what happened.
//
// It draws from the pacer's RNG, so a seeded simulation says the same things in
// the same order -- and so the chat roll and the delay draws share one
// reproducible stream.
type CannedChat struct {
	pacer  *HumanPacer
	byKind map[Kind][]string
	filler []string
	// last is the previously emitted line, so the bot does not say the same
	// thing twice in a row. Guarded by the pacer's lock indirectly: only the
	// owning seat's goroutine touches it.
	last string
}

// NewCannedChat builds a canned chat source flavoured by the pacer's profile
// name. A nil pacer yields a source that never speaks.
func NewCannedChat(p *HumanPacer) *CannedChat {
	c := &CannedChat{pacer: p, byKind: cannedByKind(), filler: cannedFiller()}
	if p != nil {
		if extra, ok := personaFiller[p.Profile().Name]; ok {
			c.filler = append(append([]string(nil), c.filler...), extra...)
		}
	}
	return c
}

// Line implements ChatSource.
func (c *CannedChat) Line(_ context.Context, _ *SafeView, m Macro, due bool) (string, bool) {
	if c == nil || c.pacer == nil {
		return "", false
	}
	if !due && !c.pacer.ChatRoll() {
		return "", false
	}
	pool := c.filler
	if lines, ok := c.byKind[m.KindOf()]; ok && len(lines) > 0 {
		// Half the time comment on the play, half the time say something
		// generic. Always commenting on the play is its own tell.
		if c.pacer.Pick(2) == 0 {
			pool = lines
		}
	}
	if len(pool) == 0 {
		return "", false
	}
	i := c.pacer.Pick(len(pool))
	if pool[i] == c.last && len(pool) > 1 {
		// Re-roll over the pool minus the repeat, by offsetting past it. The
		// same line twice in a row is the loudest tell a scripted talker has.
		i = (i + 1 + c.pacer.Pick(len(pool)-1)) % len(pool)
	}
	line := pool[i]
	c.last = line
	return line, true
}

func cannedByKind() map[Kind][]string {
	return map[Kind][]string{
		KindPlayLand:   {"Land, go.", "Just a land for me.", "Hitting my drop, at least."},
		KindCast:       {"Let's try this.", "This should do something.", "Resolving, I hope?"},
		KindAttack:     {"Swinging.", "Attacks!", "I'll send this one in."},
		KindMoveZone:   {"That has to go.", "Sorry about that one.", "Handling it."},
		KindMulligan:   {"No lands, back it goes.", "Can't keep that.", "Mull."},
		KindKeepHand:   {"Keeping.", "That'll do.", "Happy with this one."},
		KindPassTurn:   {"Pass.", "Your turn.", "Go ahead."},
		KindModifyLife: {"Paying for it.", "Worth it.", "Ouch."},
	}
}

func cannedFiller() []string {
	return []string{
		"gl hf",
		"Nice board.",
		"This is getting scary.",
		"Anyone got an answer for that?",
		"I'm not liking my odds here.",
		"Good draw.",
		"Well played.",
		"Sorry, thinking.",
	}
}

// personaFiller adds a few lines that only one persona says, so two bots at the
// same table do not sound identical.
var personaFiller = map[string][]string{
	"brisk":      {"Quick one.", "Yep, done.", "Moving."},
	"deliberate": {"Give me a second.", "Let me count this out.", "Hmm."},
}
