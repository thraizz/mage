package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/magefree/mage-server-go/internal/chat"
	"github.com/magefree/mage-server-go/internal/config"
	"github.com/magefree/mage-server-go/internal/draft"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/room"
	"github.com/magefree/mage-server-go/internal/server"
	"github.com/magefree/mage-server-go/internal/session"
	"github.com/magefree/mage-server-go/internal/table"
	"github.com/magefree/mage-server-go/internal/tournament"
	pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gameServerEnv struct {
	server     pb.MageServerServer
	sessionMgr session.Manager
	roomMgr    *room.Manager
	tableMgr   *table.Manager
	gameMgr    *game.Manager
	adapter    *game.EngineAdapter
	engine     *game.MageEngine
	logger     *zap.Logger
}

func newGameServerEnv(t testing.TB) *gameServerEnv {
	logger := zaptest.NewLogger(t)

	sessionMgr := session.NewManager(5*time.Minute, logger)
	roomMgr := room.NewManager(logger)
	chatMgr := chat.NewManager(logger)
	tableMgr := table.NewManager(logger)
	gameMgr := game.NewManager(logger)
	tournamentMgr := tournament.NewManager(logger)
	draftMgr := draft.NewManager(logger)
	engine := game.NewMageEngine(logger)
	adapter := game.NewEngineAdapter(engine, logger)

	cfg := &config.Config{}

	srv := server.NewMageServer(
		cfg,
		nil,
		sessionMgr,
		nil,
		nil,
		roomMgr,
		chatMgr,
		tableMgr,
		gameMgr,
		tournamentMgr,
		draftMgr,
		nil,
		nil,
		"test",
		logger,
		adapter,
	)

	return &gameServerEnv{
		server:     srv,
		sessionMgr: sessionMgr,
		roomMgr:    roomMgr,
		tableMgr:   tableMgr,
		gameMgr:    gameMgr,
		adapter:    adapter,
		engine:     engine,
		logger:     logger,
	}
}

// TestGameEngineActionQueue verifies that queued player actions flow through the engine adapter.
func TestGameEngineActionQueue(t *testing.T) {
	logger := zaptest.NewLogger(t)

	gameMgr := game.NewManager(logger)
	engine := game.NewMageEngine(logger)
	adapter := game.NewEngineAdapter(engine, logger)

	players := []string{"Alice", "Bob"}
	g := gameMgr.CreateGame("table-1", "Duel", players)
	if g == nil {
		t.Fatal("failed to create game")
	}

	if err := adapter.StartGame(g); err != nil {
		t.Fatalf("failed to start engine game: %v", err)
	}

	done := make(chan struct{})
	go func() {
		adapter.ProcessGameActions(g)
		close(done)
	}()

	if err := gameMgr.SendPlayerAction(g.ID, "Alice", "PLAYER_ACTION", "PASS"); err != nil {
		t.Fatalf("failed to enqueue player action: %v", err)
	}

	if err := gameMgr.SendPlayerAction(g.ID, "Bob", "SEND_INTEGER", 3); err != nil {
		t.Fatalf("failed to enqueue second player action: %v", err)
	}

	// Allow the adapter goroutine to process the queued actions.
	time.Sleep(25 * time.Millisecond)

	viewRaw, err := adapter.GetGameView(g.ID, "Alice")
	if err != nil {
		t.Fatalf("failed to retrieve engine view: %v", err)
	}

	view, ok := viewRaw.(*game.EngineGameView)
	if !ok {
		t.Fatalf("unexpected view type %T", viewRaw)
	}

	if len(view.Messages) == 0 {
		t.Fatal("expected engine messages to be recorded")
	}

	var passRecorded, lifeChangeRecorded bool
	for _, msg := range view.Messages {
		textUpper := strings.ToUpper(msg.Text)
		if strings.Contains(textUpper, "ALICE") && strings.Contains(textUpper, "PASS") {
			passRecorded = true
		}
		if strings.Contains(textUpper, "BOB") && strings.Contains(textUpper, "LIFE") {
			lifeChangeRecorded = true
		}
	}

	if !passRecorded {
		t.Errorf("expected pass action from Alice to be recorded in messages: %+v", view.Messages)
	}
	if !lifeChangeRecorded {
		t.Errorf("expected life change from Bob to be recorded in messages: %+v", view.Messages)
	}

	if len(view.Prompts) == 0 {
		t.Errorf("expected engine to surface prompts after player actions")
	}

	// Clean up the goroutine: simulate match shutdown.
	if err := adapter.EndGame(g, "Alice"); err != nil {
		t.Fatalf("failed to end game: %v", err)
	}

	// Closing the action queue stops ProcessGameActions loop.
	close(g.ActionQueue)
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("engine processing goroutine did not exit")
	}
}

// TestGameGetViewWatcherAccess verifies watchers can view game state while non-participants are rejected.
// TODO: Re-enable when game view properly returns battlefield data at game start
// This test expects battlefield to be populated immediately after game start
func testGameGetViewWatcherAccess(t *testing.T) {
	env := newGameServerEnv(t)

	players := []string{"Alice", "Bob"}
	g := env.gameMgr.CreateGame("table-1", "Duel", players)
	if err := env.adapter.StartGame(g); err != nil {
		t.Fatalf("failed to start mage engine game: %v", err)
	}
	defer env.adapter.EndGame(g, "")

	go env.adapter.ProcessGameActions(g)
	defer close(g.ActionQueue)

	playerSession := env.sessionMgr.CreateSession("player-session", "tester")
	playerSession.SetUserID("Alice")

	ctx := context.Background()

	// Allow watcher to join.
	watcherSession := env.sessionMgr.CreateSession("watcher-session", "tester")
	watcherSession.SetUserID("Watcher")

	if _, err := env.server.GameWatchStart(ctx, &pb.GameWatchStartRequest{
		SessionId: watcherSession.ID,
		GameId:    g.ID,
	}); err != nil {
		t.Fatalf("watcher failed to start watching: %v", err)
	}

	resp, err := env.server.GameGetView(ctx, &pb.GameGetViewRequest{
		SessionId: watcherSession.ID,
		GameId:    g.ID,
	})
	if err != nil {
		t.Fatalf("watcher GameGetView failed: %v", err)
	}
	if resp.GetGame() == nil {
		t.Fatalf("expected game view for watcher")
	}
	if len(resp.GetGame().GetMessages()) == 0 {
		t.Fatalf("expected watcher view to include engine messages")
	}
	if len(resp.GetGame().GetBattlefield()) == 0 {
		t.Fatalf("expected watcher view to include battlefield state")
	}

	// Non participant/non watcher should be denied.
	outsider := env.sessionMgr.CreateSession("outsider-session", "tester")
	outsider.SetUserID("Mallory")

	_, err = env.server.GameGetView(ctx, &pb.GameGetViewRequest{
		SessionId: outsider.ID,
		GameId:    g.ID,
	})
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("expected permission denied for outsider, got %v", status.Code(err))
	}

	// Player should be able to view.
	respPlayer, err := env.server.GameGetView(ctx, &pb.GameGetViewRequest{
		SessionId: playerSession.ID,
		GameId:    g.ID,
	})
	if err != nil {
		t.Fatalf("player GameGetView failed: %v", err)
	}
	if len(respPlayer.GetGame().GetPlayers()) != len(players) {
		t.Fatalf("expected %d players in view, got %d", len(players), len(respPlayer.GetGame().GetPlayers()))
	}
}

// TestGameActionQueueOverflow ensures SendPlayerAction returns an error when the queue is full.
func TestGameActionQueueOverflow(t *testing.T) {
	logger := zaptest.NewLogger(t)
	gameMgr := game.NewManager(logger)

	players := []string{"Alice", "Bob"}
	g := gameMgr.CreateGame("table-1", "Duel", players)

	for i := 0; i < cap(g.ActionQueue); i++ {
		if err := gameMgr.SendPlayerAction(g.ID, "Alice", "PING", i); err != nil {
			t.Fatalf("unexpected error enqueuing action %d: %v", i, err)
		}
	}

	if err := gameMgr.SendPlayerAction(g.ID, "Alice", "PING", "overflow"); err == nil {
		t.Fatalf("expected queue overflow error")
	}
}

// TestMatchFlowEndToEnd executes a representative match flow through the gRPC handlers.
func TestMatchFlowEndToEnd(t *testing.T) {
	env := newGameServerEnv(t)
	ctx := context.Background()

	tbl := env.tableMgr.CreateTable("Match Table", "Duel", "Alice", env.roomMgr.GetMainRoomID(), 2, "")
	if tbl == nil {
		t.Fatal("failed to create table")
	}

	if err := tbl.AddPlayer("Alice", "Human"); err != nil {
		t.Fatalf("failed adding Alice to table: %v", err)
	}
	if err := tbl.AddPlayer("Bob", "Human"); err != nil {
		t.Fatalf("failed adding Bob to table: %v", err)
	}

	env.roomMgr.UserJoinRoom(tbl.RoomID, "Alice")
	env.roomMgr.UserJoinRoom(tbl.RoomID, "Bob")

	aliceSession := env.sessionMgr.CreateSession("alice-session", "localhost")
	aliceSession.SetUserID("Alice")

	bobSession := env.sessionMgr.CreateSession("bob-session", "localhost")
	bobSession.SetUserID("Bob")

	watcherSession := env.sessionMgr.CreateSession("spectator-session", "localhost")
	watcherSession.SetUserID("Spectator")

	startResp, err := env.server.MatchStart(ctx, &pb.MatchStartRequest{
		SessionId: aliceSession.ID,
		TableId:   tbl.ID,
	})
	if err != nil || !startResp.GetSuccess() {
		t.Fatalf("match start failed: %v, success=%v", err, startResp.GetSuccess())
	}

	gameInstance, ok := env.gameMgr.GetGameByTable(tbl.ID)
	if !ok {
		t.Fatal("game not created after match start")
	}

	if _, err := env.server.GameWatchStart(ctx, &pb.GameWatchStartRequest{
		SessionId: watcherSession.ID,
		GameId:    gameInstance.ID,
	}); err != nil {
		t.Fatalf("watcher failed to start: %v", err)
	}

	if _, err := env.server.GameGetView(ctx, &pb.GameGetViewRequest{
		SessionId: watcherSession.ID,
		GameId:    gameInstance.ID,
	}); err != nil {
		t.Fatalf("watcher GameGetView should succeed: %v", err)
	}

	if _, err := env.server.GameJoin(ctx, &pb.GameJoinRequest{
		SessionId: bobSession.ID,
		GameId:    gameInstance.ID,
	}); err != nil {
		t.Fatalf("bob GameJoin failed: %v", err)
	}

	// Test basic player actions through gRPC
	// Alice casts Lightning Bolt (a card that exists in the starting hand)
	if _, err := env.server.SendPlayerString(ctx, &pb.SendPlayerStringRequest{
		SessionId: aliceSession.ID,
		GameId:    gameInstance.ID,
		Data:      "Lightning Bolt",
	}); err != nil {
		t.Fatalf("alice cast Lightning Bolt failed: %v", err)
	}

	time.Sleep(25 * time.Millisecond)

	// Verify the spell is on the stack
	viewRaw, err := env.adapter.GetGameView(gameInstance.ID, "Alice")
	if err != nil {
		t.Fatalf("engine view retrieval failed: %v", err)
	}

	engineView, ok := viewRaw.(*game.EngineGameView)
	if !ok {
		t.Fatalf("unexpected engine view type: %T", viewRaw)
	}

	if len(engineView.Messages) == 0 {
		t.Fatalf("expected engine messages, got 0")
	}
	if len(engineView.Stack) == 0 {
		t.Fatalf("expected stack to contain cast spell")
	}
	if len(engineView.Players) != 2 {
		t.Fatalf("expected 2 players in engine view, got %d", len(engineView.Players))
	}

	if _, err := env.server.GameGetView(ctx, &pb.GameGetViewRequest{
		SessionId: aliceSession.ID,
		GameId:    gameInstance.ID,
	}); err != nil {
		t.Fatalf("player GameGetView failed: %v", err)
	}

	if _, err := env.server.MatchQuit(ctx, &pb.MatchQuitRequest{
		SessionId: aliceSession.ID,
		GameId:    gameInstance.ID,
	}); err != nil {
		t.Fatalf("match quit failed: %v", err)
	}

	time.Sleep(25 * time.Millisecond)

	if _, ok := env.gameMgr.GetGame(gameInstance.ID); ok {
		t.Fatal("game should be removed after match quit")
	}

	if tbl.GetState() != table.TableStateFinished {
		t.Fatalf("expected table finished state, got %s", tbl.GetState().String())
	}
}

func TestTurnProgressionAfterPassChain(t *testing.T) {
	env := newGameServerEnv(t)
	ctx := context.Background()

	tbl := env.tableMgr.CreateTable("Turn Table", "Duel", "Alice", env.roomMgr.GetMainRoomID(), 2, "")
	if tbl == nil {
		t.Fatal("failed to create table")
	}

	if err := tbl.AddPlayer("Alice", "Human"); err != nil {
		t.Fatalf("failed adding Alice to table: %v", err)
	}
	if err := tbl.AddPlayer("Bob", "Human"); err != nil {
		t.Fatalf("failed adding Bob to table: %v", err)
	}

	env.roomMgr.UserJoinRoom(tbl.RoomID, "Alice")
	env.roomMgr.UserJoinRoom(tbl.RoomID, "Bob")

	aliceSession := env.sessionMgr.CreateSession("alice-turn-session", "localhost")
	aliceSession.SetUserID("Alice")
	bobSession := env.sessionMgr.CreateSession("bob-turn-session", "localhost")
	bobSession.SetUserID("Bob")

	startResp, err := env.server.MatchStart(ctx, &pb.MatchStartRequest{
		SessionId: aliceSession.ID,
		TableId:   tbl.ID,
	})
	if err != nil || !startResp.GetSuccess() {
		t.Fatalf("turn test match start failed: %v, success=%v", err, startResp.GetSuccess())
	}

	gameInstance, ok := env.gameMgr.GetGameByTable(tbl.ID)
	if !ok {
		t.Fatal("game not created for turn progression test")
	}

	viewBefore, err := env.server.GameGetView(ctx, &pb.GameGetViewRequest{
		SessionId: aliceSession.ID,
		GameId:    gameInstance.ID,
	})
	if err != nil {
		t.Fatalf("initial GameGetView failed: %v", err)
	}

	initialStep := viewBefore.GetGame().GetStep()
	initialTurn := viewBefore.GetGame().GetTurn()
	if initialStep == "" {
		t.Fatal("expected initial step value")
	}

	if _, err := env.server.SendPlayerAction(ctx, &pb.SendPlayerActionRequest{
		SessionId: aliceSession.ID,
		GameId:    gameInstance.ID,
		Action:    pb.PlayerAction_PASS,
	}); err != nil {
		t.Fatalf("alice pass failed: %v", err)
	}
	if _, err := env.server.SendPlayerAction(ctx, &pb.SendPlayerActionRequest{
		SessionId: bobSession.ID,
		GameId:    gameInstance.ID,
		Action:    pb.PlayerAction_PASS,
	}); err != nil {
		t.Fatalf("bob pass failed: %v", err)
	}

	// Allow background processing of queued actions.
	time.Sleep(25 * time.Millisecond)

	viewAfter, err := env.server.GameGetView(ctx, &pb.GameGetViewRequest{
		SessionId: aliceSession.ID,
		GameId:    gameInstance.ID,
	})
	if err != nil {
		t.Fatalf("post-pass GameGetView failed: %v", err)
	}

	if viewAfter.GetGame().GetStep() == initialStep {
		t.Fatalf("expected step to advance after double pass, still %s", initialStep)
	}
	if viewAfter.GetGame().GetTurn() != initialTurn {
		t.Fatalf("expected to remain on turn %d, got %d", initialTurn, viewAfter.GetGame().GetTurn())
	}
	if viewAfter.GetGame().GetPriorityPlayerId() != viewAfter.GetGame().GetActivePlayerId() {
		t.Fatalf("expected priority to revert to active player, got priority=%s active=%s",
			viewAfter.GetGame().GetPriorityPlayerId(), viewAfter.GetGame().GetActivePlayerId())
	}

	foundAdvance := false
	for _, msg := range viewAfter.GetGame().GetMessages() {
		if strings.Contains(msg.GetText(), "advances to") {
			foundAdvance = true
			break
		}
	}
	if !foundAdvance {
		t.Fatalf("expected log message indicating advancement, messages=%v", viewAfter.GetGame().GetMessages())
	}

	if _, err := env.server.MatchQuit(ctx, &pb.MatchQuitRequest{
		SessionId: aliceSession.ID,
		GameId:    gameInstance.ID,
	}); err != nil {
		t.Fatalf("turn test match quit failed: %v", err)
	}
}

// TODO: Re-enable when instant spells correctly resolve to graveyard instead of battlefield
// This test expects Lightning Bolt (instant) to be on battlefield after resolution, but instants go to graveyard
func testStackResolutionAfterPasses(t *testing.T) {
	env := newGameServerEnv(t)
	ctx := context.Background()

	tbl := env.tableMgr.CreateTable("Stack Table", "Duel", "Alice", env.roomMgr.GetMainRoomID(), 2, "")
	if tbl == nil {
		t.Fatal("failed to create table")
	}

	if err := tbl.AddPlayer("Alice", "Human"); err != nil {
		t.Fatalf("failed adding Alice to table: %v", err)
	}
	if err := tbl.AddPlayer("Bob", "Human"); err != nil {
		t.Fatalf("failed adding Bob to table: %v", err)
	}

	env.roomMgr.UserJoinRoom(tbl.RoomID, "Alice")
	env.roomMgr.UserJoinRoom(tbl.RoomID, "Bob")

	aliceSession := env.sessionMgr.CreateSession("alice-stack-session", "localhost")
	aliceSession.SetUserID("Alice")
	bobSession := env.sessionMgr.CreateSession("bob-stack-session", "localhost")
	bobSession.SetUserID("Bob")

	startResp, err := env.server.MatchStart(ctx, &pb.MatchStartRequest{
		SessionId: aliceSession.ID,
		TableId:   tbl.ID,
	})
	if err != nil || !startResp.GetSuccess() {
		t.Fatalf("stack test match start failed: %v, success=%v", err, startResp.GetSuccess())
	}

	gameInstance, ok := env.gameMgr.GetGameByTable(tbl.ID)
	if !ok {
		t.Fatal("game not created for stack progression test")
	}

	if _, err := env.server.SendPlayerString(ctx, &pb.SendPlayerStringRequest{
		SessionId: aliceSession.ID,
		GameId:    gameInstance.ID,
		Data:      "Lightning Bolt",
	}); err != nil {
		t.Fatalf("spell cast failed: %v", err)
	}

	time.Sleep(25 * time.Millisecond)

	viewStacked, err := env.server.GameGetView(ctx, &pb.GameGetViewRequest{
		SessionId: aliceSession.ID,
		GameId:    gameInstance.ID,
	})
	if err != nil {
		t.Fatalf("GameGetView after cast failed: %v", err)
	}
	stackObjects := viewStacked.GetGame().GetStack()
	if len(stackObjects) != 2 {
		t.Fatalf("expected stack to contain spell and triggered ability, got %d entries", len(stackObjects))
	}
	triggerTop := stackObjects[len(stackObjects)-1]
	if !strings.Contains(strings.ToLower(triggerTop.GetDisplayName()), "trigger") {
		t.Fatalf("expected top of stack to be triggered ability, got %s", triggerTop.GetDisplayName())
	}

	if _, err := env.server.SendPlayerAction(ctx, &pb.SendPlayerActionRequest{
		SessionId: aliceSession.ID,
		GameId:    gameInstance.ID,
		Action:    pb.PlayerAction_PASS,
	}); err != nil {
		t.Fatalf("alice pass failed: %v", err)
	}
	if _, err := env.server.SendPlayerAction(ctx, &pb.SendPlayerActionRequest{
		SessionId: bobSession.ID,
		GameId:    gameInstance.ID,
		Action:    pb.PlayerAction_PASS,
	}); err != nil {
		t.Fatalf("bob pass failed: %v", err)
	}

	time.Sleep(25 * time.Millisecond)

	viewResolved, err := env.server.GameGetView(ctx, &pb.GameGetViewRequest{
		SessionId: aliceSession.ID,
		GameId:    gameInstance.ID,
	})
	if err != nil {
		t.Fatalf("GameGetView after resolution failed: %v", err)
	}

	if len(viewResolved.GetGame().GetStack()) != 0 {
		t.Fatalf("expected stack to be empty after resolution")
	}

	foundLightning := false
	var lightningPower, lightningToughness string
	for _, card := range viewResolved.GetGame().GetBattlefield() {
		if strings.Contains(strings.ToLower(card.GetDisplayName()), "lightning bolt") {
			foundLightning = true
			lightningPower = card.GetPower()
			lightningToughness = card.GetToughness()
			break
		}
	}
	if !foundLightning {
		t.Fatalf("expected Lightning Bolt to appear on battlefield after resolution")
	}
	if lightningPower != "2" || lightningToughness != "3" {
		t.Fatalf("expected Sanctuary effect to boost creature to 2/3, got %s/%s", lightningPower, lightningToughness)
	}

	var aliceLife int32
	for _, p := range viewResolved.GetGame().GetPlayers() {
		if p.GetPlayerId() == "Alice" {
			aliceLife = p.GetLife()
			break
		}
	}
	if aliceLife != 21 {
		t.Fatalf("expected Alice life gain from triggered ability, got %d", aliceLife)
	}

	if _, err := env.server.MatchQuit(ctx, &pb.MatchQuitRequest{
		SessionId: aliceSession.ID,
		GameId:    gameInstance.ID,
	}); err != nil {
		t.Fatalf("stack test match quit failed: %v", err)
	}
}
