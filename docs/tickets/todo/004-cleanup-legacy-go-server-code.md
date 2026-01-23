# Complete Legacy Go Server Code Cleanup - Full Migration

## Status: READY
**Created**: January 23, 2026
**Priority**: HIGH
**Estimated Effort**: 2-3 days

## Problem

After completing the playtest-first migration, the Go server contains massive amounts of dead code from the old MageEngine:
- **~100+ files** of legacy rules enforcement code
- **~50,000+ LOC** of unused complexity
- MageEngine (13,786 lines) no longer needed
- Complex rules system (priority, stack, combat, triggers)

**Decision**: NO backward compatibility needed. Full migration to playtest engine only.

## Solution

**Complete deletion** of all MageEngine code and legacy systems. Simplify to single playtest engine only.

---

## Cleanup Strategy - Aggressive Full Migration

### Phase 1: Delete MageEngine Core (IMMEDIATE)

**Primary target**: Remove the massive legacy engine file.

**Files to delete**:
1. `internal/game/mage_engine.go` (13,786 lines, 430KB)
2. `internal/game/engine_priority.go`
3. `internal/game/engine_stack.go`
4. `internal/game/engine_combat.go`
5. `internal/game/engine_events.go`
6. `internal/game/engine_layers.go`
7. `internal/game/null_engine.go` (test stub)

**Action**:
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go

# Delete MageEngine and all adapter files
rm internal/game/mage_engine.go
rm internal/game/engine_priority.go
rm internal/game/engine_stack.go
rm internal/game/engine_combat.go
rm internal/game/engine_events.go
rm internal/game/engine_layers.go
rm internal/game/null_engine.go
```

**Verification**:
```bash
go build ./...
# Will fail - expected, fix imports in Phase 2
```

**Impact**: ~14,000 LOC deleted

---

### Phase 2: Remove Engine Selection Logic

**Remove all code related to choosing between engines.**

**Files to modify**:

#### File: `cmd/server/main.go`

**Delete** (lines ~145-177):
```go
// Remove entire engine_type switch
engineType := cfg.Server.EngineType
if engineType == "" {
    engineType = "mage"
}

switch engineType {
case "playtest":
    // ...
case "mage":
    // ...
}
```

**Replace with** (simple, single engine):
```go
// Create playtest engine (only engine)
engine := game.NewEngine(logger, cardRepo)
engine.SetNotificationHandler(notificationAdapter)

gameAdapter := game.NewEngineAdapter(engine)
```

**Delete** (lines ~169-192):
```go
// Remove MageEngine persistence/crash recovery
if mageEngine, ok := engine.(*game.MageEngine); ok {
    // ... persistence logic
}

// Remove game restoration
activeGames, err := activeGameRepo.ListAll(ctx)
// ... restoration logic
```

#### File: `internal/config/config.go`

**Delete** (line 33):
```go
EngineType string `mapstructure:"engine_type"`
```

**Delete** (line 205):
```go
v.SetDefault("server.engine_type", "mage")
```

**Delete** (lines 285-287):
```go
// Validate engine type
if c.Server.EngineType != "" && c.Server.EngineType != "mage" && c.Server.EngineType != "playtest" {
    return fmt.Errorf("server.engine_type must be 'mage' or 'playtest'")
}
```

#### File: `internal/game/manager.go`

**Simplify** `SetNotificationCallback`:
```go
func (a *EngineAdapter) SetNotificationCallback(callback func(GameNotification)) {
    if e, ok := a.engine.(*Engine); ok {
        // Only support playtest Engine
        adapter := &engineNotificationAdapter{callback: callback}
        e.SetNotificationHandler(adapter)
    } else {
        // Should never happen - only one engine type now
        panic("unknown engine type")
    }
}
```

Or even simpler - assume it's always Engine:
```go
func (a *EngineAdapter) SetNotificationCallback(callback func(GameNotification)) {
    e := a.engine.(*Engine)
    adapter := &engineNotificationAdapter{callback: callback}
    e.SetNotificationHandler(adapter)
}
```

**Verification**:
```bash
go build ./cmd/server/...
# Should compile now
```

**Impact**: ~100 LOC simplified

---

### Phase 3: Delete Entire Rules System

**Remove all MTG rules enforcement infrastructure.**

**Directories to delete**:
```bash
rm -rf internal/game/rules/
rm -rf internal/game/effects/
```

**Files deleted**:
- `rules/` directory: 21 files (~18,000 LOC)
  - priority.go, stack.go, state_based_actions.go
  - trigger.go, turn.go, watcher.go, legality.go
  - payment_window.go, special_action.go, mana_ability.go
  - events.go, and 10+ more
- `effects/` directory: 17 files (~4,000 LOC)
  - Layer system, replacement effects, prevention effects

**Verification**:
```bash
# Check for any remaining imports
grep -r "game/rules\|game/effects" mage-server-go/ --include="*.go"
# Expected: No results (all deleted)
```

**Impact**: ~22,000 LOC deleted

---

### Phase 4: Clean Up Abilities System

**Delete complex ability implementations, keep minimal infrastructure.**

**Files to delete** (from `internal/game/abilities/`):
```bash
cd internal/game/abilities/

# Delete complex implementations
rm triggered.go                 # Triggered abilities
rm cost_modification.go         # Cost reduction
rm dynamic_value*.go           # Dynamic calculations
rm counter_effects.go           # Counter interactions
rm modal.go                     # Modal abilities
rm kicker.go                    # Kicker mechanics
rm combat_damage.go             # Combat damage
rm combat_restrictions.go       # Combat restrictions
rm combat_special.go            # Special combat
rm keyword_abilities_combat.go  # Combat keywords
```

**Files to keep**:
- `ability.go` - Base structures
- `activated.go` - Simple activated abilities
- `static.go` - Static abilities
- `costs.go` - Cost structures
- `targets.go` - Target structures

**Verification**:
```bash
ls internal/game/abilities/
# Should only show ~5-8 files kept
```

**Impact**: ~3,000 LOC deleted

---

### Phase 5: Delete All Combat Tests

**Remove entire combat test suite.**

```bash
# Delete all combat test files
find internal/game -name "*combat*test.go" -delete

# Delete integration tests for old engine
rm internal/integration/stack_integration_test.go
rm internal/integration/watcher_game_engine_test.go
rm internal/integration/watcher_integration_test.go
```

**Files deleted**: ~45 test files (~8,500 LOC)

**Verification**:
```bash
# Run remaining tests
go test ./internal/game/... -v
go test ./internal/integration/... -v
# Should pass with only playtest-related tests
```

**Impact**: ~8,500 LOC deleted

---

### Phase 6: Delete Persistence Layer for MageEngine

**Remove database persistence and serialization (MageEngine only).**

**Files to delete**:
```bash
rm internal/game/persistence_adapter.go
rm internal/game/serialization.go
rm internal/repository/active_games.go
```

**Rationale**:
- Playtest engine is in-memory only
- No need for complex game state serialization
- No crash recovery system needed

**Verification**:
```bash
grep -r "persistence_adapter\|serialization\|active_games" internal/ --include="*.go"
# Expected: No results
```

**Impact**: ~1,000 LOC deleted

---

### Phase 7: Delete Unused Features (If Not Implemented)

**Check and delete if not actually used.**

#### Draft System
```bash
# Check if used
grep -r "draft" mage-client-web/src --include="*.ts" --include="*.svelte"

# If no usage, delete
rm -rf internal/draft/
rm internal/server/grpc_draft.go
```

#### Tournament System
```bash
# Check if used
grep -r "tournament" mage-client-web/src --include="*.ts" --include="*.svelte"

# If no usage, delete
rm -rf internal/tournament/
rm internal/server/grpc_tournament.go
```

#### Replay System
```bash
# Check if replay handlers exist
grep -r "replay\|Replay" internal/server --include="*.go"

# If not used, delete
rm internal/game/replay.go
rm internal/game/rollback.go  # Old rollback, not new engine's
```

**Verification**: Check each feature individually before deleting.

**Impact**: ~2,000 LOC deleted (if all unused)

---

### Phase 8: Update Server Handler

**Simplify gRPC conversion now that only one engine exists.**

#### File: `internal/server/grpc.go`

**Simplify** `engineViewToProto`:
```go
func (s *MageServer) engineViewToProto(view interface{}) (*pb.GameView, error) {
    // Only handle PlaytestGameView now
    playtestView, ok := view.(*game.PlaytestGameView)
    if !ok {
        return nil, fmt.Errorf("expected PlaytestGameView, got %T", view)
    }

    return s.playtestViewToProto(playtestView), nil
}
```

**Delete**: All EngineGameView handling code (lines ~646-691)

**Verification**:
```bash
go build ./internal/server/...
```

**Impact**: ~50 LOC simplified

---

### Phase 9: Final Cleanup and Verification

**Remove any remaining references and verify system works.**

#### Cleanup Tasks

1. **Search for MageEngine references**:
```bash
grep -r "MageEngine\|mage_engine" mage-server-go/ --include="*.go"
# Fix or delete any remaining references
```

2. **Search for engine_type references**:
```bash
grep -r "engine_type\|EngineType" mage-server-go/ --include="*.go"
# Should only be in config structs (can leave for legacy config files)
```

3. **Remove unused imports**:
```bash
# Many files will have unused imports after deletions
go fmt ./...
go build ./...
# Fix any import errors
```

4. **Update proto definitions** (if needed):
```bash
# Check if any proto messages reference old engine
grep -r "mage.*engine" api/proto/
# Update if necessary
```

#### Verification Checklist

- [ ] All MageEngine files deleted
- [ ] Engine selection logic removed from main.go
- [ ] Config simplified (no engine_type validation)
- [ ] Rules system deleted (rules/, effects/)
- [ ] Abilities system cleaned (complex implementations removed)
- [ ] Combat tests deleted
- [ ] Persistence layer removed
- [ ] Unused features deleted (draft/tournament/replay if not used)
- [ ] Server handlers simplified
- [ ] All tests pass: `go test ./...`
- [ ] Server compiles: `go build ./cmd/server/...`
- [ ] Server runs: `./server` starts without errors

---

## Cleanup Summary

### Total Deletions

| Category | Files | LOC |
|----------|-------|-----|
| MageEngine core | 7 | ~14,000 |
| Rules system | 21 | ~18,000 |
| Effects system | 17 | ~4,000 |
| Abilities (complex) | 10 | ~3,000 |
| Combat tests | 45 | ~8,500 |
| Persistence | 3 | ~1,000 |
| Unused features | 0-6 | 0-2,000 |
| **TOTAL** | **103-109** | **~48,500-50,500** |

### Before Cleanup
- Total files: ~241 Go files
- Backend LOC: ~85,000
- Engine options: 2 (MageEngine, PlaytestEngine)

### After Cleanup
- Total files: ~132-138 Go files
- Backend LOC: ~34,500-36,500
- Engine options: 1 (PlaytestEngine only)

**Code reduction: ~57% smaller codebase**

---

## Breaking Changes

### API Changes
- ❌ **Removed**: `engine_type` config option
- ❌ **Removed**: MageEngine support
- ❌ **Removed**: Game state persistence/serialization
- ❌ **Removed**: Crash recovery system
- ❌ **Removed**: Complex rules enforcement

### Migration Path
**NONE** - This is a clean break. Old games cannot be restored.

**Impact**: All existing games will be lost. This is acceptable since:
1. This is a single-user development project
2. Playtest-first is the new primary architecture
3. No production users to migrate

---

## Implementation Order

**Aggressive 2-3 day cleanup:**

### Day 1: Core Deletions
1. ✅ Phase 1: Delete MageEngine (will break build)
2. ✅ Phase 2: Remove engine selection logic (fixes build)
3. ✅ Phase 3: Delete rules/ and effects/ directories
4. ✅ Phase 5: Delete all combat tests
5. ✅ Verify: `go test ./...` passes

### Day 2: System Cleanup
6. ✅ Phase 4: Clean abilities system
7. ✅ Phase 6: Delete persistence layer
8. ✅ Phase 7: Delete unused features (verify each)
9. ✅ Verify: `go build ./...` succeeds

### Day 3: Final Polish
10. ✅ Phase 8: Simplify server handlers
11. ✅ Phase 9: Final cleanup and verification
12. ✅ Update documentation
13. ✅ Test complete system end-to-end

---

## Benefits

### Immediate
1. **Simplicity**: Single engine, clear architecture
2. **Code reduction**: 57% smaller codebase
3. **Clarity**: No confusion about which engine to use
4. **Performance**: Faster compilation, smaller binary

### Long-term
1. **Maintainability**: Far less code to maintain
2. **Developer experience**: New devs see only relevant code
3. **Feature velocity**: No legacy code slowing development
4. **Clean foundation**: Ready for future features

---

## Risks & Mitigation

### Risk 1: Breaking something critical
**Mitigation**:
- Run tests after each phase
- Keep git commits granular
- Can rollback individual phases if needed

### Risk 2: Losing functionality we need
**Mitigation**:
- Phase 7 explicitly verifies feature usage before deletion
- Draft/tournament only deleted if confirmed unused

### Risk 3: Frontend breaks
**Mitigation**:
- Frontend already uses multiplayer-game store (not old game store)
- gRPC API unchanged (still returns GameView)
- Only backend implementation changes

---

## Documentation Updates

After cleanup, update these docs:

1. **GAME_ARCHITECTURE.md**:
   - Remove all MageEngine references
   - Document single-engine architecture
   - Update diagrams

2. **README.md**:
   - Remove engine selection instructions
   - Simplify server configuration
   - Update feature list

3. **PLAYTEST_MIGRATION_SUMMARY.md**:
   - Add "Complete Cleanup" section
   - Document final code reduction numbers
   - Mark migration as 100% complete

4. **config/config.yaml**:
   - Remove `engine_type` option
   - Update comments

---

## Success Criteria

- ✅ MageEngine completely removed
- ✅ Single playtest engine only
- ✅ ~50,000 LOC deleted
- ✅ All tests pass
- ✅ Server compiles and runs
- ✅ Frontend works with new backend
- ✅ No backward compatibility code
- ✅ Documentation updated
- ✅ Clean git history with granular commits

---

## Post-Cleanup Verification

### Build Verification
```bash
cd mage-server-go
go build ./...
go test ./...
./server
```

### Frontend Verification
```bash
cd mage-client-web
bun run check
bun run dev
# Test game creation and playtest UI
```

### End-to-End Test
1. Start server: `./server`
2. Start frontend: `bun run dev`
3. Create 2-player game
4. Test all operations:
   - Draw card
   - Play card
   - Tap/untap
   - Move to zones
   - Life changes
   - Token creation
5. Verify state sync between players

---

## Notes

**Priority**: HIGH - This is now top priority for clean architecture.

**No backward compatibility**: This is a full migration. Accept that old games are lost.

**Git strategy**:
- Commit after each phase
- Descriptive commit messages
- Easy to rollback individual phases if needed

**Timeline**: 2-3 days of focused work

**After completion**: Codebase will be 57% smaller, dramatically simpler, and fully aligned with playtest-first architecture.

---

## Related Tickets

- ✅ `003-implement-direct-engine.md` - Original playtest migration (DONE)
- 🔄 `004-cleanup-legacy-go-server-code.md` - This ticket (IN PROGRESS)

---

## Execution Checklist

Day 1:
- [ ] Phase 1: Delete MageEngine files
- [ ] Phase 2: Simplify main.go and config
- [ ] Phase 3: Delete rules/ and effects/
- [ ] Phase 5: Delete combat tests
- [ ] Verify tests pass

Day 2:
- [ ] Phase 4: Clean abilities/
- [ ] Phase 6: Delete persistence
- [ ] Phase 7: Delete unused features
- [ ] Verify build succeeds

Day 3:
- [ ] Phase 8: Simplify handlers
- [ ] Phase 9: Final cleanup
- [ ] Update docs
- [ ] End-to-end testing
- [ ] Git commits and push

**Status after completion**: Production-ready, fully migrated playtest-first architecture with zero legacy code.
