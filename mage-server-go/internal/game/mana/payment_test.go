package mana

import (
	"testing"
)

func TestCalculatePayment(t *testing.T) {
	pool := NewManaPool()
	pool.Add(ManaWhite, 1)
	pool.Add(ManaBlue, 2)
	pool.Add(ManaGreen, 1)

	cost, _ := ParseCost("{1}{G}")
	result := CalculatePayment(cost, pool, 0)

	if !result.Success {
		t.Errorf("Expected successful payment, got: %s", result.Reason)
	}
	if result.Plan == nil {
		t.Fatal("Expected payment plan")
	}
	if result.Plan.Green != 1 {
		t.Errorf("Expected 1 green in plan, got %d", result.Plan.Green)
	}
	if result.Plan.Generic != 1 {
		t.Errorf("Expected 1 generic in plan, got %d", result.Plan.Generic)
	}
}

func TestCalculatePayment_InsufficientMana(t *testing.T) {
	pool := NewManaPool()
	pool.Add(ManaGreen, 1)

	cost, _ := ParseCost("{3}{G}")
	result := CalculatePayment(cost, pool, 0)

	if result.Success {
		t.Error("Expected payment to fail")
	}
	if result.Reason == "" {
		t.Error("Expected failure reason")
	}
}

func TestExecutePayment(t *testing.T) {
	pool := NewManaPool()
	pool.Add(ManaWhite, 2)
	pool.Add(ManaBlue, 1)

	plan := &PaymentPlan{
		White: 1,
		Blue:  1,
	}

	if !ExecutePayment(plan, pool) {
		t.Error("Expected successful payment execution")
	}
	if pool.GetRegular(ManaWhite) != 1 {
		t.Errorf("Expected 1 white mana remaining, got %d", pool.GetRegular(ManaWhite))
	}
	if pool.GetRegular(ManaBlue) != 0 {
		t.Errorf("Expected 0 blue mana remaining, got %d", pool.GetRegular(ManaBlue))
	}
}

func TestExecutePayment_GenericMana(t *testing.T) {
	pool := NewManaPool()
	pool.Add(ManaWhite, 3)

	plan := &PaymentPlan{
		Generic: 2,
	}

	if !ExecutePayment(plan, pool) {
		t.Error("Expected successful payment execution")
	}
	// Generic should be paid with white mana
	if pool.GetRegular(ManaWhite) != 1 {
		t.Errorf("Expected 1 white mana remaining (2 spent for generic), got %d", pool.GetRegular(ManaWhite))
	}
}
