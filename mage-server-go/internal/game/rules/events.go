package rules

import (
	"sync"
	"time"
)

// EventType indicates the category of a rules event.
// Mirrors Java GameEvent.EventType enum.
type EventType string

const (
	// Game/Turn events
	EventBeginning                EventType = "BEGINNING"
	EventPreventDamage            EventType = "PREVENT_DAMAGE"
	EventPreventedDamage          EventType = "PREVENTED_DAMAGE"
	EventPlayTurn                 EventType = "PLAY_TURN"
	EventExtraTurn                EventType = "EXTRA_TURN"
	EventBeginTurn                EventType = "BEGIN_TURN"
	EventChangePhase              EventType = "CHANGE_PHASE"
	EventPhaseChanged             EventType = "PHASE_CHANGED"
	EventChangeStep               EventType = "CHANGE_STEP"
	EventStepChanged              EventType = "STEP_CHANGED"
	EventBeginningPhase           EventType = "BEGINNING_PHASE"
	EventBeginningPhasePre        EventType = "BEGINNING_PHASE_PRE"
	EventBeginningPhasePost       EventType = "BEGINNING_PHASE_POST"
	EventBeginningPhaseExtra      EventType = "BEGINNING_PHASE_EXTRA"
	EventBeginningPhasePreExtra   EventType = "BEGINNING_PHASE_PRE_EXTRA"
	EventBeginningPhasePostExtra  EventType = "BEGINNING_PHASE_POST_EXTRA"
	EventUntapStepPre             EventType = "UNTAP_STEP_PRE"
	EventUntapStep                EventType = "UNTAP_STEP"
	EventUntapStepPost            EventType = "UNTAP_STEP_POST"
	EventUpkeepStepPre            EventType = "UPKEEP_STEP_PRE"
	EventUpkeepStep               EventType = "UPKEEP_STEP"
	EventUpkeepStepPost           EventType = "UPKEEP_STEP_POST"
	EventDrawStepPre              EventType = "DRAW_STEP_PRE"
	EventDrawStep                 EventType = "DRAW_STEP"
	EventDrawStepPost             EventType = "DRAW_STEP_POST"
	EventPrecombatMainPhase       EventType = "PRECOMBAT_MAIN_PHASE"
	EventPrecombatMainPhasePre    EventType = "PRECOMBAT_MAIN_PHASE_PRE"
	EventPrecombatMainPhasePost   EventType = "PRECOMBAT_MAIN_PHASE_POST"
	EventPrecombatMainStepPre     EventType = "PRECOMBAT_MAIN_STEP_PRE"
	EventPrecombatMainStep        EventType = "PRECOMBAT_MAIN_STEP"
	EventPrecombatMainStepPost    EventType = "PRECOMBAT_MAIN_STEP_POST"
	EventCombatPhase              EventType = "COMBAT_PHASE"
	EventCombatPhasePre           EventType = "COMBAT_PHASE_PRE"
	EventCombatPhasePost          EventType = "COMBAT_PHASE_POST"
	EventBeginCombatStepPre       EventType = "BEGIN_COMBAT_STEP_PRE"
	EventBeginCombatStep          EventType = "BEGIN_COMBAT_STEP"
	EventBeginCombatStepPost      EventType = "BEGIN_COMBAT_STEP_POST"
	EventDeclareAttackersStepPre  EventType = "DECLARE_ATTACKERS_STEP_PRE"
	EventDeclareAttackersStep     EventType = "DECLARE_ATTACKERS_STEP"
	EventDeclareAttackersStepPost EventType = "DECLARE_ATTACKERS_STEP_POST"
	EventDeclareBlockersStepPre   EventType = "DECLARE_BLOCKERS_STEP_PRE"
	EventDeclareBlockersStep      EventType = "DECLARE_BLOCKERS_STEP"
	EventDeclareBlockersStepPost  EventType = "DECLARE_BLOCKERS_STEP_POST"
	EventCombatDamageStep         EventType = "COMBAT_DAMAGE_STEP"
	EventCombatDamageStepPre      EventType = "COMBAT_DAMAGE_STEP_PRE"
	EventCombatDamageStepPriority EventType = "COMBAT_DAMAGE_STEP_PRIORITY"
	EventCombatDamageStepPost     EventType = "COMBAT_DAMAGE_STEP_POST"
	EventEndCombatStepPre         EventType = "END_COMBAT_STEP_PRE"
	EventEndCombatStep            EventType = "END_COMBAT_STEP"
	EventEndCombatStepPost        EventType = "END_COMBAT_STEP_POST"
	EventPostcombatMainPhase      EventType = "POSTCOMBAT_MAIN_PHASE"
	EventPostcombatMainPhasePre   EventType = "POSTCOMBAT_MAIN_PHASE_PRE"
	EventPostcombatMainPhasePost  EventType = "POSTCOMBAT_MAIN_PHASE_POST"
	EventPostcombatMainStepPre    EventType = "POSTCOMBAT_MAIN_STEP_PRE"
	EventPostcombatMainStep       EventType = "POSTCOMBAT_MAIN_STEP"
	EventPostcombatMainStepPost   EventType = "POSTCOMBAT_MAIN_STEP_POST"
	EventEndPhase                 EventType = "END_PHASE"
	EventEndPhasePre              EventType = "END_PHASE_PRE"
	EventEndPhasePost             EventType = "END_PHASE_POST"
	EventEndTurnStepPre           EventType = "END_TURN_STEP_PRE"
	EventEndTurnStep              EventType = "END_TURN_STEP"
	EventEndTurnStepPost          EventType = "END_TURN_STEP_POST"
	EventCleanupStepPre           EventType = "CLEANUP_STEP_PRE"
	EventCleanupStep              EventType = "CLEANUP_STEP"
	EventCleanupStepPost          EventType = "CLEANUP_STEP_POST"
	EventEmptyManaPool            EventType = "EMPTY_MANA_POOL"
	EventAtEndOfTurn              EventType = "AT_END_OF_TURN"

	// Zone events
	EventZoneChange      EventType = "ZONE_CHANGE"
	EventZoneChangeGroup EventType = "ZONE_CHANGE_GROUP"
	EventZoneChangeBatch EventType = "ZONE_CHANGE_BATCH" // batch event

	// Card events
	EventDrawTwoOrMoreCards           EventType = "DRAW_TWO_OR_MORE_CARDS"
	EventDrawCard                     EventType = "DRAW_CARD"
	EventDrewCard                     EventType = "DREW_CARD"
	EventExplore                      EventType = "EXPLORE"
	EventExplored                     EventType = "EXPLORED"
	EventEchoPaid                     EventType = "ECHO_PAID"
	EventMiracleCardRevealed          EventType = "MIRACLE_CARD_REVEALED"
	EventMadnessCardExiled            EventType = "MADNESS_CARD_EXILED"
	EventInvestigated                 EventType = "INVESTIGATED"
	EventKicked                       EventType = "KICKED"
	EventDiscardCard                  EventType = "DISCARD_CARD"
	EventDiscardedCard                EventType = "DISCARDED_CARD"
	EventDiscardedCards               EventType = "DISCARDED_CARDS"
	EventCycleCard                    EventType = "CYCLE_CARD"
	EventCycledCard                   EventType = "CYCLED_CARD"
	EventCycleDraw                    EventType = "CYCLE_DRAW"
	EventClash                        EventType = "CLASH"
	EventClashed                      EventType = "CLASHED"
	EventMillCards                    EventType = "MILL_CARDS"
	EventMilledCard                   EventType = "MILLED_CARD"
	EventMilledCardsBatchForOnePlayer EventType = "MILLED_CARDS_BATCH_FOR_ONE_PLAYER" // batch event
	EventMilledCardsBatchForAll       EventType = "MILLED_CARDS_BATCH_FOR_ALL"        // batch event

	// Life/Damage events
	EventDamagePlayer               EventType = "DAMAGE_PLAYER"
	EventDamagedPlayer              EventType = "DAMAGED_PLAYER"
	EventDamagedBatchForPlayers     EventType = "DAMAGED_BATCH_FOR_PLAYERS"    // batch event
	EventDamagedBatchForOnePlayer   EventType = "DAMAGED_BATCH_FOR_ONE_PLAYER" // batch event
	EventDamagedBatchBySource       EventType = "DAMAGED_BATCH_BY_SOURCE"      // batch event
	EventDamagedBatchForAll         EventType = "DAMAGED_BATCH_FOR_ALL"        // batch event
	EventDamagedBatchCouldHaveFired EventType = "DAMAGED_BATCH_COULD_HAVE_FIRED"
	EventDamageCausesLifeLoss       EventType = "DAMAGE_CAUSES_LIFE_LOSS"
	EventPlayerLifeChange           EventType = "PLAYER_LIFE_CHANGE"
	EventGainLife                   EventType = "GAIN_LIFE"
	EventGainedLife                 EventType = "GAINED_LIFE"
	EventLoseLife                   EventType = "LOSE_LIFE"
	EventLostLife                   EventType = "LOST_LIFE"
	EventLostLifeBatchForOnePlayer  EventType = "LOST_LIFE_BATCH_FOR_ONE_PLAYER" // batch event
	EventLostLifeBatch              EventType = "LOST_LIFE_BATCH"                // batch event

	// Land/Spell/Ability events
	EventPlayLand                EventType = "PLAY_LAND"
	EventLandPlayed              EventType = "LAND_PLAYED"
	EventCreatureChampioned      EventType = "CREATURE_CHAMPIONED"
	EventCrewVehicle             EventType = "CREW_VEHICLE"
	EventCrewedVehicle           EventType = "CREWED_VEHICLE"
	EventVehicleCrewed           EventType = "VEHICLE_CREWED"
	EventSaddleMount             EventType = "SADDLE_MOUNT"
	EventSaddledMount            EventType = "SADDLED_MOUNT"
	EventMountSaddled            EventType = "MOUNT_SADDLED"
	EventStationPermanent        EventType = "STATION_PERMANENT"
	EventCastSpell               EventType = "CAST_SPELL"
	EventCastSpellLate           EventType = "CAST_SPELL_LATE"
	EventSpellCast               EventType = "SPELL_CAST"
	EventActivateAbility         EventType = "ACTIVATE_ABILITY"
	EventActivatedAbility        EventType = "ACTIVATED_ABILITY"
	EventTakeSpecialAction       EventType = "TAKE_SPECIAL_ACTION"
	EventTakenSpecialAction      EventType = "TAKEN_SPECIAL_ACTION"
	EventTakeSpecialManaPayment  EventType = "TAKE_SPECIAL_MANA_PAYMENT"
	EventTakenSpecialManaPayment EventType = "TAKEN_SPECIAL_MANA_PAYMENT"
	EventTriggeredAbility        EventType = "TRIGGERED_ABILITY"
	EventResolvingAbility        EventType = "RESOLVING_ABILITY"
	EventCopyStackObject         EventType = "COPY_STACKOBJECT"
	EventCopiedStackObject       EventType = "COPIED_STACKOBJECT"
	EventAddMana                 EventType = "ADD_MANA"
	EventManaAdded               EventType = "MANA_ADDED"
	EventManaPaid                EventType = "MANA_PAID"
	EventLoses                   EventType = "LOSES"
	EventLost                    EventType = "LOST"
	EventWins                    EventType = "WINS"
	EventDrawPlayer              EventType = "DRAW_PLAYER"

	// Targeting events
	EventTarget       EventType = "TARGET"
	EventTargeted     EventType = "TARGETED"
	EventTargetsValid EventType = "TARGETS_VALID"
	EventCounter      EventType = "COUNTER"
	EventCountered    EventType = "COUNTERED"

	// Combat events
	EventDeclaringAttackers  EventType = "DECLARING_ATTACKERS"
	EventDeclaredAttackers   EventType = "DECLARED_ATTACKERS"
	EventDeclareAttacker     EventType = "DECLARE_ATTACKER"
	EventAttackerDeclared    EventType = "ATTACKER_DECLARED"
	EventDefenderAttacked    EventType = "DEFENDER_ATTACKED"
	EventDeclaringBlockers   EventType = "DECLARING_BLOCKERS"
	EventDeclaredBlockers    EventType = "DECLARED_BLOCKERS"
	EventDeclareBlocker      EventType = "DECLARE_BLOCKER"
	EventBlockerDeclared     EventType = "BLOCKER_DECLARED"
	EventCreatureBlocked     EventType = "CREATURE_BLOCKED"
	EventCreatureBlocks      EventType = "CREATURE_BLOCKS"
	EventBatchBlockNoncombat EventType = "BATCH_BLOCK_NONCOMBAT"
	EventUnblockedAttacker   EventType = "UNBLOCKED_ATTACKER"
	EventCombatDamageApplied EventType = "COMBAT_DAMAGE_APPLIED"
	EventSelectedAttacker    EventType = "SELECTED_ATTACKER"
	EventSelectedBlocker     EventType = "SELECTED_BLOCKER"
	EventRemovedFromCombat   EventType = "REMOVED_FROM_COMBAT"

	// Library events
	EventSearchLibrary   EventType = "SEARCH_LIBRARY"
	EventLibrarySearched EventType = "LIBRARY_SEARCHED"
	EventShuffleLibrary  EventType = "SHUFFLE_LIBRARY"
	EventLibraryShuffled EventType = "LIBRARY_SHUFFLED"

	// Enchantment events
	EventEnchantPlayer   EventType = "ENCHANT_PLAYER"
	EventEnchantedPlayer EventType = "ENCHANTED_PLAYER"

	// Mulligan events
	EventCanTakeMulligan EventType = "CAN_TAKE_MULLIGAN"

	// Scry/Surveil events
	EventScry         EventType = "SCRY"
	EventScried       EventType = "SCRIED"
	EventScryToBottom EventType = "SCRY_TO_BOTTOM"
	EventSurveil      EventType = "SURVEIL"
	EventSurveiled    EventType = "SURVEILED"
	EventProliferate  EventType = "PROLIFERATE"
	EventProliferated EventType = "PROLIFERATED"
	EventFatesealed   EventType = "FATESEALED"

	// Random events
	EventFlipCoin         EventType = "FLIP_COIN"
	EventFlipCoins        EventType = "FLIP_COINS"
	EventCoinFlipped      EventType = "COIN_FLIPPED"
	EventReplaceRolledDie EventType = "REPLACE_ROLLED_DIE"
	EventRollDie          EventType = "ROLL_DIE"
	EventDieRolled        EventType = "DIE_ROLLED"
	EventRollDice         EventType = "ROLL_DICE"
	EventDiceRolled       EventType = "DICE_ROLLED"

	// Planeswalk events
	EventPlaneswalk   EventType = "PLANESWALK"
	EventPlaneswalked EventType = "PLANESWALKED"

	// Upkeep events
	EventPaidCumulativeUpkeep     EventType = "PAID_CUMULATIVE_UPKEEP"
	EventDidntPayCumulativeUpkeep EventType = "DIDNT_PAY_CUMULATIVE_UPKEEP"

	// Life payment events
	EventPayLife  EventType = "PAY_LIFE"
	EventLifePaid EventType = "LIFE_PAID"

	// Cascade/Learn events
	EventCascadeLand EventType = "CASCADE_LAND"
	EventLearn       EventType = "LEARN"

	// Permanent entry events
	EventEntersTheBattlefieldSelf    EventType = "ENTERS_THE_BATTLEFIELD_SELF"
	EventEntersTheBattlefieldControl EventType = "ENTERS_THE_BATTLEFIELD_CONTROL"
	EventEntersTheBattlefieldCopy    EventType = "ENTERS_THE_BATTLEFIELD_COPY"
	EventEntersTheBattlefield        EventType = "ENTERS_THE_BATTLEFIELD"
	EventPermanentEntersBattlefield  EventType = "PERMANENT_ENTERS_BATTLEFIELD" // alias for compatibility

	// Tap/Untap events
	EventTap           EventType = "TAP"
	EventTapped        EventType = "TAPPED"
	EventTappedForMana EventType = "TAPPED_FOR_MANA"
	EventTappedBatch   EventType = "TAPPED_BATCH" // batch event
	EventUntap         EventType = "UNTAP"
	EventUntapped      EventType = "UNTAPPED"
	EventUntappedBatch EventType = "UNTAPPED_BATCH" // batch event

	// Transform/Flip events
	EventFlip             EventType = "FLIP"
	EventFlipped          EventType = "FLIPPED"
	EventTransforming     EventType = "TRANSFORMING"
	EventTransformed      EventType = "TRANSFORMED"
	EventAdapt            EventType = "ADAPT"
	EventBecomesMonstrous EventType = "BECOMES_MONSTROUS"
	EventBecomesExerted   EventType = "BECOMES_EXERTED"
	EventBecomesRenowned  EventType = "BECOMES_RENOWNED"
	EventGainsClassLevel  EventType = "GAINS_CLASS_LEVEL"
	EventCreatureEnlisted EventType = "CREATURE_ENLISTED"
	EventBecomeMonarch    EventType = "BECOME_MONARCH"
	EventBecomesMonarch   EventType = "BECOMES_MONARCH"
	EventTookInitiative   EventType = "TOOK_INITIATIVE"
	EventBecomesDayNight  EventType = "BECOMES_DAY_NIGHT"
	EventMeditated        EventType = "MEDITATED"
	EventPhaseOut         EventType = "PHASE_OUT"
	EventPhasedOut        EventType = "PHASED_OUT"
	EventPhaseIn          EventType = "PHASE_IN"
	EventPhasedIn         EventType = "PHASED_IN"
	EventTurnFaceUp       EventType = "TURN_FACE_UP"
	EventTurnedFaceUp     EventType = "TURNED_FACE_UP"
	EventTurnFaceDown     EventType = "TURN_FACE_DOWN"
	EventTurnedFaceDown   EventType = "TURNED_FACE_DOWN"
	EventManifestedDread  EventType = "MANIFESTED_DREAD"
	EventOptionUsed       EventType = "OPTION_USED"

	// Damage events
	EventDamagePermanent             EventType = "DAMAGE_PERMANENT"
	EventDamagedPermanent            EventType = "DAMAGED_PERMANENT"
	EventDamagedBatchForPermanents   EventType = "DAMAGED_BATCH_FOR_PERMANENTS"    // batch event
	EventDamagedBatchForOnePermanent EventType = "DAMAGED_BATCH_FOR_ONE_PERMANENT" // batch event
	EventRemoveDamageEot             EventType = "REMOVE_DAMAGE_EOT"

	// Destruction/Sacrifice events
	EventDestroyPermanent         EventType = "DESTROY_PERMANENT"
	EventDestroyedPermanent       EventType = "DESTROYED_PERMANENT"
	EventSacrificePermanent       EventType = "SACRIFICE_PERMANENT"
	EventSacrificedPermanent      EventType = "SACRIFICED_PERMANENT"
	EventSacrificedPermanentBatch EventType = "SACRIFICED_PERMANENT_BATCH" // batch event
	EventPermanentDies            EventType = "PERMANENT_DIES"             // alias for compatibility

	// Fight/Exploit events
	EventFoughtPermanent     EventType = "FIGHTED_PERMANENT"
	EventBatchFight          EventType = "BATCH_FIGHT"
	EventExploitedCreature   EventType = "EXPLOITED_CREATURE"
	EventEvolvedCreature     EventType = "EVOLVED_CREATURE"
	EventEmbalmedCreature    EventType = "EMBALMED_CREATURE"
	EventEternalizedCreature EventType = "ETERNALIZED_CREATURE"

	// Attachment events
	EventAttach       EventType = "ATTACH"
	EventAttached     EventType = "ATTACHED"
	EventUnattach     EventType = "UNATTACH"
	EventUnattached   EventType = "UNATTACHED"
	EventStayAttached EventType = "STAY_ATTACHED"

	// Counter events
	EventAddCounter      EventType = "ADD_COUNTER"
	EventCounterAdded    EventType = "COUNTER_ADDED"
	EventAddCounters     EventType = "ADD_COUNTERS"
	EventCountersAdded   EventType = "COUNTERS_ADDED"
	EventRemoveCounter   EventType = "REMOVE_COUNTER"
	EventCounterRemoved  EventType = "COUNTER_REMOVED"
	EventRemoveCounters  EventType = "REMOVE_COUNTERS"
	EventCountersRemoved EventType = "COUNTERS_REMOVED"

	// Control events
	EventLoseControl   EventType = "LOSE_CONTROL"
	EventLostControl   EventType = "LOST_CONTROL"
	EventGainControl   EventType = "GAIN_CONTROL"
	EventGainedControl EventType = "GAINED_CONTROL"

	// Token events
	EventCreateToken   EventType = "CREATE_TOKEN"
	EventCreatedToken  EventType = "CREATED_TOKEN"
	EventCreatedTokens EventType = "CREATED_TOKENS"

	// Regeneration events
	EventRegenerate  EventType = "REGENERATE"
	EventRegenerated EventType = "REGENERATED"

	// Color change events
	EventChangeColor  EventType = "CHANGE_COLOR"
	EventColorChanged EventType = "COLOR_CHANGED"

	// Trigger events
	EventNumberOfTriggers EventType = "NUMBER_OF_TRIGGERS"

	// Voting events
	EventVote  EventType = "VOTE"
	EventVoted EventType = "VOTED"

	// Dungeon events
	EventRoomEntered      EventType = "ROOM_ENTERED"
	EventVenture          EventType = "VENTURE"
	EventVentured         EventType = "VENTURED"
	EventDungeonCompleted EventType = "DUNGEON_COMPLETED"

	// Ring/The Lord of the Rings events
	EventTemptedByRing    EventType = "TEMPTED_BY_RING"
	EventRingBearerChosen EventType = "RING_BEARER_CHOSEN"

	// Foretell events
	EventCardForetold EventType = "CARD_FORETOLD"

	// Villainous choice events
	EventFaceVillainousChoice EventType = "FACE_VILLAINOUS_CHOICE"

	// Discover events
	EventDiscovered EventType = "DISCOVERED"

	// Craft events
	EventExiledWhileCrafting EventType = "EXILED_WHILE_CRAFTING"

	// Case solving events
	EventSolveCase  EventType = "SOLVE_CASE"
	EventCaseSolved EventType = "CASE_SOLVED"

	// Suspect events
	EventBecomeSuspected EventType = "BECOME_SUSPECTED"

	// Evidence events
	EventEvidenceCollected EventType = "EVIDENCE_COLLECTED"

	// Mentor events
	EventMentoredCreature EventType = "MENTORED_CREATURE"

	// Plot events
	EventBecomePlotted EventType = "BECOME_PLOTTED"

	// Forage events
	EventForaged EventType = "FORAGED"

	// Gift events
	EventGaveGift EventType = "GAVE_GIFT"

	// Radiation events
	EventRadiationGainLife EventType = "RADIATION_GAIN_LIFE"

	// Sacrifice cost events
	EventPaySacrificeCost EventType = "PAY_SACRIFICE_COST"

	// Bending events (Avatar: The Last Airbender)
	EventEarthbended EventType = "EARTHBENDED"
	EventAirbended   EventType = "AIRBENDED"
	EventFirebended  EventType = "FIREBENDED"
	EventWaterbended EventType = "WATERBENDED"

	// Room/Door events
	EventDoorUnlocked      EventType = "DOOR_UNLOCKED"
	EventRoomFullyUnlocked EventType = "ROOM_FULLY_UNLOCKED"

	// Custom events
	EventCustomEvent EventType = "CUSTOM_EVENT"

	// Stack events (for compatibility with existing code)
	EventStackItemResolving EventType = "STACK_ITEM_RESOLVING"
	EventStackItemResolved  EventType = "STACK_ITEM_RESOLVED"
	EventStackItemRemoved   EventType = "STACK_ITEM_REMOVED"

	// State-based actions event
	EventStateBasedActions EventType = "STATE_BASED_ACTIONS"
)

// IsBatch returns true if this event type is a batch event (combines multiple sub-events).
func (et EventType) IsBatch() bool {
	batchEvents := map[EventType]bool{
		EventZoneChangeBatch:              true,
		EventMilledCardsBatchForOnePlayer: true,
		EventMilledCardsBatchForAll:       true,
		EventDamagedBatchForPlayers:       true,
		EventDamagedBatchForOnePlayer:     true,
		EventDamagedBatchBySource:         true,
		EventDamagedBatchForAll:           true,
		EventLostLifeBatchForOnePlayer:    true,
		EventLostLifeBatch:                true,
		EventTappedBatch:                  true,
		EventUntappedBatch:                true,
		EventDamagedBatchForPermanents:    true,
		EventDamagedBatchForOnePermanent:  true,
		EventSacrificedPermanentBatch:     true,
	}
	return batchEvents[et]
}

// Event represents a state change that other subsystems may react to.
// Mirrors Java GameEvent class structure.
type Event struct {
	Type           EventType
	ID             string            // Unique event ID
	TargetID       string            // ID of the target (card, player, etc.)
	SourceID       string            // ID of the source ability/object
	Controller     string            // Player ID of the controller
	PlayerID       string            // Player ID (often same as Controller, but can differ)
	Amount         int               // Numeric value (damage, life, counters, etc.)
	Flag           bool              // Boolean flag (combat damage, effect vs cost, etc.)
	Data           string            // Additional string data
	Zone           int               // Zone the event relates to (0 = none/unused)
	Targets        []string          // Multiple targets (for multi-target events)
	Timestamp      time.Time         // When the event occurred
	Metadata       map[string]string // Additional metadata
	Description    string            // Human-readable description
	AppliedEffects []string          // IDs of replacement effects already applied
}

// Listener defines a callback that reacts to incoming events.
type Listener func(Event)

// TypedListener defines a callback that reacts to a specific event type.
type TypedListener struct {
	Handle    int
	EventType EventType
	Callback  func(Event)
}

// EventBus provides a synchronous publish/subscribe implementation with type filtering.
type EventBus struct {
	mu             sync.RWMutex
	listeners      map[int]Listener              // All listeners
	typedListeners map[EventType][]TypedListener // Listeners filtered by event type
	nextHandle     int
}

// NewEventBus constructs a fresh event bus instance.
func NewEventBus() *EventBus {
	return &EventBus{
		listeners:      make(map[int]Listener),
		typedListeners: make(map[EventType][]TypedListener),
	}
}

// Subscribe registers a listener for all events and returns a handle.
func (bus *EventBus) Subscribe(listener Listener) int {
	if listener == nil {
		return -1
	}
	bus.mu.Lock()
	defer bus.mu.Unlock()
	handle := bus.nextHandle
	bus.nextHandle++
	bus.listeners[handle] = listener
	return handle
}

// SubscribeTyped registers a listener for a specific event type.
func (bus *EventBus) SubscribeTyped(eventType EventType, callback func(Event)) int {
	if callback == nil {
		return -1
	}
	bus.mu.Lock()
	defer bus.mu.Unlock()
	handle := bus.nextHandle
	bus.nextHandle++
	listener := TypedListener{
		Handle:    handle,
		EventType: eventType,
		Callback:  callback,
	}
	bus.typedListeners[eventType] = append(bus.typedListeners[eventType], listener)
	return handle
}

// Unsubscribe removes the listener identified by the provided handle.
func (bus *EventBus) Unsubscribe(handle int) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	delete(bus.listeners, handle)
	// Remove from typed listeners by handle
	for eventType, listeners := range bus.typedListeners {
		for i := len(listeners) - 1; i >= 0; i-- {
			if listeners[i].Handle == handle {
				bus.typedListeners[eventType] = append(listeners[:i], listeners[i+1:]...)
				break
			}
		}
	}
}

// UnsubscribeTyped removes a typed listener by handle.
func (bus *EventBus) UnsubscribeTyped(handle int) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	for eventType, listeners := range bus.typedListeners {
		for i := len(listeners) - 1; i >= 0; i-- {
			if listeners[i].Handle == handle {
				bus.typedListeners[eventType] = append(listeners[:i], listeners[i+1:]...)
				break
			}
		}
	}
}

// Publish delivers the event to all registered listeners synchronously.
func (bus *EventBus) Publish(event Event) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	// Publish to all listeners
	for _, listener := range bus.listeners {
		listener(event)
	}

	// Publish to typed listeners
	if typedListeners, ok := bus.typedListeners[event.Type]; ok {
		for _, listener := range typedListeners {
			listener.Callback(event)
		}
	}
}

// PublishBatch publishes multiple events, useful for batch event types.
func (bus *EventBus) PublishBatch(events []Event) {
	for _, event := range events {
		bus.Publish(event)
	}
}

// NewEvent creates a new event with common fields populated.
func NewEvent(eventType EventType, targetID, sourceID, controllerID string) Event {
	return Event{
		Type:        eventType,
		TargetID:    targetID,
		SourceID:    sourceID,
		Controller:  controllerID,
		PlayerID:    controllerID,
		Timestamp:   time.Now(),
		Metadata:    make(map[string]string),
		AppliedEffects: make([]string, 0),
	}
}

// NewEventWithAmount creates a new event with an amount value.
func NewEventWithAmount(eventType EventType, targetID, sourceID, controllerID string, amount int) Event {
	evt := NewEvent(eventType, targetID, sourceID, controllerID)
	evt.Amount = amount
	return evt
}

// NewEventWithFlag creates a new event with a flag value.
func NewEventWithFlag(eventType EventType, targetID, sourceID, controllerID string, flag bool) Event {
	evt := NewEvent(eventType, targetID, sourceID, controllerID)
	evt.Flag = flag
	return evt
}
