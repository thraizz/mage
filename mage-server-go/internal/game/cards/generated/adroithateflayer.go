package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Adroit Hateflayer", NewAdroitHateflayer)
}

// NewAdroitHateflayer creates a Adroit Hateflayer
// {U}{B}{R} - CREATURE
// Menace
func NewAdroitHateflayer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Adroit Hateflayer")
	card.ManaCost = "{U}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"NAUTOLAN", "SITH"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordMenace)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewLoseLifeEffect(2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
