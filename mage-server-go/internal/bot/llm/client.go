// Package llm implements the Phase 5 LLM-backed bot policy.
//
// It sits behind internal/bot's Policy interface: the runner, the move
// generator, the redactor and the serializer are all unchanged, and the model
// only ever picks one of the macros LegalMoves already produced. It never emits
// an engine command (plan Phase 5, anti-pattern guard 4): Go expands the macro.
//
// THE LOOP IS MANUAL, ON PURPOSE (anti-pattern 1). The SDK ships a tool-runner
// helper for the beta messages API, and it is the wrong tool here: it rebuilds
// params.Tools from
// Name()/Description()/InputSchema() on every turn, which makes CacheControl
// and Strict unsettable on a runner-managed tool. Cached tool definitions are
// exactly what a game-long conversation needs -- tools render at position 0, so
// they are the largest permanently-cacheable prefix we have. The runner also
// fans handlers out through an errgroup and silently drops pending tool calls
// on a MaxTokens or Refusal stop. So: client.Beta.Messages.New, in a loop we
// own.
package llm

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// Tool names, ported from reference/mcp-tools.json5.
//
// mage-bench exposes seven; we port four. get_game_state and get_game_log are
// deliberately skipped: our SafeView is small enough to inline into the
// decision message, so a round trip to fetch it would only add latency and
// blocks. get_action_choices is redundant for the same reason -- the macro list
// is already in the prompt. concede/join_table are not exposed upstream either.
const (
	ToolChooseAction    = "choose_action"
	ToolGetOracleText   = "get_oracle_text"
	ToolPassPriority    = "pass_priority"
	ToolSendChatMessage = "send_chat_message"
)

// Transport is the one call this package makes against the Anthropic API.
//
// It exists so the loop can be tested without a network or a credential: every
// test in this package injects a fake. It is also the seam a future
// rate-limiter or request recorder would slot into. The signature is exactly
// (*anthropic.BetaMessageService).New (§0.7), so the real implementation is a
// one-line delegation and cannot drift from the SDK.
type Transport interface {
	New(ctx context.Context, params anthropic.BetaMessageNewParams, opts ...option.RequestOption) (*anthropic.BetaMessage, error)
}

// SDKTransport is the production Transport, backed by the Anthropic Go SDK.
type SDKTransport struct {
	client anthropic.Client
}

// NewSDKTransport builds a Transport over the real API.
//
// THE KEY IS PASSED IN, NEVER READ FROM CONFIG FILES (anti-pattern 4).
// config.BotConfig.APIKey is populated from the ANTHROPIC_API_KEY environment
// variable and from nowhere else; config/config.yaml is a hard link to
// config.dev.yaml and is checked in.
func NewSDKTransport(apiKey string, opts ...option.RequestOption) *SDKTransport {
	all := make([]option.RequestOption, 0, len(opts)+1)
	if apiKey != "" {
		all = append(all, option.WithAPIKey(apiKey))
	}
	all = append(all, opts...)
	return &SDKTransport{client: anthropic.NewClient(all...)}
}

// New implements Transport.
func (t *SDKTransport) New(ctx context.Context, params anthropic.BetaMessageNewParams, opts ...option.RequestOption) (*anthropic.BetaMessage, error) {
	return t.client.Beta.Messages.New(ctx, params, opts...)
}

// ---------------------------------------------------------------------------
// Tool definitions
// ---------------------------------------------------------------------------

// Tools returns the tool definitions, SORTED BY NAME.
//
// The sort is not cosmetic. Tools serialize at position 0 of the request, ahead
// of the system prompt and every message, so any reordering changes the cached
// prefix and invalidates the entire cache for the rest of the game -- silently,
// with no error, at roughly 10x the cost. Go map iteration is randomised per
// run, so anything that builds this list from a map must sort on the way out.
//
// CacheControl goes on the LAST definition only: one breakpoint covers every
// block before it, so marking each tool would burn four of the four
// breakpoints a request is allowed.
//
// Strict + additionalProperties:false mirror mage-bench, whose schemas are
// reflected from Java annotations and always emit additionalProperties:false.
// The SDK has no typed field for additionalProperties; it goes in
// BetaToolInputSchemaParam.ExtraFields (§0.7).
func Tools() []anthropic.BetaToolParam {
	tools := []anthropic.BetaToolParam{
		chooseActionTool(),
		getOracleTextTool(),
		passPriorityTool(),
		sendChatMessageTool(),
	}
	sort.Slice(tools, func(i, j int) bool { return tools[i].Name < tools[j].Name })
	last := len(tools) - 1
	tools[last].CacheControl = anthropic.NewBetaCacheControlEphemeralParam()
	return tools
}

func strictSchema(props map[string]any, required []string) anthropic.BetaToolInputSchemaParam {
	return anthropic.BetaToolInputSchemaParam{
		Properties:  props,
		Required:    required,
		ExtraFields: map[string]any{"additionalProperties": false},
	}
}

// chooseActionTool ports reference/mcp-tools.json5:1254 (choose_action).
//
// PORTED SUBSET, AND WHY. Upstream's choose_action takes nine optional params.
// We port the five that answer a decision -- choice, amount, amounts, pile,
// text -- and omit mana_plan, auto_tap, attackers and blockers. Those four
// address subsystems this engine does not have: there is no mana payment, no
// priority and no combat step (§0.5 -- ProcessAction is a 20-case dispatcher
// over direct state mutations, and internal/game/mana has zero importers).
// Offering a parameter the harness cannot honour teaches the model a move that
// silently does nothing, which is worse than not offering it. Phase 7 unlocks
// them alongside the legality layer.
func chooseActionTool() anthropic.BetaToolParam {
	return anthropic.BetaToolParam{
		Name: ToolChooseAction,
		Description: anthropic.String(
			"Take the action you have decided on. `choice` selects one option from the " +
				"numbered Choices list in the current decision: pass the option's id " +
				"(e.g. \"m3\"), its zero-based index (e.g. \"2\"), or yes/no where yes " +
				"means mulligan/confirm and no means keep/pass. The action is executed " +
				"immediately and the next decision follows."),
		Strict: anthropic.Bool(true),
		InputSchema: strictSchema(map[string]any{
			"choice": map[string]any{
				"type":        "string",
				"description": "Option id (\"m3\"), zero-based index (\"2\"), or yes/no. yes=mulligan/confirm, no=keep/pass.",
			},
			"amount": map[string]any{
				"type":        "integer",
				"description": "Amount value (for amount actions)",
			},
			"amounts": map[string]any{
				"type":        "array",
				"items":       map[string]any{"type": "integer"},
				"description": "Multiple amount values (for multi_amount)",
			},
			"pile": map[string]any{
				"type":        "integer",
				"description": "Pile number: 1 or 2",
			},
			"text": map[string]any{
				"type":        "string",
				"description": "Text value for a named-option choice (pick option by name)",
			},
		}, nil),
	}
}

// passPriorityTool ports reference/mcp-tools.json5:211 (pass_priority).
//
// board_cursor is omitted: it is a token optimisation for a transport that
// re-sends the board on every result, and we inline the board in the decision
// message instead. `until` is kept -- it costs nothing and reads naturally --
// but this engine has no step structure to skip to, so it is advisory only and
// resolves to the same pass macro.
func passPriorityTool() anthropic.BetaToolParam {
	return anthropic.BetaToolParam{
		Name: ToolPassPriority,
		Description: anthropic.String(
			"Decline to act and pass. Use this when none of the offered options is worth " +
				"taking. If a \"Pass the turn\" option is available it is taken; otherwise " +
				"the closest no-op is."),
		Strict: anthropic.Bool(true),
		InputSchema: strictSchema(map[string]any{
			"until": map[string]any{
				"type":        "string",
				"description": "Advisory: the step you would like to skip to. Recorded, not enforced.",
				"enum": []any{
					"upkeep", "draw", "precombat_main", "begin_combat",
					"declare_attackers", "declare_blockers", "end_combat",
					"postcombat_main", "end_of_turn", "my_turn", "stack_resolved",
				},
			},
		}, nil),
	}
}

// getOracleTextTool ports reference/mcp-tools.json5:702 (get_oracle_text).
//
// The Card Reference in the decision block prints oracle text on a card's FIRST
// appearance only (serialize.go, reference/pilot_rendering.py:61-64) -- the
// single largest token saving in the design. This tool is what makes that
// dedup safe: the model can ask again for anything it has forgotten.
func getOracleTextTool() anthropic.BetaToolParam {
	return anthropic.BetaToolParam{
		Name: ToolGetOracleText,
		Description: anthropic.String(
			"Get oracle text for cards. Use card_name/card_names for lookup by name, or " +
				"object_id/object_ids for objects already on the board."),
		Strict: anthropic.Bool(true),
		InputSchema: strictSchema(map[string]any{
			"card_name": map[string]any{
				"type":        "string",
				"description": "Card name",
			},
			"card_names": map[string]any{
				"type":        "array",
				"items":       map[string]any{"type": "string"},
				"description": "Card names (batch)",
			},
			"object_id": map[string]any{
				"type":        "string",
				"description": "In-game object ID (e.g. \"p3\")",
			},
			"object_ids": map[string]any{
				"type":        "array",
				"items":       map[string]any{"type": "string"},
				"description": "In-game object IDs (batch)",
			},
		}, nil),
	}
}

// sendChatMessageTool ports reference/mcp-tools.json5:166 (send_chat_message).
//
// message is the ONLY required parameter anywhere in the whole tool surface
// (§0.3), upstream included. The line does not go out from here: it is queued
// and drained by LLMPolicy.Line, which is a bot.ChatSource, so the runner's
// cadence rules (MaxChatMessagesPerTurn, the one-line-per-two-cycles floor)
// still govern it. A model that decides to say six things in one turn says two.
func sendChatMessageTool() anthropic.BetaToolParam {
	return anthropic.BetaToolParam{
		Name:        ToolSendChatMessage,
		Description: anthropic.String("Send a chat message to the other players at the table."),
		Strict:      anthropic.Bool(true),
		InputSchema: strictSchema(map[string]any{
			"message": map[string]any{
				"type":        "string",
				"description": "Message to send",
			},
		}, []string{"message"}),
	}
}

// ---------------------------------------------------------------------------
// Client
// ---------------------------------------------------------------------------

// Defaults taken from mage-bench's measured constants (§0.4) and the model
// table in §0.7.
const (
	// DefaultMaxTokens is upstream's MAX_TOKENS.
	DefaultMaxTokens = 20_000
	// DefaultRequestTimeout is upstream's LLM_REQUEST_TIMEOUT_SECS. Epoch #15
	// raised it from 45s because 45s was empirically not enough. It is PER
	// ATTEMPT: option.WithRequestTimeout applies to each retry, so the worst
	// case for one Complete call is timeout x (retries+1). The whole-loop bound
	// lives in policy.go.
	DefaultRequestTimeout = 120 * time.Second
	// DefaultMaxRetries matches the SDK default.
	DefaultMaxRetries = 2
)

// ClientOptions configures a Client.
type ClientOptions struct {
	// Model defaults to claude-sonnet-5. See policy.go for why Haiku is not the
	// obvious cheap default.
	Model anthropic.Model
	// MaxTokens defaults to DefaultMaxTokens.
	MaxTokens int64
	// Effort is the OutputConfig effort level. IT MUST BE EMPTY FOR HAIKU 4.5,
	// which does not support the parameter at all (§0.7); Options.normalise in
	// policy.go clears it rather than trusting the caller.
	Effort anthropic.BetaOutputConfigEffort
	// ThinkingBudget, when > 0, sets ThinkingConfigParamOfEnabled(N). This is
	// the Haiku path. On Sonnet 5 / Opus 5 a budget returns 400 -- those models
	// take adaptive thinking instead, which Adaptive selects.
	ThinkingBudget int64
	// Adaptive selects adaptive thinking (Sonnet 5 / Opus 5 only).
	Adaptive bool
	// RequestTimeout is per attempt. Defaults to DefaultRequestTimeout.
	RequestTimeout time.Duration
	// MaxRetries defaults to DefaultMaxRetries.
	MaxRetries int
	// System is the system prompt. Defaults to SystemPrompt().
	System string
}

func (o *ClientOptions) applyDefaults() {
	if o.Model == "" {
		o.Model = anthropic.ModelClaudeSonnet5
	}
	if o.MaxTokens <= 0 {
		o.MaxTokens = DefaultMaxTokens
	}
	if o.RequestTimeout <= 0 {
		o.RequestTimeout = DefaultRequestTimeout
	}
	if o.MaxRetries < 0 {
		o.MaxRetries = DefaultMaxRetries
	}
	if o.System == "" {
		o.System = SystemPrompt()
	}
}

// Client wraps one Transport with a frozen system prompt and tool list.
//
// FROZEN IS THE POINT. Both are built once, in New, and never touched again --
// no timestamp, no per-turn board summary, no map iteration order. They are the
// cacheable prefix of every request for the whole game; mutating either mid-game
// invalidates the cache on every subsequent request and nothing reports it
// except the bill. The only place to notice is Usage.CacheReadInputTokens,
// which Client accumulates for exactly that reason.
type Client struct {
	transport Transport
	opts      ClientOptions

	system []anthropic.BetaTextBlockParam
	tools  []anthropic.BetaToolUnionParam

	usage Usage
}

// Usage accumulates token counters across a game. Phase 6 turns these into
// metrics; Phase 5 needs CacheRead in particular, because a zero there after
// the first request means the cache is being rebuilt every turn.
type Usage struct {
	InputTokens         int64
	OutputTokens        int64
	CacheCreationTokens int64
	CacheReadTokens     int64
	Requests            int
}

func (u *Usage) add(m *anthropic.BetaMessage) {
	if m == nil {
		return
	}
	u.InputTokens += m.Usage.InputTokens
	u.OutputTokens += m.Usage.OutputTokens
	u.CacheCreationTokens += m.Usage.CacheCreationInputTokens
	u.CacheReadTokens += m.Usage.CacheReadInputTokens
	u.Requests++
}

// NewClient builds a Client over t.
func NewClient(t Transport, opts ClientOptions) *Client {
	opts.applyDefaults()
	tools := Tools()
	wrapped := make([]anthropic.BetaToolUnionParam, 0, len(tools))
	for i := range tools {
		wrapped = append(wrapped, anthropic.BetaToolUnionParam{OfTool: &tools[i]})
	}
	return &Client{
		transport: t,
		opts:      opts,
		tools:     wrapped,
		system: []anthropic.BetaTextBlockParam{{
			Text: opts.System,
			// The second of the four breakpoints. Tools take the first; the
			// remaining two go to messages (see context.go).
			CacheControl: anthropic.NewBetaCacheControlEphemeralParam(),
		}},
	}
}

// Usage returns the accumulated token counters.
func (c *Client) Usage() Usage { return c.usage }

// Provider implements Completer.
func (c *Client) Provider() string { return ProviderAnthropic }

// Model reports the configured model.
func (c *Client) Model() Model { return c.opts.Model }

// Complete renders the conversation and makes one request. It implements
// Completer: the SDK message is folded into the neutral Response so policy.go
// never sees an Anthropic type.
//
// ctx bounds the whole call; RequestTimeout bounds each attempt. Both are
// needed: a per-attempt timeout alone lets three attempts run to 360s, and a
// context alone gives a stuck connection no chance to be retried.
func (c *Client) Complete(ctx context.Context, conv *Conversation) (*Response, error) {
	if conv == nil {
		return nil, errors.New("llm: nil conversation")
	}
	params := anthropic.BetaMessageNewParams{
		Model:     c.opts.Model,
		MaxTokens: c.opts.MaxTokens,
		System:    c.system,
		Tools:     c.tools,
		Messages:  conv.Render(),
	}
	if c.opts.Effort != "" {
		params.OutputConfig = anthropic.BetaOutputConfigParam{Effort: c.opts.Effort}
	}
	switch {
	case c.opts.ThinkingBudget > 0:
		params.Thinking = anthropic.BetaThinkingConfigParamOfEnabled(c.opts.ThinkingBudget)
	case c.opts.Adaptive:
		params.Thinking = anthropic.BetaThinkingConfigParamUnion{
			OfAdaptive: &anthropic.BetaThinkingConfigAdaptiveParam{},
		}
	}

	msg, err := c.transport.New(ctx, params,
		option.WithRequestTimeout(c.opts.RequestTimeout),
		option.WithMaxRetries(c.opts.MaxRetries),
	)
	if err != nil {
		return nil, err
	}
	c.usage.add(msg)
	return responseFromAnthropic(msg), nil
}

// ToolUse is one tool call pulled out of a response.
type ToolUse struct {
	ID    string
	Name  string
	Input []byte
}

// ToolUses returns the tool_use blocks of a message, in order.
func ToolUses(m *anthropic.BetaMessage) []ToolUse {
	if m == nil {
		return nil
	}
	var out []ToolUse
	for _, b := range m.Content {
		if b.Type != "tool_use" {
			continue
		}
		out = append(out, ToolUse{ID: b.ID, Name: b.Name, Input: []byte(b.Input)})
	}
	return out
}

// ResponseText returns the concatenated text blocks of a message.
func ResponseText(m *anthropic.BetaMessage) string {
	if m == nil {
		return ""
	}
	s := ""
	for _, b := range m.Content {
		if b.Type == "text" {
			s += b.Text
		}
	}
	return s
}
