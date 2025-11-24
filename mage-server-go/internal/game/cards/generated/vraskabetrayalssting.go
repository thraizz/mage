package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vraska Betrayals Sting", NewVraskaBetrayalsSting)
}

// NewVraskaBetrayalsSting creates a Vraska Betrayals Sting
// {4}{B}{B/P} - PLANESWALKER
func NewVraskaBetrayalsSting(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vraska Betrayals Sting")
	card.ManaCost = "{4}{B}{B/P}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"VRASKA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "6"
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
