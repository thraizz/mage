package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Journey To Nowhere", NewJourneyToNowhere)
}

// NewJourneyToNowhere creates a Journey To Nowhere
// {1}{W} - ENCHANTMENT
func NewJourneyToNowhere(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Journey To Nowhere")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: ExileTargetForSourceEffect()
	// card.AddAbility(ability0)
	// TODO: Implement triggered ability: LeavesBattlefieldTriggeredAbility
	//   - Effect: ReturnFromExileForSourceEffect()
	// card.AddAbility(ability1)
	return card, nil
}
