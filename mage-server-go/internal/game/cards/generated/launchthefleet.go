package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Launch The Fleet", NewLaunchTheFleet)
}

// NewLaunchTheFleet creates a Launch The Fleet
// {W} - SORCERY
func NewLaunchTheFleet(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Launch The Fleet")
	card.ManaCost = "{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("SoldierToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 1, true, true)).
		AddEffect(abilities.NewGrantAbilityEffect(new AttacksTriggeredAbility(new CreateTokenEffect(token0_0, 1, true, true), false))).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 1, true, true)).
		AddEffect(abilities.NewGrantAbilityEffect(new AttacksTriggeredAbility(new CreateTokenEffect(token0_0, 1, true, true), false))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}