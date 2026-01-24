# Extract Shared Game Page Components

**Status:** Ready for Implementation
**Priority:** Medium
**Complexity:** High
**Estimated Phases:** 5

---

## Overview

Extract shared layout and component patterns from the multiplayer game page (`game/[id]/+page.svelte`) and playtest page (`playtest/+page.svelte`) to eliminate ~900 lines of duplication and create reusable game UI components.

**Current State:**
- Both pages are ~1500 lines with 60%+ identical code
- Duplicated: dialog rendering, menu overlay, drag ghost, keyboard shortcuts, game layout
- Different stores but identical UI patterns

**Target State:**
- Both pages reduced to ~400-600 lines
- Shared components in `$lib/components/game/`
- Clear separation between game logic (stores) and UI (components)

---

## Phase 0: Documentation Discovery ✅

### Findings Summary

**Svelte 5 Component Patterns Documented:**
1. **Props Pattern**: `$props()` with TypeScript interface (see Card.svelte:1-1372)
2. **State Management**: `$state()` for local UI, stores for game state
3. **Derived Values**: `$derived()` and `$derived.by()` for computations
4. **Effects**: `$effect()` for side effects and cleanup
5. **Ref Callbacks**: Used for drop zone registration (BattlefieldArea.svelte)

**Reference Components:**
- `Card.svelte` - Foundation for card display with rune-based state
- `PlayerHand.svelte` - Parent-child composition, drag detection
- `BattlefieldArea.svelte` - Drop zone ref pattern
- `Graveyard.svelte` - Modal + derived filtering pattern
- `OpponentSection.svelte` - Callback-heavy interface

**Key Anti-Patterns to Avoid:**
- Direct store mutations in components (use callbacks)
- Mixing component state with store state
- Using `onMount`/`onDestroy` (use `$effect` instead)
- Not cleaning up event listeners
- Native drag API (use custom threshold-based drag)

**Store Abstraction Pattern:**
```typescript
// Controller receives store interface, not concrete store
interface GameStoreAdapter {
  gameStore: typeof multiplayerGameStore | typeof playtestGameStore;
  getState: () => GameState;
  getLocalPlayer: () => Player | null;
  getPlayers: () => Player[];
  getBattlefield: () => CardView[];
}
```

**Documentation Sources:**
- Component examples: `mage-client-web/src/lib/components/game/`
- Store patterns: `mage-client-web/src/lib/stores/playtest-game.ts`
- UI state: `mage-client-web/src/lib/stores/game-ui-state.svelte.ts`
- Drag system: `mage-client-web/src/lib/utils/drag-drop.ts`

---

## Phase 1: Extract Pure UI Components (Dialog Collection)

**Goal:** Extract dialog components that are 95%+ identical between pages with no game logic dependencies.

### Tasks

#### 1.1 Create GameDialogs Component

Create `$lib/components/game/GameDialogs.svelte` to centralize all dialog rendering.

**Copy from:** Both pages use identical dialog structure (multiplayer:947-1067, playtest:1311-1416)

**Component Interface:**
```typescript
interface GameDialogsProps {
  // UI State
  uiState: GameUIState;

  // Game data
  gameId: string;
  me: Player | null;
  selectedCardForCountersData: CardView | null;
  deckContextMenuActions: MenuAction[];

  // Event handlers
  onCreateToken: (name: string, types: string, power: string, toughness: string, color: string) => void;
  onAddCounter: (cardId: string, counterName: string, amount: number) => void;
  onRemoveCounter: (cardId: string, counterName: string, amount: number) => void;
  onSetCounter: (cardId: string, counterName: string, amount: number) => void;
  onScryComplete: (keepOnTop: CardView[], putToBottom: CardView[]) => void;
  onNumberConfirm: (value: number) => void;

  // Variants
  keyboardShortcutsMode: 'game' | 'playtest';
  librarySearchVariant: 'multiplayer' | 'playtest';

  // Multiplayer-specific (optional)
  showGameChatOverlay?: boolean;

  // Playtest-specific (optional)
  showDebugOverlay?: boolean;
  onCopyGamestate?: () => void;
}
```

**Dialogs to Include:**
1. `TokenCreator` (conditional: `uiState.showTokenCreator`)
2. `CreateTokenDialog` (conditional: `uiState.showCreateTokenDialog`)
3. `CounterDialog` (conditional: `uiState.showCounterDialog`)
4. `DeckContextMenu` (conditional: `uiState.showDeckContextMenu`)
5. `NumberInputDialog` (conditional: `uiState.showNumberInputDialog`)
6. `ScryDialog` (conditional: `uiState.showScryDialog`)
7. `RevealTopDialog` (conditional: `uiState.showRevealTopDialog`)
8. `KeyboardShortcutsModal` (always rendered with bind)

**Library Search Handling:**
```svelte
{#if librarySearchVariant === 'multiplayer' && uiState.showDeckSearch && me}
  <LibrarySearch {gameId} {librarySearchData} ... />
{:else if librarySearchVariant === 'playtest' && uiState.showDeckSearch && me}
  <PlaytestLibrarySearch cards={me.library} ... />
{/if}
```

**Optional Components:**
```svelte
{#if showGameChatOverlay}
  <GameChatOverlay {gameId} />
{/if}

{#if showDebugOverlay && uiState.showDebugOverlay}
  <DebugOverlay {onCopyGamestate} ... />
{/if}
```

**Reference Implementations:**
- Multiplayer dialogs: `game/[id]/+page.svelte:947-1067`
- Playtest dialogs: `playtest/+page.svelte:1311-1416`

#### 1.2 Update Game Pages to Use GameDialogs

**Multiplayer page:**
```svelte
<!-- Replace lines 947-1069 with: -->
<GameDialogs
  {uiState}
  gameId={data.gameId}
  {me}
  {selectedCardForCountersData}
  {deckContextMenuActions}
  onCreateToken={(name, types, power, toughness, color) => {
    multiplayerGameStore.createToken(name, types, power, toughness, color);
    uiState.showCreateTokenDialog = false;
  }}
  onAddCounter={(cardId, name, amount) => multiplayerGameStore.addCounter(cardId, name, amount)}
  onRemoveCounter={(cardId, name, amount) => multiplayerGameStore.removeCounter(cardId, name, amount)}
  onSetCounter={(cardId, name, amount) => multiplayerGameStore.setCounter(cardId, name, amount)}
  {onScryComplete}
  onNumberConfirm={(value) => {
    uiState.numberInputDialogConfig?.onConfirm(value);
    uiState.showNumberInputDialog = false;
  }}
  keyboardShortcutsMode="game"
  librarySearchVariant="multiplayer"
  showGameChatOverlay={true}
/>
```

**Playtest page:**
```svelte
<!-- Replace lines 1311-1416 with: -->
<GameDialogs
  {uiState}
  gameId="playtest"
  {me}
  {selectedCardForCountersData}
  {deckContextMenuActions}
  onCreateToken={(name, types, power, toughness, color) => {
    playtestGameStore.createToken(name, types, power, toughness, color);
    uiState.showCreateTokenDialog = false;
  }}
  onAddCounter={(cardId, name, amount) => playtestGameStore.addCounter(cardId, name, amount)}
  onRemoveCounter={(cardId, name, amount) => playtestGameStore.removeCounter(cardId, name, amount)}
  onSetCounter={(cardId, name, amount) => playtestGameStore.setCounter(cardId, name, amount)}
  onScryComplete={handleScryComplete}
  onNumberConfirm={(value) => {
    uiState.numberInputDialogConfig?.onConfirm(value);
    uiState.showNumberInputDialog = false;
  }}
  keyboardShortcutsMode="playtest"
  librarySearchVariant="playtest"
  showDebugOverlay={true}
  onCopyGamestate={copyGamestateToClipboard}
/>
```

### Verification

**Success Criteria:**
1. Both pages compile without errors
2. All dialogs function identically to before
3. Line count reduced by ~100 lines per page
4. No runtime errors in dialog interactions

**Tests:**
```bash
# Verify no TypeScript errors
npm run check

# Verify no ESLint errors
npm run lint

# Manual UI testing checklist:
# - Open token creator dialog ✓
# - Add/remove counters on card ✓
# - Scry cards ✓
# - Search library ✓
# - Open deck context menu ✓
# - Use keyboard shortcuts modal ✓
# - (Multiplayer) Open chat overlay ✓
# - (Playtest) Open debug overlay ✓
```

**Grep Checks:**
```bash
# Verify all dialog imports removed from pages
grep -n "import.*CreateTokenDialog" src/routes/(protected)/game/\[id\]/+page.svelte
grep -n "import.*CounterDialog" src/routes/(protected)/playtest/+page.svelte
# Should return: no matches

# Verify GameDialogs imported
grep -n "import.*GameDialogs" src/routes/(protected)/game/\[id\]/+page.svelte
# Should return: import line
```

---

## Phase 2: Extract Drag Ghost Component

**Goal:** Create shared drag visualization component.

### Tasks

#### 2.1 Create DragGhost Component

Create `$lib/components/game/DragGhost.svelte`.

**Copy from:** Both pages have identical structure (multiplayer:1071-1083, playtest:1565-1577)

**Component Interface:**
```typescript
interface DragGhostProps {
  isDragging: boolean;
  cardName: string | null;
  position: { x: number; y: number };
  isOverValidDrop: boolean;
  imageSize?: 'small' | 'normal'; // Default: 'normal'
}
```

**Implementation:**
```svelte
<script lang="ts">
  import { getScryfallImageUrl } from '$lib/utils/scryfall';

  let {
    isDragging,
    cardName,
    position,
    isOverValidDrop,
    imageSize = 'normal'
  }: {
    isDragging: boolean;
    cardName: string | null;
    position: { x: number; y: number };
    isOverValidDrop: boolean;
    imageSize?: 'small' | 'normal';
  } = $props();

  const imageUrl = $derived(
    cardName ? getScryfallImageUrl(cardName, imageSize) : null
  );
</script>

{#if isDragging && cardName}
  <div class="drag-ghost" style="left: {position.x}px; top: {position.y}px;">
    <div class="drag-ghost-card" class:valid={isOverValidDrop}>
      {#if imageUrl}
        <img src={imageUrl} alt={cardName} class="drag-ghost-image" draggable="false" />
      {:else}
        <span class="drag-ghost-name">{cardName}</span>
      {/if}
    </div>
  </div>
{/if}

<style>
  /* Copy styles from game/[id]/+page.svelte:1353-1384 */
  .drag-ghost { ... }
  .drag-ghost-card { ... }
  .drag-ghost-card.valid { ... }
  .drag-ghost-image { ... }
  .drag-ghost-name { ... }
</style>
```

**Reference Implementation:**
- Multiplayer: `game/[id]/+page.svelte:1071-1083` (styles: 1353-1384)
- Playtest: `playtest/+page.svelte:1565-1577`

#### 2.2 Update Game Pages to Use DragGhost

**Multiplayer page:**
```svelte
<!-- Replace lines 1071-1083 with: -->
<DragGhost
  {isDragging}
  cardName={dragCardName}
  position={dragPos}
  isOverValidDrop={isOverValidDrop}
  imageSize="normal"
/>
```

**Playtest page:**
```svelte
<!-- Replace lines 1565-1577 with: -->
<DragGhost
  {isDragging}
  cardName={dragCardName}
  position={dragPos}
  isOverValidDrop={isOverValidDrop}
  imageSize="small"
/>
```

### Verification

```bash
# Check compilation
npm run check

# Verify drag ghost renders during drag
# Manual test: Drag card from hand, verify ghost appears

# Verify image size differences
# Multiplayer should show larger drag ghost than playtest
```

---

## Phase 3: Extract Game Menu Component

**Goal:** Create shared slide-in menu overlay with configurable sections.

### Tasks

#### 3.1 Create GameMenu Component

Create `$lib/components/game/GameMenu.svelte`.

**Copy from:** Both pages have identical menu structure (multiplayer:750-797, playtest:1017-1135)

**Component Interface:**
```typescript
interface GameMenuProps {
  // State
  isOpen: boolean;

  // Data
  isMultiplayer: boolean;
  players?: Player[];
  activeControlSeat?: string;
  turnNumber?: number;
  activePlayerName?: string;
  availableSessions?: number;

  // UI State (for playtest controls)
  showAllHands?: boolean;

  // Event handlers
  onClose: () => void;
  onBackToLobby: () => void;
  onShowKeyboardShortcuts: () => void;

  // Playtest-specific handlers (optional)
  onSwitchPlayer?: (playerId: string) => void;
  onToggleAllHands?: () => void;
  onNextTurn?: () => void;
  onShowDebug?: () => void;
  onSessionsClick?: () => void;
}
```

**Implementation Structure:**
```svelte
<script lang="ts">
  import X from '@lucide/svelte/icons/x';
  import Keyboard from '@lucide/svelte/icons/keyboard';
  import Clock from '@lucide/svelte/icons/clock';
  import Eye from '@lucide/svelte/icons/eye';
  import EyeOff from '@lucide/svelte/icons/eye-off';

  let { isOpen, isMultiplayer, onClose, ... } = $props();
</script>

{#if isOpen}
  <!-- Backdrop (identical in both) -->
  <div class="menu-backdrop" onclick={onClose} ... />

  <!-- Menu Panel (identical structure) -->
  <div class="menu-overlay open">
    <div class="menu-header">
      <h2>Menu</h2>
      <button class="menu-close-btn" onclick={onClose}>
        <X size={24} />
      </button>
    </div>

    <div class="menu-content">
      <!-- Playtest-specific: Controls Section -->
      {#if !isMultiplayer && players && activeControlSeat}
        <div class="menu-section">
          <h3 class="menu-section-title">Controls</h3>
          <div class="menu-section-content">
            <!-- Player select, all hands toggle -->
          </div>
        </div>
      {/if}

      <!-- Playtest-specific: Turn Info Section -->
      {#if !isMultiplayer}
        <div class="menu-section">
          <h3 class="menu-section-title">Turn Info</h3>
          <!-- Turn display, next turn button -->
        </div>
      {/if}

      <!-- Common: Utilities Section -->
      <div class="menu-section">
        <h3 class="menu-section-title">Utilities</h3>
        <div class="menu-section-content">
          <button class="menu-btn" onclick={...}>
            <Keyboard size={18} /> Keyboard Shortcuts
          </button>
          {#if !isMultiplayer && onShowDebug}
            <button class="menu-btn" onclick={onShowDebug}>
              🔧 Debug View
            </button>
          {/if}
        </div>
      </div>

      <!-- Common: Navigation Section -->
      <div class="menu-section">
        <h3 class="menu-section-title">Navigation</h3>
        <div class="menu-section-content">
          <button class="menu-btn" onclick={onBackToLobby}>
            ← Back to Lobby
          </button>
          {#if !isMultiplayer && availableSessions && availableSessions > 0 && onSessionsClick}
            <button class="menu-btn" onclick={onSessionsClick}>
              <Clock size={18} /> Sessions
            </button>
          {/if}
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  /* Copy styles from game/[id]/+page.svelte:1149-1235 */
  .menu-backdrop { ... }
  .menu-overlay { ... }
  .menu-header { ... }
  .menu-section { ... }
  .menu-btn { ... }
  /* ... */
</style>
```

**Reference Implementations:**
- Multiplayer: `game/[id]/+page.svelte:750-797` (styles: 1149-1235)
- Playtest: `playtest/+page.svelte:1017-1135`

#### 3.2 Update Game Pages to Use GameMenu

**Multiplayer page:**
```svelte
<!-- Replace lines 750-797 with: -->
<GameMenu
  isOpen={uiState.showMenu}
  isMultiplayer={true}
  onClose={() => (uiState.showMenu = false)}
  onBackToLobby={() => goto('/lobby')}
  onShowKeyboardShortcuts={() => {
    uiState.showKeyboardShortcuts = true;
    uiState.showMenu = false;
  }}
/>
```

**Playtest page:**
```svelte
<!-- Replace lines 1017-1135 with: -->
<GameMenu
  isOpen={uiState.showMenu}
  isMultiplayer={false}
  {players}
  {activeControlSeat}
  {turnNumber}
  {activePlayerName}
  availableSessions={availableSessions.length}
  showAllHands={uiState.showAllHands}
  onClose={() => (uiState.showMenu = false)}
  onBackToLobby={() => goto('/lobby')}
  onShowKeyboardShortcuts={() => {
    uiState.showKeyboardShortcuts = true;
    uiState.showMenu = false;
  }}
  onSwitchPlayer={(playerId) => playtestGameStore.switchControlSeat(playerId)}
  onToggleAllHands={() => (uiState.showAllHands = !uiState.showAllHands)}
  onNextTurn={handleNextTurn}
  onShowDebug={() => {
    uiState.showDebugOverlay = true;
    uiState.showMenu = false;
  }}
  onSessionsClick={() => {
    showSessionPicker = true;
    uiState.showMenu = false;
  }}
/>
```

### Verification

```bash
# Check compilation
npm run check

# Manual tests:
# - Press 'm' to open menu ✓
# - Click backdrop to close ✓
# - Click X button to close ✓
# - (Multiplayer) Verify Utilities and Navigation sections only ✓
# - (Playtest) Verify all 4 sections present ✓
# - (Playtest) Switch controlling player ✓
# - (Playtest) Toggle all hands ✓
# - (Playtest) Click next turn ✓
```

---

## Phase 4: Extract Game Shell Component

**Goal:** Create outer container with loading/error/init state handling.

### Tasks

#### 4.1 Create GameShell Component

Create `$lib/components/game/GameShell.svelte`.

**Copy from:** Both pages have identical loading/error structure (multiplayer:708-724, playtest:912-938)

**Component Interface:**
```typescript
interface GameShellProps {
  loading: boolean;
  error: string | null;
  isInitialized: boolean;

  // Event handlers
  onRetry: () => void;

  // Slots
  children?: Snippet;
}
```

**Implementation:**
```svelte
<script lang="ts">
  import type { Snippet } from 'svelte';
  import { goto } from '$app/navigation';

  let {
    loading,
    error,
    isInitialized,
    onRetry,
    children
  }: {
    loading: boolean;
    error: string | null;
    isInitialized: boolean;
    onRetry: () => void;
    children?: Snippet;
  } = $props();
</script>

<div class="game-container">
  {#if loading}
    <div class="loading-overlay">
      <div class="spinner"></div>
      <p>Loading game...</p>
    </div>
  {:else if error}
    <div class="error-overlay">
      <div class="error-icon">⚠️</div>
      <h2>Error</h2>
      <p>{error}</p>
      <button class="btn-primary" onclick={onRetry}>
        Return to Lobby
      </button>
    </div>
  {:else if !isInitialized}
    <div class="loading-overlay">
      <p>Initializing game state...</p>
    </div>
  {:else}
    {@render children?.()}
  {/if}
</div>

<style>
  /* Copy styles from game/[id]/+page.svelte:1087-1146 */
  .game-container { ... }
  .loading-overlay, .error-overlay { ... }
  .spinner { ... }
  @keyframes spin { ... }
  .btn-primary { ... }
</style>
```

**Reference Implementations:**
- Multiplayer: `game/[id]/+page.svelte:708-724` (styles: 1087-1146)
- Playtest: `playtest/+page.svelte:912-938`

#### 4.2 Update Game Pages to Use GameShell

**Multiplayer page:**
```svelte
<!-- Replace outer div and conditionals (lines 708-724, 1085) with: -->
<GameShell
  {loading}
  {error}
  {isInitialized}
  onRetry={() => goto('/lobby')}
>
  <!-- All game content (header, layout, dialogs, etc.) -->
  <PlaytestHeader ... />
  {#if uiState.showMenu}
    <GameMenu ... />
  {/if}
  <main class="game-layout">
    <!-- ... -->
  </main>
  <GameDialogs ... />
  <DragGhost ... />
</GameShell>
```

**Playtest page:**
```svelte
<!-- Replace outer div and conditionals (lines 912-938, 1572) with: -->
<GameShell
  {loading}
  {error}
  {isInitialized}
  onRetry={() => goto('/lobby')}
>
  <!-- Mulligan dialog (shown before game starts) -->
  {#if mulliganPlayerIndex !== null && !allPlayersKept}
    <MulliganDialog ... />
  {:else}
    <!-- Game content -->
    <PlaytestHeader ... />
    {#if showSessionPicker}
      <!-- Session picker modal -->
    {/if}
    {#if uiState.showMenu}
      <GameMenu ... />
    {/if}
    <!-- ... rest of game UI -->
  {/if}
</GameShell>
```

### Verification

```bash
# Check compilation
npm run check

# Manual tests:
# - Force loading state (add delay in init) ✓
# - Force error state (break gameId) ✓
# - Verify retry button redirects to lobby ✓
# - Verify normal game loads correctly ✓
```

---

## Phase 5: Extract Game Layout Component (Optional - Advanced)

**Goal:** Create shared main game area layout. This phase is optional due to complexity.

**Note:** This extraction is more complex due to:
1. Different opponent rendering logic (1v1 vs grid)
2. Playtest has integrated GameStateLog
3. Many prop requirements

**Decision Point:** Evaluate after Phase 4 whether this extraction provides sufficient value given the complexity.

### Tasks

#### 5.1 Create GameLayout Component

Create `$lib/components/game/GameLayout.svelte`.

**Component Interface:**
```typescript
interface GameLayoutProps {
  // Players and battlefield
  otherPlayers: Player[];
  selectedOpponent: Player | null;
  myBattlefieldNonlands: CardView[];
  myBattlefieldLands: CardView[];
  opponentBattlefieldNonlands: CardView[];
  opponentBattlefieldLands: CardView[];
  myCommandCards: CardView[];
  opponentCommandCards: CardView[];
  isCommanderGame: boolean;

  // My zones
  me: Player | null;
  myGraveyard: CardView[];
  exile: CardView[];
  myMana: ManaPool;

  // UI state
  uiState: GameUIState;
  isDragging: boolean;
  isOverValidDrop: boolean;
  dropZone: string | null;

  // Battlefield state for multiplayer grid
  battlefield: CardView[];
  commandCards: CardView[];

  // Event handlers (many)
  onSelectOpponent: (playerId: string) => void;
  onLifeChange: (delta: number, playerId: string) => void;
  onPoisonChange: (delta: number, playerId: string) => void;
  onCardClick: (cardId: string) => void;
  onCardMouseDown: (cardId: string, cardName: string, e: MouseEvent) => void;
  onCardContextMenu: (cardId: string, cardName: string) => void;
  onCommandCardMouseDown: (cardId: string, cardName: string, e: MouseEvent) => void;
  onCardHover: (cardId: string | null) => void;
  onSearchLibrary: () => void;
  onDeckContextMenu: (e: MouseEvent) => void;

  // Ref callbacks
  battlefieldDropZoneRef: (el: HTMLDivElement | null) => void;
  commandDropZoneRef: (el: HTMLElement | null) => void;
  libraryDropZoneRef: (el: HTMLElement | null) => void;
  graveyardDropZoneRef: (el: HTMLElement | null) => void;
  exileDropZoneRef: (el: HTMLElement | null) => void;
  handDropZoneRef: (el: HTMLElement | null) => void;
}
```

**Reference Implementations:**
- Multiplayer: `game/[id]/+page.svelte:800-945`
- Playtest: `playtest/+page.svelte:1156-1309`

#### 5.2 Evaluation Criteria

Before implementing, verify:
- [ ] Component has <30 props (currently ~40)
- [ ] Clear boundaries between layout and logic
- [ ] No conditional logic for multiplayer vs playtest
- [ ] Reduction in line count justifies complexity

**If evaluation fails:** Keep layout inline in both pages, only extract smaller pieces (opponent section, hand area).

### Verification

```bash
# If implemented:
npm run check

# Verify all opponent layouts render correctly
# Verify battlefield interaction works
# Verify hand drag/drop works
```

---

## Final Verification Checklist

After completing all phases:

### Code Quality
- [ ] No TypeScript errors (`npm run check`)
- [ ] No ESLint errors (`npm run lint`)
- [ ] No duplicate code warnings from static analysis
- [ ] All imports resolved correctly

### Functionality
- [ ] Multiplayer game loads and plays correctly
- [ ] Playtest game loads and plays correctly
- [ ] All dialogs function identically to before extraction
- [ ] Drag and drop works in both modes
- [ ] Menu overlay works in both modes
- [ ] Loading/error states display correctly
- [ ] Keyboard shortcuts work

### Code Metrics
- [ ] Multiplayer page reduced from ~1380 lines to <700 lines
- [ ] Playtest page reduced from ~1567 lines to <800 lines
- [ ] New shared components total <500 lines
- [ ] Net reduction: >600 lines of code

### Documentation
- [ ] Update component README if exists
- [ ] Document new shared components
- [ ] Update CLAUDE.md with completion

---

## Rollback Plan

If critical issues found during Phase N:

1. Revert to previous commit before Phase N
2. Document specific issue encountered
3. Re-evaluate extraction strategy
4. Consider alternative component boundaries

Each phase should be committed separately to enable granular rollback.

---

## Anti-Pattern Guards

### Do NOT:
- ❌ Mix store mutations into shared components (pass callbacks)
- ❌ Use `onMount`/`onDestroy` (use `$effect` instead)
- ❌ Create prop drilling more than 2 levels deep
- ❌ Forget cleanup in `$effect` returns
- ❌ Directly mutate props in child components

### DO:
- ✅ Use TypeScript interfaces for all props
- ✅ Keep game logic in pages/controllers, UI in components
- ✅ Provide ref callbacks for drop zones
- ✅ Use derived state for computed values
- ✅ Clean up event listeners in `$effect` cleanup

---

## Success Metrics

**Primary:**
- Both pages compile and run without errors
- All game functionality works identically to before
- Code duplication reduced by >60%

**Secondary:**
- Easier to maintain game UI (single source of truth)
- New game modes can reuse components
- Clearer separation of concerns

**Timeline:** Phases 1-4 can be completed in 1-2 sessions. Phase 5 is optional.
