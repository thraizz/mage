package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Maulfist Doorbuster", NewMaulfistDoorbuster)
}

// NewMaulfistDoorbuster creates a Maulfist Doorbuster
// {3}{R} - CREATURE
func NewMaulfistDoorbuster(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Maulfist Doorbuster")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WARRIOR"}
	card.Power = "4"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new CantBlockTargetEffect(Duration.EndOfTurn), new...)
	// card.AddAbility(ability0)
	return card, nil
}
