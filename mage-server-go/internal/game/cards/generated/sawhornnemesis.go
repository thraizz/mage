package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sawhorn Nemesis", NewSawhornNemesis)
}

// NewSawhornNemesis creates a Sawhorn Nemesis
//
//	-
func NewSawhornNemesis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sawhorn Nemesis")
	card.ManaCost = ""
	card.Subtypes = []string{"DINOSAUR"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ChoosePlayerEffect(Outcome.Detriment)
	// card.AddAbility(ability0)
	return card, nil
}
