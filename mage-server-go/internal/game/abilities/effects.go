package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Effect represents what an ability does when it resolves
type Effect interface {
	// Apply applies this effect
	Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error

	// GetDescription returns a text description of this effect
	GetDescription() string
}

// ========================================
// Damage Effects
// ========================================

// DamageEffect deals damage to targets
type DamageEffect struct {
	Amount int
}

func NewDamageEffect(amount int) *DamageEffect {
	return &DamageEffect{Amount: amount}
}

func (e *DamageEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	for _, target := range targets {
		if err := game.DealDamage(source, target, e.Amount); err != nil {
			return fmt.Errorf("failed to deal damage: %w", err)
		}
	}
	return nil
}

func (e *DamageEffect) GetDescription() string {
	return fmt.Sprintf("deals %d damage", e.Amount)
}

// ========================================
// Card Draw Effects
// ========================================

// DrawCardsEffect has a player draw cards
type DrawCardsEffect struct {
	Amount int
}

func NewDrawCardsEffect(amount int) *DrawCardsEffect {
	return &DrawCardsEffect{Amount: amount}
}

func (e *DrawCardsEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// If no targets, controller draws
	if len(targets) == 0 {
		// TODO: Get controller from source
		return fmt.Errorf("draw effect needs target player")
	}

	for _, target := range targets {
		if err := game.DrawCards(target, e.Amount); err != nil {
			return fmt.Errorf("failed to draw cards: %w", err)
		}
	}
	return nil
}

func (e *DrawCardsEffect) GetDescription() string {
	if e.Amount == 1 {
		return "draw a card"
	}
	return fmt.Sprintf("draw %d cards", e.Amount)
}

// ========================================
// Destroy Effects
// ========================================

// DestroyEffect destroys target permanents
type DestroyEffect struct {
	NoRegeneration bool
}

func NewDestroyEffect() *DestroyEffect {
	return &DestroyEffect{NoRegeneration: false}
}

func NewDestroyEffectNoRegen() *DestroyEffect {
	return &DestroyEffect{NoRegeneration: true}
}

func (e *DestroyEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	for _, target := range targets {
		if err := game.DestroyPermanent(target); err != nil {
			return fmt.Errorf("failed to destroy permanent: %w", err)
		}
	}
	return nil
}

func (e *DestroyEffect) GetDescription() string {
	if e.NoRegeneration {
		return "destroy target. It can't be regenerated"
	}
	return "destroy target"
}

// ========================================
// Life Gain/Loss Effects
// ========================================

// GainLifeEffect has a player gain life
type GainLifeEffect struct {
	Amount int
}

func NewGainLifeEffect(amount int) *GainLifeEffect {
	return &GainLifeEffect{Amount: amount}
}

func (e *GainLifeEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Implement life gain
	return fmt.Errorf("life gain not yet implemented")
}

func (e *GainLifeEffect) GetDescription() string {
	return fmt.Sprintf("gain %d life", e.Amount)
}

// LoseLifeEffect has a player lose life
type LoseLifeEffect struct {
	Amount int
}

func NewLoseLifeEffect(amount int) *LoseLifeEffect {
	return &LoseLifeEffect{Amount: amount}
}

func (e *LoseLifeEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Implement life loss
	return fmt.Errorf("life loss not yet implemented")
}

func (e *LoseLifeEffect) GetDescription() string {
	return fmt.Sprintf("lose %d life", e.Amount)
}

// ========================================
// P/T Modification Effects
// ========================================

// BoostEffect modifies power/toughness
type BoostEffect struct {
	Power     int
	Toughness int
	Duration  Duration
}

// Duration specifies how long a continuous effect lasts
// Java: mage.constants.Duration
// MTG Rules: 611 (Continuous Effects)
type Duration int

const (
	// DurationUntilEndOfTurn lasts until end of turn (most common)
	DurationUntilEndOfTurn Duration = iota

	// DurationPermanent lasts forever (until removed)
	DurationPermanent

	// DurationWhileOnBattlefield lasts while the source is on battlefield
	DurationWhileOnBattlefield

	// DurationUntilEndOfCombat lasts until end of combat
	DurationUntilEndOfCombat

	// DurationEndOfTurn alias for DurationUntilEndOfTurn
	DurationEndOfTurn = DurationUntilEndOfTurn

	// DurationEndOfCombat alias for DurationUntilEndOfCombat
	DurationEndOfCombat = DurationUntilEndOfCombat

	// DurationUntilYourNextTurn lasts until your next turn
	DurationUntilYourNextTurn Duration = iota + 3

	// DurationWhileInGraveyard lasts while in graveyard
	DurationWhileInGraveyard

	// DurationWhileInHand lasts while in hand
	DurationWhileInHand

	// DurationWhileInExile lasts while in exile
	DurationWhileInExile

	// DurationCustom for special durations
	DurationCustom
)

func NewBoostEffect(power, toughness int) *BoostEffect {
	return &BoostEffect{
		Power:     power,
		Toughness: toughness,
		Duration:  DurationUntilEndOfTurn,
	}
}

func (e *BoostEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Implement boost effect via continuous effects system
	return fmt.Errorf("boost effect not yet implemented")
}

func (e *BoostEffect) GetDescription() string {
	powerStr := formatBoost(e.Power)
	toughnessStr := formatBoost(e.Toughness)
	return fmt.Sprintf("gets %s/%s", powerStr, toughnessStr)
}

func formatBoost(value int) string {
	if value >= 0 {
		return fmt.Sprintf("+%d", value)
	}
	return fmt.Sprintf("%d", value)
}

// ========================================
// Tap/Untap Effects
// ========================================

// TapEffect taps target permanents
type TapEffect struct{}

func NewTapEffect() *TapEffect {
	return &TapEffect{}
}

func (e *TapEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	for _, target := range targets {
		if err := game.TapPermanent(target); err != nil {
			return fmt.Errorf("failed to tap permanent: %w", err)
		}
	}
	return nil
}

func (e *TapEffect) GetDescription() string {
	return "tap target"
}

// UntapEffect untaps target permanents
type UntapEffect struct{}

func NewUntapEffect() *UntapEffect {
	return &UntapEffect{}
}

func (e *UntapEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	for _, target := range targets {
		if err := game.UntapPermanent(target); err != nil {
			return fmt.Errorf("failed to untap permanent: %w", err)
		}
	}
	return nil
}

func (e *UntapEffect) GetDescription() string {
	return "untap target"
}

// ========================================
// Mana Effects
// ========================================

// AddManaEffect adds mana to a player's mana pool
type AddManaEffect struct {
	Mana *Mana
}

func NewAddManaEffect(mana *Mana) *AddManaEffect {
	return &AddManaEffect{Mana: mana}
}

func (e *AddManaEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// If no targets, add to controller's pool
	if len(targets) == 0 {
		// TODO: Get controller from source
		return fmt.Errorf("add mana effect needs target player")
	}

	for _, target := range targets {
		if err := game.AddMana(target, e.Mana); err != nil {
			return fmt.Errorf("failed to add mana: %w", err)
		}
	}
	return nil
}

func (e *AddManaEffect) GetDescription() string {
	return "add mana"
}

// ========================================
// Counter Spell Effect
// ========================================

// CounterSpellEffect counters a spell
type CounterSpellEffect struct{}

func NewCounterSpellEffect() *CounterSpellEffect {
	return &CounterSpellEffect{}
}

func (e *CounterSpellEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Implement counter spell
	return fmt.Errorf("counter spell not yet implemented")
}

func (e *CounterSpellEffect) GetDescription() string {
	return "counter target spell"
}

// ========================================
// Effect List (composite effect)
// ========================================

// EffectList is a list of effects that execute in sequence
type EffectList struct {
	Effects []Effect
}

func NewEffectList(effects ...Effect) *EffectList {
	return &EffectList{Effects: effects}
}

func (e *EffectList) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	for _, effect := range e.Effects {
		if err := effect.Apply(ctx, game, source, targets); err != nil {
			return err
		}
	}
	return nil
}

func (e *EffectList) GetDescription() string {
	if len(e.Effects) == 0 {
		return ""
	}
	if len(e.Effects) == 1 {
		return e.Effects[0].GetDescription()
	}

	desc := e.Effects[0].GetDescription()
	for i := 1; i < len(e.Effects); i++ {
		desc += ". " + e.Effects[i].GetDescription()
	}
	return desc
}

func (e *EffectList) AddEffect(effect Effect) {
	e.Effects = append(e.Effects, effect)
}
