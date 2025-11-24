# CLAUDE.md

## Project Overview

**Go port of XMage (Magic: The Gathering simulator)** with hybrid gRPC + WebSocket architecture:
- **gRPC**: 60+ request/response RPC methods (auth, room, table, game, tournament)
- **WebSocket**: Server-to-client push events (game updates, chat, real-time notifications)
- **Protocol Buffers**: Type-safe serialization for both protocols

**Goals**: 100% MTG rule coverage, full test coverage, all 30,400+ cards supported.

## Essential Commands

```bash
# Build and run
make proto              # MUST run before first build - generates protobuf code
make build              # Build server binary
make run                # Build and run with config.yaml

# Testing
make test               # All unit tests (82 tests)
make test-integration   # Integration tests only
go test -v -run TestName ./internal/package

# Code quality
make fmt                # gofmt + goimports
make lint               # golangci-lint

# Database
make migrate-up         # Apply migrations
make migrate-down       # Rollback one migration
make migrate-create NAME=feature_name
```

## Architecture

### Component Initialization Order

**CRITICAL** - Initialize in this exact order:

**Server Components**:
1. Config → Logger → Database
2. Session Manager (starts cleanup goroutine)
3. Repositories (User, Stats, Card)
4. Domain Managers (User, Room, Chat, Table, Game, Tournament, Draft)
5. External Services (Email)
6. Servers (gRPC + WebSocket)

**Game Engine Components**:
1. Core: `TurnManager`, `LayerSystem`, `AbilityRegistry`
2. Managers: `PriorityManager`, `EnhancedStackManager`, `TargetSelectionManager`, `ContinuousEffectsManager`
3. Integration: `PriorityManager.SetLayerManager(continuousEffectsMgr)`
4. High-Level: `AbilityActivationManager` (needs all of the above)

### Manager Interface Pattern

All components use this pattern:

```go
type Manager interface {
    Create(...) (*Entity, error)
    Get(id string) (*Entity, bool)
    Update(...) error
    Remove(id string)
}

type manager struct {
    data   map[string]*Entity
    mu     sync.RWMutex
    logger *zap.Logger
}
```

**Key Managers**:
- `session.Manager`: Lease-based session expiration
- `user.Manager`: Authentication, registration
- `room.Manager`: Lobby system
- `table.Manager`: Table state machine
- `game.Manager`: Game controller (integrates with engine)
- `tournament.Manager`: Swiss pairing, elimination brackets
- `draft.Manager`: Booster draft mechanics
- `chat.Manager`: Chat sessions

### Session Management

**Lease-based expiration** (not simple timeouts):
- `LeaseUntil` timestamp per session
- `UpdateActivity()` extends lease by 5 minutes
- Background goroutine cleans expired sessions
- Ping RPC keeps sessions alive

### gRPC Interceptor Chain

```
Request → RecoveryInterceptor → LoggingInterceptor → SessionValidationInterceptor
       → AdminInterceptor → MetricsInterceptor → Handler
```

### Authentication

- **Password Hashing**: Argon2id (time=1, memory=64MB, threads=4, keyLen=32)
- **Format**: `$argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>`
- **Validation**: Username 3-16 chars, Password 8+ chars with complexity

### Database

- **Driver**: pgx (not lib/pq)
- **Connection pooling**: max_conns=25, min_conns=5
- **Repositories**: UserRepository, StatsRepository, CardRepository

### Rating System

- **Algorithm**: Glicko-2
- **Defaults**: rating=1500, deviation=350, volatility=0.06
- **Location**: `internal/rating/glicko.go`

### Tournament & Draft

- **Swiss Pairing**: Score-based pairing with bye handling
- **Draft**: Pick tracking, booster passing, configurable sets

## Game Engine Integration

### Card Implementation System

**30,400+ generated cards** transpiled from Java XMage:
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go
python scripts/transpile_cards.py
```

**Structure**:
- Location: `internal/game/cards/generated/`
- Each card: Builder function returning `*game.Card`
- Registration: `init()` functions → global registry
- Complex cards: `internal/game/cards/manual/`

**Example**:
```go
func NewLightningBolt(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
    card := game.NewCard(ownerID, "Lightning Bolt")
    card.ManaCost = "{R}"
    card.Types = []string{"INSTANT"}

    ability := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
        AddTarget(abilities.NewAnyTarget()).
        AddEffect(abilities.NewDamageEffect(3)).
        Build()
    card.AddAbility(ability)
    return card, nil
}
```

### Abilities System

**6 Ability Types** (`internal/game/abilities/`):
1. **Spell Abilities**: Instants and sorceries
2. **Activated Abilities**: Cost: Effect (e.g., "{T}: Add {G}")
3. **Triggered Abilities**: When/Whenever/At triggers
4. **Static Abilities**: Continuous effects with layers (Rule 613)
5. **Mana Abilities**: Special timing, don't use stack
6. **Keyword Abilities**: Flying, trample, etc.

**Builder Pattern**:
```go
// Activated ability
ability := abilities.NewActivatedAbilityBuilder(sourceID).
    AddCost(abilities.NewManaCost("{2}{U}")).
    AddCost(abilities.NewTapCost()).
    AddTarget(abilities.NewCreatureTarget()).
    AddEffect(abilities.NewTapEffect()).
    Build()

// Triggered ability
ability := abilities.NewTriggeredAbilityBuilder(sourceID).
    SetTrigger(abilities.NewEntersBattlefieldTrigger(sourceID)).
    AddEffect(abilities.NewDrawCardsEffect(1)).
    Build()

// Static ability
ability := abilities.NewSimpleStaticAbility(sourceID, abilities.ZoneBattlefield).
    AddEffect(abilities.NewBoostEffect(1, 1)).
    Build()
```

**40+ Effects** organized by category:
- **Damage**: DamageEffect, LifeGainEffect, LoseLifeEffect
- **Card Draw**: DrawCardsEffect, DiscardEffect
- **Permanents**: DestroyEffect, TapEffect, UntapEffect
- **P/T Modification**: BoostEffect (layer 7)
- **Mana**: AddManaEffect
- **Counters**: AddCounterEffect, RemoveCounterEffect
- **Tokens**: CreateTokenEffect
- **Attachment**: AttachEffect, GainAbilityAttachedEffect
- **Special**: CounterSpellEffect, SearchLibraryEffect

**Cost System**:
- `ManaCost`: Parses "{2}{U}{U}", integrates with mana pool
- `TapCost`: Tap source permanent
- `SacrificeCost`: Sacrifice permanents with filters
- `DiscardTargetCost`: Discard cards matching filter
- `PayLifeCost`: Pay life points

**Target System**:
- `TargetFilter`: Battlefield permanents (creature, player, any, etc.)
- `CardFilter`: Hand/graveyard cards (artifact card, creature card, etc.)
- `TargetRequirement`: Min/max targets, targeting mode

**Trigger System**:
- `EntersBattlefieldTrigger`: When permanent enters
- `LeavesBattlefieldTrigger`: When permanent leaves
- `DiesTrigger`: When creature/planeswalker dies
- `BecomesTappedTrigger`: When permanent becomes tapped
- `DealsDamageTrigger`: When source deals damage
- `DrawCardTrigger`: When player draws
- `SpellCastTrigger`: When spell is cast

**Layer System** (Rule 613):
```
Layer 1: Copy effects
Layer 2: Control-changing effects
Layer 3: Text-changing effects
Layer 4: Type-changing effects
Layer 5: Color-changing effects
Layer 6: Ability-adding/removing effects
Layer 7: Power/toughness effects (7a: CDA, 7b: Set, 7c: Modify, 7d: Counters, 7e: Switch)
```

**Zone System**:
- `ZoneLibrary`, `ZoneHand`, `ZoneBattlefield`, `ZoneGraveyard`
- `ZoneStack`, `ZoneExile`, `ZoneCommand`, `ZoneOutside`

**Duration System**:
- `DurationUntilEndOfTurn`, `DurationPermanent`
- `DurationWhileOnBattlefield`, `DurationUntilEndOfCombat`
- `DurationWhileInGraveyard`, `DurationWhileInHand`, `DurationWhileInExile`

### Engine Integration Components

**Stack Management** (`internal/game/engine_stack.go`):
- `EnhancedStackManager`: Wraps rules stack with abilities integration
- `StackObject`: Stores spell/ability with full context (targets, choices)
- `PushSpell()`, `PushActivatedAbility()`, `PushTriggeredAbility()`
- `ResolveTop()`: Resolves with automatic SBA checks
- `Counter()`: Removes object from stack (counterspells)
- Implements Rule 405 (Stack), Rule 608 (Resolution)

**Ability Registry** (`internal/game/ability_registry.go`):
- `AbilityRegistry`: Thread-safe UUID → ability mapping
- `RegisterAbility()`: Register when card enters play
- `GetAbility()`: Retrieve by ID for activation
- `GetAbilitiesBySource()`: Get all abilities of a card
- `UnregisterSource()`: Cleanup when card leaves
- `AbilityMetadata`: Tracks controller, zone, index

**Target Selection** (`internal/game/engine_targeting.go`):
- `TargetSelectionManager`: Handles targeting workflow
- `ValidateTargets()`: Checks legality, count, duplicates
- `GetLegalTargets()`: Calculates from filters
- `CreateTargetRequest()`: Builds from ability
- `AutoSelectTargets()`: Optimizes single legal choice
- `CheckTargetingRestrictions()`: Hexproof/shroud/protection
- `ValidateDivision()`: For "distribute X" effects
- Implements Rule 115 (Targets)

**Layer Recalculation** (`internal/game/engine_layers.go`):
- `ContinuousEffectsManager`: Implements Rule 613
- `RecalculateAll()`: Processes all static abilities
- `processPermanent()`: Extracts static abilities
- `convertToContinuousEffect()`: Converts to layer effects
- `ApplyToCard()`: Applies all layers to card
- `RemoveSourceEffects()`: Cleanup on zone change
- Integrates abilities with effects layer system
- Called before every SBA check

**Ability Activation** (`internal/game/engine_abilities.go`):
- `AbilityActivationManager`: Manages activation workflow
- `ActivateAbility()`: Full Rule 602 sequence
- `CastSpell()`: Full Rule 601 sequence
- `ActivateManaAbility()`: Rule 605 (no stack)
- `payCosts()`: Ordered cost payment (mana first)
- `CheckActivationRestrictions()`: Timing validation
- Integrates with stack, registry, targeting

**Priority Integration** (`internal/game/engine_priority.go`):
- `PriorityManager`: Integrates SBAs + layers + turn structure
- `SetLayerManager()`: Connects layer system
- `CheckStateBasedActions()`: Now recalculates layers first
- `GivePriority()`: Checks SBAs before giving priority
- Sequence: Layers → SBAs → Execute → Repeat if needed

**Combat Integration** (`internal/game/engine_combat.go`):
- `CombatIntegrationManager`: Connects combat to triggered abilities
- `OnDeclareAttackers()`: Triggers "When X attacks"
- `OnDeclareBlockers()`: Triggers "When X blocks" and "becomes blocked"
- `OnCombatDamage()`: Triggers combat damage abilities
- `CheckCombatKeywordAbilities()`: Extracts flying, first strike, etc.
- `CanAttack()` / `CanBlock()`: Ability-based restrictions
- Added triggers: AttacksTrigger, BlocksTrigger, BecomesBlockedTrigger, DealsCombatDamageTrigger
- Implements Rule 508-510 (Combat steps)

### Game Engine Architecture

**Turn Structure** (`internal/game/rules/turn.go`):
- `TurnManager`: Tracks phase/step progression
- Full turn sequence: Untap → Upkeep → Draw → Main1 → Combat (BeginCombat, DeclareAttackers, DeclareBlockers, [FirstStrikeDamage], CombatDamage, EndCombat) → Main2 → End → Cleanup
- Dynamic first strike step insertion
- Active and priority player tracking

**Event System** (`internal/game/rules/events.go`):
- 200+ `EventType` constants
- Phase/step events, zone changes, damage/life, card play, combat
- Event adapter bridges to abilities system (`internal/game/engine_events.go`)

**State-Based Actions** (`internal/game/rules/state_based_actions.go`):
Implements Rule 704 - automatic game actions:
- Player life ≤ 0 (704.5a)
- Poison counters ≥ 10 (704.5b)
- Creature toughness ≤ 0 (704.5f)
- Lethal damage (704.5g)
- Deathtouch damage (704.5h)
- Planeswalker loyalty = 0 (704.5i)
- Legend rule (704.5j)
- Aura/Equipment attachment validity (704.5k, 704.5m)
- Counter annihilation: +1/+1 and -1/-1 (704.5q)

**Priority System** (`internal/game/engine_priority.go`):
- `PriorityManager`: Integrates SBAs with turn structure
- Checks SBAs when giving priority (Rule 117.5)
- Checks SBAs after spell/ability resolves (Rule 608.2k)
- Handles priority passing and step advancement

**Mana System** (`internal/game/mana/pool.go`):
- `ManaPool`: Thread-safe with regular + floating mana
- Automatic emptying at step/phase boundaries
- Integrated with `ManaCost.CanPay()` and `ManaCost.Pay()`

**Rules Engine** (`internal/game/rules/`):
- `Stack`: Spell/ability stack management
- `TriggerManager`: Trigger collection and resolution
- `PriorityManager`: Priority passing rounds
- `ResolutionContext`: Nested resolution tracking

**Effects System** (`internal/game/effects/`):
- `Layer`: Continuous effects layer system (Rule 613)
- `ReplacementEffect`: Replacement effects
- `StaticEffects`: Static ability tracking

### Integration Status

**✅ Completed**:
- Event adapter (rules → abilities events)
- State-based actions checker (Rule 704)
- Priority manager with SBA integration
- Mana cost payment (CanPay/Pay with mana pool)
- 40+ effect implementations
- Complete trigger system including combat triggers
- Static abilities with layers
- Cost system with filters
- **Combat integration** (Rule 508-510, triggered abilities)
- **Stack management system** (Rule 405, 608)
- **Ability retrieval registry** (UUID-based mapping)
- **Target selection and validation** (Rule 115)
- **Continuous effects layer recalculation** (Rule 613)
- **Spell casting workflow** (Rule 601)

**📋 Remaining**:
- Comprehensive integration tests
- Push triggered combat abilities to stack
- Modal spell support
- Cost modifications (increases/reductions)

### Key Files

**Abilities**:
- `internal/game/abilities/ability.go` - Core interfaces
- `internal/game/abilities/activated.go` - Activated abilities
- `internal/game/abilities/triggered.go` - Triggered abilities
- `internal/game/abilities/static.go` - Static abilities, layers, zones
- `internal/game/abilities/effects.go` - 40+ effect implementations
- `internal/game/abilities/costs.go` - Cost system with mana integration
- `internal/game/abilities/targets.go` - Target and filter system
- `internal/game/abilities/enchanted_effects.go` - Aura/Equipment effects

**Engine**:
- `internal/game/mage_engine.go` - Main game engine
- `internal/game/game_context.go` - GameContext implementation
- `internal/game/engine_events.go` - Event adapter
- `internal/game/engine_priority.go` - Priority + SBA + layer integration
- `internal/game/engine_combat.go` - Combat integration with triggered abilities
- `internal/game/engine_stack.go` - Enhanced stack with abilities integration
- `internal/game/engine_targeting.go` - Target selection and validation
- `internal/game/engine_layers.go` - Continuous effects layer recalculation
- `internal/game/engine_abilities.go` - Ability activation and spell casting
- `internal/game/ability_registry.go` - Ability ID to object mapping

**Rules**:
- `internal/game/rules/turn.go` - Turn structure
- `internal/game/rules/events.go` - 200+ event types
- `internal/game/rules/state_based_actions.go` - SBA checker (Rule 704)
- `internal/game/rules/priority.go` - Priority windows
- `internal/game/rules/stack.go` - Stack management

**Mana**:
- `internal/game/mana/pool.go` - Mana pool with GetAmount/SpendMana

**Cards**:
- `internal/game/cards/generated/` - 30,400+ transpiled cards
- `internal/game/cards/manual/` - Complex manual implementations
- Example: `internal/game/cards/generated/psychicoverload.go`

### Complex Card Example (Psychic Overload)

```go
// Enchant permanent
ability0 := abilities.NewEnchantAbility(card.ID,
    abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter()))

// Spell ability: Attach
ability1 := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
    AddTarget(abilities.NewPermanentTarget()).
    AddEffect(abilities.NewAttachEffect(abilities.OutcomeDetriment)).
    Build()

// Triggered: ETB tap enchanted permanent
ability2 := abilities.NewTriggeredAbilityBuilder(card.ID).
    SetTrigger(abilities.NewEntersBattlefieldTrigger(card.ID)).
    AddEffect(abilities.NewTapEnchantedEffect()).
    Build()

// Static: Doesn't untap
ability3 := abilities.NewSimpleStaticAbility(card.ID, abilities.ZoneBattlefield).
    AddEffect(abilities.NewDontUntapInControllersUntapStepEnchantedEffect()).
    Build()

// Static: Grant activated ability to enchanted permanent
grantedAbility := abilities.NewActivatedAbilityBuilder(card.ID).
    AddCost(abilities.NewDiscardTargetCost(2, abilities.NewArtifactCardFilter())).
    AddEffect(abilities.NewUntapSourceEffect()).
    Build()

ability4 := abilities.NewSimpleStaticAbility(card.ID, abilities.ZoneBattlefield).
    AddEffect(abilities.NewGainAbilityAttachedEffect(
        grantedAbility,
        abilities.AttachmentTypeAura,
        abilities.DurationWhileOnBattlefield,
        "Enchanted permanent has \"Discard two artifact cards: Untap this permanent.\"",
    )).
    Build()
```

## Protocol Buffers

**CRITICAL**: Run `make proto` before first build to generate Go code from `.proto` files.

**Proto files**: `api/proto/mage/v1/*.proto`
**Generated code**: `pkg/proto/mage/v1/`

**Import in Go**:
```go
import pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
```

## Code Conventions

1. **Interfaces before implementations**: Define Manager interfaces, implement as private structs
2. **Constructor pattern**: `NewManager()` returns interface
3. **Context propagation**: All I/O accepts `context.Context`
4. **Structured logging**: `zap.Logger` with typed fields
5. **Error wrapping**: `fmt.Errorf("context: %w", err)`
6. **Mutex discipline**: `sync.RWMutex` for read-heavy structures

## Testing

**82 tests total**:
- Authentication (Argon2id) - 4 tests
- Session management - 6 tests
- Rating (Glicko-2) - 9 tests
- Draft mechanics - 8 tests
- Tournament (Swiss) - 9 tests
- Abilities system - 21 tests
- Integration flows - 25 tests

**Test pattern**:
```go
func TestFeature(t *testing.T) {
    logger := zap.NewNop()
    mgr := NewManager(logger)

    result, err := mgr.DoSomething()

    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

## Configuration

**Key settings** (`config/config.yaml`):
```yaml
server:
  grpc:
    address: "0.0.0.0:17171"
  websocket:
    address: "0.0.0.0:17179"
  lease_period: 5m

database:
  host: "localhost"
  port: 5432
  max_conns: 25

auth:
  mode: "optional"  # or "required"

mail:
  provider: "smtp"  # or "mailgun" or "none"
```

**Environment overrides**: `DB_PASSWORD`, `SMTP_HOST`, `SMTP_USER`, `SMTP_PASSWORD`, `MAILGUN_DOMAIN`, `MAILGUN_API_KEY`

## Dependencies

- `google.golang.org/grpc` - gRPC framework
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/gorilla/websocket` - WebSocket server
- `go.uber.org/zap` - Structured logging
- `github.com/google/uuid` - UUID generation
- `golang.org/x/crypto` - Argon2id hashing

## References

- **MTG Rules**: `/Users/aron/dev/opensource/mage/RULES.txt`
- **Engine Gap Analysis**: `ENGINE_GAP_ANALYSIS.md` (missing systems, priority roadmap)
- **Integration Summary**: `INTEGRATION_WORK_SUMMARY.md` (completed work, 6 major systems)
- **Integration Status**: `ABILITIES_INTEGRATION_STATUS.md` (detailed status, test coverage)
- **Java XMage**: https://github.com/magefree/mage
