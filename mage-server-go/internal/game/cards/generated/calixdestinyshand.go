package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Calix Destinys Hand", NewCalixDestinysHand)
}

// NewCalixDestinysHand creates a Calix Destinys Hand
// {2}{G}{W} - PLANESWALKER
func NewCalixDestinysHand(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Calix Destinys Hand")
	card.ManaCost = "{2}{G}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"CALIX"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 4, 1, filter, PutCards.HAND, PutC...)
	// card.AddAbility(ability0)
	return card, nil
}