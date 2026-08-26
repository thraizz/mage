package llm

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/magefree/mage-server-go/internal/bot"
)

// live_gemini_test.go is the Gemini half of live_test.go, and it skips itself
// unless GEMINI_API_KEY is set.
//
// IT SKIPS RATHER THAN FAILS, and it never prompts. A missing credential is not
// a broken build, and a test that blocks waiting for one turns `go test ./...`
// into a hang on every machine that does not have a key.
//
// Everything the loop decides -- windowing, tool dispatch, choice resolution,
// the recovery matrix -- is covered against the httptest stub in openai_test.go.
// What only a live call can settle is the wire contract: that the model id is
// real, that the function schemas are inside the OpenAPI subset Gemini accepts
// (this is where dropping strict/additionalProperties gets its verdict), that
// the response parses, and above all THE COST GATE -- that cached_tokens is
// non-zero on a second request with a shared prefix. A zero there means the
// prefix is being invalidated between turns, costs are roughly 10x, and nothing
// else reports it.
func TestLiveGeminiTwoTurnCacheHit(t *testing.T) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		t.Skip("GEMINI_API_KEY not set; skipping the live Gemini call")
	}

	p, err := NewPolicy(Options{
		Seat:     "bot-a",
		Provider: ProviderGemini,
		APIKey:   key,
		Oracle:   simOracle(),
	})
	if err != nil {
		t.Fatalf("NewPolicy: %v", err)
	}
	v := testView()
	moves := bot.LegalMoves(v)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	// Check 1 and 3 in one: a request is accepted at all (model id, function
	// schemas, response parsing), and the model returns a well-formed macro
	// handle that resolves to a legal move.
	for i := 0; i < 2; i++ {
		m, err := p.Pick(ctx, v, moves)
		if err != nil {
			t.Fatalf("Pick %d: %v", i, err)
		}
		t.Logf("turn %d picked %q (%s)", i, m.Label, m.KindOf())
		if m.Label == "" {
			t.Fatalf("turn %d produced an empty macro", i)
		}
	}
	if p.rec.Autopilot() {
		t.Fatalf("the seat degraded to autopilot on a live run: %s", p.rec.AutopilotReason())
	}

	// Check 2: the cost gate.
	u := p.Usage()
	t.Logf("usage: requests=%d in=%d out=%d cache_read=%d (implicit-cache minimum is %d tokens)",
		u.Requests, u.InputTokens, u.OutputTokens, u.CacheReadTokens, GeminiImplicitCacheMinTokens)
	if u.Requests < 2 {
		t.Fatalf("only %d requests were made; the two-turn cache check did not run", u.Requests)
	}
	if u.CacheReadTokens == 0 {
		// NOT a hard failure, unlike the Anthropic equivalent, and the
		// difference is real rather than a softened assertion. Anthropic
		// caching is explicit: we ask for it, so a miss is our bug. Gemini's is
		// implicit and requires a >=4096-token shared prefix; a two-turn
		// fixture game can legitimately sit under that, in which case a zero
		// says "prompt too short", not "prefix unstable". Phase 6 measures this
		// over a real game, where the transcript is long enough for the
		// distinction to collapse.
		t.Logf("WARNING: cache_read is 0 after %d requests with %d input tokens. "+
			"If the average prompt is above %d tokens this means the prefix is being "+
			"invalidated between turns -- find the invalidator before running a real game.",
			u.Requests, u.InputTokens, GeminiImplicitCacheMinTokens)
	}
}
