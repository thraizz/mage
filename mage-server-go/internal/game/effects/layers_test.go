package effects

import "testing"

func TestLayerSystemAppliesEffects(t *testing.T) {
	system := NewLayerSystem()
	system.AddEffect(NewSimplePTBoostEffect("source", "Alice", 0, 1, false))

	snapshot := NewSnapshot("creature", "Alice", []string{"Creature"}, 2, 2, true, true)
	system.Apply(snapshot)

	if snapshot.Power != 2 || snapshot.Toughness != 3 {
		t.Fatalf("expected 2/3 after effect, got %d/%d", snapshot.Power, snapshot.Toughness)
	}
}
