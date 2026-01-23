# Playtest-First Game View Replacement Plan

## Overview

Replace the current rules-enforced game view with a multiplayer-adapted version of the playtest engine. This plan uses the proven playtest UI/state management as the foundation and adds only the minimal server synchronization needed for multiplayer.

**Core Philosophy**: Copy the working playtest implementation, add server sync, discard all rules enforcement.

---

## Phase 0: Documentation & Pattern Analysis

### Playtest Engine Reference Patterns

**Core Files to Copy From**:
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/playtest-game.ts` (1,257 lines)
  - Lines 380-401: Immutable state update pattern with `update()`
  - Lines 492-1151: All game operations (drawCards, tapCard, moveCardToZone, etc.)
  - Lines 25-59: State interface definitions

- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/utils/playtest-helpers.ts` (303 lines)
  - Lines 30-63: `findCardInState()` - zone searching
  - Lines 120-170: `removeCardFromZone()` - immutable removal
  - Lines 175-243: `addCardToZone()` - immutable addition

- `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/playtest/+page.svelte` (1,992 lines)
  - Lines 440-760: Event handler patterns
  - Lines 251-350: Derived state computations
  - Lines 955-1050: Keyboard shortcuts

**Backend Direct Action API** (Already Implemented):
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/actions.go`
  - String command parsing: `TAP:{cardId}`, `MOVE:{cardId}:{zone}`, etc.
  - All playtest operations already supported server-side
  - Single GameEngine (rules-light) architecture

**Frontend API Wrapper**:
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/api/direct-actions.ts` (264 lines)
  - All operations map to `SendPlayerString` gRPC call
  - Pattern: `sendPlayerString(gameId, playerId, "COMMAND:args")`

### Anti-Patterns to Avoid

**DO NOT**:
- ❌ Copy rules-enforcement logic from game view (priority, validation, prompts)
- ❌ Use `CallbackMethod.*` event handlers (prompts, choices, mana payment)
- ❌ Implement combat validation (DeclareAttackers/Blockers components)
- ❌ Add priority system (PriorityActionBar, hasPriority checks)
- ❌ Copy phase progression logic (PhaseIndicator, phase-based enabling)

**DO**:
- ✅ Copy playtest store operations verbatim
- ✅ Copy playtest UI component integration patterns
- ✅ Use direct-actions.ts API (already server-compatible)
- ✅ Copy keyboard shortcuts and event handlers
- ✅ Add multiplayer UI from game view (OpponentSection, PlayerInfoRow)

---

## Phase 1: Backend - Simplify to Rules-Light Engine

### Goal
Replace MageEngine with minimal state management engine that matches playtest operations.

### Implementation Steps

#### 1.1 Create New Engine Files

**File**: `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/engine.go`

Copy state structure from playtest-game.ts lines 25-59, translate to Go:
```go
type Engine struct {
    mu          sync.RWMutex
    games       map[string]*GameState
    notifyFn    NotificationHandler
}

type GameState struct {
    GameID       string
    Players      map[string]*Player
    Zones        map[string]*Zone
    Turn         int
    ActivePlayer string
    Messages     []ActionLog
    StartedAt    time.Time
}

type Card struct {
    ID           string
    Name         string
    OwnerID      string
    ControllerID string
    Zone         string
    Tapped       bool
    FaceUp       bool
    Counters     map[string]int
    Properties   map[string]interface{}
}

type Player struct {
    ID           string
    Name         string
    Life         int
    Poison       int
    Energy       int
    LibraryCount int
    HandCount    int
}
```

**Pattern Source**: Copy from `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/playtest-game.ts` lines 25-59

#### 1.2 Implement Core Operations

**File**: `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/actions.go`

Implement each operation from playtest-game.ts lines 492-1151:

```go
func (e *Engine) MoveCard(gameID, cardID, targetZone string) error {
    // Pattern: playtest-game.ts lines 574-607 moveCardToZone()
    state := e.getGameState(gameID)
    card, sourceZone := findCardInState(state, cardID)

    // Token rule: tokens leaving battlefield cease to exist
    if strings.HasPrefix(cardID, "token-") && sourceZone == "battlefield" {
        return e.removeCard(gameID, cardID)
    }

    removeCardFromZone(state, cardID, sourceZone)
    addCardToZone(state, card, targetZone)
    e.logAction(gameID, "CARD_MOVED", cardID)
    e.broadcast(gameID)
    return nil
}

func (e *Engine) TapCard(gameID, cardID string, tapped bool) error {
    // Pattern: playtest-game.ts lines 612-630 tapCard()
    state := e.getGameState(gameID)
    card := findCard(state, cardID)
    card.Tapped = tapped
    e.logAction(gameID, "CARD_TAPPED", cardID)
    e.broadcast(gameID)
    return nil
}

// ... implement all operations from playtest-game.ts:
// - drawCards (line 492)
// - playCard (line 532)
// - modifyLife (line 668)
// - setPlayerCounter (line 684)
// - shuffleLibrary (line 705)
// - createToken (line 759)
// - addCounter (line 810)
// - removeCounter (line 845)
// - setCounter (line 887)
// - millCards (line 923)
// - scryCards (line 983)
// - nextTurn (line 1076)
```

**Pattern Source**: Translate each function from playtest-game.ts lines 492-1151

#### 1.3 Implement Rollback System

**File**: `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/rollback.go`

```go
type Bookmark struct {
    ID        string
    GameID    string
    State     *GameState
    Timestamp time.Time
}

func (e *Engine) BookmarkState(gameID string) string {
    // Create snapshot of current state
    // Return bookmark ID
}

func (e *Engine) RestoreState(gameID, bookmarkID string) error {
    // Restore game to bookmarked state
    // Broadcast update to all clients
}
```

**Pattern Source**: Extend current bookmark system in mage_engine.go lines 740-772

#### 1.4 Update Game Manager

**File**: `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/manager.go`

Replace MageEngine with Engine:
```go
// Line ~401: Replace engine creation
engine := NewEngine(logger, cardRepo)

// Keep existing GameEngine interface (lines 319+)
// New Engine must implement:
// - StartGame()
// - StartGameWithDecks()
// - ProcessAction()
// - GetGameView()
// - EndGame()
```

**Pattern Source**: Keep interface from manager.go lines 319-350, wire to new Engine

#### 1.5 Delete Old Engine Files

Remove:
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/mage_engine.go` (13,786 lines)
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/engine_combat.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/engine_priority.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/engine_stack.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/engine_layers.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/engine_events.go`

### Verification Checklist

- [ ] Engine compiles and implements GameEngine interface
- [ ] All playtest operations have Go equivalents
- [ ] Action logging works (Messages array populated)
- [ ] Bookmark/restore functions correctly
- [ ] WebSocket notifications broadcast on state changes
- [ ] Game manager successfully creates games with new engine

---

## Phase 2: Backend - Server State Synchronization

### Goal
Add WebSocket broadcast and state sync for multiplayer.

### Implementation Steps

#### 2.1 Implement Broadcast System

**File**: `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/engine.go`

```go
func (e *Engine) broadcast(gameID string) {
    state := e.getGameState(gameID)

    // Send full state to all players
    for _, playerID := range state.Players {
        view := e.GetGameView(gameID, playerID)
        e.notifyFn.NotifyGameStateChange(playerID, view)
    }
}
```

**Pattern Source**: Adapt from mage_engine.go notification system (lines 797-850)

#### 2.2 Add Hidden Information Filtering

**File**: `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/view.go`

```go
func (e *Engine) GetGameView(gameID, playerID string) *GameView {
    state := e.getGameState(gameID)

    view := &GameView{
        GameID:       gameID,
        MyHand:       getPlayerCards(state, playerID, "hand"),
        MyLibrary:    getLibraryCount(state, playerID),  // Count only
        Battlefield:  getAllBattlefieldCards(state),     // Public zone
        Graveyard:    getAllGraveyardCards(state),       // Public zone
        Exile:        getAllExileCards(state),           // Public zone
        Opponents:    getOpponentViews(state, playerID), // Hidden hands
    }

    return view
}

func getOpponentViews(state *GameState, viewerID string) []*OpponentView {
    // Return hand counts, not actual cards
    // Return battlefield cards (public)
    // Return library counts (hidden)
}
```

**Pattern Source**: Similar to current GetGameView in mage_engine.go, but simplified (no prompts, no priority)

### Verification Checklist

- [ ] WebSocket broadcasts on every state change
- [ ] All players receive updates
- [ ] Opponent hands are hidden (counts only)
- [ ] Opponent libraries are hidden (counts only)
- [ ] Public zones visible to all (battlefield, graveyard, exile)

---

## Phase 3: Frontend - Copy Playtest Store for Multiplayer

### Goal
Create multiplayer game store by copying playtest store and adding server sync.

### Implementation Steps

#### 3.1 Create Multiplayer Game Store

**File**: `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/multiplayer-game.ts`

Copy entire playtest-game.ts structure:
```typescript
// Copy lines 25-59: State interface
export interface MultiplayerGameState {
  gameId: string;
  activeControlSeat: string;
  players: PlaytestPlayer[];     // Same structure
  battlefield: CardView[];       // Same structure
  exile: CardView[];
  stack: CardView[];
  command: CardView[];
  turn: number;
  activePlayerId: string;
  isInitialized: boolean;
  log: PlaytestLogEntry[];

  // NEW: Server sync fields
  isConnected: boolean;
  pendingActions: string[];
}

// Copy lines 380-401: Store creation pattern
const { subscribe, update } = writable<MultiplayerGameState>(initialState);

// Copy all operations from lines 492-1151, but send to server:
function drawCards(playerId: string, count: number): void {
  // Send to server via direct-actions.ts
  directActions.drawCards(gameId, playerId, count);

  // Optimistic update (optional)
  update((state) => {
    // Apply change locally immediately
  });

  // Server will broadcast back, which updates all clients
}
```

**Pattern Source**: Copy from `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/playtest-game.ts` lines 25-1257

#### 3.2 Add WebSocket Sync

**File**: `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/multiplayer-game.ts`

```typescript
// Subscribe to WebSocket game state updates
websocketStore.subscribe((message) => {
  if (message.type === 'GAME_UPDATE') {
    update((state) => {
      // Apply server state directly
      return {
        ...state,
        ...message.gameState,
        isConnected: true
      };
    });
  }
});
```

**Pattern Source**: Adapt from `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/game.ts` lines 140-200 (WebSocket subscription)

#### 3.3 Wire Operations to Direct Actions API

**File**: `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/multiplayer-game.ts`

```typescript
import * as directActions from '$lib/api/direct-actions';

function moveCardToZone(cardId: string, targetZone: string): void {
  // Send to server (already implemented in direct-actions.ts)
  directActions.moveCard(gameId, playerId, cardId, targetZone);
}

function tapCard(cardId: string, tapped: boolean): void {
  directActions.tapUntap(gameId, playerId, cardId, tapped);
}

function modifyLife(playerId: string, delta: number): void {
  directActions.modifyPlayerLife(gameId, requestingPlayerId, playerId, delta);
}

// ... wire all operations from playtest-game.ts to direct-actions.ts
```

**Pattern Source**: Use existing API from `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/api/direct-actions.ts` lines 1-264

### Verification Checklist

- [ ] Multiplayer store matches playtest store structure
- [ ] All operations call direct-actions.ts API
- [ ] WebSocket updates apply to store
- [ ] State changes broadcast to all players
- [ ] Optimistic updates work (optional)

---

## Phase 4: Frontend - Copy Playtest UI Components

### Goal
Replace game page with playtest page adapted for multiplayer.

### Implementation Steps

#### 4.1 Copy Playtest Page Structure

**File**: `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/game/[id]/+page.svelte`

Replace entire file with copy of playtest page:
```typescript
<script lang="ts">
  // Copy from playtest/+page.svelte lines 1-110: Imports
  import { multiplayerGameStore } from '$lib/stores/multiplayer-game';
  import PlayerHand from '$lib/components/game/PlayerHand.svelte';
  import BattlefieldArea from '$lib/components/game/BattlefieldArea.svelte';
  // ... all other imports

  // Copy from playtest/+page.svelte lines 111-250: State initialization
  let loading = $state(true);
  let error = $state<string | null>(null);
  let hoveredCardId = $state<string | null>(null);

  // Copy from playtest/+page.svelte lines 251-350: Derived state
  const players = $derived($multiplayerGameStore.players);
  const me = $derived(players.find(p => p.playerId === localPlayerId));
  const battlefield = $derived($multiplayerGameStore.battlefield);

  // Copy from playtest/+page.svelte lines 440-760: Event handlers
  function handleLifeChange(delta: number, playerId?: string) {
    multiplayerGameStore.modifyLife(playerId ?? localPlayerId, delta);
  }

  function handleDrawCard() {
    multiplayerGameStore.drawCards(localPlayerId, 1);
  }

  // ... copy all handlers
</script>

<!-- Copy from playtest/+page.svelte lines 1100+: Template -->
<div class="game-container">
  <PlaytestHeader {players} ... />

  <BattlefieldArea
    battlefieldNonlands={myBattlefieldNonlands}
    battlefieldLands={myBattlefieldLands}
    ...
  />

  <PlayerHand cards={me?.hand ?? []} ... />
</div>
```

**Pattern Source**: Copy from `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/playtest/+page.svelte` lines 1-1992

#### 4.2 Add Multiplayer Components

Add opponent views from game view (keep these, discard rest):

**Components to Copy**:
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/OpponentSection.svelte`
  - Shows opponent battlefield and info

- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/PlayerInfoRow.svelte`
  - Multi-player life totals

- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/GameChatOverlay.svelte`
  - Chat functionality

**Integration**:
```svelte
<!-- Add opponent panels -->
{#each otherPlayers as opponent}
  <OpponentSection
    {opponent}
    onLifeChange={(delta) => handleLifeChange(delta, opponent.playerId)}
  />
{/each}

<!-- Add player info row -->
<PlayerInfoRow {players} {activePlayerId} />

<!-- Add chat -->
<GameChatOverlay gameId={$multiplayerGameStore.gameId} />
```

**Pattern Source**: Copy from `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/game/[id]/+page.svelte` lines 1100-1400 (opponent rendering)

#### 4.3 Remove Rules-Enforcement UI

**Delete These Components** (don't copy to new page):
- `PriorityActionBar.svelte` - Priority system
- `DeclareAttackers.svelte` - Combat validation
- `DeclareBlockers.svelte` - Combat validation
- `AssignDamage.svelte` - Damage validation
- `ManaPayment.svelte` - Mana prompts
- `XManaSelector.svelte` - X mana prompts
- `AbilitiesPanel.svelte` - Ability prompts
- All `CallbackMethod.*` event handlers

**Keep Shared Components**:
- `Card.svelte` - Card rendering
- `PlayerHand.svelte` - Hand display
- `Graveyard.svelte` - Graveyard
- `ExileZone.svelte` - Exile
- `BattlefieldArea.svelte` - Battlefield layout
- `ScryDialog.svelte` - Scry UI
- `TokenCreator.svelte` - Token creation

**Pattern Source**: Component list from game view analysis (see exploration report)

### Verification Checklist

- [ ] Page renders with playtest UI structure
- [ ] All zones display correctly (hand, battlefield, graveyard, etc.)
- [ ] Opponent sections show other players
- [ ] Keyboard shortcuts work (copied from playtest)
- [ ] Drag-drop works (copied from playtest)
- [ ] No priority/prompt UI visible
- [ ] No combat validation UI visible

---

## Phase 5: Frontend - Keyboard Shortcuts & Polish

### Goal
Copy playtest keyboard shortcuts and ensure all interactions work.

### Implementation Steps

#### 5.1 Copy Keyboard Handler

**File**: `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/game/[id]/+page.svelte`

```typescript
// Copy from playtest/+page.svelte lines 955-1050
function handleGlobalKeydown(event: KeyboardEvent) {
  // Skip if typing in input
  if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
    return;
  }

  switch (event.key.toLowerCase()) {
    case 'c':
      handleDrawCard();
      break;
    case 'v':
      handleShuffleLibrary();
      break;
    case 'x':
      handleUntapAll();
      break;
    case 'e':
      handleNextTurn();
      break;
    case ' ':
      event.preventDefault();
      // Space for pass/priority (can keep simple version)
      break;
    // ... copy all shortcuts
  }
}

$effect(() => {
  window.addEventListener('keydown', handleGlobalKeydown);
  return () => window.removeEventListener('keydown', handleGlobalKeydown);
});
```

**Pattern Source**: Copy from `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/playtest/+page.svelte` lines 955-1050

#### 5.2 Copy Drag-Drop System

**File**: `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/game/[id]/+page.svelte`

```typescript
// Copy from playtest/+page.svelte lines 196-223: Drop zone registration
import { registerDropZone } from '$lib/utils/drag-drop';

let battlefieldDropZone: HTMLElement;
let graveyardDropZone: HTMLElement;
let exileDropZone: HTMLElement;
let handDropZone: HTMLElement;

$effect(() => {
  if (battlefieldDropZone) {
    registerDropZone(battlefieldDropZone, 'battlefield', handleBattlefieldDrop);
  }
  if (graveyardDropZone) {
    registerDropZone(graveyardDropZone, 'graveyard', (cardId) => handleZoneDrop(cardId, 'G'));
  }
  // ... register all zones
});

// Copy handlers from lines 643-780
function handleBattlefieldCardMouseDown(cardId: string, cardName: string, e: MouseEvent) {
  // Drag start logic
}

function handleBattlefieldDrop(cardId: string) {
  multiplayerGameStore.moveCardToZone(cardId, 'BATTLEFIELD');
}

function handleZoneDrop(cardId: string, zone: string) {
  multiplayerGameStore.moveCardToZone(cardId, zone);
}
```

**Pattern Source**: Copy from `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/playtest/+page.svelte` lines 196-780

#### 5.3 Copy Context Menus

**File**: `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/game/[id]/+page.svelte`

```typescript
// Copy from playtest/+page.svelte lines 499-540
let showDeckContextMenu = $state(false);
let deckContextMenuX = $state(0);
let deckContextMenuY = $state(0);

function handleDeckContextMenu(event: MouseEvent) {
  event.preventDefault();
  deckContextMenuX = event.clientX;
  deckContextMenuY = event.clientY;
  showDeckContextMenu = true;
}

function handleDrawN(count: number) {
  multiplayerGameStore.drawCards(localPlayerId, count);
  showDeckContextMenu = false;
}

function handleMill(count: number) {
  multiplayerGameStore.millCards(localPlayerId, count);
  showDeckContextMenu = false;
}

function handleScry(count: number) {
  const session = multiplayerGameStore.scryCards(localPlayerId, count);
  // Show scry dialog
}
```

**Pattern Source**: Copy from `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/playtest/+page.svelte` lines 499-540

### Verification Checklist

- [ ] Keyboard shortcuts work (C=draw, V=shuffle, X=untap, E=next turn)
- [ ] Drag-drop works between all zones
- [ ] Context menus open on right-click
- [ ] Deck context menu shows draw/mill/scry options
- [ ] Card context menus work (tap, counters, etc.)

---

## Phase 6: Integration Testing

### Goal
Verify multiplayer functionality works end-to-end.

### Test Scenarios

#### 6.1 Basic Operations (2 Players)

**Setup**: Start game with 2 players, each with a deck

**Test**:
1. Player 1 draws card → Both players see hand count update
2. Player 1 plays card → Both players see battlefield update
3. Player 2 taps card → Both players see tap animation
4. Player 1 changes life → Both players see life total update
5. Player 2 moves card to graveyard → Both players see zone update

**Expected**: All operations sync via WebSocket to both clients

#### 6.2 Hidden Information (2 Players)

**Setup**: Start game with 2 players

**Test**:
1. Player 1 views own hand → Sees all cards
2. Player 2 views Player 1's hand → Sees only count
3. Player 1 plays card from hand → Player 2 sees card revealed on battlefield
4. Player 1 library count → Both see count, Player 2 doesn't see card names

**Expected**: Opponent hands/libraries hidden, public zones visible

#### 6.3 Combat (2 Players)

**Setup**: Both players have creatures on battlefield

**Test**:
1. Player 1 declares attacker (tap creature, no validation)
2. Player 2 declares blocker (tap creature, no validation)
3. Players manually apply damage (life changes, destroy creatures)
4. Clear combat state

**Expected**: No automated combat, manual state manipulation works

#### 6.4 Rollback (2 Players)

**Setup**: Game in progress, bookmark created

**Test**:
1. Player 1 requests rollback to bookmark
2. Player 2 receives consent dialog
3. Player 2 approves
4. Both players see state restored

**Expected**: Rollback system works with consent

#### 6.5 Multiplayer (4 Players)

**Setup**: 4-player game

**Test**:
1. All 4 players see synchronized state
2. Turn order cycles correctly
3. Each player sees 3 opponent views
4. Chat messages visible to all

**Expected**: Scales to 4 players correctly

### Verification Checklist

- [ ] 2-player game works end-to-end
- [ ] Hidden information correctly filtered
- [ ] Combat works (manual damage assignment)
- [ ] Rollback consent flow works
- [ ] 4-player game scales correctly
- [ ] WebSocket syncs state reliably
- [ ] No client-side state desync issues

---

## Phase 7: Cleanup & Documentation

### Goal
Remove old game view code, update documentation.

### Tasks

#### 7.1 Delete Old Files

Remove rules-enforcement components:
```bash
rm mage-client-web/src/lib/components/game/PriorityActionBar.svelte
rm mage-client-web/src/lib/components/game/DeclareAttackers.svelte
rm mage-client-web/src/lib/components/game/DeclareBlockers.svelte
rm mage-client-web/src/lib/components/game/AssignDamage.svelte
rm mage-client-web/src/lib/components/game/ManaPayment.svelte
rm mage-client-web/src/lib/components/game/XManaSelector.svelte
rm mage-client-web/src/lib/components/game/AbilitiesPanel.svelte
rm mage-client-web/src/lib/components/game/AbilityItem.svelte
rm mage-client-web/src/lib/stores/combat.ts
rm mage-client-web/src/lib/types/combat.ts
```

Remove old game store (replace with multiplayer-game.ts):
```bash
# Keep for reference initially, rename to game.legacy.ts
mv mage-client-web/src/lib/stores/game.ts mage-client-web/src/lib/stores/game.legacy.ts
```

#### 7.2 Update Documentation

**File**: `/Users/aron/dev/opensource/mage/docs/GAME_ARCHITECTURE.md`

Document new architecture:
- Rules-light backend engine
- Playtest-based frontend
- Direct manipulation API
- WebSocket state sync
- Rollback system

**File**: `/Users/aron/dev/opensource/mage/README.md`

Update feature list:
- Remove: "Full MTG rules enforcement"
- Add: "Rules-light multiplayer engine (Untap.in style)"
- Add: "Player-controlled game state with rollback"

#### 7.3 Update Ticket

**File**: `/Users/aron/dev/opensource/mage/docs/tickets/todo/003-implement-direct-engine.md`

Rewrite to reflect playtest-first approach:
- Backend: Replace with rules-light engine matching playtest operations
- Frontend: Copy playtest UI, add multiplayer components
- Migration: Full replacement (no backward compatibility)

### Verification Checklist

- [ ] Old components deleted
- [ ] Documentation updated
- [ ] README reflects new architecture
- [ ] Ticket updated with implementation details
- [ ] Code comments reference playtest patterns

---

## Success Criteria

### Backend
- ✅ New Engine implements all playtest operations
- ✅ State syncs via WebSocket to all clients
- ✅ Rollback system works
- ✅ No rules enforcement logic
- ✅ Codebase reduced from 13,786 to ~2,500 lines

### Frontend
- ✅ Game page uses playtest UI structure
- ✅ All playtest features work (drag-drop, shortcuts, context menus)
- ✅ Multiplayer components integrated (opponents, chat)
- ✅ No priority/prompt UI
- ✅ State syncs from server correctly

### Integration
- ✅ 2-4 player games work end-to-end
- ✅ Hidden information correctly enforced
- ✅ Manual combat works (no validation)
- ✅ Rollback consent flow works
- ✅ No client/server state desync

---

## Migration Notes

### Backward Compatibility
**NONE NEEDED** - Single user project, can break existing games

### Deployment Strategy
1. Deploy backend with new engine
2. Deploy frontend with new game page
3. Old games incompatible (expected)
4. Create new games with rules-light engine

### Risk Mitigation
- Backend already supports direct-actions API (lines 2590-2740 in mage_engine.go)
- Playtest UI proven to work client-side
- Only adding server sync layer (low risk)
- Can test with local multiplayer before production

---

## File Path Reference

### Backend Files to Create
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/engine.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/state.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/actions.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/rollback.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/view.go`

### Backend Files to Modify
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/manager.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/server/grpc_game.go`

### Backend Files to Delete
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/mage_engine.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/engine_*.go` (all engine files)

### Frontend Files to Create
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/multiplayer-game.ts`

### Frontend Files to Modify
- `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/game/[id]/+page.svelte`

### Frontend Files to Keep (Copy to New Page)
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/Card.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/PlayerHand.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/BattlefieldArea.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/Graveyard.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/ExileZone.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/OpponentSection.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/PlayerInfoRow.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/GameChatOverlay.svelte`

### Frontend Files to Delete
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/PriorityActionBar.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/DeclareAttackers.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/DeclareBlockers.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/AssignDamage.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/ManaPayment.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/XManaSelector.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/AbilitiesPanel.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/combat.ts`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/types/combat.ts`
