package mana

import (
	"testing"
)

func TestParseCost(t *testing.T) {
	tests := []struct {
		input    string
		expected *ManaCost
		err      bool
	}{
		{"", &ManaCost{}, false},
		{"{1}", &ManaCost{Generic: 1}, false},
		{"{G}", &ManaCost{Green: 1}, false},
		{"{1}{G}", &ManaCost{Generic: 1, Green: 1}, false},
		{"{2}{R}{R}", &ManaCost{Generic: 2, Red: 2}, false},
		{"{X}{R}", &ManaCost{X: true, Red: 1}, false},
		{"{W}{U}{B}{R}{G}", &ManaCost{White: 1, Blue: 1, Black: 1, Red: 1, Green: 1}, false},
		{"{C}", &ManaCost{Colorless: 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseCost(tt.input)
			if tt.err {
				if err == nil {
					t.Errorf("Expected error for %s, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.input, err)
				return
			}
			if result.Generic != tt.expected.Generic {
				t.Errorf("Generic: expected %d, got %d", tt.expected.Generic, result.Generic)
			}
			if result.White != tt.expected.White {
				t.Errorf("White: expected %d, got %d", tt.expected.White, result.White)
			}
			if result.Blue != tt.expected.Blue {
				t.Errorf("Blue: expected %d, got %d", tt.expected.Blue, result.Blue)
			}
			if result.Black != tt.expected.Black {
				t.Errorf("Black: expected %d, got %d", tt.expected.Black, result.Black)
			}
			if result.Red != tt.expected.Red {
				t.Errorf("Red: expected %d, got %d", tt.expected.Red, result.Red)
			}
			if result.Green != tt.expected.Green {
				t.Errorf("Green: expected %d, got %d", tt.expected.Green, result.Green)
			}
			if result.Colorless != tt.expected.Colorless {
				t.Errorf("Colorless: expected %d, got %d", tt.expected.Colorless, result.Colorless)
			}
			if result.X != tt.expected.X {
				t.Errorf("X: expected %v, got %v", tt.expected.X, result.X)
			}
		})
	}
}

func TestManaCost_CanPay(t *testing.T) {
	pool := NewManaPool()
	pool.Add(ManaWhite, 1)
	pool.Add(ManaBlue, 2)
	pool.Add(ManaGreen, 1)

	tests := []struct {
		cost     string
		xValue   int
		canPay   bool
	}{
		{"{G}", 0, true},
		{"{U}", 0, true},
		{"{W}", 0, true},
		{"{R}", 0, false},
		{"{1}{G}", 0, true}, // 1 generic + 1 green, pool has 1 green + 2 blue + 1 white = 4 total, need 2
		{"{3}{G}", 0, true}, // 3 generic + 1 green, pool has 1 green + 2 blue + 1 white = 4 total, need 4 (1 green + 3 generic from blue/white)
		{"{X}{G}", 0, true},  // X=0 means 0 generic + 1 green
		{"{X}{G}", 1, true},  // X=1 means 1 generic + 1 green, pool has 4 total, need 2
		{"{X}{G}", 2, true}, // X=2 means 2 generic + 1 green, pool has 4 total (1 green + 2 blue + 1 white), need 3 (1 green + 2 generic from blue/white)
	}

	for _, tt := range tests {
		t.Run(tt.cost, func(t *testing.T) {
			cost, err := ParseCost(tt.cost)
			if err != nil {
				t.Fatalf("Failed to parse cost: %v", err)
			}
			result := cost.CanPay(pool, tt.xValue)
			if result != tt.canPay {
				t.Errorf("CanPay(%s, x=%d): expected %v, got %v", tt.cost, tt.xValue, tt.canPay, result)
			}
		})
	}
}

func TestManaCost_ApplyReduction(t *testing.T) {
	cost, _ := ParseCost("{3}{G}{G}")

	// Reduce generic by 1
	reduced := cost.ApplyReduction(1, nil)
	if reduced.Generic != 2 {
		t.Errorf("Expected generic 2, got %d", reduced.Generic)
	}
	if reduced.Green != 2 {
		t.Errorf("Expected green 2, got %d", reduced.Green)
	}

	// Reduce green by 1
	coloredReduction := map[ManaType]int{ManaGreen: 1}
	reduced2 := cost.ApplyReduction(0, coloredReduction)
	if reduced2.Green != 1 {
		t.Errorf("Expected green 1, got %d", reduced2.Green)
	}

	// Can't reduce below 0
	reduced3 := cost.ApplyReduction(5, nil)
	if reduced3.Generic != 0 {
		t.Errorf("Expected generic 0, got %d", reduced3.Generic)
	}
}
