package llm

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// stub_test.go is the fake Transport every test in this package runs against.
//
// THERE IS NO LIVE CALL ANYWHERE IN THIS PACKAGE'S TESTS. The loop, the
// recovery matrix and the cache-breakpoint placement are all decided on our
// side of the wire, so all of them are testable without a credential -- which
// is also why Transport exists as an interface at all. live_test.go is the one
// test that would talk to the API, and it skips itself unless
// ANTHROPIC_API_KEY is set.

// stubTransport returns scripted responses and records every request.
type stubTransport struct {
	mu      sync.Mutex
	calls   []anthropic.BetaMessageNewParams
	respond func(n int, params anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error)
}

func (s *stubTransport) New(_ context.Context, params anthropic.BetaMessageNewParams, _ ...option.RequestOption) (*anthropic.BetaMessage, error) {
	s.mu.Lock()
	n := len(s.calls)
	s.calls = append(s.calls, params)
	s.mu.Unlock()
	return s.respond(n, params)
}

func (s *stubTransport) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.calls)
}

func (s *stubTransport) last() anthropic.BetaMessageNewParams {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls[len(s.calls)-1]
}

// toolUseResponse builds a response containing one tool_use block.
func toolUseResponse(id, name string, input any) *anthropic.BetaMessage {
	raw, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	return &anthropic.BetaMessage{
		Role:       "assistant",
		StopReason: anthropic.BetaStopReasonToolUse,
		Content: []anthropic.BetaContentBlockUnion{{
			Type:  "tool_use",
			ID:    id,
			Name:  name,
			Input: json.RawMessage(raw),
		}},
	}
}

// textResponse builds a response with no tool call at all.
func textResponse(text string, stop anthropic.BetaStopReason) *anthropic.BetaMessage {
	return &anthropic.BetaMessage{
		Role:       "assistant",
		StopReason: stop,
		Content:    []anthropic.BetaContentBlockUnion{{Type: "text", Text: text}},
	}
}

// lastUserText returns the text of the newest user message in a request, which
// is how a test reads the prompt the policy actually built.
func lastUserText(params anthropic.BetaMessageNewParams) string {
	for i := len(params.Messages) - 1; i >= 0; i-- {
		if params.Messages[i].Role != anthropic.BetaMessageParamRoleUser {
			continue
		}
		var b strings.Builder
		for _, c := range params.Messages[i].Content {
			if c.OfText != nil {
				b.WriteString(c.OfText.Text)
			}
		}
		if b.Len() > 0 {
			return b.String()
		}
	}
	return ""
}
