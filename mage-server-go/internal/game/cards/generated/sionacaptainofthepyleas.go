package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Siona Captain Of The Pyleas", NewSionaCaptainOfThePyleas)
}

// NewSionaCaptainOfThePyleas creates a Siona Captain Of The Pyleas
// {1}{G}{W} - CREATURE
func NewSionaCaptainOfThePyleas(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Siona Captain Of The Pyleas")
	card.ManaCost = "{1}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(7, 1, filter, PutCards.HAND, PutCards.BOTTOM_RANDO...)
	// card.AddAbility(ability0)
	return card, nil
}
