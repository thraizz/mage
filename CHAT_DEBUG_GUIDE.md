# Chat Debug Guide

## Quick Diagnostic Steps

### Step 1: Check WebSocket Connection

1. Open browser console (F12 → Console)
2. Go to lobby
3. Look for:
   ```
   [LobbyChat] Got chat ID: room:xxx
   [LobbyChat] Joined chat
   [LobbyChat] Connected and listening for messages
   ```

**Expected**: All three messages should appear
**If missing**: WebSocket initialization failed

### Step 2: Check Server-Side Join

Look at Go server logs for:
```
{"level":"info","msg":"user joined chat room","room_id":"room:xxx","username":"your_username","total_users":1,"users":["your_username"]}
```

**Expected**: You should see your username in the users array
**If missing**: ChatJoin RPC failed or user not authenticated

### Step 3: Send a Test Message

Type a message and click send. Check:

**Browser Console**:
```
[LobbyChat] Sending message: test
[LobbyChat] Sending regular message
[LobbyChat] Message sent successfully
```

**Server Logs**:
```
{"level":"info","msg":"chat message received","room_id":"room:xxx","username":"your_username","message":"test","users_in_room":1}
{"level":"debug","msg":"broadcasted chat message","room_id":"room:xxx","username":"your_username","users":1,"sessions_sent":1}
```

**Expected**: Message reaches server and broadcast completes
**If server logs missing**: gRPC call failed
**If broadcast shows 0 sessions_sent**: User not in room or session not connected

### Step 4: Check Message Receipt

After sending, browser console should show:
```
[LobbyChat] Received WebSocket message: {...}
[LobbyChat] Parsed as ChatMessageData: {...}
[LobbyChat] Current chatId: room:xxx Message chatId: room:xxx
[LobbyChat] Proto message: {...}
[LobbyChat] Adding message to list: {...}
```

**Expected**: Message comes back via WebSocket
**If missing**: WebSocket not receiving, or message filtered out

## Common Issues & Fixes

### Issue 1: No WebSocket Connection

**Symptoms**:
- No `[LobbyChat] Connected and listening` message
- Error in console about failed to connect

**Fix**:
- Check WebSocket server is running on port 17179
- Check VITE_WEBSOCKET_URL environment variable
- Ensure session ID is valid

### Issue 2: User Not Joined to Chat Room

**Symptoms**:
- Server logs show `users_in_room: 0`
- Broadcast shows `sessions_sent: 0`

**Fix**:
- Check ChatJoin RPC completed successfully
- Verify username is set in session (not guest without username)
- Check `joinChat(chatId)` was called

### Issue 3: Messages Filtered Out

**Symptoms**:
- Server broadcasts message (sessions_sent > 0)
- Client receives WebSocket message
- But logs show "Ignoring message for different chat"

**Fix**:
- chatId mismatch between component and message
- Check the chatId format: should be `room:xxx` for lobby
- Verify both InitializeChat and SendMessage use same chatId

### Issue 4: WebSocket Data Format Mismatch

**Symptoms**:
- WebSocket receives data
- Error parsing ChatMessageData
- Message not displayed

**Fix**:
- Check proto definitions match between client and server
- Verify `data` field contains ChatMessageData, not raw ChatMessage
- Check if Any type is being unpacked correctly

## Manual Testing Procedure

### Test 1: Single User Echo

1. Login as one user
2. Go to lobby
3. Send message "test 1"
4. **Expected**: Message appears in your own chat

**Pass criteria**: You see your own message

### Test 2: Multi-User Broadcast

1. Open two browser windows
2. Login as different users (or use guest twice with different browsers)
3. Both go to lobby
4. User A sends "hello from A"
5. **Expected**: Both users see the message

**Pass criteria**: Message appears in both windows

### Test 3: WebSocket Reconnection

1. Login and go to lobby
2. Stop Go server
3. Restart Go server
4. Wait for reconnection (should happen automatically)
5. Send message
6. **Expected**: Message works after reconnect

**Pass criteria**: Automatic reconnection and continued functionality

## Server Log Levels

For maximum debugging, ensure server config has:
```yaml
logging:
  level: debug  # Shows all broadcast details
```

Or set via environment:
```bash
export LOG_LEVEL=debug
```

## Client Dev Mode

For detailed client logs:
```bash
cd mage-client-web
npm run dev
```

This enables hot reload and shows all console logs.

## Quick Fix Checklist

If chat doesn't work, verify IN THIS ORDER:

- [ ] Server is running and accessible
- [ ] WebSocket server is running on :17179
- [ ] User is authenticated (has username in session)
- [ ] User joined the chat room (ChatJoin succeeded)
- [ ] WebSocket connection is established (check browser console)
- [ ] Component subscribed to CHATMESSAGE events
- [ ] ChatId format matches (room:xxx for lobby)
- [ ] Server received the message (check server logs)
- [ ] Server broadcasted the message (sessions_sent > 0)
- [ ] Client received WebSocket event
- [ ] Client didn't filter out the message
- [ ] Svelte reactivity triggered ($state update)

## Debug Commands

### Check WebSocket in Browser Console

```javascript
// Check if WebSocket is connected
$websocketStore.state

// Check registered handlers
// (This won't work directly, but you can check in React DevTools or similar)
```

### Check Server State

```bash
# View active sessions
curl http://localhost:17171/health  # if you have health endpoint

# Check chat room state (via debugging)
# Add temporary endpoint to dump chat room state
```

## Next Steps If Still Broken

1. Share the EXACT browser console output (all [LobbyChat] lines)
2. Share the EXACT server logs (all chat-related lines)
3. Check Network tab: does ChatJoin RPC return success:true?
4. Check Network tab: does ChatSendMessage RPC return success:true?
5. Check WS tab in Network: are messages being sent over WebSocket?
