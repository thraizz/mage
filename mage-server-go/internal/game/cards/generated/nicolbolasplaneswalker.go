package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nicol Bolas Planeswalker", NewNicolBolasPlaneswalker)
}

// NewNicolBolasPlaneswalker creates a Nicol Bolas Planeswalker
// {4}{U}{B}{B}{R} - PLANESWALKER
func NewNicolBolasPlaneswalker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nicol Bolas Planeswalker")
	card.ManaCost = "{4}{U}{B}{B}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"BOLAS"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		AddEffect(abilities.NewGainControlTargetEffect(abilities.DurationCustom)).
		AddEffect(abilities.NewDamageEffect(7)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}