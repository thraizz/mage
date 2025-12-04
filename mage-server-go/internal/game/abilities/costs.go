package abilities

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

// internalCard is a forward declaration for the game engine's card type
// We need this to work with cards returned from GameContext
type internalCard struct {
	ID         string
	Name       string
	Type       string
	SubTypes   []string
	SuperTypes []string
}

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
	if c.Mana == nil {
		return true // No cost
	}

	// Get player's mana pool
	pool := game.GetManaPool(playerID)
	if pool == nil {
		return false
	}

	// Check colored mana requirements
	if pool.GetAmount("WHITE") < c.Mana.White {
		return false
	}
	if pool.GetAmount("BLUE") < c.Mana.Blue {
		return false
	}
	if pool.GetAmount("BLACK") < c.Mana.Black {
		return false
	}
	if pool.GetAmount("RED") < c.Mana.Red {
		return false
	}
	if pool.GetAmount("GREEN") < c.Mana.Green {
		return false
	}
	if pool.GetAmount("COLORLESS") < c.Mana.Colorless {
		return false
	}

	// Check generic mana requirement (can be paid with any mana)
	// Calculate total available mana after paying colored costs
	availableAfterColored := pool.GetAmount("WHITE") - c.Mana.White +
		pool.GetAmount("BLUE") - c.Mana.Blue +
		pool.GetAmount("BLACK") - c.Mana.Black +
		pool.GetAmount("RED") - c.Mana.Red +
		pool.GetAmount("GREEN") - c.Mana.Green +
		pool.GetAmount("COLORLESS") - c.Mana.Colorless

	if availableAfterColored < c.Mana.Generic {
		return false
	}

	return true
}

func (c *ManaCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error {
	if c.Mana == nil {
		return nil // No cost
	}

	// First check if we can pay
	if !c.CanPay(ctx, game, playerID) {
		return fmt.Errorf("insufficient mana to pay cost")
	}

	// Get player's mana pool
	pool := game.GetManaPool(playerID)
	if pool == nil {
		return fmt.Errorf("player has no mana pool")
	}

	// Pay colored mana first (these MUST be paid with the specific color)
	if c.Mana.White > 0 {
		if err := pool.SpendMana("WHITE", c.Mana.White); err != nil {
			return fmt.Errorf("failed to spend white mana: %w", err)
		}
	}
	if c.Mana.Blue > 0 {
		if err := pool.SpendMana("BLUE", c.Mana.Blue); err != nil {
			return fmt.Errorf("failed to spend blue mana: %w", err)
		}
	}
	if c.Mana.Black > 0 {
		if err := pool.SpendMana("BLACK", c.Mana.Black); err != nil {
			return fmt.Errorf("failed to spend black mana: %w", err)
		}
	}
	if c.Mana.Red > 0 {
		if err := pool.SpendMana("RED", c.Mana.Red); err != nil {
			return fmt.Errorf("failed to spend red mana: %w", err)
		}
	}
	if c.Mana.Green > 0 {
		if err := pool.SpendMana("GREEN", c.Mana.Green); err != nil {
			return fmt.Errorf("failed to spend green mana: %w", err)
		}
	}
	if c.Mana.Colorless > 0 {
		if err := pool.SpendMana("COLORLESS", c.Mana.Colorless); err != nil {
			return fmt.Errorf("failed to spend colorless mana: %w", err)
		}
	}

	// Pay generic mana (can be paid with any remaining mana)
	// We need to implement a strategy for paying generic costs
	if c.Mana.Generic > 0 {
		remaining := c.Mana.Generic

		// Try to pay with each color in order: W, U, B, R, G, C
		colors := []string{"WHITE", "BLUE", "BLACK", "RED", "GREEN", "COLORLESS"}
		for _, color := range colors {
			available := pool.GetAmount(color)
			if available > 0 {
				toPay := available
				if toPay > remaining {
					toPay = remaining
				}
				if err := pool.SpendMana(color, toPay); err != nil {
					return fmt.Errorf("failed to spend %s mana for generic cost: %w", color, err)
				}
				remaining -= toPay
				if remaining == 0 {
					break
				}
			}
		}

		if remaining > 0 {
			return fmt.Errorf("insufficient mana to pay generic cost")
		}
	}

	return nil
}

// ConvertedManaCost returns the total mana value (CMC)
func (c *ManaCost) ConvertedManaCost() int {
	if c.Mana == nil {
		return 0
	}
	return c.Mana.Total()
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

// TapCost represents tapping the source permanent as a cost
type TapCost struct {
	sourceID uuid.UUID
}

// NewTapCost creates a new tap cost (sourceID will be set by the ability builder)
func NewTapCost() *TapCost {
	return &TapCost{}
}

// NewTapCostWithSource creates a tap cost with a specific source permanent
func NewTapCostWithSource(sourceID uuid.UUID) *TapCost {
	return &TapCost{sourceID: sourceID}
}

// SetSource sets the source permanent ID for this tap cost
func (c *TapCost) SetSource(sourceID uuid.UUID) {
	c.sourceID = sourceID
}

// GetSource returns the source permanent ID
func (c *TapCost) GetSource() uuid.UUID {
	return c.sourceID
}

func (c *TapCost) CanPay(ctx context.Context, game GameContext, playerID uuid.UUID) bool {
	// If no source is set, we can't pay
	if c.sourceID == uuid.Nil {
		return false
	}

	// Check if the source permanent is already tapped
	return !game.IsPermanentTapped(c.sourceID)
}

func (c *TapCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error {
	// If no source is set, we can't pay
	if c.sourceID == uuid.Nil {
		return fmt.Errorf("tap cost has no source permanent set")
	}

	// Check if already tapped
	if game.IsPermanentTapped(c.sourceID) {
		return fmt.Errorf("permanent is already tapped")
	}

	// Tap the permanent
	return game.TapPermanent(c.sourceID)
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

// NewSacrificeSourceCost creates a cost for sacrificing the source permanent
func NewSacrificeSourceCost() *SacrificeCost {
	return &SacrificeCost{
		Amount: 1,
		Filter: "source",
	}
}

func (c *SacrificeCost) CanPay(ctx context.Context, game GameContext, playerID uuid.UUID) bool {
	// Get all permanents controlled by the player
	permanents, err := game.GetPermanentsControlledByPlayer(playerID)
	if err != nil {
		return false
	}

	// Special case: sacrificing source
	if c.Filter == "source" {
		// For source sacrifice, we need at least one permanent (the source itself)
		// This will be validated during actual payment
		return len(permanents) > 0
	}

	// Count permanents that match the filter
	matchCount := 0
	for _, perm := range permanents {
		card, ok := perm.(*internalCard)
		if !ok {
			continue
		}

		if permanentMatchesFilter(card, c.Filter) {
			matchCount++
		}
	}

	return matchCount >= c.Amount
}

func (c *SacrificeCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error {
	// First check if we can pay
	if !c.CanPay(ctx, game, playerID) {
		return fmt.Errorf("insufficient permanents to sacrifice")
	}

	// Get all permanents controlled by the player
	permanents, err := game.GetPermanentsControlledByPlayer(playerID)
	if err != nil {
		return fmt.Errorf("failed to get permanents: %w", err)
	}

	// Special case: sacrificing source
	// For now, we don't have a way to identify the source in this context
	// This will need to be handled by the ability activation system
	if c.Filter == "source" {
		// The ability activation system should pass the source ID somehow
		// For now, return an error indicating this needs special handling
		return fmt.Errorf("sacrifice source cost requires special handling by ability system")
	}

	// Collect permanents that match the filter
	var candidates []uuid.UUID
	for _, perm := range permanents {
		card, ok := perm.(*internalCard)
		if !ok {
			continue
		}

		if permanentMatchesFilter(card, c.Filter) {
			id, err := uuid.Parse(card.ID)
			if err != nil {
				continue
			}
			candidates = append(candidates, id)
		}
	}

	if len(candidates) < c.Amount {
		return fmt.Errorf("not enough permanents matching filter '%s'", c.Filter)
	}

	// TODO: In a real implementation, we'd need player input to choose which permanents to sacrifice
	// For now, just sacrifice the first N that match
	for i := 0; i < c.Amount && i < len(candidates); i++ {
		if err := game.SacrificePermanent(candidates[i]); err != nil {
			return fmt.Errorf("failed to sacrifice permanent: %w", err)
		}
	}

	return nil
}

func (c *SacrificeCost) String() string {
	if c.Amount == 1 {
		if c.Filter == "source" {
			return "Sacrifice this permanent"
		}
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
	// Get player's hand
	hand, err := game.GetPlayerHand(playerID)
	if err != nil {
		return false
	}

	// Check if player has enough cards to discard
	return len(hand) >= c.Amount
}

func (c *DiscardCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error {
	// First check if we can pay
	if !c.CanPay(ctx, game, playerID) {
		return fmt.Errorf("insufficient cards in hand to discard")
	}

	// Get player's hand
	hand, err := game.GetPlayerHand(playerID)
	if err != nil {
		return fmt.Errorf("failed to get player hand: %w", err)
	}

	// Collect card IDs
	var cardIDs []uuid.UUID
	for _, cardIface := range hand {
		card, ok := cardIface.(*internalCard)
		if !ok {
			continue
		}
		id, err := uuid.Parse(card.ID)
		if err != nil {
			continue
		}
		cardIDs = append(cardIDs, id)
	}

	if len(cardIDs) < c.Amount {
		return fmt.Errorf("not enough cards in hand")
	}

	// If random, shuffle the card IDs
	if c.Random {
		rand.Shuffle(len(cardIDs), func(i, j int) {
			cardIDs[i], cardIDs[j] = cardIDs[j], cardIDs[i]
		})
	}

	// TODO: In a real implementation, we'd need player input to choose which cards to discard
	// For now, just discard the first N cards (or N random cards if Random is true)
	for i := 0; i < c.Amount && i < len(cardIDs); i++ {
		if err := game.DiscardCard(playerID, cardIDs[i]); err != nil {
			return fmt.Errorf("failed to discard card: %w", err)
		}
	}

	return nil
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

// ========================================
// Discard Target Cost (with filtering)
// ========================================

// DiscardTargetCost represents discarding specific cards (with filtering) as a cost
// Java: mage.abilities.costs.common.DiscardTargetCost
// MTG Rules: 118.9 (Costs), 701.8 (Discard)
type DiscardTargetCost struct {
	amount int
	filter CardFilter // What types of cards can be discarded
}

// NewDiscardTargetCost creates a new discard target cost
// Example: NewDiscardTargetCost(2, NewArtifactCardFilter()) = "Discard two artifact cards"
func NewDiscardTargetCost(amount int, filter CardFilter) *DiscardTargetCost {
	return &DiscardTargetCost{
		amount: amount,
		filter: filter,
	}
}

// CanPay checks if the player can pay this cost
func (c *DiscardTargetCost) CanPay(ctx context.Context, game GameContext, playerID uuid.UUID) bool {
	// Get player's hand
	hand, err := game.GetPlayerHand(playerID)
	if err != nil {
		return false
	}

	// Count cards that match the filter
	matchCount := 0
	for _, cardIface := range hand {
		card, ok := cardIface.(*internalCard)
		if !ok {
			continue
		}

		if cardMatchesFilter(card, c.filter) {
			matchCount++
		}
	}

	return matchCount >= c.amount
}

// Pay pays this cost (discards the specified cards)
func (c *DiscardTargetCost) Pay(ctx context.Context, game GameContext, playerID uuid.UUID) error {
	// First check if we can pay
	if !c.CanPay(ctx, game, playerID) {
		return fmt.Errorf("insufficient cards matching filter to discard")
	}

	// Get player's hand
	hand, err := game.GetPlayerHand(playerID)
	if err != nil {
		return fmt.Errorf("failed to get player hand: %w", err)
	}

	// Collect cards that match the filter
	var candidates []uuid.UUID
	for _, cardIface := range hand {
		card, ok := cardIface.(*internalCard)
		if !ok {
			continue
		}

		if cardMatchesFilter(card, c.filter) {
			id, err := uuid.Parse(card.ID)
			if err != nil {
				continue
			}
			candidates = append(candidates, id)
		}
	}

	if len(candidates) < c.amount {
		return fmt.Errorf("not enough cards matching filter in hand")
	}

	// TODO: In a real implementation, we'd need player input to choose which cards to discard
	// For now, just discard the first N that match
	for i := 0; i < c.amount && i < len(candidates); i++ {
		if err := game.DiscardCard(playerID, candidates[i]); err != nil {
			return fmt.Errorf("failed to discard card: %w", err)
		}
	}

	return nil
}

// String returns a text representation
func (c *DiscardTargetCost) String() string {
	filterDesc := "card"
	if c.filter != nil {
		filterDesc = c.filter.GetDescription()
	}

	if c.amount == 1 {
		return fmt.Sprintf("Discard a %s", filterDesc)
	}
	return fmt.Sprintf("Discard %d %ss", c.amount, filterDesc)
}

// ========================================
// Helper Functions
// ========================================

// permanentMatchesFilter checks if a permanent matches a type filter string
func permanentMatchesFilter(card *internalCard, filter string) bool {
	if filter == "" || filter == "permanent" {
		return true
	}

	// Normalize filter to uppercase for comparison
	filterUpper := strings.ToUpper(filter)
	typeUpper := strings.ToUpper(card.Type)

	// Check main types
	if strings.Contains(typeUpper, filterUpper) {
		return true
	}

	// Check subtypes
	for _, subType := range card.SubTypes {
		if strings.EqualFold(subType, filter) {
			return true
		}
	}

	return false
}

// cardMatchesFilter checks if a card matches a card filter
// Note: This is a simplified version that checks based on card type
// The actual CardFilter.Matches() method requires a GameContext and UUID
func cardMatchesFilter(card *internalCard, filter CardFilter) bool {
	if filter == nil {
		return true
	}

	// Since we can't call filter.Matches() without a GameContext,
	// we do a simplified type check based on the filter's description
	// This is a temporary solution until we can pass GameContext through
	filterDesc := filter.GetDescription()
	typeUpper := strings.ToUpper(card.Type)

	// Check for artifact
	if strings.Contains(filterDesc, "artifact") || strings.Contains(filterDesc, "Artifact") {
		return strings.Contains(typeUpper, "ARTIFACT")
	}

	// Check for creature
	if strings.Contains(filterDesc, "creature") || strings.Contains(filterDesc, "Creature") {
		return strings.Contains(typeUpper, "CREATURE")
	}

	// Check for land
	if strings.Contains(filterDesc, "land") || strings.Contains(filterDesc, "Land") {
		return strings.Contains(typeUpper, "LAND")
	}

	// Check for instant
	if strings.Contains(filterDesc, "instant") || strings.Contains(filterDesc, "Instant") {
		return strings.Contains(typeUpper, "INSTANT")
	}

	// Check for sorcery
	if strings.Contains(filterDesc, "sorcery") || strings.Contains(filterDesc, "Sorcery") {
		return strings.Contains(typeUpper, "SORCERY")
	}

	// Check for enchantment
	if strings.Contains(filterDesc, "enchantment") || strings.Contains(filterDesc, "Enchantment") {
		return strings.Contains(typeUpper, "ENCHANTMENT")
	}

	// Check for planeswalker
	if strings.Contains(filterDesc, "planeswalker") || strings.Contains(filterDesc, "Planeswalker") {
		return strings.Contains(typeUpper, "PLANESWALKER")
	}

	// Default: accept any card
	return true
}
