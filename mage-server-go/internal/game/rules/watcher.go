package rules

import (
	"sync"
)

// WatcherScope defines the scope of a watcher's tracking.
type WatcherScope int

const (
	// WatcherScopeGame tracks events for the entire game.
	WatcherScopeGame WatcherScope = iota
	// WatcherScopePlayer tracks events for a specific player.
	WatcherScopePlayer
	// WatcherScopeCard tracks events for a specific card/permanent.
	WatcherScopeCard
)

// String returns the string representation of the watcher scope.
func (ws WatcherScope) String() string {
	switch ws {
	case WatcherScopeGame:
		return "GAME"
	case WatcherScopePlayer:
		return "PLAYER"
	case WatcherScopeCard:
		return "CARD"
	default:
		return "UNKNOWN"
	}
}

// Watcher is an interface for objects that watch game events and track conditions.
// Watchers are used to implement conditional abilities and track game state.
type Watcher interface {
	// Watch is called when an event occurs that this watcher is interested in.
	Watch(event Event)

	// Reset clears the watcher's condition and state (typically called at end of turn/phase).
	Reset()

	// ConditionMet returns true if the condition this watcher tracks has been met.
	ConditionMet() bool

	// GetScope returns the scope of this watcher.
	GetScope() WatcherScope

	// GetKey returns a unique key for this watcher instance.
	// For GAME scope: returns class name
	// For PLAYER scope: returns playerID + class name
	// For CARD scope: returns cardID + class name
	GetKey() string

	// Copy creates a deep copy of this watcher.
	Copy() Watcher
}

// BaseWatcher provides a base implementation for watchers.
type BaseWatcher struct {
	scope        WatcherScope
	controllerID string
	sourceID     string
	condition    bool
	key          string
}

// NewBaseWatcher creates a new base watcher with the specified scope.
func NewBaseWatcher(scope WatcherScope) *BaseWatcher {
	return &BaseWatcher{
		scope:     scope,
		condition: false,
	}
}

// GetScope returns the watcher's scope.
func (bw *BaseWatcher) GetScope() WatcherScope {
	return bw.scope
}

// SetControllerID sets the controller ID (for PLAYER scope watchers).
func (bw *BaseWatcher) SetControllerID(id string) {
	bw.controllerID = id
}

// GetControllerID returns the controller ID.
func (bw *BaseWatcher) GetControllerID() string {
	return bw.controllerID
}

// SetSourceID sets the source ID (for CARD scope watchers).
func (bw *BaseWatcher) SetSourceID(id string) {
	bw.sourceID = id
}

// GetSourceID returns the source ID.
func (bw *BaseWatcher) GetSourceID() string {
	return bw.sourceID
}

// ConditionMet returns whether the condition has been met.
func (bw *BaseWatcher) ConditionMet() bool {
	return bw.condition
}

// SetCondition sets the condition flag.
func (bw *BaseWatcher) SetCondition(condition bool) {
	bw.condition = condition
}

// Reset clears the condition.
func (bw *BaseWatcher) Reset() {
	bw.condition = false
}

// GetKey returns the unique key for this watcher.
func (bw *BaseWatcher) GetKey() string {
	return bw.key
}

// SetKey sets the unique key for this watcher.
func (bw *BaseWatcher) SetKey(key string) {
	bw.key = key
}

// WatcherRegistry manages watchers for a game.
type WatcherRegistry struct {
	mu       sync.RWMutex
	watchers map[string]Watcher // key -> watcher
	byScope  map[WatcherScope][]Watcher
}

// NewWatcherRegistry creates a new watcher registry.
func NewWatcherRegistry() *WatcherRegistry {
	return &WatcherRegistry{
		watchers: make(map[string]Watcher),
		byScope:  make(map[WatcherScope][]Watcher),
	}
}

// AddWatcher adds a watcher to the registry.
func (wr *WatcherRegistry) AddWatcher(watcher Watcher) {
	if watcher == nil {
		return
	}

	wr.mu.Lock()
	defer wr.mu.Unlock()

	key := watcher.GetKey()
	if key == "" {
		// Generate key based on scope
		key = wr.generateKey(watcher)
		// Try to set the key if the watcher has a SetKey method
		// This works for watchers that embed BaseWatcher
		if setter, ok := watcher.(interface{ SetKey(string) }); ok {
			setter.SetKey(key)
		}
	}

	wr.watchers[key] = watcher
	scope := watcher.GetScope()
	wr.byScope[scope] = append(wr.byScope[scope], watcher)
}

// RemoveWatcher removes a watcher from the registry.
func (wr *WatcherRegistry) RemoveWatcher(key string) {
	wr.mu.Lock()
	defer wr.mu.Unlock()

	watcher, ok := wr.watchers[key]
	if !ok {
		return
	}

	delete(wr.watchers, key)

	// Remove from byScope
	scope := watcher.GetScope()
	watchers := wr.byScope[scope]
	for i, w := range watchers {
		if w.GetKey() == key {
			wr.byScope[scope] = append(watchers[:i], watchers[i+1:]...)
			break
		}
	}
}

// GetWatcher retrieves a watcher by key.
func (wr *WatcherRegistry) GetWatcher(key string) Watcher {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	return wr.watchers[key]
}

// GetWatchersByScope returns all watchers for a given scope.
func (wr *WatcherRegistry) GetWatchersByScope(scope WatcherScope) []Watcher {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	watchers := wr.byScope[scope]
	result := make([]Watcher, len(watchers))
	copy(result, watchers)
	return result
}

// GetAllWatchers returns all registered watchers.
func (wr *WatcherRegistry) GetAllWatchers() []Watcher {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	result := make([]Watcher, 0, len(wr.watchers))
	for _, watcher := range wr.watchers {
		result = append(result, watcher)
	}
	return result
}

// ResetWatchers resets all watchers (typically called at end of turn/phase).
func (wr *WatcherRegistry) ResetWatchers() {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	for _, watcher := range wr.watchers {
		watcher.Reset()
	}
}

// ResetWatchersByScope resets all watchers for a given scope.
func (wr *WatcherRegistry) ResetWatchersByScope(scope WatcherScope) {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	for _, watcher := range wr.byScope[scope] {
		watcher.Reset()
	}
}

// generateKey generates a unique key for a watcher based on its scope.
// Uses reflection to get the type name for better uniqueness.
func (wr *WatcherRegistry) generateKey(watcher Watcher) string {
	scope := watcher.GetScope()
	
	// Use reflection to get type name
	typeName := getWatcherTypeName(watcher)
	
	switch scope {
	case WatcherScopeGame:
		return typeName
	case WatcherScopePlayer:
		if getter, ok := watcher.(interface{ GetControllerID() string }); ok {
			if controllerID := getter.GetControllerID(); controllerID != "" {
				return controllerID + "_" + typeName
			}
		}
		return typeName
	case WatcherScopeCard:
		if getter, ok := watcher.(interface{ GetSourceID() string }); ok {
			if sourceID := getter.GetSourceID(); sourceID != "" {
				return sourceID + "_" + typeName
			}
		}
		return typeName
	default:
		return typeName
	}
}

// getWatcherTypeName extracts the type name from a watcher.
// This is a simplified version - in practice you might use reflection.
func getWatcherTypeName(watcher Watcher) string {
	// Try to get a meaningful name from common watcher types
	switch watcher.(type) {
	case interface{ GetKey() string }:
		// If watcher has GetKey, use it (but avoid recursion)
		return "Watcher"
	default:
		// Use reflection would be better, but for now use a simple approach
		return "Watcher"
	}
}

// NotifyWatchers notifies all relevant watchers of an event.
func (wr *WatcherRegistry) NotifyWatchers(event Event) {
	wr.mu.RLock()
	defer wr.mu.RUnlock()

	// Notify all watchers (they can filter internally)
	for _, watcher := range wr.watchers {
		watcher.Watch(event)
	}
}
