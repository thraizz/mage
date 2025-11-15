package server

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	chatPrefixRoom       = "room"
	chatPrefixTable      = "table"
	chatPrefixGame       = "game"
	chatPrefixTournament = "tournament"
)

func formatChatRoomID(prefix, id string) string {
	return fmt.Sprintf("%s:%s", prefix, id)
}

// ChatJoin adds the current user to the requested chat channel.
func (s *mageServer) ChatJoin(ctx context.Context, req *pb.ChatJoinRequest) (*pb.ChatJoinResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.ChatJoinResponse{Success: false, Error: "session_id is required"}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.ChatJoinResponse{Success: false, Error: "session not found"}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.ChatJoinResponse{Success: false, Error: "session not associated with a user"}, nil
	}

	chatID := strings.TrimSpace(req.GetChatId())
	if chatID == "" {
		return &pb.ChatJoinResponse{Success: false, Error: "chat_id is required"}, nil
	}

	s.chatMgr.JoinRoom(chatID, username)

	s.logger.Debug("user joined chat",
		zap.String("chat_id", chatID),
		zap.String("username", username),
	)

	return &pb.ChatJoinResponse{Success: true}, nil
}

// ChatLeave removes the current user from the requested chat channel.
func (s *mageServer) ChatLeave(ctx context.Context, req *pb.ChatLeaveRequest) (*pb.ChatLeaveResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.ChatLeaveResponse{Success: false, Error: "session_id is required"}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.ChatLeaveResponse{Success: false, Error: "session not found"}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.ChatLeaveResponse{Success: false, Error: "session not associated with a user"}, nil
	}

	chatID := strings.TrimSpace(req.GetChatId())
	if chatID == "" {
		return &pb.ChatLeaveResponse{Success: false, Error: "chat_id is required"}, nil
	}

	s.chatMgr.LeaveRoom(chatID, username)

	s.logger.Debug("user left chat",
		zap.String("chat_id", chatID),
		zap.String("username", username),
	)

	return &pb.ChatLeaveResponse{Success: true}, nil
}

// ChatSendMessage sends a text message to a chat channel.
func (s *mageServer) ChatSendMessage(ctx context.Context, req *pb.ChatSendMessageRequest) (*pb.ChatSendMessageResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.ChatSendMessageResponse{Success: false, Error: "session_id is required"}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.ChatSendMessageResponse{Success: false, Error: "session not found"}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.ChatSendMessageResponse{Success: false, Error: "session not associated with a user"}, nil
	}

	chatID := strings.TrimSpace(req.GetChatId())
	if chatID == "" {
		return &pb.ChatSendMessageResponse{Success: false, Error: "chat_id is required"}, nil
	}

	message := strings.TrimSpace(req.GetMessage())
	if message == "" {
		return &pb.ChatSendMessageResponse{Success: false, Error: "message is required"}, nil
	}

	if err := s.chatMgr.SendMessage(chatID, username, message); err != nil {
		return &pb.ChatSendMessageResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.ChatSendMessageResponse{Success: true}, nil
}

// ChatFindByTable returns the chat identifier for the specified table.
func (s *mageServer) ChatFindByTable(ctx context.Context, req *pb.ChatFindByTableRequest) (*pb.ChatFindByTableResponse, error) {
	tableID := strings.TrimSpace(req.GetTableId())
	if tableID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "table_id is required")
	}

	if _, ok := s.tableMgr.GetTable(tableID); !ok {
		return nil, status.Errorf(codes.NotFound, "table not found")
	}

	return &pb.ChatFindByTableResponse{
		ChatId: formatChatRoomID(chatPrefixTable, tableID),
	}, nil
}

// ChatFindByGame returns the chat identifier for the specified game.
func (s *mageServer) ChatFindByGame(ctx context.Context, req *pb.ChatFindByGameRequest) (*pb.ChatFindByGameResponse, error) {
	gameID := strings.TrimSpace(req.GetGameId())
	if gameID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "game_id is required")
	}

	if _, ok := s.gameMgr.GetGame(gameID); !ok {
		return nil, status.Errorf(codes.NotFound, "game not found")
	}

	return &pb.ChatFindByGameResponse{
		ChatId: formatChatRoomID(chatPrefixGame, gameID),
	}, nil
}

// ChatFindByTournament returns the chat identifier for the specified tournament.
func (s *mageServer) ChatFindByTournament(ctx context.Context, req *pb.ChatFindByTournamentRequest) (*pb.ChatFindByTournamentResponse, error) {
	tournamentID := strings.TrimSpace(req.GetTournamentId())
	if tournamentID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "tournament_id is required")
	}

	if _, ok := s.tournamentMgr.GetTournament(tournamentID); !ok {
		return nil, status.Errorf(codes.NotFound, "tournament not found")
	}

	return &pb.ChatFindByTournamentResponse{
		ChatId: formatChatRoomID(chatPrefixTournament, tournamentID),
	}, nil
}

// ChatFindByRoom returns the chat identifier for the specified room.
func (s *mageServer) ChatFindByRoom(ctx context.Context, req *pb.ChatFindByRoomRequest) (*pb.ChatFindByRoomResponse, error) {
	roomID := strings.TrimSpace(req.GetRoomId())
	if roomID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "room_id is required")
	}

	if _, ok := s.roomMgr.GetRoom(roomID); !ok {
		return nil, status.Errorf(codes.NotFound, "room not found")
	}

	return &pb.ChatFindByRoomResponse{
		ChatId: formatChatRoomID(chatPrefixRoom, roomID),
	}, nil
}
