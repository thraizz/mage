package server

import (
	"context"
	"fmt"
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
		stats, statsErr := s.statsRepo.GetByUsername(ctx, user.Name)
		var userStats *pb.UserStatsView
		if statsErr == nil && stats != nil {
			userStats = &pb.UserStatsView{
				Matches:     int32(stats.Matches),
				Tournaments: int32(stats.Tournaments),
				Wins:        int32(stats.Wins),
				Losses:      int32(stats.Losses),
				Rating:      stats.Rating,
				TourneysWon: int32(stats.TournamentsWon),
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

	// Broadcast table list update to room
	if room, ok := s.roomMgr.GetRoom(s.roomMgr.GetMainRoomID()); ok {
		room.BroadcastTableUpdate()
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

	// Send broadcast to all connected sessions
	broadcastEvent := &pb.ServerEvent{
		Event: &pb.ServerEvent_ServerMessage{
			ServerMessage: &pb.ServerMessageEvent{
				Message:  fmt.Sprintf("[SERVER ANNOUNCEMENT] %s", message),
				IsSystem: true,
			},
		},
	}

	// Get all active sessions and send the broadcast
	sessionIDs := s.sessionMgr.GetAllSessionIDs()
	broadcastCount := 0
	for _, sessionID := range sessionIDs {
		if sess, ok := s.sessionMgr.GetSession(sessionID); ok {
			sess.SendEvent(broadcastEvent)
			broadcastCount++
		}
	}

	s.logger.Info("broadcast message sent",
		zap.Int("recipient_count", broadcastCount),
	)

	return &pb.AdminSendBroadcastMessageResponse{
		Success:        true,
		RecipientCount: int32(broadcastCount),
	}, nil
}
