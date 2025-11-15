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
- [~] Extend stack resolution to support triggered abilities, replacement effects, and modal choices _(trigger queue wired through event bus; replacement/modal hooks pending)_
- [ ] Implement priority windows for casting during stack resolution (e.g., mana abilities, nested responses)
- [ ] Persist stack/game events for replay and spectator synchronization
- [ ] Add comprehensive error handling and rollback when resolution fails
- [x] Implement priority retention after casting (caster retains priority by default, only passes when explicitly passing)
- [ ] Add state bookmarking and rollback mechanism for error recovery
- [ ] Implement comprehensive priority loop structure matching Java `playPriority()` pattern

## Game State & Zones
- [x] Surface battlefield/stack state via `GameGetView`
- [~] Synchronize graveyard, exile, command, and hidden zones with engine updates _(unified `moveCard` system handles all zone transitions with proper removal/addition; graveyard/exile/command zones tracked; zone change events emitted)_
- [ ] Track card ownership/controller changes (gain control, copying, phasing, etc.)
- [~] Implement continuous effects layer system (layers 1-7 per Comprehensive Rules) _(layer manager in place for basic power/toughness buffs; additional layers forthcoming)_
- [x] Handle state-based actions (lethal damage, zero loyalty, legend rule, etc.) _(life ≤ 0 auto-loss; zero/less toughness deaths; planeswalker 0 loyalty; lethal damage (damage >= toughness); legend rule (multiple legendary permanents with same name); world rule (multiple world enchantments); planeswalker uniqueness (multiple planeswalkers with same type); damage tracking system implemented)_
- [x] Support counters (loyalty, +1/+1, poison, energy, experience) _(Counter and Counters data structures; counter operations (add/remove) with event emission; boost counters integrated with layer system for power/toughness; planeswalker loyalty SBA; player counters tracked for poison/energy/experience; counter views updated across all zones)_
- [x] Provide deterministic UUID mapping for permanents, abilities, and triggers
- [x] Call `checkStateBasedActions()` before each priority (per rule 117.5)
- [x] Fix `resetPassed()` to preserve lost/left player state (`passed = loses || hasLeft()`)
- [ ] Add `canRespond()` checking in pass logic (only consider responding players in `allPassed()`)
- [ ] Ensure proper zone tracking after stack resolution (cards moved to correct zones with events)

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
- [ ] Implement triggered ability queue processing before priority (APNAP order: Active Player, Non-Active Player)
- [ ] Add `checkStateAndTriggered()` method that runs before each priority (SBA → triggers → repeat until stable)
- [ ] Handle simultaneous events between stack resolutions (process events after each resolution)

## Player Interaction & Prompts
- [x] Emit prompts when priority passes require player response
- [ ] Model blocking/attacking declarations with legality enforcement
- [ ] Support multi-choice prompts (choose mode, targets, numbers, colors) _(target selection prompts covered by "target selection flow" task in Stack & Trigger System)_
- [x] Implement mana payment flow (floating mana, cost reductions, hybrid costs) _(mana pool with regular/floating mana; cost parsing for generic, colored, X, hybrid; payment calculation and execution; cost reduction effects; floating mana empties at end of step; integrated with spell casting; comprehensive tests)_
- [ ] Recreate combat damage assignment logic (first strike, double strike, trample)
- [ ] Add concession, timeout, and match result handling aligned with rules

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
- [ ] Provide hooks for UI/websocket notifications (combat updates, triggers, log lines)
- [ ] Capture analytics/metrics for stack depth, actions per turn, average response time
- [ ] Queue triggered abilities instead of immediately pushing to stack (process via `checkTriggered()` before priority)

## Persistence, Replays & Recovery
- [ ] Store game snapshots for reconnection and spectating
- [ ] Implement replay recording/playback (step-by-step action logs)
- [ ] Ensure deterministic serialization for saved games and tournaments
- [ ] Add checksum/validation to guard against divergent game state

## Testing & Parity Validation
- [x] Add unit tests for `TurnManager` sequencing and wraparound behavior
- [x] Add unit tests for `StackManager` LIFO behavior and resolution callbacks
- [x] Extend integration tests to cover stack resolution after pass chains
- [ ] Add integration tests for combat phases (attackers/blockers, damage assignment)
- [ ] Create regression tests comparing Go vs Java engine outputs for core scenarios
- [ ] Establish rules test harness mirroring Java's JUnit suite (CR regression coverage)
- [ ] Implement fuzz/invariant tests for state-based actions and stack integrity

## Migration & Compatibility
- [ ] Provide compatibility layer for existing Java client callbacks (message equivalence)
- [ ] Translate Java replay/log formats to Go for client consumption
- [ ] Document protocol changes and migration steps for server operators
- [ ] Benchmark Go engine against Java baseline (latency, throughput, memory)

