package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/magefree/mage-server-go/internal/auth"
	"github.com/magefree/mage-server-go/internal/chat"
	"github.com/magefree/mage-server-go/internal/config"
	"github.com/magefree/mage-server-go/internal/draft"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/mail"
	_ "github.com/magefree/mage-server-go/internal/plugin" // Import to register game types
	"github.com/magefree/mage-server-go/internal/repository"
	"github.com/magefree/mage-server-go/internal/room"
	"github.com/magefree/mage-server-go/internal/server"
	"github.com/magefree/mage-server-go/internal/session"
	"github.com/magefree/mage-server-go/internal/table"
	"github.com/magefree/mage-server-go/internal/tournament"
	"github.com/magefree/mage-server-go/internal/user"
	pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	configPath = flag.String("config", "config/config.yaml", "path to configuration file")
	version    = "dev" // set via ldflags during build
)

func main() {
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := initLogger(cfg.Logging)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("starting MAGE server",
		zap.String("version", version),
		zap.String("config", *configPath),
	)

	if cfg.Auth.AdminPassword == "" {
		logger.Warn("admin password not configured; admin RPC access disabled")
	}

	// Create context that listens for termination signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Initialize database
	db, err := repository.NewDB(ctx, cfg.Database, logger)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Log database stats
	stats := db.Stats()
	logger.Info("database connection pool initialized",
		zap.Int32("total_conns", stats.TotalConns()),
		zap.Int32("idle_conns", stats.IdleConns()),
	)

	// Initialize session manager
	sessionMgr := session.NewManager(cfg.Server.LeasePeriod, logger)
	logger.Info("session manager initialized",
		zap.Duration("lease_period", cfg.Server.LeasePeriod),
	)

	// Start session cleanup goroutine
	go sessionMgr.CleanupExpiredSessions(ctx)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	statsRepo := repository.NewStatsRepository(db)
	deckRepo := repository.NewDeckRepository(db)
	cardRepo := repository.NewCardRepository(db, logger)
	matchHistoryRepo := repository.NewMatchHistoryRepository(db)

	// Initialize user manager
	userMgr := user.NewManager(userRepo, statsRepo, cfg.Validation, logger)
	logger.Info("user manager initialized")

	// Initialize auth token store
	tokenStore := auth.NewTokenStore(cfg.Auth.PasswordResetTokenTTL)
	logger.Info("auth token store initialized",
		zap.Duration("token_ttl", cfg.Auth.PasswordResetTokenTTL),
	)

	// Initialize room manager
	roomMgr := room.NewManager(logger)
	logger.Info("room manager initialized",
		zap.String("main_room_id", roomMgr.GetMainRoomID()),
	)

	// Initialize chat manager
	chatMgr := chat.NewManager(logger, sessionMgr)
	logger.Info("chat manager initialized")

	// Initialize table manager
	tableMgr := table.NewManager(logger)
	logger.Info("table manager initialized")

	// Initialize game manager
	gameMgr := game.NewManager(logger)
	logger.Info("game manager initialized")

	// Initialize game engine adapter
	mageEngine := game.NewMageEngine(logger)
	mageEngine.SetCardRepository(cardRepo) // Enable card metadata lookup
	gameAdapter := game.NewEngineAdapter(mageEngine, logger)

	// Initialize tournament manager
	tournamentMgr := tournament.NewManager(logger)
	logger.Info("tournament manager initialized")

	// Initialize draft manager
	draftMgr := draft.NewManager(logger)
	logger.Info("draft manager initialized")

	// Initialize email client
	mailClient, err := mail.NewClient(cfg.Mail, logger)
	if err != nil {
		logger.Warn("failed to initialize email client", zap.Error(err))
		mailClient = nil
	} else {
		logger.Info("email client initialized", zap.String("provider", cfg.Mail.Provider))
	}

	mageServer := server.NewMageServer(
		cfg,
		db,
		sessionMgr,
		userMgr,
		userRepo,
		statsRepo,
		deckRepo,
		cardRepo,
		matchHistoryRepo,
		roomMgr,
		chatMgr,
		tableMgr,
		gameMgr,
		tournamentMgr,
		draftMgr,
		tokenStore,
		mailClient,
		version,
		logger,
		gameAdapter,
	)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(server.ChainUnaryInterceptors(
			server.RecoveryInterceptor(logger),
			server.LoggingInterceptor(logger),
			server.SessionValidationInterceptor(sessionMgr),
			server.AdminInterceptor(sessionMgr),
			server.MetricsInterceptor(),
		)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    30 * time.Second,
			Timeout: 10 * time.Second,
		}),
		grpc.MaxConcurrentStreams(uint32(cfg.Server.GRPC.MaxConcurrentStreams)),
	)

	pb.RegisterMageServerServer(grpcServer, mageServer)

	// Create HTTP/JSON handler for browser-friendly JSON endpoints
	jsonHandler := server.NewHTTPJSONHandler(mageServer, logger)

	// Wrap gRPC server with gRPC-Web to support HTTP/JSON requests from browsers
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			// Allow all origins for development
			// TODO: Configure allowed origins for production
			return true
		}),
	)

	// Create HTTP handler that supports both JSON and gRPC-Web
	httpHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		contentType := req.Header.Get("Content-Type")

		// Route based on content type
		if contentType == "application/json" {
			// Use our custom JSON handler
			jsonHandler.ServeHTTP(resp, req)
		} else if wrappedGrpc.IsGrpcWebRequest(req) {
			// Use gRPC-Web for protobuf requests
			wrappedGrpc.ServeHTTP(resp, req)
		} else {
			// Use native gRPC for gRPC clients
			grpcServer.ServeHTTP(resp, req)
		}
	})

	// Add CORS support for browser requests
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // TODO: Configure for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(httpHandler)

	// Start HTTP server (supports JSON, gRPC-Web, and native gRPC)
	httpServer := &http.Server{
		Addr:    cfg.Server.GRPC.Address,
		Handler: corsHandler,
	}

	go func() {
		logger.Info("starting multi-protocol server (HTTP/JSON + gRPC-Web + native gRPC)",
			zap.String("address", cfg.Server.GRPC.Address))
		if serveErr := httpServer.ListenAndServe(); serveErr != nil && serveErr != http.ErrServerClosed {
			logger.Error("multi-protocol server error", zap.Error(serveErr))
		}
	}()

	// Start WebSocket server
	go func() {
		if wsErr := server.StartWebSocketServer(cfg.Server.WebSocket, sessionMgr, logger); wsErr != nil {
			logger.Error("WebSocket server error", zap.Error(wsErr))
		}
	}()

	// Start health check server if enabled
	if cfg.Health.Enabled {
		go func() {
			if healthErr := server.StartHealthCheckServer(
				cfg.Health.Address,
				db,
				sessionMgr,
				version,
				logger,
			); healthErr != nil {
				logger.Error("health check server error", zap.Error(healthErr))
			}
		}()
	}

	// TODO: Start metrics server if enabled
	// if cfg.Metrics.Enabled {
	//     go startMetricsServer(cfg.Metrics, logger)
	// }

	healthAddress := "disabled"
	if cfg.Health.Enabled {
		healthAddress = cfg.Health.Address
	}

	logger.Info("MAGE server initialized",
		zap.String("version", version),
		zap.String("grpc_address", cfg.Server.GRPC.Address),
		zap.String("websocket_address", cfg.Server.WebSocket.Address),
		zap.String("health_address", healthAddress),
		zap.Int("max_sessions", cfg.Server.MaxSessions),
	)

	// Wait for termination signal
	sig := <-sigChan
	logger.Info("received shutdown signal", zap.String("signal", sig.String()))

	// Graceful shutdown
	logger.Info("shutting down gracefully...")
	cancel()

	// Close all active sessions
	sessionMgr.CloseAll()

	// Shutdown HTTP server with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("error during HTTP server shutdown", zap.Error(err))
	}

	grpcServer.GracefulStop()

	logger.Info("MAGE server stopped")
}

// initLogger initializes the zap logger based on configuration
func initLogger(cfg config.LoggingConfig) (*zap.Logger, error) {
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	var zapCfg zap.Config
	if cfg.Format == "json" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapCfg.Level = zap.NewAtomicLevelAt(level)

	return zapCfg.Build()
}
