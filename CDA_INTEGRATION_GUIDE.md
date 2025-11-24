# CDA Integration Guide

**Status**: CDA system implemented, integration pending
**Files Created**: `/internal/game/abilities/characteristic_defining.go` (450 lines)

## What Was Implemented

### CDA System Components ✓

1. **Base Interface** - `CharacteristicDefiningAbility`
   - Methods: `DefinesPower()`, `DefinesToughness()`, `DefinesColor()`, `DefinesTypes()`
   - Calculation methods for each characteristic
   - Inherits from `Ability` interface

2. **Five CDA Implementations**:
   - `TarmogoyfCDA` - Counts card types in all graveyards
   - `LordOfExtinctionCDA` - Counts all cards in graveyards
   - `CreaturesYouControlCDA` - Counts creatures you control
   - `HandSizeCDA` - Counts cards in your hand (Maro, etc.)
   - `CountersCDA` - P/T equals counters on permanent
   - `CalculationCDA` - Generic with custom calculation functions

3. **GameContext Extensions** ✓
   - Added to `/internal/game/abilities/ability.go`:
   ```go
   GetAllCardsInZone(ctx context.Context, zone int) []CardInfo
   GetCreaturesControlledBy(ctx context.Context, playerID uuid.UUID) []CardInfo
   GetPlayerHand(ctx context.Context, playerID uuid.UUID) []CardInfo
   GetCountersOnPermanent(ctx context.Context, permanentID uuid.UUID, counterType string) int
   ```

4. **Supporting Types** ✓
   - `CardInfo` interface for minimal card data
   - Zone constants (ZoneGraveyard, ZoneBattlefield, etc.)

## Integration Steps Required

### Step 1: Modify getCreaturePower and getCreatureToughness

**Current Code** (`mage_engine.go` lines 6867-6902):
```go
func (e *MageEngine) getCreaturePower(creature *internalCard) (int, error) {
    if creature.Power == "*" || creature.Power == "X" {
        return 0, nil // TODO: Calculate dynamic power
    }
    power, err := strconv.Atoi(creature.Power)
    return power, err
}
```

**Required Changes**:

#### Option A: Add gameID parameter (Breaking Change)
```go
func (e *MageEngine) getCreaturePower(gameID string, creature *internalCard) (int, error) {
    if creature.Power == "*" || creature.Power == "X" {
        return e.calculateCDAPower(gameID, creature)
    }
    power, err := strconv.Atoi(creature.Power)
    return power, err
}

func (e *MageEngine) calculateCDAPower(gameID string, creature *internalCard) (int, error) {
    // Find CDA among creature's abilities
    for _, abilityID := range creature.Abilities {
        ability := e.getAbilityByID(gameID, abilityID)
        if cda, ok := ability.(abilities.CharacteristicDefiningAbility); ok {
            if cda.DefinesPower() {
                ctx := context.Background()
                return cda.CalculatePower(ctx, e.getGameContext(gameID))
            }
        }
    }
    return 0, nil // No CDA found, default to 0
}
```

#### Option B: Store gameID in internalCard (Non-Breaking)
```go
type internalCard struct {
    // ... existing fields ...
    GameID string // Add this field
}

func (e *MageEngine) getCreaturePower(creature *internalCard) (int, error) {
    if creature.Power == "*" || creature.Power == "X" {
        // Use gameID from creature
        return e.calculateCDAPower(creature.GameID, creature)
    }
    // ... rest unchanged ...
}
```

**Recommendation**: Use Option B to avoid breaking existing call sites.

### Step 2: Implement getAbilityByID

**Add to MageEngine**:
```go
// getAbilityByID retrieves an ability by its ID
func (e *MageEngine) getAbilityByID(gameID string, abilityID uuid.UUID) abilities.Ability {
    e.mu.RLock()
    defer e.mu.RUnlock()

    gameState, exists := e.games[gameID]
    if !exists {
        return nil
    }

    // Check ability registry
    if e.abilityRegistry != nil {
        if metadata := e.abilityRegistry.GetAbility(abilityID); metadata != nil {
            return metadata.Ability
        }
    }

    return nil
}
```

### Step 3: Implement GameContext Methods

**Add to MageEngine** (implement the CDA query methods):

```go
// GetAllCardsInZone returns all cards in a specific zone
func (e *MageEngine) GetAllCardsInZone(ctx context.Context, zone int) []abilities.CardInfo {
    // Extract gameID from context or use current game
    gameID := extractGameIDFromContext(ctx)

    e.mu.RLock()
    defer e.mu.RUnlock()

    gameState, exists := e.games[gameID]
    if !exists {
        return []abilities.CardInfo{}
    }

    result := []abilities.CardInfo{}

    switch zone {
    case abilities.ZoneGraveyard:
        for _, card := range gameState.graveyard {
            result = append(result, newCardInfoAdapter(card))
        }
    case abilities.ZoneBattlefield:
        for _, card := range gameState.battlefield {
            result = append(result, newCardInfoAdapter(card))
        }
    case abilities.ZoneHand:
        for _, player := range gameState.players {
            for _, card := range player.hand {
                result = append(result, newCardInfoAdapter(card))
            }
        }
    // ... other zones ...
    }

    return result
}

// GetCreaturesControlledBy returns all creatures controlled by a player
func (e *MageEngine) GetCreaturesControlledBy(ctx context.Context, playerID uuid.UUID) []abilities.CardInfo {
    gameID := extractGameIDFromContext(ctx)

    e.mu.RLock()
    defer e.mu.RUnlock()

    gameState, exists := e.games[gameID]
    if !exists {
        return []abilities.CardInfo{}
    }

    result := []abilities.CardInfo{}
    for _, card := range gameState.battlefield {
        if card.ControllerID == playerID.String() && card.isCreature() {
            result = append(result, newCardInfoAdapter(card))
        }
    }

    return result
}

// GetPlayerHand returns cards in a player's hand
func (e *MageEngine) GetPlayerHand(ctx context.Context, playerID uuid.UUID) []abilities.CardInfo {
    gameID := extractGameIDFromContext(ctx)

    e.mu.RLock()
    defer e.mu.RUnlock()

    gameState, exists := e.games[gameID]
    if !exists {
        return []abilities.CardInfo{}
    }

    player, exists := gameState.players[playerID.String()]
    if !exists {
        return []abilities.CardInfo{}
    }

    result := make([]abilities.CardInfo, 0, len(player.hand))
    for _, card := range player.hand {
        result = append(result, newCardInfoAdapter(card))
    }

    return result
}

// GetCountersOnPermanent returns the number of a specific counter type
func (e *MageEngine) GetCountersOnPermanent(ctx context.Context, permanentID uuid.UUID, counterType string) int {
    gameID := extractGameIDFromContext(ctx)

    e.mu.RLock()
    defer e.mu.RUnlock()

    gameState, exists := e.games[gameID]
    if !exists {
        return 0
    }

    card, exists := gameState.battlefield[permanentID.String()]
    if !exists {
        return 0
    }

    if card.Counters == nil {
        return 0
    }

    return card.Counters[counterType]
}
```

### Step 4: Create CardInfo Adapter

**Add to mage_engine.go**:
```go
// cardInfoAdapter adapts internalCard to abilities.CardInfo interface
type cardInfoAdapter struct {
    card *internalCard
}

func newCardInfoAdapter(card *internalCard) abilities.CardInfo {
    return &cardInfoAdapter{card: card}
}

func (c *cardInfoAdapter) GetID() uuid.UUID {
    id, _ := uuid.Parse(c.card.ID)
    return id
}

func (c *cardInfoAdapter) GetName() string {
    return c.card.Name
}

func (c *cardInfoAdapter) GetTypes() []string {
    // Parse CardType string into slice
    // Format is like "Creature — Human Warrior"
    types := []string{}
    parts := strings.Split(c.card.CardType, "—")
    if len(parts) > 0 {
        typePart := strings.TrimSpace(parts[0])
        types = strings.Fields(typePart)
    }
    return types
}

func (c *cardInfoAdapter) GetSubtypes() []string {
    // Parse subtypes from CardType string
    subtypes := []string{}
    parts := strings.Split(c.card.CardType, "—")
    if len(parts) > 1 {
        subtypePart := strings.TrimSpace(parts[1])
        subtypes = strings.Fields(subtypePart)
    }
    return subtypes
}

func (c *cardInfoAdapter) GetPower() int {
    power, _ := strconv.Atoi(c.card.Power)
    return power
}

func (c *cardInfoAdapter) GetToughness() int {
    toughness, _ := strconv.Atoi(c.card.Toughness)
    return toughness
}
```

### Step 5: Context Management

**Add helper to extract gameID from context**:
```go
// contextKey is a custom type for context keys
type contextKey string

const gameIDContextKey contextKey = "gameID"

// extractGameIDFromContext gets the game ID from context
func extractGameIDFromContext(ctx context.Context) string {
    if gameID, ok := ctx.Value(gameIDContextKey).(string); ok {
        return gameID
    }
    // Fallback: return empty string and let caller handle
    return ""
}

// withGameID adds game ID to context
func withGameID(ctx context.Context, gameID string) context.Context {
    return context.WithValue(ctx, gameIDContextKey, gameID)
}
```

### Step 6: Update internalCard Structure

**Add GameID field**:
```go
type internalCard struct {
    // ... existing fields ...
    GameID    string      // Add this to track which game this card belongs to
}
```

**Update card creation** to set GameID:
```go
// Whenever creating/adding a card to a game, set its GameID
card.GameID = gameID
```

## Usage Examples

### Example 1: Tarmogoyf Card Implementation

```go
package generated

import (
    "github.com/google/uuid"
    "github.com/magefree/mage-server-go/internal/game"
    "github.com/magefree/mage-server-go/internal/game/abilities"
    "github.com/magefree/mage-server-go/internal/game/cards"
)

func NewTarmogoyf(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
    card := game.NewCard(ownerID, "Tarmogoyf")
    card.ManaCost = "{1}{G}"
    card.Types = []string{"CREATURE"}
    card.Subtypes = []string{"Lhurgoyf"}
    card.Power = "*"       // Indicates dynamic power
    card.Toughness = "1+*" // Indicates dynamic toughness

    // Add Tarmogoyf's CDA
    cda := abilities.NewTarmogoyfCDA(card.ID)
    card.AddAbility(cda)

    return card, nil
}
```

### Example 2: Lord of Extinction

```go
func NewLordOfExtinction(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
    card := game.NewCard(ownerID, "Lord of Extinction")
    card.ManaCost = "{3}{B}{G}"
    card.Types = []string{"CREATURE"}
    card.Subtypes = []string{"Elemental"}
    card.Power = "*"
    card.Toughness = "*"

    // Add Lord of Extinction's CDA
    cda := abilities.NewLordOfExtinctionCDA(card.ID)
    card.AddAbility(cda)

    return card, nil
}
```

### Example 3: Custom CDA with CalculationCDA

```go
// For a card like "Multani, Yavimaya's Avatar"
// Power/toughness = lands you control + cards in your hand

func NewMultani(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
    card := game.NewCard(ownerID, "Multani, Yavimaya's Avatar")
    card.ManaCost = "{4}{G}{G}"
    card.Types = []string{"LEGENDARY", "CREATURE"}
    card.Subtypes = []string{"Elemental", "Avatar"}
    card.Power = "*"
    card.Toughness = "*"

    powerCalc := func(ctx context.Context, game abilities.GameContext) (int, error) {
        // Count lands controlled
        lands := 0
        creatures := game.GetCreaturesControlledBy(ctx, ownerID)
        for _, creature := range creatures {
            types := creature.GetTypes()
            for _, t := range types {
                if strings.EqualFold(t, "land") {
                    lands++
                }
            }
        }

        // Count cards in hand
        hand := game.GetPlayerHand(ctx, ownerID)
        return lands + len(hand), nil
    }

    cda := abilities.NewCalculationCDA(
        card.ID,
        true, true, false, false, // defines power and toughness
        "Multani's power and toughness are each equal to the total number of lands you control and cards in your hand.",
        powerCalc, powerCalc, nil, nil,
    )

    card.AddAbility(cda)
    return card, nil
}
```

## Testing Plan

### Unit Tests

**File**: `/internal/game/abilities/characteristic_defining_test.go`

```go
func TestTarmogoyfCDA(t *testing.T) {
    // Setup mock game context with cards in graveyards
    // Test that Tarmogoyf's P/T changes with graveyard composition
}

func TestLordOfExtinctionCDA(t *testing.T) {
    // Test that Lord of Extinction's P/T equals card count in graveyards
}

func TestCreaturesYouControlCDA(t *testing.T) {
    // Test that P/T changes as creatures enter/leave battlefield
}
```

### Integration Tests

**File**: `/internal/game/cda_integration_test.go`

```go
func TestTarmogoyfGameplay(t *testing.T) {
    // Full game test:
    // 1. Cast Tarmogoyf (should be 0/1)
    // 2. Cast Lightning Bolt (instant in graveyard, should be 1/2)
    // 3. Creature dies (creature in graveyard, should be 2/3)
    // 4. Play land (land in graveyard via discard, should be 3/4)
}
```

## Migration Path

### Phase 1: Foundation (Current)
- ✅ CDA system implemented
- ✅ GameContext interface extended
- ✅ CardInfo interface defined

### Phase 2: Engine Integration (1-2 days)
- Add GameID to internalCard
- Implement GameContext methods in MageEngine
- Create CardInfoAdapter
- Update getCreaturePower/getCreatureToughness

### Phase 3: Card Implementation (1 day)
- Add CDAs to Tarmogoyf, Lord of Extinction in generated cards
- Test with manual card implementations first
- Update card transpiler to recognize and generate CDAs

### Phase 4: Testing (1 day)
- Unit tests for each CDA type
- Integration tests for gameplay scenarios
- Regression tests to ensure static P/T still works

## Estimated Effort

- **Engine Integration**: 1-2 days
- **Card Updates**: 1 day
- **Testing**: 1 day
- **Total**: 3-4 days

## Success Criteria

✅ Tarmogoyf correctly counts card types in graveyards
✅ Lord of Extinction correctly counts all graveyard cards
✅ Maro correctly matches hand size
✅ Static P/T creatures still work
✅ Layer system still functions
✅ No performance degradation

## Notes

- CDAs function in all zones (Rule 604.3), but currently only needed on battlefield
- CDAs are calculated on-demand, not cached (may need optimization later)
- Context propagation is critical for game ID tracking
- CardInfo adapter avoids circular dependencies between abilities and game packages
