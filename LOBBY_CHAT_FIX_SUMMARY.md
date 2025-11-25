# Lobby Chat Fix Summary

## Problems Identified

### 1. LobbyChat Not Subscribing to WebSocket Messages
**Root Cause**: The `LobbyChat.svelte` component was only calling `fetchLobbyMessages()` on mount, which returns an empty array. It never subscribed to WebSocket events to receive real-time messages.

**Impact**: Messages sent to lobby chat were successfully transmitted to the server and broadcasted via WebSocket, but the LobbyChat component wasn't listening, so no messages appeared in the UI.

### 2. WebSocket Duplicate Connection Race Condition
**Root Cause**: When the lobby page loaded, it called `connectLobbyUpdates()` which creates a WebSocket connection. Due to a race condition, if `connect()` was called again before the first connection finished opening (state = CONNECTING), a second WebSocket would be created.

**Impact**:
- Two WebSocket connections with the same session ID
- Both reading from the same `CallbackChan` on the server
- Only one connection receives each message (channel reads are non-broadcast)
- The other connection times out and closes, causing "broken pipe" errors in server logs
- Users might not receive all messages

## Fixes Applied

### Fix 1: LobbyChat WebSocket Subscription

**File**: `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/LobbyChat.svelte`

**Changes**:
1. Added WebSocket imports:
   - `websocketStore`
   - `CallbackMethod`
   - `ChatMessageData`
   - `getMageClient`

2. Added state management:
   ```typescript
   let chatId = $state<string | null>(null);
   let connected = $state(false);
   let unsubscribeWebSocket: (() => void) | null = null;
   ```

3. Replaced `loadMessages()` with `initializeChat()`:
   ```typescript
   async function initializeChat(): Promise<void> {
       // Get main room ID
       const roomResponse = await client.getMainRoomId();

       // Find chat ID for the lobby
       const chatResponse = await client.call('ChatFindByRoom', {
           sessionId: await client.ensureSessionId(),
           roomId: roomResponse.roomId
       });

       chatId = chatResponse.chatId;

       // Join the chat room
       await joinChat(chatId);

       // Subscribe to WebSocket chat messages
       unsubscribeWebSocket = websocketStore.on(
           CallbackMethod.CHATMESSAGE,
           handleChatMessage
       );

       connected = true;
   }
   ```

4. Added `handleChatMessage()` to process incoming WebSocket messages:
   ```typescript
   function handleChatMessage(data: unknown): void {
       const messageData = data as ChatMessageData;

       // Only process messages for this chat
       if (messageData.chatId !== chatId) return;

       const message: ChatMessage = {
           id: `msg-${Date.now()}-${Math.random()}`,
           type: protoMessage.userName.toLowerCase() === 'system' ? 'system' : 'user',
           username: protoMessage.userName || 'Unknown',
           content: protoMessage.message,
           timestamp: protoMessage.time?.getTime() || Date.now()
       };

       messages.push(message);
       setTimeout(() => scrollToBottom(), 50);
   }
   ```

5. Updated `handleSendMessage()` to not push messages locally:
   ```typescript
   // Before: messages.push(message);
   // After: Messages come back via WebSocket
   ```

6. Added cleanup in `onDestroy()`:
   ```typescript
   onDestroy(() => {
       if (unsubscribeWebSocket) {
           unsubscribeWebSocket();
       }
   });
   ```

7. Added connection status indicator in header:
   - Shows "Live" when connected
   - Shows "Connecting..." when loading
   - Shows "Offline" on error

### Fix 2: WebSocket Duplicate Connection Prevention

**File**: `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/websocket.ts`

**Changes**:
Enhanced the `connect()` function to handle the CONNECTING state:

```typescript
function connect(newSessionId: string): Promise<void> {
    return new Promise((resolve, reject) => {
        // Check if already connected or connecting
        if (ws) {
            if (ws.readyState === WebSocket.OPEN) {
                resolve();
                return;
            } else if (ws.readyState === WebSocket.CONNECTING) {
                // Already connecting, wait for it to finish
                const checkConnection = setInterval(() => {
                    if (ws && ws.readyState === WebSocket.OPEN) {
                        clearInterval(checkConnection);
                        resolve();
                    } else if (!ws || ws.readyState === WebSocket.CLOSED) {
                        clearInterval(checkConnection);
                        reject(new Error('WebSocket connection failed'));
                    }
                }, 100);
                return;
            }
        }

        // Continue with connection creation...
    });
}
```

**How it works**:
1. If WebSocket is OPEN → return immediately (existing behavior)
2. If WebSocket is CONNECTING → **wait for it** instead of creating a new one
3. Poll every 100ms until state becomes OPEN or CLOSED
4. Only create a new WebSocket if none exists or previous one is CLOSED

## Message Flow (After Fixes)

### Sending a Message
1. User types message in LobbyChat and clicks send
2. `sendLobbyMessage()` calls gRPC `ChatSendMessage`
3. Server's `chat.SendMessage()` stores message and calls `broadcastMessage()`
4. Server iterates through all users in the chat room
5. For each user's session, server pushes message to `CallbackChan`
6. WebSocket goroutine reads from `CallbackChan` and sends JSON to client

### Receiving a Message
1. Client's WebSocket receives JSON message
2. `websocketStore` parses it as `ServerEvent`
3. Finds all registered handlers for `CallbackMethod.CHATMESSAGE`
4. Calls `handleChatMessage()` in LobbyChat component
5. Component filters by `chatId` (only process lobby messages)
6. Converts proto message to client `ChatMessage` type
7. Pushes to `messages` array
8. UI reactively updates and auto-scrolls

### Key Points
- **Single WebSocket per session** - no more duplicate connections
- **Multiple subscribers per event type** - both lobby-updates and LobbyChat can subscribe to CHATMESSAGE
- **Real-time bidirectional** - sender and all recipients see messages instantly
- **No local message injection** - all messages come through WebSocket for consistency

## Testing

To verify the fixes:

1. **Start the servers**:
   ```bash
   # Terminal 1: Start Go server
   cd /Users/aron/dev/opensource/mage/mage-server-go
   ./server

   # Terminal 2: Start web client
   cd /Users/aron/dev/opensource/mage/mage-client-web
   npm run dev
   ```

2. **Test with two browser windows**:
   - Open `http://localhost:5173/lobby` in two different browser windows/tabs
   - Login in both (or use guest login)
   - Send a message from one window
   - **Expected**: Message appears instantly in both windows
   - **Expected**: Server logs show only ONE WebSocket connection per session
   - **Expected**: No "broken pipe" errors in server logs

3. **Check server logs**:
   ```
   ✅ Good: {"level":"info","msg":"WebSocket connected","session":"xxx","user":"username"}
   ✅ Good: {"level":"debug","msg":"broadcasted chat message","users":2,"sessions_sent":2}
   ❌ Bad:  {"level":"error","msg":"failed to send event","error":"write tcp ... broken pipe"}
   ```

## Related Files

### Client-Side
- `/mage-client-web/src/lib/components/LobbyChat.svelte` - Main fix
- `/mage-client-web/src/lib/stores/websocket.ts` - Race condition fix
- `/mage-client-web/src/lib/api/chat.ts` - Chat API functions (unchanged)
- `/mage-client-web/src/lib/services/lobby-updates.ts` - Lobby WebSocket setup (unchanged)
- `/mage-client-web/src/routes/(protected)/lobby/+page.svelte` - Lobby page (unchanged)

### Server-Side
- `/mage-server-go/internal/chat/manager.go` - Broadcasting logic (already working, no changes)
- `/mage-server-go/internal/server/websocket.go` - WebSocket handler (unchanged)
- `/mage-server-go/internal/session/session.go` - Session with CallbackChan (unchanged)

## Future Improvements

1. **Server-Side Session WebSocket Tracking**:
   - Track individual WebSocket connections per session
   - Allow multiple connections with separate callback channels
   - Properly cleanup stale connections

2. **Message History**:
   - Implement `ChatGetMessages` RPC to return actual message history
   - Load recent messages on component mount
   - Currently returns empty array

3. **Typing Indicators**:
   - Add "User is typing..." WebSocket events
   - Show in chat UI

4. **Read Receipts**:
   - Track which messages each user has seen
   - Show read status indicators

5. **Message Persistence**:
   - Store messages in database
   - Survive server restarts
   - Support message editing/deletion

6. **Rate Limiting Server-Side**:
   - Currently only client-side rate limiting
   - Add server-side protection against spam

## Notes

- The table chat was already working correctly (as documented in `CHAT_FIX_SUMMARY.md`)
- The server broadcast logic was already implemented and working
- The issue was entirely on the client side - LobbyChat wasn't listening
- The WebSocket race condition was a subtle bug that could affect any page using WebSocket
