package game

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Replay represents a recorded game with sequential state snapshots
// Per Java mage.game.GameReplay: stores list of GameState snapshots for playback
type Replay struct {
	GameID       string
	States       []*gameStateSnapshot
	CurrentIndex int
	mu           sync.RWMutex
}

// NewReplay creates a new replay instance
func NewReplay(gameID string) *Replay {
	return &Replay{
		GameID:       gameID,
		States:       make([]*gameStateSnapshot, 0),
		CurrentIndex: 0,
	}
}

// RecordState adds a new state snapshot to the replay
// Per Java GameStates.save(): appends copy of current state
func (r *Replay) RecordState(snapshot *gameStateSnapshot) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.States = append(r.States, snapshot)
}

// Start resets the replay to the beginning
// Per Java GameReplay.start(): sets index to 0
func (r *Replay) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.CurrentIndex = 0
}

// Next moves to the next state and returns it
// Per Java GameReplay.next(): increments index and returns state
func (r *Replay) Next() *gameStateSnapshot {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.CurrentIndex < len(r.States) {
		state := r.States[r.CurrentIndex]
		r.CurrentIndex++
		return state
	}
	return nil
}

// Previous moves to the previous state and returns it
// Per Java GameReplay.previous(): decrements index and returns state
func (r *Replay) Previous() *gameStateSnapshot {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.CurrentIndex > 0 {
		r.CurrentIndex--
		return r.States[r.CurrentIndex]
	}
	return nil
}

// Skip moves forward by the specified number of states
// Per Java ReplaySession.next(int moves): skip multiple states
func (r *Replay) Skip(count int) *gameStateSnapshot {
	r.mu.Lock()
	defer r.mu.Unlock()

	newIndex := r.CurrentIndex + count
	if newIndex >= len(r.States) {
		newIndex = len(r.States) - 1
	}
	if newIndex < 0 {
		newIndex = 0
	}

	r.CurrentIndex = newIndex
	if r.CurrentIndex < len(r.States) {
		return r.States[r.CurrentIndex]
	}
	return nil
}

// Size returns the number of recorded states
func (r *Replay) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.States)
}

// GetStateAt returns the state at a specific index
func (r *Replay) GetStateAt(index int) *gameStateSnapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if index >= 0 && index < len(r.States) {
		return r.States[index]
	}
	return nil
}

// SaveToFile saves the replay to a gzipped file
// Per Java GameReplay.loadGame(): saves Game and GameStates to gzip file
func (r *Replay) SaveToFile(directory string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create directory if it doesn't exist
	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	filename := filepath.Join(directory, fmt.Sprintf("%s.replay", r.GameID))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Create gzip writer
	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	// Create gob encoder
	encoder := gob.NewEncoder(gzipWriter)

	// Encode metadata
	metadata := replayMetadata{
		GameID:     r.GameID,
		Timestamp:  time.Now(),
		Version:    1,
		StateCount: len(r.States),
	}
	if err := encoder.Encode(&metadata); err != nil {
		return fmt.Errorf("failed to encode metadata: %w", err)
	}

	// Encode all states
	for i, state := range r.States {
		if err := encoder.Encode(state); err != nil {
			return fmt.Errorf("failed to encode state %d: %w", i, err)
		}
	}

	return nil
}

// LoadFromFile loads a replay from a gzipped file
// Per Java GameReplay.loadGame(): loads Game and GameStates from gzip file
func LoadReplayFromFile(directory, gameID string) (*Replay, error) {
	filename := filepath.Join(directory, fmt.Sprintf("%s.replay", gameID))

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	// Create gob decoder
	decoder := gob.NewDecoder(gzipReader)

	// Decode metadata
	var metadata replayMetadata
	if err := decoder.Decode(&metadata); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	// Verify version
	if metadata.Version != 1 {
		return nil, fmt.Errorf("unsupported replay version: %d", metadata.Version)
	}

	// Create replay
	replay := NewReplay(metadata.GameID)

	// Decode all states
	for i := 0; i < metadata.StateCount; i++ {
		var state gameStateSnapshot
		if err := decoder.Decode(&state); err != nil {
			return nil, fmt.Errorf("failed to decode state %d: %w", i, err)
		}
		replay.States = append(replay.States, &state)
	}

	return replay, nil
}

// replayMetadata contains information about a saved replay
type replayMetadata struct {
	GameID     string
	Timestamp  time.Time
	Version    int
	StateCount int
}

// ReplayRecorder manages replay recording for the engine
type ReplayRecorder struct {
	logger  *zap.Logger
	mu      sync.RWMutex
	replays map[string]*Replay // gameID -> Replay
	enabled map[string]bool    // gameID -> whether recording is enabled
	saveDir string             // Directory to save replay files
}

// NewReplayRecorder creates a new replay recorder
func NewReplayRecorder(logger *zap.Logger, saveDir string) *ReplayRecorder {
	return &ReplayRecorder{
		logger:  logger,
		replays: make(map[string]*Replay),
		enabled: make(map[string]bool),
		saveDir: saveDir,
	}
}

// StartRecording begins recording a game
func (rr *ReplayRecorder) StartRecording(gameID string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	rr.replays[gameID] = NewReplay(gameID)
	rr.enabled[gameID] = true

	if rr.logger != nil {
		rr.logger.Info("started replay recording",
			zap.String("game_id", gameID),
		)
	}
}

// StopRecording stops recording a game
func (rr *ReplayRecorder) StopRecording(gameID string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	rr.enabled[gameID] = false

	if rr.logger != nil {
		rr.logger.Info("stopped replay recording",
			zap.String("game_id", gameID),
		)
	}
}

// RecordState records a game state snapshot if recording is enabled
func (rr *ReplayRecorder) RecordState(gameID string, snapshot *gameStateSnapshot) {
	rr.mu.RLock()
	enabled := rr.enabled[gameID]
	replay := rr.replays[gameID]
	rr.mu.RUnlock()

	if !enabled || replay == nil {
		return
	}

	replay.RecordState(snapshot)

	if rr.logger != nil {
		rr.logger.Debug("recorded replay state",
			zap.String("game_id", gameID),
			zap.Int("state_count", replay.Size()),
		)
	}
}

// GetReplay returns the replay for a game
func (rr *ReplayRecorder) GetReplay(gameID string) (*Replay, bool) {
	rr.mu.RLock()
	defer rr.mu.RUnlock()

	replay, exists := rr.replays[gameID]
	return replay, exists
}

// SaveReplay saves a replay to disk and removes it from memory
func (rr *ReplayRecorder) SaveReplay(gameID string) error {
	rr.mu.Lock()
	replay, exists := rr.replays[gameID]
	if !exists {
		rr.mu.Unlock()
		return fmt.Errorf("no replay found for game %s", gameID)
	}
	delete(rr.replays, gameID)
	delete(rr.enabled, gameID)
	rr.mu.Unlock()

	if err := replay.SaveToFile(rr.saveDir); err != nil {
		return fmt.Errorf("failed to save replay: %w", err)
	}

	if rr.logger != nil {
		rr.logger.Info("saved replay to disk",
			zap.String("game_id", gameID),
			zap.Int("state_count", replay.Size()),
			zap.String("directory", rr.saveDir),
		)
	}

	return nil
}

// LoadReplay loads a replay from disk
func (rr *ReplayRecorder) LoadReplay(gameID string) (*Replay, error) {
	replay, err := LoadReplayFromFile(rr.saveDir, gameID)
	if err != nil {
		return nil, err
	}

	if rr.logger != nil {
		rr.logger.Info("loaded replay from disk",
			zap.String("game_id", gameID),
			zap.Int("state_count", replay.Size()),
		)
	}

	return replay, nil
}

// ClearReplay removes a replay from memory without saving
func (rr *ReplayRecorder) ClearReplay(gameID string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	delete(rr.replays, gameID)
	delete(rr.enabled, gameID)

	if rr.logger != nil {
		rr.logger.Debug("cleared replay from memory",
			zap.String("game_id", gameID),
		)
	}
}

// IsRecording returns whether recording is enabled for a game
func (rr *ReplayRecorder) IsRecording(gameID string) bool {
	rr.mu.RLock()
	defer rr.mu.RUnlock()

	return rr.enabled[gameID]
}
