package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Skyreaping", NewSkyreaping)
}

// NewSkyreaping creates a Skyreaping
// {1}{G} - SORCERY
func NewSkyreaping(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Skyreaping")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(DevotionCount.G, StaticFilters.FILTER_CREATURE_FLY...)
	// card.AddAbility(ability0)
	return card, nil
}
