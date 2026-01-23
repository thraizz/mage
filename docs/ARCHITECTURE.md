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
│  │ grpc_game.go │    │ (websocket.go)   │   │
│  └──────┬───────┘    └────────┬─────────┘   │
│         │                     │             │
│         ↓                     ↓             │
│  ┌──────────────────────────────────────┐   │
│  │       Game Manager (manager.go)      │   │
│  └──────────────┬───────────────────────┘   │
│                 ↓                           │
│  ┌──────────────────────────────────────┐   │
│  │   GameEngine (game_engine.go)        │   │
│  │   - Rules-Light Game Logic           │   │
│  └──────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

## Communication Protocols

- **Client → Server**: gRPC-Web over HTTP POST to `/mage.v1.MageServer/{Method}`
- **Server → Client**: WebSocket push for real-time game updates
- **State Reading**: WebSocket events primary, `GameGetView` gRPC as fallback

## Server as Source of Truth

The server owns game state. The client is a player-controlled interface.

**Server provides:**
- Game state synchronization across all players
- Hidden information filtering (opponent hands/libraries)
- Action logging and rollback capability
- Turn and phase tracking (cosmetic, not enforced)

**Client behavior:**
- Players have direct control over game state
- No automatic rules enforcement
- Manual combat and spell resolution
- Flexible, Untap.in-style gameplay

## Data Flow

**Player Action:**
```
Web Client → HTTP/gRPC → grpc_game.go → Manager → GameEngine.ProcessAction()
```

**State Update:**
```
GameEngine.emitNotification() → handleGameNotification() → WebSocket → Web Client
```

## Key Components

| Component          | Location                          | Responsibility                                     |
| ------------------ | --------------------------------- | -------------------------------------------------- |
| `grpc.go`          | `internal/server/`                | Server initialization                              |
| `grpc_game.go`     | `internal/server/`                | Game gRPC handlers                                 |
| `manager.go`       | `internal/game/`                  | Game lifecycle management                          |
| `game_engine.go`   | `internal/game/`                  | Rules-light game engine + state management         |
| `websocket.go`     | `internal/server/`                | Real-time push to clients                          |
| `multiplayer-game.ts` | `src/lib/stores/`              | Multiplayer game state store                       |
| `direct-actions.ts`| `src/lib/api/`                   | Direct game actions API                            |

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
