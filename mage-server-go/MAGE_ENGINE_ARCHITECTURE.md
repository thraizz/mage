# Mage Engine Architecture Documentation

> **File**: `mage-server-go/internal/game/direct_engine.go`  
> **Purpose**: Rules-light Magic: The Gathering game engine - assists players without enforcing rules

## Table of Contents

1. [Overview](#overview)
2. [Design Philosophy](#design-philosophy)
3. [Core Types](#core-types)
4. [DirectEngine Structure](#directengine-structure)
5. [Game State Management](#game-state-management)
6. [Game Lifecycle](#game-lifecycle)
7. [Direct Action System](#direct-action-system)
8. [Zone Management](#zone-management)
9. [Combat Tracking](#combat-tracking)
10. [Bookmark & Rollback System](#bookmark--rollback-system)
11. [Action Log](#action-log)
12. [Persistence & Recovery](#persistence--recovery)
13. [Notification System](#notification-system)
14. [View Building](#view-building)

---

## Overview

The `DirectEngine` is a rules-light game engine that **assists** players rather than **enforces** rules. It provides:

- **State tracking**: Zones, life totals, counters, combat
- **Direct manipulation**: Players control all game state through UI
- **Action logging**: Every change is logged for review
- **Rollback support**: Any action can be undone
- **Real-time sync**: WebSocket notifications keep all clients updated

### What This Engine Does NOT Do

Unlike traditional MTG engines, we intentionally **do not**:

- Validate spell timing or mana costs
- Enforce priority or turn structure
- Check creature abilities (flying, menace, etc.)
- Auto-resolve state-based actions
- Process triggered abilities
- Validate targets
- Calculate combat damage automatically

Players are trusted to play correctly. The engine is a shared game board, not a referee.

---

## Design Philosophy

### Rules-Light Approach

Inspired by platforms like Untap.in, this engine follows the principle:

> **"Assist, don't enforce"**

| Traditional Engine | Rules-Light Engine |
|--------------------|-------------------|
| Validates every action | Logs every action |
| Rejects illegal plays | Allows any play, can rollback |
| Auto-resolves triggers | Players manually resolve |
| Enforces priority | Suggests phases, players control |
| Complex combat rules | Simple attacker/blocker tracking |

### Benefits

1. **Simplicity**: ~3,000 lines vs ~15,000 lines
2. **Flexibility**: House rules, casual play, testing
3. **Speed**: No validation overhead
4. **Reliability**: Fewer edge cases = fewer bugs
5. **Maintainability**: Easy to understand and modify

---

## Core Types

### Zone Constants

```go
const (
    ZoneLibrary     Zone = 0
    ZoneHand        Zone = 1
    ZoneBattlefield Zone = 2
    ZoneGraveyard   Zone = 3
    ZoneStack       Zone = 4  // Visual only - no resolution
    ZoneExile       Zone = 5
    ZoneCommand     Zone = 6
)
```

### Phase/Step Constants

```go
const (
    PhaseBeginning     Phase = "BEGINNING"
    PhasePrecombatMain Phase = "PRECOMBAT_MAIN"
    PhaseCombat        Phase = "COMBAT"
    PhasePostcombatMain Phase = "POSTCOMBAT_MAIN"
    PhaseEnding        Phase = "ENDING"
)

const (
    StepUntap         Step = "UNTAP"
    StepUpkeep        Step = "UPKEEP"
    StepDraw          Step = "DRAW"
    StepBeginCombat   Step = "BEGIN_COMBAT"
    StepDeclareAttackers Step = "DECLARE_ATTACKERS"
    StepDeclareBlockers  Step = "DECLARE_BLOCKERS"
    StepCombatDamage  Step = "COMBAT_DAMAGE"
    StepEndCombat     Step = "END_COMBAT"
    StepEnd           Step = "END"
    StepCleanup       Step = "CLEANUP"
)
```

---

## DirectEngine Structure

### Main Engine

```go
type DirectEngine struct {
    logger          *zap.Logger
    mu              sync.RWMutex
    games           map[string]*directGameState
    notifyHandler   atomic.Value  // NotificationHandler
    cardRepo        CardRepositoryInterface
    persistenceRepo PersistenceRepository
    
    // Rollback support
    snapshots       map[string][]*gameStateSnapshot
}
```

### Engine Initialization

```go
func NewDirectEngine(logger *zap.Logger) *DirectEngine
func (e *DirectEngine) SetNotificationHandler(handler NotificationHandler)
func (e *DirectEngine) SetCardRepository(repo CardRepositoryInterface)
func (e *DirectEngine) SetPersistenceRepository(repo PersistenceRepository)
```

---

## Game State Management

### directGameState

The internal representation of a single game:

```go
type directGameState struct {
    gameID      string
    gameType    string
    state       GameState  // ACTIVE, PAUSED, FINISHED
    
    // Players
    players     map[string]*internalPlayer
    playerOrder []string
    
    // Cards in all zones
    cards       map[string]*internalCard
    battlefield []*internalCard
    exile       []*internalCard
    command     []*internalCard
    
    // Turn tracking (suggested, not enforced)
    turnNumber    int
    activePlayer  string
    currentPhase  Phase
    currentStep   Step
    
    // Combat tracking (visual only)
    combat      *combatState
    
    // Action history
    actionLog   []ActionLogEntry
    nextActionID int
    
    // Messaging
    messages    []EngineMessage
    startedAt   time.Time
    
    mu sync.RWMutex
}
```

### internalCard

Represents a card in any zone:

```go
type internalCard struct {
    ID, Name, DisplayName   string
    ManaCost, Type          string
    SubTypes, SuperTypes    []string
    Color                   string
    Power, Toughness        string
    Loyalty                 string
    RulesText               string
    ExpansionSet, Rarity    string
    
    // State
    Zone         Zone
    Tapped       bool
    Flipped      bool
    Transformed  bool
    FaceDown     bool
    
    // Ownership
    ControllerID string
    OwnerID      string
    
    // Attachments
    AttachedTo   string    // Card this is attached to
    Attachments  []string  // Cards attached to this
    
    // Counters
    Counters     *counters.Counters
    
    // Combat (visual tracking only)
    Attacking      bool
    Blocking       bool
    AttackingWhat  string   // Defender ID
    BlockingWhat   []string // Attacker IDs
    
    // Token/Commander flags
    IsToken     bool
    IsCommander bool
    
    // Custom metadata
    Metadata    map[string]string
}
```

### internalPlayer

Represents a player in the game:

```go
type internalPlayer struct {
    PlayerID    string
    Name        string
    
    // Life and counters
    Life        int
    Poison      int
    Energy      int
    Experience  int
    
    // Zones (cards stored by reference)
    Library     []*internalCard
    Hand        []*internalCard
    Graveyard   []*internalCard
    
    // Mana pool (for display)
    ManaPool    *mana.ManaPool
    
    // Commander damage tracking
    CommanderDamage map[string]int
    
    // Game state
    Lost        bool
    Conceded    bool
}
```

---

## Game Lifecycle

### Starting a Game

```go
func (e *DirectEngine) StartGame(gameID string, players []string, gameType string) error
func (e *DirectEngine) StartGameWithDecks(gameID string, players []string, gameType string, decks map[string]DeckList) error
```

**Initialization Flow:**

1. Create `directGameState`
2. Load game type configuration (starting life, etc.)
3. Create players with starting life
4. Build decks from deck lists
5. Shuffle libraries
6. Set up command zone (for commanders)
7. Draw initial hands (7 cards)
8. Set state to `GameStateActive`
9. Determine starting player (random)
10. Log game start action

### Ending a Game

```go
func (e *DirectEngine) EndGame(gameID string, winner string) error
func (e *DirectEngine) PlayerConcede(gameID, playerID string) error
func (e *DirectEngine) CleanupGame(gameID string) error
```

---

## Direct Action System

All game state changes happen through direct actions. Each action:

1. Creates a snapshot (for rollback)
2. Modifies state
3. Logs the action
4. Emits notification

### Player State Actions

```go
// Life total
func (e *DirectEngine) SetLife(gameID, playerID string, life int) error
func (e *DirectEngine) ModifyLife(gameID, playerID string, delta int) error

// Player counters
func (e *DirectEngine) SetPlayerCounter(gameID, playerID, counterType string, amount int) error

// Mana pool (display only)
func (e *DirectEngine) AddMana(gameID, playerID string, manaType string, amount int) error
func (e *DirectEngine) ClearManaPool(gameID, playerID string) error
```

### Card Actions

```go
// Zone movement
func (e *DirectEngine) MoveCard(gameID, cardID string, targetZone Zone, controllerID string) error

// Card state
func (e *DirectEngine) TapUntap(gameID, cardID string, tapped bool) error
func (e *DirectEngine) FlipCard(gameID, cardID string) error
func (e *DirectEngine) TransformCard(gameID, cardID string) error

// Counters
func (e *DirectEngine) SetCardCounter(gameID, cardID, counterType string, amount int) error

// Attachments
func (e *DirectEngine) AttachCard(gameID, cardID, targetID string) error
func (e *DirectEngine) DetachCard(gameID, cardID string) error
```

### Library Actions

```go
func (e *DirectEngine) DrawCards(gameID, playerID string, count int) error
func (e *DirectEngine) ShuffleLibrary(gameID, playerID string) error
func (e *DirectEngine) MillCards(gameID, playerID string, count int) error
func (e *DirectEngine) RevealCards(gameID string, cardIDs []string, toPlayerIDs []string) error
func (e *DirectEngine) LookAtCards(gameID, playerID string, cardIDs []string) error
```

### Token Actions

```go
type TokenDefinition struct {
    Name       string
    Types      []string
    SubTypes   []string
    Power      string
    Toughness  string
    Colors     []string
    Abilities  []string
}

func (e *DirectEngine) CreateToken(gameID, controllerID string, token TokenDefinition, count int) ([]string, error)
func (e *DirectEngine) DestroyToken(gameID, tokenID string) error
```

### Turn/Phase Actions

```go
func (e *DirectEngine) SetPhase(gameID string, phase Phase, step Step) error
func (e *DirectEngine) NextTurn(gameID string) error
func (e *DirectEngine) SetActivePlayer(gameID, playerID string) error
func (e *DirectEngine) UntapAll(gameID, playerID string) error
```

### Combat Actions

```go
func (e *DirectEngine) DeclareAttacker(gameID, creatureID, defenderID string) error
func (e *DirectEngine) DeclareBlocker(gameID, blockerID, attackerID string) error
func (e *DirectEngine) ClearCombat(gameID string) error
```

---

## Zone Management

### Moving Cards

All zone transitions use `MoveCard`:

```go
func (e *DirectEngine) MoveCard(gameID, cardID string, targetZone Zone, controllerID string) error
```

**Examples:**

- Play a land: `MoveCard(gameID, cardID, ZoneBattlefield, playerID)`
- Cast a spell: `MoveCard(gameID, cardID, ZoneStack, playerID)` then `MoveCard(gameID, cardID, ZoneBattlefield, playerID)`
- Discard: `MoveCard(gameID, cardID, ZoneGraveyard, playerID)`
- Exile: `MoveCard(gameID, cardID, ZoneExile, playerID)`

### Zone Queries

```go
func (e *DirectEngine) GetCardsInZone(gameID string, zone Zone, playerID string) ([]*internalCard, error)
func (e *DirectEngine) GetCardByID(gameID, cardID string) (*internalCard, error)
```

---

## Combat Tracking

Combat is **visual only** - no rules enforcement.

### combatState

```go
type combatState struct {
    attackingPlayer string
    attackers       map[string]string  // creatureID -> defenderID
    blockers        map[string][]string // attackerID -> []blockerIDs
}
```

### Combat Flow (Suggested)

1. Player declares attackers via UI → `DeclareAttacker`
2. Opponent declares blockers via UI → `DeclareBlocker`
3. Players manually apply damage via `ModifyLife` / counter changes
4. After combat → `ClearCombat`

The engine tracks who is attacking/blocking for visual display only.

---

## Bookmark & Rollback System

### gameStateSnapshot

```go
type gameStateSnapshot struct {
    ID           int
    GameID       string
    Timestamp    time.Time
    
    // Complete state copy
    Players      map[string]*internalPlayer
    Cards        map[string]*internalCard
    Battlefield  []*internalCard
    Exile        []*internalCard
    Command      []*internalCard
    
    TurnNumber   int
    ActivePlayer string
    Phase        Phase
    Step         Step
    Combat       *combatState
}
```

### Snapshot Methods

```go
// Create snapshot before action
func (e *DirectEngine) createSnapshot(gameState *directGameState) *gameStateSnapshot

// Restore to snapshot
func (e *DirectEngine) restoreSnapshot(gameState *directGameState, snapshot *gameStateSnapshot) error
```

---

## Action Log

Every action is logged for review and rollback.

### ActionLogEntry

```go
type ActionLogEntry struct {
    ID          int
    Timestamp   time.Time
    PlayerID    string
    ActionType  string    // "MOVE_CARD", "SET_LIFE", "CREATE_TOKEN", etc.
    Description string    // Human-readable description
    SnapshotID  int       // Links to snapshot for rollback
    Data        map[string]interface{}  // Action-specific data
}
```

### Action Log Methods

```go
// Get action history
func (e *DirectEngine) GetActionLog(gameID string, limit, offset int) ([]ActionLogEntry, error)

// Rollback to specific action
func (e *DirectEngine) RollbackToAction(gameID string, actionID int) error
```

### Rollback Flow

1. Find action in log by ID
2. Get associated snapshot
3. Restore game state from snapshot
4. Truncate action log to that point
5. Notify all clients of rollback

---

## Persistence & Recovery

### Interface

```go
type PersistenceRepository interface {
    SaveGameState(ctx context.Context, gameID, tableID, gameType string, 
                  players []string, gameState []byte, turnNumber int, state string) error
    LoadGameState(ctx context.Context, gameID string) ([]byte, error)
    DeleteActiveGame(ctx context.Context, gameID string) error
}
```

### Methods

```go
func (e *DirectEngine) PersistGameState(gameID string) error
func (e *DirectEngine) LoadGameFromPersistence(gameID string) error
func (e *DirectEngine) DeletePersistedGame(gameID string) error
```

### Auto-Save

Game state is automatically persisted:

- After every action (debounced)
- On game pause
- On server shutdown

---

## Notification System

### GameNotification

```go
type GameNotification struct {
    Type      string                 // Event type
    GameID    string
    PlayerID  string                 // Empty for broadcast
    Timestamp time.Time
    Data      map[string]interface{}
}
```

### Notification Types

| Type | Trigger |
|------|---------|
| `GAME_STATE_CHANGE` | Any state modification |
| `CARD_MOVED` | Card changes zones |
| `LIFE_CHANGED` | Life total modified |
| `COUNTER_CHANGED` | Counter added/removed |
| `COMBAT_UPDATED` | Attackers/blockers declared |
| `TURN_CHANGED` | Turn/phase advanced |
| `ACTION_LOGGED` | New action in log |
| `ROLLBACK` | State rolled back |
| `GAME_ENDED` | Game finished |

### Notification Methods

```go
func (e *DirectEngine) emitNotification(notification GameNotification)
func (e *DirectEngine) notifyGameStateChange(gameID string, data map[string]interface{})
func (e *DirectEngine) notifyCardMoved(gameID string, cardID string, fromZone, toZone Zone)
func (e *DirectEngine) notifyLifeChanged(gameID, playerID string, oldLife, newLife int)
```

---

## View Building

### EngineGameView

The complete game state sent to clients:

```go
type EngineGameView struct {
    GameID       string
    GameType     string
    State        string
    TurnNumber   int
    ActivePlayer string
    CurrentPhase string
    CurrentStep  string
    
    Players     []EnginePlayerView
    Battlefield []EngineCardView
    Stack       []EngineCardView  // Visual only
    Exile       []EngineCardView
    Command     []EngineCardView
    Combat      *EngineCombatView
    
    ActionLog   []ActionLogEntry  // Recent actions
}
```

### View Methods

```go
func (e *DirectEngine) GetGameView(gameID, playerID string) (*EngineGameView, error)
func (e *DirectEngine) buildPlayerView(player *internalPlayer, isViewer bool) EnginePlayerView
func (e *DirectEngine) buildCardView(card *internalCard) EngineCardView
func (e *DirectEngine) buildCombatView(combat *combatState) *EngineCombatView
```

### Hidden Information

- Other players' hands are hidden (show card backs)
- Other players' libraries are hidden
- Face-down cards show only card back
- Revealed cards are visible to specified players

---

## Thread Safety

The engine uses a two-level locking strategy:

1. **Engine lock (`e.mu`)**: Protects the `games` map
2. **Game state lock (`gameState.mu`)**: Protects individual game state

**Lock ordering** (to prevent deadlocks):

1. Acquire `e.mu` (RLock)
2. Get game state reference
3. Release `e.mu`
4. Acquire `gameState.mu` (Lock)
5. Perform operations
6. Release `gameState.mu`

---

## API Integration

### gRPC Service

The `DirectGameService` exposes all direct actions:

```protobuf
service DirectGameService {
    // Card actions
    rpc MoveCard(MoveCardRequest) returns (MoveCardResponse);
    rpc TapUntap(TapUntapRequest) returns (TapUntapResponse);
    rpc SetCardCounter(SetCardCounterRequest) returns (SetCardCounterResponse);
    
    // Player actions
    rpc SetLife(SetLifeRequest) returns (SetLifeResponse);
    rpc ModifyLife(ModifyLifeRequest) returns (ModifyLifeResponse);
    rpc DrawCards(DrawCardsRequest) returns (DrawCardsResponse);
    
    // Turn actions
    rpc SetPhase(SetPhaseRequest) returns (SetPhaseResponse);
    rpc NextTurn(NextTurnRequest) returns (NextTurnResponse);
    
    // Combat
    rpc DeclareAttacker(DeclareAttackerRequest) returns (DeclareAttackerResponse);
    rpc DeclareBlocker(DeclareBlockerRequest) returns (DeclareBlockerResponse);
    rpc ClearCombat(ClearCombatRequest) returns (ClearCombatResponse);
    
    // Tokens
    rpc CreateToken(CreateTokenRequest) returns (CreateTokenResponse);
    
    // Rollback
    rpc GetActionLog(GetActionLogRequest) returns (GetActionLogResponse);
    rpc RollbackToAction(RollbackToActionRequest) returns (RollbackToActionResponse);
}
```

### WebSocket Events

Real-time updates are pushed via WebSocket:

```json
{
    "type": "GAME_STATE_CHANGE",
    "gameId": "game-123",
    "data": {
        "turnNumber": 5,
        "activePlayer": "player-1",
        "phase": "PRECOMBAT_MAIN"
    }
}
```
