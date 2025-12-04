package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Crack The Earth", NewCrackTheEarth)
}

// NewCrackTheEarth creates a Crack The Earth
// {R} - SORCERY
func NewCrackTheEarth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Crack The Earth")
	card.ManaCost = "{R}"
	card.Types = []string{"SORCERY"}
	card.Subtypes = []string{"ARCANE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(1, new FilterControlledPermanent("permanent"))
	// card.AddAbility(ability0)
	return card, nil
}
