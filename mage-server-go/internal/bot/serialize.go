package bot

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// serialize.go is a Go port of mage-bench's decision renderer:
//
//	reference/decision_renderer.py:133-248  _render_decision_block
//	reference/decision_renderer.py:254-291  _render_board
//	reference/decision_renderer.py:304-342  permanent_display
//	reference/decision_renderer.py:459-481  format_choice
//	reference/decision_renderer.py:714-772  _render_card_reference
//	reference/pilot_rendering.py:71-95      Respond / Mana pool / Recent chat
//	reference/pilot_rendering.py:61-64      oracle dedup
//
// The field list and emission order are specified in BOT_PLAYERS_PLAN.md §0.2.
// Where the local engine cannot supply a field that mage-bench's bridge does,
// the field is carried on Decision (caller-supplied) rather than invented here.

// basicLandNames is the Card Reference exclusion set. Single package-level
// definition -- these names never get an oracle entry, because every model
// already knows what a Mountain does and eleven oracle lines per render is the
// difference between a cheap prompt and an expensive one.
// reference/decision_renderer.py:38-52.
var basicLandNames = map[string]struct{}{
	"Plains":                {},
	"Island":                {},
	"Swamp":                 {},
	"Mountain":              {},
	"Forest":                {},
	"Wastes":                {},
	"Snow-Covered Plains":   {},
	"Snow-Covered Island":   {},
	"Snow-Covered Swamp":    {},
	"Snow-Covered Mountain": {},
	"Snow-Covered Forest":   {},
}

// IsBasicLandName reports whether name is excluded from the Card Reference.
func IsBasicLandName(name string) bool {
	_, ok := basicLandNames[name]
	return ok
}

// Choice is one option offered to the deciding player.
// reference/decision_renderer.py:459-481 (format_choice).
type Choice struct {
	Name        string
	Description string // fallback when Name is empty
	ID          string // short ID, e.g. "p3"
	Action      string // e.g. "cast", "activate"
	ManaCost    string
}

// MultiAmountItem is one row of a "distribute N among these" decision.
// reference/decision_renderer.py:219-233.
type MultiAmountItem struct {
	Description string
	Min         *int
	Max         *int
}

// CombatGroup is one attacker-plus-blockers unit.
// reference/decision_renderer.py:399-421 (_render_combat).
//
// The local engine tracks Attacking/Blocking as independent booleans on cards
// (internal/game/game_state.go:81-83) with no attacker-to-blocker mapping, so
// combat groups cannot be derived from a SafeView. They are supplied by the
// caller that owns the combat decision.
type CombatGroup struct {
	Attackers []string
	Blockers  []string
	Blocked   bool
	Defending string
}

// IncomingAttacker is one attacker a blocking player must answer.
// reference/decision_renderer.py:431-456.
type IncomingAttacker struct {
	Name           string
	ID             string
	PowerToughness string
}

// PilotContext is the per-decision overlay mage-bench's bridge attaches.
// Pointer fields distinguish "absent" from "zero", which the renderer relies on
// (reference/decision_renderer.py:196: has_field, not truthiness).
type PilotContext struct {
	CombatPhase       string
	IncomingAttackers []IncomingAttacker
	UntappedLands     *int
	LandDropsUsed     *int
}

// Decision is everything about the pending decision that is not board state.
type Decision struct {
	Index         int
	SnapshotIndex int
	Turn          int
	Phase         string // rendered as "PREGAME" when empty
	Player        string // display name of the deciding seat
	Message       string
	Choices       []Choice
	Items         []MultiAmountItem
	TotalMin      *int
	TotalMax      *int
	Combat        []CombatGroup
	PilotContext  *PilotContext

	// Trailing lines. reference/pilot_rendering.py:71-95.
	RespondWith  string
	ResponseType string // used only when RespondWith is empty
	RecentChat   []string
}

// Serializer renders SafeViews as mage-bench board text.
//
// It is stateful across calls on purpose: seenOracle implements the
// first-appearance-only Card Reference dedup that
// reference/pilot_rendering.py:61-64 describes, and which is the single largest
// token saving in the whole design. One Serializer per bot per game.
type Serializer struct {
	oracle     OracleLookup
	ids        *ShortIDRegistry
	seenOracle map[string]struct{}
}

// NewSerializer creates a Serializer. oracle may be nil, in which case no Card
// Reference is emitted and permanents fall back to whatever P/T the engine set.
func NewSerializer(oracle OracleLookup) *Serializer {
	return &Serializer{
		oracle:     oracle,
		ids:        NewShortIDRegistry(),
		seenOracle: make(map[string]struct{}),
	}
}

// IDs exposes the short-ID registry so a caller can resolve a model's "p7"
// back to an engine card ID.
func (s *Serializer) IDs() *ShortIDRegistry { return s.ids }

// ResetOracleDedup forgets which cards have already had their oracle text
// printed. Call it whenever the conversation context is reset, since the model
// no longer has the earlier Card Reference sections in view.
// reference/pilot_state.py resets the same set on context reset.
func (s *Serializer) ResetOracleDedup() {
	s.seenOracle = make(map[string]struct{})
}

// Render produces the full prompt fragment: an optional "## Card Reference"
// section, the "## Decision" block, and the trailing Respond / Mana pool /
// Recent chat lines.
func (s *Serializer) Render(ctx context.Context, v *SafeView, d *Decision) string {
	r := &renderState{ctx: ctx, s: s, v: v, d: d}
	r.assignShortIDs()

	var parts []string
	if ref := r.cardReference(); ref != "" {
		parts = append(parts, ref)
	}
	parts = append(parts, "## Decision\n\n"+r.decisionBlock())

	out := strings.Join(parts, "\n\n")

	// Trailing lines, reference/pilot_rendering.py:71-95.
	var tail []string
	if d.RespondWith != "" {
		respond := d.RespondWith
		if d.TotalMin != nil && d.TotalMax != nil && *d.TotalMin == *d.TotalMax {
			respond = strings.ReplaceAll(respond,
				"sum between total_min and total_max",
				fmt.Sprintf("sum must equal total (%d)", *d.TotalMin))
		}
		tail = append(tail, "  Respond: "+respond)
	} else if d.ResponseType != "" {
		tail = append(tail, "  Response type: "+d.ResponseType)
	}
	if pool := r.manaPoolLine(); pool != "" {
		tail = append(tail, pool)
	}
	if len(d.RecentChat) > 0 {
		tail = append(tail, "  Recent chat: "+strings.Join(d.RecentChat, " | "))
	}
	if len(tail) > 0 {
		out += "\n" + strings.Join(tail, "\n")
	}
	return out
}

// renderState carries the per-render scratch: the context for oracle lookups
// and the view/decision being rendered.
type renderState struct {
	ctx context.Context
	s   *Serializer
	v   *SafeView
	d   *Decision
}

// assignShortIDs gives every visible card a stable short ID before anything is
// rendered, in the deterministic (name, sequence) order ShortIDRegistry
// mandates. Doing it up front means IDs never depend on which zone the
// renderer happened to walk first.
func (r *renderState) assignShortIDs() {
	if r.v == nil {
		return
	}
	var all []*SafeCard
	all = append(all, r.v.Battlefield...)
	all = append(all, r.v.Stack...)
	all = append(all, r.v.Command...)
	all = append(all, r.v.Exile...)
	if r.v.Me != nil {
		all = append(all, r.v.Me.Hand...)
		all = append(all, r.v.Me.Graveyard...)
	}
	for _, o := range r.v.Opponents {
		all = append(all, o.Graveyard...)
		if o.TopCard != nil {
			all = append(all, o.TopCard)
		}
	}
	r.s.ids.AssignAll(all)
}

// ---------------------------------------------------------------------------
// Decision block. Port of _render_decision_block, decision_renderer.py:133-248.
// ---------------------------------------------------------------------------

func (r *renderState) decisionBlock() string {
	d := r.d

	phase := d.Phase
	if phase == "" {
		// decision_renderer.py:143-146: empty phase is legal only pregame.
		phase = "PREGAME"
	}

	lines := []string{
		fmt.Sprintf("[Decision %d, snapshot=%d] Turn %d %s - %s",
			d.Index, d.SnapshotIndex, d.Turn, phase, d.Player),
	}

	lines = append(lines, "  Board: "+r.board())

	if r.v != nil && len(r.v.Stack) > 0 {
		items := make([]string, 0, len(r.v.Stack))
		for _, c := range r.v.Stack {
			items = append(items, cardDisplay(c))
		}
		lines = append(lines, "  Stack: ["+strings.Join(items, ", ")+"]")
	}

	if len(d.Combat) > 0 {
		lines = append(lines, "  Combat: "+renderCombat(d.Combat))
	}

	var pc *PilotContext
	if d.PilotContext != nil {
		pc = d.PilotContext
	}

	if pc != nil && pc.CombatPhase != "" {
		lines = append(lines, "  Combat Phase: "+pc.CombatPhase)
	}

	if pc != nil && isDeclareBlockersPhase(pc.CombatPhase) && len(pc.IncomingAttackers) > 0 {
		lines = append(lines, "  Incoming Attackers: "+renderIncomingAttackers(pc.IncomingAttackers))
	}

	// Pointer-nil is "field absent"; a present zero still renders.
	// decision_renderer.py:196.
	if pc != nil && (pc.UntappedLands != nil || pc.LandDropsUsed != nil) {
		var ctxParts []string
		if pc.UntappedLands != nil {
			ctxParts = append(ctxParts, fmt.Sprintf("Untapped lands: %d", *pc.UntappedLands))
		}
		if pc.LandDropsUsed != nil {
			ctxParts = append(ctxParts, fmt.Sprintf("Land drops remaining: %d", 1-*pc.LandDropsUsed))
		}
		lines = append(lines, "  "+strings.Join(ctxParts, ", "))
	}

	// Message is always emitted, even when empty. decision_renderer.py:214.
	lines = append(lines, "  Message: "+d.Message)

	if len(d.Items) > 0 {
		lines = append(lines, r.itemLines()...)
	} else {
		descs := make([]string, 0, len(d.Choices))
		for _, c := range d.Choices {
			descs = append(descs, formatChoice(c))
		}
		lines = append(lines, fmt.Sprintf("  Choices (%d): %s", len(d.Choices), strings.Join(descs, ", ")))
	}

	if strings.Contains(d.Message, "Pick triggered ability") {
		lines = append(lines,
			"  NOTE: This decision only determines the order triggered abilities"+
				" are placed on the stack. Targets are chosen in separate decisions.")
	}

	return strings.Join(lines, "\n")
}

// itemLines renders a multi-amount decision. decision_renderer.py:216-233.
func (r *renderState) itemLines() []string {
	d := r.d
	header := fmt.Sprintf("  Items (%d)", len(d.Items))
	switch {
	case d.TotalMin != nil && d.TotalMax != nil && *d.TotalMin == *d.TotalMax:
		header += fmt.Sprintf(": total=%d", *d.TotalMin)
	default:
		var totals []string
		if d.TotalMin != nil {
			totals = append(totals, fmt.Sprintf("total_min=%d", *d.TotalMin))
		}
		if d.TotalMax != nil {
			totals = append(totals, fmt.Sprintf("total_max=%d", *d.TotalMax))
		}
		if len(totals) > 0 {
			header += ": " + strings.Join(totals, ", ")
		}
	}
	lines := []string{header}
	for i, item := range d.Items {
		var constraints []string
		if item.Min != nil {
			constraints = append(constraints, "min="+strconv.Itoa(*item.Min))
		}
		if item.Max != nil {
			constraints = append(constraints, "max="+strconv.Itoa(*item.Max))
		}
		suffix := ""
		if len(constraints) > 0 {
			suffix = " [" + strings.Join(constraints, ", ") + "]"
		}
		lines = append(lines, fmt.Sprintf("    %d: %s%s", i, item.Description, suffix))
	}
	return lines
}

// ---------------------------------------------------------------------------
// Board. Port of _render_board, decision_renderer.py:254-291.
// ---------------------------------------------------------------------------

// board renders every seat, separated by " | ".
//
// Seat order is sorted by name. mage-bench takes snapshot player order from the
// bridge; the local engine stores players in a map[string]*Player and
// buildGameView ranges over it (internal/game/view.go:85), so there is no order
// to inherit -- Go would hand us a different one on every call and every render
// would diff. Sorting by name is the same determinism rule ShortIDRegistry uses.
func (r *renderState) board() string {
	if r.v == nil {
		return ""
	}
	type seat struct {
		name string
		line string
	}
	seats := make([]seat, 0, 1+len(r.v.Opponents))

	if me := r.v.Me; me != nil {
		var s strings.Builder
		hand := make([]string, 0, len(me.Hand))
		for _, c := range me.Hand {
			hand = append(hand, cardDisplay(c))
		}
		if len(hand) > 0 {
			fmt.Fprintf(&s, "%s: %dhp hand=[%s]", me.Name, me.Life, strings.Join(hand, ", "))
		} else {
			fmt.Fprintf(&s, "%s: %dhp hand=0", me.Name, me.Life)
		}
		fmt.Fprintf(&s, " lib=%d", me.LibraryCount)
		s.WriteString(playerCounters(me.Poison, me.Energy))
		s.WriteString(r.zoneSuffixes(me.PlayerID, me.Graveyard))
		seats = append(seats, seat{name: me.Name, line: s.String()})
	}

	for _, o := range r.v.Opponents {
		var s strings.Builder
		fmt.Fprintf(&s, "%s: %dhp", o.Name, o.Life)
		// Opponent hands are counts, and the count is omitted at zero.
		// decision_renderer.py:270-273.
		if o.HandCount != 0 {
			fmt.Fprintf(&s, " hand=%d", o.HandCount)
		}
		fmt.Fprintf(&s, " lib=%d", o.LibraryCount)
		s.WriteString(playerCounters(o.Poison, o.Energy))
		s.WriteString(r.zoneSuffixes(o.PlayerID, o.Graveyard))
		seats = append(seats, seat{name: o.Name, line: s.String()})
	}

	sort.SliceStable(seats, func(i, j int) bool { return seats[i].name < seats[j].name })
	lines := make([]string, 0, len(seats))
	for _, st := range seats {
		lines = append(lines, st.line)
	}
	return strings.Join(lines, " | ")
}

// zoneSuffixes renders " bf=[..] gy=[..] exile=[..]" for one seat.
//
// The engine keeps battlefield and exile as single global slices
// (internal/game/game_state.go), unlike mage-bench's per-player snapshot zones,
// so they are partitioned here: battlefield by controller (who it answers to
// now), exile by owner (nothing controls an exiled card).
func (r *renderState) zoneSuffixes(playerID string, graveyard []*SafeCard) string {
	var s strings.Builder
	if bf := filterByController(r.v.Battlefield, playerID); len(bf) > 0 {
		strs := make([]string, 0, len(bf))
		for _, c := range bf {
			strs = append(strs, r.permanentDisplay(c))
		}
		fmt.Fprintf(&s, " bf=[%s]", strings.Join(strs, ", "))
	}
	if len(graveyard) > 0 {
		strs := make([]string, 0, len(graveyard))
		for _, c := range graveyard {
			strs = append(strs, cardDisplay(c))
		}
		fmt.Fprintf(&s, " gy=[%s]", strings.Join(strs, ", "))
	}
	if ex := filterByOwner(r.v.Exile, playerID); len(ex) > 0 {
		strs := make([]string, 0, len(ex))
		for _, c := range ex {
			strs = append(strs, cardDisplay(c))
		}
		fmt.Fprintf(&s, " exile=[%s]", strings.Join(strs, ", "))
	}
	return s.String()
}

func filterByController(cards []*SafeCard, playerID string) []*SafeCard {
	out := make([]*SafeCard, 0, len(cards))
	for _, c := range cards {
		if c != nil && c.ControllerID == playerID {
			out = append(out, c)
		}
	}
	return out
}

func filterByOwner(cards []*SafeCard, playerID string) []*SafeCard {
	out := make([]*SafeCard, 0, len(cards))
	for _, c := range cards {
		if c != nil && c.OwnerID == playerID {
			out = append(out, c)
		}
	}
	return out
}

// playerCounters renders player-level counters. decision_renderer.py:345-364
// takes an arbitrary counter bag; the local engine models exactly two, as
// dedicated int fields (internal/game/view.go:36-37).
func playerCounters(poison, energy int) string {
	var s strings.Builder
	if poison != 0 {
		fmt.Fprintf(&s, " poison=%d", poison)
	}
	if energy != 0 {
		fmt.Fprintf(&s, " energy=%d", energy)
	}
	return s.String()
}

// cardDisplay renders a card in hand / graveyard / exile / stack: just a name.
// decision_renderer.py:297-301.
func cardDisplay(c *SafeCard) string {
	if c == nil {
		return "?"
	}
	if c.DisplayName != "" {
		return c.DisplayName
	}
	return c.Name
}

// permanentDisplay renders a battlefield permanent with its status annotations:
// "Name P/T (tapped, sick, face_down, loyalty=N, ctr=N, copy of X, token)".
// Port of permanent_display, decision_renderer.py:304-342 -- extras are
// collected in that exact order, then P/T is appended to the NAME (not to the
// extras) before the parenthesised list.
func (r *renderState) permanentDisplay(c *SafeCard) string {
	if c == nil {
		return "?"
	}
	name := cardDisplay(c)
	var extras []string
	if c.Tapped {
		extras = append(extras, "tapped")
	}
	if c.SummoningSickness {
		extras = append(extras, "sick")
	}
	if c.FaceDown {
		extras = append(extras, "face_down")
	}
	if c.Loyalty != "" {
		extras = append(extras, "loyalty="+c.Loyalty)
	}
	for _, ctr := range c.Counters {
		n := ctr.Name
		if n == "" {
			n = "?"
		}
		extras = append(extras, fmt.Sprintf("%s=%d", n, ctr.Count))
	}
	if c.OriginalCard != "" {
		extras = append(extras, "copy of "+c.OriginalCard)
	} else if c.IsCopy {
		extras = append(extras, "copy")
	}
	if isToken(c) {
		extras = append(extras, "token")
	}

	pt := ""
	if c.Power != "" {
		pt = c.Power + "/" + c.Toughness
	} else if o, ok := r.lookupOracle(name); ok {
		// game.Card.Power/Toughness are left empty by StartGameWithDecks
		// (game_engine.go:96-107), so without this fallback no creature on the
		// board would ever show a P/T.
		pt = o.PowerToughness()
	}
	if pt != "" {
		name += " " + pt
	}

	if len(extras) > 0 {
		return name + " (" + strings.Join(extras, ", ") + ")"
	}
	return name
}

// isToken reports whether a permanent is a token. The engine has no flag for
// it; CreateToken mints IDs prefixed "token-" (internal/game/actions.go:347)
// and MoveCard keys the cease-to-exist rule off the same prefix
// (internal/game/actions.go:133), so the prefix is the engine's own marker.
func isToken(c *SafeCard) bool {
	return c != nil && strings.HasPrefix(c.ID, "token-")
}

// renderCombat ports _render_combat, decision_renderer.py:399-421.
func renderCombat(groups []CombatGroup) string {
	parts := make([]string, 0, len(groups))
	for _, g := range groups {
		part := strings.Join(g.Attackers, ", ")
		if len(g.Blockers) > 0 {
			part += " blocked by " + strings.Join(g.Blockers, ", ")
		} else if g.Blocked {
			part += " (blocked)"
		}
		if g.Defending != "" {
			part += " -> " + g.Defending
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, " | ")
}

// isDeclareBlockersPhase ports _is_declare_blockers_phase,
// decision_renderer.py:424-429.
func isDeclareBlockersPhase(combatPhase string) bool {
	switch strings.ToLower(combatPhase) {
	case "blockers", "declare_blockers":
		return true
	}
	return false
}

// renderIncomingAttackers ports _render_incoming_attackers,
// decision_renderer.py:431-456.
func renderIncomingAttackers(attackers []IncomingAttacker) string {
	parts := make([]string, 0, len(attackers))
	for _, a := range attackers {
		var extras []string
		if a.ID != "" {
			extras = append(extras, "id="+a.ID)
		}
		if a.PowerToughness != "" {
			extras = append(extras, a.PowerToughness)
		}
		if len(extras) > 0 {
			parts = append(parts, a.Name+" ["+strings.Join(extras, ", ")+"]")
		} else {
			parts = append(parts, a.Name)
		}
	}
	return strings.Join(parts, ", ")
}

// formatChoice ports format_choice, decision_renderer.py:459-481.
func formatChoice(c Choice) string {
	name := c.Name
	if name == "" {
		name = c.Description
	}
	if name == "" {
		name = "?"
	}
	var attrs []string
	if c.ID != "" {
		attrs = append(attrs, "id="+c.ID)
	}
	if c.Action != "" {
		attrs = append(attrs, c.Action)
	}
	if c.ManaCost != "" {
		attrs = append(attrs, c.ManaCost)
	}
	if len(attrs) > 0 {
		return name + " [" + strings.Join(attrs, ", ") + "]"
	}
	return name
}

// manaPoolLine ports the mana pool trailing line,
// reference/pilot_rendering.py:87-90. Emitted only when the pool is non-empty,
// and only the non-zero colours appear.
func (r *renderState) manaPoolLine() string {
	if r.v == nil || r.v.Me == nil {
		return ""
	}
	p := r.v.Me.ManaPool
	if p.Total() == 0 {
		return ""
	}
	// Fixed WUBRG+C order: the engine's pool is a struct, so unlike
	// mage-bench's dict there is no insertion order to preserve, and a stable
	// order is what makes the goldens diffable.
	pairs := []struct {
		key string
		val int
	}{
		{"W", p.White}, {"U", p.Blue}, {"B", p.Black},
		{"R", p.Red}, {"G", p.Green}, {"C", p.Colorless},
	}
	var parts []string
	for _, kv := range pairs {
		if kv.val > 0 {
			parts = append(parts, fmt.Sprintf("%s=%d", kv.key, kv.val))
		}
	}
	return "  Mana pool: " + strings.Join(parts, ", ")
}

// ---------------------------------------------------------------------------
// Card Reference. Port of _render_card_reference, decision_renderer.py:714-772,
// plus the first-appearance dedup from pilot_rendering.py:61-64.
// ---------------------------------------------------------------------------

func (r *renderState) cardReference() string {
	if r.s.oracle == nil {
		return ""
	}

	names := make(map[string]struct{})
	add := func(cards []*SafeCard) {
		for _, c := range cards {
			if c == nil {
				continue
			}
			if n := cardDisplay(c); n != "" && n != "?" {
				names[n] = struct{}{}
			}
		}
	}
	if r.v != nil {
		if r.v.Me != nil {
			add(r.v.Me.Hand)
			add(r.v.Me.Graveyard)
		}
		for _, o := range r.v.Opponents {
			add(o.Graveyard)
			if o.TopCard != nil {
				add([]*SafeCard{o.TopCard})
			}
		}
		add(r.v.Battlefield)
		add(r.v.Exile)
		add(r.v.Stack)
		add(r.v.Command) // mage-bench's "commanders" zone
	}
	for _, c := range r.d.Choices {
		if c.Name != "" {
			names[c.Name] = struct{}{}
		}
	}

	sorted := make([]string, 0, len(names))
	for n := range names {
		sorted = append(sorted, n)
	}
	sort.Strings(sorted)

	var lines []string
	for _, name := range sorted {
		if IsBasicLandName(name) {
			continue
		}
		// First appearance only. The model keeps earlier turns in context, so
		// reprinting a card's oracle text every render buys nothing and is the
		// biggest single line item in the prompt.
		if _, seen := r.s.seenOracle[name]; seen {
			continue
		}
		o, ok := r.s.oracle.Oracle(r.ctx, name)
		if !ok || o == nil {
			continue
		}
		r.s.seenOracle[name] = struct{}{}

		entry := "- " + name
		if o.ManaCost != "" {
			entry += " " + o.ManaCost
		}
		if o.TypeLine != "" {
			entry += " -- " + o.TypeLine
		}
		if pt := o.PowerToughness(); pt != "" {
			entry += " " + pt
		}
		if o.OracleText != "" {
			entry += ": " + strings.ReplaceAll(o.OracleText, "\n", " / ")
		}
		lines = append(lines, entry)
	}

	if len(lines) == 0 {
		return ""
	}
	return "## Card Reference\n" + strings.Join(lines, "\n")
}

// lookupOracle is a nil-safe oracle lookup for the P/T fallback.
func (r *renderState) lookupOracle(name string) (*OracleCard, bool) {
	if r.s.oracle == nil {
		return nil, false
	}
	return r.s.oracle.Oracle(r.ctx, name)
}
