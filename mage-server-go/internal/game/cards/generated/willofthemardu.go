package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Will Of The Mardu", NewWillOfTheMardu)
}

// NewWillOfTheMardu creates a Will Of The Mardu
// {2}{W} - INSTANT
func NewWillOfTheMardu(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Will Of The Mardu")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("RedWarriorToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, WillOfTheMarduValue.instance)).
		AddEffect(abilities.NewDamageEffect(CreaturesYouControlCount.PLURAL)).
		AddTarget(abilities.NewPlayerTargetFilter()).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
