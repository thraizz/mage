package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Barrage Of Boulders", NewBarrageOfBoulders)
}

// NewBarrageOfBoulders creates a Barrage Of Boulders
// {2}{R} - SORCERY
func NewBarrageOfBoulders(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Barrage Of Boulders")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(1, StaticFilters.FILTER_CREATURE_YOU_DONT_CONTROL)
	// card.AddAbility(ability0)
	return card, nil
}
