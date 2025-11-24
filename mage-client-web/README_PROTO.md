# Protocol Buffers Setup

This web client uses TypeScript types generated from the **same protobuf files** as the Go server.

## Source of Truth

**Proto files location**: `../mage-server-go/api/proto/mage/v1/`

The server's proto files are the single source of truth for the API contract. This ensures:

- Type safety between client and server
- No API drift or version mismatches
- Automatic updates when server API changes

## Generated Code

**Output location**: `src/lib/generated/mage/v1/`

Contains TypeScript interfaces and types for:

- All 70 RPC methods (auth, room, table, game, tournament, draft, chat, admin, replay)
- Data models (GameView, TableView, CardView, TournamentView, etc.)
- WebSocket event types (86 event definitions)
- Request/Response message types

## Regenerating Types

After any proto file changes:

```bash
npm run proto:generate
```

This runs `scripts/generate-proto.sh` which:

1. Reads proto files from `../mage-server-go/api/proto/mage/v1/`
2. Generates TypeScript code using `ts-proto`
3. Outputs to `src/lib/generated/mage/v1/`

## Configuration

Generation script uses these `ts-proto` options:

- `outputServices=generic-definitions` - Generate type definitions (no gRPC client code)
- `env=browser` - Browser-compatible code
- `useOptionals=messages` - Optional fields for message fields
- `outputIndex=true` - Generate index.ts for easy imports

## Client Usage

Use the `MageClient` wrapper for type-safe RPC calls:

```typescript
import { getMageClient } from '$lib/grpc/client';

const client = getMageClient();

// Type-safe login
const response = await client.connectUser('username', 'password');
if (response.success) {
	console.log('Session:', response.sessionId);
}

// Type-safe server state
const state = await client.getServerState();
console.log('Players:', state.serverState?.activePlayers);
```

For advanced usage, import generated types directly:

```typescript
import type { RoomCreateTableRequest, RoomCreateTableResponse } from '$lib/generated/mage/v1/table';
import type { MatchOptions } from '$lib/generated/mage/v1/models';

const request: RoomCreateTableRequest = {
	sessionId: client.getSessionId()!,
	roomId: lobbyId,
	matchOptions: {
		name: 'My Game',
		gameType: 'TwoPlayerDuel',
		deckType: 'Constructed',
		winsNeeded: 2
		// TypeScript will autocomplete all valid fields!
	}
};

const response = await client.call<RoomCreateTableRequest, RoomCreateTableResponse>(
	'RoomCreateTable',
	request
);
```

## Files Generated

- `admin.ts` - Admin operations (9 methods)
- `auth.ts` - Authentication & connection (7 methods)
- `chat.ts` - Chat messaging (7 methods)
- `draft.ts` - Booster draft (5 methods)
- `game.ts` - Game execution & replay (18 methods)
- `models.ts` - Data models (GameView, TableView, CardView, etc.)
- `room.ts` - Lobby operations (5 methods)
- `server.ts` - Service definitions (all 70 methods)
- `table.ts` - Table management (10 methods)
- `tournament.ts` - Tournament operations (4 methods)
- `websocket.ts` - WebSocket events (86 event types)

## Version Control

Generated files ARE committed to git for convenience, but can be regenerated anytime from the server proto files.

## Troubleshooting

### "Cannot find module" errors

Ensure types are generated:

```bash
npm run proto:generate
```

### Type mismatches

Regenerate both server and client code from the same proto files:

```bash
# Server
cd ../mage-server-go && make proto

# Client
cd ../mage-client-web && npm run proto:generate
```

### Import errors

The generated code is in `src/lib/generated/mage/v1/`, not `src/lib/generated/`:

```typescript
// Correct ✓
import type { GameView } from '$lib/generated/mage/v1/models';

// Wrong ✗
import type { GameView } from '$lib/generated/models';
```

## See Also

- `PROTOBUF_INFRASTRUCTURE.md` in project root - Complete documentation
- `mage-server-go/api/proto/mage/v1/` - Source proto files
- `scripts/generate-proto.sh` - Generation script
