package mana

import (
	"testing"
)

func TestManaPool_Add(t *testing.T) {
	pool := NewManaPool()

	pool.Add(ManaWhite, 2)
	if pool.GetRegular(ManaWhite) != 2 {
		t.Errorf("Expected 2 white mana, got %d", pool.GetRegular(ManaWhite))
	}

	pool.Add(ManaBlue, 1)
	if pool.GetRegular(ManaBlue) != 1 {
		t.Errorf("Expected 1 blue mana, got %d", pool.GetRegular(ManaBlue))
	}
}

func TestManaPool_Spend(t *testing.T) {
	pool := NewManaPool()
	pool.Add(ManaWhite, 3)
	pool.Add(ManaBlue, 2)

	if !pool.Spend(ManaWhite, 2) {
		t.Error("Expected to spend 2 white mana")
	}
	if pool.GetRegular(ManaWhite) != 1 {
		t.Errorf("Expected 1 white mana remaining, got %d", pool.GetRegular(ManaWhite))
	}

	if !pool.Spend(ManaBlue, 1) {
		t.Error("Expected to spend 1 blue mana")
	}
	if pool.GetRegular(ManaBlue) != 1 {
		t.Errorf("Expected 1 blue mana remaining, got %d", pool.GetRegular(ManaBlue))
	}

	// Try to spend more than available
	if pool.Spend(ManaWhite, 5) {
		t.Error("Expected to fail spending 5 white mana when only 1 available")
	}
}

func TestManaPool_FloatingMana(t *testing.T) {
	pool := NewManaPool()

	pool.AddFloating(ManaRed, 2)
	if pool.GetFloating(ManaRed) != 2 {
		t.Errorf("Expected 2 floating red mana, got %d", pool.GetFloating(ManaRed))
	}

	pool.Add(ManaRed, 1)
	if pool.GetTotal(ManaRed) != 3 {
		t.Errorf("Expected total 3 red mana, got %d", pool.GetTotal(ManaRed))
	}

	// Spend prefers regular over floating
	pool.Spend(ManaRed, 1)
	if pool.GetRegular(ManaRed) != 0 {
		t.Errorf("Expected 0 regular red mana after spending, got %d", pool.GetRegular(ManaRed))
	}
	if pool.GetFloating(ManaRed) != 2 {
		t.Errorf("Expected 2 floating red mana remaining, got %d", pool.GetFloating(ManaRed))
	}
}

func TestManaPool_Empty(t *testing.T) {
	pool := NewManaPool()
	pool.Add(ManaWhite, 2)
	pool.AddFloating(ManaBlue, 1)

	pool.Empty()
	if pool.GetRegular(ManaWhite) != 0 {
		t.Error("Expected regular mana to be empty")
	}
	if pool.GetFloating(ManaBlue) != 1 {
		t.Error("Expected floating mana to persist")
	}

	pool.EmptyFloating()
	if pool.GetFloating(ManaBlue) != 0 {
		t.Error("Expected floating mana to be empty")
	}
}
