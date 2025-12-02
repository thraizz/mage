package manual

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Martyr of Sands", NewMartyrOfSands)
}

// NewMartyrOfSands creates a Martyr of Sands
// {W} - Creature — Human Cleric - 1/1
// {1}, Reveal X white cards from your hand, Sacrifice Martyr of Sands: You gain three times X life.
func NewMartyrOfSands(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Martyr of Sands")
	card.ManaCost = "{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "CLERIC"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "CSP"
	card.Rarity = "common"
	card.RulesText = "{1}, Reveal X white cards from your hand, Sacrifice {this}: You gain three times X life."

	// Build the activated ability
	// Cost: {1}, Sacrifice this creature
	// Effect: Gain life based on white cards revealed (simplified: we'll count all white cards in hand)
	ability := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddSacrificeSourceCost().
		AddEffect(&MartyrOfSandsEffect{sourceID: card.ID}).
		Build()

	// Note: The rules text on the card already describes the ability
	// The ability text is auto-generated from costs/effects for display

	card.AddAbility(ability)

	return card, nil
}

// MartyrOfSandsEffect implements the life gain effect for Martyr of Sands
// It counts white cards in the controller's hand and gains 3x that much life
type MartyrOfSandsEffect struct {
	sourceID uuid.UUID
}

func (e *MartyrOfSandsEffect) Apply(ctx context.Context, gameCtx abilities.GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Get the controller of the source
	controllerID, err := gameCtx.GetControllerID(source)
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	// Get the hand of the controller
	hand, err := gameCtx.GetPlayerHand(controllerID)
	if err != nil {
		return fmt.Errorf("failed to get hand: %w", err)
	}

	// Count white cards in hand
	// TODO: In a full implementation, we'd prompt the player to choose which cards to reveal
	// For now, we reveal all white cards in hand
	whiteCardCount := 0
	for _, cardInterface := range hand {
		// The hand contains *internalCard which has a Color field
		// We can access it through reflection or type assertion
		// First, try to get the Color field directly via reflection
		if card, ok := cardInterface.(interface{ GetColor() string }); ok {
			color := card.GetColor()
			if isWhiteColor(color) {
				whiteCardCount++
				continue
			}
		}

		// Fallback: try to get the card ID and look up colors
		if card, ok := cardInterface.(interface{ GetID() string }); ok {
			cardUUID, err := uuid.Parse(card.GetID())
			if err != nil {
				continue
			}

			colors, err := gameCtx.GetCardColors(cardUUID)
			if err != nil {
				continue
			}
			for _, color := range colors {
				if color == "W" {
					whiteCardCount++
					break
				}
			}
		}
	}

	// Gain 3 times X life
	lifeGain := whiteCardCount * 3

	if lifeGain > 0 {
		if err := gameCtx.GainLife(controllerID, lifeGain); err != nil {
			return fmt.Errorf("failed to gain life: %w", err)
		}
	}

	return nil
}

// isWhiteColor checks if a color string indicates white
func isWhiteColor(color string) bool {
	lowerColor := strings.ToLower(color)
	return strings.Contains(lowerColor, "white") || lowerColor == "w"
}

func (e *MartyrOfSandsEffect) GetDescription() string {
	return "You gain three times X life, where X is the number of white cards revealed"
}
