package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ral Izzet Viceroy", NewRalIzzetViceroy)
}

// NewRalIzzetViceroy creates a Ral Izzet Viceroy
// {3}{U}{R} - PLANESWALKER
func NewRalIzzetViceroy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ral Izzet Viceroy")
	card.ManaCost = "{3}{U}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"RAL"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(2, 1, PutCards.HAND, PutCards.GRAVEYARD)
	// card.AddAbility(ability0)
	return card, nil
}
