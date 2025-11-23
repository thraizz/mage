package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Loot Exuberant Explorer", NewLootExuberantExplorer)
}

// NewLootExuberantExplorer creates a Loot Exuberant Explorer
// {2}{G} - CREATURE
func NewLootExuberantExplorer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Loot Exuberant Explorer")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAST", "NOBLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 6, 1, filter,                 Put...)
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}
