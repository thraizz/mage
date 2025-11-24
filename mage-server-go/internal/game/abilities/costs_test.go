package abilities

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCostGameContext is a test implementation of GameContext for cost tests
type mockCostGameContext struct {
	permanents map[uuid.UUID]*mockCostPermanent
	players    map[uuid.UUID]*mockCostPlayer
}

type mockCostPermanent struct {
	id           uuid.UUID
	name         string
	cardType     string
	subTypes     []string
	controllerID uuid.UUID
	sacrificed   bool
}

type mockCostPlayer struct {
	id   uuid.UUID
	hand []*mockCostCard
}

type mockCostCard struct {
	id       uuid.UUID
	name     string
	cardType string
}

func newMockCostGameContext() *mockCostGameContext {
	return &mockCostGameContext{
		permanents: make(map[uuid.UUID]*mockCostPermanent),
		players:    make(map[uuid.UUID]*mockCostPlayer),
	}
}

func (m *mockCostGameContext) addPermanent(id, controllerID uuid.UUID, name, cardType string, subTypes []string) {
	m.permanents[id] = &mockCostPermanent{
		id:           id,
		name:         name,
		cardType:     cardType,
		subTypes:     subTypes,
		controllerID: controllerID,
	}
}

func (m *mockCostGameContext) addPlayer(id uuid.UUID) {
	m.players[id] = &mockCostPlayer{
		id:   id,
		hand: make([]*mockCostCard, 0),
	}
}

func (m *mockCostGameContext) addCardToHand(playerID, cardID uuid.UUID, name, cardType string) {
	player := m.players[playerID]
	player.hand = append(player.hand, &mockCostCard{
		id:       cardID,
		name:     name,
		cardType: cardType,
	})
}

// Implement GameContext interface
func (m *mockCostGameContext) GetCard(id uuid.UUID) (interface{}, error) {
	return nil, nil
}

func (m *mockCostGameContext) GetPlayer(id uuid.UUID) (interface{}, error) {
	return nil, nil
}

func (m *mockCostGameContext) DealDamage(sourceID, targetID uuid.UUID, amount int) error {
	return nil
}

func (m *mockCostGameContext) DrawCards(playerID uuid.UUID, amount int) error {
	return nil
}

func (m *mockCostGameContext) DestroyPermanent(permanentID uuid.UUID) error {
	return nil
}

func (m *mockCostGameContext) AddMana(playerID uuid.UUID, mana *Mana) error {
	return nil
}

func (m *mockCostGameContext) GetManaPool(playerID uuid.UUID) ManaPoolInterface {
	return nil
}

func (m *mockCostGameContext) TapPermanent(permanentID uuid.UUID) error {
	return nil
}

func (m *mockCostGameContext) UntapPermanent(permanentID uuid.UUID) error {
	return nil
}

func (m *mockCostGameContext) SacrificePermanent(permanentID uuid.UUID) error {
	perm, ok := m.permanents[permanentID]
	if !ok {
		return nil
	}
	perm.sacrificed = true
	delete(m.permanents, permanentID)
	return nil
}

func (m *mockCostGameContext) DiscardCard(playerID uuid.UUID, cardID uuid.UUID) error {
	player, ok := m.players[playerID]
	if !ok {
		return nil
	}

	for i, card := range player.hand {
		if card.id == cardID {
			player.hand = append(player.hand[:i], player.hand[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockCostGameContext) GetPlayerHand(playerID uuid.UUID) ([]interface{}, error) {
	player, ok := m.players[playerID]
	if !ok {
		return nil, nil
	}

	hand := make([]interface{}, len(player.hand))
	for i, card := range player.hand {
		hand[i] = &internalCard{
			ID:   card.id.String(),
			Name: card.name,
			Type: card.cardType,
		}
	}
	return hand, nil
}

func (m *mockCostGameContext) GetPermanentsControlledByPlayer(playerID uuid.UUID) ([]interface{}, error) {
	var result []interface{}
	for _, perm := range m.permanents {
		if perm.controllerID == playerID {
			result = append(result, &internalCard{
				ID:       perm.id.String(),
				Name:     perm.name,
				Type:     perm.cardType,
				SubTypes: perm.subTypes,
			})
		}
	}
	return result, nil
}

// ========================================
// SacrificeCost Tests
// ========================================

func TestSacrificeCost_CanPay_WithSufficientPermanents(t *testing.T) {
	ctx := context.Background()
	game := newMockCostGameContext()
	playerID := uuid.New()

	game.addPlayer(playerID)
	game.addPermanent(uuid.New(), playerID, "Grizzly Bears", "CREATURE", []string{"Bear"})
	game.addPermanent(uuid.New(), playerID, "Bronze Sword", "ARTIFACT", []string{"Equipment"})

	cost := NewSacrificeCost(1, "creature")

	canPay := cost.CanPay(ctx, game, playerID)
	assert.True(t, canPay, "Should be able to pay with 1 creature")
}

func TestSacrificeCost_CanPay_InsufficientPermanents(t *testing.T) {
	ctx := context.Background()
	game := newMockCostGameContext()
	playerID := uuid.New()

	game.addPlayer(playerID)
	game.addPermanent(uuid.New(), playerID, "Bronze Sword", "ARTIFACT", []string{"Equipment"})

	cost := NewSacrificeCost(1, "creature")

	canPay := cost.CanPay(ctx, game, playerID)
	assert.False(t, canPay, "Should not be able to pay without creatures")
}

func TestSacrificeCost_Pay_RemovesPermanent(t *testing.T) {
	ctx := context.Background()
	game := newMockCostGameContext()
	playerID := uuid.New()
	creatureID := uuid.New()

	game.addPlayer(playerID)
	game.addPermanent(creatureID, playerID, "Grizzly Bears", "CREATURE", []string{"Bear"})

	cost := NewSacrificeCost(1, "creature")

	err := cost.Pay(ctx, game, playerID)
	require.NoError(t, err)

	_, exists := game.permanents[creatureID]
	assert.False(t, exists, "Creature should be removed from battlefield")
}

func TestSacrificeCost_String(t *testing.T) {
	tests := []struct {
		name     string
		cost     *SacrificeCost
		expected string
	}{
		{
			name:     "Single creature",
			cost:     NewSacrificeCost(1, "creature"),
			expected: "Sacrifice a creature",
		},
		{
			name:     "Multiple artifacts",
			cost:     NewSacrificeCost(2, "artifact"),
			expected: "Sacrifice 2 artifacts",
		},
		{
			name:     "Source",
			cost:     NewSacrificeSourceCost(),
			expected: "Sacrifice this permanent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cost.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ========================================
// DiscardCost Tests
// ========================================

func TestDiscardCost_CanPay_WithSufficientCards(t *testing.T) {
	ctx := context.Background()
	game := newMockCostGameContext()
	playerID := uuid.New()

	game.addPlayer(playerID)
	game.addCardToHand(playerID, uuid.New(), "Lightning Bolt", "INSTANT")
	game.addCardToHand(playerID, uuid.New(), "Giant Growth", "INSTANT")

	cost := NewDiscardCost(1)

	canPay := cost.CanPay(ctx, game, playerID)
	assert.True(t, canPay, "Should be able to pay with 2 cards in hand")
}

func TestDiscardCost_CanPay_InsufficientCards(t *testing.T) {
	ctx := context.Background()
	game := newMockCostGameContext()
	playerID := uuid.New()

	game.addPlayer(playerID)
	game.addCardToHand(playerID, uuid.New(), "Lightning Bolt", "INSTANT")

	cost := NewDiscardCost(2)

	canPay := cost.CanPay(ctx, game, playerID)
	assert.False(t, canPay, "Should not be able to pay with only 1 card")
}

func TestDiscardCost_Pay_RemovesCards(t *testing.T) {
	ctx := context.Background()
	game := newMockCostGameContext()
	playerID := uuid.New()

	game.addPlayer(playerID)
	game.addCardToHand(playerID, uuid.New(), "Lightning Bolt", "INSTANT")
	game.addCardToHand(playerID, uuid.New(), "Giant Growth", "INSTANT")

	initialHandSize := len(game.players[playerID].hand)

	cost := NewDiscardCost(1)
	err := cost.Pay(ctx, game, playerID)
	require.NoError(t, err)

	finalHandSize := len(game.players[playerID].hand)
	assert.Equal(t, initialHandSize-1, finalHandSize, "Hand size should decrease by 1")
}

func TestDiscardCost_String(t *testing.T) {
	tests := []struct {
		name     string
		cost     *DiscardCost
		expected string
	}{
		{
			name:     "Single card",
			cost:     NewDiscardCost(1),
			expected: "Discard a card",
		},
		{
			name:     "Multiple cards",
			cost:     NewDiscardCost(3),
			expected: "Discard 3 cards",
		},
		{
			name:     "Random discard",
			cost:     NewDiscardCostRandom(1),
			expected: "Discard a card at random",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cost.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ========================================
// DiscardTargetCost Tests
// ========================================

func TestDiscardTargetCost_CanPay_WithMatchingCards(t *testing.T) {
	ctx := context.Background()
	game := newMockCostGameContext()
	playerID := uuid.New()

	game.addPlayer(playerID)
	game.addCardToHand(playerID, uuid.New(), "Bronze Sword", "ARTIFACT")
	game.addCardToHand(playerID, uuid.New(), "Lightning Bolt", "INSTANT")

	filter := NewArtifactCardFilter()
	cost := NewDiscardTargetCost(1, filter)

	canPay := cost.CanPay(ctx, game, playerID)
	assert.True(t, canPay, "Should be able to pay with 1 artifact card")
}

func TestDiscardTargetCost_CanPay_WithoutMatchingCards(t *testing.T) {
	ctx := context.Background()
	game := newMockCostGameContext()
	playerID := uuid.New()

	game.addPlayer(playerID)
	game.addCardToHand(playerID, uuid.New(), "Lightning Bolt", "INSTANT")
	game.addCardToHand(playerID, uuid.New(), "Giant Growth", "INSTANT")

	filter := NewArtifactCardFilter()
	cost := NewDiscardTargetCost(1, filter)

	canPay := cost.CanPay(ctx, game, playerID)
	assert.False(t, canPay, "Should not be able to pay without artifact cards")
}

func TestDiscardTargetCost_Pay_RemovesMatchingCards(t *testing.T) {
	ctx := context.Background()
	game := newMockCostGameContext()
	playerID := uuid.New()

	game.addPlayer(playerID)
	game.addCardToHand(playerID, uuid.New(), "Bronze Sword", "ARTIFACT")
	game.addCardToHand(playerID, uuid.New(), "Lightning Bolt", "INSTANT")

	initialHandSize := len(game.players[playerID].hand)

	filter := NewArtifactCardFilter()
	cost := NewDiscardTargetCost(1, filter)

	err := cost.Pay(ctx, game, playerID)
	require.NoError(t, err)

	finalHandSize := len(game.players[playerID].hand)
	assert.Equal(t, initialHandSize-1, finalHandSize, "Hand size should decrease by 1")
}

func TestDiscardTargetCost_String(t *testing.T) {
	filter := NewArtifactCardFilter()
	cost := NewDiscardTargetCost(2, filter)

	result := cost.String()
	assert.Contains(t, result, "Discard", "Should mention discarding")
	assert.Contains(t, result, "2", "Should mention the amount")
}

// ========================================
// Helper Function Tests
// ========================================

func TestPermanentMatchesFilter(t *testing.T) {
	tests := []struct {
		name     string
		card     *internalCard
		filter   string
		expected bool
	}{
		{
			name:     "Creature matches creature filter",
			card:     &internalCard{Type: "CREATURE"},
			filter:   "creature",
			expected: true,
		},
		{
			name:     "Artifact matches artifact filter",
			card:     &internalCard{Type: "ARTIFACT"},
			filter:   "artifact",
			expected: true,
		},
		{
			name:     "Creature doesn't match artifact filter",
			card:     &internalCard{Type: "CREATURE"},
			filter:   "artifact",
			expected: false,
		},
		{
			name:     "Empty filter matches anything",
			card:     &internalCard{Type: "LAND"},
			filter:   "",
			expected: true,
		},
		{
			name:     "Subtype matches",
			card:     &internalCard{Type: "CREATURE", SubTypes: []string{"Bear"}},
			filter:   "bear",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := permanentMatchesFilter(tt.card, tt.filter)
			assert.Equal(t, tt.expected, result)
		})
	}
}
