package abilities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCombatRequirements tests the combat requirements tracker
func TestCombatRequirements(t *testing.T) {
	cr := NewCombatRequirements()

	creature1 := uuid.New()
	creature2 := uuid.New()
	creature3 := uuid.New()

	t.Run("Add and check attack restrictions", func(t *testing.T) {
		cr.AddAttackRestriction(creature1, CantAttackAny)

		assert.False(t, cr.CanAttack(creature1), "Creature should not be able to attack")
		assert.True(t, cr.CanAttack(creature2), "Unrestricted creature should be able to attack")
	})

	t.Run("Add and check attack requirements", func(t *testing.T) {
		cr.AddAttackRequirement(creature2, MustAttackIfAble)

		assert.True(t, cr.MustAttackCheck(creature2), "Creature must attack")
		assert.False(t, cr.MustAttackCheck(creature3), "Creature has no attack requirement")
	})

	t.Run("Add and check block restrictions", func(t *testing.T) {
		cr.AddBlockRestriction(creature1, CantBlockAny)

		assert.False(t, cr.CanBlock(creature1), "Creature should not be able to block")
		assert.True(t, cr.CanBlock(creature2), "Unrestricted creature should be able to block")
	})

	t.Run("Add and check evasion", func(t *testing.T) {
		cr.AddEvasion(creature3, CantBeBlockedAtAll)

		assert.False(t, cr.CanBeBlocked(creature3), "Creature should be unblockable")
		assert.True(t, cr.CanBeBlocked(creature1), "Normal creature can be blocked")
	})

	t.Run("Goad creature", func(t *testing.T) {
		player := uuid.New()
		cr.GoadCreature(creature1, player)

		assert.True(t, cr.IsGoaded(creature1), "Creature should be goaded")
		assert.Equal(t, player, cr.GetGoadingPlayer(creature1))
		assert.False(t, cr.IsGoaded(creature2), "Other creature not goaded")
	})

	t.Run("Clear requirements", func(t *testing.T) {
		cr.Clear()

		// After clear, goad persists but other effects don't
		assert.False(t, cr.CanAttack(creature1) == false, "Restrictions cleared")
		assert.False(t, cr.MustAttackCheck(creature2), "Requirements cleared")
	})
}

// TestValidateAttackDeclaration tests attack declaration validation
func TestValidateAttackDeclaration(t *testing.T) {
	cr := NewCombatRequirements()
	activePlayer := uuid.New()

	creature1 := uuid.New()
	creature2 := uuid.New()
	creature3 := uuid.New()

	t.Run("Must attack - violation", func(t *testing.T) {
		cr.AddAttackRequirement(creature1, MustAttackIfAble)

		// Declare attackers without creature1
		violations := cr.ValidateAttackDeclaration([]uuid.UUID{creature2}, activePlayer)

		require.Len(t, violations, 1)
		assert.Contains(t, violations[0], "must attack")
	})

	t.Run("Must attack - satisfied", func(t *testing.T) {
		cr.AddAttackRequirement(creature1, MustAttackIfAble)

		// Declare attackers including creature1
		violations := cr.ValidateAttackDeclaration([]uuid.UUID{creature1, creature2}, activePlayer)

		assert.Len(t, violations, 0, "No violations when must-attack creature attacks")
	})

	t.Run("Goad - violation", func(t *testing.T) {
		goadingPlayer := uuid.New()
		cr.GoadCreature(creature2, goadingPlayer)

		// creature2 not attacking (but should attack someone other than goading player)
		violations := cr.ValidateAttackDeclaration([]uuid.UUID{creature1}, activePlayer)

		require.Len(t, violations, 1)
		assert.Contains(t, violations[0], "Goaded")
	})

	t.Run("Multiple violations", func(t *testing.T) {
		cr.AddAttackRequirement(creature1, MustAttackIfAble)
		cr.AddAttackRequirement(creature3, MustAttackIfAble)

		// Neither attacking
		violations := cr.ValidateAttackDeclaration([]uuid.UUID{creature2}, activePlayer)

		assert.GreaterOrEqual(t, len(violations), 2, "Multiple violations detected")
	})
}

// TestValidateBlockDeclaration tests block declaration validation
func TestValidateBlockDeclaration(t *testing.T) {
	cr := NewCombatRequirements()

	blocker1 := uuid.New()
	blocker2 := uuid.New()
	attacker1 := uuid.New()

	t.Run("Must block - violation", func(t *testing.T) {
		cr.AddBlockRequirement(blocker1, MustBlockIfAble, attacker1)

		// blocker1 not blocking
		blocks := map[uuid.UUID][]uuid.UUID{
			attacker1: {blocker2},
		}

		violations := cr.ValidateBlockDeclaration(blocks)

		require.Len(t, violations, 1)
		assert.Contains(t, violations[0], "must block")
	})

	t.Run("Must block - satisfied", func(t *testing.T) {
		cr.AddBlockRequirement(blocker1, MustBlockIfAble, attacker1)

		// blocker1 blocking
		blocks := map[uuid.UUID][]uuid.UUID{
			attacker1: {blocker1, blocker2},
		}

		violations := cr.ValidateBlockDeclaration(blocks)

		assert.Len(t, violations, 0, "No violations when must-block creature blocks")
	})

	t.Run("No blocks", func(t *testing.T) {
		cr.AddBlockRequirement(blocker1, MustBlockIfAble, attacker1)

		blocks := map[uuid.UUID][]uuid.UUID{}

		violations := cr.ValidateBlockDeclaration(blocks)

		require.Len(t, violations, 1, "Violation when must-block creature doesn't block")
	})
}

// TestCantAttackEffect tests attack restriction effects
func TestCantAttackEffect(t *testing.T) {
	source := uuid.New()
	creature := uuid.New()

	t.Run("Can't attack any", func(t *testing.T) {
		effect := NewCantAttackEffect(source, creature, CantAttackAny, DurationPermanent)

		assert.Equal(t, CantAttackAny, effect.restriction)
		assert.Equal(t, creature, effect.targetCreature)
		assert.Contains(t, effect.GetDescription(), "Can't attack")
	})

	t.Run("Can't attack alone", func(t *testing.T) {
		effect := NewCantAttackEffect(source, creature, CantAttackAlone, DurationUntilEndOfTurn)

		assert.Equal(t, CantAttackAlone, effect.restriction)
		assert.Contains(t, effect.GetDescription(), "alone")
	})

	t.Run("Can't attack player", func(t *testing.T) {
		effect := NewCantAttackEffect(source, creature, CantAttackPlayer, DurationPermanent)

		assert.Equal(t, CantAttackPlayer, effect.restriction)
		assert.Contains(t, effect.GetDescription(), "players")
	})
}

// TestMustAttackEffect tests attack requirement effects
func TestMustAttackEffect(t *testing.T) {
	source := uuid.New()
	creature := uuid.New()

	t.Run("Must attack if able", func(t *testing.T) {
		effect := NewMustAttackEffect(source, creature, MustAttackIfAble, DurationPermanent)

		assert.Equal(t, MustAttackIfAble, effect.requirement)
		assert.Equal(t, creature, effect.targetCreature)
		assert.Contains(t, effect.GetDescription(), "Attacks each combat if able")
	})

	t.Run("Must attack each turn (goad)", func(t *testing.T) {
		effect := NewMustAttackEffect(source, creature, MustAttackEachTurn, DurationUntilYourNextTurn)

		assert.Equal(t, MustAttackEachTurn, effect.requirement)
		assert.Contains(t, effect.GetDescription(), "attacks a player other than you")
	})
}

// TestCantBlockEffect tests block restriction effects
func TestCantBlockEffect(t *testing.T) {
	source := uuid.New()
	creature := uuid.New()

	t.Run("Can't block any", func(t *testing.T) {
		effect := NewCantBlockEffect(source, creature, CantBlockAny, DurationPermanent)

		assert.Equal(t, CantBlockAny, effect.restriction)
		assert.Equal(t, creature, effect.targetCreature)
		assert.Contains(t, effect.GetDescription(), "Can't block")
	})

	t.Run("Can't block flying", func(t *testing.T) {
		effect := NewCantBlockEffect(source, creature, CantBlockFlying, DurationPermanent)

		assert.Equal(t, CantBlockFlying, effect.restriction)
		assert.Contains(t, effect.GetDescription(), "flying")
	})

	t.Run("Can block only one", func(t *testing.T) {
		effect := NewCantBlockEffect(source, creature, CantBlockMoreThanOneCreature, DurationPermanent)

		assert.Equal(t, CantBlockMoreThanOneCreature, effect.restriction)
		assert.Contains(t, effect.GetDescription(), "only one creature")
	})
}

// TestMustBlockEffect tests block requirement effects
func TestMustBlockEffect(t *testing.T) {
	source := uuid.New()
	creature := uuid.New()
	attacker := uuid.New()

	t.Run("Must block if able", func(t *testing.T) {
		effect := NewMustBlockEffect(source, creature, uuid.Nil, MustBlockIfAble, DurationUntilEndOfTurn)

		assert.Equal(t, MustBlockIfAble, effect.requirement)
		assert.Equal(t, creature, effect.targetCreature)
		assert.Contains(t, effect.GetDescription(), "Blocks if able")
	})

	t.Run("Must block specific attacker", func(t *testing.T) {
		effect := NewMustBlockEffect(source, creature, attacker, MustBlockAttacker, DurationUntilEndOfTurn)

		assert.Equal(t, MustBlockAttacker, effect.requirement)
		assert.Equal(t, attacker, effect.targetAttacker)
		assert.Contains(t, effect.GetDescription(), "Must block attacker")
	})

	t.Run("Must be blocked by all", func(t *testing.T) {
		effect := NewMustBlockEffect(source, creature, uuid.Nil, MustBlockWithAllCreatures, DurationPermanent)

		assert.Equal(t, MustBlockWithAllCreatures, effect.requirement)
		assert.Contains(t, effect.GetDescription(), "Must be blocked by all")
	})
}

// TestCantBeBlockedEffect tests evasion effects
func TestCantBeBlockedEffect(t *testing.T) {
	source := uuid.New()
	attacker := uuid.New()

	t.Run("Unblockable", func(t *testing.T) {
		effect := NewCantBeBlockedEffect(source, attacker, CantBeBlockedAtAll, DurationUntilEndOfTurn)

		assert.Equal(t, CantBeBlockedAtAll, effect.condition)
		assert.Equal(t, attacker, effect.attacker)
		assert.Contains(t, effect.GetDescription(), "Can't be blocked")
	})

	t.Run("Can't be blocked except by", func(t *testing.T) {
		effect := NewCantBeBlockedEffect(source, attacker, CantBeBlockedExceptBy, DurationPermanent)

		assert.Equal(t, CantBeBlockedExceptBy, effect.condition)
		assert.Contains(t, effect.GetDescription(), "except by specific")
	})

	t.Run("Can't be blocked by more than X", func(t *testing.T) {
		effect := NewCantBeBlockedEffect(source, attacker, CantBeBlockedByMoreThan, DurationUntilEndOfTurn)

		assert.Equal(t, CantBeBlockedByMoreThan, effect.condition)
		assert.Contains(t, effect.GetDescription(), "more than X")
	})
}

// TestProvokeAbility tests Provoke keyword
func TestProvokeAbility(t *testing.T) {
	source := uuid.New()

	ability := NewProvokeAbility(source)

	assert.Equal(t, AbilityTypeTriggered, ability.GetType())
	assert.Equal(t, "Provoke", ability.String())
	assert.False(t, ability.CanActivate(nil, nil), "Triggered abilities can't be activated")
}

// TestGoadAbility tests Goad keyword
func TestGoadAbility(t *testing.T) {
	source := uuid.New()
	goadingPlayer := uuid.New()

	ability := NewGoadAbility(source, goadingPlayer)

	assert.Equal(t, AbilityTypeStatic, ability.GetType())
	assert.Equal(t, "Goad", ability.String())
	assert.Equal(t, goadingPlayer, ability.goadingPlayer)
	assert.NotNil(t, ability.goadedCreatures)
}

// TestCombatRequirements_ComplexScenarios tests complex interaction scenarios
func TestCombatRequirements_ComplexScenarios(t *testing.T) {
	t.Run("Multiple restrictions on same creature", func(t *testing.T) {
		cr := NewCombatRequirements()
		creature := uuid.New()

		// Add multiple restrictions
		cr.AddAttackRestriction(creature, CantAttackPlayer)
		cr.AddAttackRestriction(creature, CantAttackAlone)

		// Creature has restrictions but not CantAttackAny
		assert.True(t, cr.CanAttack(creature), "Can attack (just not players, and not alone)")
	})

	t.Run("Conflicting requirements", func(t *testing.T) {
		cr := NewCombatRequirements()
		creature := uuid.New()

		// Both must attack and can't attack (rules issue)
		cr.AddAttackRestriction(creature, CantAttackAny)
		cr.AddAttackRequirement(creature, MustAttackIfAble)

		// Can't attack wins (restrictions override requirements per rules)
		assert.False(t, cr.CanAttack(creature))
		assert.True(t, cr.MustAttackCheck(creature), "Still has requirement (creates violation)")
	})

	t.Run("Multiple goading players", func(t *testing.T) {
		cr := NewCombatRequirements()
		creature := uuid.New()
		player1 := uuid.New()
		player2 := uuid.New()

		// Goad by first player
		cr.GoadCreature(creature, player1)
		assert.Equal(t, player1, cr.GetGoadingPlayer(creature))

		// Goad by second player (overwrites)
		cr.GoadCreature(creature, player2)
		assert.Equal(t, player2, cr.GetGoadingPlayer(creature))
	})

	t.Run("Evasion with restrictions", func(t *testing.T) {
		cr := NewCombatRequirements()
		attacker := uuid.New()

		// Unblockable
		cr.AddEvasion(attacker, CantBeBlockedAtAll)

		// Also can't be blocked by more than 1
		cr.AddEvasion(attacker, CantBeBlockedByMoreThan)

		// Can't be blocked at all wins
		assert.False(t, cr.CanBeBlocked(attacker))
	})
}

// BenchmarkCombatRequirements benchmarks performance with many creatures
func BenchmarkCombatRequirements(b *testing.B) {
	cr := NewCombatRequirements()

	// Create 100 creatures with various requirements
	creatures := make([]uuid.UUID, 100)
	for i := 0; i < 100; i++ {
		creatures[i] = uuid.New()

		if i%3 == 0 {
			cr.AddAttackRequirement(creatures[i], MustAttackIfAble)
		}
		if i%5 == 0 {
			cr.AddBlockRestriction(creatures[i], CantBlockFlying)
		}
		if i%7 == 0 {
			cr.AddEvasion(creatures[i], CantBeBlockedExceptBy)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Check all creatures
		for _, creature := range creatures {
			cr.CanAttack(creature)
			cr.CanBlock(creature)
			cr.MustAttackCheck(creature)
		}
	}
}
