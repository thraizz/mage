package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Needle Storm", NewNeedleStorm)
}

// NewNeedleStorm creates a Needle Storm
// {2}{G} - SORCERY
func NewNeedleStorm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Needle Storm")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(4, StaticFilters.FILTER_CREATURE_FLYING)
	// card.AddAbility(ability0)
	return card, nil
}
