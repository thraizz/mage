package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Plargg Dean Of Chaos", NewPlarggDeanOfChaos)
}

// NewPlarggDeanOfChaos creates a Plargg Dean Of Chaos
//   - CREATURE
func NewPlarggDeanOfChaos(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Plargg Dean Of Chaos")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(StaticValue.get(1), StaticValue.get(0), tappedFilter, true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(StaticValue.get(0), StaticValue.get(1), untappedFilter, true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
