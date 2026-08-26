package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/magefree/mage-server-go/internal/bot"
	"go.uber.org/zap"
)

// policy.go is the agentic loop, expressed as a bot.Policy.
//
// TWO TIMERS, AND THEY ARE NOT THE SAME TIMER (§0.4). Conflating them is the
// main design trap in this phase:
//
//   - The LLM request timeout is 120s, PER ATTEMPT. mage-bench raised it from
//     45s at harness epoch #15 because 45s was empirically too short; a request
//     that has been running 50s is usually about to succeed. It is applied via
//     option.WithRequestTimeout, which the SDK re-applies to every retry -- so
//     the worst case for one request is 120s x (retries+1).
//   - The table-stall guard is 45-60s of held priority, TOTAL. It exists for
//     the other three players, not for the bot: past it, the seat force-passes
//     regardless of what the model is doing. It is the context.WithTimeout that
//     bounds the whole Pick, which is what stops the per-attempt timeout from
//     multiplying out to six minutes of dead air.
//
// The stall guard normally fires first, and that is correct: a decision worth
// more than a minute of everyone else's time is not worth it.
//
// PACING OVERLAPS LATENCY, IT DOES NOT FOLLOW IT. Phase 4 gives each seat a
// human-like pre-action pause. Thinking time is already a pause, so serialising
// the two would make an LLM bot visibly slower than a human rather than
// human-like. Pick fires the request immediately and reports how long it took
// via LastThinkTime; the runner subtracts that from the pause it wanted and
// waits only for the remainder. Latency becomes visible only when it exceeds
// the pause we were going to take anyway.
type LLMPolicy struct {
	client *Client
	conv   *Conversation
	ser    *bot.Serializer
	rec    *Recovery
	opts   Options

	decisions int

	// chat holds lines the model produced through send_chat_message, drained by
	// Line. It is a queue rather than a single slot because a model can call
	// the tool twice in one turn, and the runner's cadence rules may only let
	// one out.
	chat []string

	// wasReset records that the last recovery step wiped the transcript, so
	// the loop knows to re-send the decision the model can no longer see.
	wasReset bool

	// mu guards chat and think, which Line and LastThinkTime read from the
	// runner's goroutine. Everything else is touched only inside Pick.
	mu    sync.Mutex
	think time.Duration
}

// Options configures an LLMPolicy.
type Options struct {
	// Seat is the bot's player ID, for logs.
	Seat string

	// Model defaults to claude-sonnet-5.
	//
	// DO NOT ASSUME HAIKU IS CHEAPER (§0.7). Its minimum cacheable prefix is
	// 4096 tokens against Sonnet 5's 1024, so a prompt that Sonnet caches from
	// turn one may never cache at all on Haiku -- and an uncached Haiku prompt
	// beats a cached Sonnet one only on sticker price. It also supports neither
	// adaptive thinking nor OutputConfig.Effort, which removes the main dial
	// for trading quality against cost. Phase 6 measures both and lets the data
	// decide.
	Model anthropic.Model

	// MaxTokens defaults to DefaultMaxTokens (upstream's 20k).
	MaxTokens int64

	// Effort is ignored on Haiku 4.5, which does not offer it. normalise()
	// clears it rather than letting a 400 teach the lesson at 3am.
	Effort anthropic.BetaOutputConfigEffort

	// ThinkingBudget sets ThinkingConfigParamOfEnabled(N). It is the Haiku
	// path; on Sonnet 5 / Opus 5 a budget is a 400, and adaptive thinking is
	// used instead.
	ThinkingBudget int64

	// RequestTimeout is per attempt. Defaults to DefaultRequestTimeout (120s).
	RequestTimeout time.Duration

	// StallTimeout bounds one whole Pick, retries included. Defaults to
	// DefaultStallTimeout.
	StallTimeout time.Duration

	// MaxRetries is passed to the SDK. Defaults to DefaultMaxRetries.
	MaxRetries int

	// MaxSteps bounds the tool-calling loop within one decision, so a model
	// that keeps asking for oracle text never becomes an unbounded spend.
	MaxSteps int

	// Oracle supplies printed characteristics for the Card Reference and for
	// get_oracle_text. May be nil.
	Oracle bot.OracleLookup

	// System overrides the system prompt. Leave empty; it exists for golden
	// tests, not for per-seat personality.
	System string

	Logger *zap.Logger
}

// Defaults for the loop itself.
const (
	// DefaultStallTimeout is the table-stall guard: the wall clock a bot may
	// hold priority before it force-passes. The plan's range is 45-60s; 50s
	// sits inside it with room for the fallback macro to be sent.
	DefaultStallTimeout = 50 * time.Second
	// DefaultMaxSteps is the tool-call budget for one decision. A decision
	// needs one call; the headroom is for oracle lookups, a chat line and one
	// correction after a malformed call.
	DefaultMaxSteps = 6
	// DefaultHaikuThinkingBudget is used when Haiku is selected without an
	// explicit budget. It must be below MaxTokens and at least 1024.
	DefaultHaikuThinkingBudget = 4096
)

func (o *Options) normalise() {
	if o.Model == "" {
		o.Model = anthropic.ModelClaudeSonnet5
	}
	if o.MaxTokens <= 0 {
		o.MaxTokens = DefaultMaxTokens
	}
	if o.RequestTimeout <= 0 {
		o.RequestTimeout = DefaultRequestTimeout
	}
	if o.StallTimeout <= 0 {
		o.StallTimeout = DefaultStallTimeout
	}
	if o.MaxSteps <= 0 {
		o.MaxSteps = DefaultMaxSteps
	}
	if o.MaxRetries < 0 {
		o.MaxRetries = DefaultMaxRetries
	}
	if o.Logger == nil {
		o.Logger = zap.NewNop()
	}
	if IsHaiku(o.Model) {
		// Haiku 4.5 supports NEITHER adaptive thinking NOR effort (§0.7).
		// Clearing Effort here rather than trusting the caller is deliberate:
		// the failure is a 400 on the first request of a long unattended run.
		o.Effort = ""
		if o.ThinkingBudget <= 0 {
			o.ThinkingBudget = DefaultHaikuThinkingBudget
		}
		if o.ThinkingBudget >= o.MaxTokens {
			o.ThinkingBudget = o.MaxTokens / 2
		}
	} else {
		// budget_tokens on Sonnet 5 / Opus 5 returns 400 (plan, Phase 5
		// anti-pattern guards). Adaptive thinking is the supported form.
		o.ThinkingBudget = 0
	}
}

// IsHaiku reports whether m is a Haiku model, which is the model family with
// the different thinking and effort rules.
func IsHaiku(m anthropic.Model) bool {
	return strings.Contains(strings.ToLower(string(m)), "haiku")
}

// New builds an LLMPolicy over t.
func New(t Transport, o Options) *LLMPolicy {
	o.normalise()
	return &LLMPolicy{
		client: NewClient(t, ClientOptions{
			Model:          o.Model,
			MaxTokens:      o.MaxTokens,
			Effort:         o.Effort,
			ThinkingBudget: o.ThinkingBudget,
			Adaptive:       !IsHaiku(o.Model),
			RequestTimeout: o.RequestTimeout,
			MaxRetries:     o.MaxRetries,
			System:         o.System,
		}),
		conv: NewConversation(),
		ser:  bot.NewSerializer(o.Oracle),
		rec:  NewRecovery(),
		opts: o,
	}
}

// Interface assertions. LLMPolicy is a drop-in for RandomPolicy: the runner,
// the harness and every Phase 3 test treat them identically.
var (
	_ bot.Policy            = (*LLMPolicy)(nil)
	_ bot.ChatSource        = (*LLMPolicy)(nil)
	_ bot.ThinkTimeReporter = (*LLMPolicy)(nil)
)

// Usage returns the accumulated token counters.
func (p *LLMPolicy) Usage() Usage { return p.client.Usage() }

// LastThinkTime implements bot.ThinkTimeReporter: how long the last Pick spent
// waiting on the model, which the runner deducts from its pacing delay.
func (p *LLMPolicy) LastThinkTime() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.think
}

// Line implements bot.ChatSource, replacing Phase 4's CannedChat with the
// model's own words.
//
// It only ever drains what the model already said; it never blocks and never
// makes a request of its own. Cadence -- the two-lines-per-turn cap and the
// once-per-two-cycles floor -- stays in the runner where it was, so a chatty
// model cannot spam the table.
func (p *LLMPolicy) Line(_ context.Context, _ *bot.SafeView, _ bot.Macro, _ bool) (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.chat) == 0 {
		return "", false
	}
	line := p.chat[0]
	p.chat = p.chat[1:]
	return line, true
}

func (p *LLMPolicy) queueChat(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	p.mu.Lock()
	p.chat = append(p.chat, line)
	p.mu.Unlock()
}

// Pick implements bot.Policy.
//
// It always returns a legal macro when one exists. Every error path -- timeout,
// outage, a model that will not call a tool -- ends in the fallback macro
// rather than an error, because an error would stop the seat and stall the
// table for the other three players.
func (p *LLMPolicy) Pick(ctx context.Context, v *bot.SafeView, moves []bot.Macro) (bot.Macro, error) {
	if len(moves) == 0 {
		return bot.Macro{}, bot.ErrNoMoves
	}
	start := time.Now()
	defer func() {
		p.mu.Lock()
		p.think = time.Since(start)
		p.mu.Unlock()
	}()

	if p.rec.Autopilot() {
		return fallbackMacro(moves), nil
	}

	// The stall guard. Everything below -- every attempt, every retry -- lives
	// inside it.
	ctx, cancel := context.WithTimeout(ctx, p.opts.StallTimeout)
	defer cancel()

	p.decisions++
	if p.conv.NeedsStateRefresh() {
		p.conv.SetStateSummary(boardSummary(v))
	}
	p.conv.AppendUserText(decisionMessage(ctx, p.ser, v, moves, p.decisions, nil))

	malformed := 0
	for step := 0; step < p.opts.MaxSteps; step++ {
		msg, err := p.client.Complete(ctx, p.conv)
		if err != nil {
			if m, done := p.handleError(err, moves); done {
				return m, nil
			}
			p.redecideIfReset(ctx, v, moves)
			continue
		}
		p.conv.AppendAssistant(msg)

		uses := ToolUses(msg)
		if len(uses) == 0 {
			kind := FailureEmpty
			if isTruncated(msg) {
				kind = FailureTruncated
			}
			if m, done := p.handleSoftFailure(kind, moves); done {
				return m, nil
			}
			p.redecideIfReset(ctx, v, moves)
			continue
		}

		picked, ok, attempted, results := p.handleToolUses(ctx, v, moves, uses)
		p.conv.AppendToolResults(results)
		if ok {
			p.rec.Success()
			return moves[picked], nil
		}
		if !attempted {
			// Info-only tools -- oracle lookups, a chat line -- are not a
			// failure to decide. Upstream draws the same line
			// (INFO_ONLY_TOOLS, Sec 0.4). Counting a chatty turn as malformed
			// would degrade a working bot to autopilot for talking.
			continue
		}
		malformed++
		if malformed >= MaxMalformedPerDecision {
			p.log("malformed tool output, force-passing", zap.Int("attempts", malformed))
			return p.forcePass(moves), nil
		}
	}

	p.log("step budget exhausted, force-passing", zap.Int("steps", p.opts.MaxSteps))
	return p.forcePass(moves), nil
}

// redecideIfReset re-sends the current decision after a context reset.
//
// Without it a reset would leave the model with a transcript that no longer
// contains the board or the options it is being asked about, and the very next
// request would be answered blind -- the failure mode a reset exists to avoid.
func (p *LLMPolicy) redecideIfReset(ctx context.Context, v *bot.SafeView, moves []bot.Macro) {
	if !p.wasReset {
		return
	}
	p.wasReset = false
	p.conv.SetStateSummary(boardSummary(v))
	p.conv.AppendUserText(decisionMessage(ctx, p.ser, v, moves, p.decisions, nil))
}

// handleError applies the recovery matrix to a transport error. It reports the
// macro to return and whether the decision is over.
func (p *LLMPolicy) handleError(err error, moves []bot.Macro) (bot.Macro, bool) {
	kind := Classify(err)
	if kind == FailurePermanent {
		p.log("permanent LLM failure",
			zap.String("reason", PermanentFailureReason(err)), zap.Error(err))
	} else {
		p.log("LLM request failed", zap.Error(err))
	}
	return p.applyAction(p.rec.Record(kind), kind, moves)
}

// handleSoftFailure applies the matrix to a response that arrived but was not
// usable: empty, or truncated on the output cap.
func (p *LLMPolicy) handleSoftFailure(kind FailureKind, moves []bot.Macro) (bot.Macro, bool) {
	return p.applyAction(p.rec.Record(kind), kind, moves)
}

func (p *LLMPolicy) applyAction(action Action, kind FailureKind, moves []bot.Macro) (bot.Macro, bool) {
	switch action {
	case ActionAutopilot:
		p.enterAutopilot()
		return fallbackMacro(moves), true
	case ActionForcePass:
		return p.forcePass(moves), true
	case ActionResetContext:
		p.log("resetting conversation context", zap.Int("failure", int(kind)))
		p.wasReset = true
		p.conv.Reset("Your context was reset. Continue playing: pick one of the options in the next message.")
		// The Card Reference dedup has to reset with it -- the model can no
		// longer see the oracle text it was printed under the old transcript.
		p.ser.ResetOracleDedup()
		return bot.Macro{}, false
	default:
		// Retry within this decision. Nudge, so the next request is not a
		// byte-identical repeat of the one that just failed.
		p.conv.AppendUserText("That did not produce an action. Call choose_action with one of the listed option ids, or pass_priority.")
		return bot.Macro{}, false
	}
}

// forcePass is the "keep the table moving" exit: say something in character,
// then take the fallback macro. It is what the other three players see instead
// of a seat that stopped responding.
func (p *LLMPolicy) forcePass(moves []bot.Macro) bot.Macro {
	if p.rec.ForcedPass() {
		p.enterAutopilot()
	} else {
		p.queueChat(StallLine)
	}
	return fallbackMacro(moves)
}

func (p *LLMPolicy) enterAutopilot() {
	p.queueChat(DegradationLine)
	p.log("degraded to autopilot", zap.String("reason", p.rec.AutopilotReason()))
}

func (p *LLMPolicy) log(msg string, fields ...zap.Field) {
	p.opts.Logger.Warn("bot/llm: "+msg,
		append([]zap.Field{zap.String("seat", p.opts.Seat)}, fields...)...)
}

// ---------------------------------------------------------------------------
// Tool dispatch
// ---------------------------------------------------------------------------

type chooseActionInput struct {
	Choice  string `json:"choice"`
	Amount  *int   `json:"amount"`
	Amounts []int  `json:"amounts"`
	Pile    *int   `json:"pile"`
	Text    string `json:"text"`
}

type oracleInput struct {
	CardName  string   `json:"card_name"`
	CardNames []string `json:"card_names"`
	ObjectID  string   `json:"object_id"`
	ObjectIDs []string `json:"object_ids"`
}

type chatInput struct {
	Message string `json:"message"`
}

// handleToolUses runs every tool call of one assistant turn and returns the
// chosen macro index, whether one was chosen, and the tool_result blocks to
// append.
//
// EVERY tool_use gets a result, including the ones that arrive after an action
// has already been chosen. The API rejects a turn with an unanswered tool_use,
// so a missing result would not degrade the next request -- it would fail it,
// for the rest of the game.
func (p *LLMPolicy) handleToolUses(ctx context.Context, v *bot.SafeView, moves []bot.Macro, uses []ToolUse) (picked int, chosen, attempted bool, results []anthropic.BetaContentBlockParamUnion) {
	picked = -1
	results = make([]anthropic.BetaContentBlockParamUnion, 0, len(uses))

	for _, u := range uses {
		switch u.Name {
		case ToolSendChatMessage:
			var in chatInput
			if err := json.Unmarshal(u.Input, &in); err != nil || strings.TrimSpace(in.Message) == "" {
				results = append(results, anthropic.NewBetaToolResultBlock(u.ID, "error: message is required", true))
				continue
			}
			p.queueChat(in.Message)
			results = append(results, anthropic.NewBetaToolResultBlock(u.ID, `{"success": true}`, false))

		case ToolGetOracleText:
			var in oracleInput
			if err := json.Unmarshal(u.Input, &in); err != nil {
				results = append(results, anthropic.NewBetaToolResultBlock(u.ID, "error: unreadable input", true))
				continue
			}
			results = append(results, anthropic.NewBetaToolResultBlock(u.ID, p.oracleText(ctx, v, in), false))

		case ToolChooseAction, ToolPassPriority:
			attempted = true
			if chosen {
				results = append(results, anthropic.NewBetaToolResultBlock(u.ID,
					"error: one action per decision; the earlier call was taken", true))
				continue
			}
			idx, ok, msg := p.resolveAction(u, moves)
			if !ok {
				results = append(results, anthropic.NewBetaToolResultBlock(u.ID, msg, true))
				continue
			}
			picked, chosen = idx, true
			results = append(results, anthropic.NewBetaToolResultBlock(u.ID, msg, false))

		default:
			// An unknown tool is a schema problem, not a chat turn: it must
			// count against this decision or the loop spins on it.
			attempted = true
			results = append(results, anthropic.NewBetaToolResultBlock(u.ID,
				"error: unknown tool "+u.Name, true))
		}
	}
	return picked, chosen, attempted, results
}

// resolveAction maps one action tool call onto a macro index.
func (p *LLMPolicy) resolveAction(u ToolUse, moves []bot.Macro) (int, bool, string) {
	if u.Name == ToolPassPriority {
		idx := fallbackIndex(moves)
		return idx, true, fmt.Sprintf("{\"success\": true, \"action_taken\": %q}", moves[idx].Label)
	}
	var in chooseActionInput
	if err := json.Unmarshal(u.Input, &in); err != nil {
		return 0, false, "error: unreadable input; call choose_action with choice=\"mN\""
	}
	idx, ok := resolveChoice(in.Choice, moves)
	if !ok {
		return 0, false, fmt.Sprintf(
			"error: %q is not one of the %d listed options; use choice=\"m1\"..%q, or yes/no",
			in.Choice, len(moves), macroID(len(moves)-1))
	}
	return idx, true, fmt.Sprintf("{\"success\": true, \"action_taken\": %q}", moves[idx].Label)
}

// resolveChoice parses the model's `choice` against the offered macros.
//
// Three accepted forms, in the order they are tried:
//
//	"m3"  -- an option id, 1-based, as printed in the Choices line.
//	"2"   -- a bare index, ZERO-based, which is mage-bench's index semantics
//	         (§0.3: "parseable int -> index"). The two are deliberately
//	         different bases because they are deliberately different syntaxes;
//	         conflating them would silently off-by-one every bare index.
//	"yes" -- yes|true means mulligan/confirm, no|false means keep/pass, exactly
//	         as the vendored system prompt documents.
func resolveChoice(choice string, moves []bot.Macro) (int, bool) {
	c := strings.ToLower(strings.TrimSpace(choice))
	if c == "" {
		return 0, false
	}
	switch c {
	case "yes", "true":
		if i, ok := indexOfKind(moves, bot.KindMulligan); ok {
			return i, true
		}
		return 0, false
	case "no", "false", "pass":
		if i, ok := indexOfKind(moves, bot.KindKeepHand); ok {
			return i, true
		}
		if i, ok := indexOfKind(moves, bot.KindPassTurn); ok {
			return i, true
		}
		return 0, false
	}
	if strings.HasPrefix(c, "m") {
		if n, err := strconv.Atoi(c[1:]); err == nil && n >= 1 && n <= len(moves) {
			return n - 1, true
		}
		return 0, false
	}
	if n, err := strconv.Atoi(c); err == nil && n >= 0 && n < len(moves) {
		return n, true
	}
	return 0, false
}

func indexOfKind(moves []bot.Macro, k bot.Kind) (int, bool) {
	for i, m := range moves {
		if m.KindOf() == k {
			return i, true
		}
	}
	return 0, false
}

// fallbackIndex is the move to take when the model cannot or will not choose.
//
// Passing the turn is the least-committal legal action and the one that keeps
// the game moving; during the mulligan phase there is no turn to pass, so
// keeping the hand is the equivalent no-op. Only if neither exists does it fall
// through to the first offered macro, which cannot be worse than doing nothing
// -- doing nothing hangs the table.
func fallbackIndex(moves []bot.Macro) int {
	if i, ok := indexOfKind(moves, bot.KindPassTurn); ok {
		return i
	}
	if i, ok := indexOfKind(moves, bot.KindKeepHand); ok {
		return i
	}
	return 0
}

func fallbackMacro(moves []bot.Macro) bot.Macro { return moves[fallbackIndex(moves)] }

// oracleText answers get_oracle_text from the same OracleLookup the serializer
// uses, so an answer here can never contradict the Card Reference.
func (p *LLMPolicy) oracleText(ctx context.Context, v *bot.SafeView, in oracleInput) string {
	names := append([]string(nil), in.CardNames...)
	if in.CardName != "" {
		names = append(names, in.CardName)
	}
	ids := append([]string(nil), in.ObjectIDs...)
	if in.ObjectID != "" {
		ids = append(ids, in.ObjectID)
	}
	for _, id := range ids {
		if c := p.cardByShortID(v, id); c != nil {
			names = append(names, c.Name)
		}
	}
	if len(names) == 0 {
		return "error: give card_name/card_names or object_id/object_ids"
	}
	if p.opts.Oracle == nil {
		return "no oracle data is available at this table"
	}
	var b strings.Builder
	for _, n := range names {
		card, ok := p.opts.Oracle.Oracle(ctx, n)
		if !ok || card == nil {
			fmt.Fprintf(&b, "%s: not found\n", n)
			continue
		}
		fmt.Fprintf(&b, "%s %s -- %s", card.Name, card.ManaCost, card.TypeLine)
		if card.Power != "" || card.Toughness != "" {
			fmt.Fprintf(&b, " %s/%s", card.Power, card.Toughness)
		}
		fmt.Fprintf(&b, ": %s\n", strings.ReplaceAll(card.OracleText, "\n", " / "))
	}
	return strings.TrimRight(b.String(), "\n")
}

// cardByShortID resolves a "p3" back to a card through the serializer's
// registry, which is the same allocation the model was shown.
func (p *LLMPolicy) cardByShortID(v *bot.SafeView, short string) *bot.SafeCard {
	id, ok := p.ser.IDs().TryResolve(short)
	if !ok || v == nil {
		return nil
	}
	for _, zone := range [][]*bot.SafeCard{v.Battlefield, v.Exile, v.Stack, v.Command} {
		for _, c := range zone {
			if c != nil && c.ID == id {
				return c
			}
		}
	}
	if v.Me != nil {
		for _, c := range append(append([]*bot.SafeCard(nil), v.Me.Hand...), v.Me.Graveyard...) {
			if c != nil && c.ID == id {
				return c
			}
		}
	}
	return nil
}
