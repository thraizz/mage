package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sentry Of The Underworld", NewSentryOfTheUnderworld)
}

// NewSentryOfTheUnderworld creates a Sentry Of The Underworld
// {3}{W}{B} - CREATURE
// Flying, Vigilance
func NewSentryOfTheUnderworld(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sentry Of The Underworld")
	card.ManaCost = "{3}{W}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GRIFFIN", "SKELETON"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	// TODO: Implement activated ability with unmapped effects
	//   - RegenerateSourceEffect()
	// card.AddAbility(ability2)
	return card, nil
}
