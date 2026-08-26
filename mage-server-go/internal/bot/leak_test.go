package bot

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/magefree/mage-server-go/internal/game"
)

// leak_test.go is the regression test that actually matters. Everything else in
// this package is a formatting concern; this file is the reason Redact exists.
//
// A bot that can see its own library draws perfectly and every result the
// harness ever produces is worthless -- and, crucially, the games still LOOK
// fine. There is no crash, no error log, no failing assertion anywhere else in
// the system. This test is the only thing standing between a refactor and a
// silently invalid benchmark.

// legitimateKnowledge returns every card name the viewing seat is entitled to
// know about: its own hand, the public zones, and any revealed card.
// Deliberately built from the ENGINE view, not from the SafeView, so that a
// SafeView which grew an extra zone fails instead of widening the whitelist.
func legitimateKnowledge(v *game.PlaytestGameView, viewerID string) map[string]struct{} {
	known := make(map[string]struct{})
	add := func(cards []*game.Card) {
		for _, c := range cards {
			if c != nil {
				known[c.Name] = struct{}{}
			}
		}
	}
	add(v.Battlefield)
	add(v.Exile)
	add(v.Stack)
	add(v.Command)
	if v.Me != nil && v.Me.PlayerID == viewerID {
		add(v.Me.Hand) // own hand: legitimate
		add(v.Me.Graveyard)
	}
	for _, o := range v.Opponents {
		add(o.Graveyard)
		if o.RevealedTopCard && o.TopCard != nil {
			add([]*game.Card{o.TopCard}) // face-up top card: legitimate
		}
	}
	return known
}

// collectSafeViewCardNames reflectively walks a SafeView and returns every card
// name reachable from it, with the field path where each was found.
//
// Reflection rather than a hand-written walk on purpose: a hand-written walk
// only checks the fields someone remembered to list, so re-adding Library to
// SafePlayerView would sail straight past it. This catches any new *SafeCard
// anywhere in the graph, including one nobody told the test about.
func collectSafeViewCardNames(v *SafeView) map[string][]string {
	found := make(map[string][]string)
	var walk func(val reflect.Value, path string)
	walk = func(val reflect.Value, path string) {
		switch val.Kind() {
		case reflect.Ptr, reflect.Interface:
			if val.IsNil() {
				return
			}
			if val.Kind() == reflect.Ptr && val.Type().Elem() == reflect.TypeOf(SafeCard{}) {
				card := val.Interface().(*SafeCard)
				found[card.Name] = append(found[card.Name], path)
				return
			}
			walk(val.Elem(), path)
		case reflect.Slice, reflect.Array:
			for i := 0; i < val.Len(); i++ {
				walk(val.Index(i), fmt.Sprintf("%s[%d]", path, i))
			}
		case reflect.Struct:
			if val.Type() == reflect.TypeOf(SafeCard{}) {
				found[val.Interface().(SafeCard).Name] = append(found[val.Interface().(SafeCard).Name], path)
				return
			}
			t := val.Type()
			for i := 0; i < val.NumField(); i++ {
				if t.Field(i).PkgPath != "" {
					continue // unexported
				}
				walk(val.Field(i), path+"."+t.Field(i).Name)
			}
		case reflect.Map:
			for _, k := range val.MapKeys() {
				walk(val.MapIndex(k), fmt.Sprintf("%s[%v]", path, k.Interface()))
			}
		}
	}
	walk(reflect.ValueOf(v), "SafeView")
	return found
}

// TestRedactLeaksNothingIntoSafeView is the structural half of the guard: no
// card outside the seat's knowledge set may be REACHABLE from the SafeView,
// whether or not the serializer happens to render it today.
func TestRedactLeaksNothingIntoSafeView(t *testing.T) {
	view := baseView()
	known := legitimateKnowledge(view, seatAlice)

	// Sanity: the fixture actually has something to leak. Without this, a
	// fixture that quietly lost its library would make the test vacuous.
	require.NotEmpty(t, view.Me.Library, "fixture must plant secrets to be worth testing")
	for _, n := range secretLibraryNames {
		require.NotContains(t, known, n, "fixture secret %q must not be legitimate knowledge", n)
	}

	safe := Redact(view, seatAlice)

	for name, paths := range collectSafeViewCardNames(safe) {
		if _, ok := known[name]; !ok {
			t.Errorf("SafeView leaks card %q, reachable at %v -- not in the viewer's knowledge set", name, paths)
		}
	}
}

// TestSerializedOutputLeaksNothing is the textual half: whatever ends up in
// front of the model must not name a card the seat cannot know about.
func TestSerializedOutputLeaksNothing(t *testing.T) {
	view := baseView()
	known := legitimateKnowledge(view, seatAlice)

	safe := Redact(view, seatAlice)
	s := NewSerializer(testOracle())
	out := s.Render(context.Background(), safe, &Decision{
		Index: 1, Turn: 3, Phase: "PRECOMBAT_MAIN", Player: "Alice",
		Message: "Play a land or cast a spell",
		Choices: []Choice{{Name: "Forest", ID: "p2", Action: "play land"}},
	})

	for _, secret := range secretLibraryNames {
		require.NotContains(t, out, secret,
			"rendered prompt names %q, which is in the viewer's own library and must never be shown", secret)
	}
	for _, secret := range opponentSecretNames {
		require.NotContains(t, out, secret,
			"rendered prompt names %q, which is hidden opponent information", secret)
	}

	// And the general form: every card name the oracle knows about that shows
	// up in the prompt must be legitimate knowledge.
	for name := range testOracle() {
		if _, ok := known[name]; ok {
			continue
		}
		require.NotContains(t, out, name, "rendered prompt names %q, outside the viewer's knowledge set", name)
	}

	require.NotEmpty(t, out)
	t.Logf("rendered prompt:\n%s", out)
}

// TestRedactHardFailsOnLeakyOpponentView proves the assertion is an assertion:
// a view whose opponent hand or library is populated is rejected, not silently
// passed through.
func TestRedactHardFailsOnLeakyOpponentView(t *testing.T) {
	t.Run("opponent hand", func(t *testing.T) {
		view := baseView()
		view.Opponents[0].Hand = []*game.Card{card("leak-0", opponentSecretNames[0], seatBob)}
		_, err := RedactErr(view, seatAlice)
		require.Error(t, err)
		require.Contains(t, err.Error(), "hand is not hidden")
		require.Panics(t, func() { Redact(view, seatAlice) })
	})

	t.Run("opponent library", func(t *testing.T) {
		view := baseView()
		view.Opponents[0].Library = []*game.Card{card("leak-1", opponentSecretNames[1], seatBob)}
		_, err := RedactErr(view, seatAlice)
		require.Error(t, err)
		require.Contains(t, err.Error(), "library is not hidden")
	})

	t.Run("wrong seat", func(t *testing.T) {
		view := baseView()
		_, err := RedactErr(view, seatBob)
		require.Error(t, err)
		require.Contains(t, err.Error(), "may only be redacted for the seat it was built for")
	})

	t.Run("unrevealed top card is dropped", func(t *testing.T) {
		view := baseView()
		view.Opponents[0].RevealedTopCard = false
		view.Opponents[0].TopCard = card("top-0", opponentSecretNames[0], seatBob)
		safe, err := RedactErr(view, seatAlice)
		require.NoError(t, err)
		require.Nil(t, safe.Opponents[0].TopCard)
	})
}

// TestRedactSharesNoPointers is anti-pattern 3 from the plan: buildGameView
// hands out the engine's own slices. Mutating the source after redaction must
// not be observable through the SafeView -- and, symmetrically, a bot writing
// through its SafeView must not be able to reach into the running game.
func TestRedactSharesNoPointers(t *testing.T) {
	view := baseView()
	safe := Redact(view, seatAlice)

	before := snapshotStrings(safe)

	// Mutate every zone buildGameView shares by pointer (view.go:74-81, 102, 119-120).
	view.Battlefield[0].Name = "MUTATED-BATTLEFIELD"
	view.Battlefield[0].DisplayName = "MUTATED-BATTLEFIELD"
	view.Battlefield[0].Tapped = !view.Battlefield[0].Tapped
	view.Battlefield[0].Counters = append(view.Battlefield[0].Counters, game.Counter{Name: "MUTATED", Count: 9})
	view.Battlefield = append(view.Battlefield, card("bf-x", "MUTATED-APPENDED", seatAlice))
	view.Exile[0].Name = "MUTATED-EXILE"
	view.Me.Hand[0].Name = "MUTATED-HAND"
	view.Me.Graveyard = append(view.Me.Graveyard, card("gy-x", "MUTATED-GY", seatAlice))
	view.Me.ManaPool.Red = 99
	view.Opponents[0].Graveyard[0].Name = "MUTATED-OPP-GY"
	view.Opponents[0].ManaPool.Blue = 99
	view.Log = append(view.Log, game.LogEntry{Kind: "x", Message: "MUTATED-LOG"})

	require.Equal(t, before, snapshotStrings(safe),
		"SafeView changed when the source view was mutated -- it is aliasing live engine state")

	// The reverse direction: writing through the SafeView must not reach the engine.
	safe.Battlefield[0].Name = "BOT-WROTE-THIS"
	safe.Me.Hand[0].Name = "BOT-WROTE-THIS"
	safe.Me.ManaPool.Green = 42
	for _, c := range view.Battlefield {
		require.NotEqual(t, "BOT-WROTE-THIS", c.Name)
	}
	for _, c := range view.Me.Hand {
		require.NotEqual(t, "BOT-WROTE-THIS", c.Name)
	}
	require.Equal(t, 0, view.Me.ManaPool.Green)
}

// snapshotStrings renders every card name and mana pool in the SafeView into a
// single comparable string.
func snapshotStrings(v *SafeView) string {
	var b strings.Builder
	names := collectSafeViewCardNames(v)
	keys := make([]string, 0, len(names))
	for k := range names {
		keys = append(keys, k)
	}
	// Deterministic order.
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[j] < keys[i] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	for _, k := range keys {
		fmt.Fprintf(&b, "%s=%v;", k, names[k])
	}
	fmt.Fprintf(&b, "pool=%+v;", v.Me.ManaPool)
	for _, o := range v.Opponents {
		fmt.Fprintf(&b, "opool=%+v;", o.ManaPool)
	}
	fmt.Fprintf(&b, "log=%d;", len(v.Log))
	for _, c := range v.Battlefield {
		fmt.Fprintf(&b, "bf:%s:%v:%v;", c.ID, c.Tapped, c.Counters)
	}
	return b.String()
}
