package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Bringer Of The Green Dawn", NewBringerOfTheGreenDawn)
}

// NewBringerOfTheGreenDawn creates a Bringer Of The Green Dawn
// {7}{G}{G} - CREATURE
// Trample
func NewBringerOfTheGreenDawn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bringer Of The Green Dawn")
	card.ManaCost = "{7}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BRINGER"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("BeastToken")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewTriggeredAbilityBuilder(card.ID).
		// TODO: Set trigger for LeavesBattlefieldAll (when any permanent you control leaves the battlefield)
		// SetTrigger(abilities.NewLeavesBattlefieldAllTrigger(card.ID, abilities.NewControlledPermanentFilter())).
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		Build()
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	token2_0, err := token.GetToken("BeastToken")
	if err != nil {
		return nil, err
	}
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token2_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
