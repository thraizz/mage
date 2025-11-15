package server

import (
	"context"
	"strings"

	"github.com/magefree/mage-server-go/internal/tournament"
	pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TournamentJoin adds a player to an existing tournament.
func (s *mageServer) TournamentJoin(ctx context.Context, req *pb.TournamentJoinRequest) (*pb.TournamentJoinResponse, error) {
	player, tourney, err := s.resolveTournamentPlayer(req.GetSessionId(), req.GetTournamentId())
	if err != nil {
		return &pb.TournamentJoinResponse{Success: false, Error: err.Error()}, nil
	}

	if err := tourney.AddPlayer(player); err != nil {
		if !strings.Contains(err.Error(), "already joined") {
			return &pb.TournamentJoinResponse{Success: false, Error: err.Error()}, nil
		}
	}

	if err := s.roomMgr.UserJoinRoom(tourney.RoomID, player); err != nil {
		s.logger.Debug("failed to ensure player joined tournament room",
			zap.String("tournament_id", tourney.ID),
			zap.String("room_id", tourney.RoomID),
			zap.String("player", player),
			zap.Error(err),
		)
	}

	s.logger.Info("player joined tournament",
		zap.String("tournament_id", tourney.ID),
		zap.String("player", player),
	)

	return &pb.TournamentJoinResponse{Success: true}, nil
}

// TournamentStart transitions a tournament into progress.
func (s *mageServer) TournamentStart(ctx context.Context, req *pb.TournamentStartRequest) (*pb.TournamentStartResponse, error) {
	player, tourney, err := s.resolveTournamentPlayer(req.GetSessionId(), req.GetTournamentId())
	if err != nil {
		return &pb.TournamentStartResponse{Success: false, Error: err.Error()}, nil
	}

	if req.GetRoomId() != "" && req.GetRoomId() != tourney.RoomID {
		return &pb.TournamentStartResponse{Success: false, Error: "tournament not found in requested room"}, nil
	}

	session, _ := s.sessionMgr.GetSession(strings.TrimSpace(req.GetSessionId()))
	if session == nil {
		return &pb.TournamentStartResponse{Success: false, Error: "session not found"}, nil
	}

	if !session.IsAdminSession() && !tourney.IsController(player) {
		return &pb.TournamentStartResponse{Success: false, Error: "only controller or admin can start tournament"}, nil
	}

	if err := tourney.Start(); err != nil {
		return &pb.TournamentStartResponse{Success: false, Error: err.Error()}, nil
	}

	s.logger.Info("tournament started",
		zap.String("tournament_id", tourney.ID),
		zap.String("controller", player),
	)

	return &pb.TournamentStartResponse{Success: true}, nil
}

// TournamentQuit removes or flags a player as having quit the tournament.
func (s *mageServer) TournamentQuit(ctx context.Context, req *pb.TournamentQuitRequest) (*pb.TournamentQuitResponse, error) {
	player, tourney, err := s.resolveTournamentPlayer(req.GetSessionId(), req.GetTournamentId())
	if err != nil {
		return &pb.TournamentQuitResponse{Success: false, Error: err.Error()}, nil
	}

	state := tourney.GetState()
	if state == tournament.TournamentStateWaiting {
		if err := tourney.RemovePlayer(player); err != nil {
			return &pb.TournamentQuitResponse{Success: false, Error: err.Error()}, nil
		}
	} else {
		if err := tourney.QuitPlayer(player); err != nil {
			return &pb.TournamentQuitResponse{Success: false, Error: err.Error()}, nil
		}
	}

	s.roomMgr.UserLeaveRoom(tourney.RoomID, player)

	if len(tourney.GetPlayers()) == 0 {
		s.tournamentMgr.RemoveTournament(tourney.ID)
	}

	s.logger.Info("player quit tournament",
		zap.String("tournament_id", tourney.ID),
		zap.String("player", player),
	)

	return &pb.TournamentQuitResponse{Success: true}, nil
}

// TournamentFindById returns the tournament view.
func (s *mageServer) TournamentFindById(ctx context.Context, req *pb.TournamentFindByIdRequest) (*pb.TournamentFindByIdResponse, error) {
	_, tourney, err := s.resolveTournamentPlayer(req.GetSessionId(), req.GetTournamentId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	snapshot := tourney.Snapshot()
	return &pb.TournamentFindByIdResponse{
		Tournament: tournamentSnapshotToProto(snapshot),
	}, nil
}

func (s *mageServer) resolveTournamentPlayer(sessionID, tournamentID string) (string, *tournament.Tournament, error) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return "", nil, status.Errorf(codes.InvalidArgument, "session_id is required")
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return "", nil, status.Errorf(codes.NotFound, "session not found")
	}

	user := sess.GetUserID()
	if user == "" {
		return "", nil, status.Errorf(codes.InvalidArgument, "session not associated with a user")
	}

	tournamentID = strings.TrimSpace(tournamentID)
	if tournamentID == "" {
		return "", nil, status.Errorf(codes.InvalidArgument, "tournament_id is required")
	}

	tourney, ok := s.tournamentMgr.GetTournament(tournamentID)
	if !ok {
		return "", nil, status.Errorf(codes.NotFound, "tournament not found")
	}

	return user, tourney, nil
}

func tournamentSnapshotToProto(s tournament.TournamentSnapshot) *pb.TournamentView {
	playerViews := make([]*pb.TournamentPlayerView, 0, len(s.Players))
	for _, p := range s.Players {
		playerViews = append(playerViews, &pb.TournamentPlayerView{
			PlayerName: p.Name,
			Points:     int32(p.Points),
			Wins:       int32(p.Wins),
			Losses:     int32(p.Losses),
			Draws:      int32(p.Draws),
			Eliminated: p.Eliminated,
			Quit:       p.Quit,
			State:      deriveTournamentPlayerState(p),
		})
	}

	roundViews := make([]*pb.RoundView, 0, len(s.Rounds))
	for _, r := range s.Rounds {
		pairingViews := make([]*pb.PairingView, 0, len(r.Pairings))
		for _, p := range r.Pairings {
			pairingViews = append(pairingViews, &pb.PairingView{
				Player1:     p.Player1,
				Player2:     p.Player2,
				TableId:     p.TableID,
				Player1Wins: int32(p.Player1Wins),
				Player2Wins: int32(p.Player2Wins),
			})
		}

		roundViews = append(roundViews, &pb.RoundView{
			RoundNumber: int32(r.Number),
			Pairings:    pairingViews,
			State:       deriveRoundState(r),
		})
	}

	view := &pb.TournamentView{
		TournamentId:   s.ID,
		TournamentName: s.Name,
		TournamentType: s.Type,
		State:          s.State.String(),
		NumPlayers:     int32(len(s.Players)),
		NumRounds:      int32(s.NumRounds),
		CurrentRound:   int32(s.CurrentRound),
		Players:        playerViews,
		Rounds:         roundViews,
	}

	if s.StartTime != nil {
		view.StartTime = timestamppb.New(*s.StartTime)
	}
	if s.EndTime != nil {
		view.EndTime = timestamppb.New(*s.EndTime)
	}

	return view
}

func deriveTournamentPlayerState(p tournament.PlayerSnapshot) string {
	switch {
	case p.Quit:
		return "QUIT"
	case p.Eliminated:
		return "ELIMINATED"
	default:
		return "ACTIVE"
	}
}

func deriveRoundState(r tournament.RoundSnapshot) string {
	switch {
	case r.Finished:
		return "FINISHED"
	case r.Started:
		return "IN_PROGRESS"
	default:
		return "PENDING"
	}
}
