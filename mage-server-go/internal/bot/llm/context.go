package llm

import (
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

// context.go manages the conversation. It is its own component rather than a
// detail of the loop because two INDEPENDENT constraints collide in it, and
// solving either one alone leaves a live bug:
//
//  1. Context size. A game runs for hundreds of decisions; sending the whole
//     transcript every turn is quadratic. mage-bench windows it
//     (reference/pilot_rendering.py:275-343): the last 40 messages verbatim,
//     the 20 before that with big tool results summarised, and a synthetic
//     "state bridge" user message between the two halves carrying a board
//     summary refreshed every 5 renders, so the model is never asked to reason
//     from a hole in its own history.
//
//  2. The 20-block cache lookback. A cache_control breakpoint searches BACKWARDS
//     AT MOST 20 CONTENT BLOCKS for a previously cached prefix. An agentic turn
//     with several tool_use/tool_result pairs blows past that, at which point
//     the cache is silently rebuilt -- no error, no warning, roughly 10x the
//     cost. So breakpoints go in every ~15 blocks, not just at the top.
//
// A request may carry at most 4 breakpoints. Two are already spoken for: the
// last tool definition and the system prompt (client.go). That leaves exactly
// two for messages, which is why placement walks backwards from the newest
// message -- the recent end is where a hit is worth the most, and it is the
// region that actually changes between turns.
const (
	// ContextRecentCount is how many trailing messages survive verbatim.
	// reference/pilot_rendering.py:29.
	ContextRecentCount = 40
	// ContextSummaryCount is how many messages before those are kept in
	// summarised form. reference/pilot_rendering.py:30.
	ContextSummaryCount = 20
	// ToolSummaryTriggerChars is the tool-result length above which a message
	// in the summary band gets compressed. reference/pilot_rendering.py:31.
	ToolSummaryTriggerChars = 200
	// RenderInterval is how often the state-bridge board summary is refreshed.
	// reference/pilot_rendering.py:32.
	RenderInterval = 5

	// CacheBlockInterval is the spacing between message breakpoints. It is 15,
	// not 20, so that a turn that adds a few blocks between renders still lands
	// inside the 20-block lookback.
	CacheBlockInterval = 15
	// MaxCacheBreakpoints is the API limit on cache_control blocks per request.
	MaxCacheBreakpoints = 4
	// ReservedBreakpoints are the ones client.go spends: tools + system prompt.
	ReservedBreakpoints = 2
	// MaxMessageBreakpoints is what is left for the message list.
	MaxMessageBreakpoints = MaxCacheBreakpoints - ReservedBreakpoints
)

// bridgeMarker is the literal mage-bench scans for when locating its cache
// breakpoint (reference/pilot_rendering.py:35). We place breakpoints by block
// count instead, but the sentence is kept verbatim: it is also the instruction
// that tells the model the offered options are pre-filtered, and golden prompts
// diff against it.
const bridgeMarker = "All cards listed are playable right now."

// Conversation is the append-only history of one bot's game.
//
// Append-only is deliberate: Render derives a bounded window from it on every
// call, so the window can widen, narrow or re-summarise without ever losing the
// real transcript, and a context reset is an explicit truncation rather than an
// accident of windowing.
//
// It is NOT safe for concurrent use. One Conversation belongs to one seat's
// goroutine, exactly like the rest of internal/bot.
type Conversation struct {
	history []anthropic.BetaMessageParam
	renders int

	// stateSummary is the board digest carried by the bridge message. It is
	// refreshed every RenderInterval renders rather than every render, because
	// changing it invalidates every cache entry after it.
	stateSummary string
}

// NewConversation returns an empty conversation.
func NewConversation() *Conversation { return &Conversation{} }

// Len reports the number of messages in the full history.
func (c *Conversation) Len() int { return len(c.history) }

// History returns the raw append-only history.
func (c *Conversation) History() []anthropic.BetaMessageParam { return c.history }

// Append adds a message to the history.
func (c *Conversation) Append(m anthropic.BetaMessageParam) { c.history = append(c.history, m) }

// AppendUserText appends a user message with a single text block.
func (c *Conversation) AppendUserText(text string) {
	c.Append(anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(text)))
}

// AppendAssistant appends the model's own reply, converted to param form.
func (c *Conversation) AppendAssistant(m *anthropic.BetaMessage) {
	if m == nil {
		return
	}
	c.Append(m.ToParam())
}

// AppendAssistantTurn appends the model's own reply from a neutral Response.
//
// The Completer that produced the Response also rendered the assistant turn,
// because only it knows what the provider actually sent -- an SDK
// BetaMessage.ToParam() on the Anthropic path, a reconstruction from
// tool_calls on the OpenAI-compatible one.
func (c *Conversation) AppendAssistantTurn(r *Response) {
	if r == nil || len(r.assistant.Content) == 0 {
		return
	}
	c.Append(r.assistant)
}

// AppendToolResults appends one user message carrying every tool result of a
// turn. They must arrive in a single message, in the order the tool_use blocks
// appeared, or the API rejects the request.
func (c *Conversation) AppendToolResults(blocks []anthropic.BetaContentBlockParamUnion) {
	if len(blocks) == 0 {
		return
	}
	c.Append(anthropic.NewBetaUserMessage(blocks...))
}

// NeedsStateRefresh reports whether the next render should be handed a fresh
// board summary. reference/pilot_rendering.py:32 (RENDER_INTERVAL).
func (c *Conversation) NeedsStateRefresh() bool {
	return c.stateSummary == "" || c.renders%RenderInterval == 0
}

// SetStateSummary updates the bridge message's board digest.
func (c *Conversation) SetStateSummary(s string) { c.stateSummary = s }

// Reset drops the transcript and starts again from one user message.
//
// This is the recovery path (reference/pilot_state.py::reset_context): after
// repeated truncations or timeouts the transcript itself is the suspect, and a
// model reasoning from 60 messages of failure keeps failing. The system prompt
// and the tool list are NOT touched -- they are the cached prefix, and
// rebuilding them would turn a recovery into a 10x cost event.
func (c *Conversation) Reset(text string) {
	c.history = nil
	c.stateSummary = ""
	c.renders = 0
	if text != "" {
		c.AppendUserText(text)
	}
}

// Render produces the bounded, CACHE-ANNOTATED message list for one request.
// It is the Anthropic strategy: explicit cache_control breakpoints, placed by
// the arithmetic below.
//
// It never mutates the history: the returned messages are copies down to the
// content blocks it annotates, so the breakpoints of one request cannot leak
// into the next.
func (c *Conversation) Render() []anthropic.BetaMessageParam {
	out := c.RenderUncached()
	placeCacheBreakpoints(out, MaxMessageBreakpoints, CacheBlockInterval)
	return out
}

// RenderUncached produces the same bounded window with NO cache annotations.
//
// It is what the OpenAI-compatible providers use. That is not a missing
// feature: Gemini's caching is implicit and enabled by default, has no request
// field to set, and asks only that the large common content stay at the front
// of the prompt and that requests with a shared prefix arrive close together
// (ai.google.dev/gemini-api/docs/caching). The windowing, the summary band and
// the state bridge all still apply, because those are about context size and
// prefix stability, which every provider has. Only the breakpoint arithmetic is
// Anthropic's.
//
// The render counter advances here rather than in Render so the state-bridge
// refresh interval is identical on both providers.
func (c *Conversation) RenderUncached() []anthropic.BetaMessageParam {
	c.renders++
	return c.window()
}

// window applies the recent/summary/bridge split.
func (c *Conversation) window() []anthropic.BetaMessageParam {
	if len(c.history) <= ContextRecentCount {
		return copyMessages(c.history)
	}

	recentStart := len(c.history) - ContextRecentCount
	// A message that opens with a tool_result cannot start a window: its
	// tool_use is in the message before it, and the API rejects an orphan.
	// reference/pilot_rendering.py:310-311 does the same walk.
	for recentStart > 0 && startsWithToolResult(c.history[recentStart]) {
		recentStart--
	}

	summaryStart := recentStart - ContextSummaryCount
	if summaryStart < 0 {
		summaryStart = 0
	}
	for summaryStart > 0 && startsWithToolResult(c.history[summaryStart]) {
		summaryStart--
	}

	out := make([]anthropic.BetaMessageParam, 0, ContextRecentCount+ContextSummaryCount+1)
	for i := summaryStart; i < recentStart; i++ {
		out = append(out, summarizeMessage(c.history, i))
	}
	out = append(out, anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(c.bridgeText())))
	out = append(out, copyMessages(c.history[recentStart:])...)
	return out
}

// bridgeText is the synthetic message that spans the summarised gap.
// reference/pilot_rendering.py:325-333.
func (c *Conversation) bridgeText() string {
	var b strings.Builder
	if c.stateSummary != "" {
		b.WriteString(c.stateSummary)
		if !strings.HasSuffix(c.stateSummary, "\n") {
			b.WriteString("\n")
		}
	}
	b.WriteString("Continue playing. The decision below is current; earlier tool results " +
		"above have been summarised. ")
	b.WriteString(bridgeMarker)
	b.WriteString(" Pick one option with choose_action(choice=\"mN\"), or pass_priority.")
	return b.String()
}

// summarizeMessage compresses the oversized tool results of one message.
func summarizeMessage(history []anthropic.BetaMessageParam, idx int) anthropic.BetaMessageParam {
	msg := copyMessage(history[idx])
	for i := range msg.Content {
		tr := msg.Content[i].OfToolResult
		if tr == nil {
			continue
		}
		text := toolResultText(tr)
		if len(text) <= ToolSummaryTriggerChars {
			continue
		}
		name := findToolName(history, idx, tr.ToolUseID)
		tr.Content = []anthropic.BetaToolResultBlockParamContentUnion{
			{OfText: &anthropic.BetaTextBlockParam{Text: summarizeToolResult(name, text)}},
		}
	}
	return msg
}

// summarizeToolResult is the Go analogue of
// reference/pilot_rendering.py::_summarize_tool_result.
//
// Upstream parses each tool's JSON result and rebuilds a field-aware digest.
// Ours are plain text -- a rendered decision block, an oracle entry, an ack --
// so the digest is structural instead: the first line, which for a decision
// block is the "[Decision N, snapshot=M] Turn ..." header and carries the turn,
// phase and deciding seat, plus a count of what was dropped. Enough for the
// model to know what happened without re-reading a board that is six turns stale.
func summarizeToolResult(name, text string) string {
	trimmed := strings.TrimSpace(text)
	first := trimmed
	if i := strings.IndexByte(trimmed, '\n'); i >= 0 {
		first = trimmed[:i]
	}
	if len(first) > ToolSummaryTriggerChars {
		first = first[:ToolSummaryTriggerChars]
	}
	label := name
	if label == "" {
		label = "tool"
	}
	dropped := len(trimmed) - len(first)
	if dropped <= 0 {
		return first
	}
	return fmt.Sprintf("%s: %s (+%d chars elided)", label, first, dropped)
}

// findToolName walks back for the assistant tool_use that a tool_result answers.
// reference/pilot_rendering.py:194-215.
func findToolName(history []anthropic.BetaMessageParam, resultIdx int, toolUseID string) string {
	for i := resultIdx - 1; i >= 0; i-- {
		if history[i].Role != anthropic.BetaMessageParamRoleAssistant {
			continue
		}
		for _, b := range history[i].Content {
			if b.OfToolUse != nil && b.OfToolUse.ID == toolUseID {
				return b.OfToolUse.Name
			}
		}
		break
	}
	return ""
}

func toolResultText(tr *anthropic.BetaToolResultBlockParam) string {
	var b strings.Builder
	for _, c := range tr.Content {
		if c.OfText != nil {
			b.WriteString(c.OfText.Text)
		}
	}
	return b.String()
}

func startsWithToolResult(m anthropic.BetaMessageParam) bool {
	return len(m.Content) > 0 && m.Content[0].OfToolResult != nil
}

// ---------------------------------------------------------------------------
// Cache breakpoints
// ---------------------------------------------------------------------------

// placeCacheBreakpoints annotates up to budget messages, at most every
// `interval` content blocks, walking backwards from the newest message.
//
// Walking backwards is what makes the chain hold turn to turn. Each request
// adds a handful of blocks at the end; a breakpoint pinned to a fixed index
// would drift out of the previous entry's 20-block lookback within two turns,
// whereas one pinned `interval` blocks from the end moves with the transcript
// and stays inside it.
func placeCacheBreakpoints(msgs []anthropic.BetaMessageParam, budget, interval int) {
	if len(msgs) == 0 || budget <= 0 {
		return
	}
	placed := 0
	since := 0
	for i := len(msgs) - 1; i >= 0 && placed < budget; i-- {
		blocks := len(msgs[i].Content)
		if placed == 0 || since >= interval {
			if setCacheControl(&msgs[i]) {
				placed++
				// The message's own blocks count toward the NEXT gap, not the
				// one just closed: the mark lands on its last block, so
				// everything before it in the message sits on the far side.
				since = blocks
				continue
			}
		}
		since += blocks
	}
}

// setCacheControl marks the last cacheable block of a message. It reports
// whether it found one -- a message of nothing but thinking blocks has no
// breakpoint to hang, and the caller must keep looking.
func setCacheControl(m *anthropic.BetaMessageParam) bool {
	for i := len(m.Content) - 1; i >= 0; i-- {
		switch {
		case m.Content[i].OfText != nil:
			m.Content[i].OfText.CacheControl = anthropic.NewBetaCacheControlEphemeralParam()
			return true
		case m.Content[i].OfToolResult != nil:
			m.Content[i].OfToolResult.CacheControl = anthropic.NewBetaCacheControlEphemeralParam()
			return true
		case m.Content[i].OfToolUse != nil:
			m.Content[i].OfToolUse.CacheControl = anthropic.NewBetaCacheControlEphemeralParam()
			return true
		}
	}
	return false
}

// BreakpointBlockIndexes returns the global content-block index of every
// cache_control breakpoint in msgs, ascending.
//
// It exists so the placement rule is testable as a property -- "no two
// breakpoints more than 20 blocks apart, never more than 4" -- rather than by
// asserting on hard-coded indexes that any change to the window would break.
func BreakpointBlockIndexes(msgs []anthropic.BetaMessageParam) []int {
	var out []int
	block := 0
	for _, m := range msgs {
		for _, c := range m.Content {
			marked := false
			switch {
			case c.OfText != nil:
				marked = !param.IsOmitted(c.OfText.CacheControl)
			case c.OfToolResult != nil:
				marked = !param.IsOmitted(c.OfToolResult.CacheControl)
			case c.OfToolUse != nil:
				marked = !param.IsOmitted(c.OfToolUse.CacheControl)
			}
			if marked {
				out = append(out, block)
			}
			block++
		}
	}
	return out
}

// copyMessages / copyMessage deep-copy far enough that annotating the copy
// cannot reach the history. The block unions hold pointers, so copying the
// slice alone would share every block.
func copyMessages(in []anthropic.BetaMessageParam) []anthropic.BetaMessageParam {
	out := make([]anthropic.BetaMessageParam, 0, len(in))
	for _, m := range in {
		out = append(out, copyMessage(m))
	}
	return out
}

func copyMessage(m anthropic.BetaMessageParam) anthropic.BetaMessageParam {
	dup := m
	dup.Content = make([]anthropic.BetaContentBlockParamUnion, len(m.Content))
	for i, c := range m.Content {
		blk := c
		switch {
		case c.OfText != nil:
			t := *c.OfText
			blk.OfText = &t
		case c.OfToolResult != nil:
			tr := *c.OfToolResult
			tr.Content = append([]anthropic.BetaToolResultBlockParamContentUnion(nil), c.OfToolResult.Content...)
			for j, cc := range tr.Content {
				if cc.OfText != nil {
					t := *cc.OfText
					tr.Content[j].OfText = &t
				}
			}
			blk.OfToolResult = &tr
		case c.OfToolUse != nil:
			tu := *c.OfToolUse
			blk.OfToolUse = &tu
		}
		dup.Content[i] = blk
	}
	return dup
}
