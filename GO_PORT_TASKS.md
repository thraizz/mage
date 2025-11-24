# Go Port Task Tracker

Minimal tracker for the tasks of migrating the Java Mage MTG server to Go. This file shall only contain tasks, no descriptions, documentation or anything. Just tasks.

Status legend:
- `[x]` Completed
- `[ ]` Pending / not yet started
- `[-]` In progress or partially implemented
- `[~]` Canceled

# ENGINE WORK

## Engine Scaffolding & Lifecycle
- [x] Wire gRPC server to `MageEngine` via `EngineAdapter`
- [x] Provide `MageEngine` core skeleton that tracks games, players, and actions
- [x] Implement `TurnManager` mirroring MTG phase/step progression and priority handoff
- [x] Introduce `StackManager` with basic push/pop mechanics and simple resolution hooks
- [x] Extend stack resolution to support triggered abilities, replacement effects, and modal choices
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
- [-] Implement comprehensive priority loop structure matching Java `playPriority()` pattern
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
- [x] Support casting spells/activating abilities while another object is resolving (linked abilities)
- [x] Implement replacement/prevention effects that modify or negate stack resolution
- [x] Ensure stack legality checks (targets available, costs paid) prior to resolution
- [x] Implement target selection flow for spells/abilities requiring targets
- [x] Add exhaustive integration tests covering multi-object stacks, counterspells, and priority loops
- [x] Resolve stack one item at a time with state-based action and triggered ability checks between each resolution
- [x] Implement triggered ability queue processing before priority (APNAP order: Active Player, Non-Active Player)
- [x] Add `checkStateAndTriggered()` method that runs before each priority (SBA → triggers → repeat until stable)
- [x] Handle simultaneous events between stack resolutions (process events after each resolution)

## Engine Integration Systems
- [x] Implement EnhancedStackManager with ability integration (Rule 405, 608)
- [x] Implement AbilityRegistry for UUID-based ability tracking and retrieval
- [x] Implement TargetSelectionManager with validation and legal target calculation (Rule 115)
- [x] Implement ContinuousEffectsManager with automatic layer recalculation (Rule 613)
- [x] Implement AbilityActivationManager for spell casting and ability activation workflows (Rule 601, 602, 605)
- [x] Implement CombatIntegrationManager connecting combat steps to triggered abilities (Rule 508-510)
- [x] Integrate PriorityManager with layer recalculation before SBA checks
- [x] Implement event adapter bridging rules events to abilities system triggers

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
- [x] **Phase 1 Complete: Card Factory Infrastructure**
  - [x] Export 30,459 cards from Java source files to CSV
  - [x] Implement Card type with full game state
  - [x] Create CardFactory interface and implementation
  - [x] Build Registry system for card builders (self-registering)
  - [x] Add CardInfo helper with type checking methods
  - [x] Write unit tests (12 tests passing)
- [x] **Phase 2 Complete: Ability Framework**
  - [x] Define core ability interfaces (Effect, Cost, Target, Ability)
  - [x] Implement 40+ effects (damage, draw, destroy, counter, boost, tap, mana, tokens, counters, mill, bounce, exile, search, scry, surveil, etc.)
  - [x] Implement 7 cost types (mana, tap, sacrifice, discard, pay life, composite)
  - [x] Implement 10+ target filters (any, creature, player, permanent, spell, etc.)
  - [x] Write builder API for fluent ability construction
  - [x] Write 8 test files with comprehensive unit tests (all passing)
- [ ] **Phase 3: Manual Test Cards** (20 cards)
  - [ ] Implement 5 basic lands
  - [ ] Implement 3 vanilla creatures
  - [ ] Implement 6 simple spells (Lightning Bolt, Murder, etc.)
  - [ ] Implement 3 keyword creatures (flying, vigilance, etc.)
  - [ ] Implement 3 activated abilities (Llanowar Elves, etc.)
  - [ ] Write 60+ integration tests for manual cards
- [-] **Phase 4: Transpiler Development** (8 weeks)
  - [x] Build Java AST parser for card files
  - [x] Create ability mapper (Java effects → Go effects, 150+ effect mappings, 70+ counter types, 711 token types)
  - [x] Implement Go code generator from AST
  - [x] Add counter support (AddCountersSourceEffect, AddCountersTargetEffect, RemoveCounterTargetEffect)
  - [x] Add token support (CreateTokenEffect with 711 token types via registry)
  - [x] Add triggered ability extraction (EntersBattlefieldAbility, EntersBattlefieldControlledTriggeredAbility)
  - [x] Implement balanced parentheses parser for nested function calls
  - [x] Add smart import detection (auto-add counters/token packages when needed)
  - [x] Test transpiler with complex cards (Yorvo Lord of Garenbrig, Regisaur Alpha) - ✅ Both transpile without TODOs
  - [x] Create batch generation pipeline
  - [x] Generate all 30,439 remaining cards (30,404 files generated = 99.8% success rate)
  - [ ] Fix unmapped effects (estimated 1000-2000 cards need manual fixes)
  - [x] Add more triggered ability types (DiesTriggeredAbility, AttacksTriggeredAbility, etc.)
  - [x] Add static ability support (GainAbilityControlledEffect, BoostControlledEffect, etc.)
  - [x] Add activated ability support (SimpleActivatedAbility with costs)
  - [ ] Manually implement complex cards (planeswalkers, transforming)
- [ ] Build automated verification to compare Java vs Go card behavior for representative samples
- [ ] **Re-enable disabled integration tests** that expect specific cards (8 tests disabled - see comments in test files for details)
- [x] Add abilities that are missing in the java implementation, e.g. face-down cards

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

### Log-Entry-Level Rollback (Fine-Grained State Management)
- [ ] Implement comprehensive game state snapshot system
  - [ ] Create serializable GameSnapshot struct with all game state
  - [ ] Include all zones (battlefield, hand, library, graveyard, exile, stack, command)
  - [ ] Include all player state (life, counters, mana pools, combat state)
  - [ ] Include continuous effects, replacement effects, and triggers
  - [ ] Include turn/phase/step information
  - [ ] Include priority state and passed flags
  - [ ] Make snapshots immutable and deep-copyable
  - [ ] Optimize snapshot size (use copy-on-write where possible)
- [ ] Implement log entry tracking system
  - [ ] Create LogEntry struct with timestamp, action type, player, description
  - [ ] Assign unique sequential IDs to each log entry
  - [ ] Track game state snapshot before each log entry
  - [ ] Store log entries in circular buffer (configurable size, e.g., last 1000 entries)
  - [ ] Support different log entry types (card played, ability activated, phase change, damage dealt, etc.)
  - [ ] Include diff information (what changed from previous state)
  - [ ] Compress old snapshots to save memory
- [ ] Implement snapshot creation triggers
  - [ ] Create snapshot before each player action (play card, activate ability, pass priority)
  - [ ] Create snapshot after each stack resolution
  - [ ] Create snapshot at each phase/step boundary
  - [ ] Create snapshot after each combat step
  - [ ] Create snapshot after state-based actions
  - [ ] Create snapshot after triggered abilities are queued
  - [ ] Skip redundant snapshots (no state change)
- [ ] Implement rollback execution
  - [ ] Add RollbackToLogEntry(gameID, logEntryID) method
  - [ ] Validate rollback target (entry must exist and be valid)
  - [ ] Restore complete game state from snapshot
  - [ ] Restore all zones with correct card positions and states
  - [ ] Restore player states (life, counters, mana pools)
  - [ ] Restore stack state
  - [ ] Restore continuous effects and layer system
  - [ ] Restore turn/phase/priority state
  - [ ] Clear any temporary state not in snapshot
  - [ ] Trigger state recalculation (layers, SBAs) after rollback
- [ ] Add rollback validation
  - [ ] Prevent rollback during resolution (must be at priority window)
  - [ ] Validate game hasn't ended
  - [ ] Ensure target log entry is accessible
  - [ ] Check admin permissions for rollback (if not player-initiated)
  - [ ] Prevent rollback past turn boundaries in competitive mode
- [ ] Implement log entry browsing
  - [ ] Add GetGameLog(gameID, startIndex, count) RPC method
  - [ ] Return paginated log entries with descriptions
  - [ ] Include snapshot availability flag per entry
  - [ ] Support filtering by player, action type, or time range
  - [ ] Add GetLogEntryCount(gameID) method
- [ ] Add rollback UI support
  - [ ] Add proto definitions for rollback RPCs
  - [ ] Add RollbackToLogEntry RPC method
  - [ ] Send log update events via WebSocket when new entries added
  - [ ] Send rollback notification event when rollback occurs
  - [ ] Include log entry in GameView for client display
- [ ] Implement snapshot persistence for replay
  - [ ] Store snapshots in match history for replay functionality
  - [ ] Compress snapshots for storage (gzip or similar)
  - [ ] Support loading game from any log entry for replay
  - [ ] Implement snapshot pruning (keep every Nth snapshot, keep all recent)
- [ ] Add snapshot memory management
  - [ ] Implement circular buffer with max size (e.g., 1000 snapshots)
  - [ ] Evict oldest snapshots when buffer full
  - [ ] Keep critical snapshots (turn boundaries) longer
  - [ ] Add configuration for snapshot buffer size
  - [ ] Monitor snapshot memory usage
  - [ ] Add metrics for snapshot creation rate and size
- [ ] Implement admin rollback controls
  - [ ] Add admin-only rollback permission check
  - [ ] Allow admin to rollback any game to any log entry
  - [ ] Log all rollback actions with admin identifier
  - [ ] Add rollback reason field for auditing
  - [ ] Broadcast rollback notification to all players/watchers
- [ ] Add rollback testing
  - [ ] Test rollback after card played
  - [ ] Test rollback after ability activated
  - [ ] Test rollback during combat (each step)
  - [ ] Test rollback after stack resolution
  - [ ] Test rollback across multiple turns
  - [ ] Test snapshot memory limits
  - [ ] Test rollback with complex board states (tokens, counters, effects)
  - [ ] Verify game state consistency after rollback
- [ ] Implement delta-based snapshots for efficiency
  - [ ] Store full snapshot at turn boundaries
  - [ ] Store deltas (changes only) for intermediate states
  - [ ] Reconstruct full state from base snapshot + deltas
  - [ ] Measure memory savings vs reconstruction cost

## Persistence, Replays & Recovery
- [x] Store game snapshots for reconnection and spectating
- [x] Implement replay recording (step-by-step action logs)
- [x] Ensure deterministic serialization for saved games and tournaments
- [x] Add checksum/validation to guard against divergent game state

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

# SERVER WORK 

## gRPC/Protocol Buffers Infrastructure
- [x] Define all 60+ RPC methods in proto files (auth, room, table, game, tournament, draft, chat, admin)
- [x] Define data models in proto (TableView, GameView, MatchOptions, etc.)
- [x] Define WebSocket event messages for server push
- [x] Generate Go code from proto files
- [x] Generate TypeScript client code from proto files

## Server Core Infrastructure
- [x] Initialize Go module and project structure
- [x] Set up Viper configuration management (full Config struct with all subsystems)
- [x] Implement comprehensive configuration file (config.yaml with 12 sections)
- [x] Implement environment variable overrides via Viper
- [x] Implement configuration validation
- [x] Create example configurations (config.example.yaml, config.dev.yaml)
- [x] Implement PostgreSQL connection pooling (pgx)
- [x] Create database migration system (4 migrations: users, cards, stats, table_records)
- [x] Implement migration up/down commands in Makefile
- [x] Implement migration creation script
- [x] Set up Zap structured logging
- [x] Implement configurable log levels (debug, info, warn, error)
- [x] Implement configurable log formats (json, console)
- [x] Implement Argon2id password hashing
- [x] Set up health check endpoint (/health, /ready, /live)
- [x] Implement health check with database status
- [x] Implement readiness probe for Kubernetes
- [x] Implement liveness probe for Kubernetes
- [x] Create comprehensive Makefile with 20+ commands
- [x] Implement proto generation script (generate_proto.sh)
- [x] Implement build system with LDFLAGS optimization
- [x] Implement test runner with race detection and coverage
- [ ] Implement metrics/Prometheus instrumentation

## Plugin System & Game Types
- [x] Implement plugin registry system (thread-safe registration)
- [x] Define GameType interface (name, min/max players, description)
- [x] Define TournamentType interface (name, description, draft support)
- [x] Define PlayerType interface (name, AI flag, description)
- [x] Implement game type registry with Get/List methods
- [x] Register 6 game types (TwoPlayerDuel, FreeForAll, CommanderFreeForAll, CommanderDuel, Brawl, CanadianHighlander)
- [x] Register 3 tournament types (Constructed, BoosterDraft, Sealed)
- [x] Register 3 player types (Human, ComputerMAX, ComputerDraft)
- [x] Implement plugin configuration in config.yaml
- [x] Auto-register all plugins via init() functions

## Session & Authentication
- [x] Implement in-memory session manager with expiration
- [x] Implement session cleanup goroutine
- [x] Add concurrent request locking per session
- [x] Implement user registration (anonymous and authenticated modes)
- [x] Implement username validation and uniqueness checks
- [x] Implement password reset token generation (6-digit)
- [x] Implement email service (SMTP/Mailgun)
- [ ] Add Redis-backed session store for production
  - [ ] Add Redis client dependency
  - [ ] Implement Redis session repository
  - [ ] Store session data in Redis with TTL
  - [ ] Support session migration between server instances
  - [ ] Add Redis connection pooling
  - [ ] Make Redis optional (fallback to in-memory for dev)

## Database Repositories
- [x] Implement user repository (CRUD + queries)
- [x] Implement card repository with full-text search
- [x] Implement card repository caching (10k card LRU cache)
- [x] Implement stats repository with Glicko rating queries
- [ ] Implement table records repository (persist table/match history)
- [x] Implement deck repository (store/retrieve user deck collections)
  - [x] Create deck table schema (user_id, deck_name, format, description, cards JSON)
  - [x] Implement CreateDeck method
  - [x] Implement GetDecksByUser method
  - [x] Implement GetDecksByUserAndFormat method
  - [x] Implement UpdateDeck method
  - [x] Implement DeleteDeck method
  - [x] Implement DeleteByUserAndID method (with ownership check)
  - [x] Implement GetDeckByID method
- [x] Implement match history repository (player game history for stats/replay)
  - [x] Create match_history table schema (with JSONB players, indexes, constraints)
  - [x] Implement SaveMatch method
  - [x] Implement GetMatchesByUser method (paginated)
  - [x] Implement GetMatchByID method
  - [x] Implement CountMatchesByUser method
  - [x] Implement GetRecentMatches method (for lobby display)
  - [x] Implement GetMatchesByGameType method
  - [x] Implement GetMatchesByTournament method

## gRPC Server & Interceptors
- [x] Implement gRPC server bootstrap with graceful shutdown
- [x] Implement session validation interceptor
- [x] Implement admin authorization interceptor
- [x] Implement logging interceptor
- [x] Implement panic recovery interceptor
- [x] Implement metrics interceptor (placeholder - needs Prometheus integration)
- [x] Implement rate limiting interceptor (placeholder)

## WebSocket Server for Push Events
- [x] Implement WebSocket server with Gorilla WebSocket
- [x] Handle WebSocket upgrade and session validation
- [x] Forward ServerEvent messages from session to WebSocket
- [x] Implement ping keep-alive (configurable interval)
- [x] Implement write deadline handling
- [x] Implement graceful connection cleanup
- [ ] Add reconnection handling (resume from last message ID)
  - [ ] Assign sequence IDs to all WebSocket messages
  - [ ] Buffer recent messages per session (circular buffer)
  - [ ] Accept reconnect with last received message ID
  - [ ] Replay missed messages on reconnection
  - [ ] Expire buffered messages after timeout
- [ ] Implement message compression for large payloads
  - [ ] Enable WebSocket permessage-deflate extension
  - [ ] Configure compression threshold (only compress >1KB)
  - [ ] Add compression level configuration
  - [ ] Measure compression ratio and performance impact

## Authentication RPC Methods
- [x] Implement ConnectUser (login with credentials)
- [x] Implement Ping (session keep-alive)
- [x] Implement AuthRegister (new user registration)
- [x] Implement AuthSendTokenToEmail (password reset)
- [x] Implement AuthResetPassword (reset with token)
- [x] Implement ConnectAdmin (admin login)
- [x] Implement ConnectSetUserData (client preferences)

## Server Info RPC Methods
- [x] Implement GetServerState (status, player count)
- [x] Implement ServerGetPromotionMessages (MOTD, announcements)
- [x] Implement ServerAddFeedbackMessage (bug reports)

## Room/Lobby RPC Methods
- [x] Implement ServerGetMainRoomId (lobby ID)
- [x] Implement RoomGetUsers (online player list)
- [x] Implement RoomGetFinishedMatches (recent matches)
- [x] Implement RoomGetAllTables (table list)
- [x] Implement RoomGetTableById (table details)

## Table Management RPC Methods
- [x] Implement RoomCreateTable (create new table)
- [x] Implement RoomJoinTable (join existing table)
- [x] Implement RoomLeaveTableOrTournament (leave table)
- [x] Implement RoomWatchTable (spectator join)
- [x] Implement TableSwapSeats (change seat)
- [x] Implement TableRemove (host closes table)
- [x] Implement TableIsOwner (check host status)
- [x] Implement RoomCreateTournament
- [x] Implement RoomJoinTournament
- [x] Implement RoomWatchTournament

## Deck Management RPC Methods
- [x] Implement DeckSubmit (submit deck for table/tournament)
- [x] Implement DeckSave (save deck to collection with format and description)
- [x] Implement DeckList (get user's decks)
  - [x] Add proto definition for DeckListRequest/Response with DeckInfo message
  - [x] Implement gRPC handler in grpc_table.go
  - [x] Call deck repository GetDecksByUser / GetDecksByUserAndFormat
  - [x] Return list of deck metadata (id, name, format, description, card counts, timestamps)
- [x] Implement DeckDelete (delete deck)
  - [x] Add proto definition for DeckDeleteRequest/Response
  - [x] Implement gRPC handler with ownership check
  - [x] Call deck repository DeleteByUserAndID
- [x] Implement DeckGet (get single deck by ID)
  - [x] Add proto definition for DeckGetRequest/Response
  - [x] Implement gRPC handler with ownership validation
  - [x] Return full deck info and card lists
- [ ] Implement DeckValidate (check format legality: Standard, Modern, Commander, etc.)
  - [ ] Add proto definition for DeckValidateRequest/Response
  - [ ] Implement format rules engine (card legality, deck size, banned list)
  - [ ] Check card count restrictions (60 minimum, 4-of limit, etc.)
  - [ ] Check format-specific rules (Commander color identity, sideboard size)
  - [ ] Return validation errors with card details
- [ ] Implement DeckImport (import from MTGO/Arena text format)
  - [ ] Add proto definition for DeckImportRequest/Response
  - [ ] Parse MTGO format (quantity + card name per line)
  - [ ] Parse Arena format (deck/sideboard sections)
  - [ ] Parse plaintext format (deck/sideboard sections)
  - [ ] Resolve card names to card IDs via repository
  - [ ] Handle missing/ambiguous card names gracefully

## Game Execution RPC Methods
- [x] Implement MatchStart (start match from table)
- [x] Implement GameJoin (join ongoing game)
- [x] Implement GameWatchStart (spectator join game)
- [x] Implement GameWatchStop (spectator leave game)
- [x] Implement GameGetView (get full game state)
- [x] Implement SendPlayerUUID (target selection)
- [x] Implement SendPlayerString (text input)
- [x] Implement SendPlayerBoolean (yes/no choice)
- [x] Implement SendPlayerInteger (number input)
- [x] Implement SendPlayerManaType (mana color choice)
- [x] Implement SendPlayerAction (pass priority, play card, etc.)
- [x] Implement MatchQuit (concede)

## Game Rollback RPC Methods
- [ ] Implement GetGameLog (retrieve game log entries)
  - [ ] Add proto definition for GetGameLogRequest/Response
  - [ ] Add LogEntry message with id, timestamp, player, action, description
  - [ ] Implement gRPC handler to query log entries
  - [ ] Support pagination (startIndex, count)
  - [ ] Support filtering by player, action type, time range
  - [ ] Include snapshot availability flag per entry
- [ ] Implement GetLogEntryCount (get total log entry count)
  - [ ] Add proto definition for GetLogEntryCountRequest/Response
  - [ ] Return total number of log entries for game
- [ ] Implement RollbackToLogEntry (rollback game state)
  - [ ] Add proto definition for RollbackToLogEntryRequest/Response
  - [ ] Validate rollback permission (admin or casual mode)
  - [ ] Validate log entry exists and is accessible
  - [ ] Execute rollback via game engine
  - [ ] Broadcast rollback event to all players/watchers
  - [ ] Return updated game view after rollback
- [ ] Implement GetSnapshotInfo (get snapshot metadata)
  - [ ] Add proto definition for GetSnapshotInfoRequest/Response
  - [ ] Return snapshot availability for log entries
  - [ ] Include snapshot size and timestamp
  - [ ] Show memory usage statistics

## Chat RPC Methods
- [x] Implement ChatJoin (join chat channel)
- [x] Implement ChatLeave (leave chat channel)
- [x] Implement ChatSendMessage (send message)
- [x] Implement ChatFindByRoom (get lobby chat)
- [x] Implement ChatFindByTable (get table chat)
- [x] Implement ChatFindByGame (get game chat)
- [x] Implement ChatFindByTournament (get tournament chat)
- [ ] Implement whisper/private message support
  - [ ] Add proto definition for ChatWhisperRequest/Response
  - [ ] Implement gRPC handler for direct messages
  - [ ] Add private message routing in chat manager
  - [ ] Store whisper history per user pair
  - [ ] Add whisper notifications via WebSocket
- [ ] Implement chat rate limiting
  - [ ] Add rate limiter per user (e.g., 5 messages per 10 seconds)
  - [ ] Track message timestamps in session
  - [ ] Return error when rate limit exceeded
  - [ ] Add admin exemption from rate limits
- [ ] Implement HTML sanitization (bluemonday)
  - [ ] Add bluemonday dependency
  - [ ] Create sanitization policy (allow basic formatting, strip scripts)
  - [ ] Sanitize all incoming chat messages
  - [ ] Sanitize whispers and broadcast messages

## Draft RPC Methods
- [x] Implement DraftJoin (join draft)
- [x] Implement SendDraftCardPick (pick card from pack)
- [x] Implement SendDraftCardMark (mark card for review)
- [x] Implement DraftSetBoosterLoaded (client ready for pack)
- [x] Implement DraftQuit (leave draft)
- [ ] Implement booster pack generation from card repository
  - [ ] Create booster pack generator interface
  - [ ] Implement set-specific pack composition (common/uncommon/rare ratios)
  - [ ] Query card repository by set code and rarity
  - [ ] Handle foils and special slots (mythic rare, double-faced cards)
  - [ ] Support custom cube drafts (user-defined card pools)
  - [ ] Cache booster configurations per set

## Tournament RPC Methods
- [x] Implement TournamentJoin (register for tournament)
- [x] Implement TournamentStart (begin tournament)
- [x] Implement TournamentQuit (drop from tournament)
- [x] Implement TournamentFindById (get tournament details)
- [x] Implement Swiss pairing algorithm
- [ ] Implement elimination bracket generation
  - [ ] Implement single elimination bracket (2^n players, bye handling)
  - [ ] Implement double elimination bracket (winners/losers brackets)
  - [ ] Add bracket seeding options (random, by rating, by record)
  - [ ] Generate bracket structure and initial pairings
  - [ ] Handle advancement logic (winner moves up, loser drops/moves to losers)
  - [ ] Add bracket visualization data for client rendering

## Admin RPC Methods
- [x] Implement AdminGetUsers (list all users)
- [x] Implement AdminDisconnectUser (kick user)
- [x] Implement AdminMuteUser (mute in chat)
- [x] Implement AdminLockUser (temporary ban)
- [x] Implement AdminActivateUser (unlock account)
- [x] Implement AdminToggleActivateUser (toggle active status)
- [x] Implement AdminEndUserSession (force disconnect)
- [x] Implement AdminTableRemove (force close table)
- [x] Implement AdminSendBroadcastMessage (server announcement)

## Room/Lobby Management
- [x] Implement GamesRoom (main lobby)
- [x] Implement lobby features (user list, table list, finished matches)
- [x] Implement room update broadcasting
- [ ] Implement real-time table updates via WebSocket
  - [ ] Send TableCreated event when table is created
  - [ ] Send TableUpdated event when player joins/leaves
  - [ ] Send TableRemoved event when table is closed
  - [ ] Send TableStateChanged event (WAITING → STARTING → DUELING)
  - [ ] Filter events per user subscription (only subscribed rooms)
- [ ] Implement real-time user join/leave notifications
  - [ ] Send UserJoinedRoom event when user connects
  - [ ] Send UserLeftRoom event when user disconnects
  - [ ] Send UserStatusChanged event (idle, in-game, in-draft)
  - [ ] Update online user count in real-time

## Table Controller
- [x] Implement TableController state machine (WAITING → STARTING → DUELING → FINISHED)
- [x] Implement player seat assignment and swapping
- [x] Implement deck validation hooks
- [x] Implement match creation logic
- [x] Implement host controls (kick player, start game)
- [ ] Integrate with game/tournament controllers
  - [ ] Trigger match creation when table starts
  - [ ] Link table to active game instance
  - [ ] Update table state based on game completion
  - [ ] Handle tournament table creation (automated pairing)

## Game Controller
- [x] Implement GameController state management
- [x] Implement player action queue and processing
- [x] Implement watcher management
- [x] Integrate with MageEngine
- [x] Implement game view generation for clients
- [ ] Implement spectator view (hide hidden information)
  - [ ] Create filtered game view for watchers
  - [ ] Hide opponent hands from spectator view
  - [ ] Hide face-down cards and unrevealed information
  - [ ] Show revealed zones (battlefield, graveyard, exile)
  - [ ] Allow spectators to see both players' perspectives optionally

## Tournament System
- [x] Implement TournamentController state machine
- [x] Implement round management
- [x] Implement tournament view generation
- [ ] Integrate with draft system
  - [ ] Support draft tournaments (draft → deck building → matches)
  - [ ] Pass drafted card pools to players
  - [ ] Enforce deck building from drafted cards only
  - [ ] Handle draft pod completion before tournament start

## User Profile & Stats
- [x] Implement user stats tracking (matches, wins, losses, tournaments)
- [x] Implement Glicko rating calculation
- [x] Implement rating update on match completion
- [x] Implement match history display
  - [x] Add proto definitions for GetMatchHistory and GetMatchById RPCs
  - [x] Query match_history repository by user ID with pagination
  - [x] Return paginated match list (players, winner, result, date, game type, duration)
  - [x] Include total count for pagination
  - [x] Implement GetMatchById with replay data support
- [ ] Implement user profile page data
  - [ ] Add proto definition for GetUserProfile RPC
  - [ ] Return user stats (W/L record, rating, tournaments played)
  - [ ] Include achievement data (if implemented)
  - [ ] Return recent match history
  - [ ] Include deck collection count

## Real-Time Updates & Streaming
- [ ] Implement gRPC streaming for lobby updates
  - [ ] Add StreamLobbyUpdates RPC with server streaming
  - [ ] Stream table created/updated/removed events
  - [ ] Stream user joined/left events
  - [ ] Stream finished match notifications
  - [ ] Support subscription filters (by room ID)
- [ ] Implement gRPC streaming for table updates
  - [ ] Add StreamTableUpdates RPC with server streaming
  - [ ] Stream player joined/left events
  - [ ] Stream deck submitted events
  - [ ] Stream table state changes
  - [ ] Stream chat messages for table
- [ ] Implement gRPC streaming for game state updates
  - [ ] Add StreamGameUpdates RPC with server streaming
  - [ ] Stream game state changes (phase, priority, stack)
  - [ ] Stream zone updates (battlefield, graveyard, hand)
  - [ ] Stream combat events
  - [ ] Replace/complement WebSocket with gRPC streaming
- [ ] Implement event filtering per client subscription
  - [ ] Add subscription management per session
  - [ ] Filter events by room/table/game ID
  - [ ] Support multiple concurrent subscriptions
  - [ ] Unsubscribe on disconnect or explicit request
- [ ] Handle stream errors and reconnection
  - [ ] Detect broken streams and clean up
  - [ ] Resume stream from last received event ID
  - [ ] Implement exponential backoff for stream reconnection
  - [ ] Send missed events on reconnection

## Connection Management
- [x] Implement connection status tracking (session expiration tracking)
- [ ] Implement auto-reconnect with exponential backoff (client-side feature)
- [x] Implement connection health check (Ping RPC + WebSocket ping)
- [x] Implement session lease-based expiration
- [x] Implement periodic expired session cleanup
- [ ] Implement max reconnection attempts (client-side feature)
- [ ] Implement state restoration after reconnect
  - [ ] Store session state checkpoints (last game view, last event ID)
  - [ ] Add GetSessionState RPC to retrieve state after reconnect
  - [ ] Resume in-progress games after disconnect
  - [ ] Replay missed events since disconnect

## Error Handling & Recovery
- [x] Implement global error interceptor
- [x] Convert gRPC errors to user-friendly messages
- [ ] Implement retry logic for transient errors
  - [ ] Add retry interceptor for idempotent operations
  - [ ] Implement exponential backoff with jitter
  - [ ] Retry on specific gRPC codes (Unavailable, DeadlineExceeded)
  - [ ] Configure max retry attempts per operation type
- [x] Add error logging and reporting
- [ ] Handle 401/403/404/500 errors appropriately
  - [ ] Return Unauthenticated for expired/invalid sessions
  - [ ] Return PermissionDenied for admin-only operations
  - [ ] Return NotFound for missing resources with helpful messages
  - [ ] Return Internal for server errors without exposing internals

## Testing & Quality
- [ ] Write unit tests for repositories (70%+ coverage)
  - [ ] Test UserRepository CRUD operations
  - [ ] Test CardRepository search and caching
  - [ ] Test StatsRepository rating queries
  - [ ] Test DeckRepository (when implemented)
  - [ ] Test MatchHistoryRepository (when implemented)
- [ ] Write unit tests for managers and controllers
  - [ ] Test SessionManager expiration logic
  - [ ] Test UserManager authentication
  - [ ] Test TableController state machine
  - [ ] Test GameController action processing
  - [ ] Test TournamentManager pairing algorithms
  - [ ] Test DraftManager pack passing
- [x] Write integration tests for complete user flows
- [x] Write integration tests for game flow
- [x] Write integration tests for tournament flow
- [x] Write integration tests for chat system
- [ ] Set up load testing (100, 500, 1000 concurrent users)
  - [ ] Set up k6 or Locust framework
  - [ ] Create load test scenarios (login, join table, play game)
  - [ ] Test concurrent game sessions
  - [ ] Measure response times under load
  - [ ] Identify performance bottlenecks
- [ ] Test WebSocket callback delivery
  - [ ] Test event delivery to single client
  - [ ] Test broadcast to multiple clients
  - [ ] Test delivery during high load
  - [ ] Test reconnection scenarios
- [ ] Test session lifecycle (connect → ping → timeout)
  - [ ] Test session creation and validation
  - [ ] Test lease extension via Ping
  - [ ] Test automatic cleanup of expired sessions
  - [ ] Test concurrent session access

## Documentation
- [ ] Write API documentation (generated from proto)
  - [ ] Generate Markdown docs from proto files using protoc-gen-doc
  - [ ] Document all 70+ RPC methods with examples
  - [ ] Document WebSocket event types and payloads
  - [ ] Add authentication flow diagrams
  - [ ] Add game flow sequence diagrams
- [ ] Document server configuration options
  - [ ] Document all config.yaml sections
  - [ ] Document environment variable overrides
  - [ ] Add example configurations (dev, staging, prod)
  - [ ] Document database connection options
  - [ ] Document session/auth settings
- [ ] Write deployment guide
  - [ ] Docker deployment instructions
  - [ ] Kubernetes deployment manifests
  - [ ] Database setup and migration guide
  - [ ] TLS/HTTPS configuration
  - [ ] Reverse proxy setup (nginx/Caddy)
  - [ ] Production hardening checklist
- [ ] Write client integration guide
  - [ ] gRPC client setup (Go, TypeScript, Python)
  - [ ] WebSocket connection flow
  - [ ] Authentication and session management
  - [ ] Error handling best practices
  - [ ] Example client implementations
- [ ] Write troubleshooting guide
  - [ ] Common errors and solutions
  - [ ] Connection issues (gRPC, WebSocket)
  - [ ] Database connection problems
  - [ ] Session expiration handling
  - [ ] Performance tuning tips
- [ ] Document protocol changes from Java server
  - [ ] List breaking changes in RPC signatures
  - [ ] Document new features (lease-based sessions, etc.)
  - [ ] Migration guide for existing clients
  - [ ] Compatibility notes

## Deployment & Operations
- [x] Create multi-stage Dockerfile
- [x] Create Docker Compose for local development
- [x] Set up CI/CD pipeline (GitHub Actions)
- [ ] Configure Prometheus metrics scraping
  - [ ] Implement Prometheus metrics endpoint (/metrics)
  - [ ] Export RPC call counters (by method, status)
  - [ ] Export RPC latency histograms (by method)
  - [ ] Export active session count gauge
  - [ ] Export database connection pool metrics
  - [ ] Export game/table/tournament counts
  - [ ] Add Prometheus scrape configuration
- [ ] Create Grafana dashboards
  - [ ] Create server overview dashboard (requests, latency, errors)
  - [ ] Create session management dashboard
  - [ ] Create database performance dashboard
  - [ ] Create game activity dashboard (active games, players)
  - [ ] Create tournament dashboard
- [ ] Set up alerting rules
  - [ ] Alert on high error rate (>5% for 5 minutes)
  - [ ] Alert on high latency (p99 >1s for 5 minutes)
  - [ ] Alert on database connection exhaustion
  - [ ] Alert on low available sessions
  - [ ] Alert on server down/unreachable
- [x] Implement graceful shutdown (drain connections)
- [ ] Add distributed tracing (optional: Jaeger)
  - [ ] Add OpenTelemetry instrumentation
  - [ ] Trace RPC calls end-to-end
  - [ ] Trace database queries
  - [ ] Export traces to Jaeger
  - [ ] Add trace context propagation

