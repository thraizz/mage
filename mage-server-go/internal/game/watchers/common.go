package watchers

import (
	"github.com/magefree/mage-server-go/internal/game/rules"
)

// SpellsCastWatcher tracks spells cast by players.
type SpellsCastWatcher struct {
	*rules.BaseWatcher
	spellsCast map[string][]string // playerID -> list of spell IDs
}

// NewSpellsCastWatcher creates a new spells cast watcher.
func NewSpellsCastWatcher() *SpellsCastWatcher {
	w := &SpellsCastWatcher{
		BaseWatcher: rules.NewBaseWatcher(rules.WatcherScopeGame),
		spellsCast:  make(map[string][]string),
	}
	w.SetKey("SpellsCastWatcher")
	return w
}

// Watch implements the Watcher interface.
func (w *SpellsCastWatcher) Watch(event rules.Event) {
	if event.Type != rules.EventSpellCast {
		return
	}
	playerID := event.PlayerID
	if playerID == "" {
		playerID = event.Controller
	}
	if playerID == "" {
		return
	}
	spellID := event.TargetID
	if spellID == "" {
		spellID = event.SourceID
	}
	if spellID == "" {
		return
	}
	w.spellsCast[playerID] = append(w.spellsCast[playerID], spellID)
	w.SetCondition(true)
}

// Reset clears the watcher's state.
func (w *SpellsCastWatcher) Reset() {
	w.BaseWatcher.Reset()
	w.spellsCast = make(map[string][]string)
}

// GetSpellsCast returns the list of spell IDs cast by a player.
func (w *SpellsCastWatcher) GetSpellsCast(playerID string) []string {
	return w.spellsCast[playerID]
}

// GetCount returns the number of spells cast by a player.
func (w *SpellsCastWatcher) GetCount(playerID string) int {
	return len(w.spellsCast[playerID])
}

// Copy creates a copy of this watcher.
func (w *SpellsCastWatcher) Copy() rules.Watcher {
	copy := NewSpellsCastWatcher()
	copy.SetControllerID(w.GetControllerID())
	copy.SetSourceID(w.GetSourceID())
	copy.SetCondition(w.ConditionMet())
	// Deep copy spells cast map
	copy.spellsCast = make(map[string][]string)
	for k, v := range w.spellsCast {
		copy.spellsCast[k] = append([]string(nil), v...)
	}
	return copy
}

// CreaturesDiedWatcher tracks creatures that died (went to graveyard from battlefield).
type CreaturesDiedWatcher struct {
	*rules.BaseWatcher
	creaturesDiedByController map[string]int // controllerID -> count
	creaturesDiedByOwner      map[string]int // ownerID -> count
}

// NewCreaturesDiedWatcher creates a new creatures died watcher.
func NewCreaturesDiedWatcher() *CreaturesDiedWatcher {
	w := &CreaturesDiedWatcher{
		BaseWatcher:              rules.NewBaseWatcher(rules.WatcherScopeGame),
		creaturesDiedByController: make(map[string]int),
		creaturesDiedByOwner:      make(map[string]int),
	}
	w.SetKey("CreaturesDiedWatcher")
	return w
}

// Watch implements the Watcher interface.
func (w *CreaturesDiedWatcher) Watch(event rules.Event) {
	if event.Type != rules.EventPermanentDies {
		return
	}
	// Check if it's a creature (would need to check card type from metadata)
	// For now, assume all permanent dies events are creatures
	controllerID := event.Controller
	ownerID := event.Metadata["owner_id"]
	if ownerID == "" {
		ownerID = controllerID
	}
	if controllerID != "" {
		w.creaturesDiedByController[controllerID]++
	}
	if ownerID != "" {
		w.creaturesDiedByOwner[ownerID]++
	}
	w.SetCondition(true)
}

// Reset clears the watcher's state.
func (w *CreaturesDiedWatcher) Reset() {
	w.BaseWatcher.Reset()
	w.creaturesDiedByController = make(map[string]int)
	w.creaturesDiedByOwner = make(map[string]int)
}

// GetAmountByController returns the number of creatures that died for a controller.
func (w *CreaturesDiedWatcher) GetAmountByController(controllerID string) int {
	return w.creaturesDiedByController[controllerID]
}

// GetAmountByOwner returns the number of creatures that died for an owner.
func (w *CreaturesDiedWatcher) GetAmountByOwner(ownerID string) int {
	return w.creaturesDiedByOwner[ownerID]
}

// GetTotalAmount returns the total number of creatures that died.
func (w *CreaturesDiedWatcher) GetTotalAmount() int {
	total := 0
	for _, count := range w.creaturesDiedByController {
		total += count
	}
	return total
}

// Copy creates a copy of this watcher.
func (w *CreaturesDiedWatcher) Copy() rules.Watcher {
	copy := NewCreaturesDiedWatcher()
	copy.SetControllerID(w.GetControllerID())
	copy.SetSourceID(w.GetSourceID())
	copy.SetCondition(w.ConditionMet())
	// Deep copy maps
	copy.creaturesDiedByController = make(map[string]int)
	for k, v := range w.creaturesDiedByController {
		copy.creaturesDiedByController[k] = v
	}
	copy.creaturesDiedByOwner = make(map[string]int)
	for k, v := range w.creaturesDiedByOwner {
		copy.creaturesDiedByOwner[k] = v
	}
	return copy
}

// CardsDrawnWatcher tracks cards drawn by players.
type CardsDrawnWatcher struct {
	*rules.BaseWatcher
	cardsDrawn map[string]int // playerID -> count
}

// NewCardsDrawnWatcher creates a new cards drawn watcher.
func NewCardsDrawnWatcher() *CardsDrawnWatcher {
	w := &CardsDrawnWatcher{
		BaseWatcher: rules.NewBaseWatcher(rules.WatcherScopeGame),
		cardsDrawn:  make(map[string]int),
	}
	w.SetKey("CardsDrawnWatcher")
	return w
}

// Watch implements the Watcher interface.
func (w *CardsDrawnWatcher) Watch(event rules.Event) {
	if event.Type != rules.EventDrewCard {
		return
	}
	playerID := event.PlayerID
	if playerID == "" {
		playerID = event.Controller
	}
	if playerID == "" {
		return
	}
	w.cardsDrawn[playerID]++
	w.SetCondition(true)
}

// Reset clears the watcher's state.
func (w *CardsDrawnWatcher) Reset() {
	w.BaseWatcher.Reset()
	w.cardsDrawn = make(map[string]int)
}

// GetCount returns the number of cards drawn by a player.
func (w *CardsDrawnWatcher) GetCount(playerID string) int {
	return w.cardsDrawn[playerID]
}

// Copy creates a copy of this watcher.
func (w *CardsDrawnWatcher) Copy() rules.Watcher {
	copy := NewCardsDrawnWatcher()
	copy.SetControllerID(w.GetControllerID())
	copy.SetSourceID(w.GetSourceID())
	copy.SetCondition(w.ConditionMet())
	// Deep copy map
	copy.cardsDrawn = make(map[string]int)
	for k, v := range w.cardsDrawn {
		copy.cardsDrawn[k] = v
	}
	return copy
}

// PermanentsEnteredWatcher tracks permanents that entered the battlefield.
type PermanentsEnteredWatcher struct {
	*rules.BaseWatcher
	permanentsEntered map[string][]string // controllerID -> list of permanent IDs
}

// NewPermanentsEnteredWatcher creates a new permanents entered watcher.
func NewPermanentsEnteredWatcher() *PermanentsEnteredWatcher {
	w := &PermanentsEnteredWatcher{
		BaseWatcher:       rules.NewBaseWatcher(rules.WatcherScopeGame),
		permanentsEntered: make(map[string][]string),
	}
	w.SetKey("PermanentsEnteredWatcher")
	return w
}

// Watch implements the Watcher interface.
func (w *PermanentsEnteredWatcher) Watch(event rules.Event) {
	if event.Type != rules.EventEntersTheBattlefield {
		return
	}
	controllerID := event.Controller
	if controllerID == "" {
		return
	}
	permanentID := event.TargetID
	if permanentID == "" {
		permanentID = event.SourceID
	}
	if permanentID == "" {
		return
	}
	w.permanentsEntered[controllerID] = append(w.permanentsEntered[controllerID], permanentID)
	w.SetCondition(true)
}

// Reset clears the watcher's state.
func (w *PermanentsEnteredWatcher) Reset() {
	w.BaseWatcher.Reset()
	w.permanentsEntered = make(map[string][]string)
}

// GetPermanentsEntered returns the list of permanent IDs that entered for a controller.
func (w *PermanentsEnteredWatcher) GetPermanentsEntered(controllerID string) []string {
	return w.permanentsEntered[controllerID]
}

// Copy creates a copy of this watcher.
func (w *PermanentsEnteredWatcher) Copy() rules.Watcher {
	copy := NewPermanentsEnteredWatcher()
	copy.SetControllerID(w.GetControllerID())
	copy.SetSourceID(w.GetSourceID())
	copy.SetCondition(w.ConditionMet())
	// Deep copy map
	copy.permanentsEntered = make(map[string][]string)
	for k, v := range w.permanentsEntered {
		copy.permanentsEntered[k] = append([]string(nil), v...)
	}
	return copy
}
