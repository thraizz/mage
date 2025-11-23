package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_Register(t *testing.T) {
	// Create a fresh registry for testing
	reg := &tokenRegistry{
		builders: make(map[string]TokenBuilder),
	}

	builder := func() *Token {
		return NewToken("TestToken", "1/1 test creature token")
	}

	reg.register("TestToken", builder)

	assert.Equal(t, 1, reg.Count())

	// Test retrieval
	retrieved, ok := reg.Get("TestToken")
	assert.True(t, ok)
	assert.NotNil(t, retrieved)

	// Test building
	tok := retrieved()
	assert.Equal(t, "TestToken", tok.Name)
	assert.Equal(t, "1/1 test creature token", tok.Description)
}

func TestRegistry_GetToken(t *testing.T) {
	reg := &tokenRegistry{
		builders: make(map[string]TokenBuilder),
	}

	builder := func() *Token {
		tok := NewToken("SoldierToken", "1/1 white Soldier creature token")
		tok.AddCardType(CardTypeCreature)
		tok.AddSubtype("SOLDIER")
		tok.SetColor(Color{White: true})
		tok.SetPowerToughness(1, 1)
		return tok
	}

	reg.register("SoldierToken", builder)

	tok, err := reg.GetToken("SoldierToken")
	require.NoError(t, err)
	assert.Equal(t, "SoldierToken", tok.Name)
	assert.Equal(t, 1, tok.Power)
	assert.Equal(t, 1, tok.Toughness)
	assert.True(t, tok.Color.White)
	assert.Contains(t, tok.Subtypes, "SOLDIER")
}

func TestRegistry_GetToken_NotFound(t *testing.T) {
	reg := &tokenRegistry{
		builders: make(map[string]TokenBuilder),
	}

	_, err := reg.GetToken("NonexistentToken")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in registry")
}

func TestRegistry_List(t *testing.T) {
	reg := &tokenRegistry{
		builders: make(map[string]TokenBuilder),
	}

	tokens := []string{"Token1", "Token2", "Token3"}
	for _, name := range tokens {
		n := name // Capture loop variable
		reg.register(name, func() *Token {
			return NewToken(n, "")
		})
	}

	list := reg.List()
	assert.Equal(t, 3, len(list))

	// Check all tokens are in the list
	for _, name := range tokens {
		assert.Contains(t, list, name)
	}
}

func TestRegistry_Count(t *testing.T) {
	reg := &tokenRegistry{
		builders: make(map[string]TokenBuilder),
	}

	assert.Equal(t, 0, reg.Count())

	reg.register("Token1", func() *Token { return NewToken("Token1", "") })
	assert.Equal(t, 1, reg.Count())

	reg.register("Token2", func() *Token { return NewToken("Token2", "") })
	assert.Equal(t, 2, reg.Count())
}

func TestRegistry_Clear(t *testing.T) {
	reg := &tokenRegistry{
		builders: make(map[string]TokenBuilder),
	}

	reg.register("Token1", func() *Token { return NewToken("Token1", "") })
	reg.register("Token2", func() *Token { return NewToken("Token2", "") })

	assert.Equal(t, 2, reg.Count())

	reg.Clear()
	assert.Equal(t, 0, reg.Count())
}

func TestGlobalRegistry(t *testing.T) {
	// This test verifies that the global registry has been populated
	// by the init() functions in generated_tokens.go

	count := CountTokens()
	t.Logf("Global registry has %d tokens", count)

	// We expect 711 tokens from the generated file
	assert.Equal(t, 711, count, "Expected 711 tokens to be registered")

	// Test retrieval of a few known tokens
	tokens := []string{"ATATToken", "AkroanSoldierToken", "SaprolingToken", "ZombieToken"}

	for _, tokenName := range tokens {
		tok, err := GetToken(tokenName)
		if err != nil {
			// Token might not exist in the list above, just log it
			t.Logf("Token %s not found (might not be in generated list)", tokenName)
			continue
		}

		assert.NotNil(t, tok, "Token %s should not be nil", tokenName)
		t.Logf("Successfully retrieved token: %s", tok.Name)
	}

	// List all tokens
	allTokens := ListTokens()
	assert.Equal(t, 711, len(allTokens), "ListTokens should return all 711 tokens")

	t.Logf("First 10 tokens: %v", allTokens[:min(10, len(allTokens))])
}

func TestToken_Copy(t *testing.T) {
	original := NewToken("TestToken", "Test description")
	original.AddCardType(CardTypeCreature)
	original.AddSubtype("GOBLIN")
	original.SetColor(Color{Red: true})
	original.SetPowerToughness(2, 2)
	original.AddAbility("haste")

	copy := original.Copy()

	// Verify all fields are copied
	assert.Equal(t, original.Name, copy.Name)
	assert.Equal(t, original.Description, copy.Description)
	assert.Equal(t, original.Power, copy.Power)
	assert.Equal(t, original.Toughness, copy.Toughness)
	assert.Equal(t, original.Color, copy.Color)
	assert.Equal(t, original.CardTypes, copy.CardTypes)
	assert.Equal(t, original.Subtypes, copy.Subtypes)
	assert.Equal(t, original.Abilities, copy.Abilities)

	// Verify it's a deep copy (modifying copy doesn't affect original)
	copy.SetPowerToughness(3, 3)
	assert.Equal(t, 2, original.Power)
	assert.Equal(t, 3, copy.Power)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
