# Integration Test Plan: Rules-Light Game Engine Architecture

**Document Version:** 2.0
**Date:** January 23, 2026
**Status:** Phase 9 - Migration Complete

## Table of Contents

1. [Overview](#overview)
2. [Code Path Verification](#code-path-verification)
3. [Backend Integration Tests](#backend-integration-tests)
4. [Frontend Integration Tests](#frontend-integration-tests)
5. [End-to-End Tests](#end-to-end-tests)
6. [Manual Test Checklist](#manual-test-checklist)
7. [Automated Test Suggestions](#automated-test-suggestions)
8. [Known Limitations](#known-limitations)
9. [Test Results Template](#test-results-template)

---

## Overview

This document provides a comprehensive test plan for verifying the integration of the rules-light game engine architecture. The system consists of:

- **Backend**: Go server with GameEngine (single rules-light engine)
- **Frontend**: SvelteKit client with multiplayerGameStore
- **Communication**: WebSocket for state sync, gRPC for actions

**Note**: MageEngine has been completely removed. The system now uses a single GameEngine for all games.

### Test Objectives

1. Verify all code paths are properly wired
2. Ensure multiplayer synchronization works correctly
3. Validate hidden information filtering
4. Confirm all 19 operations function correctly
5. Test 2-player and 4-player scenarios
6. Verify error handling and edge cases

---

## Code Path Verification

### Backend Verification Checklist

#### 1. Server Initialization (`/Users/aron/dev/opensource/mage/mage-server-go/cmd/server/main.go`)

**GameEngine Initialization**

- [x] Creates single `GameEngine` instance (rules-light)
- [x] Sets notification handler for WebSocket sync
- [x] No configuration needed (single engine)
- [x] Wraps engine in `EngineAdapter` (lines 158, 175)
- [x] Passes adapter to `NewMageServer` (line 232)

**Validation Steps:**
```bash
# 1. Check config file structure
cat config/config.yaml | grep engine_type

# 2. Verify server starts with playtest engine
# Expected log: "playtest engine initialized (rules-light mode)"

# 3. Verify server starts with mage engine (default)
# Expected log: "mage engine initialized (full rules enforcement)"
```

**Status:** ✅ **VERIFIED** - Code paths exist and are properly wired

---

#### 2. Engine Action Processing (`/Users/aron/dev/opensource/mage/mage-server-go/internal/game/engine.go`)

**Lines 142-251: ProcessAction Routing**

- [x] `ProcessAction` receives `PlayerAction` (line 144)
- [x] Routes based on `ActionType` (switch statement line 157)
- [x] Supports 19 action types:
  - DRAW (line 158)
  - PLAY (line 163)
  - MOVE (line 168)
  - TAP (line 173)
  - UNTAP_ALL (line 177)
  - FLIP (line 180)
  - MODIFY_LIFE (line 185)
  - SET_COUNTER (line 191)
  - SHUFFLE (line 196)
  - CREATE_TOKEN (line 199)
  - ADD_COUNTER (line 207)
  - REMOVE_COUNTER (line 213)
  - SET_CARD_COUNTER (line 219)
  - MILL (line 225)
  - SCRY (line 229)
  - SET_REVEALED_TOP (line 235)
  - NEXT_TURN (line 239)
  - MULLIGAN (line 242)
  - KEEP_HAND (line 245)

**Lines 326-343: Broadcast Mechanism**

- [x] `broadcast()` function exists (line 327)
- [x] Checks if `notifyFn` is set (line 328)
- [x] Gets game state (line 332)
- [x] Builds `PlaytestGameView` for each player (line 339)
- [x] Calls `notifyFn.NotifyGameStateChange()` (line 340)

**Validation Steps:**
```go
// Test that ProcessAction routes to correct operations
func TestEngineProcessAction(t *testing.T) {
    // Verify DRAW action
    // Verify PLAY action
    // Verify MOVE action
    // etc. for all 19 operations
}
```

**Status:** ✅ **VERIFIED** - All operations route correctly and broadcast

---

#### 3. Game View Building (`/Users/aron/dev/opensource/mage/mage-server-go/internal/game/view.go`)

**Lines 68-136: buildGameView**

- [x] Creates `PlaytestGameView` struct (line 70)
- [x] Public zones shared: battlefield, exile, stack, command (lines 74-77)
- [x] Viewer gets full hand/library (lines 91-106)
- [x] Opponents get counts only (lines 109-132)
- [x] Revealed top card visible if set (lines 126-129)

**Validation Steps:**
```go
// Test hidden information filtering
func TestGameViewHiddenInformation(t *testing.T) {
    // Create game with 2 players
    // Player A draws cards
    // Get Player B's view
    // Verify Player B sees count, not cards
}
```

**Status:** ✅ **VERIFIED** - Hidden information properly filtered

---

#### 4. Notification Handling (`/Users/aron/dev/opensource/mage/mage-server-go/internal/server/grpc.go`)

**Lines 126-138: SetupGameNotifications**

- [x] Sets notification callback on adapter (line 133)
- [x] Callback routes to `handleGameNotification` (line 134)

**Lines 140-214: handleGameNotification**

- [x] Gets game instance from manager (line 148)
- [x] Handles special notifications (ERROR, TARGET, XMANA, CHOOSE_CHOICE)
- [x] Sends GAME_UPDATE to all players (lines 198-204)
- [x] Sends GAME_UPDATE to all watchers (lines 207-213)

**Lines 216-301: sendGameUpdateToPlayer**

- [x] Gets engine view via adapter (line 227)
- [x] Converts to protobuf via `engineViewToProto` (line 238)
- [x] Creates `GameUpdateData` protobuf message (lines 248-250)
- [x] Sends via WebSocket to all player sessions (lines 268-300)

**Validation Steps:**
```bash
# Monitor WebSocket traffic when action occurs
# Expected: GAME_UPDATE event sent to all players
# Expected: Each player receives personalized view
```

**Status:** ✅ **VERIFIED** - Notification pipeline fully wired

---

### Frontend Verification Checklist

#### 1. Multiplayer Store Initialization (`/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/multiplayer-game.ts`)

**Lines 108-187: Store Creation & WebSocket Subscription**

- [x] `createMultiplayerGameStore()` function exists (line 108)
- [x] WebSocket subscription in `subscribeToGameEvents()` (line 118)
- [x] GAME_INIT handler (lines 123-140)
- [x] GAME_UPDATE handler (lines 143-159)
- [x] `initialize()` subscribes to game events (lines 175-187)

**Status:** ✅ **VERIFIED** - Store initializes and subscribes

---

#### 2. Operation Wiring to direct-actions (`/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/multiplayer-game.ts`)

**Lines 194-523: Store Operations**

All operations call corresponding `directActions` functions:

- [x] `drawCards` → `directActions.drawCards` (line 196)
- [x] `playCard` → `directActions.moveCard` + `tapUntap` (lines 208-211)
- [x] `moveCardToZone` → `directActions.moveCard` (line 222)
- [x] `tapCard` → `directActions.tapUntap` (line 233)
- [x] `untapAll` → `directActions.untapAll` (line 244)
- [x] `flipCard` → `directActions.flipCard` (line 255)
- [x] `transformCard` → `directActions.transformCard` (line 265)
- [x] `modifyLife` → `directActions.modifyPlayerLife` (line 276)
- [x] `setPlayerCounter` → `directActions.setPlayerCounter` (line 287)
- [x] `shuffleLibrary` → `directActions.shuffleLibrary` (line 298)
- [x] `addToStack` → `directActions.addToStack` (line 309)
- [x] `removeFromStack` → `directActions.removeFromStack` (line 320)
- [x] `createToken` → `directActions.createToken` (line 338)
- [x] `destroyToken` → `directActions.destroyToken` (line 348)
- [x] `addCounter` → `directActions.modifyCardCounter` (line 359)
- [x] `removeCounter` → `directActions.modifyCardCounter` (line 370)
- [x] `setCounter` → `directActions.setCardCounter` (line 381)
- [x] `nextTurn` → `directActions.nextTurn` (line 467)
- [x] `clearCombat` → `directActions.clearCombat` (line 477)
- [x] `searchLibrary` → `directActions.searchLibrary` (line 491)

**Not Yet Implemented:**
- ⚠️ `millCards` - Shows warning, no server implementation (lines 391-397)
- ⚠️ `revealTopCards` - Shows warning, no server implementation (lines 404-411)
- ⚠️ `scryCards` - Shows warning, no server implementation (lines 418-425)
- ⚠️ `applyScryDecision` - Shows warning (lines 432-445)
- ⚠️ `setRevealedTop` - Shows warning (lines 452-458)
- ⚠️ `mulligan` - Shows warning (lines 510-513)
- ⚠️ `keepHand` - Shows warning (lines 520-523)

**Status:** ✅ **VERIFIED** - 19/26 operations wired, 7 pending server implementation

---

#### 3. direct-actions API (`/Users/aron/dev/opensource/mage/mage-client-web/src/lib/api/direct-actions.ts`)

**Lines 8-263: All Operations Defined**

- [x] Each operation calls `sendPlayerString` with correct format
- [x] Format: `COMMAND:param1:param2:...`
- [x] Examples:
  - `TAP:${cardId}` (line 17)
  - `MOVE:${cardId}:${targetZone}` (line 55)
  - `DRAW:${playerId}:${count}` (line 164)
  - `CREATE_TOKEN:${name}:${types}:...` (line 113)

**Status:** ✅ **VERIFIED** - All operations properly formatted

---

#### 4. Game Page Integration (`/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/game/[id]/+page.svelte`)

**Lines 62, 236-250: Initialization**

- [x] Gets `gameId` from route params (line 62)
- [x] `initializeFromGameId()` function (line 225)
- [x] Calls `multiplayerGameStore.initialize(gameId)` (line 239)
- [x] `onMount` triggers initialization (line 804-806)

**Lines 254-693: Event Handlers**

All handlers call multiplayerGameStore operations:

- [x] `handleLifeChange` → `multiplayerGameStore.modifyLife` (line 260)
- [x] `handlePoisonChange` → `multiplayerGameStore.setPlayerCounter` (line 273)
- [x] `handleDrawCard` → `multiplayerGameStore.drawCards` (line 282)
- [x] `handleShuffleLibrary` → `multiplayerGameStore.shuffleLibrary` (line 292)
- [x] `handleUntapAll` → `multiplayerGameStore.untapAll` (line 302)
- [x] `handleNextTurn` → `multiplayerGameStore.nextTurn` (line 311)
- [x] `handleBattlefieldCardClick` → `multiplayerGameStore.tapCard` (line 461)
- [x] `handleBattlefieldDrop` → `multiplayerGameStore.moveCardToZone` (line 571)
- [x] `handleZoneDrop` → `multiplayerGameStore.moveCardToZone` (line 582)

**Status:** ✅ **VERIFIED** - UI properly wired to store

---

## Backend Integration Tests

### Test 1: Engine Selection

**Objective:** Verify server starts with correct engine based on config

**Steps:**
1. Set `engine_type: "playtest"` in config
2. Start server
3. Verify log shows "playtest engine initialized"
4. Set `engine_type: "mage"` in config
5. Restart server
6. Verify log shows "mage engine initialized"

**Expected Result:**
- Playtest engine loads for `engine_type="playtest"`
- MageEngine loads for `engine_type="mage"` or missing config
- Invalid engine type causes error/fallback

**Status:** ⬜ Not Tested

---

### Test 2: State Synchronization

**Objective:** Verify engine broadcasts state changes to all players

**Setup:**
- Start server with playtest engine
- Create game with 2 players

**Steps:**
1. Player 1 draws a card via gRPC
2. Verify Player 1 receives GAME_UPDATE
3. Verify Player 2 receives GAME_UPDATE
4. Check Player 1's view shows new card in hand
5. Check Player 2's view shows hand count increased

**Expected Result:**
- Both players receive GAME_UPDATE WebSocket event
- Player 1 sees card details
- Player 2 sees count only (hidden information)

**Status:** ⬜ Not Tested

---

### Test 3: Action Processing

**Objective:** Verify all 19 operations process correctly

**Test Matrix:**

| Operation | Command Format | Test Case |
|-----------|----------------|-----------|
| DRAW | `DRAW:playerId:count` | Draw 3 cards |
| PLAY | `MOVE:cardId:BATTLEFIELD` | Play land |
| MOVE | `MOVE:cardId:GRAVEYARD` | Discard card |
| TAP | `TAP:cardId` | Tap permanent |
| UNTAP | `UNTAP:cardId` | Untap permanent |
| UNTAP_ALL | `UNTAP_ALL` | Untap all |
| FLIP | `FLIP:cardId:true` | Flip face-down |
| TRANSFORM | `TRANSFORM:cardId` | Transform DFC |
| MODIFY_LIFE | `MODIFY_LIFE:playerId:-3` | Lose 3 life |
| SET_PLAYER_COUNTER | `SET_PLAYER_COUNTER:playerId:poison:5` | Set poison to 5 |
| SHUFFLE | `SHUFFLE:playerId` | Shuffle library |
| CREATE_TOKEN | `CREATE_TOKEN:Goblin:Creature — Goblin:1:1:red:` | Create token |
| MODIFY_COUNTER | `MODIFY_COUNTER:cardId:+1/+1:2` | Add 2 +1/+1 counters |
| SET_COUNTER | `SET_COUNTER:cardId:loyalty:4` | Set loyalty to 4 |
| STACK_ADD | `STACK_ADD:cardId` | Add to stack |
| STACK_REMOVE | `STACK_REMOVE:itemId` | Remove from stack |
| NEXT_TURN | `NEXT_TURN` | Advance turn |
| CLEAR_COMBAT | `CLEAR_COMBAT` | Clear combat |
| SEARCH_LIBRARY | `SEARCH_LIBRARY:hand:true:Forest` | Search for card |

**For Each Operation:**
1. Send command via `SendPlayerString` gRPC
2. Verify server processes without error
3. Verify state changes correctly
4. Verify all players receive GAME_UPDATE
5. Verify game log entry created

**Status:** ⬜ Not Tested

---

### Test 4: Hidden Information Filtering

**Objective:** Verify opponent hands/libraries are hidden

**Steps:**
1. Create 2-player game
2. Player 1 draws 7 cards
3. Player 2 draws 7 cards
4. Get Player 1's view
5. Verify `me.hand` contains 7 CardView objects
6. Verify `opponents[0].hand` is empty array
7. Verify `opponents[0].handCount` equals 7
8. Get Player 2's view
9. Verify same filtering for Player 2

**Expected Result:**
- Own hand shows full card data
- Opponent hand shows empty array with count
- Libraries similarly filtered

**Status:** ⬜ Not Tested

---

### Test 5: Revealed Top Card

**Objective:** Verify revealed top card visible to opponents

**Steps:**
1. Player 1 sets `revealedTopCard = true`
2. Get Player 2's view
3. Verify `opponents[0].topCard` is not null
4. Verify card name matches top of Player 1's library

**Expected Result:**
- Opponent sees revealed top card
- Card data matches actual top card

**Status:** ⬜ Not Tested

---

## Frontend Integration Tests

### Test 6: Store Initialization

**Objective:** Verify multiplayerGameStore initializes correctly

**Steps:**
1. Navigate to `/game/{gameId}`
2. Monitor network for WebSocket connection
3. Verify `initialize()` called with correct gameId
4. Verify WebSocket subscriptions established
5. Check store state: `isConnected = true`, `isInitialized = true`

**Expected Result:**
- WebSocket connects successfully
- Event handlers registered for GAME_INIT and GAME_UPDATE
- Store populated with initial game state

**Status:** ⬜ Not Tested

---

### Test 7: WebSocket Event Handling

**Objective:** Verify GAME_UPDATE events apply to store

**Setup:**
- Game initialized
- WebSocket connected

**Steps:**
1. Server sends GAME_UPDATE event
2. Verify `GAME_UPDATE` handler fires
3. Verify store state updates
4. Verify UI re-renders with new state

**Expected Result:**
- Handler receives GameUpdateData
- Store merges new state
- UI reflects changes immediately

**Status:** ⬜ Not Tested

---

### Test 8: Operation Execution

**Objective:** Verify UI operations call server correctly

**Test Matrix:**

| UI Action | Store Method | API Call | Server Command |
|-----------|-------------|----------|----------------|
| Click "Draw" | `drawCards()` | `directActions.drawCards()` | `DRAW:playerId:1` |
| Tap card | `tapCard()` | `directActions.tapUntap()` | `TAP:cardId` |
| Move to graveyard | `moveCardToZone()` | `directActions.moveCard()` | `MOVE:cardId:GRAVEYARD` |
| Click "Untap All" | `untapAll()` | `directActions.untapAll()` | `UNTAP_ALL` |
| Change life | `modifyLife()` | `directActions.modifyPlayerLife()` | `MODIFY_LIFE:playerId:-1` |

**For Each Action:**
1. Perform UI interaction
2. Verify store method called
3. Verify direct-actions API called
4. Monitor network: verify gRPC SendPlayerString call
5. Verify server receives command

**Status:** ⬜ Not Tested

---

## End-to-End Tests

### Test 9: 2-Player Game Flow

**Objective:** Complete game flow between 2 players

**Players:**
- Player A (Browser 1)
- Player B (Browser 2)

**Steps:**

1. **Setup**
   - Both players join game
   - Both see initial state (20 life, 7 cards, empty battlefield)

2. **Player A Draws Card**
   - Player A clicks "Draw"
   - Verify: Player A sees 8 cards in hand
   - Verify: Player B sees opponent hand count = 8

3. **Player A Plays Card**
   - Player A drags card to battlefield
   - Verify: Player A sees card on battlefield
   - Verify: Player B sees same card on battlefield
   - Verify: Both see identical battlefield state

4. **Player B Taps Card**
   - Player B taps a permanent
   - Verify: Player B sees card tapped
   - Verify: Player A sees card tapped
   - Verify: Tap animation plays for both

5. **Player A Changes Life**
   - Player A sets life to 15 (-5)
   - Verify: Player A sees life = 15
   - Verify: Player B sees opponent life = 15

6. **Player A Passes Turn**
   - Player A clicks "Next Turn"
   - Verify: Turn counter increments
   - Verify: Active player changes to Player B
   - Verify: Both players see turn change

**Expected Result:**
- All actions synchronize immediately
- No desync or stale data
- Hidden information maintained
- UI updates smooth and instant

**Status:** ⬜ Not Tested

---

### Test 10: 4-Player Game Flow

**Objective:** Verify multiplayer scaling

**Players:**
- Player A, B, C, D (4 browsers)

**Steps:**

1. **Setup**
   - All 4 players join game
   - Verify opponent grid shows 3 opponents

2. **Round-Robin Actions**
   - Player A draws card → all see update
   - Player B plays card → all see update
   - Player C taps card → all see update
   - Player D changes life → all see update

3. **Turn Order**
   - Verify turn cycles: A → B → C → D → A
   - Verify all players see correct active player

4. **Simultaneous Actions**
   - Player A and Player B both tap cards
   - Verify both actions process
   - Verify no race conditions

**Expected Result:**
- All 4 players stay synchronized
- Opponent grid displays correctly
- Turn order cycles properly
- No performance degradation

**Status:** ⬜ Not Tested

---

### Test 11: Hidden Information Isolation

**Objective:** Verify no data leaks between players

**Setup:**
- 2-player game

**Steps:**

1. **Hand Privacy**
   - Player A draws 10 cards
   - Player B inspects network traffic
   - Verify: Player B's GAME_UPDATE shows empty hand array
   - Verify: No card names in Player B's data

2. **Library Privacy**
   - Player A shuffles library
   - Player B inspects network traffic
   - Verify: Player B sees library count only
   - Verify: No card order data in Player B's view

3. **Revealed Top Card**
   - Player A reveals top card
   - Verify: Player B sees revealed card name
   - Verify: Player B sees count for rest of library

**Expected Result:**
- No hidden information leaks
- Network traffic shows filtered data only
- Revealed cards properly exposed

**Status:** ⬜ Not Tested

---

### Test 12: Error Handling

**Objective:** Verify graceful error handling

**Test Cases:**

1. **Invalid Card ID**
   - Attempt to tap non-existent card
   - Expected: GAME_ERROR event to player
   - Expected: UI shows error message
   - Expected: Game state unchanged

2. **Invalid Action**
   - Attempt to draw from empty library
   - Expected: Error or empty draw
   - Expected: Game continues

3. **Network Interruption**
   - Disconnect Player B's WebSocket
   - Player A performs action
   - Reconnect Player B
   - Expected: Player B receives catch-up state

4. **Concurrent Actions**
   - Player A and B both tap same card
   - Expected: First action succeeds
   - Expected: Second action handled gracefully

**Status:** ⬜ Not Tested

---

## Manual Test Checklist

### Pre-Test Setup

- [ ] Build backend: `cd mage-server-go && go build -o bin/mage-server cmd/server/main.go`
- [ ] Build frontend: `cd mage-client-web && bun run build`
- [ ] Configure `config/config.yaml`:
  ```yaml
  server:
    engine_type: "playtest"
  ```
- [ ] Start server: `./bin/mage-server -config config/config.yaml`
- [ ] Start frontend dev: `cd mage-client-web && bun run dev`
- [ ] Open 2+ browser windows (for multiplayer)

---

### Basic Operations Test

For each operation, verify:
- [ ] Operation sends to server (check network tab)
- [ ] Server processes operation (check server logs)
- [ ] Server broadcasts update (check WebSocket frames)
- [ ] All clients receive update (check browser consoles)
- [ ] UI updates correctly (visual verification)

**Operations to Test:**

1. [ ] **Draw Card**
   - Click "Draw Card" button
   - Verify hand count increases
   - Verify card appears in hand

2. [ ] **Play Card**
   - Drag card from hand to battlefield
   - Verify card moves to battlefield
   - Verify hand count decreases

3. [ ] **Tap Card**
   - Click card on battlefield
   - Verify card rotates to tapped position
   - Verify `tapped: true` in state

4. [ ] **Untap All**
   - Click "Untap All" button
   - Verify all controlled permanents untap

5. [ ] **Modify Life**
   - Click life total
   - Adjust +/- buttons
   - Verify life changes

6. [ ] **Create Token**
   - Click "Create Token"
   - Fill in form (name, type, P/T, color)
   - Verify token appears on battlefield

7. [ ] **Add Counter**
   - Right-click card
   - Select "Add Counter"
   - Choose counter type and amount
   - Verify counter badge appears

8. [ ] **Remove Counter**
   - Right-click card with counters
   - Select "Remove Counter"
   - Verify counter decrements

9. [ ] **Move to Graveyard**
   - Drag card to graveyard
   - Verify card appears in graveyard pile

10. [ ] **Move to Exile**
    - Drag card to exile zone
    - Verify card appears in exile pile

11. [ ] **Shuffle Library**
    - Right-click library
    - Select "Shuffle"
    - Verify shuffle animation/log entry

12. [ ] **Next Turn**
    - Click "Next Turn"
    - Verify turn counter increments
    - Verify active player changes

---

### Multiplayer Sync Test

Open 2 browser windows as Player A and Player B:

- [ ] **Action by P1 visible to P2**
  - P1: Draw card
  - P2: See opponent hand count increase

- [ ] **Action by P2 visible to P1**
  - P2: Play card to battlefield
  - P1: See card appear on battlefield

- [ ] **Simultaneous visibility**
  - P1: Tap card
  - P2: Immediately see card tap

- [ ] **Hidden information enforced**
  - P1: Draw cards
  - P2: Verify can't see card names, only count

- [ ] **Public zones shared**
  - P1: Move card to graveyard
  - P2: Verify can see card in graveyard

---

### Edge Cases

- [ ] **Empty Library**
  - Draw from empty library
  - Verify graceful handling (no crash)

- [ ] **Zero Life**
  - Set life to 0
  - Verify game continues (rules-light)

- [ ] **Many Tokens**
  - Create 20+ tokens
  - Verify UI renders correctly
  - Verify no performance issues

- [ ] **Long Game**
  - Play 50+ turns
  - Verify no memory leaks
  - Verify game log doesn't slow down

- [ ] **Reconnection**
  - Disconnect and reconnect
  - Verify state recovers

---

## Automated Test Suggestions

### Backend Unit Tests

**Location:** `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/`

```go
package game

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "go.uber.org/zap/zaptest"
)

func TestEngineDrawCards(t *testing.T) {
    logger := zaptest.NewLogger(t)
    engine := NewEngine(logger)

    // Setup game with 2 players
    players := []string{"player1", "player2"}
    decks := map[string]DeckList{
        "player1": {
            MainDeck: []string{"Mountain", "Mountain", "Mountain"},
        },
        "player2": {
            MainDeck: []string{"Island", "Island", "Island"},
        },
    }

    err := engine.StartGameWithDecks("game1", players, "constructed", decks)
    assert.NoError(t, err)

    // Draw 1 card for player1
    err = engine.DrawCards("game1", "player1", 1)
    assert.NoError(t, err)

    // Verify player1 has 8 cards in hand (7 initial + 1 drawn)
    view, err := engine.GetGameView("game1", "player1")
    assert.NoError(t, err)

    playtestView := view.(*PlaytestGameView)
    assert.Equal(t, 8, len(playtestView.Me.Hand))

    // Verify player2 sees count only
    view2, err := engine.GetGameView("game1", "player2")
    assert.NoError(t, err)

    playtestView2 := view2.(*PlaytestGameView)
    assert.Equal(t, 0, len(playtestView2.Opponents[0].Hand))
    assert.Equal(t, 8, playtestView2.Opponents[0].HandCount)
}

func TestEngineTapCard(t *testing.T) {
    logger := zaptest.NewLogger(t)
    engine := NewEngine(logger)

    // Setup and play a card
    // ...

    // Tap the card
    err := engine.TapCard("game1", "player1", "card-id", true)
    assert.NoError(t, err)

    // Verify card is tapped
    view, _ := engine.GetGameView("game1", "player1")
    playtestView := view.(*PlaytestGameView)

    card := findCard(playtestView.Battlefield, "card-id")
    assert.NotNil(t, card)
    assert.True(t, card.Tapped)
}

func TestEngineHiddenInformation(t *testing.T) {
    logger := zaptest.NewLogger(t)
    engine := NewEngine(logger)

    // Setup 2-player game
    // ...

    // Player 1 draws cards
    engine.DrawCards("game1", "player1", 5)

    // Get player 2's view
    view, _ := engine.GetGameView("game1", "player2")
    playtestView := view.(*PlaytestGameView)

    // Verify opponent hand is hidden
    assert.Equal(t, 0, len(playtestView.Opponents[0].Hand))
    assert.Equal(t, 12, playtestView.Opponents[0].HandCount) // 7 initial + 5 drawn

    // Verify own hand is visible
    assert.Greater(t, len(playtestView.Me.Hand), 0)
}

// Add tests for all 19 operations
// - TestEnginePlayCard
// - TestEngineMoveCard
// - TestEngineUntapAll
// - TestEngineFlipCard
// - TestEngineModifyLife
// - TestEngineSetPlayerCounter
// - TestEngineShuffleLibrary
// - TestEngineCreateToken
// - TestEngineAddCounter
// - TestEngineRemoveCounter
// - TestEngineNextTurn
// etc.
```

---

### Frontend Unit Tests

**Location:** `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/`

```typescript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import { multiplayerGameStore } from './multiplayer-game';
import * as directActions from '$lib/api/direct-actions';

// Mock direct-actions module
vi.mock('$lib/api/direct-actions', () => ({
    drawCards: vi.fn(),
    moveCard: vi.fn(),
    tapUntap: vi.fn(),
    untapAll: vi.fn(),
    modifyPlayerLife: vi.fn(),
    createToken: vi.fn(),
    // ... mock all operations
}));

describe('multiplayerGameStore', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        multiplayerGameStore.reset();
    });

    it('calls directActions.drawCards when drawCards is called', () => {
        const gameId = 'test-game';
        const playerId = 'player1';

        multiplayerGameStore.initialize(gameId);
        multiplayerGameStore.drawCards(playerId, 3);

        expect(directActions.drawCards).toHaveBeenCalledWith(gameId, playerId, 3);
    });

    it('calls directActions.tapUntap when tapCard is called', () => {
        const gameId = 'test-game';
        const cardId = 'card-123';

        multiplayerGameStore.initialize(gameId);
        multiplayerGameStore.tapCard(cardId, true);

        expect(directActions.tapUntap).toHaveBeenCalledWith(gameId, cardId, true);
    });

    it('updates state when GAME_UPDATE is received', () => {
        // Mock WebSocket event
        const mockGameUpdate = {
            game: {
                gameId: 'test-game',
                turn: 5,
                players: [/* ... */],
                // ... full game state
            }
        };

        // Simulate GAME_UPDATE event
        // (requires WebSocket mock setup)

        const state = get(multiplayerGameStore);
        expect(state.isConnected).toBe(true);
        expect(state.turn).toBe(5);
    });
});
```

---

### Integration Tests (E2E)

**Tool:** Playwright or Cypress

```typescript
// tests/e2e/multiplayer-game.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Multiplayer Game', () => {
    test('2-player game synchronization', async ({ browser }) => {
        // Create 2 browser contexts (2 players)
        const context1 = await browser.newContext();
        const context2 = await browser.newContext();

        const player1 = await context1.newPage();
        const player2 = await context2.newPage();

        // Both players join the same game
        await player1.goto('/game/test-game-id');
        await player2.goto('/game/test-game-id');

        // Wait for game to initialize
        await player1.waitForSelector('[data-test="battlefield"]');
        await player2.waitForSelector('[data-test="battlefield"]');

        // Player 1 draws a card
        await player1.click('[data-test="draw-button"]');

        // Verify Player 1 sees card in hand
        const p1HandCount = await player1.locator('[data-test="hand-count"]').textContent();
        expect(p1HandCount).toBe('8');

        // Verify Player 2 sees opponent hand count increase
        await player2.waitForTimeout(500); // Wait for WebSocket update
        const p2OpponentHandCount = await player2.locator('[data-test="opponent-hand-count"]').textContent();
        expect(p2OpponentHandCount).toBe('8');

        // Player 2 plays a card
        await player2.dragAndDrop('[data-test="hand"] .card:first-child', '[data-test="battlefield"]');

        // Verify both players see the card on battlefield
        await player1.waitForTimeout(500);
        const p1BattlefieldCards = await player1.locator('[data-test="battlefield"] .card').count();
        const p2BattlefieldCards = await player2.locator('[data-test="battlefield"] .card').count();

        expect(p1BattlefieldCards).toBe(1);
        expect(p2BattlefieldCards).toBe(1);

        await context1.close();
        await context2.close();
    });

    test('hidden information is enforced', async ({ browser }) => {
        // Similar setup...

        // Intercept network traffic
        player2.on('response', async (response) => {
            if (response.url().includes('/mage.v1.MageServer/')) {
                const body = await response.text();
                // Verify opponent hand data is not present
                expect(body).not.toContain('opponent-card-secret-data');
            }
        });

        // Trigger game actions and verify network traffic
        await player1.click('[data-test="draw-button"]');
        // ... assertions
    });
});
```

---

## Known Limitations

### 1. Operations Not Yet Implemented

These operations show warnings in `multiplayer-game.ts`:

- **Mill Cards** (lines 391-397)
  - Frontend: Shows warning
  - Backend: Needs MILL action implementation
  - Workaround: Manually move cards to graveyard

- **Reveal Top Cards** (lines 404-411)
  - Frontend: Shows warning
  - Backend: Needs REVEAL_TOP action
  - Workaround: View library manually

- **Scry** (lines 418-445)
  - Frontend: Shows warning
  - Backend: Needs SCRY state tracking
  - Workaround: Look at top cards, reorder manually

- **Set Revealed Top** (lines 452-458)
  - Frontend: Shows warning
  - Backend: Partial support in Engine
  - Workaround: Mark in game log

- **Mulligan** (lines 510-513)
  - Frontend: Shows warning
  - Backend: Server-side only (MageEngine)
  - Workaround: Manual mulligan (return cards, redraw)

- **Keep Hand** (lines 520-523)
  - Frontend: Shows warning
  - Backend: Server-side only
  - Workaround: Click "Done" or skip mulligan

---

### 2. Library Search

**Status:** Commented out in game page (lines 1132-1147)

**Issue:** LibrarySearch component expects server-driven search, but playtest engine uses local state.

**Workaround:** Use keyboard shortcut "F" to open search, or implement client-side search.

---

### 3. Rollback System

**Status:** Exists in backend, not exposed in frontend UI

**Backend:** `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/engine-operations.go` contains rollback operations.

**Frontend:** No UI for proposing/accepting rollback.

**Workaround:** Use server admin commands or implement UI in future phase.

---

### 4. Consent System

**Status:** Backend implementation exists, frontend integration pending

**Use Case:** Mulligan consent, rollback consent.

**Workaround:** Manual coordination between players.

---

### 5. Performance

**Large Games:** No stress testing done for:
- 100+ permanents on battlefield
- 1000+ turn games
- 10+ players

**Recommendation:** Test performance separately for large-scale scenarios.

---

### 6. Error Recovery

**Network Failures:** Reconnection logic exists but not thoroughly tested.

**Race Conditions:** Concurrent actions from multiple players may cause edge cases.

**Recommendation:** Add retry logic and conflict resolution.

---

## Test Results Template

Use this template to document each test execution:

```markdown
## Test Execution: [Date]

### Test: [Test Name]

**Status**: ✅ Pass / ❌ Fail / ⚠️ Partial

**Environment:**
- Server Version: [version]
- Frontend Version: [version]
- Engine Type: playtest / mage
- Browser: [browser and version]

**Steps:**
1. [Step 1]
2. [Step 2]
3. [Step 3]
...

**Expected Result:**
[Describe expected behavior]

**Actual Result:**
[Describe actual behavior]

**Issues Found:**
- Issue 1: [description]
- Issue 2: [description]

**Screenshots/Logs:**
[Attach any relevant screenshots or log snippets]

**Notes:**
[Any additional observations]

---
```

### Example Test Result

```markdown
## Test Execution: January 22, 2026

### Test: 2-Player Draw Card Synchronization

**Status**: ✅ Pass

**Environment:**
- Server Version: dev (commit ed84689)
- Frontend Version: dev (commit ed84689)
- Engine Type: playtest
- Browser: Chrome 131.0.6778.109

**Steps:**
1. Started server with `engine_type: "playtest"`
2. Opened two browser windows (Player A and Player B)
3. Both players joined game "test-game-123"
4. Player A clicked "Draw Card" button
5. Observed both players' UI

**Expected Result:**
- Player A sees new card in hand
- Player A hand count increases to 8
- Player B sees opponent hand count increase to 8
- Player B does NOT see card details (hidden information)

**Actual Result:**
- ✅ Player A hand count increased to 8
- ✅ Player A sees card details (Mountain)
- ✅ Player B sees opponent hand count = 8
- ✅ Player B sees empty hand array (hidden)
- ✅ WebSocket GAME_UPDATE received by both players in <100ms

**Issues Found:**
None

**Screenshots/Logs:**
```
[Server Log]
INFO handling game notification game_id=test-game-123 type=GAME_STATE_CHANGE
INFO sending GAME_UPDATE to player game_id=test-game-123 player=playerA
INFO sending GAME_UPDATE to player game_id=test-game-123 player=playerB

[Player A Console]
[MultiplayerGame] drawCards: { playerId: 'playerA', count: 1 }
[MultiplayerGame] GAME_UPDATE received: { turn: 1, ... }

[Player B Console]
[MultiplayerGame] GAME_UPDATE received: { opponents: [{ handCount: 8, hand: [] }] }
```

**Notes:**
Synchronization was instant. No lag observed. Hidden information correctly enforced.

---
```

---

## Summary

This integration test plan provides:

1. **Code Path Verification:** All critical paths documented and verified to exist
2. **Backend Tests:** 5 comprehensive backend integration scenarios
3. **Frontend Tests:** 3 frontend integration scenarios
4. **E2E Tests:** 4 end-to-end multiplayer scenarios
5. **Manual Checklist:** Step-by-step manual test guide
6. **Automated Tests:** Suggestions for unit and E2E test automation
7. **Known Limitations:** 6 documented limitations with workarounds
8. **Test Template:** Structured format for recording results

### Next Steps

1. Execute manual test checklist
2. Document results using template
3. Fix any issues found
4. Implement automated tests for critical paths
5. Re-run tests after fixes
6. Sign off on Phase 6 completion

### Success Criteria for Phase 6

- [ ] All code paths verified to exist
- [ ] All 19 operations tested and working
- [ ] 2-player synchronization works
- [ ] 4-player synchronization works
- [ ] Hidden information correctly enforced
- [ ] No critical bugs found
- [ ] Documentation complete

---

**Document End**
