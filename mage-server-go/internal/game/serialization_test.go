package game

import (
	"fmt"
	"testing"
	"time"

	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestComputeChecksum verifies that checksums are computed correctly
func TestComputeChecksum(t *testing.T) {
	snapshot := createTestSnapshot()

	checksum, err := snapshot.ComputeChecksum()
	require.NoError(t, err)
	assert.NotEmpty(t, checksum.Hash)
	assert.Equal(t, 1, checksum.Version)
	assert.NotEmpty(t, checksum.Timestamp)
}

// TestDeterministicChecksum verifies that identical snapshots produce identical checksums
// regardless of map iteration order (which is randomized in Go)
func TestDeterministicChecksum(t *testing.T) {
	// Create same snapshot multiple times
	checksums := make([]string, 10)
	for i := 0; i < 10; i++ {
		snapshot := createTestSnapshot()
		checksum, err := snapshot.ComputeChecksum()
		require.NoError(t, err)
		checksums[i] = checksum.Hash
	}

	// All checksums should be identical
	expected := checksums[0]
	for i := 1; i < len(checksums); i++ {
		assert.Equal(t, expected, checksums[i],
			"Checksum %d differs from checksum 0 - not deterministic", i)
	}
}

// TestChecksumDifferentStates verifies that different states produce different checksums
func TestChecksumDifferentStates(t *testing.T) {
	snapshot1 := createTestSnapshot()
	checksum1, err := snapshot1.ComputeChecksum()
	require.NoError(t, err)

	// Modify state
	snapshot2 := createTestSnapshot()
	snapshot2.TurnNumber = 5
	checksum2, err := snapshot2.ComputeChecksum()
	require.NoError(t, err)

	assert.NotEqual(t, checksum1.Hash, checksum2.Hash,
		"Different game states must produce different checksums")
}

// TestChecksumIgnoresTimestamp verifies that timestamp differences don't affect checksum
// Only the game state matters, not when it was captured
func TestChecksumIgnoresTimestamp(t *testing.T) {
	snapshot1 := createTestSnapshot()
	snapshot1.Timestamp = time.Now()
	checksum1, err := snapshot1.ComputeChecksum()
	require.NoError(t, err)

	snapshot2 := createTestSnapshot()
	snapshot2.Timestamp = time.Now().Add(1 * time.Hour)
	checksum2, err := snapshot2.ComputeChecksum()
	require.NoError(t, err)

	// Checksums should be the same because only timestamp differs
	assert.Equal(t, checksum1.Hash, checksum2.Hash,
		"Timestamp should not affect checksum")
}

// TestChecksumDetectsPlayerChanges verifies that player state changes affect checksum
func TestChecksumDetectsPlayerChanges(t *testing.T) {
	snapshot1 := createTestSnapshot()
	checksum1, err := snapshot1.ComputeChecksum()
	require.NoError(t, err)

	// Change player life
	snapshot2 := createTestSnapshot()
	snapshot2.Players["player1"].Life = 10
	checksum2, err := snapshot2.ComputeChecksum()
	require.NoError(t, err)

	assert.NotEqual(t, checksum1.Hash, checksum2.Hash,
		"Player life change must affect checksum")

	// Change player poison
	snapshot3 := createTestSnapshot()
	snapshot3.Players["player1"].Poison = 5
	checksum3, err := snapshot3.ComputeChecksum()
	require.NoError(t, err)

	assert.NotEqual(t, checksum1.Hash, checksum3.Hash,
		"Player poison change must affect checksum")
}

// TestChecksumDetectsCardChanges verifies that card state changes affect checksum
func TestChecksumDetectsCardChanges(t *testing.T) {
	snapshot1 := createTestSnapshot()
	checksum1, err := snapshot1.ComputeChecksum()
	require.NoError(t, err)

	// Change card tapped status
	snapshot2 := createTestSnapshot()
	snapshot2.Cards["card1"].Tapped = true
	checksum2, err := snapshot2.ComputeChecksum()
	require.NoError(t, err)

	assert.NotEqual(t, checksum1.Hash, checksum2.Hash,
		"Card tapped status change must affect checksum")

	// Change card damage
	snapshot3 := createTestSnapshot()
	snapshot3.Cards["card1"].Damage = 3
	checksum3, err := snapshot3.ComputeChecksum()
	require.NoError(t, err)

	assert.NotEqual(t, checksum1.Hash, checksum3.Hash,
		"Card damage change must affect checksum")

	// Change card zone
	snapshot4 := createTestSnapshot()
	snapshot4.Cards["card1"].Zone = zoneGraveyard
	checksum4, err := snapshot4.ComputeChecksum()
	require.NoError(t, err)

	assert.NotEqual(t, checksum1.Hash, checksum4.Hash,
		"Card zone change must affect checksum")
}

// TestSerializeDeserialize verifies basic serialization roundtrip
func TestSerializeDeserialize(t *testing.T) {
	snapshot := createTestSnapshot()

	// Serialize
	data, err := snapshot.SerializeToBytes()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Deserialize
	deserialized, err := DeserializeFromBytes(data)
	require.NoError(t, err)

	// Verify basic fields
	assert.Equal(t, snapshot.GameID, deserialized.GameID)
	assert.Equal(t, snapshot.TurnNumber, deserialized.TurnNumber)
	assert.Equal(t, snapshot.ActivePlayer, deserialized.ActivePlayer)
	assert.Equal(t, len(snapshot.Players), len(deserialized.Players))
	assert.Equal(t, len(snapshot.Cards), len(deserialized.Cards))
}

// TestValidateSerializationRoundtrip verifies that serialization preserves checksums
func TestValidateSerializationRoundtrip(t *testing.T) {
	snapshot := createTestSnapshot()

	err := ValidateSerializationRoundtrip(snapshot)
	assert.NoError(t, err, "Serialization roundtrip should preserve state")
}

// TestVerifyChecksum verifies checksum verification function
func TestVerifyChecksum(t *testing.T) {
	snapshot := createTestSnapshot()
	checksum, err := snapshot.ComputeChecksum()
	require.NoError(t, err)

	// Should match
	matches, err := snapshot.VerifyChecksum(checksum)
	require.NoError(t, err)
	assert.True(t, matches, "Checksum should match original")

	// Modify snapshot
	snapshot.TurnNumber = 999

	// Should not match
	matches, err = snapshot.VerifyChecksum(checksum)
	require.NoError(t, err)
	assert.False(t, matches, "Checksum should not match after modification")
}

// TestChecksumWithEmptyState verifies checksum works with minimal state
func TestChecksumWithEmptyState(t *testing.T) {
	snapshot := &gameStateSnapshot{
		GameID:   "empty-game",
		GameType: "test",
		State:    GameStateStarting,
		Players:  make(map[string]*internalPlayer),
		Cards:    make(map[string]*internalCard),
	}

	checksum, err := snapshot.ComputeChecksum()
	require.NoError(t, err)
	assert.NotEmpty(t, checksum.Hash)
}

// TestChecksumWithComplexState verifies checksum handles complex game state
func TestChecksumWithComplexState(t *testing.T) {
	snapshot := createComplexTestSnapshot()

	checksum1, err := snapshot.ComputeChecksum()
	require.NoError(t, err)

	// Serialize and deserialize
	data, err := snapshot.SerializeToBytes()
	require.NoError(t, err)

	deserialized, err := DeserializeFromBytes(data)
	require.NoError(t, err)

	checksum2, err := deserialized.ComputeChecksum()
	require.NoError(t, err)

	assert.Equal(t, checksum1.Hash, checksum2.Hash,
		"Complex state should survive serialization roundtrip")
}

// TestChecksumCardOrder verifies that card order in zones doesn't affect checksum
// (except stack where order matters)
func TestChecksumCardOrder(t *testing.T) {
	snapshot1 := createTestSnapshot()
	snapshot1.Battlefield = []*internalCard{
		snapshot1.Cards["card1"],
		snapshot1.Cards["card2"],
	}
	checksum1, err := snapshot1.ComputeChecksum()
	require.NoError(t, err)

	snapshot2 := createTestSnapshot()
	snapshot2.Battlefield = []*internalCard{
		snapshot2.Cards["card2"],
		snapshot2.Cards["card1"],
	}
	checksum2, err := snapshot2.ComputeChecksum()
	require.NoError(t, err)

	// Should be the same because battlefield order doesn't matter
	// (we sort by card ID in the deterministic representation)
	assert.Equal(t, checksum1.Hash, checksum2.Hash,
		"Battlefield card order should not affect checksum")
}

// TestChecksumMapOrder verifies that player/card map order doesn't affect checksum
func TestChecksumMapOrder(t *testing.T) {
	// Create snapshots with players added in different orders
	// In Go, map iteration order is randomized, but our checksum should be deterministic
	checksums := make([]string, 5)
	for i := 0; i < 5; i++ {
		snapshot := &gameStateSnapshot{
			GameID:      "game1",
			GameType:    "test",
			State:       GameStateInProgress,
			TurnNumber:  1,
			Players:     make(map[string]*internalPlayer),
			Cards:       make(map[string]*internalCard),
			PlayerOrder: []string{"p1", "p2", "p3"},
		}

		// Add players in different orders each iteration
		players := []string{"p1", "p2", "p3"}
		for _, pid := range players {
			snapshot.Players[pid] = &internalPlayer{
				PlayerID:  pid,
				Name:      "Player " + pid,
				Life:      20,
				Hand:      []*internalCard{},
				Library:   []*internalCard{},
				Graveyard: []*internalCard{},
			}
		}

		checksum, err := snapshot.ComputeChecksum()
		require.NoError(t, err)
		checksums[i] = checksum.Hash
	}

	// All checksums should be identical despite map iteration randomness
	expected := checksums[0]
	for i := 1; i < len(checksums); i++ {
		assert.Equal(t, expected, checksums[i],
			"Checksums must be deterministic regardless of map order")
	}
}

// createTestSnapshot creates a simple test snapshot
func createTestSnapshot() *gameStateSnapshot {
	player1 := &internalPlayer{
		PlayerID:  "player1",
		Name:      "Alice",
		Life:      20,
		Poison:    0,
		Energy:    0,
		Passed:    false,
		Lost:      false,
		Hand:      []*internalCard{},
		Library:   []*internalCard{},
		Graveyard: []*internalCard{},
	}

	card1 := &internalCard{
		ID:                "card1",
		Name:              "Grizzly Bears",
		OwnerID:           "player1",
		ControllerID:      "player1",
		Zone:              zoneBattlefield,
		Type:              "Creature",
		SubTypes:          []string{"Bear"},
		Power:             "2",
		Toughness:         "2",
		Damage:            0,
		Tapped:            false,
		SummoningSickness: false,
		Counters:          counters.NewCounters(),
		Abilities:         []EngineAbilityView{},
	}

	card2 := &internalCard{
		ID:                "card2",
		Name:              "Forest",
		OwnerID:           "player1",
		ControllerID:      "player1",
		Zone:              zoneBattlefield,
		Type:              "Land",
		SubTypes:          []string{},
		Power:             "",
		Toughness:         "",
		Damage:            0,
		Tapped:            false,
		SummoningSickness: false,
		Counters:          counters.NewCounters(),
		Abilities:         []EngineAbilityView{},
	}

	return &gameStateSnapshot{
		GameID:         "game1",
		GameType:       "test",
		State:          GameStateInProgress,
		TurnNumber:     3,
		ActivePlayer:   "player1",
		PriorityPlayer: "player1",
		Players: map[string]*internalPlayer{
			"player1": player1,
		},
		Cards: map[string]*internalCard{
			"card1": card1,
			"card2": card2,
		},
		Battlefield: []*internalCard{card1, card2},
		Exile:       []*internalCard{},
		Command:     []*internalCard{},
		StackItems:  []rules.StackItem{},
		Messages:    []EngineMessage{},
		Prompts:     []EnginePrompt{},
		PlayerOrder: []string{"player1"},
		Timestamp:   time.Now(),
	}
}

// createComplexTestSnapshot creates a more complex test snapshot with multiple players and cards
func createComplexTestSnapshot() *gameStateSnapshot {
	snapshot := &gameStateSnapshot{
		GameID:         "complex-game",
		GameType:       "standard",
		State:          GameStateInProgress,
		TurnNumber:     10,
		ActivePlayer:   "player1",
		PriorityPlayer: "player2",
		Players:        make(map[string]*internalPlayer),
		Cards:          make(map[string]*internalCard),
		Battlefield:    []*internalCard{},
		Exile:          []*internalCard{},
		Command:        []*internalCard{},
		StackItems:     []rules.StackItem{},
		Messages:       []EngineMessage{},
		Prompts:        []EnginePrompt{},
		PlayerOrder:    []string{"player1", "player2"},
		Timestamp:      time.Now(),
	}

	// Add two players
	for _, pid := range []string{"player1", "player2"} {
		player := &internalPlayer{
			PlayerID:  pid,
			Name:      "Player " + pid,
			Life:      15,
			Poison:    1,
			Energy:    3,
			Passed:    false,
			Lost:      false,
			Hand:      []*internalCard{},
			Library:   []*internalCard{},
			Graveyard: []*internalCard{},
		}
		snapshot.Players[pid] = player
	}

	// Add multiple cards
	for i := 1; i <= 5; i++ {
		cardID := fmt.Sprintf("card%d", i)
		card := &internalCard{
			ID:                cardID,
			Name:              fmt.Sprintf("Test Card %d", i),
			OwnerID:           "player1",
			ControllerID:      "player1",
			Zone:              zoneBattlefield,
			Type:              "Creature",
			SubTypes:          []string{"Test"},
			Power:             fmt.Sprintf("%d", i),
			Toughness:         fmt.Sprintf("%d", i),
			Damage:            0,
			Tapped:            i%2 == 0,
			SummoningSickness: false,
			Counters:          counters.NewCounters(),
			Abilities:         []EngineAbilityView{},
		}
		if i%3 == 0 {
			card.Counters.AddCounter(counters.NewCounter("+1/+1", 2))
		}
		snapshot.Cards[cardID] = card
		snapshot.Battlefield = append(snapshot.Battlefield, card)
	}

	return snapshot
}
