package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mendicant Core Guidelight", NewMendicantCoreGuidelight)
}

// NewMendicantCoreGuidelight creates a Mendicant Core Guidelight
// {W}{U} - ARTIFACT CREATURE
func NewMendicantCoreGuidelight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mendicant Core Guidelight")
	card.ManaCost = "{W}{U}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(copyEffect, new ManaCostsImpl<>("{1}"))
	// card.AddAbility(ability0)
	return card, nil
}
