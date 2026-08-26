package bot

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func land(id, name, owner string) *SafeCard {
	return &SafeCard{ID: id, Name: name, DisplayName: name, OwnerID: owner,
		ControllerID: owner, Zone: ZoneBattlefield, Type: "Basic Land"}
}

func TestParseManaCost(t *testing.T) {
	cases := []struct {
		in      string
		generic int
		pips    []string
	}{
		{"", 0, nil},
		{"{G}", 0, []string{"G"}},
		{"{3}{W}{W}", 3, []string{"W", "W"}},
		{"{2}{U/R}", 2, []string{"UR"}},
		{"{X}{R}", 0, []string{"R"}},
		{"{10}", 10, nil},
		{"{B/P}", 0, []string{"B"}},
	}
	for _, c := range cases {
		got := ParseManaCost(c.in)
		if got.Generic != c.generic || !reflect.DeepEqual(got.Pips, c.pips) {
			t.Errorf("ParseManaCost(%q) = {generic:%d pips:%v}, want {generic:%d pips:%v}",
				c.in, got.Generic, got.Pips, c.generic, c.pips)
		}
	}
}

func TestProducedColors(t *testing.T) {
	cases := []struct {
		name, rules, want string
	}{
		{"Forest", "", "G"},
		{"Island", "", "U"},
		{"Wastes", "", "C"},
		{"Command Tower", "{T}: Add one mana of any color in your commander's color identity.", "WUBRGC"},
		{"Llanowar Elves", "{T}: Add {G}.", "G"},
		{"Birds of Paradise", "{T}: Add one mana of any color.", "WUBRGC"},
		{"Bogus", "Flying", ""},
	}
	for _, c := range cases {
		got := producedColors(&SafeCard{Name: c.name, RulesText: c.rules})
		if got != c.want {
			t.Errorf("producedColors(%q) = %q, want %q", c.name, got, c.want)
		}
	}
}

func TestPayManaSolvesColoredRequirements(t *testing.T) {
	// One Forest, one Island, one Mountain. {1}{G} must take the Forest plus
	// exactly one other; a greedy solver that grabs the first two sources in
	// sorted order (Forest, Island) happens to be right here, so the real test
	// is the next case.
	sources := []*SafeCard{
		land("l1", "Forest", "me"), land("l2", "Island", "me"), land("l3", "Mountain", "me"),
	}
	got, ok := PayMana("{1}{G}", sources)
	if !ok {
		t.Fatal("PayMana({1}{G}) with G/U/R available reported unpayable")
	}
	if len(got) != 2 || !containsName(got, "Forest") {
		t.Fatalf("PayMana({1}{G}) = %v, want a Forest plus one other", names(got))
	}

	// The backtracking case: {G}{U} with only Forest+Island available means the
	// generic-first ordering must not eat the Island.
	got, ok = PayMana("{G}{U}", []*SafeCard{land("a", "Island", "me"), land("b", "Forest", "me")})
	if !ok || len(got) != 2 {
		t.Fatalf("PayMana({G}{U}) = %v, ok=%v; want both lands", names(got), ok)
	}

	// Unpayable: no green source at all.
	if _, ok := PayMana("{G}", []*SafeCard{land("a", "Island", "me")}); ok {
		t.Fatal("PayMana({G}) with only an Island reported payable")
	}
	// Unpayable: not enough sources.
	if _, ok := PayMana("{5}", []*SafeCard{land("a", "Forest", "me")}); ok {
		t.Fatal("PayMana({5}) with one land reported payable")
	}
	// Free spells are always payable and tap nothing.
	if got, ok := PayMana("", nil); !ok || len(got) != 0 {
		t.Fatalf("PayMana(\"\") = %v, ok=%v; want (nil, true)", got, ok)
	}
}

// TestPayManaBacktracks is the case a naive first-fit solver gets wrong: the
// only white source is also the only source that can pay the generic.
func TestPayManaBacktracks(t *testing.T) {
	// Cost {1}{W}{W}. Sources: Plains, Plains, Forest.
	// Both Plains must go to the W pips and the Forest must pay the generic.
	sources := []*SafeCard{
		land("p1", "Plains", "me"), land("p2", "Plains", "me"), land("f1", "Forest", "me"),
	}
	got, ok := PayMana("{1}{W}{W}", sources)
	if !ok {
		t.Fatal("PayMana({1}{W}{W}) reported unpayable")
	}
	if len(got) != 3 {
		t.Fatalf("tapped %d sources, want 3: %v", len(got), names(got))
	}
}

func TestPayManaIsDeterministic(t *testing.T) {
	sources := []*SafeCard{
		land("l1", "Forest", "me"), land("l2", "Forest", "me"),
		land("l3", "Island", "me"), land("l4", "Mountain", "me"),
	}
	first, _ := PayMana("{2}{G}", sources)
	for i := 0; i < 50; i++ {
		got, _ := PayMana("{2}{G}", sources)
		if !reflect.DeepEqual(names(got), names(first)) {
			t.Fatalf("run %d: %v, want %v", i, names(got), names(first))
		}
	}
}

// TestLegalMovesUsesOnlyImplementedVerbs is the anti-pattern 8 guard at the
// vocabulary level: the runtime harness proves the steps land, this proves no
// macro can ever name a verb ParseAndExecuteStringCommand does not implement.
func TestLegalMovesUsesOnlyImplementedVerbs(t *testing.T) {
	// The complete verb list from internal/game/game_engine.go:414 (§0.6).
	implemented := map[string]bool{
		"TAP": true, "UNTAP": true, "MOVE": true, "FLIP": true, "DRAW": true,
		"MODIFY_LIFE": true, "SET_COUNTER": true, "SHUFFLE": true,
		"CREATE_TOKEN": true, "ADD_COUNTER": true, "REMOVE_COUNTER": true,
		"SET_CARD_COUNTER": true, "MILL": true, "SCRY": true, "REVEAL_TOP": true,
		"NEXT_TURN": true, "MULLIGAN": true, "KEEP_HAND": true,
	}
	zones := map[string]bool{
		ZoneLibrary: true, ZoneHand: true, ZoneBattlefield: true,
		ZoneGraveyard: true, ZoneExile: true, ZoneStack: true, ZoneCommand: true,
	}

	check := func(moves []Macro) {
		t.Helper()
		if len(moves) == 0 {
			t.Fatal("no moves to check")
		}
		for _, m := range moves {
			if m.Label == "" {
				t.Errorf("macro with empty label: %v", m.Steps)
			}
			if len(m.Steps) == 0 {
				t.Errorf("macro %q has no steps", m.Label)
			}
			for _, step := range m.Steps {
				verb := strings.SplitN(step, ":", 2)[0]
				if !implemented[verb] {
					t.Errorf("macro %q emits unimplemented verb %q (step %q)", m.Label, verb, step)
				}
				if verb == "MOVE" {
					parts := strings.Split(step, ":")
					if len(parts) != 3 || !zones[parts[2]] {
						t.Errorf("macro %q emits bad MOVE step %q", m.Label, step)
					}
				}
			}
		}
	}

	check(LegalMoves(mulliganView()))
	check(LegalMoves(midGameView()))
}

func TestLegalMovesIsDeterministic(t *testing.T) {
	v := midGameView()
	want := labels(LegalMoves(v))
	for i := 0; i < 100; i++ {
		if got := labels(LegalMoves(midGameView())); !reflect.DeepEqual(got, want) {
			t.Fatalf("run %d differs:\n got %v\nwant %v", i, got, want)
		}
	}
}

func TestLegalMovesMulliganPhaseOffersOnlyKeepOrMulligan(t *testing.T) {
	moves := LegalMoves(mulliganView())
	if len(moves) != 2 {
		t.Fatalf("got %d macros, want 2: %v", len(moves), labels(moves))
	}
	if moves[0].KindOf() != KindKeepHand || moves[1].KindOf() != KindMulligan {
		t.Fatalf("wrong kinds: %v %v", moves[0].KindOf(), moves[1].KindOf())
	}
}

func TestLegalMovesOffersCastWithSolvedTaps(t *testing.T) {
	moves := LegalMoves(midGameView())
	var cast *Macro
	for i := range moves {
		if moves[i].KindOf() == KindCast && strings.Contains(moves[i].Label, "Grizzly Bears") {
			cast = &moves[i]
		}
	}
	if cast == nil {
		t.Fatalf("no Grizzly Bears cast offered: %v", labels(moves))
	}
	// {1}{G}: two TAP steps, then stack, then battlefield.
	if len(cast.Steps) != 4 {
		t.Fatalf("cast steps = %v, want 2 taps + stack + battlefield", cast.Steps)
	}
	if !strings.HasPrefix(cast.Steps[0], "TAP:") || !strings.HasPrefix(cast.Steps[1], "TAP:") {
		t.Fatalf("cast does not begin with solved taps: %v", cast.Steps)
	}
	if !strings.HasSuffix(cast.Steps[2], ":"+ZoneStack) {
		t.Fatalf("cast does not go through the stack: %v", cast.Steps)
	}
	if !strings.HasSuffix(cast.Steps[3], ":"+ZoneBattlefield) {
		t.Fatalf("permanent did not resolve to the battlefield: %v", cast.Steps)
	}
}

func TestLegalMovesInstantResolvesToGraveyard(t *testing.T) {
	moves := LegalMoves(midGameView())
	for _, m := range moves {
		if m.KindOf() == KindCast && strings.Contains(m.Label, "Lightning Bolt") {
			last := m.Steps[len(m.Steps)-1]
			if !strings.HasSuffix(last, ":"+ZoneGraveyard) {
				t.Fatalf("instant resolved to %q, want the graveyard", last)
			}
			return
		}
	}
	t.Fatalf("no Lightning Bolt cast offered: %v", labels(moves))
}

func TestLegalMovesNeverTargetsDeadOpponents(t *testing.T) {
	v := midGameView()
	v.Opponents[0].Life = 0
	dead := v.Opponents[0].PlayerID
	for _, m := range LegalMoves(v) {
		for _, step := range m.Steps {
			if strings.HasPrefix(step, "MODIFY_LIFE:"+dead+":-") {
				t.Fatalf("macro %q attacks an eliminated seat: %q", m.Label, step)
			}
		}
	}
}

func TestEnrichFillsNameOnlyCards(t *testing.T) {
	v := &SafeView{
		ViewerID: "me",
		Me: &SafePlayerView{
			PlayerID: "me", KeptHand: true,
			Hand: []*SafeCard{{ID: "c1", Name: "Grizzly Bears"}},
		},
	}
	Enrich(context.Background(), v, simOracle())
	c := v.Me.Hand[0]
	if c.ManaCost != "{1}{G}" || !strings.Contains(c.Type, "Creature") || c.Power != "2" {
		t.Fatalf("Enrich left card unpopulated: %+v", c)
	}
}

// --- fixtures ---------------------------------------------------------------

func mulliganView() *SafeView {
	return &SafeView{
		GameID: "g", ViewerID: "me", ActivePlayerID: "me", Turn: 1,
		Me: &SafePlayerView{
			PlayerID: "me", Name: "Me", Life: 20, LibraryCount: 93, KeptHand: false,
			Hand: []*SafeCard{{ID: "h1", Name: "Forest", Type: "Basic Land"}},
		},
		Opponents: []*SafeOpponentView{{PlayerID: "opp", Name: "Opp", Life: 20}},
	}
}

// midGameView: the viewer's turn, two Forests untapped, a Mountain untapped,
// a summoning-sick Llanowar Elves, a castable creature and a castable instant.
func midGameView() *SafeView {
	return &SafeView{
		GameID: "g", ViewerID: "me", ActivePlayerID: "me", Turn: 4,
		Me: &SafePlayerView{
			PlayerID: "me", Name: "Me", Life: 20, LibraryCount: 80, KeptHand: true,
			Hand: []*SafeCard{
				{ID: "h1", Name: "Grizzly Bears", ManaCost: "{1}{G}", Type: "Creature — Bear", Power: "2", Toughness: "2"},
				{ID: "h2", Name: "Lightning Bolt", ManaCost: "{R}", Type: "Instant"},
				{ID: "h3", Name: "Forest", Type: "Basic Land — Forest"},
				{ID: "h4", Name: "Serra Angel", ManaCost: "{3}{W}{W}", Type: "Creature — Angel", Power: "4", Toughness: "4"},
			},
		},
		Battlefield: []*SafeCard{
			{ID: "b1", Name: "Forest", ControllerID: "me", Type: "Basic Land — Forest"},
			{ID: "b2", Name: "Forest", ControllerID: "me", Type: "Basic Land — Forest"},
			{ID: "b3", Name: "Mountain", ControllerID: "me", Type: "Basic Land — Mountain"},
			{ID: "b4", Name: "Llanowar Elves", ControllerID: "me", Type: "Creature — Elf Druid",
				Power: "1", Toughness: "1", RulesText: "{T}: Add {G}.", SummoningSickness: true},
			{ID: "b5", Name: "Hill Giant", ControllerID: "me", Type: "Creature — Giant", Power: "3", Toughness: "3"},
			{ID: "b6", Name: "Air Elemental", ControllerID: "opp", Type: "Creature — Elemental", Power: "4", Toughness: "4"},
		},
		Opponents: []*SafeOpponentView{
			{PlayerID: "opp", Name: "Opp", Life: 20, LibraryCount: 80},
			{PlayerID: "opp2", Name: "Opp2", Life: 17, LibraryCount: 80},
		},
	}
}

func labels(ms []Macro) []string {
	out := make([]string, 0, len(ms))
	for _, m := range ms {
		out = append(out, m.Label+" => "+strings.Join(m.Steps, " ; "))
	}
	return out
}

func names(cs []*SafeCard) []string {
	out := make([]string, 0, len(cs))
	for _, c := range cs {
		out = append(out, c.Name+"/"+c.ID)
	}
	return out
}

func containsName(cs []*SafeCard, name string) bool {
	for _, c := range cs {
		if c.Name == name {
			return true
		}
	}
	return false
}

// TestRandomPolicyIsReproducible pins the injectable-seed requirement: a bot
// simulation that cannot be replayed from a seed is not a debugging tool.
func TestRandomPolicyIsReproducible(t *testing.T) {
	moves := LegalMoves(midGameView())
	pick := func(seed int64) []string {
		p := NewRandomPolicy(seed)
		out := make([]string, 0, 50)
		for i := 0; i < 50; i++ {
			m, err := p.Pick(context.Background(), nil, moves)
			if err != nil {
				t.Fatalf("Pick: %v", err)
			}
			out = append(out, m.Label)
		}
		return out
	}
	if !reflect.DeepEqual(pick(42), pick(42)) {
		t.Fatal("same seed produced different choices")
	}
	if reflect.DeepEqual(pick(42), pick(43)) {
		t.Fatal("different seeds produced identical choices")
	}
	if _, err := NewRandomPolicy(1).Pick(context.Background(), nil, nil); err != ErrNoMoves {
		t.Fatalf("empty move set: got %v, want ErrNoMoves", err)
	}
}
