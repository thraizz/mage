package llm

import (
	"context"
	"encoding/json"
	"sort"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

func TestToolsAreSortedAndStable(t *testing.T) {
	a := Tools()
	b := Tools()

	names := make([]string, len(a))
	for i, tool := range a {
		names[i] = tool.Name
	}
	if !sort.StringsAreSorted(names) {
		t.Fatalf("tools are not sorted by name: %v", names)
	}
	// Tools serialize at position 0. Two builds that differ by one swapped
	// entry would rebuild the whole cache every game, silently.
	for i := range a {
		if a[i].Name != b[i].Name {
			t.Fatalf("tool order is not stable: %v vs %v", a[i].Name, b[i].Name)
		}
	}
	want := []string{ToolChooseAction, ToolGetOracleText, ToolPassPriority, ToolSendChatMessage}
	if len(names) != len(want) {
		t.Fatalf("tool set changed: %v", names)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("tool set changed: got %v want %v", names, want)
		}
	}
}

func TestToolsStrictAndClosed(t *testing.T) {
	tools := Tools()
	for _, tool := range tools {
		if !tool.Strict.Valid() || !tool.Strict.Value {
			t.Errorf("%s: Strict not set", tool.Name)
		}
		extra := tool.InputSchema.ExtraFields["additionalProperties"]
		if v, ok := extra.(bool); !ok || v {
			t.Errorf("%s: additionalProperties is %v, want false", tool.Name, extra)
		}
	}
	// Only send_chat_message.message is required anywhere (plan Sec 0.3): the
	// model is never forced to supply a field it cannot reason about.
	for _, tool := range tools {
		switch tool.Name {
		case ToolSendChatMessage:
			if len(tool.InputSchema.Required) != 1 || tool.InputSchema.Required[0] != "message" {
				t.Errorf("send_chat_message required = %v", tool.InputSchema.Required)
			}
		default:
			if len(tool.InputSchema.Required) != 0 {
				t.Errorf("%s has required params %v", tool.Name, tool.InputSchema.Required)
			}
		}
	}
}

func TestCacheControlOnLastToolOnly(t *testing.T) {
	tools := Tools()
	for i, tool := range tools {
		marked := !param.IsOmitted(tool.CacheControl)
		want := i == len(tools)-1
		if marked != want {
			t.Errorf("%s: cache_control=%v, want %v", tool.Name, marked, want)
		}
	}
}

func TestSystemPromptCarriesCacheControl(t *testing.T) {
	st := &stubTransport{respond: func(int, anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		return textResponse("ok", anthropic.BetaStopReasonEndTurn), nil
	}}
	c := NewClient(st, ClientOptions{})
	conv := NewConversation()
	conv.AppendUserText("hello")
	if _, err := c.Complete(context.Background(), conv); err != nil {
		t.Fatalf("Complete: %v", err)
	}
	p := st.last()
	if len(p.System) != 1 {
		t.Fatalf("system blocks = %d", len(p.System))
	}
	if param.IsOmitted(p.System[0].CacheControl) {
		t.Fatal("system prompt has no cache_control breakpoint")
	}
	if p.Model != anthropic.ModelClaudeSonnet5 {
		t.Errorf("default model = %q, want claude-sonnet-5", p.Model)
	}
	if p.MaxTokens != DefaultMaxTokens {
		t.Errorf("max tokens = %d", p.MaxTokens)
	}
	if len(p.Tools) != 4 {
		t.Errorf("tools = %d", len(p.Tools))
	}
}

func TestRequestNeverExceedsFourBreakpoints(t *testing.T) {
	st := &stubTransport{respond: func(int, anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		return textResponse("ok", anthropic.BetaStopReasonEndTurn), nil
	}}
	c := NewClient(st, ClientOptions{})
	conv := NewConversation()
	for i := 0; i < 120; i++ {
		conv.AppendUserText("decision")
		conv.Append(anthropic.BetaMessageParam{
			Role: anthropic.BetaMessageParamRoleAssistant,
			Content: []anthropic.BetaContentBlockParamUnion{
				{OfText: &anthropic.BetaTextBlockParam{Text: "thinking"}},
				{OfToolUse: &anthropic.BetaToolUseBlockParam{ID: "t", Name: ToolChooseAction, Input: map[string]any{}}},
			},
		})
	}
	if _, err := c.Complete(context.Background(), conv); err != nil {
		t.Fatalf("Complete: %v", err)
	}
	p := st.last()

	total := len(BreakpointBlockIndexes(p.Messages))
	for _, tool := range p.Tools {
		if tool.OfTool != nil && !param.IsOmitted(tool.OfTool.CacheControl) {
			total++
		}
	}
	for _, sys := range p.System {
		if !param.IsOmitted(sys.CacheControl) {
			total++
		}
	}
	if total > MaxCacheBreakpoints {
		t.Fatalf("request carries %d cache_control blocks, API limit is %d", total, MaxCacheBreakpoints)
	}
}

func TestThinkingAndEffortPerModel(t *testing.T) {
	cases := []struct {
		name         string
		opts         Options
		wantEffort   anthropic.BetaOutputConfigEffort
		wantBudget   bool
		wantAdaptive bool
	}{
		{
			name:         "sonnet gets adaptive thinking and no budget",
			opts:         Options{Model: anthropic.ModelClaudeSonnet5, Effort: "medium"},
			wantEffort:   "medium",
			wantAdaptive: true,
		},
		{
			// Haiku 4.5 supports NEITHER adaptive thinking NOR effort (Sec 0.7).
			// Effort must be cleared even when the caller asked for it, or the
			// first request of an overnight run is a 400.
			name:       "haiku loses effort and gets a budget",
			opts:       Options{Model: anthropic.ModelClaudeHaiku4_5, Effort: "high"},
			wantEffort: "",
			wantBudget: true,
		},
		{
			name:         "opus gets adaptive thinking",
			opts:         Options{Model: anthropic.ModelClaudeOpus5},
			wantAdaptive: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			st := &stubTransport{respond: func(int, anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
				return textResponse("ok", anthropic.BetaStopReasonEndTurn), nil
			}}
			p := New(st, tc.opts)
			conv := NewConversation()
			conv.AppendUserText("hi")
			if _, err := p.client.Complete(context.Background(), conv); err != nil {
				t.Fatalf("Complete: %v", err)
			}
			got := st.last()
			if got.OutputConfig.Effort != tc.wantEffort {
				t.Errorf("effort = %q, want %q", got.OutputConfig.Effort, tc.wantEffort)
			}
			hasBudget := got.Thinking.OfEnabled != nil
			if hasBudget != tc.wantBudget {
				t.Errorf("thinking budget set = %v, want %v", hasBudget, tc.wantBudget)
			}
			if hasBudget && got.Thinking.OfEnabled.BudgetTokens >= got.MaxTokens {
				t.Errorf("thinking budget %d >= max tokens %d",
					got.Thinking.OfEnabled.BudgetTokens, got.MaxTokens)
			}
			adaptive := got.Thinking.OfAdaptive != nil
			if adaptive != tc.wantAdaptive {
				t.Errorf("adaptive = %v, want %v", adaptive, tc.wantAdaptive)
			}
		})
	}
}

func TestToolSchemasAreValidJSON(t *testing.T) {
	for _, tool := range Tools() {
		b, err := json.Marshal(tool)
		if err != nil {
			t.Fatalf("%s: %v", tool.Name, err)
		}
		var round map[string]any
		if err := json.Unmarshal(b, &round); err != nil {
			t.Fatalf("%s: %v", tool.Name, err)
		}
		schema, ok := round["input_schema"].(map[string]any)
		if !ok {
			t.Fatalf("%s: no input_schema in %s", tool.Name, b)
		}
		if schema["additionalProperties"] != false {
			t.Errorf("%s: additionalProperties survived as %v", tool.Name, schema["additionalProperties"])
		}
		if schema["type"] != "object" {
			t.Errorf("%s: schema type = %v", tool.Name, schema["type"])
		}
	}
}
