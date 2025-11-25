package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tropical Storm", NewTropicalStorm)
}

// NewTropicalStorm creates a Tropical Storm
// {X}{G} - SORCERY
func NewTropicalStorm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tropical Storm")
	card.ManaCost = "{X}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(GetXValue.instance, StaticFilters.FILTER_CREATURE_...)
	// card.AddAbility(ability0)
	return card, nil
}
