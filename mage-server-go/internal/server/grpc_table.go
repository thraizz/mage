package server

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/magefree/mage-server-go/internal/repository"
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

	// Convert DeckCard messages to card name strings for internal representation
	mainDeckNames := make([]string, 0)
	for _, card := range deck.GetMainDeck() {
		for i := int32(0); i < card.GetQuantity(); i++ {
			mainDeckNames = append(mainDeckNames, card.GetName())
		}
	}

	sideboardNames := make([]string, 0)
	for _, card := range deck.GetSideboard() {
		for i := int32(0); i < card.GetQuantity(); i++ {
			sideboardNames = append(sideboardNames, card.GetName())
		}
	}

	// Validate all card names exist in the database
	allCardNames := append(mainDeckNames, sideboardNames...)
	if err := s.validateCardNames(ctx, allCardNames); err != nil {
		return &pb.DeckSubmitResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	deckList := table.DeckList{
		MainDeck:  mainDeckNames,
		Sideboard: sideboardNames,
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

// DeckSave saves a deck list to the database for later reuse.
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

	format := strings.TrimSpace(req.GetFormat())
	if format == "" {
		format = "Unknown" // Default format if not specified
	}

	deck := req.GetDeck()
	if deck == nil || (len(deck.GetMainDeck()) == 0 && len(deck.GetSideboard()) == 0) {
		return &pb.DeckSaveResponse{
			Success: false,
			Error:   "deck list is required",
		}, nil
	}

	// Get user ID from repository
	user, err := s.userRepo.GetByName(ctx, username)
	if err != nil {
		s.logger.Error("failed to get user for deck save",
			zap.String("username", username),
			zap.Error(err),
		)
		return &pb.DeckSaveResponse{
			Success: false,
			Error:   "user not found",
		}, nil
	}

	// Convert DeckCard messages to card name strings for storage
	var mainDeckNames []string
	for _, card := range deck.GetMainDeck() {
		for i := int32(0); i < card.GetQuantity(); i++ {
			mainDeckNames = append(mainDeckNames, card.GetName())
		}
	}

	var sideboardNames []string
	for _, card := range deck.GetSideboard() {
		for i := int32(0); i < card.GetQuantity(); i++ {
			sideboardNames = append(sideboardNames, card.GetName())
		}
	}

	var commanderNames []string
	for _, card := range deck.GetCommanders() {
		for i := int32(0); i < card.GetQuantity(); i++ {
			commanderNames = append(commanderNames, card.GetName())
		}
	}

	// Validate all card names exist in the database
	allCardNames := append(append(mainDeckNames, sideboardNames...), commanderNames...)
	if err := s.validateCardNames(ctx, allCardNames); err != nil {
		return &pb.DeckSaveResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Create deck in repository
	deckModel := &repository.Deck{
		UserID:      user.ID,
		Name:        deckName,
		Format:      format,
		Description: strings.TrimSpace(req.GetDescription()),
		MainDeck:    mainDeckNames,
		Sideboard:   sideboardNames,
		Commanders:  commanderNames,
	}

	if err := s.deckRepo.Create(ctx, deckModel); err != nil {
		s.logger.Error("failed to save deck to database",
			zap.String("username", username),
			zap.String("deck_name", deckName),
			zap.Error(err),
		)
		return &pb.DeckSaveResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to save deck: %v", err),
		}, nil
	}

	s.logger.Info("deck saved to database",
		zap.String("username", username),
		zap.Int64("deck_id", deckModel.ID),
		zap.String("deck_name", deckName),
		zap.String("format", format),
		zap.Int("main_count", len(deckModel.MainDeck)),
		zap.Int("sideboard_count", len(deckModel.Sideboard)),
	)

	return &pb.DeckSaveResponse{
		Success: true,
		DeckId:  deckModel.ID,
	}, nil
}

// DeckList retrieves all decks for the current user, optionally filtered by format.
func (s *mageServer) DeckList(ctx context.Context, req *pb.DeckListRequest) (*pb.DeckListResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.DeckListResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.DeckListResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.DeckListResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	// Get user ID from repository
	user, err := s.userRepo.GetByName(ctx, username)
	if err != nil {
		s.logger.Error("failed to get user for deck list",
			zap.String("username", username),
			zap.Error(err),
		)
		return &pb.DeckListResponse{
			Success: false,
			Error:   "user not found",
		}, nil
	}

	// Get decks from repository
	var decks []*repository.Deck
	format := strings.TrimSpace(req.GetFormat())
	if format != "" {
		decks, err = s.deckRepo.GetByUserAndFormat(ctx, user.ID, format)
	} else {
		decks, err = s.deckRepo.GetByUser(ctx, user.ID)
	}

	if err != nil {
		s.logger.Error("failed to retrieve decks from database",
			zap.String("username", username),
			zap.String("format", format),
			zap.Error(err),
		)
		return &pb.DeckListResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to retrieve decks: %v", err),
		}, nil
	}

	// Convert to proto format
	deckInfos := make([]*pb.DeckInfo, len(decks))
	for i, deck := range decks {
		deckInfos[i] = &pb.DeckInfo{
			Id:             deck.ID,
			Name:           deck.Name,
			Format:         deck.Format,
			Description:    deck.Description,
			MainDeckCount:  int32(deck.MainDeckCount()),
			SideboardCount: int32(deck.SideboardCount()),
			CreatedAt:      deck.CreatedAt.Unix(),
			UpdatedAt:      deck.UpdatedAt.Unix(),
		}
	}

	s.logger.Info("deck list retrieved",
		zap.String("username", username),
		zap.String("format", format),
		zap.Int("count", len(deckInfos)),
	)

	return &pb.DeckListResponse{
		Success: true,
		Decks:   deckInfos,
	}, nil
}

// DeckDelete deletes a deck from the database.
func (s *mageServer) DeckDelete(ctx context.Context, req *pb.DeckDeleteRequest) (*pb.DeckDeleteResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.DeckDeleteResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.DeckDeleteResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.DeckDeleteResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	if req.GetDeckId() == 0 {
		return &pb.DeckDeleteResponse{
			Success: false,
			Error:   "deck_id is required",
		}, nil
	}

	// Get user ID from repository
	user, err := s.userRepo.GetByName(ctx, username)
	if err != nil {
		s.logger.Error("failed to get user for deck delete",
			zap.String("username", username),
			zap.Error(err),
		)
		return &pb.DeckDeleteResponse{
			Success: false,
			Error:   "user not found",
		}, nil
	}

	// Delete deck (with ownership check)
	if err := s.deckRepo.DeleteByUserAndID(ctx, user.ID, req.GetDeckId()); err != nil {
		s.logger.Error("failed to delete deck",
			zap.String("username", username),
			zap.Int64("deck_id", req.GetDeckId()),
			zap.Error(err),
		)
		return &pb.DeckDeleteResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to delete deck: %v", err),
		}, nil
	}

	s.logger.Info("deck deleted",
		zap.String("username", username),
		zap.Int64("deck_id", req.GetDeckId()),
	)

	return &pb.DeckDeleteResponse{
		Success: true,
	}, nil
}

// DeckGet retrieves a specific deck by ID.
func (s *mageServer) DeckGet(ctx context.Context, req *pb.DeckGetRequest) (*pb.DeckGetResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.DeckGetResponse{
			Success: false,
			Error:   "session_id is required",
		}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.DeckGetResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.DeckGetResponse{
			Success: false,
			Error:   "session not associated with a user",
		}, nil
	}

	if req.GetDeckId() == 0 {
		return &pb.DeckGetResponse{
			Success: false,
			Error:   "deck_id is required",
		}, nil
	}

	// Get user ID from repository
	user, err := s.userRepo.GetByName(ctx, username)
	if err != nil {
		s.logger.Error("failed to get user for deck get",
			zap.String("username", username),
			zap.Error(err),
		)
		return &pb.DeckGetResponse{
			Success: false,
			Error:   "user not found",
		}, nil
	}

	// Get deck from repository
	deck, err := s.deckRepo.GetByID(ctx, req.GetDeckId())
	if err != nil {
		s.logger.Error("failed to get deck",
			zap.String("username", username),
			zap.Int64("deck_id", req.GetDeckId()),
			zap.Error(err),
		)
		return &pb.DeckGetResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to get deck: %v", err),
		}, nil
	}

	// Verify ownership
	if deck.UserID != user.ID {
		s.logger.Warn("user attempted to access deck they don't own",
			zap.String("username", username),
			zap.Int64("deck_id", req.GetDeckId()),
			zap.Int64("deck_user_id", deck.UserID),
			zap.Int64("user_id", user.ID),
		)
		return &pb.DeckGetResponse{
			Success: false,
			Error:   "deck not found",
		}, nil
	}

	s.logger.Info("deck retrieved",
		zap.String("username", username),
		zap.Int64("deck_id", req.GetDeckId()),
	)

	// Convert card names to DeckCard messages with metadata
	mainDeckCards, err := s.buildDeckCardsWithMetadata(ctx, deck.MainDeck)
	if err != nil {
		s.logger.Error("failed to build main deck cards with metadata",
			zap.String("username", username),
			zap.Int64("deck_id", req.GetDeckId()),
			zap.Error(err),
		)
		return &pb.DeckGetResponse{
			Success: false,
			Error:   "failed to fetch card metadata",
		}, nil
	}

	sideboardCards, err := s.buildDeckCardsWithMetadata(ctx, deck.Sideboard)
	if err != nil {
		s.logger.Error("failed to build sideboard cards with metadata",
			zap.String("username", username),
			zap.Int64("deck_id", req.GetDeckId()),
			zap.Error(err),
		)
		return &pb.DeckGetResponse{
			Success: false,
			Error:   "failed to fetch card metadata",
		}, nil
	}

	commanderCards, err := s.buildDeckCardsWithMetadata(ctx, deck.Commanders)
	if err != nil {
		s.logger.Error("failed to build commander cards with metadata",
			zap.String("username", username),
			zap.Int64("deck_id", req.GetDeckId()),
			zap.Error(err),
		)
		return &pb.DeckGetResponse{
			Success: false,
			Error:   "failed to fetch card metadata",
		}, nil
	}

	return &pb.DeckGetResponse{
		Success: true,
		Info: &pb.DeckInfo{
			Id:             deck.ID,
			Name:           deck.Name,
			Format:         deck.Format,
			Description:    deck.Description,
			MainDeckCount:  int32(deck.MainDeckCount()),
			SideboardCount: int32(deck.SideboardCount()),
			CreatedAt:      deck.CreatedAt.Unix(),
			UpdatedAt:      deck.UpdatedAt.Unix(),
		},
		Deck: &pb.DeckCardLists{
			MainDeck:   mainDeckCards,
			Sideboard:  sideboardCards,
			Commanders: commanderCards,
		},
	}, nil
}

// normalizeCardName normalizes a card name to match database format
// Removes punctuation (commas, apostrophes, hyphens) and normalizes spacing
func normalizeCardName(name string) string {
	// Trim whitespace
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	// Replace common punctuation with spaces
	name = strings.ReplaceAll(name, ",", " ")
	name = strings.ReplaceAll(name, "'", "")
	name = strings.ReplaceAll(name, "’", "") // Unicode apostrophe
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "/", " ")
	name = strings.ReplaceAll(name, "//", " ")

	// Normalize multiple spaces to single space
	spaceRegex := regexp.MustCompile(`\s+`)
	name = spaceRegex.ReplaceAllString(name, " ")

	return strings.TrimSpace(name)
}

// validateCardNames checks that all card names exist in the database
func (s *mageServer) validateCardNames(ctx context.Context, cardNames []string) error {
	if len(cardNames) == 0 {
		return nil
	}

	// Normalize and get unique card names to avoid duplicate database queries
	// Map: normalized name -> original name (for error messages)
	uniqueNames := make(map[string]string)
	for _, name := range cardNames {
		normalized := normalizeCardName(name)
		if normalized != "" {
			// Store mapping from normalized to original for error messages
			// If multiple cards normalize to the same name, keep the first original
			if _, exists := uniqueNames[normalized]; !exists {
				uniqueNames[normalized] = name
			}
		}
	}

	// Validate each unique normalized card name (case-insensitive)
	var invalidCards []string
	for normalizedName, originalName := range uniqueNames {
		cards, err := s.cardRepo.GetByNameCaseInsensitive(ctx, normalizedName)
		if err != nil {
			s.logger.Warn("failed to validate card",
				zap.String("card_name", originalName),
				zap.String("normalized", normalizedName),
				zap.Error(err),
			)
			invalidCards = append(invalidCards, originalName)
			continue
		}

		if len(cards) == 0 {
			s.logger.Info("card not found in database",
				zap.String("card_name", originalName),
				zap.String("normalized", normalizedName),
				zap.String("normalized_length", fmt.Sprintf("%d", len(normalizedName))),
				zap.String("normalized_bytes", fmt.Sprintf("%v", []byte(normalizedName))),
			)
			invalidCards = append(invalidCards, originalName)
		} else {
			s.logger.Debug("card found in database",
				zap.String("card_name", originalName),
				zap.String("normalized", normalizedName),
				zap.Int("matches", len(cards)),
			)
		}
	}

	if len(invalidCards) > 0 {
		return fmt.Errorf("invalid card names: %s", strings.Join(invalidCards, ", "))
	}

	return nil
}

// buildDeckCardsWithMetadata converts card names to DeckCard messages with full metadata
func (s *mageServer) buildDeckCardsWithMetadata(ctx context.Context, cardNames []string) ([]*pb.DeckCard, error) {
	// Count card quantities
	cardQuantities := make(map[string]int32)
	for _, cardName := range cardNames {
		cardQuantities[cardName]++
	}

	// Build DeckCard messages with metadata
	var deckCards []*pb.DeckCard
	for cardName, quantity := range cardQuantities {
		// Look up card metadata from database
		cards, err := s.cardRepo.GetByName(ctx, cardName)
		if err != nil {
			s.logger.Warn("failed to get card metadata",
				zap.String("card_name", cardName),
				zap.Error(err),
			)
			// Continue with minimal data if card not found
			deckCards = append(deckCards, &pb.DeckCard{
				Name:     cardName,
				Quantity: quantity,
			})
			continue
		}

		if len(cards) == 0 {
			// Card not in database, use minimal data
			deckCards = append(deckCards, &pb.DeckCard{
				Name:     cardName,
				Quantity: quantity,
			})
			continue
		}

		// Use first printing (could be enhanced to allow set selection)
		cardData := cards[0]

		// Parse types from card type string (e.g., "Creature - Human Wizard" -> ["CREATURE"])
		types := parseCardTypes(cardData.CardType)

		// Parse colors from mana cost (e.g., "{2}{U}{U}" -> ["U"])
		colors := parseColorsFromManaCost(cardData.ManaCost)

		deckCards = append(deckCards, &pb.DeckCard{
			Name:      cardData.Name,
			ManaCost:  cardData.ManaCost,
			CardType:  cardData.CardType,
			Types:     types,
			Colors:    colors,
			Power:     cardData.Power,
			Toughness: cardData.Toughness,
			Quantity:  quantity,
		})
	}

	return deckCards, nil
}

// parseCardTypes extracts main card types from the card type string
func parseCardTypes(cardType string) []string {
	// Simple implementation: extract basic types
	// e.g., "Creature - Human Wizard" -> ["CREATURE"]
	// e.g., "Legendary Artifact" -> ["ARTIFACT"]
	var types []string

	cardTypeUpper := strings.ToUpper(cardType)

	basicTypes := []string{"CREATURE", "INSTANT", "SORCERY", "ENCHANTMENT", "ARTIFACT", "PLANESWALKER", "LAND", "BATTLE"}
	for _, t := range basicTypes {
		if strings.Contains(cardTypeUpper, t) {
			types = append(types, t)
		}
	}

	return types
}

// parseColorsFromManaCost extracts color identity from mana cost
func parseColorsFromManaCost(manaCost string) []string {
	// Extract color symbols from mana cost
	// e.g., "{2}{U}{U}" -> ["U"]
	// e.g., "{W}{U}{B}{R}{G}" -> ["W", "U", "B", "R", "G"]
	colorMap := make(map[string]bool)

	// Look for single-letter color symbols
	for _, color := range []string{"W", "U", "B", "R", "G"} {
		if strings.Contains(manaCost, color) {
			colorMap[color] = true
		}
	}

	// Convert map to slice
	var colors []string
	for color := range colorMap {
		colors = append(colors, color)
	}

	return colors
}
