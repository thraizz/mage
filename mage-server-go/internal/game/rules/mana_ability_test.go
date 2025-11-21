package rules

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManaAbilityActivationContext(t *testing.T) {
	t.Run("basic activation tracking", func(t *testing.T) {
		maac := NewManaAbilityActivationContext()

		abilityID := "forest-tap"

		// Can activate initially
		assert.True(t, maac.CanActivate(abilityID))
		assert.False(t, maac.IsActivating(abilityID))

		// Begin activation
		err := maac.BeginActivation(abilityID)
		require.NoError(t, err)

		// Cannot activate again while activating (Rule 605.3c)
		assert.False(t, maac.CanActivate(abilityID))
		assert.True(t, maac.IsActivating(abilityID))

		// Trying to activate again should fail
		err = maac.BeginActivation(abilityID)
		assert.Error(t, err)

		// End activation
		maac.EndActivation(abilityID)

		// Can activate again
		assert.True(t, maac.CanActivate(abilityID))
		assert.False(t, maac.IsActivating(abilityID))
	})

	t.Run("resolution count tracking", func(t *testing.T) {
		maac := NewManaAbilityActivationContext()

		abilityID := "forest-tap"

		// No resolutions initially
		assert.Equal(t, 0, maac.GetResolutionCount(abilityID))

		// Activate and resolve
		maac.BeginActivation(abilityID)
		maac.EndActivation(abilityID)
		assert.Equal(t, 1, maac.GetResolutionCount(abilityID))

		// Activate again
		maac.BeginActivation(abilityID)
		maac.EndActivation(abilityID)
		assert.Equal(t, 2, maac.GetResolutionCount(abilityID))

		// Reset window
		maac.ResetWindow()
		assert.Equal(t, 0, maac.GetResolutionCount(abilityID))
	})

	t.Run("multiple abilities", func(t *testing.T) {
		maac := NewManaAbilityActivationContext()

		ability1 := "forest-tap"
		ability2 := "mountain-tap"

		// Activate both
		maac.BeginActivation(ability1)
		maac.BeginActivation(ability2)

		assert.True(t, maac.IsActivating(ability1))
		assert.True(t, maac.IsActivating(ability2))

		// End both
		maac.EndActivation(ability1)
		maac.EndActivation(ability2)

		assert.False(t, maac.IsActivating(ability1))
		assert.False(t, maac.IsActivating(ability2))
		assert.Equal(t, 1, maac.GetResolutionCount(ability1))
		assert.Equal(t, 1, maac.GetResolutionCount(ability2))
	})
}

func TestManaAbilityManager(t *testing.T) {
	t.Run("activate mana ability - immediate resolution", func(t *testing.T) {
		mam := NewManaAbilityManager()

		activated := false
		ability := ManaAbility{
			ID:           "forest-tap",
			SourceID:     "forest-1",
			ControllerID: "player-1",
			Text:         "{T}: Add {G}",
			Activate: func() error {
				activated = true
				return nil
			},
		}

		// Activate the ability (Rule 605.3b: resolves immediately)
		err := mam.ActivateManaAbility(ability)
		require.NoError(t, err)
		assert.True(t, activated)

		// Can activate again (it has resolved)
		activated = false
		err = mam.ActivateManaAbility(ability)
		require.NoError(t, err)
		assert.True(t, activated)
	})

	t.Run("cannot reactivate while activating", func(t *testing.T) {
		mam := NewManaAbilityManager()

		// Create ability that tries to activate itself
		var ability ManaAbility
		ability = ManaAbility{
			ID:           "infinite-mana",
			SourceID:     "lotus-1",
			ControllerID: "player-1",
			Text:         "{T}: Add {C}",
			Activate: func() error {
				// Try to activate again from within activation
				// This should fail per Rule 605.3c
				err := mam.ActivateManaAbility(ability)
				if err != nil {
					return err
				}
				return nil
			},
		}

		// First activation should fail because inner activation fails
		err := mam.ActivateManaAbility(ability)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already activating")
	})

	t.Run("triggered mana ability - immediate resolution", func(t *testing.T) {
		mam := NewManaAbilityManager()

		resolved := false
		triggered := TriggeredManaAbility{
			ID:           "triggered-1",
			SourceID:     "lotus-bloom",
			ControllerID: "player-1",
			TriggerID:    "suspend-trigger",
			Text:         "Add three mana of any one color",
			Resolve: func() error {
				resolved = true
				return nil
			},
		}

		// Queue the triggered ability
		mam.QueueTriggeredManaAbility(triggered)
		assert.True(t, mam.HasPendingTriggeredAbilities())

		// Resolve all triggered abilities (Rule 605.4a: immediately)
		err := mam.ResolveTriggeredManaAbilities()
		require.NoError(t, err)
		assert.True(t, resolved)
		assert.False(t, mam.HasPendingTriggeredAbilities())
	})

	t.Run("cascading triggered mana abilities", func(t *testing.T) {
		mam := NewManaAbilityManager()

		count := 0

		// First triggered ability
		triggered1 := TriggeredManaAbility{
			ID:           "triggered-1",
			SourceID:     "source-1",
			ControllerID: "player-1",
			Resolve: func() error {
				count++
				// Trigger another mana ability
				mam.QueueTriggeredManaAbility(TriggeredManaAbility{
					ID:           "triggered-2",
					SourceID:     "source-2",
					ControllerID: "player-1",
					Resolve: func() error {
						count++
						return nil
					},
				})
				return nil
			},
		}

		mam.QueueTriggeredManaAbility(triggered1)

		// Resolve all (should handle cascade)
		err := mam.ResolveTriggeredManaAbilities()
		require.NoError(t, err)
		assert.Equal(t, 2, count)
		assert.False(t, mam.HasPendingTriggeredAbilities())
	})

	t.Run("activation permission during cast and resolve", func(t *testing.T) {
		mam := NewManaAbilityManager()

		// Can activate during cast by default (Rule 117.1d)
		assert.True(t, mam.CanActivate(true, false))

		// Cannot activate during resolve by default
		assert.False(t, mam.CanActivate(false, true))

		// Enable during resolve (for payment window)
		mam.SetCanActivateDuringResolve(true)
		assert.True(t, mam.CanActivate(false, true))

		// Disable during cast (unusual but possible)
		mam.SetCanActivateDuringCast(false)
		assert.False(t, mam.CanActivate(true, false))
	})
}

func TestManaAbilityIntegrationWithPayment(t *testing.T) {
	t.Run("activate mana abilities during payment window", func(t *testing.T) {
		mam := NewManaAbilityManager()
		rc := NewResolutionContext()

		manaPool := 0

		// Create Forest tap ability
		forestTap := ManaAbility{
			ID:           "forest-tap",
			SourceID:     "forest-1",
			ControllerID: "player-1",
			Text:         "{T}: Add {G}",
			Activate: func() error {
				manaPool++
				return nil
			},
		}

		// Start resolving (casting) Lightning Bolt
		rc.BeginResolution("lightning-bolt")
		rc.SetAllowManaAbilities(true)

		// Activate mana ability during payment
		mam.SetCanActivateDuringResolve(true)
		err := mam.ActivateManaAbility(forestTap)
		require.NoError(t, err)
		assert.Equal(t, 1, manaPool)

		// End resolution
		rc.EndResolution("lightning-bolt")
		assert.Equal(t, 1, manaPool)
	})

	t.Run("mana ability with triggered mana ability", func(t *testing.T) {
		mam := NewManaAbilityManager()

		totalMana := 0

		// Main ability that triggers another
		mainAbility := ManaAbility{
			ID:           "nykthos-tap",
			SourceID:     "nykthos",
			ControllerID: "player-1",
			Text:         "{T}: Add mana",
			Activate: func() error {
				totalMana += 3

				// Trigger another mana ability
				// (simulating devotion triggering additional mana)
				mam.QueueTriggeredManaAbility(TriggeredManaAbility{
					ID:           "devotion-trigger",
					SourceID:     "nykthos",
					ControllerID: "player-1",
					TriggerID:    "nykthos-tap",
					Resolve: func() error {
						totalMana += 2
						return nil
					},
				})
				return nil
			},
		}

		// Activate main ability
		err := mam.ActivateManaAbility(mainAbility)
		require.NoError(t, err)
		assert.Equal(t, 3, totalMana)
		assert.True(t, mam.HasPendingTriggeredAbilities())

		// Resolve triggered abilities (Rule 605.4a: immediately)
		err = mam.ResolveTriggeredManaAbilities()
		require.NoError(t, err)
		assert.Equal(t, 5, totalMana)
	})

	t.Run("error during mana ability activation", func(t *testing.T) {
		mam := NewManaAbilityManager()

		errorAbility := ManaAbility{
			ID:           "broken-land",
			SourceID:     "land-1",
			ControllerID: "player-1",
			Activate: func() error {
				return fmt.Errorf("cannot produce mana")
			},
		}

		err := mam.ActivateManaAbility(errorAbility)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to activate")

		// Should be able to try again (error means it didn't start activating)
		err = mam.ActivateManaAbility(errorAbility)
		assert.Error(t, err)
	})
}
