package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rowans Grim Search", NewRowansGrimSearch)
}

// NewRowansGrimSearch creates a Rowans Grim Search
// {2}{B} - INSTANT
func NewRowansGrimSearch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rowans Grim Search")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(4, 2, PutCards.TOP_ANY, PutCards.GRAVEYARD, true)
	// card.AddAbility(ability0)
	return card, nil
}
