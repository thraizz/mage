package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Belisarius Cawl", NewBelisariusCawl)
}

// NewBelisariusCawl creates a Belisarius Cawl
// {2}{W}{U} - ARTIFACT CREATURE
func NewBelisariusCawl(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Belisarius Cawl")
	card.ManaCost = "{2}{W}{U}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"HUMAN"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("WhiteAstartesWarriorToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(GetXValue.instance, 1, StaticFilters.FILTER_CARD_A...)
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
