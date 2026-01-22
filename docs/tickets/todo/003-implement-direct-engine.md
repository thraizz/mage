# Replace Game View with Playtest-Based Multiplayer Engine

## Problem

The current game view has two major issues:

1. **Backend**: `mage_engine.go` (13,787 lines) implements full MTG rules enforcement
   - Validates spell timing, mana costs, priority
   - Enforces complex rules (stack, triggers, state-based actions)
   - Auto-resolves combat damage and effects
   - Over-engineered for our use case

2. **Frontend**: Game view (`game/[id]/+page.svelte`) duplicates playtest functionality
   - Re-implements zones, card management, UI patterns
   - Adds rules-enforcement UI (priority bar, combat dialogs, mana prompts)
   - Creates unnecessary complexity and maintenance burden

Meanwhile, **playtest mode works perfectly** as a client-side rules-light engine:
- Clean state management (1,257 lines)
- Proven UI patterns (1,992 lines)
- All game operations implemented (tap, move, counters, life, etc.)
- Excellent UX (keyboard shortcuts, drag-drop, context menus)

## Solution

**Use playtest as the foundation, adapt for multiplayer**:

1. **Backend**: Replace MageEngine with rules-light engine matching playtest operations
2. **Frontend**: Replace game view with playtest UI + multiplayer components
3. **Integration**: Add minimal server sync layer for state synchronization

**No rules enforcement, no backward compatibility** - clean replacement.

---

## Architecture

### Backend: Rules-Light Engine

**Copy state structure from playtest** (`playtest-game.ts` lines 25-59):

```go
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

**Implement all playtest operations** (`playtest-game.ts` lines 492-1151):

```go
func (e *Engine) MoveCard(gameID, cardID, targetZone string) error
func (e *Engine) TapCard(gameID, cardID string, tapped bool) error
func (e *Engine) DrawCards(gameID, playerID string, count int) error
func (e *Engine) ModifyLife(gameID, playerID string, delta int) error
func (e *Engine) AddCounter(gameID, cardID, counterType string, amount int) error
func (e *Engine) CreateToken(gameID, name, types, power, toughness, color string) error
func (e *Engine) ShuffleLibrary(gameID, playerID string) error
func (e *Engine) ScryCards(gameID, playerID string, count int) ([]Card, error)
func (e *Engine) MillCards(gameID, playerID string, count int) error
func (e *Engine) NextTurn(gameID string) error
```

**Use existing direct-actions API** (already implemented in `mage_engine.go` lines 2590-2740):
- Backend already parses string commands: `TAP:{cardId}`, `MOVE:{cardId}:{zone}`, etc.
- Frontend already uses `direct-actions.ts` to send these commands
- No proto changes needed initially (can migrate to typed RPCs later)

### Frontend: Playtest UI + Multiplayer

**Copy playtest page structure** (`playtest/+page.svelte` lines 1-1992):
- State management patterns
- Event handlers (lines 440-760)
- Keyboard shortcuts (lines 955-1050)
- Drag-drop system (lines 196-780)
- Context menus (lines 499-540)

**Add multiplayer components** from game view:
- `OpponentSection.svelte` - Opponent battlefield/info
- `PlayerInfoRow.svelte` - Multi-player life totals
- `GameChatOverlay.svelte` - Chat

**Remove rules-enforcement UI**:
- Delete `PriorityActionBar`, `DeclareAttackers`, `DeclareBlockers`
- Delete `ManaPayment`, `XManaSelector`, `AbilitiesPanel`
- Delete all `CallbackMethod.*` event handlers

### State Synchronization

**Create multiplayer store** (`multiplayer-game.ts`):

```typescript
// Copy playtest store structure
export interface MultiplayerGameState {
  gameId: string;
  players: PlaytestPlayer[];  // Same as playtest
  battlefield: CardView[];
  exile: CardView[];
  // ... all zones from playtest

  // NEW: Server sync
  isConnected: boolean;
  pendingActions: string[];
}

// Copy all operations, send to server
function drawCards(playerId: string, count: number): void {
  directActions.drawCards(gameId, playerId, count);
  // Server broadcasts update back
}
```

**WebSocket sync pattern**:

```typescript
websocketStore.subscribe((message) => {
  if (message.type === 'GAME_UPDATE') {
    multiplayerGameStore.applyServerState(message.gameState);
  }
});
```

---

## Implementation Plan

### Phase 1: Backend - Rules-Light Engine

**Files to Create**:
- `mage-server-go/internal/game/engine.go` (~500-700 lines)
  - Engine struct, state management
  - Copy state structure from playtest-game.ts lines 25-59

- `mage-server-go/internal/game/actions.go` (~400-600 lines)
  - Implement all operations from playtest-game.ts lines 492-1151
  - Translate TypeScript logic to Go

- `mage-server-go/internal/game/rollback.go` (~300-400 lines)
  - Bookmark/restore system
  - State snapshots at turn boundaries

- `mage-server-go/internal/game/view.go` (~200-300 lines)
  - Hidden information filtering (opponent hands/libraries)
  - Public zone visibility (battlefield, graveyard, exile)

**Files to Modify**:
- `mage-server-go/internal/game/manager.go`
  - Replace MageEngine with Engine
  - Keep GameEngine interface (lines 319-350)

**Files to Delete**:
- `mage-server-go/internal/game/mage_engine.go` (13,786 lines)
- `mage-server-go/internal/game/engine_*.go` (all engine files)

**Reference Patterns**:
- State structure: `playtest-game.ts` lines 25-59
- Operations: `playtest-game.ts` lines 492-1151
- Helpers: `playtest-helpers.ts` lines 30-243
- Existing direct-action parsing: `mage_engine.go` lines 2590-2740

### Phase 2: Frontend - Multiplayer Store

**Files to Create**:
- `mage-client-web/src/lib/stores/multiplayer-game.ts`
  - Copy entire `playtest-game.ts` structure (lines 1-1257)
  - Wire operations to `direct-actions.ts` API
  - Add WebSocket subscription for server updates
  - Remove localStorage persistence (server is source of truth)

**Reference Patterns**:
- Store structure: `playtest-game.ts` lines 380-401
- Operations: `playtest-game.ts` lines 492-1151
- WebSocket subscription: `game.ts` lines 140-200
- Direct actions API: `direct-actions.ts` lines 1-264

### Phase 3: Frontend - Replace Game Page

**Files to Modify**:
- `mage-client-web/src/routes/(protected)/game/[id]/+page.svelte`
  - Replace with copy of `playtest/+page.svelte` (lines 1-1992)
  - Use `multiplayerGameStore` instead of `playtestGameStore`
  - Add opponent views: `OpponentSection`, `PlayerInfoRow`
  - Add chat: `GameChatOverlay`
  - Keep all playtest UI patterns (keyboard, drag-drop, context menus)

**Components to Keep**:
- Copy from playtest: All core components (Card, PlayerHand, BattlefieldArea, etc.)
- Copy from game view: Multiplayer components (OpponentSection, PlayerInfoRow, GameChatOverlay)

**Components to Delete**:
- `PriorityActionBar.svelte`
- `DeclareAttackers.svelte`
- `DeclareBlockers.svelte`
- `AssignDamage.svelte`
- `ManaPayment.svelte`
- `XManaSelector.svelte`
- `AbilitiesPanel.svelte`
- `combat.ts` store
- `combat.ts` types

**Reference Patterns**:
- Page structure: `playtest/+page.svelte` lines 1-1992
- Event handlers: `playtest/+page.svelte` lines 440-760
- Keyboard shortcuts: `playtest/+page.svelte` lines 955-1050
- Drag-drop: `playtest/+page.svelte` lines 196-780
- Opponent views: `game/[id]/+page.svelte` lines 1100-1400

### Phase 4: Integration Testing

**Test Scenarios**:

1. **Basic Operations (2 Players)**:
   - Draw card → Both see update
   - Play card → Both see battlefield
   - Tap card → Both see animation
   - Change life → Both see life total
   - Move to graveyard → Both see zone update

2. **Hidden Information**:
   - Own hand → See all cards
   - Opponent hand → See count only
   - Play from hand → Opponent sees card revealed
   - Library → Count visible, cards hidden

3. **Combat (Manual)**:
   - Declare attackers (tap creatures, no validation)
   - Declare blockers (tap creatures, no validation)
   - Apply damage manually (life changes, destroy creatures)
   - Clear combat state

4. **Rollback**:
   - Bookmark state
   - Request rollback
   - Consent dialog
   - State restored

5. **Multiplayer (4 Players)**:
   - All players synchronized
   - Turn order cycles
   - 3 opponent views per player
   - Chat visible to all

### Phase 5: Cleanup

**Delete**:
- Old backend engine files
- Old frontend rules-enforcement components
- Old game store (`game.ts` → rename to `game.legacy.ts`)

**Update Documentation**:
- `docs/GAME_ARCHITECTURE.md` - Document new architecture
- `README.md` - Update feature list (rules-light, Untap.in style)

---

## File Reference

### Backend Files

**To Create**:
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/engine.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/state.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/actions.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/rollback.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/view.go`

**To Modify**:
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/manager.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/server/grpc_game.go`

**To Delete**:
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/mage_engine.go`
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/engine_*.go`

### Frontend Files

**To Create**:
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/multiplayer-game.ts`

**To Modify**:
- `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/game/[id]/+page.svelte`

**To Keep (Copy to New Page)**:
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/Card.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/PlayerHand.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/BattlefieldArea.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/OpponentSection.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/PlayerInfoRow.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/GameChatOverlay.svelte`

**To Delete**:
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/PriorityActionBar.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/DeclareAttackers.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/DeclareBlockers.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/AssignDamage.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/ManaPayment.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/XManaSelector.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/components/game/AbilitiesPanel.svelte`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/combat.ts`
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/types/combat.ts`

### Reference Files (Copy Patterns From)

**Playtest Engine** (Copy state & operations):
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/playtest-game.ts` (1,257 lines)
  - Lines 25-59: State interface
  - Lines 380-401: Store creation
  - Lines 492-1151: All operations

**Playtest Helpers** (Copy state helpers):
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/utils/playtest-helpers.ts` (303 lines)
  - Lines 30-63: Find card in state
  - Lines 120-170: Remove from zone
  - Lines 175-243: Add to zone

**Playtest UI** (Copy event handlers & layout):
- `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/playtest/+page.svelte` (1,992 lines)
  - Lines 440-760: Event handlers
  - Lines 955-1050: Keyboard shortcuts
  - Lines 196-780: Drag-drop system

**Direct Actions API** (Use existing):
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/api/direct-actions.ts` (264 lines)

**Backend Direct Actions** (Reference existing implementation):
- `/Users/aron/dev/opensource/mage/mage-server-go/internal/game/mage_engine.go` (lines 2590-2740)

---

## What It Does NOT Do

- ❌ Validate spell timing or mana costs
- ❌ Enforce priority or turn structure
- ❌ Check creature abilities (flying, menace, etc.)
- ❌ Auto-resolve state-based actions
- ❌ Process triggered abilities
- ❌ Validate targets
- ❌ Calculate combat damage automatically
- ❌ Show priority/phase UI
- ❌ Show mana payment prompts
- ❌ Show ability selection prompts

---

## Benefits

1. **Simplicity**:
   - Backend: 13,786 lines → ~2,500 lines (82% reduction)
   - Frontend: Unified codebase (no playtest/game duplication)

2. **Maintainability**:
   - Single UI pattern (playtest-based)
   - Simple state management
   - No rules engine complexity

3. **Proven UX**:
   - Playtest UI already works perfectly
   - Keyboard shortcuts, drag-drop, context menus
   - Fast, responsive, intuitive

4. **Flexibility**:
   - Supports house rules and casual play
   - Players control game state
   - Rollback for mistakes

5. **Speed**:
   - No validation overhead
   - Direct state manipulation
   - Real-time WebSocket sync

---

## Migration Strategy

**Full Replacement** - No backward compatibility:

- Delete old backend engine (13,786 lines)
- Delete old frontend rules UI
- Replace with playtest-based multiplayer
- Existing games incompatible (expected, single user)
- Clean, simple codebase

**Deployment**:
1. Deploy backend with new engine
2. Deploy frontend with new game page
3. Create new games (old games won't work)

**Risk Mitigation**:
- Backend already supports direct-actions API
- Playtest UI proven to work
- Only adding server sync layer
- Test locally before production

---

## Success Criteria

### Backend
- ✅ Engine implements all playtest operations
- ✅ State syncs via WebSocket
- ✅ Rollback system works
- ✅ No rules enforcement
- ✅ Codebase reduced to ~2,500 lines

### Frontend
- ✅ Game page uses playtest UI
- ✅ All playtest features work (drag-drop, shortcuts, menus)
- ✅ Multiplayer components integrated
- ✅ No priority/prompt UI
- ✅ State syncs correctly

### Integration
- ✅ 2-4 player games work
- ✅ Hidden information enforced
- ✅ Manual combat works
- ✅ Rollback consent works
- ✅ No state desync

---

## Related Documentation

- `/Users/aron/dev/opensource/mage/docs/PLAYTEST_REPLACEMENT_PLAN.md` - Detailed implementation plan
- `/Users/aron/dev/opensource/mage/mage-client-web/src/lib/stores/playtest-game.ts` - Reference implementation
- `/Users/aron/dev/opensource/mage/mage-client-web/src/routes/(protected)/playtest/+page.svelte` - Reference UI

---

## Priority

**High** - Core architectural change that simplifies codebase and unifies playtest/game modes into single multiplayer rules-light engine.
