package game

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// rollback.go implements bookmark/restore system for game state snapshots
// Allows players to roll back game state to previous points

// Bookmark represents a saved game state snapshot
type Bookmark struct {
	ID        string     `json:"id"`
	GameID    string     `json:"gameId"`
	State     *GameState `json:"state"`
	Timestamp time.Time  `json:"timestamp"`
	Label     string     `json:"label"`
	CreatedBy string     `json:"createdBy"`
}

// bookmarks stores all bookmarks for all games
type bookmarkManager struct {
	bookmarks map[string][]*Bookmark // gameID -> list of bookmarks
}

var globalBookmarkManager = &bookmarkManager{
	bookmarks: make(map[string][]*Bookmark),
}

// BookmarkState creates a snapshot of the current game state
func (e *GameEngine) BookmarkState(gameID, playerID, label string) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return "", fmt.Errorf("game not found: %s", gameID)
	}

	// Create deep copy of state
	stateCopy, err := deepCopyState(state)
	if err != nil {
		return "", fmt.Errorf("failed to copy state: %w", err)
	}

	// Generate bookmark ID
	bookmarkID := fmt.Sprintf("bookmark-%s-%d", gameID, time.Now().UnixNano())

	bookmark := &Bookmark{
		ID:        bookmarkID,
		GameID:    gameID,
		State:     stateCopy,
		Timestamp: time.Now(),
		Label:     label,
		CreatedBy: playerID,
	}

	// Store bookmark
	if globalBookmarkManager.bookmarks[gameID] == nil {
		globalBookmarkManager.bookmarks[gameID] = make([]*Bookmark, 0)
	}
	globalBookmarkManager.bookmarks[gameID] = append(globalBookmarkManager.bookmarks[gameID], bookmark)

	e.logger.Info("bookmark created",
		zap.String("game_id", gameID),
		zap.String("bookmark_id", bookmarkID),
		zap.String("label", label),
		zap.String("created_by", playerID))

	msg := fmt.Sprintf("%s created bookmark: %s", e.getPlayerName(state, playerID), label)
	e.appendLog(state, "bookmark", msg)
	e.broadcast(gameID)

	return bookmarkID, nil
}

// RestoreState restores game state from a bookmark
func (e *GameEngine) RestoreState(gameID, bookmarkID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Find bookmark
	bookmarks := globalBookmarkManager.bookmarks[gameID]
	var bookmark *Bookmark
	for _, b := range bookmarks {
		if b.ID == bookmarkID {
			bookmark = b
			break
		}
	}

	if bookmark == nil {
		return fmt.Errorf("bookmark not found: %s", bookmarkID)
	}

	// Create deep copy of bookmarked state
	restoredState, err := deepCopyState(bookmark.State)
	if err != nil {
		return fmt.Errorf("failed to restore state: %w", err)
	}

	// Restore the state
	e.games[gameID] = restoredState

	e.logger.Info("state restored",
		zap.String("game_id", gameID),
		zap.String("bookmark_id", bookmarkID),
		zap.String("label", bookmark.Label))

	msg := fmt.Sprintf("Game state restored to: %s", bookmark.Label)
	e.appendLog(restoredState, "restore", msg)
	e.broadcast(gameID)

	return nil
}

// ListBookmarks returns all bookmarks for a game
func (e *GameEngine) ListBookmarks(gameID string) ([]*Bookmark, error) {
	bookmarks := globalBookmarkManager.bookmarks[gameID]
	if bookmarks == nil {
		return make([]*Bookmark, 0), nil
	}

	// Return copies without full state data (to save memory)
	result := make([]*Bookmark, len(bookmarks))
	for i, b := range bookmarks {
		result[i] = &Bookmark{
			ID:        b.ID,
			GameID:    b.GameID,
			Timestamp: b.Timestamp,
			Label:     b.Label,
			CreatedBy: b.CreatedBy,
			// Don't include State in list response
		}
	}

	return result, nil
}

// DeleteBookmark deletes a bookmark
func (e *GameEngine) DeleteBookmark(gameID, bookmarkID string) error {
	bookmarks := globalBookmarkManager.bookmarks[gameID]
	if bookmarks == nil {
		return fmt.Errorf("no bookmarks found for game: %s", gameID)
	}

	// Find and remove bookmark
	newBookmarks := make([]*Bookmark, 0, len(bookmarks))
	found := false
	for _, b := range bookmarks {
		if b.ID == bookmarkID {
			found = true
		} else {
			newBookmarks = append(newBookmarks, b)
		}
	}

	if !found {
		return fmt.Errorf("bookmark not found: %s", bookmarkID)
	}

	globalBookmarkManager.bookmarks[gameID] = newBookmarks

	e.logger.Info("bookmark deleted",
		zap.String("game_id", gameID),
		zap.String("bookmark_id", bookmarkID))

	return nil
}

// CleanupBookmarks removes all bookmarks for a game (called when game ends)
func (e *GameEngine) CleanupBookmarks(gameID string) {
	delete(globalBookmarkManager.bookmarks, gameID)

	e.logger.Info("bookmarks cleaned up",
		zap.String("game_id", gameID))
}

// AutoBookmark creates automatic bookmarks at turn boundaries
func (e *GameEngine) AutoBookmark(gameID, playerID string) error {
	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	label := fmt.Sprintf("Turn %d - Start of %s's turn", state.Turn, e.getPlayerName(state, playerID))
	_, err := e.BookmarkState(gameID, playerID, label)
	return err
}

// deepCopyState creates a deep copy of game state using JSON serialization
// This ensures complete independence between original and copy
func deepCopyState(state *GameState) (*GameState, error) {
	// Serialize to JSON
	data, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	// Deserialize to new object
	var copy GameState
	if err := json.Unmarshal(data, &copy); err != nil {
		return nil, err
	}

	return &copy, nil
}

// GetBookmark retrieves a specific bookmark
func (e *GameEngine) GetBookmark(gameID, bookmarkID string) (*Bookmark, error) {
	bookmarks := globalBookmarkManager.bookmarks[gameID]
	if bookmarks == nil {
		return nil, fmt.Errorf("no bookmarks found for game: %s", gameID)
	}

	for _, b := range bookmarks {
		if b.ID == bookmarkID {
			return b, nil
		}
	}

	return nil, fmt.Errorf("bookmark not found: %s", bookmarkID)
}

// BookmarkWithConsent creates a bookmark that requires player consent before restore
type ConsentRequest struct {
	BookmarkID  string          `json:"bookmarkId"`
	RequestedBy string          `json:"requestedBy"`
	Timestamp   time.Time       `json:"timestamp"`
	Consents    map[string]bool `json:"consents"` // playerID -> consented
	Required    int             `json:"required"` // Number of consents required
}

var consentRequests = make(map[string]*ConsentRequest) // bookmarkID -> consent request

// RequestRestore requests restoration of a bookmark with player consent
func (e *GameEngine) RequestRestore(gameID, bookmarkID, playerID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	// Check if bookmark exists
	_, err := e.GetBookmark(gameID, bookmarkID)
	if err != nil {
		return err
	}

	// Create consent request
	required := len(state.Players) - 1 // All players except requester
	if required <= 0 {
		// Single player - no consent needed
		return e.RestoreState(gameID, bookmarkID)
	}

	consentRequests[bookmarkID] = &ConsentRequest{
		BookmarkID:  bookmarkID,
		RequestedBy: playerID,
		Timestamp:   time.Now(),
		Consents:    make(map[string]bool),
		Required:    required,
	}

	msg := fmt.Sprintf("%s requests to restore game state (requires %d consents)",
		e.getPlayerName(state, playerID), required)
	e.appendLog(state, "consentRequest", msg)
	e.broadcast(gameID)

	return nil
}

// ConsentToRestore records a player's consent to restore
func (e *GameEngine) ConsentToRestore(gameID, bookmarkID, playerID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	request := consentRequests[bookmarkID]
	if request == nil {
		return fmt.Errorf("no consent request found for bookmark: %s", bookmarkID)
	}

	// Record consent
	request.Consents[playerID] = true

	msg := fmt.Sprintf("%s consents to restore", e.getPlayerName(state, playerID))
	e.appendLog(state, "consent", msg)

	// Check if we have enough consents
	if len(request.Consents) >= request.Required {
		// All consents received - restore state
		delete(consentRequests, bookmarkID)
		return e.RestoreState(gameID, bookmarkID)
	}

	e.broadcast(gameID)
	return nil
}

// DenyRestore denies a restore request
func (e *GameEngine) DenyRestore(gameID, bookmarkID, playerID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.games[gameID]
	if state == nil {
		return fmt.Errorf("game not found: %s", gameID)
	}

	request := consentRequests[bookmarkID]
	if request == nil {
		return fmt.Errorf("no consent request found for bookmark: %s", bookmarkID)
	}

	// Remove consent request
	delete(consentRequests, bookmarkID)

	msg := fmt.Sprintf("%s denies restore request", e.getPlayerName(state, playerID))
	e.appendLog(state, "denyRestore", msg)
	e.broadcast(gameID)

	return nil
}
