package bot

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShortIDAssignsFromOneWithPPrefix(t *testing.T) {
	r := NewShortIDRegistry()
	require.Equal(t, "p1", r.GetOrAssign("card-a"))
	require.Equal(t, "p2", r.GetOrAssign("card-b"))
	require.Equal(t, "p1", r.GetOrAssign("card-a"), "IDs must be stable for the same card")
	require.Equal(t, 3, r.PeekNextID())
}

// TestShortIDAssignmentIsByNameNotByCardID pins the invariant from
// reference/ShortIdRegistry.java:26-33.
func TestShortIDAssignmentIsByNameNotByCardID(t *testing.T) {
	cards := []*SafeCard{
		{ID: "zzz-9", Name: "Ancestral Vision"},
		{ID: "aaa-1", Name: "Zealous Conscripts"},
		{ID: "mmm-5", Name: "Mox Diamond"},
	}
	r := NewShortIDRegistry()
	r.AssignAll(cards)

	// Sorted by NAME: Ancestral Vision, Mox Diamond, Zealous Conscripts.
	require.Equal(t, "p1", r.GetOrAssign("zzz-9"))
	require.Equal(t, "p2", r.GetOrAssign("mmm-5"))
	require.Equal(t, "p3", r.GetOrAssign("aaa-1"))
}

// TestShortIDAssignmentIsOrderIndependent: shuffling the input slice must not
// change any assignment.
func TestShortIDAssignmentIsOrderIndependent(t *testing.T) {
	mk := func() []*SafeCard {
		return []*SafeCard{
			{ID: "c1", Name: "Forest"},
			{ID: "c2", Name: "Birds of Paradise"},
			{ID: "c3", Name: "Forest"},
			{ID: "c4", Name: "Wrath of God"},
		}
	}
	a := NewShortIDRegistry()
	a.AssignAll(mk())

	shuffled := mk()
	shuffled[0], shuffled[3] = shuffled[3], shuffled[0]
	shuffled[1], shuffled[2] = shuffled[2], shuffled[1]
	b := NewShortIDRegistry()
	b.AssignAll(shuffled)

	require.Equal(t, a.DumpAssignments(), b.DumpAssignments())
}

func TestShortIDStableAcrossZoneChanges(t *testing.T) {
	c := &SafeCard{ID: "h-1", Name: "Grizzly Bears", Zone: "HAND"}
	r := NewShortIDRegistry()
	r.AssignAll([]*SafeCard{c})
	before := r.GetOrAssign(c.ID)

	c.Zone = "BATTLEFIELD"
	c.Tapped = true
	r.AssignAll([]*SafeCard{c, {ID: "h-2", Name: "Aardvark"}})

	require.Equal(t, before, r.GetOrAssign(c.ID), "short ID must survive a zone change")
}

func TestShortIDSupersededIDsRemainResolvable(t *testing.T) {
	r := NewShortIDRegistry()
	local := r.GetOrAssign("card-a") // p1
	r.Register("card-a", "p7")       // server overrides

	require.Equal(t, "p7", r.GetOrAssign("card-a"))

	// The old ID must still resolve -- an LLM's mana plan from three decisions
	// ago may still reference it.
	got, err := r.Resolve(local)
	require.NoError(t, err)
	require.Equal(t, "card-a", got)

	got, err = r.Resolve("p7")
	require.NoError(t, err)
	require.Equal(t, "card-a", got)

	require.Equal(t, 8, r.PeekNextID(), "register must advance the counter past the registered ID")
	require.Equal(t, []string{"p1", "p7"}, r.SnapshotShortIDs())
}

func TestShortIDResolveUnknown(t *testing.T) {
	r := NewShortIDRegistry()
	_, err := r.Resolve("p99")
	require.Error(t, err)
	_, ok := r.TryResolve("p99")
	require.False(t, ok)
}

func TestShortIDClear(t *testing.T) {
	r := NewShortIDRegistry()
	r.GetOrAssign("a")
	r.Clear()
	require.Equal(t, 1, r.PeekNextID())
	require.Equal(t, "p1", r.GetOrAssign("b"))
}

func TestParseSequence(t *testing.T) {
	n, err := ParseSequence("p6")
	require.NoError(t, err)
	require.Equal(t, 6, n)

	_, err = ParseSequence("p")
	require.Error(t, err)
	_, err = ParseSequence("pX")
	require.Error(t, err)
}

func TestShortIDConcurrentAccess(t *testing.T) {
	r := NewShortIDRegistry()
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.GetOrAssign("shared")
			r.Sequence("shared")
			r.SnapshotShortIDs()
		}()
	}
	wg.Wait()
	require.Equal(t, "p1", r.GetOrAssign("shared"))
}
