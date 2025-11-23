package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ========================================
// AttachEffect - Attaches an Aura or Equipment to a permanent
// ========================================

// AttachEffect handles the attachment of Auras and Equipment
// Java: mage.abilities.effects.common.AttachEffect
type AttachEffect struct {
	outcome Outcome // What type of effect this is (Benefit, Protect, etc.)
}

// Outcome represents the expected outcome of an effect
type Outcome string

const (
	OutcomeBenefit       Outcome = "Benefit"
	OutcomeBoostCreature Outcome = "BoostCreature"
	OutcomeAddAbility    Outcome = "AddAbility"
	OutcomeProtect       Outcome = "Protect"
	OutcomeDetriment     Outcome = "Detriment"
	OutcomeNeutral       Outcome = "Neutral"
)

// NewAttachEffect creates a new attach effect
func NewAttachEffect(outcome Outcome) *AttachEffect {
	return &AttachEffect{outcome: outcome}
}

// Apply attaches the source to the target
func (e *AttachEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	if len(targets) == 0 {
		return fmt.Errorf("no targets for attach effect")
	}

	// TODO: Implement actual attachment logic
	// This requires:
	// 1. Mark the source card as attached to the target
	// 2. Update target's attached list
	// 3. Handle state-based actions if target is invalid

	// For now, this is a placeholder
	_ = source
	_ = targets

	return nil
}

// GetDescription returns a description of the effect
func (e *AttachEffect) GetDescription() string {
	return "attach"
}

// ========================================
// EnchantAbility - Represents "Enchant X" ability
// ========================================

// EnchantAbility represents the "Enchant creature", "Enchant artifact", etc. ability
// Java: mage.abilities.keyword.EnchantAbility
type EnchantAbility struct {
	baseAbility
	target *TargetRequirement
}

// NewEnchantAbility creates a new enchant ability
func NewEnchantAbility(sourceID uuid.UUID, target *TargetRequirement) *EnchantAbility {
	text := fmt.Sprintf("Enchant %s", target.Description)
	return &EnchantAbility{
		baseAbility: newBaseAbility(sourceID, text),
		target:      target,
	}
}

// GetType returns the ability type
func (a *EnchantAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

// CanActivate always returns true (enchant is a static ability)
func (a *EnchantAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

// Resolve does nothing (enchant is a static ability)
func (a *EnchantAbility) Resolve(ctx context.Context, game GameContext) error {
	return nil
}

// GetTarget returns the target requirement
func (a *EnchantAbility) GetTarget() *TargetRequirement {
	return a.target
}

// ========================================
// EquipAbility - Represents "Equip {cost}" ability
// ========================================

// EquipAbility represents the "Equip {cost}" ability on Equipment
// Java: mage.abilities.keyword.EquipAbility
type EquipAbility struct {
	baseAbility
	cost         *ManaCost
	sorcerySpeed bool // true if can only activate at sorcery speed
}

// NewEquipAbility creates a new equip ability
// cost is the mana cost to equip (e.g., "{2}" for "Equip {2}")
// sorcerySpeed should be true for normal equip (false for instant-speed equip like "Equip Knight")
func NewEquipAbility(sourceID uuid.UUID, cost string, sorcerySpeed bool) (*EquipAbility, error) {
	manaCost, err := ParseManaCost(cost)
	if err != nil {
		return nil, fmt.Errorf("failed to parse equip cost: %w", err)
	}

	text := fmt.Sprintf("Equip %s", cost)
	return &EquipAbility{
		baseAbility:  newBaseAbility(sourceID, text),
		cost:         manaCost,
		sorcerySpeed: sorcerySpeed,
	}, nil
}

// GetType returns the ability type
func (a *EquipAbility) GetType() AbilityType {
	return AbilityTypeActivated
}

// CanActivate checks if this ability can be activated
func (a *EquipAbility) CanActivate(ctx context.Context, game GameContext) bool {
	// TODO: Check if:
	// 1. You control the equipment
	// 2. You control a creature to equip it to
	// 3. It's your main phase (if sorcerySpeed)
	// 4. You can pay the cost
	return true
}

// Resolve equips the equipment to a target creature
func (a *EquipAbility) Resolve(ctx context.Context, game GameContext) error {
	// TODO: Implement equip logic
	// This requires selecting a target creature and attaching to it
	return fmt.Errorf("equip ability not yet implemented")
}

// GetCost returns the equip cost
func (a *EquipAbility) GetCost() *ManaCost {
	return a.cost
}

// ========================================
// BoostEnchantedEffect - Boosts enchanted creature
// ========================================

// BoostEnchantedEffect modifies power/toughness of enchanted creature
// Java: mage.abilities.effects.common.continuous.BoostEnchantedEffect
type BoostEnchantedEffect struct {
	power     int
	toughness int
}

// NewBoostEnchantedEffect creates a new boost enchanted effect
func NewBoostEnchantedEffect(power, toughness int) *BoostEnchantedEffect {
	return &BoostEnchantedEffect{
		power:     power,
		toughness: toughness,
	}
}

// Apply applies the boost effect
func (e *BoostEnchantedEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Implement boost enchanted logic
	// This requires finding the enchanted creature and applying a continuous effect
	return fmt.Errorf("boost enchanted effect not yet implemented")
}

// GetDescription returns a description
func (e *BoostEnchantedEffect) GetDescription() string {
	powerStr := formatBoost(e.power)
	toughnessStr := formatBoost(e.toughness)
	return fmt.Sprintf("enchanted creature gets %s/%s", powerStr, toughnessStr)
}

// ========================================
// BoostEquippedEffect - Boosts equipped creature
// ========================================

// BoostEquippedEffect modifies power/toughness of equipped creature
// Java: mage.abilities.effects.common.continuous.BoostEquippedEffect
type BoostEquippedEffect struct {
	power     int
	toughness int
}

// NewBoostEquippedEffect creates a new boost equipped effect
func NewBoostEquippedEffect(power, toughness int) *BoostEquippedEffect {
	return &BoostEquippedEffect{
		power:     power,
		toughness: toughness,
	}
}

// Apply applies the boost effect
func (e *BoostEquippedEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Implement boost equipped logic
	// This requires finding the equipped creature and applying a continuous effect
	return fmt.Errorf("boost equipped effect not yet implemented")
}

// GetDescription returns a description
func (e *BoostEquippedEffect) GetDescription() string {
	powerStr := formatBoost(e.power)
	toughnessStr := formatBoost(e.toughness)
	return fmt.Sprintf("equipped creature gets %s/%s", powerStr, toughnessStr)
}

// ========================================
// GainAbilityAttachedEffect - Grants ability to attached creature
// ========================================

// AttachmentType specifies what type of attachment this is
type AttachmentType string

const (
	AttachmentTypeAura      AttachmentType = "AURA"
	AttachmentTypeEquipment AttachmentType = "EQUIPMENT"
)

// GainAbilityAttachedEffect grants an ability to the attached permanent
// Java: mage.abilities.effects.common.continuous.GainAbilityAttachedEffect
type GainAbilityAttachedEffect struct {
	abilityID      string         // The ability being granted (e.g., "TrampleAbility")
	attachmentType AttachmentType // Whether this is from an Aura or Equipment
}

// NewGainAbilityAttachedEffect creates a new gain ability attached effect
func NewGainAbilityAttachedEffect(abilityID string, attachmentType AttachmentType) *GainAbilityAttachedEffect {
	return &GainAbilityAttachedEffect{
		abilityID:      abilityID,
		attachmentType: attachmentType,
	}
}

// Apply applies the effect
func (e *GainAbilityAttachedEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Implement gain ability attached logic
	// This requires finding the attached permanent and granting the ability
	return fmt.Errorf("gain ability attached effect not yet implemented")
}

// GetDescription returns a description
func (e *GainAbilityAttachedEffect) GetDescription() string {
	// Convert ability ID to readable text
	abilityName := e.abilityID
	if len(abilityName) > 7 && abilityName[len(abilityName)-7:] == "Ability" {
		abilityName = abilityName[:len(abilityName)-7]
	}

	target := "attached"
	if e.attachmentType == AttachmentTypeAura {
		target = "enchanted"
	} else if e.attachmentType == AttachmentTypeEquipment {
		target = "equipped"
	}

	return fmt.Sprintf("%s creature has %s", target, abilityName)
}
