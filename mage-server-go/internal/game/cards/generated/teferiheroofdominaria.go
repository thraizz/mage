package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Teferi Hero Of Dominaria", NewTeferiHeroOfDominaria)
}

// NewTeferiHeroOfDominaria creates a Teferi Hero Of Dominaria
// {3}{W}{U} - PLANESWALKER
func NewTeferiHeroOfDominaria(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Teferi Hero Of Dominaria")
	card.ManaCost = "{3}{W}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"TEFERI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
