package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpecialActionRestrictions(t *testing.T) {
	tests := []struct {
		name               string
		actionType         SpecialActionType
		wantMainPhase      bool
		wantEmptyStack     bool
		wantOwnTurn        bool
		wantOncePerGame    bool
	}{
		{
			name:           "play land",
			actionType:     SpecialActionPlayLand,
			wantMainPhase:  true,
			wantEmptyStack: true,
		},
		{
			name:       "turn face up",
			actionType: SpecialActionTurnFaceUp,
		},
		{
			name:       "end effect",
			actionType: SpecialActionEndEffect,
		},
		{
			name:       "ignore static ability",
			actionType: SpecialActionIgnoreStaticAbility,
		},
		{
			name:       "suspend",
			actionType: SpecialActionSuspend,
		},
		{
			name:            "companion",
			actionType:      SpecialActionCompanion,
			wantMainPhase:   true,
			wantEmptyStack:  true,
			wantOncePerGame: true,
		},
		{
			name:        "foretell",
			actionType:  SpecialActionForetell,
			wantOwnTurn: true,
		},
		{
			name:           "plot",
			actionType:     SpecialActionPlot,
			wantEmptyStack: true,
			wantOwnTurn:    true,
		},
		{
			name:           "unlock",
			actionType:     SpecialActionUnlock,
			wantMainPhase:  true,
			wantEmptyStack: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restrictions := GetRestrictions(tt.actionType)

			assert.Equal(t, tt.wantMainPhase, restrictions.RequiresMainPhase, "main phase")
			assert.Equal(t, tt.wantEmptyStack, restrictions.RequiresEmptyStack, "empty stack")
			assert.Equal(t, tt.wantOwnTurn, restrictions.RequiresOwnTurn, "own turn")
			assert.Equal(t, tt.wantOncePerGame, restrictions.OncePerGame, "once per game")
			assert.True(t, restrictions.RequiresPriority, "all require priority")
		})
	}
}

func TestSpecialActionManager(t *testing.T) {
	t.Run("play land - requires main phase and empty stack", func(t *testing.T) {
		sam := NewSpecialActionManager()

		executed := false
		action := SpecialAction{
			Type:     SpecialActionPlayLand,
			PlayerID: "player-1",
			SourceID: "forest-1",
			Execute: func() error {
				executed = true
				return nil
			},
		}

		// Cannot take during non-main phase
		canTake := sam.CanTakeAction(action, true, false, true, true)
		assert.False(t, canTake)

		// Cannot take with non-empty stack
		canTake = sam.CanTakeAction(action, true, true, false, true)
		assert.False(t, canTake)

		// Can take during main phase with empty stack
		canTake = sam.CanTakeAction(action, true, true, true, true)
		assert.True(t, canTake)

		// Execute the action
		err := sam.TakeAction(action)
		require.NoError(t, err)
		assert.True(t, executed)

		// Should track that land was played
		count := sam.GetActionsTakenThisTurn("player-1", SpecialActionPlayLand)
		assert.Equal(t, 1, count)
	})

	t.Run("turn face up - can take anytime with priority", func(t *testing.T) {
		sam := NewSpecialActionManager()

		executed := false
		action := SpecialAction{
			Type:     SpecialActionTurnFaceUp,
			PlayerID: "player-1",
			SourceID: "creature-1",
			Execute: func() error {
				executed = true
				return nil
			},
		}

		// Can take during combat
		canTake := sam.CanTakeAction(action, true, false, false, true)
		assert.True(t, canTake)

		// Can take with non-empty stack
		canTake = sam.CanTakeAction(action, true, false, false, true)
		assert.True(t, canTake)

		// Cannot take without priority
		canTake = sam.CanTakeAction(action, false, false, false, true)
		assert.False(t, canTake)

		err := sam.TakeAction(action)
		require.NoError(t, err)
		assert.True(t, executed)
	})

	t.Run("companion - once per game restriction", func(t *testing.T) {
		sam := NewSpecialActionManager()

		executed := 0
		action := SpecialAction{
			Type:     SpecialActionCompanion,
			PlayerID: "player-1",
			SourceID: "companion-1",
			Execute: func() error {
				executed++
				return nil
			},
		}

		// Can take first time
		canTake := sam.CanTakeAction(action, true, true, true, true)
		assert.True(t, canTake)

		err := sam.TakeAction(action)
		require.NoError(t, err)
		assert.Equal(t, 1, executed)

		// Cannot take again (once per game)
		canTake = sam.CanTakeAction(action, true, true, true, true)
		assert.False(t, canTake)
	})

	t.Run("foretell - requires own turn", func(t *testing.T) {
		sam := NewSpecialActionManager()

		action := SpecialAction{
			Type:     SpecialActionForetell,
			PlayerID: "player-1",
			SourceID: "card-1",
			Execute:  func() error { return nil },
		}

		// Cannot take during opponent's turn
		canTake := sam.CanTakeAction(action, true, false, false, false)
		assert.False(t, canTake)

		// Can take during own turn
		canTake = sam.CanTakeAction(action, true, false, false, true)
		assert.True(t, canTake)
	})

	t.Run("custom restrictions", func(t *testing.T) {
		sam := NewSpecialActionManager()

		allowed := true
		action := SpecialAction{
			Type:     SpecialActionSuspend,
			PlayerID: "player-1",
			SourceID: "card-1",
			Execute:  func() error { return nil },
			CanTake: func() bool {
				return allowed
			},
		}

		// Custom check passes
		canTake := sam.CanTakeAction(action, true, false, false, false)
		assert.True(t, canTake)

		// Custom check fails
		allowed = false
		canTake = sam.CanTakeAction(action, true, false, false, false)
		assert.False(t, canTake)
	})

	t.Run("reset turn clears per-turn tracking", func(t *testing.T) {
		sam := NewSpecialActionManager()

		action := SpecialAction{
			Type:     SpecialActionPlayLand,
			PlayerID: "player-1",
			SourceID: "land-1",
			Execute:  func() error { return nil },
		}

		sam.TakeAction(action)
		assert.Equal(t, 1, sam.GetActionsTakenThisTurn("player-1", SpecialActionPlayLand))

		sam.ResetTurn()
		assert.Equal(t, 0, sam.GetActionsTakenThisTurn("player-1", SpecialActionPlayLand))
	})

	t.Run("multiple players tracking", func(t *testing.T) {
		sam := NewSpecialActionManager()

		action1 := SpecialAction{
			Type:     SpecialActionPlayLand,
			PlayerID: "player-1",
			SourceID: "land-1",
			Execute:  func() error { return nil },
		}

		action2 := SpecialAction{
			Type:     SpecialActionPlayLand,
			PlayerID: "player-2",
			SourceID: "land-2",
			Execute:  func() error { return nil },
		}

		sam.TakeAction(action1)
		sam.TakeAction(action2)

		assert.Equal(t, 1, sam.GetActionsTakenThisTurn("player-1", SpecialActionPlayLand))
		assert.Equal(t, 1, sam.GetActionsTakenThisTurn("player-2", SpecialActionPlayLand))
	})
}

func TestSpecialActionsRuleCoverage(t *testing.T) {
	t.Run("Rule 116.3 - player receives priority after special action", func(t *testing.T) {
		// This test documents that per Rule 116.3, the player who takes a special action
		// receives priority afterward. The actual priority handling is in the game engine,
		// but we verify the special action executes successfully.

		sam := NewSpecialActionManager()

		playerGotPriority := false
		action := SpecialAction{
			Type:     SpecialActionTurnFaceUp,
			PlayerID: "player-1",
			SourceID: "morph-creature",
			Execute: func() error {
				// After this executes, player-1 should receive priority
				playerGotPriority = true
				return nil
			},
		}

		err := sam.TakeAction(action)
		require.NoError(t, err)
		assert.True(t, playerGotPriority)
	})
}
