# MAGE Server Go Implementation - Definitive Plan

## Overview
This document outlines the implementation of a new MAGE (Magic Another Game Engine) server in Go using **gRPC for RPC calls** and **WebSocket for server push events**. The server will provide the same API as the existing Java server, allowing existing clients to connect with minimal changes.

## Current Progress (Nov 2025)
- Core scaffolding, proto definitions, migrations, and configuration loader are in place and exercised by `go test`.
- Session management (including preference storage and cleanup) plus Argon2id-based authentication and token flow are implemented with unit coverage.
- gRPC bootstrap now runs end-to-end with interceptors, WebSocket bridge, admin login, and initial RPC coverage for auth, server info, and lobby/table creation.
- User, stats, and card repositories/managers are active; additional controllers and advanced RPC categories remain planned work.
- MageEngine now drives gameplay for the Go server: MatchStart wires into the real engine, and `GameGetView` returns full zone, prompt, and message state for players and watchers; integration tests cover the richer views.
- Rules scaffolding initiated in Go: a `TurnManager` now models the full MTG phase/step progression and feeds priority/turn updates back into `GameGetView`.
- A `StackManager` now tracks stack items, resolves spells after full pass cycles, and surfaces battlefield transitions/log messages via `GameGetView`.
- Event bus and trigger manager deliver triggered abilities onto the stack with life-gain sample coverage; integration tests assert ordering and resolution behaviour.
- Continuous effect layer system applies basic power/toughness buffs during normalization, laying groundwork for full layer 1-7 implementation.

## Architecture Decision Summary

### Protocol Architecture: gRPC + WebSocket Hybrid
- **gRPC**: All request/response RPC methods (60+ methods from MageServer interface)
- **WebSocket**: Server-to-client push events (callbacks from ClientCallback system)
- **Rationale**: gRPC provides type-safe, efficient RPC with built-in Protocol Buffers. WebSocket handles real-time push notifications. This combination requires minimal client changes (swap JBoss Remoting client for gRPC + WebSocket client).

### Database: PostgreSQL
- **Choice**: PostgreSQL 15+ with pgx driver (pure Go)
- **Rationale**: Best performance, scalability, and tooling. Full-text search built-in. Production-ready.

### Password Hashing: Argon2id
- **Choice**: Argon2id (modern, secure, Go-native)
- **Rationale**: Industry-standard secure password hashing

### Session Storage: In-Memory (with Redis-ready interface)
- **Initial**: In-memory map with sync.RWMutex
- **Production**: Redis for distributed deployment (interface-based design allows easy swap)
- **Rationale**: Simplicity for initial development, scalability path for production

### Plugin System: Pre-compiled Registry Pattern
- **Choice**: All game types compiled into binary, registry-based instantiation
- **Rationale**: Go doesn't support cross-platform dynamic loading. Pre-compiled approach is reliable, performant, and simplifies deployment.
- **Extensibility**: New game types added via PR and recompilation (matches OSS workflow)

### Configuration: Viper with YAML
- **Choice**: Viper for YAML config with environment variable overrides
- **Rationale**: Industry standard, supports complex config structures

### Logging: Zap
- **Choice**: Uber's Zap (structured, high-performance)
- **Rationale**: Best-in-class performance, structured logging for debugging complex game states

### Email: Dual Provider Support
- **SMTP**: gomail.v2 for direct SMTP
- **Mailgun**: mailgun-go SDK
- **Configuration-driven selection**

### Cache Layer: In-memory with groupcache
- **Card data**: groupcache for distributed caching
- **Rationale**: Cards are read-heavy, immutable data perfect for caching

## Client Integration Strategy

### Minimal Client Changes Required

**Java Client Changes (Single PR)**:
1. Replace `JBoss Remoting TransporterClient` with `gRPC MageServerClient`
2. Replace `InvokerCallbackHandler` with `WebSocket Client`
3. Update serialization (JBoss Serialization → Protocol Buffers)
4. Keep all UI, game logic, views unchanged

**Implementation**:
```java
// Before (JBoss Remoting)
TransporterClient client = new TransporterClient(...);
client.setSessionId(sessionId);
MageServer server = (MageServer) client.getTarget();
server.connectUser(userName, password, sessionId, ...);

// After (gRPC + WebSocket)
MageServerGrpc.MageServerBlockingStub server = MageServerGrpc.newBlockingStub(channel);
ConnectUserResponse response = server.connectUser(
    ConnectUserRequest.newBuilder()
        .setUserName(userName)
        .setPassword(password)
        .setSessionId(sessionId)
        .build()
);

WebSocketClient wsClient = new WebSocketClient(wsUri);
wsClient.addMessageHandler(this::handleServerCallback);
```

**Compatibility Shim**: Create a Java adapter class that wraps gRPC/WebSocket to match the original `MageServer` interface signature, minimizing code changes in client.

## Project Structure

```
mage-server-go/
├── cmd/
│   └── server/
│       └── main.go                      # Server entry point
│
├── api/
│   └── proto/
│       └── mage/
│           └── v1/
│               ├── server.proto         # Main RPC service definition
│               ├── auth.proto           # Authentication messages
│               ├── room.proto           # Room/lobby messages
│               ├── table.proto          # Table management messages
│               ├── game.proto           # Game execution messages
│               ├── tournament.proto     # Tournament messages
│               ├── draft.proto          # Draft messages
│               ├── chat.proto           # Chat messages
│               ├── admin.proto          # Admin messages
│               ├── models.proto         # Shared data models (TableView, GameView, etc.)
│               └── websocket.proto      # WebSocket callback messages
│
├── internal/
│   ├── server/
│   │   ├── grpc.go                      # gRPC server implementation
│   │   ├── websocket.go                 # WebSocket server for callbacks
│   │   ├── interceptors.go              # gRPC interceptors (auth, logging, metrics)
│   │   └── middleware.go                # Session validation middleware
│   │
│   ├── session/
│   │   ├── session.go                   # Session struct and methods
│   │   ├── manager.go                   # SessionManager interface + impl
│   │   ├── store.go                     # Session storage (in-memory or Redis)
│   │   └── lease.go                     # Lease/ping mechanism
│   │
│   ├── user/
│   │   ├── user.go                      # User domain model
│   │   ├── manager.go                   # UserManager interface + impl
│   │   ├── repository.go                # User database operations
│   │   └── validator.go                 # Username/password validation
│   │
│   ├── auth/
│   │   ├── password.go                  # Password hashing (Argon2id)
│   │   ├── token.go                     # Password reset token generation
│   │   └── service.go                   # Auth service (registration, reset)
│   │
│   ├── table/
│   │   ├── controller.go                # TableController
│   │   ├── manager.go                   # TableManager
│   │   ├── state.go                     # Table state machine
│   │   └── seat.go                      # Player seat management
│   │
│   ├── game/
│   │   ├── controller.go                # GameController
│   │   ├── manager.go                   # GameManager
│   │   ├── player_session.go            # GameSessionPlayer
│   │   ├── watcher_session.go           # GameSessionWatcher
│   │   ├── worker.go                    # Game execution worker pool
│   │   ├── replay.go                    # Replay system
│   │   └── view.go                      # GameView generation
│   │
│   ├── tournament/
│   │   ├── controller.go                # TournamentController
│   │   ├── manager.go                   # TournamentManager
│   │   ├── session.go                   # TournamentSession
│   │   ├── pairing.go                   # Swiss/elimination pairing
│   │   └── types.go                     # Tournament type registry
│   │
│   ├── draft/
│   │   ├── controller.go                # DraftController
│   │   ├── manager.go                   # DraftManager
│   │   ├── session.go                   # DraftSession
│   │   ├── cube.go                      # Cube draft logic
│   │   └── booster.go                   # Booster pack generation
│   │
│   ├── room/
│   │   ├── room.go                      # Room interface + base impl
│   │   ├── games_room.go                # GamesRoomImpl (main lobby)
│   │   └── manager.go                   # GamesRoomManager
│   │
│   ├── chat/
│   │   ├── chat.go                      # ChatSession
│   │   ├── manager.go                   # ChatManager
│   │   ├── message.go                   # ChatMessage handling
│   │   └── sanitizer.go                 # HTML sanitization
│   │
│   ├── repository/
│   │   ├── db.go                        # Database connection management
│   │   ├── cards.go                     # CardRepository
│   │   ├── users.go                     # AuthorizedUserRepository
│   │   ├── stats.go                     # UserStatsRepository
│   │   └── records.go                   # TableRecordRepository
│   │
│   ├── rating/
│   │   ├── glicko.go                    # Glicko rating implementation
│   │   └── calculator.go                # Rating calculation
│   │
│   ├── mail/
│   │   ├── client.go                    # MailClient interface
│   │   ├── smtp.go                      # SMTP implementation
│   │   ├── mailgun.go                   # Mailgun implementation
│   │   └── templates.go                 # Email templates
│   │
│   ├── plugin/
│   │   ├── registry.go                  # Plugin registry pattern
│   │   ├── game_types.go                # Game type definitions
│   │   ├── tournament_types.go          # Tournament type definitions
│   │   └── player_types.go              # Player type definitions
│   │
│   ├── config/
│   │   ├── config.go                    # Config structs
│   │   ├── loader.go                    # Viper-based config loading
│   │   └── validator.go                 # Config validation
│   │
│   ├── cache/
│   │   └── cards.go                     # Card cache implementation
│   │
│   └── util/
│       ├── compress.go                  # Compression utilities (for Protocol Buffers)
│       ├── uuid.go                      # UUID handling
│       └── errors.go                    # Error types
│
├── pkg/
│   ├── proto/                           # Generated Go protobuf code
│   │   └── mage/
│   │       └── v1/
│   │           ├── server.pb.go
│   │           ├── server_grpc.pb.go
│   │           └── ...
│   │
│   └── models/                          # Shared models (if needed outside internal/)
│       └── version.go                   # Version handling
│
├── migrations/                          # SQL migration files
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_cards_table.up.sql
│   ├── 002_create_cards_table.down.sql
│   ├── 003_create_stats_table.up.sql
│   ├── 003_create_stats_table.down.sql
│   ├── 004_create_table_records.up.sql
│   └── 004_create_table_records.down.sql
│
├── config/
│   ├── config.yaml                      # Default server config
│   ├── config.example.yaml              # Example config with comments
│   └── config.dev.yaml                  # Development overrides
│
├── scripts/
│   └── generate_proto.sh                # Regenerate protobuf code
│
├── test/
│   ├── integration/                     # Integration tests
│   │   ├── session_test.go
│   │   ├── game_flow_test.go
│   │   └── tournament_test.go
│   └── testdata/                        # Test fixtures
│       └── cards.json
│
├── .proto                               # Protobuf generation config
├── Dockerfile                           # Multi-stage Docker build
├── docker-compose.yml                   # Local dev environment (server + postgres)
├── go.mod
├── go.sum
├── Makefile                             # Build targets
└── README.md                            # Go server documentation
```

## Phase-by-Phase Implementation

### Phase 1: Foundation & Protobuf Definitions (Week 1-3)

#### 1.1 Project Initialization
- [x] Initialize Go module: `go mod init github.com/magefree/mage-server-go`
- [x] Set up Makefile with common targets (build, test, proto, run)
- [x] Configure protobuf tooling (buf or protoc)
- [ ] Set up CI/CD pipeline (GitHub Actions)

#### 1.2 Protocol Buffer Definitions
**Map all 60+ RPC methods to gRPC services**

**Action Items**:
- [ ] Define all 60+ request/response message types in separate .proto files
- [ ] Define data models (TableView, GameView, MatchOptions, etc.) in models.proto
- [ ] Define WebSocket event messages in websocket.proto
- [ ] Generate Go code: `buf generate` or `make proto`
- [ ] Commit generated code to git

#### 1.3 Database Schema

**Action Items**:
- [x] Create all migration files (up and down)
- [ ] Set up golang-migrate: `brew install golang-migrate`
- [ ] Create initial data seed scripts (cards, expansions, etc.)
- [ ] Test migrations against PostgreSQL

#### 1.4 Configuration Management

**Action Items**:
- [x] Implement config loading with Viper
- [x] Support environment variable overrides
- [x] Add config validation
- [x] Create example configs for dev/staging/prod
- [ ] Document all config options

### Phase 2: Core Infrastructure (Week 4-6)

#### 2.1 Database Layer

**Action Items**:
- [x] Implement connection pooling with pgx
- [ ] Create repository interfaces for each entity
- [x] Implement user repository (CRUD + queries)
- [x] Implement card repository with full-text search
- [x] Implement stats repository with Glicko rating queries
- [ ] Implement table records repository
- [ ] Add database health check
- [ ] Add query logging for debugging
- [ ] Write repository unit tests with test database

#### 2.2 Session Management

**Action Items**:
- [x] Implement Session struct with all methods
- [x] Implement SessionManager with in-memory storage
- [x] Add session expiration cleanup goroutine
- [x] Implement session restoration for reconnects
- [x] Add concurrent request locking per session
- [x] Write session manager tests
- [ ] (Future) Add Redis-backed session store interface

#### 2.3 Authentication & Security

**Action Items**:
- [x] Implement Argon2id password hashing
- [x] Implement password reset token generation (6-digit)
- [x] Add token storage (in-memory cache with TTL)
- [x] Write password hashing tests
- [ ] Write auth service tests

### Phase 3: gRPC Server Implementation (Week 7-10)

#### 3.1 gRPC Server Bootstrap

**Action Items**:
- [x] Implement main entry point with graceful shutdown
- [x] Set up gRPC server with keepalive
- [x] Implement session validation interceptor
- [x] Implement logging interceptor
- [x] Implement panic recovery interceptor
- [ ] Implement metrics interceptor (Prometheus)
- [ ] Add health check service

#### 3.2 RPC Method Implementation (60+ methods)

**Action Items**:
- [x] Implement all Authentication methods (6 methods)
- [x] Implement all Server Info methods (3 methods)
- [x] Implement all Room/Lobby methods (5 methods)
- [x] Implement all Table Management methods (10 methods)
- [x] Implement all Deck Management methods (2 methods)
- [ ] Implement all Game Execution methods (15 methods) _(match lifecycle now backed by MageEngine with zones/prompts/messages exposed via GameGetView; remaining advanced prompts/mana/replay hooks pending)_
- [ ] Implement all Draft methods (6 methods) _(initial join/pick/mark/booster loaded/quit handlers implemented)_
- [ ] Implement all Tournament methods (4 methods) _(initial join/start/quit/find handlers implemented)_
- [x] Implement all Chat methods (7 methods)
- [ ] Implement all Replay methods (6 methods)
- [ ] Implement all Admin methods (9 methods)
- [ ] Add error handling and validation for each method
- [ ] Write integration tests for each category

### Phase 4: WebSocket Server (Week 11-12)

#### 4.1 WebSocket Implementation

**Action Items**:
- [x] Implement WebSocket server with Gorilla WebSocket
- [x] Handle WebSocket upgrade from HTTP
- [x] Implement session validation for WebSocket connections
- [x] Forward ServerEvent messages from session.CallbackChan to WebSocket
- [x] Implement ping/pong keep-alive
- [ ] Add reconnection handling (resume from last message ID)
- [ ] Implement compression for large messages
- [ ] Write WebSocket integration tests
- [ ] Document WebSocket protocol for client developers

### Phase 5: Business Logic - Managers & Controllers (Week 13-18)

#### 5.1 User Management
**Action Items**:
- [x] Implement User domain model
- [x] Implement UserManager interface and implementation
- [x] Implement user registration (anonymous and authenticated modes)
- [x] Implement username validation and uniqueness checks
- [x] Implement multiple connection detection
- [x] Implement lock/mute/activate operations
- [ ] Write user manager tests

#### 5.2 Table Controller
**Action Items**:
- [ ] Implement TableController
- [ ] Implement table state machine (WAITING → STARTING → DUELING → FINISHED)
- [ ] Implement player seat assignment and swapping
- [ ] Implement deck validation hooks
- [ ] Implement match/tournament creation logic
- [ ] Integrate with game/tournament controllers
- [ ] Write table controller tests

#### 5.3 Game Controller
**Action Items**:
- [ ] Implement GameController
- [ ] Implement game state management
- [ ] Implement player action queue and processing
- [ ] Implement GameView generation for clients
- [ ] Implement watcher management
- [ ] Implement replay recording
- [ ] Define game engine integration interface
- [ ] Write game controller tests (with mock game engine)

#### 5.4 Tournament System
**Action Items**:
- [ ] Implement TournamentController
- [ ] Implement tournament state machine
- [ ] Implement round management
- [ ] Implement Swiss pairing algorithm
- [ ] Implement elimination bracket generation
- [ ] Implement TournamentView generation
- [ ] Integrate with draft system
- [ ] Write tournament controller tests

#### 5.5 Draft System
**Action Items**:
- [ ] Implement DraftController
- [ ] Implement draft state machine (pick phase, pass direction)
- [ ] Implement booster pack generation from card repository
- [ ] Implement card pick handling
- [ ] Implement hidden card management (cards not yet seen)
- [ ] Implement DraftPickView generation
- [ ] Write draft controller tests

#### 5.6 Chat System
**Action Items**:
- [ ] Implement ChatManager
- [ ] Implement chat room management (game, tournament, table, lobby)
- [ ] Implement message broadcasting to all room members
- [ ] Implement message history (configurable limit)
- [ ] Implement user join/leave notifications
- [ ] Implement HTML sanitization (bluemonday)
- [ ] Write chat manager tests

#### 5.7 Room Management
**Action Items**:
- [ ] Implement GamesRoom (main lobby)
- [ ] Implement lobby features (user list, table list, finished matches)
- [ ] Implement room update broadcasting
- [ ] Write room manager tests

### Phase 6: Plugin System & Game Types (Week 19-21)

#### 6.1 Plugin Registry Pattern

**Action Items**:
- [ ] Implement plugin registry pattern
- [ ] Define GameType, TournamentType, PlayerType interfaces
- [ ] Implement game types:
  - [ ] TwoPlayerDuel
  - [ ] FreeForAll
  - [ ] CommanderFreeForAll
  - [ ] CommanderDuel
  - [ ] Brawl variants
  - [ ] Canadian Highlander
  - [ ] Momir variants
  - [ ] Oathbreaker variants
  - [ ] Penny Dreadful Commander
  - [ ] Tiny Leaders
- [ ] Implement tournament types:
  - [ ] Constructed
  - [ ] BoosterDraft
  - [ ] Sealed
- [ ] Implement player types:
  - [ ] Human
  - [ ] ComputerMAX (AI)
  - [ ] ComputerDraft (Draft AI)
- [ ] Register all types in init() functions
- [ ] Write plugin registry tests
- [ ] Port Java gameplay engine rules into Go (`MageEngine` parity)
  - [ ] Inventory existing Java engine modules (zones, layers, effects, abilities, watchers) and map them onto Go packages
  - [ ] Recreate the comprehensive turn/phase/step state machine, including priority passing, upkeep triggers, and cleanup handling
  - [ ] Port stack resolution, replacement/prevention effects, triggered abilities, continuous effect layers, and dependency management
  - [ ] Mirror event bus, watcher, and log infrastructure used by Java clients (game logs, game events, state change notifications)
  - [ ] Generate or translate card definitions/ability scripts so Go runtime can load the canonical card and token database
  - [ ] Build compatibility shims to ingest existing Java-serialized game states until all clients can consume Go-native ones
  - [ ] Establish regression harness comparing Java vs Go engine outputs (per-turn snapshots, rules tests, card-specific scenarios)

### Phase 7: Supporting Services (Week 22-23)

#### 7.1 Email Service
**Action Items**:
- [x] Implement MailClient interface
- [x] Implement SMTP client with gomail.v2
- [x] Implement Mailgun client with mailgun-go
- [ ] Create email templates (password reset, welcome, etc.)
- [ ] Add retry logic for failed emails
- [ ] Write email service tests (with mock SMTP server)

#### 7.2 Rating System
**Action Items**:
- [x] Implement Glicko rating calculation
- [x] Implement rating update on match completion
- [x] Write rating system tests with known inputs/outputs

#### 7.3 Card Repository & Caching
**Action Items**:
- [ ] Implement full-text search for cards (PostgreSQL trgm)
- [ ] Implement card caching with groupcache
- [ ] Add cache warming on startup
- [ ] Write card repository tests

#### 7.4 Logging & Metrics
**Action Items**:
- [ ] Configure Zap logger with structured fields
- [ ] Add contextual logging (session ID, user ID, game ID)
- [ ] Implement Prometheus metrics:
  - [ ] Active sessions gauge
  - [ ] Active games gauge
  - [ ] RPC request counter
  - [ ] RPC latency histogram
  - [ ] Database query latency
  - [ ] WebSocket connections gauge
- [ ] Create Grafana dashboard JSON
- [ ] Write metrics tests

### Phase 8: Testing & Quality (Week 24-26)

#### 8.1 Unit Tests
**Action Items**:
- [ ] Achieve 70%+ code coverage
- [ ] Test all repositories with test database
- [ ] Test all managers with mocks
- [ ] Test all controllers with mocks
- [ ] Test authentication flows
- [ ] Test session management
- [ ] Test rating calculations
- [ ] Set up test fixtures and helpers

#### 8.2 Integration Tests
**Action Items**:
- [ ] Test complete user registration → login → game flow
- [ ] Test session lifecycle (connect → ping → timeout)
- [ ] Test WebSocket callback delivery
- [ ] Test multi-player game creation and joining
- [ ] Test tournament creation and pairing
- [ ] Test draft flow (join → pick → complete)
- [ ] Test chat message broadcasting
- [ ] Test admin operations
- [ ] Set up integration test environment (Docker Compose)

#### 8.3 Performance & Load Testing
**Action Items**:
- [ ] Create load testing tool (simulated clients)
- [ ] Test 100 concurrent users
- [ ] Test 500 concurrent users
- [ ] Test 1000 concurrent users
- [ ] Profile with pprof (CPU, memory, goroutines)
- [ ] Optimize database queries (add indexes where needed)
- [ ] Optimize hot paths identified by profiling
- [ ] Benchmark critical functions

#### 8.4 Client Compatibility Testing
**Action Items**:
- [ ] Create Java client adapter for gRPC/WebSocket
- [ ] Test all 60+ RPC methods from Java client
- [ ] Test callback delivery to Java client
- [ ] Test complete game flow with Java client
- [ ] Document client integration guide
- [ ] Create example client code (Java and potentially others)

### Phase 9: Deployment & Documentation (Week 27-28)

#### 9.1 Containerization
**Action Items**:
- [ ] Create multi-stage Dockerfile
- [ ] Create Docker Compose for local development
- [ ] Add health check endpoint
- [ ] Implement graceful shutdown (drain connections)
- [ ] Test container deployment

#### 9.2 Monitoring & Observability
**Action Items**:
- [ ] Set up Prometheus scraping
- [ ] Create Grafana dashboards:
  - [ ] Server health (CPU, memory, goroutines)
  - [ ] Business metrics (active users, games, sessions)
  - [ ] RPC performance (latency, errors)
  - [ ] Database performance
- [ ] Set up alerting rules:
  - [ ] High error rate
  - [ ] High latency
  - [ ] Database connection failures
  - [ ] Memory leaks (goroutine growth)
- [ ] Add distributed tracing (optional: Jaeger)

#### 9.3 Documentation
**Action Items**:
- [ ] Write README.md for Go server
- [ ] Document architecture and design decisions
- [ ] Document configuration options
- [ ] Write deployment guide
- [ ] Write client integration guide (Java client → gRPC/WebSocket)
- [ ] Document API (generate from protobuf)
- [ ] Create troubleshooting guide
- [ ] Write developer onboarding guide

## Technology Stack Summary

### Core Technologies
- **Language**: Go 1.21+
- **RPC**: gRPC (google.golang.org/grpc)
- **Real-time**: WebSocket (github.com/gorilla/websocket)
- **Serialization**: Protocol Buffers (google.golang.org/protobuf)
- **Database**: PostgreSQL 15+ (github.com/jackc/pgx/v5)

### Libraries
- **Config**: Viper (github.com/spf13/viper)
- **Logging**: Zap (go.uber.org/zap)
- **Password**: Argon2id (golang.org/x/crypto/argon2)
- **Email**: gomail.v2 + mailgun-go/v4
- **Cache**: groupcache (github.com/golang/groupcache)
- **Validation**: validator/v10 (github.com/go-playground/validator)
- **Testing**: testify (github.com/stretchr/testify)
- **Metrics**: Prometheus (github.com/prometheus/client_golang)
- **Sanitization**: bluemonday (github.com/microcosm-cc/bluemonday)
- **Database migrations**: golang-migrate (github.com/golang-migrate/migrate)

### Development Tools
- **Protobuf**: buf (buf.build) or protoc
- **Linting**: golangci-lint
- **Formatting**: gofmt, goimports
- **Build**: Make
- **Container**: Docker, Docker Compose

## Timeline Estimate

### Detailed Timeline
- **Phase 1**: Foundation & Protobuf (3 weeks) - Weeks 1-3
- **Phase 2**: Core Infrastructure (3 weeks) - Weeks 4-6
- **Phase 3**: gRPC Server (4 weeks) - Weeks 7-10
- **Phase 4**: WebSocket Server (2 weeks) - Weeks 11-12
- **Phase 5**: Business Logic (6 weeks) - Weeks 13-18
- **Phase 6**: Plugin System (3 weeks) - Weeks 19-21
- **Phase 7**: Supporting Services (2 weeks) - Weeks 22-23
- **Phase 8**: Testing & Quality (3 weeks) - Weeks 24-26
- **Phase 9**: Deployment (2 weeks) - Weeks 27-28

**Total**: 28 weeks (~7 months) with 2-3 developers

### Parallel Work Opportunities
- **Weeks 1-3**: Dev 1 (Protobuf), Dev 2 (Database schemas), Dev 3 (Config)
- **Weeks 4-6**: Dev 1 (Repositories), Dev 2 (Session mgmt), Dev 3 (Auth)
- **Weeks 7-18**: Dev 1 (Game/Table), Dev 2 (Tournament/Draft), Dev 3 (RPC implementation)

## Success Metrics

- [ ] All 60+ RPC methods implemented and tested
- [ ] WebSocket push events working for all callback types
- [ ] Java client successfully connects and plays games
- [ ] Performance: Handle 1000+ concurrent users
- [ ] Latency: p95 < 100ms for game actions, p99 < 500ms
- [ ] Memory: < 2GB RAM for 500 concurrent games
- [ ] All game types functional
- [ ] 70%+ code coverage
- [ ] Load test: 1000 concurrent users for 1 hour without crashes

## Risk Mitigation

### High-Risk Areas
1. **Game Engine Integration**: Server is separate from game engine
2. **Callback Delivery**: Real-time push is critical for UX
3. **Session Management**: Lease mechanism must be rock-solid
4. **Client Compatibility**: Java client must work with minimal changes

### Mitigation Strategies
- **Feature Flags**: Gradual rollout of features
- **Extensive Testing**: Integration tests for all critical flows
- **Client Adapter**: Create compatibility layer in Java client to minimize code changes
- **Performance Testing**: Early and continuous load testing

## Next Steps for Development Agent

### Immediate Actions (Week 1)
1. Initialize Go project: `go mod init github.com/magefree/mage-server-go`
2. Set up project structure (create directories per structure above)
3. Create Makefile with targets: `build`, `test`, `proto`, `run`
4. Define all protobuf schemas in `api/proto/mage/v1/`
5. Generate Go code: `buf generate` or `make proto`
6. Create initial config.yaml with all settings
7. Implement config loader with Viper
8. Set up PostgreSQL schema migrations (create all .up.sql and .down.sql files)
9. Implement database connection with pgx
10. Create basic main.go that loads config and connects to database

### Development Priorities
1. **Critical Path**: Protobuf → Database → Session → Auth → Core RPC methods
2. **Early Wins**: Get ConnectUser, Ping, and basic session management working first
3. **Incremental**: Implement one RPC category at a time, test thoroughly
4. **Client Feedback Loop**: Create simple test client early to validate RPC responses

This plan is complete and ready for implementation. The coding agent can proceed with Phase 1 immediately.
