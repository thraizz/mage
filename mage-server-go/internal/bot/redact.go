// Package bot turns the engine's per-player game view into the token-efficient
// board text that mage-bench feeds to an LLM pilot.
//
// The package is deliberately split into two halves:
//
//	redact.go     the security boundary: engine view -> SafeView (deep copy, hidden info stripped)
//	serialize.go  the presentation layer: SafeView + Decision -> "## Decision" text
//
// Nothing downstream of Redact ever sees a *game.PlaytestGameView. That type
// boundary IS the security boundary: SafeView is a distinct type, not an alias
// or an embedding, so it is impossible to "accidentally" hand a bot the live
// engine view.
package bot

import (
	"fmt"

	"github.com/magefree/mage-server-go/internal/game"
)

// SafeView is the redacted, fully-owned copy of a player's game view.
//
// It is a NEW type, not an alias of *game.PlaytestGameView, on purpose. Two
// invariants hold for every SafeView produced by Redact:
//
//  1. Nothing in it aliases live engine state. game.GameEngine.buildGameView
//     (internal/game/view.go:74-81, 102, 119-120) shares Battlefield, Exile,
//     Stack, Command, graveyards, hands and mana pools *by pointer* with the
//     authoritative GameState. A consumer holding that view can mutate the
//     running game through those slices. protojson saves the websocket clients
//     from this because it copies on marshal; a bot in-process has no such luck.
//     Redact therefore deep-copies every slice and every *game.Card.
//  2. It contains no information the viewer is not entitled to. In particular
//     the viewer's own library is dropped entirely -- see SafePlayerView.
type SafeView struct {
	GameID            string
	ViewerID          string
	ActiveControlSeat string
	Me                *SafePlayerView
	Opponents         []*SafeOpponentView
	Battlefield       []*SafeCard
	Exile             []*SafeCard
	Stack             []*SafeCard
	Command           []*SafeCard
	Turn              int
	ActivePlayerID    string
	IsInitialized     bool
	Log               []SafeLogEntry
	MulliganType      string
	FreeMulligans     int
}

// SafePlayerView is the viewing seat. Note what is absent.
type SafePlayerView struct {
	PlayerID     string
	Name         string
	Life         int
	Poison       int
	Energy       int
	LibraryCount int
	HandCount    int
	Hand         []*SafeCard
	Graveyard    []*SafeCard
	ManaPool     SafeManaPool

	// Library is DELIBERATELY ABSENT.
	//
	// game.PlaytestPlayerView.Library (internal/game/view.go:41) carries the
	// viewer's own library in full, in order, because the web client needs it to
	// render search / scry / "look at the top N" UI. A human sitting in front of
	// that UI only opens it when the game tells them to. A bot handed the same
	// slice reads the top card of its deck on every single decision and plays a
	// game no human could play -- it is a perfect-information cheat that would
	// silently invalidate every result the bots ever produce.
	//
	// LibraryCount above is the only library information a bot gets.
	//
	// When scry / surveil / tutor decisions are wired up (they are not in this
	// phase), the fix is a narrow, explicitly-populated Revealed []*SafeCard
	// slice scoped to that one decision -- never this field coming back.

	KeptHand        bool
	MulliganCount   int
	RevealedTopCard bool
}

// SafeOpponentView is an opposing seat: counts, public zones, and nothing else.
type SafeOpponentView struct {
	PlayerID        string
	Name            string
	Life            int
	Poison          int
	Energy          int
	LibraryCount    int
	HandCount       int
	TopCard         *SafeCard // non-nil only when that player's top card is face up
	Graveyard       []*SafeCard
	ManaPool        SafeManaPool
	KeptHand        bool
	MulliganCount   int
	RevealedTopCard bool
}

// SafeCard is a value-owned copy of game.Card.
type SafeCard struct {
	ID           string
	Name         string
	DisplayName  string
	OwnerID      string
	ControllerID string

	ManaCost   string
	Type       string
	SubTypes   string
	SuperTypes string
	Color      string
	Power      string
	Toughness  string
	Loyalty    string
	RulesText  string

	Zone        string
	Tapped      bool
	Flipped     bool
	Transformed bool
	FaceDown    bool
	Counters    []SafeCounter

	Attacking         bool
	Blocking          bool
	SummoningSickness bool

	AttachedTo []string

	// OriginalCard and IsCopy back the "copy of X" / "copy" annotations the
	// mage-bench permanent display emits (reference/decision_renderer.py:333-337).
	// The local engine has no copy mechanic -- game.Card has no equivalent field
	// -- so Redact never populates these and they are always zero today. They
	// exist so the renderer implements the full documented format rather than a
	// subset that quietly loses information the day copies are wired up.
	OriginalCard string
	IsCopy       bool
}

// SafeCounter mirrors game.Counter.
type SafeCounter struct {
	Name  string
	Count int
}

// SafeManaPool mirrors game.ManaPool. It is stored by value rather than by
// pointer specifically so that it cannot alias the live pool.
type SafeManaPool struct {
	White     int
	Blue      int
	Black     int
	Red       int
	Green     int
	Colorless int
}

// Total reports the number of mana in the pool.
func (p SafeManaPool) Total() int {
	return p.White + p.Blue + p.Black + p.Red + p.Green + p.Colorless
}

// LeakError reports that an engine view handed to Redact contained information
// the viewer must not see. It always indicates a bug in the engine's view
// builder (or that Redact was called with the wrong viewer), never bad input
// from a player, and it is never recoverable by retrying.
type LeakError struct {
	ViewerID string
	Detail   string
}

func (e *LeakError) Error() string {
	return fmt.Sprintf("bot: refusing to build SafeView for viewer %q: %s", e.ViewerID, e.Detail)
}

// SafeLogEntry mirrors game.LogEntry without the timestamp, which is volatile
// and would poison golden comparisons.
type SafeLogEntry struct {
	Kind    string
	Message string
}

// Redact deep-copies v into a SafeView for viewerID, dropping everything the
// viewer is not entitled to see.
//
// It PANICS if the input view leaks hidden information (an opponent hand or
// library that is not empty, or a view built for a different seat). That is
// deliberate: a leak here is an engine bug that silently invalidates every game
// the bots play, and a returned error would eventually get logged and ignored.
// Callers that need to survive the condition use RedactErr.
func Redact(v *game.PlaytestGameView, viewerID string) *SafeView {
	sv, err := RedactErr(v, viewerID)
	if err != nil {
		panic(err)
	}
	return sv
}

// RedactErr is Redact with the invariant violations returned as a *LeakError
// instead of panicking. It never returns a partially-built SafeView alongside
// an error -- on any violation the result is nil.
func RedactErr(v *game.PlaytestGameView, viewerID string) (*SafeView, error) {
	if v == nil {
		return nil, &LeakError{ViewerID: viewerID, Detail: "nil game view"}
	}
	if v.ViewerID != viewerID {
		return nil, &LeakError{
			ViewerID: viewerID,
			Detail:   fmt.Sprintf("view was built for viewer %q; a view may only be redacted for the seat it was built for", v.ViewerID),
		}
	}
	if v.Me == nil {
		return nil, &LeakError{ViewerID: viewerID, Detail: "view has no Me seat"}
	}
	if v.Me.PlayerID != viewerID {
		return nil, &LeakError{
			ViewerID: viewerID,
			Detail:   fmt.Sprintf("view.Me is seat %q", v.Me.PlayerID),
		}
	}

	out := &SafeView{
		GameID:            v.GameID,
		ViewerID:          v.ViewerID,
		ActiveControlSeat: v.ActiveControlSeat,
		Battlefield:       copyCards(v.Battlefield),
		Exile:             copyCards(v.Exile),
		Stack:             copyCards(v.Stack),
		Command:           copyCards(v.Command),
		Turn:              v.Turn,
		ActivePlayerID:    v.ActivePlayerID,
		IsInitialized:     v.IsInitialized,
		MulliganType:      v.MulliganType,
		FreeMulligans:     v.FreeMulligans,
	}

	if len(v.Log) > 0 {
		out.Log = make([]SafeLogEntry, 0, len(v.Log))
		for _, e := range v.Log {
			out.Log = append(out.Log, SafeLogEntry{Kind: e.Kind, Message: e.Message})
		}
	}

	out.Me = &SafePlayerView{
		PlayerID:     v.Me.PlayerID,
		Name:         v.Me.Name,
		Life:         v.Me.Life,
		Poison:       v.Me.Poison,
		Energy:       v.Me.Energy,
		LibraryCount: v.Me.LibraryCount,
		HandCount:    v.Me.HandCount,
		Hand:         copyCards(v.Me.Hand),
		Graveyard:    copyCards(v.Me.Graveyard),
		ManaPool:     copyManaPool(v.Me.ManaPool),
		// v.Me.Library is intentionally not copied. See SafePlayerView.
		KeptHand:        v.Me.KeptHand,
		MulliganCount:   v.Me.MulliganCount,
		RevealedTopCard: v.Me.RevealedTopCard,
	}

	out.Opponents = make([]*SafeOpponentView, 0, len(v.Opponents))
	for _, o := range v.Opponents {
		if o == nil {
			return nil, &LeakError{ViewerID: viewerID, Detail: "nil opponent view"}
		}
		// Assert, do not trust. buildGameView sets these to empty slices
		// (internal/game/view.go:117-118); if that ever regresses -- a new zone
		// helper, a refactor that copies the Player struct wholesale -- the bot
		// would read every opponent's hand and library and nothing else in the
		// system would notice.
		if len(o.Hand) != 0 {
			return nil, &LeakError{
				ViewerID: viewerID,
				Detail:   fmt.Sprintf("opponent %q hand is not hidden: %d card(s) present", o.PlayerID, len(o.Hand)),
			}
		}
		if len(o.Library) != 0 {
			return nil, &LeakError{
				ViewerID: viewerID,
				Detail:   fmt.Sprintf("opponent %q library is not hidden: %d card(s) present", o.PlayerID, len(o.Library)),
			}
		}
		// TopCard is legitimate knowledge only while that player's top card is
		// actually face up. Drop it otherwise rather than trusting the producer.
		var top *SafeCard
		if o.RevealedTopCard {
			top = copyCard(o.TopCard)
		}
		out.Opponents = append(out.Opponents, &SafeOpponentView{
			PlayerID:        o.PlayerID,
			Name:            o.Name,
			Life:            o.Life,
			Poison:          o.Poison,
			Energy:          o.Energy,
			LibraryCount:    o.LibraryCount,
			HandCount:       o.HandCount,
			TopCard:         top,
			Graveyard:       copyCards(o.Graveyard),
			ManaPool:        copyManaPool(o.ManaPool),
			KeptHand:        o.KeptHand,
			MulliganCount:   o.MulliganCount,
			RevealedTopCard: o.RevealedTopCard,
		})
	}

	return out, nil
}

func copyManaPool(p *game.ManaPool) SafeManaPool {
	if p == nil {
		return SafeManaPool{}
	}
	return SafeManaPool{
		White:     p.White,
		Blue:      p.Blue,
		Black:     p.Black,
		Red:       p.Red,
		Green:     p.Green,
		Colorless: p.Colorless,
	}
}

func copyCards(in []*game.Card) []*SafeCard {
	if in == nil {
		return nil
	}
	out := make([]*SafeCard, 0, len(in))
	for _, c := range in {
		if c == nil {
			continue
		}
		out = append(out, copyCard(c))
	}
	return out
}

func copyCard(c *game.Card) *SafeCard {
	if c == nil {
		return nil
	}
	out := &SafeCard{
		ID:                c.ID,
		Name:              c.Name,
		DisplayName:       c.DisplayName,
		OwnerID:           c.OwnerID,
		ControllerID:      c.ControllerID,
		ManaCost:          c.ManaCost,
		Type:              c.Type,
		SubTypes:          c.SubTypes,
		SuperTypes:        c.SuperTypes,
		Color:             c.Color,
		Power:             c.Power,
		Toughness:         c.Toughness,
		Loyalty:           c.Loyalty,
		RulesText:         c.RulesText,
		Zone:              c.Zone,
		Tapped:            c.Tapped,
		Flipped:           c.Flipped,
		Transformed:       c.Transformed,
		FaceDown:          c.FaceDown,
		Attacking:         c.Attacking,
		Blocking:          c.Blocking,
		SummoningSickness: c.SummoningSickness,
	}
	if c.Counters != nil {
		out.Counters = make([]SafeCounter, 0, len(c.Counters))
		for _, ctr := range c.Counters {
			out.Counters = append(out.Counters, SafeCounter{Name: ctr.Name, Count: ctr.Count})
		}
	}
	if c.AttachedTo != nil {
		out.AttachedTo = append([]string(nil), c.AttachedTo...)
	}
	return out
}
