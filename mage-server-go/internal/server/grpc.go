package server

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/auth"
	"github.com/magefree/mage-server-go/internal/chat"
	"github.com/magefree/mage-server-go/internal/config"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/mail"
	"github.com/magefree/mage-server-go/internal/repository"
	"github.com/magefree/mage-server-go/internal/room"
	"github.com/magefree/mage-server-go/internal/session"
	"github.com/magefree/mage-server-go/internal/table"
	"github.com/magefree/mage-server-go/internal/user"
	pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// mageServer implements the MageServer gRPC service
type mageServer struct {
	pb.UnimplementedMageServerServer

	config        *config.Config
	logger        *zap.Logger
	serverVersion string

	sessionMgr       session.Manager
	userMgr          user.Manager
	userRepo         *repository.UserRepository
	statsRepo        *repository.StatsRepository
	deckRepo         *repository.DeckRepository
	cardRepo         *repository.CardRepository
	matchHistoryRepo *repository.MatchHistoryRepository
	roomMgr          *room.Manager
	chatMgr          *chat.Manager
	tableMgr         *table.Manager
	gameMgr          *game.Manager

	tokenStore *auth.TokenStore
	mailClient mail.Client
	db         *repository.DB

	savedDecks   map[string][]savedDeck
	savedDecksMu sync.RWMutex
	gameAdapter  *game.EngineAdapter
}

type savedDeck struct {
	Name string
	Deck table.DeckList
}

// NewMageServer creates a new MAGE server instance
func NewMageServer(
	cfg *config.Config,
	db *repository.DB,
	sessionMgr session.Manager,
	userMgr user.Manager,
	userRepo *repository.UserRepository,
	statsRepo *repository.StatsRepository,
	deckRepo *repository.DeckRepository,
	cardRepo *repository.CardRepository,
	matchHistoryRepo *repository.MatchHistoryRepository,
	roomMgr *room.Manager,
	chatMgr *chat.Manager,
	tableMgr *table.Manager,
	gameMgr *game.Manager,
	tokenStore *auth.TokenStore,
	mailClient mail.Client,
	serverVersion string,
	logger *zap.Logger,
	gameAdapter *game.EngineAdapter,
) *mageServer {
	s := &mageServer{
		config:           cfg,
		logger:           logger,
		serverVersion:    serverVersion,
		sessionMgr:       sessionMgr,
		userMgr:          userMgr,
		userRepo:         userRepo,
		statsRepo:        statsRepo,
		deckRepo:         deckRepo,
		cardRepo:         cardRepo,
		matchHistoryRepo: matchHistoryRepo,
		roomMgr:          roomMgr,
		chatMgr:          chatMgr,
		tableMgr:         tableMgr,
		gameMgr:          gameMgr,
		tokenStore:       tokenStore,
		mailClient:       mailClient,
		db:               db,
		savedDecks:       make(map[string][]savedDeck),
		gameAdapter:      gameAdapter,
	}

	// Set up game notifications to push updates via WebSocket
	s.SetupGameNotifications()

	return s
}

// SetupGameNotifications wires the game engine notifications to WebSocket push
func (s *mageServer) SetupGameNotifications() {
	if s.gameAdapter == nil {
		s.logger.Warn("game adapter not configured, notifications disabled")
		return
	}

	s.gameAdapter.SetNotificationCallback(func(notification game.GameNotification) {
		s.handleGameNotification(notification)
	})

	s.logger.Info("game notifications configured")
}

// handleGameNotification processes a game notification and pushes to relevant clients
func (s *mageServer) handleGameNotification(notification game.GameNotification) {
	gameID := notification.GameID
	if gameID == "" {
		s.logger.Warn("received game notification with empty game ID")
		return
	}

	gameInstance, ok := s.gameMgr.GetGame(gameID)
	if !ok {
		s.logger.Warn("game not found for notification", zap.String("game_id", gameID))
		return
	}

	s.logger.Info("handling game notification",
		zap.String("game_id", gameID),
		zap.String("type", notification.Type),
		zap.String("player_id", notification.PlayerID),
		zap.Int("num_players", len(gameInstance.Players)),
		zap.Any("data", notification.Data),
	)

	// Handle GAME_ERROR notifications - send only to the specific player
	if notification.Type == "GAME_ERROR" {
		if notification.PlayerID != "" {
			if errorMsg, exists := notification.Data["error"]; exists {
				s.sendGameErrorToPlayer(gameID, notification.PlayerID, fmt.Sprintf("%v", errorMsg))
			}
		}
		return
	}

	// Handle GAME_TARGET notifications - send only to the specific player
	if notification.Type == "GAME_TARGET" {
		if notification.PlayerID != "" {
			s.handleTargetNotification(gameID, notification.PlayerID, notification.Data)
		}
		return
	}

	// Handle GAME_XMANA notifications - send X value selection prompt to specific player
	if notification.Type == "GAME_XMANA" {
		if notification.PlayerID != "" {
			s.handleXManaNotification(gameID, notification.PlayerID, notification.Data)
		}
		return
	}

	// Handle GAME_CHOOSE_CHOICE notifications - send choice prompt to specific player
	// Used for combat declarations (declare attackers, declare blockers) and other choices
	if notification.Type == "GAME_CHOOSE_CHOICE" {
		if notification.PlayerID != "" {
			s.handleChoiceNotification(gameID, notification.PlayerID, notification.Data)
		}
		return
	}

	// GAME_UPDATE notifications come with a pre-built view for a specific player.
	// The engine's broadcast() already iterates over all players and sends one notification per player.
	// We should NOT loop through all players here - that would cause N² notifications.
	// More importantly, calling GetGameView() here would deadlock since broadcast() is called
	// while holding the engine's write lock, and GetGameView() needs a read lock.

	// Extract the pre-built view from the notification
	view, hasView := notification.Data["view"]
	if !hasView || view == nil {
		s.logger.Warn("GAME_UPDATE notification missing view data",
			zap.String("game_id", gameID),
			zap.String("player_id", notification.PlayerID),
		)
		return
	}

	// Send to the specific player this notification is for
	playerName := notification.PlayerID
	if playerName == "" {
		s.logger.Warn("GAME_UPDATE notification missing player_id",
			zap.String("game_id", gameID),
		)
		return
	}

	s.logger.Info("sending GAME_UPDATE to player",
		zap.String("game_id", gameID),
		zap.String("player", playerName),
	)
	s.sendGameUpdateWithView(gameID, playerName, view)

	// TODO: Support watchers. Currently the engine's broadcast() only sends to players.
	// Watchers are tracked in Game.Watchers (manager level), not in GameState.Players (engine level).
	// To properly support watchers, either:
	// 1. Have the engine also broadcast to watchers (would need access to watcher list), or
	// 2. Queue watcher updates to be sent after the engine lock is released
}

// sendGameUpdateWithView sends a GAME_UPDATE event to a specific player using a pre-built view.
// This avoids calling GetGameView() which would acquire a lock and cause deadlock when called
// from within broadcast() which already holds the write lock.
func (s *mageServer) sendGameUpdateWithView(gameID, playerName string, engineView interface{}) {
	// Convert engine view to protobuf
	gameView := s.engineViewToProto(engineView, playerName)
	if gameView == nil {
		s.logger.Warn("engineViewToProto returned nil",
			zap.String("game_id", gameID),
			zap.String("player", playerName),
		)
		return
	}

	// Create the GAME_UPDATE event
	updateData := &pb.GameUpdateData{
		Game: gameView,
	}

	event := &pb.ServerEvent{
		ObjectId: gameID,
		Method:   pb.CallbackMethod_GAME_UPDATE,
	}

	anyData, err := anypb.New(updateData)
	if err != nil {
		s.logger.Error("failed to marshal GameUpdateData",
			zap.String("game_id", gameID),
			zap.Error(err),
		)
		return
	}
	event.Data = anyData

	// Send to all sessions for this player
	s.logger.Info("looking up WebSocket sessions for player",
		zap.String("game_id", gameID),
		zap.String("player", playerName),
	)

	sessions := s.sessionMgr.GetSessionsByUser(playerName)
	if len(sessions) == 0 {
		s.logger.Warn("❌ NO WEBSOCKET SESSIONS FOUND - player will not receive update",
			zap.String("game_id", gameID),
			zap.String("player", playerName),
		)
		return
	}

	s.logger.Info("✓ found WebSocket sessions - sending GAME_UPDATE",
		zap.String("game_id", gameID),
		zap.String("player", playerName),
		zap.Int("session_count", len(sessions)),
		zap.Int32("turn", gameView.Turn),
		zap.String("phase", gameView.Phase),
		zap.String("priority_player", gameView.PriorityPlayerId),
	)

	for _, sess := range sessions {
		if !sess.SendCallback(event) {
			s.logger.Warn("failed to send GAME_UPDATE to session",
				zap.String("game_id", gameID),
				zap.String("player", playerName),
				zap.String("session_id", sess.ID),
			)
		} else {
			s.logger.Info("successfully sent GAME_UPDATE to session",
				zap.String("game_id", gameID),
				zap.String("player", playerName),
				zap.String("session_id", sess.ID),
			)
		}
	}
}

// sendGameUpdateToPlayer sends a GAME_UPDATE event to a specific player by fetching the view.
// NOTE: This function calls GetGameView() which acquires a read lock. Do NOT call this from
// within handleGameNotification when processing broadcast notifications, as that would deadlock.
// Use sendGameUpdateWithView instead for broadcast notifications.
func (s *mageServer) sendGameUpdateToPlayer(gameID, playerName string) {
	// Get the player's view of the game
	if s.gameAdapter == nil {
		s.logger.Warn("game adapter is nil, cannot send update",
			zap.String("game_id", gameID),
			zap.String("player", playerName),
		)
		return
	}

	engineView, err := s.gameAdapter.GetGameView(gameID, playerName)
	if err != nil {
		s.logger.Warn("failed to get game view for notification",
			zap.String("game_id", gameID),
			zap.String("player", playerName),
			zap.Error(err),
		)
		return
	}

	// Convert engine view to protobuf
	gameView := s.engineViewToProto(engineView, playerName)
	if gameView == nil {
		s.logger.Warn("engineViewToProto returned nil",
			zap.String("game_id", gameID),
			zap.String("player", playerName),
		)
		return
	}

	// Create the GAME_UPDATE event
	updateData := &pb.GameUpdateData{
		Game: gameView,
	}

	event := &pb.ServerEvent{
		ObjectId: gameID,
		Method:   pb.CallbackMethod_GAME_UPDATE,
	}

	anyData, err := anypb.New(updateData)
	if err != nil {
		s.logger.Error("failed to marshal GameUpdateData",
			zap.String("game_id", gameID),
			zap.Error(err),
		)
		return
	}
	event.Data = anyData

	// Send to all sessions for this player
	s.logger.Info("looking up WebSocket sessions for player",
		zap.String("game_id", gameID),
		zap.String("player", playerName),
	)

	sessions := s.sessionMgr.GetSessionsByUser(playerName)
	if len(sessions) == 0 {
		s.logger.Warn("❌ NO WEBSOCKET SESSIONS FOUND - player will not receive update",
			zap.String("game_id", gameID),
			zap.String("player", playerName),
		)
		return
	}

	s.logger.Info("✓ found WebSocket sessions - sending GAME_UPDATE",
		zap.String("game_id", gameID),
		zap.String("player", playerName),
		zap.Int("session_count", len(sessions)),
		zap.Int32("turn", gameView.Turn),
		zap.String("phase", gameView.Phase),
		zap.String("priority_player", gameView.PriorityPlayerId),
	)

	for _, sess := range sessions {
		if !sess.SendCallback(event) {
			s.logger.Warn("failed to send GAME_UPDATE to session",
				zap.String("game_id", gameID),
				zap.String("player", playerName),
				zap.String("session_id", sess.ID),
			)
		} else {
			s.logger.Info("successfully sent GAME_UPDATE to session",
				zap.String("game_id", gameID),
				zap.String("player", playerName),
				zap.String("session_id", sess.ID),
			)
		}
	}
}

// sendGameErrorToPlayer sends a GAME_ERROR event to a specific player
func (s *mageServer) sendGameErrorToPlayer(gameID, playerName, errorMsg string) {
	// Create the GAME_ERROR event
	errorData := &pb.GameErrorData{
		Error: errorMsg,
	}

	event := &pb.ServerEvent{
		ObjectId: gameID,
		Method:   pb.CallbackMethod_GAME_ERROR,
	}

	anyData, err := anypb.New(errorData)
	if err != nil {
		s.logger.Error("failed to marshal GameErrorData",
			zap.String("game_id", gameID),
			zap.Error(err),
		)
		return
	}
	event.Data = anyData

	// Send to all sessions for this player
	sessions := s.sessionMgr.GetSessionsByUser(playerName)
	if len(sessions) == 0 {
		s.logger.Warn("no sessions found for player to send error",
			zap.String("game_id", gameID),
			zap.String("player", playerName),
		)
		return
	}

	s.logger.Info("sending GAME_ERROR via WebSocket",
		zap.String("game_id", gameID),
		zap.String("player", playerName),
		zap.Int("session_count", len(sessions)),
		zap.String("error", errorMsg),
	)

	for _, sess := range sessions {
		if !sess.SendCallback(event) {
			s.logger.Warn("failed to send GAME_ERROR to session",
				zap.String("game_id", gameID),
				zap.String("player", playerName),
				zap.String("session_id", sess.ID),
			)
		} else {
			s.logger.Info("successfully sent GAME_ERROR to session",
				zap.String("game_id", gameID),
				zap.String("player", playerName),
				zap.String("session_id", sess.ID),
			)
		}
	}
}

// sendGameTargetToPlayer sends a GAME_TARGET event to a specific player for target selection
// This is triggered when a spell or ability requires the player to choose targets
func (s *mageServer) sendGameTargetToPlayer(gameID, playerName string, targetData *pb.GameTargetData) {
	// Create the GAME_TARGET event
	event := &pb.ServerEvent{
		ObjectId: gameID,
		Method:   pb.CallbackMethod_GAME_TARGET,
	}

	anyData, err := anypb.New(targetData)
	if err != nil {
		s.logger.Error("failed to marshal GameTargetData",
			zap.String("game_id", gameID),
			zap.Error(err),
		)
		return
	}
	event.Data = anyData

	// Send to all sessions for this player
	sessions := s.sessionMgr.GetSessionsByUser(playerName)
	if len(sessions) == 0 {
		s.logger.Warn("no sessions found for player to send target request",
			zap.String("game_id", gameID),
			zap.String("player", playerName),
		)
		return
	}

	s.logger.Info("sending GAME_TARGET via WebSocket",
		zap.String("game_id", gameID),
		zap.String("player", playerName),
		zap.Int("session_count", len(sessions)),
		zap.String("message", targetData.Message),
		zap.Int("target_count", len(targetData.Targets)),
		zap.Bool("required", targetData.Required),
	)

	for _, sess := range sessions {
		if !sess.SendCallback(event) {
			s.logger.Warn("failed to send GAME_TARGET to session",
				zap.String("game_id", gameID),
				zap.String("player", playerName),
				zap.String("session_id", sess.ID),
			)
		} else {
			s.logger.Info("successfully sent GAME_TARGET to session",
				zap.String("game_id", gameID),
				zap.String("player", playerName),
				zap.String("session_id", sess.ID),
			)
		}
	}
}

// handleTargetNotification processes a GAME_TARGET notification and sends to the player
func (s *mageServer) handleTargetNotification(gameID, playerName string, data map[string]interface{}) {
	// Extract target data from the notification
	message, _ := data["message"].(string)
	required, _ := data["required"].(bool)

	// Extract options
	options := make(map[string]string)
	if minTargets, ok := data["min_targets"].(int); ok {
		options["minTargets"] = strconv.Itoa(minTargets)
	}
	if maxTargets, ok := data["max_targets"].(int); ok {
		options["maxTargets"] = strconv.Itoa(maxTargets)
	}
	if sourceID, ok := data["source_id"].(string); ok {
		options["sourceId"] = sourceID
	}

	// Build CardView targets from the notification data
	var targets []*pb.CardView
	if targetsData, ok := data["targets"].([]map[string]interface{}); ok {
		for _, t := range targetsData {
			cardView := &pb.CardView{
				Id:   t["id"].(string),
				Name: t["name"].(string),
			}
			if cardType, ok := t["type"].(string); ok {
				cardView.Type = cardType
			}
			if zone, ok := t["zone"].(string); ok {
				// Convert zone string to int32 (zone codes are stored as strings in notifications)
				switch zone {
				case "LIBRARY":
					cardView.Zone = 0
				case "HAND":
					cardView.Zone = 1
				case "BATTLEFIELD":
					cardView.Zone = 2
				case "GRAVEYARD":
					cardView.Zone = 3
				case "STACK":
					cardView.Zone = 4
				case "EXILE":
					cardView.Zone = 5
				case "COMMAND":
					cardView.Zone = 6
				}
			}
			if controller, ok := t["controller"].(string); ok {
				cardView.ControllerId = controller
			}
			targets = append(targets, cardView)
		}
	}

	// Create the target data
	targetData := &pb.GameTargetData{
		Message:  message,
		Targets:  targets,
		Required: required,
		Options:  options,
	}

	// Send to the player
	s.sendGameTargetToPlayer(gameID, playerName, targetData)
}

// handleXManaNotification processes a GAME_XMANA notification and sends to the player
func (s *mageServer) handleXManaNotification(gameID, playerName string, data map[string]interface{}) {
	// Extract X mana data from the notification
	message, _ := data["message"].(string)
	available := 0
	if avail, ok := data["available"].(int); ok {
		available = avail
	}

	// Create the X mana data
	xManaData := &pb.GamePlayXManaData{
		Message:   message,
		Available: int32(available),
	}

	// Send to the player
	s.sendGameXManaToPlayer(gameID, playerName, xManaData)
}

// sendGameXManaToPlayer sends a GAME_PLAY_XMANA event to a specific player
func (s *mageServer) sendGameXManaToPlayer(gameID, playerName string, xManaData *pb.GamePlayXManaData) {
	// Create the GAME_PLAY_XMANA event
	event := &pb.ServerEvent{
		ObjectId: gameID,
		Method:   pb.CallbackMethod_GAME_PLAY_XMANA,
	}

	// Marshal the X mana data into the Any field
	anyData, err := anypb.New(xManaData)
	if err != nil {
		s.logger.Error("failed to marshal X mana data",
			zap.String("game_id", gameID),
			zap.String("player", playerName),
			zap.Error(err),
		)
		return
	}
	event.Data = anyData

	// Send to all sessions for this player
	sessions := s.sessionMgr.GetSessionsByUser(playerName)
	if len(sessions) == 0 {
		s.logger.Warn("no sessions found for player to send X value request",
			zap.String("game_id", gameID),
			zap.String("player", playerName),
		)
		return
	}

	s.logger.Info("sending GAME_PLAY_XMANA via WebSocket",
		zap.String("game_id", gameID),
		zap.String("player", playerName),
		zap.Int("session_count", len(sessions)),
		zap.String("message", xManaData.Message),
		zap.Int32("available", xManaData.Available),
	)

	for _, sess := range sessions {
		if !sess.SendCallback(event) {
			s.logger.Warn("failed to send GAME_PLAY_XMANA to session",
				zap.String("game_id", gameID),
				zap.String("player", playerName),
				zap.String("session_id", sess.ID),
			)
		} else {
			s.logger.Info("successfully sent GAME_PLAY_XMANA to session",
				zap.String("game_id", gameID),
				zap.String("player", playerName),
				zap.String("session_id", sess.ID),
			)
		}
	}
}

// handleChoiceNotification processes a GAME_CHOOSE_CHOICE notification and sends to the player
// This is used for combat declarations (attackers, blockers) and other player choices
func (s *mageServer) handleChoiceNotification(gameID, playerName string, data map[string]interface{}) {
	// Extract choice data from the notification
	message, _ := data["message"].(string)

	// Extract choices array
	var choices []string
	if choicesData, ok := data["choices"].([]string); ok {
		choices = choicesData
	} else if choicesInterface, ok := data["choices"].([]interface{}); ok {
		// Handle case where choices come through as []interface{}
		for _, c := range choicesInterface {
			if str, ok := c.(string); ok {
				choices = append(choices, str)
			}
		}
	}

	// Create the choice data
	choiceData := &pb.GameChoiceData{
		Message: message,
		Choices: choices,
	}

	// Send to the player
	s.sendGameChoiceToPlayer(gameID, playerName, choiceData)
}

// sendGameChoiceToPlayer sends a GAME_CHOOSE_CHOICE event to a specific player
func (s *mageServer) sendGameChoiceToPlayer(gameID, playerName string, choiceData *pb.GameChoiceData) {
	// Create the GAME_CHOOSE_CHOICE event
	event := &pb.ServerEvent{
		ObjectId: gameID,
		Method:   pb.CallbackMethod_GAME_CHOOSE_CHOICE,
	}

	// Marshal the choice data into the Any field
	anyData, err := anypb.New(choiceData)
	if err != nil {
		s.logger.Error("failed to marshal choice data",
			zap.String("game_id", gameID),
			zap.String("player", playerName),
			zap.Error(err),
		)
		return
	}
	event.Data = anyData

	// Send to all sessions for this player
	sessions := s.sessionMgr.GetSessionsByUser(playerName)
	if len(sessions) == 0 {
		s.logger.Warn("no sessions found for player to send choice prompt",
			zap.String("game_id", gameID),
			zap.String("player", playerName),
		)
		return
	}

	s.logger.Info("sending GAME_CHOOSE_CHOICE via WebSocket",
		zap.String("game_id", gameID),
		zap.String("player", playerName),
		zap.Int("session_count", len(sessions)),
		zap.String("message", choiceData.Message),
		zap.Int("choice_count", len(choiceData.Choices)),
	)

	for _, sess := range sessions {
		if !sess.SendCallback(event) {
			s.logger.Warn("failed to send GAME_CHOOSE_CHOICE to session",
				zap.String("game_id", gameID),
				zap.String("player", playerName),
				zap.String("session_id", sess.ID),
			)
		} else {
			s.logger.Info("successfully sent GAME_CHOOSE_CHOICE to session",
				zap.String("game_id", gameID),
				zap.String("player", playerName),
				zap.String("session_id", sess.ID),
			)
		}
	}
}

// engineViewToProto converts an engine view to protobuf GameView
// Phase 8: Simplified to only handle PlaytestGameView (MageEngine removed)
func (s *mageServer) engineViewToProto(engineView interface{}, playerID string) *pb.GameView {
	if engineView == nil {
		return nil
	}

	// Only handle PlaytestGameView now (MageEngine support removed)
	if playtestData, ok := engineView.(*game.PlaytestGameView); ok {
		return s.playtestViewToProto(playtestData, playerID)
	}

	// Unknown view type
	s.logger.Warn("unknown engine view type",
		zap.String("type", fmt.Sprintf("%T", engineView)),
	)
	return nil
}

// playtestViewToProto converts a PlaytestGameView to protobuf GameView
func (s *mageServer) playtestViewToProto(data *game.PlaytestGameView, playerID string) *pb.GameView {
	view := &pb.GameView{
		GameId:            data.GameID,
		State:             "IN_PROGRESS", // Playtest is always in progress
		Phase:             "",            // No phase tracking in playtest
		Step:              "",            // No step tracking in playtest
		Turn:              int32(data.Turn),
		ActivePlayerId:    data.ActivePlayerID,
		ActiveControlSeat: playerID, // Viewing player's ID (from Task 1.2)
		Battlefield:       playtestEngineCardsToProto(data.Battlefield),
		Stack:             playtestEngineCardsToProto(data.Stack),
		Exile:             playtestEngineCardsToProto(data.Exile),
		Command:           playtestEngineCardsToProto(data.Command),
	}

	// Convert players
	players := make([]*pb.PlayerView, 0)

	// Add the viewing player first
	if data.Me != nil {
		players = append(players, &pb.PlayerView{
			PlayerId:      data.Me.PlayerID,
			Name:          data.Me.Name,
			Life:          int32(data.Me.Life),
			Poison:        int32(data.Me.Poison),
			Energy:        int32(data.Me.Energy),
			LibraryCount:  int32(data.Me.LibraryCount),
			HandCount:     int32(data.Me.HandCount),
			Hand:          playtestEngineCardsToProto(data.Me.Hand),
			Library:       playtestEngineCardsToProto(data.Me.Library), // Task 1.6: Include library for viewing player
			Graveyard:     playtestEngineCardsToProto(data.Me.Graveyard),
			KeptHand:      data.Me.KeptHand,
			MulliganCount: int32(data.Me.MulliganCount),
		})
	}

	// Add opponents
	for _, opponent := range data.Opponents {
		players = append(players, &pb.PlayerView{
			PlayerId:      opponent.PlayerID,
			Name:          opponent.Name,
			Life:          int32(opponent.Life),
			Poison:        int32(opponent.Poison),
			Energy:        int32(opponent.Energy),
			LibraryCount:  int32(opponent.LibraryCount),
			HandCount:     int32(opponent.HandCount),
			Hand:          playtestEngineCardsToProto(opponent.Hand), // Empty for opponents
			Graveyard:     playtestEngineCardsToProto(opponent.Graveyard),
			KeptHand:      opponent.KeptHand,
			MulliganCount: int32(opponent.MulliganCount),
		})
	}

	view.Players = players
	return view
}

// playtestEngineCardsToProto converts playtest engine cards to protobuf
func playtestEngineCardsToProto(cards []*game.Card) []*pb.CardView {
	if cards == nil {
		return []*pb.CardView{}
	}

	result := make([]*pb.CardView, 0, len(cards))
	for _, card := range cards {
		if card == nil {
			continue
		}

		cardView := &pb.CardView{
			Id:           card.ID,
			Name:         card.Name,
			DisplayName:  card.DisplayName,
			ManaCost:     card.ManaCost,
			Type:         card.Type,
			SubTypes:     card.SubTypes,
			SuperTypes:   card.SuperTypes,
			Color:        card.Color,
			Power:        card.Power,
			Toughness:    card.Toughness,
			Loyalty:      card.Loyalty,
			OwnerId:      card.OwnerID,
			ControllerId: card.ControllerID,
			Tapped:       card.Tapped,
			Flipped:      card.Flipped,
			Transformed:  card.Transformed,
			FaceDown:     card.FaceDown,
			RulesText:    card.RulesText,
		}

		// Convert counters
		counters := make([]*pb.CounterView, 0, len(card.Counters))
		for _, counter := range card.Counters {
			counters = append(counters, &pb.CounterView{
				Name:  counter.Name,
				Count: int32(counter.Count),
			})
		}
		cardView.Counters = counters

		result = append(result, cardView)
	}

	return result
}

// ==================== Authentication & Connection Methods ====================

// ConnectUser handles user connection
func (s *mageServer) ConnectUser(ctx context.Context, req *pb.ConnectUserRequest) (*pb.ConnectUserResponse, error) {
	if req.GetUserName() == "" || req.GetPassword() == "" {
		return &pb.ConnectUserResponse{
			Success: false,
			Error:   "username and password are required",
		}, nil
	}

	host := extractHostFromContext(ctx)

	u, err := s.userMgr.Authenticate(ctx, req.GetUserName(), req.GetPassword())
	if err != nil {
		s.logger.Warn("connect user failed",
			zap.String("username", req.GetUserName()),
			zap.Error(err),
		)
		return &pb.ConnectUserResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	var sess *session.Session
	if restoreID := req.GetRestoreSessionId(); restoreID != "" {
		if existing, ok := s.sessionMgr.GetSession(restoreID); ok {
			sess = existing
		}
	}

	if sess == nil {
		sessionID := req.GetSessionId()
		if sessionID == "" {
			sessionID = uuid.NewString()
		} else {
			// Ensure no stale session exists with the same ID
			s.sessionMgr.RemoveSession(sessionID)
		}
		sess = s.sessionMgr.CreateSession(sessionID, host)
	}

	sess.SetUserID(u.Name)
	sess.SetAdmin(false)
	sess.UpdateActivity()

	s.userMgr.UserConnect(ctx, u.Name, sess.ID)

	mainRoomID := s.roomMgr.GetMainRoomID()
	_ = s.roomMgr.UserJoinRoom(mainRoomID, u.Name)
	s.chatMgr.JoinRoom(mainRoomID, u.Name)

	s.logger.Info("user connected",
		zap.String("username", u.Name),
		zap.String("session_id", sess.ID),
		zap.String("host", host),
	)

	userIDStr := strconv.FormatInt(u.ID, 10)

	return &pb.ConnectUserResponse{
		Success:   true,
		SessionId: sess.ID,
		UserId:    userIDStr,
	}, nil
}

// ConnectAdmin handles admin connection
func (s *mageServer) ConnectAdmin(ctx context.Context, req *pb.ConnectAdminRequest) (*pb.ConnectAdminResponse, error) {
	if s.config.Auth.AdminPassword == "" {
		return &pb.ConnectAdminResponse{
			Success: false,
			Error:   "admin access not configured",
		}, nil
	}

	if req.GetPassword() != s.config.Auth.AdminPassword {
		s.logger.Warn("admin authentication failed", zap.String("session_id", req.GetSessionId()))
		return &pb.ConnectAdminResponse{
			Success: false,
			Error:   "invalid admin password",
		}, nil
	}

	host := extractHostFromContext(ctx)

	sessionID := req.GetSessionId()
	if sessionID == "" {
		sessionID = uuid.NewString()
	} else {
		s.sessionMgr.RemoveSession(sessionID)
	}

	sess := s.sessionMgr.CreateSession(sessionID, host)
	sess.SetAdmin(true)
	sess.SetUserID("admin")
	sess.UpdateActivity()

	s.logger.Info("admin connected",
		zap.String("session_id", sess.ID),
		zap.String("host", host),
	)

	return &pb.ConnectAdminResponse{
		Success: true,
	}, nil
}

// Ping keeps session alive
func (s *mageServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	if req.GetSessionId() == "" {
		return &pb.PingResponse{Success: false}, nil
	}

	s.sessionMgr.UpdateActivity(req.GetSessionId())
	return &pb.PingResponse{Success: true}, nil
}

// AuthRegister registers a new user
func (s *mageServer) AuthRegister(ctx context.Context, req *pb.AuthRegisterRequest) (*pb.AuthRegisterResponse, error) {
	if req.GetUserName() == "" || req.GetPassword() == "" {
		return &pb.AuthRegisterResponse{
			Success: false,
			Error:   "username and password are required",
		}, nil
	}

	if s.config.Auth.RequireEmail && req.GetEmail() == "" {
		return &pb.AuthRegisterResponse{
			Success: false,
			Error:   "email is required",
		}, nil
	}

	if err := s.userMgr.Register(ctx, req.GetUserName(), req.GetPassword(), req.GetEmail()); err != nil {
		s.logger.Warn("user registration failed",
			zap.String("username", req.GetUserName()),
			zap.Error(err),
		)
		return &pb.AuthRegisterResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	if s.mailClient != nil && req.GetEmail() != "" {
		if err := s.mailClient.SendWelcomeEmail(req.GetEmail(), req.GetUserName()); err != nil {
			s.logger.Warn("failed to send welcome email",
				zap.String("username", req.GetUserName()),
				zap.String("email", req.GetEmail()),
				zap.Error(err),
			)
		}
	}

	s.logger.Info("user registered", zap.String("username", req.GetUserName()))

	return &pb.AuthRegisterResponse{Success: true}, nil
}

// ConnectSetUserData updates a user's preferences
func (s *mageServer) ConnectSetUserData(ctx context.Context, req *pb.ConnectSetUserDataRequest) (*pb.ConnectSetUserDataResponse, error) {
	sess, ok := s.sessionMgr.GetSession(req.GetSessionId())
	if !ok {
		return &pb.ConnectSetUserDataResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	prefs := session.Preferences{
		AvatarID:                 req.GetAvatarId(),
		ShowAbsoluteAbilities:    req.GetShowAbsoluteAbilities(),
		AllowRequestsFromFriends: req.GetAllowRequestsFromFriends(),
		ConfirmEmptyManaPool:     req.GetConfirmEmptyManaPool(),
		UserGroup:                req.GetUserGroup(),
		SkipPrioritySteps:        req.GetUserSkipPrioritySteps(),
		FlagsName:                req.GetFlagsName(),
		AskMoveToGraveOrder:      req.GetAskMoveToGraveOrder(),
	}

	sess.SetPreferences(prefs)

	s.logger.Debug("updated session preferences",
		zap.String("session_id", sess.ID),
		zap.String("user", sess.GetUserID()),
	)

	return &pb.ConnectSetUserDataResponse{
		Success: true,
	}, nil
}

// AuthSendTokenToEmail sends password reset token to email
func (s *mageServer) AuthSendTokenToEmail(ctx context.Context, req *pb.AuthSendTokenToEmailRequest) (*pb.AuthSendTokenToEmailResponse, error) {
	if req.GetEmail() == "" {
		return &pb.AuthSendTokenToEmailResponse{
			Success: false,
			Error:   "email is required",
		}, nil
	}

	u, err := s.userRepo.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		// Do not leak existence of email addresses; log and return success
		s.logger.Warn("password reset requested for unknown email",
			zap.String("email", req.GetEmail()),
			zap.Error(err),
		)
		return &pb.AuthSendTokenToEmailResponse{Success: true}, nil
	}

	token, err := s.tokenStore.GenerateToken(req.GetEmail())
	if err != nil {
		s.logger.Error("failed to generate password reset token",
			zap.String("email", req.GetEmail()),
			zap.Error(err),
		)
		return &pb.AuthSendTokenToEmailResponse{
			Success: false,
			Error:   "failed to generate token",
		}, nil
	}

	if s.mailClient != nil {
		if err := s.mailClient.SendPasswordResetEmail(req.GetEmail(), u.Name, token); err != nil {
			s.logger.Error("failed to send password reset email",
				zap.String("email", req.GetEmail()),
				zap.Error(err),
			)
			return &pb.AuthSendTokenToEmailResponse{
				Success: false,
				Error:   "failed to send email",
			}, nil
		}
	} else {
		s.logger.Info("password reset email not sent (mail disabled)",
			zap.String("email", req.GetEmail()),
			zap.String("username", u.Name),
		)
	}

	s.logger.Info("password reset token generated",
		zap.String("email", req.GetEmail()),
		zap.String("username", u.Name),
	)

	return &pb.AuthSendTokenToEmailResponse{Success: true}, nil
}

// AuthResetPassword resets user password with token
func (s *mageServer) AuthResetPassword(ctx context.Context, req *pb.AuthResetPasswordRequest) (*pb.AuthResetPasswordResponse, error) {
	if req.GetEmail() == "" || req.GetToken() == "" || req.GetNewPassword() == "" {
		return &pb.AuthResetPasswordResponse{
			Success: false,
			Error:   "email, token, and new password are required",
		}, nil
	}

	if !s.tokenStore.ConsumeToken(req.GetEmail(), req.GetToken()) {
		return &pb.AuthResetPasswordResponse{
			Success: false,
			Error:   "invalid or expired token",
		}, nil
	}

	u, err := s.userRepo.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		s.logger.Warn("password reset failed: user not found",
			zap.String("email", req.GetEmail()),
			zap.Error(err),
		)
		return &pb.AuthResetPasswordResponse{
			Success: false,
			Error:   "user not found",
		}, nil
	}

	passwordHash, err := auth.HashPassword(req.GetNewPassword())
	if err != nil {
		s.logger.Error("failed to hash new password",
			zap.String("username", u.Name),
			zap.Error(err),
		)
		return &pb.AuthResetPasswordResponse{
			Success: false,
			Error:   "failed to hash password",
		}, nil
	}

	if err := s.userRepo.UpdatePassword(ctx, u.Name, passwordHash); err != nil {
		s.logger.Error("failed to update user password",
			zap.String("username", u.Name),
			zap.Error(err),
		)
		return &pb.AuthResetPasswordResponse{
			Success: false,
			Error:   "failed to update password",
		}, nil
	}

	s.logger.Info("password reset successful",
		zap.String("username", u.Name),
		zap.String("email", req.GetEmail()),
	)

	return &pb.AuthResetPasswordResponse{Success: true}, nil
}

// ==================== Server Info Methods ====================

// GetServerState returns server state information
func (s *mageServer) GetServerState(ctx context.Context, req *pb.GetServerStateRequest) (*pb.GetServerStateResponse, error) {
	serverState := &pb.ServerState{
		ActivePlayers:     int32(s.sessionMgr.GetActiveSessions()),
		ActiveGames:       int32(s.gameMgr.GetActiveGameCount()),
		ActiveTournaments: 0, // Tournament feature removed
		ActiveTables:      int32(s.tableMgr.GetActiveTableCount()),
		NumberOfThreads:   int32(runtime.NumGoroutine()),
		ServerVersion:     s.serverVersion,
		ServerTime:        timestamppb.Now(),
	}

	return &pb.GetServerStateResponse{
		ServerState: serverState,
	}, nil
}

// ServerGetPromotionMessages returns promotion messages (if any)
func (s *mageServer) ServerGetPromotionMessages(ctx context.Context, req *pb.ServerGetPromotionMessagesRequest) (*pb.ServerGetPromotionMessagesResponse, error) {
	// Promotion messages will eventually be driven from configuration or storage.
	return &pb.ServerGetPromotionMessagesResponse{
		Messages: []string{},
	}, nil
}

// ServerAddFeedbackMessage logs feedback from clients
func (s *mageServer) ServerAddFeedbackMessage(ctx context.Context, req *pb.ServerAddFeedbackMessageRequest) (*pb.ServerAddFeedbackMessageResponse, error) {
	s.logger.Info("feedback received",
		zap.String("session_id", req.GetSessionId()),
		zap.String("user_name", req.GetUserName()),
		zap.String("title", req.GetTitle()),
		zap.String("type", req.GetFeedbackType()),
		zap.String("email", req.GetEmail()),
	)

	return &pb.ServerAddFeedbackMessageResponse{
		Success: true,
	}, nil
}

// ServerGetMainRoomId returns the main room ID
func (s *mageServer) ServerGetMainRoomId(ctx context.Context, req *pb.ServerGetMainRoomIdRequest) (*pb.ServerGetMainRoomIdResponse, error) {
	return &pb.ServerGetMainRoomIdResponse{
		RoomId: s.roomMgr.GetMainRoomID(),
	}, nil
}

// RoomGetUsers returns users in a room
func (s *mageServer) RoomGetUsers(ctx context.Context, req *pb.RoomGetUsersRequest) (*pb.RoomGetUsersResponse, error) {
	roomID := req.GetRoomId()
	if roomID == "" {
		roomID = s.roomMgr.GetMainRoomID()
	}

	usernames := s.roomMgr.GetRoomUsers(roomID)
	users := make([]*pb.UserView, 0, len(usernames))

	for _, username := range usernames {
		users = append(users, &pb.UserView{
			UserName: username,
			State:    "ONLINE",
		})
	}

	return &pb.RoomGetUsersResponse{
		Users: users,
	}, nil
}

// RoomGetFinishedMatches returns finished matches for a room (currently empty placeholder)
func (s *mageServer) RoomGetFinishedMatches(ctx context.Context, req *pb.RoomGetFinishedMatchesRequest) (*pb.RoomGetFinishedMatchesResponse, error) {
	return &pb.RoomGetFinishedMatchesResponse{
		FinishedMatches: []*pb.MatchView{},
	}, nil
}

// RoomGetAllTables returns all tables in a room
func (s *mageServer) RoomGetAllTables(ctx context.Context, req *pb.RoomGetAllTablesRequest) (*pb.RoomGetAllTablesResponse, error) {
	roomID := req.GetRoomId()
	if roomID == "" {
		roomID = s.roomMgr.GetMainRoomID()
	}

	tables := s.tableMgr.GetTablesByRoom(roomID)
	tableViews := make([]*pb.TableView, 0, len(tables))
	for _, tbl := range tables {
		tableViews = append(tableViews, s.tableToProto(tbl))
	}

	return &pb.RoomGetAllTablesResponse{
		Tables: tableViews,
	}, nil
}

// RoomGetTableById returns a table by ID
func (s *mageServer) RoomGetTableById(ctx context.Context, req *pb.RoomGetTableByIdRequest) (*pb.RoomGetTableByIdResponse, error) {
	tableID := req.GetTableId()
	if tableID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "table_id is required")
	}

	tbl, ok := s.tableMgr.GetTable(tableID)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "table not found")
	}

	if req.GetRoomId() != "" && tbl.RoomID != req.GetRoomId() {
		return nil, status.Errorf(codes.NotFound, "table not found in room")
	}

	return &pb.RoomGetTableByIdResponse{
		Table: s.tableToProto(tbl),
	}, nil
}

// ==================== Table Management Methods ====================

// RoomCreateTable creates a new game table
func (s *mageServer) RoomCreateTable(ctx context.Context, req *pb.RoomCreateTableRequest) (*pb.RoomCreateTableResponse, error) {
	sess, ok := s.sessionMgr.GetSession(req.GetSessionId())
	if !ok {
		return &pb.RoomCreateTableResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	controller := sess.GetUserID()
	if controller == "" {
		return &pb.RoomCreateTableResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	roomID := req.GetRoomId()
	if roomID == "" {
		roomID = s.roomMgr.GetMainRoomID()
	}

	if _, exists := s.roomMgr.GetRoom(roomID); !exists {
		s.roomMgr.CreateRoom(roomID, fmt.Sprintf("Room %s", roomID))
	}

	matchOptions := req.GetMatchOptions()

	tableName := fmt.Sprintf("%s's table", controller)
	if matchOptions != nil && strings.TrimSpace(matchOptions.GetName()) != "" {
		tableName = matchOptions.GetName()
	}

	gameType := "Duel"
	if matchOptions != nil && strings.TrimSpace(matchOptions.GetGameType()) != "" {
		gameType = matchOptions.GetGameType()
	}

	numSeats := deriveSeatCount(gameType)
	password := ""

	newTable := s.tableMgr.CreateTable(tableName, gameType, controller, roomID, numSeats, password)

	if err := newTable.AddPlayer(controller, "Human"); err != nil {
		s.logger.Debug("failed to add controller to table",
			zap.String("table_id", newTable.ID),
			zap.String("controller", controller),
			zap.Error(err),
		)
	}

	if err := s.roomMgr.UserJoinRoom(roomID, controller); err != nil {
		s.logger.Debug("failed to ensure controller present in room",
			zap.String("room_id", roomID),
			zap.String("controller", controller),
			zap.Error(err),
		)
	}

	s.logger.Info("table created",
		zap.String("table_id", newTable.ID),
		zap.String("room_id", roomID),
		zap.String("controller", controller),
		zap.String("game_type", gameType),
	)

	return &pb.RoomCreateTableResponse{
		Success: true,
		TableId: newTable.ID,
	}, nil
}

// NOTE: TableIsOwner and ChatFindByTournament implementations moved to
// grpc_table.go and grpc_chat.go respectively

// ==================== Helper Functions ====================

func (s *mageServer) tableToProto(t *table.Table) *pb.TableView {
	seats := make([]*pb.SeatView, len(t.Seats))
	for _, seat := range t.Seats {
		if seat.Number >= 0 && seat.Number < len(seats) {
			seats[seat.Number] = &pb.SeatView{
				SeatNumber: int32(seat.Number),
				PlayerName: seat.PlayerName,
				PlayerType: seat.PlayerType,
				Locked:     seat.Locked,
			}
		}
	}

	matchOptions := &pb.MatchOptions{
		Name:     t.Name,
		GameType: t.GameType,
	}

	return &pb.TableView{
		TableId:              t.ID,
		GameType:             t.GameType,
		TableName:            t.Name,
		ControllerName:       t.ControllerName,
		TableStateText:       t.GetState().String(),
		NumSeats:             int32(t.NumSeats),
		Seats:                seats,
		MatchOptions:         matchOptions,
		CreateTime:           timestamppb.New(t.CreateTime),
		IsTournament:         t.Tournament,
		TournamentId:         t.TournamentID,
		SpecTatorshipAllowed: true,
		Password:             "",
	}
}

func deriveSeatCount(gameType string) int {
	if gameType == "" {
		return 2
	}

	lower := strings.ToLower(gameType)
	switch {
	case strings.Contains(lower, "commander"),
		strings.Contains(lower, "freeforall"),
		strings.Contains(lower, "free-for-all"),
		strings.Contains(lower, "brawl"),
		strings.Contains(lower, "oathbreaker"):
		return 4
	default:
		return 2
	}
}

// Helper function to extract host from context
func extractHostFromContext(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok && p.Addr != net.Addr(nil) {
		if host, _, err := net.SplitHostPort(p.Addr.String()); err == nil {
			return host
		}
		return p.Addr.String()
	}
	return "unknown"
}
