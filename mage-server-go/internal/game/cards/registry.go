package cards

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
)

// CardBuilder is a function that creates a new card instance
type CardBuilder func(ownerID uuid.UUID, cardInfo *CardInfo) (*game.Card, error)

// Registry is the global card implementation registry
// Cards self-register via init() functions in their package
var Registry = &cardRegistry{
	builders: make(map[string]CardBuilder),
}

type cardRegistry struct {
	builders map[string]CardBuilder
	mu       sync.RWMutex
}

// Register registers a card builder function
// This is typically called from init() functions in card implementation packages
func Register(cardName string, builder CardBuilder) {
	Registry.register(cardName, builder)
}

func (r *cardRegistry) register(cardName string, builder CardBuilder) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.builders[cardName]; exists {
		// Allow re-registration during tests
		fmt.Printf("Warning: Card %s is being re-registered\n", cardName)
	}

	r.builders[cardName] = builder
}

// Get retrieves a card builder by name
func (r *cardRegistry) Get(cardName string) (CardBuilder, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	builder, ok := r.builders[cardName]
	return builder, ok
}

// GetByClassName retrieves a card builder by Java class name
// e.g., "mage.cards.l.LightningBolt" -> "Lightning Bolt"
func (r *cardRegistry) GetByClassName(className string) (CardBuilder, bool) {
	// For now, just try to extract the simple name
	// TODO: Build a proper class name -> card name mapping
	return r.Get(className)
}

// IsImplemented returns true if a card implementation exists
func (r *cardRegistry) IsImplemented(cardName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.builders[cardName]
	return ok
}

// Count returns the number of registered cards
func (r *cardRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.builders)
}

// ListAll returns all registered card names
func (r *cardRegistry) ListAll() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.builders))
	for name := range r.builders {
		names = append(names, name)
	}
	return names
}

// Clear removes all registered cards (useful for testing)
func (r *cardRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.builders = make(map[string]CardBuilder)
}

// BuildCard creates a Card from the registry by card name
// This function is designed to be passed to MageEngine.SetCardBuilder
// to avoid import cycles (game cannot import cards, but cards imports game)
// Returns (nil, nil) if card is not implemented in registry (not an error)
func BuildCard(cardName string, ownerID uuid.UUID) (*game.Card, error) {
	builder, ok := Registry.Get(cardName)
	if !ok {
		// Card not in registry - this is not an error, just not implemented
		return nil, nil
	}
	return builder(ownerID, nil)
}
