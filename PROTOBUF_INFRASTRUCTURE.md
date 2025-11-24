# Protocol Buffers Infrastructure

This document describes the complete gRPC/Protocol Buffers infrastructure for the MAGE server and client.

## Overview

The MAGE project uses **Protocol Buffers (protobuf)** for type-safe, efficient communication between the Go server and web clients. The same `.proto` files are used to generate both:

- **Go server code** - gRPC service implementations
- **TypeScript client code** - Type-safe client bindings

## Architecture

### Hybrid gRPC + WebSocket Design

- **gRPC**: All 60+ request/response RPC methods (authentication, game operations, etc.)
- **WebSocket**: Server-to-client push events (real-time game updates, chat messages)
- **Protocol Buffers**: Type-safe serialization for both protocols

### Proto File Organization

Location: `mage-server-go/api/proto/mage/v1/`

| File | Description | Lines |
|------|-------------|-------|
| `server.proto` | Main MageServer service definition (60+ RPC methods) | 250 |
| `auth.proto` | Authentication messages (login, register, password reset) | 122 |
| `room.proto` | Lobby and room management | 67 |
| `table.proto` | Table creation, joining, deck submission | 175 |
| `game.proto` | Game execution, player actions, replay | 212 |
| `tournament.proto` | Tournament management | 50 |
| `draft.proto` | Booster draft mechanics | 61 |
| `chat.proto` | Chat rooms and messaging | 75 |
| `admin.proto` | Admin operations (user management, moderation) | 101 |
| `models.proto` | Shared data models (GameView, TableView, etc.) | 319 |
| `websocket.proto` | WebSocket event definitions (86 event types) | 295 |

**Total: 11 proto files, 1,727 lines**

## RPC Method Coverage

The `server.proto` file defines all 60+ RPC methods across 9 categories:

### 1. Authentication & Connection (7 methods)
- `AuthRegister` - Register new user
- `AuthSendTokenToEmail` - Send password reset token
- `AuthResetPassword` - Reset password with token
- `ConnectUser` - User login
- `ConnectAdmin` - Admin login
- `ConnectSetUserData` - Set user preferences
- `Ping` - Keep session alive (lease renewal)

### 2. Server Info (3 methods)
- `GetServerState` - Get server status (players, games, etc.)
- `ServerGetPromotionMessages` - Get MOTD/announcements
- `ServerAddFeedbackMessage` - Submit bug reports

### 3. Room/Lobby (5 methods)
- `ServerGetMainRoomId` - Get lobby ID
- `RoomGetUsers` - Get online players
- `RoomGetFinishedMatches` - Get recent match history
- `RoomGetAllTables` - List all tables
- `RoomGetTableById` - Get specific table details

### 4. Table Management (10 methods)
- `RoomCreateTable` - Create new game table
- `RoomCreateTournament` - Create tournament
- `RoomJoinTable` - Join table
- `RoomJoinTournament` - Join tournament
- `RoomLeaveTableOrTournament` - Leave table/tournament
- `RoomWatchTable` - Spectate table
- `RoomWatchTournament` - Spectate tournament
- `TableSwapSeats` - Change seat
- `TableRemove` - Close table (host only)
- `TableIsOwner` - Check host status

### 5. Deck Management (2 methods)
- `DeckSubmit` - Submit deck for table/tournament
- `DeckSave` - Save deck to collection

### 6. Game Execution (12 methods)
- `GameJoin` - Join game as player
- `GameWatchStart` - Start spectating
- `GameWatchStop` - Stop spectating
- `GameGetView` - Get current game state
- `SendPlayerUUID` - Select card/permanent/spell
- `SendPlayerString` - Send text input
- `SendPlayerBoolean` - Send yes/no choice
- `SendPlayerInteger` - Send number
- `SendPlayerManaType` - Select mana color
- `SendPlayerAction` - Send action (pass, undo, concede)
- `MatchStart` - Start match from table
- `MatchQuit` - Concede match

### 7. Draft (5 methods)
- `DraftJoin` - Join draft
- `SendDraftCardPick` - Pick card from booster
- `SendDraftCardMark` - Mark card for tracking
- `DraftSetBoosterLoaded` - Notify booster loaded
- `DraftQuit` - Leave draft

### 8. Tournament (4 methods)
- `TournamentJoin` - Join tournament
- `TournamentStart` - Start tournament (owner)
- `TournamentQuit` - Drop from tournament
- `TournamentFindById` - Get tournament details

### 9. Chat (7 methods)
- `ChatJoin` - Join chat room
- `ChatLeave` - Leave chat room
- `ChatSendMessage` - Send message
- `ChatFindByTable` - Get table chat
- `ChatFindByGame` - Get game chat
- `ChatFindByTournament` - Get tournament chat
- `ChatFindByRoom` - Get lobby chat

### 10. Replay (6 methods)
- `ReplayInit` - Initialize replay
- `ReplayStart` - Start playback
- `ReplayStop` - Stop playback
- `ReplayNext` - Next action
- `ReplayPrevious` - Previous action
- `ReplaySkipForward` - Skip N actions

### 11. Admin (9 methods)
- `AdminGetUsers` - List all users
- `AdminDisconnectUser` - Kick user
- `AdminMuteUser` - Mute in chat
- `AdminLockUser` - Temporary ban
- `AdminActivateUser` - Unlock account
- `AdminToggleActivateUser` - Toggle active status
- `AdminEndUserSession` - Force disconnect
- `AdminTableRemove` - Force close table
- `AdminSendBroadcastMessage` - Server announcement

**Total: 70 RPC methods**

## WebSocket Event Types

The `websocket.proto` file defines 86 callback event types for real-time updates:

### Categories
- **Chat & Messages** (3 events): CHATMESSAGE, SHOW_USERMESSAGE, SERVER_MESSAGE
- **Table Events** (2 events): JOINED_TABLE, TABLE_WAITING
- **Tournament Events** (4 events): START_TOURNAMENT, TOURNAMENT_INIT, TOURNAMENT_UPDATE, TOURNAMENT_OVER
- **Draft Events** (7 events): START_DRAFT, SIDEBOARD, CONSTRUCT, DRAFT_OVER, DRAFT_INIT, DRAFT_PICK, DRAFT_UPDATE
- **Watch Events** (2 events): SHOW_TOURNAMENT, WATCHGAME
- **Deck View Events** (2 events): VIEW_LIMITED_DECK, VIEW_SIDEBOARD
- **User Interaction** (2 events): USER_REQUEST_DIALOG, GAME_REDRAW_GUI
- **Game Events** (17 events): START_GAME, GAME_INIT, GAME_UPDATE, GAME_TARGET, GAME_CHOOSE_ABILITY, etc.
- **Replay Events** (4 events): REPLAY_GAME, REPLAY_INIT, REPLAY_UPDATE, REPLAY_DONE

## Data Models

The `models.proto` file defines comprehensive data structures:

### Core Models
- `TableView` - Game table representation
- `GameView` - Complete game state
- `PlayerView` - Player state (life, mana, zones)
- `CardView` - Card with abilities, counters, attachments
- `TournamentView` - Tournament state and standings
- `DraftPickView` - Draft pick state
- `UserView` - User info and statistics
- `ChatMessage` - Chat message with metadata

### Game State Components
- `CombatView` - Combat state with combat groups
- `ManaPoolView` - Mana pool (WUBRG + colorless)
- `AbilityView` - Card ability representation
- `CounterView` - Counters on permanents
- `RevealedView` / `LookedAtView` - Hidden information

## Code Generation

### Go Server Code

**Location**: `mage-server-go/pkg/proto/mage/v1/`

**Command**:
```bash
cd mage-server-go
make proto
```

**Script**: `scripts/generate_proto.sh`

**Generated files**:
- `*.pb.go` - Message structs and serialization
- `*_grpc.pb.go` - gRPC service stubs

**Requirements**:
- `protoc` - Protocol buffer compiler
- `protoc-gen-go` - Go plugin for protoc
- `protoc-gen-go-grpc` - gRPC plugin for protoc

**Install tools**:
```bash
cd mage-server-go
make tools
```

### TypeScript Client Code

**Location**: `mage-client-web/src/lib/generated/mage/v1/`

**Command**:
```bash
cd mage-client-web
npm run proto:generate
```

**Script**: `scripts/generate-proto.sh`

**Generated files**:
- `*.ts` - TypeScript interfaces, types, and client implementations
- Uses `ts-proto` for modern TypeScript generation

**Configuration**:
- `outputServices=grpc-js` - Generate grpc-js compatible service definitions
- `env=browser` - Browser-compatible code
- `useOptionals=messages` - Optional fields for message fields
- `outputIndex=true` - Generate index.ts for easier imports

**Requirements**:
- `protoc` - Protocol buffer compiler
- `ts-proto` - TypeScript plugin for protoc (installed via npm)

## Client Usage Examples

### Basic Setup

```typescript
import { getMageClient } from '$lib/grpc/client';

const client = getMageClient({
  serverUrl: 'http://localhost:17171'
});
```

### Authentication

```typescript
// Login
const response = await client.connectUser('username', 'password');
if (response.success) {
  console.log('Session ID:', response.sessionId);
}

// Register
await client.register('newuser', 'password123', 'user@example.com');

// Keep session alive
setInterval(() => client.ping(), 60000); // Ping every minute
```

### Lobby Operations

```typescript
// Get main lobby
const lobbyResponse = await client.getMainRoomId();
const roomId = lobbyResponse.roomId;

// List tables
const tablesResponse = await client.getAllTables(roomId);
console.log('Tables:', tablesResponse.tables);

// Get online users
const usersResponse = await client.rpc.RoomGetUsers({
  sessionId: client.getSessionId()!,
  roomId
});
```

### Game Operations

```typescript
// Get game state
const gameView = await client.getGameView(gameId, playerId);
console.log('Turn:', gameView.game?.turn);
console.log('Phase:', gameView.game?.phase);

// Send player action
await client.rpc.SendPlayerAction({
  sessionId: client.getSessionId()!,
  gameId,
  action: PlayerAction.PASS
});
```

### Direct RPC Access

For methods not wrapped in convenience functions:

```typescript
// Chat
await client.rpc.ChatSendMessage({
  sessionId: client.getSessionId()!,
  chatId: 'lobby',
  message: 'Hello world!'
});

// Create table
await client.rpc.RoomCreateTable({
  sessionId: client.getSessionId()!,
  roomId,
  matchOptions: {
    name: 'My Game',
    gameType: 'TwoPlayerDuel',
    deckType: 'Constructed',
    winsNeeded: 2,
    // ... other options
  }
});
```

## Type Safety Benefits

### Request/Response Typing

```typescript
// TypeScript knows the exact shape of requests and responses
const response: ConnectUserResponse = await client.connectUser('user', 'pass');

// Autocomplete and type checking
if (response.success) {
  const sessionId: string = response.sessionId; // Type: string
}
```

### Enum Support

```typescript
import { PlayerAction } from '$lib/generated/mage/v1/game';

// Type-safe enums
await client.rpc.SendPlayerAction({
  sessionId: client.getSessionId()!,
  gameId,
  action: PlayerAction.CONCEDE // Autocomplete available!
});
```

### Nested Types

```typescript
import type { GameView, CardView, PlayerView } from '$lib/generated/mage/v1/models';

const gameView: GameView = await client.getGameView(gameId, playerId);
const players: PlayerView[] = gameView.players || [];
const battlefield: CardView[] = gameView.battlefield || [];
```

## File Structure

```
mage-server-go/
├── api/proto/mage/v1/          # Source proto files (11 files)
│   ├── server.proto            # Main service definition
│   ├── auth.proto              # Auth messages
│   ├── room.proto              # Room messages
│   ├── table.proto             # Table messages
│   ├── game.proto              # Game messages
│   ├── tournament.proto        # Tournament messages
│   ├── draft.proto             # Draft messages
│   ├── chat.proto              # Chat messages
│   ├── admin.proto             # Admin messages
│   ├── models.proto            # Data models
│   └── websocket.proto         # WebSocket events
├── pkg/proto/mage/v1/          # Generated Go code (14 files)
│   ├── *.pb.go                 # Message structs
│   └── *_grpc.pb.go            # gRPC service stubs
└── scripts/
    └── generate_proto.sh       # Go code generation script

mage-client-web/
├── src/lib/generated/mage/v1/  # Generated TypeScript code (11 files)
│   ├── admin.ts
│   ├── auth.ts
│   ├── chat.ts
│   ├── draft.ts
│   ├── game.ts
│   ├── models.ts
│   ├── room.ts
│   ├── server.ts               # MageServerClientImpl
│   ├── table.ts
│   ├── tournament.ts
│   └── websocket.ts
├── src/lib/grpc/
│   └── client.ts               # MageClient wrapper
└── scripts/
    └── generate-proto.sh       # TypeScript code generation script
```

## Workflow

### After Updating Proto Files

When you modify any `.proto` file, regenerate both Go and TypeScript code:

```bash
# Go server
cd mage-server-go
make proto

# TypeScript client
cd mage-client-web
npm run proto:generate
```

### Development Cycle

1. **Define API**: Update `.proto` files in `mage-server-go/api/proto/mage/v1/`
2. **Generate Code**: Run `make proto` (Go) and `npm run proto:generate` (TypeScript)
3. **Implement**: Implement gRPC handlers in `mage-server-go/internal/server/grpc.go`
4. **Use Client**: Call methods via `MageClient` in web app

### Version Control

**Committed to git**:
- ✅ Source proto files (`api/proto/`)
- ✅ Generated Go code (`pkg/proto/`) - for convenience
- ✅ Generated TypeScript code (`src/lib/generated/`) - for convenience

While generated code is typically excluded, we include it for easier onboarding.

## Migration from Java Server

The protobuf definitions maintain compatibility with the Java XMage server API:

- **Message names** match Java class names (e.g., `ConnectUserRequest`)
- **Field names** match Java field names (camelCase)
- **RPC methods** match Java MageServer interface methods

This ensures existing Java clients can migrate gradually.

## Future Enhancements

### Planned Additions

- [ ] Streaming RPCs for real-time game state updates
- [ ] gRPC-Web proxy for direct browser access (currently using custom fetch wrapper)
- [ ] Compression for large messages (e.g., game state)
- [ ] Request/response interceptors for logging and metrics
- [ ] OpenAPI/REST gateway for HTTP clients

### Performance Optimizations

- [ ] Message batching for high-frequency updates
- [ ] Delta compression for game state changes
- [ ] Client-side caching with cache invalidation
- [ ] Connection pooling and reuse

## Troubleshooting

### "proto files not found" error

```bash
# Ensure proto files exist
ls mage-server-go/api/proto/mage/v1/

# Regenerate code
cd mage-server-go && make proto
cd mage-client-web && npm run proto:generate
```

### TypeScript import errors

```bash
# Ensure generation completed
cd mage-client-web
npm run proto:generate

# Check generated files
ls src/lib/generated/mage/v1/
```

### gRPC connection failed

Check that:
1. Server is running: `curl http://localhost:17171`
2. URL is correct in `.env`: `VITE_GRPC_SERVER_URL=http://localhost:17171`
3. CORS is configured if using different origin

### Type mismatches

Ensure both Go and TypeScript code are generated from the **same** proto files. The web client should NOT have separate proto files - it uses the server's proto files as the source of truth.

## References

- [Protocol Buffers Documentation](https://protobuf.dev/)
- [gRPC Documentation](https://grpc.io/)
- [ts-proto GitHub](https://github.com/stephenh/ts-proto)
- [gRPC-Web](https://github.com/grpc/grpc-web)

## Summary

✅ **11 proto files** defining complete API
✅ **70 RPC methods** covering all server operations
✅ **86 WebSocket events** for real-time updates
✅ **Go server code** generated and integrated
✅ **TypeScript client code** generated with type safety
✅ **Unified client wrapper** (`MageClient`) for ease of use
✅ **Single source of truth** - server proto files used for both Go and TypeScript
