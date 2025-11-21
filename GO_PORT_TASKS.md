# Go Port Task Tracker

Minimal tracker for the tasks of migrating the Java Mage MTG server to Go. This file shall only contain tasks, no descriptions, documentation or anything. Just tasks.

Status legend:
- `[x]` Completed
- `[ ]` Pending / not yet started
- `[-]` In progress or partially implemented
- `[~]` Canceled

## Engine Scaffolding & Lifecycle
- [x] Wire gRPC server to `MageEngine` via `EngineAdapter`
- [x] Provide `MageEngine` core skeleton that tracks games, players, and actions
- [x] Implement `TurnManager` mirroring MTG phase/step progression and priority handoff
- [x] Introduce `StackManager` with basic push/pop mechanics and simple resolution hooks
- [-] Extend stack resolution to support triggered abilities, replacement effects, and modal choices
- [x] Implement priority windows for casting during stack resolution (e.g., mana abilities, nested responses)
  - [x] Implement mana ability activation during spell/ability resolution (Rule 117.1d, 605.3a)
    - [x] Add `ActivateManaAbility()` method that can be called during resolution
    - [x] Track resolution context (which spell/ability is currently resolving)
    - [x] Allow mana abilities when: (1) player has priority, (2) casting spell/activating ability that needs mana, (3) rule/effect asks for mana
    - [x] Ensure mana abilities resolve immediately without going on stack (Rule 605.3b)
    - [x] Implement triggered mana abilities that resolve immediately after triggering mana ability (Rule 605.4a)
    - [x] Prevent mana ability re-activation until current activation resolves (Rule 605.3c)
  - [x] Implement special actions during resolution (Rule 116, 117.1c)
    - [x] Track which special actions are allowed during resolution vs. only during main phase
    - [x] Implement special action execution that doesn't use the stack
    - [x] Grant priority to player after special action (Rule 116.3)
    - [x] Handle special actions that can be taken "any time you have priority": face-down creatures (116.2b), ending effects (116.2c), ignoring static abilities (116.2d)
    - [x] Handle special actions restricted to main phase + empty stack: playing lands (116.2a), companion (116.2g), plot (116.2k), unlock (116.2m)
  - [x] Implement nested spell/ability casting during resolution (Rule 117.2e, 608.2)
    - [x] Support casting copies of spells during resolution (Rule 707.12 - "cast while another spell or ability is resolving")
    - [x] Implement linked abilities that allow casting/activating during resolution
    - [x] Track nested resolution depth to prevent infinite recursion
    - [x] Ensure proper priority handling within nested resolution context
    - [x] Handle mana payment flow during nested casting (allow mana abilities)
  - [x] Implement resolution payment/choice windows
    - [x] Add player choice prompts during resolution (modes, targets, X values, etc.) per Rule 608.2
    - [x] Implement APNAP order for multi-player choices during resolution (Rule 608.2e)
    - [x] Allow mana ability activation during cost payment within resolution
    - [x] Support special action activation for cost payment (e.g., Convoke, delve)
    - [x] Track payment state (which costs paid, which remain) during resolution
  - [x] Add comprehensive testing for priority windows
    - [x] Test mana ability activation while paying for spell during another spell's resolution
    - [x] Test special actions (morph, etc.) during priority windows
    - [x] Test nested spell casting (copy effects like Isochron Sceptron)
    - [x] Test linked abilities that cast during resolution
    - [x] Test priority retention and passing during nested resolution
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
- [x] Implement continuous effects layer system (layers 1-7 per Comprehensive Rules) - fully integrated
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

## Combat System
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
- [x] Implement attacker tapping logic (tap unless vigilance)
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
- [x] Support effects that grant vigilance during combat

### Flying & Reach (P1 - High Priority)
- [x] Implement flying restriction (can only be blocked by flying/reach)
- [x] Implement reach ability (can block flying)
- [x] Add `CanBlock` validation for flying/reach interactions
- [ ] Implement dragon blocking exception (requires subtype system and AsThough effects)
- [x] Support effects that grant/remove flying during combat

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
- [x] Queue triggers on attacker declared (e.g., "Whenever - attacks")
- [x] Queue triggers on blocker declared (e.g., "Whenever - blocks")
- [x] Queue triggers on creature becomes blocked (e.g., "Whenever - becomes blocked")
- [x] Queue triggers on combat damage dealt (e.g., "Whenever - deals combat damage")
- [x] Queue triggers on creature dies in combat
- [x] Process triggers via existing `checkStateAndTriggered()` system

### Special Combat Mechanics (P2 - Medium Priority)
- [x] Implement menace (must be blocked by 2+ creatures)
- [x] Implement deathtouch (any damage is lethal) - integrated with trample
- [x] Implement lifelink (gain life equal to damage dealt)
- [x] Implement defender (can't attack)
- [x] Implement "can't be blocked" effects
- [x] Implement "must be blocked if able" effects (lure)
- [x] Implement "attacks each combat if able" effects

### Planeswalker & Battle Combat (P2 - Medium Priority)
- [x] Support attacking planeswalkers (planeswalkers added to defenders)
- [x] Implement damage-to-loyalty conversion (Rule 306.8, 120.3c)
- [x] Update lethal damage calculation for planeswalkers (loyalty-based)
- [x] Support lifelink with planeswalker damage
- [x] Support deathtouch with planeswalker damage
- [x] Block attacks on own planeswalkers
- [x] Track which planeswalkers are attacked for triggers (PlayersAttackedThisTurnWatcher)
- [x] Handle planeswalker removal during combat (graceful damage handling)
- [ ] Support attacking battles
- [ ] Implement planeswalker damage redirection rules (pre-2018 deprecated rules)

### Damage Division (P2 - Medium Priority) - Modern Rules (No Ordering)
- [x] Implement damage division for attacker with multiple blockers (Rule 510.1c)
- [x] Implement damage division for blocker blocking multiple attackers (Rule 510.1d)
- [x] Add AssignAttackerDamage() API for player damage choices
- [x] Add AssignBlockerDamage() API for player damage choices
- [x] Implement default damage division (even split, lethal for trample)
- [x] Validate damage assignments (total equals power, valid targets)
- [x] Handle blockers in multiple combat groups correctly
- [ ] UI for damage division prompts (multi-amount dialog)
- [ ] Support "you choose damage order" effects (Defensive Formation, etc.)

### Banding (P3 - Low Priority, Complex)

- [x] Add ability constant for banding detection
- [x] Add band tracking fields to internalCard (BandedCards)
- [x] Damage assignment control - defending player assigns (Rule 702.22j)
- [x] Damage assignment control - attacking player assigns (Rule 702.22k)
- [x] Update AssignAttackerDamage/AssignBlockerDamage APIs with player validation
- [x] Comprehensive tests for damage assignment control
- [~] Band formation during attack declaration (bidirectional tracking)
- [~] Block propagation across band members (Rule 702.22h)
- [~] "Bands with other" variants (by subtype/supertype/name)
- [~] Edge cases (removal, banding lost mid-combat, multiple bands)
- [~] Band formation UI/API

### Combat Removal & Interruption (P2 - Medium Priority)
- [x] Implement `RemoveFromCombat(gameID, creatureID)` - remove during combat
- [x] Handle creature removal during attacker declaration
- [x] Handle creature removal during blocker declaration
- [x] Handle creature removal during damage assignment
- [x] Update combat groups when creatures removed
- [x] Implement `CheckForRemoveFromCombat()` - automatic removal when creatures lose creature type
- [x] Integrate CheckForRemoveFromCombat into all combat steps (declare attackers, declare blockers, damage steps)
- [x] Add comprehensive tests for automatic type-loss removal (6 tests)
- [ ] Handle blink/flicker during combat (removed and returns as new object)
- [ ] Handle phase out during combat

### Combat Integration with Turn Structure (P0 - Critical)
- [x] Wire `ResetCombat()` to beginning of combat step
- [x] Wire `SetAttacker()` and `SetDefenders()` to beginning of combat
- [x] Wire attacker declaration to declare attackers step
- [x] Wire blocker declaration to declare blockers step
- [x] Wire first strike damage to first strike damage step
- [x] Wire normal damage to combat damage step
- [x] Wire `EndCombat()` to end of combat step
- [x] Add combat damage steps to turn structure if first strike exists

### Combat Testing (P0 - Critical)
- [x] Test single attacker, no blockers (damage to player)
- [x] Test single attacker, single blocker (damage to creatures)
- [x] Test multiple attackers, no blockers
- [x] Test multiple attackers, multiple blockers
- [x] Test multiple blockers on single attacker (damage ordering)
- [x] Test creature death from combat damage
- [x] Test player death from combat damage
- [x] Test vigilance (no tap on attack)
- [x] Test first strike damage (kill before normal damage)
- [x] Test double strike damage (damage in both steps)
- [x] Test trample damage (overflow to player)
- [x] Test flying/reach restrictions
- [x] Test combat triggers firing
- [x] Test combat events firing
- [x] Test removal during combat (all phases)
- [x] Test blocker requirements/restrictions
- [x] Test attacker requirements/restrictions

### Combat View & Display (P1 - High Priority)
- [x] Populate `EngineCombatView` with actual combat state
- [x] Populate `EngineCombatGroupView` for each combat group
- [x] Show attacking creatures in game view
- [x] Show blocking creatures in game view
- [x] Show damage assignments in game view
- [x] Update combat view after each declaration/assignment

## Player Interaction & Prompts
- [x] Emit prompts when priority passes require player response
- [ ] Support multi-choice prompts (choose mode, targets, numbers, colors)
- [x] Implement mana payment flow (floating mana, cost reductions, hybrid costs)
- [x] Add concession, timeout, and match result handling aligned with rules

## Card Database & Ability Port
- [ ] Inventory Java ability/card modules and map to Go packages
- [ ] Create complete list of all cards and abilities we have to port
- [ ] Generate Go card definitions from existing Java card data (expansions, tokens, abilities) based off the list, one by one
- [ ] Translate ability scripts (activated, triggered, static) into Go equivalents
- [ ] Port keyword ability handlers (flying, deathtouch, scry, etc., check RULES.txt and Java implementation)
- [ ] Implement effect infrastructure (replacement effects, static ability watchers, continuous effects)
- [ ] Build automated verification to compare Java vs Go card behavior for representative samples
- [ ] **Re-enable disabled integration tests** that expect specific cards (8 tests disabled - see comments in test files for details)
- [ ] Add abilities that are missing in the java implementation, e.g. face-down cards

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
- [ ] Translate Java replay/log formats to Go for client consumption
- [ ] Document protocol changes and migration steps for server operators
- [ ] Benchmark Go engine against Java baseline (latency, throughput, memory, stability)

