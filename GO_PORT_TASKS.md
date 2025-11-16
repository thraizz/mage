# Go Port Task Tracker

Minimal tracker for the tasks of migrating the Java Mage MTG server to Go. This file shall only contain tasks, no descriptions, documentation or anything. Just tasks.

Status legend:
- `[x]` Completed
- `[ ]` Pending / not yet started
- `[~]` In progress or partially implemented

## Engine Scaffolding & Lifecycle
- [x] Wire gRPC server to `MageEngine` via `EngineAdapter`
- [x] Provide `MageEngine` core skeleton that tracks games, players, and actions
- [x] Implement `TurnManager` mirroring MTG phase/step progression and priority handoff
- [x] Introduce `StackManager` with basic push/pop mechanics and simple resolution hooks
- [~] Extend stack resolution to support triggered abilities, replacement effects, and modal choices
- [ ] Implement priority windows for casting during stack resolution (e.g., mana abilities, nested responses)
- [ ] Persist stack/game events for replay and spectator synchronization
- [x] Add comprehensive error handling and rollback when resolution fails
- [x] Implement priority retention after casting (caster retains priority by default, only passes when explicitly passing)
- [x] Add state bookmarking and rollback mechanism for error recovery
- [ ] Implement comprehensive priority loop structure matching Java `playPriority()` pattern
- [x] Implement mulligan system
- [x] Implement game cleanup and resource disposal
- [x] Add complete lifecycle state validation

## Game State & Zones
- [x] Surface battlefield/stack state via `GameGetView`
- [x] Synchronize graveyard, exile, command, and hidden zones with engine updates
- [x] Track card ownership/controller changes (gain control, copying, phasing, etc.)
- [~] Implement continuous effects layer system (layers 1-7 per Comprehensive Rules)
- [x] Handle state-based actions (lethal damage, zero loyalty, legend rule, etc.)
- [x] Support counters (loyalty, +1/+1, poison, energy, experience)
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
- [x] Ensure stack legality checks (targets available, costs paid) prior to resolution
- [ ] Implement target selection flow for spells/abilities requiring targets
- [x] Add exhaustive integration tests covering multi-object stacks, counterspells, and priority loops
- [x] Resolve stack one item at a time with state-based action and triggered ability checks between each resolution
- [x] Implement triggered ability queue processing before priority (APNAP order: Active Player, Non-Active Player)
- [x] Add `checkStateAndTriggered()` method that runs before each priority (SBA → triggers → repeat until stable)
- [x] Handle simultaneous events between stack resolutions (process events after each resolution)

## Combat System 🚧 IN PROGRESS (~5% coverage, ~2,500 lines needed)
### Core Combat Infrastructure (P0 - Critical)
- [x] Implement `combatState` struct tracking attackers, blockers, groups, defenders, tapped creatures
- [x] Implement `combatGroup` struct for attacker-blocker-defender groupings with damage ordering
- [x] Add combat fields to `internalCard`: `Attacking`, `Blocking`, `AttackingWhat`, `BlockingWhat`
- [x] Add `combat *combatState` to `engineGameState`
- [x] Implement `ResetCombat(gameID)` - clear combat state at beginning of combat
- [x] Implement `SetAttacker(gameID, playerID)` - set attacking player
- [x] Implement `SetDefenders(gameID)` - identify all valid defenders (players, planeswalkers, battles)

### Attacker Declaration (P0 - Critical)
- [x] Implement `DeclareAttacker(gameID, creatureID, defenderID, playerID)` - declare single attacker
- [x] Implement `CanAttack(gameID, creatureID)` - validate creature can attack (summoning sickness, tapped, restrictions)
- [x] Implement `CanAttackDefender(gameID, creatureID, defenderID)` - validate can attack specific defender
- [~] Implement attacker tapping logic (tap unless vigilance)
- [x] Create/update combat groups when attackers declared
- [x] Implement `RemoveAttacker(gameID, attackerID)` - undo attacker declaration
- [x] Fire `EventAttackerDeclared` per attacker and `EventDeclaredAttackers` when complete

### Blocker Declaration (P0 - Critical)
- [x] Implement `DeclareBlocker(gameID, blockerID, attackerID, playerID)` - declare single blocker
- [x] Implement `CanBlock(gameID, blockerID, attackerID)` - validate creature can block (tapped, flying, restrictions)
- [x] Add blocker to combat group, update blocked status
- [x] Implement `AcceptBlockers(gameID)` - finalize blockers, check requirements/restrictions
- [x] Implement blocker ordering for multiple blockers on same attacker
- [x] Implement `RemoveBlocker(gameID, blockerID)` - undo blocker declaration
- [x] Fire `EventBlockerDeclared` per blocker and `EventDeclaredBlockers` when complete

### Damage Assignment & Application (P0 - Critical)
- [x] Implement `AssignCombatDamage(gameID, firstStrike bool)` - assign damage for attackers and blockers
- [x] Implement `combatGroup.assignDamageToBlockers()` - attacker damage to blockers with ordering
- [x] Implement `combatGroup.assignDamageToAttackers()` - blocker damage to attackers with ordering
- [x] Implement `ApplyCombatDamage(gameID)` - apply all assigned damage
- [x] Implement `combatGroup.applyDamage()` - mark damage on creatures and players
- [x] Handle unblocked attacker damage to defending player/permanent
- [x] Fire `EventCombatDamageAssigned` and `EventCombatDamageApplied` events

### Combat Cleanup (P0 - Critical)
- [x] Implement `EndCombat(gameID)` - move groups to formerGroups, clear current combat
- [x] Clear `Attacking` and `Blocking` flags on all creatures
- [x] Keep attacker tracking for "attacked this turn" queries
- [x] Fire `EventEndCombat` event
- [x] Implement `GetAttackedThisTurn(gameID, creatureID)` - check if creature attacked this turn
- [x] Clear damage tracking on creatures after combat

### First Strike & Double Strike (P1 - High Priority)
- [x] Implement `HasFirstOrDoubleStrike(gameID)` - check if any creature has first/double strike
- [x] Add first strike combat damage step before normal damage step
- [x] Implement `combatGroup.hasFirstOrDoubleStrike()` per group
- [x] Handle double strike creatures dealing damage in both steps
- [x] Prevent normal damage from creatures killed by first strike
- [x] Add ability constants (FirstStrikeAbility, DoubleStrikeAbility)
- [x] Track first strikers in combat state
- [x] Implement `dealsDamageThisStep()` logic for first/double strike

### Trample (P1 - High Priority)
- [x] Implement trample damage calculation (excess damage to defender)
- [x] Add `canDamageDefenderDirectly` flag to combat groups
- [x] Validate lethal damage assigned to blockers before overflow
- [x] Implement deathtouch + trample interaction (1 damage is lethal)
- [ ] Implement player damage assignment choice (requires UI system)
- [ ] Handle trample damage to planeswalkers/battles (requires planeswalker system)
- [ ] Support "trample over planeswalkers" rule

### Vigilance (P1 - High Priority)
- [x] Check for vigilance ability before tapping attacker
- [x] Track which attackers were tapped by attack in `attackersTappedByAttack`
- [ ] Support effects that grant vigilance during combat

### Flying & Reach (P1 - High Priority)
- [x] Implement flying restriction (can only be blocked by flying/reach)
- [x] Implement reach ability (can block flying)
- [x] Add `CanBlock` validation for flying/reach interactions
- [ ] Implement dragon blocking exception (requires subtype system and AsThough effects)
- [ ] Support effects that grant/remove flying during combat (requires continuous effects)

### Combat Events (P1 - High Priority)
- [x] Add `EventBeginCombat` - beginning of combat step
- [x] Add `EventDeclareAttackersStepPre` - before attacker declaration
- [x] Add `EventAttackerDeclared` - per attacker declared
- [x] Add `EventDeclaredAttackers` - all attackers declared
- [x] Add `EventDeclareBlockersStepPre` - before blocker declaration
- [x] Add `EventBlockerDeclared` - per blocker declared
- [x] Add `EventDeclaredBlockers` - all blockers declared
- [x] Add `EventCombatDamageStepPre` - before damage assignment
- [x] Add `EventCombatDamageApplied` - damage applied
- [x] Add `EventEndCombatStepPre` - before end of combat
- [x] Add `EventEndCombat` - combat ended
- [x] Add `EventUnblockedAttacker` - unblocked attacker after blockers declared
- [x] Add `EventRemovedFromCombat` - creature removed from combat

### Combat Validation & Requirements (P1 - High Priority)
- [x] Implement `CheckBlockRequirements(gameID, playerID)` - must block if able
- [x] Implement `CheckBlockRestrictions(gameID, playerID)` - can't block restrictions
- [x] Implement forced attack tracking (`creaturesForcedToAttack` map)
- [x] Implement must block tracking (`creatureMustBlockAttackers` map)
- [x] Validate minimum/maximum attacker counts
- [x] Validate minimum/maximum blocker counts per attacker

### Combat Triggers (P1 - High Priority)
- [x] Queue triggers on attacker declared (e.g., "Whenever ~ attacks")
- [x] Queue triggers on blocker declared (e.g., "Whenever ~ blocks")
- [x] Queue triggers on creature becomes blocked (e.g., "Whenever ~ becomes blocked")
- [ ] Queue triggers on combat damage dealt (e.g., "Whenever ~ deals combat damage")
- [ ] Queue triggers on creature dies in combat
- [x] Process triggers via existing `checkStateAndTriggered()` system

### Special Combat Mechanics (P2 - Medium Priority)
- [ ] Implement menace (must be blocked by 2+ creatures)
- [x] Implement deathtouch (any damage is lethal) - integrated with trample
- [x] Implement lifelink (gain life equal to damage dealt)
- [x] Implement defender (can't attack)
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
- [ ] Support multi-choice prompts (choose mode, targets, numbers, colors)
- [x] Implement mana payment flow (floating mana, cost reductions, hybrid costs)
- [x] Add concession, timeout, and match result handling aligned with rules

## Card Database & Ability Port
- [ ] Inventory Java ability/card modules and map to Go packages
- [ ] Generate Go card definitions from existing Java card data (expansions, tokens, abilities)
- [ ] Translate ability scripts (activated, triggered, static) into Go equivalents
- [ ] Port keyword ability handlers (flying, deathtouch, scry, etc.)
- [ ] Implement effect infrastructure (replacement effects, static ability watchers, continuous effects)
- [ ] Build automated verification to compare Java vs Go card behavior for representative samples

## Event System & Watchers
- [x] Mirror Java event bus for game events
- [x] Port watcher/listener infrastructure to track conditional abilities
- [x] Provide hooks for UI/websocket notifications (combat updates, triggers, log lines)
- [x] Capture analytics/metrics for stack depth, actions per turn, average response time
- [x] Queue triggered abilities instead of immediately pushing to stack (process via `checkTriggered()` before priority)

## Undo/Redo & State Management
- [x] Implement per-player stored bookmarks for action undo
- [x] Add player-initiated undo command
- [x] Implement strategic bookmark placement in game flow
- [x] Add bookmark invalidation rules
- [x] Implement turn rollback system with turn-level snapshots
- [x] Integrate undo/redo with error recovery system

## Persistence, Replays & Recovery
- [x] Store game snapshots for reconnection and spectating
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

