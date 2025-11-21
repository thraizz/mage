package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolutionContext(t *testing.T) {
	t.Run("basic resolution tracking", func(t *testing.T) {
		rc := NewResolutionContext()

		assert.False(t, rc.IsResolving())
		assert.Equal(t, 0, rc.GetDepth())

		// Begin first resolution
		err := rc.BeginResolution("spell-1")
		require.NoError(t, err)
		assert.True(t, rc.IsResolving())
		assert.Equal(t, 1, rc.GetDepth())
		assert.Equal(t, "spell-1", rc.GetCurrentResolvingID())

		// End resolution
		err = rc.EndResolution("spell-1")
		require.NoError(t, err)
		assert.False(t, rc.IsResolving())
		assert.Equal(t, 0, rc.GetDepth())
	})

	t.Run("nested resolution", func(t *testing.T) {
		rc := NewResolutionContext()

		// Begin first resolution
		err := rc.BeginResolution("spell-1")
		require.NoError(t, err)

		// Begin nested resolution (casting copy)
		err = rc.BeginResolution("spell-2")
		require.NoError(t, err)
		assert.Equal(t, 2, rc.GetDepth())
		assert.Equal(t, "spell-2", rc.GetCurrentResolvingID())

		// End nested resolution
		err = rc.EndResolution("spell-2")
		require.NoError(t, err)
		assert.Equal(t, 1, rc.GetDepth())
		assert.Equal(t, "spell-1", rc.GetCurrentResolvingID())

		// End first resolution
		err = rc.EndResolution("spell-1")
		require.NoError(t, err)
		assert.Equal(t, 0, rc.GetDepth())
	})

	t.Run("maximum depth limit", func(t *testing.T) {
		rc := NewResolutionContext()

		// Fill to max depth
		for i := 0; i < rc.maxDepth; i++ {
			err := rc.BeginResolution("spell")
			require.NoError(t, err)
		}

		// Exceeding max depth should fail
		err := rc.BeginResolution("spell")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "maximum resolution depth")
	})

	t.Run("mana ability permission", func(t *testing.T) {
		rc := NewResolutionContext()

		assert.False(t, rc.CanActivateManaAbilities())

		rc.SetAllowManaAbilities(true)
		assert.True(t, rc.CanActivateManaAbilities())

		rc.SetAllowManaAbilities(false)
		assert.False(t, rc.CanActivateManaAbilities())
	})

	t.Run("reset", func(t *testing.T) {
		rc := NewResolutionContext()

		rc.BeginResolution("spell-1")
		rc.SetAllowManaAbilities(true)
		rc.SetAllowSpecialActions(true)

		rc.Reset()

		assert.False(t, rc.IsResolving())
		assert.Equal(t, 0, rc.GetDepth())
		assert.False(t, rc.CanActivateManaAbilities())
		assert.False(t, rc.CanTakeSpecialActions())
	})
}

func TestPriorityWindowManager(t *testing.T) {
	t.Run("open and close window", func(t *testing.T) {
		pwm := NewPriorityWindowManager()

		window := PriorityWindow{
			Type:     PriorityWindowManaPayment,
			PlayerID: "player-1",
			Context:  "Paying for Lightning Bolt",
			AllowedActions: []ActionType{ActionActivateMana},
		}

		err := pwm.OpenWindow(window)
		require.NoError(t, err)

		active := pwm.GetActiveWindow()
		require.NotNil(t, active)
		assert.Equal(t, PriorityWindowManaPayment, active.Type)

		pwm.CloseWindow()
		assert.Nil(t, pwm.GetActiveWindow())
	})

	t.Run("cannot open multiple windows", func(t *testing.T) {
		pwm := NewPriorityWindowManager()

		window1 := PriorityWindow{
			Type:     PriorityWindowManaPayment,
			PlayerID: "player-1",
		}

		err := pwm.OpenWindow(window1)
		require.NoError(t, err)

		window2 := PriorityWindow{
			Type:     PriorityWindowChoice,
			PlayerID: "player-1",
		}

		err = pwm.OpenWindow(window2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already open")
	})

	t.Run("action allowed check", func(t *testing.T) {
		pwm := NewPriorityWindowManager()

		window := PriorityWindow{
			Type:     PriorityWindowManaPayment,
			PlayerID: "player-1",
			AllowedActions: []ActionType{
				ActionActivateMana,
				ActionSpecialAction,
			},
		}

		pwm.OpenWindow(window)

		assert.True(t, pwm.IsActionAllowed(ActionActivateMana))
		assert.True(t, pwm.IsActionAllowed(ActionSpecialAction))
		assert.False(t, pwm.IsActionAllowed(ActionCastSpell))
	})
}

func TestResolutionWithPriorityWindows(t *testing.T) {
	t.Run("mana ability during payment", func(t *testing.T) {
		rc := NewResolutionContext()
		pwm := NewPriorityWindowManager()

		// Not resolving yet, can't open payment window
		assert.False(t, rc.IsResolving())

		// Start resolving a spell
		rc.BeginResolution("lightning-bolt")
		rc.SetAllowManaAbilities(true)

		// Open mana payment window
		window := PriorityWindow{
			Type:     PriorityWindowManaPayment,
			PlayerID: "player-1",
			Context:  "Paying {R} for Lightning Bolt",
			AllowedActions: []ActionType{ActionActivateMana},
		}
		pwm.OpenWindow(window)

		// Can activate mana abilities in this window
		assert.True(t, rc.CanActivateManaAbilities())
		assert.True(t, pwm.IsActionAllowed(ActionActivateMana))

		// Close window and end resolution
		pwm.CloseWindow()
		rc.EndResolution("lightning-bolt")

		assert.False(t, rc.IsResolving())
	})

	t.Run("nested resolution depth tracking", func(t *testing.T) {
		rc := NewResolutionContext()

		// Resolve ability that creates copy
		rc.BeginResolution("isochron-scepter-ability")
		assert.Equal(t, 1, rc.GetDepth())

		// Nested: cast the copy (Rule 707.12)
		rc.BeginResolution("lightning-bolt-copy")
		assert.Equal(t, 2, rc.GetDepth())

		// End copy resolution
		rc.EndResolution("lightning-bolt-copy")
		assert.Equal(t, 1, rc.GetDepth())

		// End original resolution
		rc.EndResolution("isochron-scepter-ability")
		assert.Equal(t, 0, rc.GetDepth())
	})
}
