# WebSocket Chat Fix Summary

## Problems Identified

Based on browser console logs and server logs, there were **3 critical issues** preventing chat from working:

### 1. WebSocket Connection Not Persisting
**Symptom**: Server logs showed `WebSocket closed` immediately after `WebSocket connected`

**Root Cause**: The WebSocket server wasn't reading from the client connection. In Go's gorilla/websocket, you MUST actively read from the connection to process control frames (like PONG responses to PING messages). Without a read loop, the connection would close when the read deadline expired.

**Fix** (`mage-server-go/internal/server/websocket.go`):
- Added `readHandler()` goroutine that:
  - Sets up a PONG handler to reset read deadline
  - Continuously reads messages (even though client doesn't send any)
  - Properly processes WebSocket control frames
- Now the connection stays open indefinitely

### 2. Chat Messages Not Reaching Client UI
**Symptom**: Browser console showed messages received but `chatId` was `undefined`, causing the client to ignore the message

**Root Cause**: The server sends `ServerEvent` messages with `data` field as a protobuf `Any` type. When marshaled to JSON by `protojson`, this becomes:
```json
{
  "typeUrl": "type.googleapis.com/mage.v1.ChatMessageData",
  "value": "CiEKB3RocmFpenoSBH..." // base64-encoded protobuf bytes
}
```

The client was passing this raw object to handlers, which expected `ChatMessageData` with a `chatId` field.

**Fix** (`mage-client-web/src/lib/stores/websocket.ts`):
- Added protobuf `Any` type detection (checks for `typeUrl` and `value` fields)
- Decodes base64 `value` to binary protobuf bytes
- Uses generated `ChatMessageDataCodec.decode()` to deserialize into proper TypeScript object
- Now handlers receive correctly decoded `ChatMessageData` with `chatId` populated

### 3. WebSocket Read Loop Implementation
**Before**: The server had no mechanism to process incoming WebSocket frames from the client

**After**: Added complete read/write separation:
- **Write loop** (main goroutine): Reads from `CallbackChan` and sends events to client
- **Ping goroutine**: Sends periodic PING messages to keep connection alive
- **Read goroutine** (NEW): Processes incoming messages and control frames (PONG)

## Changes Made

### Server (Go)
**File**: `mage-server-go/internal/server/websocket.go`

```go
// NEW: Read handler to process client messages and control frames
func (ws *WebSocketServer) readHandler(conn *websocket.Conn, done chan struct{}, sessionID string) {
    // Set pong handler to reset read deadline
    conn.SetPongHandler(func(string) error {
        ws.logger.Debug("received pong from client", zap.String("session", sessionID))
        conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })

    conn.SetReadDeadline(time.Now().Add(60 * time.Second))

    // Read loop - necessary to process control frames
    for {
        select {
        case <-done:
            return
        default:
            _, _, err := conn.ReadMessage()
            if err != nil {
                if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
                    ws.logger.Warn("websocket read error", zap.String("session", sessionID), zap.Error(err))
                }
                close(done)
                return
            }
        }
    }
}

// Start read handler in handleConnection()
go ws.readHandler(conn, done, sessionID)
```

### Client (TypeScript)
**File**: `mage-client-web/src/lib/stores/websocket.ts`

```typescript
// NEW: Decode protobuf Any types
function decodeChatMessageData(bytes: Uint8Array): ChatMessageData {
    const reader = new BinaryReader(bytes);
    return ChatMessageDataCodec.decode(reader);
}

// UPDATED: WebSocket message handler
ws.onmessage = (event) => {
    const serverEvent = JSON.parse(event.data) as ServerEvent;

    // Convert method string to enum
    if (typeof serverEvent.method === 'string') {
        serverEvent.method = callbackMethodFromJSON(serverEvent.method);
    }

    // NEW: Unpack protobuf Any type
    let eventData = serverEvent.data;
    if (eventData && typeof eventData === 'object' && 'typeUrl' in eventData && 'value' in eventData) {
        try {
            const base64Value = (eventData as any).value;
            const binaryString = atob(base64Value);
            const bytes = new Uint8Array(binaryString.length);
            for (let i = 0; i < binaryString.length; i++) {
                bytes[i] = binaryString.charCodeAt(i);
            }

            const typeUrl = (eventData as any).typeUrl;
            if (typeUrl === 'type.googleapis.com/mage.v1.ChatMessageData') {
                eventData = decodeChatMessageData(bytes);
            }
        } catch (err) {
            console.error('[WebSocket] Failed to decode Any type:', err);
        }
    }

    // Call handlers with decoded data
    const methodHandlers = handlers.get(serverEvent.method);
    if (methodHandlers) {
        methodHandlers.forEach((handler) => handler(eventData));
    }
};
```

## Testing

### Expected Behavior
1. **Connection**: WebSocket connects and stays open (no immediate disconnect)
2. **Lobby Chat**: Messages sent in lobby appear in all clients' UI
3. **Table Chat**: Messages sent at a table appear only for players at that table
4. **Filtering**: Each chat component only shows messages for its `chatId`

### Test Steps
1. Open two browser windows
2. Login as two different users
3. Go to lobby in both - test lobby chat
4. Join same table - test table chat
5. Verify messages appear in real-time
6. Check browser console - no errors about undefined `chatId`

## Architecture Notes

### Why We Need a Read Loop
WebSocket is a bidirectional protocol. Even if your application only sends messages from server to client, you MUST read from the connection to:
1. Process control frames (PING, PONG, CLOSE)
2. Detect connection failures
3. Reset read deadlines to prevent timeouts

Without reading, the connection will timeout when `SetReadDeadline()` expires.

### Protobuf Any Type
Protocol Buffers' `Any` type allows embedding arbitrary message types. When using `protojson` to marshal to JSON:
- `Any` becomes `{"typeUrl": "...", "value": "<base64>"}`
- The `value` is the protobuf binary representation, base64-encoded
- JavaScript/TypeScript must decode this back to the original message type

This is why we need the `decodeChatMessageData()` function - to reverse this encoding.

## Files Modified
- `mage-server-go/internal/server/websocket.go` - Added read handler
- `mage-client-web/src/lib/stores/websocket.ts` - Added protobuf Any decoding

## Status
✅ **COMPLETE** - Both lobby chat and table chat should now work correctly with real-time message delivery.
