package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Invasion Of Ixalan", NewInvasionOfIxalan)
}

// NewInvasionOfIxalan creates a Invasion Of Ixalan
// {1}{G} - BATTLE
func NewInvasionOfIxalan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invasion Of Ixalan")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"BATTLE"}
	card.Subtypes = []string{"SIEGE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 5, 1, StaticFilters.FILTER_CARD_A...)
	// card.AddAbility(ability0)
	return card, nil
}