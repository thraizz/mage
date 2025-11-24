package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ashiok Dream Render", NewAshiokDreamRender)
}

// NewAshiokDreamRender creates a Ashiok Dream Render
// {1}{U/B}{U/B} - PLANESWALKER
func NewAshiokDreamRender(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ashiok Dream Render")
	card.ManaCost = "{1}{U/B}{U/B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"ASHIOK"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewMillCardsTargetEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}