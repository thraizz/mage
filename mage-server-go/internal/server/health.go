package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/magefree/mage-server-go/internal/repository"
	"github.com/magefree/mage-server-go/internal/session"
	"go.uber.org/zap"
)

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status     string            `json:"status"`
	Version    string            `json:"version"`
	Timestamp  time.Time         `json:"timestamp"`
	Uptime     float64           `json:"uptime_seconds"`
	Database   DatabaseHealth    `json:"database"`
	Sessions   SessionHealth     `json:"sessions"`
	Components map[string]string `json:"components"`
}

// DatabaseHealth represents database health metrics
type DatabaseHealth struct {
	Status       string `json:"status"`
	TotalConns   int32  `json:"total_connections"`
	IdleConns    int32  `json:"idle_connections"`
	AcquireCount int64  `json:"acquire_count"`
}

// SessionHealth represents session health metrics
type SessionHealth struct {
	ActiveSessions int `json:"active_sessions"`
}

// HealthCheckHandler provides health check endpoints
type HealthCheckHandler struct {
	db         *repository.DB
	sessionMgr session.Manager
	version    string
	logger     *zap.Logger
	startTime  time.Time
}

// NewHealthCheckHandler creates a new health check handler
func NewHealthCheckHandler(
	db *repository.DB,
	sessionMgr session.Manager,
	version string,
	logger *zap.Logger,
) *HealthCheckHandler {
	return &HealthCheckHandler{
		db:         db,
		sessionMgr: sessionMgr,
		version:    version,
		logger:     logger,
		startTime:  time.Now(),
	}
}

// ServeHTTP implements http.Handler
func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := h.buildHealthCheckResponse()

	w.Header().Set("Content-Type", "application/json")

	// Determine HTTP status based on health
	statusCode := http.StatusOK
	if response.Status == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode health check response", zap.Error(err))
	}
}

// buildHealthCheckResponse builds the health check response
func (h *HealthCheckHandler) buildHealthCheckResponse() HealthCheckResponse {
	// Check database health
	dbStats := h.db.Stats()
	dbHealth := DatabaseHealth{
		Status:       "healthy",
		TotalConns:   dbStats.TotalConns(),
		IdleConns:    dbStats.IdleConns(),
		AcquireCount: dbStats.AcquireCount(),
	}

	// If no connections, database is unhealthy
	if dbStats.TotalConns() == 0 {
		dbHealth.Status = "unhealthy"
	}

	// Check session manager health
	sessionHealth := SessionHealth{
		ActiveSessions: h.sessionMgr.ActiveSessionCount(),
	}

	// Determine overall status
	overallStatus := "healthy"
	components := map[string]string{
		"database": dbHealth.Status,
		"sessions": "healthy",
	}

	if dbHealth.Status == "unhealthy" {
		overallStatus = "degraded"
	}

	return HealthCheckResponse{
		Status:     overallStatus,
		Version:    h.version,
		Timestamp:  time.Now(),
		Uptime:     time.Since(h.startTime).Seconds(),
		Database:   dbHealth,
		Sessions:   sessionHealth,
		Components: components,
	}
}

// ReadinessHandler provides a simple readiness check
type ReadinessHandler struct {
	db     *repository.DB
	logger *zap.Logger
}

// NewReadinessHandler creates a new readiness handler
func NewReadinessHandler(db *repository.DB, logger *zap.Logger) *ReadinessHandler {
	return &ReadinessHandler{
		db:     db,
		logger: logger,
	}
}

// ServeHTTP implements http.Handler for readiness checks
func (h *ReadinessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if database is reachable
	stats := h.db.Stats()
	if stats.TotalConns() == 0 {
		http.Error(w, "Database not ready", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// LivenessHandler provides a simple liveness check
type LivenessHandler struct{}

// NewLivenessHandler creates a new liveness handler
func NewLivenessHandler() *LivenessHandler {
	return &LivenessHandler{}
}

// ServeHTTP implements http.Handler for liveness checks
func (h *LivenessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// StartHealthCheckServer starts the health check HTTP server
func StartHealthCheckServer(
	address string,
	db *repository.DB,
	sessionMgr session.Manager,
	version string,
	logger *zap.Logger,
) error {
	mux := http.NewServeMux()

	// Health check endpoint (detailed)
	healthHandler := NewHealthCheckHandler(db, sessionMgr, version, logger)
	mux.Handle("/health", healthHandler)

	// Readiness probe (Kubernetes)
	readinessHandler := NewReadinessHandler(db, logger)
	mux.Handle("/ready", readinessHandler)

	// Liveness probe (Kubernetes)
	livenessHandler := NewLivenessHandler()
	mux.Handle("/live", livenessHandler)

	logger.Info("starting health check server", zap.String("address", address))

	server := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server.ListenAndServe()
}
