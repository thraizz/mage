package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Drinker Of Sorrow", NewDrinkerOfSorrow)
}

// NewDrinkerOfSorrow creates a Drinker Of Sorrow
// {2}{B} - CREATURE
func NewDrinkerOfSorrow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Drinker Of Sorrow")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HORROR"}
	card.Power = "5"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeControllerEffect(StaticFilters.FILTER_PERMANENT_A, 1, "")
	// card.AddAbility(ability0)
	return card, nil
}