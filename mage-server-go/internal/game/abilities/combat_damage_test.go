package abilities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTrampleDamageCalculator tests trample damage calculations
func TestTrampleDamageCalculator(t *testing.T) {
	tests := []struct {
		name               string
		attackerPower      int
		hasDeathtouch      bool
		blockers           []BlockerInfo
		expectedAssignment map[uuid.UUID]int
		expectedTrample    int
	}{
		{
			name:          "Simple trample - one blocker",
			attackerPower: 5,
			hasDeathtouch: false,
			blockers: []BlockerInfo{
				{BlockerID: uuid.New(), Toughness: 3, DamageMarked: 0, HasIndestruc: false},
			},
			expectedTrample: 2, // 5 - 3 = 2 tramples through
		},
		{
			name:          "Trample with deathtouch - one blocker",
			attackerPower: 5,
			hasDeathtouch: true,
			blockers: []BlockerInfo{
				{BlockerID: uuid.New(), Toughness: 3, DamageMarked: 0, HasIndestruc: false},
			},
			expectedTrample: 4, // Only 1 damage needed (deathtouch), 4 tramples
		},
		{
			name:          "Trample with deathtouch - multiple blockers",
			attackerPower: 5,
			hasDeathtouch: true,
			blockers: []BlockerInfo{
				{BlockerID: uuid.New(), Toughness: 3, DamageMarked: 0, HasIndestruc: false},
				{BlockerID: uuid.New(), Toughness: 4, DamageMarked: 0, HasIndestruc: false},
			},
			expectedTrample: 3, // 1 to each blocker, 3 tramples
		},
		{
			name:          "Trample with indestructible blocker",
			attackerPower: 5,
			hasDeathtouch: false,
			blockers: []BlockerInfo{
				{BlockerID: uuid.New(), Toughness: 4, DamageMarked: 0, HasIndestruc: true},
			},
			expectedTrample: 1, // Still need 4 damage to indestructible
		},
		{
			name:          "Trample with already damaged blocker",
			attackerPower: 5,
			hasDeathtouch: false,
			blockers: []BlockerInfo{
				{BlockerID: uuid.New(), Toughness: 5, DamageMarked: 3, HasIndestruc: false},
			},
			expectedTrample: 3, // Only 2 more damage needed, 3 tramples
		},
		{
			name:          "No trample - not enough power",
			attackerPower: 3,
			hasDeathtouch: false,
			blockers: []BlockerInfo{
				{BlockerID: uuid.New(), Toughness: 5, DamageMarked: 0, HasIndestruc: false},
			},
			expectedTrample: 0, // Not enough to kill blocker
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attackerID := uuid.New()
			calc := NewTrampleDamageCalculator(attackerID, tt.attackerPower, tt.hasDeathtouch)

			// Add blockers
			for _, blocker := range tt.blockers {
				calc.AddBlocker(blocker.BlockerID, blocker.Toughness, blocker.DamageMarked, blocker.HasIndestruc)
			}

			// Calculate damage
			assignment, trampleDamage := calc.CalculateTrampleDamage()

			// Verify trample damage
			assert.Equal(t, tt.expectedTrample, trampleDamage, "Incorrect trample damage")

			// Verify total assignment doesn't exceed power
			totalAssigned := 0
			for _, damage := range assignment {
				totalAssigned += damage
			}
			assert.LessOrEqual(t, totalAssigned, tt.attackerPower, "Assigned more damage than power")
		})
	}
}

// TestTrampleDamageCalculator_Validation tests damage assignment validation
func TestTrampleDamageCalculator_Validation(t *testing.T) {
	attackerID := uuid.New()
	blocker1 := uuid.New()
	blocker2 := uuid.New()

	calc := NewTrampleDamageCalculator(attackerID, 5, false)
	calc.AddBlocker(blocker1, 3, 0, false) // Needs 3 lethal
	calc.AddBlocker(blocker2, 2, 0, false) // Needs 2 lethal

	tests := []struct {
		name        string
		assignment  map[uuid.UUID]int
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid - exactly lethal to both",
			assignment: map[uuid.UUID]int{
				blocker1: 3,
				blocker2: 2,
			},
			expectError: false,
		},
		{
			name: "Valid - more than lethal to first, trample remaining",
			assignment: map[uuid.UUID]int{
				blocker1: 4,
				blocker2: 1,
			},
			expectError: true, // Not enough to blocker2
			errorMsg:    "must assign at least",
		},
		{
			name: "Invalid - not enough to first blocker",
			assignment: map[uuid.UUID]int{
				blocker1: 2, // Need 3
				blocker2: 2,
			},
			expectError: true,
			errorMsg:    "must assign at least 3 lethal damage",
		},
		{
			name: "Invalid - too much total damage",
			assignment: map[uuid.UUID]int{
				blocker1: 4,
				blocker2: 2,
			},
			expectError: true,
			errorMsg:    "cannot assign more damage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := calc.ValidateTrampleAssignment(tt.assignment)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCalculateLethalDamage tests lethal damage calculations
func TestCalculateLethalDamage(t *testing.T) {
	tests := []struct {
		name              string
		toughness         int
		markedDamage      int
		hasDeathtouch     bool
		hasIndestructible bool
		expected          int
	}{
		{
			name:              "Undamaged creature",
			toughness:         3,
			markedDamage:      0,
			hasDeathtouch:     false,
			hasIndestructible: false,
			expected:          3,
		},
		{
			name:              "Partially damaged creature",
			toughness:         5,
			markedDamage:      2,
			hasDeathtouch:     false,
			hasIndestructible: false,
			expected:          3,
		},
		{
			name:              "Already lethal damage",
			toughness:         3,
			markedDamage:      3,
			hasDeathtouch:     false,
			hasIndestructible: false,
			expected:          0,
		},
		{
			name:              "Deathtouch - undamaged",
			toughness:         10,
			markedDamage:      0,
			hasDeathtouch:     true,
			hasIndestructible: false,
			expected:          1, // Deathtouch makes 1 damage lethal
		},
		{
			name:              "Deathtouch - already damaged",
			toughness:         5,
			markedDamage:      1,
			hasDeathtouch:     true,
			hasIndestructible: false,
			expected:          0, // Already has deathtouch damage
		},
		{
			name:              "Indestructible - normal calculation",
			toughness:         4,
			markedDamage:      1,
			hasDeathtouch:     false,
			hasIndestructible: true,
			expected:          3, // Indestructible still needs full toughness
		},
		{
			name:              "Indestructible with deathtouch",
			toughness:         5,
			markedDamage:      0,
			hasDeathtouch:     true,
			hasIndestructible: true,
			expected:          5, // Indestructible ignores deathtouch
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateLethalDamage(tt.toughness, tt.markedDamage, tt.hasDeathtouch, tt.hasIndestructible)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCombatDamageContext tests damage event processing
func TestCombatDamageContext(t *testing.T) {
	t.Run("No modifications", func(t *testing.T) {
		ctx := NewCombatDamageContext(false)

		source := uuid.New()
		target := uuid.New()

		ctx.AddDamageEvent(CombatDamageEvent{
			sourceID:    source,
			targetID:    target,
			amount:      5,
			isCombat:    true,
			finalAmount: 5,
		})

		processed := ctx.ProcessDamageEvents()
		require.Len(t, processed, 1)
		assert.Equal(t, 5, processed[0].finalAmount)
		assert.Equal(t, 0, processed[0].prevented)
	})

	t.Run("Prevention effect", func(t *testing.T) {
		ctx := NewCombatDamageContext(false)

		source := uuid.New()
		target := uuid.New()

		ctx.AddDamageEvent(CombatDamageEvent{
			sourceID:    source,
			targetID:    target,
			amount:      5,
			isCombat:    true,
			finalAmount: 5,
		})

		// Add prevention effect that prevents 3 damage
		ctx.AddPreventionEffect(PreventionEffect{
			effectID: uuid.New(),
			sourceID: uuid.New(),
			applies: func(event *CombatDamageEvent) bool {
				return event.targetID == target
			},
			preventAmount: 3,
			usedOnce:      false,
		})

		processed := ctx.ProcessDamageEvents()
		require.Len(t, processed, 1)
		assert.Equal(t, 2, processed[0].finalAmount, "Should prevent 3 damage")
		assert.Equal(t, 3, processed[0].prevented)
	})

	t.Run("Prevent all damage", func(t *testing.T) {
		ctx := NewCombatDamageContext(false)

		source := uuid.New()
		target := uuid.New()

		ctx.AddDamageEvent(CombatDamageEvent{
			sourceID:    source,
			targetID:    target,
			amount:      5,
			isCombat:    true,
			finalAmount: 5,
		})

		// Prevent all damage (-1)
		ctx.AddPreventionEffect(PreventionEffect{
			effectID: uuid.New(),
			sourceID: uuid.New(),
			applies: func(event *CombatDamageEvent) bool {
				return true
			},
			preventAmount: -1, // All damage
			usedOnce:      false,
		})

		processed := ctx.ProcessDamageEvents()
		assert.Len(t, processed, 0, "All damage prevented - no events")
	})

	t.Run("Replacement effect - doubling", func(t *testing.T) {
		ctx := NewCombatDamageContext(false)

		source := uuid.New()
		target := uuid.New()

		ctx.AddDamageEvent(CombatDamageEvent{
			sourceID:    source,
			targetID:    target,
			amount:      5,
			isCombat:    true,
			finalAmount: 5,
		})

		// Double damage replacement
		ctx.AddReplacementEffect(ReplacementEffect{
			effectID: uuid.New(),
			sourceID: uuid.New(),
			applies: func(event *CombatDamageEvent) bool {
				return true
			},
			replace: func(event *CombatDamageEvent) *CombatDamageEvent {
				event.amount = event.amount * 2
				return event
			},
			usedOnce: false,
			maxUses:  -1,
		})

		processed := ctx.ProcessDamageEvents()
		require.Len(t, processed, 1)
		assert.Equal(t, 10, processed[0].finalAmount, "Damage should be doubled")
	})

	t.Run("Replacement then prevention", func(t *testing.T) {
		ctx := NewCombatDamageContext(false)

		source := uuid.New()
		target := uuid.New()

		ctx.AddDamageEvent(CombatDamageEvent{
			sourceID:    source,
			targetID:    target,
			amount:      5,
			isCombat:    true,
			finalAmount: 5,
		})

		// First: double damage (replacement)
		ctx.AddReplacementEffect(ReplacementEffect{
			effectID: uuid.New(),
			sourceID: uuid.New(),
			applies: func(event *CombatDamageEvent) bool {
				return true
			},
			replace: func(event *CombatDamageEvent) *CombatDamageEvent {
				event.amount = event.amount * 2
				return event
			},
			usedOnce: false,
			maxUses:  -1,
		})

		// Then: prevent 7 damage (prevention)
		ctx.AddPreventionEffect(PreventionEffect{
			effectID: uuid.New(),
			sourceID: uuid.New(),
			applies: func(event *CombatDamageEvent) bool {
				return true
			},
			preventAmount: 7,
			usedOnce:      false,
		})

		processed := ctx.ProcessDamageEvents()
		require.Len(t, processed, 1)
		// 5 doubled = 10, then prevent 7 = 3 final
		assert.Equal(t, 3, processed[0].finalAmount)
		assert.Equal(t, 7, processed[0].prevented)
	})
}

// TestDamageAssignmentOrder tests damage ordering
func TestDamageAssignmentOrder(t *testing.T) {
	source := uuid.New()
	target1 := uuid.New()
	target2 := uuid.New()
	target3 := uuid.New()

	t.Run("Create and get order", func(t *testing.T) {
		order := NewDamageAssignmentOrder(source, []uuid.UUID{target1, target2, target3}, true)
		targets := order.GetTargets()

		assert.Len(t, targets, 3)
		assert.Equal(t, target1, targets[0])
		assert.Equal(t, target2, targets[1])
		assert.Equal(t, target3, targets[2])
	})

	t.Run("Set new order", func(t *testing.T) {
		order := NewDamageAssignmentOrder(source, []uuid.UUID{target1, target2, target3}, true)

		// Reverse order
		err := order.SetOrder([]uuid.UUID{target3, target2, target1})
		require.NoError(t, err)

		targets := order.GetTargets()
		assert.Equal(t, target3, targets[0])
		assert.Equal(t, target2, targets[1])
		assert.Equal(t, target1, targets[2])
	})

	t.Run("Invalid order - wrong count", func(t *testing.T) {
		order := NewDamageAssignmentOrder(source, []uuid.UUID{target1, target2, target3}, true)

		err := order.SetOrder([]uuid.UUID{target1, target2}) // Missing one
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must specify order for all")
	})

	t.Run("Invalid order - wrong target", func(t *testing.T) {
		order := NewDamageAssignmentOrder(source, []uuid.UUID{target1, target2, target3}, true)

		wrongTarget := uuid.New()
		err := order.SetOrder([]uuid.UUID{target1, target2, wrongTarget})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not in original target list")
	})
}

// TestPreventCombatDamageEffect tests prevention effects
func TestPreventCombatDamageEffect(t *testing.T) {
	source := uuid.New()
	target := uuid.New()

	t.Run("Prevent fixed amount", func(t *testing.T) {
		effect := NewPreventCombatDamageEffect(source, target, 3, DurationUntilEndOfTurn)

		assert.Equal(t, 3, effect.amount)
		assert.Equal(t, target, effect.targetID)
		assert.Equal(t, DurationUntilEndOfTurn, effect.duration)
	})

	t.Run("Prevent all damage", func(t *testing.T) {
		effect := NewPreventCombatDamageEffect(source, target, -1, DurationUntilEndOfTurn)

		assert.Equal(t, -1, effect.amount)
		assert.Contains(t, effect.GetDescription(), "all combat damage")
	})
}

// TestReplaceCombatDamageEffect tests replacement effects
func TestReplaceCombatDamageEffect(t *testing.T) {
	source := uuid.New()
	target := uuid.New()

	t.Run("Double damage", func(t *testing.T) {
		effect := NewReplaceCombatDamageEffect(source, target, 2.0, 0, DurationPermanent)

		assert.Equal(t, 2.0, effect.multiplier)
		assert.Equal(t, target, effect.targetID)
	})

	t.Run("Add damage", func(t *testing.T) {
		effect := NewReplaceCombatDamageEffect(source, target, 1.0, 3, DurationUntilEndOfTurn)

		assert.Equal(t, 1.0, effect.multiplier)
		assert.Equal(t, 3, effect.addAmount)
	})
}

// TestRedirectCombatDamageEffect tests damage redirection
func TestRedirectCombatDamageEffect(t *testing.T) {
	source := uuid.New()
	fromTarget := uuid.New()
	toTarget := uuid.New()

	t.Run("Redirect all damage", func(t *testing.T) {
		effect := NewRedirectCombatDamageEffect(source, fromTarget, toTarget, -1, DurationUntilEndOfTurn)

		assert.Equal(t, -1, effect.maxAmount)
		assert.Equal(t, fromTarget, effect.fromTarget)
		assert.Equal(t, toTarget, effect.toTarget)
		assert.Contains(t, effect.GetDescription(), "all damage")
	})

	t.Run("Redirect limited damage", func(t *testing.T) {
		effect := NewRedirectCombatDamageEffect(source, fromTarget, toTarget, 5, DurationUntilEndOfTurn)

		assert.Equal(t, 5, effect.maxAmount)
		assert.Contains(t, effect.GetDescription(), "up to 5 damage")
	})
}
