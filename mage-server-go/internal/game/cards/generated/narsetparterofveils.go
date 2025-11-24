package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Narset Parter Of Veils", NewNarsetParterOfVeils)
}

// NewNarsetParterOfVeils creates a Narset Parter Of Veils
// {1}{U}{U} - PLANESWALKER
func NewNarsetParterOfVeils(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Narset Parter Of Veils")
	card.ManaCost = "{1}{U}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"NARSET"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(4, 1, filter, PutCards.HAND, PutCards.BOTTOM_RANDO...)
	// card.AddAbility(ability0)
	return card, nil
}