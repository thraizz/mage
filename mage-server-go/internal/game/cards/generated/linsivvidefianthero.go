package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lin Sivvi Defiant Hero", NewLinSivviDefiantHero)
}

// NewLinSivviDefiantHero creates a Lin Sivvi Defiant Hero
// {1}{W}{W} - CREATURE
func NewLinSivviDefiantHero(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lin Sivvi Defiant Hero")
	card.ManaCost = "{1}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "REBEL"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - LinSivviDefiantHeroEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - PutOnLibraryTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{3}")
	// card.AddAbility(ability1)
	return card, nil
}
