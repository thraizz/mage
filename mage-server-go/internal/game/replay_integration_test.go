package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestReplayRecordingDuringGameplay(t *testing.T) {
	logger := zap.NewNop()
	engine := NewMageEngine(logger)

	gameID := "replay-test-game"
	players := []string{"player1", "player2"}

	// Start a game
	err := engine.StartGame(gameID, players, "standard")
	require.NoError(t, err)

	// Enable replay recording
	err = engine.StartReplayRecording(gameID)
	require.NoError(t, err)
	assert.True(t, engine.IsRecordingReplay(gameID))

	// Initial state should be recorded from StartGame
	replay, exists := engine.GetReplay(gameID)
	require.True(t, exists)
	initialSize := replay.Size()
	assert.Greater(t, initialSize, 0, "Initial game state should be recorded")

	// Pass priority a few times (should record states)
	for i := 0; i < 3; i++ {
		currentPlayer := engine.getPriorityPlayerForTest(gameID)
		err = engine.ProcessAction(gameID, PlayerAction{
			PlayerID:   currentPlayer,
			ActionType: "PLAYER_ACTION",
			Data:       "PASS",
		})
		require.NoError(t, err)
	}

	// Verify more states were recorded
	replay, exists = engine.GetReplay(gameID)
	require.True(t, exists)
	assert.Greater(t, replay.Size(), initialSize, "States should be recorded during gameplay")

	// Stop recording
	engine.StopReplayRecording(gameID)
	assert.False(t, engine.IsRecordingReplay(gameID))

	sizeBeforeStop := replay.Size()

	// Pass priority again - should not record
	currentPlayer := engine.getPriorityPlayerForTest(gameID)
	err = engine.ProcessAction(gameID, PlayerAction{
		PlayerID:   currentPlayer,
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
	})
	require.NoError(t, err)

	// Replay size should not have changed
	replay, exists = engine.GetReplay(gameID)
	require.True(t, exists)
	assert.Equal(t, sizeBeforeStop, replay.Size(), "No states should be recorded after stopping")
}

func TestReplayPlayback(t *testing.T) {
	logger := zap.NewNop()
	engine := NewMageEngine(logger)

	gameID := "replay-playback-test"
	players := []string{"player1", "player2"}

	// Start a game and enable recording
	err := engine.StartGame(gameID, players, "standard")
	require.NoError(t, err)

	err = engine.StartReplayRecording(gameID)
	require.NoError(t, err)

	// Play several turns
	for i := 0; i < 5; i++ {
		currentPlayer := engine.getPriorityPlayerForTest(gameID)
		err = engine.ProcessAction(gameID, PlayerAction{
			PlayerID:   currentPlayer,
			ActionType: "PLAYER_ACTION",
			Data:       "PASS",
		})
		require.NoError(t, err)
	}

	// Get the replay
	replay, exists := engine.GetReplay(gameID)
	require.True(t, exists)
	require.Greater(t, replay.Size(), 1)

	// Test playback navigation
	replay.Start()
	assert.Equal(t, 0, replay.CurrentIndex)

	// Navigate forward
	state1 := replay.Next()
	require.NotNil(t, state1)
	assert.Equal(t, gameID, state1.GameID)

	state2 := replay.Next()
	require.NotNil(t, state2)
	assert.Equal(t, gameID, state2.GameID)

	// Navigate backward
	statePrev := replay.Previous()
	require.NotNil(t, statePrev)
	assert.Equal(t, state1.GameID, statePrev.GameID)

	// Skip forward
	stateSkip := replay.Skip(3)
	require.NotNil(t, stateSkip)
	assert.Equal(t, gameID, stateSkip.GameID)
}

func TestReplaySaveAndLoadFullGame(t *testing.T) {
	logger := zap.NewNop()
	tempDir := t.TempDir()

	// Create engine with custom replay directory
	engine := NewMageEngine(logger)
	engine.replayRecorder = NewReplayRecorder(logger, tempDir)

	gameID := "full-game-test"
	players := []string{"player1", "player2"}

	// Start game and recording
	err := engine.StartGame(gameID, players, "standard")
	require.NoError(t, err)

	err = engine.StartReplayRecording(gameID)
	require.NoError(t, err)

	// Play through several actions
	for i := 0; i < 10; i++ {
		currentPlayer := engine.getPriorityPlayerForTest(gameID)
		err = engine.ProcessAction(gameID, PlayerAction{
			PlayerID:   currentPlayer,
			ActionType: "PLAYER_ACTION",
			Data:       "PASS",
		})
		require.NoError(t, err)
	}

	// Get replay and save it
	replay, exists := engine.GetReplay(gameID)
	require.True(t, exists)
	originalSize := replay.Size()
	require.Greater(t, originalSize, 1)

	err = engine.SaveReplayToFile(gameID)
	require.NoError(t, err)

	// Replay should be removed from memory after save
	_, exists = engine.GetReplay(gameID)
	assert.False(t, exists)

	// Load replay from file
	loadedReplay, err := engine.LoadReplayFromFile(gameID)
	require.NoError(t, err)
	require.NotNil(t, loadedReplay)

	// Verify loaded replay matches original
	assert.Equal(t, gameID, loadedReplay.GameID)
	assert.Equal(t, originalSize, loadedReplay.Size())

	// Verify we can navigate the loaded replay
	loadedReplay.Start()
	state := loadedReplay.Next()
	require.NotNil(t, state)
	assert.Equal(t, gameID, state.GameID)
}

func TestReplayRecordingDisabledByDefault(t *testing.T) {
	logger := zap.NewNop()
	engine := NewMageEngine(logger)

	gameID := "no-recording-test"
	players := []string{"player1", "player2"}

	// Start game without enabling recording
	err := engine.StartGame(gameID, players, "standard")
	require.NoError(t, err)

	// Recording should be disabled
	assert.False(t, engine.IsRecordingReplay(gameID))

	// Play a few actions
	for i := 0; i < 3; i++ {
		currentPlayer := engine.getPriorityPlayerForTest(gameID)
		err = engine.ProcessAction(gameID, PlayerAction{
			PlayerID:   currentPlayer,
			ActionType: "PLAYER_ACTION",
			Data:       "PASS",
		})
		require.NoError(t, err)
	}

	// No replay should exist
	_, exists := engine.GetReplay(gameID)
	assert.False(t, exists)
}

func TestReplayCleanupOnGameEnd(t *testing.T) {
	logger := zap.NewNop()
	engine := NewMageEngine(logger)

	gameID := "cleanup-test"
	players := []string{"player1", "player2"}

	// Start game and recording
	err := engine.StartGame(gameID, players, "standard")
	require.NoError(t, err)

	err = engine.StartReplayRecording(gameID)
	require.NoError(t, err)

	// Play a few actions
	for i := 0; i < 3; i++ {
		currentPlayer := engine.getPriorityPlayerForTest(gameID)
		err = engine.ProcessAction(gameID, PlayerAction{
			PlayerID:   currentPlayer,
			ActionType: "PLAYER_ACTION",
			Data:       "PASS",
		})
		require.NoError(t, err)
	}

	// Verify replay exists
	_, exists := engine.GetReplay(gameID)
	require.True(t, exists)

	// Cleanup game
	err = engine.CleanupGame(gameID)
	require.NoError(t, err)

	// Replay should be cleared from memory
	_, exists = engine.GetReplay(gameID)
	assert.False(t, exists)
}

// Helper function to get priority player from engine (for testing)
func (e *MageEngine) getPriorityPlayerForTest(gameID string) string {
	e.mu.RLock()
	gameState, exists := e.games[gameID]
	e.mu.RUnlock()

	if !exists {
		return ""
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	return gameState.turnManager.PriorityPlayer()
}
