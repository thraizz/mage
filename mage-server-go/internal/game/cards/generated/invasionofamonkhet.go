package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Invasion Of Amonkhet", NewInvasionOfAmonkhet)
}

// NewInvasionOfAmonkhet creates a Invasion Of Amonkhet
// {1}{U}{B} - BATTLE
func NewInvasionOfAmonkhet(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invasion Of Amonkhet")
	card.ManaCost = "{1}{U}{B}"
	card.Types = []string{"BATTLE"}
	card.Subtypes = []string{"SIEGE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: MillCardsEachPlayerEffect()
	// card.AddAbility(ability0)
	return card, nil
}
