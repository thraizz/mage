# TableSetReady Fix Summary

## Problem

The web client was attempting to call a non-existent gRPC method `TableSetReady`, resulting in the error:
```
Method not supported in HTTP/JSON mode: TableSetReady
```

## Root Cause

The `toggleReady()` function in `mage-client-web/src/lib/api/table.ts` was calling a `TableSetReady` RPC that doesn't exist in the proto definitions.

After inspecting the proto files (`mage-server-go/api/proto/mage/v1/`), the available table-related RPCs are:
- `RoomGetAllTables`
- `RoomGetTableById`
- `RoomCreateTable`
- `RoomJoinTable`
- `RoomLeaveTableOrTournament`
- `RoomWatchTable`
- `TableSwapSeats`
- `TableRemove`
- `TableIsOwner`

**No `TableSetReady` method exists.**

## How XMage Handles Ready Status

In XMage, ready status is controlled by **deck submission**, not a separate ready toggle:
1. Players join a table via `RoomJoinTable` (with deck list)
2. Players become "ready" when they submit their deck via `DeckSubmit`
3. When all players are ready, the host can start the game

## Solution

Modified `toggleReady()` function to:
1. Throw a clear error message explaining that ready status is controlled by deck submission
2. Removed the non-functional RPC call attempt
3. Removed unused imports (`TableSetReadyRequest`, `TableSetReadyResponse`)

## Code Changes

### File: `mage-client-web/src/lib/api/table.ts`

**Before:**
```typescript
export async function toggleReady(tableId: string, isReady: boolean): Promise<void> {
	const client = getMageClient();
	// ... session setup ...

	try {
		const request: TableSetReadyRequest = {
			sessionId,
			roomId: roomResponse.roomId,
			tableId,
			ready: isReady
		};

		const response = await client.call<TableSetReadyRequest, TableSetReadyResponse>(
			'TableSetReady',  // ❌ This RPC doesn't exist
			request
		);
		// ...
	} catch (error) {
		console.warn('TableSetReady not available:', error);
	}
}
```

**After:**
```typescript
export async function toggleReady(tableId: string, isReady: boolean): Promise<void> {
	// TableSetReady RPC doesn't exist - ready status is controlled by deck submission
	// For now, throw an error to make it clear this functionality isn't implemented
	throw new Error('Ready status is controlled by deck submission. Please submit your deck to mark yourself as ready.');
}
```

## Impact

- ✅ **Error eliminated**: No more "Method not supported" errors
- ✅ **Clear messaging**: Users now see a helpful error message explaining how ready status works
- ✅ **Reduced code complexity**: Removed unnecessary session/room lookup code
- ✅ **Type safety**: Removed unused type imports

## Future Implementation

To properly implement ready status in the web client:

### Option 1: Implement TableSetReady RPC (Server-side)
1. Add `TableSetReady` RPC to `mage-server-go/api/proto/mage/v1/server.proto`
2. Define request/response messages in `table.proto`
3. Implement handler in `mage-server-go/internal/server/`
4. Run `make proto` to regenerate code
5. Update client to use the new RPC

### Option 2: Use Deck Submission (Current XMage approach)
1. Implement deck builder UI
2. Add `DeckSubmit` RPC call when player selects deck
3. Remove ready toggle button, replace with "Submit Deck" button
4. Server marks player as ready upon deck submission

**Recommendation**: Follow Option 2 to maintain compatibility with XMage's design.

## Files Modified

- `mage-client-web/src/lib/api/table.ts` - Updated `toggleReady()` function and removed unused imports

## Testing

The error no longer occurs when:
1. Users view the table lobby page (`/table/[id]`)
2. The `toggleReady` function is called (though it now throws a helpful error)

The fix preserves the UI button but provides clear feedback that the feature requires deck submission implementation.
