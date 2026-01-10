package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/session"
	"github.com/magefree/mage-server-go/internal/table"
	pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MatchStart promotes a table to an active game.
func (s *mageServer) MatchStart(ctx context.Context, req *pb.MatchStartRequest) (*pb.MatchStartResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.MatchStartResponse{Success: false, Error: "session_id is required"}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.MatchStartResponse{Success: false, Error: "session not found"}, nil
	}

	tableID := strings.TrimSpace(req.GetTableId())
	if tableID == "" {
		return &pb.MatchStartResponse{Success: false, Error: "table_id is required"}, nil
	}

	tbl, ok := s.tableMgr.GetTable(tableID)
	if !ok {
		return &pb.MatchStartResponse{Success: false, Error: "table not found"}, nil
	}

	if !sess.IsAdminSession() && !tbl.IsController(sess.GetUserID()) {
		return &pb.MatchStartResponse{Success: false, Error: "only table controller or admin can start match"}, nil
	}

	if _, exists := s.gameMgr.GetGameByTable(tableID); exists {
		return &pb.MatchStartResponse{Success: false, Error: "game already active for table"}, nil
	}

	if tbl.GetState() != table.TableStateWaiting && tbl.GetState() != table.TableStateStarting {
		return &pb.MatchStartResponse{Success: false, Error: "table not ready to start"}, nil
	}

	s.logger.Info("[MATCH START] Beginning match start",
		zap.String("table_id", tbl.ID),
		zap.Int("seat_count", len(tbl.Seats)),
	)

	players := make([]string, 0, len(tbl.Seats))
	for i, seat := range tbl.Seats {
		s.logger.Debug("[MATCH START] Checking seat",
			zap.String("table_id", tbl.ID),
			zap.Int("seat_index", i),
			zap.String("player_name", seat.PlayerName),
			zap.Bool("deck_valid", seat.DeckValid),
		)
		if seat.PlayerName != "" {
			players = append(players, seat.PlayerName)
			if !seat.DeckValid {
				s.logger.Warn("[MATCH START] Player has NO valid deck submitted!",
					zap.String("table_id", tbl.ID),
					zap.String("player", seat.PlayerName),
				)
			} else {
				s.logger.Info("[MATCH START] Player has valid deck",
					zap.String("table_id", tbl.ID),
					zap.String("player", seat.PlayerName),
				)
			}
		}
	}

	if len(players) < 2 {
		return &pb.MatchStartResponse{Success: false, Error: "not enough players to start match"}, nil
	}

	// Collect submitted decks from the table
	s.logger.Info("[MATCH START] Collecting submitted decks",
		zap.String("table_id", tbl.ID),
		zap.Strings("players", players),
	)

	decks := make(map[string]game.DeckList)
	for _, playerName := range players {
		if tableDeck, ok := tbl.GetSubmittedDeck(playerName); ok {
			decks[playerName] = game.DeckList{
				MainDeck:   tableDeck.MainDeck,
				Sideboard:  tableDeck.Sideboard,
				Commanders: tableDeck.Commanders,
			}
			s.logger.Info("[MATCH START] Got deck for player",
				zap.String("player", playerName),
				zap.Int("main_deck_size", len(tableDeck.MainDeck)),
				zap.Int("sideboard_size", len(tableDeck.Sideboard)),
				zap.Int("commander_count", len(tableDeck.Commanders)),
				zap.Strings("first_5_cards", firstN(tableDeck.MainDeck, 5)),
			)
		} else {
			s.logger.Warn("[MATCH START] NO DECK found for player!",
				zap.String("player", playerName),
			)
		}
	}

	s.logger.Info("[MATCH START] Decks collection complete",
		zap.String("table_id", tbl.ID),
		zap.Int("decks_collected", len(decks)),
		zap.Int("players_count", len(players)),
	)

	gameInstance := s.gameMgr.CreateGame(tbl.ID, tbl.GameType, players)
	tbl.RecordMatch(gameInstance.ID)
	tbl.SetState(table.TableStateDueling)

	if s.gameAdapter != nil {
		// Use StartGameWithDecks to pass player decks to the engine
		if err := s.gameAdapter.StartGameWithDecks(gameInstance, decks); err != nil {
			s.logger.Warn("failed to start game engine",
				zap.String("game_id", gameInstance.ID),
				zap.Error(err),
			)
		}
		go s.gameAdapter.ProcessGameActions(gameInstance)
	}

	s.logger.Info("match started",
		zap.String("table_id", tbl.ID),
		zap.String("game_id", gameInstance.ID),
		zap.Strings("players", players),
	)

	// Send START_GAME WebSocket event to all players in the table
	s.notifyGameStart(gameInstance.ID, players)

	return &pb.MatchStartResponse{Success: true, GameId: gameInstance.ID}, nil
}

// notifyGameStart sends START_GAME WebSocket events to all players
func (s *mageServer) notifyGameStart(gameID string, playerNames []string) {
	// Create the START_GAME event data
	startGameData := &pb.StartGameData{
		GameId:      gameID,
		PlayerNames: playerNames,
	}

	// Create the ServerEvent
	event := &pb.ServerEvent{
		ObjectId: gameID,
		Method:   pb.CallbackMethod_START_GAME,
	}

	// Marshal the data into Any
	anyData, err := anypb.New(startGameData)
	if err != nil {
		s.logger.Error("failed to marshal StartGameData",
			zap.String("game_id", gameID),
			zap.Error(err),
		)
		return
	}
	event.Data = anyData

	// Send to all players
	for _, playerName := range playerNames {
		sessions := s.sessionMgr.GetSessionsByUser(playerName)
		for _, sess := range sessions {
			if !sess.SendCallback(event) {
				s.logger.Warn("failed to send START_GAME event to player",
					zap.String("game_id", gameID),
					zap.String("player", playerName),
				)
			} else {
				s.logger.Debug("sent START_GAME event to player",
					zap.String("game_id", gameID),
					zap.String("player", playerName),
				)
			}
		}
	}
}

// GameJoin registers a player to an active game session.
func (s *mageServer) GameJoin(ctx context.Context, req *pb.GameJoinRequest) (*pb.GameJoinResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.GameJoinResponse{Success: false, Error: "session_id is required"}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.GameJoinResponse{Success: false, Error: "session not found"}, nil
	}

	gameID := strings.TrimSpace(req.GetGameId())
	if gameID == "" {
		return &pb.GameJoinResponse{Success: false, Error: "game_id is required"}, nil
	}

	game, ok := s.gameMgr.GetGame(gameID)
	if !ok {
		// Game not found in active games - check if it's a completed game
		if s.matchHistoryRepo != nil {
			exists, err := s.matchHistoryRepo.ExistsByGameID(ctx, gameID)
			if err != nil {
				s.logger.Warn("failed to check match history",
					zap.String("game_id", gameID),
					zap.Error(err),
				)
			} else if exists {
				return &pb.GameJoinResponse{Success: false, Error: "game has ended"}, nil
			}
		}
		return &pb.GameJoinResponse{Success: false, Error: "game not found"}, nil
	}

	user := sess.GetUserID()
	if user == "" {
		return &pb.GameJoinResponse{Success: false, Error: "session not associated with a user"}, nil
	}

	if !game.IsPlayer(user) {
		return &pb.GameJoinResponse{Success: false, Error: "player not part of this game"}, nil
	}

	s.logger.Info("player joined game",
		zap.String("game_id", game.ID),
		zap.String("player", user),
	)

	return &pb.GameJoinResponse{Success: true}, nil
}

// GameWatchStart registers a watcher for a running game.
func (s *mageServer) GameWatchStart(ctx context.Context, req *pb.GameWatchStartRequest) (*pb.GameWatchStartResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.GameWatchStartResponse{Success: false, Error: "session_id is required"}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.GameWatchStartResponse{Success: false, Error: "session not found"}, nil
	}

	gameID := strings.TrimSpace(req.GetGameId())
	if gameID == "" {
		return &pb.GameWatchStartResponse{Success: false, Error: "game_id is required"}, nil
	}

	game, ok := s.gameMgr.GetGame(gameID)
	if !ok {
		return &pb.GameWatchStartResponse{Success: false, Error: "game not found"}, nil
	}

	user := sess.GetUserID()
	if user == "" {
		return &pb.GameWatchStartResponse{Success: false, Error: "session not associated with a user"}, nil
	}

	game.AddWatcher(user)
	s.logger.Info("watcher added to game",
		zap.String("game_id", game.ID),
		zap.String("username", user),
	)

	return &pb.GameWatchStartResponse{Success: true}, nil
}

// GameWatchStop removes a watcher from a game.
func (s *mageServer) GameWatchStop(ctx context.Context, req *pb.GameWatchStopRequest) (*pb.GameWatchStopResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.GameWatchStopResponse{Success: false, Error: "session_id is required"}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.GameWatchStopResponse{Success: false, Error: "session not found"}, nil
	}

	gameID := strings.TrimSpace(req.GetGameId())
	if gameID == "" {
		return &pb.GameWatchStopResponse{Success: false, Error: "game_id is required"}, nil
	}

	game, ok := s.gameMgr.GetGame(gameID)
	if !ok {
		return &pb.GameWatchStopResponse{Success: false, Error: "game not found"}, nil
	}

	user := sess.GetUserID()
	if user == "" {
		return &pb.GameWatchStopResponse{Success: false, Error: "session not associated with a user"}, nil
	}

	game.RemoveWatcher(user)
	s.logger.Info("watcher removed from game",
		zap.String("game_id", game.ID),
		zap.String("username", user),
	)

	return &pb.GameWatchStopResponse{Success: true}, nil
}

// GameGetView returns a simplified game snapshot.
func (s *mageServer) GameGetView(ctx context.Context, req *pb.GameGetViewRequest) (*pb.GameGetViewResponse, error) {
	s.logger.Info("GameGetView called",
		zap.String("session_id", req.GetSessionId()),
		zap.String("game_id", req.GetGameId()),
		zap.String("player_id", req.GetPlayerId()),
	)

	sess, gameInstance, err := s.resolveGameAccess(req.GetSessionId(), req.GetGameId(), true)
	if err != nil {
		s.logger.Error("GameGetView resolveGameAccess failed", zap.Error(err))
		return nil, err
	}

	s.logger.Info("GameGetView resolved game access",
		zap.String("game_id", gameInstance.ID),
		zap.Int("num_players", len(gameInstance.Players)),
	)

	gamePlayers := make([]*pb.PlayerView, 0, len(gameInstance.Players))
	for _, name := range gameInstance.Players {
		gamePlayers = append(gamePlayers, &pb.PlayerView{
			PlayerId: name,
			Name:     name,
		})
	}

	view := &pb.GameView{
		GameId:           gameInstance.ID,
		State:            gameInstance.GetState().String(),
		Players:          gamePlayers,
		ActivePlayerId:   gameInstance.ActivePlayerID,
		PriorityPlayerId: gameInstance.PriorityPlayer,
		Turn:             int32(gameInstance.Turn),
		StartTime:        timestamppb.New(gameInstance.StartTime),
	}

	playerID := strings.TrimSpace(req.GetPlayerId())
	if playerID == "" && sess != nil {
		playerID = sess.GetUserID()
	}

	if s.gameAdapter != nil {
		s.logger.Info("GameGetView calling gameAdapter.GetGameView",
			zap.String("game_id", gameInstance.ID),
			zap.String("player_id", playerID),
		)
		if engineView, engineErr := s.gameAdapter.GetGameView(gameInstance.ID, playerID); engineErr == nil && engineView != nil {
			s.logger.Info("GameGetView got engine view")
			switch data := engineView.(type) {
			case *game.EngineGameView:
				if data.GameID != "" {
					view.GameId = data.GameID
				}
				view.State = data.State.String()
				view.Phase = data.Phase
				view.Step = data.Step
				view.Turn = int32(data.Turn)
				view.ActivePlayerId = data.ActivePlayerID
				view.PriorityPlayerId = data.PriorityPlayer
				view.Players = enginePlayersToProto(data.Players)
				view.Battlefield = engineCardsToProto(data.Battlefield)
				view.Stack = engineCardsToProto(data.Stack)
				view.Exile = engineCardsToProto(data.Exile)
				view.Command = engineCardsToProto(data.Command)
				view.Revealed = engineRevealedToProto(data.Revealed)
				view.LookedAt = engineLookedAtToProto(data.LookedAt)
				if combat := engineCombatToProto(data.Combat); combat != nil {
					view.Combat = combat
				}
				if !data.StartedAt.IsZero() {
					view.StartTime = timestamppb.New(data.StartedAt)
				}

				// Add pre-computed display values (server source of truth)
				view.ActivePlayerName = data.ActivePlayerName
				view.PriorityPlayerName = data.PriorityPlayerName
				view.GameFormat = data.GameFormat
				view.IsMulliganPhase = data.IsMulliganPhase
				view.LandsPlayedThisTurn = int32(data.LandsPlayedThisTurn)
				view.LandsAllowedThisTurn = int32(data.LandsAllowedThisTurn)

				nextID := int32(len(view.Messages) + 1)
				engineMessages := engineMessagesToProto(data.Messages, nextID)
				view.Messages = append(view.Messages, engineMessages...)
				nextID += int32(len(engineMessages))
				view.Messages = append(view.Messages, enginePromptsToMessages(data.Prompts, nextID)...)
			case game.NullGameView:
				for idx, action := range data.Actions {
					view.Messages = append(view.Messages, &pb.GameMessage{
						Id:   int32(idx + 1),
						Text: fmt.Sprintf("%s %s %v", action.PlayerID, action.ActionType, action.Data),
						Time: timestamppb.New(action.Timestamp),
					})
				}
			}
		}
	}

	return &pb.GameGetViewResponse{
		Game: view,
	}, nil
}

// MatchQuit ends a game for a player or admin.
func (s *mageServer) MatchQuit(ctx context.Context, req *pb.MatchQuitRequest) (*pb.MatchQuitResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.MatchQuitResponse{Success: false, Error: "session_id is required"}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.MatchQuitResponse{Success: false, Error: "session not found"}, nil
	}

	gameID := strings.TrimSpace(req.GetGameId())
	if gameID == "" {
		return &pb.MatchQuitResponse{Success: false, Error: "game_id is required"}, nil
	}

	game, ok := s.gameMgr.GetGame(gameID)
	if !ok {
		return &pb.MatchQuitResponse{Success: false, Error: "game not found"}, nil
	}

	user := sess.GetUserID()
	if user == "" {
		return &pb.MatchQuitResponse{Success: false, Error: "session not associated with a user"}, nil
	}

	if !sess.IsAdminSession() && !game.IsPlayer(user) {
		return &pb.MatchQuitResponse{Success: false, Error: "player not part of this game"}, nil
	}

	if s.gameAdapter != nil {
		if err := s.gameAdapter.EndGame(game, user); err != nil {
			s.logger.Debug("failed to end game via engine",
				zap.String("game_id", gameID),
				zap.Error(err),
			)
		}
	}

	s.gameMgr.RemoveGame(gameID)

	if tbl, ok := s.tableMgr.GetTable(game.TableID); ok {
		tbl.SetState(table.TableStateFinished)
	}

	s.logger.Info("match ended",
		zap.String("game_id", gameID),
		zap.String("ended_by", user),
	)

	return &pb.MatchQuitResponse{Success: true}, nil
}

// SendPlayerUUID forwards a UUID selection to the game engine.
func (s *mageServer) SendPlayerUUID(ctx context.Context, req *pb.SendPlayerUUIDRequest) (*pb.SendPlayerUUIDResponse, error) {
	player, gameInstance, errMsg := s.resolveGamePlayer(req.GetSessionId(), req.GetGameId())
	if errMsg != "" {
		return &pb.SendPlayerUUIDResponse{Success: false, Error: errMsg}, nil
	}

	if req.GetUuid() == "" {
		return &pb.SendPlayerUUIDResponse{Success: false, Error: "uuid is required"}, nil
	}

	if err := s.gameMgr.SendPlayerAction(gameInstance.ID, player, "SEND_UUID", req.GetUuid()); err != nil {
		return &pb.SendPlayerUUIDResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.SendPlayerUUIDResponse{Success: true}, nil
}

// SendPlayerString forwards string data to the game engine.
func (s *mageServer) SendPlayerString(ctx context.Context, req *pb.SendPlayerStringRequest) (*pb.SendPlayerStringResponse, error) {
	player, gameInstance, errMsg := s.resolveGamePlayer(req.GetSessionId(), req.GetGameId())
	if errMsg != "" {
		return &pb.SendPlayerStringResponse{Success: false, Error: errMsg}, nil
	}

	if req.GetData() == "" {
		return &pb.SendPlayerStringResponse{Success: false, Error: "data is required"}, nil
	}

	if err := s.gameMgr.SendPlayerAction(gameInstance.ID, player, "SEND_STRING", req.GetData()); err != nil {
		return &pb.SendPlayerStringResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.SendPlayerStringResponse{Success: true}, nil
}

// SendPlayerBoolean forwards boolean data to the game engine.
func (s *mageServer) SendPlayerBoolean(ctx context.Context, req *pb.SendPlayerBooleanRequest) (*pb.SendPlayerBooleanResponse, error) {
	player, gameInstance, errMsg := s.resolveGamePlayer(req.GetSessionId(), req.GetGameId())
	if errMsg != "" {
		return &pb.SendPlayerBooleanResponse{Success: false, Error: errMsg}, nil
	}

	if err := s.gameMgr.SendPlayerAction(gameInstance.ID, player, "SEND_BOOLEAN", req.GetData()); err != nil {
		return &pb.SendPlayerBooleanResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.SendPlayerBooleanResponse{Success: true}, nil
}

// SendPlayerInteger forwards integer data to the game engine.
func (s *mageServer) SendPlayerInteger(ctx context.Context, req *pb.SendPlayerIntegerRequest) (*pb.SendPlayerIntegerResponse, error) {
	player, gameInstance, errMsg := s.resolveGamePlayer(req.GetSessionId(), req.GetGameId())
	if errMsg != "" {
		return &pb.SendPlayerIntegerResponse{Success: false, Error: errMsg}, nil
	}

	if err := s.gameMgr.SendPlayerAction(gameInstance.ID, player, "SEND_INTEGER", req.GetData()); err != nil {
		return &pb.SendPlayerIntegerResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.SendPlayerIntegerResponse{Success: true}, nil
}

// SendPlayerManaType forwards a mana selection to the game engine.
func (s *mageServer) SendPlayerManaType(ctx context.Context, req *pb.SendPlayerManaTypeRequest) (*pb.SendPlayerManaTypeResponse, error) {
	player, gameInstance, errMsg := s.resolveGamePlayer(req.GetSessionId(), req.GetGameId())
	if errMsg != "" {
		return &pb.SendPlayerManaTypeResponse{Success: false, Error: errMsg}, nil
	}

	payload := map[string]string{
		"mana_type":     req.GetManaType(),
		"mana_type_str": req.GetManaTypeStr(),
	}

	if err := s.gameMgr.SendPlayerAction(gameInstance.ID, player, "SEND_MANA_TYPE", payload); err != nil {
		return &pb.SendPlayerManaTypeResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.SendPlayerManaTypeResponse{Success: true}, nil
}

// SendPlayerAction forwards a high-level player action to the game engine.
func (s *mageServer) SendPlayerAction(ctx context.Context, req *pb.SendPlayerActionRequest) (*pb.SendPlayerActionResponse, error) {
	player, gameInstance, errMsg := s.resolveGamePlayer(req.GetSessionId(), req.GetGameId())
	if errMsg != "" {
		return &pb.SendPlayerActionResponse{Success: false, Error: errMsg}, nil
	}

	action := req.GetAction()
	if action == pb.PlayerAction_PLAYER_ACTION_UNSPECIFIED {
		return &pb.SendPlayerActionResponse{Success: false, Error: "action is required"}, nil
	}

	if err := s.gameMgr.SendPlayerAction(gameInstance.ID, player, "PLAYER_ACTION", action.String()); err != nil {
		return &pb.SendPlayerActionResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.SendPlayerActionResponse{Success: true}, nil
}

// SendSpecialAction handles special actions (play land, foretell, etc.)
func (s *mageServer) SendSpecialAction(ctx context.Context, req *pb.SendSpecialActionRequest) (*pb.SendSpecialActionResponse, error) {
	player, gameInstance, errMsg := s.resolveGamePlayer(req.GetSessionId(), req.GetGameId())
	if errMsg != "" {
		return &pb.SendSpecialActionResponse{Success: false, Error: errMsg}, nil
	}

	actionType := req.GetActionType()
	if actionType == pb.SpecialActionType_SPECIAL_ACTION_UNSPECIFIED {
		return &pb.SendSpecialActionResponse{Success: false, Error: "action_type is required"}, nil
	}

	sourceID := req.GetSourceId()
	// source_id is required for most actions, but not for ADVANCE_PHASE
	if sourceID == "" && actionType != pb.SpecialActionType_ADVANCE_PHASE {
		return &pb.SendSpecialActionResponse{Success: false, Error: "source_id is required"}, nil
	}

	payload := map[string]interface{}{
		"action_type": actionType.String(),
		"source_id":   sourceID,
	}

	if err := s.gameMgr.SendPlayerAction(gameInstance.ID, player, "SPECIAL_ACTION", payload); err != nil {
		return &pb.SendSpecialActionResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.SendSpecialActionResponse{Success: true}, nil
}

// ActivateAbility activates an activated ability on a permanent
func (s *mageServer) ActivateAbility(ctx context.Context, req *pb.ActivateAbilityRequest) (*pb.ActivateAbilityResponse, error) {
	player, gameInstance, errMsg := s.resolveGamePlayer(req.GetSessionId(), req.GetGameId())
	if errMsg != "" {
		return &pb.ActivateAbilityResponse{Success: false, Error: errMsg}, nil
	}

	cardID := req.GetCardId()
	if cardID == "" {
		return &pb.ActivateAbilityResponse{Success: false, Error: "card_id is required"}, nil
	}

	abilityID := req.GetAbilityId()
	if abilityID == "" {
		return &pb.ActivateAbilityResponse{Success: false, Error: "ability_id is required"}, nil
	}

	payload := map[string]interface{}{
		"card_id":    cardID,
		"ability_id": abilityID,
		"targets":    req.GetTargets(),
	}

	if err := s.gameMgr.SendPlayerAction(gameInstance.ID, player, "ACTIVATE_ABILITY", payload); err != nil {
		return &pb.ActivateAbilityResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.ActivateAbilityResponse{Success: true}, nil
}

// helper to resolve session/game/player for action RPCs
func (s *mageServer) resolveGamePlayer(sessionID, gameID string) (string, *game.Game, string) {
	sess, gameInstance, err := s.resolveGameAccess(sessionID, gameID, false)
	if err != nil {
		return "", nil, err.Error()
	}

	player := sess.GetUserID()
	if !gameInstance.IsPlayer(player) {
		return "", nil, "player not part of this game"
	}

	return player, gameInstance, ""
}

func (s *mageServer) resolveGameAccess(sessionID, gameID string, allowWatcher bool) (*session.Session, *game.Game, error) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil, nil, status.Errorf(codes.InvalidArgument, "session_id is required")
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return nil, nil, status.Errorf(codes.NotFound, "session not found")
	}

	user := sess.GetUserID()
	if user == "" {
		return nil, nil, status.Errorf(codes.InvalidArgument, "session not associated with a user")
	}

	gameID = strings.TrimSpace(gameID)
	if gameID == "" {
		return nil, nil, status.Errorf(codes.InvalidArgument, "game_id is required")
	}

	gameInstance, ok := s.gameMgr.GetGame(gameID)
	if !ok {
		return nil, nil, status.Errorf(codes.NotFound, "game not found")
	}

	if !gameInstance.IsPlayer(user) {
		watcherAllowed := false
		if allowWatcher {
			for _, watcher := range gameInstance.GetWatchers() {
				if watcher == user {
					watcherAllowed = true
					break
				}
			}
		}

		if !watcherAllowed && !sess.IsAdminSession() {
			return nil, nil, status.Errorf(codes.PermissionDenied, "user not part of this game")
		}
	}

	return sess, gameInstance, nil
}

func enginePlayersToProto(players []game.EnginePlayerView) []*pb.PlayerView {
	if len(players) == 0 {
		return nil
	}

	result := make([]*pb.PlayerView, 0, len(players))
	for _, p := range players {
		playerView := &pb.PlayerView{
			PlayerId:     p.PlayerID,
			Name:         p.Name,
			Life:         int32(p.Life),
			Poison:       int32(p.Poison),
			Energy:       int32(p.Energy),
			LibraryCount: int32(p.LibraryCount),
			HandCount:    int32(p.HandCount),
			Hand:         engineCardsToProto(p.Hand),
			Graveyard:    engineCardsToProto(p.Graveyard),
			ManaPool: &pb.ManaPoolView{
				White:     int32(p.ManaPool.White),
				Blue:      int32(p.ManaPool.Blue),
				Black:     int32(p.ManaPool.Black),
				Red:       int32(p.ManaPool.Red),
				Green:     int32(p.ManaPool.Green),
				Colorless: int32(p.ManaPool.Colorless),
			},
			HasPriority:         p.HasPriority,
			Passed:              p.Passed,
			StateOrdinal:        int32(p.StateOrdinal),
			Lost:                p.Lost,
			Left:                p.Left,
			Wins:                int32(p.Wins),
			KeptHand:            p.KeptHand,
			HasAvailableActions: p.HasAvailableActions,
		}
		result = append(result, playerView)
	}
	return result
}

func engineCardsToProto(cards []game.EngineCardView) []*pb.CardView {
	if len(cards) == 0 {
		return nil
	}

	result := make([]*pb.CardView, 0, len(cards))
	for _, card := range cards {
		cardView := &pb.CardView{
			Id:               card.ID,
			Name:             card.Name,
			DisplayName:      card.DisplayName,
			ManaCost:         card.ManaCost,
			Type:             card.Type,
			SubTypes:         strings.Join(card.SubTypes, " "),
			SuperTypes:       strings.Join(card.SuperTypes, " "),
			Color:            card.Color,
			Power:            card.Power,
			Toughness:        card.Toughness,
			Loyalty:          card.Loyalty,
			CardNumber:       int32(card.CardNumber),
			ExpansionSetCode: card.ExpansionSet,
			Rarity:           card.Rarity,
			RulesText:        card.RulesText,
			Tapped:           card.Tapped,
			Flipped:          card.Flipped,
			Transformed:      card.Transformed,
			FaceDown:         card.FaceDown,
			Zone:             int32(card.Zone),
			ControllerId:     card.ControllerID,
			OwnerId:          card.OwnerID,
			AttachedTo:       append([]string(nil), card.AttachedToCard...),
		}

		if len(card.Abilities) > 0 {
			abilities := make([]*pb.AbilityView, 0, len(card.Abilities))
			for _, ability := range card.Abilities {
				abilities = append(abilities, &pb.AbilityView{
					Id:   ability.ID,
					Text: ability.Text,
					Rule: ability.Rule,
				})
			}
			cardView.Abilities = abilities
		}

		if len(card.Counters) > 0 {
			counters := make([]*pb.CounterView, 0, len(card.Counters))
			for _, counter := range card.Counters {
				counters = append(counters, &pb.CounterView{
					Name:  counter.Name,
					Count: int32(counter.Count),
				})
			}
			cardView.Counters = counters
		}

		// Add summoning sickness status
		cardView.SummoningSickness = card.SummoningSickness

		// Add available actions (server source of truth)
		if len(card.AvailableActions) > 0 {
			cardView.AvailableActions = engineCardActionsToProto(card.AvailableActions)
		}

		result = append(result, cardView)
	}

	return result
}

func engineRevealedToProto(entries []game.EngineRevealedView) []*pb.RevealedView {
	if len(entries) == 0 {
		return nil
	}
	result := make([]*pb.RevealedView, 0, len(entries))
	for _, entry := range entries {
		result = append(result, &pb.RevealedView{
			Name:  entry.Name,
			Cards: engineCardsToProto(entry.Cards),
		})
	}
	return result
}

func engineLookedAtToProto(entries []game.EngineLookedAtView) []*pb.LookedAtView {
	if len(entries) == 0 {
		return nil
	}
	result := make([]*pb.LookedAtView, 0, len(entries))
	for _, entry := range entries {
		result = append(result, &pb.LookedAtView{
			Name:  entry.Name,
			Cards: engineCardsToProto(entry.Cards),
		})
	}
	return result
}

func engineCombatToProto(combat game.EngineCombatView) *pb.CombatView {
	if combat.AttackingPlayerID == "" && len(combat.Groups) == 0 {
		return nil
	}
	groups := make([]*pb.CombatGroupView, 0, len(combat.Groups))
	for _, group := range combat.Groups {
		groups = append(groups, &pb.CombatGroupView{
			Attackers:         append([]string(nil), group.Attackers...),
			Blockers:          append([]string(nil), group.Blockers...),
			DefendingPlayerId: group.DefendingPlayerID,
		})
	}
	return &pb.CombatView{
		AttackingPlayerId: combat.AttackingPlayerID,
		Groups:            groups,
	}
}

func engineMessagesToProto(messages []game.EngineMessage, startID int32) []*pb.GameMessage {
	if len(messages) == 0 {
		return nil
	}

	result := make([]*pb.GameMessage, 0, len(messages))
	nextID := startID
	for _, message := range messages {
		msg := &pb.GameMessage{
			Id:                nextID,
			Text:              message.Text,
			Color:             engineColorToString(message.Color),
			BookmarkId:        int32(message.BookmarkID),
			RollbackAvailable: message.RollbackAvailable,
		}
		if !message.Timestamp.IsZero() {
			msg.Time = timestamppb.New(message.Timestamp)
		}
		result = append(result, msg)
		nextID++
	}
	return result
}

func enginePromptsToMessages(prompts []game.EnginePrompt, startID int32) []*pb.GameMessage {
	if len(prompts) == 0 {
		return nil
	}

	result := make([]*pb.GameMessage, 0, len(prompts))
	nextID := startID
	for _, prompt := range prompts {
		text := prompt.Text
		if len(prompt.Options) > 0 {
			text = fmt.Sprintf("%s (options: %s)", prompt.Text, strings.Join(prompt.Options, ", "))
		}
		msg := &pb.GameMessage{
			Id:    nextID,
			Text:  fmt.Sprintf("Prompt for %s: %s", prompt.PlayerID, text),
			Color: "YELLOW",
		}
		if !prompt.Timestamp.IsZero() {
			msg.Time = timestamppb.New(prompt.Timestamp)
		}
		result = append(result, msg)
		nextID++
	}
	return result
}

func engineColorToString(color string) string {
	switch strings.ToLower(strings.TrimSpace(color)) {
	case "action":
		return "ORANGE"
	case "prompt":
		return "YELLOW"
	case "life":
		return "GREEN"
	case "mana":
		return "BLUE"
	case "status":
		return "BLACK"
	default:
		return "BLACK"
	}
}

func engineCardActionsToProto(actions []game.EngineCardAction) []*pb.CardAction {
	if len(actions) == 0 {
		return nil
	}

	result := make([]*pb.CardAction, 0, len(actions))
	for _, action := range actions {
		result = append(result, &pb.CardAction{
			ActionType:     stringToCardActionType(action.ActionType),
			ActionId:       action.ActionID,
			DisplayText:    action.DisplayText,
			IsEnabled:      action.IsEnabled,
			DisabledReason: action.DisabledReason,
		})
	}
	return result
}

func stringToCardActionType(actionType string) pb.CardActionType {
	switch actionType {
	case "CAST_SPELL":
		return pb.CardActionType_CARD_ACTION_CAST_SPELL
	case "PLAY_LAND":
		return pb.CardActionType_CARD_ACTION_PLAY_LAND
	case "ACTIVATE_ABILITY":
		return pb.CardActionType_CARD_ACTION_ACTIVATE_ABILITY
	case "ACTIVATE_MANA_ABILITY":
		return pb.CardActionType_CARD_ACTION_ACTIVATE_MANA_ABILITY
	default:
		return pb.CardActionType_CARD_ACTION_UNSPECIFIED
	}
}

// ==================== Replay Methods (Stubs) ====================

// ReplayInit initializes a replay session for a completed game
func (s *mageServer) ReplayInit(ctx context.Context, req *pb.ReplayInitRequest) (*pb.ReplayInitResponse, error) {
	// TODO: Implement replay functionality
	// For now, return not implemented
	s.logger.Debug("ReplayInit called (not yet implemented)",
		zap.String("session_id", req.GetSessionId()),
		zap.String("game_id", req.GetGameId()),
	)
	return &pb.ReplayInitResponse{
		Success: false,
		Error:   "Replay functionality not yet implemented",
	}, nil
}

// ReplayStart starts playing a replay
func (s *mageServer) ReplayStart(ctx context.Context, req *pb.ReplayStartRequest) (*pb.ReplayStartResponse, error) {
	// TODO: Implement replay functionality
	s.logger.Debug("ReplayStart called (not yet implemented)",
		zap.String("session_id", req.GetSessionId()),
		zap.String("game_id", req.GetGameId()),
	)
	return &pb.ReplayStartResponse{
		Success: false,
		Error:   "Replay functionality not yet implemented",
	}, nil
}

// ReplayStop stops playing a replay
func (s *mageServer) ReplayStop(ctx context.Context, req *pb.ReplayStopRequest) (*pb.ReplayStopResponse, error) {
	// TODO: Implement replay functionality
	s.logger.Debug("ReplayStop called (not yet implemented)",
		zap.String("session_id", req.GetSessionId()),
		zap.String("game_id", req.GetGameId()),
	)
	return &pb.ReplayStopResponse{
		Success: false,
		Error:   "Replay functionality not yet implemented",
	}, nil
}

// ReplayNext advances replay to next step
func (s *mageServer) ReplayNext(ctx context.Context, req *pb.ReplayNextRequest) (*pb.ReplayNextResponse, error) {
	// TODO: Implement replay functionality
	s.logger.Debug("ReplayNext called (not yet implemented)",
		zap.String("session_id", req.GetSessionId()),
		zap.String("game_id", req.GetGameId()),
	)
	return &pb.ReplayNextResponse{
		Success: false,
		Error:   "Replay functionality not yet implemented",
	}, nil
}

// ReplayPrevious goes back to previous step in replay
func (s *mageServer) ReplayPrevious(ctx context.Context, req *pb.ReplayPreviousRequest) (*pb.ReplayPreviousResponse, error) {
	// TODO: Implement replay functionality
	s.logger.Debug("ReplayPrevious called (not yet implemented)",
		zap.String("session_id", req.GetSessionId()),
		zap.String("game_id", req.GetGameId()),
	)
	return &pb.ReplayPreviousResponse{
		Success: false,
		Error:   "Replay functionality not yet implemented",
	}, nil
}

// ReplaySkipForward skips forward by specified number of steps
func (s *mageServer) ReplaySkipForward(ctx context.Context, req *pb.ReplaySkipForwardRequest) (*pb.ReplaySkipForwardResponse, error) {
	// TODO: Implement replay functionality
	s.logger.Debug("ReplaySkipForward called (not yet implemented)",
		zap.String("session_id", req.GetSessionId()),
		zap.String("game_id", req.GetGameId()),
		zap.Int32("steps", req.GetSteps()),
	)
	return &pb.ReplaySkipForwardResponse{
		Success: false,
		Error:   "Replay functionality not yet implemented",
	}, nil
}

// ==================== Active Game Recovery ====================

// GetMyActiveGames returns all active games the current user is participating in
// This enables reconnection after disconnection or server restart
func (s *mageServer) GetMyActiveGames(ctx context.Context, req *pb.GetMyActiveGamesRequest) (*pb.GetMyActiveGamesResponse, error) {
	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return &pb.GetMyActiveGamesResponse{Games: nil}, nil
	}

	sess, ok := s.sessionMgr.GetSession(sessionID)
	if !ok {
		return &pb.GetMyActiveGamesResponse{Games: nil}, nil
	}

	username := sess.GetUserID()
	if username == "" {
		return &pb.GetMyActiveGamesResponse{Games: nil}, nil
	}

	// Query active games from database
	if s.activeGameRepo == nil {
		s.logger.Warn("GetMyActiveGames called but activeGameRepo is nil")
		return &pb.GetMyActiveGamesResponse{Games: nil}, nil
	}

	activeGames, err := s.activeGameRepo.GetActiveGamesForPlayer(ctx, username)
	if err != nil {
		s.logger.Error("failed to get active games for player",
			zap.String("username", username),
			zap.Error(err),
		)
		return &pb.GetMyActiveGamesResponse{Games: nil}, nil
	}

	// Convert to proto format
	games := make([]*pb.ActiveGameInfo, 0, len(activeGames))
	for _, ag := range activeGames {
		gameInfo := &pb.ActiveGameInfo{
			GameId:     ag.GameID,
			TableId:    ag.TableID,
			GameType:   ag.GameType,
			Players:    ag.Players,
			TurnNumber: int32(ag.TurnNumber),
			State:      ag.State,
			CreatedAt:  ag.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:  ag.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
		games = append(games, gameInfo)
	}

	s.logger.Info("returned active games for player",
		zap.String("username", username),
		zap.Int("count", len(games)),
	)

	return &pb.GetMyActiveGamesResponse{Games: games}, nil
}

// ==================== Rollback Methods ====================

// RequestRollback initiates a rollback request to a specific message in the game log.
// This requires opponent consent in multiplayer games.
func (s *mageServer) RequestRollback(ctx context.Context, req *pb.RequestRollbackRequest) (*pb.RequestRollbackResponse, error) {
	player, gameInstance, errMsg := s.resolveGamePlayer(req.GetSessionId(), req.GetGameId())
	if errMsg != "" {
		return &pb.RequestRollbackResponse{Success: false, Error: errMsg}, nil
	}

	messageID := req.GetMessageId()
	if messageID <= 0 {
		return &pb.RequestRollbackResponse{Success: false, Error: "message_id is required"}, nil
	}

	// Get the underlying MageEngine from the adapter
	mageEngine := s.gameAdapter.GetMageEngine()
	if mageEngine == nil {
		return &pb.RequestRollbackResponse{Success: false, Error: "game engine not available"}, nil
	}

	requestID, err := mageEngine.RequestRollback(gameInstance.ID, player, int(messageID))
	if err != nil {
		s.logger.Warn("rollback request failed",
			zap.String("game_id", gameInstance.ID),
			zap.String("player", player),
			zap.Int32("message_id", messageID),
			zap.Error(err),
		)
		return &pb.RequestRollbackResponse{Success: false, Error: err.Error()}, nil
	}

	s.logger.Info("rollback request initiated",
		zap.String("game_id", gameInstance.ID),
		zap.String("player", player),
		zap.Int32("message_id", messageID),
		zap.String("request_id", requestID),
	)

	return &pb.RequestRollbackResponse{
		Success:   true,
		RequestId: requestID,
	}, nil
}

// RespondToRollback handles a player's response to a pending rollback request.
func (s *mageServer) RespondToRollback(ctx context.Context, req *pb.RespondToRollbackRequest) (*pb.RespondToRollbackResponse, error) {
	player, gameInstance, errMsg := s.resolveGamePlayer(req.GetSessionId(), req.GetGameId())
	if errMsg != "" {
		return &pb.RespondToRollbackResponse{Success: false, Error: errMsg}, nil
	}

	requestID := req.GetRequestId()
	if requestID == "" {
		return &pb.RespondToRollbackResponse{Success: false, Error: "request_id is required"}, nil
	}

	// Get the underlying MageEngine from the adapter
	mageEngine := s.gameAdapter.GetMageEngine()
	if mageEngine == nil {
		return &pb.RespondToRollbackResponse{Success: false, Error: "game engine not available"}, nil
	}

	err := mageEngine.RespondToRollback(gameInstance.ID, player, requestID, req.GetApproved())
	if err != nil {
		s.logger.Warn("rollback response failed",
			zap.String("game_id", gameInstance.ID),
			zap.String("player", player),
			zap.String("request_id", requestID),
			zap.Bool("approved", req.GetApproved()),
			zap.Error(err),
		)
		return &pb.RespondToRollbackResponse{Success: false, Error: err.Error()}, nil
	}

	s.logger.Info("rollback response processed",
		zap.String("game_id", gameInstance.ID),
		zap.String("player", player),
		zap.String("request_id", requestID),
		zap.Bool("approved", req.GetApproved()),
	)

	return &pb.RespondToRollbackResponse{Success: true}, nil
}

// CancelRollback cancels a pending rollback request.
func (s *mageServer) CancelRollback(ctx context.Context, req *pb.CancelRollbackRequest) (*pb.CancelRollbackResponse, error) {
	_, gameInstance, errMsg := s.resolveGamePlayer(req.GetSessionId(), req.GetGameId())
	if errMsg != "" {
		return &pb.CancelRollbackResponse{Success: false, Error: errMsg}, nil
	}

	// Get the underlying MageEngine from the adapter
	mageEngine := s.gameAdapter.GetMageEngine()
	if mageEngine == nil {
		return &pb.CancelRollbackResponse{Success: false, Error: "game engine not available"}, nil
	}

	err := mageEngine.CancelRollbackRequest(gameInstance.ID)
	if err != nil {
		return &pb.CancelRollbackResponse{Success: false, Error: err.Error()}, nil
	}

	s.logger.Info("rollback request cancelled",
		zap.String("game_id", gameInstance.ID),
	)

	return &pb.CancelRollbackResponse{Success: true}, nil
}
