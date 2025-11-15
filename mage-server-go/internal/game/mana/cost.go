package mana

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ManaCost represents a parsed mana cost.
type ManaCost struct {
	Generic   int
	White     int
	Blue      int
	Black     int
	Red       int
	Green     int
	Colorless int
	X         bool // X in cost (e.g., {X}{R})
	Hybrid    []HybridCost
}

// HybridCost represents a hybrid mana cost (e.g., {W/U}, {2/B}).
type HybridCost struct {
	Options [][]ManaType // Each option is a list of mana types that can pay for it
}

// ParseCost parses a mana cost string (e.g., "{1}{G}", "{2}{R}{R}", "{X}{R}").
// Supports:
// - Generic: {1}, {2}, {3}, etc.
// - Colored: {W}, {U}, {B}, {R}, {G}, {C}
// - X costs: {X}
// - Hybrid: {W/U}, {2/B}, etc. (basic support)
func ParseCost(costStr string) (*ManaCost, error) {
	if costStr == "" {
		return &ManaCost{}, nil
	}

	cost := &ManaCost{}

	// Pattern to match mana symbols: {1}, {G}, {X}, {W/U}, etc.
	pattern := regexp.MustCompile(`\{([^}]+)\}`)
	matches := pattern.FindAllStringSubmatch(costStr, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		symbol := strings.ToUpper(strings.TrimSpace(match[1]))

		switch symbol {
		case "X":
			cost.X = true
		case "W":
			cost.White++
		case "U":
			cost.Blue++
		case "B":
			cost.Black++
		case "R":
			cost.Red++
		case "G":
			cost.Green++
		case "C":
			cost.Colorless++
		default:
			// Check if it's a number (generic mana)
			if num, err := strconv.Atoi(symbol); err == nil {
				cost.Generic += num
			} else if strings.Contains(symbol, "/") {
				// Hybrid mana: {W/U}, {2/B}, etc.
				hybrid := parseHybridCost(symbol)
				if hybrid != nil {
					cost.Hybrid = append(cost.Hybrid, *hybrid)
				}
			} else {
				return nil, fmt.Errorf("unknown mana symbol: {%s}", symbol)
			}
		}
	}

	return cost, nil
}

// parseHybridCost parses a hybrid mana symbol like "W/U" or "2/B".
func parseHybridCost(symbol string) *HybridCost {
	parts := strings.Split(symbol, "/")
	if len(parts) != 2 {
		return nil
	}

	left := strings.ToUpper(strings.TrimSpace(parts[0]))
	right := strings.ToUpper(strings.TrimSpace(parts[1]))

	hybrid := &HybridCost{Options: [][]ManaType{}}

	// Parse left side
	leftTypes := parseManaTypes(left)
	rightTypes := parseManaTypes(right)

	// Each option can be paid with either left or right
	if len(leftTypes) > 0 {
		hybrid.Options = append(hybrid.Options, leftTypes)
	}
	if len(rightTypes) > 0 {
		hybrid.Options = append(hybrid.Options, rightTypes)
	}

	return hybrid
}

// parseManaTypes parses a mana type string (e.g., "W", "2", "B").
func parseManaTypes(s string) []ManaType {
	var types []ManaType

	switch s {
	case "W":
		types = append(types, ManaWhite)
	case "U":
		types = append(types, ManaBlue)
	case "B":
		types = append(types, ManaBlack)
	case "R":
		types = append(types, ManaRed)
	case "G":
		types = append(types, ManaGreen)
	case "C":
		types = append(types, ManaColorless)
	default:
		// If it's a number, it can be paid with generic mana
		if num, err := strconv.Atoi(s); err == nil && num > 0 {
			// Generic mana can be paid with any type
			types = append(types, ManaGeneric)
		}
	}

	return types
}

// String returns a string representation of the mana cost.
func (mc *ManaCost) String() string {
	var parts []string

	if mc.X {
		parts = append(parts, "{X}")
	}

	for i := 0; i < mc.Generic; i++ {
		parts = append(parts, fmt.Sprintf("{%d}", i+1))
	}
	for i := 0; i < mc.White; i++ {
		parts = append(parts, "{W}")
	}
	for i := 0; i < mc.Blue; i++ {
		parts = append(parts, "{U}")
	}
	for i := 0; i < mc.Black; i++ {
		parts = append(parts, "{B}")
	}
	for i := 0; i < mc.Red; i++ {
		parts = append(parts, "{R}")
	}
	for i := 0; i < mc.Green; i++ {
		parts = append(parts, "{G}")
	}
	for i := 0; i < mc.Colorless; i++ {
		parts = append(parts, "{C}")
	}

	for _, hybrid := range mc.Hybrid {
		// Simple representation - full implementation would show both options
		if len(hybrid.Options) > 0 && len(hybrid.Options[0]) > 0 {
			parts = append(parts, fmt.Sprintf("{%s}", hybrid.Options[0][0]))
		}
	}

	return strings.Join(parts, "")
}

// GetTotalGeneric returns the total generic mana required (including hybrid costs).
func (mc *ManaCost) GetTotalGeneric() int {
	total := mc.Generic
	// Hybrid costs add to generic requirement
	for range mc.Hybrid {
		total++
	}
	return total
}

// CanPay checks if a mana pool can pay for this cost (without X costs).
// X costs require special handling during payment.
func (mc *ManaCost) CanPay(pool *ManaPool, xValue int) bool {
	// For X costs, we need the X value to be set
	if mc.X && xValue < 0 {
		return false
	}

	// Check if we have enough colored mana
	if pool.GetTotal(ManaWhite) < mc.White {
		return false
	}
	if pool.GetTotal(ManaBlue) < mc.Blue {
		return false
	}
	if pool.GetTotal(ManaBlack) < mc.Black {
		return false
	}
	if pool.GetTotal(ManaRed) < mc.Red {
		return false
	}
	if pool.GetTotal(ManaGreen) < mc.Green {
		return false
	}
	if pool.GetTotal(ManaColorless) < mc.Colorless {
		return false
	}

	// Pay hybrid costs (simplified - full implementation would try all combinations)
	// For now, check if we can pay at least one option for each hybrid
	for _, hybrid := range mc.Hybrid {
		canPayHybrid := false
		for _, option := range hybrid.Options {
			if len(option) > 0 {
				// Check if we can pay with this option
				canPay := false
				for _, mt := range option {
					if mt == ManaGeneric {
						// Generic can be paid with any mana
						canPay = pool.GetTotalMana() > 0
					} else {
						canPay = pool.GetTotal(mt) > 0
					}
					if canPay {
						break
					}
				}
				if canPay {
					canPayHybrid = true
					break
				}
			}
		}
		if !canPayHybrid {
			return false
		}
	}

	// Calculate total generic requirement (including hybrid costs)
	totalGeneric := mc.Generic
	if mc.X {
		totalGeneric += xValue
	}
	// Each hybrid cost adds 1 to generic requirement
	totalGeneric += len(mc.Hybrid)

	// Calculate available mana after paying colored requirements
	// We need to ensure we have enough total mana AND enough of each specific color
	totalRequired := mc.White + mc.Blue + mc.Black + mc.Red + mc.Green + mc.Colorless + len(mc.Hybrid) + totalGeneric
	totalAvailable := pool.GetTotalMana()

	if totalAvailable < totalRequired {
		return false
	}

	// Verify we can actually pay: colored requirements must be satisfied, then generic can use any remaining
	// This is a simplified check - full implementation would simulate actual payment
	return true
}

// ApplyReduction applies a cost reduction to this mana cost.
func (mc *ManaCost) ApplyReduction(genericReduction int, coloredReduction map[ManaType]int) *ManaCost {
	reduced := &ManaCost{
		Generic:   mc.Generic,
		White:     mc.White,
		Blue:      mc.Blue,
		Black:     mc.Black,
		Red:       mc.Red,
		Green:     mc.Green,
		Colorless: mc.Colorless,
		X:         mc.X,
		Hybrid:    mc.Hybrid, // Hybrid costs don't get reduced
	}

	// Apply generic reduction
	reduced.Generic -= genericReduction
	if reduced.Generic < 0 {
		reduced.Generic = 0
	}

	// Apply colored reductions
	if coloredReduction != nil {
		if red, ok := coloredReduction[ManaWhite]; ok {
			reduced.White -= red
			if reduced.White < 0 {
				reduced.White = 0
			}
		}
		if red, ok := coloredReduction[ManaBlue]; ok {
			reduced.Blue -= red
			if reduced.Blue < 0 {
				reduced.Blue = 0
			}
		}
		if red, ok := coloredReduction[ManaBlack]; ok {
			reduced.Black -= red
			if reduced.Black < 0 {
				reduced.Black = 0
			}
		}
		if red, ok := coloredReduction[ManaRed]; ok {
			reduced.Red -= red
			if reduced.Red < 0 {
				reduced.Red = 0
			}
		}
		if red, ok := coloredReduction[ManaGreen]; ok {
			reduced.Green -= red
			if reduced.Green < 0 {
				reduced.Green = 0
			}
		}
		if red, ok := coloredReduction[ManaColorless]; ok {
			reduced.Colorless -= red
			if reduced.Colorless < 0 {
				reduced.Colorless = 0
			}
		}
	}

	return reduced
}
