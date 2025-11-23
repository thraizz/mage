package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Arcum Dagsson", NewArcumDagsson)
}

// NewArcumDagsson creates a Arcum Dagsson
// {3}{U} - CREATURE
func NewArcumDagsson(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arcum Dagsson")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - SacrificeTargetEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}
