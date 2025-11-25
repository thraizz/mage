# Hybrid Update Strategy: WebSocket Events + Polling

## Problem Statement

The application needs to keep UI state synchronized with server state. However:

1. **Not all state changes are event-driven** (e.g., player ready status)
2. **WebSocket can disconnect** and needs a fallback
3. **Events can be missed** during reconnection
4. **Tab visibility** affects update frequency needs

## Solution: Hybrid Approach

### Three-Layer Strategy

```
┌─────────────────────────────────────────────────────────┐
│ Layer 1: WebSocket Events (Primary, Event-Driven)      │
│ - Instant updates when available                        │
│ - Low latency, efficient                                │
│ - Handles: Chat, table creation/deletion, game events   │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ Layer 2: Smart Polling (Safety Net, Always On)         │
│ - Catches state changes not covered by events           │
│ - Slower interval (5-10s) to avoid overhead             │
│ - Handles: Ready status, player joins/leaves            │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ Layer 3: Manual Refresh (User Control)                 │
│ - Instant fetch on user action                          │
│ - Bypasses all timers                                    │
│ - Handles: User wants fresh data NOW                     │
└─────────────────────────────────────────────────────────┘
```

## Implementation by Component

### Lobby Page

**What Updates:**
- Table list (created, deleted, player count changes)
- Online players list

**How:**
- **Events**: `TABLE_WAITING` for table updates
- **Polling**: 5s (tables), 10s (players) - always on as safety net
- **Why Both**: Table player counts and states may change without events

```ts
// WebSocket: instant updates for table changes
const unsubscribe = subscribeLobbyUpdates(handleTableUpdate);

// Polling: catches missed updates and non-evented changes
const { refresh } = usePolling(loadTables, {
  interval: 5000,
  pollWhenConnected: true  // ALWAYS poll as safety net
});
```

### Table Page

**What Updates:**
- Player ready status
- Player joins/leaves
- Game starting

**How:**
- **Events**: `TABLE_WAITING` for table state updates
- **Polling**: 5s always on - critical for ready status
- **Why Both**: Ready status changes aren't reliably event-driven

```ts
// WebSocket: instant updates when table state changes
const unsubscribe = subscribeTableUpdates(tableId, handleTableUpdate);

// Polling: ensures ready status stays fresh
const { refresh } = usePolling(loadTable, {
  interval: 5000,
  pollWhenConnected: true  // ALWAYS poll - ready status needs this
});
```

## Polling Configuration Philosophy

### Always Poll (pollWhenConnected: true)

**Use When:**
- State changes aren't fully event-driven
- Critical UI updates that can't be missed
- Small payload, low server cost

**Examples:**
- Table page (ready status)
- Lobby page (table list as safety net)

### Poll Only as Fallback (pollWhenConnected: false)

**Use When:**
- All state changes have reliable WebSocket events
- Large payload or expensive queries
- Non-critical data

**Examples:**
- (Future) Game replays list
- (Future) Tournament brackets (if fully event-driven)

## Performance Characteristics

### Network Request Frequency

| Component | WS Connected | WS Disconnected | Tab Hidden (WS Connected) |
|-----------|--------------|-----------------|---------------------------|
| Lobby (tables) | 1 req / 5s | 1 req / 5s | 1 req / 30s |
| Lobby (players) | 1 req / 10s | 1 req / 10s | 1 req / 60s |
| Table (state) | 1 req / 5s | 1 req / 5s | 1 req / 30s |

### Latency Comparison

| Update Type | WebSocket Event | Polling (5s) | Manual Refresh |
|-------------|-----------------|--------------|----------------|
| Latency | ~50-100ms | 0-5000ms (avg 2.5s) | ~200-500ms |
| Reliability | 95% | 100% | 100% |
| User Control | No | No | Yes |

## Event-Driven vs Polling

### Currently Event-Driven (WebSocket)

✅ **Fully Covered:**
- Chat messages (`CHATMESSAGE`)
- Table creation (`TABLE_WAITING`)
- Table deletion (`TABLE_WAITING`)
- Game start (`START_GAME`)
- Game state updates (`GAME_UPDATE`)

⚠️ **Partially Covered:**
- Player ready status (inconsistent `TABLE_WAITING` events)
- Player joins/leaves (may miss during reconnect)

### Currently Polling-Dependent

These require polling because they're not reliably event-driven:

- Player ready status changes
- Online player count (as safety net)
- Table player count verification

## Migration Path to Full Event-Driven

### Phase 1 (Current)
- Hybrid: Events + polling
- Polling always on at 5-10s
- Safe but slightly inefficient

### Phase 2 (Future - Server Improvements)
- Server sends `TABLE_WAITING` on ALL player state changes
- Server sends `PLAYER_READY_CHANGED` event
- Server sends `PLAYER_JOINED` / `PLAYER_LEFT` events

### Phase 3 (Future - Client Optimization)
- Switch to `pollWhenConnected: false`
- Reduce intervals to 30s+ (rare safety checks)
- Full event-driven architecture

## Testing Scenarios

### Scenario 1: Normal Operation (WebSocket Connected)
1. Player A clicks "Ready" → WebSocket event → Instant update
2. 5 seconds pass → Polling refresh → Confirms state (safety net)
3. Player B clicks "Ready" → WebSocket event → Instant update

**Result**: Users see instant updates, polling confirms accuracy

### Scenario 2: WebSocket Disconnected
1. Player A clicks "Ready" → No WebSocket event
2. 2 seconds pass → Polling refresh → Updates for Player A & B
3. WebSocket reconnects → Events resume

**Result**: Max 5s delay, but data stays fresh

### Scenario 3: Missed Event (Network Glitch)
1. Player A clicks "Ready" → WebSocket event sent but dropped by network
2. 3 seconds pass → Polling refresh → Catches Player A ready status
3. UI updates correctly

**Result**: Polling catches the missed state change

### Scenario 4: Tab Hidden
1. User switches to another tab
2. Polling slows to 30s (6x slower)
3. User returns to tab → Immediate refresh
4. Polling returns to 5s

**Result**: Reduced server load, instant update on return

## Best Practices

### DO:
✅ Always poll for critical state that affects user actions
✅ Use longer intervals (5-10s) when also using WebSocket
✅ Slow down significantly when tab hidden (6x or more)
✅ Provide manual refresh for user control
✅ Log WebSocket events vs polling updates for debugging

### DON'T:
❌ Poll more frequently than 3s (excessive server load)
❌ Disable polling just because WebSocket is connected (events aren't 100% reliable)
❌ Use same interval for hidden tabs (wastes resources)
❌ Forget to clean up subscriptions on unmount
❌ Show loading spinners for polling refreshes (jarring UX)

## Future Improvements

### Short Term
1. **Event Monitoring**: Log which updates come from events vs polling
2. **Adaptive Intervals**: Increase interval if data rarely changes
3. **Error Handling**: Exponential backoff on repeated poll failures

### Medium Term
1. **Server Events**: Add `PLAYER_READY_CHANGED`, `PLAYER_JOINED_TABLE`, etc.
2. **Event Acknowledgment**: Confirm events received, request resync if missed
3. **Differential Updates**: Poll only returns what changed since last fetch

### Long Term
1. **Full Event-Driven**: All state changes have reliable events
2. **Rare Polling**: 30-60s intervals only as safety net
3. **Optimistic Updates**: Update UI immediately, confirm with server
