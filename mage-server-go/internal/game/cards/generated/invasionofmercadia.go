package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Invasion Of Mercadia", NewInvasionOfMercadia)
}

// NewInvasionOfMercadia creates a Invasion Of Mercadia
// {1}{R} - BATTLE
func NewInvasionOfMercadia(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invasion Of Mercadia")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"BATTLE"}
	card.Subtypes = []string{"SIEGE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new DrawCardSourceControllerEffect(2), new Discard...)
	// card.AddAbility(ability0)
	return card, nil
}
