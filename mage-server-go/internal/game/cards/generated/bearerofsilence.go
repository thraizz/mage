package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bearer Of Silence", NewBearerOfSilence)
}

// NewBearerOfSilence creates a Bearer Of Silence
// {1}{B} - CREATURE
// Flying
func NewBearerOfSilence(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bearer Of Silence")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: CastSourceTriggeredAbility
	//   - Effect: DoIfCostPaid(new SacrificeEffect(StaticFilters.FILTER_PERMANENT...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewOpponentTargetFilter())
	// card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new SacrificeEffect(StaticFilters.FILTER_PERMANENT...)
	// card.AddAbility(ability2)
	return card, nil
}
