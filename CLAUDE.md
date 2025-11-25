# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **port of XMage (Magic: The Gathering game simulator) from Java to Go and Svelte**. The goal is to achieve everything the original Java XMage offers and more, with 100% rule coverage and full test coverage for all 28,000+ unique cards.

**Key Principles:**
- Never create summary or example documentation
- Build the MVP first before optimizing structure and architecture
- 100% MTG rule coverage is required
- Full test coverage is required
- All cards and spells must be supported
- Always use absolute paths for cd commands

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

1. Generate protobuf code before building:
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

Always style with CSS classes.
Always keep code as DRY as possible by reusing and extracting to well scoped typescript files.

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

## Docker Compose

```bash
cd /Users/aron/dev/opensource/mage

# Start full stack
docker compose up

# Start specific services
docker compose up postgres
docker compose up mage-server


# Rebuild after code changes
docker compose up --build
```

## Key Files

among other things:

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

### Session Management

Uses **lease-based expiration** (not simple timeouts):
- Each session has `LeaseUntil` timestamp
- `UpdateActivity()` extends lease by 5 minutes (configurable)
- Background goroutine cleans expired sessions
- Ping RPC keeps sessions alive

### Database

- **PostgreSQL 15+** with pgx driver (pure Go)
- Connection pooling with configurable limits
- Migrations via golang-migrate
- Environment variable overrides: `DB_PASSWORD`, etc.

## Configuration

**Main config:** `mage-server-go/config/config.yaml`

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


