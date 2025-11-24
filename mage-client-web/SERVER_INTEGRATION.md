# Server Integration Analysis

## Overview

The mage-client-web communicates with the Go server using **gRPC/Protobuf over HTTP** (for RPC calls) and **WebSocket** (for real-time streaming updates).

---

## 🔌 Connection Methods

### 1. **gRPC Client** (`src/lib/grpc/client.ts`)

The main interface for communicating with the Go server.

**Server URL Configuration:**

```typescript
serverUrl: import.meta.env.VITE_GRPC_SERVER_URL || 'http://localhost:17171';
```

**Key Features:**

- Type-safe RPC method calls
- Session ID management (stored after login)
- JWT token authentication
- Timeout and retry logic
- Error handling with automatic logout on auth errors

### 2. **WebSocket Client** (`src/lib/websocket.ts`)

Generic WebSocket client for real-time bidirectional communication.

**Features:**

- Automatic reconnection (up to 5 attempts)
- Message type routing
- Event callbacks (onConnect, onDisconnect)
- Connection state management

**Note:** Currently has a basic implementation but is designed to handle real-time game state updates from the server.

---

## 📡 Server Communication Points

### **Authentication & Connection**

#### 1. Login (`src/routes/login/+page.svelte`)

**Status:** ⚠️ Mock Implementation

```typescript
// Currently simulated - will be replaced with:
const client = getMageClient();
const response = await client.connectUser(username, password);
```

**Server Method:** `ConnectUser`

- **Request:** `ConnectUserRequest` (userName, password, clientVersion)
- **Response:** `ConnectUserResponse` (success, sessionId, userId, userName)

#### 2. Registration (`src/routes/register/+page.svelte`)

**Status:** ⚠️ Mock Implementation

```typescript
// Will be replaced with:
const client = getMageClient();
const response = await client.register(userName, password, email);
```

**Server Method:** `AuthRegister`

- **Request:** `AuthRegisterRequest` (userName, password, email)
- **Response:** `AuthRegisterResponse` (success, error, userId)

#### 3. Ping/Keepalive (`src/lib/grpc/client.ts`)

**Status:** ✅ Implemented

```typescript
await client.ping();
```

**Server Method:** `Ping`

- Keeps session alive
- Checks connection health
- Returns server timestamp and latency

---

### **Lobby & Room Operations**

#### 1. Get Main Room ID (`src/lib/grpc/client.ts`)

**Status:** ✅ Implemented

```typescript
const lobby = await client.getMainRoomId();
```

**Server Method:** `ServerGetMainRoomId`

- **Request:** `ServerGetMainRoomIdRequest` (sessionId)
- **Response:** `ServerGetMainRoomIdResponse` (roomId)

#### 2. Get All Tables (`src/lib/grpc/client.ts`)

**Status:** ✅ Implemented

```typescript
const tables = await client.getAllTables(roomId);
```

**Server Method:** `RoomGetAllTables`

- **Request:** `RoomGetAllTablesRequest` (sessionId, roomId)
- **Response:** `RoomGetAllTablesResponse` (tables[])

#### 3. Fetch Tables (`src/lib/api/lobby.ts`)

**Status:** ⚠️ Mock Implementation (uses hard-coded data)

```typescript
// Currently returns MOCK_TABLES
// Will be replaced with actual gRPC call
```

**Usage:** `src/routes/(protected)/lobby/+page.svelte`

#### 4. Create Table (`src/lib/api/lobby.ts`)

**Status:** ⚠️ Mock Implementation

**Server Method (when implemented):** `RoomCreateTable`

- **Request:** `RoomCreateTableRequest` (sessionId, roomId, matchOptions)
- **Response:** `RoomCreateTableResponse` (tableId)

#### 5. Join Table (`src/lib/api/lobby.ts`)

**Status:** ⚠️ Mock Implementation

**Server Method (when implemented):** `RoomJoinTable`

#### 6. Leave Table (`src/lib/api/lobby.ts`)

**Status:** ⚠️ Mock Implementation

**Server Method (when implemented):** `RoomLeaveTable`

---

### **Deck Management**

#### 1. Fetch User Decks (`src/lib/api/decks.ts`)

**Status:** ✅ Implemented

```typescript
const client = getMageClient();
const response = await client.call('DeckList', request);
```

**Server Method:** `DeckList`

- **Request:** `DeckListRequest` (sessionId, format?)
- **Response:** `DeckListResponse` (success, decks[], error?)

#### 2. Get Deck Details (`src/lib/api/decks.ts`)

**Status:** ✅ Implemented

```typescript
await client.call('DeckGet', { sessionId, deckId });
```

**Server Method:** `DeckGet`

- **Request:** `DeckGetRequest` (sessionId, deckId)
- **Response:** `DeckGetResponse` (success, info, deck, error?)

#### 3. Upload/Save Deck (`src/lib/api/decks.ts`)

**Status:** ✅ Implemented

```typescript
await client.call('DeckSave', saveRequest);
```

**Server Method:** `DeckSave`

- **Request:** `DeckSaveRequest` (sessionId, deckName, deck, format, description)
- **Response:** `DeckSaveResponse` (success, deckId, error?)

#### 4. Delete Deck (`src/lib/api/decks.ts`)

**Status:** ✅ Implemented

```typescript
await client.call('DeckDelete', request);
```

**Server Method:** `DeckDelete`

- **Request:** `DeckDeleteRequest` (sessionId, deckId)
- **Response:** `DeckDeleteResponse` (success, error?)

---

### **Chat System**

#### 1. Send Chat Message (`src/lib/grpc/client.ts`)

**Status:** ✅ Implemented

```typescript
await client.sendChatMessage(chatId, message);
```

**Server Method:** `ChatSendMessage`

- **Request:** `ChatSendMessageRequest` (sessionId, chatId, message)
- **Response:** `ChatSendMessageResponse` (success, messageId)

#### 2. Fetch Lobby Messages (`src/lib/api/chat.ts`)

**Status:** ⚠️ Mock Implementation (uses MOCK_MESSAGES)

**Server Method (when implemented):** `ChatGetMessages` or streaming via WebSocket

#### 3. Send Whisper (`src/lib/api/chat.ts`)

**Status:** ⚠️ Mock Implementation

**Server Method (when implemented):** `ChatSendWhisper`

---

### **Game Operations**

#### 1. Get Game View (`src/lib/grpc/client.ts`)

**Status:** ✅ Implemented

```typescript
await client.getGameView(gameId, playerId);
```

**Server Method:** `GameGetView`

- **Request:** `GameGetViewRequest` (sessionId, gameId, playerId)
- **Response:** `GameGetViewResponse` (gameView)

#### 2. Start Game (`src/lib/api/table.ts`)

**Status:** ⚠️ Mock Implementation

**Server Method (when implemented):** `MatchStart`

#### 3. Game Actions

**Status:** Planned for WebSocket implementation

**Server Methods (when implemented):**

- `GamePassPriority`
- `GameConcede`
- `GamePlayCard`
- `GameActivateAbility`
- `GameDeclareAttackers`
- `GameDeclareBlockers`
- etc.

---

### **Real-time Updates (WebSocket)**

The client is designed to receive real-time updates via WebSocket using callback methods defined in `src/lib/generated/mage/v1/websocket.ts`.

#### **Callback Methods** (Server → Client)

##### Chat & Messages

- `CHATMESSAGE` - New chat message
- `SHOW_USERMESSAGE` - User-specific message
- `SERVER_MESSAGE` - Server announcement

##### Table Events

- `JOINED_TABLE` - Player joined table
- `TABLE_WAITING` - Table status update

##### Tournament Events

- `START_TOURNAMENT` - Tournament started
- `TOURNAMENT_INIT` - Tournament initialization
- `TOURNAMENT_UPDATE` - Tournament state update
- `TOURNAMENT_OVER` - Tournament ended

##### Draft Events

- `START_DRAFT` - Draft started
- `DRAFT_INIT` - Draft initialization
- `DRAFT_PICK` - Card picked in draft
- `DRAFT_UPDATE` - Draft state update
- `DRAFT_OVER` - Draft completed

##### Game Events

- `START_GAME` - Game started
- `GAME_INIT` - Game initialization
- `GAME_UPDATE` - Game state update
- `GAME_UPDATE_AND_INFORM` - Game state + info message
- `GAME_INFORM_PERSONAL` - Personal game info
- `GAME_ERROR` - Game error
- `GAME_TARGET` - Target selection required
- `GAME_CHOOSE_ABILITY` - Ability choice required
- `GAME_CHOOSE_PILE` - Pile choice required
- `GAME_CHOOSE_CHOICE` - Generic choice required
- `GAME_ASK` - Yes/No question
- `GAME_SELECT` - Select cards/permanents
- `GAME_PLAY_MANA` - Mana payment required
- `GAME_PLAY_XMANA` - X mana payment required
- `GAME_GET_AMOUNT` - Amount input required
- `GAME_GET_MULTI_AMOUNT` - Multiple amount inputs required
- `GAME_OVER` - Game ended
- `END_GAME_INFO` - End game statistics

##### Watch Events

- `SHOW_TOURNAMENT` - Spectate tournament
- `WATCHGAME` - Spectate game

##### Replay Events

- `REPLAY_GAME` - Replay started
- `REPLAY_INIT` - Replay initialization
- `REPLAY_UPDATE` - Replay state update
- `REPLAY_DONE` - Replay completed

##### Deck View Events

- `VIEW_LIMITED_DECK` - View limited deck
- `VIEW_SIDEBOARD` - View sideboard

##### User Interaction

- `USER_REQUEST_DIALOG` - User dialog required
- `GAME_REDRAW_GUI` - GUI refresh needed

---

## 🗂️ Generated Proto Types

All protocol buffer definitions are generated into `src/lib/generated/`:

### **Available Services:**

1. **`admin.ts`** - Admin operations (disconnect user, lock user, mute user, etc.)
2. **`auth.ts`** - Authentication & authorization (login, register, password reset)
3. **`chat.ts`** - Chat operations (join, leave, send message)
4. **`draft.ts`** - Draft operations (join, pick card, mark card)
5. **`game.ts`** - Game operations (join, watch, replay, actions)
6. **`models.ts`** - Shared data models (Card, Player, Game, etc.)
7. **`room.ts`** - Room/lobby operations (get tables, create room)
8. **`server.ts`** - Main MageServer service definition (70+ RPC methods)
9. **`table.ts`** - Table operations (create, join, start match, deck management)
10. **`tournament.ts`** - Tournament operations (join, leave, get info)
11. **`websocket.ts`** - WebSocket callback definitions

### **Total RPC Methods:** 70+

---

## 🔄 Connection Flow

### **1. Initial Connection**

```
User opens app
  ↓
Load auth from localStorage
  ↓
If token valid → Auto-connect
  ↓
Connect to gRPC server (http://localhost:17171)
  ↓
Ping server to verify connection
```

### **2. Login Flow**

```
User enters credentials
  ↓
client.connectUser(username, password)
  ↓
Server validates credentials
  ↓
Returns sessionId + JWT token
  ↓
Store sessionId in MageClient
  ↓
Store JWT in localStorage
  ↓
Redirect to /lobby
```

### **3. Lobby Flow**

```
User enters lobby
  ↓
client.getMainRoomId()
  ↓
client.getAllTables(roomId)
  ↓
Display tables
  ↓
User creates/joins table
  ↓
Establish WebSocket connection for real-time updates
```

### **4. Game Flow**

```
Table ready to start
  ↓
Host starts match
  ↓
Server creates game instance
  ↓
Sends START_GAME callback via WebSocket
  ↓
Client navigates to /game/[id]
  ↓
client.getGameView(gameId, playerId)
  ↓
Render game state
  ↓
Listen for GAME_UPDATE callbacks
  ↓
User makes actions → Send RPC calls
  ↓
Receive updates via WebSocket
```

---

## 🔐 Authentication

### **Session Management**

1. **Session ID** - Stored in `MageClient` instance, required for most RPC calls
2. **JWT Token** - Stored in `localStorage` under key `mage_auth_token`
3. **Auto-restore** - On app load, JWT is validated and session restored if valid

### **Authentication Flow**

```typescript
// Metadata injection in service-factory.ts
const metadata = createGrpcMetadata({
	authorization: `Bearer ${token}`
});

// Automatic logout on auth errors
if (isAuthError(grpcError)) {
	auth.logout();
	toast.error('Session expired. Please log in again.');
}
```

---

## 📊 Current Implementation Status

| Feature               | Status         | Notes                                                 |
| --------------------- | -------------- | ----------------------------------------------------- |
| **Authentication**    | ⚠️ Mock        | Uses simulated login, needs real gRPC integration     |
| **Deck Management**   | ✅ Implemented | Fully integrated with server                          |
| **Lobby Tables**      | ⚠️ Mock        | Uses hard-coded data, needs gRPC integration          |
| **Chat**              | ⚠️ Mock        | Uses hard-coded messages, needs WebSocket integration |
| **Game View**         | ✅ Implemented | Can fetch game state via gRPC                         |
| **Game Actions**      | 🔜 Planned     | Needs WebSocket bidirectional communication           |
| **Real-time Updates** | 🔜 Planned     | WebSocket callbacks defined, not yet connected        |
| **Connection Health** | ✅ Implemented | Ping, reconnection, health checks working             |
| **Error Handling**    | ✅ Implemented | Comprehensive gRPC error handling                     |

**Legend:**

- ✅ Implemented - Fully working with real server integration
- ⚠️ Mock - Interface exists but uses simulated data
- 🔜 Planned - Types/interfaces defined, implementation pending

---

## 🛠️ Development Setup

### **Environment Variables**

Create `.env` file (optional):

```bash
VITE_GRPC_SERVER_URL=http://localhost:17171
```

If not set, defaults to `http://localhost:17171`.

### **Start Server**

```bash
cd ../mage-server-go
go run cmd/server/main.go
```

Server should start on port `17171`.

### **Start Client**

```bash
cd mage-client-web
bun install
bun run dev
```

Client will be available at `http://localhost:5173`.

---

## 🔍 Files That Communicate with Server

### **Direct Server Communication:**

1. **`src/lib/grpc/client.ts`** - Main gRPC client, all RPC methods
2. **`src/lib/websocket.ts`** - WebSocket client wrapper
3. **`src/lib/api/decks.ts`** - Deck management API (✅ real server calls)
4. **`src/lib/grpc/service-factory.ts`** - gRPC call wrapper with auth

### **Planned Server Communication:**

5. **`src/lib/api/lobby.ts`** - Lobby operations (currently mocked)
6. **`src/lib/api/chat.ts`** - Chat operations (currently mocked)
7. **`src/lib/api/table.ts`** - Table operations (currently mocked)
8. **`src/lib/api/match_history.ts`** - Match history (to be implemented)

### **Server Connection Management:**

9. **`src/lib/stores/auth.ts`** - Authentication state
10. **`src/lib/stores/connection.ts`** - Connection state & health checks

### **Pages Using Server:**

11. **`src/routes/login/+page.svelte`** - Login (needs connectUser)
12. **`src/routes/register/+page.svelte`** - Registration (needs register)
13. **`src/routes/(protected)/lobby/+page.svelte`** - Lobby (needs getAllTables)
14. **`src/routes/(protected)/decks/+page.svelte`** - Deck management (✅ uses real server)
15. **`src/routes/(protected)/table/[id]/+page.svelte`** - Table lobby
16. **`src/routes/(protected)/game/[id]/+page.svelte`** - Game view

---

## 🚀 Next Steps for Full Integration

### **Priority 1: Authentication**

- Replace mock login with `client.connectUser()`
- Replace mock register with `client.register()`
- Test session persistence and token refresh

### **Priority 2: Lobby Integration**

- Replace `fetchTables()` with `client.getAllTables()`
- Implement `createTable()` with `RoomCreateTable` RPC
- Implement `joinTable()` with `RoomJoinTable` RPC
- Implement `leaveTable()` with `RoomLeaveTable` RPC

### **Priority 3: WebSocket Integration**

- Connect WebSocket client to server
- Implement callback message routing
- Handle real-time game state updates
- Handle real-time chat messages
- Handle table state updates

### **Priority 4: Game Actions**

- Implement game action RPCs
- Handle game state updates via WebSocket
- Implement user interaction prompts (target selection, choices, etc.)

### **Priority 5: Chat System**

- Replace mock chat with WebSocket-based chat
- Implement whisper messages
- Handle server announcements

---

## 📝 Summary

The mage-client-web is designed to communicate with the Go server using:

1. **gRPC/Protobuf over HTTP** for RPC calls (request/response)
2. **WebSocket** for real-time bidirectional communication (server push)

**Current Status:**

- ✅ Infrastructure is fully set up and working
- ✅ Deck management is fully integrated
- ⚠️ Auth, lobby, and chat are mocked but have interfaces ready
- 🔜 Game actions and real-time updates need WebSocket integration

All proto definitions are generated and type-safe. The client is ready to replace mock implementations with real server calls by simply swapping out the mock functions with the already-implemented gRPC client methods.
