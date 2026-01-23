# Game Architecture Documentation

## Overview

The Mage game engine uses a **rules-light, player-controlled architecture** inspired by Untap.in. This design prioritizes flexibility, simplicity, and player agency over strict rules enforcement.

## Core Philosophy

**Players control the game state directly** - no automatic validation, no forced turn structure, no rules enforcement. The engine provides:

- Direct manipulation of game objects (tap, move, counters, life totals)
- Rollback capability for mistake recovery
- WebSocket synchronization for multiplayer
- Hidden information filtering (opponent hands/libraries)

**What it does NOT do:**

- Validate spell timing or mana costs
- Enforce priority or turn structure
- Check creature abilities (flying, menace, etc.)
- Auto-resolve state-based actions
- Process triggered abilities
- Validate targets
- Calculate combat damage automatically

## Architecture Layers

### 1. Backend: Rules-Light Engine

**Location**: `/mage-server-go/internal/game/`

**Key Components:**

#### Engine Selection (manager.go)

The game manager supports two engine implementations:

```go
// Engine selection via config
if config.UsePlaytestEngine {
    engine := NewPlaytestEngine(logger, cardRepo)
} else {
    engine := NewMageEngine(logger, cardRepo) // Legacy rules-enforced
}
```

Both engines implement the `GameEngine` interface, allowing seamless switching via configuration.

#### Playtest Engine (playtest_engine.go)

The rules-light engine with minimal state management:

```go
type PlaytestEngine struct {
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
```

**Operations** (playtest_actions.go):
- `TapCard(cardId, tapped)` - Toggle tap state
- `MoveCard(cardId, zone)` - Move between zones
- `DrawCards(playerId, count)` - Draw from library
- `ModifyLife(playerId, delta)` - Change life total
- `AddCounter(cardId, type, amount)` - Add/remove counters
- `CreateToken(...)` - Create token on battlefield
- `ShuffleLibrary(playerId)` - Randomize library
- `ScryCards(playerId, count)` - Look at top N cards
- `MillCards(playerId, count)` - Move top N to graveyard
- `NextTurn()` - Advance turn counter

#### Hidden Information Filtering (playtest_view.go)

Server filters game state per player:

```go
func (e *PlaytestEngine) GetGameView(gameID, playerID string) *GameView {
    view := &GameView{
        MyHand:       getPlayerCards(state, playerID, "hand"),      // Full cards
        MyLibrary:    getLibraryCount(state, playerID),             // Count only
        Battlefield:  getAllBattlefieldCards(state),                // Public zone
        Graveyard:    getAllGraveyardCards(state),                  // Public zone
        Exile:        getAllExileCards(state),                      // Public zone
        Opponents:    getOpponentViews(state, playerID),            // Hidden hands
    }
}

func getOpponentViews(state *GameState, viewerID string) []*OpponentView {
    // Return hand counts, not actual cards
    // Return library counts, not actual cards
    // Return battlefield cards (public)
}
```

#### Rollback System (playtest_rollback.go)

State snapshots for mistake recovery:

```go
type Bookmark struct {
    ID        string
    GameID    string
    State     *GameState  // Deep copy
    Timestamp time.Time
}

func (e *PlaytestEngine) BookmarkState(gameID string) string
func (e *PlaytestEngine) RestoreState(gameID, bookmarkID string) error
```

**Automatic bookmarks**: Created at turn boundaries
**Manual bookmarks**: Players can request snapshots
**Consent flow**: Multiplayer rollbacks require approval

### 2. Frontend: Playtest-Based UI

**Location**: `/mage-client-web/src/`

#### Game Store (stores/multiplayer-game.ts)

Manages local game state and server synchronization:

```typescript
export interface MultiplayerGameState {
  gameId: string;
  activeControlSeat: string;
  players: PlaytestPlayer[];
  battlefield: CardView[];
  exile: CardView[];
  stack: CardView[];
  command: CardView[];
  graveyard: CardView[];
  turn: number;
  activePlayerId: string;
  isInitialized: boolean;
  log: PlaytestLogEntry[];

  // Server sync fields
  isConnected: boolean;
  pendingActions: string[];
}
```

**Store operations** call backend via direct-actions API:

```typescript
function drawCards(playerId: string, count: number): void {
  directActions.drawCards(gameId, playerId, count);
  // Server broadcasts update back to all clients
}

function moveCardToZone(cardId: string, targetZone: string): void {
  directActions.moveCard(gameId, playerId, cardId, targetZone);
}

function tapCard(cardId: string, tapped: boolean): void {
  directActions.tapUntap(gameId, playerId, cardId, tapped);
}
```

#### WebSocket Synchronization

Server broadcasts state changes to all players:

```typescript
websocketStore.subscribe((message) => {
  if (message.type === 'GAME_UPDATE') {
    multiplayerGameStore.update((state) => ({
      ...state,
      ...message.gameState,
      isConnected: true
    }));
  }
});
```

#### UI Components

**Core Components** (used by both playtest and multiplayer game):
- `Card.svelte` - Card rendering with hover/drag
- `PlayerHand.svelte` - Hand display with drag-drop
- `BattlefieldArea.svelte` - Battlefield layout (lands/nonlands)
- `Graveyard.svelte` - Graveyard zone
- `ExileZone.svelte` - Exile zone
- `ScryDialog.svelte` - Scry interface
- `TokenCreator.svelte` - Token creation
- `CounterDialog.svelte` - Counter management

**Multiplayer Components** (game view only):
- `OpponentSection.svelte` - Opponent battlefield and info
- `PlayerInfoRow.svelte` - Multi-player life totals
- `GameChatOverlay.svelte` - Chat functionality

**Deleted Components** (rules-enforcement, no longer used):
- `PriorityActionBar.svelte` - Priority system
- `DeclareAttackers.svelte` - Combat validation
- `DeclareBlockers.svelte` - Combat validation
- `AssignDamage.svelte` - Damage validation
- `ManaPayment.svelte` - Mana prompts
- `XManaSelector.svelte` - X mana prompts
- `AbilitiesPanel.svelte` - Ability prompts

#### Interaction Patterns

**Keyboard Shortcuts** (from playtest):
- `C` - Draw card
- `V` - Shuffle library
- `X` - Untap all
- `E` - Next turn
- `Space` - Pass priority (cosmetic)
- `T` - Tap selected card
- `+`/`-` - Modify life

**Drag-Drop System**:
- Cards draggable between zones
- Drop zones highlight on drag start
- Valid zones determined by card location
- Token rule: tokens leaving battlefield cease to exist

**Context Menus**:
- Deck: Draw N, Mill N, Scry N, Shuffle
- Card: Tap/Untap, Add Counter, Remove Counter, Move to Zone
- Life Total: +1, -1, +5, -5, Custom

### 3. API Layer

#### Direct Actions API (api/direct-actions.ts)

String-based command interface to backend:

```typescript
export async function tapUntap(
  gameId: string,
  playerId: string,
  cardId: string,
  tapped: boolean
): Promise<void> {
  const command = tapped ? `TAP:${cardId}` : `UNTAP:${cardId}`;
  await sendPlayerString(gameId, playerId, command);
}

export async function moveCard(
  gameId: string,
  playerId: string,
  cardId: string,
  targetZone: string
): Promise<void> {
  const command = `MOVE:${cardId}:${targetZone}`;
  await sendPlayerString(gameId, playerId, command);
}
```

**Backend parsing** (mage_engine.go lines 2590-2740):
```go
// ProcessPlayerString handles direct action commands
func (e *Engine) ProcessPlayerString(gameID, playerID, command string) error {
    parts := strings.Split(command, ":")

    switch parts[0] {
    case "TAP":
        return e.TapCard(gameID, parts[1], true)
    case "UNTAP":
        return e.TapCard(gameID, parts[1], false)
    case "MOVE":
        return e.MoveCard(gameID, parts[1], parts[2])
    // ... all other operations
    }
}
```

## Data Flow

### Game Start Flow

1. **Client** → `POST /api/games/start` → **Server**
2. **Server** → `PlaytestEngine.StartGame()` → Creates GameState
3. **Server** → Broadcast initial state → **All Clients**
4. **Clients** → `multiplayerGameStore.update()` → Render UI

### Player Action Flow

1. **Client** → User taps card → `tapCard(cardId, true)`
2. **Client** → `directActions.tapUntap()` → `SendPlayerString("TAP:cardId")`
3. **Server** → Parse command → `PlaytestEngine.TapCard()`
4. **Server** → Update state → Log action → Broadcast
5. **All Clients** → WebSocket message → `multiplayerGameStore.update()`
6. **All Clients** → Reactive UI updates → Show tap animation

### Rollback Flow

1. **Client A** → Request rollback → `POST /api/games/{id}/rollback`
2. **Server** → Create consent request → Broadcast to all players
3. **Client B,C,D** → Show consent dialog
4. **Client B,C,D** → Approve → `POST /api/games/{id}/rollback/approve`
5. **Server** → All approved → `RestoreState(bookmarkId)`
6. **Server** → Broadcast restored state → **All Clients**
7. **All Clients** → `multiplayerGameStore.update()` → Render restored state

## State Synchronization

### Server-Authoritative Model

**Server is source of truth**:
- Clients send commands to server
- Server updates state
- Server broadcasts to all clients
- Clients apply updates reactively

**No client-side prediction** (for now):
- Simple request-response model
- Adds ~50-200ms latency per action
- Acceptable for turn-based game
- Could add optimistic updates later

### Hidden Information

**Server filters per player**:

| Zone | Visibility |
|------|------------|
| My Hand | Full card data |
| My Library | Full card data (ordered) |
| Opponent Hand | Count only |
| Opponent Library | Count only |
| Battlefield | All players see all cards |
| Graveyard | All players see all cards |
| Exile | All players see all cards |
| Stack | All players see all cards |
| Command | All players see all cards |

### Zones and Card Movement

**Zone Types**:

| Zone ID | Name | Owner | Ordered | Hidden |
|---------|------|-------|---------|--------|
| `HAND` | Hand | Player | No | Yes |
| `LIBRARY` | Library | Player | Yes | Yes |
| `BATTLEFIELD` | Battlefield | Shared | No | No |
| `GRAVEYARD` | Graveyard | Player | Yes | No |
| `EXILE` | Exile | Shared | No | No |
| `STACK` | Stack | Shared | Yes | No |
| `COMMAND` | Command | Player | No | No |

**Movement Rules**:
- Cards can move freely between any zones
- Tokens leaving battlefield cease to exist
- Library order preserved unless shuffled
- Graveyard order preserved (can be important)

## Engine Selection

### Configuration

Backend supports both engines simultaneously:

```yaml
# config.yaml
game:
  default_engine: "playtest"  # or "mage"
```

**PlaytestEngine**: Rules-light, player-controlled (default)
**MageEngine**: Rules-enforced, automatic validation (legacy)

### Interface Compatibility

Both engines implement `GameEngine` interface:

```go
type GameEngine interface {
    StartGame(gameId string, players []Player, config GameConfig) error
    StartGameWithDecks(gameId string, decks []Deck) error
    ProcessAction(gameId, playerId string, action Action) error
    GetGameView(gameId, playerId string) (*GameView, error)
    EndGame(gameId string) error
}
```

This allows:
- Seamless engine switching via config
- Per-game engine selection (future)
- A/B testing between engines
- Gradual migration path

### Migration Path

**Current State** (Phase 7 Complete):
1. Both engines exist in backend
2. Frontend unified on playtest UI
3. Default config uses PlaytestEngine
4. Old game store renamed to `game.legacy.ts`

**Future Options**:
1. Delete MageEngine entirely (simplify codebase)
2. Keep both engines for different game modes
3. Add engine selection in game creation UI

## Code Organization

### Backend Files

**Playtest Engine**:
- `playtest_engine.go` - Core engine and state management
- `playtest_actions.go` - All game operations
- `playtest_rollback.go` - Bookmark/restore system
- `playtest_view.go` - Hidden information filtering
- `playtest_state.go` - State structures and helpers

**Legacy Rules Engine**:
- `mage_engine.go` - Full rules enforcement
- `engine_combat.go` - Combat system
- `engine_priority.go` - Priority system
- `engine_stack.go` - Stack resolution
- `engine_layers.go` - Continuous effects
- `engine_events.go` - Triggered abilities

**Shared**:
- `manager.go` - Game lifecycle, engine selection
- `grpc_game.go` - gRPC API handlers

### Frontend Files

**Stores**:
- `multiplayer-game.ts` - Multiplayer game state (playtest-based)
- `playtest-game.ts` - Solo playtest state
- `game.legacy.ts` - Old rules-enforced store (deprecated)

**Routes**:
- `routes/(protected)/game/[id]/+page.svelte` - Multiplayer game view
- `routes/(protected)/playtest/+page.svelte` - Solo playtest view

**Components** (shared):
- `components/game/Card.svelte`
- `components/game/PlayerHand.svelte`
- `components/game/BattlefieldArea.svelte`
- `components/game/OpponentSection.svelte`
- `components/game/PlayerInfoRow.svelte`
- `components/game/ScryDialog.svelte`
- `components/game/TokenCreator.svelte`
- `components/game/CounterDialog.svelte`

**API**:
- `api/direct-actions.ts` - String command API
- `api/game.ts` - Game lifecycle API

## Benefits of This Architecture

### Simplicity

**Backend**: ~2,500 lines (playtest engine) vs 13,786 lines (mage engine)
- 82% code reduction
- Easier to understand and maintain
- Fewer bugs and edge cases

**Frontend**: Single UI pattern (playtest-based)
- No duplication between playtest and game views
- Consistent UX across solo and multiplayer
- Simpler state management

### Flexibility

**No rules enforcement**:
- Supports house rules and casual play
- Players can undo mistakes (rollback)
- No "stuck in invalid state" issues
- Faster gameplay (no validation overhead)

**Player agency**:
- Full control over game state
- Manual combat resolution
- Direct manipulation of all game objects
- Can play any format (Commander, Cube, Limited)

### Performance

**No validation overhead**:
- Instant state updates (no rules checking)
- Simple state structure (no layers, triggers, stack)
- Fast WebSocket broadcasts
- Scales to 4+ player games easily

### Maintainability

**Single source of truth**:
- Playtest UI proven to work
- Backend mirrors playtest operations
- Frontend and backend stay in sync
- Easy to add new operations

**Clear separation of concerns**:
- Backend: State management + sync
- Frontend: UI + user input
- API: Simple command strings
- No complex business logic

## Testing Strategy

### Unit Tests

**Backend**:
- Test each operation in isolation
- Verify state updates correctly
- Check hidden information filtering
- Test rollback system

**Frontend**:
- Test store operations
- Verify WebSocket handling
- Check UI component rendering
- Test keyboard shortcuts

### Integration Tests

**2-Player Scenarios**:
1. Basic operations sync correctly
2. Hidden information properly filtered
3. Manual combat works
4. Rollback flow functions

**4-Player Scenarios**:
1. All players synchronized
2. Turn order cycles correctly
3. Each player sees 3 opponent views
4. Chat visible to all

### End-to-End Tests

**Real Game Flow**:
1. Create game
2. Multiple players join
3. Play cards, change life, tap creatures
4. Request rollback
5. Other players approve
6. State restored correctly
7. Game continues

## Future Enhancements

### Potential Additions

1. **Optimistic Updates**: Apply changes locally before server confirmation
2. **Action History**: Undo/redo individual actions (not just rollback)
3. **Spectator Mode**: Join games as observer
4. **Replay System**: Save and replay entire games
5. **Mobile UI**: Touch-optimized controls
6. **Voice Chat**: Integrated voice communication
7. **Advanced Bookmarks**: Named bookmarks, bookmark browser
8. **Automatic Bookmarks**: Configurable auto-bookmark triggers

### Performance Optimizations

1. **Delta Updates**: Send only changed state, not full state
2. **Binary Protocol**: Replace string commands with binary format
3. **State Compression**: Compress large game states
4. **Client Caching**: Cache card data, reduce bandwidth

### Rules Assistance (Optional)

**Non-blocking helper features**:
- Phase indicator (cosmetic only)
- Mana calculator (suggestion, not enforcement)
- Combat damage calculator (manual override)
- Card rules lookup (oracle text display)

**Important**: These would be **hints**, not **enforcement**. Players always have final control.

## Related Documentation

- `/docs/PLAYTEST_REPLACEMENT_PLAN.md` - Implementation plan
- `/docs/PLAYTEST_MIGRATION_SUMMARY.md` - Migration summary
- `/docs/tickets/done/003-implement-direct-engine.md` - Completed ticket
- `/mage-client-web/src/lib/stores/playtest-game.ts` - Store reference
- `/mage-client-web/src/routes/(protected)/playtest/+page.svelte` - UI reference

## Version History

- **v1.0** (2026-01-23): Initial architecture documentation
  - Playtest engine backend implemented
  - Frontend unified on playtest UI
  - Both engines coexist via config
  - Phase 7 cleanup completed
