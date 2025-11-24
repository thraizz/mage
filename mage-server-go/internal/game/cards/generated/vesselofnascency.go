package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vessel Of Nascency", NewVesselOfNascency)
}

// NewVesselOfNascency creates a Vessel Of Nascency
// {G} - ENCHANTMENT
func NewVesselOfNascency(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vessel Of Nascency")
	card.ManaCost = "{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - RevealLibraryPickControllerEffect(4, 1, filter, PutCards.HAND, PutCards.GRAVEYARD)
	//
	// Costs:
	//   - AddSacrificeSourceCost()
	// card.AddAbility(ability0)
	return card, nil
}