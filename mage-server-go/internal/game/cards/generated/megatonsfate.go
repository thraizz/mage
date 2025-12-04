package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Megatons Fate", NewMegatonsFate)
}

// NewMegatonsFate creates a Megatons Fate
// {5}{R} - SORCERY
func NewMegatonsFate(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Megatons Fate")
	card.ManaCost = "{5}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(8, StaticFilters.FILTER_PERMANENT_CREATURE)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewArtifactTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
