package bot

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// ShortIDRegistry is a Go port of reference/ShortIdRegistry.java (upstream
// Mage/src/main/java/mage/util/ShortIdRegistry.java). The class doc there
// states the invariants this type must preserve; they are reproduced with the
// code that implements each one.
//
// Purpose: engine card IDs are long ("game-1-alice-17") and burn tokens on
// every board render. Short IDs ("p1", "p2", ...) are what the model sees and
// what it sends back in choice / attackers / blockers / mana_plan parameters.
//
// Invariants:
//
//   - Monotonic counter from 1, prefix "p" (Java: nextId = new AtomicInteger(1),
//     ShortIdRegistry() -> this("p")).
//   - IDs are stable as a card changes zones. This falls out of keying on the
//     card ID rather than on anything zone-dependent; game.Card.ID is assigned
//     once at deck load (internal/game/game_engine.go:98) and never rewritten.
//   - Superseded IDs stay resolvable as aliases. See Register.
//   - Assignment order is DETERMINISTIC, sorted by (name, sequence). Never by
//     card ID. Java class doc: "Never use UUID as a sort key." Our card IDs
//     embed an insertion index rather than a UUID, but sorting by them would
//     still make short IDs depend on deck order and library shuffles, which is
//     exactly the nondeterminism the invariant exists to prevent.
//
// Thread-safe, like the Java original (ConcurrentHashMap + AtomicInteger there,
// one mutex here -- the maps are small and never hot).
type ShortIDRegistry struct {
	mu       sync.Mutex
	prefix   string
	idToItem map[string]string // card ID -> short ID
	shortTo  map[string]string // short ID -> card ID, including resolve-only aliases
	nextID   int
}

// DefaultShortIDPrefix is the prefix the server uses (Java: the no-arg ctor).
// The bridge client uses "l" for locally-assigned fallback IDs; we have no
// separate client, so only "p" is produced here.
const DefaultShortIDPrefix = "p"

// NewShortIDRegistry creates a registry with the default "p" prefix.
func NewShortIDRegistry() *ShortIDRegistry {
	return NewShortIDRegistryWithPrefix(DefaultShortIDPrefix)
}

// NewShortIDRegistryWithPrefix creates a registry with a custom prefix.
func NewShortIDRegistryWithPrefix(prefix string) *ShortIDRegistry {
	if prefix == "" {
		panic("bot: short ID prefix must not be empty")
	}
	return &ShortIDRegistry{
		prefix:   prefix,
		idToItem: make(map[string]string),
		shortTo:  make(map[string]string),
		nextID:   1,
	}
}

// GetOrAssign returns the short ID for a card ID, assigning one on first
// encounter. Port of ShortIdRegistry.getOrAssign.
func (r *ShortIDRegistry) GetOrAssign(cardID string) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.getOrAssignLocked(cardID)
}

func (r *ShortIDRegistry) getOrAssignLocked(cardID string) string {
	if existing, ok := r.idToItem[cardID]; ok {
		return existing
	}
	short := r.prefix + strconv.Itoa(r.nextID)
	r.nextID++
	r.idToItem[cardID] = short
	r.shortTo[short] = cardID
	return short
}

// Sequence returns the numeric part of a card's short ID, or math.MaxInt when
// the card has not been assigned one. Side-effect free, so it is safe inside a
// comparator. Port of ShortIdRegistry.getSequence.
func (r *ShortIDRegistry) Sequence(cardID string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	short, ok := r.idToItem[cardID]
	if !ok {
		return unassignedSequence
	}
	n, err := ParseSequence(short)
	if err != nil {
		return unassignedSequence
	}
	return n
}

// unassignedSequence sorts unassigned cards last, matching Java's
// Integer.MAX_VALUE sentinel.
const unassignedSequence = int(^uint(0) >> 1)

// Resolve maps a short ID back to a card ID. Port of ShortIdRegistry.resolve.
func (r *ShortIDRegistry) Resolve(shortID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cardID, ok := r.shortTo[shortID]
	if !ok {
		return "", fmt.Errorf("bot: unknown short ID: %s", shortID)
	}
	return cardID, nil
}

// TryResolve is the non-erroring variant of Resolve.
func (r *ShortIDRegistry) TryResolve(shortID string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cardID, ok := r.shortTo[shortID]
	return cardID, ok
}

// Register records an externally-assigned short ID for a card.
//
// The external assignment is authoritative: if the card already had a
// locally-assigned ID, the mapping is updated, but the OLD short ID stays in
// the reverse map as a resolve-only alias. That is the point -- a model that
// referenced "p4" three decisions ago (in a mana plan, say) must still resolve
// it to the same card after the ID was superseded. Port of
// ShortIdRegistry.register.
func (r *ShortIDRegistry) Register(cardID, shortID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.idToItem[cardID]; ok {
		if existing == shortID {
			return
		}
		// Old short ID deliberately left in r.shortTo as an alias.
		r.idToItem[cardID] = shortID
		r.shortTo[shortID] = cardID
		r.advanceNextIDLocked(shortID)
		return
	}

	if existingCard, ok := r.shortTo[shortID]; ok && existingCard != cardID {
		// Same short ID claimed by two different cards. Java logs SEVERE here
		// and evicts; we do the same, minus the logger dependency.
		delete(r.idToItem, existingCard)
		delete(r.shortTo, shortID)
	}

	r.idToItem[cardID] = shortID
	r.shortTo[shortID] = cardID
	r.advanceNextIDLocked(shortID)
}

func (r *ShortIDRegistry) advanceNextIDLocked(shortID string) {
	n, err := ParseSequence(shortID)
	if err != nil {
		return // non-standard format, ignore (Java swallows NumberFormatException)
	}
	if n+1 > r.nextID {
		r.nextID = n + 1
	}
}

// AssignAll assigns short IDs to every not-yet-assigned card in cards, in the
// deterministic (name, sequence) order the class doc mandates. Call it once per
// render, before serializing, so that IDs do not depend on slice order.
//
// Already-assigned cards keep their IDs (stability across zones), and sort
// ahead of new cards because their sequence is lower than the unassigned
// sentinel.
func (r *ShortIDRegistry) AssignAll(cards []*SafeCard) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ordered := make([]*SafeCard, 0, len(cards))
	for _, c := range cards {
		if c != nil {
			ordered = append(ordered, c)
		}
	}
	sort.SliceStable(ordered, func(i, j int) bool {
		return r.lessLocked(ordered[i], ordered[j])
	})
	for _, c := range ordered {
		r.getOrAssignLocked(c.ID)
	}
}

// lessLocked implements the canonical sort key (name, sequence). Ties on both
// are broken by card ID only to keep the sort total -- never as a primary key.
func (r *ShortIDRegistry) lessLocked(a, b *SafeCard) bool {
	an, bn := a.Name, b.Name
	if an != bn {
		return an < bn
	}
	as, bs := r.sequenceLocked(a.ID), r.sequenceLocked(b.ID)
	if as != bs {
		return as < bs
	}
	return a.ID < b.ID
}

func (r *ShortIDRegistry) sequenceLocked(cardID string) int {
	short, ok := r.idToItem[cardID]
	if !ok {
		return unassignedSequence
	}
	n, err := ParseSequence(short)
	if err != nil {
		return unassignedSequence
	}
	return n
}

// SnapshotShortIDs returns every currently-known short ID, aliases included.
// Port of ShortIdRegistry.snapshotShortIds.
func (r *ShortIDRegistry) SnapshotShortIDs() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, 0, len(r.shortTo))
	for s := range r.shortTo {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// PeekNextID reports the next-ID counter, for diagnostics.
func (r *ShortIDRegistry) PeekNextID() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.nextID
}

// DumpAssignments renders "[p1=..., p2=...]" sorted by sequence, for logging
// when diagnosing nondeterministic assignment. Port of
// ShortIdRegistry.dumpAssignments.
func (r *ShortIDRegistry) DumpAssignments() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	type pair struct {
		short, card string
		seq         int
	}
	pairs := make([]pair, 0, len(r.idToItem))
	for card, short := range r.idToItem {
		seq, err := ParseSequence(short)
		if err != nil {
			seq = unassignedSequence
		}
		pairs = append(pairs, pair{short: short, card: card, seq: seq})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].seq != pairs[j].seq {
			return pairs[i].seq < pairs[j].seq
		}
		return pairs[i].short < pairs[j].short
	})
	var b strings.Builder
	b.WriteByte('[')
	for i, p := range pairs {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%s=%s", p.short, p.card)
	}
	b.WriteByte(']')
	return b.String()
}

// Clear resets every mapping. Call on game start. Port of ShortIdRegistry.clear.
func (r *ShortIDRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.idToItem = make(map[string]string)
	r.shortTo = make(map[string]string)
	r.nextID = 1
}

// ParseSequence extracts the numeric suffix of a short ID ("p6" -> 6).
// Port of ShortIdRegistry.parseSequence, which throws where this returns an error.
func ParseSequence(shortID string) (int, error) {
	if len(shortID) < 2 {
		return 0, fmt.Errorf("bot: malformed short ID: %q", shortID)
	}
	n, err := strconv.Atoi(shortID[1:])
	if err != nil {
		return 0, fmt.Errorf("bot: malformed short ID: %q", shortID)
	}
	return n, nil
}
