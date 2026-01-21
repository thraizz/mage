# Code Extraction Plan: Slim Down Game & Playtest Pages

**Goal:** Extract duplicated script logic into reusable utilities/composables to reduce page complexity by 25-30% (300-400 lines per page).

**Current State:**
- Playtest page: 1,992 lines
- Game page: 3,554 lines
- Identified duplication: ~600-800 lines of extractable patterns

**Target State:**
- Playtest page: ~1,400-1,600 lines
- Game page: ~2,500-2,800 lines
- New utility modules: 6-8 files (~800 lines total)

---

## Phase 0: Documentation Discovery ✅

### Existing Infrastructure Audit

**Utility Files (`$lib/utils/`):**
- `drag-drop.ts` (358 lines) - Complete drag/drop system with store, derived values, zone registration
- `playtest-helpers.ts` (303 lines) - Zone mutations, card lookups, player updates
- `zones.ts` (151 lines) - Zone enum/mapping layer
- `game-api.ts` (77 lines) - Session management wrapper
- `jwt.ts` (80+ lines) - Token handling
- Others: auth-guard, chat, deck-parser, grpc-errors, polling, scryfall, session-error-handler

**Store Files (`$lib/stores/`):**
- `game.ts` (1,014 lines) - 27 derived stores, 24 methods, WebSocket subscriptions
- `playtest-game.ts` (1,257 lines) - 30+ methods, localStorage persistence, 7 derived stores
- `combat.ts` (562 lines) - Combat phase state machine, 11 derived stores
- `websocket.ts` (431 lines) - Connection + protobuf decoding
- `connection.ts` (476 lines) - Health checks + reconnection
- Others: auth, toast, confirm, visual-stack

**Key Findings:**
1. Drag-drop infrastructure is already well-factored (no extraction needed)
2. Zone utilities exist but are underutilized (playtest-helpers.ts has reusable patterns)
3. No keyboard shortcut infrastructure exists (100% inline duplication)
4. No action handler wrapper exists (async/loading/error pattern repeats 8+ times per page)
5. Dialog state management is repeated but could use composable pattern
6. Clipboard operations are duplicated with fallback logic

### Allowed APIs (Verified from Source)

**PlaytestGameStore Methods (client-only):**
```typescript
// From $lib/stores/playtest-game.ts (lines 1188-1225)
initialize(gameId, players, options?)
drawCards(playerId, count)
moveCardToZone(cardId, zone)
tapCard(cardId, tapped)
untapAll(playerId)
modifyLife(playerId, delta)
setPlayerCounter(playerId, type, value)
shuffleLibrary(playerId)
createToken(name, types, power, toughness, color)
addCounter/removeCounter(cardId, name, amount)
nextTurn()
mulligan/keepHand(playerId)
```

**DragDropStore API:**
```typescript
// From $lib/utils/drag-drop.ts (lines 1-358)
dragDropStore.startDrag(cardId, cardName, sourceZone, x, y, validZones)
dragDropStore.registerDropZone(config: DropZoneConfig)
// Derived stores:
$isDragging, $draggedCardId, $draggedCardName, $dragPosition,
$isOverValidDropZone, $currentDropZone

// Utilities:
isDragThresholdMet(start, current, threshold=5)
getAllValidDropZones(sourceZone)
```

**Component Props (from component source):**
```typescript
// PlayerInfoRow.svelte
{ player, graveyard, exile, mana, showLifeMenu,
  onLifeChange, onPoisonChange, onToggleLifeMenu, onDeckContextMenu,
  libraryDropZoneRef, graveyardDropZoneRef, exileDropZoneRef }

// BattlefieldArea.svelte
{ battlefieldNonlands, battlefieldLands, commandCards, isCommanderGame,
  isDragging, isOverValidDrop, dropZone, hoveredCardId,
  onCardClick, onCardMouseDown, onCardContextMenu, onCommandCardMouseDown, onCardHover,
  battlefieldDropZoneRef, commandDropZoneRef }

// OpponentSection.svelte (note: includes playerId in callbacks)
{ opponent, otherPlayers, battlefieldNonlands, battlefieldLands, commandCards,
  isCommanderGame, showLifeMenu,
  onSelectOpponent, onLifeChange(delta, playerId), onPoisonChange(delta, playerId),
  onToggleLifeMenu, onCardContextMenu }

// DeckContextMenu.svelte
{ position, deckCount, playerName, onClose, actions: MenuAction[] }
```

### Anti-Patterns to Avoid

❌ **DO NOT** create new drag-drop infrastructure - use existing `drag-drop.ts`
❌ **DO NOT** create new zone utilities - extend `zones.ts` or `playtest-helpers.ts`
❌ **DO NOT** duplicate store subscriptions - use derived stores
❌ **DO NOT** create inline WebSocket handlers - use store methods
❌ **DO NOT** hardcode zone names - use constants from `zones.ts`

### Documentation Sources

All patterns verified from:
- `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/playtest/+page.svelte` (lines 1-1992)
- `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/game/[id]/+page.svelte` (lines 1-3554)
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/utils/drag-drop.ts` (lines 1-358)
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/playtest-game.ts` (lines 1-1257)
- Component files in `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/`

---

## Phase 1: Extract Keyboard Shortcuts System

**Goal:** Replace ~100 lines of duplicated keyboard event handling with config-driven system.

### Documentation References

**Current Pattern (Playtest page lines 955-1057, Game page lines 1072-1198):**
```typescript
function handleGlobalKeydown(event: KeyboardEvent): void {
  if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
    return;
  }

  const key = event.key.toLowerCase();

  switch (key) {
    case '?': showKeyboardShortcuts = !showKeyboardShortcuts; break;
    case 'x': handleUntapAll(); break;
    case 'c': handleDrawCard(); break;
    case 'v': handleShuffleLibrary(); break;
    // ... 8-12 more cases
  }

  // Context-sensitive shortcuts with hoveredCard
  if (hoveredCardId) {
    switch (key) {
      case 'd': moveCardToZone(cardId, 'GRAVEYARD'); break;
      // ... 3-4 more cases
    }
  }
}
```

### Implementation Steps

1. **Create utility file** `$lib/utils/keyboard-shortcuts.ts`:

```typescript
export interface KeyboardShortcut {
  key: string;
  handler: () => void | Promise<void>;
  requiresHoveredCard?: boolean;
  description: string;
  category?: 'general' | 'card' | 'game' | 'ui';
}

export function createKeyboardHandler(
  shortcuts: KeyboardShortcut[],
  getHoveredCard?: () => string | null
): (event: KeyboardEvent) => void {
  return (event: KeyboardEvent) => {
    // Ignore if typing in input
    if (
      event.target instanceof HTMLInputElement ||
      event.target instanceof HTMLTextAreaElement
    ) {
      return;
    }

    const key = event.key.toLowerCase();
    const hoveredCard = getHoveredCard?.();

    for (const shortcut of shortcuts) {
      if (shortcut.key === key) {
        // Check if requires hovered card
        if (shortcut.requiresHoveredCard && !hoveredCard) {
          continue;
        }

        event.preventDefault();
        shortcut.handler();
        return;
      }
    }
  };
}

// Helper to group shortcuts by category for help display
export function groupShortcutsByCategory(
  shortcuts: KeyboardShortcut[]
): Record<string, KeyboardShortcut[]> {
  const grouped: Record<string, KeyboardShortcut[]> = {
    general: [],
    card: [],
    game: [],
    ui: []
  };

  for (const shortcut of shortcuts) {
    const category = shortcut.category || 'general';
    grouped[category].push(shortcut);
  }

  return grouped;
}
```

2. **Update playtest page** to use new system (replace lines 955-1057):

```typescript
import { createKeyboardHandler, type KeyboardShortcut } from '$lib/utils/keyboard-shortcuts';

// Define shortcuts array
const shortcuts: KeyboardShortcut[] = [
  // UI
  { key: '?', handler: () => showKeyboardShortcuts = !showKeyboardShortcuts, description: 'Show shortcuts', category: 'ui' },
  { key: 'escape', handler: handleEscape, description: 'Close dialogs', category: 'ui' },

  // Game actions
  { key: 'x', handler: handleUntapAll, description: 'Untap all', category: 'game' },
  { key: 'c', handler: handleDrawCard, description: 'Draw card', category: 'game' },
  { key: 'v', handler: handleShuffleLibrary, description: 'Shuffle library', category: 'game' },
  { key: 'e', handler: handleNextTurn, description: 'Next turn', category: 'game' },
  { key: 'w', handler: () => showCreateTokenDialog = true, description: 'Create token', category: 'game' },
  { key: 'f', handler: () => showDeckSearch = true, description: 'Search library', category: 'game' },

  // Card actions (require hovered card)
  {
    key: 'd',
    handler: () => {
      if (hoveredCardId) moveCardToZone(hoveredCardId, 'GRAVEYARD');
    },
    requiresHoveredCard: true,
    description: 'Move to graveyard',
    category: 'card'
  },
  {
    key: 's',
    handler: () => {
      if (hoveredCardId) moveCardToZone(hoveredCardId, 'EXILE');
    },
    requiresHoveredCard: true,
    description: 'Move to exile',
    category: 'card'
  },
  {
    key: 'r',
    handler: () => {
      if (hoveredCardId) moveCardToZone(hoveredCardId, 'HAND');
    },
    requiresHoveredCard: true,
    description: 'Return to hand',
    category: 'card'
  },
  {
    key: 't',
    handler: () => {
      if (hoveredCardId) moveCardToZone(hoveredCardId, 'LIBRARY');
    },
    requiresHoveredCard: true,
    description: 'To top of library',
    category: 'card'
  },
  {
    key: 'k',
    handler: () => {
      if (hoveredCardId) {
        const card = battlefield.find(c => c.id === hoveredCardId);
        if (card) {
          selectedCardForCounters = { id: card.id, name: card.name };
          showCounterDialog = true;
        }
      }
    },
    requiresHoveredCard: true,
    description: 'Add/remove counters',
    category: 'card'
  }
];

// Create handler
const handleGlobalKeydown = createKeyboardHandler(shortcuts, () => hoveredCardId);
```

3. **Update game page similarly** (adapt for game API differences).

4. **Export shortcuts array** for KeyboardShortcutsModal to display dynamically.

### Verification Checklist

- [ ] All shortcuts from playtest page work identically
- [ ] All shortcuts from game page work identically
- [ ] Shortcuts don't fire when typing in inputs
- [ ] Hovered card shortcuts only work with card present
- [ ] '?' key opens shortcuts help modal
- [ ] Escape key closes modals in order
- [ ] Shortcuts array can be passed to help modal for display

### Anti-Pattern Guards

- ✅ Use existing `hoveredCardId` state (don't create new tracking)
- ✅ Keep handler functions in page (only extract the switch/if structure)
- ✅ Use existing modal state variables (showKeyboardShortcuts, etc.)
- ❌ DO NOT implement new keyboard event capture - use svelte:window

**Estimated Savings:** ~80-90 lines per page

---

## Phase 2: Extract Async Action Handler

**Goal:** Replace 8-12 duplicated async/loading/error wrapper patterns with single utility.

### Documentation References

**Current Pattern (repeats 8+ times per page):**

```typescript
// Playtest page lines 461-465, 470-474, 479-483, etc.
async function handleDrawCard(): Promise<void> {
  if (!me) return;
  isActionLoading = true;
  try {
    await playtestGameStore.drawCards(me.playerId, 1);
    toast.success('Drew a card');
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Unknown error';
    console.error('Failed to draw card:', err);
    toast.error(`Failed to draw: ${errorMessage}`);
  } finally {
    isActionLoading = false;
  }
}
```

### Implementation Steps

1. **Create utility file** `$lib/utils/game-actions.ts`:

```typescript
import { toast } from '$lib/stores/toast';

export interface ActionOptions {
  loadingState?: { set: (value: boolean) => void };
  successMessage?: string;
  errorPrefix?: string;
  silent?: boolean; // Don't show toasts
  onSuccess?: () => void;
  onError?: (error: string) => void;
}

export async function executeGameAction(
  action: () => Promise<void>,
  options: ActionOptions = {}
): Promise<boolean> {
  const {
    loadingState,
    successMessage,
    errorPrefix = 'Failed',
    silent = false,
    onSuccess,
    onError
  } = options;

  if (loadingState) {
    loadingState.set(true);
  }

  try {
    await action();

    if (!silent && successMessage) {
      toast.success(successMessage);
    }

    onSuccess?.();
    return true;
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Unknown error';
    console.error(`${errorPrefix}:`, err);

    if (!silent) {
      toast.error(`${errorPrefix}: ${message}`);
    }

    onError?.(message);
    return false;
  } finally {
    if (loadingState) {
      loadingState.set(false);
    }
  }
}

// Convenience wrapper for playtest actions
export function createPlaytestAction(
  action: () => Promise<void>,
  successMessage?: string
) {
  return () => executeGameAction(action, { successMessage });
}

// Convenience wrapper for game API actions
export function createGameAction(
  action: () => Promise<void>,
  loadingState: { set: (value: boolean) => void },
  successMessage?: string
) {
  return () => executeGameAction(action, { loadingState, successMessage });
}
```

2. **Update playtest page** to use wrapper (example for handleDrawCard):

```typescript
import { executeGameAction } from '$lib/utils/game-actions';

async function handleDrawCard(): Promise<void> {
  if (!me) return;

  await executeGameAction(
    () => playtestGameStore.drawCards(me.playerId, 1),
    { successMessage: 'Drew a card' }
  );
}

async function handleShuffleLibrary(): Promise<void> {
  if (!me) return;

  await executeGameAction(
    () => playtestGameStore.shuffleLibrary(me.playerId),
    { successMessage: 'Shuffled library' }
  );
}

async function handleUntapAll(): Promise<void> {
  if (!me) return;

  await executeGameAction(
    () => playtestGameStore.untapAll(me.playerId),
    { successMessage: 'Untapped all permanents' }
  );
}

// Can also use inline for simple cases
const handleDrawN = (count: number) =>
  executeGameAction(
    () => {
      if (!me) return Promise.reject('No player');
      return playtestGameStore.drawCards(me.playerId, count);
    },
    { successMessage: `Drew ${count} card(s)` }
  );
```

3. **Update game page similarly** with loading state:

```typescript
let isActionLoading = $state(false);

async function handleDrawCard(): Promise<void> {
  const me = $localPlayer;
  if (!me || !gameId) return;

  await executeGameAction(
    () => drawCards(gameId, me.playerId, 1),
    {
      loadingState: { set: (v) => isActionLoading = v },
      successMessage: 'Drew a card'
    }
  );
}
```

### Verification Checklist

- [ ] All async actions show loading state correctly
- [ ] Success toasts appear for all actions
- [ ] Error toasts appear on failures
- [ ] Console errors still logged
- [ ] Loading state resets on both success and error
- [ ] Silent mode works (no toasts when needed)

### Anti-Pattern Guards

- ✅ Use existing toast store (don't create new notification system)
- ✅ Keep validation logic in handler (me check, gameId check)
- ✅ Use existing loading state variables
- ❌ DO NOT swallow errors silently - always log to console

**Estimated Savings:** ~40-60 lines per page

---

## Phase 3: Extract Drop Zone Registration Pattern

**Goal:** Replace 5-6 repeated $effect blocks with single registration helper.

### Documentation References

**Current Pattern (repeats 5-6 times per page):**

```typescript
// Playtest page lines 1062-1090 (battlefield), 1092-1120 (graveyard), etc.
let battlefieldDropZoneEl: HTMLElement | null = $state(null);
let dropZoneUnregister: (() => void) | null = null;

$effect(() => {
  if (battlefieldDropZoneEl && !dropZoneUnregister) {
    dropZoneUnregister = dragDropStore.registerDropZone({
      id: 'battlefield',
      type: 'battlefield',
      element: battlefieldDropZoneEl,
      accepts: (cardId, sourceZone) => sourceZone !== 'battlefield',
      onDrop: (cardId) => {
        const dragState = get(dragDropStore);
        if (dragState.sourceZone) {
          playtestGameStore.moveCardToZone(cardId, 'BATTLEFIELD');
        }
      }
    });
  }

  return () => {
    if (dropZoneUnregister) {
      dropZoneUnregister();
      dropZoneUnregister = null;
    }
  };
});
```

### Implementation Steps

1. **Create utility file** `$lib/utils/drop-zone-helpers.ts`:

```typescript
import { dragDropStore, type DropZone, type SourceZone } from './drag-drop';
import type { DropZoneConfig } from './drag-drop';

export interface DropZoneSetup {
  zoneId: string;
  zoneType: DropZone;
  accepts: (cardId: string, sourceZone: SourceZone) => boolean;
  onDrop: (cardId: string, sourceZone: SourceZone) => void;
}

/**
 * Helper to setup drop zone registration with automatic cleanup.
 * Use in $effect with element ref.
 */
export function registerDropZone(
  element: HTMLElement | null,
  unregisterRef: { current: (() => void) | null },
  setup: DropZoneSetup
): () => void {
  if (element && !unregisterRef.current) {
    unregisterRef.current = dragDropStore.registerDropZone({
      id: setup.zoneId,
      type: setup.zoneType,
      element,
      accepts: setup.accepts,
      onDrop: setup.onDrop
    });
  }

  return () => {
    if (unregisterRef.current) {
      unregisterRef.current();
      unregisterRef.current = null;
    }
  };
}

/**
 * Batch register multiple drop zones with single effect.
 * Returns cleanup function.
 */
export function registerMultipleDropZones(
  zones: Array<{
    element: HTMLElement | null;
    setup: DropZoneSetup;
  }>
): () => void {
  const unregisterFns: Array<(() => void) | null> = [];

  for (const { element, setup } of zones) {
    if (element) {
      unregisterFns.push(
        dragDropStore.registerDropZone({
          id: setup.zoneId,
          type: setup.zoneType,
          element,
          accepts: setup.accepts,
          onDrop: setup.onDrop
        })
      );
    }
  }

  return () => {
    for (const unregister of unregisterFns) {
      unregister?.();
    }
  };
}
```

2. **Update playtest page** to use helper:

```typescript
import { registerDropZone } from '$lib/utils/drop-zone-helpers';

// Element refs (keep these)
let battlefieldDropZoneEl: HTMLElement | null = $state(null);
let graveyardDropZoneEl: HTMLElement | null = $state(null);
let exileDropZoneEl: HTMLElement | null = $state(null);
// ... etc

// Unregister refs
const battlefieldUnregister = { current: null as (() => void) | null };
const graveyardUnregister = { current: null as (() => void) | null };
const exileUnregister = { current: null as (() => void) | null };
// ... etc

// Register battlefield
$effect(() =>
  registerDropZone(battlefieldDropZoneEl, battlefieldUnregister, {
    zoneId: 'battlefield',
    zoneType: 'battlefield',
    accepts: (cardId, sourceZone) => sourceZone !== 'battlefield',
    onDrop: (cardId) => {
      playtestGameStore.moveCardToZone(cardId, 'BATTLEFIELD');
    }
  })
);

// Register graveyard
$effect(() =>
  registerDropZone(graveyardDropZoneEl, graveyardUnregister, {
    zoneId: 'graveyard',
    zoneType: 'graveyard',
    accepts: (cardId, sourceZone) => sourceZone !== 'graveyard',
    onDrop: (cardId) => {
      playtestGameStore.moveCardToZone(cardId, 'GRAVEYARD');
    }
  })
);

// ... repeat for other zones
```

**Alternative:** Single effect with batch registration:

```typescript
import { registerMultipleDropZones } from '$lib/utils/drop-zone-helpers';

$effect(() => {
  const zones = [
    {
      element: battlefieldDropZoneEl,
      setup: {
        zoneId: 'battlefield',
        zoneType: 'battlefield' as const,
        accepts: (cardId, sourceZone) => sourceZone !== 'battlefield',
        onDrop: (cardId) => playtestGameStore.moveCardToZone(cardId, 'BATTLEFIELD')
      }
    },
    {
      element: graveyardDropZoneEl,
      setup: {
        zoneId: 'graveyard',
        zoneType: 'graveyard' as const,
        accepts: (cardId, sourceZone) => sourceZone !== 'graveyard',
        onDrop: (cardId) => playtestGameStore.moveCardToZone(cardId, 'GRAVEYARD')
      }
    },
    // ... other zones
  ];

  return registerMultipleDropZones(zones);
});
```

### Verification Checklist

- [ ] All drop zones accept cards correctly
- [ ] Drop zone highlights show on drag
- [ ] Cards move to correct zones on drop
- [ ] Drop zones clean up on unmount
- [ ] No memory leaks (check with multiple page navigations)
- [ ] Accepts() validation still works

### Anti-Pattern Guards

- ✅ Use existing dragDropStore.registerDropZone (don't create new registration)
- ✅ Keep element refs in component (they're needed for DOM binding)
- ✅ Keep accepts() and onDrop() logic inline (zone-specific)
- ❌ DO NOT create global drop zone registry - use per-component registration

**Estimated Savings:** ~60-80 lines per page

---

## Phase 4: Extract Clipboard Utilities

**Goal:** Deduplicate clipboard copy logic with fallback handling.

### Documentation References

**Current Pattern (Playtest page lines 770-799, 914-950):**

```typescript
async function handleCopyLog(): Promise<void> {
  const logText = playtestGameStore.buildLogText($playtestGameStore);

  try {
    if (typeof navigator !== 'undefined' && navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(logText);
      toast.success('Game log copied to clipboard!');
      return;
    }

    // Fallback for older browsers
    const textarea = document.createElement('textarea');
    textarea.value = logText;
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.select();

    const successful = document.execCommand('copy');
    document.body.removeChild(textarea);

    if (successful) {
      toast.success('Game log copied to clipboard!');
    } else {
      throw new Error('execCommand failed');
    }
  } catch (err) {
    console.error('Failed to copy log to clipboard:', err);
    toast.error('Failed to copy log');
  }
}
```

### Implementation Steps

1. **Create utility file** `$lib/utils/clipboard.ts`:

```typescript
import { toast } from '$lib/stores/toast';

export interface ClipboardOptions {
  successMessage?: string;
  errorMessage?: string;
  silent?: boolean;
}

/**
 * Copy text to clipboard with fallback for older browsers.
 * Returns true on success, false on failure.
 */
export async function copyToClipboard(
  text: string,
  options: ClipboardOptions = {}
): Promise<boolean> {
  const {
    successMessage = 'Copied to clipboard',
    errorMessage = 'Failed to copy',
    silent = false
  } = options;

  try {
    // Modern API (supported in all modern browsers)
    if (typeof navigator !== 'undefined' && navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text);

      if (!silent) {
        toast.success(successMessage);
      }

      return true;
    }

    // Fallback for older browsers (IE11, Safari < 13.1)
    const textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.style.position = 'fixed';
    textarea.style.left = '-9999px';
    textarea.style.top = '-9999px';
    textarea.style.opacity = '0';
    textarea.setAttribute('readonly', '');

    document.body.appendChild(textarea);

    // iOS requires contentEditable for selection
    if (navigator.userAgent.match(/ipad|ipod|iphone/i)) {
      textarea.contentEditable = 'true';
      textarea.readOnly = false;
      const range = document.createRange();
      range.selectNodeContents(textarea);
      const selection = window.getSelection();
      selection?.removeAllRanges();
      selection?.addRange(range);
      textarea.setSelectionRange(0, textarea.value.length);
    } else {
      textarea.select();
    }

    const successful = document.execCommand('copy');
    document.body.removeChild(textarea);

    if (successful) {
      if (!silent) {
        toast.success(successMessage);
      }
      return true;
    } else {
      throw new Error('execCommand("copy") returned false');
    }
  } catch (err) {
    console.error('Failed to copy to clipboard:', err);

    if (!silent) {
      toast.error(errorMessage);
    }

    return false;
  }
}

/**
 * Convenience function to copy JSON with formatting.
 */
export async function copyJSON(
  data: unknown,
  options?: ClipboardOptions
): Promise<boolean> {
  const text = JSON.stringify(data, null, 2);
  return copyToClipboard(text, options);
}

/**
 * Convenience function to copy array of lines.
 */
export async function copyLines(
  lines: string[],
  options?: ClipboardOptions
): Promise<boolean> {
  const text = lines.join('\n');
  return copyToClipboard(text, options);
}
```

2. **Update playtest page** to use utility:

```typescript
import { copyToClipboard } from '$lib/utils/clipboard';

async function handleCopyLog(): Promise<void> {
  const logText = playtestGameStore.buildLogText($playtestGameStore);
  await copyToClipboard(logText, {
    successMessage: 'Game log copied to clipboard!',
    errorMessage: 'Failed to copy log'
  });
}

async function handleCopyGamestate(): Promise<void> {
  const gamestate = JSON.stringify($playtestGameStore, null, 2);
  await copyToClipboard(gamestate, {
    successMessage: 'Game state copied!',
    errorMessage: 'Failed to copy game state'
  });
}
```

### Verification Checklist

- [ ] Copy works in modern browsers (Chrome, Firefox, Safari, Edge)
- [ ] Copy works in older browsers (IE11 if supported)
- [ ] Copy works on iOS Safari
- [ ] Success toast shows on successful copy
- [ ] Error toast shows on failure
- [ ] Console logs error details
- [ ] Silent mode works (no toasts)

### Anti-Pattern Guards

- ✅ Use existing toast store for notifications
- ✅ Use navigator.clipboard API as primary (fastest)
- ✅ Keep execCommand fallback (some enterprise browsers)
- ❌ DO NOT use Flash-based clipboard libraries

**Estimated Savings:** ~30-40 lines per page

---

## Phase 5: Extract Drag Threshold Handlers

**Goal:** Reduce 40+ line drag handlers to ~10 lines using reusable pattern.

### Documentation References

**Current Pattern (Playtest lines 643-689, Game lines 921-967):**

```typescript
let battlefieldDragStartPosition = $state<{ x: number; y: number } | null>(null);
let battlefieldIsDragPending = $state(false);
const DRAG_THRESHOLD = 5;

function handleBattlefieldCardMouseDown(
  cardId: string,
  cardName: string,
  event: MouseEvent
): void {
  if (event.button !== 0) return;
  event.preventDefault();
  event.stopPropagation();

  battlefieldDragStartPosition = { x: event.clientX, y: event.clientY };
  battlefieldIsDragPending = true;

  const handleMouseMove = (moveEvent: MouseEvent) => {
    if (!battlefieldDragStartPosition || !battlefieldIsDragPending) return;

    const dx = moveEvent.clientX - battlefieldDragStartPosition.x;
    const dy = moveEvent.clientY - battlefieldDragStartPosition.y;
    const distance = Math.sqrt(dx * dx + dy * dy);

    if (distance >= DRAG_THRESHOLD) {
      battlefieldIsDragPending = false;
      const validZones = getAllValidDropZones('battlefield' as SourceZone);
      dragDropStore.startDrag(
        cardId,
        cardName,
        'battlefield' as SourceZone,
        moveEvent.clientX,
        moveEvent.clientY,
        validZones
      );

      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    }
  };

  const handleMouseUp = () => {
    battlefieldIsDragPending = false;
    battlefieldDragStartPosition = null;
    document.removeEventListener('mousemove', handleMouseMove);
    document.removeEventListener('mouseup', handleMouseUp);
  };

  document.addEventListener('mousemove', handleMouseMove);
  document.addEventListener('mouseup', handleMouseUp);
}
```

### Implementation Steps

1. **Extend** `$lib/utils/drag-drop.ts` with threshold helper:

```typescript
// Add to existing drag-drop.ts file (around line 340)

export interface DragThresholdState {
  startPosition: { x: number; y: number } | null;
  isPending: boolean;
}

/**
 * Creates a mousedown handler with drag threshold to prevent
 * accidental drags from clicks.
 */
export function createThresholdDragHandler(
  sourceZone: SourceZone,
  threshold: number = 5
): {
  state: DragThresholdState;
  handleMouseDown: (
    cardId: string,
    cardName: string,
    event: MouseEvent
  ) => void;
} {
  const state: DragThresholdState = {
    startPosition: null,
    isPending: false
  };

  function handleMouseDown(
    cardId: string,
    cardName: string,
    event: MouseEvent
  ): void {
    if (event.button !== 0) return; // Only left click
    event.preventDefault();
    event.stopPropagation();

    state.startPosition = { x: event.clientX, y: event.clientY };
    state.isPending = true;

    const handleMouseMove = (moveEvent: MouseEvent) => {
      if (!state.startPosition || !state.isPending) return;

      // Calculate distance from start
      if (isDragThresholdMet(state.startPosition, {
        x: moveEvent.clientX,
        y: moveEvent.clientY
      }, threshold)) {
        state.isPending = false;

        // Start actual drag
        const validZones = getAllValidDropZones(sourceZone);
        dragDropStore.startDrag(
          cardId,
          cardName,
          sourceZone,
          moveEvent.clientX,
          moveEvent.clientY,
          validZones
        );

        cleanup();
      }
    };

    const handleMouseUp = () => {
      state.isPending = false;
      state.startPosition = null;
      cleanup();
    };

    const cleanup = () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };

    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
  }

  return { state, handleMouseDown };
}
```

2. **Update playtest page** to use helper:

```typescript
import { createThresholdDragHandler } from '$lib/utils/drag-drop';

// Create handlers for each source zone
const battlefieldDragHandler = createThresholdDragHandler('battlefield');
const commandDragHandler = createThresholdDragHandler('command');

// Use in component props
const handleBattlefieldCardMouseDown = battlefieldDragHandler.handleMouseDown;
const handleCommandCardMouseDown = commandDragHandler.handleMouseDown;
```

3. **Update game page similarly**.

### Verification Checklist

- [ ] Clicking card doesn't start drag (<5px movement)
- [ ] Dragging card >5px starts drag correctly
- [ ] Drag starts from correct source zone
- [ ] Valid drop zones calculated correctly
- [ ] Event listeners cleaned up after drag
- [ ] Works for all zones (battlefield, command, hand, graveyard, etc.)

### Anti-Pattern Guards

- ✅ Use existing `isDragThresholdMet()` utility
- ✅ Use existing `getAllValidDropZones()` utility
- ✅ Use existing `dragDropStore.startDrag()` method
- ❌ DO NOT create new drag state store - use local state

**Estimated Savings:** ~30-40 lines per handler × 2 zones = ~60-80 lines per page

---

## Phase 6: Extract Derived State Patterns

**Goal:** Create composable for battlefield card filtering (100% duplicated logic).

### Documentation References

**Current Pattern (Playtest lines 153-172, Game lines 278-307):**

```typescript
function isLandPermanent(cardType?: string | null): boolean {
  return !!cardType && /\bland\b/i.test(cardType);
}

const myBattlefield = $derived(
  $battlefield.filter((c) => c.controllerId === activeControlSeat)
);

const myBattlefieldNonlands = $derived(
  myBattlefield.filter((c) => !isLandPermanent(c.type))
);

const myBattlefieldLands = $derived(
  myBattlefield.filter((c) => isLandPermanent(c.type))
);

const opponentBattlefield = $derived.by(() => {
  const opponent = selectedOpponent;
  return opponent ? $battlefield.filter((c) => c.controllerId === opponent.playerId) : [];
});

const opponentBattlefieldNonlands = $derived.by(() =>
  opponentBattlefield.filter((c) => !isLandPermanent(c.type))
);

const opponentBattlefieldLands = $derived.by(() =>
  opponentBattlefield.filter((c) => isLandPermanent(c.type))
);
```

### Implementation Steps

1. **Create composable file** `$lib/composables/use-battlefield.svelte.ts`:

```typescript
import type { CardView } from '$lib/generated/mage/v1/models';

/**
 * Check if card type includes "Land"
 */
export function isLandPermanent(cardType?: string | null): boolean {
  return !!cardType && /\bland\b/i.test(cardType);
}

/**
 * Composable for filtering battlefield cards by controller and type.
 * Returns reactive derived values for lands and nonlands.
 */
export function useBattlefieldCards(
  battlefield: CardView[],
  controllerId: string | null
) {
  const myCards = $derived(
    controllerId ? battlefield.filter((c) => c.controllerId === controllerId) : []
  );

  const nonlands = $derived(
    myCards.filter((c) => !isLandPermanent(c.type))
  );

  const lands = $derived(
    myCards.filter((c) => isLandPermanent(c.type))
  );

  return {
    all: myCards,
    nonlands,
    lands
  };
}

/**
 * Composable for opponent battlefield with selection support.
 */
export function useOpponentBattlefield(
  battlefield: CardView[],
  selectedOpponentId: string | null
) {
  const opponentCards = $derived(
    selectedOpponentId
      ? battlefield.filter((c) => c.controllerId === selectedOpponentId)
      : []
  );

  const nonlands = $derived(
    opponentCards.filter((c) => !isLandPermanent(c.type))
  );

  const lands = $derived(
    opponentCards.filter((c) => isLandPermanent(c.type))
  );

  return {
    all: opponentCards,
    nonlands,
    lands
  };
}
```

2. **Update playtest page** to use composables:

```typescript
import { useBattlefieldCards, useOpponentBattlefield } from '$lib/composables/use-battlefield.svelte';

// Replace inline derived state with composables
const myBattlefieldCards = useBattlefieldCards($battlefield, activeControlSeat);
const opponentBattlefieldCards = useOpponentBattlefield($battlefield, selectedOpponentId);

// Access in template:
// myBattlefieldCards.nonlands
// myBattlefieldCards.lands
// opponentBattlefieldCards.nonlands
// opponentBattlefieldCards.lands
```

3. **Update game page similarly**.

### Verification Checklist

- [ ] My battlefield filters by local player ID
- [ ] Opponent battlefield filters by selected opponent
- [ ] Lands and nonlands separated correctly
- [ ] Reactivity works (updates on card moves)
- [ ] Works with empty battlefield
- [ ] Works with multiple opponents

### Anti-Pattern Guards

- ✅ Use existing CardView type from generated protobuf
- ✅ Keep $derived for reactivity (Svelte 5 pattern)
- ✅ Use regex from original implementation (case-insensitive)
- ❌ DO NOT change card type detection logic - keep exact regex

**Estimated Savings:** ~20-30 lines per page

---

## Phase 7: Extract Dialog State Management

**Goal:** Create reusable dialog state pattern to reduce boilerplate.

### Documentation References

**Current Pattern (repeated for 5+ dialogs per page):**

```typescript
// State
let showNumberInputDialog = $state(false);
let numberInputDialogConfig = $state<{
  title: string;
  defaultValue: number;
  min: number;
  max: number;
  onConfirm: (value: number) => void;
} | null>(null);

// Helper
function showNumberInput(
  title: string,
  defaultValue: number,
  onConfirm: (value: number) => void
): void {
  numberInputDialogConfig = {
    title,
    defaultValue,
    min: 1,
    max: 99,
    onConfirm: (value) => {
      onConfirm(value);
      showNumberInputDialog = false;
      numberInputDialogConfig = null;
    }
  };
  showNumberInputDialog = true;
}

// Template
{#if showNumberInputDialog && numberInputDialogConfig}
  <NumberInputDialog
    {...numberInputDialogConfig}
    onCancel={() => {
      showNumberInputDialog = false;
      numberInputDialogConfig = null;
    }}
  />
{/if}
```

### Implementation Steps

1. **Create composable** `$lib/composables/use-dialog.svelte.ts`:

```typescript
/**
 * Composable for managing dialog open/close state with config.
 */
export function useDialog<TConfig = void>() {
  let isOpen = $state(false);
  let config = $state<TConfig | null>(null);

  const open = (dialogConfig?: TConfig) => {
    config = dialogConfig ?? (null as TConfig | null);
    isOpen = true;
  };

  const close = () => {
    isOpen = false;
    config = null;
  };

  const withClose = <TArgs extends unknown[]>(
    handler: (...args: TArgs) => void
  ) => {
    return (...args: TArgs) => {
      handler(...args);
      close();
    };
  };

  return {
    get isOpen() {
      return isOpen;
    },
    get config() {
      return config;
    },
    open,
    close,
    withClose
  };
}

// Type-safe dialog configs
export interface NumberInputConfig {
  title: string;
  defaultValue: number;
  min: number;
  max: number;
  onConfirm: (value: number) => void;
}

export interface CounterDialogConfig {
  cardId: string;
  cardName: string;
}

export interface ScryConfig {
  sessionId: string;
  playerId: string;
  cards: CardView[];
}
```

2. **Update playtest page** to use composable:

```typescript
import { useDialog, type NumberInputConfig } from '$lib/composables/use-dialog.svelte';

// Replace state variables with dialog instances
const numberInputDialog = useDialog<NumberInputConfig>();
const counterDialog = useDialog<{ cardId: string; cardName: string }>();
const scryDialog = useDialog<ScrySession>();

// Helper function becomes simpler
function showNumberInput(
  title: string,
  defaultValue: number,
  onConfirm: (value: number) => void
): void {
  numberInputDialog.open({
    title,
    defaultValue,
    min: 1,
    max: 99,
    onConfirm: numberInputDialog.withClose(onConfirm)
  });
}

// Template becomes cleaner
{#if numberInputDialog.isOpen && numberInputDialog.config}
  <NumberInputDialog
    {...numberInputDialog.config}
    onCancel={numberInputDialog.close}
  />
{/if}

{#if counterDialog.isOpen && counterDialog.config}
  <CounterDialog
    cardId={counterDialog.config.cardId}
    cardName={counterDialog.config.cardName}
    onClose={counterDialog.close}
  />
{/if}
```

### Verification Checklist

- [ ] Dialogs open correctly
- [ ] Dialogs close on cancel
- [ ] Dialogs close on confirm (with withClose)
- [ ] Config resets on close
- [ ] Multiple dialogs can coexist
- [ ] Type safety maintained

### Anti-Pattern Guards

- ✅ Use Svelte 5 $state for reactivity
- ✅ Keep dialog components unchanged (only state management extracted)
- ✅ Preserve existing dialog prop interfaces
- ❌ DO NOT create global dialog manager - use per-page instances

**Estimated Savings:** ~15-20 lines per dialog × 5 dialogs = ~75-100 lines per page

---

## Phase 8: Create Shared Type Definitions

**Goal:** Extract repeated type definitions to shared module.

### Implementation Steps

1. **Create type file** `$lib/types/game-ui.ts`:

```typescript
import type { CardView } from '$lib/generated/mage/v1/models';

/**
 * Menu action for context menus (already exported by DeckContextMenu)
 */
export interface MenuAction {
  label?: string;
  icon?: string;
  divider?: boolean;
  submenu?: MenuAction[];
  onClick?: () => void;
  disabled?: boolean;
}

/**
 * Dialog configurations
 */
export interface NumberInputConfig {
  title: string;
  defaultValue: number;
  min: number;
  max: number;
  onConfirm: (value: number) => void;
}

export interface CounterDialogConfig {
  cardId: string;
  cardName: string;
}

export interface TokenCreatorConfig {
  onCreateToken: (token: {
    name: string;
    types: string[];
    power: string;
    toughness: string;
    color: string;
  }) => void;
}

/**
 * Scry session state
 */
export interface ScrySession {
  sessionId: string;
  playerId: string;
  cards: CardView[];
}

/**
 * Card selection for dialogs
 */
export interface SelectedCard {
  id: string;
  name: string;
}

/**
 * Drop zone reference wrapper
 */
export interface DropZoneRefs {
  battlefield: HTMLElement | null;
  graveyard: HTMLElement | null;
  exile: HTMLElement | null;
  hand: HTMLElement | null;
  library: HTMLElement | null;
  command: HTMLElement | null;
  stack: HTMLElement | null;
}
```

2. **Update both pages** to import shared types:

```typescript
import type {
  MenuAction,
  NumberInputConfig,
  CounterDialogConfig,
  ScrySession,
  SelectedCard
} from '$lib/types/game-ui';

// Remove local type definitions
// Use imported types instead
```

### Verification Checklist

- [ ] All type imports resolve correctly
- [ ] Type checking passes (npm run check)
- [ ] No duplicate type definitions
- [ ] IDE autocomplete works for all types

### Anti-Pattern Guards

- ✅ Import CardView from generated protobuf types
- ✅ Keep MenuAction compatible with DeckContextMenu export
- ❌ DO NOT duplicate types from generated code

**Estimated Savings:** ~20-30 lines per page

---

## Phase 9: Final Cleanup & Verification

**Goal:** Verify all extractions work together and measure final line reduction.

### Verification Steps

1. **Run type checking:**
```bash
cd mage-client-web
npm run check
```

2. **Test playtest page:**
- [ ] Start new playtest game
- [ ] Test all keyboard shortcuts
- [ ] Test drag-drop between zones
- [ ] Test deck context menu
- [ ] Test all dialogs (number input, counter, scry, token)
- [ ] Test copy log and copy gamestate
- [ ] Test battlefield filtering (lands/nonlands)

3. **Test game page:**
- [ ] Join online game
- [ ] Test all keyboard shortcuts
- [ ] Test drag-drop between zones
- [ ] Test deck context menu
- [ ] Test all dialogs
- [ ] Test battlefield filtering

4. **Measure line reduction:**
```bash
# Before (from Phase 0)
# Playtest: 1,992 lines
# Game: 3,554 lines

# Count lines after extraction
wc -l src/routes/\(protected\)/playtest/+page.svelte
wc -l src/routes/\(protected\)/game/\[id\]/+page.svelte

# Count new utility lines
wc -l src/lib/utils/keyboard-shortcuts.ts
wc -l src/lib/utils/game-actions.ts
wc -l src/lib/utils/drop-zone-helpers.ts
wc -l src/lib/utils/clipboard.ts
wc -l src/lib/composables/use-battlefield.svelte.ts
wc -l src/lib/composables/use-dialog.svelte.ts
wc -l src/lib/types/game-ui.ts
```

5. **Anti-pattern verification:**
```bash
# Check for playtest store usage in game page
grep -n "playtestGameStore\." src/routes/\(protected\)/game/\[id\]/+page.svelte

# Check for duplicated keyboard handlers
grep -n "handleGlobalKeydown" src/routes/\(protected\)/playtest/+page.svelte
grep -n "handleGlobalKeydown" src/routes/\(protected\)/game/\[id\]/+page.svelte

# Check for duplicated clipboard code
grep -n "execCommand('copy')" src/routes/\(protected\)/playtest/+page.svelte

# Verify all imports resolve
npm run build
```

### Success Criteria

- [ ] Both pages compile without errors
- [ ] Type checking passes (npm run check)
- [ ] All features work identically to before
- [ ] Playtest page reduced by 300-400 lines (20-25%)
- [ ] Game page reduced by 300-400 lines (8-12%)
- [ ] New utility files total ~800 lines
- [ ] No anti-patterns detected
- [ ] No console errors in browser
- [ ] All tests pass (if tests exist)

### Rollback Plan

Each phase is independently revertable:

- **Phase 8:** Delete `$lib/types/game-ui.ts`, restore inline types
- **Phase 7:** Remove useDialog composable, restore inline state
- **Phase 6:** Remove use-battlefield composable, restore inline derived state
- **Phase 5:** Remove threshold drag helpers, restore inline handlers
- **Phase 4:** Remove clipboard utility, restore inline copy functions
- **Phase 3:** Remove drop-zone-helpers, restore inline $effect blocks
- **Phase 2:** Remove game-actions utility, restore inline try/catch
- **Phase 1:** Remove keyboard-shortcuts utility, restore inline switch statement

---

## Summary

**Total Phases:** 9 (0 = Discovery, 1-8 = Implementation, 9 = Verification)

**Extraction Targets:**
1. Keyboard shortcuts system (~80-90 lines/page)
2. Async action handler (~40-60 lines/page)
3. Drop zone registration (~60-80 lines/page)
4. Clipboard utilities (~30-40 lines/page)
5. Drag threshold handlers (~60-80 lines/page)
6. Battlefield filtering (~20-30 lines/page)
7. Dialog state management (~75-100 lines/page)
8. Shared type definitions (~20-30 lines/page)

**Total Estimated Savings:** 385-510 lines per page (19-26% reduction for playtest, 11-14% for game)

**New Utility Files:**
- `$lib/utils/keyboard-shortcuts.ts` (~80 lines)
- `$lib/utils/game-actions.ts` (~60 lines)
- `$lib/utils/drop-zone-helpers.ts` (~70 lines)
- `$lib/utils/clipboard.ts` (~80 lines)
- Extension to `$lib/utils/drag-drop.ts` (~40 lines)
- `$lib/composables/use-battlefield.svelte.ts` (~60 lines)
- `$lib/composables/use-dialog.svelte.ts` (~40 lines)
- `$lib/types/game-ui.ts` (~50 lines)

**Total New Code:** ~480 lines (reusable across both pages and future features)

**Net Reduction:** ~770-1,020 lines (playtest) + ~770-1,020 lines (game) - ~480 new = **1,060-1,560 lines eliminated**

**Key Success Factors:**
1. Each phase is independently testable
2. All APIs verified from existing code
3. No breaking changes to component interfaces
4. Type safety maintained throughout
5. Existing functionality preserved exactly
