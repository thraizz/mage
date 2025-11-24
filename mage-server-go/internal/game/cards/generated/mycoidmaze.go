package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mycoid Maze", NewMycoidMaze)
}

// NewMycoidMaze creates a Mycoid Maze
//   - LAND
func NewMycoidMaze(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mycoid Maze")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"CAVE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 4, 1, StaticFilters.FILTER_CARD_C...)
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
