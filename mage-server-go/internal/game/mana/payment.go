package mana

import (
	"fmt"
)

// PaymentPlan represents a plan for paying a mana cost.
type PaymentPlan struct {
	White     int
	Blue      int
	Black     int
	Red       int
	Green     int
	Colorless int
	Generic   int // Generic mana can be paid with any type
	XValue    int // Value chosen for X costs
}

// PaymentResult represents the result of a payment attempt.
type PaymentResult struct {
	Success   bool
	Plan      *PaymentPlan
	Remaining *ManaCost // Remaining cost if payment failed
	Reason    string
}

// CalculatePayment calculates a payment plan for a mana cost.
// This is a simplified version - full implementation would allow player choice.
func CalculatePayment(cost *ManaCost, pool *ManaPool, xValue int) *PaymentResult {
	if cost == nil {
		return &PaymentResult{
			Success: true,
			Plan:    &PaymentPlan{},
		}
	}

	plan := &PaymentPlan{
		XValue: xValue,
	}

	// Create a copy of the pool to simulate payment
	testPool := pool.Copy()

	// Pay colored mana first (exact match required)
	plan.White = cost.White
	if !testPool.Spend(ManaWhite, cost.White) {
		return &PaymentResult{
			Success: false,
			Reason:  fmt.Sprintf("insufficient white mana (need %d)", cost.White),
		}
	}

	plan.Blue = cost.Blue
	if !testPool.Spend(ManaBlue, cost.Blue) {
		return &PaymentResult{
			Success: false,
			Reason:  fmt.Sprintf("insufficient blue mana (need %d)", cost.Blue),
		}
	}

	plan.Black = cost.Black
	if !testPool.Spend(ManaBlack, cost.Black) {
		return &PaymentResult{
			Success: false,
			Reason:  fmt.Sprintf("insufficient black mana (need %d)", cost.Black),
		}
	}

	plan.Red = cost.Red
	if !testPool.Spend(ManaRed, cost.Red) {
		return &PaymentResult{
			Success: false,
			Reason:  fmt.Sprintf("insufficient red mana (need %d)", cost.Red),
		}
	}

	plan.Green = cost.Green
	if !testPool.Spend(ManaGreen, cost.Green) {
		return &PaymentResult{
			Success: false,
			Reason:  fmt.Sprintf("insufficient green mana (need %d)", cost.Green),
		}
	}

	plan.Colorless = cost.Colorless
	if !testPool.Spend(ManaColorless, cost.Colorless) {
		return &PaymentResult{
			Success: false,
			Reason:  fmt.Sprintf("insufficient colorless mana (need %d)", cost.Colorless),
		}
	}

	// Pay hybrid costs (simplified - try first option)
	for _, hybrid := range cost.Hybrid {
		paid := false
		for _, option := range hybrid.Options {
			if len(option) > 0 {
				if testPool.Spend(option[0], 1) {
					// Track which type was used (simplified)
					switch option[0] {
					case ManaWhite:
						plan.White++
					case ManaBlue:
						plan.Blue++
					case ManaBlack:
						plan.Black++
					case ManaRed:
						plan.Red++
					case ManaGreen:
						plan.Green++
					case ManaColorless:
						plan.Colorless++
					case ManaGeneric:
						plan.Generic++
					}
					paid = true
					break
				}
			}
		}
		if !paid {
			return &PaymentResult{
				Success: false,
				Reason:  "cannot pay hybrid mana cost",
			}
		}
	}

	// Pay generic + X costs (can use any remaining mana)
	totalGeneric := cost.Generic
	if cost.X {
		totalGeneric += xValue
	}

	plan.Generic = totalGeneric
	remainingMana := testPool.GetTotalMana()
	if remainingMana < totalGeneric {
		return &PaymentResult{
			Success: false,
			Reason:  fmt.Sprintf("insufficient mana for generic cost (need %d, have %d)", totalGeneric, remainingMana),
		}
	}

	// Spend generic mana (prefer colorless, then any color)
	genericRemaining := totalGeneric
	if genericRemaining > 0 {
		// Spend colorless first
		if testPool.GetRegular(ManaColorless) > 0 {
			spend := genericRemaining
			if spend > testPool.GetRegular(ManaColorless) {
				spend = testPool.GetRegular(ManaColorless)
			}
			testPool.Spend(ManaColorless, spend)
			genericRemaining -= spend
		}

		// Then spend any colored mana
		types := []ManaType{ManaWhite, ManaBlue, ManaBlack, ManaRed, ManaGreen}
		for _, mt := range types {
			if genericRemaining <= 0 {
				break
			}
			available := testPool.GetRegular(mt)
			if available > 0 {
				spend := genericRemaining
				if spend > available {
					spend = available
				}
				testPool.Spend(mt, spend)
				genericRemaining -= spend
			}
		}

		// Finally use floating mana
		for _, mt := range types {
			if genericRemaining <= 0 {
				break
			}
			available := testPool.GetFloating(mt)
			if available > 0 {
				spend := genericRemaining
				if spend > available {
					spend = available
				}
				testPool.Spend(mt, spend)
				genericRemaining -= spend
			}
		}

		if genericRemaining > 0 {
			return &PaymentResult{
				Success: false,
				Reason:  fmt.Sprintf("insufficient mana for generic cost (need %d more)", genericRemaining),
			}
		}
	}

	return &PaymentResult{
		Success: true,
		Plan:    plan,
	}
}

// ExecutePayment executes a payment plan against a mana pool.
func ExecutePayment(plan *PaymentPlan, pool *ManaPool) bool {
	if plan == nil {
		return true
	}

	// Pay colored mana
	if !pool.Spend(ManaWhite, plan.White) {
		return false
	}
	if !pool.Spend(ManaBlue, plan.Blue) {
		return false
	}
	if !pool.Spend(ManaBlack, plan.Black) {
		return false
	}
	if !pool.Spend(ManaRed, plan.Red) {
		return false
	}
	if !pool.Spend(ManaGreen, plan.Green) {
		return false
	}
	if !pool.Spend(ManaColorless, plan.Colorless) {
		return false
	}

	// Pay generic mana (can use any type)
	genericRemaining := plan.Generic
	if genericRemaining > 0 {
		// Prefer colorless
		if pool.GetRegular(ManaColorless) > 0 {
			spend := genericRemaining
			if spend > pool.GetRegular(ManaColorless) {
				spend = pool.GetRegular(ManaColorless)
			}
			pool.Spend(ManaColorless, spend)
			genericRemaining -= spend
		}

		// Then any colored mana
		types := []ManaType{ManaWhite, ManaBlue, ManaBlack, ManaRed, ManaGreen}
		for _, mt := range types {
			if genericRemaining <= 0 {
				break
			}
			available := pool.GetRegular(mt)
			if available > 0 {
				spend := genericRemaining
				if spend > available {
					spend = available
				}
				pool.Spend(mt, spend)
				genericRemaining -= spend
			}
		}

		// Finally floating mana
		for _, mt := range types {
			if genericRemaining <= 0 {
				break
			}
			available := pool.GetFloating(mt)
			if available > 0 {
				spend := genericRemaining
				if spend > available {
					spend = available
				}
				pool.Spend(mt, spend)
				genericRemaining -= spend
			}
		}

		if genericRemaining > 0 {
			return false
		}
	}

	return true
}
