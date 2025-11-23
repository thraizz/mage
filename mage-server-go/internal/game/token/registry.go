package token

import (
	"fmt"
	"sync"
)

// TokenBuilder is a function that creates a new token instance.
type TokenBuilder func() *Token

// tokenRegistry holds all registered token builders.
type tokenRegistry struct {
	builders map[string]TokenBuilder
	mu       sync.RWMutex
}

// Registry is the global token registry.
var Registry = &tokenRegistry{
	builders: make(map[string]TokenBuilder),
}

// Register adds a token builder to the registry.
// This is called by generated token init() functions.
func Register(tokenName string, builder TokenBuilder) {
	Registry.register(tokenName, builder)
}

// register is the internal method to add a token builder.
func (r *tokenRegistry) register(tokenName string, builder TokenBuilder) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.builders[tokenName]; exists {
		// Allow re-registration for testing purposes
		fmt.Printf("Warning: Token '%s' is being re-registered\n", tokenName)
	}

	r.builders[tokenName] = builder
}

// Get retrieves a token builder by name.
func (r *tokenRegistry) Get(tokenName string) (TokenBuilder, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	builder, exists := r.builders[tokenName]
	return builder, exists
}

// GetToken creates a new token instance by name.
func (r *tokenRegistry) GetToken(tokenName string) (*Token, error) {
	builder, exists := r.Get(tokenName)
	if !exists {
		return nil, fmt.Errorf("token '%s' not found in registry", tokenName)
	}

	return builder(), nil
}

// List returns all registered token names.
func (r *tokenRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.builders))
	for name := range r.builders {
		names = append(names, name)
	}
	return names
}

// Count returns the number of registered tokens.
func (r *tokenRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.builders)
}

// Clear removes all registered tokens (useful for testing).
func (r *tokenRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.builders = make(map[string]TokenBuilder)
}

// GetToken is a convenience function to create a token by name.
func GetToken(tokenName string) (*Token, error) {
	return Registry.GetToken(tokenName)
}

// ListTokens returns all registered token names.
func ListTokens() []string {
	return Registry.List()
}

// CountTokens returns the number of registered tokens.
func CountTokens() int {
	return Registry.Count()
}
