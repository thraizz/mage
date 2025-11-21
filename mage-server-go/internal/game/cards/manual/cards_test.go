package manual_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/cards"
	_ "github.com/magefree/mage-server-go/internal/game/cards/manual" // Import to register cards
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicLands(t *testing.T) {
	ownerID := uuid.New()

	t.Run("Create Plains", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Plains")
		require.True(t, ok, "Plains should be registered")
		require.NotNil(t, builder)

		plains, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Plains", plains.Name)
		assert.Equal(t, []string{"LAND"}, plains.Types)
		assert.Equal(t, []string{"PLAINS"}, plains.Subtypes)
		assert.Equal(t, []string{"BASIC"}, plains.Supertypes)
		assert.Len(t, plains.Abilities, 1, "Plains should have 1 mana ability")
	})

	t.Run("Create Island", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Island")
		require.True(t, ok)

		island, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Island", island.Name)
		assert.Equal(t, []string{"LAND"}, island.Types)
		assert.Len(t, island.Abilities, 1)
	})

	t.Run("Create Swamp", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Swamp")
		require.True(t, ok)

		swamp, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Swamp", swamp.Name)
		assert.Len(t, swamp.Abilities, 1)
	})

	t.Run("Create Mountain", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Mountain")
		require.True(t, ok)

		mountain, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Mountain", mountain.Name)
		assert.Len(t, mountain.Abilities, 1)
	})

	t.Run("Create Forest", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Forest")
		require.True(t, ok)

		forest, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Forest", forest.Name)
		assert.Len(t, forest.Abilities, 1)
	})
}

func TestVanillaCreatures(t *testing.T) {
	ownerID := uuid.New()

	t.Run("Create Grizzly Bears", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Grizzly Bears")
		require.True(t, ok)

		bear, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Grizzly Bears", bear.Name)
		assert.Equal(t, "{1}{G}", bear.ManaCost)
		assert.Equal(t, []string{"CREATURE"}, bear.Types)
		assert.Equal(t, []string{"BEAR"}, bear.Subtypes)
		assert.Equal(t, "2", bear.Power)
		assert.Equal(t, "2", bear.Toughness)
		assert.Len(t, bear.Abilities, 0, "Vanilla creature has no abilities")
	})

	t.Run("Create Savannah Lions", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Savannah Lions")
		require.True(t, ok)

		lion, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Savannah Lions", lion.Name)
		assert.Equal(t, "{W}", lion.ManaCost)
		assert.Equal(t, "2", lion.Power)
		assert.Equal(t, "1", lion.Toughness)
	})

	t.Run("Create Hill Giant", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Hill Giant")
		require.True(t, ok)

		giant, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Hill Giant", giant.Name)
		assert.Equal(t, "{3}{R}", giant.ManaCost)
		assert.Equal(t, "3", giant.Power)
		assert.Equal(t, "3", giant.Toughness)
	})
}

func TestSimpleSpells(t *testing.T) {
	ownerID := uuid.New()

	t.Run("Create Lightning Bolt", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Lightning Bolt")
		require.True(t, ok)

		bolt, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Lightning Bolt", bolt.Name)
		assert.Equal(t, "{R}", bolt.ManaCost)
		assert.Equal(t, []string{"INSTANT"}, bolt.Types)
		assert.Len(t, bolt.Abilities, 1, "Lightning Bolt has 1 spell ability")
	})

	t.Run("Create Shock", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Shock")
		require.True(t, ok)

		shock, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Shock", shock.Name)
		assert.Len(t, shock.Abilities, 1)
	})

	t.Run("Create Giant Growth", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Giant Growth")
		require.True(t, ok)

		growth, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Giant Growth", growth.Name)
		assert.Equal(t, "{G}", growth.ManaCost)
		assert.Len(t, growth.Abilities, 1)
	})

	t.Run("Create Divination", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Divination")
		require.True(t, ok)

		div, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Divination", div.Name)
		assert.Equal(t, "{2}{U}", div.ManaCost)
		assert.Equal(t, []string{"SORCERY"}, div.Types)
		assert.Len(t, div.Abilities, 1)
	})

	t.Run("Create Murder", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Murder")
		require.True(t, ok)

		murder, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Murder", murder.Name)
		assert.Equal(t, "{1}{B}{B}", murder.ManaCost)
		assert.Len(t, murder.Abilities, 1)
	})

	t.Run("Create Counterspell", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Counterspell")
		require.True(t, ok)

		counter, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Counterspell", counter.Name)
		assert.Equal(t, "{U}{U}", counter.ManaCost)
		assert.Len(t, counter.Abilities, 1)
	})
}

func TestActivatedAbilities(t *testing.T) {
	ownerID := uuid.New()

	t.Run("Create Llanowar Elves", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Llanowar Elves")
		require.True(t, ok)

		elves, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Llanowar Elves", elves.Name)
		assert.Equal(t, "{G}", elves.ManaCost)
		assert.Equal(t, "1", elves.Power)
		assert.Equal(t, "1", elves.Toughness)
		assert.Len(t, elves.Abilities, 1, "Llanowar Elves has 1 mana ability")
	})

	t.Run("Create Prodigal Pyromancer", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Prodigal Pyromancer")
		require.True(t, ok)

		pyro, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Prodigal Pyromancer", pyro.Name)
		assert.Equal(t, "{2}{R}", pyro.ManaCost)
		assert.Equal(t, "1", pyro.Power)
		assert.Equal(t, "1", pyro.Toughness)
		assert.Len(t, pyro.Abilities, 1, "Prodigal Pyromancer has 1 activated ability")
	})

	t.Run("Create Prodigal Sorcerer", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Prodigal Sorcerer")
		require.True(t, ok)

		tim, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Prodigal Sorcerer", tim.Name)
		assert.Equal(t, "{2}{U}", tim.ManaCost)
		assert.Equal(t, "1", tim.Power)
		assert.Equal(t, "1", tim.Toughness)
		assert.Len(t, tim.Abilities, 1, "Prodigal Sorcerer has 1 activated ability")
	})
}

func TestKeywordCreatures(t *testing.T) {
	ownerID := uuid.New()

	t.Run("Create Serra Angel", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Serra Angel")
		require.True(t, ok)

		angel, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Serra Angel", angel.Name)
		assert.Equal(t, "{3}{W}{W}", angel.ManaCost)
		assert.Equal(t, "4", angel.Power)
		assert.Equal(t, "4", angel.Toughness)
		assert.Len(t, angel.Abilities, 2, "Serra Angel has 2 keyword abilities (flying, vigilance)")
	})

	t.Run("Create Serra's Guardian", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Serra's Guardian")
		require.True(t, ok)

		guardian, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Serra's Guardian", guardian.Name)
		assert.Equal(t, "{4}{W}", guardian.ManaCost)
		assert.Equal(t, "3", guardian.Power)
		assert.Equal(t, "4", guardian.Toughness)
		assert.Len(t, guardian.Abilities, 2, "Serra's Guardian has 2 keyword abilities")
	})

	t.Run("Create Vampire Nighthawk", func(t *testing.T) {
		builder, ok := cards.Registry.Get("Vampire Nighthawk")
		require.True(t, ok)

		nighthawk, err := builder(ownerID, nil)
		require.NoError(t, err)
		assert.Equal(t, "Vampire Nighthawk", nighthawk.Name)
		assert.Equal(t, "{1}{B}{B}", nighthawk.ManaCost)
		assert.Equal(t, "2", nighthawk.Power)
		assert.Equal(t, "3", nighthawk.Toughness)
		assert.Len(t, nighthawk.Abilities, 3, "Vampire Nighthawk has 3 keyword abilities")
	})
}

func TestCardRegistry(t *testing.T) {
	t.Run("All 20 cards are registered", func(t *testing.T) {
		expectedCards := []string{
			// Basic lands (5)
			"Plains", "Island", "Swamp", "Mountain", "Forest",
			// Vanilla creatures (3)
			"Grizzly Bears", "Savannah Lions", "Hill Giant",
			// Simple spells (6)
			"Lightning Bolt", "Shock", "Giant Growth",
			"Divination", "Murder", "Counterspell",
			// Activated abilities (3)
			"Llanowar Elves", "Prodigal Pyromancer", "Prodigal Sorcerer",
			// Keyword creatures (3)
			"Serra Angel", "Serra's Guardian", "Vampire Nighthawk",
		}

		for _, cardName := range expectedCards {
			t.Run(cardName, func(t *testing.T) {
				_, ok := cards.Registry.Get(cardName)
				assert.True(t, ok, "%s should be registered", cardName)
			})
		}

		// Verify count
		assert.GreaterOrEqual(t, cards.Registry.Count(), 20,
			"At least 20 cards should be registered")
	})
}
