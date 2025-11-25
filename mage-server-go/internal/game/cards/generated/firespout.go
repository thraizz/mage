package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Firespout", NewFirespout)
}

// NewFirespout creates a Firespout
// {2}{R/G} - SORCERY
func NewFirespout(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Firespout")
	card.ManaCost = "{2}{R/G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(3, filter1)
	//   - DamageAllEffect(3, StaticFilters.FILTER_CREATURE_FLYING)
	// card.AddAbility(ability0)
	return card, nil
}
