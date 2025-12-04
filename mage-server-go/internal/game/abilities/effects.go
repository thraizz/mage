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
	Amount       DynamicValue
	CombatDamage bool // true if this is combat damage (for lifelink, deathtouch, etc.)
}

// NewDamageEffect creates a damage effect with optional combat damage flag
// Java signature: new DamageTargetEffect(amount) or new DamageTargetEffect(amount, true) for combat damage
// amount can be int, DynamicValue, or any value that implements DynamicValue
// Can be called with no arguments for placeholder (transpiler compatibility)
func NewDamageEffect(args ...interface{}) *DamageEffect {
	effect := &DamageEffect{
		Amount:       NewStaticValue(0), // Default placeholder
		CombatDamage: false,
	}
	for i, arg := range args {
		switch v := arg.(type) {
		case int, DynamicValue:
			if i == 0 {
				effect.Amount = toDynamicValue(v)
			}
		case bool:
			effect.CombatDamage = v
		}
	}
	return effect
}

// NewCombatDamageEffect creates a damage effect that counts as combat damage
func NewCombatDamageEffect(amount interface{}) *DamageEffect {
	return &DamageEffect{Amount: toDynamicValue(amount), CombatDamage: true}
}

func (e *DamageEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	amount := e.Amount.Calculate(ctx, game, source)
	for _, target := range targets {
		if err := game.DealDamage(source, target, amount); err != nil {
			return fmt.Errorf("failed to deal damage: %w", err)
		}
	}
	return nil
}

func (e *DamageEffect) GetDescription() string {
	if e.CombatDamage {
		return fmt.Sprintf("deals %s combat damage", e.Amount.GetMessage())
	}
	return fmt.Sprintf("deals %s damage", e.Amount.GetMessage())
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
	Amount DynamicValue
}

// NewGainLifeEffect creates a gain life effect.
// Can be called with no arguments (returns placeholder), an int, or a DynamicValue.
func NewGainLifeEffect(amount ...interface{}) *GainLifeEffect {
	if len(amount) == 0 {
		return &GainLifeEffect{Amount: NewStaticValue(0)} // Placeholder
	}
	return &GainLifeEffect{Amount: toDynamicValue(amount[0])}
}

func (e *GainLifeEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Get the controller of the source to gain life
	controllerID, err := game.GetControllerID(source)
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	amount := e.Amount.Calculate(ctx, game, source)
	return game.GainLife(controllerID, amount)
}

func (e *GainLifeEffect) GetDescription() string {
	return fmt.Sprintf("gain %s life", e.Amount.GetMessage())
}

// LoseLifeEffect has a player lose life
type LoseLifeEffect struct {
	Amount DynamicValue
}

// NewLoseLifeEffect creates a lose life effect.
// Can be called with no arguments (returns placeholder), an int, or a DynamicValue.
func NewLoseLifeEffect(amount ...interface{}) *LoseLifeEffect {
	if len(amount) == 0 {
		return &LoseLifeEffect{Amount: NewStaticValue(0)} // Placeholder
	}
	return &LoseLifeEffect{Amount: toDynamicValue(amount[0])}
}

func (e *LoseLifeEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Implement life loss
	return fmt.Errorf("life loss not yet implemented")
}

func (e *LoseLifeEffect) GetDescription() string {
	return fmt.Sprintf("lose %s life", e.Amount.GetMessage())
}

// ========================================
// P/T Modification Effects
// ========================================

// BoostEffect modifies power/toughness
// Supports both static int values and dynamic values
type BoostEffect struct {
	Power          int
	Toughness      int
	PowerValue     DynamicValue // Alternative to static Power
	ToughnessValue DynamicValue // Alternative to static Toughness
	Duration       Duration
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

	// DurationWhileControlled lasts while controlled
	DurationWhileControlled

	// DurationWhileSourceCosts applies while source card costs something
	DurationWhileSourceCosts

	// DurationEndOfGame lasts until end of game
	DurationEndOfGame

	// DurationUntilSourceLeavesBattlefield lasts until source leaves battlefield
	DurationUntilSourceLeavesBattlefield
)

// NewBoostEffect creates a boost effect with static values
// Java: new BoostSourceEffect(power, toughness, Duration.EndOfTurn)
// Can be called with no arguments (placeholder) or with power/toughness values
// Additional optional params are ignored (for compatibility with generated code)
func NewBoostEffect(args ...interface{}) *BoostEffect {
	effect := &BoostEffect{
		Duration: DurationUntilEndOfTurn,
	}

	// Handle variadic args: power, toughness, and optional extras
	if len(args) >= 2 {
		// Handle power - can be int or DynamicValue
		switch p := args[0].(type) {
		case int:
			effect.Power = p
		case DynamicValue:
			effect.PowerValue = p
		}

		// Handle toughness - can be int or DynamicValue
		switch t := args[1].(type) {
		case int:
			effect.Toughness = t
		case DynamicValue:
			effect.ToughnessValue = t
		}
	}
	// Other args (duration, condition, etc.) ignored for now

	return effect
}

// NewBoostEffectWithDuration creates a boost effect with explicit duration
func NewBoostEffectWithDuration(power, toughness int, duration Duration) *BoostEffect {
	return &BoostEffect{
		Power:     power,
		Toughness: toughness,
		Duration:  duration,
	}
}

func (e *BoostEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// TODO: Implement boost effect via continuous effects system
	return fmt.Errorf("boost effect not yet implemented")
}

func (e *BoostEffect) GetDescription() string {
	// If using dynamic values, show a descriptive message
	if e.PowerValue != nil || e.ToughnessValue != nil {
		return "gets power/toughness boost based on dynamic values"
	}
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
type TapEffect struct {
	filter TargetFilter // Optional filter for targets
}

// NewTapEffect creates a new tap effect with optional filter.
// Java: new TapAllEffect(filter) or TapTargetEffect()
// Accepts TargetFilter or string (string values are ignored for compatibility)
func NewTapEffect(args ...interface{}) *TapEffect {
	e := &TapEffect{}
	for _, arg := range args {
		if f, ok := arg.(TargetFilter); ok {
			e.filter = f
			break
		}
		// Ignore string arguments - these are text descriptions from transpiler
	}
	return e
}

func (e *TapEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	for _, target := range targets {
		// TODO: Apply filter if present
		if err := game.TapPermanent(target); err != nil {
			return fmt.Errorf("failed to tap permanent: %w", err)
		}
	}
	return nil
}

func (e *TapEffect) GetDescription() string {
	if e.filter != nil {
		return fmt.Sprintf("tap each %s", e.filter.GetDescription())
	}
	return "tap target"
}

// UntapEffect untaps target permanents
type UntapEffect struct {
	text string
}

// NewUntapEffect creates a new untap effect with optional text description.
// Java: new UntapTargetEffect() or new UntapTargetEffect("untap it")
// Accepts any type but only uses string values for text
func NewUntapEffect(args ...interface{}) *UntapEffect {
	e := &UntapEffect{}
	for _, arg := range args {
		if s, ok := arg.(string); ok && s != "" {
			e.text = s
			break
		}
	}
	return e
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
	if e.text != "" {
		return e.text
	}
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

// NewCounterSpellEffect creates a counter spell effect.
// Optional args are ignored (for transpiler compatibility).
func NewCounterSpellEffect(args ...interface{}) *CounterSpellEffect {
	_ = args
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
