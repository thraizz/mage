# Mage Architecture

## Overview

This document describes how the key components of the Mage system connect and communicate.

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
│                 │                           │
│                 ↓                           │
│  ┌──────────────────────────────────────┐   │
│  │     Engine Adapter (EngineAdapter)   │   │
│  └──────────────┬───────────────────────┘   │
│                 │                           │
│                 ↓                           │
│  ┌──────────────────────────────────────┐   │
│  │   MageEngine (mage_engine.go)        │   │
│  │   - Game Logic & Rules Engine        │   │
│  └──────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

## Communication Flow

### 1. Web Client → Server (Actions)

The web client sends game actions via **gRPC-Web over HTTP**:

```typescript
// mage-client-web/src/lib/api/game.ts
export async function sendPlayerAction(
  gameId: string,
  action: PlayerAction
): Promise<void> {
  const client = getMageClient();
  const sessionId = await client.ensureSessionId();

  const request: SendPlayerActionRequest = {
    sessionId,
    gameId,
    action,
  };

  const response = await client.call<
    SendPlayerActionRequest,
    SendPlayerActionResponse
  >("SendPlayerAction", request);
}
```

The client makes HTTP POST requests to endpoints like `/mage.v1.MageServer/SendPlayerAction`.

### 2. grpc_table.go - Table/Lobby Management

`grpc_table.go` handles table-level operations like joining tables, submitting decks, and starting games:

```go
// mage-server-go/internal/server/grpc_table.go
func (s *mageServer) RoomJoinTable(ctx context.Context, req *pb.RoomJoinTableRequest) (*pb.RoomJoinTableResponse, error) {
    // ... validation ...
    if err := tbl.AddPlayer(username, playerType); err != nil {
        // ...
    }
    // Parse and submit deck if provided
    if deckListText != "" {
        tbl.SubmitDeck(username, deckList)
    }
}
```

It acts as the **entry point** for table/lobby operations before games start.

### 3. Game Manager (manager.go) - Game Lifecycle

The `Manager` in `manager.go` tracks active games and routes actions to the engine:

```go
// mage-server-go/internal/game/manager.go
type Manager struct {
    games        map[string]*Game
    gamesByTable map[string]string // tableID -> gameID
    mu           sync.RWMutex
    logger       *zap.Logger
}

func (m *Manager) CreateGame(tableID, gameType string, players []string) *Game {
    game := NewGame(tableID, gameType, players)
    m.games[game.ID] = game
    m.gamesByTable[tableID] = game.ID
    return game
}
```

It also defines the **GameEngine interface** that `MageEngine` implements:

```go
type GameEngine interface {
    StartGame(gameID string, players []string, gameType string) error
    StartGameWithDecks(gameID string, players []string, gameType string, decks map[string]DeckList) error
    ProcessAction(gameID string, action PlayerAction) error
    GetGameView(gameID, playerID string) (interface{}, error)
    EndGame(gameID string, winner string) error
    PauseGame(gameID string) error
    ResumeGame(gameID string) error
}
```

### 4. MageEngine (mage_engine.go) - The Rules Engine

`MageEngine` is the **core game rules engine** that implements Magic: The Gathering rules:

```go
// mage-server-go/internal/game/mage_engine.go
type MageEngine struct {
    logger              *zap.Logger
    mu                  sync.RWMutex
    games               map[string]*engineGameState
    notificationHandler NotificationHandler // Optional handler for UI/websocket notifications
    // ... bookmarks, turn rollback, replay recording ...
    replacementEffects map[string]*effects.ReplacementManager
}

func (e *MageEngine) SetNotificationHandler(handler NotificationHandler) {
    e.mu.Lock()
    defer e.mu.Unlock()
    e.notificationHandler = handler
}
```

When processing actions, it routes by type:

```go
func (e *MageEngine) ProcessAction(gameID string, action PlayerAction) (err error) {
    // ... get game state and create bookmark for error recovery ...

    switch action.ActionType {
    case "PLAYER_ACTION":
        return e.handlePlayerAction(gameState, action)
    case "SEND_STRING":
        return e.handleStringAction(gameState, action)
    case "SEND_UUID":
        return e.handleUUIDAction(gameState, action)
    // ... other action types ...
    }
}
```

### 5. Engine → Server → WebSocket (Notifications)

When game state changes, `MageEngine` emits notifications via a callback:

```go
// mage-server-go/internal/game/mage_engine.go
func (e *MageEngine) emitNotification(notification GameNotification) {
    e.mu.RLock()
    handler := e.notificationHandler
    e.mu.RUnlock()

    if handler != nil {
        go handler(notification)  // Async to avoid blocking game logic
    }
}
```

The server wires this up in `grpc.go`:

```go
// mage-server-go/internal/server/grpc.go
func (s *mageServer) SetupGameNotifications() {
    s.gameAdapter.SetNotificationCallback(func(notification game.GameNotification) {
        s.handleGameNotification(notification)
    })
}

func (s *mageServer) handleGameNotification(notification game.GameNotification) {
    // Send game update to all players
    for _, playerName := range gameInstance.Players {
        s.sendGameUpdateToPlayer(gameID, playerName)
    }
    // Also send to watchers
    for _, watcher := range gameInstance.GetWatchers() {
        s.sendGameUpdateToPlayer(gameID, watcher)
    }
}
```

The WebSocket server in `websocket.go` then pushes events to connected clients:

```go
// mage-server-go/internal/server/websocket.go
for {
    select {
    case event, ok := <-sess.CallbackChan:
        if err := ws.sendEvent(conn, event); err != nil {
            // handle error
        }
    case <-done:
        return
    }
}
```

### 6. WebSocket → Web Client (Real-time Updates)

The web client connects via WebSocket and receives updates:

```typescript
// mage-client-web/src/lib/stores/websocket.ts
function connect(newSessionId: string): Promise<void> {
  const url = `${WEBSOCKET_URL}?sessionId=${encodeURIComponent(sessionId)}`;
  ws = new WebSocket(url);

  ws.onmessage = (event) => {
    // Decode and route message by type (GAME_UPDATE, CHAT_MESSAGE, etc.)
  };
}
```

### 7. Frontend Game State Reading

The frontend uses a **hybrid approach** to read game state:

**Primary: WebSocket Events (Real-time)**

- On game initialization, the client: (1) connects WebSocket, (2) calls `joinGame()` via gRPC, (3) subscribes to WebSocket events
- The game store (`game.ts`) subscribes to events: `GAME_INIT`, `GAME_UPDATE`, `GAME_UPDATE_AND_INFORM`
- Each event contains the full `GameView` which updates the Svelte store

```typescript
// mage-client-web/src/lib/stores/game.ts
websocketStore.on(CallbackMethod.GAME_UPDATE, (data) => {
  const normalized = normalizeGameView(updateData.game);
  update((s) => ({ ...s, gameView: normalized }));
});
```

**Fallback: HTTP/gRPC (Initial Load)**

- `fetchGameView()` calls `GameGetView` gRPC method to fetch initial state
- Used as fallback if WebSocket `GAME_INIT` event is missed

```typescript
// mage-client-web/src/lib/api/game.ts
export async function fetchGameView(
  gameId: string,
  playerId?: string
): Promise<GameView> {
  const response = await client.call<GameGetViewRequest, GameGetViewResponse>(
    "GameGetView",
    { sessionId, gameId, playerId }
  );
  return response.game;
}
```

## Data Flow Summary

### Player Action Flow

```
Web Client → HTTP/gRPC → grpc.go → EngineAdapter → MageEngine.ProcessAction()
```

### State Update Flow

```
MageEngine.emitNotification() → handleGameNotification() → WebSocket → Web Client
```

### Game State Reading Flow

```
Web Client: Connect WebSocket → joinGame() → Subscribe to GAME_UPDATE events
                                 ↓
                            fetchGameView() (fallback)
                                 ↓
                         gameStore updates GameView
```

### Key Components

| Component        | Location                          | Responsibility                                     |
| ---------------- | --------------------------------- | -------------------------------------------------- |
| `grpc_table.go`  | `mage-server-go/internal/server/` | Table/lobby management (pre-game)                  |
| `manager.go`     | `mage-server-go/internal/game/`   | Game lifecycle management + `GameEngine` interface |
| `mage_engine.go` | `mage-server-go/internal/game/`   | MTG rules implementation + state management        |
| `EngineAdapter`  | `mage-server-go/internal/game/`   | Bridges manager and engine, handles notifications  |
| `websocket.go`   | `mage-server-go/internal/server/` | Real-time push channel to clients                  |
| `client.ts`      | `mage-client-web/src/lib/grpc/`   | gRPC-Web client for server communication           |
| `websocket.ts`   | `mage-client-web/src/lib/stores/` | WebSocket store for real-time updates              |
| `game.ts`        | `mage-client-web/src/lib/stores/` | Game state store, subscribes to WebSocket events   |
| `game.ts` (API)  | `mage-client-web/src/lib/api/`    | Game API functions (`fetchGameView`, `joinGame`)   |

```

---
```
