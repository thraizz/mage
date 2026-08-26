package llm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
)

// syntheticTurn appends one agentic turn: a user decision, an assistant
// message with a thinking-sized text block plus a tool_use, and the matching
// tool_result. That is the block shape the 20-block lookback actually has to
// survive -- three to four content blocks per decision, not one.
func syntheticTurn(c *Conversation, n int) {
	c.AppendUserText(fmt.Sprintf("decision %d", n))
	c.Append(anthropic.BetaMessageParam{
		Role: anthropic.BetaMessageParamRoleAssistant,
		Content: []anthropic.BetaContentBlockParamUnion{
			{OfText: &anthropic.BetaTextBlockParam{Text: fmt.Sprintf("reasoning %d", n)}},
			{OfToolUse: &anthropic.BetaToolUseBlockParam{
				ID: fmt.Sprintf("tu%d", n), Name: ToolChooseAction,
				Input: map[string]any{"choice": "m1"},
			}},
		},
	})
	c.AppendToolResults([]anthropic.BetaContentBlockParamUnion{
		anthropic.NewBetaToolResultBlock(fmt.Sprintf("tu%d", n),
			strings.Repeat("board line; ", 40), false),
	})
}

func TestCacheBreakpointsBoundedAndSpaced(t *testing.T) {
	c := NewConversation()
	for i := 0; i < 200; i++ {
		syntheticTurn(c, i)
	}
	msgs := c.Render()

	bps := BreakpointBlockIndexes(msgs)
	if len(bps) == 0 {
		t.Fatal("no cache breakpoints placed at all")
	}
	// Two of the four are already spent on the tool list and the system prompt
	// (client.go), so the message list may never use more than the rest.
	if len(bps) > MaxMessageBreakpoints {
		t.Fatalf("%d message breakpoints, budget is %d", len(bps), MaxMessageBreakpoints)
	}

	blocks := 0
	for _, m := range msgs {
		blocks += len(m.Content)
	}
	// The invariant that matters: a breakpoint looks back AT MOST 20 content
	// blocks for a cached prefix. The newest one must be within 20 blocks of
	// the end of the request, and consecutive ones within 20 of each other, or
	// the cache is rebuilt every turn at roughly 10x the cost -- silently.
	const lookback = 20
	if last := bps[len(bps)-1]; blocks-last > lookback {
		t.Errorf("last breakpoint is %d blocks from the end of %d, lookback is %d",
			blocks-last, blocks, lookback)
	}
	for i := 1; i < len(bps); i++ {
		if gap := bps[i] - bps[i-1]; gap > lookback {
			t.Errorf("breakpoints %d and %d are %d blocks apart, lookback is %d",
				bps[i-1], bps[i], gap, lookback)
		}
	}
}

func TestRenderDoesNotMutateHistory(t *testing.T) {
	c := NewConversation()
	for i := 0; i < 5; i++ {
		syntheticTurn(c, i)
	}
	_ = c.Render()
	// A breakpoint that leaked into the history would be re-sent next turn at a
	// stale position, and would accumulate until every request was over the
	// four-breakpoint limit.
	if bps := BreakpointBlockIndexes(c.History()); len(bps) != 0 {
		t.Fatalf("Render annotated the history in place: %v", bps)
	}
}

func TestWindowKeepsRecentVerbatimAndSummarisesOlder(t *testing.T) {
	c := NewConversation()
	for i := 0; i < 60; i++ {
		syntheticTurn(c, i)
	}
	c.SetStateSummary("Board: Turn 12")
	msgs := c.Render()

	// +4 of slack: both band boundaries walk backwards off a message that
	// opens with a tool_result, so the window can overshoot by a message or two
	// rather than ship an orphan.
	cap := ContextRecentCount + ContextSummaryCount + 1 + 4
	if len(msgs) > cap {
		t.Fatalf("window is %d messages, cap is %d", len(msgs), cap)
	}
	if c.Len() != 180 {
		t.Fatalf("history was truncated: %d messages", c.Len())
	}

	var bridge string
	for _, m := range msgs {
		for _, blk := range m.Content {
			if blk.OfText != nil && strings.Contains(blk.OfText.Text, bridgeMarker) {
				bridge = blk.OfText.Text
			}
		}
	}
	if bridge == "" {
		t.Fatal("no state-bridge message between the summarised and verbatim halves")
	}
	if !strings.Contains(bridge, "Board: Turn 12") {
		t.Errorf("bridge lost the board summary: %q", bridge)
	}

	// The last decision must survive verbatim -- it is the one being answered.
	tail := msgs[len(msgs)-1]
	if len(tail.Content) == 0 {
		t.Fatal("empty tail message")
	}
	// And an old tool result must have been compressed below the trigger.
	summarised := false
	for _, m := range msgs {
		for _, blk := range m.Content {
			if blk.OfToolResult == nil {
				continue
			}
			if strings.Contains(toolResultText(blk.OfToolResult), "elided") {
				summarised = true
			}
		}
	}
	if !summarised {
		t.Error("no tool result was summarised in the summary band")
	}
}

func TestWindowNeverStartsWithAnOrphanToolResult(t *testing.T) {
	c := NewConversation()
	for i := 0; i < 100; i++ {
		syntheticTurn(c, i)
	}
	msgs := c.Render()
	// A tool_result whose tool_use was windowed away is a 400 from the API, for
	// every remaining request of the game.
	if startsWithToolResult(msgs[0]) {
		t.Fatal("window opens on an orphaned tool_result")
	}
	ids := map[string]bool{}
	for _, m := range msgs {
		for _, blk := range m.Content {
			if blk.OfToolUse != nil {
				ids[blk.OfToolUse.ID] = true
			}
			if blk.OfToolResult != nil && !ids[blk.OfToolResult.ToolUseID] {
				t.Fatalf("tool_result %s has no tool_use in the window", blk.OfToolResult.ToolUseID)
			}
		}
	}
}

func TestStateSummaryRefreshInterval(t *testing.T) {
	c := NewConversation()
	if !c.NeedsStateRefresh() {
		t.Fatal("a fresh conversation should want a state summary")
	}
	c.SetStateSummary("Board: Turn 1")
	for i := 1; i < RenderInterval; i++ {
		c.Render()
		if c.NeedsStateRefresh() {
			t.Fatalf("refresh asked for after %d renders, interval is %d", i, RenderInterval)
		}
	}
	c.Render()
	if !c.NeedsStateRefresh() {
		t.Fatalf("no refresh after %d renders", RenderInterval)
	}
}

func TestResetClearsTranscriptOnly(t *testing.T) {
	c := NewConversation()
	for i := 0; i < 10; i++ {
		syntheticTurn(c, i)
	}
	c.Reset("continue")
	if c.Len() != 1 {
		t.Fatalf("history = %d messages after reset", c.Len())
	}
	if c.stateSummary != "" {
		t.Error("reset kept a stale board summary")
	}
}

func TestSummarizeToolResultKeepsTheHeader(t *testing.T) {
	body := "[Decision 7, snapshot=7] Turn 3 MAIN - bot-a\n" + strings.Repeat("x", 500)
	got := summarizeToolResult(ToolChooseAction, body)
	if !strings.Contains(got, "Decision 7") || !strings.Contains(got, "Turn 3") {
		t.Fatalf("summary lost the decision header: %q", got)
	}
	if len(got) >= len(body) {
		t.Fatalf("summary is not shorter: %d vs %d", len(got), len(body))
	}
}
