package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jace Ingenious Mind Mage", NewJaceIngeniousMindMage)
}

// NewJaceIngeniousMindMage creates a Jace Ingenious Mind Mage
// {4}{U}{U} - PLANESWALKER
func NewJaceIngeniousMindMage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jace Ingenious Mind Mage")
	card.ManaCost = "{4}{U}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"JACE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainControlTargetEffect(abilities.DurationCustom)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
