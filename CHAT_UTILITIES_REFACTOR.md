# Chat Utilities Refactor Summary

## Overview
Extracted common chat message handling logic from `LobbyChat.svelte` and `TableChat.svelte` into a shared TypeScript utility module for reusability and maintainability.

## New File Created

### `mage-client-web/src/lib/utils/chat.ts`

A comprehensive utility module with the following exports:

#### Functions

1. **`extractTimestamp(protoTime)`**
   - Handles multiple timestamp formats from protobuf messages
   - Supports: `Date`, `number`, protobuf `Timestamp` object
   - Returns: Unix timestamp in milliseconds
   - Fallback: `Date.now()` if undefined/null

2. **`convertProtoMessageToClientMessage(protoMessage)`**
   - Converts protobuf `ChatMessage` to client `ChatMessage` type
   - Handles timestamp extraction
   - Determines message type (system vs user)
   - Generates unique message ID
   - Returns ready-to-use client message object

3. **`isSystemMessage(protoMessage)`**
   - Checks if a message is from the system
   - Returns boolean

4. **`formatMessageTime(timestamp)`**
   - Formats Unix timestamp as HH:MM (24-hour format)
   - Used for message display in UI

5. **`parseWhisperCommand(content)`**
   - Parses `/w username message` whisper commands
   - Returns object with: `{ isWhisper, username?, message? }`

6. **`validateWhisperCommand(whisperInfo, currentUsername?)`**
   - Validates parsed whisper command
   - Checks for proper format
   - Prevents self-whispers
   - Returns error message string or `null` if valid

#### Classes

**`RateLimiter`**
- Configurable rate limiting for chat messages
- Default: 10 messages per 60 seconds
- Methods:
  - `isLimited()`: Check if currently rate limited
  - `recordMessage()`: Record a sent message
  - `getCooldownSeconds()`: Get remaining cooldown time
  - `reset()`: Clear all recorded timestamps

## Files Refactored

### `LobbyChat.svelte`
**Before**: 330 lines with duplicated utility functions
**After**: ~280 lines using shared utilities

**Changes**:
- Removed: `formatTime()`, `isRateLimited()`, `getCooldownSeconds()`, `recordMessageSent()`, `parseWhisperCommand()`
- Replaced manual timestamp handling with `convertProtoMessageToClientMessage()`
- Replaced inline rate limiting logic with `RateLimiter` class
- Updated whisper validation to use `validateWhisperCommand()`
- Updated time display to use `formatMessageTime()`

### `TableChat.svelte`
**Before**: 290 lines with duplicated utility functions
**After**: ~240 lines using shared utilities

**Changes**:
- Removed: `formatTime()`, `isRateLimited()`, `getCooldownSeconds()`, `recordMessageSent()`
- Replaced manual timestamp handling with `convertProtoMessageToClientMessage()`
- Replaced inline rate limiting logic with `RateLimiter` class
- Updated time display to use `formatMessageTime()`

## Benefits

### 1. **DRY (Don't Repeat Yourself)**
- Eliminated ~100 lines of duplicated code across components
- Single source of truth for chat message handling

### 2. **Maintainability**
- Bug fixes and improvements now apply to all chat components
- Easier to add new chat features (e.g., new message types, formatting)

### 3. **Testability**
- Utilities can be unit tested independently
- Components have less logic to test

### 4. **Reusability**
- Easy to add new chat components (e.g., DM chat, tournament chat)
- `RateLimiter` class can be used for other rate-limited features

### 5. **Type Safety**
- Centralized TypeScript types and interfaces
- Proper handling of protobuf timestamp edge cases

## Usage Example

```typescript
import {
	convertProtoMessageToClientMessage,
	formatMessageTime,
	parseWhisperCommand,
	validateWhisperCommand,
	RateLimiter
} from '$lib/utils/chat';

// Convert protobuf message
const clientMessage = convertProtoMessageToClientMessage(protoMessage);

// Format timestamp for display
const timeString = formatMessageTime(clientMessage.timestamp); // "14:30"

// Rate limiting
const rateLimiter = new RateLimiter(10, 60000); // 10 msgs per minute
if (rateLimiter.isLimited()) {
	const seconds = rateLimiter.getCooldownSeconds();
	console.log(`Wait ${seconds} seconds`);
} else {
	// Send message
	rateLimiter.recordMessage();
}

// Whisper command handling
const whisperInfo = parseWhisperCommand('/w alice hello');
const error = validateWhisperCommand(whisperInfo, 'bob');
if (error) {
	console.error(error);
}
```

## Future Enhancements

### Potential Additions
1. **Message Filtering**
   - Filter by user, date range, message type
   - Search functionality

2. **Message Formatting**
   - Markdown support
   - Emoji parsing
   - URL detection and linking

3. **Enhanced Rate Limiting**
   - Different limits for different user roles
   - Burst allowance with token bucket algorithm

4. **Message History**
   - Pagination helpers
   - Load more messages
   - Message caching

5. **Notification Helpers**
   - Unread message counting
   - Mention detection (@username)
   - Sound notification triggers

## Testing Recommendations

### Unit Tests to Add
```typescript
// chat.test.ts
describe('extractTimestamp', () => {
  it('should handle Date objects', () => {
    const date = new Date('2025-01-01T12:00:00Z');
    expect(extractTimestamp(date)).toBe(date.getTime());
  });

  it('should handle numbers', () => {
    const timestamp = 1704110400000;
    expect(extractTimestamp(timestamp)).toBe(timestamp);
  });

  it('should fallback to Date.now() for invalid input', () => {
    const result = extractTimestamp(undefined);
    expect(result).toBeCloseTo(Date.now(), -2);
  });
});

describe('RateLimiter', () => {
  it('should allow messages under limit', () => {
    const limiter = new RateLimiter(3, 1000);
    limiter.recordMessage();
    limiter.recordMessage();
    expect(limiter.isLimited()).toBe(false);
  });

  it('should block messages over limit', () => {
    const limiter = new RateLimiter(2, 1000);
    limiter.recordMessage();
    limiter.recordMessage();
    expect(limiter.isLimited()).toBe(true);
  });
});
```

## Migration Checklist

✅ Create `src/lib/utils/chat.ts` with all utilities
✅ Refactor `LobbyChat.svelte` to use utilities
✅ Refactor `TableChat.svelte` to use utilities
✅ Test lobby chat functionality
✅ Test table chat functionality
✅ Verify rate limiting works correctly
✅ Verify timestamp formatting displays correctly

## Files Modified
- ✅ `mage-client-web/src/lib/utils/chat.ts` (NEW)
- ✅ `mage-client-web/src/lib/components/LobbyChat.svelte`
- ✅ `mage-client-web/src/lib/components/TableChat.svelte`

## Status
✅ **COMPLETE** - All chat components now use shared utilities for message handling, rate limiting, and formatting.
