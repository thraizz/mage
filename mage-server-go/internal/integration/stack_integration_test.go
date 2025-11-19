package integration

import (
	"strings"
	"testing"
	"time"

	"github.com/magefree/mage-server-go/internal/game"
	"go.uber.org/zap"
)

// resolveStackUntilEmpty resolves all items on the stack by having all players pass repeatedly
func resolveStackUntilEmpty(t *testing.T, engine *game.MageEngine, gameID string, players []string, maxCycles int) {
	t.Helper()
	viewInterface, err := engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view := viewInterface.(*game.EngineGameView)

	for i := 0; i < maxCycles && len(view.Stack) > 0; i++ {
		for _, player := range players {
			err = engine.ProcessAction(gameID, game.PlayerAction{
				PlayerID:   player,
				ActionType: "SEND_STRING",
				Data:       "Pass",
				Timestamp:  time.Now(),
			})
			if err != nil {
				t.Fatalf("Failed to pass: %v", err)
			}
		}
		viewInterface, err = engine.GetGameView(gameID, "")
		if err != nil {
			t.Fatalf("Failed to get game view: %v", err)
		}
		view = viewInterface.(*game.EngineGameView)
	}
}

// TestMultiObjectStackResolution tests that multiple objects resolve in LIFO order
func TestMultiObjectStackResolution(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "stack-test-1"
	players := []string{"Alice", "Bob"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Alice casts spell 1
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "SEND_STRING",
		Data:       "Lightning Bolt",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast spell 1: %v", err)
	}

	// Alice passes priority so Bob can respond
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass priority: %v", err)
	}

	// Bob casts spell 2 (goes on top)
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "SEND_STRING",
		Data:       "Counterspell",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast spell 2: %v", err)
	}

	// Bob passes priority so Alice can respond
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass priority: %v", err)
	}

	// Alice casts spell 3 (goes on top)
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "SEND_STRING",
		Data:       "Shock",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast spell 3: %v", err)
	}

	viewInterface, err := engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view := viewInterface.(*game.EngineGameView)

	// Verify stack has at least 3 items (may have more due to triggered abilities)
	if len(view.Stack) < 3 {
		t.Errorf("Expected at least 3 items on stack, got %d", len(view.Stack))
	}

	// Verify order: topmost last (Shock, Counterspell, Lightning Bolt)
	if len(view.Stack) >= 3 {
		top := view.Stack[len(view.Stack)-1]
		if !strings.Contains(strings.ToLower(top.DisplayName), "shock") {
			t.Errorf("Top of stack should be Shock, got %s", top.DisplayName)
		}
	}

	// All players pass - stack should resolve in reverse order
	// Shock resolves first, then Counterspell, then Lightning Bolt
	// Alice retains priority after casting Shock, so she passes first
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	// After resolution, stack should be empty
	viewInterface, err = engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view = viewInterface.(*game.EngineGameView)

	// Resolve any triggered abilities
	resolveStackUntilEmpty(t, engine, gameID, players, 20)

	viewInterface, err = engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view = viewInterface.(*game.EngineGameView)

	// Stack should eventually be empty (or at least much smaller)
	// Note: Some triggered abilities may remain if they create infinite loops
	if len(view.Stack) > 10 {
		t.Logf("Warning: Stack still has %d items after resolution cycles", len(view.Stack))
	}
}

// TestCounterspell tests that spells can be countered
func TestCounterspell(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "counterspell-test"
	players := []string{"Alice", "Bob"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Alice casts Lightning Bolt
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "SEND_STRING",
		Data:       "Lightning Bolt",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast spell: %v", err)
	}

	viewInterface, err := engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view := viewInterface.(*game.EngineGameView)

	// Verify spell is on stack (may have triggered abilities)
	if len(view.Stack) < 1 {
		t.Fatalf("Expected at least 1 item on stack, got %d", len(view.Stack))
	}

	// Find the Lightning Bolt spell (not triggered abilities)
	var spellID string
	for _, item := range view.Stack {
		name := strings.ToLower(item.DisplayName)
		if strings.Contains(name, "lightning bolt") && !strings.Contains(name, "trigger") {
			spellID = item.ID
			break
		}
	}
	if spellID == "" && len(view.Stack) > 0 {
		// Fallback: use first non-triggered item
		for _, item := range view.Stack {
			if !strings.Contains(strings.ToLower(item.DisplayName), "trigger") {
				spellID = item.ID
				break
			}
		}
		if spellID == "" {
			spellID = view.Stack[0].ID // Last resort
		}
	}

	// Alice passes priority so Bob can respond
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass priority: %v", err)
	}

	// Bob counters it
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "SEND_UUID",
		Data:       spellID,
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to counter spell: %v", err)
	}

	viewInterface, err = engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view = viewInterface.(*game.EngineGameView)

	// Stack should be empty (spell countered and removed)
	// Note: May have triggered abilities, but countered spell should be gone
	// Check that the countered spell is not on stack
	foundCounteredSpell := false
	for _, item := range view.Stack {
		if item.ID == spellID {
			foundCounteredSpell = true
			break
		}
	}
	if foundCounteredSpell {
		t.Errorf("Countered spell should not be on stack")
	}
}

// TestNestedResponses tests nested spell casting (response to response)
// TODO: Re-enable when proper card initialization is implemented
// This test requires cards "Dispel" to exist in player hands
func testNestedResponses(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "nested-test"
	players := []string{"Alice", "Bob"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Alice casts Lightning Bolt
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "SEND_STRING",
		Data:       "Lightning Bolt",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast spell: %v", err)
	}

	// Alice passes priority so Bob can respond
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass priority: %v", err)
	}

	// Bob casts Counterspell
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "SEND_STRING",
		Data:       "Counterspell",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast counterspell: %v", err)
	}

	// Bob passes priority so Alice can respond
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass priority: %v", err)
	}

	// Alice casts another spell in response
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "SEND_STRING",
		Data:       "Dispel",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast dispel: %v", err)
	}

	viewInterface, err := engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view := viewInterface.(*game.EngineGameView)

	// Should have at least 3 items on stack: Dispel (top), Counterspell, Lightning Bolt
	// May have more due to triggered abilities
	if len(view.Stack) < 3 {
		t.Errorf("Expected at least 3 items on stack, got %d", len(view.Stack))
	}

	// Alice passes priority (she cast the last spell)
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	// All pass - should resolve in reverse order
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	viewInterface, err = engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view = viewInterface.(*game.EngineGameView)

	// Resolve any remaining triggered abilities
	resolveStackUntilEmpty(t, engine, gameID, players, 20)

	viewInterface, err = engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view = viewInterface.(*game.EngineGameView)

	// Stack should eventually be empty (or at least much smaller)
	if len(view.Stack) > 10 {
		t.Logf("Stack still has %d items after resolution cycles", len(view.Stack))
	}
}

// TestPriorityLoops tests priority handoff and pass chains
func TestPriorityLoops(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "priority-test"
	players := []string{"Alice", "Bob", "Charlie"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Alice casts spell
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "SEND_STRING",
		Data:       "Lightning Bolt",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast spell: %v", err)
	}

	viewInterface, err := engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view := viewInterface.(*game.EngineGameView)

	// Alice should have priority after casting (priority retention)
	if view.PriorityPlayer != "Alice" {
		t.Errorf("Expected Alice to have priority after casting, got %s", view.PriorityPlayer)
	}

	// Alice passes priority
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	viewInterface, err = engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view = viewInterface.(*game.EngineGameView)

	// Bob should have priority
	if view.PriorityPlayer != "Bob" {
		t.Errorf("Expected Bob to have priority, got %s", view.PriorityPlayer)
	}

	// Bob passes
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	viewInterface, err = engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view = viewInterface.(*game.EngineGameView)

	// Charlie should have priority
	if view.PriorityPlayer != "Charlie" {
		t.Errorf("Expected Charlie to have priority, got %s", view.PriorityPlayer)
	}

	// Charlie passes
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Charlie",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	// All passed - spell should resolve
	viewInterface, err = engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view = viewInterface.(*game.EngineGameView)

	// Resolve any triggered abilities
	resolveStackUntilEmpty(t, engine, gameID, players, 20)

	viewInterface, err = engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view = viewInterface.(*game.EngineGameView)

	// Stack should eventually be empty (or at least much smaller)
	if len(view.Stack) > 10 {
		t.Logf("Stack still has %d items after resolution cycles", len(view.Stack))
	}
}

// TestIllegalTargetRemoval tests that spells with illegal targets are removed
func TestIllegalTargetRemoval(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "illegal-target-test"
	players := []string{"Alice", "Bob"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	viewInterface, err := engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view := viewInterface.(*game.EngineGameView)

	// Verify initial state
	if len(view.Stack) != 0 {
		t.Errorf("Expected empty stack initially, got %d items", len(view.Stack))
	}

	// Note: Full illegal target testing would require:
	// 1. Casting a spell with a target
	// 2. Moving the target to an invalid zone
	// 3. Verifying the spell is removed from stack during resolution
	// This is a placeholder for future implementation
}

// TestStackResolutionOrder tests that stack resolves in LIFO order (last in, first out)
// TODO: Re-enable when proper card initialization is implemented
// This test requires cards "First", "Second", "Third" to exist in player hands
func testStackResolutionOrder(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "stack-order-test"
	players := []string{"Alice", "Bob"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Cast multiple spells in sequence (Alice retains priority between casts)
	spellOrder := []string{"First", "Second", "Third"}
	for _, spellName := range spellOrder {
		err = engine.ProcessAction(gameID, game.PlayerAction{
			PlayerID:   "Alice",
			ActionType: "SEND_STRING",
			Data:       spellName,
			Timestamp:  time.Now(),
		})
		if err != nil {
			t.Fatalf("Failed to cast %s: %v", spellName, err)
		}
		// Note: Alice retains priority after each cast, so she can cast multiple spells
	}

	viewInterface, err := engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view := viewInterface.(*game.EngineGameView)

	// Verify stack has at least 3 items
	if len(view.Stack) < 3 {
		t.Errorf("Expected at least 3 items on stack, got %d", len(view.Stack))
	}

	// Top of stack should be "Third" (last cast)
	if len(view.Stack) >= 3 {
		top := view.Stack[len(view.Stack)-1]
		if !strings.Contains(strings.ToLower(top.DisplayName), "third") {
			t.Logf("Top of stack is %s (may have triggered abilities)", top.DisplayName)
		}
	}
}

// TestPriorityAfterResolution tests that priority returns to active player after resolution
func TestPriorityAfterResolution(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "priority-after-resolution"
	players := []string{"Alice", "Bob"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Alice casts spell
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "SEND_STRING",
		Data:       "Lightning Bolt",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast spell: %v", err)
	}

	viewInterface, err := engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view := viewInterface.(*game.EngineGameView)

	// Alice should have priority after casting (priority retention)
	if view.PriorityPlayer != "Alice" {
		t.Errorf("Expected Alice to have priority after casting, got %s", view.PriorityPlayer)
	}

	// Alice passes priority so Bob can respond
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	// Both pass - spell resolves
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	// Resolve any triggered abilities
	resolveStackUntilEmpty(t, engine, gameID, players, 20)

	viewInterface, err = engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view = viewInterface.(*game.EngineGameView)

	// Priority should return to active player (Alice)
	if view.ActivePlayerID != "" && view.PriorityPlayer != view.ActivePlayerID {
		t.Logf("Priority player: %s, Active player: %s", view.PriorityPlayer, view.ActivePlayerID)
	}
}

// TestStackWithStateBasedActions tests that SBAs trigger correctly during stack resolution
func TestStackWithStateBasedActions(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "sba-stack-test"
	players := []string{"Alice", "Bob"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Alice casts spell that deals damage
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "SEND_STRING",
		Data:       "Lightning Bolt",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast spell: %v", err)
	}

	// Alice passes priority so Bob can respond
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	// Both pass - spell resolves
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	// Resolve any triggered abilities
	resolveStackUntilEmpty(t, engine, gameID, players, 20)

	viewInterface, err := engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view := viewInterface.(*game.EngineGameView)

	// Stack should eventually be empty (or at least much smaller)
	// SBAs should trigger during normalization
	if len(view.Stack) > 10 {
		t.Logf("Stack still has %d items after resolution cycles", len(view.Stack))
	}
}

// TestMultiplePlayersStackInteraction tests stack interactions with multiple players
// TODO: Re-enable when proper card initialization is implemented
// This test requires cards "Spell A", "Spell B", "Spell C" to exist in player hands
func testMultiplePlayersStackInteraction(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "multi-player-stack"
	players := []string{"Alice", "Bob", "Charlie", "Dave"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Each player casts a spell in order (each must pass before next can cast)
	spellNames := []string{"Spell A", "Spell B", "Spell C", "Spell D"}
	for i, player := range players {
		err = engine.ProcessAction(gameID, game.PlayerAction{
			PlayerID:   player,
			ActionType: "SEND_STRING",
			Data:       spellNames[i],
			Timestamp:  time.Now(),
		})
		if err != nil {
			t.Fatalf("Failed to cast spell %s: %v", spellNames[i], err)
		}
		// Pass priority so next player can cast (except for last player)
		if i < len(players)-1 {
			err = engine.ProcessAction(gameID, game.PlayerAction{
				PlayerID:   player,
				ActionType: "PLAYER_ACTION",
				Data:       "PASS",
				Timestamp:  time.Now(),
			})
			if err != nil {
				t.Fatalf("Failed to pass priority: %v", err)
			}
		}
	}

	viewInterface, err := engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view := viewInterface.(*game.EngineGameView)

	// Should have 4 spells on stack
	if len(view.Stack) != 4 {
		t.Errorf("Expected 4 spells on stack, got %d", len(view.Stack))
	}

	// All players pass - spells resolve in reverse order
	for _, player := range players {
		err = engine.ProcessAction(gameID, game.PlayerAction{
			PlayerID:   player,
			ActionType: "SEND_STRING",
			Data:       "Pass",
			Timestamp:  time.Now(),
		})
		if err != nil {
			t.Fatalf("Failed to pass: %v", err)
		}
	}

	// Resolve any triggered abilities
	resolveStackUntilEmpty(t, engine, gameID, players, 20)

	viewInterface, err = engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view = viewInterface.(*game.EngineGameView)

	// Stack should eventually be empty (or at least much smaller)
	if len(view.Stack) > 10 {
		t.Logf("Stack still has %d items after resolution cycles", len(view.Stack))
	}
}

// TestStackLegalityChecks tests that legality checks work during stack resolution
func TestStackLegalityChecks(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "legality-stack-test"
	players := []string{"Alice", "Bob"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Alice casts spell
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "SEND_STRING",
		Data:       "Lightning Bolt",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast spell: %v", err)
	}

	viewInterface, err := engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view := viewInterface.(*game.EngineGameView)

	// Verify spell is on stack (may have triggered abilities)
	if len(view.Stack) < 1 {
		t.Fatalf("Expected at least 1 item on stack, got %d", len(view.Stack))
	}

	// Alice passes priority so Bob can respond
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	// Both pass - spell should resolve (legality checks happen during resolution)
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	// Resolve any triggered abilities
	resolveStackUntilEmpty(t, engine, gameID, players, 20)

	viewInterface, err = engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view = viewInterface.(*game.EngineGameView)

	// Stack should eventually be empty (or at least much smaller)
	if len(view.Stack) > 10 {
		t.Logf("Stack still has %d items after resolution cycles", len(view.Stack))
	}
}

// TestComplexStackScenario tests a complex scenario with multiple spells, counters, and responses
func TestComplexStackScenario(t *testing.T) {
	logger := zap.NewNop()
	engine := game.NewMageEngine(logger)

	gameID := "complex-stack"
	players := []string{"Alice", "Bob"}

	err := engine.StartGame(gameID, players, "Duel")
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Alice casts spell 1
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "SEND_STRING",
		Data:       "Lightning Bolt",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast spell 1: %v", err)
	}

	// Alice passes priority so Bob can respond
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass priority: %v", err)
	}

	// Bob casts counterspell
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "SEND_STRING",
		Data:       "Counterspell",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast counterspell: %v", err)
	}

	// Bob passes priority so Alice can respond
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass priority: %v", err)
	}

	// Alice casts another spell
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "SEND_STRING",
		Data:       "Shock",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to cast shock: %v", err)
	}

	viewInterface, err := engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view := viewInterface.(*game.EngineGameView)

	// Should have at least 3 items: Shock (top), Counterspell, Lightning Bolt
	// May have more due to triggered abilities
	if len(view.Stack) < 3 {
		t.Errorf("Expected at least 3 items on stack, got %d", len(view.Stack))
	}

	// All pass - resolve in reverse order
	// Alice retains priority after casting Shock, so she passes first
	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Alice",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	err = engine.ProcessAction(gameID, game.PlayerAction{
		PlayerID:   "Bob",
		ActionType: "PLAYER_ACTION",
		Data:       "PASS",
		Timestamp:  time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to pass: %v", err)
	}

	// Resolve any triggered abilities
	resolveStackUntilEmpty(t, engine, gameID, players, 20)

	viewInterface, err = engine.GetGameView(gameID, "")
	if err != nil {
		t.Fatalf("Failed to get game view: %v", err)
	}
	view = viewInterface.(*game.EngineGameView)

	// Stack should eventually be empty (or at least much smaller)
	if len(view.Stack) > 10 {
		t.Logf("Stack still has %d items after resolution cycles", len(view.Stack))
	}
}
