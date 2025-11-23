package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Oreplate Pangolin", NewOreplatePangolin)
}

// NewOreplatePangolin creates a Oreplate Pangolin
// {1}{R} - ARTIFACT CREATURE
func NewOreplatePangolin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Oreplate Pangolin")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT", "PANGOLIN"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                         new AddCountersSourceEffe...)
	// card.AddAbility(ability0)
	return card, nil
}
