package server

import (
	"context"
	"strings"

	"github.com/magefree/mage-server-go/internal/draft"
	pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
	"go.uber.org/zap"
)

// DraftJoin adds a player to a draft session.
func (s *mageServer) DraftJoin(ctx context.Context, req *pb.DraftJoinRequest) (*pb.DraftJoinResponse, error) {
	player, draftSession, errMsg := s.resolveDraftPlayer(req.GetSessionId(), req.GetDraftId())
	if errMsg != "" {
		return &pb.DraftJoinResponse{Success: false, Error: errMsg}, nil
	}

	if err := draftSession.AddPlayer(player); err != nil {
		if !strings.Contains(err.Error(), "already in draft") {
			return &pb.DraftJoinResponse{Success: false, Error: err.Error()}, nil
		}
	}

	s.logger.Info("player joined draft",
		zap.String("draft_id", draftSession.ID),
		zap.String("player", player),
	)

	return &pb.DraftJoinResponse{Success: true}, nil
}

// SendDraftCardPick records a player's pick.
func (s *mageServer) SendDraftCardPick(ctx context.Context, req *pb.SendDraftCardPickRequest) (*pb.SendDraftCardPickResponse, error) {
	player, draftSession, errMsg := s.resolveDraftPlayer(req.GetSessionId(), req.GetDraftId())
	if errMsg != "" {
		return &pb.SendDraftCardPickResponse{Success: false, Error: errMsg}, nil
	}

	if req.GetCardId() == "" {
		return &pb.SendDraftCardPickResponse{Success: false, Error: "card_id is required"}, nil
	}

	if err := draftSession.PickCard(player, req.GetCardId()); err != nil {
		return &pb.SendDraftCardPickResponse{Success: false, Error: err.Error()}, nil
	}

	draftSession.SetState(draft.DraftStatePicking)

	if len(req.GetHiddenCards()) > 0 {
		s.logger.Debug("hidden cards submitted with pick",
			zap.String("draft_id", draftSession.ID),
			zap.String("player", player),
			zap.Int("hidden_count", len(req.GetHiddenCards())),
		)
	}

	s.logger.Info("card picked",
		zap.String("draft_id", draftSession.ID),
		zap.String("player", player),
		zap.String("card_id", req.GetCardId()),
	)

	return &pb.SendDraftCardPickResponse{Success: true}, nil
}

// SendDraftCardMark toggles a card mark for the player.
func (s *mageServer) SendDraftCardMark(ctx context.Context, req *pb.SendDraftCardMarkRequest) (*pb.SendDraftCardMarkResponse, error) {
	player, draftSession, errMsg := s.resolveDraftPlayer(req.GetSessionId(), req.GetDraftId())
	if errMsg != "" {
		return &pb.SendDraftCardMarkResponse{Success: false, Error: errMsg}, nil
	}

	if req.GetCardId() == "" {
		return &pb.SendDraftCardMarkResponse{Success: false, Error: "card_id is required"}, nil
	}

	if err := draftSession.MarkCard(player, req.GetCardId()); err != nil {
		return &pb.SendDraftCardMarkResponse{Success: false, Error: err.Error()}, nil
	}

	s.logger.Debug("card mark toggled",
		zap.String("draft_id", draftSession.ID),
		zap.String("player", player),
		zap.String("card_id", req.GetCardId()),
	)

	return &pb.SendDraftCardMarkResponse{Success: true}, nil
}

// DraftSetBoosterLoaded marks that a player loaded their booster.
func (s *mageServer) DraftSetBoosterLoaded(ctx context.Context, req *pb.DraftSetBoosterLoadedRequest) (*pb.DraftSetBoosterLoadedResponse, error) {
	player, draftSession, errMsg := s.resolveDraftPlayer(req.GetSessionId(), req.GetDraftId())
	if errMsg != "" {
		return &pb.DraftSetBoosterLoadedResponse{Success: false, Error: errMsg}, nil
	}

	if err := draftSession.SetBoosterLoaded(player); err != nil {
		return &pb.DraftSetBoosterLoadedResponse{Success: false, Error: err.Error()}, nil
	}

	if draftSession.AllBoostersLoaded() {
		if err := draftSession.PassBoosters(); err != nil {
			return &pb.DraftSetBoosterLoadedResponse{Success: false, Error: err.Error()}, nil
		}
		s.logger.Info("boosters passed",
			zap.String("draft_id", draftSession.ID),
			zap.Int("current_pack", draftSession.CurrentPack),
			zap.Int("current_pick", draftSession.CurrentPick),
		)
	}

	return &pb.DraftSetBoosterLoadedResponse{Success: true}, nil
}

// DraftQuit removes a player from the draft.
func (s *mageServer) DraftQuit(ctx context.Context, req *pb.DraftQuitRequest) (*pb.DraftQuitResponse, error) {
	player, draftSession, errMsg := s.resolveDraftPlayer(req.GetSessionId(), req.GetDraftId())
	if errMsg != "" {
		return &pb.DraftQuitResponse{Success: false, Error: errMsg}, nil
	}

	if err := draftSession.RemovePlayer(player); err != nil {
		return &pb.DraftQuitResponse{Success: false, Error: err.Error()}, nil
	}

	if len(draftSession.PlayerOrder) == 0 {
		s.draftMgr.RemoveDraft(draftSession.ID)
	}

	s.logger.Info("player left draft",
		zap.String("draft_id", draftSession.ID),
		zap.String("player", player),
	)

	return &pb.DraftQuitResponse{Success: true}, nil
}

// resolveDraftPlayer validates session and draft membership.
func (s *mageServer) resolveDraftPlayer(sessionID, draftID string) (string, *draft.Draft, string) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return "", nil, "session_id is required"
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return "", nil, "session not found"
	}

	user := sess.GetUserID()
	if user == "" {
		return "", nil, "session not associated with a user"
	}

	draftID = strings.TrimSpace(draftID)
	if draftID == "" {
		return "", nil, "draft_id is required"
	}

	draftSession, ok := s.draftMgr.GetDraft(draftID)
	if !ok {
		return "", nil, "draft not found"
	}

	return user, draftSession, ""
}
