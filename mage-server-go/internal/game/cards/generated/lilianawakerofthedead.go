package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Liliana Waker Of The Dead", NewLilianaWakerOfTheDead)
}

// NewLilianaWakerOfTheDead creates a Liliana Waker Of The Dead
// {2}{B}{B} - PLANESWALKER
func NewLilianaWakerOfTheDead(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Liliana Waker Of The Dead")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"LILIANA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(xValue, xValue)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}