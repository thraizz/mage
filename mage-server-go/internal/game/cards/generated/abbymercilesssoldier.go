package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Abby Merciless Soldier", NewAbbyMercilessSoldier)
}

// NewAbbyMercilessSoldier creates a Abby Merciless Soldier
// {1}{R}{G} - CREATURE
func NewAbbyMercilessSoldier(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abby Merciless Soldier")
	card.ManaCost = "{1}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SURVIVOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("CordycepsInfectedToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, ManaSpentToCastCount.instance)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
