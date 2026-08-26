package bot

import (
	"github.com/magefree/mage-server-go/internal/game"
)

// fixtures_test.go builds *game.PlaytestGameView values shaped exactly the way
// GameEngine.buildGameView (internal/game/view.go:68) shapes them, so that the
// serializer and leak tests exercise the same input the engine produces --
// including the pointer sharing that Redact exists to defeat.

const (
	seatAlice = "alice"
	seatBob   = "bob"
)

// secretLibraryNames are the cards planted in the VIEWER'S OWN library. They
// appear nowhere else in any fixture, so any one of them surfacing in a
// SafeView or in rendered text is unambiguous proof that the library strip in
// Redact stopped working.
var secretLibraryNames = []string{
	"Ancestral Recall",
	"Black Lotus",
	"Time Walk",
}

// opponentSecretNames are planted in the opponent's Player-side hand/library in
// the engine state. buildGameView must never copy them into the view; the
// fixtures reproduce that by leaving the opponent view's Hand/Library empty.
var opponentSecretNames = []string{
	"Timetwister",
	"Mox Sapphire",
}

func card(id, name, owner string) *game.Card {
	return &game.Card{
		ID:           id,
		Name:         name,
		DisplayName:  name,
		OwnerID:      owner,
		ControllerID: owner,
		Counters:     []game.Counter{},
		AttachedTo:   []string{},
	}
}

func inZone(c *game.Card, zone string) *game.Card {
	c.Zone = zone
	return c
}

// testOracle is the stub OracleLookup every test uses. It stands in for
// internal/repository.ScryfallCardRepository so no test needs a database.
func testOracle() MapOracle {
	return MapOracle{
		"Llanowar Elves": {
			Name: "Llanowar Elves", ManaCost: "{G}",
			TypeLine:   "Creature — Elf Druid",
			OracleText: "{T}: Add {G}.",
			Power:      "1", Toughness: "1",
		},
		"Grizzly Bears": {
			Name: "Grizzly Bears", ManaCost: "{1}{G}",
			TypeLine: "Creature — Bear",
			Power:    "2", Toughness: "2",
		},
		"Serra Angel": {
			Name: "Serra Angel", ManaCost: "{3}{W}{W}",
			TypeLine:   "Creature — Angel",
			OracleText: "Flying, vigilance",
			Power:      "4", Toughness: "4",
		},
		"Lightning Bolt": {
			Name: "Lightning Bolt", ManaCost: "{R}",
			TypeLine:   "Instant",
			OracleText: "Lightning Bolt deals 3 damage to any target.",
		},
		"Counterspell": {
			Name: "Counterspell", ManaCost: "{U}{U}",
			TypeLine:   "Instant",
			OracleText: "Counter target spell.",
		},
		"Wall of Omens": {
			Name: "Wall of Omens", ManaCost: "{1}{W}",
			TypeLine:   "Creature — Wall",
			OracleText: "Defender\nWhen Wall of Omens enters, draw a card.",
			Power:      "0", Toughness: "4",
		},
		"Chandra, Torch of Defiance": {
			Name: "Chandra, Torch of Defiance", ManaCost: "{2}{R}{R}",
			TypeLine:   "Legendary Planeswalker — Chandra",
			OracleText: "+1: Exile the top card of your library.\n−3: Chandra deals 4 damage to target creature.",
			Loyalty:    "4",
		},
		"Sol Ring": {
			Name: "Sol Ring", ManaCost: "{1}",
			TypeLine:   "Artifact",
			OracleText: "{T}: Add {C}{C}.",
		},
		// Basic lands are in the oracle so the tests prove they are excluded by
		// the Card Reference filter, not merely missing from the lookup.
		"Forest":   {Name: "Forest", TypeLine: "Basic Land — Forest", OracleText: "({T}: Add {G}.)"},
		"Mountain": {Name: "Mountain", TypeLine: "Basic Land — Mountain", OracleText: "({T}: Add {R}.)"},
		"Plains":   {Name: "Plains", TypeLine: "Basic Land — Plains", OracleText: "({T}: Add {W}.)"},
	}
}

// baseView returns a mid-game view from Alice's seat, with the same pointer
// sharing buildGameView produces.
func baseView() *game.PlaytestGameView {
	aliceLibrary := make([]*game.Card, 0, len(secretLibraryNames))
	for i, n := range secretLibraryNames {
		aliceLibrary = append(aliceLibrary, inZone(card("lib-alice-"+itoa(i), n, seatAlice), game.ZoneLibraryStr))
	}

	aliceHand := []*game.Card{
		inZone(card("h-a-0", "Lightning Bolt", seatAlice), game.ZoneHandStr),
		inZone(card("h-a-1", "Forest", seatAlice), game.ZoneHandStr),
		inZone(card("h-a-2", "Serra Angel", seatAlice), game.ZoneHandStr),
	}

	battlefield := []*game.Card{
		inZone(card("bf-a-0", "Llanowar Elves", seatAlice), game.ZoneBattlefieldStr),
		inZone(card("bf-a-1", "Forest", seatAlice), game.ZoneBattlefieldStr),
		inZone(card("bf-b-0", "Grizzly Bears", seatBob), game.ZoneBattlefieldStr),
	}
	battlefield[0].SummoningSickness = true
	battlefield[1].Tapped = true
	battlefield[2].Counters = []game.Counter{{Name: "+1/+1", Count: 2}}

	view := &game.PlaytestGameView{
		GameID:         "game-1",
		ViewerID:       seatAlice,
		Battlefield:    battlefield,
		Exile:          []*game.Card{inZone(card("ex-b-0", "Sol Ring", seatBob), game.ZoneExileStr)},
		Stack:          []*game.Card{},
		Command:        []*game.Card{},
		Turn:           3,
		ActivePlayerID: seatAlice,
		IsInitialized:  true,
		Log:            []game.LogEntry{{Kind: "play", Message: "Alice plays Forest."}},
		Me: &game.PlaytestPlayerView{
			PlayerID:     seatAlice,
			Name:         "Alice",
			Life:         20,
			LibraryCount: len(aliceLibrary),
			HandCount:    len(aliceHand),
			Hand:         aliceHand,
			Library:      aliceLibrary, // buildGameView shares this by pointer
			Graveyard:    []*game.Card{},
			ManaPool:     &game.ManaPool{},
		},
		Opponents: []*game.PlaytestOpponentView{
			{
				PlayerID:     seatBob,
				Name:         "Bob",
				Life:         18,
				Poison:       1,
				LibraryCount: 40,
				HandCount:    5,
				Hand:         []*game.Card{}, // hidden, per view.go:117
				Library:      []*game.Card{}, // hidden, per view.go:118
				Graveyard:    []*game.Card{inZone(card("gy-b-0", "Counterspell", seatBob), game.ZoneGraveyardStr)},
				ManaPool:     &game.ManaPool{},
			},
		},
	}
	return view
}

func itoa(i int) string {
	const digits = "0123456789"
	if i == 0 {
		return "0"
	}
	var b []byte
	for i > 0 {
		b = append([]byte{digits[i%10]}, b...)
		i /= 10
	}
	return string(b)
}

func intp(i int) *int { return &i }
