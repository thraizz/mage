# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **port of XMage (Magic: The Gathering game simulator) from Java to Go and Svelte**. The goal is to achieve everything the original Java XMage offers and more, with 100% rule coverage and full test coverage for all 28,000+ unique cards.

**Key Principles:**
- Build the MVP first before optimizing structure and architecture
- 100% MTG rule coverage is required
- Full test coverage is required
- All cards and spells must be supported
- Always use absolute paths for cd commands
- Never create new markdown files unless explicitly requested

## Architecture Overview

### Hybrid gRPC + WebSocket Design

The server uses a **hybrid protocol approach**:
- **gRPC**: All 60+ request/response RPC methods (authentication, room/table/game/tournament operations)
- **WebSocket**: Server-to-client push events (game updates, chat messages, real-time notifications)
- **Protocol Buffers**: Type-safe serialization for both protocols

This design provides efficient RPC with type safety while maintaining real-time push capabilities.

### Project Structure

```
mage-server-go/          # Go backend server
├── cmd/server/          # Server entry point
├── api/proto/           # Protocol Buffer definitions (11 .proto files)
├── internal/            # Core business logic
│   ├── server/          # gRPC + WebSocket servers, interceptors
│   ├── session/         # Lease-based session management
│   ├── user/            # User authentication and management
│   ├── auth/            # Argon2id password hashing
│   ├── table/           # Table state machine
│   ├── game/            # Game controller and engine
│   ├── tournament/      # Swiss/elimination pairing algorithms
│   ├── draft/           # Booster draft mechanics
│   ├── room/            # Lobby system
│   ├── chat/            # Chat sessions
│   ├── repository/      # PostgreSQL database layer (pgx)
│   ├── rating/          # Glicko-2 rating system
│   └── plugin/          # Game type registry
├── migrations/          # Database migrations
└── config/              # YAML configuration

mage-client-web/         # Svelte web client
├── src/
│   ├── lib/
│   │   ├── components/  # Reusable Svelte components
│   │   ├── stores/      # Reactive state management (auth, etc.)
│   │   ├── grpc/        # gRPC client setup
│   │   └── types/       # TypeScript type definitions
│   └── routes/          # SvelteKit file-based routing
│       ├── (protected)/ # Auth-required routes (lobby, game, etc.)
│       ├── login/       # Login page
│       └── register/    # Registration page
```

## Development Commands

### Go Server

**CRITICAL FIRST STEP** - Generate protobuf code before building:
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go
make proto              # MUST run before first build
```

**Common commands:**
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go

make build              # Build binary to bin/mage-server
make run                # Build and run with config.yaml
make test               # Run all 82 tests with coverage
make test-integration   # Run integration tests only
make lint               # Run golangci-lint
make fmt                # Run gofmt + goimports
make deps               # Download Go dependencies
make tools              # Install development tools (protoc-gen-go, etc.)
```

**Database:**
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go

# Start PostgreSQL via Docker
docker-compose up -d postgres

# Run migrations
make migrate-up         # Apply all migrations
make migrate-down       # Rollback last migration
```

**Running server manually:**
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go
go run cmd/server/main.go -config config/config.yaml
```

### Svelte Web Client

```bash
cd /Users/aron/dev/opensource/mage/mage-client-web

npm install             # Install dependencies
npm run dev             # Start dev server on :5173 (hot reload)
npm run build           # Production build
npm run preview         # Preview production build on :4173
npm run check           # TypeScript + Svelte type checking
npm run lint            # ESLint
npm run format          # Prettier formatting
```

### Testing

**Go tests:**
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go

# All tests
go test ./...

# With coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Single package
go test -v ./internal/rating

# Integration tests only
go test -tags=integration ./internal/integration/...

# Specific test
go test -v -run TestSessionLeaseExpiration ./internal/session
```

**Svelte tests:**
```bash
cd /Users/aron/dev/opensource/mage/mage-client-web
npm run test            # Run Vitest tests
```

### Docker Compose

```bash
cd /Users/aron/dev/opensource/mage

# Start full stack
docker-compose up

# Start specific services
docker-compose up postgres
docker-compose up mage-server

# With monitoring (Prometheus + Grafana)
docker-compose --profile monitoring up

# Rebuild after code changes
docker-compose up --build
```

## Card Implementation System

### Auto-Generated Cards (30,400+ files!)

Cards are **transpiled from Java XMage source** to Go using `scripts/transpile_cards.py`:

```bash
cd /Users/aron/dev/opensource/mage/mage-server-go
python scripts/transpile_cards.py
```

**Generated card structure:**
- Location: `internal/game/cards/generated/`
- Each card is a builder function returning `*game.Card`
- Registered in global registry via `init()` functions
- Abilities built using builder pattern

**Example generated card:**
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

### Manual Card Implementations

Complex cards requiring custom logic go in `internal/game/cards/manual/`:
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go/internal/game/cards/manual
```

## Abilities System

The abilities system (`internal/game/abilities/`) implements MTG's comprehensive ability framework following the Comprehensive Rules.

### Ability Types

**Six ability types** (Rule 113):
- **Spell Abilities**: Instants and sorceries on the stack
- **Activated Abilities**: Cost: Effect format (`{T}: Draw a card`)
- **Triggered Abilities**: When/Whenever/At patterns (`When ~ enters, ...`)
- **Static Abilities**: Continuous effects always active (`~ has flying`)
- **Mana Abilities**: Special activated abilities that produce mana
- **Keyword Abilities**: Flying, Trample, Haste, etc.

### Builder Pattern

All abilities use fluent builders for construction:

```go
// Spell ability (instant/sorcery)
spell := abilities.NewSpellAbilityBuilder(cardID, "{2}{U}").
    AddTarget(abilities.NewCreatureTargetFilter()).
    AddEffect(abilities.NewDamageEffect(3)).
    AddEffect(abilities.NewDrawCardsEffect(1)).
    Build()

// Activated ability
activated := abilities.NewActivatedAbilityBuilder(cardID).
    AddTapCost().
    AddManaCost("{1}").
    AddEffect(abilities.NewDrawCardsEffect(1)).
    Build()

// Triggered ability
triggered := abilities.NewTriggeredAbilityBuilder(cardID).
    SetTrigger(abilities.NewEntersBattlefieldTrigger(cardID)).
    AddEffect(abilities.NewDamageEffect(2)).
    AddTarget(abilities.NewAnyTargetFilter()).
    SetOptional(false).
    Build()

// Static ability (continuous effects)
static := abilities.NewSimpleStaticAbility(cardID, abilities.ZoneBattlefield).
    AddEffect(abilities.NewBoostEnchantedEffect(2, 2)).
    Build()
```

### Effects System

**40+ effect types** organized by category:

**Basic Effects**:
- `DamageEffect`: Deal damage to targets
- `DrawCardsEffect`: Player draws N cards
- `DestroyEffect`: Destroy target permanent
- `GainLifeEffect` / `LoseLifeEffect`: Modify life totals
- `TapEffect` / `UntapEffect`: Tap/untap permanents

**Card Movement**:
- `ReturnToHandTargetEffect`: Bounce to hand
- `ExileTargetEffect` / `ExileSourceEffect`: Exile permanents
- `MillCardsTargetEffect`: Mill from library
- `SearchLibraryPutInHandEffect`: Tutor effects

**Token & Counters**:
- `CreateTokenEffect`: Create creature tokens
- `AddCountersSourceEffect` / `AddCountersTargetEffect`: Add counters
- 100+ counter types (P1P1, M1M1, Loyalty, Poison, Energy, etc.)

**Ability Modification**:
- `GrantAbilityEffect`: Grant abilities to targets
- `GainAbilityAttachedEffect`: Enchanted/equipped creature gains abilities

**Continuous Effects**:
- `BoostEffect`: Modify power/toughness until end of turn
- `BoostEnchantedEffect` / `BoostEquippedEffect`: Static P/T modification
- `DontUntapInControllersUntapStepEnchantedEffect`: Prevention effects

### Cost System

**Comprehensive cost types**:
- `ManaCost`: Mana payment (`{2}{U}{U}`)
- `TapCost`: Tap symbol (`{T}`)
- `SacrificeSourceCost`: Sacrifice this permanent
- `SacrificeCost`: Sacrifice other permanents
- `DiscardCost` / `DiscardTargetCost`: Discard cards (with optional filtering)
- `PayLifeCost`: Pay life
- `CompositeCost`: Multiple costs combined

### Trigger System

**Common triggers**:
- `EntersBattlefieldTrigger`: When this enters (ETB)
- `LeavesBattlefieldTrigger`: When this leaves (LTB)
- `DiesTrigger`: When creature dies
- `BecomesTappedTrigger` / `BecomesUntappedTrigger`: Tap state changes
- `DealsDamageTrigger`: When this deals damage

### Target & Filter System

**Target Filters** (permanents on battlefield):
- `AnyTargetFilter`: Any legal target (creature, player, planeswalker, or battle)
- `CreatureTargetFilter`: Target creature (with optional subtype)
- `PlayerTargetFilter`: Target player or opponent
- `PermanentTargetFilter`: Any permanent (with optional type filter)

**Card Filters** (cards in hand/graveyard):
- `ArtifactCardFilter`: Artifact cards
- `CreatureCardFilter`: Creature cards
- `LandCardFilter`: Land cards
- `AnyCardFilter`: Any card

### Zone & Duration System

**Zones** (Rule 400):
```go
const (
    ZoneLibrary
    ZoneHand
    ZoneBattlefield
    ZoneGraveyard
    ZoneStack
    ZoneExile
    ZoneCommand
    ZoneOutside
)
```

**Durations** (Rule 611 - Continuous Effects):
```go
const (
    DurationUntilEndOfTurn      // Most common
    DurationPermanent           // Lasts forever
    DurationWhileOnBattlefield  // While source is on battlefield
    DurationUntilEndOfCombat    // Until end of combat
    DurationUntilYourNextTurn   // Until your next turn
    DurationWhileInGraveyard    // While in graveyard
    DurationWhileInHand         // While in hand
    DurationWhileInExile        // While in exile
)
```

### Layer System (Rule 613)

Continuous effects apply in specific layers to determine interaction order:
```go
const (
    LayerCopyEffects             // Layer 1: Clone effects
    LayerControlChanging         // Layer 2: Control change
    LayerTextChanging            // Layer 3: Text modification
    LayerTypeChanging            // Layer 4: Type changes
    LayerColorChanging           // Layer 5: Color changes
    LayerAbilityAddingRemoving   // Layer 6: Ability grants
    LayerPowerToughnessEffects   // Layer 7: P/T modification
)
```

### Attachment System

**Auras & Equipment** (Rules 303.4, 301.5):
- `EnchantAbility`: "Enchant creature", "Enchant permanent", etc.
- `EquipAbility`: "Equip {cost}" activated ability
- `AttachEffect`: Handles attachment resolution
- `GainAbilityAttachedEffect`: Grants abilities to attached permanent (accepts full Ability objects for complex abilities)
- `TapEnchantedEffect` / `UntapEnchantedEffect`: Tap/untap attached permanent
- `BoostEnchantedEffect` / `BoostEquippedEffect`: Modify P/T

### Complex Ability Example

**Psychic Overload** - Demonstrates nested abilities:
```go
// Enchant permanent
enchant := abilities.NewEnchantAbility(
    cardID,
    abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter()),
)

// ETB trigger: tap enchanted permanent
etb := abilities.NewTriggeredAbilityBuilder(cardID).
    SetTrigger(abilities.NewEntersBattlefieldTrigger(cardID)).
    AddEffect(abilities.NewTapEnchantedEffect()).
    Build()

// Static: Enchanted permanent doesn't untap
dontUntap := abilities.NewSimpleStaticAbility(cardID, abilities.ZoneBattlefield).
    AddEffect(abilities.NewDontUntapInControllersUntapStepEnchantedEffect()).
    Build()

// Static: Grant complex activated ability to enchanted permanent
grantedAbility := abilities.NewActivatedAbilityBuilder(cardID).
    AddCost(abilities.NewDiscardTargetCost(2, abilities.NewArtifactCardFilter())).
    AddEffect(abilities.NewUntapSourceEffect()).
    Build()

grantAbility := abilities.NewSimpleStaticAbility(cardID, abilities.ZoneBattlefield).
    AddEffect(abilities.NewGainAbilityAttachedEffect(
        grantedAbility,
        abilities.AttachmentTypeAura,
        abilities.DurationWhileOnBattlefield,
        "Enchanted permanent has \"Discard two artifact cards: Untap this permanent.\"",
    )).
    Build()
```

### Implementation Status

**✅ Fully Implemented**:
- Base ability interface & types
- Spell abilities with targeting
- Activated abilities with costs
- Triggered abilities with event system
- Static abilities with zone awareness
- 40+ effects (damage, draw, destroy, tokens, counters, mill, bounce, exile, etc.)
- Cost system (mana, tap, sacrifice, discard, life)
- Target filters & card filters
- Attachment system (Auras & Equipment)
- Duration & layer system
- Keyword abilities (Flying, Trample, Haste, Vigilance, etc.)

**⚠️ Stub/TODO**:
- Full layer system application (currently stubs)
- Integration with turn structure and state-based actions
- Comprehensive game event system
- Mana pool management
- Complete rules engine for resolution

### Key Files

- `internal/game/abilities/ability.go`: Core interfaces
- `internal/game/abilities/activated.go`: Activated abilities
- `internal/game/abilities/triggered.go`: Triggered abilities & events
- `internal/game/abilities/static.go`: Static abilities, zones, durations, layers
- `internal/game/abilities/spell.go`: Spell abilities
- `internal/game/abilities/effects.go`: Basic effects & Duration
- `internal/game/abilities/enchanted_effects.go`: Aura/Equipment effects
- `internal/game/abilities/attach.go`: Attachment system
- `internal/game/abilities/costs.go`: Cost system
- `internal/game/abilities/targets.go`: Target & card filters
- `internal/game/abilities/keyword.go`: Keyword abilities
- `internal/game/abilities/counter_effects.go`: Counter manipulation
- `internal/game/abilities/token_effects.go`: Token creation
- `internal/game/abilities/grant_ability_effect.go`: Ability granting
- `internal/game/abilities/bounce.go`: Return to hand effects
- `internal/game/abilities/exile.go`: Exile effects
- `internal/game/abilities/mill.go`: Mill effects
- `internal/game/abilities/search_library.go`: Library search effects
- `internal/game/abilities/scry_surveil.go`: Scry & surveil
- `internal/game/abilities/gain_control.go`: Control-changing effects

## Key Architectural Patterns

### Manager Interface Pattern

All components (User, Session, Table, Game, etc.) follow this pattern:
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

### Builder Pattern for Abilities

Abilities are constructed using fluent builders:
```go
ability := abilities.NewSpellAbilityBuilder(cardID, manaCost).
    AddTarget(abilities.NewCreatureTarget()).
    AddEffect(abilities.NewDamageEffect(3)).
    AddEffect(abilities.NewDrawCardsEffect(2)).
    Build()
```

### Repository Pattern

Clean separation between domain logic and database:
- Repositories handle PostgreSQL operations (using pgx driver)
- Managers handle business logic
- Domain models are database-agnostic

### Registry Pattern

Pre-compiled "plugin" system (no dynamic loading):
- **Card Registry**: 30,400+ cards via `init()` registration
- **Game Type Registry**: TwoPlayerDuel, CommanderFreeForAll, etc.
- **Tournament Type Registry**: Constructed, BoosterDraft, Sealed
- **Player Type Registry**: Human, ComputerMAX, ComputerDraft

## Critical Implementation Details

### Component Initialization Order

**MUST follow this order** to avoid dependency issues:
1. Config → Logger → Database
2. Session Manager (starts background cleanup goroutine)
3. Repositories (User, Stats, Card)
4. Domain Managers (User, Room, Chat, Table, Game, Tournament, Draft)
5. External Services (Email)
6. Servers (gRPC + WebSocket)

### Session Management

Uses **lease-based expiration** (not simple timeouts):
- Each session has `LeaseUntil` timestamp
- `UpdateActivity()` extends lease by 5 minutes (configurable)
- Background goroutine cleans expired sessions
- Ping RPC keeps sessions alive

### Password Security

Uses **Argon2id** (not bcrypt):
- Parameters: time=1, memory=64MB, threads=4, keyLen=32
- Format: `$argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>`

### Rating System

Implements **Glicko-2** (not Elo):
- Default: rating=1500, deviation=350, volatility=0.06
- Handles inactivity (deviation increases over time)
- Full implementation in `internal/rating/`

### Database

- **PostgreSQL 15+** with pgx driver (pure Go)
- Connection pooling with configurable limits
- Migrations via golang-migrate
- Environment variable overrides: `DB_PASSWORD`, etc.

## Configuration

**Main config:** `mage-server-go/config/config.yaml`

**Environment variable overrides:**
```bash
export DB_PASSWORD="secure_password"
export SMTP_HOST="smtp.gmail.com"
export SMTP_USER="user@gmail.com"
export SMTP_PASSWORD="app_password"
export MAILGUN_DOMAIN="mg.example.com"
export MAILGUN_API_KEY="key-xxx"
export VITE_GRPC_SERVER_URL="http://localhost:17171"  # For web client
```

**Key settings:**
- `server.grpc.address`: gRPC server address (default: 0.0.0.0:17171)
- `server.websocket.address`: WebSocket address (default: 0.0.0.0:17179)
- `server.lease_period`: Session timeout (default: 5m)
- `database.max_conns`: PostgreSQL connection pool size (default: 25)
- `auth.mode`: "optional" or "required"
- `mail.provider`: "smtp", "mailgun", or "none"

## Web Client Routing

**Public routes:**
- `/` - Landing page
- `/login` - Login with guest option
- `/register` - User registration

**Protected routes** (require authentication):
- `/lobby` - Main game lobby
- `/game/[id]` - Active game view
- `/table/[id]` - Table/match setup
- `/decks` - Deck builder
- `/profile` - User profile

Authentication is handled via Svelte stores (`lib/stores/auth.ts`) with JWT tokens stored in localStorage.

## Test Coverage

**82 test files** covering:
- **Unit tests** (36 tests): Auth (Argon2id), Session (lease expiration), Rating (Glicko-2), Draft, Tournament, Abilities
- **Integration tests** (21 tests): End-to-end workflows, authentication flows, session management, game flows

**Critical test areas:**
- `internal/auth/`: Password hashing and validation
- `internal/session/`: Lease-based expiration, concurrent access
- `internal/rating/`: Glicko-2 rating calculations
- `internal/draft/`: Pick tracking, booster pack passing
- `internal/tournament/`: Swiss pairing, elimination brackets, tiebreakers
- `internal/game/abilities/`: Effect resolution, token creation, counters

## Common Development Workflows

### Adding a New Card

1. **Auto-generate from Java source:**
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go
python scripts/transpile_cards.py --card-name "Card Name"
```

2. **Or manually create:**
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go/internal/game/cards/manual
# Create new .go file with card implementation
```

3. **Test the card:**
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go
go test ./internal/game/cards/...
```

### Adding a New Ability Type

1. Create ability in `internal/game/abilities/`
2. Add builder methods if needed
3. Update `AbilityType` constants
4. Write unit tests in `*_test.go`
5. Update card transpiler if Java mapping needed

### Adding a New API Endpoint

1. **Define in protobuf:**
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go/api/proto
# Edit relevant .proto file
```

2. **Regenerate code:**
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go
make proto
```

3. **Implement handler in server:**
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go/internal/server
# Implement gRPC method
```

4. **Add web client call:**
```bash
cd /Users/aron/dev/opensource/mage/mage-client-web/src/lib/grpc
# Add client method
```

### Running Integration Tests

```bash
cd /Users/aron/dev/opensource/mage/mage-server-go

# Ensure database is running
docker-compose up -d postgres
make migrate-up

# Run integration tests
make test-integration

# Or manually with verbose output
go test -v -tags=integration ./internal/integration/...
```

## Troubleshooting

### "proto files not found" error
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go
make proto  # Must run before building
```

### gRPC/WebSocket servers won't start
Check that protobuf generation completed successfully and imports are correct in `cmd/server/main.go`

### Database connection failed
```bash
# Ensure PostgreSQL is running
docker-compose up -d postgres

# Check connection settings in config.yaml
# Verify DB_PASSWORD environment variable if set
```

### Card not found in registry
Cards must be registered via `init()` functions. Check that:
1. Card file exists in `internal/game/cards/generated/` or `manual/`
2. File has proper `init()` registration
3. Package is imported (Go may optimize out unused packages)

### Tests failing with "database: relation does not exist"
```bash
cd /Users/aron/dev/opensource/mage/mage-server-go
make migrate-up  # Apply all migrations
```

### Web client can't connect to server
1. Check `VITE_GRPC_SERVER_URL` environment variable
2. Verify server is running: `curl http://localhost:17171`
3. Check browser console for CORS errors
4. Ensure protobuf types match between client and server
