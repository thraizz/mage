# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go reimplementation of the MAGE (Magic Another Game Engine) server using a **hybrid gRPC + WebSocket architecture**:
- **gRPC**: All 60+ request/response RPC methods (authentication, room management, table/game/tournament operations)
- **WebSocket**: Server-to-client push events (game updates, chat messages, callbacks)
- **Protocol Buffers**: Type-safe serialization for both protocols

The server maintains API compatibility with the existing Java MAGE server to support existing clients.

## Essential Commands

```bash
# Build and run
make build              # Build server binary to bin/mage-server
make run                # Build and run with config/config.yaml

# Testing
make test               # Run all unit tests (57 tests)
make test-integration   # Run integration tests

# Code quality
make fmt                # Format with gofmt + goimports
make lint               # Run golangci-lint
make vet                # Run go vet

# Protocol Buffers (REQUIRED before first build)
make proto              # Generate Go code from .proto files

# Database
make migrate-up         # Apply all migrations
make migrate-down       # Rollback one migration
make migrate-create NAME=feature_name  # Create new migration files

# Development tools
make tools              # Install protoc-gen-go, golangci-lint, goimports
make deps               # Install Go dependencies
```

### Running Specific Tests

```bash
# Single test
go test -v -run TestSessionCreation ./internal/session

# Package tests
go test -v ./internal/rating

# With coverage
go test -v -coverprofile=coverage.out ./internal/rating
go tool cover -html=coverage.out
```

## Architecture & Design Patterns

### Component Initialization Order (cmd/server/main.go)

**CRITICAL**: Components must be initialized in this exact order to avoid dependency issues:

```go
1. Config → Logger → Database
2. Session Manager (starts background cleanup goroutine)
3. Repositories (User, Stats, Card)
4. Domain Managers (User, Room, Chat, Table, Game, Tournament, Draft)
5. External Services (Email client)
6. Servers (gRPC + WebSocket) - currently commented pending protobuf generation
```

### Manager Interface Pattern

All components follow a **Manager interface pattern**:

```go
// 1. Define interface in internal/<domain>/manager.go
type Manager interface {
    Create(...) (*Entity, error)
    Get(id string) (*Entity, bool)
    Update(...) error
    Remove(id string)
}

// 2. Implement as private struct with logger
type manager struct {
    data   map[string]*Entity
    mu     sync.RWMutex
    logger *zap.Logger
}

// 3. Constructor returns interface
func NewManager(logger *zap.Logger) Manager {
    return &manager{
        data:   make(map[string]*Entity),
        logger: logger,
    }
}
```

**Key managers and their responsibilities**:
- `session.Manager`: Session lifecycle with lease-based expiration (internal/session/manager.go)
- `user.Manager`: User authentication, validation, registration (internal/user/manager.go)
- `room.Manager`: Lobby and room system (internal/room/manager.go)
- `table.Manager`: Table state and player management (internal/table/manager.go)
- `game.Manager`: Game controller interface for game engine integration (internal/game/manager.go)
- `tournament.Manager`: Swiss pairing algorithm, round management (internal/tournament/manager.go)
- `draft.Manager`: Booster passing and pick handling (internal/draft/manager.go)
- `chat.Manager`: Chat rooms and message history (internal/chat/manager.go)

### Session Management with Lease Mechanism

Sessions use a **lease-based expiration** instead of simple timeouts:

```go
type Session struct {
    ID           string
    UserID       string
    IsAdmin      bool
    LeaseUntil   time.Time      // Key field: lease expiration timestamp
    CallbackChan chan interface{} // For WebSocket push events
}
```

**How it works**:
1. Each session has a `LeaseUntil` timestamp (not just LastActivity)
2. `UpdateActivity()` extends the lease by `leasePeriod` (default 5 minutes)
3. Background goroutine in `session.Manager` cleans up expired sessions
4. Ping RPC method keeps sessions alive by calling `UpdateActivity()`

**Location**: `internal/session/session.go:102-115` (lease logic), `internal/session/manager.go:87-104` (cleanup)

### gRPC Interceptor Chain

The gRPC server uses **ordered interceptors** for cross-cutting concerns (internal/server/interceptors.go):

```
Request → RecoveryInterceptor → LoggingInterceptor → SessionValidationInterceptor
       → AdminInterceptor → MetricsInterceptor → Handler
```

**Order matters**:
1. **RecoveryInterceptor**: Panic recovery with stack traces (outermost - catches everything)
2. **LoggingInterceptor**: Request/response logging with duration
3. **SessionValidationInterceptor**: Validates session ID in metadata (fails early for invalid sessions)
4. **AdminInterceptor**: Checks admin privileges for admin methods
5. **MetricsInterceptor**: Prometheus metrics (request count, latency)

### WebSocket Callback System

Real-time event delivery from server → client:

```
1. Client connects to /ws?sessionId=<id>
2. Server validates session and retrieves CallbackChan from Session
3. Server goroutine listens on CallbackChan
4. Any component pushes events to session.CallbackChan
5. WebSocket server marshals to JSON and sends to client
```

**Event types** (api/proto/mage/v1/websocket.proto): Chat messages, game updates, tournament updates, draft updates, table events

**Location**: `internal/server/websocket.go:93-127` (connection handler), `internal/session/session.go:47` (CallbackChan field)

### Authentication & Security

**Password Hashing** (internal/auth/password.go):
- **Argon2id** (not bcrypt) with parameters: time=1, memory=64MB, threads=4, keyLen=32
- Format: `$argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>`

**Password Reset** (internal/auth/token.go):
- 6-digit numeric tokens with TTL (default 1 hour)
- In-memory token store (TokenStore)
- Single-use tokens (consumed on use)

**Validation** (internal/user/validator.go:17-67):
- Username: 3-16 chars, alphanumeric + underscore/hyphen
- Password: 8+ chars, requires uppercase, lowercase, digit, special char
- Email: RFC 5322 format

### Database Layer (Repository Pattern)

Uses **pgx** (not lib/pq) with connection pooling:

```go
type DB struct {
    Pool *pgxpool.Pool  // pgx connection pool
}
```

**Repositories**:
- `repository.UserRepository`: User CRUD operations (internal/repository/users.go)
- `repository.StatsRepository`: Glicko rating system, win/loss tracking (internal/repository/stats.go)
- `repository.CardRepository`: Card data with caching layer (internal/repository/cards.go)

**Connection pooling** (config.yaml):
```yaml
database:
  max_conns: 25
  min_conns: 5
  max_conn_lifetime: 1h
```

### Rating System (Glicko-2)

**Implementation** (internal/rating/glicko.go):
- Default: rating=1500, deviation=350, volatility=0.06
- Updates after each match with opponent rating/deviation
- Handles inactivity (deviation increases over time)

**Key functions**:
- `g(φ)`: Deviation impact on rating change (internal/rating/glicko.go:113-120)
- `E(μ, μj, φj)`: Expected score against opponent (internal/rating/glicko.go:123-130)
- `Δ`: Performance-based rating change (internal/rating/glicko.go:133-147)

### Tournament System (Swiss Pairing)

**Algorithm** (internal/tournament/manager.go:154-231):
1. Players paired by score (highest vs highest in each score bracket)
2. Bye handling for odd players (lowest-rated unpaired player gets bye)
3. Match results tracked (win/loss/draw)
4. Standings calculation with tiebreakers

### Draft System

**Mechanics** (internal/draft/manager.go):
- Configurable sets, packs per player, cards per pack
- Pick tracking per player
- Booster passing direction alternates by pack
- Completion detection: `CurrentPack > NumPacks`

## Protocol Buffers

### Generating Protobuf Code

**CRITICAL**: Run this before first build or after any `.proto` changes:
```bash
make proto
```

This runs `scripts/generate_proto.sh` which:
1. Generates Go structs from `.proto` files in `api/proto/mage/v1/`
2. Outputs to `pkg/proto/mage/v1/`
3. Generates gRPC service stubs

**Proto organization**:
- `server.proto`: Main MageServer service (60+ RPC methods)
- Domain-specific: `auth.proto`, `room.proto`, `table.proto`, `game.proto`, `tournament.proto`, `draft.proto`
- `chat.proto`, `admin.proto`: Chat and admin messages
- `models.proto`: Shared data models
- `websocket.proto`: WebSocket event types

**Import path** in Go code:
```go
import pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
```

### Current Protobuf Status

**IMPORTANT**: The server structure is complete but gRPC/WebSocket servers are **not started** until protobuf code is generated. See commented code in `cmd/server/main.go:151-200` for server initialization that activates after `make proto`.

## Common Workflows

### Adding a New RPC Method

1. Add method to `api/proto/mage/v1/server.proto` service definition
2. Define request/response messages in appropriate proto file
3. Run `make proto` to regenerate code
4. Implement method in `internal/server/grpc.go`
5. Add business logic to domain managers if needed

### Adding a New WebSocket Event

1. Define event message in `api/proto/mage/v1/websocket.proto`
2. Run `make proto`
3. Send event via `session.CallbackChan` from any component:
   ```go
   if sess, ok := s.sessionMgr.GetSession(sessionID); ok {
       sess.CallbackChan <- &pb.GameUpdateEvent{...}
   }
   ```

### Adding a Database Migration

```bash
make migrate-create NAME=add_feature
# Edit migrations/NNNN_add_feature.up.sql
# Edit migrations/NNNN_add_feature.down.sql
make migrate-up
```

### Adding a New Manager

1. Define interface in `internal/<domain>/manager.go`
2. Implement struct with logger and data structures
3. Add to initialization in `cmd/server/main.go` (respect initialization order)
4. Wire into gRPC server methods as needed

## Code Style & Conventions

**Critical conventions for this codebase**:

1. **Interfaces before implementations**: Define Manager interfaces, implement as private structs
2. **Constructor pattern**: `NewManager()` functions for all components
3. **Context propagation**: All I/O operations accept `context.Context`
4. **Structured logging**: Use `zap.Logger` with typed fields
   ```go
   logger.Info("user connected",
       zap.String("user", userName),
       zap.String("session", sessionID))
   ```
5. **Error wrapping**: Use `fmt.Errorf("context: %w", err)` for error chains
6. **Mutex discipline**: Use `sync.RWMutex` for read-heavy data structures
   ```go
   m.mu.RLock()         // For reads
   defer m.mu.RUnlock()

   m.mu.Lock()          // For writes
   defer m.mu.Unlock()
   ```

## Testing

### Test Organization

- **Unit tests** (36 tests): `*_test.go` next to implementation
- **Integration tests** (21 tests): `internal/integration/*_test.go`
- Integration tests use `-tags=integration` flag

**Test coverage** (57 tests total):
- ✅ Authentication (Argon2id) - 4 tests
- ✅ Session management - 6 tests
- ✅ Rating (Glicko-2) - 9 tests
- ✅ Draft mechanics - 8 tests
- ✅ Tournament (Swiss) - 9 tests
- ✅ Integration flows - 21 tests

**Components without tests** (require DB/network):
- User Manager (needs database)
- Repository Layer (needs database)
- Server/gRPC (needs protobuf generation)
- Table/Game/Room/Chat Managers
- Email Service
- WebSocket Server

### Writing Tests

**Unit test pattern**:
```go
func TestFeature(t *testing.T) {
    // Setup
    logger := zap.NewNop()
    mgr := NewManager(logger)

    // Execute
    result, err := mgr.DoSomething()

    // Assert
    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

**Integration test pattern** (internal/integration/*_test.go):
```go
func TestCompleteFlow(t *testing.T) {
    // Setup managers
    sessionMgr := session.NewManager(5 * time.Second)
    userMgr := user.NewManager(/* ... */)

    // Test complete workflow
    sess := sessionMgr.CreateSession(id, host)
    user := userMgr.Register(/* ... */)
    // ... test full flow
}
```

## Configuration

Configuration loaded from YAML with environment variable overrides via **Viper**.

**Key settings** (config/config.yaml):
```yaml
server:
  grpc:
    address: "0.0.0.0:50051"
  websocket:
    address: "0.0.0.0:50052"
  max_sessions: 1000
  lease_period: 5m          # Session lease duration

database:
  host: "localhost"
  port: 5432
  user: "mage"
  password: "mage"          # Override with DB_PASSWORD env var

auth:
  mode: "db"                # or "none" for testing
  password_reset_token_ttl: 1h

mail:
  provider: "smtp"          # or "mailgun" or "none"
```

**Environment variable overrides**:
- `DB_PASSWORD`: Database password
- `SMTP_HOST`, `SMTP_USER`, `SMTP_PASSWORD`: SMTP settings
- `MAILGUN_DOMAIN`, `MAILGUN_API_KEY`: Mailgun settings

## Dependencies

**Core**:
- `google.golang.org/grpc`: gRPC framework
- `google.golang.org/protobuf`: Protocol buffers
- `github.com/gorilla/websocket`: WebSocket server
- `github.com/jackc/pgx/v5`: PostgreSQL driver (not lib/pq)
- `github.com/spf13/viper`: Configuration management
- `go.uber.org/zap`: Structured logging (not logrus)
- `golang.org/x/crypto`: Argon2id password hashing
- `github.com/google/uuid`: UUID generation

## Known Issues & Next Steps

### Current Limitations

1. **No protobuf code generated**: Run `make proto` before first build
2. **gRPC/WebSocket servers not started**: Uncomment code in `main.go` after protobuf generation
3. **Game engine interface not implemented**: Game controller is a stub awaiting integration
4. **No database integration tests**: Requires testcontainers setup
5. **Card repository not populated**: Requires card data import

### Implementation Priority

1. Generate protobuf code (`make proto`)
2. Implement remaining gRPC methods (currently stubbed)
3. Integrate game engine (external library or custom)
4. Add database integration tests with testcontainers
5. Performance testing (load testing, profiling)

## Additional Resources

- **Implementation Plan**: See `../GO_SERVER_IMPLEMENTATION.md` for complete 28-week plan
- **Test Coverage Report**: See `TEST_COVERAGE.md` for detailed test breakdown
- **Protocol Buffers**: `api/proto/mage/v1/*.proto` for API definitions
- **Java Server Reference**: https://github.com/magefree/mage
