package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Windstorm", NewWindstorm)
}

// NewWindstorm creates a Windstorm
// {X}{G} - INSTANT
func NewWindstorm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Windstorm")
	card.ManaCost = "{X}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(GetXValue.instance, StaticFilters.FILTER_CREATURE_...)
	// card.AddAbility(ability0)
	return card, nil
}
