package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fast Furious", NewFastFurious)
}

// NewFastFurious creates a Fast Furious
// {2}{R} - INSTANT
func NewFastFurious(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fast Furious")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(3, filter)
	//   - DiscardControllerEffect(1)
	// card.AddAbility(ability0)
	return card, nil
}
