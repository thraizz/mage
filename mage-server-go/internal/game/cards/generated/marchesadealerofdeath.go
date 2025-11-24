package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Marchesa Dealer Of Death", NewMarchesaDealerOfDeath)
}

// NewMarchesaDealerOfDeath creates a Marchesa Dealer Of Death
// {U}{B}{R} - CREATURE
func NewMarchesaDealerOfDeath(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Marchesa Dealer Of Death")
	card.ManaCost = "{U}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                         2, 1, PutCards.HAND, PutC...)
	//   - DoIfCostPaid(                 new LookLibraryAndPickControllerE...)
	// card.AddAbility(ability0)
	return card, nil
}
