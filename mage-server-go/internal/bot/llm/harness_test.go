package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/magefree/mage-server-go/internal/bot"
	"github.com/magefree/mage-server-go/internal/game"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// harness_test.go is the interchangeability proof: the SAME headless game, the
// same runner, the same assertions, played once by RandomPolicy and once by
// LLMPolicy over a stubbed transport.
//
// This is the check that says the Phase 5 seam actually held. LLMPolicy is a
// bot.Policy and nothing else -- no extra hook into the runner, no special case
// in the seat loop -- so if the LLM run needs a different harness, the
// abstraction leaked and the failure would surface as a mysterious difference
// in a real game instead of here.
//
// The stub reads the prompt the policy built and answers it, which also makes
// it a live check on the rendering: if the option list ever stops carrying
// "Choices (N)" and "[id=mK", the stub cannot answer and the game stalls.

var choicesRE = regexp.MustCompile(`Choices \((\d+)\)`)

// promptDrivenStub answers each decision by picking a uniformly random option
// out of the prompt it was sent -- which makes it behaviourally a RandomPolicy
// wearing the whole LLM code path: prompt rendering, tool dispatch, choice
// resolution, conversation growth and cache annotation all run for real.
type promptDrivenStub struct {
	mu  sync.Mutex
	rng *rand.Rand
	n   int
}

func (s *promptDrivenStub) New(_ context.Context, params anthropic.BetaMessageNewParams, _ ...option.RequestOption) (*anthropic.BetaMessage, error) {
	text := lastUserText(params)
	m := choicesRE.FindStringSubmatch(text)
	if m == nil {
		return nil, fmt.Errorf("prompt carried no Choices line: %q", text)
	}
	count, err := strconv.Atoi(m[1])
	if err != nil || count <= 0 {
		return nil, fmt.Errorf("unusable choice count %q", m[1])
	}

	s.mu.Lock()
	pick := s.rng.Intn(count) + 1
	s.n++
	id := fmt.Sprintf("tu%d", s.n)
	s.mu.Unlock()

	raw, _ := json.Marshal(map[string]any{"choice": fmt.Sprintf("m%d", pick)})
	return &anthropic.BetaMessage{
		Role:       "assistant",
		StopReason: anthropic.BetaStopReasonToolUse,
		Content: []anthropic.BetaContentBlockUnion{{
			Type: "tool_use", ID: id, Name: ToolChooseAction, Input: json.RawMessage(raw),
		}},
	}, nil
}

func simOracle() bot.MapOracle {
	return bot.MapOracle{
		"Forest":   {Name: "Forest", TypeLine: "Basic Land — Forest", OracleText: "({T}: Add {G}.)"},
		"Mountain": {Name: "Mountain", TypeLine: "Basic Land — Mountain", OracleText: "({T}: Add {R}.)"},
		"Llanowar Elves": {
			Name: "Llanowar Elves", ManaCost: "{G}", TypeLine: "Creature — Elf Druid",
			OracleText: "{T}: Add {G}.", Power: "1", Toughness: "1",
		},
		"Grizzly Bears": {
			Name: "Grizzly Bears", ManaCost: "{1}{G}", TypeLine: "Creature — Bear",
			Power: "2", Toughness: "2",
		},
		"Lightning Bolt": {
			Name: "Lightning Bolt", ManaCost: "{R}", TypeLine: "Instant",
			OracleText: "Lightning Bolt deals 3 damage to any target.",
		},
		"Marath, Will of the Wild": {
			Name: "Marath, Will of the Wild", ManaCost: "{R}{G}{W}",
			TypeLine: "Legendary Creature — Elemental Beast", Power: "0", Toughness: "0",
		},
	}
}

func simDeck(seat int) game.DeckList {
	lands := []string{"Forest", "Mountain"}
	spells := []string{"Llanowar Elves", "Grizzly Bears", "Lightning Bolt"}
	main := make([]string, 0, 99)
	for i := 0; i < 99; i++ {
		if i%5 < 2 {
			main = append(main, lands[(i+seat)%len(lands)])
		} else {
			main = append(main, spells[(i+seat)%len(spells)])
		}
	}
	return game.DeckList{MainDeck: main, Commanders: []string{"Marath, Will of the Wild"}}
}

// runHarness plays one headless game with policies built by newPolicy.
func runHarness(t *testing.T, seed int64, newPolicy func(seat string, i int) bot.Policy) (bool, bot.Stats) {
	t.Helper()
	logger := zaptest.NewLogger(t, zaptest.Level(zap.ErrorLevel))

	seats := []string{"bot-a", "bot-b", "bot-c", "bot-d"}
	engine := game.NewGameEngine(logger)
	adapter := game.NewEngineAdapter(engine, logger)
	mgr := game.NewManager(logger)
	g := mgr.CreateGame(fmt.Sprintf("table-%d", seed), "Commander Free For All", seats)
	decks := make(map[string]game.DeckList, len(seats))
	for i, s := range seats {
		decks[s] = simDeck(i)
	}
	if err := adapter.StartGameWithDecks(g, decks); err != nil {
		t.Fatalf("StartGameWithDecks: %v", err)
	}
	go adapter.ProcessGameActions(g)
	defer mgr.RemoveGame(g.ID)

	runner := bot.NewBotRunner(bot.RunnerConfig{
		GameID:   g.ID,
		Actions:  mgr,
		Views:    adapter,
		Oracle:   simOracle(),
		Pacing:   bot.Pacing{}, // zero: no delays, this is a headless sim
		MaxTurns: 200,
		Logger:   logger,
	})
	for i, s := range seats {
		runner.AddSeat(s, newPolicy(s, i))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	done, err := runner.Run(ctx)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	return done, runner.Stats()
}

// promptDrivenGemini is the same random-answer stub behind the
// OpenAI-compatible wire format instead of the SDK one: an httptest server that
// reads the rendered prompt out of the JSON body and answers with a tool call.
//
// It exists so the interchangeability proof covers the SECOND PROVIDER TOO. The
// seam this phase added is below Policy, so the claim being tested is that a
// whole game plays identically through either Completer -- which is only worth
// anything if a real game is actually played through both.
func newPromptDrivenGemini(t *testing.T, seed int64) string {
	t.Helper()
	var mu sync.Mutex
	rng := rand.New(rand.NewSource(seed))
	n := 0

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req openAIRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		text := lastUserContent(req)
		m := choicesRE.FindStringSubmatch(text)
		if m == nil {
			http.Error(w, "prompt carried no Choices line", http.StatusBadRequest)
			return
		}
		count, err := strconv.Atoi(m[1])
		if err != nil || count <= 0 {
			http.Error(w, "unusable choice count", http.StatusBadRequest)
			return
		}
		mu.Lock()
		pick := rng.Intn(count) + 1
		n++
		id := fmt.Sprintf("tu%d", n)
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, toolCallBody(id, ToolChooseAction,
			map[string]any{"choice": fmt.Sprintf("m%d", pick)}))
	}))
	t.Cleanup(srv.Close)
	return srv.URL
}

func TestPolicyImplementationsAreInterchangeable(t *testing.T) {
	if testing.Short() {
		t.Skip("headless game, skipped under -short")
	}
	const seed = 7

	randomDone, randomStats := runHarness(t, seed, func(_ string, i int) bot.Policy {
		return bot.NewRandomPolicy(seed*1000 + int64(i))
	})
	llmDone, llmStats := runHarness(t, seed, func(seat string, i int) bot.Policy {
		return New(&promptDrivenStub{rng: rand.New(rand.NewSource(seed*1000 + int64(i)))},
			Options{Seat: seat, Oracle: simOracle()})
	})
	geminiDone, geminiStats := runHarness(t, seed, func(seat string, i int) bot.Policy {
		p, err := NewPolicy(Options{
			Seat:     seat,
			Provider: ProviderGemini,
			APIKey:   "test-key",
			BaseURL:  newPromptDrivenGemini(t, seed*1000+int64(i)),
			Oracle:   simOracle(),
		})
		if err != nil {
			t.Fatalf("NewPolicy: %v", err)
		}
		return p
	})

	t.Logf("random:    completed=%v turns=%d macros=%d failed=%d",
		randomDone, randomStats.Turns, randomStats.MacrosExecuted, randomStats.MacrosFailed)
	t.Logf("anthropic: completed=%v turns=%d macros=%d failed=%d",
		llmDone, llmStats.Turns, llmStats.MacrosExecuted, llmStats.MacrosFailed)
	t.Logf("gemini:    completed=%v turns=%d macros=%d failed=%d",
		geminiDone, geminiStats.Turns, geminiStats.MacrosExecuted, geminiStats.MacrosFailed)

	for _, tc := range []struct {
		name  string
		done  bool
		stats bot.Stats
	}{
		{"RandomPolicy", randomDone, randomStats},
		{"LLMPolicy/anthropic", llmDone, llmStats},
		{"LLMPolicy/gemini", geminiDone, geminiStats},
	} {
		if !tc.done {
			t.Errorf("%s: game did not reach a terminal state", tc.name)
		}
		if tc.stats.MacrosExecuted == 0 {
			t.Errorf("%s: no macros executed", tc.name)
		}
		if tc.stats.MacrosFailed != 0 {
			t.Errorf("%s: %d macros failed to land: %v",
				tc.name, tc.stats.MacrosFailed, tc.stats.FailedMacros)
		}
	}
}
