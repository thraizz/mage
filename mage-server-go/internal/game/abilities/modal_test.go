package abilities

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestModalSpellCreation tests basic modal spell creation
func TestModalSpellCreation(t *testing.T) {
	sourceID := uuid.New()

	// Create a simple modal spell (like a Charm)
	builder := NewModalSpellBuilder(sourceID, "{2}{U}")
	builder.AddMode("Counter target spell", []Effect{
		&MockEffect{text: "Counter target spell"},
	})
	builder.AddMode("Draw a card", []Effect{
		&MockEffect{text: "Draw a card"},
	})
	builder.AddMode("Return target permanent to its owner's hand", []Effect{
		&MockEffect{text: "Return to hand"},
	})

	spell, err := builder.Build()
	require.NoError(t, err)
	require.NotNil(t, spell)

	assert.Equal(t, 3, len(spell.Modes))
	assert.Equal(t, 1, spell.MinModes)
	assert.Equal(t, 1, spell.MaxModes)
}

// TestModalSpellMultipleModes tests choosing multiple modes
func TestModalSpellMultipleModes(t *testing.T) {
	sourceID := uuid.New()

	// Create Cryptic Command-like spell (choose two)
	builder := NewModalSpellBuilder(sourceID, "{1}{U}{U}{U}")
	builder.SetModeRange(2, 2) // Choose exactly 2
	builder.AddMode("Counter target spell", []Effect{
		&MockEffect{text: "Counter target spell"},
	})
	builder.AddMode("Return target permanent to its owner's hand", []Effect{
		&MockEffect{text: "Return to hand"},
	})
	builder.AddMode("Tap all creatures your opponents control", []Effect{
		&MockEffect{text: "Tap creatures"},
	})
	builder.AddMode("Draw a card", []Effect{
		&MockEffect{text: "Draw a card"},
	})

	spell, err := builder.Build()
	require.NoError(t, err)

	assert.Equal(t, 4, len(spell.Modes))
	assert.Equal(t, 2, spell.MinModes)
	assert.Equal(t, 2, spell.MaxModes)

	// Test valid mode selection (2 modes)
	err = spell.SetChosenModes([]int{1, 4}) // Counter + Draw
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 4}, spell.GetChosenModes())

	// Test invalid: too few modes
	err = spell.SetChosenModes([]int{1})
	assert.Error(t, err)

	// Test invalid: too many modes
	err = spell.SetChosenModes([]int{1, 2, 3})
	assert.Error(t, err)

	// Test invalid: duplicate modes (should fail by default)
	err = spell.SetChosenModes([]int{1, 1})
	assert.Error(t, err)
}

// TestModalSpellAllowDuplicates tests allowing duplicate mode selection
func TestModalSpellAllowDuplicates(t *testing.T) {
	sourceID := uuid.New()

	// Some cards allow choosing the same mode multiple times
	builder := NewModalSpellBuilder(sourceID, "{2}{R}")
	builder.SetModeRange(2, 2)
	builder.AllowDuplicateModes()
	builder.AddMode("Deal 2 damage to any target", []Effect{
		&MockEffect{text: "Deal 2 damage"},
	})
	builder.AddMode("Destroy target artifact", []Effect{
		&MockEffect{text: "Destroy artifact"},
	})

	spell, err := builder.Build()
	require.NoError(t, err)

	// Should allow choosing the same mode twice
	err = spell.SetChosenModes([]int{1, 1})
	assert.NoError(t, err)
}

// TestModalSpellAvailableModes tests mode restrictions
func TestModalSpellAvailableModes(t *testing.T) {
	sourceID := uuid.New()

	// Create spell with conditional modes
	builder := NewModalSpellBuilder(sourceID, "{3}{G}")

	// Mode 1: Always available
	builder.AddMode("Put a +1/+1 counter on target creature", []Effect{
		&MockEffect{text: "+1/+1 counter"},
	})

	// Mode 2: Conditional (only if condition is met)
	builder.AddModeWithRestriction(
		"Destroy target artifact or enchantment",
		[]Effect{&MockEffect{text: "Destroy artifact/enchantment"}},
		&ConditionalRestriction{
			Condition: func(ctx context.Context, game GameContext, source uuid.UUID) bool {
				// Example: only available if there are artifacts/enchantments
				return true // Simplified for test
			},
			Text: "Only if you control three or more creatures",
		},
	)

	spell, err := builder.Build()
	require.NoError(t, err)

	// Get available modes
	ctx := context.Background()
	available := spell.GetAvailableModes(ctx, nil)

	// Both should be available in this test
	assert.Equal(t, 2, len(available))
}

// TestModalSpellWithTargets tests modes with different target requirements
func TestModalSpellWithTargets(t *testing.T) {
	sourceID := uuid.New()

	builder := NewModalSpellBuilder(sourceID, "{1}{W}")

	// Mode 1: Requires creature target
	builder.AddModeWithTargets(
		"Destroy target creature",
		[]Effect{&MockEffect{text: "Destroy"}},
		NewTargetRequirement(1, 1, NewCreatureTargetFilter()),
	)

	// Mode 2: No targets
	builder.AddMode("Create a 1/1 white Soldier creature token", []Effect{
		&MockEffect{text: "Create token"},
	})

	spell, err := builder.Build()
	require.NoError(t, err)

	// Get targets for mode 1
	targets, err := spell.GetTargetsForMode(1)
	assert.NoError(t, err)
	assert.NotNil(t, targets)

	// Get targets for mode 2 (should be nil)
	targets, err = spell.GetTargetsForMode(2)
	assert.NoError(t, err)
	assert.Nil(t, targets)
}

// TestModalSpellCanActivate tests CanActivate with mode restrictions
func TestModalSpellCanActivate(t *testing.T) {
	sourceID := uuid.New()

	builder := NewModalSpellBuilder(sourceID, "{2}{B}")
	builder.SetModeRange(1, 2) // Choose one or two

	// Add mode that's always available
	builder.AddMode("Each opponent loses 2 life", []Effect{
		&MockEffect{text: "Lose life"},
	})

	// Add mode with restriction
	builder.AddModeWithRestriction(
		"Return target creature card from your graveyard to your hand",
		[]Effect{&MockEffect{text: "Return from graveyard"}},
		&ConditionalRestriction{
			Condition: func(ctx context.Context, game GameContext, source uuid.UUID) bool {
				return false // Not available
			},
			Text: "Only if a creature died this turn",
		},
	)

	spell, err := builder.Build()
	require.NoError(t, err)

	ctx := context.Background()

	// Should be able to activate (at least 1 mode available)
	canActivate := spell.CanActivate(ctx, nil)
	assert.True(t, canActivate)

	// But can only choose mode 1
	err = spell.SetChosenModes([]int{1})
	assert.NoError(t, err)

	// Cannot choose mode 2 (restricted)
	available := spell.GetAvailableModes(ctx, nil)
	assert.Equal(t, 1, len(available))
	assert.Equal(t, 1, available[0].ID)
}

// TestModalSpellResolve tests resolution of chosen modes
func TestModalSpellResolve(t *testing.T) {
	sourceID := uuid.New()

	effectsApplied := []string{}

	builder := NewModalSpellBuilder(sourceID, "{2}{U}{U}")
	builder.SetModeRange(2, 2)

	builder.AddMode("Counter target spell", []Effect{
		&MockEffectWithCallback{
			text: "Counter",
			callback: func() {
				effectsApplied = append(effectsApplied, "Counter")
			},
		},
	})

	builder.AddMode("Draw a card", []Effect{
		&MockEffectWithCallback{
			text: "Draw",
			callback: func() {
				effectsApplied = append(effectsApplied, "Draw")
			},
		},
	})

	builder.AddMode("Return target permanent to its owner's hand", []Effect{
		&MockEffectWithCallback{
			text: "Bounce",
			callback: func() {
				effectsApplied = append(effectsApplied, "Bounce")
			},
		},
	})

	spell, err := builder.Build()
	require.NoError(t, err)

	// Choose modes 1 and 2
	err = spell.SetChosenModes([]int{1, 2})
	require.NoError(t, err)

	// Resolve
	ctx := context.Background()
	err = spell.Resolve(ctx, nil)
	assert.NoError(t, err)

	// Both effects should have been applied in order
	assert.Equal(t, []string{"Counter", "Draw"}, effectsApplied)
}

// TestModalSpellInvalidModeID tests error handling for invalid mode IDs
func TestModalSpellInvalidModeID(t *testing.T) {
	sourceID := uuid.New()

	builder := NewModalSpellBuilder(sourceID, "{1}{G}")
	builder.AddMode("Target creature gets +3/+3 until end of turn", []Effect{
		&MockEffect{text: "+3/+3"},
	})

	spell, err := builder.Build()
	require.NoError(t, err)

	// Try to choose non-existent mode
	err = spell.SetChosenModes([]int{5})
	assert.Error(t, err)

	// Try to get non-existent mode
	mode, err := spell.GetMode(99)
	assert.Error(t, err)
	assert.Nil(t, mode)
}

// TestModalSpellString tests string representation
func TestModalSpellString(t *testing.T) {
	sourceID := uuid.New()

	builder := NewModalSpellBuilder(sourceID, "{2}{U}")
	builder.AddMode("Counter target spell", []Effect{
		&MockEffect{text: "Counter"},
	})
	builder.AddMode("Draw a card", []Effect{
		&MockEffect{text: "Draw"},
	})

	spell, err := builder.Build()
	require.NoError(t, err)

	str := spell.String()
	assert.Contains(t, str, "Choose one")
	assert.Contains(t, str, "Counter target spell")
	assert.Contains(t, str, "Draw a card")
}

// Mock effect for testing
type MockEffect struct {
	text string
}

func (e *MockEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	return nil
}

func (e *MockEffect) String() string {
	return e.text
}

// Mock effect with callback for testing resolution
type MockEffectWithCallback struct {
	text     string
	callback func()
}

func (e *MockEffectWithCallback) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	if e.callback != nil {
		e.callback()
	}
	return nil
}

func (e *MockEffectWithCallback) String() string {
	return e.text
}
