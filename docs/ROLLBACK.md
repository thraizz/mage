# Rollback System

## Current State

| Feature | Status |
|---------|--------|
| Engine rollback infrastructure | ✅ Complete |
| Turn snapshots (last 4 turns) | ✅ Complete |
| Action-level bookmarks | ✅ Complete |
| DB persistence for crash recovery | ✅ Complete |
| `GetMyActiveGames` RPC + UI | ✅ Complete |
| `UNDO` player action handler | ❌ Not wired |
| `CONCEDE` player action handler | ❌ Not wired |
| `PASS_TURN` player action handler | ❌ Not wired |
| `RollbackTurns` RPC endpoint | ❌ Not exposed |
| `rollback_turns_allowed` config usage | ❌ Not honored |

## Key Implementation Details

### Snapshot System
- `gameStateSnapshot` captures full game state (players, cards, zones, stack, messages)
- Turn snapshots: stored per turn in `turnSnapshots[gameID][turnNumber]`, max 4 turns
- Action bookmarks: stored in `bookmarks[gameID][]` for player undo

### Persistence Strategy
- Automatic save at turn boundaries via `SaveTurnSnapshot()`
- Game deleted from `active_games` table on completion
- Server startup runs `restoreActiveGames()` to reload all active games

### Missing Wiring in `handlePlayerAction()`
The switch statement at ~line 1345 needs cases for `UNDO`, `CONCEDE`, `PASS_TURN`.

## Related Files

### Engine
- `mage-server-go/internal/game/mage_engine.go` — Core rollback functions
- `mage-server-go/internal/game/mage_engine_test.go` — Unit tests
- `mage-server-go/internal/game/persistence_adapter.go` — DB bridge

### Repository
- `mage-server-go/internal/repository/active_games.go`
- `mage-server-go/internal/repository/models.go`
- `mage-server-go/migrations/008_create_active_games_table.up.sql`

### Proto
- `mage-server-go/api/proto/mage/v1/game.proto` — `PlayerAction` enum, `GetMyActiveGames`
- `mage-server-go/api/proto/mage/v1/models.proto` — `MatchOptions.rollback_turns_allowed`

### Client
- `mage-client-web/src/lib/api/lobby.ts` — `fetchMyActiveGames()`
- `mage-client-web/src/routes/(protected)/lobby/+page.svelte` — Active games UI
