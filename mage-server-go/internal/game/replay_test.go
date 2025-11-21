package game

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewReplay(t *testing.T) {
	replay := NewReplay("game-123")
	assert.Equal(t, "game-123", replay.GameID)
	assert.Equal(t, 0, replay.CurrentIndex)
	assert.Equal(t, 0, len(replay.States))
}

func TestReplayRecordState(t *testing.T) {
	replay := NewReplay("game-123")

	// Create a test snapshot
	snapshot := &gameStateSnapshot{
		GameID:     "game-123",
		TurnNumber: 1,
		State:      GameStateInProgress,
	}

	replay.RecordState(snapshot)

	assert.Equal(t, 1, replay.Size())
	assert.Equal(t, snapshot, replay.States[0])
}

func TestReplayNavigation(t *testing.T) {
	replay := NewReplay("game-123")

	// Record 5 states
	for i := 0; i < 5; i++ {
		snapshot := &gameStateSnapshot{
			GameID:     "game-123",
			TurnNumber: i + 1,
			State:      GameStateInProgress,
		}
		replay.RecordState(snapshot)
	}

	assert.Equal(t, 5, replay.Size())

	// Test Start - resets to beginning
	replay.Start()
	assert.Equal(t, 0, replay.CurrentIndex)

	// Test Next
	state := replay.Next()
	assert.NotNil(t, state)
	assert.Equal(t, 1, state.TurnNumber)
	assert.Equal(t, 1, replay.CurrentIndex)

	state = replay.Next()
	assert.NotNil(t, state)
	assert.Equal(t, 2, state.TurnNumber)
	assert.Equal(t, 2, replay.CurrentIndex)

	// Test Previous
	// After next() twice, CurrentIndex is 2. Previous() should decrement to 1 and return state[1] (turn 2)
	state = replay.Previous()
	assert.NotNil(t, state)
	assert.Equal(t, 2, state.TurnNumber)
	assert.Equal(t, 1, replay.CurrentIndex)

	// Call Previous() again to get turn 1
	state = replay.Previous()
	assert.NotNil(t, state)
	assert.Equal(t, 1, state.TurnNumber)
	assert.Equal(t, 0, replay.CurrentIndex)

	// Test Previous at beginning
	replay.Start()
	state = replay.Previous()
	assert.Nil(t, state)
	assert.Equal(t, 0, replay.CurrentIndex)

	// Test Next at end
	for i := 0; i < 10; i++ {
		replay.Next()
	}
	state = replay.Next()
	assert.Nil(t, state)
}

func TestReplaySkip(t *testing.T) {
	replay := NewReplay("game-123")

	// Record 10 states
	for i := 0; i < 10; i++ {
		snapshot := &gameStateSnapshot{
			GameID:     "game-123",
			TurnNumber: i + 1,
			State:      GameStateInProgress,
		}
		replay.RecordState(snapshot)
	}

	// Start at beginning
	replay.Start()

	// Skip forward 3
	state := replay.Skip(3)
	assert.NotNil(t, state)
	assert.Equal(t, 4, state.TurnNumber) // Index 3 = turn 4
	assert.Equal(t, 3, replay.CurrentIndex)

	// Skip forward 5 more
	state = replay.Skip(5)
	assert.NotNil(t, state)
	assert.Equal(t, 9, state.TurnNumber)
	assert.Equal(t, 8, replay.CurrentIndex)

	// Skip forward beyond end
	state = replay.Skip(100)
	assert.NotNil(t, state)
	assert.Equal(t, 10, state.TurnNumber)
	assert.Equal(t, 9, replay.CurrentIndex)

	// Skip backward
	state = replay.Skip(-5)
	assert.NotNil(t, state)
	assert.Equal(t, 5, state.TurnNumber)
	assert.Equal(t, 4, replay.CurrentIndex)

	// Skip backward beyond beginning
	state = replay.Skip(-100)
	assert.NotNil(t, state)
	assert.Equal(t, 1, state.TurnNumber)
	assert.Equal(t, 0, replay.CurrentIndex)
}

func TestReplayGetStateAt(t *testing.T) {
	replay := NewReplay("game-123")

	// Record 5 states
	for i := 0; i < 5; i++ {
		snapshot := &gameStateSnapshot{
			GameID:     "game-123",
			TurnNumber: i + 1,
			State:      GameStateInProgress,
		}
		replay.RecordState(snapshot)
	}

	// Valid indices
	state := replay.GetStateAt(0)
	assert.NotNil(t, state)
	assert.Equal(t, 1, state.TurnNumber)

	state = replay.GetStateAt(4)
	assert.NotNil(t, state)
	assert.Equal(t, 5, state.TurnNumber)

	// Invalid indices
	state = replay.GetStateAt(-1)
	assert.Nil(t, state)

	state = replay.GetStateAt(5)
	assert.Nil(t, state)
}

func TestReplaySaveAndLoad(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	replay := NewReplay("game-123")

	// Record 5 states
	for i := 0; i < 5; i++ {
		snapshot := &gameStateSnapshot{
			GameID:     "game-123",
			GameType:   "standard",
			TurnNumber: i + 1,
			State:      GameStateInProgress,
		}
		replay.RecordState(snapshot)
	}

	// Save to file
	err := replay.SaveToFile(tempDir)
	require.NoError(t, err)

	// Verify file exists
	filename := filepath.Join(tempDir, "game-123.replay")
	_, err = os.Stat(filename)
	require.NoError(t, err)

	// Load from file
	loadedReplay, err := LoadReplayFromFile(tempDir, "game-123")
	require.NoError(t, err)
	require.NotNil(t, loadedReplay)

	// Verify loaded replay matches original
	assert.Equal(t, replay.GameID, loadedReplay.GameID)
	assert.Equal(t, replay.Size(), loadedReplay.Size())

	for i := 0; i < replay.Size(); i++ {
		originalState := replay.GetStateAt(i)
		loadedState := loadedReplay.GetStateAt(i)

		assert.Equal(t, originalState.GameID, loadedState.GameID)
		assert.Equal(t, originalState.GameType, loadedState.GameType)
		assert.Equal(t, originalState.TurnNumber, loadedState.TurnNumber)
		assert.Equal(t, originalState.State, loadedState.State)
	}
}

func TestReplaySaveNonexistentDirectory(t *testing.T) {
	replay := NewReplay("game-123")

	// Record a state
	snapshot := &gameStateSnapshot{
		GameID:     "game-123",
		TurnNumber: 1,
		State:      GameStateInProgress,
	}
	replay.RecordState(snapshot)

	// Save to nonexistent directory (should create it)
	tempDir := filepath.Join(t.TempDir(), "subdir", "another")
	err := replay.SaveToFile(tempDir)
	require.NoError(t, err)

	// Verify directory and file exist
	filename := filepath.Join(tempDir, "game-123.replay")
	_, err = os.Stat(filename)
	require.NoError(t, err)
}

func TestReplayLoadNonexistentFile(t *testing.T) {
	tempDir := t.TempDir()

	// Try to load nonexistent replay
	_, err := LoadReplayFromFile(tempDir, "nonexistent")
	assert.Error(t, err)
}

func TestReplayRecorder(t *testing.T) {
	logger := zap.NewNop()
	tempDir := t.TempDir()

	recorder := NewReplayRecorder(logger, tempDir)
	require.NotNil(t, recorder)

	gameID := "game-123"

	// Start recording
	recorder.StartRecording(gameID)
	assert.True(t, recorder.IsRecording(gameID))

	// Record some states
	for i := 0; i < 5; i++ {
		snapshot := &gameStateSnapshot{
			GameID:     gameID,
			TurnNumber: i + 1,
			State:      GameStateInProgress,
		}
		recorder.RecordState(gameID, snapshot)
	}

	// Get replay
	replay, exists := recorder.GetReplay(gameID)
	require.True(t, exists)
	require.NotNil(t, replay)
	assert.Equal(t, 5, replay.Size())

	// Stop recording
	recorder.StopRecording(gameID)
	assert.False(t, recorder.IsRecording(gameID))

	// Recording should still exist in memory
	replay, exists = recorder.GetReplay(gameID)
	require.True(t, exists)
	assert.Equal(t, 5, replay.Size())

	// States recorded after stopping should be ignored
	snapshot := &gameStateSnapshot{
		GameID:     gameID,
		TurnNumber: 6,
		State:      GameStateInProgress,
	}
	recorder.RecordState(gameID, snapshot)

	replay, exists = recorder.GetReplay(gameID)
	require.True(t, exists)
	assert.Equal(t, 5, replay.Size()) // Still 5, not 6

	// Save replay
	err := recorder.SaveReplay(gameID)
	require.NoError(t, err)

	// Replay should be removed from memory after save
	_, exists = recorder.GetReplay(gameID)
	assert.False(t, exists)

	// Load replay
	loadedReplay, err := recorder.LoadReplay(gameID)
	require.NoError(t, err)
	require.NotNil(t, loadedReplay)
	assert.Equal(t, 5, loadedReplay.Size())
}

func TestReplayRecorderClear(t *testing.T) {
	logger := zap.NewNop()
	tempDir := t.TempDir()

	recorder := NewReplayRecorder(logger, tempDir)
	gameID := "game-123"

	// Start recording and record states
	recorder.StartRecording(gameID)
	for i := 0; i < 3; i++ {
		snapshot := &gameStateSnapshot{
			GameID:     gameID,
			TurnNumber: i + 1,
		}
		recorder.RecordState(gameID, snapshot)
	}

	// Verify replay exists
	_, exists := recorder.GetReplay(gameID)
	assert.True(t, exists)

	// Clear replay
	recorder.ClearReplay(gameID)

	// Verify replay is gone
	_, exists = recorder.GetReplay(gameID)
	assert.False(t, exists)

	// Recording should also be disabled
	assert.False(t, recorder.IsRecording(gameID))
}

func TestReplayRecorderMultipleGames(t *testing.T) {
	logger := zap.NewNop()
	tempDir := t.TempDir()

	recorder := NewReplayRecorder(logger, tempDir)

	// Start recording multiple games
	game1 := "game-1"
	game2 := "game-2"
	game3 := "game-3"

	recorder.StartRecording(game1)
	recorder.StartRecording(game2)
	recorder.StartRecording(game3)

	// Record different numbers of states for each
	for i := 0; i < 3; i++ {
		recorder.RecordState(game1, &gameStateSnapshot{GameID: game1, TurnNumber: i + 1})
	}
	for i := 0; i < 5; i++ {
		recorder.RecordState(game2, &gameStateSnapshot{GameID: game2, TurnNumber: i + 1})
	}
	for i := 0; i < 7; i++ {
		recorder.RecordState(game3, &gameStateSnapshot{GameID: game3, TurnNumber: i + 1})
	}

	// Verify each game's replay
	replay1, exists := recorder.GetReplay(game1)
	require.True(t, exists)
	assert.Equal(t, 3, replay1.Size())

	replay2, exists := recorder.GetReplay(game2)
	require.True(t, exists)
	assert.Equal(t, 5, replay2.Size())

	replay3, exists := recorder.GetReplay(game3)
	require.True(t, exists)
	assert.Equal(t, 7, replay3.Size())

	// Save and verify each independently
	require.NoError(t, recorder.SaveReplay(game1))
	require.NoError(t, recorder.SaveReplay(game2))
	require.NoError(t, recorder.SaveReplay(game3))

	// Load and verify
	loaded1, err := recorder.LoadReplay(game1)
	require.NoError(t, err)
	assert.Equal(t, 3, loaded1.Size())

	loaded2, err := recorder.LoadReplay(game2)
	require.NoError(t, err)
	assert.Equal(t, 5, loaded2.Size())

	loaded3, err := recorder.LoadReplay(game3)
	require.NoError(t, err)
	assert.Equal(t, 7, loaded3.Size())
}
