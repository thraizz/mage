# Polling Implementation Summary

## Overview

Implemented smart periodic fetching across the lobby and table pages using the new `usePolling` utility. This provides fallback data loading when WebSocket is disconnected, with optimized performance for hidden tabs.

## Changes Made

### 1. Core Utility (`/src/lib/utils/polling.ts`)

Created a reusable polling utility with these features:

- **WebSocket-Aware**: Only polls when WebSocket is disconnected (by default)
- **Visibility Detection**: Reduces polling frequency when tab is hidden (6x slower)
- **Automatic Cleanup**: Stops polling on component unmount
- **Manual Refresh**: Provides user-triggered refresh capability
- **Configurable**: Flexible options for different use cases

### 2. Lobby Page (`/routes/(protected)/lobby/+page.svelte`)

**Changes:**
- Added `usePolling` import
- Setup polling for `loadTables()` - 5s visible / 30s hidden
- Setup polling for `loadOnlinePlayers()` - 10s visible / 60s hidden
- Connected `handleRefresh()` to `refreshTables()`

**Behavior:**
- WebSocket connected → No polling (real-time updates only)
- WebSocket disconnected → Polls every 5s (tables) and 10s (players)
- Tab hidden → Polls every 30s (tables) and 60s (players)
- User clicks refresh → Immediate fetch + reset timer

**Code:**
```ts
const { refresh: refreshTables } = usePolling(loadTables, {
  interval: 5000,
  intervalWhenHidden: 30000,
  pollWhenConnected: false,
  immediate: false
});

const { refresh: refreshPlayers } = usePolling(loadOnlinePlayers, {
  interval: 10000,
  intervalWhenHidden: 60000,
  pollWhenConnected: false,
  immediate: false
});
```

### 3. Table Page (`/routes/(protected)/table/[id]/+page.svelte`)

**Changes:**
- Added `usePolling` import
- Setup polling for `loadTable()` - 3s visible / 15s hidden
- Modified `loadTable()` to not show spinner on polling refreshes
- Polling only activates when WebSocket disconnected

**Behavior:**
- WebSocket connected → No polling (real-time updates only)
- WebSocket disconnected → Polls every 3s for table state
- Tab hidden → Polls every 15s
- Initial load → Shows loading spinner
- Polling refresh → Updates silently without spinner

**Code:**
```ts
const { refresh: refreshTable } = usePolling(loadTable, {
  interval: 3000,
  intervalWhenHidden: 15000,
  pollWhenConnected: false,
  immediate: false
});
```

### 4. Documentation

Created comprehensive guides:
- **POLLING_GUIDE.md**: Full usage guide with examples
- **POLLING_IMPLEMENTATION.md**: This file

## Polling Strategy by Component

| Component | Visible Interval | Hidden Interval | Poll When WS Connected |
|-----------|------------------|-----------------|------------------------|
| Lobby (tables) | 5s | 30s | No |
| Lobby (players) | 10s | 60s | No |
| Table (state) | 3s | 15s | No |
| Game (not impl) | - | - | - |

## Performance Impact

### Network Requests (WebSocket Connected)
- **Before**: No polling
- **After**: No polling (same - only uses WebSocket)

### Network Requests (WebSocket Disconnected)
- **Before**: No fallback - stale data
- **After**: Smart polling - data stays fresh

### Tab Hidden
- **Before**: N/A
- **After**: 6x slower polling (reduced server load)

### Browser Resources
- Single timeout per component
- Automatically cleaned up on unmount
- Pauses when not needed

## Testing Checklist

### Lobby Page
- [ ] Load page → should NOT poll while WebSocket connected
- [ ] Disconnect WebSocket → should start polling tables (5s) and players (10s)
- [ ] Hide tab → should slow down to 30s (tables) and 60s (players)
- [ ] Show tab → should fetch immediately and resume normal intervals
- [ ] Click refresh → should fetch immediately and reset timer
- [ ] Navigate away → should stop polling (no memory leaks)
- [ ] WebSocket reconnects → should stop polling

### Table Page
- [ ] Load page → should NOT poll while WebSocket connected
- [ ] Disconnect WebSocket → should start polling (3s)
- [ ] Hide tab → should slow down to 15s
- [ ] Show tab → should fetch immediately and resume 3s interval
- [ ] Initial load → should show loading spinner
- [ ] Polling refresh → should update silently (no spinner)
- [ ] Navigate away → should stop polling (no memory leaks)
- [ ] WebSocket reconnects → should stop polling

## Future Enhancements

1. **Exponential Backoff on Errors**: If fetch fails repeatedly, slow down polling
2. **User Preference**: Allow users to configure polling intervals
3. **Game Page Integration**: Add polling when game state API is implemented
4. **Polling Indicator**: Show subtle UI indicator when actively polling
5. **Network Status Detection**: Pause polling when offline
6. **Smart Intervals**: Adjust based on data change frequency

## Migration Guide for Future Components

1. Import the utility:
```ts
import { usePolling } from '$lib/utils/polling';
```

2. Setup polling:
```ts
const { refresh } = usePolling(yourFetchFunction, {
  interval: 5000,
  intervalWhenHidden: 30000,
  pollWhenConnected: false
});
```

3. Connect to manual refresh button:
```ts
<button onclick={refresh}>Refresh</button>
```

4. That's it! Cleanup is automatic.

## Troubleshooting

**Problem**: Polling doesn't stop when WebSocket connects
- Check that `pollWhenConnected: false` is set
- Verify WebSocket state in store

**Problem**: Too many requests
- Increase `interval` value
- Check that multiple pollers aren't running
- Verify cleanup on unmount

**Problem**: Stale data
- Decrease `interval` value
- Check that `enabled: true` (default)
- Verify fetch function is working

**Problem**: Memory leak
- Check that component is properly unmounted
- Verify no circular references in fetch function
