package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Makindi Sliderunner", NewMakindiSliderunner)
}

// NewMakindiSliderunner creates a Makindi Sliderunner
// {1}{R} - CREATURE
// Trample
func NewMakindiSliderunner(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Makindi Sliderunner")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAST"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
