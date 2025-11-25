package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gates Ablaze", NewGatesAblaze)
}

// NewGatesAblaze creates a Gates Ablaze
// {2}{R} - SORCERY
func NewGatesAblaze(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gates Ablaze")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(xValue, StaticFilters.FILTER_PERMANENT_CREATURE)
	// card.AddAbility(ability0)
	return card, nil
}
