package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Sarkhan Vol", NewSarkhanVol)
}

// NewSarkhanVol creates a Sarkhan Vol
// {2}{R}{G} - PLANESWALKER
func NewSarkhanVol(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sarkhan Vol")
	card.ManaCost = "{2}{R}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"SARKHAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 1)).
		AddEffect(abilities.NewGrantAbilityEffect("HasteAbility", effects.DurationEndOfTurn)).
		AddEffect(abilities.NewGainControlTargetEffect(abilities.DurationEndOfTurn)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("DragonToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0, 5)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
