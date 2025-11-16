package game

// This file demonstrates a simple ability system design for the Go engine.
// This is a REFERENCE IMPLEMENTATION showing one possible approach.
// Actual implementation should be tailored to your specific needs.

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

// ============================================================================
// ABILITY INTERFACES
// ============================================================================

// Ability is the base interface for all card abilities
type Ability interface {
	// GetID returns unique identifier for this ability instance
	GetID() uuid.UUID
	
	// GetType returns the ability type (activated, triggered, static, etc.)
	GetType() AbilityType
	
	// CanExecute checks if this ability can be executed in current game state
	CanExecute(ctx context.Context, game *GameEngine, source *CardInstance) bool
	
	// Execute performs the ability's effect
	Execute(ctx context.Context, game *GameEngine, source *CardInstance, targets []Target) error
}

// AbilityType defines the category of ability
type AbilityType string

const (
	AbilityTypeActivated AbilityType = "activated"
	AbilityTypeTriggered AbilityType = "triggered"
	AbilityTypeStatic    AbilityType = "static"
	AbilityTypeSpell     AbilityType = "spell"
	AbilityTypeMana      AbilityType = "mana"
)

// ============================================================================
// ACTIVATED ABILITIES
// ============================================================================

// ActivatedAbility represents abilities that can be activated by paying a cost
// Example: "{T}: Deal 1 damage to any target"
type ActivatedAbility struct {
	ID          uuid.UUID
	Cost        Cost
	Effects     []Effect
	Targets     []TargetRequirement
	Timing      ActivationTiming
	UsesStack   bool
}

func (a *ActivatedAbility) GetID() uuid.UUID {
	return a.ID
}

func (a *ActivatedAbility) GetType() AbilityType {
	return AbilityTypeActivated
}

func (a *ActivatedAbility) CanExecute(ctx context.Context, game *GameEngine, source *CardInstance) bool {
	// Check if cost can be paid
	if !a.Cost.CanPay(game, source.ControllerID) {
		return false
	}
	
	// Check timing restrictions
	if !a.Timing.IsValid(game) {
		return false
	}
	
	// Check if targets are available
	for _, targetReq := range a.Targets {
		if !targetReq.HasValidTargets(game) {
			return false
		}
	}
	
	return true
}

func (a *ActivatedAbility) Execute(ctx context.Context, game *GameEngine, source *CardInstance, targets []Target) error {
	// Pay costs
	if err := a.Cost.Pay(game, source.ControllerID); err != nil {
		return fmt.Errorf("failed to pay cost: %w", err)
	}
	
	// If uses stack, create stack object
	if a.UsesStack {
		stackObj := &StackObject{
			ID:         uuid.New(),
			Source:     source,
			Ability:    a,
			Targets:    targets,
			Controller: source.ControllerID,
		}
		game.Stack.Push(stackObj)
		return nil
	}
	
	// Otherwise, resolve immediately (e.g., mana abilities)
	return a.Resolve(ctx, game, source, targets)
}

func (a *ActivatedAbility) Resolve(ctx context.Context, game *GameEngine, source *CardInstance, targets []Target) error {
	// Execute each effect in order
	for _, effect := range a.Effects {
		if err := effect.Apply(ctx, game, source, targets); err != nil {
			return fmt.Errorf("effect failed: %w", err)
		}
	}
	return nil
}

// ============================================================================
// TRIGGERED ABILITIES
// ============================================================================

// TriggeredAbility represents abilities that trigger on game events
// Example: "Whenever ~ attacks, draw a card"
type TriggeredAbility struct {
	ID          uuid.UUID
	Trigger     Trigger
	Condition   Condition
	Effects     []Effect
	Targets     []TargetRequirement
	Optional    bool
}

func (a *TriggeredAbility) GetID() uuid.UUID {
	return a.ID
}

func (a *TriggeredAbility) GetType() AbilityType {
	return AbilityTypeTriggered
}

func (a *TriggeredAbility) CanExecute(ctx context.Context, game *GameEngine, source *CardInstance) bool {
	// Triggered abilities are queued automatically, not manually activated
	return false
}

func (a *TriggeredAbility) Execute(ctx context.Context, game *GameEngine, source *CardInstance, targets []Target) error {
	// Check condition
	if a.Condition != nil && !a.Condition.Check(game, source) {
		return nil // Condition not met, don't trigger
	}
	
	// If optional, ask controller
	if a.Optional {
		// TODO: Prompt player
	}
	
	// Create triggered ability on stack
	stackObj := &StackObject{
		ID:         uuid.New(),
		Source:     source,
		Ability:    a,
		Targets:    targets,
		Controller: source.ControllerID,
	}
	game.Stack.Push(stackObj)
	
	return nil
}

func (a *TriggeredAbility) Resolve(ctx context.Context, game *GameEngine, source *CardInstance, targets []Target) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(ctx, game, source, targets); err != nil {
			return err
		}
	}
	return nil
}

// ============================================================================
// STATIC ABILITIES
// ============================================================================

// StaticAbility represents continuous effects
// Example: "Creatures you control get +1/+1"
type StaticAbility struct {
	ID          uuid.UUID
	Effect      ContinuousEffect
	Condition   Condition
	Layer       EffectLayer
}

func (a *StaticAbility) GetID() uuid.UUID {
	return a.ID
}

func (a *StaticAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *StaticAbility) CanExecute(ctx context.Context, game *GameEngine, source *CardInstance) bool {
	// Static abilities are always active, not executed
	return false
}

func (a *StaticAbility) Execute(ctx context.Context, game *GameEngine, source *CardInstance, targets []Target) error {
	// Static abilities don't execute, they apply continuously
	return fmt.Errorf("static abilities cannot be executed")
}

func (a *StaticAbility) Apply(game *GameEngine, source *CardInstance) {
	// Check condition
	if a.Condition != nil && !a.Condition.Check(game, source) {
		return
	}
	
	// Apply continuous effect
	a.Effect.Apply(game, source)
}

// ============================================================================
// EFFECTS
// ============================================================================

// Effect represents a game action that modifies state
type Effect interface {
	Apply(ctx context.Context, game *GameEngine, source *CardInstance, targets []Target) error
}

// DamageEffect deals damage to targets
type DamageEffect struct {
	Amount       int
	SourceDamage bool // If true, use source's power
}

func (e *DamageEffect) Apply(ctx context.Context, game *GameEngine, source *CardInstance, targets []Target) error {
	amount := e.Amount
	if e.SourceDamage {
		amount = source.Power
	}
	
	for _, target := range targets {
		if err := target.DealDamage(game, amount, source); err != nil {
			return err
		}
	}
	
	return nil
}

// DrawCardsEffect draws cards
type DrawCardsEffect struct {
	Amount int
}

func (e *DrawCardsEffect) Apply(ctx context.Context, game *GameEngine, source *CardInstance, targets []Target) error {
	for _, target := range targets {
		player, ok := target.(*PlayerTarget)
		if !ok {
			continue
		}
		
		for i := 0; i < e.Amount; i++ {
			if err := game.DrawCard(player.PlayerID); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// CounterSpellEffect counters a spell on the stack
type CounterSpellEffect struct{}

func (e *CounterSpellEffect) Apply(ctx context.Context, game *GameEngine, source *CardInstance, targets []Target) error {
	for _, target := range targets {
		stackTarget, ok := target.(*StackTarget)
		if !ok {
			continue
		}
		
		// Remove from stack and put in graveyard
		game.Stack.Remove(stackTarget.StackObjectID)
		// TODO: Move to graveyard
	}
	
	return nil
}

// ============================================================================
// COSTS
// ============================================================================

// Cost represents the cost to activate an ability
type Cost interface {
	CanPay(game *GameEngine, playerID uuid.UUID) bool
	Pay(game *GameEngine, playerID uuid.UUID) error
	String() string
}

// ManaCost represents a mana payment cost
type ManaCost struct {
	Generic   int
	White     int
	Blue      int
	Black     int
	Red       int
	Green     int
	Colorless int
}

func (c *ManaCost) CanPay(game *GameEngine, playerID uuid.UUID) bool {
	player := game.GetPlayer(playerID)
	if player == nil {
		return false
	}
	
	// Check if player has enough mana in pool
	return player.ManaPool.CanPay(c)
}

func (c *ManaCost) Pay(game *GameEngine, playerID uuid.UUID) error {
	player := game.GetPlayer(playerID)
	if player == nil {
		return fmt.Errorf("player not found")
	}
	
	return player.ManaPool.Pay(c)
}

func (c *ManaCost) String() string {
	result := ""
	if c.Generic > 0 {
		result += fmt.Sprintf("{%d}", c.Generic)
	}
	for i := 0; i < c.White; i++ {
		result += "{W}"
	}
	for i := 0; i < c.Blue; i++ {
		result += "{U}"
	}
	for i := 0; i < c.Black; i++ {
		result += "{B}"
	}
	for i := 0; i < c.Red; i++ {
		result += "{R}"
	}
	for i := 0; i < c.Green; i++ {
		result += "{G}"
	}
	for i := 0; i < c.Colorless; i++ {
		result += "{C}"
	}
	return result
}

// TapCost represents tapping the source permanent
type TapCost struct{}

func (c *TapCost) CanPay(game *GameEngine, playerID uuid.UUID) bool {
	// Check if source is untapped (handled by caller)
	return true
}

func (c *TapCost) Pay(game *GameEngine, playerID uuid.UUID) error {
	// Tap the source (handled by caller)
	return nil
}

func (c *TapCost) String() string {
	return "{T}"
}

// CompositeCost combines multiple costs
type CompositeCost struct {
	Costs []Cost
}

func (c *CompositeCost) CanPay(game *GameEngine, playerID uuid.UUID) bool {
	for _, cost := range c.Costs {
		if !cost.CanPay(game, playerID) {
			return false
		}
	}
	return true
}

func (c *CompositeCost) Pay(game *GameEngine, playerID uuid.UUID) error {
	for _, cost := range c.Costs {
		if err := cost.Pay(game, playerID); err != nil {
			return err
		}
	}
	return nil
}

func (c *CompositeCost) String() string {
	result := ""
	for _, cost := range c.Costs {
		result += cost.String()
	}
	return result
}

// ============================================================================
// TARGETS
// ============================================================================

// Target represents something that can be targeted
type Target interface {
	IsValid(game *GameEngine) bool
	DealDamage(game *GameEngine, amount int, source *CardInstance) error
}

// CardTarget represents a card as a target
type CardTarget struct {
	CardID uuid.UUID
}

func (t *CardTarget) IsValid(game *GameEngine) bool {
	card := game.GetCard(t.CardID)
	return card != nil
}

func (t *CardTarget) DealDamage(game *GameEngine, amount int, source *CardInstance) error {
	card := game.GetCard(t.CardID)
	if card == nil {
		return fmt.Errorf("target card not found")
	}
	
	card.Damage += amount
	return nil
}

// PlayerTarget represents a player as a target
type PlayerTarget struct {
	PlayerID uuid.UUID
}

func (t *PlayerTarget) IsValid(game *GameEngine) bool {
	player := game.GetPlayer(t.PlayerID)
	return player != nil
}

func (t *PlayerTarget) DealDamage(game *GameEngine, amount int, source *CardInstance) error {
	player := game.GetPlayer(t.PlayerID)
	if player == nil {
		return fmt.Errorf("target player not found")
	}
	
	player.Life -= amount
	return nil
}

// StackTarget represents a spell/ability on the stack
type StackTarget struct {
	StackObjectID uuid.UUID
}

func (t *StackTarget) IsValid(game *GameEngine) bool {
	return game.Stack.Contains(t.StackObjectID)
}

func (t *StackTarget) DealDamage(game *GameEngine, amount int, source *CardInstance) error {
	return fmt.Errorf("cannot deal damage to stack object")
}

// ============================================================================
// EXAMPLE CARD DEFINITIONS
// ============================================================================

// Example: Lightning Bolt
// {R}
// Instant
// Lightning Bolt deals 3 damage to any target.
func NewLightningBolt(id uuid.UUID, setInfo CardSetInfo) *CardInstance {
	card := &CardInstance{
		ID:        id,
		Name:      "Lightning Bolt",
		ManaCost:  "{R}",
		Types:     []CardType{TypeInstant},
		Abilities: []Ability{},
	}
	
	// Spell ability
	spellAbility := &ActivatedAbility{
		ID:   uuid.New(),
		Cost: &ManaCost{Red: 1},
		Effects: []Effect{
			&DamageEffect{Amount: 3},
		},
		Targets: []TargetRequirement{
			{
				Type:     TargetTypeAny,
				MinCount: 1,
				MaxCount: 1,
			},
		},
		Timing:    ActivationTimingStack,
		UsesStack: true,
	}
	
	card.Abilities = append(card.Abilities, spellAbility)
	
	return card
}

// Example: Prodigal Pyromancer (Tim)
// {2}{R}
// Creature — Human Wizard
// {T}: Prodigal Pyromancer deals 1 damage to any target.
// 1/1
func NewProdigalPyromancer(id uuid.UUID, setInfo CardSetInfo) *CardInstance {
	card := &CardInstance{
		ID:        id,
		Name:      "Prodigal Pyromancer",
		ManaCost:  "{2}{R}",
		Types:     []CardType{TypeCreature},
		Subtypes:  []string{"Human", "Wizard"},
		Power:     1,
		Toughness: 1,
		Abilities: []Ability{},
	}
	
	// Activated ability: {T}: Deal 1 damage
	ability := &ActivatedAbility{
		ID: uuid.New(),
		Cost: &CompositeCost{
			Costs: []Cost{&TapCost{}},
		},
		Effects: []Effect{
			&DamageEffect{Amount: 1},
		},
		Targets: []TargetRequirement{
			{
				Type:     TargetTypeAny,
				MinCount: 1,
				MaxCount: 1,
			},
		},
		Timing:    ActivationTimingMain,
		UsesStack: true,
	}
	
	card.Abilities = append(card.Abilities, ability)
	
	return card
}

// Example: Divination
// {2}{U}
// Sorcery
// Draw two cards.
func NewDivination(id uuid.UUID, setInfo CardSetInfo) *CardInstance {
	card := &CardInstance{
		ID:        id,
		Name:      "Divination",
		ManaCost:  "{2}{U}",
		Types:     []CardType{TypeSorcery},
		Abilities: []Ability{},
	}
	
	// Spell ability
	spellAbility := &ActivatedAbility{
		ID: uuid.New(),
		Cost: &ManaCost{
			Generic: 2,
			Blue:    1,
		},
		Effects: []Effect{
			&DrawCardsEffect{Amount: 2},
		},
		Timing:    ActivationTimingSorcery,
		UsesStack: true,
	}
	
	card.Abilities = append(card.Abilities, spellAbility)
	
	return card
}

// ============================================================================
// CARD REGISTRY
// ============================================================================

// CardFactory is a function that creates a card instance
type CardFactory func(id uuid.UUID, setInfo CardSetInfo) *CardInstance

// CardRegistry maps card class names to factory functions
var CardRegistry = map[string]CardFactory{
	"mage.cards.l.LightningBolt":      NewLightningBolt,
	"mage.cards.p.ProdigalPyromancer": NewProdigalPyromancer,
	"mage.cards.d.Divination":         NewDivination,
}

// RegisterCard adds a card factory to the registry
func RegisterCard(className string, factory CardFactory) {
	CardRegistry[className] = factory
}

// CreateCard instantiates a card by class name
func CreateCard(className string, id uuid.UUID, setInfo CardSetInfo) (*CardInstance, error) {
	factory, ok := CardRegistry[className]
	if !ok {
		return nil, fmt.Errorf("card class not found: %s", className)
	}
	
	return factory(id, setInfo), nil
}

// ============================================================================
// SUPPORTING TYPES (referenced above)
// ============================================================================

type CardSetInfo struct {
	Name       string
	SetCode    string
	CardNumber string
	Rarity     string
}

type ActivationTiming int

const (
	ActivationTimingAny ActivationTiming = iota
	ActivationTimingMain
	ActivationTimingSorcery
	ActivationTimingStack
	ActivationTimingInstant
)

func (t ActivationTiming) IsValid(game *GameEngine) bool {
	// TODO: Check game state for timing restrictions
	return true
}

type TargetRequirement struct {
	Type     TargetType
	MinCount int
	MaxCount int
	Filter   func(*CardInstance) bool
}

func (r *TargetRequirement) HasValidTargets(game *GameEngine) bool {
	// TODO: Check if valid targets exist
	return true
}

type TargetType int

const (
	TargetTypeAny TargetType = iota
	TargetTypeCreature
	TargetTypePlayer
	TargetTypePermanent
	TargetTypeSpell
)

type Trigger interface {
	Check(event GameEvent) bool
}

type Condition interface {
	Check(game *GameEngine, source *CardInstance) bool
}

type ContinuousEffect interface {
	Apply(game *GameEngine, source *CardInstance)
}

type EffectLayer int

const (
	LayerCopyEffects EffectLayer = iota
	LayerControlChanging
	LayerTextChanging
	LayerTypeChanging
	LayerColorChanging
	LayerAbilityAddingRemoving
	LayerPowerToughnessChanging
)
