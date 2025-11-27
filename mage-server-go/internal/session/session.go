package session

import (
	"sync"
	"time"
)

// Preferences captures user-specific settings stored in the session.
type Preferences struct {
	AvatarID                 int32
	ShowAbsoluteAbilities    bool
	AllowRequestsFromFriends bool
	ConfirmEmptyManaPool     bool
	UserGroup                string
	SkipPrioritySteps        []string
	FlagsName                string
	AskMoveToGraveOrder      int32
}

// Session represents a user session
type Session struct {
	ID           string
	UserID       string
	Host         string
	IsAdmin      bool
	Connected    bool
	LastActivity time.Time
	LeasePeriod  time.Duration
	CallbackChan chan interface{} // Channel for WebSocket callbacks
	wsCloseChan  chan struct{}    // Signal to close current WebSocket
	wsCloseFunc  func()           // Safe close function for the WebSocket (uses sync.Once internally)
	preferences  *Preferences
	mu           sync.RWMutex
	reqMu        sync.Mutex // Prevents concurrent requests for same session
}

// NewSession creates a new session
func NewSession(id, host string, leasePeriod time.Duration) *Session {
	return &Session{
		ID:           id,
		Host:         host,
		Connected:    true,
		LastActivity: time.Now(),
		LeasePeriod:  leasePeriod,
		CallbackChan: make(chan interface{}, 100),
		wsCloseChan:  nil, // Will be set when WebSocket connects
	}
}

// SetWebSocketCloseFunc sets the close function for the active WebSocket.
// The closeFunc should be a function that safely closes the done channel using sync.Once.
// Returns the old close function if one exists (to close previous connection safely).
func (s *Session) SetWebSocketCloseFunc(closeChan chan struct{}, closeFunc func()) func() {
	s.mu.Lock()
	defer s.mu.Unlock()
	oldFunc := s.wsCloseFunc
	s.wsCloseChan = closeChan
	s.wsCloseFunc = closeFunc
	return oldFunc
}

// ClearWebSocketCloseFunc clears the WebSocket close function only if the channel matches
func (s *Session) ClearWebSocketCloseFunc(closeChan chan struct{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Only clear if this is still the active close channel
	if s.wsCloseChan == closeChan {
		s.wsCloseChan = nil
		s.wsCloseFunc = nil
	}
}

// UpdateActivity updates the last activity timestamp
func (s *Session) UpdateActivity() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LastActivity = time.Now()
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return time.Since(s.LastActivity) > s.LeasePeriod
}

// SendCallback sends a callback event to the session
// Returns false if the channel is full (timeout after 5 seconds)
func (s *Session) SendCallback(event interface{}) bool {
	select {
	case s.CallbackChan <- event:
		return true
	case <-time.After(5 * time.Second):
		return false // Timeout, client not reading
	}
}

// SetUserID sets the user ID for the session
func (s *Session) SetUserID(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.UserID = userID
}

// GetUserID gets the user ID for the session
func (s *Session) GetUserID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.UserID
}

// SetAdmin sets the admin flag for the session
func (s *Session) SetAdmin(isAdmin bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IsAdmin = isAdmin
}

// IsAdminSession checks if the session is an admin session
func (s *Session) IsAdminSession() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.IsAdmin
}

// Disconnect marks the session as disconnected
func (s *Session) Disconnect() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Connected = false
}

// IsConnected checks if the session is connected
func (s *Session) IsConnected() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Connected
}

// Lock locks the session for exclusive request processing
func (s *Session) Lock() {
	s.reqMu.Lock()
}

// Unlock unlocks the session
func (s *Session) Unlock() {
	s.reqMu.Unlock()
}

// SetPreferences updates the session preferences.
func (s *Session) SetPreferences(p Preferences) {
	s.mu.Lock()
	defer s.mu.Unlock()

	steps := make([]string, len(p.SkipPrioritySteps))
	copy(steps, p.SkipPrioritySteps)
	p.SkipPrioritySteps = steps

	s.preferences = &p
}

// GetPreferences returns a copy of the session preferences if available.
func (s *Session) GetPreferences() *Preferences {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.preferences == nil {
		return nil
	}

	// Return a shallow copy with duplicated slice to avoid external mutation.
	p := *s.preferences
	if len(s.preferences.SkipPrioritySteps) > 0 {
		steps := make([]string, len(s.preferences.SkipPrioritySteps))
		copy(steps, s.preferences.SkipPrioritySteps)
		p.SkipPrioritySteps = steps
	}
	return &p
}
