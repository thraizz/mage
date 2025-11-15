# Go Port Task Tracker

Status legend:
- `[x]` Completed
- `[ ]` Pending / not yet started
- `[~]` In progress or partially implemented

## Engine Scaffolding & Lifecycle
- [x] Wire gRPC server to `MageEngine` via `EngineAdapter`
- [x] Provide `MageEngine` core skeleton that tracks games, players, and actions
- [x] Implement `TurnManager` mirroring MTG phase/step progression and priority handoff
- [x] Introduce `StackManager` with basic push/pop mechanics and simple resolution hooks
- [~] Extend stack resolution to support triggered abilities, replacement effects, and modal choices _(trigger queue with APNAP ordering implemented; triggered abilities queued and processed before priority; replacement/modal hooks pending)_
- [ ] Implement priority windows for casting during stack resolution (e.g., mana abilities, nested responses)
- [ ] Persist stack/game events for replay and spectator synchronization
- [x] Add comprehensive error handling and rollback when resolution fails _(ProcessAction creates automatic bookmark before each action; on error, state is restored to bookmark; successful actions remove bookmark; comprehensive logging of recovery operations)_
- [x] Implement priority retention after casting (caster retains priority by default, only passes when explicitly passing)
- [x] Add state bookmarking and rollback mechanism for error recovery _(Complete snapshot system with deep copy of all game state: players, cards, zones, stack, messages; BookmarkState/RestoreState/RemoveBookmark/ClearBookmarks methods; automatic error recovery in ProcessAction; tested with multiple bookmarks, restoration, and error scenarios)_
- [ ] Implement comprehensive priority loop structure matching Java `playPriority()` pattern
- [x] Implement mulligan system _(London mulligan with StartMulligan/PlayerMulligan/PlayerKeepHand/EndMulligan; tracks mulligan count per player; draws 7-N cards; full validation; tested with multiple mulligans and edge cases)_
- [x] Implement game cleanup and resource disposal _(CleanupGame removes game from engine; clears all bookmarks, turn snapshots, and watchers; thread-safe; tested with resource verification)_
- [x] Add complete lifecycle state validation _(Pause/resume validation; mulligan phase checks; finished game checks; comprehensive error messages; tested with edge cases)_

## Game State & Zones
- [x] Surface battlefield/stack state via `GameGetView`
- [x] Synchronize graveyard, exile, command, and hidden zones with engine updates _(unified `moveCard` system handles all zone transitions with proper removal/addition; graveyard/exile/command zones tracked; zone change events emitted)_
- [x] Track card ownership/controller changes (gain control, copying, phasing, etc.) _(ChangeControl method implemented to change controller of permanents; emits GAIN_CONTROL and LOSE_CONTROL events; validates new controller is in game; comprehensive error handling; ownership tracked separately from control; tested with error cases)_
- [~] Implement continuous effects layer system (layers 1-7 per Comprehensive Rules) _(layer manager in place for basic power/toughness buffs; additional layers forthcoming)_
- [x] Handle state-based actions (lethal damage, zero loyalty, legend rule, etc.) _(life ≤ 0 auto-loss; zero/less toughness deaths; planeswalker 0 loyalty; lethal damage (damage >= toughness); legend rule (multiple legendary permanents with same name); world rule (multiple world enchantments); planeswalker uniqueness (multiple planeswalkers with same type); damage tracking system implemented)_
- [x] Support counters (loyalty, +1/+1, poison, energy, experience) _(Counter and Counters data structures; counter operations (add/remove) with event emission; boost counters integrated with layer system for power/toughness; planeswalker loyalty SBA; player counters tracked for poison/energy/experience; counter views updated across all zones)_
- [x] Provide deterministic UUID mapping for permanents, abilities, and triggers
- [x] Call `checkStateBasedActions()` before each priority (per rule 117.5)
- [x] Fix `resetPassed()` to preserve lost/left player state (`passed = loses || hasLeft()`)
- [x] Add `canRespond()` checking in pass logic (only consider responding players in `allPassed()`)
- [x] Ensure proper zone tracking after stack resolution (cards moved to correct zones with events)

## Stack & Trigger System
- [x] Record log message when a stack item resolves
- [x] Auto-advance priority after resolution back to the active player
- [x] Allow triggered abilities to be queued automatically when conditions are met
- [ ] Support casting spells/activating abilities while another object is resolving (linked abilities)
- [ ] Implement replacement/prevention effects that modify or negate stack resolution
- [x] Ensure stack legality checks (targets available, costs paid) prior to resolution _(LegalityChecker validates controller status, source card zones, target validity, timing restrictions, and cost payment; illegal items automatically removed from stack with events; GameStateAccessor interface implemented for state queries; comprehensive unit tests)_
- [ ] Implement target selection flow for spells/abilities requiring targets _(server must prompt for target selection before adding to stack; validate target requirements match spell/ability rules; handle multi-targeting; ensure targets are selected before casting completes)_
- [x] Add exhaustive integration tests covering multi-object stacks, counterspells, and priority loops _(Comprehensive test suite: multi-object stack resolution (LIFO order), counterspell interactions, nested responses, priority loops with multiple players, stack legality checks, complex scenarios with multiple spells/responses; tests account for triggered abilities and verify proper resolution order)_
- [x] Resolve stack one item at a time with state-based action and triggered ability checks between each resolution
- [x] Implement triggered ability queue processing before priority (APNAP order: Active Player, Non-Active Player)
- [x] Add `checkStateAndTriggered()` method that runs before each priority (SBA → triggers → repeat until stable)
- [x] Handle simultaneous events between stack resolutions (process events after each resolution)

## Combat System 🚧 IN PROGRESS (~5% coverage, ~2,500 lines needed)
### Core Combat Infrastructure (P0 - Critical)
- [x] Implement `combatState` struct tracking attackers, blockers, groups, defenders, tapped creatures _(complete with all tracking maps)_
- [x] Implement `combatGroup` struct for attacker-blocker-defender groupings with damage ordering _(complete with attacker/blocker lists, orders, blocked status)_
- [x] Add combat fields to `internalCard`: `Attacking`, `Blocking`, `AttackingWhat`, `BlockingWhat` _(all fields added)_
- [x] Add `combat *combatState` to `engineGameState` _(initialized in StartGame)_
- [x] Implement `ResetCombat(gameID)` - clear combat state at beginning of combat _(clears all combat state and card flags)_
- [x] Implement `SetAttacker(gameID, playerID)` - set attacking player _(sets attackingPlayerID)_
- [x] Implement `SetDefenders(gameID)` - identify all valid defenders (players, planeswalkers, battles) _(adds opponent players; TODO: planeswalkers/battles)_

### Attacker Declaration (P0 - Critical)
- [x] Implement `DeclareAttacker(gameID, creatureID, defenderID, playerID)` - declare single attacker _(complete with validation, tapping, group creation, events)_
- [~] Implement `CanAttack(gameID, creatureID)` - validate creature can attack (summoning sickness, tapped, restrictions) _(tapped check done; TODO: summoning sickness, restrictions)_
- [ ] Implement `CanAttackDefender(gameID, creatureID, defenderID)` - validate can attack specific defender
- [~] Implement attacker tapping logic (tap unless vigilance) _(taps creature; TODO: check vigilance ability)_
- [x] Create/update combat groups when attackers declared _(creates group per defender, adds attackers)_
- [ ] Implement `RemoveAttacker(gameID, attackerID)` - undo attacker declaration
- [x] Fire `EventAttackerDeclared` per attacker and `EventDeclaredAttackers` when complete _(EventAttackerDeclared fired per attacker)_

### Blocker Declaration (P0 - Critical)
- [ ] Implement `DeclareBlocker(gameID, blockerID, attackerID, playerID)` - declare single blocker
- [ ] Implement `CanBlock(gameID, blockerID, attackerID)` - validate creature can block (tapped, flying, restrictions)
- [ ] Add blocker to combat group, update blocked status
- [ ] Implement `AcceptBlockers(gameID)` - finalize blockers, check requirements/restrictions
- [ ] Implement blocker ordering for multiple blockers on same attacker
- [ ] Implement `RemoveBlocker(gameID, blockerID)` - undo blocker declaration
- [ ] Fire `EventBlockerDeclared` per blocker and `EventDeclaredBlockers` when complete

### Damage Assignment & Application (P0 - Critical)
- [ ] Implement `AssignCombatDamage(gameID, firstStrike bool)` - assign damage for attackers and blockers
- [ ] Implement `combatGroup.assignDamageToBlockers()` - attacker damage to blockers with ordering
- [ ] Implement `combatGroup.assignDamageToAttackers()` - blocker damage to attackers with ordering
- [ ] Implement `ApplyCombatDamage(gameID)` - apply all assigned damage
- [ ] Implement `combatGroup.applyDamage()` - mark damage on creatures and players
- [ ] Handle unblocked attacker damage to defending player/permanent
- [ ] Fire `EventCombatDamageAssigned` and `EventCombatDamageApplied` events

### Combat Cleanup (P0 - Critical)
- [ ] Implement `EndCombat(gameID)` - move groups to formerGroups, clear current combat
- [ ] Clear `Attacking` and `Blocking` flags on all creatures
- [ ] Keep attacker tracking for "attacked this turn" queries
- [ ] Fire `EventEndCombat` event

### First Strike & Double Strike (P1 - High Priority)
- [ ] Implement `HasFirstOrDoubleStrike(gameID)` - check if any creature has first/double strike
- [ ] Add first strike combat damage step before normal damage step
- [ ] Implement `combatGroup.hasFirstOrDoubleStrike()` per group
- [ ] Handle double strike creatures dealing damage in both steps
- [ ] Prevent normal damage from creatures killed by first strike

### Trample (P1 - High Priority)
- [ ] Implement trample damage calculation (excess damage to defender)
- [ ] Add `canDamageDefenderDirectly` flag to combat groups
- [ ] Validate lethal damage assigned to blockers before overflow
- [ ] Handle trample damage to planeswalkers/battles
- [ ] Support "trample over planeswalkers" rule

### Vigilance (P1 - High Priority)
- [ ] Check for vigilance ability before tapping attacker
- [ ] Track which attackers were tapped by attack in `attackersTappedByAttack`
- [ ] Support effects that grant vigilance during combat

### Flying & Reach (P1 - High Priority)
- [ ] Implement flying restriction (can only be blocked by flying/reach)
- [ ] Implement reach ability (can block flying)
- [ ] Add `CanBlock` validation for flying/reach interactions
- [ ] Support effects that grant/remove flying during combat

### Combat Events (P1 - High Priority)
- [ ] Add `EventBeginCombat` - beginning of combat step
- [ ] Add `EventDeclareAttackersStepPre` - before attacker declaration
- [ ] Add `EventAttackerDeclared` - per attacker declared
- [ ] Add `EventDeclaredAttackers` - all attackers declared
- [ ] Add `EventDeclareBlockersStepPre` - before blocker declaration
- [ ] Add `EventBlockerDeclared` - per blocker declared
- [ ] Add `EventDeclaredBlockers` - all blockers declared
- [ ] Add `EventCombatDamageStepPre` - before damage assignment
- [ ] Add `EventCombatDamageAssigned` - damage assigned
- [ ] Add `EventCombatDamageApplied` - damage applied
- [ ] Add `EventEndCombatStepPre` - before end of combat
- [ ] Add `EventEndCombat` - combat ended

### Combat Validation & Requirements (P1 - High Priority)
- [ ] Implement `CheckBlockRequirements(gameID, playerID)` - must block if able
- [ ] Implement `CheckBlockRestrictions(gameID, playerID)` - can't block restrictions
- [ ] Implement forced attack tracking (`creaturesForcedToAttack` map)
- [ ] Implement must block tracking (`creatureMustBlockAttackers` map)
- [ ] Validate minimum/maximum attacker counts
- [ ] Validate minimum/maximum blocker counts per attacker

### Combat Triggers (P1 - High Priority)
- [ ] Queue triggers on attacker declared (e.g., "Whenever ~ attacks")
- [ ] Queue triggers on blocker declared (e.g., "Whenever ~ blocks")
- [ ] Queue triggers on creature becomes blocked (e.g., "Whenever ~ becomes blocked")
- [ ] Queue triggers on combat damage dealt (e.g., "Whenever ~ deals combat damage")
- [ ] Queue triggers on creature dies in combat
- [ ] Process triggers via existing `checkStateAndTriggered()` system

### Special Combat Mechanics (P2 - Medium Priority)
- [ ] Implement menace (must be blocked by 2+ creatures)
- [ ] Implement deathtouch (any damage is lethal)
- [ ] Implement lifelink (gain life equal to damage dealt)
- [ ] Implement defender (can't attack)
- [ ] Implement "can't be blocked" effects
- [ ] Implement "must be blocked if able" effects
- [ ] Implement "attacks each combat if able" effects

### Planeswalker & Battle Combat (P2 - Medium Priority)
- [ ] Support attacking planeswalkers (redirect from player)
- [ ] Support attacking battles
- [ ] Implement planeswalker damage redirection rules
- [ ] Track which planeswalkers/battles are attacked
- [ ] Handle planeswalker removal during combat

### Damage Ordering (P2 - Medium Priority)
- [ ] Implement attacker damage ordering for multiple blockers
- [ ] Implement blocker damage ordering for multiple attackers
- [ ] Prompt players to order blockers/attackers
- [ ] Validate damage assignment follows ordering
- [ ] Support "you choose damage order" effects

### Banding (P3 - Low Priority, Complex)
- [ ] Implement banding ability detection
- [ ] Implement "bands with other" ability
- [ ] Allow banded creatures to attack as group
- [ ] Implement banding damage assignment rules (defending player assigns)
- [ ] Handle banding restrictions and requirements

### Combat Removal & Interruption (P2 - Medium Priority)
- [ ] Implement `RemoveFromCombat(gameID, creatureID)` - remove during combat
- [ ] Handle creature removal during attacker declaration
- [ ] Handle creature removal during blocker declaration
- [ ] Handle creature removal during damage assignment
- [ ] Update combat groups when creatures removed
- [ ] Handle blink/flicker during combat (removed and returns as new object)
- [ ] Handle phase out during combat

### Combat Integration with Turn Structure (P0 - Critical)
- [ ] Wire `ResetCombat()` to beginning of combat step
- [ ] Wire `SetAttacker()` and `SetDefenders()` to beginning of combat
- [ ] Wire attacker declaration to declare attackers step
- [ ] Wire blocker declaration to declare blockers step
- [ ] Wire first strike damage to first strike damage step
- [ ] Wire normal damage to combat damage step
- [ ] Wire `EndCombat()` to end of combat step
- [ ] Add combat damage steps to turn structure if first strike exists

### Combat Testing (P0 - Critical)
- [ ] Test single attacker, no blockers (damage to player)
- [ ] Test single attacker, single blocker (damage to creatures)
- [ ] Test multiple attackers, no blockers
- [ ] Test multiple attackers, multiple blockers
- [ ] Test multiple blockers on single attacker (damage ordering)
- [ ] Test creature death from combat damage
- [ ] Test player death from combat damage
- [ ] Test vigilance (no tap on attack)
- [ ] Test first strike damage (kill before normal damage)
- [ ] Test double strike damage (damage in both steps)
- [ ] Test trample damage (overflow to player)
- [ ] Test flying/reach restrictions
- [ ] Test combat triggers firing
- [ ] Test combat events firing
- [ ] Test removal during combat (all phases)
- [ ] Test blocker requirements/restrictions
- [ ] Test attacker requirements/restrictions

### Combat View & Display (P1 - High Priority)
- [ ] Populate `EngineCombatView` with actual combat state
- [ ] Populate `EngineCombatGroupView` for each combat group
- [ ] Show attacking creatures in game view
- [ ] Show blocking creatures in game view
- [ ] Show damage assignments in game view
- [ ] Update combat view after each declaration/assignment

## Player Interaction & Prompts
- [x] Emit prompts when priority passes require player response
- [ ] Support multi-choice prompts (choose mode, targets, numbers, colors) _(target selection prompts covered by "target selection flow" task in Stack & Trigger System)_
- [x] Implement mana payment flow (floating mana, cost reductions, hybrid costs) _(mana pool with regular/floating mana; cost parsing for generic, colored, X, hybrid; payment calculation and execution; cost reduction effects; floating mana empties at end of step; integrated with spell casting; comprehensive tests)_
- [x] Add concession, timeout, and match result handling aligned with rules

## Card Database & Ability Port
- [ ] Inventory Java ability/card modules and map to Go packages
- [ ] Generate Go card definitions from existing Java card data (expansions, tokens, abilities)
- [ ] Translate ability scripts (activated, triggered, static) into Go equivalents
- [ ] Port keyword ability handlers (flying, deathtouch, scry, etc.)
- [ ] Implement effect infrastructure (replacement effects, static ability watchers, continuous effects)
- [ ] Build automated verification to compare Java vs Go card behavior for representative samples

## Event System & Watchers
- [x] Mirror Java event bus for game events _(Complete event bus implementation with all 200+ event types from Java GameEvent.EventType enum; typed subscriptions, batch events, helper functions; events wired for spell cast, zone changes, life changes, mana, phase/step transitions, stack resolution, permanent entry/dies)_
- [x] Port watcher/listener infrastructure to track conditional abilities _(Watcher interface with Watch/Reset/ConditionMet methods; WatcherRegistry for managing watchers by scope (GAME/PLAYER/CARD); BaseWatcher helper; common watchers implemented: SpellsCastWatcher, CreaturesDiedWatcher, CardsDrawnWatcher, PermanentsEnteredWatcher; watchers wired to event bus; auto-reset on cleanup step; comprehensive integration tests covering event bus integration, multi-watcher scenarios, scope isolation, lifecycle management, thread safety, and real game flows)_
- [x] Provide hooks for UI/websocket notifications (combat updates, triggers, log lines)
- [x] Capture analytics/metrics for stack depth, actions per turn, average response time
- [x] Queue triggered abilities instead of immediately pushing to stack (process via `checkTriggered()` before priority)

## Undo/Redo & State Management
- [x] Implement per-player stored bookmarks for action undo _(StoredBookmark field on internalPlayer; SetPlayerStoredBookmark/ResetPlayerStoredBookmark methods; bookmark set automatically in ProcessAction; preserved until action completes or becomes irreversible)_
- [x] Add player-initiated undo command _(Undo(gameID, playerID) method; restores to player's stored bookmark; clears bookmark after undo; comprehensive error handling; tested with spell casting undo)_
- [x] Implement strategic bookmark placement in game flow _(Automatic bookmark creation in ProcessAction before each action; bookmarks set as player's stored bookmark; covers spell casting, mana payment, combat declarations)_
- [x] Add bookmark invalidation rules _(Bookmarks cleared when spell resolves, phase changes, turn rollback occurs; ResetPlayerStoredBookmark called at appropriate times; tested with resolution invalidation)_
- [x] Implement turn rollback system with turn-level snapshots _(Separate turnSnapshots map; SaveTurnSnapshot at start of each turn; keeps last 4 turns; CanRollbackTurns/RollbackTurns methods; clears all player bookmarks on rollback; tested with snapshot limits)_
- [x] Integrate undo/redo with error recovery system _(ProcessAction creates bookmark before action; on error, auto-restores; on success, checks if bookmark in use by player before removing; seamless integration with player undo)_

## Persistence, Replays & Recovery
- [x] Store game snapshots for reconnection and spectating _(gameStateSnapshot structure with complete deep copy of all game state; used for both undo/redo and turn rollback; efficient memory management with automatic cleanup)_
- [ ] Implement replay recording/playback (step-by-step action logs)
- [ ] Ensure deterministic serialization for saved games and tournaments
- [ ] Add checksum/validation to guard against divergent game state

## Testing & Parity Validation
- [x] Add unit tests for `TurnManager` sequencing and wraparound behavior
- [x] Add unit tests for `StackManager` LIFO behavior and resolution callbacks
- [x] Extend integration tests to cover stack resolution after pass chains
- [x] Add comprehensive lifecycle tests (42 tests: start, mulligan, pause, resume, end, cleanup, loss conditions)
- [ ] Create regression tests comparing Go vs Java engine outputs for core scenarios
- [ ] Establish rules test harness mirroring Java's JUnit suite (CR regression coverage)
- [ ] Implement fuzz/invariant tests for state-based actions and stack integrity

## Migration & Compatibility
- [ ] Provide compatibility layer for existing Java client callbacks (message equivalence)
- [ ] Translate Java replay/log formats to Go for client consumption
- [ ] Document protocol changes and migration steps for server operators
- [ ] Benchmark Go engine against Java baseline (latency, throughput, memory)

