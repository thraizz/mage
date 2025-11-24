package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Lady Of Otaria", NewTheLadyOfOtaria)
}

// NewTheLadyOfOtaria creates a The Lady Of Otaria
// {3}{R}{G} - CREATURE
func NewTheLadyOfOtaria(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Lady Of Otaria")
	card.ManaCost = "{3}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RevealLibraryPickControllerEffect(                         4, Integer.MAX_VALUE, fil...)
	// card.AddAbility(ability0)
	return card, nil
}
