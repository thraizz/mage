package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
)

// openai.go is the OpenAI-compatible Completer. Gemini is the default provider
// behind it; any other OpenAI-compatible endpoint works by setting BaseURL.
//
// WHY RAW net/http AND NOT github.com/openai/openai-go. The whole surface this
// package needs is one route -- POST /chat/completions -- with one request
// shape and one response shape, non-streaming, no files, no assistants, no
// embeddings, no batches. Against that, a hand-rolled client is a couple of
// hundred lines of structs that say exactly what goes on the wire, and it buys
// three things a library actively takes away:
//
//  1. Exact control of the emitted function schema. Gemini accepts only "a
//     subset of the OpenAPI schema" for function declarations
//     (ai.google.dev/gemini-api/docs/function-calling), and the two keywords
//     Phase 5 puts on every Anthropic tool -- additionalProperties:false and
//     strict:true -- are outside the documented subset. They are dropped here
//     (see openAIToolsFrom), which is a decision that has to be expressible.
//  2. Reading usage.prompt_tokens_details.cached_tokens, which is how the cost
//     gate for this provider is verified at all, and which is the field the
//     OpenAI-compat layer maps Gemini's cached_content_token_count onto.
//  3. No second opinionated retry/timeout stack fighting the two-timer design
//     in policy.go.
//
// The cost is that we own the wire format. That is a genuine cost, and it is
// paid down by openai_test.go asserting the emitted JSON against an httptest
// server rather than against a mock of a library.
//
// PROMPT CACHING HERE IS IMPLICIT AND TAKES NO REQUEST FIELD.
// "Implicit caching is enabled by default for all Gemini 2.5 and newer models"
// (ai.google.dev/gemini-api/docs/caching). The minimum prefix for a hit is
// 4,096 tokens on every Gemini 3.x Flash model. Google's only stated caller-side
// guidance is to put "large and common contents at the beginning of your
// prompt" and to send "requests with similar prefix in a short amount of time".
// This package already does exactly that and has since Phase 5, for the
// Anthropic path's sake: the tool definitions are sorted and frozen, the system
// prompt is a constant, and both are built once in the constructor. So the
// correct implementation of Gemini caching is to send NO cache field and to
// break nothing -- which is why RenderUncached exists and why this file never
// touches cache_control. Explicit caching (CachedContent) is deliberately not
// used: it is a separate resource with its own lifecycle and TTL billing, and
// a game-long conversation whose prefix never changes is the exact workload
// implicit caching is for.
const (
	// GeminiOpenAIBaseURL is Google's OpenAI-compatible endpoint
	// (ai.google.dev/gemini-api/docs/openai).
	GeminiOpenAIBaseURL = "https://generativelanguage.googleapis.com/v1beta/openai/"

	// GeminiImplicitCacheMinTokens is the smallest prefix that can produce an
	// implicit cache hit on Gemini 3.x Flash. Recorded so Phase 6 can tell a
	// cache MISS ("the prefix changed") apart from a prompt that is simply too
	// short to cache -- an early-game decision may be under it, and that is not
	// a bug.
	GeminiImplicitCacheMinTokens = 4096

	// openAIRetryBaseDelay is the first backoff step. The SDK gives the
	// Anthropic path this for free; here it is ours to provide.
	openAIRetryBaseDelay = 500 * time.Millisecond
)

// OpenAICompleter talks to an OpenAI-compatible /chat/completions endpoint.
//
// Like Client, its system prompt and tool list are BUILT ONCE and never
// touched. That is the whole cache strategy for this provider (see the file
// comment): implicit caching keys on the request prefix, and the prefix is
// tools + system prompt + the oldest surviving messages.
type OpenAICompleter struct {
	http    *http.Client
	baseURL string
	apiKey  string
	opts    ClientOptions

	system openAIMessage
	tools  []openAITool

	usage Usage
}

// NewOpenAICompleter builds a Completer over an OpenAI-compatible endpoint.
//
// THE KEY IS PASSED IN, NEVER READ FROM A CONFIG FILE (anti-pattern 4).
// config.BotConfig.APIKey is populated from GEMINI_API_KEY and from nowhere
// else; config/config.yaml is a hard link to config.dev.yaml and is checked in.
func NewOpenAICompleter(co CompleterOptions, opts ClientOptions) (*OpenAICompleter, error) {
	if opts.Model == "" {
		opts.Model = DefaultGeminiModel
	}
	opts.applyDefaults()

	base := strings.TrimSuffix(co.BaseURL, "/")
	if base == "" {
		base = strings.TrimSuffix(GeminiOpenAIBaseURL, "/")
	}
	hc := co.HTTPClient
	if hc == nil {
		// No Timeout on the client itself: the deadline is per attempt and
		// comes from the context in complete(), so a retry gets a fresh one.
		hc = &http.Client{}
	}
	return &OpenAICompleter{
		http:    hc,
		baseURL: base,
		apiKey:  co.APIKey,
		opts:    opts,
		system:  openAIMessage{Role: "system", Content: opts.System},
		tools:   openAIToolsFrom(Tools()),
	}, nil
}

// Provider implements Completer.
func (c *OpenAICompleter) Provider() string { return ProviderGemini }

// Usage returns the accumulated token counters.
func (c *OpenAICompleter) Usage() Usage { return c.usage }

// Model reports the configured model id.
func (c *OpenAICompleter) Model() Model { return c.opts.Model }

// ---------------------------------------------------------------------------
// Wire types
// ---------------------------------------------------------------------------

type openAIRequest struct {
	Model      string          `json:"model"`
	Messages   []openAIMessage `json:"messages"`
	Tools      []openAITool    `json:"tools,omitempty"`
	ToolChoice string          `json:"tool_choice,omitempty"`
	MaxTokens  int64           `json:"max_tokens,omitempty"`
	// ReasoningEffort maps onto Gemini's thinking_level/thinking_budget
	// (ai.google.dev/gemini-api/docs/openai). Omitted when empty, which leaves
	// the model's own default thinking in place.
	ReasoningEffort string `json:"reasoning_effort,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content,omitempty"`
	// ToolCalls is set on assistant turns being replayed.
	ToolCalls []openAIToolCall `json:"tool_calls,omitempty"`
	// ToolCallID is set on role:"tool" messages, answering one call.
	ToolCallID string `json:"tool_call_id,omitempty"`
}

type openAIToolCall struct {
	ID       string             `json:"id"`
	Type     string             `json:"type"`
	Function openAIToolCallFunc `json:"function"`
}

type openAIToolCallFunc struct {
	Name string `json:"name"`
	// Arguments is a JSON *string*, not an object: that is the OpenAI wire
	// contract, and it is what ToolUse.Input already expects (raw bytes).
	Arguments string `json:"arguments"`
}

type openAITool struct {
	Type     string         `json:"type"`
	Function openAIToolFunc `json:"function"`
}

type openAIToolFunc struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters"`
}

type openAIResponse struct {
	Choices []struct {
		FinishReason string        `json:"finish_reason"`
		Message      openAIMessage `json:"message"`
	} `json:"choices"`
	Usage openAIUsage `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

type openAIUsage struct {
	PromptTokens        int64 `json:"prompt_tokens"`
	CompletionTokens    int64 `json:"completion_tokens"`
	PromptTokensDetails struct {
		CachedTokens int64 `json:"cached_tokens"`
	} `json:"prompt_tokens_details"`
	// CachedContentTokenCount is Gemini's NATIVE spelling. The OpenAI-compat
	// layer is documented as still in beta and has historically not always
	// populated prompt_tokens_details.cached_tokens, so both are read and the
	// larger wins. Getting this wrong would make a working cache look broken
	// and send someone hunting an invalidator that is not there.
	CachedContentTokenCount int64 `json:"cached_content_token_count"`
}

func (u openAIUsage) cachedTokens() int64 {
	if u.CachedContentTokenCount > u.PromptTokensDetails.CachedTokens {
		return u.CachedContentTokenCount
	}
	return u.PromptTokensDetails.CachedTokens
}

// ---------------------------------------------------------------------------
// Tool translation
// ---------------------------------------------------------------------------

// openAIToolsFrom converts the frozen Anthropic tool definitions into OpenAI
// function declarations, PRESERVING ORDER.
//
// Order is not cosmetic here either: tools serialize at the front of the
// request, so they are the head of the prefix implicit caching keys on. Tools()
// already sorts by name; this must not re-sort, re-map or otherwise shuffle.
//
// TWO KEYWORDS ARE DROPPED ON PURPOSE. Phase 5 puts strict:true and
// additionalProperties:false on every tool because Anthropic honours both and
// mage-bench's reflected schemas emit the latter. Gemini accepts only "a subset
// of the OpenAPI schema" for function declarations
// (ai.google.dev/gemini-api/docs/function-calling) and neither keyword is in
// the documented subset. Sending them risks a 400 on the first request of an
// unattended run in exchange for a constraint the model is not going to be
// told about anyway -- and the loop does not rely on either: policy.go
// tolerates unknown fields (it unmarshals into a fixed struct) and answers an
// unresolvable choice with an is_error tool_result rather than trusting the
// schema to have prevented it.
func openAIToolsFrom(tools []anthropic.BetaToolParam) []openAITool {
	out := make([]openAITool, 0, len(tools))
	for _, t := range tools {
		props, _ := t.InputSchema.Properties.(map[string]any)
		params := map[string]any{
			"type":       "object",
			"properties": props,
		}
		if len(t.InputSchema.Required) > 0 {
			params["required"] = t.InputSchema.Required
		}
		fn := openAIToolFunc{Name: t.Name, Parameters: params}
		if t.Description.Valid() {
			fn.Description = t.Description.Value
		}
		out = append(out, openAITool{Type: "function", Function: fn})
	}
	return out
}

// ---------------------------------------------------------------------------
// Message translation
// ---------------------------------------------------------------------------

// renderMessages converts the conversation IR into OpenAI chat messages.
//
// The mapping is total in this direction (see completer.go): a text block is
// content, a tool_use block is a tool_call on the assistant turn, and a
// tool_result block is its own role:"tool" message carrying the id of the call
// it answers. Thinking blocks have no OpenAI equivalent and are dropped -- they
// are not replayable across providers and Gemini does not accept them.
//
// TOOL RESULTS COME OUT AS SEPARATE MESSAGES, IN ORDER. The IR packs every
// result of a turn into one user message (the Anthropic API requires that);
// OpenAI requires the opposite -- one message per call, immediately after the
// assistant turn that made them. Getting the order wrong is a 400 for the rest
// of the game, which is why openai_test.go asserts the pairing.
func renderMessages(system openAIMessage, msgs []anthropic.BetaMessageParam) []openAIMessage {
	out := make([]openAIMessage, 0, len(msgs)+2)
	out = append(out, system)

	for _, m := range msgs {
		role := "user"
		if m.Role == anthropic.BetaMessageParamRoleAssistant {
			role = "assistant"
		}

		var text strings.Builder
		var calls []openAIToolCall
		var results []openAIMessage

		for _, blk := range m.Content {
			switch {
			case blk.OfText != nil:
				if text.Len() > 0 {
					text.WriteString("\n")
				}
				text.WriteString(blk.OfText.Text)
			case blk.OfToolUse != nil:
				args, err := json.Marshal(blk.OfToolUse.Input)
				if err != nil {
					args = []byte("{}")
				}
				calls = append(calls, openAIToolCall{
					ID: blk.OfToolUse.ID, Type: "function",
					Function: openAIToolCallFunc{Name: blk.OfToolUse.Name, Arguments: string(args)},
				})
			case blk.OfToolResult != nil:
				results = append(results, openAIMessage{
					Role:       "tool",
					ToolCallID: blk.OfToolResult.ToolUseID,
					Content:    toolResultText(blk.OfToolResult),
				})
			}
		}

		// Results first: they belong to the assistant turn already emitted,
		// and a user text block in the same IR message (there is never one in
		// practice) is new input that follows them.
		out = append(out, results...)
		if text.Len() > 0 || len(calls) > 0 {
			out = append(out, openAIMessage{Role: role, Content: text.String(), ToolCalls: calls})
		}
	}
	return out
}

// assistantParam converts one OpenAI assistant message back into the IR, so
// the next request replays it in the same shape every other turn has.
func assistantParam(m openAIMessage) anthropic.BetaMessageParam {
	blocks := make([]anthropic.BetaContentBlockParamUnion, 0, 1+len(m.ToolCalls))
	if strings.TrimSpace(m.Content) != "" {
		blocks = append(blocks, anthropic.NewBetaTextBlock(m.Content))
	}
	for _, tc := range m.ToolCalls {
		var input any = map[string]any{}
		if tc.Function.Arguments != "" {
			var decoded any
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &decoded); err == nil {
				input = decoded
			}
		}
		blocks = append(blocks, anthropic.BetaContentBlockParamUnion{
			OfToolUse: &anthropic.BetaToolUseBlockParam{
				ID: tc.ID, Name: tc.Function.Name, Input: input,
			},
		})
	}
	if len(blocks) == 0 {
		// An assistant turn must not be empty or the next request is a 400.
		blocks = append(blocks, anthropic.NewBetaTextBlock("(no content)"))
	}
	return anthropic.BetaMessageParam{
		Role:    anthropic.BetaMessageParamRoleAssistant,
		Content: blocks,
	}
}

// toolUsesFrom pulls the neutral ToolUse list out of an assistant message.
func toolUsesFrom(m openAIMessage) []ToolUse {
	var out []ToolUse
	for _, tc := range m.ToolCalls {
		args := tc.Function.Arguments
		if strings.TrimSpace(args) == "" {
			args = "{}"
		}
		out = append(out, ToolUse{ID: tc.ID, Name: tc.Function.Name, Input: []byte(args)})
	}
	return out
}

// reasoningEffort maps ClientOptions.Effort onto the values the endpoint
// accepts.
//
// Gemini's OpenAI layer takes none|low|medium|high. Anthropic additionally
// offers xhigh and max; those clamp to high rather than being sent through and
// 400ing. Empty stays empty, which omits the field and leaves the model's
// default thinking alone.
func reasoningEffort(e anthropic.BetaOutputConfigEffort) string {
	switch strings.ToLower(string(e)) {
	case "":
		return ""
	case "none", "low", "medium", "high":
		return strings.ToLower(string(e))
	case "xhigh", "max":
		return "high"
	default:
		return ""
	}
}

// ---------------------------------------------------------------------------
// The request
// ---------------------------------------------------------------------------

// Complete implements Completer.
//
// ctx bounds the whole call -- it is policy.go's stall guard. RequestTimeout
// bounds each attempt, and MaxRetries governs how many attempts there are, so
// the worst case is RequestTimeout x (MaxRetries+1) unless ctx cuts it short
// first. That is the same two-timer arrangement the Anthropic path gets from
// the SDK, reproduced here rather than assumed.
func (c *OpenAICompleter) Complete(ctx context.Context, conv *Conversation) (*Response, error) {
	if conv == nil {
		return nil, fmt.Errorf("llm: nil conversation")
	}
	// RenderUncached, not Render: cache_control breakpoints are an Anthropic
	// concept and this endpoint has no field for them. Gemini's caching is
	// implicit (see the file comment).
	msgs := renderMessages(c.system, conv.RenderUncached())

	req := openAIRequest{
		Model:           string(c.opts.Model),
		Messages:        msgs,
		Tools:           c.tools,
		ToolChoice:      "auto",
		MaxTokens:       c.opts.MaxTokens,
		ReasoningEffort: reasoningEffort(c.opts.Effort),
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("llm: encoding request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= c.opts.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := openAIRetryBaseDelay << (attempt - 1)
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
			}
		}
		resp, err := c.attempt(ctx, body)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		// Only a retryable class is worth another attempt. A 400 will be a 400
		// again, and a cancelled stall guard is over.
		if Classify(err) != FailureRetryable {
			return nil, err
		}
	}
	return nil, lastErr
}

func (c *OpenAICompleter) attempt(ctx context.Context, body []byte) (*Response, error) {
	attemptCtx, cancel := context.WithTimeout(ctx, c.opts.RequestTimeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(attemptCtx, http.MethodPost,
		c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	// 1 MiB is far above any legitimate response and stops a misbehaving
	// endpoint from being an OOM.
	raw, err := io.ReadAll(io.LimitReader(httpResp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode > 299 {
		return nil, &APIError{
			StatusCode: httpResp.StatusCode,
			Provider:   ProviderGemini,
			// NEVER the request: it carries the whole transcript, and the
			// header carries the key.
			Message: truncateForLog(string(raw)),
		}
	}

	var parsed openAIResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("llm: unreadable response: %w", err)
	}
	if parsed.Error != nil {
		return nil, &APIError{StatusCode: httpResp.StatusCode, Provider: ProviderGemini,
			Message: truncateForLog(parsed.Error.Message)}
	}
	if len(parsed.Choices) == 0 {
		// Not an error: it is an empty response, which the recovery matrix
		// already has a rule for (FailureEmpty -> nudge, then reset, then
		// autopilot). Returning an error here would misclassify it as a
		// transport problem and retry it forever.
		c.addUsage(parsed.Usage)
		return &Response{assistant: assistantParam(openAIMessage{Role: "assistant"})}, nil
	}

	choice := parsed.Choices[0]
	c.addUsage(parsed.Usage)
	return &Response{
		Text:      choice.Message.Content,
		ToolUses:  toolUsesFrom(choice.Message),
		Truncated: choice.FinishReason == "length" || choice.FinishReason == "max_tokens",
		assistant: assistantParam(choice.Message),
	}, nil
}

func (c *OpenAICompleter) addUsage(u openAIUsage) {
	c.usage.InputTokens += u.PromptTokens
	c.usage.OutputTokens += u.CompletionTokens
	// CacheReadTokens is the cross-provider name for "input tokens that were a
	// cache hit". Implicit caching has no creation step and no creation charge,
	// so CacheCreationTokens stays zero on this provider -- which is a fact
	// about Gemini, not a missing implementation.
	c.usage.CacheReadTokens += u.cachedTokens()
	c.usage.Requests++
}

func truncateForLog(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 512 {
		return s[:512] + "..."
	}
	return s
}
