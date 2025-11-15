package mana

import (
	"sync"
)

// ManaType represents a type of mana.
type ManaType string

const (
	ManaWhite     ManaType = "WHITE"
	ManaBlue      ManaType = "BLUE"
	ManaBlack     ManaType = "BLACK"
	ManaRed       ManaType = "RED"
	ManaGreen     ManaType = "GREEN"
	ManaColorless ManaType = "COLORLESS"
	ManaGeneric   ManaType = "GENERIC" // Generic mana can be paid with any type
)

// ManaPool represents a player's mana pool (both regular and floating).
type ManaPool struct {
	mu sync.RWMutex

	// Regular mana pool (empties at end of step)
	White     int
	Blue      int
	Black     int
	Red       int
	Green     int
	Colorless int

	// Floating mana (persists until end of step/phase)
	FloatingWhite     int
	FloatingBlue      int
	FloatingBlack     int
	FloatingRed       int
	FloatingGreen     int
	FloatingColorless int
}

// NewManaPool creates a new empty mana pool.
func NewManaPool() *ManaPool {
	return &ManaPool{}
}

// Add adds mana to the regular pool.
func (mp *ManaPool) Add(manaType ManaType, amount int) {
	if amount <= 0 {
		return
	}
	mp.mu.Lock()
	defer mp.mu.Unlock()

	switch manaType {
	case ManaWhite:
		mp.White += amount
	case ManaBlue:
		mp.Blue += amount
	case ManaBlack:
		mp.Black += amount
	case ManaRed:
		mp.Red += amount
	case ManaGreen:
		mp.Green += amount
	case ManaColorless:
		mp.Colorless += amount
	}
}

// AddFloating adds mana to the floating pool.
func (mp *ManaPool) AddFloating(manaType ManaType, amount int) {
	if amount <= 0 {
		return
	}
	mp.mu.Lock()
	defer mp.mu.Unlock()

	switch manaType {
	case ManaWhite:
		mp.FloatingWhite += amount
	case ManaBlue:
		mp.FloatingBlue += amount
	case ManaBlack:
		mp.FloatingBlack += amount
	case ManaRed:
		mp.FloatingRed += amount
	case ManaGreen:
		mp.FloatingGreen += amount
	case ManaColorless:
		mp.FloatingColorless += amount
	}
}

// GetTotal returns the total amount of a specific mana type (regular + floating).
func (mp *ManaPool) GetTotal(manaType ManaType) int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	switch manaType {
	case ManaWhite:
		return mp.White + mp.FloatingWhite
	case ManaBlue:
		return mp.Blue + mp.FloatingBlue
	case ManaBlack:
		return mp.Black + mp.FloatingBlack
	case ManaRed:
		return mp.Red + mp.FloatingRed
	case ManaGreen:
		return mp.Green + mp.FloatingGreen
	case ManaColorless:
		return mp.Colorless + mp.FloatingColorless
	default:
		return 0
	}
}

// GetRegular returns the regular (non-floating) amount of a mana type.
func (mp *ManaPool) GetRegular(manaType ManaType) int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	switch manaType {
	case ManaWhite:
		return mp.White
	case ManaBlue:
		return mp.Blue
	case ManaBlack:
		return mp.Black
	case ManaRed:
		return mp.Red
	case ManaGreen:
		return mp.Green
	case ManaColorless:
		return mp.Colorless
	default:
		return 0
	}
}

// GetFloating returns the floating amount of a mana type.
func (mp *ManaPool) GetFloating(manaType ManaType) int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	switch manaType {
	case ManaWhite:
		return mp.FloatingWhite
	case ManaBlue:
		return mp.FloatingBlue
	case ManaBlack:
		return mp.FloatingBlack
	case ManaRed:
		return mp.FloatingRed
	case ManaGreen:
		return mp.FloatingGreen
	case ManaColorless:
		return mp.FloatingColorless
	default:
		return 0
	}
}

// Spend attempts to spend mana from the pool.
// Returns true if successful, false if insufficient mana.
// Prefers regular mana over floating mana.
func (mp *ManaPool) Spend(manaType ManaType, amount int) bool {
	if amount <= 0 {
		return true
	}
	mp.mu.Lock()
	defer mp.mu.Unlock()

	var available int
	switch manaType {
	case ManaWhite:
		available = mp.White + mp.FloatingWhite
	case ManaBlue:
		available = mp.Blue + mp.FloatingBlue
	case ManaBlack:
		available = mp.Black + mp.FloatingBlack
	case ManaRed:
		available = mp.Red + mp.FloatingRed
	case ManaGreen:
		available = mp.Green + mp.FloatingGreen
	case ManaColorless:
		available = mp.Colorless + mp.FloatingColorless
	default:
		return false
	}

	if available < amount {
		return false
	}

	// Spend from regular pool first, then floating
	switch manaType {
	case ManaWhite:
		spendFromRegular := amount
		if spendFromRegular > mp.White {
			spendFromRegular = mp.White
		}
		mp.White -= spendFromRegular
		mp.FloatingWhite -= (amount - spendFromRegular)
	case ManaBlue:
		spendFromRegular := amount
		if spendFromRegular > mp.Blue {
			spendFromRegular = mp.Blue
		}
		mp.Blue -= spendFromRegular
		mp.FloatingBlue -= (amount - spendFromRegular)
	case ManaBlack:
		spendFromRegular := amount
		if spendFromRegular > mp.Black {
			spendFromRegular = mp.Black
		}
		mp.Black -= spendFromRegular
		mp.FloatingBlack -= (amount - spendFromRegular)
	case ManaRed:
		spendFromRegular := amount
		if spendFromRegular > mp.Red {
			spendFromRegular = mp.Red
		}
		mp.Red -= spendFromRegular
		mp.FloatingRed -= (amount - spendFromRegular)
	case ManaGreen:
		spendFromRegular := amount
		if spendFromRegular > mp.Green {
			spendFromRegular = mp.Green
		}
		mp.Green -= spendFromRegular
		mp.FloatingGreen -= (amount - spendFromRegular)
	case ManaColorless:
		spendFromRegular := amount
		if spendFromRegular > mp.Colorless {
			spendFromRegular = mp.Colorless
		}
		mp.Colorless -= spendFromRegular
		mp.FloatingColorless -= (amount - spendFromRegular)
	}

	return true
}

// Empty empties the regular mana pool (floating mana persists).
func (mp *ManaPool) Empty() {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.White = 0
	mp.Blue = 0
	mp.Black = 0
	mp.Red = 0
	mp.Green = 0
	mp.Colorless = 0
}

// EmptyFloating empties the floating mana pool.
func (mp *ManaPool) EmptyFloating() {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.FloatingWhite = 0
	mp.FloatingBlue = 0
	mp.FloatingBlack = 0
	mp.FloatingRed = 0
	mp.FloatingGreen = 0
	mp.FloatingColorless = 0
}

// EmptyAll empties both regular and floating mana pools.
func (mp *ManaPool) EmptyAll() {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.White = 0
	mp.Blue = 0
	mp.Black = 0
	mp.Red = 0
	mp.Green = 0
	mp.Colorless = 0
	mp.FloatingWhite = 0
	mp.FloatingBlue = 0
	mp.FloatingBlack = 0
	mp.FloatingRed = 0
	mp.FloatingGreen = 0
	mp.FloatingColorless = 0
}

// GetTotalMana returns the total mana count across all types.
func (mp *ManaPool) GetTotalMana() int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return mp.White + mp.Blue + mp.Black + mp.Red + mp.Green + mp.Colorless +
		mp.FloatingWhite + mp.FloatingBlue + mp.FloatingBlack + mp.FloatingRed + mp.FloatingGreen + mp.FloatingColorless
}

// Copy creates a deep copy of the mana pool.
func (mp *ManaPool) Copy() *ManaPool {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return &ManaPool{
		White:            mp.White,
		Blue:             mp.Blue,
		Black:            mp.Black,
		Red:              mp.Red,
		Green:            mp.Green,
		Colorless:        mp.Colorless,
		FloatingWhite:    mp.FloatingWhite,
		FloatingBlue:     mp.FloatingBlue,
		FloatingBlack:    mp.FloatingBlack,
		FloatingRed:      mp.FloatingRed,
		FloatingGreen:    mp.FloatingGreen,
		FloatingColorless: mp.FloatingColorless,
	}
}
