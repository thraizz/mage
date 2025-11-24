package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dack Fayden", NewDackFayden)
}

// NewDackFayden creates a Dack Fayden
// {1}{U}{R} - PLANESWALKER
func NewDackFayden(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dack Fayden")
	card.ManaCost = "{1}{U}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"DACK"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(2)
	// card.AddAbility(ability0)
	return card, nil
}