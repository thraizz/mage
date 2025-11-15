package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	chatMgr := chat.NewManager(logger)
	logger.Info("chat manager initialized")

	// Initialize table manager
	tableMgr := table.NewManager(logger)
	logger.Info("table manager initialized")

	// Initialize game manager
	gameMgr := game.NewManager(logger)
	logger.Info("game manager initialized")

	// Initialize game engine adapter
	mageEngine := game.NewMageEngine(logger)
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

	lis, err := net.Listen("tcp", cfg.Server.GRPC.Address)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	// Start gRPC server
	go func() {
		logger.Info("starting gRPC server", zap.String("address", cfg.Server.GRPC.Address))
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			logger.Error("gRPC server error", zap.Error(serveErr))
		}
	}()

	// Start WebSocket server
	go func() {
		if wsErr := server.StartWebSocketServer(cfg.Server.WebSocket, sessionMgr, logger); wsErr != nil {
			logger.Error("WebSocket server error", zap.Error(wsErr))
		}
	}()

	// TODO: Start metrics server if enabled
	// if cfg.Metrics.Enabled {
	//     go startMetricsServer(cfg.Metrics, logger)
	// }

	logger.Info("MAGE server initialized",
		zap.String("version", version),
		zap.String("grpc_address", cfg.Server.GRPC.Address),
		zap.String("websocket_address", cfg.Server.WebSocket.Address),
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
