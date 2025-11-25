package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Daretti Rocketeer Engineer", NewDarettiRocketeerEngineer)
}

// NewDarettiRocketeerEngineer creates a Daretti Rocketeer Engineer
// {4}{R} - CREATURE
func NewDarettiRocketeerEngineer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Daretti Rocketeer Engineer")
	card.ManaCost = "{4}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldOrAttacksSourceTriggeredAbility
	//   - Effect: ReturnFromGraveyardToBattlefieldTargetEffect()
	//   - Effect: DoIfCostPaid(                 new ReturnFromGraveyardToBattlefi...)
	// card.AddAbility(ability0)
	return card, nil
}
