package bot

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func render(t *testing.T, s *Serializer, v *SafeView, d *Decision) string {
	t.Helper()
	return s.Render(context.Background(), v, d)
}

func TestEmptyPhaseRendersAsPregame(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	out := render(t, NewSerializer(nil), safe, &Decision{
		Index: 0, SnapshotIndex: 0, Turn: 0, Phase: "", Player: "Alice",
	})
	require.Contains(t, out, "[Decision 0, snapshot=0] Turn 0 PREGAME - Alice")
}

func TestMessageLineAlwaysEmittedEvenWhenEmpty(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	out := render(t, NewSerializer(nil), safe, &Decision{Player: "Alice"})
	require.Contains(t, out, "\n  Message: \n")
}

func TestChoicesBranchVersusItemsBranch(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	s := NewSerializer(nil)

	choices := render(t, s, safe, &Decision{
		Player:  "Alice",
		Choices: []Choice{{Name: "Forest", ID: "p2", Action: "play land"}},
	})
	require.Contains(t, choices, "  Choices (1): Forest [id=p2, play land]")
	require.NotContains(t, choices, "Items (")

	items := render(t, s, safe, &Decision{
		Player:   "Alice",
		Choices:  []Choice{{Name: "ignored"}},
		Items:    []MultiAmountItem{{Description: "Bears", Min: intp(0), Max: intp(4)}},
		TotalMin: intp(4), TotalMax: intp(4),
	})
	require.Contains(t, items, "  Items (1): total=4")
	require.Contains(t, items, "    0: Bears [min=0, max=4]")
	require.NotContains(t, items, "Choices (")
}

func TestItemsHeaderWithDistinctTotals(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	out := render(t, NewSerializer(nil), safe, &Decision{
		Player:   "Alice",
		Items:    []MultiAmountItem{{Description: "A"}},
		TotalMin: intp(1), TotalMax: intp(3),
	})
	require.Contains(t, out, "  Items (1): total_min=1, total_max=3")
}

func TestLandDropsAndUntappedLandsDistinguishZeroFromAbsent(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	s := NewSerializer(nil)

	present := render(t, s, safe, &Decision{
		Player:       "Alice",
		PilotContext: &PilotContext{UntappedLands: intp(0), LandDropsUsed: intp(0)},
	})
	require.Contains(t, present, "  Untapped lands: 0, Land drops remaining: 1")

	absent := render(t, s, safe, &Decision{Player: "Alice", PilotContext: &PilotContext{}})
	require.NotContains(t, absent, "Untapped lands")
}

func TestIncomingAttackersOnlyAtDeclareBlockers(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	s := NewSerializer(nil)
	attackers := []IncomingAttacker{{Name: "Serra Angel", ID: "p9", PowerToughness: "4/4"}}

	blocking := render(t, s, safe, &Decision{
		Player:       "Alice",
		PilotContext: &PilotContext{CombatPhase: "declare_blockers", IncomingAttackers: attackers},
	})
	require.Contains(t, blocking, "  Incoming Attackers: Serra Angel [id=p9, 4/4]")

	attacking := render(t, s, safe, &Decision{
		Player:       "Alice",
		PilotContext: &PilotContext{CombatPhase: "declare_attackers", IncomingAttackers: attackers},
	})
	require.NotContains(t, attacking, "Incoming Attackers")
}

func TestTriggeredAbilityNote(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	out := render(t, NewSerializer(nil), safe, &Decision{
		Player: "Alice", Message: "Pick triggered ability to put on the stack",
	})
	require.Contains(t, out, "  NOTE: This decision only determines the order triggered abilities")
}

func TestBoardLineShapes(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	out := render(t, NewSerializer(testOracle()), safe, &Decision{Player: "Alice"})
	board := boardLine(t, out)

	// Own seat: full hand. Opponent: count only, plus a player counter.
	require.Contains(t, board, "Alice: 20hp hand=[Lightning Bolt, Forest, Serra Angel] lib=3")
	require.Contains(t, board, "Bob: 18hp hand=5 lib=40 poison=1")
	// Seats are sorted by name so the line is stable across renders.
	require.True(t, strings.Index(board, "Alice:") < strings.Index(board, "Bob:"))
	// Battlefield is partitioned by controller, exile by owner.
	require.Contains(t, board, "bf=[Llanowar Elves 1/1 (sick), Forest (tapped)]")
	require.Contains(t, board, "bf=[Grizzly Bears 2/2 (+1/+1=2)]")
	require.Contains(t, board, "gy=[Counterspell]")
	require.Contains(t, board, "exile=[Sol Ring]")
}

func TestOwnEmptyHandRendersAsZero(t *testing.T) {
	v := baseView()
	v.Me.Hand = nil
	v.Me.HandCount = 0
	safe := Redact(v, seatAlice)
	board := boardLine(t, render(t, NewSerializer(nil), safe, &Decision{Player: "Alice"}))
	require.Contains(t, board, "Alice: 20hp hand=0 lib=3")
}

func TestOpponentZeroHandCountOmitsHand(t *testing.T) {
	v := baseView()
	v.Opponents[0].HandCount = 0
	safe := Redact(v, seatAlice)
	board := boardLine(t, render(t, NewSerializer(nil), safe, &Decision{Player: "Alice"}))
	require.Contains(t, board, "Bob: 18hp lib=40")
	require.NotContains(t, board, "hand=0 lib=40")
}

func TestManaPoolLineOnlyWhenNonZero(t *testing.T) {
	v := baseView()
	safe := Redact(v, seatAlice)
	require.NotContains(t, render(t, NewSerializer(nil), safe, &Decision{Player: "Alice"}), "Mana pool:")

	v2 := baseView()
	v2.Me.ManaPool.Red = 2
	v2.Me.ManaPool.Colorless = 1
	safe2 := Redact(v2, seatAlice)
	require.Contains(t, render(t, NewSerializer(nil), safe2, &Decision{Player: "Alice"}),
		"  Mana pool: R=2, C=1")
}

func TestRespondWithTotalSubstitution(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	out := render(t, NewSerializer(nil), safe, &Decision{
		Player:      "Alice",
		Items:       []MultiAmountItem{{Description: "A"}},
		TotalMin:    intp(4),
		TotalMax:    intp(4),
		RespondWith: "amounts=[N,N,...] — one per item, sum between total_min and total_max",
	})
	require.Contains(t, out, "sum must equal total (4)")
	require.NotContains(t, out, "sum between total_min and total_max")
}

func TestResponseTypeUsedOnlyWhenRespondWithAbsent(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	s := NewSerializer(nil)
	both := render(t, s, safe, &Decision{Player: "Alice", RespondWith: "do a thing", ResponseType: "priority"})
	require.Contains(t, both, "  Respond: do a thing")
	require.NotContains(t, both, "Response type:")

	only := render(t, s, safe, &Decision{Player: "Alice", ResponseType: "priority"})
	require.Contains(t, only, "  Response type: priority")
}

func TestRecentChat(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	out := render(t, NewSerializer(nil), safe, &Decision{
		Player: "Alice", RecentChat: []string{"Bob: gg", "Alice: gg"},
	})
	require.Contains(t, out, "  Recent chat: Bob: gg | Alice: gg")
}

func TestCardReferenceExcludesBasicLandsAndCondensesOracleText(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	out := render(t, NewSerializer(testOracle()), safe, &Decision{Player: "Alice"})

	ref := out[:strings.Index(out, "## Decision")]
	require.Contains(t, ref, "- Lightning Bolt {R} -- Instant: Lightning Bolt deals 3 damage to any target.")
	require.Contains(t, ref, "- Serra Angel {3}{W}{W} -- Creature — Angel 4/4: Flying, vigilance")
	require.NotContains(t, ref, "- Forest", "basic lands are excluded from the Card Reference")

	// Multi-line oracle text is condensed with " / ".
	v := baseView()
	v.Me.Hand = append(v.Me.Hand, inZone(card("h-a-9", "Wall of Omens", seatAlice), "HAND"))
	out2 := render(t, NewSerializer(testOracle()), Redact(v, seatAlice), &Decision{Player: "Alice"})
	require.Contains(t, out2, "- Wall of Omens {1}{W} -- Creature — Wall 0/4: Defender / When Wall of Omens enters, draw a card.")
	require.NotContains(t, out2, "Defender\nWhen")
}

func TestCardReferenceIncludesChoiceNames(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	out := render(t, NewSerializer(testOracle()), safe, &Decision{
		Player:  "Alice",
		Choices: []Choice{{Name: "Chandra, Torch of Defiance", ID: "p8", Action: "cast"}},
	})
	require.Contains(t, out, "- Chandra, Torch of Defiance {2}{R}{R} -- Legendary Planeswalker — Chandra:")
}

func TestNoCardReferenceWithoutOracle(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	out := render(t, NewSerializer(nil), safe, &Decision{Player: "Alice"})
	require.False(t, strings.HasPrefix(out, "## Card Reference"))
	require.True(t, strings.HasPrefix(out, "## Decision"))
}

func TestResetOracleDedup(t *testing.T) {
	safe := Redact(baseView(), seatAlice)
	s := NewSerializer(testOracle())
	require.Contains(t, render(t, s, safe, &Decision{Player: "Alice"}), "- Lightning Bolt")
	require.NotContains(t, render(t, s, safe, &Decision{Player: "Alice"}), "- Lightning Bolt")
	s.ResetOracleDedup()
	require.Contains(t, render(t, s, safe, &Decision{Player: "Alice"}), "- Lightning Bolt")
}

func TestPermanentDisplayExtras(t *testing.T) {
	r := &renderState{ctx: context.Background(), s: NewSerializer(testOracle())}

	require.Equal(t, "Grizzly Bears 2/2", r.permanentDisplay(&SafeCard{Name: "Grizzly Bears"}))
	require.Equal(t, "Grizzly Bears 2/2 (tapped, sick, face_down)", r.permanentDisplay(&SafeCard{
		Name: "Grizzly Bears", Tapped: true, SummoningSickness: true, FaceDown: true,
	}))
	require.Equal(t, "Chandra, Torch of Defiance (loyalty=4)", r.permanentDisplay(&SafeCard{
		Name: "Chandra, Torch of Defiance", Loyalty: "4",
	}))
	require.Equal(t, "Grizzly Bears 2/2 (+1/+1=3, stun=1)", r.permanentDisplay(&SafeCard{
		Name: "Grizzly Bears", Counters: []SafeCounter{{Name: "+1/+1", Count: 3}, {Name: "stun", Count: 1}},
	}))
	require.Equal(t, "Bear 3/3 (token)", r.permanentDisplay(&SafeCard{
		ID: "token-123-4", Name: "Bear", Power: "3", Toughness: "3",
	}))
	require.Equal(t, "Grizzly Bears 2/2 (copy of Serra Angel)", r.permanentDisplay(&SafeCard{
		Name: "Grizzly Bears", OriginalCard: "Serra Angel",
	}))
	require.Equal(t, "Grizzly Bears 2/2 (copy)", r.permanentDisplay(&SafeCard{
		Name: "Grizzly Bears", IsCopy: true,
	}))
	// Card-supplied P/T wins over the oracle fallback.
	require.Equal(t, "Grizzly Bears 5/5", r.permanentDisplay(&SafeCard{
		Name: "Grizzly Bears", Power: "5", Toughness: "5",
	}))
}

func TestFormatChoice(t *testing.T) {
	require.Equal(t, "Forest", formatChoice(Choice{Name: "Forest"}))
	require.Equal(t, "Forest [id=p2]", formatChoice(Choice{Name: "Forest", ID: "p2"}))
	require.Equal(t, "Bolt [id=p4, cast, {R}]", formatChoice(Choice{Name: "Bolt", ID: "p4", Action: "cast", ManaCost: "{R}"}))
	require.Equal(t, "a description", formatChoice(Choice{Description: "a description"}))
	require.Equal(t, "?", formatChoice(Choice{}))
}

func TestRenderCombat(t *testing.T) {
	require.Equal(t, "Bears blocked by Elves -> Alice", renderCombat([]CombatGroup{
		{Attackers: []string{"Bears"}, Blockers: []string{"Elves"}, Defending: "Alice"},
	}))
	require.Equal(t, "Bears (blocked)", renderCombat([]CombatGroup{
		{Attackers: []string{"Bears"}, Blocked: true},
	}))
	require.Equal(t, "Bears, Angel -> Bob | Elves -> Bob", renderCombat([]CombatGroup{
		{Attackers: []string{"Bears", "Angel"}, Defending: "Bob"},
		{Attackers: []string{"Elves"}, Defending: "Bob"},
	}))
}

func TestStackLineOnlyWhenNonEmpty(t *testing.T) {
	v := baseView()
	safe := Redact(v, seatAlice)
	require.NotContains(t, render(t, NewSerializer(nil), safe, &Decision{Player: "Alice"}), "  Stack:")

	v2 := baseView()
	bolt := v2.Me.Hand[0]
	bolt.Zone = "STACK"
	v2.Stack = append(v2.Stack, bolt)
	out := render(t, NewSerializer(nil), Redact(v2, seatAlice), &Decision{Player: "Alice"})
	require.Contains(t, out, "  Stack: [Lightning Bolt]")
}

func boardLine(t *testing.T, out string) string {
	t.Helper()
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "  Board: ") {
			return line
		}
	}
	t.Fatalf("no Board line in:\n%s", out)
	return ""
}
