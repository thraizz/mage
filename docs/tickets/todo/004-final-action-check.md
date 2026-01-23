# Game Engine Final Action Check & Cleanup

## Overview

This ticket verifies complete feature support in the rules-light game engine (both solo playtest and multiplayer modes use the same engine) and performs final cleanup after MageEngine removal.

**Context**: The playtest engine (`engine.go`) handles both solo playtest (client-side state) and multiplayer games (server-synced state). After removing MageEngine, we need to rename files and types for clarity.

---

## Part 1: Feature Support Verification

We need to check if we offer full support for:

## Zone & Card Movement
- [ ] Move cards between zones (Library, Hand, Battlefield, Graveyard, Exile, Command Zone, Stack, Sideboard)
- [ ] Draw X cards
- [ ] Discard specific card
- [ ] Discard random card
- [ ] Search library (with filters by type, cost, etc.)
- [ ] Shuffle library
- [ ] Mill X cards (move from library to graveyard)
- [ ] Scry/Surveil/Explore (look at top X, reorder to top/bottom/graveyard)
- [ ] Look at top X cards
- [ ] Reveal cards (temporary or persistent)
- [ ] Reveal hand
- [ ] Exile face-down
- [ ] Peek at face-down cards or opponent's hand

## Card States
- [ ] Tap/Untap (individual or all)
- [ ] Flip/Transform (double-faced cards, Day/Night)
- [ ] Face up/Face down (Morph)
- [ ] Phase in/Phase out
- [ ] Attach/Detach (Auras, Equipment, Fortifications)
- [ ] Meld (combine two cards)

## Counters on Permanents
- [ ] Add/remove +1/+1 counters
- [ ] Add/remove -1/-1 counters
- [ ] Add/remove loyalty counters (Planeswalkers)
- [ ] Add/remove keyword counters (Flying, Lifelink, Deathtouch, Shield, etc.)
- [ ] Add/remove utility counters (Charge, Oil, Lore, Time, etc.)

## Counters on Players
- [ ] Add/remove poison counters
- [ ] Add/remove energy counters
- [ ] Add/remove experience counters
- [ ] Add/remove rad counters
- [ ] Add/remove ticket counters

## Battlefield Actions
- [ ] Create tokens (custom or predefined)
- [ ] Create emblems
- [ ] Copy spells on stack
- [ ] Copy permanents on battlefield
- [ ] Change control of permanent
- [ ] Track ownership vs control

## Combat
- [ ] Declare attackers (choose creatures and targets)
- [ ] Declare blockers
- [ ] Assign combat damage
- [ ] Apply trample logic
- [ ] Track commander damage per opponent (21-damage tracker)
- [ ] Visual targeting indicators (arrows/lines)
- [ ] Mark creatures as attacking/blocking

## Player Resources & State
- [ ] Change life total
- [ ] Add mana to pool (W, U, B, R, G, C) via optional modal. Players can also just tap lands.
- [ ] Track floating mana
- [ ] Track monarch status
- [ ] Track initiative status
- [ ] Track City's Blessing
- [ ] Track Day/Night cycle

## Turn & Phase Management
- [ ] Progress through turn phases (Untap, Upkeep, Draw, Main 1, Combat, Main 2, End, Cleanup)
- [ ] Pass priority
- [ ] Pass turn
- [ ] Take extra turn
- [ ] Mulligan

## Stack & Spell Resolution
- [ ] Add spells/abilities to stack (LIFO)
- [ ] Resolve stack items
- [ ] Target objects or players
- [ ] Make modal choices

## Commander-Specific
- [ ] Command Zone management
- [ ] Track commander tax (+2 per cast)
- [ ] Commander damage matrix (track per commander to each player)
- [ ] Partner/Background support
- [ ] Dungeon progression tracking

## Randomization
- [ ] Roll dice (D6, D20, N-sided)
- [ ] Flip coins

## UI & Game Tools
- [ ] Game log/action history
- [ ] Priority indicator
- [ ] Turn/phase indicator
- [ ] Card preview on hover
- [ ] Undo functionality
- [ ] Spectator mode
- [ ] Ping/Point at cards
- [ ] Track temporary effects ("until end of turn")
- [ ] Track Storm count/spells cast this turn
- [ ] Concede

---

## Part 2: Post-MageEngine Cleanup & Renaming

After MageEngine removal, simplify naming to avoid confusion now that only one engine exists.

### Backend File Cleanup

**Delete Legacy MageEngine Files:**
- [ ] Delete `mage-server-go/internal/game/mage_engine.go` (13,786 LOC)
- [ ] Delete `mage-server-go/internal/game/mage_engine_test.go`
- [ ] Delete `mage-server-go/internal/game/engine_stack.go` (legacy stack resolution)
- [ ] Delete `mage-server-go/internal/game/engine_priority.go` (legacy priority system)
- [ ] Delete `mage-server-go/internal/game/engine_events.go` (legacy triggered abilities)
- [ ] Delete `mage-server-go/internal/game/engine_layers.go` (legacy continuous effects)
- [ ] Delete `mage-server-go/internal/game/engine_combat.go` (legacy combat system)
- [ ] Delete `mage-server-go/internal/game/null_engine.go` (if exists, was placeholder)

**Rename Rules-Light Engine Files:**
- [ ] Rename `engine.go` → `game_engine.go`
- [ ] Keep `state.go` (game state structures)
- [ ] Keep `actions.go` (game actions)
- [ ] Keep `view.go` (view filtering for hidden info)
- [ ] Keep `rollback.go` (rollback system)

### Type Name Simplification

**Remove `Engine` prefix from types in `state.go`** (now that only one engine exists):
- [ ] `EngineGameState` → `GameState`
- [ ] `EnginePlayer` → `Player`
- [ ] `EngineCard` → `Card`
- [ ] `EngineCounter` → `Counter`
- [ ] `EngineManaPool` → `ManaPool`
- [ ] `EngineLogEntry` → `LogEntry`
- [ ] `EngineNotificationHandler` → `NotificationHandler` (in game_engine.go)

**Keep these names as-is** (from mage_engine.go, used in views):
- `EngineGameView` (view struct for clients)
- `EnginePlayerView` (view struct for clients)
- `EngineCardView` (view struct for clients)
- `EngineCounterView` (view struct for clients)
- `EngineManaPoolView` (view struct for clients)

### Constructor & Interface Updates

**In `game_engine.go`:**
- [ ] Keep `NewEngine()` constructor name (returns `*Engine`)
- [ ] Update all internal references to use new type names

**In `manager.go`:**
- [ ] Remove `NewEngineAdapter()` if only wrapping one engine
- [ ] Simplify `GameEngine` interface if needed
- [ ] Update references to use new type names

### Configuration Cleanup

**In `internal/config/config.go`:**
- [ ] Remove `EngineType string` field (no longer needed)
- [ ] Remove engine type validation
- [ ] Update config docs/comments

**In `config/config.yaml` (example config):**
- [ ] Remove `engine_type: "playtest"` setting
- [ ] Update comments to reflect single engine

### Main Server Simplification

**In `cmd/server/main.go`:**
- [ ] Remove engine type selection logic (lines ~145-177)
- [ ] Replace with simple: `gameEngine := game.NewEngine(logger)`
- [ ] Remove conditional MageEngine initialization
- [ ] Remove conditional active game restoration (was MageEngine-only)
- [ ] Simplify logging to just "game engine initialized"

### Frontend Store Cleanup

**In `mage-client-web/src/lib/stores/`:**
- [ ] Rename or delete `game.legacy.ts` (old rules-enforced store)
- [ ] Ensure `multiplayer-game.ts` is primary multiplayer store
- [ ] Ensure `playtest-game.ts` is primary solo playtest store
- [ ] Update imports in components if needed

### Documentation Updates

- [ ] Update `docs/GAME_ARCHITECTURE.md` to remove dual-engine references
- [ ] Update type names in architecture docs
- [ ] Remove engine selection section
- [ ] Update file organization section with new names
- [ ] Add migration note about single-engine architecture

### Test Updates

**In `internal/integration/game_flow_test.go`:**
- [ ] Update `NewMageEngine()` → `NewEngine()` calls
- [ ] Update type assertions for new names
- [ ] Remove engine-specific test cases if any

**In other test files:**
- [ ] Update all type references
- [ ] Update constructor calls
- [ ] Verify all tests pass with new names

### Import Updates

**Across entire codebase:**
- [ ] Search for `EngineGameState` and replace with `GameState`
- [ ] Search for `EnginePlayer` and replace with `Player`
- [ ] Search for `EngineCard` and replace with `Card`
- [ ] Verify no broken imports after file renames

### Final Verification

- [ ] Run `go build ./...` (verify compilation)
- [ ] Run `go test ./...` (verify all tests pass)
- [ ] Run `npm run build` in mage-client-web (verify frontend builds)
- [ ] Start server and verify it initializes correctly
- [ ] Start a multiplayer game and verify state sync works
- [ ] Start a solo playtest and verify it works
- [ ] Check logs for any deprecated warnings

---

## Success Criteria

### Feature Support
- All checklist items verified as implemented or marked as not needed
- Missing features documented with implementation plan

### Cleanup Complete
- Zero references to MageEngine in codebase
- Single `GameEngine` type with clear naming
- No dual-engine selection logic
- All tests passing
- Both solo playtest and multiplayer modes working

### Code Quality
- Simplified architecture with single engine
- Clear, non-confusing naming conventions
- Updated documentation reflecting current state
- No legacy technical debt remaining