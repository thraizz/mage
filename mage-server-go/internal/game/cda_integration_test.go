package game

import (
	"testing"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestCDAIntegration_TarmogoyfPowerToughness tests that Tarmogoyf's power and toughness
// are calculated correctly based on card types in all graveyards
func TestCDAIntegration_TarmogoyfPowerToughness(t *testing.T) {
	logger := zap.NewNop()
	engine := NewMageEngine(logger)

	player1ID := uuid.New().String()
	player2ID := uuid.New().String()

	// Start a game
	gameID := uuid.New().String()
	err := engine.StartGame(gameID, []string{player1ID, player2ID}, "TwoPlayerDuel")
	require.NoError(t, err)

	// Get game state
	engine.mu.RLock()
	gameState, exists := engine.games[gameID]
	engine.mu.RUnlock()
	require.True(t, exists, "Game should exist")

	// Create a Tarmogoyf card with CDA
	tarmogoyfID := uuid.New()
	tarmogoyfCDA := abilities.NewTarmogoyfCDA(tarmogoyfID)

	// Register the CDA with the ability registry
	gameState.mu.Lock()
	gameState.abilityRegistry.RegisterAbility(
		tarmogoyfCDA,
		uuid.MustParse(player1ID),
		0,
		abilities.ZoneBattlefield,
	)

	// Create the Tarmogoyf card
	tarmogoyf := &internalCard{
		ID:           tarmogoyfID.String(),
		Name:         "Tarmogoyf",
		ManaCost:     "{1}{G}",
		Type:         "Creature",
		SubTypes:     []string{"Lhurgoyf"},
		Color:        "G",
		Power:        "*",
		Toughness:    "1+*",
		Zone:         zoneBattlefield,
		ControllerID: player1ID,
		OwnerID:      player1ID,
		Abilities: []EngineAbilityView{
			{
				ID:   tarmogoyfCDA.GetID().String(),
				Text: "Tarmogoyf's power is equal to the number of card types among cards in all graveyards and its toughness is equal to that number plus 1.",
			},
		},
	}

	// Add Tarmogoyf to the battlefield
	gameState.cards[tarmogoyf.ID] = tarmogoyf
	gameState.battlefield = append(gameState.battlefield, tarmogoyf)
	gameState.mu.Unlock()

	// Test 1: With empty graveyards, Tarmogoyf should be 0/1
	t.Run("EmptyGraveyards", func(t *testing.T) {
		power, err := engine.getCreaturePower(gameState, tarmogoyf)
		assert.NoError(t, err)
		assert.Equal(t, 0, power, "Tarmogoyf should have 0 power with empty graveyards")

		toughness, err := engine.getCreatureToughness(gameState, tarmogoyf)
		assert.NoError(t, err)
		assert.Equal(t, 1, toughness, "Tarmogoyf should have 1 toughness with empty graveyards")
	})

	// Test 2: Add an instant to player1's graveyard (1 card type)
	t.Run("OneCardType_Instant", func(t *testing.T) {
		gameState.mu.Lock()
		instant := &internalCard{
			ID:           uuid.New().String(),
			Name:         "Lightning Bolt",
			Type:         "Instant",
			Zone:         zoneGraveyard,
			ControllerID: player1ID,
			OwnerID:      player1ID,
		}
		gameState.cards[instant.ID] = instant
		player1 := gameState.players[player1ID]
		player1.Graveyard = append(player1.Graveyard, instant)
		gameState.mu.Unlock()

		power, err := engine.getCreaturePower(gameState, tarmogoyf)
		assert.NoError(t, err)
		assert.Equal(t, 1, power, "Tarmogoyf should have 1 power with 1 card type (Instant)")

		toughness, err := engine.getCreatureToughness(gameState, tarmogoyf)
		assert.NoError(t, err)
		assert.Equal(t, 2, toughness, "Tarmogoyf should have 2 toughness with 1 card type")
	})

	// Test 3: Add a creature to player2's graveyard (2 card types total)
	t.Run("TwoCardTypes_InstantAndCreature", func(t *testing.T) {
		gameState.mu.Lock()
		creature := &internalCard{
			ID:           uuid.New().String(),
			Name:         "Grizzly Bears",
			Type:         "Creature",
			SubTypes:     []string{"Bear"},
			Zone:         zoneGraveyard,
			ControllerID: player2ID,
			OwnerID:      player2ID,
		}
		gameState.cards[creature.ID] = creature
		player2 := gameState.players[player2ID]
		player2.Graveyard = append(player2.Graveyard, creature)
		gameState.mu.Unlock()

		power, err := engine.getCreaturePower(gameState, tarmogoyf)
		assert.NoError(t, err)
		assert.Equal(t, 2, power, "Tarmogoyf should have 2 power with 2 card types (Instant, Creature)")

		toughness, err := engine.getCreatureToughness(gameState, tarmogoyf)
		assert.NoError(t, err)
		assert.Equal(t, 3, toughness, "Tarmogoyf should have 3 toughness with 2 card types")
	})

	// Test 4: Add a sorcery (3 card types total)
	t.Run("ThreeCardTypes", func(t *testing.T) {
		gameState.mu.Lock()
		sorcery := &internalCard{
			ID:           uuid.New().String(),
			Name:         "Wrath of God",
			Type:         "Sorcery",
			Zone:         zoneGraveyard,
			ControllerID: player1ID,
			OwnerID:      player1ID,
		}
		gameState.cards[sorcery.ID] = sorcery
		player1 := gameState.players[player1ID]
		player1.Graveyard = append(player1.Graveyard, sorcery)
		gameState.mu.Unlock()

		power, err := engine.getCreaturePower(gameState, tarmogoyf)
		assert.NoError(t, err)
		assert.Equal(t, 3, power, "Tarmogoyf should have 3 power with 3 card types")

		toughness, err := engine.getCreatureToughness(gameState, tarmogoyf)
		assert.NoError(t, err)
		assert.Equal(t, 4, toughness, "Tarmogoyf should have 4 toughness with 3 card types")
	})

	// Test 5: Add duplicate card type (should still be 3 types)
	t.Run("DuplicateCardType", func(t *testing.T) {
		gameState.mu.Lock()
		anotherCreature := &internalCard{
			ID:           uuid.New().String(),
			Name:         "Runeclaw Bear",
			Type:         "Creature",
			SubTypes:     []string{"Bear"},
			Zone:         zoneGraveyard,
			ControllerID: player1ID,
			OwnerID:      player1ID,
		}
		gameState.cards[anotherCreature.ID] = anotherCreature
		player1 := gameState.players[player1ID]
		player1.Graveyard = append(player1.Graveyard, anotherCreature)
		gameState.mu.Unlock()

		power, err := engine.getCreaturePower(gameState, tarmogoyf)
		assert.NoError(t, err)
		assert.Equal(t, 3, power, "Tarmogoyf should still have 3 power (duplicate Creature type doesn't count twice)")

		toughness, err := engine.getCreatureToughness(gameState, tarmogoyf)
		assert.NoError(t, err)
		assert.Equal(t, 4, toughness, "Tarmogoyf should still have 4 toughness")
	})

	// Test 6: Add artifact, enchantment, and land (6 card types total)
	t.Run("SixCardTypes", func(t *testing.T) {
		gameState.mu.Lock()

		// Add artifact
		artifact := &internalCard{
			ID:           uuid.New().String(),
			Name:         "Sol Ring",
			Type:         "Artifact",
			Zone:         zoneGraveyard,
			ControllerID: player2ID,
			OwnerID:      player2ID,
		}
		gameState.cards[artifact.ID] = artifact
		player2 := gameState.players[player2ID]
		player2.Graveyard = append(player2.Graveyard, artifact)

		// Add enchantment
		enchantment := &internalCard{
			ID:           uuid.New().String(),
			Name:         "Pacifism",
			Type:         "Enchantment",
			Zone:         zoneGraveyard,
			ControllerID: player1ID,
			OwnerID:      player1ID,
		}
		gameState.cards[enchantment.ID] = enchantment
		player1 := gameState.players[player1ID]
		player1.Graveyard = append(player1.Graveyard, enchantment)

		// Add land
		land := &internalCard{
			ID:           uuid.New().String(),
			Name:         "Forest",
			Type:         "Land",
			SubTypes:     []string{"Forest"},
			Zone:         zoneGraveyard,
			ControllerID: player1ID,
			OwnerID:      player1ID,
		}
		gameState.cards[land.ID] = land
		player1.Graveyard = append(player1.Graveyard, land)

		gameState.mu.Unlock()

		power, err := engine.getCreaturePower(gameState, tarmogoyf)
		assert.NoError(t, err)
		assert.Equal(t, 6, power, "Tarmogoyf should have 6 power with 6 card types (Instant, Creature, Sorcery, Artifact, Enchantment, Land)")

		toughness, err := engine.getCreatureToughness(gameState, tarmogoyf)
		assert.NoError(t, err)
		assert.Equal(t, 7, toughness, "Tarmogoyf should have 7 toughness with 6 card types")
	})
}

// TestCDAIntegration_MaroPowerToughness tests Maro-type CDAs (power/toughness equals hand size)
func TestCDAIntegration_MaroPowerToughness(t *testing.T) {
	logger := zap.NewNop()
	engine := NewMageEngine(logger)

	player1ID := uuid.New().String()
	player2ID := uuid.New().String()

	// Start a game
	gameID := uuid.New().String()
	err := engine.StartGame(gameID, []string{player1ID, player2ID}, "TwoPlayerDuel")
	require.NoError(t, err)

	// Get game state
	engine.mu.RLock()
	gameState, exists := engine.games[gameID]
	engine.mu.RUnlock()
	require.True(t, exists)

	// Create a Maro-like card (P/T = hand size)
	maroID := uuid.New()
	player1UUID := uuid.MustParse(player1ID)
	maroCDA := abilities.NewHandSizeCDA(maroID, player1UUID)

	// Register the CDA
	gameState.mu.Lock()
	gameState.abilityRegistry.RegisterAbility(
		maroCDA,
		player1UUID,
		0,
		abilities.ZoneBattlefield,
	)

	// Create the Maro card
	maro := &internalCard{
		ID:           maroID.String(),
		Name:         "Maro",
		ManaCost:     "{2}{G}{G}",
		Type:         "Creature",
		SubTypes:     []string{"Elemental"},
		Color:        "G",
		Power:        "*",
		Toughness:    "*",
		Zone:         zoneBattlefield,
		ControllerID: player1ID,
		OwnerID:      player1ID,
		Abilities: []EngineAbilityView{
			{
				ID:   maroCDA.GetID().String(),
				Text: "Maro's power and toughness are each equal to the number of cards in your hand.",
			},
		},
	}

	gameState.cards[maro.ID] = maro
	gameState.battlefield = append(gameState.battlefield, maro)

	// Player1 starts with 7 cards in hand (from StartGame)
	player1 := gameState.players[player1ID]
	initialHandSize := len(player1.Hand)
	gameState.mu.Unlock()

	// Test: Maro's P/T should equal hand size
	power, err := engine.getCreaturePower(gameState, maro)
	assert.NoError(t, err)
	assert.Equal(t, initialHandSize, power, "Maro's power should equal hand size")

	toughness, err := engine.getCreatureToughness(gameState, maro)
	assert.NoError(t, err)
	assert.Equal(t, initialHandSize, toughness, "Maro's toughness should equal hand size")
}
