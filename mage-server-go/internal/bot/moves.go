package bot

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// moves.go turns a SafeView into the set of macros a seat can perform.
//
// HONOR SYSTEM, NOT LEGALITY. Read Phase 0 §0.5: internal/game has no rules
// engine. ProcessAction is a dispatcher over twenty direct state mutations --
// no priority, no stack resolution, no combat, no mana payment, no turn
// structure beyond a counter -- and every package that could supply legality
// (mana, targeting, abilities, cards) is orphaned with zero external importers.
//
// So LegalMoves does NOT compute what the rules permit. It computes the manual
// action set a human playtester sitting in front of this client has: pick a
// card up, put it somewhere, tap something, change a life total, end the turn.
// A bot cheats exactly as easily as a human playtester does, and that is the
// deliberate state of the world until Phase 7 wires a real legality layer.
//
// The one thing that is NOT a choice here is mana payment. Which lands to tap
// for a spell is a solver problem with a right answer, so payMana solves it and
// bakes the taps into the macro's steps. A policy -- random today, an LLM in
// Phase 5 -- never gets asked "which Forest?".

// Macro is one atomic decision offered to a Policy.
//
// Label is the only thing a policy sees: it must describe the whole macro in
// one human-readable line, because in Phase 5 it is what goes in the prompt.
//
// Steps are ParseAndExecuteStringCommand strings, executed in order, each sent
// as its own SendPlayerAction(..., "SEND_STRING", step). Every verb used here
// is one the engine actually implements (§0.6); anti-pattern 8 is the reason
// that list is checked rather than assumed -- unknown action types fall into
// ProcessAction's default branch and are silently swallowed.
type Macro struct {
	Label string
	Steps []string

	// kind classifies the macro for the runner's turn-structure rules (one land
	// drop per turn; "pass turn" ends the seat's action loop). It is unexported
	// so that Macro's public shape stays exactly the two fields the design
	// specifies, and so nothing outside this package is tempted to branch on it
	// instead of on Label.
	kind Kind
}

// Kind classifies a macro so the runner can apply turn-structure rules (one
// land drop per turn, "pass turn" ends the seat's action loop) without parsing
// Label strings.
type Kind string

const (
	KindPlayLand   Kind = "play_land"
	KindCast       Kind = "cast"
	KindTap        Kind = "tap"
	KindUntap      Kind = "untap"
	KindMoveZone   Kind = "move_zone"
	KindAttack     Kind = "attack"
	KindModifyLife Kind = "modify_life"
	KindDraw       Kind = "draw"
	KindMulligan   Kind = "mulligan"
	KindKeepHand   Kind = "keep_hand"
	KindPassTurn   Kind = "pass_turn"
)

// KindOf reports the macro's class.
func (m Macro) KindOf() Kind { return m.kind }

func macro(k Kind, label string, steps ...string) Macro {
	return Macro{Label: label, Steps: steps, kind: k}
}

// LegalMoves returns every macro the viewing seat can perform right now.
//
// It is a pure function of the SafeView: no engine access, no oracle lookup, no
// randomness, and the output order is deterministic for a given view. That
// matters because a seeded RandomPolicy indexes into this slice -- a
// nondeterministic order would make bot simulations irreproducible even with a
// fixed seed.
//
// Printed characteristics (type line, mana cost, P/T, "{T}: Add ...") come off
// the SafeCard fields. The engine leaves those blank -- StartGameWithDecks
// creates name-only cards (internal/game/game_engine.go:94-114) -- so callers
// run Enrich first. Without it, only basic lands are recognised, by name.
func LegalMoves(v *SafeView) []Macro {
	if v == nil || v.Me == nil {
		return nil
	}
	me := v.Me.PlayerID

	// Mulligan decisions pre-empt everything: until a seat has kept, it has no
	// board and no turn.
	if !v.Me.KeptHand {
		return []Macro{
			macro(KindKeepHand, fmt.Sprintf("Keep this hand (%d cards)", len(v.Me.Hand)),
				"KEEP_HAND:"+me),
			macro(KindMulligan, fmt.Sprintf("Mulligan (%d so far)", v.Me.MulliganCount),
				"MULLIGAN:"+me),
		}
	}

	var out []Macro

	mine := controlledBy(v.Battlefield, me)
	sources := untappedManaSources(mine)

	out = append(out, handMoves(v, me, sources)...)
	out = append(out, battlefieldMoves(mine)...)
	out = append(out, combatMoves(v, mine)...)
	out = append(out, lifeMoves(v, me)...)

	if v.Me.LibraryCount > 0 {
		out = append(out, macro(KindDraw,
			fmt.Sprintf("Draw a card (%d left in library)", v.Me.LibraryCount),
			"DRAW:"+me+":1"))
	}

	if v.ActivePlayerID == me {
		out = append(out, macro(KindPassTurn, "Pass the turn", "NEXT_TURN:"+me))
	}

	return out
}

// handMoves covers everything that starts in hand: land drops, casts, and
// discards.
func handMoves(v *SafeView, me string, sources []*SafeCard) []Macro {
	var out []Macro
	for _, c := range v.Me.Hand {
		switch {
		case IsLand(c):
			// A land drop is a move to the battlefield, no cost.
			out = append(out, macro(KindPlayLand,
				"Play land: "+displayName(c),
				"MOVE:"+c.ID+":"+ZoneBattlefield))

		default:
			// "Cast" in a rules-light engine: pay by tapping, put the card on
			// the stack so the move is visible to the other seats, then resolve
			// it to wherever it belongs. Two moves rather than one is
			// deliberate -- each step broadcasts, so a watching client sees the
			// card hit the stack before it resolves.
			taps, ok := PayMana(c.ManaCost, sources)
			if !ok {
				continue
			}
			dest := ZoneBattlefield
			if isInstantOrSorcery(c) {
				dest = ZoneGraveyard
			}
			steps := make([]string, 0, len(taps)+2)
			names := make([]string, 0, len(taps))
			for _, land := range taps {
				steps = append(steps, "TAP:"+land.ID)
				names = append(names, displayName(land))
			}
			steps = append(steps, "MOVE:"+c.ID+":"+ZoneStack, "MOVE:"+c.ID+":"+dest)

			label := "Cast " + displayName(c)
			if c.ManaCost != "" {
				label += " " + c.ManaCost
			}
			if len(names) > 0 {
				label += " (tapping " + strings.Join(names, ", ") + ")"
			} else {
				label += " (free)"
			}
			out = append(out, macro(KindCast, label, steps...))
		}

		// Discarding is always available to a playtester.
		out = append(out, macro(KindMoveZone,
			"Discard "+displayName(c),
			"MOVE:"+c.ID+":"+ZoneGraveyard))
	}
	return out
}

// battlefieldMoves covers tapping, untapping and destroying one's own
// permanents.
func battlefieldMoves(mine []*SafeCard) []Macro {
	var out []Macro
	for _, c := range mine {
		if c.Tapped {
			out = append(out, macro(KindUntap, "Untap "+displayName(c), "UNTAP:"+c.ID))
		} else {
			out = append(out, macro(KindTap, "Tap "+displayName(c), "TAP:"+c.ID))
		}
		out = append(out, macro(KindMoveZone,
			"Put "+displayName(c)+" into the graveyard",
			"MOVE:"+c.ID+":"+ZoneGraveyard))
	}
	return out
}

// combatMoves offers a single all-in attack per living opponent.
//
// The engine has no combat: Card.Attacking and Card.Blocking are independent
// booleans with no attacker-to-blocker mapping and nothing ever reads them
// (internal/game/game_state.go:81-83). What a playtester actually does is tap
// their creatures and subtract the damage by hand, so that is what this emits:
// tap every untapped, non-summoning-sick creature, then MODIFY_LIFE the
// defender by the negated total power.
//
// Blocking is therefore not modelled at all. That is a Phase 7 problem; here it
// only needs to be a damage source strong enough that random games terminate.
func combatMoves(v *SafeView, mine []*SafeCard) []Macro {
	attackers := make([]*SafeCard, 0, len(mine))
	power := 0
	for _, c := range mine {
		if c.Tapped || c.SummoningSickness || !IsCreature(c) {
			continue
		}
		p := parsePower(c.Power)
		if p <= 0 {
			continue
		}
		attackers = append(attackers, c)
		power += p
	}
	if len(attackers) == 0 || power == 0 {
		return nil
	}

	var out []Macro
	for _, o := range livingOpponents(v) {
		steps := make([]string, 0, len(attackers)+1)
		for _, c := range attackers {
			steps = append(steps, "TAP:"+c.ID)
		}
		steps = append(steps, "MODIFY_LIFE:"+o.PlayerID+":"+strconv.Itoa(-power))
		out = append(out, macro(KindAttack,
			fmt.Sprintf("Attack %s with %d creature(s) for %d", o.Name, len(attackers), power),
			steps...))
	}
	return out
}

// lifeMoves covers the manual life-total adjustments a playtester makes for
// effects the engine does not implement.
func lifeMoves(v *SafeView, me string) []Macro {
	out := []Macro{
		macro(KindModifyLife, "Pay 1 life", "MODIFY_LIFE:"+me+":-1"),
		macro(KindModifyLife, "Gain 1 life", "MODIFY_LIFE:"+me+":1"),
	}
	for _, o := range livingOpponents(v) {
		out = append(out, macro(KindModifyLife,
			fmt.Sprintf("Deal 1 damage to %s", o.Name),
			"MODIFY_LIFE:"+o.PlayerID+":-1"))
	}
	return out
}

// livingOpponents returns opponents at a positive life total, in a stable
// order.
//
// The sort is load-bearing. GameEngine.buildGameView fills Opponents by ranging
// over the Players map (internal/game/view.go:88), so the slice arrives in a
// different order every call. Without this, a seeded policy picking from
// LegalMoves would still produce different games run to run.
func livingOpponents(v *SafeView) []*SafeOpponentView {
	out := make([]*SafeOpponentView, 0, len(v.Opponents))
	for _, o := range v.Opponents {
		if o != nil && o.Life > 0 {
			out = append(out, o)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].PlayerID < out[j].PlayerID })
	return out
}

func controlledBy(cards []*SafeCard, playerID string) []*SafeCard {
	out := make([]*SafeCard, 0, len(cards))
	for _, c := range cards {
		if c != nil && c.ControllerID == playerID {
			out = append(out, c)
		}
	}
	return out
}

func displayName(c *SafeCard) string {
	if c.DisplayName != "" {
		return c.DisplayName
	}
	return c.Name
}

// Zone names, mirrored from internal/game/game_state.go:113-121 so this package
// does not have to import the engine just to spell "BATTLEFIELD".
const (
	ZoneLibrary     = "LIBRARY"
	ZoneHand        = "HAND"
	ZoneBattlefield = "BATTLEFIELD"
	ZoneGraveyard   = "GRAVEYARD"
	ZoneExile       = "EXILE"
	ZoneStack       = "STACK"
	ZoneCommand     = "COMMAND"
)

// IsLand reports whether c is a land, by type line when one is present and by
// basic-land name otherwise.
func IsLand(c *SafeCard) bool {
	if c == nil {
		return false
	}
	if c.Type != "" {
		return strings.Contains(strings.ToLower(c.Type), "land")
	}
	return IsBasicLandName(c.Name)
}

// IsCreature reports whether c is a creature.
func IsCreature(c *SafeCard) bool {
	return c != nil && strings.Contains(strings.ToLower(c.Type), "creature")
}

func isInstantOrSorcery(c *SafeCard) bool {
	t := strings.ToLower(c.Type)
	return strings.Contains(t, "instant") || strings.Contains(t, "sorcery")
}

func parsePower(s string) int {
	// Power can be "*", "1+*", "" -- anything not a plain integer counts as
	// zero rather than guessing.
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return n
}

// Enrich fills in the printed characteristics the engine never sets.
//
// StartGameWithDecks builds cards with nothing but Name/DisplayName/OwnerID
// (internal/game/game_engine.go:94-114): ManaCost, Type, Power, Toughness and
// RulesText are all empty. LegalMoves needs all of them -- it cannot tell a
// land from a creature, price a spell, or know what a land taps for otherwise.
//
// Enrich mutates the SafeView in place. That is safe precisely because Redact
// already deep-copied everything (anti-pattern 3): nothing here can reach live
// engine state. Fields the engine did populate are left alone.
func Enrich(ctx context.Context, v *SafeView, o OracleLookup) {
	if v == nil || o == nil {
		return
	}
	seen := make(map[*SafeCard]bool)
	apply := func(cards []*SafeCard) {
		for _, c := range cards {
			if c == nil || seen[c] {
				continue
			}
			seen[c] = true
			oc, ok := o.Oracle(ctx, c.Name)
			if !ok {
				continue
			}
			if c.ManaCost == "" {
				c.ManaCost = oc.ManaCost
			}
			if c.Type == "" {
				c.Type = oc.TypeLine
			}
			if c.Power == "" {
				c.Power = oc.Power
			}
			if c.Toughness == "" {
				c.Toughness = oc.Toughness
			}
			if c.Loyalty == "" {
				c.Loyalty = oc.Loyalty
			}
			if c.RulesText == "" {
				c.RulesText = oc.OracleText
			}
		}
	}
	apply(v.Battlefield)
	apply(v.Exile)
	apply(v.Stack)
	apply(v.Command)
	if v.Me != nil {
		apply(v.Me.Hand)
		apply(v.Me.Graveyard)
	}
	for _, opp := range v.Opponents {
		if opp == nil {
			continue
		}
		apply(opp.Graveyard)
		if opp.TopCard != nil {
			apply([]*SafeCard{opp.TopCard})
		}
	}
}

// ---------------------------------------------------------------------------
// Mana solver
// ---------------------------------------------------------------------------

// ManaColors is the set of colors a source can produce. Order is WUBRG then
// colorless, matching the mana-cost symbol order.
const manaColors = "WUBRGC"

// ManaCost is a parsed mana cost.
type ManaCost struct {
	Generic int
	// Pips holds one entry per colored/colorless symbol required, each a set of
	// acceptable colors: a plain {G} is "G", a hybrid {G/W} is "GW". A hybrid
	// pip is satisfied by any one of its colors.
	Pips []string
}

// ParseManaCost parses a printed mana cost such as "{3}{G}{G}" or "{2}{U/R}".
//
// Unparseable symbols degrade to one generic each rather than failing, and {X}
// is treated as zero: this is a playtest sandbox, and refusing to offer a cast
// because a symbol was exotic is worse than offering a slightly mispriced one.
func ParseManaCost(s string) ManaCost {
	var mc ManaCost
	for _, sym := range manaSymbols(s) {
		up := strings.ToUpper(sym)
		switch {
		case up == "X" || up == "Y" || up == "Z":
			// {X} in a cost the caster chooses; zero is a legal choice.
		case isDigits(up):
			n, _ := strconv.Atoi(up)
			mc.Generic += n
		default:
			// Split hybrid / phyrexian into acceptable colors.
			var accept []byte
			for _, part := range strings.Split(up, "/") {
				part = strings.TrimSuffix(part, "P")
				if len(part) == 1 && strings.ContainsAny(part, manaColors) {
					accept = append(accept, part[0])
				} else if isDigits(part) {
					// {2/W}: the generic half is handled by treating the pip as
					// any color; close enough for a playtest solver.
					accept = append(accept, []byte(manaColors)...)
				}
			}
			if len(accept) == 0 {
				mc.Generic++
				continue
			}
			mc.Pips = append(mc.Pips, string(accept))
		}
	}
	return mc
}

// CMC is the converted mana cost.
func (m ManaCost) CMC() int { return m.Generic + len(m.Pips) }

func manaSymbols(s string) []string {
	var out []string
	for {
		i := strings.IndexByte(s, '{')
		if i < 0 {
			return out
		}
		j := strings.IndexByte(s[i:], '}')
		if j < 0 {
			return out
		}
		out = append(out, s[i+1:i+j])
		s = s[i+j+1:]
	}
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// manaSource is an untapped permanent that can be tapped for mana, with the
// colors it produces.
type manaSource struct {
	card   *SafeCard
	colors string // subset of manaColors; "" means it produces nothing usable
}

// untappedManaSources finds the seat's untapped mana producers, in a
// deterministic order.
func untappedManaSources(mine []*SafeCard) []*SafeCard {
	out := make([]*SafeCard, 0, len(mine))
	for _, c := range mine {
		if c.Tapped {
			continue
		}
		if producedColors(c) == "" {
			continue
		}
		out = append(out, c)
	}
	// Battlefield order is already stable (append-only), but sort anyway so
	// the solver's choice does not depend on play order.
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].ID < out[j].ID
	})
	return out
}

// basicLandColor maps the basic land names to what they tap for.
var basicLandColor = map[string]string{
	"Plains": "W", "Snow-Covered Plains": "W",
	"Island": "U", "Snow-Covered Island": "U",
	"Swamp": "B", "Snow-Covered Swamp": "B",
	"Mountain": "R", "Snow-Covered Mountain": "R",
	"Forest": "G", "Snow-Covered Forest": "G",
	"Wastes": "C",
}

// producedColors reports which colors c can add, as a subset of manaColors.
//
// Basic lands are resolved by name; everything else is read out of the rules
// text, which is where Enrich put the oracle text. This is a text scan, not a
// rules implementation: it looks for "add" followed by mana symbols, which
// covers the overwhelming majority of real mana sources and quietly ignores the
// rest. A source it cannot read produces nothing and is never tapped for a
// cast, which is the safe direction to be wrong in.
func producedColors(c *SafeCard) string {
	if c == nil {
		return ""
	}
	if col, ok := basicLandColor[c.Name]; ok {
		return col
	}
	text := strings.ToLower(c.RulesText)
	var set []byte
	for _, idx := range addOffsets(text) {
		// Scan the symbols in the clause following "add", stopping at the end
		// of the sentence.
		clause := text[idx:]
		if end := strings.IndexAny(clause, ".;\n"); end >= 0 {
			clause = clause[:end]
		}
		for _, sym := range manaSymbols(clause) {
			up := strings.ToUpper(sym)
			for i := 0; i < len(up); i++ {
				if strings.IndexByte(manaColors, up[i]) >= 0 && !bytesContain(set, up[i]) {
					set = append(set, up[i])
				}
			}
		}
		if strings.Contains(clause, "any color") || strings.Contains(clause, "one mana of any") {
			set = []byte(manaColors)
			break
		}
	}
	sort.Slice(set, func(i, j int) bool {
		return strings.IndexByte(manaColors, set[i]) < strings.IndexByte(manaColors, set[j])
	})
	return string(set)
}

func addOffsets(text string) []int {
	var out []int
	for i := 0; ; {
		j := strings.Index(text[i:], "add ")
		if j < 0 {
			return out
		}
		out = append(out, i+j)
		i += j + 4
	}
}

func bytesContain(b []byte, want byte) bool {
	for _, x := range b {
		if x == want {
			return true
		}
	}
	return false
}

// PayMana solves for a set of sources whose combined output pays cost.
//
// It returns the sources to tap, in the order they should be tapped, and
// whether a payment exists at all. A free spell (empty or zero cost) returns
// (nil, true).
//
// This is deliberately a solver and not a decision. Which Forest pays for a
// {G} has exactly one right answer class and asking a policy -- or an LLM --
// to choose is pure token waste plus a source of illegal plays. Colored pips
// are assigned first, most-constrained-first, with backtracking; whatever is
// left over pays the generic portion.
func PayMana(cost string, available []*SafeCard) ([]*SafeCard, bool) {
	mc := ParseManaCost(cost)
	if mc.CMC() == 0 {
		return nil, true
	}
	if len(available) < mc.CMC() {
		return nil, false
	}

	srcs := make([]manaSource, 0, len(available))
	for _, c := range available {
		srcs = append(srcs, manaSource{card: c, colors: producedColors(c)})
	}

	// Most-constrained pip first: a pip with one acceptable color must be
	// placed before a hybrid that any of several sources could cover.
	order := make([]int, len(mc.Pips))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(a, b int) bool {
		return len(mc.Pips[order[a]]) < len(mc.Pips[order[b]])
	})

	used := make([]bool, len(srcs))
	chosen := make([]*SafeCard, 0, mc.CMC())

	var assign func(k int) bool
	assign = func(k int) bool {
		if k == len(order) {
			// Colored pips are covered; pay the generic from anything left.
			need := mc.Generic
			for i := range srcs {
				if need == 0 {
					break
				}
				if used[i] {
					continue
				}
				used[i] = true
				chosen = append(chosen, srcs[i].card)
				need--
			}
			return need == 0
		}
		want := mc.Pips[order[k]]
		for i := range srcs {
			if used[i] || !anyColorIn(srcs[i].colors, want) {
				continue
			}
			used[i] = true
			chosen = append(chosen, srcs[i].card)
			if assign(k + 1) {
				return true
			}
			used[i] = false
			chosen = chosen[:len(chosen)-1]
		}
		return false
	}

	if !assign(0) {
		return nil, false
	}
	// Tap in a stable order so identical situations produce identical steps.
	sort.SliceStable(chosen, func(i, j int) bool {
		if chosen[i].Name != chosen[j].Name {
			return chosen[i].Name < chosen[j].Name
		}
		return chosen[i].ID < chosen[j].ID
	})
	return chosen, true
}

func anyColorIn(have, want string) bool {
	for i := 0; i < len(want); i++ {
		if strings.IndexByte(have, want[i]) >= 0 {
			return true
		}
	}
	return false
}
