package bot

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/magefree/mage-server-go/internal/chat"
	"github.com/magefree/mage-server-go/internal/game"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	i := int(p * float64(len(sorted)))
	if i >= len(sorted) {
		i = len(sorted) - 1
	}
	return sorted[i]
}

func drawMany(h *HumanPacer, k Kind, n int) []time.Duration {
	out := make([]time.Duration, n)
	for i := range out {
		out[i] = h.PreActionDelay(k)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// ---------------------------------------------------------------------------
// Zero value
// ---------------------------------------------------------------------------

// TestZeroPacingIsFree is the gate on Phase 3's completion-rate test: a
// zero-value Pacing must not cost measurable time, or `make bot-sim` goes from
// seconds to hours.
func TestZeroPacingIsFree(t *testing.T) {
	var p Pacing

	if got := p.preActionDelay(KindMulligan); got != 0 {
		t.Errorf("zero Pacing preActionDelay = %v, want 0", got)
	}
	if got := p.stepDelay(); got != 0 {
		t.Errorf("zero Pacing stepDelay = %v, want 0", got)
	}
	if got := p.pollDelay(); got != 0 {
		t.Errorf("zero Pacing pollDelay = %v, want 0", got)
	}

	// 30k waits is more than a 100-game sim performs. If the zero path ever
	// starts sleeping, even for a microsecond, this blows up by three orders of
	// magnitude.
	const n = 30000
	ctx := context.Background()
	start := time.Now()
	for i := 0; i < n; i++ {
		if !p.wait(ctx, p.preActionDelay(KindCast)) {
			t.Fatal("wait reported cancellation on a live context")
		}
	}
	elapsed := time.Since(start)
	t.Logf("%d zero-value waits in %v (%v each)", n, elapsed, elapsed/n)
	if elapsed > 500*time.Millisecond {
		t.Fatalf("%d zero-value waits took %v; the zero path is sleeping", n, elapsed)
	}
}

// TestNilPacerIsSafe: a nil *HumanPacer must behave like "no pacing" rather
// than panic, so a caller can leave the field unset.
func TestNilPacerIsSafe(t *testing.T) {
	var h *HumanPacer
	if got := h.PreActionDelay(KindCast); got != 0 {
		t.Errorf("nil pacer PreActionDelay = %v, want 0", got)
	}
	if got := h.StepDelay(); got != 0 {
		t.Errorf("nil pacer StepDelay = %v, want 0", got)
	}
	if got := h.PollDelay(); got != 0 {
		t.Errorf("nil pacer PollDelay = %v, want 0", got)
	}
	if h.Hesitate() || h.ChatRoll() {
		t.Error("nil pacer should never roll true")
	}
	if h.Pick(5) != 0 {
		t.Error("nil pacer Pick should be 0")
	}
}

// ---------------------------------------------------------------------------
// Distribution shape
// ---------------------------------------------------------------------------

// TestDelayIsLogNormal checks the BODY of the distribution: with the tank
// component disabled, log(delay) must look normal -- symmetric around its mean,
// with the sigma the band asked for. Bands, not exact values.
func TestDelayIsLogNormal(t *testing.T) {
	p := AveragePace()
	p.TankChance = 0
	h := NewHumanPacer(p, 20240826)

	const n = 4000
	samples := drawMany(h, KindCast, n)

	logs := make([]float64, n)
	for i, d := range samples {
		logs[i] = math.Log(float64(d) / float64(time.Second))
	}

	var mean float64
	for _, x := range logs {
		mean += x
	}
	mean /= float64(n)

	var m2, m3 float64
	for _, x := range logs {
		d := x - mean
		m2 += d * d
		m3 += d * d * d
	}
	m2 /= float64(n)
	m3 /= float64(n)
	sd := math.Sqrt(m2)
	skew := m3 / (sd * sd * sd)

	medLog := math.Log(float64(percentile(samples, 0.50)) / float64(time.Second))
	wantMedLog := math.Log(2.4)

	t.Logf("log-scale: mean=%.3f median=%.3f sd=%.3f skew=%.3f (want mean~median~%.3f, sd~0.42, |skew|<0.3)",
		mean, medLog, sd, skew, wantMedLog)

	// A log-normal has mean(log) == median(log) exactly, and both equal
	// log(Median). Allow a fifth of a sigma of sampling slop.
	if math.Abs(mean-wantMedLog) > 0.2*sd {
		t.Errorf("log-mean %.3f is not at log(median) %.3f", mean, wantMedLog)
	}
	if math.Abs(medLog-mean) > 0.2*sd {
		t.Errorf("log-median %.3f differs from log-mean %.3f: not symmetric in log space", medLog, mean)
	}
	if sd < 0.34 || sd > 0.50 {
		t.Errorf("log-sd %.3f outside [0.34,0.50] for Sigma=0.42", sd)
	}
	if math.Abs(skew) > 0.3 {
		t.Errorf("log-skew %.3f: log(delay) is not symmetric, so delay is not log-normal", skew)
	}

	// And on the linear scale it must be right-skewed: mean > median. That is
	// the property a uniform distribution does not have and the reason uniform
	// reads as robotic.
	var lin float64
	for _, d := range samples {
		lin += float64(d)
	}
	linMean := time.Duration(lin / float64(n))
	if linMean <= percentile(samples, 0.50) {
		t.Errorf("mean %v <= median %v: distribution is not right-skewed", linMean, percentile(samples, 0.50))
	}
}

// TestDelayLongTail is the one that matters by eye: ~5% of decisions have to
// blow past 25 seconds, or the bot never once looks like it is actually
// thinking.
func TestDelayLongTail(t *testing.T) {
	h := NewHumanPacer(AveragePace(), 7)

	const n = 1000
	samples := drawMany(h, KindCast, n)

	p50 := percentile(samples, 0.50)
	p95 := percentile(samples, 0.95)
	max := samples[n-1]

	tanked := 0
	for _, d := range samples {
		if d >= 20*time.Second {
			tanked++
		}
	}
	frac := float64(tanked) / float64(n)

	t.Logf("KindCast n=%d: p50=%v p95=%v max=%v  >=20s: %d (%.1f%%)",
		n, p50.Round(time.Millisecond), p95.Round(time.Millisecond),
		max.Round(time.Millisecond), tanked, frac*100)

	if p50 < 1600*time.Millisecond || p50 > 3600*time.Millisecond {
		t.Errorf("p50 %v outside [1.6s,3.6s] for a 2.4s median", p50)
	}
	if p95 < 4*time.Second || p95 > 25*time.Second {
		t.Errorf("p95 %v outside [4s,25s]", p95)
	}
	if max < 25*time.Second {
		t.Errorf("max %v < 25s: the long tail never happened", max)
	}
	if max > AveragePace().MaxDelay {
		t.Errorf("max %v exceeds MaxDelay %v", max, AveragePace().MaxDelay)
	}
	if frac < 0.015 || frac > 0.10 {
		t.Errorf("tank fraction %.3f outside [0.015,0.10]", frac)
	}
}

// TestDelayBandsByKind: the plan's ordering -- land drop fastest, then routine
// cast, then combat, then mulligan. Asserted on p50 with the tail off, so this
// tests the weighting and not the mixture.
func TestDelayBandsByKind(t *testing.T) {
	p := AveragePace()
	p.TankChance = 0

	type want struct {
		kind   Kind
		lo, hi time.Duration
	}
	// Bands from Phase 4: land 0.5-1.5s, cast 1.5-4s, combat 4-10s,
	// mulligan 8-20s. Asserted on p05/p95, which is where the band lives.
	cases := []want{
		{KindPlayLand, 400 * time.Millisecond, 1700 * time.Millisecond},
		{KindCast, 1200 * time.Millisecond, 5 * time.Second},
		{KindAttack, 3500 * time.Millisecond, 11 * time.Second},
		{KindMoveZone, 3500 * time.Millisecond, 11 * time.Second},
		{KindMulligan, 6 * time.Second, 23 * time.Second},
	}

	var prev time.Duration
	for _, c := range cases {
		h := NewHumanPacer(p, 99)
		s := drawMany(h, c.kind, 2000)
		p05, p50, p95 := percentile(s, 0.05), percentile(s, 0.50), percentile(s, 0.95)
		t.Logf("%-12s p05=%v p50=%v p95=%v", c.kind,
			p05.Round(time.Millisecond), p50.Round(time.Millisecond), p95.Round(time.Millisecond))
		if p05 < c.lo {
			t.Errorf("%s p05 %v below band floor %v", c.kind, p05, c.lo)
		}
		if p95 > c.hi {
			t.Errorf("%s p95 %v above band ceiling %v", c.kind, p95, c.hi)
		}
		if c.kind != KindMoveZone && p50 <= prev {
			t.Errorf("%s p50 %v is not slower than the previous class (%v)", c.kind, p50, prev)
		}
		if c.kind != KindMoveZone {
			prev = p50
		}
	}
}

// TestPersonaSpeeds: personas must actually differ, or every seat at the table
// moves in lockstep.
func TestPersonaSpeeds(t *testing.T) {
	median := func(p PacingProfile) time.Duration {
		p.TankChance = 0
		return percentile(drawMany(NewHumanPacer(p, 3), KindCast, 2000), 0.50)
	}
	brisk, avg, slow := median(BriskPace()), median(AveragePace()), median(DeliberatePace())
	t.Logf("brisk=%v average=%v deliberate=%v", brisk, avg, slow)
	if !(brisk < avg && avg < slow) {
		t.Errorf("persona speeds not ordered: brisk=%v average=%v deliberate=%v", brisk, avg, slow)
	}
}

// ---------------------------------------------------------------------------
// Reproducibility
// ---------------------------------------------------------------------------

// script is a representative decision sequence: a mulligan, a keep, then turns
// of untap / draw / land / cast / attack / pass. Feeding the same script to two
// pacers seeded alike must produce byte-identical timing, hesitation and chat
// decisions -- that is what makes a paced sim replayable.
func script() []Kind {
	s := []Kind{KindMulligan, KindKeepHand}
	for i := 0; i < 40; i++ {
		s = append(s, KindUntap, KindDraw, KindPlayLand, KindCast, KindAttack, KindPassTurn)
	}
	return s
}

type trace struct {
	pre  []time.Duration
	step []time.Duration
	poll []time.Duration
	hes  []bool
	chat []string
}

func runScript(seed int64) trace {
	h := NewHumanPacer(AveragePace(), seed)
	c := NewCannedChat(h)
	var tr trace
	for _, k := range script() {
		tr.hes = append(tr.hes, h.Hesitate())
		tr.pre = append(tr.pre, h.PreActionDelay(k))
		tr.step = append(tr.step, h.StepDelay(), h.StepDelay())
		tr.poll = append(tr.poll, h.PollDelay())
		if line, ok := c.Line(context.Background(), nil, macro(k, "x"), false); ok {
			tr.chat = append(tr.chat, line)
		}
	}
	return tr
}

func TestPacedSequenceIsReproducible(t *testing.T) {
	a, b := runScript(4242), runScript(4242)

	if len(a.pre) == 0 {
		t.Fatal("empty trace")
	}
	if fmt.Sprint(a.pre) != fmt.Sprint(b.pre) {
		t.Fatal("pre-action delay sequence differs between two runs with the same seed")
	}
	if fmt.Sprint(a.step) != fmt.Sprint(b.step) {
		t.Fatal("step stagger sequence differs between two runs with the same seed")
	}
	if fmt.Sprint(a.poll) != fmt.Sprint(b.poll) {
		t.Fatal("poll sequence differs between two runs with the same seed")
	}
	if fmt.Sprint(a.hes) != fmt.Sprint(b.hes) {
		t.Fatal("hesitation sequence differs between two runs with the same seed")
	}
	if fmt.Sprint(a.chat) != fmt.Sprint(b.chat) {
		t.Fatal("chat sequence differs between two runs with the same seed")
	}
	t.Logf("%d pre-action delays, %d chat lines reproduced exactly; first five delays: %v",
		len(a.pre), len(a.chat), a.pre[:5])

	// ...and a different seed must produce a different sequence, or the test
	// above is vacuous.
	if fmt.Sprint(a.pre) == fmt.Sprint(runScript(4243).pre) {
		t.Fatal("two different seeds produced the same delay sequence")
	}
}

// TestPacerAcceptsInjectedRand mirrors RandomPolicy's constructor pair.
func TestPacerAcceptsInjectedRand(t *testing.T) {
	a := NewHumanPacerWithRand(AveragePace(), rand.New(rand.NewSource(11)))
	b := NewHumanPacer(AveragePace(), 11)
	for i := 0; i < 50; i++ {
		if a.PreActionDelay(KindCast) != b.PreActionDelay(KindCast) {
			t.Fatalf("injected rand diverged at draw %d", i)
		}
	}
}

// ---------------------------------------------------------------------------
// Chat
// ---------------------------------------------------------------------------

// The real chat manager must satisfy ChatSink as written. This is the whole
// "wire internal/chat" claim, checked by the compiler rather than by eye:
// chat.Manager.SendMessage(roomID, username, text string) error
// (internal/chat/manager.go:152).
var _ ChatSink = (*chat.Manager)(nil)

type fakeSink struct {
	mu   sync.Mutex
	sent []string
	room []string
	user []string
}

func (f *fakeSink) SendMessage(roomID, username, text string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.room = append(f.room, roomID)
	f.user = append(f.user, username)
	f.sent = append(f.sent, text)
	return nil
}

func (f *fakeSink) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.sent)
}

func chatSeat(t *testing.T, sink ChatSink, profile PacingProfile) *seat {
	t.Helper()
	r := NewBotRunner(RunnerConfig{
		GameID: "g",
		Chat:   ChatConfig{Sink: sink, RoomID: "room-1"},
		Logger: zaptest.NewLogger(t),
	})
	h := NewHumanPacer(profile, 5)
	return &seat{r: r, botID: "bot-a", pacing: Pacing{Human: h, Chat: NewCannedChat(h)}}
}

func chatView(turn int) *SafeView {
	return &SafeView{
		Turn:      turn,
		Me:        &SafePlayerView{PlayerID: "bot-a"},
		Opponents: []*SafeOpponentView{{}, {}, {}},
	}
}

// TestChatCapPerTurn is mage-bench's MAX_CHAT_MESSAGES_PER_TURN.
func TestChatCapPerTurn(t *testing.T) {
	p := AveragePace()
	p.ChatChance = 1 // always wants to talk; the cap is the only brake
	sink := &fakeSink{}
	s := chatSeat(t, sink, p)

	for i := 0; i < 8; i++ {
		s.maybeChat(context.Background(), chatView(1), macro(KindCast, "x"))
	}
	if got := sink.count(); got != MaxChatMessagesPerTurn {
		t.Fatalf("turn 1: sent %d messages, want the cap of %d", got, MaxChatMessagesPerTurn)
	}
	for i := 0; i < 8; i++ {
		s.maybeChat(context.Background(), chatView(2), macro(KindCast, "x"))
	}
	if got := sink.count(); got != 2*MaxChatMessagesPerTurn {
		t.Fatalf("after turn 2: sent %d total, want %d", got, 2*MaxChatMessagesPerTurn)
	}
	if sink.room[0] != "room-1" || sink.user[0] != "bot-a" {
		t.Errorf("sent to room=%q as user=%q, want room-1/bot-a", sink.room[0], sink.user[0])
	}
}

// TestChatCadenceFloor is the "at least once every 2 turn cycles" rule: a bot
// that never feels like talking still has to talk.
func TestChatCadenceFloor(t *testing.T) {
	p := AveragePace()
	p.ChatChance = 0 // never volunteers; only the cadence floor can fire
	sink := &fakeSink{}
	s := chatSeat(t, sink, p)

	// Four seats, so a cycle is 4 turns and the floor is 8.
	for turn := 1; turn <= 7; turn++ {
		s.maybeChat(context.Background(), chatView(turn), macro(KindPassTurn, "x"))
	}
	if got := sink.count(); got != 0 {
		t.Fatalf("spoke %d times before the cadence floor; want silence", got)
	}
	s.maybeChat(context.Background(), chatView(8), macro(KindPassTurn, "x"))
	if got := sink.count(); got != 1 {
		t.Fatalf("sent %d at the cadence floor, want exactly 1 (%q)", got, sink.sent)
	}
	t.Logf("cadence floor fired at turn 8: %q", sink.sent[0])
}

// TestChatDisabledWithoutSink: a headless sim has no chat manager, and that
// must be silent rather than a nil dereference.
func TestChatDisabledWithoutSink(t *testing.T) {
	r := NewBotRunner(RunnerConfig{GameID: "g", Logger: zap.NewNop()})
	h := NewHumanPacer(AveragePace(), 1)
	s := &seat{r: r, botID: "bot-a", pacing: Pacing{Human: h, Chat: NewCannedChat(h)}}
	s.maybeChat(context.Background(), chatView(1), macro(KindCast, "x"))
	// No panic is the assertion.
}

// TestCannedChatIsVaried guards against the single most obvious bot tell: the
// same line twice in a row.
func TestCannedChatIsVaried(t *testing.T) {
	p := AveragePace()
	p.ChatChance = 1
	h := NewHumanPacer(p, 17)
	c := NewCannedChat(h)

	seen := map[string]int{}
	prev := ""
	for i := 0; i < 200; i++ {
		line, ok := c.Line(context.Background(), nil, macro(KindCast, "x"), true)
		if !ok {
			t.Fatal("due chat produced no line")
		}
		if line == prev {
			t.Fatalf("repeated %q back to back at %d", line, i)
		}
		prev = line
		seen[line]++
	}
	if len(seen) < 5 {
		t.Fatalf("only %d distinct lines in 200 draws", len(seen))
	}
	t.Logf("%d distinct canned lines over 200 draws", len(seen))
}

// TestPersonaFlavouredChat: two personas must not sound identical.
func TestPersonaFlavouredChat(t *testing.T) {
	pool := func(p PacingProfile) map[string]bool {
		p.ChatChance = 1
		h := NewHumanPacer(p, 8)
		c := NewCannedChat(h)
		out := map[string]bool{}
		for i := 0; i < 500; i++ {
			if line, ok := c.Line(context.Background(), nil, macro(KindPassTurn, "x"), true); ok {
				out[line] = true
			}
		}
		return out
	}
	brisk, slow := pool(BriskPace()), pool(DeliberatePace())
	onlySlow := 0
	for line := range slow {
		if !brisk[line] {
			onlySlow++
		}
	}
	if onlySlow == 0 {
		t.Fatal("deliberate persona has no lines the brisk persona lacks")
	}
	t.Logf("%d lines unique to the deliberate persona", onlySlow)
}

// ---------------------------------------------------------------------------
// Hesitation
// ---------------------------------------------------------------------------

// TestFidgetIsStateNeutral: a hesitation must untap exactly what it tapped and
// touch nothing else, or it is a cheat with extra steps.
func TestFidgetIsStateNeutral(t *testing.T) {
	logger := zaptest.NewLogger(t)
	seats := fourSeats()
	sg := newSimGame(t, logger, "table-fidget", seats)
	defer sg.close()

	me := seats[0]
	viewOf := func() *SafeView {
		t.Helper()
		raw, err := sg.adapter.GetGameView(sg.g.ID, me)
		if err != nil {
			t.Fatalf("view: %v", err)
		}
		sv, err := RedactErr(raw.(*game.PlaytestGameView), me)
		if err != nil {
			t.Fatalf("redact: %v", err)
		}
		Enrich(context.Background(), sv, simOracle())
		return sv
	}

	r := NewBotRunner(RunnerConfig{
		GameID:  sg.g.ID,
		Actions: sg.mgr,
		Views:   sg.adapter,
		Oracle:  simOracle(),
		Logger:  logger,
	})
	p := AveragePace()
	p.Speed = 0.0001 // the delays are not what is under test here
	p.HesitateChance = 1
	s := &seat{r: r, botID: me, pacing: Pacing{Human: NewHumanPacer(p, 1)}}

	// Keep, then get a land onto the battlefield to fidget with.
	if !s.exec(context.Background(), macro(KindKeepHand, "keep", "KEEP_HAND:"+me)) {
		t.Fatal("keep hand did not land")
	}
	v := viewOf()
	var land *SafeCard
	for _, c := range v.Me.Hand {
		if IsLand(c) {
			land = c
			break
		}
	}
	if land == nil {
		t.Skip("no land in the opening hand for this fixture")
	}
	if !s.exec(context.Background(), macro(KindPlayLand, "land", "MOVE:"+land.ID+":BATTLEFIELD")) {
		t.Fatal("land drop did not land")
	}

	before := tapState(viewOf(), me)
	if len(before) == 0 {
		t.Fatal("no permanents to fidget with")
	}
	for i := 0; i < 5; i++ {
		if !s.fidget(context.Background(), viewOf()) {
			t.Fatal("fidget reported cancellation")
		}
	}
	after := tapState(viewOf(), me)

	if len(before) != len(after) {
		t.Fatalf("fidget changed the board: %d permanents before, %d after", len(before), len(after))
	}
	for id, tapped := range before {
		if after[id] != tapped {
			t.Errorf("fidget left %s tapped=%v, was %v", id, after[id], tapped)
		}
	}
	if s.stats.MacrosFailed != 0 {
		t.Errorf("fidget macros failed: %v", s.stats.FailedMacros)
	}
	t.Logf("5 fidgets over %d permanents left tap state unchanged", len(before))
}

func tapState(v *SafeView, playerID string) map[string]bool {
	out := map[string]bool{}
	for _, c := range controlledBy(v.Battlefield, playerID) {
		out[c.ID] = c.Tapped
	}
	return out
}

// ---------------------------------------------------------------------------
// A paced game, end to end
// ---------------------------------------------------------------------------

// TestPacedGameCompletes runs the Phase 3 harness with the Phase 4 pacer
// switched on -- personas, hesitation and chat included -- and asserts that
// pacing changes nothing except the clock. The profile is scaled down hard
// (Speed=0.002) so the test runs in seconds while keeping the distribution's
// shape: this is the same dial an operator uses to make a headless sim fast.
func TestPacedGameCompletes(t *testing.T) {
	logger := zaptest.NewLogger(t)
	seats := fourSeats()
	sg := newSimGame(t, logger, "table-paced", seats)
	defer sg.close()

	sink := &fakeSink{}
	runner := NewBotRunner(RunnerConfig{
		GameID:   sg.g.ID,
		Actions:  sg.mgr,
		Views:    sg.adapter,
		Oracle:   simOracle(),
		Chat:     ChatConfig{Sink: sink, RoomID: "table-paced"},
		MaxTurns: 200,
		Logger:   logger,
	})

	profiles := []PacingProfile{AveragePace(), BriskPace(), DeliberatePace(), AveragePace()}
	for i, s := range seats {
		p := profiles[i]
		p.Speed *= 0.002
		p.HesitateChance = 0.25 // exercised harder than production
		h := NewHumanPacer(p, int64(1000+i))
		runner.AddSeatWithPacing(s, NewRandomPolicy(int64(i+1)),
			Pacing{Human: h, Chat: NewCannedChat(h)})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	start := time.Now()
	done, err := runner.Run(ctx)
	elapsed := time.Since(start)
	stats := runner.Stats()

	t.Logf("paced game: completed=%v turns=%d macros=%d failed=%d steps=%d chat=%d in %v",
		done, stats.Turns, stats.MacrosExecuted, stats.MacrosFailed, stats.StepsSent,
		sink.count(), elapsed.Round(time.Millisecond))

	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !done {
		t.Fatal("paced game did not reach a terminal state")
	}
	if stats.MacrosFailed != 0 {
		t.Fatalf("%d macro(s) failed under pacing: %v", stats.MacrosFailed, stats.FailedMacros)
	}
	if sink.count() == 0 {
		t.Error("no chat messages were sent during a full paced game")
	}
	if elapsed == 0 {
		t.Error("paced game took no measurable time; pacing did not apply")
	}
}
