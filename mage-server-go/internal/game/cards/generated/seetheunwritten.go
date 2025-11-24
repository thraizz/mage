package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("See The Unwritten", NewSeeTheUnwritten)
}

// NewSeeTheUnwritten creates a See The Unwritten
// {4}{G}{G} - SORCERY
func NewSeeTheUnwritten(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "See The Unwritten")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RevealLibraryPickControllerEffect(8, 2, StaticFilters.FILTER_CARD_CREATURE_A, PutCar...)
	//   - RevealLibraryPickControllerEffect(8, 1, StaticFilters.FILTER_CARD_CREATURE_A, PutCar...)
	// card.AddAbility(ability0)
	return card, nil
}
