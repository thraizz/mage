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

	players := make([]string, 0, len(tbl.Seats))
	for _, seat := range tbl.Seats {
		if seat.PlayerName != "" {
			players = append(players, seat.PlayerName)
			if !seat.DeckValid {
				s.logger.Debug("player starting without submitted deck",
					zap.String("table_id", tbl.ID),
					zap.String("player", seat.PlayerName),
				)
			}
		}
	}

	if len(players) < 2 {
		return &pb.MatchStartResponse{Success: false, Error: "not enough players to start match"}, nil
	}

	game := s.gameMgr.CreateGame(tbl.ID, tbl.GameType, players)
	tbl.RecordMatch(game.ID)
	tbl.SetState(table.TableStateDueling)

	if s.gameAdapter != nil {
		if err := s.gameAdapter.StartGame(game); err != nil {
			s.logger.Warn("failed to start game engine",
				zap.String("game_id", game.ID),
				zap.Error(err),
			)
		}
		go s.gameAdapter.ProcessGameActions(game)
	}

	s.logger.Info("match started",
		zap.String("table_id", tbl.ID),
		zap.String("game_id", game.ID),
		zap.Strings("players", players),
	)

	return &pb.MatchStartResponse{Success: true}, nil
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
	sess, gameInstance, err := s.resolveGameAccess(req.GetSessionId(), req.GetGameId(), true)
	if err != nil {
		return nil, err
	}

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
		if engineView, engineErr := s.gameAdapter.GetGameView(gameInstance.ID, playerID); engineErr == nil && engineView != nil {
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
			HasPriority:  p.HasPriority,
			Passed:       p.Passed,
			StateOrdinal: int32(p.StateOrdinal),
			Lost:         p.Lost,
			Left:         p.Left,
			Wins:         int32(p.Wins),
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
			Id:    nextID,
			Text:  message.Text,
			Color: engineColorToString(message.Color),
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
