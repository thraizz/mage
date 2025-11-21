package abilities

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Cost represents a cost that must be paid to activate an ability
type Cost interface {
	// CanPay checks if the player can pay this cost
	CanPay(ctx context.Context, game GameContext, playerID uuid.UUID) bool

	// Pay pays this cost
	Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error

	// String returns a text representation of the cost
	String() string
}

// ========================================
// Mana Cost
// ========================================

// ManaCost represents a mana cost
type ManaCost struct {
	Mana *Mana
}

func NewManaCost(mana *Mana) *ManaCost {
	return &ManaCost{Mana: mana}
}

// ParseManaCost parses a mana cost string like "{2}{U}{U}"
func ParseManaCost(costStr string) (*ManaCost, error) {
	mana := NewMana()

	// Remove braces and parse
	costStr = strings.ReplaceAll(costStr, "{", "")
	parts := strings.Split(costStr, "}")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		switch part {
		case "W":
			mana.White++
		case "U":
			mana.Blue++
		case "B":
			mana.Black++
		case "R":
			mana.Red++
		case "G":
			mana.Green++
		case "C":
			mana.Colorless++
		default:
			// Try to parse as generic mana
			var amount int
			if _, err := fmt.Sscanf(part, "%d", &amount); err == nil {
				mana.Generic += amount
			} else {
				return nil, fmt.Errorf("unknown mana symbol: %s", part)
			}
		}
	}

	return NewManaCost(mana), nil
}

func (c *ManaCost) CanPay(ctx context.Context, game GameContext, playerID uuid.UUID) bool {
	// TODO: Check if player has enough mana
	return true
}

func (c *ManaCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error {
	// TODO: Implement mana payment
	return fmt.Errorf("mana payment not yet implemented")
}

func (c *ManaCost) String() string {
	parts := []string{}

	if c.Mana.Generic > 0 {
		parts = append(parts, fmt.Sprintf("{%d}", c.Mana.Generic))
	}
	for i := 0; i < c.Mana.White; i++ {
		parts = append(parts, "{W}")
	}
	for i := 0; i < c.Mana.Blue; i++ {
		parts = append(parts, "{U}")
	}
	for i := 0; i < c.Mana.Black; i++ {
		parts = append(parts, "{B}")
	}
	for i := 0; i < c.Mana.Red; i++ {
		parts = append(parts, "{R}")
	}
	for i := 0; i < c.Mana.Green; i++ {
		parts = append(parts, "{G}")
	}
	for i := 0; i < c.Mana.Colorless; i++ {
		parts = append(parts, "{C}")
	}

	if len(parts) == 0 {
		return "{0}"
	}

	return strings.Join(parts, "")
}

// ========================================
// Tap Cost
// ========================================

// TapCost represents tapping the source as a cost
type TapCost struct{}

func NewTapCost() *TapCost {
	return &TapCost{}
}

func (c *TapCost) CanPay(ctx context.Context, game GameContext, playerID uuid.UUID) bool {
	// TODO: Check if source is tapped
	return true
}

func (c *TapCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error {
	// TODO: Tap the source
	return fmt.Errorf("tap cost not yet implemented")
}

func (c *TapCost) String() string {
	return "{T}"
}

// ========================================
// Sacrifice Cost
// ========================================

// SacrificeCost represents sacrificing permanents as a cost
type SacrificeCost struct {
	Amount int
	Filter string // e.g., "creature", "artifact", etc.
}

func NewSacrificeCost(amount int, filter string) *SacrificeCost {
	return &SacrificeCost{
		Amount: amount,
		Filter: filter,
	}
}

func (c *SacrificeCost) CanPay(ctx context.Context, game GameContext, playerID uuid.UUID) bool {
	// TODO: Check if player has enough permanents to sacrifice
	return true
}

func (c *SacrificeCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error {
	// TODO: Sacrifice permanents
	return fmt.Errorf("sacrifice cost not yet implemented")
}

func (c *SacrificeCost) String() string {
	if c.Amount == 1 {
		if c.Filter != "" {
			return fmt.Sprintf("Sacrifice a %s", c.Filter)
		}
		return "Sacrifice a permanent"
	}
	if c.Filter != "" {
		return fmt.Sprintf("Sacrifice %d %ss", c.Amount, c.Filter)
	}
	return fmt.Sprintf("Sacrifice %d permanents", c.Amount)
}

// ========================================
// Discard Cost
// ========================================

// DiscardCost represents discarding cards as a cost
type DiscardCost struct {
	Amount int
	Random bool
}

func NewDiscardCost(amount int) *DiscardCost {
	return &DiscardCost{Amount: amount, Random: false}
}

func NewDiscardCostRandom(amount int) *DiscardCost {
	return &DiscardCost{Amount: amount, Random: true}
}

func (c *DiscardCost) CanPay(ctx context.Context, game GameContext, playerID uuid.UUID) bool {
	// TODO: Check if player has enough cards to discard
	return true
}

func (c *DiscardCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error {
	// TODO: Discard cards
	return fmt.Errorf("discard cost not yet implemented")
}

func (c *DiscardCost) String() string {
	if c.Amount == 1 {
		if c.Random {
			return "Discard a card at random"
		}
		return "Discard a card"
	}
	if c.Random {
		return fmt.Sprintf("Discard %d cards at random", c.Amount)
	}
	return fmt.Sprintf("Discard %d cards", c.Amount)
}

// ========================================
// Pay Life Cost
// ========================================

// PayLifeCost represents paying life as a cost
type PayLifeCost struct {
	Amount int
}

func NewPayLifeCost(amount int) *PayLifeCost {
	return &PayLifeCost{Amount: amount}
}

func (c *PayLifeCost) CanPay(ctx context.Context, game GameContext, playerID uuid.UUID) bool {
	// TODO: Check if player has enough life
	return true
}

func (c *PayLifeCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error {
	// TODO: Pay life
	return fmt.Errorf("pay life cost not yet implemented")
}

func (c *PayLifeCost) String() string {
	if c.Amount == 1 {
		return "Pay 1 life"
	}
	return fmt.Sprintf("Pay %d life", c.Amount)
}

// ========================================
// Composite Cost
// ========================================

// CompositeCost is a list of costs that must all be paid
type CompositeCost struct {
	Costs []Cost
}

func NewCompositeCost(costs ...Cost) *CompositeCost {
	return &CompositeCost{Costs: costs}
}

func (c *CompositeCost) CanPay(ctx context.Context, game GameContext, playerID uuid.UUID) bool {
	for _, cost := range c.Costs {
		if !cost.CanPay(ctx, game, playerID) {
			return false
		}
	}
	return true
}

func (c *CompositeCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error {
	for _, cost := range c.Costs {
		if err := cost.Pay(ctx, game, playerID); err != nil {
			return err
		}
	}
	return nil
}

func (c *CompositeCost) String() string {
	if len(c.Costs) == 0 {
		return "{0}"
	}

	parts := make([]string, len(c.Costs))
	for i, cost := range c.Costs {
		parts[i] = cost.String()
	}

	return strings.Join(parts, ", ")
}

func (c *CompositeCost) AddCost(cost Cost) {
	c.Costs = append(c.Costs, cost)
}

// ========================================
// No Cost
// ========================================

// NoCost represents an ability with no cost
type NoCost struct{}

func NewNoCost() *NoCost {
	return &NoCost{}
}

func (c *NoCost) CanPay(ctx context.Context, game GameContext, playerID uuid.UUID) bool {
	return true
}

func (c *NoCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error {
	return nil
}

func (c *NoCost) String() string {
	return ""
}
