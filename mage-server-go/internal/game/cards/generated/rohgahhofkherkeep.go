package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rohgahh Of Kher Keep", NewRohgahhOfKherKeep)
}

// NewRohgahhOfKherKeep creates a Rohgahh Of Kher Keep
// {2}{B}{B}{R}{R} - CREATURE
func NewRohgahhOfKherKeep(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rohgahh Of Kher Keep")
	card.ManaCost = "{2}{B}{B}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KOBOLD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
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
		AddEffect(abilities.NewBoostEffect(2, 2, filter, false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
