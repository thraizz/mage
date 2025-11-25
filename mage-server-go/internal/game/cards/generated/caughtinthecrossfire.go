package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Caught In The Crossfire", NewCaughtInTheCrossfire)
}

// NewCaughtInTheCrossfire creates a Caught In The Crossfire
// {R}{R} - INSTANT
func NewCaughtInTheCrossfire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Caught In The Crossfire")
	card.ManaCost = "{R}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(2, filter)
	//   - DamageAllEffect(2, filter2)
	// card.AddAbility(ability0)
	return card, nil
}
