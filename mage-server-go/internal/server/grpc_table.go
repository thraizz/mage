package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/magefree/mage-server-go/internal/table"
	pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RoomJoinTable allows a player to join an existing table.
func (s *mageServer) RoomJoinTable(ctx context.Context, req *pb.RoomJoinTableRequest) (*pb.RoomJoinTableResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.RoomJoinTableResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.RoomJoinTableResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	username := strings.TrimSpace(req.GetPlayerName())
	if username == "" {
		username = sess.GetUserID()
	}
	if username == "" {
		return &pb.RoomJoinTableResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	tableID := strings.TrimSpace(req.GetTableId())
	if tableID == "" {
		return &pb.RoomJoinTableResponse{
			Success: false,
			Error:   "table_id is required",
		}, nil
	}

	tbl, ok := s.tableMgr.GetTable(tableID)
	if !ok {
		return &pb.RoomJoinTableResponse{
			Success: false,
			Error:   "table not found",
		}, nil
	}

	if requestedRoom := strings.TrimSpace(req.GetRoomId()); requestedRoom != "" && tbl.RoomID != requestedRoom {
		return &pb.RoomJoinTableResponse{
			Success: false,
			Error:   "table not found in requested room",
		}, nil
	}

	if tbl.Password != "" && tbl.Password != req.GetPassword() && !sess.IsAdminSession() && !tbl.IsController(sess.GetUserID()) {
		return &pb.RoomJoinTableResponse{
			Success: false,
			Error:   "invalid table password",
		}, nil
	}

	playerType := strings.TrimSpace(req.GetPlayerType())
	if playerType == "" {
		playerType = "Human"
	}

	if err := tbl.AddPlayer(username, playerType); err != nil {
		return &pb.RoomJoinTableResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	if err := s.roomMgr.UserJoinRoom(tbl.RoomID, username); err != nil {
		s.logger.Debug("failed to add user to room after table join",
			zap.String("room_id", tbl.RoomID),
			zap.String("table_id", tbl.ID),
			zap.String("username", username),
			zap.Error(err),
		)
	}

	s.logger.Info("user joined table",
		zap.String("table_id", tbl.ID),
		zap.String("room_id", tbl.RoomID),
		zap.String("username", username),
		zap.String("player_type", playerType),
	)

	return &pb.RoomJoinTableResponse{
		Success: true,
	}, nil
}

// RoomLeaveTableOrTournament removes a player from a table or tournament.
func (s *mageServer) RoomLeaveTableOrTournament(ctx context.Context, req *pb.RoomLeaveTableOrTournamentRequest) (*pb.RoomLeaveTableOrTournamentResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.RoomLeaveTableOrTournamentResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.RoomLeaveTableOrTournamentResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.RoomLeaveTableOrTournamentResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	targetID := strings.TrimSpace(req.GetTableId())
	if targetID == "" {
		return &pb.RoomLeaveTableOrTournamentResponse{
			Success: false,
			Error:   "table_id is required",
		}, nil
	}

	if tbl, ok := s.tableMgr.GetTable(targetID); ok {
		if err := tbl.RemovePlayer(username); err != nil {
			// The user might be a spectator; remove silently.
			tbl.RemoveSpectator(username)
		}

		if tbl.RoomID != "" {
			s.roomMgr.UserLeaveRoom(tbl.RoomID, username)
		}

		s.logger.Info("user left table",
			zap.String("table_id", tbl.ID),
			zap.String("room_id", tbl.RoomID),
			zap.String("username", username),
		)

		return &pb.RoomLeaveTableOrTournamentResponse{Success: true}, nil
	}

	if tournament, ok := s.tournamentMgr.GetTournament(targetID); ok {
		leftAsPlayer := true
		if err := tournament.RemovePlayer(username); err != nil {
			if !tournament.RemoveWatcher(username) {
				return &pb.RoomLeaveTableOrTournamentResponse{
					Success: false,
					Error:   err.Error(),
				}, nil
			}
			leftAsPlayer = false
		}

		if tournament.RoomID != "" {
			s.roomMgr.UserLeaveRoom(tournament.RoomID, username)
		}

		if leftAsPlayer {
			s.logger.Info("user left tournament",
				zap.String("tournament_id", tournament.ID),
				zap.String("username", username),
			)
		} else {
			s.logger.Info("user stopped watching tournament",
				zap.String("tournament_id", tournament.ID),
				zap.String("username", username),
			)
		}

		return &pb.RoomLeaveTableOrTournamentResponse{Success: true}, nil
	}

	return &pb.RoomLeaveTableOrTournamentResponse{
		Success: false,
		Error:   "table or tournament not found",
	}, nil
}

// RoomWatchTable registers a user as a spectator for a table.
func (s *mageServer) RoomWatchTable(ctx context.Context, req *pb.RoomWatchTableRequest) (*pb.RoomWatchTableResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.RoomWatchTableResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.RoomWatchTableResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.RoomWatchTableResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	tableID := strings.TrimSpace(req.GetTableId())
	if tableID == "" {
		return &pb.RoomWatchTableResponse{
			Success: false,
			Error:   "table_id is required",
		}, nil
	}

	tbl, ok := s.tableMgr.GetTable(tableID)
	if !ok {
		return &pb.RoomWatchTableResponse{
			Success: false,
			Error:   "table not found",
		}, nil
	}

	if requestedRoom := strings.TrimSpace(req.GetRoomId()); requestedRoom != "" && tbl.RoomID != requestedRoom {
		return &pb.RoomWatchTableResponse{
			Success: false,
			Error:   "table not found in requested room",
		}, nil
	}

	tbl.AddSpectator(username)

	s.logger.Info("user watching table",
		zap.String("table_id", tbl.ID),
		zap.String("room_id", tbl.RoomID),
		zap.String("username", username),
	)

	return &pb.RoomWatchTableResponse{
		Success: true,
	}, nil
}

// TableSwapSeats allows the controller (or admin) to swap seats at a table.
func (s *mageServer) TableSwapSeats(ctx context.Context, req *pb.TableSwapSeatsRequest) (*pb.TableSwapSeatsResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.TableSwapSeatsResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.TableSwapSeatsResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	tableID := strings.TrimSpace(req.GetTableId())
	if tableID == "" {
		return &pb.TableSwapSeatsResponse{
			Success: false,
			Error:   "table_id is required",
		}, nil
	}

	tbl, ok := s.tableMgr.GetTable(tableID)
	if !ok {
		return &pb.TableSwapSeatsResponse{
			Success: false,
			Error:   "table not found",
		}, nil
	}

	if !sess.IsAdminSession() && !tbl.IsController(sess.GetUserID()) {
		return &pb.TableSwapSeatsResponse{
			Success: false,
			Error:   "only table controller or admin can swap seats",
		}, nil
	}

	seat1 := int(req.GetSeatNum1())
	seat2 := int(req.GetSeatNum2())

	if seat1 == seat2 {
		return &pb.TableSwapSeatsResponse{Success: true}, nil
	}

	if err := tbl.SwapSeats(seat1, seat2); err != nil {
		return &pb.TableSwapSeatsResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	s.logger.Info("table seats swapped",
		zap.String("table_id", tbl.ID),
		zap.Int("seat1", seat1),
		zap.Int("seat2", seat2),
		zap.String("username", sess.GetUserID()),
	)

	return &pb.TableSwapSeatsResponse{Success: true}, nil
}

// TableRemove removes a table. Only the controller or admin may remove.
func (s *mageServer) TableRemove(ctx context.Context, req *pb.TableRemoveRequest) (*pb.TableRemoveResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.TableRemoveResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.TableRemoveResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	tableID := strings.TrimSpace(req.GetTableId())
	if tableID == "" {
		return &pb.TableRemoveResponse{
			Success: false,
			Error:   "table_id is required",
		}, nil
	}

	tbl, ok := s.tableMgr.GetTable(tableID)
	if !ok {
		return &pb.TableRemoveResponse{
			Success: false,
			Error:   "table not found",
		}, nil
	}

	if !sess.IsAdminSession() && !tbl.IsController(sess.GetUserID()) {
		return &pb.TableRemoveResponse{
			Success: false,
			Error:   "only table controller or admin can remove table",
		}, nil
	}

	s.tableMgr.RemoveTable(tableID)

	s.logger.Info("table removed",
		zap.String("table_id", tableID),
		zap.String("room_id", tbl.RoomID),
		zap.String("username", sess.GetUserID()),
	)

	return &pb.TableRemoveResponse{Success: true}, nil
}

// TableIsOwner verifies if the caller controls the table.
func (s *mageServer) TableIsOwner(ctx context.Context, req *pb.TableIsOwnerRequest) (*pb.TableIsOwnerResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "session_id is required")
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "session not found")
	}

	tableID := strings.TrimSpace(req.GetTableId())
	if tableID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "table_id is required")
	}

	tbl, ok := s.tableMgr.GetTable(tableID)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "table not found")
	}

	isOwner := tbl.IsController(sess.GetUserID()) || sess.IsAdminSession()

	return &pb.TableIsOwnerResponse{
		IsOwner: isOwner,
	}, nil
}

func (s *mageServer) RoomCreateTournament(ctx context.Context, req *pb.RoomCreateTournamentRequest) (*pb.RoomCreateTournamentResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.RoomCreateTournamentResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.RoomCreateTournamentResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	controller := sess.GetUserID()
	if controller == "" {
		return &pb.RoomCreateTournamentResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	roomID := strings.TrimSpace(req.GetRoomId())
	if roomID == "" {
		roomID = s.roomMgr.GetMainRoomID()
	}

	if _, exists := s.roomMgr.GetRoom(roomID); !exists {
		s.roomMgr.CreateRoom(roomID, fmt.Sprintf("Room %s", roomID))
	}

	opts := req.GetTournamentOptions()

	name := fmt.Sprintf("%s's tournament", controller)
	if opts != nil && strings.TrimSpace(opts.GetName()) != "" {
		name = opts.GetName()
	}

	tournamentType := "Constructed"
	if opts != nil && strings.TrimSpace(opts.GetTournamentType()) != "" {
		tournamentType = opts.GetTournamentType()
	}

	numRounds := int32(3)
	if opts != nil && opts.GetNumRounds() > 0 {
		numRounds = opts.GetNumRounds()
	}

	winsRequired := int32(2)
	if opts != nil && opts.GetNumWins() > 0 {
		winsRequired = opts.GetNumWins()
	}

	tournament := s.tournamentMgr.CreateTournament(
		name,
		tournamentType,
		controller,
		roomID,
		int(numRounds),
		int(winsRequired),
	)

	if err := tournament.AddPlayer(controller); err != nil {
		s.logger.Warn("failed to add controller to tournament",
			zap.String("tournament_id", tournament.ID),
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

	s.logger.Info("tournament created",
		zap.String("tournament_id", tournament.ID),
		zap.String("room_id", roomID),
		zap.String("controller", controller),
		zap.String("name", tournament.Name),
		zap.String("type", tournament.TournamentTypeStr),
	)

	return &pb.RoomCreateTournamentResponse{
		Success:      true,
		TournamentId: tournament.ID,
	}, nil
}

func (s *mageServer) RoomJoinTournament(ctx context.Context, req *pb.RoomJoinTournamentRequest) (*pb.RoomJoinTournamentResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.RoomJoinTournamentResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.RoomJoinTournamentResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.RoomJoinTournamentResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	tournamentID := strings.TrimSpace(req.GetTournamentId())
	if tournamentID == "" {
		return &pb.RoomJoinTournamentResponse{
			Success: false,
			Error:   "tournament_id is required",
		}, nil
	}

	tournament, ok := s.tournamentMgr.GetTournament(tournamentID)
	if !ok {
		return &pb.RoomJoinTournamentResponse{
			Success: false,
			Error:   "tournament not found",
		}, nil
	}

	if requestedRoom := strings.TrimSpace(req.GetRoomId()); requestedRoom != "" && tournament.RoomID != requestedRoom {
		return &pb.RoomJoinTournamentResponse{
			Success: false,
			Error:   "tournament not found in requested room",
		}, nil
	}

	if err := tournament.AddPlayer(username); err != nil {
		return &pb.RoomJoinTournamentResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	if err := s.roomMgr.UserJoinRoom(tournament.RoomID, username); err != nil {
		s.logger.Debug("failed to ensure player present in room",
			zap.String("room_id", tournament.RoomID),
			zap.String("username", username),
			zap.Error(err),
		)
	}

	s.logger.Info("user joined tournament",
		zap.String("tournament_id", tournament.ID),
		zap.String("room_id", tournament.RoomID),
		zap.String("username", username),
		zap.String("player_type", strings.TrimSpace(req.GetPlayerType())),
	)

	return &pb.RoomJoinTournamentResponse{
		Success: true,
	}, nil
}

func (s *mageServer) RoomWatchTournament(ctx context.Context, req *pb.RoomWatchTournamentRequest) (*pb.RoomWatchTournamentResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.RoomWatchTournamentResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.RoomWatchTournamentResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.RoomWatchTournamentResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	tournamentID := strings.TrimSpace(req.GetTournamentId())
	if tournamentID == "" {
		return &pb.RoomWatchTournamentResponse{
			Success: false,
			Error:   "tournament_id is required",
		}, nil
	}

	tournament, ok := s.tournamentMgr.GetTournament(tournamentID)
	if !ok {
		return &pb.RoomWatchTournamentResponse{
			Success: false,
			Error:   "tournament not found",
		}, nil
	}

	if requestedRoom := strings.TrimSpace(req.GetRoomId()); requestedRoom != "" && tournament.RoomID != requestedRoom {
		return &pb.RoomWatchTournamentResponse{
			Success: false,
			Error:   "tournament not found in requested room",
		}, nil
	}

	tournament.AddWatcher(username)

	if err := s.roomMgr.UserJoinRoom(tournament.RoomID, username); err != nil {
		s.logger.Debug("failed to ensure watcher present in room",
			zap.String("room_id", tournament.RoomID),
			zap.String("username", username),
			zap.Error(err),
		)
	}

	s.logger.Info("user watching tournament",
		zap.String("tournament_id", tournament.ID),
		zap.String("room_id", tournament.RoomID),
		zap.String("username", username),
	)

	return &pb.RoomWatchTournamentResponse{
		Success: true,
	}, nil
}

// DeckSubmit validates and stores a player's deck against a table.
func (s *mageServer) DeckSubmit(ctx context.Context, req *pb.DeckSubmitRequest) (*pb.DeckSubmitResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.DeckSubmitResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.DeckSubmitResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.DeckSubmitResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	tableID := strings.TrimSpace(req.GetTableId())
	if tableID == "" {
		return &pb.DeckSubmitResponse{
			Success: false,
			Error:   "table_id is required",
		}, nil
	}

	tbl, ok := s.tableMgr.GetTable(tableID)
	if !ok {
		return &pb.DeckSubmitResponse{
			Success: false,
			Error:   "table not found",
		}, nil
	}

	deck := req.GetDeck()
	if deck == nil || (len(deck.GetMainDeck()) == 0 && len(deck.GetSideboard()) == 0) {
		return &pb.DeckSubmitResponse{
			Success: false,
			Error:   "deck list is required",
		}, nil
	}

	deckList := table.DeckList{
		MainDeck:  append([]string(nil), deck.GetMainDeck()...),
		Sideboard: append([]string(nil), deck.GetSideboard()...),
	}

	if err := tbl.SubmitDeck(username, deckList); err != nil {
		return &pb.DeckSubmitResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	s.logger.Info("deck submitted",
		zap.String("table_id", tbl.ID),
		zap.String("username", username),
		zap.Int("main_count", len(deckList.MainDeck)),
		zap.Int("sideboard_count", len(deckList.Sideboard)),
	)

	return &pb.DeckSubmitResponse{
		Success: true,
	}, nil
}

const maxSavedDecksPerUser = 20

// DeckSave saves a deck list for later reuse.
func (s *mageServer) DeckSave(ctx context.Context, req *pb.DeckSaveRequest) (*pb.DeckSaveResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.DeckSaveResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.DeckSaveResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.DeckSaveResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	deckName := strings.TrimSpace(req.GetDeckName())
	if deckName == "" {
		return &pb.DeckSaveResponse{
			Success: false,
			Error:   "deck_name is required",
		}, nil
	}

	deck := req.GetDeck()
	if deck == nil || (len(deck.GetMainDeck()) == 0 && len(deck.GetSideboard()) == 0) {
		return &pb.DeckSaveResponse{
			Success: false,
			Error:   "deck list is required",
		}, nil
	}

	entry := savedDeck{
		Name: deckName,
		Deck: table.DeckList{
			MainDeck:  append([]string(nil), deck.GetMainDeck()...),
			Sideboard: append([]string(nil), deck.GetSideboard()...),
		},
	}

	s.savedDecksMu.Lock()
	defer s.savedDecksMu.Unlock()

	decks := append(s.savedDecks[username], entry)
	if len(decks) > maxSavedDecksPerUser {
		decks = decks[len(decks)-maxSavedDecksPerUser:]
	}
	s.savedDecks[username] = decks

	s.logger.Info("deck saved",
		zap.String("username", username),
		zap.String("deck_name", deckName),
		zap.Int("main_count", len(entry.Deck.MainDeck)),
		zap.Int("sideboard_count", len(entry.Deck.Sideboard)),
		zap.Int("total_saved", len(decks)),
	)

	return &pb.DeckSaveResponse{
		Success: true,
	}, nil
}
