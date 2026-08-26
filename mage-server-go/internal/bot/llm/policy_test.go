package llm

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/magefree/mage-server-go/internal/bot"
)

// testView is the smallest SafeView the prompt builder needs.
func testView() *bot.SafeView {
	return &bot.SafeView{
		GameID:         "g1",
		ViewerID:       "bot-a",
		Turn:           3,
		ActivePlayerID: "bot-a",
		Me: &bot.SafePlayerView{
			PlayerID: "bot-a", Name: "bot-a", Life: 20, LibraryCount: 90,
			KeptHand: true,
			Hand:     []*bot.SafeCard{{ID: "c1", Name: "Forest", Type: "Land"}},
		},
		Opponents: []*bot.SafeOpponentView{
			{PlayerID: "bot-b", Name: "bot-b", Life: 20, HandCount: 7, LibraryCount: 90, KeptHand: true},
		},
	}
}

// testMoves mirrors what LegalMoves offers mid-turn: something to do, and a
// pass. The pass is the fallback every failure path must land on.
func testMoves() []bot.Macro {
	v := testView()
	moves := bot.LegalMoves(v)
	if len(moves) < 2 {
		panic("fixture produced too few moves")
	}
	return moves
}

func indexOfPass(moves []bot.Macro) int {
	for i, m := range moves {
		if m.KindOf() == bot.KindPassTurn {
			return i
		}
	}
	panic("no pass macro in fixture")
}

func TestPickToolCallBecomesMacro(t *testing.T) {
	moves := testMoves()
	st := &stubTransport{respond: func(n int, p anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		// The prompt must actually carry the enumerated options, or the id the
		// model returns is a guess.
		if txt := lastUserText(p); !strings.Contains(txt, "[id=m2") {
			t.Errorf("decision prompt has no option ids:\n%s", txt)
		}
		return toolUseResponse("tu1", ToolChooseAction, map[string]any{"choice": "m2"}), nil
	}}
	p := New(st, Options{Seat: "bot-a"})

	got, err := p.Pick(context.Background(), testView(), moves)
	if err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if got.Label != moves[1].Label {
		t.Fatalf("picked %q, want %q", got.Label, moves[1].Label)
	}
	if st.count() != 1 {
		t.Fatalf("made %d requests for one clean decision", st.count())
	}
}

func TestPickMalformedRetriesOnceThenFallsBack(t *testing.T) {
	moves := testMoves()
	st := &stubTransport{respond: func(n int, _ anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		// "p7" is a card short id, not an option id: exactly the plausible
		// wrong answer, and one that must never resolve to a macro by accident.
		return toolUseResponse("tu1", ToolChooseAction, map[string]any{"choice": "p7"}), nil
	}}
	p := New(st, Options{Seat: "bot-a"})

	got, err := p.Pick(context.Background(), testView(), moves)
	if err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if st.count() != MaxMalformedPerDecision {
		t.Fatalf("made %d requests, want %d (one correction, then out)",
			st.count(), MaxMalformedPerDecision)
	}
	if got.KindOf() != bot.KindPassTurn {
		t.Fatalf("fallback was %q (%s), want the pass macro", got.Label, got.KindOf())
	}
	// The correction must reach the model as an is_error TOOL RESULT -- the
	// block the API pairs with the offending tool_use -- not as a fresh user
	// message and not as a silent retry of a byte-identical request.
	corrected := false
	for _, m := range st.last().Messages {
		for _, c := range m.Content {
			if c.OfToolResult == nil {
				continue
			}
			if strings.Contains(toolResultText(c.OfToolResult), "not one of the") {
				if !c.OfToolResult.IsError.Valid() || !c.OfToolResult.IsError.Value {
					t.Error("correction was not flagged is_error")
				}
				corrected = true
			}
		}
	}
	if !corrected {
		t.Error("no correction tool_result in the retry request")
	}
}

func TestPickTimeoutForcePasses(t *testing.T) {
	moves := testMoves()
	st := &stubTransport{respond: func(int, anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		return nil, context.DeadlineExceeded
	}}
	p := New(st, Options{Seat: "bot-a", StallTimeout: 200 * time.Millisecond})

	start := time.Now()
	got, err := p.Pick(context.Background(), testView(), moves)
	if err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if got.KindOf() != bot.KindPassTurn {
		t.Fatalf("timeout produced %q (%s), want a forced pass", got.Label, got.KindOf())
	}
	// A timeout must not become a stall: the table has three other players in it.
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("timeout recovery took %v", elapsed)
	}
	if st.count() != 1 {
		t.Fatalf("retried a timeout %d times within one decision", st.count())
	}
}

func TestRepeatedTimeoutsDegradeToAutopilotInCharacter(t *testing.T) {
	moves := testMoves()
	st := &stubTransport{respond: func(int, anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		return nil, context.DeadlineExceeded
	}}
	p := New(st, Options{Seat: "bot-a", StallTimeout: 200 * time.Millisecond})

	for i := 0; i < MaxConsecutiveTimeouts; i++ {
		if _, err := p.Pick(context.Background(), testView(), moves); err != nil {
			t.Fatalf("Pick %d: %v", i, err)
		}
	}
	if !p.rec.Autopilot() {
		t.Fatalf("still not on autopilot after %d timeouts", MaxConsecutiveTimeouts)
	}
	before := st.count()
	got, err := p.Pick(context.Background(), testView(), moves)
	if err != nil {
		t.Fatalf("Pick on autopilot: %v", err)
	}
	if got.KindOf() != bot.KindPassTurn {
		t.Fatalf("autopilot picked %q", got.Label)
	}
	if st.count() != before {
		t.Fatal("autopilot is still making API requests")
	}

	// Degradation is in character: it says something, through the ordinary
	// ChatSource path, and only then stops playing properly.
	var lines []string
	for {
		line, ok := p.Line(context.Background(), testView(), moves[0], true)
		if !ok {
			break
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		t.Fatal("degraded silently")
	}
	if !containsLine(lines, DegradationLine) {
		t.Fatalf("no autopilot announcement in %v", lines)
	}
}

func containsLine(lines []string, want string) bool {
	for _, l := range lines {
		if l == want {
			return true
		}
	}
	return false
}

func TestPassPriorityResolvesToPass(t *testing.T) {
	moves := testMoves()
	st := &stubTransport{respond: func(int, anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		return toolUseResponse("tu1", ToolPassPriority, map[string]any{}), nil
	}}
	p := New(st, Options{Seat: "bot-a"})
	got, err := p.Pick(context.Background(), testView(), moves)
	if err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if got.KindOf() != bot.KindPassTurn {
		t.Fatalf("pass_priority gave %q (%s)", got.Label, got.KindOf())
	}
	if want := moves[indexOfPass(moves)].Label; got.Label != want {
		t.Fatalf("pass macro = %q, want %q", got.Label, want)
	}
}

func TestChatToolFeedsChatSource(t *testing.T) {
	moves := testMoves()
	st := &stubTransport{respond: func(n int, _ anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		if n == 0 {
			return toolUseResponse("tu1", ToolSendChatMessage,
				map[string]any{"message": "nice board"}), nil
		}
		return toolUseResponse("tu2", ToolChooseAction, map[string]any{"choice": "m1"}), nil
	}}
	p := New(st, Options{Seat: "bot-a"})
	if _, err := p.Pick(context.Background(), testView(), moves); err != nil {
		t.Fatalf("Pick: %v", err)
	}
	line, ok := p.Line(context.Background(), testView(), moves[0], true)
	if !ok || line != "nice board" {
		t.Fatalf("chat line = %q, ok=%v", line, ok)
	}
	if _, ok := p.Line(context.Background(), testView(), moves[0], true); ok {
		t.Fatal("chat queue did not drain")
	}
}

func TestOracleToolAnswersFromTheSameLookup(t *testing.T) {
	moves := testMoves()
	oracle := bot.MapOracle{"Grizzly Bears": {
		Name: "Grizzly Bears", ManaCost: "{1}{G}", TypeLine: "Creature — Bear",
		Power: "2", Toughness: "2",
	}}
	var answer string
	st := &stubTransport{respond: func(n int, p anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		if n == 0 {
			return toolUseResponse("tu1", ToolGetOracleText,
				map[string]any{"card_name": "Grizzly Bears"}), nil
		}
		for _, m := range p.Messages {
			for _, c := range m.Content {
				if c.OfToolResult != nil && c.OfToolResult.ToolUseID == "tu1" {
					answer = toolResultText(c.OfToolResult)
				}
			}
		}
		return toolUseResponse("tu2", ToolChooseAction, map[string]any{"choice": "m1"}), nil
	}}
	p := New(st, Options{Seat: "bot-a", Oracle: oracle})
	if _, err := p.Pick(context.Background(), testView(), moves); err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if !strings.Contains(answer, "Grizzly Bears") || !strings.Contains(answer, "2/2") {
		t.Fatalf("oracle answer = %q", answer)
	}
}

func TestEmptyResponseIsNudgedThenPassed(t *testing.T) {
	moves := testMoves()
	st := &stubTransport{respond: func(int, anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		return textResponse("I am thinking about it.", anthropic.BetaStopReasonEndTurn), nil
	}}
	p := New(st, Options{Seat: "bot-a", MaxSteps: 3})
	got, err := p.Pick(context.Background(), testView(), moves)
	if err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if got.KindOf() != bot.KindPassTurn {
		t.Fatalf("got %q, want a pass", got.Label)
	}
	if st.count() != 3 {
		t.Fatalf("made %d requests, want the step budget of 3", st.count())
	}
}

func TestPermanentFailureStopsCallingTheAPI(t *testing.T) {
	moves := testMoves()
	apiErr := &anthropic.Error{StatusCode: 401}
	st := &stubTransport{respond: func(int, anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		return nil, apiErr
	}}
	p := New(st, Options{Seat: "bot-a"})
	if _, err := p.Pick(context.Background(), testView(), moves); err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if st.count() != 1 {
		t.Fatalf("retried a 401 %d times", st.count())
	}
	before := st.count()
	if _, err := p.Pick(context.Background(), testView(), moves); err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if st.count() != before {
		t.Fatal("kept calling the API after a permanent failure")
	}
}

func TestNoMovesIsAnError(t *testing.T) {
	st := &stubTransport{respond: func(int, anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		t.Fatal("should not have called the API")
		return nil, nil
	}}
	p := New(st, Options{Seat: "bot-a"})
	if _, err := p.Pick(context.Background(), testView(), nil); !errors.Is(err, bot.ErrNoMoves) {
		t.Fatalf("err = %v, want ErrNoMoves", err)
	}
}

func TestResolveChoiceForms(t *testing.T) {
	mullView := testView()
	mullView.Me.KeptHand = false
	mullMoves := bot.LegalMoves(mullView)

	moves := testMoves()
	cases := []struct {
		choice string
		in     []bot.Macro
		want   bot.Kind
		ok     bool
	}{
		{"m1", moves, moves[0].KindOf(), true},
		{"M1", moves, moves[0].KindOf(), true},
		{" m2 ", moves, moves[1].KindOf(), true},
		{"0", moves, moves[0].KindOf(), true},
		{"yes", mullMoves, bot.KindMulligan, true},
		{"no", mullMoves, bot.KindKeepHand, true},
		{"no", moves, bot.KindPassTurn, true},
		{"m0", moves, "", false},
		{"m999", moves, "", false},
		{"p3", moves, "", false},
		{"", moves, "", false},
		{"maybe", moves, "", false},
	}
	for _, tc := range cases {
		idx, ok := resolveChoice(tc.choice, tc.in)
		if ok != tc.ok {
			t.Errorf("resolveChoice(%q) ok = %v, want %v", tc.choice, ok, tc.ok)
			continue
		}
		if ok && tc.in[idx].KindOf() != tc.want {
			t.Errorf("resolveChoice(%q) = %s, want %s", tc.choice, tc.in[idx].KindOf(), tc.want)
		}
	}
}

func TestThinkTimeIsReportedForPacingOverlap(t *testing.T) {
	moves := testMoves()
	st := &stubTransport{respond: func(int, anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		time.Sleep(20 * time.Millisecond)
		return toolUseResponse("tu1", ToolChooseAction, map[string]any{"choice": "m1"}), nil
	}}
	p := New(st, Options{Seat: "bot-a"})
	if _, err := p.Pick(context.Background(), testView(), moves); err != nil {
		t.Fatalf("Pick: %v", err)
	}
	// The runner subtracts this from the persona's pause so thinking time and
	// pacing overlap instead of stacking.
	if got := p.LastThinkTime(); got < 20*time.Millisecond {
		t.Fatalf("LastThinkTime = %v, want at least the request time", got)
	}
}

func TestConversationGrowsAcrossDecisions(t *testing.T) {
	moves := testMoves()
	st := &stubTransport{respond: func(int, anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		return toolUseResponse("tu1", ToolChooseAction, map[string]any{"choice": "m1"}), nil
	}}
	p := New(st, Options{Seat: "bot-a"})
	for i := 0; i < 3; i++ {
		if _, err := p.Pick(context.Background(), testView(), moves); err != nil {
			t.Fatalf("Pick %d: %v", i, err)
		}
	}
	// user decision + assistant + tool results, three times over.
	if got := p.conv.Len(); got != 9 {
		t.Fatalf("history = %d messages after 3 decisions, want 9", got)
	}
	// The system prompt and tool list must be byte-identical across the game --
	// they are the cached prefix, and a mutation costs a full rebuild per turn.
	first, last := st.calls[0], st.last()
	if first.System[0].Text != last.System[0].Text {
		t.Fatal("system prompt changed mid-game")
	}
	if len(first.Tools) != len(last.Tools) {
		t.Fatal("tool list changed mid-game")
	}
	for i := range first.Tools {
		if first.Tools[i].OfTool.Name != last.Tools[i].OfTool.Name {
			t.Fatalf("tool %d changed: %s -> %s", i,
				first.Tools[i].OfTool.Name, last.Tools[i].OfTool.Name)
		}
	}
}

func TestInfoOnlyToolsDoNotCountAsFailures(t *testing.T) {
	moves := testMoves()
	st := &stubTransport{respond: func(n int, _ anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		// Three chat-only turns, then a real decision. A chatty model must not
		// talk itself into autopilot (upstream's INFO_ONLY_TOOLS).
		if n < 3 {
			return toolUseResponse("tu"+strconv.Itoa(n), ToolSendChatMessage,
				map[string]any{"message": "thinking out loud"}), nil
		}
		return toolUseResponse("tuX", ToolChooseAction, map[string]any{"choice": "m1"}), nil
	}}
	p := New(st, Options{Seat: "bot-a", MaxSteps: 6})
	got, err := p.Pick(context.Background(), testView(), moves)
	if err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if got.Label != moves[0].Label {
		t.Fatalf("picked %q, want %q", got.Label, moves[0].Label)
	}
	if p.rec.Autopilot() {
		t.Fatal("chatting degraded the seat to autopilot")
	}
}

func TestContextResetResendsTheDecision(t *testing.T) {
	moves := testMoves()
	var promptAfterReset string
	st := &stubTransport{respond: func(n int, p anthropic.BetaMessageNewParams) (*anthropic.BetaMessage, error) {
		// Truncate repeatedly to force a reset, then look at what the model is
		// handed next: a reset that drops the board would leave it answering
		// blind.
		if n < MaxConsecutiveTruncations {
			return textResponse("...", anthropic.BetaStopReasonMaxTokens), nil
		}
		promptAfterReset = lastUserText(p)
		return toolUseResponse("tu1", ToolChooseAction, map[string]any{"choice": "m1"}), nil
	}}
	pol := New(st, Options{Seat: "bot-a", MaxSteps: 8})
	if _, err := pol.Pick(context.Background(), testView(), moves); err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if !strings.Contains(promptAfterReset, "Choices (") || !strings.Contains(promptAfterReset, "[id=m1") {
		t.Fatalf("post-reset prompt has no decision:\n%s", promptAfterReset)
	}
}
