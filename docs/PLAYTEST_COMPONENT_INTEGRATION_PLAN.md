# Playtest Component Integration Plan

**Goal:** Integrate improved UI components from playtest mode into the normal online game page.

**Approach:** Copy-based integration using existing component APIs documented below. Each phase is independently executable with verification steps.

---

## Phase 0: Documentation Discovery ✅

### Allowed APIs

All APIs documented below are verified from source code with line references.

#### Component APIs (All in `/lib/components/game/`)

**OpponentSection.svelte** - Modular opponent display
```typescript
interface Props {
  opponent: { playerId: string; name: string; life: number; poison: number; handCount: number; libraryCount: number; graveyard: CardView[] };
  otherPlayers: Opponent[];  // For multiplayer cycling dropdown
  battlefieldNonlands: CardView[];
  battlefieldLands: CardView[];
  commandCards: CardView[];
  isCommanderGame: boolean;
  showLifeMenu: boolean;
  onSelectOpponent?: (playerId: string) => void;
  onLifeChange: (delta: number, playerId: string) => void;
  onPoisonChange: (delta: number, playerId: string) => void;
  onToggleLifeMenu: () => void;
  onCardContextMenu: (cardId: string, cardName: string) => void;
}
```

**PlayerInfoRow.svelte** - Consolidated player stats bar
```typescript
interface Props {
  player: { name: string; life: number; poison: number; libraryCount: number };
  graveyard: CardView[];
  exile: CardView[];
  mana: ManaPoolData;
  showLifeMenu: boolean;
  onLifeChange: (delta: number) => void;
  onPoisonChange: (delta: number) => void;
  onToggleLifeMenu: () => void;
  onSearchLibrary: () => void;
  onDeckContextMenu: (e: MouseEvent) => void;
  libraryDropZoneRef?: (el: HTMLElement | null) => void;
  graveyardDropZoneRef?: (el: HTMLElement | null) => void;
  exileDropZoneRef?: (el: HTMLElement | null) => void;
}
```

**BattlefieldArea.svelte** - Battlefield with land/nonland separation
```typescript
interface Props {
  battlefieldNonlands: CardView[];
  battlefieldLands: CardView[];
  commandCards: CardView[];
  isCommanderGame: boolean;
  isDragging: boolean;
  isOverValidDrop: boolean;
  dropZone: string | null;
  hoveredCardId: string | null;
  onCardClick: (cardId: string) => void;
  onCardMouseDown: (cardId: string, cardName: string, e: MouseEvent) => void;
  onCardContextMenu: (cardId: string, cardName: string) => void;
  onCommandCardMouseDown: (cardId: string, cardName: string, e: MouseEvent) => void;
  onCardHover: (cardId: string | null) => void;
  battlefieldDropZoneRef?: (el: HTMLDivElement | null) => void;
  commandDropZoneRef?: (el: HTMLDivElement | null) => void;
}
```

**DeckContextMenu.svelte** - Nested context menu system
```typescript
export interface MenuAction {
  label?: string;
  icon?: string;
  divider?: boolean;
  submenu?: MenuAction[];
  onClick?: () => void;
  disabled?: boolean;
}

interface Props {
  position: { x: number; y: number };
  deckCount: number;
  playerName: string;
  onClose: () => void;
  actions: MenuAction[];
}
```

#### Game Store API (`/lib/stores/game.ts`)

**Derived Stores (Read-Only):**
```typescript
gameView: Readable<GameView | null>
players: Readable<PlayerView[]>
localPlayer: Readable<PlayerView | null>
opponents: Readable<PlayerView[]>
hasPriority: Readable<boolean>
currentPhase: Readable<string>
currentStep: Readable<string>
battlefield: Readable<CardView[]>  // Sorted by UUID ascending
myHand: Readable<CardView[]>
myGraveyard: Readable<CardView[]>
exile: Readable<CardView[]>
command: Readable<CardView[]>
myManaPool: Readable<ManaPoolView>
```

**Helper Functions:**
```typescript
getCardById(cardId: string): CardView | null  // Searches all zones
```

#### Game API Functions (`/lib/api/game.ts`)

**Direct Actions (Available in game mode):**
```typescript
// DO NOT USE these - they don't exist in game mode:
// - playtestGameStore.moveCardToZone()
// - playtestGameStore.tapCard()
// - playtestGameStore.createToken()
// - playtestGameStore.addCounter()

// USE THESE instead - they send actions to server:
sendPlayerAction(gameId: string, action: PlayerAction): Promise<void>
sendPlayerUUID(gameId: string, uuid: string): Promise<void>
passPriority(gameId: string): Promise<void>
```

#### Drag-Drop API (`/lib/utils/drag-drop.ts`)

**Store:**
```typescript
dragDropStore: {
  registerDropZone(config: {
    id: string;
    type: DropZone;
    element: HTMLElement;
    accepts: (cardId: string, sourceZone: SourceZone) => boolean;
    onDrop: (cardId: string) => void;
  }): () => void;  // Returns unregister function
}

// Derived stores:
isDragging: Readable<boolean>
draggedCardId: Readable<string | null>
draggedCardName: Readable<string | null>
dragPosition: Readable<{ x: number; y: number }>
isOverValidDropZone: Readable<boolean>
currentDropZone: Readable<DropZone>
```

**Helper:**
```typescript
isDragThresholdMet(start: {x, y}, current: {x, y}, threshold?: number): boolean
// Default threshold: 5px
```

### Anti-Patterns to Avoid

❌ **DO NOT** use playtest store methods in game mode:
- `playtestGameStore.moveCardToZone()` - use game API instead
- `playtestGameStore.tapCard()` - use `tapUntap()` API
- `playtestGameStore.createToken()` - tokens come from server
- `playtestGameStore.addCounter()` - use server actions

❌ **DO NOT** assume free zone movement:
- Playtest allows dragging any card anywhere
- Game mode requires rule-compliant actions only
- `accepts()` functions must check game state

❌ **DO NOT** call non-existent APIs:
- No `gameStore.switchControlSeat()` - online games have fixed player
- No `gameStore.scryCards()` - scry uses server prompts
- No `gameStore.millCards()` - use card abilities/effects

❌ **DO NOT** replace GameHeader entirely:
- GameHeader has phase indicator (important for rules)
- GameHeader has priority display (important for online play)
- PlaytestHeader lacks these critical features

### Documentation Sources

All APIs verified from:
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/*.svelte` (lines documented above)
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/game.ts` (lines 1-1027)
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/api/game.ts` (lines 1-468)
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/utils/drag-drop.ts` (lines 1-337)
- `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/playtest/+page.svelte` (lines 1-1793)

---

## Phase 1: Integrate PlayerInfoRow Component

**Goal:** Replace scattered player stat display with consolidated PlayerInfoRow component.

### Documentation References

**Component Location:** `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/PlayerInfoRow.svelte`

**Usage Pattern (playtest page lines 1474-1494):**
```svelte
{#if me}
  <PlayerInfoRow
    player={{ name: me.name, life: me.life, poison: me.poison, libraryCount: me.libraryCount }}
    graveyard={myGrave}
    exile={exile}
    mana={myMana}
    showLifeMenu={showLifeMenu}
    onLifeChange={handleLifeChange}
    onPoisonChange={handlePoisonChange}
    onToggleLifeMenu={() => (showLifeMenu = !showLifeMenu)}
    onSearchLibrary={() => (showDeckSearch = true)}
    onDeckContextMenu={handleDeckContextMenu}
    libraryDropZoneRef={(el) => (libraryDropZoneEl = el)}
    graveyardDropZoneRef={(el) => (graveyardDropZoneEl = el)}
    exileDropZoneRef={(el) => (exileDropZoneEl = el)}
  />
{/if}
```

### Implementation Steps

1. **Locate game page file:**
   ```bash
   find /Users/aron/dev/opensource/mage/mage-client-web/src/routes -name "+page.svelte" | grep "game/\[id\]"
   ```
   Read the file to understand current player stat layout.

2. **Add PlayerInfoRow import:**
   ```svelte
   import PlayerInfoRow from '$lib/components/game/PlayerInfoRow.svelte';
   ```

3. **Copy handler functions from playtest (lines 432-449):**
   ```typescript
   function handleLifeChange(delta: number): void {
     const playerId = me?.playerId;
     if (!playerId) return;
     // Call game API - NOT playtest store
     modifyLife(gameId, playerId, delta);
   }

   function handlePoisonChange(delta: number): void {
     const playerId = me?.playerId;
     if (!playerId) return;
     setPlayerCounter(gameId, playerId, 'poison', (me?.poison || 0) + delta);
   }
   ```

4. **Add state variables:**
   ```typescript
   let showLifeMenu = $state(false);
   let libraryDropZoneEl: HTMLElement | null = $state(null);
   let graveyardDropZoneEl: HTMLElement | null = $state(null);
   let exileDropZoneEl: HTMLElement | null = $state(null);
   ```

5. **Replace existing player stat markup with PlayerInfoRow:**
   - Find current life display, mana pool, graveyard, exile components
   - Replace entire section with PlayerInfoRow component
   - Wire up event handlers and drop zone refs

6. **Register drop zones (copy pattern from playtest lines 912-982):**
   ```typescript
   $effect(() => {
     if (libraryDropZoneEl && !libraryDropZoneUnregister) {
       libraryDropZoneUnregister = dragDropStore.registerDropZone({
         id: 'library',
         type: 'library',
         element: libraryDropZoneEl,
         accepts: (cardId, sourceZone) => sourceZone !== 'library' && hasPriority,
         onDrop: (cardId) => {
           // Use game API to move card
           moveCardToZone(gameId, cardId, 'LIBRARY');
         }
       });
     }
     return () => {
       if (libraryDropZoneUnregister) {
         libraryDropZoneUnregister();
         libraryDropZoneUnregister = null;
       }
     };
   });
   ```

### Verification Checklist

- [ ] PlayerInfoRow renders with correct player data
- [ ] Life +/-1 buttons work and update server
- [ ] Life menu opens/closes correctly
- [ ] Poison counter displays and updates when changed
- [ ] Mana pool displays correctly
- [ ] Library, graveyard, exile zones are visible and interactive
- [ ] Drop zones accept cards (test by dragging from hand)
- [ ] Right-click on library opens context menu (if implemented)

### Anti-Pattern Guards

- ✅ Use `modifyLife(gameId, playerId, delta)` API - NOT playtest store
- ✅ Use `setPlayerCounter(gameId, playerId, 'poison', value)` API - NOT playtest store
- ✅ Drop zone `accepts()` checks `hasPriority` from game store
- ✅ `onDrop` uses `moveCardToZone(gameId, cardId, zone)` API - NOT playtest store

---

## Phase 2: Integrate BattlefieldArea Component

**Goal:** Replace battlefield display with BattlefieldArea component for land/nonland separation.

### Documentation References

**Component Location:** `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/BattlefieldArea.svelte`

**Usage Pattern (playtest page lines 1452-1471):**
```svelte
<BattlefieldArea
  battlefieldNonlands={myBattlefieldNonlands}
  battlefieldLands={myBattlefieldLands}
  commandCards={myCommandCards}
  isCommanderGame={isCommanderGame}
  isDragging={isDragging}
  isOverValidDrop={isOverValidDrop}
  dropZone={dropZone}
  hoveredCardId={hoveredCardId}
  onCardClick={handleBattlefieldCardClick}
  onCardMouseDown={handleBattlefieldCardMouseDown}
  onCardContextMenu={(cardId, cardName) => {
    selectedCardForCounters = { id: cardId, name: cardName };
    showCounterDialog = true;
  }}
  onCommandCardMouseDown={handleCommandCardMouseDown}
  onCardHover={(cardId) => (hoveredCardId = cardId)}
  battlefieldDropZoneRef={(el) => (battlefieldDropZoneEl = el)}
  commandDropZoneRef={(el) => (commandDropZoneEl = el)}
/>
```

**Derived State Pattern (playtest page lines 165-184):**
```typescript
function isLandPermanent(cardType?: string | null): boolean {
  return !!cardType && /\bland\b/i.test(cardType);
}

const myBattlefield = $derived(battlefield.filter((c) => c.controllerId === localPlayerId));
const myBattlefieldNonlands = $derived(myBattlefield.filter((c) => !isLandPermanent(c.type)));
const myBattlefieldLands = $derived(myBattlefield.filter((c) => isLandPermanent(c.type)));
const myCommandCards = $derived(commandCards.filter((c) => (c.ownerId || c.controllerId) === localPlayerId));
const isCommanderGame = $derived(commandCards.length > 0);
```

**Drag Threshold Pattern (playtest page lines 636-682):**
```typescript
let battlefieldDragStartPosition = $state<{ x: number; y: number } | null>(null);
let battlefieldIsDragPending = $state(false);
const DRAG_THRESHOLD = 5;

function handleBattlefieldCardMouseDown(cardId: string, cardName: string, event: MouseEvent): void {
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
      dragDropStore.startDrag(cardId, cardName, 'battlefield' as SourceZone, moveEvent.clientX, moveEvent.clientY, validZones);

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

1. **Add BattlefieldArea import:**
   ```svelte
   import BattlefieldArea from '$lib/components/game/BattlefieldArea.svelte';
   ```

2. **Copy helper function and derived state:**
   ```typescript
   function isLandPermanent(cardType?: string | null): boolean {
     return !!cardType && /\bland\b/i.test(cardType);
   }

   const myBattlefield = $derived($battlefield.filter((c) => c.controllerId === localPlayerId));
   const myBattlefieldNonlands = $derived(myBattlefield.filter((c) => !isLandPermanent(c.type)));
   const myBattlefieldLands = $derived(myBattlefield.filter((c) => isLandPermanent(c.type)));
   const myCommandCards = $derived($command.filter((c) => (c.ownerId || c.controllerId) === localPlayerId));
   const isCommanderGame = $derived($command.length > 0);
   ```

3. **Add drag state variables:**
   ```typescript
   let battlefieldDragStartPosition = $state<{ x: number; y: number } | null>(null);
   let battlefieldIsDragPending = $state(false);
   let commandDragStartPosition = $state<{ x: number; y: number } | null>(null);
   let commandIsDragPending = $state(false);
   let hoveredCardId = $state<string | null>(null);
   const DRAG_THRESHOLD = 5;
   ```

4. **Copy drag handlers (adapt for game API):**
   ```typescript
   function handleBattlefieldCardClick(cardId: string): void {
     const card = $battlefield.find((c) => c.id === cardId);
     if (!card) return;
     // Use game API to tap/untap
     tapUntap(gameId, cardId, !card.tapped);
   }

   // Copy handleBattlefieldCardMouseDown from playtest (shown above)
   // Copy handleCommandCardMouseDown similarly
   ```

5. **Register battlefield drop zone:**
   ```typescript
   $effect(() => {
     if (battlefieldDropZoneEl && !dropZoneUnregister) {
       dropZoneUnregister = dragDropStore.registerDropZone({
         id: 'battlefield',
         type: 'battlefield',
         element: battlefieldDropZoneEl,
         accepts: (cardId, sourceZone) => {
           // Only allow playing from hand if have priority
           if (sourceZone === 'hand') return $hasPriority;
           // Don't allow dragging battlefield card onto battlefield
           if (sourceZone === 'battlefield') return false;
           // Other zones require priority
           return $hasPriority;
         },
         onDrop: (cardId) => {
           const dragState = $dragDropStore;
           if (dragState.sourceZone === 'hand') {
             // Playing from hand - use game API
             // This will be handled by cast/play logic elsewhere
             console.log('Playing card from hand:', cardId);
           } else if (dragState.sourceZone) {
             // Moving from other zone - use move API
             moveCardToZone(gameId, cardId, 'BATTLEFIELD');
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

6. **Replace existing battlefield markup:**
   - Find current battlefield card rendering
   - Replace with BattlefieldArea component
   - Wire up all event handlers and refs

### Verification Checklist

- [ ] Battlefield displays with lands on bottom, nonlands on top
- [ ] Cards tap/untap on click
- [ ] Drag threshold prevents accidental drags (must drag >5px)
- [ ] Dragging from hand to battlefield works
- [ ] Command zone displays (if Commander format)
- [ ] Card hover highlights card
- [ ] Right-click opens context menu
- [ ] Drop zones only accept cards when have priority

### Anti-Pattern Guards

- ✅ Use `tapUntap(gameId, cardId, tapped)` API - NOT playtest store
- ✅ Use `moveCardToZone(gameId, cardId, zone)` API - NOT playtest store
- ✅ Drop zone `accepts()` checks `$hasPriority` from game store
- ✅ Use drag threshold (5px) before starting drag - prevents accidental drags from clicks

---

## Phase 3: Integrate OpponentSection Component

**Goal:** Replace OpponentPanel with OpponentSection for improved multi-opponent display.

### Documentation References

**Component Location:** `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/OpponentSection.svelte`

**1v1 Usage Pattern (playtest page lines 1376-1397):**
```svelte
{#if otherPlayers.length === 1}
  {#if selectedOpponent()}
    {@const opponent = selectedOpponent()!}
    <OpponentSection
      opponent={opponent}
      otherPlayers={otherPlayers}
      battlefieldNonlands={opponentBattlefieldNonlands()}
      battlefieldLands={opponentBattlefieldLands()}
      commandCards={opponentCommandCards()}
      isCommanderGame={isCommanderGame}
      showLifeMenu={showOpponentLifeMenu}
      onSelectOpponent={(playerId) => (selectedOpponentId = playerId)}
      onLifeChange={handleLifeChange}
      onPoisonChange={handlePoisonChange}
      onToggleLifeMenu={() => (showOpponentLifeMenu = !showOpponentLifeMenu)}
      onCardContextMenu={(cardId, cardName) => {
        selectedCardForCounters = { id: cardId, name: cardName };
        showCounterDialog = true;
      }}
    />
  {/if}
{/if}
```

**Multi-Opponent Grid Layout (playtest page lines 1399-1449):**
```svelte
{:else}
  <!-- Grid layout for large screens -->
  <div class="opponents-grid opponents-grid-large">
    {#each otherPlayers as opponent (opponent.playerId)}
      {@const oppBattlefield = battlefield.filter((c) => c.controllerId === opponent.playerId)}
      {@const oppBattlefieldNonlands = oppBattlefield.filter((c) => !isLandPermanent(c.type))}
      {@const oppBattlefieldLands = oppBattlefield.filter((c) => isLandPermanent(c.type))}
      {@const oppCommandCards = commandCards.filter((c) => (c.ownerId || c.controllerId) === opponent.playerId)}
      <OpponentSection
        opponent={opponent}
        otherPlayers={[]}
        battlefieldNonlands={oppBattlefieldNonlands}
        battlefieldLands={oppBattlefieldLands}
        commandCards={oppCommandCards}
        isCommanderGame={isCommanderGame}
        showLifeMenu={false}
        onSelectOpponent={undefined}
        onLifeChange={handleLifeChange}
        onPoisonChange={handlePoisonChange}
        onToggleLifeMenu={() => {}}
        onCardContextMenu={(cardId, cardName) => {
          selectedCardForCounters = { id: cardId, name: cardName };
          showCounterDialog = true;
        }}
      />
    {/each}
  </div>

  <!-- Single opponent with cycling for small screens -->
  <div class="opponents-grid-small">
    {#if selectedOpponent()}
      {@const opponent = selectedOpponent()!}
      <OpponentSection
        opponent={opponent}
        otherPlayers={otherPlayers}
        battlefieldNonlands={opponentBattlefieldNonlands()}
        battlefieldLands={opponentBattlefieldLands()}
        commandCards={opponentCommandCards()}
        isCommanderGame={isCommanderGame}
        showLifeMenu={showOpponentLifeMenu}
        onSelectOpponent={(playerId) => (selectedOpponentId = playerId)}
        onLifeChange={handleLifeChange}
        onPoisonChange={handlePoisonChange}
        onToggleLifeMenu={() => (showOpponentLifeMenu = !showOpponentLifeMenu)}
        onCardContextMenu={(cardId, cardName) => {
          selectedCardForCounters = { id: cardId, name: cardName };
          showCounterDialog = true;
        }}
      />
    {/if}
  </div>
{/if}
```

**CSS Grid Classes (global.css lines added in playtest):**
```css
.opponents-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
  padding: 1rem;
}

.opponents-grid-large {
  display: grid;
}

.opponents-grid-small {
  display: none;
}

@media (max-width: 1023px) {
  .opponents-grid-large {
    display: none;
  }

  .opponents-grid-small {
    display: flex;
    flex-direction: column;
  }
}
```

**Derived State (playtest page lines 154-184):**
```typescript
const selectedOpponent = $derived(() => {
  if (otherPlayers.length === 0) return null;
  if (!selectedOpponentId || !otherPlayers.find((p) => p.playerId === selectedOpponentId)) {
    return otherPlayers[0];
  }
  return otherPlayers.find((p) => p.playerId === selectedOpponentId) || otherPlayers[0];
});

const opponentBattlefield = $derived(() => {
  const opponent = selectedOpponent();
  return opponent ? battlefield.filter((c) => c.controllerId === opponent.playerId) : [];
});

const opponentBattlefieldNonlands = $derived(() =>
  opponentBattlefield().filter((c) => !isLandPermanent(c.type))
);

const opponentBattlefieldLands = $derived(() =>
  opponentBattlefield().filter((c) => isLandPermanent(c.type))
);

const opponentCommandCards = $derived(() => {
  const opponent = selectedOpponent();
  return opponent ? commandCards.filter((c) => (c.ownerId || c.controllerId) === opponent.playerId) : [];
});
```

### Implementation Steps

1. **Add OpponentSection import:**
   ```svelte
   import OpponentSection from '$lib/components/game/OpponentSection.svelte';
   ```

2. **Add state variables:**
   ```typescript
   let selectedOpponentId = $state<string | null>(null);
   let showOpponentLifeMenu = $state(false);
   ```

3. **Copy derived state (reuse isLandPermanent from Phase 2):**
   ```typescript
   const otherPlayers = $derived($opponents);

   const selectedOpponent = $derived(() => {
     if (otherPlayers.length === 0) return null;
     if (!selectedOpponentId || !otherPlayers.find((p) => p.playerId === selectedOpponentId)) {
       return otherPlayers[0];
     }
     return otherPlayers.find((p) => p.playerId === selectedOpponentId) || otherPlayers[0];
   });

   const opponentBattlefield = $derived(() => {
     const opponent = selectedOpponent();
     return opponent ? $battlefield.filter((c) => c.controllerId === opponent.playerId) : [];
   });

   const opponentBattlefieldNonlands = $derived(() =>
     opponentBattlefield().filter((c) => !isLandPermanent(c.type))
   );

   const opponentBattlefieldLands = $derived(() =>
     opponentBattlefield().filter((c) => isLandPermanent(c.type))
   );

   const opponentCommandCards = $derived(() => {
     const opponent = selectedOpponent();
     return opponent ? $command.filter((c) => (c.ownerId || c.controllerId) === opponent.playerId) : [];
   });
   ```

4. **Add CSS classes to global.css:**
   ```css
   /* Copy from playtest global.css (search for .opponents-grid) */
   .opponents-grid {
     display: grid;
     grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
     gap: 1rem;
     padding: 1rem;
   }

   .opponents-grid-large {
     display: grid;
   }

   .opponents-grid-small {
     display: none;
   }

   @media (max-width: 1023px) {
     .opponents-grid-large {
       display: none;
     }

     .opponents-grid-small {
       display: flex;
       flex-direction: column;
     }
   }
   ```

5. **Replace existing opponent display:**
   - Find current OpponentPanel usage
   - Replace with conditional layout (1v1 vs multi-opponent)
   - Copy exact markup from playtest page lines 1376-1449
   - Wire up event handlers

6. **Adapt life/poison handlers for opponents:**
   ```typescript
   function handleLifeChange(delta: number, playerId?: string): void {
     const targetPlayerId = playerId || localPlayerId;
     if (!targetPlayerId) return;
     modifyLife(gameId, targetPlayerId, delta);
   }

   function handlePoisonChange(delta: number, playerId?: string): void {
     const targetPlayerId = playerId || localPlayerId;
     if (!targetPlayerId) return;
     const player = $players.find((p) => p.playerId === targetPlayerId);
     if (!player) return;
     setPlayerCounter(gameId, targetPlayerId, 'poison', (player.poison || 0) + delta);
   }
   ```

### Verification Checklist

- [ ] OpponentSection renders for single opponent (1v1)
- [ ] OpponentSection grid renders for 3-4 players (multiplayer)
- [ ] Grid shows all opponents on large screens (>1024px)
- [ ] Single opponent cycling shows on small screens (<1024px)
- [ ] Opponent dropdown selector works in cycling mode
- [ ] Opponent life +/-1 buttons work
- [ ] Opponent life menu opens/closes
- [ ] Opponent poison counter displays correctly
- [ ] Opponent battlefield separates lands/nonlands
- [ ] Opponent command zone displays (if Commander)
- [ ] Right-click on opponent cards opens context menu

### Anti-Pattern Guards

- ✅ Use `modifyLife(gameId, playerId, delta)` with opponent's playerId
- ✅ Use `setPlayerCounter(gameId, playerId, 'poison', value)` with opponent's playerId
- ✅ Opponent cards are read-only (no tap/untap, only context menu)
- ✅ Life/poison changes target opponent, not local player

---

## Phase 4: Integrate DeckContextMenu Component

**Goal:** Add right-click context menu to library zone with nested submenus.

### Documentation References

**Component Location:** `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/DeckContextMenu.svelte`

**Usage Pattern (playtest page lines 88-99, 566-620, 1576-1584):**

**State Variables:**
```typescript
let showDeckContextMenu = $state(false);
let deckContextMenuPosition = $state<{ x: number; y: number }>({ x: 0, y: 0 });
```

**Action Definitions:**
```typescript
const deckContextMenuActions = $derived<MenuAction[]>(
  !me ? [] : [
    {
      label: 'Draw Cards',
      submenu: [
        { label: '1 Card', onClick: () => handleDrawN(1) },
        { label: '2 Cards', onClick: () => handleDrawN(2) },
        { label: '3 Cards', onClick: () => handleDrawN(3) },
        { label: '7 Cards', onClick: () => handleDrawN(7) },
        { label: 'Custom...', onClick: () => showNumberInput('Draw N Cards', 1, handleDrawN) }
      ]
    },
    { divider: true },
    {
      label: 'Search Library',
      onClick: () => { showDeckSearch = true; }
    },
    {
      label: 'Shuffle Library',
      onClick: handleShuffleLibrary
    }
  ]
);
```

**Handler Functions:**
```typescript
function handleDeckContextMenu(event: MouseEvent): void {
  if (!me) return;
  deckContextMenuPosition = { x: event.clientX, y: event.clientY };
  showDeckContextMenu = true;
}

function handleDrawN(count: number): void {
  if (!me) return;
  drawCards(gameId, me.playerId, count);
  toast.success(`Drew ${count} card(s)`);
}

function handleShuffleLibrary(): void {
  if (!me) return;
  shuffleLibrary(gameId, me.playerId);
  toast.success('Shuffled library');
}
```

**Component Usage:**
```svelte
{#if showDeckContextMenu}
  <DeckContextMenu
    position={deckContextMenuPosition}
    deckCount={me?.libraryCount || 0}
    playerName={me?.name || 'You'}
    onClose={() => (showDeckContextMenu = false)}
    actions={deckContextMenuActions}
  />
{/if}
```

### Implementation Steps

1. **Add DeckContextMenu import:**
   ```svelte
   import DeckContextMenu from '$lib/components/game/DeckContextMenu.svelte';
   import type { MenuAction } from '$lib/components/game/DeckContextMenu.svelte';
   ```

2. **Add state variables:**
   ```typescript
   let showDeckContextMenu = $state(false);
   let deckContextMenuPosition = $state<{ x: number; y: number }>({ x: 0, y: 0 });
   ```

3. **Add handler function to PlayerInfoRow's onDeckContextMenu:**
   ```typescript
   function handleDeckContextMenu(event: MouseEvent): void {
     const me = $localPlayer;
     if (!me) return;
     deckContextMenuPosition = { x: event.clientX, y: event.clientY };
     showDeckContextMenu = true;
   }
   ```

4. **Copy action definitions (adapt for game API):**
   ```typescript
   const deckContextMenuActions = $derived<MenuAction[]>(() => {
     const me = $localPlayer;
     if (!me) return [];

     return [
       {
         label: 'Draw Cards',
         submenu: [
           { label: '1 Card', onClick: () => handleDrawN(1) },
           { label: '2 Cards', onClick: () => handleDrawN(2) },
           { label: '3 Cards', onClick: () => handleDrawN(3) },
           { label: '7 Cards', onClick: () => handleDrawN(7) }
         ]
       },
       { divider: true },
       {
         label: 'Search Library',
         onClick: () => {
           searchLibrary(gameId, me.playerId);
         }
       },
       {
         label: 'Shuffle Library',
         onClick: () => {
           shuffleLibrary(gameId, me.playerId);
           toast.success('Shuffled library');
         }
       }
     ];
   });

   function handleDrawN(count: number): void {
     const me = $localPlayer;
     if (!me) return;
     drawCards(gameId, me.playerId, count);
     toast.success(`Drew ${count} card(s)`);
   }
   ```

5. **Add DeckContextMenu component:**
   ```svelte
   {#if showDeckContextMenu}
     <DeckContextMenu
       position={deckContextMenuPosition}
       deckCount={$localPlayer?.libraryCount || 0}
       playerName={$localPlayer?.name || 'You'}
       onClose={() => (showDeckContextMenu = false)}
       actions={deckContextMenuActions()}
     />
   {/if}
   ```

6. **Wire up PlayerInfoRow's onDeckContextMenu prop:**
   ```svelte
   <PlayerInfoRow
     ...
     onDeckContextMenu={handleDeckContextMenu}
   />
   ```

### Verification Checklist

- [ ] Right-click on library zone opens context menu
- [ ] Context menu shows at cursor position
- [ ] Draw Cards submenu opens on hover
- [ ] Draw N actions work and update game state
- [ ] Search Library action works (if API exists)
- [ ] Shuffle Library action works
- [ ] Context menu closes on outside click
- [ ] Context menu closes on Escape key
- [ ] Submenus flip left if too close to right edge

### Anti-Pattern Guards

- ✅ Use `drawCards(gameId, playerId, count)` API - exists in game API
- ✅ Use `shuffleLibrary(gameId, playerId)` API - exists in game API
- ✅ Use `searchLibrary(gameId, playerId)` API - exists in game API
- ❌ DO NOT implement mill or scry actions - these are playtest-only (no game API)
- ❌ DO NOT use NumberInputDialog for custom counts in Phase 4 - add in Phase 5

---

## Phase 5: Add Supporting Dialogs

**Goal:** Integrate NumberInputDialog for enhanced UX.

### Documentation References

**NumberInputDialog Location:** `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/NumberInputDialog.svelte`

**Usage Pattern (playtest page lines 91-98, 547-564, 1587-1599):**
```typescript
let showNumberInputDialog = $state(false);
let numberInputDialogConfig = $state<{
  title: string;
  defaultValue: number;
  min: number;
  max: number;
  onConfirm: (value: number) => void;
} | null>(null);

function showNumberInput(title: string, defaultValue: number, onConfirm: (value: number) => void): void {
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
```

```svelte
{#if showNumberInputDialog && numberInputDialogConfig}
  <NumberInputDialog
    title={numberInputDialogConfig.title}
    defaultValue={numberInputDialogConfig.defaultValue}
    min={numberInputDialogConfig.min}
    max={numberInputDialogConfig.max}
    onConfirm={numberInputDialogConfig.onConfirm}
    onCancel={() => {
      showNumberInputDialog = false;
      numberInputDialogConfig = null;
    }}
  />
{/if}
```

### Implementation Steps

1. **Add imports:**
   ```svelte
   import NumberInputDialog from '$lib/components/game/NumberInputDialog.svelte';
   ```

2. **Add state:**
   ```typescript
   let showNumberInputDialog = $state(false);
   let numberInputDialogConfig = $state<{
     title: string;
     defaultValue: number;
     min: number;
     max: number;
     onConfirm: (value: number) => void;
   } | null>(null);
   ```

3. **Add helper function:**
   ```typescript
   function showNumberInput(title: string, defaultValue: number, onConfirm: (value: number) => void): void {
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
   ```

4. **Update deckContextMenuActions to include Custom:**
   ```typescript
   {
     label: 'Draw Cards',
     submenu: [
       { label: '1 Card', onClick: () => handleDrawN(1) },
       { label: '2 Cards', onClick: () => handleDrawN(2) },
       { label: '3 Cards', onClick: () => handleDrawN(3) },
       { label: '7 Cards', onClick: () => handleDrawN(7) },
       { label: 'Custom...', onClick: () => showNumberInput('Draw N Cards', 1, handleDrawN) }
     ]
   }
   ```

5. **Add component:**
   ```svelte
   {#if showNumberInputDialog && numberInputDialogConfig}
     <NumberInputDialog
       title={numberInputDialogConfig.title}
       defaultValue={numberInputDialogConfig.defaultValue}
       min={numberInputDialogConfig.min}
       max={numberInputDialogConfig.max}
       onConfirm={numberInputDialogConfig.onConfirm}
       onCancel={() => {
         showNumberInputDialog = false;
         numberInputDialogConfig = null;
       }}
     />
   {/if}
   ```

### Verification Checklist

- [ ] NumberInputDialog opens when clicking "Custom..." in context menu
- [ ] Can input number with keyboard or +/- buttons
- [ ] Enter key confirms, Escape key cancels
- [ ] Number validation enforces min/max bounds
- [ ] Draw N Cards with custom count works
- [ ] Dialog closes after confirming or canceling

### Anti-Pattern Guards

- ✅ NumberInputDialog is UI-only, doesn't call any APIs directly
- ✅ `onConfirm` callback uses game API (drawCards, etc.)

---

## Phase 6: Add Keyboard Shortcuts

**Goal:** Add keyboard shortcuts for common game actions.

### Documentation References

**Shortcut Handler Location:** Playtest page lines 794-889

**Pattern:**
```typescript
function handleGlobalKeydown(event: KeyboardEvent): void {
  // Ignore if typing in input/textarea
  if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
    return;
  }

  const key = event.key.toLowerCase();

  switch (key) {
    case 'escape':
      // Close modals in priority order
      if (showMenu) {
        showMenu = false;
        event.preventDefault();
      }
      break;
    case '?':
      showKeyboardShortcuts = !showKeyboardShortcuts;
      event.preventDefault();
      break;
    case 'f':
      // Search library
      showDeckSearch = true;
      event.preventDefault();
      break;
    case 'x':
      // Untap all
      handleUntapAll();
      event.preventDefault();
      break;
    case 'c':
      // Draw card
      handleDrawCard();
      event.preventDefault();
      break;
    case 'v':
      // Shuffle library
      handleShuffleLibrary();
      event.preventDefault();
      break;
  }

  // Context-sensitive shortcuts when hovering over card
  if (hoveredCard) {
    switch (key) {
      case 'd':
        // Move to graveyard
        moveCardToZone(gameId, hoveredCard.id, 'GRAVEYARD');
        event.preventDefault();
        break;
      case 's':
        // Move to exile
        moveCardToZone(gameId, hoveredCard.id, 'EXILE');
        event.preventDefault();
        break;
      case 'r':
        // Move to hand
        moveCardToZone(gameId, hoveredCard.id, 'HAND');
        event.preventDefault();
        break;
      case 't':
        // Move to top of library
        moveCardToZone(gameId, hoveredCard.id, 'LIBRARY');
        event.preventDefault();
        break;
    }
  }
}
```

**Window Event Binding:**
```svelte
<svelte:window onkeydown={handleGlobalKeydown} />
```

### Implementation Steps

1. **Copy handleGlobalKeydown function (adapt for game API):**
   ```typescript
   function handleGlobalKeydown(event: KeyboardEvent): void {
     // Ignore if typing in input/textarea
     if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
       return;
     }

     const key = event.key.toLowerCase();

     switch (key) {
       case 'escape':
         // Close modals in priority order
         if (showDeckContextMenu) {
           showDeckContextMenu = false;
           event.preventDefault();
         } else if (showNumberInputDialog) {
           showNumberInputDialog = false;
           numberInputDialogConfig = null;
           event.preventDefault();
         }
         break;
       case '?':
         showKeyboardShortcuts = !showKeyboardShortcuts;
         event.preventDefault();
         break;
       case 'f':
         // Search library
         const me = $localPlayer;
         if (me) {
           searchLibrary(gameId, me.playerId);
           event.preventDefault();
         }
         break;
       case 'x':
         // Untap all
         if ($hasPriority) {
           untapAll(gameId, localPlayerId);
           event.preventDefault();
         }
         break;
       case 'c':
         // Draw card
         const me2 = $localPlayer;
         if (me2 && $hasPriority) {
           drawCards(gameId, me2.playerId, 1);
           toast.success('Drew a card');
           event.preventDefault();
         }
         break;
       case 'v':
         // Shuffle library
         const me3 = $localPlayer;
         if (me3) {
           shuffleLibrary(gameId, me3.playerId);
           toast.success('Shuffled library');
           event.preventDefault();
         }
         break;
     }

     // Context-sensitive shortcuts when hovering over card
     const hoveredCard = hoveredCardId ? getCardById(hoveredCardId) : null;
     if (hoveredCard && $hasPriority) {
       switch (key) {
         case 'd':
           moveCardToZone(gameId, hoveredCard.id, 'GRAVEYARD');
           event.preventDefault();
           break;
         case 's':
           moveCardToZone(gameId, hoveredCard.id, 'EXILE');
           event.preventDefault();
           break;
         case 'r':
           moveCardToZone(gameId, hoveredCard.id, 'HAND');
           event.preventDefault();
           break;
         case 't':
           moveCardToZone(gameId, hoveredCard.id, 'LIBRARY');
           event.preventDefault();
           break;
       }
     }
   }
   ```

2. **Add window event binding:**
   ```svelte
   <svelte:window onkeydown={handleGlobalKeydown} />
   ```

3. **Ensure hoveredCardId is tracked:**
   - Already implemented in Phase 2 (BattlefieldArea's onCardHover)

4. **Add keyboard shortcuts help (optional):**
   ```svelte
   import KeyboardShortcutsModal from '$lib/components/game/KeyboardShortcutsModal.svelte';

   let showKeyboardShortcuts = $state(false);

   <KeyboardShortcutsModal bind:open={showKeyboardShortcuts} mode="game" />
   ```

### Verification Checklist

- [ ] `?` key toggles keyboard shortcuts modal
- [ ] `Escape` closes modals in priority order
- [ ] `F` key triggers library search
- [ ] `X` key untaps all permanents (only when have priority)
- [ ] `C` key draws a card (only when have priority)
- [ ] `V` key shuffles library
- [ ] `D` moves hovered card to graveyard (only when hovering + have priority)
- [ ] `S` moves hovered card to exile
- [ ] `R` moves hovered card to hand
- [ ] `T` moves hovered card to top of library
- [ ] Shortcuts don't fire when typing in input fields
- [ ] All shortcuts call game API (not playtest store)

### Anti-Pattern Guards

- ✅ Use game API functions: `drawCards()`, `shuffleLibrary()`, `untapAll()`, `moveCardToZone()`
- ✅ Check `$hasPriority` before actions that require priority
- ✅ Use `getCardById(cardId)` to look up hovered card
- ❌ DO NOT call playtest store methods
- ❌ DO NOT allow shortcuts when typing in input fields

---

## Phase 7: Visual Polish & Responsive Layout

**Goal:** Add CSS for multi-opponent grid, drag-drop visual feedback, and responsive breakpoints.

### Documentation References

**CSS Classes from Playtest (global.css):**

**Opponent Grid (already added in Phase 3):**
```css
.opponents-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
  padding: 1rem;
}

.opponents-grid-large {
  display: grid;
}

.opponents-grid-small {
  display: none;
}

@media (max-width: 1023px) {
  .opponents-grid-large {
    display: none;
  }

  .opponents-grid-small {
    display: flex;
    flex-direction: column;
  }
}
```

**Drag-Drop Visual Feedback:**
```css
.drag-active {
  outline: 2px dashed var(--color-primary);
  outline-offset: 4px;
  background-color: var(--color-primary-alpha-10);
}

.drag-valid {
  outline-color: var(--color-success);
  background-color: var(--color-success-alpha-10);
}

.drag-ghost {
  position: fixed;
  pointer-events: none;
  z-index: 10000;
  transform: translate(-50%, -50%);
  transition: none;
}

.drag-ghost-card {
  width: 180px;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
  transform: rotate(-5deg);
  opacity: 0.9;
}

.drag-ghost-card.valid {
  transform: rotate(0deg);
  opacity: 1;
}

.drag-ghost-image {
  width: 100%;
  height: auto;
  border-radius: 8px;
}
```

**Drag Ghost Component (playtest page lines 1779-1790):**
```svelte
{#if isDragging && dragCardName}
  {@const dragImageUrl = getScryfallImageUrl(dragCardName, 'small')}
  <div class="drag-ghost" style="left: {dragPos.x}px; top: {dragPos.y}px;">
    <div class="drag-ghost-card" class:valid={isOverValidDrop}>
      {#if dragImageUrl}
        <img src={dragImageUrl} alt={dragCardName} class="drag-ghost-image" draggable="false" />
      {:else}
        <span class="drag-ghost-name">{dragCardName}</span>
      {/if}
    </div>
  </div>
{/if}
```

### Implementation Steps

1. **Add drag-drop CSS to global.css:**
   ```bash
   # Open /Users/aron/dev/opensource/mage/mage-client-web/src/lib/styles/global.css
   # Search for existing drag classes or add new section
   # Copy .drag-active, .drag-valid, .drag-ghost classes from above
   ```

2. **Add drag ghost component to game page:**
   ```svelte
   {#if $isDragging && $draggedCardName}
     {@const dragImageUrl = getScryfallImageUrl($draggedCardName, 'small')}
     <div class="drag-ghost" style="left: {$dragPosition.x}px; top: {$dragPosition.y}px;">
       <div class="drag-ghost-card" class:valid={$isOverValidDropZone}>
         {#if dragImageUrl}
           <img src={dragImageUrl} alt={$draggedCardName} class="drag-ghost-image" draggable="false" />
         {:else}
           <span class="drag-ghost-name">{$draggedCardName}</span>
         {/if}
       </div>
     </div>
   {/if}
   ```

3. **Ensure drop zones have drag-active/drag-valid classes:**
   - Verify battlefield, graveyard, exile, hand, library zones apply classes
   - Check that components use `class:drag-active={isDragging}` and `class:drag-valid={isDragging && isOverValidDrop && dropZone === 'zone-name'}`

4. **Test responsive layout:**
   - Open game in browser
   - Resize to <1024px width
   - Verify opponent grid switches to cycling mode
   - Verify opponent dropdown appears

### Verification Checklist

- [ ] Drag ghost displays when dragging card
- [ ] Drag ghost follows cursor position
- [ ] Drag ghost changes appearance when over valid drop zone
- [ ] Drop zones highlight with dashed outline when dragging
- [ ] Drop zones show green outline when valid, blue when invalid
- [ ] Opponent grid displays correctly on large screens
- [ ] Opponent cycling displays correctly on small screens
- [ ] Responsive breakpoint works at 1024px width
- [ ] All CSS variables resolve correctly

### Anti-Pattern Guards

- ✅ Use existing CSS variables (--color-primary, --color-success, etc.)
- ✅ Use drag-drop store derived values ($isDragging, $dragPosition, etc.)
- ✅ Apply `draggable="false"` to drag ghost image
- ✅ Use `pointer-events: none` on drag ghost (prevents interference)

---

## Phase 8: Final Verification & Testing

**Goal:** Verify all integrations work together and follow game rules.

### Verification Checklist

#### Component Integration
- [ ] PlayerInfoRow displays and updates correctly
- [ ] BattlefieldArea displays and updates correctly
- [ ] OpponentSection displays and updates correctly
- [ ] DeckContextMenu opens and functions correctly
- [ ] NumberInputDialog works for custom counts
- [ ] All components use Lucide icons (not emoji)

#### Drag-Drop System
- [ ] Drag threshold (5px) works - clicks don't start drags
- [ ] Dragging from hand to battlefield works
- [ ] Dragging between zones works (graveyard, exile, hand, library)
- [ ] Drop zones only accept cards when have priority
- [ ] Drop zones reject invalid moves
- [ ] Drag ghost displays correctly
- [ ] Drop zones highlight correctly

#### Game API Integration
- [ ] All actions use game API (not playtest store)
- [ ] Life changes sync to server
- [ ] Poison changes sync to server
- [ ] Card movements sync to server
- [ ] Tap/untap syncs to server
- [ ] Draw cards syncs to server
- [ ] Shuffle library syncs to server
- [ ] Library search syncs to server
- [ ] No console errors about undefined methods

#### Keyboard Shortcuts
- [ ] All shortcuts documented in Phase 6 work
- [ ] Shortcuts respect priority
- [ ] Shortcuts don't fire in input fields
- [ ] `?` key shows shortcuts help
- [ ] Escape key closes modals

#### Multi-Opponent Layout
- [ ] 1v1 games show single opponent
- [ ] 3-4 player games show grid on large screens
- [ ] 3-4 player games show cycling on small screens
- [ ] Opponent selector dropdown works
- [ ] Responsive breakpoint works at 1024px

#### Visual Polish
- [ ] CSS grid layout works
- [ ] Drag-drop visual feedback works
- [ ] All colors and theming consistent
- [ ] No layout shifts or jumpiness
- [ ] Mobile responsive design works

### Anti-Pattern Verification

Run these grep commands to verify no anti-patterns exist:

```bash
# Check for playtest store usage in game page
grep -n "playtestGameStore\." src/routes/\(protected\)/game/\[id\]/+page.svelte

# Check for invented APIs
grep -n "gameStore\.switchControlSeat\|gameStore\.scryCards\|gameStore\.createToken" src/routes/\(protected\)/game/\[id\]/+page.svelte

# Check for direct zone manipulation
grep -n "playtestGameStore\.moveCardToZone\|playtestGameStore\.tapCard" src/routes/\(protected\)/game/\[id\]/+page.svelte

# Verify game API usage
grep -n "modifyLife\|setPlayerCounter\|drawCards\|shuffleLibrary\|tapUntap\|moveCardToZone" src/routes/\(protected\)/game/\[id\]/+page.svelte
```

All checks should show game API usage (from `/lib/api/game.ts` and `/lib/api/direct-actions.ts`), NOT playtest store methods.

### Testing Scenarios

1. **1v1 Game:**
   - Start 1v1 game
   - Verify single opponent displays
   - Test life changes (self and opponent)
   - Test drag-drop from hand to battlefield
   - Test keyboard shortcuts
   - Test deck context menu

2. **Multiplayer Game (3-4 players):**
   - Start multiplayer game
   - Verify opponent grid on large screen
   - Verify opponent cycling on small screen
   - Test responsive breakpoint (resize browser)
   - Test opponent selector dropdown

3. **Commander Game:**
   - Start Commander format game
   - Verify command zone displays
   - Verify command zone accepts drops
   - Test dragging commander from command zone

4. **Drag-Drop Edge Cases:**
   - Click card (shouldn't drag - <5px movement)
   - Drag and release outside valid zone (should cancel)
   - Drag when don't have priority (should be rejected)
   - Drag to same zone (should be rejected)

5. **Keyboard Shortcuts Edge Cases:**
   - Type in input field (shortcuts shouldn't fire)
   - Press shortcuts without priority (should be ignored)
   - Hover card and press zone shortcuts (d/s/r/t)
   - Press Escape multiple times (closes modals in order)

### Rollback Plan

If issues arise, phases can be reverted independently:

- **Phase 6 (Keyboard Shortcuts):** Comment out `<svelte:window onkeydown={...} />`
- **Phase 5 (Dialogs):** Remove dialog components and state variables
- **Phase 4 (DeckContextMenu):** Remove component and state, restore simple right-click
- **Phase 3 (OpponentSection):** Restore OpponentPanel component
- **Phase 2 (BattlefieldArea):** Restore previous battlefield markup
- **Phase 1 (PlayerInfoRow):** Restore previous player stat layout

Each phase is independently revertable without affecting others.

---

## Summary

**Total Phases:** 8 (0 = Discovery, 1-7 = Implementation, 8 = Verification)

**Copy-Based Approach:**
- All phases copy existing patterns from playtest
- All APIs documented with sources
- All anti-patterns explicitly listed
- All verification steps included

**Key Success Factors:**
1. Use game API functions, not playtest store methods
2. Check `$hasPriority` before actions requiring priority
3. Use drag threshold (5px) for better UX
4. Support multi-opponent layout with responsive design
5. Follow existing component APIs exactly

**Documentation Complete:** All component props, store APIs, and game APIs documented with file paths and line numbers.
