package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cloudthresher", NewCloudthresher)
}

// NewCloudthresher creates a Cloudthresher
// {2}{G}{G}{G}{G} - CREATURE
// Flash, Reach
func NewCloudthresher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cloudthresher")
	card.ManaCost = "{2}{G}{G}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(2, "it", StaticFilters.FILTER_CREATURE_FLYING)
	// card.AddAbility(ability2)
	return card, nil
}
