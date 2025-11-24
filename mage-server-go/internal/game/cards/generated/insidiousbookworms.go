package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Insidious Bookworms", NewInsidiousBookworms)
}

// NewInsidiousBookworms creates a Insidious Bookworms
// {B} - CREATURE
func NewInsidiousBookworms(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Insidious Bookworms")
	card.ManaCost = "{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WORM"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1, true)
	//   - DoIfCostPaid(new DiscardTargetEffect(1, true), new ManaCostsImp...)
	// card.AddAbility(ability0)
	return card, nil
}
