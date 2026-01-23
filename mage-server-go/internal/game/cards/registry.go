package cards

import (
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
)

// CardBuilder is a function that creates a new card instance
type CardBuilder func(ownerID uuid.UUID, cardInfo *CardInfo) (*game.LegacyCard, error)

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

// RegisterSplitCard registers a split/DFC card under both its combined name and individual face names
// Example: RegisterSplitCard("Fire // Ice", []string{"Fire", "Ice"}, builder)
func RegisterSplitCard(fullName string, faceNames []string, builder CardBuilder) {
	Registry.registerSplitCard(fullName, faceNames, builder)
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

func (r *cardRegistry) registerSplitCard(fullName string, faceNames []string, builder CardBuilder) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Register under full name
	r.builders[fullName] = builder

	// Register under each face name
	for _, faceName := range faceNames {
		if faceName != "" {
			r.builders[faceName] = builder
		}
	}
}

// Get retrieves a card builder by name
// Handles split/DFC card names by trying multiple variations:
// 1. Exact match (e.g., "Fire // Ice")
// 2. First face only (e.g., "Fire" for "Fire // Ice")
// 3. Second face only (e.g., "Ice" for "Fire // Ice")
func (r *cardRegistry) Get(cardName string) (CardBuilder, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Try exact match first
	if builder, ok := r.builders[cardName]; ok {
		return builder, true
	}

	// Try normalizing split/DFC card names
	// Handle "Fire // Ice" or "Fire / Ice" formats
	if strings.Contains(cardName, "//") {
		parts := strings.Split(cardName, "//")
		if len(parts) >= 2 {
			// Try first face
			firstFace := strings.TrimSpace(parts[0])
			if builder, ok := r.builders[firstFace]; ok {
				return builder, true
			}
			// Try second face
			secondFace := strings.TrimSpace(parts[1])
			if builder, ok := r.builders[secondFace]; ok {
				return builder, true
			}
		}
	}

	// Try single slash format
	if strings.Contains(cardName, " / ") {
		normalized := strings.ReplaceAll(cardName, " / ", " // ")
		if builder, ok := r.builders[normalized]; ok {
			return builder, true
		}
		// Also try individual faces
		parts := strings.Split(cardName, " / ")
		if len(parts) >= 2 {
			firstFace := strings.TrimSpace(parts[0])
			if builder, ok := r.builders[firstFace]; ok {
				return builder, true
			}
		}
	}

	return nil, false
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
func BuildCard(cardName string, ownerID uuid.UUID) (*game.LegacyCard, error) {
	builder, ok := Registry.Get(cardName)
	if !ok {
		// Card not in registry - this is not an error, just not implemented
		return nil, nil
	}
	return builder(ownerID, nil)
}
