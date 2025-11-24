package server

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetMatchHistory retrieves a user's match history (paginated)
func (s *mageServer) GetMatchHistory(ctx context.Context, req *pb.GetMatchHistoryRequest) (*pb.GetMatchHistoryResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.GetMatchHistoryResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.GetMatchHistoryResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.GetMatchHistoryResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	// Get user ID from repository
	user, err := s.userRepo.GetByName(ctx, username)
	if err != nil {
		s.logger.Error("failed to get user for match history",
			zap.String("username", username),
			zap.Error(err),
		)
		return &pb.GetMatchHistoryResponse{
			Success: false,
			Error:   "user not found",
		}, nil
	}

	// Set default and max limits
	limit := req.GetLimit()
	if limit <= 0 {
		limit = 50 // Default
	}
	if limit > 100 {
		limit = 100 // Max
	}

	offset := req.GetOffset()
	if offset < 0 {
		offset = 0
	}

	// Get matches from repository
	matches, err := s.matchHistoryRepo.GetMatchesByUser(ctx, user.ID, int(limit), int(offset))
	if err != nil {
		s.logger.Error("failed to retrieve match history",
			zap.String("username", username),
			zap.Int64("user_id", user.ID),
			zap.Error(err),
		)
		return &pb.GetMatchHistoryResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to retrieve match history: %v", err),
		}, nil
	}

	// Get total count for pagination
	totalCount, err := s.matchHistoryRepo.CountMatchesByUser(ctx, user.ID)
	if err != nil {
		s.logger.Error("failed to count matches",
			zap.String("username", username),
			zap.Error(err),
		)
		totalCount = len(matches) // Fallback to partial count
	}

	// Convert to proto format
	matchEntries := make([]*pb.MatchHistoryEntry, len(matches))
	for i, match := range matches {
		players := make([]*pb.MatchPlayerInfo, len(match.Players))
		for j, player := range match.Players {
			players[j] = &pb.MatchPlayerInfo{
				UserId:   player.UserID,
				Username: player.Username,
				Deck:     player.Deck,
				Result:   player.Result,
			}
		}

		entry := &pb.MatchHistoryEntry{
			Id:              match.ID,
			GameId:          match.GameID,
			TableId:         match.TableID,
			TournamentId:    match.TournamentID,
			Players:         players,
			GameType:        match.GameType,
			StartTime:       timestamppb.New(match.StartTime),
			EndTime:         timestamppb.New(match.EndTime),
			DurationSeconds: int32(match.DurationSeconds),
			WinnerName:      match.WinnerName,
			CreatedAt:       timestamppb.New(match.CreatedAt),
		}

		// Only set winner_id if not null
		if match.WinnerID != nil {
			entry.WinnerId = *match.WinnerID
		}

		matchEntries[i] = entry
	}

	s.logger.Info("match history retrieved",
		zap.String("username", username),
		zap.Int("count", len(matchEntries)),
		zap.Int("total_count", totalCount),
		zap.Int32("limit", limit),
		zap.Int32("offset", offset),
	)

	return &pb.GetMatchHistoryResponse{
		Success:    true,
		Matches:    matchEntries,
		TotalCount: int32(totalCount),
	}, nil
}

// GetMatchById retrieves a specific match by ID
func (s *mageServer) GetMatchById(ctx context.Context, req *pb.GetMatchByIdRequest) (*pb.GetMatchByIdResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.GetMatchByIdResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.GetMatchByIdResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.GetMatchByIdResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	if req.GetMatchId() == 0 {
		return &pb.GetMatchByIdResponse{
			Success: false,
			Error:   "match_id is required",
		}, nil
	}

	// Get match from repository
	match, err := s.matchHistoryRepo.GetMatchByID(ctx, req.GetMatchId())
	if err != nil {
		s.logger.Error("failed to get match",
			zap.String("username", username),
			zap.Int64("match_id", req.GetMatchId()),
			zap.Error(err),
		)
		return &pb.GetMatchByIdResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to get match: %v", err),
		}, nil
	}

	// Convert players to proto format
	players := make([]*pb.MatchPlayerInfo, len(match.Players))
	for i, player := range match.Players {
		players[i] = &pb.MatchPlayerInfo{
			UserId:   player.UserID,
			Username: player.Username,
			Deck:     player.Deck,
			Result:   player.Result,
		}
	}

	entry := &pb.MatchHistoryEntry{
		Id:              match.ID,
		GameId:          match.GameID,
		TableId:         match.TableID,
		TournamentId:    match.TournamentID,
		Players:         players,
		GameType:        match.GameType,
		StartTime:       timestamppb.New(match.StartTime),
		EndTime:         timestamppb.New(match.EndTime),
		DurationSeconds: int32(match.DurationSeconds),
		WinnerName:      match.WinnerName,
		CreatedAt:       timestamppb.New(match.CreatedAt),
	}

	// Only set winner_id if not null
	if match.WinnerID != nil {
		entry.WinnerId = *match.WinnerID
	}

	s.logger.Info("match retrieved",
		zap.String("username", username),
		zap.Int64("match_id", req.GetMatchId()),
		zap.String("game_type", match.GameType),
	)

	return &pb.GetMatchByIdResponse{
		Success:    true,
		Match:      entry,
		ReplayData: match.ReplayData, // Include replay data if available
	}, nil
}
