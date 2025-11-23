package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Resonance Technician", NewResonanceTechnician)
}

// NewResonanceTechnician creates a Resonance Technician
// {3}{U/R}{U/R} - CREATURE
// Flying
func NewResonanceTechnician(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Resonance Technician")
	card.ManaCost = "{3}{U/R}{U/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WEIRD", "DETECTIVE"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new InvestigateEffect(2), new DiscardCardCost())
	// card.AddAbility(ability1)
	return card, nil
}
