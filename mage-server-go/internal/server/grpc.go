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
	"github.com/magefree/mage-server-go/internal/draft"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/mail"
	"github.com/magefree/mage-server-go/internal/repository"
	"github.com/magefree/mage-server-go/internal/room"
	"github.com/magefree/mage-server-go/internal/session"
	"github.com/magefree/mage-server-go/internal/table"
	"github.com/magefree/mage-server-go/internal/tournament"
	"github.com/magefree/mage-server-go/internal/user"
	pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// mageServer implements the MageServer gRPC service
type mageServer struct {
	pb.UnimplementedMageServerServer

	config        *config.Config
	logger        *zap.Logger
	serverVersion string

	sessionMgr    session.Manager
	userMgr       user.Manager
	userRepo      *repository.UserRepository
	roomMgr       *room.Manager
	chatMgr       *chat.Manager
	tableMgr      *table.Manager
	gameMgr       *game.Manager
	tournamentMgr *tournament.Manager
	draftMgr      *draft.Manager

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
	roomMgr *room.Manager,
	chatMgr *chat.Manager,
	tableMgr *table.Manager,
	gameMgr *game.Manager,
	tournamentMgr *tournament.Manager,
	draftMgr *draft.Manager,
	tokenStore *auth.TokenStore,
	mailClient mail.Client,
	serverVersion string,
	logger *zap.Logger,
	gameAdapter *game.EngineAdapter,
) *mageServer {
	return &mageServer{
		config:        cfg,
		logger:        logger,
		serverVersion: serverVersion,
		sessionMgr:    sessionMgr,
		userMgr:       userMgr,
		userRepo:      userRepo,
		roomMgr:       roomMgr,
		chatMgr:       chatMgr,
		tableMgr:      tableMgr,
		gameMgr:       gameMgr,
		tournamentMgr: tournamentMgr,
		draftMgr:      draftMgr,
		tokenStore:    tokenStore,
		mailClient:    mailClient,
		db:            db,
		savedDecks:    make(map[string][]savedDeck),
		gameAdapter:   gameAdapter,
	}
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
		ActiveTournaments: int32(s.tournamentMgr.GetActiveTournamentCount()),
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
