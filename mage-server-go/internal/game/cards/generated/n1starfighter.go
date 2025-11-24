package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("N1 Starfighter", NewN1Starfighter)
}

// NewN1Starfighter creates a N1 Starfighter
// {1}{W/U}{W/U} - ARTIFACT CREATURE
func NewN1Starfighter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "N1 Starfighter")
	card.ManaCost = "{1}{W/U}{W/U}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"STARSHIP"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new ExileTargetForSourceEffect(), new ManaCostsImp...)
	// card.AddAbility(ability0)
	return card, nil
}