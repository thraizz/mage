package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Crystalline Crawler", NewCrystallineCrawler)
}

// NewCrystallineCrawler creates a Crystalline Crawler
// {4} - ARTIFACT CREATURE
func NewCrystallineCrawler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Crystalline Crawler")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"CONSTRUCT"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1), ColorsOfManaSpentToCastCount.getInstance(), true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
