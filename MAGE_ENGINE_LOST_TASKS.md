# Lost MageEngine Implementation - Task Analysis

## Summary
The file `mage-server-go/internal/game/mage_engine.go` was overwritten and now only contains target-related helper methods. This document identifies which tasks from `GO_PORT_TASKS.md` are now incomplete due to this loss.

## Critical Missing Components

### 1. Core Engine Structure (CRITICAL - Blocks Everything)
**Status in GO_PORT_TASKS.md:** `[x]` Completed  
**Current Status:** ❌ **LOST - Must be reimplemented**

Missing components:
- `MageEngine` struct that implements `GameEngine` interface
- `NewMageEngine(logger *zap.Logger) *MageEngine` constructor function
- `EngineGameView` struct with fields:
  - `GameID string`
  - `State GameState` (with String() method)
  - `Phase string`
  - `Step string`
  - `Turn int`
  - `ActivePlayerID string`
  - `PriorityPlayer string`
  - `Players []EnginePlayerView`
  - `Battlefield []EngineCardView`
  - `Stack []EngineCardView`
  - `Exile []EngineCardView`
  - `Command []EngineCardView`
  - `Revealed []EngineRevealedView`
  - `LookedAt []EngineLookedAtView`
  - `Combat *EngineCombatView`
  - `StartedAt time.Time`
  - `Messages []EngineMessage`
  - `Prompts []EnginePrompt`

- `EnginePlayerView` struct with fields (from grpc_game.go):
  - `PlayerID string`
  - `Name string`
  - `Life int`
  - `Poison int`
  - `Energy int`
  - `LibraryCount int`
  - `HandCount int`
  - `Hand []EngineCardView`
  - `Graveyard []EngineCardView`
  - `ManaPool EngineManaPoolView`
  - `HasPriority bool`
  - `Passed bool`
  - `StateOrdinal int`
  - `Lost bool`
  - `Left bool`
  - `Wins int`

- `EngineCardView` struct with fields (from grpc_game.go):
  - `ID string`
  - `Name string`
  - `DisplayName string`
  - `ManaCost string`
  - `Type string`
  - `SubTypes []string`
  - `SuperTypes []string`
  - `Color string`
  - `Power string`
  - `Toughness string`
  - `Loyalty string`
  - `CardNumber int`
  - `ExpansionSet string`
  - `Rarity string`
  - `RulesText string`
  - `Tapped bool`
  - `Flipped bool`
  - `Transformed bool`
  - `FaceDown bool`
  - `Zone int`
  - `ControllerID string`
  - `OwnerID string`
  - `AttachedToCard []string`
  - `Abilities []EngineAbilityView`
  - `Counters []EngineCounterView`

- `EngineAbilityView` struct:
  - `ID string`
  - `Text string`
  - `Rule string`

- `EngineCounterView` struct:
  - `Name string`
  - `Count int`

- `EngineManaPoolView` struct:
  - `White int`
  - `Blue int`
  - `Black int`
  - `Red int`
  - `Green int`
  - `Colorless int`

- `EngineRevealedView` struct:
  - `Name string`
  - `Cards []EngineCardView`

- `EngineLookedAtView` struct:
  - `Name string`
  - `Cards []EngineCardView`

- `EngineCombatView` struct:
  - `AttackingPlayerID string`
  - `Groups []EngineCombatGroupView`

- `EngineCombatGroupView` struct:
  - `Attackers []string`
  - `Blockers []string`
  - `DefendingPlayerID string`

- `EngineMessage` struct:
  - `Text string`
  - `Color string`
  - `Timestamp time.Time`

- `EnginePrompt` struct:
  - `PlayerID string`
  - `Text string`
  - `Options []string`
  - `Timestamp time.Time`

### 2. GameEngine Interface Implementation (CRITICAL)
**Status in GO_PORT_TASKS.md:** `[x]` Completed  
**Current Status:** ❌ **LOST - Must be reimplemented**

Missing methods:
- `StartGame(gameID string, players []string, gameType string) error`
- `ProcessAction(gameID string, action PlayerAction) error`
- `GetGameView(gameID, playerID string) (interface{}, error)` - must return `*EngineGameView`
- `EndGame(gameID string, winner string) error`
- `PauseGame(gameID string) error`
- `ResumeGame(gameID string) error`

### 3. Game State Management (CRITICAL)
**Status in GO_PORT_TASKS.md:** `[x]` Completed  
**Current Status:** ❌ **LOST - Must be reimplemented**

Missing functionality:
- Game state tracking (games map, mutex protection)
- Player management within games
- Game lifecycle (start, pause, resume, end)
- Integration with `TurnManager` from `rules/turn.go`
- Integration with `StackManager` from `rules/stack.go`
- Integration with event bus from `rules/events.go`

### 4. Game View Generation (CRITICAL)
**Status in GO_PORT_TASKS.md:** `[x]` Surface battlefield/stack state via `GameGetView`  
**Current Status:** ❌ **LOST - Must be reimplemented**

Missing functionality:
- Building `EngineGameView` from current game state
- Converting internal game state to view format
- Player-specific view filtering (hand visibility, etc.)
- Battlefield state aggregation
- Stack state aggregation
- Zone state aggregation (graveyard, exile, command)
- Message log aggregation
- Prompt generation

### 5. Action Processing (CRITICAL)
**Status in GO_PORT_TASKS.md:** `[x]` Provide `MageEngine` core skeleton that tracks games, players, and actions  
**Current Status:** ❌ **LOST - Must be reimplemented**

Missing functionality:
- Action routing based on `ActionType`
- Spell casting (`SEND_STRING` actions)
- Priority passing (`PLAYER_ACTION` with "PASS")
- Integer input (`SEND_INTEGER` actions)
- Integration with stack for spell casting
- Integration with turn manager for priority
- Integration with mana system for cost payment

### 6. Integration Points (CRITICAL)
**Status in GO_PORT_TASKS.md:** Various `[x]` Completed items  
**Current Status:** ❌ **LOST - Must be reimplemented**

Missing integrations:
- Event bus integration (from `rules/events.go`)
- Watcher system integration (from `watchers/`)
- Counter system integration (from `counters/`)
- Mana system integration (from `mana/`)
- Layer system integration (from `effects/layers.go`)
- Targeting system integration (from `targeting/`)
- Legality checking integration (from `rules/legality.go`)

## Tasks Now Incomplete (Previously Completed)

Based on the overwrite, these tasks marked as `[x]` in `GO_PORT_TASKS.md` are now incomplete:

### Engine Scaffolding & Lifecycle
- ❌ **Task 10:** Provide `MageEngine` core skeleton that tracks games, players, and actions
  - **Impact:** Core engine structure completely missing
  - **Dependencies:** Blocks all other engine functionality

### Game State & Zones  
- ❌ **Task 19:** Surface battlefield/stack state via `GameGetView`
  - **Impact:** Cannot retrieve game views, breaks gRPC `GameGetView` endpoint
  - **Dependencies:** Requires EngineGameView struct and view generation logic

### Stack & Trigger System
- ❌ **Task 28:** Record log message when a stack item resolves
  - **Impact:** No message logging during stack resolution
  - **Dependencies:** Requires action processing and event integration

- ❌ **Task 29:** Auto-advance priority after resolution back to the active player
  - **Impact:** Priority system not functioning
  - **Dependencies:** Requires turn manager integration

- ❌ **Task 30:** Allow triggered abilities to be queued automatically when conditions are met
  - **Impact:** Triggered abilities not working
  - **Dependencies:** Requires event bus and watcher integration

- ❌ **Task 33:** Ensure stack legality checks (targets available, costs paid) prior to resolution
  - **Impact:** Stack legality checking not integrated
  - **Dependencies:** Requires legality checker integration

- ❌ **Task 35:** Add exhaustive integration tests covering multi-object stacks, counterspells, and priority loops
  - **Impact:** Tests exist but cannot run without engine implementation
  - **Dependencies:** Requires full engine implementation

### Player Interaction & Prompts
- ❌ **Task 38:** Emit prompts when priority passes require player response
  - **Impact:** No prompts generated for players
  - **Dependencies:** Requires prompt generation in view

- ❌ **Task 41:** Implement mana payment flow (floating mana, cost reductions, hybrid costs)
  - **Impact:** Mana system exists but not integrated into engine
  - **Dependencies:** Requires action processing integration

### Event System & Watchers
- ❌ **Task 54:** Mirror Java event bus for game events
  - **Impact:** Event bus exists but not wired into engine
  - **Dependencies:** Requires engine to emit events

- ❌ **Task 55:** Port watcher/listener infrastructure to track conditional abilities
  - **Impact:** Watchers exist but not integrated into engine
  - **Dependencies:** Requires event bus integration

### Testing & Parity Validation
- ❌ **Task 66:** Add unit tests for `TurnManager` sequencing and wraparound behavior
  - **Impact:** Tests exist but engine not available to test
  - **Dependencies:** Requires engine implementation

- ❌ **Task 67:** Add unit tests for `StackManager` LIFO behavior and resolution callbacks
  - **Impact:** Tests exist but engine not available to test
  - **Dependencies:** Requires engine implementation

- ❌ **Task 68:** Extend integration tests to cover stack resolution after pass chains
  - **Impact:** Integration tests exist but cannot run
  - **Dependencies:** Requires full engine implementation

## What Still Exists (Not Lost)

The following components are still intact and can be reused:

✅ **Supporting Systems (Still Intact):**
- `rules/turn.go` - TurnManager implementation
- `rules/stack.go` - StackManager implementation  
- `rules/events.go` - Event bus implementation
- `rules/legality.go` - LegalityChecker implementation
- `rules/trigger.go` - Trigger system
- `rules/watcher.go` - Watcher infrastructure
- `watchers/common.go` - Common watcher implementations
- `counters/` - Counter system
- `mana/` - Mana payment system
- `effects/layers.go` - Layer system
- `targeting/` - Targeting system
- `manager.go` - Game Manager (separate from engine)
- `null_engine.go` - NullEngine stub implementation

✅ **Test Files (Still Intact):**
- `mage_engine_test.go` - Unit tests for MageEngine
- `internal/integration/game_flow_test.go` - Integration tests
- `internal/integration/stack_integration_test.go` - Stack tests
- `internal/integration/watcher_*_test.go` - Watcher tests

✅ **gRPC Integration (Still Intact):**
- `internal/server/grpc_game.go` - Knows how to convert EngineGameView to proto
- `cmd/server/main.go` - Creates MageEngine instance

## Recovery Strategy

### Phase 1: Core Structure (Critical Path)
1. Recreate `MageEngine` struct with:
   - Games map with mutex protection
   - Logger reference
   - References to supporting systems (turn manager, stack manager, event bus)
2. Implement `NewMageEngine` constructor
3. Define `EngineGameView` struct with all required fields
4. Implement basic `StartGame` method (initialize game state)

### Phase 2: View Generation
1. Implement `GetGameView` method
2. Build view from game state
3. Integrate with turn manager for phase/step/turn info
4. Integrate with stack manager for stack state
5. Aggregate messages and prompts

### Phase 3: Action Processing
1. Implement `ProcessAction` method
2. Route actions by type
3. Integrate spell casting with stack
4. Integrate priority passing with turn manager
5. Integrate mana payment

### Phase 4: Integration
1. Wire event bus for game events
2. Integrate watcher system
3. Integrate counter system
4. Integrate layer system
5. Integrate targeting system
6. Integrate legality checking

### Phase 5: Lifecycle
1. Implement `EndGame` method
2. Implement `PauseGame` method
3. Implement `ResumeGame` method
4. Clean up resources

## Test Coverage Impact

All existing tests that depend on `MageEngine` will fail until reimplementation:
- `TestCardIDConsistencyAcrossZones` in `mage_engine_test.go`
- All tests in `internal/integration/game_flow_test.go`
- All tests in `internal/integration/stack_integration_test.go`
- All tests in `internal/integration/watcher_*_test.go`

## Estimated Effort

Based on the complexity and dependencies:
- **Phase 1 (Core Structure):** ~4-6 hours
- **Phase 2 (View Generation):** ~3-4 hours
- **Phase 3 (Action Processing):** ~4-6 hours
- **Phase 4 (Integration):** ~6-8 hours
- **Phase 5 (Lifecycle):** ~1-2 hours

**Total Estimated Recovery Time:** ~18-26 hours

## Notes

- The target-related methods in the current `mage_engine.go` file (`FindCardForTarget`, `FindPlayerForTarget`, `GetStackItemsForTarget`) suggest these were part of a `TargetGameStateAccessor` implementation, which may still be needed but is incomplete without the main engine structure.

- The supporting systems (turn manager, stack manager, event bus, etc.) are all intact and well-tested, so the recovery effort can focus on wiring them together rather than rebuilding them.

- The gRPC server code already knows how to handle `EngineGameView`, so once the view is generated correctly, the server integration should work immediately.
