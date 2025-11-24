package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gather The Pack", NewGatherThePack)
}

// NewGatherThePack creates a Gather The Pack
// {1}{G} - SORCERY
func NewGatherThePack(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gather The Pack")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RevealLibraryPickControllerEffect(5, 2, StaticFilters.FILTER_CARD_CREATURE_A, PutCar...)
	//   - RevealLibraryPickControllerEffect(5, 1, StaticFilters.FILTER_CARD_CREATURE_A, PutCar...)
	// card.AddAbility(ability0)
	return card, nil
}