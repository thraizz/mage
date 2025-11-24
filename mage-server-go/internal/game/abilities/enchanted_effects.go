package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ========================================
// TapEnchantedEffect
// ========================================

// TapEnchantedEffect taps the enchanted permanent
// Java: mage.abilities.effects.common.TapEnchantedEffect
// MTG Rules: 303.4m (enchanted [object or player])
type TapEnchantedEffect struct {
	description string // "permanent" or "creature"
}

// NewTapEnchantedEffect creates a new tap enchanted effect
func NewTapEnchantedEffect() *TapEnchantedEffect {
	return &TapEnchantedEffect{
		description: "permanent",
	}
}

// Apply taps the enchanted permanent
func (e *TapEnchantedEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Implement actual tap logic
	// This requires:
	// 1. Find what the source is attached to
	// 2. Tap that permanent

	// For now, this is a placeholder
	_ = source
	_ = targets

	return nil
}

// GetDescription returns a description of the effect
func (e *TapEnchantedEffect) GetDescription() string {
	return fmt.Sprintf("tap enchanted %s", e.description)
}

// ========================================
// UntapEnchantedEffect
// ========================================

// UntapEnchantedEffect untaps the enchanted permanent
// Java: mage.abilities.effects.common.UntapEnchantedEffect
type UntapEnchantedEffect struct {
	description string
}

// NewUntapEnchantedEffect creates a new untap enchanted effect
func NewUntapEnchantedEffect() *UntapEnchantedEffect {
	return &UntapEnchantedEffect{
		description: "permanent",
	}
}

// Apply untaps the enchanted permanent
func (e *UntapEnchantedEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Implement actual untap logic
	return nil
}

// GetDescription returns a description
func (e *UntapEnchantedEffect) GetDescription() string {
	return fmt.Sprintf("untap enchanted %s", e.description)
}

// ========================================
// DontUntapInControllersUntapStepEnchantedEffect
// ========================================

// DontUntapInControllersUntapStepEnchantedEffect prevents the enchanted permanent from untapping
// Java: mage.abilities.effects.common.DontUntapInControllersUntapStepEnchantedEffect
// MTG Rules: 502.1 (Untap Step)
type DontUntapInControllersUntapStepEnchantedEffect struct {
	baseContinuousEffect
	description string
}

// NewDontUntapInControllersUntapStepEnchantedEffect creates a new don't untap effect
func NewDontUntapInControllersUntapStepEnchantedEffect() *DontUntapInControllersUntapStepEnchantedEffect {
	return &DontUntapInControllersUntapStepEnchantedEffect{
		baseContinuousEffect: baseContinuousEffect{
			layer:    LayerAbilityAddingRemoving,
			duration: DurationWhileOnBattlefield,
		},
		description: "permanent",
	}
}

// Apply applies the effect (prevents untapping)
func (e *DontUntapInControllersUntapStepEnchantedEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Implement actual don't untap logic
	// This requires:
	// 1. Find what the source is attached to
	// 2. Mark that permanent as "doesn't untap during untap step"
	// 3. The untap step logic needs to check this flag

	// For now, this is a placeholder
	_ = source
	_ = targets

	return nil
}

// GetDescription returns a description
func (e *DontUntapInControllersUntapStepEnchantedEffect) GetDescription() string {
	return fmt.Sprintf("enchanted %s doesn't untap during its controller's untap step", e.description)
}

// ========================================
// UntapSourceEffect
// ========================================

// UntapSourceEffect untaps the source permanent
// Java: mage.abilities.effects.common.UntapSourceEffect
// MTG Rules: 701.22b (Untap)
type UntapSourceEffect struct{}

// NewUntapSourceEffect creates a new untap source effect
func NewUntapSourceEffect() *UntapSourceEffect {
	return &UntapSourceEffect{}
}

// Apply untaps the source permanent
func (e *UntapSourceEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Use the game context's UntapPermanent method
	return game.UntapPermanent(source)
}

// GetDescription returns a description
func (e *UntapSourceEffect) GetDescription() string {
	return "untap this permanent"
}

// ========================================
// TapSourceEffect
// ========================================

// TapSourceEffect taps the source permanent
// Java: mage.abilities.effects.common.TapSourceEffect
// MTG Rules: 701.21a (Tap)
type TapSourceEffect struct{}

// NewTapSourceEffect creates a new tap source effect
func NewTapSourceEffect() *TapSourceEffect {
	return &TapSourceEffect{}
}

// Apply taps the source permanent
func (e *TapSourceEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Use the game context's TapPermanent method
	return game.TapPermanent(source)
}

// GetDescription returns a description
func (e *TapSourceEffect) GetDescription() string {
	return "tap this permanent"
}

// ========================================
// DontUntapInControllersUntapStepSourceEffect
// ========================================

// DontUntapInControllersUntapStepSourceEffect prevents the source from untapping
// Java: mage.abilities.effects.common.DontUntapInControllersUntapStepSourceEffect
type DontUntapInControllersUntapStepSourceEffect struct {
	baseContinuousEffect
}

// NewDontUntapInControllersUntapStepSourceEffect creates a new don't untap source effect
func NewDontUntapInControllersUntapStepSourceEffect() *DontUntapInControllersUntapStepSourceEffect {
	return &DontUntapInControllersUntapStepSourceEffect{
		baseContinuousEffect: baseContinuousEffect{
			layer:    LayerAbilityAddingRemoving,
			duration: DurationWhileOnBattlefield,
		},
	}
}

// Apply applies the effect (prevents untapping)
func (e *DontUntapInControllersUntapStepSourceEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Implement actual don't untap logic
	// Mark the source as "doesn't untap during untap step"
	return nil
}

// GetDescription returns a description
func (e *DontUntapInControllersUntapStepSourceEffect) GetDescription() string {
	return "this permanent doesn't untap during your untap step"
}
