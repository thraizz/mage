package abilities

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================================
// Mock Game Context for Testing
// ========================================

type mockCounterGameContext struct {
	permanents map[uuid.UUID]*mockPermanent
	players    map[uuid.UUID]*mockPlayer
	cards      map[uuid.UUID]*mockCard
	messages   []string
}

type mockPermanent struct {
	id       uuid.UUID
	counters *counters.Counters
}

type mockPlayer struct {
	id       uuid.UUID
	counters *counters.Counters
}

type mockCard struct {
	id       uuid.UUID
	counters *counters.Counters
}

func newMockCounterGameContext() *mockCounterGameContext {
	return &mockCounterGameContext{
		permanents: make(map[uuid.UUID]*mockPermanent),
		players:    make(map[uuid.UUID]*mockPlayer),
		cards:      make(map[uuid.UUID]*mockCard),
		messages:   make([]string, 0),
	}
}

func (m *mockCounterGameContext) GetCard(id uuid.UUID) (interface{}, error) {
	if card, ok := m.cards[id]; ok {
		return card, nil
	}
	return nil, fmt.Errorf("card not found")
}

func (m *mockCounterGameContext) GetPlayer(id uuid.UUID) (interface{}, error) {
	if player, ok := m.players[id]; ok {
		return player, nil
	}
	return nil, fmt.Errorf("player not found")
}

func (m *mockCounterGameContext) GetPermanent(id uuid.UUID) (interface{}, error) {
	if perm, ok := m.permanents[id]; ok {
		return perm, nil
	}
	return nil, fmt.Errorf("permanent not found")
}

func (m *mockCounterGameContext) AddCountersToPermanent(permanent interface{}, counter *counters.Counter) error {
	perm, ok := permanent.(*mockPermanent)
	if !ok {
		return fmt.Errorf("invalid permanent type")
	}
	perm.counters.AddCounter(counter)
	return nil
}

func (m *mockCounterGameContext) AddCountersToPlayer(player interface{}, counter *counters.Counter) error {
	p, ok := player.(*mockPlayer)
	if !ok {
		return fmt.Errorf("invalid player type")
	}
	p.counters.AddCounter(counter)
	return nil
}

func (m *mockCounterGameContext) AddCountersToCard(card interface{}, counter *counters.Counter) error {
	c, ok := card.(*mockCard)
	if !ok {
		return fmt.Errorf("invalid card type")
	}
	c.counters.AddCounter(counter)
	return nil
}

func (m *mockCounterGameContext) GetAllPermanents() ([]interface{}, error) {
	result := make([]interface{}, 0, len(m.permanents))
	for _, perm := range m.permanents {
		result = append(result, perm)
	}
	return result, nil
}

func (m *mockCounterGameContext) InformPlayers(message string) {
	m.messages = append(m.messages, message)
}

// Unused GameContext methods
func (m *mockCounterGameContext) DealDamage(sourceID, targetID uuid.UUID, amount int) error {
	return nil
}

func (m *mockCounterGameContext) DrawCards(playerID uuid.UUID, amount int) error {
	return nil
}

func (m *mockCounterGameContext) DestroyPermanent(permanentID uuid.UUID) error {
	return nil
}

func (m *mockCounterGameContext) AddMana(playerID uuid.UUID, mana *Mana) error {
	return nil
}

func (m *mockCounterGameContext) TapPermanent(permanentID uuid.UUID) error {
	return nil
}

func (m *mockCounterGameContext) UntapPermanent(permanentID uuid.UUID) error {
	return nil
}

// ========================================
// Tests for AddCountersSourceEffect
// ========================================

func TestAddCountersSourceEffect_Apply(t *testing.T) {
	tests := []struct {
		name          string
		counter       *counters.Counter
		informPlayers bool
		expectedCount int
		expectedName  string
		expectMessage bool
	}{
		{
			name:          "Add single +1/+1 counter",
			counter:       counters.NewBoostCounter(1, 1, 1).Counter,
			informPlayers: true,
			expectedCount: 1,
			expectedName:  "+1/+1",
			expectMessage: true,
		},
		{
			name:          "Add multiple +1/+1 counters",
			counter:       counters.NewBoostCounter(1, 1, 3).Counter,
			informPlayers: true,
			expectedCount: 3,
			expectedName:  "+1/+1",
			expectMessage: true,
		},
		{
			name:          "Add lore counter",
			counter:       counters.CounterTypeLore.CreateInstance(1),
			informPlayers: false,
			expectedCount: 1,
			expectedName:  "lore",
			expectMessage: false,
		},
		{
			name:          "Add shield counter",
			counter:       counters.CounterTypeShield.CreateInstance(1),
			informPlayers: true,
			expectedCount: 1,
			expectedName:  "shield",
			expectMessage: true,
		},
		{
			name:          "Add divinity counters",
			counter:       counters.CounterTypeDivinity.CreateInstance(2),
			informPlayers: true,
			expectedCount: 2,
			expectedName:  "divinity",
			expectMessage: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			game := newMockCounterGameContext()

			// Create a source permanent
			sourceID := uuid.New()
			sourcePerm := &mockPermanent{
				id:       sourceID,
				counters: counters.NewCounters(),
			}
			game.permanents[sourceID] = sourcePerm

			// Create the effect
			effect := NewAddCountersSourceEffectInform(tt.counter, tt.informPlayers)

			// Apply the effect
			err := effect.Apply(ctx, game, sourceID, nil)
			require.NoError(t, err)

			// Verify counters were added
			assert.Equal(t, tt.expectedCount, sourcePerm.counters.GetCount(tt.expectedName))

			// Verify message if expected
			if tt.expectMessage {
				assert.NotEmpty(t, game.messages)
			}
		})
	}
}

func TestAddCountersSourceEffect_GetDescription(t *testing.T) {
	tests := []struct {
		name         string
		counter      *counters.Counter
		expectedDesc string
	}{
		{
			name:         "Single counter",
			counter:      counters.CounterTypeP1P1.CreateInstance(1),
			expectedDesc: "put a +1/+1 counter on {this}",
		},
		{
			name:         "Multiple counters",
			counter:      counters.CounterTypeP1P1.CreateInstance(3),
			expectedDesc: "put 3 +1/+1 counters on {this}",
		},
		{
			name:         "Lore counter",
			counter:      counters.CounterTypeLore.CreateInstance(1),
			expectedDesc: "put a lore counter on {this}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			effect := NewAddCountersSourceEffect(tt.counter)
			desc := effect.GetDescription()
			assert.Equal(t, tt.expectedDesc, desc)
		})
	}
}

// ========================================
// Tests for AddCountersTargetEffect
// ========================================

func TestAddCountersTargetEffect_ApplyToPermanent(t *testing.T) {
	ctx := context.Background()
	game := newMockCounterGameContext()

	// Create target permanents
	target1 := uuid.New()
	target2 := uuid.New()
	perm1 := &mockPermanent{id: target1, counters: counters.NewCounters()}
	perm2 := &mockPermanent{id: target2, counters: counters.NewCounters()}
	game.permanents[target1] = perm1
	game.permanents[target2] = perm2

	// Create effect
	counter := counters.CounterTypeP1P1.CreateInstance(2)
	effect := NewAddCountersTargetEffect(counter)

	// Apply to targets
	err := effect.Apply(ctx, game, uuid.New(), []uuid.UUID{target1, target2})
	require.NoError(t, err)

	// Verify both permanents got counters
	assert.Equal(t, 2, perm1.counters.GetCount("+1/+1"))
	assert.Equal(t, 2, perm2.counters.GetCount("+1/+1"))
}

func TestAddCountersTargetEffect_ApplyToPlayer(t *testing.T) {
	ctx := context.Background()
	game := newMockCounterGameContext()

	// Create target player
	targetID := uuid.New()
	player := &mockPlayer{id: targetID, counters: counters.NewCounters()}
	game.players[targetID] = player

	// Create effect (poison counter)
	counter := counters.CounterTypePoison.CreateInstance(3)
	effect := NewAddCountersTargetEffect(counter)

	// Apply to player
	err := effect.Apply(ctx, game, uuid.New(), []uuid.UUID{targetID})
	require.NoError(t, err)

	// Verify player got counters
	assert.Equal(t, 3, player.counters.GetCount("poison"))
}

func TestAddCountersTargetEffect_ApplyToCard(t *testing.T) {
	ctx := context.Background()
	game := newMockCounterGameContext()

	// Create target card (e.g., in graveyard)
	targetID := uuid.New()
	card := &mockCard{id: targetID, counters: counters.NewCounters()}
	game.cards[targetID] = card

	// Create effect
	counter := counters.CounterTypeAge.CreateInstance(1)
	effect := NewAddCountersTargetEffect(counter)

	// Apply to card
	err := effect.Apply(ctx, game, uuid.New(), []uuid.UUID{targetID})
	require.NoError(t, err)

	// Verify card got counters
	assert.Equal(t, 1, card.counters.GetCount("age"))
}

func TestAddCountersTargetEffect_GetDescription(t *testing.T) {
	tests := []struct {
		name         string
		counter      *counters.Counter
		expectedDesc string
	}{
		{
			name:         "Single counter",
			counter:      counters.CounterTypeP1P1.CreateInstance(1),
			expectedDesc: "put a +1/+1 counter on target",
		},
		{
			name:         "Multiple counters",
			counter:      counters.CounterTypeShield.CreateInstance(2),
			expectedDesc: "put 2 shield counters on target",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			effect := NewAddCountersTargetEffect(tt.counter)
			desc := effect.GetDescription()
			assert.Equal(t, tt.expectedDesc, desc)
		})
	}
}

// ========================================
// Tests for AddCountersAllEffect
// ========================================

func TestAddCountersAllEffect_Apply(t *testing.T) {
	ctx := context.Background()
	game := newMockCounterGameContext()

	// Create multiple permanents
	perm1 := &mockPermanent{id: uuid.New(), counters: counters.NewCounters()}
	perm2 := &mockPermanent{id: uuid.New(), counters: counters.NewCounters()}
	perm3 := &mockPermanent{id: uuid.New(), counters: counters.NewCounters()}
	game.permanents[perm1.id] = perm1
	game.permanents[perm2.id] = perm2
	game.permanents[perm3.id] = perm3

	// Create effect (no filter, affects all)
	counter := counters.CounterTypeP1P1.CreateInstance(1)
	effect := NewAddCountersAllEffect(counter, nil, "each permanent")

	// Apply effect
	err := effect.Apply(ctx, game, uuid.New(), nil)
	require.NoError(t, err)

	// Verify all permanents got counters
	assert.Equal(t, 1, perm1.counters.GetCount("+1/+1"))
	assert.Equal(t, 1, perm2.counters.GetCount("+1/+1"))
	assert.Equal(t, 1, perm3.counters.GetCount("+1/+1"))
}

func TestAddCountersAllEffect_WithFilter(t *testing.T) {
	ctx := context.Background()
	game := newMockCounterGameContext()

	// Create multiple permanents
	perm1 := &mockPermanent{id: uuid.New(), counters: counters.NewCounters()}
	perm2 := &mockPermanent{id: uuid.New(), counters: counters.NewCounters()}
	perm3 := &mockPermanent{id: uuid.New(), counters: counters.NewCounters()}
	game.permanents[perm1.id] = perm1
	game.permanents[perm2.id] = perm2
	game.permanents[perm3.id] = perm3

	// Create filter that only accepts first two permanents
	acceptedIDs := map[uuid.UUID]bool{
		perm1.id: true,
		perm2.id: true,
	}
	filter := func(permanent interface{}) bool {
		p := permanent.(*mockPermanent)
		return acceptedIDs[p.id]
	}

	// Create effect with filter
	counter := counters.CounterTypeLore.CreateInstance(1)
	effect := NewAddCountersAllEffect(counter, filter, "each matching permanent")

	// Apply effect
	err := effect.Apply(ctx, game, uuid.New(), nil)
	require.NoError(t, err)

	// Verify only filtered permanents got counters
	assert.Equal(t, 1, perm1.counters.GetCount("lore"))
	assert.Equal(t, 1, perm2.counters.GetCount("lore"))
	assert.Equal(t, 0, perm3.counters.GetCount("lore"))
}

func TestAddCountersAllEffect_GetDescription(t *testing.T) {
	tests := []struct {
		name         string
		counter      *counters.Counter
		description  string
		expectedDesc string
	}{
		{
			name:         "All permanents, single counter",
			counter:      counters.CounterTypeP1P1.CreateInstance(1),
			description:  "each permanent",
			expectedDesc: "put a +1/+1 counter on each permanent",
		},
		{
			name:         "All creatures, multiple counters",
			counter:      counters.CounterTypeP1P1.CreateInstance(2),
			description:  "each creature",
			expectedDesc: "put 2 +1/+1 counters on each creature",
		},
		{
			name:         "Custom description",
			counter:      counters.CounterTypeDivinity.CreateInstance(3),
			description:  "each creature you control",
			expectedDesc: "put 3 divinity counters on each creature you control",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			effect := NewAddCountersAllEffect(tt.counter, nil, tt.description)
			desc := effect.GetDescription()
			assert.Equal(t, tt.expectedDesc, desc)
		})
	}
}

// ========================================
// Integration Tests
// ========================================

func TestCounterEffects_Integration(t *testing.T) {
	ctx := context.Background()
	game := newMockCounterGameContext()

	// Setup: Create a permanent and a player
	permID := uuid.New()
	perm := &mockPermanent{id: permID, counters: counters.NewCounters()}
	game.permanents[permID] = perm

	playerID := uuid.New()
	player := &mockPlayer{id: playerID, counters: counters.NewCounters()}
	game.players[playerID] = player

	// Test 1: Add +1/+1 counters to permanent
	t.Run("Add +1/+1 to permanent", func(t *testing.T) {
		effect := NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(2))
		err := effect.Apply(ctx, game, permID, nil)
		require.NoError(t, err)
		assert.Equal(t, 2, perm.counters.GetCount("+1/+1"))
	})

	// Test 2: Add more +1/+1 counters (should stack)
	t.Run("Add more +1/+1 counters", func(t *testing.T) {
		effect := NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1))
		err := effect.Apply(ctx, game, permID, nil)
		require.NoError(t, err)
		assert.Equal(t, 3, perm.counters.GetCount("+1/+1"))
	})

	// Test 3: Add poison counters to player
	t.Run("Add poison to player", func(t *testing.T) {
		effect := NewAddCountersTargetEffect(counters.CounterTypePoison.CreateInstance(5))
		err := effect.Apply(ctx, game, uuid.New(), []uuid.UUID{playerID})
		require.NoError(t, err)
		assert.Equal(t, 5, player.counters.GetCount("poison"))
	})

	// Test 4: Add shield counter to permanent
	t.Run("Add shield counter", func(t *testing.T) {
		effect := NewAddCountersTargetEffect(counters.CounterTypeShield.CreateInstance(1))
		err := effect.Apply(ctx, game, uuid.New(), []uuid.UUID{permID})
		require.NoError(t, err)
		assert.Equal(t, 1, perm.counters.GetCount("shield"))
	})

	// Test 5: Verify all counters exist
	t.Run("Verify all counters", func(t *testing.T) {
		assert.Equal(t, 3, perm.counters.GetCount("+1/+1"))
		assert.Equal(t, 1, perm.counters.GetCount("shield"))
		assert.Equal(t, 5, player.counters.GetCount("poison"))
	})
}
