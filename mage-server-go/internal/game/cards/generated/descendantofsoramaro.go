package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Descendant Of Soramaro", NewDescendantOfSoramaro)
}

// NewDescendantOfSoramaro creates a Descendant Of Soramaro
// {3}{U} - CREATURE
func NewDescendantOfSoramaro(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Descendant Of Soramaro")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryControllerEffect(CardsInControllerHandCount.ANY)
	// card.AddAbility(ability0)
	return card, nil
}
