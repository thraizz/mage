package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentState(t *testing.T) {
	t.Run("basic payment tracking", func(t *testing.T) {
		costs := []Cost{
			{Type: CostTypeMana, Amount: 3, Description: "{3}"},
			{Type: CostTypeTap, Description: "{T}"},
		}

		ps := NewPaymentState("lightning-bolt", "player-1", costs)

		// Initial state
		assert.Equal(t, PaymentStepBefore, ps.GetCurrentStep())
		assert.True(t, ps.CanActivateManaAbilities())
		assert.False(t, ps.IsFullyPaid())

		// Move to normal payment
		ps.SetCurrentStep(PaymentStepNormal)
		assert.Equal(t, PaymentStepNormal, ps.GetCurrentStep())

		// Pay mana
		ps.SetManaRemaining(3)
		ps.AddManaPaid(3)
		assert.Equal(t, 3, ps.GetTotalManaPaid())
		assert.Equal(t, 0, ps.GetManaRemaining())

		// Mark costs as paid
		ps.MarkCostPaid(CostTypeMana)
		ps.MarkCostPaid(CostTypeTap)

		assert.True(t, ps.IsCostPaid(CostTypeMana))
		assert.True(t, ps.IsCostPaid(CostTypeTap))
	})

	t.Run("payment step progression", func(t *testing.T) {
		costs := []Cost{
			{Type: CostTypeMana, Amount: 5, Description: "{3}{G}{G}"},
		}

		ps := NewPaymentState("convoke-spell", "player-1", costs)

		// Start at BEFORE (for Convoke)
		assert.Equal(t, PaymentStepBefore, ps.GetCurrentStep())
		assert.True(t, ps.CanActivateManaAbilities())

		// After using Convoke, move to AFTER
		ps.SetCurrentStep(PaymentStepAfter)

		// Per Java implementation: normal mana abilities blocked after special payment
		assert.False(t, ps.CanActivateManaAbilities())
	})

	t.Run("partial payment tracking", func(t *testing.T) {
		costs := []Cost{
			{Type: CostTypeMana, Amount: 5},
		}

		ps := NewPaymentState("spell", "player-1", costs)
		ps.SetManaRemaining(5)

		// Pay 3 mana
		ps.AddManaPaid(3)
		assert.Equal(t, 3, ps.GetTotalManaPaid())
		assert.Equal(t, 2, ps.GetManaRemaining())
		assert.False(t, ps.IsFullyPaid())

		// Pay remaining 2 mana
		ps.AddManaPaid(2)
		assert.Equal(t, 5, ps.GetTotalManaPaid())
		assert.Equal(t, 0, ps.GetManaRemaining())
	})
}

func TestPaymentWindowManager(t *testing.T) {
	t.Run("open and close payment window", func(t *testing.T) {
		pwm := NewPaymentWindowManager()

		costs := []Cost{{Type: CostTypeMana, Amount: 1}}
		state := NewPaymentState("bolt", "player-1", costs)

		assert.False(t, pwm.IsPaymentInProgress())

		err := pwm.BeginPayment(state)
		require.NoError(t, err)
		assert.True(t, pwm.IsPaymentInProgress())

		active := pwm.GetActivePayment()
		require.NotNil(t, active)
		assert.Equal(t, "bolt", active.spellOrAbilityID)

		err = pwm.EndPayment("bolt")
		require.NoError(t, err)
		assert.False(t, pwm.IsPaymentInProgress())
	})

	t.Run("cannot open multiple payment windows", func(t *testing.T) {
		pwm := NewPaymentWindowManager()

		state1 := NewPaymentState("spell-1", "player-1", nil)
		state2 := NewPaymentState("spell-2", "player-1", nil)

		err := pwm.BeginPayment(state1)
		require.NoError(t, err)

		err = pwm.BeginPayment(state2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already in progress")
	})

	t.Run("payment mismatch error", func(t *testing.T) {
		pwm := NewPaymentWindowManager()

		state := NewPaymentState("spell-1", "player-1", nil)
		pwm.BeginPayment(state)

		err := pwm.EndPayment("spell-2")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mismatch")
	})
}

func TestChoiceManager(t *testing.T) {
	t.Run("make choice during resolution", func(t *testing.T) {
		cm := NewChoiceManager()

		choice := Choice{
			Type:       ChoiceTypeMode,
			PlayerID:   "player-1",
			Prompt:     "Choose a mode",
			Options:    []string{"draw 2", "deal 2"},
			MinChoices: 1,
			MaxChoices: 1,
		}

		cm.AddChoice(choice)
		assert.True(t, cm.HasPendingChoices())

		// Get the choice
		next := cm.GetNextChoice()
		require.NotNil(t, next)
		assert.Equal(t, ChoiceTypeMode, next.Type)

		// Make the choice
		err := cm.MakeChoice([]string{"draw 2"})
		require.NoError(t, err)
		assert.False(t, cm.HasPendingChoices())
	})

	t.Run("choice validation - too few", func(t *testing.T) {
		cm := NewChoiceManager()

		choice := Choice{
			Type:       ChoiceTypeCard,
			PlayerID:   "player-1",
			Prompt:     "Choose 2 cards",
			MinChoices: 2,
			MaxChoices: 2,
		}

		cm.AddChoice(choice)
		cm.GetNextChoice()

		// Try to make only 1 choice
		err := cm.MakeChoice([]string{"card-1"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too few")
	})

	t.Run("choice validation - too many", func(t *testing.T) {
		cm := NewChoiceManager()

		choice := Choice{
			Type:       ChoiceTypeCard,
			PlayerID:   "player-1",
			Prompt:     "Choose up to 1 card",
			MinChoices: 0,
			MaxChoices: 1,
		}

		cm.AddChoice(choice)
		cm.GetNextChoice()

		// Try to make 2 choices
		err := cm.MakeChoice([]string{"card-1", "card-2"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too many")
	})

	t.Run("multiple choices in sequence", func(t *testing.T) {
		cm := NewChoiceManager()

		choice1 := Choice{
			Type:       ChoiceTypeMode,
			PlayerID:   "player-1",
			Prompt:     "Choose mode",
			MinChoices: 1,
			MaxChoices: 1,
		}

		choice2 := Choice{
			Type:       ChoiceTypeColor,
			PlayerID:   "player-1",
			Prompt:     "Choose color",
			MinChoices: 1,
			MaxChoices: 1,
		}

		cm.AddChoice(choice1)
		cm.AddChoice(choice2)

		// Make first choice
		cm.GetNextChoice()
		err := cm.MakeChoice([]string{"mode-1"})
		require.NoError(t, err)

		// Make second choice
		cm.GetNextChoice()
		err = cm.MakeChoice([]string{"red"})
		require.NoError(t, err)

		assert.False(t, cm.HasPendingChoices())
	})

	t.Run("no choice in progress error", func(t *testing.T) {
		cm := NewChoiceManager()

		err := cm.MakeChoice([]string{"option"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no choice in progress")
	})
}

func TestPaymentIntegrationWithManaAbilities(t *testing.T) {
	t.Run("pay for spell with mana abilities", func(t *testing.T) {
		pwm := NewPaymentWindowManager()
		mam := NewManaAbilityManager()

		// Create payment state for Lightning Bolt (R)
		costs := []Cost{
			{Type: CostTypeMana, Amount: 1, Description: "{R}"},
		}
		paymentState := NewPaymentState("lightning-bolt", "player-1", costs)
		paymentState.SetManaRemaining(1)

		// Begin payment
		err := pwm.BeginPayment(paymentState)
		require.NoError(t, err)

		// Activate mana ability during payment (Rule 117.1d)
		mam.SetCanActivateDuringResolve(true)

		paid := 0
		forestTap := ManaAbility{
			ID:           "mountain-tap",
			SourceID:     "mountain-1",
			ControllerID: "player-1",
			Text:         "{T}: Add {R}",
			Activate: func() error {
				paid++
				return nil
			},
		}

		err = mam.ActivateManaAbility(forestTap)
		require.NoError(t, err)
		assert.Equal(t, 1, paid)

		// Mark payment complete
		paymentState.AddManaPaid(1)
		paymentState.MarkCostPaid(CostTypeMana)

		// End payment
		err = pwm.EndPayment("lightning-bolt")
		require.NoError(t, err)

		assert.True(t, paymentState.IsFullyPaid())
	})

	t.Run("convoke blocks normal mana abilities", func(t *testing.T) {
		// This simulates the Java implementation where using Convoke
		// (or other special mana payments) blocks normal mana abilities

		paymentState := NewPaymentState("convoke-spell", "player-1", nil)

		// Start in BEFORE step (can use special mana like Convoke)
		assert.Equal(t, PaymentStepBefore, paymentState.GetCurrentStep())
		assert.True(t, paymentState.CanActivateManaAbilities())

		// Use Convoke (special mana payment)
		// This moves to AFTER step
		paymentState.SetCurrentStep(PaymentStepAfter)

		// Normal mana abilities are now blocked
		assert.False(t, paymentState.CanActivateManaAbilities())
	})
}

func TestResolutionChoices(t *testing.T) {
	t.Run("APNAP order for multiplayer choices", func(t *testing.T) {
		// Per Rule 608.2e: Multi-player choices are made in APNAP order
		// This test documents the expected behavior

		cm := NewChoiceManager()

		// Active player choice
		choice1 := Choice{
			Type:       ChoiceTypeCard,
			PlayerID:   "active-player",
			Prompt:     "Choose a card to discard",
			MinChoices: 1,
			MaxChoices: 1,
		}

		// Non-active player choice
		choice2 := Choice{
			Type:       ChoiceTypeCard,
			PlayerID:   "non-active-player",
			Prompt:     "Choose a card to discard",
			MinChoices: 1,
			MaxChoices: 1,
		}

		// Choices should be queued in APNAP order
		cm.AddChoice(choice1) // Active player first
		cm.AddChoice(choice2) // Non-active player second

		// Process active player choice
		next := cm.GetNextChoice()
		require.NotNil(t, next)
		assert.Equal(t, "active-player", next.PlayerID)
		cm.MakeChoice([]string{"card-1"})

		// Process non-active player choice
		next = cm.GetNextChoice()
		require.NotNil(t, next)
		assert.Equal(t, "non-active-player", next.PlayerID)
		cm.MakeChoice([]string{"card-2"})
	})
}
