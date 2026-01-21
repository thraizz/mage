# Mage Architecture

## System Components

```
┌─────────────────────┐
│   Web Client        │  (SvelteKit + TypeScript)
│   mage-client-web   │
└─────────┬───────────┘
          │ HTTP/gRPC-Web (actions) + WebSocket (events)
          ↓
┌─────────────────────────────────────────────┐
│   Go Server (mage-server-go)                │
│  ┌──────────────┐    ┌──────────────────┐   │
│  │ grpc.go &    │───→│ WebSocket Server │   │
│  │ grpc_table.go│    │ (websocket.go)   │   │
│  └──────┬───────┘    └────────┬─────────┘   │
│         │                     │             │
│         ↓                     ↓             │
│  ┌──────────────────────────────────────┐   │
│  │       Game Manager (manager.go)      │   │
│  └──────────────┬───────────────────────┘   │
│                 ↓                           │
│  ┌──────────────────────────────────────┐   │
│  │   MageEngine (mage_engine.go)        │   │
│  │   - Game Logic & Rules Engine        │   │
│  └──────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

## Communication Protocols

- **Client → Server**: gRPC-Web over HTTP POST to `/mage.v1.MageServer/{Method}`
- **Server → Client**: WebSocket push for real-time game updates
- **State Reading**: WebSocket events primary, `GameGetView` gRPC as fallback

## Server as Source of Truth

The server owns all game logic. The client is a thin presentation layer.

**Server provides:**
- `CardView.availableActions[]` - Per-card actions with `isEnabled` and `disabledReason`
- `GameView.activePlayerName`, `priorityPlayerName`, `gameFormat` - Pre-computed display values
- `GameView.landsPlayedThisTurn`, `landsAllowedThisTurn` - Rule state for UI feedback

**Client behavior:**
- Uses `availableActions` to determine available buttons and show disabled reasons
- Never infers game rules (no `card.type.includes('land')` checks)

## Data Flow

**Player Action:**
```
Web Client → HTTP/gRPC → grpc.go → Manager → MageEngine.ProcessAction()
```

**State Update:**
```
MageEngine.emitNotification() → handleGameNotification() → WebSocket → Web Client
```

## Key Components

| Component        | Location                          | Responsibility                                     |
| ---------------- | --------------------------------- | -------------------------------------------------- |
| `grpc_table.go`  | `internal/server/`                | Table/lobby management (pre-game)                  |
| `manager.go`     | `internal/game/`                  | Game lifecycle + `GameEngine` interface            |
| `mage_engine.go` | `internal/game/`                  | MTG rules engine + state management                |
| `websocket.go`   | `internal/server/`                | Real-time push to clients                          |
| `game.ts`        | `src/lib/stores/`                 | Game state store, WebSocket event subscription     |
| `game.ts`        | `src/lib/api/`                    | Game API functions                                 |

## GameEngine Interface

```go
type GameEngine interface {
    StartGameWithDecks(gameID string, players []string, gameType string, decks map[string]DeckList) error
    ProcessAction(gameID string, action PlayerAction) error
    GetGameView(gameID, playerID string) (interface{}, error)
    EndGame(gameID string, winner string) error
}
```

Engine emits notifications via `NotificationHandler` callback, which the server routes to WebSocket clients.
