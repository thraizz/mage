# Apply Logic Integration Plan - Complete

## Summary

This document provides a comprehensive plan for implementing the runtime `Apply()` logic for search library effects and integrating them with the Go MAGE game engine.

## Current State

### Completed
✅ **SearchLibraryPutInHandEffect** - Structure created in `search_library.go:15-55`
✅ **SearchLibraryPutInPlayEffect** - Structure created in `search_library.go:63-120`
✅ **SearchLibraryPutOnTopEffect** - Structure created in `search_library.go:128-166`
✅ **Filter parsing** - StaticFilters mapping implemented in transpiler
✅ **Activated ability extraction** - SimpleActivatedAbility cost and effect extraction
✅ **Cost extraction** - TapSourceCost, SacrificeSourceCost, GenericManaCost
✅ **ActivatedAbilityBuilder** - Builder pattern with cost helpers

### TODO Placeholders

All three search effects have placeholder `Apply()` methods that return `nil`:

```go
func (e *SearchLibraryPutInHandEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
    // TODO: Implement actual search logic
    return nil
}
```

## Architecture Overview

### Game Context Interface

The `GameContext` interface (from `abilities/effects.go`) provides these methods for runtime operations:

```go
type GameContext interface {
    GetCard(id uuid.UUID) *game.Card
    GetPlayer(id uuid.UUID) *game.Player
    // TODO: Add more methods as needed:
    // - MoveCard(cardID, from, to Zone)
    // - ShuffleLibrary(playerID)
    // - RevealCard(cardID)
    // - GetController(cardID) *Player
}
```

## Integration Steps

### Phase 1: Core Zone Management (Week 1-2)

**Goal**: Implement basic zone movement and library operations

#### 1.1 Zone Enumeration
Add zone types to `internal/game/zones.go`:

```go
type Zone int

const (
    ZoneLibrary Zone = iota
    ZoneHand
    ZoneBattlefield
    ZoneGraveyard
    ZoneExile
    ZoneStack
    ZoneCommand
)
```

#### 1.2 Player Library Methods
Add to `internal/game/player.go`:

```go
type Player struct {
    ID      uuid.UUID
    Library []*Card
    Hand    []*Card
    // ... other zones
}

// SearchLibrary searches the library for cards matching a filter
func (p *Player) SearchLibrary(filter TargetFilter, maxResults int) []*Card {
    matching := []*Card{}
    for _, card := range p.Library {
        if filter.IsValid(card) {
            matching = append(matching, card)
            if len(matching) >= maxResults {
                break
            }
        }
    }
    return matching
}

// ShuffleLibrary randomizes the library order
func (p *Player) ShuffleLibrary() {
    rand.Shuffle(len(p.Library), func(i, j int) {
        p.Library[i], p.Library[j] = p.Library[j], p.Library[i]
    })
}

// MoveCardToHand moves a card from library to hand
func (p *Player) MoveCardToHand(card *Card) error {
    // Remove from library
    for i, c := range p.Library {
        if c.ID == card.ID {
            p.Library = append(p.Library[:i], p.Library[i+1:]...)
            break
        }
    }
    // Add to hand
    p.Hand = append(p.Hand, card)
    return nil
}

// MoveCardToBattlefield moves a card from library to battlefield
func (p *Player) MoveCardToBattlefield(card *Card, tapped bool) error {
    // Implementation depends on battlefield structure
    return nil
}

// PutOnTopOfLibrary puts cards on top of library
func (p *Player) PutOnTopOfLibrary(cards []*Card) {
    p.Library = append(cards, p.Library...)
}
```

#### 1.3 Game Context Implementation
Update `internal/game/game_context.go`:

```go
func (gc *gameContext) MoveCard(cardID uuid.UUID, from, to Zone) error {
    // Implementation
}

func (gc *gameContext) ShuffleLibrary(playerID uuid.UUID) error {
    player := gc.GetPlayer(playerID)
    if player == nil {
        return fmt.Errorf("player not found")
    }
    player.ShuffleLibrary()
    return nil
}

func (gc *gameContext) GetController(sourceID uuid.UUID) (*Player, error) {
    card := gc.GetCard(sourceID)
    if card == nil {
        return nil, fmt.Errorf("card not found")
    }
    return gc.GetPlayer(card.ControllerID), nil
}
```

### Phase 2: Target Filter Implementation (Week 2-3)

**Goal**: Implement filter validation logic

#### 2.1 TargetFilter Interface Extension
Update `internal/game/abilities/targets.go`:

```go
type TargetFilter interface {
    IsValid(card *game.Card) bool
    GetDescription() string
}

// Implement for each filter type
type LandTargetFilter struct{}

func (f *LandTargetFilter) IsValid(card *game.Card) bool {
    for _, t := range card.Types {
        if t == "LAND" {
            return true
        }
    }
    return false
}

type CreatureTargetFilter struct{}

func (f *CreatureTargetFilter) IsValid(card *game.Card) bool {
    for _, t := range card.Types {
        if t == "CREATURE" {
            return true
        }
    }
    return false
}

// ... implement for all filter types
```

#### 2.2 Composite Filters
For complex filters (e.g., "creature with power 2 or less"), create filter builders:

```go
type CompositeFilter struct {
    filters []TargetFilter
    mode    FilterMode  // AND or OR
}

type FilterMode int
const (
    FilterModeAnd FilterMode = iota
    FilterModeOr
)

func (f *CompositeFilter) IsValid(card *game.Card) bool {
    if f.mode == FilterModeAnd {
        for _, filter := range f.filters {
            if !filter.IsValid(card) {
                return false
            }
        }
        return true
    }
    // OR mode
    for _, filter := range f.filters {
        if filter.IsValid(card) {
            return true
        }
    }
    return false
}
```

### Phase 3: Search Effect Implementation (Week 3-4)

**Goal**: Complete Apply() methods for all search effects

#### 3.1 SearchLibraryPutInHandEffect.Apply()

```go
func (e *SearchLibraryPutInHandEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
    // 1. Get the controller
    controller, err := game.GetController(source)
    if err != nil {
        return fmt.Errorf("failed to get controller: %w", err)
    }

    // 2. Search library with target filter
    maxResults := e.target.Max
    if maxResults == 0 {
        maxResults = 1
    }

    matching := controller.SearchLibrary(e.target.Filter, maxResults)

    if len(matching) == 0 {
        // No cards found - still shuffle
        game.ShuffleLibrary(controller.ID)
        return nil
    }

    // 3. Let player choose which cards to take (for now, take all up to max)
    chosen := matching
    if len(chosen) > maxResults {
        chosen = chosen[:maxResults]
    }

    // 4. Reveal cards if required
    if e.reveal {
        for _, card := range chosen {
            // TODO: Add reveal event to game log
            _ = card
        }
    }

    // 5. Move found cards to hand
    for _, card := range chosen {
        if err := controller.MoveCardToHand(card); err != nil {
            return fmt.Errorf("failed to move card to hand: %w", err)
        }
    }

    // 6. Shuffle library
    game.ShuffleLibrary(controller.ID)

    return nil
}
```

#### 3.2 SearchLibraryPutInPlayEffect.Apply()

```go
func (e *SearchLibraryPutInPlayEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
    // 1. Get the controller
    controller, err := game.GetController(source)
    if err != nil {
        return fmt.Errorf("failed to get controller: %w", err)
    }

    // 2. If optional, ask whether to search
    if e.optional {
        // TODO: Implement player choice system
        // For now, always search
    }

    // 3. Search library with target filter
    maxResults := e.target.Max
    if maxResults == 0 {
        maxResults = 1
    }

    matching := controller.SearchLibrary(e.target.Filter, maxResults)

    if len(matching) == 0 {
        game.ShuffleLibrary(controller.ID)
        return nil
    }

    // 4. Let player choose which cards to take
    chosen := matching
    if len(chosen) > maxResults {
        chosen = chosen[:maxResults]
    }

    // 5. Move found cards to battlefield (tapped or untapped)
    for _, card := range chosen {
        if err := controller.MoveCardToBattlefield(card, e.tapped); err != nil {
            return fmt.Errorf("failed to move card to battlefield: %w", err)
        }
    }

    // 6. Shuffle library
    game.ShuffleLibrary(controller.ID)

    return nil
}
```

#### 3.3 SearchLibraryPutOnTopEffect.Apply()

```go
func (e *SearchLibraryPutOnTopEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
    // 1. Get the controller
    controller, err := game.GetController(source)
    if err != nil {
        return fmt.Errorf("failed to get controller: %w", err)
    }

    // 2. Search library with target filter
    maxResults := e.target.Max
    if maxResults == 0 {
        maxResults = 1
    }

    matching := controller.SearchLibrary(e.target.Filter, maxResults)

    if len(matching) == 0 {
        game.ShuffleLibrary(controller.ID)
        return nil
    }

    // 3. Let player choose which card to take
    chosen := matching
    if len(chosen) > maxResults {
        chosen = chosen[:maxResults]
    }

    // 4. Reveal cards if required
    if e.reveal {
        for _, card := range chosen {
            // TODO: Add reveal event to game log
            _ = card
        }
    }

    // 5. Shuffle library (before putting card on top!)
    game.ShuffleLibrary(controller.ID)

    // 6. Put found card on top of library
    controller.PutOnTopOfLibrary(chosen)

    return nil
}
```

### Phase 4: Cost Payment Implementation (Week 4-5)

**Goal**: Implement cost payment for activated abilities

#### 4.1 TapCost.Pay()

```go
func (c *TapCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID, sourceID uuid.UUID) error {
    card := game.GetCard(sourceID)
    if card == nil {
        return fmt.Errorf("source card not found")
    }

    if card.Tapped {
        return fmt.Errorf("card is already tapped")
    }

    card.Tapped = true
    // TODO: Trigger tapped events
    return nil
}

func (c *TapCost) CanPay(ctx context.Context, game GameContext, playerID uuid.UUID, sourceID uuid.UUID) bool {
    card := game.GetCard(sourceID)
    if card == nil {
        return false
    }
    return !card.Tapped
}
```

#### 4.2 SacrificeSourceCost.Pay()

```go
func (c *SacrificeCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID, sourceID uuid.UUID) error {
    if c.Filter == "source" {
        // Sacrifice the source permanent
        card := game.GetCard(sourceID)
        if card == nil {
            return fmt.Errorf("source card not found")
        }

        // Move to graveyard
        return game.MoveCard(sourceID, ZoneBattlefield, ZoneGraveyard)
    }

    // Other sacrifice costs (sacrifice a creature, etc.)
    // TODO: Implement target selection for sacrifice
    return fmt.Errorf("non-source sacrifice not yet implemented")
}
```

#### 4.3 ManaCost.Pay()

```go
func (c *ManaCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error {
    player := game.GetPlayer(playerID)
    if player == nil {
        return fmt.Errorf("player not found")
    }

    // Check if player has enough mana
    if !c.CanPay(ctx, game, playerID) {
        return fmt.Errorf("not enough mana")
    }

    // Deduct mana from player's mana pool
    // TODO: Implement mana pool system
    player.ManaPool.Subtract(c.Mana)

    return nil
}
```

### Phase 5: Activated Ability Activation (Week 5-6)

**Goal**: Implement full activated ability activation flow

#### 5.1 ActivatedAbility.Activate()

```go
func (a *ActivatedAbility) Activate(ctx context.Context, game GameContext, playerID uuid.UUID) error {
    // 1. Check if ability can be activated
    if !a.CanActivate(ctx, game) {
        return fmt.Errorf("ability cannot be activated")
    }

    // 2. Check timing restrictions
    // TODO: Implement timing rules (sorcery speed, etc.)

    // 3. Pay costs
    for _, cost := range a.Costs {
        if err := cost.Pay(ctx, game, playerID); err != nil {
            return fmt.Errorf("failed to pay cost: %w", err)
        }
    }

    // 4. Put ability on stack
    // TODO: Implement stack system

    // 5. Choose targets if needed
    if a.Targets != nil {
        // TODO: Implement target selection
    }

    return nil
}
```

#### 5.2 Stack Resolution

```go
type StackItem struct {
    ID         uuid.UUID
    Ability    Ability
    Controller uuid.UUID
    Targets    []uuid.UUID
}

func (g *Game) ResolveStack() error {
    if len(g.Stack) == 0 {
        return nil
    }

    // Pop top item from stack
    item := g.Stack[len(g.Stack)-1]
    g.Stack = g.Stack[:len(g.Stack)-1]

    // Resolve the ability
    return item.Ability.Resolve(g.Context)
}
```

### Phase 6: Player Interaction System (Week 6-7)

**Goal**: Implement player choices for search effects

#### 6.1 Choice Interface

```go
type Choice interface {
    GetOptions() []string
    IsValid(response string) bool
}

type CardChoice struct {
    Cards     []*Card
    MinChoose int
    MaxChoose int
}

func (c *CardChoice) GetOptions() []string {
    options := make([]string, len(c.Cards))
    for i, card := range c.Cards {
        options[i] = card.Name
    }
    return options
}
```

#### 6.2 Player.MakeChoice()

```go
func (p *Player) MakeChoice(ctx context.Context, choice Choice) ([]int, error) {
    // This will integrate with the UI/network layer
    // For now, return default choice
    // TODO: Implement choice system with UI integration
    return []int{0}, nil
}
```

#### 6.3 Integration with Search Effects

Update SearchLibraryPutInHandEffect.Apply() to use choices:

```go
// 3. Let player choose which cards to take
choice := &CardChoice{
    Cards:     matching,
    MinChoose: e.target.Min,
    MaxChoose: e.target.Max,
}

chosenIndices, err := controller.MakeChoice(ctx, choice)
if err != nil {
    return fmt.Errorf("failed to make choice: %w", err)
}

chosen := make([]*Card, len(chosenIndices))
for i, idx := range chosenIndices {
    chosen[i] = matching[idx]
}
```

## Testing Strategy

### Unit Tests

```go
// internal/game/abilities/search_library_test.go

func TestSearchLibraryPutInHandEffect(t *testing.T) {
    // Setup
    game := NewMockGameContext()
    player := NewTestPlayer()

    // Add cards to library
    player.Library = []*Card{
        {ID: uuid.New(), Types: []string{"LAND"}, Subtypes: []string{"PLAINS"}},
        {ID: uuid.New(), Types: []string{"CREATURE"}},
        {ID: uuid.New(), Types: []string{"LAND"}, Subtypes: []string{"ISLAND"}},
    }

    // Create effect
    effect := NewSearchLibraryPutInHandEffect(
        NewTargetRequirement(0, 1, NewLandTargetFilter()),
        true, // reveal
    )

    // Execute
    err := effect.Apply(context.Background(), game, player.ID, nil)

    // Assert
    require.NoError(t, err)
    assert.Equal(t, 1, len(player.Hand))
    assert.Equal(t, 2, len(player.Library))
    assert.True(t, isLand(player.Hand[0]))
}
```

### Integration Tests

```go
// internal/integration/search_effects_test.go

func TestRenegadeMapActivation(t *testing.T) {
    // Create game
    game := NewGame()
    player := game.AddPlayer("Player1")

    // Create RenegadeMap
    card, err := generated.NewRenegadeMap(player.ID, nil)
    require.NoError(t, err)

    // Put on battlefield
    game.AddToBattlefield(card)

    // Add basic land to library
    player.Library = []*Card{
        NewBasicLand("Plains"),
    }

    // Activate ability
    ability := card.Abilities[0].(*ActivatedAbility)
    err = ability.Activate(context.Background(), game, player.ID)

    // Assert
    require.NoError(t, err)
    assert.Equal(t, 1, len(player.Hand))
    assert.Equal(t, "Plains", player.Hand[0].Name)
}
```

## Timeline

| Week | Phase | Deliverables |
|------|-------|--------------|
| 1-2  | Zone Management | Zone enum, Player library methods, basic GameContext |
| 2-3  | Target Filters | TargetFilter implementation, composite filters |
| 3-4  | Search Effects | Complete Apply() for all 3 search effects |
| 4-5  | Cost Payment | TapCost, SacrificeCost, ManaCost payment |
| 5-6  | Ability Activation | Activate() flow, stack system |
| 6-7  | Player Choices | Choice system, UI integration |
| 7-8  | Testing | Unit tests, integration tests, bug fixes |

## Dependencies

### External
- **UI/Network Layer**: For player choice system (Phase 6)
- **Event System**: For card reveal, zone change events
- **Stack System**: For ability resolution order

### Internal
- **Card Structure**: Need to add `Tapped` field, `ControllerID`
- **Player Structure**: Need `ManaPool` field
- **Game State**: Need stack, priority system

## Success Criteria

✅ RenegadeMap can be activated and search for basic lands
✅ RamosianSergeant can search for Rebel permanents
✅ All costs can be paid correctly (Tap, Sacrifice, Mana)
✅ Library is shuffled after search
✅ Cards are revealed when required
✅ Player can choose from multiple matching cards
✅ All search effects have >80% test coverage

## Next Steps

1. **Immediate**: Implement Phase 1 (Zone Management)
2. **Week 2**: Start Phase 2 (Target Filters)
3. **Week 3**: Begin search effect Apply() methods
4. **Document**: Update SEARCH_LIBRARY_IMPLEMENTATION.md with progress

## References

- `internal/game/abilities/search_library.go` - Effect structures
- `internal/game/abilities/activated.go` - Activated ability structure
- `internal/game/abilities/costs.go` - Cost implementations
- `scripts/transpile_cards.py` - Card transpiler with filter parsing
