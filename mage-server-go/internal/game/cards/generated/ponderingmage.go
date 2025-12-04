package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pondering Mage", NewPonderingMage)
}

// NewPonderingMage creates a Pondering Mage
// {3}{U}{U} - CREATURE
func NewPonderingMage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pondering Mage")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: LookLibraryControllerEffect(3)
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryControllerEffect(3)
	// card.AddAbility(ability1)
	return card, nil
}
