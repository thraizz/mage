package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Commander Liara Portyr", NewCommanderLiaraPortyr)
}

// NewCommanderLiaraPortyr creates a Commander Liara Portyr
// {3}{R}{W} - CREATURE
func NewCommanderLiaraPortyr(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Commander Liara Portyr")
	card.ManaCost = "{3}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: AttacksWithCreaturesTriggeredAbility
	//   - Effect: CommanderLiaraPortyrCostEffect()
	// card.AddAbility(ability0)
	return card, nil
}
