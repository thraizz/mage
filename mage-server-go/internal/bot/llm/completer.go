package llm

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

// completer.go is the multi-provider seam (plan Phase 5b).
//
// THE SEAM IS BELOW Policy, NOT BESIDE IT. There is exactly one LLMPolicy and
// there will stay exactly one: the enumerated-macro design, the board
// serialization, choice resolution, the failure matrix in recovery.go, the chat
// cadence, the two-timer logic and ThinkTimeReporter are all provider-agnostic
// by construction -- they are decided on our side of the wire. Duplicating them
// per provider would fork the part of this package that actually plays Magic in
// order to vary the part that only moves bytes.
//
// What genuinely differs is confined to a Completer:
//
//   - Prompt caching. Anthropic needs explicit cache_control breakpoints,
//     at most 4 per request, each with a 20-block lookback -- the arithmetic in
//     context.go. Gemini has implicit caching that is on by default and takes
//     no request field at all (see openai.go). One of these is a placement
//     problem and the other is a "do nothing, but do not break the prefix"
//     problem; there is no shared abstraction over them that is not a lie, so
//     each Completer renders the conversation with its own strategy.
//   - Thinking/reasoning configuration and the names of the token-accounting
//     fields.
//
// THE CONVERSATION IR IS THE ANTHROPIC BLOCK MODEL, deliberately and with a
// known wart. Conversation stores []anthropic.BetaMessageParam. That model
// (text / tool_use / tool_result blocks, results carrying the id of the call
// they answer) is a strict superset of the OpenAI chat model, so translating
// out of it in openai.go is total and mechanical, while translating the other
// way would not be. Reusing it also means Phase 5's windowing, summarisation
// and orphan-tool_result guards -- and every test that covers them -- carry
// over to the new provider unchanged rather than being reimplemented. The wart
// is that a Gemini-only build still links the Anthropic SDK's param types.

// Provider identifies which Completer to build. These are the values
// config.BotConfig.Provider accepts.
const (
	// ProviderGemini is the default: Gemini through Google's
	// OpenAI-compatible endpoint.
	ProviderGemini = "gemini"
	// ProviderAnthropic is the native Anthropic SDK path from Phase 5.
	ProviderAnthropic = "anthropic"
)

// Model is the model id.
//
// It is an ALIAS for anthropic.Model, not a new type, so that every Phase 5
// call site -- Options{Model: anthropic.ModelClaudeSonnet5} and friends --
// keeps compiling untouched. anthropic.Model is a bare string type with no
// behaviour, so nothing about a Gemini id is misrepresented by it beyond the
// name.
type Model = anthropic.Model

// Gemini model ids, verified against the OpenAI-compatible endpoint (Aug 2026).
//
// Pricing per MTok at the time of writing (ai.google.dev/gemini-api/docs/pricing),
// recorded here for the Phase 6 cost comparison:
//
//	gemini-3.7-flash       in $0.75  out $3.75  cached-in $0.075  (until 2026-12-31; 2x after)
//	gemini-3.6-flash       in $0.75  out $3.75  cached-in $0.075  (until 2026-12-31; 2x after)
//	gemini-3.5-flash       in $1.50  out $9.00  cached-in $0.15
//	gemini-3.5-flash-lite  in $0.30  out $2.50  cached-in $0.03
//
// For comparison, Anthropic (plan Sec 0.7): Sonnet 5 in $2 / out $10, Opus 5
// in $5 / out $25, Haiku 4.5 in $1 / out $5.
//
// DefaultGeminiModel is 3.7-flash: it is the newest, it is tied for the
// cheapest of the non-lite Flash line, and its implicit-cache threshold (4096
// tokens) is the same as every other 3.x Flash, so nothing is traded away by
// taking the newest. Flash-lite is cheaper still but is a smaller model, and
// Phase 6 -- not this phase -- is where a quality/cost trade gets measured.
const (
	ModelGemini37Flash     Model = "gemini-3.7-flash"
	ModelGemini36Flash     Model = "gemini-3.6-flash"
	ModelGemini35Flash     Model = "gemini-3.5-flash"
	ModelGemini35FlashLite Model = "gemini-3.5-flash-lite"

	// DefaultGeminiModel is the model used when the Gemini provider is
	// selected without an explicit id.
	DefaultGeminiModel = ModelGemini37Flash
	// DefaultAnthropicModel preserves Phase 5's default.
	DefaultAnthropicModel = anthropic.ModelClaudeSonnet5
)

// GeminiModels lists the model ids the Gemini provider accepts, so config
// validation and this package cannot drift apart.
func GeminiModels() []Model {
	return []Model{ModelGemini37Flash, ModelGemini36Flash, ModelGemini35Flash, ModelGemini35FlashLite}
}

// AnthropicModels lists the model ids the Anthropic provider accepts.
func AnthropicModels() []Model {
	return []Model{anthropic.ModelClaudeSonnet5, anthropic.ModelClaudeOpus5, anthropic.ModelClaudeHaiku4_5}
}

// Providers lists the accepted provider ids, gemini first because it is the
// default.
func Providers() []string { return []string{ProviderGemini, ProviderAnthropic} }

// DefaultModelFor reports the default model id for a provider.
func DefaultModelFor(provider string) Model {
	if provider == ProviderAnthropic {
		return DefaultAnthropicModel
	}
	return DefaultGeminiModel
}

// IsGeminiModel reports whether m is one of the Gemini ids.
func IsGeminiModel(m Model) bool {
	return strings.HasPrefix(strings.ToLower(string(m)), "gemini-")
}

// Response is the provider-neutral result of one completion.
//
// It is what the loop in policy.go reasons about, and it carries only the four
// things the loop actually uses: what the model said, what it called, whether
// it ran out of output budget, and the assistant turn to append to the
// transcript. Everything provider-shaped -- stop-reason spellings, usage field
// names, cache accounting -- is folded away by the Completer that produced it.
type Response struct {
	// Text is the concatenated text the model produced, for logs.
	Text string
	// ToolUses are the tool calls, in the order the model made them.
	ToolUses []ToolUse
	// Truncated reports that generation stopped on the output cap rather than
	// because the model was finished. It drives FailureTruncated.
	Truncated bool

	// assistant is the model's turn rendered back into the conversation IR, so
	// the next request replays it. The Completer owns this rendering because
	// only it knows what the provider sent.
	assistant anthropic.BetaMessageParam
}

// Completer is the provider seam: one request, one response, plus the token
// counters the caller needs to know whether caching is working.
//
// Complete takes the whole Conversation rather than a rendered message list
// because rendering IS the per-provider part: the Anthropic implementation
// annotates cache_control breakpoints on the way out, and the
// OpenAI-compatible one must not (see openai.go).
type Completer interface {
	// Complete renders conv and makes one request.
	Complete(ctx context.Context, conv *Conversation) (*Response, error)
	// Usage returns the counters accumulated so far this game.
	Usage() Usage
	// Provider reports which provider id built this Completer, for logs.
	Provider() string
}

var (
	_ Completer = (*Client)(nil)
	_ Completer = (*OpenAICompleter)(nil)
)

// CompleterOptions is what NewCompleter needs beyond ClientOptions.
type CompleterOptions struct {
	// Provider is ProviderGemini (default) or ProviderAnthropic.
	Provider string
	// APIKey is the credential. IT COMES FROM THE ENVIRONMENT AND NOWHERE ELSE
	// (anti-pattern 4): config.BotConfig.APIKey is tagged mapstructure:"-" and
	// is filled from GEMINI_API_KEY or ANTHROPIC_API_KEY by config.Load.
	APIKey string
	// BaseURL overrides the OpenAI-compatible endpoint. Empty means Gemini's.
	// It exists so openai_test.go can point at an httptest server, and so a
	// future OpenAI-compatible provider needs no code change at all.
	BaseURL string
	// HTTPClient overrides the HTTP client used by the OpenAI-compatible path.
	HTTPClient *http.Client
}

// NewCompleter builds the Completer for a provider.
func NewCompleter(co CompleterOptions, opts ClientOptions) (Completer, error) {
	switch co.Provider {
	case ProviderAnthropic:
		return NewClient(NewSDKTransport(co.APIKey), opts), nil
	case ProviderGemini, "":
		return NewOpenAICompleter(co, opts)
	default:
		return nil, fmt.Errorf("llm: unknown provider %q, want one of %v", co.Provider, Providers())
	}
}

// responseFromAnthropic folds an SDK message into the neutral Response.
func responseFromAnthropic(m *anthropic.BetaMessage) *Response {
	return &Response{
		Text:      ResponseText(m),
		ToolUses:  ToolUses(m),
		Truncated: isTruncated(m),
		assistant: m.ToParam(),
	}
}
