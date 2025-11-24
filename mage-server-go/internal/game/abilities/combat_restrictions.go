package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Combat Restrictions and Requirements (Rules 508-509)
// This file implements effects that restrict or require combat actions:
// - Attack restrictions (can't attack, must attack, attacks if able)
// - Block restrictions (can't block, must block, blocks if able)
// - Evasion abilities (can't be blocked except by...)
// - Combat requirements (goad, provoke, etc.)

// ===== Attack Restrictions =====

// CantAttackEffect prevents a creature from attacking
// Rule 508.1d: Restrictions checked before declaring attackers
type CantAttackEffect struct {
	description    string
	source         uuid.UUID
	targetCreature uuid.UUID
	restriction    AttackRestriction
	duration       Duration
}

// AttackRestriction specifies what the creature can't attack
type AttackRestriction int

const (
	CantAttackAny                    AttackRestriction = iota // Can't attack at all
	CantAttackPlayer                                          // Can't attack players
	CantAttackPlaneswalker                                    // Can't attack planeswalkers
	CantAttackDefender                                        // Can't attack specific defender
	CantAttackAlone                                           // Can't attack alone
	CantAttackIfDefenderControlsType                          // Can't attack if defender controls specific type
)

// NewCantAttackEffect creates an attack restriction
func NewCantAttackEffect(source, creature uuid.UUID, restriction AttackRestriction, duration Duration) *CantAttackEffect {
	return &CantAttackEffect{
		description:    "Can't attack",
		source:         source,
		targetCreature: creature,
		restriction:    restriction,
		duration:       duration,
	}
}

// Apply implements the Effect interface
func (e *CantAttackEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Restriction effects are checked during attack declaration, not applied here
	return nil
}

// GetDescription returns effect description
func (e *CantAttackEffect) GetDescription() string {
	switch e.restriction {
	case CantAttackAny:
		return "Can't attack"
	case CantAttackPlayer:
		return "Can't attack players"
	case CantAttackPlaneswalker:
		return "Can't attack planeswalkers"
	case CantAttackAlone:
		return "Can't attack alone"
	default:
		return "Attack restricted"
	}
}

// ===== Attack Requirements =====

// MustAttackEffect forces a creature to attack
// Rule 508.1d: Requirements checked during attack declaration
type MustAttackEffect struct {
	description    string
	source         uuid.UUID
	targetCreature uuid.UUID
	requirement    AttackRequirement
	duration       Duration
}

// AttackRequirement specifies attack requirements
type AttackRequirement int

const (
	MustAttackIfAble       AttackRequirement = iota // Attacks each combat if able
	MustAttackPlayer                                // Must attack specific player if able
	MustAttackPlaneswalker                          // Must attack specific planeswalker if able
	MustAttackDefender                              // Must attack specific defender if able
	MustAttackEachTurn                              // Must attack each turn (goad)
)

// NewMustAttackEffect creates an attack requirement
func NewMustAttackEffect(source, creature uuid.UUID, requirement AttackRequirement, duration Duration) *MustAttackEffect {
	return &MustAttackEffect{
		description:    "Must attack",
		source:         source,
		targetCreature: creature,
		requirement:    requirement,
		duration:       duration,
	}
}

// Apply implements the Effect interface
func (e *MustAttackEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Requirements are checked during attack declaration
	return nil
}

// GetDescription returns effect description
func (e *MustAttackEffect) GetDescription() string {
	switch e.requirement {
	case MustAttackIfAble:
		return "Attacks each combat if able"
	case MustAttackEachTurn:
		return "Must attack each turn if able and attacks a player other than you"
	default:
		return "Must attack"
	}
}

// ===== Block Restrictions =====

// CantBlockEffect prevents blocking
// Rule 509.1b: Restrictions checked before declaring blockers
type CantBlockEffect struct {
	description    string
	source         uuid.UUID
	targetCreature uuid.UUID
	restriction    BlockRestriction
	duration       Duration
}

// BlockRestriction specifies what can't be blocked
type BlockRestriction int

const (
	CantBlockAny                 BlockRestriction = iota // Can't block at all (like Defender)
	CantBlockFlying                                      // Can't block creatures with flying
	CantBlockCreatureWithPower                           // Can't block creatures with power >= X
	CantBlockMoreThanOneCreature                         // Can block only one creature
	CantBlockAlone                                       // Can't block alone
)

// NewCantBlockEffect creates a block restriction
func NewCantBlockEffect(source, creature uuid.UUID, restriction BlockRestriction, duration Duration) *CantBlockEffect {
	return &CantBlockEffect{
		description:    "Can't block",
		source:         source,
		targetCreature: creature,
		restriction:    restriction,
		duration:       duration,
	}
}

// Apply implements the Effect interface
func (e *CantBlockEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Restrictions are checked during blocker declaration
	return nil
}

// GetDescription returns effect description
func (e *CantBlockEffect) GetDescription() string {
	switch e.restriction {
	case CantBlockAny:
		return "Can't block"
	case CantBlockFlying:
		return "Can't block creatures with flying"
	case CantBlockMoreThanOneCreature:
		return "Can block only one creature"
	case CantBlockAlone:
		return "Can't block alone"
	default:
		return "Block restricted"
	}
}

// ===== Block Requirements =====

// MustBlockEffect forces blocking
// Rule 509.1b: Requirements checked during blocker declaration
type MustBlockEffect struct {
	description    string
	source         uuid.UUID
	targetCreature uuid.UUID
	requirement    BlockRequirement
	targetAttacker uuid.UUID // Specific attacker to block
	duration       Duration
}

// BlockRequirement specifies block requirements
type BlockRequirement int

const (
	MustBlockIfAble           BlockRequirement = iota // Blocks if able
	MustBlockAttacker                                 // Must block specific attacker
	MustBlockWithAllCreatures                         // All creatures must block this
	MustBlockAlone                                    // Must block alone (provoke)
)

// NewMustBlockEffect creates a block requirement
func NewMustBlockEffect(source, creature, attacker uuid.UUID, requirement BlockRequirement, duration Duration) *MustBlockEffect {
	return &MustBlockEffect{
		description:    "Must block",
		source:         source,
		targetCreature: creature,
		requirement:    requirement,
		targetAttacker: attacker,
		duration:       duration,
	}
}

// Apply implements the Effect interface
func (e *MustBlockEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Requirements are checked during blocker declaration
	return nil
}

// GetDescription returns effect description
func (e *MustBlockEffect) GetDescription() string {
	switch e.requirement {
	case MustBlockIfAble:
		return "Blocks if able"
	case MustBlockAttacker:
		return fmt.Sprintf("Must block attacker %s if able", e.targetAttacker)
	case MustBlockWithAllCreatures:
		return "Must be blocked by all creatures that are able to block it"
	case MustBlockAlone:
		return "Must be blocked alone if any creature blocks it"
	default:
		return "Must block"
	}
}

// ===== Evasion Abilities =====

// CantBeBlockedEffect makes a creature unblockable or conditionally blockable
// Rule 509.1b: Evasion abilities
type CantBeBlockedEffect struct {
	description string
	source      uuid.UUID
	attacker    uuid.UUID
	condition   EvasionCondition
	duration    Duration
}

// EvasionCondition specifies blocking conditions
type EvasionCondition int

const (
	CantBeBlockedAtAll      EvasionCondition = iota // Unblockable
	CantBeBlockedExceptBy                           // Can only be blocked by X
	CantBeBlockedByMoreThan                         // Can't be blocked by more than X creatures
	CantBeBlockedByLessThan                         // Can't be blocked by fewer than X creatures
	CantBeBlockedByColor                            // Can't be blocked by color
	CantBeBlockedByType                             // Can't be blocked by creature type
)

// NewCantBeBlockedEffect creates an evasion effect
func NewCantBeBlockedEffect(source, attacker uuid.UUID, condition EvasionCondition, duration Duration) *CantBeBlockedEffect {
	return &CantBeBlockedEffect{
		description: "Can't be blocked",
		source:      source,
		attacker:    attacker,
		condition:   condition,
		duration:    duration,
	}
}

// Apply implements the Effect interface
func (e *CantBeBlockedEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Evasion effects are checked during blocker declaration
	return nil
}

// GetDescription returns effect description
func (e *CantBeBlockedEffect) GetDescription() string {
	switch e.condition {
	case CantBeBlockedAtAll:
		return "Can't be blocked"
	case CantBeBlockedExceptBy:
		return "Can't be blocked except by specific creatures"
	case CantBeBlockedByMoreThan:
		return "Can't be blocked by more than X creatures"
	case CantBeBlockedByColor:
		return "Can't be blocked by creatures of specific color"
	default:
		return "Evasion ability"
	}
}

// ===== Special Combat Abilities =====

// ProvokeAbility represents the Provoke keyword
// Rule 702.39: Provoke
type ProvokeAbility struct {
	baseAbility
	targetCreature uuid.UUID
}

// NewProvokeAbility creates a Provoke ability
func NewProvokeAbility(source uuid.UUID) *ProvokeAbility {
	return &ProvokeAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *ProvokeAbility) GetType() AbilityType {
	return AbilityTypeTriggered
}

func (a *ProvokeAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return false // Triggered, not activated
}

func (a *ProvokeAbility) Resolve(ctx context.Context, game GameContext) error {
	// When this attacks, target creature blocks it this combat if able
	// The targeted creature must block alone
	return nil
}

func (a *ProvokeAbility) String() string {
	return "Provoke"
}

// GoadAbility represents the Goad keyword action effect
// Rule 701.38: Goad
type GoadAbility struct {
	baseAbility
	goadedCreatures map[uuid.UUID]bool
	goadingPlayer   uuid.UUID
}

// NewGoadAbility creates a Goad ability
func NewGoadAbility(source, goadingPlayer uuid.UUID) *GoadAbility {
	return &GoadAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		goadedCreatures: make(map[uuid.UUID]bool),
		goadingPlayer:   goadingPlayer,
	}
}

func (a *GoadAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *GoadAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return false
}

func (a *GoadAbility) Resolve(ctx context.Context, game GameContext) error {
	// Goaded creatures must attack each combat if able
	// Must attack a player other than the goading player if able
	return nil
}

func (a *GoadAbility) String() string {
	return "Goad"
}

// ===== Combat Requirement Tracking =====

// CombatRequirements tracks all combat restrictions and requirements
type CombatRequirements struct {
	// Attack restrictions
	cantAttack       map[uuid.UUID][]AttackRestriction
	mustAttack       map[uuid.UUID][]AttackRequirement
	mustAttackPlayer map[uuid.UUID]uuid.UUID // creature -> player it must attack

	// Block restrictions
	cantBlock         map[uuid.UUID][]BlockRestriction
	mustBlock         map[uuid.UUID][]BlockRequirement
	mustBlockAttacker map[uuid.UUID]uuid.UUID // creature -> attacker it must block

	// Evasion
	cantBeBlocked         map[uuid.UUID][]EvasionCondition
	cantBeBlockedExceptBy map[uuid.UUID][]string // attacker -> types/colors that can block

	// Special requirements
	mustBeBlockedByAll map[uuid.UUID]bool // attacker -> must be blocked by all able
	mustBeBlockedAlone map[uuid.UUID]bool // attacker -> must be blocked alone
	maxBlockers        map[uuid.UUID]int  // attacker -> max number of blockers
	minBlockers        map[uuid.UUID]int  // attacker -> min number of blockers

	// Goad tracking
	goadedCreatures map[uuid.UUID]uuid.UUID // creature -> goading player
}

// NewCombatRequirements creates a new combat requirements tracker
func NewCombatRequirements() *CombatRequirements {
	return &CombatRequirements{
		cantAttack:            make(map[uuid.UUID][]AttackRestriction),
		mustAttack:            make(map[uuid.UUID][]AttackRequirement),
		mustAttackPlayer:      make(map[uuid.UUID]uuid.UUID),
		cantBlock:             make(map[uuid.UUID][]BlockRestriction),
		mustBlock:             make(map[uuid.UUID][]BlockRequirement),
		mustBlockAttacker:     make(map[uuid.UUID]uuid.UUID),
		cantBeBlocked:         make(map[uuid.UUID][]EvasionCondition),
		cantBeBlockedExceptBy: make(map[uuid.UUID][]string),
		mustBeBlockedByAll:    make(map[uuid.UUID]bool),
		mustBeBlockedAlone:    make(map[uuid.UUID]bool),
		maxBlockers:           make(map[uuid.UUID]int),
		minBlockers:           make(map[uuid.UUID]int),
		goadedCreatures:       make(map[uuid.UUID]uuid.UUID),
	}
}

// AddAttackRestriction adds an attack restriction
func (cr *CombatRequirements) AddAttackRestriction(creature uuid.UUID, restriction AttackRestriction) {
	cr.cantAttack[creature] = append(cr.cantAttack[creature], restriction)
}

// AddAttackRequirement adds an attack requirement
func (cr *CombatRequirements) AddAttackRequirement(creature uuid.UUID, requirement AttackRequirement) {
	cr.mustAttack[creature] = append(cr.mustAttack[creature], requirement)
}

// AddBlockRestriction adds a block restriction
func (cr *CombatRequirements) AddBlockRestriction(creature uuid.UUID, restriction BlockRestriction) {
	cr.cantBlock[creature] = append(cr.cantBlock[creature], restriction)
}

// AddBlockRequirement adds a block requirement
func (cr *CombatRequirements) AddBlockRequirement(creature uuid.UUID, requirement BlockRequirement, attacker uuid.UUID) {
	cr.mustBlock[creature] = append(cr.mustBlock[creature], requirement)
	if attacker != uuid.Nil {
		cr.mustBlockAttacker[creature] = attacker
	}
}

// AddEvasion adds an evasion condition
func (cr *CombatRequirements) AddEvasion(attacker uuid.UUID, condition EvasionCondition) {
	cr.cantBeBlocked[attacker] = append(cr.cantBeBlocked[attacker], condition)
}

// GoadCreature goads a creature
func (cr *CombatRequirements) GoadCreature(creature, goadingPlayer uuid.UUID) {
	cr.goadedCreatures[creature] = goadingPlayer
}

// CanAttack checks if a creature can attack
func (cr *CombatRequirements) CanAttack(creature uuid.UUID) bool {
	restrictions, exists := cr.cantAttack[creature]
	if !exists {
		return true
	}

	for _, restriction := range restrictions {
		if restriction == CantAttackAny {
			return false
		}
	}

	return true
}

// MustAttackCheck checks if a creature must attack
func (cr *CombatRequirements) MustAttackCheck(creature uuid.UUID) bool {
	requirements, exists := cr.mustAttack[creature]
	if !exists {
		return false
	}

	for _, requirement := range requirements {
		if requirement == MustAttackIfAble || requirement == MustAttackEachTurn {
			return true
		}
	}

	return false
}

// CanBlock checks if a creature can block
func (cr *CombatRequirements) CanBlock(creature uuid.UUID) bool {
	restrictions, exists := cr.cantBlock[creature]
	if !exists {
		return true
	}

	for _, restriction := range restrictions {
		if restriction == CantBlockAny {
			return false
		}
	}

	return true
}

// CanBeBlocked checks if an attacker can be blocked
func (cr *CombatRequirements) CanBeBlocked(attacker uuid.UUID) bool {
	conditions, exists := cr.cantBeBlocked[attacker]
	if !exists {
		return true
	}

	for _, condition := range conditions {
		if condition == CantBeBlockedAtAll {
			return false
		}
	}

	return true
}

// IsGoaded checks if a creature is goaded
func (cr *CombatRequirements) IsGoaded(creature uuid.UUID) bool {
	_, goaded := cr.goadedCreatures[creature]
	return goaded
}

// GetGoadingPlayer returns the player who goaded this creature
func (cr *CombatRequirements) GetGoadingPlayer(creature uuid.UUID) uuid.UUID {
	return cr.goadedCreatures[creature]
}

// ValidateAttackDeclaration validates that attack declaration follows all requirements
func (cr *CombatRequirements) ValidateAttackDeclaration(attackers []uuid.UUID, activePlayer uuid.UUID) []string {
	violations := make([]string, 0)

	// Check creatures that must attack
	for creature, requirements := range cr.mustAttack {
		foundAttacking := false
		for _, attacker := range attackers {
			if attacker == creature {
				foundAttacking = true
				break
			}
		}

		if !foundAttacking {
			for _, req := range requirements {
				if req == MustAttackIfAble || req == MustAttackEachTurn {
					violations = append(violations,
						fmt.Sprintf("Creature %s must attack if able", creature))
					break
				}
			}
		}
	}

	// Check goaded creatures
	for creature, goadingPlayer := range cr.goadedCreatures {
		if goadingPlayer != activePlayer {
			foundAttacking := false
			for _, attacker := range attackers {
				if attacker == creature {
					foundAttacking = true
					break
				}
			}

			if !foundAttacking {
				violations = append(violations,
					fmt.Sprintf("Goaded creature %s must attack a player other than %s", creature, goadingPlayer))
			}
		}
	}

	return violations
}

// ValidateBlockDeclaration validates that block declaration follows all requirements
func (cr *CombatRequirements) ValidateBlockDeclaration(blocks map[uuid.UUID][]uuid.UUID) []string {
	violations := make([]string, 0)

	// Check creatures that must block
	for creature, requirements := range cr.mustBlock {
		foundBlocking := false
		for _, blockers := range blocks {
			for _, blocker := range blockers {
				if blocker == creature {
					foundBlocking = true
					break
				}
			}
			if foundBlocking {
				break
			}
		}

		if !foundBlocking {
			for _, req := range requirements {
				if req == MustBlockIfAble {
					violations = append(violations,
						fmt.Sprintf("Creature %s must block if able", creature))
					break
				}
			}
		}
	}

	// Check must-be-blocked-by-all requirements
	for attacker := range cr.mustBeBlockedByAll {
		blockerCount := 0
		if blockers, exists := blocks[attacker]; exists {
			blockerCount = len(blockers)
		}

		// TODO: Check against number of able blockers
		_ = blockerCount
	}

	return violations
}

// Clear clears all combat requirements (called at end of combat)
func (cr *CombatRequirements) Clear() {
	cr.cantAttack = make(map[uuid.UUID][]AttackRestriction)
	cr.mustAttack = make(map[uuid.UUID][]AttackRequirement)
	cr.mustAttackPlayer = make(map[uuid.UUID]uuid.UUID)
	cr.cantBlock = make(map[uuid.UUID][]BlockRestriction)
	cr.mustBlock = make(map[uuid.UUID][]BlockRequirement)
	cr.mustBlockAttacker = make(map[uuid.UUID]uuid.UUID)
	cr.cantBeBlocked = make(map[uuid.UUID][]EvasionCondition)
	cr.cantBeBlockedExceptBy = make(map[uuid.UUID][]string)
	cr.mustBeBlockedByAll = make(map[uuid.UUID]bool)
	cr.mustBeBlockedAlone = make(map[uuid.UUID]bool)
	cr.maxBlockers = make(map[uuid.UUID]int)
	cr.minBlockers = make(map[uuid.UUID]int)
	// Note: Goad persists until next turn
}

// ===== Integration with Existing Systems =====

// These combat restrictions integrate with:
// - Combat system: Checked during attack/block declaration (Rules 508-509)
// - Static abilities: Restrictions are continuous effects
// - Triggered abilities: Provoke triggers on attack
// - State-based actions: Some restrictions cause illegal blocks
// - Turn structure: Goad persists across turns

// Important rules:
// - Rule 508.1d: Attack restrictions and requirements
// - Rule 509.1b: Block restrictions and requirements
// - Rule 509.1c: Evasion abilities ("can't be blocked except by")
// - Rule 701.38: Goad
// - Rule 702.39: Provoke
// - Rule 702.67: Menace (can't be blocked except by 2+)
// - Rule 509.2: Multiple block restrictions

// Example usage for "can't attack" (Pacifism):
// 1. Enchanted creature can't attack or block
// 2. Create CantAttackEffect with restriction=CantAttackAny
// 3. Create CantBlockEffect with restriction=CantBlockAny
// 4. During declare attackers, check canAttack()
// 5. During declare blockers, check canBlock()

// Example usage for Provoke:
// 1. Creature with Provoke attacks
// 2. Trigger: Choose target creature
// 3. Create MustBlockEffect for target with requirement=MustBlockAttacker
// 4. Add mustBeBlockedAlone flag
// 5. During declare blockers, enforce requirement

// Example usage for Goad:
// 1. Effect: "Goad target creature"
// 2. Create GoadAbility, track goading player
// 3. Create MustAttackEffect with requirement=MustAttackEachTurn
// 4. Add restriction: can't attack goading player
// 5. Effect persists until goaded creature's controller's next turn
