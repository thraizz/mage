package server

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/magefree/mage-server-go/internal/config"
	"github.com/magefree/mage-server-go/internal/session"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Add proper origin validation
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebSocketServer handles WebSocket connections for server-to-client callbacks
type WebSocketServer struct {
	sessionMgr session.Manager
	logger     *zap.Logger
	config     config.WebSocketConfig
}

// StartWebSocketServer starts the WebSocket server
func StartWebSocketServer(cfg config.WebSocketConfig, sessionMgr session.Manager, logger *zap.Logger) error {
	ws := &WebSocketServer{
		sessionMgr: sessionMgr,
		logger:     logger,
		config:     cfg,
	}

	http.HandleFunc("/ws", ws.handleConnection)

	logger.Info("starting WebSocket server", zap.String("address", cfg.Address))

	return http.ListenAndServe(cfg.Address, nil)
}

// handleConnection handles a WebSocket connection
func (ws *WebSocketServer) handleConnection(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from query params
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		http.Error(w, "missing sessionId", http.StatusBadRequest)
		return
	}

	// Validate session
	sess, ok := ws.sessionMgr.GetSession(sessionID)
	if !ok {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.logger.Error("failed to upgrade connection", zap.Error(err))
		return
	}
	defer conn.Close()

	// Create close channel for this connection
	done := make(chan struct{})

	// Set this as the active WebSocket and get the old one
	oldCloseChan := sess.SetWebSocketCloseChan(done)
	if oldCloseChan != nil {
		// Close the previous WebSocket connection
		ws.logger.Info("closing previous WebSocket connection",
			zap.String("session", sessionID),
			zap.String("user", sess.GetUserID()),
		)
		close(oldCloseChan)
	}

	// Clear the close channel when this connection ends (only if still active)
	defer sess.ClearWebSocketCloseChan(done)

	ws.logger.Info("WebSocket connected",
		zap.String("session", sessionID),
		zap.String("user", sess.GetUserID()),
	)

	// Start ping handler in background
	go ws.pingHandler(conn, done)

	// Start read handler to process pong messages from client
	go ws.readHandler(conn, done, sessionID)

	// Read callback channel and send to client
	for {
		select {
		case event, ok := <-sess.CallbackChan:
			if !ok {
				// Channel closed, connection terminated
				ws.logger.Info("callback channel closed", zap.String("session", sessionID))
				close(done)
				return
			}

			ws.logger.Info("sending event via WebSocket",
				zap.String("session", sessionID),
				zap.Any("event_type", event),
			)

			if err := ws.sendEvent(conn, event); err != nil {
				ws.logger.Error("failed to send event",
					zap.Error(err),
					zap.String("session", sessionID),
				)
				close(done)
				return
			}

			ws.logger.Info("event sent successfully via WebSocket",
				zap.String("session", sessionID),
			)

		case <-done:
			ws.logger.Info("WebSocket closed", zap.String("session", sessionID))
			return
		}
	}
}

// sendEvent sends an event to the client via WebSocket
func (ws *WebSocketServer) sendEvent(conn *websocket.Conn, event interface{}) error {
	// Convert event to JSON using protobuf JSON marshaler (handles Any types properly)
	protoMsg, ok := event.(proto.Message)
	if !ok {
		ws.logger.Error("event is not a proto.Message", zap.Any("event", event))
		return nil
	}

	marshaler := protojson.MarshalOptions{
		EmitUnpopulated: false,
		UseProtoNames:   false,
	}

	jsonData, err := marshaler.Marshal(protoMsg)
	if err != nil {
		ws.logger.Error("failed to marshal event to JSON", zap.Error(err))
		return err
	}

	ws.logger.Debug("marshaled WebSocket event",
		zap.Int("json_length", len(jsonData)),
	)

	// Set write deadline
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// Send message
	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		ws.logger.Error("failed to write WebSocket message", zap.Error(err))
		return err
	}

	return nil
}

// readHandler reads messages from the client (handles pong responses)
func (ws *WebSocketServer) readHandler(conn *websocket.Conn, done chan struct{}, sessionID string) {
	// Set pong handler
	conn.SetPongHandler(func(string) error {
		ws.logger.Debug("received pong from client", zap.String("session", sessionID))
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Set initial read deadline
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// Read loop - necessary to process control frames (pong)
	for {
		select {
		case <-done:
			return
		default:
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					ws.logger.Warn("websocket read error",
						zap.String("session", sessionID),
						zap.Error(err),
					)
				}
				close(done)
				return
			}
			// Client shouldn't send messages, but if they do, ignore them
		}
	}
}

// pingHandler sends periodic ping messages to keep connection alive
func (ws *WebSocketServer) pingHandler(conn *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(ws.config.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				ws.logger.Debug("ping failed, closing connection", zap.Error(err))
				close(done)
				return
			}

		case <-done:
			return
		}
	}
}
