Summary: Finalizing Multiplayer Interface

Current State: TODO
Date: Jan 24, 2026

Server (Go):

- PlaytestGameView structure exists in view.go with proper hidden information filtering
- Conversion to protobuf GameView happens in grpc.go:playtestViewToProto()
- Server sends GAME_UPDATE via WebSocket with the converted view

Frontend (TypeScript):

- multiplayer-game.ts store has correct interface expecting all the right data
- Game page components expect this data structure
- BUT: GAME_UPDATE handler receives GameView but doesn't apply it to store state (lines 162-179 just set  
  isConnected: true)

What Needs to Be Done

1. GameView State Mapping (CRITICAL)

Problem: Server's PlaytestGameView has nested player structure (Me/Opponents) but proto GameView flattens it to  
Players[]. Frontend expects PlaytestPlayer[].

Location: multiplayer-game.ts lines 162-179

Fix: Add conversion function in GAME_UPDATE handler:

websocketStore.on(CallbackMethod.GAME_UPDATE, (data) => {  
 const updateData = data as GameUpdateData;  
 if (updateData.game) {  
 update((state) => ({  
 ...state,  
 gameId: updateData.game.gameId,  
 turn: Number(updateData.game.turn),  
 activePlayerId: updateData.game.activePlayerId || '',  
 battlefield: updateData.game.battlefield || [],  
 exile: updateData.game.exile || [],  
 stack: updateData.game.stack || [],  
 command: updateData.game.command || [],  
 players: (updateData.game.players || []).map(convertPlayerViewToPlaytestPlayer),  
 isConnected: true,  
 isInitialized: true  
 }));  
 }  
});

function convertPlayerViewToPlaytestPlayer(pv: PlayerView): PlaytestPlayer {  
 return {  
 playerId: pv.playerId,  
 name: pv.name,  
 life: Number(pv.life),  
 poison: Number(pv.poison),  
 energy: Number(pv.energy),  
 libraryCount: Number(pv.libraryCount),  
 handCount: Number(pv.handCount),  
 hand: pv.hand || [],  
 library: [], // Hidden for opponents, full for you  
 graveyard: pv.graveyard || [],  
 manaPool: pv.manaPool || { white: 0, blue: 0, black: 0, red: 0, green: 0, colorless: 0 },  
 keptHand: pv.keptHand || false,  
 mulliganCount: 0,  
 revealedTopCard: false  
 };  
}

2. Server GameView Extension

Check: grpc.go:playtestViewToProto() lines 648-697 - already maps:

- ✅ Me.Library → PlayerView.Hand (for viewing player)
- ✅ Opponents.Hand → empty array (hidden)
- ❌ MISSING: Me.Library field not sent to frontend

Fix needed in grpc.go:667-677:  
players = append(players, &pb.PlayerView{  
 PlayerId: data.Me.PlayerID,  
 // ... existing fields ...  
 Hand: playtestEngineCardsToProto(data.Me.Hand),  
 Library: playtestEngineCardsToProto(data.Me.Library), // ADD THIS  
 Graveyard: playtestEngineCardsToProto(data.Me.Graveyard),  
})

Also check proto definition - models.proto:114-134 has library field but it's not in the mapping!

3. Missing Server Operations

These operations log warnings - need server implementation:

- MILL command - move top N from library to graveyard atomically
- SCRY command - reveal top N, choose order/bottom
- REVEAL_TOP - temporarily show top cards
- MULLIGAN / KEEP_HAND - mulligan phase handling
- REVEAL_TOP_PERMANENT - permanently reveal top card

Files to modify:

- mage-server-go/internal/game/actions.go - add new action handlers
- mage-client-web/src/lib/api/direct-actions.ts - add API calls

4. Initial State Application

Location: multiplayer-game.ts lines 219-235

Currently fetches \_gameView but doesn't use it. Apply same conversion as GAME_UPDATE:

const gameView = await fetchGameView(gameId, playerId);  
update((state) => ({  
 ...state,  
 ...convertGameViewToState(gameView), // Same logic as GAME_UPDATE  
 isConnected: true,  
 isInitialized: true  
}));

5. Derived Stores Verification

Location: multiplayer-game.ts lines 673-699

These already derive correctly IF the players array is populated:

- ✅ multiplayerPlayers - derives from $game.players
- ✅ multiplayerLocalPlayer - finds by activeControlSeat
- ✅ multiplayerOpponents - filters out local player

Should work once step 1 is complete.

Priority Order

1. Fix GameView mapping (Step 1) - unblocks everything
2. Add Library field to server (Step 2) - enables full hand/library view
3. Apply initial state (Step 4) - enables page load
4. Implement missing ops (Step 3) - enables full gameplay
5. Verify derived stores (Step 5) - should just work

The core issue is that the WebSocket GAME_UPDATE is being received but ignored. Fix the mapping first, then
everything else should fall into place.

## Additional Critical Missing Pieces (6-10)

### 6. Direct Action Command Parsing

**Problem**: Client sends string commands like `"TAP:cardId"` via `sendPlayerString()`, but server's `SendPlayerString`
(grpc_game.go:454) passes it as `"SEND_STRING"` action type. The GameEngine's `ProcessAction` (game_engine.go:148-255)
doesn't handle `"SEND_STRING"` - it only knows about structured action types like `"TAP"`, `"MOVE"`, etc.

**Files affected**:
- `mage-server-go/internal/game/manager.go:291-310` - SendPlayerAction sends action to queue
- `mage-server-go/internal/game/game_engine.go:146-255` - ProcessAction switch statement
- `mage-client-web/src/lib/api/direct-actions.ts` - Sends string commands

**Fix needed**: Add command parser in ProcessAction that converts `"SEND_STRING"` with data like `"TAP:cardId"` into
structured actions:

```go
case "SEND_STRING":
    // Parse direct action commands from string
    command := getStringFromData(data, "command", "")
    parts := strings.Split(command, ":")
    if len(parts) == 0 {
        return fmt.Errorf("invalid command format")
    }

    switch parts[0] {
    case "TAP", "UNTAP":
        cardID := parts[1]
        tapped := parts[0] == "TAP"
        return e.TapCard(gameID, playerID, cardID, tapped)
    case "MOVE":
        return e.MoveCard(gameID, playerID, parts[1], parts[2])
    case "FLIP":
        faceDown, _ := strconv.ParseBool(parts[2])
        return e.FlipCard(gameID, playerID, parts[1], faceDown)
    case "CREATE_TOKEN":
        // Parse CREATE_TOKEN:name:types:power:toughness:color
        return e.CreateToken(gameID, playerID, parts[1], parts[2], parts[3], parts[4], parts[5])
    // ... handle all direct-actions.ts commands
    }
```

### 7. Action Queue Processing Missing

**Problem**: Actions are sent to `game.ActionQueue` channel (manager.go:302-308) but **nothing reads from it**. The
engine needs a goroutine processing this queue. Code in grpc_game.go:131 references `ProcessGameActions` but it's not
implemented.

**Files affected**:
- `mage-server-go/internal/game/manager.go:291-310` - Sends to queue but nothing consumes it
- `mage-server-go/internal/server/grpc_game.go:125-132` - Should start processor
- `mage-server-go/internal/game/manager.go` - Missing ProcessGameActions implementation

**Fix needed**: Add ProcessGameActions goroutine:

```go
// In grpc_game.go after StartGameWithDecks (line 131):
go s.gameAdapter.ProcessGameActions(gameInstance)

// In manager.go, add new method to EngineAdapter:
func (ea *EngineAdapter) ProcessGameActions(game *Game) {
    ea.logger.Info("starting action processor", zap.String("game_id", game.ID))

    for action := range game.ActionQueue {
        ea.logger.Debug("processing action",
            zap.String("game_id", game.ID),
            zap.String("player_id", action.PlayerID),
            zap.String("action_type", action.ActionType))

        if err := ea.engine.ProcessAction(game.ID, action); err != nil {
            ea.logger.Error("action processing failed",
                zap.Error(err),
                zap.String("game_id", game.ID),
                zap.String("action_type", action.ActionType))
        }
    }

    ea.logger.Info("action processor stopped", zap.String("game_id", game.ID))
}
```

### 8. Card Metadata Enrichment

**Problem**: Cards are created with only `Name` field (game_engine.go:96-107). No type, mana cost, power/toughness,
rules text. Frontend components need this metadata for proper display (card tooltips, type-based filtering, etc).

**Files affected**:
- `mage-server-go/internal/game/game_engine.go:84-131` - Card creation from deck lists
- `mage-client-web/src/lib/utils/scryfall.ts` - Card data fetching utilities
- `mage-client-web/src/lib/components/game/Card.svelte` - Needs full card data

**Options**:

**Option 1 - Server-side enrichment** (recommended):
- Load card data from Scryfall API during StartGameWithDecks
- Cache card data in CardRepository
- Pros: Single source of truth, immediate availability
- Cons: Server API calls, startup delay

**Option 2 - Client-side enrichment**:
- Fetch card data when receiving GameView
- Cache in browser localStorage
- Pros: No server changes needed
- Cons: Network requests per client, inconsistent data

**Implementation for Option 2** (quick fix):

```typescript
// In multiplayer-game.ts after applying GameView:
async function enrichCardsWithMetadata(cards: CardView[]): Promise<CardView[]> {
    const enriched = await Promise.all(cards.map(async (card) => {
        if (!card.type || !card.manaCost) {
            try {
                const metadata = await getScryfallCard(card.name);
                return {
                    ...card,
                    type: metadata.type_line,
                    manaCost: metadata.mana_cost,
                    power: metadata.power,
                    toughness: metadata.toughness,
                    rulesText: metadata.oracle_text,
                    color: metadata.colors?.join('') || 'colorless'
                };
            } catch (err) {
                console.warn(`Failed to enrich card: ${card.name}`, err);
                return card;
            }
        }
        return card;
    }));
    return enriched;
}

// Apply after GameView conversion:
battlefield: await enrichCardsWithMetadata(updateData.game.battlefield || []),
```

### 9. Library Field Not Transmitted

**Problem**: `playtestViewToProto` (grpc.go:675) sends `Hand` but **NOT** `Library` for viewing player. Proto definition
(models.proto:114-134) has the library field but the Go mapping omits it.

**Files affected**:
- `mage-server-go/internal/server/grpc.go:666-677` - Player view conversion (missing library)
- `mage-server-go/api/proto/mage/v1/models.proto:123` - Proto has library field
- `mage-client-web/src/lib/stores/multiplayer-game.ts:31` - Expects library array

**Fix**: Add library field to server conversion (already identified in TODO #2 but worth repeating):

```go
// In grpc.go:playtestViewToProto, around line 675:
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
```

Then update frontend conversion to use it:

```typescript
function convertPlayerViewToPlaytestPlayer(pv: PlayerView): PlaytestPlayer {
    return {
        // ... existing fields ...
        library: pv.library || [], // Use server data instead of empty array
        // ... rest of fields ...
    };
}
```

### 10. Missing activeControlSeat Field

**Problem**: GameView proto doesn't include `activeControlSeat` field, but multiplayer-game.ts expects it (line 70, 73,
132). Server needs to set this to indicate which player perspective the view is for. Currently derived stores fail
because they can't determine "me" vs "opponents".

**Files affected**:
- `mage-server-go/api/proto/mage/v1/models.proto:72-103` - GameView proto definition
- `mage-server-go/internal/server/grpc.go:648-696` - playtestViewToProto conversion
- `mage-client-web/src/lib/stores/multiplayer-game.ts:70,132,676` - Expects activeControlSeat

**Fix needed**:

**1. Add field to proto** (models.proto):
```proto
message GameView {
    string game_id = 1;
    string state = 2;
    repeated PlayerView players = 3;
    string active_player_id = 4;
    string priority_player_id = 5;
    int32 turn = 6;
    // ... existing fields ...

    // Which player perspective this view is for (their own player ID)
    string active_control_seat = 19;
}
```

**2. Set in server conversion** (grpc.go:playtestViewToProto):
```go
view := &pb.GameView{
    GameId:            data.GameID,
    State:             "IN_PROGRESS",
    Turn:              int32(data.Turn),
    ActivePlayerId:    data.ActivePlayerID,
    ActiveControlSeat: playerID, // ADD THIS - each player gets their own ID as control seat
    Battlefield:       playtestEngineCardsToProto(data.Battlefield),
    // ... rest of fields ...
}
```

**3. Apply in frontend GAME_UPDATE handler** (multiplayer-game.ts):
```typescript
update((state) => ({
    ...state,
    activeControlSeat: updateData.game.activeControlSeat || playerID, // Set from GameView
    // ... rest of state ...
}));
```

**4. Regenerate protos**:
```bash
cd mage-server-go && make proto
cd ../mage-client-web && bun run proto:generate
```

## Complete Implementation Order

### Phase 1: Core State Synchronization (Unblocks UI)
1. **TODO #10** - Add activeControlSeat field (proto + server + client)
2. **TODO #1** - Fix GameView mapping in GAME_UPDATE handler
3. **TODO #4** - Apply initial state from fetchGameView
4. **TODO #9** - Add Library field transmission

### Phase 2: Action Processing (Enables Gameplay)
5. **TODO #6** - Implement direct action command parsing
6. **TODO #7** - Add action queue processing goroutine
7. **TODO #2** - Verify all proto mappings complete

### Phase 3: Polish & Advanced Features
8. **TODO #8** - Card metadata enrichment (client or server-side)
9. **TODO #3** - Implement missing operations (MILL, SCRY, MULLIGAN)
10. **TODO #5** - Verify derived stores work correctly

### Testing After Each Phase:
- **Phase 1**: Page loads, shows players/zones, no errors
- **Phase 2**: Can tap cards, move cards, modify life
- **Phase 3**: All UI features work, cards display properly

## Summary

**First 5 TODOs** (State Sync):
1. GameView state mapping in GAME_UPDATE
2. Server GameView completeness
3. Missing server operations
4. Initial state application
5. Derived stores verification

**Next 5 TODOs** (Action Processing & Polish):
6. Direct action command parsing
7. Action queue processing
8. Card metadata enrichment
9. Library field transmission
10. activeControlSeat field

These 10 items create a **fully functional multiplayer rules-light game engine**.
