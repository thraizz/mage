package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Day Of The Dragons", NewDayOfTheDragons)
}

// NewDayOfTheDragons creates a Day Of The Dragons
// {4}{U}{U}{U} - ENCHANTMENT
func NewDayOfTheDragons(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Day Of The Dragons")
	card.ManaCost = "{4}{U}{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LeavesBattlefieldTriggeredAbility
	//   - Effect: SacrificeAllControllerEffect()
	// card.AddAbility(ability0)
	return card, nil
}
