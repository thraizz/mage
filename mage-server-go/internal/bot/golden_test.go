package bot

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/magefree/mage-server-go/internal/game"
)

// golden_test.go pins the rendered board text. These are the first golden files
// in the repo; the format follows reference/json5_utils.py so that a change to
// the renderer shows up as a line-by-line diff of the board, not as one
// modified 4KB escaped string.
//
// Regenerate with:  UPDATE_GOLDEN=1 go test ./internal/bot/...
//
// Regenerating is not the same as approving. Read the diff. A golden that
// changed for a reason you cannot state is a bug you just enshrined.

const goldenDir = "testdata/golden"

type goldenCase struct {
	name  string
	build func() (*game.PlaytestGameView, *Decision)
	// after mutates the redacted view. It exists for the handful of SafeCard
	// fields the engine has no equivalent for (OriginalCard / IsCopy), which
	// therefore cannot be set on the input *game.Card.
	after func(*SafeView)
}

func goldenCases() []goldenCase {
	return []goldenCase{
		{
			name: "opening_mulligan",
			build: func() (*game.PlaytestGameView, *Decision) {
				v := baseView()
				// Pregame: nothing on the battlefield, no turn yet.
				v.Turn = 0
				v.Battlefield = nil
				v.Exile = nil
				v.IsInitialized = false
				v.Me.ManaPool = &game.ManaPool{}
				v.Opponents[0].HandCount = 7
				v.Opponents[0].Graveyard = nil
				v.Opponents[0].Poison = 0
				return v, &Decision{
					Index:         0,
					SnapshotIndex: 0,
					Turn:          0,
					Phase:         "", // renders as PREGAME
					Player:        "Alice",
					Message:       "Mulligan?",
					Choices: []Choice{
						{Name: "Mulligan", ID: "yes"},
						{Name: "Keep", ID: "no"},
					},
					RespondWith: "choice=yes to mulligan, or choice=no to keep",
				}
			},
		},
		{
			name: "land_drop_main_phase",
			build: func() (*game.PlaytestGameView, *Decision) {
				v := baseView()
				return v, &Decision{
					Index:         12,
					SnapshotIndex: 7,
					Turn:          3,
					Phase:         "PRECOMBAT_MAIN",
					Player:        "Alice",
					Message:       "Play a land or cast a spell",
					Choices: []Choice{
						{Name: "Forest", ID: "p2", Action: "play land"},
						{Name: "Lightning Bolt", ID: "p4", Action: "cast", ManaCost: "{R}"},
						{Name: "Serra Angel", ID: "p6", Action: "cast", ManaCost: "{3}{W}{W}"},
					},
					PilotContext: &PilotContext{
						UntappedLands: intp(0),
						LandDropsUsed: intp(0),
					},
					RespondWith: "choice=pN to play, or choice=no to pass",
				}
			},
		},
		{
			name: "stack_non_empty",
			build: func() (*game.PlaytestGameView, *Decision) {
				v := baseView()
				bolt := v.Me.Hand[0]
				v.Me.Hand = v.Me.Hand[1:]
				v.Me.HandCount = len(v.Me.Hand)
				bolt.Zone = game.ZoneStackStr
				v.Stack = []*game.Card{bolt}
				v.Me.ManaPool = &game.ManaPool{Red: 1, Green: 2}
				return v, &Decision{
					Index:         21,
					SnapshotIndex: 11,
					Turn:          3,
					Phase:         "PRECOMBAT_MAIN",
					Player:        "Alice",
					Message:       "Respond to Lightning Bolt?",
					Choices: []Choice{
						{Name: "Pass priority", ID: "no"},
					},
					RespondWith: "choice=pN to respond, or choice=no to pass",
					RecentChat:  []string{"Bob: nice bolt", "Alice: thanks"},
				}
			},
		},
		{
			name: "declare_attackers",
			build: func() (*game.PlaytestGameView, *Decision) {
				v := baseView()
				v.Battlefield[0].SummoningSickness = false
				v.Battlefield[0].Attacking = true
				return v, &Decision{
					Index:         30,
					SnapshotIndex: 16,
					Turn:          5,
					Phase:         "COMBAT",
					Player:        "Alice",
					Message:       "Declare attackers",
					Choices: []Choice{
						{Name: "Llanowar Elves", ID: "p3", Action: "attack"},
					},
					Combat: []CombatGroup{
						{Attackers: []string{"Llanowar Elves"}, Defending: "Bob"},
					},
					PilotContext: &PilotContext{CombatPhase: "declare_attackers"},
					RespondWith:  "attackers=p3,p4 or attackers=all, or choice=no to attack with nothing",
				}
			},
		},
		{
			name: "declare_blockers_incoming",
			build: func() (*game.PlaytestGameView, *Decision) {
				v := baseView()
				// Bob is attacking Alice; Alice is the one blocking.
				v.Battlefield[2].Attacking = true
				v.Battlefield = append(v.Battlefield,
					inZone(card("bf-b-1", "Serra Angel", seatBob), game.ZoneBattlefieldStr))
				v.Battlefield[3].Attacking = true
				v.Battlefield[0].SummoningSickness = false
				return v, &Decision{
					Index:         41,
					SnapshotIndex: 22,
					Turn:          6,
					Phase:         "COMBAT",
					Player:        "Alice",
					Message:       "Declare blockers",
					Choices: []Choice{
						{Name: "Llanowar Elves", ID: "p3", Action: "block"},
					},
					Combat: []CombatGroup{
						{Attackers: []string{"Grizzly Bears", "Serra Angel"}, Defending: "Alice"},
					},
					PilotContext: &PilotContext{
						CombatPhase: "declare_blockers",
						IncomingAttackers: []IncomingAttacker{
							{Name: "Grizzly Bears", ID: "p5", PowerToughness: "4/4"},
							{Name: "Serra Angel", ID: "p7", PowerToughness: "4/4"},
						},
					},
					RespondWith: "blockers=blocker:attacker pairs, or choice=no to block with nothing",
				}
			},
		},
		{
			name: "multi_amount_items",
			build: func() (*game.PlaytestGameView, *Decision) {
				v := baseView()
				return v, &Decision{
					Index:         55,
					SnapshotIndex: 30,
					Turn:          7,
					Phase:         "POSTCOMBAT_MAIN",
					Player:        "Alice",
					Message:       "Distribute 4 damage among target creatures",
					Items: []MultiAmountItem{
						{Description: "Grizzly Bears", Min: intp(0), Max: intp(4)},
						{Description: "Serra Angel", Min: intp(0), Max: intp(4)},
						{Description: "Wall of Omens"},
					},
					TotalMin:    intp(4),
					TotalMax:    intp(4),
					RespondWith: "amounts=[N,N,...] — one per item, sum between total_min and total_max",
				}
			},
		},
		{
			name: "multi_amount_items_range",
			build: func() (*game.PlaytestGameView, *Decision) {
				v := baseView()
				return v, &Decision{
					Index:         56,
					SnapshotIndex: 31,
					Turn:          7,
					Phase:         "POSTCOMBAT_MAIN",
					Player:        "Alice",
					Message:       "Pick triggered ability order",
					Items: []MultiAmountItem{
						{Description: "Llanowar Elves trigger", Min: intp(0)},
						{Description: "Wall of Omens trigger", Max: intp(3)},
					},
					TotalMin:     intp(1),
					TotalMax:     intp(3),
					ResponseType: "multi_amount",
				}
			},
		},
		{
			name: "planeswalker_and_token_permanents",
			build: func() (*game.PlaytestGameView, *Decision) {
				v := baseView()
				pw := inZone(card("bf-a-2", "Chandra, Torch of Defiance", seatAlice), game.ZoneBattlefieldStr)
				pw.Loyalty = "4"
				tok := inZone(card("token-1700000000-42", "Beast", seatAlice), game.ZoneBattlefieldStr)
				tok.Power, tok.Toughness = "3", "3"
				fd := inZone(card("bf-a-3", "Wall of Omens", seatAlice), game.ZoneBattlefieldStr)
				fd.FaceDown = true
				fd.Tapped = true
				v.Battlefield = append(v.Battlefield, pw, tok, fd)
				return v, &Decision{
					Index:         60,
					SnapshotIndex: 33,
					Turn:          8,
					Phase:         "PRECOMBAT_MAIN",
					Player:        "Alice",
					Message:       "",
					Choices:       []Choice{},
					ResponseType:  "priority",
				}
			},
			after: func(sv *SafeView) {
				for _, c := range sv.Battlefield {
					if c.ID == "bf-a-3" {
						c.OriginalCard = "Grizzly Bears"
					}
				}
			},
		},
	}
}

func TestGoldenRenders(t *testing.T) {
	for _, tc := range goldenCases() {
		t.Run(tc.name, func(t *testing.T) {
			view, decision := tc.build()
			safe, err := RedactErr(view, seatAlice)
			require.NoError(t, err)
			if tc.after != nil {
				tc.after(safe)
			}

			s := NewSerializer(testOracle())
			out := s.Render(context.Background(), safe, decision)

			assertGolden(t, tc.name, map[string]any{
				"scenario": tc.name,
				"prompt":   out,
			})
		})
	}
}

// TestGoldenOracleDedup pins the first-appearance-only behaviour: the same
// board rendered twice must carry a Card Reference the first time and none of
// the same entries the second.
func TestGoldenOracleDedup(t *testing.T) {
	view, decision := goldenCases()[1].build()
	safe, err := RedactErr(view, seatAlice)
	require.NoError(t, err)

	s := NewSerializer(testOracle())
	ctx := context.Background()
	first := s.Render(ctx, safe, decision)

	// Second decision, same board plus one new card in hand. Only the new card
	// may appear in the Card Reference.
	view2, decision2 := goldenCases()[1].build()
	view2.Me.Hand = append(view2.Me.Hand,
		inZone(card("h-a-3", "Wall of Omens", seatAlice), game.ZoneHandStr))
	view2.Me.HandCount = len(view2.Me.Hand)
	decision2.Index++
	safe2, err := RedactErr(view2, seatAlice)
	require.NoError(t, err)
	second := s.Render(ctx, safe2, decision2)

	require.Contains(t, first, "- Lightning Bolt {R} -- Instant:")
	require.NotContains(t, second, "- Lightning Bolt",
		"oracle text must be emitted on a card's first appearance only")
	require.Contains(t, second, "- Wall of Omens {1}{W} -- Creature — Wall 0/4:")

	assertGolden(t, "oracle_dedup", map[string]any{
		"scenario":     "oracle_dedup",
		"first_render": first,
		"render_2":     second,
	})
}

// assertGolden compares payload against testdata/golden/<name>.json5, or
// rewrites it when UPDATE_GOLDEN=1.
func assertGolden(t *testing.T, name string, payload map[string]any) {
	t.Helper()

	got, err := dumpsJSON5(payload)
	require.NoError(t, err)
	got += "\n"

	path := filepath.Join(goldenDir, name+".json5")
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		require.NoError(t, os.MkdirAll(goldenDir, 0o755))
		require.NoError(t, os.WriteFile(path, []byte(got), 0o644))
		t.Logf("updated golden %s", path)
		return
	}

	want, err := os.ReadFile(path)
	require.NoError(t, err, "missing golden file %s; run UPDATE_GOLDEN=1 go test ./internal/bot/...", path)
	require.Equal(t, string(want), got,
		"golden %s is stale.\nIf the change is intended, re-read the diff line by line, then run:\n  UPDATE_GOLDEN=1 go test ./internal/bot/...", path)
}
