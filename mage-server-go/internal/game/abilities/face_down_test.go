package abilities

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFaceDownState tests the FaceDownState structure
func TestFaceDownState(t *testing.T) {
	permanentID := uuid.New()
	actualCard := uuid.New()
	ownerID := uuid.New()
	controllerID := uuid.New()

	t.Run("creates face-down state with correct defaults", func(t *testing.T) {
		fds := NewFaceDownState(permanentID, actualCard, ownerID, controllerID, FaceDownMorph)

		assert.Equal(t, permanentID, fds.permanentID)
		assert.Equal(t, actualCard, fds.actualCard)
		assert.Equal(t, ownerID, fds.ownerID)
		assert.Equal(t, controllerID, fds.controllerID)
		assert.Equal(t, FaceDownMorph, fds.faceDownType)

		// Rule 708.2: Face-down permanents are 2/2 creatures with no characteristics
		assert.True(t, fds.isCreature)
		assert.Equal(t, 2, fds.power)
		assert.Equal(t, 2, fds.toughness)
		assert.True(t, fds.hasNoText)
		assert.True(t, fds.hasNoName)
		assert.True(t, fds.hasNoSubtypes)
		assert.True(t, fds.hasNoManaCost)
		assert.True(t, fds.hasNoColor)
	})

	t.Run("GetPower returns 2", func(t *testing.T) {
		fds := NewFaceDownState(permanentID, actualCard, ownerID, controllerID, FaceDownMorph)
		assert.Equal(t, 2, fds.GetPower())
	})

	t.Run("GetToughness returns 2", func(t *testing.T) {
		fds := NewFaceDownState(permanentID, actualCard, ownerID, controllerID, FaceDownMorph)
		assert.Equal(t, 2, fds.GetToughness())
	})

	t.Run("IsFaceDown always returns true", func(t *testing.T) {
		fds := NewFaceDownState(permanentID, actualCard, ownerID, controllerID, FaceDownMorph)
		assert.True(t, fds.IsFaceDown())
	})

	t.Run("GetFaceDownType returns correct type", func(t *testing.T) {
		testCases := []struct {
			name string
			typ  FaceDownType
		}{
			{"Morph", FaceDownMorph},
			{"Manifest", FaceDownManifest},
			{"Megamorph", FaceDownMegamorph},
			{"Cloak", FaceDownCloak},
			{"Disguise", FaceDownDisguise},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				fds := NewFaceDownState(permanentID, actualCard, ownerID, controllerID, tc.typ)
				assert.Equal(t, tc.typ, fds.GetFaceDownType())
			})
		}
	})
}

// TestMorphAbility tests the Morph keyword ability
func TestMorphAbility(t *testing.T) {
	sourceID := uuid.New()
	cost, _ := ParseManaCost("{2}{U}")

	t.Run("creates morph ability", func(t *testing.T) {
		morph := NewMorphAbility(sourceID, cost)

		require.NotNil(t, morph)
		assert.NotEqual(t, uuid.Nil, morph.id)
		assert.Equal(t, sourceID, morph.sourceID)
		assert.Equal(t, cost, morph.morphCost)
		assert.False(t, morph.isMegamorph)
	})

	t.Run("GetType returns static", func(t *testing.T) {
		morph := NewMorphAbility(sourceID, cost)
		assert.Equal(t, AbilityTypeStatic, morph.GetType())
	})

	t.Run("CanActivate returns true", func(t *testing.T) {
		morph := NewMorphAbility(sourceID, cost)
		ctx := context.Background()
		assert.True(t, morph.CanActivate(ctx, nil))
	})

	t.Run("GetMorphCost returns cost", func(t *testing.T) {
		morph := NewMorphAbility(sourceID, cost)
		assert.Equal(t, cost, morph.GetMorphCost())
	})

	t.Run("IsMegamorph returns false for morph", func(t *testing.T) {
		morph := NewMorphAbility(sourceID, cost)
		assert.False(t, morph.IsMegamorph())
	})

	t.Run("String formats correctly", func(t *testing.T) {
		morph := NewMorphAbility(sourceID, cost)
		expected := "Morph {2}{U}"
		assert.Equal(t, expected, morph.String())
	})
}

// TestMegamorphAbility tests the Megamorph keyword ability
func TestMegamorphAbility(t *testing.T) {
	sourceID := uuid.New()
	cost, _ := ParseManaCost("{1}{G}")

	t.Run("creates megamorph ability", func(t *testing.T) {
		megamorph := NewMegamorphAbility(sourceID, cost)

		require.NotNil(t, megamorph)
		assert.NotEqual(t, uuid.Nil, megamorph.id)
		assert.Equal(t, sourceID, megamorph.sourceID)
		assert.Equal(t, cost, megamorph.morphCost)
		assert.True(t, megamorph.isMegamorph)
	})

	t.Run("IsMegamorph returns true", func(t *testing.T) {
		megamorph := NewMegamorphAbility(sourceID, cost)
		assert.True(t, megamorph.IsMegamorph())
	})

	t.Run("String formats correctly", func(t *testing.T) {
		megamorph := NewMegamorphAbility(sourceID, cost)
		expected := "Megamorph {1}{G}"
		assert.Equal(t, expected, megamorph.String())
	})
}

// TestManifestEffect tests the Manifest effect
func TestManifestEffect(t *testing.T) {
	sourceID := uuid.New()

	t.Run("creates manifest effect from library", func(t *testing.T) {
		effect := NewManifestEffect(sourceID, ManifestFromLibrary, 1)

		require.NotNil(t, effect)
		assert.Equal(t, "Manifest", effect.description)
		assert.Equal(t, sourceID, effect.source)
		assert.Equal(t, ManifestFromLibrary, effect.cardSource)
		assert.Equal(t, 1, effect.count)
	})

	t.Run("creates manifest effect from hand", func(t *testing.T) {
		effect := NewManifestEffect(sourceID, ManifestFromHand, 2)

		assert.Equal(t, ManifestFromHand, effect.cardSource)
		assert.Equal(t, 2, effect.count)
	})

	t.Run("creates manifest effect from graveyard", func(t *testing.T) {
		effect := NewManifestEffect(sourceID, ManifestFromGraveyard, 3)

		assert.Equal(t, ManifestFromGraveyard, effect.cardSource)
		assert.Equal(t, 3, effect.count)
	})

	t.Run("GetDescription formats correctly", func(t *testing.T) {
		effect := NewManifestEffect(sourceID, ManifestFromLibrary, 1)
		expected := "Manifest 1 card(s)"
		assert.Equal(t, expected, effect.GetDescription())
	})

	t.Run("GetDescription with multiple cards", func(t *testing.T) {
		effect := NewManifestEffect(sourceID, ManifestFromLibrary, 5)
		expected := "Manifest 5 card(s)"
		assert.Equal(t, expected, effect.GetDescription())
	})
}

// TestCloakEffect tests the Cloak effect
func TestCloakEffect(t *testing.T) {
	sourceID := uuid.New()
	cost, _ := ParseManaCost("{3}")

	t.Run("creates cloak effect", func(t *testing.T) {
		effect := NewCloakEffect(sourceID, cost)

		require.NotNil(t, effect)
		assert.Equal(t, "Cloak", effect.description)
		assert.Equal(t, sourceID, effect.source)
		assert.Equal(t, cost, effect.cloakCost)
	})

	t.Run("GetDescription formats correctly", func(t *testing.T) {
		effect := NewCloakEffect(sourceID, cost)
		expected := "Cloak {3}"
		assert.Equal(t, expected, effect.GetDescription())
	})
}

// TestTurnFaceUpAction tests the turn face up special action
func TestTurnFaceUpAction(t *testing.T) {
	permanentID := uuid.New()
	playerID := uuid.New()
	cost, _ := ParseManaCost("{2}{U}")

	t.Run("creates turn face up action for morph", func(t *testing.T) {
		action := NewTurnFaceUpAction(permanentID, playerID, cost, false)

		require.NotNil(t, action)
		assert.Equal(t, permanentID, action.permanentID)
		assert.Equal(t, playerID, action.player)
		assert.Equal(t, cost, action.cost)
		assert.False(t, action.isMegamorph)
	})

	t.Run("creates turn face up action for megamorph", func(t *testing.T) {
		action := NewTurnFaceUpAction(permanentID, playerID, cost, true)

		assert.True(t, action.isMegamorph)
	})

	t.Run("GetPermanentID returns correct ID", func(t *testing.T) {
		action := NewTurnFaceUpAction(permanentID, playerID, cost, false)
		assert.Equal(t, permanentID, action.GetPermanentID())
	})

	t.Run("GetCost returns correct cost", func(t *testing.T) {
		action := NewTurnFaceUpAction(permanentID, playerID, cost, false)
		assert.Equal(t, cost, action.GetCost())
	})

	t.Run("IsMegamorph returns correct value", func(t *testing.T) {
		actionMorph := NewTurnFaceUpAction(permanentID, playerID, cost, false)
		actionMegamorph := NewTurnFaceUpAction(permanentID, playerID, cost, true)

		assert.False(t, actionMorph.IsMegamorph())
		assert.True(t, actionMegamorph.IsMegamorph())
	})
}

// TestCastFaceDownOption tests the face-down casting option
func TestCastFaceDownOption(t *testing.T) {
	cardID := uuid.New()
	cost, _ := ParseManaCost("{3}")

	t.Run("creates cast face down option for morph", func(t *testing.T) {
		option := NewCastFaceDownOption(cardID, FaceDownMorph, cost)

		require.NotNil(t, option)
		assert.Equal(t, cardID, option.cardID)
		assert.Equal(t, FaceDownMorph, option.faceDownType)
		assert.Equal(t, cost, option.cost)
	})

	t.Run("creates cast face down option for disguise", func(t *testing.T) {
		option := NewCastFaceDownOption(cardID, FaceDownDisguise, cost)

		assert.Equal(t, FaceDownDisguise, option.faceDownType)
	})

	t.Run("GetCardID returns correct ID", func(t *testing.T) {
		option := NewCastFaceDownOption(cardID, FaceDownMorph, cost)
		assert.Equal(t, cardID, option.GetCardID())
	})

	t.Run("GetFaceDownType returns correct type", func(t *testing.T) {
		option := NewCastFaceDownOption(cardID, FaceDownMorph, cost)
		assert.Equal(t, FaceDownMorph, option.GetFaceDownType())
	})

	t.Run("GetCost returns correct cost", func(t *testing.T) {
		option := NewCastFaceDownOption(cardID, FaceDownMorph, cost)
		assert.Equal(t, cost, option.GetCost())
	})
}

// TestFaceDownTypes tests all face-down type constants
func TestFaceDownTypes(t *testing.T) {
	testCases := []struct {
		name string
		typ  FaceDownType
	}{
		{"FaceDownMorph", FaceDownMorph},
		{"FaceDownManifest", FaceDownManifest},
		{"FaceDownMegamorph", FaceDownMegamorph},
		{"FaceDownCloak", FaceDownCloak},
		{"FaceDownDisguise", FaceDownDisguise},
	}

	t.Run("all types have unique values", func(t *testing.T) {
		seen := make(map[FaceDownType]bool)
		for _, tc := range testCases {
			assert.False(t, seen[tc.typ], "duplicate face-down type: %s", tc.name)
			seen[tc.typ] = true
		}
	})
}

// TestManifestSources tests all manifest source constants
func TestManifestSources(t *testing.T) {
	testCases := []struct {
		name   string
		source ManifestSource
	}{
		{"ManifestFromLibrary", ManifestFromLibrary},
		{"ManifestFromHand", ManifestFromHand},
		{"ManifestFromGraveyard", ManifestFromGraveyard},
	}

	t.Run("all sources have unique values", func(t *testing.T) {
		seen := make(map[ManifestSource]bool)
		for _, tc := range testCases {
			assert.False(t, seen[tc.source], "duplicate manifest source: %s", tc.name)
			seen[tc.source] = true
		}
	})
}

// TestFaceDownStateWithMorphCost tests face-down state with morph cost
func TestFaceDownStateWithMorphCost(t *testing.T) {
	permanentID := uuid.New()
	actualCard := uuid.New()
	ownerID := uuid.New()
	controllerID := uuid.New()
	morphCost, _ := ParseManaCost("{2}{U}")

	t.Run("face-down state can store morph cost", func(t *testing.T) {
		fds := NewFaceDownState(permanentID, actualCard, ownerID, controllerID, FaceDownMorph)
		fds.morphCost = morphCost
		fds.canTurnFaceUp = true

		assert.Equal(t, morphCost, fds.morphCost)
		assert.True(t, fds.canTurnFaceUp)
	})

	t.Run("face-down state can store megamorph flag", func(t *testing.T) {
		fds := NewFaceDownState(permanentID, actualCard, ownerID, controllerID, FaceDownMegamorph)
		fds.morphCost = morphCost
		fds.isMegamorph = true
		fds.canTurnFaceUp = true

		assert.True(t, fds.isMegamorph)
		assert.True(t, fds.canTurnFaceUp)
	})
}

// TestFaceDownStateCharacteristics tests Rule 708.2 compliance
func TestFaceDownStateCharacteristics(t *testing.T) {
	permanentID := uuid.New()
	actualCard := uuid.New()
	ownerID := uuid.New()
	controllerID := uuid.New()

	t.Run("Rule 708.2: face-down permanent has only specified characteristics", func(t *testing.T) {
		fds := NewFaceDownState(permanentID, actualCard, ownerID, controllerID, FaceDownMorph)

		// Rule 708.2a: 2/2 creature
		assert.Equal(t, 2, fds.GetPower(), "must be 2/2")
		assert.Equal(t, 2, fds.GetToughness(), "must be 2/2")
		assert.True(t, fds.isCreature, "must be a creature")

		// Rule 708.2b: No text, name, subtypes, mana cost
		assert.True(t, fds.hasNoText, "must have no text")
		assert.True(t, fds.hasNoName, "must have no name")
		assert.True(t, fds.hasNoSubtypes, "must have no subtypes")
		assert.True(t, fds.hasNoManaCost, "must have no mana cost")

		// Rule 708.2c: Colorless
		assert.True(t, fds.hasNoColor, "must be colorless")
	})
}

// Example test showing integration with game context
func TestFaceDownIntegration(t *testing.T) {
	t.Run("example workflow: casting with morph", func(t *testing.T) {
		// 1. Card has morph ability
		cardID := uuid.New()
		morphCost, _ := ParseManaCost("{2}{U}")
		morphAbility := NewMorphAbility(cardID, morphCost)

		assert.NotNil(t, morphAbility)
		assert.Equal(t, morphCost, morphAbility.GetMorphCost())

		// 2. Player chooses to cast face down
		faceDownCost, _ := ParseManaCost("{3}")
		castOption := NewCastFaceDownOption(cardID, FaceDownMorph, faceDownCost)

		assert.Equal(t, cardID, castOption.GetCardID())
		assert.Equal(t, FaceDownMorph, castOption.GetFaceDownType())

		// 3. Card enters battlefield face down
		permanentID := uuid.New()
		ownerID := uuid.New()
		controllerID := uuid.New()
		faceDownState := NewFaceDownState(permanentID, cardID, ownerID, controllerID, FaceDownMorph)
		faceDownState.morphCost = morphCost
		faceDownState.canTurnFaceUp = true

		assert.True(t, faceDownState.IsFaceDown())
		assert.Equal(t, 2, faceDownState.GetPower())
		assert.Equal(t, 2, faceDownState.GetToughness())
		assert.True(t, faceDownState.CanTurnFaceUp())

		// 4. Later, player turns it face up
		turnFaceUp := NewTurnFaceUpAction(permanentID, ownerID, morphCost, false)

		assert.Equal(t, permanentID, turnFaceUp.GetPermanentID())
		assert.Equal(t, morphCost, turnFaceUp.GetCost())
		assert.False(t, turnFaceUp.IsMegamorph())
	})

	t.Run("example workflow: manifest from library", func(t *testing.T) {
		// 1. Spell manifests top card
		sourceID := uuid.New()
		manifestEffect := NewManifestEffect(sourceID, ManifestFromLibrary, 1)

		assert.NotNil(t, manifestEffect)
		assert.Equal(t, ManifestFromLibrary, manifestEffect.cardSource)
		assert.Equal(t, 1, manifestEffect.count)

		// 2. Card enters battlefield face down
		permanentID := uuid.New()
		cardID := uuid.New() // The actual card from library (hidden)
		ownerID := uuid.New()
		controllerID := uuid.New()
		faceDownState := NewFaceDownState(permanentID, cardID, ownerID, controllerID, FaceDownManifest)

		assert.True(t, faceDownState.IsFaceDown())
		assert.Equal(t, FaceDownManifest, faceDownState.GetFaceDownType())

		// 3. If it's a creature card, it can be turned face up by paying mana cost
		// If it has morph, it can be turned face up by paying morph cost
		// If it's not a creature, it stays face down
	})

	t.Run("example workflow: megamorph", func(t *testing.T) {
		// 1. Card has megamorph ability
		cardID := uuid.New()
		megamorphCost, _ := ParseManaCost("{1}{G}")
		megamorphAbility := NewMegamorphAbility(cardID, megamorphCost)

		assert.True(t, megamorphAbility.IsMegamorph())

		// 2. Cast face down
		permanentID := uuid.New()
		ownerID := uuid.New()
		controllerID := uuid.New()
		faceDownState := NewFaceDownState(permanentID, cardID, ownerID, controllerID, FaceDownMegamorph)
		faceDownState.morphCost = megamorphCost
		faceDownState.isMegamorph = true
		faceDownState.canTurnFaceUp = true

		assert.True(t, faceDownState.IsMegamorph())

		// 3. Turn face up (will add +1/+1 counter)
		turnFaceUp := NewTurnFaceUpAction(permanentID, ownerID, megamorphCost, true)

		assert.True(t, turnFaceUp.IsMegamorph())
		// When resolved, this should add a +1/+1 counter to the permanent
	})
}
