# Implementation Plan: Multiplayer Interface TODOs

**Created**: 2026-01-24
**Source**: docs/20260124-missing-pieces-multiplayer-interface.md
**Status**: Ready for execution

---

## Phase 0: Documentation Discovery (COMPLETE)

### Gathered Documentation

All necessary APIs, patterns, and structures have been documented through comprehensive codebase analysis:

#### 0.1 Protobuf Patterns

**Sources Consulted**:

- `mage-server-go/api/proto/mage/v1/models.proto` (GameView message lines 72-103, PlayerView message lines 115-134)
- `mage-server-go/internal/server/grpc.go` (conversion functions lines 628-746)
- `mage-server-go/Makefile` (proto generation command line 49)

**Key Findings**:

- Proto regeneration: `cd mage-server-go && make proto`
- Field numbering: 1-19 core fields, 20-29 display values, 30+ feature fields
- Conversion pattern: `PlaytestGameView` → `playtestViewToProto()` → `pb.GameView`
- Helper pattern: `playtestEngineCardsToProto()` for card arrays

**Missing Fields Identified**:

- `activeControlSeat` (field 19) - NOT in proto, needed by client
- `library` - NOT transmitted in `playtestViewToProto()` despite proto having the field
- `manaPool` - Proto type exists but not populated in conversion

#### 0.2 WebSocket State Sync Patterns

**Sources Consulted**:

- `mage-client-web/src/lib/stores/multiplayer-game.ts` (GAME_UPDATE handler lines 162-179)
- `mage-client-web/src/lib/stores/game.legacy.ts` (working example lines 176-226)
- `mage-server-go/internal/server/grpc.go` (notification flow lines 115-290)
- `mage-server-go/internal/game/view.go` (buildGameView lines 68-136)
- `mage-server-go/internal/game/game_engine.go` (broadcast mechanism lines 330-346)

**Key Findings**:

- Complete data flow: `GameEngine.broadcast()` → `buildGameView()` → `NotifyGameStateChange()` → `handleGameNotification()` → `sendGameUpdateToPlayer()` → `engineViewToProto()` → WebSocket → Client `GAME_UPDATE` handler
- Working example exists in `game.legacy.ts` showing proper mapping pattern
- Current handler only updates connection flags, doesn't map GameView to state
- Hidden information filtering works correctly (Me gets full hand/library, Opponents get empty arrays)

**Required Pattern** (from game.legacy.ts):

```typescript
websocketStore.on(CallbackMethod.GAME_UPDATE, (data) => {
  const updateData = data as GameUpdateData;
  if (updateData.game) {
    const normalized = normalizeGameView(updateData.game);
    update((s) => ({
      ...s,
      gameView: normalized,
      // Clear optimistic updates
    }));
  }
});
```

#### 0.3 Action Processing Architecture

**Sources Consulted**:

- `mage-server-go/internal/game/game_engine.go` (ProcessAction switch lines 148-255)
- `mage-server-go/internal/game/actions.go` (all action handlers)
- `mage-server-go/internal/game/manager.go` (SendPlayerAction lines 291-310, ProcessGameActions lines 312-327)
- `mage-server-go/internal/server/grpc_game.go` (SendPlayerString RPC line 454, goroutine spawn line 131)
- `mage-client-web/src/lib/api/direct-actions.ts` (client command API)
- `mage-client-web/src/lib/api/game.ts` (sendPlayerString function)

**Key Findings**:

- ActionQueue consumer goroutine **IS running** (spawned at line 131 in grpc_game.go)
- ProcessAction uses **structured data**, NOT string parsing
- Current flow: `SendPlayerString(data: "TAP:card-123")` → `SendPlayerAction(type: "SEND_STRING", data: "TAP:card-123")` → ProcessAction **NO HANDLER for "SEND_STRING"** → error
- All existing handlers use structured data from `map[string]interface{}`
- 18 action types already implemented: DRAW, PLAY, MOVE, TAP, FLIP, MODIFY_LIFE, SET_COUNTER, SHUFFLE, CREATE_TOKEN, ADD_COUNTER, REMOVE_COUNTER, SET_CARD_COUNTER, MILL, SCRY, SET_REVEALED_TOP, NEXT_TURN, MULLIGAN, KEEP_HAND

**Critical Gap**: Client sends string commands but server expects structured actions

**Missing Operations** (from TODO #3): These handlers **ALREADY EXIST** in actions.go:

- ✅ MILL - `MillCards()` at lines 923-957
- ✅ SCRY - `ScryCards()` at lines 1016-1053
- ✅ MULLIGAN - `Mulligan()` at lines 1095-1146
- ✅ KEEP_HAND - `KeepHand()` at lines 1151-1163
- ✅ SET_REVEALED_TOP - `SetRevealedTop()` at lines 1058-1071

**Action NOT needed**: Server operations exist; only need command parser.

#### 0.4 Card Data Structures

**Sources Consulted**:

- `mage-server-go/internal/game/game_state.go` (Card struct lines 48-87)
- `mage-server-go/api/proto/mage/v1/models.proto` (CardView message lines 136-167)
- `mage-client-web/src/lib/types/game.ts` (GameCard interface lines 5-48)
- `mage-server-go/internal/game/game_engine.go` (card creation lines 96-106)
- `mage-client-web/src/lib/utils/scryfall.ts` (image URL generation)
- `mage-server-go/internal/repository/cards.go` (Scryfall database schema)

**Key Findings**:

- Card struct has all metadata fields (ManaCost, Type, Power, Toughness, RulesText, etc.)
- Cards created with **name only** (lines 96-106 in game_engine.go)
- Database has Scryfall data including image URIs, but **NOT integrated with game engine**
- Client generates Scryfall URLs on-demand: `getScryfallImageUrl(cardName, version)`
- Service worker caches Scryfall images for 7 days
- No imageUrl field in proto or server Card struct (client-only concern)

**Current Pattern** (client-side enrichment):

```typescript
const effectiveImageUrl = $derived(
  imageUrl || getScryfallImageUrl(cardName, getScryfallVersionForSize(size)),
);
```

**Enrichment Options**:

1. **Server-side** (wanted): Fetch from CardRepository during StartGameWithDecks
2. **Client-side** (current): Generate URLs on-demand, fetch metadata as needed

---

## Phase 1: Core State Synchronization (UNBLOCKS UI)

**Goal**: Make the multiplayer game page load and display current game state correctly.

**Dependencies**: None - uses existing proto structure and WebSocket infrastructure.

**Testing After Phase**: Page loads, shows players with life totals, displays cards in zones, no console errors.

---

### Task 1.1: Add activeControlSeat Field to Proto

**What to Implement**:
Add the `active_control_seat` field to the `GameView` proto message so each player knows which view is theirs.

**Documentation Reference**:

- Proto pattern: `mage-server-go/api/proto/mage/v1/models.proto` lines 72-103
- Field numbering convention: Use field 19 (display values range 20-29, so 19 fits as identity field)
- Example: `LibrarySearchView pending_library_search = 30;` added in previous feature

**Files to Modify**:

1. `mage-server-go/api/proto/mage/v1/models.proto`

**Implementation Steps**:

1. Edit the proto file and add field to GameView message (after field 18):

```protobuf
message GameView {
  string game_id = 1;
  string state = 2;
  repeated PlayerView players = 3;
  // ... existing fields 4-18 ...

  // Field 19: Which player perspective this view is for (their own player ID)
  string active_control_seat = 19;

  // Pre-computed display values (server source of truth)
  string active_player_name = 20;
  // ... rest of fields 20-30 ...
}
```

2. Regenerate proto files:

```bash
cd mage-server-go && make proto
```

3. Verify TypeScript types are regenerated (should happen automatically if build is running)

**Verification**:

- ✅ Proto generation succeeds without errors
- ✅ TypeScript client types include `activeControlSeat?: string` in GameView
- ✅ Go code compiles without errors

**Anti-Patterns to Avoid**:

- ❌ Don't use field numbers already in use
- ❌ Don't skip regeneration step (proto changes won't take effect)
- ❌ Don't modify generated Go or TypeScript files manually

---

### Task 1.2: Set activeControlSeat in Server Conversion

**What to Implement**:
Populate the new `active_control_seat` field when converting PlaytestGameView to proto, using the `playerID` parameter (which is the viewing player's ID).

**Documentation Reference**:

- Conversion function: `mage-server-go/internal/server/grpc.go` lines 647-697
- Pattern: Set field in `&pb.GameView{...}` struct initialization
- Source data: `playerID` parameter is the viewing player's ID

**Files to Modify**:

1. `mage-server-go/internal/server/grpc.go`

**Implementation Steps**:

1. Locate `playtestViewToProto` function (line 647)

2. Add `ActiveControlSeat` field to the `pb.GameView` initialization (around line 650):

```go
func (s *mageServer) playtestViewToProto(data *game.PlaytestGameView, playerID string) *pb.GameView {
	view := &pb.GameView{
		GameId:            data.GameID,
		State:             "IN_PROGRESS",
		Phase:             "",
		Step:              "",
		Turn:              int32(data.Turn),
		ActivePlayerId:    data.ActivePlayerID,
		ActiveControlSeat: playerID, // ADD THIS LINE - viewing player's ID
		Battlefield:       playtestEngineCardsToProto(data.Battlefield),
		Stack:             playtestEngineCardsToProto(data.Stack),
		Exile:             playtestEngineCardsToProto(data.Exile),
		Command:           playtestEngineCardsToProto(data.Command),
	}
	// ... rest of function
}
```

**Verification**:

- ✅ Server compiles without errors
- ✅ Start a test game and check WebSocket GAME_UPDATE payload includes `activeControlSeat`
- ✅ `activeControlSeat` matches the player ID receiving the update

**Anti-Patterns to Avoid**:

- ❌ Don't use `data.ActiveControlSeat` (that's the engine's internal seat, not the viewing player)
- ❌ Don't hardcode a player ID
- ❌ Don't set it to `data.ActivePlayerID` (that's whose turn it is, not who's viewing)

---

### Task 1.3: Create Proto GameView to State Mapping Function

**What to Implement**:
Create a conversion function that maps the proto `GameView` structure to the client's `MultiplayerGameState` structure.

**Documentation Reference**:

- Target state structure: `mage-client-web/src/lib/stores/multiplayer-game.ts` lines 43-84
- Proto GameView structure: Generated TypeScript types from proto
- Working example: `mage-client-web/src/lib/stores/game.legacy.ts` lines 176-226 (normalizeGameView pattern)
- PlayerView array needs mapping to PlaytestPlayer objects

**Files to Modify**:

1. `mage-client-web/src/lib/stores/multiplayer-game.ts`

**Implementation Steps**:

1. Add helper function to convert proto PlayerView to PlaytestPlayer (add after interface definitions, before store creation):

```typescript
/**
 * Converts a proto PlayerView to a PlaytestPlayer object.
 * Handles differences between proto structure and client expectations.
 */
function convertPlayerViewToPlaytestPlayer(pv: PlayerView): PlaytestPlayer {
  return {
    playerId: pv.playerId || "",
    name: pv.name || "",
    life: Number(pv.life) || 0,
    poison: Number(pv.poison) || 0,
    energy: Number(pv.energy) || 0,
    libraryCount: Number(pv.libraryCount) || 0,
    handCount: Number(pv.handCount) || 0,
    hand: pv.hand || [],
    library: pv.library || [], // Will be populated for viewing player, empty for opponents
    graveyard: pv.graveyard || [],
    manaPool: pv.manaPool || {
      white: 0,
      blue: 0,
      black: 0,
      red: 0,
      green: 0,
      colorless: 0,
    },
    keptHand: pv.keptHand || false,
    mulliganCount: 0, // Not in proto yet, default to 0
    revealedTopCard: false, // Not in proto yet, default to false
  };
}

/**
 * Normalizes a proto GameView by ensuring all array fields are initialized.
 * Proto messages omit empty arrays, so we need to provide defaults.
 */
function normalizeProtoGameView(game: GameView): GameView {
  return {
    ...game,
    players: game.players || [],
    battlefield: game.battlefield || [],
    stack: game.stack || [],
    exile: game.exile || [],
    command: game.command || [],
  };
}

/**
 * Maps a normalized proto GameView to MultiplayerGameState structure.
 * Converts proto PlayerView[] to PlaytestPlayer[] and maps all fields.
 */
function mapProtoGameViewToState(
  game: GameView,
): Partial<MultiplayerGameState> {
  return {
    gameId: game.gameId || "",
    turn: Number(game.turn) || 0,
    activePlayerId: game.activePlayerId || "",
    activeControlSeat: game.activeControlSeat || "", // From Phase 1 Task 1.1
    battlefield: game.battlefield || [],
    exile: game.exile || [],
    stack: game.stack || [],
    command: game.command || [],
    players: (game.players || []).map(convertPlayerViewToPlaytestPlayer),
    mulliganType: "london", // Default, could be in proto later
    freeMulligans: 0, // Default, could be in proto later
    log: [], // Not in current proto, default to empty
  };
}
```

2. Add these helper functions to the file (they'll be used in Task 1.4)

**Verification**:

- ✅ TypeScript compiles without errors
- ✅ All functions are properly typed
- ✅ No linting errors

**Anti-Patterns to Avoid**:

- ❌ Don't skip normalization step (proto omits empty arrays)
- ❌ Don't assume all proto fields are present (use `||` defaults)
- ❌ Don't mutate input objects (use spread operator)
- ❌ Don't hardcode player IDs or array indices

---

### Task 1.4: Update GAME_UPDATE Handler to Map State

**What to Implement**:
Replace the current placeholder GAME_UPDATE handler with one that actually maps the proto GameView to the client state using the functions from Task 1.3.

**Documentation Reference**:

- Current handler: `mage-client-web/src/lib/stores/multiplayer-game.ts` lines 162-179
- Working example: `mage-client-web/src/lib/stores/game.legacy.ts` lines 176-226
- Helper functions from Task 1.3

**Files to Modify**:

1. `mage-client-web/src/lib/stores/multiplayer-game.ts`

**Implementation Steps**:

1. Locate the GAME_UPDATE handler (line 162)

2. Replace the current implementation:

```typescript
// GAME_UPDATE - State update (from game.legacy.ts lines 176-226)
unsubscribers.push(
  websocketStore.on(CallbackMethod.GAME_UPDATE, (data) => {
    const updateData = data as GameUpdateData;
    console.log("[MultiplayerGame] GAME_UPDATE received:", updateData);

    if (updateData.game) {
      // Normalize proto GameView (handle empty arrays)
      const normalized = normalizeProtoGameView(updateData.game);

      // Map proto structure to our state structure
      const mappedState = mapProtoGameViewToState(normalized);

      // Apply server state to store
      update((state) => ({
        ...state,
        ...mappedState,
        isConnected: true,
        isInitialized: true,
        // Clear pending actions (server is source of truth)
        pendingActions: [],
      }));
    }
  }),
);
```

**Verification**:

- ✅ TypeScript compiles without errors
- ✅ Start a test multiplayer game
- ✅ GAME_UPDATE arrives via WebSocket
- ✅ Check browser devtools: state.players array is populated
- ✅ Check browser devtools: state.battlefield, state.turn, etc. have correct values
- ✅ Console log shows mapped state

**Anti-Patterns to Avoid**:

- ❌ Don't skip normalization (proto omits empty arrays)
- ❌ Don't clear entire state (use spread to preserve connection flags)
- ❌ Don't mutate state directly (use update function)
- ❌ Don't ignore errors (add try-catch if mapping could fail)

---

### Task 1.5: Apply Initial State from fetchGameView

**What to Implement**:
Update the `loadGameState` function to apply the fetched GameView to the store state, not just set connection flags.

**Documentation Reference**:

- Current implementation: `mage-client-web/src/lib/stores/multiplayer-game.ts` lines 219-235
- Same mapping functions as Task 1.3 and 1.4
- fetchGameView API: `mage-client-web/src/lib/api/game.ts`

**Files to Modify**:

1. `mage-client-web/src/lib/stores/multiplayer-game.ts`

**Implementation Steps**:

1. Locate the `loadGameState` function (line 213)

2. Update to apply fetched GameView using the same mapping pattern:

```typescript
async function loadGameState(gameId: string, playerId: string): Promise<void> {
  try {
    // Fetch initial game view from server
    const gameView = await fetchGameView(gameId, playerId);

    if (gameView) {
      // Normalize and map the fetched GameView
      const normalized = normalizeProtoGameView(gameView);
      const mappedState = mapProtoGameViewToState(normalized);

      // Apply initial state
      update((state) => ({
        ...state,
        ...mappedState,
        isConnected: true,
        isInitialized: true,
        pendingActions: [],
      }));
    } else {
      // GameView fetch returned null/undefined
      update((state) => ({
        ...state,
        isConnected: true,
        isInitialized: false, // No game state available
      }));
    }
  } catch (error) {
    console.error("[MultiplayerGame] Failed to load game state:", error);
    update((state) => ({
      ...state,
      isConnected: false,
      isInitialized: false,
    }));
  }
}
```

**Verification**:

- ✅ TypeScript compiles without errors
- ✅ Load game page directly (not via WebSocket join flow)
- ✅ Check browser devtools: state is populated from fetchGameView
- ✅ Players, battlefield, turn counter all display correctly
- ✅ No errors in console

**Anti-Patterns to Avoid**:

- ❌ Don't ignore null/undefined gameView (handle gracefully)
- ❌ Don't skip error handling (network could fail)
- ❌ Don't assume gameView structure (normalize first)
- ❌ Don't duplicate mapping logic (reuse helper functions)

---

### Task 1.6: Add Library Field to Proto Conversion

**What to Implement**:
Update the `playtestViewToProto` conversion to include the `Library` field for the viewing player (currently omitted despite proto having the field).

**Documentation Reference**:

- Current conversion: `mage-server-go/internal/server/grpc.go` lines 666-677
- Proto definition: `mage-server-go/api/proto/mage/v1/models.proto` line 123 (library field exists)
- Source data: `data.Me.Library` contains the viewing player's library
- Pattern: Same as Hand field (line 675)

**Files to Modify**:

1. `mage-server-go/internal/server/grpc.go`

**Implementation Steps**:

1. Locate the viewing player conversion in `playtestViewToProto` (around line 666)

2. Add the Library field mapping (after Hand field):

```go
// Add the viewing player first
if data.Me != nil {
	players = append(players, &pb.PlayerView{
		PlayerId:     data.Me.PlayerID,
		Name:         data.Me.Name,
		Life:         int32(data.Me.Life),
		Poison:       int32(data.Me.Poison),
		Energy:       int32(data.Me.Energy),
		LibraryCount: int32(data.Me.LibraryCount),
		HandCount:    int32(data.Me.HandCount),
		Hand:         playtestEngineCardsToProto(data.Me.Hand),
		Library:      playtestEngineCardsToProto(data.Me.Library), // ADD THIS LINE
		Graveyard:    playtestEngineCardsToProto(data.Me.Graveyard),
	})
}
```

3. Verify opponent conversion does NOT include Library (should remain empty, line 683):

```go
// Add opponents
for _, opponent := range data.Opponents {
	players = append(players, &pb.PlayerView{
		PlayerId:     opponent.PlayerID,
		Name:         opponent.Name,
		Life:         int32(opponent.Life),
		Poison:       int32(opponent.Poison),
		Energy:       int32(opponent.Energy),
		LibraryCount: int32(opponent.LibraryCount),
		HandCount:    int32(opponent.HandCount),
		Hand:         playtestEngineCardsToProto(opponent.Hand), // Empty for opponents
		// NO Library field - hidden information
		Graveyard:    playtestEngineCardsToProto(opponent.Graveyard),
	})
}
```

**Verification**:

- ✅ Server compiles without errors
- ✅ Start a test game and draw cards
- ✅ Check WebSocket GAME_UPDATE payload: viewing player's PlayerView has `library` array
- ✅ Opponent's PlayerView does NOT have `library` field (or it's empty)
- ✅ Client state shows library array for local player

**Anti-Patterns to Avoid**:

- ❌ Don't add Library to opponent conversion (breaks hidden information)
- ❌ Don't forget the proto field exists (it's already defined, just not populated)
- ❌ Don't use wrong helper function (use `playtestEngineCardsToProto`, not a generic one)

---

### Task 1.7: Verify Derived Stores Work

**What to Implement**:
Test that the existing derived stores (`multiplayerPlayers`, `multiplayerLocalPlayer`, `multiplayerOpponents`) work correctly once the players array is populated.

**Documentation Reference**:

- Derived stores: `mage-client-web/src/lib/stores/multiplayer-game.ts` lines 673-699
- Dependencies: Rely on `players` array and `activeControlSeat` being set correctly

**Files to Verify** (no changes needed):

1. `mage-client-web/src/lib/stores/multiplayer-game.ts` (lines 673-699)

**Verification Steps**:

1. After completing Tasks 1.1-1.6, start a multiplayer game with 2+ players

2. Open browser devtools console and check derived store values:

```javascript
// In console:
// Check that players array is populated
$multiplayerGame.players; // Should show PlaytestPlayer[]

// Check derived stores
$multiplayerPlayers; // Should show all players with derived fields (isLocal, isActive)
$multiplayerLocalPlayer; // Should show viewing player's data
$multiplayerOpponents; // Should show other players
```

3. Verify derived fields:

- `isLocal` is true for one player (the viewing player)
- `isActive` is true for the active player (whose turn it is)
- `multiplayerLocalPlayer` matches `activeControlSeat`
- `multiplayerOpponents` excludes the local player

4. Test UI components that use these stores:

- Player life totals display
- Hand cards show for local player, card count for opponents
- Active player indicator highlights correctly

**Expected Behavior**:

- ✅ `multiplayerPlayers` derives correctly from `players` array
- ✅ `multiplayerLocalPlayer` finds player matching `activeControlSeat`
- ✅ `multiplayerOpponents` filters out local player
- ✅ UI components display player data without errors

**Anti-Patterns to Avoid**:

- ❌ Don't modify derived store logic unless verification fails
- ❌ Don't skip manual testing in browser
- ❌ Don't assume it works without checking (derived stores can have subtle bugs)

---

### Phase 1 Completion Checklist

**Before moving to Phase 2, verify**:

- [ ] Proto regenerated successfully with `activeControlSeat` field
- [ ] Server sets `activeControlSeat` to viewing player's ID
- [ ] GAME_UPDATE handler maps proto GameView to state
- [ ] fetchGameView applies initial state correctly
- [ ] Library field transmitted for viewing player only
- [ ] Derived stores (`multiplayerPlayers`, `multiplayerLocalPlayer`, `multiplayerOpponents`) work
- [ ] Game page loads and displays:
  - [ ] All players with names and life totals
  - [ ] Cards in battlefield, hand, graveyard
  - [ ] Turn counter
  - [ ] Active player indicator
- [ ] No console errors or TypeScript compilation errors
- [ ] WebSocket GAME_UPDATE payloads contain all expected data

---

## Phase 2: Action Processing (ENABLES GAMEPLAY)

**Goal**: Enable players to perform game actions (tap cards, modify life, draw cards, etc.) through the UI.

**Dependencies**: Phase 1 complete (state sync working).

**Testing After Phase**: Can tap/untap cards, modify life totals, draw cards, move cards between zones, and see updates in real-time for all players.

---

### Task 2.1: Implement String Command Parser

**What to Implement**:
Add a `ParseAndExecuteStringCommand` method to the GameEngine that converts string commands (e.g., "TAP:card-123") into structured action calls. Add a case in ProcessAction to route "SEND_STRING" actions to this parser.

**Documentation Reference**:

- ProcessAction switch: `mage-server-go/internal/game/game_engine.go` lines 148-255
- Existing action handlers: `mage-server-go/internal/game/actions.go` (all handlers)
- String command format: `mage-client-web/src/lib/api/direct-actions.ts` (all functions)
- SendPlayerString flow: `mage-server-go/internal/server/grpc_game.go` line 454

**String Command Formats** (from direct-actions.ts):

```
TAP:cardId
UNTAP:cardId
MOVE:cardId:targetZone
FLIP:cardId:true|false
DRAW:playerId:count
MODIFY_LIFE:targetPlayerId:delta
SET_COUNTER:targetPlayerId:counterType:value
SHUFFLE:playerId
CREATE_TOKEN:name:types:power:toughness:color
ADD_COUNTER:cardId:counterName:amount
REMOVE_COUNTER:cardId:counterName:amount
SET_CARD_COUNTER:cardId:counterName:amount
MILL:playerId:count
SCRY:playerId:scryCount
REVEAL_TOP:playerId:true|false
NEXT_TURN:playerId
MULLIGAN:playerId
KEEP_HAND:playerId
```

**Files to Modify**:

1. `mage-server-go/internal/game/game_engine.go` (add case to ProcessAction, add new method)
2. `mage-server-go/internal/game/actions.go` (or create `command_parser.go` for the parser)

**Implementation Steps**:

1. Add a case for "SEND_STRING" in ProcessAction (in game_engine.go, add before the `default` case):

```go
case "SEND_STRING":
	// Parse string command from direct-actions.ts
	cmdStr, ok := action.Data.(string)
	if !ok {
		return fmt.Errorf("SEND_STRING action data must be a string")
	}
	return e.ParseAndExecuteStringCommand(gameID, playerID, cmdStr)
```

2. Create `ParseAndExecuteStringCommand` method (add to game_engine.go or create command_parser.go):

```go
// ParseAndExecuteStringCommand parses a colon-delimited command string and executes the corresponding action.
// Command format: "ACTION:param1:param2:..."
// Used by direct-actions.ts client API for rules-light gameplay.
func (e *GameEngine) ParseAndExecuteStringCommand(gameID, playerID, command string) error {
	if command == "" {
		return fmt.Errorf("empty command string")
	}

	parts := strings.Split(command, ":")
	if len(parts) == 0 {
		return fmt.Errorf("invalid command format")
	}

	actionType := parts[0]

	switch actionType {
	case "TAP":
		if len(parts) < 2 {
			return fmt.Errorf("TAP command requires cardId parameter")
		}
		cardID := parts[1]
		return e.TapCard(gameID, playerID, cardID, true)

	case "UNTAP":
		if len(parts) < 2 {
			return fmt.Errorf("UNTAP command requires cardId parameter")
		}
		cardID := parts[1]
		return e.TapCard(gameID, playerID, cardID, false)

	case "MOVE":
		if len(parts) < 3 {
			return fmt.Errorf("MOVE command requires cardId and targetZone parameters")
		}
		cardID := parts[1]
		targetZone := parts[2]
		return e.MoveCard(gameID, playerID, cardID, targetZone)

	case "FLIP":
		if len(parts) < 3 {
			return fmt.Errorf("FLIP command requires cardId and faceDown parameters")
		}
		cardID := parts[1]
		faceDown, _ := strconv.ParseBool(parts[2])
		return e.FlipCard(gameID, playerID, cardID, faceDown)

	case "DRAW":
		if len(parts) < 3 {
			return fmt.Errorf("DRAW command requires playerId and count parameters")
		}
		targetPlayerID := parts[1]
		count, _ := strconv.Atoi(parts[2])
		return e.DrawCards(gameID, targetPlayerID, count)

	case "MODIFY_LIFE":
		if len(parts) < 3 {
			return fmt.Errorf("MODIFY_LIFE command requires targetPlayerId and delta parameters")
		}
		targetPlayerID := parts[1]
		delta, _ := strconv.Atoi(parts[2])
		return e.ModifyLife(gameID, targetPlayerID, delta)

	case "SET_COUNTER":
		if len(parts) < 4 {
			return fmt.Errorf("SET_COUNTER command requires targetPlayerId, counterType, and value parameters")
		}
		targetPlayerID := parts[1]
		counterType := parts[2]
		value, _ := strconv.Atoi(parts[3])
		return e.SetPlayerCounter(gameID, targetPlayerID, counterType, value)

	case "SHUFFLE":
		if len(parts) < 2 {
			return fmt.Errorf("SHUFFLE command requires playerId parameter")
		}
		targetPlayerID := parts[1]
		return e.ShuffleLibrary(gameID, targetPlayerID)

	case "CREATE_TOKEN":
		if len(parts) < 6 {
			return fmt.Errorf("CREATE_TOKEN command requires name, types, power, toughness, and color parameters")
		}
		name := parts[1]
		types := parts[2]
		power := parts[3]
		toughness := parts[4]
		color := parts[5]
		return e.CreateToken(gameID, playerID, name, types, power, toughness, color)

	case "ADD_COUNTER":
		if len(parts) < 4 {
			return fmt.Errorf("ADD_COUNTER command requires cardId, counterName, and amount parameters")
		}
		cardID := parts[1]
		counterName := parts[2]
		amount, _ := strconv.Atoi(parts[3])
		return e.AddCounter(gameID, playerID, cardID, counterName, amount)

	case "REMOVE_COUNTER":
		if len(parts) < 4 {
			return fmt.Errorf("REMOVE_COUNTER command requires cardId, counterName, and amount parameters")
		}
		cardID := parts[1]
		counterName := parts[2]
		amount, _ := strconv.Atoi(parts[3])
		return e.RemoveCounter(gameID, playerID, cardID, counterName, amount)

	case "SET_CARD_COUNTER":
		if len(parts) < 4 {
			return fmt.Errorf("SET_CARD_COUNTER command requires cardId, counterName, and amount parameters")
		}
		cardID := parts[1]
		counterName := parts[2]
		amount, _ := strconv.Atoi(parts[3])
		return e.SetCounter(gameID, playerID, cardID, counterName, amount)

	case "MILL":
		if len(parts) < 3 {
			return fmt.Errorf("MILL command requires playerId and count parameters")
		}
		targetPlayerID := parts[1]
		count, _ := strconv.Atoi(parts[2])
		return e.MillCards(gameID, targetPlayerID, count)

	case "SCRY":
		if len(parts) < 3 {
			return fmt.Errorf("SCRY command requires playerId and scryCount parameters")
		}
		targetPlayerID := parts[1]
		scryCount, _ := strconv.Atoi(parts[2])
		// For simple string command, assume keep all on top (full scry UI would need more params)
		return e.ScryCards(gameID, targetPlayerID, scryCount, []string{}, []string{})

	case "REVEAL_TOP":
		if len(parts) < 3 {
			return fmt.Errorf("REVEAL_TOP command requires playerId and revealed parameters")
		}
		targetPlayerID := parts[1]
		revealed, _ := strconv.ParseBool(parts[2])
		return e.SetRevealedTop(gameID, targetPlayerID, revealed)

	case "NEXT_TURN":
		if len(parts) < 2 {
			return fmt.Errorf("NEXT_TURN command requires playerId parameter")
		}
		targetPlayerID := parts[1]
		return e.NextTurn(gameID, targetPlayerID)

	case "MULLIGAN":
		if len(parts) < 2 {
			return fmt.Errorf("MULLIGAN command requires playerId parameter")
		}
		targetPlayerID := parts[1]
		return e.Mulligan(gameID, targetPlayerID)

	case "KEEP_HAND":
		if len(parts) < 2 {
			return fmt.Errorf("KEEP_HAND command requires playerId parameter")
		}
		targetPlayerID := parts[1]
		return e.KeepHand(gameID, targetPlayerID)

	default:
		return fmt.Errorf("unknown string command type: %s", actionType)
	}
}
```

3. Add import for `strconv` package at top of file if not already present:

```go
import (
	"strconv"
	"strings"
	// ... other imports
)
```

**Verification**:

- ✅ Server compiles without errors
- ✅ Start a test game
- ✅ Use browser console to call direct-actions API:
  ```javascript
  // Should tap a card
  await tapUntap(gameId, cardId, true);
  ```
- ✅ Check server logs: "processing action" with type "SEND_STRING"
- ✅ Check server logs: Action executes successfully
- ✅ GAME_UPDATE broadcast shows card is now tapped
- ✅ No errors in server logs or client console

**Anti-Patterns to Avoid**:

- ❌ Don't skip parameter validation (check parts length before accessing)
- ❌ Don't ignore type conversion errors (strconv.Atoi, strconv.ParseBool)
- ❌ Don't assume command format (validate before parsing)
- ❌ Don't forget error messages (help debugging)

---

### Task 2.2: Test All Direct Action Commands

**What to Implement**:
Systematically test each command from direct-actions.ts to ensure the parser routes correctly and actions execute.

**Documentation Reference**:

- Direct actions client API: `mage-client-web/src/lib/api/direct-actions.ts`
- Parser implementation from Task 2.1

**Files to Test** (no changes):

1. All functions in `mage-client-web/src/lib/api/direct-actions.ts`

**Testing Script**:

Create a test script in browser console after starting a multiplayer game:

```javascript
// Test script for direct actions
const gameId = "your-game-id"; // Replace with actual game ID
const playerId = "your-player-id"; // Replace with actual player ID

async function testDirectActions() {
  console.log("Testing direct actions...");

  // Get a card ID from battlefield
  const cards = $multiplayerGame.battlefield;
  if (cards.length === 0) {
    console.error("No cards on battlefield to test with");
    return;
  }
  const cardId = cards[0].id;

  try {
    // Test TAP
    console.log("Testing TAP...");
    await tapUntap(gameId, cardId, true);
    console.log("✅ TAP successful");

    // Test UNTAP
    console.log("Testing UNTAP...");
    await tapUntap(gameId, cardId, false);
    console.log("✅ UNTAP successful");

    // Test DRAW
    console.log("Testing DRAW...");
    await drawCards(gameId, playerId, 1);
    console.log("✅ DRAW successful");

    // Test MODIFY_LIFE
    console.log("Testing MODIFY_LIFE...");
    await modifyLife(gameId, playerId, -1);
    console.log("✅ MODIFY_LIFE successful");

    // Test SHUFFLE
    console.log("Testing SHUFFLE...");
    await shuffleLibrary(gameId, playerId);
    console.log("✅ SHUFFLE successful");

    // Test MOVE
    console.log("Testing MOVE...");
    await moveCard(gameId, cardId, "GRAVEYARD");
    console.log("✅ MOVE successful");

    // Test CREATE_TOKEN
    console.log("Testing CREATE_TOKEN...");
    await createToken(
      gameId,
      playerId,
      "Goblin",
      "Creature - Goblin",
      "1",
      "1",
      "red",
    );
    console.log("✅ CREATE_TOKEN successful");

    console.log("All tests passed!");
  } catch (error) {
    console.error("Test failed:", error);
  }
}

// Run tests
testDirectActions();
```

**Verification Checklist**:

- [ ] TAP command taps card (visual update in UI)
- [ ] UNTAP command untaps card
- [ ] DRAW command adds card to hand
- [ ] MODIFY_LIFE changes player life total
- [ ] SHUFFLE shuffles library (check LibraryCount stays same)
- [ ] MOVE command moves card to graveyard
- [ ] CREATE_TOKEN creates token on battlefield
- [ ] FLIP toggles card face-down state
- [ ] ADD_COUNTER adds counter to card
- [ ] REMOVE_COUNTER removes counter from card
- [ ] SET_COUNTER sets counter to specific value
- [ ] SET_PLAYER_COUNTER sets player counter (poison, energy)
- [ ] MILL moves cards from library to graveyard
- [ ] NEXT_TURN advances turn counter
- [ ] MULLIGAN replaces hand
- [ ] KEEP_HAND marks hand as kept
- [ ] All actions trigger GAME_UPDATE broadcast
- [ ] All players receive updates in real-time

**Anti-Patterns to Avoid**:

- ❌ Don't skip testing any command (all must work)
- ❌ Don't test only in single-player (test with 2+ players)
- ❌ Don't ignore server logs (watch for parsing errors)
- ❌ Don't assume success (verify UI updates correctly)

---

### Task 2.3: Verify ActionQueue Processing

**What to Implement**:
Confirm that the ProcessGameActions goroutine is consuming actions from the queue and that error handling works correctly.

**Documentation Reference**:

- Goroutine spawn: `mage-server-go/internal/server/grpc_game.go` line 131
- ProcessGameActions implementation: `mage-server-go/internal/game/manager.go` lines 312-327
- ActionQueue channel: `mage-server-go/internal/game/manager.go` line 40

**Files to Verify** (no changes needed):

1. `mage-server-go/internal/server/grpc_game.go` (goroutine spawn)
2. `mage-server-go/internal/game/manager.go` (queue processor)

**Verification Steps**:

1. Add debug logging to track action queue processing (optional, for verification):
   - In `manager.go` ProcessGameActions, add more logging around action processing
   - Check server logs show "processing action" messages

2. Start a game and perform multiple actions rapidly:
   - Tap 5 cards in quick succession
   - Check server logs show all 5 actions processed in order

3. Test queue overflow (100 action limit):
   - Write a script to send 150 actions rapidly
   - Verify error "action queue full" appears for actions 101-150
   - Verify first 100 actions still process correctly

4. Test error handling:
   - Send an invalid action (e.g., TAP a non-existent card)
   - Check server logs: "failed to process action" error logged
   - Verify game doesn't crash
   - Verify queue continues processing subsequent valid actions

**Expected Behavior**:

- ✅ Goroutine runs for entire game lifetime
- ✅ Actions processed in FIFO order
- ✅ Errors logged but don't crash game
- ✅ Queue handles burst traffic (up to 100 pending actions)
- ✅ Queue full errors handled gracefully

**Anti-Patterns to Avoid**:

- ❌ Don't skip stress testing (queue could deadlock under load)
- ❌ Don't ignore error logs (errors should be logged, not silent)
- ❌ Don't assume FIFO ordering (verify with timestamps)

---

### Phase 2 Completion Checklist

**Before moving to Phase 3, verify**:

- [ ] String command parser implemented and compiles
- [ ] "SEND_STRING" case added to ProcessAction
- [ ] All 16 direct action commands tested and working:
  - [ ] TAP / UNTAP
  - [ ] MOVE
  - [ ] FLIP
  - [ ] DRAW
  - [ ] MODIFY_LIFE
  - [ ] SET_COUNTER (player)
  - [ ] SHUFFLE
  - [ ] CREATE_TOKEN
  - [ ] ADD_COUNTER / REMOVE_COUNTER / SET_CARD_COUNTER
  - [ ] MILL
  - [ ] SCRY
  - [ ] REVEAL_TOP
  - [ ] NEXT_TURN
  - [ ] MULLIGAN / KEEP_HAND
- [ ] ActionQueue goroutine verified running
- [ ] Actions process in correct order
- [ ] Error handling works (invalid actions don't crash game)
- [ ] All actions trigger real-time GAME_UPDATE to all players
- [ ] UI updates correctly for all action types

---

## Phase 3: Polish & Advanced Features

**Goal**: Add optional enhancements for better UX and completeness.

**Dependencies**: Phase 2 complete (actions working).

**Testing After Phase**: Cards display with full metadata, UI polish features work, all edge cases handled.

---

### Task 3.1: Client-Side Card Metadata Enrichment

**What to Implement** (OPTIONAL):
Add metadata enrichment for cards that arrive with minimal data (name only). This improves card tooltips and type-based filtering.

**Documentation Reference**:

- Current pattern: `mage-client-web/src/lib/components/game/Card.svelte` (uses getScryfallImageUrl)
- Scryfall utilities: `mage-client-web/src/lib/utils/scryfall.ts`
- Card types: `mage-client-web/src/lib/types/game.ts`

**Files to Modify**:

1. `mage-client-web/src/lib/stores/multiplayer-game.ts` (add enrichment helper)
2. `mage-client-web/src/lib/utils/scryfall.ts` (add metadata fetch function if not exists)

**Implementation Steps**:

1. Add Scryfall card data fetch function (if not already exists):

```typescript
// In scryfall.ts
export interface ScryfallCard {
  name: string;
  mana_cost?: string;
  type_line?: string;
  oracle_text?: string;
  power?: string;
  toughness?: string;
  colors?: string[];
  image_uris?: {
    small?: string;
    normal?: string;
    large?: string;
  };
}

/**
 * Fetch full card metadata from Scryfall API
 */
export async function getScryfallCard(
  cardName: string,
): Promise<ScryfallCard | null> {
  try {
    const encodedName = encodeURIComponent(cardName);
    const response = await fetch(
      `https://api.scryfall.com/cards/named?exact=${encodedName}`,
    );

    if (!response.ok) {
      console.warn(`Failed to fetch Scryfall data for ${cardName}`);
      return null;
    }

    return await response.json();
  } catch (error) {
    console.error(`Error fetching Scryfall data for ${cardName}:`, error);
    return null;
  }
}
```

2. Add enrichment function in multiplayer-game.ts:

```typescript
/**
 * Enriches cards with Scryfall metadata if missing
 * Only enriches cards that don't have type/manaCost data
 */
async function enrichCardsWithMetadata(cards: CardView[]): Promise<CardView[]> {
  const enrichPromises = cards.map(async (card) => {
    // Skip if card already has metadata
    if (card.type && card.manaCost) {
      return card;
    }

    try {
      const metadata = await getScryfallCard(card.name);
      if (!metadata) {
        return card;
      }

      // Merge Scryfall data with card
      return {
        ...card,
        type: metadata.type_line || card.type,
        manaCost: metadata.mana_cost || card.manaCost,
        power: metadata.power || card.power,
        toughness: metadata.toughness || card.toughness,
        rulesText: metadata.oracle_text || card.rulesText,
        color: metadata.colors?.join("") || card.color,
      };
    } catch (error) {
      console.warn(`Failed to enrich card ${card.name}:`, error);
      return card;
    }
  });

  return await Promise.all(enrichPromises);
}
```

3. Use enrichment in GAME_UPDATE handler (OPTIONAL - only if needed):

```typescript
// In GAME_UPDATE handler, after mapping:
const mappedState = mapProtoGameViewToState(normalized);

// Optionally enrich battlefield cards (or all cards)
const enrichedBattlefield = await enrichCardsWithMetadata(
  mappedState.battlefield || [],
);

update((state) => ({
  ...state,
  ...mappedState,
  battlefield: enrichedBattlefield,
  // ... rest of update
}));
```

**Note**: This is **optional** because:

- Card images already work (getScryfallImageUrl)
- Adds API call overhead (rate limiting concerns)
- Server-side enrichment would be more efficient (future enhancement)

**Verification** (if implemented):

- ✅ Cards without metadata get enriched on first load
- ✅ Enrichment doesn't block rendering (async)
- ✅ Tooltips show mana cost, type, rules text
- ✅ No excessive Scryfall API calls (cache or deduplicate)

**Anti-Patterns to Avoid**:

- ❌ Don't enrich every card on every update (cache results)
- ❌ Don't block rendering waiting for enrichment (async)
- ❌ Don't spam Scryfall API (rate limit 10 req/sec)
- ❌ Don't enrich cards that already have data

---

### Task 3.2: Add Missing Client API Functions

**What to Implement** (OPTIONAL):
The server has additional operations that the client doesn't expose. Add client wrappers if needed.

**Documentation Reference**:

- Server action handlers: `mage-server-go/internal/game/actions.go`
- Existing client API: `mage-client-web/src/lib/api/direct-actions.ts`

**Missing Client Functions** (already exist on server):

- None - all server operations have client wrappers

**Files to Check**:

1. `mage-client-web/src/lib/api/direct-actions.ts`

**Verification**:

- ✅ All server action types have corresponding client functions
- ✅ No additional client wrappers needed

**Note**: If future server actions are added, create corresponding client functions following the pattern:

```typescript
export async function newAction(gameId: string, param: string): Promise<void> {
  return sendPlayerString(gameId, `NEW_ACTION:${param}`);
}
```

---

### Task 3.3: Add Optimistic UI Updates (OPTIONAL)

**What to Implement** (OPTIONAL):
Add optimistic updates for better perceived performance. Show action result immediately, then reconcile with server state.

**Documentation Reference**:

- Example: `mage-client-web/src/lib/stores/game.legacy.ts` (pendingCardPlays pattern)

**Files to Modify**:

1. `mage-client-web/src/lib/stores/multiplayer-game.ts`
2. `mage-client-web/src/lib/api/direct-actions.ts` (wrap actions)

**Implementation Pattern**:

```typescript
// In multiplayer-game.ts, add to state:
pendingTaps: Map<string, boolean>; // cardId -> tapped state

// Wrap direct action:
export async function tapUntapOptimistic(
  gameId: string,
  cardId: string,
  tapped: boolean,
): Promise<void> {
  // Apply optimistic update
  multiplayerGame.update((state) => {
    const pendingTaps = new Map(state.pendingTaps);
    pendingTaps.set(cardId, tapped);
    return { ...state, pendingTaps };
  });

  try {
    // Send to server
    await tapUntap(gameId, cardId, tapped);
  } catch (error) {
    // Rollback on error
    multiplayerGame.update((state) => {
      const pendingTaps = new Map(state.pendingTaps);
      pendingTaps.delete(cardId);
      return { ...state, pendingTaps };
    });
    throw error;
  }
}

// In GAME_UPDATE handler, clear pending when server confirms:
update((state) => {
  const newPendingTaps = new Map(state.pendingTaps);
  // Remove confirmed taps
  for (const card of mappedState.battlefield || []) {
    newPendingTaps.delete(card.id);
  }
  return {
    ...state,
    ...mappedState,
    pendingTaps: newPendingTaps,
  };
});
```

**Note**: This is **optional** because:

- Network latency is typically low enough
- Adds complexity for marginal UX improvement
- Server state sync is already fast (WebSocket)

**Verification** (if implemented):

- ✅ Actions appear instant in UI
- ✅ Server state reconciles correctly
- ✅ Rollback works on error
- ✅ No visual glitches from race conditions

---

### Phase 3 Completion Checklist

**Optional enhancements completed as desired**:

- [ ] Card metadata enrichment (if implemented)
- [ ] All client API functions present (verified)
- [ ] Optimistic UI updates (if implemented)
- [ ] Any custom polish features

---

## Phase 4: Verification & Testing

**Goal**: Comprehensive testing to ensure all features work correctly in production scenarios.

**Dependencies**: All previous phases complete.

**Testing Checklist**:

### 4.1 State Synchronization Tests

- [ ] Start 2-player game
- [ ] Both players see same battlefield state
- [ ] Each player sees their own hand/library
- [ ] Opponents' hands/libraries are hidden (counts only)
- [ ] Turn counter syncs correctly
- [ ] Active player indicator correct
- [ ] Refresh page: state reloads correctly from server

### 4.2 Action Tests

- [ ] Player 1 taps card → Player 2 sees tap
- [ ] Player 2 draws card → Player 1 sees hand count increase
- [ ] Modify life → all players see updated life total
- [ ] Create token → all players see new token
- [ ] Move card → all players see zone change
- [ ] Shuffle library → library count stays same, order changes
- [ ] Mulligan → hand replaced, graveyard updated
- [ ] Next turn → turn counter increments, active player changes

### 4.3 Edge Cases

- [ ] Rapid actions (5 actions in 1 second) → all process correctly
- [ ] Concurrent actions (both players tap cards simultaneously) → both resolve
- [ ] Invalid action (tap non-existent card) → error logged, game continues
- [ ] Network disconnect → reconnect and sync state
- [ ] Empty library → draw doesn't crash
- [ ] 100+ cards on battlefield → performance acceptable

### 4.4 Multi-Player Tests

- [ ] 3-player game: all players receive updates
- [ ] 4-player game: turn order correct
- [ ] Watcher joins: sees game state correctly
- [ ] Player leaves: game continues for remaining players

### 4.5 UI/UX Tests

- [ ] Cards display with images
- [ ] Tapped cards show rotation
- [ ] Counters display correctly
- [ ] Life totals update smoothly
- [ ] No console errors
- [ ] No TypeScript compilation errors
- [ ] No Go compilation errors

### 4.6 Performance Tests

- [ ] Game with 60-card decks loads in <2 seconds
- [ ] GAME_UPDATE processing <100ms
- [ ] No memory leaks (check after 100+ actions)
- [ ] WebSocket connection stable (30+ minutes)

---

## Anti-Pattern Guards

### Common Mistakes to Avoid

**Proto Changes**:

- ❌ Forgetting to regenerate protos after editing .proto files
- ❌ Using field numbers already in use
- ❌ Modifying generated files manually

**WebSocket State**:

- ❌ Mutating state directly instead of using update()
- ❌ Skipping normalization of proto GameView (missing arrays)
- ❌ Not handling null/undefined GameView gracefully

**Action Processing**:

- ❌ Sending structured data where string expected (or vice versa)
- ❌ Not validating command parameters before parsing
- ❌ Blocking action queue (synchronous processing)
- ❌ Not logging errors (silent failures)

**Card Data**:

- ❌ Assuming all cards have metadata (some are name-only)
- ❌ Spamming Scryfall API without caching
- ❌ Blocking render waiting for metadata fetch

**Testing**:

- ❌ Testing only single-player (multi-player has race conditions)
- ❌ Not testing edge cases (empty zones, invalid actions)
- ❌ Skipping performance tests (queue could deadlock)

---

## Rollback Plan

If issues arise during implementation:

### Phase 1 Rollback

- Revert proto changes: `git checkout mage-server-go/api/proto/mage/v1/models.proto`
- Regenerate: `cd mage-server-go && make proto`
- Revert GAME_UPDATE handler to original placeholder

### Phase 2 Rollback

- Remove "SEND_STRING" case from ProcessAction
- Remove ParseAndExecuteStringCommand method
- Client actions will return errors but won't crash

### Phase 3 Rollback

- Optional features can be disabled without affecting core functionality

---

## Success Criteria

**Phase 1 Success**:

- ✅ Game page loads without errors
- ✅ All players display with correct data
- ✅ Cards appear in zones
- ✅ Turn counter shows

**Phase 2 Success**:

- ✅ All direct actions work from browser console
- ✅ UI components can trigger actions
- ✅ Real-time updates for all players

**Phase 3 Success**:

- ✅ Optional enhancements work as designed
- ✅ No regressions in core functionality

**Complete Success**:

- ✅ All 10 TODOs from docs/20260124-missing-pieces-multiplayer-interface.md resolved
- ✅ Multiplayer rules-light game fully functional
- ✅ No console errors
- ✅ All tests pass
- ✅ Performance acceptable
- ✅ Ready for production use

---

## Appendix: File Reference

### Server Files (Go)

- `mage-server-go/api/proto/mage/v1/models.proto` - Proto definitions
- `mage-server-go/internal/server/grpc.go` - WebSocket notifications, proto conversion
- `mage-server-go/internal/server/grpc_game.go` - RPC handlers, goroutine spawn
- `mage-server-go/internal/game/game_engine.go` - ProcessAction, broadcast
- `mage-server-go/internal/game/actions.go` - Action handler implementations
- `mage-server-go/internal/game/manager.go` - ActionQueue, ProcessGameActions
- `mage-server-go/internal/game/view.go` - buildGameView (hidden information filtering)
- `mage-server-go/internal/game/game_state.go` - Card, Player, GameState structs

### Client Files (TypeScript)

- `mage-client-web/src/lib/stores/multiplayer-game.ts` - Main game state store
- `mage-client-web/src/lib/stores/game.legacy.ts` - Working example (reference)
- `mage-client-web/src/lib/api/game.ts` - RPC wrappers (sendPlayerString, etc.)
- `mage-client-web/src/lib/api/direct-actions.ts` - Direct action commands
- `mage-client-web/src/lib/types/game.ts` - TypeScript type definitions
- `mage-client-web/src/lib/utils/scryfall.ts` - Scryfall API utilities
- `mage-client-web/src/lib/components/game/Card.svelte` - Card rendering

### Build Commands

- `cd mage-server-go && make proto` - Regenerate proto files
- `cd mage-client-web && bun run build` - Build client
- `cd mage-server-go && go build ./cmd/server` - Build server

---

## Conclusion

This plan provides a complete, step-by-step implementation guide for all 10 TODOs from the multiplayer interface documentation. Each phase builds on the previous, with clear verification steps and anti-pattern guards to prevent common mistakes.

**Estimated Implementation Time**:

- Phase 1 (State Sync): 2-3 hours
- Phase 2 (Actions): 2-3 hours
- Phase 3 (Polish): 1-2 hours (optional)
- Phase 4 (Testing): 1-2 hours

**Total**: 6-10 hours for complete implementation and testing.

The plan can be executed in parallel across multiple sessions by assigning phases to different contexts, as long as Phase 1 completes before Phase 2 begins.
