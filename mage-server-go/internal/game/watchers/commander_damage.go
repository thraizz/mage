package watchers

import (
	"sync"

	"github.com/magefree/mage-server-go/internal/game/rules"
)

// CommanderDamageWatcher tracks combat damage dealt by a specific commander to each player.
// This matches Java's CommanderInfoWatcher pattern.
// Per rule 903.10a: A player that's been dealt 21 or more combat damage by the same commander
// over the course of the game loses the game.
//
// Unlike other watchers that reset per turn, this watcher tracks cumulative damage
// across the entire game (WatcherScopeCard scope since each commander has its own watcher).
type CommanderDamageWatcher struct {
	*rules.BaseWatcher
	mu sync.RWMutex

	// CommanderID is the ID of the commander this watcher tracks
	CommanderID string

	// CommanderName is the name of the commander for display purposes
	CommanderName string

	// DamageToPlayer maps playerID -> total combat damage dealt by this commander
	damageToPlayer map[string]int

	// CheckCommanderDamage determines if this watcher should track damage.
	// Some formats (like Duel Commander since Nov 2016) don't use the commander damage rule.
	checkCommanderDamage bool
}

// NewCommanderDamageWatcher creates a new watcher for tracking a specific commander's damage.
// commanderID: The card ID of the commander
// commanderName: The name of the commander for display purposes
// checkDamage: Whether to track commander damage (false for formats that don't use this rule)
func NewCommanderDamageWatcher(commanderID, commanderName string, checkDamage bool) *CommanderDamageWatcher {
	w := &CommanderDamageWatcher{
		BaseWatcher:          rules.NewBaseWatcher(rules.WatcherScopeCard),
		CommanderID:          commanderID,
		CommanderName:        commanderName,
		damageToPlayer:       make(map[string]int),
		checkCommanderDamage: checkDamage,
	}
	w.SetKey("CommanderDamageWatcher_" + commanderID)
	w.SetSourceID(commanderID)
	return w
}

// Watch handles damage events and tracks commander combat damage.
func (w *CommanderDamageWatcher) Watch(event rules.Event) {
	if !w.checkCommanderDamage {
		return
	}

	// Only interested in damage events from this commander
	if event.Type != rules.EventDamagedPlayer {
		return
	}

	// Check if this damage is from our commander
	if event.SourceID != w.CommanderID {
		return
	}

	// Check if this is combat damage (stored in metadata as string)
	isCombat := event.Metadata["is_combat"] == "true" || event.Flag
	if !isCombat {
		return
	}

	// Get the player who was damaged
	playerID := event.TargetID
	if playerID == "" {
		playerID = event.PlayerID
	}
	if playerID == "" {
		return
	}

	// Get the damage amount - prefer Amount field, fall back to metadata
	amount := event.Amount
	if amount <= 0 {
		return
	}

	w.mu.Lock()
	w.damageToPlayer[playerID] += amount
	w.mu.Unlock()
	w.SetCondition(true)
}

// Reset clears the watcher state. For commander damage, this is typically only
// called when the game restarts (e.g., Karn Liberated).
func (w *CommanderDamageWatcher) Reset() {
	w.BaseWatcher.Reset()
	w.mu.Lock()
	w.damageToPlayer = make(map[string]int)
	w.mu.Unlock()
}

// GetDamageToPlayer returns the total commander damage dealt to a specific player.
func (w *CommanderDamageWatcher) GetDamageToPlayer(playerID string) int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.damageToPlayer[playerID]
}

// GetAllDamage returns a copy of all damage dealt to each player.
func (w *CommanderDamageWatcher) GetAllDamage() map[string]int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	result := make(map[string]int)
	for k, v := range w.damageToPlayer {
		result[k] = v
	}
	return result
}

// ChecksCommanderDamage returns whether this watcher is tracking damage.
func (w *CommanderDamageWatcher) ChecksCommanderDamage() bool {
	return w.checkCommanderDamage
}

// Copy creates a deep copy of this watcher.
func (w *CommanderDamageWatcher) Copy() rules.Watcher {
	w.mu.RLock()
	defer w.mu.RUnlock()

	copy := NewCommanderDamageWatcher(w.CommanderID, w.CommanderName, w.checkCommanderDamage)
	copy.SetControllerID(w.GetControllerID())
	copy.SetCondition(w.ConditionMet())

	// Deep copy damage map
	copy.damageToPlayer = make(map[string]int)
	for k, v := range w.damageToPlayer {
		copy.damageToPlayer[k] = v
	}

	return copy
}

// CommanderPlaysCountWatcher tracks how many times a commander has been cast from command zone.
// This is used for the commander tax calculation (per rule 903.8).
type CommanderPlaysCountWatcher struct {
	*rules.BaseWatcher
	mu sync.RWMutex

	// playsCount maps commanderID -> number of times cast from command zone
	playsCount map[string]int
}

// NewCommanderPlaysCountWatcher creates a watcher to track commander cast counts.
func NewCommanderPlaysCountWatcher() *CommanderPlaysCountWatcher {
	w := &CommanderPlaysCountWatcher{
		BaseWatcher: rules.NewBaseWatcher(rules.WatcherScopeGame),
		playsCount:  make(map[string]int),
	}
	w.SetKey("CommanderPlaysCountWatcher")
	return w
}

// Watch tracks when commanders are cast from the command zone.
func (w *CommanderPlaysCountWatcher) Watch(event rules.Event) {
	// We look for spell cast events where the source is a commander from command zone
	if event.Type != rules.EventSpellCast {
		return
	}

	// Check if cast from command zone (stored in metadata as string)
	fromCommandZone := event.Metadata["from_command_zone"] == "true"
	if !fromCommandZone {
		return
	}

	// Check if this is a commander (stored in metadata as string)
	isCommander := event.Metadata["is_commander"] == "true"
	if !isCommander {
		return
	}

	commanderID := event.SourceID
	if commanderID == "" {
		return
	}

	w.mu.Lock()
	w.playsCount[commanderID]++
	w.mu.Unlock()
	w.SetCondition(true)
}

// Reset clears the watcher state.
func (w *CommanderPlaysCountWatcher) Reset() {
	w.BaseWatcher.Reset()
	w.mu.Lock()
	w.playsCount = make(map[string]int)
	w.mu.Unlock()
}

// GetPlaysCount returns how many times a commander has been cast from command zone.
func (w *CommanderPlaysCountWatcher) GetPlaysCount(commanderID string) int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.playsCount[commanderID]
}

// Copy creates a deep copy of this watcher.
func (w *CommanderPlaysCountWatcher) Copy() rules.Watcher {
	w.mu.RLock()
	defer w.mu.RUnlock()

	copy := NewCommanderPlaysCountWatcher()
	copy.SetControllerID(w.GetControllerID())
	copy.SetSourceID(w.GetSourceID())
	copy.SetCondition(w.ConditionMet())

	// Deep copy plays count map
	copy.playsCount = make(map[string]int)
	for k, v := range w.playsCount {
		copy.playsCount[k] = v
	}

	return copy
}
