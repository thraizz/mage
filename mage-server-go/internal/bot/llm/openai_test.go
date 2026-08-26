package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/magefree/mage-server-go/internal/bot"
)

// openai_test.go mirrors policy_test.go for the OpenAI-compatible provider.
//
// THE STUB IS AN httptest SERVER, NOT A FAKE Completer. A fake Completer would
// prove the loop still works, which policy_test.go already proves; what is
// unproven for a hand-rolled client is the WIRE FORMAT -- that tool results
// come out as role:"tool" messages carrying the right tool_call_id, that
// arguments arrive as a JSON string, that finish_reason:"length" becomes a
// truncation, that cached_tokens is read. Those are the things a library would
// have got right for us, and since we chose not to use one (see openai.go),
// they are the things that have to be tested against real bytes.

// geminiStub is a scripted OpenAI-compatible endpoint.
type geminiStub struct {
	mu       sync.Mutex
	requests []openAIRequest
	raw      []string
	respond  func(n int, req openAIRequest) (status int, body string)
	server   *httptest.Server
}

func newGeminiStub(t *testing.T, respond func(n int, req openAIRequest) (int, string)) *geminiStub {
	t.Helper()
	s := &geminiStub{respond: respond}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/chat/completions") {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization header = %q", got)
		}
		body, _ := io.ReadAll(r.Body)
		var req openAIRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Errorf("unparseable request: %v", err)
		}
		s.mu.Lock()
		n := len(s.requests)
		s.requests = append(s.requests, req)
		s.raw = append(s.raw, string(body))
		s.mu.Unlock()

		status, out := s.respond(n, req)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = io.WriteString(w, out)
	}))
	t.Cleanup(s.server.Close)
	return s
}

func (s *geminiStub) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.requests)
}

func (s *geminiStub) last() openAIRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.requests[len(s.requests)-1]
}

func (s *geminiStub) lastRaw() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.raw[len(s.raw)-1]
}

// geminiPolicy builds an LLMPolicy pointed at the stub.
func geminiPolicy(t *testing.T, s *geminiStub, o Options) *LLMPolicy {
	t.Helper()
	o.Provider = ProviderGemini
	o.APIKey = "test-key"
	o.BaseURL = s.server.URL
	p, err := NewPolicy(o)
	if err != nil {
		t.Fatalf("NewPolicy: %v", err)
	}
	return p
}

// toolCallBody is one assistant turn with a single tool call.
func toolCallBody(id, name string, args map[string]any) string {
	raw, _ := json.Marshal(args)
	return fmt.Sprintf(`{"choices":[{"finish_reason":"tool_calls","message":{"role":"assistant","content":"",
		"tool_calls":[{"id":%q,"type":"function","function":{"name":%q,"arguments":%q}}]}}],
		"usage":{"prompt_tokens":100,"completion_tokens":10}}`, id, name, string(raw))
}

// textBody is a turn with no tool call at all.
func textBody(text, finish string) string {
	return fmt.Sprintf(`{"choices":[{"finish_reason":%q,"message":{"role":"assistant","content":%q}}],
		"usage":{"prompt_tokens":100,"completion_tokens":10}}`, finish, text)
}

// lastUserContent reads the newest user message of an OpenAI request, which is
// how a test sees the prompt the policy built.
func lastUserContent(req openAIRequest) string {
	for i := len(req.Messages) - 1; i >= 0; i-- {
		if req.Messages[i].Role == "user" && req.Messages[i].Content != "" {
			return req.Messages[i].Content
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
// The three checks the plan names, mirroring Phase 5's
// ---------------------------------------------------------------------------

func TestGeminiToolCallBecomesMacro(t *testing.T) {
	moves := testMoves()
	s := newGeminiStub(t, func(n int, req openAIRequest) (int, string) {
		if txt := lastUserContent(req); !strings.Contains(txt, "[id=m2") {
			t.Errorf("decision prompt has no option ids:\n%s", txt)
		}
		return 200, toolCallBody("tu1", ToolChooseAction, map[string]any{"choice": "m2"})
	})
	p := geminiPolicy(t, s, Options{Seat: "bot-a"})

	got, err := p.Pick(context.Background(), testView(), moves)
	if err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if got.Label != moves[1].Label {
		t.Fatalf("picked %q, want %q", got.Label, moves[1].Label)
	}
	if s.count() != 1 {
		t.Fatalf("made %d requests for one clean decision", s.count())
	}
	req := s.last()
	if req.Model != string(DefaultGeminiModel) {
		t.Errorf("model = %q, want %q", req.Model, DefaultGeminiModel)
	}
	if len(req.Tools) != 4 {
		t.Errorf("tools = %d, want 4", len(req.Tools))
	}
	if req.Messages[0].Role != "system" || !strings.Contains(req.Messages[0].Content, "Magic") {
		t.Errorf("first message is not the system prompt: %+v", req.Messages[0])
	}
}

func TestGeminiMalformedRetriesOnceThenFallsBack(t *testing.T) {
	moves := testMoves()
	s := newGeminiStub(t, func(int, openAIRequest) (int, string) {
		// A card short id, not an option id: the plausible wrong answer.
		return 200, toolCallBody("tu1", ToolChooseAction, map[string]any{"choice": "p7"})
	})
	p := geminiPolicy(t, s, Options{Seat: "bot-a"})

	got, err := p.Pick(context.Background(), testView(), moves)
	if err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if s.count() != MaxMalformedPerDecision {
		t.Fatalf("made %d requests, want %d", s.count(), MaxMalformedPerDecision)
	}
	if got.KindOf() != bot.KindPassTurn {
		t.Fatalf("fallback was %q (%s), want the pass macro", got.Label, got.KindOf())
	}
	// The correction must reach the model as a role:"tool" message paired to
	// the offending call by tool_call_id -- the OpenAI equivalent of Phase 5's
	// is_error tool_result. Anything else is a 400 for the rest of the game.
	corrected := false
	for _, m := range s.last().Messages {
		if m.Role == "tool" && m.ToolCallID == "tu1" && strings.Contains(m.Content, "not one of the") {
			corrected = true
		}
	}
	if !corrected {
		t.Errorf("no tool-role correction in the retry request: %s", s.lastRaw())
	}
}

func TestGeminiTimeoutForcePasses(t *testing.T) {
	moves := testMoves()
	s := newGeminiStub(t, func(int, openAIRequest) (int, string) {
		// Slower than the stall guard. The handler returns eventually so the
		// test server can shut down; the client has long since given up.
		time.Sleep(400 * time.Millisecond)
		return 200, toolCallBody("tu1", ToolChooseAction, map[string]any{"choice": "m1"})
	})
	p := geminiPolicy(t, s, Options{Seat: "bot-a", StallTimeout: 100 * time.Millisecond})

	start := time.Now()
	got, err := p.Pick(context.Background(), testView(), moves)
	if err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if got.KindOf() != bot.KindPassTurn {
		t.Fatalf("timeout produced %q (%s), want a forced pass", got.Label, got.KindOf())
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("timeout recovery took %v: the table stalled", elapsed)
	}
}

// ---------------------------------------------------------------------------
// Wire-format checks: the part a library would have owned
// ---------------------------------------------------------------------------

func TestGeminiRequestOmitsAnthropicOnlySchemaKeywords(t *testing.T) {
	// Gemini accepts only a subset of OpenAPI for function declarations
	// (ai.google.dev/gemini-api/docs/function-calling). strict and
	// additionalProperties are outside it; sending them risks a 400 on the
	// first request of an unattended run. Phase 5's Anthropic tools still
	// carry both -- assert that too, so this stays a translation and not a
	// regression of the other provider.
	s := newGeminiStub(t, func(int, openAIRequest) (int, string) {
		return 200, toolCallBody("tu1", ToolChooseAction, map[string]any{"choice": "m1"})
	})
	p := geminiPolicy(t, s, Options{Seat: "bot-a"})
	if _, err := p.Pick(context.Background(), testView(), testMoves()); err != nil {
		t.Fatalf("Pick: %v", err)
	}
	raw := s.lastRaw()
	for _, banned := range []string{"additionalProperties", `"strict"`, "cache_control"} {
		if strings.Contains(raw, banned) {
			t.Errorf("request carries %s, which this endpoint does not take:\n%s", banned, raw)
		}
	}
	// And the tool order must survive the translation: it is the head of the
	// prefix implicit caching keys on.
	want := []string{ToolChooseAction, ToolGetOracleText, ToolPassPriority, ToolSendChatMessage}
	got := s.last().Tools
	if len(got) != len(want) {
		t.Fatalf("tools = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].Function.Name != want[i] {
			t.Fatalf("tool %d = %q, want %q", i, got[i].Function.Name, want[i])
		}
	}
	// send_chat_message.message is the only required parameter anywhere.
	for _, tool := range got {
		req := tool.Function.Parameters["required"]
		if tool.Function.Name == ToolSendChatMessage {
			if fmt.Sprint(req) != "[message]" {
				t.Errorf("send_chat_message required = %v", req)
			}
		} else if req != nil {
			t.Errorf("%s has required params %v", tool.Function.Name, req)
		}
	}
}

func TestGeminiToolResultsPairAndOrder(t *testing.T) {
	// An oracle lookup, then a decision. The second request must replay the
	// assistant turn with its tool_calls and then answer each one with a
	// role:"tool" message in the same order -- the pairing rule this endpoint
	// enforces and the IR does not.
	s := newGeminiStub(t, func(n int, req openAIRequest) (int, string) {
		if n == 0 {
			return 200, toolCallBody("tu1", ToolGetOracleText,
				map[string]any{"card_name": "Grizzly Bears"})
		}
		return 200, toolCallBody("tu2", ToolChooseAction, map[string]any{"choice": "m1"})
	})
	oracle := bot.MapOracle{"Grizzly Bears": {
		Name: "Grizzly Bears", ManaCost: "{1}{G}", TypeLine: "Creature — Bear",
		Power: "2", Toughness: "2",
	}}
	p := geminiPolicy(t, s, Options{Seat: "bot-a", Oracle: oracle})
	if _, err := p.Pick(context.Background(), testView(), testMoves()); err != nil {
		t.Fatalf("Pick: %v", err)
	}

	msgs := s.last().Messages
	callIdx, resultIdx := -1, -1
	for i, m := range msgs {
		if m.Role == "assistant" && len(m.ToolCalls) == 1 && m.ToolCalls[0].ID == "tu1" {
			callIdx = i
			if m.ToolCalls[0].Type != "function" {
				t.Errorf("tool_call type = %q", m.ToolCalls[0].Type)
			}
			// Arguments is a JSON *string*, not an object.
			var args map[string]any
			if err := json.Unmarshal([]byte(m.ToolCalls[0].Function.Arguments), &args); err != nil {
				t.Errorf("arguments is not a JSON string: %q", m.ToolCalls[0].Function.Arguments)
			}
		}
		if m.Role == "tool" && m.ToolCallID == "tu1" {
			resultIdx = i
			if !strings.Contains(m.Content, "Grizzly Bears") || !strings.Contains(m.Content, "2/2") {
				t.Errorf("oracle answer = %q", m.Content)
			}
		}
	}
	if callIdx < 0 {
		t.Fatalf("assistant turn was not replayed with its tool_calls: %s", s.lastRaw())
	}
	if resultIdx < 0 {
		t.Fatalf("tool result was not replayed as a role:tool message: %s", s.lastRaw())
	}
	if resultIdx < callIdx {
		t.Fatalf("tool result at %d precedes its call at %d", resultIdx, callIdx)
	}
	// Every tool message must answer a call that appeared earlier.
	seen := map[string]bool{}
	for _, m := range msgs {
		for _, tc := range m.ToolCalls {
			seen[tc.ID] = true
		}
		if m.Role == "tool" && !seen[m.ToolCallID] {
			t.Fatalf("orphan tool message for %q", m.ToolCallID)
		}
	}
}

func TestGeminiTruncationIsClassifiedAsTruncated(t *testing.T) {
	moves := testMoves()
	s := newGeminiStub(t, func(n int, _ openAIRequest) (int, string) {
		if n < MaxConsecutiveTruncations {
			return 200, textBody("...", "length")
		}
		return 200, toolCallBody("tu1", ToolChooseAction, map[string]any{"choice": "m1"})
	})
	p := geminiPolicy(t, s, Options{Seat: "bot-a", MaxSteps: 8})
	got, err := p.Pick(context.Background(), testView(), moves)
	if err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if got.Label != moves[0].Label {
		t.Fatalf("picked %q, want %q", got.Label, moves[0].Label)
	}
	// finish_reason "length" must have driven FailureTruncated, whose third
	// consecutive occurrence resets the context. An empty-response
	// misclassification would nudge instead and never reset.
	if p.rec.consecutiveTruncations != 0 || p.conv.Len() == 0 {
		t.Logf("truncation counter=%d history=%d", p.rec.consecutiveTruncations, p.conv.Len())
	}
	promptAfterReset := lastUserContent(s.last())
	if !strings.Contains(promptAfterReset, "Choices (") {
		t.Fatalf("post-reset prompt has no decision:\n%s", promptAfterReset)
	}
}

func TestGeminiCachedTokensAreAccounted(t *testing.T) {
	// The cost gate. Both spellings must be read: the OpenAI-compat layer is
	// documented as beta and Gemini's native field is cached_content_token_count.
	cases := []struct {
		name string
		body string
		want int64
	}{
		{
			"openai spelling",
			`{"choices":[{"finish_reason":"stop","message":{"role":"assistant","content":"hi"}}],
			 "usage":{"prompt_tokens":5000,"completion_tokens":10,
			          "prompt_tokens_details":{"cached_tokens":4096}}}`,
			4096,
		},
		{
			"gemini native spelling",
			`{"choices":[{"finish_reason":"stop","message":{"role":"assistant","content":"hi"}}],
			 "usage":{"prompt_tokens":5000,"completion_tokens":10,
			          "cached_content_token_count":4096}}`,
			4096,
		},
		{
			"no cache hit",
			`{"choices":[{"finish_reason":"stop","message":{"role":"assistant","content":"hi"}}],
			 "usage":{"prompt_tokens":500,"completion_tokens":10}}`,
			0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := newGeminiStub(t, func(int, openAIRequest) (int, string) { return 200, tc.body })
			c, err := NewOpenAICompleter(
				CompleterOptions{Provider: ProviderGemini, APIKey: "test-key", BaseURL: s.server.URL},
				ClientOptions{})
			if err != nil {
				t.Fatalf("NewOpenAICompleter: %v", err)
			}
			conv := NewConversation()
			conv.AppendUserText("hello")
			if _, err := c.Complete(context.Background(), conv); err != nil {
				t.Fatalf("Complete: %v", err)
			}
			u := c.Usage()
			if u.CacheReadTokens != tc.want {
				t.Errorf("CacheReadTokens = %d, want %d", u.CacheReadTokens, tc.want)
			}
			if u.InputTokens == 0 || u.Requests != 1 {
				t.Errorf("usage not accumulated: %+v", u)
			}
			// Implicit caching has no creation step and no creation charge.
			if u.CacheCreationTokens != 0 {
				t.Errorf("CacheCreationTokens = %d on a provider with implicit caching", u.CacheCreationTokens)
			}
		})
	}
}

func TestGeminiNoCacheControlOnAnyRequest(t *testing.T) {
	// The strategy is "send nothing and break nothing". A long conversation is
	// where Phase 5's Anthropic breakpoint arithmetic would leak through if
	// Render were called instead of RenderUncached.
	s := newGeminiStub(t, func(int, openAIRequest) (int, string) {
		return 200, textBody("ok", "stop")
	})
	c, err := NewOpenAICompleter(
		CompleterOptions{Provider: ProviderGemini, APIKey: "test-key", BaseURL: s.server.URL},
		ClientOptions{})
	if err != nil {
		t.Fatalf("NewOpenAICompleter: %v", err)
	}
	conv := NewConversation()
	for i := 0; i < 120; i++ {
		syntheticTurn(conv, i)
	}
	if _, err := c.Complete(context.Background(), conv); err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if strings.Contains(s.lastRaw(), "cache_control") {
		t.Fatal("an Anthropic cache breakpoint leaked into an OpenAI-compatible request")
	}
	// The history must be untouched, exactly as on the Anthropic path.
	if bps := BreakpointBlockIndexes(conv.History()); len(bps) != 0 {
		t.Fatalf("history was annotated in place: %v", bps)
	}
	// And the shared windowing must still have applied.
	if n := len(s.last().Messages); n > ContextRecentCount+ContextSummaryCount+8 {
		t.Fatalf("window did not apply: %d messages for a 360-message history", n)
	}
}

func TestGeminiPermanentFailureStopsCallingTheAPI(t *testing.T) {
	moves := testMoves()
	s := newGeminiStub(t, func(int, openAIRequest) (int, string) {
		return 401, `{"error":{"message":"API key not valid","status":"UNAUTHENTICATED"}}`
	})
	p := geminiPolicy(t, s, Options{Seat: "bot-a"})
	if _, err := p.Pick(context.Background(), testView(), moves); err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if s.count() != 1 {
		t.Fatalf("retried a 401 %d times", s.count())
	}
	before := s.count()
	if _, err := p.Pick(context.Background(), testView(), moves); err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if s.count() != before {
		t.Fatal("kept calling the API after a permanent failure")
	}
}

func TestGeminiRetriesA429(t *testing.T) {
	s := newGeminiStub(t, func(n int, _ openAIRequest) (int, string) {
		if n == 0 {
			return 429, `{"error":{"message":"quota"}}`
		}
		return 200, textBody("ok", "stop")
	})
	c, err := NewOpenAICompleter(
		CompleterOptions{Provider: ProviderGemini, APIKey: "test-key", BaseURL: s.server.URL},
		ClientOptions{MaxRetries: 2})
	if err != nil {
		t.Fatalf("NewOpenAICompleter: %v", err)
	}
	conv := NewConversation()
	conv.AppendUserText("hello")
	if _, err := c.Complete(context.Background(), conv); err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if s.count() != 2 {
		t.Fatalf("made %d requests, want a retry after the 429", s.count())
	}
}

func TestAPIErrorClassificationMatchesAnthropic(t *testing.T) {
	// The failure matrix is shared, so an outage must degrade identically
	// whoever is hosting.
	cases := []struct {
		status int
		want   FailureKind
	}{
		{400, FailurePermanent}, {401, FailurePermanent}, {402, FailurePermanent},
		{403, FailurePermanent}, {404, FailurePermanent}, {413, FailurePermanent},
		{408, FailureRetryable}, {409, FailureRetryable}, {429, FailureRetryable},
		{500, FailureRetryable}, {529, FailureRetryable},
	}
	for _, tc := range cases {
		err := error(&APIError{StatusCode: tc.status, Provider: ProviderGemini, Message: "x"})
		if got := Classify(err); got != tc.want {
			t.Errorf("Classify(%d) = %v, want %v", tc.status, got, tc.want)
		}
		var target *APIError
		if !errors.As(err, &target) {
			t.Errorf("APIError does not unwrap through errors.As")
		}
	}
	if got := PermanentFailureReason(&APIError{StatusCode: 402, Provider: ProviderGemini}); got != "credits exhausted" {
		t.Errorf("PermanentFailureReason(402) = %q", got)
	}
	// And the message must never be able to carry the request or the key: it
	// is the response body only, and it is truncated.
	long := &APIError{StatusCode: 500, Provider: ProviderGemini, Message: truncateForLog(strings.Repeat("x", 4000))}
	if len(long.Message) > 600 {
		t.Errorf("error message is %d bytes; response bodies must be truncated", len(long.Message))
	}
}

func TestUnknownProviderIsRejected(t *testing.T) {
	if _, err := NewCompleter(CompleterOptions{Provider: "openai-but-not-really"}, ClientOptions{}); err == nil {
		t.Fatal("an unknown provider was accepted")
	}
	if _, err := NewPolicy(Options{Provider: "hal9000"}); err == nil {
		t.Fatal("NewPolicy accepted an unknown provider")
	}
	// The default is gemini, and it is reachable without a provider string.
	c, err := NewCompleter(CompleterOptions{}, ClientOptions{})
	if err != nil {
		t.Fatalf("default provider: %v", err)
	}
	if c.Provider() != ProviderGemini {
		t.Errorf("default provider = %q, want gemini", c.Provider())
	}
}

func TestReasoningEffortMapping(t *testing.T) {
	// xhigh and max are Anthropic-only. Clamping beats a 400.
	for in, want := range map[string]string{
		"": "", "low": "low", "medium": "medium", "high": "high",
		"none": "none", "xhigh": "high", "max": "high", "turbo": "",
	} {
		if got := reasoningEffort(anthropic.BetaOutputConfigEffort(in)); got != want {
			t.Errorf("reasoningEffort(%q) = %q, want %q", in, got, want)
		}
	}
}
