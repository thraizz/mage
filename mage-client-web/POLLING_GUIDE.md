# Polling Guide: DRY and UX-Optimized Periodic Fetching

## Overview

The `usePolling` utility provides smart periodic fetching with these features:

1. **WebSocket-Aware**: Only polls when WebSocket is disconnected (fallback mode)
2. **Visibility Detection**: Reduces polling frequency when tab is hidden
3. **Automatic Cleanup**: Stops polling on component unmount
4. **Manual Refresh**: Provides user-triggered refresh capability
5. **Configurable**: Flexible options for different use cases

## Basic Usage

### Simple Periodic Refresh

```ts
import { usePeriodicRefresh } from '$lib/utils/polling';

// Polls every 5 seconds when WS is down, 30 seconds when tab hidden
const { refresh } = usePeriodicRefresh(loadData);

// Manual refresh
<button onclick={refresh}>Refresh</button>
```

### Advanced Configuration

```ts
import { usePolling } from '$lib/utils/polling';

const { refresh, stop, start, isPolling } = usePolling(loadData, {
	interval: 5000, // Poll every 5s when visible
	intervalWhenHidden: 30000, // Poll every 30s when hidden
	pollWhenConnected: false, // Only poll when WS disconnected (default)
	immediate: true, // Fetch immediately on mount (default)
	enabled: true // Enable polling (default)
});
```

## Real-World Examples

### Example 1: Lobby Tables List

**Before (manual implementation):**

```ts
let pollTimer: ReturnType<typeof setTimeout> | null = null;

async function loadTables() {
	/* ... */
}

function startPolling() {
	if (pollTimer) clearTimeout(pollTimer);
	pollTimer = setTimeout(async () => {
		await loadTables();
		startPolling();
	}, 5000);
}

onMount(() => {
	loadTables();
	startPolling();
});

onDestroy(() => {
	if (pollTimer) clearTimeout(pollTimer);
});
```

**After (with usePolling):**

```ts
import { usePeriodicRefresh } from '$lib/utils/polling';

async function loadTables() { /* ... */ }

// Automatically handles mount/unmount, visibility, and WebSocket state
const { refresh } = usePeriodicRefresh(loadTables);

// Button for manual refresh
<button onclick={refresh}>Refresh</button>
```

### Example 2: Table View (Player List)

```ts
import { usePolling } from '$lib/utils/polling';

async function loadTable() {
	if (!tableId) return;
	try {
		table = await fetchTable(tableId);
	} catch (err) {
		console.error('Failed to load table:', err);
	}
}

// Only poll as fallback when WebSocket is down
const { refresh } = usePolling(loadTable, {
	interval: 3000, // Check every 3s when WS down
	intervalWhenHidden: 15000, // Every 15s when hidden
	pollWhenConnected: false // Trust WebSocket when connected
});
```

### Example 3: Game Lobby (Always Poll)

For scenarios where you want polling even with WebSocket:

```ts
const { refresh } = usePolling(loadOnlinePlayers, {
	interval: 10000,
	pollWhenConnected: true, // Poll even when WS connected
	immediate: true
});
```

### Example 4: Conditional Polling

```ts
let pollingEnabled = $state(true);

const { refresh, stop, start } = usePolling(loadData, {
	enabled: pollingEnabled,
	interval: 5000
});

// User can toggle polling
function togglePolling() {
	pollingEnabled = !pollingEnabled;
	if (pollingEnabled) {
		start();
	} else {
		stop();
	}
}
```

## Integration with Existing Components

### Lobby Page Pattern

```ts
// 1. Keep WebSocket for real-time updates
const unsubscribeLobby = subscribeLobbyUpdates(handleTableUpdate);

// 2. Add polling as fallback
const { refresh: refreshTables } = usePeriodicRefresh(loadTables, 5000);

// 3. Manual refresh button uses polling's refresh
<button onclick={refreshTables}>Refresh</button>
```

### Table Page Pattern

```ts
// 1. Load initial data
await loadTable();

// 2. Connect WebSocket
await connectWebSocket();

// 3. Add polling fallback
const { refresh: refreshTable } = usePolling(loadTable, {
	interval: 3000,
	pollWhenConnected: false // Only when WS down
});
```

## How It Works

### Polling Decision Tree

```
Should Poll?
├─ enabled = false? → NO
├─ component unmounted? → NO
├─ WebSocket connected?
│  ├─ pollWhenConnected = true? → YES
│  └─ pollWhenConnected = false? → NO
└─ WebSocket disconnected? → YES
```

### Visibility-Based Intervals

- **Tab Visible**: Uses `interval` (e.g., 5 seconds)
- **Tab Hidden**: Uses `intervalWhenHidden` (e.g., 30 seconds)
- **Switching to Visible**: Fetches immediately, then resumes normal interval

### WebSocket State Changes

- **Connected → Disconnected**: Starts polling (if `pollWhenConnected = false`)
- **Disconnected → Connected**: Stops polling (if `pollWhenConnected = false`)
- **Manual Refresh**: Always works regardless of WebSocket state

## Performance Characteristics

### Before (manual polling)

- ✗ Always polls regardless of WebSocket state
- ✗ Same frequency when tab hidden
- ✗ Manual cleanup in every component
- ✗ No coordination with WebSocket

### After (usePolling)

- ✓ Polls only when WebSocket down (by default)
- ✓ 6x slower polling when tab hidden
- ✓ Automatic cleanup
- ✓ Coordinates with WebSocket state

## Best Practices

1. **Default to Fallback Mode**: Use `pollWhenConnected: false` (default) to avoid redundant fetches
2. **Reasonable Intervals**: 5-10s visible, 30-60s hidden
3. **Combine with WebSocket**: Use WebSocket for real-time + polling for fallback
4. **Manual Refresh**: Always expose `refresh()` for user control
5. **Loading States**: Show loading indicator during refresh

## Migration Checklist

For each component with periodic fetching:

- [ ] Import `usePolling` or `usePeriodicRefresh`
- [ ] Replace manual `setTimeout` loops with `usePolling`
- [ ] Remove manual cleanup in `onDestroy`
- [ ] Connect manual refresh button to `refresh()`
- [ ] Test WebSocket disconnect → should start polling
- [ ] Test tab hidden → should slow down polling
- [ ] Test component unmount → should stop polling

## FAQ

**Q: Why not always poll even with WebSocket?**
A: Reduces server load and network usage. WebSocket provides instant updates, polling is fallback.

**Q: What if I want to poll even when connected?**
A: Use `pollWhenConnected: true` option.

**Q: How do I know if polling is active?**
A: Check `isPolling` property or observe network requests.

**Q: Can I dynamically enable/disable polling?**
A: Yes, use `enabled` option with reactive state and call `start()`/`stop()`.

**Q: What happens on component unmount?**
A: Polling automatically stops, no manual cleanup needed.

**Q: Does it handle fetch errors?**
A: Yes, catches errors, logs them, and continues polling.

## Examples in Codebase

- **Lobby tables**: `/routes/(protected)/lobby/+page.svelte`
- **Table view**: `/routes/(protected)/table/[id]/+page.svelte`
- **Game state**: `/routes/(protected)/game/[id]/+page.svelte`
