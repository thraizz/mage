package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Rule 702: Keyword Abilities
// This file implements the mechanical effects of critical keyword abilities.
// Basic keyword abilities (Flying, Trample, etc.) are tracked via KeywordAbility
// but their effects are implemented in combat/targeting systems.

// ===== Flash (Rule 702.8) =====

// FlashAbility represents the Flash keyword
// Rule 702.8a: Flash is a static ability that allows casting at instant speed
type FlashAbility struct {
	baseAbility
}

// NewFlashAbility creates a Flash keyword ability
func NewFlashAbility(source uuid.UUID) *FlashAbility {
	return &FlashAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *FlashAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *FlashAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true // Always active
}

func (a *FlashAbility) Resolve(ctx context.Context, game GameContext) error {
	// Flash doesn't resolve - it modifies casting timing
	// Checked during Rule 601.2 (casting process)
	return nil
}

func (a *FlashAbility) String() string {
	return "Flash"
}

// HasFlash checks if a card has Flash
func HasFlash(card interface{}) bool {
	// This would check the card's abilities for Flash
	// Implementation depends on card structure
	return false // Placeholder
}

// ===== Haste (Rule 702.10) =====

// HasteAbility represents the Haste keyword
// Rule 702.10a: Haste removes summoning sickness
type HasteAbility struct {
	baseAbility
}

// NewHasteAbility creates a Haste keyword ability
func NewHasteAbility(source uuid.UUID) *HasteAbility {
	return &HasteAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *HasteAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *HasteAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true // Always active
}

func (a *HasteAbility) Resolve(ctx context.Context, game GameContext) error {
	// Haste doesn't resolve - it's a static ability
	// Checked during attack/activate ability declaration
	return nil
}

func (a *HasteAbility) String() string {
	return "Haste"
}

// HasHaste checks if a permanent has Haste
func HasHaste(permanentID uuid.UUID, game GameContext) bool {
	// This would check the permanent's abilities for Haste
	// Used in attack declaration and tap ability activation
	return false // Placeholder
}

// CanAttackOrActivate checks if a permanent can attack/tap despite summoning sickness
func CanAttackOrActivate(permanentID uuid.UUID, game GameContext, enteredThisTurn bool) bool {
	// If it didn't enter this turn, it can attack/tap
	if !enteredThisTurn {
		return true
	}

	// If it has Haste, it can attack/tap even if it entered this turn
	return HasHaste(permanentID, game)
}

// ===== Indestructible (Rule 702.12) =====

// IndestructibleAbility represents the Indestructible keyword
// Rule 702.12a: Indestructible permanents can't be destroyed
// Rule 702.12b: Damage can't cause them to be destroyed
type IndestructibleAbility struct {
	baseAbility
}

// NewIndestructibleAbility creates an Indestructible keyword ability
func NewIndestructibleAbility(source uuid.UUID) *IndestructibleAbility {
	return &IndestructibleAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *IndestructibleAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *IndestructibleAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true // Always active
}

func (a *IndestructibleAbility) Resolve(ctx context.Context, game GameContext) error {
	// Indestructible doesn't resolve - it's checked during:
	// - Destroy effects (Rule 701.7a)
	// - State-based actions for lethal damage (Rule 704.5g)
	// Implementation via IndestructiblePreventionEffect in effects/prevention.go
	return nil
}

func (a *IndestructibleAbility) String() string {
	return "Indestructible"
}

// HasIndestructible checks if a permanent has Indestructible
func HasIndestructible(permanentID uuid.UUID, game GameContext) bool {
	// This would check the permanent's abilities for Indestructible
	// Used in destroy effects and state-based actions
	return false // Placeholder
}

// ===== Ward (Rule 702.21) =====

// WardAbility represents the Ward keyword
// Rule 702.21a: Ward [cost] - Whenever this becomes the target of a spell or
// ability an opponent controls, counter it unless that player pays [cost]
type WardAbility struct {
	baseAbility
	cost     *ManaCost
	genericN int    // For "Ward {N}"
	costText string // For other costs like "Ward—Discard a card"
}

// NewWardAbility creates a Ward keyword ability with a mana cost
func NewWardAbility(source uuid.UUID, cost *ManaCost) *WardAbility {
	return &WardAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		cost: cost,
	}
}

// NewWardAbilityGeneric creates a Ward ability with a generic mana cost (Ward {N})
func NewWardAbilityGeneric(source uuid.UUID, amount int) *WardAbility {
	return &WardAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		genericN: amount,
		costText: fmt.Sprintf("{%d}", amount),
	}
}

// NewWardAbilityCustom creates a Ward ability with a custom cost
func NewWardAbilityCustom(source uuid.UUID, costText string) *WardAbility {
	return &WardAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		costText: costText,
	}
}

func (a *WardAbility) GetType() AbilityType {
	return AbilityTypeTriggered
}

func (a *WardAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true // Always active
}

func (a *WardAbility) Resolve(ctx context.Context, game GameContext) error {
	// Ward triggers when the permanent becomes a target
	// The opponent must pay the cost or the spell/ability is countered
	// Implementation:
	// 1. Trigger when source becomes target of opponent's spell/ability
	// 2. Put triggered ability on stack
	// 3. On resolution: opponent chooses to pay or not
	// 4. If not paid: counter the spell/ability that targeted this
	return nil
}

func (a *WardAbility) GetCost() *ManaCost {
	return a.cost
}

func (a *WardAbility) GetGenericAmount() int {
	return a.genericN
}

func (a *WardAbility) GetCostText() string {
	if a.costText != "" {
		return a.costText
	}
	if a.cost != nil {
		return a.cost.String()
	}
	return ""
}

func (a *WardAbility) String() string {
	if a.costText != "" {
		return fmt.Sprintf("Ward—%s", a.costText)
	}
	if a.genericN > 0 {
		return fmt.Sprintf("Ward {%d}", a.genericN)
	}
	if a.cost != nil {
		return fmt.Sprintf("Ward %s", a.cost.String())
	}
	return "Ward"
}

// WardTrigger represents the triggered ability created by Ward
type WardTrigger struct {
	*TriggeredAbility
	wardSource              uuid.UUID
	targetingSpellOrAbility uuid.UUID
	opponent                uuid.UUID
	cost                    *ManaCost
	costText                string
}

// NewWardTrigger creates a Ward trigger
func NewWardTrigger(
	wardSource uuid.UUID,
	targetingSpellOrAbility uuid.UUID,
	opponent uuid.UUID,
	cost *ManaCost,
	costText string,
) *WardTrigger {
	return &WardTrigger{
		wardSource:              wardSource,
		targetingSpellOrAbility: targetingSpellOrAbility,
		opponent:                opponent,
		cost:                    cost,
		costText:                costText,
	}
}

// ResolveWardTrigger resolves a Ward trigger
func (wt *WardTrigger) ResolveWardTrigger(ctx context.Context, game GameContext) error {
	// 1. Ask opponent if they want to pay the cost
	// 2. If they choose not to pay or can't pay, counter the spell/ability
	// 3. If they pay, nothing happens (spell/ability continues normally)

	// Placeholder implementation - would need UI interaction
	// For now, assume opponent doesn't pay
	paid := false

	if !paid {
		// Counter the spell or ability
		// This would call the CounterSpellEffect or similar
		return fmt.Errorf("ward trigger: spell/ability countered by ward")
	}

	return nil
}

// ===== Helper Functions =====

// CheckKeywordAbility checks if a permanent has a specific keyword ability
func CheckKeywordAbility(permanentID uuid.UUID, keyword KeywordType, game GameContext) bool {
	// This would iterate through the permanent's abilities
	// and check for the keyword
	// Implementation depends on GameContext structure
	return false // Placeholder
}

// GetKeywordAbilities returns all keyword abilities of a permanent
func GetKeywordAbilities(permanentID uuid.UUID, game GameContext) []KeywordType {
	// This would collect all keyword abilities from a permanent
	// Used for combat calculations, targeting restrictions, etc.
	return []KeywordType{} // Placeholder
}

// ===== Integration with Existing Systems =====

// These keyword abilities integrate with:
// - Flash: Checked during spell casting (Rule 601.2, timing restrictions)
// - Haste: Checked during attack declaration and tap ability activation
// - Indestructible: Creates IndestructiblePreventionEffect (already implemented in effects/prevention.go)
// - Ward: Triggered ability that goes on stack when permanent becomes a target

// Example usage for Indestructible:
// When a permanent gains Indestructible, create an IndestructiblePreventionEffect
// and register it with the ReplacementManager

// Example usage for Ward:
// When a permanent becomes a target of an opponent's spell/ability,
// check if it has Ward. If yes, create a WardTrigger and push to stack.
