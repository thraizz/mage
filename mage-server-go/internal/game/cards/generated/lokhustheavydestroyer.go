package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lokhust Heavy Destroyer", NewLokhustHeavyDestroyer)
}

// NewLokhustHeavyDestroyer creates a Lokhust Heavy Destroyer
// {1}{B}{B}{B} - ARTIFACT CREATURE
// Flying
func NewLokhustHeavyDestroyer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lokhust Heavy Destroyer")
	card.ManaCost = "{1}{B}{B}{B}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"NECRON"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(StaticFilters.FILTER_PERMANENT_CREATURE)
	// card.AddAbility(ability1)
	return card, nil
}
