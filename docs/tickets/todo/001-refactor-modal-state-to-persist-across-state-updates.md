# Refactor Modal State to Persist Across Game State Updates

## Problem

When game state updates arrive from the server (via WebSocket `GAME_UPDATE` events), modals that are currently open lose their local state. This happens because:

1. **Game state updates trigger full reactivity**: When `gameStore.setGameView()` is called, all `$derived` values in the game page recompute
2. **Conditional rendering causes remounting**: Modals are conditionally rendered (e.g., `{#if prompt && prompt.type === 'mana'}`), and when the parent component structure changes during reactive updates, Svelte may unmount and remount the modal component
3. **Local state is ephemeral**: When a component remounts, all local `$state` variables reset to their initial values

### Example Scenario

User is uploading a deck in `DeckUploadModal`:
- They've entered a deck name, selected format, and typed in a deck list
- A game state update arrives (e.g., opponent plays a card)
- The game page rerenders due to store updates
- The modal component gets remounted
- All form state is lost: `deckName = ''`, `deckList = ''`, `structuredCards = []`, etc.

## Affected Components

### Game Page Modals (High Priority)
These modals are rendered conditionally based on game state and are most likely to be affected:

1. **DeckUploadModal** (`src/lib/components/DeckUploadModal.svelte`)
   - State: `deckName`, `selectedFormat`, `deckList`, `structuredCards`, `viewMode`, `errors`, `loading`
   - Used in: `/decks` page (but could be used elsewhere)
   - Impact: **HIGH** - User loses all entered deck data

2. **TokenCreator** (`src/lib/components/game/TokenCreator.svelte`)
   - State: `name`, `types`, `power`, `toughness`, `color`, `customAbilities`
   - Used in: Game page (`showTokenCreator` state)
   - Impact: **MEDIUM** - User loses token configuration

3. **LibrarySearch** (`src/lib/components/game/LibrarySearch.svelte`)
   - State: `selectedCardId`, `searchQuery`, `filterType`, `selectedDestination`, drag state
   - Used in: Game page (prompt-based, `prompt.type === 'librarySearch'`)
   - Impact: **MEDIUM** - User loses search query and selection

4. **ManaPayment** (`src/lib/components/game/ManaPayment.svelte`)
   - State: `selectedMana`, `isLoading`, `error`
   - Used in: Game page (prompt-based, `prompt.type === 'mana'`)
   - Impact: **LOW** - Usually quick interaction, but could be annoying

5. **XManaSelector** (`src/lib/components/game/XManaSelector.svelte`)
   - State: Likely has local state for X value selection
   - Used in: Game page (prompt-based, `prompt.type === 'xmana'`)
   - Impact: **LOW** - Usually quick interaction

6. **DeclareAttackers** (`src/lib/components/game/DeclareAttackers.svelte`)
   - State: Likely has local state for selected attackers
   - Used in: Game page (combat phase)
   - Impact: **MEDIUM** - User loses attack declarations

7. **DeclareBlockers** (`src/lib/components/game/DeclareBlockers.svelte`)
   - State: Likely has local state for selected blockers
   - Used in: Game page (combat phase)
   - Impact: **MEDIUM** - User loses block declarations

8. **AssignDamage** (`src/lib/components/game/AssignDamage.svelte`)
   - State: Likely has local state for damage assignments
   - Used in: Game page (combat phase)
   - Impact: **MEDIUM** - User loses damage assignments

### Other Modals (Lower Priority)
These may also be affected but are less critical:

- **DebugOverlay** - Usually read-only, less critical
- **KeyboardShortcutsModal** - No user input state
- **ActionLogOverlay** - May have scroll position, but less critical
- **GameChatOverlay** - May have input state, but chat is usually short-lived

## Solution Approach

### Option 1: Move State to Stores (Recommended)

Create dedicated stores for each modal's state that needs to persist:

**Pros:**
- State persists across component remounts
- Can be accessed from anywhere
- Clear separation of concerns
- Easy to debug (store state is visible in devtools)

**Cons:**
- More boilerplate
- Need to manage store lifecycle (clear on close)

**Implementation:**
```typescript
// src/lib/stores/modal-state.ts
import { writable } from 'svelte/store';

// Deck upload state
export const deckUploadState = writable<{
  deckName: string;
  selectedFormat: string;
  deckList: string;
  structuredCards: CardEntry[];
  viewMode: 'text' | 'structured';
}>({
  deckName: '',
  selectedFormat: 'Standard',
  deckList: '',
  structuredCards: [],
  viewMode: 'text'
});

export function clearDeckUploadState() {
  deckUploadState.set({
    deckName: '',
    selectedFormat: 'Standard',
    deckList: '',
    structuredCards: [],
    viewMode: 'text'
  });
}

// Token creator state
export const tokenCreatorState = writable<{
  name: string;
  types: string;
  power: string;
  toughness: string;
  color: string;
  customAbilities: string;
}>({
  name: '',
  types: 'Creature',
  power: '1',
  toughness: '1',
  color: 'colorless',
  customAbilities: ''
});

export function clearTokenCreatorState() {
  tokenCreatorState.set({
    name: '',
    types: 'Creature',
    power: '1',
    toughness: '1',
    color: 'colorless',
    customAbilities: ''
  });
}

// Library search state
export const librarySearchState = writable<{
  searchQuery: string;
  filterType: string;
  selectedDestination: 'hand' | 'battlefield' | 'graveyard' | 'exile';
}>({
  searchQuery: '',
  filterType: 'all',
  selectedDestination: 'hand'
});

export function clearLibrarySearchState() {
  librarySearchState.set({
    searchQuery: '',
    filterType: 'all',
    selectedDestination: 'hand'
  });
}
```

**Usage in component:**
```svelte
<script lang="ts">
  import { deckUploadState, clearDeckUploadState } from '$lib/stores/modal-state';
  
  // Use store instead of local state
  const state = $derived($deckUploadState);
  
  function handleClose() {
    clearDeckUploadState();
    onclose?.();
  }
</script>

<input bind:value={$deckUploadState.deckName} />
```

### Option 2: Use `{#key}` Block

Use Svelte's `{#key}` block to prevent remounting when the key doesn't change:

**Pros:**
- Minimal code changes
- Component stays mounted if key is stable

**Cons:**
- Doesn't solve the root cause
- Key management can be tricky
- May not work if parent structure changes significantly

**Implementation:**
```svelte
{#key 'deck-upload-modal'}
  {#if showUploadModal}
    <DeckUploadModal
      open={showUploadModal}
      onclose={() => (showUploadModal = false)}
    />
  {/if}
{/key}
```

### Option 3: Render Modals Outside Reactive Context

Move modals to a stable location in the component tree that doesn't depend on game state:

**Pros:**
- Prevents remounting
- Clean separation

**Cons:**
- May require significant refactoring
- Modals might need to be rendered at app level

**Implementation:**
```svelte
<!-- Render modals at the end, outside conditional blocks -->
<div class="game-container">
  <!-- Game content with conditionals -->
  {#if isMulliganPhase}
    <!-- ... -->
  {:else}
    <!-- ... -->
  {/if}
</div>

<!-- Modals always rendered, controlled by open prop -->
<DeckUploadModal bind:open={showUploadModal} />
<TokenCreator bind:open={showTokenCreator} />
```

### Option 4: Hybrid Approach (Recommended for Complex Cases)

For modals with complex state (like `DeckUploadModal`), use stores. For simpler modals, use `{#key}` blocks or move outside reactive context.

## Implementation Plan

### Phase 1: High-Impact Modals (Priority)
1. **DeckUploadModal** - Move to store (highest impact)
2. **TokenCreator** - Move to store
3. **LibrarySearch** - Move to store

### Phase 2: Combat Modals
4. **DeclareAttackers** - Move to store or use `{#key}`
5. **DeclareBlockers** - Move to store or use `{#key}`
6. **AssignDamage** - Move to store or use `{#key}`

### Phase 3: Prompt-Based Modals
7. **ManaPayment** - Use `{#key}` (simple state)
8. **XManaSelector** - Use `{#key}` (simple state)

### Phase 4: Verification
9. Test all modals during active game state updates
10. Verify state persists correctly
11. Ensure state is cleared when modals close

## Testing Checklist

For each modal:
- [ ] Open modal and enter data
- [ ] Trigger game state update (e.g., opponent action)
- [ ] Verify modal state is preserved
- [ ] Complete modal action successfully
- [ ] Verify state is cleared on close
- [ ] Test with multiple rapid game updates

## Files to Modify

### New Files
- `src/lib/stores/modal-state.ts` - Modal state stores

### Modified Files
- `src/lib/components/DeckUploadModal.svelte`
- `src/lib/components/game/TokenCreator.svelte`
- `src/lib/components/game/LibrarySearch.svelte`
- `src/lib/components/game/DeclareAttackers.svelte` (if needed)
- `src/lib/components/game/DeclareBlockers.svelte` (if needed)
- `src/lib/components/game/AssignDamage.svelte` (if needed)
- `src/routes/(protected)/game/[id]/+page.svelte` - Update modal rendering

## Notes

- Consider using a single `modalState` store with namespaced keys if there are many modals
- For modals that are prompt-based (from server), the server state might be the source of truth - consider preserving only user input state
- Some modals might benefit from auto-save to localStorage as a backup
- Consider adding a visual indicator when modal state is preserved (e.g., "Draft saved" message)

## Related Issues

- Game state updates trigger unnecessary rerenders
- Modal components should be more resilient to parent updates
- Consider using Svelte 5 runes more effectively for state management
