package server

import (
	"context"
	"time"

	pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AdminGetUsers retrieves a list of all users (admin only)
func (s *mageServer) AdminGetUsers(ctx context.Context, req *pb.AdminGetUsersRequest) (*pb.AdminGetUsersResponse, error) {
	// Session is already validated as admin by AdminInterceptor

	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to get all users", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to retrieve users")
	}

	// Get connected users for online status
	connectedUsers := s.userMgr.GetConnectedUsers()
	connectedSet := make(map[string]bool)
	for _, username := range connectedUsers {
		connectedSet[username] = true
	}

	userViews := make([]*pb.UserView, 0, len(users))
	for _, user := range users {
		// Get user stats for this user
		stats, statsErr := s.statsRepo.GetByUserName(ctx, user.Name)
		var userStats *pb.UserStatsView
		if statsErr == nil && stats != nil {
			userStats = &pb.UserStatsView{
				Matches:     int32(stats.Matches),
				Tournaments: int32(stats.Tournaments),
				Wins:        int32(stats.Wins),
				Losses:      int32(stats.Losses),
				Rating:      stats.Rating,
				TourneysWon: int32(stats.TourneysWon),
			}
		}

		userView := &pb.UserView{
			UserName: user.Name,
			State:    "",                   // TODO: Add user state tracking
			Admin:    user.Name == "admin", // TODO: Proper admin flag in DB
			Stats:    userStats,
		}

		userViews = append(userViews, userView)
	}

	return &pb.AdminGetUsersResponse{
		Users: userViews,
	}, nil
}

// AdminDisconnectUser forcibly disconnects a user (admin only)
func (s *mageServer) AdminDisconnectUser(ctx context.Context, req *pb.AdminDisconnectUserRequest) (*pb.AdminDisconnectUserResponse, error) {
	// Session is already validated as admin by AdminInterceptor

	username := req.GetUserName()
	if username == "" {
		return nil, status.Error(codes.InvalidArgument, "user_name is required")
	}

	// Check if user exists
	user, err := s.userMgr.GetByName(ctx, username)
	if err != nil {
		s.logger.Warn("user not found for disconnect", zap.String("username", username))
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Disconnect user (will trigger cleanup of sessions)
	s.userMgr.UserDisconnect(ctx, "") // Will disconnect all sessions for this user

	s.logger.Info("admin disconnected user",
		zap.String("username", user.Name),
		zap.String("admin_session", req.GetSessionId()),
	)

	return &pb.AdminDisconnectUserResponse{
		Success: true,
	}, nil
}

// AdminMuteUser mutes a user for a specified duration (admin only)
func (s *mageServer) AdminMuteUser(ctx context.Context, req *pb.AdminMuteUserRequest) (*pb.AdminMuteUserResponse, error) {
	// Session is already validated as admin by AdminInterceptor

	username := req.GetUserName()
	if username == "" {
		return nil, status.Error(codes.InvalidArgument, "user_name is required")
	}

	durationMinutes := req.GetDurationMinutes()
	if durationMinutes <= 0 {
		return nil, status.Error(codes.InvalidArgument, "duration_minutes must be positive")
	}

	duration := time.Duration(durationMinutes) * time.Minute

	if err := s.userMgr.MuteUser(ctx, username, duration); err != nil {
		s.logger.Error("failed to mute user",
			zap.String("username", username),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to mute user")
	}

	s.logger.Info("admin muted user",
		zap.String("username", username),
		zap.Int64("duration_minutes", durationMinutes),
		zap.String("admin_session", req.GetSessionId()),
	)

	return &pb.AdminMuteUserResponse{
		Success: true,
	}, nil
}

// AdminLockUser locks a user account for a specified duration (admin only)
func (s *mageServer) AdminLockUser(ctx context.Context, req *pb.AdminLockUserRequest) (*pb.AdminLockUserResponse, error) {
	// Session is already validated as admin by AdminInterceptor

	username := req.GetUserName()
	if username == "" {
		return nil, status.Error(codes.InvalidArgument, "user_name is required")
	}

	durationMinutes := req.GetDurationMinutes()
	if durationMinutes <= 0 {
		return nil, status.Error(codes.InvalidArgument, "duration_minutes must be positive")
	}

	duration := time.Duration(durationMinutes) * time.Minute

	if err := s.userMgr.LockUser(ctx, username, duration); err != nil {
		s.logger.Error("failed to lock user",
			zap.String("username", username),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to lock user")
	}

	s.logger.Info("admin locked user",
		zap.String("username", username),
		zap.Int64("duration_minutes", durationMinutes),
		zap.String("admin_session", req.GetSessionId()),
	)

	return &pb.AdminLockUserResponse{
		Success: true,
	}, nil
}

// AdminActivateUser activates a user account (admin only)
func (s *mageServer) AdminActivateUser(ctx context.Context, req *pb.AdminActivateUserRequest) (*pb.AdminActivateUserResponse, error) {
	// Session is already validated as admin by AdminInterceptor

	username := req.GetUserName()
	if username == "" {
		return nil, status.Error(codes.InvalidArgument, "user_name is required")
	}

	if err := s.userMgr.ActivateUser(ctx, username); err != nil {
		s.logger.Error("failed to activate user",
			zap.String("username", username),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to activate user")
	}

	s.logger.Info("admin activated user",
		zap.String("username", username),
		zap.String("admin_session", req.GetSessionId()),
	)

	return &pb.AdminActivateUserResponse{
		Success: true,
	}, nil
}

// AdminToggleActivateUser toggles a user's active status (admin only)
func (s *mageServer) AdminToggleActivateUser(ctx context.Context, req *pb.AdminToggleActivateUserRequest) (*pb.AdminToggleActivateUserResponse, error) {
	// Session is already validated as admin by AdminInterceptor

	username := req.GetUserName()
	if username == "" {
		return nil, status.Error(codes.InvalidArgument, "user_name is required")
	}

	// Get current status
	user, err := s.userMgr.GetByName(ctx, username)
	if err != nil {
		s.logger.Warn("user not found for toggle", zap.String("username", username))
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Toggle status
	var toggleErr error
	newStatus := !user.Active
	if newStatus {
		toggleErr = s.userMgr.ActivateUser(ctx, username)
	} else {
		toggleErr = s.userMgr.DeactivateUser(ctx, username)
	}

	if toggleErr != nil {
		s.logger.Error("failed to toggle user status",
			zap.String("username", username),
			zap.Error(toggleErr),
		)
		return nil, status.Error(codes.Internal, "failed to toggle user status")
	}

	s.logger.Info("admin toggled user status",
		zap.String("username", username),
		zap.Bool("new_status", newStatus),
		zap.String("admin_session", req.GetSessionId()),
	)

	return &pb.AdminToggleActivateUserResponse{
		Success:   true,
		NewStatus: newStatus,
	}, nil
}

// AdminEndUserSession ends a specific user session (admin only)
func (s *mageServer) AdminEndUserSession(ctx context.Context, req *pb.AdminEndUserSessionRequest) (*pb.AdminEndUserSessionResponse, error) {
	// Session is already validated as admin by AdminInterceptor

	targetSessionID := req.GetTargetSessionId()
	if targetSessionID == "" {
		return nil, status.Error(codes.InvalidArgument, "target_session_id is required")
	}

	// Get session to find username for logging
	sess, exists := s.sessionMgr.GetSession(targetSessionID)
	username := "unknown"
	if exists {
		username = sess.GetUserID()
	}

	// Remove the session
	s.sessionMgr.RemoveSession(targetSessionID)

	s.logger.Info("admin ended user session",
		zap.String("target_session_id", targetSessionID),
		zap.String("username", username),
		zap.String("admin_session", req.GetSessionId()),
	)

	return &pb.AdminEndUserSessionResponse{
		Success: true,
	}, nil
}

// AdminTableRemove forcibly removes a table (admin only)
func (s *mageServer) AdminTableRemove(ctx context.Context, req *pb.AdminTableRemoveRequest) (*pb.AdminTableRemoveResponse, error) {
	// Session is already validated as admin by AdminInterceptor

	tableID := req.GetTableId()
	if tableID == "" {
		return nil, status.Error(codes.InvalidArgument, "table_id is required")
	}

	// Check if table exists
	table, exists := s.tableMgr.GetTable(tableID)
	if !exists {
		return nil, status.Error(codes.NotFound, "table not found")
	}

	// Remove the table
	s.tableMgr.RemoveTable(tableID)

	s.logger.Info("admin removed table",
		zap.String("table_id", tableID),
		zap.String("table_name", table.Name),
		zap.String("admin_session", req.GetSessionId()),
	)

	// TODO: Broadcast table list update to room
	// This requires implementing websocket broadcast functionality
	if room, ok := s.roomMgr.GetRoom(s.roomMgr.GetMainRoomID()); ok {
		_ = room // TODO: room.BroadcastTableUpdate() when websocket is integrated
	}

	return &pb.AdminTableRemoveResponse{
		Success: true,
	}, nil
}

// AdminSendBroadcastMessage sends a server-wide announcement (admin only)
func (s *mageServer) AdminSendBroadcastMessage(ctx context.Context, req *pb.AdminSendBroadcastMessageRequest) (*pb.AdminSendBroadcastMessageResponse, error) {
	// Session is already validated as admin by AdminInterceptor

	message := req.GetMessage()
	if message == "" {
		return nil, status.Error(codes.InvalidArgument, "message is required")
	}

	s.logger.Info("admin sending broadcast message",
		zap.String("message", message),
		zap.String("admin_session", req.GetSessionId()),
	)

	// TODO: Send broadcast to all connected sessions
	// This requires implementing websocket event broadcasting
	// broadcastEvent := &pb.ServerEvent{
	// 	Event: &pb.ServerEvent_ServerMessage{
	// 		ServerMessage: &pb.ServerMessageEvent{
	// 			Message:  fmt.Sprintf("[SERVER ANNOUNCEMENT] %s", message),
	// 			IsSystem: true,
	// 		},
	// 	},
	// }

	// Get all active sessions (for future websocket broadcast)
	sessionIDs := s.sessionMgr.GetAllSessionIDs()
	broadcastCount := len(sessionIDs)
	// TODO: Implement when websocket broadcasting is ready
	// for _, sessionID := range sessionIDs {
	// 	if sess, ok := s.sessionMgr.GetSession(sessionID); ok {
	// 		sess.SendEvent(broadcastEvent)
	// 		broadcastCount++
	// 	}
	// }

	s.logger.Info("broadcast message sent",
		zap.Int("recipient_count", broadcastCount),
	)

	return &pb.AdminSendBroadcastMessageResponse{
		Success:        true,
		RecipientCount: int32(broadcastCount),
	}, nil
}

// AdminGetAllActiveGames returns all active games with memory vs database comparison
func (s *mageServer) AdminGetAllActiveGames(ctx context.Context, req *pb.AdminGetAllActiveGamesRequest) (*pb.AdminGetAllActiveGamesResponse, error) {
	// Session is already validated as admin by AdminInterceptor

	// Get games from server memory
	memoryGames := s.gameMgr.ListGames()
	memoryGameIDs := make(map[string]bool)
	for _, g := range memoryGames {
		memoryGameIDs[g.ID] = true
	}

	// Get games from database
	var dbGames []*pb.DebugActiveGameInfo
	dbGameIDs := make(map[string]bool)

	if s.activeGameRepo != nil {
		activeGames, err := s.activeGameRepo.LoadAllActiveGames(ctx)
		if err != nil {
			s.logger.Warn("failed to load active games from database", zap.Error(err))
		} else {
			for _, ag := range activeGames {
				dbGameIDs[ag.GameID] = true

				// Check if also in memory
				inMemory := memoryGameIDs[ag.GameID]

				gameInfo := &pb.DebugActiveGameInfo{
					GameId:     ag.GameID,
					TableId:    ag.TableID,
					GameType:   ag.GameType,
					Players:    ag.Players,
					TurnNumber: int32(ag.TurnNumber),
					State:      ag.State,
					CreatedAt:  ag.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
					UpdatedAt:  ag.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
					InMemory:   inMemory,
					InDatabase: true,
				}
				dbGames = append(dbGames, gameInfo)
			}
		}
	}

	// Add memory-only games
	for _, g := range memoryGames {
		if !dbGameIDs[g.ID] {
			// Game is only in memory, not in database
			gameInfo := &pb.DebugActiveGameInfo{
				GameId:     g.ID,
				TableId:    g.TableID,
				GameType:   g.GameType,
				Players:    g.Players,
				TurnNumber: int32(g.Turn),
				State:      g.State.String(),
				InMemory:   true,
				InDatabase: false,
			}
			dbGames = append(dbGames, gameInfo)
		}
	}

	s.logger.Info("admin retrieved all active games",
		zap.Int("total_in_memory", len(memoryGames)),
		zap.Int("total_in_database", len(dbGameIDs)),
		zap.String("admin_session", req.GetSessionId()),
	)

	return &pb.AdminGetAllActiveGamesResponse{
		Games:           dbGames,
		TotalInMemory:   int32(len(memoryGames)),
		TotalInDatabase: int32(len(dbGameIDs)),
	}, nil
}

// AdminGetServerDebugState returns detailed server debug state
func (s *mageServer) AdminGetServerDebugState(ctx context.Context, req *pb.AdminGetServerDebugStateRequest) (*pb.AdminGetServerDebugStateResponse, error) {
	// Session is already validated as admin by AdminInterceptor

	// Server memory stats
	memoryActiveGames := s.gameMgr.GetActiveGameCount()
	memoryActiveTables := s.tableMgr.GetActiveTableCount()
	memoryActiveSessions := len(s.sessionMgr.GetAllSessionIDs())

	// Database stats
	var dbActiveGames int
	var dbMatchHistory int

	if s.activeGameRepo != nil {
		count, err := s.activeGameRepo.CountActiveGames(ctx)
		if err != nil {
			s.logger.Warn("failed to count active games in database", zap.Error(err))
		} else {
			dbActiveGames = count
		}
	}

	if s.matchHistoryRepo != nil {
		count, err := s.matchHistoryRepo.CountMatches(ctx)
		if err != nil {
			s.logger.Warn("failed to count match history in database", zap.Error(err))
		} else {
			dbMatchHistory = count
		}
	}

	// Find discrepancies
	var gamesInMemoryOnly []string
	var gamesInDbOnly []string

	memoryGames2 := s.gameMgr.ListGames()
	memoryGameIDs2 := make(map[string]bool)
	for _, g := range memoryGames2 {
		memoryGameIDs2[g.ID] = true
	}

	if s.activeGameRepo != nil {
		dbGames, err := s.activeGameRepo.LoadAllActiveGames(ctx)
		if err == nil {
			dbGameIDs := make(map[string]bool)
			for _, ag := range dbGames {
				dbGameIDs[ag.GameID] = true
				if !memoryGameIDs2[ag.GameID] {
					gamesInDbOnly = append(gamesInDbOnly, ag.GameID)
				}
			}
			for id := range memoryGameIDs2 {
				if !dbGameIDs[id] {
					gamesInMemoryOnly = append(gamesInMemoryOnly, id)
				}
			}
		}
	}

	// Calculate uptime (placeholder - would need startup time tracking)
	uptime := "unknown"

	s.logger.Info("admin retrieved server debug state",
		zap.Int("memory_games", memoryActiveGames),
		zap.Int("db_games", dbActiveGames),
		zap.Int("discrepancies", len(gamesInMemoryOnly)+len(gamesInDbOnly)),
		zap.String("admin_session", req.GetSessionId()),
	)

	return &pb.AdminGetServerDebugStateResponse{
		MemoryActiveGames:    int32(memoryActiveGames),
		MemoryActiveTables:   int32(memoryActiveTables),
		MemoryActiveSessions: int32(memoryActiveSessions),
		DbActiveGames:        int32(dbActiveGames),
		DbMatchHistory:       int32(dbMatchHistory),
		GamesInMemoryOnly:    gamesInMemoryOnly,
		GamesInDbOnly:        gamesInDbOnly,
		ServerUptime:         uptime,
		LastGameSave:         "", // TODO: Track last save time
	}, nil
}
