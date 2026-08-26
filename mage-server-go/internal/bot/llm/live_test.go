package llm

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/magefree/mage-server-go/internal/bot"
)

// live_test.go is the ONLY test in this package that would talk to the API, and
// it skips itself unless ANTHROPIC_API_KEY is set.
//
// It skips rather than fails on purpose: a missing credential is not a broken
// build. Everything the loop decides -- windowing, breakpoint placement, tool
// dispatch, choice resolution, the recovery matrix -- is decided on our side of
// the wire and is covered against the stub transport. What only a live call can
// confirm is the wire contract itself: that the request is accepted, that
// Strict schemas are honoured, and above all that CacheReadInputTokens is
// NON-ZERO on the second request. A zero there means something in the prefix is
// changing between turns, costs are roughly 10x, and nothing else reports it
// (plan Phase 5 verification).
func TestLiveTwoTurnCacheHit(t *testing.T) {
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" {
		t.Skip("ANTHROPIC_API_KEY not set; skipping the only test that makes a live call")
	}

	p := New(NewSDKTransport(key), Options{Seat: "bot-a", Oracle: simOracle()})
	v := testView()
	moves := bot.LegalMoves(v)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	for i := 0; i < 2; i++ {
		m, err := p.Pick(ctx, v, moves)
		if err != nil {
			t.Fatalf("Pick %d: %v", i, err)
		}
		t.Logf("turn %d picked %q", i, m.Label)
	}

	u := p.Usage()
	t.Logf("usage: requests=%d in=%d out=%d cache_create=%d cache_read=%d",
		u.Requests, u.InputTokens, u.OutputTokens, u.CacheCreationTokens, u.CacheReadTokens)
	if u.CacheReadTokens == 0 {
		t.Fatal("CacheReadInputTokens is 0 on the second request: the prompt prefix " +
			"is being invalidated between turns. Find the invalidator (a timestamp, " +
			"an unsorted map, a mutated tool list) before running a real game.")
	}
}
