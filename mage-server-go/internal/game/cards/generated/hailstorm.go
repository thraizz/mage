package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hail Storm", NewHailStorm)
}

// NewHailStorm creates a Hail Storm
// {1}{G}{G} - INSTANT
func NewHailStorm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hail Storm")
	card.ManaCost = "{1}{G}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(2, new FilterAttackingCreature())
	//   - DamageAllEffect(1, StaticFilters.FILTER_PERMANENT_CREATURE_CONTROL...)
	// card.AddAbility(ability0)
	return card, nil
}
