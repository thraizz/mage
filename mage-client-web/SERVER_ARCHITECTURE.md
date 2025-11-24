# Server Architecture & Communication Flow

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                      MAGE WEB CLIENT                            │
│                     (SvelteKit + TypeScript)                    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │
        ┌─────────────────────┴─────────────────────┐
        │                                           │
        ▼                                           ▼
┌──────────────────┐                     ┌──────────────────┐
│   gRPC Client    │                     │ WebSocket Client │
│  (HTTP/JSON)     │                     │  (Binary/JSON)   │
└──────────────────┘                     └──────────────────┘
        │                                           │
        │ Port 17171                                │ TBD
        │ http://localhost:17171/mage.v1.*         │ ws://localhost:???
        │                                           │
        └─────────────────────┬─────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      MAGE GO SERVER                             │
│                    (gRPC + WebSocket)                           │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📡 Communication Layers

### **Layer 1: Transport**

```
┌────────────────────────────────────────────────────────┐
│                    CLIENT SIDE                         │
├────────────────────────────────────────────────────────┤
│                                                        │
│  ┌──────────────────┐        ┌──────────────────┐    │
│  │  fetch() API     │        │  WebSocket API   │    │
│  │  (Browser)       │        │  (Browser)       │    │
│  └──────────────────┘        └──────────────────┘    │
│         │                            │                │
│         │ HTTP POST                  │ WS Binary      │
│         │ JSON Body                  │ or JSON        │
│         │                            │                │
└─────────┼────────────────────────────┼────────────────┘
          │                            │
          │                            │
┌─────────┼────────────────────────────┼────────────────┐
│         │                            │                │
│         ▼                            ▼                │
│  ┌──────────────────┐        ┌──────────────────┐    │
│  │  gRPC Handler    │        │  WS Handler      │    │
│  │  (Go)            │        │  (Go)            │    │
│  └──────────────────┘        └──────────────────┘    │
│                                                        │
├────────────────────────────────────────────────────────┤
│                    SERVER SIDE                         │
└────────────────────────────────────────────────────────┘
```

### **Layer 2: API Clients**

```
┌────────────────────────────────────────────────────────┐
│                  API CLIENT LAYER                      │
├────────────────────────────────────────────────────────┤
│                                                        │
│  src/lib/grpc/client.ts                               │
│  ┌──────────────────────────────────────────────┐    │
│  │            MageClient Class                   │    │
│  │                                               │    │
│  │  • connectUser()                              │    │
│  │  • register()                                 │    │
│  │  • ping()                                     │    │
│  │  • getServerState()                           │    │
│  │  • getMainRoomId()                            │    │
│  │  • getAllTables()                             │    │
│  │  • getGameView()                              │    │
│  │  • sendChatMessage()                          │    │
│  │  • call<TReq, TRes>(method, request)          │    │
│  │                                               │    │
│  │  Session Management:                          │    │
│  │  • setSessionId(sessionId)                    │    │
│  │  • getSessionId()                             │    │
│  │  • clearSession()                             │    │
│  └──────────────────────────────────────────────┘    │
│                                                        │
│  src/lib/websocket.ts                                 │
│  ┌──────────────────────────────────────────────┐    │
│  │         MageWebSocket Class                   │    │
│  │                                               │    │
│  │  • connect()                                  │    │
│  │  • disconnect()                               │    │
│  │  • send(message)                              │    │
│  │  • on(type, handler)                          │    │
│  │  • onConnect(callback)                        │    │
│  │  • onDisconnect(callback)                     │    │
│  │  • isConnected()                              │    │
│  │                                               │    │
│  │  Auto-reconnect:                              │    │
│  │  • Max 5 attempts                             │    │
│  │  • Exponential backoff                        │    │
│  └──────────────────────────────────────────────┘    │
└────────────────────────────────────────────────────────┘
```

### **Layer 3: API Wrappers**

```
┌────────────────────────────────────────────────────────┐
│                API WRAPPER LAYER                       │
├────────────────────────────────────────────────────────┤
│                                                        │
│  src/lib/api/                                         │
│                                                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │
│  │  decks.ts   │  │  lobby.ts   │  │  chat.ts    │  │
│  │             │  │             │  │             │  │
│  │ ✅ Real     │  │ ⚠️  Mock    │  │ ⚠️  Mock    │  │
│  │ Server      │  │ Data        │  │ Data        │  │
│  │ Integration │  │             │  │             │  │
│  └─────────────┘  └─────────────┘  └─────────────┘  │
│                                                        │
│  ┌─────────────┐  ┌─────────────┐                    │
│  │  table.ts   │  │  match_     │                    │
│  │             │  │  history.ts │                    │
│  │ ⚠️  Mock    │  │             │                    │
│  │ Data        │  │ 🔜 TBD      │                    │
│  │             │  │             │                    │
│  └─────────────┘  └─────────────┘                    │
│                                                        │
└────────────────────────────────────────────────────────┘
```

---

## 🔄 Data Flow Examples

### **Example 1: Login Flow**

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   User      │     │   Login     │     │  MageClient │     │   Server    │
│             │     │   Page      │     │             │     │             │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                    │                    │
       │ Enter credentials │                    │                    │
       │──────────────────>│                    │                    │
       │                   │                    │                    │
       │                   │ connectUser()      │                    │
       │                   │───────────────────>│                    │
       │                   │                    │                    │
       │                   │                    │ POST /mage.v1.     │
       │                   │                    │ MageServer/        │
       │                   │                    │ ConnectUser        │
       │                   │                    │───────────────────>│
       │                   │                    │                    │
       │                   │                    │ ConnectUserResponse│
       │                   │                    │ { sessionId, ... } │
       │                   │                    │<───────────────────│
       │                   │                    │                    │
       │                   │ { success, token } │                    │
       │                   │<───────────────────│                    │
       │                   │                    │                    │
       │                   │ Store sessionId    │                    │
       │                   │ in MageClient      │                    │
       │                   │───────────────────>│                    │
       │                   │                    │                    │
       │                   │ auth.login(token)  │                    │
       │                   │─────────┐          │                    │
       │                   │         │          │                    │
       │                   │ Store token in     │                    │
       │                   │ localStorage       │                    │
       │                   │<────────┘          │                    │
       │                   │                    │                    │
       │ Navigate to lobby │                    │                    │
       │<──────────────────│                    │                    │
       │                   │                    │                    │
```

### **Example 2: Fetch Decks Flow**

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   User      │     │   Decks     │     │  decks.ts   │     │  MageClient │
│             │     │   Page      │     │  API        │     │             │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                    │                    │
       │ Open decks page   │                    │                    │
       │──────────────────>│                    │                    │
       │                   │                    │                    │
       │                   │ onMount()          │                    │
       │                   │────────┐           │                    │
       │                   │        │           │                    │
       │                   │ fetchUserDecks()   │                    │
       │                   │───────────────────>│                    │
       │                   │                    │                    │
       │                   │                    │ client.call(       │
       │                   │                    │  'DeckList',       │
       │                   │                    │  { sessionId }     │
       │                   │                    │ )                  │
       │                   │                    │───────────────────>│
       │                   │                    │                    │
       │                   │                    │                    │ ┌─────────┐
       │                   │                    │                    │ │ Server  │
       │                   │                    │                    │ └────┬────┘
       │                   │                    │                    │      │
       │                   │                    │ POST /mage.v1.MageServer/ │
       │                   │                    │      DeckList              │
       │                   │                    │─────────────────────────>│
       │                   │                    │                    │      │
       │                   │                    │ DeckListResponse   │      │
       │                   │                    │ { decks: [...] }   │      │
       │                   │                    │<─────────────────────────│
       │                   │                    │                    │
       │                   │ Convert to Deck[]  │                    │
       │                   │<───────────────────│                    │
       │                   │                    │                    │
       │ Display deck list │                    │                    │
       │<──────────────────│                    │                    │
       │                   │                    │                    │
```

### **Example 3: WebSocket Real-time Updates (Planned)**

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   User      │     │   Game      │     │  WebSocket  │     │   Server    │
│             │     │   Page      │     │  Client     │     │             │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                    │                    │
       │ Join game         │                    │                    │
       │──────────────────>│                    │                    │
       │                   │                    │                    │
       │                   │ ws.connect()       │                    │
       │                   │───────────────────>│                    │
       │                   │                    │                    │
       │                   │                    │ WS Handshake       │
       │                   │                    │───────────────────>│
       │                   │                    │                    │
       │                   │                    │ Connected          │
       │                   │                    │<───────────────────│
       │                   │                    │                    │
       │                   │ Connected!         │                    │
       │                   │<───────────────────│                    │
       │                   │                    │                    │
       │ Cast spell        │                    │                    │
       │──────────────────>│                    │                    │
       │                   │                    │                    │
       │                   │ ws.send({          │                    │
       │                   │   type: 'game_     │                    │
       │                   │   action',         │                    │
       │                   │   data: {...}      │                    │
       │                   │ })                 │                    │
       │                   │───────────────────>│                    │
       │                   │                    │                    │
       │                   │                    │ Send action        │
       │                   │                    │───────────────────>│
       │                   │                    │                    │
       │                   │                    │                    │ ┌──────┐
       │                   │                    │                    │ │Process│
       │                   │                    │                    │ │Action│
       │                   │                    │                    │ └───┬──┘
       │                   │                    │                    │     │
       │                   │                    │ GAME_UPDATE        │<────┘
       │                   │                    │ callback           │
       │                   │                    │<───────────────────│
       │                   │                    │                    │
       │                   │ Update game state  │                    │
       │                   │<───────────────────│                    │
       │                   │                    │                    │
       │ Render updated    │                    │                    │
       │ game state        │                    │                    │
       │<──────────────────│                    │                    │
       │                   │                    │                    │
```

---

## 🗂️ Proto Service Definitions

All services are defined in `src/lib/generated/mage/v1/server.ts`:

### **MageServer Service (Main Service)**

```typescript
export interface MageServer {
  // Authentication (4 methods)
  connectUser(request: ConnectUserRequest): Promise<ConnectUserResponse>;
  authRegister(request: AuthRegisterRequest): Promise<AuthRegisterResponse>;
  ping(request: PingRequest): Promise<PingResponse>;
  getServerState(request: GetServerStateRequest): Promise<GetServerStateResponse>;

  // Room Management (5 methods)
  serverGetMainRoomId(request: ServerGetMainRoomIdRequest): Promise<ServerGetMainRoomIdResponse>;
  roomGetAllTables(request: RoomGetAllTablesRequest): Promise<RoomGetAllTablesResponse>;
  roomJoinTable(request: RoomJoinTableRequest): Promise<RoomJoinTableResponse>;
  roomLeaveTable(request: RoomLeaveTableRequest): Promise<RoomLeaveTableResponse>;
  roomCreateTable(request: RoomCreateTableRequest): Promise<RoomCreateTableResponse>;

  // Table Management (10+ methods)
  tableJoinTournament(...): ...;
  tableWatchTable(...): ...;
  tableSetSeats(...): ...;
  tableStartMatch(...): ...;
  tableSubmitDeck(...): ...;
  // ... more table methods

  // Deck Management (4 methods)
  deckList(request: DeckListRequest): Promise<DeckListResponse>;
  deckGet(request: DeckGetRequest): Promise<DeckGetResponse>;
  deckSave(request: DeckSaveRequest): Promise<DeckSaveResponse>;
  deckDelete(request: DeckDeleteRequest): Promise<DeckDeleteResponse>;

  // Game Operations (20+ methods)
  gameGetView(request: GameGetViewRequest): Promise<GameGetViewResponse>;
  gameJoin(request: GameJoinRequest): Promise<GameJoinResponse>;
  gameWatchStart(...): ...;
  gameWatchStop(...): ...;
  matchStart(...): ...;
  matchQuit(...): ...;
  // ... more game methods

  // Chat Operations (5 methods)
  chatSendMessage(request: ChatSendMessageRequest): Promise<ChatSendMessageResponse>;
  chatJoin(request: ChatJoinRequest): Promise<ChatJoinResponse>;
  chatLeave(request: ChatLeaveRequest): Promise<ChatLeaveResponse>;
  chatFindByRoom(...): ...;
  chatFindByTable(...): ...;

  // Draft Operations (5 methods)
  draftJoin(...): ...;
  draftQuit(...): ...;
  sendDraftCardPick(...): ...;
  sendDraftCardMark(...): ...;
  // ... more draft methods

  // Tournament Operations (5 methods)
  tournamentJoin(...): ...;
  tournamentQuit(...): ...;
  tournamentGetView(...): ...;
  // ... more tournament methods

  // Admin Operations (10+ methods)
  adminGetUsers(...): ...;
  adminLockUser(...): ...;
  adminMuteUser(...): ...;
  adminDisconnectUser(...): ...;
  // ... more admin methods

  // ... Total: 70+ RPC methods
}
```

---

## 🔐 Security & Session Management

### **Session Flow**

```
┌──────────────────────────────────────────────────────────┐
│                  SESSION MANAGEMENT                      │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  1. User logs in                                         │
│     ↓                                                    │
│  2. Server returns sessionId + JWT token                 │
│     ↓                                                    │
│  3. Client stores:                                       │
│     • sessionId → MageClient instance (memory)           │
│     • JWT token → localStorage (persistent)              │
│     ↓                                                    │
│  4. All subsequent RPC calls include:                    │
│     • sessionId in request body                          │
│     • JWT in Authorization header (optional)             │
│     ↓                                                    │
│  5. Server validates session                             │
│     ↓                                                    │
│  6. If session expired:                                  │
│     • Server returns auth error                          │
│     • Client auto-logs out                               │
│     • Client redirects to login                          │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

### **Request Headers**

```typescript
// gRPC call in client.ts
const response = await fetch(`${serverUrl}/mage.v1.MageServer/${method}`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    // Optional: JWT authentication
    // 'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify(request) // includes sessionId
});
```

---

## 📊 State Management

### **Client-side State**

```
┌───────────────────────────────────────────────────────────┐
│                  SVELTE STORES                            │
├───────────────────────────────────────────────────────────┤
│                                                           │
│  src/lib/stores/auth.ts                                   │
│  ┌─────────────────────────────────────────────────┐     │
│  │  AuthState {                                     │     │
│  │    isAuthenticated: boolean                      │     │
│  │    token: string | null                          │     │
│  │    user: { id, username, email } | null          │     │
│  │  }                                               │     │
│  │                                                   │     │
│  │  Methods:                                        │     │
│  │  • login(token, user)                            │     │
│  │  • logout()                                      │     │
│  │  • loadAuthFromStorage()                         │     │
│  │  • checkTokenValidity()                          │     │
│  └─────────────────────────────────────────────────┘     │
│                                                           │
│  src/lib/stores/connection.ts                             │
│  ┌─────────────────────────────────────────────────┐     │
│  │  ConnectionState {                               │     │
│  │    status: 'connected' | 'connecting' |          │     │
│  │            'disconnected' | 'reconnecting'       │     │
│  │    lastConnected: number | null                  │     │
│  │    lastDisconnected: number | null               │     │
│  │    reconnectAttempt: number                      │     │
│  │    error: string | null                          │     │
│  │    latency: number | null                        │     │
│  │  }                                               │     │
│  │                                                   │     │
│  │  Methods:                                        │     │
│  │  • connect()                                     │     │
│  │  • disconnect()                                  │     │
│  │  • reconnect()                                   │     │
│  │  • pong() (health check response)                │     │
│  └─────────────────────────────────────────────────┘     │
│                                                           │
│  src/lib/stores/game.ts (planned)                         │
│  ┌─────────────────────────────────────────────────┐     │
│  │  GameState {                                     │     │
│  │    gameId: string | null                         │     │
│  │    gameView: GameView | null                     │     │
│  │    hand: Card[]                                  │     │
│  │    battlefield: Card[]                           │     │
│  │    stack: Card[]                                 │     │
│  │    currentPhase: string                          │     │
│  │    priorityPlayerId: string | null               │     │
│  │  }                                               │     │
│  └─────────────────────────────────────────────────┘     │
│                                                           │
└───────────────────────────────────────────────────────────┘
```

---

## 🧪 Testing & Development

### **Mock vs Real Implementations**

| Component | Current | Target | How to Switch |
|-----------|---------|--------|---------------|
| **Login** | Mock | Real | Replace `simulateLogin()` with `client.connectUser()` |
| **Register** | Mock | Real | Replace `simulateRegister()` with `client.register()` |
| **Fetch Tables** | Mock | Real | Replace `MOCK_TABLES` with `client.getAllTables()` |
| **Create Table** | Mock | Real | Call `client.call('RoomCreateTable', request)` |
| **Join Table** | Mock | Real | Call `client.call('RoomJoinTable', request)` |
| **Chat Messages** | Mock | Real | Use WebSocket callbacks |
| **Fetch Decks** | ✅ Real | ✅ Real | Already using `client.call('DeckList')` |
| **Save Deck** | ✅ Real | ✅ Real | Already using `client.call('DeckSave')` |
| **Delete Deck** | ✅ Real | ✅ Real | Already using `client.call('DeckDelete')` |

### **Environment Configuration**

```bash
# Development
VITE_GRPC_SERVER_URL=http://localhost:17171

# Staging
VITE_GRPC_SERVER_URL=https://staging-api.mage.example.com

# Production
VITE_GRPC_SERVER_URL=https://api.mage.example.com
```

---

## 🚦 Error Handling

### **Error Flow**

```
┌─────────────────────────────────────────────────────────────┐
│                    ERROR HANDLING                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  RPC Call Fails                                             │
│    ↓                                                        │
│  toGrpcError(error)   (utils/grpc-errors.ts)               │
│    ↓                                                        │
│  ┌─────────────────────────────────────┐                   │
│  │ Is Auth Error?                      │                   │
│  │ (UNAUTHENTICATED, PERMISSION_DENIED)│                   │
│  └─────────────┬───────────────────────┘                   │
│                │                                            │
│         YES ───┴──> auth.logout()                           │
│                     toast.error('Session expired')          │
│                     Redirect to /login                      │
│                                                             │
│                │                                            │
│         NO  ───┴──> Is Retryable?                           │
│                     (UNAVAILABLE, DEADLINE_EXCEEDED,        │
│                      RESOURCE_EXHAUSTED)                    │
│                     │                                       │
│                     └─> YES → Retry with backoff            │
│                         NO  → Show error to user            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### **gRPC Status Codes**

```typescript
export enum GrpcStatusCode {
  OK = 0,
  CANCELLED = 1,
  UNKNOWN = 2,
  INVALID_ARGUMENT = 3,
  DEADLINE_EXCEEDED = 4,
  NOT_FOUND = 5,
  ALREADY_EXISTS = 6,
  PERMISSION_DENIED = 7,
  RESOURCE_EXHAUSTED = 8,
  FAILED_PRECONDITION = 9,
  ABORTED = 10,
  OUT_OF_RANGE = 11,
  UNIMPLEMENTED = 12,
  INTERNAL = 13,
  UNAVAILABLE = 14,
  DATA_LOSS = 15,
  UNAUTHENTICATED = 16,
}
```

---

## 📈 Performance Considerations

### **Connection Health**

```
┌────────────────────────────────────────────────────────┐
│           CONNECTION HEALTH MONITORING                 │
├────────────────────────────────────────────────────────┤
│                                                        │
│  Every 30 seconds:                                     │
│    ↓                                                   │
│  Send Ping (client.ping())                             │
│    ↓                                                   │
│  Wait for Pong (max 5 seconds)                         │
│    ↓                                                   │
│  ┌────────────────────────┐                           │
│  │ Pong received?         │                           │
│  └────────┬───────────────┘                           │
│           │                                            │
│    YES ───┴──> Update latency metric                   │
│               connection.latency = responseTime        │
│                                                        │
│           │                                            │
│    NO  ───┴──> Connection lost!                        │
│               connection.simulateError(...)            │
│               Trigger reconnection logic               │
│                                                        │
└────────────────────────────────────────────────────────┘
```

### **Reconnection Strategy**

```
Attempt 1: Wait 1 second   (1000ms)
Attempt 2: Wait 2 seconds  (2000ms)
Attempt 3: Wait 4 seconds  (4000ms)
Attempt 4: Wait 8 seconds  (8000ms)
Attempt 5: Wait 16 seconds (16000ms)
...
Max delay: 30 seconds      (30000ms)
Max attempts: 10

If all attempts fail:
  → Show "Unable to connect" error
  → User can manually retry
```

---

## 🎯 Summary

The mage-client-web uses a **two-channel architecture**:

1. **gRPC over HTTP** - Request/response RPCs for actions and queries
2. **WebSocket** - Real-time streaming updates for game state, chat, etc.

The infrastructure is **fully in place**, with:
- ✅ Type-safe proto definitions generated
- ✅ gRPC client implemented and working
- ✅ WebSocket client implemented (needs connection to server)
- ✅ Error handling, retry logic, and health checks
- ✅ Session management and authentication flows
- ✅ Deck management fully integrated with server
- ⚠️ Auth, lobby, and chat using mock data (easy to replace)

**Next steps:** Replace mock implementations with real server calls by swapping function implementations in `src/lib/api/*.ts` files.

