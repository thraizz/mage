package cards

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/repository"
	"go.uber.org/zap"
)

// Factory creates card instances from card names and database metadata
type Factory interface {
	// CreateCard creates a new card instance by name
	CreateCard(ctx context.Context, name string, ownerID uuid.UUID) (*game.LegacyCard, error)

	// IsImplemented returns true if a card has an implementation
	IsImplemented(name string) bool

	// Stats returns implementation statistics
	Stats() FactoryStats
}

// FactoryStats contains statistics about card implementations
type FactoryStats struct {
	TotalCardsInDB     int
	ImplementedCards   int
	UnimplementedCards int
	PercentImplemented float64
}

// factory implements the Factory interface
type factory struct {
	cardRepo *repository.CardRepository
	registry *cardRegistry
	logger   *zap.Logger
}

// NewFactory creates a new card factory
func NewFactory(cardRepo *repository.CardRepository, logger *zap.Logger) Factory {
	return &factory{
		cardRepo: cardRepo,
		registry: Registry,
		logger:   logger,
	}
}

// CreateCard creates a new card instance by name
func (f *factory) CreateCard(ctx context.Context, name string, ownerID uuid.UUID) (*game.LegacyCard, error) {
	// 1. Check if card is implemented
	builder, ok := f.registry.Get(name)
	if !ok {
		return nil, fmt.Errorf("card not implemented: %s", name)
	}

	// 2. Get card metadata from database
	cards, err := f.cardRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get card metadata: %w", err)
	}

	if len(cards) == 0 {
		return nil, fmt.Errorf("card not found in database: %s", name)
	}

	// Use first printing (could be enhanced to select specific set)
	cardData := cards[0]

	// 3. Convert repository.Card to cards.CardInfo
	cardInfo := &CardInfo{
		ID:         cardData.ID,
		CardNumber: cardData.CardNumber,
		SetCode:    cardData.SetCode,
		Name:       cardData.Name,
		CardType:   cardData.CardType,
		ManaCost:   cardData.ManaCost,
		Power:      cardData.Power,
		Toughness:  cardData.Toughness,
		RulesText:  cardData.RulesText,
		FlavorText: nullStringToString(cardData.FlavorText),
		Rarity:     cardData.Rarity,
		// TODO: Parse types, subtypes, supertypes, colors from CardType string
	}

	// 4. Build the card using the registered builder
	card, err := builder(ownerID, cardInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to build card: %w", err)
	}

	f.logger.Debug("created card",
		zap.String("name", name),
		zap.String("owner", ownerID.String()),
	)

	return card, nil
}

// IsImplemented returns true if a card has an implementation
func (f *factory) IsImplemented(name string) bool {
	return f.registry.IsImplemented(name)
}

// Stats returns implementation statistics
func (f *factory) Stats() FactoryStats {
	// TODO: Query database for total card count
	// For now, return basic stats
	implementedCount := f.registry.Count()

	return FactoryStats{
		TotalCardsInDB:     30459, // From our export
		ImplementedCards:   implementedCount,
		UnimplementedCards: 30459 - implementedCount,
		PercentImplemented: float64(implementedCount) / 30459.0 * 100.0,
	}
}

// nullStringToString converts sql.NullString to string
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
