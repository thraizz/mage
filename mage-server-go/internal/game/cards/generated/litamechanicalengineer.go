package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Lita Mechanical Engineer", NewLitaMechanicalEngineer)
}

// NewLitaMechanicalEngineer creates a Lita Mechanical Engineer
// {2}{W} - ARTIFACT CREATURE
// Vigilance
func NewLitaMechanicalEngineer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lita Mechanical Engineer")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("ZeppelinToken")
	if err != nil {
		return nil, err
	}
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		Build()
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewUntapEffect(filter)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
