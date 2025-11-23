package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Horn Of The Mark", NewHornOfTheMark)
}

// NewHornOfTheMark creates a Horn Of The Mark
// {2} - ARTIFACT
func NewHornOfTheMark(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Horn Of The Mark")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(5, 1, StaticFilters.FILTER_CARD_CREATURE_A, PutCar...)
	// card.AddAbility(ability0)
	return card, nil
}
