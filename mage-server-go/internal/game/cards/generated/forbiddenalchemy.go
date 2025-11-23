package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Forbidden Alchemy", NewForbiddenAlchemy)
}

// NewForbiddenAlchemy creates a Forbidden Alchemy
// {2}{U} - INSTANT
func NewForbiddenAlchemy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Forbidden Alchemy")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(4, 1, PutCards.HAND, PutCards.GRAVEYARD)
	// card.AddAbility(ability0)
	return card, nil
}
