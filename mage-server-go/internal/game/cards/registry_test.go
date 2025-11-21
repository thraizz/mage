package cards_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	// Clear registry before tests
	cards.Registry.Clear()

	t.Run("Register and Get card", func(t *testing.T) {
		// Register a test card
		testBuilder := func(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
			return &game.Card{ID: uuid.New(), Name: "Test Card"}, nil
		}

		cards.Register("Test Card", testBuilder)

		// Verify it's registered
		builder, ok := cards.Registry.Get("Test Card")
		require.True(t, ok, "Card should be registered")
		require.NotNil(t, builder, "Builder should not be nil")

		// Create a card using the builder
		card, err := builder(uuid.New(), &cards.CardInfo{Name: "Test Card"})
		require.NoError(t, err)
		assert.Equal(t, "Test Card", card.Name)
	})

	t.Run("IsImplemented", func(t *testing.T) {
		cards.Registry.Clear()

		// Not implemented initially
		assert.False(t, cards.Registry.IsImplemented("Lightning Bolt"))

		// Register it
		cards.Register("Lightning Bolt", func(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
			return &game.Card{ID: uuid.New(), Name: "Lightning Bolt"}, nil
		})

		// Now it's implemented
		assert.True(t, cards.Registry.IsImplemented("Lightning Bolt"))
	})

	t.Run("Count", func(t *testing.T) {
		cards.Registry.Clear()

		assert.Equal(t, 0, cards.Registry.Count())

		// Register some cards
		cards.Register("Card 1", dummyBuilder)
		assert.Equal(t, 1, cards.Registry.Count())

		cards.Register("Card 2", dummyBuilder)
		assert.Equal(t, 2, cards.Registry.Count())

		cards.Register("Card 3", dummyBuilder)
		assert.Equal(t, 3, cards.Registry.Count())
	})

	t.Run("ListAll", func(t *testing.T) {
		cards.Registry.Clear()

		cards.Register("Card A", dummyBuilder)
		cards.Register("Card B", dummyBuilder)
		cards.Register("Card C", dummyBuilder)

		list := cards.Registry.ListAll()
		assert.Len(t, list, 3)
		assert.Contains(t, list, "Card A")
		assert.Contains(t, list, "Card B")
		assert.Contains(t, list, "Card C")
	})

	t.Run("Get non-existent card", func(t *testing.T) {
		cards.Registry.Clear()

		_, ok := cards.Registry.Get("Non Existent Card")
		assert.False(t, ok, "Should return false for non-existent card")
	})

	t.Run("Re-registration prints warning", func(t *testing.T) {
		cards.Registry.Clear()

		// Register once
		cards.Register("Duplicate Card", dummyBuilder)

		// Register again (should print warning but not fail)
		cards.Register("Duplicate Card", dummyBuilder)

		// Should still be registered
		assert.True(t, cards.Registry.IsImplemented("Duplicate Card"))
	})
}

// dummyBuilder is a simple builder for testing
func dummyBuilder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	return &game.Card{
		ID:   uuid.New(),
		Name: info.Name,
	}, nil
}
